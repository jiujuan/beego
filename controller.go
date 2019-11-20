package beego

import (
    "bytes"
    "encoding/json"
    "html/template"
    "io/ioutil"
    "net/http"
    "path"
    "strconv"
)

type ControllerInterface interface {
    Init(ctx *Context, cn string)
    Prepare()
    Get()
    Post()
    Delete()
    Put()
    Head()
    Options()
    Finish()
    Render() error
}

type Controller struct {
    Ctx        *Context
    Tpl       *template.Template
    Data      map[interface{}]interface{}
    ChildName string
    TplNames  string
    Layout    string
    TplExt    string
}

func (c *Controller) Init(ct *Context, cn string) {
    c.Data = make(map[interface{}]interface{})
    c.Tpl = template.New(cn + ct.Request.Method)
    c.Layout = ""
    c.TplNames = ""
    c.ChildName = cn
    c.Ctx = ct
    c.TplExt = "tpl"

}

func (c *Controller) Prepare() {

}

func (c *Controller) Finish() {

}

func (c *Controller) Get() {
    http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Post() {
    http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Delete() {
    http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Put() {
    http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Head() {
    http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Patch() {
    http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Options() {
    http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Render() error {
    if c.Layout != "" {
        if c.TplNames == "" {
            c.TplNames = c.ChildName + "/" + c.Ctx.Request.Method + "." + c.TplExt
        }
        t, err := c.Tpl.ParseFiles(path.Join(ViewsPath, c.TplNames), path.Join(ViewsPath, c.Layout))
        if err != nil {
            Trace("template ParseFiles err:", err)
        }
        _, file := path.Split(c.TplNames)
        newbytes := bytes.NewBufferString("")
        t.ExecuteTemplate(newbytes, file, c.Data)
        tplcontent, _ := ioutil.ReadAll(newbytes)
        c.Data["LayoutContent"] = template.HTML(string(tplcontent))
        _, file = path.Split(c.Layout)
        err = t.ExecuteTemplate(c.Ctx.ResponseWriter, file, c.Data)
        if err != nil {
            Trace("template Execute err:", err)
        }
    } else {
        if c.TplNames == "" {
            //c.TplNames = c.ChildName + "/" + c.Ctx.Request.Method + "." + c.TplExt
            return nil
        }
        t, err := c.Tpl.ParseFiles(path.Join(ViewsPath, c.TplNames))
        if err != nil {
            Trace("template ParseFiles err:", err)
        }
        _, file := path.Split(c.TplNames)
        err = t.ExecuteTemplate(c.Ctx.ResponseWriter, file, c.Data)
        if err != nil {
            Trace("template Execute err:", err)
        }
    }
    return nil
}

func (c *Controller) Redirect(url string, code int) {
    c.Ctx.Redirect(code, url)
}

func (c *Controller) ServeJson() {
    content, err := json.MarshalIndent(c.Data, "", "  ")
    if err != nil {
        http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
        return
    }
    c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
    c.Ctx.ContentType("json")
    c.Ctx.ResponseWriter.Write(content)
}