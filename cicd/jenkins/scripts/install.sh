#!/bin/bash

# 確保腳本在錯誤時停止
set -e

# 檢查 Docker 和 Docker Compose 是否已安裝
if ! command -v docker &> /dev/null || ! command -v docker-compose &> /dev/null
then
    echo "Docker 或 Docker Compose 未安裝。請先安裝它們。"
    exit 1
fi

# 創建必要的目錄
mkdir -p config

# 如果 jenkins-casc.yaml 不存在，創建一個基本版本
if [ ! -f config/jenkins-casc.yaml ]; then
    cat > config/jenkins-casc.yaml <<EOL
jenkins:
  systemMessage: "Jenkins 已通過 Docker 配置"
  numExecutors: 2
  scmCheckoutRetryCount: 2
  mode: NORMAL
EOL
fi


# 構建並啟動 Jenkins
docker-compose up -d --build

echo "Jenkins 正在啟動。請稍候..."
echo "初始管理員密碼將在 Jenkins 完全啟動後顯示。"

# 等待 Jenkins 啟動
while ! docker-compose logs jenkins | grep -q "Jenkins is fully up and running"; do
    sleep 5
done

# 顯示初始管理員密碼
echo "Jenkins 已啟動。初始管理員密碼是："
docker-compose exec jenkins cat /var/jenkins_home/secrets/initialAdminPassword

echo "請訪問 http://localhost:8080 來完成 Jenkins 設置。"