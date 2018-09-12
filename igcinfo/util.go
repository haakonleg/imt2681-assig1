package igcinfo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ParseJSONRequest parses a json request
func ParseJSONRequest(r *http.Request, result interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		fmt.Println(err)
		return errors.New("Error parsing body")
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		fmt.Println(err)
		return errors.New("Invalid JSON")
	}

	return nil
}

// WriteError writes an error in the HTTP response
func WriteError(w http.ResponseWriter, err string, statusCode int) {
	http.Error(w, "{\"error\":\""+err+"\"}", statusCode)
}

// WriteMessage writes a message in the HTTP response
func WriteMessage(w http.ResponseWriter, msg string) {
	fmt.Fprint(w, "{\"msg\":\""+msg+"\"}")
}
