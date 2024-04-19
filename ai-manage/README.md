# AI prompts manager

### 1. [下载go环境](https://go.dev/) 并安装IDE
### 2. 数据存储
    1. 安装mongo
    2. 创建 database `ai`
    3. 初始化数据
        1. 在 `ai` 数据库中分别创建 `prompt` `project` Collection
        2. 使用 `../docs/{project.json,prompt.json}`文件在上一步的Collection中初始化数据
            * 注：可以使用 `MongoDB Compass` IDE 的工具`Export Collection`来初始化数据
### 配置项目环境
    1. ../app.yaml 中填写 openai的key和mongo的地址
    2. 根据启动环境分别填写不同环境下的存储地址
        1. 系统环境：VERSE_ACTIVE=develop
        2. 文件环境：app.yaml-> app.active: "develop"
        3. **优先级**：系统环境大于文件环境(无需考虑命令行方式)
        4. 现在仅提供 `develop` 和 `us` 两个环境
    3. VPN 代理在国内(开发环境)需要自行使用代理，否则无法访问海外环境
        1. 配置地址：app-develop.yaml -> proxy:{  http:,  socket:}

### 启动方式
#### 本地启动
```shell
go run main.go
# debug使用IDE调试
```
#### docker

```shell
docker build -t $IMAGE_NAME .

ENV=develop && \
docker run -d -v /usr/share/zoneinfo/Asia/Shanghai:/etc/localtime:ro \
	--net=host --restart=always \
	-e VERSE_ACTIVE=$ENV --name $PROJECT_NAME \
	 $IMAGE_NAME

# $IMAGE_NAME和$PROJECT_NAME都可以按照自己项目管理取名
# 如果觉得麻烦可以通过上面的shell写一个docker-compose来管理,可以更好的管理控制 network，port，volume
```