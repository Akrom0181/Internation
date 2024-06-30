package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	us "user_service/genproto/user_service"
	"user_service/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type branchRepo struct {
	db *pgxpool.Pool
}

func NewBranchRepo(db *pgxpool.Pool) storage.BranchRepoI {
	return &branchRepo{
		db: db,
	}
}

// Create implements storage.BranchRepoI.
func (b *branchRepo) Create(ctx context.Context, req *us.CreateBranch) (*us.Branch, error) {
	id := uuid.NewString()

	_, err := b.db.Exec(ctx, `
		INSERT INTO "branch" (
			id,
			name,
			address,
			phone
		) VALUES (
			$1, $2, $3, $4
		)`, id, req.Name, req.Address, req.Phone)

	if err != nil {
		log.Println("error while creating branch in storage", err)
		return nil, err
	}

	branch, err := b.GetByID(ctx, &us.BranchPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting branch by id after creating", err)
		return nil, err
	}
	return branch, nil
}

// GetByID implements storage.BranchRepoI.
func (b *branchRepo) GetByID(ctx context.Context, req *us.BranchPrimaryKey) (*us.Branch, error) {
	resp := &us.Branch{}

	var (
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := b.db.QueryRow(ctx, `
	        SELECT id,
			name,
			address,
			phone,
	        created_at,
	        updated_at
	        FROM "branch"
	    WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.Name, &resp.Address, &resp.Phone, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting branch by id", err)
		return nil, err
	}

	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

// GetList implements storage.BranchRepoI.
func (b *branchRepo) GetList(ctx context.Context, req *us.GetListBranchRequest) (*us.GetListBranchResponse, error) {
	resp := &us.GetListBranchResponse{}
	var (
		filter     string
		created_at sql.NullString
		updated_at sql.NullString
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` AND (name ILIKE '%%%v%%' OR address ILIKE '%%%v%%' OR phone ILIKE '%%%v%%')`, req.Search, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := b.db.Query(ctx, `
        SELECT
            id,
            name,
            address,
            phone,
            created_at,
            updated_at
        FROM "branch" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all branches:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var branch us.Branch
		count++
		err = rows.Scan(&branch.Id, &branch.Name, &branch.Address, &branch.Phone, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning branches:", err)
			return nil, err
		}
		branch.CreatedAt = created_at.String
		branch.UpdatedAt = updated_at.String

		resp.Branches = append(resp.Branches, &branch)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.BranchRepoI.
func (b *branchRepo) Update(ctx context.Context, req *us.UpdateBranch) (*us.Branch, error) {

	_, err := b.db.Exec(ctx, `
        UPDATE "branch" SET
			name=$1,
			address=$2,
			phone=$3,
            updated_at = NOW()
        WHERE id = $4`, req.Name, req.Address, req.Phone, req.Id)

	if err != nil {
		log.Println("error while updating branch in storage", err)
		return nil, err
	}

	branch, err := b.GetByID(ctx, &us.BranchPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated branch by id", err)
		return nil, err
	}

	return branch, nil
}

// Delete implements storage.BranchRepoI.
func (b *branchRepo) Delete(ctx context.Context, req *us.BranchPrimaryKey) error {
	_, err := b.db.Exec(ctx, `
		UPDATE "branch" SET 
			deleted_at = 1
		WHERE id = $1
	`, req.Id)

	if err != nil {
		log.Println("error while deleting branch")
		return err
	}

	return nil
}
