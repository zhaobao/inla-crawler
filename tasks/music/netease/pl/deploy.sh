#!/usr/bin/env bash

# upload DB
aws s3 cp db.sqlite s3://inla-assets-center/music/3383308c/db.sqlite \
  --profile zb

# upload RES
aws s3 cp output/ s3://inla-assets-center/music/3383308c --recursive \
  --exclude "*" \
  --include "*.png" \
  --include "*.jpg" \
  --include "*.html" \
  --include "*.mp3" \
  --include "*.txt" \
  --profile zb

#aws s3 cp output/cd187ea3/m.mp3 s3://inla-assets-center/music/3383308c/cd187ea3/m.mp3 --profile zb
