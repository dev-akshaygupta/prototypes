package main

import (
	"fmt"

	"github.com/spaolacci/murmur3"
)

// var objects []string

type BloomFilter struct {
	bitSizeArr []int8
	size       uint
	hashes     uint
}

func NewBloomFilter(size uint, hashes uint) *BloomFilter {
	return &BloomFilter{
		bitSizeArr: make([]int8, size),
		size:       size,
		hashes:     hashes,
	}
}

func (bf *BloomFilter) addObject(obj string) {
	for range bf.hashes {
		objByte := []byte(obj)
		idx := int(murmur3.Sum32(objByte)) % int(bf.size)
		bf.bitSizeArr[idx] = 1
	}
}

func (bf *BloomFilter) checkObject(obj string) bool {
	for range bf.hashes {
		objByte := []byte(obj)
		idx := int(murmur3.Sum32(objByte)) % int(bf.size)
		if bf.bitSizeArr[idx] == 0 {
			return false
		}
	}
	return true
}

// func readFile(fileName string, bf *BloomFilter) []string {
// 	file, err := os.Open(fileName)
// 	if err != nil {
// 		fmt.Printf("Unable to open file: %s\n", err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)

// 	for scanner.Scan() {
// 		line := strings.ToLower((strings.TrimSpace(scanner.Text())))
// 		if line == "" {
// 			continue
// 		}

// 		bf.addObject(line)
// 		objects = append(objects, line)
// 	}

// 	return objects
// }

func main() {
	bf := NewBloomFilter(60, 7)

	// Initialize Object List and Bit Size Array
	// objectList := readFile("objectList.txt", bf)
	// fmt.Println(objectList)
	// fmt.Println()
	// fmt.Println(bf.bitSizeArr)

	var choice int8
	var obj string

outer:
	for {
		fmt.Print("Select from following options:\n 1: Check object\n 2: Add object\n 3: Exit\n")
		fmt.Println("Note: Enter \"back\" to switch between the choices!")
		fmt.Print("Choice: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Println("==== Checking object ====")
			for {
				fmt.Print("Enter object: ")
				fmt.Scanln(&obj)
				if obj == "back" {
					break
				}
				if bf.checkObject(obj) {
					fmt.Printf("%s is in the list!\n\n", obj)
				} else {
					fmt.Printf("%s not found!\n\n", obj)
				}
			}
		case 2:
			fmt.Println("==== Adding object ====")
			for {
				fmt.Print("Enter object: ")
				fmt.Scanln(&obj)
				if obj == "back" {
					break
				}
				bf.addObject(obj)
				// objectList = append(objectList, obj)
			}
		case 3:
			fmt.Println("Ok...Bye!!!")
			break outer
		default:
			fmt.Println("Enter a valid option!")
		}
	}
}
