package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

type FileStructure struct {
	Data     []byte
	Filename string
}

func main() {
	folderIn := `C:\Users\jacob\OneDrive\Desktop\Move From`
	folderOut := `C:\Users\jacob\OneDrive\Desktop\Move To`

	conn, err := net.Dial("tcp", "127.0.0.1:20")
	if err != nil {
		println("Error: ", err.Error())
	}

	client(conn, folderIn, folderOut)
}

func client(conn net.Conn, folderIn string, folderOut string) {
	dirSize := receiveData(conn)
	var size int
	json.Unmarshal(dirSize, &size)

	for i := 0; i < size; i++ {
		b := receiveData(conn)

		var file FileStructure
		json.Unmarshal(b, &file)

		writeFile(file.Data, folderOut, file.Filename)
	}
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

func receiveData(conn net.Conn) []byte {
	// Receive Size information to set the size  of the buffer
	var buff [32]byte
	_, err := conn.Read(buff[0:])
	if err != nil {
		println("Error: Could not Read size of incomeing INT: ", err.Error())
	}

	// Create size int and have it set so we can read the
	// next bunch of packets.
	var size int32
	b := bytes.NewReader(buff[0:])
	err2 := binary.Read(b, binary.LittleEndian, &size)
	if err2 != nil {
		println("Error: Could Create size int: ", err2.Error())
	}

	// Set the size of the buffer and prepare to read data in
	buff2 := make([]byte, size)
	_, err4 := conn.Read(buff2[0:])
	if err4 != nil {
		println("Error: Could not Read Data into Buffer: ", err4.Error())
	}

	return buff2[:]
}

func readDir(folder string) []os.FileInfo {
	arr, err := ioutil.ReadDir(folder)

	if err != nil {
		println("readDir Error: ", err.Error())
	}
	return arr
}

func readFile(folder string, filename string) []byte {
	b, err := ioutil.ReadFile(folder + "\\" + filename)

	if err != nil {
		println("readFile Error: ", err.Error())
	}

	return b
}

func writeFile(content []byte, folder string, filename string) {
	err := ioutil.WriteFile(folder+"\\"+filename, content, 0777)

	if err != nil {
		print("writeFile Error: ", err.Error())
	}
}

func deleteFile(folder string, filename string) {
	err := os.Remove(folder + "\\" + filename)

	if err != nil {
		println("deleteFile Error: ", err.Error())
	}
}
