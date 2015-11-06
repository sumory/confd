### confd

confd是一个统一配置管理工具，idea参考自kelseyhightower的[confd](https://github.com/kelseyhightower/confd)。  
由于规划的功能会与[confd](https://github.com/kelseyhightower/confd)有很大不同，所以没有fork它。  
confd目前仍在开发中。

##### Features

- [x] 基础解析：配置文件使用toml格式，模板数据分离
- [x] 多种存储支持：local file、redis、zookeeper
- [x] 两种运行模式：Debug或Daemon
- [ ] cli工具：修改store，批量更新配置
- [ ] 更友好的交互方式，比如通过web界面
- [ ] 分离server和client，提供客户端cli或API供拉取指定配置
- [ ] 加密支持：store中存储的配置可加密，防止泄露

##### Usage

参考安装脚本[install.sh](./docs/install.sh)