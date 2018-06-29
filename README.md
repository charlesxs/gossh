## gossh ##

一个简易的类似于Ansible ssh远程批量执行工具


**安装**

	1. 安装 go 和 gossh

	    yum -y install golang && git clone https://github.com/charlesxs/gossh.git $GOPATH/src/gossh

	2. 安装 glide, glide 默认安装在 $GOPATH/bin 目录下

        curl https://glide.sh/get | sh

        export PATH=$PATH:$GOPATH/bin

	3. 编译 (若要安装执行 make install, 则默认安装到 $GOPATH/bin 目录下)

	    cd $GOPATH/src/gossh/ && make


**使用**

	build 之后会生成一个 gossh 的二进制可执行文件, 然后执行此文件即可。

	使用方法:

	1. ./gossh command

	2. 在gossh所在目录编辑 hosts.conf 配置文件, 当此文件不存在时 执行 ./gossh command 会打印出配置文件模板


