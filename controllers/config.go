package controllers

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

/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                               Sorting Options                              */
/* -------------------------------------------------------------------------- */
const (
	ByRecent = "commented_at DESC, id DESC"
	ByNew    = "created_at  DESC, id DESC"
	ByHot    = "comments_count DESC, commented_at DESC"
)
