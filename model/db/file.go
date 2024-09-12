package db

type File struct {
	BaseField        // (实质上当前表并没有update_time字段)
	Path      string `db:"path" json:"path" info:"文件路径"`
	Ext       string `db:"ext" json:"ext" info:"文件类型,如jpg/mp4等"`
	Type      int8   `db:"type" json:"type" info:"类型;0:本地;1:远程(如cdn)"`
}

func (File) TableName() string {
	return `file`
}
