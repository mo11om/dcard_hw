package controllers

import (
	"api/src/model"
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

var redisClient *Client

func Init_redis() error {
	cfg := &Config{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	var err error
	redisClient, err = NewClient(cfg)
	if err != nil {
		return err // Handle error more gracefully in production
	}
	fmt.Println("redis conn success")
	return nil
}

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
			flushAllAdsCache()
		}
	}()
	if ad.Conditions[0].AgeEnd == 0 || ad.Conditions[0].AgeStart == 0 {
		err := Db.Exec("INSERT INTO Ads (title, start_at, end_at) VALUES (?, ?, ?)", ad.Title, ad.StartAt, ad.EndAt).Error
		if err != nil {
			return err
		}
	} else {
		// fmt.Println(ad)
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
func flushAllAdsCache() error {

	// Implement logic to connect to Redis and flush the cache (e.g., using FlushAll)
	// Consider alternative cache invalidation strategies for better performance

	return redisClient.FlushAll(context.Background()) // Consider setting a suitable expiration time // Replace with actual implementation
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

	wg.Wait() // Wait for all goroutines to finish fetching data

	// platformClause = get_clause(condition.Platform)
	// genderClause = get_clause(condition.Gender)
	// countryClause = get_clause(condition.Country)
	// fmt.Println(genderClause, countryClause, platformClause)
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

func find_adBysql(condition model.Search_Condition) ([]model.Result, error) {

	num_param := 0
	params := []interface{}{} // Initial parameter

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
		sql_query += whereClause + `GROUP BY ad_id 
		HAVING COUNT(*) = ?` // Adjust for number of conditions

		params = append(params, num_param)

	} else {

		sql_query += whereClause
	}
	//order by asc
	sql_query += ` ORDER BY a.end_at asc`

	// end of sql_query
	sql_query += ";"
	// sql_query = `

	// select * from Ads;
	// `
	// fmt.Println(sql_query)
	// fmt.Println(params...)

	var ads []model.Result

	err := Db.Raw(sql_query, params...).Scan(&ads).Error
	if err != nil {

		return nil, err
	}
	err = redisSetAd(getSearchString(condition), ads)
	if err != nil {
		fmt.Println("cache fail")
	}
	return ads, err
}

func redisSetAd(key string, ads []model.Result) error {
	data, err := json.Marshal(ads)

	if err != nil {
		return err
	}
	return redisClient.Set(context.Background(), key, data, 0) // Consider setting a suitable expiration time
}
func redisGetAd(key string) ([]model.Result, error) {
	result, err := redisClient.Get(context.Background(), key)
	// fmt.Println("get error", err)
	if result == "" {
		// fmt.Println("leave nil")
		return nil, err // Key not found in cache

	} else if err != nil {
		// fmt.Println("leave != nil")

		return nil, err
	}
	// fmt.Println(" success get res")
	var ads []model.Result
	if err := json.Unmarshal([]byte(result), &ads); err != nil {
		return nil, err
	}
	// fmt.Println(key, "get ", ads)

	return ads, nil
}

// /
func getSearchString(condition model.Search_Condition) string {
	return fmt.Sprintf("%d_%s_%s_%s", condition.Age, condition.Gender, condition.Country, condition.Platform)
}

func find_adByredis(condition model.Search_Condition) ([]model.Result, error) {
	key := getSearchString(condition)
	ads, err := redisGetAd(key)
	if err != nil {
		return nil, err
	}
	return ads, nil
	// redisClient.S

}

// Package  ads

// FunctionName explains the function's purpose and behavior.
//
// Args:
//   * param1 model.Search_Condition - input  condition
//   * param2 (type) - Description of the second parameter with type information.
//
// Returns:
//
//   * []model.Result
//   * error
//
//

func Find_ad(condition model.Search_Condition) ([]model.Result, error) {

	var (
		err error
		ads []model.Result
	)
	ads, err = find_adByredis(condition)
	if err != nil {

		return nil, err
	}
	if ads == nil {
		fmt.Println("redis fail")
		ads, err = find_adBysql(condition)
		if err != nil {

			return nil, err
		}
	} else {
		fmt.Println("redis success")
	}
	// fmt.Println(ads)
	return getSlice(ads, condition.Limit, condition.Offset), err

}
func getSlice(data []model.Result, limit, offset int) []model.Result {
	if limit == 0 {
		return nil // Return empty slice for limit 0
	}
	start := offset
	end := min(start+limit, len(data))
	// fmt.Println(start, end, data[start:end])
	return data[start:end]
}
