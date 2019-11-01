package main

import (
   "os"
   "github.com/asccclass/staticfileserver"
)

func main() {
   port := os.Getenv("PORT")
   if port == "" {
      port = "11009"
   }
   documentRoot := os.Getenv("DocumentRoot")
   if documentRoot == "" {
      documentRoot = "www"
   }

   server, err := SherryServer.NewServer(":" + port, documentRoot)
   if err != nil {
      panic(err)
   }
   server.Server.Handler = NewRouter(server, documentRoot
   server.Start()
}
