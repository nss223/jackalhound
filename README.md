# Jackalhound

一个 hyperledger/fabric 原型
* 网络一键启停
* 用户与账户管理
* 虚拟资产管理
* 积分系统
* 数据流转平台
* Java 客户端

## Getting Start

### Prerequisites
* server
	- docker
	- golang
	- hyperledger/fabric binaries
* client
	- jre
	- fabric-java-sdk

### Installing
* server: see [hyperledger/fabric prerequisites](http://hyperledger-fabric.readthedocs.io/en/release-1.1/prereqs.html)
* client: see [hyperledger/fabric-java-sdk](https://github.com/hyperledger/fabric-sdk-java)

## Running the tests
* server
	- 安装
	- 拉取仓库，并执行
```
cd jackalhound
cd basic-network

# edit `fabric.sh' as need
./restart
```
	- 等待网络启动
	- 调用链码：`./query`或`./invoke`，详见
	- 关闭网络：`./teardown.sh`
更多问题详见 [/basic-network/README.md](basic-network/README.md)

* client
	- 安装
	- 拉取仓库，以`/hx_jclient`为根导入 maven 项目
	- 执行`maven install`
	- 编译
	- 运行 test
更多问题详见 [/hx_jclient/README.md](hx_jclient/README.md)

## Authors
* **TeemoGuo** - *架构设计，编写设计文档*
* **ToricDong** - *架构设计，编写设计文档*
* **houndpan** - *链码设计，链码实现*
* **QwertyJack** - *设计实现，网络调试，java 客户端*

## Copyright
Copyleft (C) 2018, LAB2528
All rights reserved.