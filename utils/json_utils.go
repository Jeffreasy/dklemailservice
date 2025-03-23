package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// JSONResponse stuurt een JSON response naar de client
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		response, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}

		w.Write(response)
	}
}

// ParseJSONBody parseert de JSON body van een request
func ParseJSONBody(r *http.Request, target interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.Unmarshal(body, target)
}
