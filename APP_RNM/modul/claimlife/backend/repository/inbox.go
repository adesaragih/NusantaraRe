package repository

// Kotak masuk klaim - F0.4.
//
// Untuk apa berkas ini: membaca antrian kerja per tahap, berbatas dan
// berurut, beserta cacah totalnya untuk lencana tab.
//
// ⛔ BENTUKNYA DARI XML, bukan dari selera:
//
//	`[terverifikasi]` `Claim Life/ReportDefinition/InboxPremiumList.xml`
//	  592  `<pyPageSize>50</pyPageSize>`        - 50 baris per halaman
//	  942  `<pyMaxRecords>500</pyMaxRecords>`   - maksimum 500
//	  733  `<pySortType>DESC</pySortType>`      - menurun
//	  736  `<pyFieldLabel>Create Date/Time`     - pada waktu BUAT
//	  721  `Case ID` · 751 `Create Operator Name` · 765 `Work Status`
//
// ⚠️ Diurutkan `TGL_CREATE`, BUKAN `TGL_UPDATE` (butir au). `TGL_UPDATE`
// ditimpa tiap perpindahan, sehingga daftar yang diurutkan dengannya
// melompat-lompat setiap kali seseorang memindah kasus lain.
//
// Dibaca sesudah: klaimlife.go.

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
)

// Batas halaman - dari `InboxPremiumList.xml`.
const (
	// UkuranHalamanBawaan - `pyPageSize` baris 592.
	UkuranHalamanBawaan = 50
	// UkuranHalamanMaksimum - `pyMaxRecords` baris 942.
	//
	// ⛔ Plafon, bukan saran. Permintaan yang meminta lebih DIPANGKAS, tidak
	// ditolak: pemanggil yang meminta 10.000 baris hampir selalu salah tulis,
	// dan menolaknya membuat layar kosong alih-alih menampilkan 500 pertama.
	UkuranHalamanMaksimum = 500
)

// BarisInbox adalah satu kasus di antrian.
//
// ⚠️ Nol medan uang. Nilai klaim tidak ikut di daftar (ADR-U-0003); layar
// yang memerlukannya membuka kasusnya.
type BarisInbox struct {
	// CaseID mengisi kolom `Case ID` - `InboxPremiumList.xml:721`.
	CaseID string
	// WorkID adalah pengenal baris work; dipakai membuka kasusnya.
	WorkID string
	// Tahap adalah nama assignment VERBATIM (butir at).
	Tahap string
	// CreateOpName mengisi `Create Operator Name` - baris 751.
	//
	// ⛔ DATA ORANG. Ia ditampilkan di layar orang yang berhak melihat
	// antriannya, dan tidak disalin ke mana pun lagi.
	CreateOpName string
	// TglCreate mengisi `Create Date/Time` - baris 736.
	TglCreate time.Time
	// StatusKlaim mengisi `Work Status` - baris 765.
	StatusKlaim string
	NomorKlaim  string
	NomorPolis  string
	NamaBisnis  string
	MataUang    string
}

// SaringInbox adalah penyaring satu pembacaan antrian.
type SaringInbox struct {
	// Tahap - nama assignment VERBATIM; wajib.
	Tahap string
	// PeranTahap - nilai `PY_POSITION` yang setara, untuk baris LAMA yang
	// kolom `TAHAP`-nya masih kosong.
	PeranTahap string
	// AkunID menyaring worklist ke kasus milik pelaku sendiri. Kosong berarti
	// workbasket - seluruh kasus pada tahap itu.
	AkunID string
	Offset int
	Ukuran int
}

// sqlInboxWhere merakit klausa penyaring yang DIPAKAI BERSAMA oleh pembacaan
// baris dan pencacahan total.
//
// ⛔ Satu tempat, bukan dua. Dua salinan klausa penyaring berarti lencana
// dapat menyebut angka yang tidak cocok dengan isi daftarnya - dan yang salah
// justru yang lebih dipercaya pemakai, sebab ia lebih ringkas.
//
// ⚠️ `NVL(w.TAHAP, :2)` membuat baris LAMA (TAHAP kosong) ikut terbaca lewat
// `PY_POSITION`-nya. Tanpa itu seluruh kasus yang sudah ada sebelum migrasi
// 016 menghilang dari kotak masuk tanpa satu pun galat.
func sqlInboxWhere(pakaiAkun bool) string {
	w := ` WHERE NVL(w.TAHAP, :tahapCadangan) = :tahap`
	if pakaiAkun {
		w += ` AND w.CREATE_OP = :akun`
	}
	return w
}

// sqlInbox merakit pembacaan barisnya.
func sqlInbox(work, header string, pakaiAkun bool) string {
	return fmt.Sprintf(`SELECT w.CASE_ID, w.ID, NVL(w.TAHAP, :tahapCadangan),
			   w.CREATE_OP_NAME, w.TGL_CREATE,
			   h.STS_REJECT, h.CLAIM_NO, h.POLICY_NO, h.BUSINESS_NAME, h.CURRENCY
		  FROM %s w
		  JOIN %s h ON h.ID = w.ID%s
		 ORDER BY w.TGL_CREATE DESC, w.ID DESC
		 OFFSET :offset ROWS FETCH NEXT :ukuran ROWS ONLY`,
		work, header, sqlInboxWhere(pakaiAkun))
}

// sqlCacahInbox merakit pencacahannya.
func sqlCacahInbox(work, header string, pakaiAkun bool) string {
	return fmt.Sprintf(`SELECT COUNT(*)
		  FROM %s w
		  JOIN %s h ON h.ID = w.ID%s`,
		work, header, sqlInboxWhere(pakaiAkun))
}

// AmbilInbox membaca satu halaman antrian beserta cacah totalnya.
//
// ⚠️ Cacah total dibaca TERPISAH dan disengaja: lencana tab menyebut seluruh
// kasus pada tahap itu, sedangkan daftarnya hanya satu halaman. Menurunkan
// totalnya dari panjang halaman akan membuat lencana berkata "50" selamanya.
func (r *KlaimLife) AmbilInbox(ctx context.Context, s SaringInbox) (
	[]BarisInbox, int, error) {

	work, err := r.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return nil, 0, err
	}
	header, err := r.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return nil, 0, err
	}
	pakaiAkun := s.AkunID != ""

	// ⛔ Urutan argumen mengikuti urutan MUNCULNYA penanda di teks, sebab
	// driver Oracle mengikat penanda bernama secara berurutan di jalur ini.
	// Cadangan tahap muncul dua kali di pembacaan baris (SELECT dan WHERE).
	argBaris := []any{s.Tahap, s.PeranTahap, s.Tahap}
	argCacah := []any{s.PeranTahap, s.Tahap}
	if pakaiAkun {
		argBaris = append(argBaris, s.AkunID)
		argCacah = append(argCacah, s.AkunID)
	}

	qCacah := sqlCacahInbox(work, header, pakaiAkun)
	if err := db.PeriksaSQL(qCacah); err != nil {
		return nil, 0, err
	}
	var total int
	if err := r.db.QueryRowContext(ctx, qCacah, argCacah...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository: mencacah kotak masuk: %w", err)
	}

	q := sqlInbox(work, header, pakaiAkun)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, 0, err
	}
	argBaris = append(argBaris, s.Offset, s.Ukuran)
	baris, err := r.db.QueryContext(ctx, q, argBaris...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: membaca kotak masuk: %w", err)
	}
	defer func() { _ = baris.Close() }()

	var hasil []BarisInbox
	for baris.Next() {
		var (
			b                                       BarisInbox
			caseID, workID, tahap, opName           sql.NullString
			status, nomorKlaim, polis, bisnis, kurs sql.NullString
			tglCreate                               sql.NullTime
		)
		if err := baris.Scan(&caseID, &workID, &tahap, &opName, &tglCreate,
			&status, &nomorKlaim, &polis, &bisnis, &kurs); err != nil {
			return nil, 0, fmt.Errorf("repository: membaca baris kotak masuk: %w", err)
		}
		b.CaseID, b.WorkID, b.Tahap = caseID.String, workID.String, tahap.String
		b.CreateOpName, b.TglCreate = opName.String, tglCreate.Time
		b.StatusKlaim, b.NomorKlaim = status.String, nomorKlaim.String
		b.NomorPolis, b.NamaBisnis, b.MataUang = polis.String, bisnis.String, kurs.String
		hasil = append(hasil, b)
	}
	if err := baris.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository: menelusuri kotak masuk: %w", err)
	}
	return hasil, total, nil
}
