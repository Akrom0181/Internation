package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
	us "user_service/genproto/user_service"
	"user_service/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type supportTeacherRepo struct {
	db *pgxpool.Pool
}

func NewSupportTeacherRepo(db *pgxpool.Pool) storage.SupportTeacherRepoI {
	return &supportTeacherRepo{
		db: db,
	}
}

// Create implements storage.SupportTeacherRepoI.
func (s *supportTeacherRepo) Create(ctx context.Context, req *us.CreateSupportTeacher) (*us.SupportTeacher, error) {
	id := uuid.NewString()

	loginlast, err := s.GetLastLogin(ctx)
	if err != nil {
		log.Println("error while creating Sp teacher:", err)
		return nil, err
	}
	login := GenerateNewLoginSPTeacher(loginlast)

	_, err = s.db.Exec(ctx, `
		INSERT INTO "support_teacher" (
			id,
			login,
			fullname,
			phone,
			password,
			salary,
			ieltsScore,
			ieltsAttemptCount,
			branchId
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)`, id, login, req.Fullname, req.Phone, req.Password, req.Salary, req.IeltsScore, req.IeltsAttemptCount, req.BranchId)

	if err != nil {
		log.Println("error while creating support teacher in storage", err)
		return nil, err
	}

	supportTeacher, err := s.GetByID(ctx, &us.SupportTeacherPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting support teacher by id after creating", err)
		return nil, err
	}
	return supportTeacher, nil
}

// GetByID implements storage.SupportTeacherRepoI.
func (s *supportTeacherRepo) GetByID(ctx context.Context, req *us.SupportTeacherPrimaryKey) (*us.SupportTeacher, error) {
	resp := &us.SupportTeacher{}

	var (
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := s.db.QueryRow(ctx, `
	        SELECT id,
			login,
			fullname,
			phone,
			password,
			salary,
			ieltsScore,
			ieltsAttemptCount,
			branchId,
	        created_at,
	        updated_at
	        FROM "support_teacher"
	    WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.Login, &resp.Fullname, &resp.Phone, &resp.Password, &resp.Salary, &resp.IeltsScore, &resp.IeltsAttemptCount, &resp.BranchId, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting support teacher by id", err)
		return nil, err
	}

	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

// GetList implements storage.SupportTeacherRepoI.
func (s *supportTeacherRepo) GetList(ctx context.Context, req *us.GetListSupportTeacherRequest) (*us.GetListSupportTeacherResponse, error) {
	resp := &us.GetListSupportTeacherResponse{}
	var (
		filter     string
		created_at sql.NullString
		updated_at sql.NullString
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` AND (login ILIKE '%%%v%%' OR fullname ILIKE '%%%v%%' OR phone ILIKE '%%%v%%')`, req.Search, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := s.db.Query(ctx, `
        SELECT
            id,
            login,
            fullname,
            phone,
            password,
            salary,
            ieltsScore,
            ieltsAttemptCount,
            branchId,
            created_at,
            updated_at
        FROM "support_teacher" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all support teachers:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var supportTeacher us.SupportTeacher
		count++
		err = rows.Scan(&supportTeacher.Id, &supportTeacher.Login, &supportTeacher.Fullname, &supportTeacher.Phone, &supportTeacher.Password, &supportTeacher.Salary, &supportTeacher.IeltsScore, &supportTeacher.IeltsAttemptCount, &supportTeacher.BranchId, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning support teachers:", err)
			return nil, err
		}
		supportTeacher.CreatedAt = created_at.String
		supportTeacher.UpdatedAt = updated_at.String

		resp.SupportTeachers = append(resp.SupportTeachers, &supportTeacher)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.SupportTeacherRepoI.
func (s *supportTeacherRepo) Update(ctx context.Context, req *us.UpdateSupportTeacher) (*us.SupportTeacher, error) {

	_, err := s.db.Exec(ctx, `
        UPDATE "support_teacher" SET
			login=$1,
			fullname=$2,
			phone=$3,
			password=$4,
			salary=$5,
			ieltsScore=$6,
			ieltsAttemptCount=$7,
			branchId=$8,
            updated_at = NOW()
        WHERE id = $9`, req.Login, req.Fullname, req.Phone, req.Password, req.Salary, req.IeltsScore, req.IeltsAttemptCount, req.BranchId, req.Id)

	if err != nil {
		log.Println("error while updating support teacher in storage", err)
		return nil, err
	}

	supportTeacher, err := s.GetByID(ctx, &us.SupportTeacherPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated support teacher by id", err)
		return nil, err
	}

	return supportTeacher, nil
}

// Delete implements storage.SupportTeacherRepoI.
func (s *supportTeacherRepo) Delete(ctx context.Context, req *us.SupportTeacherPrimaryKey) error {
	_, err := s.db.Exec(ctx, `
		UPDATE "support_teacher" SET 
			deleted_at = 1
		WHERE id = $1
	`, req.Id)

	if err != nil {
		log.Println("error while deleting support teacher")
		return err
	}

	return nil
}

func (s *supportTeacherRepo) GetLastLogin(ctx context.Context) (string, error) {
	var login string
	err := s.db.QueryRow(ctx, `
        SELECT login FROM "support_teacher"
        ORDER BY login DESC LIMIT 1
    `).Scan(&login)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "no rows in result set" {
			log.Println("No rows found, returning default login: ST00000")
			return "ST00000", nil
		}
		// Log the unexpected error
		log.Printf("Unexpected database error: %v", err)
		return "", fmt.Errorf("database query error: %w", err)
	}

	log.Printf("Retrieved login: %s", login)
	return login, nil
}

func GenerateNewLoginSPTeacher(log string) string {
	prefix := "ST"
	numbStr := log[2:]
	num, _ := strconv.Atoi(numbStr)
	newNum := num + 1
	return fmt.Sprintf("%s%05d", prefix, newNum)
}

func (s *supportTeacherRepo) GetByLogin(ctx context.Context, login string) (*us.SupportTeacher, error) {
	var resp us.SupportTeacher
	var updatedAt, branchId, created_at sql.NullString

	err := s.db.QueryRow(ctx, `
		SELECT id,login, fullname, phone, password, salary, ieltsScore, ieltsAttemptCount, branchId, created_at, updated_at
		FROM "support_teacher" 
		WHERE id = $1 AND  deleted_at=0
	`, login).Scan(&resp.Id,
		&resp.Login, &resp.Fullname, &resp.Phone, &resp.Password, &resp.Salary, &resp.IeltsScore, &resp.IeltsAttemptCount,
		&branchId, &created_at, &updatedAt,
	)
	if err != nil {
		log.Println("error while getting teacher by id:", err)
		return nil, err
	}
	resp.CreatedAt = created_at.String
	resp.BranchId = branchId.String
	resp.UpdatedAt = updatedAt.String
	return &resp, nil
}

func (r *supportTeacherRepo) GetReportList(ctx context.Context, req *us.GetReportListSupportTeacherRequest) (*us.GetReportListSupportTeacherResponse, error) {
	var resp us.GetReportListSupportTeacherResponse
	var updatedAt, branchId, created_at sql.NullString

	offset := (req.Page - 1) * req.Limit

	filter := ""

	if req.Search != "" {
		filter = fmt.Sprintf(`
			WHERE (login ILIKE '%%%v%%'
			OR fullname ILIKE '%%%v%%'
			OR phone ILIKE '%%%v%%') AND deleted_at=0
		`, req.Search, req.Search, req.Search)
	}

	query := fmt.Sprintf(`
		SELECT id, login, fullname, phone, password, salary, ieltsScore, ieltsAttemptCount, branchId, created_at, updated_at
		FROM support_teacher
		%s
		LIMIT $1 OFFSET $2
	`, filter)

	rows, err := r.db.Query(ctx, query, req.Limit, offset)
	if err != nil {
		log.Println("error while getting teacher list:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var teacher us.GetReportSupportTeacherResponse
		err := rows.Scan(&teacher.Id, &teacher.Login, &teacher.Fullname, &teacher.Phone, &teacher.Password, &teacher.Salary, &teacher.IeltsScore, &teacher.IeltsAttemptCount, &branchId, &created_at, &updatedAt)
		if err != nil {
			log.Println("error while scanning teacher row:", err)
			return nil, err
		}
		teacher.BranchId = branchId.String
		teacher.CreatedAt = created_at.String
		teacher.UpdatedAt = updatedAt.String

		startDateStr := teacher.CreatedAt[:10]
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			log.Println("Error parsing start date:", err)
			return nil, err
		}

		daysWorked := time.Since(startDate).Hours() / 24

		monthlyRate := float64(teacher.Salary) 

		totalSum := int64(daysWorked * monthlyRate)

		teacher.Totalsum = fmt.Sprintf("%d", totalSum)


		resp.GetSupportTeachers = append(resp.GetSupportTeachers, &teacher)
	}
	if err := rows.Err(); err != nil {
		log.Println("error after iterating over teacher rows:", err)
		return nil, err
	}
	return &resp, nil
}
