package Base

import (
	"net/http"
	"fmt"
	"regexp"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter,r *http.Request)
}

type MyMux struct {
	Routers map[string]func(w http.ResponseWriter,r *http.Request)
}

func(ro *MyMux) ServeHTTP(w http.ResponseWriter,r *http.Request){
	for path, f := range ro.Routers{
		if ok, _ :=regexp.MatchString("^" + path + "$", r.URL.Path); ok{
		//if r.URL.Path == path{ //简单方法
			f(w,r)
			return
		}
	}
	fmt.Fprintf(w,"错误:访问路径不正确Url: '%s'",r.URL.Path)
}

func (ro *MyMux) MyMuxInit(){
	ro.Routers = make(map[string]func(w http.ResponseWriter,r *http.Request))

}
