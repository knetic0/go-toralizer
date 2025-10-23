package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"unsafe"
)

const (
	PROXY     string  = "127.0.0.1"
	PROXYPORT string  = "9050"
	USERNAME  string  = "toraliz"
	RESSIZE   uintptr = unsafe.Sizeof(Res{})
)

/*
SOCKS4 CONNECT REQUEST FORMAT:
+----+----+----+----+----+----+----+----+----+----+....+----+
| VN | CD | DSTPORT |      DSTIP        | USERID       |NULL|
+----+----+----+----+----+----+----+----+----+----+....+----+

	1    1      2              4           variable       1
*/
type proxyRequest struct {
	VN      uint8
	CD      uint8
	DstPort uint16
	DstIP   uint32
}

/*
SOCKS4 RESPONSE FORMAT:
+----+----+----+----+----+----+----+----+
| VN | CD | DSTPORT |      DSTIP        |
+----+----+----+----+----+----+----+----+

	1    1      2              4
*/
type proxyResponse struct {
	VN      uint8
	CD      uint8
	DstPort uint16
	DstIP   uint32
}

type Req proxyRequest
type Res proxyResponse

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <host> <port>", os.Args[0])
	}

	host := os.Args[1]
	var port int
	fmt.Sscanf(os.Args[2], "%d", &port)

	remoteAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(PROXY, PROXYPORT))
	if err != nil {
		log.Fatalf("Failed to resolve proxy address: %v", err)
	}

	conn, err := net.DialTCP("tcp", nil, remoteAddr)
	if err != nil {
		log.Fatalf("Failed to connect to proxy: %v", err)
	}
	defer conn.Close()

	log.Printf("Connected to proxy %s:%s", PROXY, PROXYPORT)

	reqBytes := buildRequest(host, port)

	if _, err := conn.Write(reqBytes); err != nil {
		log.Fatalf("Failed to send SOCKS4 request: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	buffer := make([]byte, RESSIZE)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}
	if n < int(RESSIZE) {
		log.Printf("Warning: response size smaller than expected (%d bytes)", n)
	}

	resp := &Res{}
	reader := bytes.NewReader(buffer)
	if err := binary.Read(reader, binary.BigEndian, resp); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	if resp.CD != 90 {
		log.Fatalf("Connection failed, code: %d", resp.CD)
	}

	log.Printf("✅ SOCKS4 connection established to %s:%d through Tor proxy", host, port)

	httpRequest := fmt.Sprintf("GET / HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", host)
	if _, err := conn.Write([]byte(httpRequest)); err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	response := new(bytes.Buffer)
	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)
		if n > 0 {
			response.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	log.Println("== HTTP Response ==")
	fmt.Println(response.String())
}

func buildRequest(dstip string, dstport int) []byte {
	req := &Req{
		VN:      4,
		CD:      1,
		DstPort: uint16(dstport),
	}

	ip := net.ParseIP(dstip).To4()
	if ip == nil {
		log.Fatalf("Invalid IPv4 address: %s", dstip)
	}
	req.DstIP = binary.BigEndian.Uint32(ip)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, req); err != nil {
		log.Fatalf("Failed to serialize request: %v", err)
	}

	buf.WriteString(USERNAME)
	buf.WriteByte(0x00)

	return buf.Bytes()
}
