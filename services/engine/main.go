package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/samyosm/astatine/services/engine/tokenizer"
)

type term = string

func count_freq(terms []term) map[term]int {
	m := make(map[term]int)
	for _, term := range terms {
		m[term] += 1
	}

	return m
}

var wg sync.WaitGroup

func indexDocument(document string, ctx context.Context, driver neo4j.DriverWithContext) {
	defer wg.Done()
	document = strings.ToLower(document)
	for term, count := range count_freq(tokenizer.Wordpunkt(document)) {
		_, err := neo4j.ExecuteQuery(ctx, driver,
			`merge (document:Document {content: $doc_content}) merge (term:Term { content: $term_content}) create (term) -[:APPEAR_IN { count: $count}]->(document)`,
			map[string]any{
				"doc_content":  document,
				"term_content": term,
				"count":        count,
			}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	ctx := context.Background()
	dbUri := "neo4j://localhost:7687"
	dbUser := "neo4j"
	dbPassword := "testing1234"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	/* Indexing */
	file, _ := os.Open("../../quotes.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		wg.Add(1)
		go indexDocument(scanner.Text(), ctx, driver)
	}

	wg.Wait()
}
