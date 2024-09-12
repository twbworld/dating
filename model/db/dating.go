package db

import "encoding/json"

type Dating struct {
	BaseField
	CreateUserId uint   `db:"create_user_id" json:"create_user_id" info:"会面创建者"`
	Result       string `db:"result" json:"result" info:"推荐结果"`
	Status       int8   `db:"status" json:"status" info:"会面状态; 0:结束;1:进行中;"`
	AddTime      int64  `db:"add_time" json:"-"`
}

func (Dating) TableName() string {
	return `dating`
}

// json转结构体
func (d *Dating) ResultUnmarshal() *DatingResult {
	result := &DatingResult{
		false,
		[]string{},
	}
	if d.Result == "" {
		return result
	}
	var res DatingResult
	if json.Unmarshal([]byte(d.Result), &res) != nil {
		return result
	}

	if res.Date == nil {
		//避免转json后,属性为"null"
		res.Date = make([]string, 0)
	}
	return &res
}

// Result的数据
type DatingResult struct {
	Res  bool     `json:"r" info:"匹配是否成功"`
	Date []string `json:"d"`
}

// 结构体转json
func (d *DatingResult) Marshal() string {
	if d.Date == nil {
		//避免转json后,属性为"null"
		d.Date = make([]string, 0)
	}

	result, err := json.Marshal(d)
	if err != nil {
		return ``
	}
	return string(result)
}
