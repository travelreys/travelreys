package moodboard

import (
	"github.com/google/uuid"
	"github.com/otiai10/opengraph/v2"
	"github.com/travelreys/travelreys/pkg/common"
)

const (
	TypePinBase = "ogp"
)

type Moodboard struct {
	ID    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
	Pins  PinMap `json:"pins" bson:"pins"`

	Labels common.Labels `json:"labels" bson:"labels"`
	Tags   common.Tags   `json:"tags" bson:"tags"`
}

func NewMoodboard(id string) Moodboard {
	return Moodboard{
		ID:     id,
		Pins:   PinMap{},
		Title:  "",
		Labels: common.Labels{},
		Tags:   common.Tags{},
	}
}

func (mb *Moodboard) AddPin(pin Pin) {
	mb.Pins[pin.ID] = pin
}

// Pin represents a content embed to the moodboard (e.g Notion embeds)
type Pin struct {
	ID    string `json:"id" bson:"id"`
	Type  string `json:"type" bson:"type"`
	Notes string `json:"notes" bson:"notes"`

	OGP opengraph.OpenGraph `json:"ogp"`

	Labels common.Labels `json:"labels" bson:"labels"`
	Tags   common.Tags   `json:"tags" bson:"tags"`
}

type PinMap map[string]Pin

func PinFromOGP(ogp *opengraph.OpenGraph) Pin {
	return Pin{
		ID:     uuid.New().String(),
		Type:   TypePinBase,
		Notes:  "",
		OGP:    *ogp,
		Labels: common.Labels{},
		Tags:   common.Tags{},
	}
}
