## 启动步骤

1. 启动 etcd
2. 启动 zookeeper, kafka：
```shell
cd D:\Tools\kafka\kafka_2.12-3.5.1

# 启动 zookeeper
bin\windows\zookeeper-server-start.bat config\zookeeper.properties

# 管理员权限启动 kafka 
bin\windows\kafka-server-start.bat config\server.properties
```
3. 启动各项服务