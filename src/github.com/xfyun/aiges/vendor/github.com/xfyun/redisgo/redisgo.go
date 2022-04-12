/*
 *基于garyburd的redisgo进行简单的封装
 *参阅redisgo的godoc
 */
package redisgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var(
	default_idleTimeout = 3600 //默认空闲超时时间（单位：s）
)

//配置选项
type RedisCfgOption struct {
	RedisHost      string        //redis主机(必传)
	RedisPasswd    string        //redis密码
	MaxActive      int           //redis最大连接数(maxActive=0代表没连接限制)
	MaxIdle        int           //redis最大空闲实例数
	Db             int           //redis数据库
	IdleTimeOut    time.Duration //redis空闲实例超时设置(单位:s)
	UseTwemproxy   bool          //使用twem代理
}

type RedisCfgOpt func(*RedisCfgOption)

//配置代理使用选项
func WithUseTwemproxy(use bool) RedisCfgOpt {
	return func(rc *RedisCfgOption) {
		rc.UseTwemproxy = use
	}
}

func WithRedisHost(host string) RedisCfgOpt {
	return func(rc *RedisCfgOption) {
		rc.RedisHost = host
	}
}

func WithRedisPwd(pwd string) RedisCfgOpt {
	return func(rc *RedisCfgOption) {
		rc.RedisPasswd = pwd
	}
}

func WithMaxactive(maxactive int) RedisCfgOpt {
	return func(rc *RedisCfgOption) {
		rc.MaxActive = maxactive
	}
}

func WithMaxIdle(maxidle int) RedisCfgOpt {
	return func(rc *RedisCfgOption) {
		rc.MaxIdle =  maxidle
	}
}

func WithDb(db int) RedisCfgOpt {
	return func(rc *RedisCfgOption) {
		rc.Db =  db
	}
}

func WithIdleTimeout(idletimeout time.Duration) RedisCfgOpt {
	return func(rc *RedisCfgOption) {
		rc.IdleTimeOut =  idletimeout
	}
}

type Redis struct {
	Pool *redis.Pool
}

//创建一个redis实例
func NewRedisInst(o ...RedisCfgOpt) (redisInst *Redis, err error) {
	var rcOpt = RedisCfgOption{
		IdleTimeOut:time.Second*time.Duration(default_idleTimeout),
	}
		
	for _, opt := range o {
		opt(&rcOpt)
	}

	if rcOpt.RedisHost==""{
		err =errors.New("redis host is null")
		return
	}
	
	
	Pool := &redis.Pool{
		MaxActive:   rcOpt.MaxActive,   //设置的最大连接数
		MaxIdle:     rcOpt.MaxIdle,     // 最大的空闲连接数，即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
		IdleTimeout: rcOpt.IdleTimeOut, //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭

		//建立连接
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", rcOpt.RedisHost)
			if err != nil {
				return nil, err
			}

			//password为空，将不进行权限验证
			if rcOpt.RedisPasswd != "" {
				if _, err := c.Do("AUTH", rcOpt.RedisPasswd); err != nil {
					c.Close()
					return nil, err
				}
			}

			//默认db使用0
			if !rcOpt.UseTwemproxy{
				if _, err := c.Do("SELECT", rcOpt.Db); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if !rcOpt.UseTwemproxy{
				_, err := c.Do("PING")
				return err
			}
			return nil
		},
	}

	redisInst = &Redis{
		Pool: Pool,
	}

	return redisInst, err
}

//Do(commandName string, args ...interface{}) (reply interface{}, err error)
func (r *Redis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}

func (r *Redis) Send(commandName string, args ...interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()
	return conn.Send(commandName, args...)
}

func (r *Redis) Flush() error {
	conn := r.Pool.Get()
	defer conn.Close()
	return conn.Flush()
}

/*--------------------------------------------key操作-----------------------------------*/

//EXISTS key
//检查键值是否存在
func (r *Redis) Exists(key string) (bool, error) {
	return redis.Bool(r.Do("EXISTS", key))
}

//DEL key
//当key值存在时删除key
func (r *Redis) Del(key string) (err error) {
	_, err = r.Do("DEL", key)
	return
}

//EXPIRE key seconds
// Expire设置键过期时间,expire的单位为秒
func (r *Redis) Expire(key string, expire int) (err error) {
	_, err = r.Do("EXPIRE", key, expire)
	return
}

//PEXPIRE key milliseconds
//对一个key设置过期时间，单位是毫秒
func (r *Redis) PExpire(key string, expired int) error {
	_, err := r.Do("PEXPIRE", key, expired)
	return err
}

//PERSIST key
//移除key的过期时间,key保持永久
func (r *Redis) Persist(key string) (err error) {
	_, err = r.Do("PERSIST", key)
	return
}

//TTL key
// TTL以秒为单位。当 key 不存在时，返回-2.当 key存在但没有设置剩余生存时间时，返回-1
func (r *Redis) Ttl(key string) (ttl int64, err error) {
	ttl, err = redis.Int64(r.Do("TTL", key))
	return
}

//PTTL key s
// PTL以毫秒秒为单位。当 key不存在时，返回-2.当 key存在但没有设置剩余生存时间时，返回-1
func (r *Redis) PTtl(key string) (ttl int64, err error) {
	return redis.Int64(r.Do("PTTL", key))
}

//RENAME key newkey
//对key值重新命名
func (r *Redis) Rename(key, newKey string) (err error) {
	_, err = r.Do("RENAME", key, newKey)
	return
}

/*--------------------------------------------string操作-----------------------------------*/
//GET key
func (r *Redis) Get(key string) (interface{}, error) {
	return r.Do("GET", key)
}

//GET操作解析string类型的value
func (r *Redis) GetString(key string) (string, error) {
	return redis.String(r.Do("GET", key))
}

//GET操作解析int类型的value
func (r *Redis) GetInt(key string) (int, error) {
	return redis.Int(r.Do("GET", key))
}

//GET操作解析int64类型的value
func (r *Redis) GetInt64(key string) (int64, error) {
	return redis.Int64(r.Do("GET", key))
}

//GET操作解析int64类型的value
func (r *Redis) GetBool(key string) (bool, error) {
	return redis.Bool(r.Do("GET", key))
}

//GET操作解析json封装的value
func (r *Redis) GetObject(key string, val interface{}) error {
	reply, err := r.GetString(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(reply), val)
}

//SET key value
//SET操作，设置为永不过期
func (r *Redis) Set(key string, val interface{}) (interface{}, error) {
	var value interface{}
	switch v := val.(type) {
	case string, int, uint, int8, int16, int32, int64, float32, float64, bool:
		value = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		value = string(b)
	}
	return r.Do("SET", key, value)
}

//SETEX key seconds value
//SetEx操作,设置过期时间(单位为秒)
func (r *Redis) SetEx(key string, val interface{}, expire int) (interface{}, error) {
	var value interface{}
	switch v := val.(type) {
	case string, int, uint, int8, int16, int32, int64, float32, float64, bool:
		value = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		value = string(b)
	}
	return r.Do("SETEX", key, expire, value)
}

//PSETEX key milliseconds value
//PSetEx操作,设置过期时间(单位为毫秒)
func (r *Redis) PSetEx(key string, val interface{}, expire int) (interface{}, error) {
	var value interface{}
	switch v := val.(type) {
	case string, int, uint, int8, int16, int32, int64, float32, float64, bool:
		value = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		value = string(b)
	}
	return r.Do("PSETEX", key, expire, value)
}

//INCR key
//将 key 中储存的数字值增一
func (r *Redis) Incr(key string) (val int64, err error) {
	return redis.Int64(r.Do("INCR", key))
}

//INCRBY key increment
//将 key 所储存的值加上给定的增量值
func (r *Redis) IncrBy(key string, amount int) (val int64, err error) {
	return redis.Int64(r.Do("INCRBY", key, amount))
}

//DECR key
//将 key 中储存的数字值减一
func (r *Redis) Decr(key string) (val int64, err error) {
	return redis.Int64(r.Do("DECR", key))
}

//DECRBY key decrement
//将 key 所储存的值减去给定的减量值
func (r *Redis) DecrBy(key string, amount int) (val int64, err error) {
	return redis.Int64(r.Do("DECRBY", key, amount))
}

/*--------------------------------------------hash操作-----------------------------------*/

//HSET key field value
//将哈希表 key 中的字段 field 的值设为 value
func (r *Redis) Hset(key, field string, val interface{}) (interface{}, error) {
	var value interface{}
	switch v := val.(type) {
	case string, int, uint, int8, int16, int32, int64, float32, float64, bool:
		value = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		value = string(b)
	}
	return r.Do("HSET", key, field, value)
}

//HMSET key field1 value1 [field2 value2 ]
//同时将多个 field-value (域-值)对设置到哈希表 key 中
func (r *Redis) Hmset(key string, keyFieldStruct interface{}) (err error) {
	_, err = r.Do("HMSET", redis.Args{}.Add(key).AddFlat(keyFieldStruct)...)
	return err
}

//HEXISTS key field
//查看哈希表 key 中，指定的字段是否存在
func (r *Redis) Hexists(key, field string) (exists bool, err error) {
	exists, err = redis.Bool(r.Do("HEXISTS", key, field))
	return
}

//HGET key field
//获取存储在哈希表中指定字段的值
func (r *Redis) Hget(key, field string) (reply interface{}, err error) {
	reply, err = r.Do("HGET", key, field)
	return
}

//hash获取value为string类型的值
func (r *Redis) HgetString(key, field string) (reply string, err error) {
	reply, err = redis.String(r.Do("HGET", key, field))
	return
}

//hash获取value为int类型的
func (r *Redis) HgetInt(key, field string) (reply int, err error) {
	reply, err = redis.Int(r.Do("HGET", key, field))
	return
}

//hash获取value为int64类型的
func (r *Redis) HgetInt64(key, field string) (reply int64, err error) {
	reply, err = redis.Int64(r.Do("HGET", key, field))
	return
}

//hash获取value为bool类型的
func (r *Redis) HgetBool(key, field string) (reply bool, err error) {
	reply, err = redis.Bool(r.Do("HGET", key, field))
	return
}

//hash获取value为一个对象类型的
func (r *Redis) HgetObject(key, field string, val interface{}) error {
	reply, err := r.HgetString(key, field)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(reply), val)
}

// Hmget 用法：cache.Redis.HgetAll("key", &val)
func (r *Redis) HgetAll(key string, val interface{}) error {
	v, err := redis.Values(r.Do("HGETALL", key))
	if err != nil {
		return err
	}

	err = redis.ScanStruct(v, val)
	return err
}

//HDEL key field1
//删除一个或多个哈希表字段
func (r *Redis) Hdel(key string, field string) (err error) {
	_, err = r.Do("HDEL", key, field)
	return
}

/*--------------------------------------------list操作-----------------------------------*/

//LPUSH key value1 [value2]
//Redis Lpush 命令将一个或多个值插入到列表头部。 如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。 当 key 存在但不是列表类型时，返回一个错误。
//注意：在Redis 2.4版本以前的 LPUSH 命令，都只接受单个 value 值。
func (r *Redis) Lpush(args ...interface{}) (err error) {
	_, err = r.Do("LPUSH", args...)
	return
}

//LPOP key
//移出并获取列表的第一个元素(value值为string类型)
func (r *Redis) LpopString(key string) (elem string, err error) {
	elem, err = redis.String(r.Do("LPOP", key))
	return
}

//LLEN key
//获取列表长度
func (r *Redis) Llen(key string) (listLen int64, err error) {
	listLen, err = redis.Int64(r.Do("LLEN", key))
	return
}

//LRANGE key start stop
//获取列表指定范围内的元素(value值为string类型)
func (r *Redis) LrangeString(key string, start, end int) (elems []string, err error) {
	elems, err = redis.Strings(r.Do("LRANGE", key, start, end))
	return
}

//Rpush key
//Redis Rpush 命令用于将一个或多个值插入到列表的尾部(最右边)。
//如果列表不存在，一个空列表会被创建并执行 RPUSH 操作。 当列表存在但不是列表类型时，返回一个错误。
func (r *Redis) Rpush(args ...interface{}) (err error) {
	_, err = r.Do("RPUSH", args...)
	return
}

//RPOP key
//移除并获取列表最后一个元素
func (r *Redis) RpopString(key string) (elem string, err error) {
	elem, err = redis.String(r.Do("RPOP", key))
	return
}

/*--------------------------------------------set操作-----------------------------------*/

//SADD key member1 [member2]
//向集合添加一个或多个成员
func (r *Redis) Sadd(args ...interface{}) (err error) {
	_, err = r.Do("SADD", args...)
	return
}

//SCARD key
//获取集合的成员数
func (r *Redis) Scard(key string) (memberNum int64, err error) {
	memberNum, err = redis.Int64(r.Do("SCARD", key))
	return
}

//SISMEMBER key member
//判断 member 元素是否是集合 key 的成员
func (r *Redis) Sismember(key string, member string) (isMember bool, err error) {
	isMember, err = redis.Bool(r.Do("SISMEMBER", key, member))
	return
}

//SMEMBERS key
//返回集合中的所有成员
func (r *Redis) Smembers(key string) (members []string, err error) {
	members, err = redis.Strings(r.Do("SMEMBERS", key))
	return
}

//SPOP key
//移除并返回集合中的一个随机元素
func (r *Redis) Spop(key string) (elem string, err error) {
	elem, err = redis.String(r.Do("RPOP", key))
	return
}

//SREM key member1 [member2]
//移除集合中一个或多个成员
func (r *Redis) Srem(args ...interface{}) (err error) {
	_, err = r.Do("SREM", args...)
	return
}

//SUNION key1 [key2]
//返回所有给定集合的并集
func (r *Redis) Sunion(args ...interface{}) (unionMembers []string, err error) {
	unionMembers, err = redis.Strings(r.Do("SUNION", args...))
	return
}

//SUNIONSTORE destination key1 [key2]
//所有给定集合的并集存储在 destination 集合中
func (r *Redis) Suionstrore(args ...interface{}) (err error) {
	_, err = r.Do("SUNIONSTORE", args...)
	return
}

//SDIFF key1 [key2]
//返回给定所有集合的差集
func (r *Redis) Sdiff(args ...interface{}) (diffMembers []string, err error) {
	diffMembers, err = redis.Strings(r.Do("SDIFF", args...))
	return
}

//SDIFFSTORE destination key1 [key2]
//返回给定所有集合的差集并存储在 destination 中
func (r *Redis) Sdiffstore(args ...interface{}) (err error) {
	_, err = r.Do("SDIFFSTORE", args...)
	return
}

//SMOVE source destination member
//将 member 元素从 source 集合移动到 destination 集合
func (r *Redis) Smove(args ...interface{}) (err error) {
	_, err = r.Do("SMOVE", args...)
	return
}

/*--------------------------------------------zset操作-----------------------------------*/

// ZADD key score1 member1 [score2 member2]
//向有序集合添加一个或多个成员，或者更新已存在成员的分数
func (r *Redis) Zadd(args ...interface{}) (err error) {
	_, err = r.Do("ZADD", args...)
	return
}

//ZCARD key
//获取有序集合的成员数
func (r *Redis) Zcard(key string) (memberNum int64, err error) {
	memberNum, err = redis.Int64(r.Do("ZCARD", key))
	return
}

//ZCOUNT key min max
//计算在有序集合中指定区间分数的成员数
func (r *Redis) Zcount(key string, min, max int64) (memberNum int64, err error) {
	memberNum, err = redis.Int64(r.Do("ZCOUNT", key, min, max))
	return
}

//ZRANGE key start stop [WITHSCORES]
// Zrange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递增(从小到大)来排序。具有相同分数值的成员按字典序(lexicographical order )来排列。
// 以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。或 以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (r *Redis) Zrange(key string, from, to int64) (map[string]int64, error) {
	return redis.Int64Map(r.Do("ZRANGE", key, from, to, "WITHSCORES"))
}

//ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT]
// ZrangeByScore 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
// 具有相同分数值的成员按字典序来排列
func (r *Redis) ZrangeByScore(key string, from, to, offset int64, count int) (map[string]int64, error) {
	return redis.Int64Map(r.Do("ZRANGEBYSCORE", key, from, to, "WITHSCORES", "LIMIT", offset, count))
}

//ZREVRANGE key start stop [WITHSCORES]
// Zrevrange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递减(从大到小)来排列。具有相同分数值的成员按字典序(lexicographical order )来排列。
// 以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。或 以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (r *Redis) Zrevrange(key string, from, to int64) (map[string]int64, error) {
	return redis.Int64Map(r.Do("ZREVRANGE", key, from, to, "WITHSCORES"))
}

//ZREVRANGEBYSCORE key max min [WITHSCORES]
// ZrevrangeByScore 返回有序集中指定分数区间内的所有的成员。有序集成员按分数值递减(从大到小)的次序排列。
// 具有相同分数值的成员按字典序来排列
func (r *Redis) ZrevrangeByScore(key string, from, to, offset int64, count int) (map[string]int64, error) {
	return redis.Int64Map(r.Do("ZREVRANGEBYSCORE", key, from, to, "WITHSCORES", "LIMIT", offset, count))
}

//ZSCORE key member
// Zscore 返回有序集中，成员的分数值。 如果成员元素不是有序集 key 的成员，或 key 不存在，返回 nil
func (r *Redis) Zscore(key string, member string) (int64, error) {
	return redis.Int64(r.Do("ZSCORE", key, member))
}

//ZRANK key member
// Zrank 返回有序集中指定成员的排名。其中有序集成员按分数值递增(从小到大)顺序排列。score 值最小的成员排名为 0
func (r *Redis) Zrank(key, member string) (int64, error) {
	return redis.Int64(r.Do("ZRANK", key, member))
}

//ZREVRANK key member
// Zrevrank 返回有序集中成员的排名。其中有序集成员按分数值递减(从大到小)排序。分数值最大的成员排名为 0 。
func (r *Redis) Zrevrank(key, member string) (int64, error) {
	return redis.Int64(r.Do("ZREVRANK", key, member))
}

//ZREM key member [member ...]
//移除有序集合中的一个或多个成员
func (r *Redis) Zrem(args ...interface{}) (err error) {
	_, err = r.Do("ZREM", args...)
	return
}

//ZREMRANGEBYRANK key start stop
//移除有序集合中给定的排名区间的所有成员
func (r *Redis) ZremRangeByRank(key string, start, end int64) (err error) {
	_, err = r.Do("ZREMRANGEBYRANK", key, start, end)
	return
}

//ZREMRANGEBYSCORE key min max
//移除有序集合中给定的分数区间的所有成员
func (r *Redis) ZremRangeByScore(key string, min, max int64) (err error) {
	_, err = r.Do("ZREMRANGEBYSCORE", key, min, max)
	return
}

/*--------------------------------------------队列消息操作-----------------------------------*/

//PUBLISH channel message
// Publish 将信息发送到指定的频道，返回接收到信息的订阅者数量
func (r *Redis) Publish(channel, message string) (int, error) {
	return redis.Int(r.Do("PUBLISH", channel, message))
}

//SUBSCRIBE channel [channel ...]
//订阅给定的一个或多个频道的信息
func (r *Redis) Subscribe(args ...interface{}) (psc redis.PubSubConn) {
	conn := r.Pool.Get()
	defer conn.Close()

	psc = redis.PubSubConn{Conn: conn}
	err := psc.Subscribe(args...)
	if err != nil {
		return
	}

	for {
		switch n := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("Message: %s %s\n", n.Channel, n.Data)
		case redis.PMessage:
			fmt.Printf("PMessage: %s %s %s\n", n.Pattern, n.Channel, n.Data)
		case redis.Subscription:
			fmt.Printf("Subscription: %s %s %d\n", n.Kind, n.Channel, n.Count)
			if n.Count == 0 {
				return
			}
		case error:
			fmt.Printf("error: %v\n", n)
			return
		}
	}

	return
}
