package storage

import (
	"context"
	us "schedule_service/genproto/schedule_service"
)

type StorageI interface {
	CloseDB()
	EventStudent() EventStudentRepoI
	Event() EventRepoI
	Group() GroupRepoI
	Journal() JournalRepoI
	Schedule() ScheduleRepoI
	StudentTask() StudentTaskRepoI
	Task() TaskRepoI
	StudentPayment() StudentPaymentRepoI
}

type EventStudentRepoI interface {
	Create(ctx context.Context, req *us.CreateEventStudent) (*us.GetEventStudent, error)
	GetByID(ctx context.Context, req *us.EventStudentPrimaryKey) (*us.GetEventStudent, error)
	GetList(ctx context.Context, req *us.GetListEventStudentRequest) (*us.GetListEventStudentResponse, error)
	Update(ctx context.Context, req *us.UpdateEventStudent) (*us.GetEventStudent, error)
	Delete(ctx context.Context, req *us.EventStudentPrimaryKey) error
    GetStudentWithEventsByID(ctx context.Context, req *us.EventStudentPrimaryKey) (*us.GetStudentWithEventsResponse, error)

}

type EventRepoI interface {
	Create(ctx context.Context, req *us.CreateEvent) (*us.GetEvent, error)
	GetByID(ctx context.Context, req *us.EventPrimaryKey) (*us.GetEvent, error)
	GetList(ctx context.Context, req *us.GetListEventRequest) (*us.GetListEventResponse, error)
	Update(ctx context.Context, req *us.UpdateEvent) (*us.GetEvent, error)
	Delete(ctx context.Context, req *us.EventPrimaryKey) error
}

type GroupRepoI interface {
	Create(ctx context.Context, req *us.CreateGroup) (*us.GetGroup, error)
	GetByID(ctx context.Context, req *us.GroupPrimaryKey) (*us.GetGroup, error)
	GetList(ctx context.Context, req *us.GetListGroupRequest) (*us.GetListGroupResponse, error)
	Update(ctx context.Context, req *us.UpdateGroup) (*us.GetGroup, error)
	Delete(ctx context.Context, req *us.GroupPrimaryKey) error
	GetByIDTeacher(ctx context.Context, req *us.TeacherID) (*us.GetGroup, error)
}

type JournalRepoI interface {
	Create(ctx context.Context, req *us.CreateJournal) (*us.GetJournal, error)
	GetByID(ctx context.Context, req *us.JournalPrimaryKey) (*us.GetJournal, error)
	GetList(ctx context.Context, req *us.GetListJournalRequest) (*us.GetListJournalResponse, error)
	Update(ctx context.Context, req *us.UpdateJournal) (*us.GetJournal, error)
	Delete(ctx context.Context, req *us.JournalPrimaryKey) error
	GetByGroupID(ctx context.Context, req *us.JournalPrimaryKey) (*us.GetJournal, error)
}

type ScheduleRepoI interface {
	Create(ctx context.Context, req *us.CreateSchedule) (*us.GetSchedule, error)
	GetByID(ctx context.Context, req *us.SchedulePrimaryKey) (*us.GetSchedule, error)
	GetList(ctx context.Context, req *us.GetListScheduleRequest) (*us.GetListScheduleResponse, error)
	Update(ctx context.Context, req *us.UpdateSchedule) (*us.GetSchedule, error)
	Delete(ctx context.Context, req *us.SchedulePrimaryKey) error
	GetScheduleForWeek(ctx context.Context, teacherId string, weekStartDate, weekEndDate string) (*us.GetListScheduleResponse, error)
	GetScheduleForMonth(ctx context.Context, teacherId string, monthStartDate, monthEndDate string) (*us.GetListScheduleResponse, error)
}

type StudentTaskRepoI interface {
	Create(ctx context.Context, req *us.CreateStudentTask) (*us.GetStudentTask, error)
	GetByID(ctx context.Context, req *us.StudentTaskPrimaryKey) (*us.GetStudentTask, error)
	GetList(ctx context.Context, req *us.GetListStudentTaskRequest) (*us.GetListStudentTaskResponse, error)
	Update(ctx context.Context, req *us.UpdateStudentTask) (*us.GetStudentTask, error)
	Delete(ctx context.Context, req *us.StudentTaskPrimaryKey) error
	UpdateScoreforTeacher(ctx context.Context, req *us.UpdateStudentScoreRequest) (*us.GetStudentTask, error)
	UpdateScoreforStudent(ctx context.Context, req *us.UpdateStudentScoreRequest) (*us.GetStudentTask, error)
}

type TaskRepoI interface {
	Create(ctx context.Context, req *us.CreateTask) (*us.GetTask, error)
	GetByID(ctx context.Context, req *us.TaskPrimaryKey) (*us.GetTask, error)
	GetList(ctx context.Context, req *us.GetListTaskRequest) (*us.GetListTaskResponse, error)
	Update(ctx context.Context, req *us.UpdateTask) (*us.GetTask, error)
	Delete(ctx context.Context, req *us.TaskPrimaryKey) error
}

type StudentPaymentRepoI interface {
	Create(ctx context.Context, req *us.CreateStudentPayment) (*us.GetStudentPayment, error)
	GetByID(ctx context.Context, req *us.StudentPaymentPrimaryKey) (*us.GetStudentPayment, error)
	GetList(ctx context.Context, req *us.GetListStudentPaymentRequest) (*us.GetListStudentPaymentResponse, error)
	Update(ctx context.Context, req *us.UpdateStudentPayment) (*us.GetStudentPayment, error)
	Delete(ctx context.Context, req *us.StudentPaymentPrimaryKey) error
}
