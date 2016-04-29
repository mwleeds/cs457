package main

import (
    "fmt"
    "net/http"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
    "log"
)

const db_filename string = "bike_shop.db"

func main() {
    // create a sqlite db with fake data
    create_db()
    // create an http server with handlers for each function
    http.HandleFunc("/", mainHandler)
    http.HandleFunc("/get_table/", getTableHandler)
    http.HandleFunc("/find_expensive_bikes/", findExpensiveBikesHandler)
    http.HandleFunc("/bikes_by_state/", bikesByStateHandler)
    http.HandleFunc("/insert_bike/", insertBikeHandler)
    http.HandleFunc("/exec_select/", execSelectHandler)
    http.ListenAndServe(":8080", nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
    // HTML for the home page with a form for each function
    fmt.Fprintf(w, `<h1>CS 457 HW 5 by Matthew Leeds</h1>
                    <br>
                    <ol><li><form action="/get_table/" method="POST">
                    SELECT * FROM
                    <select name="table">
                      <option value="Employee">Employee</option>
                      <option value="Teaches">Teaches</option>
                      <option value="Repairs">Repairs</option>
                      <option value="Bike">Bike</option>
                      <option value="Customer">Customer</option>
                      <option value="Customer_Phone_Number">Customer_Phone_Number</option>
                    </select>
                    &nbsp;&nbsp;<input type="submit" value="Execute">
                    </form></li>
                    <br>
                    <li><form action="/find_expensive_bikes/" method="POST">
                    SELECT * FROM Bike WHERE Bike.Dollar_value > 
                    <input type="text" name="bike_value_floor">
                    &nbsp;&nbsp;<input type="submit" value="Execute">
                    <br>(Enter a non-negative decimal value.)
                    </form></li>
                    <br>
                    <li><form action="/bikes_by_state/" method="POST">
                    SELECT Bike.ID, Bike.Dollar_value, Bike.Purchaser_Email, Customer.State<br>
                    FROM Bike JOIN Customer ON Bike.Purchaser_Email = Customer.Email<br>
                    WHERE Customer.State = 
                    <input type="text" name="bike_purchaser_state">
                    &nbsp;&nbsp;<input type="submit" value="Execute">
                    <br>(Enter a two letter state code like AL or CA.)
                    </form></li>
                    <br>
                    <li><form action="/insert_bike/" method="POST">
                    INSERT INTO Bike VALUES (<br>
                    <input type="text" name="bike_id">,&nbsp;&nbsp;(Enter an integer primary key not already used by another bike.)<br>
                    <input type="text" name="bike_status">,&nbsp;&nbsp;(Enter an explanation of the bike's status in less than 255 characters.)<br>
                    <input type="text" name="bike_dollar_value">,&nbsp;&nbsp;(Enter the bike's dollar value as a decimal number.)<br>
                    <input type="text" name="bike_purchase_time">,&nbsp;&nbsp;(Enter the bike's purchase time as YYYY-MM-DD HH:MM:SS.)<br>
                    <input type="text" name="bike_purchaser_email">,&nbsp;&nbsp;(Enter the bike's purchaser email address foreign key.)<br>
                    )<br>
                    &nbsp;&nbsp;<input type="submit" value="Execute">
                    </form></li>
                    <br>
                    <li><form action="/exec_select/" method="POST">
                    SELECT <input type="text" name="select_query">
                    &nbsp;&nbsp;<input type="submit" value="Execute">
                    </form></li>
                    </ol>`)
}

// execute a SELECT * on the specified table,
// and format the results as an HTML table
func getTableHandler(w http.ResponseWriter, r *http.Request) {
    table_name := r.FormValue("table")

	db, err := sql.Open("sqlite3", db_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    // of course if this were production code we would sanitize the input
	sqlStmt := "SELECT * FROM " + table_name + ";"

	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

    html := `<style>table { border-collapse: collapse; border-spacing: 0px; } 
                    table, th, td { padding: 5px; border: 1px solid black; }</style>
             <h1>CS 457 HW 5 by Matthew Leeds</h1>
             <h3>Table: ` + table_name + `</h3>
             <table>`

    // Ideally the below code would be generalized to work for any table.
    switch table_name {
        case "Employee":
            html += "<tr><th>SSN</th><th>Name</th><th>start_date</th><th>end_date</th></tr>"
            for rows.Next() {
                var ssn string
                var name string
                var start_date string
                var end_date sql.NullString
                err := rows.Scan(&ssn, &name, &start_date, &end_date)
                if err != nil {
                    log.Fatal(err)
                }
                html += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
                                    ssn, name, start_date, end_date.String)
            }

        case "Teaches":
            html += "<tr><th>Teacher_SSN</th><th>Student_Email</th></tr>"
            for rows.Next() {
                var teacher_ssn string
                var student_email string
                err := rows.Scan(&teacher_ssn, &student_email)
                if err != nil {
                    log.Fatal(err)
                }
                html += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", teacher_ssn, student_email)
            }

        case "Repairs":
            html += "<tr><th>Bike_ID</th><th>Employee_SSN</th></tr>"
            for rows.Next() {
                var bike_id int
                var employee_ssn string
                err := rows.Scan(&bike_id, &employee_ssn)
                if err != nil {
                    log.Fatal(err)
                }
                html += fmt.Sprintf("<tr><td>%d</td><td>%s</td></tr>", bike_id, employee_ssn)
            }

        case "Bike":
            html += "<tr><th>ID</th><th>Status</th><th>Dollar_value</th><th>Purchase_Time</th><th>Purchaser_Email</th></tr>"
            for rows.Next() {
                var id int
                var status string
                var dollar_value float32
                var purchase_time string
                var purchaser_email string
                err := rows.Scan(&id, &status, &dollar_value, &purchase_time, &purchaser_email)
                if err != nil {
                    log.Fatal(err)
                }
                html += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%.2f</td><td>%s</td><td>%s</td></tr>",
                                    id, status, dollar_value, purchase_time, purchaser_email)
            }

        case "Customer":
            html += "<tr><th>Email</th><th>Name</th><th>House_Number</th><th>Street</th><th>City</th><th>State</th><th>Zip</th></tr>"
            for rows.Next() {
                var email string
                var name string
                var house_number string
                var street string
                var city string
                var state string
                var zip string
                err := rows.Scan(&email, &name, &house_number, &street, &city, &state, &zip)
                if err != nil {
                    log.Fatal(err)
                }
                html += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
                                    email, name, house_number, street, city, state, zip)
            }

        case "Customer_Phone_Number":
            html += "<tr><th>Phone_Number</th><th>Email</th></tr>"
            for rows.Next() {
                var phone_number string
                var email string
                err := rows.Scan(&phone_number, &email)
                if err != nil {
                    log.Fatal(err)
                }
                html += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", phone_number, email)
            }

        default:
            log.Printf("Unrecognized table name: %s", table_name)
    }

    html += `</table>
             <br>
             <a href="/">Back</a>`

    fmt.Fprint(w, html)
}

// execute a query to find bikes with values above a specified threshold,
// and format the results as an HTML table
func findExpensiveBikesHandler(w http.ResponseWriter, r *http.Request) {
    bike_value_floor := r.FormValue("bike_value_floor")

	db, err := sql.Open("sqlite3", db_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    // of course if this were production code we would sanitize the input
	sqlStmt := "SELECT * FROM Bike WHERE Bike.Dollar_value > " + bike_value_floor + ";"

	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

    html := `<style>table { border-collapse: collapse; border-spacing: 0px; } 
                    table, th, td { padding: 5px; border: 1px solid black; }</style>
             <h1>CS 457 HW 5 by Matthew Leeds</h1>
             <h3>Bikes with dollar values greater than ` + bike_value_floor + `</h3>
             <table>
             <tr><th>ID</th><th>Status</th><th>Dollar_value</th><th>Purchase_Time</th><th>Purchaser_Email</th></tr>`

    for rows.Next() {
        var id int
        var status string
        var dollar_value float32
        var purchase_time string
        var purchaser_email string
        err := rows.Scan(&id, &status, &dollar_value, &purchase_time, &purchaser_email)
        if err != nil {
            log.Fatal(err)
        }
        html += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%.2f</td><td>%s</td><td>%s</td></tr>",
                            id, status, dollar_value, purchase_time, purchaser_email)
    }

    html += `</table>
             <br>
             <a href="/">Back</a>`

    fmt.Fprint(w, html)
}

// execute a query to find bikes from the specified state,
// and format the results as an HTML table
func bikesByStateHandler(w http.ResponseWriter, r *http.Request) {
    bike_purchaser_state := r.FormValue("bike_purchaser_state")

	db, err := sql.Open("sqlite3", db_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    // of course if this were production code we would sanitize the input
	sqlStmt := "SELECT Bike.ID, Bike.Dollar_value, Bike.Purchaser_Email, Customer.State FROM Bike JOIN Customer ON Bike.Purchaser_Email = Customer.Email WHERE Customer.State = '" + bike_purchaser_state + "';"

	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

    html := `<style>table { border-collapse: collapse; border-spacing: 0px; } 
                    table, th, td { padding: 5px; border: 1px solid black; }</style>
             <h1>CS 457 HW 5 by Matthew Leeds</h1>
             <h3>Bikes purchased by people from ` + bike_purchaser_state + `</h3>
             <table>
             <tr><th>ID</th><th>Dollar_value</th><th>Purchaser_Email</th><th>State</th></tr>`

    for rows.Next() {
        var id int
        var dollar_value float32
        var purchaser_email string
        var state string
        err := rows.Scan(&id, &dollar_value, &purchaser_email, &state)
        if err != nil {
            log.Fatal(err)
        }
        html += fmt.Sprintf("<tr><td>%d</td><td>%.2f</td><td>%s</td><td>%s</td></tr>",
                            id, dollar_value, purchaser_email, state)
    }

    html += `</table>
             <br>
             <a href="/">Back</a>`

    fmt.Fprint(w, html)
}

// execute an INSERT for the Bike table using the specified values,
// and format the results as an HTML table
func insertBikeHandler(w http.ResponseWriter, r *http.Request) {
    bike_id := r.FormValue("bike_id")
    bike_status := r.FormValue("bike_status")
    bike_dollar_value := r.FormValue("bike_dollar_value")
    bike_purchase_time := r.FormValue("bike_purchase_time")
    bike_purchaser_email := r.FormValue("bike_purchaser_email")

	db, err := sql.Open("sqlite3", db_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    // of course if this were production code we would sanitize the input
	sqlStmt := fmt.Sprintf("INSERT INTO Bike VALUES (%s, '%s', %s, '%s', '%s');",
                           bike_id, bike_status, bike_dollar_value, bike_purchase_time, bike_purchaser_email)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

    html := `<style>table { border-collapse: collapse; border-spacing: 0px; } 
                    table, th, td { padding: 5px; border: 1px solid black; }</style>
             <h1>CS 457 HW 5 by Matthew Leeds</h1>
             <h3>INSERT query successful!</h3>
             <br>
             <a href="/">Back</a>`

    fmt.Fprint(w, html)
}

// execute a SELECT as specified,
// and format the results as an HTML table
func execSelectHandler(w http.ResponseWriter, r *http.Request) {
    select_query := r.FormValue("select_query")

	db, err := sql.Open("sqlite3", db_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    // of course if this were production code we would sanitize the input
	sqlStmt := fmt.Sprintf("SELECT " + select_query + ";")

	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

    html := `<style>table { border-collapse: collapse; border-spacing: 0px; } 
                    table, th, td { padding: 5px; border: 1px solid black; }</style>
             <h1>CS 457 HW 5 by Matthew Leeds</h1>
             <h3>SELECT query results:</h3>
             <table><tr>`

    columns, err := rows.Columns()
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

    for _, column := range columns {
        html += fmt.Sprintf("<th>%s</th>", column)
    }
    html += "</tr>"

    // move the result into a string array
    var result [][]string
    for rows.Next() {
        pointers := make([]interface{}, len(columns))
        container := make([]string, len(columns))
        for i, _ := range pointers {
            pointers[i] = &container[i]
        }
        rows.Scan(pointers...)
        result = append(result, container)
    }

    for i := range result {
        html += "<tr>"
        for j := range columns {
            html += fmt.Sprintf("<td>%s</td>", result[i][j])
        }
        html += "</tr>"
    }

    html += `</table>
             <br>
             <a href="/">Back</a>`

    fmt.Fprint(w, html)
}
