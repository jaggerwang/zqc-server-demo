package test

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

func EmptyDb() (err error) {
	cmd := exec.Command("docker-compose", "-p", "zqc", "exec", "server", "zqc", "db", "empty")
	cmd.Dir = filepath.Join(os.Getenv("GOPATH"), "src", "zqc")
	err = cmd.Run()
	return err
}

func PostResp(path string, f url.Values, token string) (resp *http.Response) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("http://server:1323%s", path)
	body := strings.NewReader(f.Encode())
	req, err := http.NewRequest("POST", url, body)
	So(err, ShouldBeNil)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	So(err, ShouldBeNil)

	return resp
}

func Post(path string, f url.Values, token string) (body []byte) {
	resp := PostResp(path, f, token)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	So(err, ShouldBeNil)

	return body
}

func PostResult(path string, f url.Values, token string) (result map[string]interface{}) {
	body := Post(path, f, token)

	err := json.Unmarshal(body, &result)
	So(err, ShouldBeNil)

	return result
}
