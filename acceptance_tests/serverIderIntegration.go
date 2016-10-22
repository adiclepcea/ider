package main

import (
	"encoding/json"
	"fmt"
	"ider"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

//MockRequest is used as the base request form
type MockRequest struct {
	ID         int64  `json:"ID"`
	ItemName   string `json:"name"`
	ItemYear   int    `json:"year"`
	ItemAuthor string `json:"author"`
}

var port int
var serverID int
var idProvider *ider.Ider
var dir string
var file *os.File
var writer chan string
var ln net.Listener

func initServer() error {
	var err error
	idProvider, err = ider.NewIder(uint(serverID))
	if err != nil {
		return fmt.Errorf("Error while creating the id generator: %s\n", err.Error())
	}

	writer = make(chan string, 100000)

	return nil
}

//this is to mock a single request
func insertOne(response http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	mockRequest := MockRequest{}
	err := decoder.Decode(&mockRequest)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Incorrect request"))
		return
	}
	//we now have the request with the id incorporated
	//the next step would be to store it
	mockRequest.ID = idProvider.GenerateID()

	response.Write([]byte(fmt.Sprintf("%d", mockRequest.ID)))

	writer <- fmt.Sprintf("%d\n", mockRequest.ID)
}

//this is to mock the insert of a batch of items/requests
func insertMany(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var mockRequests []MockRequest
	mockRequests = make([]MockRequest, 0)
	err := decoder.Decode(&mockRequests)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Incorrect request"))
		return
	}

	for i := range mockRequests {
		mockRequests[i].ID = idProvider.GenerateID()
		response.Write([]byte(fmt.Sprintf("%d,", mockRequests[i].ID)))
		writer <- fmt.Sprintf("%d\n", mockRequests[i].ID)
	}
}

func stop(response http.ResponseWriter, request *http.Request) {

	response.Write([]byte("OK"))
	request.Body.Close()
	ln.Close()

	close(writer)

	file.Close()

	log.Print("Shutting down ...")

	os.Exit(0)
}

func validateParams() error {
	var err error

	envServerID := os.Getenv("SERVER_ID")
	envServerPORT := os.Getenv("SERVER_PORT")

	serverID, err = strconv.Atoi(envServerID)
	if err != nil || envServerID == "" {
		return fmt.Errorf("Invalid server id: %s. Must be an integer\n", envServerID)
	}

	port, err = strconv.Atoi(envServerPORT)
	if err != nil || envServerPORT == "" {
		return fmt.Errorf("Invalid server port: %s. Must be an integer\n", envServerPORT)
	}

	if len(os.Args) < 2 {
		return fmt.Errorf("Please provide an output folder i.e.\n %s /data\n", os.Args[0])
	}

	if port > 65535 || port < 80 {
		if err != nil {
			return fmt.Errorf("Invalid server port %s. Must be between 80 and 65535\n", os.Args[1])
		}
	}
	dir = os.Args[1]
	return nil
}

func main() {
	var err error
	if err = validateParams(); err != nil {
		log.Fatal(err.Error())
		return
	}
	if err = initServer(); err != nil {
		log.Fatal(err.Error())
		return
	}

	file, err = os.OpenFile(filepath.Join(dir, fmt.Sprintf("%d.txt", serverID)), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)

	if err != nil {
		log.Fatal(err.Error())
		return
	}
	//we write as soon as we can the data in the file
	//in a production environment a Redis like server would be more appropriate
	go func() {
		for {
			s := <-writer

			_, err = file.Write([]byte(s))
			if err != nil {
				log.Printf(err.Error())
			}
			log.Printf("read %s\n", s)

		}
	}()

	http.HandleFunc("/insertone", insertOne)
	http.HandleFunc("/insertmany", insertMany)
	http.HandleFunc("/stop", stop)
	ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	fmt.Printf("Runing on port %d", port)
	http.Serve(ln, nil)
}
