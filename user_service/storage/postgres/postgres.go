package postgres

import (
	"context"
	"fmt"
	"log"
	"user_service/config"
	"user_service/storage"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db             *pgxpool.Pool
	administration storage.AdministrationRepoI
	branch         storage.BranchRepoI
	manager        storage.ManagerRepoI
	student        storage.StudentRepoI
	supportTeacher storage.SupportTeacherRepoI
	teacher        storage.TeacherRepoI
}

func NewPostgres(ctx context.Context, cfg config.Config) (storage.StorageI, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	))
	if err != nil {
		return nil, err
	}

	config.MaxConns = cfg.PostgresMaxConnections

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: pool,
	}, err
}

func (s *Store) CloseDB() {
	s.db.Close()
}

func (l *Store) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	args := make([]interface{}, 0, len(data)+2) // making space for arguments + level + msg
	args = append(args, level, msg)
	for k, v := range data {
		args = append(args, fmt.Sprintf("%s=%v", k, v))
	}
	log.Println(args...)
}

func (s *Store) Administration() storage.AdministrationRepoI {
	if s.administration == nil {
		s.administration = NewAdministrationRepo(s.db)
	}

	return s.administration
}

// Branch implements storage.StorageI.
func (s *Store) Branch() storage.BranchRepoI {
	if s.branch == nil {
		s.branch = NewBranchRepo(s.db)
	}

	return s.branch
}

// Manager implements storage.StorageI.
func (s *Store) Manager() storage.ManagerRepoI {
	if s.manager == nil {
		s.manager = NewManagerRepo(s.db)
	}

	return s.manager
}

// Student implements storage.StorageI.
func (s *Store) Student() storage.StudentRepoI {
	if s.student == nil {
		s.student = NewStudentRepo(s.db)
	}

	return s.student
}

// SupportTeacher implements storage.StorageI.
func (s *Store) SupportTeacher() storage.SupportTeacherRepoI {
	if s.supportTeacher == nil {
		s.supportTeacher = NewSupportTeacherRepo(s.db)
	}

	return s.supportTeacher
}

// Teacher implements storage.StorageI.
func (s *Store) Teacher() storage.TeacherRepoI {
	if s.teacher == nil {
		s.teacher = NewTeacherRepo(s.db)
	}

	return s.teacher
}
