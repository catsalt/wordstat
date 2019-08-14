// corE.go
package filestat

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

func aEreadCam(filePath string) (strs SS) {
	temp, _ := ioutil.ReadFile(filePath)
	strs = strings.Split(string(temp), "</>\r\n")
	return strs
}
func aEreadFile(filePath string) (strs SS) {
	temp, _ := ioutil.ReadFile(filePath)
	strs = strings.Split(string(temp), "\r\n")
	return strs
}
func bEtidyCam(filePathA, filePathB, dirOut string) {
	strs := aEreadCam(filePathA)
	clude := aEreadFile(filePathB)
	// fmt.Println(len(strs))
	// bErelate(strs, dirOut)
	include, exclude := bEclude(strs, clude)
	bEgroup(include, dirOut)
	fmt.Println(len(include), len(exclude))
}

//bEclude - 1st step, included or excluded part in Cambridge file;
func bEclude(strs, clude SS) (include, exclude SS) {
	for _, v := range strs {
		if ok, id := cEindex(v, "\r\n"); ok {
			if cEhas(v[:id], clude) {
				include = append(include, v)
			} else {
				exclude = append(exclude, v)
			}
		}
	}
	return include, exclude
}

//bEgroup - 2nd step, Relate and Ioslate of clude part;
func bEgroup(strs SS, dirOut string) {
	relate, isolate := bEcludeRelate(strs)
	relate.jFwriteTxt(dirOut, fmt.Sprintf("Cambridge 6749 Part Relate %d.txt", len(relate)))
	isolate.jFwriteTxt(dirOut, fmt.Sprintf("Cambridge 6749 Part Isolate %d.txt", len(isolate)))
	bEpos(relate, dirOut, "Relate")
	fmt.Println("---------------")
	bEpos(isolate, dirOut, "Isolate")
}
func bEcludeRelate(strs SS) (include, exclude SS) {
	for _, v := range strs {
		if v == "" {
			continue
		}
		r := false
		for _, p := range []string{"Verbs:", "Adjectives:", "Nouns:", "Adverbs:"} {
		// for _, p := range []string{"Nouns:"} {
			if strings.Count(v, p) > 0 {
				r = true
				if strings.Count(v, p) > 1 {
					fmt.Println(v[:strings.Index(v, "\r\n")], " bErelate More than one ", p)
					return
				}
			}
		}
		if r {
			include = append(include, v+"</>")
		} else {
			exclude = append(exclude, v+"</>")
		}
	}
	return include, exclude
}

//bEpos - 3rd step, Print isolate Head, relate Head,Group,List(Head+Group);
func bEpos(strs SS, dirOut, late string) {
	fn := "bEpos"
	for _, p := range []string{
		"verb", "adjective", "noun", "adverb",
		"auxiliary verb", "modal verb", "phrasal verb",
		"pronoun", "ordinal number", "number", "infinitive marker",
		"determiner", "preposition", "conjunction", "exclamation"} {
		partA, partB := bEcludePos(strs, p)
		if len(partA) == 0 {
			continue
		}
		if late == "Isolate" {
			iso := bEisolate(partA, dirOut)
			iso.jFwriteTxt(dirOut, fmt.Sprintf("Cambridge 6749 %s %s %s Head %d.txt", fn, late, p, len(iso)))
			strs = partB
			continue
		}
		headS, relateS, relate := bErelate(partA, dirOut)
		_, partC := bEclude(partB, relate)
		headS.jFwriteTxt(dirOut, fmt.Sprintf("Cambridge 6749 %s %s %s Head %d.txt", fn, late, p, len(headS)))
		relateS.jFwriteTxt(dirOut, fmt.Sprintf("Cambridge 6749 %s %s %s Group %d.txt", fn, late, p, len(relateS)))
		relate.jFwriteTxt(dirOut, fmt.Sprintf("Cambridge 6749 %s %s %s List %d.txt", fn, late, p, len(relate)))
		strs = partC
	}
}
func bEcludePos(strs SS, pos string) (include, exclude SS) {
	for _, v := range strs {
		if strings.Count(v, "▶ <b>"+pos+"<") > 0 {
			include = append(include, v)
			if strings.Count(v, "▶ <b>"+pos+"<") > 1 {
				fmt.Println(v[:strings.Index(v, "\r\n")], " bEposGroup More than one ", pos)
			}
		} else {
			exclude = append(exclude, v)
		}
	}
	return include, exclude
}
func bEisolate(strs SS, dirOut string) (headS SS) {
	for _, v := range strs {
		if strings.Index(v, "\r\n") == -1 {
			fmt.Printf("bEisolate (%s) has no head or relate.", v)
			continue
		}
		headS = append(headS, v[:strings.Index(v, "\r\n")])
	}
	return headS
}
func bErelate(strs SS, dirOut string) (headS, relateS, x SS) {
	fn := "cErelate"
	for _, v := range strs {
		if strings.Index(v, "\r\n") == -1 || strings.Index(v, "<br>▶") == -1 {
			fmt.Printf("%s: (%s) has no head or relate.", fn, v)
			continue
		}
		head := v[:strings.Index(v, "\r\n")]
		temp := v[:strings.Index(v, "<br>▶")]
		if strings.LastIndex(temp, "<br>") != -1 {
			temp = temp[strings.LastIndex(temp, "<br>")+4:]
			temp = cErmtag(temp, '<', '>')
			for _, v := range []string{"Nouns:", "Verbs:", "Adjectives:", "Adverbs:"} {
				temp = strings.ReplaceAll(temp, v, ",")
			}
			temp = cErmSpace(temp)
		} else {
			fmt.Printf("%s: %s has no relate", fn, head)
		}
		headS = append(headS, head)
		relateS = cEappendNo(relateS, temp)
		for _, v := range strings.Split(temp, ",") {
			x = cEappendNo(x, v)
		}
	}
	sort.SliceStable(relateS, func(i, j int) bool {
		return relateS[i] < relateS[j]
	})
	return headS, relateS, x
}
func cEappendNo(strs SS, str string) SS {
	if str == "" {
		return strs
	}
	if !cEhas(str, strs) {
		strs = append(strs, str)
	}
	return strs
}
func cEhas(str string, strs SS) (ok bool) {
	for _, v := range strs {
		if v == str {
			ok = true
			break
		}
	}
	return ok
}
func cErmSpace(strA string) (strB string) {
	for _, v := range strA {
		if v != rune(' ') {
			strB += string(v)
		}
	}
	return strB
}
func cErmtag(strA string, a, b byte) (strB string) {
	add := true
	var str []rune
	for _, v := range strA {
		if v == rune(a) {
			str = append(str, rune(' '))
			add = false
		}
		if add {
			str = append(str, v)
		}
		if v == rune(b) {
			add = true
		}
	}
	strB = string(str)
	return strB
}
func cEindex(str, substr string) (ok bool, id int) {
	if strings.Index(str, substr) == -1 {
		fmt.Printf("cEchkId (%s) has no %s.", str, substr)
		ok = false
	} else {
		ok = true
		id = strings.Index(str, substr)
	}
	return ok, id
}

func bElistPos(strsA SS, dirOut string) {
	var item SS
	for _, v := range strsA {
		item = append(item, bElistPosA(v))
	}
	item.jFwriteTxt(dirOut, fmt.Sprintf("Cambridge 6749 bElistPos %d.txt", len(item)))
}
func bElistPosA(strA string) (strB string) {
	var str []rune
	add := false
	for _, v := range strA {
		if v == rune('▶') {
			add = true
		}
		if add {
			str = append(str, v)
		}
		if v == rune('/') {
			add = false
		}
	}
	strB = strings.ReplaceAll(string(str), "▶ <b>", "")
	strB = strings.ReplaceAll(strB, "</", ",")
	return strB
}
func bEsumPos(strsA SS) (num []int, total int) {
	num = make([]int, 5)
	for _, v := range strsA {
		n := 0
		for _, p := range []string{"pronoun", "noun", "ordinal number", "number", "infinitive marker",
			"auxiliary verb", "modal verb", "phrasal verb", "adverb", "verb",
			"adjective", "determiner", "preposition", "conjunction", "exclamation"} {
			n += strings.Count(v, "▶ <b>"+p+"<")
		}
		switch n {
		case 0:
			fmt.Println(v[:strings.Index(v, "\r\n")], n)
		case 1:
			num[0]++
		case 2:
			num[1]++
		case 3:
			num[2]++
		case 4:
			num[3]++
		default:
			fmt.Println(v[:strings.Index(v, "\r\n")], n)
			num[4]++
		}
	}
	for k, v := range num {
		total += v * (k + 1)
	}
	fmt.Printf("Cambridge 6749 bEsumPos total %d, \r\n", total)
	fmt.Println(num)
	return num, total
}
func bEsumPosA(strsA SS) (num []int, total int) {
	for _, p := range []string{"pronoun", "noun", "ordinal number", "number", "infinitive marker",
		"auxiliary verb", "modal verb", "phrasal verb", "adverb", "verb",
		"adjective", "determiner", "preposition", "conjunction", "exclamation"} {
		n := 0
		for _, v := range strsA {
			n += strings.Count(v, "▶ <b>"+p+"<")
		}
		num = append(num, n)
	}
	for _, v := range num {
		total += v
	}
	fmt.Printf("Cambridge 6749 bEsumPosA total %d, \r\n", total)
	fmt.Println(num)
	return num, total
}
func bEsumGroup(strsA SS) (num []int, total int) {
	num = make([]int, 6)
	for _, v := range strsA {
		n := 0
		for _, p := range []string{"Verbs:", "Adjectives:", "Nouns:", "Adverbs:"} {
			if strings.Contains(v, p) {
				n++
			}
			if strings.Count(v, p) > 1 {
				fmt.Println(v[:strings.Index(v, "\r\n")])
			}
		}
		switch n {
		case 0:
			num[0]++
		case 1:
			num[1]++
		case 2:
			num[2]++
		case 3:
			num[3]++
		case 4:
			num[4]++
		default:
			fmt.Println(v[:strings.Index(v, "\r\n")], n)
			num[5]++
		}
	}
	for k, v := range num {
		total += v * (k + 1)
	}
	fmt.Printf("Cambridge 6749 bEsumGroup total %d, \r\n", total)
	fmt.Println(num)
	return num, total
}

// for BNC COCA 25k tidy Head words
// - tidy -
func bncCoca(dir, dirOut string) {
	fmt.Println("Hello World!")
	temp, err := ioutil.ReadDir(dir)
	iFtoRmDir(dirOut)
	checkErr(err, "eSstat 1")
	for _, v := range temp {
		file, err := ioutil.ReadFile(dir + "/" + v.Name())
		checkErr(err, "eSstat 2")
		// - tidy original file
		// str := strings.ReplaceAll(string(file), "0\r\n\t", "")
		// str = strings.ReplaceAll(str, "0", "")
		str := strings.Split(string(file), "\r\n")
		var strs strings.Builder
		for _, v := range str {
			// s := strings.SplitN(v, 2)[0]
			strs.WriteString(strings.SplitN(v, " ", 2)[0] + "\r\n")
		}
		kFwriteStat(dirOut, "BncCoca HeadWords 25k "+v.Name(), strs.String())
	}
}
