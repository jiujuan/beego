package beego

import (
    "mime"
    "net/http"
    "strings"
)

type Context struct {
    ResponseWriter http.ResponseWriter
    Request *http.Request
    Params  map[string]string
}

func (ctx *Context) WriteString(content string) {
    ctx.ResponseWriter.Write([]byte(content))
}

func (ctx *Context) Abort(status int, body string) {
    ctx.ResponseWriter.WriteHeader(status)
    ctx.ResponseWriter.Write([]byte(body))
}

func (ctx *Context) Redirect(status int, url string) {
    ctx.ResponseWriter.Header().Set("Location", url)
    ctx.ResponseWriter.WriteHeader(status)
    ctx.ResponseWriter.Write([]byte("Redirecting to: " + url))
}

func (ctx *Context) NotFound(message string) {
    ctx.ResponseWriter.WriteHeader(404)
    ctx.ResponseWriter.Write([]byte(message))
}

func (ctx *Context) ContentType(ext string) {
    if !strings.HasPrefix(ext, ".") {
        ext = "." + ext
    }
    ctype := mime.TypeByExtension(ext)
    if ctype != "" {
        ctx.ResponseWriter.Header().Set("Content-Type", ctype)
    }
}

func (ctx *Context) SetHeader(hdr string, val string, unique bool) {
    if unique {
        ctx.ResponseWriter.Header().Set(hdr, val)
    } else {
        ctx.ResponseWriter.Header().Add(hdr, val)
    }
}