package controllers

import (
	"fmt"
	"pluto_remastered/helpers"

	"github.com/gofiber/fiber/v2"
)

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

	returnData, err := helpers.NewExecuteQuery(query1)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika eksekusi query",
			Success: false,
		})
	}

	if len(returnData) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data stok tidak ditemukan",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data stok berhasil diambil",
		Data:    returnData,
		Success: true,
	})
}
