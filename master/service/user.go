package service

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"go-job/internal/model"
	"go-job/master/repo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUsernameDuplicate     = errors.New("duplicate username")
	ErrInvalidUserOrPassword = errors.New("用户名或密码不对")
	ErrDuplicateEmail        = errors.New("邮箱已被占用")
)

type IUserService interface {
	GetUser(id int) (model.DomainUser, error)
	Login(username string, password string) (model.DomainUser, error)
	GetUserList(page model.Page) (model.Page, error)
	AddUser(user model.User) (model.User, error)
	DeleteUser(id int) error
	UpdateUser(user model.DomainUser) error

	UserBind(id int, email string) error
}

type UserService struct {
	UserRepo repo.IUserRepo
}

func (s *UserService) GetUser(id int) (model.DomainUser, error) {
	if id <= 0 {
		return model.DomainUser{}, errors.New("user id is zero")
	}
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
		Email:       user.Email,
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
	var me *mysql.MySQLError
	switch {
	case err == nil:
		return user, nil
	case errors.As(err, &me):
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return user, ErrUsernameDuplicate
		}
		return user, err
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

func (s *UserService) UserBind(id int, email string) error {
	err := s.UserRepo.Update(model.User{Id: id, Email: email})
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return ErrDuplicateEmail
		}
	}
	return err
}

func NewUserService(userRepo repo.IUserRepo) IUserService {
	return &UserService{
		UserRepo: userRepo,
	}
}
