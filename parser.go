package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type InstagramProfile struct {
	Rank  int
	Name  string
	Nick  string
	Score float64
}

func main() {
	// Открываем файл для записи CSV
	file, err := os.Create("instagram_profiles.csv")
	if err != nil {
		log.Fatal("Error creating CSV file:", err)
	}
	defer file.Close()

	// Создаем записыватель CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Записываем заголовки CSV
	err = writer.Write([]string{"Rank", "Name", "Nick", "Score"})
	if err != nil {
		log.Fatal("Error writing CSV headers:", err)
	}

	// URL для парсинга
	url := "https://hypeauditor.com/top-instagram-all-russia/"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal("Error loading HTML:", err)
	}

	// Парсим данные и записываем в CSV
	doc.Find(".table-profile tbody tr").Each(func(i int, s *goquery.Selection) {
		rank := i + 1
		name := strings.TrimSpace(s.Find("td:eq(1)").Text())
		nick := strings.TrimSpace(s.Find("td:eq(2)").Text())
		scoreStr := strings.TrimSpace(s.Find("td:eq(3)").Text())

		// Преобразуем строку в число
		score, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			log.Printf("Error parsing score for %s: %v", name, err)
			return
		}

		// Создаем объект InstagramProfile
		profile := InstagramProfile{
			Rank:  rank,
			Name:  name,
			Nick:  nick,
			Score: score,
		}

		// Записываем данные в CSV
		err = writer.Write([]string{fmt.Sprintf("%d", profile.Rank), profile.Name, profile.Nick, fmt.Sprintf("%.2f", profile.Score)})
		if err != nil {
			log.Printf("Error writing CSV for %s: %v", profile.Name, err)
		}
	})

	fmt.Println("Scraping completed. Check instagram_profiles.csv")
}
