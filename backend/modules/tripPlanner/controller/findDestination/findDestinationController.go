package findDestinationController

import (
	globalFuctions "trip-planner-backend/utilities/globalFunctions"

	"github.com/gin-gonic/gin"
)

func FindDestination(ginCtx *gin.Context) {
	globalFuctions.WriteJsonLogs(ginCtx, "arvind_testing", map[string]any{"arvind": "k"})
	ginCtx.Status(200)
}
