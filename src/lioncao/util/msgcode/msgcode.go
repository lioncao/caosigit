package msgcode

const (
	SUCCESS = 0
	FAILED  = -1

	///////////////////////////////////////////////////////////////////////////
	// 服务器层面
	INVALID_USERID = 10000 // 无效的用户id

	///////////////////////////////////////////////////////////////////////////
	// 业务层面
	INVALID_CHAINID          = 20000 // 无效的底层链id
	INVALID_AREAID           = 20001 // 无效的领域id
	INVALID_TECHID           = 20002 // 无效的技术特性id
	INVALID_AREAID_OR_TECHID = 20003 // 无效的领域或技术特性
	INVALID_PROJECT_NAME     = 20004 // 无效的项目名字
	NOT_ENOUGH_GOLD          = 20005 // 金币不足
	NOT_ENOUGH_DIAMOND       = 20006 // 钻石不足
)
