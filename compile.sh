#!/bin/bash

ARCH=$1


wails build -platform=darwin/${ARCH}
NAME="build/bin/Ztty-1.0.0-${ARCH}.dmg"

cd build/bin

rm -rf Ztty-1.0.0-${ARCH}.dmg

create-dmg \
  --volname "Ztty" \
  --window-size 600 400 \
  --icon-size 80 \
  --icon "Ztty.app" 100 120 \
  --app-drop-link 400 120 \
  "Ztty-1.0.0-${ARCH}.dmg" \
  "Ztty.app"
