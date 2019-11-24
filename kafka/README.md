
http://zhongmingmao.me/2018/10/08/kafka-install-cluster-docker/

## 1.修改host
```bash
127.0.0.1 zk1 zk2 zk3

127.0.0.1 kafka1 kafka2 kafka3
```


## 2.创建主题

```bash
# 进入kafka1
$ docker exec -it kafka1 bash

# 创建主题
root@kafka1:/# kafka-topics --zookeeper zk1:2181,zk2:2181,zk3:2181 --replication-factor 1 --partitions 1 --create --topic zhongmingmao
Created topic "test".
root@kafka1:/# kafka-topics --zookeeper zk1:2181,zk2:2181,zk3:2181 --describe --topic test

# 发送消息
root@kafka1:/# kafka-console-producer --broker-list kafka1:9092,kafka2:9092,kafka3:9092 --topic=test
>hello
```

## 3.读取消息

```bash
# 进入kafka1
# 进入kafka2
$ docker exec -it kafka2 bash

# 读取消息
root@kafka2:/# kafka-console-consumer --bootstrap-server kafka1:9092,kafka2:9092,kafka3:9092 --topic test --from-beginning
hello
zhongmingmao
```