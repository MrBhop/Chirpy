package handlers

import (
	"sync/atomic"
	"github.com/MrBhop/Chirpy/internal/database"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	Db             *database.Queries
	Platform       string
	Secret         string
}
