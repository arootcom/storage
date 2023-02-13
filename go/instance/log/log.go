package log

import(
    "os"
    "fmt"
    "sync"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var once sync.Once
var instance log.Logger

//
func getInstance() log.Logger {
    once.Do(func() {
        instance = log.NewLogfmtLogger(os.Stdout)
        instance = log.With(instance, "ts", log.DefaultTimestampUTC)

        switch os.Getenv("STORAGE_LOG_LEVEL") {
            case "All":
	            instance = level.NewFilter(instance, level.AllowAll())
                break
            case "Debug":
	            instance = level.NewFilter(instance, level.AllowDebug())
                break
            case "Error":
	            instance = level.NewFilter(instance, level.AllowError())
                break
            case "Info":
	            instance = level.NewFilter(instance, level.AllowInfo())
                break
            case "None" :
	            instance = level.NewFilter(instance, level.AllowNone())
                break
            default:
	            instance = level.NewFilter(instance, level.AllowInfo())
                break
        }
    })
    return instance
}

//
func Debug (key string, notes ...any) {
    logger := getInstance()
    level.Debug(logger).Log(key, fmt.Sprint(notes))
}

//
func Info (key string, notes ...any) {
    logger := getInstance()
    level.Info(logger).Log(key, fmt.Sprint(notes))
}

//
func Warn (key string, notes ...any) {
    logger := getInstance()
    level.Warn(logger).Log(key, fmt.Sprint(notes))
}

//
func Error (key string, notes ...any) {
    logger := getInstance()
    level.Error(logger).Log(key, fmt.Sprint(notes))
}
