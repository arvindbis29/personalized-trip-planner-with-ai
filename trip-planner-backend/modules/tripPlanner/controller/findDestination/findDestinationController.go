package findDestinationController

import (
	"trip-planner/utilities/globalFunctions"

	"github.com/gin-gonic/gin"
)

func FindDestination(ginCtx *gin.Context) {
	globalFunctions.WriteJsonLogs(ginCtx, "arvind_testing", map[string]any{"arvind": "k"})
	ginCtx.Status(200)
}
