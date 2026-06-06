package gnucash

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type GnuCashFile struct {
	XMLName xml.Name `xml:"gnc-v2"`
	Book    Book     `xml:"book"`
}

type Book struct {
	Accounts     []Account     `xml:"account"`
	Transactions []Transaction `xml:"transaction"`
}

type Account struct {
	Name      string    `xml:"name"`
	ID        string    `xml:"id"`
	Type      string    `xml:"type"`
	Commodity Commodity `xml:"commodity"`
	Parent    string    `xml:"parent"`
}

type Commodity struct {
	ID string `xml:"id"`
}

type Transaction struct {
	ID          string    `xml:"id"`
	DatePosted  string    `xml:"date-posted>date"`
	Description string    `xml:"description"`
	Currency    Commodity `xml:"currency"`
	Splits      Splits    `xml:"splits"`
}

type Splits struct {
	Split []Split `xml:"split"`
}

type Split struct {
	Account  string `xml:"account"`
	Value    string `xml:"value"`
	Quantity string `xml:"quantity"`
}

type Parser struct {
	file       *GnuCashFile
	accountMap map[string]*Account
	Verbose    bool
}

func New() *Parser {
	return &Parser{
		accountMap: make(map[string]*Account),
	}
}

func (p *Parser) Parse(filename string) error {
	var data []byte
	var err error

	if filename == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(filename)
	}

	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	p.file = &GnuCashFile{}
	if err := xml.Unmarshal(data, p.file); err != nil {
		return fmt.Errorf("failed to parse XML: %w", err)
	}

	p.buildAccountMap()
	return nil
}

func (p *Parser) buildAccountMap() {
	for i := range p.file.Book.Accounts {
		account := &p.file.Book.Accounts[i]
		p.accountMap[account.ID] = account
	}
}

func (p *Parser) GetAccountName(accountID string) string {
	return p.getAccountNameRecursive(accountID, "")
}

func (p *Parser) getAccountNameRecursive(accountID, base string) string {
	account, exists := p.accountMap[accountID]
	if !exists {
		return base
	}

	if account.Type == "ROOT" {
		return base
	}

	current := account.Name
	if base != "" {
		current = current + ":" + base
	}

	if account.Parent != "" {
		return p.getAccountNameRecursive(account.Parent, current)
	}

	return current
}

func (p *Parser) GetTransactions() []Transaction {
	return p.file.Book.Transactions
}

func (p *Parser) GetAccount(accountID string) *Account {
	return p.accountMap[accountID]
}

func ParseFraction(fraction string) float64 {
	parts := strings.Split(fraction, "/")

	if len(parts) == 1 {
		val, _ := strconv.ParseFloat(fraction, 64)
		return val
	}

	if len(parts) == 2 {
		numerator, _ := strconv.ParseFloat(parts[0], 64)
		denominator, _ := strconv.ParseFloat(parts[1], 64)

		if denominator == 0 {
			return 0
		}

		return numerator / denominator
	}

	panic(fmt.Sprintf("invalid fraction format: %s (has %d parts)", fraction, len(parts)))
}
