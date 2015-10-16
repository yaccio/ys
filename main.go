package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type handle func(http.ResponseWriter, *http.Request)

func main() {
	addr := flag.String("addr", ":1111", "Set the address of ys")
	pwd := flag.String("pwd", "", "The password used to protect the content servered by ys")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Please provide file or folder to serve:")
		fmt.Println("ys <dir/file>")
		return
	}

	filename := flag.Arg(0)

	http.HandleFunc("/", auth(*pwd, serve(filename)))

	http.ListenAndServe(*addr, nil)
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
			}
		}
	}
}

func serve(f string) handle {
	return func(res http.ResponseWriter, req *http.Request) {
		file, err := os.Open(f)
		E(err)

		stat, err := file.Stat()
		E(err)
		if stat.IsDir() {
			res.Header().Add("Content-Disposition",
				fmt.Sprintf("attachment; filename=\"%s.zip\"", f))
			res.Write(zipDir(file))
		} else {
			b, err := ioutil.ReadAll(file)
			E(err)
			res.Header().Add("Content-Disposition",
				fmt.Sprintf("attachment; filename=\"%s\"", f))
			res.Write(b)
		}
	}
}

func zipDir(file *os.File) []byte {
	panic("Not implemented")
}

func E(e error) {
	if e != nil {
		panic(e)
	}
}
