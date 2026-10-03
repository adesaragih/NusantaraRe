// Package services memuat aturan dagang modul Marketing Officer.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak pernah mengimpor handlers dan tidak
// pernah mengimpor modul lain - yang bersama datang dari `inti/`.
//
// Sumber aturan (rancangan 03-10-2026, keputusan work owner):
//   - Form Pega `InputMarketingOfficer` + `SaveMarketingOfficer_Act` + prosedur `PEGA_MARKETINGOFFICER` (NB FacIn,
//     RNW Fac In, Endorsment Fac In): leader = `CLIENTID2` `LEADER` dan `MOLEADER` nama sendiri; anggota =
//     `CLIENTID2` ID leader dan `MOLEADER` nama leader; `TEAMGROUP`/`BRANCHDETAILNAME` dari cabang; nama dan leader
//     wajib. Prosedurnya TIDAK dipanggil - logikanya ditiru di sini dan di repository.
//   - `AKSES_LOGIN` = `M_LOGIN_GO.LOGIN_ID` (di Pega Operator ID; `SendEmailPolicy` mencari email MO dan leader).
//   - `CLIENTID` (Marketing Code) = `M_LOGIN_GO.CONTACT_ID`, atau `CLIENTID` lama bila akun itu sudah punya baris
//     MO; tidak pernah berubah sesudah dibuat. `CLIENTNAME` disalin dari `M_LOGIN_GO.NAME` saat dibuat.
//   - Satu baris AKTIF per Marketing Code dan per akun (DEV: maks. satu aktif per CLIENTID); nol hapus - nonaktif =
//     `MOSTATUS` 2. Nol tabel baru, nol DDL.
//
// Transaksi dibuka di sini (`inti.Dasar.DalamTransaksi`), satu per permintaan, nol COMMIT di SQL.
package services

import (
	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/marketingofficer/backend/repository"
)

// Service adalah akar layanan Marketing Officer.
type Service struct {
	*inti.Dasar
}

// DariDasar membuat Service di atas akar bersama yang disetel `cmd/api`.
func DariDasar(d *inti.Dasar) *Service { return &Service{Dasar: d} }

// PunyaDatabase aman dipanggil pada nil.
func (s *Service) PunyaDatabase() bool { return s != nil && s.Dasar.PunyaDatabase() }

// LayananOracle menyusun Layanan di atas Oracle - satu-satunya penyusun yang dipakai handlers (handlers tidak
// mengimpor repository).
func LayananOracle(s *Service) *Layanan {
	return BaruLayanan(repository.Baru(s.DB()), s.DalamTransaksi)
}

// pastikan Gudang Oracle memenuhi antarmuka layanan.
var _ Gudang = (*repository.Gudang)(nil)

// dbTx - alias supaya tanda tangan antarmuka terbaca.
type dbTx = db.Tx
