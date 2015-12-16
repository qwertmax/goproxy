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
	body, err := MakeRequest("GET", "http://52.34.228.148:8123/v1/services/_"+name+"._tcp.marathon.mesos")
	if err != nil {
		panic(err.Error())
	}

	var apps []App
	err = json.Unmarshal(body, &apps)
	if err != nil {
		panic(err.Error())
	}

	return apps
}

func Route(apps []App, path string) ([]byte, error) {
	item := rand.Intn(len(apps))
	url := "http://" + apps[item].IP + ":" + apps[item].Port

	if len(path) > 0 {
		url = url + "/" + path
	}

	body, err := MakeRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func handleApplications(w http.ResponseWriter, r *http.Request) {
	url := strings.SplitN(strings.TrimLeft(r.URL.Path, "/"), "/", 2)

	fmt.Printf("%#v\n", url)
	if len(url) == 0 || len(url[0]) == 0 {
		fmt.Fprintf(w, "%v\n", struct {
			Type    string `json:"error"`
			Message string `json:"message"`
		}{
			"error",
			"no app name",
		})
		return
	}

	applicationName := url[0]

	var path string
	if len(url) == 1 {
		path = "/"
	} else {
		path = url[1]
	}

	// Router logic starts here
	app := GetEndpoint(applicationName)

	fmt.Fprintf(w, "Application: %s\n", applicationName)

	resp, err := Route(app, path)
	if err != nil {
		fmt.Fprintf(w, "error: %s\n", err.Error())
	}

	fmt.Fprintf(w, "%s", resp)
}

func main() {

	// http.HandleFunc("/app1", func(w http.ResponseWriter, r *http.Request) {
	// 	app1 := GetEndpoint("maxapp1")

	// 	fmt.Fprintf(w, "app1:\n")

	// 	resp, err := Route(app1, "")
	// 	if err != nil {
	// 		fmt.Fprintf(w, "error: %s\n", err.Error())
	// 	}

	// 	fmt.Fprintf(w, "%s\n", resp)
	// })

	// http.HandleFunc("/app2", func(w http.ResponseWriter, r *http.Request) {
	// 	app2 := GetEndpoint("maxapp2")

	// 	fmt.Fprintf(w, "app2:\n")

	// 	resp, err := Route(app2, "")
	// 	if err != nil {
	// 		fmt.Fprintf(w, "error: %s\n", err.Error())
	// 	}

	// 	fmt.Fprintf(w, "%s\n", resp)
	// })

	// http.HandleFunc("/from2", func(w http.ResponseWriter, r *http.Request) {
	// 	app1 := GetEndpoint("maxapp1")

	// 	fmt.Fprintf(w, "app1:\n")

	// 	resp, err := Route(app1, "/from2")
	// 	if err != nil {
	// 		fmt.Fprintf(w, "error: %s\n", err.Error())
	// 	}

	// 	fmt.Fprintf(w, "%s\n", resp)
	// })

	// http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	// 	app1 := GetEndpoint("maxapp1")
	// 	app2 := GetEndpoint("maxapp2")

	// 	fmt.Fprintf(w, "app1:\n")
	// 	fmt.Fprintf(w, "%v\n", app1)
	// 	fmt.Fprintf(w, "app2:\n")
	// 	fmt.Fprintf(w, "%v\n", app2)
	// })
	println("ready")
	log.Fatal(http.ListenAndServe(":3000", http.HandlerFunc(handleApplications)))
}

func MakeRequest(httpType, url string) ([]byte, error) {
	reader := strings.NewReader("")
	request, err := http.NewRequest(httpType, url, reader)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
