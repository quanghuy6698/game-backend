package tetGame

const insertPlayerSQL = `INSERT INTO player (id, public_id, name) VALUES (?, ?, ?)`
const insertRankingSQL = `INSERT INTO ranking (player_id, score, skin) VALUES (?, ?, ?)`
const getLeaderboardSQL = `
	SELECT 
		p.public_id,
		p.name,
		s.score,
		s.skin 
	FROM score s 
	JOIN player p 
	ON s.player_id = p.id 
	WHERE s.skin = ?
	ORDER BY s.score DESC 
	LIMIT 20;
`
