#!/bin/bash

# 設置備份目錄
BACKUP_DIR="/path/to/jenkins/backups"
BACKUP_NAME="jenkins_backup_$(date +%Y%m%d_%H%M%S).tar.gz"

# 創建備份目錄（如果不存在）
mkdir -p $BACKUP_DIR

# 停止 Jenkins 容器
docker-compose stop jenkins

# 創建備份
docker run --rm --volumes-from $(docker-compose ps -q jenkins) \
    -v $BACKUP_DIR:/backup ubuntu tar czf /backup/$BACKUP_NAME -C /var/jenkins_home .

# 重新啟動 Jenkins 容器
docker-compose start jenkins

echo "Jenkins 備份已創建：$BACKUP_DIR/$BACKUP_NAME"