package initialize

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/colynn/pontus/config"

	"github.com/colynn/pontus/tools"

	"github.com/spf13/viper"
	log "unknwon.dev/clog/v2"
)

type loggerConf struct {
	Buffer int64
	Config interface{}
}

type logConf struct {
	RootPath string
	Modes    []string
	Configs  []*loggerConf
}

// Log settings
var Log *logConf

// initLogConf returns parsed logging configuration from given viper conf file.
// NOTE: Because we always create a console logger as the primary logger at init time,
// we need to remove it in case the user doesn't configure to use it after the logging
// service is initalized.
func initLogConf(cfg *viper.Viper) (_ *logConf, hasConsole bool, _ error) {
	rootPath := cfg.GetString("log.rootPath")

	modes := strings.Split(cfg.GetString("log.mode"), ",")
	lc := &logConf{
		RootPath: tools.EnsureAbs(rootPath),
		Modes:    make([]string, 0, len(modes)),
		Configs:  make([]*loggerConf, 0, len(modes)),
	}

	// Iterate over [log.*] sections to initialize individual logger.
	levelMappings := map[string]log.Level{
		"trace": log.LevelTrace,
		"info":  log.LevelInfo,
		"warn":  log.LevelWarn,
		"error": log.LevelError,
		"fatal": log.LevelFatal,
	}

	for i := range modes {
		modes[i] = strings.ToLower(strings.TrimSpace(modes[i]))
		secName := "log." + modes[i]
		// TODO: verify conf whether exist or not
		// sec := cfg.GetStringMapString(secName)
		// if err != nil {
		// 	return nil, hasConsole, errors.Errorf("missing configuration section [%s] for %q logger", secName, modes[i])
		// }
		modeLogLevel := cfg.GetString(fmt.Sprintf("%s.level", secName))
		level := levelMappings[strings.ToLower(modeLogLevel)]
		buffer := int64(100)
		var c *loggerConf
		switch modes[i] {
		case log.DefaultConsoleName:
			hasConsole = true
			c = &loggerConf{
				Buffer: buffer,
				Config: log.ConsoleConfig{
					Level: level,
				},
			}

		case log.DefaultFileName:
			fileName := cfg.GetString("log.file.fileName")
			logPath := filepath.Join(lc.RootPath, fileName)
			c = &loggerConf{
				Buffer: buffer,
				Config: log.FileConfig{
					Level:    level,
					Filename: logPath,
					FileRotationConfig: log.FileRotationConfig{
						Rotate:   true,
						Daily:    true,
						MaxSize:  1 << uint(28),
						MaxLines: 1000000,
						MaxDays:  7,
					},
				},
			}
		default:
			continue
		}

		lc.Modes = append(lc.Modes, modes[i])
		lc.Configs = append(lc.Configs, c)
	}

	return lc, hasConsole, nil
}

// InitLogging initializes the logging service of the application.
func InitLogging(cfg *viper.Viper) {
	logConf, hasConsole, err := initLogConf(cfg)
	if err != nil {
		log.Fatal("Failed to init logging configuration: %v", err)
	}
	defer func() {
		Log = logConf
	}()

	err = os.MkdirAll(logConf.RootPath, os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create log directory: %v", err)
	}

	for i, mode := range logConf.Modes {
		c := logConf.Configs[i]

		var err error
		var level log.Level
		switch mode {
		case log.DefaultConsoleName:
			level = c.Config.(log.ConsoleConfig).Level
			err = log.NewConsole(c.Buffer, c.Config)
		case log.DefaultFileName:
			level = c.Config.(log.FileConfig).Level
			err = log.NewFile(c.Buffer, c.Config)
		default:
			panic("unreachable")
		}

		if err != nil {
			log.Fatal("Failed to init %s logger: %v", mode, err)
			return
		}
		log.Trace("Log mode: %s (%s)", strings.Title(mode), strings.Title(strings.ToLower(level.String())))
	}

	// ⚠️ WARNING: It is only safe to remove the primary logger until
	// there are other loggers that are initialized. Otherwise, the
	// application will print errors to nowhere.
	if !hasConsole {
		log.Remove(log.DefaultConsoleName)
	}
}

// InitLogger 初始化日志模块
func InitLogger() {
	c := config.GetConfig()
	logLevel := c.GetInt("log.level")
	logPath := c.GetString("log.fileName")
	err := log.NewFile(1000, log.FileConfig{
		Level:    log.Level(logLevel),
		Filename: logPath,
		FileRotationConfig: log.FileRotationConfig{
			Rotate: true,
			Daily:  true,
		},
	})
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}

	err = log.NewConsole(
		1000,
		log.ConsoleConfig{
			Level: log.Level(logLevel),
		},
	)
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}
