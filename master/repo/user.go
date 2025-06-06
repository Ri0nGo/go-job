package repo

import (
	"go-job/internal/model"
	"go-job/internal/pkg/paginate"
	"gorm.io/gorm"
)

type IUserRepo interface {
	QueryById(id int) (model.User, error)
	QueryByUsername(username string) (model.User, error)
	Insert(*model.User) error
	Update(model.User) error
	Delete(id int) error
	QueryList(page model.Page) (model.Page, error)

	// oauth2
	QueryByAuth(authType model.AuthType, identity string) (model.User, error)
	CreateUserAndAuth(user *model.User, authModel *model.AuthIdentity) error
}

type UserRepo struct {
	mysqlDB *gorm.DB
}

func (j *UserRepo) QueryById(id int) (model.User, error) {
	var user model.User
	err := j.mysqlDB.First(&user, id).Error
	return user, err
}

func (j *UserRepo) Insert(user *model.User) error {
	return j.mysqlDB.Create(user).Error
}

func (j *UserRepo) Update(user model.User) error {
	if user.Id == 0 {
		return ErrorIDIsZero
	}
	return j.mysqlDB.Updates(&user).Error
}

func (j *UserRepo) Delete(id int) error {
	return j.mysqlDB.Where("id = ?", id).Delete(&model.User{}).Error
}

func (j *UserRepo) QueryList(page model.Page) (model.Page, error) {
	return paginate.PaginateList[model.DomainUser](j.mysqlDB, page)
}

func (j *UserRepo) QueryByUsername(username string) (model.User, error) {
	var user model.User
	err := j.mysqlDB.Where("username = ?", username).First(&user).Error
	return user, err
}

func (j *UserRepo) QueryByAuth(authType model.AuthType, identity string) (model.User, error) {
	// select t1.* from user t1 join auth_identity t2 on t1.id = t2.user_id and t2.identity = identity and t2.type = type
	var user model.User
	err := j.mysqlDB.Table("user AS t1").
		Select("t1.*").
		Joins("JOIN auth_identity AS t2 ON t1.id = t2.user_id").
		Where("t2.type = ? and t2.identity = ?", authType, identity).
		First(&user).Error
	return user, err
}

func (j *UserRepo) CreateUserAndAuth(user *model.User, authModel *model.AuthIdentity) error {
	return j.mysqlDB.Transaction(func(tx *gorm.DB) error {
		if err := j.Insert(user); err != nil {
			return err
		}
		authModel.UserID = user.Id
		if err := j.mysqlDB.Create(authModel).Error; err != nil {
			return err
		}
		return nil
	})
}

func NewUserRepo(mysqlDB *gorm.DB) IUserRepo {
	return &UserRepo{
		mysqlDB: mysqlDB,
	}
}
