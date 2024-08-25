CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(255) PRIMARY KEY,
    status int(11) DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS loan (
    id int(11) PRIMARY KEY,
    username varchar(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pay_loan(
    id int(11) PRIMARY KEY,
    loan_id int(11) NOT NULL,
    amount int(11) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);