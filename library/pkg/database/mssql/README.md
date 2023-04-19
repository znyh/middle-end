#### database/mssql

##### 项目简介
sqlserver数据库驱动，进行封装加入了链路追踪和统计。

如果需要SQL级别的超时管理 可以在业务代码里面使用context.WithDeadline实现 推荐超时配置放到application.toml里面 方便热加载

##### 依赖包
1. [Go-MSSQL-Driver](https://github.com/denisenkom/go-mssqldb)
