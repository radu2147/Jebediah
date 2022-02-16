package frontend

import (
	"fmt"
	"sync"
	"time"
)

func DisplayFrontend(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Scanning for Windows 11 compatibility...")
	time.Sleep(2 * time.Second)
	fmt.Println("System is not compatible with Windows 11")
}
