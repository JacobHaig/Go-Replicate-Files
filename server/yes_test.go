package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/*
func TestBytesToInt(t *testing.T) {

	str := "hello"
	str2 := []byte(str)

	val := int32(568)
	val2 := new(bytes.Buffer)

	err := binary.Write(val2, binary.LittleEndian, val)

	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	t.Log(str, "=>", str2)
	t.Log(val, "=>", val2.Bytes())
}*/

func rreadDir(folder string) []os.FileInfo {
	arr, err := ioutil.ReadDir(folder)

	if err != nil {
		println("Error: ", err.Error())
	}
	return arr
}

func WalkAllFilesInDir(root string) ([]string, error) {
	var listFiles []string

	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			listFiles = append(listFiles, path)
			return nil
		})

	if err != nil {
		panic(err)
	}

	return listFiles, nil
}

func TestDirFiles(t *testing.T) {
	root := `C:\Users\jacob\OneDrive\Desktop\Move From`
	dirs, _ := WalkAllFilesInDir(root)

	// Clean up the Dirs so that they do not include the rootfolder
	for e, i := range dirs {
		dirs[e] = strings.ReplaceAll(i, root, "")
	}
}
