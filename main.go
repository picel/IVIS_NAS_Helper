package main

import (
	"net/http"

	files "ivis_nas/files"
	login "ivis_nas/login"
	token_process "ivis_nas/login/token_process"
	logout "ivis_nas/logout"
	serve_dir "ivis_nas/serve_dir"
	serve_file "ivis_nas/serve_file"
)

func initPage(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in
	_, err := token_process.VerifyToken(r)
	if err != nil {
		// show index.html
		http.ServeFile(w, r, "index.html")
	} else {
		// redirect to /files
		http.Redirect(w, r, "/files", http.StatusSeeOther)
	}
}

func main() {
	http.HandleFunc("/", initPage)
	http.HandleFunc("/files/", files.FilesHandler)
	http.HandleFunc("/serve_file/", serve_file.FileServer)
	http.HandleFunc("/download_file/", serve_file.DownloadFile)
	http.HandleFunc("/download_dir/", serve_dir.DownloadDir)
	http.HandleFunc("/login/", login.Login)
	http.HandleFunc("/remote_login/", login.RemoteLogin)
	http.HandleFunc("/remote_token_verify/", login.RemoteTokenVerify)
	http.HandleFunc("/logout/", logout.Logout)
	http.ListenAndServe(":2222", nil)
}
