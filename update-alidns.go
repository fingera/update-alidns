package main

import (
	"flag"
	"fmt"
	"os"
	"net/http"
	"io/ioutil"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

func main() {
	var region string
	var ak string
	var sk string
	var dn string
	var rr string
	var t string
	var v string

	flag.StringVar(&region, "region", "cn-hangzhou", "AK")
	flag.StringVar(&ak, "ak", os.Getenv("AK"), "AK")
	flag.StringVar(&sk, "sk", os.Getenv("SK"), "SK")
	flag.StringVar(&dn, "dn", os.Getenv("DomainName"), "DomainName")
	flag.StringVar(&rr, "rr", os.Getenv("RR"), "RR")
	flag.StringVar(&t, "t", "A", "Type")
	flag.StringVar(&v, "v", os.Getenv("IP"), "Value")

	flag.Parse()

	if ak == "" {
		fmt.Println("Error: no AK!")
		return
	}
	// fmt.Printf("AK: %s\n", ak)
	if sk == "" {
		fmt.Println("Error: no SK!")
		return
	}
	if dn == "" {
		fmt.Println("Error: no DomainName!")
		return
	}
	fmt.Printf("DomainName: %s\n", dn)
	if rr == "" {
		fmt.Println("Error: no RR!")
		return
	}
	fmt.Printf("RR: %s\n", rr)
	if t == "" {
		fmt.Println("Error: no Type!")
		return
	}
	fmt.Printf("Type: %s\n", t)
	if v == "" {
		res, err := http.Get("https://ifconfig.me/ip")
		if err != nil || res.StatusCode != 200 {
			fmt.Printf("Error: making http request: %s\n", err)
			return
		}
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error: could not read response body: %s\n", err)
			return
		}
		v = string(resBody)
	}
	fmt.Printf("Value: %s\n", v)

	client, err := alidns.NewClientWithAccessKey(region, ak, sk)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	desc := alidns.CreateDescribeDomainRecordsRequest()
	desc.DomainName = dn
	desc.SearchMode = "EXACT"
	desc.KeyWord = rr
	existed, err := client.DescribeDomainRecords(desc)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	if existed.TotalCount == 0 {
		request := alidns.CreateAddDomainRecordRequest()
		request.Scheme = "https"
		request.DomainName = dn
		request.RR = rr
		request.Type = t
		request.Value = v
		response, err := client.AddDomainRecord(request)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
		fmt.Printf("OK: %s\n", response)
	} else {
		fmt.Printf("TotalCount: %#v\n", existed.TotalCount)
		var rid string = ""
		for _, r := range existed.DomainRecords.Record {
			if r.RR == rr {
				rid = r.RecordId
				break
			}
		}
		if len(rid) == 0 {
			fmt.Printf("Error: %s not found!\n", rr)
			return
		}

		request := alidns.CreateUpdateDomainRecordRequest()
		request.Scheme = "https"
		request.RecordId = rid
		request.RR = rr
		request.Type = t
		request.Value = v
		response, err := client.UpdateDomainRecord(request)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
		fmt.Printf("OK: %s\n", response)
	}
}
