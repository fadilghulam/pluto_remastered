package controllers

import (
	"bytes"
	"fmt"
	"pluto_remastered/helpers"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gofiber/fiber/v2"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func GetProductTrends(c *fiber.Ctx) error {

	branchId := helpers.ParamArray(c.Context().QueryArgs().PeekMulti("branchId[]"))
	produkId := helpers.ParamArray(c.Context().QueryArgs().PeekMulti("produkId[]"))
	month, _ := strconv.Atoi(c.Query("month", "12"))
	groupingOption := c.Query("groupingOption", "PRODUK")

	var qSelect, qGroup, qOrder string
	switch groupingOption {
	case "PRODUK":
		qSelect = "p.code as product, "
		qGroup = "p.id"
		qOrder = "p.order ASC"
	case "PRODUK BRAND":
		qSelect = "pbr.name as product, "
		qGroup = "pbr.id"
		qOrder = "pbr.name ASC"
	default:
		qSelect = "p.code as product, "
		qGroup = "p.id"
		qOrder = "p.order ASC"
	}

	var qWhereProdukId string
	if len(produkId) > 0 {
		qWhereProdukId = " AND produk_id IN (" + strings.Join(produkId, ",") + ")"
		// qOnProdukId = qWhereProdukId
	}

	var qWhereBranchId string
	if len(branchId) > 0 {
		qWhereBranchId = " AND branch_id IN (" + strings.Join(branchId, ",") + ")"
		// qOnBranchId = qWhereBranchId
	}

	minDate := time.Now().AddDate(0, -month, 0).Format("2006-01-02")

	selectMonth := ""
	// var months []string

	for i := 1; i <= month; i++ {
		// Get the first date of the month, 12 months before current date, and add 'i' months
		dateTemp := time.Now().AddDate(0, -month, 0)
		currentDate := dateTemp.AddDate(0, i, 0).Format("2006-01")
		monthName := dateTemp.AddDate(0, i, 0).Format("Jan")
		yearName := strconv.Itoa(dateTemp.AddDate(0, i, 0).Year())

		// Add to months slice
		// months = append(months, currentDate)

		// Build the selectMonth string
		selectMonth += ", COALESCE(SUM(CASE WHEN year_month= DATE('" + currentDate + "-01') AND transaction_type ='PENJUALAN' THEN jumlah ELSE 0 END),0) - COALESCE(SUM(CASE WHEN year_month= DATE('" + currentDate + "-01') AND transaction_type ='RETUR' THEN jumlah ELSE 0 END),0) AS \"" + monthName + "-" + yearName + "\""
	}

	// Remove leading comma
	if len(selectMonth) > 0 {
		selectMonth = selectMonth[1:]
	}

	const sqlTemplate = `
	SELECT 
		{{.QSelect}}
		{{.SelectMonth}}
	FROM
	(
		SELECT 
			pd.produk_id AS produk_id,
			COALESCE(SUM(pd.jumlah),0) AS jumlah,
			DATE(date_part('year', DATE(p.tanggal_penjualan))||'-'||date_part('month', DATE(p.tanggal_penjualan))||'-'||'01') AS year_month,
			'PENJUALAN' AS transaction_type
		FROM penjualan p
		JOIN penjualan_detail pd ON pd.penjualan_id = p.id
		JOIN produk pr ON pr.id = pd.produk_id
		WHERE pd.harga>0 AND DATE(p.tanggal_penjualan) >= DATE('{{.MinDate}}') 
		{{.QWhereSrId}} {{.QWhereRayonId}} {{.QWhereBranchId}} {{.QWhereAreaId}} {{.QWhereSalesmanId}} {{.QWhereProductCategory}}
		GROUP BY date_part('year', DATE(p.tanggal_penjualan)), date_part('month', DATE(p.tanggal_penjualan)), pd.produk_id

		UNION ALL

		SELECT
			pd.produk_id AS produk_id,
			COALESCE(SUM(pd.jumlah),0) AS jumlah,
			DATE(date_part('year', p.tanggal_pengembalian)||'-'||date_part('month', p.tanggal_pengembalian)||'-'||'01') AS year_month,
			'RETUR' AS transaction_type
		FROM pengembalian p
		JOIN pengembalian_detail pd ON pd.pengembalian_id = p.id
		JOIN produk pr ON pr.id = pd.produk_id
		WHERE DATE(p.tanggal_pengembalian) >= DATE('{{.MinDate}}') 
		{{.QWhereSrId}} {{.QWhereRayonId}} {{.QWhereBranchId}} {{.QWhereAreaId}} {{.QWhereSalesmanId}} {{.QWhereProductCategory}}
		GROUP BY date_part('year', p.tanggal_pengembalian), date_part('month', p.tanggal_pengembalian), pd.produk_id
	) AS data
	JOIN produk p ON p.id = data.produk_id
	JOIN (SELECT produk_id FROM produk_branch pb JOIN produk p ON p.id = pb.produk_id WHERE TRUE {{.QWhereProdukId}} {{.QWhereProductCategoryIsi}} GROUP BY produk_id) pb ON pb.produk_id = p.id
	JOIN produk_brand pbr ON pbr.id = p.brand_id
	GROUP BY {{.QGroup}}
	ORDER BY {{.QOrder}}
	`

	// fmt.Println(selectMonth)

	// Create a map of data to replace the placeholders in the template
	templateQuery := map[string]interface{}{
		"QSelect":        qSelect,
		"SelectMonth":    selectMonth,
		"MinDate":        minDate,
		"QWhereBranchId": qWhereBranchId,
		"QWhereProdukId": qWhereProdukId,
		"QGroup":         qGroup,
		"QOrder":         qOrder,
	}

	tmpl, err := template.New("sqlQuery").Parse(sqlTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal build query",
			Success: false,
		})
	}

	var queryBuffer bytes.Buffer
	err = tmpl.Execute(&queryBuffer, templateQuery)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal build query",
			Success: false,
		})
	}

	finalQuery := strings.Replace(queryBuffer.String(), "<no value>", "", -1)

	data, err := helpers.ExecuteQuery2(finalQuery, "")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal execute query",
			Success: false,
		})
	}

	labels := []string{}
	for pair := data[0].Oldest(); pair != nil; pair = pair.Next() {
		if pair.Key != "product" {
			labels = append(labels, pair.Key)
		}
	}

	param := []string{}
	for _, result := range data {
		for pair := result.Oldest(); pair != nil; pair = pair.Next() {
			if pair.Key == "product" {
				param = append(param, pair.Value.(string))
			}
		}
	}

	tempData := make([]map[string]interface{}, len(data))
	for j, result := range data {

		tempData[j] = make(map[string]interface{})

		tempData[j]["name"] = param[j]
		tempData[j]["data"] = []float64{}

		for pair := result.Oldest(); pair != nil; pair = pair.Next() {
			if pair.Key != "product" {
				tempData[j]["data"] = append(tempData[j]["data"].([]float64), pair.Value.(float64))
			}
		}

		if helpers.ArraySum(tempData[j]["data"].([]float64)) < 1 {
			data = append(data[:j], data[j+1:]...)
		}
	}

	statistik := fiber.Map{
		"labels": labels,
		"values": tempData,
	}

	type customReturn struct {
		Message     string                                        `json:"message"`
		Success     bool                                          `json:"success"`
		Data        []*orderedmap.OrderedMap[string, interface{}] `json:"data"`
		Statatistik fiber.Map                                     `json:"statatistik"`
	}

	return c.Status(fiber.StatusOK).JSON(customReturn{
		Success:     true,
		Message:     "Data has been loaded",
		Data:        data,
		Statatistik: statistik,
	})
}

func GetUserBranch(c *fiber.Ctx) error {

	type TemplateInputUser struct {
		DateStart *string `json:"dateStart"`
		DateEnd   *string `json:"dateEnd"`
		Date      *string `json:"date"`
		BranchId  *string `json:"branchId"`
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

	if inputUser.DateStart == nil && inputUser.DateEnd == nil && inputUser.Date == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	if inputUser.BranchId == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	var qWhere string

	if inputUser.Date == nil {
		qWhere = fmt.Sprintf(" AND ulb.start_date >= DATE('%s') AND COALESCE(ulb.end_date, ulb.last_visit_date) <= DATE('%s')", *inputUser.DateStart, *inputUser.DateEnd)
	} else {
		qWhere = fmt.Sprintf(" AND DATE('%s') BETWEEN ulb.start_date AND COALESCE(ulb.end_date, ulb.last_visit_date)", *inputUser.Date)
	}

	// userLogBranch := []

	// datas := db.DB.Where("branch_id = ? "+qWhere, *inputUser.BranchId).
	// 				Find(&userLogBranch).
	// 				Joins("JOIN public.user ON user.id = user_log_branch.user_id").
	// 				Select(`public.user.full_name, user_log_branch.user_id,
	// 						user_log_branch.user_id_subtitute, user_log_branch.branch_id,
	// 						user_log_branch.start_date, user_log_branch.end_date,
	// 						user_log_branch.last_visit_date`)

	datas, err := helpers.NewExecuteQuery(fmt.Sprintf(`SELECT u.full_name, ulb.user_id,
							CASE WHEN ulb.user_id_subtitute = -1 THEN NULL ELSE ulb.user_id_subtitute END as user_id_subtitute,
							ulb.branch_id, 
							ulb.start_date, ulb.end_date, 
							ulb.last_visit_date
							FROM public.user_log_branch ulb
							JOIN public.user u
								ON u.id = ulb.user_id
							WHERE ulb.branch_id IN ('%s') %s`, *inputUser.BranchId, qWhere))

	if err != nil {
		fmt.Println(err.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Success: true,
		Message: "Success",
		Data:    datas,
	})
}

func TestQuery(c *fiber.Ctx) error {

	result, err := helpers.NewExecuteQuery(`SELECT sq.date, JSON_AGG(sq.x), JSON_AGG(sq.x2), 'test5'
											FROM (
											SELECT DATE('2021-08-01'), JSONB_BUILD_OBJECT('test', 'value') as x, 'test2' as x2
											UNION
											SELECT DATE('2021-08-01'), JSONB_BUILD_OBJECT('test3', 'value3'), 'test3' as x2
											) sq
											GROUP BY sq.date`)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal execute query",
			Success: false,
		})
	}

	// result2, err := helpers.ExecuteQuery2(`SELECT 1, 2,JSONB_BUILD_OBJECT('customer', JSONB_BUILD_OBJECT('test', 'value'), 'test2', '132') as customer2, ARRAY[3,4,5], 'test2'`, "")

	// if err != nil {
	// 	fmt.Println(err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
	// 		Message: "Gagal execute query",
	// 		Success: false,
	// 	})
	// }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": result,
		// "result2": result2,
	})
}

func GetDashboardOmzet(c *fiber.Ctx) error {

	start := time.Now()

	branchId := helpers.ParamArray(c.Context().QueryArgs().PeekMulti("branchId[]"))
	date := c.Query("date")

	if date == "" {
		date = "CURRENT_DATE"
	} else {
		date = " DATE('" + date + "') "
	}

	var qWhereBranchId, QWhereBranchHolderId string
	if len(branchId) > 0 {
		qWhereBranchId = " AND p.branch_id IN (" + strings.Join(branchId, ",") + ")"
		QWhereBranchHolderId = " AND bh.id IN (" + strings.Join(branchId, ",") + ")"
		// qOnBranchId = qWhereBranchId
	}

	templateReplaceQuery := map[string]interface{}{
		"QWherePbranchId":      qWhereBranchId,
		"QWhereBranchHolderId": QWhereBranchHolderId,
		"QDate":                date,
	}

	queryGetOmzet := `WITH penjualan_this_month as (
							SELECT SUM((pd.harga - pd.diskon) * pd.jumlah) as total_penjualan,
											SUM(pd.jumlah) as total_pack,
											SUM(pd.jumlah) FILTER (WHERE pd.harga <> 0) as total_pack_omzet,
											SUM(pd.jumlah) FILTER (WHERE pd.harga = 0) as total_pack_bonus
							FROM penjualan p
							JOIN penjualan_detail pd
								ON p.id = pd.penjualan_id
							WHERE 
								DATE(p.tanggal_penjualan) 
									BETWEEN DATE(date_trunc('month', {{.QDate}})) 
											AND DATE(date_trunc('month', {{.QDate}}) + '1 month'::interval - '1 day'::interval) 
								{{.QWherePbranchId}}
						), penjualan_last_month as (
							SELECT SUM((pd.harga - pd.diskon) * pd.jumlah) as total_penjualan, 
										SUM(pd.jumlah) as total_pack,
										SUM(pd.jumlah) FILTER (WHERE pd.harga <> 0) as total_pack_omzet,
										SUM(pd.jumlah) FILTER (WHERE pd.harga = 0) as total_pack_bonus
							FROM penjualan p
							JOIN penjualan_detail pd
								ON p.id = pd.penjualan_id
							WHERE 
								DATE(p.tanggal_penjualan) 
									BETWEEN DATE((date_trunc('month', {{.QDate}}) - '1 month'::interval)) 
											AND DATE((date_trunc('month', {{.QDate}})- '1 month'::interval) + '1 month'::interval - '1 day'::interval)

								{{.QWherePbranchId}}
						), pengembalian_this_month as (
							SELECT SUM(pd.harga * pd.jumlah) as total_penjualan,
										SUM(pd.jumlah) as total_pack,
										SUM(pd.jumlah) FILTER (WHERE pd.harga <> 0) as total_pack_omzet,
										SUM(pd.jumlah) FILTER (WHERE pd.harga = 0) as total_pack_bonus
							FROM pengembalian p
							JOIN pengembalian_detail pd
								ON p.id = pd.pengembalian_id
							WHERE 
								DATE(p.tanggal_pengembalian) 
									BETWEEN DATE(date_trunc('month', {{.QDate}})) 
											AND DATE(date_trunc('month', {{.QDate}}) + '1 month'::interval - '1 day'::interval)
								{{.QWherePbranchId}}
								
						), pengembalian_last_month as (
							SELECT SUM(pd.harga * pd.jumlah) as total_penjualan,
										SUM(pd.jumlah) as total_pack,
										SUM(pd.jumlah) FILTER (WHERE pd.harga <> 0) as total_pack_omzet,
										SUM(pd.jumlah) FILTER (WHERE pd.harga = 0) as total_pack_bonus
							FROM pengembalian p
							JOIN pengembalian_detail pd
								ON p.id = pd.pengembalian_id
							WHERE 
								DATE(p.tanggal_pengembalian) 
									BETWEEN DATE((date_trunc('month', {{.QDate}}) - '1 month'::interval)) 
											AND DATE((date_trunc('month', {{.QDate}})- '1 month'::interval) + '1 month'::interval - '1 day'::interval)
								{{.QWherePbranchId}}
						)

						SELECT data.otm as this_month, 
										data.olm as last_month,
										ROUND((((data.otm - data.olm)  / CASE WHEN data.olm = 0 THEN 1 ELSE data.olm END) * 100)::numeric,2) as growth
						FROM (
							SELECT COALESCE(MAX(sq.total_penjualan) FILTER (WHERE sq.flag = 'ptm'),0) - 
											COALESCE(MAX(sq.total_penjualan) FILTER (WHERE sq.flag = 'pgtm'),0) as otm,
											COALESCE(MAX(sq.total_penjualan) FILTER (WHERE sq.flag = 'plm'),0) - 
											COALESCE(MAX(sq.total_penjualan) FILTER (WHERE sq.flag = 'pglm'),0) as olm
							FROM (
								SELECT *, 'ptm' as flag
								FROM penjualan_this_month ptm

								UNION ALL

								SELECT *, 'plm' as flag
								FROM penjualan_last_month plm

								UNION ALL

								SELECT *, 'pgtm' as flag
								FROM pengembalian_this_month pgtm

								UNION ALL

								SELECT *, 'pglm' as flag
								FROM pengembalian_last_month pglm
							) sq
						) data
						`

	query1, err := helpers.PrepareQuery(queryGetOmzet, templateReplaceQuery)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	dataOmzet, err := helpers.ExecuteQuery(query1)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika mengambil data omzet",
			Success: false,
		})
	}

	queryGetPiutang := `SELECT SUM(data.pitm) as piutang_this_month, COALESCE(SUM(data.pilm),0) as piutang_last_month, ROUND((((SUM(data.pitm) - COALESCE(SUM(data.pilm),0))  / CASE WHEN SUM(data.pilm) = 0 THEN 1 ELSE SUM(data.pilm) END) * 100)::numeric,2) as growth
						FROM (
						SELECT
							SUM( CASE WHEN (DATE_PART('day',{{.QDate}}::timestamp -DATE(pi.tanggal_piutang)::timestamp) > 90 )
								THEN total_piutang-COALESCE(ppd.nominal,0) ELSE 0 END ) AS pitm, 
		--                     SUM( CASE WHEN (DATE_PART('day',{{.QDate}}::timestamp -DATE(pi.tanggal_piutang)::timestamp) > 90 )
		--                          THEN total_piutang-COALESCE(ppd.nominal,0) ELSE 0 END ) AS pilm,
												0 as pilm,
							SUM(total_piutang-COALESCE(ppd.nominal,0)) AS total
								
							FROM
								piutang pi
							LEFT JOIN
								(SELECT piutang_id, SUM(nominal) as nominal 
								FROM pembayaran_piutang pp
								JOIN pembayaran_piutang_detail ppd 
								ON ppd.pembayaran_piutang_id = pp.id
								WHERE DATE(pp.tanggal_pembayaran) <= {{.QDate}}
								GROUP BY piutang_id) ppd 
								ON ppd.piutang_id = pi.id
							JOIN penjualan p
							ON p.id = pi.penjualan_id
							JOIN customer c
							ON c.id = p.customer_id  AND c.is_kasus IN ( 0 )
							JOIN salesman se
							ON se.id = p.salesman_id

							LEFT JOIN
							area ae
							ON ae.id = p.area_id 
							LEFT JOIN branch be
							ON be.id = p.branch_id 
							LEFT JOIN rayon re 
							ON re.id = p.rayon_id 
							LEFT JOIN sr sre
							ON sre.id = p.sr_id 

							JOIN salesman sh
							ON c.salesman_id = sh.id
							LEFT JOIN area ah
							ON ah.id = ANY(sh.area_id) AND ah.id = p.area_id
							LEFT JOIN branch bh
							ON sh.branch_id = bh.id 
							LEFT JOIN rayon rh
							ON rh.id = bh.rayon_id 
							LEFT JOIN sr srh 
							ON srh.id = rh.sr_id 

							WHERE
								DATE(pi.tanggal_piutang) <= {{.QDate}} 
								AND c.is_kasus IN ( 0 )
								AND rh.id <> 501
								{{.QWhereBranchHolderId}}

							UNION ALL

							SELECT
								0 as pitm,
								SUM( CASE WHEN (DATE_PART('day',({{.QDate}}-'1 month'::interval)::timestamp -DATE(pi.tanggal_piutang)::timestamp) > 90 )
									THEN total_piutang-COALESCE(ppd.nominal,0) ELSE 0 END ) AS pilm,
								SUM(total_piutang-COALESCE(ppd.nominal,0)) AS total
								
							FROM
								piutang pi
							LEFT JOIN
								(SELECT piutang_id, SUM(nominal) as nominal 
								FROM pembayaran_piutang pp
								JOIN pembayaran_piutang_detail ppd 
								ON ppd.pembayaran_piutang_id = pp.id
								WHERE DATE(pp.tanggal_pembayaran) <= ({{.QDate}}-'1 month'::interval)
								GROUP BY piutang_id) ppd 
								ON ppd.piutang_id = pi.id
							JOIN penjualan p
							ON p.id = pi.penjualan_id
							JOIN customer c
							ON c.id = p.customer_id  AND c.is_kasus IN ( 0 )
							JOIN salesman se
							ON se.id = p.salesman_id

							LEFT JOIN
							area ae
							ON ae.id = p.area_id 
							LEFT JOIN branch be
							ON be.id = p.branch_id 
							LEFT JOIN rayon re 
							ON re.id = p.rayon_id 
							LEFT JOIN sr sre
							ON sre.id = p.sr_id 

							JOIN salesman sh
							ON c.salesman_id = sh.id
							LEFT JOIN area ah
							ON ah.id = ANY(sh.area_id) AND ah.id = p.area_id
							LEFT JOIN branch bh
							ON sh.branch_id = bh.id 
							LEFT JOIN rayon rh
							ON rh.id = bh.rayon_id 
							LEFT JOIN sr srh 
							ON srh.id = rh.sr_id 

							WHERE
								DATE(pi.tanggal_piutang) <= ({{.QDate}}-'1 month'::interval)
								AND c.is_kasus IN ( 0 )
								AND rh.id <> 501
								{{.QWhereBranchHolderId}}
						) data`

	query2, err := helpers.PrepareQuery(queryGetPiutang, templateReplaceQuery)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	dataPiutang, err := helpers.ExecuteQuery(query2)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika mengambil data piutang",
			Success: false,
		})
	}

	returnData := make(map[string]interface{})

	if len(dataOmzet) > 0 {
		returnData["omzet"] = dataOmzet[0]
	}

	if len(dataPiutang) > 0 {
		returnData["receiveable"] = dataPiutang[0]
	}

	elapsed := time.Since(start)

	type Response struct {
		Message string        `json:"message"`
		Success bool          `json:"success"`
		Data    interface{}   `json:"data"`
		Elapsed time.Duration `json:"elapsed"`
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Message: "Success",
		Success: true,
		Data:    returnData,
		Elapsed: time.Duration(elapsed.Seconds()),
	})

}

func GetReceiveableDetail(c *fiber.Ctx) error {

	start := time.Now()

	branchId := helpers.ParamArray(c.Context().QueryArgs().PeekMulti("branchId[]"))
	date := c.Query("date")

	if date == "" {
		date = "CURRENT_DATE"
	} else {
		date = " DATE('" + date + "') "
	}

	var qWhereBranchId, QWhereBranchHolderId string
	if len(branchId) > 0 {
		qWhereBranchId = " AND p.branch_id IN (" + strings.Join(branchId, ",") + ")"
		QWhereBranchHolderId = " AND bh.id IN (" + strings.Join(branchId, ",") + ")"
		// qOnBranchId = qWhereBranchId
	}

	templateReplaceQuery := map[string]interface{}{
		"QWherePbranchId":      qWhereBranchId,
		"QWhereBranchHolderId": QWhereBranchHolderId,
		"QDate":                date,
	}

	queryGetPiutang := `SELECT 
							/*
							JSONB_BUILD_OBJECT(
								'id', b.id,
								'name', b.name
							) as branch,
							JSONB_BUILD_OBJECT(
								'pack', data.pack_last_month,
								'nominal', data.piutang_last_month
							) as last_month,
							JSONB_BUILD_OBJECT(
								'pack', data.pack_this_month,
								'nominal', data.piutang_this_month
							) as this_month,
							JSONB_BUILD_OBJECT(
								'pack', data.growth_pack,
								'nominal', data.growth
							) as growth_percentage
							 */
							b.id as branch_id,
							b.name as branch_name,
							data.piutang_last_month as last_month,
							data.piutang_this_month as this_month,
							data.growth as growth
						FROM (
							SELECT COALESCE(SUM(data.pitm),0) as piutang_this_month, 
									--COALESCE(SUM(data.packtm),0) as pack_this_month,
									COALESCE(SUM(data.pilm),0) as piutang_last_month, 
									--COALESCE(SUM(data.packlm),0) as pack_last_month,
									ROUND((((COALESCE(SUM(data.pitm),0) - COALESCE(SUM(data.pilm),0))  / CASE WHEN SUM(data.pilm) = 0 THEN 1 ELSE SUM(data.pilm) END) * 100)::numeric,2) as growth,
									--ROUND((((COALESCE(SUM(data.packtm),0) - COALESCE(SUM(data.packlm),0)) / CASE WHEN SUM(data.packlm) = 0 THEN 1 ELSE SUM(data.packlm) END) * 100)::numeric,2) as growth_pack,
									data.branch_id
					
								FROM (
								SELECT
									SUM( CASE WHEN (DATE_PART('day',{{.QDate}}::timestamp -DATE(pi.tanggal_piutang)::timestamp) > 90 )
										THEN total_piutang-COALESCE(ppd.nominal,0) ELSE 0 END ) AS pitm, 
				--                     SUM( CASE WHEN (DATE_PART('day',{{.QDate}}::timestamp -DATE(pi.tanggal_piutang)::timestamp) > 90 )
				--                          THEN total_piutang-COALESCE(ppd.nominal,0) ELSE 0 END ) AS pilm,
														0 as pilm,
									SUM(total_piutang-COALESCE(ppd.nominal,0)) AS total,
									--SUM(pd.jumlah) as packtm,
									--0 as packlm,
									sh.branch_id
										
									FROM
										piutang pi
									LEFT JOIN
										(SELECT piutang_id, SUM(nominal) as nominal 
										FROM pembayaran_piutang pp
										JOIN pembayaran_piutang_detail ppd 
										ON ppd.pembayaran_piutang_id = pp.id
										WHERE DATE(pp.tanggal_pembayaran) <= {{.QDate}}
										GROUP BY piutang_id) ppd 
										ON ppd.piutang_id = pi.id
									JOIN penjualan p
									ON p.id = pi.penjualan_id
									--JOIN penjualan_detail pd
									--ON p.id = pd.penjualan_id
									JOIN customer c
									ON c.id = p.customer_id  AND c.is_kasus IN ( 0 )
									JOIN salesman se
									ON se.id = p.salesman_id

									LEFT JOIN
									area ae
									ON ae.id = p.area_id 
									LEFT JOIN branch be
									ON be.id = p.branch_id 
									LEFT JOIN rayon re 
									ON re.id = p.rayon_id 
									LEFT JOIN sr sre
									ON sre.id = p.sr_id 

									JOIN salesman sh
									ON c.salesman_id = sh.id
									LEFT JOIN area ah
									ON ah.id = ANY(sh.area_id) AND ah.id = p.area_id
									LEFT JOIN branch bh
									ON sh.branch_id = bh.id 
									LEFT JOIN rayon rh
									ON rh.id = bh.rayon_id 
									LEFT JOIN sr srh 
									ON srh.id = rh.sr_id 

									WHERE
										DATE(pi.tanggal_piutang) <= {{.QDate}} 
										AND c.is_kasus IN ( 0 )
										AND rh.id <> 501
										{{.QWhereBranchHolderId}}

									GROUP BY sh.branch_id

									UNION ALL

									SELECT
										0 as pitm,
										SUM( CASE WHEN (DATE_PART('day',({{.QDate}}-'1 month'::interval)::timestamp -DATE(pi.tanggal_piutang)::timestamp) > 90 )
											THEN total_piutang-COALESCE(ppd.nominal,0) ELSE 0 END ) AS pilm,
										SUM(total_piutang-COALESCE(ppd.nominal,0)) AS total,
										--0 as packtm,
										--SUM(pd.jumlah) as packlm,
										sh.branch_id
										
									FROM
										piutang pi
									LEFT JOIN
										(SELECT piutang_id, SUM(nominal) as nominal 
										FROM pembayaran_piutang pp
										JOIN pembayaran_piutang_detail ppd 
										ON ppd.pembayaran_piutang_id = pp.id
										WHERE DATE(pp.tanggal_pembayaran) <= ({{.QDate}}-'1 month'::interval)
										GROUP BY piutang_id) ppd 
										ON ppd.piutang_id = pi.id
									JOIN penjualan p
									ON p.id = pi.penjualan_id
									--JOIN penjualan_detail pd
									--ON p.id = pd.penjualan_id
									JOIN customer c
									ON c.id = p.customer_id  AND c.is_kasus IN ( 0 )
									JOIN salesman se
									ON se.id = p.salesman_id

									LEFT JOIN
									area ae
									ON ae.id = p.area_id 
									LEFT JOIN branch be
									ON be.id = p.branch_id 
									LEFT JOIN rayon re 
									ON re.id = p.rayon_id 
									LEFT JOIN sr sre
									ON sre.id = p.sr_id 

									JOIN salesman sh
									ON c.salesman_id = sh.id
									LEFT JOIN area ah
									ON ah.id = ANY(sh.area_id) AND ah.id = p.area_id
									LEFT JOIN branch bh
									ON sh.branch_id = bh.id 
									LEFT JOIN rayon rh
									ON rh.id = bh.rayon_id 
									LEFT JOIN sr srh 
									ON srh.id = rh.sr_id 

									WHERE
										DATE(pi.tanggal_piutang) <= ({{.QDate}}-'1 month'::interval)
										AND c.is_kasus IN ( 0 )
										AND rh.id <> 501
										{{.QWhereBranchHolderId}}

									GROUP BY sh.branch_id
								) data
								GROUP BY data.branch_id
							) data
							JOIN branch b
								ON data.branch_id = b.id`

	query2, err := helpers.PrepareQuery(queryGetPiutang, templateReplaceQuery)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	// fmt.Println(query2)

	dataPiutang, err := helpers.ExecuteQuery(query2)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika mengambil data piutang",
			Success: false,
		})
	}

	if len(dataPiutang) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data tidak ditemukan",
			Success: true,
		})
	}

	elapsed := time.Since(start)

	type Response struct {
		Message string        `json:"message"`
		Success bool          `json:"success"`
		Data    interface{}   `json:"data"`
		Elapsed time.Duration `json:"elapsed"`
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Message: "Berhasil mendapatkan data",
		Success: true,
		Data:    dataPiutang,
		Elapsed: time.Duration(elapsed.Seconds()),
	})
}

func GetCustomerAR(c *fiber.Ctx) error {

	start := time.Now()

	branchId := helpers.ParamArray(c.Context().QueryArgs().PeekMulti("branchId[]"))
	date := c.Query("date")

	if date == "" {
		date = "CURRENT_DATE"
	} else {
		date = " DATE('" + date + "') "
	}

	var qWhereBranchId, QWhereBranchHolderId string
	if len(branchId) > 0 {
		qWhereBranchId = " AND p.branch_id IN (" + strings.Join(branchId, ",") + ")"
		QWhereBranchHolderId = " AND bh.id IN (" + strings.Join(branchId, ",") + ")"
		// qOnBranchId = qWhereBranchId
	}

	templateReplaceQuery := map[string]interface{}{
		"QWherePbranchId":      qWhereBranchId,
		"QWhereBranchHolderId": QWhereBranchHolderId,
		"QDate":                date,
	}

	queryGetPiutang := `SELECT
							data.*
						FROM
							(SELECT
								DATE(pi.tanggal_piutang) AS tanggal,
								pay.tipe as tipe_nota,
								JSON_BUILD_OBJECT('nomor', COALESCE(p.no_nota,'UNKNOWN'),'foto', p.image_nota) AS nota,
								DATE_PART('day',{{.QDate}}::timestamp -DATE(pi.tanggal_piutang)::timestamp) as umur_piutang,
								JSON_BUILD_OBJECT('id', srh.id, 'name', srh.name, 'alias', srh.alias) AS sr_pemilik,
								JSON_BUILD_OBJECT('id', rh.id, 'name', rh.name, 'alias', rh.name ||' ('||srh.alias||')') AS rayon_pemilik,
								JSON_BUILD_OBJECT('id', bh.id, 'name', bh.name, 'type', bh.type) AS branch_pemilik,
								JSON_BUILD_OBJECT('id', ah.id, 'name', ah.name) AS area_pemilik,
								JSON_BUILD_OBJECT('id', MAX(sh.id),'nama', MAX(sh.name)) AS salesman_pemilik,
								JSON_BUILD_OBJECT('id', sre.id, 'name', sre.name, 'alias', sre.alias) AS sr_pelaksana,
								JSON_BUILD_OBJECT('id', re.id, 'name', re.name, 'alias', re.name ||' ('||sre.alias||')') AS rayon_pelaksana,
								JSON_BUILD_OBJECT('id', be.id, 'name', be.name, 'type', be.type) AS branch_pelaksana,
								JSON_BUILD_OBJECT('id', ae.id, 'name', ae.name) AS area_pelaksana,
								JSON_BUILD_OBJECT('id', MAX(se.id),'nama', MAX(se.name)) AS salesman_pelaksana,
								c.id||'' AS customer_id,
								JSON_BUILD_OBJECT('id',c.id,
									'nama', LOWER(c.name),
									'toko', LOWER(c.outlet_name)) AS customer,
								(CASE WHEN c.tipe= 2 THEN 'sub grosir' WHEN c.tipe = 1 THEN 'grosir' ELSE 'retail' END) AS tipe,
								c.phone AS telepon, 
								LOWER(c.alamat) AS alamat,
								LOWER(c.kelurahan) AS kelurahan, 
								LOWER(c.kecamatan) AS kecamatan, 
								LOWER(c.kabupaten) AS kota,
								LOWER(c.provinsi) AS provinsi,
								c.image_toko AS foto,
								SUM(total_piutang) as total_piutang,
								SUM(COALESCE(ppd.nominal,0)) as angsuran,
								SUM(total_piutang-COALESCE(ppd.nominal,0)) as sisa_piutang,
								JSON_BUILD_OBJECT('coordinate', c.latitude_longitude, 'is_verif', c.is_verifikasi_lokasi) AS lokasi,
								p.id||'' AS penjualan_id,
								pi.id||'' AS piutang_id                    
								-- JSON_AGG(JSON_BUILD_OBJECT('metode pembayaran','Unknown')) AS informasi.
								-- JSON_AGG(array_to_string(ppd.information, ',',null)) AS information
								-- JSON_AGG(JSON_OBJECT(tipe_key,tipe_value)) AS information

							FROM
								piutang pi
							LEFT JOIN
								(SELECT piutang_id, SUM(ppd.nominal) as nominal,
								ARRAY_AGG(CASE WHEN py.id IS NOT NULL THEN py.tipe WHEN py.id IS NULL AND pp.id IS NOT NULL THEN 'CASH' END) AS information
								FROM pembayaran_piutang pp
								JOIN
								customer c
								ON c.id = pp.customer_id
								JOIN pembayaran_piutang_detail ppd 
								ON ppd.pembayaran_piutang_id = pp.id
								LEFT JOIN payment py
								ON py.id = pp.payment_id
								WHERE DATE(pp.tanggal_pembayaran) <= {{.QDate}}
								GROUP BY piutang_id) ppd 
								ON ppd.piutang_id = pi.id
							JOIN
							penjualan p
							ON p.id = pi.penjualan_id
							LEFT JOIN payment pay
							ON p.id = pay.penjualan_id
							AND pay.tipe = 'BILYET GIRO'
							JOIN customer c
							ON c.id = p.customer_id  AND c.is_kasus IN ( 0 )
							JOIN salesman se
							ON se.id = p.salesman_id 

							LEFT JOIN
							area ae
							ON ae.id = p.area_id 
							LEFT JOIN branch be
							ON be.id = p.branch_id 
							LEFT JOIN rayon re 
							ON re.id = p.rayon_id 
							LEFT JOIN sr sre
							ON sre.id = p.sr_id 

							JOIN salesman sh
							ON c.salesman_id = sh.id 
							LEFT JOIN area ah
							ON ah.id = ANY(sh.area_id) AND ah.id = p.area_id 
							LEFT JOIN branch bh
							ON sh.branch_id = bh.id  {{.QWhereBranchHolderId}}
							LEFT JOIN rayon rh
							ON rh.id = bh.rayon_id 
							LEFT JOIN sr srh 
							ON srh.id = rh.sr_id 
							WHERE
								DATE(pi.tanggal_piutang) <= {{.QDate}}
								AND rh.id <> 501
								{{.QWhereBranchHolderId}}  
								AND c.is_kasus IN ( 0 ) AND DATE_PART('day',{{.QDate}}::timestamp - DATE(pi.tanggal_piutang)::timestamp) > 90
							GROUP BY
								srh.id, rh.id, bh.id, ah.id, sh.id, sre.id, re.id, be.id, ae.id, se.id , c.id, p.id, pi.id, pay.id) as data

						WHERE
							data.sisa_piutang != 0
						ORDER BY umur_piutang DESC`

	query2, err := helpers.PrepareQuery(queryGetPiutang, templateReplaceQuery)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	// fmt.Println(query2)

	dataPiutang, err := helpers.ExecuteQuery(query2)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika mengambil data piutang",
			Success: false,
		})
	}

	if len(dataPiutang) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data tidak ditemukan",
			Success: true,
		})
	}

	elapsed := time.Since(start)

	type Response struct {
		Message string        `json:"message"`
		Success bool          `json:"success"`
		Data    interface{}   `json:"data"`
		Elapsed time.Duration `json:"elapsed"`
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Message: "Berhasil mendapatkan data",
		Success: true,
		Data:    dataPiutang,
		Elapsed: time.Duration(elapsed.Seconds()),
	})
}
