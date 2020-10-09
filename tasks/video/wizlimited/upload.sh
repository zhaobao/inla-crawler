#!/usr/bin/env bash

# 上传视频和缩略图
aws s3 cp /Volumes/extend/crawler/video/wiz/6089c8e4 s3://inla-assets-center/video/f8663b49/6089c8e4 --recursive \
  --exclude "*" \
  --include "*.png" \
  --include "*.jpg" \
  --include "*.html" \
  --include "*.mp3" \
  --include "*.mp4" \
  --profile zb >> progress.txt
