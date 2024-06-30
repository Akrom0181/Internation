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

type administrationRepo struct {
	db *pgxpool.Pool
}

func NewAdministrationRepo(db *pgxpool.Pool) storage.AdministrationRepoI {
	return &administrationRepo{
		db: db,
	}
}

// Create implements storage.AdministrationRepoI.
func (a *administrationRepo) Create(ctx context.Context, req *us.CreateAdministration) (*us.Administration, error) {
	id := uuid.NewString()

	loginlast, err := a.GetLastLogin(ctx)
	if err != nil {
		log.Println("error while creating administration:", err)
		return nil, err
	}
	login := GenerateNewLogin(loginlast)

	_, err = a.db.Exec(ctx, `
		INSERT INTO "administration" (
			id,
			login,
			phone,
			fullname,
			password,
			salary,
			ieltsScore,
			branchId
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) `, id, login, req.Phone, req.Fullname, req.Password, req.Salary, req.IeltsScore, req.BranchId)

	if err != nil {
		log.Println("error while creating administration in storage", err)
		return nil, err
	}

	branch, err := a.GetByID(ctx, &us.AdministrationPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting administration by id after creating", err)
		return nil, err
	}
	return branch, nil
}

// GetByID implements storage.AdministrationRepoI.
func (a *administrationRepo) GetByID(ctx context.Context, req *us.AdministrationPrimaryKey) (*us.Administration, error) {
	resp := &us.Administration{}

	var (
		created_at sql.NullString
		updated_at sql.NullString
		// salary     sql.NullInt32
	)

	err := a.db.QueryRow(ctx, `
	        SELECT 
			id,
	        login,
			phone,
			fullname,
			password,
			salary,
			ieltsScore,
			branchId,
	        created_at,
	        updated_at
	        FROM "administration"
	    WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.Login, &resp.Phone, &resp.Fullname, &resp.Password, &resp.Salary, &resp.IeltsScore, &resp.BranchId, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting administration by id", err)
		return nil, err
	}
	// resp.Salary   = salary.Int32
	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

// GetList implements storage.AdministrationRepoI.
func (a *administrationRepo) GetList(ctx context.Context, req *us.GetListAdministrationRequest) (*us.GetListAdministrationResponse, error) {
	resp := &us.GetListAdministrationResponse{}
	var (
		filter     string
		created_at sql.NullString
		updated_at sql.NullString
		branchId   sql.NullString
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` AND (fullname ILIKE '%%%v%%' OR phone ILIKE '%%%v%%' OR salary ILIKE '%%%v%%')`, req.Search, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := a.db.Query(ctx, `
        SELECT
            id,
			login,
            fullname,
            phone,
			password,
			salary,
			ieltsScore,
			branchId,
            created_at,
            updated_at
        FROM "administration" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all adinistrations:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var adiministration us.Administration
		count++
		err = rows.Scan(&adiministration.Id, &adiministration.Login, &adiministration.Fullname, &adiministration.Phone, &adiministration.Password, &adiministration.Salary,
			&adiministration.IeltsScore, &branchId, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning adinistrations:", err)
			return nil, err
		}
		adiministration.BranchId = branchId.String
		adiministration.CreatedAt = created_at.String
		adiministration.UpdatedAt = updated_at.String

		resp.Administrations = append(resp.Administrations, &adiministration)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.AdministrationRepoI.
func (a *administrationRepo) Update(ctx context.Context, req *us.UpdateAdministration) (*us.Administration, error) {

	_, err := a.db.Exec(ctx, `
        UPDATE "admininstration" SET
			fullname=$1,
			phone=$2,
			password=$3,
			salary=$4,
			ieltsScore=$5,
			branchId=$6,
            updated_at = NOW()
        WHERE id = $1`, req.Id, req.Fullname, req.Phone, req.Password,
		req.Salary, req.IeltsScore, req.BranchId)

	if err != nil {
		log.Println("error while updating admininstration in storage", err)
		return nil, err
	}

	customer, err := a.GetByID(ctx, &us.AdministrationPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated admininstration by id", err)
		return nil, err
	}

	return customer, nil
}

// Delete implements storage.AdministrationRepoI.
func (a *administrationRepo) Delete(ctx context.Context, req *us.AdministrationPrimaryKey) error {
	_, err := a.db.Exec(ctx, `
		UPDATE "administration" SET 
			deleted_at = 1
		WHERE id = $1
	`, req.Id)

	if err != nil {
		log.Println("error while deleting admininstration")
		return err
	}

	return nil
}

func (r *administrationRepo) GetLastLogin(ctx context.Context) (string, error) {
	var login string
	err := r.db.QueryRow(ctx, `
        SELECT login FROM "administration" 
        ORDER BY login DESC LIMIT 1
    `).Scan(&login)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "no rows in result set" {
			log.Println("No rows found, returning default login: A00000")
			return "A00000", nil
		}
		// Log the unexpected error
		log.Printf("Unexpected database error: %v", err)
		return "", fmt.Errorf("database query error: %w", err)
	}

	log.Printf("Retrieved login: %s", login)
	return login, nil
}

func GenerateNewLogin(log string) string {
	prefix := "A"
	numbStr := log[1:]
	num, _ := strconv.Atoi(numbStr)
	newNum := num + 1
	return fmt.Sprintf("%s%05d", prefix, newNum)
}

func (r *administrationRepo) GetByLogin(ctx context.Context, login string) (*us.Administration, error) {
	var resp us.Administration
	var updated_at, created_at sql.NullString

	err := r.db.QueryRow(ctx, `
        SELECT id,login, fullname, phone, password, salary, ieltsScore, branchId, created_at, updated_at
        FROM "administration" 
        WHERE login = $1 AND deleted_at=0
    `, login).Scan(&resp.Id,
		&resp.Login, &resp.Fullname, &resp.Phone, &resp.Password, &resp.Salary, &resp.IeltsScore,
		&resp.BranchId, &created_at, &updated_at,
	)

	resp.UpdatedAt = updated_at.String
	resp.CreatedAt = created_at.String

	if err != nil {
		log.Println("error while getting administration by id:", err)
		return nil, err
	}
	return &resp, nil
}

func (r *administrationRepo) GetReportList(ctx context.Context, req *us.GetReportListAdministrationRequest) (*us.GetReportListAdministrationResponse, error) {
	var resp us.GetReportListAdministrationResponse
	var updated_at, branchId, created_at sql.NullString

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
        SELECT id, login, fullname, phone, password, salary, ieltsScore, branchId, created_at, updated_at
        FROM "administration"
        %s
        LIMIT $1 OFFSET $2
    `, filter)

	rows, err := r.db.Query(ctx, query, req.Limit, offset)
	if err != nil {
		log.Println("error while getting administration list:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var administration us.GetReportAdministrationResponse
		err := rows.Scan(&administration.Id, &administration.Login, &administration.Fullname, &administration.Phone, &administration.Password, &administration.Salary, &administration.IeltsScore, &branchId, &created_at, &updated_at)
		if err != nil {
			log.Println("error while scanning administration row:", err)
			return nil, err
		}

		administration.UpdatedAt = updated_at.String
		administration.BranchId = branchId.String
		administration.CreatedAt = created_at.String

		// Parse created date to calculate days worked
		startDateStr := administration.CreatedAt[:10]
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			log.Println("Error parsing start date:", err)
			return nil, err
		}

		daysWorked := time.Since(startDate).Hours() / 24

		monthlyRate := float64(administration.Salary / 30)
		totalSum := int64(daysWorked * monthlyRate)

		administration.Totalsum = fmt.Sprintf("%d", totalSum)

		resp.Getadministrations = append(resp.Getadministrations, &administration)
	}
	if err := rows.Err(); err != nil {
		log.Println("error after iterating over administration rows:", err)
		return nil, err
	}
	resp.Count = int64(len(resp.Getadministrations))
	return &resp, nil
}
