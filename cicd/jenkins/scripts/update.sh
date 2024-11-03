#!/bin/bash

# 確保腳本在錯誤時停止
set -e

# 檢查是否有未提交的更改
if [ -n "$(git status --porcelain)" ]; then
    echo "警告：有未提交的更改。請先提交或存儲這些更改。"
    exit 1
fi

# 獲取最新的 LTS 版本標籤
LATEST_LTS=$(curl -s https://api.github.com/repos/jenkinsci/jenkins/releases | grep '"tag_name":' | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | sort -V | tail -n 1)

# 更新 Dockerfile
sed -i "s/FROM jenkins\/jenkins:.*$/FROM jenkins\/jenkins:$LATEST_LTS-jdk17/" Dockerfile

# 構建新映像
docker-compose build

# 停止並刪除當前容器
docker-compose down

# 啟動新容器
docker-compose up -d

echo "Jenkins 已更新到版本 $LATEST_LTS 並重新啟動。"
echo "請檢查 http://localhost:8080 確保一切正常。"