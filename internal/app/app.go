package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"sonymimic1/Golang_server/checkRTP/config"
	"sonymimic1/Golang_server/checkRTP/pkg/datasource/redis/conn"
	"sonymimic1/Golang_server/checkRTP/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

type App struct {
	redisClient *redis.ClusterClient // Redis連接
	ctx         context.Context      // 上下文
	router      *gin.Engine          // 路由
	cfg         config.Config        //配置
	cron        *cron.Cron           // 定時任務
	logger      *zap.Logger          // 日誌
}

func NewApp(ctx context.Context, cfg config.Config) *App {
	return &App{
		redisClient: conn.SetupRedisConnection(ctx, cfg.Redis),
		ctx:         ctx,
		router:      gin.Default(),
		cfg:         cfg,
		cron:        cron.New(),
		logger: logger.NewLogger(&logger.Config{
			Enable:     cfg.Log.Enable,
			Level:      cfg.Log.Level,
			MaxSize:    cfg.Log.FileSizeMega,
			MaxBackups: cfg.Log.FileCount,
			MaxAge:     cfg.Log.KeepDays,
			FilenPath:  cfg.Log.Path,
			Filename:   "checkRTP.log",
		}, nil),
	}

}

// 啟動應用程式
func (app *App) Run() error {

	// 設定路由服務項目
	if err := app.setService(); err != nil {
		return err
	}

	// 啟動服務
	app.cron.Start()
	app.startHTTPServer()

	return nil
}

// 啟動HTTP服務
func (app *App) startHTTPServer() {

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(app.cfg.App.Port),
		Handler: app.router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	conn.CloseRedisConnection(app.redisClient)
	app.cron.Stop()

	log.Println("Server exiting")
}
