package controllers

import (
	db "pluto_remastered/config"
	"pluto_remastered/helpers"
	"sync"

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
	db.DB.Exec(`INSERT INTO (user_id, user_id_subtitute, branch_id, start_date, end_date)
											SELECT sq.user_id,
													COALESCE(sq.user_id_subtitute, -1),
													sq.branch_id,
													sq.min_date,
													DATE(LEAD(sq.min_date, 1) OVER (PARTITION BY sq.user_id ORDER BY sq.user_id, sq.min_date)::date - INTERVAL '1 day')
											FROM (
												SELECT branch_id, salesman_id, merchandiser_id, teamleader_id, user_id, user_id_subtitute, MIN(DATE(tanggal_kunjungan)) as min_date
												FROM kunjungan 
												GROUP BY branch_id, salesman_id, merchandiser_id, teamleader_id, user_id, user_id_subtitute
											) sq
											WHERE sq.user_id IS NOT NULL
											ORDER BY sq.user_id, sq.min_date`)

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		Message: "Cron success",
		Success: true,
	})
}
