package customerimporter

// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// customerImporter represents a CSV file importer for customer data.
type customerImporter struct {
	csvFilePath    string
	emailFieldName string
}

// emailDomain represents the structure holding the email domain and the number
// of customers with email addresses for that domain.
type emailDomain struct {
	Domain        string
	CustomerCount int
}

// NewCustomerImporter creates a new customerImporter instance.
func NewCustomerImporter(csvFilePath string, emailFieldName string) (*customerImporter, error) {
	if !strings.HasSuffix(csvFilePath, ".csv") {
		return nil, fmt.Errorf("invalid file path: %s, expecting a '.csv' file", csvFilePath)
	}

	if _, err := os.Stat(csvFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("the file does not exist at the specified path: %s", csvFilePath)
	}

	if emailFieldName == "" {
		return nil, fmt.Errorf("the email field name is empty")
	}

	return &customerImporter{
		csvFilePath:    csvFilePath,
		emailFieldName: emailFieldName,
	}, nil
}

// GetDomainCounts reads the CSV file, extracts email domains, and returns a
// sorted list of email domains with customer counts.
func (importer *customerImporter) GetDomainCounts() ([]emailDomain, error) {
	file, err := os.Open(importer.csvFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open the file %s: %w", importer.csvFilePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	emailDomainCounts := make(map[string]int)

	fileHeaders, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read the headers row from the CSV file: %w", err)
	}

	emailIndex := indexOf(fileHeaders, importer.emailFieldName)
	if emailIndex == -1 {
		return nil, fmt.Errorf("failed to find the field '%s' in the headers", importer.emailFieldName)
	}

	//mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	emailAddresses := make(chan string)

	go func() {
		wg.Add(1)

		for {
			record, err := reader.Read()

			// Break the loop if we reach the end of the file
			if err == io.EOF {
				wg.Done()
				close(emailAddresses)

				break
			}

			// Handle other errors
			if err != nil {
				log.Println(err)
				continue
			}

			email := record[emailIndex]
			if !isValidEmail(email) {
				log.Printf("invalid email address found in row %v, skipping", record)
				continue
			}

			emailAddresses <- email
		}

		wg.Done()
	}()

	//for i := 0; i < 4; i++ {
	go func() {
		wg.Add(1)

		for address := range emailAddresses {
			domain, err := extractEmailDomain(address)
			if err != nil {
				log.Println(err)
				continue
			}

			//mutex.Lock()
			emailDomainCounts[domain]++
			//mutex.Unlock()
		}
		wg.Done()
	}()
	//}

	wg.Wait()

	sortedDomains := sortEmailDomainsByCount(emailDomainCounts)

	return sortedDomains, nil
}

// isValidEmail checks if the provided email address is valid.
func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}

// extractEmailDomain extracts the domain from the given email address.
func extractEmailDomain(email string) (string, error) {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1], nil
	}

	return "", fmt.Errorf("unable to extract domain from email: %s", email)
}

// sortEmailDomainsByCount sorts the email domains by customer count in descending order.
func sortEmailDomainsByCount(m map[string]int) []emailDomain {
	var sorted []emailDomain
	for k, v := range m {
		sorted = append(sorted, emailDomain{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CustomerCount > sorted[j].CustomerCount
	})
	return sorted
}

// indexOf finds the index of the given value in the slice.
func indexOf(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}
