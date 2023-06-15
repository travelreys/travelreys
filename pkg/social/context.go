package social

import (
	"context"
	"errors"

	"github.com/travelreys/travelreys/pkg/common"
)

type FriendRequestInfo struct {
	Req FriendRequest
}

func ContextWithFriendRequestInfo(ctx context.Context, req FriendRequest) context.Context {
	return context.WithValue(ctx, common.ContextKeyFriendRequestInfo, FriendRequestInfo{Req: req})
}

func FriendRequestInfoFromCtx(ctx context.Context) (FriendRequestInfo, error) {
	val := ctx.Value(common.ContextKeyFriendRequestInfo)
	if val == nil {
		return FriendRequestInfo{}, errors.New("no info set")
	}
	fi, _ := val.(FriendRequestInfo)
	return fi, nil
}
