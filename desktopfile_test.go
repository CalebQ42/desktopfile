package desktopfile

import (
	"fmt"
	"os"
	"testing"
)

func TestDesktopFile(t *testing.T) {
	testFil, err := os.Open("testing/testing.desktop")
	if err != nil {
		t.Fatal(err)
	}
	fil, err := OpenWithOptions(testFil, Options{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fil.DefaultGroup().GetEntry("Name").GetValue())
	t.Fatal("HI")
}
