package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/twbworld/dating/global"
)

func TestBase64EncodeDecode(t *testing.T) {
	original := "hello world"
	encoded := Base64Encode(original)
	decoded := Base64Decode(encoded)
	assert.Equal(t, original, decoded)
}

func TestHash(t *testing.T) {
	original := "hello world"
	hashed := Hash(original)
	expected := "2f05477fc24bb4faefd86517156dafdecec45b8ad3cf2522a563582b"
	assert.Equal(t, expected, hashed)
}

func TestTimeFormat(t *testing.T) {
	global.Tz, _ = time.LoadLocation("Asia/Shanghai")
	timestamp := int64(1700000000)
	formatted := TimeFormat(timestamp)
	expected := time.Unix(timestamp, 0).In(global.Tz).Format(time.DateTime)
	assert.Equal(t, expected, formatted)
}

func TestNumberFormat(t *testing.T) {
	number := 123.456789
	formatted := NumberFormat(number, 2)
	expected := 123.46
	assert.Equal(t, expected, formatted)
}

func TestFileExist(t *testing.T) {
	path := "testfile.txt"
	file, err := os.Create(path)
	assert.NoError(t, err)
	file.Close()
	defer os.Remove(path)
	assert.True(t, FileExist(path))
}

func TestMkdirAndCreateFile(t *testing.T) {
	dir := "testdir"
	filePath := filepath.Join(dir, "testfile.txt")
	defer os.RemoveAll(dir)

	err := Mkdir(filePath)
	assert.NoError(t, err)
	assert.True(t, FileExist(dir))

	err = CreateFile(filePath)
	assert.NoError(t, err)
	assert.True(t, FileExist(filePath))
}

func TestListToMap(t *testing.T) {
	type Item struct {
		Key   string
		Value string
	}
	list := []Item{
		{Key: "a", Value: "1"},
		{Key: "b", Value: "2"},
	}
	result := ListToMap(list, "Key")
	expected := map[string]interface{}{
		"a": Item{Key: "a", Value: "1"},
		"b": Item{Key: "b", Value: "2"},
	}
	assert.Equal(t, expected, result)
}

func TestInSlice(t *testing.T) {
	slice := []string{"a", "b", "c"}
	assert.Equal(t, 1, InSlice(slice, "b"))
	assert.Equal(t, -1, InSlice(slice, "d"))
}

func TestUnixGroup(t *testing.T) {
	times := []int{1700000000, 1700003600, 1700080000}
	result := UnixGroup(times)
	expected := [][]int{
		{1700000000, 1700003600},
		{1700080000},
	}
	assert.Equal(t, expected, result)
}

func TestSpreadPeriodToHour(t *testing.T) {
	start := int64(1700000000)
	end := int64(1700007200)
	result := SpreadPeriodToHour(start, end)
	expected := []int64{1700000000, 1700003600}
	assert.Equal(t, expected, result)
}

func TestReadyFile(t *testing.T) {
	dir, file := ReadyFile(".txt")
	assert.True(t, strings.HasPrefix(dir, global.Config.StaticDir))
	assert.True(t, strings.HasSuffix(file, ".txt"))
	assert.Equal(t, 10, len(file)-len(".txt"))

	dir, file = ReadyFile()
	assert.True(t, strings.HasPrefix(dir, global.Config.StaticDir))
	assert.Equal(t, 10, len(file))
}
