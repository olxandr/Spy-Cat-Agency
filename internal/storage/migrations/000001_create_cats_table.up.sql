CREATE TABLE cats (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    years_of_experience INT NOT NULL CHECK (years_of_experience >= 0),
    breed VARCHAR(50) NOT NULL,
    salary NUMERIC(10,2) NOT NULL CHECK (salary >= 0)
);

