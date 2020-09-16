#!/usr/bin/env bash

aws s3 cp input/ s3://inla-assets-center/comic/b24a0424 --recursive \
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