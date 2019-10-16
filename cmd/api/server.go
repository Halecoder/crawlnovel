package api

import (
	"crawlnovel/pkg/api/router"
	"crawlnovel/pkg/common/cache"
	"crawlnovel/pkg/common/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

var (
	config   string
	port     string
	loglevel uint8
	//StartCmd : set up restful api server
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "Start download API server",
		Example: "download server config/in-local.yaml",
		PreRun: func(cmd *cobra.Command, args []string) {
			usage()
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/in-local.yaml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "Tcp port server listening on")
	StartCmd.PersistentFlags().Uint8VarP(&loglevel, "loglevel", "l", 0, "Log level")
}

func usage() {
	usageStr := `_                       _ 
                       | |                     | |
  ____  ____ ____ _ _ _| |____   ___ _   _ ____| |
 / ___)/ ___) _  | | | | |  _ \ / _ \ | | / _  ) |
( (___| |  ( ( | | | | | | | | | |_| \ V ( (/ /| |
 \____)_|   \_||_|\____|_|_| |_|\___/ \_/ \____)_|
`
	fmt.Printf("%s\n", usageStr)
}

func setup() {
	//1.Set up log level
	zerolog.SetGlobalLevel(zerolog.Level(loglevel))
	//2.Set up configuration
	viper.SetConfigFile(config)
	content, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}
	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}
	//3.Set up run mode
	mode := viper.GetString("mode")
	gin.SetMode(mode)
	//4.Set up database connection
	db.Setup()
	//5.Set up cache
	cache.SetUp()
	//6.set up mongo
	db.SetUpMongo()
	//7.auto migrate
	db.AutoMigrate()
}

func run() error {
	engine := gin.Default()
	router.Init(engine)
	return engine.Run(":" + port)
}
