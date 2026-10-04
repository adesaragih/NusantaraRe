// Package services memuat aturan dagang modul Company Detail.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak pernah mengimpor handlers dan tidak
// pernah mengimpor modul lain - yang bersama datang dari `inti/`.
//
// Sumber aturan (diskusi work owner 03/04-10-2026; layar SFAGIS Company Detail tidak ada di korpus XML - acuannya
// screenshot Pega dan katalog DEV):
//   - Tanpa dokumen Pega: `CLIENT`, `CLIENT_PICLIST`, `CLIENT_ADDRESS` adalah sumber datanya ("aku tidak mau ada json
//     lagi"); prosedur `RDBINSERTCLIENT` TIDAK dipanggil, kolom barisnya ditiru.
//   - Layar: Company Detail (NPWP, Parent organization, COUNTRY*, Title, Organization Name*, Business Field*, Note),
//     grid PIC (Name*, Position*, Gender, Email, Date of birth, Phone number), grid Address (Type, Address, Phone and
//     Fax), tombol Create. Client Status, Established Date, Number of Employees, dan Owner DIHAPUS (work owner).
//   - Nomor ORG baru dari sequence `SEQ_CLIENT_ORG`; migrasi 810 memulainya dari nomor ORG tertinggi Pega (CLIENT
//     dan M_CLIENT) + 1 ("buat seq aja, start-nya dari id max+1"). Tanpa indeks unik (work owner: "jangan bantah").
//   - COUNTRY dari `NATION` (OLDID, seperti seluruh data lama); dropdown lain dari `M_ENUMERASI`.
//   - Phone and Fax: satu baris `CLIENT_ADDRESS` per nomor ("CLIENT_ADDRESS itu kan list bisa banyak row").
//   - Baris PIC dan alamat boleh dihapus dari grid (Pega pun menghapusnya lewat trigger); organisasi tidak dihapus.
//
// Transaksi dibuka di sini (`inti.Dasar.DalamTransaksi`), satu per permintaan, nol COMMIT di SQL.
package services

import (
	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/companydetail/backend/repository"
)

// Service adalah akar layanan Company Detail.
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
