package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const MASTER_LISTENS_ON_PORT = "8050"

var MASTER_HOSTNAME = flag.String("master", "", "Hostname of master server")

type Server struct {
	Name         string
	IPs          []string
	LastModified time.Time
}

type Services []string

func main() {
	flag.Parse()

	if *MASTER_HOSTNAME == "" {
		log.Fatal("USAGE: node --master localhost")
	}

	run()

	services, err := GetServices()

	if err == nil {
		fmt.Println(services.HasApache())
		fmt.Println(services.HasNginx())
		fmt.Println(services.HasPostgresql())
	}
}

func run() {
	servername, err := os.Hostname()
	if err != nil {
		servername = ""
	}

	ips := make([]string, 0)
	allAddrs, _ := net.InterfaceAddrs()
	for i := 0; i < len(allAddrs); i++ {
		addr := allAddrs[i]
		ip := addr.String()
		match, _ := regexp.MatchString("^[0-9]+", strings.Split(ip, ".")[0])
		if match && !strings.Contains(ip, "127.0.0.1") {
			ips = append(ips, strings.Split(ip, "/")[0])
		}
	}

	s := Server{}
	s.Name = servername
	s.IPs = ips
	s.LastModified = time.Now()

	jsonOut, _ := json.Marshal(s)

	sendMessage(string(jsonOut))
}

func sendMessage(data string) {
	client := &http.Client{
		Timeout: time.Second * 3,
	}

	client.PostForm(getUrl(), url.Values{"data": {data}})
}

func getUrl() string {
	return fmt.Sprintf("http://%s:%s/data", *MASTER_HOSTNAME, MASTER_LISTENS_ON_PORT)
}

func GetServices() (Services, error) {
	var err error
	dir := "/etc/init.d"
	services := Services{}
	if _, err = os.Stat(dir); err == nil {
		contents, _ := ioutil.ReadDir(dir)
		for _, content := range contents {
			services = append(services, content.Name())
		}
	}

	return services, err
}

func (arr Services) Contains(s string) bool {
	for _, item := range arr {
		if item == s {
			return true
		}
	}
	return false
}

func (arr Services) HasNginx() bool {
	return arr.Contains("nginx")
}

func (arr Services) HasApache() bool {
	return arr.Contains("apache2")
}

func (arr Services) HasPostgresql() bool {
	return arr.Contains("postgresql")
}
