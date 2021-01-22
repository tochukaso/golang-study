package env

import "net/http"

func GetSameSitePolicy() http.SameSite {
	sameSitePolicy := http.SameSiteStrictMode
	switch GetEnv().CookieSameSite {
	case "Strinct":
		sameSitePolicy = http.SameSiteStrictMode
	case "Lax":
		sameSitePolicy = http.SameSiteLaxMode
	case "None":
		sameSitePolicy = http.SameSiteNoneMode
	case "Default":
		sameSitePolicy = http.SameSiteDefaultMode
	}

	return sameSitePolicy
}
