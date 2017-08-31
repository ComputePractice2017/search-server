package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	//"github.com/opesun/goquery"
)

var (
	WORKERS   int    = 2
	URL_FILE  string = "url.txt"
	TEXT_FILE string = "text.txt"
)

func init() {
	flag.IntVar(&WORKERS, "w", WORKERS, "количество потоков")
	flag.StringVar(&URL_FILE, "hf", URL_FILE, "файл url")
	flag.StringVar(&TEXT_FILE, "qf", TEXT_FILE, "файл text")
	flag.Parse()
}

func _check(err error) {
	if err != nil {
		panic(err)
	}
}

func par() <-chan string {
	c := make(chan string)
	for i := 0; i < WORKERS; i++ {
		go func() {
			x, err := goquery.NewDocument("https://habrahabr.ru/post/150134/")
			_check(err)
			x.Find("a").Each(func(i int, s *goquery.Selection) {
				link, _ := s.Attr("href")
				c <- link
			})
			time.Sleep(100 * time.Millisecond)
		}()
	}
	fmt.Println("Запущено потоков: ", WORKERS)
	return c
}

func main() {
	url_file, err := os.OpenFile(URL_FILE, os.O_APPEND|os.O_CREATE, 0666)
	_check(err)
	defer url_file.Close()
	url_count := 0
	url_chan := par()
	//for {
	url := <-url_chan
	url_count++
	url_file.WriteString(url + "\n\n\n")
	//}
}
