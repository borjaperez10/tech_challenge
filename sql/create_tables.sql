
CREATE TABLE ISSUER(
	NIF VARCHAR PRIMARY KEY,
	NAME VARCHAR
);

CREATE SEQUENCE user_id_seq;

CREATE TABLE INVOICE(
	ID INTEGER  DEFAULT nextval('user_id_seq'),
	NAME VARCHAR,
	ISSUER VARCHAR,
	TOTAL FLOAT,
	TORECEIVE FLOAT,
	CLOSED VARCHAR,
	FOREIGN KEY (ISSUER) REFERENCES ISSUER (NIF),
	PRIMARY KEY (ID)
);


CREATE TABLE INVESTOR(
	DNI VARCHAR PRIMARY KEY,
	NAME VARCHAR,
	TOTAL_MONEY FLOAT,
	RETAINED_MONEY FLOAT
);

CREATE SEQUENCE part_id_seq;

CREATE TABLE PART_INVOICE(
	INVOICE_ID INTEGER,
	INVOICE_PART INTEGER DEFAULT nextval('part_id_seq'),
	TOTAL FLOAT,
	AMOUNT FLOAT,
	BUYER VARCHAR,
	SELLER VARCHAR,
	
	FOREIGN KEY (INVOICE_ID) REFERENCES INVOICE(ID),
	FOREIGN KEY (SELLER) REFERENCES ISSUER(NIF),
	FOREIGN KEY (BUYER) REFERENCES INVESTOR(DNI),
	PRIMARY KEY (INVOICE_PART)
);

CREATE SEQUENCE debt_id_seq;
CREATE TABLE DEBT(
	ID INTEGER DEFAULT nextval('debt_id_seq'),
	INVOICE_PART INTEGER,
	AMOUNT FLOAT,
	CREDITOR VARCHAR,
	DEBTOR VARCHAR,
	
	FOREIGN KEY (INVOICE_PART) REFERENCES PART_INVOICE(INVOICE_PART),
	FOREIGN KEY (CREDITOR) REFERENCES INVESTOR(DNI),
	FOREIGN KEY (DEBTOR) REFERENCES ISSUER(NIF),
	
	PRIMARY KEY (ID)
);


