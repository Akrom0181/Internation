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

type studentPaymentRepo struct {
	db *pgxpool.Pool
}

func NewStudentPaymentRepo(db *pgxpool.Pool) storage.StudentPaymentRepoI {
	return &studentPaymentRepo{
		db: db,
	}
}

func (s *studentPaymentRepo) Create(ctx context.Context, req *schedule_service.CreateStudentPayment) (*schedule_service.GetStudentPayment, error) {
	id := uuid.NewString()

	_, err := s.db.Exec(ctx, `
        INSERT INTO "student_payment" (
            id,
            student_id,
            group_id,
            paidsum,
            administration_id
        ) VALUES (
            $1, $2, $3, $4, $5
        )`, id, req.StudentId, req.GroupId, req.Paidsum, req.AdministrationId)

	if err != nil {
		log.Println("error while creating student payment in storage", err)
		return nil, err
	}

	studentPayment, err := s.GetByID(ctx, &schedule_service.StudentPaymentPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting student payment by id after creating", err)
		return nil, err
	}
	return studentPayment, nil
}

func (s *studentPaymentRepo) GetByID(ctx context.Context, req *schedule_service.StudentPaymentPrimaryKey) (*schedule_service.GetStudentPayment, error) {
	resp := &schedule_service.GetStudentPayment{}

	var (
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := s.db.QueryRow(ctx, `
            SELECT id,
            student_id,
            group_id,
            paidsum,
            administration_id,
            created_at,
            updated_at
            FROM "student_payment"
        WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.StudentId, &resp.GroupId, &resp.Paidsum, &resp.AdministrationId, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting student payment by id", err)
		return nil, err
	}

	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

func (s *studentPaymentRepo) GetList(ctx context.Context, req *schedule_service.GetListStudentPaymentRequest) (*schedule_service.GetListStudentPaymentResponse, error) {
	resp := &schedule_service.GetListStudentPaymentResponse{}
	var (
		filter     string
		created_at sql.NullString
		updated_at sql.NullString
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` AND (student_id ILIKE '%%%v%%' OR group_id ILIKE '%%%v%%' OR paidsum ILIKE '%%%v%%' OR administration_id ILIKE '%%%v%%')`, req.Search, req.Search, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := s.db.Query(ctx, `
        SELECT
            id,
            student_id,
            group_id,
            paidsum,
            administration_id,
            created_at,
            updated_at
        FROM "student_payment" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all student payments:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var studentPayment schedule_service.StudentPayment
		count++
		err = rows.Scan(&studentPayment.Id, &studentPayment.StudentId, &studentPayment.GroupId, &studentPayment.Paidsum, &studentPayment.AdministrationId, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning student payments:", err)
			return nil, err
		}
		studentPayment.CreatedAt = created_at.String
		studentPayment.UpdatedAt = updated_at.String

		resp.StudentPayments = append(resp.StudentPayments, &studentPayment)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

func (s *studentPaymentRepo) Update(ctx context.Context, req *schedule_service.UpdateStudentPayment) (*schedule_service.GetStudentPayment, error) {

	_, err := s.db.Exec(ctx, `
        UPDATE "student_payment" SET
            student_id=$1,
            group_id=$2,
            paidsum=$3,
            administration_id=$4,
            updated_at = NOW()
        WHERE id = $5`, req.StudentId, req.GroupId, req.Paidsum, req.AdministrationId, req.Id)

	if err != nil {
		log.Println("error while updating student payment in storage", err)
		return nil, err
	}

	studentPayment, err := s.GetByID(ctx, &schedule_service.StudentPaymentPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated student payment by id", err)
		return nil, err
	}

	return studentPayment, nil
}

func (s *studentPaymentRepo) Delete(ctx context.Context, req *schedule_service.StudentPaymentPrimaryKey) error {
	_, err := s.db.Exec(ctx, `
        UPDATE "student_payment" SET 
            deleted_at = 1
        WHERE id = $1
    `, req.Id)

	if err != nil {
		log.Println("error while deleting student payment")
		return err
	}

	return nil
}
