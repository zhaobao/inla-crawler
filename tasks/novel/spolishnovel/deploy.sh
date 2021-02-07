#!/usr/bin/env bash

# upload DB
aws s3 cp db.sqlite s3://inla-assets-center/novel/68a0f95f/db.sqlite \
  --profile zb

# upload RES
aws s3 cp output/ s3://inla-assets-center/novel/68a0f95f --recursive \
  --exclude "*" \
  --include "*.png" \
  --include "*.jpg" \
  --include "*.html" \
  --include "*.mp3" \
  --include "*.txt" \
  --profile zb
