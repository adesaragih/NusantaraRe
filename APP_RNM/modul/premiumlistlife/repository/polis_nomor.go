package repository

// Penyimpanan `PL_NUMBER` - tiket 03 PremiumList Life.
//
// Untuk apa berkas ini: membaca bahan nomor premium list dari tabel header,
// membaca nomor yang sudah terbit, dan menuliskan nomor baru ke seluruh baris
// peserta polis itu.
//
// ⛔ NOMORNYA TINGGAL DI BARIS PESERTA, BUKAN DI HEADER. Di Pega ia properti
// halaman kerja (`.PremiumListSummary.PL_NUMBER`) dan ikut tersimpan di work
// object; di skema kami satu-satunya kolom bernama `PL_NUMBER` ada di tabel
// peserta (migrasi 052). Tabel ringkasan (055) TIDAK punya kolom itu, dan
// tabel header (051) hanya punya `PL_NUMBER_EDM` - yang lain sama sekali.
//
// Akibatnya satu, dan dinyatakan alih-alih disiasati: polis yang belum punya
// satu pun baris peserta TIDAK DAPAT DINOMORI, sebab nomornya tidak punya
// tempat tinggal. Itu bukan pembatasan yang kami karang - `SubmitPremiumList`
// memang berjalan SESUDAH rincian terisi, dan unggahan CSV (tiket 04) yang
// mengisinya. Menerbitkan nomor lebih dahulu berarti menggerakkan penghitung
// lalu membuang hasilnya.
//
// ⚠️ `[terbuka - pemilik kerja]` Bila kelak polis perlu bernomor sebelum
// pesertanya ada, kolom `PL_NUMBER` di tabel header-lah jawabannya - dan itu
// keputusan skema, yang menuntut migrasi baru dari keputusan yang TERCATAT.
//
// ⛔ Setiap query menyebut skemanya lewat `Qualify` (ADR-U-0033).
//
// Dibaca sesudah: penomor.go (penghitungnya), polis_inbox.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/db"
)

var (
	// ErrPolisTakDitemukan - tidak ada baris header dengan ID itu.
	ErrPolisTakDitemukan = errors.New(
		"repository: polis tidak ditemukan")

	// ErrPolisTanpaPeserta - polis belum punya baris peserta.
	//
	// ⛔ Pesannya MENYEBUT jalan keluarnya. Galat yang hanya berkata "tidak
	// dapat dinomori" membuat orang mencari kerusakan yang tidak ada.
	ErrPolisTanpaPeserta = errors.New(
		"repository: polis belum punya satu pun baris peserta, sehingga " +
			"nomor premium list belum punya tempat tersimpan; unggah rincian " +
			"peserta lebih dahulu")

	// ErrNomorPLBerbedaAntarPeserta - baris peserta tidak sepakat nomornya.
	//
	// ⛔ TIDAK dirapikan diam-diam. Dua nomor pada satu polis berarti salah
	// satu rombongan baris milik polis lain, atau satu penomoran pernah
	// berjalan separuh. Menimpanya dengan salah satu nomor membuang bukti
	// satu-satunya tentang mana yang benar.
	ErrNomorPLBerbedaAntarPeserta = errors.New(
		"repository: baris peserta satu polis memuat lebih dari satu PL_NUMBER")

	// ErrNomorPLTerbitBersamaan - permintaan lain menomori polis ini lebih
	// dahulu, di antara pembacaan gerbang dan penulisan nomor.
	//
	// ⚠️ Transaksinya DIBATALKAN saat galat ini kembali, jadi kenaikan
	// penghitung ikut batal dan nomor yang terlanjur dirakit tidak hilang
	// dari deret - ia tidak pernah terpakai sama sekali. Pemanggil cukup
	// membaca ulang nomornya.
	ErrNomorPLTerbitBersamaan = errors.New(
		"repository: PL_NUMBER polis ini baru saja diterbitkan permintaan lain; " +
			"baca ulang nomornya")
)

// IdentitasPolis adalah bahan nomor yang tinggal di tabel header.
type IdentitasPolis struct {
	// Tipe adalah kolom `TYPE` - QR / QP / TP / TR.
	Tipe string
	// KodeBisnis adalah kolom `BUSINESS_CODE` - COB.
	KodeBisnis string
}

// KeadaanNomorPL adalah apa yang sudah tersimpan di baris peserta.
type KeadaanNomorPL struct {
	// Nomor adalah nomor yang sudah terbit; kosong berarti belum.
	Nomor string
	// CacahPeserta adalah jumlah SELURUH baris peserta polis itu.
	CacahPeserta int
	// CacahBernomor adalah jumlah baris yang sudah memuat nomor.
	CacahBernomor int
}

// Bernomor menjawab apakah polis ini sudah punya nomor.
func (k KeadaanNomorPL) Bernomor() bool { return k.Nomor != "" }

// Utuh menjawab apakah SETIAP baris peserta sudah memuat nomornya.
//
// ⛔ `MIN` DAN `MAX` MELEWATI NULL, dan di situlah lubangnya. Polis yang
// separuh barisnya bernomor dan separuh lagi kosong menghasilkan `MIN == MAX`
// - yaitu "sudah bernomor" - sehingga gerbang lahir-sekali kembali lebih awal
// dan baris yang kosong TIDAK PERNAH terisi. Selisih `CacahBernomor` dengan
// `CacahPeserta` adalah satu-satunya yang menunjukkannya, dan selisih yang
// dibaca tetapi tidak pernah dibandingkan sama saja dengan tidak dibaca.
//
// ⚠️ Keadaan ini BUKAN kerusakan: unggahan CSV (tiket 04) dapat menambah
// peserta SESUDAH polis bernomor, dan peserta baru itu memang lahir tanpa
// nomor. Yang benar adalah memberi mereka nomor yang SUDAH ada - bukan nomor
// baru, dan bukan galat.
func (k KeadaanNomorPL) Utuh() bool { return k.CacahBernomor == k.CacahPeserta }

// NomorPolis membaca dan menuliskan `PL_NUMBER`.
type NomorPolis struct{ db *db.DB }

// NewNomorPolis menyusunnya.
func NewNomorPolis(db *db.DB) *NomorPolis { return &NomorPolis{db: db} }

// sqlIdentitasPolis merakit pembacaan bahan nomor dari header.
func sqlIdentitasPolis(polis string) string {
	return fmt.Sprintf(
		`SELECT p.TYPE, p.BUSINESS_CODE FROM %s p WHERE p.ID = :1`, polis)
}

// Identitas membaca tipe dan kode bisnis satu polis.
func (r *NomorPolis) Identitas(ctx context.Context, tx *db.Tx, polisID string) (
	IdentitasPolis, error) {

	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return IdentitasPolis{}, err
	}
	q := sqlIdentitasPolis(polis)
	if err := db.PeriksaSQL(q); err != nil {
		return IdentitasPolis{}, err
	}
	var tipe, kode sql.NullString
	if err := tx.QueryRowContext(ctx, q, polisID).Scan(&tipe, &kode); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return IdentitasPolis{}, fmt.Errorf("%w: %q", ErrPolisTakDitemukan, polisID)
		}
		return IdentitasPolis{}, fmt.Errorf("repository: membaca identitas polis: %w", err)
	}
	return IdentitasPolis{
		Tipe:       strings.TrimSpace(tipe.String),
		KodeBisnis: strings.TrimSpace(kode.String),
	}, nil
}

// sqlKeadaanNomorPL merakit pembacaan keadaan nomor.
//
// ⛔ `MIN` DAN `MAX` KEDUANYA DIBACA, bukan `MAX` saja. Satu polis bernomor
// satu; bila kedua ujungnya berbeda, barisnya tidak sepakat - dan pembacaan
// yang hanya mengambil `MAX` mengubur ketidaksepakatan itu di balik satu
// nilai yang kelihatan meyakinkan.
//
// ⚠️ `COUNT(d.PL_NUMBER)` mencacah baris yang kolomnya BUKAN null - berbeda
// dari `COUNT(*)`. Selisih keduanya yang memberi tahu penomoran pernah
// berjalan separuh.
func sqlKeadaanNomorPL(detail string) string {
	return fmt.Sprintf(`SELECT MIN(d.PL_NUMBER), MAX(d.PL_NUMBER),
	        COUNT(*), COUNT(d.PL_NUMBER)
	   FROM %s d WHERE d.PREMIUM_LIST_ID = :1`, detail)
}

// Keadaan membaca apakah polis sudah bernomor, dan seberapa utuh.
func (r *NomorPolis) Keadaan(ctx context.Context, tx *db.Tx, polisID string) (
	KeadaanNomorPL, error) {

	detail, err := r.db.Qualify("T_PREMIUM_LIST_DETAIL")
	if err != nil {
		return KeadaanNomorPL{}, err
	}
	q := sqlKeadaanNomorPL(detail)
	if err := db.PeriksaSQL(q); err != nil {
		return KeadaanNomorPL{}, err
	}
	var terkecil, terbesar sql.NullString
	var cacah, bernomor int
	if err := tx.QueryRowContext(ctx, q, polisID).Scan(
		&terkecil, &terbesar, &cacah, &bernomor); err != nil {
		return KeadaanNomorPL{}, fmt.Errorf(
			"repository: membaca keadaan nomor polis: %w", err)
	}
	kecil := strings.TrimSpace(terkecil.String)
	besar := strings.TrimSpace(terbesar.String)
	if kecil != besar {
		return KeadaanNomorPL{}, fmt.Errorf("%w: %q dan %q",
			ErrNomorPLBerbedaAntarPeserta, kecil, besar)
	}
	return KeadaanNomorPL{
		Nomor:         besar,
		CacahPeserta:  cacah,
		CacahBernomor: bernomor,
	}, nil
}

// RingkasPolis adalah kepala polis beserta keadaan nomornya, dalam satu baca.
type RingkasPolis struct {
	Identitas IdentitasPolis
	Nomor     KeadaanNomorPL
}

// sqlRingkasPolis merakit pembacaan kepala polis beserta nomornya.
//
// ⛔ SATU QUERY, BUKAN DUA. Dua pembacaan terpisah di luar transaksi dapat
// melihat dua keadaan yang berbeda - kepala polis sebelum penomoran, nomor
// sesudahnya - dan layar lalu menampilkan pasangan yang tidak pernah ada.
// Menyatukannya membuat transaksi pembaca tidak diperlukan sama sekali.
//
// ⛔ `LEFT JOIN`, bukan `JOIN`: polis yang belum punya satu pun peserta
// justru yang paling perlu terbaca - ia yang belum dapat dinomori, dan layar
// harus dapat mengatakannya. `JOIN` biasa membuatnya tampak tidak ada sama
// sekali, yaitu 404 untuk polis yang jelas ada.
//
// ⚠️ `COUNT(d.ID)`, bukan `COUNT(*)`. Atas `LEFT JOIN` tanpa pasangan,
// `COUNT(*)` mencacah baris hasil join - yaitu SATU - sedangkan `COUNT(d.ID)`
// mencacah peserta yang benar-benar ada, yaitu nol. Selisih itu persis
// gerbang "polis belum punya peserta".
func sqlRingkasPolis(polis, detail string) string {
	return fmt.Sprintf(`SELECT p.TYPE, p.BUSINESS_CODE,
	        MIN(d.PL_NUMBER), MAX(d.PL_NUMBER),
	        COUNT(d.ID), COUNT(d.PL_NUMBER)
	   FROM %s p LEFT JOIN %s d ON d.PREMIUM_LIST_ID = p.ID
	  WHERE p.ID = :1
	  GROUP BY p.TYPE, p.BUSINESS_CODE`, polis, detail)
}

// Ringkas membaca kepala polis beserta keadaan nomornya, tanpa transaksi.
func (r *NomorPolis) Ringkas(ctx context.Context, polisID string) (RingkasPolis, error) {
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return RingkasPolis{}, err
	}
	detail, err := r.db.Qualify("T_PREMIUM_LIST_DETAIL")
	if err != nil {
		return RingkasPolis{}, err
	}
	q := sqlRingkasPolis(polis, detail)
	if err := db.PeriksaSQL(q); err != nil {
		return RingkasPolis{}, err
	}
	var tipe, kode, terkecil, terbesar sql.NullString
	var cacah, bernomor int
	if err := r.db.QueryRowContext(ctx, q, polisID).Scan(
		&tipe, &kode, &terkecil, &terbesar, &cacah, &bernomor); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return RingkasPolis{}, fmt.Errorf("%w: %q", ErrPolisTakDitemukan, polisID)
		}
		return RingkasPolis{}, fmt.Errorf("repository: membaca ringkas polis: %w", err)
	}
	kecil := strings.TrimSpace(terkecil.String)
	besar := strings.TrimSpace(terbesar.String)
	if kecil != besar {
		return RingkasPolis{}, fmt.Errorf("%w: %q dan %q",
			ErrNomorPLBerbedaAntarPeserta, kecil, besar)
	}
	return RingkasPolis{
		Identitas: IdentitasPolis{
			Tipe:       strings.TrimSpace(tipe.String),
			KodeBisnis: strings.TrimSpace(kode.String),
		},
		Nomor: KeadaanNomorPL{
			Nomor:         besar,
			CacahPeserta:  cacah,
			CacahBernomor: bernomor,
		},
	}, nil
}

// sqlTulisNomorPL merakit penulisan nomor ke baris peserta.
//
// ⛔ HANYA baris yang BELUM bernomor yang disentuh. `WHERE ... PL_NUMBER IS
// NULL` membuat penulisan ini tidak dapat menimpa nomor yang sudah terbit,
// bahkan bila pemanggilnya keliru - dan gerbang "lahir sekali" karena itu
// berdiri di DUA tempat: di layanan, dan di kalimat `WHERE` ini.
func sqlTulisNomorPL(detail string) string {
	return fmt.Sprintf(`UPDATE %s SET PL_NUMBER = :1
	  WHERE PREMIUM_LIST_ID = :2 AND PL_NUMBER IS NULL`, detail)
}

// TulisNomor menuliskan nomor ke seluruh baris peserta yang belum bernomor.
//
// Mengembalikan cacah baris yang tersentuh.
func (r *NomorPolis) TulisNomor(ctx context.Context, tx *db.Tx,
	polisID, nomor string) (int, error) {

	if strings.TrimSpace(nomor) == "" {
		return 0, errors.New("repository: menolak menulis PL_NUMBER kosong")
	}
	detail, err := r.db.Qualify("T_PREMIUM_LIST_DETAIL")
	if err != nil {
		return 0, err
	}
	q := sqlTulisNomorPL(detail)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	hasil, err := tx.ExecContext(ctx, q, nomor, polisID)
	if err != nil {
		return 0, fmt.Errorf("repository: menulis PL_NUMBER: %w", err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("repository: membaca cacah baris ter-update: %w", err)
	}
	// ⛔ Nol baris tersentuh sesudah gerbang lolos berarti keadaannya BERUBAH
	// di antara pembacaan dan penulisan. Nomor sudah diambil dari penghitung
	// pada titik ini; membiarkannya lewat berarti menaikkan penghitung tanpa
	// menyimpan nomornya, dan nomor yang hilang itu tidak akan pernah terpakai.
	//
	// ⛔ SEBABNYA DIBEDAKAN, bukan ditebak. Ada DUA cara sampai ke sini, dan
	// keduanya menuntut jawaban yang berbeda:
	//
	//	nol baris peserta   -> polis memang belum siap dinomori
	//	ada baris, semuanya sudah bernomor -> permintaan LAIN mendahului kita
	//
	// Menjawab keduanya dengan "belum punya peserta" mengirim orang mencari
	// peserta yang sebenarnya ada - dan menyembunyikan satu-satunya petunjuk
	// bahwa dua permintaan berjalan bersamaan atas polis yang sama.
	if n == 0 {
		cacah, errCacah := r.cacahPeserta(ctx, tx, detail, polisID)
		if errCacah != nil {
			return 0, errCacah
		}
		if cacah == 0 {
			return 0, ErrPolisTanpaPeserta
		}
		return 0, fmt.Errorf("%w: polis %q", ErrNomorPLTerbitBersamaan, polisID)
	}
	return int(n), nil
}

// cacahPeserta mencacah baris peserta polis di dalam transaksi.
//
// ⛔ Query-nya DIPINJAM dari `sqlCacahPeserta` (polis_detail.go, paket yang
// sama), tidak diketik ulang. Ronde pertama menyalinnya utuh - dan dua
// salinan satu query berarti penyaring yang suatu hari hanya diperbaiki di
// satu tempat, yaitu dua jawaban berbeda untuk pertanyaan yang sama.
func (r *NomorPolis) cacahPeserta(ctx context.Context, tx *db.Tx,
	detail, polisID string) (int, error) {

	q := sqlCacahPeserta(detail)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var cacah int
	if err := tx.QueryRowContext(ctx, q, polisID).Scan(&cacah); err != nil {
		return 0, fmt.Errorf("repository: mencacah peserta polis: %w", err)
	}
	return cacah, nil
}
