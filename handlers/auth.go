package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	pb "github.com/iochti/auth-service/proto"
	"github.com/iochti/gateway-service/helpers"
)

// AuthHandler wrapper
type AuthHandler struct {
	AuthSvc pb.AuthSvcClient
	Store   sessions.Store
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// HandleLoginURLRequest returns the URL for the user to login
func (a *AuthHandler) HandleLoginURLRequest(w http.ResponseWriter, r *http.Request) {
	ctx := helpers.GetContext(r)
	state := randToken()

	sessions, err := a.Store.Get(r, "state-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessions.Values["state"] = state
	rsp, err := a.AuthSvc.GetLoginURL(ctx, &pb.LoginURLRequest{State: state})
	if err != nil {
		fmt.Println("authsvc error :", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessions.Save(r, w)
	w.Write([]byte(rsp.GetUrl()))
}

func (a *AuthHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	ctx := helpers.GetContext(r)
	session, err := a.Store.Get(r, "state-session")
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
	w.Write([]byte(rsp.GetUser()))
}
