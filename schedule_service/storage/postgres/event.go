package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"schedule_service/genproto/schedule_service"
	"schedule_service/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type eventRepo struct {
	db *pgxpool.Pool
}

func NewEventRepo(db *pgxpool.Pool) storage.EventRepoI {
	return &eventRepo{
		db: db,
	}
}

// Create implements storage.EventRepoI.
func (e *eventRepo) Create(ctx context.Context, req *schedule_service.CreateEvent) (*schedule_service.GetEvent, error) {
	id := uuid.NewString()

	_, err := e.db.Exec(ctx, `
        INSERT INTO "event" (
            id,
            assignStudent,
            topic,
            startTime,
            date,
            branchId
        ) VALUES (
            $1, $2, $3, $4, $5, $6
        )`, id, req.AssignStudent, req.Topic, req.StartTime, req.Date, req.BranchId)

	if err != nil {
		log.Println("error while creating event in storage", err)
		return nil, err
	}

	event, err := e.GetByID(ctx, &schedule_service.EventPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting event by id after creating", err)
		return nil, err
	}
	return event, nil
}

// GetByID implements storage.EventRepoI.
func (e *eventRepo) GetByID(ctx context.Context, req *schedule_service.EventPrimaryKey) (*schedule_service.GetEvent, error) {
	resp := &schedule_service.GetEvent{}

	var (
		startTime  sql.NullString
		date       sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := e.db.QueryRow(ctx, `
            SELECT id,
            assignStudent,
            topic,
            startTime,
            date,
            branchId,
            created_at,
            updated_at
            FROM "event"
        WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.AssignStudent, &resp.Topic, &startTime, &date, &resp.BranchId, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting event by id", err)
		return nil, err
	}

	resp.StartTime = startTime.String
	resp.Date = date.String
	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

// GetList implements storage.EventRepoI.
func (e *eventRepo) GetList(ctx context.Context, req *schedule_service.GetListEventRequest) (*schedule_service.GetListEventResponse, error) {
	resp := &schedule_service.GetListEventResponse{}
	var (
		filter     string
		date       sql.NullString
		startTime  sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` AND (assignStudent ILIKE '%%%v%%' OR topic ILIKE '%%%v%%' OR startTime ILIKE '%%%v%%' OR date ILIKE '%%%v%%' OR branchId::text ILIKE '%%%v%%')`, req.Search, req.Search, req.Search, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := e.db.Query(ctx, `
        SELECT
            id,
            assignStudent,
            topic,
            startTime,
            date,
            branchId,
            created_at,
            updated_at
        FROM "event" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all events:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var event schedule_service.Event
		count++
		err = rows.Scan(&event.Id, &event.AssignStudent, &event.Topic, &startTime, &date, &event.BranchId, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning events:", err)
			return nil, err
		}

		event.Date = date.String
		event.StartTime = startTime.String
		event.CreatedAt = created_at.String
		event.UpdatedAt = updated_at.String

		resp.Events = append(resp.Events, &event)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.EventRepoI.
func (e *eventRepo) Update(ctx context.Context, req *schedule_service.UpdateEvent) (*schedule_service.GetEvent, error) {

	_, err := e.db.Exec(ctx, `
        UPDATE "event" SET
            assignStudent=$1,
            topic=$2,
            startTime=$3,
            date=$4,
            branchId=$5,
            updated_at = NOW()
        WHERE id = $6`, req.AssignStudent, req.Topic, req.StartTime, req.Date, req.BranchId, req.Id)

	if err != nil {
		log.Println("error while updating event in storage", err)
		return nil, err
	}

	event, err := e.GetByID(ctx, &schedule_service.EventPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated event by id", err)
		return nil, err
	}

	return event, nil
}

// Delete implements storage.EventRepoI.
func (e *eventRepo) Delete(ctx context.Context, req *schedule_service.EventPrimaryKey) error {
	_, err := e.db.Exec(ctx, `
        UPDATE "event" SET 
            deleted_at = 1
        WHERE id = $1
    `, req.Id)

	if err != nil {
		log.Println("error while deleting event")
		return err
	}

	return nil
}
