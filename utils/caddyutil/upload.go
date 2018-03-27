// gettoken
package caddyutil

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"zyjsxy/go-nginx/utils/aes"
	str "zyjsxy/go-nginx/utils/strings"
)

func UpLoad(w http.ResponseWriter, r *http.Request, config *Config) {
	w.Header().Set("Access-Control-Allow-Origin", config.AllowOrigin)
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	imgSrc := r.FormValue("imgSrc")
	if imgSrc == "" {
		UpLoadFile(w, r, config)
	} else {
		UpLoadImg(w, r, config)
	}
}

func UpLoadImg(w http.ResponseWriter, r *http.Request, config *Config) {
	authStr := r.Header.Get("Authorization")
	aesEnc := aes.AesEncrypt{}
	tokenori, _ := base64.StdEncoding.DecodeString(authStr)
	tokenstr, _ := aesEnc.Decrypt(tokenori)
	id := handleStr(tokenstr, "id:", ",ip:")
	file, handle, err := r.FormFile("files")
	imgSrc := r.FormValue("imgSrc")
	if strings.Index(imgSrc, "?") != -1 {
		imgSrc = str.SubStr(imgSrc, 0, strings.Index(imgSrc, "?"))
	}
	suffix := str.GetSuffix(handle.Filename)
	imgName := str.NTos(id) + suffix
	if str.GetFileName(imgSrc) != str.NTos(id)+suffix && imgSrc != "/img/user.jpg" {
		os.Remove(config.AvatarSrc + imgSrc)
	}

	f, err := os.OpenFile(config.AvatarSrc+"/img/"+imgName, os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)
	if err != nil {
		w.Write([]byte("{\"meta\":{\"code\":20000,\"success\":false,\"message\":" + err.Error() + "},\"data\":null}"))
	} else {
		w.Write([]byte("{\"meta\":{\"code\":20000,\"success\":true,\"message\":\"ok!\"},\"data\":\"" + "/img/" + imgName + "\"}"))
	}
	defer f.Close()
	defer file.Close()
}

func UpLoadFile(w http.ResponseWriter, r *http.Request, config *Config) {
	file, handle, err := r.FormFile("files")
	suffix := str.GetSuffix(handle.Filename)
	now := time.Now()
	mscond := now.Hour()*60*60*1000 + now.Minute()*60*1000 + now.Second()*1000 + int(now.Nanosecond()/1000000)
	fileName := strconv.Itoa(mscond) + suffix
	dir := config.FileSrc + "/file/" + strconv.Itoa(now.Year()) + "-" + strconv.Itoa(int(now.Month())) + "-" + strconv.Itoa(now.Day())
	err = os.MkdirAll(dir, 0755)
	f, err := os.OpenFile(dir+"/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)
	if err != nil {
		w.Write([]byte("{\"meta\":{\"code\":20000,\"success\":false,\"message\":" + err.Error() + "},\"data\":null}"))
	} else {
		w.Write([]byte("{\"meta\":{\"code\":20000,\"success\":true,\"message\":\"ok!\"},\"data\":\"" + "/file/" + strconv.Itoa(now.Year()) + "-" + strconv.Itoa(int(now.Month())) + "-" + strconv.Itoa(now.Day()) + "/" + fileName + "\"}"))
	}
	defer f.Close()
	defer file.Close()
}
