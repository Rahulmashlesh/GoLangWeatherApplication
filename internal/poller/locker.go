package poller

import "context"

type Locker interface {
	Lock(context.Context) bool // (bool, error)
//	Unlock(context.Context) error
}
