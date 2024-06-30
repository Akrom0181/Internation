package storage

import (
	"context"
	us "user_service/genproto/user_service"
)

type StorageI interface {
	CloseDB()
	Administration() AdministrationRepoI
	Branch() BranchRepoI
	Manager() ManagerRepoI
	Student() StudentRepoI
	SupportTeacher() SupportTeacherRepoI
	Teacher() TeacherRepoI
}

type AdministrationRepoI interface {
	Create(ctx context.Context, req *us.CreateAdministration) (*us.Administration, error)
	GetByID(ctx context.Context, req *us.AdministrationPrimaryKey) (*us.Administration, error)
	GetList(ctx context.Context, req *us.GetListAdministrationRequest) (*us.GetListAdministrationResponse, error)
	Update(ctx context.Context, req *us.UpdateAdministration) (*us.Administration, error)
	Delete(ctx context.Context, req *us.AdministrationPrimaryKey) error
	GetByLogin(ctx context.Context, login string) (*us.Administration, error)
	GetReportList(ctx context.Context, req *us.GetReportListAdministrationRequest) (*us.GetReportListAdministrationResponse, error)
}

type BranchRepoI interface {
	Create(ctx context.Context, req *us.CreateBranch) (*us.Branch, error)
	GetByID(ctx context.Context, req *us.BranchPrimaryKey) (*us.Branch, error)
	GetList(ctx context.Context, req *us.GetListBranchRequest) (*us.GetListBranchResponse, error)
	Update(ctx context.Context, req *us.UpdateBranch) (*us.Branch, error)
	Delete(ctx context.Context, req *us.BranchPrimaryKey) error
}

type ManagerRepoI interface {
	Create(ctx context.Context, req *us.CreateManager) (*us.Manager, error)
	GetByID(ctx context.Context, req *us.ManagerPrimaryKey) (*us.Manager, error)
	GetList(ctx context.Context, req *us.GetListManagerRequest) (*us.GetListManagerResponse, error)
	Update(ctx context.Context, req *us.UpdateManager) (*us.Manager, error)
	Delete(ctx context.Context, req *us.ManagerPrimaryKey) error
	GetByLogin(ctx context.Context, login string) (*us.Manager, error)
}

type StudentRepoI interface {
	Create(ctx context.Context, req *us.CreateStudent) (*us.Student, error)
	GetByID(ctx context.Context, req *us.StudentPrimaryKey) (*us.Student, error)
	GetList(ctx context.Context, req *us.GetListStudentRequest) (*us.GetListStudentResponse, error)
	Update(ctx context.Context, req *us.UpdateStudent) (*us.Student, error)
	Delete(ctx context.Context, req *us.StudentPrimaryKey) error
	GetByLogin(ctx context.Context, login string) (*us.Student, error)
	GetReportList(ctx context.Context, req *us.GetReportListStudentRequest) (*us.GetReportListStudentResponse, error)
}

type SupportTeacherRepoI interface {
	Create(ctx context.Context, req *us.CreateSupportTeacher) (*us.SupportTeacher, error)
	GetByID(ctx context.Context, req *us.SupportTeacherPrimaryKey) (*us.SupportTeacher, error)
	GetList(ctx context.Context, req *us.GetListSupportTeacherRequest) (*us.GetListSupportTeacherResponse, error)
	Update(ctx context.Context, req *us.UpdateSupportTeacher) (*us.SupportTeacher, error)
	Delete(ctx context.Context, req *us.SupportTeacherPrimaryKey) error
	GetByLogin(ctx context.Context, login string) (*us.SupportTeacher, error)
	GetReportList(ctx context.Context, req *us.GetReportListSupportTeacherRequest) (*us.GetReportListSupportTeacherResponse, error)
}

type TeacherRepoI interface {
	Create(ctx context.Context, req *us.CreateTeacher) (*us.Teacher, error)
	GetByID(ctx context.Context, req *us.TeacherPrimaryKey) (*us.Teacher, error)
	GetList(ctx context.Context, req *us.GetListTeacherRequest) (*us.GetListTeacherResponse, error)
	Update(ctx context.Context, req *us.UpdateTeacher) (*us.Teacher, error)
	Delete(ctx context.Context, req *us.TeacherPrimaryKey) error
	GetByLogin(ctx context.Context, login string) (*us.Teacher, error)
	GetReportList(ctx context.Context, req *us.GetReportListTeacherRequest) (*us.GetReportListTeacherResponse, error)
}
