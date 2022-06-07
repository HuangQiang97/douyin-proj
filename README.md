# douyin-proj

## 框架设计
```markdown
|__ bin #项目需要引用的文件                   
|__ resources  
|   |__ application.ini #配置文件  
|__ src  
|   |__ config          #全局参数配置  
|   |__ controller      #控制层  
|   |__ database        #数据库的管理模块  
|   |__ global          #全局需要的组件  
|   |__ repository      #数据模型与数据库操作  
|   |__ server          #gin中间件与路由配置  
|   |__ service         #负责数据层与控制层之间的逻辑  
|   |__ types           #接口声明  
|__ upload #本地保存的视频与封面文件
```

## 配置文件

- application.ini

```ini
[server]
HTTP_PORT = *******
HTTP_HOST = *******
MODE = debug

[mysql]
TYPE = mysql
USER = root
PASSWORD = ********
DB_HOST = ********
DB_PORT = 3306
DB_NAME = douyin
CHARSET = utf8mb4
ParseTime = true
MaxIdleConns = 20
MaxOpenConns = 100
Loc = Local

[jwt]
secretKey = *********

[crypto]
salt = *********
```

## video
* ffmpeg依赖: https://ffbinaries.com/downloads ,根据系统下载对应的ffmpeg二进制文件(只需要ffmpeg，不需要ffprobe,ffplay,ffserver)，并放入GoPath的`bin`目录下。
* IP地址: 校园网只能使用内网地址，外网地址无法访问，可手动指定IP。