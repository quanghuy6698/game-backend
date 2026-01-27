package tetGame

const insertPlayerSQL = `INSERT INTO player (id, public_id, name) VALUES (?, ?, ?)`
const insertRankingSQL = `INSERT INTO ranking (player_id, score, skin) VALUES (?, ?, ?)`
const getRanking100SQL = `
	SELECT 
		p.public_id,
		p.name,
		r.score,
		r.skin 
	FROM ranking r 
	JOIN player p 
	ON r.player_id = p.id 
	WHERE r.skin = ? 
	ORDER BY r.score DESC 
	LIMIT 100
`
