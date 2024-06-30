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

type taskRepo struct {
	db *pgxpool.Pool
}

func NewTaskRepo(db *pgxpool.Pool) storage.TaskRepoI {
	return &taskRepo{
		db: db,
	}
}

// Create implements storage.TaskRepoI.
func (t *taskRepo) Create(ctx context.Context, req *schedule_service.CreateTask) (*schedule_service.GetTask, error) {
	id := uuid.NewString()

	_, err := t.db.Exec(ctx, `
		INSERT INTO "task" (
			id,
			scheduleId,
			label,
			deadlineDate,
			deadlineTime,
			score
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)`, id, req.ScheduleId, req.Label, req.DeadlineDate, req.DeadlineTime, req.Score)

	if err != nil {
		log.Println("error while creating task in storage", err)
		return nil, err
	}

	task, err := t.GetByID(ctx, &schedule_service.TaskPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting task by id after creating", err)
		return nil, err
	}
	return task, nil
}

// GetByID implements storage.TaskRepoI.
func (t *taskRepo) GetByID(ctx context.Context, req *schedule_service.TaskPrimaryKey) (*schedule_service.GetTask, error) {
	resp := &schedule_service.GetTask{}

	var (
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := t.db.QueryRow(ctx, `
	        SELECT id,
			scheduleId,
			label,
			deadlineDate,
			deadlineTime,
			score,
	        created_at,
	        updated_at
	        FROM "task"
	    WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.ScheduleId, &resp.Label, &resp.DeadlineDate, &resp.DeadlineTime, &resp.Score, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting task by id", err)
		return nil, err
	}

	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

// GetList implements storage.TaskRepoI.
func (t *taskRepo) GetList(ctx context.Context, req *schedule_service.GetListTaskRequest) (*schedule_service.GetListTaskResponse, error) {
	resp := &schedule_service.GetListTaskResponse{}
	var (
		filter     string
		created_at sql.NullString
		updated_at sql.NullString
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` AND (label ILIKE '%%%v%%')`, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := t.db.Query(ctx, `
        SELECT
            id,
            scheduleId,
            label,
            deadlineDate,
            deadlineTime,
            score,
            created_at,
            updated_at
        FROM "task" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all tasks:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var task schedule_service.Task
		count++
		err = rows.Scan(&task.Id, &task.ScheduleId, &task.Label, &task.DeadlineDate, &task.DeadlineTime, &task.Score, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning tasks:", err)
			return nil, err
		}
		task.CreatedAt = created_at.String
		task.UpdatedAt = updated_at.String

		resp.Tasks = append(resp.Tasks, &task)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.TaskRepoI.
func (t *taskRepo) Update(ctx context.Context, req *schedule_service.UpdateTask) (*schedule_service.GetTask, error) {

	_, err := t.db.Exec(ctx, `
        UPDATE "task" SET
			scheduleId=$1,
			label=$2,
			deadlineDate=$3,
			deadlineTime=$4,
			score=$5,
            updated_at = NOW()
        WHERE id = $6`, req.ScheduleId, req.Label, req.DeadlineDate, req.DeadlineTime, req.Score, req.Id)

	if err != nil {
		log.Println("error while updating task in storage", err)
		return nil, err
	}

	task, err := t.GetByID(ctx, &schedule_service.TaskPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated task by id", err)
		return nil, err
	}

	return task, nil
}

// Delete implements storage.TaskRepoI.
func (t *taskRepo) Delete(ctx context.Context, req *schedule_service.TaskPrimaryKey) error {
	_, err := t.db.Exec(ctx, `
		UPDATE "task" SET 
			deleted_at = 1
		WHERE id = $1
	`, req.Id)

	if err != nil {
		log.Println("error while deleting task")
		return err
	}

	return nil
}
