package api

import (
	"net/http"
	"user_api_gateway/api/handler"
	"user_api_gateway/config"
	"user_api_gateway/pkg/grpc_client"
	"user_api_gateway/pkg/logger"

	_ "user_api_gateway/api/docs" //for swagger

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Config ...
type Config struct {
	Logger     logger.Logger
	GrpcClient *grpc_client.GrpcClient
	Cfg        config.Config
}

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(cnf Config) *gin.Engine {
	r := gin.New()

	r.Static("/images", "./static/images")

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "*")
	// config.AllowOrigins = cnf.Cfg.AllowOrigins
	r.Use(cors.New(config))

	handler := handler.New(&handler.HandlerConfig{
		Logger:     cnf.Logger,
		GrpcClient: cnf.GrpcClient,
		Cfg:        cnf.Cfg,
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Api gateway"})
	})

	// Administration
	r.POST("/CreateAdministration", handler.CreateAdministration)
	r.GET("/GetListAdministration", handler.GetListAdministration)
	r.GET("/GetByIdAdministration/:id", handler.GetAdministrationByID)
	r.PUT("/UpdateAdministration/:id", handler.UpdateAdministrationr)
	r.DELETE("/DeleteAdministration/:id", handler.DeleteAdministration)

	// Branch
	r.POST("/CreateBranch", handler.CreateBranch)
	r.GET("/GetListBranch", handler.GetListBranch)
	r.GET("/GetByIdBranch/:id", handler.GetBranchByID)
	r.PUT("/UpdateBranch/:id", handler.UpdateBranch)
	r.DELETE("/DeleteBranch/:id", handler.DeleteBranch)

	// Manager
	r.POST("/CreateManager", handler.CreateManager)
	r.GET("/GetListManager", handler.GetListManager)
	r.GET("/GetByIdManager/:id", handler.GetManagerByID)
	r.PUT("/UpdateManager/:id", handler.UpdateManager)
	r.DELETE("/DeleteManager/:id", handler.DeleteManager)

	// Student
	r.POST("/CreateStudent", handler.CreateStudent)
	r.GET("/GetListStudent", handler.GetListStudent)
	r.GET("/GetByIdStudent/:id", handler.GetStudentByID)
	r.PUT("/UpdateStuddent/:id", handler.UpdateStudent)
	r.DELETE("/DeleteStudent/:id", handler.DeleteStudent)

	// SupportTeacher
	r.POST("/CreateSupportTeacher", handler.CreateSupportTeacher)
	r.GET("/GetListSupportTeacher", handler.GetListSupportTeacher)
	r.GET("/GetByIdSupportTeacher/:id", handler.GetSupportTeacherByID)
	r.PUT("/UpdateSupportTeacher/:id", handler.UpdateSupportTeacher)
	r.DELETE("/DeleteSupportTeacher/:id", handler.DeleteSupportTeacher)

	// Teacher
	r.POST("/CreateTeacher", handler.CreateTeacher)
	r.GET("/GetListTeacher", handler.GetListTeacher)
	r.GET("/GetByIdTeacher/:id", handler.GetTeacherByID)
	r.PUT("/UpdateTeacher/:id", handler.UpdateStudent)
	r.DELETE("/DeleteTeacher/:id", handler.DeleteStudent)

	// Login
	r.POST("/LoginAdministration", handler.AdministarationLogin)
	r.POST("/LoginManager", handler.ManagerLogin)
	r.POST("/LoginStudent", handler.StudentLogin)
	r.POST("/LoginSupportTeacher", handler.SupportTeacherLogin)
	r.POST("/LoginTeacher", handler.TeacherLogin)
	r.POST("/LoginSuperAdmin", handler.SuperAdminLogin)

	// EventStudent
	r.POST("/CreateEventStudent", handler.CreateEventStudent)
	r.GET("/GetListEventStudent", handler.GetListEventStudent)
	r.GET("/GetByIdEventStudent/:id", handler.GetEventStudentByID)
	r.PUT("/UpdateEventStudent/:id", handler.UpdateEventStudent)
	r.DELETE("/DeleteEventStudent/:id", handler.DeleteEventStudent)
	r.GET("/EventStudent/:id", handler.GetStudentWithEventsByID)

	// Event
	r.POST("/CreateEvent", handler.CreateEvent)
	r.GET("/GetListEvent", handler.GetListEvent)
	r.GET("/GetByIdEvent/:id", handler.GetEventByID)
	r.PUT("/UpdateEvent/:id", handler.UpdateEvent)
	r.DELETE("/DeleteEvent/:id", handler.DeleteEvent)

	// Group
	r.POST("/CreateGroup", handler.CreateGroup)
	r.GET("/GetListGroup", handler.GetListGroup)
	r.GET("/GetByIdGroup/:id", handler.GetGroupByID)
	r.GET("/GroupTeacher/:id", handler.GetGroupByIDTeacher)
	r.PUT("/UpdateGroup/:id", handler.UpdateGroup)
	r.DELETE("/DeleteGroup/:id", handler.DeleteGroup)

	// Journal
	r.POST("/CreateJournal", handler.CreateJournal)
	r.GET("/GetListJournal", handler.GetListJournal)
	r.GET("/GetByIdJournal/:id", handler.GetJournalByID)
	r.GET("/GetJurnalsStudent/:id", handler.GetJurnalByIDStudent)
	r.PUT("/UpdateJournal/:id", handler.UpdateJournal)
	r.DELETE("/DeleteJournal/:id", handler.DeleteJournal)

	// Schedule 
	r.POST("/CreateSchedule", handler.CreateSchedule)
	r.GET("/GetListSchedule", handler.GetListSchedule)
	r.GET("/GetByIdSchedule/:id", handler.GetScheduleByID)
	r.PUT("/UpdateSchedule/:id", handler.UpdateSchedule)
	r.DELETE("/DeleteSchedule/:id", handler.DeleteSchedule)
	r.GET("/GetScheduleForWeek", handler.GetScheduleForWeek)
	r.GET("/GetScheduleForMonth", handler.GetScheduleForMonth)

	// StudentPayment
	r.POST("/CreateStudentPayment", handler.CreateStudentPayment)
	r.GET("/GetListStudentPayment", handler.GetListStudentPayment)
	r.GET("/GetByIdStudentPayment/:id", handler.GetStudentPaymentByID)
	r.PUT("/UpdateStudentPayment/:id", handler.UpdateStudentPayment)
	r.DELETE("/DeleteStudentPayment/:id", handler.DeleteStudentPayment)

	// StudentTask 
	r.POST("/CreateStudentTask", handler.CreateStudentTask)
	r.GET("/GetListStudentTask", handler.GetListStudentTask)
	r.GET("/GetByIdStudentTask/:id", handler.GetStudentTaskByID)
	r.PUT("/UpdateStudentTask/:id", handler.UpdateStudentTask)
	r.DELETE("/DeleteStudentTask/:id", handler.DeleteStudentTask)

	// Task 
	r.POST("/CreateTask", handler.CreateTask)
	r.GET("/GetListTask", handler.GetListTask)
	r.GET("/GetByIdTask/:id", handler.GetTaskByID)
	r.PUT("/UpdateTask/:id", handler.UpdateTask)
	r.DELETE("/DeleteTask/:id", handler.DeleteTask)

	// Report
	r.GET("/AdministrationReportList", handler.GetReportListAdministration)
	r.GET("/TeacherReportList", handler.GetReportListTeacher)
	r.GET("/SupportTeacherReportList", handler.GetReportListSupportTeacher)
	r.GET("/StudentReportList", handler.GetReportListStudent)

	// Shipper endpoints
	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return r
}
