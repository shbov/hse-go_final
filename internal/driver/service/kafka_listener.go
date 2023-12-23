package service

import "context"

type Listener interface {
	Run(ctx context.Context, locationURL string)
}
