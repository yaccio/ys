package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type handle func(http.ResponseWriter, *http.Request)

func main() {
	hostname := flag.String("hostname", ":8080", "Set the hostname used by ys")
	pwd := flag.String("pwd", "", "The password used to protect the content servered by ys")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Please provide file or folder to serve:")
		fmt.Println("ys <dir/file>")
		return
	}

	filename := flag.Arg(0)

	StartHTTPS(YsHandler{*hostname, *pwd, filename})
}

func StartHTTPS(ys YsHandler) {
	log.Println("Starting ys on", ys.Hostname)
	var err error
	certfile := os.Getenv("YS_CERT")
	privkey := os.Getenv("YS_PRIVKEY")
	if certfile == "" || privkey == "" {
		cert, err := generateKeys()
		errPanic(err)
		server := &http.Server{
			Addr: ys.Hostname,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{*cert},
			},
			Handler: ys,
		}
		err = server.ListenAndServeTLS("", "")
	} else {
		err = http.ListenAndServeTLS(ys.Hostname, certfile, privkey, ys)
	}

	if err != nil {
		log.Println("HTTPS Listening error:", err.Error())
	}
}

type YsHandler struct {
	Hostname string
	Password string
	Filename string
}

func (ys YsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	(rootOnly(logging(auth(ys.Password, serve(ys.Filename)))))(res, req)
}

func rootOnly(h handle) handle {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/" {
			h(res, req)
		}
	}
}

func logging(h handle) handle {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("[%s] %s %s\n", req.Method, req.RemoteAddr, req.RequestURI)
		h(res, req)
	}
}

func auth(pwd string, h handle) handle {
	if pwd == "" {
		return h
	} else {
		return func(res http.ResponseWriter, req *http.Request) {
			res.Header().Add("WWW-Authenticate", "Basic realm=\"ys:\"")
			_, p, _ := req.BasicAuth()
			if pwd == p {
				h(res, req)
				return
			} else {
				res.WriteHeader(http.StatusUnauthorized)
				res.Write([]byte("Unauthorized"))
				log.Println("Denied request from %s", req.RemoteAddr)
			}
		}
	}
}

func serve(f string) handle {
	return func(res http.ResponseWriter, req *http.Request) {
		file, err := os.Open(f)
		errPanic(err)

		stat, err := file.Stat()
		errPanic(err)
		if stat.IsDir() {
			res.Header().Add("Content-Disposition",
				fmt.Sprintf("inline; filename=\"%s.zip\"", f))
			err := zipDir(f, res)
			errPanic(err)
		} else {
			b, err := ioutil.ReadAll(file)
			errPanic(err)
			res.Header().Add("Content-Disposition",
				fmt.Sprintf("inline; filename=\"%s\"", f))
			res.Write(b)
		}
	}
}

func errPanic(e error) {
	if e != nil {
		panic(e)
	}
}
