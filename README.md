# douyin-proj

video
* ffmpeg依赖: https://ffbinaries.com/downloads ,根据系统下载对应的ffmpeg二进制文件(只需要ffmpeg，不需要ffprobe,ffplay,ffserver)，并放入GoPath的`bin`目录下。
* IP地址: 校园网只能使用内网地址，外网地址无法访问，可在`servce/videoService.go`第121行手动指定IP。