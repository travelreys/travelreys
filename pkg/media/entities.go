package media

import (
	"fmt"
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

	Type   string            `json:"type" bson:"type"`
	TripID string            `json:"tripID" bson:"tripID"`
	UserID string            `json:"userID" bson:"userID"`
	URLs   MediaPresignedUrl `json:"urls" bson:"-"`
}

func (item MediaItem) OptimizedPath() string {
	if item.Type == MediaTypePicture {
		return fmt.Sprintf("%s.jpeg", item.Path)
	}
	return item.Path
}

func (item MediaItem) UploadPreviewPath() string {
	return item.Path + "-preview"
}

func (item MediaItem) PreviewPath() string {
	if item.Type == MediaTypePicture {
		return fmt.Sprintf("%s.jpeg", item.Path)
	}
	return fmt.Sprintf("%s-preview.jpeg", item.Path)
}

func (item MediaItem) VideoH264Path() string {
	inFile := item.Path
	return inFile[0:len(inFile)-len(filepath.Ext(inFile))] + ".h264.mp4"
}

func (item MediaItem) VideoH265Path() string {
	inFile := item.Path
	return inFile[0:len(inFile)-len(filepath.Ext(inFile))] + ".h265.mp4"
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
	ContentURL string `json:"contentURL" bson:"-"`

	Image ImagePresignedUrls `json:"image"`
	Video VideoPresignedUrls `json:"video"`
}

type MediaPresignedUrlList []MediaPresignedUrl

type ImagePresignedUrls struct {
	OptimizedURL string `json:"optimizedURL" bson:"-"`
}

type VideoPresignedUrls struct {
	PreviewURL string        `json:"previewURL" bson:"-"`
	Sources    []VideoSource `json:"sources" bson:"-"`
}

type VideoSource struct {
	Source string `json:"source" bson:"-"`
	Codecs string `json:"codecs" bson:"-"`
}
