package service

import (
	"errors"
	"go-job/internal/model"
	"go-job/master/repo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUsernameDuplicate     = errors.New("duplicate username")
	ErrInvalidUserOrPassword = errors.New("用户名或密码不对")
)

type IUserService interface {
	GetUser(id int) (model.DomainUser, error)
	Login(username string, password string) (model.DomainUser, error)
	GetUserList(page model.Page) (model.Page, error)
	AddUser(job model.User) (model.User, error)
	DeleteUser(id int) error
	UpdateUser(job model.DomainUser) error
}

type UserService struct {
	UserRepo repo.IUserRepo
}

func (s *UserService) GetUser(id int) (model.DomainUser, error) {
	userDao, err := s.UserRepo.QueryById(id)
	if err != nil {
		return model.DomainUser{}, err
	}
	return s.userToDomainUser(userDao), nil
}

func (s *UserService) userToDomainUser(user model.User) model.DomainUser {
	return model.DomainUser{
		Id:          user.Id,
		Username:    user.Username,
		Nickname:    user.Nickname,
		CreatedTime: user.CreatedTime,
		About:       user.About,
	}
}

func (s *UserService) GetUserList(page model.Page) (model.Page, error) {
	return s.UserRepo.QueryList(page)
}

func (s *UserService) AddUser(user model.User) (model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	user.Password = string(hash)

	err = s.UserRepo.Insert(&user)
	switch {
	case err == nil:
		return user, nil
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return user, ErrUsernameDuplicate
	default:
		return user, err
	}
}

func (s *UserService) DeleteUser(id int) error {
	return s.UserRepo.Delete(id)
}

func (s *UserService) UpdateUser(user model.DomainUser) error {
	return s.UserRepo.Update(s.domainUserToUser(user))
}

func (s *UserService) domainUserToUser(user model.DomainUser) model.User {
	return model.User{
		Id:       user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		About:    user.About,
	}
}

func (s *UserService) Login(username, password string) (model.DomainUser, error) {
	user, err := s.UserRepo.QueryByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.userToDomainUser(user), ErrInvalidUserOrPassword
	}
	if err != nil {
		return s.userToDomainUser(user), err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return s.userToDomainUser(user), ErrInvalidUserOrPassword
	}
	return s.userToDomainUser(user), nil
}

func NewUserService(userRepo repo.IUserRepo) IUserService {
	return &UserService{
		UserRepo: userRepo,
	}
}
