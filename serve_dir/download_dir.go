package serve_dir

import (
	"fmt"
	"net/http"
	"os"
	"archive/zip"
	"io"
	"io/ioutil"
	"time"
	token_process "ivis_nas/login/token_process"
)

func DownloadDir(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in
	_, err := token_process.VerifyToken(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	date := time.Now().Format("2006-01-02")

	dirPath := "/mnt/ivis_nas" + r.URL.Path[13:] + "/"
	outFile, err := os.Create("/tmp/" + date + ".zip")
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)

	// Add files in the directory recursively to the zip file
	addFilesToZip(zipWriter, dirPath, "")

	// Don't forget to check the error on Close.
	err = zipWriter.Close()
	if err != nil {
		fmt.Println(err)
	}

	// Download the zip file
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename="+date+".zip")
	http.ServeFile(w, r, "/tmp/"+date+".zip")
}

func addFilesToZip(zipWriter *zip.Writer, basePath string, baseInZip string) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range files {
		if file.Name()[0] != '.' {
			if file.IsDir() {
				addFilesToZip(zipWriter, basePath+file.Name()+"/", baseInZip+file.Name()+"/")
			} else {
				filePath := basePath + file.Name()
				fmt.Println(filePath)
				fileToZip, err := os.Open(filePath)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer fileToZip.Close()

				info, err := fileToZip.Stat()
				if err != nil {
					fmt.Println(err)
					return
				}

				header, err := zip.FileInfoHeader(info)
				if err != nil {
					fmt.Println(err)
					return
				}
				header.Name = baseInZip + file.Name()
				header.Method = zip.Deflate

				writer, err := zipWriter.CreateHeader(header)
				if err != nil {
					fmt.Println(err)
					return
				}
				_, err = fileToZip.Seek(0, 0)
				if err != nil {
					fmt.Println(err)
					return
				}
				_, err = io.Copy(writer, fileToZip)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}