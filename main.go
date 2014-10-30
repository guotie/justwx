package main

import (
	"flag"
	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/smtc/glog"
	"github.com/smtc/wxutils"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"net/http"
)

var (
	configFn = flag.String("config", "./config.json", "config file path")
	wxAuth   *wxutils.WXAuth
)

func main() {
	flag.Parse()
	config.ReadCfg(*configFn)

	glog.InitLogger(glog.PRO, map[string]interface{}{"typ": "file"})

	deferinit.InitAll()

	wxAuth = wxutils.CreateWXAuth(config.GetStringDefault("token", "weixin-token"))

	run()

	glog.Close()
}

func run() {
	weixinHandler := weixinMux()
	goji.Handle("/weixin/*", weixinHandler)
	goji.Get("/weixin", http.RedirectHandler("/weixin/", 301))

	goji.Get("/assets/*", http.FileServer(http.Dir("./")))
	goji.Serve()
}

func weixinMux() *web.Mux {
	mux := web.New()

	mux.Get("/weixin/", indexHandler)
	//mux.Get(regexp.MustCompile(`^/weixin/(?P<fn>.+).html$`), tplHandler)

	return mux
}

func indexHandler(c web.C, w http.ResponseWriter, req *http.Request) {
	var (
		err       error
		signature string
		timestamp string
		nonce     string
		echostr   string
	)

	signature = c.URLParams["signature"]
	timestamp = c.URLParams["timestamp"]
	nonce = c.URLParams["nonce"]
	echostr = c.URLParams["echostr"]

	err = wxAuth.CheckSignature(signature, timestamp, nonce)
	if err != nil {
		glog.Error("Check Signature failed: %v\n", err)
	}

	w.Write([]byte(echostr))
}
