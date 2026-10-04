// Package models memuat tipe data modul NB Fac In yang dipakai lintas lapisan
// (repository → services → handlers).
package models

import "github.com/cockroachdb/apd/v3"

// BarisLimitA - satu baris tabel limit akseptasi bentuk A (lima tabel
// `POOLDATA.M_LIMIT_*`, tiket 11), hanya kolom yang menentukan tangga. Nilai
// `NUMBER` diurai desimal; NULL = nil, bukan nol (ADR-U-0027).
type BarisLimitA struct {
	Tabel        string // nama tabel asal, mis. M_LIMIT_PROPERTYY
	Jabatan      string // JABATAN
	TeamGroup    string // TEAM_GROUP
	LimitBottom  *apd.Decimal
	LimitBottom2 *apd.Decimal
}

// BarisLimitB - satu baris `POOLDATA.M_LIMIT_FINANCIALINS` (bentuk B, tiket 12).
type BarisLimitB struct {
	Jabatan        string // JABATAN (berspasi, mis. "DIREKTUR TEKNIK")
	LimitBond      *apd.Decimal
	LimitCreditCL  *apd.Decimal
	LimitCreditNCL *apd.Decimal
}
