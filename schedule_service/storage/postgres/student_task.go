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

type studentTaskRepo struct {
	db *pgxpool.Pool
}

func NewStudentTaskRepo(db *pgxpool.Pool) storage.StudentTaskRepoI {
	return &studentTaskRepo{
		db: db,
	}
}

// Create implements storage.StudentTaskRepoI.
func (s *studentTaskRepo) Create(ctx context.Context, req *schedule_service.CreateStudentTask) (*schedule_service.GetStudentTask, error) {
	id := uuid.NewString()

	_, err := s.db.Exec(ctx, `
        INSERT INTO "student_task" (
            id,
            taskId,
            studentId
        ) VALUES (
            $1, $2, $3
        )`, id, req.TaskId, req.StudentId)

	if err != nil {
		log.Println("error while creating student task in storage", err)
		return nil, err
	}

	studentTask, err := s.GetByID(ctx, &schedule_service.StudentTaskPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting student task by id after creating", err)
		return nil, err
	}
	return studentTask, nil
}

// GetByID implements storage.StudentTaskRepoI.
func (s *studentTaskRepo) GetByID(ctx context.Context, req *schedule_service.StudentTaskPrimaryKey) (*schedule_service.GetStudentTask, error) {
	resp := &schedule_service.GetStudentTask{}

	var (
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := s.db.QueryRow(ctx, `
            SELECT id,
            taskId,
            studentId,
            created_at,
            updated_at
            FROM "student_task"
        WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.TaskId, &resp.StudentId, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting student task by id", err)
		return nil, err
	}

	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

// GetList implements storage.StudentTaskRepoI.
func (s *studentTaskRepo) GetList(ctx context.Context, req *schedule_service.GetListStudentTaskRequest) (*schedule_service.GetListStudentTaskResponse, error) {
	resp := &schedule_service.GetListStudentTaskResponse{}
	var (
		filter     string
		created_at sql.NullString
		updated_at sql.NullString
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` AND (taskId ILIKE '%%%v%%' OR studentId ILIKE '%%%v%%')`, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := s.db.Query(ctx, `
        SELECT
            id,
            taskId,
            studentId,
            created_at,
            updated_at
        FROM "student_task" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all student tasks:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var studentTask schedule_service.StudentTask
		count++
		err = rows.Scan(&studentTask.Id, &studentTask.TaskId, &studentTask.StudentId, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning student tasks:", err)
			return nil, err
		}
		studentTask.CreatedAt = created_at.String
		studentTask.UpdatedAt = updated_at.String

		resp.StudentTasks = append(resp.StudentTasks, &studentTask)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.StudentTaskRepoI.
func (s *studentTaskRepo) Update(ctx context.Context, req *schedule_service.UpdateStudentTask) (*schedule_service.GetStudentTask, error) {

	_, err := s.db.Exec(ctx, `
        UPDATE "student_task" SET
            taskId=$1,
            studentId=$2,
            updated_at = NOW()
        WHERE id = $3`, req.TaskId, req.StudentId, req.Id)

	if err != nil {
		log.Println("error while updating student task in storage", err)
		return nil, err
	}

	studentTask, err := s.GetByID(ctx, &schedule_service.StudentTaskPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated student task by id", err)
		return nil, err
	}

	return studentTask, nil
}

// Delete implements storage.StudentTaskRepoI.
func (s *studentTaskRepo) Delete(ctx context.Context, req *schedule_service.StudentTaskPrimaryKey) error {
	_, err := s.db.Exec(ctx, `
        UPDATE "student_task" SET 
            deleted_at = 1
        WHERE id = $1
    `, req.Id)

	if err != nil {
		log.Println("error while deleting student task")
		return err
	}

	return nil
}

func (s *studentTaskRepo) UpdateScoreforTeacher(ctx context.Context, req *schedule_service.UpdateStudentScoreRequest) (*schedule_service.GetStudentTask, error) {
    _, err := s.db.Exec(ctx, `
        UPDATE "student_task"
        SET score = $1
        WHERE id = $2
    `, req.Score, req.Id)
    if err != nil {
        log.Println("error while updating student_task:", err)
        return nil, err
    }

    studentTask, err := s.GetByID(ctx, &schedule_service.StudentTaskPrimaryKey{Id: req.Id})
    if err != nil {
        log.Println("error while getting student_task by id after update:", err)
        return nil, err
    }

    return studentTask, nil
}

func (s *studentTaskRepo) UpdateScoreforStudent(ctx context.Context, req *schedule_service.UpdateStudentScoreRequest) (*schedule_service.GetStudentTask, error) {
    _, err := s.db.Exec(ctx, `
        UPDATE "student_task"
        SET score = $1
        WHERE id = $2
    `, req.Score, req.Id)
    if err != nil {
        log.Println("error while updating student_task:", err)
        return nil, err
    }

    studentTask, err := s.GetByID(ctx, &schedule_service.StudentTaskPrimaryKey{Id: req.Id})
    if err != nil {
        log.Println("error while getting student_task by id after update:", err)
        return nil, err
    }
    return studentTask, nil
}