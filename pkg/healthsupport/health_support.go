package healthsupport

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandlerFunction(w http.ResponseWriter, _ *http.Request) {
	data, _ := json.Marshal(map[string]string{"status": "pass"})
	log.Printf("we're looking healty.")
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
