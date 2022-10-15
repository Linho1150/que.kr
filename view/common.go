package view

import (
	"net/http"
	"net/url"
	"strings"
)

func ValidationURL(targetUrl string) bool{
	parseUrl, err := url.Parse(targetUrl)
	if err != nil {
		return true
	}
	if !(parseUrl.Scheme == "http" || parseUrl.Scheme == "https"){
		return true
	}
	return false
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
			IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
			IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func ReadUserReferer(r *http.Request) string {
	return r.Header.Get("referer")
}

func ReadDeivceMobile(r *http.Request) bool  {
	userAgent:=r.Header.Get("User-Agent")
	return strings.Contains(userAgent,"Mobi")
}