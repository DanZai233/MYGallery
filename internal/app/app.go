package app

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/mygallery/mygallery/internal/config"
	"github.com/mygallery/mygallery/internal/database"
	"github.com/mygallery/mygallery/internal/router"
	"github.com/mygallery/mygallery/internal/storage"
)

// Application wires together configuration, infrastructure, and HTTP routing.
type Application struct {
	cfg         *config.Config
	router      *gin.Engine
	logger      *log.Logger
	warnings    []string
	shutdownFns []func(context.Context) error
	options     options
}

// New constructs a fully initialised Application instance.
func New(opts ...Option) (*Application, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	cfg, err := loadConfig(options)
	if err != nil {
		return nil, err
	}

	application := &Application{
		cfg:     cfg,
		logger:  options.logger,
		options: options,
	}

	if err := application.bootstrapInfrastructure(); err != nil {
		return nil, err
	}

	application.router = router.SetupRouter(cfg)

	return application, nil
}

// Run launches the HTTP server and blocks until it exits.
func (a *Application) Run() error {
	address := fmt.Sprintf("%s:%d", a.cfg.Server.Host, a.cfg.Server.Port)
	return a.router.Run(address)
}

// Logger returns the logger associated with the application.
func (a *Application) Logger() *log.Logger {
	return a.logger
}

// Config exposes the loaded configuration.
func (a *Application) Config() *config.Config {
	return a.cfg
}

// Router exposes the configured Gin router.
func (a *Application) Router() *gin.Engine {
	return a.router
}

// Warnings returns bootstrap warnings (non-fatal issues surfaced to the caller).
func (a *Application) Warnings() []string {
	out := make([]string, len(a.warnings))
	copy(out, a.warnings)
	return out
}

// RegisterShutdown adds a cleanup hook that runs with a context on shutdown.
func (a *Application) RegisterShutdown(fn func(context.Context) error) {
	a.shutdownFns = append(a.shutdownFns, fn)
}

// Shutdown executes registered cleanup hooks.
func (a *Application) Shutdown(ctx context.Context) error {
	var combined error
	for i := len(a.shutdownFns) - 1; i >= 0; i-- {
		if err := a.shutdownFns[i](ctx); err != nil {
			combined = errors.Join(combined, err)
		}
	}
	return combined
}

func (a *Application) bootstrapInfrastructure() error {
	if err := database.InitDB(a.cfg); err != nil {
		return fmt.Errorf("database initialisation failed: %w", err)
	}
	a.logger.Println("✓ 数据库初始化成功")

	if err := storage.InitStorage(a.cfg); err != nil {
		return fmt.Errorf("storage initialisation failed: %w", err)
	}
	a.logger.Println("✓ 存储初始化成功")

	if err := database.CreateDefaultAdmin(a.cfg); err != nil {
		warning := fmt.Sprintf("创建管理员账号失败: %v", err)
		a.logger.Printf("⚠️ %s", warning)
		a.warnings = append(a.warnings, warning)
	} else {
		a.logger.Println("✓ 管理员账号已就绪")
	}

	return nil
}

func loadConfig(opts options) (*config.Config, error) {
	if opts.configOverride != nil {
		return opts.configOverride, nil
	}

	path := opts.configPath
	if path == "" {
		path = "config.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			opts.logger.Printf("警告: 未找到配置文件 %s, 使用默认配置", path)
			return config.DefaultConfig(), nil
		}
		return nil, fmt.Errorf("读取配置文件 %s 失败: %w", path, err)
	}

	cfg, err := config.FromBytes(data)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件 %s 失败: %w", path, err)
	}

	return cfg, nil
}
