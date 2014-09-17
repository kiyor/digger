/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 09-16-2014

* Last Modified : Wed 17 Sep 2014 12:28:32 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"fmt"
	"github.com/miekg/unbound"
	"github.com/wsxiaoys/terminal/color"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type reslov struct {
	note string
	file string
	init bool
	u    *unbound.Unbound
	res  string
}

var domain string

func init() {
	if len(os.Args) == 1 {
		fmt.Println("use digger domain.com")
		os.Exit(1)
	}
	domain = os.Args[1]
}

func main() {
	reslovs := getReslovs("/usr/local/etc/reslov/")
	ch := make(chan reslov)
	for _, r := range reslovs {
		go func(r reslov) {
			r.u = unbound.New()
			defer r.u.Destroy()
			r.u.ResolvConf(r.file)
			r.digger(domain)
			ch <- r
		}(r)
	}
	t := time.Tick(time.Second * 20)
	for i := 0; i < len(reslovs); i++ {
		select {
		case r := <-ch:
			fmt.Println(r.res)
		case <-t:
		}
	}
}

func (r *reslov) digger(domain string) {
	if r.init {
		r.res += color.Sprintf("@{r}DIG@{|}   @{g}%s@{|} in %s\n", domain, r.note)
		r.init = false
	}
	cname, err := r.u.LookupCNAME(domain)
	if err != nil {
		log.Fatalf("error %s\n", err.Error())
	}
	if len(cname) != 0 {
		r.res += color.Sprintf("@{y}CNAME@{|} %s\n", cname)
		r.digger(cname)
		return
	}
	a, err := r.u.LookupHost(domain)
	if err != nil {
		log.Fatalf("error %s\n", err.Error())
	}
	for _, a1 := range a {
		r.res += color.Sprintf("@{c}A@{|}     %s\n", a1)
	}
	ns, err := r.u.LookupNS(domain)
	if err != nil {
		log.Fatalf("error %s\n", err.Error())
	}
	for _, n := range ns {
		r.res += color.Sprintf("@{g}NS@{|}    %s\n", n)
	}
}

func getReslovs(dir string) (rs []reslov) {
	filepath.Walk(dir, func(path string, _ os.FileInfo, _ error) error {
		defer func() {
			if r := recover(); r != nil {
				log.Fatalln("recover in ", path)
			}
		}()
		info, err := os.Stat(path)
		if info.IsDir() || err != nil {
			return nil
		}
		var r reslov
		r.file = path
		r.note = strings.Split(info.Name(), ".")[0]
		r.init = true
		rs = append(rs, r)
		return nil
	})
	return
}
