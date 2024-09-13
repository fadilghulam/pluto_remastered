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
