package constants

import (
	"time"
)

const ACCESS_TOKEN_EXP_TIME = time.Minute * 20
const REFRESH_TOKEN_EXP_TIME = time.Hour * 1 
const JWT_SIGNING_ALGORITHM = "HS256"