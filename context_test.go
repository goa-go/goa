package goa

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSet(t *testing.T) {
	c := &Context{}
	c.Set("key", "value")
	value, _ := c.Get("key")

	assert.Equal(t, "value", value)
}

func TestQuery(t *testing.T) {
	c := &Context{}
	req, _ := http.NewRequest("GET", "/?name=nicholascao", nil)
	c.Request = req
	name := c.Query("name")
	name2 := c.Query("name2")

	assert.Equal(t, "nicholascao", name)
	assert.Equal(t, "", name2)
}

func TestPostForm(t *testing.T) {
	c := &Context{}
	req, _ := http.NewRequest("POST", "/", strings.NewReader("key=value"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	c.Request = req
	value := c.PostForm("key")

	assert.Equal(t, "value", value)
}

func TestFormFile(t *testing.T) {
	c := &Context{}
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", "test")
	if assert.NoError(t, err) {
		_, err = w.Write([]byte("test"))
		assert.NoError(t, err)
	}
	mw.Close()

	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())

	file, fh, err := c.FormFile("file")

	if assert.NoError(t, err) && assert.NotNil(t, file) {
		assert.Equal(t, "test", fh.Filename)
	}
}

func TestFormFileFailed(t *testing.T) {
	c := &Context{}
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	mw.Close()

	c.Request, _ = http.NewRequest("POST", "/", nil)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())

	file, fh, err := c.FormFile("file")
	assert.Error(t, err)
	assert.Nil(t, file)
	assert.Nil(t, fh)
}

func TestParam(t *testing.T) {
	c := &Context{}
	c.Params = Params{Param{
		Key:   "key",
		Value: "value",
	}}

	assert.Equal(t, "value", c.Param("key"))
	assert.Equal(t, "", c.Param("key2"))
}

type obj struct {
	Key string `json:"key" xml:"key" query:"key" form:"key"`
}

func TestParseJSON(t *testing.T) {
	c := &Context{}
	c.Request, _ = http.NewRequest("POST", "/", strings.NewReader("{\"key\":\"value\"}"))
	ptr := &obj{}
	err := c.ParseJSON(ptr)

	assert.Nil(t, err)
	assert.Equal(t, "value", ptr.Key)
}

func TestParseXML(t *testing.T) {
	c := &Context{}
	c.Request, _ = http.NewRequest("POST", "/", strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?>
	<root>
		<key>value</key>
	</root>`))
	ptr := &obj{}
	err := c.ParseXML(ptr)

	assert.Nil(t, err)
	assert.Equal(t, "value", ptr.Key)
}

func TestParseString(t *testing.T) {
	c := &Context{}
	c.Request, _ = http.NewRequest("POST", "/", strings.NewReader("string"))
	str, err := c.ParseString()

	assert.Nil(t, err)
	assert.Equal(t, "string", str)
}

func TestParseQuery(t *testing.T) {
	c := &Context{}
	c.Request, _ = http.NewRequest("GET", "/?key=value", nil)
	ptr := &obj{}
	err := c.ParseQuery(ptr)

	assert.Nil(t, err)
	assert.Equal(t, "value", ptr.Key)
}

func TestParseForm(t *testing.T) {
	c := &Context{}
	req, _ := http.NewRequest("POST", "/", strings.NewReader("key=value"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	c.Request = req
	ptr := &obj{}
	err := c.ParseForm(ptr)

	assert.Nil(t, err)
	assert.Equal(t, "value", ptr.Key)
}

func TestGetCookie(t *testing.T) {
	c := &Context{}
	c.ResponseWriter = httptest.NewRecorder()
	c.Request, _ = http.NewRequest("GET", "/getCookie", nil)
	c.Request.Header.Set("Cookie", "user=goa")
	cookie, _ := c.Cookie("user")
	assert.Equal(t, "goa", cookie)

	_, err := c.Cookie("nokey")
	assert.Error(t, err)
}

func TestSetCookie(t *testing.T) {
	c := &Context{}
	c.ResponseWriter = httptest.NewRecorder()
	c.SetCookie(&http.Cookie{
		Name:     "user",
		Value:    "goa",
		MaxAge:   1,
		Path:     "/",
		Domain:   "localhost",
		Secure:   true,
		HttpOnly: true,
	})
	assert.Equal(t, "user=goa; Path=/; Domain=localhost; Max-Age=1; HttpOnly; Secure", c.ResponseWriter.Header().Get("Set-Cookie"))
}

func TestStatus(t *testing.T) {
	c := &Context{}
	c.Status(200)
	assert.Equal(t, 200, c.GetStatus())
	c.Status(300)
	assert.Equal(t, 300, c.GetStatus())
	c.Status(400)
	assert.Equal(t, 400, c.GetStatus())
	c.Status(500)
	assert.Equal(t, 500, c.GetStatus())
}

func TestStatusFailed(t *testing.T) {
	c := &Context{}
	assert.Panics(t, func() { c.Status(10) })
}

func TestRespondJSON(t *testing.T) {
	c := &Context{}
	w := httptest.NewRecorder()
	c.ResponseWriter = w
	c.JSON(M{
		"key": "value",
	})
	c.writeContentType(c.ct)
	c.ResponseWriter.WriteHeader(c.status)
	c.respond(c.responser)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"key\":\"value\"}\n", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRespondXML(t *testing.T) {
	c := &Context{}
	w := httptest.NewRecorder()
	c.ResponseWriter = w
	c.XML(obj{"value"})
	c.writeContentType(c.ct)
	c.ResponseWriter.WriteHeader(c.status)
	c.respond(c.responser)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<obj><key>value</key></obj>", w.Body.String())
	assert.Equal(t, "application/xml; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRespondString(t *testing.T) {
	c := &Context{}
	w := httptest.NewRecorder()
	c.ResponseWriter = w
	c.String("string")
	c.writeContentType(c.ct)
	c.ResponseWriter.WriteHeader(c.status)
	c.respond(c.responser)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "string", w.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRespondHTML(t *testing.T) {
	c := &Context{}
	w := httptest.NewRecorder()
	c.ResponseWriter = w
	c.HTML("<html>html</html>")
	c.writeContentType(c.ct)
	c.ResponseWriter.WriteHeader(c.status)
	c.respond(c.responser)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<html>html</html>", w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRespondWithCustomStatus(t *testing.T) {
	c := &Context{}
	w := httptest.NewRecorder()
	c.ResponseWriter = w
	c.Status(http.StatusNotFound)
	c.String("Not Found")
	c.writeContentType(c.ct)
	c.ResponseWriter.WriteHeader(c.status)
	c.respond(c.responser)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Not Found", w.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRedirect(t *testing.T) {
	c := &Context{}
	c.Request, _ = http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	c.ResponseWriter = w

	assert.Panics(t, func() { c.Redirect(299, "/wrongpath") })
	assert.Panics(t, func() { c.Redirect(309, "/wrongpath") })

	c.Redirect(http.StatusMovedPermanently, "/path")
	assert.Equal(t, http.StatusMovedPermanently, w.Code)
	assert.Equal(t, "/path", w.Header().Get("Location"))
}

func TestSetHeader(t *testing.T) {
	c := &Context{}
	w := httptest.NewRecorder()
	c.ResponseWriter = w

	c.SetHeader("Content-Type", "text/plain")
	c.SetHeader("Custom-Header", "value")
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	assert.Equal(t, "value", w.Header().Get("Custom-Header"))

	c.SetHeader("Content-Type", "text/html")
	c.SetHeader("Custom-Header", "")
	assert.Equal(t, "text/html", w.Header().Get("Content-Type"))
	assert.Equal(t, "", w.Header().Get("Custom-Header"))
}

func TestContextError(t *testing.T) {
	c := &Context{}
	assert.Panics(t, func() { c.Error(404, http.StatusText(404)) })
	assert.Panics(t, func() { c.Error(500, http.StatusText(500)) })

	defer func() {
		err := recover()
		assert.Equal(t, Error{
			Code: 500,
			Msg:  http.StatusText(500),
		}, err.(Error))
	}()

	c.Error(500, http.StatusText(500))
}
