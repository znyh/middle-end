package log

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    slog "log"
    "os"
    "strings"
    "testing"

    "github.com/go-kratos/kratos/v2/log"
    "github.com/sirupsen/logrus"
)

func TestLoggerLog(t *testing.T) {
    tests := map[string]struct {
        level     logrus.Level
        formatter logrus.Formatter
        logLevel  log.Level
        kvs       []interface{}
        want      string
    }{
        "json format": {
            level:     logrus.InfoLevel,
            formatter: &logrus.JSONFormatter{},
            logLevel:  log.LevelInfo,
            kvs:       []interface{}{"case", "json format", "msg", "1"},
            want:      `{"case":"json format","level":"info","msg":"1"`,
        },
        "level unmatch": {
            level:     logrus.InfoLevel,
            formatter: &logrus.JSONFormatter{},
            logLevel:  log.LevelDebug,
            kvs:       []interface{}{"case", "level unmatch", "msg", "1"},
            want:      "",
        },
        "no tags": {
            level:     logrus.InfoLevel,
            formatter: &logrus.JSONFormatter{},
            logLevel:  log.LevelInfo,
            kvs:       []interface{}{"msg", "1"},
            want:      `{"level":"info","msg":"1"`,
        },
    }

    for name, test := range tests {
        t.Run(name, func(t *testing.T) {
            output := new(bytes.Buffer)
            logger := NewLogrusLogger(Level(test.level), Formatter(test.formatter), Output(output))
            _ = logger.Log(test.logLevel, test.kvs...)
            if !strings.HasPrefix(output.String(), test.want) {
                t.Errorf("strings.HasPrefix(output.String(), test.want) got %v want: %v", strings.HasPrefix(output.String(), test.want), true)
            }
        })
    }
}

type LogFormatter struct {
}

func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    //timestamp := //日期格式实现
    //        entry.Caller //调用者信息
    //entry.Data //withfiled方式传入的数据
    //msg := 自己想要的格式输出内容
    //return []byte(msg), nil

    return nil, nil
}

func TestFormatter(t *testing.T) {
    logrus.Infof("hello wolrd")

    logrus.SetFormatter(&logrus.TextFormatter{
        ForceColors:               true,
        DisableColors:             false,
        ForceQuote:                false,
        DisableQuote:              false,
        EnvironmentOverrideColors: true,
        DisableTimestamp:          false,
        FullTimestamp:             true,
        TimestampFormat:           "",
        DisableSorting:            false,
        SortingFunc:               nil,
        DisableLevelTruncation:    false,
        PadLevelText:              false,
        QuoteEmptyFields:          true,
        FieldMap:                  nil,
        CallerPrettyfier:          nil,
    })
    logrus.SetOutput(os.Stdout)

    logrus.WithFields(logrus.Fields{"animal": "walrus"}).Info("A walrus appears")

    logrus.SetReportCaller(true)

    logrus.WithFields(logrus.Fields{"animal": "walrus"}).Info("A walrus appears")

    err := errors.New("json err")
    logrus.WithError(err).Errorf("playerId:%+v", 123)

    logrus.WithField("msg", 123541).Infof("playerId:%+v", 123)

    //entry := logrus.WithFields(map[string]interface{}{
    //    "playerID":          1254,
    //    "ctrType":           "fdsfdsaf",
    //    "totalLossWinScore": "lossWinScore",
    //    "fee":               "fee",
    //    "key":               "key",
    //    "isHitSupply":       "param.isHitSupply",
    //})
    //entry.Infof("检查更新玩家控牌统计数据...")

    slog.SetFlags(slog.Llongfile) //slog.Llongfile ,

    slog.Printf("helllllll")
}

type location struct {
    name string
    age  int
}

func Test_A(t *testing.T) {
    v := &location{
        name: "abc",
        age:  10,
    }

    b, _ := json.Marshal(v)
    fmt.Println(string(b))

    temp := &location{}
    s := json.Unmarshal(b, temp)
}
