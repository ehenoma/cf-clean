package main

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/manifoldco/promptui"
	"log"
	"os"
)

func main() {
	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	zoneName := promptInput("What is the Name of the Zone you want to clear?")
	id, err := api.ZoneIDByName(zoneName)
	if err != nil {
		log.Fatal(err)
	}
	recordType := promptInput("What type of record do you want to clear?")
	log.Printf("Listing %s records. This may take a few seconds...", recordType)
	records, err := api.DNSRecords(ctx, id, cloudflare.DNSRecord{Type: recordType})
	if err != nil {
		log.Fatal(err)
	}
	confirm := fmt.Sprintf("Found %d records. Continue with delete?", len(records))
	if !promptConfirm(confirm) {
		return
	}
	for index, record := range records {
		if err := api.DeleteDNSRecord(ctx, id, record.ID); err != nil {
			log.Fatalf("Failed to delete record %s: %v", record.Name, err)
		}
		log.Printf("Deleted %s (%d/%d)", record.Name, index, len(records))
	}
}

func promptInput(text string) string {
	p := promptui.Prompt{Label: text}
	value, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func promptConfirm(question string) bool {
	prompt := promptui.Select{
		Label: question,
		Items: []string{"Yes", "No"},
	}
	_, answer, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}
	return answer == "Yes"
}
