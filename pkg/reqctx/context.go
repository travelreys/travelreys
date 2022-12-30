package reqctx

import "context"

type CallerInfo struct {
	RawToken  string
	UserID    string
	UserEmail string

	Err error
}

type Context struct {
	context.Context
	CallerInfo CallerInfo
}

func NewContextFromContext(ctx context.Context) Context {
	return Context{Context: ctx}
}
