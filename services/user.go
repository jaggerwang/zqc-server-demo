package services

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"zqc/models"
	"zqc/util"
)

const (
	UserGenderMale   = "m"
	UserGenderFemale = "f"
)

type User struct {
	Id         bson.ObjectId `json:"id"`
	Mobile     string        `json:"mobile"`
	Password   string        `json:"-"`
	Nickname   string        `json:"nickname"`
	Gender     string        `json:"gender"`
	CreateTime *time.Time    `json:"createTime"`
	UpdateTime *time.Time    `json:"updateTime"`
}

func NewUserFromModel(m models.User) (user User) {
	user = User{
		Id:         m.Id,
		Mobile:     m.Mobile,
		Password:   m.Password,
		Nickname:   m.Nickname,
		Gender:     m.Gender,
		CreateTime: m.CreateTime,
		UpdateTime: m.UpdateTime,
	}

	return user
}

func CreateUser(mobile string, password string, nickname string, gender string) (user User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return user, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	salt := util.RandString(16, nil)
	password = util.Md5WithSalt(password, salt)
	t := time.Now()
	m := models.User{
		Id:         bson.NewObjectId(),
		Mobile:     mobile,
		Password:   password,
		Salt:       salt,
		Nickname:   nickname,
		Gender:     gender,
		CreateTime: &t,
	}
	err = c.Insert(m)
	if err != nil {
		return user, NewServiceError(ErrCodeDuplicated, err.Error())
	}

	err = c.FindId(m.Id).One(&m)
	if err != nil {
		return user, NewServiceError(ErrCodeNotFound, err.Error())
	}

	return NewUserFromModel(m), nil
}

func UpdateUser(id bson.ObjectId, update bson.M) (user User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return user, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	var m models.User
	err = c.FindId(id).One(&m)
	if err != nil {
		return user, NewServiceError(ErrCodeNotFound, err.Error())
	}

	if password, ok := update["password"]; ok {
		update["password"] = util.Md5WithSalt(password.(string), m.Salt)
	}

	update["updateTime"] = time.Now()
	err = c.UpdateId(id, bson.M{
		"$set": update,
	})
	if err != nil {
		code := ErrCodeSystem
		if strings.HasPrefix(err.Error(), "E11000 ") {
			code = ErrCodeDuplicated
		}
		return user, NewServiceError(code, err.Error())
	}

	err = c.FindId(id).One(&m)
	if err != nil {
		return user, NewServiceError(ErrCodeNotFound, err.Error())
	}

	return NewUserFromModel(m), nil
}

func GetUser(id bson.ObjectId) (user User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return user, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	var m models.User
	err = c.FindId(id).One(&m)
	if err != nil {
		return user, NewServiceError(ErrCodeNotFound, err.Error())
	}

	return NewUserFromModel(m), nil
}

func GetUsers(ids []bson.ObjectId) (users []User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return users, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	ms := make([]models.User, 0, len(ids))
	err = c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&ms)
	if err != nil {
		return users, NewServiceError(ErrCodeNotFound, err.Error())
	}

	users = make([]User, 0, len(ids))
	for _, m := range ms {
		users = append(users, NewUserFromModel(m))
	}

	return users, nil
}

func GetUserByMobile(mobile string) (user User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return user, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	var m models.User
	err = c.Find(bson.M{"mobile": mobile}).One(&m)
	if err != nil {
		return user, NewServiceError(ErrCodeNotFound, err.Error())
	}

	return NewUserFromModel(m), nil
}

func VerifyUserPassword(id bson.ObjectId, password string) (user User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return user, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	var m models.User
	err = c.Find(bson.M{"_id": id}).One(&m)
	if err != nil {
		return user, NewServiceError(ErrCodeNotFound, err.Error())
	}

	if util.Md5WithSalt(password, m.Salt) != m.Password {
		return user, NewServiceError(ErrCodeWrongPassword, "")
	}

	return NewUserFromModel(m), nil
}
