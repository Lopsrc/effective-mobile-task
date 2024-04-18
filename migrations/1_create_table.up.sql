CREATE TABLE IF NOT EXISTS person(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255) NOT NULL,
    del BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS car(
    id SERIAL PRIMARY KEY,
    reg_number VARCHAR(255) UNIQUE NOT NULL,
    mark VARCHAR(255) DEFAULT 'not stated',
    model VARCHAR(255) DEFAULT 'not stated',
    year VARCHAR(255) DEFAULT 'not stated',
    owner_id INTEGER NOT NULL,
    del BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(owner_id) REFERENCES person(id)
);
