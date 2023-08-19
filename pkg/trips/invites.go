package trips

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/email"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Invite struct {
	ID       string `json:"id"`
	AuthorID string `json:"authorID"`
	TripID   string `json:"tripID"`
	UserID   string `json:"userID"`

	Labels common.Labels `json:"labels"`
}

func NewInvite(tripID, authorID, userID string) Invite {
	return Invite{
		ID:       uuid.NewString(),
		AuthorID: authorID,
		TripID:   tripID,
		UserID:   userID,
	}
}

type InviteList []Invite

const (
	SyncMsgWaitInterval = 500 * time.Millisecond
)

type InviteService interface {
	Send(ctx context.Context, tripID, authorID, userID string) error
	Accept(ctx context.Context, ID string) error
	Decline(ctx context.Context, ID string) error
	Read(ctx context.Context, ID string) (Invite, error)
	List(ctx context.Context, ff ListInvitesFilter) (InviteList, error)
}

type inviteService struct {
	syncSvc  SyncService
	emailSvc email.Service
	store    InviteStore
	logger   *zap.Logger
}

func NewInviteService(
	syncSvc SyncService,
	emailSvc email.Service,
	store InviteStore,
	logger *zap.Logger,
) InviteService {
	return &inviteService{
		syncSvc,
		emailSvc,
		store,
		logger,
	}
}

func (svc *inviteService) Send(
	ctx context.Context,
	tripID,
	authorID,
	userID string,
) error {
	invite := NewInvite(tripID, authorID, userID)
	err := svc.store.Save(ctx, invite)
	if err != nil {
		// send email

	}
	return err
}

func (svc *inviteService) Accept(ctx context.Context, ID string) error {
	invite, err := svc.store.Read(ctx, ID)
	if err != nil {
		return err
	}

	connID := uuid.NewString()
	joinMsg := MakeSyncMsgTOBTopicJoin(
		connID,
		invite.TripID,
		invite.AuthorID,
	)
	if err := svc.syncSvc.Join(ctx, &joinMsg); err != nil {
		return err
	}
	time.Sleep(SyncMsgWaitInterval)

	member := NewMember(invite.UserID, MemberRoleCollaborator)
	addMemMsg := MakeSyncMsgTOBTopicUpdate(
		connID,
		invite.TripID,
		invite.AuthorID,
		SyncMsgTOBUpdateOpUpdateTripMembers,
		MakeSyncMsgTOBUpdateOpUpdateTripMembersOps(member),
	)
	if err := svc.syncSvc.Update(ctx, &addMemMsg); err != nil {
		return err
	}

	return svc.store.Delete(ctx, ID)
}

func (svc *inviteService) Decline(ctx context.Context, ID string) error {
	return svc.store.Delete(ctx, ID)
}

func (svc *inviteService) Read(ctx context.Context, ID string) (Invite, error) {
	return svc.store.Read(ctx, ID)
}

func (svc *inviteService) List(ctx context.Context, ff ListInvitesFilter) (InviteList, error) {
	return svc.store.List(ctx, ff)
}

var (
	ErrInviteNotFound = errors.New("trips.ErrInviteNotFound")
)

type InviteStore interface {
	List(ctx context.Context, ff ListInvitesFilter) (InviteList, error)
	Read(ctx context.Context, ID string) (Invite, error)
	Save(ctx context.Context, invite Invite) error
	Delete(ctx context.Context, ID string) error
}

type inviteStore struct {
	db   *mongo.Database
	coll *mongo.Collection

	logger *zap.Logger
}

func NewInviteStore(
	ctx context.Context,
	db *mongo.Database,
	logger *zap.Logger,
) InviteStore {
	coll := db.Collection("trip_invites")
	coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{bsonKeyID: 1}},
		{Keys: bson.M{"tripID": 1}},
		{Keys: bson.M{"userID": 1}},
	})
	return &inviteStore{db, coll, logger.Named("trips.inviteStore")}
}

func (s *inviteStore) List(
	ctx context.Context,
	ff ListInvitesFilter,
) (InviteList, error) {
	list := InviteList{}
	bsonM, _ := ff.toBSON()
	cursor, err := s.coll.Find(ctx, bsonM)
	if err != nil {
		s.logger.Error("ListInvitesForTrip", zap.Error(err))
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (s *inviteStore) Read(ctx context.Context, ID string) (Invite, error) {
	var invite Invite
	err := s.coll.FindOne(ctx, bson.M{bsonKeyID: ID}).Decode(&invite)
	if err == mongo.ErrNoDocuments {
		return Invite{}, ErrInviteNotFound
	}
	if err != nil {
		s.logger.Error("Read", zap.String("id", ID), zap.Error(err))
		return Invite{}, ErrUnexpectedStoreError
	}
	return invite, err
}

func (s *inviteStore) Save(ctx context.Context, invite Invite) error {
	saveFF := bson.M{bsonKeyID: invite.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := s.coll.ReplaceOne(ctx, saveFF, invite, opts)
	if err != nil {
		s.logger.Error("Save", zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}

func (s *inviteStore) Delete(ctx context.Context, ID string) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{ID: ID})
	if err != nil {
		s.logger.Error("Delete", zap.String("id", ID), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return err
}
