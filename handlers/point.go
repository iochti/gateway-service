package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	proto "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/iochti/gateway-service/helpers"
	pb "github.com/iochti/point-service/proto"
)

type PointHandler struct {
	ContentType string
	PointSvc    pb.PointSvcClient
}

type createReq struct {
	Client string               `json:"client"`
	Tags   map[string]string    `json:"tags"`
	Fields map[string]fieldType `json:"fields"`
}

type fieldType struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// HandleCreate handles point creation
func (p *PointHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", p.ContentType)
	ctx := helpers.GetContext(r)
	req := createReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fields := map[string]*any.Any{}
	for k, v := range req.Fields {
		value, err := fieldTypeToProtoByte(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fields[k] = &any.Any{
			TypeUrl: "iochti.com/" + v.Type,
			Value:   value,
		}
	}
	_, err := p.PointSvc.CreatePoint(ctx, &pb.Point{User: req.Client, Tags: req.Tags, Fields: fields})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("{\"success\": true}"))
}

// HandleGetFromThings get datas from a thing
func (p *PointHandler) HandleGetFromThings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", p.ContentType)
	ctx := helpers.GetContext(r)
	params := r.URL.Query()
	startTime, err := strconv.Atoi(params.Get("start"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	endTime, err := strconv.Atoi(params.Get("end"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	start := timestamp.Timestamp{Seconds: int64(startTime)}
	end := timestamp.Timestamp{Seconds: int64(endTime)}

	rsp, err := p.PointSvc.GetPointsByThing(ctx, &pb.ThingId{Start: &start, End: &end, ThingId: params.Get("thing_id"), User: params.Get("user")})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(rsp.GetItem())
}

// HandleGetFromThings get datas from a thing
func (p *PointHandler) HandleGetFromGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", p.ContentType)
	ctx := helpers.GetContext(r)
	params := r.URL.Query()
	startTime, err := strconv.Atoi(params.Get("start"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	endTime, err := strconv.Atoi(params.Get("end"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	start := timestamp.Timestamp{Seconds: int64(startTime)}
	end := timestamp.Timestamp{Seconds: int64(endTime)}

	rsp, err := p.PointSvc.GetPointsByGroup(ctx, &pb.GroupId{Start: &start, End: &end, GroupId: params.Get("group_id"), User: params.Get("user")})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(rsp.GetItem())
}

func getProtoMessageName(fieldType string) string {
	switch fieldType {
	case "string":
		return proto.MessageName(&pb.StringPoint{})
	case "float":
		return proto.MessageName(&pb.FloatPoint{})
	case "int":
		return proto.MessageName(&pb.IntegerPoint{})
	case "bool":
		return proto.MessageName(&pb.BoolPoint{})
	case "duration":
		return proto.MessageName(&pb.DurationPoint{})
	case "date-time":
		return proto.MessageName(&pb.DateTimePoint{})
	}
	return ""
}

func fieldTypeToProtoByte(elmt fieldType) ([]byte, error) {
	var value proto.Message
	switch elmt.Type {
	case "string":
		value = &pb.StringPoint{Value: elmt.Value.(string)}
	case "float":
		value = &pb.FloatPoint{Value: float32(elmt.Value.(float64))}
	case "int":
		value = &pb.IntegerPoint{Value: elmt.Value.(int32)}
	case "bool":
		value = &pb.BoolPoint{Value: elmt.Value.(bool)}
	case "duration":
		value = &pb.DurationPoint{Value: elmt.Value.(*duration.Duration)}
	case "date-time":
		value = &pb.DateTimePoint{Value: elmt.Value.(*timestamp.Timestamp)}
	}
	return proto.Marshal(value)
}
