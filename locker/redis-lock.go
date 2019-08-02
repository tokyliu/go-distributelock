package locker

//redis distribute lock support features:
//1. expire time support
//2. non-blocking
//3. support  Reenter

import (
	"github.com/gomodule/redigo/redis"
	"fmt"
	"time"
	"strings"
	"strconv"
)

type RedisDistributeLock struct {
	concurrentId string
	timeoutMilSeconds int64
	redisCon redis.Conn
}


func NewRedisDstLock(concurrentId string, timeoutMilSec int64, redisConn redis.Conn) RedisDistributeLock {
	return RedisDistributeLock{
		concurrentId: concurrentId,
		timeoutMilSeconds: timeoutMilSec,
		redisCon: redisConn,
	}
}


func (l RedisDistributeLock)Lock(lockName string) bool{
	reply, err := l.redisCon.Do("EXISTS", lockName)
	if err != nil {
		dealError(err)
		return false
	}
	//lock is unlock, lock it
	if reply == nil {
		v := fmt.Sprintf("%s:%d", l.concurrentId, time.Now().UnixNano()/1e6+l.timeoutMilSeconds)
		reply, err = l.redisCon.Do("GETSET", lockName, v)
		if err != nil {
			dealError(err)
			return false
		}
		if reply == nil {
			l.redisCon.Do("PEXPIRE", lockName, l.timeoutMilSeconds)
			return true
		}else{
			l.redisCon.Do("PSETEX", lockName, l.timeoutMilSeconds, reply)
			return false
		}
	}

	nowMilSeconds := time.Now().UnixNano()/1e6
	tmpReplyArr := strings.Split(string(reply.([]byte)), ":")
	if len(tmpReplyArr) != 2 {
		return false
	}
	//lock hold by self, reentrant lock success
	if tmpReplyArr[0] == l.concurrentId {
		return true
	}

	expireTimestamp,_ := strconv.ParseInt(tmpReplyArr[1], 10, 64)
	if nowMilSeconds < expireTimestamp {
		return false
	}

	value := fmt.Sprintf("%s:%d", l.concurrentId, nowMilSeconds+l.timeoutMilSeconds)
	reply1, err1 := l.redisCon.Do("GETSET", lockName, value)
	if err1 != nil {
		dealError(err1)
		return false
	}
	if reply1 == reply {
		l.redisCon.Do("PEXPIRE", lockName, l.timeoutMilSeconds)
		return true
	}else{
		l.redisCon.Do("PSETEX", lockName, l.timeoutMilSeconds, reply1)
		return false
	}
}


func (l RedisDistributeLock)Unlock(lockName string) bool {
	reply, err := l.redisCon.Do("EXISTS", lockName)
	if err != nil {
		dealError(err)
		return false
	}
	if reply == nil {
		return true
	}
	tmpReplyArr := strings.Split(string(reply.([]byte)), ":")
	if len(tmpReplyArr) != 2 || tmpReplyArr[0] != l.concurrentId{
		return false
	}

	reply, err = l.redisCon.Do("DEL", lockName)
	if err != nil {
		dealError(err)
		return false
	}
	return true
}







