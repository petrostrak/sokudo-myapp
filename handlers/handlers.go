package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"myapp/data"

	"github.com/CloudyKit/jet/v6"
	"github.com/petrostrak/sokudo"
	"github.com/petrostrak/sokudo/filesystems"
	"github.com/petrostrak/sokudo/filesystems/miniofilesystem"
)

// Handlers is the type for handlers, and gives access to Celeritas and models
type Handlers struct {
	App    *sokudo.Sokudo
	Models data.Models
}

// Home is the handler to render the home page
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "home", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) ListFS(w http.ResponseWriter, r *http.Request) {
	var (
		fs   filesystems.FS
		list []filesystems.Listing
	)

	fsType := ""
	if r.URL.Query().Get("fs-type") != "" {
		fsType = r.URL.Query().Get("fs-type")
	}

	curPath := "/"
	if r.URL.Query().Get("curPath") != "" {
		curPath = r.URL.Query().Get("curPath")
		curPath, _ = url.QueryUnescape(curPath)
	}

	if fsType != "" {
		switch fsType {
		case "MINIO":
			f := h.App.FileSystems["MINIO"].(miniofilesystem.Minio)
			fs = &f
			fsType = "MINIO"
		}

		l, err := fs.List(curPath)
		if err != nil {
			h.App.ErrorLog.Println(err)
			return
		}

		list = l
	}

	vars := make(jet.VarMap)
	vars.Set("list", list)
	vars.Set("fs_type", fsType)
	vars.Set("curPath", curPath)

	err := h.render(w, r, "list-fs", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) UploadToFS(w http.ResponseWriter, r *http.Request) {
	fsType := r.URL.Query().Get("type")

	vars := make(jet.VarMap)
	vars.Set("fs_type", fsType)

	err := h.render(w, r, "upload", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) PostUploadToFS(w http.ResponseWriter, r *http.Request) {
	fileName, err := getFileToUpload(r, "formFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uploadType := r.Form.Get("upload-type")

	switch uploadType {
	case "MINIO":
		fs := h.App.FileSystems["MINIO"].(miniofilesystem.Minio)
		err = fs.Put(fileName, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	h.App.Session.Put(r.Context(), "flash", "File uploaded!")
	http.Redirect(w, r, "/files/upload?type="+uploadType, http.StatusSeeOther)
}

func getFileToUpload(r *http.Request, fieldName string) (string, error) {
	err := r.ParseMultipartForm(10 << 22)
	if err != nil {
		return "", err
	}

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	dst, err := os.Create(fmt.Sprintf("./tmp/%s", header.Filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("./tmp/%s", header.Filename), nil
}
