package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"goweb/goRedis"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	m3u8File = "/home/ubuntu/www/protected/demo/bwf.m3u8"
	// m3u8File  = "protected/demo/bwf.m3u8"
	redisAddr = "goredis:6379"
)

func getHashString(in string) string {
	h := sha1.New()
	h.Write([]byte(in))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
	// hasher := sha256.New()
	// hasher.Write([]byte(in))
	// sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	// fmt.Println(in, sha)
}
func genM3u8Content(filename, timstr string) string {
	var out string
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		if strings.Contains(scanner.Text(), ".ts") {
			name := strings.TrimSpace(scanner.Text())
			ticket := getHashString(fmt.Sprintf("%v-%v", name, timstr))
			out = fmt.Sprintf("%v%v?timestamp=%v&ticket=%v\n", out, name, timstr, ticket)
		} else {
			out = fmt.Sprintf("%v%v\n", out, strings.TrimSpace(scanner.Text()))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return out
}

func genM3u8Handler(w http.ResponseWriter, r *http.Request) {
	timestr := fmt.Sprintf("%v", time.Now().UTC().Unix())
	rkey := getHashString(fmt.Sprintf("%v%v", timestr, m3u8File))
	red := goRedis.GetRedis(redisAddr)
	err := red.Set(rkey, m3u8File, 60*time.Second).Err()
	if err != nil {
		http.Error(w, "File not found.", 404)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"path": fmt.Sprintf("http://localhost:55555/play/%v.m3u8", rkey),
		"code": 0,
	})
}

func playM3u8Handler(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile("/(\\w+?)\\.m3u8")
	s := re.FindStringSubmatch(r.URL.Path)
	log.Println(s)
	if len(s) > 1 {
		red := goRedis.GetRedis(redisAddr)
		val, err := red.Get(s[1]).Result()
		log.Println(val)
		if err != nil {
			http.Error(w, "File not found.", 404)
			return
		}
		fcontent := genM3u8Content(val, fmt.Sprintf("%v", time.Now().UTC().Unix()))
		//Send the headers
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		w.Write([]byte(fcontent))
		// io.Copy(w, ofile) //'Copy' the file to the client
		return
	}
	http.Error(w, "File disappear.", 404)

	// fmt.Println(s)
}

func getDownloadFile(w http.ResponseWriter, r *http.Request) {
	mm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "missing parameter.", 404)
		return
	}
	re := regexp.MustCompile("/(\\w+?)\\.ts")
	s := re.FindStringSubmatch(r.URL.Path)
	if len(s) > 1 {
		if mm["ticket"][0] == getHashString(fmt.Sprintf("%v.ts-%v", s[1], mm["timestamp"][0])) {
			log.Println(s)
			aliasedFile := fmt.Sprintf("/download/demo/%v.ts", s[1])               //this is the nginx alias of the file path
			realFile := fmt.Sprintf("/home/ubuntu/www/protected/demo/%v.ts", s[1]) //this is the physical file path
			filename := fmt.Sprintf("%v.ts", s[1])                                 //this is the file name user will get
			file, err := os.Open(realFile)
			if err != nil {
				http.Error(w, "File not found a.", 404)
				return
			}
			defer file.Close()
			stat, err := file.Stat()
			if err != nil {
				http.Error(w, "File not found b.", 404)
				return
			}
			w.Header().Set("Cache-Control", "public, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Content-Type", "video/MP2T")
			w.Header().Set("Content-Length", fmt.Sprintf("%v", stat.Size()))
			w.Header().Set("Content-Disposition: attachment; filename=", filename)
			w.Header().Set("Content-Transfer-Encoding", "binary")
			w.Header().Set("X-Accel-Redirect", aliasedFile)
			return
		}
		http.Error(w, "File not found c.", 404)
		return
	}
	http.Error(w, "File not found d.", 404)
	return

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.Contains(r.URL.Path, "/ask"):
		{
			genM3u8Handler(w, r)
		}
	case strings.Contains(r.URL.Path, ".m3u8"):
		{
			playM3u8Handler(w, r)
		}
	case strings.Contains(r.URL.Path, ".ts"):
		{
			getDownloadFile(w, r)
		}
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe("0.0.0.0:55550", nil)
}
