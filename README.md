### confd

confd是一个统一配置管理工具，目前仍在开发中，基于本地文件存储的部署方式已经可以用于生产环境中。

##### Features

- [x] 基础解析：配置文件使用toml格式，模板数据分离
- [x] 多种存储支持：local file、redis、zookeeper
- [x] 两种运行模式：Debug或Daemon
- [x] cli工具：修改store，批量更新配置
- [ ] 更友好的交互方式，比如通过web界面
- [ ] 分离server和client，提供客户端cli或API供拉取指定配置
- [ ] 加密支持：store中存储的配置可加密，防止泄露

##### Usage

安装：

```
#假设安装路径为/data/server/confd
#go get获取依赖的第三方库 

sh build.sh linux /data/server/confd
#执行以上构建脚本后，在docs目录下生成了安装所需的文件
#若build.sh增加了参数$2，则默认加载$2/data/config.toml作为confd运行所需的配置文件

cd docs
sh install.sh /data/server/confd
cd /data/server/confd 
#注意修改${path}/data/config.toml里的ConfDir和ConnectAddr

#然后使用confd、confd-cli即可
```