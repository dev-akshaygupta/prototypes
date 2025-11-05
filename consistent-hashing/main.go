package main

import (
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
Problem Statement: Implement Consistent Hashing
  - Take user object as input
  - Assign server to a user data - map of users with keys as servers
    . random
    . hashing
    . consistent hashing (without virutal node - complex)
  - Re-shuffle data, if new server is -
  		- added
		- deleted // TODO
*/

type User struct {
	id    int
	name  string
	phone string
	hash  uint64
}

var servers = make(map[int][]User)
var serverRange = map[int][]int{
	0: {0, 5, 10, 15, 20},
	1: {25, 30, 35, 40},
	2: {45, 50, 55, 60},
	3: {65, 70, 75, 80, 85, 90, 95, 100},
}
var visitedIds = make(map[string]bool)

const ringSize = 100

var method string

func listServerData() {
	for server, users := range servers {
		fmt.Printf("\nServer %d: \n", server)
		for _, user := range users {
			fmt.Printf("- uid: %d, name: %s, phone: %s, hash: %d\n", user.id, user.name, user.phone, user.hash)
		}
		if len(users) == 0 {
			fmt.Println(" (no users yet)")
		}
	}
}

// Get server count randomly
// Output will be unpredictable/indeterministic
func getServerByRandomNumber(serverCount int) (randomServer int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(serverCount)
}

// Get server by only hashing
// Output will be predictable/deterministic
func getServerByHashing(id int, serverCount int) (randomServer int) {
	idStr := strconv.Itoa(id)
	hash := md5.Sum([]byte(idStr))
	numHash := binary.BigEndian.Uint64(hash[:8])
	return int(numHash % uint64(serverCount))
}

// Get server by consistent hashing
// Output will be predictable/d eterministic
func getServerByConsistentHashing(id int) (randomServer int, numHash uint64) {
	idStr := strconv.Itoa(id)
	hash := md5.Sum([]byte(idStr))
	numHash = binary.BigEndian.Uint64(hash[:8])
	ringNumber := int(numHash % 100)

	keys := make([]int, 0, len(serverRange))
	for k := range serverRange {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		val := serverRange[key]
		for _, v := range val {
			if v > ringNumber {
				return key, uint64(ringNumber)
			}
		}
	}
	return 0, 0
}

// Re-shuffling the server
func reshufflingServers(method string, addServersCount int) {
	newServerCount := len(servers) + addServersCount
	tempServerMap := make(map[int][]User)

	// create map to later use for visited IDs
	for _, users := range servers {
		for _, user := range users {
			visitedIds[strconv.Itoa(user.id)] = false
		}
	}

	switch method {
	case "hashing":
		start := time.Now()

		for _, users := range servers {
			for _, user := range users {
				if !visitedIds[strconv.Itoa(user.id)] {
					serverId := getServerByHashing(user.id, newServerCount)
					tempServerMap[serverId] = append(tempServerMap[serverId], User{id: user.id, name: user.name, phone: user.phone, hash: 0})
					visitedIds[strconv.Itoa(user.id)] = true
				}
			}
		}

		fmt.Println(time.Since(start).Nanoseconds())
		servers = tempServerMap
		clear(visitedIds)
	case "consistent_hashing":
		serverRange = map[int][]int{
			0: {0, 5, 10, 15, 20},    // 0
			1: {25, 30},              // 1
			2: {35, 40},              // new
			3: {45, 50, 55, 60},      // 2
			4: {65, 70, 75},          // new
			5: {80, 85, 90, 95, 100}, // 3
		}
		start := time.Now()
		for _, users := range servers {
			for _, user := range users {
				if !visitedIds[strconv.Itoa(user.id)] {
					serverId, hash := getServerByConsistentHashing(user.id)
					tempServerMap[serverId] = append(tempServerMap[serverId], User{id: user.id, name: user.name, phone: user.phone, hash: hash})
					visitedIds[strconv.Itoa(user.id)] = true
				}
			}
		}
		fmt.Println(time.Since(start).Nanoseconds())
		servers = tempServerMap
		clear(visitedIds)
	default:
		// TODO
	}
}

// Load users from file
func readUsersFromFile(filename string) ([]User, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []User
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // skip blank lines
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Println("Skipping invalid line:", line)
			continue
		}

		name := parts[0]
		phone := parts[1]

		users = append(users, User{name: name, phone: phone})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func main() {
	var uid = 0
	var choice int8

	// Initialize few servers
	for i := range 4 {
		servers[i] = []User{}
		// fmt.Println(servers[i])
	}

outer: // label the loop, to use with break
	for {
		fmt.Print("Choice of operation:")
		fmt.Println(" 1 = Add User, 2 = Retrieve User, 3 = Add New Server, 9 = Exit App")
		fmt.Scanf("%d", &choice)
		switch choice {
		case 1:
			users, err := readUsersFromFile("users.txt")
			if err != nil {
				fmt.Println("Error: ", err)
				break outer
			}

			var hash uint64 = 0
			for _, u := range users {
				uid += 1

				// serverId := getServerByRandomNumber(len(servers))
				// method = "random"

				serverId := getServerByHashing(uid, len(servers))
				method = "hashing"

				// serverId, hash := getServerByConsistentHashing(uid)
				// method = "consistent_hashing"

				servers[serverId] = append(servers[serverId], User{id: uid, name: u.name, phone: u.phone, hash: hash})
			}
		case 2:
			listServerData()
		case 3:
			reshufflingServers(method, 2)
		default:
			break outer
		}
	}
}
