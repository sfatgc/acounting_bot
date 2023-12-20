package accounting_bot

import (
	"fmt"
	"io"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("dispatchMessages", dispatchMessages)
}

/* type tg_response struct {
	method string
} */

func dispatchMessages(w http.ResponseWriter, r *http.Request) {
	/* var d struct {
		UpdateID string `json:"update_id"`
		Message  string `json:"message"`
	} */

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	/* 	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
	   		log.Printf("unexpected message received: %v", r.Body)
	   		return
	   	}
	*/
	if err != nil {
		fmt.Printf("Error reading request body. %v", r.URL)
	} else {
		fmt.Printf("Message received: %v", string(value))
	}
	/* var data tg_response{}
	json.NewEncoder(w).Encode(data)
	*/
}
