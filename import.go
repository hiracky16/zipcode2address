package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "path"
  "encoding/csv"
  "golang.org/x/text/transform"
  "golang.org/x/text/encoding/japanese"
  "strings"

	"github.com/artdarek/go-unzip"
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/guregu/dynamo"
)

type AddressData struct {
  Zipcode string `dynamo:"zipcode"`
  Address string
  IsOffce bool
}

func download() string {
  var url = os.Getenv("CSV_SOURCE")

  // ファイルをダウンロード
  fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
    fmt.Println(err)
    os.Exit(1)
	}
	fmt.Println("status:", response.Status)

	body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  // file に書き込み
  _, filename := path.Split(url)
  filepath := "/tmp/" + filename
  file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0777)
  if err != nil {
		fmt.Println(err)
		os.Exit(1)
  }
  defer file.Close()
  file.Write(body)
	return filepath
}

func defrost(file string) string {
  // s3 への接続情報
  var bucket = os.Getenv("BUCKET")
  svc := s3.New(session.New(), &aws.Config{
    Region: aws.String(endpoints.ApNortheast1RegionID),
  })
  uz := unzip.New(file, "/tmp/")
  unzipErr := uz.Extract()
  if unzipErr != nil {
    fmt.Println(unzipErr)
    os.Exit(1)
  }
  // FIXME: ファイル名
  filename := "zenkoku.csv"
  f, err := os.Open("/tmp/" + filename)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  defer f.Close()
  _, err = svc.PutObject(&s3.PutObjectInput{
    Body: f,
    Bucket: aws.String(bucket),
    Key: aws.String(filename),
    ACL: aws.String("private"),
    ServerSideEncryption: aws.String("AES256"),
  })
  return "/tmp/" + filename
}

func parse(file string) {
  db := dynamo.New(session.New(), &aws.Config{
    Region: aws.String(endpoints.ApNortheast1RegionID),
  })
  var tableName = os.Getenv("TABLE")
  table := db.Table(tableName)

  f, err := os.Open(file)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  defer f.Close()
  reader := csv.NewReader(f)
  var line []string

  // NOTE: 行数を絞っている
  for i := 0; i < 100; i++ {
    line, err = reader.Read()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
    // NOTE: 最初はヘッダーなので飛ばす
    if i == 0 {
      continue
    }
    // put item
    isOffce := line[5] == "1"
    var address string
    if isOffce {
      address = fmt.Sprintf("%s%s%s", line[7], line[9], line[20])
    } else {
      address = fmt.Sprintf("%s%s%s%s", line[7], line[9], line[11], line[15])
    }
    address = tranform(address)
    zipcode := tranform(line[4])
    zipcode = strings.Replace(zipcode, "-", "", -1)
    w := AddressData{Zipcode: zipcode, Address: address, IsOffce: isOffce}
    fmt.Println(w)
    putErr := table.Put(w).Run()
    if putErr != nil {
      fmt.Println(putErr)
      os.Exit(1)
    }
  }
}

func tranform(str string) string {
  utf, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), []byte(str))
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  return string(utf)
}

func Handler() {
	file := download()
  fmt.Println(fmt.Sprintf("download %s", file))
  csv := defrost(file)
  parse(csv)
}

func main() {
  lambda.Start(Handler)
}