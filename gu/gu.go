package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var port = flag.String("port", "50508", "Hostname of the server")

func main() {
	flag.Parse()
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, address := range addrs {
		ipNet, ok := address.(*net.IPNet)
		if !ok {
			continue
		}
		if ipNet.IP.IsLoopback() {
			handleLoopback(ipNet.IP.String())
			continue
		}
		if ipNet.IP.To4() != nil {
			handleIpV4(ipNet.IP.String())
		}
	}
}

func handleLoopback(ip string) {
	if ret, err := exist(ip + ":" + *port); err == nil {
		fmt.Println(ret + ":" + ip)
	}
}

func handleIpV4(ip string) {
	ipSplit := strings.Split(ip, ".")
	ipBase := strings.Join(ipSplit[:3], ".") + "."
	wg := &sync.WaitGroup{}
	for i := 2; i < 255; i++ {
		ip := ipBase + strconv.Itoa(i)
		fmt.Println(ip)
		wg.Add(1)
		go func() {
			if msg, err := exist(ip + ":" + *port); err == nil {
				fmt.Println("gua: " + msg + ":" + ip)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func exist(ip4 string) (string, error) {
	conn, err := net.DialTimeout("tcp", ip4, 3*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	if _, err = conn.Write([]byte("gu\n")); err != nil {
		return "", err
	}

	buf := make([]byte, 1024)
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if _, err = conn.Read(buf); err != nil {
		return "", err
	}
	msg := string(buf)
	if !strings.HasPrefix(msg, "gua") || len(msg) < 4 {
		return "", errors.New("this is an old six")
	}
	return msg[4:], nil
}
