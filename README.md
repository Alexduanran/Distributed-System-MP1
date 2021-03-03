# Distributed-System-MP1
MP1 for Distributed System Spring 2021

A simple network unicast communication simulation with simulated delay

## How To Run
Format the configaration file *config.txt* as follow

```txt
minDelay(ms) maxDelay(ms)
ID1 IP1 port1
ID2 IP2 port2
```
with the minimum and maximum delay (in milliseconds) between each message specified in the first line, follow by a number of processes, each on a separate line, with unique ID and Port as well as an IP address of your choice. 

Then, in separate terminals, start each process with its unique ID assigned in *config.txt*
```bash
go run main.go -ID 1
```
```bash
go run main.go -ID 2
```
```bash
...
```
#
In order to send message between each process, in the command line
```bash
send [id] [message]
```
where id is the unique ID of the receiver process.

## Package Design
### Main
A **process** struct for each process sepcified in *config.txt*
```go
type process struct {
  Ip string      // IP address obtained from config.txt
  Port string    // Port number obtained from config.txt
  Conn net.Conn  // The client connection from the current process to this process. If no connection has been built, Conn is nil
}
```
Each process holds a map of **process** with key ID that contains information to all the processes in *config.txt*
```go
var processes map[string]process
```
Each process only creates a client connection with another process in the first *send* command.

The communications in the latter rounds use **Conn** stored in each **process** object.

```go
func main() {
  // a goroutine to handle user inputs to send messages between processes
  go handleMessages()

  // start server listening
  tcp.MultiThreadedServer(processes[_id].Port, handleConnection)
}
```

### Msg
```go
type Message struct {
  Id string   \\ the ID of the process the message is sent from
  Msg string  \\ the content of the message
}
```

### TCP
Seperates TCP logic from the main process.
Supports building servers, connecting clients, and sending/encoding/decoding messages using ` "encoding/gob" ` in between. 
