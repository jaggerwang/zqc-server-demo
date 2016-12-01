package services

import (
	"math/rand"
	"time"

	"gopkg.in/mgo.v2/bson"

	"jaggerwang.net/zqcserverdemo/models"
	"jaggerwang.net/zqcserverdemo/utils"
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
	Id             bson.ObjectId `json:"id"`
	Username       string        `json:"username"`
	Nickname       string        `json:"nickname"`
	Gender         string        `json:"gender"`
	Mobile         string        `json:"mobile"`
	AvatarType     string        `json:"avatarType"`
	AvatarName     string        `json:"avatarName"`
	AvatarId       bson.ObjectId `json:"avatarId"`
	Email          string        `json:"email"`
	Intro          string        `json:"intro"`
	BackgroundType string        `json:"backgroundType"`
	BackgroundName string        `json:"backgroundName"`
	BackgroundId   bson.ObjectId `json:"backgroundId"`
	Location       *Location     `json:"location"`
	CreateTime     *time.Time    `json:"createTime"`
	UpdateTime     *time.Time    `json:"updateTime"`
	AvatarFile     *File         `json:"avatarFile"`
	BackgroundFile *File         `json:"backgroundFile"`
}

func NewUserFromModel(m *models.User) (user *User) {
	user = &User{
		Id:             m.Id,
		Username:       m.Username,
		Nickname:       m.Nickname,
		Gender:         m.Gender,
		Mobile:         m.Mobile,
		AvatarType:     m.AvatarType,
		AvatarName:     m.AvatarName,
		AvatarId:       m.AvatarId,
		Email:          m.Email,
		Intro:          m.Intro,
		BackgroundType: m.BackgroundType,
		BackgroundName: m.BackgroundName,
		BackgroundId:   m.BackgroundId,
		CreateTime:     m.CreateTime,
		UpdateTime:     m.UpdateTime,
	}

	if m.Location != nil {
		user.Location = &Location{m.Location.Coordinates[0], m.Location.Coordinates[1]}
	}

	if user.AvatarType == AvatarTypeCustom {
		if file, err := GetFile(user.AvatarId); err == nil {
			user.AvatarFile = file
		}
	}

	if user.BackgroundType == BackgroundTypeCustom {
		if file, err := GetFile(user.BackgroundId); err == nil {
			user.BackgroundFile = file
		}
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
		Id:             bson.NewObjectId(),
		Username:       "_" + utils.RandString(8, nil),
		Password:       password,
		Salt:           salt,
		Gender:         UserGenderMale,
		Mobile:         mobile,
		AvatarType:     AvatarTypeBuiltin,
		AvatarName:     builtinAvatars[rand.Intn(len(builtinAvatars))],
		BackgroundType: BackgroundTypeBuiltin,
		BackgroundName: builtinBackgrounds[rand.Intn(len(builtinBackgrounds))],
		CreateTime:     &t,
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
	if location, ok := update["location"]; ok {
		loc := location.(*Location)
		update["location"] = models.NewPoint(loc.Longitude, loc.Latitude)
	}

	update["updateTime"] = time.Now()
	err = c.UpdateId(id, bson.M{
		"$set": update,
	})
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
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

func GetUserByUsername(username string) (user *User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	var m models.User
	err = c.Find(bson.M{"username": username}).One(&m)
	if err != nil {
		return nil, NewServiceError(ErrCodeNotFound, err.Error())
	}

	return NewUserFromModel(&m), nil
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

func GetUserByEmail(email string) (user *User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	var m models.User
	err = c.Find(bson.M{"email": email}).One(&m)
	if err != nil {
		return nil, NewServiceError(ErrCodeNotFound, err.Error())
	}

	return NewUserFromModel(&m), nil
}

func VerifyUserPassword(username string, password string) (user *User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	var m models.User
	err = c.Find(bson.M{"username": username}).One(&m)
	if err != nil {
		return nil, NewServiceError(ErrCodeNotFound, err.Error())
	}

	if utils.Md5WithSalt(password, m.Salt) != m.Password {
		return nil, NewServiceError(ErrCodeWrongPassword, "")
	}

	return NewUserFromModel(&m), nil
}

func NearbyUsers(loc *Location, dist int, limit int) (users []*User, err error) {
	c, err := models.NewUserColl()
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	var ms []models.User
	err = c.Find(bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float32{loc.Longitude, loc.Latitude},
				},
				"$maxDistance": dist,
			},
		},
	}).Limit(limit).All(&ms)
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}

	users = make([]*User, 0, limit)
	for _, m := range ms {
		users = append(users, NewUserFromModel(&m))
	}

	return users, nil
}
