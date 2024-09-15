package user

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"
	"unicode/utf8"

	"mime/multipart"

	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/model/common"
)

type fc struct {
	ext  string
	size int64
}

var fileExt []fc = []fc{
	{"jpg", 2},
	{"jpeg", 2},
	{"png", 2},
	{"bmp", 2},
	{"gif", 2},
	{"mp4", 32},
	{"mov", 32},
	{"mkv", 32},
	{"avi", 32},
}

type Validator struct{}

// 检验LoginPost参数
func (v *Validator) ValidatorLoginPost(data *common.LoginPost) error {
	if data.Code == "" || len(data.Code) < 5 {
		return errors.New("参数错误[dotsd]")
	}
	return nil
}

// 检验UserAddPost参数
func (v *Validator) ValidatorUserAddPost(data *common.UserInfoPost) error {
	//nick_name可能为空, 不做判断
	if data.Code == "" || len(data.Code) < 5 {
		return errors.New("参数错误[dosfs0d]")
	}
	//判断文件是否存在
	return nil
}

// 检验FeedbackPost参数
func (v *Validator) ValidatorFeedbackPost(data *common.FeedbackPost) error {
	if len(data.Desc) < 1 || utf8.RuneCountInString(data.Desc) > 100 {
		return errors.New("请详细描述[dofsd]")
	}
	return nil
}

// 检验上传的文件
func (v *Validator) ValidatorUpload(file *multipart.FileHeader) error {
	ext := strings.TrimPrefix(path.Ext(file.Filename), ".")
	if ext == "" {
		return errors.New("文件异常[rkapd]")
	}
	for _, v := range fileExt {
		if v.ext == ext {
			if file.Size > v.size<<20 {
				return fmt.Errorf("文件限制%dMB", v.size)
			}
			return nil
		}
	}

	return errors.New("文件类型不支持[rkaopd]")
}

// 检验InfoPost参数
func (v *Validator) ValidatorInfo(data *common.InfoPost) error {
	if len(data.Time) < 1 {
		return errors.New("参数错误[doifj]")
	}

	allowedStartDate, allowedEndDate, dates := time.Now().AddDate(-1, 0, 0), time.Now().AddDate(1, 0, 0), make(map[string]bool) //允许两年时间跨度
	for _, value := range data.Time {
		te, err := time.ParseInLocation(time.DateTime, value[1], global.Tz)
		if err != nil {
			return errors.New("时间选择错误[oilsng]")
		}
		ts, err := time.ParseInLocation(time.DateTime, value[0], global.Tz)
		if err != nil {
			return errors.New("时间选择错误[oildng]")
		}

		if ts.Format("04:05") != "00:00" || te.Format("04:05") != "00:00" || ts.After(te) || ts.Before(allowedStartDate) || te.After(allowedEndDate) {
			return errors.New("时间错误[odfibj]")
		}
		if ts.Hour() == te.Hour() {
			return errors.New("前后时间不能相同[dojinv]")
		}
		if !(ts.Hour() >= minTime && ts.Hour() < maxTime) {
			return fmt.Errorf("开始时间必须%d点 - %d点[oifjd]", minTime, maxTime-1)
		}
		if !(te.Hour() > minTime && te.Hour() <= maxTime) {
			return fmt.Errorf("结束时间必须%d点 - %d点[oifgjka]", minTime+1, maxTime)
		}
		if fts := ts.Format(time.DateOnly); fts != te.Format(time.DateOnly) {
			return errors.New("时间选择错误[odfkng]")
		} else if _, ok := dates[fts]; ok {
			return errors.New("时间选择错误[odfdfskng]")
		} else {
			dates[fts] = true
		}
	}
	return nil
}
