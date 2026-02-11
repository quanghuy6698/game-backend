package tetGame

const insertPlayerSQL = `INSERT INTO player (id, name) VALUES (?, ?)`
const insertScoreSQL = `INSERT INTO score (player_id, score, skin) VALUES (?, ?, ?)`
const getLeaderboardSQL = `
	SELECT 
		p.id,
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
