package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iochti/gateway-service/helpers"
	tpb "github.com/iochti/thing-service/proto"
)

// ThingHandler handles /thing REST requests
type ThingHandler struct {
	ThingSvc    tpb.ThingSvcClient
	ContentType string
}

// HandleGetThing handles GET:/thing/:id requests
func (t *ThingHandler) HandleGetThing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	ctx := helpers.GetContext(r)
	vars := mux.Vars(r)

	if vars["id"] == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad ID parameter (not an number)", http.StatusBadRequest)
		return
	}

	rsp, err := t.ThingSvc.GetThing(ctx, &tpb.ThingIDRequest{ID: int32(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(tb)
}

//HandleCreateThing handles thing creation on POST:/thing
func (t *ThingHandler) HandleCreateThing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	ctx := helpers.GetContext(r)
	var user tpb.Thing
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rsp, err := t.ThingSvc.CreateThing(ctx, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(tb)
}

// HandleUpdateThing handles thing update on PUT#/thing
func (t *ThingHandler) HandleUpdateThing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
}

// HandleDeleteOneThing handles thing delete on DELETE#/thing/one/:id
func (t *ThingHandler) HandleDeleteOneThing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
}

// HandleDeleteManyThings handles thing bulk delete on DELETE#/thing/many
func (t *ThingHandler) HandleDeleteManyThings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
}
