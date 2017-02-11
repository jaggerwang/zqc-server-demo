# Zaiqiuchang server demo

Zaiqiuchang is a mobile app developed using React Native, both iOS and Android are supported. This project is the lite version of the api service code, which developed in Go, and deploy using Docker. The lite version is only including user related api, but it should be a good start point for developing any api service. 

### Packages

|Package|Description|
|-------|-----------|
|[github.com/labstack/echo](https://echo.labstack.com/)|Http server and framework.|
|[github.com/spf13/cobra](https://github.com/spf13/cobra)|CLI app framework.|
|[github.com/spf13/viper](https://github.com/spf13/viper)|Manage app configuration.|
|[github.com/asaskevich/govalidator](https://github.com/asaskevich/govalidator)|Validate input.|
|[github.com/Sirupsen/logrus](https://github.com/Sirupsen/logrus)|Logging.|
|[gopkg.in/mgo.v2](https://labix.org/mgo)|Mongodb driver.|
|[github.com/garyburd/redigo/redis](https://github.com/garyburd/redigo/)|Redis driver.|
|[github.com/gorilla/sessions](https://github.com/gorilla/sessions)|Manage user session.|
|[github.com/smartystreets/goconvey/convey](https://github.com/smartystreets/goconvey)|Test framework.|

### How to deploy

You need install [docker engine](https://docs.docker.com/engine/installation/) first.

**run in dev mode with auto detecting code change**

```
> git clone git@github.com:jaggerwang/zqc-server-demo.git && cd zqc-server-demo
> mkdir -p ~/data/projects/zqc-server-demo # create directory for data volumes
> ./deploy.sh # pull images and run containers
> ./fswatch.sh # watching code change, fswatch needed
```

The data and log of server, mongodb and redis will be saved at host's path "~/data/projects/zqc-server-demo", which mounted at container's path "/data".

**run in prod mode**

```
> git clone git@github.com:jaggerwang/zqc-server-demo.git && cd zqc-server-demo
> mkdir -p /data/zqc-server-demo # create directory for data volumes
> ./deploy-prod.sh
```

The data and log of server, mongodb and redis will be saved at host's path "/data/zqc-server-demo", which mounted at container's path "/data".

**run test**

```
> cd zqc-server-demo
> ./test.sh
jagger-mbp 2017-02-11 16:00:44 INFO deploy test begin ...
Pulling mongodb (daocloud.io/jaggerwang/mongodb:latest)...
...
Pulling redis (daocloud.io/jaggerwang/redis:latest)...
...
Pulling server (daocloud.io/jaggerwang/zqc-server-demo:latest)...
...
zqcserverdemotest_mongodb_1 is up-to-date
zqcserverdemotest_redis_1 is up-to-date
zqcserverdemotest_server_1 is up-to-date
jagger-mbp 2017-02-11 16:00:50 INFO deploy test end
jagger-mbp 2017-02-11 16:00:50 INFO empty db begin ...
using config file ./config.json
empty db ok
jagger-mbp 2017-02-11 16:00:50 INFO empty db end
jagger-mbp 2017-02-11 16:00:50 INFO create db indexes begin ...
using config file ./config.json
create indexes ok
jagger-mbp 2017-02-11 16:00:50 INFO create db indexes end
jagger-mbp 2017-02-11 16:00:51 INFO run unittest begin ...
2017/02/11 16:00:51 goconvey.go:61: Initial configuration: [host: 0.0.0.0] [port: 8080] [poll: 250ms] [cover: true]
2017/02/11 16:00:51 tester.go:19: Now configured to test 10 packages concurrently.
2017/02/11 16:00:51 goconvey.go:192: Serving HTTP at: http://[::]:8080
2017/02/11 16:00:51 integration.go:122: File system state modified, publishing current folders... 0 40135146353
2017/02/11 16:00:51 goconvey.go:118: Received request from watcher to execute tests...
2017/02/11 16:00:51 executor.go:69: Executor status: 'executing'
...
2017/02/11 16:01:03 parser.go:24: [passed]: zqc/test
2017/02/11 16:01:03 executor.go:69: Executor status: 'idle'
```

Script `tesh.sh` will run a new group of docker containers to run unittests, including server, mongodb and redis. The new group has separated volumes and ports, to avoid empty the existing data. We use `goconvey` to write and run unittests. It supplied a web console to show the test result, and will run test automatically when detected code change.

Open url "http://localhost:10402/" to view test result.
<img src="https://zqc.cdn.zaiqiuchang.com/github/goconvey-result.png?x-oss-process=style/w-1080" />

**build image of your own**

```
> cd zqc-server-demo
> ./docker-build.sh
```

### Command

**help**

```
> cd zqc-server-demo
> docker-compose -p zqc-server-demo exec server zqc
Zai qiu chang app.

Usage:
  zqc [command]

Available Commands:
  db          Database admin
  server      Run server

Flags:
  -c, --config string              config file (default "./config.json")
      --dir.data string            directory for saving runtime data
      --env string                 deployment environment
      --log.level string           log filter level
      --mongodb.zqc.addrs string   address of zqc db

Use "zqc [command] --help" for more information about a command.
```

**create mongodb index**

```
> cd zqc-server-demo
> docker-compose -p zqc-server-demo exec server zqc db createIndexes
```
When deploy, it will auto run this command to create mongodb index. So normally you do not need to do this by your own.

### API

The server container exposed port 1323, and it mapped to port 10400 of the host. So you can use domain "http://localhost:10400" to access the following api.

Path|Method|Description
----|------|-----------
/register|POST|Register account.
/login|GET|Login.
/isLogined|GET|Check whether logined.
/logout|GET|Logout.
/account/edit|POST|Edit account profile.
/account/info|GET|Get current account info.
/user/info|GET|Get user info by id.
/user/infos|GET|Get user info by ids.

### FAQ

**How to change image repository?**

> Search and replace all "daocloud.io/jaggerwang/zqc-server-demo" to your own.

**How can I build the base images of this project, including go, mongodb and redis?**

> The dockerfiles of the base images can be found at "https://github.com/jaggerwang/jw-dockerfiles".

### Other resources

* [技术文章 - Go + Docker API服务开发和部署](https://jaggerwang.net/develop-and-deploy-api-service-with-go-and-docker-intro/)
* [在球场官网](https://www.zaiqiuchang.com)
