package diako

import (
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"
	"github.com/gookit/event"
)

var cachedMessage []MessageRequest
var thresholdMessage float64 = 3

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

	log.Println(request)

	if !shallSendTheMessage(request) {
		log.Println("Message will not be forwarded")
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

func shallSendTheMessage(messageRequest MessageRequest) bool {
	clearCachedMessage()
	timestamp := time.Now().UTC()
	var existingMessage []MessageRequest

	linq.From(cachedMessage).WhereT(func(mr MessageRequest) bool {
		return mr.Message == messageRequest.Message && mr.Sender == messageRequest.Sender && math.Abs(messageRequest.GetTime().Sub(timestamp).Seconds()) <= thresholdMessage
	}).ToSlice(&existingMessage)

	if len(existingMessage) > 0 {
		return false
	}

	cachedMessage = append(cachedMessage, messageRequest)
	return true
}

func clearCachedMessage() {
	var copyCachedMessage []MessageRequest
	timestamp := time.Now().UTC()

	for _, v := range cachedMessage {
		if math.Abs(v.GetTime().Sub(timestamp).Seconds()) <= thresholdMessage {
			copyCachedMessage = append(copyCachedMessage, v)
		}
	}

	cachedMessage = copyCachedMessage
}
