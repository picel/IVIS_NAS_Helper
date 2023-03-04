package logout

import (
	"fmt"
	"net/http"
	"ivis_nas/login/token_process"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	// get token value from cookie
	tokenString, err := r.Cookie("jwt_token")
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// delete token from redis
	err = token_process.RemoveToken(tokenString.Value)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// delete cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt_token",
		Value:   "",
		HttpOnly: true,
		Path: "/",
	})

	// redirect to login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
