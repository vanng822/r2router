package r2router

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func httpTestHandler(w http.ResponseWriter, r *http.Request, p Params) {}

func TestMain(m *testing.M) {
	fmt.Println("Test starting")

	retCode := m.Run()

	fmt.Println("Test ending")
	os.Exit(retCode)
}
