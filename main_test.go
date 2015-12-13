package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppJSON(t *testing.T) {
	t.Log("testing App struct... (expected Service + Host + IP + Port)")

	app := App{
		Service: "Test Service",
		Host:    "Test Host",
		IP:      "Test IP",
		Port:    "Test Port",
	}

	jsonString, err := json.Marshal(app)
	if err != nil {
		t.Error(err.Error())
	}

	jsonExpected := "{\"service\":\"Test Service\",\"host\":\"Test Host\",\"ip\":\"Test IP\",\"port\":\"Test Port\"}"

	if string(jsonString) != jsonExpected {
		t.Error("json Marshal not working")
	}

}

func TestGetEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "[{\"service\":\"_maxapp1._tcp.marathon.mesos\",\"host\":\"maxapp1-36554-s1.marathon.slave.mesos.\",\"ip\":\"10.0.3.65\",\"port\":\"11926\"},{\"service\":\"_maxapp1._tcp.marathon.mesos\",\"host\":\"maxapp1-54137-s3.marathon.slave.mesos.\",\"ip\":\"10.0.3.67\",\"port\":\"25897\"},{\"service\":\"_maxapp1._tcp.marathon.mesos\",\"host\":\"maxapp1-33389-s5.marathon.slave.mesos.\",\"ip\":\"10.0.3.68\",\"port\":\"27832\"}]")
	}))
	defer server.Close()

	body, err := MakeRequest("GET", server.URL)
	if err != nil {
		t.Error(err.Error())
	}

	expectedJSON := "[{\"service\":\"_maxapp1._tcp.marathon.mesos\",\"host\":\"maxapp1-36554-s1.marathon.slave.mesos.\",\"ip\":\"10.0.3.65\",\"port\":\"11926\"},{\"service\":\"_maxapp1._tcp.marathon.mesos\",\"host\":\"maxapp1-54137-s3.marathon.slave.mesos.\",\"ip\":\"10.0.3.67\",\"port\":\"25897\"},{\"service\":\"_maxapp1._tcp.marathon.mesos\",\"host\":\"maxapp1-33389-s5.marathon.slave.mesos.\",\"ip\":\"10.0.3.68\",\"port\":\"27832\"}]"
	if string(body) != string(expectedJSON) {
		t.Error("JSON responce is not correct")
	}
}

func TestUnmarshal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "[{\"service\":\"_maxapp1._tcp.marathon.mesos\",\"host\":\"maxapp1-36554-s1.marathon.slave.mesos.\",\"ip\":\"10.0.3.65\",\"port\":\"11926\"},{\"service\":\"_maxapp1._tcp.marathon.mesos\",\"host\":\"maxapp1-54137-s3.marathon.slave.mesos.\",\"ip\":\"10.0.3.67\",\"port\":\"25897\"},{\"service\":\"_maxapp1._tcp.marathon.mesos\",\"host\":\"maxapp1-33389-s5.marathon.slave.mesos.\",\"ip\":\"10.0.3.68\",\"port\":\"27832\"}]")
	}))
	defer server.Close()

	body, err := MakeRequest("GET", server.URL)
	if err != nil {
		t.Error(err.Error())
	}

	var app []App
	err = json.Unmarshal(body, &app)
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("expected []App will be 3")
	if len(app) != 3 {
		t.Error("Unmarshal error, len(app) = %d", len(app))
	}

	t.Log("expected app[0].Service will be '_maxapp1._tcp.marathon.mesos'")
	if app[0].Service != "_maxapp1._tcp.marathon.mesos" {
		t.Error("Unmarshal error, app[0].Service = %s", app[0].Service)
	}

	t.Log("expected app[0].Host will be 'maxapp1-36554-s1.marathon.slave.mesos.'")
	if app[0].Host != "maxapp1-36554-s1.marathon.slave.mesos." {
		t.Error("Unmarshal error, app[0].Host = %s", app[0].Host)
	}

	t.Log("expected app[0].IP will be '10.0.3.65'")
	if app[0].IP != "10.0.3.65" {
		t.Error("Unmarshal error, app[0].IP = %s", app[0].IP)
	}

	t.Log("expected app[0].Port will be '11926'")
	if app[0].Port != "11926" {
		t.Error("Unmarshal error, app[0].Port = %s", app[0].Port)
	}

}
