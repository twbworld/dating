package db

type Feedback struct {
	BaseField
	Desc   string `db:"desc" json:"desc" info:"反馈详情"`
	UserId uint   `db:"user_id" json:"user_id" info:"反馈人"`
	FileId string `db:"file_id" json:"file_id" info:"file表id, 逗号相隔"`
}

func (Feedback) TableName() string {
	return `feedback`
}
