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
	typeDemo  = "demo"
	typeLocal = "localdemo"
	typeEncry = "encdemo"
)

type environment struct {
	folderName, port, m3u8FileGeneral, m3u8FileLocal, m3u8FileEncry, redisAddr string
	isLocal                                                                    bool
}

func (d *environment) setParameter(local bool) {
	d.isLocal = local
	if local {
		d.port = "55550"
		d.folderName = "protected/"
		d.m3u8FileGeneral = "protected/demo/bwf.m3u8"
		d.m3u8FileLocal = "protected/localdemo/localbwf.m3u8"
		d.m3u8FileEncry = "protected/encdemo/encbwflocal.m3u8"
		d.redisAddr = "localhost:6379"
	} else {
		d.port = "55555"
		d.folderName = "/home/ubuntu/www/protected/"
		d.m3u8FileGeneral = "/home/ubuntu/www/protected/demo/bwf.m3u8"
		d.m3u8FileLocal = "/home/ubuntu/www/protected/localdemo/localbwf.m3u8"
		d.m3u8FileEncry = "/home/ubuntu/www/protected/encdemo/encbwf.m3u8"
		d.redisAddr = "goredis:6379"
	}
}

var env environment

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
func genM3u8Content(filename, timstr string, para map[string]interface{}) string {
	var out string
	var appendstr string
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for kk, vv := range para {
		appendstr = fmt.Sprintf("%v&%v=%v", appendstr, kk, vv)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		if strings.Contains(scanner.Text(), ".ts") {
			name := strings.TrimSpace(scanner.Text())
			ticket := getHashString(fmt.Sprintf("%v-%v", name, timstr))
			out = fmt.Sprintf("%v%v?timestamp=%v&ticket=%v%v\n", out, name, timstr, ticket, appendstr)
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
	filetype := typeDemo
	m3u8File := env.m3u8FileGeneral
	mm, err := url.ParseQuery(r.URL.RawQuery)
	fmt.Println(mm)
	if len(mm) > 0 {
		if item, ok := mm["type"]; ok {
			switch item[0] {
			case typeLocal:
				{
					filetype = typeLocal
					m3u8File = env.m3u8FileLocal
				}
			case typeEncry:
				{
					filetype = typeEncry
					m3u8File = env.m3u8FileEncry
				}
			}
		}
	}
	timestr := fmt.Sprintf("%v", time.Now().UTC().Unix())
	rkey := getHashString(fmt.Sprintf("%v%v", timestr, m3u8File))
	red := goRedis.GetRedis(env.redisAddr)
	err = red.Set(rkey, m3u8File, 60*time.Second).Err()
	if err != nil {
		http.Error(w, "File not found.", 404)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"path": fmt.Sprintf("http://localhost:%v/play/%v.m3u8?type=%v", env.port, rkey, filetype),
		"code": 0,
	})
}

func playM3u8Handler(w http.ResponseWriter, r *http.Request) {
	mm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "missing parameter.", 404)
		return
	}
	re := regexp.MustCompile("/(\\w+?)\\.m3u8")
	s := re.FindStringSubmatch(r.URL.Path)
	log.Println(s)
	if len(s) > 1 {
		red := goRedis.GetRedis(env.redisAddr)
		val, err := red.Get(s[1]).Result()
		log.Println(val)
		if err != nil {
			http.Error(w, "File not found.", 404)
			return
		}
		fcontent := genM3u8Content(val,
			fmt.Sprintf("%v", time.Now().UTC().Unix()),
			map[string]interface{}{"type": mm["type"][0]})
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
			filename := fmt.Sprintf("%v.ts", s[1]) //this is the file name user will get
			if env.isLocal {
				http.Redirect(w, r,
					fmt.Sprintf("http://localhost:%v/download/%v/%v", env.port, mm["type"][0], filename), http.StatusMovedPermanently)
				return
			}
			aliasedFile := fmt.Sprintf("/download/%v/%v", mm["type"][0], filename)      //this is the nginx alias of the file path
			realFile := fmt.Sprintf("%v%v/%v", env.folderName, mm["type"][0], filename) //this is the physical file path

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

func getStaticFile(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile("/(\\w+?)\\.key")
	s := re.FindStringSubmatch(r.URL.Path)
	if len(s) > 1 {
		filename := fmt.Sprintf("%v.key", s[1])
		aliasedFile := fmt.Sprintf("/download/key/%v", filename)
		realFile := fmt.Sprintf("%vkey/%v", env.folderName, filename)

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

		w.Header().Set("Content-Length", fmt.Sprintf("%v", stat.Size()))
		w.Header().Set("Content-Disposition: attachment; filename=", filename)
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("X-Accel-Redirect", aliasedFile)
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
	case strings.Contains(r.URL.Path, "/key"):
		{
			getStaticFile(w, r)
		}
	default:
		{
			http.Error(w, `Invalid operation!`, http.StatusUnsupportedMediaType)
		}
	}
}

func main() {
	env.setParameter(false)
	http.HandleFunc("/", indexHandler)
	if env.isLocal {
		fs := http.FileServer(http.Dir("protected/"))
		http.Handle("/download/", http.StripPrefix("/download/", fs))
	}
	http.ListenAndServe("0.0.0.0:55550", nil)
}
