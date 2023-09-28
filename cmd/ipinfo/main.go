package main

import (
    "flag"
    "log"
    "strings"

    "github.com/gryffyn/ipinfo/http"
    "github.com/gryffyn/ipinfo/iputil"
    "github.com/gryffyn/ipinfo/iputil/geo"
)

type multiValueFlag []string

func (f *multiValueFlag) String() string {
    return strings.Join(*f, ", ")
}

func (f *multiValueFlag) Set(v string) error {
    *f = append(*f, v)
    return nil
}

func init() {
    log.SetPrefix("ipinfo: ")
    log.SetFlags(log.Lshortfile)
}

func main() {
    countryFile := flag.String("f", "", "Path to GeoIP country database")
    cityFile := flag.String("c", "", "Path to GeoIP city database")
    asnFile := flag.String("a", "", "Path to GeoIP ASN database")
    listen := flag.String("l", ":8080", "Listening address")
    reverseLookup := flag.Bool("r", false, "Perform reverse hostname lookups")
    portLookup := flag.Bool("p", false, "Enable port lookup")
    cacheSize := flag.Int("C", 0, "Size of response cache. Set to 0 to disable")
    profile := flag.Bool("P", false, "Enables profiling handlers")
    sponsor := flag.Bool("s", false, "Show sponsor logo")
    var headers multiValueFlag
    flag.Var(&headers, "H", "Header to trust for remote IP, if present (e.g. X-Forwarded-For)")
    flag.Parse()
    if len(flag.Args()) != 0 {
        flag.Usage()
        return
    }

    r, err := geo.Open(*countryFile, *cityFile, *asnFile)
    if err != nil {
        log.Fatal(err)
    }
    cache := http.NewCache(*cacheSize)
    server := http.New(r, cache, *profile)
    server.IPHeaders = headers
    if *reverseLookup {
        log.Println("Enabling reverse lookup")
        server.LookupAddr = iputil.LookupAddr
    }
    if *portLookup {
        log.Println("Enabling port lookup")
        server.LookupPort = iputil.LookupPort
    }
    if *sponsor {
        log.Println("Enabling sponsor logo")
        server.Sponsor = *sponsor
    }
    if len(headers) > 0 {
        log.Printf("Trusting remote IP from header(s): %s", headers.String())
    }
    if *cacheSize > 0 {
        log.Printf("Cache capacity set to %d", *cacheSize)
    }
    if *profile {
        log.Printf("Enabling profiling handlers")
    }
    log.Printf("Listening on http://%s", *listen)
    if err := server.ListenAndServe(*listen); err != nil {
        log.Fatal(err)
    }
}
