// ioFile.go
// dir dirOut fileName fileNames txt file files word char words chars str strs
package wordstat

import (
	"io/ioutil"
	"os"
	"strings"
)

//MSI - map string int - for ioSSR
type MSI map[string]int

// SSR - slice slice rune - for ioSSR
// type SSR [][]rune

// SR - slice rune - for ioSSR
type SR []rune

// SS - slice string
type SS []string

//SSS - slice slice string
type SSS [][]string

//SSB - slice slice byte
// type SSB [][]byte

// SB - slice byte
// type SB []byte

//fnNo- function Name and Oders. - for all
func checkErr(err error, fnNo string) {
	if err != nil {
		panic(fnNo + " :>>>>" + err.Error())
	}
}

func aFreadDir(dir string) (filePathS string) {
	temp, err := ioutil.ReadDir(dir)
	checkErr(err, "aFreadDir 1")
	for _, v := range temp {
		filePathS += dir + "/" + v.Name() + "\r\n"
	}
	return filePathS
}

// bFreadDir - for the bSstat
func bFreadDir(filePathS string) (files SS, fileNameS string) {
	files, _ = cFreadDir(filePathS)
	fileNameS = filePathS
	return files, fileNameS
}

// cFreadDir
func cFreadDir(filePathS string) (files, fileNameSs SS) {
	fileNameSs = strings.Split(filePathS, "\r\n")
	for _, f := range fileNameSs {
		if f == "" {
			continue
		}
		temp, err := ioutil.ReadFile(f)
		checkErr(err, "cFreadDir 1")
		files = append(files, string(temp))
	}
	return files, fileNameSs
}

//rename
func hFtoRnFileS(filePathS, partOld, partNew string) {
	for _, f := range strings.Split(filePathS, "\r\n") {
		if f == "" {
			continue
		}
		os.Rename(f, strings.Replace(f, partOld, partNew, 1))
	}
}
func iFtoRmDir(dirOut string) {
	os.RemoveAll(dirOut)
}

// jFwrite - defer f.Close()  after jFwrite()
func jFwrite(dirOut, fileName, fnName string) (f *os.File) {
	os.MkdirAll(dirOut, 0777)
	f, err := os.OpenFile(dirOut+"/"+fileName, os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0777)
	checkErr(err, fnName)
	return f
}
func (strs SS) jFwriteTxt(dirOut, fileName string) {
	f := jFwrite(dirOut, fileName, "jFwriteTxt 1")
	defer f.Close()
	for _, v := range strs {
		f.WriteString(v + "\r\n")
	}
}
func kFwriteStat(dirOut, fileName, statResault string) {
	f := jFwrite(dirOut, fileName, "kFwriteStat 1")
	f.WriteString("\r\n" + statResault + "\r\n")
	f.Close()
}
