package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iochti/gateway-service/helpers"
	"github.com/iochti/thing-service/models"
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

	rsp, err := t.ThingSvc.GetThing(ctx, &tpb.ThingIDRequest{ID: vars["id"]})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(rsp.GetItem())
}

//HandleCreateThing handles thing creation on POST:/thing
func (t *ThingHandler) HandleCreateThing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	ctx := helpers.GetContext(r)
	var thing models.Thing
	if err := json.NewDecoder(r.Body).Decode(&thing); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(thing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rsp, err := t.ThingSvc.CreateThing(ctx, &tpb.Thing{Item: bytes})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(rsp.GetItem())
}

// HandleUpdateThing handles thing update on PUT#/thing
func (t *ThingHandler) HandleUpdateThing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	ctx := helpers.GetContext(r)
	thing := new(models.Thing)
	if err := json.NewDecoder(r.Body).Decode(thing); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(thing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rsp, err := t.ThingSvc.UpdateThing(ctx, &tpb.Thing{Item: bytes})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(rsp.GetItem())
}

// HandleDeleteOneThing handles thing delete on DELETE#/thing/one/:id
func (t *ThingHandler) HandleDeleteOneThing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	ctx := helpers.GetContext(r)
	vars := mux.Vars(r)

	if vars["id"] == "" {
		http.Error(w, "Error: missing id on DELETE", http.StatusBadRequest)
		return
	}

	// Empty response
	_, err := t.ThingSvc.DeleteThing(ctx, &tpb.ThingIDRequest{ID: vars["id"]})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	deleteResponse := struct {
		ID string `json:"id_deleted"`
	}{vars["id"]}
	br, err := json.Marshal(deleteResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(br)
}

type deleteManyRequest struct {
	Ids []string `json:"ids"`
}

// HandleDeleteManyThings handles thing bulk delete on DELETE#/thing/many
func (t *ThingHandler) HandleDeleteManyThings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	ctx := helpers.GetContext(r)
	rq := new(deleteManyRequest)
	if err := json.NewDecoder(r.Body).Decode(rq); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := t.ThingSvc.BulkDeleteThing(ctx, &tpb.ThingIDArray{Things: rq.Ids}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	deleteResponse := struct {
		Ids []string `json:"ids_deleted"`
	}{rq.Ids}
	br, err := json.Marshal(deleteResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(br)
}
