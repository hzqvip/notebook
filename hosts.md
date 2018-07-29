# hosts文件简介

hosts文件对于普通互联网用户可能比较陌生，但对于开发者而言或多或少总能遇见，特别是IT&运维人员更是常年接触。那么它到底是什么呢？对上网有什么作用呢？

域名到内容的简单路由过程：
>  输入域名 -- DNS（Domain Name Server）服务器解析 -- 路由到对应 IP -- 访问服务器返回内容

在上述过程中还有一些步骤未体现出来且，其中一项就与系统hosts文件有关。

当我们输入域名到电脑去请求资源时，系统会优先去本地 hosts 文件里查询一次是否有与域名对应的 IP，以便快速访问或指定访问。如果存在则不经过 DNS 解析，直接访问对应 IP 这样就节省了解析时间（有时DNS服务解析可能经过几次跳转，多层代理），如果经过 DNS 解析中间可能就会出现其他情况。 例如

>  DNS 服务器本身不稳定，超时
>  DNS 劫持，污染
>  GFW 还可以利用了DNS做很多文章

### 小结
hosts文件对于域名解析有着重要重做，解析优先级 hosts > dns > apigateway 

## hosts文件“外貌”

hosts是一个纯文本文件，使用记事本即可编辑。我们可以先看一下windows下面的hosts文件初始内容是什么？ C:\windows\system32\drivers\etc（在系统盘下）

注意可以看到，前面几行文字就是微软介绍他的 hosts 文件是干什么用的。# 号是注释。

> 102.54.94.97  rhino.acme.com # 中间一般用 tab 隔开


## hosts文件运用场景

### 指定访问 IP

如果一个服务有多个服务器IP的话，那么可以把该网址指定到访问某台服务器上。

``` json
153.35.175.112   www.163.com  
127.0.0.1        www.ly.com
```

如果，我把 www.ly.com 的域名绑定到我本地了，所以访问的其实就是我本地，这样就永远访问不了此域名。 (这种方式就可以屏蔽恶意网站)

### 加速域名解析
通过绑定指定 IP ，则解析过程中就不会访问 dns 服务器，则节省多次跳转，访问到目标服务器上。

### 绕过DNS污染&劫持
绕过DNS污染，我们在国内可能访问一些特定网站会出现无法访问的情况 （eg Facebook YouTube），这里有一层原因就是利用的 dns 污染屏蔽掉了该网址后面真正提供服务的地址，在之前我们还可以通过找到 google 后台服务器地址，通过修改 hosts 直接访问，如今已经不行了 gfw 已经过程多种方式让国民健康上网。但有一些国外比较冷门的服务商可能还行。

还有一种就是劫持了，我们一般普通网民的DNS服务器ISP（*Internet Service Provider*）提供商都是电信，联通（网通），移动（铁通）这三家提供，那么有趣的一个现象就是，我们有时输入一个网址，得到的却是说到不到网站，然后跳出一个114网页，里面一大堆的垃圾广告，恶心的要死，关键是有些网站明明能够访问，却被半道劫持了，给我们返回一个我们所不期望的 IP，以前有段时间如果你输入 www.google.com ，神奇的是它会给你返回 www.baidu.com 的网址。所有hosts文件对于访问一些特定的网站非常有用。

### 恶意软件利用绑定导流

我们的hosts文件也会被恶意病毒所利用，把你的很多域名绑定到指定非所期望的 IP 地址上，例如我们想访问 www.baidu.com 但是被病毒劫持到另外一台服务器上了，可能就是一个带有搜索框的广告页面。所以hosts文件的读写权限一般也是恶意软件常驻之地。当然我们也可以利用一些杀毒软件隔离一些恶意网站和控制 hosts 文件的读写权限。

### 屏蔽特定网站

如果我不想孩子访问游戏网站或看直播等，则把 游戏网站的网址 绑定到本地 127.0.0.1 即可

### 方便局域网用户

有时局域网中因为机器比较少，没有资源搭建本地 dns 服务器，那么这个时候想到相互访问需要输入一大串 IP 地址就显得非常不方便，那么就可以通过维护hosts文件来维护你想要到的地方。既方便又快捷。但是当机器多了，维护成本也将变得较高。

## 如何合理的使用hosts文件

虽然hosts文件有这么多好处，但是如果当绑定的量较大的时候手动维护成本就会非常的高了，而且要随时注意你所绑定的IP是否能够随时访问，因为，一旦不能访问，你就需要手动去修改它，所以市面上有一些软件可以自动管理hosts文件的，有时我们还要观察一下我们的hosts文件是否正常。如果发现你的权限不能写入或者发现有大量你不清楚情况的网址出现在了上面，那么你应该考虑考虑杀杀毒了。

## Windows & Linux hosts文件位置&修改&生效

说了这么多，那么我们的hosts文件到底在什么地方呢？怎么修改呢？怎么生效呢？

### windows 

- 位置： C:\windows\system32\drivers\etc（在系统盘下） 
- 修改：通过记事本打开，填写好需要绑定的关系，保存。 保存的时候有可能提醒你 unable save，这是就要注意你的系统当前用户是否具有此文件的写入权限。 可以百度 如果切换系统管理员，还可以看是否杀毒软件隔离了权限，以及通过杀毒软件查询是否恶意软件修改了此文件权限。
- 生效：上述步骤操作成功后，为什么没有生效，因为window系统为了加快访问，对你经常访问的域名解析是做了缓存的，所以我们这里要清一下 DNS 服务本地缓存。步骤是先关掉浏览器（在打开浏览器那一刻，一般浏览器就读取了你当前本地缓存）。然后点开运行（win + r）输入  ipconfig /flushdns （注意有两词之间有一个空格）回车。一般就会生效，可以通过浏览器或 ping 命令验证
  
### Linux | unix | macos
- 位置：/etc/hosts
- 修改：vim /etc/hosts 
- 生效：一般修改后立即生效，如果没有生效，重启一下即ok。

## 小结

当网络访问未达到预期的值，可以先检查一下 hosts 文件。当想访问特定的服务，可以尝试先修改 hosts 达到目的。

有些网站被 dns 服务商劫持或者污染了，我们可以尝试修改我们的 dns 服务器地址，如下

``` sh
# 国际
8.8.8.8
8.8.4.4
1.1.1.1

# 国内阿里云
223.5.5.5
223.6.6.6
``` 
