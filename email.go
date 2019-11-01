package customer 
/*
   客戶資料管理
*/

import(
   "fmt"
   "net/http"
   "io/ioutil"
   "strconv"
   // "reflect"
   "encoding/json"
   "database/sql"
   "github.com/gorilla/mux"
   "github.com/asccclass/sherrytime"
   "github.com/asccclass/sherrydb/mysql"
   "github.com/asccclass/sherryschema"
)

// 輸出web error
func Error2Web(w http.ResponseWriter, err error) {
   w.Header().Set("Content-Type", "application/json;charset=UTF-8")
   w.WriteHeader(http.StatusOK)
   fmt.Fprintf(w, "{\"errMsg\": \"%s(server)\"}", err.Error())
}

func NewSherryMail()(*SherryMail, error) {
}

func(c *Customer) AddRouter(router *mux.Router) {
/*
   router.HandleFunc("/customer/{customerID}", c.Read).Methods("GET")   	//取得單筆資料內容
   router.HandleFunc("/customer", c.AddCustomers).Methods("POST")		//新增資料
   router.HandleFunc("/customer/{customerID}", c.Update).Methods("PUT") 	//修改單筆
   router.HandleFunc("/customer/{customerID}", c.Delete).Methods("DELETE") 	//刪除單筆
*/
}

