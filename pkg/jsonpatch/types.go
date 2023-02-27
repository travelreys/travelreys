package jsonpatch

type Op struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func MakeAddOp(path string, val interface{}) Op {
	return Op{"add", path, val}
}

func MakeRemoveOp(path string, val interface{}) Op {
	return Op{"remove", path, val}
}

func MakeRepOp(path string, val interface{}) Op {
	return Op{"replace", path, val}
}
