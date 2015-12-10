package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/creack/goproxy"
	"github.com/creack/goproxy/registry"
)

type App struct {
	Service string `json:"service"`
	Host    string `json:"host"`
	IP      string `json:"ip"`
	Port    string `json:"port"`
}

var ServiceRegistry = registry.DefaultRegistry{
	"service1": {
		"v1": {
			"localhost:9091",
			"localhost:9092",
		},
	},
}

func GetEndpoint(name string) App {
	reader := strings.NewReader("")
	request, err := http.NewRequest("GET", "http://52.34.228.148:8123/v1/services/_"+name+"._tcp.marathon.mesos", reader)
	if err != nil {
		panic(err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var app []App
	json.Unmarshal(body, &app)

	return app[0]
}

func main() {
	app1 := GetEndpoint("maxapp1")
	app2 := GetEndpoint("maxapp2")

	registry := registry.DefaultRegistry{
		"max": {
			"v1": {
				string(app1.IP + ":" + app1.Port),
				string(app1.IP + ":" + app2.Port),
			},
		},
	}

	// http.HandleFunc("/", goproxy.NewMultipleHostReverseProxy(ServiceRegistry))
	http.HandleFunc("/", goproxy.NewMultipleHostReverseProxy(registry))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "%v\n", ServiceRegistry)
		fmt.Fprintf(w, "%v\n", registry)
	})
	println("ready")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
