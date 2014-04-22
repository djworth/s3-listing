package main

import (
	"flag"
	"fmt"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/s3"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

var bucket string

func init() {
	flag.StringVar(&bucket, "bucket", "", "Name of the bucket in S3 to generate a listing page")
}

func executeListingPage(listing *s3.ListResp, list *template.Template, filename string) error {
	f, err := os.Create(filename)
	defer f.Close()

	if err != nil {
		return err
	}

	list.Execute(f, listing)

	return nil
}

func main() {

	flag.Parse()
	if bucket == "" {
		flag.Usage()
		return
	}

	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatalln(err)
	}

	client := s3.New(auth, aws.USEast)

	b := client.Bucket(bucket)

	listing, err := b.List("", "", "", 1000)

	if err != nil {
		log.Fatalln(err)
	}

	list, err := template.ParseFiles("list.html")
	if err != nil {
		log.Fatalln(err)
	}

	err = executeListingPage(listing, list, "index.html")
	if err != nil {
		log.Fatalln(err)
	}
	data, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Uploading index.html to S3...")

	err = b.Put("index.html", data, "text/html", s3.PublicRead, s3.Options{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Finished")
}
