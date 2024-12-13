package app

import (
	"fmt"
	"log"
	"os"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"meta-node-ficam/internal/config"
	"meta-node-ficam/internal/database"
	"meta-node-ficam/internal/handler"
	"meta-node-ficam/internal/repositories"
	"meta-node-ficam/internal/service"
	"meta-node-ficam/internal/route"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/cmd/client"
	"github.com/meta-node-blockchain/meta-node/types"
	c_config "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type App struct {
	Config      *config.AppConfig
	ApiApp      *gin.Engine
	StorageClient *client.Client
	StopChan    chan bool
	EventChan   chan types.EventLogs
	EventHandler     *handler.EventHandler
}

func NewApp(
	configPath string,
	loglevel int,
) (*App, error) {
	var loggerConfig = &logger.LoggerConfig{
		Flag:    loglevel,
		Outputs: []*os.File{os.Stdout},
	}
	logger.SetConfig(loggerConfig)
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("can not load config", err)
	}
	app := &App{}
	app.StorageClient, err = client.NewStorageClient(
		&c_config.ClientConfig{
		  Version_:                config.MetaNodeVersion,
		  PrivateKey_:             config.PrivateKey_,
		  ParentAddress:           config.StorageAddress,
		  ParentConnectionAddress: config.NodeConnectionAddress,
		  DnsLink_:                config.DnsLink(),
		},
		[]common.Address{
			common.HexToAddress(config.FicamAddress),
		},
	)
	if err != nil {
		logger.Error(fmt.Sprintf("error when create storage client %v", err))
		return nil, err
	}
	//
	app.EventChan = app.StorageClient.GetEventLogsChan()
	// create customdomain abi
	reader, err := os.Open(config.FicamABIPath)
	if err != nil {
		logger.Error("Error occured while read resolver abi")
		return nil, err
	}
	defer reader.Close()

	ficamAbi, err := abi.JSON(reader)
	if err != nil {
		logger.Error("Error occured while parse resolver smart contract abi")
		return nil, err
	}
	engine := gin.Default()
	database.StartMySQL(config)
	db := database.GetMySqlConn()
	emailRepo := repositories.NewEmailRepository(db)
	emailService := service.NewEmailService(emailRepo)
	app.Config = config
	app.EventHandler = handler.NewEventHandler(&ficamAbi)
	handler := handler.NewEmailHandler(
		emailService,
	)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	//
	engine.Use(cors.New(corsConfig))

	route.InitialRoutes(
		engine,
		handler,
	)
	app.ApiApp = engine
	return app, nil
}

func (app *App) Run() {
	go func() {
		app.ApiApp.Run(app.Config.API_PORT)
	}()
	for {
		select {
		case <-app.StopChan:
			return
		case eventLogs := <-app.EventChan:
			app.EventHandler.HandleEvent(eventLogs)
		}
	}
}
  