package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	asm "github.com/iochti/auth-service/models"
	pb "github.com/iochti/auth-service/proto"
	"github.com/iochti/gateway-service/helpers"
	usm "github.com/iochti/user-service/models"
	userpb "github.com/iochti/user-service/proto"
)

// AuthHandler wrapper
type AuthHandler struct {
	AuthSvc     pb.AuthSvcClient
	UserSvc     userpb.UserSvcClient
	Store       sessions.Store
	ContentType string
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// HandleLoginURLRequest returns the URL for the user to login
func (a *AuthHandler) HandleLoginURLRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", a.ContentType)
	ctx := helpers.GetContext(r)
	state := randToken()
	session, err := a.Store.Get(r, "state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session.Values["auth"] != nil {
		userMsg, fetchErr := a.UserSvc.GetUser(ctx, &userpb.UserRequest{Categ: "login", Value: session.Values["auth"].(string)})
		if fetchErr != nil {
			http.Error(w, fetchErr.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(userMsg.GetUser()))
		return
	}
	session.Values["state"] = state
	rsp, err := a.AuthSvc.GetLoginURL(ctx, &pb.LoginURLRequest{State: state})
	if err != nil {
		fmt.Println("authsvc error :", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytesResp, err := json.Marshal(struct {
		ConnectionURL string `json:"connection_url"`
	}{rsp.GetUrl()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Save(r, w)
	w.Write(bytesResp)
}

// HandleAuth handles GET:/auth request
func (a *AuthHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", a.ContentType)
	ctx := helpers.GetContext(r)
	session, err := a.Store.Get(r, "state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	retrievedState := session.Values["state"]

	// Get params
	params := r.URL.Query()
	state := params.Get("state")
	code := params.Get("code")
	if retrievedState != state {
		fmt.Println(state, session.Values)
		http.Error(w, fmt.Errorf("Expected response state to match stored state").Error(), http.StatusInternalServerError)
		return
	}
	rsp, err := a.AuthSvc.HandleAuth(ctx, &pb.AuthRequest{State: state, Code: code})
	if err != nil {
		fmt.Println("authsvc error :", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ghubUser := asm.GhubUser{}
	if err = json.Unmarshal(rsp.GetUser(), &ghubUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	account := usm.User{}
	account.Email = ghubUser.Email
	account.Name = ghubUser.Name
	account.AvatarURL = ghubUser.AvatarURL
	account.Login = ghubUser.Login

	// Try to get the user in the database
	userMsg, err := a.UserSvc.GetUser(ctx, &userpb.UserRequest{Categ: "login", Value: ghubUser.Login})

	// If the user was not found
	if err != nil && strings.Contains(err.Error(), "not found") {
		var data []byte
		data, err = account.ToByte()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userMsg, err = a.UserSvc.CreateUser(ctx, &userpb.UserMessage{User: data})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["auth"] = ghubUser.Login
	session.Save(r, w)
	w.Write([]byte(userMsg.GetUser()))
}
