package webserver

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/AryalKTM/monitor/monitor/system"
	"github.com/AryalKTM/monitor/core/utilities"
	"github.com/AryalKTM/monitor/core/models/response"
	"github.com/AryalKTM/monitor/monitor/models/data"
	"github.com/AryalKTM/monitor/monitor/models/config"
	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/rakyll/statik/fs"
)

const (
	ROOM_SERVICES_LISTENERS = "serviceListners"
)

const (
	EVENT_SERVICE_CHANGE = "serviceChange"
)

var serverConfig *config.ServerConfig
var lastSystemInfo system.SystemInfo
var lastTsSystemInfo time.Time

type WebServer struct{
	Address	string
	Port int
	SSL bool
	Router *gin.Engine
	InputChannel chan response.Response
	ServiceMap *map[string]data.ServiceData
	ServerSocket *socketio.Server
}

func NewWebServer(serverConfPath string) *WebServer {
	sConf, err := config.ServerConfigFromFile(serverConfPath)

	if err != nil {
		log.Fatal(err.Error())
	}
	serverConfig = sConf

	ServerSocket := newServerSocket()
	go func() {
		err := ServerSocket.Serve()

		if err != nil {
			log.Fatal(err)
		}
	}()

	//defer serverSocket.Close()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = serverConfig.AllowOrigins
	configCors.AllowCredentials = true
	router.Use(cors.New(configCors))
	// router.Use(GinMiddleware("localhost:8080"))
	// router.Use(cors.Default())
	// router.Use(CORSMiddleware())
	// router.Use(gin.Logger())
	router.Use(gin.Recovery())
	inputChannel := make(chan response.Response)
	serviceMap := make(map[string]data.ServiceData, 0)

	ws := WebServer{
		Address: serverConfig.Address,
		Port: serverConfig.Port,
		SSL: serverConfig.SSL,
		Router: router,
		ServiceMap: &serviceMap,
		InputChannel: inputChannel,
		ServerSocket: ServerSocket,
	}
	lastSystemInfo, err = system.GetSystemInfo()
	if err != nil {
		log.Fatal(err)
		
	}
	//Register REST
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	staticGroup := router.Group("/dashboard")
	staticGroup.StaticFS("/", statikFS)

	apiGroup := router.Group("/api")
	apiGroup.GET("/services", func(context *gin.Context) {getServicesList(context, &ws)})
	apiGroup.GET("/system", func(context *gin.Context) {getSystemInformation(context)})

	router.GET("/socket.io/*any", gin.WrapH(ServerSocket))
	router.POST("/socket.io/*any", gin.WrapH(ServerSocket))

	router.GET("/", func(context *gin.Context){context.Redirect(http.StatusMovedPermanently, "dashboard")})

	lastSystemInfo, err = system.GetSystemInfo()
	if err != nil {
		log.Fatal(err)
	}
	lastTsSystemInfo = time.Now()
	return &ws
	}

func (ws *WebServer) Start() {
	go func () {
		log.Printf("[%s] %s", utilities.CreateColorString("Web Panel", color.FgHiCyan), "Availaible on Port: "+ strconv.Itoa(ws.Port))
		err := ws.Router.Run(ws.Address + ":" + strconv.Itoa(ws.Port))
		log.Fatal(err.Error())
	} ()

	go func () {
		for {
			ws.listenInputChannel()
		}
	} ()
}

func (ws *WebServer) listenInputChannel() {
	resp := <-ws.InputChannel;
		if serviceData, exists := (*ws.ServiceMap) [resp.ServiceName]; exists {
			protocolExists := false
			sendToWs := false

			for i := 0; i < len(serviceData.Protocols); i++ {
				protocolData := &serviceData.Protocols[i]
				if protocolData.Protocol.Type == resp.Protocol.Type && protocolData.Protocol.Server == resp.Protocol.Server && protocolData.Protocol.Port == resp.Protocol.Port {
					if protocolData.Err != nil && resp.Error != nil {
						if protocolData.Err.Error() != resp.Error.Error() {
							protocolData.Err = resp.Error
							sendToWs = true
						}
					} else {
						if protocolData.Err != resp.Error {
							protocolData.Err = resp.Error
							sendToWs = true
						}
					}
					protocolExists = true
				}
			}
			if !protocolExists {
				serviceData.Protocols = append(serviceData.Protocols, data.ProtocolData {
					Protocol: resp.Protocol,
					Err: resp.Error,
				})
			}
			if sendToWs {
				ws.ServerSocket.BroadcastToRoom("/", ROOM_SERVICES_LISTENERS, EVENT_SERVICE_CHANGE, serviceData)
			}
			(*ws.ServiceMap)[resp.ServiceName] = serviceData
		} else {
			//Service Does Not Exist,
			firstProtocol := data.ProtocolData {
				Protocol: resp.Protocol,
				Err: resp.Error,
			}
			protocolsList := make([]data.ProtocolData, 1)
			protocolsList[0] = firstProtocol
			auxData := data.ServiceData {
				Name: resp.ServiceName,
				Protocols: protocolsList,
			}
			(*ws.ServiceMap)[resp.ServiceName] = auxData
		}
	}

func newServerSocket() *socketio.Server {
	serverSocket := socketio.NewServer(nil)

	serverSocket.OnConnect("/", func (s socketio.Conn) error{
		s.Join(ROOM_SERVICES_LISTENERS)
		s.SetContext("")
		return nil
	})

	serverSocket.OnEvent("/", "Notice", func (s socketio.Conn, msg string)  {
		fmt.Println("Notice: ", msg)
		s.Emit("Reply", "Have" +msg)
	})
	
	serverSocket.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("Meet Error: ", e)
	})

	serverSocket.OnDisconnect("/", func (s socketio.Conn, reason string)  {
		s.Close()
	})
	return serverSocket
}

func getServicesList(context *gin.Context, ws *WebServer) {
	context.JSON(http.StatusOK, ws.ServiceMap)
}

func getSystemInformation(context *gin.Context) {
	var err error
	timeNow := time.Now()

	if timeNow.After(lastTsSystemInfo.Add(time.Millisecond * time.Duration(serverConfig.IntervalSystemInfoMs))) {
		lastSystemInfo, err = system.GetSystemInfo()
		lastTsSystemInfo = timeNow
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, err)
	} else {
		context.JSON(http.StatusOK, lastSystemInfo)
	}
}

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Request.Header.Del("Origin")
		c.Next()
	}
}