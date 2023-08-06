package trips

import "context"

type AsyncService interface {
	SendJoinTripEmail(ctx context.Context, impersonatee, email string)
	JoinTripFromEmail(ctx context.Context, impersonatee, email string)

	SendJoinTripMsg(ctx context.Context, impersonatee, targetUserID string)
	JoinTripFromMsg(ctx context.Context, impersonatee, targetUserID string)
}

type asyncService struct{}

// need to check if user has an account or not on frontend
// if have, once click

func (svc *asyncService) SendJoinTripEmail(ctx context.Context, impersonatee, email string) {

}
