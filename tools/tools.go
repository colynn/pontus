package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/colynn/pontus/internal/pkg/customerror"

	"github.com/google/uuid"
	log "unknwon.dev/clog/v2"
)

// StrToInt ..
func StrToInt(err error, index string) int {
	result, err := strconv.Atoi(index)
	if err != nil {
		customerror.HasError(err, "string to int error"+err.Error(), -1)
	}
	return result
}

// StringToInt64 ..
func StringToInt64(e string) (int64, error) {
	return strconv.ParseInt(e, 10, 64)
}

// StringToFloat32 ..
func StringToFloat32(e string) (float32, error) {
	float64, err := strconv.ParseFloat(e, 32)
	return float32(float64), err
}

// StringToInt ..
func StringToInt(e string) (int, error) {
	return strconv.Atoi(e)
}

// IntToString ..
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// GetUniqString ..
func GetUniqString() string {
	timeStr := time.Now().Format("2006-0102-1504")
	u := uuid.New()
	return fmt.Sprintf("%s_%d", timeStr, u.ID())
}

// ParseStrToDate ..
func ParseStrToDate(timestr string) (time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	timeValue, err := time.ParseInLocation("2006-01-02", timestr, loc)
	return timeValue, err
}

// GetCurrentTime ..
func GetCurrentTime(loc string) time.Time {
	if loc == "" {
		loc = "Asia/Shanghai"
	}
	locd, _ := time.LoadLocation(loc)
	return time.Now().In(locd)
}

// ParseStrToDateTime ..
func ParseStrToDateTime(timestr string) (time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	timeValue, err := time.ParseInLocation("2006-01-02 15:04:05", timestr, loc)
	return timeValue, err
}

// EnsureAbs prepends the WorkDir to the given path if it is not an absolute path.
func EnsureAbs(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(WorkDir(), path)
}

var (
	workDir     string
	workDirOnce sync.Once
)

// WorkDir returns the absolute path of work directory. It reads the value of envrionment
// variable GOGS_WORK_DIR. When not set, it uses the directory where the application's
// binary is located.
func WorkDir() string {
	workDirOnce.Do(func() {
		workDir = filepath.Dir(AppPath())
	})

	return workDir
}

var (
	appPath     string
	appPathOnce sync.Once
)

// AppPath returns the absolute path of the application's binary.
func AppPath() string {
	appPathOnce.Do(func() {
		var err error
		appPath, err = exec.LookPath(os.Args[0])
		if err != nil {
			panic("look executable path: " + err.Error())
		}
		log.Trace("appPath: %s", appPath)

		appPath, err = filepath.Abs(appPath)
		if err != nil {
			panic("get absolute executable path: " + err.Error())
		}
	})

	return appPath
}
