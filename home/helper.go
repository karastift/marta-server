package home

import (
	"fmt"
	"time"
)

func Log(str string) {
	fmt.Println("[" + time.Now().Format(time.ANSIC) + "] " + str)
}
