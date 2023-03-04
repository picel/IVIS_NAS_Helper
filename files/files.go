package files

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	token_process "ivis_nas/login/token_process"
)

type PageVariables struct {
	PageTitle    string
	ParentFolder string
	PageFiles    []File
	IsRoot       bool
}

type File struct {
	Name    string
	Path    string
	Size    int64
	Unit    string
	ModTime string
	IsDir   bool
}

func FilesHandler(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in
	_, err := token_process.VerifyToken(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//show file list in /mnt/ivis_nas with html template in template folder except for hidden files
	files, _ := ioutil.ReadDir("/mnt/ivis_nas" + r.URL.Path[6:] + "/")
	var fileList []File
	for _, file := range files {
		if file.Name()[0] != '.' {
			//file size as Size
			var size int64
			var unit string
			if file.Size() < 1024 {
				size = file.Size()
				unit = "B"
			} else if file.Size() < 1048576 {
				size = file.Size() / 1024
				unit = "KB"
			} else if file.Size() < 1073741824 {
				size = file.Size() / 1048576
				unit = "MB"
			} else {
				size = file.Size() / 1073741824
				unit = "GB"
			}
			currentTime := time.Now()
			modTime := file.ModTime()
			// if file is modified in the last 24 hours, show difference in hours
			if currentTime.Sub(modTime).Hours() < 1.0 {
				fileList = append(fileList, File{file.Name(), r.URL.Path[6:] + "/" + file.Name(), size, unit, "조금 전", file.IsDir()})
			} else if currentTime.Sub(modTime).Hours() < 24 {
				fileList = append(fileList, File{file.Name(), r.URL.Path[6:] + "/" + file.Name(), size, unit, fmt.Sprintf("%.0f", currentTime.Sub(modTime).Hours()) + "시간 전", file.IsDir()})
			} else if currentTime.Sub(modTime).Hours() < 720 {
				// if file is modified in the last 7 days, show difference in days
				fileList = append(fileList, File{file.Name(), r.URL.Path[6:] + "/" + file.Name(), size, unit, fmt.Sprintf("%.0f", currentTime.Sub(modTime).Hours()/24) + "일 전", file.IsDir()})
			} else if currentTime.Sub(modTime).Hours() < 8790 {
				// if file is modified in the last 365 days, show difference in months
				fileList = append(fileList, File{file.Name(), r.URL.Path[6:] + "/" + file.Name(), size, unit, fmt.Sprintf("%.0f", currentTime.Sub(modTime).Hours()/24/30) + "개월 전", file.IsDir()})
			} else {
				fileList = append(fileList, File{file.Name(), r.URL.Path[6:] + "/" + file.Name(), size, unit, "오래 전", file.IsDir()})
			}
		}
	}
	var parentFolder string
	if r.URL.Path[6:] == "" {
		parentFolder = ""
	} else {
		// split path by / and remove last element
		parentFolder = r.URL.Path[6:]
		parentFolder = parentFolder[:len(parentFolder)-len(parentFolder[strings.LastIndex(parentFolder, "/"):])]
	}
	HomePageVars := PageVariables{
		PageTitle:    r.URL.Path[6:],
		ParentFolder: parentFolder,
		PageFiles:    fileList,
		IsRoot:       !(r.URL.Path[6:] == "/"),
	}
	t, _ := template.ParseFiles("files/template.html")
	t.Execute(w, HomePageVars)
}