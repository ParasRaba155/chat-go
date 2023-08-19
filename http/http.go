// http utility
package http

import (
	"encoding/json"
	"net/http"
	"time"

	"app/user/queries"
)

// response the generic response struct
type response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  bool        `json:"status"`
}

// SendResponse will send the response with given info
func SendResponse(w http.ResponseWriter, msg string, data interface{}, statuscode int, success bool) {
	r := response{
		Data:    data,
		Message: msg,
		Status:  success,
	}

	by, err := json.Marshal(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"some error could not send response\"}"))
	}

	w.WriteHeader(statuscode)
	w.Write(by)
}

func GetUserFromRequestContext(r *http.Request) queries.User {
	userCtxKey := struct{}{}
	user, ok := r.Context().Value(userCtxKey).(queries.User)

	if !ok {
		return queries.User{}
	}

	return user
}

func CreateCookie(name, value string, maxAgeSec int) *http.Cookie {
	return &http.Cookie{
		Name:   name,
		Value:  value,
		Path:   "/",
		MaxAge: maxAgeSec,
	}
}

func DeleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})
}
