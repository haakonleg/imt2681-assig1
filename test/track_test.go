package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/haakonleg/imt2681-assig1/igcinfo"
	"github.com/haakonleg/imt2681-assig1/track"
)

const listenPort = "8080"

// Start the server in the background to use for the test functions
func init() {
	// Configure and start the API
	go func() {
		app := igcinfo.App{
			ListenPort: listenPort}
		app.StartServer()
	}()

	// Ensure server is started before continuing
	time.Sleep(1000 * time.Millisecond)
}

// Tests the API path POST /api/igc
func TestRegisterTrack(t *testing.T) {
	fmt.Println("Running test TestRegisterTrack")

	request := struct {
		URL string `json:"url"`
	}{"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"}
	var response struct {
		ID int `json:"id"`
	}

	// Add two tracks
	if err := sendPostRequest("/api/igc", request, &response); err != nil {
		t.Fatalf(err.Error())
	}

	if response.ID != 0 {
		t.Fatalf("Expected ID to be 0")
	}

	if err := sendPostRequest("/api/igc", request, &response); err != nil {
		t.Fatalf(err.Error())
	}

	if response.ID != 1 {
		t.Fatalf("Expected ID to be 1")
	}
}

// Tests the API path GET /api/igc
func TestGetTracks(t *testing.T) {
	fmt.Println("Running test TestGetTracks")

	var response []int

	if err := sendGetRequest("/api/igc", &response, true); err != nil {
		t.Fatalf(err.Error())
	}

	if len(response) != 2 {
		t.Fatalf("Expected length of track IDs array to be 2")
	}
}

// Tests the API path GET /api/igc/{id}
func TestGetTrack(t *testing.T) {
	fmt.Println("Running test TestGetTrack")

	var response track.Track

	if err := sendGetRequest("/api/igc/1", &response, true); err != nil {
		t.Fatalf(err.Error())
	}

	expect := track.Track{
		HDate:       "2016-02-19 00:00:00 +0000 UTC",
		Pilot:       "Miguel Angel Gordillo",
		Glider:      "RV8",
		GliderID:    "EC-XLL",
		TrackLength: "1h39m21s"}

	if equal := reflect.DeepEqual(response, expect); !equal {
		t.Fatalf("Response and expected response do not match")
	}
}

// Tests the API path GET /api/igc/{id}/{field}
func TestGetTrackField(t *testing.T) {
	fmt.Println("Running test TestGetTrackField")

	var response string
	var expect string

	if err := sendGetRequest("/api/igc/1/pilot", &response, false); err != nil {
		t.Fatalf(err.Error())
	}

	expect = "Miguel Angel Gordillo"
	if response != expect {
		t.Fatalf("Expected %s, got %s", expect, response)
	}

	if err := sendGetRequest("/api/igc/1/glider", &response, false); err != nil {
		t.Fatalf(err.Error())
	}

	expect = "RV8"
	if response != expect {
		t.Fatalf("Expected %s, got %s", expect, response)
	}

	if err := sendGetRequest("/api/igc/1/glider_id", &response, false); err != nil {
		t.Fatalf(err.Error())
	}

	expect = "EC-XLL"
	if response != expect {
		t.Fatalf("Expected %s, got %s", expect, response)
	}

	if err := sendGetRequest("/api/igc/1/track_length", &response, false); err != nil {
		t.Fatalf(err.Error())
	}

	expect = "1h39m21s"
	if response != expect {
		t.Fatalf("Expected %s, got %s", expect, response)
	}

	if err := sendGetRequest("/api/igc/1/H_date", &response, false); err != nil {
		t.Fatalf(err.Error())
	}

	expect = "2016-02-19 00:00:00 +0000 UTC"
	if response != expect {
		t.Fatalf("Expected %s, got %s", expect, response)
	}
}

func sendPostRequest(path string, requestBody interface{}, responseBody interface{}) error {
	reqBytes, _ := json.Marshal(requestBody)
	body := bytes.NewBuffer(reqBytes)

	resp, err := http.Post("http://:"+listenPort+"/igcinfo"+path, "application/json", body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Got status code" + strconv.Itoa(resp.StatusCode))
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBytes, responseBody)
	if err != nil {
		return err
	}
	return nil
}

func sendGetRequest(path string, responseBody interface{}, isJSON bool) error {
	resp, err := http.Get("http://:" + listenPort + "/igcinfo" + path)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(path + " got status code" + strconv.Itoa(resp.StatusCode))
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if isJSON {
		err = json.Unmarshal(respBytes, responseBody)
		if err != nil {
			return err
		}
	} else {
		*responseBody.(*string) = string(respBytes)
	}

	return nil
}
