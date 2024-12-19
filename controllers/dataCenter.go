package controllers

import (
	"fmt"
	"pluto_remastered/helpers"
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

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Data:    datas,
		Message: "Berhasil mendapatkan data",
		Success: true,
	})
}
