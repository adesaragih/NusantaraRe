package repository

// Sumber tiga dropdown layar Register - A3 kelompok Register.
//
// Untuk apa berkas ini: mengisi ketiga isian ber-autocomplete di
// `Section/InputRegisterClaimLife.xml`, dari tabel yang report definition
// korpus sebut - bukan dari daftar yang dikarang.
//
// `[terverifikasi]` kelas dan medan tiap report definition:
//
//	ReportDefinition/BrowseCedingCoLife_RD.xml    `ASM-FW-GISFW-Int-AGENT`
//	  medan `.ID` `.ClientName` `.Leader0` `.ChildCount`
//	ReportDefinition/BrowseBusinessLife_RD.xml    `ASM-FW-GISFW-Int-BUSINESS`
//	  medan `.Note` `.NoteINA` `.ContentNote` `.PolicyCost` `.MaxDisc`
//	ReportDefinition/BrowseMarketingOfficer_RD.xml
//	  `ASM-FW-GISFW-Int-marketingofficer`
//	  medan `.ID` `.ClientID` `.BranchDetailID` `.BranchDetailName`
//	  `.MOLeader` `.MOStatus`
//
// `[data DBA — katalog DEV 27-09-2026]` ketiganya tabel nyata di `POOLDATA`:
// `AGENT` 28 kolom 429 baris, `BUSINESS` 18 kolom 177 baris,
// `MARKETINGOFFICER` 14 kolom 74 baris. Nama kolom = nama medan RD dalam
// HURUF BESAR.
//
// ⛔ DATA ORANG DAN MITRA. `AGENT.CLIENTNAME` dan
// `MARKETINGOFFICER.CLIENTNAME` adalah nama; keduanya dibaca SAAT JALAN untuk
// mengisi dropdown, dan NOL baris disalin ke fixture, tiket, maupun log.
//
// ⛔ SETIAP kueri BERBATAS. Dropdown tanpa batas menarik 429 baris ke layar
// pada setiap pengetikan huruf - dan tabel yang tumbuh membuatnya makin buruk
// tanpa satu pun galat.
//
// Dibaca sesudah: pesertapolis.go.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// batasRujukan adalah plafon baris satu pembacaan dropdown.
//
// ⚠️ Bukan ukuran halaman: dropdown tidak berhalaman. Ia plafon supaya
// pencarian yang terlalu luas tetap menjawab cepat, dan pemakai mempersempit
// ketikannya alih-alih menunggu.
const batasRujukan = 50

// BarisRujukan adalah satu pilihan dropdown.
//
// ⚠️ DUA medan saja - pengenal dan yang tampil. Kolom lain di ketiga tabel itu
// tidak dipakai layar Register, dan kolom yang dibaca cenderung ikut tercatat
// di suatu tempat pada akhirnya.
type BarisRujukan struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

// sqlRujukan merakit pembacaan berbatas sebuah tabel rujukan.
//
// ⛔ Nama tabel dan kolom datang dari PEMANGGIL yang sudah memvalidasinya,
// bukan dari permintaan HTTP. Menyusun nama objek dari masukan pemakai adalah
// injeksi lewat pintu yang tidak dijaga `:1`.
func sqlRujukan(tabel, kolomID, kolomNama string) string {
	return fmt.Sprintf(`SELECT %s, %s FROM %s
		 WHERE UPPER(%s) LIKE UPPER(:1)
		 ORDER BY %s
		 FETCH FIRST %d ROWS ONLY`,
		kolomID, kolomNama, tabel, kolomNama, kolomNama, batasRujukan)
}

// bacaRujukan menjalankan satu pembacaan dropdown.
func (r *PohonKlaim) bacaRujukan(ctx context.Context,
	objek, kolomID, kolomNama, cari string) ([]BarisRujukan, error) {

	tabel, err := r.db.Qualify(objek)
	if err != nil {
		return nil, err
	}
	q := sqlRujukan(tabel, kolomID, kolomNama)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	// ⚠️ Pola LIKE dirakit di Go, bukan di SQL: menyusunnya di teks SQL
	// berarti tanda `%` pemakai menjadi bagian pernyataan.
	pola := "%" + strings.TrimSpace(cari) + "%"
	baris, err := r.db.sql.QueryContext(ctx, q, pola)
	if err != nil {
		// ⛔ Galat dibungkus TANPA barisnya: pesan driver dapat memuat nilai
		// kolom, dan kolom ini memuat nama orang.
		return nil, fmt.Errorf("repository: membaca rujukan %s: %w", objek, err)
	}
	defer func() { _ = baris.Close() }()

	var hasil []BarisRujukan
	for baris.Next() {
		var id, nama sql.NullString
		if err := baris.Scan(&id, &nama); err != nil {
			return nil, fmt.Errorf("repository: membaca baris rujukan: %w", err)
		}
		hasil = append(hasil, BarisRujukan{ID: id.String, Nama: nama.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: menelusuri rujukan: %w", err)
	}
	return hasil, nil
}

// CariCeding membaca pilihan isian `Ceding`.
//
// `[terverifikasi]` `BrowseCedingCoLife_RD.xml` kelas `…Int-AGENT`, medan
// `.ID` dan `.ClientName`.
func (r *PohonKlaim) CariCeding(ctx context.Context, cari string) ([]BarisRujukan, error) {
	return r.bacaRujukan(ctx, "AGENT", "ID", "CLIENTNAME", cari)
}

// CariBisnis membaca pilihan isian `Class of Business`.
//
// `[terverifikasi]` `BrowseBusinessLife_RD.xml` kelas `…Int-BUSINESS`.
//
// ⚠️ Yang ditampilkan `NOTE`, bukan `NOTEINA`: keduanya ada di tabelnya, dan
// mana yang layar Pega tampilkan belum terbaca dari pohon section. Dicatat
// `[terbuka]` di `PARITAS-LAYAR-DAN-AKSI.md`, bukan ditebak diam-diam.
func (r *PohonKlaim) CariBisnis(ctx context.Context, cari string) ([]BarisRujukan, error) {
	return r.bacaRujukan(ctx, "BUSINESS", "ID", "NOTE", cari)
}

// CariMarketing membaca pilihan isian `Marketing Officer`.
//
// `[terverifikasi]` `BrowseMarketingOfficer_RD.xml` kelas
// `…Int-marketingofficer`, medan `.ID` dan `.ClientName`.
func (r *PohonKlaim) CariMarketing(ctx context.Context, cari string) ([]BarisRujukan, error) {
	return r.bacaRujukan(ctx, "MARKETINGOFFICER", "ID", "CLIENTNAME", cari)
}
