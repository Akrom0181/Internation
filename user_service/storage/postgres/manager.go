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

type managerRepo struct {
	db *pgxpool.Pool
}

func NewManagerRepo(db *pgxpool.Pool) storage.ManagerRepoI {
	return &managerRepo{
		db: db,
	}
}

// Create implements storage.ManagerRepoI.
func (m *managerRepo) Create(ctx context.Context, req *us.CreateManager) (*us.Manager, error) {
	id := uuid.NewString()

	loginlast, err := m.GetLastLogin(ctx)
	if err != nil {
		log.Println("error while creating menager:", err)
		return nil, err
	}
	login:=GenerateNewLoginManager(loginlast)

	_, err = m.db.Exec(ctx, `
		INSERT INTO "manager" (
			id,
			login,
			fullname,
			salary,
			phone,
			password,
			branchId
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)`, id, login, req.Fullname, req.Salary, req.Phone, req.Password, req.BranchId)

	if err != nil {
		log.Println("error while creating manager in storage", err)
		return nil, err
	}

	manager, err := m.GetByID(ctx, &us.ManagerPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting manager by id after creating", err)
		return nil, err
	}
	return manager, nil
}

// GetByID implements storage.ManagerRepoI.
func (m *managerRepo) GetByID(ctx context.Context, req *us.ManagerPrimaryKey) (*us.Manager, error) {
	resp := &us.Manager{}

	var (
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := m.db.QueryRow(ctx, `
	        SELECT id,
			login,
			fullname,
			salary,
			phone,
			password,
			branchId,
	        created_at,
	        updated_at
	        FROM "manager"
	    WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.Login, &resp.Fullname, &resp.Salary, &resp.Phone, &resp.Password, &resp.BranchId, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting manager by id", err)
		return nil, err
	}

	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

// GetList implements storage.ManagerRepoI.
func (m *managerRepo) GetList(ctx context.Context, req *us.GetListManagerRequest) (*us.GetListManagerResponse, error) {
	resp := &us.GetListManagerResponse{}
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

	rows, err := m.db.Query(ctx, `
        SELECT
            id,
            login,
            fullname,
            salary,
            phone,
            password,
            branchId,
            created_at,
            updated_at
        FROM "manager" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all managers:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var manager us.Manager
		count++
		err = rows.Scan(&manager.Id, &manager.Login, &manager.Fullname, &manager.Salary, &manager.Phone, &manager.Password, &manager.BranchId, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning managers:", err)
			return nil, err
		}
		manager.CreatedAt = created_at.String
		manager.UpdatedAt = updated_at.String

		resp.Managers = append(resp.Managers, &manager)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.ManagerRepoI.
func (m *managerRepo) Update(ctx context.Context, req *us.UpdateManager) (*us.Manager, error) {

	_, err := m.db.Exec(ctx, `
        UPDATE "manager" SET
			login=$1,
			fullname=$2,
			salary=$3,
			phone=$4,
			password=$5,
			branchId=$6,
            updated_at = NOW()
        WHERE id = $7`, req.Login, req.Fullname, req.Salary, req.Phone, req.Password, req.BranchId, req.Id)

	if err != nil {
		log.Println("error while updating manager in storage", err)
		return nil, err
	}

	manager, err := m.GetByID(ctx, &us.ManagerPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated manager by id", err)
		return nil, err
	}

	return manager, nil
}

// Delete implements storage.ManagerRepoI.
func (m *managerRepo) Delete(ctx context.Context, req *us.ManagerPrimaryKey) error {
	_, err := m.db.Exec(ctx, `
		UPDATE "manager" SET 
			deleted_at = 1
		WHERE id = $1
	`, req.Id)

	if err != nil {
		log.Println("error while deleting manager")
		return err
	}

	return nil
}

func (r *managerRepo) GetLastLogin(ctx context.Context) (string, error) {
	var login string
	err := r.db.QueryRow(ctx, `
        SELECT login FROM "manager"
        ORDER BY login DESC LIMIT 1
    `).Scan(&login)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "no rows in result set" {
			log.Println("No rows found, returning default login: M00000")
			return "M00000", nil
		}
		// Log the unexpected error
		log.Printf("Unexpected database error: %v", err)
		return "", fmt.Errorf("database query error: %w", err)
	}

	log.Printf("Retrieved login: %s", login)
	return login, nil
}

func GenerateNewLoginManager(log string) string {
	prefix := "M"
	numbStr := log[1:]
	num, _ := strconv.Atoi(numbStr)
	newNum := num + 1
	return fmt.Sprintf("%s%05d", prefix, newNum)
}

func (r *managerRepo) GetByLogin(ctx context.Context, login string) (*us.Manager, error) {
	var resp us.Manager
	var created_at, updated_at sql.NullString
	err := r.db.QueryRow(ctx, `
        SELECT 
		id,
		login, 
		fullname, 
		phone, 
		password, 
		salary, 
		branchId, 
		created_at, 
		updated_at
        FROM "manager"
        WHERE login = $1 and deleted_at=0
    `, login).Scan(&resp.Id,
		&resp.Login, &resp.Fullname, &resp.Phone, &resp.Password, &resp.Salary, &resp.BranchId,
		&created_at, &updated_at,
	)
	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	if err != nil {
		log.Println("error while getting manager by id:", err)
		return nil, err
	}
	return &resp, nil
}
