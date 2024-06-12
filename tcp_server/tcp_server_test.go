package tcp_server

import (
	"fmt"
	"testing"
)

func TestSliceIndex(t *testing.T) {
	c := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(c[:2])
}
