// test-server
package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/goa-go/goa"
	"github.com/goa-go/goa/router"
)

func logger(c *goa.Context, next func()) {
	start := time.Now()

	fmt.Printf("[%s] <-- %s %s\n", start.Format("2006-6-2 15:04:05"), c.Method, c.URL)
	next()
	fmt.Printf("[%s] --> %s %s %d%s\n", time.Now().Format("2006-6-2 15:04:05"), c.Method, c.URL, time.Since(start).Nanoseconds()/1e6, "ms")
}

func xml(c *goa.Context) {
	type Address struct {
		City, State string
	}
	type Person struct {
		Id        int     `xml:"id,attr"`
		FirstName string  `xml:"name>first"`
		LastName  string  `xml:"name>last"`
		Age       int     `xml:"age"`
		Height    float32 `xml:"height,omitempty"`
		Married   bool
		Address
		Comment string `xml:",comment"`
	}
	xml := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
	xml.Comment = " Need more details. "
	xml.Address = Address{"Hanga Roa", "Easter Island"}
	c.XML(xml)
}

func json(c *goa.Context) {
	c.JSON(goa.M{
		"string": "string",
		"int":    1,
		"json": goa.M{
			"key": "value",
		},
	})
}

func setStatus(c *goa.Context) {
	code := c.Param("code")
	int, err := strconv.Atoi(code)
	if err != nil {
		c.Status(400).String("plz input int")
	} else {
		c.Status(int).String("ok")
	}
}

func hello(c *goa.Context) {
	name := c.Query("name")
	c.String("hello " + name)
}

func postForm(c *goa.Context) {
	value := c.PostForm("key")
	c.String("key: " + value)
}

func main() {
	app := goa.New()
	router := router.New()

	router.GET("/", func(c *goa.Context) {
		c.String("hello world")
	})
	router.GET("/xml", xml)
	router.GET("/json", json)
	router.GET("/redirect", func(c *goa.Context) {
		c.Redirect(301, "http://github.com")
	})
	router.GET("/status/:code", setStatus)
	router.GET("/hello", hello)
	router.POST("/postForm", postForm)

	app.Use(logger)
	app.Use(router.Routes())
	app.Listen(":3000")
}
