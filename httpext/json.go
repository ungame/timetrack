package httpext

import (
	"encoding/json"
	"github.com/ungame/timetrack/types"
	"log"
	"net/http"
)

const (
	HeaderContentType = "Content-Type"
	HeaderLocation    = "Location"
	HeaderEntity      = "Entity"

	MimeJSON = "application/json"
)

func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set(HeaderContentType, MimeJSON)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("write json failed with error:", err.Error())
	}
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, types.NewError(err))
}
