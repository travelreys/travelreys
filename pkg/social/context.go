package social

import (
	"context"
	"errors"

	"github.com/travelreys/travelreys/pkg/common"
)

var (
	ErrNoInfoSet = errors.New("social.ErrNoInfoSet")
)

type FollowRequestInfo struct {
	Req FollowRequest
}

func ContextWithFollowRequestInfo(ctx context.Context, req FollowRequest) context.Context {
	return context.WithValue(ctx, common.ContextKeyFollowRequestInfo, FollowRequestInfo{Req: req})
}

func FollowRequestInfoFromCtx(ctx context.Context) (FollowRequestInfo, error) {
	val := ctx.Value(common.ContextKeyFollowRequestInfo)
	if val == nil {
		return FollowRequestInfo{}, ErrNoInfoSet
	}
	fi, _ := val.(FollowRequestInfo)
	return fi, nil
}
