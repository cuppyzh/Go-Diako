package diako

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gookit/event"
)

func Start() {

}

// Used to init Gin engine along with Diako router
func InitRouter() *gin.Engine {
	router := gin.Default()
	setAuthenticationApi(&router.RouterGroup)
	return router
}

// Used to add Diako router to existing Gin engine
func SetupRouter(router *gin.Engine) {
	setAuthenticationApi(&router.RouterGroup)
}

func setAuthenticationApi(router *gin.RouterGroup) {

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		os.Getenv("DIAKO_AUTH_USERNAME"): os.Getenv("DIAKO_AUTH_PASSWORD"),
	}))

	authorized.POST("/api/diako/message", apiDiakoMessageHandler)
}

func apiDiakoMessageHandler(context *gin.Context) {
	var request MessageRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	eventData := &event.BasicEvent{}
	eventData.SetName("diako.message.recieved")
	eventData.SetData(event.M{
		"Sender":  request.Sender,
		"Message": request.Message,
	})

	event.FireEvent(eventData)

	context.String(http.StatusOK, "")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
