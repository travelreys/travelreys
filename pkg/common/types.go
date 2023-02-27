package common

const (
	DefaultChSize = 512
)

type Labels map[string]string
type Tags map[string]string

type GenericJSON map[string]interface{}

func UInt64Ptr(i uint64) *uint64 { return &i }
func Int64Ptr(i int64) *int64    { return &i }
func StringPtr(i string) *string { return &i }
func BoolPtr(i bool) *bool       { return &i }

type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Positioning struct {
	Name    string `json:"name"`
	Address string `json:"address"`

	Continent string `json:"continent"`
	Country   string `json:"country"`
	State     string `json:"state"`
	City      string `json:"city"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`

	Labels Labels `json:"labels"`
}

type FileObject struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Size     float64 `json:"size"`
	URL      string  `json:"url"`
	MIMEType string  `json:"mimetype"`
}