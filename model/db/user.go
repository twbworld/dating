package db

type User struct {
	BaseField
	NickName    string `db:"nick_name" json:"nick_name" info:"昵称"`
	Avatar      uint   `db:"avatar" json:"-" info:"头像, 关联file表"`
	Gender      int8   `db:"gender" json:"gender" info:"性别;0:未知;1:男;2:女"`
	PhoneNumber string `db:"phone_number" json:"phoneNumber" info:"手机号"`
	OpenId      string `db:"openid" json:"-" info:"小程序平台的用户识别码"`
	UnionId     string `db:"unionid" json:"-" info:"微信用户识别码"`
	SessionKey  string `db:"session_key" json:"-" info:"微信Session_Key"`
}

func (User) TableName() string {
	return `user`
}
