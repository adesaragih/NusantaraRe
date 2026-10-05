// Package services memuat aturan dagang modul Aggregate.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak pernah mengimpor handlers dan tidak
// pernah mengimpor modul lain - yang bersama datang dari `inti/`.
//
// Sumber aturan: folder korpus `D:\XML\RNM_BRD\Aggregate` (kelas ASM-FW-GISFW-Int-AGGREGATE) -
// `ShowAggregateList` (layar unggah), `GridDasbordAgg` (daftar), `ChooseMasterID` + `GetMasterIDAgg_Act` +
// `SetMasterID` (pilih Master ID), `UploadCSVAggregate_Act` (baca CSV), `SaveAggregate_Act` (simpan). Keputusan work
// owner 04-10-2026: modul Aggregate di kelompok MASTER TREATY; langkah XML diikuti apa adanya (Treaty Year dari
// periode As At, baris periode pertama tanpa urutan, beberapa Master ID: share dijumlah dan ID digabung `;`);
// COMMENCEMENT tidak diisi; chart batang bertingkat dirancang sendiri.
//
// Transaksi dibuka di sini (`inti.Dasar.DalamTransaksi`), satu per permintaan, nol COMMIT di SQL.
package services

import (
	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/aggregate/backend/repository"
)

// Service adalah akar layanan Aggregate.
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
