package server

const (
	// Create User table
	CreateTableUsers = `
	CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		sessions TEXT NOT NULL
    );`
	// Create Post table
	CreateTablePost = `
	CREATE TABLE IF NOT EXISTS Post (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		categories TEXT,
		date_created TEXT,
		title TEXT NOT NULL,
		content TEXT,
		image TEXT,
		id_user INTEGER,
		FOREIGN KEY (id_user) REFERENCES User(ID)
	);`
	// Create Engagement table
	CreateTableEngagement = `
	CREATE TABLE IF NOT EXISTS Engagement (
    	ID INTEGER PRIMARY KEY AUTOINCREMENT,
    	like INTEGER DEFAULT 0,
    	dislike INTEGER DEFAULT 0,
    	category TEXT,
    	id_post INTEGER,
    	FOREIGN KEY (id_post) REFERENCES Post(ID)
	);`
	LoginQuery     = `SELECT password FROM users WHERE username = ? OR email = ?`
	AddUserQuery   = `INSERT INTO users(username, password, email) VALUES (?, ?, ?)`
	ValidUserQuery = `SELECT username FROM users WHERE username = ?`
)

// func Select()string{

// 	return ""
// }
