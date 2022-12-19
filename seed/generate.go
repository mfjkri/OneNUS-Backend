package seed

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/mfjkri/One-NUS-Backend/controllers"
	"github.com/mfjkri/One-NUS-Backend/database"
	"github.com/mfjkri/One-NUS-Backend/models"
	"golang.org/x/crypto/bcrypt"
)


func GenerateUser() models.User {
	password, _ :=bcrypt.GenerateFromPassword([]byte("password"), 10)
	return models.User{Username: strings.ToLower(faker.LastName(options.WithRandomStringLength(10))), Password: password}
}

func GenerateUsers(number int) []models.User {
	users := make([]models.User, number)

	for index := 1; index <= number; index++ {
		user := GenerateUser()
		users = append(users, user)
		database.DB.Create(&user)
	}

	return users
}

func ChooseRandomTag() string {
	return models.ValidTags[rand.Intn(len(models.ValidTags))]
}

func GeneratePost(user models.User, creationTime time.Time) Post {
	return Post{
		CreatedAt: creationTime,
		UpdatedAt: creationTime,

		Title: faker.Sentence(options.WithRandomStringLength(uint(controllers.MAX_POST_TITLE_CHAR))),
		Tag: ChooseRandomTag(),
		Text: faker.Paragraph(options.WithRandomStringLength(uint(controllers.MAX_POST_TEXT_CHAR))),

		Author: user.Username,
		User: user,

		CommentsCount: 0,
		CommentedAt: time.Unix(0, 0),
		StarsCount: uint(rand.Intn((100))),
	}
}

func GeneratePosts(number int, user models.User, startTime time.Time) {
	initialTime := startTime.Add(time.Duration(rand.Intn(100)))

	for index := 1; index <= number; index++ {
		post := GeneratePost(user, initialTime)
		database.DB.Create(&post)
		initialTime.Add(
			(time.Millisecond * time.Duration(rand.Intn(100))) + 
			(time.Second * time.Duration(rand.Intn(60))) + 
			(time.Minute * time.Duration(rand.Intn(60))))
	}
}

func GeneratePostsForEachUser() {
	var users []models.User
	database.DB.Find(&users)

	initialTime := time.Now().Add(-time.Hour * 24 * 20)
	for _, user := range(users) {
		GeneratePosts(
			rand.Intn(12),
			user,
			initialTime,
		)
	}
}

func GenerateComment(user models.User, post models.Post, creationTime time.Time) Comment {
	return Comment{
		CreatedAt: creationTime,
		UpdatedAt: creationTime,

		Text: faker.Paragraph(options.WithRandomStringLength(uint(controllers.MAX_COMMENT_TEXT_CHAR))),

		Author: user.Username,
		User: user,

		Post: post,
	}
}

func GenerateComments(number int, user models.User, post models.Post) {
	initialTime := post.CreatedAt.Add(time.Duration(rand.Intn(100)))

	for index := 1; index <= number; index++ {
		comment := GenerateComment(user, post, initialTime)
		database.DB.Create(&comment)
		initialTime = initialTime.Add(
			(time.Millisecond * time.Duration(rand.Intn(100))) + 
			(time.Second * time.Duration(rand.Intn(60))) + 
			(time.Minute * time.Duration(rand.Intn(60))))
	}
}

func GenerateCommentsForEachPost() {
	var users []models.User
	database.DB.Find(&users)

	var posts []models.Post
	database.DB.Find(&posts)

	for _, user := range(users) {
		shouldComment := rand.Float64() > 0.4

		if shouldComment {
			for _, post := range(posts) {
				GenerateComments(
					rand.Intn(3),
					user,
					post,
				)
			}
		}
	}
}

func GenerateData() {
	fmt.Println("Seeding...")

	GenerateUsers(20)
	GeneratePostsForEachUser()
	GenerateCommentsForEachPost()

	fmt.Println("Seeding complete!")
}