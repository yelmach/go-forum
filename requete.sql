-- Users table
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
		);

		-- Sessions
		CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		session_id TEXT NOT NULL UNIQUE,
		FOREIGN KEY (user_id) REFERENCES users(id)
		);

		-- Posts table
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
		
		-- Comments table
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id)
		);
		
		-- Likes table
	CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		is_dislike BOOLEAN NOT NULL,
		is_comment BOOLEAN NOT NULL,
		user_id INTEGER NOT NULL,
		post_id INTEGER,
		comment_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		CHECK ((post_id IS NOT NULL AND comment_id IS NULL) OR (post_id IS NULL AND comment_id IS NOT NULL))
		);
		
	CREATE TABLE IF NOT EXISTS categories(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		categori TEXT NOT NULL,
		post_id INTEGER,
		FOREIGN KEY (post_id) REFERENCES posts(id));

-- Sample data for posts table
INSERT INTO posts (title, content, user_id) VALUES 
    ("Exploring the Blue Pearl of Morocco", "Chefchaouen, known as the Blue Pearl, is a beautiful city in Morocco painted in shades of blue. It attracts photographers and travelers from all over the world.", 1),
    ("A Journey Through the Sahara", "The Sahara Desert offers a breathtaking experience, with camel treks and unforgettable sunsets over the golden dunes.", 2),
    ("Top 5 Must-See Places in Marrakech", "From the vibrant souks to the serene Majorelle Garden, Marrakech has it all. Here’s a guide to the top spots you must visit!", 3),
    ("Hiking the High Atlas Mountains", "The High Atlas Mountains offer some of the most stunning hiking trails in Morocco. Toubkal, the highest peak, is a popular destination for hikers.", 4),
    ("Discovering the Ancient Medina of Fes", "Fes is home to one of the world's oldest medinas. It's a maze of narrow alleys, historic buildings, and cultural heritage.", 5),
    ("Culinary Delights of Morocco", "Moroccan cuisine is rich and flavorful. From tagine to couscous, every dish tells a story of tradition and culture.", 1);
-- Sample data for comments table
INSERT INTO comments (content, user_id, post_id) VALUES 
    ("This is amazing! Chefchaouen has been on my bucket list for ages.", 2, 1),
    ("I've been to the Sahara! The desert at night is unforgettable.", 3, 2),
    ("Thanks for the tips! Marrakech sounds like a dream.", 4, 3),
    ("I want to hike Toubkal someday! Great post.", 5, 4),
    ("The Medina of Fes is truly a cultural treasure. Can't wait to visit again.", 1, 5),
    ("Moroccan food is the best! Can’t get enough of it.", 3, 6),
    ("Thanks for sharing your experience in Chefchaouen!", 4, 1),
    ("The desert feels so serene, your description really takes me back.", 5, 2),
    ("Marrakech sounds incredible! Definitely adding it to my travel plans.", 1, 3);

-- Sample data for categories table
INSERT INTO categories (categori, post_id) VALUES 
    ("Travel", 1),
    ("Adventure", 2),
    ("Guides", 3),
    ("Nature", 4),
    ("Culture", 5),
    ("Food", 6),
    ("Photography", 1),
    ("History", 5),
    ("Tips", 3);

