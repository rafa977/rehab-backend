CREATE TABLE accounts (
	user_id serial PRIMARY KEY,
	username VARCHAR ( 50 ) UNIQUE NOT NULL,
	password VARCHAR ( 50 ) NOT NULL,
	email VARCHAR ( 255 ) UNIQUE NOT NULL,
	firstname VARCHAR ( 50 ) NOT NULL,
	lastname VARCHAR ( 50 ) NOT NULL,
	address VARCHAR ( 50 ) NOT NULL,
	job VARCHAR ( 50 ) NOT NULL,
	age VARCHAR ( 50 ) NOT NULL,
	amka VARCHAR ( 50 ) NOT NULL,
	created_on TIMESTAMP NOT NULL,
        last_login TIMESTAMP
);
