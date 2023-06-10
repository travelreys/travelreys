package media

import (
	"context"
	"errors"

	"github.com/travelreys/travelreys/pkg/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	defaultPageSize = 10

	bsonKeyID       = "id"
	bsonKeyUserID   = "userID"
	mediaItemColl   = "media"
	storeLoggerName = "media.store"
)

var (
	ErrUnexpectedStoreError = errors.New("media.store.unexpected-error")
)

type ListMediaFilter struct {
	UserID *string  `json:"userID" bson:"userID"`
	TripID *string  `json:"tripID" bson:"tripID"`
	IDs    []string `json:"ids"`
}

func (ff ListMediaFilter) toBSON() bson.M {
	bsonM := bson.M{}
	if ff.UserID != nil && *ff.UserID != "" {
		bsonM["userID"] = *ff.UserID
	}
	if ff.TripID != nil && *ff.TripID != "" {
		bsonM["tripID"] = *ff.TripID
	}
	if len(ff.IDs) > 0 {
		bsonM["id"] = bson.M{"$in": ff.IDs}
	}
	return bsonM
}

type ListMediaPagination struct {
	Page    *uint64 `json:"page" bson:"page"`
	StartId *string `json:"startId" bson:"startId"`
}

type DeleteMediaFilter struct {
	IDs []string `json:"ids"`
}

func (ff DeleteMediaFilter) toBSON() bson.M {
	bsonM := bson.M{}
	if len(ff.IDs) > 0 {
		bsonM["id"] = bson.M{"$in": ff.IDs}
	}
	return bsonM
}

type Store interface {
	Save(ctx context.Context, items MediaItemList) error
	Delete(ctx context.Context, ff DeleteMediaFilter) error
	List(ctx context.Context, ff ListMediaFilter, pg ListMediaPagination) (MediaItemList, string, error)
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
		{Keys: bson.M{bsonKeyUserID: 1, bsonKeyID: 1}},
	})
	return &store{db, coll, logger.Named(storeLoggerName)}
}

func (str *store) Save(ctx context.Context, items MediaItemList) error {
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
		str.logger.Error("Save", zap.Error(err))
	}
	return err
}

func (str *store) List(ctx context.Context, ff ListMediaFilter, pg ListMediaPagination) (MediaItemList, string, error) {
	list := MediaItemList{}
	bsonFF := ff.toBSON()

	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(defaultPageSize)
	if pg.Page != nil {
		numSkip := *pg.Page * defaultPageSize
		opts.SetSkip(int64(numSkip))
	}
	if pg.StartId != nil {
		objID, err := primitive.ObjectIDFromHex(*pg.StartId)
		if err != nil {
			return nil, "", err
		}
		bsonFF["_id"] = bson.M{"$lt": objID}
	}

	cursor, err := str.coll.Find(ctx, bsonFF, opts)
	if err != nil {
		str.logger.Error("List", zap.String("ff", common.FmtString(ff)), zap.Error(err))
		return list, "", err
	}
	var lastId string
	for cursor.Next(ctx) {
		bsonItem := struct {
			ObjectID  primitive.ObjectID `bson:"_id"`
			MediaItem `bson:"inline"`
		}{}
		if err := cursor.Decode(&bsonItem); err != nil {
			return nil, "", err
		}
		lastId = bsonItem.ObjectID.Hex()
		list = append(list, bsonItem.MediaItem)
	}
	return list, lastId, nil
}

func (str *store) Delete(ctx context.Context, ff DeleteMediaFilter) error {
	_, err := str.coll.DeleteMany(ctx, ff.toBSON())
	if err != nil {
		str.logger.Error("Delete", zap.String("ff", common.FmtString(ff)), zap.Error(err))
		return ErrUnexpectedStoreError
	}
	return nil
}
