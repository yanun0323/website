package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"website/internal/app"
	"website/pkg/config"

	"github.com/labstack/echo/v4"
)

var (
	l   *log.Logger
	ctx context.Context
)

func main2() {
	l = log.Default()
	ctx = context.Background()
	if err := config.Init("config"); err != nil {
		l.Fatalf("init config failed %s", err)
		return
	}

	ch := make(chan *echo.Echo, 1)

	go func() {
		e := app.Run()
		ch <- e
		l.Fatal(e.Start(":80"))
		// l.Fatal(e.StartAutoTLS(":443"))
	}()

	for {
		time.Sleep(10 * time.Minute)
		go func() {
			c := <-ch
			e := app.Run()
			ch <- e
			l.Fatal(c.Shutdown(ctx))
			l.Fatal(e.Start(":80"))
			// l.Fatal(e.StartAutoTLS(":443"))
		}()
	}
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() //解析參數，預設是不會解析的

	fmt.Println(r.Form) //這些資訊是輸出到伺服器端的列印資訊

	fmt.Println("path", r.URL.Path)

	fmt.Println("scheme", r.URL.Scheme)

	fmt.Println(r.Form["url_long"])

	for k, v := range r.Form {

		fmt.Println("key:", k)

		fmt.Println("val:", strings.Join(v, ""))

	}

	fmt.Fprintf(w, "Hello astaxie!") //這個寫入到 w 的是輸出到客戶端的
}

func main() {

	http.HandleFunc("/", sayhelloName) //設定存取的路由

	err := http.ListenAndServe(":8080", nil) //設定監聽的埠

	if err != nil {
		fmt.Println(err)
		log.Fatal("ListenAndServe: ", err)
	}

}
