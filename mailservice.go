package main

import (
   "os"
   "github.com/asccclass/staticfileserver"
)

<<<<<<< HEAD:server.go
=======
const (
   SecretKey = "Welcome to Sinica ITs@2018"
)

var (
   mail *SherryMail
)

type Token struct {
   Token	string	`json:"token"`
}

type UserCredentials struct {
   Username	string	`json:"username"`
   Password	string	`json:"password"`
}

type UserInfo struct {
   ID		int	`json:"id"`
   Name		string	`json:"name"`
   Username 	string	`json:"username"`
   Password	string	`json:"password"`
}

type Exception struct {
   Message string `json:"message"`
   Status  string `json:"status"`
}

type ExceptionError struct {
   ErrMsg string `json:"errMsg"`
}

// 資料庫連線設定
type DBConnect struct {
    DbServer string
    DbName   string
    DbLogin  string
    DbPasswd string
}

var dbconnect DBConnect

// 輸出 Json
func JsonResponse(s interface{}, w http.ResponseWriter) {
   json, err := json.Marshal(s)
   if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
   }
   w.Header().Set("Content-Type", "application/json")
   w.Write(json)
}

// 建立JWT
func (user *UserCredentials)CreateJWT(secretKey string)(Token, error) {
   token := jwt.New(jwt.SigningMethodHS256)
   claims := make(jwt.MapClaims)
   claims["username"] = user.Username
   claims["password"] = user.Password
   claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
   claims["iat"] = time.Now().Unix()
   token.Claims = claims

   tokenString, err := token.SignedString([]byte(secretKey))
   if err != nil {
      return Token{}, err
   }
   return Token{tokenString}, nil
}

func chkLoginFromJSON(w http.ResponseWriter, r *http.Request) {
   if err := r.ParseForm(); err != nil {
      fmt.Fprintf(w, "%v", err)
      return
   }
   b, err := ioutil.ReadAll(r.Body)
   defer r.Body.Close()
   if err != nil {
      fmt.Fprintf(w, "Error in parse request body.%v", err)
      return
   }
   var user UserCredentials
   if err := json.Unmarshal(b, &user); err != nil {
      w.WriteHeader(http.StatusForbidden)
      fmt.Fprintf(w, "%v", b)
      fmt.Fprintf(w, "Error in Ummarshal request body.%v", err)
      return
   }
   // 檢查帳號密碼
    ip := IPAddress.GetIPAdress(r)
    dorelogin, err := dorelogin.NewDorelogin(dbconnect.DbName, dbconnect.DbLogin, dbconnect.DbPasswd, dbconnect.DbServer, ip)
    if err != nil {
        log.Printf("Connect DBMS error: %v", err)
        w.WriteHeader(http.StatusUnauthorized)
        JsonResponse(Exception{Message: err.Error()}, w)
        return
    }
    if err = dorelogin.Chklogin(user.Username, user.Password, ""); err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        JsonResponse(Exception{Message: err.Error()}, w)
        return
    }
    // 產生JWT
    response, err := user.CreateJWT(SecretKey)
    if err != nil {
        w.WriteHeader(http.StatusForbidden)
        fmt.Fprintf(w, "Error while signing the token.%v", err)
        return
    }
    w.WriteHeader(http.StatusOK)
    JsonResponse(response, w)
}

func chkLoginFromWeb(w http.ResponseWriter, r *http.Request) {
   if err := r.ParseForm(); err != nil {
      fmt.Fprintf(w, "%v", err)
      return
   }
   // Get web params
   var user UserCredentials
   user.Username = html.EscapeString(r.FormValue("username"))
   user.Password = html.EscapeString(r.FormValue("password"))
   if user.Username == "" || user.Password == "" {
      w.WriteHeader(http.StatusForbidden)
      fmt.Fprintf(w, "Wrong Username or Password.")
      return
   }

   if strings.ToLower(user.Username) != "eplusapi" || strings.ToLower(user.Password) != "pass@word" {
      w.WriteHeader(http.StatusForbidden)
      fmt.Fprintf(w, "Error in request.")
      return
   }

   response, err := user.CreateJWT(SecretKey)
   if err != nil {
      w.WriteHeader(http.StatusForbidden)
      fmt.Fprintf(w, "Error while signing the token.")
      return
   }
   JsonResponse(response, w)
}

func IsValid(next http.HandlerFunc)  http.HandlerFunc {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      authorizationHeader := r.Header.Get("authorization")
      if authorizationHeader != "" {
         bearerToken := strings.Split(authorizationHeader, " ")
         if len(bearerToken) == 2 {
            token, err := jwt.Parse(string(bearerToken[1]), func(token *jwt.Token) (interface{}, error) {
               if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                  return nil, fmt.Errorf("There was an error.")
               } 
               return []byte(SecretKey), nil
            })
            if err != nil {
               json.NewEncoder(w).Encode(Exception{Message: err.Error()})
               return
            }
            if token.Valid {
               gcontext.Set(r, "decoded", token.Claims)
               next(w, r)
            }
         } else {  // 格式錯誤
            json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token."})
         }
      } else { // 尚未登入
            json.NewEncoder(w).Encode(Exception{Message: "Authorization Token is required.Please login First."})
      } 
   })
}

// API 新增資料
func SendEmail(w http.ResponseWriter, r *http.Request) {
   w.Header().Set("Content-Type", "application/json")
   w.WriteHeader(http.StatusOK)
   body, err := ioutil.ReadAll(r.Body)
   if err != nil {
      JsonResponse(ExceptionError{ErrMsg: err.Error()}, w)
      return
   }
   cs := []CourseInfo.Course{}
   err = json.Unmarshal(body, &cs)
   if err != nil {
      JsonResponse(ExceptionError{ErrMsg: err.Error()}, w)
      return
   }

   // Get user's information
   decoded := gcontext.Get(r, "decoded")
   _, ok := decoded.(jwt.MapClaims)
   if !ok {
       JsonResponse(ExceptionError{ErrMsg:"JWT decode error when insert."}, w)
       return
   }
   // 送信
   templateData := struct {
      Name string
      URL  string
   }{
      Name: "許功蓋",
      URL:  "http://www.justdrink.com.tw",
   }
   mail.SetRequest("andyliu@sinica.edu.tw", []string{"justgps@gmail.com"}, "測試信件標題! ", "Hello, World!")
   if err := mail.Request.ParseTemplate("template/template.html", templateData); err == nil {
      ok, err := mail.SendEmail()
      if err != nil {
         fmt.Printf("Send Email error:%v\n", err)
      } else {
         fmt.Printf("Send Email result:%v\n", ok)
      }
   } else {
      fmt.Println("Send email ok.")
   }
   // json.NewEncoder(w).Encode(Exception{Message: "ok", Status: string(b)})
}

// API 取得個人資料
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
   decoded := gcontext.Get(r, "decoded")
   if user, ok := decoded.(jwt.MapClaims); ok {
      json.NewEncoder(w).Encode(user)
   } else {
      fmt.Printf("no ok: %v", ok) 
   }
}

func CheckHealth(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
}


>>>>>>> d87084c8217f722f2196626c35b64db11368de18:mailservice.go
func main() {
   port := os.Getenv("PORT")
   if port == "" {
      port = "11009"
   }
   documentRoot := os.Getenv("DocumentRoot")
   if documentRoot == "" {
      documentRoot = "www"
   }
<<<<<<< HEAD:server.go
=======
   mailServer := os.Getenv("MailServer")
   if mailServer == "" {
      log.Printf("須設定Email Server IP/DNS.")
      os.Exit(0)
   }
   mailServerPort := os.Getenv("MailServerPort")
   if mailServerPort == "" {
      mailServerPort = "25"
   }

   port := os.Getenv("PORT")
   if port == "" {
      port = "80"
   }

   mail = NewSherryMail(account, password, mailServer,mailServerPort)
   router := mux.NewRouter()
   // step1 取得Token
   router.HandleFunc("/webauthenticate", chkLoginFromWeb)
   router.HandleFunc("/dorelogin", chkLoginFromJSON)
   // step2 執行API 
   router.HandleFunc("/sendmail", IsValid(SendEmail)).Methods("POST")
   // health check
   router.HandleFunc("/chkhealth", CheckHealth)

   // other
   // router.PathPrefix("/www/").Handler(http.StripPrefix("/www/", http.FileServer(http.Dir("./www/"))))
>>>>>>> d87084c8217f722f2196626c35b64db11368de18:mailservice.go

   server, err := SherryServer.NewServer(":" + port, documentRoot)
   if err != nil {
      panic(err)
   }
   server.Server.Handler = NewRouter(server, documentRoot
   server.Start()
}
