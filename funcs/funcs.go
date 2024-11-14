package funcs

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func Luhn(cardNum string) bool {
	sum := 0
	alt := false

	for i := len(cardNum) - 1; i >= 0; i-- {
		n, err := strconv.Atoi(string(cardNum[i]))
		if err != nil {
			return false
		}
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}

func Validate(cardNums []string) []bool {
	res := make([]bool, len(cardNums))
	for i, num := range cardNums {
		if len(num) >= 13 && Luhn(num) {
			res[i] = true
		} else {
			res[i] = false
		}
	}
	return res
}

func GenerateCardNums(cardNum string, pick bool) {
	cardNum = strings.TrimSpace(cardNum)
	fmt.Println("Received card template:", cardNum)

	starsCount := len(cardNum) - len(strings.TrimRight(cardNum, "*"))
	if starsCount != 4 {
		fmt.Fprintln(os.Stderr, "Error: the number of asterisks at the end must be exactly 4.")
		os.Exit(1)
	}

	if !strings.HasSuffix(cardNum, "****") {
		fmt.Fprintln(os.Stderr, "Error: asterisks must be at the end.")
		os.Exit(1)
	}

	prefix := strings.TrimRight(cardNum, "*")
	fmt.Println("Card prefix:", prefix)

	var generatedNums []string
	for i := 0; i < 10000; i++ {
		generatedNum := fmt.Sprintf("%s%04d", prefix, i)
		if Luhn(generatedNum) {
			generatedNums = append(generatedNums, generatedNum)
		}
	}

	if pick {
		if len(generatedNums) > 0 {
			rand.Seed(time.Now().UnixNano())
			randomIndex := rand.Intn(len(generatedNums))
			fmt.Println("Randomly generated number:", generatedNums[randomIndex])
		} else {
			fmt.Println("No valid numbers generated.")
		}
	} else {
		for _, num := range generatedNums {
			fmt.Println(num)
		}
	}
}

func LoadBrands(filename string) (map[string]string, error) {
	brands := make(map[string]string)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			brand := strings.TrimSpace(parts[0])
			prefix := strings.TrimSpace(parts[1])
			fmt.Printf("Загружен бренд: %s с префиксом: %s\n", brand, prefix)
			brands[brand] = prefix
		}
	}

	return brands, scanner.Err()
}

func LoadIssuers(filename string) (map[string][]string, error) {
	issuers := make(map[string][]string)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			issuer := strings.TrimSpace(parts[0])
			prefix := strings.TrimSpace(parts[1])
			fmt.Printf("Загружен эмитент: %s с префиксом: %s\n", issuer, prefix)
			issuers[issuer] = append(issuers[issuer], prefix)
		}
	}

	return issuers, scanner.Err()
}

type CardInfo struct {
	Brand  string
	Issuer string
}

func GetCardInformation(cardNum string, brands map[string]string, issuers map[string][]string) (info CardInfo) {
	if !Luhn(cardNum) {
		info.Brand = "-"
		info.Issuer = "-"
		return
	}

	// Проверка на бренд
	for brand, prefix := range brands {
		fmt.Printf("Checking brand: %s with prefix: %s for card: %s\n", brand, prefix, cardNum)
		if strings.HasPrefix(cardNum, prefix) {
			info.Brand = brand
			break
		}
	}

	if info.Brand == "" {
		info.Brand = "-"
	}

	// Проверка на эмитента
	var longestPrefix string
	for issuer, prefixes := range issuers {
		for _, prefix := range prefixes {
			fmt.Printf("Checking issuer: %s with prefix: %s for card: %s\n", issuer, prefix, cardNum)
			if strings.HasPrefix(cardNum, prefix) && len(prefix) > len(longestPrefix) {
				longestPrefix = prefix
				info.Issuer = issuer // Сохраняем только имя эмитента без префикса
			}
		}
	}

	if info.Issuer == "" {
		info.Issuer = "-"
	}

	return
}
