package main

import (
	"fmt"
	"github.com/cloudcapsule/cdp/pkg/plugin"
	"github.com/cloudcapsule/cdp/pkg/task"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

type param struct {
	name      string
	shorthand string
	value     interface{}
	usage     string
	required  bool
}

var (
	Version    string
	Build      string
	rootParams = []param{
		{name: "verbose", shorthand: "", value: false, usage: "enable verbose logs"},
	}
	startParams = []param{
		{name: "addr", shorthand: "a", value: "0.0.0.0:50052", usage: "bind grpc server address"},
	}
	rootCmd = &cobra.Command{
		Use:   "cdp",
		Short: "cdp - cloud capsule data plugin",
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print cdp version and build sha",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("🐾 versionCmd: %s build: %s \n", Version, Build)
		},
	}
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start cdp server",
		Run: func(cmd *cobra.Command, args []string) {
			sigCh := make(chan os.Signal)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
			dt := []task.DataTask{&task.PGTask{}}
			svc := plugin.NewDataPluginService(task.NewExecutor(dt))
			svc.Serve()
			<-sigCh
		},
	}
)

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CDP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	setupLogging()
}

func init() {
	cobra.OnInitialize(initConfig)
	setParams(rootParams, rootCmd)
	setParams(startParams, startCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
}

func setParams(params []param, command *cobra.Command) {
	for _, param := range params {
		switch v := param.value.(type) {
		case int:
			command.PersistentFlags().IntP(param.name, param.shorthand, v, param.usage)
		case string:
			command.PersistentFlags().StringP(param.name, param.shorthand, v, param.usage)
		case bool:
			command.PersistentFlags().BoolP(param.name, param.shorthand, v, param.usage)
		}
		if err := viper.BindPFlag(param.name, command.PersistentFlags().Lookup(param.name)); err != nil {
			panic(err)
		}
	}
}

func setupLogging() {

	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := strings.TrimSuffix(filepath.Base(frame.File), filepath.Ext(frame.File))
			line := strconv.Itoa(frame.Line)
			return "", fmt.Sprintf("%s:%s", fileName, line)
		},
	})
	// Logs are always goes to STDOUT
	log.SetOutput(os.Stdout)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
