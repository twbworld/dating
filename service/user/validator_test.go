package user

import (
	"mime/multipart"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/model/common"
)

func TestValidatorLoginPost(t *testing.T) {
	validator := &Validator{}

	t.Run("Valid Code", func(t *testing.T) {
		data := &common.LoginPost{Code: "validCode"}
		err := validator.ValidatorLoginPost(data)
		assert.NoError(t, err)
	})

	t.Run("Empty Code", func(t *testing.T) {
		data := &common.LoginPost{Code: ""}
		err := validator.ValidatorLoginPost(data)
		assert.Error(t, err)
	})
}

func TestValidatorUserAddPost(t *testing.T) {
	validator := &Validator{}

	t.Run("Valid Code", func(t *testing.T) {
		data := &common.UserInfoPost{}
		data.Code = "validCode"
		err := validator.ValidatorUserAddPost(data)
		assert.NoError(t, err)
	})

	t.Run("Empty Code", func(t *testing.T) {
		data := &common.UserInfoPost{}
		data.Code = ""
		err := validator.ValidatorUserAddPost(data)
		assert.Error(t, err)
	})
}

func TestValidatorFeedbackPost(t *testing.T) {
	validator := &Validator{}

	t.Run("Valid Desc", func(t *testing.T) {
		data := &common.FeedbackPost{Desc: "valid description"}
		err := validator.ValidatorFeedbackPost(data)
		assert.NoError(t, err)
	})

	t.Run("Empty Desc", func(t *testing.T) {
		data := &common.FeedbackPost{Desc: ""}
		err := validator.ValidatorFeedbackPost(data)
		assert.Error(t, err)
	})

	t.Run("Desc Exceeds 100 Characters", func(t *testing.T) {
		data := &common.FeedbackPost{Desc: strings.Repeat("a", 101)}
		err := validator.ValidatorFeedbackPost(data)
		assert.Error(t, err)
	})
}

func TestValidatorUpload(t *testing.T) {
	validator := &Validator{}

	t.Run("Valid File", func(t *testing.T) {
		file := &multipart.FileHeader{Filename: "test.jpg", Size: 1 << 20}
		err := validator.ValidatorUpload(file)
		assert.NoError(t, err)
	})

	t.Run("Unsupported File Extension", func(t *testing.T) {
		file := &multipart.FileHeader{Filename: "test.txt", Size: 1 << 20}
		err := validator.ValidatorUpload(file)
		assert.Error(t, err)
	})

	t.Run("File Size Exceeds Limit", func(t *testing.T) {
		file := &multipart.FileHeader{Filename: "test.jpg", Size: 3 << 20}
		err := validator.ValidatorUpload(file)
		assert.Error(t, err)
	})
}

func TestValidatorInfo(t *testing.T) {
	validator := &Validator{}
	global.Tz, _ = time.LoadLocation("UTC")
	n := time.Now().In(global.Tz)

	t.Run("Valid Time", func(t *testing.T) {
		data := &common.InfoPost{Time: [][2]string{{time.Date(n.Year(), n.Month(), n.Day(), 8, 0, 0, 0, global.Tz).Format(time.DateTime), time.Date(n.Year(), n.Month(), n.Day(), 23, 0, 0, 0, global.Tz).Format(time.DateTime)}}}
		err := validator.ValidatorInfo(data)
		assert.NoError(t, err)
	})

	t.Run("Empty Time", func(t *testing.T) {
		data := &common.InfoPost{Time: [][2]string{}}
		err := validator.ValidatorInfo(data)
		assert.Error(t, err)
	})

	t.Run("Invalid Date Format", func(t *testing.T) {
		data := &common.InfoPost{Time: [][2]string{{"invalid", "2023-01-01 23:00:00"}}}
		err := validator.ValidatorInfo(data)
		assert.Error(t, err)
	})

	t.Run("Start Time After End Time", func(t *testing.T) {
		data := &common.InfoPost{Time: [][2]string{{time.Date(n.Year(), n.Month(), n.Day(), 23, 0, 0, 0, global.Tz).Format(time.DateTime), time.Date(n.Year(), n.Month(), n.Day(), 8, 0, 0, 0, global.Tz).Format(time.DateTime)}}}
		err := validator.ValidatorInfo(data)
		assert.Error(t, err)
	})

	t.Run("Start Time After End Time", func(t *testing.T) {
		data := &common.InfoPost{Time: [][2]string{{time.Date(n.Year()-2, n.Month(), n.Day(), 8, 0, 0, 0, global.Tz).Format(time.DateTime), time.Date(n.Year()-2, n.Month(), n.Day(), 23, 0, 0, 0, global.Tz).Format(time.DateTime)}}}
		err := validator.ValidatorInfo(data)
		assert.Error(t, err)
	})

	t.Run("Start Time After End Time", func(t *testing.T) {
		data := &common.InfoPost{Time: [][2]string{{time.Date(n.Year(), n.Month(), n.Day(), 7, 0, 0, 0, global.Tz).Format(time.DateTime), time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, global.Tz).Format(time.DateTime)}}}
		err := validator.ValidatorInfo(data)
		assert.Error(t, err)
	})

}
