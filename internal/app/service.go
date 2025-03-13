package app

import (
	"log"
	v1 "sonymimic1/Golang_server/checkRTP/internal/controller/http/v1"
	"sonymimic1/Golang_server/checkRTP/internal/repo"
	"sonymimic1/Golang_server/checkRTP/internal/usecase"

	"go.uber.org/zap"
)

func (app *App) setService() error {
	//logger := app.logger.

	// 查RTP值服務
	rtpRepo := repo.NewRTPRepository(app.ctx, app.redisClient, app.logger)
	rtpUseCase := usecase.NewRTPUseCase(rtpRepo, app.cfg, app.logger)
	rtpHandler := v1.NewRTPHandler(rtpUseCase, app.cfg)
	healthHandler := v1.NewHealthHandler(app.redisClient)

	rtpsRoutes := app.router.Group("rtp")
	{
		rtpsRoutes.GET("/all", rtpHandler.GetAllRTP)
		rtpsRoutes.GET("/", rtpHandler.GetRTPByGameCode) // rtp?gamecode=AB3
	}

	// 清除BET/WIN
	clearRoutes := app.router.Group("clear")
	{
		clearRoutes.PUT("/all", rtpHandler.ClearRTPAll)
		clearRoutes.PUT("/", rtpHandler.ClearRTPByGameCode) // clear?gamecode=AB3
	}

	// 健康檢查
	healthRoutes := app.router.Group("health")
	{
		healthRoutes.GET("/", healthHandler.CheckHealth)
		healthRoutes.GET("/redis", healthHandler.CheckReidsConnect)
	}

	// 定時排程任務（清除BET/WIN <redis>）
	{
		_, err := app.cron.AddFunc(app.cfg.App.ScheduleCleanRTP, func() {

			logPrefix := "rtpUseCase.ClearRTPsAll"

			_, err := rtpUseCase.ClearRTPsAll()
			if err != nil {
				app.logger.Error(logPrefix, zap.Error(err))
			}

			//must(rtpUseCase.ClearRTPsAll())
		})
		if err != nil {
			logPrefix := "app.cron.AddFunc"
			app.logger.Error(logPrefix, zap.Error(err))
			return err
		}
	}

	return nil
}

// --------------------
// golang write skill
// --------------------
// function for error handling
func Must[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)
	}

	return v
}

// Variadic for optional parameters
func GetRequest(address string, optTimeout ...int) int {
	//Default value
	timeout := 30
	//If the value is set we override the default
	if len(optTimeout) > 0 {
		timeout = optTimeout[0]
	}
	//Do something with the timeout
	return timeout
}
