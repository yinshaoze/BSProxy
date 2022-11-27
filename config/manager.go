package config

import (
	"encoding/json"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/yinshaoze/BSProxy/common/set"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

var (
	Config     configMain
	Lists      map[string]*set.StringSet
	reloadLock sync.Mutex
)

func LoadConfig() {
	configFile, err := os.ReadFile("BSProxy.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Configuration file is not exists. Generating a new one...")
			generateDefaultConfig()
			goto success
		} else {
			log.Panic(color.HiRedString("Unexpected error when loading config: %s", err.Error()))
		}
	}

	err = json.Unmarshal(configFile, &Config)
	if err != nil {
		log.Panic(color.HiRedString("Config format error: %s", err.Error()))
	}

success:
	LoadLists(false)
	log.Println(color.HiYellowString("Successfully loaded config from file."))
}

func generateDefaultConfig() {
	file, err := os.Create("BSProxy.json")
	if err != nil {
		log.Panic("Failed to create configuration file:", err.Error())
	}
	Config = configMain{
		Services: []*ConfigProxyService{
			{
				Name:          "HypixelDefault",
				TargetAddress: "mc.hypixel.net",
				TargetPort:    25565,
				Listen:        25565,
				Flow:          "auto",
				Minecraft: minecraft{
					EnableHostnameRewrite: true,
					IgnoreFMLSuffix:       true,
					OnlineCount: onlineCount{
						Max:            114514,
						Online:         -1,
						EnableMaxLimit: false,
					},
					MotdFavicon:     "{DEFAULT_MOTD}",
					MotdDescription: "§r§c§lProxy for §6§nhypixel.net:25565§r\\n§r§9QQ:3325395619",
				},
			},
		},
		Lists: map[string][]string{
			//"test": {"foo", "bar"},
		},
	}
	newConfig, _ := json.MarshalIndent(Config, "", "    ")
	_, err = file.WriteString(strings.ReplaceAll(string(newConfig), "\n", "\r\n"))
	file.Close()
	if err != nil {
		log.Panic("Failed to save configuration file:", err.Error())
	}
}

func LoadLists(isReload bool) bool {
	reloadLock.Lock()
	defer reloadLock.Unlock()
	if isReload {
		configFile, err := os.ReadFile("BSProxy.json")
		if err != nil {
			if os.IsNotExist(err) {
				log.Println(color.HiRedString("Fail to reload : Configuration file is not exists."))
			} else {
				log.Println(color.HiRedString("Unexpected error when reloading config: %s", err.Error()))
			}
			return false
		}

		err = json.Unmarshal(configFile, &Config)
		if err != nil {
			log.Println(color.HiRedString("Fail to reload : Config format error: %s", err.Error()))
			return false
		}
	}
	// log.Println("Lists:", Config.Lists)
	if l := len(Config.Lists); l == 0 { // if nothing in Lists
		Lists = map[string]*set.StringSet{} // empty map
	} else {
		Lists = make(map[string]*set.StringSet, l) // map size init
		for k, v := range Config.Lists {
			// log.Println("List: Loading", k, "value:", v)
			list := set.NewStringSetFromSlice(v)
			Lists[k] = &list
		}
	}
	Config.Lists = nil // free memory
	debug.FreeOSMemory()
	return true
}

func MonitorConfig(watcher *fsnotify.Watcher) error {
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op.Has(fsnotify.Write) { // config reload
					log.Println(color.HiMagentaString("Config Reload : File change detected. Reloading..."))
					if LoadLists(true) { // reload success
						log.Println(color.HiMagentaString("Config Reload : Successfully reloaded Lists."))
					} else {
						log.Println(color.HiMagentaString("Config Reload : Failed to reload Lists."))
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println(color.HiRedString("Config Reload Error : ", err))
			}
		}
	}()
	return watcher.Add("BSProxy.json")
}
