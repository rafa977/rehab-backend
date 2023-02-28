------------------ PERSONAL DATA -----------------------
CREATE TABLE accounts (
	user_id SERIAL PRIMARY KEY,
	username VARCHAR ( 50 ) UNIQUE NOT NULL,
	password VARCHAR ( 250 ) NOT NULL,
	email VARCHAR ( 255 ) UNIQUE NOT NULL,
	user_type VARCHAR (10),
	firstname VARCHAR ( 50 ) ,
	lastname VARCHAR ( 50 ) ,
	address VARCHAR ( 50 ) ,
	job VARCHAR ( 50 ) ,
	age VARCHAR ( 50 ) ,
	amka VARCHAR ( 50 ) ,
	created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP
);

------------------ HOW DID YOU FIND US -----------------
CREATE TABLE reach_out (
	reach_id SERIAL PRIMARY KEY,
	user_id INT,
	how_know VARCHAR (250), 
	person BOOLEAN,
	person_firstname VARCHAR (50),
	person_lastname VARCHAR (50),
	website VARCHAR (50),
	created_on TIMESTAMP NOT NULL,
	CONSTRAINT fk_user_id
	FOREIGN KEY(user_id)
	REFERENCES accounts(user_id)
	ON DELETE NO ACTION
);


------------------ THERAPIES OF ANY TYPE ---------------
CREATE TABLE therapies (
	therapy_id SERIAL PRIMARY KEY,
	user_id INT,
	therapy_title VARCHAR (50),
	therapy_description VARCHAR (100),
	diagnosis VARCHAR (250),
	from_date TIMESTAMP,
	to_date TIMESTAMP,
	quantity SMALLINT,
	frequency SMALLINT,
	created_on TIMESTAMP NOT NULL,
	CONSTRAINT fk_user_id
	FOREIGN KEY(user_id)
	REFERENCES accounts(user_id)
	ON DELETE NO ACTION
);

CREATE TABLE injuries (
	injury_id SERIAL PRIMARY KEY,
	user_id INT,
	injury_title VARCHAR (50),
	injury_description VARCHAR (100),
	injury_date TIMESTAMP,
	bposition_id INT,
	created_on TIMESTAMP NOT NULL,
	CONSTRAINT fk_user_id
	FOREIGN KEY(user_id)
	REFERENCES accounts(user_id)
	ON DELETE NO ACTION
);

CREATE TABLE medical_therapies (
	medical_therapy_id SERIAL PRIMARY KEY,
	user_id INT,
	medical_therapy_title VARCHAR (50),
	medical_therapy_description VARCHAR (100),
	from_date TIMESTAMP,
	to_date TIMESTAMP,
	quantity SMALLINT,
	frequency SMALLINT,
	created_on TIMESTAMP NOT NULL,
	CONSTRAINT fk_user_id
	FOREIGN KEY(user_id)
	REFERENCES accounts(user_id)
	ON DELETE NO ACTION
);

CREATE TABLE personal_allergies (
	personal_allergies_id SERIAL PRIMARY KEY,
	allergy_id INT,
	user_id INT,
	diagnosed_time TIMESTAMP,
	created_on TIMESTAMP NOT NULL,
	CONSTRAINT fk_allergy_id
	FOREIGN KEY(allergy_id)
	REFERENCES allergies(allergy_id)
	ON DELETE NO ACTION,
	CONSTRAINT fk_user_id
	FOREIGN KEY(user_id)
	REFERENCES accounts(user_id)
	ON DELETE NO ACTION
);

CREATE TABLE drug_treatment (
	drug_treatment_id SERIAL PRIMARY KEY,
	drug_id INT,
	user_id INT,
	from_date TIMESTAMP,
	to_date TIMESTAMP,
	quantity SMALLINT,
	frequency SMALLINT,
	created_on TIMESTAMP NOT NULL,
	CONSTRAINT fk_drug_id
	FOREIGN KEY(drug_id)
	REFERENCES drugs(drug_id)
	ON DELETE NO ACTION,
	CONSTRAINT fk_user_id
	FOREIGN KEY(user_id)
	REFERENCES accounts(user_id)
	ON DELETE NO ACTION
);

CREATE TABLE personal_disorders (
	personal_disorder_id SERIAL PRIMARY KEY,
	disorder_id INT,
	user_id INT,
	from_date TIMESTAMP,
	to_date TIMESTAMP,
	quantity SMALLINT,
	frequency SMALLINT,
	created_on TIMESTAMP NOT NULL,
	CONSTRAINT fk_disorder_id
	FOREIGN KEY(disorder_id)
	REFERENCES disorders(disorder_id)
	ON DELETE NO ACTION,
	CONSTRAINT fk_user_id
	FOREIGN KEY(user_id)
	REFERENCES accounts(user_id)
	ON DELETE NO ACTION
);
--------------------------------------------------------

CREATE TABLE allergies (
	allergy_id SERIAL PRIMARY KEY,
	allergy_title VARCHAR (50),
	allergy_description VARCHAR (100),
	allergy_notes VARCHAR (200)
);

CREATE TABLE body_positions (
	bposition_id SERIAL PRIMARY KEY,
	bposition_title VARCHAR (50),
	bposition_description VARCHAR (100),
	bposition_notes VARCHAR (150)
);

CREATE TABLE drugs (
	drug_id SERIAL PRIMARY KEY,
	drug_title VARCHAR (50),
	drug_description VARCHAR (100),
	drug_notes VARCHAR (200)
);

CREATE TABLE disorders (
	disorder_id SERIAL PRIMARY KEY,
	disorder_title VARCHAR (50),
	disorder_description VARCHAR (100),
	disorder_notes VARCHAR (200)
);
