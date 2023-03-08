package debug

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"unicode/utf8"
)

type debug struct {
	debug *Debug
}

type Area struct {
	File  string
	Func  string
	Alias string
}

type Debug struct {
	areas   []*Area
	aread   map[string]*Area
	current int
	label   string
	alias   string
}

var lock = &sync.Mutex{}

var singleDebugInstance *debug

func getInstance() *debug {
	if singleDebugInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleDebugInstance == nil {
			singleDebugInstance = &debug{}
		} else {
		}
	} else {
	}
	return singleDebugInstance
}

func NewDebug() {
	di := getInstance()
	lock.Lock()
	defer lock.Unlock()
	if di.debug == nil {
		di.debug = &Debug{}
		di.debug.aread = make(map[string]*Area)
	} else {
	}
}

func SetArea(areas ...Area) {
	di := getInstance()
	lock.Lock()
	defer lock.Unlock()
	if di.debug == nil {
	} else {
		for i, _ := range areas {
			di.debug.areas = append(di.debug.areas, &areas[i])
			if len(areas[i].File) > 0 {
				di.debug.aread[areas[i].File+"_file"] = &areas[i]
			}
			if len(areas[i].Func) > 0 {
				di.debug.aread[areas[i].Func+"_func"] = &areas[i]
			}
			if len(areas[i].Alias) > 0 {
				di.debug.aread[areas[i].Alias+"_alias"] = &areas[i]
			}
		}
	}
}

func LoadFromFile(fileName string) error {
	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	if bs[0] == 0xEF && bs[1] == 0xBB && bs[2] == 0xBF {
		bs = bs[3:]
	}
	str := string(bs)
	lines_list := strings.Split(str, "\n")
	areas := []Area{}
	for i := 0; i < len(lines_list); i++ {
		ln := lines_list[i]
		line := strings.TrimSpace(ln)
		if len(line) > 0 {
			ch, _ := utf8.DecodeRune([]byte(line))
			if ch == '#' {
			} else {
				strList := strings.Split(line, ":")
				if len(strList) == 2 {
					tag := strings.Trim(strList[0], " ")
					tag = strings.ToLower(tag)
					value := strings.Trim(strList[1], " ")
					switch tag {
					case "alias":
						a := Area{Alias: value}
						areas = append(areas, a)
					case "func":
						a := Area{Func: value}
						areas = append(areas, a)
					case "file":
						a := Area{File: value}
						areas = append(areas, a)
					}
				}
			}
		}
	}
	if len(areas) > 0 {
		SetArea(areas...)
	}
	return nil
}

func Printf(fmts string, args ...interface{}) *Debug {
	di := getInstance()
	lock.Lock()
	defer lock.Unlock()
	if di.debug == nil {
	} else {
		di.debug.current = 0
		di.debug.label = ""
		di.debug.alias = ""
		return di.debug.Printf(fmts, args...)
	}
	return nil
}

func trace(level int) (string, int, string) {
	pc, file, line, ok := runtime.Caller(3 + level)
	if !ok {
		return "?", 0, "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return file, line, "?"
	}

	return file, line, fn.Name()
}

func (d *Debug) Printf(fmts string, args ...interface{}) *Debug {
	if d != nil {
		f, _, n := trace(d.current)
		_, file := filepath.Split(f)
		flag := false
		fnl := strings.Split(n, ".")
		fn := ""
		if len(fnl) > 1 {
			fn = fnl[len(fnl)-1]
		} else {
			fn = n
		}

		_, ok := d.aread[file+"_file"]
		if ok {
			flag = true
		}
		_, ok = d.aread[fn+"_func"]
		if ok {
			flag = true
		}
		_, ok = d.aread[d.alias+"_alias"]
		if ok {
			flag = true
		}
		if flag {
			argss := []interface{}{file, fn}
			argss = append(argss, args...)
			fmt.Printf("%v %v "+fmts, argss...)
		}
	}
	d.alias = ""
	d.current = d.current - 1
	return d
}

func Label(l string) *Debug {
	di := getInstance()
	lock.Lock()
	defer lock.Unlock()
	if di.debug == nil {
	} else {
		di.debug.current = 0
		return di.debug.Label(l)
	}
	return nil
}

func (d *Debug) Label(l string) *Debug {
	d.current = d.current - 1
	d.label = l
	return d
}

func Alias(n string) *Debug {
	di := getInstance()
	lock.Lock()
	defer lock.Unlock()
	if di.debug == nil {
	} else {
		di.debug.current = 0
		return di.debug.Alias(n)
	}
	return nil
}

func (d *Debug) Alias(n string) *Debug {
	d.current = d.current - 1
	d.alias = n
	return d
}
