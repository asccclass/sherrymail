package sherrymail

import (
   "fmt"
   "bytes"
   "strings"
   "io/ioutil"
   "html/template"
   "net/http"
   "net/smtp"
   "github.com/gorilla/mux"
)

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

// 取得Email相關訊息，透過emailrep.io提供的服務
func(sm *SherryMail) CheckEmailValid(email string)(interface{}, error) {
   if email == "" {
      return nil, fmt.Errorf("No Email.")
   }
   resp, err := http.Get("emailrep.io/" + email)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode/100 == 4 || resp.StatusCode/100 == 5 {
      return nil, fmt.Errorf("the response was returned with ")
   }
   data, _ := ioutil.ReadAll(resp.Body)
   return data, nil
}

func(sm *SherryMail) SetRequest(from string, tos []string, subject, body string)(error) {
   sm.Request = NewRequest(from, tos, subject, body)
   return nil
}

func(sm *SherryMail) SendEmail(req *SendInfos) (*SendInfos, error) {
   var body strings.Builder
   body.WriteString("Subject: ")
   body.WriteString(req.Subject)
   body.WriteString("\r\n")
   if req.ReplyTo.Email != "" {
      body.WriteString("Return-Path: <")
      body.WriteString(req.ReplyTo.Email)
      body.WriteString(">\r\n")
   }
   body.WriteString("MIME-version: 1.0;\nContent-Type: ")
   body.WriteString(req.Typez)
   body.WriteString("; charset=\"UTF-8\";\n\n\n")
   body.WriteString(req.Content)

   for i, e := range req.Receiver {
      to := []string{e.Email}
      if err := smtp.SendMail(sm.MailServer, sm.Auth, req.Sender.Email, to, []byte(body.String())); err != nil {
         req.Receiver[i].Result = err.Error()
      } else {
         req.Receiver[i].Result = "ok"
      }
   }
   return req, nil
}

// AddCustomerRouter
func(c *SherryMail) AddCustomerRouter(router *mux.Router) {
   router.HandleFunc("/sendmail", c.SendEmailFromWeb).Methods("POST")           // 送信
}
