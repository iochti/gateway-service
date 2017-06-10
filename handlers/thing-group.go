package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iochti/gateway-service/helpers"
	pb "github.com/iochti/thing-group-service/proto"
)

// ThingGroupHandler handles REST calls on /group
type ThingGroupHandler struct {
	ThingGroupSvc pb.ThingGroupSvcClient
	ContentType   string
}

// HandleGetGroup handles GET#/group/:id and sends a group as JSON
func (t *ThingGroupHandler) HandleGetGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	ctx := helpers.GetContext(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rsp, err := t.ThingGroupSvc.GetGroup(ctx, &pb.GroupIDRequest{ID: int32(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	byteResp, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(byteResp)
}

// HandleCreateGroup handles POST#/group and sends the created group as JSON
func (t *ThingGroupHandler) HandleCreateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	ctx := helpers.GetContext(r)
	group := new(pb.ThingGroup)
	if err := json.NewDecoder(r.Body).Decode(group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(group)
	rsp, err := t.ThingGroupSvc.CreateGroup(ctx, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytesResp, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytesResp)
}

// HandleUpdateGroup handles PUT#/group and sends the updated group as JSON
func (t *ThingGroupHandler) HandleUpdateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	ctx := helpers.GetContext(r)
	group := new(pb.ThingGroup)
	if err := json.NewDecoder(r.Body).Decode(group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rsp, err := t.ThingGroupSvc.UpdateGroup(ctx, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytesResp, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytesResp)
}

// HandleDeleteGroup handles DELETE#/group and sends the deleted id as JSON
func (t *ThingGroupHandler) HandleDeleteGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	ctx := helpers.GetContext(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_, err = t.ThingGroupSvc.DeleteGroup(ctx, &pb.GroupIDRequest{ID: int32(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytesResp, err := json.Marshal(struct {
		ID int `json:"id_deleted"`
	}{id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytesResp)
}
