package repo

import (
	"context"
	"errors"

	"sonymimic1/Golang_server/checkRTP/internal/entity"
	"sonymimic1/Golang_server/checkRTP/pkg/datasource/redis/redlock"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RTPRepository interface {
	FindBetWinByGameCode(gameCode string) (entity.BetWin, error)
	ClearBetWinValueByGameCode(gameCode string) error
}

type rtpRepository struct {
	redisClient *redis.ClusterClient
	redlock     *redlock.RedLock
	ctx         context.Context
	log         *zap.Logger
}

func NewRTPRepository(ctx context.Context, redisClient *redis.ClusterClient, log *zap.Logger) RTPRepository {
	return &rtpRepository{
		redisClient: redisClient,
		redlock:     redlock.NewRediLock(redisClient),
		ctx:         ctx,
		log:         log,
	}
}

func (r *rtpRepository) FindBetWinByGameCode(gameCode string) (entity.BetWin, error) {

	logPrefix := "rtpRepository.FindBetWinByGameCode"

	GameCode_BET := gameCode + "-BET"
	GameCode_WIN := gameCode + "-WIN"
	betWin := entity.BetWin{}

	exists, err := r.redisClient.Exists(r.ctx, GameCode_BET).Result()
	if err != nil {

		r.log.Error(logPrefix,
			zap.String("Key", GameCode_BET),
			zap.Error(err))

		return betWin, err
	}
	if exists > 0 {
		betWin.Bet, _ = r.redisClient.Get(r.ctx, GameCode_BET).Float64()
	} else {
		r.log.Warn(logPrefix,
			zap.String("Key", GameCode_BET),
			zap.Error(errors.New(gameCode+"-BET not found")))

		return betWin, errors.New(gameCode + "-BET not found")
	}

	exists, err = r.redisClient.Exists(r.ctx, GameCode_WIN).Result()
	if err != nil {

		r.log.Error(logPrefix,
			zap.String("Key", GameCode_WIN),
			zap.Error(err))

		return betWin, err
	}
	if exists > 0 {
		betWin.Win, _ = r.redisClient.Get(r.ctx, GameCode_WIN).Float64()
	} else {
		r.log.Warn(logPrefix,
			zap.String("Key", GameCode_WIN),
			zap.Error(errors.New(gameCode+"-WIN not found")))
		return betWin, errors.New(gameCode + "-WIN not found")
	}

	return betWin, nil
}

func (r *rtpRepository) ClearBetWinValueByGameCode(gameCode string) error {

	logPrefix := "rtpRepository.ClearBetWinValueByGameCode"

	GameCode_BET := gameCode + "-BET"
	GameCode_WIN := gameCode + "-WIN"

	//BET
	exists, err := r.redisClient.Exists(r.ctx, GameCode_BET).Result()
	if err != nil {

		r.log.Error(logPrefix,
			zap.String("IsExists Key", GameCode_BET),
			zap.Error(err))

		return err
	}
	if exists > 0 {
		err := r.redisClient.Set(r.ctx, GameCode_BET, 0, 0).Err()
		if err != nil {

			r.log.Error(logPrefix,
				zap.String("Set Key", GameCode_BET),
				zap.Error(err))

			return err
		}
	} else {

		r.log.Warn(logPrefix,
			zap.String("Key", GameCode_BET),
			zap.Error(errors.New("not found")))

		return errors.New(GameCode_BET + " not found")
	}

	// WIN
	exists, err = r.redisClient.Exists(r.ctx, GameCode_WIN).Result()
	if err != nil {

		r.log.Error(logPrefix,
			zap.String("IsExists Key", GameCode_WIN),
			zap.Error(err))

		return err
	}
	if exists > 0 {
		err = r.redisClient.Set(r.ctx, GameCode_WIN, 0, 0).Err()
		if err != nil {

			r.log.Error(logPrefix,
				zap.String("Set Key", GameCode_WIN),
				zap.Error(err))

			return err
		}
	} else {

		r.log.Warn(logPrefix,
			zap.String("Key", GameCode_WIN),
			zap.Error(errors.New("not found")))

		return errors.New(GameCode_WIN + " not found")
	}

	return nil
}
