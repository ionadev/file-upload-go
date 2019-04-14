package main

import (
	"io"
	"os"

	"github.com/kataras/iris"
)

const maxSize = 5 << 30

func main() {
	app := iris.New()

	app.RegisterView(iris.HTML("./templates", ".html"))

	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})

	app.Post("/upload", iris.LimitRequestBodySize(maxSize), func(ctx iris.Context) {
		file, info, err := ctx.FormFile("file[]")

		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.HTML("<p> An Error Occured: " + err.Error() + "</p>")
			return
		}

		defer file.Close()

		fileName := info.Filename

		out, err := os.OpenFile("./uploads/"+fileName,
			os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.HTML("<p> An Error Occured: " + err.Error() + "</p>")
			return
		}

		defer out.Close()

		io.Copy(out, file)

		ctx.JSON(iris.Map{
			"success":  true,
			"fileName": fileName})
	})
	app.Run(iris.Addr(":241"))
}
