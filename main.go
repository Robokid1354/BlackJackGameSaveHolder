package main

import (
    "log"
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
  "io"
  "path/filepath"
)

type Data struct {
    //SteamID string `json:"SteamID"`
    Name string  `json:"Name"`
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.Error(w, "404 not found.", http.StatusNotFound)
    return
  }

  switch r.Method {
    case "GET":
      queryDataSteamID, ok := r.URL.Query()["requestDataSteamID"]

      if !ok || len(queryDataSteamID[0]) < 1 /*|| len(queryDataName[0]) < 1*/{
        queryDataSpecial, ok := r.URL.Query()["requestDataSpecial"]
        if !ok || len(queryDataSpecial[0]) < 1 {
          log.Println("Url Param 'key' is missing")
          http.ServeFile(w,r, "form.html")
          return
        }
        requestDataSpecial := queryDataSpecial[0]
        if requestDataSpecial == "all" {
          var out string
          var files []string

          root := "Saves"
          err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
            files = append(files, info.Name())
            return nil
          })
          if err != nil {
            log.Fatal(err)
          }

          for _, file := range files {
            if file != "File_To_Delete.md" && file !="Saves" {out = out + file}
          }

          fmt.Fprintln(w, out)
          log.Println(out)

          return
        }
      }

      // Query()["key"] will return an array of items,
      // we only want the single item.
	    requestDataSteamID := queryDataSteamID[0]

      log.Println("The Key Thing, lets output that real fast: " + string(requestDataSteamID))

      _, err := os.Stat("Saves\\"+requestDataSteamID+".json")

      if err != nil {
        fmt.Fprintln(w, "nil")
        log.Println("Save not Found! :: " + string(requestDataSteamID))
        return
      }

      file, err := ioutil.ReadFile("Saves\\"+requestDataSteamID+".json")

      if err != nil {
        log.Fatal(err)
      }

      fmt.Fprintln(w, string(file))
    case "POST":
      // Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
      if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
      }
      fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
      steamID := r.FormValue("steamID")
      json := r.FormValue("json")
      WriteToFile("Saves\\"+steamID+".json",json)
    default:
      fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
  }
}



// GO HERE     http://localhost:8080/?requestDataSteamID=6345642654


func WriteToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
  }
