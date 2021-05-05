package main

import (
	`bufio`
	`encoding/csv`
	`fmt`
	`io`
	`os`
	`strconv`
	`strings`

	`github.com/gookit/color`
	`webrequestsingo.com/libs/bannerLib`
	`webrequestsingo.com/libs/probeLib`
)

var proxyHost = false
var (
	fileInfo *os.FileInfo
	err      error
)

type portConfig struct {
	http []string `yaml:"http"`
	https []string `yaml:"https"`
	sql []string `yaml:"sql"`
	shell_services []string `yaml:"shell_services"`
	mail []string `yaml:"mail"`
}

func main() {
	var user_pw string
	var cred_len int
	bannerLib.PrintMainBanner()
	if len(os.Args[1:]) < 4 {
		user_pw = "admin:123456"
		color.Style{color.BgWhite, color.FgRed, color.Bold}.Printf("\n\n[ !! ] Hey buddy, you need to wake up." +
			"Need at least 4 args to make this thing work awesomely. Since thats not happening, we are going to " +
			"default to %s [ !! ]", user_pw)
		cred_len = 0
	} else {
		cred_len = 1
		user_pw = os.Args[4]
	}
	color.Style{color.BgBlack, color.FgGreen, color.Bold}.Println("\n[ !! ] Starting. [ !! ]\n")
	var schema = ""
	var rPort, _ = strconv.Atoi(os.Args[2])
	var rHost = os.Args[1]
	if strings.Contains(os.Args[2], "443") {
		schema = "https"
		color.Style{color.BgBlack, color.FgGreen, color.Bold}.Println("[ !! ] Selecting HTTPS [ !! ]")
	} else {
		schema = "http"
		color.Style{color.BgBlack, color.FgGreen, color.Bold}.Println("[ !! ] Selecting HTTP [ !! ]")
	}
	if os.Args[3] == "false" {
		proxyHost = false
		color.Style{color.BgRed, color.FgWhite, color.Bold}.Println("[ !! ] Not proxying any requests! [ !! ]")
	} else {
		proxyHost = true
		color.Style{color.BgBlack, color.FgGreen, color.Bold}.Println("[ !! ] Proxying all requests! [ !! ]")
	}
	if probeLib.TestHost(rHost, rPort, schema, proxyHost) {
		if _, err := os.Stat(user_pw); os.IsNotExist(err) {
			color.Style{color.BgRed, color.FgWhite, color.Bold}.Println(
				"[ !! ] Supplied arg was not a file, using as the username/password[ !! ]")
			probeLib.Brute401(rHost, rPort, proxyHost, user_pw, schema, cred_len)
		} else {
			file, _ := os.Stat(user_pw)
			if file.IsDir() {

			} else {
				if strings.Contains(user_pw, "txt") {
					read_file, err := os.Open(user_pw)
					if err == io.EOF{
						color.Style{color.BgRed, color.FgWhite, color.Bold}.Printf(
							"[ !! ] Reached end of file. [ !! ]\n")
					}
					defer read_file.Close()
					scann := bufio.NewScanner(read_file)
					for scann.Scan() {
						probeLib.Brute401(rHost, rPort, proxyHost, scann.Text(), schema, cred_len)
					}
				} else if strings.Contains(user_pw, "csv") {
					file, err := os.Open(user_pw)
					defer file.Close()
					if err == io.EOF{
						color.Style{color.BgRed, color.FgWhite, color.Bold}.Printf(
							"[ !! ] Reached end of file. [ !! ]\n")
					}
					reader,_ := csv.NewReader(file).ReadAll()
					for value,_ := range reader{
						if reader[value][0] != "" {
							fmt.Printf("%s : %s\n", reader[value][0], reader[value][1])
						}
					}
				}
			}
		}
	}
}
