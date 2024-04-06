package controllers

import (
	"api/src/model"
	"fmt"
	"sync"
)

// Database connection (replace with your actual connection logic)
// var db *sql.DB
// db, err :=

func Create_ad(ad model.Ad) error {
	var err error
	//Insert ad
	tx := Db.Begin() // Begin a transaction

	defer func() {
		if err != nil { // Defer closing the transaction until the end
			fmt.Println("rollback")
			tx.Rollback() // Rollback if any error occurs
		} else {
			fmt.Println("success")

			tx.Commit() // Commit if everything succeeds
		}
	}()
	if ad.Conditions[0].AgeEnd == 0 || ad.Conditions[0].AgeStart == 0 {
		err := Db.Exec("INSERT INTO Ads (title, start_at, end_at) VALUES (?, ?, ?)", ad.Title, ad.StartAt, ad.EndAt).Error
		if err != nil {
			return err
		}
	} else {
		fmt.Println(ad)
		err := Db.Exec("INSERT INTO Ads (title, start_at, end_at ,min_age,max_age) VALUES (?, ?, ?,?,?)",
			ad.Title, ad.StartAt, ad.EndAt, ad.Conditions[0].AgeStart, ad.Conditions[0].AgeEnd).Error
		if err != nil {
			return err
		}

	}
	var adID int
	err = Db.Raw("SELECT @@IDENTITY AS adID;").Scan(&adID).Error

	if err != nil {
		return err
	}
	err = Create_condition(ad.Conditions[0], adID)
	if err != nil {
		return err
	}
	// fmt.Println("add_ad")
	return nil

}

func get_clause(data []string) string {

	Clause := ""
	for _, Code := range data {
		Clause += "'" + Code + "',"
	}
	Clause = Clause[:len(Clause)-1] // Remove trailing comma
	return Clause
}
func Create_condition(condition model.Condition, adId int) error {

	var (
		err                                         error
		countryClause, platformClause, genderClause string
		wg                                          sync.WaitGroup
	)

	wg.Add(3)

	// Execute get_clause concurrently for faster data fetching
	go func() {
		defer wg.Done()
		countryClause = get_clause(condition.Country)
	}()

	go func() {
		defer wg.Done()
		platformClause = get_clause(condition.Platform)
	}()

	go func() {
		defer wg.Done()
		genderClause = get_clause(condition.Gender)
	}()

	// wg.Wait() // Wait for all goroutines to finish fetching data

	// platformClause = get_clause(condition.Platform)
	// genderClause = get_clause(condition.Gender)
	// countryClause = get_clause(condition.Country)

	sql_query := `INSERT INTO AdConditions (ad_id, condition_id)
		SELECT ?, c.id
	FROM Conditions c
	WHERE   (
	    c.type = 'gender' AND c.value IN (` + genderClause + `)
	) OR (
	    c.type = 'country' AND c.value IN (` + countryClause + `)
	) OR (
	    c.type = 'platform' AND c.value IN (` + platformClause + `)
	);`

	fmt.Println("Create condition")
	// fmt.Println(sql_query)
	err = Db.Exec(sql_query, adId).Error
	if err != nil {
		return err
	}

	return nil
}
func Find_ad(condition model.Search_Condition) ([]model.Result, error) {
	// now := time.Now() // Get current time
	num_param := 0
	params := []interface{}{} // Initial parameter

	// ... (rest of the code for building parameters similar to the previous response)

	whereClause := ("WHERE a.start_at <= NOW() AND a.end_at >= NOW() ") // Initial WHERE clause
	typeClause := "AND ("

	if condition.Age != 0 {

		whereClause += "AND (a.min_age <= ? AND a.max_age >= ?) "
		// str := fmt.Sprintf("%d", condition.Age)
		params = append(params, condition.Age, condition.Age)

	}

	// Add conditions based on Search_Condition
	if len(condition.Gender) > 0 {
		num_param++

		typeClause += "(c.type = 'gender' AND c.value = (?)) "
		params = append(params, condition.Gender)
	} //else {
	// 	num_param++
	// 	whereClause += "OR (c.type = 'gender' AND c.value IN (?, ?)) "
	// 	params = append(params, "M", "F")
	// }
	if len(condition.Country) > 0 {
		num_param++
		// print("how much", len(condition.Country)*2-1)
		// countryPlaceholders := strings.Repeat("?,", len(condition.Country))[:len(condition.Country)*2-1]
		// whereClause += fmt.Sprintf("OR (c.type = 'country' AND c.value IN (%s)) ", countryPlaceholders)
		typeClause += "OR (c.type = 'country' AND c.value = (?)) "
		params = append(params, condition.Country)

		// for _, country := range condition.Country {
		// 	params = append(params, country)
		// }
	}
	if len(condition.Platform) != 0 {
		num_param++

		// platformPlaceholders := strings.Repeat("?,", len(condition.Platform))[:len(condition.Platform)*2-1]
		// fmt.Println(platformPlaceholders)
		// whereClause += fmt.Sprintf("OR (c.type = 'platform' AND c.value IN (%s)) ", platformPlaceholders)
		typeClause += "OR (c.type = 'platform' AND c.value = (?)) "

		params = append(params, condition.Platform)

		// for _, platform := range condition.Platform {
		// 	params = append(params, platform)
		// }
	}

	var sql_query string
	sql_query = "SELECT a.title,a.end_at" +
		" FROM Ads a " +
		" INNER JOIN AdConditions ac ON a.id = ac.ad_id " +
		" INNER JOIN Conditions c ON ac.condition_id = c.id "
	if num_param != 0 {
		typeClause += ")"
		// fmt.Println("type \n", typeClause)
		whereClause += typeClause
		sql_query += whereClause + "GROUP BY ad_id " + "HAVING COUNT(*) = ?;" // Adjust for number of conditions

		params = append(params, num_param)

	} else {

		sql_query += whereClause
	}

	// sql_query += "order by end_at asc"
	// sql_query = `

	// select * from Ads;
	// `
	// fmt.Println(sql_query)
	// fmt.Println(params...)

	var ads []model.Result
	// err := Db.Exec("SELECT a.title,a.end_at" +
	// 	" FROM Ads a ").Scan(&ads).Error

	err := Db.Raw(sql_query, params...).Scan(&ads).Error
	if err != nil {

		return nil, err
	}
	return ads, err

}
