package seed

import (
	"fmt"

	"github.com/mfjkri/OneNUS-Backend/database"
)

func DeleteUsers() {
	fmt.Println("Deleting users")
	database.DB.Migrator().DropTable("users")
}

func DeletePosts() {
	fmt.Println("Deleting posts")
	database.DB.Migrator().DropTable("posts")
}

func DeleteComments() {
	fmt.Println("Deleting comments")
	database.DB.Migrator().DropTable("comments")
}

func DeleteAll() {
	fmt.Println("RESETTING DATABASE")
	DeleteUsers()
	DeletePosts()
	DeleteComments()
	fmt.Println("DATABASE RESETTED")
}
