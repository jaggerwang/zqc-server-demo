# Zaiqiuchang server demo

Zaiqiuchang is a mobile app developed using React Native, both iOS and Android are supported. This project is the lite version of the server code, which written in Go, and deployed by Docker. Full feature app can be downloaded at [在球场](https://www.zaiqiuchang.com) . The lite version is only including user and file related api, but it should be a good start point for developing api service. 

### Screenshot

<img alt="Nearby screenshot" src="https://zqc.oss-cn-shanghai.aliyuncs.com/screenshot/ios/screenshot-nearby-720.jpg" width="360" height="640" />

### How to deploy

**run in prod mode**

Of course you should install docker engine first.
```
> git clone git@github.com:jaggerwang/zqc-server-demo.git && cd zqc-server-demo
> ./deploy-prod.sh
```
The deploy script use docker-compose to run all needed containers, including server, mongodb and redis.  
The data and log of server, mongodb and redis will be saved at host's directory "/data/zqc-server-demo", which mounted at container's path "/data". You can change the data dir to your own, but the same change should be made to docker-compose file.

**run in dev mode with auto detecting code changing**

```
> git clone git@github.com:jaggerwang/zqc-server-demo.git && cd zqc-server-demo
> ./deploy.sh
> ./fswatch.sh # fswatch command needed, you can use brew to install it on macOS
```
The data and log of server, mongodb and redis will be saved at host's directory "~/data/projects/zqc-server-demo/server", which mounted at container's path "/data". You can change the data dir to your own, but the same change should be made to docker-compose file.

**build image of your own**

```
> git clone git@github.com:jaggerwang/zqc-server-demo.git && cd zqc-server-demo
> ./docker-build.sh
```

### API

The server container exposed port 1323, and it mapped to port 10400 of the host. So you can use domain "http://localhost:10040" to access the following api.

Path|Method|Description
----|------|-----------
/register|POST|Register account.
/login|GET|Login.
/resetPassword|GET|Reset password by mobile verify code.
/isLogined|GET|Check whether logined.
/logout|GET|Logout.
/security/sendVerifyCode|GET|Send mobile verify code. It's a faked api, you should implement it by yourself.
/account/edit|POST|Edit account profile.
/account/info|GET|Get current account info.
/user/info|GET|Get user info by id.
/user/nearby|GET|Get nearby users.
/file/upload|POST|Upload file, saved on local disk.
/file/info|GET|Get uploaded file info by id.
/upload/a.jpg|GET|Access uploaded file. The path can be found in the response of file upload api.

### FAQ

**How to change image repository?**

> Search and replace all "daocloud.io/jaggerwang/zqc-server-demo" to your own.

**How can I build the base images of this project, including go, mongodb and redis?**

> The dockerfiles of the base images can be found at "https://github.com/jaggerwang/dockerfiles".

### Other resources

* [技术文章 - Go + Docker API服务开发和部署](https://jaggerwang.net/develop-and-deploy-api-service-with-go-and-docker-intro/)
* [在球场官网](https://www.zaiqiuchang.com)
