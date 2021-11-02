/* ASCII ART YNOV INFORMATIQUE 2020 */
/* Copyright INGREMEAU, CLAMADIEU-THARAUD, QUESNOY 2020 */

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	var args []string = os.Args[1:]
	var argTable []string
	if len(args) < 1 {
		fmt.Println("ERROR : No arguments.")
		return
	}
	terminalWidth := terminalWidthF()
	start := 0
	words := []rune(args[0])
	var wordList []string
	for i := 0; i < len(words); i++ {
		if words[i] == 92 && words[i+1] == 110 {
			wordList = append(wordList, string(words[start:i]))
			start = i + 2
			i += 2
		}
	}
	tooLong := false
	wordList = append(wordList, string(words[start:]))
	for !tooLong {
		wordList, tooLong = verifLen(wordList, int(terminalWidth))
	}
	police := "standard.txt"
	align := "left"
	mode := -1
	argTable = findArgument(args, argTable, &police)

	for i := 0; i < len(wordList); i++ {
		tabColor := strings.Split(wordList[i], "")
		verifArgument(argTable, &tabColor, &police, &mode, &align)
		word := []rune(wordList[i])
		result := fillArray(word, openFiles(police))
		printResult(result, argTable, word, mode, align, tabColor)
	}
}

func terminalWidthF() uint {
	terminalWidth, err := Width()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return terminalWidth
}

func fileExist(name string, police *string) bool {
	fileName := name + ".txt"
	// If file exist
	fileExist, err := os.Stat(fileName)
	if err != nil && fileExist == nil {
		/*if os.IsNotExist(err) {
			return false
		}*/
		return false
	}
	*police = fileName
	return true

}

func fileExist2(name string) bool {
	// If file exist
	fileExist, err := os.Stat(name)
	if err != nil && fileExist == nil {
		return false
	}
	return true

}

func openFiles(filename string) []string {
	// Open file
	file, err := os.Open(filename)
	if err != nil && file == nil {
		fmt.Println(string("ERROR OPENING FILE"))
		os.Exit(0)
	}
	var vals []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		vals = append(vals, scanner.Text())
	}
	return vals
}

func fillArray(word []rune, vals []string) [][]string {
	//Fill array
	var lines = make([][]int, 8)
	for i := range lines {
		lines[i] = make([]int, len(word))
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < len(word); j++ {
			lines[i][j] = findLine(word, j, i)
		}
	}

	var result = make([][]string, 8)
	for i := range result {
		result[i] = make([]string, len(word))
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < len(word); j++ {
			result[i][j] = vals[lines[i][j]-1]
		}
	}
	return result
}

func findLine(word []rune, nb int, line int) int {
	return int((word[nb]-31)*8+(word[nb]-31)-6) + line - 1
}

func isInArray(array [][]int, value int, word []rune) (bool, int, int) {
	for i := 0; i < 8; i++ {
		for j := 0; j < len(word); j++ {
			if int(array[i][j]) == value {
				return true, i, j
			}
		}
	}
	return false, 0, 0
}

func printResult(array [][]string, argTable []string, word []rune, mode int, align string, tabColor []string) {
	// Mode == -1 -> Print in console
	// Else -> Print in file
	if mode == -1 && align == "left" {
		for i := 0; i < 8; i++ {
			for j := 0; j < len(word); j++ {
				if len(tabColor[j]) > 1 {
					fmt.Print(tabColor[j], array[i][j], "\033[0m")
				} else {
					fmt.Print(array[i][j])
				}
				if j == len(word)-1 {
					fmt.Print("\n")
				}
			}
		}
	} else if mode >= 0 {
		// If file exist
		fileExist, err := os.Stat(argTable[mode+1])
		if err != nil && fileExist == nil {
			if os.IsNotExist(err) {
				// Create file
				ff, err := os.Create(argTable[mode+1])
				if err != nil && ff == nil {
					fmt.Println(err)
					return
				}
			}
		}

		// Open file
		f, err := os.OpenFile(argTable[mode+1], os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			fmt.Println(err)
			return
		}

		//Write in file
		for i := 0; i < 8; i++ {
			line := ""
			for j := 0; j < len(word); j++ {
				line += array[i][j]
			}
			writer := bufio.NewWriter(f)
			fmt.Fprintln(writer, line)
			writer.Flush()
		}
	}
	if align == "center" {
		wordlen, spaces := wordLen(array)
		terminalWidth, err := Width()
		spaceNb := (int(terminalWidth) / 2) - (wordlen / 2) + spaces - spaces
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		for i := 0; i < 8; i++ {
			for k := 0; k < spaceNb; k++ {
				fmt.Print(" ")
			}
			for j := 0; j < len(word); j++ {
				if len(tabColor[j]) > 1 {
					fmt.Print(tabColor[j], array[i][j], "\033[0m")
				} else {
					fmt.Print(array[i][j])
				}
				if j == len(word)-1 {
					fmt.Print("\n")
				}
			}
		}
	} else if align == "right" {
		wordlen, spaces := wordLen(array)
		terminalWidth, err := Width()
		spaceNb := int(terminalWidth) - wordlen + spaces - spaces
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		for i := 0; i < 8; i++ {
			for k := 0; k < spaceNb; k++ {
				fmt.Print(" ")
			}
			for j := 0; j < len(word); j++ {
				if len(tabColor[j]) > 1 {
					fmt.Print(tabColor[j], array[i][j], "\033[0m")
				} else {
					fmt.Print(array[i][j])
				}
				if j == len(word)-1 {
					fmt.Print("\n")
				}
			}
		}
	} else if align == "justify" {
		wordlen, spaces := wordLen(array)
		terminalWidth, err := Width()
		spaceNb := 0
		if spaces > 1 {
			spaceNb = (int(terminalWidth) - wordlen) / (spaces)
		} else {
			spaceNb = (int(terminalWidth) - wordlen)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		for i := 0; i < 8; i++ {
			for j := 0; j < len(word); j++ {
				if len(tabColor[j]) > 1 {
					fmt.Print(tabColor[j], array[i][j], "\033[0m")
				} else {
					fmt.Print(array[i][j])
				}
				if word[j] == 32 {
					for k := 0; k < spaceNb; k++ {
						fmt.Print(" ")
					}
				}
				if j == len(word)-1 {
					fmt.Print("\n")
				}
			}
		}
	}
}

func wordLen(array [][]string) (int, int) {
	result := 0
	spaces := 0
	for i := 0; i < len(array[0]); i++ {
		result += len(array[0][i])
		if array[0][i] == "      " {
			spaces++
		}
	}
	return result, spaces
}

// ARGUMENTS FUNCTIONS

func findArgument(args []string, argTable []string, police *string) []string {
	for i := 0; i < len(args); i++ {
		if len(args[i]) > 8 && args[i][:8] == "--color=" {
			argTable = append(argTable, args[i][:8])
			argTable = append(argTable, args[i][8:])
			if i != len(args)-1 {
				argTable = append(argTable, args[i+1])
			}
			argTable = append(argTable, "")
		} else if len(args[i]) > 9 && args[i][:9] == "--output=" {
			argTable = append(argTable, args[i][:9])
			argTable = append(argTable, args[i][9:])
			argTable = append(argTable, "")
			if fileExist2(args[i][9:]) {
				e := os.Remove(args[i][9:])
				if e != nil {
					fmt.Println(e)
					os.Exit(0)
				}
			}
		} else if len(args[i]) > 8 && args[i][:8] == "--align=" {
			argTable = append(argTable, args[i][:8])
			argTable = append(argTable, args[i][8:])
			argTable = append(argTable, "")
		} else if len(args[i]) > 10 && args[i][:10] == "--reverse=" {
			argTable = append(argTable, args[i][:10])
			argTable = append(argTable, args[i][10:])
			argTable = append(argTable, "")
		} else if len(args[i]) > 6 && args[i][:7] == "--help=" {
			argTable = append(argTable, args[i][:7])
			argTable = append(argTable, args[i][7:])
			argTable = append(argTable, "")
		} else {
			if i != 0 {
				fileExist(args[i], police)
			}
			argTable = append(argTable, args[i])
			argTable = append(argTable, "")
		}
	}
	return argTable
}

func verifArgument(argTable []string, tabColor *[]string, police *string, mode *int, align *string) {
	isReversed := false
	index := 0
	for i := 0; i < len(argTable); i++ {
		if argTable[i] == "--color=" {
			color(argTable[i+1], argTable[i+2], tabColor)
		} else if argTable[i] == "--output=" {
			//argTable[i+1] = fichier
			*mode = i
		} else if argTable[i] == "--align=" {
			//argTable[i+1] = paramètre d'alignement
			*align = argTable[i+1]
		} else if argTable[i] == "--reverse=" {
			//argTable[i+1] = fichier
			isReversed = true
			index = i + 1
		} else if argTable[i] == "--help=" {
			help(argTable[i+1])
		} else if argTable[i] == "--help" {
			help(argTable[i+1])
		} else {
			if i != 0 {
				fileExist(argTable[i], police)
			}
		}
	}
	if isReversed {
		asciiReverse(argTable[index], *police, argTable)
	}
}

//END OF ARGUMENTS FUNCTIONS

// COLOR FUNCTIONS

func color(name string, letters string, tabColor *[]string) string {

	switch name {
	case "white":
		changeColor(letters, tabColor, "\033[37m")
	case "black":
		changeColor(letters, tabColor, "\033[30m")
	case "red":
		changeColor(letters, tabColor, "\033[31m")
	case "yellow":
		changeColor(letters, tabColor, "\033[33m")
	case "green":
		changeColor(letters, tabColor, "\033[32m")
	case "blue":
		changeColor(letters, tabColor, "\033[34m")
	case "pink":
		changeColor(letters, tabColor, "\033[35m")
	}
	return ""

}

func changeColor(letters string, tabColor *[]string, color string) {
	if len(letters) > 1 {
		if letters[0:2] == "--" || fileExist(letters, &letters) {
			for j := range *tabColor {
				if len((*tabColor)[j]) < 2 {
					(*tabColor)[j] = color
				}
			}
		}
	} else if len(letters) == 0 {
		for j := range *tabColor {
			if len((*tabColor)[j]) < 2 {
				(*tabColor)[j] = color
			}
		}
	}
	for _, i := range letters {
		for j, letter := range *tabColor {
			if string(letter) == string(i) {
				(*tabColor)[j] = color
			}
		}
	}

}

//END OF COLOR FUNCTIONS

// REVERSE FUNCTIONS

func asciiReverse(fileName string, police string, argTable []string) {
	txt, err := readFile(police)
	if err != nil {
		fmt.Println(err)
		return
	}
	signes := getEachSign(txt)
	fileTODO, err := readFile(fileName)
	if err != nil {
		fmt.Printf("Can`t find file: %s. Please make sure that the file name is right\n", fileName)
		return
	}
	var wordList = make([][]string, (len(fileTODO)-1)/8)
	for i := 0; i < len(wordList); i++ {
		wordList[i] = make([]string, 8)
		for j := 0; j < 8; j++ {
			wordList[i][j] = fileTODO[j+i*8]
		}
	}
	for i := 0; i < len(wordList); i++ {
		result := solver(signes, wordList[i])
		tabColor := strings.Split(result, "")
		align := "left"
		verifArgumentReverse(argTable, &tabColor, &align)
		printReverse(result, tabColor, align)
	}
	os.Exit(0)
}

func printReverse(word string, tabColor []string, align string) {
	if align == "left" {
		for j := 0; j < len(word); j++ {
			if len(tabColor[j]) > 1 {
				fmt.Print(tabColor[j], string(word[j]), "\033[0m")
			} else {
				fmt.Print(string(word[j]))
			}
		}
		println("")
	}
	if align == "center" {
		wordlen, spaces := wordLenReverse(word)
		terminalWidth, err := Width()
		spaceNb := (int(terminalWidth) / 2) - (wordlen / 2) + spaces - spaces
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		for k := 0; k < spaceNb; k++ {
			fmt.Print(" ")
		}
		for j := 0; j < len(word); j++ {
			if len(tabColor[j]) > 1 {
				fmt.Print(tabColor[j], string(word[j]), "\033[0m")
			} else {
				fmt.Print(string(word[j]))
			}
		}
		println("")
	} else if align == "right" {
		wordlen, spaces := wordLenReverse(word)
		terminalWidth, err := Width()
		spaceNb := int(terminalWidth) - wordlen + spaces - spaces
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		for k := 0; k < spaceNb; k++ {
			fmt.Print(" ")
		}
		for j := 0; j < len(word); j++ {
			if len(tabColor[j]) > 1 {
				fmt.Print(tabColor[j], string(word[j]), "\033[0m")
			} else {
				fmt.Print(string(word[j]))
			}
		}
		println("")

	} else if align == "justify" {
		wordlen, spaces := wordLenReverse(word)
		terminalWidth, err := Width()
		spaceNb := 0
		if spaces > 1 {
			spaceNb = (int(terminalWidth) - wordlen) / (spaces)
		} else {
			spaceNb = (int(terminalWidth) - wordlen)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		for j := 0; j < len(word); j++ {
			if len(tabColor[j]) > 1 {
				fmt.Print(tabColor[j], string(word[j]), "\033[0m")
			} else {
				fmt.Print(string(word[j]))
			}
			if word[j] == 32 {
				for k := 0; k < spaceNb; k++ {
					fmt.Print(" ")
				}
			}
			if j == len(word)-1 {
				fmt.Print("\n")
			}
		}
	}
}

func solver(signes [][]string, fileTODO []string) string {
	var result string
	for len(fileTODO[0]) > 0 {
		for i, v := range signes {
			if findIndex(v, fileTODO) {
				result += string(rune(i + 32))
				fileTODO = removeSign(len(v[0]), fileTODO)
			}
		}
	}
	return result
}

func readFile(str string) ([]string, error) {
	bytes, err := ioutil.ReadFile(str)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(bytes), "\n")
	return lines, nil
}

func getEachSign(template []string) [][]string {
	var res [][]string
	for i := 0; i < len(template)-1; i += 9 {
		tmp := template[i+1 : i+9]
		res = append(res, tmp)
	}
	return res
}

func findIndex(sign, word []string) bool {
	if len(sign[0]) > len(word[0]) {
		return false
	}
	for i, v := range word {
		if sign[i] != v[:len(sign[i])] {
			return false
		}
	}
	return true
}

func removeSign(length int, word []string) []string {
	for i, v := range word {
		word[i] = v[length:]
	}
	return word
}

func wordLenReverse(array string) (int, int) {
	spaces := 0
	for i := 0; i < len(array); i++ {
		if array[i] == ' ' {
			spaces++
		}
	}
	return len(array), spaces
}

func verifArgumentReverse(argTable []string, tabColor *[]string, align *string) {
	for i := 0; i < len(argTable); i++ {
		if argTable[i] == "--color=" {
			color(argTable[i+1], argTable[i+2], tabColor)
		} else if argTable[i] == "--align=" {
			//argTable[i+1] = paramètre d'alignement
			*align = argTable[i+1]
		}
	}
}

// END OF REVERSE FUNCTIONS

// TERMINAL SIZE FUNCTIONS

func size() (string, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	return string(out), err
}

func Width() (uint, error) {
	output, err := size()
	if err != nil {
		return 0, err
	}
	_, width, err := parse(output)
	return width, err
}

func parse(input string) (uint, uint, error) {
	parts := strings.Split(input, " ")
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	y, err := strconv.Atoi(strings.Replace(parts[1], "\n", "", 1))
	if err != nil {
		return 0, 0, err
	}
	return uint(x), uint(y), nil
}

// END OF TERMINAL SIZE FUNCTIONS

// BONUS FONCTIONS

func help(arg string) {
	if arg == "color" {
		fmt.Println("")
		fmt.Println("Change the color of the designated letters ")
		fmt.Println("The command is --color== followed by a color and the letters")
		fmt.Println("The different colors are :")
		fmt.Println("white")
		fmt.Println("black")
		fmt.Println("red")
		fmt.Println("yellow")
		fmt.Println("green")
		fmt.Println("blue")
		fmt.Println("pink")
		fmt.Println("")
	} else if arg == "output" {
		fmt.Println("")
		fmt.Println("Writing the result into a file ")
		fmt.Println("The command is --output= followed by the name of the .txt")
		fmt.Println("")
	} else if arg == "output" {
		fmt.Println("")
		fmt.Println("Recovers an ascii text in a .txt file and displays it as normal text.")
		fmt.Println("The command is --reverse= followed by the name of the .txt")
		fmt.Println("")
	} else if arg == "align" {
		fmt.Println("")
		fmt.Println("The representation is formatted using a flag")
		fmt.Println("The command is --align= followed by the flag")
		fmt.Println("The different flag are :")
		fmt.Println("left")
		fmt.Println("right")
		fmt.Println("center")
		fmt.Println("justify")
		fmt.Println("")
	} else if arg == "fs" {
		fmt.Println("")
		fmt.Println("Change the ascii police")
		fmt.Println("It's not a command juste write the name of a file of police, but not with the .txt")
		fmt.Println("")
	} else {
		fmt.Println("")
		fmt.Println("For more information on an element enter --help= followed by the key word designating it \n")
		fmt.Println("color     Change the color of the designated letters ")
		fmt.Println("output    Writing the result into a file ")
		fmt.Println("reverse   Recovers an ascii text in a .txt file and displays it as normal text.")
		fmt.Println("align     The representation is formatted using a flag ")
		fmt.Println("fs        Change the ascii police")
		fmt.Println("")
	}
	os.Exit(0)
}

func verifLen(wordList []string, terminalWidth int) ([]string, bool) {
	//Vérifie la taille du mot et informe l'utilisateur si il est trop grand
	for i := 0; i < len(wordList); i++ {
		if len(wordList[i]) > int(terminalWidth)/8 {
			wordList = insert(wordList, i, "PLACEHOLDER")
			wordList[i] = (wordList[i+1])[:terminalWidth/8]
			wordList[i+1] = (wordList[i+1])[terminalWidth/8:]
			return wordList, false
		}
	}
	return wordList, true
}

func insert(array []string, position int, value string) []string {
	array = append(array, "")
	copy((array)[position+1:], (array)[position:])
	(array)[position] = value
	return array
}

// END BONUS FONCTIONS
