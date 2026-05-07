package response

import (
	"encoding/json"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if msg, ok := data.(proto.Message); ok {
			marshaller := protojson.MarshalOptions{
				EmitUnpopulated: true,
				UseProtoNames:   true,
			}
			b, err := marshaller.Marshal(msg)
			if err == nil {
				_, _ = w.Write(b)
				return
			}
		}
		_ = json.NewEncoder(w).Encode(data)
	}
}

func Error(w http.ResponseWriter, status int, code, message string) {
	resp := ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	JSON(w, status, resp)
}
