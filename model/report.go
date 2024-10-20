package model

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/FakJeongTeeNhoi/report-service/service"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Participant struct {
	Role    string `json:"role"`
	Faculty string `json:"faculty"`
}

type Report struct {
	Id            string        `json:"id"`
	ReservationId string        `json:"reservation_id"`
	RoomId        string        `json:"room_id"`
	RoomName      string        `json:"room_name"`
	SpaceID       string        `json:"space_id"`
	SpaceName     string        `json:"space_name"`
	Status        string        `json:"status"`
	StartDatetime time.Time     `json:"start_datetime"`
	EndDatetime   time.Time     `json:"end_datetime"`
	Participant   []Participant `json:"participant"`
}

type Reserve struct {
	Id            string        `json:"id"`
	Participant   []Participant `json:"participant"`
	Status        string        `json:"status"`
	Approver      string        `json:"approver"`
	StartDatetime time.Time     `json:"start_datetime"`
	EndDatetime   time.Time     `json:"end_datetime"`
	RoomId        string        `json:"room_id"`
	RoomName      string        `json:"room_name"`
	SpaceID       string        `json:"space_id"`
	SpaceName     string        `json:"space_name"`
}

func (r Report) ArrayOfString() []string {
	var res []string
	res = append(
		res,
		r.Id,
		r.ReservationId,
		r.RoomId,
		r.RoomName,
		r.SpaceID,
		r.SpaceName,
		r.Status,
		r.StartDatetime.Format("2006-01-02 15:04:05"),
		r.EndDatetime.Format("2006-01-02 15:04:05"),
	)
	return res
}

func ParseDateTime(raw primitive.DateTime) time.Time {
	return raw.Time()
}

func ParseParticipant(raw primitive.A) []Participant {
	var participants []Participant
	for _, item := range raw {
		if participantMap, ok := item.(bson.M); ok {
			participant := Participant{
				Faculty: participantMap["faculty"].(string),
				Role:    participantMap["role"].(string),
			}
			participants = append(participants, participant)
		}
	}
	return participants
}

func GetReportsBySpace(spaceID string) ([]Report, error) {
	collection := service.DB.Client().Database(os.Getenv("MONGO_DB_NAME")).Collection("reports")

	filter := bson.M{"space_id": spaceID}
	fmt.Println(filter)

	// Execute the query
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	// Iterate over the cursor and decode documents into a slice of Report
	var reports []Report

	for cursor.Next(context.TODO()) {
		var report Report
		var raw bson.M
		if err := cursor.Decode(&raw); err != nil {
			return nil, err
		}
		fmt.Println(raw)

		// Parse UUID fields
		report.Id = raw["id"].(string)
		report.ReservationId = raw["reservation_id"].(string)
		report.RoomId = raw["room_id"].(string)

		// Map other fields
		report.RoomName = raw["room_name"].(string)
		report.SpaceID = raw["space_id"].(string)
		report.SpaceName = raw["space_name"].(string)
		report.Status = raw["status"].(string)

		// Parse time fields
		report.StartDatetime = ParseDateTime(raw["start_datetime"].(primitive.DateTime))
		report.EndDatetime = ParseDateTime(raw["end_datetime"].(primitive.DateTime))

		// Parse participants
		report.Participant = ParseParticipant(raw["participant"].(primitive.A))

		reports = append(reports, report)
	}

	// Check for errors after iterating
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func AddReportFromReserve(reserve Reserve) error {
	collection := service.DB.Client().Database(os.Getenv("MONGO_DB_NAME")).Collection("reports")

	_, err := collection.InsertOne(context.TODO(), bson.M{
		"id":             uuid.New().String(),
		"reservation_id": reserve.Id,
		"room_id":        reserve.RoomId,
		"room_name":      reserve.RoomName,
		"space_id":       reserve.SpaceID,
		"space_name":     reserve.SpaceName,
		"status":         reserve.Status,
		"start_datetime": reserve.StartDatetime,
		"end_datetime":   reserve.EndDatetime,
		"participant":    reserve.Participant,
	})
	if err != nil {
		return err
	}
	return nil
}
