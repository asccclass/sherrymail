## Sherrymail 寄信公用程式


### Install
* get library

```
go get github.com/asccclass/sherrymail
```

* create envfile (vi envfile)

```
SystemName=mail公用程式
MailServer=smtp.myemail.server
MailServerPort=25
MailAccount=myaccount
MailPassword=mypassword
DocumentRoot=www/html
PORT=80

TemplateRoot=www/template
```

### API 發送範例
* 網址：https://myserver/email/send

* 資料格式 json

```
[
   {
      "minetype":"text/html",
      "subject":"2020YoungerBoss 活動訊息通知",
      "content":"丁OO您好<br><br>國立臺灣戲曲學院 燈光音響體驗課程 12:00(已錄取)",
      "from":{
         "Name":"YoungerBoss",
         "email":"andyliu@com.tw"
      },
      "to":[
         {
            "Name":"丁OO",
            "email":"ps@mail.com"
         }
      ],
      "replyto":{
         "Name":"Andy",
         "Email":"ps@gmail.com"
      }
   },
   {
      "minetype":"text/html",
      "subject":"2020YoungerBoss 活動訊息通知",
      "content":"曾OO 您好<br><br>國立臺灣戲曲學院 燈光音響體驗課程 14:30(已錄取)",
      "from":{
         "Name":"YoungerBoss",
         "email":"andyliu@com.tw"
      },
      "to":[
         {
            "Name":"曾OO",
            "email":"ju@gmail.com"
         }
      ],
      "replyto":{
         "Name":"Andy",
         "Email":"just@gmail.com"
      }
   }
]
```

* Go 程式碼

```
package main

import (
  "fmt"
  "strings"
  "net/http"
  "io/ioutil"
)

func main() {

  url := "https://ascare.sinica.edu.tw/email/send"
  method := "POST"

  payload := strings.NewReader("[\n   {\n      \"minetype\":\"text/html\",\n      \"subject\":\"2020YoungerBoss 活動訊息通知\",\n      \"content\":\"丁世軒您好<br><br>國立臺灣戲曲學院 燈光音響體驗課程 12:00(已錄取)\",\n      \"from\":{\n         \"Name\":\"YoungerBoss\",\n         \"email\":\"andyliu@sinica.edu.tw\"\n      },\n      \"to\":[\n         {\n            \"Name\":\"丁世軒\",\n            \"email\":\"justgps@gmail.com\"\n         }\n      ],\n      \"replyto\":{\n         \"Name\":\"Andy\",\n         \"Email\":\"justgps@gmail.com\"\n      }\n   },\n   {\n      \"minetype\":\"text/html\",\n      \"subject\":\"2020YoungerBoss 活動訊息通知\",\n      \"content\":\"曾楚晴 您好<br><br>國立臺灣戲曲學院 燈光音響體驗課程 14:30(已錄取)\",\n      \"from\":{\n         \"Name\":\"YoungerBoss\",\n         \"email\":\"andyliu@sinica.edu.tw\"\n      },\n      \"to\":[\n         {\n            \"Name\":\"曾楚晴\",\n            \"email\":\"justgps@gmail.com\"\n         }\n      ],\n      \"replyto\":{\n         \"Name\":\"Andy\",\n         \"Email\":\"justgps@gmail.com\"\n      }\n   }\n]")

  client := &http.Client { }
  req, err := http.NewRequest(method, url, payload)

  if err != nil {
    fmt.Println(err)
  }
  req.Header.Add("Content-Type", "application/json")

  res, err := client.Do(req)
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)

  fmt.Println(string(body))
}
```

### TODP
* libs/smtp.go 尚未處理完成


### 參考文件
* [將建好的 Image 發佈到 Docker Hub](https://justjii.justdrink.com.tw/%e5%b0%87%e5%bb%ba%e5%a5%bd%e7%9a%84-image-%e7%99%bc%e4%bd%88%e5%88%b0-docker-hub/)
* https://github.com/jordan-wright/email
* [notify](https://github.com/nikoksr/notify)
