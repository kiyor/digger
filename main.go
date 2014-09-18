/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 09-16-2014

* Last Modified : Thu 18 Sep 2014 07:54:05 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/unbound"
	"github.com/wsxiaoys/terminal/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type reslov struct {
	note   string
	file   string
	init   bool
	u      *unbound.Unbound
	res    string
	tstart time.Time
	tend   time.Time
}

var (
	domain   string
	reDomain = regexp.MustCompile(`([a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,6})`)
)

func init() {
	if len(os.Args) == 1 {
		fmt.Println("use digger domain.com")
		os.Exit(1)
	}
	domain = os.Args[1]
	if reDomain.MatchString(domain) {
		val := reDomain.FindStringSubmatch(domain)
		domain = val[0]
	}
}

func main() {
	reslovs := getReslovs("/usr/local/etc/reslov/")
	ch := make(chan reslov)
	for _, r := range reslovs {
		go func(r reslov) {
			r.tstart = time.Now()
			r.u = unbound.New()
			defer r.u.Destroy()
			r.u.ResolvConf(r.file)
			r.digger(domain)
			r.tend = time.Now()
			ch <- r
		}(r)
	}
	t := time.Tick(time.Second * 100)
	for i := 0; i < len(reslovs); i++ {
		select {
		case r := <-ch:
			color.Printf("@{r}DIG@{|}   @{g}%s@{|} in @{g}%s@{|} using @{y}%v@{|}\n", domain, strings.ToUpper(r.note), r.tend.Sub(r.tstart))
			fmt.Println(r.res)
		case <-t:
		}
	}
}

func (r *reslov) digger(domain string) {
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
		go func(a1 string) {
			loc := ip2loc(a1)
			r.res += color.Sprintf("@{c}A@{|}     %-18s %s %s\n", a1, loc.CountryCode, loc.City)
			wg.Done()
		}(a1)
	}
	wg.Wait()
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

type ipLoc struct {
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	RegionName  string `json:"region_name"`
	City        string `json:"city"`
}

func ip2loc(ip string) (loc ipLoc) {
	res, err := http.Get(fmt.Sprintf("http://freegeoip.net/json/%s", ip))
	if err != nil {
		log.Fatalf("error %s\n", err.Error())
	}
	b, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(b, &loc)
	return
}
