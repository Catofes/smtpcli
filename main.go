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
	"io/ioutil"
)

func main() {
	loginInfo := os.Getenv("MAILINFO")
	to := flag.String("to", "", "receipt address")
	from := flag.String("from", "", "from address")
	fromName := flag.String("name", "", "from name")
	title := flag.String("title", "Hello", "email title")
	body := flag.String("body", "", "email body")
	attach := flag.String("attach", "", "email attachment")
	flag.Parse()
	if *body == "" {
		tmp, _ := ioutil.ReadAll(os.Stdin)
		tmpString := string(tmp)
		body = &tmpString
	}
	m := gomail.NewMessage()
	m.SetHeader("Sender", *from)
	m.SetHeader("From", m.FormatAddress(*from, *fromName))
	m.SetHeader("To", *to)
	m.SetHeader("Subject", *title)
	m.SetBody("text/html", *body)
	if *attach != "" {
		m.Attach(*attach)
	}
	info, err := url.Parse(loginInfo)
	if err != nil {
		log.Panic(err)
	}
	if info.Scheme != "smtp" {
		log.Panic(errors.New("login info should be smtp://username:password@host:port"))
	}
	host, p, err := net.SplitHostPort(info.Host)
	if err!= nil{
		log.Fatal(err)
	}
	port, err := strconv.Atoi(p)
	if err!=nil {
		log.Fatal(err)
	}
	password, _ := info.User.Password()
	d := gomail.NewPlainDialer(host, port, info.User.Username(), password)
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	fmt.Printf("Sending mail from [%s] to [%s] using server [%s:%d] with username [%s].\n", *from, *to,
		host, port, info.User.Username())
	if err := d.DialAndSend(m); err != nil {
		log.Panic(err)
	}
	fmt.Println("Done.")
}
