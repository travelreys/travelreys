package flights

import "context"

type AirlinesStore interface {
	Get(context.Context, []string) AirlinesList
}
