package basicthreads

import (
	_ "basicthreads/internal/database"
)

type Response struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// path: /sms
