package config

import "time"

/* -------------------------------------------------------------------------- */
/*                              Config variables                              */
/* -------------------------------------------------------------------------- */
var MAX_PER_PAGE = float64(50)

var MAX_POST_TITLE_CHAR = 100
var MAX_POST_TEXT_CHAR = 5000
var USER_POST_COOLDOWN = time.Second * 45

var MAX_COMMENT_TEXT_CHAR = 1000
var USER_COMMENT_COOLDOWN = time.Second * 20

var MAX_USER_BIO_LENGTH = 100

/* -------------------------------------------------------------------------- */
/*                                 USER ROLES                                 */
/* -------------------------------------------------------------------------- */
const (
	USER_ROLE_ADMIN  = "admin"
	USER_ROLE_MEMBER = "member"
)

/* -------------------------------------------------------------------------- */
/*                               Sorting Options                              */
/* -------------------------------------------------------------------------- */
const (
	SORT_BYRECENT = "commented_at DESC, id DESC"
	SORT_BYNEW    = "created_at  DESC, id DESC"
	SORT_BYHOT    = "comments_count DESC, commented_at DESC"
)
