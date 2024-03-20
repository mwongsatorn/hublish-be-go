package database

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"

	"hublish-be-go/internal/models"
	"hublish-be-go/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedDatabase() {
	var (
		users       []models.User
		usersMap    []map[string]string
		articles    []models.Article
		articlesMap []map[string]string
	)

	db := DB

	var countData int64
	if countResult := db.Model(&models.User{}).Count(&countData); countResult.Error != nil {
		panic(err)
	}
	if countData != 0 {
		fmt.Println("\n\nDatabase already has the data")
		return
	}

	usersFile, err := os.Open("internal/database/users.json")
	if err != nil {
		panic(err)
	}
	articlesFile, err := os.Open("internal/database/articles.json")
	if err != nil {
		panic(err)
	}

	defer usersFile.Close()
	defer articlesFile.Close()

	usersBytes, err := io.ReadAll(usersFile)
	if err != nil {
		panic(err)
	}
	articlesBytes, err := io.ReadAll(articlesFile)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(usersBytes, &usersMap)
	json.Unmarshal(articlesBytes, &articlesMap)

	for i := 0; i < len(usersMap); i++ {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usersMap[i]["password"]), 10)
		if err != nil {
			panic(err)
		}
		newUser := models.User{
			Username: usersMap[i]["username"],
			Email:    usersMap[i]["email"],
			Name:     usersMap[i]["name"],
			Password: string(hashedPassword),
			Image:    usersMap[i]["image"],
		}
		db.Create(&newUser)
		users = append(users, newUser)
	}

	tagList := [][]string{
		{"Lorem", "ipsum", "dolor"},
		{"consectetur", "adipiscing"},
		{"Mauris", "rutrum", "ligula", "et"},
		{"Vivamus"},
		{"Nulla", "porttitor"},
		{"Etiam", "pretium"},
		{"molestie"},
		{"Suspendisse", "eugiat ornare", "cursus"},
	}

	for i := 0; i < len(articlesMap); i++ {
		slug := utils.GenerateSlug(articlesMap[i]["title"])

		newArticle := models.Article{
			Title:    articlesMap[i]["title"],
			Slug:     slug,
			Content:  articlesMap[i]["content"],
			Tags:     tagList[rand.Intn(8)],
			AuthorID: users[rand.Intn(10)].ID,
		}
		db.Create(&newArticle)
		articles = append(articles, newArticle)
	}

	for i := 0; i < len(users); i++ {
		for j := 0; j < len(users); j++ {

			if i == j {
				continue
			}

			follow := rand.Float32() < 0.5
			if !follow {
				continue
			}

			db.Transaction(func(tx *gorm.DB) error {

				tx.Create(&models.Follow{
					FollowingID: users[i].ID,
					FollowerID:  users[j].ID,
				})

				tx.Model(&models.User{}).
					Where("id = ?", users[j].ID).
					Update("following_count", gorm.Expr("following_count + 1"))

				tx.Model(&models.User{}).
					Where("id = ?", users[i].ID).
					Update("follower_count", gorm.Expr("follower_count + 1"))

				return nil
			})
		}
	}

	for i := 0; i < len(users); i++ {
		for j := 0; j < len(articles); j++ {
			favourite := rand.Float32() < 0.5
			if !favourite {
				continue
			}

			db.Transaction(func(tx *gorm.DB) error {

				tx.Create(&models.Favourite{
					UserID:    users[i].ID,
					ArticleID: articles[j].ID,
				})

				tx.Model(&models.Article{}).
					Where("id = ?", articles[j].ID).
					Update("favourite_count", gorm.Expr("favourite_count + 1"))

				return nil
			})
		}
	}

}
