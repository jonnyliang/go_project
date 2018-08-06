// ipmi project main.go
package main

import (
	//	"log"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"ipmi/db"
	"ipmi/my_redis"
	"ipmi/status"
	"ipmi/work"

	"github.com/garyburd/redigo/redis"
)

var RedisClient *redis.Pool

const (
	maxGoroutines = 50 //number of process
	timeout       = 2 * time.Minute
	ipmi_on       = "ipmi_on"
	ipmi_unkown   = "ipmi_unkown"
	ipmi_mark     = "ipmi_mark"
)

func init() {
	RedisClient = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   maxGoroutines,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "X.X.X.X:6379")
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}

// ip 地址段
type ip_range struct {
	name    string
	ipmi    *status.Ipmi_status
	timeout <-chan time.Time
}

func (ips *ip_range) Task() {
	ips.timeout = time.After(timeout)
	var out []byte
	var err error
	var cmd_str string
	var ipmi_on_slice, ipmi_unkown_slice, ipmi_443_slice, ipmi_623_slice []string

	cmd_str = fmt.Sprintf("nmap -sS -p 443 %s | grep -B 3 open | grep report | awk '{print $5}'", ips.name)
	cmd_1 := exec.Command("/bin/sh", "-c", cmd_str)
	if out, err = cmd_1.Output(); err != nil {
		fmt.Println(err)
	}
	ipmi_443_slice = strings.Split(string(out), "\n")

	cmd_str = fmt.Sprintf("nmap -sS -p 623 %s | grep -B 3 open | grep report | awk '{print $5}'", ips.name)
	cmd_2 := exec.Command("/bin/sh", "-c", cmd_str)
	if out, err = cmd_2.Output(); err != nil {
		fmt.Println(err)
	}
	ipmi_623_slice = strings.Split(string(out), "\n")

	ipmi_on_slice, ipmi_unkown_slice = Get_duplicate(ipmi_443_slice, ipmi_623_slice)
	ips.ipmi.On_State(*RedisClient, ipmi_on_slice, ipmi_on)
	ips.ipmi.Unkonw_State(*RedisClient, ipmi_unkown_slice, ipmi_unkown)
	select {
	case <-ips.timeout:
		return
	default:
	}
}

func main() {
	t1 := time.Now()
	the_time := t1.Format("2006-01-02 15:04:05")
	nmap_pool := work.New(maxGoroutines)
	ipmi := status.New()
	//start mark set
	my_redis.SetMark("0_"+the_time, "0", ipmi_mark, ipmi_on, ipmi_unkown)
	//delete the T status of nmap commond last time
	exec.Command("/bin/sh", "-c", "ps aux | grep 'nmap -sS' | grep T | awk '{print $2}' | xargs kill -9")

	ips_slice := db.Get_ips()
	for _, ips := range ips_slice {
		ipr := ip_range{
			name: ips,
			ipmi: ipmi,
		}
		nmap_pool.Run(&ipr)

	}
	nmap_pool.Shutdown()
	ipmi.Shutdown()
	fmt.Println(time.Since(t1))
	fmt.Println("done!")
	my_redis.SetMark("1_"+time.Now().Format("2006-01-02 15:04:05"), "1", ipmi_mark, ipmi_on, ipmi_unkown)
}

//取交集 和 差集
func Get_duplicate(slice1 []string, slice2 []string) ([]string, []string) {
	list := append(slice1, slice2[0:]...)
	var x []string = []string{} //unsame
	var y []string = []string{} // same
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					y = append(y, i)
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x, y
}
