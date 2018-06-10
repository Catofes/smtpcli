package main

import (
	"flag"
	"os"
	"gopkg.in/gomail.v2"
	"net/url"
	"errors"
	"log"
	"strconv"
	"net"
	"fmt"
	"crypto/tls"
)

func main() {
	loginInfo := os.Getenv("MAILINFO")
	to := flag.String("to", "", "receipt address")
	from := flag.String("from", "", "from address")
	title := flag.String("title", "Hello", "email title")
	body := flag.String("body", "Wow", "email body")
	flag.Parse()
	m := gomail.NewMessage()
	m.SetHeader("From", *from)
	m.SetHeader("To", *to)
	m.SetHeader("Subject", *title)
	m.SetBody("text/html", *body)
	info, err := url.Parse(loginInfo)
	if err != nil {
		log.Panic(err)
	}
	if info.Scheme != "smtp" {
		log.Panic(errors.New("login info should be smtp://username:password@host:port"))
	}
	host, p, _ := net.SplitHostPort(info.Host)
	port, _ := strconv.Atoi(p)
	password, _ := info.User.Password()
	d := gomail.NewPlainDialer(host, port, info.User.Username(), password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	fmt.Printf("Sending mail from [%s] to [%s] using server [%s:%d] with username [%s].\n", *from, *to,
		host, port, info.User.Username())
	if err := d.DialAndSend(m); err != nil {
		log.Panic(err)
	}
	fmt.Println("Done.")
}
