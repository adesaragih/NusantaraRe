package repository

import (
	"context"
	"database/sql"
	"fmt"

	"time"

	"nusantarare/internal/models"
	"nusantarare/pkg/utils"
)

// KlaimLife membaca klaim Life beserta peserta dan baris adjustment-nya.
//
// Satu antarmuka per agregat. Paket ini memuat SELURUH SQL dan tidak memuat
// satu pun aturan dagang - perakitan agregat milik services.
type KlaimLife struct {
	db *DB
}

// NewKlaimLife membuat pembaca klaim Life.
func NewKlaimLife(db *DB) *KlaimLife { return &KlaimLife{db: db} }

// fmtDesimal memaksa Oracle menyerahkan angka sebagai TEKS, bukan sebagai
// bilangan pecahan biner.
//
// ⛔ Uang tidak pernah melewati float (ADR-U-0003, ADR-U-0016). Membiarkan
// driver menyerahkan NUMBER sebagai float64 melanggar aturan itu di tempat yang
// paling sulit terlihat - karena itu konversinya dilakukan di dalam SQL.
//
// 'TM9' memberi bentuk desimal terpendek tanpa notasi ilmiah. Argumen NLS
// memaksa titik sebagai pemisah desimal: tanpa itu, sesi ber-NLS Indonesia
// mengembalikan koma dan pembacaannya gagal senyap.
const fmtDesimal = `TO_CHAR(%s, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')`

// AmbilHeader membaca satu baris header klaim. Mengembalikan nil bila tidak ada.
func (r *KlaimLife) AmbilHeader(ctx context.Context, id string) (*models.Klaim, error) {
	tabel, err := r.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return nil, err
	}
	// CLAIM_RETRO dibaca lewat TO_CHAR ber-argumen NLS: ia uang (butir w2), dan
	// uang tidak pernah lewat float maupun bergantung setelan sesi
	// (ADR-U-0003, ADR-U-0016). CURRENCY menyertainya sejak butir z1.
	q := fmt.Sprintf(
		`SELECT ID, CLAIM_NO, POLICY_NO, BUSINESS_NAME, STS_REJECT, CURRENCY, `+fmtDesimal+
			` FROM %s WHERE ID = :1`, "CLAIM_RETRO", tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}

	var (
		kID                              string
		nomorKlaim, nomorPolis, bisnis   sql.NullString
		kodeStatus, mataUang, claimRetro sql.NullString
	)
	err = r.db.sql.QueryRowContext(ctx, q, id).
		Scan(&kID, &nomorKlaim, &nomorPolis, &bisnis, &kodeStatus, &mataUang, &claimRetro)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository: membaca header klaim: %w", err)
	}
	// ✅ Mata uangnya datang dari kolom CURRENCY header sejak butir z1
	// (26-09-2026). Sampai kolom itu ada, nilai ini bermata uang kosong -
	// benar secara mekanis, tidak berguna bagi pembacanya.
	retro, err := uraiUang(kID, "CLAIM_RETRO", claimRetro, mataUang.String)
	if err != nil {
		return nil, err
	}
	return &models.Klaim{
		ID:         kID,
		NomorKlaim: nomorKlaim.String,
		NomorPolis: nomorPolis.String,
		NamaBisnis: bisnis.String,
		KodeStatus: kodeStatus.String,
		ClaimRetro: retro,
	}, nil
}

// AmbilPeserta membaca seluruh peserta milik satu klaim, terurut stabil.
func (r *KlaimLife) AmbilPeserta(ctx context.Context, klaimID string) ([]models.Peserta, error) {
	tabel, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return nil, err
	}
	// Daftar kolomnya datang dari kolompeserta.go, sama dengan yang dipakai
	// saat menulis. Angka dan tanggal dibungkus TO_CHAR di sana.
	q := fmt.Sprintf(`SELECT ID, %s FROM %s WHERE CLAIM_ID = :1 ORDER BY ID`,
		selectPeserta(), tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.sql.QueryContext(ctx, q, klaimID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca peserta: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := []models.Peserta{}
	for rows.Next() {
		var id string
		sel := make([]sql.NullString, len(kolomPeserta))
		tujuan := make([]any, 0, len(sel)+1)
		tujuan = append(tujuan, &id)
		for i := range sel {
			tujuan = append(tujuan, &sel[i])
		}
		if err := rows.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: membaca baris peserta: %w", err)
		}
		ps, err := rakitPeserta(id, sel)
		if err != nil {
			return nil, err
		}
		out = append(out, ps)
	}
	return out, rows.Err()
}

// AmbilBaris membaca SELURUH baris adjustment milik satu klaim, dikelompokkan
// per peserta.
//
// Tiket 01 AC-1: seluruh baris, bukan hanya yang terakhir. Sistem lama terbiasa
// menyentuh `AdjustmentList(<LAST>)` saja; pembacaan di sini tidak pernah
// membatasi jumlah baris.
func (r *KlaimLife) AmbilBaris(ctx context.Context, klaimID string) (map[string][]models.BarisAdjustment, error) {
	adj, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT")
	if err != nil {
		return nil, err
	}
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return nil, err
	}
	// Kedelapan kolom warisan ikut dibaca sejak 26-09-2026: menulisnya tanpa
	// membacanya membuat separuh pewarisan mati tanpa satu pun test gagal -
	// cacat yang persis sama sudah terjadi pada CLAIM_RETRO di tiket 14.
	q := fmt.Sprintf(
		`SELECT a.PREMIUM_LIST_DETAIL_ID, a.ID, `+fmtDesimal+`, a.CURRENCY,
		        a.STS_REJECT, a.ACCEPTED_NO, a.ACCEPTATION_DATE, a.KOMITE_ID,
		        `+fmtDesimal+`, `+fmtDesimal+`, `+fmtDesimal+`, `+fmtDesimal+`,
		        `+fmtDesimal+`, `+fmtDesimal+`, a.CURRENCY_ID
		   FROM %s a JOIN %s p ON p.ID = a.PREMIUM_LIST_DETAIL_ID
		  WHERE p.CLAIM_ID = :1
		  ORDER BY a.PREMIUM_LIST_DETAIL_ID, a.ID`,
		"a.CLAIM_AMOUNT",
		"a.SHARE_NUSANTARA_RE", "a.CEDING_RETENTION", "a.SUM_REASURED",
		"a.SUM_INSURED", "a.SHARE_RETRO", "a.RETROCEDED_SHARE",
		adj, pes)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.sql.QueryContext(ctx, q, klaimID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca baris adjustment: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := map[string][]models.BarisAdjustment{}
	for rows.Next() {
		var pesertaID, id string
		var jumlah, mataUang, kodeStatus, nomorAksep, komiteID sql.NullString
		var tglAksep sql.NullTime
		var warisan [6]sql.NullString
		var mataUangID sql.NullString
		if err := rows.Scan(&pesertaID, &id, &jumlah, &mataUang,
			&kodeStatus, &nomorAksep, &tglAksep, &komiteID,
			&warisan[0], &warisan[1], &warisan[2], &warisan[3],
			&warisan[4], &warisan[5], &mataUangID); err != nil {
			return nil, fmt.Errorf("repository: membaca satu baris adjustment: %w", err)
		}

		b := models.BarisAdjustment{
			ID:             id,
			KodeStatus:     kodeStatus.String,
			NomorAkseptasi: nomorAksep.String,
			KomiteID:       komiteID.String,
			// Mata uang dibawa walau jumlahnya kosong: keduanya kolom
			// terpisah, dan membuang mata uang hanya karena jumlahnya NULL
			// menghilangkan fakta yang tersimpan.
			JumlahKlaim: models.Money{Currency: mataUang.String},
			CurrencyID:  mataUangID.String,
		}
		// Urutan tujuan mengikuti urutan kolom di SELECT di atas, persis.
		tujuanWarisan := []*models.Money{
			&b.ShareNusantaraRe, &b.CedingRetention, &b.SumReasured,
			&b.SumInsured, &b.ShareRetro, &b.RetrocededShare,
		}
		namaWarisan := []string{
			"SHARE_NUSANTARA_RE", "CEDING_RETENTION", "SUM_REASURED",
			"SUM_INSURED", "SHARE_RETRO", "RETROCEDED_SHARE",
		}
		for i := range tujuanWarisan {
			m, err := uraiUang(id, namaWarisan[i], warisan[i], mataUang.String)
			if err != nil {
				return nil, err
			}
			*tujuanWarisan[i] = m
		}
		// Kolom kosong tetap kosong; ia tidak menjadi nol (ADR-U-0027).
		if jumlah.Valid && jumlah.String != "" {
			d, err := utils.ParseDecimal(jumlah.String)
			if err != nil {
				return nil, fmt.Errorf("repository: baris %s: %w", id, err)
			}
			b.JumlahKlaim.Amount = d
		}
		if tglAksep.Valid {
			b.TanggalAkseptasi = tglAksep.Time
		}
		out[pesertaID] = append(out[pesertaID], b)
	}
	return out, rows.Err()
}

// TypeKlaim membaca `Type` klaim dari baris work object-nya.
//
// ⛔ Satu-satunya pembaca `Type`, dan ia membacanya dari satu-satunya tempat
// yang menyimpannya - `T_WORK_CLAIM.TYPE`. Tiket 06 AC 5: validasi dan
// penurunan jenis klaim membaca satu field yang sama, bukan dua salinan
// seperti di Pega.
//
// Header klaim dan baris work object berbagi kunci utama (`isiIdentitas`),
// jadi pengenal klaim dapat dipakai apa adanya.
func (r *KlaimLife) TypeKlaim(ctx context.Context, klaimID string) (string, error) {
	tabel, err := r.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf(`SELECT TYPE FROM %s WHERE ID = :1`, tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	var tipe sql.NullString
	err = r.db.sql.QueryRowContext(ctx, q, klaimID).Scan(&tipe)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("repository: work object %q tidak ada", klaimID)
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca Type klaim: %w", err)
	}
	return tipe.String, nil
}

// PerbaruiTanggalKejadian menulis DATE_OF_LOSS seorang peserta.
//
// ⛔ Tanggalnya ditulis lewat TO_DATE berformat tetap, sama dengan jalur tulis
// kolom tanggal peserta yang lain - bentuknya tidak pernah bergantung
// NLS_DATE_FORMAT sesi (ADR-U-0022).
func (r *KlaimLife) PerbaruiTanggalKejadian(ctx context.Context, tx *Tx,
	pesertaID string, dol time.Time) error {
	tabel, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(
		`UPDATE %s SET DATE_OF_LOSS = TO_DATE(:1, 'YYYY-MM-DD HH24:MI:SS') WHERE ID = :2`,
		tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, utils.FormatTanggalWaktu(dol), pesertaID)
	if err != nil {
		return fmt.Errorf("repository: menulis DATE_OF_LOSS: %w", err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: mencacah baris tersentuh: %w", err)
	}
	if n != 1 {
		return fmt.Errorf("repository: DATE_OF_LOSS peserta %q menyentuh %d baris, mau 1",
			pesertaID, n)
	}
	return nil
}

// PerbaruiStatusBaris menulis status ke baris adjustment DAN ke pesertanya.
//
// Header klaim TIDAK ditulis di sini - sumbernya baris yang berbeda; lihat
// CerminkanHeader. Pemanggil menjalankan keduanya dalam satu transaksi.
//
// `[terverifikasi]` sasarannya diturunkan dari sensus penulis `STS_REJECT` di
// seluruh korpus Claim Life - ENAM Property-Set di LIMA rule:
//
//	baris adjustment : SaveOutStandingLife_Act  `.STS_REJECT = 0`
//	                   RejectOSClaimLife_Act    `.STS_REJECT = 2`
//	                   SaveAdjustment_Act       `.AdjustmentList(<LAST>)` = 1,
//	                                            beserta ACCEPTEDNO dan
//	                                            ACCEPTATION_DATE
//	peserta          : RejectOSClaimLife_Act
//	                   `PremiumListDetail(idx).STS_REJECT = 2`
//	header klaim     : serviceInsertArasapasClaimLife_act
//	                   `pyWorkPage.ClaimData.STS_REJECT = .STS_REJECT`
//	diagnosa         : SetSTS_Reject atas `.DiagnoseList` - TIDAK ditiru,
//	                   tabelnya belum ada (butir al masih `[USULAN]`)
//
// ⛔ SYARAT KODE LAMA ikut di WHERE. Baris dibaca di luar transaksi, jadi ia
// dapat berubah di antara baca dan tulis. Tanpa syarat itu, kefinalan hanya
// berlaku di dalam satu proses dan dua permintaan serentak dapat sama-sama
// menang - yang kedua menimpa keputusan yang pertama tanpa jejak.
// ⚠️ Kedua tingkat TIDAK mencerminkan baris yang sama, dan itu bacaan XML
// bukan penyederhanaan kami: peserta mengikuti baris yang BERUBAH
// (`RejectOSClaimLife_Act` memakai `local.IndexPremium`, yaitu peserta pemilik
// baris itu), sedangkan header mengikuti baris TERAKHIR yang diulang
// (`serviceInsertArasapasClaimLife_act`, putaran bersarang tanpa henti).
// Karena itu header punya methodnya sendiri, `CerminkanHeader`.
func (r *KlaimLife) PerbaruiStatusBaris(ctx context.Context, tx *Tx,
	pesertaID, adjID, kodeLama, kode, nomorAksep string, tglAksep time.Time) error {

	adj, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT")
	if err != nil {
		return err
	}
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return err
	}

	// Peserta tidak punya kolom ACCEPTED_NO - hanya baris dan header punya.
	langkah := []struct {
		nama string
		q    string
		args []any
	}{
		// ACCEPTATION_DATE ditulis bersama status, seperti `SaveAdjustment_Act`
		// yang menyetel ketiganya dalam satu Property-Set.
		{"baris adjustment", fmt.Sprintf(
			`UPDATE %s SET STS_REJECT = :1, ACCEPTED_NO = :2, ACCEPTATION_DATE = :3
			  WHERE ID = :4 AND STS_REJECT = :5`, adj),
			[]any{kosongJadiNil(kode), kosongJadiNil(nomorAksep),
				waktuJadiNil(tglAksep), adjID, kosongJadiNil(kodeLama)}},
		{"peserta", fmt.Sprintf(
			`UPDATE %s SET STS_REJECT = :1 WHERE ID = :2`, pes),
			[]any{kosongJadiNil(kode), pesertaID}},
	}
	for _, l := range langkah {
		if err := PeriksaSQL(l.q); err != nil {
			return err
		}
		hasil, err := tx.tx.ExecContext(ctx, l.q, l.args...)
		if err != nil {
			return fmt.Errorf("repository: mencerminkan status ke %s: %w", l.nama, err)
		}
		n, err := hasil.RowsAffected()
		if err != nil {
			return fmt.Errorf("repository: mencacah baris %s: %w", l.nama, err)
		}
		if n != 1 {
			return fmt.Errorf("repository: pencerminan ke %s menyentuh %d baris, mau 1 "+
				"(baris yang statusnya sudah berubah sejak dibaca tidak ditimpa)",
				l.nama, n)
		}
	}
	return nil
}

// CerminkanHeader menyalin status BARIS TERAKHIR ke header klaim.
//
// Dipisah dari PerbaruiStatusBaris karena sumbernya baris yang BERBEDA - lihat
// komentar di atas. Pemanggilnya wajib menjalankan keduanya dalam satu
// transaksi yang sama.
func (r *KlaimLife) CerminkanHeader(ctx context.Context, tx *Tx,
	klaimID, kode, nomorAksep string) error {
	hdr, err := r.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(
		`UPDATE %s SET STS_REJECT = :1, ACCEPTED_NO = :2 WHERE ID = :3`, hdr)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q,
		kosongJadiNil(kode), kosongJadiNil(nomorAksep), klaimID)
	if err != nil {
		return fmt.Errorf("repository: mencerminkan status ke header klaim: %w", err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: mencacah baris header: %w", err)
	}
	if n != 1 {
		return fmt.Errorf("repository: pencerminan ke header menyentuh %d baris, mau 1", n)
	}
	return nil
}

// PasangPenandaDipilih menyetel IS_CHECK peserta menjadi dipilih-untuk-diklaim.
//
// ⛔ Kebalikan `CabutPenandaDipilih`, dan ia WAJIB dipanggil ketika putaran
// adjustment berikutnya lahir. `[terverifikasi]` `Komite Claim Life/Activity/
// KomitePostAdjustment.xml` pecahan baris 1900, 2117, 4896, 5341, 7622, dan
// 8067: keenam prasyaratnya menuntut `.IsCheck = true`. Peserta yang
//
// ⚠️ SENSUS, jendelanya DINAMAI. Di berkas pecahan `KomitePostAdjustment.xml`
// kata `IsCheck` muncul di **13** baris; **4** di antaranya tag
// `<pyStepsPreCondParamsWhen>` langsung (4896, 5341, 7622, 8067), dan **2**
// lagi tersimpan sebagai `rowdata REPEATINGINDEX="pyStepsPreCondParamsWhen"`
// (1900, 2117) - masing-masing berpasangan dengan kembaran `<pyExpression>`
// (1897, 2114) yang berisi kondisi yang sama.
//
// Jendela yang dipakai: **LANGKAH PRASYARAT YANG BERBEDA** → **6**. Bila
// yang dihitung ELEMEN XML pembawa kondisi, angkanya **8**; bila seluruh
// penyebutan, **13**. Ketiganya benar untuk pertanyaan yang berbeda, dan
// angka 6 tidak berarti apa-apa tanpa kalimat ini.
//
// Audit: `py` + pemecah `><` → `>\n<`, lalu cacah baris ber-`IsCheck` yang
// juga ber-`pyStepsPreCondParamsWhen` (4) dan yang ber-`rowdata` (2).
// penandanya masih tercabut sejak penolakan sebelumnya tidak akan pernah
// diambil Komite - putarannya lahir, lalu mati diam-diam.
func (r *KlaimLife) PasangPenandaDipilih(ctx context.Context, tx *Tx,
	pesertaID string) error {
	return r.setelPenandaDipilih(ctx, tx, pesertaID, "true")
}

// CabutPenandaDipilih menyetel IS_CHECK peserta menjadi tidak-dipilih.
//
// `[terverifikasi]` `RejectOSClaimLife_Act` langkah 2 (pecahan baris 517-518)
// menyetel `PremiumListDetail(local.IndexPremium).IsCheck = "false"` bersama
// kedua penulisan STS_REJECT-nya. Peserta yang barisnya dibatalkan berhenti
// terhitung "dipilih untuk diklaim", sehingga ia dapat dipilih ulang dengan
// baris pengganti.
//
// Nilainya ditulis "false" persis seperti rule-nya - teks, bukan bilangan
// (ADR-U-0022), dan bukan NULL: tidak-dipilih adalah pernyataan, sedangkan
// NULL berarti belum pernah diputuskan.
func (r *KlaimLife) CabutPenandaDipilih(ctx context.Context, tx *Tx, pesertaID string) error {
	return r.setelPenandaDipilih(ctx, tx, pesertaID, "false")
}

// setelPenandaDipilih menulis IS_CHECK peserta.
//
// ⛔ Nilainya TEKS - "true"/"false" apa adanya seperti kolomnya, bukan boolean
// Go yang diformat (ADR-U-0022).
func (r *KlaimLife) setelPenandaDipilih(ctx context.Context, tx *Tx,
	pesertaID, nilai string) error {

	tabel, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`UPDATE %s SET IS_CHECK = :1 WHERE ID = :2`, tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, nilai, pesertaID)
	if err != nil {
		return fmt.Errorf("repository: menyetel IS_CHECK menjadi %q: %w", nilai, err)
	}
	return pastikanSatuBaris(hasil, "penyetelan IS_CHECK")
}

// CaseIDKlaim membaca CASE_ID baris work object sebuah klaim.
//
// ⛔ DIBACA, tidak diandaikan sama dengan pengenal klaim. Butir ae1 memang
// memutuskan `CASEID = pengenal work object` untuk klaim yang sistem ini
// buat sendiri, tetapi itu keputusan pengisian - bukan jaminan bentuk. Klaim
// yang kelak dimigrasikan (tiket 13) membawa CASE_ID warisannya sendiri, dan
// kode yang mengandaikan keduanya sama akan mencacah baris milik klaim lain.
func (r *KlaimLife) CaseIDKlaim(ctx context.Context, klaimID string) (string, error) {
	tabel, err := r.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf(`SELECT CASE_ID FROM %s WHERE ID = :1`, tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	var caseID sql.NullString
	err = r.db.sql.QueryRowContext(ctx, q, klaimID).Scan(&caseID)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("repository: work object %q tidak ada", klaimID)
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca CASE_ID: %w", err)
	}
	return caseID.String, nil
}

// PerbaruiTahap memindahkan kasus di tangga kerja.
//
// ⛔ SATU pernyataan, dan ia TIDAK menyentuh STS_REJECT mana pun. Perpindahan
// tahap memindahkan pekerjaan, bukan memutuskan klaim (ADR-U-0011).
//
// Kolomnya sudah ada di migrasi 001: PY_POSITION, SENDTO_ADMIN, SENDTO_MEDICAL.
// Nol kolom baru, nol langkah migrasi.
//
// ⛔ Tahap ASAL ikut sebagai syarat WHERE. Ia dibaca di luar transaksi, jadi
// dapat berubah di antara baca dan tulis - dan dua pemindahan serentak dapat
// sama-sama menang, yang kedua menimpa yang pertama tanpa jejak. Penjaga yang
// sama sudah dipasang pada perubahan status (tiket 04).
func (r *KlaimLife) PerbaruiTahap(ctx context.Context, tx *Tx,
	klaimID, tahapAsal, tahapTujuan, peranTujuan, sendtoAdmin, sendtoMedical string,
	saat time.Time) error {
	tabel, err := r.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	// ⛔ BUTIR at: `TAHAP` ikut ditulis, dan penjaga optimisnya kini
	// memakai TAHAP - bukan PY_POSITION. Alasannya memaksa: pada
	// perpindahan Input Register ⇄ Outstanding Claim `PY_POSITION` TIDAK
	// BERUBAH (keduanya ReasLifeAdmin), sehingga penjaga ber-PY_POSITION
	// tidak dapat mendeteksi kasus yang sudah dipindahkan orang lain
	// sejak dibaca.
	//
	// ⚠️ `NVL` dipakai untuk baris LAMA yang `TAHAP`-nya masih kosong:
	// tanpa itu setiap perpindahan pertama baris lama akan menyentuh nol
	// baris dan gagal, padahal tidak ada yang salah dengannya.
	q := fmt.Sprintf(`UPDATE %s SET PY_POSITION = :1, SENDTO_ADMIN = :2,
		 SENDTO_MEDICAL = :3, TGL_UPDATE = :4, TAHAP = :5
		 WHERE ID = :6 AND NVL(TAHAP, :7) = :8`, tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, kosongJadiNil(peranTujuan),
		kosongJadiNil(sendtoAdmin), kosongJadiNil(sendtoMedical),
		waktuJadiNil(saat), kosongJadiNil(tahapTujuan),
		klaimID, kosongJadiNil(tahapAsal), kosongJadiNil(tahapAsal))
	if err != nil {
		return fmt.Errorf("repository: memperbarui tahap: %w", err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: mencacah baris work: %w", err)
	}
	if n != 1 {
		return fmt.Errorf("repository: perpindahan tahap menyentuh %d baris, mau 1 "+
			"(kasus yang tahapnya sudah berubah sejak dibaca tidak ditimpa)", n)
	}
	return nil
}

// TahapKlaim membaca PY_POSITION baris work object sebuah klaim.
func (r *KlaimLife) TahapKlaim(ctx context.Context, klaimID string) (string, error) {
	_, peran, err := r.TahapDanPeran(ctx, klaimID)
	return peran, err
}

// TahapDanPeran membaca TAHAP dan PY_POSITION sebuah kasus sekaligus.
//
// ⛔ BUTIR at. Keduanya dibaca BERSAMA, dalam satu kueri, karena keduanya
// menjawab pertanyaan yang berbeda dan pemanggil memerlukan keduanya:
// `TAHAP` mengatakan di anak tangga mana kasusnya berdiri, `PY_POSITION`
// mengatakan peran siapa yang memegangnya. Dua kueri berarti ada jendela
// ketika keduanya dibaca dari keadaan yang berbeda.
//
// ⚠️ `TAHAP` KOSONG dikembalikan kosong, bukan ditebak di sini. Yang
// menebaknya `services` lewat `models.TahapDariPeran`, dan ia menebak
// dengan sadar: baris lama yang sebenarnya di Input Register akan tampak
// Outstanding. Menebak di repository menyembunyikan bahwa itu tebakan.
func (r *KlaimLife) TahapDanPeran(ctx context.Context, klaimID string) (
	tahap, peran string, err error) {

	tabel, err := r.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return "", "", err
	}
	q := fmt.Sprintf(`SELECT TAHAP, PY_POSITION FROM %s WHERE ID = :1`, tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", "", err
	}
	var kolomTahap, posisi sql.NullString
	err = r.db.sql.QueryRowContext(ctx, q, klaimID).Scan(&kolomTahap, &posisi)
	if err == sql.ErrNoRows {
		return "", "", fmt.Errorf("repository: work object %q tidak ada", klaimID)
	}
	if err != nil {
		return "", "", fmt.Errorf("repository: membaca tahap kasus: %w", err)
	}
	return kolomTahap.String, posisi.String, nil
}

// PerbaruiKomiteID menautkan baris adjustment ke kasus Komite yang baru lahir.
//
// Tiket 10 AC 61. Rujukan, BUKAN salinan: roster dan keputusan per anggota
// tetap milik konteks Komite Claim Life.
//
// KOMITE_ID lama ikut di WHERE dan wajib NULL. Baris yang sudah tertaut sejak
// dibaca tidak ditimpa - tanpa klausa itu dua penyerahan bersamaan menghasilkan
// dua kasus komite dan hanya satu yang teringat.
func (r *KlaimLife) PerbaruiKomiteID(ctx context.Context, tx *Tx,
	adjID, komiteID string) error {

	tabel, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(
		`UPDATE %s SET KOMITE_ID = :1 WHERE ID = :2 AND KOMITE_ID IS NULL`, tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, kosongJadiNil(komiteID), adjID)
	if err != nil {
		return fmt.Errorf("repository: menautkan baris ke Komite: %w", err)
	}
	return pastikanSatuBaris(hasil, "penautan baris ke Komite")
}

// pastikanSatuBaris menuntut sebuah pernyataan menyentuh tepat satu baris.
func pastikanSatuBaris(hasil sql.Result, nama string) error {
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: mencacah baris %s: %w", nama, err)
	}
	if n != 1 {
		return fmt.Errorf("repository: %s menyentuh %d baris, mau 1", nama, n)
	}
	return nil
}
