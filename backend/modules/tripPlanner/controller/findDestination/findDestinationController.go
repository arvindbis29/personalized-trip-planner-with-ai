package findDestinationController

import (
	globalFuctions "trip-planner-backend/utilities/globalFunctions"

	"github.com/gin-gonic/gin"
)

func FindDestination(ginCtx *gin.Context) {
	globalFuctions.WirteJsonLogs(ginCtx, "arvind_testing", map[string]interface{}{"arvind": "k"})
	ginCtx.Status(200)
}
