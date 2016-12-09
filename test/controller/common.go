package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func deploy() (err error) {
	cmd := exec.Command("./deploy-prod.sh")
	cmd.Dir = filepath.Join(os.Getenv("GOPATH"), "src", "zqc")
	err = cmd.Run()
	return err
}

func destroy() (err error) {
	cmd := exec.Command("docker-compose", "-p", "zqctest", "down", "-v")
	cmd.Dir = filepath.Join(os.Getenv("GOPATH"), "src", "zqc")
	return cmd.Run()
}

func createDbIndexes() (err error) {
	cmd := exec.Command("docker-compose", "-p", "zqctest", "exec", "server", "zqc", "db", "createindexes")
	cmd.Dir = filepath.Join(os.Getenv("GOPATH"), "src", "zqc")
	return cmd.Run()
}

func emptyDb() (err error) {
	cmd := exec.Command("docker-compose", "-p", "zqctest", "exec", "server", "zqc", "db", "empty")
	cmd.Dir = filepath.Join(os.Getenv("GOPATH"), "src", "zqc")
	err = cmd.Run()
	return err
}

func postResp(path string, f url.Values, token string) (resp *http.Response) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("http://127.0.0.1:1324%s", path)
	body := strings.NewReader(f.Encode())
	req, err := http.NewRequest("POST", url, body)
	So(err, ShouldBeNil)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	So(err, ShouldBeNil)

	return resp
}

func post(path string, f url.Values, token string) (body []byte) {
	resp := postResp(path, f, token)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	So(err, ShouldBeNil)

	return body
}

func postResult(path string, f url.Values, token string) (result map[string]interface{}) {
	body := post(path, f, token)

	err := json.Unmarshal(body, &result)
	So(err, ShouldBeNil)

	return result
}
