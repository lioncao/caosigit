package redis

import (
	"buddy/util/tools"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type RedisCon struct {
	rs   redis.Conn
	addr string
}

func (this *RedisCon) Init(ip, port string) error {
	var err error
	this.addr = fmt.Sprintf("%s:%s", ip, port)
	this.rs, err = redis.Dial("tcp", this.addr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//选择数据库
func (this *RedisCon) Select(dbid int32) error {
	n, err := this.Do("SELECT", dbid)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(n)
	return nil
}

func (this *RedisCon) Set(key, value string) (string, error) {
	n, err := this.Do("SET", key, value)
	if err != nil {
		return "", err
	}
	return n.(string), err
}

func (this *RedisCon) Get(key string) ([]byte, error) {
	n, err := this.Do("GET", key)
	if err != nil {
		return n.([]byte), err
	} else {
		if n == nil {
			return []byte{}, fmt.Errorf("cont't find")
		}
	}
	return n.([]byte), err
}

func (this *RedisCon) Close() {
	this.rs.Close()
}

//有序列表增加
func (this *RedisCon) ZAdd(key string, score int32, member string) (int32, error) {
	n, err := this.Do("ZADD", key, score, member)
	return int32(n.(int64)), err
}

//有序列表删除
func (this *RedisCon) ZRem(key string, member string) (int32, error) {
	n, err := this.Do("ZREM", key, member)
	return int32(n.(int64)), err
}

//有序列表限制数量增加
func (this *RedisCon) ZAddLimit(key string, score int32, member string, limit int32) (int32, error) {
	var (
		n   interface{}
		num interface{}
		err error
	)
	num, err = this.Do("ZADD", key, score, member)
	if err != nil {
		return 0, err
	}
	n, err = this.Do("ZCARD", key)
	if err != nil {
		return 1, err
	}
	if int32(n.(int64)) > limit {
		n, err = this.Do("zremrangebyrank", key, limit, n)
		return 1, err
	}
	return int32(num.(int64)), err
}

//有序列表取得排名
func (this *RedisCon) ZRevrank(key string, member string) (int32, error) {
	n, err := this.Do("ZREVRANK", key, member)
	if err != nil {
		return 0, err
	}
	if n == nil {
		return -1, err
	}
	return int32(n.(int64)), err
}

func (this *RedisCon) Zrange(key string, startRank, endRank int32) ([]interface{}, error) {
	n, err := this.rs.Do("zrange", key, startRank, endRank)
	ar := n.([]interface{})
	return ar, err
}

//hashset
func (this *RedisCon) Hset(key, field, value string) (int32, error) {
	n, err := this.Do("HSET", key, field, value)
	return int32(n.(int64)), err
}

//hashget
func (this *RedisCon) Hget(key, field string) ([]byte, error) {
	n, err := this.Do("HGET", key, field)
	if err != nil {
		return []byte{}, err
	}
	if n == nil {
		return []byte{}, fmt.Errorf("can't find")
	}
	fmt.Println(string(n.([]byte)))
	return n.([]byte), err
}

func (this *RedisCon) HDel(key, field string) (int32, error) {
	n, err := this.Do("HDEL", key, field)
	return int32(n.(int64)), err
}

func (this *RedisCon) Do(cmd string, args ...interface{}) (interface{}, error) {
	n, err := this.rs.Do(cmd, args...)
	if err != nil {
		if err.Error() == "can't find" {
			return n, err
		}
		//重连
		tools.GetLog().LogError("do redis cmd:%s err:%s start reconnect", cmd, err)
		this.rs, err = redis.Dial("tcp", this.addr)
		if err != nil {
			tools.GetLog().LogError("reconnect redis failed %s", this.addr)
			return n, err
		} else {
			tools.GetLog().LogError("reconnect redis success %s", this.addr)
			return this.rs.Do(cmd, args...)
		}
	}
	return n, err
}
