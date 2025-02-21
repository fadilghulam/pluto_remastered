package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	db "pluto_remastered/config"
	"pluto_remastered/helpers"
	"pluto_remastered/structs"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetListOrder(c *fiber.Ctx) error {
	type TemplateInputUser struct {
		UserId   *string `json:"userId"`
		ChildId  *string `json:"childId"`
		ParentId *string `json:"parentId"`
		Type     *string `json:"type"`
	}

	inputUser := new(TemplateInputUser)
	err := c.QueryParser(inputUser)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	where := ""
	if inputUser.UserId != nil {
		where = where + " AND ssr.user_id IN ( " + *inputUser.UserId + " )"
	}
	if inputUser.ChildId != nil {
		where = where + " AND ssr.id IN ( " + *inputUser.ChildId + " )"
	}
	if inputUser.ParentId != nil {
		where = where + " AND ssr.parent_id IN ( " + *inputUser.ParentId + " )"
	}
	if inputUser.Type != nil {
		where = where + " AND UPPER(ssr.aksi) = UPPER('" + *inputUser.Type + "')"
	}

	page := c.Query("page")
	pageSize := c.Query("pageSize")

	var qLimit, qPage string
	iPage, _ := strconv.Atoi(page)
	iPageSize, _ := strconv.Atoi(pageSize)

	if pageSize != "" {
		qLimit = " LIMIT " + pageSize
	} else {
		qLimit = " LIMIT 20"
		iPageSize = 20
	}

	if page == "" {
		iPage = 0
	} else {
		iPage = iPage - 1
	}

	tempQ := strconv.Itoa(iPage * iPageSize)
	qPage = " OFFSET " + tempQ

	query := fmt.Sprintf(`SELECT sq.* FROM (
		SELECT ssr.parent_id as order_id, 
			CASE WHEN MIN(ssr.is_validate) = 1 THEN 'Approve'
				WHEN MIN(ssr.is_validate) = 0 AND ssr.confirm_key IS NOT NULL THEN 'Processed'
				WHEN MIN(ssr.is_validate) = 0 AND ssr.confirm_key IS NULL THEN 'Pending' END as status,
			ssr.condition,
			ssr.gudang_id,
			ssr.user_id,
			to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS') as tanggal_riwayat,
			ssr.aksi,
			JSONB_AGG( DISTINCT
				JSONB_BUILD_OBJECT(
					'id_order_child', ssr.id,
					'produk', JSONB_BUILD_OBJECT(
											'id_produk', p.id,
											'code', p.code,
											'name', p.name,
											'foto', p.foto
										),
					'jumlah', ssr.jumlah,
					'condition', ssr.condition,
					'pita', ssr.pita,
					'gudang_id', ssr.gudang_id
				) --ORDER BY ssr.pita DESC, ssr.id
			) as datas
		FROM stok_salesman_riwayat ssr
		JOIN produk p
			ON ssr.produk_id = p.id
		JOIN stok_gudang sg
			ON p.id = sg.produk_id
			AND ssr.condition = sg.condition
			AND ssr.pita = sg.pita
		WHERE TRUE AND ssr.parent_id IS NOT NULL AND DATE(tanggal_riwayat) BETWEEN CURRENT_DATE -'1 month'::interval AND CURRENT_DATE %s
		GROUP BY ssr.parent_id, ssr.confirm_key, ssr.condition, ssr.gudang_id, ssr.user_id, to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS'), ssr.aksi
		ORDER BY to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS') DESC
		) sq

		UNION ALL

		SELECT sq.* FROM (
		SELECT ssr.parent_id as order_id, 
			CASE WHEN MIN(ssr.is_validate) = 1 THEN 'Approve'
				WHEN MIN(ssr.is_validate) = 0 AND ssr.confirm_key IS NOT NULL THEN 'Processed'
				WHEN MIN(ssr.is_validate) = 0 AND ssr.confirm_key IS NULL THEN 'Pending' END as status,
				null as condition,
			ssr.gudang_id,
			ssr.user_id,
			to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS') as tanggal_riwayat,
			ssr.aksi,
			JSONB_AGG( DISTINCT
				JSONB_BUILD_OBJECT(
					'id_order_child', ssr.id,
					'item', JSONB_BUILD_OBJECT(
											'id_produk', p.id,
											'code', p.code,
											'name', p.name
										),
					'jumlah', ssr.jumlah,
					'gudang_id', ssr.gudang_id
				) --ORDER BY ssr.pita DESC, ssr.id
			) as datas
		FROM md.stok_merchandiser_riwayat ssr
		JOIN md.item p
			ON ssr.item_id = p.id
		JOIN md.stok_gudang_item sg
			ON p.id = sg.item_id
		WHERE TRUE AND ssr.parent_id IS NOT NULL AND DATE(tanggal_riwayat) BETWEEN CURRENT_DATE -'1 month'::interval AND CURRENT_DATE %s
		GROUP BY ssr.parent_id, ssr.confirm_key, ssr.gudang_id, ssr.user_id, to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS'), ssr.aksi
		ORDER BY to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS') DESC
		) sq
		ORDER BY tanggal_riwayat DESC`, where, where)

	var wg sync.WaitGroup
	resultsChan := make(chan map[int][]map[string]interface{}, 2)

	queries := []string{
		query,
		query + qPage + qLimit,
	}

	tempResults := make([][]map[string]interface{}, len(queries))

	// Launch concurrent Goroutines
	for i, query := range queries {
		wg.Add(1)
		go helpers.ExecuteGORMQuery(query, resultsChan, i, &wg)
	}

	// Wait for all Goroutines to finish
	wg.Wait()
	close(resultsChan)

	for result := range resultsChan {
		for index, res := range result {
			tempResults[index] = res
		}
	}

	if len(tempResults) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data not found",
			Success: true,
		})
	}

	type newResponseDataMultiple struct {
		Message    string      `json:"message"`
		Success    bool        `json:"success"`
		Data       interface{} `json:"datas"`
		TotalPages int         `json:"total_pages"`
	}

	var tempTotalPages int
	if len(tempResults[0]) < iPageSize {
		tempTotalPages = 1
	} else {
		tempTotalPages = int(math.Ceil(float64(len(tempResults[0])) / float64(iPageSize)))
	}

	return c.Status(fiber.StatusOK).JSON(newResponseDataMultiple{
		Message:    "Data has been loaded successfully",
		Success:    true,
		Data:       tempResults[1],
		TotalPages: tempTotalPages,
	})
}

func GetListOrderMD(c *fiber.Ctx) error {
	type TemplateInputUser struct {
		UserId   *string `json:"userId"`
		ChildId  *string `json:"childId"`
		ParentId *string `json:"parentId"`
		Type     *string `json:"type"`
	}

	inputUser := new(TemplateInputUser)
	err := c.QueryParser(inputUser)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	where := ""
	if inputUser.UserId != nil {
		where = where + " AND ssr.user_id IN ( " + *inputUser.UserId + " )"
	}
	if inputUser.ChildId != nil {
		where = where + " AND ssr.id IN ( " + *inputUser.ChildId + " )"
	}
	if inputUser.ParentId != nil {
		where = where + " AND ssr.parent_id IN ( " + *inputUser.ParentId + " )"
	}
	if inputUser.Type != nil {
		where = where + " AND UPPER(ssr.aksi) = UPPER('" + *inputUser.Type + "')"
	}

	page := c.Query("page")
	pageSize := c.Query("pageSize")

	var qLimit, qPage string
	iPage, _ := strconv.Atoi(page)
	iPageSize, _ := strconv.Atoi(pageSize)

	if pageSize != "" {
		qLimit = " LIMIT " + pageSize
	} else {
		qLimit = " LIMIT 20"
		iPageSize = 20
	}

	if page == "" {
		iPage = 0
	} else {
		iPage = iPage - 1
	}

	tempQ := strconv.Itoa(iPage * iPageSize)
	qPage = " OFFSET " + tempQ

	query := fmt.Sprintf(`SELECT ssr.parent_id as order_id, 
							CASE WHEN MIN(ssr.is_validate) = 1 THEN 'Approve'
								WHEN MIN(ssr.is_validate) = 0 AND ssr.confirm_key IS NOT NULL THEN 'Processed'
								WHEN MIN(ssr.is_validate) = 0 AND ssr.confirm_key IS NULL THEN 'Pending' END as status,
							ssr.gudang_id,
							ssr.merchandiser_id,
							to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS') as tanggal_riwayat,
							ssr.aksi,
							JSONB_AGG( DISTINCT
								JSONB_BUILD_OBJECT(
									'id_order_child', ssr.id,
									'item', JSONB_BUILD_OBJECT(
												'id_produk', i.id,
												'code', i.code,
												'name', i.name,
												'category_id', ic.id,
												'category_name', ic.name,
												'brand_id', pb.id,
												'brand_name', pb.name
											),
									'jumlah', ssr.jumlah,
									'gudang_id', ssr.gudang_id
								) --ORDER BY ssr.pita DESC, ssr.id
							) as datas
						FROM md.stok_merchandiser_riwayat ssr
						JOIN md.item i
							ON ssr.item_id = i.id
						JOIN md.item_category ic
							ON i.category_id = ic.id
						JOIN produk_brand pb
							ON i.brand_id = pb.id
						JOIN md.stok_gudang_item sg
							ON i.id = sg.item_id
						WHERE TRUE AND ssr.parent_id IS NOT NULL AND DATE(tanggal_riwayat) BETWEEN CURRENT_DATE -'1 month'::interval AND CURRENT_DATE %s
						GROUP BY ssr.parent_id, ssr.confirm_key, ssr.gudang_id, ssr.merchandiser_id, to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS'), ssr.aksi
						ORDER BY to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS') DESC`, where)

	var wg sync.WaitGroup
	resultsChan := make(chan map[int][]map[string]interface{}, 2)

	queries := []string{
		query,
		query + qPage + qLimit,
	}

	tempResults := make([][]map[string]interface{}, len(queries))

	// Launch concurrent Goroutines
	for i, query := range queries {
		wg.Add(1)
		go helpers.ExecuteGORMQuery(query, resultsChan, i, &wg)
	}

	// Wait for all Goroutines to finish
	wg.Wait()
	close(resultsChan)

	for result := range resultsChan {
		for index, res := range result {
			tempResults[index] = res
		}
	}

	if len(tempResults) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data not found",
			Success: true,
		})
	}

	type newResponseDataMultiple struct {
		Message    string      `json:"message"`
		Success    bool        `json:"success"`
		Data       interface{} `json:"datas"`
		TotalPages int         `json:"total_pages"`
	}

	var tempTotalPages int
	if len(tempResults[0]) < iPageSize {
		tempTotalPages = 1
	} else {
		tempTotalPages = int(math.Ceil(float64(len(tempResults[0])) / float64(iPageSize)))
	}

	return c.Status(fiber.StatusOK).JSON(newResponseDataMultiple{
		Message:    "Data has been loaded successfully",
		Success:    true,
		Data:       tempResults[1],
		TotalPages: tempTotalPages,
	})
}

func PostOrder(c *fiber.Ctx) error {

	type Products struct {
		Id        *int         `json:"id"`
		Qty       *int         `json:"qty"`
		Pita      *interface{} `json:"pita"`
		Condition *string      `json:"condition"`
		Aksi      *string      `json:"aksi"`
	}

	type TemplateInputUser struct {
		Date     *string    `json:"date"`
		GudangId *int       `json:"gudang_id"`
		UserId   *int       `json:"user_id"`
		Type     *string    `json:"type"`
		Products []Products `json:"products"`
	}

	inputUser := new(TemplateInputUser)
	err := c.BodyParser(inputUser)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	var stokSalesmanRiwayat []structs.StokSalesmanRiwayat

	parentId := int64(*inputUser.UserId) + time.Now().Unix()
	for i := 0; i < len(inputUser.Products); i++ {
		// 	tempString := *inputUser.Products[i].Id + "-" + *inputUser.Products[i].Pita + "-" + *inputUser.Products[i].Condition + "-" + *inputUser.Products[i].Aksi + *inputUser.UserId
		tempString := ""
		switch pita := (*inputUser.Products[i].Pita).(type) {
		case int:
			tempString = strconv.Itoa(pita)
			// ...
		default:
			// handle the case where Pita is not an int
			tempString = fmt.Sprintf("%v", pita)
		}
		// 	if
		stokSalesmanRiwayat = append(stokSalesmanRiwayat, structs.StokSalesmanRiwayat{
			ProdukId:       int16(*inputUser.Products[i].Id),
			Jumlah:         int32(*inputUser.Products[i].Qty),
			Pita:           tempString,
			Condition:      *inputUser.Products[i].Condition,
			IsValidate:     0,
			GudangId:       int16(*inputUser.GudangId),
			UserId:         int32(*inputUser.UserId),
			TanggalRiwayat: helpers.ParseDate(*inputUser.Date).Format("2006-01-02T15:04:05"), // or any other format you need
			ParentId:       parentId,
			Aksi:           *inputUser.Products[i].Aksi,
		})
	}

	tx := db.DB.Begin()

	// fmt.Println(stokSalesmanRiwayat)

	err = tx.Create(&stokSalesmanRiwayat).Error
	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	tx.Commit()

	var whereId string
	for i := 0; i < len(stokSalesmanRiwayat); i++ {
		whereId = whereId + strconv.Itoa(int(stokSalesmanRiwayat[i].ID)) + ","
	}

	whereId = whereId[:len(whereId)-1]

	dataReturn, err := helpers.NewExecuteQuery(fmt.Sprintf(`SELECT ssr.parent_id as order_id, 
	CASE WHEN MIN(ssr.is_validate) = 1 THEN 'Approve'
		WHEN MIN(ssr.is_validate) = 0 AND ssr.confirm_key IS NOT NULL THEN 'Processed'
		WHEN MIN(ssr.is_validate) = 0 AND ssr.confirm_key IS NULL THEN 'Pending' END as status,
	ssr.condition,
	ssr.gudang_id,
	ssr.user_id,
	to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS') as tanggal_riwayat,
	ssr.aksi,
	JSONB_AGG( DISTINCT
		JSONB_BUILD_OBJECT(
			'id_order_child', ssr.id,
			'produk', JSONB_BUILD_OBJECT(
									'id_produk', p.id,
									'code', p.code,
									'name', p.name,
									'foto', p.foto
								),
			'jumlah', ssr.jumlah,
			'condition', ssr.condition,
			'pita', ssr.pita,
			'gudang_id', ssr.gudang_id
		) --ORDER BY ssr.pita DESC, ssr.id
	) as datas
FROM stok_salesman_riwayat ssr
JOIN produk p
	ON ssr.produk_id = p.id
JOIN stok_gudang sg
	ON p.id = sg.produk_id
	AND ssr.condition = sg.condition
	AND ssr.pita = sg.pita
WHERE TRUE AND ssr.id IN (%s)
GROUP BY ssr.parent_id, ssr.confirm_key, ssr.condition, ssr.gudang_id, ssr.user_id, to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS'), ssr.aksi
ORDER BY to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS') DESC`, whereId))

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan return data",
			Success: false,
		})
	}

	params := map[string]interface{}{
		"warehouseId": *inputUser.GudangId,
	}

	dataSend, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Error marshaling JSON",
			Success: false,
		})
	}

	_, err = helpers.SendCurl(dataSend, "POST", "https://api.gudangku.pt-bks.com/order/new-order-product-notify")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengirim notification",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Success insert data",
		Success: true,
		Data:    dataReturn,
	})
}

func PostOrderMD(c *fiber.Ctx) error {

	type Items struct {
		Id   *int    `json:"id"`
		Qty  *int    `json:"qty"`
		Aksi *string `json:"aksi"`
	}

	type TemplateInputUser struct {
		Date     *string `json:"date"`
		GudangId *int    `json:"gudang_id"`
		UserId   *int    `json:"user_id"`
		Type     *string `json:"type"`
		Items    []Items `json:"items"`
	}

	inputUser := new(TemplateInputUser)
	err := c.BodyParser(inputUser)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	var StokMerchandiserRiwayat []structs.StokMerchandiserRiwayat

	parentId := 5 + int64(*inputUser.UserId) + time.Now().Unix()
	for i := 0; i < len(inputUser.Items); i++ {
		// 	tempString := *inputUser.Items[i].Id + "-" + *inputUser.Items[i].Pita + "-" + *inputUser.Items[i].Condition + "-" + *inputUser.Items[i].Aksi + *inputUser.UserId

		// 	if
		StokMerchandiserRiwayat = append(StokMerchandiserRiwayat, structs.StokMerchandiserRiwayat{
			ItemId:         int16(*inputUser.Items[i].Id),
			Jumlah:         int32(*inputUser.Items[i].Qty),
			IsValidate:     0,
			GudangId:       int16(*inputUser.GudangId),
			UserId:         int32(*inputUser.UserId),
			TanggalRiwayat: helpers.ParseDate(*inputUser.Date).Format("2006-01-02T15:04:05"), // or any other format you need,
			ParentId:       parentId,
			Aksi:           *inputUser.Items[i].Aksi,
		})
	}

	tx := db.DB.Begin()

	err = tx.Create(&StokMerchandiserRiwayat).Error
	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	tx.Commit()

	var whereId string
	for i := 0; i < len(StokMerchandiserRiwayat); i++ {
		whereId = whereId + strconv.Itoa(int(StokMerchandiserRiwayat[i].ID)) + ","
	}

	whereId = whereId[:len(whereId)-1]

	dataReturn, err := helpers.NewExecuteQuery(fmt.Sprintf(`SELECT ssr.parent_id as order_id, 
			CASE WHEN MIN(ssr.is_validate) = 1 THEN 'Approve'
				WHEN MIN(ssr.is_validate) = 0 AND ssr.confirm_key IS NOT NULL THEN 'Processed'
				WHEN MIN(ssr.is_validate) = 0 AND ssr.confirm_key IS NULL THEN 'Pending' END as status,
			ssr.gudang_id,
			ssr.user_id,
			to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS') as tanggal_riwayat,
			ssr.aksi,
			JSONB_AGG( DISTINCT
				JSONB_BUILD_OBJECT(
					'id_order_child', ssr.id,
					'produk', JSONB_BUILD_OBJECT(
											'id_produk', p.id,
											'code', p.code,
											'name', p.name
										),
					'jumlah', ssr.jumlah,
					'gudang_id', ssr.gudang_id
				) --ORDER BY ssr.pita DESC, ssr.id
			) as datas
		FROM md.stok_merchandiser_riwayat ssr
		JOIN md.item p
			ON ssr.item_id = p.id
		JOIN md.stok_gudang_item sg
			ON p.id = sg.item_id
		WHERE TRUE AND ssr.id IN (%s)
		GROUP BY ssr.parent_id, ssr.confirm_key, ssr.gudang_id, ssr.user_id, to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS'), ssr.aksi
		ORDER BY to_char(ssr.tanggal_riwayat, 'YYYY-MM-DD HH24:MI:SS') DESC`, whereId))

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan return data",
			Success: false,
		})
	}

	params := map[string]interface{}{
		"warehouseId": *inputUser.GudangId,
	}

	dataSend, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Error marshaling JSON",
			Success: false,
		})
	}

	_, err = helpers.SendCurl(dataSend, "POST", "https://api.gudangku.pt-bks.com/order/new-order-item-notify")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengirim notification",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Success insert data",
		Success: true,
		Data:    dataReturn,
	})
}
