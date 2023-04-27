package media

import (
	"context"
	"errors"
	"fmt"

	"github.com/travelreys/travelreys/pkg/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	bsonKeyID       = "id"
	bsonKeyUserID   = "userID"
	mediaItemColl   = "media"
	storeLoggerName = "media.store"
)

var (
	ErrUnexpectedStoreError = errors.New("media.store.unexpected-error")
)

type ListMediaFilter struct {
	UserID *string `json:"userID" bson:"userID"`
}

type DeleteMediaFilter struct {
	UserID *string  `json:"userID"`
	IDs    []string `json:"ids"`
}

func (ff DeleteMediaFilter) toBSON() bson.M {
	bsonM := bson.M{}
	if ff.UserID != nil || *ff.UserID != "" {
		bsonM["userID"] = ff.UserID
	}
	if len(ff.IDs) > 0 {
		bsonM["ids"] = bson.M{"$in": ff.IDs}
	}
	return bsonM
}

type Store interface {
	SaveForUser(ctx context.Context, userID string, items MediaItemList) error
	List(ctx context.Context, ff ListMediaFilter) (MediaItemList, error)
	Delete(ctx context.Context, ff DeleteMediaFilter) error
}

type store struct {
	db   *mongo.Database
	coll *mongo.Collection

	logger *zap.Logger
}

func NewStore(ctx context.Context, db *mongo.Database, logger *zap.Logger) Store {
	ctx, cancel := context.WithTimeout(ctx, common.DbReqTimeout)
	defer cancel()

	coll := db.Collection(mediaItemColl)
	coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{bsonKeyID: 1, bsonKeyUserID: 1}},
	})
	return &store{db, coll, logger.Named(storeLoggerName)}
}

func (str *store) SaveForUser(ctx context.Context, userID string, items MediaItemList) error {
	writeModels := make([]mongo.WriteModel, len(items))
	for idx, i := range items {
		replaceFF := bson.M{bsonKeyID: i.ID}
		writeModels[idx] = mongo.NewReplaceOneModel().
			SetFilter(replaceFF).
			SetUpsert(true).
			SetReplacement(i)
	}
	_, err := str.coll.BulkWrite(ctx, writeModels, options.BulkWrite())
	if err != nil {
		str.logger.Error("Save", zap.String("userID", userID), zap.Error(err))
	}
	return err
}

func (str *store) List(ctx context.Context, ff ListMediaFilter) (MediaItemList, error) {
	fmt.Println(ff)
	list := MediaItemList{}
	cursor, err := str.coll.Find(ctx, ff)
	if err != nil {
		str.logger.Error("List", zap.String("ff", common.FmtString(ff)), zap.Error(err))
		return list, err
	}
	err = cursor.All(ctx, &list)
	return list, err
}

func (str *store) Delete(ctx context.Context, ff DeleteMediaFilter) error {
	_, err := str.coll.DeleteMany(ctx, ff.toBSON())
	if err != nil {
		str.logger.Error("Delete", zap.String("ff", common.FmtString(ff)), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}
