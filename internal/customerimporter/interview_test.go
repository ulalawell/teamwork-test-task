package customerimporter

import (
	"os"
	"testing"
)

func TestNewCustomerImporter(t *testing.T) {
	_, err := NewCustomerImporter("customers_test.csv", "email")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test invalid file path
	_, err = NewCustomerImporter("invalid.txt", "email")
	if err == nil {
		t.Error("Expected error for invalid file path, got nil")
	}

	// Test non-existent file
	_, err = NewCustomerImporter("nonexistent.csv", "email")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	// Test empty email field name
	_, err = NewCustomerImporter("example.csv", "")
	if err == nil {
		t.Error("Expected error for empty email field name, got nil")
	}
}

func TestGetDomainInformation_FileOpenError(t *testing.T) {
	// Test when there is an error opening the file
	importer := &customerImporter{csvFilePath: "nonexistent.csv", emailFieldName: "email"}
	domainInfo, err := importer.GetDomainCounts()
	if err == nil {
		t.Error("Expected error for file open failure, got nil")
	}

	// Since the file cannot be opened, there should be no domain information
	var expected []emailDomain
	if !isEqual(domainInfo, expected) {
		t.Errorf("Expected %v, got %v", expected, domainInfo)
	}
}

func TestGetDomainInformation_EmptyFile(t *testing.T) {
	// Test when there is an error reading the headers row from the CSV file
	file, err := os.CreateTemp("", "read_headers_error*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	file.Close()

	importer := &customerImporter{csvFilePath: file.Name(), emailFieldName: "email"}
	domainInfo, err := importer.GetDomainCounts()
	if err == nil {
		t.Error("Expected error for read headers failure, got nil")
	}

	// Since there is an error reading headers, there should be no domain information
	var expected []emailDomain
	if !isEqual(domainInfo, expected) {
		t.Errorf("Expected %v, got %v", expected, domainInfo)
	}
}

func TestGetDomainInformation_EmailFieldNotFound(t *testing.T) {
	// Test when the specified email field is not found in the headers
	file, err := os.CreateTemp("", "email_field_not_found*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	csvHeader := "first_name,last_name,invalid_field,gender,ip_address"
	_, err = file.WriteString(csvHeader)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	importer := &customerImporter{csvFilePath: file.Name(), emailFieldName: "email"}
	domainInfo, err := importer.GetDomainCounts()
	if err == nil {
		t.Error("Expected error for email field not found, got nil")
	}

	// Since the email field is not found, there should be no domain information
	var expected []emailDomain
	if !isEqual(domainInfo, expected) {
		t.Errorf("Expected %v, got %v", expected, domainInfo)
	}
}

func TestGetDomainInformation_OnlyHeader(t *testing.T) {
	// Create a temporary empty CSV file for testing
	file, err := os.CreateTemp("", "only_header*.csv")
	if err != nil {
		t.Fatal(err)
	}

	csvHeader := "first_name,last_name,email,gender,ip_address"
	_, err = file.WriteString(csvHeader)

	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	file.Close()

	importer, err := NewCustomerImporter(file.Name(), "email")
	if err != nil {
		t.Fatal(err)
	}

	domainInfo, err := importer.GetDomainCounts()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Since the file is empty, there should be no domain information
	var expected []emailDomain
	if !isEqual(domainInfo, expected) {
		t.Errorf("Expected %v, got %v", expected, domainInfo)
	}
}

func TestGetDomainInformation_InvalidEmails(t *testing.T) {
	// Create a temporary CSV file for testing
	file, err := os.CreateTemp("", "example_invalid_emails*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Write sample CSV data with invalid emails to the file
	csvHeader := "first_name,last_name,email,gender,ip_address"
	firstRow := "Mildred,Hernandez,invalid_email,Female,38.194.51.128"
	secondRow := "Bonnie,Ortiz,bortiz1@example.com,Female,197.54.209.129"

	_, err = file.WriteString(csvHeader + "\n" + firstRow + "\n" + secondRow)

	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	importer, err := NewCustomerImporter(file.Name(), "email")
	if err != nil {
		t.Fatal(err)
	}

	domainInfo, err := importer.GetDomainCounts()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Since the second email is invalid, it should not contribute to the domain count
	expected := []emailDomain{
		{"example.com", 1},
	}
	if !isEqual(domainInfo, expected) {
		t.Errorf("Expected %v, got %v", expected, domainInfo)
	}
}

func TestGetDomainInformation(t *testing.T) {
	// Create a temporary CSV file for testing
	file, err := os.CreateTemp("", "example*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Write sample CSV data with invalid emails to the file
	csvHeader := "first_name,last_name,email,gender,ip_address"
	firstRow := "Mildred,Hernandez,bortiz2@example.com,Female,38.194.51.128"
	secondRow := "Bonnie,Ortiz,bortiz1@example.com,Female,197.54.209.129"

	_, err = file.WriteString(csvHeader + "\n" + firstRow + "\n" + secondRow)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	importer, err := NewCustomerImporter(file.Name(), "email")
	if err != nil {
		t.Fatal(err)
	}

	domainInfo, err := importer.GetDomainCounts()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Add assertions for the expected results based on the sample CSV data
	expected := []emailDomain{
		{"example.com", 2},
	}
	if !isEqual(domainInfo, expected) {
		t.Errorf("Expected %v, got %v", expected, domainInfo)
	}
}

func isEqual(a, b []emailDomain) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
