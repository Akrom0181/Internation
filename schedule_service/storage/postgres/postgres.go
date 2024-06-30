package postgres

import (
	"context"
	"fmt"
	"log"
	"schedule_service/config"
	"schedule_service/storage"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db             *pgxpool.Pool
	event_student  storage.EventStudentRepoI
	event          storage.EventRepoI
	group          storage.GroupRepoI
	journal        storage.JournalRepoI
	schedule       storage.ScheduleRepoI
	studentTask    storage.StudentTaskRepoI
	task           storage.TaskRepoI
	studentPayment storage.StudentPaymentRepoI
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

func (s *Store) EventStudent() storage.EventStudentRepoI {
	if s.event_student == nil {
		s.event_student = NewEventStudentRepo(s.db)
	}

	return s.event_student
}

// Event implements storage.StorageI.
func (s *Store) Event() storage.EventRepoI {
	if s.event == nil {
		s.event = NewEventRepo(s.db)
	}

	return s.event
}

// Group implements storage.StorageI.
func (s *Store) Group() storage.GroupRepoI {
	if s.group == nil {
		s.group = NewGroupRepo(s.db)
	}

	return s.group
}

// Journal implements storage.StorageI.
func (s *Store) Journal() storage.JournalRepoI {
	if s.journal == nil {
		s.journal = NewJournalRepo(s.db)
	}

	return s.journal
}

// Schedule implements storage.StorageI.
func (s *Store) Schedule() storage.ScheduleRepoI {
	if s.schedule == nil {
		s.schedule = NewScheduleRepo(s.db)
	}

	return s.schedule
}

// StudentTask implements storage.StorageI.
func (s *Store) StudentTask() storage.StudentTaskRepoI {
	if s.studentTask == nil {
		s.studentTask = NewStudentTaskRepo(s.db)
	}

	return s.studentTask
}

// Task implements storage.StorageI.
func (s *Store) Task() storage.TaskRepoI {
	if s.task == nil {
		s.task = NewTaskRepo(s.db)
	}

	return s.task
}

// StudentPayment implements storage.StorageI.
func (s *Store) StudentPayment() storage.StudentPaymentRepoI {
	if s.studentPayment == nil {
		s.studentPayment = NewStudentPaymentRepo(s.db)
	}

	return s.studentPayment
}
