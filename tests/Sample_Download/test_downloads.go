package main

import (
    "fmt"
    "time"
)

func DL_Song(c chan string) {
    for s := range(c){
		fmt.Println("\t~ Downloading song...",s)
		time.Sleep(500 * time.Millisecond)
		fmt.Println("\t~ Done!", s)
    }
    fmt.Println("closed")
    close(c)
}

func DL_PL(pl []string, c chan string, name string) {
    fmt.Println("Download PL", name)
    for _, song := range pl {
		fmt.Println("\t* Queue song", song, len(c))
        c <- song
        //time.Sleep(100 * time.Millisecond)
    }
    fmt.Println("***Async***")
}

func main() {
    PL0 := []string{"A", "B", "C", "D"}
    PL1 := []string{"E", "F"}
    PL2 := []string{"G", "H", "I"}
	PL3 := []string{}
    c := make(chan string, 10)
    
    go DL_Song(c)
    DL_PL(PL0,c,"0")
    time.Sleep(2 * time.Second)
    DL_PL(PL1,c,"1")
    DL_PL(PL2,c,"2")
	DL_PL(PL3,c,"3")
	fmt.Println("----------------------------")
    fmt.Println("No More Download Requests...")
	fmt.Println("----------------------------")
	time.Sleep(5 * time.Second)
	fmt.Println("Program End")
}