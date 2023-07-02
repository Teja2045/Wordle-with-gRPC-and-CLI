package errors

import "errors"

var Err_UNKNOWN_ERROR = errors.New("unknown error")
var Err_RANKS_NO_ONE = errors.New("no ranks today")
var Err_DB_EMPTY = errors.New("database is not initialized")
var Err_RANK_GAME_LOST = errors.New("you lost and not ranked")
var Err_RANK_GAME_PENDING = errors.New("your game is still pending")
var Err_RANK_GAME_NOT_STARTED = errors.New("your game is yet started")
var Err_TRIES_EXHAUSTED = errors.New("no more tries left for the user")
var Err_GAME_ALREADY_WON = errors.New("games is already completed")
var Err_USER_DOES_NOT_EXIST = errors.New("the user with given user name doesn't exist")
var Err_USER_ALREADY_EXIST = errors.New("the user with given user name is already exists")
