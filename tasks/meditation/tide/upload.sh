#!/usr/bin/env bash

aws s3 cp output/ s3://inla-assets-center/meditation/dcd04d26 --recursive \
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
