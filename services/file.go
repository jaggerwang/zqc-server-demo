package services

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"

	"zqc/models"
)

type File struct {
	Id         bson.ObjectId `json:"id"`
	Url        string        `json:"url"`
	Size       int           `json:"size"`
	Filename   string        `json:"filename,omitempty"`
	UploaderId bson.ObjectId `json:"uploaderId"`
	CreateTime *time.Time    `json:"createTime"`
	UpdateTime *time.Time    `json:"updateTime,omitempty"`
}

func NewFileFromModel(m *models.File) (file *File) {
	urlBase := viper.GetString("storage.local.urlBase")
	if m.Bucket != "" {
		urlBase = viper.GetString(fmt.Sprintf("storage.buckets.%s.urlBase", m.Bucket))
	}

	file = &File{
		Id:         m.Id,
		Url:        urlBase + m.ObjectKey,
		Size:       m.Size,
		Filename:   m.Filename,
		UploaderId: m.UploaderId,
		CreateTime: m.CreateTime,
		UpdateTime: m.UpdateTime,
	}

	return file
}

func UploadFileToCloudStorage(bucket string, content []byte, objectKey string) (objectKeyReturn string, err error) {
	// TODO saved loally now, upload to cloud storage instead
	if objectKey == "" {
		ext := ""
		exts, _ := mime.ExtensionsByType(http.DetectContentType(content))
		if exts != nil {
			ext = exts[0]
		}
		objectKey = bson.NewObjectId().Hex() + ext
	}

	filename := path.Join(viper.GetString("storage.local.dir"), objectKey)
	err = ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		return "", err
	}

	return objectKey, nil
}

func CreateFile(uploaderId bson.ObjectId, bucket string, objectKey string,
	size int, filename string) (file *File, err error) {
	c, err := models.NewFileColl()
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}

	service := ""
	region := ""
	if bucket != "" {
		service = viper.GetString(fmt.Sprintf("storage.buckets.%s.%s", bucket, "service"))
		region = viper.GetString(fmt.Sprintf("storage.buckets.%s.%s", bucket, "region"))
	}
	t := time.Now()
	m := models.File{
		Id:         bson.NewObjectId(),
		Service:    service,
		Region:     region,
		Bucket:     bucket,
		ObjectKey:  objectKey,
		Size:       size,
		Filename:   filename,
		UploaderId: uploaderId,
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

	return NewFileFromModel(&m), nil
}

func GetFile(id bson.ObjectId) (file *File, err error) {
	c, err := models.NewFileColl()
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	var m models.File
	err = c.FindId(id).One(&m)
	if err != nil {
		return nil, NewServiceError(ErrCodeNotFound, err.Error())
	}

	return NewFileFromModel(&m), nil
}

func GetFiles(ids []bson.ObjectId) (files []*File, err error) {
	c, err := models.NewFileColl()
	if err != nil {
		return nil, NewServiceError(ErrCodeSystem, err.Error())
	}
	defer c.Close()

	ms := make([]models.File, 0, len(ids))
	err = c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&ms)
	if err != nil {
		return nil, NewServiceError(ErrCodeNotFound, err.Error())
	}

	files = make([]*File, 0, len(ids))
	for _, m := range ms {
		files = append(files, NewFileFromModel(&m))
	}

	return files, nil
}
