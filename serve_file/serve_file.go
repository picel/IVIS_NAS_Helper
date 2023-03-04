package serve_file

import (
	"net/http"
	token_process "ivis_nas/login/token_process"
)

func FileServer(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in
	_, err := token_process.VerifyToken(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.ServeFile(w, r, "/mnt/ivis_nas"+r.URL.Path[11:])
}