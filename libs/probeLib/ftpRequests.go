package probeLib

import (
	`bufio`
	`fmt`
	`log`
	`net`
	`strings`
)

func ftpOpenAnonAccess(rHost string, rPort int){
	conn, err := net.Dial("TCP", fmt.Sprintf("%s:%d",rHost, rPort))
	if err != nil{
		panic("Cannot continue.")
	}
	fmt.Fprintf(conn, "USER anonymous\n")
	stats, _ := bufio.NewReader(conn).ReadString('\n')
	if strings.Contains(stats, "230") {
		log.Println("Anonymous login enabled, we got a 230 response. Please ensure that this is not a honeypot.")
		return
	}else {
		fmt.Fprintf(conn, "USER guest")
		stats2, _ := bufio.NewReader(conn).ReadString('\n')
		if strings.Contains(stats2, "230") {
			log.Println("User guest was successful!")
		}
	}
	defer conn.Close()

}


func bruteForceFTP(rHost string, rPort int){

}