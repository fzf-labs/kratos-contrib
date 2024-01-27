package errx

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/fzf-labs/fpkg/util/fileutil"
	"github.com/go-kratos/kratos/v2/errors"
	http2 "github.com/go-kratos/kratos/v2/transport/http"
)

const (
	HeaderLang = "lang"  // 语言类型
	ZhCN       = "zh-CN" // zh_CN 简体中文-中国
	EnUS       = "en-US" // en_US 英文-美国
)

type Error errors.Error

type ErrorManager struct {
	errs    map[string]*Error
	reasons []string
	i18n    map[string]map[string]string
}

// NewErrorManager 创建错误管理器
func NewErrorManager(opts ...Option) *ErrorManager {
	e := &ErrorManager{
		errs:    make(map[string]*Error),
		reasons: make([]string, 0),
		i18n:    make(map[string]map[string]string),
	}
	if len(opts) > 0 {
		for _, v := range opts {
			v(e)
		}
	}
	return e
}

// Option 配置选项
type Option func(gen *ErrorManager)

// WithI18n 设置国际化
func WithI18n(lang string, m map[string]string) Option {
	return func(e *ErrorManager) {
		e.i18n[lang] = m
	}
}

// New 创建错误
func (e *ErrorManager) New(code int, reason, message string) *Error {
	_, ok := e.errs[reason]
	if ok {
		panic(fmt.Sprintf("reason %s is exsit, please change one", reason))
	}
	err := (*Error)(errors.New(code, reason, message))
	e.errs[reason] = err
	e.reasons = append(e.reasons, reason)
	return err
}

// WithFmtMsg 设置自定义消息
func (e *Error) WithFmtMsg(msg ...string) *Error {
	if len(msg) == 0 {
		return e
	}
	err := (*errors.Error)(e)
	metadata := err.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}
	metadata["fmt"] = strings.Join(msg, ",")
	err = errors.Clone(err).WithMetadata(metadata)
	return (*Error)(err)
}

// WithError 设置错误
func (e *Error) WithError(n error) *Error {
	if n == nil {
		return e
	}
	err := (*errors.Error)(e)
	metadata := err.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}
	metadata["cause"] = n.Error()
	err = err.WithCause(n).WithMetadata(metadata)
	return (*Error)(err)
}

// Err 错误
func (e *Error) Err() *errors.Error {
	err := (*errors.Error)(e)
	metadata := err.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}
	metadata["line"] = getFileLine()
	err = err.WithMetadata(metadata)
	return err
}

// setErrMetadata 设置错误元数据
func setErrMetadata(err error) map[string]string {
	return map[string]string{
		"cause": err.Error(),
	}
}

// getFileLine 获取行数
func getFileLine() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	return file + ":" + strconv.Itoa(line)
}

// getLanguages 获取语言列表
func (e *ErrorManager) getLanguages() []string {
	languages := make([]string, 0)
	for k := range e.i18n {
		languages = append(languages, k)
	}
	sort.Strings(languages)
	return languages
}

// HTTPErrorEncoder 错误编码
func HTTPErrorEncoder(m *ErrorManager) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		se := errors.FromError(err)
		if se != nil {
			if se.Code == http.StatusInternalServerError && se.Reason == "" {
				se = errors.New(http.StatusInternalServerError, "InternalServerError", "internal server error").WithCause(err).WithMetadata(setErrMetadata(err))
			}
			message := GetMessage(m, se, r.Header.Get(HeaderLang))
			if message != "" {
				se.Message = message
			}
		}
		codec, _ := http2.CodecForRequest(r, "Accept")
		body, err := codec.Marshal(se)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/"+codec.Name())
		w.WriteHeader(int(se.Code))
		_, _ = w.Write(body)
	}
}

// GetMessage 获取错误消息
func GetMessage(manager *ErrorManager, err *errors.Error, lang string) string {
	if lang == "" {
		lang = ZhCN
	}
	msg := ""
	if lang == ZhCN {
		msg = err.GetMessage()
	} else {
		msg = manager.i18n[lang][err.Reason]
	}
	if err.Metadata["fmt"] != "" {
		msg = fmt.Sprintf(msg, err.Metadata["fmt"])
		delete(err.Metadata, "fmt")
	}
	return msg
}

// Export 导出错误码
func (e *ErrorManager) Export(path string) {
	list := make([]map[string]string, 0)
	for _, v := range e.reasons {
		m := make(map[string]string)
		m["http_code"] = strconv.Itoa(int(e.errs[v].GetCode()))
		m["reason"] = e.errs[v].GetReason()
		m["message"] = e.errs[v].GetMessage()
		for _, lang := range e.getLanguages() {
			m[lang] = GetMessage(e, e.errs[v].Err(), lang)
		}
		list = append(list, m)
	}
	if len(list) > 0 {
		first := `|http_code|reason|message|`
		second := `|--|--|--|`
		for _, v := range e.getLanguages() {
			first += v + `|`
			second += `--|`
		}
		str := fmt.Sprintln(first) + fmt.Sprintln(second)
		for _, m := range list {
			tmpStr := `|` + m["http_code"] + `|` + m["reason"] + `|` + m["message"] + `|`
			for _, v := range e.getLanguages() {
				tmpStr += m[v] + `|`
			}
			str += fmt.Sprintln(tmpStr)
		}
		err := fileutil.WriteContentCover(filepath.Join(path, "code.md"), str)
		if err != nil {
			return
		}
		slog.Info("错误码MARKDOWN导出成功")
	} else {
		slog.Info("错误码不存在!!!")
	}
}
