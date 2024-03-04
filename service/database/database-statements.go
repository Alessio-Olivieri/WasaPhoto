package database

const (
	usersTableCreationStatement = `
	CREATE TABLE "Users" (
		"user_id"	INTEGER NOT NULL UNIQUE,
		"username"	TEXT NOT NULL UNIQUE,
		"date of birth"	TEXT,
		PRIMARY KEY("user_id" AUTOINCREMENT),
		CHECK (length("username") >= 3 AND length("username") <= 20)
	);
	`
	photosTableCreationStatement = `
	CREATE TABLE "Photos" (
		"photo_id"	INTEGER UNIQUE,
		"user_id"	INTEGER NOT NULL,
		"text"	TEXT CHECK(length("text")<=5000),
		"like_count"	INTEGER DEFAULT 0,
		"image"	BLOB,
		"date"	TEXT,
		PRIMARY KEY("photo_id" AUTOINCREMENT),
		FOREIGN KEY("user_id") REFERENCES "Users"("user_id")
	);
	`
	commentsTableCreationStatement = `
	CREATE TABLE "Comments" (
		"comment_id"	INTEGER UNIQUE,
		"photo_id"	INTEGER NOT NULL,
		"user_id"	INTEGER NOT NULL,
		"text"	TEXT CHECK(length("text")<=500),
		"date"	TEXT,
		PRIMARY KEY("comment_id" AUTOINCREMENT)
		FOREIGN KEY("user_id") REFERENCES "Users"("user_id")
		FOREIGN KEY("photo_id") REFERENCES "Photos"("photo_id")
	);
	`
	likesTableCreationStatement = `
	CREATE TABLE "Likes" (
		"photo_id"	INTEGER NOT NULL,
		"user_id"	INTEGER NOT NULL,
		PRIMARY KEY("photo_id", "user_id")
		FOREIGN KEY("user_id") REFERENCES "Users"("user_id")
		FOREIGN KEY("photo_id") REFERENCES "Photos"("photo_id")
	);
	`
	followersTableCreationStatement = `
	CREATE TABLE "Followers" (
		"follower_id"	INTEGER UNIQUE,
		"user_id"	INTEGER NOT NULL,
		"follower_user_id"	INTEGER NOT NULL,
		PRIMARY KEY("follower_id" AUTOINCREMENT)
		FOREIGN KEY("user_id") REFERENCES "Users"("user_id")
		FOREIGN KEY("follower_user_id") REFERENCES "Users"("user_id")
	);
	`
)
