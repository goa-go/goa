# Goa
Goa is a using middleware web framework for golang.
goa = go + [koajs](https://github.com/koajs/koa)

##  Hello Goa
```go
func main() {
	app := New()

	app.Use(func(w http.ResponseWriter, r *http.Request, next func()) {
    w.Write([]byte("Hello Goa!"))
    next()
	})
	app.Listen(":3000")
}
```
