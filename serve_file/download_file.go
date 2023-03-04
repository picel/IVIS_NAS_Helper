package serve_file

import (
	"fmt"
	"net/http"
	token_process "ivis_nas/login/token_process"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in
	_, err := token_process.VerifyToken(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fmt.Println("DownloadFile")
	w.Header().Set("Content-Disposition", "attachment; filename="+r.URL.Path[14:])
	http.ServeFile(w, r, "/mnt/ivis_nas"+r.URL.Path[14:])
}