#!/usr/bin/env bash

# usage sh thumb.sh [输入视频文件] [截取第几秒] [输出截图文件]
ffmpeg -i $1 -ss 00:00:$2 -vframes 1 $3