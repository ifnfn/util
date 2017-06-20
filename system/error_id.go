package system

/***************************************************************
 * 通用错误码
 *
 * 只有通用的，常用的错误码放在这里，每个接口私有的错误码，放在接口内部
 ***************************************************************/

const (
	SqlErrorID    = 9
	TokenErrorID  = 101
	DelErrorID    = 7
	ModifyErrorID = 7
)

var (
	SuccessStatus          = &QError{Code: 0, Message: "成功"}
	AccessTokenErrorStatus = &QError{Code: 3, Message: "授权码已过期或者无效"}
	SqlErrorStatus         = &QError{Code: 1, Message: "SQL 语法错误,系统操作数据库错误"}
	ParameterErrorStatus   = &QError{Code: 2, Message: "参数解析错误"}

	ParameterKeyErrorStatus = &QError{Code: 4, Message: "参数关键字段错误或者不存在"}
	NoFoundStatus           = &QError{Code: 5, Message: "未找到匹配的数据"}
	PermissionErrorStatus   = &QError{Code: 6, Message: "没有该功能操作权限"}
	NoUserErrorStatus       = &QError{Code: 7, Message: "账户不存在或者密码错误"}
)
