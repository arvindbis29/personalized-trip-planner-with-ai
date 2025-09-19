package processItineraryController

import (
	imagemedia "trip-planner-backend/utilities/imageMedia"

	"github.com/gin-gonic/gin"
)

func ProcessItinerary(ginCtx *gin.Context) {
	image := imagemedia.FetchMedia("taj mahal")

	ginCtx.JSON(200, map[string]any{
		"url": image,
	})
}
