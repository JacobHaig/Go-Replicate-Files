package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"
)

type FileStructure struct {
	Data     []byte
	Filename string
}

func main() {
	folderIn := `C:\Users\jacob\OneDrive\Desktop\Move From`
	folderOut := `C:\Users\jacob\OneDrive\Desktop\Move To`

	listen, err := net.Listen("tcp", ":20")
	if err != nil {
		println("Error: ", err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			println("Error: ", err.Error())
		}

		println("Connection made to client: ", conn.RemoteAddr().String())
		server(conn, folderIn, folderOut)
		println("Connection closed: ", conn.RemoteAddr().String())
	}
}

func server(conn net.Conn, folderIn string, folderOut string) {
	dirs := readDir(folderIn)
	count := 0
	for _, file := range dirs {
		if !file.IsDir() {
			count++
		}
	}

	dirSize, err := json.Marshal(count)
	if err != nil {
		println("Error: ", err.Error())
	}

	sendData(conn, dirSize)

	for _, file := range dirs {
		if !file.IsDir() {
			println(file.Name())

			data := readFile(folderIn, file.Name())
			f := FileStructure{data, file.Name()}
			encoded, err := json.Marshal(f)

			if err != nil {
				println("Error: ", err.Error())
			}

			println("data_bytes: ", len(f.Data), " encoded_bytes: ", len(encoded))
			sendData(conn, encoded)

			deleteFile(folderIn, file.Name())
		}
	}
	time.Sleep(1 * time.Second)
}

func bytesSizeAsBytes(data []byte) []byte {
	size := int32(len(data))
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, size)

	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func sendData(conn net.Conn, data []byte) {
	sizeData := bytesSizeAsBytes(data)

	conn.Write(sizeData)
	conn.Write(data)
}

func readDir(folder string) []os.FileInfo {
	arr, err := ioutil.ReadDir(folder)

	if err != nil {
		println("Error: ", err.Error())
	}
	return arr
}

func readFile(folder string, filename string) []byte {
	b, err := ioutil.ReadFile(folder + "\\" + filename)

	if err != nil {
		println("Error: ", err)
	}

	return b
}

func writeFile(content []byte, folder string, filename string) {
	err := ioutil.WriteFile(folder+"\\"+filename, content, 0777)

	if err != nil {
		print("Error: ", err)
	}
}

func deleteFile(folder string, filename string) {
	err := os.Remove(folder + "\\" + filename)

	if err != nil {
		println("Error: ", err)
	}
}
