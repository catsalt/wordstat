package stat

import (
	"fmt"
	"testing"
	"time"
)

var dir = "F:/inputX Test"
var dirA = "F:/inputX Test/BncCoca25kv3 HeadWords"
var dirB = "F:/inputX Test/BncCoca25kv4"
var filePathsA = "F:/inputX Test/Data Structures & Algorithms.txt"
var filePathsB = "F:/inputX Test/tidy Academic 0-5k.txt"
var dirOut = "F:/outX Test"
var fileA = "F:/inputX Test/Trader.txt"
var fileB = "F:/inputX Test/Cambridge 6749 Word List.txt\r\n"

func TestSS(t *testing.T) {
	fmt.Println("Testing SS.go------------------------1")
	start := time.Now()

	// hFtoRnFileS(aFreadDir(dirA), "HeadWords 25k basewrd", "25k HeadWords ")
	// ZsCompare(fileA, fileB, dirOut, false, true)
	ZsGrade(fileA, fileB, dirOut, true, false)
	// ZsTidy(fileB, dirOut, false, true)
	fmt.Println(time.Since(start).String())
}

func TestFF(t *testing.T) {
	start := time.Now()
	fmt.Println("Testing SS.go------------------------2")
	// iFtoRmDir(dirOut)
	// fileA := "F:/inputX Test/CambridgeEnglishVocabularyProfile_v8.txt"
	// fileB := "F:/inputX Test/Academic 810 finished.txt"
	// fileB := "F:/inputX Test/Cambridge 6749 bEpos Relate verb Head 739.txt"
	// exfile := "F:/inputX Test/Cambridge 298.txt"
	// cambridgeA(camprofile, exfile, dirOut)
	// cambridgeB(camprofile, exfile, dirOut)
	// cambridgeC(camprofile, exfile, dirOut)
	// bEtidyCam(fileA, fileB, dirOut)
	fmt.Println(time.Now().Sub(start).String())
}
