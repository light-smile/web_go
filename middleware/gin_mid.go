package middleware

import (
	"bytes"
	"dnds_go/common"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func disableEscapeHtml(data interface{}) (string, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(data); err != nil {
		return "", err
	}
	return bf.String(), nil
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		// 是否需要显示reqData和resData
		showReqAndRes := viper.GetBool("Logger.ShowReqAndRes")
		// var bodyLogWriter bodyL
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		var reqData []byte
		if showReqAndRes {
			var err error
			c.Writer = bodyLogWriter
			reqData, err = ioutil.ReadAll(c.Request.Body)
			if err != nil {
				logger.Sugar().Errorf("log request body err: %v", err)
			}
			// 回写body,body只可读一次，必须回写，不然业务代码无法拿到body
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqData))
		}

		c.Next()

		if showReqAndRes {
			var responseCode int
			var responseMsg string
			var responseData interface{}
			var jsonData []byte
			responseBody := bodyLogWriter.body.String()

			response := common.Res{}
			if responseBody != "" {
				err := json.Unmarshal([]byte(responseBody), &response)
				if err == nil {
					responseCode = response.Code
					responseMsg = response.Msg
					responseData = response.Data
				}
			}
			jsonData, err := json.MarshalIndent(responseData, "", "\t")
			if err != nil {
				logger.Sugar().Errorf("response data marshal err: %s", err.Error())
			}
			cost := time.Since(start)
			logger.Info(path+"\nreq_data: "+string(reqData)+"\nres_data: "+string(jsonData),
				zap.String("cost", cost.String()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.Int("status", c.Writer.Status()),
				zap.String("query", query),
				// zap.String("body", string(data)),
				zap.String("ip", c.ClientIP()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Int("res_code", responseCode),
				zap.String("res_msg", responseMsg),
			)

		} else {
			cost := time.Since(start)
			logger.Info(path,
				zap.String("cost", cost.String()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.Int("status", c.Writer.Status()),
				zap.String("query", query),
				// zap.String("body", string(data)),
				zap.String("ip", c.ClientIP()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			)

		}

	}
}
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}
				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
