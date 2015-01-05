package main

import (
	"cm.local/jsonstore"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const DATA_FILE_PATH = "hoarder.data"
const WEBSERVER_LISTENS_ON_PORT = "8050"

var DF jsonstore.DataFile

type Server struct {
	Name         string
	IPs          []string
	LastModified time.Time
}

var jsonData []Server

func main() {
	log.Println("Serving on localhost:" + WEBSERVER_LISTENS_ON_PORT)

	DF = jsonstore.DataFile{}
	DF.Path = DATA_FILE_PATH
	setJsonData()

	http.HandleFunc("/data", dataHandler)
	http.ListenAndServe(":"+WEBSERVER_LISTENS_ON_PORT, nil)
}

func dataHandler(rw http.ResponseWriter, req *http.Request) {
	log.Println(req.URL)

	data := req.FormValue("data")

	server := Server{}

	dec := json.NewDecoder(strings.NewReader(data))
	if err := dec.Decode(&server); err != nil {
		log.Fatal(err)
	}

	server.LastModified = time.Now()

	serverExists := false
	for i := 0; i < len(jsonData); i++ {
		log.Println(jsonData[i])
		if jsonData[i].Name == server.Name {
			jsonData[i] = server
			serverExists = true
		}
	}
	if !serverExists {
		jsonData = append(jsonData, server)
	}

	saveJsonData()

	rw.Header().Add("content-type", "text/plain")
	fmt.Fprintln(rw, "")
}

func setJsonData() {
	//fmt.Println("---setting jsonData")
	rawdata := DF.Read()
	if rawdata != "" {
		dec := json.NewDecoder(strings.NewReader(rawdata))
		if err := dec.Decode(&jsonData); err == io.EOF {
			log.Fatal(err)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

func saveJsonData() {
	//fmt.Println("---saving jsonData:")
	//fmt.Println(jsonData)
	jsonOut, _ := json.Marshal(jsonData)
	out := string(jsonOut)
	DF.Write(out)
	setJsonData()
}
