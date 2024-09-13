package controllers

import (
	"fmt"
	db "pluto_remastered/config"
	"pluto_remastered/helpers"
	"pluto_remastered/structs"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	// orderedmap "github.com/wk8/go-ordered-map/v2"
	newOrderedmap "github.com/iancoleman/orderedmap"
)

func GetSalesmanDailySales(c *fiber.Ctx) error {

	date := c.Query("date")
	salesmanId := c.Query("salesmanId")
	isMd, _ := strconv.Atoi(c.Query("isMd", "0"))
	branchId := helpers.ParamArray(c.Context().QueryArgs().PeekMulti("branchId[]"))
	// dateNewRule := "2021-08-01"

	queries := []string{}

	masterBranchIdDemo := []int16{151, 152, 153, 154, 156, 157}
	mappingBranchIdDemo := map[int]int{151: 1, 152: 2, 153: 3, 154: 4, 156: 6, 157: 7}

	// var qWhereBranchId string
	// if len(branchId) > 0 {
	// 	// branchId = " AND b.id IN (" + strings.Join(branchId, ",") + ")"
	// 	qWhereBranchId = " AND branch_id IN (" + strings.Join(branchId, ",") + ")"
	// }

	var tempBranch int16
	var merchandiserId, tempSalesmanId int32
	// var selectFrom string

	if isMd == 0 {

		salesman := new(structs.Salesman)
		err := db.DB.Where("id = ?", salesmanId).First(&salesman).Error

		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika mengambil data salesmen",
				Success: false,
			})
		}

		tempSalesmanId = salesman.ID
		tempBranch = salesman.BranchID
		// selectFrom = " salesman"
		merchandiserId = 0
	} else {
		id, _ := strconv.Atoi(salesmanId)
		merchandiserId = int32(id)

		merchandiser := structs.Merchandiser{}
		err := db.DB.Where("id = ?", merchandiserId).First(&merchandiser).Error
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika mengambil data merchandiser",
				Success: false,
			})
		}

		tempBranch = merchandiser.BranchID
		// selectFrom = " md.merchandiser"
		tempSalesmanId = 0
	}

	var queryProduk string
	if slices.Contains(masterBranchIdDemo, tempBranch) {
		srId := mappingBranchIdDemo[int(tempBranch)]
		queryProduk = fmt.Sprintf(`SELECT 
						DISTINCT ON(pb.produk_id)
						p.id, LOWER(p.name) AS name, LOWER(p.code) AS code 
						FROM produk_branch pb
						JOIN branch b 
						ON b.id = pb.branch_id
						JOIN rayon r
						ON r.id = b.rayon_id
						JOIN produk p
						ON p.id = pb.produk_id
						WHERE r.sr_id = %v
						ORDER BY pb.produk_id, p.order ASC`, srId)

	} else {
		queryProduk = fmt.Sprintf(`SELECT produk.id, LOWER(name) AS name, LOWER(code) AS code 
								FROM produk_branch 
								JOIN produk ON produk.id = produk_branch.produk_id 
								WHERE is_aktif = 1 AND branch_id IN (` + strings.Join(branchId, ",") + `) 
								GROUP BY produk.id
								ORDER BY produk.order ASC`)
	}

	// fmt.Println(queryProduk)

	produks, err := helpers.ExecuteQuery2(queryProduk, "")
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika mengambil data produk",
			Success: false,
		})
	}

	var qProd string
	tempArray := make(map[string]int)
	var produkHeader []string
	for _, result := range produks {
		tempProdId, _ := result.Get("id")
		tempProdCode, _ := result.Get("code")
		qProd = qProd + fmt.Sprintf(`,COALESCE(SUM(case when (data.tunai > 0 OR data.kredit > 0) AND data.produk_id = %v then data.jumlah else 0 end),0) %v`, tempProdId, tempProdCode)
		tempArray[tempProdCode.(string)] = 0
		produkHeader = append(produkHeader, tempProdCode.(string))
	}

	qProd = qProd[1:]

	const sqlTemplate1 = `SELECT 
							DATE(tanggal) AS tanggal,

							json_build_object('id', COALESCE(c.id||'', outlet.id||'', '-1'), 'name', INITCAP(COALESCE(c.name, outlet.name, 'End User')), 'outlet_name', COALESCE(c.outlet_name, outlet.outlet_name, 'End User'), 'type', MAX(ct.name)) AS customer,
							INITCAP(COALESCE(c.kelurahan, outlet.kelurahan)) as kelurahan,
							(select json_agg(a ORDER BY a ASC) from unnest(array_agg(DISTINCT CASE WHEN no_nota != '' THEN no_nota END))a WHERE a IS NOT NULL) AS invoice,
							{{.QProd}},
							-- (select json_agg(a) from unnest(array_agg( case when (data.retur > 0) AND data.produk_id > 0 THEN
							--   (json_build_object('kode', prd.code,'jumlah', data.jumlah)) END))a where a is not null) AS produk_retur,
							-- (select json_agg(a) from unnest(array_agg( case when ( data.tunai = 0 AND data.kredit = 0 AND data.retur = 0 AND pembayaran = 0) AND data.produk_id > 0 THEN
							--   (json_build_object('kode', prd.code,'jumlah', data.jumlah)) END))a where a is not null) AS bonus,
							(select json_agg(a) from unnest(array_agg( case when (data.pengambalian_id > 0) AND data.produk_id > 0 THEN
							(json_build_object('kode', prd.code,'jumlah', data.jumlah)) END))a where a is not null) AS produk_retur,
							-- (select json_agg(a) from unnest(array_agg( case when (data.retur > 0) AND data.produk_id > 0 THEN
							-- (json_build_object('kode', prd.code,'jumlah', data.jumlah)) END))a where a is not null) AS produk_retur,

							(select json_agg(a) from unnest(array_agg( case when ( data.tunai = 0 AND data.kredit = 0 AND data.retur = 0 AND pembayaran = 0 AND data.penjualan_id <> 9999999) AND data.produk_id > 0 THEN
							(json_build_object('kode', prd.code,'jumlah', data.jumlah)) END))a where a is not null) AS bonus,
							SUM(data.kredit) AS kredit,
							SUM(data.tunai) AS tunai,
								SUM(pembayaran) AS pembayaran,
							SUM(dp) AS dp, 
							SUM(data.retur) AS retur, 
							SUM(adjustment) AS adjustment,
								SUM(data.tunai) + SUM(pembayaran) + SUM(dp) - SUM(retur) + SUM(adjustment) as setoran,
								JSONB_BUILD_OBJECT(
									'checkin', to_char(MIN(kl.checkin_at) , 'HH24:MI'),
									'checkout', to_char(MAX(kl.checkout_at), 'HH24:MI'),
									'durasi', MIN(kl.checkout_at) - MAX(kl.checkin_at),
							'accuracy_in', MAX(kl.accuracy_in),
							'accuracy_out', MAX(kl.accuracy_out),
							'distance_in', MAX(kl.distance_in),
							'distance_out', MAX(kl.distance_out),
							-- 'trans_tunai', COUNT(data.penjualan_id) FILTER (WHERE data.penjualan_id <> 9999999 AND data.penjualan_id <> 0),
									-- 'trans_kredit', COUNT(data.penjualan_id) FILTER (WHERE data.penjualan_id <> 9999999 AND data.penjualan_id <> 0),
									-- 'trans_retur', COUNT(data.pengambalian_id) FILTER (WHERE data.pengambalian_id <> 0),
									-- 'trans_payment', COUNT(data.pembayaran_id) FILTER (WHERE data.pembayaran_id <> 0),
									-- 'qr', COUNT(qr.id),
									-- 'kunjungan_log', COUNT(DISTINCT kl.id),
							'jumlah_aktifitas', COUNT(DISTINCT kl.id) FILTER (WHERE kl.checkin_at IS NOT NULL) + 
												COUNT(DISTINCT kl.id) FILTER (WHERE kl.checkout_at IS NOT NULL) + 
												COUNT(DISTINCT data.penjualan_id) FILTER (WHERE data.penjualan_id <> 9999999 AND data.penjualan_id <> 0) + 
												-- COUNT(DISTINCT data.penjualan_id) FILTER (WHERE data.penjualan_id <> 9999999 AND data.penjualan_id <> 0) + 
												-- COUNT(DISTINCT data.penjualan_id) FILTER (WHERE data.penjualan_id <> 9999999 AND data.penjualan_id <> 0) + 
												COUNT(DISTINCT data.pengambalian_id) FILTER (WHERE data.pengambalian_id <> 0) + 
																			COUNT(DISTINCT data.pembayaran_id) FILTER (WHERE data.pembayaran_id <> 0) + 
												COUNT(DISTINCT qr.id)

								) as informasi

							FROM
							(SELECT MAX(p.id) AS penjualan_id, 0 AS pengambalian_id, 0 AS pembayaran_id, 0 AS payment_id, 0 AS kunjungan_id, p.customer_id, pd.produk_id, pd.condition, pd.pita, SUM (pd.jumlah) AS jumlah, 
							MAX(pd.harga) AS harga, SUM(CASE WHEN p.is_kredit = 0 THEN (pd.harga-pd.diskon)*pd.jumlah ELSE 0 END) AS tunai, SUM(CASE WHEN p.is_kredit = 1 THEN (pd.harga-pd.diskon)*pd.jumlah ELSE 0 END) AS kredit, 
							0 AS retur, 0 AS pembayaran, 0 AS adjustment, 0 AS dp, p.tanggal_penjualan as tanggal, p.tanggal_penjualan::varchar AS waktu_penjualan, null AS waktu_pengembalian, null AS waktu_pembayaran, null AS waktu_kunjungan, 
							null AS waktu_payment, MAX(p.no_nota) as no_nota, MAX(p.image_nota) AS image_nota, MAX(pd.diskon) AS diskon, COALESCE(MAX(p.salesman_id), MAX(p.merchandiser_id)) as salesman_id, null AS payment_tipe, null AS payment_image, 
							0 AS payment_nominal, null AS pembayaran_tgl_piutang, -1 AS is_lunas, p.outlet_id
							FROM penjualan p
							LEFT JOIN penjualan_detail pd ON p.id = pd.penjualan_id
							WHERE DATE(p.tanggal_penjualan) = DATE('{{.QDate}}') AND (p.salesman_id = {{.QSalesmanId}} OR p.merchandiser_id = {{.QMerchandiserId}})
							GROUP BY  pd.id, pd.condition, pd.pita, p.id, pd.harga

							UNION ALL

							SELECT 9999999 AS penjualan_id, MAX(pg.id) AS pengambalian_id, 0 AS pembayaran_id, 0 AS payment_id, 0 AS kunjungan_id, MAX(p.customer_id) AS customer_id, pgd.produk_id, pgd.condition, pgd.pita, SUM(pgd.jumlah) AS jumlah, 0 AS harga, 0 AS tunai,
							0 AS kredit, SUM(pgd.harga*pgd.jumlah) AS retur, 0 AS pembayaran, 0 AS adjustment, 0 AS dp, pg.tanggal_pengembalian as tanggal, null AS waktu_penjualan, pg.tanggal_pengembalian::varchar AS waktu_pengembalian, null AS waktu_pembayaran, null AS waktu_kunjungan, null AS waktu_payment, '' AS no_nota 
							, MAX(pg.image_nota) AS image_nota, 0 AS diskon, COALESCE(MAX(pg.salesman_id), MAX(pg.merchandiser_id)) as salesman_id, null AS payment_tipe, null AS payment_image, 0 AS payment_nominal, null AS pembayaran_tgl_piutang, -1 AS is_lunas, null as outlet_id
							FROM pengembalian pg
							LEFT JOIN penjualan p ON p.id = pg.penjualan_id
							LEFT JOIN pengembalian_detail pgd ON pg.id = pgd.pengembalian_id
							WHERE DATE(pg.tanggal_pengembalian) = DATE('{{.QDate}}') AND (pg.salesman_id = {{.QSalesmanId}} OR pg.merchandiser_id = {{.QMerchandiserId}})
							GROUP BY pgd.id, pgd.condition, pgd.pita, pg.id, pgd.id

							UNION ALL

							SELECT 9999999 AS penjualan_id, 0 AS pengambalian_id, pp.id AS pembayaran_id, py.id AS payment_id, 0 AS kunjungan_id, pp.customer_id , 0 AS produk_id, '' AS condition, '' AS pita, 0 AS jumlah, 0 AS harga, 0 as tunai, 0 as kredit, 0 as retur, 
							(CASE WHEN pp.tipe_pelunasan = 0 AND pp.pengembalian_id IS NULL then pp.total_pembayaran ELSE 0 END) AS pembayaran,
							(CASE WHEN pp.tipe_pelunasan = 0 AND pp.pengembalian_id IS NOT NULL then pp.total_pembayaran ELSE 0 END) AS adjustment,
							(CASE WHEN tipe_pelunasan = 1 then ppd.nominal ELSE 0 END) as dp, tanggal_pembayaran as tanggal , null AS waktu_penjualan, null AS waktu_pengembalian, tanggal_pembayaran::varchar AS waktu_pembayaran, 
							null AS waktu_kunjungan, null AS waktu_payment, '' AS no_nota,image_nota AS image_nota, 0 AS diskon, COALESCE(pp.salesman_id, pp.merchandiser_id) as salesman_id, py.tipe AS payment_tipe, null AS payment_image,  COALESCE(py.nominal,0) AS payment_nominal, p.tanggal_piutang::varchar AS pembayaran_tgl_piutang, p.is_lunas AS is_lunas,
							null as outlet_id
							FROM pembayaran_piutang pp
							JOIN pembayaran_piutang_detail ppd ON ppd.pembayaran_piutang_id = pp.id
							LEFT JOIN piutang p ON p.id = ppd.piutang_id
							LEFT JOIN payment py ON py.id =  pp.payment_id AND py.tipe = 'BILYET GIRO'
							WHERE DATE(tanggal_pembayaran) = DATE('{{.QDate}}') AND (pp.salesman_id = {{.QSalesmanId}} OR pp.merchandiser_id = {{.QMerchandiserId}})

							UNION ALL

							SELECT 9999999 AS penjualan_id, 0 AS pengambalian_id, 0 AS pembayaran_id, 0 AS payment_id, k.id AS kunjungan_id, customer_id, 0 AS produk_id, '' AS condition, '' AS pita, 0 AS jumlah, 0 AS harga, 0 AS tunai, 0 AS kredit, 0 AS retur, 0 AS pembayaran, 0 AS adjustment, 0 AS dp, 
							tanggal_kunjungan AS tanggal, null AS waktu_penjualan, null AS waktu_pengembalian, null AS waktu_pembayaran, tanggal_kunjungan::varchar AS waktu_kunjungan, null AS waktu_payment, '' AS no_nota, null AS image_nota, 0 AS diskon, COALESCE(k.salesman_id, k.merchandiser_id) as salesman_id, 
							null AS payment_tipe, null AS payment_image,
							0 AS payment_nominal, null AS pembayaran_tgl_piutang, -1 AS is_lunas, k.outlet_id as outlet_id
							FROM kunjungan k
							WHERE DATE(tanggal_kunjungan) = DATE('{{.QDate}}') AND (k.salesman_id = {{.QSalesmanId}} OR k.merchandiser_id = {{.QMerchandiserId}})     

							UNION ALL

							SELECT 9999999 AS penjualan_id, 0 AS pengambalian_id, 0 AS pembayaran_id, (CASE WHEN tipe != 'CASH' THEN id ELSE 0 END ) AS payment_id, 0 AS kunjungan_id, customer_id, 0 AS produk_id, '' AS condition, '' AS pita, 0 AS jumlah, 0 AS harga, 0 AS tunai, 0 AS kredit, 0 AS retur, 0 AS pembayaran, 0 AS adjustment, 0 AS dp,
							tanggal_transaksi AS tanggal, null AS waktu_penjualan, null AS waktu_pengembalian, null AS waktu_pembayaran, null AS waktu_kunjungan, tanggal_transaksi::varchar AS waktu_payment,'' AS no_nota, null AS image_nota, 0 AS diskon, COALESCE(p.salesman_id, p.merchandiser_id) as salesman_id,  
							(CASE WHEN tipe = 'BILYET GIRO' THEN 'OPEN ' ELSE '' END)||tipe AS payment_tipe, bukti_bayar AS payment_image,
							nominal AS payment_nominal, null AS pembayaran_tgl_piutang, -1 AS is_lunas, null as outlet_id
							FROM payment p
							WHERE (p.salesman_id = {{.QSalesmanId}} OR p.merchandiser_id = {{.QMerchandiserId}}) AND (DATE(tanggal_transaksi) = DATE('{{.QDate}}') )
							) data
							LEFT JOIN produk prd ON prd.id = data.produk_id
							LEFT JOIN customer c ON data.customer_id = c.id
							LEFT JOIN md.outlet outlet ON outlet.id = data.outlet_id AND data.customer_id IS NULL
							LEFT JOIN customer_type ct ON ct.id = CASE WHEN data.customer_id IS NULL THEN outlet.type ELSE c.tipe END
							LEFT JOIN verifikasi_setoran vs ON date(data.tanggal) = date(vs.tanggal_setoran) AND data.customer_id = vs.customer_id AND data.salesman_id = vs.salesman_id
							LEFT JOIN public.user u ON vs.user_id = u.id
							LEFT JOIN stok_salesman ss ON ss.salesman_id = data.salesman_id AND DATE(ss.tanggal_stok) = DATE(data.tanggal) AND ss.produk_id = data.produk_id AND data.condition = ss.condition AND data.pita = ss.pita
							LEFT JOIN (SELECT DISTINCT ON (customer_id) * 
										FROM kunjungan_log WHERE (salesman_id = {{.QSalesmanId}} OR merchandiser_id = {{.QMerchandiserId}}) AND date(checkin_at) = DATE('{{.QDate}}') 
										ORDER BY customer_id, checkin_at) kl ON c.id = kl.customer_id AND DATE(kl.checkin_at) = DATE('{{.QDate}}')
							LEFT JOIN qr_code_history qr ON c.id = qr.customer_id AND qr.datetime = DATE('{{.QDate}}')
							GROUP BY c.id, date(data.tanggal), outlet.id--, kl.id
							ORDER BY MIN(data.tanggal), MIN(penjualan_id) ASC`

	templateQuery := map[string]interface{}{
		"QProd":           qProd,
		"QDate":           date,
		"QSalesmanId":     tempSalesmanId,
		"QMerchandiserId": merchandiserId,
	}

	query1, err := helpers.PrepareQuery(sqlTemplate1, templateQuery)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	const queryCheckSO = `SELECT id FROM stok_salesman ss WHERE ss.tanggal_stok::DATE = DATE('{{.QDate}}') AND is_complete = 0 AND (salesman_id = {{.QSalesmanId}} OR merchandiser_id = {{.QMerchandiserId}})`

	finalQueryCheckSO, err := helpers.PrepareQuery(queryCheckSO, templateQuery)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	var queryGetReq string
	if isMd == 0 {
		queryGetReq = `SELECT to_char(DATE('{{.QDate}}'), 'DD Mon YYYY') as tanggal,
							JSONB_BUILD_OBJECT(
							'name', INITCAP(s.name),
							'type', INITCAP(st.name)
							) as salesman 
						FROM salesman s
						LEFT JOIN salesman_type2 st 
						ON s.salesman_type_id = st.id
						WHERE s.id = {{.QSalesmanId}}`
	} else {
		queryGetReq = `SELECT to_char(DATE('{{.QDate}}'), 'DD Mon YYYY') as tanggal,
										  JSONB_BUILD_OBJECT(
											'name', INITCAP(s.name),
											'type', INITCAP(st.name)
										  ) as salesman 
									  FROM md.merchandiser s
									  LEFT JOIN salesman_type2 st 
										ON 5 = st.id
									  WHERE s.id = {{.QMerchandiserId}}`
	}

	finalQueryGetReq, err := helpers.PrepareQuery(queryGetReq, templateQuery)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	queries = append(queries, query1)
	queries = append(queries, finalQueryCheckSO)
	queries = append(queries, finalQueryGetReq)

	var wg2 sync.WaitGroup
	resultsChan2 := make(chan map[int][]*newOrderedmap.OrderedMap, len(queries))

	tempResults := make([][]*newOrderedmap.OrderedMap, len(queries))

	// Launch concurrent Goroutines
	for i, query := range queries {
		wg2.Add(1)
		go helpers.ExecuteGORMQueryOrdered(query, resultsChan2, i, &wg2)
	}

	// Wait for all Goroutines to finish
	wg2.Wait()
	close(resultsChan2)

	for result := range resultsChan2 {
		for index, res := range result {
			tempResults[index] = res
		}
	}

	rowTable := tempResults[0]
	resCheckSO := false
	//result so
	if len(tempResults[1]) == 0 {
		resCheckSO = true
	}

	var returnTanggal string
	tempReqSalesman := newOrderedmap.New()
	for i := 0; i < len(tempResults[2]); i++ {
		for _, k := range tempResults[2][i].Keys() {
			val, _ := tempResults[2][i].Get(k)
			if k == "tanggal" {
				returnTanggal = val.(string)
			}

			if k == "salesman" {
				tempReqSalesman.Set("salesman", val)
			}
		}
	}

	returnSalesman, _ := tempReqSalesman.Get("salesman")

	tempArrayProdukHeader := make(map[string]interface{}, len(produkHeader))
	for i := 0; i < len(produkHeader); i++ {
		tempArrayProdukHeader[produkHeader[i]] = 0
		for j := 0; j < len(rowTable); j++ {
			for _, k := range rowTable[j].Keys() {
				value, _ := rowTable[j].Get(k)
				if value, ok := value.(float64); ok {
					if k == produkHeader[i] {
						tempArrayProdukHeader[produkHeader[i]] = tempArrayProdukHeader[produkHeader[i]].(int) + int(value)
					}
				}
			}
		}
	}

	tempArrayProdukHeader, remainingKeysHeader, removedKeysHeader := helpers.RemoveZeroFromMap(tempArrayProdukHeader)

	indexToRemove := []string{"produk_retur", "bonus"}
	headerUang := newOrderedmap.New()

	headerUang.Set("tunai", 0)
	headerUang.Set("kredit", 0)
	headerUang.Set("retur", 0)
	headerUang.Set("pembayaran", 0)
	headerUang.Set("dp", 0)
	headerUang.Set("adjustment", 0)
	headerUang.Set("setoran", 0)

	sumUang := make(map[string]interface{}, len(headerUang.Keys()))

	for i := 0; i < len(rowTable); i++ {
		for _, val := range headerUang.Keys() {
			value, _ := rowTable[i].Get(val)

			if sumUang[val] != nil {
				sumUang[val] = sumUang[val].(int) + int(value.(float64))
			} else {
				sumUang[val] = int(value.(float64))
			}
		}
	}

	sumUang, remainingKeysUang, removedKeysUang := helpers.RemoveZeroFromMap(sumUang)

	for i := 0; i < len(rowTable); i++ {
		for j := 0; j < len(indexToRemove); j++ {
			rowTable[i].Delete(indexToRemove[j])
		}

		for j := 0; j < len(removedKeysHeader); j++ {
			rowTable[i].Delete(removedKeysHeader[j])
		}

		for j := 0; j < len(removedKeysUang); j++ {
			rowTable[i].Delete(removedKeysUang[j])
		}
	}

	type newResponse struct {
		Message       string                      `json:"message"`
		Success       bool                        `json:"success"`
		IsStokOpname  bool                        `json:"isStokOpname"`
		Date          string                      `json:"date"`
		Salesman      interface{}                 `json:"salesman"`
		Data          []*newOrderedmap.OrderedMap `json:"data"`
		Nominal       []string                    `json:"nominal"`
		Header        []string                    `json:"header"`
		Produk        []string                    `json:"produk"`
		TotalProduk   map[string]interface{}      `json:"total_produk"`
		TotalNominal  map[string]interface{}      `json:"total_nominal"`
		ProdukLength  int                         `json:"produk_length"`
		NominalLength int                         `json:"nominal_length"`
	}

	return c.Status(fiber.StatusOK).JSON(newResponse{
		Message:       "Data has been loaded",
		Success:       true,
		IsStokOpname:  resCheckSO,
		Date:          returnTanggal,
		Salesman:      returnSalesman,
		Data:          rowTable,
		Nominal:       remainingKeysUang,
		Produk:        remainingKeysHeader,
		TotalProduk:   tempArrayProdukHeader,
		TotalNominal:  sumUang,
		ProdukLength:  len(remainingKeysHeader),
		NominalLength: len(remainingKeysUang),
	})

	// return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 	"Data":           rowTable,
	// 	"isStokOpname":   resCheckSO,
	// 	"Salesman":       returnSalesman,
	// 	"date":           returnTanggal,
	// 	"nominal":        sumUang,
	// 	"header":         indexToRemove,
	// 	"produk":         tempArrayProdukHeader,
	// 	"produk_length":  len(tempArrayProdukHeader),
	// 	"nominal_length": len(sumUang),
	// 	"Success":        true,
	// 	"Message":        "Success",
	// })
}
