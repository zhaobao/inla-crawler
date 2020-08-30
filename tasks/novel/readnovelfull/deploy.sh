#!/usr/bin/env bash

aws s3 cp output/ s3://inla-assets-center/novel/330706f3 --recursive \
     --exclude "*" \
     --include "*.png" \
     --include "*.jpg" \
     --include "*.html" \
     --include "*.mp3" \
     --profile zb