#!/usr/bin/env bash

echo "run style check for import pkg order"

for file in $(find . -name '*.go'); do goimports-reviser -project-name none-pjn $file; done

GIT_STATUS=`git status | grep "Changes not staged for commit"`;
if [ "$GIT_STATUS" = "" ]; then
  echo "code already go formatted";
else
  echo "style check failed, please format your code using goimports-reviser";
  echo "ref: github.com/incu6us/goimports-reviser";
  echo "git diff files:";
  git diff --stat | tee;
  echo "git diff details: ";
  git diff | tee;
  exit 1;
fi