package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	us "user_service/genproto/user_service"
	"user_service/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type studentRepo struct {
	db *pgxpool.Pool
}

func NewStudentRepo(db *pgxpool.Pool) storage.StudentRepoI {
	return &studentRepo{
		db: db,
	}
}

// Create implements storage.StudentRepoI.
func (s *studentRepo) Create(ctx context.Context, req *us.CreateStudent) (*us.Student, error) {
	id := uuid.NewString()

	loginlast, err := s.GetLastLogin(ctx)
	if err != nil {
		log.Println("error while creating student:", err)
		return nil, err
	}

	login:=GenerateNewLoginStudent(loginlast)

	_, err = s.db.Exec(ctx, `
		INSERT INTO "student" (
			id,
			login,
			fullname,
			phone,
			password,
			groupName,
			branchId
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)`, id, login, req.Fullname, req.Phone, req.Password, req.GroupName, req.BranchId)

	if err != nil {
		log.Println("error while creating student in storage", err)
		return nil, err
	}

	student, err := s.GetByID(ctx, &us.StudentPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting student by id after creating", err)
		return nil, err
	}
	return student, nil
}

// GetByID implements storage.StudentRepoI.
func (s *studentRepo) GetByID(ctx context.Context, req *us.StudentPrimaryKey) (*us.Student, error) {
	resp := &us.Student{}

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
			groupName,
			branchId,
	        created_at,
	        updated_at
	        FROM "student"
	    WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.Login, &resp.Fullname, &resp.Phone, &resp.Password, &resp.GroupName, &resp.BranchId, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting student by id", err)
		return nil, err
	}

	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

// GetList implements storage.StudentRepoI.
func (s *studentRepo) GetList(ctx context.Context, req *us.GetListStudentRequest) (*us.GetListStudentResponse, error) {
	resp := &us.GetListStudentResponse{}
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
            groupName,
            branchId,
            created_at,
            updated_at
        FROM "student" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all students:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var student us.Student
		count++
		err = rows.Scan(&student.Id, &student.Login, &student.Fullname, &student.Phone, &student.Password, &student.GroupName, &student.BranchId, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning students:", err)
			return nil, err
		}
		student.CreatedAt = created_at.String
		student.UpdatedAt = updated_at.String

		resp.Students = append(resp.Students, &student)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.StudentRepoI.
func (s *studentRepo) Update(ctx context.Context, req *us.UpdateStudent) (*us.Student, error) {

	_, err := s.db.Exec(ctx, `
        UPDATE "student" SET
			login=$1,
			fullname=$2,
			phone=$3,
			password=$4,
			groupName=$5,
			branchId=$6,
            updated_at = NOW()
        WHERE id = $7`, req.Login, req.Fullname, req.Phone, req.Password, req.GroupName, req.BranchId, req.Id)

	if err != nil {
		log.Println("error while updating student in storage", err)
		return nil, err
	}

	student, err := s.GetByID(ctx, &us.StudentPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated student by id", err)
		return nil, err
	}

	return student, nil
}

// Delete implements storage.StudentRepoI.
func (s *studentRepo) Delete(ctx context.Context, req *us.StudentPrimaryKey) error {
	_, err := s.db.Exec(ctx, `
		UPDATE "student" SET 
			deleted_at = 1
		WHERE id = $1
	`, req.Id)

	if err != nil {
		log.Println("error while deleting student")
		return err
	}

	return nil
}

func (r *studentRepo) GetLastLogin(ctx context.Context) (string, error) {
	var login string
	err := r.db.QueryRow(ctx, `
        SELECT login FROM "student"
        ORDER BY login DESC LIMIT 1
    `).Scan(&login)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "no rows in result set" {
			log.Println("No rows found, returning default login: S00000")
			return "S00000", nil
		}
		// Log the unexpected error
		log.Printf("Unexpected database error: %v", err)
		return "", fmt.Errorf("database query error: %w", err)
	}

	log.Printf("Retrieved login: %s", login)
	return login, nil
}

func GenerateNewLoginStudent(log string) string {
	prefix := "S"
	numbStr := log[1:]
	num, _ := strconv.Atoi(numbStr)
	newNum := num + 1
	return fmt.Sprintf("%s%05d", prefix, newNum)
}

func (r *studentRepo) GetByLogin(ctx context.Context, login string) (*us.Student, error) {
	var resp us.Student
	var created_at, updated_at sql.NullString
	err := r.db.QueryRow(ctx, `
        SELECT id,login, fullname, phone, password, groupName, branchId, created_at, updated_at
        FROM "student"
        WHERE login = $1
    `, login).Scan(&resp.Id,
		&resp.Login, &resp.Fullname, &resp.Phone, &resp.Password, &resp.GroupName, &resp.BranchId,
		&created_at, &updated_at,
	)
	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	if err != nil {
		log.Println("error while getting student by id:", err)
		return nil, err
	}
	return &resp, nil
}

func (r *studentRepo) GetReportList(ctx context.Context, req *us.GetReportListStudentRequest) (*us.GetReportListStudentResponse, error) {
	var resp us.GetReportListStudentResponse
	var created_at, updated_at sql.NullString

	// Calculate OFFSET based on page and limit
	offset := (req.Page - 1) * req.Limit

	// Initialize filter string
	filter := ""

	// Add search condition if req.Search is not empty
	if req.Search != "" {
		filter = fmt.Sprintf(`
            WHERE (s.login ILIKE '%%%v%%'
            OR s.fullname ILIKE '%%%v%%'
            OR s.phone ILIKE '%%%v%%'
            OR s.groupName ILIKE '%%%v%%') AND s.deleted_at=0
        `, req.Search, req.Search, req.Search, req.Search)
	} else {
		filter = " WHERE s.deleted_at IS NULL "
	}

	// Construct the final query with filter, limit, and offset
	query := fmt.Sprintf(`
        SELECT 
            s.id,
            s.login,
            s.fullname,
            s.phone,
            s.password,
            s.groupName,
            s.branchId,
            s.created_at,
            s.updated_at,
            COALESCE(SUM(sp.paidSum), 0) AS total_paid_sum
        FROM 
            student s
        LEFT JOIN 
            student_payment sp ON s.id = sp.studentId
        %s
        GROUP BY 
            s.id,
            s.login,
            s.fullname,
            s.phone,
            s.password,
            s.groupName,
            s.branchId,
            s.created_at,
            s.updated_at
        LIMIT $1 OFFSET $2
    `, filter)

	// Execute the query
	rows, err := r.db.Query(ctx, query, req.Limit, offset)
	if err != nil {
		log.Println("error while getting student list:", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and populate response
	for rows.Next() {
		var student us.GetReportStudentResponse
		err := rows.Scan(&student.Id, &student.Login, &student.Fullname, &student.Phone, &student.Password, &student.GroupName,
			&student.BranchId, &created_at, &updated_at, &student.Paidsum)
		if err != nil {
			log.Println("error while scanning student row:", err)
			return nil, err
		}
		student.CreatedAt = created_at.String
		student.UpdatedAt = updated_at.String
		resp.Students = append(resp.Students, &student)
	}
	if err := rows.Err(); err != nil {
		log.Println("error after iterating over student rows:", err)
		return nil, err
	}
	return &resp, nil
}
