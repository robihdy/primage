package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kyokomi/cloudinary"
	"github.com/robihid/primage/primitive"
	"golang.org/x/net/context"
)

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/upload", upload)
	fs := http.FileServer(http.Dir("./images/"))
	http.Handle("/images/", http.StripPrefix("/images/", fs))
	server.ListenAndServe()
}

func upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	shape := r.FormValue("shape")
	numShapes := r.FormValue("numShapes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)[1:]

	n, _ := strconv.Atoi(numShapes)
	a, err := genImage(file, ext, n, shape)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(a)
	fmt.Println(a)

	ctx := context.Background()
	ctx = cloudinary.NewContext(ctx, "cloudinary://286849985535788:gGbYhBPmIEwqoPi4uPOLLPguik4@robihid")

	data, _ := ioutil.ReadFile(a)
	cloudinary.UploadStaticImage(ctx, a, bytes.NewBuffer(data))

	res, _ := json.Marshal(a)

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func genImage(img io.Reader, ext string, n int, m string) (string, error) {
	out, err := primitive.Transform(img, ext, n, m)
	if err != nil {
		return "", err
	}

	o, err := tempFile(ext)
	if err != nil {
		return "", err
	}
	defer o.Close()

	io.Copy(o, out)
	return o.Name(), nil
}

func tempFile(ext string) (*os.File, error) {
	tmp, err := ioutil.TempFile("./images", "img_")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())
	return os.Create(fmt.Sprintf("%s.%s", tmp.Name(), ext))
}
