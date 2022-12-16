package controllers

import "time"

/* -------------------------------------------------------------------------- */
/*                              Config variables                              */
/* -------------------------------------------------------------------------- */
var MAX_TITLE_CHAR = 100
var MAX_PER_PAGE = float64(50)
var USER_POST_COOLDOWN = -time.Minute * 1
var USER_COMMENT_COOLDOWN = -time.Minute * 1
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                               Sorting Options                              */
/* -------------------------------------------------------------------------- */
const (
	ByRecent 	= "commented_at DESC, id DESC"
	ByNew 		= "created_at  DESC, id DESC"
	ByHot 		= "comments_count DESC, commented_at DESC"
)