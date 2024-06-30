package service

import (
	"context"
	"errors"
	"user_service/config"
	"user_service/genproto/user_service"
	"user_service/grpc/client"
	"user_service/pkg/jwt"
	"user_service/pkg/password"
	"user_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type LoginService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewLoginService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *LoginService {
	return &LoginService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}
					   
func (s *LoginService) AdministarationLogin(ctx context.Context, req *user_service.LoginPasswors) (*user_service.Token, error) {
	s.log.Info("---LoginAdministration--->>>", logger.Any("req", req))

	resp, err := s.strg.Administration().GetByLogin(ctx, req.Login)
	if err != nil {
		s.log.Error("---LoginAdministration--->>>", logger.Error(err))
		return &user_service.Token{}, err
	}

	if err = password.CompareHashAndPassword(resp.Password, req.Password); err != nil {
		s.log.Error("error while comparing password", logger.Error(err))
		return &user_service.Token{}, err
	}

	m := make(map[interface{}]interface{})

	m["user_id"] = resp.Id
	m["user_role"] = "Administration"

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		s.log.Error("error while generating tokens for User login", logger.Error(err))
		return &user_service.Token{}, err
	}

	token := user_service.Token{}
	token.AccessToken = accessToken
	token.RefreshToken = refreshToken

	return &token, nil
}

func (s *LoginService) ManagerLogin(ctx context.Context, req *user_service.LoginPasswors) (*user_service.Token, error) {
	s.log.Info("---LoginManager--->>>", logger.Any("req", req))

	resp, err := s.strg.Manager().GetByLogin(ctx, req.Login)
	if err != nil {
		s.log.Error("---LoginManager--->>>", logger.Error(err))
		return &user_service.Token{}, err
	}

	if err = password.CompareHashAndPassword(resp.Password, req.Password); err != nil {
		s.log.Error("error while comparing password", logger.Error(err))
		return &user_service.Token{}, err
	}

	////
	m := make(map[interface{}]interface{})

	m["user_id"] = resp.Id
	m["user_role"] = "Manager"

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		s.log.Error("error while generating tokens for User login", logger.Error(err))
		return &user_service.Token{}, err
	}

	token := user_service.Token{}
	token.AccessToken = accessToken
	token.RefreshToken = refreshToken

	return &token, nil
}

func (s *LoginService) StudentLogin(ctx context.Context, req *user_service.LoginPasswors) (*user_service.Token, error) {
	s.log.Info("---LoginStudent--->>>", logger.Any("req", req))

	resp, err := s.strg.Student().GetByLogin(ctx, req.Login)
	if err != nil {
		s.log.Error("---LoginStudent--->>>", logger.Error(err))
		return &user_service.Token{}, err
	}

	if err = password.CompareHashAndPassword(resp.Password, req.Password); err != nil {
		s.log.Error("error while comparing password", logger.Error(err))
		return &user_service.Token{}, err
	}

	m := make(map[interface{}]interface{})

	m["user_id"] = resp.Id
	m["user_role"] = "Student"

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		s.log.Error("error while generating tokens for User login", logger.Error(err))
		return &user_service.Token{}, err
	}

	token := user_service.Token{}
	token.AccessToken = accessToken
	token.RefreshToken = refreshToken

	return &token, nil
}

func (s *LoginService) SupportTeacherLogin(ctx context.Context, req *user_service.LoginPasswors) (*user_service.Token, error) {
	s.log.Info("---LoginSupportTeacher--->>>", logger.Any("req", req))

	resp, err := s.strg.SupportTeacher().GetByLogin(ctx, req.Login)
	if err != nil {
		s.log.Error("---LoginSupportTeacher--->>>", logger.Error(err))
		return &user_service.Token{}, err
	}
	if err = password.CompareHashAndPassword(resp.Password, req.Password); err != nil {
		s.log.Error("error while comparing password", logger.Error(err))
		return &user_service.Token{}, err
	}

	////
	m := make(map[interface{}]interface{})

	m["user_id"] = resp.Id
	m["user_role"] = "SupportTeacher"

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		s.log.Error("error while generating tokens for User login", logger.Error(err))
		return &user_service.Token{}, err
	}

	token := user_service.Token{}
	token.AccessToken = accessToken
	token.RefreshToken = refreshToken

	return &token, nil
}

func (s *LoginService) TeacherLogin(ctx context.Context, req *user_service.LoginPasswors) (*user_service.Token, error) {
	s.log.Info("---LoginTeacher--->>>", logger.Any("req", req))

	resp, err := s.strg.Teacher().GetByLogin(ctx, req.Login)
	if err != nil {
		s.log.Error("---LoginTeacher--->>>", logger.Error(err))
		return &user_service.Token{}, err
	}

	if err = password.CompareHashAndPassword(resp.Password, req.Password); err != nil {
		s.log.Error("error while comparing password", logger.Error(err))
		return &user_service.Token{}, err
	}

	////
	m := make(map[interface{}]interface{})

	m["user_id"] = resp.Id
	m["user_role"] = "Teacher"

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		s.log.Error("error while generating tokens for User login", logger.Error(err))
		return &user_service.Token{}, err
	}

	token := user_service.Token{}
	token.AccessToken = accessToken
	token.RefreshToken = refreshToken

	return &token, nil
}

func (s *LoginService) SuperAdminLogin(ctx context.Context, req *user_service.LoginPasswors) (*user_service.Token, error) {
	s.log.Info("---LoginTeacher--->>>", logger.Any("req", req))

	if req.Login != "SuperAdmin" || req.Password != "SuperAdmin1!" {
		err := errors.New("incorrect login or password")
		s.log.Error("---LoginSuperAdmin--->>>", logger.Error(err))
		return &user_service.Token{}, err
	}

	////
	m := make(map[interface{}]interface{})
	Id := "e924cb31-e068-4062-a3b9-66790722e68a"
	m["user_id"] = Id
	m["user_role"] = "SuperAdmin"

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		s.log.Error("error while generating tokens for SuperAdmin login", logger.Error(err))
		return &user_service.Token{}, err
	}

	token := user_service.Token{}
	token.AccessToken = accessToken
	token.RefreshToken = refreshToken

	return &token, nil
}
