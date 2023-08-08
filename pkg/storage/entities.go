package storage

import (
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/minio/minio-go/v7"
	"github.com/travelreys/travelreys/pkg/common"
)

type Object struct {
	ID           string        `json:"id" bson:"id" msgpack:"id"`
	Name         string        `json:"name" bson:"name" msgpack:"name"`
	Bucket       string        `json:"bucket" bson:"bucket" msgpack:"bucket"`
	Size         int64         `json:"size" bson:"size" msgpack:"size"`
	Path         string        `json:"path" bson:"path" msgpack:"path"`
	MIMEType     string        `json:"mimetype" bson:"mimetype" msgpack:"mimetype"`
	LastModified time.Time     `json:"lastModified" bson:"lastModified" msgpack:"lastModified"`
	Labels       common.Labels `json:"labels" bson:"labels" msgpack:"labels"`
	Tags         common.Tags   `json:"tags" bson:"tags" msgpack:"tags"`
}

func ObjectFromObjectInfo(info minio.ObjectInfo, bucket string) Object {
	keyTkns := strings.Split(info.Key, "/")
	return Object{
		ID:           keyTkns[len(keyTkns)-1],
		Name:         keyTkns[len(keyTkns)-1],
		Bucket:       bucket,
		Path:         info.Key,
		Size:         info.Size,
		MIMEType:     info.ContentType,
		LastModified: info.LastModified,
		Labels:       info.UserTags,
		Tags:         common.Tags{},
	}
}

func ObjectFromAttrs(attrs *storage.ObjectAttrs) Object {
	keyTkns := strings.Split(attrs.Name, "/")
	return Object{
		ID:           keyTkns[len(keyTkns)-1],
		Name:         keyTkns[len(keyTkns)-1],
		Bucket:       attrs.Bucket,
		Size:         attrs.Size,
		Path:         attrs.Name,
		MIMEType:     attrs.ContentType,
		LastModified: attrs.Updated,
		Labels:       attrs.Metadata,
		Tags:         common.Tags{},
	}
}

type ObjectList []Object
