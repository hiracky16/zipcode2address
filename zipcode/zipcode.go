package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "path"
  "encoding/csv"
  "github.com/artdarek/go-unzip"
)

const url = "http://jusyo.jp/downloads/new/csv/csv_zenkoku.zip";

func download() string {
  // ファイルをダウンロード
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

  // ファイルを書き込む
  _, filename := path.Split(url)
  filepath := "downloads/" + filename
  file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0666)
  if err != nil {
    fmt.Println(err)
  }
  defer file.Close()
  file.Write(body)

  return filepath
}

func defrost(file string) {
  uz := unzip.New(file, "downloads")
  err := uz.Extract()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func putItems(file string) {
  f, err := os.Open(file)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  defer f.Close()
  reader := csv.NewReader(f)
    var line []string

    for {
      line, err = reader.Read()
      if err != nil {
          break
      }
      fmt.Println(line[0])
      fmt.Println(line[1])
      os.Exit(1)
    }
}

func main() {
  file := download()
  fmt.Println(fmt.Sprintf("download %s", file))
  defrost(file)
  csv := "downloads/zenkoku.csv"
  putItems(csv)
}