package apns

import "context"

type Service interface {
	Push(ctx context.Context, udid string) (string, error)
}
