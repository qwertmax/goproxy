package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

type App struct {
	Service string `json:"service"`
	Host    string `json:"host"`
	IP      string `json:"ip"`
	Port    string `json:"port"`
}

func GetEndpoint(name string) []App {
	reader := strings.NewReader("")
	request, err := http.NewRequest("GET", "http://52.34.228.148:8123/v1/services/_"+name+"._tcp.marathon.mesos", reader)
	if err != nil {
		panic(err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var apps []App
	json.Unmarshal(body, &apps)

	return apps
}

func Route(apps []App, path string) ([]byte, error) {
	item := rand.Intn(len(apps))
	url := "http://" + apps[item].IP + ":" + apps[item].Port

	if len(path) > 0 {
		url = url + "/" + path
	}

	reader := strings.NewReader("")
	request, err := http.NewRequest("GET", url, reader)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {
	app1 := GetEndpoint("maxapp1")
	app2 := GetEndpoint("maxapp2")

	http.HandleFunc("/app1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "app1:\n")

		resp, err := Route(app1, "")
		if err != nil {
			fmt.Fprintf(w, "error: %s\n", err.Error())
		}

		fmt.Fprintf(w, "%s\n", resp)
	})

	http.HandleFunc("/app2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "app2:\n")

		resp, err := Route(app2, "")
		if err != nil {
			fmt.Fprintf(w, "error: %s\n", err.Error())
		}

		fmt.Fprintf(w, "%s\n", resp)
	})

	http.HandleFunc("/from2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "app1:\n")

		resp, err := Route(app1, "")
		if err != nil {
			fmt.Fprintf(w, "error: %s\n", err.Error())
		}

		fmt.Fprintf(w, "%s\n", resp)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "app1:\n")
		fmt.Fprintf(w, "%v\n", app1)
		fmt.Fprintf(w, "app2:\n")
		fmt.Fprintf(w, "%v\n", app2)
	})
	println("ready")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
