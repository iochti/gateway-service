package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iochti/gateway-service/helpers"
	usm "github.com/iochti/user-service/models"
	pb "github.com/iochti/user-service/proto"
)

// UserHandler handle /user REST requests
type UserHandler struct {
	UserSvc     pb.UserSvcClient
	ContentType string
}

type deletionResponse struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// HandleGetUser is used on GET:/user and returns a user as JSON
func (u *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", u.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := helpers.GetContext(r)
	params := r.URL.Query()
	categ := params.Get("categ")
	value := params.Get("value")
	if categ == "" || value == "" {
		http.Error(w, "categ & value params must not be empty", http.StatusBadRequest)
		return
	}
	rsp, err := u.UserSvc.GetUser(ctx, &pb.UserRequest{Categ: categ, Value: value})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(rsp.GetUser())
}

// HandleCreateUser is used on POST:/user and returns the created user as JSON
func (u *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", u.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := helpers.GetContext(r)
	var user usm.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userBytes, err := user.ToByte()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	rsp, err := u.UserSvc.CreateUser(ctx, &pb.UserMessage{User: userBytes})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(rsp.GetUser())
}

// HandleDeleteUser is used on DELETE:/user/:id and returns the deletion state & deleted id
func (u *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", u.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := helpers.GetContext(r)
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	rsp, err := u.UserSvc.DeleteUser(ctx, &pb.UserID{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	delResp := deletionResponse{ID: rsp.GetId(), Deleted: rsp.GetDeleted()}
	delRespBytes, err := json.Marshal(delResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(delRespBytes)
}
