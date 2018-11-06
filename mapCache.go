package utils

import (
	"fmt"
	"sync"
	"time"
)

type Application struct {
	smCache    sync.Map
	smCacheExt sync.Map
	expire     int
}

func NewCache(expire int) *Application {
	app := &Application{expire: expire}
	return app
}

func (app *Application) Store(key string, value interface{}) {
	i := getCTime()
	app.smCache.Store(key, value)
	app.smCacheExt.Store(key, i)

}

func (app *Application) Get(key string) (interface{}, bool) {
	return app.smCache.Load(key)
}

func (app *Application) ClearExpired() {
	cTime := getCTime()

	app.smCacheExt.Range(func(k, v interface{}) bool {
		vTime := v.(int)
		if vTime+60 < cTime {
			app.smCacheExt.Delete(k)
			app.smCache.Delete(k)
			fmt.Println("deleting key:", k)
		}
		return true
	})

}

func getCTime() int {

	return int(time.Now().Unix())
}
