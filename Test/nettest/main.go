package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, err := http.Get("https://blog.csdn.net/sszgg2006/article/details/73342566/")
	if err != nil {
		fmt.Printf(" 获取网页失败，err %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error Code: ", resp.StatusCode)
		return
	}
	ecoding, err := determineEcoding(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	utf8Reader := transform.NewReader(resp.Body, ecoding.NewDecoder())
	bytes, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bytes))
}

func determineEcoding(r io.Reader) (encoding.Encoding, error) {
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		return nil, err
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e, nil
}
