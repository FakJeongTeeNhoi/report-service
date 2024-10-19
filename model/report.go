package model

import (
	"context"
	"time"

	"github.com/FakJeongTeeNhoi/report-service/service"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Participant struct {
	Type    string `json:"type"`
	Faculty string `json:"faculty"`
}

type Report struct {
	Id            uuid.UUID     `json:"id"`
	ReservationId uuid.UUID     `json:"reservation_id"`
	RoomId        uuid.UUID     `json:"room_id"`
	SpaceName     string        `json:"space_name"`
	Status        string        `json:"status"`
	StartDatetime time.Time     `json:"start_datetime"`
	EndDatetime   time.Time     `json:"end_datetime"`
	Participant   []Participant `json:"participant"`
}

type Reserve struct {
	Id            uuid.UUID     `json:"id"`
	Participant   []Participant `json:"participant"`
	Status        string        `json:"status"`
	Approver      string        `json:"approver"`
	StartDatetime time.Time     `json:"start_datetime"`
	EndDatetime   time.Time     `json:"end_datetime"`
	RoomId        uuid.UUID     `json:"room_id"`
	SpaceName     string        `json:"space_name"`
}

func ParseUUID(raw primitive.Binary) uuid.UUID {

	uuidValue, err := uuid.FromBytes(raw.Data)
	if err != nil {
		return uuid.Nil
	}
	return uuidValue
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
				Type:    participantMap["type"].(string),
			}
			participants = append(participants, participant)
		}
	}
	return participants
}

func GetReportsBySpace(spaceName string) ([]Report, error) {
	collection := service.DB.Client().Database("ReportSystem").Collection("report")

	filter := bson.M{"space_name": spaceName}
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

		// Parse UUID fields
		report.Id = ParseUUID(raw["id"].(primitive.Binary))
		report.ReservationId = ParseUUID(raw["reservation_id"].(primitive.Binary))
		report.RoomId = ParseUUID(raw["room_id"].(primitive.Binary))

		// Map other fields
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
	collection := service.DB.Client().Database("ReportSystem").Collection("report")

	_, err := collection.InsertOne(context.TODO(), bson.M{
		"id":             uuid.New(),
		"reservation_id": reserve.Id,
		"room_id":        reserve.RoomId,
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
