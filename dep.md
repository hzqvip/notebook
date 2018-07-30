## dep golang 包管理使用记录

dep 是 golang 项目依赖管理之一，是官方的实验项目，目前更新很频繁处于高速发展期，所以选 dep 作为 golang 的依赖管理器是比较靠谱的。（已知 glide 仅支持不再开发新功能）

目前 dep v0.5.0 release 已经发布，最新的 changelog 显示只支持 golang 1.9+ 以上的版本 

golang 最原始的依赖管理是 go get ，执行命令后会拉取代码放入 src 下面，但是它是作为 GOPATH 下全局的依赖，并且 go get 还不能版本控制，以及隔离项目的包依赖在没有依赖管理工具的时候，golang 项目有一种目录结构比较流行如下：

``` golang
.
└── src
    ├── demo
    │   └── main.go
    ├── github.com
    ├── golang.org
    └── gopkg.in

```

这样做的话就是每一个项目一个 GOPATH 则上面的 GOPATH=/xx/xx/src 这样设置后也是可以编译的，且项目依赖都是 src 下的包，与全局的无关联， go get
获取依赖，也必须在项目 GOPATH 下。可以看到这样的目录结构还是有很多缺陷的。特别是 import 包的包目录有时会很奇怪，没有统一的风格，且 ide 或 编辑器支持也不够理想。

所以目前 golang 引入了 vendor 目录作为依赖管理目录，且 ide 或 golang 编辑插件目前都能很好的支持 例如 gogland 索引依赖包时会优先查找项目根目录下的 vendor 目录， vscode 的 go 插件也是。那么目前比较流行的目录结构如下：

```golang
.
├── Gopkg.lock
├── Gopkg.toml
├── main.go
└── vendor
    ├── github.com
    │   ├── gin-contrib
    │   ├── gin-gonic
    │   ├── golang
    │   ├── mattn
    │   └── ugorji
    ├── golang.org
    │   └── x
    └── gopkg.in
        ├── go-playground
        └── yaml.v2

```
项目目录 $GOPATH/src/projectname/.

### Getting Started

#### install / uninstall

##### macOS

``` sh
brew install dep
brew uninstall dep
# /usr/local/bin/dep
```

##### windows

``` sh
# 推荐 go get 安装
go get -u github.com/golang/dep/cmd/dep
# $GOPATH/bin/dep.exe
```

##### Arch linux

```sh
pacman -S dep
# 删除 pacman -R dep
```

##### 二进制安装

``` sh
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
# $GOPATH/bin/dep
# 删除 > rm $GOPATH/bin/dep
```
##### 源码安装

``` sh
go get -d -u github.com/golang/dep
cd $(go env GOPATH)/src/github.com/golang/dep
DEP_LATEST=$(git describe --abbrev=0 --tags)
git checkout $DEP_LATEST
go install -ldflags="-X main.version=$DEP_LATEST" ./cmd/dep
git checkout master
```

推荐安装方式为 各系统的快捷安装，最方便且不容易出错， 如果想使用最新/指定版本，推荐源码安装

#### 初始化项目

推荐使用上述第二种目录结构 

``` sh
mkdir $GOPATH/src/example
cd $GOPATH/src/example
dep init

# 生成 vendor/ 目录下 Gopkg.toml Gopkg.lock
```


