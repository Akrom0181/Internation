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

type journalRepo struct {
	db *pgxpool.Pool
}

func NewJournalRepo(db *pgxpool.Pool) storage.JournalRepoI {
	return &journalRepo{
		db: db,
	}
}

// Create implements storage.JournalRepoI.
func (j *journalRepo) Create(ctx context.Context, req *schedule_service.CreateJournal) (*schedule_service.GetJournal, error) {
	id := uuid.NewString()

	_, err := j.db.Exec(ctx, `
        INSERT INTO "journal" (
            id,
            fromDate,
            toDate,
            groupId,
            studentsCount
        ) VALUES (
            $1, $2, $3, $4, $5
        )`, id, req.FromDate, req.ToDate, req.GroupId, req.StudentsCount)

	if err != nil {
		log.Println("error while creating journal in storage", err)
		return nil, err
	}

	journal, err := j.GetByID(ctx, &schedule_service.JournalPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting journal by id after creating", err)
		return nil, err
	}
	return journal, nil
}

// GetByID implements storage.JournalRepoI.
func (j *journalRepo) GetByID(ctx context.Context, req *schedule_service.JournalPrimaryKey) (*schedule_service.GetJournal, error) {
	resp := &schedule_service.GetJournal{}

	var (
		created_at sql.NullString
		updated_at sql.NullString
		fromDate   sql.NullString
		toDate     sql.NullString
	)

	err := j.db.QueryRow(ctx, `
            SELECT id,
            fromDate,
            toDate,
            groupId,
            studentsCount,
            created_at,
            updated_at
            FROM "journal"
        WHERE id=$1`, req.Id).Scan(&resp.Id, &fromDate, &toDate, &resp.GroupId, &resp.StudentsCount, &created_at, &updated_at)

	if err != nil {
		log.Println("error while getting journal by id", err)
		return nil, err
	}

	resp.FromDate = fromDate.String
	resp.ToDate = toDate.String
	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	return resp, nil
}

// GetList implements storage.JournalRepoI.
func (j *journalRepo) GetList(ctx context.Context, req *schedule_service.GetListJournalRequest) (*schedule_service.GetListJournalResponse, error) {
	resp := &schedule_service.GetListJournalResponse{}
	var (
		filter     string
		fromDate   sql.NullString
		toDate     sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter = fmt.Sprintf(` AND (fromDate ILIKE '%%%v%%' OR toDate ILIKE '%%%v%%' OR groupId::text ILIKE '%%%v%%' OR studentsCount::text ILIKE '%%%v%%')`, req.Search, req.Search, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	rows, err := j.db.Query(ctx, `
        SELECT
            id,
            fromDate,
            toDate,
            groupId,
            studentsCount,
            created_at,
            updated_at
        FROM "journal" WHERE deleted_at=0
    `+filter)

	if err != nil {
		log.Println("error while getting all journals:", err)
		return nil, err
	}

	defer rows.Close()

	var count int64

	for rows.Next() {
		var journal schedule_service.Journal
		count++
		err = rows.Scan(&journal.Id, &fromDate, &toDate, &journal.GroupId, &journal.StudentsCount, &created_at, &updated_at)

		if err != nil {
			log.Println("error while scanning journals:", err)
			return nil, err
		}

		journal.FromDate = fromDate.String
		journal.ToDate = toDate.String
		journal.CreatedAt = created_at.String
		journal.UpdatedAt = updated_at.String

		resp.Journals = append(resp.Journals, &journal)
	}

	if err = rows.Err(); err != nil {
		log.Println("rows iteration error:", err)
		return nil, err
	}

	resp.Count = count

	return resp, nil
}

// Update implements storage.JournalRepoI.
func (j *journalRepo) Update(ctx context.Context, req *schedule_service.UpdateJournal) (*schedule_service.GetJournal, error) {

	_, err := j.db.Exec(ctx, `
        UPDATE "journal" SET
            fromDate=$1,
            toDate=$2,
            groupId=$3,
            studentsCount=$4,
            updated_at = NOW()
        WHERE id = $5`, req.FromDate, req.ToDate, req.GroupId, req.StudentsCount, req.Id)

	if err != nil {
		log.Println("error while updating journal in storage", err)
		return nil, err
	}

	journal, err := j.GetByID(ctx, &schedule_service.JournalPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting updated journal by id", err)
		return nil, err
	}

	return journal, nil
}

// Delete implements storage.JournalRepoI.
func (j *journalRepo) Delete(ctx context.Context, req *schedule_service.JournalPrimaryKey) error {
	_, err := j.db.Exec(ctx, `
        UPDATE "journal" SET 
            deleted_at = 1
        WHERE id = $1
    `, req.Id)

	if err != nil {
		log.Println("error while deleting journal")
		return err
	}

	return nil
}

// Get implements storage.JournalRepozI
func (r *journalRepo) GetByGroupID(ctx context.Context, req *schedule_service.JournalPrimaryKey) (*schedule_service.GetJournal, error) {
	var resp schedule_service.GetJournal
	var created_at, updated_at sql.NullString

	err := r.db.QueryRow(ctx, `
        SELECT 
			id,
			groupId, 
			fromDate, 
			toDate, 
			studentsCount, 
			created_at, 
			updated_at
        FROM "journal"
        WHERE groupId = $1
    `, req.Id).Scan(
		&resp.Id, &resp.GroupId, &resp.FromDate, &resp.ToDate, &resp.StudentsCount,
		&created_at, &updated_at,
	)

	resp.CreatedAt = created_at.String
	resp.UpdatedAt = updated_at.String

	if err != nil {
		log.Println("error while getting journal by id:", err)
		return nil, err
	}
	return &resp, nil
}
