package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed showmovie.html
var showmovie string


//go:embed artplayer.js
//go:embed jquery-3.5.1.min.js
var ef embed.FS

func main() {
	//gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	app.StaticFS("/v", http.Dir("./"))
	app.StaticFS("/e", http.FS(ef))
	//app.LoadHTMLFiles("showmovie.html")
	//app.StaticFile("/artplayer.js", "./artplayer.js")
	//app.StaticFile("/jquery-3.5.1.min.js", "./jquery-3.5.1.min.js")
	html := template.Must(template.New("showmovie.html").Parse(showmovie))
	app.SetHTMLTemplate(html)
	app.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "showmovie.html", gin.H{
			"title": "Main website",
		})

	})

	app.GET("/list", func(ctx *gin.Context) {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		dir=strings.ReplaceAll(dir, "\\", "/")
		if !strings.HasSuffix(dir,"/"){
			dir=dir+"/"
		}

		list,_:=WalkDir(dir,"mp4")
		liststr:=""
		for _,v:=range(list) {
			tmp:=strings.ReplaceAll(v, dir, "")
			liststr+=tmp+","
			liststr="/v/"+liststr
		}

		ctx.String(http.StatusOK,liststr)

	})

	servePort:="80"
	if len(os.Args) == 2 {
		_,errport:=strconv.Atoi(os.Args[1])
		if errport==nil {
			servePort=os.Args[1]
		}
	}
	app.Run(":"+servePort)
}

func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil {
			fmt.Println(err.Error())
		}

		if fi.IsDir() { // 忽略目录
			return nil
		}

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			fname:=strings.ReplaceAll(filename, "\\", "/")
			files = append(files, fname)
		}

		return nil
	})

	return files, err
}