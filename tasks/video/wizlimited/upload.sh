#!/usr/bin/env bash

# 上传视频和缩略图
aws s3 cp /Volumes/extend/crawler/video/wiz/eccf485d s3://inla-assets-center/video/f8663b49/eccf485d --recursive \
  --exclude "*" \
  --include "*.png" \
  --include "*.jpg" \
  --include "*.html" \
  --include "*.mp3" \
  --include "*.mp4" \
  --profile zb >>progress.txt
