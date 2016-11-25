package main

import (
  "fmt"
  "encoding/base64"
  "html/template"
  "bytes"
  "os"
  "os/exec"
  "net/http"
  "image/png"
  "strings"
  "log"
)

type FileStruct struct{
  Generation string
}

func test(w http.ResponseWriter, r *http.Request){
  testTorch := exec.Command("th", "test_torch.lua")
  if output, err := testTorch.Output(); err != nil {
    print(err)
  } else {
   fmt.Fprintf(w, string(output))
  }
}

func generate(w http.ResponseWriter, r *http.Request){
    //command for Torch model to generate a single image
    generated := exec.Command("bash", "-c", "th $GOPATH/src/github.com/dreamlets/oneira_art/main.lua -m ~/dcgan_vae_torch/checkpoints/34000_net_G.t7 -o /home/ubuntu/golang/src/github.com/dreamlets/oneira_art/generations/")

    //run command, then send the generated image to client via HTML
    if _, err := generated.Output(); err != nil {
        fmt.Fprint(w, "unable to open file.")
        fmt.Fprint(w, err)
    } else {
        var ImageTemplate string = `<!DOCTYPE html>
          <html lang="en"><head></head>
          <body><img src="data:image/jpg;base64,{{.Generation}}"></body>`
        file, err := os.Open("../oneira_art/generations/generation.png")
        if err != nil {
            log.Println("unable to open file")
        }
        defer file.Close()

        buffer := new(bytes.Buffer)
        img, err := png.Decode(file)
        if err != nil {
            log.Println("unable to encode image.")
        }
        if err := png.Encode(buffer, img); err != nil {
            log.Fatalln("unable to encode image")
        }
        str := base64.StdEncoding.EncodeToString(buffer.Bytes())
        if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil{
          log.Println("unable to parse image template")
        } else {
          data := FileStruct{str}
          if err = tmpl.Execute(w, data); err != nil{
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
