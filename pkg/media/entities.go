package media

import (
	"os"
	"time"

	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/storage"
)

const (
	MediaTypePicture = "p"
	MediaTypeVideo   = "v"

	UserMediaPathPrefix = "users"
	LabelMediaURL       = "mediaURL"
	LabelPreviewURL     = "previewURL"
	LabelOptimizedURL   = "optimizedURL"
	LabelWidth          = "width"
	LabelHeight         = "height"
)

var (
	MediaItemBucket          = os.Getenv("TRAVELREYS_MEDIA_BUCKET")
	MediaItemOptimizedBucket = os.Getenv("TRAVELREYS_MEDIA_OPTIMIZED_BUCKET")
)

type NewMediaItemParams struct {
	Type     string `json:"type"`
	Hash     string `json:"hash"`
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
}

type MediaItem struct {
	storage.Object `bson:"inline"`

	Type   string `json:"type" bson:"type"`
	TripID string `json:"tripID" bson:"tripID"`
	UserID string `json:"userID" bson:"userID"`
}

func (item MediaItem) PreviewPath() string {
	return item.Path + "-preview"
}

type MediaItemList []MediaItem
type MediaItemMap map[string]MediaItem

func NewMediaItem(tripID, userID, objectPath string, param NewMediaItemParams) MediaItem {
	return MediaItem{
		Object: storage.Object{
			ID:           param.Hash,
			Name:         param.Name,
			Bucket:       MediaItemBucket,
			Path:         objectPath,
			Size:         param.Size,
			MIMEType:     param.MimeType,
			LastModified: time.Now(),
			Labels:       common.Labels{},
			Tags:         common.Tags{},
		},
		UserID: userID,
		Type:   param.Type,
	}
}

type MediaPresignedUrl struct {
	ContentURL   string `json:"contentURL"`
	PreviewURL   string `json:"previewURL"`
	OptimizedURL string `json:"optimizedURL"`
}

type MediaPresignedUrlList []MediaPresignedUrl
