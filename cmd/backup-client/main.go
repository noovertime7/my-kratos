package main

import (
	"backup-client/conf"
	"backup-client/pkg/logger"
	"flag"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/log"
	"net/url"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
}

// 使用配置文件中提供的地址
func buildEndpoints(conf *conf.Bootstrap) []*url.URL {
	grpcEndpoint := &url.URL{
		Scheme: "grpc",
		Host:   conf.Registry.GetGrpcServer(),
	}
	httpEndPoint := &url.URL{
		Scheme: "http",
		Host:   conf.Registry.GetHttpServer(),
	}
	return []*url.URL{grpcEndpoint, httpEndPoint}
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, registry *etcd.Registry, conf *conf.Bootstrap) *kratos.App {
	if conf.Name == "" {
		panic("please enter the service name")
	}
	return kratos.New(
		kratos.ID(id),
		kratos.Name(conf.Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{
			"app-name": "backup-client",
		}),
		kratos.Logger(logger),
		kratos.Registrar(registry),
		kratos.Endpoint(buildEndpoints(conf)...),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func registerLog(cfg *conf.Bootstrap) (log.Logger, func() error) {
	zapLogger := logger.NewLogger(logger.NewZapLogger(cfg))
	return log.With(zapLogger,
		//"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		//"service.version", Version,
		//"traceID", tracing.TraceID(),
		//"spanID", tracing.SpanID(),
	), zapLogger.Sync
}

func main() {
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		log.Fatalf("加载配置文件出错，请检查! [%v]", err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		log.Fatalf("加载配置文件出错，请检查! [%v]", err)
	}

	// 注册logger
	lg, syncFunc := registerLog(&bc)
	defer syncFunc()

	app, cleanup, err := wireApp(&bc, bc.Server, bc.Data, lg)
	if err != nil {
		log.Fatalf("程序初始化失败，请检查! [%v]", err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		log.Fatalf("程序启动失败，请检查! [%v]", err)
	}
}
