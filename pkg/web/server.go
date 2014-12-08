package web

import (
	"net/http"
	"os"

	"github.com/go-martini/martini"

	"github.com/telehash/interoper/pkg/web/static"
)

func Run(report []byte) {
	m := martini.Classic()

	m.Use(staticFiles())

	m.Get("/dump.json", func(rw http.ResponseWriter) {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(200)
		rw.Write(report)
	})

	m.Get("/**", func(rw http.ResponseWriter, req *http.Request) {
		http.Redirect(rw, req, "/#"+req.URL.Path, 301)
	})

	m.Run()
}

func staticFiles() martini.Handler {
	fs := static.FS(false)

	return func(ctx martini.Context, rw http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			ctx.Next()
			return
		}

		path := req.URL.Path
		if path == "" || path == "/" {
			path = "/index.html"
		}

		f, err := fs.Open(path)
		if os.IsNotExist(err) {
			ctx.Next()
			return
		}
		if err != nil {
			panic(err)
		}
		f.Close()

		http.FileServer(fs).ServeHTTP(rw, req)
	}
}
