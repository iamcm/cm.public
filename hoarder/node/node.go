package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	Name           string
	IPs            []string
	HasNginx       bool
	HasApache      bool
	HasPostgresql  bool
	HasIIS         bool
	HasMsSqlServer bool
	LastModified   time.Time
}

type Services []string

func main() {
	flag.Parse()

	if *MASTER_HOSTNAME == "" {
		*MASTER_HOSTNAME = "localhost"
		//log.Fatal("USAGE: node --master localhost")
	}

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
		if match && !strings.Contains(ip, "127.0.0.1") && !strings.Contains(ip, "0.0.0.0") {
			ips = append(ips, strings.Split(ip, "/")[0])
		}
	}

	s := Server{}
	s.Name = servername
	s.IPs = ips
	s.LastModified = time.Now()

	services := GetServices()
	s.HasNginx = s.HasService(services, "nginx")
	s.HasApache = s.HasService(services, "apache2")
	s.HasPostgresql = s.HasService(services, "postgresql")
	s.HasIIS = s.HasService(services, "IIS")
	s.HasMsSqlServer = s.HasService(services, "Microsoft SQL Server")

	jsonOut, _ := json.Marshal(s)

	sendMessage(string(jsonOut))
}

func sendMessage(data string) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	client.PostForm(getUrl(), url.Values{"data": {data}})
}

func getUrl() string {
	return fmt.Sprintf("http://%s:%s/savedata", *MASTER_HOSTNAME, MASTER_LISTENS_ON_PORT)
}

func GetServices() Services {
	var err error
	services := Services{}
	dirs := []string{"/etc/init.d", "c:/program files", "c:/program files (x86)"}

	for _, dir := range dirs {
		if _, err = os.Stat(dir); err == nil {
			contents, _ := ioutil.ReadDir(dir)
			for _, content := range contents {
				services = append(services, content.Name())
			}
		}
	}

	return services
}

func (arr Services) ContainsPartOfText(s string) bool {
	for _, item := range arr {
		if strings.Contains(item, s) {
			return true
		}
	}
	return false
}

func (s Server) HasService(arr Services, service string) bool {
	return arr.ContainsPartOfText(service)
}
