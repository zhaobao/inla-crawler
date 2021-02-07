#!/usr/bin/env bash

# https://a.inlamob.com
aws s3 cp dist/ s3://inla-assets-center/wallpaper/vgals --recursive \
  --exclude "*" \
  --include "*.css" \
  --include "*.eot" \
  --include "*.eot?" \
  --include "*.woff" \
  --include "*.woff2" \
  --include "*.ttf" \
  --include "*.ico" \
  --include "*.png" \
  --include "*.svg" \
  --include "*.jpg" \
  --include "*.js" \
  --include "*.map" \
  --include "*.html" \
  --include "*.mp3" \
  --profile zb
