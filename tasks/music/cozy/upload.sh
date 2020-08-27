#!/usr/bin/env bash

aws s3 cp output/ s3://inla-assets-center/music/1e8c9cc8 --recursive \
     --exclude "*" \
     --include "*.png" \
     --include "*.mp3" \
     --profile zb