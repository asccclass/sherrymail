// router.go
package main

import(
   "github.com/gorilla/mux"
   "github.com/asccclass/staticfileserver"
   "github.com/asccclass/serverstatus"
)

// Create your Router function
func NewRouter(srv *SherryServer.ShryServer, documentRoot string)(*mux.Router) {
   router := mux.NewRouter()
/*
   // Rank router
   rank, err := rank.NewRank(dbconnect)
   if err != nil {
      panic(err)
   }
   rank.AddRankRouter(router)			// 客戶等級 Router
*/

   // Static File server
   staticfileserver := SherryServer.StaticFileServer{documentRoot, "index.html"}
   router.PathPrefix("/").Handler(staticfileserver)

   return router
}
