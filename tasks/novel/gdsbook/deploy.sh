#!/usr/bin/env bash

# upload DB
aws s3 cp db.sqlite s3://inla-assets-center/novel/ba4491c2/db.min.sqlite

# upload RES
aws s3 cp output/data s3://inla-assets-center/novel/ba4491c2 --recursive \
     --exclude "*" \
     --include "*.epub" \
     --include "*.png" \
     --include "*.jpg" \
     --include "*.html" \
     --include "*.mp3" \
     --profile zb