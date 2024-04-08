package controllers

import (
	"api/src/model"
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

var redisClient *Client

// init redis client
func Init_redis() error {
	cfg := &Config{
		Addr:     "redis:6379",
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

// redis operation
// flush all ad record in redis
func flushAllAdsCache() error {

	// Implement logic to connect to Redis and flush the cache (e.g., using FlushAll)
	// Consider alternative cache invalidation strategies for better performance

	return redisClient.FlushAll(context.Background()) // Consider setting a suitable expiration time // Replace with actual implementation
}

// set model.result with key  "Age_Gender_Country_Platform"
func redisSetAd(key string, ads []model.Result) error {
	data, err := json.Marshal(ads)

	if err != nil {
		return err
	}
	return redisClient.Set(context.Background(), key, data, 0) // Consider setting a suitable expiration time
}

// get model.result  with key  "Age_Gender_Country_Platform"

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

//
// clause
//

// get_clause constructs a comma-separated string of single-quoted values from a slice of strings.
// It's typically used to generate SQL IN clauses for conditions.
//
// It takes a slice of strings as input and returns a single string in the format:
// 'value1', 'value2', ..., 'valueN'
func get_clause(data []string) string {

	Clause := ""
	for _, Code := range data {
		Clause += "'" + Code + "',"
	}
	Clause = Clause[:len(Clause)-1] // Remove trailing comma
	return Clause
}

// getSearchString generates a search string by combining condition fields.
// It takes a Search_Condition struct and returns a string suitable for search purposes.
//
// The generated string format is: "<age>_<gender>_<country>_<platform>".
// - `<age>`: The value of the `Age` field in the `Search_Condition` struct.
// - `<gender>`: The value of the `Gender` field in the `Search_Condition` struct.
// - `<country>`: The value of the `Country` field in the `Search_Condition` struct.
// - `<platform>`: The value of the `Platform` field in the `Search_Condition` struct.
//
// This function is likely used to create a unique identifier for search combinations
// or for building search queries based on the provided conditions.
func getSearchString(condition model.Search_Condition) string {
	return fmt.Sprintf("%d_%s_%s_%s", condition.Age, condition.Gender, condition.Country, condition.Platform)
}

//
// database dealing
//

//create ads

// CreateAd creates a new ad record in the database and its associated conditions.
// It takes an Ad struct containing the ad details and returns an error if any occur during creation.
//
//	in example/insert sql
func Create_ad(ad model.Ad) error {
	var err error
	//Insert ad
	// Begin a transaction to ensure data consistency
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
	// Insert the ad record with or without age fields
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
	// Retrieve the newly created ad ID
	var adID int
	err = Db.Raw("SELECT @@IDENTITY AS adID;").Scan(&adID).Error

	if err != nil {
		return err
	}
	// Create the ad's conditions using a separate function
	err = Create_condition(ad.Conditions[0], adID)
	if err != nil {
		return err
	}
	// fmt.Println("add_ad")
	return nil

}

// CreateAdConditions creates associations between an ad and its target conditions in the database.
// It takes a Condition struct and the ad ID and returns an error if any occur during creation.
func Create_condition(condition model.Condition, adId int) error {
	// Concurrently fetch clauses for faster data retrieval (if applicable
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

	// Construct the SQL query dynamically
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

	// fmt.Println("Create condition")
	// Execute the query to create ad condition records
	err = Db.Exec(sql_query, adId).Error
	if err != nil {
		return err
	}

	return nil
}

//
//find ad
//

// eample of search sql
// ..
//
// find_adBysql retrieves ads matching given search conditions from the database.
// It dynamically constructs a SQL query based on the provided criteria,
// handling various combinations of age, gender, country, and platform conditions.
//
// It also caches the retrieved results in Redis for potential future performance gains.
//
// Args:
//
//	condition (model.Search_Condition): The search criteria to match against ads.
//
// Returns:
//
//	[]model.Result (nil on error): An array of matching ad results.
//	error: Any error encountered during database retrieval or caching.
func find_adBysql(condition model.Search_Condition) ([]model.Result, error) {
	// Initialize variables for dynamic query construction
	num_param := 0
	params := []interface{}{} //Parameters for query execution

	whereClause := ("WHERE a.start_at <= NOW() AND a.end_at >= NOW() ") // Base WHERE clause
	typeClause := "AND ("                                               // Clause for handling multiple conditions
	// Add age condition if specified
	if condition.Age != 0 {

		whereClause += "AND (a.min_age <= ? AND a.max_age >= ?) "

		params = append(params, condition.Age, condition.Age)

	}

	// Add gender condition if specified
	if len(condition.Gender) > 0 {
		num_param++

		typeClause += "(c.type = 'gender' AND c.value = (?)) "
		params = append(params, condition.Gender)
	}
	// Add country condition if specified, handling "ALL" value
	if len(condition.Country) > 0 {
		num_param++
		typeClause += "OR (c.type = 'country' AND (c.value = (? ) or c.value ='ALL')) "
		params = append(params, condition.Country)

	}

	// Add platform condition if specified
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
	// Construct the final SQL query
	var sql_query string
	sql_query = "SELECT a.title,a.end_at" +
		" FROM Ads a " +
		" INNER JOIN AdConditions ac ON a.id = ac.ad_id " +
		" INNER JOIN Conditions c ON ac.condition_id = c.id "

		// Complete conditional clauses and grouping if necessary
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
	// Add ordering and finalization
	sql_query += ` ORDER BY a.end_at asc`
	sql_query += ";"

	// Execute the query, handle errors, and cache results
	var ads []model.Result
	err := Db.Raw(sql_query, params...).Scan(&ads).Error
	if err != nil {

		return nil, err
	}
	// Cache results in Redis
	err = redisSetAd(getSearchString(condition), ads)
	if err != nil {
		fmt.Println("cache fail")
	}
	return ads, err
}

// find_adByredis attempts to retrieve ads matching given search conditions from Redis.
// It uses a uniquely generated key based on the condition to fetch data from Redis.
//
// Args:
//
//	condition (model.Search_Condition): The search criteria to match against ads.
//
// Returns:
//
//	[]model.Result (nil on error): An array of matching ad results if found in Redis.
//	error: Any error encountered during Redis retrieval.
func find_adByredis(condition model.Search_Condition) ([]model.Result, error) {
	key := getSearchString(condition)
	ads, err := redisGetAd(key)
	if err != nil {
		return nil, err
	}
	return ads, nil
	// redisClient.S

}

// Find_ad searches for ads matching the given search conditions.
// It prioritizes retrieving results from Redis for faster response times,
// falling back to a database query if Redis is unavailable or the data is missing.
//
// Args:
//
//	condition (model.Search_Condition): The search criteria to match against ads.
//
// Returns:
//
//	[]model.Result (nil on error): An array of matching ad results.
//	error: Any error encountered during the search process.
func Find_ad(condition model.Search_Condition) ([]model.Result, error) {

	var (
		err error
		ads []model.Result
	)
	// Attempt to fetch ads from Redis for faster retrieval
	ads, err = find_adByredis(condition)
	if err != nil {

		return nil, err
	}
	// If Redis fails or data is missing, fallback to database query
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

// getSlice extracts a slice from a given slice of model.Result objects,
// applying limit and offset for pagination purposes.
//
// Args:
//
//	data ([]model.Result): The original slice of data to be sliced.
//	limit (int): The maximum number of elements to include in the returned slice.
//	offset (int): The starting index for the slice, indicating where to begin extracting elements.
//
// Returns:
//
//	[]model.Result: The sliced portion of the data, adhering to the given limit and offset.
//	nil: If the limit is specified as 0, an empty slice is returned.
func getSlice(data []model.Result, limit, offset int) []model.Result {
	if limit == 0 {
		return nil // Return empty slice for limit 0
	}
	start := offset
	end := min(start+limit, len(data))
	// fmt.Println(start, end, data[start:end])
	return data[start:end]
}
