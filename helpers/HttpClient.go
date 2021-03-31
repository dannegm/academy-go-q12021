package helpers

import (
	"encoding/json"
	"net/http"
	"time"
)

var HttpClient = &http.Client{Timeout: 10 * time.Second}

// Fetch JSON Data
func GetJson(url string, target interface{}) error {
	r, err := HttpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
