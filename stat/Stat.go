// ZsCompare and ZsTidy
// statistic the English txt file for English vocabulary learning
// the Frequency of words in one file, the same words of two files.

package wordstat

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type kv struct {
	k string
	v int
}

// aSpick - bSstat - channel MSI concurrency
// pick true aSpickA (tidy), false aSpickB
func aSpick(files SS, wordDic, charDic MSI, pick bool) {
	var wg, wga sync.WaitGroup
	charCh, wordCh := make(chan MSI), make(chan MSI)
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			var wordDic, charDic MSI
			if pick {
				wordDic, charDic = aSpickA(file)
			} else {
				wordDic, charDic = aSpickB(file)
			}
			wordCh <- wordDic
			charCh <- charDic
		}(file)
	}
	wga.Add(3)
	go func() {
		defer wga.Done()
		wg.Wait()
		close(charCh)
		close(wordCh)
	}()
	go func() {
		defer wga.Done()
		for dic := range charCh {
			for k, v := range dic {
				charDic[k] = charDic[k] + v
			}
		}
	}()
	go func() {
		defer wga.Done()
		for dic := range wordCh {
			for k, v := range dic {
				wordDic[k] = wordDic[k] + v
			}
		}
	}()
	wga.Wait()
}
func aSpickA(file string) (wordDic, charDic MSI) {
	wordDic, charDic = make(MSI), make(MSI)
	var temp SR
	inDic := func(temp SR) {
		switch len(temp) {
		case 0:
		case 1:
			charDic[string(temp)]++
		default:
			wordDic[string(temp)]++
		}
	}
	for _, r := range file {
		switch {
		case r > 64 && r < 91:
			temp = append(temp, r)
		case r > 96 && r < 123:
			temp = append(temp, unicode.ToUpper(r))
		case r == '-':
			charDic[string(r)]++
		default:
			charDic[string(r)]++
			inDic(temp)
			temp = nil
		}
	}

	return wordDic, charDic
}
func aSpickB(file string) (wordDic, charDic MSI) {
	wordDic, charDic = make(MSI), make(MSI)
	for _, v := range strings.Split(file, "\r\n") {
		if len(v) > 1 {
			wordDic[strings.ToUpper(v)]++
		} else {
			charDic[v]++
		}
	}
	return wordDic, charDic
}

// bSstat - remove the redundency - use Upper
// pick true aSpickA (tidy), false aSpickB
func bSstat(filePathS string, pick bool) (words, chars SS, statDir string) {
	var wg sync.WaitGroup
	var wordStat, charStat string
	var sWord, sChar []kv
	wordDic, charDic := make(MSI), make(MSI)
	totalChars, totalWords := 0, 0
	files, statDir := bFreadDir(filePathS)
	aSpick(files, wordDic, charDic, pick)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for k, v := range charDic {
			sChar = append(sChar, kv{k, v})
		}
		chars, charStat, totalChars = bSstatA(sChar)
	}()
	for k, v := range wordDic {
		sWord = append(sWord, kv{k, v})
	}
	words, wordStat, totalWords = bSstatA(sWord)
	wg.Wait()
	statDir = bSstatB(statDir, wordStat, charStat, totalWords, len(words), totalChars, len(chars))
	return words, chars, statDir
}
func bSstatA(skv []kv) (strs SS, strStat string, total int) {
	bSsortK(skv)
	bSsortV(skv)
	var str strings.Builder
	str.Grow(20 * len(skv))
	for i, kv := range skv {
		str.WriteString(strconv.Itoa(i) + "\t-\t" + kv.k + "\t-\t" + strconv.Itoa(kv.v) + "\r\n")
		strs = append(strs, kv.k)
		total += kv.v
	}
	return strs, str.String(), total
}
func bSsortK(skv []kv) {
	sort.SliceStable(skv, func(i, j int) bool { return skv[i].k < skv[j].k })
}
func bSsortV(skv []kv) {
	sort.SliceStable(skv, func(i, j int) bool { return skv[i].v < skv[j].v })
}
func bSstatB(fileName, wordStat, charStat string, totalWords, lnWords, totalChars, lnChars int) string {
	fileName += fmt.Sprintf("\r\n%d Words sorted by frequency into %d\r\n%d Chars sorted by frequency into %d\r\n%d Tokens\r\n",
		totalWords, lnWords, totalChars, lnChars, totalWords+totalChars)
	fileName =
		"begin<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\r\n" + fileName +
			"\r\nRank\t-\tWords\t-\tFrequency\r\n" + wordStat +
			"\r\nRank\t-\tChars\t-\tFrequency\r\n" + charStat +
			"<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<over\r\n"
	return fileName
}

// cSstatA - for ZsGrade to read graded files - use Upper - use tidied files
// func cSstatA(filePathS string) (wordss SSS, fileNames SS) {
// 	files, fileNames := cFreadDir(filePathS)
// 	for _, v := range files {
// 		temp := strings.Split(strings.ToUpper(v), "\r\n")
// 		wordss = append(wordss, temp)
// 	}
// 	return wordss, fileNames
// }

// cSstat - aSpickA (tidy) - for ZsGrade
func cSstat(filePathS string, pick bool) (wordss SSS, fileNames SS) {
	// fileNames = strings.Split(filePathS, "\r\n")
	for _, f := range strings.Split(filePathS, "\r\n") {
		if f == "" {
			continue
		}
		temp, _, _ := bSstat(f, pick)
		wordss = append(wordss, temp)
		fileNames = append(fileNames, f)
	}
	return wordss, fileNames
}

// gSsyll - fileNames  begin from syOnes,
// - if false not split. else num words list for dictionary Learning.
func gSsyll(dirOut string, fileNames, strs SS, group bool) (statSy string) {
	var wg sync.WaitGroup
	lf := len(fileNames)
	syll := make([]SS, lf)
	for _, word := range strs {
		if gSvowel(word) >= lf || gSvowel(word) == 0 {
			syll[lf-1] = append(syll[lf-1], word)
		} else {
			syll[gSvowel(word)-1] = append(syll[gSvowel(word)-1], word)
		}
	}
	statSy = fmt.Sprintf("Total %d words Divided by syllable:\r\n", len(strs))
	if group {
		// 100 one group
		num := 100
		for k, v := range fileNames {
			for i, j := 0, 0; i < len(syll[k]); i += num {
				if i+num < len(syll[k]) {
					j = i + num
				} else {
					j = len(syll[k])
				}
				wg.Add(1)
				go func(k, i, j int, v string) {
					defer wg.Done()
					syll[k][i:j].jFwriteTxt(dirOut, fmt.Sprintf("%s %d %d-%d.txt", strings.TrimSuffix(v, ".txt"), len(syll[k]), i, j))
				}(k, i, j, v)
			}
			statSy += fmt.Sprintf("%d syllable words %d\r\n", k+1, len(syll[k]))
		}
	} else {
		for k, v := range fileNames {
			wg.Add(1)
			go func(k int, v string) {
				defer wg.Done()
				syll[k].jFwriteTxt(dirOut, fmt.Sprintf("%s %d.txt", strings.TrimSuffix(v, ".txt"), len(syll[k])))
			}(k, v)
			statSy += fmt.Sprintf("%d syllable words %d\r\n", k, len(syll[k]))
		}
	}
	wg.Wait()
	return statSy
}

func gSvowel(str string) (num int) {
	for _, letter := range str {
		switch letter {
		case 'A', 'E', 'I', 'O', 'U', 'a', 'e', 'i', 'o', 'u':
			num++
		}
	}
	return num
}

// ZsTidy - aSpick
func ZsTidy(filePathS, dirOut string, group, pick bool) {
	start := time.Now()
	iFtoRmDir(dirOut)
	words, _, statResultDir := bSstat(filePathS, pick)
	words.jFwriteTxt(dirOut, "tidyFileA.txt")
	statResultSy := gSsyll(dirOut, SS{"syOnes.txt", "syTwos.txt", "syTris.txt", "syOthers.txt"}, words, group)
	kFwriteStat(dirOut, "statDirA.txt", statResultSy+statResultDir)
	kFwriteStat(dirOut, "timeUsed.txt", time.Since(start).String())
}

// oScompare - must execute after bSsatat or remove the redundency
// - simply compare with qSequalA() in.
func oScompare(wordsA, wordsB SS) (sameB, onlyA, onlyB SS) {
	_, sameB, onlyA, onlyB = qSequal(wordsA, wordsB, qSequalA)
	return
}

// pSgrade -use upper strings,compare with qSequalB() in
func pSgrade(wordsA, wordsB SS) (sameA, sameB, onlyA SS) {
	sameA, sameB, onlyA, _ = qSequal(wordsA, wordsB, qSequalB)
	return
}

// qSequal - result deponds on equal function, qSequalB fix, qSequalA simply equal.
// - sameA :same wordFamily, same as wordB + fix; sameB :same headWord, same as wordB;
// - onlyA,onlyB :differ totally;
func qSequal(wordsA, wordsB SS, equal func(a, b string) bool) (sameA, sameB, onlyA, onlyB SS) {
	var wg, wga, wgb sync.WaitGroup
	sameDicB := make(MSI)
	sameChA, sameChB, aCh := make(chan string, 2), make(chan string, 2), make(chan string, 2)
	for _, a := range wordsA {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			for i, b := range wordsB {
				if strings.EqualFold(a, b) {
					sameChB <- b
					break
				} else if equal(a, b) {
					sameChA <- a
					sameChB <- b
					break
				} else if i == len(wordsB)-1 {
					aCh <- a
				}
			}
		}(a)
	}
	wgb.Add(1)
	go func() {
		defer wgb.Done()
		for v := range sameChB {
			sameDicB[v]++
		}
	}()
	wga.Add(3)
	go func() {
		defer wga.Done()
		for v := range sameChA {
			sameA = append(sameA, v)
		}
	}()
	go func() {
		defer wga.Done()
		for v := range aCh {
			onlyA = append(onlyA, v)
		}
	}()
	go func() {
		defer wga.Done()
		wg.Wait()
		close(sameChA)
		close(sameChB)
		close(aCh)
	}()
	wgb.Wait()
	for _, b := range wordsB {
		if _, ok := sameDicB[b]; ok {
			sameB = append(sameB, b)
		} else {
			onlyB = append(onlyB, b)
		}
	}
	wga.Wait()
	return sameA, sameB, onlyA, onlyB
}
func qSequalA(wordA, wordB string) bool {
	return strings.EqualFold(wordA, wordB)
}

// oSequalB - use upper string
func qSequalB(wordA, wordB string) bool {
	equalSuffix := func(wordA, wordB, suffixB, suffix string) bool {
		return strings.EqualFold(wordA, strings.TrimSuffix(wordB, suffixB)+suffix)
	}
	for _, v := range []string{"S", "ES", "ER", "EST", "ED", "ING", "LY"} {
		if strings.HasSuffix(wordA, v) && len(wordB) > 0 {
			switch {
			case strings.EqualFold(wordA, wordB+v):
				return true
			case v != "S" && v != "ING" && strings.HasSuffix(wordB, "Y"):
				return equalSuffix(wordA, wordB, "Y", "I"+v)
			case v != "S" && v != "ES" && strings.HasSuffix(wordB, "E"):
				switch {
				case v == "ING" && strings.HasSuffix(wordB, "IE"):
					return equalSuffix(wordA, wordB, "IE", "YING")
				case v == "LY" && strings.HasSuffix(wordB, "LE"):
					return equalSuffix(wordA, wordB, "E", "Y")
				default:
					return equalSuffix(wordA, wordB, "E", v)
				}
			case v != "S" && v != "ES" && strings.EqualFold(wordA, wordB+wordB[len(wordB)-1:]+v):
				return true
			}
		}
	}
	return false
}

// ZsCompare - oScompare - stat and then compare
func ZsCompare(filePathsA, filePathsB, dirOut string, fileApick, fileBpick bool) {
	start := time.Now()
	iFtoRmDir(dirOut)
	var wg sync.WaitGroup
	var wordsA, wordsB SS
	var resultA, resultB string
	wg.Add(2)
	go func() {
		defer wg.Done()
		wordsA, _, resultA = bSstat(filePathsA, fileApick)
	}()
	go func() {
		defer wg.Done()
		wordsB, _, resultB = bSstat(filePathsB, fileBpick)
	}()
	wg.Wait()
	sameAB, inA, inB := oScompare(wordsA, wordsB)
	wg.Add(5)
	go func() {
		defer wg.Done()
		wordsA.jFwriteTxt(dirOut, "tidyFileA.txt")
	}()
	go func() {
		defer wg.Done()
		wordsB.jFwriteTxt(dirOut, "tidyFileB.txt")
	}()
	go func() {
		defer wg.Done()
		inA.jFwriteTxt(dirOut, "onlyInA.txt")
	}()
	go func() {
		defer wg.Done()
		inB.jFwriteTxt(dirOut, "onlyInB.txt")
	}()
	go func() {
		defer wg.Done()
		sameAB.jFwriteTxt(dirOut, "sameAB.txt")
	}()
	resultC := "Compared tidied FileA and FileB [-two parts below-]:\r\n"
	resultC += fmt.Sprintf("FileA: wordsA %d,  sameA Words %d, differA Words %d \r\n",
		len(wordsA), len(sameAB), len(inA))
	resultC += fmt.Sprintf("FileB: wordsB %d,  sameB Words %d, differB Words %d \r\n\r\n",
		len(wordsB), len(sameAB), len(inB))
	kFwriteStat(dirOut, "statCompare.txt", resultC+resultA+resultB)
	wg.Wait()
	kFwriteStat(dirOut, "timeUsed.txt", time.Since(start).String())
	// fmt.Println(time.Since(start).String())
}

// ZsGrade - stat A, compare with B (divided already in grades)
// cSstat(dirB) - read grade list
func ZsGrade(filePathsA, filePathsB, dirOut string, fileApick, fileBpick bool) {
	start := time.Now()
	iFtoRmDir(dirOut)
	var wg sync.WaitGroup
	var resultA string
	resultB := "the Selected word files below Compared with BaseWord GradeLists: \r\n\r\n"
	var sameA, sameB, onlyA, fileNames SS
	var wordssB SSS
	wg.Add(2)
	go func() {
		defer wg.Done()
		onlyA, _, resultA = bSstat(filePathsA, fileApick)
	}()
	go func() {
		defer wg.Done()
		wordssB, fileNames = cSstat(filePathsB, fileBpick)
	}()
	wg.Wait()
	for k, v := range wordssB {
		if len(v) == 0 {
			continue
		}
		sameA, sameB, onlyA = pSgrade(onlyA, v)
		wg.Add(2)
		go func(sameB SS, k int) {
			defer wg.Done()
			sameB.jFwriteTxt(dirOut, "Grade "+strconv.Itoa(k)+" base words.txt")
		}(sameB, k)
		go func(sameA SS, k int) {
			defer wg.Done()
			sameA.jFwriteTxt(dirOut, "Grade "+strconv.Itoa(k)+" words family.txt")
		}(sameA, k)
		resultB += fmt.Sprintf("Grade %d: Basewords %d,  Wordsfamily %d, %d || %d, BaseWord of %s \r\n",
			k, len(sameB), len(sameA), len(sameB)+len(sameA), len(v), fileNames[k])
	}
	onlyA.jFwriteTxt(dirOut, "onlyInA.txt")
	resultB += strconv.Itoa(len(onlyA)) + "\tWords not Found in BaseWord GradeLists \r\n"
	kFwriteStat(dirOut, "gradeCompare.txt", resultB+"\r\n"+resultA)
	wg.Wait()
	kFwriteStat(dirOut, "timeUsed.txt", time.Since(start).String())
}
