package repo

import (
	"context"
	"reflect"
	"testing"

	"sonymimic1/Golang_server/checkRTP/internal/entity"

	"github.com/redis/go-redis/v9"
)

func Test_rtpRepository_FindBetWinByGameCode(t *testing.T) {
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"127.0.0.1:6371",
			"127.0.0.1:6372",
			"127.0.0.1:6373",
			"127.0.0.1:6374",
			"127.0.0.1:6375",
			"127.0.0.1:6376"},
	})
	ctx := context.Background()

	r := &rtpRepository{
		redisClient: redisClient,
		ctx:         ctx,
	}

	gameCode := "AT01"

	// Set up test data in Redis
	redisClient.Set(ctx, gameCode+"-BET", "10.5", 0)
	redisClient.Set(ctx, gameCode+"-WIN", "5.2", 0)

	want := entity.BetWin{
		Bet: 10.5,
		Win: 5.2,
	}

	got, err := r.FindBetWinByGameCode(gameCode)
	if err != nil {
		t.Errorf("rtpRepository.FindBetWinByGameCode() error = %v", err)
		return
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("rtpRepository.FindBetWinByGameCode() = %v, want %v", got, want)
	}
}
