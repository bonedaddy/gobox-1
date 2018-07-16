package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego/context"
	//"github.com/muesli/cache2go"
	"github.com/satori/go.uuid"
)

// !TODO // !DEBUG

var CSRF *CSRFManager

var ERR_INVALID_CSRF_TOKEN error = errors.New("Invalid or expired CSRF token!")

type CSRFManager struct {
	//Cache *cache2go.CacheTable
}

func (c *CSRFManager) getTokenData(scope string, ctx *context.Context) (string, time.Time) {

	token, err1 := ctx.Input.CruSession.Get(fmt.Sprintf("csrf/%s/token", scope)).(string)
	expire, err2 := ctx.Input.CruSession.Get(fmt.Sprintf("csrf/%s/expire", scope)).(time.Time)

	return token, expire
}

func (c *CSRFManager) SetToken(scope string, lifetime time.Duration, ctx *context.Context) string {

	if int(lifetime) == 0 {
		lifetime = (24 * 7) * time.Hour
	}

	etoken, eexpire := c.getTokenData(scope, ctx)

	if etoken != "" {
		ctx.Input.CruSession.Set(fmt.Sprintf("csrf/%s/expire", scope), time.Now().Add(lifetime))
		return etoken
	}
	tid, _ := uuid.NewV4()
	token := tid.String()
	expire := time.Now().Add(lifetime)
	err := ctx.Input.CruSession.Set(fmt.Sprintf("csrf/%s/token", scope), token)
	fmt.Println(err)
	err = ctx.Input.CruSession.Set(fmt.Sprintf("csrf/%s/expire", scope), expire)
	fmt.Println(err)

	return token
}

func (c *CSRFManager) ValidateToken(scope string, inputToken string, ctx *context.Context) error {

	token, expire := c.getTokenData(scope, ctx)
	if time.Now().After(expire) {
		return ERR_INVALID_CSRF_TOKEN
	}
	if token != inputToken {
		return ERR_INVALID_CSRF_TOKEN
	}
	return nil
}
