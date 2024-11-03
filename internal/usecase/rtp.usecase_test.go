package usecase

import (
	"math"
	"reflect"
	"testing"

	"sonymimic1/Golang_server/checkRTP/config"
	"sonymimic1/Golang_server/checkRTP/internal/dto"
	"sonymimic1/Golang_server/checkRTP/internal/entity"
	"sonymimic1/Golang_server/checkRTP/pkg/logger"

	"github.com/stretchr/testify/mock"
)

// mockRTPRepository 是 RTPRepository 接口的一個模擬實現
type mockRTPRepository struct {
	mock.Mock
}

// FindBetWinByGameCode is a method that needs to be implemented
func (m *mockRTPRepository) FindBetWinByGameCode(gameCode string) (entity.BetWin, error) {
	args := m.Called(gameCode)
	return args.Get(0).(entity.BetWin), args.Error(1)
}

// FindBetWinByGameCode is a method that needs to be implemented
func (m *mockRTPRepository) ClearBetWinValueByGameCode(gameCode string) error {
	args := m.Called(gameCode)
	return args.Error(1)
}

func Test_rtpUseCase_FindBetWinByGameCode(t *testing.T) {
	logger := logger.NewLogger(&logger.Config{}, nil)
	cfg, err := config.LoadConfig("config")
	if err != nil {
		panic(err)
	}
	// 創建一個 mockRTPRepository 的實例
	mockRepo := &mockRTPRepository{}
	mockRepo.Mock.On("FindBetWinByGameCode", "AT01").Return(entity.BetWin{
		Bet: 10,
		Win: 5,
	}, nil)

	usecase := NewRTPUseCase(mockRepo, cfg, logger) // var rtpUseCase *rtpUseCase
	// useCase := &rtpUseCase{
	// 	rtpRepo: mockRepo,
	// }

	gameCode := "AT01"
	mockRTP := &entity.BetWin{
		Bet: 10,
		Win: 5,
	}

	want := &dto.RTPsResponse{
		RTPs: []dto.RTP{
			{
				GameCode: gameCode,
				RTP:      math.Round((mockRTP.Win/mockRTP.Bet*100)*100) / 100,
			},
		},
	}

	got, err := usecase.FindRTPByGameCode(gameCode)
	if err != nil {
		t.Errorf("rtpUseCase.FindBetWinByGameCode() error = %v", err)
		return
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("rtpUseCase.FindBetWinByGameCode() = %v, want %v", got, want)
	}

}
