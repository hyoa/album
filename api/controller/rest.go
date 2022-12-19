package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyoa/album/api/internal/media"
)

type CloudconvertHookPayload struct {
	Event string `json:"event"`
	Job   struct {
		ID    string `json:"id"`
		Tasks []struct {
			Result struct {
				Files []struct {
					Filename string `json:"filename"`
					URL      string `json:"url"`
				} `json:"files"`
			} `json:"result"`
		} `json:"tasks"`
	} `json:"job"`
}

type RestController struct {
	mediaManager media.MediaManager
}

func CreateRestController(mediaManager media.MediaManager) RestController {
	return RestController{
		mediaManager: mediaManager,
	}
}

func (rc *RestController) AcknowledgeCloudconvertCall(ctx *gin.Context) {
	var request CloudconvertHookPayload

	errBind := ctx.ShouldBindJSON(&request)

	if errBind != nil {
		ctx.JSON(http.StatusNoContent, gin.H{})
		return
	}

	if request.Event != "job.finished" || len(request.Job.Tasks) == 0 || len(request.Job.Tasks[0].Result.Files) == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{})
		return
	}

	rc.mediaManager.AcknowledgeVideoConversion(request.Job.Tasks[0].Result.Files[0].Filename)

	ctx.JSON(http.StatusAccepted, gin.H{})
}
