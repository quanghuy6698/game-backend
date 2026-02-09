package tetGame

import (
	"log"
	"net/http"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./data/tet-game.db?_busy_timeout=5000&_journal_mode=WAL")
	if err != nil {
		log.Println("ERROR INIT DB:", err.Error())
		return err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	return db.Ping()
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

/* === Create Player === */
func CreatePlayer(ctx *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    gin.H{},
			"message": "bad_request",
		})
		return
	}

	playerID := uuid.NewString()

	_, err := db.Exec(
		insertPlayerSQL,
		playerID,
		req.Name,
	)

	if err != nil {
		log.Println("ERROR EXECUTE SQL CREATE PLAYER:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": 	false,
			"data": 	gin.H{},
			"message": "internal_server_error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id": playerID,
			"name": req.Name,
		},
		"message": "ok",
	})
}

/* === Get Leaderboard === */
func GetLeaderboard(c *gin.Context) {
	skin := c.Query("skin")
	if skin == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 	false,
			"data": 	gin.H{},
			"message": "skin_required",
		})
		return
	}

	rows, err := db.Query(getLeaderboardSQL, skin)
	if err != nil {
		log.Println("ERROR EXECUTE SQL GET LEADERBOARD:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": 	false,
			"data": 	gin.H{},
			"message": "internal_server_error",
		})
		return
	}
	defer rows.Close()

	result := make([]LeaderboardItem, 0)

	for rows.Next() {
		var item LeaderboardItem
		if err := rows.Scan(
			&item.Name,
			&item.Score,
			&item.Skin,
		); err != nil {
			log.Println("ERROR SCAN SQL RESULT LEADERBOARD: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": 	false,
				"data": 	gin.H{},
				"message": "internal_server_error",
			})
			return
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		log.Println("ERROR GET LEADERBOARD: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": 	false,
			"data": 	gin.H{},
			"message": "internal_server_error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": result,
		"message": "ok",
	})
}

/* === Save Score === */
func SaveScore(ctx *gin.Context) {
	var req struct {
		PlayerID string `json:"public_id" binding:"required"`
		Score    int64 	`json:"score" binding:"required"`
		Skin	 string `json:"skin" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": 	false,
			"data": 	gin.H{},
			"message": "bad_request",
		})
		return
	}

	_, err := db.Exec(
		insertScoreSQL,
		req.PlayerID,
		req.Score,
		req.Skin,
	)

	if err != nil {
		log.Println("ERROR EXECUTE GEN PLAYER ID:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    gin.H{},
			"message": "internal_server_error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    gin.H{},
		"message": "ok",
	})
}
