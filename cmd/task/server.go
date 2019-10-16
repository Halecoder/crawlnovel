package task

import (
	"crawlnovel/pkg/api/dto"
	"crawlnovel/pkg/api/service"
	"crawlnovel/pkg/common/cache"
	"crawlnovel/pkg/common/db"
	_ "crawlnovel/pkg/crawler/novels/aoyuge"
	_ "crawlnovel/pkg/crawler/novels/dingdian"
	_ "crawlnovel/pkg/crawler/novels/fanfan"
	"fmt"
	"github.com/robfig/cron"
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
	types    int
	loglevel uint8
	//StartCmd : set up restful api server
	StartCmd = &cobra.Command{
		Use:     "task",
		Short:   "task API server",
		Example: "download task",
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
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "./config/in-local.yaml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().IntVar(&types, "type", 1, "type")
}

func usage() {
	usageStr := `       _                       _ 
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
	db.Setup()
	cache.SetUp()
}

func run() error {
	switch types {
	case 1:
		// 启动服务时，先检查task相关状态
		go service.CheckTask()
		// 管理task状态(如task运行完成之后需要将状态标为已完成)
		go service.ManageTaskStatus()
		select {}
		break
	case 2:
		i := 0
		c := cron.New()
		spec := "*/60 * * * * ?"
		c.AddFunc(spec, func() {
			i++
			log.Println("cron running:", i)
		})
		c.Start()
		select {}
		break
	case 3:
		var taskService = service.TaskService{}
		gDto := dto.GeneralGetDto{}
		gDto.Id = 3
		taskService.ExecTask(gDto)
		break
	case 4:
		var taskService = service.TaskService{}
		gDto := dto.GeneralGetDto{}
		gDto.Id = 2
		taskService.ExecTask(gDto)
		break
	default:
		break
	}
	return nil
}
