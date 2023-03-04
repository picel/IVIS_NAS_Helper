package login

import (
	"fmt"
	"net/http"
	"net/url"

	token_process "ivis_nas/login/token_process"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// get id and password from post request
	id := r.FormValue("id")
	pw := r.FormValue("pw")

	// send post request to 192.168.195.1/loginCheck and get response code
	resp, err := http.PostForm("http://192.168.195.1/loginCheck", url.Values{"id": {id}, "pw": {pw}})
	if err != nil {
		fmt.Println(err)
		// redirect to login page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	defer resp.Body.Close()

	// if response code is 200, login success
	if resp.StatusCode == 200 {
		// create jwt token
		ts, err := token_process.CreateToken(id)
		if err != nil {
			fmt.Println(err)
			// redirect to login page
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		// set token to cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "jwt_token",
			Value:   ts,
			HttpOnly: true,
			Path: "/",
		})
		
		// redirect to file list page
		http.Redirect(w, r, "/files/", http.StatusSeeOther)
		return
	} else {
		// redirect to login page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}