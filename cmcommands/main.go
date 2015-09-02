package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

var DATAFILEPATH string

func main() {
	user, _ := user.Current()
	DATAFILEPATH = user.HomeDir + getPathSep() + "cm.data"

	if len(os.Args) == 1 {
		returnUsage()
	}

	args := os.Args[1:]

	if args[0] == "add" {
		if len(args) == 1 {
			returnUsage()
		}

		strdata := getStringData()

		strdata += strings.Join(args[1:], " ") + "\n"
		err := ioutil.WriteFile(DATAFILEPATH, []byte(strdata), 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else if args[0] == "list" {
		fmt.Println("")
		for _, cmd := range strings.Split(getStringData(), "\n") {
			parts := strings.Split(cmd, "__")
			if len(parts) > 1 {
				fmt.Println(parts[0] + "[" + strings.TrimSpace(parts[1]) + "]")
			} else {
				fmt.Println(parts[0])
			}
		}
		fmt.Println("")
	} else {
		fmt.Println("")
		searchterm := args[0]
		for _, cmd := range strings.Split(getStringData(), "\n") {
			if strings.Contains(cmd, searchterm) {
				parts := strings.Split(cmd, "__")
				if len(parts) > 1 {
					fmt.Println(parts[0] + "[" + strings.TrimSpace(parts[1]) + "]")
				} else {
					fmt.Println(parts[0])
				}
			}
		}
		fmt.Println("")
	}
}

func returnUsage() {
	msg := "\nUSAGE:\n"
	msg += " \n"
	msg += " # Enter command and optional description seperated by two underscores to add an item\n"
	msg += " cm add <cmd> __ <description>\n"
	msg += " \n"
	msg += " # Enter a single searchterm to search commands and descriptions\n"
	msg += " cm <searchterm>\n"
	msg += " \n"
	msg += " # list all items\n"
	msg += " cm list\n"
	msg += " \n"
	msg += " ** Data stored in " + DATAFILEPATH + "\n"
	log.Fatal(msg)
}

func getStringData() string {
	var strdata string
	data, err := ioutil.ReadFile(DATAFILEPATH)
	if err == nil {
		strdata = string(data)
	} else {
		strdata = ""
	}
	return strdata
}
