package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	db "pluto_remastered/config"
	"pluto_remastered/helpers"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

func executeWithoutResult(query string, wg *sync.WaitGroup) {
	defer wg.Done()

	db.DB.Exec(query)
}

func TestGenerateUserId(c *fiber.Ctx) error {

	queries := []string{
		`UPDATE penjualan SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
					COALESCE(s.user_id, mm.user_id, tl.user_id) as user_id,
					CASE WHEN k.teamleader_id IS NULL THEN NULL ELSE COALESCE(tl.user_id, k.teamleader_id) END as user_id_subtitute
			FROM penjualan k
			LEFT JOIN salesman s
				ON k.salesman_id = s.id
			LEFT JOIN md.merchandiser mm
				ON k.merchandiser_id = mm.id
			LEFT JOIN teamleader tl
				ON k.teamleader_id = tl.id
			WHERE k.user_id IS NULL AND k.salesman_id = 781
			) data
			WHERE penjualan.id = data.id`,

		`UPDATE kunjungan SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
					COALESCE(s.user_id, mm.user_id, tl.user_id) as user_id,
					CASE WHEN k.teamleader_id IS NULL THEN NULL ELSE COALESCE(tl.user_id, k.teamleader_id) END as user_id_subtitute
			FROM kunjungan k
			LEFT JOIN salesman s
				ON k.salesman_id = s.id
			LEFT JOIN md.merchandiser mm
				ON k.merchandiser_id = mm.id
			LEFT JOIN teamleader tl
				ON k.teamleader_id = tl.id
			WHERE k.user_id IS NULL AND k.salesman_id = 781
			) data
			WHERE kunjungan.id = data.id`,
	}

	var wg sync.WaitGroup

	// Launch concurrent Goroutines
	for _, query := range queries {
		wg.Add(1)
		go executeWithoutResult(query, &wg)
	}

	// Wait for all Goroutines to finish
	wg.Wait()

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		Message: "Success",
		Success: true,
	})
}

func GenerateTransactionsUserId(c *fiber.Ctx) error {

	queries := []string{
		`UPDATE penjualan SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
							COALESCE(s.user_id, mm.user_id, tl.user_id) as user_id,
							CASE WHEN k.teamleader_id IS NULL THEN NULL ELSE COALESCE(tl.user_id, k.teamleader_id) END as user_id_subtitute
			FROM penjualan k
			LEFT JOIN salesman s
				ON k.salesman_id = s.id
			LEFT JOIN md.merchandiser mm
				ON k.merchandiser_id = mm.id
			LEFT JOIN teamleader tl
				ON k.teamleader_id = tl.id
			WHERE k.user_id IS NULL
			) data
			WHERE penjualan.id = data.id`,

		`UPDATE kunjungan SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
							COALESCE(s.user_id, mm.user_id, tl.user_id) as user_id,
							CASE WHEN k.teamleader_id IS NULL THEN NULL ELSE COALESCE(tl.user_id, k.teamleader_id) END as user_id_subtitute
			FROM kunjungan k
			LEFT JOIN salesman s
				ON k.salesman_id = s.id
			LEFT JOIN md.merchandiser mm
				ON k.merchandiser_id = mm.id
			LEFT JOIN teamleader tl
				ON k.teamleader_id = tl.id
			WHERE k.user_id IS NULL
			) data
			WHERE kunjungan.id = data.id`,

		`UPDATE payment SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
							COALESCE(s.user_id, mm.user_id, tl.user_id) as user_id,
							CASE WHEN k.teamleader_id IS NULL THEN NULL ELSE COALESCE(tl.user_id, k.teamleader_id) END as user_id_subtitute
			FROM payment k
			LEFT JOIN salesman s
				ON k.salesman_id = s.id
			LEFT JOIN md.merchandiser mm
				ON k.merchandiser_id = mm.id
			LEFT JOIN teamleader tl
				ON k.teamleader_id = tl.id
			WHERE k.user_id IS NULL
			) data
			WHERE payment.id = data.id`,

		`UPDATE piutang SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
							tl.user_id as user_id,
							CASE WHEN k.teamleader_id IS NULL THEN NULL ELSE COALESCE(tl.user_id, k.teamleader_id) END as user_id_subtitute
			FROM piutang k
			LEFT JOIN teamleader tl
				ON k.teamleader_id = tl.id
			WHERE k.user_id IS NULL
			) data
			WHERE piutang.id = data.id`,

		`UPDATE pengembalian SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
							COALESCE(s.user_id, mm.user_id, tl.user_id) as user_id,
							CASE WHEN k.teamleader_id IS NULL THEN NULL ELSE COALESCE(tl.user_id, k.teamleader_id) END as user_id_subtitute
			FROM pengembalian k
			LEFT JOIN salesman s
				ON k.salesman_id = s.id
			LEFT JOIN md.merchandiser mm
				ON k.merchandiser_id = mm.id
			LEFT JOIN teamleader tl
				ON k.teamleader_id = tl.id
			WHERE k.user_id IS NULL
			) data
			WHERE pengembalian.id = data.id`,

		`UPDATE pembayaran_piutang SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
							COALESCE(s.user_id, mm.user_id, tl.user_id) as user_id,
							CASE WHEN k.teamleader_id IS NULL THEN NULL ELSE COALESCE(tl.user_id, k.teamleader_id) END as user_id_subtitute
			FROM pembayaran_piutang k
			LEFT JOIN salesman s
				ON k.salesman_id = s.id
			LEFT JOIN md.merchandiser mm
				ON k.merchandiser_id = mm.id
			LEFT JOIN teamleader tl
				ON k.teamleader_id = tl.id
			WHERE k.user_id IS NULL
			) data
			WHERE pembayaran_piutang.id = data.id`,

		`UPDATE stok_salesman_riwayat SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
							COALESCE(s.user_id, mm.user_id, tl.user_id) as user_id,
							CASE WHEN k.teamleader_id IS NULL THEN NULL ELSE COALESCE(tl.user_id, k.teamleader_id) END as user_id_subtitute
			FROM stok_salesman_riwayat k
			LEFT JOIN salesman s
				ON k.salesman_id = s.id
			LEFT JOIN md.merchandiser mm
				ON k.merchandiser_id = mm.id
			LEFT JOIN teamleader tl
				ON k.teamleader_id = tl.id
			WHERE k.user_id IS NULL
			) data
			WHERE stok_salesman_riwayat.id = data.id`,

		`UPDATE stok_salesman SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
							COALESCE(s.user_id, mm.user_id) as user_id,
							NULL as user_id_subtitute
			FROM stok_salesman k
			LEFT JOIN salesman s
				ON k.salesman_id = s.id
			LEFT JOIN md.merchandiser mm
				ON k.merchandiser_id = mm.id
			WHERE k.user_id IS NULL 
			) data
			WHERE stok_salesman.id = data.id`,

		`UPDATE md.stok_merchandiser_riwayat SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
							mm.user_id as user_id,
							NULL as user_id_subtitute
			FROM md.stok_merchandiser_riwayat k
			LEFT JOIN md.merchandiser mm
				ON k.merchandiser_id = mm.id
			WHERE k.user_id IS NULL 
			) data
			WHERE md.stok_merchandiser_riwayat.id = data.id`,

		`UPDATE md.stok_merchandiser SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
							mm.user_id as user_id,
							NULL as user_id_subtitute
			FROM md.stok_merchandiser k
			LEFT JOIN md.merchandiser mm
				ON k.merchandiser_id = mm.id
			WHERE k.user_id IS NULL 
			) data
			WHERE md.stok_merchandiser.id = data.id`,

		`UPDATE kunjungan_log SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
							COALESCE(s.user_id, mm.user_id, tl.user_id) as user_id,
							CASE WHEN k.teamleader_id IS NULL THEN NULL ELSE COALESCE(tl.user_id, k.teamleader_id) END as user_id_subtitute
			FROM kunjungan_log k
			LEFT JOIN salesman s
				ON k.salesman_id = s.id
			LEFT JOIN md.merchandiser mm
				ON k.merchandiser_id = mm.id
			LEFT JOIN teamleader tl
				ON k.teamleader_id = tl.id
			WHERE k.user_id IS NULL
			) data
			WHERE kunjungan_log.id = data.id`,

		`UPDATE md.transaction SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
			FROM (
			SELECT k.id, 
							COALESCE(s.user_id, mm.user_id) as user_id,
							NULL as user_id_subtitute
			FROM md.transaction k
			LEFT JOIN salesman s
				ON k.salesman_id = s.id
			LEFT JOIN md.merchandiser mm
				ON k.merchandiser_id = mm.id
			WHERE k.user_id IS NULL
			) data
			WHERE md.transaction.id = data.id`,
	}
	var wg sync.WaitGroup

	// Launch concurrent Goroutines
	for _, query := range queries {
		wg.Add(1)
		go executeWithoutResult(query, &wg)
	}

	// Wait for all Goroutines to finish
	wg.Wait()

	// for _, query := range queries {
	// 	fmt.Println(query) // Prints with proper formatting
	// 	fmt.Println("------")
	// }

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		Message: "Data has been loaded successfully",
		Success: true,
	})

}

func GenerateUserLog(c *fiber.Ctx) error {
	db.DB.Exec(`INSERT INTO user_log_branch (user_id, user_id_subtitute, branch_id, start_date, end_date, last_visit_date)
				SELECT sq.user_id,
						COALESCE(sq.user_id_subtitute, -1),
						sq.branch_id,
						sq.min_date,
						DATE(LEAD(sq.min_date, 1) OVER (PARTITION BY sq.user_id ORDER BY sq.user_id, sq.min_date)::date - INTERVAL '1 day'),
						CASE WHEN 
							DATE(LEAD(sq.min_date, 1) OVER (PARTITION BY sq.user_id ORDER BY sq.user_id, sq.min_date)::date - INTERVAL '1 day') IS NOT NULL 
								THEN DATE(LEAD(sq.min_date, 1) OVER (PARTITION BY sq.user_id ORDER BY sq.user_id, sq.min_date)::date - INTERVAL '1 day') 
								ELSE sq.max_date 
						END
				FROM (
					WITH last_visit as (SELECT user_id, MAX(DATE(tanggal_kunjungan)) as tgl_kunjungan FROM kunjungan GROUP BY user_id)
				
					SELECT k.branch_id, 
									k.salesman_id, 
									k.merchandiser_id, 
									k.teamleader_id, 
									k.user_id, 
									k.user_id_subtitute, 
									MIN(DATE(k.tanggal_kunjungan)) as min_date,
									MAX(lv.tgl_kunjungan) as max_date
					FROM kunjungan k
					LEFT JOIN last_visit lv
						ON k.user_id = lv.user_id
					GROUP BY k.branch_id, k.salesman_id, k.merchandiser_id, k.teamleader_id, k.user_id, k.user_id_subtitute
				) sq
				WHERE sq.user_id IS NOT NULL AND sq.user_id > 0 AND sq.branch_id IS NOT NULL
				ORDER BY sq.user_id, sq.min_date
				ON CONFLICT (user_id, branch_id, start_date, user_id_subtitute, last_visit_date)
				DO NOTHING`)

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		Message: "Cron success",
		Success: true,
	})
}

func GenerateFlag(c *fiber.Ctx) error {

	start := time.Now()
	type GetFlagRequest struct {
		FlagTable *string `json:"flagTable"`
	}
	var flagReq GetFlagRequest
	if err := c.QueryParser(&flagReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	dataSend, err := json.Marshal(flagReq)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	responseData, err := helpers.SendCurl(dataSend, "GET", "https://rest.pt-bks.com/olympus/flagGenerateQuery")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	queries := []string{}
	for _, val := range responseData["data"].([]interface{}) {
		queries = append(queries, val.(string))
	}

	var wg sync.WaitGroup

	for _, query := range queries {
		wg.Add(1)
		// go executeGORMQuery(query, resultsChan, i, &wg)
		go helpers.ExecuteGORMQueryWithoutResult(query, &wg)
	}

	// Wait for all Goroutines to finish
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Function took %s", elapsed)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"success": true,
		"elapsed": elapsed,
	})
}

func GetData(c *fiber.Ctx) error {

	start := time.Now()
	type GetFlagRequest struct {
		FlagID   string `json:"flagId"`
		UserID   int    `json:"userId"`
		BranchId int    `json:"branchId"`
	}
	var flagReq GetFlagRequest
	if err := c.QueryParser(&flagReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	dataSend, err := json.Marshal(flagReq.FlagID)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	urls := "https://rest.pt-bks.com/pluto-mobile/getData2?flagId=" + flagReq.FlagID

	responseData, err := helpers.SendCurl(dataSend, "GET", urls)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	queries := []string{}
	keyQuery := []string{}

	// fmt.Println("=============")
	// fmt.Println(responseData)
	tempUserID := strconv.Itoa(flagReq.UserID)
	tempBranchId := strconv.Itoa(flagReq.BranchId)
	for key, val := range responseData["data"].(map[string]interface{}) {
		r := strings.NewReplacer("userId", tempUserID, "branchId", tempBranchId)
		tempString := r.Replace(val.(string))
		queries = append(queries, tempString)

		keyQuery = append(keyQuery, key)
	}

	var wg sync.WaitGroup
	resultsChan := make(chan map[string][]map[string]interface{}, len(queries))
	// tempResults := make([][]map[string]interface{}, len(queries))
	tempResults := make([]map[string]interface{}, 0, len(queries))

	for i, query := range queries {
		wg.Add(1)
		// go executeGORMQuery(query, resultsChan, i, &wg)
		go helpers.ExecuteGORMQueryIndexString(query, resultsChan, keyQuery[i], &wg)
	}

	// Wait for all Goroutines to finish
	wg.Wait()
	close(resultsChan)

	for result := range resultsChan {
		for key, res := range result {
			tempResults = append(tempResults, map[string]interface{}{
				key: res,
			})
		}
	}
	finalResult := make(map[string]interface{})
	for _, val := range tempResults {
		for key, res := range val {
			// fmt.Println(key) //nama table
			// fmt.Println(res) //data each table
			finalResult[key] = res
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Function took %s", elapsed)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    finalResult,
		"keyData": keyQuery,
		"elapsed": elapsed,
	})
}

func GetDataToday(c *fiber.Ctx) error {
	start := time.Now()
	type GetFlagRequest struct {
		FlagTable *string `json:"flagTable"`
		UserID    int     `json:"userId"`
		BranchId  int     `json:"branchId"`
	}
	var flagReq GetFlagRequest
	if err := c.QueryParser(&flagReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	dataSend, err := json.Marshal(flagReq)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	responseData, err := helpers.SendCurl(dataSend, "GET", "https://rest.pt-bks.com/pluto-mobile/getDataTodayGenerate")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	queries := []string{}
	keyQuery := []string{}

	// fmt.Println("=============")
	// fmt.Println(responseData)
	tempUserID := strconv.Itoa(flagReq.UserID)
	tempBranchId := strconv.Itoa(flagReq.BranchId)
	for key, val := range responseData["data"].(map[string]interface{}) {
		r := strings.NewReplacer("userId", tempUserID, "branchId", tempBranchId)
		tempString := r.Replace(val.(string))
		queries = append(queries, tempString)

		keyQuery = append(keyQuery, key)
	}

	var wg sync.WaitGroup
	resultsChan := make(chan map[string][]map[string]interface{}, len(queries))
	// tempResults := make([][]map[string]interface{}, len(queries))
	tempResults := make([]map[string]interface{}, 0, len(queries))

	for i, query := range queries {
		wg.Add(1)
		// go executeGORMQuery(query, resultsChan, i, &wg)
		go helpers.ExecuteGORMQueryIndexString(query, resultsChan, keyQuery[i], &wg)
	}

	// Wait for all Goroutines to finish
	wg.Wait()
	close(resultsChan)

	for result := range resultsChan {
		for key, res := range result {
			tempResults = append(tempResults, map[string]interface{}{
				key: res,
			})
		}
	}
	finalResult := make(map[string]interface{})
	for _, val := range tempResults {
		for key, res := range val {
			// fmt.Println(key) //nama table
			// fmt.Println(res) //data each table
			finalResult[key] = res
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Function took %s", elapsed)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    finalResult,
		"keyData": keyQuery,
		"elapsed": elapsed,
	})
}
