# douyin-proj

框架设计  
|-- bin #项目需要引用的文件                   
|-- resources  
     |-- application.ini #配置文件  
|-- src
     |-- config          #全局参数配置  
     |-- controller      #控制层  
     |-- database        #数据库的管理模块  
     |-- global          #全局需要的组件  
     |-- repository      #数据模型与数据库操作  
     |-- server          #gin中间件与路由配置  
     |-- service         #负责数据层与控制层之间的逻辑  
     |-- types           #接口声明  
|-- upload #本地保存的视频与封面文件

video
* ffmpeg依赖: https://ffbinaries.com/downloads ,根据系统下载对应的ffmpeg二进制文件(只需要ffmpeg，不需要ffprobe,ffplay,ffserver)，并放入GoPath的`bin`目录下。
* IP地址: 校园网只能使用内网地址，外网地址无法访问，可在`servce/videoService.go`第121行手动指定IP。