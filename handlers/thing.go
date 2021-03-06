package handlers

import (
	"encoding/json"
	"io"
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

// HandleListThings returns a list of things as JSON fetched from thing service
func (t *ThingHandler) HandleListThings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := helpers.GetContext(r)
	groupID := r.URL.Query().Get("groupid")
	if groupID == "" {
		http.Error(w, "Error missing groupid, try: URL?groupdid=MY_GROUP_ID", http.StatusBadRequest)
	}
	stream, err := t.ThingSvc.ListGroupThings(ctx, &tpb.GroupRequest{ID: groupID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	streaming := true
	fetched := []models.Thing{}
	for streaming {
		rsp, err := stream.Recv()
		if err == io.EOF {
			streaming = false
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			streaming = false
		} else {
			temp := models.Thing{}
			if err := json.Unmarshal(rsp.GetItem(), &temp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				streaming = false
			}
			fetched = append(fetched, temp)
		}
	}
	bytes, err := json.Marshal(fetched)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(bytes)
}

//HandleCreateThing handles thing creation on POST:/thing
func (t *ThingHandler) HandleCreateThing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
