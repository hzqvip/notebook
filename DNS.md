## CentOS 修改 DNS 服务配置
ping www.baidu.com 报错 name or service not known
多数情况为新安装的系统未添加 DNS 服务器

## 添加方式

1. vi /etc/resolv.conf
2. 在文件中添加 dns 服务器地址
``` sh
    nameserver 8.8.8.8
    nameserver 8.8.4.4
```
3. 保存退出，重启服务器

## 小结
如果没有生效
1. vi /etc/sysconfig/network-scprits/ifcfg-ens33 (ens33 可能数字不同 ifconfig 查看)
2. ONBOOT=no 改成 ONBOOT=yes
3. service network restart