package helper

import (
	"fmt"
)

func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
