#!/usr/bin/env bash

ffmpeg -i "$1" -vf scale="$2" "$3"