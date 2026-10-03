package repository

// Pembaca efek keluar kasus Komite - tiket 08 Komite Claim Life.
//
// Membaca outbox `T_LOG_SERVICE_RNM` baris `MODUL = KOMITELIFE` - per kasus
// (`RUJUKAN` = pengenal kasus komite), dan seluruh yang gagal permanen untuk
// laporan harian. Nol tulisan.
//
// ⚠️ `GALAT_TERAKHIR` sengaja TIDAK dibaca: pesan galat dapat menyebut nama
// objek basis data atau kunci kategori; layar cukup tahu efek MANA dan SEJAK
// KAPAN (AC tiket 08). Rinciannya tetap di tabel untuk operator.

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
)

// ModulOutboxKomite = `services.ModulKomiteLife` - diuji sama.
const ModulOutboxKomite = "KOMITELIFE"

// EfekKasusKomite adalah satu baris efek.
type EfekKasusKomite struct {
	KasusID   string
	Jenis     string
	Status    string
	Percobaan int
	Dibuat    time.Time
	// Diperbarui kosong bila baris belum pernah dipungut.
	Diperbarui sql.NullTime
}

func sqlEfekKasusKomite(tabel string) string {
	return fmt.Sprintf(`SELECT RUJUKAN, JENIS_EFEK, STATUS, PERCOBAAN, DIBUAT, DIPERBARUI
	  FROM %s WHERE MODUL = :1 AND (RUJUKAN = :2 OR RUJUKAN LIKE :3)
	 ORDER BY DIBUAT, ID`, tabel)
}

func sqlEfekPerluIntervensi(tabel string) string {
	return fmt.Sprintf(`SELECT RUJUKAN, JENIS_EFEK, STATUS, PERCOBAAN, DIBUAT, DIPERBARUI
	  FROM %s WHERE MODUL = :1 AND STATUS = :2 ORDER BY DIPERBARUI, ID`, tabel)
}

func (r *InboxKomite) bacaEfek(ctx context.Context, q string, args ...any) ([]EfekKasusKomite, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca efek komite: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []EfekKasusKomite
	for rows.Next() {
		var e EfekKasusKomite
		var p sql.NullInt64
		if err := rows.Scan(&e.KasusID, &e.Jenis, &e.Status, &p, &e.Dibuat, &e.Diperbarui); err != nil {
			return nil, fmt.Errorf("repository: memindai efek komite: %w", err)
		}
		e.Percobaan = int(p.Int64)
		out = append(out, e)
	}
	return out, rows.Err()
}

// EfekKasus membaca seluruh efek satu kasus.
func (r *InboxKomite) EfekKasus(ctx context.Context, kasusID string) ([]EfekKasusKomite, error) {
	tabel, err := r.db.Qualify("T_LOG_SERVICE_RNM")
	if err != nil {
		return nil, err
	}
	// Email berujukan `kasus#Tn` (services.RujukanEfekKomite).
	return r.bacaEfek(ctx, sqlEfekKasusKomite(tabel), ModulOutboxKomite, kasusID, kasusID+"#%")
}

// EfekPerluIntervensi membaca seluruh efek Komite yang gagal permanen.
func (r *InboxKomite) EfekPerluIntervensi(ctx context.Context) ([]EfekKasusKomite, error) {
	tabel, err := r.db.Qualify("T_LOG_SERVICE_RNM")
	if err != nil {
		return nil, err
	}
	return r.bacaEfek(ctx, sqlEfekPerluIntervensi(tabel), ModulOutboxKomite, outbox.StatusEfekGagalPermanen)
}
