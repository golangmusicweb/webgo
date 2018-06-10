package restful

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config struct {
	App        map[string]string
	Session    map[string]string
	Datasource map[string](map[string]string)
	Log        map[string]string
	Static     map[string]string
	Redis      map[string]string
	All        map[string]string
}

func (cfg *Config) LoadConfig() {
	cfg.App = make(map[string]string)
	cfg.Session = make(map[string]string)
	cfg.Datasource = make(map[string](map[string]string))
	cfg.Log = make(map[string]string)
	cfg.Static = make(map[string]string)
	cfg.Redis = make(map[string]string)
	cfg.All = make(map[string]string)

	fl, err := os.Open("config/settings.conf")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		panic(err)
		return
	}

	defer fl.Close()

	bf := bufio.NewReader(fl)

	for {
		line, err := bf.ReadString('\n')
		if err == io.EOF {
			break
		} else {
			tmp := strings.TrimLeft(line, " ")
			tmp = strings.TrimRight(tmp, " ")
			tmp = strings.TrimRight(tmp, "\n")
			if len(tmp) == 0 || strings.Index(tmp, "#") == 0 {
				continue
			}

			value := strings.Split(tmp, "=")
			//fmt.Println(tmp)
			if len(value) == 2 {
				cfg.All[strings.Trim(value[0], " ")] = strings.Trim(value[1], " ")
			}
		}
	}

	for k, v := range cfg.All {
		if strings.Index(k, "webgo.app.") == 0 {
			tmp := strings.TrimPrefix(k, "webgo.app.")
			cfg.App[tmp] = v
		}
		if strings.Index(k, "webgo.session.") == 0 {
			tmp := strings.TrimPrefix(k, "webgo.session.")
			cfg.Session[tmp] = v
		}
		if strings.Index(k, "webgo.datasource.") == 0 {
			var sd = strings.Split(k, ".")
			_, exists := cfg.Datasource[sd[2]]
			if exists == false {
				cfg.Datasource[sd[2]] = make(map[string]string)
				cfg.Datasource[sd[2]][sd[3]] = v
			} else {
				cfg.Datasource[sd[2]][sd[3]] = v
			}
		}
		if strings.Index(k, "webgo.log.") == 0 {
			tmp := strings.TrimPrefix(k, "webgo.log.")
			cfg.Log[tmp] = v
		}
		if strings.Index(k, "webgo.static.") == 0 {
			tmp := strings.TrimPrefix(k, "webgo.static.")
			cfg.Static[tmp] = v
		}
		if strings.Index(k, "webgo.redis.") == 0 {
			tmp := strings.TrimPrefix(k, "webgo.redis.")
			cfg.Redis[tmp] = v
		}
	}
}
