package main

import (
   "bytes"
   "fmt"
   "io/ioutil"
   "html/template"
   "net/smtp"
   "net/textproto"
)

//Request struct
type Request struct {
   Typez   	string
   Headers	textproto.MIMEHeader
   ReplyTo	[]string
   from		string
   to		[]string
   Cc		[]string
   Bcc		[]string
   SubjectText	[]byte		// Plaintext message (optional)
   subject	string
   body		string
}

// 設定信件格式
func (r *Request)SetMailType(mailtype string) {
   r.Typez = mailtype;
}

// 處理信件內容
func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
   t, err := template.ParseFiles(templateFileName)
   if err != nil {
      return err
   }
   buf := new(bytes.Buffer)
   if err = t.Execute(buf, data); err != nil {
      return err
   }
   r.body = buf.String()
   return nil
}

func NewRequest(from string, to []string, subject, body string) *Request {
   return &Request{
      Typez:	"text/html",
      Headers:	textproto.MIMEHeader{},
      from:	from,
      to:	to,
      subject:	subject,
      body:	body,
   }
}

type SherryMail struct {
   Auth			smtp.Auth
   MailServer		string
   MailServerPort	string
   Request		*Request
}

// 取得Email相關訊息，透過emailrep.io提供的服務
func(sm *SherryMail) CheckEmailValid(email string)(interface{}, error) {
   if email == "" {
      return nil, fmt.Errorf("No Email.")
   }
   resp, err := http.Get("emailrep.io/" + email)
   if err != nil {
      retrun nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode/100 == 4 || resp.StatusCode/100 == 5 {
      return nil, fmt.Errorf("the response was returned with a %d", res.StatusCode)
   }
   data, _ := ioutil.ReadAll(resp.Body)
   return data, nil
}

func(sm *SherryMail) SetRequest(from string, tos []string, subject, body string)(error) {
   sm.Request = NewRequest(from, tos, subject, body)
   return nil
}

func(sm *SherryMail) SendEmail() (bool, error) {
   mime := "MIME-version: 1.0;\nContent-Type: " + sm.Request.Typez + "; charset=\"UTF-8\";\n\n"
   subject := "Subject: " + sm.Request.subject + "!\n"
   msg := []byte(subject + mime + "\n" + sm.Request.body)

   if err := smtp.SendMail(sm.MailServer, sm.Auth, sm.Request.from, sm.Request.to, msg); err != nil {
      return false, fmt.Errorf("Send email error: %v", err)
   }
   return true, nil
}

func NewSherryMail(account, password, mailserver,port string)(*SherryMail) {
   auth := smtp.PlainAuth("", account, password, mailserver)
   return &SherryMail {
      Auth: auth,
      MailServer: mailserver + ":" + port, 
   }
}

/*
func main() {
   account := os.Getenv("MailAccount")
   if  account == "" {
      log.Printf("須設定Email Server登入帳號")
      os.Exit(0)
   }
   password := os.Getenv("MailPassword")
   if password == "" {
      log.Printf("須設定Email Server登入密碼")
      os.Exit(0)
   }
   mailServer := os.Getenv("MailServer")
   if mailServer == "" {
      log.Printf("須設定Email Server IP/DNS.")
      os.Exit(0)
   }
   mailServerPort := os.Getenv("MailServerPort")
   if mailServerPort == "" {
      log.Printf("須設定Email Server PORT.")
      os.Exit(0)
   }

   sm := NewSherryMail(account, password, mailServer, mailServerPort)
   templateData := struct {
      Name string
      URL  string
   }{
      Name: "許功蓋",
      URL:  "http://www.justdrink.com.tw",
   }
   sm.SetRequest("andyliu@sinica.edu.tw", []string{"justgps@gmail.com"}, "測試信件標題! ", "Hello, World!")
   if err := sm.Request.ParseTemplate("template/template.html", templateData); err == nil {
      ok, err := sm.SendEmail()
      if err != nil {
         fmt.Printf("Send Email error:%v\n", err)
      } else {
         fmt.Printf("Send Email result:%v\n", ok)
      }
   } else {
      fmt.Println("Send email ok.")
   }
}
*/
