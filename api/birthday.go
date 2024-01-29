package handler

import (
	"fmt"
	"net/http"
	"os"
)

func Birthday(w http.ResponseWriter, r *http.Request) {
	fmt.Println(os.Getenv("EDGE_CONFIG"))
}
