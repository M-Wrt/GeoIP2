package main

import (
	"bufio"
	"flag"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	log "github.com/sirupsen/logrus"
	"os"
)
		//"country": mmdbtype.Map{
		//"continent": mmdbtype.Map{
		//"code": mmdbtype.String("AS"),
		//"geoname_id": mmdbtype.Uint32(6255147),
		//"names":mmdbtype.Map{"de":mmdbtype.String("Asien"),
		//"en":mmdbtype.String("Asia"),
		//"es":mmdbtype.String("Asia"),
		//"fr":mmdbtype.String("Asie"),
		//"ja":mmdbtype.String("アジア"),
		//"pt-BR":mmdbtype.String("Ásia"),
		//"zh":mmdbtype.String("亞洲")},
		//},
		//"country": mmdbtype.Map{
		//"geoname_id":mmdbtype.Uint32(1668284),
		//"is_in_european_union":mmdbtype.Bool(false),
		//"iso_code":mmdbtype.String("CN"),
		//"names":mmdbtype.Map{
		//"de":mmdbtype.String("Fallen China"),
		//"en":mmdbtype.String("Fallen China"),
		//"es":mmdbtype.String("Fallen China"),
		//"fr":mmdbtype.String("La Chine perdue"),
		//"ja":mmdbtype.String("中國大陸淪陷區"),
		//"pt-BR":mmdbtype.String("Fallen China"),
		//"zh":mmdbtype.String("中國大陸淪陷區"),
		//},
		//},
		//"registered_country": mmdbtype.Map{
		//"geoname_id":mmdbtype.Uint32(1668284),
		//"is_in_european_union":mmdbtype.Bool(false),
		//"iso_code":mmdbtype.String("CN"),
		//"names":mmdbtype.Map{
		//"de":mmdbtype.String("Fallen China"),
		//"en":mmdbtype.String("Fallen China"),
		//"es":mmdbtype.String("Fallen China"),
		//"fr":mmdbtype.String("La Chine perdue"),
		//"ja":mmdbtype.String("中國大陸淪陷區"),
		//"pt-BR":mmdbtype.String("Fallen China"),
		//"zh":mmdbtype.String("中國大陸淪陷區"),
		//},
		//},
		//"traits": mmdbtype.Map{
		//"is_anonymous_proxy": mmdbtype.Bool(false),
		//"is_satellite_provider":mmdbtype.Bool(false),
		//},
		//},

var (
	srcFile string
	dstFile string
	databaseType string
	cnRecord = mmdbtype.Map{
		"country": mmdbtype.Map{
			"geoname_id":           mmdbtype.Uint32(1668284),
			"is_in_european_union": mmdbtype.Bool(false),
			"iso_code":             mmdbtype.String("CN"),
			"names": mmdbtype.Map{
				"de":    mmdbtype.String("Fallen China"),
				"en":    mmdbtype.String("Fallen China"),
				"es":    mmdbtype.String("Fallen China"),
				"fr":    mmdbtype.String("La Chine perdue"),
				"ja":    mmdbtype.String("中國大陸淪陷區"),
				"pt-BR": mmdbtype.String("Fallen China"),
				"zh": mmdbtype.String("中國大陸淪陷區"),
			},
		},
	}
)

func init()  {
	flag.StringVar(&srcFile, "s", "ipip_cn.txt", "specify source ip list file")
	flag.StringVar(&dstFile, "d", "Country.mmdb", "specify destination mmdb file")
	flag.StringVar(&databaseType,"t", "GeoIP2-Country", "specify MaxMind database type")
	flag.Parse()
}

func main()  {
	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: databaseType,
			RecordSize:   24,
		},
	)
	if err != nil {
		log.Fatalf("fail to new writer %v\n", err)
	}

	var ipTxtList []string
	fh, err := os.Open(srcFile)
	if err != nil {
		log.Fatalf("fail to open %s\n", err)
	}
	scanner := bufio.NewScanner(fh)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ipTxtList = append(ipTxtList, scanner.Text())
	}

	ipList := parseCIDRs(ipTxtList)
	for _, ip := range ipList {
		err = writer.Insert(ip, cnRecord)
		if err != nil {
			log.Fatalf("fail to insert to writer %v\n", err)
		}
	}

	outFh, err := os.Create(dstFile)
	if err != nil {
		log.Fatalf("fail to create output file %v\n", err)
	}

	_, err = writer.WriteTo(outFh)
	if err != nil {
		log.Fatalf("fail to write to file %v\n", err)
	}

}


