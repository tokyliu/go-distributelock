package locker

type Lock interface {
	Lock(lockName string) bool
	Unlock(lockName string) bool
}





