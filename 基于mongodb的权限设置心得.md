# 基于 mongodb 设计灵活后台管理权限
    mongodb 是一款基于 文档 结构的 nosql 数据库，目前社区比较火，在文档的存储灵活性有着天然的优势（同集合可以任何形式的行数据，当然我们不会这样去烂用这种灵活性 :) ），且有不俗的性能表现,以及副本集的高可用！
    此款权限粒度存在范围与个体的概念，可以理解为 功能性权限 和 数据归属性权限。数据权限是功能权限的再次细分，作为子集。下面，我会一点点剖析，下面的案例会给出所有表结构，和关键操作的 sql 语句（基于 golang 语言的 (mgo mongodb driver)）和设计图。当然这可能并不是最佳实践，只是我自己设计的一点心得，已便解决实际问题。
## 背景
  
之前在做一款早期 serverless 架构的雏形（现在只能算模版系统 :) ）后台管理系统的时候，公司的人员结构很稳定，又为了快速内部上线提供使用，我们想到一款最简单的基于部门实现泛权限的功能权限管理，意思就是，这个部门下面的组拥有这个组创建的所有查看权限，这个组拥有的功能权限由角色定义，如果为资深开发者角色，那么管理员会分配 模版编辑权限，模版清缓存/发布等一套可以从代码web ide 上编辑到外网访问发布的整个流程权限，初级可能只有编辑权限等。

但是我们很快发现问题，公司虽然有统一的权限系统，但是控制粒度太粗，所以开始我们是把人员的权限放入到自己的库里维护，时隔几月，随着公司的大举扩张，人员结构，部门等发生翻天覆地的变化，人事变更意味着之前基于部门设计的权限变成了脏数据，部门之间的数据没法自动实现访问，只能通过我们简单的数据权限操作机制，让某一人拥有多个部门的权限，这样处理因为不方便的原因可能存在越权处理，因为在 A 部门你拥有发布权限， 在 B 部门我只想让你有查询功能等，这样有一段时间我们就陷入了不断为人调整权限的售后工作！而且在跨部门处理问题都会面临不方便。

结构图

## 思考

当时我们面临一个元素种类比较多的后台系统，且权限组合种类比较多，那么怎样能够设计一款能够自由组合，能够自驱玩转权限，权限的下方不再由单人控制，以及在做工作交接的时候权限与资源能够一键交接。所以列了如下几点特性。

1. 权限控制下放，开发人员不再费力维护权限。
2. 权限有角色管理控制，有数据权限的细分概念。
3. 权限以及资源能够一键转让。
4. 权限基于人为个体，细分到，某个人下的某条数据。
5. 权限可以临时组成，不同开发人员，可以同时拥有一套或多套权限，方便跨部门协作。
6. 使用要简单，与用户基本信息隔离
7. 性能要出众。可以做成 web 服务 / sdk

## 实践

接下来就是进行表结构的梳理,所有结构都会给出 **最简洁关键字段及注解**

要想细分到数据权限，那么人/组必须要有一个元素与资源进行绑定，在此我们用 signkey 作为一种签名，某个人拥有这个签名就有对这个资源有一定的操作权限。一定要有中间签名的纽带来建立联系，不然数据量的存储，转让，变动都会出现问题。 如果说把资源属于哪些人来处理的话，可能一个字段就有很多工号来识别。

所以我们权限表里面的 user 信息表这样设计

``` go
type UserDoc struct {
	Name      string          // 姓名
	UserId    string          // 工号 唯一标识
	SignKey   map[string]string // 签名 key是签名 24字符id mongodb _id (可换成任意 唯一 字符串)，value是签名的描述。
}
```
可以看到一个人自己可以创建多个私有 签名
GroupName+UserId 作联合唯一索引

那么我们的角色权限是怎么设计的呢？一个角色拥有多种职能，我们把每一种职能看作一个接口实现了不同的方法，那么角色权限，就是由一个或多个接口实现的组合，接口实现又体现在 url + method 上面，然而问题是，有些接口涉及到操作数据，对数据有写操作，有的接口，只是单纯的某一次事件，不涉及到任何数据查询和变更。那么在这里就会划分此接口是否开启对数据做认证。通过不同的接口实现组合，这样就拥有了角色权限的维护。

``` go
type RoleDoc struct {
	RoleName        string          // 角色名称
	Desc            string          // 角色描述
	IsDefault       bool            // 是否为默认角色，就是用户登陆进来自带的角色权限
	UserIds         []string        // 角色下面有哪些用户，  这里是一个优化点，可以使用外键表关联。 视系统大小可以调整
	PathsDataVerify map[string]bool // 该接口是否开启数据验证权限，用于 sql 快速查询数据
	PathMethods     []PathMethod    // 该角色有哪些功能即 API 的组合， 详细数据，是一份冗余数据，方便配置时使用 通过方法把详细数据转化为快速查询数据
	Typ             int             // 角色类型，用于做特权处理， 例如 超级管理员角色，无需验证权限 则 type = 0
}

// restful API 风格设计 一个 path 有多个 method，
// method 使用 int 代替 GET=1 POST=2 PUT=4 DELETE=8 HEAD=16... 可以用到 mongodb 位运算
// 如果一个 path 的 get  post 方法都赋予了这个角色， 则 MethodInt=3
type PathMethod struct {
	Path      string // 请求路径
	MethodInt int // 方法的 int 和   9 代表 GET 和 DELETE 方法
}
```
上述表结构 PathsDataVerify key="1_/template" value=true 因为 1 表示 GET 所以这种表示描述就为 **模版的查询方法，需要开启数据认证权限**
roleName+Typ 作联合唯一索引

有了 Path 与角色的关联，就自然会有 Path 的详情，以方便管理，作为组成角色的权限的源，同时可以设置该接口是否开启数据验证权限。

``` go
type RouterDoc struct {
	Path      string                   // 请求路径  /template
	Desc      string                   // 路径实现功能的大类， 例如 模块的 CRUD
	MethodMap map[string]MethoedDetail // key 方法 具体表现  "1"  "2" ..
}

type MethoedDetail struct {
	DataVerify bool   // 该方法是否开启 数据验证
	Desc       string // 模版的删除 则 MethodMap key = 8
}
```

通过上面的简单表述，就可以建立 url + method 等于方法的实现，可以很方便的进行权限管理
Path 作唯一索引

那么我们的数据权限验证又是怎么实现的呢。它与角色权限又是什么样的关系。与人又怎么绑定的，注解如下

``` go
type SignDoc struct {
	SignKey      string         // 签名, 某人的私有签名 + 用户唯一标识  可以理解成 这个人拥有这个签名的哪些数据权限
	CreateUserId string         // 签名的创建者
	UserId       string         // 这个成员也拥有这个签名的一部分/全部权限， 相当于创建者把自己的签名共享出，用于多人使用
	RouterMap    map[string]int // key=Path value=Method   $bitsAllSet 计算  {"/template":14} 表述 这个人拥有此签名下面的模版数据 增/改/删的权限
}
```
SignKey+UserId 作联合唯一索引

这样创建者可分配的权限为自己的角色权限的子集，可以把一个私有key 暴露成公共的授权个他人，那么他人就对此签名下面所有的资源，都有相应你配置的权限，RouterMap 代表他拥有哪些权限，这样你就可以把自己的签名，分给不同的人，不同的人，拥有对此签名不同的权限。完全实现了，自己的数据自己做主管理。用于多种协调工作的场景。比如临时有一个修改活动逻辑的任务，有3个人做，那么他们组负责人可以创建一个公共的 key ，赋予这三个人，用于此次任务的所有权限操作。

### 所有的表结构都注解了，那么它们是怎么关联协作的，下面我们将从一个用户的角度使用权限系统里面所涉及的 sql 语句做简单解析，以便熟悉整个请求的流转。


#### 配置权限流程

初始化路由 Method 描述以及是否开启数据权限 ---> 创建角色 ---> 把才创建好的路由与创建的角色进行绑定 ---> 为角色添加用户 ---> 角色创建自己的签名 ---> 每一次请求携带上签名 ---> 创建资源的时候携带上此签名与资源绑定

查看个人自己拥有什么权限

``` go
    // userid 查询 RoleDoc 然后把所有返回的 PathMethods 做聚合
    // 通过 PathMethods 也可以结合 RouterDoc 查询具体的数据验证权限
```
查看别人授予的权限

``` go
    // userid 查询 SignDoc  然后对返回，做聚合即可
```

在权限的配置和管理上面，没有任何复杂的 sql ，当然在实际项目当中有很多的 sql 组合来满足不同的配置方式。总体下来是很方便的。

复杂分配举例，以对资源的 CRUD 为例
1. u1 u2 u3, u1 拥有 u2 签名 u2-s-1 下的 ceph资源读取权限， u2 拥有 u1-s-3 下的 template CRU 权限，u3-s-1 下的ceph CU 权限。 u3拥有 u1-s-3 所有数据的 R 权限

数据储存表现

``` json
// RoleDoc
[
    {"userId":"u1", "signKey":{"u1-s-3":"ceph kv manager"}},
    {"userId":"u2", "signKey":{"u2-s-1":"template manager"}},
    {"userId":"u3", "signKey":{"u3-s-1":"template manager"}}
]
```
``` json
// RouterDoc
[
    {"path":"/ceph","methodMap":{
        "1":{"dataVerify":true,"desc":"query ceph value"},
        "2":{"dataVerify":true,"desc":"create ceph value"},
        "4":{"dataVerify":true,"desc":"update ceph value"},
        "8":{"dataVerify":true,"desc":"delete ceph value"},
        }
    },
    {"path":"/template","methodMap":{
        "1":{"dataVerify":true,"desc":"query template value"},
        "2":{"dataVerify":true,"desc":"create template value"},
        "4":{"dataVerify":true,"desc":"update template value"},
        "8":{"dataVerify":true,"desc":"delete template value"},
        }
    }
]
```

``` json
    // RoleDoc 此角色没有删除权限
[
    {"reoleName":"ceph&template manager","userIds":["u1","u2","u3"],
    "pathMethods":[
        {"path":"/ceph","methodInt":7},
        {"path":"/template","methodInt":7}
    ]}
]

```
``` json
    // SignDoc 签名数据权限关联
[
    {"signKey":"u1-s-3","createUserId":"u1","userId":"u2","routerMap":{"/template":7}},
    {"signKey":"u3-s-1","createUserId":"u3","userId":"u2","routerMap":{"/ceph":6}},
    {"signKey":"u2-s-1","createUserId":"u2","userId":"u1","routerMap":{"/ceph":1}},
    {"signKey":"u1-s-3","createUserId":"u1","userId":"u3","routerMap":{"/template":1,"/ceph":1}}
]

```
当数据量特大的嵌套文档，都可以外键单独行存储。

#### 验证权限 

用户（统一单点）登陆后 ---> 判断是否为初始用户（默认角色） ---> 返回角色类型，拥有哪些权限（前端可以根据拥有的权限对 ui 进行渲染控制）

用户请求 携带 userid + signkey 请求获取 path + method （路由如果支持 动态路由， 用 httprouter 解析成后台配置的路由）
使用 userid + path + method 查询 roleDoc 表 验证是否存在角色权限（如果角色为超级管理员，则后面无需验证）
根据角色权限 PathsDataVerify 字段  value 值可以判断是否验证数据权限。
如果需要判断数据权限，则 携带 userid + signkey + path + method 查询 signDoc 表，是否存在 有权限则把相应的 权限信息写入 context 供后续方法使用

如果是查询接口 用户无需携带 signkey ， 需要查询的 signkey 由权限处理中间方法计算组成 signkey 数组 ， mongodb 使用 $in 查询数据
signkey 数组： 自己的私有 key ， 别人授权的数据权限 key 通过 userid + path + method 查询到所有授权的 signDoc 把里面的 signkey 放到数组里

## 总结

这一套权限验证是解决之前固有依赖部门设计的权限而提出的，此权限不存在与其他的依赖，只取决与签名与人与路由之间的绑定关系。可以任意组合！
数据权限是角色权限的子集，必须要满足角色权限。比如 u1 没有删除 /template 的角色权限， u2 授权 u1 拥有 u2-s-1 的 /template 的删除权限，也是不可以的。
通过这套权限把之前的数据进行了拆分，结构便于理解。当然此套表结构也是适用于其他数据库的。如果对性能有要求是可以用 redis 对权限做缓存的。当然目前我们基于200 + 用户， 200+ 路由配置，权限过滤和所有的管理结构都是在 20ms 下的。对于后台系统是可以接受的。

### 设计不好，仅供参考