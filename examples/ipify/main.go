package main

import (
	"fmt"
	"github.com/lithdew/nicehttp"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func exampleDownloadBytes(client *nicehttp.Client) {
	buf, err := client.DownloadBytes(nil, "https://api.ipify.org")
	check(err)

	fmt.Printf("Your IP is: %q\n", string(buf))
}

func exampleDownloadSerially(client *nicehttp.Client) {
	var b strings.Builder
	check(client.DownloadSerially(&b, "https://api.ipify.org"))

	fmt.Printf("Your IP is: %q\n", b.String())
}

func main() {
	client := nicehttp.NewClient()
	exampleDownloadBytes(&client)
	exampleDownloadSerially(&client)
}
