package main

import (
	"os"
	"fmt"
	"bufio"
	"io/ioutil"
)

func main() {
	folderpath := "../data"
	filename := "sample.txt"
	filepath := fmt.Sprintf("%s/%s",folderpath, filename)
	
	missingfilepath := fmt.Sprintf("%s/missingfile.txt",folderpath)
	checkExist(missingfilepath)
	checkExist("../data2")

	fmt.Println()
	listFiles(folderpath)
	fmt.Println()
	readFile(filepath)
}

func checkExist(path string) (bool, error){
	_, err := os.Stat(path)
	if err == nil {
		fmt.Printf("Path %s exists!\n", path)
		return true, nil
	} else if os.IsNotExist(err) {
		fmt.Printf("Path %s does not exist...\n", path)
	}
	return false, err
}

func listFiles(dirpath string) error {
	checkExist(dirpath)
	fmt.Println("Listing files under", dirpath)
	fileInfos, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return err
	}
	for _, value := range fileInfos {
		fmt.Println(value.Name())
	}
	return nil
}

func readFile(filepath string) {
	checkExist(filepath)
	fmt.Printf("Reading file [%s]:\n", filepath)
	ff, _ := os.Open(filepath) 
    f := bufio.NewReader(ff) 
	c := 0
    for { 
        read_line, _ := f.ReadString('\n') 
		if (read_line == "") { break } 
        fmt.Printf("%d:%s",c,read_line) 
		c++
    } 
    ff.Close()
}
