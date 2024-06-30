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

type scheduleRepo struct {
	db *pgxpool.Pool
}

func NewScheduleRepo(db *pgxpool.Pool) storage.ScheduleRepoI {
	return &scheduleRepo{
		db: db,
	}
}

// Create implements storage.ScheduleRepoI.
func (s *scheduleRepo) Create(ctx context.Context, req *schedule_service.CreateSchedule) (*schedule_service.GetSchedule, error) {
	id := uuid.NewString()

	_, err := s.db.Exec(ctx, `
        INSERT INTO "schedule" (
            id,
            journalId,
            date,
            startTime,
            endTime,
            lesson
        ) VALUES (
            $1, $2, $3, $4, $5, $6
        )`, id, req.JournalId, req.Date, req.StartTime, req.EndTime, req.Lesson)

	if err != nil {
		log.Println("error while creating schedule in storage", err)
		return nil, err
	}

	schedule, err := s.GetByID(ctx, &schedule_service.SchedulePrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting schedule by id after creating", err)
		return nil, err
	}
	return schedule, nil
}

// GetByID implements storage.ScheduleRepoI.
func (s *scheduleRepo) GetByID(ctx context.Context, req *schedule_service.SchedulePrimaryKey) (*schedule_service.GetSchedule, error) {
	resp := &schedule_service.GetSchedule{}

	var (
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := s.db.QueryRow(ctx, `
            SELECT id,
            journalId,
            date,
            startTime,
            endTime,
            lesson,
            created_at,
            updated_at
            FROM "schedule"
        WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.JournalId, &resp.Date, &resp.StartTime, &resp.EndTime, &resp.Lesson, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting schedule by id", err)
		return nil, err
	}

	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

// GetList implements storage.ScheduleRepoI.
func (s *scheduleRepo) GetList(ctx context.Context, req *schedule_service.GetListScheduleRequest) (*schedule_service.GetListScheduleResponse, error) {
	resp := &schedule_service.GetListScheduleResponse{}
	var (
		filter     string
		created_at sql.NullString
		updated_at sql.NullString
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` AND (date ILIKE '%%%v%%' OR startTime ILIKE '%%%v%%' OR endTime ILIKE '%%%v%%' OR lesson ILIKE '%%%v%%' OR journalId::text ILIKE '%%%v%%')`, req.Search, req.Search, req.Search, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := s.db.Query(ctx, `
        SELECT
            id,
            journalId,
            date,
            startTime,
            endTime,
            lesson,
            created_at,
            updated_at
        FROM "schedule" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all schedules:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var schedule schedule_service.Schedule
		count++
		err = rows.Scan(&schedule.Id, &schedule.JournalId, &schedule.Date, &schedule.StartTime, &schedule.EndTime, &schedule.Lesson, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning schedules:", err)
			return nil, err
		}
		schedule.CreatedAt = created_at.String
		schedule.UpdatedAt = updated_at.String

		resp.Schedules = append(resp.Schedules, &schedule)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.ScheduleRepoI.
func (s *scheduleRepo) Update(ctx context.Context, req *schedule_service.UpdateSchedule) (*schedule_service.GetSchedule, error) {

	_, err := s.db.Exec(ctx, `
        UPDATE "schedule" SET
            journalId=$1,
            date=$2,
            startTime=$3,
            endTime=$4,
            lesson=$5,
            updated_at = NOW()
        WHERE id = $6`, req.JournalId, req.Date, req.StartTime, req.EndTime, req.Lesson, req.Id)

	if err != nil {
		log.Println("error while updating schedule in storage", err)
		return nil, err
	}

	schedule, err := s.GetByID(ctx, &schedule_service.SchedulePrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated schedule by id", err)
		return nil, err
	}

	return schedule, nil
}

// Delete implements storage.ScheduleRepoI.
func (s *scheduleRepo) Delete(ctx context.Context, req *schedule_service.SchedulePrimaryKey) error {
	_, err := s.db.Exec(ctx, `
        UPDATE "schedule" SET 
            deleted_at = 1
        WHERE id = $1
    `, req.Id)

	if err != nil {
		log.Println("error while deleting schedule")
		return err
	}

	return nil
}

// GetScheduleForWeek retrieves schedules for a specific week.
func (s *scheduleRepo) GetScheduleForWeek(ctx context.Context, teacherId string, weekStartDate, weekEndDate string) (*schedule_service.GetListScheduleResponse, error) {
	resp := &schedule_service.GetListScheduleResponse{}
	var (
		createdAt sql.NullString
		updatedAt sql.NullString
	)

	query := `
		SELECT id,
			   journalId,
			   date,
			   startTime,
			   endTime,
			   lesson,
			   created_at,
			   updated_at
		FROM "schedule"
		WHERE date BETWEEN $1 AND $2
		AND (teacherId = $3 OR supportTeacherId = $3)
		AND deleted_at = 0
	`
	rows, err := s.db.Query(ctx, query, weekStartDate, weekEndDate, teacherId)
	if err != nil {
		log.Printf("Error while getting schedules for the week: %v", err)
		return nil, err
	}
	defer rows.Close()

	var count int64
	for rows.Next() {
		var schedule schedule_service.Schedule
		count++
		err = rows.Scan(&schedule.Id, &schedule.JournalId, &schedule.Date, &schedule.StartTime, &schedule.EndTime, &schedule.Lesson, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("Error while scanning schedules: %v", err)
			return nil, err
		}
		schedule.CreatedAt = createdAt.String
		schedule.UpdatedAt = updatedAt.String
		resp.Schedules = append(resp.Schedules, &schedule)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// GetScheduleForMonth retrieves schedules for a specific month.
func (s *scheduleRepo) GetScheduleForMonth(ctx context.Context, teacherId string, monthStartDate, monthEndDate string) (*schedule_service.GetListScheduleResponse, error) {
	resp := &schedule_service.GetListScheduleResponse{}
	var (
		createdAt sql.NullString
		updatedAt sql.NullString
	)

	query := `
		SELECT id,
			   journalId,
			   date,
			   startTime,
			   endTime,
			   lesson,
			   created_at,
			   updated_at
		FROM "schedule"
		WHERE date BETWEEN $1 AND $2
		AND (teacherId = $3 OR supportTeacherId = $3)
		AND deleted_at=0
	`
	rows, err := s.db.Query(ctx, query, monthStartDate, monthEndDate, teacherId)
	if err != nil {
		log.Printf("Error while getting schedules for the month: %v", err)
		return nil, err
	}
	defer rows.Close()

	var count int64
	for rows.Next() {
		var schedule schedule_service.Schedule
		count++
		err = rows.Scan(&schedule.Id, &schedule.JournalId, &schedule.Date, &schedule.StartTime, &schedule.EndTime, &schedule.Lesson, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("Error while scanning schedules: %v", err)
			return nil, err
		}
		schedule.CreatedAt = createdAt.String
		schedule.UpdatedAt = updatedAt.String
		resp.Schedules = append(resp.Schedules, &schedule)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}
