package middleware

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// Logging middleware logs request details with context info
func Logging(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if tr, ok := transport.FromServerContext(ctx); ok {
				kind = tr.Kind().String()
				operation = tr.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level := log.LevelInfo
			if err != nil {
				level = log.LevelError
			}

			_ = log.WithContext(ctx, logger).Log(level,
				"kind", kind,
				"operation", operation,
				"args", sanitizeArgs(req),
				"code", code,
				"reason", reason,
				"latency", time.Since(startTime).String(),
				"user_id", GetUserID(ctx),
				"tenant_id", GetTenantID(ctx),
			)
			return reply, err
		}
	}
}

func sanitizeArgs(req interface{}) interface{} {
	return sanitizeValue(reflect.ValueOf(req), "")
}

func sanitizeValue(v reflect.Value, fieldName string) interface{} {
	if !v.IsValid() {
		return nil
	}
	if isSensitiveField(fieldName) {
		if v.IsZero() {
			return ""
		}
		return "[REDACTED]"
	}
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		out := make(map[string]interface{})
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				continue
			}
			out[field.Name] = sanitizeValue(v.Field(i), field.Name)
		}
		if len(out) > 0 {
			return out
		}
	case reflect.Map, reflect.Slice, reflect.Array:
		return fmt.Sprintf("%s(len=%d)", v.Type().String(), v.Len())
	default:
		if v.CanInterface() {
			return v.Interface()
		}
	}
	return fmt.Sprintf("<%s>", v.Type().String())
}

func isSensitiveField(name string) bool {
	lower := strings.ToLower(name)
	for _, part := range []string{"password", "token", "secret", "credential", "authorization", "pat", "apikey", "api_key", "privatekey", "private_key"} {
		if strings.Contains(lower, part) {
			return true
		}
	}
	return false
}
