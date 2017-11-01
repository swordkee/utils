package utils

import (
	"sync"
	"time"
	"fmt"
)

type Application struct {
	smCache     sync.Map
	smCache_ext sync.Map
	expire      int
}

func NewCache(expire int) *Application {
	app := &Application{expire: expire}
	return app
}

func (app *Application) Store(key string, value interface{}) {
	ctime := getCTime()
	app.smCache.Store(key, value)
	app.smCache_ext.Store(key, ctime)

}

func (app *Application) Get(key string) (interface{}, bool) {
	return app.smCache.Load(key)
}

func (app *Application) ClearExpired() {
	cTime := getCTime()

	app.smCache_ext.Range(func(k, v interface{}) bool {
		vTime := v.(int)
		if vTime+60 < cTime {
			app.smCache_ext.Delete(k)
			app.smCache.Delete(k)
			fmt.Println("deleting key:", k)
		}
		return true
	})

}

func getCTime() (int) {

	return int(time.Now().Unix())
}

