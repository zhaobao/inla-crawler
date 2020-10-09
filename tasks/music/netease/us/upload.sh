#!/usr/bin/env bash

# 上传音频、歌词、封面
aws s3 cp output/ s3://inla-assets-center/music/fc154b81 --recursive \
  --exclude "*" \
  --include "*.png" \
  --include "*.jpg" \
  --include "*.html" \
  --include "*.txt" \
  --include "*.mp3" \
  --include "*.mp4" \
  --profile zb >>progress.txt
