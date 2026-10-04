// Package repository adalah SATU-SATUNYA lapisan yang menyentuh Oracle.
//
// Arah ketergantungan: handlers -> services -> repository. Tidak terbalik,
// tidak memotong. Paket ini tidak pernah mengimpor handlers atau services.
// Pintu koneksi dan transaksinya tinggal di `inti/backend/db` (refactor bentuk B).
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/claimlife/backend/models"
)

// KlaimLife membaca klaim Life beserta peserta dan baris adjustment-nya.
//
// Satu antarmuka per agregat. Paket ini memuat SELURUH SQL dan tidak memuat
// satu pun aturan dagang - perakitan agregat milik services.
type KlaimLife struct {
	db *db.DB
}

// NewKlaimLife membuat pembaca klaim Life.
func NewKlaimLife(db *db.DB) *KlaimLife { return &KlaimLife{db: db} }

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
		`SELECT ID, CLAIM_NO, POLICY_NO, BUSINESS_NAME, STS_REJECT, CURRENCY, `+db.FmtDesimal+
			` FROM %s WHERE ID = :1`, "CLAIM_RETRO", tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}

	var (
		kID                              string
		nomorKlaim, nomorPolis, bisnis   sql.NullString
		kodeStatus, mataUang, claimRetro sql.NullString
	)
	err = r.db.QueryRowContext(ctx, q, id).
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
	// OQ-M6: peserta tercabut (`STS_HAPUS`) tidak dibaca lagi.
	q := fmt.Sprintf(`SELECT ID, %s FROM %s WHERE CLAIM_ID = :1 AND STS_HAPUS IS NULL ORDER BY ID`,
		selectPeserta(), tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, klaimID)
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
		`SELECT a.PREMIUM_LIST_DETAIL_ID, a.ID, `+db.FmtDesimal+`, a.CURRENCY,
		        a.STS_REJECT, a.ACCEPTED_NO, a.ACCEPTATION_DATE, a.KOMITE_ID,
		        `+db.FmtDesimal+`, `+db.FmtDesimal+`, `+db.FmtDesimal+`, `+db.FmtDesimal+`,
		        `+db.FmtDesimal+`, `+db.FmtDesimal+`, a.CURRENCY_ID
		   FROM %s a JOIN %s p ON p.ID = a.PREMIUM_LIST_DETAIL_ID
		  WHERE p.CLAIM_ID = :1 AND p.STS_HAPUS IS NULL
		  ORDER BY a.PREMIUM_LIST_DETAIL_ID, a.ID`,
		"a.CLAIM_AMOUNT",
		"a.SHARE_NUSANTARA_RE", "a.CEDING_RETENTION", "a.SUM_REASURED",
		"a.SUM_INSURED", "a.SHARE_RETRO", "a.RETROCEDED_SHARE",
		adj, pes)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, klaimID)
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
			JumlahKlaim: uang.Money{Currency: mataUang.String},
			CurrencyID:  mataUangID.String,
		}
		// Urutan tujuan mengikuti urutan kolom di SELECT di atas, persis.
		tujuanWarisan := []*uang.Money{
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

// TypeKlaim membaca `Type` klaim dari header klaimnya.
//
// ⛔ Satu-satunya pembaca `Type`, dan ia membacanya dari satu-satunya tempat
// yang menyimpannya - `T_GENERAL_CLAIM.TYPE` sejak migrasi 023 (pindah dari
// `T_WORK_CLAIM`, keputusan work owner 01-10-2026). Tiket 06 AC 5: validasi
// dan penurunan jenis klaim membaca satu field yang sama, bukan dua salinan
// seperti di Pega.
//
// Header klaim dan baris work object berbagi kunci utama (`isiIdentitas`),
// jadi pengenal klaim dapat dipakai apa adanya.
func (r *KlaimLife) TypeKlaim(ctx context.Context, klaimID string) (string, error) {
	tabel, err := r.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf(`SELECT TYPE FROM %s WHERE ID = :1`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var tipe sql.NullString
	err = r.db.QueryRowContext(ctx, q, klaimID).Scan(&tipe)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("repository: header klaim %q tidak ada", klaimID)
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
func (r *KlaimLife) PerbaruiTanggalKejadian(ctx context.Context, tx *db.Tx,
	pesertaID string, dol time.Time) error {
	tabel, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(
		`UPDATE %s SET DATE_OF_LOSS = TO_DATE(:1, 'YYYY-MM-DD HH24:MI:SS') WHERE ID = :2 AND STS_HAPUS IS NULL`,
		tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, utils.FormatTanggalWaktu(dol), pesertaID)
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

func sqlTanggalKlaim(tabel string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET CLAIM_RECEIVED_DATE = TO_DATE(:1, 'YYYY-MM-DD HH24:MI:SS'),
	       COMPLETE_DATE       = TO_DATE(:2, 'YYYY-MM-DD HH24:MI:SS'),
	       CONFIRMATION_DATE   = TO_DATE(:3, 'YYYY-MM-DD HH24:MI:SS')
	 WHERE ID = :4 AND CLAIM_ID = :5 AND STS_HAPUS IS NULL`, tabel)
}

// argTanggal menulis satu tanggal untuk TO_DATE; nil menjadi NULL.
func argTanggal(t *time.Time) any {
	if t == nil {
		return nil
	}
	return utils.FormatTanggalWaktu(*t)
}

// PerbaruiTanggalKlaim menulis ketiga tanggal klaim seorang peserta SEKALIGUS.
//
// Satu pernyataan, sebab satu tombol: `Save` b1910 menyimpan seluruh isian
// dialog bersama-sama (`UpdateDateClaimLife_Act` b252-b384).
//
// ⛔ `CLAIM_ID` ikut di WHERE: peserta klaim lain tidak dapat tersentuh
// walau pengenal pesertanya tertukar.
func (r *KlaimLife) PerbaruiTanggalKlaim(ctx context.Context, tx *db.Tx,
	klaimID, pesertaID string, t models.TanggalKlaim) error {
	tabel, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return err
	}
	q := sqlTanggalKlaim(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, argTanggal(t.TerimaKlaim),
		argTanggal(t.DokumenLengkap), argTanggal(t.Konfirmasi), pesertaID, klaimID)
	if err != nil {
		return fmt.Errorf("repository: menulis tanggal klaim: %w", err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: mencacah baris tersentuh: %w", err)
	}
	if n != 1 {
		return fmt.Errorf("repository: tanggal klaim peserta %q menyentuh %d baris, mau 1",
			pesertaID, n)
	}
	return nil
}

// IsiNomorKlaimKosong menulis CLAIM_NO header klaim yang MASIH kosong.
//
// `Save to RNM` langkah 16-20 bergerbang `CLAIM_NO==""` di setiap langkahnya:
// nomor ditulis hanya sekali. ⛔ `CLAIM_NO IS NULL` ikut di WHERE, dan nol
// baris tersentuh adalah GALAT - klaim yang dinomori pihak lain di antara baca
// dan tulis tidak boleh ditimpa nomor kedua.
func (r *KlaimLife) IsiNomorKlaimKosong(ctx context.Context, tx *db.Tx, klaimID, nomor string) error {
	tabel, err := r.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`UPDATE %s SET CLAIM_NO = :1 WHERE ID = :2 AND CLAIM_NO IS NULL`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, nomor, klaimID)
	if err != nil {
		return fmt.Errorf("repository: menulis CLAIM_NO: %w", err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: mencacah baris tersentuh: %w", err)
	}
	if n != 1 {
		return fmt.Errorf("repository: CLAIM_NO klaim %q menyentuh %d baris, mau 1 (sudah bernomor?)",
			klaimID, n)
	}
	return nil
}

// sqlSudahSaveRNM - penanda turunan "klaim ini sudah pernah di-Save to RNM"
// (OQ-M1, GILIRAN-17): ada baris adjustment klaim itu yang `STS_REJECT`-nya
// tidak NULL. Hasilnya 0 atau 1 (`ROWNUM = 1` sebelum COUNT).
//
// ⚠️ Klaim tanpa satu pun baris adjustment tidak pernah terkunci - Save to
// RNM tidak menulis apa pun untuknya, dan di Pega pun ia tidak bernomor lewat
// jalur baris. Dicatat, bukan ditebak.
func sqlSudahSaveRNM(adj, pes string) string {
	return fmt.Sprintf(`SELECT COUNT(*)
	   FROM %s a JOIN %s p ON p.ID = a.PREMIUM_LIST_DETAIL_ID
	  WHERE p.CLAIM_ID = :1 AND p.STS_HAPUS IS NULL AND a.STS_REJECT IS NOT NULL AND ROWNUM = 1`, adj, pes)
}

// SudahSaveRNM menjawab apakah klaim ini sudah pernah di-Save to RNM - lihat
// `sqlSudahSaveRNM` dan `models.BolehUbahTanggalKlaim`.
func (r *KlaimLife) SudahSaveRNM(ctx context.Context, klaimID string) (bool, error) {
	adj, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT")
	if err != nil {
		return false, err
	}
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return false, err
	}
	q := sqlSudahSaveRNM(adj, pes)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var n int
	if err := r.db.QueryRowContext(ctx, q, klaimID).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: membaca penanda Save to RNM: %w", err)
	}
	return n > 0, nil
}

// sqlCabutPeserta - OQ-M6 (GILIRAN-17): menandai, tidak menghapus
// (ADR-U-0031); hanya peserta milik klaim itu, hanya sekali.
//
// ⛔ `NOT EXISTS` memeriksa ULANG "belum Save to RNM" (b18082) DI DALAM
// pernyataannya (temuan /code-review): gerbang layanan dibaca di luar
// transaksi, dan Save to RNM yang menyela di antaranya tidak boleh diikuti
// pencabutan.
func sqlCabutPeserta(tabel, adj string) string {
	return fmt.Sprintf(`UPDATE %s SET STS_HAPUS = '1'
	 WHERE ID = :1 AND CLAIM_ID = :2 AND STS_HAPUS IS NULL
	   AND NOT EXISTS (SELECT 1 FROM %s a JOIN %s p2 ON p2.ID = a.PREMIUM_LIST_DETAIL_ID
	                    WHERE p2.CLAIM_ID = :3 AND p2.STS_HAPUS IS NULL AND a.STS_REJECT IS NOT NULL)`,
		tabel, adj, tabel)
}

// CabutPeserta menandai seorang peserta dicabut dari klaimnya - tombol
// `DELETE` `InputOSClaimLife` b17865. Nol baris = galat (peserta bukan milik
// klaim, atau sudah tercabut).
func (r *KlaimLife) CabutPeserta(ctx context.Context, tx *db.Tx, klaimID, pesertaID string) error {
	tabel, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return err
	}
	adj, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT")
	if err != nil {
		return err
	}
	q := sqlCabutPeserta(tabel, adj)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, pesertaID, klaimID, klaimID)
	if err != nil {
		return fmt.Errorf("repository: mencabut peserta: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "pencabutan peserta (STS_HAPUS)")
}

// sqlTandaiBarisOutstanding - `Save to RNM` langkah 22.1.3.2.
const sqlTandaiBarisOutstanding = `UPDATE %s SET STS_REJECT = :1 WHERE ID = :2 AND STS_REJECT IS NULL`

// TandaiBarisOutstanding menulis `STS_REJECT = 0` ke SATU baris adjustment
// yang belum berstatus - dan HANYA itu.
//
// `[terverifikasi]` `SaveOutStandingLife_Act` langkah 22.1.3.2 menyetel
// `.STS_REJECT = 0` pada baris adjustment saja (sensus di PerbaruiStatusBaris):
// status peserta, `ACCEPTED_NO`, dan `ACCEPTATION_DATE` tidak disentuh. Karena
// itu bukan PerbaruiStatusBaris, yang menulis keempatnya.
//
// ⛔ `STS_REJECT IS NULL` di WHERE, dan nol baris tersentuh adalah GALAT -
// baris yang diberi status pihak lain sejak dibaca tidak ditimpa.
func (r *KlaimLife) TandaiBarisOutstanding(ctx context.Context, tx *db.Tx, adjID string) error {
	adj, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(sqlTandaiBarisOutstanding, adj)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, kontrak.KodeOutstanding, adjID)
	if err != nil {
		return fmt.Errorf("repository: menandai baris Outstanding: %w", err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: mencacah baris tersentuh: %w", err)
	}
	if n != 1 {
		return fmt.Errorf("repository: menandai baris %q menyentuh %d baris, mau 1 "+
			"(baris yang sudah berstatus sejak dibaca tidak ditimpa)", adjID, n)
	}
	return nil
}

// sqlStatusBarisAdjustment merakit UPDATE baris adjustment bersyarat kode lama.
//
// ⛔ Kode lama KOSONG menjadi `STS_REJECT IS NULL`, bukan `= NULL`. Ronde
// sebelumnya selalu menulis `STS_REJECT = :5`, dan `= NULL` tidak pernah benar
// di SQL - baris yang belum berstatus karena itu TIDAK PERNAH dapat ditulis,
// dan setiap percobaan berakhir "menyentuh 0 baris". (`Save to RNM` kini
// memakai TandaiBarisOutstanding; perbaikan ini tetap, sebab `= NULL` salah
// bagi pemanggil mana pun.)
func sqlStatusBarisAdjustment(adj string, lamaKosong bool) string {
	syarat := "STS_REJECT = :5"
	if lamaKosong {
		syarat = "STS_REJECT IS NULL"
	}
	return fmt.Sprintf(`UPDATE %s SET STS_REJECT = :1, ACCEPTED_NO = :2, ACCEPTATION_DATE = :3
			  WHERE ID = :4 AND %s`, adj, syarat)
}

func argStatusBarisAdjustment(kode, nomorAksep string, tglAksep time.Time, adjID, kodeLama string) []any {
	arg := []any{db.KosongJadiNil(kode), db.KosongJadiNil(nomorAksep), waktuJadiNil(tglAksep), adjID}
	if !kosong(kodeLama) {
		arg = append(arg, kodeLama)
	}
	return arg
}

func kosong(s string) bool { return strings.TrimSpace(s) == "" }

// PerbaruiStatusBaris menulis status ke baris adjustment DAN ke pesertanya.
//
// Header klaim TIDAK ditulis di sini - sumbernya baris yang berbeda; lihat
// CerminkanHeader. Pemanggil menjalankan keduanya dalam satu transaksi.
//
// `[terverifikasi]` sasarannya diturunkan dari sensus penulis `STS_REJECT` di
// seluruh korpus Claim Life - ENAM Property-Set di LIMA rule:
//
//	baris adjustment : SaveOutStandingLife_Act  `.STS_REJECT = 0` (-> TandaiBarisOutstanding)
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
func (r *KlaimLife) PerbaruiStatusBaris(ctx context.Context, tx *db.Tx,
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
		{"baris adjustment", sqlStatusBarisAdjustment(adj, kosong(kodeLama)),
			argStatusBarisAdjustment(kode, nomorAksep, tglAksep, adjID, kodeLama)},
		{"peserta", fmt.Sprintf(
			`UPDATE %s SET STS_REJECT = :1 WHERE ID = :2 AND STS_HAPUS IS NULL`, pes),
			[]any{db.KosongJadiNil(kode), pesertaID}},
	}
	for _, l := range langkah {
		if err := db.PeriksaSQL(l.q); err != nil {
			return err
		}
		hasil, err := tx.ExecContext(ctx, l.q, l.args...)
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
	// Temuan /code-review GILIRAN-18 (OQ-N13): cermin warisan mengikuti status
	// baris - `sqlIkutkanStatusCermin`. ⚠️ Nol baris BUKAN galat: baris putaran
	// Komite tidak punya cermin, dan cermin yang sudah ditulis jalur lain tidak
	// ditimpa.
	lama, err := r.db.Qualify(namaTabelLama)
	if err != nil {
		return err
	}
	// CASEID cermin = CLAIM_ID peserta: kolom CASE_ID dibuang migrasi 023,
	// CASEID = ID klaim.
	q := sqlIkutkanStatusCermin(lama, pes, kosong(kodeLama))
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	arg := []any{db.KosongJadiNil(kode), adjID, pesertaID}
	if !kosong(kodeLama) {
		arg = append(arg, kodeLama)
	}
	if _, err := tx.ExecContext(ctx, q, arg...); err != nil {
		return fmt.Errorf("repository: mencerminkan status ke cermin warisan: %w", err)
	}
	return nil
}

// CerminkanHeader menyalin status BARIS TERAKHIR ke header klaim.
//
// Dipisah dari PerbaruiStatusBaris karena sumbernya baris yang BERBEDA - lihat
// komentar di atas. Pemanggilnya wajib menjalankan keduanya dalam satu
// transaksi yang sama.
func (r *KlaimLife) CerminkanHeader(ctx context.Context, tx *db.Tx,
	klaimID, kode, nomorAksep string) error {
	hdr, err := r.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(
		`UPDATE %s SET STS_REJECT = :1, ACCEPTED_NO = :2 WHERE ID = :3`, hdr)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q,
		db.KosongJadiNil(kode), db.KosongJadiNil(nomorAksep), klaimID)
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
func (r *KlaimLife) PasangPenandaDipilih(ctx context.Context, tx *db.Tx,
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
func (r *KlaimLife) CabutPenandaDipilih(ctx context.Context, tx *db.Tx, pesertaID string) error {
	return r.setelPenandaDipilih(ctx, tx, pesertaID, "false")
}

// setelPenandaDipilih menulis IS_CHECK peserta.
//
// ⛔ Nilainya TEKS - "true"/"false" apa adanya seperti kolomnya, bukan boolean
// Go yang diformat (ADR-U-0022).
func (r *KlaimLife) setelPenandaDipilih(ctx context.Context, tx *db.Tx,
	pesertaID, nilai string) error {

	tabel, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`UPDATE %s SET IS_CHECK = :1 WHERE ID = :2 AND STS_HAPUS IS NULL`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, nilai, pesertaID)
	if err != nil {
		return fmt.Errorf("repository: menyetel IS_CHECK menjadi %q: %w", nilai, err)
	}
	return db.PastikanSatuBaris(hasil, "penyetelan IS_CHECK")
}

// PerbaruiTahap memindahkan kasus di tangga kerja.
//
// ⛔ SATU pernyataan, dan ia TIDAK menyentuh STS_REJECT mana pun. Perpindahan
// tahap memindahkan pekerjaan, bukan memutuskan klaim (ADR-U-0011).
//
// Kolomnya: POSITION (migrasi 023; PY_POSITION di 001), SENDTO_ADMIN, SENDTO_MEDICAL.
// Nol kolom baru, nol langkah migrasi.
//
// ⛔ Tahap ASAL ikut sebagai syarat WHERE. Ia dibaca di luar transaksi, jadi
// dapat berubah di antara baca dan tulis - dan dua pemindahan serentak dapat
// sama-sama menang, yang kedua menimpa yang pertama tanpa jejak. Penjaga yang
// sama sudah dipasang pada perubahan status (tiket 04).
func (r *KlaimLife) PerbaruiTahap(ctx context.Context, tx *db.Tx,
	klaimID, tahapAsal, tahapTujuan, peranTujuan, sendtoAdmin, sendtoMedical string,
	saat time.Time) error {
	tabel, err := r.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	// ⛔ BUTIR at: `TAHAP` ikut ditulis, dan penjaga optimisnya kini
	// memakai TAHAP - bukan POSITION. Alasannya memaksa: pada
	// perpindahan Input Register ⇄ Outstanding Claim `POSITION` TIDAK
	// BERUBAH (keduanya ReasLifeAdmin), sehingga penjaga ber-POSITION
	// tidak dapat mendeteksi kasus yang sudah dipindahkan orang lain
	// sejak dibaca.
	//
	// ⚠️ `NVL` dipakai untuk baris LAMA yang `TAHAP`-nya masih kosong:
	// tanpa itu setiap perpindahan pertama baris lama akan menyentuh nol
	// baris dan gagal, padahal tidak ada yang salah dengannya.
	q := fmt.Sprintf(`UPDATE %s SET POSITION = :1, SENDTO_ADMIN = :2,
		 SENDTO_MEDICAL = :3, TGL_UPDATE = :4, TAHAP = :5
		 WHERE ID = :6 AND NVL(TAHAP, :7) = :8`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, db.KosongJadiNil(peranTujuan),
		db.KosongJadiNil(sendtoAdmin), db.KosongJadiNil(sendtoMedical),
		waktuJadiNil(saat), db.KosongJadiNil(tahapTujuan),
		klaimID, db.KosongJadiNil(tahapAsal), db.KosongJadiNil(tahapAsal))
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

// TahapKlaim membaca POSITION baris work object sebuah klaim.
func (r *KlaimLife) TahapKlaim(ctx context.Context, klaimID string) (string, error) {
	_, peran, err := r.TahapDanPeran(ctx, klaimID)
	return peran, err
}

// TahapDanPeran membaca TAHAP dan POSITION sebuah kasus sekaligus.
//
// ⛔ BUTIR at. Keduanya dibaca BERSAMA, dalam satu kueri, karena keduanya
// menjawab pertanyaan yang berbeda dan pemanggil memerlukan keduanya:
// `TAHAP` mengatakan di anak tangga mana kasusnya berdiri, `POSITION`
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
	q := fmt.Sprintf(`SELECT TAHAP, POSITION FROM %s WHERE ID = :1`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", "", err
	}
	var kolomTahap, posisi sql.NullString
	err = r.db.QueryRowContext(ctx, q, klaimID).Scan(&kolomTahap, &posisi)
	if err == sql.ErrNoRows {
		return "", "", fmt.Errorf("%w: %q", ErrWorkTidakAda, klaimID)
	}
	if err != nil {
		return "", "", fmt.Errorf("repository: membaca tahap kasus: %w", err)
	}
	return kolomTahap.String, posisi.String, nil
}

// ErrWorkTidakAda - baris `T_WORK_CLAIM` klaim itu tidak ada.
//
// ⚠️ Bersentinel sendiri, sebab pemanggil harus dapat MEMBEDAKANNYA dari
// galat basis data. Aplikasi ini tidak pernah MENYISIPKAN ke T_WORK_CLAIM -
// baris itu lahir di sistem lama - sehingga klaim tanpa baris work adalah
// keadaan yang nyata, bukan kerusakan. Layar Detail memperlakukannya sebagai
// "tahap tidak diketahui" dan karena itu tidak menawarkan Close Claim;
// galat lain tetap menggagalkan pembacaan.
var ErrWorkTidakAda = errors.New("repository: work object tidak ada")

// StatusWorkKlaim membaca `STATUS_WORK` baris work object sebuah klaim.
//
// Butir bb. Kosong berarti kasusnya BELUM ditutup - bukan tidak diketahui
// (ADR-U-0027), dan seluruh baris yang sudah ada memang belum ditutup.
func (r *KlaimLife) StatusWorkKlaim(ctx context.Context, klaimID string) (string, error) {
	tabel, err := r.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf(`SELECT STATUS_WORK FROM %s WHERE ID = :1`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var status sql.NullString
	err = r.db.QueryRowContext(ctx, q, klaimID).Scan(&status)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("%w: %q", ErrWorkTidakAda, klaimID)
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca status kerja kasus: %w", err)
	}
	return status.String, nil
}

// TutupKasus menutup kasus: status kerja diisi, TAHAP DIKOSONGKAN.
//
// Butir bb. Padanan `Call FinishAssignment` (`ProtectCloseClaim_act` b838)
// yang seluruh parameternya kosong, ditambah `pyWorkStatus` shape End1
// (`Register_Flow.xml` b899).
//
// ⛔ `TAHAP = NULL`, bukan tahap kelima. Penugasannya memang SELESAI, dan
// kotak masuk adalah worklist - kasus yang tertutup hilang dari keempat tab.
// Tahap kelima "Selesai" akan menjadi antrean yang tidak pernah dikerjakan
// siapa pun, dan XML tidak menyebutnya.
//
// ⛔ `STATUS_WORK IS NULL` ikut di WHERE. Tanpa itu, dua penutupan bersamaan
// sama-sama berhasil dan yang kedua menimpa TGL_UPDATE penutupan pertama -
// jejak yang menunjuk waktu yang salah. Dengan klausa itu yang kedua
// menyentuh nol baris dan berkata jelas.
//
// ⚠️ `NVL(TAHAP, :n)` sama seperti PerbaruiTahap: baris LAMA yang TAHAP-nya
// masih kosong tetap dapat ditutup dari tahap yang disimpulkan POSITION.
func (r *KlaimLife) TutupKasus(ctx context.Context, tx *db.Tx,
	klaimID, tahapAsal, statusWork string, saat time.Time) error {

	tabel, err := r.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`UPDATE %s SET STATUS_WORK = :1, TAHAP = NULL, TGL_UPDATE = :2
		 WHERE ID = :3 AND NVL(TAHAP, :4) = :5 AND STATUS_WORK IS NULL`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, db.KosongJadiNil(statusWork),
		waktuJadiNil(saat), klaimID,
		db.KosongJadiNil(tahapAsal), db.KosongJadiNil(tahapAsal))
	if err != nil {
		return fmt.Errorf("repository: menutup kasus: %w", err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: mencacah baris work: %w", err)
	}
	if n != 1 {
		return fmt.Errorf("repository: penutupan kasus menyentuh %d baris, mau 1 "+
			"(kasus yang tahapnya sudah berubah sejak dibaca, atau yang sudah "+
			"tertutup, tidak ditimpa)", n)
	}
	return nil
}

// PerbaruiKomiteID menautkan baris adjustment ke kasus Komite yang baru lahir.
//
// Tiket 10 AC 61. Rujukan, BUKAN salinan: roster dan keputusan per anggota
// tetap milik konteks Komite Claim Life.
//
// KOMITE_ID lama ikut di WHERE dan wajib NULL. Baris yang sudah tertaut sejak
// dibaca tidak ditimpa - tanpa klausa itu dua penyerahan bersamaan menghasilkan
// dua kasus komite dan hanya satu yang teringat.
func (r *KlaimLife) PerbaruiKomiteID(ctx context.Context, tx *db.Tx,
	adjID, komiteID string) error {

	tabel, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(
		`UPDATE %s SET KOMITE_ID = :1 WHERE ID = :2 AND KOMITE_ID IS NULL`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, db.KosongJadiNil(komiteID), adjID)
	if err != nil {
		return fmt.Errorf("repository: menautkan baris ke Komite: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "penautan baris ke Komite")
}

// AmbilDokumen membaca seluruh dokumen pendukung satu klaim, per peserta.
//
// Meniru `Activity/LoadDocumentLife_ACT.xml`, dibaca sebagai pohon
// 27-09-2026 (nomor baris hasil `sed -e 's/></>\n</g'`):
//
//	langkah 1     b320
//	  1.1 b369  `Param.inskey` = `@If(pyWorkCover.pzInsKey="",pyWorkPage.pzInsKey,
//	            pyWorkCover.pzInsKey)` b388-389
//	  1.2 b495  `Obj-Browse` - `ObjClass` `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM`,
//	            `PageName` `DOCUMENT_CLAIM` (nama halaman warisan), `RowKey` `ID`,
//	            saring `Field .KATEGORI_1` `Condition =` `Value .DOCUMENT`
//	  1.3 b778  `.DocumentClaimList` = `DOCUMENT_CLAIM.pxResults` (warisan) b797
//	  (langkah 1 sendiri: prasyarat b1404 `.DOCUMENT==""` WhenTrue=3 -> LEWATI,
//	  MENGULANG b1423)
//	  1.4 b904  ⛔ TER-REMARK (`//` b913), beserta 1.4.1 dan 1.4.2 di bawah
//	    `Local.ImageID` = `.T_STORAGE_ID` b924-925; `DataImage.URLImage` = "" b970
//	    1.4.1 b1011 `Call GetUrlGoogleStorage_Act`
//	    1.4.2 b1129 prasyarat b1310 `DataImage.URLImage==""` WhenTrue=3 -> LEWATI
//	          `…DocumentClaimList(<APPEND>).IMAGEID`  = `Local.ImageID`   b1149
//	          `…DocumentClaimList(<LAST>).URLPUBLIC`  = `DataImage.URLImage` b1195
//	          `…DocumentClaimList(<LAST>).PNOTE`      = `.NAMAFILE`       b1216
//	          `…DocumentClaimList(<LAST>).pyMemoo`    = `.NAMAFILE`       b1237
//
// ⚠️ Nomor baris di brief lanjutan 12 §2 berbeda jauh dari bacaan ini
// (b755/b861/b1164/b1377/b1516/b1562/b1583). Yang dipakai adalah bacaan ini,
// dan selisihnya DICATAT - aturannya "bila bacaan Anda berbeda, yang menang
// adalah bacaan Anda".
//
// ⛔ RALAT 28-09-2026 (sensus remark GILIRAN-12): yang berjalan adalah 1.3
// b797 - `.DocumentClaimList` = SELURUH baris browse, tanpa URL. Bab ini dulu
// menyebut DUA penyimpangan sadar; yang kedua (baris ber-URL kosong dibuang,
// b1310) milik 1.4.2, langkah mati. Mengembalikan seluruh baris adalah
// PARITAS, bukan penyimpangan.
//
// ⛔ SATU PENYIMPANGAN SADAR, dan ia dinyatakan.
//
//  1. Saringan Pega adalah `KATEGORI_1 = <.DOCUMENT peserta>` pada halaman
//     peserta - BUKAN pengenal kasus. Kolom `DOCUMENT` itu tidak ada di
//     `T_CLAIMLF_PREMIUMLIST_DETAIL`, sehingga saringan itu tidak dapat
//     ditiru apa adanya. Yang dipakai: FK `PREMIUM_LIST_DETAIL_ID` - relasi
//     yang di model baru MEMANG memiliki dokumen itu. Dilaporkan OQ-J.
func (r *KlaimLife) AmbilDokumen(ctx context.Context, klaimID string) (
	map[string][]models.Dokumen, error) {

	dok, err := r.db.Qualify("T_CLAIMLF_DOCUMENT")
	if err != nil {
		return nil, err
	}
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(
		`SELECT d.PREMIUM_LIST_DETAIL_ID, d.ID, d.NAMA_FILE, d.MIME,
		        d.KATEGORI_1, d.KATEGORI_2, d.T_STORAGE_ID, d.TANGGAL,
		        d.PAYMENT_DATE
		   FROM %s d JOIN %s p ON p.ID = d.PREMIUM_LIST_DETAIL_ID
		  WHERE p.CLAIM_ID = :1 AND p.STS_HAPUS IS NULL
		  ORDER BY d.PREMIUM_LIST_DETAIL_ID, d.ID`, dok, pes)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, q, klaimID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca dokumen klaim: %w", err)
	}
	defer baris.Close()

	keluar := map[string][]models.Dokumen{}
	for baris.Next() {
		var (
			pesertaID                       sql.NullString
			id                              sql.NullInt64
			nama, mime, kat1, kat2, storage sql.NullString
			tanggal, bayar                  sql.NullTime
		)
		if err := baris.Scan(&pesertaID, &id, &nama, &mime, &kat1, &kat2,
			&storage, &tanggal, &bayar); err != nil {
			return nil, fmt.Errorf("repository: memindai dokumen: %w", err)
		}
		d := models.Dokumen{
			ID:         id.Int64,
			PesertaID:  pesertaID.String,
			NamaFile:   nama.String,
			Mime:       mime.String,
			Kategori1:  kat1.String,
			Kategori2:  kat2.String,
			TStorageID: storage.String,
		}
		// ⛔ NULL tetap nil, bukan tanggal nol. Tanggal nol adalah tahun 1
		// Masehi, dan pembaca hilir tidak dapat membedakannya dari kolom
		// yang memang kosong (ADR-U-0027).
		if tanggal.Valid {
			t := tanggal.Time
			d.Tanggal = &t
		}
		if bayar.Valid {
			t := bayar.Time
			d.PaymentDate = &t
		}
		keluar[d.PesertaID] = append(keluar[d.PesertaID], d)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca dokumen klaim: %w", err)
	}
	return keluar, nil
}
