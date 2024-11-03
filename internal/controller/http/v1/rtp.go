package v1

import (
	"net/http"

	"sonymimic1/Golang_server/checkRTP/config"
	"sonymimic1/Golang_server/checkRTP/internal/usecase"

	"github.com/gin-gonic/gin"
)

type RTPHandler interface {
	GetAllRTP(ctx *gin.Context)
	GetRTPByGameCode(ctx *gin.Context)
	ClearRTPAll(ctx *gin.Context)
	ClearRTPByGameCode(ctx *gin.Context)
}

type rtpHandler struct {
	rtpUseCase usecase.RTPUseCase
	cfg        config.Config
}

func NewRTPHandler(rtpUseCase usecase.RTPUseCase, cfg config.Config) RTPHandler {
	return &rtpHandler{
		rtpUseCase: rtpUseCase,
		cfg:        cfg,
	}
}

func (h *rtpHandler) GetAllRTP(ctx *gin.Context) {

	rtpsResponse, err := h.rtpUseCase.FindRTPsAll()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, rtpsResponse)
}

func (h *rtpHandler) GetRTPByGameCode(ctx *gin.Context) {

	gameCode := ctx.Query("gamecode")

	rtpsResponse, err := h.rtpUseCase.FindRTPByGameCode(gameCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, rtpsResponse)
}

func (h *rtpHandler) ClearRTPAll(ctx *gin.Context) {

	clearResponse, err := h.rtpUseCase.ClearRTPsAll()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, clearResponse)

}

func (h *rtpHandler) ClearRTPByGameCode(ctx *gin.Context) {

	gameCode := ctx.Query("gamecode")

	clearResponse, err := h.rtpUseCase.ClearRTPsByGameCode(gameCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, clearResponse)

}
