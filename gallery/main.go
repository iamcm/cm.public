package main

import (
	datastore "cm.local/db"
	"cm.local/util"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	PORT     = "8010"
	APPROOT  = "/home/chris/hugo/fotodelic"
	IMAGEDIR = APPROOT + "/images"
	DBFILE   = APPROOT + "/db"
)

var (
	db datastore.DataFile
)

type Image struct {
	Sysname  string
	Nicename string
	CatId    string
}

type Cat struct {
	Id   string
	Name string
}

func main() {
	_, err := os.Stat(APPROOT)
	if err != nil {
		log.Fatal(err)
	}

	db = datastore.New(DBFILE)

	http.HandleFunc("/images", imagesHandler)
	http.HandleFunc("/cats", catsHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/cat", catHandler)
	fmt.Println("Serving: http://127.0.0.1:" + PORT)
	http.ListenAndServe(":"+PORT, nil)
}

func uploadHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	rw.Header().Add("Access-Control-Allow-Headers", "accept, cache-control, content-type, x-requested-with")
	rw.Header().Add("Access-Control-Allow-Methods", "POST")

	if req.Method != "POST" {
		return
	}

	// the FormFile function takes in the POST input id file
	file, header, err := req.FormFile("file")

	i := Image{}
	i.Sysname = ""
	for i.Sysname == "" || pathExists(IMAGEDIR+"/"+i.Sysname) {
		i.Sysname = util.RandString(6) + "." + getExtension(header.Filename)
	}
	i.Nicename = header.Filename
	i.CatId = req.FormValue("catId")

	ims := make([]Image, 0)

	rawImages := db.Get("images")
	if len(rawImages) != 0 {
		dec := json.NewDecoder(strings.NewReader(rawImages))
		err = dec.Decode(&ims)
		if err != nil {
			fmt.Fprintln(rw, err)
			return
		}
	}

	ims = append(ims, i)
	bytes, err := json.Marshal(ims)

	fmt.Println(i.Sysname + ":" + string(bytes))

	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}

	db.Set("images", string(bytes))

	defer file.Close()

	out, err := os.Create(IMAGEDIR + "/" + i.Sysname)
	if err != nil {
		fmt.Fprintf(rw, "Unable to create the file for writing. Check your write access privilege")
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(rw, err)
	}

	fmt.Fprintf(rw, "File uploaded successfully : ")
	fmt.Fprintf(rw, header.Filename)
}

func imagesHandler(rw http.ResponseWriter, req *http.Request) {
	ims := make([]Image, 0)
	data := db.Get("images")

	dec := json.NewDecoder(strings.NewReader(data))
	err := dec.Decode(&ims)
	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(rw, data)
}

func catHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	rw.Header().Add("Access-Control-Allow-Headers", "accept, cache-control, content-type, x-requested-with")
	rw.Header().Add("Access-Control-Allow-Methods", "POST")

	if req.Method != "POST" {
		return
	}

	cat := req.FormValue("cat")
	if cat == "" {
		return
	}

	c := Cat{}
	c.Name = cat
	c.Id = md5String(cat)

	cats := make([]Cat, 0)

	data := db.Get("cats")
	if len(data) != 0 {
		dec := json.NewDecoder(strings.NewReader(data))
		err := dec.Decode(&cats)
		if err != nil {
			fmt.Fprintln(rw, err)
			return
		}
	}

	cats = append(cats, c)
	bytes, err := json.Marshal(cats)

	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}

	db.Set("cats", string(bytes))

	http.Redirect(rw, req, req.FormValue("returnTo"), 302)
}

func catsHandler(rw http.ResponseWriter, req *http.Request) {
	cats := make([]Cat, 0)
	data := db.Get("cats")

	dec := json.NewDecoder(strings.NewReader(data))
	err := dec.Decode(&cats)
	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(rw, data)
}

func getExtension(s string) string {
	a := strings.Split(s, ".")
	return a[len(a)-1]
}

func pathExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func md5String(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//add category

//list categories

//list photos in category
