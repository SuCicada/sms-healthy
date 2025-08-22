android 受到sms  ->  自动转发给（bark） ->  iphone

问题：
流程线路可能失败

测试：
手动发短信  -》 android -》 转发  -》 目标



1. send sms
2. check telegram, if not, send alert


我们默认认为，只要网络是连通的，那么bark就一定能发送成功。 
而只要telegram能收到消息，那么就能证明网络是连通的。 

这样就不会使用我自己的服务，避免了我自己的网络异常带来的影响

