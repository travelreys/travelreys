package ogp

import "github.com/otiai10/opengraph/v2"

/**
 * Structured Properties.
 * https://ogp.me/#structured
 */

// Image represents a structure of "og:image".
// "og:image" might have following properties:
//   - og:image:url
//   - og:image:secure_url
//   - og:image:type
//   - og:image:width
//   - og:image:height
//   - og:image:alt
type Image struct {
	URL       string `json:"url" bson:"url" msgpack:"url"`
	SecureURL string `json:"secure_url" bson:"secure_url" msgpack:"secure_url"`
	Type      string `json:"type" bson:"type" msgpack:"type"`
	Width     int    `json:"width" bson:"width" msgpack:"width"`
	Height    int    `json:"height" bson:"height" msgpack:"height"`
	Alt       string `json:"alt" bson:"alt" msgpack:"alt"`
}

// Video represents a structure of "og:video".
// "og:video" might have following properties:
//   - og:video:url
//   - og:video:secure_url
//   - og:video:type
//   - og:video:width
//   - og:video:height
type Video struct {
	URL       string `json:"url" bson:"url" msgpack:"url"`
	SecureURL string `json:"secure_url" bson:"secure_url" msgpack:"secure_url"`
	Type      string `json:"type" bson:"type" msgpack:"type"` // Content-Type
	Width     int    `json:"width" bson:"width" msgpack:"width"`
	Height    int    `json:"height" bson:"height" msgpack:"height"`
	Duration  int    `json:"duration" bson:"duration" msgpack:"duration"`
}

// Audio represents a structure of "og:audio".
// "og:audio" might have following properties:
//   - og:audio:url
//   - og:audio:secure_url
//   - og:audio:type
type Audio struct {
	URL       string `json:"url" bson:"url" msgpack:"url"`
	SecureURL string `json:"secure_url" bson:"secure_url" msgpack:"secure_url"`
	Type      string `json:"type" bson:"type"  msgpack:"type"` // Content-Type
}

// Favicon represents an extra structure for "shortcut icon".
type Favicon struct {
	URL string `json:"url" bson:"url" msgpack:"url"`
}

type Opengraph struct {
	// Basic Metadata
	// https://ogp.me/#metadata
	Title string  `json:"title"`
	Type  string  `json:"type"`
	Image []Image `json:"image"`
	URL   string  `json:"url"`

	// Optional Metadata
	// https://ogp.me/#optional
	Audio       []Audio  `json:"audio"`
	Description string   `json:"description"`
	Determiner  string   `json:"determiner"`
	Locale      string   `json:"locale"`
	LocaleAlt   []string `json:"locale_alternate"`
	SiteName    string   `json:"site_name"`
	Video       []Video  `json:"video"`

	// Additional (unofficial)
	Favicon Favicon `json:"favicon"`
}

func OpengraphFromRawGraph(graph *opengraph.OpenGraph) Opengraph {
	result := Opengraph{
		Title:       graph.Title,
		Type:        graph.Type,
		Image:       []Image{},
		URL:         graph.URL,
		Audio:       []Audio{},
		Description: graph.Description,
		Determiner:  graph.Determiner,
		Locale:      graph.Locale,
		LocaleAlt:   graph.LocaleAlt,
		SiteName:    graph.SiteName,
		Video:       []Video{},
		Favicon:     Favicon(graph.Favicon),
	}

	for _, item := range graph.Image {
		result.Image = append(result.Image, Image{
			URL:       item.URL,
			SecureURL: item.SecureURL,
			Type:      item.Type,
			Width:     item.Width,
			Height:    item.Height,
			Alt:       item.Alt,
		})
	}
	for _, item := range graph.Audio {
		result.Audio = append(result.Audio, Audio{
			URL:       item.URL,
			SecureURL: item.SecureURL,
			Type:      item.Type,
		})
	}
	for _, item := range graph.Video {
		result.Video = append(result.Video, Video{
			URL:       item.URL,
			SecureURL: item.SecureURL,
			Type:      item.Type,
			Width:     item.Width,
			Height:    item.Height,
			Duration:  item.Duration,
		})
	}

	return result
}
