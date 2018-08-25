package notebook

type RoleDoc struct {
	RoleName        string          // 角色名称
	Desc            string          // 角色描述
	GroupName       string          // 多项目用到的分组
	IsDefault       bool            // 是否为默认角色，就是用户登陆进来自带的角色权限
	UserIds         []string        // 角色下面有哪些用户，  这里是一个优化点，可以使用外键表关联。 视系统大小可以调整
	PathsDataVerify map[string]bool // 该接口是否开启数据验证权限，用于 sql 快速查询
	PathMethods     []PathMethod    // 该角色有哪些功能即 API 的组合
	Typ             int             // 角色类型，用于做特权处理， 例如 超级管理员角色，无需验证权限 则 type = 0
}

// restful API 风格 一个 path 有多个 method，
// method 使用 int 代替 GET=1 POST=2 PUT=4 DELETE=8 HEAD=16...
// 如果一个 path 的 get  post 方法都赋予了这个角色， 则 MethodInt=3
type PathMethod struct {
	Path      string
	MethodInt int
}

type RouterDoc struct {
	Path      string                   // 请求路径  /template
	Desc      string                   // 路径实现功能的大类， 例如 模块的 CRUD
	MethodMap map[string]MethoedDetail // key 方法 具体表现  "1"  "2" ..
}

type MethoedDetail struct {
	DataVerify bool   // 该方法是否开启 数据验证
	Desc       string // 模版的删除 则 MethodMap key = 8
}

type SignDoc struct {
	SignKey      string         // 签名, 某人的私有签名 + 工号  可以理解成 这个人拥有这个签名的哪些数据权限
	CreateUserId string         // 签名的创建者
	UserId       string         // 这个成员也拥有这个签名的一部分/全部权限， 相当于创建者把自己的签名共享出，用于多人使用
	RouterMap    map[string]int // key=Path value=Method   $bitsAllSet 计算  {"/template":14} 表述 这个人拥有此签名下面的模版数据 增/改/删的权限
}
