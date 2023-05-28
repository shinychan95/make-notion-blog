#!/bin/bash

# 실행하려는 바이너리 파일 경로
BINARY_PATH="/Users/user/github/make-notion-blog/make-notion-blog"
CONFIG_PATH="/Users/user/github/make-notion-blog/config.json"
REPO_PATH="/Users/user/github/shinychan95.github.io"

# 바이너리 실행 전 모든 post 삭제 (이미지는 유지)
rm $REPO_PATH/_post/*

# 바이너리 실행 (config.json 파일 경로를 파라미터로 전달)
$BINARY_PATH -config $CONFIG_PATH

# git 저장소 경로로 이동
cd $REPO_PATH

# git add, commit, push
git add .
git commit -m "Update blog content"
git push
