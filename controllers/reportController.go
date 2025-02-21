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

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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
		tempBranch = *salesman.BranchID
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

		tempBranch = *merchandiser.BranchID
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
}

func GetUserDailyReport(c *fiber.Ctx) error {

	date := c.Query("date")
	userId := c.Query("userId")
	isQuery := c.Query("isQuery")
	requestId := c.Query("requestId")
	branchId := helpers.ParamArray(c.Context().QueryArgs().PeekMulti("branchId[]"))

	if date == "" {
		date = "CURRENT_DATE"
	} else {
		date = "DATE('" + date + "') "
	}

	templateQuery := map[string]interface{}{
		"QUserID": userId,
		"QDate":   date,
	}

	queryProduk, err := helpers.PrepareQuery(`SELECT pr.id, pr.code FROM penjualan p
					JOIN penjualan_detail pd
						ON p.id = pd.penjualan_id
					JOIN produk pr
						ON pd.produk_id = pr.id
					WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = {{.QDate}}
					GROUP BY pr.id`, templateQuery)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	produks, err := helpers.ExecuteQuery2(queryProduk, "")
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika mengambil data produk",
			Success: false,
		})
	}

	var qProd, qProdItem, qSelectProd string
	tempArray := make(map[string]int)
	var produkHeader []string
	var produkIds []int

	if len(produks) > 0 {

		for _, result := range produks {
			tempProdId, _ := result.Get("id")
			tempProdCode, _ := result.Get("code")
			qProd = qProd + fmt.Sprintf(`,COALESCE(SUM(case when (pj.tunai > 0 OR pj.kredit > 0) AND pj.produk_id = %v then pj.jumlah else 0 end),0) %v`, tempProdId, tempProdCode)
			qProdItem = qProdItem + fmt.Sprintf(`,0 %v`, tempProdCode)
			qSelectProd = qSelectProd + fmt.Sprintf(`,MAX(sq.%v) as %v`, tempProdCode, tempProdCode)
			tempArray[tempProdCode.(string)] = 0
			produkHeader = append(produkHeader, tempProdCode.(string))
			produkIds = append(produkIds, int(tempProdId.(float64)))
		}

		qProd = qProd[1:]
		qProd = qProd + ","
		qProdItem = qProdItem[1:]
		qProdItem = qProdItem + ","
		qSelectProd = qSelectProd[1:]
		qSelectProd = qSelectProd + ","
	}

	templateEndQuery := map[string]interface{}{
		"QUserID":     userId,
		"QDate":       date,
		"QSelectProd": qSelectProd,
		"QProd":       qProd,
		"QProdItem":   qProdItem,
		"QAsUserID":   userId + " as user_id",
		"QBranchID":   strings.Join(branchId, ","),
	}

	if requestId != "" {

		endQuery := `WITH visits as (
						SELECT id, tanggal_kunjungan, customer_id, user_id, subject_type_id, image_kunjungan
						FROM kunjungan 
						WHERE user_id = {{.QUserID}} 
							AND DATE(tanggal_kunjungan) = {{.QDate}}
						ORDER BY tanggal_kunjungan
					), penjualan_cus as (
						SELECT p.no_nota,
										p.id,
										CASE WHEN COALESCE(p.image_nota_print, p.image_nota) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/nota/'||COALESCE(p.image_nota_print, p.image_nota) ELSE NULL END as image_nota,
										CASE WHEN p.image_bukti_serah IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/nota/'||p.image_bukti_serah ELSE NULL END as image_bukti_serah,
										p.user_id,
										p.customer_id,
										DATE(p.tanggal_penjualan),
										COALESCE(SUM(CASE WHEN p.is_kredit = 0 THEN ((pd.harga - pd.diskon) * pd.jumlah) ELSE 0 END),0) as tunai,
										COALESCE(SUM(CASE WHEN p.is_kredit = 1 THEN ((pd.harga - pd.diskon) * pd.jumlah) ELSE 0 END),0) as kredit,
										pd.produk_id as produk_id,
										COALESCE(SUM(pd.jumlah),0) as jumlah
						FROM penjualan p
						JOIN penjualan_detail pd
							ON p.id = pd.penjualan_id
						WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = {{.QDate}}
						GROUP BY p.id, pd.produk_id
					), pembayaran_cus as (
						SELECT p.customer_id,
										p.user_id,
										COALESCE(SUM(pd.nominal),0) as pembayaran,
										COALESCE((CASE WHEN p.tipe_pelunasan = 0 AND p.pengembalian_id IS NOT NULL then p.total_pembayaran ELSE 0 END),0) AS adjustment,
										COALESCE(SUM((CASE WHEN tipe_pelunasan = 1 then pd.nominal ELSE 0 END)),0) as dp
						FROM pembayaran_piutang p
						JOIN pembayaran_piutang_detail pd
							ON p.id = pd.pembayaran_piutang_id
						LEFT JOIN payment py ON py.id =  p.payment_id AND py.tipe = 'BILYET GIRO'
						WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pembayaran) = {{.QDate}}
						GROUP BY p.id
					), pengembalian_cus as (
						SELECT p.customer_id,
										p.user_id,
										COALESCE(SUM(pd.jumlah * pd.harga),0) as pengembalian
						FROM pengembalian p
						JOIN pengembalian_detail pd
							ON p.id = pd.pengembalian_id
					WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pengembalian) = {{.QDate}}
						GROUP BY p.id
					), payments as (
						SELECT UPPER((CASE WHEN tipe = 'BILYET GIRO' THEN 'OPEN ' ELSE '' END)||tipe) AS payment_tipe, 
								COALESCE(nominal,0) AS payment_nominal,
										p.user_id,
										p.customer_id,
										CASE WHEN p.bukti_bayar = '' THEN NULL ELSE p.bukti_bayar END as payment_image,
																				p.penjualan_id as penjualan_id
								FROM payment p
								WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_transaksi) = {{.QDate}} 

						UNION

						SELECT 'CASH' as payment_tipe,
								SUM(((pd.harga - pd.diskon) * pd.jumlah)) as payment_nominal,
								p.user_id,
								p.customer_id,
								null as payment_image,
								p.id as penjualan_id
						FROM penjualan p
						LEFT JOIN penjualan_detail pd
								ON p.id = pd.penjualan_id
						LEFT JOIN payment py
								ON p.id = py.penjualan_id
								AND py.id IS NULL 
						WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_penjualan) = {{.QDate}} AND p.is_kredit = 0 
						GROUP BY p.id
					), kunjunganlogs as (
						SELECT DISTINCT ON (customer_id) * 
						FROM kunjungan_log WHERE user_id = {{.QUserID}} AND date(checkin_at) = {{.QDate}}
						ORDER BY customer_id, checkin_at
					), transactions as (
						SELECT t.user_id,
										t.subject_type_id,
										t.customer_id,
											JSONB_BUILD_OBJECT(
												'name', i.name,
												'quantity', COALESCE(SUM(td.qty),0)
										) as items,
										CASE WHEN ic.id = 191 AND t.transaction_type_id <> 1 
											THEN JSONB_BUILD_OBJECT(
												'photo_before', CASE WHEN MAX(t.photos_before) = '{}' OR MAX(t.photos_before) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos_before)) image
																	) END,
												'photo_after', CASE WHEN MAX(t.photos_after) = '{}' OR MAX(t.photos_after) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos_after)) image
																	) END
											)
											ELSE NULL
										END as sampling,
										CASE WHEN ic.id = 171 AND t.transaction_type_id <> 1 
											THEN JSONB_BUILD_OBJECT(
												'photo_before', CASE WHEN MAX(t.photos_before) = '{}' OR MAX(t.photos_before) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos_before)) image
																		) END,
												'photo_after', CASE WHEN MAX(t.photos_after) = '{}' OR MAX(t.photos_after) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos_after)) image
																		) END
											)
											ELSE NULL
										END as posm,
										CASE WHEN ic.id = 161 AND t.transaction_type_id <> 1 
											THEN JSONB_BUILD_OBJECT(
												'photo_before', CASE WHEN MAX(t.photos_before) = '{}' OR MAX(t.photos_before) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos_before)) image
																		) END,
												'photo_after', CASE WHEN MAX(t.photos_after) = '{}' OR MAX(t.photos_after) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos_after)) image
																		) END
											)
											ELSE NULL
										END as merchandise
										
						FROM md.transaction t
						JOIN md.transaction_detail td
							ON t.id = td.transaction_id
						JOIN md.item i
							ON td.item_id = i.id
						JOIN md.item_category ic
							ON i.category_id = ic.id
						WHERE t.user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}}
						GROUP BY t.user_id, t.customer_id, ic.id, t.transaction_type_id, t.subject_type_id, i.id
					)

					SELECT st.name as subject_type, 
							sq.customer, 
							MAX(invoice) as invoice,  
							{{.QSelectProd}}
							JSONB_AGG(sq.items) FILTER (WHERE sq.items IS NOT NULL) as items,
							COALESCE(SUM(sq.tunai),0) as tunai,
							COALESCE(SUM(sq.kredit),0) as kredit,
							COALESCE(SUM(sq.pengembalian), 0) as retur,
							COALESCE(SUM(sq.pembayaran),0) as pembayaran,
							COALESCE(SUM(sq.adjustment),0) as adjustment,
							COALESCE(SUM(sq.tunai),0) + 
							COALESCE(SUM(sq.pembayaran),0) + 
							COALESCE(SUM(sq.dp),0) - 
							COALESCE(SUM(sq.pengembalian),0) + 
							COALESCE(SUM(sq.adjustment),0) as setoran,
							JSONB_AGG(sq.payment_information) FILTER (WHERE sq.payment_information IS NOT NULL)->0 as payment_information,
							JSONB_AGG(sq.time_info) FILTER (WHERE sq.time_info IS NOT NULL)->0 as time_info,
							JSONB_AGG(sq.nota) FILTER (WHERE sq.nota IS NOT NULL)->0 as nota,
							JSONB_AGG(sq.penyerahan_produk) FILTER (WHERE sq.penyerahan_produk IS NOT NULL)->0 as penyerahan_produk,
							JSONB_AGG(sq.posm) FILTER (WHERE sq.posm IS NOT NULL)->0 as posm,
							JSONB_AGG(sq.merchandise) FILTER (WHERE sq.merchandise IS NOT NULL)->0 as merchandise,
							JSONB_AGG(sq.sampling) FILTER (WHERE sq.sampling IS NOT NULL)->0 as sampling
					FROM (
						SELECT v.subject_type_id,
								COALESCE(c.id,v.customer_id, pj.customer_id, kl.customer_id) as customer_id,
								JSONB_BUILD_OBJECT(
										'name', COALESCE(c.name, 'End User'),
										'contact', COALESCE(c.outlet_name, '-'),
										'location', COALESCE(c.kelurahan, '-'),
										'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
								) as customer,
								pj.no_nota as invoice,
								{{.QProd}}
								null::jsonb as items,
								COALESCE(SUM(pj.tunai),0) as tunai,
								COALESCE(SUM(pj.kredit),0) as kredit,
								0 as pembayaran,
								0 as pengembalian,
								0 as adjustment,
								0 as dp,
								NULL as payment_information,
								MIN(kl.checkin_at) as param,
								JSONB_BUILD_OBJECT(
									'checkin', to_char(MIN(kl.checkin_at) , 'HH24:MI'),
									'checkout', to_char(MAX(kl.checkout_at), 'HH24:MI'),
									'durasi', MAX(kl.checkout_at) - MIN(kl.checkin_at)
								) as time_info,
								pj.image_nota as nota,
								pj.image_bukti_serah as penyerahan_produk,
								null::jsonb as posm,
								null::jsonb as merchandise,
								null::jsonb as sampling
							FROM visits v
							FULL JOIN penjualan_cus pj
								ON v.user_id = pj.user_id
								AND v.customer_id = pj.customer_id
							LEFT JOIN kunjunganlogs kl
								ON v.user_id = kl.user_id
								AND v.customer_id = kl.customer_id
							LEFT JOIN customer c
								ON COALESCE(v.customer_id, pj.customer_id, kl.customer_id) = c.id
							GROUP BY v.subject_type_id, v.image_kunjungan, c.id, COALESCE(c.id,v.customer_id, pj.customer_id, kl.customer_id), pj.no_nota, pj.image_nota, pj.image_bukti_serah
							
							UNION

							SELECT v.subject_type_id,
									COALESCE(c.id,v.customer_id, pem.customer_id, pg.customer_id, py.customer_id, kl.customer_id) as customer_id,
									JSONB_BUILD_OBJECT(
											'name', COALESCE(c.name, 'End User'),
											'contact', COALESCE(c.outlet_name, '-'),
											'location', COALESCE(c.kelurahan, '-'),
											'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
									) as customer,
									NULL as invoice,
									{{.QProdItem}}
									null::jsonb as items,
									0 as tunai,
									0 as kredit,
									COALESCE(SUM(pem.pembayaran),0) as pembayaran,
									COALESCE(SUM(pg.pengembalian),0) as pengembalian,
									COALESCE(SUM(pem.adjustment),0) as adjustment,
									COALESCE(SUM(pem.dp),0) as dp,
									JSONB_BUILD_OBJECT(
									'cash', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'CASH'), 0) - COALESCE(SUM(pg.pengembalian),0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'CASH' AND py.payment_image <> '')
										),
									'transfer', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'TRANSFER'), 0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'TRANSFER' AND py.payment_image <> '')
										),
									'bilyet_giro_cair', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'BILYET GIRO'), 0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'BILYET GIRO' AND py.payment_image <> '')
										),
									'bilyet_giro_open', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'OPEN BILYET GIRO'), 0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'OPEN BILYET GIRO' AND py.payment_image <> '')
										)
									) as payment_information,
									MIN(kl.checkin_at) as param,
									JSONB_BUILD_OBJECT(
										'checkin', to_char(MIN(kl.checkin_at) , 'HH24:MI'),
										'checkout', to_char(MAX(kl.checkout_at), 'HH24:MI'),
										'durasi', MAX(kl.checkout_at) - MIN(kl.checkin_at)
									) as time_info,
									null as nota,
									null as penyerahan_produk,
									null::jsonb as posm,
									null::jsonb as merchandise,
									null::jsonb as sampling
							FROM visits v
							FULL JOIN pembayaran_cus pem
								ON v.user_id = pem.user_id
								AND v.customer_id = pem.customer_id
							FULL JOIN pengembalian_cus pg
								ON v.user_id = pg.user_id
								AND v.customer_id = pg.customer_id
							FULL JOIN payments py
								ON v.user_id = py.user_id
								AND v.customer_id = py.customer_id
							LEFT JOIN kunjunganlogs kl
								ON v.user_id = kl.user_id
								AND v.customer_id = kl.customer_id
							LEFT JOIN customer c
								ON COALESCE(v.customer_id, pem.customer_id, pg.customer_id, py.customer_id, kl.customer_id) = c.id
							GROUP BY v.subject_type_id, v.image_kunjungan, c.id, COALESCE(c.id,v.customer_id, pem.customer_id, pg.customer_id, py.customer_id, kl.customer_id)
							
							UNION
							
								SELECT COALESCE(v.subject_type_id, t.subject_type_id) as subject_type_id,
									COALESCE(c.id, v.customer_id, t.customer_id, kl.customer_id) as customer_id,
									JSONB_BUILD_OBJECT(
											'name', COALESCE(c.name, 'End User'),
											'contact', COALESCE(c.outlet_name, '-'),
											'location', COALESCE(c.kelurahan, '-'),
											'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
									) as customer,
									NULL as invoice,
									{{.QProdItem}}
									t.items as items,
									0 as tunai,
									0 as kredit,
									0 as pembayaran,
									0 as pengembalian,
									0 as adjustment,
									0 as dp,
									NULL as payment_information,
									MIN(kl.checkin_at) as param,
									JSONB_BUILD_OBJECT(
										'checkin', to_char(MIN(kl.checkin_at) , 'HH24:MI'),
										'checkout', to_char(MAX(kl.checkout_at), 'HH24:MI'),
										'durasi', MAX(kl.checkout_at) - MIN(kl.checkin_at)
									) as time_info,
									null as nota,
									null as penyerahan_produk,
									t.posm as posm,
									t.merchandise as merchandise,
									t.sampling as sampling
							FROM visits v
							FULL JOIN transactions t
								ON v.user_id = t.user_id
								AND v.customer_id = t.customer_id
							LEFT JOIN kunjunganlogs kl
								ON v.user_id = kl.user_id
								AND v.customer_id = kl.customer_id
							LEFT JOIN customer c
								ON COALESCE(v.customer_id, t.customer_id, kl.customer_id) = c.id
							GROUP BY COALESCE(v.subject_type_id, t.subject_type_id), v.image_kunjungan, c.id, COALESCE(c.id, v.customer_id, t.customer_id, kl.customer_id), t.items, t.posm, t.merchandise, t.sampling
					) sq
					JOIN md.subject_type st
						ON COALESCE(sq.subject_type_id, CASE WHEN sq.customer_id < 0 THEN 3 ELSE 1 END) = st.id
					GROUP BY st.id, sq.customer
					ORDER BY st.id`
		finalQuery, err := helpers.PrepareQuery(endQuery, templateEndQuery)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		if isQuery != "" && isQuery == "1" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"query": finalQuery,
			})
		}

		// fmt.Println(finalQuery)

		data, err := helpers.NewExecuteQuery(finalQuery)
		if err != nil {
			fmt.Println("Error executing query:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal execute query",
				Success: false,
			})
		}

		getRequester, err := helpers.ExecuteQuery(fmt.Sprintf(`SELECT to_char(%s, 'DD Mon YYYY') as tanggal,
															JSONB_BUILD_OBJECT(
																'name', INITCAP(full_name),
																'type', INITCAP('user')
															) as requester
														FROM public.user WHERE id = %s`, date, requestId))
		if err != nil {
			fmt.Println("Error executing query request:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal execute query request",
				Success: false,
			})
		}

		implodedProdukHeader := helpers.SplitToString(produkIds, ",")

		templateEndQuery["QProdukIds"] = implodedProdukHeader

		queryStokProduk := `SELECT sq.* 
							FROM (
							WITH all_products as (SELECT * FROM produk WHERE id IN ({{.QProdukIds}}) ORDER BY produk.order)
							
							SELECT sq.* FROM (
								SELECT p.id, p.order,
										MAX(p.code) as code, 
										COALESCE(SUM(ss.stok_awal),0) - COALESCE(SUM(ssr.jumlah),0) as packs, 
										'stok awal' as tag,
										1 as orders
								FROM all_products p
								LEFT JOIN stok_salesman ss
								ON ss.produk_id = p.id
								AND ss.user_id = {{.QUserID}}
								AND DATE(ss.tanggal_stok) = {{.QDate}} 
								AND ss.produk_id IN ({{.QProdukIds}})
								LEFT JOIN
								(
								SELECT produk_id, condition, pita, SUM(COALESCE(jumlah,0)) AS jumlah 
								FROM stok_salesman_riwayat 
								WHERE is_validate = 1 
									AND DATE(tanggal_riwayat) = {{.QDate}} 
									AND user_id = {{.QUserID}} AND aksi='ORDER'
								GROUP BY produk_id, condition, pita
								) ssr 
										ON ssr.produk_id = ss.produk_id 
								AND ssr.condition = ss.condition 
								AND ssr.pita = ss.pita
								GROUP BY p.id,ss.produk_id, p.order
								ORDER BY p.order
							) sq
							
							UNION
							
							SELECT sq.* FROM (
								SELECT p.id, p.order,
										MAX(p.code) as code, 
										COALESCE(SUM(jumlah),0) as packs, 
										'order' as tag,
										2
								FROM all_products p
								LEFT JOIN stok_salesman_riwayat ssr
								ON ssr.produk_id = p.id
								AND ssr.is_validate = 1
								AND ssr.user_id = {{.QUserID}} 
								AND DATE(tanggal_riwayat) = {{.QDate}} 
								AND produk_id IN ({{.QProdukIds}}) 
								AND UPPER(aksi) = 'ORDER'
								GROUP BY p.id,produk_id, p.order
								ORDER BY p.order
							) sq
							
							UNION
							
							SELECT sq.* FROM (
								SELECT p.id, p.order, 
										MAX(p.code) as code, 
										COALESCE(SUM(sq.total_sales_pack),0) as packs, 
										'penjualan' as tag,
										3
								FROM all_products p
								LEFT JOIN (
									SELECT p.customer_id, 
											pd.produk_id, 
											SUM(pd.jumlah* (pd.harga-pd.diskon)) AS total_sales, 
											SUM(CASE WHEN (pd.harga - pd.diskon) <> 0 THEN pd.jumlah ELSE 0 END) AS total_sales_pack
										FROM penjualan p
										JOIN penjualan_detail pd
										ON p.id = pd.penjualan_id
										WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = {{.QDate}} AND pd.produk_id IN ({{.QProdukIds}})
										GROUP BY p.customer_id, pd.produk_id
									) sq
								ON p.id = sq.produk_id
								GROUP BY p.id,produk_id, p.order
								ORDER BY p.order
							) sq

							UNION
							
							SELECT sq.* FROM (
								SELECT p.id, p.order, 
										MAX(p.code) as code, 
										COALESCE(SUM(sq.total_sales_pack),0) as packs, 
										'program' as tag,
										4
								FROM all_products p
								LEFT JOIN (
									SELECT p.customer_id, 
											pd.produk_id, 
											SUM(pd.jumlah* (pd.harga-pd.diskon)) AS total_sales, 
											SUM(CASE WHEN (pd.harga - pd.diskon) = 0 THEN pd.jumlah ELSE 0 END) AS total_sales_pack
										FROM penjualan p
										JOIN penjualan_detail pd
										ON p.id = pd.penjualan_id
										WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = {{.QDate}} AND pd.produk_id IN ({{.QProdukIds}})
										GROUP BY p.customer_id, pd.produk_id
									) sq
								ON p.id = sq.produk_id
								GROUP BY p.id,produk_id, p.order
								ORDER BY p.order
							) sq
							
							-- UNION
							
							-- SELECT sq.* FROM (
							--   SELECT p.id, p.order, 
							--           MAX(p.code) as code, 
							--           COALESCE(SUM(pjd.jumlah) FILTER (WHERE pjd.produk_id = p.id),0) as packs, 
							--           'program' as tag,
							--           5
							--   FROM penjualan pj
							--   JOIN penjualan_detail pjd
							--     ON pjd.harga = 0
							--     AND pjd.penjualan_id = pj.id
							--     AND (salesman_id = $salesmanId OR p.merchandiser_id = $merchandiserId) 
							--     AND DATE(tanggal_penjualan) = {{.QDate}}
							--   CROSS JOIN all_products p
							--   GROUP BY p.id,produk_id, p.order
							--   ORDER BY p.order
							-- ) sq
							
							UNION
							
							SELECT sq.* FROM (
								SELECT p.id, p.order, 
										MAX(p.code) as code, 
										COALESCE(SUM(pgd.jumlah) FILTER (WHERE DATE(tanggal_pengembalian) = {{.QDate}}),0) as packs, 
										'retur customer' as tag,
										6
								FROM all_products p
								LEFT JOIN pengembalian_detail pgd
								ON pgd.produk_id = p.id
								LEFT JOIN pengembalian pg
								ON pgd.pengembalian_id = pg.id
								AND pg.user_id = {{.QUserID}}
								AND DATE(tanggal_pengembalian) = {{.QDate}} 
								AND produk_id IN ({{.QProdukIds}})
								GROUP BY p.id,produk_id, p.order
								ORDER BY p.order
							) sq
							
							UNION
							
							SELECT sq.* FROM (
								SELECT p.id, p.order, 
										MAX(p.code) as code, 
										COALESCE(SUM(jumlah),0) as packs, 
										'retur gudang' as tag,
										7
								FROM all_products p
								LEFT JOIN stok_salesman_riwayat ssr
								ON ssr.produk_id = p.id
								AND (ssr.user_id = {{.QUserID}}) 
								AND DATE(tanggal_riwayat) = {{.QDate}} 
								AND produk_id IN ({{.QProdukIds}}) 
								AND UPPER(aksi) = 'RETUR'
								AND ssr.is_validate = 1
								GROUP BY p.id,produk_id, p.order
								ORDER BY p.order
							) sq
							
							UNION
							
							SELECT sq.* FROM (
								SELECT p.id, p.order, 
										MAX(p.code) as code, 
										COALESCE(SUM(stok_akhir),0) as packs, 
										'stok akhir' as tag,
										8
								FROM all_products p
								LEFT JOIN stok_salesman ss
								ON ss.produk_id = p.id
								AND (ss.user_id = {{.QUserID}}) 
								AND DATE(tanggal_stok) = {{.QDate}} 
								AND produk_id IN ({{.QProdukIds}})
								GROUP BY p.id,produk_id, p.order
								ORDER BY p.order
							) sq

							UNION
							
							SELECT sq.* FROM (
									SELECT p.id, p.order, 
									MAX(p.code) as code, 
									COUNT(DISTINCT sq.customer_id) FILTER (WHERE total_sales > 0 AND customer_id > 0 AND total_sales_pack >=2) AS total_effective_call,
									'ec brand 2 pack' as tag,
									9999
									FROM all_products p
									LEFT JOIN (
										SELECT p.customer_id, 
												pd.produk_id, 
												SUM(pd.jumlah* (pd.harga-pd.diskon)) AS total_sales, 
												SUM(CASE WHEN pd.harga > 0 THEN pd.jumlah ELSE 0 END) AS total_sales_pack
											FROM penjualan p
											JOIN penjualan_detail pd
											ON p.id = pd.penjualan_id
											WHERE (p.user_id = {{.QUserID}}) AND DATE(p.tanggal_penjualan) = {{.QDate}} AND pd.produk_id IN ({{.QProdukIds}})
											GROUP BY p.customer_id, pd.produk_id
										) sq
									ON p.id = sq.produk_id
									GROUP BY p.id,produk_id, p.order
									ORDER BY p.order
								) sq 
							
								-- SELECT id, all_products.order, MAX(code) as code, 0, 'ec_brand_2_pack' as tag, 99999 FROM all_products GROUP BY id, all_products.order
							-- ) sq
							) sq
							ORDER BY sq.orders, sq.order`

		queryStokProdukExec, err := helpers.PrepareQuery(queryStokProduk, templateEndQuery)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		// fmt.Println(queryStokProdukExec)
		restructuredArray := make(map[string]map[string]interface{})
		testHeader := []string{}

		if implodedProdukHeader != "" {

			dataStokProduk, err := helpers.ExecuteQuery(queryStokProdukExec)
			if err != nil {
				fmt.Println("Error executing query 1:", err)
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal execute query",
					Success: false,
				})
			}

			// Process each item
			for _, item := range dataStokProduk {
				tag := cases.Title(language.English, cases.NoLower).String(item["tag"].(string))
				code := item["code"].(string)
				packs := item["packs"]

				// Initialize map if the key does not exist
				if _, exists := restructuredArray[tag]; !exists {
					restructuredArray[tag] = make(map[string]interface{})
				}

				// Assign the packs value
				restructuredArray[tag][code] = packs
			}

			testHeader = append(testHeader, "Stok Awal")
			testHeader = append(testHeader, "Order")
			testHeader = append(testHeader, "Penjualan")
			testHeader = append(testHeader, "Program")
			testHeader = append(testHeader, "Retur Customer")
			testHeader = append(testHeader, "Retur Gudang")
			testHeader = append(testHeader, "Stok Akhir")
		}

		queryStokItem := `WITH ordersreturs as (
								SELECT user_id, item_id, aksi, SUM(jumlah) as jumlah
								FROM md.stok_merchandiser_riwayat
								WHERE user_id = {{.QUserID}} AND DATE(tanggal_riwayat) = {{.QDate}}
								GROUP BY user_id, item_id, aksi
							), stoks as (
								SELECT user_id, item_id, stok_awal, stok_akhir
								FROM md.stok_merchandiser
								WHERE user_id = {{.QUserID}} AND DATE(tanggal_stok) = {{.QDate}}
							), transactions as (
								SELECT t.user_id, td.item_id, SUM(qty) as jumlah
								FROM md.transaction t
								JOIN md.transaction_detail td
									ON t.id = td.transaction_id
								WHERE t.user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}}
								GROUP BY t.user_id, td.item_id
							)

							SELECT i.name as item_name,
									JSONB_BUILD_OBJECT(
										'Stok Awal', s.stok_awal,
										'Order', CASE WHEN o.aksi = 'ORDER' THEN o.jumlah ELSE 0 END,
										'Transaksi', COALESCE(t.jumlah,0),
										'Retur Gudang', CASE WHEN o.aksi = 'RETUR' THEN o.jumlah ELSE 0 END,
										'Stok Akhir', s.stok_akhir
									) as datas
							FROM stoks s
							FULL JOIN ordersreturs o
								ON s.item_id = o.item_id
							FULL JOIN transactions t
								ON s.item_id = t.item_id
							FULL JOIN md.merchandiser m
								ON m.user_id = COALESCE(s.user_id, o.user_id, t.user_id)
							LEFT JOIN md.item i
								ON s.item_id = i.id
							WHERE m.user_id = {{.QUserID}}
							ORDER BY i.name`

		queryStokItemExec, err := helpers.PrepareQuery(queryStokItem, templateEndQuery)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		dataStokItem, err := helpers.ExecuteQuery(queryStokItemExec)
		if err != nil {
			fmt.Println("Error executing query 2:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal execute query",
				Success: false,
			})
		}

		testHeaderItem := []string{}
		testHeaderItem = append(testHeaderItem, "Stok Awal")
		testHeaderItem = append(testHeaderItem, "Order")
		testHeaderItem = append(testHeaderItem, "Transaksi")
		testHeaderItem = append(testHeaderItem, "Retur Gudang")
		testHeaderItem = append(testHeaderItem, "Stok Akhir")

		querySummary := `WITH tunais as (
                          SELECT SUM(sum_cash) as total_tunai,
                          {{.QAsUserID}}
                          FROM (
                              SELECT COALESCE(SUM(CASE WHEN p.is_kredit = 0 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END),0) as sum_cash, 1 as param
                              FROM penjualan p
                              LEFT JOIN penjualan_detail pd
                                      ON p.id = pd.penjualan_id
                              LEFT JOIN payment py
                                      ON p.id = py.penjualan_id
                              WHERE py.id IS NULL 
                                      AND p.is_kredit = 0
                                      AND p.user_id = {{.QUserID}}
                                      AND DATE(p.tanggal_penjualan) = {{.QDate}}
                                      
                              UNION ALL

                              SELECT SUM(py.nominal) as sum_cash
                                -- COALESCE(SUM(CASE WHEN p.is_kredit = 0 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END),0) as sum_cash
                                , 2
                              FROM penjualan p
                              JOIN payment py
									ON p.id = py.penjualan_id
                              LEFT JOIN pembayaran_piutang pp 
									ON py.id = pp.payment_id
                              WHERE UPPER(py.tipe) = 'CASH'
                                      AND pp.id IS NULL
                                      AND p.user_id = {{.QUserID}}
                                      AND DATE(p.tanggal_penjualan) = {{.QDate}}
                                      
                              UNION ALL

                              SELECT COALESCE(pp.total_pembayaran,0) as sum_cash, 3
                              FROM pembayaran_piutang pp
                              LEFT JOIN payment py 
                                      ON pp.payment_id = py.id
                              WHERE py.id IS NULL
                                      AND pp.user_id = {{.QUserID}}
                                      AND DATE(pp.tanggal_pembayaran) = {{.QDate}}
                                      
                              UNION ALL

                              SELECT COALESCE(pp.total_pembayaran,0) as sum_cash, 4
                              FROM pembayaran_piutang pp
                              JOIN payment py 
                                      ON pp.payment_id = py.id
                              WHERE UPPER(py.tipe) = 'CASH'
                                      AND pp.user_id = {{.QUserID}}
                                      AND DATE(pp.tanggal_pembayaran) = {{.QDate}}
                      ) sq
                  ), total_call as (
                      SELECT COUNT(kunjungan.id) AS total_call, {{.QAsUserID}}
                          FROM (
                                          SELECT            
                                          DISTINCT ON(k.customer_id) k.id, k.salesman_id
                                          FROM kunjungan k 
                                          -- LEFT JOIN penjualan p 
                                          -- ON p.customer_id = k.customer_id AND DATE(p.tanggal_penjualan) =  DATE({{.QDate}}) 
                                          -- AND (p.salesman_id = $salesmanId OR p.merchandiser_id = $merchandiserId)
                                          WHERE k.user_id = {{.QUserID}}
                                           AND DATE(tanggal_kunjungan) = DATE({{.QDate}}) AND UPPER(k.status_toko) = 'BUKA' AND k.customer_id > 0
                                          ORDER BY k.customer_id, tanggal_kunjungan ASC
                          ) kunjungan
                  ), sales as (
                      SELECT
                              SUM(CASE WHEN pd.harga > 0 THEN pd.jumlah ELSE 0 END) AS total_sales_pack,
                              SUM(pd.jumlah* (pd.harga-pd.diskon)) AS total_sales, 
                              COUNT(DISTINCT p.customer_id) FILTER (WHERE p.customer_id > 0 AND pd.jumlah >= 2 AND pd.harga > 0) AS total_effective_call,
                              SUM(CASE WHEN p.is_kredit = 0 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END) AS total_cash,
                              SUM(CASE WHEN p.is_kredit = 1 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END) AS total_credit,
                              SUM(CASE WHEN pd.harga > 0 AND p.is_kredit = 0 THEN pd.jumlah ELSE 0 END) AS total_cash_pack,
                              SUM(CASE WHEN pd.harga > 0 AND p.is_kredit = 1 THEN pd.jumlah ELSE 0 END) AS total_credit_pack,
                              SUM(CASE WHEN pd.harga = 0 THEN pd.jumlah ELSE 0 END) AS total_program_pack,
                              {{.QAsUserID}}
                              FROM penjualan p
                              JOIN penjualan_detail pd
                              ON p.id = pd.penjualan_id
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = DATE({{.QDate}})
                  ), retur as (
                      SELECT 
                              COALESCE(SUM(pd.jumlah * pd.harga),0) AS total_return,
                              COALESCE(SUM(pd.jumlah),0) AS total_return_pack,
                              {{.QAsUserID}}
                      FROM 
                              pengembalian p
                              JOIN pengembalian_detail pd
                              ON p.id = pd.pengembalian_id
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pengembalian) = DATE({{.QDate}})
                  ),customer_register as (
                      SELECT 
                              COUNT(p.id) as customer_register,
                              {{.QAsUserID}}
                      FROM 
                      customer p
                      WHERE p.user_id_holder = {{.QUserID}} AND DATE(p.dtm_crt) = DATE({{.QDate}})
                  ), timecall as (
                      SELECT 
                              SUM(AGE(p.checkout_at,p.checkin_at))/CASE WHEN count(distinct customer_id) = 0 THEN 1 ELSE count(distinct customer_id) END as time_call,
                              {{.QAsUserID}}
                      FROM 
                              kunjungan_log p
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.checkin_at) = DATE({{.QDate}})
                  ), payments as (
                      SELECT
                          SUM(CASE WHEN UPPER(p.tipe) = 'CASH' THEN nominal ELSE 0 END) AS total_payment_cash,
                          SUM(CASE WHEN UPPER(p.tipe) = 'TRANSFER' THEN nominal ELSE 0 END) AS total_payment_transfer,
                          SUM(CASE WHEN UPPER(p.tipe) = 'CEK' THEN nominal ELSE 0 END) AS total_payment_check,
                          SUM(CASE WHEN UPPER(p.tipe) = 'BILYET GIRO' AND (DATE(tanggal_cair) != DATE({{.QDate}}) OR tanggal_cair IS NULL) THEN nominal ELSE 0 END) AS total_payment_open_bilyet_giro,
                          SUM(CASE WHEN UPPER(p.tipe) = 'BILYET GIRO' AND is_cair = 1 AND DATE(tanggal_cair) = DATE({{.QDate}})  THEN nominal ELSE 0 END) AS total_payment_bilyet_giro,
                          {{.QAsUserID}}
                      FROM
                          payment p
                      JOIN penjualan pj
                      ON pj.id = p.penjualan_id 
                      WHERE p.user_id = {{.QUserID}} AND p.is_verif>-1 --AND pj.is_kredit = 0 
                      AND (
                          DATE(p.tanggal_transaksi) = DATE({{.QDate}}) 
                          OR (DATE(p.tanggal_cair) = DATE({{.QDate}}) AND p.is_cair = 1)
                      )
                  ) , 
                  pembayaran_piutang as (
                              SELECT COALESCE(SUM(p.total_pembayaran),0) as total_payment_and_dp,
                              {{.QAsUserID}}
                              
                              FROM
                              pembayaran_piutang p
                              JOIN pembayaran_piutang_detail pd 
                              ON p.id = pd.pembayaran_piutang_id
                              LEFT JOIN piutang pi
                              ON pi.id = pd.piutang_id
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pembayaran) = DATE({{.QDate}})
                  )

                  SELECT 
                              SUM(datas.total_call) as total_call_buka,
                              SUM(datas.total_effective_call) AS total_effective_call_2pack,
                              SUM(datas.customer_register) as total_register_customer,
                              MAX(datas.time_call) as average_call,
                              SUM(datas.total_sales_pack) AS omzet,
                              SUM(datas.total_cash_pack) AS omzet_tunai,
                              SUM(datas.total_credit_pack) AS omzet_kredit,
                              SUM(datas.total_return_pack) AS retur,
                              --COALESCE(SUM(datas.total_payment_cash) - SUM(datas.total_return),0) AS pembayaran_tunai,
                              COALESCE(MAX(datas.total_tunai) - SUM(datas.total_return),0) AS pembayaran_tunai,
                              COALESCE(SUM(datas.total_payment_transfer),0) AS pembayaran_transfer,
                              COALESCE(SUM(datas.total_payment_check),0) AS cek,
                              COALESCE(SUM(datas.total_payment_open_bilyet_giro),0) AS bilyet_giro_baru,
                              COALESCE(SUM(datas.total_payment_bilyet_giro),0) AS bilyet_giro_cair,
                              SUM(datas.total_cash)+ SUM(datas.total_payment_and_dp)
                                -SUM(datas.total_return)-COALESCE(SUM(datas.total_payment_transfer),0)
                                -COALESCE(SUM(datas.total_payment_check),0)
                                -COALESCE(SUM(datas.total_payment_open_bilyet_giro),0)
                                -COALESCE(SUM(datas.total_payment_bilyet_giro),0) AS total_setoran,
                              MAX(datas.total_tunai)-SUM(datas.total_return) AS total_setoran_tunai
                  FROM ( SELECT tc.*, sales.*, retur.*, cr.*, timecall.*, payments.*, pp.*, tunais.*
                                  FROM public.user s
                                  LEFT JOIN total_call tc
                                      ON s.id = tc.user_id
                                  LEFT JOIN sales
                                      ON s.id = sales.user_id
                                  LEFT JOIN retur
                                      ON s.id = retur.user_id
                                  LEFT JOIN customer_register cr
                                      ON s.id = cr.user_id
                                  LEFT JOIN timecall
                                      ON s.id = timecall.user_id
                                  LEFT JOIN payments
                                      ON s.id = payments.user_id
                                  LEFT JOIN pembayaran_piutang pp
                                      ON s.id = pp.user_id
                                  LEFT JOIN tunais
                                      ON s.id = tunais.user_id
                  ) datas`

		querySummaryExec, err := helpers.PrepareQuery(querySummary, templateEndQuery)

		// fmt.Println(querySummaryExec)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		dataSummary, err := helpers.ExecuteQuery(querySummaryExec)
		if err != nil {
			fmt.Println("Error executing query 3:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal execute query",
				Success: false,
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":             "success",
			"success":             true,
			"data":                data,
			"tanggal":             getRequester[0]["tanggal"],
			"user":                getRequester[0]["requester"],
			"table_produk":        restructuredArray,
			"table_produk_header": testHeader,
			"table_item":          dataStokItem,
			"table_item_header":   testHeaderItem,
			"summary":             dataSummary,
		})
	} else {

		endQuery := `WITH visits as (
						SELECT id, tanggal_kunjungan, customer_id, user_id, subject_type_id, image_kunjungan
						FROM kunjungan 
						WHERE user_id = {{.QUserID}} 
							AND DATE(tanggal_kunjungan) = {{.QDate}}
						ORDER BY tanggal_kunjungan
					), penjualan_cus as (
						SELECT p.no_nota,
										p.id,
										CASE WHEN COALESCE(p.image_nota_print, p.image_nota) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/nota/'||COALESCE(p.image_nota_print, p.image_nota) ELSE NULL END as image_nota,
										CASE WHEN p.image_bukti_serah IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/nota/'||p.image_bukti_serah ELSE NULL END as image_bukti_serah,
										p.user_id,
										p.customer_id,
										DATE(p.tanggal_penjualan),
										COALESCE(SUM(CASE WHEN p.is_kredit = 0 THEN ((pd.harga - pd.diskon) * pd.jumlah) ELSE 0 END),0) as tunai,
										COALESCE(SUM(CASE WHEN p.is_kredit = 1 THEN ((pd.harga - pd.diskon) * pd.jumlah) ELSE 0 END),0) as kredit,
										pd.produk_id as produk_id,
										COALESCE(SUM(pd.jumlah),0) as jumlah
						FROM penjualan p
						JOIN penjualan_detail pd
							ON p.id = pd.penjualan_id
						WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = {{.QDate}}
						GROUP BY p.id, pd.produk_id
					), pembayaran_cus as (
						SELECT p.customer_id,
										p.user_id,
										COALESCE(SUM(pd.nominal),0) as pembayaran,
										COALESCE((CASE WHEN p.tipe_pelunasan = 0 AND p.pengembalian_id IS NOT NULL then p.total_pembayaran ELSE 0 END),0) AS adjustment,
										COALESCE(SUM((CASE WHEN tipe_pelunasan = 1 then pd.nominal ELSE 0 END)),0) as dp
						FROM pembayaran_piutang p
						JOIN pembayaran_piutang_detail pd
							ON p.id = pd.pembayaran_piutang_id
						LEFT JOIN payment py ON py.id =  p.payment_id AND py.tipe = 'BILYET GIRO'
						WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pembayaran) = {{.QDate}}
						GROUP BY p.id
					), pengembalian_cus as (
						SELECT p.customer_id,
										p.user_id,
										COALESCE(SUM(pd.jumlah * pd.harga),0) as pengembalian
						FROM pengembalian p
						JOIN pengembalian_detail pd
							ON p.id = pd.pengembalian_id
					WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pengembalian) = {{.QDate}}
						GROUP BY p.id
					), payments as (
						SELECT UPPER((CASE WHEN tipe = 'BILYET GIRO' THEN 'OPEN ' ELSE '' END)||tipe) AS payment_tipe, 
								COALESCE(nominal,0) AS payment_nominal,
										p.user_id,
										p.customer_id,
										CASE WHEN p.bukti_bayar = '' THEN NULL ELSE p.bukti_bayar END as payment_image,
																				p.penjualan_id as penjualan_id
								FROM payment p
								WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_transaksi) = {{.QDate}} 

						UNION

						SELECT 'CASH' as payment_tipe,
								SUM(((pd.harga - pd.diskon) * pd.jumlah)) as payment_nominal,
								p.user_id,
								p.customer_id,
								null as payment_image,
								p.id as penjualan_id
						FROM penjualan p
						LEFT JOIN penjualan_detail pd
								ON p.id = pd.penjualan_id
						LEFT JOIN payment py
								ON p.id = py.penjualan_id
								AND py.id IS NULL 
						WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_penjualan) = {{.QDate}} AND p.is_kredit = 0 
						GROUP BY p.id
					), kunjunganlogs as (
						SELECT DISTINCT ON (customer_id) * 
						FROM kunjungan_log WHERE user_id = {{.QUserID}} AND date(checkin_at) = {{.QDate}}
						ORDER BY customer_id, checkin_at
					), transactions as (
						SELECT t.user_id,
										t.subject_type_id,
										t.customer_id,
											JSONB_BUILD_OBJECT(
												'name', i.name,
												'quantity', COALESCE(SUM(td.qty),0)
										) as items,
										CASE WHEN ic.id = 191 AND t.transaction_type_id <> 1 
											THEN JSONB_BUILD_OBJECT(
												'photo_before', CASE WHEN MAX(t.photos_before) = '{}' OR MAX(t.photos_before) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos_before)) image
																	) END,
												'photo_after', CASE WHEN MAX(t.photos_after) = '{}' OR MAX(t.photos_after) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos_after)) image
																	) END
											)
											ELSE NULL
										END as sampling,
										CASE WHEN ic.id = 171 AND t.transaction_type_id <> 1 
											THEN JSONB_BUILD_OBJECT(
												'photo_before', CASE WHEN MAX(t.photos_before) = '{}' OR MAX(t.photos_before) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos_before)) image
																		) END,
												'photo_after', CASE WHEN MAX(t.photos_after) = '{}' OR MAX(t.photos_after) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos_after)) image
																		) END
											)
											ELSE NULL
										END as posm,
										CASE WHEN ic.id = 161 AND t.transaction_type_id <> 1 
											THEN JSONB_BUILD_OBJECT(
												'photo_before', CASE WHEN MAX(t.photos_before) = '{}' OR MAX(t.photos_before) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos_before)) image
																		) END,
												'photo_after', CASE WHEN MAX(t.photos_after) = '{}' OR MAX(t.photos_after) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos_after)) image
																		) END
											)
											ELSE NULL
										END as merchandise
										
						FROM md.transaction t
						JOIN md.transaction_detail td
							ON t.id = td.transaction_id
						JOIN md.item i
							ON td.item_id = i.id
						JOIN md.item_category ic
							ON i.category_id = ic.id
						WHERE t.user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}}
						GROUP BY t.user_id, t.customer_id, ic.id, t.transaction_type_id, t.subject_type_id, i.id
					)

					SELECT st.name as subject_type, 
							sq.customer, 
							MAX(invoice) as invoice,  
							{{.QSelectProd}}
							JSONB_AGG(sq.items) FILTER (WHERE sq.items IS NOT NULL) as items,
							COALESCE(SUM(sq.tunai),0) as tunai,
							COALESCE(SUM(sq.kredit),0) as kredit,
							COALESCE(SUM(sq.pengembalian), 0) as retur,
							COALESCE(SUM(sq.pembayaran),0) as pembayaran,
							COALESCE(SUM(sq.adjustment),0) as adjustment,
							COALESCE(SUM(sq.tunai),0) + 
							COALESCE(SUM(sq.pembayaran),0) + 
							COALESCE(SUM(sq.dp),0) - 
							COALESCE(SUM(sq.pengembalian),0) + 
							COALESCE(SUM(sq.adjustment),0) as setoran,
							JSONB_AGG(sq.payment_information) FILTER (WHERE sq.payment_information IS NOT NULL)->0 as payment_information,
							JSONB_AGG(sq.time_info) FILTER (WHERE sq.time_info IS NOT NULL)->0 as time_info,
							JSONB_AGG(sq.nota) FILTER (WHERE sq.nota IS NOT NULL)->0 as nota,
							JSONB_AGG(sq.penyerahan_produk) FILTER (WHERE sq.penyerahan_produk IS NOT NULL)->0 as penyerahan_produk,
							JSONB_AGG(sq.posm) FILTER (WHERE sq.posm IS NOT NULL)->0 as posm,
							JSONB_AGG(sq.merchandise) FILTER (WHERE sq.merchandise IS NOT NULL)->0 as merchandise,
							JSONB_AGG(sq.sampling) FILTER (WHERE sq.sampling IS NOT NULL)->0 as sampling
					FROM (
						SELECT v.subject_type_id,
								COALESCE(c.id,v.customer_id, pj.customer_id, kl.customer_id) as customer_id,
								JSONB_BUILD_OBJECT(
										'name', COALESCE(c.name, 'End User'),
										'contact', COALESCE(c.outlet_name, '-'),
										'location', COALESCE(c.kelurahan, '-'),
										'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
								) as customer,
								pj.no_nota as invoice,
								{{.QProd}}
								null::jsonb as items,
								COALESCE(SUM(pj.tunai),0) as tunai,
								COALESCE(SUM(pj.kredit),0) as kredit,
								0 as pembayaran,
								0 as pengembalian,
								0 as adjustment,
								0 as dp,
								NULL as payment_information,
								MIN(kl.checkin_at) as param,
								JSONB_BUILD_OBJECT(
									'checkin', to_char(MIN(kl.checkin_at) , 'HH24:MI'),
									'checkout', to_char(MAX(kl.checkout_at), 'HH24:MI'),
									'durasi', MAX(kl.checkout_at) - MIN(kl.checkin_at)
								) as time_info,
								pj.image_nota as nota,
								pj.image_bukti_serah as penyerahan_produk,
								null::jsonb as posm,
								null::jsonb as merchandise,
								null::jsonb as sampling
							FROM visits v
							FULL JOIN penjualan_cus pj
								ON v.user_id = pj.user_id
								AND v.customer_id = pj.customer_id
							LEFT JOIN kunjunganlogs kl
								ON v.user_id = kl.user_id
								AND v.customer_id = kl.customer_id
							LEFT JOIN customer c
								ON COALESCE(v.customer_id, pj.customer_id, kl.customer_id) = c.id
							GROUP BY v.subject_type_id, v.image_kunjungan, c.id, COALESCE(c.id,v.customer_id, pj.customer_id, kl.customer_id), pj.no_nota, pj.image_nota, pj.image_bukti_serah
							
							UNION

							SELECT v.subject_type_id,
									COALESCE(c.id,v.customer_id, pem.customer_id, pg.customer_id, py.customer_id, kl.customer_id) as customer_id,
									JSONB_BUILD_OBJECT(
											'name', COALESCE(c.name, 'End User'),
											'contact', COALESCE(c.outlet_name, '-'),
											'location', COALESCE(c.kelurahan, '-'),
											'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
									) as customer,
									NULL as invoice,
									{{.QProdItem}}
									null::jsonb as items,
									0 as tunai,
									0 as kredit,
									COALESCE(SUM(pem.pembayaran),0) as pembayaran,
									COALESCE(SUM(pg.pengembalian),0) as pengembalian,
									COALESCE(SUM(pem.adjustment),0) as adjustment,
									COALESCE(SUM(pem.dp),0) as dp,
									JSONB_BUILD_OBJECT(
									'cash', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'CASH'), 0) - COALESCE(SUM(pg.pengembalian),0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'CASH' AND py.payment_image <> '')
										),
									'transfer', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'TRANSFER'), 0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'TRANSFER' AND py.payment_image <> '')
										),
									'bilyet_giro_cair', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'BILYET GIRO'), 0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'BILYET GIRO' AND py.payment_image <> '')
										),
									'bilyet_giro_open', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'OPEN BILYET GIRO'), 0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'OPEN BILYET GIRO' AND py.payment_image <> '')
										)
									) as payment_information,
									MIN(kl.checkin_at) as param,
									JSONB_BUILD_OBJECT(
										'checkin', to_char(MIN(kl.checkin_at) , 'HH24:MI'),
										'checkout', to_char(MAX(kl.checkout_at), 'HH24:MI'),
										'durasi', MAX(kl.checkout_at) - MIN(kl.checkin_at)
									) as time_info,
									null as nota,
									null as penyerahan_produk,
									null::jsonb as posm,
									null::jsonb as merchandise,
									null::jsonb as sampling
							FROM visits v
							FULL JOIN pembayaran_cus pem
								ON v.user_id = pem.user_id
								AND v.customer_id = pem.customer_id
							FULL JOIN pengembalian_cus pg
								ON v.user_id = pg.user_id
								AND v.customer_id = pg.customer_id
							FULL JOIN payments py
								ON v.user_id = py.user_id
								AND v.customer_id = py.customer_id
							LEFT JOIN kunjunganlogs kl
								ON v.user_id = kl.user_id
								AND v.customer_id = kl.customer_id
							LEFT JOIN customer c
								ON COALESCE(v.customer_id, pem.customer_id, pg.customer_id, py.customer_id, kl.customer_id) = c.id
							GROUP BY v.subject_type_id, v.image_kunjungan, c.id, COALESCE(c.id,v.customer_id, pem.customer_id, pg.customer_id, py.customer_id, kl.customer_id)
							
							UNION
							
								SELECT v.subject_type_id,
									COALESCE(c.id, v.customer_id, t.customer_id, kl.customer_id) as customer_id,
									JSONB_BUILD_OBJECT(
											'name', COALESCE(c.name, 'End User'),
											'contact', COALESCE(c.outlet_name, '-'),
											'location', COALESCE(c.kelurahan, '-'),
											'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
									) as customer,
									NULL as invoice,
									{{.QProdItem}}
									t.items as items,
									0 as tunai,
									0 as kredit,
									0 as pembayaran,
									0 as pengembalian,
									0 as adjustment,
									0 as dp,
									NULL as payment_information,
									MIN(kl.checkin_at) as param,
									JSONB_BUILD_OBJECT(
										'checkin', to_char(MIN(kl.checkin_at) , 'HH24:MI'),
										'checkout', to_char(MAX(kl.checkout_at), 'HH24:MI'),
										'durasi', MAX(kl.checkout_at) - MIN(kl.checkin_at)
									) as time_info,
									null as nota,
									null as penyerahan_produk,
									t.posm as posm,
									t.merchandise as merchandise,
									t.sampling as sampling
							FROM visits v
							FULL JOIN transactions t
								ON v.user_id = t.user_id
								AND v.customer_id = t.customer_id
							LEFT JOIN kunjunganlogs kl
								ON v.user_id = kl.user_id
								AND v.customer_id = kl.customer_id
							LEFT JOIN customer c
								ON COALESCE(v.customer_id, t.customer_id, kl.customer_id) = c.id
							GROUP BY v.subject_type_id, v.image_kunjungan, c.id, COALESCE(c.id, v.customer_id, t.customer_id, kl.customer_id), t.items, t.posm, t.merchandise, t.sampling
					) sq
					JOIN md.subject_type st
						ON COALESCE(sq.subject_type_id, CASE WHEN sq.customer_id < 0 THEN 3 ELSE 1 END) = st.id
					GROUP BY sq.param, st.id, sq.customer
					ORDER BY sq.param`

		finalQuery, err := helpers.PrepareQuery(endQuery, templateEndQuery)

		// fmt.Println(finalQuery)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		data, err := helpers.NewExecuteQuery(finalQuery)
		if err != nil {
			fmt.Println("Error executing query 1:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal execute query",
				Success: false,
			})
		}

		queryStokProduk := `SELECT 
								JSON_BUILD_OBJECT('produk_id',p.id, 'kode', p.code, 'nama', p.name, 'foto', p.foto) AS produk,
								SUM(COALESCE(ss.stok_awal,0)-COALESCE(ssr_order.jumlah,0)) AS stok_awal,
								SUM(COALESCE(ssr_order.jumlah,0)) AS order,
								SUM(COALESCE(pj.jumlah,0)) AS penjualan,
								SUM(COALESCE(prg.jumlah,0)) AS program,
								SUM(COALESCE(pg.jumlah,0)) AS retur_customer,
								SUM(COALESCE(ssr_retur.jumlah,0)) AS retur_gudang,
								SUM(COALESCE(ss.stok_akhir,0)) AS stok_akhir
							FROM produk p  
							JOIN produk_branch pb
							ON pb.produk_id = p.id
							LEFT JOIN 
								stok_salesman ss 
								ON ss.produk_id = p.id AND ss.user_id = {{.QUserID}} AND DATE(ss.tanggal_stok) = {{.QDate}}
							LEFT JOIN
							(
								SELECT produk_id, condition, pita, SUM(COALESCE(jumlah,0)) AS jumlah 
								FROM stok_salesman_riwayat 
								WHERE is_validate = 1 AND DATE(tanggal_riwayat) = {{.QDate}} 
								AND user_id = {{.QUserID}} AND aksi='ORDER'
								GROUP BY produk_id, condition, pita
							) ssr_order ON ssr_order.produk_id = ss.produk_id AND ssr_order.condition = ss.condition AND ssr_order.pita = ss.pita
							LEFT JOIN
							( 
								SELECT produk_id, condition, pita, SUM(COALESCE(jumlah,0)) AS jumlah 
								FROM stok_salesman_riwayat 
								WHERE is_validate = 1 AND DATE(tanggal_riwayat) = {{.QDate}} 
								AND user_id = {{.QUserID}} AND aksi='RETUR'
								GROUP BY produk_id, condition, pita
							) ssr_retur ON ssr_retur.produk_id = p.id AND ssr_retur.condition = ss.condition AND ssr_retur.pita = ss.pita
							LEFT JOIN 
							(
								SELECT pd.produk_id, pd.condition, pd.pita, SUM(COALESCE(pd.jumlah,0)) AS jumlah
								FROM penjualan p 
								JOIN penjualan_detail pd 
								ON pd.penjualan_id = p.id
								WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_penjualan) = {{.QDate}} AND pd.harga >0
								GROUP BY pd.produk_id, pd.condition, pd.pita
							) pj ON pj.produk_id = p.id AND pj.condition = ss.condition AND pj.pita = ss.pita
							LEFT JOIN
							(
								SELECT pd.produk_id, pd.condition, pd.pita, SUM(COALESCE(pd.jumlah,0)) AS jumlah
								FROM penjualan p 
								JOIN penjualan_detail pd 
								ON pd.penjualan_id = p.id
								WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_penjualan) = {{.QDate}} AND pd.harga =0
								GROUP BY pd.produk_id, pd.condition, pd.pita
							) prg ON prg.produk_id = p.id AND prg.condition = ss.condition AND prg.pita = ss.pita
							LEFT JOIN
							(
								SELECT pd.produk_id, pd.condition, pd.pita, SUM(COALESCE(pd.jumlah,0)) AS jumlah
								FROM pengembalian p 
								JOIN pengembalian_detail pd 
								ON pd.pengembalian_id = p.id
								WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_pengembalian) = {{.QDate}}
								GROUP BY pd.produk_id, pd.condition, pd.pita
							) pg ON pg.produk_id = p.id AND pg.condition = ss.condition AND pg.pita = ss.pita
							WHERE pb.branch_id IN ( {{.QBranchID}} )
							GROUP BY p.id
							ORDER BY p.order`

		queryStokProdukExec, err := helpers.PrepareQuery(queryStokProduk, templateEndQuery)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		// fmt.Println(queryStokProdukExec)

		dataStokProduk, err := helpers.NewExecuteQuery(queryStokProdukExec)
		if err != nil {
			fmt.Println("Error executing query 2:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal execute query",
				Success: false,
			})
		}

		testHeader := []string{}
		testHeader = append(testHeader, "Stok Awal")
		testHeader = append(testHeader, "Order")
		testHeader = append(testHeader, "Penjualan")
		testHeader = append(testHeader, "Program")
		testHeader = append(testHeader, "Retur Customer")
		testHeader = append(testHeader, "Retur Gudang")
		testHeader = append(testHeader, "Stok Akhir")

		queryStokItem := `WITH ordersreturs as (
							SELECT user_id, item_id, aksi, SUM(jumlah) as jumlah
							FROM md.stok_merchandiser_riwayat
							WHERE user_id = {{.QUserID}} AND DATE(tanggal_riwayat) = {{.QDate}}
							GROUP BY user_id, item_id, aksi
						), stoks as (
							SELECT user_id, item_id, stok_awal, stok_akhir
							FROM md.stok_merchandiser
							WHERE user_id = {{.QUserID}} AND DATE(tanggal_stok) = {{.QDate}}
						), transactions as (
							SELECT t.user_id, td.item_id, SUM(qty) as jumlah
							FROM md.transaction t
							JOIN md.transaction_detail td
								ON t.id = td.transaction_id
							WHERE t.user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}}
							GROUP BY t.user_id, td.item_id
						)

						SELECT i.name as item_name,
								JSONB_BUILD_OBJECT(
									'Stok Awal', s.stok_awal,
									'Order', CASE WHEN o.aksi = 'ORDER' THEN o.jumlah ELSE 0 END,
									'Transaksi', COALESCE(t.jumlah,0),
									'Retur Gudang', CASE WHEN o.aksi = 'RETUR' THEN o.jumlah ELSE 0 END,
									'Stok Akhir', s.stok_akhir
								) as datas
						FROM stoks s
						FULL JOIN ordersreturs o
							ON s.item_id = o.item_id
						FULL JOIN transactions t
							ON s.item_id = t.item_id
						FULL JOIN md.merchandiser m
							ON m.user_id = COALESCE(s.user_id, o.user_id, t.user_id)
						LEFT JOIN md.item i
							ON s.item_id = i.id
						WHERE m.user_id = {{.QUserID}}
						ORDER BY i.name`

		queryStokItemExec, err := helpers.PrepareQuery(queryStokItem, templateEndQuery)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		dataStokItem, err := helpers.ExecuteQuery(queryStokItemExec)
		if err != nil {
			fmt.Println("Error executing query 3:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal execute query",
				Success: false,
			})
		}

		testHeaderItem := []string{}
		testHeaderItem = append(testHeaderItem, "Stok Awal")
		testHeaderItem = append(testHeaderItem, "Order")
		testHeaderItem = append(testHeaderItem, "Transaksi")
		testHeaderItem = append(testHeaderItem, "Retur Gudang")
		testHeaderItem = append(testHeaderItem, "Stok Akhir")

		queryVisit := `WITH data_aktifitas as (
						SELECT sq.customer_id, sq.date, COUNT(id) FROM (
						SELECT id, customer_id, DATE(tanggal_penjualan) FROM penjualan WHERE user_id = {{.QUserID}} AND DATE(tanggal_penjualan) = {{.QDate}}
						UNION
						SELECT id, customer_id, DATE(tanggal_pengembalian) FROM pengembalian WHERE user_id = {{.QUserID}} AND DATE(tanggal_pengembalian) = {{.QDate}}
						UNION
						SELECT id, customer_id, DATE(tanggal_pembayaran) FROM pembayaran_piutang WHERE user_id = {{.QUserID}} AND DATE(tanggal_pembayaran) = {{.QDate}}
						UNION
						SELECT id, customer_id, DATE(datetime) FROM qr_code_history WHERE user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}}
						UNION
						SELECT sq.* FROM (
							SELECT DISTINCT ON (customer_id) id, customer_id, DATE(checkin_at)FROM kunjungan_log WHERE user_id = {{.QUserID}} AND DATE(checkin_at) = {{.QDate}} ORDER BY customer_id, checkin_at DESC
						) sq
						UNION
						SELECT sq.* FROM (
							SELECT DISTINCT ON (customer_id) id, customer_id, DATE(checkout_at) FROM kunjungan_log WHERE user_id = {{.QUserID}} AND DATE(checkout_at) = {{.QDate}} ORDER BY customer_id, checkout_at DESC
						) sq
						) sq
						GROUP BY sq.customer_id, sq.date
						)


						SELECT DISTINCT ON (c.id, DATE(tanggal_kunjungan)) 
							ROW_NUMBER() OVER(ORDER BY tanggal_kunjungan) as no,
							DATE(k.tanggal_kunjungan),
							s.name,
							JSONB_BUILD_OBJECT(
								'id', c.id||'',
								'name', c.name,
								'outlet_name', c.outlet_name,
								'type', COALESCE(k.customer_tipe, ct.name)
							) as customer,
							c.alamat,
							k.status_toko,
							k.keterangan,
							JSONB_BUILD_OBJECT(
								'checkin', to_char(kl.checkin_at, 'HH24:MI'),
								'checkout', to_char(kl.checkout_at, 'HH24:MI'),
								'durasi', AGE(kl.checkout_at, kl.checkin_at),
								'jumlah_aktifitas', da.count
							) as informasi,
							k.latitude_longitude as lokasi,
							k.image_kunjungan,
							CASE WHEN p.id IS NOT NULL THEN 1 ELSE 0 END as is_ec,
							CASE WHEN DATE(c.dtm_crt) = DATE(k.tanggal_kunjungan) THEN 1 ELSE 0 END as new_register
					FROM kunjungan k
					JOIN salesman s
						ON k.salesman_id = s.id
					JOIN customer c
						ON k.customer_id = c.id
					JOIN customer_type ct
						ON c.tipe = ct.id
					LEFT JOIN (
						SELECT ssq.customer_id, MAX(checkin_at) as checkin_at, MAX(checkout_at) as checkout_at
						FROM (
							SELECT sq.* FROM (
								SELECT DISTINCT ON (customer_id) customer_id, checkin_at, NULL::timestamp as checkout_at 
									FROM kunjungan_log 
									WHERE user_id = {{.QUserID}} 
									AND DATE(checkin_at) = {{.QDate}} 
									ORDER BY customer_id, checkin_at DESC
							) sq
							UNION
							SELECT sq.* FROM (
								SELECT DISTINCT ON (customer_id) customer_id, NULL::timestamp, checkout_at 
									FROM kunjungan_log 
									WHERE user_id = {{.QUserID}} 
									AND DATE(checkout_at) = {{.QDate}} 
									ORDER BY customer_id, checkout_at DESC
							) sq
						) ssq
						GROUP BY ssq.customer_id
					) kl
						ON k.customer_id = kl.customer_id
					LEFT JOIN data_aktifitas da
						ON k.customer_id = da.customer_id
						AND DATE(k.tanggal_kunjungan) = da.date
					LEFT JOIN penjualan p
						ON k.customer_id = p.customer_id
						AND DATE(k.tanggal_kunjungan) = DATE(p.tanggal_penjualan)
					WHERE k.user_id IN ({{.QUserID}}) 
						AND DATE(k.tanggal_kunjungan) = {{.QDate}}
					ORDER BY c.id, DATE(tanggal_kunjungan)`

		queryVisitExec, err := helpers.PrepareQuery(queryVisit, templateEndQuery)

		// fmt.Println(querySummaryExec)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		dataVisit, err := helpers.ExecuteQuery(queryVisitExec)
		if err != nil {
			fmt.Println("Error executing query 4:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal execute query",
				Success: false,
			})
		}

		querySummary := `WITH tunais as (
                          SELECT SUM(sum_cash) as total_tunai,
                          {{.QAsUserID}}
                          FROM (
                              SELECT COALESCE(SUM(CASE WHEN p.is_kredit = 0 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END),0) as sum_cash, 1 as param
                              FROM penjualan p
                              LEFT JOIN penjualan_detail pd
                                      ON p.id = pd.penjualan_id
                              LEFT JOIN payment py
                                      ON p.id = py.penjualan_id
                              WHERE py.id IS NULL 
                                      AND p.is_kredit = 0
                                      AND p.user_id = {{.QUserID}}
                                      AND DATE(p.tanggal_penjualan) = {{.QDate}}
                                      
                              UNION ALL

                              SELECT SUM(py.nominal) as sum_cash
                                -- COALESCE(SUM(CASE WHEN p.is_kredit = 0 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END),0) as sum_cash
                                , 2
                              FROM penjualan p
                              JOIN payment py
									ON p.id = py.penjualan_id
                              LEFT JOIN pembayaran_piutang pp 
									ON py.id = pp.payment_id
                              WHERE UPPER(py.tipe) = 'CASH'
                                      AND pp.id IS NULL
                                      AND p.user_id = {{.QUserID}}
                                      AND DATE(p.tanggal_penjualan) = {{.QDate}}
                                      
                              UNION ALL

                              SELECT COALESCE(pp.total_pembayaran,0) as sum_cash, 3
                              FROM pembayaran_piutang pp
                              LEFT JOIN payment py 
                                      ON pp.payment_id = py.id
                              WHERE py.id IS NULL
                                      AND pp.user_id = {{.QUserID}}
                                      AND DATE(pp.tanggal_pembayaran) = {{.QDate}}
                                      
                              UNION ALL

                              SELECT COALESCE(pp.total_pembayaran,0) as sum_cash, 4
                              FROM pembayaran_piutang pp
                              JOIN payment py 
                                      ON pp.payment_id = py.id
                              WHERE UPPER(py.tipe) = 'CASH'
                                      AND pp.user_id = {{.QUserID}}
                                      AND DATE(pp.tanggal_pembayaran) = {{.QDate}}
                      ) sq
                  ), total_call as (
                      SELECT COUNT(kunjungan.id) AS total_call, {{.QAsUserID}}
                          FROM (
                                          SELECT            
                                          DISTINCT ON(k.customer_id) k.id, k.salesman_id
                                          FROM kunjungan k 
                                          -- LEFT JOIN penjualan p 
                                          -- ON p.customer_id = k.customer_id AND DATE(p.tanggal_penjualan) =  DATE({{.QDate}}) 
                                          -- AND (p.salesman_id = $salesmanId OR p.merchandiser_id = $merchandiserId)
                                          WHERE k.user_id = {{.QUserID}}
                                           AND DATE(tanggal_kunjungan) = DATE({{.QDate}}) AND UPPER(k.status_toko) = 'BUKA' AND k.customer_id > 0
                                          ORDER BY k.customer_id, tanggal_kunjungan ASC
                          ) kunjungan
                  ), sales as (
                      SELECT
                              SUM(CASE WHEN pd.harga > 0 THEN pd.jumlah ELSE 0 END) AS total_sales_pack,
                              SUM(pd.jumlah* (pd.harga-pd.diskon)) AS total_sales, 
                              COUNT(DISTINCT p.customer_id) FILTER (WHERE p.customer_id > 0 AND pd.jumlah >= 2 AND pd.harga > 0) AS total_effective_call,
                              SUM(CASE WHEN p.is_kredit = 0 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END) AS total_cash,
                              SUM(CASE WHEN p.is_kredit = 1 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END) AS total_credit,
                              SUM(CASE WHEN pd.harga > 0 AND p.is_kredit = 0 THEN pd.jumlah ELSE 0 END) AS total_cash_pack,
                              SUM(CASE WHEN pd.harga > 0 AND p.is_kredit = 1 THEN pd.jumlah ELSE 0 END) AS total_credit_pack,
                              SUM(CASE WHEN pd.harga = 0 THEN pd.jumlah ELSE 0 END) AS total_program_pack,
                              {{.QAsUserID}}
                              FROM penjualan p
                              JOIN penjualan_detail pd
                              ON p.id = pd.penjualan_id
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = DATE({{.QDate}})
                  ), retur as (
                      SELECT 
                              COALESCE(SUM(pd.jumlah * pd.harga),0) AS total_return,
                              COALESCE(SUM(pd.jumlah),0) AS total_return_pack,
                              {{.QAsUserID}}
                      FROM 
                              pengembalian p
                              JOIN pengembalian_detail pd
                              ON p.id = pd.pengembalian_id
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pengembalian) = DATE({{.QDate}})
                  ),customer_register as (
                      SELECT 
                              COUNT(p.id) as customer_register,
                              {{.QAsUserID}}
                      FROM 
                      customer p
                      WHERE p.user_id_holder = {{.QUserID}} AND DATE(p.dtm_crt) = DATE({{.QDate}})
                  ), timecall as (
                      SELECT 
                              SUM(AGE(p.checkout_at,p.checkin_at))/CASE WHEN count(distinct customer_id) = 0 THEN 1 ELSE count(distinct customer_id) END as time_call,
                              {{.QAsUserID}}
                      FROM 
                              kunjungan_log p
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.checkin_at) = DATE({{.QDate}})
                  ), payments as (
                      SELECT
                          SUM(CASE WHEN UPPER(p.tipe) = 'CASH' THEN nominal ELSE 0 END) AS total_payment_cash,
                          SUM(CASE WHEN UPPER(p.tipe) = 'TRANSFER' THEN nominal ELSE 0 END) AS total_payment_transfer,
                          SUM(CASE WHEN UPPER(p.tipe) = 'CEK' THEN nominal ELSE 0 END) AS total_payment_check,
                          SUM(CASE WHEN UPPER(p.tipe) = 'BILYET GIRO' AND (DATE(tanggal_cair) != DATE({{.QDate}}) OR tanggal_cair IS NULL) THEN nominal ELSE 0 END) AS total_payment_open_bilyet_giro,
                          SUM(CASE WHEN UPPER(p.tipe) = 'BILYET GIRO' AND is_cair = 1 AND DATE(tanggal_cair) = DATE({{.QDate}})  THEN nominal ELSE 0 END) AS total_payment_bilyet_giro,
                          {{.QAsUserID}}
                      FROM
                          payment p
                      JOIN penjualan pj
                      ON pj.id = p.penjualan_id 
                      WHERE p.user_id = {{.QUserID}} AND p.is_verif>-1 --AND pj.is_kredit = 0 
                      AND (
                          DATE(p.tanggal_transaksi) = DATE({{.QDate}}) 
                          OR (DATE(p.tanggal_cair) = DATE({{.QDate}}) AND p.is_cair = 1)
                      )
                  ) , 
                  pembayaran_piutang as (
                              SELECT COALESCE(SUM(p.total_pembayaran),0) as total_payment_and_dp,
                              {{.QAsUserID}}
                              
                              FROM
                              pembayaran_piutang p
                              JOIN pembayaran_piutang_detail pd 
                              ON p.id = pd.pembayaran_piutang_id
                              LEFT JOIN piutang pi
                              ON pi.id = pd.piutang_id
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pembayaran) = DATE({{.QDate}})
                  )

                  SELECT 
                              SUM(datas.total_call) as total_call_buka,
                              SUM(datas.total_effective_call) AS total_effective_call_2pack,
                              SUM(datas.customer_register) as total_register_customer,
                              MAX(datas.time_call) as average_call,
                              SUM(datas.total_sales_pack) AS omzet,
                              SUM(datas.total_cash_pack) AS omzet_tunai,
                              SUM(datas.total_credit_pack) AS omzet_kredit,
                              SUM(datas.total_return_pack) AS retur,
                              --COALESCE(SUM(datas.total_payment_cash) - SUM(datas.total_return),0) AS pembayaran_tunai,
                              COALESCE(MAX(datas.total_tunai) - SUM(datas.total_return),0) AS pembayaran_tunai,
                              COALESCE(SUM(datas.total_payment_transfer),0) AS pembayaran_transfer,
                              COALESCE(SUM(datas.total_payment_check),0) AS cek,
                              COALESCE(SUM(datas.total_payment_open_bilyet_giro),0) AS bilyet_giro_baru,
                              COALESCE(SUM(datas.total_payment_bilyet_giro),0) AS bilyet_giro_cair,
                              SUM(datas.total_cash)+ SUM(datas.total_payment_and_dp)
                                -SUM(datas.total_return)-COALESCE(SUM(datas.total_payment_transfer),0)
                                -COALESCE(SUM(datas.total_payment_check),0)
                                -COALESCE(SUM(datas.total_payment_open_bilyet_giro),0)
                                -COALESCE(SUM(datas.total_payment_bilyet_giro),0) AS total_setoran,
                              MAX(datas.total_tunai)-SUM(datas.total_return) AS total_setoran_tunai
                  FROM ( SELECT tc.*, sales.*, retur.*, cr.*, timecall.*, payments.*, pp.*, tunais.*
                                  FROM public.user s
                                  LEFT JOIN total_call tc
                                      ON s.id = tc.user_id
                                  LEFT JOIN sales
                                      ON s.id = sales.user_id
                                  LEFT JOIN retur
                                      ON s.id = retur.user_id
                                  LEFT JOIN customer_register cr
                                      ON s.id = cr.user_id
                                  LEFT JOIN timecall
                                      ON s.id = timecall.user_id
                                  LEFT JOIN payments
                                      ON s.id = payments.user_id
                                  LEFT JOIN pembayaran_piutang pp
                                      ON s.id = pp.user_id
                                  LEFT JOIN tunais
                                      ON s.id = tunais.user_id
                  ) datas`

		querySummaryExec, err := helpers.PrepareQuery(querySummary, templateEndQuery)

		// fmt.Println(querySummaryExec)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		dataSummary, err := helpers.ExecuteQuery(querySummaryExec)
		if err != nil {
			fmt.Println("Error executing query 5:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal execute query",
				Success: false,
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":             "success",
			"success":             true,
			"data":                data,
			"table_produk":        dataStokProduk,
			"table_produk_header": testHeader,
			"table_item":          dataStokItem,
			"table_item_header":   testHeaderItem,
			"visit":               dataVisit,
			"summary":             dataSummary,
		})
	}
}

func GetUserDailyReport2(c *fiber.Ctx) error {

	date := c.Query("date")
	userId := c.Query("userId")
	isQuery := c.Query("isQuery")
	requestId := c.Query("requestId")
	branchId := helpers.ParamArray(c.Context().QueryArgs().PeekMulti("branchId[]"))

	queries := []string{}
	// keyQuery := []string{}

	if date == "" {
		date = "CURRENT_DATE"
	} else {
		date = "DATE('" + date + "') "
	}

	templateQuery := map[string]interface{}{
		"QUserID": userId,
		"QDate":   date,
	}

	resultsForReturn := make(map[string]interface{})

	queryProduk, err := helpers.PrepareQuery(`SELECT pr.id, pr.code FROM penjualan p
					JOIN penjualan_detail pd
						ON p.id = pd.penjualan_id
					JOIN produk pr
						ON pd.produk_id = pr.id
					WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = {{.QDate}}
					GROUP BY pr.id`, templateQuery)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	produks, err := helpers.ExecuteQuery2(queryProduk, "")
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika mengambil data produk",
			Success: false,
		})
	}

	var qProd, qProdItem, qSelectProd string
	tempArray := make(map[string]int)
	var produkHeader []string
	var produkIds []int

	if len(produks) > 0 {

		for _, result := range produks {
			tempProdId, _ := result.Get("id")
			tempProdCode, _ := result.Get("code")
			qProd = qProd + fmt.Sprintf(`,COALESCE(SUM(case when (pj.tunai > 0 OR pj.kredit > 0) AND pj.produk_id = %v then pj.jumlah else 0 end),0) %v`, tempProdId, tempProdCode)
			qProdItem = qProdItem + fmt.Sprintf(`,0 %v`, tempProdCode)
			qSelectProd = qSelectProd + fmt.Sprintf(`,MAX(sq.%v) as %v`, tempProdCode, tempProdCode)
			tempArray[tempProdCode.(string)] = 0
			produkHeader = append(produkHeader, tempProdCode.(string))
			produkIds = append(produkIds, int(tempProdId.(float64)))
		}

		qProd = qProd[1:]
		qProd = qProd + ","
		qProdItem = qProdItem[1:]
		qProdItem = qProdItem + ","
		qSelectProd = qSelectProd[1:]
		qSelectProd = qSelectProd + ","
	}

	templateEndQuery := map[string]interface{}{
		"QUserID":     userId,
		"QDate":       date,
		"QSelectProd": qSelectProd,
		"QProd":       qProd,
		"QProdItem":   qProdItem,
		"QAsUserID":   userId + " as user_id",
		"QBranchID":   strings.Join(branchId, ","),
	}

	if requestId != "" {

		endQuery := `WITH visits as (
						SELECT id, tanggal_kunjungan, customer_id, user_id, subject_type_id, image_kunjungan
						FROM kunjungan 
						WHERE user_id = {{.QUserID}} 
							AND DATE(tanggal_kunjungan) = {{.QDate}}
						ORDER BY tanggal_kunjungan
					), penjualan_cus as (
						SELECT p.no_nota,
										p.id,
										CASE WHEN COALESCE(p.image_nota_print, p.image_nota) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/nota/'||COALESCE(p.image_nota_print, p.image_nota) ELSE NULL END as image_nota,
										CASE WHEN p.image_bukti_serah IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/nota/'||p.image_bukti_serah ELSE NULL END as image_bukti_serah,
										p.user_id,
										p.customer_id,
										DATE(p.tanggal_penjualan),
										COALESCE(SUM(CASE WHEN p.is_kredit = 0 THEN ((pd.harga - pd.diskon) * pd.jumlah) ELSE 0 END),0) as tunai,
										COALESCE(SUM(CASE WHEN p.is_kredit = 1 THEN ((pd.harga - pd.diskon) * pd.jumlah) ELSE 0 END),0) as kredit,
										pd.produk_id as produk_id,
										COALESCE(SUM(pd.jumlah),0) as jumlah
						FROM penjualan p
						JOIN penjualan_detail pd
							ON p.id = pd.penjualan_id
						WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = {{.QDate}}
						GROUP BY p.id, pd.produk_id
					), pembayaran_cus as (
						SELECT p.customer_id,
										p.user_id,
										COALESCE(SUM(pd.nominal),0) as pembayaran,
										COALESCE((CASE WHEN p.tipe_pelunasan = 0 AND p.pengembalian_id IS NOT NULL then p.total_pembayaran ELSE 0 END),0) AS adjustment,
										COALESCE(SUM((CASE WHEN tipe_pelunasan = 1 then pd.nominal ELSE 0 END)),0) as dp
						FROM pembayaran_piutang p
						JOIN pembayaran_piutang_detail pd
							ON p.id = pd.pembayaran_piutang_id
						LEFT JOIN payment py ON py.id =  p.payment_id AND py.tipe = 'BILYET GIRO'
						WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pembayaran) = {{.QDate}}
						GROUP BY p.id
					), pengembalian_cus as (
						SELECT p.customer_id,
										p.user_id,
										COALESCE(SUM(pd.jumlah * pd.harga),0) as pengembalian
						FROM pengembalian p
						JOIN pengembalian_detail pd
							ON p.id = pd.pengembalian_id
					WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pengembalian) = {{.QDate}}
						GROUP BY p.id
					), payments as (
						SELECT UPPER((CASE WHEN tipe = 'BILYET GIRO' THEN 'OPEN ' ELSE '' END)||tipe) AS payment_tipe, 
								COALESCE(nominal,0) AS payment_nominal,
										p.user_id,
										p.customer_id,
										CASE WHEN p.bukti_bayar = '' THEN NULL ELSE p.bukti_bayar END as payment_image,
																				p.penjualan_id as penjualan_id
								FROM payment p
								WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_transaksi) = {{.QDate}} 

						UNION

						SELECT 'CASH' as payment_tipe,
								p.total_penjualan as payment_nominal,
								p.user_id,
								p.customer_id,
								null as payment_image,
								p.id as penjualan_id
						FROM penjualan p
						LEFT JOIN payment py
								ON p.id = py.penjualan_id
								AND py.id IS NULL 
						WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_penjualan) = {{.QDate}} AND p.is_kredit = 0
						GROUP BY p.id
					), kunjunganlogs as (
						SELECT DISTINCT ON (customer_id) * 
						FROM kunjungan_log WHERE user_id = {{.QUserID}} AND date(checkin_at) = {{.QDate}}
						ORDER BY customer_id, checkin_at
					), transactions as (
						SELECT t.user_id,
										t.subject_type_id,
										t.customer_id,
											JSONB_BUILD_OBJECT(
												'name', i.name,
												'quantity', COALESCE(SUM(td.qty),0)
										) as items,
										i.id as item_id,
										i.name as item_name,
										COALESCE(SUM(td.qty),0) as item_jumlah,
										CASE WHEN ic.id = 191 AND t.transaction_type_id <> 1 
											THEN JSONB_BUILD_OBJECT(
												'photo_before', CASE WHEN MAX(t.photos_before) = '{}' OR MAX(t.photos_before) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos_before)) image
																	) END,
												'photo_after', CASE WHEN MAX(t.photos_after) = '{}' OR MAX(t.photos_after) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos_after)) image
																	) END
											)
											ELSE NULL
										END as sampling,
										CASE WHEN ic.id = 171 AND t.transaction_type_id <> 1 
											THEN JSONB_BUILD_OBJECT(
												'photo_before', CASE WHEN MAX(t.photos_before) = '{}' OR MAX(t.photos_before) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos_before)) image
																		) END,
												'photo_after', CASE WHEN MAX(t.photos_after) = '{}' OR MAX(t.photos_after) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos_after)) image
																		) END
											)
											ELSE NULL
										END as posm,
										CASE WHEN ic.id = 161 AND t.transaction_type_id <> 1 
											THEN JSONB_BUILD_OBJECT(
												'photo_before', CASE WHEN MAX(t.photos_before) = '{}' OR MAX(t.photos_before) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' || image 
																			FROM unnest(MAX(t.photos_before)) image
																		) END,
												'photo_after', CASE WHEN MAX(t.photos_after) = '{}' OR MAX(t.photos_after) IS NULL 
																	THEN ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' || image 
																			FROM unnest(MAX(t.photos_after)) image
																		) END
											)
											ELSE NULL
										END as merchandise
										
						FROM md.transaction t
						JOIN md.transaction_detail td
							ON t.id = td.transaction_id
						JOIN md.item i
							ON td.item_id = i.id
						JOIN md.item_category ic
							ON i.category_id = ic.id
						WHERE t.user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}}
						GROUP BY t.user_id, t.customer_id, ic.id, t.transaction_type_id, t.subject_type_id, i.id
					)

					SELECT st.name as subject_type, 
							sq.customer, 
							MAX(invoice) as invoice,  
							{{.QSelectProd}}
							--JSONB_AGG(sq.items) FILTER (WHERE sq.items IS NOT NULL) as items,
							JSONB_AGG(
									JSONB_BUILD_OBJECT(
										'name', sq.item_name,
										'quantity', sq.item_jumlah
								)
							) FILTER (WHERE sq.item_name IS NOT NULL) as items,
							COALESCE(SUM(sq.tunai),0) as tunai,
							COALESCE(SUM(sq.kredit),0) as kredit,
							COALESCE(SUM(sq.pengembalian), 0) as retur,
							COALESCE(SUM(sq.pembayaran),0) as pembayaran,
							COALESCE(SUM(sq.adjustment),0) as adjustment,
							COALESCE(SUM(sq.tunai),0) + 
							COALESCE(SUM(sq.pembayaran),0) + 
							COALESCE(SUM(sq.dp),0) - 
							COALESCE(SUM(sq.pengembalian),0) + 
							COALESCE(SUM(sq.adjustment),0) as setoran,
							JSONB_BUILD_OBJECT(
										'checkin', MIN(sq.checkin),
										'checkout', MAX(sq.checkout),
										'durasi', MAX(sq.durasi)
							) as time_info
							/*,JSONB_AGG(sq.payment_information) FILTER (WHERE sq.payment_information IS NOT NULL)->0 as payment_information,
							JSONB_AGG(sq.nota) FILTER (WHERE sq.nota IS NOT NULL)->0 as nota,
							JSONB_AGG(sq.penyerahan_produk) FILTER (WHERE sq.penyerahan_produk IS NOT NULL)->0 as penyerahan_produk,
							JSONB_AGG(sq.posm) FILTER (WHERE sq.posm IS NOT NULL)->0 as posm,
							JSONB_AGG(sq.merchandise) FILTER (WHERE sq.merchandise IS NOT NULL)->0 as merchandise,
							JSONB_AGG(sq.sampling) FILTER (WHERE sq.sampling IS NOT NULL)->0 as sampling*/
					FROM (
						SELECT v.subject_type_id,
								--COALESCE(c.id,v.customer_id, pj.customer_id, kl.customer_id) as customer_id,
								JSONB_BUILD_OBJECT(
										'name', COALESCE(c.name, 'End User'),
										'contact', COALESCE(c.outlet_name, '-'),
										'location', COALESCE(c.kelurahan, '-'),
										'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
								) as customer,
								pj.no_nota as invoice,
								{{.QProd}}
								null as item_name,
								0 as item_jumlah,
								COALESCE(SUM(pj.tunai),0) as tunai,
								COALESCE(SUM(pj.kredit),0) as kredit,
								0 as pembayaran,
								0 as pengembalian,
								0 as adjustment,
								0 as dp,
								NULL as payment_information,
								--MIN(kl.checkin_at) as param,
								to_char(COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan)) , 'HH24:MI') as checkin,
								to_char(COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)), 'HH24:MI') as checkout,
								COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)) - COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan)) as durasi
								/*,
								pj.image_nota as nota,
								pj.image_bukti_serah as penyerahan_produk,
								null::jsonb as posm,
								null::jsonb as merchandise,
								null::jsonb as sampling*/
							FROM visits v
							FULL JOIN penjualan_cus pj
								ON v.user_id = pj.user_id
								AND v.customer_id = pj.customer_id
							LEFT JOIN kunjunganlogs kl
								ON v.user_id = kl.user_id
								AND v.customer_id = kl.customer_id
							LEFT JOIN customer c
								ON COALESCE(v.customer_id, pj.customer_id, kl.customer_id) = c.id
							GROUP BY v.subject_type_id, v.image_kunjungan, c.id, COALESCE(c.id,v.customer_id, pj.customer_id, kl.customer_id), pj.no_nota, pj.image_nota, pj.image_bukti_serah
							
							UNION

							SELECT v.subject_type_id,
									--COALESCE(c.id,v.customer_id, pem.customer_id, kl.customer_id) as customer_id,
									JSONB_BUILD_OBJECT(
											'name', COALESCE(c.name, 'End User'),
											'contact', COALESCE(c.outlet_name, '-'),
											'location', COALESCE(c.kelurahan, '-'),
											'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
									) as customer,
									NULL as invoice,
									{{.QProdItem}}
									null as item_name,
									0 as item_jumlah,
									0 as tunai,
									0 as kredit,
									COALESCE(SUM(pem.pembayaran),0) as pembayaran,
									0 as pengembalian,
									COALESCE(SUM(pem.adjustment),0) as adjustment,
									COALESCE(SUM(pem.dp),0) as dp,
									/*JSONB_BUILD_OBJECT(
									'cash', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'CASH'), 0) - COALESCE(SUM(pg.pengembalian),0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'CASH' AND py.payment_image <> '')
										),
									'transfer', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'TRANSFER'), 0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'TRANSFER' AND py.payment_image <> '')
										),
									'bilyet_giro_cair', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'BILYET GIRO'), 0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'BILYET GIRO' AND py.payment_image <> '')
										),
									'bilyet_giro_open', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'OPEN BILYET GIRO'), 0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'OPEN BILYET GIRO' AND py.payment_image <> '')
										)
									) as payment_information,*/
									null as payment_information,
									--MIN(kl.checkin_at) as param,
									to_char(COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan)) , 'HH24:MI') as checkin,
									to_char(COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)), 'HH24:MI') as checkout,
									COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)) - COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan)) as durasi
									/*,
									null as nota,
									null as penyerahan_produk,
									null::jsonb as posm,
									null::jsonb as merchandise,
									null::jsonb as sampling*/
							FROM visits v
							FULL JOIN pembayaran_cus pem
								ON v.user_id = pem.user_id
								AND v.customer_id = pem.customer_id
							LEFT JOIN kunjunganlogs kl
								ON v.user_id = kl.user_id
								AND v.customer_id = kl.customer_id
							LEFT JOIN customer c
								ON COALESCE(v.customer_id, pem.customer_id, kl.customer_id) = c.id
							GROUP BY v.subject_type_id, v.image_kunjungan, c.id, COALESCE(c.id,v.customer_id, pem.customer_id, kl.customer_id)

							UNION

							SELECT v.subject_type_id,
									--COALESCE(c.id,v.customer_id, pg.customer_id, py.customer_id, kl.customer_id) as customer_id,
									JSONB_BUILD_OBJECT(
											'name', COALESCE(c.name, 'End User'),
											'contact', COALESCE(c.outlet_name, '-'),
											'location', COALESCE(c.kelurahan, '-'),
											'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
									) as customer,
									NULL as invoice,
									{{.QProdItem}}
									null as item_name,
									0 as item_jumlah,
									0 as tunai,
									0 as kredit,
									0 as pembayaran,
									COALESCE(SUM(pg.pengembalian),0) as pengembalian,
									0 as adjustment,
									0 as dp,
									/*JSONB_BUILD_OBJECT(
									'cash', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'CASH'), 0) - COALESCE(SUM(pg.pengembalian),0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'CASH' AND py.payment_image <> '')
										),
									'transfer', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'TRANSFER'), 0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'TRANSFER' AND py.payment_image <> '')
										),
									'bilyet_giro_cair', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'BILYET GIRO'), 0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'BILYET GIRO' AND py.payment_image <> '')
										),
									'bilyet_giro_open', JSONB_BUILD_OBJECT(
											'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'OPEN BILYET GIRO'), 0),
											'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'OPEN BILYET GIRO' AND py.payment_image <> '')
										)
									) as payment_information,*/
									null as payment_information,
									--MIN(kl.checkin_at) as param,
									to_char(COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan)) , 'HH24:MI') as checkin,
									to_char(COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)), 'HH24:MI') as checkout,
									COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)) - COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan)) as durasi
									/*,
									null as nota,
									null as penyerahan_produk,
									null::jsonb as posm,
									null::jsonb as merchandise,
									null::jsonb as sampling*/
							FROM visits v
							FULL JOIN pengembalian_cus pg
								ON v.user_id = pg.user_id
								AND v.customer_id = pg.customer_id
							FULL JOIN payments py
								ON v.user_id = py.user_id
								AND v.customer_id = py.customer_id
							LEFT JOIN kunjunganlogs kl
								ON v.user_id = kl.user_id
								AND v.customer_id = kl.customer_id
							LEFT JOIN customer c
								ON COALESCE(v.customer_id, pg.customer_id, py.customer_id, kl.customer_id) = c.id
							GROUP BY v.subject_type_id, v.image_kunjungan, c.id, COALESCE(c.id,v.customer_id, pg.customer_id, py.customer_id, kl.customer_id)
							
							UNION
							
							SELECT sq.subject_type_id,
									sq.customer,
									null as invoice,
									{{.QProdItem}}
									sq.item_name as item_name,
									SUM(sq.item_jumlah) as item_jumlah,
									SUM(sq.tunai) as tunai,
									SUM(sq.kredit) as kredit,
									SUM(sq.pembayaran) as pembayaran,
									SUM(sq.pengembalian) as pengembalian,
									SUM(sq.adjustment) as adjustment,
									SUM(sq.dp) as dp,
									NULL as payment_information,
									MIN(sq.checkin) as checkin,
									MAX(sq.checkout) as checkout,
									MAX(sq.durasi) as durasi
								FROM (
										SELECT COALESCE(v.subject_type_id, t.subject_type_id) as subject_type_id,
											COALESCE(c.id, v.customer_id, t.customer_id, kl.customer_id) as customer_id,
											JSONB_BUILD_OBJECT(
													'name', COALESCE(c.name, 'End User'),
													'contact', COALESCE(c.outlet_name, '-'),
													'location', COALESCE(c.kelurahan, '-'),
													'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
											) as customer,
											NULL as invoice,
											t.items as items,
											t.item_id,
											t.item_name as item_name,
											SUM(t.item_jumlah) as item_jumlah,
											0 as tunai,
											0 as kredit,
											0 as pembayaran,
											0 as pengembalian,
											0 as adjustment,
											0 as dp,
											NULL as payment_information,
											MIN(kl.checkin_at) as param,
											to_char(COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan)) , 'HH24:MI') as checkin,
											to_char(COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)), 'HH24:MI') as checkout,
											COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)) - COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan)) as durasi
											/*,JSONB_BUILD_OBJECT(
												'checkin', to_char(MIN(kl.checkin_at) , 'HH24:MI'),
												'checkout', to_char(MAX(kl.checkout_at), 'HH24:MI'),
												'durasi', MAX(kl.checkout_at) - MIN(kl.checkin_at)
											) as time_info,
											null as nota,
											null as penyerahan_produk,
											t.posm as posm,
											t.merchandise as merchandise,
											t.sampling as sampling*/
									FROM visits v
									FULL JOIN transactions t
										ON v.user_id = t.user_id
										AND v.customer_id = t.customer_id
									LEFT JOIN kunjunganlogs kl
										ON v.user_id = kl.user_id
										AND v.customer_id = kl.customer_id
									LEFT JOIN customer c
										ON COALESCE(v.customer_id, t.customer_id, kl.customer_id) = c.id
									GROUP BY COALESCE(v.subject_type_id, t.subject_type_id), v.image_kunjungan, c.id, COALESCE(c.id, v.customer_id, t.customer_id, kl.customer_id), t.items, t.posm, t.merchandise, t.sampling, t.item_id, t.item_name
								) sq
								GROUP BY sq.subject_type_id, sq.customer, sq.item_id, sq.item_name
					) sq
					JOIN md.subject_type st
						ON COALESCE(sq.subject_type_id, 3) = st.id
					GROUP BY st.id, sq.customer
					ORDER BY st.id`
		finalQuery, err := helpers.PrepareQuery(endQuery, templateEndQuery)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		if isQuery != "" && isQuery == "1" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"query": finalQuery,
			})
		}

		queries = append(queries, finalQuery)

		// fmt.Println(finalQuery)

		// data, err := helpers.NewExecuteQuery(finalQuery)
		// if err != nil {
		// 	fmt.Println("Error executing query:", err)
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal execute query",
		// 		Success: false,
		// 	})
		// }

		getSubject, err := helpers.ExecuteQuery(fmt.Sprintf(`SELECT to_char(%s, 'DD Mon YYYY') as tanggal,
															JSONB_BUILD_OBJECT(
																'name', INITCAP(full_name),
																'type', INITCAP('user')
															) as requester
														FROM public.user WHERE id = %s`, date, userId))
		if err != nil {
			fmt.Println("Error executing query request:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal execute query request",
				Success: false,
			})
		}

		getRequester, err := helpers.ExecuteQuery(fmt.Sprintf(`SELECT to_char(%s, 'DD Mon YYYY') as tanggal,
															JSONB_BUILD_OBJECT(
																'name', INITCAP(full_name),
																'type', INITCAP('user')
															) as requester
														FROM public.user WHERE id = %s`, date, requestId))
		if err != nil {
			fmt.Println("Error executing query request:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal execute query request",
				Success: false,
			})
		}

		implodedProdukHeader := helpers.SplitToString(produkIds, ",")

		templateEndQuery["QProdukIds"] = implodedProdukHeader

		queryStokProduk := `SELECT sq.* 
							FROM (
							WITH all_products as (SELECT * FROM produk WHERE id IN ({{.QProdukIds}}) ORDER BY produk.order)
							
							SELECT sq.* FROM (
								SELECT p.id, p.order,
										MAX(p.code) as code, 
										COALESCE(SUM(ss.stok_awal),0) - COALESCE(SUM(ssr.jumlah),0) as packs, 
										'stok awal' as tag,
										1 as orders
								FROM all_products p
								LEFT JOIN stok_salesman ss
								ON ss.produk_id = p.id
								AND ss.user_id = {{.QUserID}}
								AND DATE(ss.tanggal_stok) = {{.QDate}} 
								AND ss.produk_id IN ({{.QProdukIds}})
								LEFT JOIN
								(
								SELECT produk_id, condition, pita, SUM(COALESCE(jumlah,0)) AS jumlah 
								FROM stok_salesman_riwayat 
								WHERE is_validate = 1 
									AND DATE(tanggal_riwayat) = {{.QDate}} 
									AND user_id = {{.QUserID}} AND aksi='ORDER'
								GROUP BY produk_id, condition, pita
								) ssr 
										ON ssr.produk_id = ss.produk_id 
								AND ssr.condition = ss.condition 
								AND ssr.pita = ss.pita
								GROUP BY p.id,ss.produk_id, p.order
								ORDER BY p.order
							) sq
							
							UNION
							
							SELECT sq.* FROM (
								SELECT p.id, p.order,
										MAX(p.code) as code, 
										COALESCE(SUM(jumlah),0) as packs, 
										'order' as tag,
										2
								FROM all_products p
								LEFT JOIN stok_salesman_riwayat ssr
								ON ssr.produk_id = p.id
								AND ssr.is_validate = 1
								AND ssr.user_id = {{.QUserID}} 
								AND DATE(tanggal_riwayat) = {{.QDate}} 
								AND produk_id IN ({{.QProdukIds}}) 
								AND UPPER(aksi) = 'ORDER'
								GROUP BY p.id,produk_id, p.order
								ORDER BY p.order
							) sq
							
							UNION
							
							SELECT sq.* FROM (
								SELECT p.id, p.order, 
										MAX(p.code) as code, 
										COALESCE(SUM(sq.total_sales_pack),0) as packs, 
										'penjualan' as tag,
										3
								FROM all_products p
								LEFT JOIN (
									SELECT p.customer_id, 
											pd.produk_id, 
											SUM(pd.jumlah* (pd.harga-pd.diskon)) AS total_sales, 
											SUM(CASE WHEN (pd.harga - pd.diskon) <> 0 THEN pd.jumlah ELSE 0 END) AS total_sales_pack
										FROM penjualan p
										JOIN penjualan_detail pd
										ON p.id = pd.penjualan_id
										WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = {{.QDate}} AND pd.produk_id IN ({{.QProdukIds}})
										GROUP BY p.customer_id, pd.produk_id
									) sq
								ON p.id = sq.produk_id
								GROUP BY p.id,produk_id, p.order
								ORDER BY p.order
							) sq

							UNION
							
							SELECT sq.* FROM (
								SELECT p.id, p.order, 
										MAX(p.code) as code, 
										COALESCE(SUM(sq.total_sales_pack),0) as packs, 
										'program' as tag,
										4
								FROM all_products p
								LEFT JOIN (
									SELECT p.customer_id, 
											pd.produk_id, 
											SUM(pd.jumlah* (pd.harga-pd.diskon)) AS total_sales, 
											SUM(CASE WHEN (pd.harga - pd.diskon) = 0 THEN pd.jumlah ELSE 0 END) AS total_sales_pack
										FROM penjualan p
										JOIN penjualan_detail pd
										ON p.id = pd.penjualan_id
										WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = {{.QDate}} AND pd.produk_id IN ({{.QProdukIds}})
										GROUP BY p.customer_id, pd.produk_id
									) sq
								ON p.id = sq.produk_id
								GROUP BY p.id,produk_id, p.order
								ORDER BY p.order
							) sq
							
							-- UNION
							
							-- SELECT sq.* FROM (
							--   SELECT p.id, p.order, 
							--           MAX(p.code) as code, 
							--           COALESCE(SUM(pjd.jumlah) FILTER (WHERE pjd.produk_id = p.id),0) as packs, 
							--           'program' as tag,
							--           5
							--   FROM penjualan pj
							--   JOIN penjualan_detail pjd
							--     ON pjd.harga = 0
							--     AND pjd.penjualan_id = pj.id
							--     AND (salesman_id = $salesmanId OR p.merchandiser_id = $merchandiserId) 
							--     AND DATE(tanggal_penjualan) = {{.QDate}}
							--   CROSS JOIN all_products p
							--   GROUP BY p.id,produk_id, p.order
							--   ORDER BY p.order
							-- ) sq
							
							UNION
							
							SELECT sq.* FROM (
								SELECT p.id, p.order, 
										MAX(p.code) as code, 
										COALESCE(SUM(pgd.jumlah) FILTER (WHERE DATE(tanggal_pengembalian) = {{.QDate}}),0) as packs, 
										'retur customer' as tag,
										6
								FROM all_products p
								LEFT JOIN pengembalian_detail pgd
								ON pgd.produk_id = p.id
								LEFT JOIN pengembalian pg
								ON pgd.pengembalian_id = pg.id
								AND pg.user_id = {{.QUserID}}
								AND DATE(tanggal_pengembalian) = {{.QDate}} 
								AND produk_id IN ({{.QProdukIds}})
								GROUP BY p.id,produk_id, p.order
								ORDER BY p.order
							) sq
							
							UNION
							
							SELECT sq.* FROM (
								SELECT p.id, p.order, 
										MAX(p.code) as code, 
										COALESCE(SUM(jumlah),0) as packs, 
										'retur gudang' as tag,
										7
								FROM all_products p
								LEFT JOIN stok_salesman_riwayat ssr
								ON ssr.produk_id = p.id
								AND (ssr.user_id = {{.QUserID}}) 
								AND DATE(tanggal_riwayat) = {{.QDate}} 
								AND produk_id IN ({{.QProdukIds}}) 
								AND UPPER(aksi) = 'RETUR'
								AND ssr.is_validate = 1
								GROUP BY p.id,produk_id, p.order
								ORDER BY p.order
							) sq
							
							UNION
							
							SELECT sq.* FROM (
								SELECT p.id, p.order, 
										MAX(p.code) as code, 
										COALESCE(SUM(stok_akhir),0) as packs, 
										'stok akhir' as tag,
										8
								FROM all_products p
								LEFT JOIN stok_salesman ss
								ON ss.produk_id = p.id
								AND (ss.user_id = {{.QUserID}}) 
								AND DATE(tanggal_stok) = {{.QDate}} 
								AND produk_id IN ({{.QProdukIds}})
								GROUP BY p.id,produk_id, p.order
								ORDER BY p.order
							) sq

							UNION
							
							SELECT sq.* FROM (
									SELECT p.id, p.order, 
									MAX(p.code) as code, 
									COUNT(DISTINCT sq.customer_id) FILTER (WHERE total_sales > 0 AND customer_id > 0 AND total_sales_pack >=2) AS total_effective_call,
									'ec brand 2 pack' as tag,
									9999
									FROM all_products p
									LEFT JOIN (
										SELECT p.customer_id, 
												pd.produk_id, 
												SUM(pd.jumlah* (pd.harga-pd.diskon)) AS total_sales, 
												SUM(CASE WHEN pd.harga > 0 THEN pd.jumlah ELSE 0 END) AS total_sales_pack
											FROM penjualan p
											JOIN penjualan_detail pd
											ON p.id = pd.penjualan_id
											WHERE (p.user_id = {{.QUserID}}) AND DATE(p.tanggal_penjualan) = {{.QDate}} AND pd.produk_id IN ({{.QProdukIds}})
											GROUP BY p.customer_id, pd.produk_id
										) sq
									ON p.id = sq.produk_id
									GROUP BY p.id,produk_id, p.order
									ORDER BY p.order
								) sq 
							
								-- SELECT id, all_products.order, MAX(code) as code, 0, 'ec_brand_2_pack' as tag, 99999 FROM all_products GROUP BY id, all_products.order
							-- ) sq
							) sq
							ORDER BY sq.orders, sq.order`

		queryStokProdukExec, err := helpers.PrepareQuery(queryStokProduk, templateEndQuery)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		// fmt.Println(queryStokProdukExec)
		// restructuredArray := make(map[string]map[string]interface{})
		testHeader := []string{"Stok Awal", "Order", "Penjualan", "Program", "Retur Customer", "Retur Gudang", "Stok Akhir"}

		queryStokItem := `WITH ordersreturs as (
								SELECT user_id, item_id, aksi, SUM(jumlah) as jumlah
								FROM md.stok_merchandiser_riwayat
								WHERE user_id = {{.QUserID}} AND DATE(tanggal_riwayat) = {{.QDate}}
								GROUP BY user_id, item_id, aksi
							), stoks as (
								SELECT user_id, item_id, stok_awal, stok_akhir
								FROM md.stok_merchandiser
								WHERE user_id = {{.QUserID}} AND DATE(tanggal_stok) = {{.QDate}}
							), transactions as (
								SELECT t.user_id, td.item_id, SUM(qty) as jumlah
								FROM md.transaction t
								JOIN md.transaction_detail td
									ON t.id = td.transaction_id
								WHERE t.user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}}
								GROUP BY t.user_id, td.item_id
							)

							SELECT CASE WHEN pb.id IS NOT NULL THEN CONCAT(i.name, ' - ', pb.name) ELSE i.name END as item_name,
									JSONB_BUILD_OBJECT(
										'Stok Awal', MAX(s.stok_awal) - SUM(CASE WHEN o.aksi = 'ORDER' THEN o.jumlah ELSE 0 END),
										'Order', SUM(CASE WHEN o.aksi = 'ORDER' THEN o.jumlah ELSE 0 END),
										'Transaksi', MAX(COALESCE(t.jumlah,0)),
										'Retur Gudang', SUM(CASE WHEN o.aksi = 'RETUR' THEN o.jumlah ELSE 0 END),
										'Stok Akhir', MAX(s.stok_akhir)
									) as datas
							FROM stoks s
							FULL JOIN ordersreturs o
								ON s.item_id = o.item_id
							FULL JOIN transactions t
								ON s.item_id = t.item_id
							FULL JOIN md.merchandiser m
								ON m.user_id = COALESCE(s.user_id, o.user_id, t.user_id)
							LEFT JOIN md.item i
								ON s.item_id = i.id
							LEFT JOIN produk_brand pb
								ON i.brand_id = pb.id
							WHERE m.user_id = {{.QUserID}}
							GROUP BY i.id, pb.id
							ORDER BY i.name`

		queryStokItemExec, err := helpers.PrepareQuery(queryStokItem, templateEndQuery)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		queries = append(queries, queryStokItemExec)

		// dataStokItem, err := helpers.ExecuteQuery(queryStokItemExec)
		// if err != nil {
		// 	fmt.Println("Error executing query 2:", err)
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal execute query",
		// 		Success: false,
		// 	})
		// }

		testHeaderItem := []string{"Stok Awal", "Order", "Transaksi", "Retur Gudang", "Stok Akhir"}

		querySummary := `WITH tunais as (
                          SELECT SUM(sum_cash) as total_tunai,
                          {{.QAsUserID}}
                          FROM (
                              SELECT COALESCE(SUM(CASE WHEN p.is_kredit = 0 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END),0) as sum_cash, 1 as param
                              FROM penjualan p
                              LEFT JOIN penjualan_detail pd
                                      ON p.id = pd.penjualan_id
                              LEFT JOIN payment py
                                      ON p.id = py.penjualan_id
                              WHERE py.id IS NULL 
                                      AND p.is_kredit = 0
                                      AND p.user_id = {{.QUserID}}
                                      AND DATE(p.tanggal_penjualan) = {{.QDate}}
                                      
                              UNION ALL

                              SELECT SUM(py.nominal) as sum_cash
                                -- COALESCE(SUM(CASE WHEN p.is_kredit = 0 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END),0) as sum_cash
                                , 2
                              FROM penjualan p
                              JOIN payment py
									ON p.id = py.penjualan_id
                              LEFT JOIN pembayaran_piutang pp 
									ON py.id = pp.payment_id
                              WHERE UPPER(py.tipe) = 'CASH'
                                      AND pp.id IS NULL
                                      AND p.user_id = {{.QUserID}}
                                      AND DATE(p.tanggal_penjualan) = {{.QDate}}
                                      
                              UNION ALL

                              SELECT COALESCE(pp.total_pembayaran,0) as sum_cash, 3
                              FROM pembayaran_piutang pp
                              LEFT JOIN payment py 
                                      ON pp.payment_id = py.id
                              WHERE py.id IS NULL
                                      AND pp.user_id = {{.QUserID}}
                                      AND DATE(pp.tanggal_pembayaran) = {{.QDate}}
                                      
                              UNION ALL

                              SELECT COALESCE(pp.total_pembayaran,0) as sum_cash, 4
                              FROM pembayaran_piutang pp
                              JOIN payment py 
                                      ON pp.payment_id = py.id
                              WHERE UPPER(py.tipe) = 'CASH'
                                      AND pp.user_id = {{.QUserID}}
                                      AND DATE(pp.tanggal_pembayaran) = {{.QDate}}
                      ) sq
                  ), total_call as (
                      SELECT COUNT(kunjungan.id) AS total_call, {{.QAsUserID}}
                          FROM (
                                          SELECT            
                                          DISTINCT ON(k.customer_id) k.id, k.salesman_id
                                          FROM kunjungan k 
                                          -- LEFT JOIN penjualan p 
                                          -- ON p.customer_id = k.customer_id AND DATE(p.tanggal_penjualan) =  DATE({{.QDate}}) 
                                          -- AND (p.salesman_id = $salesmanId OR p.merchandiser_id = $merchandiserId)
                                          WHERE k.user_id = {{.QUserID}}
                                           AND DATE(tanggal_kunjungan) = DATE({{.QDate}}) AND UPPER(k.status_toko) = 'BUKA' AND k.customer_id > 0
                                          ORDER BY k.customer_id, tanggal_kunjungan ASC
                          ) kunjungan
                  ), sales as (
                      SELECT
                              SUM(CASE WHEN pd.harga > 0 THEN pd.jumlah ELSE 0 END) AS total_sales_pack,
                              SUM(pd.jumlah* (pd.harga-pd.diskon)) AS total_sales, 
                              COUNT(DISTINCT p.customer_id) FILTER (WHERE p.customer_id > 0 AND pd.jumlah >= 2 AND pd.harga > 0) AS total_effective_call,
                              SUM(CASE WHEN p.is_kredit = 0 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END) AS total_cash,
                              SUM(CASE WHEN p.is_kredit = 1 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END) AS total_credit,
                              SUM(CASE WHEN pd.harga > 0 AND p.is_kredit = 0 THEN pd.jumlah ELSE 0 END) AS total_cash_pack,
                              SUM(CASE WHEN pd.harga > 0 AND p.is_kredit = 1 THEN pd.jumlah ELSE 0 END) AS total_credit_pack,
                              SUM(CASE WHEN pd.harga = 0 THEN pd.jumlah ELSE 0 END) AS total_program_pack,
                              {{.QAsUserID}}
                              FROM penjualan p
                              JOIN penjualan_detail pd
                              ON p.id = pd.penjualan_id
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = DATE({{.QDate}})
                  ), retur as (
                      SELECT 
                              COALESCE(SUM(pd.jumlah * pd.harga),0) AS total_return,
                              COALESCE(SUM(pd.jumlah),0) AS total_return_pack,
                              {{.QAsUserID}}
                      FROM 
                              pengembalian p
                              JOIN pengembalian_detail pd
                              ON p.id = pd.pengembalian_id
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pengembalian) = DATE({{.QDate}})
                  ),customer_register as (
                      SELECT 
                              COUNT(p.id) as customer_register,
                              {{.QAsUserID}}
                      FROM 
                      customer p
                      WHERE p.user_id_holder = {{.QUserID}} AND DATE(p.dtm_crt) = DATE({{.QDate}})
                  ), timecall as (
                      SELECT 
                              SUM(AGE(p.checkout_at,p.checkin_at))/CASE WHEN count(distinct customer_id) = 0 THEN 1 ELSE count(distinct customer_id) END as time_call,
                              {{.QAsUserID}}
                      FROM 
                              kunjungan_log p
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.checkin_at) = DATE({{.QDate}})
                  ), payments as (
                      SELECT
                          SUM(CASE WHEN UPPER(p.tipe) = 'CASH' THEN nominal ELSE 0 END) AS total_payment_cash,
                          SUM(CASE WHEN UPPER(p.tipe) = 'TRANSFER' THEN nominal ELSE 0 END) AS total_payment_transfer,
                          SUM(CASE WHEN UPPER(p.tipe) = 'CEK' THEN nominal ELSE 0 END) AS total_payment_check,
                          SUM(CASE WHEN UPPER(p.tipe) = 'BILYET GIRO' AND (DATE(tanggal_cair) != DATE({{.QDate}}) OR tanggal_cair IS NULL) THEN nominal ELSE 0 END) AS total_payment_open_bilyet_giro,
                          SUM(CASE WHEN UPPER(p.tipe) = 'BILYET GIRO' AND is_cair = 1 AND DATE(tanggal_cair) = DATE({{.QDate}})  THEN nominal ELSE 0 END) AS total_payment_bilyet_giro,
                          {{.QAsUserID}}
                      FROM
                          payment p
                      JOIN penjualan pj
                      ON pj.id = p.penjualan_id 
                      WHERE p.user_id = {{.QUserID}} AND p.is_verif>-1 --AND pj.is_kredit = 0 
                      AND (
                          DATE(p.tanggal_transaksi) = DATE({{.QDate}}) 
                          OR (DATE(p.tanggal_cair) = DATE({{.QDate}}) AND p.is_cair = 1)
                      )
                  ) , 
                  pembayaran_piutang as (
                              SELECT COALESCE(SUM(p.total_pembayaran),0) as total_payment_and_dp,
                              {{.QAsUserID}}
                              
                              FROM
                              pembayaran_piutang p
                              JOIN pembayaran_piutang_detail pd 
                              ON p.id = pd.pembayaran_piutang_id
                              LEFT JOIN piutang pi
                              ON pi.id = pd.piutang_id
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pembayaran) = DATE({{.QDate}})
                  )

                  SELECT 
                              SUM(datas.total_call) as total_call_buka,
                              SUM(datas.total_effective_call) AS total_effective_call_2pack,
                              SUM(datas.customer_register) as total_register_customer,
                              MAX(datas.time_call) as average_call,
                              SUM(datas.total_sales_pack) AS omzet,
                              SUM(datas.total_cash_pack) AS omzet_tunai,
                              SUM(datas.total_credit_pack) AS omzet_kredit,
                              SUM(datas.total_return_pack) AS retur,
                              --COALESCE(SUM(datas.total_payment_cash) - SUM(datas.total_return),0) AS pembayaran_tunai,
                              COALESCE(MAX(datas.total_tunai) - SUM(datas.total_return),0) AS pembayaran_tunai,
                              COALESCE(SUM(datas.total_payment_transfer),0) AS pembayaran_transfer,
                              COALESCE(SUM(datas.total_payment_check),0) AS cek,
                              COALESCE(SUM(datas.total_payment_open_bilyet_giro),0) AS bilyet_giro_baru,
                              COALESCE(SUM(datas.total_payment_bilyet_giro),0) AS bilyet_giro_cair,
                              SUM(datas.total_cash)+ SUM(datas.total_payment_and_dp)
                                -SUM(datas.total_return)-COALESCE(SUM(datas.total_payment_transfer),0)
                                -COALESCE(SUM(datas.total_payment_check),0)
                                -COALESCE(SUM(datas.total_payment_open_bilyet_giro),0)
                                -COALESCE(SUM(datas.total_payment_bilyet_giro),0) AS total_setoran,
                              MAX(datas.total_tunai)-SUM(datas.total_return) AS total_setoran_tunai
                  FROM ( SELECT tc.*, sales.*, retur.*, cr.*, timecall.*, payments.*, pp.*, tunais.*
                                  FROM public.user s
                                  LEFT JOIN total_call tc
                                      ON s.id = tc.user_id
                                  LEFT JOIN sales
                                      ON s.id = sales.user_id
                                  LEFT JOIN retur
                                      ON s.id = retur.user_id
                                  LEFT JOIN customer_register cr
                                      ON s.id = cr.user_id
                                  LEFT JOIN timecall
                                      ON s.id = timecall.user_id
                                  LEFT JOIN payments
                                      ON s.id = payments.user_id
                                  LEFT JOIN pembayaran_piutang pp
                                      ON s.id = pp.user_id
                                  LEFT JOIN tunais
                                      ON s.id = tunais.user_id
                  ) datas`

		querySummaryExec, err := helpers.PrepareQuery(querySummary, templateEndQuery)

		// fmt.Println(querySummaryExec)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		queries = append(queries, querySummaryExec)

		// dataSummary, err := helpers.ExecuteQuery(querySummaryExec)
		// if err != nil {
		// 	fmt.Println("Error executing query 3:", err)
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal execute query",
		// 		Success: false,
		// 	})
		// }

		var wg sync.WaitGroup
		queryNames := []string{"data", "table_item", "summary"}
		mu := sync.Mutex{} // Prevent race conditions
		errCh := make(chan error, len(queries))

		for i, query := range queries {
			wg.Add(1)
			go func(label, sql string) {
				defer wg.Done()
				// var data []map[string]interface{}

				data, err := helpers.NewExecuteQuery(sql)
				if err != nil {
					errCh <- fmt.Errorf("%s query failed: %v", label, err)
					return
				}
				mu.Lock()
				resultsForReturn[label] = data
				mu.Unlock()
			}(queryNames[i], query) // Pass correct query name
		}
		wg.Add(1)
		go func(label, sql string) {
			defer wg.Done()

			restructuredArray := make(map[string]map[string]interface{})

			if implodedProdukHeader != "" {

				dataStokProduk, err := helpers.ExecuteQuery(queryStokProdukExec)
				if err != nil {
					errCh <- fmt.Errorf("%s query failed: %v", label, err)
					return
				}

				// Process each item
				for _, item := range dataStokProduk {
					tag := cases.Title(language.English, cases.NoLower).String(item["tag"].(string))
					code := item["code"].(string)
					packs := item["packs"]

					// Initialize map if the key does not exist
					if _, exists := restructuredArray[tag]; !exists {
						restructuredArray[tag] = make(map[string]interface{})
					}

					// Assign the packs value
					restructuredArray[tag][code] = packs
				}
			}

			mu.Lock()
			resultsForReturn["table_produk"] = restructuredArray
			mu.Unlock()
		}("table_produk", queryStokProdukExec)

		// Wait for all goroutines to finish
		wg.Wait()
		close(errCh)

		// Handle errors if any
		if len(errCh) > 0 {
			for err := range errCh {
				fmt.Println("Error:", err)
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Error fetching data"})
		}

		resultsForReturn["message"] = "Success Fetching Data"
		resultsForReturn["success"] = true
		// resultsForReturn["table_produk"] = restructuredArray
		resultsForReturn["table_produk_header"] = testHeader
		resultsForReturn["table_item_header"] = testHeaderItem
		resultsForReturn["user"] = getRequester[0]["requester"]
		resultsForReturn["tanggal"] = getRequester[0]["tanggal"]
		resultsForReturn["subject"] = getSubject[0]["requester"]

		return c.Status(fiber.StatusOK).JSON(resultsForReturn)

		// return c.Status(fiber.StatusOK).JSON(fiber.Map{
		// 	"message": "success",
		// 	"success": true,
		// 	"data":    data,
		// 	"tanggal": getRequester[0]["tanggal"],
		// 	"table_produk": restructuredArray,
		// 	"table_produk_header": testHeader,
		// 	"table_item": dataStokItem,
		// 	"table_item_header": testHeaderItem,
		// 	"summary" : dataSummary,
		// })
	} else {

		endQuery := `WITH visits as (
						SELECT id, tanggal_kunjungan, customer_id, user_id, subject_type_id, image_kunjungan
						FROM kunjungan 
						WHERE user_id = {{.QUserID}} 
							AND DATE(tanggal_kunjungan) = {{.QDate}}
						ORDER BY tanggal_kunjungan
					), penjualan_cus as (
						SELECT p.no_nota,
										p.id,
										CASE WHEN COALESCE(p.image_nota_print, p.image_nota) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/nota/'||COALESCE(p.image_nota_print, p.image_nota) ELSE NULL END as image_nota,
										CASE WHEN p.image_bukti_serah IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/nota/'||p.image_bukti_serah ELSE NULL END as image_bukti_serah,
										p.user_id,
										p.customer_id,
										DATE(p.tanggal_penjualan),
										COALESCE(SUM(CASE WHEN p.is_kredit = 0 THEN ((pd.harga - pd.diskon) * pd.jumlah) ELSE 0 END),0) as tunai,
										COALESCE(SUM(CASE WHEN p.is_kredit = 1 THEN ((pd.harga - pd.diskon) * pd.jumlah) ELSE 0 END),0) as kredit,
										pd.produk_id as produk_id,
										COALESCE(SUM(pd.jumlah),0) as jumlah
						FROM penjualan p
						JOIN penjualan_detail pd
							ON p.id = pd.penjualan_id
						WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = {{.QDate}}
						GROUP BY p.id, pd.produk_id
					), pembayaran_cus as (
						SELECT p.customer_id,
										p.user_id,
										COALESCE(SUM(pd.nominal),0) as pembayaran,
										COALESCE((CASE WHEN p.tipe_pelunasan = 0 AND p.pengembalian_id IS NOT NULL then p.total_pembayaran ELSE 0 END),0) AS adjustment,
										COALESCE(SUM((CASE WHEN tipe_pelunasan = 1 then pd.nominal ELSE 0 END)),0) as dp
						FROM pembayaran_piutang p
						JOIN pembayaran_piutang_detail pd
							ON p.id = pd.pembayaran_piutang_id
						LEFT JOIN payment py ON py.id =  p.payment_id AND py.tipe = 'BILYET GIRO'
						WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pembayaran) = {{.QDate}}
						GROUP BY p.id
					), pengembalian_cus as (
						SELECT p.customer_id,
										p.user_id,
										COALESCE(SUM(pd.jumlah * pd.harga),0) as pengembalian
						FROM pengembalian p
						JOIN pengembalian_detail pd
							ON p.id = pd.pengembalian_id
					WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pengembalian) = {{.QDate}}
						GROUP BY p.id
					), payments as (
						SELECT UPPER((CASE WHEN tipe = 'BILYET GIRO' THEN 'OPEN ' ELSE '' END)||tipe) AS payment_tipe, 
								COALESCE(nominal,0) AS payment_nominal,
										p.user_id,
										p.customer_id,
										CASE WHEN p.bukti_bayar = '' THEN NULL ELSE p.bukti_bayar END as payment_image,
																				p.penjualan_id as penjualan_id
								FROM payment p
								WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_transaksi) = {{.QDate}} 

						UNION

						SELECT 'CASH' as payment_tipe,
								p.total_penjualan as payment_nominal,
								p.user_id,
								p.customer_id,
								null as payment_image,
								p.id as penjualan_id
						FROM penjualan p
						LEFT JOIN payment py
								ON p.id = py.penjualan_id
								AND py.id IS NULL 
						WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_penjualan) = {{.QDate}} AND p.is_kredit = 0 
						GROUP BY p.id
					), kunjunganlogs as (
						SELECT DISTINCT ON (customer_id) * 
						FROM kunjungan_log WHERE user_id = {{.QUserID}} AND date(checkin_at) = {{.QDate}}
						ORDER BY customer_id, checkin_at
					), transactions as (
						SELECT t.user_id,
										t.subject_type_id,
										t.customer_id,
											JSONB_BUILD_OBJECT(
												'name', i.name,
												'quantity', COALESCE(SUM(td.qty),0)
										) as items,
										CASE WHEN ic.id = 191 AND t.transaction_type_id <> 1 
											THEN JSONB_BUILD_OBJECT(
												'photo_before', CASE WHEN MAX(t.photos_before) = '{}' OR MAX(t.photos_before) IS NULL 
																	THEN ARRAY(
																			SELECT CASE WHEN MAX(su.id) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/transaction/' ELSE 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo/' END|| image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT CASE WHEN MAX(su.id) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/transaction/' ELSE 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' END || image 
																			FROM unnest(MAX(t.photos_before)) image
																	) END,
												'photo_after', CASE WHEN MAX(t.photos_after) = '{}' OR MAX(t.photos_after) IS NULL 
																	THEN ARRAY(
																			SELECT CASE WHEN MAX(su.id) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/transaction/' ELSE 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo/' END || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT CASE WHEN MAX(su.id) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/transaction/' ELSE 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' END || image 
																			FROM unnest(MAX(t.photos_after)) image
																	) END
											)
											ELSE NULL
										END as sampling,
										CASE WHEN ic.id = 171 AND t.transaction_type_id <> 1 
											THEN JSONB_BUILD_OBJECT(
												'photo_before', CASE WHEN MAX(t.photos_before) = '{}' OR MAX(t.photos_before) IS NULL 
																	THEN ARRAY(
																			SELECT CASE WHEN MAX(su.id) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/transaction/' ELSE 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo/' END || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT CASE WHEN MAX(su.id) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/transaction/' ELSE 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' END || image 
																			FROM unnest(MAX(t.photos_before)) image
																		) END,
												'photo_after', CASE WHEN MAX(t.photos_after) = '{}' OR MAX(t.photos_after) IS NULL 
																	THEN ARRAY(
																			SELECT CASE WHEN MAX(su.id) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/transaction/' ELSE 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo/' END || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT CASE WHEN MAX(su.id) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/transaction/' ELSE 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' END || image 
																			FROM unnest(MAX(t.photos_after)) image
																		) END
											)
											ELSE NULL
										END as posm,
										CASE WHEN ic.id = 161 
											THEN JSONB_BUILD_OBJECT(
												'photo_before', CASE WHEN MAX(t.photos_before) = '{}' OR MAX(t.photos_before) IS NULL 
																	THEN ARRAY(
																			SELECT CASE WHEN MAX(su.id) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/transaction/' ELSE 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo/' END || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT CASE WHEN MAX(su.id) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/transaction/' ELSE 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_before/' END || image 
																			FROM unnest(MAX(t.photos_before)) image
																		) END,
												'photo_after', CASE WHEN MAX(t.photos_after) = '{}' OR MAX(t.photos_after) IS NULL 
																	THEN ARRAY(
																			SELECT CASE WHEN MAX(su.id) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/transaction/' ELSE 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo/' END || image 
																			FROM unnest(MAX(t.photos)) image
																		) 
																	ELSE ARRAY(
																			SELECT CASE WHEN MAX(su.id) IS NOT NULL THEN 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/transaction/' ELSE 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/photo_after/' END || image 
																			FROM unnest(MAX(t.photos_after)) image
																		) END
											)
											ELSE NULL
										END as merchandise
										
						FROM md.transaction t
						JOIN md.transaction_detail td
							ON t.id = td.transaction_id
						JOIN md.item i
							ON td.item_id = i.id
						JOIN md.item_category ic
							ON i.category_id = ic.id
						LEFT JOIN stok_user su
							ON su.user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}}
						WHERE t.user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}}
						GROUP BY t.user_id, t.customer_id, ic.id, t.transaction_type_id, t.subject_type_id, i.id
					)

					SELECT st.name as subject_type, 
							sq.customer, 
							MAX(invoice) as invoice,  
							{{.QSelectProd}}
							JSONB_AGG(sq.items) FILTER (WHERE sq.items IS NOT NULL) as items,
							COALESCE(SUM(sq.tunai),0) as tunai,
							COALESCE(SUM(sq.kredit),0) as kredit,
							COALESCE(SUM(sq.pengembalian), 0) as retur,
							COALESCE(SUM(sq.pembayaran),0) as pembayaran,
							COALESCE(SUM(sq.adjustment),0) as adjustment,
							COALESCE(SUM(sq.tunai),0) + 
							COALESCE(SUM(sq.pembayaran),0) + 
							COALESCE(SUM(sq.dp),0) - 
							COALESCE(SUM(sq.pengembalian),0) + 
							COALESCE(SUM(sq.adjustment),0) as setoran,
							JSONB_AGG(sq.payment_information) FILTER (WHERE sq.payment_information IS NOT NULL)->0 as payment_information,
							JSONB_AGG(sq.time_info) FILTER (WHERE sq.time_info IS NOT NULL)->0 as time_info,
							JSONB_AGG(sq.nota) FILTER (WHERE sq.nota IS NOT NULL)->0 as nota,
							JSONB_AGG(sq.penyerahan_produk) FILTER (WHERE sq.penyerahan_produk IS NOT NULL)->0 as penyerahan_produk,
							JSONB_AGG(sq.posm) FILTER (WHERE sq.posm IS NOT NULL)->0 as posm,
							JSONB_AGG(sq.merchandise) FILTER (WHERE sq.merchandise IS NOT NULL)->0 as merchandise,
							JSONB_AGG(sq.sampling) FILTER (WHERE sq.sampling IS NOT NULL)->0 as sampling
					FROM (
						SELECT v.subject_type_id,
								COALESCE(c.id,v.customer_id, pj.customer_id, kl.customer_id) as customer_id,
								JSONB_BUILD_OBJECT(
										'name', COALESCE(c.name, 'End User'),
										'contact', COALESCE(c.outlet_name, '-'),
										'location', COALESCE(c.kelurahan, '-'),
										'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
								) as customer,
								pj.no_nota as invoice,
								{{.QProd}}
								null::jsonb as items,
								COALESCE(SUM(pj.tunai),0) as tunai,
								COALESCE(SUM(pj.kredit),0) as kredit,
								0 as pembayaran,
								0 as pengembalian,
								0 as adjustment,
								0 as dp,
								NULL as payment_information,
								MIN(kl.checkin_at) as param,
								JSONB_BUILD_OBJECT(
									'checkin', to_char(COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan)) , 'HH24:MI'),
										'checkout', to_char(COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)), 'HH24:MI'),
										'durasi', COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)) - COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan))
								) as time_info,
								pj.image_nota as nota,
								pj.image_bukti_serah as penyerahan_produk,
								null::jsonb as posm,
								null::jsonb as merchandise,
								null::jsonb as sampling
							FROM visits v
							FULL JOIN penjualan_cus pj
								ON v.user_id = pj.user_id
								AND v.customer_id = pj.customer_id
							LEFT JOIN kunjunganlogs kl
								ON v.user_id = kl.user_id
								AND v.customer_id = kl.customer_id
							LEFT JOIN customer c
								ON COALESCE(v.customer_id, pj.customer_id, kl.customer_id) = c.id
							GROUP BY v.subject_type_id, v.image_kunjungan, c.id, COALESCE(c.id,v.customer_id, pj.customer_id, kl.customer_id), pj.no_nota, pj.image_nota, pj.image_bukti_serah
							
							UNION

							SELECT v.subject_type_id,
									COALESCE(c.id,v.customer_id, pem.customer_id, kl.customer_id) as customer_id,
									JSONB_BUILD_OBJECT(
													'name', COALESCE(c.name, 'End User'),
													'contact', COALESCE(c.outlet_name, '-'),
													'location', COALESCE(c.kelurahan, '-'),
													'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
									) as customer,
									NULL as invoice,
									{{.QProdItem}}
									null::jsonb as items,
									0 as tunai,
									0 as kredit,
									COALESCE(SUM(pem.pembayaran),0) as pembayaran,
									0 as pengembalian,
									COALESCE(SUM(pem.adjustment),0) as adjustment,
									COALESCE(SUM(pem.dp),0) as dp,
									NULL::jsonb as payment_information,
									MIN(kl.checkin_at) as param,
									JSONB_BUILD_OBJECT(
											'checkin', to_char(COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan)) , 'HH24:MI'),
											'checkout', to_char(COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)), 'HH24:MI'),
											'durasi', COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)) - COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan))
									) as time_info,
									null as nota,
									null as penyerahan_produk,
									null::jsonb as posm,
									null::jsonb as merchandise,
									null::jsonb as sampling
								FROM visits v
								FULL JOIN pembayaran_cus pem
										ON v.user_id = pem.user_id
										AND v.customer_id = pem.customer_id
								--FULL JOIN pengembalian_cus pg
								--        ON v.user_id = pg.user_id
								--        AND v.customer_id = pg.customer_id
								--FULL JOIN payments py
								--        ON v.user_id = py.user_id
								--        AND v.customer_id = py.customer_id
								LEFT JOIN kunjunganlogs kl
										ON v.user_id = kl.user_id
										AND v.customer_id = kl.customer_id
								LEFT JOIN customer c
										ON COALESCE(v.customer_id, pem.customer_id, kl.customer_id) = c.id
								GROUP BY v.subject_type_id, v.image_kunjungan, c.id, COALESCE(c.id,v.customer_id, pem.customer_id, kl.customer_id)
																					
							UNION
																					
								SELECT v.subject_type_id,
										COALESCE(c.id,v.customer_id, pg.customer_id, py.customer_id, kl.customer_id) as customer_id,
										JSONB_BUILD_OBJECT(
														'name', COALESCE(c.name, 'End User'),
														'contact', COALESCE(c.outlet_name, '-'),
														'location', COALESCE(c.kelurahan, '-'),
														'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
										) as customer,
										NULL as invoice,
										{{.QProdItem}}
										null::jsonb as items,
										0 as tunai,
										0 as kredit,
										0 as pembayaran,
										COALESCE(SUM(pg.pengembalian),0) as pengembalian,
										0 as adjustment,
										0 as dp,
										JSONB_BUILD_OBJECT(
										'cash', JSONB_BUILD_OBJECT(
														'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'CASH'), 0) 
																																			- COALESCE(SUM(pg.pengembalian),0),
														'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'CASH' AND py.payment_image <> '')
												),
										'transfer', JSONB_BUILD_OBJECT(
														'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'TRANSFER'), 0),
														'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'TRANSFER' AND py.payment_image <> '')
												),
										'bilyet_giro_cair', JSONB_BUILD_OBJECT(
														'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'BILYET GIRO'), 0),
														'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'BILYET GIRO' AND py.payment_image <> '')
												),
										'bilyet_giro_open', JSONB_BUILD_OBJECT(
														'value', COALESCE(SUM(py.payment_nominal) FILTER (WHERE py.payment_tipe = 'OPEN BILYET GIRO'), 0),
														'attachments', JSONB_AGG(py.payment_image) FILTER (WHERE py.payment_tipe = 'OPEN BILYET GIRO' AND py.payment_image <> '')
												)
										) as payment_information,
										MIN(kl.checkin_at) as param,
										JSONB_BUILD_OBJECT(
												'checkin', to_char(COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan)) , 'HH24:MI'),
												'checkout', to_char(COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)), 'HH24:MI'),
												'durasi', COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)) - COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan))
										) as time_info,
										null as nota,
										null as penyerahan_produk,
										null::jsonb as posm,
										null::jsonb as merchandise,
										null::jsonb as sampling
								FROM visits v
								-- FULL JOIN pembayaran_cus pem
								--      ON v.user_id = pem.user_id
								--      AND v.customer_id = pem.customer_id
								FULL JOIN pengembalian_cus pg
										ON v.user_id = pg.user_id
										AND v.customer_id = pg.customer_id
								FULL JOIN payments py
										ON v.user_id = py.user_id
										AND v.customer_id = py.customer_id
								LEFT JOIN kunjunganlogs kl
										ON v.user_id = kl.user_id
										AND v.customer_id = kl.customer_id
								LEFT JOIN customer c
										ON COALESCE(v.customer_id, pg.customer_id, py.customer_id, kl.customer_id) = c.id
								GROUP BY v.subject_type_id, v.image_kunjungan, c.id, COALESCE(c.id,v.customer_id, pg.customer_id, py.customer_id, kl.customer_id)
							
							UNION
							
								SELECT v.subject_type_id,
									COALESCE(c.id, v.customer_id, t.customer_id, kl.customer_id) as customer_id,
									JSONB_BUILD_OBJECT(
											'name', COALESCE(c.name, 'End User'),
											'contact', COALESCE(c.outlet_name, '-'),
											'location', COALESCE(c.kelurahan, '-'),
											'photo', 'https://assets-sales.s3.ap-southeast-3.amazonaws.com/kunjungan/'||v.image_kunjungan
									) as customer,
									NULL as invoice,
									{{.QProdItem}}
									t.items as items,
									0 as tunai,
									0 as kredit,
									0 as pembayaran,
									0 as pengembalian,
									0 as adjustment,
									0 as dp,
									NULL as payment_information,
									MIN(kl.checkin_at) as param,
									JSONB_BUILD_OBJECT(
										'checkin', to_char(COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan)) , 'HH24:MI'),
										'checkout', to_char(COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)), 'HH24:MI'),
										'durasi', COALESCE(MAX(kl.checkout_at), MIN(v.tanggal_kunjungan)) - COALESCE(MIN(kl.checkin_at), MIN(v.tanggal_kunjungan))
									) as time_info,
									null as nota,
									null as penyerahan_produk,
									t.posm as posm,
									t.merchandise as merchandise,
									t.sampling as sampling
							FROM visits v
							FULL JOIN transactions t
								ON v.user_id = t.user_id
								AND v.customer_id = t.customer_id
							LEFT JOIN kunjunganlogs kl
								ON v.user_id = kl.user_id
								AND v.customer_id = kl.customer_id
							LEFT JOIN customer c
								ON COALESCE(v.customer_id, t.customer_id, kl.customer_id) = c.id
							GROUP BY v.subject_type_id, v.image_kunjungan, c.id, COALESCE(c.id, v.customer_id, t.customer_id, kl.customer_id), t.items, t.posm, t.merchandise, t.sampling
					) sq
					JOIN md.subject_type st
						ON COALESCE(sq.subject_type_id, CASE WHEN sq.customer_id < 0 THEN 3 ELSE 1 END) = st.id
					GROUP BY sq.param, st.id, sq.customer
					ORDER BY sq.param`

		finalQuery, err := helpers.PrepareQuery(endQuery, templateEndQuery)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		// data, err := helpers.NewExecuteQuery(finalQuery)
		// if err != nil {
		// 	fmt.Println("Error executing query 1:", err)
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal execute query",
		// 		Success: false,
		// 	})
		// }

		// fmt.Println(finalQuery)

		queries = append(queries, finalQuery)

		queryStokProduk := `SELECT 
								JSON_BUILD_OBJECT('produk_id',p.id, 'kode', p.code, 'nama', p.name, 'foto', p.foto) AS produk,
								SUM(COALESCE(ss.stok_awal,0)-COALESCE(ssr_order.jumlah,0)) AS stok_awal,
								SUM(COALESCE(ssr_order.jumlah,0)) AS order,
								SUM(COALESCE(pj.jumlah,0)) AS penjualan,
								SUM(COALESCE(prg.jumlah,0)) AS program,
								SUM(COALESCE(pg.jumlah,0)) AS retur_customer,
								SUM(COALESCE(ssr_retur.jumlah,0)) AS retur_gudang,
								SUM(COALESCE(ss.stok_akhir,0)) AS stok_akhir
							FROM produk p  
							JOIN produk_branch pb
							ON pb.produk_id = p.id
							LEFT JOIN 
								stok_salesman ss 
								ON ss.produk_id = p.id AND ss.user_id = {{.QUserID}} AND DATE(ss.tanggal_stok) = {{.QDate}}
							LEFT JOIN
							(
								SELECT produk_id, condition, pita, SUM(COALESCE(jumlah,0)) AS jumlah 
								FROM stok_salesman_riwayat 
								WHERE is_validate = 1 AND DATE(tanggal_riwayat) = {{.QDate}} 
								AND user_id = {{.QUserID}} AND aksi='ORDER'
								GROUP BY produk_id, condition, pita
							) ssr_order ON ssr_order.produk_id = ss.produk_id AND ssr_order.condition = ss.condition AND ssr_order.pita = ss.pita
							LEFT JOIN
							( 
								SELECT produk_id, condition, pita, SUM(COALESCE(jumlah,0)) AS jumlah 
								FROM stok_salesman_riwayat 
								WHERE is_validate = 1 AND DATE(tanggal_riwayat) = {{.QDate}} 
								AND user_id = {{.QUserID}} AND aksi='RETUR'
								GROUP BY produk_id, condition, pita
							) ssr_retur ON ssr_retur.produk_id = p.id AND ssr_retur.condition = ss.condition AND ssr_retur.pita = ss.pita
							LEFT JOIN 
							(
								SELECT pd.produk_id, pd.condition, pd.pita, SUM(COALESCE(pd.jumlah,0)) AS jumlah
								FROM penjualan p 
								JOIN penjualan_detail pd 
								ON pd.penjualan_id = p.id
								WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_penjualan) = {{.QDate}} AND pd.harga >0
								GROUP BY pd.produk_id, pd.condition, pd.pita
							) pj ON pj.produk_id = p.id AND pj.condition = ss.condition AND pj.pita = ss.pita
							LEFT JOIN
							(
								SELECT pd.produk_id, pd.condition, pd.pita, SUM(COALESCE(pd.jumlah,0)) AS jumlah
								FROM penjualan p 
								JOIN penjualan_detail pd 
								ON pd.penjualan_id = p.id
								WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_penjualan) = {{.QDate}} AND pd.harga =0
								GROUP BY pd.produk_id, pd.condition, pd.pita
							) prg ON prg.produk_id = p.id AND prg.condition = ss.condition AND prg.pita = ss.pita
							LEFT JOIN
							(
								SELECT pd.produk_id, pd.condition, pd.pita, SUM(COALESCE(pd.jumlah,0)) AS jumlah
								FROM pengembalian p 
								JOIN pengembalian_detail pd 
								ON pd.pengembalian_id = p.id
								WHERE p.user_id = {{.QUserID}} AND DATE(tanggal_pengembalian) = {{.QDate}}
								GROUP BY pd.produk_id, pd.condition, pd.pita
							) pg ON pg.produk_id = p.id AND pg.condition = ss.condition AND pg.pita = ss.pita
							WHERE pb.branch_id IN ( {{.QBranchID}} )
							GROUP BY p.id
							ORDER BY p.order`

		queryStokProdukExec, err := helpers.PrepareQuery(queryStokProduk, templateEndQuery)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		queries = append(queries, queryStokProdukExec)

		// fmt.Println(queryStokProdukExec)

		// dataStokProduk, err := helpers.NewExecuteQuery(queryStokProdukExec)
		// if err != nil {
		// 	fmt.Println("Error executing query 2:", err)
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal execute query",
		// 		Success: false,
		// 	})
		// }

		// testHeader := []string{}
		// testHeader = append(testHeader, "Stok Awal")
		// testHeader = append(testHeader, "Order")
		// testHeader = append(testHeader, "Penjualan")
		// testHeader = append(testHeader, "Program")
		// testHeader = append(testHeader, "Retur Customer")
		// testHeader = append(testHeader, "Retur Gudang")
		// testHeader = append(testHeader, "Stok Akhir")
		testHeader := []string{"Stok Awal", "Order", "Penjualan", "Program", "Retur Customer", "Retur Gudang", "Stok Akhir"}

		queryStokItem := `WITH ordersreturs as (
							SELECT user_id, item_id, aksi, SUM(jumlah) as jumlah
							FROM md.stok_merchandiser_riwayat
							WHERE user_id = {{.QUserID}} AND DATE(tanggal_riwayat) = {{.QDate}}
							GROUP BY user_id, item_id, aksi
						), stoks as (
							SELECT user_id, item_id, stok_awal, stok_akhir
							FROM md.stok_merchandiser
							WHERE user_id = {{.QUserID}} AND DATE(tanggal_stok) = {{.QDate}}
						), transactions as (
							SELECT t.user_id, td.item_id, SUM(qty) as jumlah
							FROM md.transaction t
							JOIN md.transaction_detail td
								ON t.id = td.transaction_id
							WHERE t.user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}}
							GROUP BY t.user_id, td.item_id
						)

						SELECT CASE WHEN pb.id IS NOT NULL THEN CONCAT(i.name, ' - ', pb.name) ELSE i.name END as item_name,
								JSONB_BUILD_OBJECT(
									'Stok Awal', MAX(s.stok_awal) - SUM(CASE WHEN o.aksi = 'ORDER' THEN o.jumlah ELSE 0 END),
									'Order', SUM(CASE WHEN o.aksi = 'ORDER' THEN o.jumlah ELSE 0 END),
									'Transaksi', MAX(COALESCE(t.jumlah,0)),
									'Retur Gudang', SUM(CASE WHEN o.aksi = 'RETUR' THEN o.jumlah ELSE 0 END),
									'Stok Akhir', MAX(s.stok_akhir)
								) as datas
						FROM stoks s
						FULL JOIN ordersreturs o
							ON s.item_id = o.item_id
						FULL JOIN transactions t
							ON s.item_id = t.item_id
						FULL JOIN md.merchandiser m
							ON m.user_id = COALESCE(s.user_id, o.user_id, t.user_id)
						LEFT JOIN md.item i
							ON s.item_id = i.id
						LEFT JOIN produk_brand pb
							ON i.brand_id = pb.id
						WHERE m.user_id = {{.QUserID}}
						GROUP BY i.id, pb.id
						ORDER BY i.name`

		queryStokItemExec, err := helpers.PrepareQuery(queryStokItem, templateEndQuery)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		// dataStokItem, err := helpers.ExecuteQuery(queryStokItemExec)
		// if err != nil {
		// 	fmt.Println("Error executing query 3:", err)
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal execute query",
		// 		Success: false,
		// 	})
		// }

		queries = append(queries, queryStokItemExec)

		testHeaderItem := []string{"Stok Awal", "Order", "Transaksi", "Retur Gudang", "Stok Akhir"}

		queryVisit := `WITH data_aktifitas as (
							SELECT ssq.customer_id, ssq.date, COALESCE(SUM(DISTINCT ssq.count),0) as count, COALESCE(SUM(DISTINCT ssq.is_sampling),0) as is_sampling, COALESCE(SUM(ssq.is_posm),0) as is_posm, COALESCE(SUM(ssq.is_merchandise),0) as is_merchandise
							FROM (
								SELECT sq.customer_id, sq.date, COUNT(sq.id), 0 as is_sampling, 0 as is_posm, 0 as is_merchandise
								FROM (
									SELECT id, customer_id, DATE(tanggal_penjualan) FROM penjualan WHERE user_id = {{.QUserID}} AND DATE(tanggal_penjualan) = {{.QDate}}
									UNION
									SELECT id, customer_id, DATE(tanggal_pengembalian) FROM pengembalian WHERE user_id = {{.QUserID}} AND DATE(tanggal_pengembalian) = {{.QDate}}
									UNION
									SELECT id, customer_id, DATE(tanggal_pembayaran) FROM pembayaran_piutang WHERE user_id = {{.QUserID}} AND DATE(tanggal_pembayaran) = {{.QDate}}
									UNION
									SELECT id, customer_id, DATE(datetime) FROM qr_code_history WHERE user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}}
									UNION
									SELECT sq.* FROM (
										SELECT DISTINCT ON (customer_id) id, customer_id, DATE(checkin_at)FROM kunjungan_log WHERE user_id = {{.QUserID}} AND DATE(checkin_at) = {{.QDate}} ORDER BY customer_id, checkin_at DESC
									) sq
									UNION
									SELECT sq.* FROM (
										SELECT DISTINCT ON (customer_id) id, customer_id, DATE(checkout_at) FROM kunjungan_log WHERE user_id = {{.QUserID}} AND DATE(checkout_at) = {{.QDate}} ORDER BY customer_id, checkout_at DESC
									) sq
								) sq
								GROUP BY sq.customer_id, sq.date
							
							UNION ALL

								SELECT sq.customer_id, sq.date, 0, COALESCE(SUM(DISTINCT sq.is_sampling),0) as is_sampling, COALESCE(SUM(DISTINCT sq.is_posm),0) as is_posm, COALESCE(SUM(DISTINCT sq.is_merchandise),0) as is_merchandise
								FROM (
								
										SELECT t.id, customer_id, DATE(datetime), 1 as is_sampling, 0 as is_posm, 0 as is_merchandise FROM md.transaction t
										JOIN md.transaction_detail td
											ON t.id = td.transaction_id
										JOIN md.item i
											ON td.item_id = i.id
										JOIN md.item_category ic
											ON i.category_id = ic.id
										WHERE user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}} AND ic.id = 191

										UNION ALL

										SELECT t.id, customer_id, DATE(datetime), 0 as is_sampling, 1 as is_posm, 0 as is_merchandise FROM md.transaction t
										JOIN md.transaction_detail td
											ON t.id = td.transaction_id
										JOIN md.item i
											ON td.item_id = i.id
										JOIN md.item_category ic
											ON i.category_id = ic.id
										WHERE user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}} AND ic.id = 171

										UNION ALL

										SELECT t.id, customer_id, DATE(datetime), 0 as is_sampling, 0 as is_posm, 1 as is_merchandise FROM md.transaction t
										JOIN md.transaction_detail td
											ON t.id = td.transaction_id
										JOIN md.item i
											ON td.item_id = i.id
										JOIN md.item_category ic
											ON i.category_id = ic.id
										WHERE user_id = {{.QUserID}} AND DATE(datetime) = {{.QDate}} AND ic.id = 161
															
									) sq
									GROUP BY sq.customer_id, sq.date
								) ssq
								GROUP BY ssq.customer_id, ssq.date
						)


						SELECT DISTINCT ON (COALESCE(c.id, da.customer_id), DATE(tanggal_kunjungan)) 
							ROW_NUMBER() OVER(ORDER BY tanggal_kunjungan) as no,
							DATE(k.tanggal_kunjungan),
							u.full_name,
							CASE WHEN c.id IS NOT NULL THEN
							JSONB_BUILD_OBJECT(
								'id', c.id||'',
								'name', c.name,
								'outlet_name', c.outlet_name,
								'type', COALESCE(k.customer_tipe, ct.name)
							)
							ELSE 
							JSONB_BUILD_OBJECT(
								'id', da.customer_id,
								'name', 'SMOKER',
								'outlet_name', '-',
								'type', 'SMOKER'
							) END as customer,
							st.name as subject_type,
							c.alamat,
							k.status_toko,
							k.keterangan,
							JSONB_BUILD_OBJECT(
								'checkin', to_char(kl.checkin_at, 'HH24:MI'),
								'checkout', to_char(kl.checkout_at, 'HH24:MI'),
								'durasi', AGE(kl.checkout_at, kl.checkin_at),
								'jumlah_aktifitas', da.count
							) as informasi,
							k.latitude_longitude as lokasi,
							k.image_kunjungan,
							CASE WHEN p.id IS NOT NULL THEN 1 ELSE 0 END as is_ec,
							CASE WHEN DATE(c.dtm_crt) = DATE(k.tanggal_kunjungan) THEN 1 ELSE 0 END as new_register,
							da.is_sampling,
							da.is_posm,
							da.is_merchandise
					FROM kunjungan k
					LEFT JOIN public.user u
						ON k.user_id = u.id
					LEFT JOIN customer c
						ON k.customer_id = c.id
					LEFT JOIN customer_type ct
						ON c.tipe = ct.id
					LEFT JOIN (
						SELECT ssq.customer_id, MAX(checkin_at) as checkin_at, MAX(checkout_at) as checkout_at
						FROM (
							SELECT sq.* FROM (
								SELECT DISTINCT ON (customer_id) customer_id, checkin_at, NULL::timestamp as checkout_at 
									FROM kunjungan_log 
									WHERE user_id = {{.QUserID}} 
									AND DATE(checkin_at) = {{.QDate}} 
									ORDER BY customer_id, checkin_at DESC
							) sq
							UNION
							SELECT sq.* FROM (
								SELECT DISTINCT ON (customer_id) customer_id, NULL::timestamp, checkout_at 
									FROM kunjungan_log 
									WHERE user_id = {{.QUserID}} 
									AND DATE(checkout_at) = {{.QDate}} 
									ORDER BY customer_id, checkout_at DESC
							) sq
						) ssq
						GROUP BY ssq.customer_id
					) kl
						ON k.customer_id = kl.customer_id
					LEFT JOIN data_aktifitas da
						ON k.customer_id = da.customer_id
						AND DATE(k.tanggal_kunjungan) = da.date
					LEFT JOIN penjualan p
						ON k.customer_id = p.customer_id
						AND DATE(k.tanggal_kunjungan) = DATE(p.tanggal_penjualan)
					LEFT JOIN md.subject_type st
						ON CASE WHEN k.subject_type_id IS NULL AND k.customer_id < 0 THEN 3 
								WHEN k.subject_type_id IS NULL AND k.customer_id > 0 THEN 1
							ELSE k.subject_type_id END = st.id
					WHERE k.user_id IN ({{.QUserID}}) 
						AND DATE(k.tanggal_kunjungan) = {{.QDate}}
					ORDER BY COALESCE(c.id, da.customer_id), DATE(tanggal_kunjungan)`

		queryVisitExec, err := helpers.PrepareQuery(queryVisit, templateEndQuery)

		// fmt.Println(queryVisitExec)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		// dataVisit, err := helpers.ExecuteQuery(queryVisitExec)
		// if err != nil {
		// 	fmt.Println("Error executing query 4:", err)
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal execute query",
		// 		Success: false,
		// 	})
		// }

		queries = append(queries, queryVisitExec)

		querySummary := `WITH tunais as (
                          SELECT SUM(sum_cash) as total_tunai,
                          {{.QAsUserID}}
                          FROM (
                              SELECT COALESCE(SUM(CASE WHEN p.is_kredit = 0 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END),0) as sum_cash, 1 as param
                              FROM penjualan p
                              LEFT JOIN penjualan_detail pd
                                      ON p.id = pd.penjualan_id
                              LEFT JOIN payment py
                                      ON p.id = py.penjualan_id
                              WHERE py.id IS NULL 
                                      AND p.is_kredit = 0
                                      AND p.user_id = {{.QUserID}}
                                      AND DATE(p.tanggal_penjualan) = {{.QDate}}
                                      
                              UNION ALL

                              SELECT SUM(py.nominal) as sum_cash
                                -- COALESCE(SUM(CASE WHEN p.is_kredit = 0 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END),0) as sum_cash
                                , 2
                              FROM penjualan p
                              JOIN payment py
									ON p.id = py.penjualan_id
                              LEFT JOIN pembayaran_piutang pp 
									ON py.id = pp.payment_id
                              WHERE UPPER(py.tipe) = 'CASH'
                                      AND pp.id IS NULL
                                      AND p.user_id = {{.QUserID}}
                                      AND DATE(p.tanggal_penjualan) = {{.QDate}}
                                      
                              UNION ALL

                              SELECT COALESCE(pp.total_pembayaran,0) as sum_cash, 3
                              FROM pembayaran_piutang pp
                              LEFT JOIN payment py 
                                      ON pp.payment_id = py.id
                              WHERE py.id IS NULL
                                      AND pp.user_id = {{.QUserID}}
                                      AND DATE(pp.tanggal_pembayaran) = {{.QDate}}
                                      
                              UNION ALL

                              SELECT COALESCE(pp.total_pembayaran,0) as sum_cash, 4
                              FROM pembayaran_piutang pp
                              JOIN payment py 
                                      ON pp.payment_id = py.id
                              WHERE UPPER(py.tipe) = 'CASH'
                                      AND pp.user_id = {{.QUserID}}
                                      AND DATE(pp.tanggal_pembayaran) = {{.QDate}}
                      ) sq
                  ), total_call as (
                      SELECT COUNT(kunjungan.id) AS total_call, {{.QAsUserID}}
                          FROM (
                                          SELECT            
                                          DISTINCT ON(k.customer_id) k.id, k.salesman_id
                                          FROM kunjungan k 
                                          -- LEFT JOIN penjualan p 
                                          -- ON p.customer_id = k.customer_id AND DATE(p.tanggal_penjualan) =  DATE({{.QDate}}) 
                                          -- AND (p.salesman_id = $salesmanId OR p.merchandiser_id = $merchandiserId)
                                          WHERE k.user_id = {{.QUserID}}
                                           AND DATE(tanggal_kunjungan) = DATE({{.QDate}}) AND UPPER(k.status_toko) = 'BUKA' AND k.customer_id > 0
                                          ORDER BY k.customer_id, tanggal_kunjungan ASC
                          ) kunjungan
                  ), sales as (
                      SELECT
                              SUM(CASE WHEN pd.harga > 0 THEN pd.jumlah ELSE 0 END) AS total_sales_pack,
                              SUM(pd.jumlah* (pd.harga-pd.diskon)) AS total_sales, 
                              COUNT(DISTINCT p.customer_id) FILTER (WHERE p.customer_id > 0 AND pd.jumlah >= 2 AND pd.harga > 0) AS total_effective_call,
                              SUM(CASE WHEN p.is_kredit = 0 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END) AS total_cash,
                              SUM(CASE WHEN p.is_kredit = 1 THEN pd.jumlah* (pd.harga-pd.diskon) ELSE 0 END) AS total_credit,
                              SUM(CASE WHEN pd.harga > 0 AND p.is_kredit = 0 THEN pd.jumlah ELSE 0 END) AS total_cash_pack,
                              SUM(CASE WHEN pd.harga > 0 AND p.is_kredit = 1 THEN pd.jumlah ELSE 0 END) AS total_credit_pack,
                              SUM(CASE WHEN pd.harga = 0 THEN pd.jumlah ELSE 0 END) AS total_program_pack,
                              {{.QAsUserID}}
                              FROM penjualan p
                              JOIN penjualan_detail pd
                              ON p.id = pd.penjualan_id
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_penjualan) = DATE({{.QDate}})
                  ), retur as (
                      SELECT 
                              COALESCE(SUM(pd.jumlah * pd.harga),0) AS total_return,
                              COALESCE(SUM(pd.jumlah),0) AS total_return_pack,
                              {{.QAsUserID}}
                      FROM 
                              pengembalian p
                              JOIN pengembalian_detail pd
                              ON p.id = pd.pengembalian_id
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pengembalian) = DATE({{.QDate}})
                  ),customer_register as (
                      SELECT 
                              COUNT(p.id) as customer_register,
                              {{.QAsUserID}}
                      FROM 
                      customer p
                      WHERE p.user_id_holder = {{.QUserID}} AND DATE(p.dtm_crt) = DATE({{.QDate}})
                  ), timecall as (
                      SELECT 
                              SUM(AGE(p.checkout_at,p.checkin_at))/CASE WHEN count(distinct customer_id) = 0 THEN 1 ELSE count(distinct customer_id) END as time_call,
                              {{.QAsUserID}}
                      FROM 
                              kunjungan_log p
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.checkin_at) = DATE({{.QDate}})
                  ), payments as (
                      SELECT
                          SUM(CASE WHEN UPPER(p.tipe) = 'CASH' THEN nominal ELSE 0 END) AS total_payment_cash,
                          SUM(CASE WHEN UPPER(p.tipe) = 'TRANSFER' THEN nominal ELSE 0 END) AS total_payment_transfer,
                          SUM(CASE WHEN UPPER(p.tipe) = 'CEK' THEN nominal ELSE 0 END) AS total_payment_check,
                          SUM(CASE WHEN UPPER(p.tipe) = 'BILYET GIRO' AND (DATE(tanggal_cair) != DATE({{.QDate}}) OR tanggal_cair IS NULL) THEN nominal ELSE 0 END) AS total_payment_open_bilyet_giro,
                          SUM(CASE WHEN UPPER(p.tipe) = 'BILYET GIRO' AND is_cair = 1 AND DATE(tanggal_cair) = DATE({{.QDate}})  THEN nominal ELSE 0 END) AS total_payment_bilyet_giro,
                          {{.QAsUserID}}
                      FROM
                          payment p
                      JOIN penjualan pj
                      ON pj.id = p.penjualan_id 
                      WHERE p.user_id = {{.QUserID}} AND p.is_verif>-1 --AND pj.is_kredit = 0 
                      AND (
                          DATE(p.tanggal_transaksi) = DATE({{.QDate}}) 
                          OR (DATE(p.tanggal_cair) = DATE({{.QDate}}) AND p.is_cair = 1)
                      )
                  ) , 
                  pembayaran_piutang as (
                              SELECT COALESCE(SUM(p.total_pembayaran),0) as total_payment_and_dp,
                              {{.QAsUserID}}
                              
                              FROM
                              pembayaran_piutang p
                              JOIN pembayaran_piutang_detail pd 
                              ON p.id = pd.pembayaran_piutang_id
                              LEFT JOIN piutang pi
                              ON pi.id = pd.piutang_id
                              WHERE p.user_id = {{.QUserID}} AND DATE(p.tanggal_pembayaran) = DATE({{.QDate}})
                  )

                  SELECT 
                              SUM(datas.total_call) as total_call_buka,
                              SUM(datas.total_effective_call) AS total_effective_call_2pack,
                              SUM(datas.customer_register) as total_register_customer,
                              MAX(datas.time_call) as average_call,
                              SUM(datas.total_sales_pack) AS omzet,
                              SUM(datas.total_cash_pack) AS omzet_tunai,
                              SUM(datas.total_credit_pack) AS omzet_kredit,
                              SUM(datas.total_return_pack) AS retur,
                              --COALESCE(SUM(datas.total_payment_cash) - SUM(datas.total_return),0) AS pembayaran_tunai,
                              COALESCE(MAX(datas.total_tunai) - SUM(datas.total_return),0) AS pembayaran_tunai,
                              COALESCE(SUM(datas.total_payment_transfer),0) AS pembayaran_transfer,
                              COALESCE(SUM(datas.total_payment_check),0) AS cek,
                              COALESCE(SUM(datas.total_payment_open_bilyet_giro),0) AS bilyet_giro_baru,
                              COALESCE(SUM(datas.total_payment_bilyet_giro),0) AS bilyet_giro_cair,
                              SUM(datas.total_cash)+ SUM(datas.total_payment_and_dp)
                                -SUM(datas.total_return)-COALESCE(SUM(datas.total_payment_transfer),0)
                                -COALESCE(SUM(datas.total_payment_check),0)
                                -COALESCE(SUM(datas.total_payment_open_bilyet_giro),0)
                                -COALESCE(SUM(datas.total_payment_bilyet_giro),0) AS total_setoran,
                              MAX(datas.total_tunai)-SUM(datas.total_return) AS total_setoran_tunai
                  FROM ( SELECT tc.*, sales.*, retur.*, cr.*, timecall.*, payments.*, pp.*, tunais.*
                                  FROM public.user s
                                  LEFT JOIN total_call tc
                                      ON s.id = tc.user_id
                                  LEFT JOIN sales
                                      ON s.id = sales.user_id
                                  LEFT JOIN retur
                                      ON s.id = retur.user_id
                                  LEFT JOIN customer_register cr
                                      ON s.id = cr.user_id
                                  LEFT JOIN timecall
                                      ON s.id = timecall.user_id
                                  LEFT JOIN payments
                                      ON s.id = payments.user_id
                                  LEFT JOIN pembayaran_piutang pp
                                      ON s.id = pp.user_id
                                  LEFT JOIN tunais
                                      ON s.id = tunais.user_id
                  ) datas`

		querySummaryExec, err := helpers.PrepareQuery(querySummary, templateEndQuery)

		// fmt.Println(querySummaryExec)

		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Terjadi kesalahan ketika generate query",
				Success: false,
			})
		}

		// dataSummary, err := helpers.ExecuteQuery(querySummaryExec)
		// if err != nil {
		// 	fmt.Println("Error executing query 5:", err)
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal execute query",
		// 		Success: false,
		// 	})
		// }

		queries = append(queries, querySummaryExec)

		var wg sync.WaitGroup
		queryNames := []string{"data", "table_produk", "table_item", "visit", "summary"}
		mu := sync.Mutex{} // Prevent race conditions
		errCh := make(chan error, len(queries))

		for i, query := range queries {
			wg.Add(1)
			go func(label, sql string) {
				defer wg.Done()
				// var data []map[string]interface{}

				data, err := helpers.NewExecuteQuery(sql)
				if err != nil {
					errCh <- fmt.Errorf("%s query failed: %v", label, err)
					return
				}
				mu.Lock()
				resultsForReturn[label] = data
				mu.Unlock()
			}(queryNames[i], query) // Pass correct query name
		}

		// Wait for all goroutines to finish
		wg.Wait()
		close(errCh)

		// Handle errors if any
		if len(errCh) > 0 {
			for err := range errCh {
				fmt.Println("Error:", err)
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Error fetching data"})
		}

		resultsForReturn["message"] = "Success Fetching Data"
		resultsForReturn["success"] = true
		// resultsForReturn["table_produk"] = restructuredArray
		resultsForReturn["table_produk_header"] = testHeader
		resultsForReturn["table_item_header"] = testHeaderItem

		return c.Status(fiber.StatusOK).JSON(resultsForReturn)
	}
}
