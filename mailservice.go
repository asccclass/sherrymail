package main

import (
   "fmt"
   "os"
   "github.com/asccclass/staticfileserver"
)

func main() {
   port := os.Getenv("PORT")
   if port == "" {
      port = "80"
   }

   documentRoot := os.Getenv("DocumentRoot")
   if documentRoot == "" {
      documentRoot = "www"
   }
   if os.Getenv("MailServer") == "" {
      fmt.Printf("須設定Email Server IP/DNS.")
      os.Exit(0)
   }
   templateRoot := os.Getenv("TemplateRoot")
   if templateRoot == "" {
      templateRoot = "www/template"
   }

   server, err := SherryServer.NewServer(":" + port, documentRoot, templateRoot)
   if err != nil {
      panic(err)
   }
   server.Server.Handler = NewRouter(server, documentRoot)
   server.Start()
}
