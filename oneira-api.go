package main

import (
  "fmt"
  "bytes"
  "encoding/base64"
  "os"
  "html/template"
  "os/exec"
  "net/http"
  "image/png"
  "strings"
  "log"
)

func test(w http.ResponseWriter, r *http.Request){
  testTorch := exec.Command("th", "test_torch.lua")
  if output, err := testTorch.Output(); err != nil {
    print(err)
  } else {
   fmt.Fprintf(w, string(output))
  }
}

func generate(w http.ResponseWriter, r *http.Request){
    generated := exec.Command("th", "../oneira_art/main.lua -i ~/dcgan_vae_torch/checkpoints_for_prod/save_cpu_model.lua -o generations/" )
    if _, err := generated.Output(); err != nil {
        fmt.Fprint(w, "Error. Please try again.")
    } else {
        if img, err = os.Open("../oneira_art/generations/generation.png"); err != nil {
            log.Println("unable to open file")
        }
        defer img.Close()
        
        var ImageTemplate string = `<DOCTYPE html>
            <html lang="en"><head></head>
<body><img src="data:image/jpg;base64,{{.Image}}"></body>`
        
        buffer := new(bytes.Buffer)
        if err := png.Encode(buffer, *img); err != nil {
            log.Println("unable to encode image.")
        }
        str := base64.StdEncoding.EncodeToString(buffer.Bytes())
        if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
            log.Println("unable to parse image template")
        } else {
            data := map[string]interface{}{"Image": str}
            if err = templ.Execute(w, data); err != nil {
                log.Println("unable to execute template.")
            }
        }
    }
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
