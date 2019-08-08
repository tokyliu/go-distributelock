package locker

import (
	"github.com/samuel/go-zookeeper/zk"
)

type ZkDistributeLock struct {
	concurrentId string
	block bool
	timeoutMilSeconds int64
	callback func(args ...interface{})
	zkConn *zk.Conn
}


func NewZkDistributeLock(concurrentId string, block bool, toMilSecs int64, cbfunc func(args ...interface{}), conn *zk.Conn) ZkDistributeLock {
	return ZkDistributeLock{
		concurrentId: concurrentId,
		block: block,
		timeoutMilSeconds: toMilSecs,
		callback: cbfunc,
		zkConn: conn,
	}
}


func (l ZkDistributeLock)Lock(lockName string) bool{

}


func (l ZkDistributeLock)Unlock(lockName string) {

	
}


