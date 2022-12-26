package seed

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/mfjkri/OneNUS-Backend/controllers"
	"github.com/mfjkri/OneNUS-Backend/database"
	"github.com/mfjkri/OneNUS-Backend/models"
	"golang.org/x/crypto/bcrypt"
)

var NEW_USERS_COUNT int = 10
var MAX_POST_PER_USER = 5
var MAX_COMMENT_PER_USER_PER_POST = 2
var POST_CREATION_TIME_OFFSET_HOURS = time.Duration(24 * 5)

func LoadGenerateConfig() {
	newUsersCount, err := strconv.Atoi(os.Getenv("GENERATE_NEW_USERS_COUNT"))
	if err == nil {
		NEW_USERS_COUNT = newUsersCount
	}

	maxPostPerUser, err := strconv.Atoi(os.Getenv("GENERATE_MAX_POST_PER_USER"))
	if err == nil {
		MAX_POST_PER_USER = maxPostPerUser
	}

	maxCommentPerUserPerPost, err := strconv.Atoi(os.Getenv("GENERATE_MAX_COMMENT_PER_USER_PER_POST"))
	if err == nil {
		MAX_COMMENT_PER_USER_PER_POST = maxCommentPerUserPerPost
	}

	postCreationTimeOffset, err := strconv.Atoi(os.Getenv("GENERATE_POST_CREATION_TIME_OFFSET_HOURS"))
	if err == nil {
		POST_CREATION_TIME_OFFSET_HOURS = time.Duration(postCreationTimeOffset)
	}
}

func FastForwardTime(currentTime time.Time) time.Time {
	return currentTime.Add(
		(time.Millisecond * time.Duration(rand.Intn(100))) +
			(time.Second * time.Duration(rand.Intn(60))) +
			(time.Minute * time.Duration(rand.Intn(60))))
}

/* -------------------------------------------------------------------------- */
/*                               Generate Users                               */
/* -------------------------------------------------------------------------- */
func GenerateUser() models.User {
	password, _ := bcrypt.GenerateFromPassword([]byte("2%2$iK66&*R#S38MY9tJ*5UZ6f!7f"), 10)
	return models.User{
		Username: strings.ToLower(faker.LastName(options.WithRandomStringLength(10), options.WithGenerateUniqueValues(true))),
		Password: password,
		Role:     "member",
	}
}

func GenerateUsers(number int) {
	fmt.Println("Generating users", number)

	for index := 1; index <= number; index++ {
		user := GenerateUser()
		database.DB.Create(&user)
	}

	fmt.Println("Users generated.")
}

/* -------------------------------------------------------------------------- */
/*                               Generate Posts                               */
/* -------------------------------------------------------------------------- */
func ChooseRandomTag() string {
	return models.ValidTags[rand.Intn(len(models.ValidTags))]
}

func GeneratePost(user models.User) models.Post {
	return models.Post{
		Title: faker.Sentence(options.WithRandomStringLength(uint(controllers.MAX_POST_TITLE_CHAR))),
		Tag:   ChooseRandomTag(),
		Text:  faker.Paragraph(options.WithRandomStringLength(uint(controllers.MAX_POST_TEXT_CHAR))),

		Author: user.Username,
		User:   user,

		CommentsCount: 0,
		CommentedAt:   time.Unix(0, 0),
		StarsCount:    uint(rand.Intn((100))),
	}
}

func GeneratePosts(number int, user models.User, creationTime time.Time) {
	for index := 1; index <= number; index++ {
		post := GeneratePost(user)
		database.DB.Create(&post)
		database.DB.Model(&post).Update("created_at", creationTime)
		database.DB.Model(&post).Update("updated_at", creationTime)

		creationTime = FastForwardTime(creationTime)
	}
}

func GeneratePostsForEachUser(maxPostPerUser int, initialTime time.Time) {
	fmt.Println("Generating posts for each generated user with max post of:", maxPostPerUser)

	var users []models.User
	database.DB.Find(&users)

	for _, user := range users {
		GeneratePosts(
			rand.Intn(maxPostPerUser),
			user,
			initialTime,
		)
		initialTime = FastForwardTime(initialTime)
	}
	fmt.Println("Posts generated.")
}

/* -------------------------------------------------------------------------- */
/*                              Generate Comments                             */
/* -------------------------------------------------------------------------- */
func GenerateComment(user models.User, post models.Post) models.Comment {
	return models.Comment{
		Text: faker.Paragraph(options.WithRandomStringLength(uint(controllers.MAX_COMMENT_TEXT_CHAR))),

		Author: user.Username,
		User:   user,

		Post: post,
	}
}

func GenerateComments(number int, user models.User, post models.Post) {
	creationTime := FastForwardTime(post.CreatedAt)

	for index := 1; index <= number; index++ {
		comment := GenerateComment(user, post)
		database.DB.Create(&comment)
		database.DB.Model(&comment).Update("created_at", creationTime)
		database.DB.Model(&comment).Update("updated_at", creationTime)

		creationTime = FastForwardTime(FastForwardTime(creationTime))
	}
}

func GenerateCommentsForEachPost(maxRandomCommentsPerUserPerPost int) {
	fmt.Println("Generating comments for each generated post with max comments per user per post of:", maxRandomCommentsPerUserPerPost)

	var users []models.User
	database.DB.Find(&users)

	var posts []models.Post
	database.DB.Find(&posts)

	for _, user := range users {
		shouldComment := rand.Float64() > 0.4
		if shouldComment {
			for _, post := range posts {
				GenerateComments(
					rand.Intn(maxRandomCommentsPerUserPerPost),
					user,
					post,
				)
			}
		}
	}

	fmt.Println("Comments generated.")
}

/* -------------------------------------------------------------------------- */

func GenerateData() {
	fmt.Println("Seeding...")

	LoadGenerateConfig()

	GenerateUsers(NEW_USERS_COUNT)
	GeneratePostsForEachUser(MAX_POST_PER_USER, time.Now().Add(time.Hour*POST_CREATION_TIME_OFFSET_HOURS))
	GenerateCommentsForEachPost(MAX_COMMENT_PER_USER_PER_POST)

	fmt.Println("Seeding complete!")
}
