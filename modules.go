/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : modules.go

* Purpose :

* Creation Date : 09-22-2014

* Last Modified : Tue 23 Sep 2014 01:11:19 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type check struct {
	code   int
	header []string
	// 	doing  *bool
}

var (
	checkres map[string]check
)

type ipLoc struct {
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	RegionName  string `json:"region_name"`
	City        string `json:"city"`
}

func (loc *ipLoc) ip2loc(ip string) {
	res, err := http.Get(fmt.Sprintf("http://freegeoip.net/json/%s", ip))
	if err != nil {
		log.Fatalf("error %s\n", err.Error())
	}
	b, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(b, &loc)
}

func (r *reslov) headerCheck(url, ip string) {
	if _, ok := checkres[ip]; ok {
		return
	}
	var c check
	// 	b := true
	// 	c.doing = &b
	// 	checkres[ip] = c
	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s", ip), nil)
	if err != nil {
		panic(err)
	}
	req.Host = strings.Split(url, "/")[0]
	type tresp struct {
		resp *http.Response
		err  error
	}
	myresp := make(chan tresp)
	go func(myresp chan tresp) {
		resp, err := client.Do(req)
		myresp <- tresp{resp, err}
	}(myresp)
	select {
	case p := <-myresp:
		if err != nil {
			c.header = append(c.header, err.Error())
			return
		}
		c.code = p.resp.StatusCode
		for _, v := range p.resp.Header {
			for _, v1 := range v {
				if strings.Contains(v1, "HIT") || strings.Contains(v1, "MISS") {
					c.header = append(c.header, v1)
				}
			}
		}
		checkres[ip] = c
	case <-time.After(3 * time.Second):
		c.code = 502
		c.header = append(c.header, "timeout")
		checkres[ip] = c
		return
	}
}
