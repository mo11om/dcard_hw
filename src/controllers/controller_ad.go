package controllers

import (
	"api/src/model"
	"fmt"
)

// Database connection (replace with your actual connection logic)
// var db *sql.DB
// db, err :=

func Create_ad(ad model.Ad) error {

	// Insert ad
	// if ad.Conditions.AgeEnd == 0 || ad.Conditions.AgeStart == 0 {
	// 	err := Db.Raw("INSERT INTO Ads (title, start_at, end_at) VALUES (?, ?, ?)", ad.Title, ad.StartAt, ad.EndAt).Error
	// 	if err != nil {
	// 		return err
	// 	}
	// } else {
	// 	err := Db.Raw("INSERT INTO Ads (title, start_at, end_at ,min_age,max_age) VALUES (?, ?, ?)",
	// 	 ad.Title, ad.StartAt, ad.EndAt, ad.Conditions.AgeStart, ad.Conditions.AgeEnd).Error
	// 	if err != nil {
	// 		return err
	// 	}

	// }

	// adID, err := Db.LastInsertId()
	// if err != nil {
	// 	return 0, err
	// }

	// countryClause := ""
	// for _, countryCode := range countryCodes {
	// 	countryClause += "'" + countryCode + "',"
	// }
	// countryClause = countryClause[:len(countryClause)-1] // Remove trailing comma

	// platformClause := ""
	// for _, platformName := range platformNames {
	// 	platformClause += "'" + platformName + "',"
	// }
	// platformClause = platformClause[:len(platformClause)-1] // Remove trailing comma

	// sql_query := `INSERT INTO AdConditions (ad_id, condition_id)
	// 	SELECT ?, c.id
	// FROM Conditions c
	// WHERE   (
	//     c.type = 'gender' AND c.value IN ('M', 'F')
	// ) OR (
	//     c.type = 'country' AND c.value IN (` + countryClause + `)
	// ) OR (
	//     c.type = 'platform' AND c.value IN (` + platformClause + `)
	// );`

	fmt.Println("add_ad")
	return nil

}
func Create_condition() {
	fmt.Println("Create condition")
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
	// sql_query = `

	// select * from Ads;
	// `
	fmt.Println(sql_query)
	fmt.Println(params...)

	var ads []model.Result
	// err := Db.Raw("SELECT a.title,a.end_at" +
	// 	" FROM Ads a ").Scan(&ads).Error

	err := Db.Raw(sql_query, params...).Scan(&ads).Error
	if err != nil {

		return nil, err
	}
	return ads, err

}
