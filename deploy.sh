#!/usr/bin/env bash

#aws s3 cp "$1" s3://inla-assets-center/comic/b24a0424/"$2" --profile zb
aws s3 cp "$1" s3://inla-assets-center/novel/b24a0424/"$2" --profile zb

#aws s3 cp /Volumes/extend/crawler/comic/imgs s3://inla-assets-center/comic/b24a0424 --recursive \
#  --exclude "*" \
#  --include "*.png" \
#  --include "*.svg" \
#  --include "*.jpg" \
#  --profile zb
