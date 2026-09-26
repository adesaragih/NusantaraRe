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
