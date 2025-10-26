package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// BoltProxyTester handles testing the bolt-proxy connection
type BoltProxyTester struct {
	driver neo4j.Driver
	ctx    context.Context
}

// NewBoltProxyTester creates a new tester instance
func NewBoltProxyTester(uri, username, password string) (*BoltProxyTester, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	return &BoltProxyTester{
		driver: driver,
		ctx:    context.Background(),
	}, nil
}

// Close closes the driver connection
func (t *BoltProxyTester) Close() error {
	return t.driver.Close()
}

// TestConnection tests basic connectivity
func (t *BoltProxyTester) TestConnection() error {
	fmt.Println("Testing connection...")
	session := t.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run("RETURN 1 AS num", nil)
	if err != nil {
		return fmt.Errorf("✗ Connection failed: %w", err)
	}

	if result.Next() {
		num := result.Record().Values[0]
		if num.(int64) == 1 {
			fmt.Println("✓ Connection successful")
			return nil
		}
	}

	return fmt.Errorf("✗ Unexpected result")
}

// TestCreateNode tests creating a node
func (t *BoltProxyTester) TestCreateNode() error {
	fmt.Println("\nTesting node creation...")
	session := t.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(
		"CREATE (n:TestNode {name: $name, timestamp: timestamp()}) RETURN n",
		map[string]interface{}{
			"name": "BoltProxyTest",
		},
	)
	if err != nil {
		return fmt.Errorf("✗ Node creation failed: %w", err)
	}

	if result.Next() {
		node := result.Record().Values[0]
		fmt.Printf("✓ Created node: %v\n", node)
		return nil
	}

	return fmt.Errorf("✗ No node returned")
}

// TestQueryNode tests querying nodes
func (t *BoltProxyTester) TestQueryNode() error {
	fmt.Println("\nTesting node query...")
	session := t.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(
		"MATCH (n:TestNode {name: $name}) RETURN n LIMIT 5",
		map[string]interface{}{
			"name": "BoltProxyTest",
		},
	)
	if err != nil {
		return fmt.Errorf("✗ Query failed: %w", err)
	}

	count := 0
	for result.Next() {
		count++
		node := result.Record().Values[0]
		fmt.Printf("  Found node: %v\n", node)
	}

	if count > 0 {
		fmt.Printf("✓ Found %d node(s)\n", count)
		return nil
	}

	return fmt.Errorf("✗ No nodes found")
}

// TestCleanup cleans up test data
func (t *BoltProxyTester) TestCleanup() error {
	fmt.Println("\nCleaning up test data...")
	session := t.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.Run(
		"MATCH (n:TestNode {name: $name}) DELETE n",
		map[string]interface{}{
			"name": "BoltProxyTest",
		},
	)
	if err != nil {
		return fmt.Errorf("✗ Cleanup failed: %w", err)
	}

	fmt.Println("✓ Cleanup successful")
	return nil
}

// RunAllTests executes all test cases
func (t *BoltProxyTester) RunAllTests() bool {
	fmt.Println(strings("=", 60))
	fmt.Println("Bolt Proxy Test Suite")
	fmt.Println(strings("=", 60))

	tests := []struct {
		name string
		fn   func() error
	}{
		{"Connection", t.TestConnection},
		{"Create Node", t.TestCreateNode},
		{"Query Node", t.TestQueryNode},
		{"Cleanup", t.TestCleanup},
	}

	passed := 0
	failed := 0

	for _, test := range tests {
		if err := test.fn(); err != nil {
			log.Printf("Test '%s' failed: %v", test.name, err)
			failed++
		} else {
			passed++
		}
		time.Sleep(100 * time.Millisecond) // Small delay between tests
	}

	fmt.Println()
	fmt.Println(strings("=", 60))
	fmt.Printf("Results: %d passed, %d failed\n", passed, failed)
	fmt.Println(strings("=", 60))

	return failed == 0
}

func strings(char string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += char
	}
	return result
}

func main() {
	host := flag.String("host", "localhost", "Bolt-proxy host")
	port := flag.String("port", "8080", "Bolt-proxy port")
	username := flag.String("user", "", "Username for authentication")
	password := flag.String("password", "", "Password for authentication")

	flag.Parse()

	uri := fmt.Sprintf("bolt://%s:%s", *host, *port)
	fmt.Printf("Connecting to: %s\n", uri)
	if *username != "" {
		fmt.Printf("User: %s\n\n", *username)
	} else {
		fmt.Println("User: (none)\n")
	}

	tester, err := NewBoltProxyTester(uri, *username, *password)
	if err != nil {
		log.Fatalf("Failed to create tester: %v", err)
	}
	defer tester.Close()

	success := tester.RunAllTests()
	if !success {
		os.Exit(1)
	}
}
