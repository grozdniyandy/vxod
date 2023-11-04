package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/PuerkitoBio/goquery"
    "crypto/tls"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go domain1.com domain2.com ...")
        return
    }

    domains := os.Args[1:]

    for _, domain := range domains {
        fmt.Printf("Checking %s...\n", domain)
        checkForInputFields("http://" + domain)
        checkForInputFields("https://" + domain)
    }
}

func checkForInputFields(url string) {
    httpClient := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }

    resp, err := httpClient.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    hasInputFields := false

    doc.Find("input").Each(func(i int, s *goquery.Selection) {
        inputType, exists := s.Attr("type")
        if !exists || (inputType != "hidden" && inputType != "submit" && inputType != "reset") {
            hasInputFields = true
            return
        }
    })

    if hasInputFields {
        fmt.Printf("%s contains input fields.\n", url)
    } else {
        fmt.Printf("%s does not contain input fields.\n", url)
    }
}
