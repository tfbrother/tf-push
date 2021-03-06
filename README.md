# tf-push
基于golang版本的websocket推送服务

# 存在的问题

* 连接管理器:
    * ~~假设有100万连接时，推送消息每次要遍历100万连接很费时，还有锁，还包含上下线。~~
    * 100万连接，只有一个goroutine处理推送消息，很慢。
    * ~~消息的发送都是同步的，比如批量推送的接口，要等全部推送完了才返回，要提高效率，需要改造成异步。~~初步解决
    * ~~缺少指标监控~~

# 解决办法
* 连接管理器:
    * 类似memcache，增加一层bucket，每个bucket来管理连接，连接管理器来管理bucket
    * 思路：在做海量服务架构设计的时候，一个很有用的思路就是：大拆小。
        > 既然100万连接放在一个集合里导致锁粒度太大，那么我们就可以把连接通过哈希的方式散列到多个集合中，每个集合有自己的锁。当我们推送的时候，可以通过多个线程，分别负责若干个集合的推送任务。因为推送的集合不同，所以线程之间没有锁竞争关系。而对于同一个集合并发推送多条不同的消息，我们可以把互斥锁换成读写锁，从而支持多线程并发遍历同一个集合发送不同的消息。
    * 消息推送改成异步，连接管理器增加接收和发送的队列，bucket也增加，连接管理器的发送和接收消息采用多个goroutine。
# 说明
* 查看监控指标
    * http://127.0.0.1:7777/stats