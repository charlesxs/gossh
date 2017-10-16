## gossh ##

一个简易的类似于Ansible ssh远程批量执行工具


**安装**

	1. yum -y install go && git clone https://github.com/charlesxs/gossh.git $GOPATH/src/gossh

	2. go get golang.org/x/crypto/ssh

	3. cd $GOPATH/src/gossh/ && go build


**使用**

	build 之后会生成一个 gossh 的二进制可执行文件, 然后执行此文件即可。

	使用方法:

	1. ./gossh command

	2. 在gossh所在目录编辑 hosts.conf 配置文件, 当此文件不存在时 执行 ./gossh command 会打印出配置文件模板


