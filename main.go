package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/glesys/glesys-go/v8"
	"github.com/rwhelan/gozone"
	"github.com/xyproto/randomstring"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseRecordtoTF(record gozone.Record) (string, error) {
	r := randomstring.HumanFriendlyString(6)
	id := "dns-" + r

	datat := strings.Join(record.Data, " ")
	if datat == record.Origin {
		datat = "@"
	}
	fmt.Println("DEBUG: domainname", record.DomainName)
	return fmt.Sprintf(`
resource "glesys_dnsdomain_record" "%s" {
	domain = "%s"
	data   = "%s"
	host   = "%s"
	ttl    = %d
	type   = "%s"
} `,
		id, record.Origin, datat, record.DomainName, record.TimeToLive, record.Type), nil

}

func main() {
	agent := "norrland-dns-export/0.0.1"

	userid := os.Getenv("GLESYS_USERID")
	token := os.Getenv("GLESYS_TOKEN")
	apiurl := os.Getenv("GLESYS_APIURL")
	if len(apiurl) == 0 {
		apiurl = "https://api.glesys.com"
	}

	client := glesys.NewClient(userid, token, agent)
	client.SetBaseURL(apiurl)

	now := time.Now()
	ts := now.Unix()
	timestamp := strconv.FormatInt(ts, 10)

	domain := flag.String("domain", "example.com", "Domain name")
	inFile := flag.String("in", "", "input zone file")

	outFile := flag.String("out", "domain-out.txt", "Output file for domain export")
	outTf := flag.String("tf", "dns.tf", "Terraform output")
	//diskType := os.Args[3]
	flag.Parse()

	tf, err := os.Create(fmt.Sprintf("./%s", *outTf))
	check(err)
	wTf := bufio.NewWriter(tf)

	if len(*inFile) > 0 {
		fmt.Printf("Checking file on disk %s\n", *inFile)
		stream, err := os.Open(*inFile)
		check(err)
		var record gozone.Record
		scanner := gozone.NewScanner(stream)

		tfData := []string{}
		for {
			err := scanner.Next(&record)
			if err != nil {
				break
			}

			//fmt.Printf("hejhej %s\n", record.String())
			if record.Type.String() == "SOA" {
				continue
			} else {
				rec, err := parseRecordtoTF(record)
				check(err)
				tfData = append(tfData, rec)
			}

		}
		tfData1 := strings.Join(tfData[:], "\n")

		w1, err := wTf.WriteString(tfData1)
		check(err)

		fmt.Printf("wrote tfout - %d bytes\n", w1)
		wTf.Flush()

	} else {
		// If no zonefile provided, run export from glesys dns
		f, err := os.Create(fmt.Sprintf("./%s-%s", *outFile, timestamp))
		check(err)
		w := bufio.NewWriter(f)
		domainData, err := client.DNSDomains.Export(context.Background(), *domain)
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(20)
		}
		n1, err := w.WriteString(domainData)
		if err != nil {
			fmt.Println("error writing data")
		}
		fmt.Printf("wrote %d bytes\n", n1)
		w.Flush()
	}
}
