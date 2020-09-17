#!/usr/bin/env bash

input=$1
eval $(ffprobe -v error -of flat=s=_ -select_streams v:0 -show_entries stream=height,width ${input})
size=${streams_stream_0_width}x${streams_stream_0_height}
echo $size