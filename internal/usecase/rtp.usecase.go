package usecase

import (
	"errors"
	"math"

	"sonymimic1/Golang_server/checkRTP/config"
	"sonymimic1/Golang_server/checkRTP/internal/dto"
	"sonymimic1/Golang_server/checkRTP/internal/entity"
	"sonymimic1/Golang_server/checkRTP/internal/repo"

	"go.uber.org/zap"
)

type RTPUseCase interface {
	FindRTPByGameCode(gameCode string) (*dto.RTPsResponse, error)
	FindRTPsAll() (*dto.RTPsResponse, error)
	ClearRTPsAll() (*dto.ClearResponse, error)
	ClearRTPsByGameCode(gameCode string) (*dto.ClearResponse, error)
}

type rtpUseCase struct {
	rtpRepo     repo.RTPRepository
	cfg         config.Config
	log         *zap.Logger
	cfgGameList map[string]bool // 先將config的gameCode用map存起來
}

func NewRTPUseCase(userRepo repo.RTPRepository, cfg config.Config, log *zap.Logger) RTPUseCase {
	rtpUsecase := &rtpUseCase{
		rtpRepo:     userRepo,
		cfg:         cfg,
		log:         log,
		cfgGameList: make(map[string]bool),
	}

	for _, code := range cfg.App.SearchGameCode {
		rtpUsecase.cfgGameList[code] = true
	}

	return rtpUsecase
}

func (u *rtpUseCase) FindRTPByGameCode(gameCode string) (*dto.RTPsResponse, error) {

	logPrefix := "rtpUseCase.FindRTPByGameCode"

	if _, exists := u.cfgGameList[gameCode]; !exists {
		u.log.Error(logPrefix,
			zap.String("gameCode", gameCode),
			zap.Error(errors.New("not found gameCode in config")))

		return nil, errors.New("not found gameCode in config")
	}

	u.log.Debug(logPrefix, zap.String("gamecode", gameCode))

	betWin, err := u.rtpRepo.FindBetWinByGameCode(gameCode)
	if err != nil {
		return nil, err
	}

	return &dto.RTPsResponse{
		RTPs: []dto.RTP{
			{
				GameCode: gameCode,
				RTP:      u.caculatorRTP(betWin),
			},
		},
	}, nil
}

func (u *rtpUseCase) FindRTPsAll() (*dto.RTPsResponse, error) {

	logPrefix := "rtpUseCase.FindRTPsAll"

	gameCodes := u.cfg.App.SearchGameCode
	if len(gameCodes) == 0 || gameCodes[0] == "" {
		u.log.Error(logPrefix,
			zap.Error(errors.New("server config app.searchGameCode is empty")))

		return nil, errors.New("server config app.searchGameCode is empty")
	}

	u.log.Debug(logPrefix, zap.Strings("gamecode", gameCodes))

	resp := &dto.RTPsResponse{}

	for _, gameCode := range gameCodes {

		betWin, err := u.rtpRepo.FindBetWinByGameCode(gameCode)
		if err != nil {
			// 找不到key pass
		} else {
			resp.RTPs = append(resp.RTPs, dto.RTP{
				GameCode: gameCode,
				RTP:      u.caculatorRTP(betWin),
			})
		}

	}

	return resp, nil
}

func (u *rtpUseCase) caculatorRTP(betWin entity.BetWin) (rtp float64) {

	if betWin.Bet == 0 {
		return 0
	} else {
		return math.Round((betWin.Win/betWin.Bet*100)*100) / 100
	}

}

func (u *rtpUseCase) ClearRTPsAll() (*dto.ClearResponse, error) {
	logPrefix := "rtpUseCase.ClearRTPsAll"

	gameCodes := u.cfg.App.SearchGameCode
	if len(gameCodes) == 0 || gameCodes[0] == "" {
		u.log.Error(logPrefix,
			zap.Error(errors.New("server config app.searchGameCode is empty")))

		return nil, errors.New("server config app.searchGameCode is empty")
	}
	u.log.Debug(logPrefix, zap.Strings("gamecode", gameCodes))

	resp := &dto.ClearResponse{}

	for _, gameCode := range gameCodes {

		message := "Clear Success"

		err := u.rtpRepo.ClearBetWinValueByGameCode(gameCode)
		if err != nil {
			message = err.Error()
		}
		resp.GameInfos = append(resp.GameInfos, dto.ClearInfo{
			GameCode: gameCode,
			Mesage:   message,
		})

	}
	return resp, nil
}

func (u *rtpUseCase) ClearRTPsByGameCode(gameCode string) (*dto.ClearResponse, error) {

	logPrefix := "rtpUseCase.ClearRTPsByGameCode"

	if _, exists := u.cfgGameList[gameCode]; !exists {
		u.log.Error(logPrefix,
			zap.String("gameCode", gameCode),
			zap.Error(errors.New("not found gameCode in config")))

		return nil, errors.New("not found gameCode in config")
	}

	u.log.Debug(logPrefix, zap.String("gamecode", gameCode))

	resp := &dto.ClearResponse{}

	message := "Clear Success"

	err := u.rtpRepo.ClearBetWinValueByGameCode(gameCode)
	if err != nil {
		message = err.Error()
	}
	resp.GameInfos = append(resp.GameInfos, dto.ClearInfo{
		GameCode: gameCode,
		Mesage:   message,
	})

	return resp, nil
}
