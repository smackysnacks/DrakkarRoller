package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var outputFile *os.File

func init() {
	f, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}

	outputFile = f
}

type Stats struct {
	str    int
	int    int
	wis    int
	wil    int
	con    int
	agi    int
	cha    int
	health int
	gold   int
	luck   int
}

func (s *Stats) Sum() int {
	return s.str + s.int + s.wil + s.agi + s.cha + s.luck
}

func statsFromData(buf []byte) Stats {
	if len(buf) < 0x0b {
		panic("buf not big enuf")
	}

	return Stats{
		str:    int(buf[0x01]),
		int:    int(buf[0x03]),
		wis:    int(buf[0x04]),
		wil:    int(buf[0x05]),
		con:    int(buf[0x06]),
		agi:    int(buf[0x07]),
		cha:    int(buf[0x09]),
		health: int(buf[0x0f]),
		gold:   int(buf[0x0b]),
		luck:   int(buf[0x0a]),
	}
}

func goodStats(stats Stats) bool {
	return stats.agi >= 17 && stats.cha >= 17 && stats.int >= 17 && stats.str >= 17 && stats.wil >= 17 && stats.luck >= 17
}

func exit(code int, err error) {
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func heartbeat(con net.Conn) {
	for {
		time.Sleep(time.Second * 15)
		_, _ = con.Write([]byte{0x2a})
	}
}

func roll(con net.Conn) Stats {
	buffer := make([]byte, 64)

	_, _ = con.Write([]byte("CM10Smacks"))
	_, _ = con.Write([]byte{0x00})
	_, _ = con.Write([]byte{0x53})

	n, err := con.Read(buffer)
	if err != nil {
		exit(1, err)
	}

	return statsFromData(buffer[0:n])
}

func accept(con net.Conn) {
	buffer := make([]byte, 64)
	_, _ = con.Write([]byte{0x41})
	_, err := con.Read(buffer)
	if err != nil {
		exit(1, err)
	}
}

func main() {
	time.Sleep(time.Millisecond * 250)

	con, err := net.Dial("tcp", "127.0.0.1:25042")
	if err != nil {
		exit(1, err)
	}
	go heartbeat(con)

	_, _ = con.Write([]byte{0x52})
	_, _ = con.Write([]byte{0x53})

	buffer := make([]byte, 64)
	n, err := con.Read(buffer)
	if err != nil {
		exit(1, err)
	}
	fmt.Printf("read %d bytes...\n", n)

	bestSum := 0
	for {
		stats := roll(con)
		sum := stats.Sum()
		if sum >= bestSum {
			bestSum = sum
			_, _ = outputFile.WriteString(fmt.Sprintf("%d %+v\n", sum, stats))
		}
		fmt.Printf("%d %+v\n", sum, stats)
		time.Sleep(1 * time.Millisecond)

		if goodStats(stats) {
			accept(con)
			time.Sleep(time.Millisecond * 500)
			exit(0, nil)
		}
	}

	time.Sleep(time.Second * 10)
}
