package main

import (
  "fmt"
  "net/http"
  "strings"
  "log"
  "os/exec"
)

func test(w http.ResponseWriter, r *http.Request){
  testTorch := exec.Command("th", "test_torch.lua")
  if output, err := testTorch.Output(); err != nil{
    print(err)
  } else {
   fmt.Fprintf(w, string(output))
  }
}

func generate(w http.ResponseWriter, r *http.Request){
  //code to generate images and shit
}

func sayHelloFoucault(w http.ResponseWriter, r *http.Request){
  r.ParseForm()
  fmt.Println(r.Form)
  fmt.Println("path", r.URL.Path)
  fmt.Println("scheme", r.URL.Scheme)
  fmt.Println(r.Form["url_long"])
  for k, v := range r.Form{
    fmt.Println("key:", k)
    fmt.Println("val:", strings.Join(v, ""))
  }
  fmt.Fprintf(w, "all is walking or dreams, truth or error, the light of being or the nothingness of shadow")
}

func main(){
  http.HandleFunc("/", sayHelloFoucault)
  http.HandleFunc("/generate", generate)
  http.HandleFunc("/test", test)
  if err := http.ListenAndServe(":9090", nil); err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
