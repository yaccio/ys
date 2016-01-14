#ys (Yacc Serve)
`ys` is a simple command line utility that allows you to serve a single file or
directory over HTTPS. The file/directory will be downloadable on
https://address:8080/ and can be fetched using a browser or wget. Directories
will be served directly as a single .zip file.

#Installation
Currently only available as src, is compiled with go so can be installed as
follows:

> go get github.com/yaccio/ys

#Usage
You simply call ys on a file:

> ys <filename>


`ys` takes the following optional flags:

- `--hostname` set the hostname for ys (example: example.com:1337, :8080).
- `--pwd` sets a password used for basic authentication.


ys looks for the following optional env variables:

- YS_CERT and YS_PRIVKEY, if set will be use for SSL. If unavailable
a temporary generated certificate and private key will be used.
