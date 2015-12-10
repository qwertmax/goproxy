package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"

	// "github.com/creack/goproxy"
	// "github.com/creack/goproxy/registry"
)

type App struct {
	Service string `json:"service"`
	Host    string `json:"host"`
	IP      string `json:"ip"`
	Port    string `json:"port"`
}

// var ServiceRegistry = registry.DefaultRegistry{
//     "service1": {
//         "v1": {
//             "localhost:9091",
//             "localhost:9092",
//         },
//     },
// 	"max": {
// 		"v1": {
// 			"localhost:9091",
// 			"localhost:9092",
// 		},
// 	},
// }

func GetEndpoint(name string) []App {
	reader := strings.NewReader("")
	// request, err := http.NewRequest("GET", "http://master.mesos:8123/v1/services/_"+name+"._tcp.marathon.mesos", reader)
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

func Route(apps []App) ([]byte, error) {
	item := rand.Intn(len(apps))
	url := apps[item].IP + ":" + apps[item].Port

	reader := strings.NewReader("")
	request, err := http.NewRequest("GET", url, reader)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {
	app1 := GetEndpoint("maxapp1")
	app2 := GetEndpoint("maxapp2")

	// registry := registry.DefaultRegistry{
	// 	"max": {
	// 		"v1": {
	// 			string(app1.IP + ":" + app1.Port),
	// 			string(app1.IP + ":" + app2.Port),
	// 		},
	// 	},
	// }

	http.HandleFunc("/app1", func(w http.ResponseWriter, r *http.Request) {
		resp, err := Route(app1)
		if err != nil {
			fmt.Fprintf(w, "error: %s\n", err.Error())
		}

		fmt.Fprintf(w, "%s\n", resp)
	})

	http.HandleFunc("/app2", func(w http.ResponseWriter, r *http.Request) {
		resp, err := Route(app2)
		if err != nil {
			fmt.Fprintf(w, "error: %s\n", err.Error())
		}

		fmt.Fprintf(w, "%s\n", resp)
	})

	// http.HandleFunc("/", goproxy.NewMultipleHostReverseProxy(ServiceRegistry))
	// http.HandleFunc("/", goproxy.NewMultipleHostReverseProxy(registry))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "%v\n", ServiceRegistry)
		// fmt.Fprintf(w, "%v\n", registry)
		fmt.Fprintf(w, "app1:\n")
		fmt.Fprintf(w, "%v\n", app1)
		fmt.Fprintf(w, "app2:\n")
		fmt.Fprintf(w, "%v\n", app2)
	})
	println("ready")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
