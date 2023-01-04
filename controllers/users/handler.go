package users

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mfjkri/OneNUS-Backend/config"
	"github.com/mfjkri/OneNUS-Backend/controllers/auth"
	"github.com/mfjkri/OneNUS-Backend/database"
	"github.com/mfjkri/OneNUS-Backend/models"
	"github.com/mfjkri/OneNUS-Backend/utils"
)

/* -------------------------------------------------------------------------- */
/*                GetUserFromID | route: /users/getbyid/:userId               */
/* -------------------------------------------------------------------------- */
type GetUserFromIDRequest struct {
	UserID uint `uri:"userId" binding:"required"`
}

func GetUserFromID(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := auth.VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json GetUserFromIDRequest
	if err := c.ShouldBindUri(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	targetUser, found := auth.FindUserFromID(c, json.UserID)

	if found == false {
		return
	}

	fmt.Printf("%s has requested for: %s\n", user.Username, targetUser.Username)

	// Return fetch User
	c.JSON(http.StatusAccepted, CreateUserResponse(&targetUser))
}

/* -------------------------------------------------------------------------- */
/*                     UpdateBio | route: /users/updatebio                    */
/* -------------------------------------------------------------------------- */
type UpdateBioRequest struct {
	Bio string `json:"bio" binding:"required"`
}

func UpdateBio(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := auth.VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json UpdateBioRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Check that new bio does not contain illegal characters
	if !utils.ContainsValidCharactersOnly(json.Bio) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Bio contains illegal characters."})
		return
	}

	// Update Bio and save
	user.Bio = utils.TrimString(strings.TrimSpace(json.Bio), config.MAX_USER_BIO_LENGTH)
	database.DB.Save(&user)

	fmt.Printf("Updated %s`s bio.\n\tBio: %s\n", user.Username, user.Bio)

	c.JSON(http.StatusAccepted, CreateUserResponse(&user))
}

/* -------------------------------------------------------------------------- */
/*                     DeleteUser | route : /users/delete                     */
/* -------------------------------------------------------------------------- */
func DeleteUser(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := auth.VerifyAuth(c)
	if found == false {
		return
	}

	// Delete all of users Comments (to update existing posts commentsCount correctly)
	// Deletion of user Posts will be handled by CascadeDelete
	var comments []models.Comment
	database.DB.Table("comments").Where("user_id = ?", user.ID).Find(&comments)
	for _, comment := range comments {
		database.DB.Delete(&comment)
	}

	// Delete the RequestUser from database
	database.DB.Delete(&user)

	fmt.Printf("Deleted user: %s.\n", user.Username)

	// Success, user deleted
	c.JSON(http.StatusAccepted, CreateUserResponse(&user))
}
