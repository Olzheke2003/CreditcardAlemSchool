package main

import (
	"bufio"
	"fmt"
	"main/funcs"
	"os"
	"strings"
)

func outputResults(results []bool) {
	exitCode := 0
	for _, result := range results {
		if result {
			fmt.Println("OK")
		} else {
			fmt.Fprintln(os.Stderr, "INCORRECT")
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Incorrect input")
		return
	}
	flag := os.Args[1]

	switch flag {
	case "validate":
		if len(os.Args) > 2 && os.Args[2] == "--stdin" {
			scanner := bufio.NewScanner(os.Stdin)
			var cardNums []string
			for scanner.Scan() {
				cardNums = append(cardNums, strings.Fields(scanner.Text())...)
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
				os.Exit(1)
			}
			res := funcs.Validate(cardNums)
			outputResults(res)
		} else {
			res := funcs.Validate(os.Args[2:])
			outputResults(res)
		}
	case "generate":
		if len(os.Args) > 3 && os.Args[2] == "--pick" {
			funcs.GenerateCardNums(os.Args[3], true)
		} else {
			funcs.GenerateCardNums(os.Args[2], false)
		}
	case "information":
		var brandsFile, issuersFile string
		var cardNums []string

		// Обработка аргументов для информации о картах
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if strings.HasPrefix(arg, "--brands=") {
				brandsFile = strings.TrimPrefix(arg, "--brands=")
			} else if strings.HasPrefix(arg, "--issuers=") {
				issuersFile = strings.TrimPrefix(arg, "--issuers=")
			} else if arg == "--stdin" {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					cardNums = append(cardNums, strings.Fields(scanner.Text())...)
				}
				if err := scanner.Err(); err != nil {
					fmt.Fprintln(os.Stderr, "Error reading input:", err)
					os.Exit(1)
				}
			} else {
				cardNums = append(cardNums, arg)
			}
		}

		if brandsFile == "" || issuersFile == "" || len(cardNums) == 0 {
			fmt.Fprintln(os.Stderr, "Error: you must specify brand and issuer files, and card numbers.")
			os.Exit(1)
		}

		// Загрузка брендов и эмитентов
		brands, err := funcs.LoadBrands(brandsFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error loading brands:", err)
			os.Exit(1)
		}

		issuers, err := funcs.LoadIssuers(issuersFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error loading issuers:", err)
			os.Exit(1)
		}

		// Вывод информации по каждому номеру карты
		for _, cardNum := range cardNums {
			info := funcs.GetCardInformation(cardNum, brands, issuers)
			fmt.Println(cardNum)
			if funcs.Luhn(cardNum) {
				fmt.Println("Correct: yes")
			} else {
				fmt.Println("Correct: no")
			}
			fmt.Println("Card Brand:", info.Brand)
			fmt.Println("Card Issuer:", info.Issuer)
		}
	case "issue":
		var brandsFile, issuersFile, brand, issuer string

		// Обработка аргументов для выпуска карты
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if strings.HasPrefix(arg, "--brands=") {
				brandsFile = strings.TrimPrefix(arg, "--brands=")
			} else if strings.HasPrefix(arg, "--issuers=") {
				issuersFile = strings.TrimPrefix(arg, "--issuers=")
			} else if strings.HasPrefix(arg, "--brand=") {
				brand = strings.TrimPrefix(arg, "--brand=")
			} else if strings.HasPrefix(arg, "--issuer=") {
				issuer = strings.TrimPrefix(arg, "--issuer=")
			}
		}

		if brandsFile == "" || issuersFile == "" || brand == "" || issuer == "" {
			fmt.Fprintln(os.Stderr, "Error: you must specify brand and issuer files, brand, and issuer.")
			os.Exit(1)
		}

		// Загрузка брендов и эмитентов
		brands, err := funcs.LoadBrands(brandsFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error loading brands:", err)
			os.Exit(1)
		}

		issuers, err := funcs.LoadIssuers(issuersFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error loading issuers:", err)
			os.Exit(1)
		}

		// Проверка наличия бренда
		prefix, brandExists := brands[brand]
		if !brandExists {
			fmt.Fprintln(os.Stderr, "Error: brand not found")
			os.Exit(1)
		}

		// Проверка наличия эмитента
		_, issuerExists := issuers[issuer]
		if !issuerExists {
			fmt.Fprintln(os.Stderr, "Error: issuer not found")
			os.Exit(1)
		}

		// Генерация номера карты
		funcs.GenerateCardNums(prefix, false) // нет необходимости захватывать возвращаемое значение

	default:
		fmt.Fprintln(os.Stderr, "ERROR: unknown flag")
		os.Exit(1)
	}
}
