package MTGoLogger

import (
	"context"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

func NewLoggerWithoutConfigFile(lev string, outputPaths []string, appends []string) (*AssetLog, error) {
	level, err := getLevel(lev)
	if err != nil {
		return nil, err
	}

	zapWriteSyncers := make([]zapcore.WriteSyncer, 0)

	for _, fileName := range outputPaths {
		zapWriteSyncers = append(zapWriteSyncers, getWriteSyncers(fileName))
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zap.CombineWriteSyncers(zapWriteSyncers...), level)

	logger := zap.New(core)
	sugarLogger := logger.Sugar()
	lc := &LoggerConfig{
		Level:       lev,
		OutputPaths: outputPaths,
		Appends:     appends,
	}
	al := &AssetLog{sugarLogger: sugarLogger, conf: lc}
	return al, nil
}

func NewLoggerWithConfigFile(env string, loggerFilePath string) (*AssetLog, error) {
	if env == "" {
		return nil, errors.New("Env Not set")
	}

	if env != "prod" && env != "test" {
		env = "dev"
	}

	loggerConfigs := new(LoggerConfigs)
	yamlData, err := ioutil.ReadFile(loggerFilePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlData, loggerConfigs)
	if err != nil {
		return nil, err
	}

	conf := loggerConfigs.Configs[env]

	if conf == nil {
		return nil, errors.New("Error in logger config file")
	}

	level, err := getLevel(conf.Level)
	if err != nil {
		return nil, err
	}

	zapWriteSyncers := make([]zapcore.WriteSyncer, 0)

	for _, fileName := range conf.OutputPaths {
		zapWriteSyncers = append(zapWriteSyncers, getWriteSyncers(fileName))
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zap.CombineWriteSyncers(zapWriteSyncers...), level)

	logger := zap.New(core)
	sugarLogger := logger.Sugar()
	return &AssetLog{sugarLogger: sugarLogger, conf: conf}, nil
}

func getWriteSyncers(fileName string) zapcore.WriteSyncer {
	if strings.Compare(fileName, "stdout") == 0 {
		return zapcore.AddSync(os.Stdout)
	}

	if strings.Compare(fileName, "stderr") == 0 {
		return zapcore.AddSync(os.Stderr)
	}

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    50, // megabytes
		MaxBackups: 1,
		MaxAge:     1, // days
	})
}

func getLevel(level string) (zapcore.Level, error) {
	level = strings.ToUpper(level)
	switch level {
	case "DEBUG":
		return zap.DebugLevel, nil
	case "INFO":
		return zap.InfoLevel, nil
	case "WARN":
		return zap.WarnLevel, nil
	case "ERROR":
		return zap.ErrorLevel, nil
	case "FATAL":
		return zap.FatalLevel, nil
	default:
		return zapcore.PanicLevel, errors.New("unsupported level")
	}
}

func getAppendMapFromContext(ctx context.Context, appends []string) map[string]string {
	if ctx == nil {
		return nil
	}
	appendMap := make(map[string]string)
	for _, append := range appends {
		var finalVal string
		ctxVal := ctx.Value(append)
		if ctxVal == nil {
			finalVal = ""
		} else {
			finalVal = ctxVal.(string)
		}
		appendMap[append] = finalVal
	}
	return appendMap
}

func getAppendedLogger(logger *zap.SugaredLogger, appendMap map[string]string) *zap.SugaredLogger {
	for k, v := range appendMap {
		logger = logger.With(k, v)
	}
	return logger
}

func (logger AssetLog) Debug(ctx context.Context, args ...interface{}) {
	appendMap := getAppendMapFromContext(ctx, logger.conf.Appends)
	appendedLogger := getAppendedLogger(logger.sugarLogger, appendMap)
	appendedLogger.Debug(args...)
}

func (logger AssetLog) Debugf(ctx context.Context, format string, args ...interface{}) {
	appendMap := getAppendMapFromContext(ctx, logger.conf.Appends)
	appendedLogger := getAppendedLogger(logger.sugarLogger, appendMap)
	appendedLogger.Debugf(format, args...)
}

func (logger AssetLog) Info(ctx context.Context, args ...interface{}) {
	appendMap := getAppendMapFromContext(ctx, logger.conf.Appends)
	appendedLogger := getAppendedLogger(logger.sugarLogger, appendMap)
	appendedLogger.Info(args...)
}

func (logger AssetLog) Infof(ctx context.Context, format string, args ...interface{}) {
	appendMap := getAppendMapFromContext(ctx, logger.conf.Appends)
	appendedLogger := getAppendedLogger(logger.sugarLogger, appendMap)
	appendedLogger.Infof(format, args...)
}

func (logger AssetLog) Fatal(ctx context.Context, args ...interface{}) {
	appendMap := getAppendMapFromContext(ctx, logger.conf.Appends)
	appendedLogger := getAppendedLogger(logger.sugarLogger, appendMap)
	appendedLogger.Fatal(args...)
}

func (logger AssetLog) Fatalf(ctx context.Context, format string, args ...interface{}) {
	appendMap := getAppendMapFromContext(ctx, logger.conf.Appends)
	appendedLogger := getAppendedLogger(logger.sugarLogger, appendMap)
	appendedLogger.Fatalf(format, args...)
}

func (logger AssetLog) Warn(ctx context.Context, args ...interface{}) {
	appendMap := getAppendMapFromContext(ctx, logger.conf.Appends)
	appendedLogger := getAppendedLogger(logger.sugarLogger, appendMap)
	appendedLogger.Warn(args...)
}

func (logger AssetLog) Warningf(ctx context.Context, format string, args ...interface{}) {
	appendMap := getAppendMapFromContext(ctx, logger.conf.Appends)
	appendedLogger := getAppendedLogger(logger.sugarLogger, appendMap)
	appendedLogger.Warnf(format, args...)
}

func (logger AssetLog) Error(ctx context.Context, args ...interface{}) {
	appendMap := getAppendMapFromContext(ctx, logger.conf.Appends)
	appendedLogger := getAppendedLogger(logger.sugarLogger, appendMap)
	appendedLogger.Error(args...)
}

func (logger AssetLog) Errorf(ctx context.Context, format string, args ...interface{}) {
	appendMap := getAppendMapFromContext(ctx, logger.conf.Appends)
	appendedLogger := getAppendedLogger(logger.sugarLogger, appendMap)
	appendedLogger.Errorf(format, args...)
}
