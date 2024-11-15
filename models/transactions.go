package models

import (
	db "pluto_remastered/config"
	"pluto_remastered/structs"
)

// func GetStructTransactions(tableName string, customerId string, dateStart string, dateEnd string) (map[string]interface{}, error) {
// 	// Define struct types with their corresponding is_parent values
// 	tableStructMap := map[string][]struct {
// 		Datas interface{}
// 	}{
// 		structs.TableNamePenjualan: {
// 			{db.DB.Debug().
// 				Where("customer_id = ? AND DATE(tanggal_penjualan) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
// 				Find(&[]structs.Penjualan{})},
// 			{db.DB.Debug().
// 				Where("penjualan.customer_id = ? AND DATE(tanggal_penjualan) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
// 				Joins("JOIN penjualan ON penjualan.id = penjualan_detail.penjualan_id").
// 				Find(&[]structs.PenjualanDetail{})},
// 		},
// 		structs.TableNamePengembalian: {
// 			{db.DB.
// 				Where("customer_id = ? AND DATE(tanggal_pengembalian) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
// 				Find(&[]structs.Pengembalian{})},
// 			{db.DB.
// 				Where("pengembalian.customer_id = ? AND DATE(tanggal_pengembalian) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
// 				Joins("JOIN pengembalian ON pengembalian.id = pengembalian_detail.pengembalian_id").
// 				Find(&[]structs.PengembalianDetail{})},
// 		},
// 		structs.TableNamePembayaranPiutang: {
// 			{db.DB.
// 				Where("customer_id = ? AND DATE(tanggal_pembayaran) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
// 				Find(&[]structs.PembayaranPiutang{})},
// 			{db.DB.
// 				Where("pembayaran_piutang.customer_id = ? AND DATE(tanggal_pembayaran) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
// 				Joins("pembayaran_piutang pp ON pp.id = pembayaran_piutang_detail.pembayaran_piutang_id").
// 				Find(&[]structs.PembayaranPiutangDetail{})},
// 		},
// 		structs.TableNamePayment: {
// 			{db.DB.
// 				Where("customer_id = ? AND DATE(tanggal_transaksi) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
// 				Find(&[]structs.Payment{})},
// 		},
// 		structs.TableNameKunjungan: {
// 			{db.DB.
// 				Where("customer_id = ? AND DATE(tanggal_kunjungan) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
// 				Find(&[]structs.Kunjungan{})},
// 			{db.DB.
// 				Where("customer_id = ? AND DATE(checkin_at) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
// 				Find(&[]structs.KunjunganLog{})},
// 		},
// 		structs.TableNamePiutang: {
// 			{db.DB.
// 				Where("customer_id = ? AND DATE(tanggal_piutang) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
// 				Find(&[]structs.Piutang{})},
// 		},
// 		structs.TableNameMdTransaction: {
// 			{db.DB.
// 				Where("customer_id = ? AND DATE(datetime) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
// 				Find(&[]structs.MdTransaction{})},
// 			{db.DB.
// 				Where("customer_id = ? AND DATE() BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
// 				Joins("JOIN md_transaction ON md_transaction.id = md_transaction_detail.md_transaction_id").
// 				Find(&[]structs.MdTransactionDetail{})},
// 		},
// 	}

// 	if structDefs, exists := tableStructMap[tableName]; exists {
// 		result := make([]interface{}, len(structDefs))

// 		for i, def := range structDefs {
// 			result[i] = def.Datas
// 		}
// 		return map[string]interface{}{tableName: result}, nil
// 	}

// 	return nil, fmt.Errorf("no struct found for table name: %s", tableName)
// }

func GetStructTransactions(tableName []string, customerId string, dateStart string, dateEnd string) map[string]interface{} {
	// Define struct types with their corresponding queries and data storage
	tableStructMap := map[string][]interface{}{
		structs.TableNamePenjualan: {
			&[]structs.Penjualan{}, // Parent table
		},
		structs.TableNamePenjualanDetail: {
			&[]structs.PenjualanDetail{},
		},
		structs.TableNamePengembalian: {
			&[]structs.Pengembalian{},
		},
		structs.TableNamePengembalianDetail: {
			&[]structs.PengembalianDetail{},
		},
		structs.TableNamePembayaranPiutang: {
			&[]structs.PembayaranPiutang{},
		},
		structs.TableNamePembayaranPiutangDetail: {
			&[]structs.PembayaranPiutangDetail{},
		},
		structs.TableNamePayment: {
			&[]structs.Payment{},
		},
		structs.TableNameKunjungan: {
			&[]structs.Kunjungan{},
		},
		structs.TableNameKunjunganLog: {
			&[]structs.KunjunganLog{},
		},
		structs.TableNamePiutang: {
			&[]structs.Piutang{},
		},
		structs.TableNameMdTransaction: {
			&[]structs.MdTransaction{},
		},
		structs.TableNameMdTransactionDetail: {
			&[]structs.MdTransactionDetail{},
		},
	}

	returnData := make(map[string]interface{})
	for _, tableName := range tableName {
		if structDefs, exists := tableStructMap[tableName]; exists {
			result := make([]interface{}, 0, len(structDefs))

			// Iterate over each table and fetch the data
			for _, dataStruct := range structDefs {
				switch v := dataStruct.(type) {
				case *[]structs.Penjualan:
					db.DB.
						Where("customer_id = ? AND DATE(tanggal_penjualan) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
						Order("tanggal_penjualan DESC").
						Find(v)
				case *[]structs.PenjualanDetail:
					db.DB.
						Where("penjualan.customer_id = ? AND DATE(tanggal_penjualan) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
						Order("tanggal_penjualan DESC").
						Joins("JOIN penjualan ON penjualan.id = penjualan_detail.penjualan_id").
						Find(v)
				case *[]structs.Pengembalian:
					db.DB.
						Where("customer_id = ? AND DATE(tanggal_pengembalian) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
						Order("tanggal_pengembalian DESC").
						Find(v)
				case *[]structs.PengembalianDetail:
					db.DB.
						Where("pengembalian.customer_id = ? AND DATE(tanggal_pengembalian) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
						Order("tanggal_pengembalian DESC").
						Joins("JOIN pengembalian ON pengembalian.id = pengembalian_detail.pengembalian_id").
						Find(v)
				case *[]structs.PembayaranPiutang:
					db.DB.
						Where("customer_id = ? AND DATE(tanggal_pembayaran_piutang) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
						Order("tanggal_pembayaran_piutang DESC").
						Find(v)
				case *[]structs.PembayaranPiutangDetail:
					db.DB.
						Where("pembayaran_piutang.customer_id = ? AND DATE(tanggal_pembayaran_piutang) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
						Order("tanggal_pembayaran_piutang DESC").
						Joins("JOIN pembayaran_piutang ON pembayaran_piutang.id = pembayaran_piutang_detail.pembayaran_piutang_id").
						Find(v)
				case *[]structs.Payment:
					db.DB.
						Where("customer_id = ? AND DATE(tanggal_transaksi) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
						Order("tanggal_transaksi DESC").
						Find(v)
				case *[]structs.Kunjungan:
					db.DB.
						Where("customer_id = ? AND DATE(tanggal_kunjungan) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
						Order("tanggal_kunjungan DESC").
						Find(v)
				case *[]structs.KunjunganLog:
					db.DB.
						Where("customer_id = ? AND DATE(checkin_at) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
						Order("checkin_at DESC").
						Find(v)
				case *[]structs.Piutang:
					db.DB.
						Where("customer_id = ? AND DATE(tanggal_piutang) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
						Order("tanggal_piutang DESC").
						Find(v)
				case *[]structs.MdTransaction:
					db.DB.
						Where("customer_id = ? AND DATE(datetime) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
						Order("datetime DESC").
						Find(v)
				case *[]structs.MdTransactionDetail:
					db.DB.
						Where("md_transaction.customer_id = ? AND DATE(datetime) BETWEEN DATE(?) AND DATE(?)", customerId, dateStart, dateEnd).
						Order("datetime DESC").
						Joins("JOIN md_transaction ON md_transaction.id = md_transaction_detail.md_transaction_id").
						Find(v)
				}

				// Append the fetched data to the result slice
				result = append(result, dataStruct)
			}
			// return map[string]interface{}{tableName: result}, nil
			returnData[tableName] = result[0]
		}
	}
	return returnData
}
