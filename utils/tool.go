package utils

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram/auth/response"
	"github.com/twbworld/dating/global"
)

type timeNumber interface {
	~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(strings.Trim(str, "\n"))))
}
func Base64Decode(str string) string {
	bstr, err := base64.StdEncoding.DecodeString(strings.TrimSpace(strings.Trim(str, "\n")))
	if err != nil {
		return str
	}
	return string(bstr)
}
func Hash(str string) string {
	b := sha256.Sum224([]byte(str))
	return hex.EncodeToString(b[:])
}

func TimeFormat[T timeNumber](t T) string {
	return time.Unix(int64(t), 0).In(global.Tz).Format(time.DateTime)
}

// 四舍五入保留小数位
func NumberFormat[T ~float32 | ~float64](f T, n ...uint) float64 {
	num := uint(2)
	if len(n) > 0 {
		num = n[0]
	}
	nu := math.Pow(10, float64(num))
	return math.Round(float64(f)*nu) / nu
}

// 文件是否存在
func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// 创建目录
func Mkdir(path string) error {
	// 从路径中取目录
	dir := filepath.Dir(path)
	// 获取信息, 即判断是否存在目录
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 生成目录
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// 创建文件
// 可能存在跨越目录创建文件的风险
func CreateFile(path string) error {
	if FileExist(path) {
		return nil
	}

	if err := Mkdir(path); err != nil {
		return err
	}

	fi, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fi.Close()

	return nil
}

// 类似php的array_column($a, null, 'key')
func ListToMap(list interface{}, key string) map[string]interface{} {
	v := reflect.ValueOf(list)
	if v.Kind() != reflect.Slice {
		return nil
	}

	res := make(map[string]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		itemValue := reflect.ValueOf(item)
		keyValue := itemValue.FieldByName(key)
		if keyValue.IsValid() && keyValue.Kind() == reflect.String {
			res[keyValue.String()] = item
		}
	}

	return res
}

// 判断字符串是否在切片中
func InSlice(slice []string, value string) int {
	//上层尽量使用map, 会更快;

	for i, item := range slice {
		if item == value {
			return i
		}
	}
	return -1
}

// 时间戳按日期分组; 例: {[1707100000,1707100000], [170720000]}
func UnixGroup(times []int) [][]int {
	if len(times) == 0 {
		return [][]int{}
	}
	dateTime, tk := make(map[int][]int), make([]int, 0, len(times)) //用于替map做排序
	for _, val := range times {
		if val < 1 {
			continue
		}
		dt := time.Unix(int64(val), 0).In(global.Tz)
		d := int(time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, global.Tz).Unix())
		if _, ok := dateTime[d]; !ok {
			dateTime[d], tk = make([]int, 0, len(times)), append(tk, d)
		}
		dateTime[d] = append(dateTime[d], val)
	}
	if len(dateTime) == 0 {
		return [][]int{}
	}
	sort.Ints(tk)
	unixGroup := make([][]int, 0, len(tk))
	for _, v := range tk {
		sort.Ints(dateTime[v])
		unixGroup = append(unixGroup, dateTime[v])
	}

	return unixGroup
}

// 打散时间段(粒度为1小时) ; 如: "1-4点" 转为 ["1点", "2点", "3点"] 三个时间段
func SpreadPeriodToHour[T timeNumber](start, end T) []T {
	add := T(3600)
	res := make([]T, 0, (end-start)/add+1)
	for start < end {
		//这不使用"<=", 不算最后的时间戳,是因为: 往后的一个时间戳值,代表当前时间戳的后一小时, 而不是当前秒
		res = append(res, start)
		start += add
	}
	return res
}

// 用Code向微信官方换取openid等信息
// 该函数可能会运行较慢
func AuthWxCode(code string) (rs *response.ResponseCode2Session, err error) {
	if code == "" {
		err = errors.New("系统错误[iodgj]")
		return
	}

	if global.Config.Weixin.XcxAppid == "" {
		err = errors.New("没配置小程序Appid[nbvkpl]")
		return
	}

	rs, err = global.MiniProgramApp.Auth.Session(context.Background(), code)
	if err != nil {
		return
	}
	if rs.OpenID == "" {
		err = errors.New("请刷新[nb09]")
		return
	}
	return
}

// 生成文件路径和文件名
func ReadyFile(fileExt ...string) (string, string) {
	ext := ""
	if len(fileExt) > 0 {
		ext = fileExt[0]
	}

	return filepath.Join("static", time.Now().In(global.Tz).Format("2006/01/")), Hash(strconv.FormatInt(time.Now().UnixNano()+rand.Int63n(100), 10))[:10] + ext
}
