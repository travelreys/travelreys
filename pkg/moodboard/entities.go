package moodboard

import "github.com/travelreys/travelreys/pkg/common"

const (
	TypePinBase = "ogp"
)

type Moodboard struct {
	ID   string `json:"id" bson:"id"`
	Pins PinMap `json:"pins" bson:"pins"`
}

// Pin represents a content embed to the moodboard (e.g Notion embeds)
type Pin struct {
	ID    string `json:"id" bson:"id"`
	Type  string `json:"type" bson:"type"`
	Notes string `json:"notes" bson:"notes"`

	Labels common.Labels `json:"labels" bson:"labels"`
	Tags   common.Tags   `json:"tags" bson:"tags"`
}

type PinMap map[string]Pin
