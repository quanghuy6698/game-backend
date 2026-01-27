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

/* === Create Player === */
func CreatePlayer(ctx *gin.Context) {
	var req struct {
		PublicID string `json:"public_id" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    gin.H{},
			"error":   "Bad Request",
		})
		return
	}

	playerID := uuid.NewString()

	_, err := db.Exec(
		insertPlayerSQL,
		playerID,
		req.PublicID,
		req.Name,
	)

	if err != nil {
		log.Println("ERROR EXECUTE SQL CREATE PLAYER:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    gin.H{},
			"error":   "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    gin.H{},
	})
}

/* === Get top 100 === */
func GetTop100(c *gin.Context) {
	skin := c.Query("skin")
	if skin == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    gin.H{},
			"error":   "Skin is required",
		})
		return
	}

	rows, err := db.Query(getRanking100SQL, skin)
	if err != nil {
		log.Println("ERROR EXECUTE SQL GET TOP 100:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    gin.H{},
			"error":   "Internal Server Error",
		})
		return
	}
	defer rows.Close()

	result := make([]RankingItem, 0)

	for rows.Next() {
		var item RankingItem
		if err := rows.Scan(
			&item.PublicID,
			&item.Name,
			&item.Score,
			&item.Skin,
		); err != nil {
			log.Println("ERROR SCAN SQL RESULT TOP 100: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"data":    gin.H{},
				"error":   "Internal Server Error",
			})
			return
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		log.Println("ERROR GET TOP 100: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    gin.H{},
			"error":   "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

/* === Gen ID === */
func SaveScore(ctx *gin.Context) {
	var req struct {
		PlayerID string `json:"public_id" binding:"required"`
		Score    int64 	`json:"score" binding:"required"`
		Skin	 string `json:"skin" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    gin.H{},
			"error":   "Bad Request",
		})
		return
	}

	_, err := db.Exec(
		insertRankingSQL,
		req.PlayerID,
		req.Score,
		req.Skin,
	)

	if err != nil {
		log.Println("ERROR EXECUTE GEN PLAYER ID:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    gin.H{},
			"error":   "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    gin.H{},
	})
}
