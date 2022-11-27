package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/yinshaoze/BSProxy/config"
	"github.com/yinshaoze/BSProxy/console"
	"github.com/yinshaoze/BSProxy/service"
	"github.com/yinshaoze/BSProxy/version"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

func main() {
	log.SetOutput(color.Output)
	console.SetTitle(fmt.Sprintf("BSProxy %v | Running...", version.Version))
	console.Println(color.HiRedString(`｀半　　　半　　　半　　　　　　　　山
　　半　　半　　　半　　　山　　　　山　　　　山
　　半　　半　　半　　　　山　　　　山　　　　山
　　　　　半　　　　　　　山　　　　山　　　　山
　半半半半半半半半半　　　山　　　　山　　　　山
　　　　　半　　　　　　　山　　　　山　　　　山
　　　　　半　　　　　　　山　　　　山　　　　山
半半半半半半半半半半半　　山　　　　山　　　　山
　　　　　半　　　　　　　山　　　　山　　　　山
　　　　　半　　　　　　　山　　　　山　　　　山
　　　　　半　　　　　　　山山山山山山山山山山山`))
	color.HiGreen("Welcome to BSProxy %s!\n", version.Version)
	color.HiBlack("Build Information: %s, %s/%s\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
	go version.CheckUpdate()

	config.LoadConfig()

	for _, s := range config.Config.Services {
		go service.StartNewService(s)
	}

	// hot reload
	// use inotify on Linux
	// use Win32 ReadDirectoryChangesW on Windows
	{
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Panic(err)
		}
		defer watcher.Close()
		err = config.MonitorConfig(watcher)
		if err != nil {
			log.Panic("Config Reload Error : ", err)
		}
	}

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
		// stop the program
		// sometimes after the program exits on Windows, the ports are still occupied and "listening".
		// so manually closes these listeners when the program exits.
		for _, listener := range service.ListenerArray {
			if listener != nil { // avoid null pointers
				listener.Close()
			}
		}
	}
}
