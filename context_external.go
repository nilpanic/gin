package gin

import (
	"github.com/nilpanic/gin/internal/rsp"
	"github.com/nilpanic/gin/internal/valid/binding"
	r "github.com/nilpanic/gin/rsp"
	"io"
	"net/http"
)

func (c *Context) Valid(obj any) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}

// ValidJSON is a shortcut for c.ValidWith(obj, binding.JSON).
func (c *Context) ValidJSON(obj any) error {
	return c.ShouldBindWith(obj, binding.JSON)
}

// ValidQuery is a shortcut for c.ValidWith(obj, binding.Query).
func (c *Context) ValidQuery(obj any) error {
	return c.ShouldBindWith(obj, binding.Query)
}

// ValidHeader is a shortcut for c.ShouldBindWith(obj, binding.Header).
func (c *Context) ValidHeader(obj any) error {
	return c.ShouldBindWith(obj, binding.Header)
}

// ValidWith binds the passed struct pointer using the specified binding engine.
// See the binding package.
func (c *Context) ValidWith(obj any, b binding.Binding) error {
	return b.Bind(c.Request, obj)
}

// ValidBodyWith is similar with ValidWith, but it stores the request
// body into the context, and reuse when it is called again.
//
// NOTE: This method reads the body before binding. So you should use
// ValidBodyWith for better performance if you need to call only once.
func (c *Context) ValidBodyWith(obj any, bb binding.BindingBody) (err error) {
	var body []byte
	if cb, ok := c.Get(BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			body = cbb
		}
	}
	if body == nil {
		body, err = io.ReadAll(c.Request.Body)
		if err != nil {
			return err
		}
		c.Set(BodyBytesKey, body)
	}
	return bb.BindBody(body, obj)
}

// JSONOk response json success data.
func (c *Context) JSONOk(val ...rsp.JSVal) {
	rel := &rsp.JSONVal{
		Code: r.CodeOK,
		Msg:  r.MsgSuccess,
	}

	for _, jsonVal := range val {
		jsonVal(rel)
	}

	c.JSON(http.StatusOK, rel)
}

// JSONErr response json error data.
func (c *Context) JSONErr(val ...rsp.JSVal) {
	rel := &rsp.JSONVal{
		Code: r.CodeErr,
		Msg:  r.MsgFailed,
	}

	for _, jsonVal := range val {
		jsonVal(rel)
	}

	c.JSON(http.StatusOK, rel)
}

// JSONPOk response jsonp success data.
func (c *Context) JSONPOk(val ...rsp.JSVal) {
	rel := &rsp.JSONVal{
		Code: r.CodeOK,
		Msg:  r.MsgSuccess,
	}

	for _, jsonVal := range val {
		jsonVal(rel)
	}

	c.JSONP(http.StatusOK, rel)
}

// JSONPErr response jsonp error data.
func (c *Context) JSONPErr(val ...rsp.JSVal) {
	rel := &rsp.JSONVal{
		Code: r.CodeErr,
		Msg:  r.MsgFailed,
		Data: nil,
	}

	for _, jsonVal := range val {
		jsonVal(rel)
	}

	c.JSONP(http.StatusOK, rel)
}
