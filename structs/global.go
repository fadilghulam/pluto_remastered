package structs

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

func GetStructInstanceByTableName(tableName string) (interface{}, error) {
	// Map table names to struct types
	tableStructMap := map[string]reflect.Type{
		TableNameCustomerHistory:           reflect.TypeOf([]CustomerHistory{}),
		TableNameCustomerKtp:               reflect.TypeOf([]CustomerKtp{}),
		TableNameCustomerScoreRaw:          reflect.TypeOf([]CustomerScoreRaw{}),
		TableNameCustomerTokoku:            reflect.TypeOf([]CustomerTokoku{}),
		TableNameCustomer:                  reflect.TypeOf([]Customer{}),
		TableNameKunjunganLog:              reflect.TypeOf([]KunjunganLog{}),
		TableNameKunjungan:                 reflect.TypeOf([]Kunjungan{}),
		TableNameLocationLog:               reflect.TypeOf([]LocationLog{}),
		TableNamePayment:                   reflect.TypeOf([]Payment{}),
		TableNamePembayaranPiutangDetail:   reflect.TypeOf([]PembayaranPiutangDetail{}),
		TableNamePembayaranPiutang:         reflect.TypeOf([]PembayaranPiutang{}),
		TableNamePengembalianDetail:        reflect.TypeOf([]PengembalianDetail{}),
		TableNamePengembalian:              reflect.TypeOf([]Pengembalian{}),
		TableNamePenjualanDetail:           reflect.TypeOf([]PenjualanDetail{}),
		TableNamePenjualan:                 reflect.TypeOf([]Penjualan{}),
		TableNamePiutang:                   reflect.TypeOf([]Piutang{}),
		TableNameQrCodeHistory:             reflect.TypeOf([]QrCodeHistory{}),
		TableNameRemainProductDetail:       reflect.TypeOf([]RemainProductDetail{}),
		TableNameRemainProduct:             reflect.TypeOf([]RemainProduct{}),
		TableNameSurveyProdukKompetitor:    reflect.TypeOf([]SurveyProdukKompetitor{}),
		TableNameSurveyProgramKompetitor:   reflect.TypeOf([]SurveyProgramKompetitor{}),
		TableNameValidateCustomer:          reflect.TypeOf([]ValidateCustomer{}),
		TableNameValidateKunjungan:         reflect.TypeOf([]ValidateKunjungan{}),
		TableNameValidatePembayaranPiutang: reflect.TypeOf([]ValidatePembayaranPiutang{}),
		TableNameValidatePengembalian:      reflect.TypeOf([]ValidatePengembalian{}),
		TableNameValidatePenjualan:         reflect.TypeOf([]ValidatePenjualan{}),
		TableNameValidateTransaksi:         reflect.TypeOf([]ValidateTransaksi{}),
		TableNameMdTransactionDetail:       reflect.TypeOf([]MdTransactionDetail{}),
		TableNameMdTransaction:             reflect.TypeOf([]MdTransaction{}),
		TableNameMdOutlet:                  reflect.TypeOf([]MdOutlet{}),
		// TableNameCustomerMoveRequest: reflect.TypeOf([]CustomerMoveRequestHiperion{}),
		TableNameCheckinRequest:            reflect.TypeOf([]CheckinRequest{}),
		TableNameCustomerAccessVisitExtra:  reflect.TypeOf([]CustomerAccessVisitExtra{}),
		TableNameCustomerAccess:            reflect.TypeOf([]CustomerAccess{}),
		TableNameCustomerMoveRequest:       reflect.TypeOf([]CustomerMoveRequest{}),
		TableNameCustomerPlafonOverRequest: reflect.TypeOf([]CustomerPlafonOverRequest{}),
		TableNameCustomerRelocation:        reflect.TypeOf([]CustomerRelocation{}),
		TableNameCustomerTypeRequest:       reflect.TypeOf([]CustomerTypeRequest{}),
		TableNameDeleteKunjunganRequest:    reflect.TypeOf([]DeleteKunjunganRequest{}),
		TableNameRuteMoveRequest:           reflect.TypeOf([]RuteMoveRequest{}),
		TableNameSalesmanAccessKunjungan:   reflect.TypeOf([]SalesmanAccessKunjungan{}),
		TableNameSalesmanAccess:            reflect.TypeOf([]SalesmanAccess{}),
		TableNameSalesmanRequestSo:         reflect.TypeOf([]SalesmanRequestSo{}),
		TableNameSalesmanRequest:           reflect.TypeOf([]SalesmanRequest{}),
	}

	if structType, exists := tableStructMap[tableName]; exists {
		// Create a new instance of the struct and return it
		return reflect.New(structType).Interface(), nil
	}
	return nil, fmt.Errorf("no struct found for table name: %s", tableName)
}

type Int32Array []int32

// Value converts Int32Array to a PostgreSQL array-compatible format.
func (a Int32Array) Value() (driver.Value, error) {
	// Convert []int32 to []interface{}
	var arr = make([]interface{}, len(a))
	for i, v := range a {
		arr[i] = v
	}
	return arr, nil
}

func (a Int32Array) Value2() (driver.Value, error) {
	// Convert to []int32 to []int64 for pq.Array
	int64Array := make([]int64, len(a))
	for i, v := range a {
		int64Array[i] = int64(v)
	}
	return pq.Array(int64Array), nil
}

// Scan converts a PostgreSQL array to Int32Array.
func (a *Int32Array) Scan(value interface{}) error {
	var ints []int32

	switch v := value.(type) {
	case string:
		// Handle the case where the array is returned as a string
		trimmed := strings.Trim(v, "{}")
		if len(trimmed) == 0 {
			*a = []int32{}
			return nil
		}
		strElements := strings.Split(trimmed, ",")
		for _, strElem := range strElements {
			i, err := strconv.Atoi(strElem)
			if err != nil {
				return err
			}
			ints = append(ints, int32(i))
		}
	case []byte:
		// Handle the case where the array is returned as []byte
		trimmed := strings.Trim(string(v), "{}")
		if len(trimmed) == 0 {
			*a = []int32{}
			return nil
		}
		strElements := strings.Split(trimmed, ",")
		for _, strElem := range strElements {
			i, err := strconv.Atoi(strElem)
			if err != nil {
				return err
			}
			ints = append(ints, int32(i))
		}
	default:
		return fmt.Errorf("unsupported data type: %T", v)
	}

	*a = ints
	return nil
}

type StringArray []string

// Value converts StringArray to a PostgreSQL array-compatible format.
func (a StringArray) Value() (driver.Value, error) {
	return fmt.Sprintf("{%s}", strings.Join(a, ",")), nil
}

// Scan converts a PostgreSQL array to StringArray.
func (a *StringArray) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		trimmed := strings.Trim(v, "{}")
		if len(trimmed) == 0 {
			*a = []string{}
			return nil
		}
		*a = strings.Split(trimmed, ",")
	case []byte:
		trimmed := strings.Trim(string(v), "{}")
		if len(trimmed) == 0 {
			*a = []string{}
			return nil
		}
		*a = strings.Split(trimmed, ",")
	default:
		return fmt.Errorf("unsupported data type: %T", v)
	}

	return nil
}
