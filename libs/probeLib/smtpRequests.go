package probeLib

import (
	`bufio`
	`fmt`
	`io/fs`
	`net`
	`strings`
)

func testOpenRelay(Host string, rPort int, Hostlist fs.File, UserList fs.File){
	req,err := net.Dial("TCP", fmt.Sprintf("%s:%d", Host, rPort))
	if err != nil{
		panic("Could not create dailer object.")
	}
	status1, _ := bufio.NewReader(req).ReadString('\n')
	if strings.Contains(status1, ""){

	}
}
