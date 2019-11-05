// router.go
package main

import(
   "github.com/gorilla/mux"
   "github.com/asccclass/staticfileserver"
   // "github.com/asccclass/serverstatus"
   "github.com/asccclass/sherrymail/libs/sherrymail"
)

// Create your Router function
func NewRouter(srv *SherryServer.ShryServer, documentRoot string)(*mux.Router) {
   router := mux.NewRouter()

   // Rank router
   mail := sherrymail.NewSherryMail(srv)
   mail.AddRouter(router)

   //logger
   router.Use(SherryServer.ZapLogger(srv.Logger))

   // Static File server
   staticfileserver := SherryServer.StaticFileServer{documentRoot, "index.html"}
   router.PathPrefix("/").Handler(staticfileserver)

   return router
}
