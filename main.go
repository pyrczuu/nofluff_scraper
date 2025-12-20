package main

import (
	"context"
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pfczx/jobscraper/iternal"
	"github.com/pfczx/jobscraper/iternal/scraper"
	"github.com/pfczx/jobscraper/iternal/scraper/scrapers"
)

func main() {
	db, err := sql.Open("sqlite3", "./database/jobs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//read only for js backend
	//_, err = db.Exec("PRAGMA journal_mode=WAL;")

	ctx := context.Background()

	//urls, err := scrapers.ScrollAndRead(ctx)
	//if err != nil {
	//	log.Fatal(err)
	//}
	urls := []string{
		"https://nofluffjobs.com/pl/job/c-software-engineer-motorola-solutions-systems-polska-krakow-13",
		"https://nofluffjobs.com/pl/job/junior-python-engineer-acaisoft-remote-1",
		"https://nofluffjobs.com/pl/job/tester-manualny-astek-polska-warszawa-4"}

	noFluffScraper := scrapers.NewNoFluffScraper(urls)

	scrapersList := []scraper.Scraper{noFluffScraper}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		iternal.StartCollector(ctx, db, scrapersList)
	}()

	wg.Wait()
	log.Println("-------------------")
}
