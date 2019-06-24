package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func read(r *http.Request, out interface{}) error {
	body := r.Body
	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, out)
}
