package client

import (
	"context"
	"fmt"
	"time"

	actProto "proto/activity"
)

type ActivityClientDirectStub struct{}

var memActivities = []*actProto.DailyActivityRecord{}

func (c *ActivityClientDirectStub) CreateActivity(ctx context.Context, req *actProto.CreateActivityRequest) (*actProto.CreateActivityResponse, error) {
	newID := int32(len(memActivities) + 1)
	actDate := req.ActivityDate
	if actDate == "" {
		actDate = time.Now().Format("2006-01-02")
	}

	record := &actProto.DailyActivityRecord{
		Id:           newID,
		UserId:       req.UserId,
		UserNip:      req.UserNip,
		UserName:     req.UserName,
		Title:        req.Title,
		Description:  req.Description,
		ActivityDate: actDate,
		PhotoUrl:     req.PhotoUrl,
		Status:       "PENDING",
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}

	memActivities = append([]*actProto.DailyActivityRecord{record}, memActivities...)

	return &actProto.CreateActivityResponse{
		Success: true,
		Message: "Kegiatan harian berhasil dikirim untuk peninjauan admin",
		Record:  record,
	}, nil
}

func (c *ActivityClientDirectStub) GetActivities(ctx context.Context, req *actProto.GetActivitiesRequest) (*actProto.GetActivitiesResponse, error) {
	var filtered []*actProto.DailyActivityRecord
	for _, item := range memActivities {
		if req.UserId > 0 && item.UserId != req.UserId {
			continue
		}
		if req.UserNip != "" && item.UserNip != req.UserNip {
			continue
		}
		if req.Status != "" && item.Status != req.Status {
			continue
		}
		filtered = append(filtered, item)
	}

	return &actProto.GetActivitiesResponse{
		Success:    true,
		Activities: filtered,
	}, nil
}

func (c *ActivityClientDirectStub) UpdateActivityStatus(ctx context.Context, req *actProto.UpdateActivityStatusRequest) (*actProto.UpdateActivityStatusResponse, error) {
	for _, item := range memActivities {
		if item.Id == req.Id {
			item.Status = req.Status
			item.UpdatedAt = time.Now().Format(time.RFC3339)
			return &actProto.UpdateActivityStatusResponse{
				Success: true,
				Message: fmt.Sprintf("Status kegiatan #%d diubah menjadi %s", req.Id, req.Status),
			}, nil
		}
	}

	return &actProto.UpdateActivityStatusResponse{
		Success: false,
		Error:   "Kegiatan harian tidak ditemukan",
	}, nil
}
