package main

import (
	"context"
	"fmt"
)

func insertCustomer(customer Customer) (int64, error) {

	result, err := conn.Exec(context.Background(), "INSERT INTO public.customers (uuid, first_name, last_name, address, phone_number, email, employer, "+
		"annual_income, request_credit_amount, additional_information) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);", customer.UUID, customer.FirstName,
		customer.LastName, customer.Address, customer.PhoneNumber, customer.Email, customer.Employer, customer.AnnualIncome, customer.RequestedCreditAmount, customer.AdditionalInformation)
	if err != nil {
		return 0, fmt.Errorf("unable to execute the query. %v", err)
	}

	return result.RowsAffected(), nil
}

// func fetchAllCustomers(conn *pgxpool.Pool) {
// 	var customers []Customer
// 	rows, err := conn.Query(context.Background(), "SELECT * FROM customers")
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
// 		os.Exit(1)
// 	}

// 	for rows.Next() {
// 		var customer Customer
// 		err := rows.Scan(&customer.ID, &customer.UUID, &customer.FirstName, &customer.LastName, &customer.Address,
// 			&customer.PhoneNumber, &customer.Email, &customer.Employer, &customer.AnnualIncome, &customer.RequestedCreditAmount, &customer.AdditionalInformation)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		customers = append(customers, customer)
// 	}

// 	fmt.Println(customers)
// }
