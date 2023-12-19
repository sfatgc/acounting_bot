package accounting_bot

import (
	"encoding/json"
	"fmt"
	"log"
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
	var d struct {
		UpdateID string `json:"update_id"`
		Message  string `json:"message"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		log.Printf("unexpected message received: %v", r.Body)
		return
	}

	fmt.Printf("Message received: %v", r.Body)

	/* var data tg_response{}
	json.NewEncoder(w).Encode(data)
	*/
}
