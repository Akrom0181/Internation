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

type groupRepo struct {
	db *pgxpool.Pool
}

func NewGroupRepo(db *pgxpool.Pool) storage.GroupRepoI {
	return &groupRepo{
		db: db,
	}
}

// Create implements storage.GroupRepoI.
func (g *groupRepo) Create(ctx context.Context, req *schedule_service.CreateGroup) (*schedule_service.GetGroup, error) {
	id := uuid.NewString()

	_, err := g.db.Exec(ctx, `
        INSERT INTO "group" (
            id,
            teacherId,
            supportTeacherId,
            branchId,
            type
        ) VALUES (
            $1, $2, $3, $4, $5
        )`, id, req.TeacherId, req.SuppportTeacherId, req.BranchId, req.Type)

	if err != nil {
		log.Println("error while creating group in storage", err)
		return nil, err
	}

	group, err := g.GetByID(ctx, &schedule_service.GroupPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting group by id after creating", err)
		return nil, err
	}
	return group, nil
}

// GetByID implements storage.GroupRepoI.
func (g *groupRepo) GetByID(ctx context.Context, req *schedule_service.GroupPrimaryKey) (*schedule_service.GetGroup, error) {
	resp := &schedule_service.GetGroup{}

	var (
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := g.db.QueryRow(ctx, `
            SELECT id,
            teacherId,
            supportTeacherId,
            branchId,
            type,
            created_at,
            updated_at
            FROM "group"
        WHERE id=$1`, req.Id).Scan(&resp.Id, &resp.TeacherId, &resp.SuppportTeacherId, &resp.BranchId, &resp.Type, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting group by id", err)
		return nil, err
	}

	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

// GetList implements storage.GroupRepoI.
func (g *groupRepo) GetList(ctx context.Context, req *schedule_service.GetListGroupRequest) (*schedule_service.GetListGroupResponse, error) {
	resp := &schedule_service.GetListGroupResponse{}
	var (
		filter     string
		created_at sql.NullString
		updated_at sql.NullString
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` AND (type ILIKE '%%%v%%')`, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := g.db.Query(ctx, `
        SELECT
            id,
            teacherId,
            supportTeacherId,
            branchId,
            type,
            created_at,
            updated_at
        FROM "group" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all groups:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var group schedule_service.Group
		count++
		err = rows.Scan(&group.Id, &group.TeacherId, &group.SuppportTeacherId, &group.BranchId, &group.Type, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning groups:", err)
			return nil, err
		}
		group.CreatedAt = created_at.String
		group.UpdatedAt = updated_at.String

		resp.Groups = append(resp.Groups, &group)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.GroupRepoI.
func (g *groupRepo) Update(ctx context.Context, req *schedule_service.UpdateGroup) (*schedule_service.GetGroup, error) {

	_, err := g.db.Exec(ctx, `
        UPDATE "group" SET
            teacherId=$1,
            supportTeacherId=$2,
            branchId=$3,
            type=$4,
            updated_at = NOW()
        WHERE id = $5`, req.TeacherId, req.SuppportTeacherId, req.BranchId, req.Type, req.Id)

	if err != nil {
		log.Println("error while updating group in storage", err)
		return nil, err
	}

	group, err := g.GetByID(ctx, &schedule_service.GroupPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated group by id", err)
		return nil, err
	}

	return group, nil
}

// Delete implements storage.GroupRepoI.
func (g *groupRepo) Delete(ctx context.Context, req *schedule_service.GroupPrimaryKey) error {
	_, err := g.db.Exec(ctx, `
        UPDATE "group" SET 
            deleted_at = 1
        WHERE id = $1
    `, req.Id)

	if err != nil {
		log.Println("error while deleting group")
		return err
	}

	return nil
}

func (r *groupRepo) GetByIDTeacher(ctx context.Context, req *schedule_service.TeacherID) (*schedule_service.GetGroup, error) {
	var resp schedule_service.GetGroup
	var created_at, updated_at sql.NullString

	err := r.db.QueryRow(ctx, `
        SELECT id, teacherId, supportTeacherId, branchId, type, created_at, updated_at
        FROM "group"
        WHERE teacherId = $1
    `, req.Id).Scan(
		&resp.Id, &resp.TeacherId, &resp.SuppportTeacherId, &resp.BranchId, &resp.Type,
		&created_at, &updated_at,
	)
	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String
	if err != nil {
		log.Println("error while getting group by id:", err)
		return nil, err
	}
	return &resp, nil
}
