package media

import (
	"os"
	"path/filepath"
	"time"

	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/storage"
)

const (
	MediaTypePicture = "p"
	MediaTypeVideo   = "v"

	UserMediaPathPrefix = "users"
	LabelMediaURL       = "mediaURL"
	LabelWidth          = "width"
	LabelHeight         = "height"
)

var (
	MediaItemBucket = os.Getenv("TRAVELREYS_MEDIA_BUCKET")
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

	UserID string `json:"userID" bson:"userID"`
	Type   string `json:"type" bson:"type"`
}

type MediaItemList []MediaItem
type MediaItemMap map[string]MediaItem

func NewMediaItem(userID string, param NewMediaItemParams) MediaItem {
	objectPath := filepath.Join(UserMediaPathPrefix, userID, param.Hash)

	return MediaItem{
		Object: storage.Object{
			ID:           param.Hash,
			Name:         param.Name,
			Bucket:       MediaItemBucket,
			Path:         objectPath,
			MIMEType:     param.MimeType,
			LastModified: time.Now(),
			Labels:       common.Labels{},
			Tags:         common.Tags{},
		},
		UserID: userID,
		Type:   param.Type,
	}
}
