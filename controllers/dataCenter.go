package controllers

import (
	"fmt"
	db "pluto_remastered/config"
	"pluto_remastered/helpers"
	"strconv"

	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetDataOmzetSetoran(c *fiber.Ctx) error {
	branchId := helpers.ParamArray(c.Context().QueryArgs().PeekMulti("branchId[]"))
	dateStart := c.Query("dateStart")
	dateEnd := c.Query("dateEnd")
	groupingOption := c.Query("groupingOption", "branch")
	groupingOption = strings.ToLower(groupingOption)

	qWhereBranchId := ""
	if len(branchId) > 0 {
		qWhereBranchId = " AND p.branch_id IN (" + strings.Join(branchId, ",") + ")"
		// qOnBranchId = qWhereBranchId
	}

	selectOption := ""
	groupingQuery := ""
	if groupingOption == "branch" {
		selectOption = `JSONB_BUILD_OBJECT('id', sr.id, 'name', sr.name) as sr,
						JSONB_BUILD_OBJECT('id', rayon.id, 'name', rayon.name) as rayon,
						JSONB_BUILD_OBJECT('id', branch.id, 'name', branch.name) as branch,`
		groupingQuery = "GROUP BY sr.id, rayon.id, branch.id"
	} else if groupingOption == "rayon" {
		selectOption = `JSONB_BUILD_OBJECT('id', sr.id, 'name', sr.name) as sr,
						JSONB_BUILD_OBJECT('id', rayon.id, 'name', rayon.name) as rayon,`
		groupingQuery = "GROUP BY sr.id, rayon.id"
	} else {
		selectOption = `JSONB_BUILD_OBJECT('id', sr.id, 'name', sr.name) as sr,`
		groupingQuery = "GROUP BY sr.id"
	}

	templateReplaceQuery := map[string]interface{}{
		"QWhereBranchId":  qWhereBranchId,
		"QDateStart":      dateStart,
		"QDateEnd":        dateEnd,
		"QSelectOption":   selectOption,
		"QGroupingOption": groupingQuery,
	}

	queries := `WITH penjualans as (
					SELECT p.sr_id, 
									p.rayon_id, 
									p.branch_id,
									COALESCE(SUM((pd.harga - pd.diskon) * pd.jumlah) FILTER (WHERE p.is_kredit = 0),0) as total_penjualan_tunai,
									COALESCE(SUM((pd.harga - pd.diskon) * pd.jumlah) FILTER (WHERE p.is_kredit = 1),0) as total_penjualan_kredit,
									COALESCE(SUM(pd.jumlah) FILTER (WHERE p.is_kredit = 0),0) as pack_tunai,
									COALESCE(SUM(pd.jumlah) FILTER (WHERE p.is_kredit = 1),0) as pack_kredit
					FROM penjualan p
					JOIN penjualan_detail pd
						ON p.id = pd.penjualan_id
					WHERE DATE(p.tanggal_penjualan) BETWEEN DATE('{{.QDateStart}}') AND DATE('{{.QDateEnd}}') {{.QWhereBranchId}}
					GROUP BY p.sr_id, p.rayon_id, p.branch_id
				), pengembalians as (
					SELECT p.sr_id, 
									p.rayon_id, 
									p.branch_id,
									COALESCE(SUM(pd.harga * pd.jumlah),0) as total_pengembalian,
									COALESCE(SUM(pd.jumlah),0) as pack_pengembalian
					FROM pengembalian p
					JOIN pengembalian_detail pd
						ON p.id = pd.pengembalian_id
					JOIN penjualan pj
						ON p.penjualan_id = pj.id
					WHERE DATE(p.tanggal_pengembalian) BETWEEN DATE('{{.QDateStart}}') AND DATE('{{.QDateEnd}}') {{.QWhereBranchId}}
					GROUP BY p.sr_id, p.rayon_id, p.branch_id
				), pembayarans as (
					SELECT p.sr_id, 
									p.rayon_id, 
									p.branch_id,
									COALESCE(SUM(p.total_pembayaran),0) as total_pembayaran
					FROM pembayaran_piutang p
					WHERE DATE(p.tanggal_pembayaran) BETWEEN DATE('{{.QDateStart}}') AND DATE('{{.QDateEnd}}') {{.QWhereBranchId}}
					GROUP BY p.sr_id, p.rayon_id, p.branch_id
				)

				SELECT 
					{{.QSelectOption}}
					SUM(data.penjualan_tunai) as penjualan_tunai,
					SUM(data.pack_tunai) as pack_tunai,
					SUM(data.penjualan_kredit) as penjualan_kredit,
					SUM(data.pack_kredit) as pack_kredit,
					SUM(data.total_pengembalian) as total_pengembalian,
					SUM(data.pack_pengembalian) as pack_pengembalian,
					SUM(data.total_pembayaran) as total_pembayaran,
					SUM(data.penjualan_tunai) + SUM(data.penjualan_kredit) - SUM(data.total_pengembalian) as omzet,
					SUM(data.pack_tunai) + SUM(data.pack_kredit) - SUM(data.pack_pengembalian) as omzet_pack,
					SUM(data.penjualan_tunai) + SUM(data.total_pembayaran) - SUM(data.total_pengembalian) as setoran
				FROM (
					SELECT p.sr_id,
									p.rayon_id,
									p.branch_id,
									COALESCE(p.total_penjualan_tunai,0) as penjualan_tunai,
									COALESCE(p.pack_tunai,0) as pack_tunai,
									COALESCE(p.total_penjualan_kredit,0) as penjualan_kredit,
									COALESCE(p.pack_kredit,0) as pack_kredit,
									COALESCE(pg.total_pengembalian,0) as total_pengembalian,
									COALESCE(pg.pack_pengembalian,0) as pack_pengembalian,
									COALESCE(pb.total_pembayaran,0) as total_pembayaran
					FROM penjualans p
					FULL JOIN pengembalians pg
						ON p.sr_id = pg.sr_id
						AND p.rayon_id = pg.rayon_id
						AND p.branch_id = pg.branch_id
					FULL JOIN pembayarans pb
						ON p.sr_id = pb.sr_id
						AND p.rayon_id = pb.rayon_id
						AND p.branch_id = pb.branch_id
				) data
				JOIN sr
					ON data.sr_id = sr.id
				JOIN rayon
					ON data.rayon_id = rayon.id
				JOIN branch
					ON data.branch_id = branch.id
				{{.QGroupingOption}}`

	query1, err := helpers.PrepareQuery(queries, templateReplaceQuery)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	datas, err := helpers.ExecuteQuery(query1)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data",
			Success: false,
		})
	}

	if len(datas) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data tidak ditemukan",
			Success: false,
		})
	}

	if c.Query("isQuery") == "true" {
		return c.Status(fiber.StatusOK).JSON(helpers.Response{
			Data:    query1,
			Message: "Berhasil debug",
			Success: true,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Data:    datas,
		Message: "Berhasil mendapatkan data",
		Success: true,
	})
}

func GetKnowledgeBase(c *fiber.Ctx) error {
	resultSr, err := helpers.ExecuteQueryBot(`SELECT sr.*, JSONB_AGG(CASE WHEN sr.id = 8 THEN 100 ELSE b.id END) as branches
												FROM sr
												LEFT JOIN rayon r
													ON sr.id = r.sr_id
												LEFT JOIN branch b
													ON r.id = b.rayon_id
												GROUP BY sr.id`)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data sr",
			Success: false,
		})
	}

	resultRayon, err := helpers.ExecuteQueryBot(`SELECT r.id, r.name, JSONB_AGG(CASE WHEN r.id = 3 THEN 100 ELSE b.id END) as branches
													FROM rayon r
													LEFT JOIN branch b
														ON r.id = b.rayon_id
													GROUP BY r.id`)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data rayon",
			Success: false,
		})
	}

	resultBranch, err := helpers.ExecuteQueryBot(`SELECT b.id, b.name
													FROM branch b`)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data branch",
			Success: false,
		})
	}

	returnData := make(map[string]interface{})
	returnData["sr"] = resultSr
	returnData["rayon"] = resultRayon
	returnData["branch"] = resultBranch

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Data:    returnData,
		Message: "Berhasil mendapatkan data",
		Success: true,
	})
}

func GetDictionaries(c *fiber.Ctx) error {
	results, err := helpers.ExecuteQueryBot("SELECT * FROM dictionary")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data",
			Success: false,
		})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data tidak ditemukan",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Data:    results,
		Message: "Berhasil mendapatkan data",
		Success: true,
	})
}

func GetDictionariesText(c *fiber.Ctx) error {
	results, err := helpers.ExecuteQueryBot("SELECT * FROM dictionary")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data",
			Success: false,
		})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data tidak ditemukan",
			Success: false,
		})
	}

	var responseString []string

	for _, data := range results {
		var resultString strings.Builder
		resultString.WriteString("[function name]\n")
		resultString.WriteString("- [endpoint url]\n")
		resultString.WriteString("- [method]\n")
		resultString.WriteString("- [tags]\n")
		resultString.WriteString("- [description]\n")
		resultString.WriteString("- [required_parameters]\n")
		resultString.WriteString("- [not_required_parameters]\n")
		resultString.WriteString("\n")
		resultString.WriteString(fmt.Sprintf("%s\n", data["name"].(string)))
		resultString.WriteString(fmt.Sprintf("- %s\n", data["url"].(string)))
		resultString.WriteString(fmt.Sprintf("- %s\n", data["method"].(string))) // Assuming method is GET
		tags := make([]string, len(data["tags"].([]interface{})))
		for i, tag := range data["tags"].([]interface{}) {
			tags[i] = fmt.Sprintf("%v", tag)
		}
		resultString.WriteString(fmt.Sprintf("- %s\n", strings.Join(tags, ", ")))
		resultString.WriteString(fmt.Sprintf("- %s\n", data["description"].(string)))

		// Generate required parameters string
		var requiredParams []string
		var notRequiredParams []string
		for _, param := range data["params"].([]interface{}) {
			paramMap, ok := param.(map[string]interface{})
			if !ok {
				// Handle the case where param is not a map
				fmt.Println("param is not a map")
				continue
			}
			// fmt.Println(paramMap["name"].(string))
			// fmt.Println(paramMap["dataType"].(string))
			if paramMap["is_required"].(bool) {
				requiredParams = append(requiredParams, fmt.Sprintf("%s (%s)", paramMap["name"].(string), paramMap["dataType"].(string)))
			} else {
				notRequiredParams = append(notRequiredParams, fmt.Sprintf("%s (%s)", paramMap["name"].(string), paramMap["dataType"].(string)))
			}
		}
		resultString.WriteString(fmt.Sprintf("- %s\n", strings.Join(requiredParams, ", ")))
		resultString.WriteString(fmt.Sprintf("- %s\n", strings.Join(notRequiredParams, ", ")))

		responseString = append(responseString, resultString.String())
	}

	c.Set("Content-Type", "text/plain")
	return c.Status(fiber.StatusOK).SendString(strings.Join(responseString, "\n\n"))
}

func GetIncompleteData(c *fiber.Ctx) error {

	paramDate := c.Query("date")
	if paramDate == "" {
		paramDate = "CURRENT_DATE"
	} else {
		paramDate = "DATE('" + paramDate + "')"
	}

	qWhereSalesman := " AND p.user_id <> 933 "
	paramSalesman := c.Query("salesmanId")
	if paramSalesman != "" {
		qWhereSalesman += fmt.Sprintf(" AND p.salesman_id = %v", paramSalesman)
	}

	query := `SELECT p.user_id as subject_id, 'user' as subject_name, p.id||'' as id, 'penjualan tanpa detail' as description
											FROM penjualan p
											LEFT JOIN penjualan_detail pd
												ON p.id = pd.penjualan_id
											WHERE DATE(p.tanggal_penjualan) = {{.QParamDate}} AND pd.id IS NULL {{.QWhereSalesman}}

											UNION ALL

											SELECT p.user_id as subject_id, 'user' as subject_name, pd.id||'' as id, 'detail tanpa penjualan' as description
											FROM penjualan_detail pd
											LEFT JOIN penjualan p
												ON pd.penjualan_id = p.id
											WHERE DATE(pd.dtm_crt) = {{.QParamDate}} AND p.id IS NULL {{.QWhereSalesman}}

											UNION ALL

											SELECT p.user_id as subject_id, 'user' as subject_name, p.id||'' as id, 'payment tanpa transaksi (penjualan)' as description
											FROM payment p
											LEFT JOIN penjualan pj
												ON p.penjualan_id = pj.id
											WHERE DATE(p.tanggal_transaksi) = {{.QParamDate}} AND pj.id IS NULL {{.QWhereSalesman}}

											UNION ALL

											SELECT p.user_id as subject_id, 'user' as subject_name, p.id||'' as id, 'payment tanpa transaksi (pembayaran)' as description
											FROM payment p
											LEFT JOIN pembayaran_piutang pp
												ON p.id = pp.payment_id
											LEFT JOIN penjualan pj
												ON p.penjualan_id = pj.id
											WHERE DATE(p.tanggal_transaksi) = {{.QParamDate}} AND pp.id IS NULL AND pj.is_kredit = 1 {{.QWhereSalesman}}

											UNION ALL

											SELECT p.user_id as subject_id, 'user' as subject_name, pj.id||'' as id, 'payment totalnya beda dengan transaksi (penjualan non kredit)' as description
											FROM payment p
											LEFT JOIN
											( SELECT p.id, p.total_penjualan, SUM((pjd.harga - pjd.diskon) * pjd.jumlah) as total_detail
												FROM penjualan p
												LEFT JOIN penjualan_detail pjd
													ON p.id = pjd.penjualan_id
												WHERE DATE(p.tanggal_penjualan) = {{.QParamDate}} AND p.is_kredit = 0 {{.QWhereSalesman}}
												GROUP BY p.id
											) pj
											ON p.penjualan_id = pj.id
											WHERE DATE(p.tanggal_transaksi) = {{.QParamDate}} {{.QWhereSalesman}}
											GROUP BY p.user_id, pj.id
											HAVING SUM(p.nominal) <> MAX(pj.total_penjualan) OR SUM(p.nominal) <> MAX(pj.total_detail)

											/*UNION ALL

											SELECT p.user_id as subject_id, 'user' as subject_name, p.id||'' as id, 'payment totalnya beda dengan transaksi (penjualan kredit)' as description
											FROM payment p
											LEFT JOIN penjualan pj
												ON p.penjualan_id = pj.id
											LEFT JOIN penjualan_detail pjd
												ON pj.id = pjd.penjualan_id
											WHERE DATE(p.tanggal_transaksi) = {{.QParamDate}} AND pj.is_kredit = 1 {{.QWhereSalesman}}
											GROUP BY p.user_id, p.id, pj.id
											HAVING p.nominal <> pj.total_penjualan OR p.nominal <> SUM((pjd.harga - pjd.diskon) * pjd.jumlah)*/

											UNION ALL

											SELECT p.user_id as subject_id, 'user' as subject_name, p.id||'' as id, 'payment totalnya beda dengan transaksi (pembayaran)' as description
											FROM payment p
											LEFT JOIN pembayaran_piutang pp
												ON p.id = pp.payment_id
											LEFT JOIN pembayaran_piutang_detail ppd
												ON pp.id = ppd.pembayaran_piutang_id
											WHERE p.pengembalian_id IS NULL AND DATE(pp.tanggal_pembayaran) = {{.QParamDate}} {{.QWhereSalesman}}
											GROUP BY p.user_id, p.id, pp.id
											HAVING p.nominal <> pp.total_pembayaran OR p.nominal <> SUM(ppd.nominal)

											UNION ALL

											SELECT p.user_id as subject_id, 'user' as subject_name, p.id||'' as id, 'stok minus (sales)' as description
											FROM stok_salesman p
											WHERE DATE(p.tanggal_stok) = {{.QParamDate}} AND p.stok_akhir < 0 {{.QWhereSalesman}}

											/*UNION ALL

											SELECT p.user_id as subject_id, 'user' as subject_name, p.id||'' as id, 'stok minus (md)' as description
											FROM md.stok_merchandiser p
											WHERE DATE(p.tanggal_stok) = {{.QParamDate}} AND p.stok_akhir < 0 {{.QWhereSalesman}}

											UNION ALL

											SELECT p.gudang_id as subject_id, 'gudang' as subject_name, p.id||'' as id, 'stok minus (stok gudang laporan)' as description
											FROM stok_gudang_laporan p
											WHERE DATE(p.tanggal_laporan) = {{.QParamDate}} AND p.stok_akhir < 0

											UNION ALL

											SELECT p.gudang_id as subject_id, 'gudang' as subject_name, p.id||'' as id, 'stok minus (stok item laporan)' as description
											FROM md.stok_gudang_item_laporan p
											WHERE DATE(p.tanggal_laporan) = {{.QParamDate}} AND p.stok_akhir < 0

											UNION ALL

											SELECT p.gudang_id as subject_id, 'gudang' as subject_name, p.id||'' as id, 'stok minus (stok gudang)' as description
											FROM stok_gudang p
											WHERE p.jumlah < 0

											UNION ALL

											SELECT p.gudang_id as subject_id, 'gudang' as subject_name, p.id||'' as id, 'stok minus (stok gudang item)' as description
											FROM md.stok_gudang_item p
											WHERE p.jumlah < 0*/

											UNION ALL

											SELECT p.user_id as subject_id, 'user' as subject_name, p.id||'' as id, 'penjualan is kredit tanpa piutang' as description
											FROM penjualan p
											LEFT JOIN piutang pi
												ON p.id = pi.penjualan_id
											WHERE DATE(p.tanggal_penjualan) = {{.QParamDate}} AND pi.id IS NULL AND p.is_kredit = 1 {{.QWhereSalesman}}

											UNION ALL

											SELECT pj.user_id as subject_id, 'user' as subject_name, p.id||'' as id, 'piutang tanpa penjualan' as description
											FROM piutang p
											LEFT JOIN penjualan pj
												ON p.penjualan_id = pj.id
											WHERE DATE(p.tanggal_piutang) = {{.QParamDate}} AND pj.id IS NULL {{.QWhereSalesman}}

											UNION ALL

											SELECT p.user_id as subject_id, 'user' as subject_name, p.id||'' as id, 'pembayaran piutang detail tanpa parent' as description
											FROM pembayaran_piutang_detail p
											LEFT JOIN pembayaran_piutang pp
												ON p.pembayaran_piutang_id = pp.id
											WHERE DATE(p.dtm_crt) = {{.QParamDate}} AND pp.id IS NULL {{.QWhereSalesman}}

											UNION ALL

											SELECT p.user_id as subject_id, 'user' as subject_name, p.id||'' as id, 'pembayaran piutang tanpa detail' as description
											FROM pembayaran_piutang p
											LEFT JOIN pembayaran_piutang_detail ppd
												ON ppd.pembayaran_piutang_id = p.id
											WHERE DATE(p.tanggal_pembayaran) = {{.QParamDate}} AND ppd.id IS NULL {{.QWhereSalesman}}

											UNION ALL

											SELECT p.user_id as subject_id, 'user' as subject_name, p.id||'' as id, 'pembayaran piutang detail totalnya tidak sama dengan parent' as description
											FROM pembayaran_piutang p
											LEFT JOIN pembayaran_piutang_detail ppd
												ON ppd.pembayaran_piutang_id = p.id
											WHERE DATE(p.tanggal_pembayaran) = {{.QParamDate}} {{.QWhereSalesman}}
											GROUP BY p.id
											HAVING p.total_pembayaran <> SUM(ppd.nominal)

											UNION ALL

											SELECT p.user_id as subject_id, 'user' as subject_name, p.id||'' as id, 'total pembayaran lebih besar dari sisa piutang' as description
											FROM piutang p
											LEFT JOIN pembayaran_piutang_detail ppd
												ON ppd.piutang_id = p.id
											WHERE DATE(p.tanggal_piutang) = {{.QParamDate}} {{.QWhereSalesman}}
											GROUP BY p.id
											HAVING SUM(ppd.nominal) > p.total_piutang`

	templateReplaceQuery := map[string]interface{}{
		"QParamDate":     paramDate,
		"QWhereSalesman": qWhereSalesman,
	}

	queryFix, err := helpers.PrepareQuery(query, templateReplaceQuery)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data",
			Success: false,
		})
	}

	results, err := helpers.ExecuteQuery(queryFix)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data",
			Success: false,
		})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Tidak ada data incomplete",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Berhasil mendapatkan data",
		Success: true,
		Data:    results,
	})

}

func GetPengiriman(c *fiber.Ctx) error {

	type Input struct {
		GudangId   int    `json:"gudangId"`
		DateStart  string `json:"dateStart"`
		DateEnd    string `json:"dateEnd"`
		SuratJalan string `json:"suratJalan"`
		Page       string `json:"page"`
		PageSize   string `json:"pageSize"`
	}

	input := Input{}

	if err := c.QueryParser(&input); err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengambil input data",
			Success: false,
		})
	}

	qWhere := ""

	if input.SuratJalan == "" {
		qWhere = " AND sgp.surat_jalan = '" + input.SuratJalan + "'"
	}

	if input.GudangId != 0 {
		qWhere = fmt.Sprintf(" AND sgp.gudang_id = %d", input.GudangId)
	}

	if input.DateStart != "" && input.DateEnd != "" {
		qWhere = " AND sgp.tanggal_kirim BETWEEN DATE('" + input.DateStart + "') AND DATE('" + input.DateEnd + "')"
	}

	var qLimit, qPage string
	iPage, _ := strconv.Atoi(input.Page)
	iPageSize, _ := strconv.Atoi(input.PageSize)

	if input.PageSize != "" {
		qLimit = " LIMIT " + input.PageSize
	} else {
		qLimit = " LIMIT 20"
		iPageSize = 20
	}

	if input.Page == "" {
		iPage = 0
	} else {
		iPage = iPage - 1
	}

	tempQ := strconv.Itoa(iPage * iPageSize)
	qPage = " OFFSET " + tempQ

	queries := fmt.Sprintf(`WITH detail_riwayat as (
									SELECT
											sgp.id,
											p.id as produk_id,
											sgpd.harga_beli,
											p.code,
											SUM(sgpd.jumlah) as qty,
											sgpd.pita,
											ARRAY_AGG(sgpd.*) AS stok_gudang_riwayats
									FROM stok_gudang_pengiriman sgp
									JOIN stok_gudang_riwayat sgpd
										ON sgp.id = sgpd.stok_gudang_pengiriman_id
									JOIN produk p
										ON sgpd.produk_id = p.id
									WHERE sgpd.aksi = 'KIRIM' AND sgp.id = 6075
									GROUP BY sgp.id, p.id, sgpd.harga_beli, sgpd.pita
							)

							SELECT
								JSONB_BUILD_OBJECT(
									'send', DATE(sgp.tanggal_kirim),
									'receive', DATE(sgp.tanggal_terima)
								) as date,
								JSONB_BUILD_OBJECT(
									'system', sgp.surat_jalan,
									'vendor', sgp.surat_jalan_vendor
								) as ref_number,
								g.name as sender,
								gr.name as receiver,
								sgp.tag as tag,
								JSONB_AGG(
									JSONB_BUILD_OBJECT(
										'id', dr.produk_id,
										'code', dr.code,
										'price', dr.harga_beli,
										'qty', dr.qty,
										'pita', dr.pita,
										'detail', TO_JSON(dr.stok_gudang_riwayats)
									)
								) as product,
								sgp.update_price_at as latest_update,
								'asd' as status
							FROM stok_gudang_pengiriman sgp
							JOIN gudang g
								ON sgp.gudang_id_asal = g.id
							JOIN gudang gr
								ON sgp.gudang_id_tujuan = gr.id
							JOIN detail_riwayat dr
								ON sgp.id = dr.id
							JOIN branch b
								ON g.branch_id = b.id
							WHERE TRUE %s
							GROUP BY sgp.id, g.id, gr.id
							%s %s`, qWhere, qLimit, qPage)

	result, err := helpers.ExecuteQuery(queries)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengambil data",
			Success: false,
		})
	}

	if len(result) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data tidak ditemukan",
			Success: false,
		})
	}

	totalPage, _ := helpers.ExecuteQuery(queries)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message":   "Data berhasil diambil",
		"Success":   true,
		"Data":      result[0],
		"TotalPage": len(totalPage),
	})
}

func GetFilterUser(c *fiber.Ctx) error {

	type Input struct {
		Date     string `json:"date"`
		BranchId int    `json:"branchId"`
	}

	input := Input{}

	if err := c.QueryParser(&input); err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengambil input data",
			Success: false,
		})
	}

	templateQuery := fmt.Sprintf(`SELECT u.id, u.full_name, CASE WHEN su.id IS NOT NULL THEN 1 ELSE 0 END as is_active
							FROM user_log_branch ulb
							JOIN public.user u
								ON ulb.user_id = u.id
							LEFT JOIN gudang g
								ON g.branch_id = ulb.branch_id
							LEFT JOIN stok_user su
								ON u.id = su.user_id
								AND g.id = su.gudang_id
								AND DATE(su.tanggal_stok) = DATE('{{.QDate}}')
							WHERE DATE('{{.QDate}}') BETWEEN ulb.start_date AND COALESCE(ulb.end_date, CURRENT_DATE) 
								AND ulb.branch_id = {{.QBranchID}}
							ORDER BY 3 DESC, u.full_name`)

	templateParamQuery := map[string]interface{}{
		"QDate":     input.Date,
		"QBranchID": input.BranchId,
	}

	query1, err := helpers.PrepareQuery(templateQuery, templateParamQuery)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengambil data",
			Success: false,
		})
	}

	templateResult := []struct {
		ID       int    `json:"id"`
		FullName string `json:"full_name"`
		IsActive bool   `json:"is_active"`
	}{}

	if err := db.DB.Exec(query1).Scan(&templateResult).Error; err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengambil data",
			Success: false,
		})
	}

	if len(templateResult) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data tidak ditemukan",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data berhasil diambil",
		Success: true,
		Data:    templateResult,
	})
}
