package jsonpatch

type Op struct {
	Op    string      `json:"op"  msgpack:"op"`
	Path  string      `json:"path"  msgpack:"path"`
	Value interface{} `json:"value"  msgpack:"value"`
	From  string      `json:"from"  msgpack:"from"`
}

func MakeAddOp(path string, val interface{}) Op {
	return Op{"add", path, val, ""}
}

func MakeRemoveOp(path string, val interface{}) Op {
	return Op{"remove", path, val, ""}
}

func MakeRepOp(path string, val interface{}) Op {
	return Op{"replace", path, val, ""}
}
