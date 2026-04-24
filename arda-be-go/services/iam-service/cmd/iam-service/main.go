package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	// Mặc định đọc configs/dev để local chạy "kratos run" không bị lỗi
	// Khi deploy K8s, Dockerfile CMD đã ghi đè bằng "-conf /data/conf"
	flag.StringVar(&flagconf, "conf", "configs/config.yaml", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()

	// Tự động chuyển đổi flagconf sang đường dẫn tuyệt đối nếu nó là tương đối
	if !filepath.IsAbs(flagconf) {
		if cwd, err := os.Getwd(); err == nil {
			path := filepath.Join(cwd, flagconf)
			// Nếu không tìm thấy ở thư mục hiện tại (ví dụ khi chạy trong cmd/iam-service)
			// thì thử tìm ở thư mục cha
			if _, err := os.Stat(path); os.IsNotExist(err) {
				path = filepath.Join(filepath.Dir(cwd), flagconf)
				// Thử thêm một cấp nữa nếu vẫn không thấy
				if _, err := os.Stat(path); os.IsNotExist(err) {
					path = filepath.Join(filepath.Dir(filepath.Dir(cwd)), flagconf)
				}
			}
			flagconf = path
		}
	}

	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
			env.NewSource(""),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// Override from Environment Variables for local development flexibility
	if envDB := os.Getenv("DATABASE_URL"); envDB != "" {
		bc.Data.Database.Source = envDB
	}
	if envRedis := os.Getenv("REDIS_ADDR"); envRedis != "" {
		bc.Data.Redis.Addr = envRedis
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.Jwt, bc.Zitadel, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
