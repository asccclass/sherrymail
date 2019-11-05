package sherrymail

import(
   "os"
   "net/textproto"
   "net/smtp"
   "github.com/gorilla/mux"
   "github.com/asccclass/staticfileserver"
)

// Request struct
type Request struct {
   Typez        string
   Headers      textproto.MIMEHeader
   ReplyTo      []string
   from         string
   to           []string
   Cc           []string
   Bcc          []string
   SubjectText  []byte          // Plaintext message (optional)
   subject      string
   body         string
}

type SherryMail struct {
   Auth                 smtp.Auth
   MailServer           string
   MailServerPort       string
   Request              *Request
   Server		*SherryServer.ShryServer
}

// 收／送信件資訊
type Receiver struct {
   Name	string	`json:"name"`
   Email	string		`json:"email"`
   Result	string		`json:"result"`
}

// 要被取代的參數
type MailParams struct {
   Key	string		`json:"name"`
   Value	string	`josn:"value"`
}

// web 送過來要寄送的封包 josn
type SendInfos struct{
   Typez	string		`json:"minetype"`
   Subject	string		`json:"subject"`
   Content	string		`json:"content"`
   Params	[]MailParams	`json:"params"`
   Receiver	[]Receiver	`json:"to"`
   Sender	Receiver	`json:"from"`
   ReplyTo	Receiver	`json:"replyto"`
   Template	string		`json:"template"`
   Headers      textproto.MIMEHeader
}

func NewRequest(from string, to []string, subject, body string) *Request {
   return &Request{
      Typez:    "text/html",
      Headers:  textproto.MIMEHeader{},
      from:     from,
      to:       to,
      subject:  subject,
      body:     body,
   }
}

func NewSherryMail(srv *SherryServer.ShryServer)(*SherryMail) {
   mailserver := os.Getenv("MailServer")
   port := os.Getenv("MailServerPort")
   if port == "" {
      port = "25"
   }

   auth := smtp.PlainAuth("", os.Getenv("MailAccount"), os.Getenv("MailPassword"), mailserver)
   return &SherryMail {
      Auth: auth,
      MailServer: mailserver + ":" + port,
      Server: srv,
   }
}

func(c *SherryMail) AddRouter(router *mux.Router) {
   router.HandleFunc("/send", c.SendEmailFromWeb).Methods("POST")   	// 送信
}
