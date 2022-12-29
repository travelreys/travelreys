package common

type FileObject struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Size     float64 `json:"size"`
	URL      string  `json:"url"`
	MIMEType string  `json:"mimetype"`
}
