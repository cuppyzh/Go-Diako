package diako

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start() {

}

// Used to init Gin engine along with Diako router
func InitRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/api/diako/message", ApiDiakoMessageHandler)
	return router
}

// Used to add Diako router to existing Gin engine
func SetupRouter(router *gin.Engine) {
	router.POST("/api/diako/message", ApiDiakoMessageHandler)
}

func ApiDiakoMessageHandler(context *gin.Context) {
	var request MessageRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Println("Test: " + request.Sender)
	context.String(http.StatusOK, "")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
