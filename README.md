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

**安装**

```
#假设安装路径为/data/server/confd
#go get获取依赖的第三方库 

sh build.sh linux /data/server/confd
#执行以上构建脚本后，在docs目录下生成了安装所需的文件
#若build.sh增加了参数$2，则默认加载$2/data/config.toml作为confd运行所需的配置文件

cd docs
sh install.sh /data/server/confd
cd /data/server/confd 
#注意修改${path}/data/config.toml里的ConfDir和ConnectAddr，前缀为/data/server/confd

#然后使用confd、confd-cli即可
```

**使用**

```
.
├── confd
├── confd-cli
├── data
│   ├── config.toml
│   └── filestore.toml
├── meta
│   ├── example_1.toml
│   └── example_2.toml
└── templates
    └── example.tmpl
```


**confd**的使用：

- 首先了解：
    - confd的配置可参看[config.toml](./docs/files/config.toml), config.toml里的配置均可通过运行时指定flag来覆盖默认配置。
    - confid.toml中指定的配置数据存储为“file”形式，即数据存储在同目录下filestore.toml文件中。
    - 根据需要编辑模板，参考[example.tmpl](docs/files/example.tmpl)，模板中使用的变量目前只支持k/v形式。
    - 根据需要编辑meta文件，参考[example_1.toml](docs/files/example_1.toml)，meta文件制定了生成最终配置文件时需要的模板文件、数据、最终文件地址等。
- 命令
    - `./confd`, 在各个meta文件指定的目的地址生成了需要的配置文件
    - `./confd --debug=false`, 默认confd在后台运行，每10分钟重新生成一次全部的配置文件

**confd-cli**的使用：

- `confd-cli`是操作confd的命令行程序
- 目前支持的子命令

	<table>
    	<tr>
        	<td width="50%">./confd-cli getall</td>
        	<td width="50%">获取当前所有配置需要的数据</td>
    	</tr>
    	<tr>
        	<td width="50%">./confd-cli get key1</td>
        	<td width="50%">获取key1现在的值</td>
    	</tr>
    	<tr>
        	<td width="50%">./confd-cli set key1 value1</td>
        	<td width="50%">设置key1值为value1</td>
    	</tr>
    	<tr>
        	<td width="50%">./confd-cli delete key1</td>
        	<td width="50%">删除key1</td>
    	</tr>
	</table>

