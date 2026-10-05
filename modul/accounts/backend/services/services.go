// Package services memuat aturan modul Accounts.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak pernah mengimpor handlers dan tidak
// pernah mengimpor modul lain - yang bersama datang dari `inti/`.
//
// Sumber aturan (keputusan work owner 04-10-2026; layar Pega SFAGIS Account dari tangkapan layar - TIDAK ada di
// korpus XML):
//   - "Insured Name" * dipilih dari organisasi `CLIENT` (`FLAG` `Org`): `INSUREDID` = `CLIENT.ID`, `INSUREDNAME` =
//     `CLIENT.NAME` saat disimpan; "Org ID" = `ORG-n` organisasi itu.
//   - "Group Business" * dari `BUSINESSGROUP`: `GROUPBUSINESSID` = `ID`, `GROUPBUSINESS` = `NOTE`.
//   - "Description" bebas; "Owner" = `CREATEOP` (akun pelaku saat dibuat) dan "Create Date" = `CREATEDATE`
//     (SYSDATE) - keduanya tidak berubah saat Edit. "Territory" dihapus.
//   - `ID` = `ASM-SFAGIS-WORK-ACCOUNT ACC-<SEQ_T_M_ACCOUNT>`; nomor yang sudah terpakai dilewati.
//   - Pasangan Insured + Group Business yang sama dengan akun LAIN ditolak (data lama memuat 1.243 pasangan ganda;
//     yang baru tidak ditambah).
//   - Nol hapus ("Delete tidak").
//
// Transaksi dibuka di sini (`inti.Dasar.DalamTransaksi`), satu per permintaan, nol COMMIT di SQL.
package services

import (
	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/accounts/backend/repository"
)

// Service adalah akar layanan Accounts.
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
