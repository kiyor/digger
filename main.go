/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 09-16-2014

* Last Modified : Tue 26 Jul 2016 11:58:35 AM PDT

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"flag"
	"fmt"
	"github.com/miekg/unbound"
	"github.com/wsxiaoys/terminal/color"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

type resolv struct {
	note string
	file string
	init bool
	u    *unbound.Unbound
	res  string
	dur  time.Duration
}

var (
	domain   string
	reDomain = regexp.MustCompile(`([a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,6})`)

	dir *string = flag.String("d", "/usr/local/etc/resolv", "use custom directory of resolv")
)

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(flag.Args()) == 0 {
		fmt.Println(os.Args[0], "[options] domain.com")
		os.Exit(1)
	}
	domain = flag.Arg(0)
	if reDomain.MatchString(domain) {
		val := reDomain.FindStringSubmatch(domain)
		domain = val[0]
	}
	checkres = make(map[string]check)
}

func main() {
	resolvs := getResolvs(*dir)
	ch := make(chan resolv)
	for _, r := range resolvs {
		go func(r resolv) {
			tstart := time.Now()
			r.u = unbound.New()
			defer r.u.Destroy()
			r.u.ResolvConf(r.file)
			r.digger(domain)
			r.dur = time.Now().Sub(tstart)
			ch <- r
		}(r)
	}
	t := time.Tick(time.Second * 100)
	for i := 0; i < len(resolvs); i++ {
		select {
		case r := <-ch:
			color.Printf("@{r}DIG@{|}   @{g}%s@{|} in @{g}%s@{|} using @{y}%v@{|}\n", domain, strings.ToUpper(r.note), r.dur)
			fmt.Println(r.res)
		case <-t:
		}
	}
}

func (r *resolv) digger(domain string) {
	cname, err := r.u.LookupCNAME(domain)
	if err != nil {
		log.Fatalf("error %s\n", err.Error())
	}
	if len(cname) != 0 {
		r.res += color.Sprintf("@{y}CNAME@{|} @{m}%s@{|}\n", cname)
		r.digger(cname)
		return
	}
	a, err := r.u.LookupHost(domain)
	if err != nil {
		log.Fatalf("error %s\n", err.Error())
	}
	r.res += color.Sprintf("@{y}[%-3d]@{|} ips\n", len(a))
	var wg sync.WaitGroup
	for _, a1 := range a {
		wg.Add(1)
		go func(a1 string, r *resolv) {
			r.res += color.Sprintf("@{c}A@{|}     %-18s\n", a1)
			wg.Done()
		}(a1, r)
	}
	wg.Wait()
}

func getResolvs(dir string) (rs []resolv) {
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
		var r resolv
		r.file = path
		r.note = strings.Split(info.Name(), ".")[0]
		r.init = true
		rs = append(rs, r)
		return nil
	})
	return
}
