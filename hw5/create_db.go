package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func create_db() {
	os.Remove(db_filename)

	db, err := sql.Open("sqlite3", db_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    // create the necessary tables
	sqlStmt := `
    CREATE TABLE Employee
    (
      SSN CHAR(9) NOT NULL,
      Name VARCHAR(255) NOT NULL,
      start_date DATE NOT NULL,
      end_date DATE,
      PRIMARY KEY (SSN)
    );

    CREATE TABLE Customer
    (
      Email VARCHAR(255) NOT NULL,
      Name VARCHAR(255) NOT NULL,
      House_Number VARCHAR(10) NOT NULL,
      Street VARCHAR(255) NOT NULL,
      City VARCHAR(255) NOT NULL,
      State CHAR(2) NOT NULL,
      Zip CHAR(5) NOT NULL,
      PRIMARY KEY (Email)
    );

    CREATE TABLE Customer_Phone_Number
    (
      Phone_Number CHAR(10) NOT NULL,
      Email VARCHAR(255) NOT NULL,
      PRIMARY KEY (Phone_Number, Email),
      FOREIGN KEY (Email) REFERENCES Customer(Email)
    );

    CREATE TABLE Teaches
    (
      Teacher_SSN CHAR(9) NOT NULL,
      Student_Email VARCHAR(255) NOT NULL,
      PRIMARY KEY (Teacher_SSN, Student_Email),
      FOREIGN KEY (Teacher_SSN) REFERENCES Employee(SSN),
      FOREIGN KEY (Student_Email) REFERENCES Customer(Email)
    );

    CREATE TABLE Bike
    (
      ID INT NOT NULL,
      Status VARCHAR(255) NOT NULL,
      Dollar_value FLOAT NOT NULL,
      Purchase_Time DATETIME NOT NULL,
      Purchaser_Email VARCHAR(255) NOT NULL,
      PRIMARY KEY (ID),
      FOREIGN KEY (Purchaser_Email) REFERENCES Customer(Email)
    );

    CREATE TABLE Repairs
    (
      Bike_ID INT NOT NULL,
      Employee_SSN CHAR(9) NOT NULL,
      PRIMARY KEY (Bike_ID, Employee_SSN),
      FOREIGN KEY (Bike_ID) REFERENCES Bike(ID),
      FOREIGN KEY (Employee_SSN) REFERENCES Employee(SSN)
    );
	`

    // create a customer view and a management view
    sqlStmt += `
    CREATE VIEW Customer_View AS
    SELECT * FROM Customer
    JOIN Customer_Phone_Number ON Customer.Email = Customer_Phone_Number.Email
    JOIN Bike ON Customer.Email = Bike.Purchaser_Email;

    CREATE VIEW Management_View AS
    SELECT * FROM Employee
    JOIN Teaches ON Employee.SSN = Teaches.Teacher_SSN
    JOIN Repairs ON Employee.SSN = Repairs.Employee_SSN;
    `

    // insert some dummy data
    sqlStmt += `
    INSERT INTO Employee VALUES ('123456789', 'Matthew Leeds', '2014-01-20', NULL);
    INSERT INTO Employee VALUES ('123456788', 'Andrew Leeds', '2013-01-20', NULL);
    INSERT INTO Employee VALUES ('123456787', 'Tait Wayland', '2015-05-20', NULL);
    INSERT INTO Employee VALUES ('123456786', 'Justin Harrison', '2011-09-10', NULL);

    INSERT INTO Customer VALUES ('james.baker@protonmail.com', 'James Baker', '128', 'Hacker St', 'Birmingham', 'AL', '35223');
    INSERT INTO Customer VALUES ('elon.musk@protonmail.com', 'Elon Musk', '255', 'Hacker St', 'San Francisco', 'CA', '94101');
    INSERT INTO Customer VALUES ('george.hotz@protonmail.com', 'George Hotz', '42', 'Comma St', 'San Francisco', 'CA', '94101');
    INSERT INTO Customer VALUES ('parker.higgins@protonmail.com', 'Parker Higgins', '2', 'Activist St', 'Oakland', 'CA', '94601');
    INSERT INTO Customer VALUES ('cory.doctorow@protonmail.com', 'Cory Doctorow', '8', 'Activist St', 'Huntsville', 'AL', '94101');
    INSERT INTO Customer VALUES ('satoshi.nakamoto@protonmail.com', 'Satoshi Nakamoto', '4', 'Innovator St', 'Westford', 'MA', '94101');
    INSERT INTO Customer VALUES ('timothy.lee@protonmail.com', 'Timothy B Lee', '9', 'Innovator St', 'San Francisco', 'CA', '94101');
    INSERT INTO Customer VALUES ('roger.ver@protonmail.com', 'Roger K Ver', '10', 'Bitcoin St', 'San Francisco', 'CA', '94101');

    INSERT INTO Customer_Phone_Number VALUES ('9996084744', 'james.baker@protonmail.com');
    INSERT INTO Customer_Phone_Number VALUES ('3652994596', 'elon.musk@protonmail.com');
    INSERT INTO Customer_Phone_Number VALUES ('2115403124', 'george.hotz@protonmail.com');
    INSERT INTO Customer_Phone_Number VALUES ('8417249108', 'parker.higgins@protonmail.com');
    INSERT INTO Customer_Phone_Number VALUES ('1414876813', 'cory.doctorow@protonmail.com');
    INSERT INTO Customer_Phone_Number VALUES ('3438499396', 'satoshi.nakamoto@protonmail.com');
    INSERT INTO Customer_Phone_Number VALUES ('4057720448', 'timothy.lee@protonmail.com');
    INSERT INTO Customer_Phone_Number VALUES ('1952411513', 'roger.ver@protonmail.com');

    INSERT INTO Teaches VALUES ('123456789', 'satoshi.nakamoto@protonmail.com');
    INSERT INTO Teaches VALUES ('123456789', 'roger.ver@protonmail.com');
    INSERT INTO Teaches VALUES ('123456788', 'elon.musk@protonmail.com');
    INSERT INTO Teaches VALUES ('123456786', 'george.hotz@protonmail.com');
    INSERT INTO Teaches VALUES ('123456786', 'james.baker@protonmail.com');
    INSERT INTO Teaches VALUES ('123456788', 'james.baker@protonmail.com');
    INSERT INTO Teaches VALUES ('123456789', 'james.baker@protonmail.com');

    INSERT INTO Bike VALUES (1, 'broken gear', 500, '2010-04-12 07:00:00', 'cory.doctorow@protonmail.com');
    INSERT INTO Bike VALUES (2, 'good', 1200, '2015-09-12 07:00:00', 'cory.doctorow@protonmail.com');
    INSERT INTO Bike VALUES (3, 'brake needs replacement', 1500, '2013-04-12 07:00:00', 'parker.higgins@protonmail.com');
    INSERT INTO Bike VALUES (4, 'flat tire', 300, '2013-04-12 08:00:00', 'parker.higgins@protonmail.com');
    INSERT INTO Bike VALUES (5, 'flat tire', 1800, '2013-04-12 07:00:00', 'elon.musk@protonmail.com');
    INSERT INTO Bike VALUES (6, 'flat tire', 1800, '2009-04-12 08:00:00', 'roger.ver@protonmail.com');
    INSERT INTO Bike VALUES (7, 'flat tire', 1400, '2014-04-12 12:00:00', 'timothy.lee@protonmail.com');
    INSERT INTO Bike VALUES (8, 'flat tire', 10000, '2008-04-02 11:00:00', 'timothy.lee@protonmail.com');
    INSERT INTO Bike VALUES (9, 'flat tire', 900, '2013-04-12 07:00:00', 'satoshi.nakamoto@protonmail.com');
    INSERT INTO Bike VALUES (10, 'mangled', 700, '2013-04-12 07:00:00', 'george.hotz@protonmail.com');
    INSERT INTO Bike VALUES (11, 'good', 550, '2013-04-12 07:00:00', 'james.baker@protonmail.com');
    INSERT INTO Bike VALUES (12, 'rusty chain', 600, '2013-01-12 10:30:00', 'james.baker@protonmail.com');

    INSERT INTO Repairs VALUES (1, '123456788');
    INSERT INTO Repairs VALUES (3, '123456788');
    INSERT INTO Repairs VALUES (4, '123456788');
    INSERT INTO Repairs VALUES (5, '123456787');
    INSERT INTO Repairs VALUES (6, '123456787');
    INSERT INTO Repairs VALUES (7, '123456789');
    INSERT INTO Repairs VALUES (8, '123456789');
    INSERT INTO Repairs VALUES (9, '123456786');
    INSERT INTO Repairs VALUES (10, '123456786');
    INSERT INTO Repairs VALUES (10, '123456787');
    INSERT INTO Repairs VALUES (10, '123456788');
    INSERT INTO Repairs VALUES (10, '123456789');
    INSERT INTO Repairs VALUES (12, '123456788');
    `

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
