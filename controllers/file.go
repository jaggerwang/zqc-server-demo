package controllers

import (
	"io/ioutil"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo"

	"zqc/middlewares"
	"zqc/services"
)

type UploadFileParams struct {
	Bucket string `valid:"stringlength(1|20),optional"`
}

func UploadFile(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	params := UploadFileParams{
		Bucket: cc.FormValue("bucket"),
	}
	if ok, err := valid.ValidateStruct(params); !ok {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}

	file, err := cc.FormFile("file")
	if err != nil {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}
	r, err := file.Open()
	if err != nil {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}
	defer r.Close()
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}

	objectKey, err := services.UploadFileToCloudStorage(params.Bucket, content, "")
	if err != nil {
		return services.NewServiceError(services.ErrCodeUploadFileToCloudStorage, err.Error())
	}

	userId := cc.SessionUserId()
	f, err := services.CreateFile(userId, params.Bucket, objectKey,
		len(content), file.Filename)
	if err != nil {
		return err
	}

	return ResponseJSON(http.StatusOK, Response{
		Data: map[string]interface{}{
			"file": f,
		},
	}, cc)
}

type FileInfoParams struct {
	Id string `valid:"objectidhex"`
}

func FileInfo(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	params := FileInfoParams{
		Id: cc.FormValue("id"),
	}
	if ok, err := valid.ValidateStruct(params); !ok {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}
	id, err := ParseObjectId(params.Id)
	if err != nil {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}

	file, err := services.GetFile(id)
	if err != nil {
		return err
	}

	return ResponseJSON(http.StatusOK, Response{
		Data: map[string]interface{}{
			"file": file,
		},
	}, cc)
}
