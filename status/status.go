package status

import (
	"log"
	//	"ipmi_redis/redis_pk"
	"sync"

	"github.com/garyburd/redigo/redis"
)

type Ipmi_status struct {
	ChOn     chan func()
	Chunkonw chan func()
	wg       sync.WaitGroup
}

func New() *Ipmi_status {
	ipmi := Ipmi_status{
		ChOn:     make(chan func(), 30),
		Chunkonw: make(chan func(), 30),
	}
	ipmi.wg.Add(2)
	go func() {
		defer ipmi.wg.Done()
		for f := range ipmi.ChOn {
			f()
		}

	}()
	go func() {
		defer ipmi.wg.Done()
		for f := range ipmi.Chunkonw {
			f()
		}
	}()
	return &ipmi
}

func (ipmi *Ipmi_status) On_State(p redis.Pool, ip_on []string, key string) {
	ipmi.ChOn <- func() {
		c := p.Get()
		defer c.Close()
		for _, ip_str := range ip_on {

			if ip_str == "0" || ip_str == "" {
				continue
			}
			_, err := c.Do("lpush", key, ip_str)
			if err != nil {
				log.Println("redis set failed:", err)
			}
		}

	}
}
func (ipmi *Ipmi_status) Unkonw_State(p redis.Pool, ip_unkonw []string, key string) {
	ipmi.Chunkonw <- func() {
		c := p.Get()
		defer c.Close()
		for _, ip_str := range ip_unkonw {

			if ip_str == "0" || ip_str == "" {
				continue
			}
			_, err := c.Do("lpush", key, ip_str)
			if err != nil {
				log.Println("redis set failed:", err)
			}
		}

	}
}

func (ipmi *Ipmi_status) Shutdown() {
	close(ipmi.ChOn)
	close(ipmi.Chunkonw)
	ipmi.wg.Wait()
}
