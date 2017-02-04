package services

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"zqc/models"
	"zqc/utils"
)

const (
	UserGenderMale   = "m"
	UserGenderFemale = "f"
)

const (
	AvatarTypeBuiltin = "builtin"
	AvatarTypeCustom  = "custom"
)

var builtinAvatars = []string{
	"american-football-player-1",
	"american-football-player",
	"baseball-player",
	"basketball-player",
	"bodybuilder",
	"cricket-player",
	"cyclist-1",
	"cyclist",
	"fencer",
	"football-player",
	"formula-1",
	"golfer",
	"gymnast",
	"hockey-player",
	"horsewoman",
	"karate",
	"kickboxer",
	"kudo",
	"motorcyclist",
	"pilot",
	"rowing",
	"shooter",
	"skier-1",
	"skier",
	"sumotori",
	"swimmer",
	"taekwondo",
	"tennis-player",
	"volleyball-player",
	"weightlifter",
}

const (
	BackgroundTypeBuiltin = "builtin"
	BackgroundTypeCustom  = "custom"
)

var builtinBackgrounds = []string{
	"light-circle",
	"juhua",
	"pugongying",
}

type User struct {
	Id         bson.ObjectId `json:"id"`
	Mobile     string        `json:"mobile"`
	Nickname   string        `json:"nickname"`
	Gender     string        `json:"gender"`
	CreateTime *time.Time    `json:"createTime"`
	UpdateTime *time.Time    `json:"updateTime"`
}

func NewUserFromModel(m *models.User) (user *User) {
	user = &User{
		Id:         m.Id,
		Mobile:     m.Mobile,
		Nickname:   m.Nickname,
		Gender:     m.Gender,
		CreateTime: m.CreateTime,
		UpdateTime: m.UpdateTime,
	}

	return user
}

func CreateUser(mobile string, password string) (user *User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	salt := utils.RandString(16, nil)
	password = utils.Md5WithSalt(password, salt)
	t := time.Now()
	m := models.User{
		Id:         bson.NewObjectId(),
		Mobile:     mobile,
		Password:   password,
		Salt:       salt,
		CreateTime: &t,
	}
	err = c.Insert(m)
	if err != nil {
		return nil, NewServiceError(ErrCodeDuplicated, err.Error())
	}

	err = c.FindId(m.Id).One(&m)
	if err != nil {
		return nil, NewServiceError(ErrCodeNotFound, err.Error())
	}

	return NewUserFromModel(&m), nil
}

func UpdateUser(id bson.ObjectId, update bson.M) (user *User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	var m models.User
	err = c.FindId(id).One(&m)
	if err != nil {
		return nil, NewServiceError(ErrCodeNotFound, err.Error())
	}

	if password, ok := update["password"]; ok {
		update["password"] = utils.Md5WithSalt(password.(string), m.Salt)
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
		return nil, NewServiceError(code, err.Error())
	}

	err = c.FindId(id).One(&m)
	if err != nil {
		return nil, NewServiceError(ErrCodeNotFound, err.Error())
	}

	return NewUserFromModel(&m), nil
}

func GetUser(id bson.ObjectId) (user *User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	var m models.User
	err = c.FindId(id).One(&m)
	if err != nil {
		return nil, NewServiceError(ErrCodeNotFound, err.Error())
	}

	return NewUserFromModel(&m), nil
}

func GetUsers(ids []bson.ObjectId) (users []*User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	ms := make([]models.User, 0, len(ids))
	err = c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&ms)
	if err != nil {
		return nil, NewServiceError(ErrCodeNotFound, err.Error())
	}

	users = make([]*User, 0, len(ids))
	for _, m := range ms {
		users = append(users, NewUserFromModel(&m))
	}

	return users, nil
}

func GetUserByMobile(mobile string) (user *User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	var m models.User
	err = c.Find(bson.M{"mobile": mobile}).One(&m)
	if err != nil {
		return nil, NewServiceError(ErrCodeNotFound, err.Error())
	}

	return NewUserFromModel(&m), nil
}

func VerifyUserPassword(id bson.ObjectId, password string) (user *User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	var m models.User
	err = c.Find(bson.M{"_id": id}).One(&m)
	if err != nil {
		return nil, NewServiceError(ErrCodeNotFound, err.Error())
	}

	if utils.Md5WithSalt(password, m.Salt) != m.Password {
		return nil, NewServiceError(ErrCodeWrongPassword, "")
	}

	return NewUserFromModel(&m), nil
}
