#!/usr/bin/env bash

if [ -z "$1" ]; then
  echo "First parameter should be the new tag."
  exit 127
fi

FILES=templates/*.gitlab-ci.yml

for FILE in $FILES; do
  sed -i "s/^\(.*PIPES_VERSION: \)v.*/\1$1/g" $FILE
done
