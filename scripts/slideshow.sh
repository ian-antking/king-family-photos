#! /bin/bash

aws s3 sync s3://king-family-photos-dev-display /home/pi/Pictures --delete
feh \
  --recursive \
  --randomize \
  --fullscreen \
  --quiet \
  --hide-pointer \
  --slideshow-delay 60 \
  /home/pi/Pictures