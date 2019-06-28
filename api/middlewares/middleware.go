package middleware

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/kataras/iris"
)

// Queue is the struct we use for this case
type Queue struct {
	Domain   string
	Weight   int
	Priority int
}

// Que declaration
var Que []string

// FinalQueue delaration
var FinalQueue []*Queue

// RepoResponse declaration
var RepoResponse []*Queue

// Repository should implement common methods
type Repository interface {
	Read() []*Queue
}

func (q *Queue) Read() []*Queue {
	path, _ := filepath.Abs("")
	file, err := os.Open(path + "/api/middlewares/domain.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var final []*Queue
	tmp := &Queue{}
	countLine := 0

	for scanner.Scan() {

		if scanner.Text() == "" { // Empty line is how the domains are separated
			countLine = 0
			continue
		}

		switch countLine {
		case 0: // Line 0 is the domain
			tmp.Domain = scanner.Text()
		case 1: // Line 1 is the weight
			tmp.Weight = parseWeightAndPriority(scanner.Text())
		case 2: // Line 2 is the
			tmp.Priority = parseWeightAndPriority(scanner.Text())
			final = append(final, tmp)
			tmp = &Queue{}
		}
		countLine++
	}
	return final
}

// parseWeightAndPriority gets a string of text and returns the integer of the weight or priority
func parseWeightAndPriority(s string) int {
	r := strings.Split(s, ":")[1]
	res, _ := strconv.Atoi(r)
	return res
}

func getDomainScore(domainElement *Queue) int {
	return (domainElement.Weight + domainElement.Priority) / 2
}

// ProxyMiddleware should queue our incoming requests
func ProxyMiddleware(c iris.Context) {

	Que = nil
	domainFound := false
	if len(RepoResponse) == 0 {
		var repo Repository
		repo = &Queue{}
		RepoResponse = repo.Read() // This contains the information that we need to prioritize requests
	}

	domain := c.GetHeader("domain")
	if len(domain) == 0 {
		c.JSON(iris.Map{"status": 400, "result": "no domain received"})
		return
	}

	for _, current := range RepoResponse {
		if current.Domain == domain {
			newRequest := Queue{current.Domain, current.Weight, current.Priority}
			FinalQueue = append(FinalQueue, &newRequest)
			domainFound = true
			break
		}
	}

	if !domainFound {
		c.JSON(iris.Map{"status": 404, "result": "domain not found"})
		return
	}

	sort.Slice(FinalQueue, func(i, j int) bool {
		return getDomainScore(FinalQueue[i]) > getDomainScore(FinalQueue[j])
	})

	for _, row := range FinalQueue {
		Que = append(Que, row.Domain)
	}

	c.Next()
}
