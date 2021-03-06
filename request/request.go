package request

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ResponseType int
type Method int

const (
	JSON ResponseType = iota
	TEXT

	GET Method = iota
	POST
	DELETE
)

// Request contains the context of the HTTP request, it provides some helper methods
// such as sending a JSON response, errors, and parsing JSON requests
type Request struct {
	Method Method

	w http.ResponseWriter
	r *http.Request
}

// CreateRequest creates a new Request object
func CreateRequest(w http.ResponseWriter, r *http.Request, method string) *Request {
	var req Request
	req.w = w
	req.r = r
	switch r.Method {
	case "GET":
		req.Method = GET
	case "POST":
		req.Method = POST
	case "DELETE":
		req.Method = DELETE
	}
	return &req
}

// SetResponseType sets the response type of the HTTP response
func (req *Request) SetResponseType(responseType ResponseType) {
	switch responseType {
	case JSON:
		req.w.Header().Set("Content-Type", "application/json")
	case TEXT:
		req.w.Header().Set("Content-Type", "text/plain")
	}
}

// SendJSON sends a json response, parameter jsonStruct is a struct
// that contains the json fields
func (req *Request) SendJSON(jsonStruct interface{}, statusCode int) {
	req.SetResponseType(JSON)
	res, err := json.Marshal(jsonStruct)
	if err != nil {
		log.Fatal(err)
	}
	req.w.WriteHeader(statusCode)
	req.w.Write(res)
}

// SendText sends a plain text message as response to the request
func (req *Request) SendText(text string) {
	req.SetResponseType(TEXT)
	fmt.Fprint(req.w, text)
}

// ParseJSONRequest parses the JSON contents of a POST request, takes
// a struct as parameter with JSON fields and writes the contents to it
func (req *Request) ParseJSONRequest(jsonStruct interface{}) error {
	body, err := ioutil.ReadAll(req.r.Body)
	defer req.r.Body.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = json.Unmarshal(body, jsonStruct)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// SendError sends a JSON error message in response to the request in case an error occured
func (req *Request) SendError(message string, statusCode int) {
	err := struct {
		Error string `json:"error"`
	}{Error: message}
	req.SendJSON(&err, statusCode)
}
