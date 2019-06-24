#### 下载安装swagger-ui

```go
# https://hub.docker.com/r/swaggerapi/swagger-ui/
docker pull swaggerapi/swagger-ui:3.17.0
docker tag docker.io/swaggerapi/swagger-ui:3.17.0 harbor.m8.ccs/centos/swagger-ui:3.17.0
docker push harbor.m8.ccs/centos/swagger-ui:3.17.0

docker run -d \
	--name=swagger-ui-3.17.0 \
	-p 80:8080 \
	harbor.m8.ccs/centos/swagger-ui:3.17.0

浏览器直接访问：
http://ip:80
```

#### 下载安装swagger-editor

```go
# https://hub.docker.com/r/swaggerapi/swagger-editor/
docker pull swaggerapi/swagger-editor:3.5.7
docker tag docker.io/swaggerapi/swagger-editor:3.5.7 harbor.m8.ccs/centos/swagger-editor:3.5.7
docker push harbor.m8.ccs/centos/swagger-editor:3.5.7

docker run -d \
	--name=swagger-editor-3.5.7 \
	-p 80:8080 \
	harbor.m8.ccs/centos/swagger-editor:3.5.7

浏览器直接访问：
http://ip:80
```

