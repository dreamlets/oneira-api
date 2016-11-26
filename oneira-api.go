package main

import (
    "fmt"
    "io/ioutil"
    "os/exec"
    "github.com/minio/minio-go"
    "net/http"
    "encoding/json"
    "strings"
    "strconv"
    "log"
)

type FileStruct struct{
    Generation string
}

type urlStruct struct{
    url string
}

func deleteFiles(){
    del := exec.Command("bash", "-c", "cd ../../../../../lua/oneira_generator/generations/; rm -rf *.png")
    if _, err := del.Output(); err != nil {
        log.Fatalln(err)
    }
}

func generate(w http.ResponseWriter, r *http.Request){
    //clear generated files once we are done
    defer deleteFiles()
    //read JSON body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Println(err.Error())
    }
    //decode JSON and snag the size parameter
    var res map[string]interface{}
    json.Unmarshal([]byte(body),&res)
    size, err := strconv.Atoi(res["size"].(string))
    if err != nil {
        log.Fatalln("size cannot be converted to number")
    }
    //set headeer type and send back 200 before generating images
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    generated := exec.Command("bash", "-c", "th ../../../../../lua/oneira_generator/main.lua -m ../../../../../lua/oneira_generator/CPU_prod.t7 -o ../../../../../lua/oneira_generator/generations/ -s " + fmt.Sprint(size))
    //run command, then send the generated image to client via HTML
    _, err = generated.Output() 
    if err != nil {
        log.Fatalln(w, err)
    } 
    //connect to our s3 instance via Minio
    var urls []string 
    ssl := true
    s3Client, err := minio.New("s3.amazonaws.com", "<AWS KEY>", "<AWS SECRET>", ssl)
    if err != nil {
        log.Fatalln(err) 
    }
    //get list of previously generated files
    doneCh := make(chan struct{})
    defer close(doneCh)
    isRecursive := true
    objectCh := s3Client.ListObjects("oneira-project-generations", "generated", isRecursive, doneCh)
    //keep count of our previously generated files
    //TODO: don't do it this way
    total_count := 1 
    for range objectCh {
        total_count = total_count + 1
    }
    //get files and send them to our s3 instance
    files, err := ioutil.ReadDir("../../../../../lua/oneira_generator/generations")
    if err != nil {
        log.Fatalln(err) 
    }
    for i, file := range files {
        objectName := fmt.Sprintf("generated-%d.png", i + total_count)
        filePath := "../../../../../lua/oneira_generator/generations/" + file.Name()
        contentType := "image/png"
        _, err := s3Client.FPutObject("oneira-project-generations", objectName, filePath, contentType)
        if err != nil {
            log.Fatalln(err)
        }
        url := "s3-us-west-2.amazonaws.com/oneira-project-generations/" + objectName
        urls = append(urls, url)
    }
    json.NewEncoder(w).Encode(urls) 
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
    if err := http.ListenAndServe(":9090", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
