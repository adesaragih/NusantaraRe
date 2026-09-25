package repository

import (
	"context"
	"database/sql"
	"fmt"

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
	q := fmt.Sprintf(
		`SELECT ID, CLAIM_NO, POLICY_NO, BUSINESS_NAME, STS_REJECT FROM %s WHERE ID = :1`, tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}

	var (
		kID                            string
		nomorKlaim, nomorPolis, bisnis sql.NullString
		kodeStatus                     sql.NullString
	)
	err = r.db.sql.QueryRowContext(ctx, q, id).Scan(&kID, &nomorKlaim, &nomorPolis, &bisnis, &kodeStatus)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository: membaca header klaim: %w", err)
	}
	return &models.Klaim{
		ID:         kID,
		NomorKlaim: nomorKlaim.String,
		NomorPolis: nomorPolis.String,
		NamaBisnis: bisnis.String,
		KodeStatus: kodeStatus.String,
	}, nil
}

// AmbilPeserta membaca seluruh peserta milik satu klaim, terurut stabil.
func (r *KlaimLife) AmbilPeserta(ctx context.Context, klaimID string) ([]models.Peserta, error) {
	tabel, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(
		`SELECT ID, PL_NUMBER, POLICY_NO, CERTIFICATE_NO, CURRENCY
		   FROM %s WHERE CLAIM_ID = :1 ORDER BY ID`, tabel)
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
		var pl, polis, sertifikat, mataUang sql.NullString
		if err := rows.Scan(&id, &pl, &polis, &sertifikat, &mataUang); err != nil {
			return nil, fmt.Errorf("repository: membaca baris peserta: %w", err)
		}
		out = append(out, models.Peserta{
			ID:              id,
			NomorPremiList:  pl.String,
			NomorPolis:      polis.String,
			NomorSertifikat: sertifikat.String,
			MataUang:        mataUang.String,
			Baris:           []models.BarisAdjustment{},
		})
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
	q := fmt.Sprintf(
		`SELECT a.PREMIUM_LIST_DETAIL_ID, a.ID, `+fmtDesimal+`, a.CURRENCY,
		        a.STS_REJECT, a.ACCEPTED_NO, a.ACCEPTATION_DATE, a.KOMITE_ID
		   FROM %s a JOIN %s p ON p.ID = a.PREMIUM_LIST_DETAIL_ID
		  WHERE p.CLAIM_ID = :1
		  ORDER BY a.PREMIUM_LIST_DETAIL_ID, a.ID`, "a.CLAIM_AMOUNT", adj, pes)
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
		if err := rows.Scan(&pesertaID, &id, &jumlah, &mataUang,
			&kodeStatus, &nomorAksep, &tglAksep, &komiteID); err != nil {
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
