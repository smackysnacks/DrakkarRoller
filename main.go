package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

var outputFile *os.File
var config Config

func init() {
	var err error

	outputFile, err = os.Create("log.txt")
	if err != nil {
		panic(err)
	}

	bytes, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}

	p, _ := json.MarshalIndent(config, "", "  ")
	fmt.Fprintf(outputFile, "Configuration:\n%s\n", string(p))
}

type Config struct {
	CharacterName   string `json:"CharacterName"`
	CharacterIsMale bool   `json:"CharacterIsMale"`
	MinimumStr      int    `json:"MinimumStr"`
	MinimumInt      int    `json:"MinimumInt"`
	MinimumWis      int    `json:"MinimumWis"`
	MinimumWil      int    `json:"MinimumWil"`
	MinimumCon      int    `json:"MinimumCon"`
	MinimumAgi      int    `json:"MinimumAgi"`
	MinimumCha      int    `json:"MinimumCha"`
	MinimumLuck     int    `json:"MinimumLuck"`
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
	return stats.str >= config.MinimumStr &&
		stats.int >= config.MinimumInt &&
		stats.wis >= config.MinimumWis &&
		stats.wil >= config.MinimumWil &&
		stats.con >= config.MinimumCon &&
		stats.agi >= config.MinimumAgi &&
		stats.cha >= config.MinimumCha &&
		stats.luck >= config.MinimumLuck
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

func characterGenderSymbol(isMale bool) rune {
	if isMale {
		return 'M'
	}
	return 'F'
}

func roll(con net.Conn) Stats {
	buffer := make([]byte, 64)

	_, _ = con.Write([]byte(fmt.Sprintf("C%c10%s", characterGenderSymbol(config.CharacterIsMale), config.CharacterName)))
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
			_, _ = fmt.Fprintf(outputFile, "%d, %+v\n", sum, stats)
		}
		fmt.Printf("%d %+v\n", sum, stats)
		time.Sleep(1 * time.Millisecond)

		if goodStats(stats) {
			accept(con)
			time.Sleep(time.Millisecond * 500)
			exit(0, nil)
		}
	}
}
