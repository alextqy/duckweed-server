package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetMap(r *http.Request) map[string][]string {
	return r.Form
}

func PostMap(r *http.Request) map[string][]string {
	return r.PostForm
}

func Get(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func Post(r *http.Request, key string) string {
	return r.PostFormValue(key)
}

func FormFile(w http.ResponseWriter, r *http.Request, key string) (bool, string) {
	f, fheader, err := r.FormFile("file")
	defer func(f io.Closer) {
		if err := f.Close(); err != nil {
			fmt.Printf("defer close file err: %v", err.Error())
		}
	}(f)
	if err != nil {
		return false, err.Error()
	}

	newf, err := os.OpenFile("./Temp/upload/"+fheader.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return false, err.Error()
	}
	defer func(f io.Closer) {
		if err := f.Close(); err != nil {
			fmt.Printf("defer close file err: %v", err.Error())
		}
	}(newf)
	if err != nil {
		return false, err.Error()
	}

	_, err = io.Copy(newf, f)
	if err != nil {
		return false, err.Error()
	}

	return true, newf.Name()
}

func HttpWrite(w http.ResponseWriter, Data interface{}) (int, error) {
	j, err := json.Marshal(Data)
	if err != nil {
		return 0, err
	}
	return w.Write(j)
}
