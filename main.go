package main

import (
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

	http.HandleFunc("/", rootOnly(logging(auth(*pwd, serve(filename)))))

	log.Println("Starting ys on", *hostname)
	err := http.ListenAndServe(*hostname, nil)
	if err != nil {
		log.Println("Listening error:", err.Error())
	}
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
			_, p, _ := req.BasicAuth()
			if pwd == p {
				h(res, req)
				return
			} else {
				res.WriteHeader(http.StatusUnauthorized)
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
			res.Write(zipDir(file))
		} else {
			b, err := ioutil.ReadAll(file)
			errPanic(err)
			res.Header().Add("Content-Disposition",
				fmt.Sprintf("inline; filename=\"%s\"", f))
			res.Write(b)
		}
	}
}

func zipDir(file *os.File) []byte {
	panic("Not implemented")
}

func errPanic(e error) {
	if e != nil {
		panic(e)
	}
}
