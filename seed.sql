CREATE DATABASE employees;

\c employees

CREATE TABLE employees_data (
    id SERIAL PRIMARY KEY,
    name TEXT,
    surname TEXT,
    birth_date DATE,
    entry_date DATE
);

INSERT INTO employees_data (name, surname, birth_date, entry_date) VALUES
('Alice', 'Johnson', '1985-03-10', '2010-07-15'),
('Bob', 'Smith', '1980-01-22', '2008-04-12'),
('Carol', 'White', '1992-11-05', '2019-03-19'),
('David', 'Brown', '1988-06-01', '2011-11-23'),
('Eve', 'Black', '1990-08-14', '2016-09-07'),
('Frank', 'Miller', '1983-05-17', '2009-01-02'),
('Grace', 'Clark', '1995-02-28', '2022-05-14'),
('Hank', 'Adams', '1982-12-09', '2006-10-19'),
('Ivy', 'Davis', '1991-07-25', '2018-06-20'),
('Jack', 'Wilson', '1987-03-03', '2013-12-30');
