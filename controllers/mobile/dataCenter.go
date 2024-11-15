package controllers

import (
	"fmt"
	"math"
	"net/http"
	db "pluto_remastered/config"
	"pluto_remastered/helpers"
	"pluto_remastered/models"
	"pluto_remastered/structs"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetAppVersioning(c *fiber.Ctx) error {
	datas, err := helpers.NewExecuteQuery(`SELECT app_name,
                                            JSONB_BUILD_OBJECT(
                                                'current_version', current_version,
                                                'minimum_version', minimum_version,
                                                'force_update', force_update,
                                                'changelog', changelog,
                                                'android_url', android_url,
                                                'ios_url', ios_url
                                            ) as version,
                                            JSONB_BUILD_OBJECT(
                                                'is_maintenance', CASE WHEN is_maintenance = 1 THEN TRUE ELSE FALSE END,
                                                'message', message
                                            ) as maintenance,
                                            JSONB_BUILD_OBJECT(
                                                'minimum_android_version', minimum_android_version,
                                                'minimum_ios_version', minimum_ios_version
                                            ) as os_compatibility,
                                            JSONB_BUILD_OBJECT(
                                                'api_base_url', api_base_url,
                                                'request_timeout', request_timeout,
                                                'max_upload_size_mb', max_upload_size_mb
                                            ) as configurations,
                                            JSONB_BUILD_OBJECT(
                                                'terms_and_conditions', terms_and_conditions,
                                                'privacy_policy', privacy_policy,
                                                'landing', landing
                                            ) as urls,
                                            JSONB_BUILD_OBJECT(
                                                'server_time', now(),
                                                'timezone', current_setting('TIMEZONE')
                                            ) as server_info
                                            FROM app_versioning
                                            ORDER BY id DESC`)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika eksekusi query",
			Success: false,
		})
	}

	if len(datas) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(helpers.ResponseWithoutData{
			Message: "Data tidak ditemukan",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data berhasil ditemukan",
		Success: true,
		Data:    datas,
	})
}

func GetDataRequests(c *fiber.Ctx) error {
	type TemplateInputUser struct {
		UserId *string `json:"userId"`
		Date   *string `json:"date"`
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

	whereDate := ""

	if inputUser.Date != nil {
		whereDate = " AND DATE('" + *inputUser.Date + "')"
	} else {
		whereDate = " CURRENT_DATE"
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

	templateQuery := `SELECT 'public.checkin_request' as ref_table, 
                        NULL as ref_table_child,
                        'Akses Checkin Ulang' as tag,
                        cr.id as ref_id,
                        cr.is_approve,
                        to_char(cr.datetime, 'YYYY-MM-DD HH24:MI:SS') as datetime,
                        cr.note,
                        ARRAY_AGG(k.customer_id) as customer_id,
                        JSONB_AGG(
                                JSONB_BUILD_OBJECT(
                                    'id', c.id||'',
                                    'name', c.name,
                                    'outlet_name', c.outlet_name,
                                    'tipe', ct.name
                                )
                            ) as customers
                    FROM checkin_request cr 
					LEFT JOIN hr.employee e
						ON cr.employee_id = e.id
                    LEFT JOIN kunjungan k
                        ON cr.kunjungan_id = k.id
                    LEFT JOIN customer c
                        ON k.customer_id = c.id
                    LEFT JOIN customer_type ct
                        ON c.tipe = ct.id
                    WHERE e.user_id = {{.QUserId}} AND DATE(cr.datetime) = {{.QDate}}
                    GROUP BY cr.id

                UNION ALL 

                SELECT 'public.customer_move_request', 
                        NULL as ref_table_child,
                        'Perpindahan Customer' as tag,  
                        cm.id as ref_id,
                        cm.is_approve,
                        to_char(cm.request_at, 'YYYY-MM-DD HH24:MI:SS'),
                        cm.note,
                        cm.customer_id,
                        JSONB_AGG(
                            JSONB_BUILD_OBJECT(
                                'id', c.id||'',
                                'name', c.name,
                                'outlet_name', c.outlet_name,
                                'tipe', ct.name
                            )
                        ) as customers
                    FROM customer_move_request cm
					LEFT JOIN hr.employee e
						ON cm.requested_id = e.id
                    JOIN customer c
                        ON c.id = ANY(cm.customer_id)
                    JOIN customer_type ct
                        ON c.tipe = ct.id
                    WHERE e.user_id = {{.QUserId}} AND DATE(cm.request_at) = {{.QDate}}
                    GROUP BY cm.id

                UNION ALL

                SELECT 'public.salesman_access',
                        'public.salesman_access_detail' as ref_table_child,
                        'Akses Retur' as tag,        
                        sa.id as ref_id,
                        sa.is_approve,
                        to_char(sa.request_at, 'YYYY-MM-DD HH24:MI:SS'),
                        sa.note,
                        null,
                        null
                    FROM salesman_access sa
					LEFT JOIN hr.employee e
						ON sa.requested_id = e.id
					WHERE e.user_id = {{.QUserId}} AND sa.access_type = 'retur' AND DATE(sa.request_at) = {{.QDate}}

                UNION ALL

                SELECT 'public.salesman_access',
                        'public.salesman_access_detail' as ref_table_child,
                        'Akses Kredit' as tag,        
                        sa.id as ref_id,
                        sa.is_approve,
                        to_char(sa.request_at, 'YYYY-MM-DD HH24:MI:SS'),
                        sa.note,
                        null,
                        null
                    FROM salesman_access sa 
					LEFT JOIN hr.employee e
						ON sa.requested_id = e.id
					WHERE e.user_id = {{.QUserId}} AND sa.access_type = 'kredit' AND DATE(sa.request_at) = {{.QDate}}

                UNION ALL

                SELECT 'public.customer_access',
                        NULL as ref_table_child,
                        'Double Kredit Customer' as tag,        
                        ca.id as ref_id,
                        ca.is_approve,
                        to_char(ca.request_at, 'YYYY-MM-DD HH24:MI:SS'),
                        ca.note,
                        ca.customer_id,
                        JSONB_AGG(
                            JSONB_BUILD_OBJECT(
                                'id', c.id||'',
                                'name', c.name,
                                'outlet_name', c.outlet_name,
                                'tipe', ct.name
                            )
                        ) as customers
                    FROM customer_access ca 
					LEFT JOIN hr.employee e
						ON ca.requested_id = e.id
                    JOIN customer c
                        ON c.id = ANY(ca.customer_id)
                    JOIN customer_type ct
                        ON c.tipe = ct.id
                    WHERE e.user_id = {{.QUserId}} AND ca.access_type = 'DOUBLE CREDIT' AND DATE(ca.request_at) = {{.QDate}}
                    GROUP BY ca.id

                UNION ALL

                SELECT 'public.customer_access_visit_extra',
                        NULL as ref_table_child,
                        'Visit Extra' as tag,        
                        ca.id as ref_id,
                        ca.is_approve,
                        to_char(ca.request_at, 'YYYY-MM-DD HH24:MI:SS'),
                        ca.note,
                        ca.customer_id,
                        JSONB_AGG(
                            JSONB_BUILD_OBJECT(
                                'id', c.id||'',
                                'name', c.name,
                                'outlet_name', c.outlet_name,
                                'tipe', ct.name
                            )
                        ) as customers
                    FROM customer_access_visit_extra ca
					LEFT JOIN hr.employee e
						ON ca.requested_id = e.id
                    JOIN customer c
                        ON c.id = ANY(ca.customer_id)
                    JOIN customer_type ct
                        ON c.tipe = ct.id
                    WHERE e.user_id = {{.QUserId}} AND ca.access_type = 'VISIT EXTRA' AND DATE(ca.request_at) = {{.QDate}}
                    GROUP BY ca.id

                UNION ALL

                SELECT 'public.customer_type_request', 
                        NULL as ref_table_child,
                        'Perubahan Tipe Customer' as tag,
                        ctr.id as ref_id,
                        ctr.is_approve,
                        to_char(ctr.request_at, 'YYYY-MM-DD HH24:MI:SS'),
                        ctr.note,
                        ctr.customer_id,
                        JSONB_AGG(
                            JSONB_BUILD_OBJECT(
                                'id', c.id||'',
                                'name', c.name,
                                'outlet_name', c.outlet_name,
                                'tipe', ct.name
                            )
                        ) as customers
                    FROM customer_type_request ctr 
					LEFT JOIN hr.employee e
						ON ctr.requested_id = e.id
                    JOIN customer c
                        ON c.id = ANY(ctr.customer_id)
                    JOIN customer_type ct
                        ON c.tipe = ct.id
                    WHERE e.user_id = {{.QUserId}} AND DATE(ctr.request_at) = {{.QDate}}
                    GROUP BY ctr.id

                UNION ALL

                SELECT 'public.delete_kunjungan_request',
                        NULL as ref_table_child,
                        'Delete Data Kunjungan',
                        dkr.id as ref_id,
                        dkr.is_approve,
                        to_char(dkr.datetime, 'YYYY-MM-DD HH24:MI:SS'),
                        dkr.note,
                        ARRAY_AGG(k.customer_id),
                        JSONB_AGG(
                                JSONB_BUILD_OBJECT(
                                    'id', c.id||'',
                                    'name', c.name,
                                    'outlet_name', c.outlet_name,
                                    'tipe', ct.name
                                )
                            ) as customers
                    FROM delete_kunjungan_request dkr 
					LEFT JOIN hr.employee e
						ON dkr.employee_id = e.id
                    LEFT JOIN kunjungan k
                        ON dkr.kunjungan_id = k.id
                    LEFT JOIN customer c
                        ON k.customer_id = c.id
                    LEFT JOIN customer_type ct
                        ON c.tipe = ct.id
                    WHERE e.user_id = {{.QUserId}} AND DATE(dkr.datetime) = {{.QDate}}
                    GROUP BY dkr.id

                UNION ALL

                SELECT 'public.customer_relocation',
                        NULL as ref_table_child,
                        'Perubahan Titik Customer',
                        cr.id as ref_id,
                        cr.is_approve,
                        to_char(cr.request_at, 'YYYY-MM-DD HH24:MI:SS'),
                        cr.note,
                        ARRAY_AGG(cr.customer_id),
                        JSONB_AGG(
                                JSONB_BUILD_OBJECT(
                                    'id', c.id||'',
                                    'name', c.name,
                                    'outlet_name', c.outlet_name,
                                    'tipe', ct.name
                                )
                            ) as customers
                    FROM customer_relocation cr
					LEFT JOIN hr.employee e
						ON cr.employee_id = e.id
                    JOIN customer c
                        ON c.id = cr.customer_id
                    JOIN customer_type ct
                        ON c.tipe = ct.id
                    WHERE e.user_id = {{.QUserId}} AND DATE(cr.request_at) = {{.QDate}}
                    GROUP BY cr.id

                UNION ALL

                SELECT 'public.salesman_request',
                        NULL as ref_table_child,
                        'Akses Login Salesman',
                        sr.id as ref_id,
                        sr.is_approve,
                        to_char(sr.request_at, 'YYYY-MM-DD HH24:MI:SS'),
                        null,
                        null,
                        null
                    FROM salesman_request sr
					LEFT JOIN hr.employee e
						ON sr.requested_id = e.id
					WHERE e.user_id = {{.QUserId}} AND DATE(sr.request_at) = {{.QDate}}

                UNION ALL

                SELECT 'public.salesman_request_so',
                        NULL as ref_table_child,
                        'Akses Buka SO Salesman',
                        sro.id as ref_id,
                        sro.is_approve,
                        to_char(sro.request_at, 'YYYY-MM-DD HH24:MI:SS'),
                        null,
                        null,
                        null
                    FROM salesman_request_so sro
					LEFT JOIN hr.employee e
						ON sro.requested_id = e.id
					WHERE e.user_id = {{.QUserId}} AND DATE(sro.request_at) = {{.QDate}}

                UNION ALL

                SELECT 'public.salesman_access_kunjungan',
                        NULL as ref_table_child,
                        'Akses Kunjungan Salesman',
                        sak.id as ref_id,
                        sak.is_approve,
                        to_char(sak.request_at, 'YYYY-MM-DD HH24:MI:SS'),
                        sak.note,
                        null,
                        null
                    FROM salesman_access_kunjungan sak
					LEFT JOIN hr.employee e
						ON sak.requested_id = e.id
					WHERE e.user_id = {{.QUserId}} AND DATE(sak.request_at) = {{.QDate}}

                    UNION ALL

                    SELECT 'public.customer_plafon_over_request',
                            'publuc.customer_plafon_over_request_detail' as ref_table_child,
                            'Request Over Plafon Customer',
                            cpor.id as ref_id,
                            cpor.is_approve,
                            to_char(cpor.created_at, 'YYYY-MM-DD HH24:MI:SS'),
                            cpor.note,
                            ARRAY_AGG(cpord.customer_id),
                            JSONB_AGG(
                                    JSONB_BUILD_OBJECT(
                                        'id', c.id||'',
                                        'name', c.name,
                                        'outlet_name', c.outlet_name,
                                        'tipe', ct.name
                                    )
                                ) as customers
                        FROM customer_plafon_over_request cpor 
						LEFT JOIN hr.employee e
							ON cpor.requested_id = e.id
                        JOIN customer_plafon_over_request_detail cpord
                            ON cpor.id = cpord.customer_plafon_over_request_id
                        LEFT JOIN customer c
                            ON cpord.customer_id = c.id
                        LEFT JOIN customer_type ct
                            ON c.tipe = ct.id
                        WHERE e.user_id = {{.QUserId}} AND DATE(cpor.created_at) = {{.QDate}}
                        GROUP BY cpor.id`

	templateParamQuery := map[string]interface{}{
		"QUserId": *inputUser.UserId,
		"QDate":   whereDate,
	}

	query1, err := helpers.PrepareQuery(templateQuery, templateParamQuery)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	var wg sync.WaitGroup
	resultsChan := make(chan map[int][]map[string]interface{}, 2)

	queries := []string{
		query1,
		query1 + qPage + qLimit,
	}

	tempResults := make([][]map[string]interface{}, len(queries))

	for i, query := range queries {
		wg.Add(1)
		go helpers.ExecuteGORMQuery(query, resultsChan, i, &wg)
	}

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

func GetPermission(c *fiber.Ctx) error {

	userIdstr := c.Query("userId")

	userId, _ := strconv.Atoi(userIdstr)

	datas, _ := getPermissions(int32(userId), c)

	return c.Status(fiber.StatusOK).JSON(datas)
}

func SetStock(c *fiber.Ctx) error {

	stokUser := new(structs.StokUser)
	gudang := new(structs.Gudang)

	if err := c.BodyParser(stokUser); err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengambil input data",
			Success: false,
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengambil input data",
			Success: false,
		})
	}

	tx := db.DB.Begin()

	qWhere := ""
	if stokUser.UserIdSubtitute != 0 {
		qWhere = " AND user_id_subtitute = " + strconv.Itoa(int(stokUser.UserIdSubtitute))
	} else {
		qWhere = " AND user_id_subtitute = 0"
	}

	if stokUser.TanggalStok == "" {
		stokUser.TanggalStok = fmt.Sprintf("%s", time.Now().Format("2006-01-02"))
	}

	if err := tx.Where("branch_id = ? ", data["branchId"]).First(&gudang).Error; err != nil && err.Error() != "record not found" {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data gudang",
			Success: false,
		})
	}

	stokUser.GudangId = gudang.ID

	if err := tx.Where("DATE(tanggal_stok) = DATE(?) AND user_id = ? AND gudang_id = ? "+qWhere, stokUser.TanggalStok, stokUser.UserId, gudang.ID).First(&stokUser).Error; err != nil && err.Error() != "record not found" {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data stok",
			Success: false,
		})
	}

	templateParamQuery := map[string]interface{}{
		"QUserId":         stokUser.UserId,
		"QInsert":         ", user_id_subtitute",
		"QSelect":         ", " + strconv.Itoa(int(stokUser.UserIdSubtitute)) + " as user_id_subtitute",
		"QInsertParentID": ", stok_user_id",
		"QSelectParentID": ", " + strconv.Itoa(int(stokUser.ID)) + " as stok_user_id",
	}

	sUserID := strconv.Itoa(int(stokUser.UserId))
	sGudangID := strconv.Itoa(int(gudang.ID))
	var sUserIDSubtitute string
	if &stokUser.UserIdSubtitute != nil {
		sUserIDSubtitute = strconv.Itoa(int(stokUser.UserIdSubtitute))
	} else {
		sUserIDSubtitute = ""
	}

	tanggalStokStr := fmt.Sprintf("%s", stokUser.TanggalStok)
	datas, err := getStokParent(&sUserID, &tanggalStokStr, &sUserIDSubtitute, &sGudangID, c)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return err
	}

	if stokUser.ID != 0 {
		tx.Rollback()
		return c.Status(http.StatusOK).JSON(helpers.Response{
			Message: "Data sudah ada",
			Success: true,
			Data:    datas[0],
		})
	}

	// fmt.Println(stokUser)
	if err := tx.Create(&stokUser).Error; err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal menyimpan data",
			Success: false,
		})
	}

	templateQuery := `INSERT INTO stok_salesman (stok_gudang_id, user_id, produk_id, stok_awal, stok_akhir, tanggal_stok, confirm_key, is_complete, tanggal_so, so_admin_gudang_id, condition, pita {{.QInsert}} {{.QInsertParentID}})

                        SELECT 0 as stok_gudang_id, 
                                sq.user_id, 
                                sq.produk_id, 
                                SUM(sq.stok_awal) as stok_awal, 
                                SUM(sq.stok_akhir) as stok_akhir, 
                                NOW() as tanggal_stok,
                                NULL as confirm_key, 
                                0 as is_complete, 
                                NULL as tanggal_so, 
                                NULL as so_admin_gudang_id, 
                                condition, 
                                pita
                                {{.QSelect}}
                                {{.QSelectParentID}}
                        FROM (
                                SELECT user_id, produk_id, COALESCE(lss.stok_akhir,0) as stok_awal, COALESCE(lss.stok_akhir,0) as stok_akhir, condition, pita
                                FROM (
                                        SELECT DISTINCT ON (user_id, produk_id, condition, pita) 
                                                    id, 
                                                    user_id, 
                                                    produk_id, 
                                                    condition, 
                                                    pita, 
                                                    stok_awal, 
                                                    stok_akhir, 
                                                    confirm_key, 
                                                    is_complete, 
                                                    tanggal_so,
                                                    so_admin_gudang_id
                                        FROM stok_salesman 
                                        WHERE user_id = {{.QUserId}} AND DATE(tanggal_stok) < CURRENT_DATE 
                                        ORDER BY user_id, produk_id, condition, pita, tanggal_stok DESC
                                ) lss
                                WHERE stok_akhir <> 0
                        ) sq
                        GROUP BY sq.user_id, sq.produk_id, sq.condition, sq.pita
                        ORDER BY sq.produk_id, sq.condition, sq.pita`

	templateQueryMD := `INSERT INTO md.stok_merchandiser (stok_gudang_id, user_id, item_id, stok_awal, stok_akhir, tanggal_stok, confirm_key, is_complete, tanggal_so, so_admin_gudang_id {{.QInsert}} {{.QInsertParentID}})

                        SELECT 0 as stok_gudang_id, 
                                sq.user_id, 
                                sq.item_id, 
                                SUM(sq.stok_awal) as stok_awal, 
                                SUM(sq.stok_akhir) as stok_akhir, 
                                NOW() as tanggal_stok,
                                NULL as confirm_key, 
                                0 as is_complete, 
                                NULL as tanggal_so, 
                                NULL as so_admin_gudang_id
                                {{.QSelect}}
                                {{.QSelectParentID}}
                        FROM (
                                SELECT user_id, item_id, COALESCE(lss.stok_akhir,0) as stok_awal, COALESCE(lss.stok_akhir,0) as stok_akhir
                                FROM (
                                    SELECT DISTINCT ON (user_id, item_id) 
                                        id, 
                                        user_id, 
                                        item_id, 
                                        stok_awal, 
                                        stok_akhir, 
                                        confirm_key, 
                                        is_complete, 
                                        tanggal_so,
                                        so_admin_gudang_id
                                    FROM md.stok_merchandiser 
                                    WHERE user_id = {{.QUserId}} AND DATE(tanggal_stok) < CURRENT_DATE 
                                    ORDER BY user_id, item_id, tanggal_stok DESC
                                ) lss
                                WHERE stok_akhir <> 0
                        ) sq
                        GROUP BY sq.user_id, sq.item_id
                        ORDER BY sq.item_id`

	templateParamQuery = map[string]interface{}{
		"QUserId":         stokUser.UserId,
		"QInsert":         ", user_id_subtitute",
		"QSelect":         ", " + strconv.Itoa(int(stokUser.UserIdSubtitute)) + " as user_id_subtitute",
		"QInsertParentID": ", stok_user_id",
		"QSelectParentID": ", " + strconv.Itoa(int(stokUser.ID)) + " as stok_user_id",
	}

	query1, err := helpers.PrepareQuery(templateQuery, templateParamQuery)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	if err := tx.Exec(query1).Error; err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal menyimpan data",
			Success: false,
		})
	}
	// fmt.Println(query1)

	query2, err := helpers.PrepareQuery(templateQueryMD, templateParamQuery)

	if err := tx.Exec(query2).Error; err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal menyimpan data",
			Success: false,
		})
	}

	if err := tx.Where("DATE(tanggal_stok) = DATE(?) AND user_id = ? AND gudang_id = ? "+qWhere, stokUser.TanggalStok, stokUser.UserId, gudang.ID).First(&stokUser).Error; err != nil && err.Error() != "record not found" {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data stok",
			Success: false,
		})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan menyimpan data",
			Success: false,
		})
	}

	sUserID = strconv.Itoa(int(stokUser.UserId))
	sGudangID = strconv.Itoa(int(stokUser.GudangId))

	if &stokUser.UserIdSubtitute != nil {
		sUserIDSubtitute = strconv.Itoa(int(stokUser.UserIdSubtitute))
	} else {
		sUserIDSubtitute = ""
	}

	tanggalStokStr = fmt.Sprintf("%s", stokUser.TanggalStok)
	datas, err = getStokParent(&sUserID, &tanggalStokStr, &sUserIDSubtitute, &sGudangID, c)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data stok telah berhasil dibuat",
		Success: true,
		Data:    datas[0],
	})
}

func getPermissions(userId int32, c *fiber.Ctx) ([]map[string]interface{}, error) {

	if userId == 0 {
		userId = int32(helpers.ParseInt(c.Query("userId")))
	}

	templateQuery := ` WITH role_permission as (
                            SELECT sq.role_id,
                                    sq.role_name,
                                    sq.app_id, 
                                    sq.app_name,
                                    JSONB_AGG(sq.permissions ORDER BY sq.permission_id) as permissions
                            FROM (
                                    SELECT r.id as role_id,
                                            r.name as role_name,
                                            app.id as app_id,
                                            app.name as app_name,
                                            p.id as permission_id,
                                            JSONB_BUILD_OBJECT(
                                                'id', p.id,
                                                'name', p.name,
                                                'modules',
                                                    JSONB_AGG(
                                                        JSONB_BUILD_OBJECT(
                                                                'module_id', m.id,
                                                                'module_name', m.name
                                                        )
                                                    )
                                            ) as permissions
                                    FROM rpm.user_role ur
                                    JOIN rpm.role r
                                        ON ur.role_id = r.id
                                    JOIN rpm.permission_role pr
                                        ON r.id = pr.role_id
                                    JOIN rpm.permission p
                                        ON pr.permission_id = p.id
                                    JOIN public.app
                                        ON p.app_id = app.id
                                    JOIN rpm.module m
                                        ON m.id = ANY(pr.module_ids)
                                    WHERE ur.user_id = {{.QUserId}} AND p.app_id = 16
                                    GROUP BY r.id, app.id, p.id
                            ) sq
	                        GROUP BY sq.role_id, sq.role_name, sq.app_id, sq.app_name
                    ), subject_profiles as (
					        SELECT r.id as role_id,
									r.name as role_name, 
									sp.name as subject_profile_name, 
									st.name as subject_type_name,
									JSONB_AGG(
                                        JSONB_BUILD_OBJECT(
                                            'id', mr.id,
                                            'name', mr.name,
                                            'value',  CASE WHEN spd.value IS NOT NULL 
                                                        THEN rpm.jsonb_dyntype(spd.value, spd.type_data)
                                                        ELSE 
                                                            JSONB_BUILD_OBJECT(
                                                                'min', rpm.jsonb_dyntype(spd.value_min, spd.type_data),
                                                                'max', rpm.jsonb_dyntype(spd.value_max, spd.type_data)
                                                            )
                                                        END
                                            ) ORDER BY mr.id
                                    ) as master_info
                            FROM rpm.user_role ur
                            JOIN rpm.role r
                                ON ur.role_id = r.id
                            JOIN rpm.subject_profile sp
                                ON r.subject_profile = sp.name
                            JOIN rpm.subject_type st
                                ON sp.subject_type_id = st.id
                            JOIN rpm.subject_profile_detail spd
                                ON sp.id = spd.subject_profile_id
                            JOIN rpm.master_rule mr
                                ON spd.master_rule_id = mr.id
                                AND mr.id = ANY(st.master_rule_ids)
                            WHERE ur.user_id = {{.QUserId}}
                            GROUP BY r.id, sp.id, st.id
                    )

                    SELECT ur.user_id,
                            rp.app_id,
                            rp.app_name,
                            JSONB_AGG(
                                JSONB_BUILD_OBJECT(
                                    'role_name', r.name,
                                    'role_id', ur.role_id,
                                    'subject_profile_name', sp.subject_profile_name,
                                    'subject_type_name', sp.subject_type_name,
                                    'permission', rp.permissions,
                                    'master_info', sp.master_info
                                ) ORDER BY ur.role_id
                            ) as user_info
                    FROM rpm.user_role ur
                    JOIN rpm.role r
                        ON ur.role_id = r.id
                    LEFT JOIN subject_profiles sp
                        ON r.id = sp.role_id
                    LEFT JOIN role_permission rp
                        ON r.id = rp.role_id
                    WHERE ur.user_id  = {{.QUserId}} AND rp.app_id = 16
                    GROUP BY ur.user_id, rp.app_id, rp.app_name`

	templateParamQuery := map[string]interface{}{
		"QUserId": userId,
	}

	query1, err := helpers.PrepareQuery(templateQuery, templateParamQuery)

	if err != nil {
		fmt.Println(err)
		return nil, c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	datas, err := helpers.ExecuteQuery(query1)

	if err != nil {
		fmt.Println(err)
		return nil, c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika mengambil data",
			Success: false,
		})
	}

	if len(datas) == 0 {
		return nil, c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data tidak ditemukan",
			Success: true,
		})
	}

	return datas, c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data berhasil diambil",
		Success: true,
		Data:    datas,
	})

}

func GetCheckSO(c *fiber.Ctx) error {

	userId := c.Query("userId")
	branchId := c.Query("branchId")

	gudang := structs.Gudang{}

	if err := db.DB.Where("branch_id = ? ", branchId).First(&gudang).Error; err != nil && err.Error() != "record not found" {
		fmt.Println(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data gudang",
			Success: false,
		})
	}

	StokUser := structs.StokUser{}

	err := db.DB.
		Where("user_id = ? AND is_complete = 0 AND gudang_id = ?", userId, gudang.ID).
		Order("DATE(tanggal_stok) ASC").
		First(&StokUser).Error
	if err != nil && err.Error() != "record not found" {
		fmt.Println(err.Error())
	}

	var tempTanggalStok time.Time
	tempTanggalStok, _ = time.Parse("2006-01-02T15:04:05Z", StokUser.TanggalStok)

	returnData := make(map[string]interface{})

	if !tempTanggalStok.IsZero() {
		returnData["pending_so"] = fmt.Sprintf("%s", tempTanggalStok.Format("2006-01-02"))
	} else {
		returnData["pending_so"] = nil
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data berhasil diambil",
		Success: true,
		Data:    returnData,
	})
}

func GetCustomerTransaction(c *fiber.Ctx) error {
	type Input struct {
		CustomerID string  `json:"customerId"`
		DateStart  string  `json:"dateStart"`
		DateEnd    string  `json:"dateEnd"`
		Type       *string `json:"type"`
	}

	input := Input{}

	if err := c.QueryParser(&input); err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengambil input data",
			Success: false,
		})
	}

	var forType []string
	if input.Type != nil {
		switch *input.Type {
		case "penjualan":
			forType = append(forType, "penjualan")
			forType = append(forType, "penjualan_detail")
		case "pengembalian":
			forType = append(forType, "pengembalian")
			forType = append(forType, "pengembalian_detail")
		case "pembayaran_piutang":
			forType = append(forType, "pembayaran_piutang")
			forType = append(forType, "pembayaran_piutang_detail")
		case "payment":
			forType = append(forType, "payment")
		case "kunjungan":
			forType = append(forType, "kunjungan")
			forType = append(forType, "kunjungan_log")
		case "piutang":
			forType = append(forType, "piutang")
		case "md.transaction":
			forType = append(forType, "md_transaction")
			forType = append(forType, "md_transaction_detail")
		default:
			forType = append(forType, "penjualan")
			forType = append(forType, "penjualan_detail")
			forType = append(forType, "pengembalian")
			forType = append(forType, "pengembalian_detail")
			forType = append(forType, "pembayaran_piutang")
			forType = append(forType, "pembayaran_piutang_detail")
			forType = append(forType, "payment")
			forType = append(forType, "kunjungan")
			forType = append(forType, "kunjungan_log")
			forType = append(forType, "piutang")
			forType = append(forType, "md_transaction")
			forType = append(forType, "md_transaction_detail")
		}
	} else {
		forType = append(forType, "penjualan")
		forType = append(forType, "penjualan_detail")
		forType = append(forType, "pengembalian")
		forType = append(forType, "pengembalian_detail")
		forType = append(forType, "pembayaran_piutang")
		forType = append(forType, "pembayaran_piutang_detail")
		forType = append(forType, "payment")
		forType = append(forType, "kunjungan")
		forType = append(forType, "kunjungan_log")
		forType = append(forType, "piutang")
		forType = append(forType, "md_transaction")
		forType = append(forType, "md_transaction_detail")
	}

	var returnData []interface{}

	dynamicGet := models.GetStructTransactions(forType, input.CustomerID, input.DateStart, input.DateEnd)

	returnData = append(returnData, dynamicGet)

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data berhasil diambil",
		Success: true,
		Data:    returnData[0],
	})
}
