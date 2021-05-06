package probeLib

import (
	`crypto/tls`
	`encoding/base64`
	`fmt`
	`io/ioutil`
	`math/rand`
	`net/http`
	`net/url`
	`time`

	`github.com/gookit/color`
	`gopkg.in/yaml.v2`
)


var client = &http.Client{
	Timeout: 10 * time.Second,
}

//defaulting to tor.
type proxyConfig struct {
	Type string `yaml:"Type"`
	Address string `yaml:"Address"`
	Port int `yaml:"Port"`
	Schema string `yaml:"Schema"`
}

func (c *proxyConfig) loadConfigs(data []byte) error{
	return yaml.Unmarshal(data, c)
}


func TestHost(RemoteHost string, RemotePort int, schema string, proxyHost bool) bool {
	conn, _ := http.NewRequest("GET", fmt.Sprintf("%s://%s:%d", schema, RemoteHost, RemotePort),
		nil)
	if proxyHost{
		// defaulting to Tor.
		data, _ := ioutil.ReadFile("configs/proxies/proxies.yml")
		var config proxyConfig
		if err := config.loadConfigs(data); err != nil{
			color.Style{color.BgWhite, color.FgRed, color.Bold}.Printf("%w", err)
		}
		proxyStringer := fmt.Sprintf("%s://%s:%d", config.Schema, config.Address, config.Port)
		urlified, _ := url.Parse(proxyStringer)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy: http.ProxyURL(urlified),
		}
		client.Transport = tr
	}
	conn.Header.Set("User-Agent",
		"Baiduspider+(+http://www.baidu.com/search/spider.htm);" +
		"googlebot|baiduspider|baidu|spider|sogou|bingbot|bot|yahoo|soso|sosospider|360spider|youdaobot|jikeSpider;)")
	conn.Header.Set("Connection", "close")
	conn.Header.Set("Accept", "*/*")
	conn.Header.Set("Content-Type", "text/plain")
	conn.Header.Set("Referer", fmt.Sprintf("https://www.google.com/search?client=firefox-b-d&q='%s'",
		RemoteHost))
	resp, errs := client.Do(conn)
	if errs != nil{
		color.Style{color.BgWhite, color.FgRed, color.Bold}.Printf("%w", errs)
	}
	defer resp.Body.Close()
	if resp.Status != ""{
		switch resp.StatusCode{
		case http.StatusOK:
			color.Style{color.BgWhite, color.FgRed, color.Bold}.Printf("[ !! ] Returned Status: %s [ !! ]", resp.Status)
			return false
		case http.StatusForbidden:
			color.Style{color.BgWhite, color.FgRed, color.Bold}.Println("[ !! ] Got a forbidden, might not be able to brute force it. [ !! ]")
			return false
		case http.StatusUnauthorized:
			color.Style{
			color.BgGreen,
			color.FgWhite,
			color.Bold}.Printf("[ !! ] Going to start brute force! Got a 401! [ !! ]\n")
			if resp.Header["Www-Authenticate"] != nil {
				color.Style{color.BgGreen, color.FgWhite, color.Bold}.Printf("[ !! ] Identified header: %s [ !! ]\n",
					resp.Header["Www-Authenticate"])
				return true
			}else {
				color.Style{
					color.BgWhite,
					color.FgRed,
					color.Bold}.Printf("[ !! ] We did get a 401, but were not able to identify an auth realm." +
						"This could mean that the 401 is a false 401, OR that we need to supply a client cert. Aborting." +
						"[ !! ]\n")
				return false
			}
		case http.StatusRequestTimeout:
			color.Style{color.BgWhite, color.FgRed, color.Bold}.Println("[ !! ] Connection timed out :( [ !! ]")
			return false
		case http.StatusBadGateway:
			color.Style{color.BgWhite, color.FgRed, color.Bold}.Println(
				"[ !! ] Need to retry connection, got a bad gateway. Could be you, could be them. [ !! ]")
			return false
		case http.StatusTeapot:
			color.Style{color.BgWhite, color.FgRed, color.Bold}.Println(
				"[ !! ] Ok, so aparently the server is a teapot. But we will still continue. [ !! ]")
			return true
		case http.StatusMethodNotAllowed:
			color.Style{color.BgWhite, color.FgRed, color.Bold}.Println(
				"[ !! ] Going to have to retry the request method, the endpoint is there, but " +
					"does not allow GET requests. [ !! ]")
			return false
		case 444:
			color.Style{color.BgWhite, color.FgRed, color.Bold}.Println(
				"[ !! ] Appears to be NGINX or a ward status code of 444 to block bots. will try to continue. " +
					"to see what we can find, starting directory brute forcing. [ !! ]")
			return false
		default:
			color.Style{color.BgRed, color.FgWhite, color.Bold}.Printf(
				"[ ** ] Unable to determine status code. %s[ ** ]",
				resp.Status)
			return false
		}
	}else{
		return false
	}
}

func Brute401(rHost string, rPort int, proxyHost bool, creds string, schema string, credlen int) bool {
	conn, _ := http.NewRequest("POST", fmt.Sprintf("%s://%s:%d/", schema, rHost, rPort), nil)
	if proxyHost{
		data, err := ioutil.ReadFile("configs/proxies/proxies.yml")
		if err != nil{
			color.Style{color.BgBlack, color.FgRed, color.Bold}.Printf("%w", err)
		}
		var config proxyConfig
		if err := config.loadConfigs(data); err != nil{
			color.Style{color.BgWhite, color.FgRed, color.Bold}.Printf("%v", err)
		}
		proxyStringer := fmt.Sprintf("%s://%s:%d", config.Schema, config.Address, config.Port)
		urlified,_ := url.Parse(proxyStringer)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy: http.ProxyURL(urlified),
		}
		client.Transport = tr
	}
	conn.Header.Set("User-Agent",
		"Baiduspider+(+http://www.baidu.com/search/spider.htm);" +
		"googlebot|baiduspider|baidu|spider|sogou|bingbot|bot|yahoo|soso|sosospider|360spider|youdaobot|jikeSpider;)")
	conn.Header.Set("Connection", "close")
	conn.Header.Set("Accept", "*/*")
	conn.Header.Set("Content-Type", "text/plain")
	conn.Header.Set("Referer", fmt.Sprintf("https://www.google.com/search?client=firefox-b-d&q='%s'", rHost))
	counter_key := 0
	if credlen  != 0 {
		for i := 0; i <= len(creds); i++ {
			color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] \nTrying: %s \n"+
				"base64 encoded: %s\n[ !! ]\n", creds, base64.RawURLEncoding.EncodeToString([]byte(creds)))
			conn.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.RawURLEncoding.EncodeToString([]byte(creds))))
			respi, _ := client.Do(conn)
			defer respi.Body.Close()
			stats := respi.StatusCode
			switch stats {
			case http.StatusOK:
				color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] Creds found: %s [ !! ]",
					creds)
				fmt.Println()
				color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] Status Code: %d [ !! ]",
					respi.StatusCode)
				fmt.Println()
				color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] Headers: %s [ !! ]",
					respi.Header)
				fmt.Println()
				if respi.Cookies() != nil {
					for _, c := range respi.Cookies() {
						color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] Cookies:\nName: %s\nValue: %s\nDomain: %s\nDoes it Expire? %s\nPath: %s \n[ !! ]",
							c.Name, c.Value, c.Domain, c.Expires, c.Path)
					}
				}
				fmt.Println()
				color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] Host: %v [ !! ]",
					respi.Request.Host)
				fmt.Println()
				color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] Auth: %v [ !! ]",
					respi.Request.Header["Authorization"])
				return true
			case http.StatusForbidden:
				color.Style{color.BgWhite, color.FgRed, color.Bold}.Println(
					"[ ** ] Bad credentials, need to continue. [ ** ]")
				continue
			case http.StatusUnauthorized:
				color.Style{color.BgWhite, color.FgRed, color.Bold}.Printf(
					"\n[ ** ] Bad credentials, need to continue [ ** ]\n")
				continue
			case http.StatusRequestTimeout:
				counter_key += 1
				if counter_key >= 5 {
					color.Style{color.BgWhite, color.FgRed, color.Bold}.Print(
						"\nMax time out number hit, quitting.\n")
					return false
				}
				randSleeper := rand.Intn(150) / 2
				color.Style{color.BgWhite, color.FgRed, color.Bold}.Printf(
					"\n[ !! ] Request timed out!! Sleeping for: %d [ !! ]\n", randSleeper)
				color.Style{color.BgWhite, color.FgRed, color.Bold}.Printf(
					"\n[ !!! ] Number of timeouts left: %d [ !!! ]\n", counter_key-5)
				time.Sleep(time.Duration(randSleeper))
				continue
			default:
				fmt.Printf("")
				continue
			}
		}
	}else{
		color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] \nTrying the only creds we have" +
			": %s \nbase64 encoded: %s\n[ !! ]\n", creds, base64.RawURLEncoding.EncodeToString([]byte(creds)))
		conn.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.RawURLEncoding.EncodeToString([]byte(creds))))
		respi, _ := client.Do(conn)
		defer respi.Body.Close()
		stats := respi.StatusCode
		switch stats {
		case http.StatusOK:
			color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] Creds found: %s [ !! ]",
				creds)
			fmt.Println()
			color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] Status Code: %d [ !! ]",
				respi.StatusCode)
			fmt.Println()
			color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] Headers: %s [ !! ]",
				respi.Header)
			fmt.Println()
			if respi.Cookies() != nil {
				for _, c := range respi.Cookies() {
					color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] Cookies:\nName: %s\nValue: %s\nDomain: %s\nDoes it Expire? %s\nPath: %s \n[ !! ]",
						c.Name, c.Value, c.Domain, c.Expires, c.Path)
				}
			}
			fmt.Println()
			color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] Host: %v [ !! ]",
				respi.Request.Host)
			fmt.Println()
			color.Style{color.BgBlue, color.FgWhite, color.Bold}.Printf("[ !! ] Auth: %v [ !! ]",
				respi.Request.Header["Authorization"])
			return true
		case http.StatusForbidden:
			color.Style{color.BgWhite, color.FgRed, color.Bold}.Println(
				"[ ** ] Bad credentials, need to continue. [ ** ]")
			break
		case http.StatusUnauthorized:
			color.Style{color.BgWhite, color.FgRed, color.Bold}.Printf(
				"\n[ ** ] Bad credentials, need to continue [ ** ]\n")
			break
		case http.StatusRequestTimeout:
			color.Style{color.BgWhite, color.FgRed, color.Bold}.Printf(
				"\n[ !! ] Request timed out!! [ !! ]\n")
			break
		default:
			fmt.Printf("")
			break
		}
	}
	return false
}

func directoryBrute(rHost string, rPort int, wordlist []string){

}