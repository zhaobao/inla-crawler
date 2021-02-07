#!/usr/bin/env bash

convert "$1" -scale 1x1\! -format '%[pixel:u]' info:-
