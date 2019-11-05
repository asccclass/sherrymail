package sherrymail

import(
   "fmt"
   "bytes"
   "io/ioutil"
   "net/http"
   "encoding/json"
   "net/textproto"
   "html/template"
)

// WebPrint
func(s *SherryMail) WebPrint(w http.ResponseWriter, message string) {
   w.Header().Set("Content-Type", "application/json;charset=UTF-8")
   w.WriteHeader(http.StatusOK)
   fmt.Fprintf(w, message)
}

// 輸出web error
func(s *SherryMail) Error2Web(w http.ResponseWriter, err error) {
   w.Header().Set("Content-Type", "application/json;charset=UTF-8")
   w.WriteHeader(http.StatusOK)
   fmt.Fprintf(w, "{\"errMsg\": \"%s(server)\"}", err.Error())
}

// 處理信件樣板
func(sm *SherryMail) ParseTemplateFromWeb(sf SendInfos)(string, error) {
   if sf.Template == "" {
      return  "", fmt.Errorf("no template message")
   }
   t, err := template.ParseFiles(sf.Template)
   if err != nil {
      return "", err
   }
   buf := new(bytes.Buffer)
   if err = t.Execute(buf, sf); err != nil {
      return "", err
   }
   return buf.String(), nil
}

// SendEmailFromWeb 透過web介面送信
func (s *SherryMail) SendEmailFromWeb(w http.ResponseWriter, r *http.Request) {
   body, err := ioutil.ReadAll(r.Body)
   if err != nil {
      s.Error2Web(w, err) 
      return
   }
   cs := SendInfos{}
   err = json.Unmarshal(body, &cs)
   if err != nil {
      s.Error2Web(w, err)
      return
   }
   cs.Headers = textproto.MIMEHeader{}
   if cs.Typez == "" {
      cs.Typez = "text/html"
   }

   if cs.Template != "" {
      cs.Content, err = s.ParseTemplateFromWeb(cs)
      if err != nil {
         s.Error2Web(w, err)
         return
      }
   }
   c, err := s.SendEmail(&cs)
   if err != nil {
      s.Error2Web(w, err)
      return
   }
   b, err := json.Marshal(&c)
   if err != nil {
      s.Error2Web(w, err)
      return
   }
   s.WebPrint(w, string(b))
}
