package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/Alexduanran/Distributed-System-MP1/msg"
	"github.com/Alexduanran/Distributed-System-MP1/tcp"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type process struct {
	Ip, Port string
	Conn net.Conn // Conn is nil when the current process does not have a client connecting with the server at Port
}				  // else, Conn stores the connection

var _id string // id of the current process
var minDelay, maxDelay int // filled in in readConfig()
var processes map[string]process // stores all information of the processes in the config file

// init reads the user input and the config file
func init() {
	// read user input ID
	flag.StringVar(&_id, "ID", "1", "help message for _id")
	flag.Parse()

	// init processes and read config.txt
	processes = make(map[string]process)
	readConfig()
}

func main() {
	// a goroutine to handle user inputs to send messages between processes
	go handleMessages()

	// start server listening
	tcp.MultiThreadedServer(processes[_id].Ip, processes[_id].Port, handleConnection)
}

func handleConnection(conn net.Conn, listen *bool) {
	// keeps on listening for incoming messages sent from conn
	for {
		var msg msg.Message
		tcp.UnicastReceive(conn, &msg)
		fmt.Printf("<<< Received “%v” from process %v, system time is %v\n\n", msg.Msg, msg.Id, time.Now().Format("15:04:05.000"))
	}
}

// handleMessage reads in user input and makes appropriate action
func handleMessages() {
	scanner := bufio.NewScanner(os.Stdin)

	// keeps on waiting for user inputs
	for scanner.Scan() {
		go func() {
			input := strings.Split(scanner.Text(), " ")

			// if user did not input 3 defined inputs
			if len(input) < 3 {
				fmt.Println("Not enough input\n")
				return
			}

			send, idStr, message := input[0], input[1], strings.Join(input[2:], " ")
			idInt, err := strconv.Atoi(idStr)

			// if the input user enters is invalid
			if err != nil || send != "send" || idInt < 1 || idInt >= len(processes) {
				fmt.Println("Invalid input\n")
				return
			}

			// if this is the first time messaging with idStr,
			// build connection with it first and stores the connection in processes
			if proc, ok := processes[idStr]; ok && proc.Conn == nil {
				proc.Conn, err = tcp.Connect(processes[idStr].Ip, processes[idStr].Port)
				if err != nil {
					fmt.Println("Process not yet started\n")
					return
				}
				processes[idStr] = proc
			}

			if err != nil {
				fmt.Println("Process not yet started\n")
				return
			}

			fmt.Printf(">>> Sent “%v” to process %v, system time is %v\n\n", message, idStr, time.Now().Format("15:04:05.000"))

			// simulate network delay
			source := rand.NewSource(time.Now().UnixNano())
			random := rand.New(source)
			duration := time.Duration(random.Intn(maxDelay-minDelay)+minDelay) * time.Millisecond
			time.Sleep(duration)

			// send message to the destination conn
			newMsg := msg.Message{_id, message}
			tcp.UnicastSend(processes[idStr].Conn, newMsg)
		}()
	}
}

func readConfig() {
	f, err := os.Open("config.txt")
	checkError(err, "Open error")
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// read in minDelay and maxDelay from the config file
	scanner.Scan()
	line := strings.Split(scanner.Text(), " ")
	minDelay, err = strconv.Atoi(line[0])
	checkError(err, "String->Int convert error")
	maxDelay, err = strconv.Atoi(line[1])
	checkError(err, "String->Int convert error")

	// read in information about the processes line by line and stores them into processes as a map of process
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		idStr, ip, port := line[0], line[1], line[2]
		processes[idStr] = process{ip, port, nil}
	}
}

// checkError checks for error in err, exits
// and prints given error message errMSG to the console if err is not nil
func checkError(err error, errMsg string) {
	if err != nil {
		if errMsg == "" {
			errMsg = "Fatal error"
		}
		fmt.Println(errMsg, ":", err.Error())
		os.Exit(1)
	}
}