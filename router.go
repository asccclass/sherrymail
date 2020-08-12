// router.go
package main

import(
   "os"
   "github.com/gorilla/mux"
   "github.com/asccclass/staticfileserver"
   "github.com/asccclass/serverstatus"
   "github.com/asccclass/sherrymail/libs/sherrymail"
   "github.com/asccclass/staticfileserver/libs/googlesheet"
)

// Create your Router function
func NewRouter(srv *SherryServer.ShryServer, documentRoot string)(*mux.Router) {
   router := mux.NewRouter()

   // Rank router
   mail := sherrymail.NewSherryMail(srv)
   mail.AddRouter(router)

   // GS Reader
   gs, _ := SherryGoogleSheet.NewSryGoogleSheet(srv)
   gs.AddRouter(router)

   //logger
   systemName := os.Getenv("SystemName")
   m := serverstatus.NewServerStatus(systemName)
   router.HandleFunc("/healthz", m.Healthz).Methods("GET")

   // Static File server
   staticfileserver := SherryServer.StaticFileServer{documentRoot, "index.html"}
   router.PathPrefix("/").Handler(staticfileserver)

   return router
}
