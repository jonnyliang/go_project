package my_redis

import (
	//"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

func MyRedis() redis.Conn {
	c, err := redis.Dial("tcp", "X.X.X.X:6379")

	if err != nil {
		log.Fatalf("Connect to redis error %s\n", err)
	}

	return c
}

func GetLrange(keyname string) []string {
	var r []string
	c := MyRedis()
	defer c.Close()
	len_l, err := redis.Int64(c.Do("LLEN", keyname))
	if err != nil {
		log.Fatal(err)
	}
	values, _ := redis.Values(c.Do("lrange", keyname, '0', len_l))
	for _, v := range values {
		r = append(r, string(v.([]byte)))

	}
	return r

}

func SetMark(mark_value, start, mark_key, on_key, unkown_key string) {

	c := MyRedis()
	defer c.Close()
	if start == "0" {
		_, err := c.Do("DEL", on_key, unkown_key)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err := c.Do("SET", mark_key, mark_value)
	if err != nil {
		log.Fatal(err)
	}

}
