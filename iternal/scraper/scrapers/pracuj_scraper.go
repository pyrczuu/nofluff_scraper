package scrapers

import (
	"context"
	"github.com/gocolly/colly"
	"github.com/pfczx/jobscraper/iternal/scraper"
	"log"
	"time"
)

const (
	titleSelector         = "h1[data-scroll-id='job-title']"
	companySelector       = "h2[data-scroll-id='employer-name']"
	locationSelector      = "div[data-test='offer-badge-title']"
	descriptionSelector   = `ul[data-test="text-about-project"]`                                                         //concat in code
	skillsSelector        = `span[data-test="item-technologies-expected"], span[data-test="item-technologies-optional"]` //concat in code
	salarySectionSelector = `div[data-test="section-salaryPerContractType"]`
	salaryAmountSelector  = `div[data-test="text-earningAmount"]`
	contractTypeSelector  = `span[data-test="text-contractTypeName"]`
)

type PracujScraper struct {
	timeoutBetweenScraps time.Duration
	collector            *colly.Collector
	urls                 []string
}

// controls
func NewPracujScraper(urls []string) *PracujScraper {
	c := colly.NewCollector(
		colly.AllowedDomains("www.pracuj.pl", "pracuj.pl"),
		//colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*pracuj.pl*",
		Parallelism: 2,
		RandomDelay: 2 * time.Second,
	})

	return &PracujScraper{
		timeoutBetweenScraps: 10 * time.Second,
		collector:            c,
	}
}

func (*PracujScraper) Source() string {
	return "pracuj.pl"
}

func (p *PracujScraper) Scrape(ctx context.Context, q chan<- scraper.JobOffer) error {
	log.Println("=== PracujScraper started ===")
	p.collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
	// Dodaj debugowe logi requestów i błędów
	p.collector.OnRequest(func(r *colly.Request) {
		log.Printf("[DEBUG] Requesting URL: %s", r.URL.String())
	})
	p.collector.OnError(func(r *colly.Response, err error) {
		log.Printf("[DEBUG] Visit failed: %v (URL: %s)", err, r.Request.URL.String())
	})
	p.collector.OnResponse(func(r *colly.Response) {
		log.Printf("Response: %d bytes from %s", len(r.Body), r.Request.URL.String())
		log.Println(string(r.Body[:min(500, len(r.Body))])) // debug snippet
	})

	// Rejestracja callbacku HTML
	p.collector.OnHTML("html", func(e *colly.HTMLElement) {
		log.Printf("[DEBUG] HTML parsed for page: %s", e.Request.URL.String())

		select {
		case <-ctx.Done():
			log.Println("[DEBUG] Context cancelled, returning")
			return
		default:
		}

		var job scraper.JobOffer
		job.URL = e.Request.URL.String()
		job.Source = p.Source()
		job.Title = e.ChildText(titleSelector)
		job.Company = e.ChildText(companySelector)
		job.Location = e.ChildText(locationSelector)

		// description
		e.ForEach(descriptionSelector+" li", func(_ int, el *colly.HTMLElement) {
			if text := el.Text; text != "" {
				job.Description += text + "\n"
			}
		})

		// skills
		var skills []string
		e.ForEach(skillsSelector, func(_ int, el *colly.HTMLElement) {
			if text := el.Text; text != "" {
				skills = append(skills, text)
			}
		})
		job.Skills = skills

		// salary
		e.ForEach(salarySectionSelector, func(_ int, el *colly.HTMLElement) {
			amount := el.ChildText(salaryAmountSelector)
			ctype := el.ChildText(contractTypeSelector)
			switch ctype {
			case "umowa o pracę":
				job.SalaryEmployment = amount
			case "umowa zlecenie":
				job.SalaryContract = amount
			case "kontrakt B2B":
				job.SalaryB2B = amount
			}
		})

		select {
		case <-ctx.Done():
			log.Println("[DEBUG] Context cancelled before sending job")
			return
		case q <- job:
			log.Printf("[DEBUG] Job sent to channel: %s at %s", job.Title, job.Company)
		}
	})

	// Pętla po URL-ach
	for _, url := range p.urls {
		log.Printf("[DEBUG] Visiting URL: %s", url)
		time.Sleep(p.timeoutBetweenScraps)

		if err := p.collector.Visit(url); err != nil {
			log.Printf("[DEBUG] Visit error: %v", err)
			return err
		}
	}

	// Czekamy na zakończenie requestów
	p.collector.Wait()

	log.Println("=== PracujScraper finished ===")
	return nil
}
