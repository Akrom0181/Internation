package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	es "schedule_service/genproto/schedule_service"
	"schedule_service/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type eventStudentRepo struct {
	db *pgxpool.Pool
}

func NewEventStudentRepo(db *pgxpool.Pool) storage.EventStudentRepoI {
	return &eventStudentRepo{
		db: db,
	}
}

// Create implements storage.EventStudentRepoI.
func (e *eventStudentRepo) Create(ctx context.Context, req *es.CreateEventStudent) (*es.GetEventStudent, error) {
	id := uuid.NewString()

	_, err := e.db.Exec(ctx, `
        INSERT INTO "event_student" (
            id,
            eventId,
            studentId
        ) VALUES (
            $1, $2, $3
        )`, id, req.EventId, req.StudentId)

	if err != nil {
		log.Println("error while creating event student in storage", err)
		return nil, err
	}

	eventStudent, err := e.GetByID(ctx, &es.EventStudentPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting event student by id after creating", err)
		return nil, err
	}
	return eventStudent, nil
}

// GetByID implements storage.EventStudentRepoI.
func (e *eventStudentRepo) GetByID(ctx context.Context, req *es.EventStudentPrimaryKey) (*es.GetEventStudent, error) {
	resp := &es.GetEventStudent{}

	var (
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := e.db.QueryRow(ctx, `
            SELECT id,
            eventId,
            studentId,
            created_at,
            updated_at
            FROM "event_student"
        WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.EventId, &resp.StudentId, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting event student by id", err)
		return nil, err
	}

	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

// GetList implements storage.EventStudentRepoI.
func (e *eventStudentRepo) GetList(ctx context.Context, req *es.GetListEventStudentRequest) (*es.GetListEventStudentResponse, error) {
	resp := &es.GetListEventStudentResponse{}
	var (
		filter     string
		created_at sql.NullString
		updated_at sql.NullString
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` AND (eventId::text ILIKE '%%%v%%' OR studentId::text ILIKE '%%%v%%')`, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := e.db.Query(ctx, `
        SELECT
            id,
            eventId,
            studentId,
            created_at,
            updated_at
        FROM "event_student" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all event students:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var eventStudent es.EventStudent
		count++
		err = rows.Scan(&eventStudent.Id, &eventStudent.EventId, &eventStudent.StudentId, &created_at, &updated_at, &eventStudent.DeletedAt)

		if err != nil {
			log.Println("error while scanning event students:", err)
			return nil, err
		}
		eventStudent.CreatedAt = created_at.String
		eventStudent.UpdatedAt = updated_at.String

		resp.EventStudents = append(resp.EventStudents, &eventStudent)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.EventStudentRepoI.
func (e *eventStudentRepo) Update(ctx context.Context, req *es.UpdateEventStudent) (*es.GetEventStudent, error) {

	_, err := e.db.Exec(ctx, `
        UPDATE "event_student" SET
            eventId=$1,
            studentId=$2,
            updated_at = NOW()
        WHERE id = $3`, req.EventId, req.StudentId, req.Id)

	if err != nil {
		log.Println("error while updating event student in storage", err)
		return nil, err
	}

	eventStudent, err := e.GetByID(ctx, &es.EventStudentPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated event student by id", err)
		return nil, err
	}

	return eventStudent, nil
}

// Delete implements storage.EventStudentRepoI.
func (e *eventStudentRepo) Delete(ctx context.Context, req *es.EventStudentPrimaryKey) error {
	_, err := e.db.Exec(ctx, `
        UPDATE "event_student" SET 
            deleted_at = 1
        WHERE id = $1
    `, req.Id)

	if err != nil {
		log.Println("error while deleting event student")
		return err
	}

	return nil
}

func (r *eventStudentRepo) GetStudentWithEventsByID(ctx context.Context, req *es.EventStudentPrimaryKey) (*es.GetStudentWithEventsResponse, error) {
	var studentResp es.GetStudentWithEventsResponse

	rows, err := r.db.Query(ctx, `
        SELECT s.id, s.name, s.email, s.phone,
               es.id as event_student_id, es.eventId, e.assignStudent, e.topic, e.startTime, e.date, e.branchId, e.created_at as event_created_at, e.updated_at as event_updated_at
        FROM student s
        LEFT JOIN event_student es ON s.id = es.studentId
        LEFT JOIN event e ON es.eventId = e.id
        WHERE s.id = $1
    `, req.Id)
	if err != nil {
		log.Println("error while getting student with events:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var studentID, name, phone string
		var eventStudentID, eventID, assignStudent, topic, startTime, date, branchID, eventCreatedAt, eventUpdatedAt sql.NullString

		err := rows.Scan(&studentID, &name, &phone,
			&eventStudentID, &eventID, &assignStudent, &topic, &startTime, &date, &branchID, &eventCreatedAt, &eventUpdatedAt)
		if err != nil {
			log.Println("error while scanning student with events row:", err)
			return nil, err
		}

		if studentResp.Id == "" {
			studentResp.Id = studentID
			studentResp.Name = name
			studentResp.Phone = phone
		}

		eventResponse := &es.EventStudentResponse{
			Id:            eventID.String,
			EventId:       eventID.String,
			StudentId:     studentID,
			CreatedAt:     eventCreatedAt.String,
			UpdatedAt:     eventUpdatedAt.String,
			AssignStudent: assignStudent.String,
			Topic:         topic.String,
			StartTime:     startTime.String,
			Date:          date.String,
			BranchId:      branchID.String,
		}

		studentResp.Events = append(studentResp.Events, eventResponse)
	}

	if err := rows.Err(); err != nil {
		log.Println("error after iterating over student with events rows:", err)
		return nil, err
	}

	return &studentResp, nil
}
