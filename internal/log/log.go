package log

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

//This is to avoid golint error
type contextKey string

const projPath = "github.com/keikoproj/ice-kube"

var requestID contextKey = "request_id"

//New function will set the level
func New(debug bool) {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetReportCaller(false)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			//repopath := fmt.Sprintf(projPath, os.Getenv("GOPATH"))
			filename := strings.SplitAfter(f.File, projPath)
			fnName := strings.SplitAfter(f.Function, projPath)
			return fmt.Sprintf("%s()", fnName[len(fnName)-1]), fmt.Sprintf("%s:%d", filename[len(filename)-1], f.Line)
		},
		FullTimestamp: true,
	})
}

//Logger function provides default entries for every logrus log statement
func Logger(ctx ...context.Context) *logrus.Entry {
	lr := &logrus.Entry{}
	if len(ctx) != 0 {
		lr = logrus.WithFields(logrus.Fields{"request_id": ctx[0].Value(requestID)})
	} else {
		lr = logrus.WithFields(logrus.Fields{})
	}
	return lr
}

// NewContext returns a new Context that carries value u.
func NewContext(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, requestID, value)
}
