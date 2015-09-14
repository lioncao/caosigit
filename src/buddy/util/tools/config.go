package tools

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type configMgr struct {
	MainData    map[string]map[string]string
	ServiceData []map[string]string
	curSection  string
}

func (this *configMgr) Load(path string) error {
	this.MainData = make(map[string]map[string]string)
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return err
	}
	r := bufio.NewReader(f)
	for {
		s, err1 := r.ReadString('\n')
		if err1 != nil {
			this.PraseString(s)
			break
		}
		this.PraseString(s)
	}
	return nil
}

func (this *configMgr) GetSection(s string) map[string]string {
	se, ok := this.MainData[s]
	if ok {
		return se
	} else {
		return nil
	}
}

func (this *configMgr) Get(section, key string) (string, bool) {
	s, ok := this.MainData[section]
	if ok {
		v, ok := s[key]
		return v, ok

	}
	return "", false
}

func (this *configMgr) GetInt(section, key string) (int64, bool) {
	s, ok := this.Get(section, key)
	if ok {
		i, e := strconv.ParseInt(s, 0, 64)
		if e != nil {
			return 0, false
		}
		return i, true

	}
	return 0, false
}

func (this *configMgr) PraseString(content string) {
	realStr := strings.TrimSpace(content)
	if realStr == "" || realStr[0] == '#' {
		return
	} else {
		if strings.ContainsAny(realStr, "[&]") {
			section := realStr[1 : len(realStr)-1]
			if strings.EqualFold(section, "Service") {
				sd := make(map[string]string)
				this.ServiceData = append(this.ServiceData, sd)
			} else {
				_, ok := this.MainData[section]
				if !ok {
					this.MainData[section] = make(map[string]string)
				}
			}
			this.curSection = section
		} else {
			if strings.Contains(realStr, "=") {
				kv := strings.Split(realStr, "=")
				if len(kv) == 2 {
					if strings.EqualFold(this.curSection, "service") {
						sm := this.ServiceData[len(this.ServiceData)-1]
						//fmt.Printf("serveric push s:%s key:%s v:%s\n", this.curSection, kv[0], kv[1])
						sm[kv[0]] = kv[1]
					} else {
						datamap := this.MainData[this.curSection]
						//fmt.Printf("push s:%s key:%s v:%s\n", this.curSection, kv[0], kv[1])
						datamap[kv[0]] = kv[1]
					}
				}

			}
		}

	}
}

var config_intanse *configMgr

func GetConfigMgr() *configMgr {
	if config_intanse == nil {
		config_intanse = new(configMgr)
	}
	return config_intanse
}
