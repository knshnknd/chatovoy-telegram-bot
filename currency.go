package main

 import (
   "encoding/xml"
   "fmt"
   "io/ioutil"
   "log"
   "net/http"
 )

 type Envelope struct {
         Cube []struct {
                 Date  string `xml:"time,attr"`
                 Rates []struct {
                         Currency string `xml:"currency,attr"`
                         Rate     string `xml:"rate,attr"`
                 } `xml:"Cube"`
         } `xml:"Cube>Cube"`
 }

 func getCurrency() string {
   resp, err := http.Get("http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")

   if err != nil {
     log.Fatal(err)
   }

   defer resp.Body.Close()

   xmlCurrenciesData, err := ioutil.ReadAll(resp.Body)

   if err != nil {
     log.Fatal(err)
   }

   var env Envelope
   err = xml.Unmarshal(xmlCurrenciesData, &env)

   if err != nil {
     log.Fatal(err)
   }

   fmt.Println("Date ", env.Cube[0].Date)

   var reply string

   for _, v := range env.Cube[0].Rates {
     if v.Currency == "USD" || v.Currency == "JPY" || v.Currency == "GBP" {
       reply = reply + "1 евро = " + v.Rate + " " +v.Currency + "\n"
     }
   }

   return reply
 }
