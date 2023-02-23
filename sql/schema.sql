------------------ PERSONAL DATA -----------------------
CREATE TABLE accounts (
	user_id INT,
	username VARCHAR ( 50 ) UNIQUE NOT NULL,
	password VARCHAR ( 250 ) NOT NULL,
	email VARCHAR ( 255 ) UNIQUE NOT NULL,
	firstname VARCHAR ( 50 ) NOT NULL,
	lastname VARCHAR ( 50 ) NOT NULL,
	address VARCHAR ( 50 ) NOT NULL,
	job VARCHAR ( 50 ) NOT NULL,
	age VARCHAR ( 50 ) NOT NULL,
	amka VARCHAR ( 50 ) NOT NULL,
	created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP,
	PRIMARY KEY(user_id)
);

------------------ HOW DID YOU FIND US -----------------
CREATE TABLE reach_out (
	reach_id INT,
	user_id INT,
	how_know VARCHAR (250), 
	person BOOLEAN,
	person_firstname VARCHAR (50),
	person_lastname VARCHAR (50),
	website VARCHAR (50),
	created_on TIMESTAMP NOT NULL,
	PRIMARY KEY(reach_id),
	CONSTRAINT fk_user_id
	FOREIGN KEY(user_id)
	REFERENCES accounts(user_id)
	ON DELETE NO ACTION
);


------------------ THERAPIES OF ANY TYPE ---------------
CREATE TABLE therapies (
	therapy_id INT,
	user_id INT,
	therapy_title VARCHAR (50),
	therapy_description VARCHAR (100),
	diagnosis VARCHAR (250),
	from_date TIMESTAMP,
	to_date TIMESTAMP,
	quantity SMALLINT,
	frequency SMALLINT,
	created_on TIMESTAMP NOT NULL,
	PRIMARY KEY(therapy_id),
	CONSTRAINT fk_user_id
	FOREIGN KEY(user_id)
	REFERENCES accounts(user_id)
	ON DELETE NO ACTION
);

CREATE TABLE injuries (
	injury_id INT,
	user_id INT,
	injury_title VARCHAR (50),
	injury_description VARCHAR (100),
	injury_date TIMESTAMP,
	bposition_id INT,
	created_on TIMESTAMP NOT NULL,
	PRIMARY KEY(injury_id),
	CONSTRAINT fk_user_id
	FOREIGN KEY(user_id)
	REFERENCES accounts(user_id)
	ON DELETE NO ACTION
);

CREATE TABLE medical_therapies (
	medical_therapy_id INT,
	user_id INT,
	medical_therapy_title VARCHAR (50),
	medical_therapy_description VARCHAR (100),
	from_date TIMESTAMP,
	to_date TIMESTAMP,
	quantity SMALLINT,
	frequency SMALLINT,
	created_on TIMESTAMP NOT NULL,
	PRIMARY KEY(medical_therapy_id),
	CONSTRAINT fk_user_id
	FOREIGN KEY(user_id)
	REFERENCES accounts(user_id)
	ON DELETE NO ACTION
);

CREATE TABLE personal_allergies (
	personal_allergies_id INT,
	allergy_id INT,
	user_id INT,
	diagnosed_time TIMESTAMP,
	created_on TIMESTAMP NOT NULL,
	PRIMARY KEY(personal_allergies_id),
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
	drug_treatment_id INT,
	drug_id INT,
	user_id INT,
	from_date TIMESTAMP,
	to_date TIMESTAMP,
	quantity SMALLINT,
	frequency SMALLINT,
	created_on TIMESTAMP NOT NULL,
	PRIMARY KEY(drug_treatment_id),
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
	personal_disorder_id INT,
	disorder_id INT,
	user_id INT,
	from_date TIMESTAMP,
	to_date TIMESTAMP,
	quantity SMALLINT,
	frequency SMALLINT,
	created_on TIMESTAMP NOT NULL,
	PRIMARY KEY(personal_disorder_id),
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
	allergy_id INT,
	allergy_title VARCHAR (50),
	allergy_description VARCHAR (100),
	allergy_notes VARCHAR (200),
	PRIMARY KEY(allergy_id)
);

CREATE TABLE body_positions (
	bposition_id INT,
	bposition_title VARCHAR (50),
	bposition_description VARCHAR (100),
	bposition_notes VARCHAR (150),
	PRIMARY KEY(bposition_id)
);

CREATE TABLE drugs (
	drug_id INT,
	drug_title VARCHAR (50),
	drug_description VARCHAR (100),
	drug_notes VARCHAR (200),
	PRIMARY KEY(drug_id),
);

CREATE TABLE disorders (
	disorder_id INT,
	disorder_title VARCHAR (50),
	disorder_description VARCHAR (100),
	disorder_notes VARCHAR (200),
	PRIMARY KEY(disorder_id)
);
