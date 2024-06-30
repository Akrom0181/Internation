package grpc_client

import (
	"fmt"
	"log"
	sc "user_api_gateway/genproto/schedule_service"
	pc "user_api_gateway/genproto/user_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"user_api_gateway/config"
)

// GrpcClientI ..90.
type GrpcClientI interface {
	AdministrationService() pc.AdministrationServiceClient
	BranchService() pc.BranchServiceClient
	ManagerService() pc.ManagerServiceClient
	StudentService() pc.StudentServiceClient
	SupportTeacherService() pc.SupportTeacherServiceClient
	TeacherService() pc.TeacherServiceClient
	LoginService() pc.LoginServiceClient
	EventStudentService() sc.EventServiceClient
	EventService() sc.EventServiceClient
	GroupService() sc.GroupServiceClient
	JournalService() sc.JournalServiceClient
	ScheduleService() sc.ScheduleServiceClient
	StudentPaymentService() sc.StudentPaymentServiceClient
	StudentTaskService() sc.StudentTaskServiceClient
	TaskService() sc.TaskServiceClient
}

// GrpcClient ...
type GrpcClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

// New ...
func New(cfg config.Config) (*GrpcClient, error) {

	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("user service dial host: %s port:%s err: %s",
			cfg.UserServiceHost, cfg.UserServicePort, err)
	}

	connSchedule, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.ScheduleServiceHost, cfg.ScheduleServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("user service dial host: %s port:%s err: %s",
			cfg.UserServiceHost, cfg.UserServicePort, err)
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"administration_service": pc.NewAdministrationServiceClient(connUser),
			"branch_service":         pc.NewBranchServiceClient(connUser),
			"manager_service":        pc.NewManagerServiceClient(connUser),
			"student_service":        pc.NewStudentServiceClient(connUser),
			"supportTeacher_service": pc.NewSupportTeacherServiceClient(connUser),
			"teacher_service":        pc.NewTeacherServiceClient(connUser),
			"login_service":          pc.NewLoginServiceClient(connUser),
			"event_student":          sc.NewEventStudentServiceClient(connSchedule),
			"event":                  sc.NewEventServiceClient(connSchedule),
			"group":                  sc.NewGroupServiceClient(connSchedule),
			"journal":                sc.NewJournalServiceClient(connSchedule),
			"schedule":               sc.NewScheduleServiceClient(connSchedule),
			"student_payment":        sc.NewStudentPaymentServiceClient(connSchedule),
			"student_task":           sc.NewStudentTaskServiceClient(connSchedule),
			"task":                   sc.NewTaskServiceClient(connSchedule),
		},
	}, nil
}

// AdministrationService returns the AdministrationServiceClient
func (g *GrpcClient) AdministrationService() pc.AdministrationServiceClient {
	client, ok := g.connections["administration_service"].(pc.AdministrationServiceClient)
	if !ok {
		log.Println("failed to assert type for administration")
		return nil
	}
	return client
}

// BranchService returns the BranchServiceClient
func (g *GrpcClient) BranchService() pc.BranchServiceClient {
	client, ok := g.connections["branch_service"].(pc.BranchServiceClient)
	if !ok {
		log.Println("failed to assert type for branch")
		return nil
	}
	return client
}

// ManagerService returns the ManagerServiceClient
func (g *GrpcClient) ManagerService() pc.ManagerServiceClient {
	client, ok := g.connections["manager_service"].(pc.ManagerServiceClient)
	if !ok {
		log.Println("failed to assert type for manager")
		return nil
	}
	return client
}

// StudentService returns the StudentServiceClient
func (g *GrpcClient) StudentService() pc.StudentServiceClient {
	client, ok := g.connections["student_service"].(pc.StudentServiceClient)
	if !ok {
		log.Println("failed to assert type for student")
		return nil
	}
	return client
}

func (g *GrpcClient) SupportTeacherService() pc.SupportTeacherServiceClient {
	client, ok := g.connections["supportTeacher_service"].(pc.SupportTeacherServiceClient)
	if !ok {
		log.Println("failed to assert type for supportTeacher")
		return nil
	}
	return client
}

func (g *GrpcClient) TeacherService() pc.TeacherServiceClient {
	client, ok := g.connections["teacher_service"].(pc.TeacherServiceClient)
	if !ok {
		log.Println("failed to assert type for teacher")
		return nil
	}
	return client
}

func (g *GrpcClient) LoginService() pc.LoginServiceClient {
	client, ok := g.connections["login_service"].(pc.LoginServiceClient)
	if !ok {
		log.Println("failed to assert type for login")
		return nil
	}
	return client
}

// AdministrationService returns the AdministrationServiceClient
func (g *GrpcClient) EventStudentService() sc.EventStudentServiceClient {
	client, ok := g.connections["event_student"].(sc.EventStudentServiceClient)
	if !ok {
		log.Println("failed to assert type for administration")
		return nil
	}
	return client
}

// BranchService returns the BranchServiceClient
func (g *GrpcClient) EventService() sc.EventServiceClient {
	client, ok := g.connections["event"].(sc.EventServiceClient)
	if !ok {
		log.Println("failed to assert type for event")
		return nil
	}
	return client
}

// ManagerService returns the ManagerServiceClient
func (g *GrpcClient) GroupService() sc.GroupServiceClient {
	client, ok := g.connections["group"].(sc.GroupServiceClient)
	if !ok {
		log.Println("failed to assert type for group")
		return nil
	}
	return client
}

// StudentService returns the StudentServiceClient
func (g *GrpcClient) JournalService() sc.JournalServiceClient {
	client, ok := g.connections["journal"].(sc.JournalServiceClient)
	if !ok {
		log.Println("failed to assert type for journal")
		return nil
	}
	return client
}

func (g *GrpcClient) ScheduleService() sc.ScheduleServiceClient {
	client, ok := g.connections["schedule"].(sc.ScheduleServiceClient)
	if !ok {
		log.Println("failed to assert type for supportTeacher")
		return nil
	}
	return client
}

func (g *GrpcClient) StudentPaymentService() sc.StudentPaymentServiceClient {
	client, ok := g.connections["student_payment"].(sc.StudentPaymentServiceClient)
	if !ok {
		log.Println("failed to assert type for teacher")
		return nil
	}
	return client
}

func (g *GrpcClient) StudentTaskService() sc.StudentTaskServiceClient {
	client, ok := g.connections["student_task"].(sc.StudentTaskServiceClient)
	if !ok {
		log.Println("failed to assert type for student_task")
		return nil
	}
	return client
}

func (g *GrpcClient) TaskService() sc.TaskServiceClient {
	client, ok := g.connections["task"].(sc.TaskServiceClient)
	if !ok {
		log.Println("failed to assert type for task")
		return nil
	}
	return client
}

func (g *GrpcClient) CloseConnections() {
	for key, conn := range g.connections {
		if c, ok := conn.(*grpc.ClientConn); ok {
			err := c.Close()
			if err != nil {
				log.Printf("failed to close connection for %s: %v", key, err)
			}
		}
	}
}
