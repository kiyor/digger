/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : modules.go

* Purpose :

* Creation Date : 09-22-2014

* Last Modified : Mon 15 Dec 2014 07:35:53 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type check struct {
	code   int
	header []string
}

var (
	checkres map[string]check
)

func (r *reslov) headerCheck(url, ip string) {
	if _, ok := checkres[ip]; ok {
		return
	}
	var c check
	checkres[ip] = c
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
	case <-time.After(20 * time.Second):
		c.code = 502
		c.header = append(c.header, "timeout")
		checkres[ip] = c
		return
	}
}
