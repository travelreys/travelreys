package common

type Labels map[string]string
type Tags map[string]string

type GenericJSON map[string]interface{}

type JSONPatchOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func UInt64Ptr(i uint64) *uint64 { return &i }
func Int64Ptr(i int64) *int64    { return &i }
func StringPtr(i string) *string { return &i }
func BoolPtr(i bool) *bool       { return &i }
