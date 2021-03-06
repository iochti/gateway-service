package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iochti/gateway-service/helpers"
	"github.com/iochti/thing-group-service/models"
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := helpers.GetContext(r)
	vars := mux.Vars(r)
	rsp, err := t.ThingGroupSvc.GetGroup(ctx, &pb.GroupIDRequest{ID: vars["id"]})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(rsp.GetItem())
}

func (t *ThingGroupHandler) HandleListGroupsByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := helpers.GetContext(r)
	accountID := r.URL.Query().Get("accountid")
	if accountID == "" {
		http.Error(w, "Error missing groupid, try: URL?accountdid=MY_GROUP_ID", http.StatusBadRequest)
	}
	stream, err := t.ThingGroupSvc.ListUserGroups(ctx, &pb.UserIDRequest{UserId: accountID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	streaming := true
	fetched := []models.ThingGroup{}
	for streaming {
		rsp, err := stream.Recv()
		if err == io.EOF {
			streaming = false
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			streaming = false
		} else {
			temp := models.ThingGroup{}
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

// HandleCreateGroup handles POST#/group and sends the created group as JSON
func (t *ThingGroupHandler) HandleCreateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := helpers.GetContext(r)
	group := new(models.ThingGroup)
	if err := json.NewDecoder(r.Body).Decode(group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bytes, err := json.Marshal(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rsp, err := t.ThingGroupSvc.CreateGroup(ctx, &pb.ThingGroup{Item: bytes})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(rsp.GetItem())
}

// HandleUpdateGroup handles PUT#/group and sends the updated group as JSON
func (t *ThingGroupHandler) HandleUpdateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := helpers.GetContext(r)
	group := models.ThingGroup{}
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bytes, err := json.Marshal(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	rsp, err := t.ThingGroupSvc.UpdateGroup(ctx, &pb.ThingGroup{Item: bytes})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(rsp.GetItem())
}

// HandleDeleteGroup handles DELETE#/group and sends the deleted id as JSON
func (t *ThingGroupHandler) HandleDeleteGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", t.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := helpers.GetContext(r)
	vars := mux.Vars(r)
	_, err := t.ThingGroupSvc.DeleteGroup(ctx, &pb.GroupIDRequest{ID: vars["id"]})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytesResp, err := json.Marshal(struct {
		ID string `json:"id_deleted"`
	}{vars["id"]})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytesResp)
}
