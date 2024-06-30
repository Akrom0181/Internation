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

type teacherRepo struct {
	db *pgxpool.Pool
}

func NewTeacherRepo(db *pgxpool.Pool) storage.TeacherRepoI {
	return &teacherRepo{
		db: db,
	}
}

// Create implements storage.TeacherRepoI.
func (t *teacherRepo) Create(ctx context.Context, req *us.CreateTeacher) (*us.Teacher, error) {
	id := uuid.NewString()

	loginlast, err := t.GetLastLogin(ctx)
	if err != nil {
		log.Println("error while creating teacher:", err)
		return nil, err
	}
	login := GenerateNewLoginTeacher(loginlast)

	_, err = t.db.Exec(ctx, `
		INSERT INTO "teacher" (
			id,
			login,
			fullname,
			phone,
			password,
			salary,
			ieltsScore,
			ieltsAttemptCount,
			supportTeacherId,
			branchId
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)`, id, login, req.Fullname, req.Phone, req.Password, req.Salary, req.IeltsScore, req.IeltsAttemptCount, req.SupportTeacherId, req.BranchId)

	if err != nil {
		log.Println("error while creating teacher in storage", err)
		return nil, err
	}

	teacher, err := t.GetByID(ctx, &us.TeacherPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting teacher by id after creating", err)
		return nil, err
	}
	return teacher, nil
}

// GetByID implements storage.TeacherRepoI.
func (t *teacherRepo) GetByID(ctx context.Context, req *us.TeacherPrimaryKey) (*us.Teacher, error) {
	resp := &us.Teacher{}

	var (
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := t.db.QueryRow(ctx, `
	        SELECT id,
			login,
			fullname,
			phone,
			password,
			salary,
			ieltsScore,
			ieltsAttemptCount,
			supportTeacherId,
			branchId,
	        created_at,
	        updated_at
	        FROM "teacher"
	    WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.Login, &resp.Fullname, &resp.Phone, &resp.Password, &resp.Salary, &resp.IeltsScore, &resp.IeltsAttemptCount, &resp.SupportTeacherId, &resp.BranchId, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting teacher by id", err)
		return nil, err
	}

	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

// GetList implements storage.TeacherRepoI.
func (t *teacherRepo) GetList(ctx context.Context, req *us.GetListTeacherRequest) (*us.GetListTeacherResponse, error) {
	resp := &us.GetListTeacherResponse{}
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

	rows, err := t.db.Query(ctx, `
        SELECT
            id,
            login,
            fullname,
            phone,
            password,
            salary,
            ieltsScore,
            ieltsAttemptCount,
            supportTeacherId,
            branchId,
            created_at,
            updated_at
        FROM "teacher" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all teachers:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var teacher us.Teacher
		count++
		err = rows.Scan(&teacher.Id, &teacher.Login, &teacher.Fullname, &teacher.Phone, &teacher.Password, &teacher.Salary, &teacher.IeltsScore, &teacher.IeltsAttemptCount, &teacher.SupportTeacherId, &teacher.BranchId, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning teachers:", err)
			return nil, err
		}
		teacher.CreatedAt = created_at.String
		teacher.UpdatedAt = updated_at.String

		resp.Teachers = append(resp.Teachers, &teacher)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.TeacherRepoI.
func (t *teacherRepo) Update(ctx context.Context, req *us.UpdateTeacher) (*us.Teacher, error) {

	_, err := t.db.Exec(ctx, `
        UPDATE "teacher" SET
			login=$1,
			fullname=$2,
			phone=$3,
			password=$4,
			salary=$5,
			ieltsScore=$6,
			ieltsAttemptCount=$7,
			supportTeacherId=$8,
			branchId=$9,
            updated_at = NOW()
        WHERE id = $10`, req.Login, req.Fullname, req.Phone, req.Password, req.Salary, req.IeltsScore, req.IeltsAttemptCount, req.SupportTeacherId, req.BranchId, req.Id)

	if err != nil {
		log.Println("error while updating teacher in storage", err)
		return nil, err
	}

	teacher, err := t.GetByID(ctx, &us.TeacherPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated teacher by id", err)
		return nil, err
	}

	return teacher, nil
}

// Delete implements storage.TeacherRepoI.
func (t *teacherRepo) Delete(ctx context.Context, req *us.TeacherPrimaryKey) error {
	_, err := t.db.Exec(ctx, `
		UPDATE "teacher" SET 
			deleted_at = 1
		WHERE id = $1
	`, req.Id)

	if err != nil {
		log.Println("error while deleting teacher")
		return err
	}

	return nil
}

func (r *teacherRepo) GetLastLogin(ctx context.Context) (string, error) {
	var login string
	err := r.db.QueryRow(ctx, `
        SELECT login FROM "teacher" 
        ORDER BY login DESC LIMIT 1
    `).Scan(&login)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "no rows in result set" {
			log.Println("No rows found, returning default login: T00000")
			return "T00000", nil
		}
		// Log the unexpected error
		log.Printf("Unexpected database error: %v", err)
		return "", fmt.Errorf("database query error: %w", err)
	}

	log.Printf("Retrieved login: %s", login)
	return login, nil
}

func GenerateNewLoginTeacher(log string) string {
	prefix := "T"
	numbStr := log[1:]
	num, _ := strconv.Atoi(numbStr)
	newNum := num + 1
	return fmt.Sprintf("%s%05d", prefix, newNum)
}

func (r *teacherRepo) GetByLogin(ctx context.Context, login string) (*us.Teacher, error) {
	var resp us.Teacher
	var updated_at, branchId, created_at sql.NullString
	err := r.db.QueryRow(ctx, `
        SELECT id, login, fullname, phone, password, salary, ieltsScore, ieltsAttemptCount, supportTeacherId, branchId, created_at, updated_at
        FROM "teacher"
        WHERE login = $1 AND deleted_at=0
    `, login).Scan(&resp.Id,
		&resp.Login, &resp.Fullname, &resp.Phone, &resp.Password, &resp.Salary, &resp.IeltsScore, &resp.IeltsAttemptCount, &resp.SupportTeacherId,
		&branchId, &created_at, &updated_at,
	)
	resp.BranchId = branchId.String
	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String
	if err != nil {
		log.Println("error while getting teacher by id:", err)
		return nil, err
	}
	return &resp, nil
}

func (r *teacherRepo) GetReportList(ctx context.Context, req *us.GetReportListTeacherRequest) (*us.GetReportListTeacherResponse, error) {
	var resp us.GetReportListTeacherResponse
	var (
	updated_at, branchId, created_at, login, phone sql.NullString
	ieltsAttemptCount sql.NullInt32
	salary sql.NullFloat64
	ieltsScore sql.NullFloat64
	)

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
        SELECT id, login, fullname, phone, salary,
         ieltsScore, ieltsAttemptCount, supportTeacherId, branchId, created_at, updated_at
        FROM "teacher"
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
		var teacher us.GetReportTeacherResponse
		err := rows.Scan(&teacher.Id, &login, &teacher.Fullname, &phone, &teacher.Salary, &ieltsScore, &ieltsAttemptCount, &teacher.SupportTeacherId, &branchId, &created_at, &updated_at)
		if err != nil {
			log.Println("error while scanning teacher row:", err)
			return nil, err
		}

		teacher.Phone = phone.String
		teacher.Login = login.String
		teacher.IeltsAttemptCount = ieltsAttemptCount.Int32
		teacher.UpdatedAt = updated_at.String
		teacher.IeltsScore = float32(ieltsScore.Float64)
		teacher.Salary = float32(salary.Float64)
		teacher.BranchId = branchId.String
		teacher.CreatedAt = created_at.String

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

		resp.GetTeacherResponse = append(resp.GetTeacherResponse, &teacher)
	}

	if err := rows.Err(); err != nil {
		log.Println("error after iterating over teacher rows:", err)
		return nil, err
	}

	return &resp, nil
}
