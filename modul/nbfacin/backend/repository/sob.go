package repository

// Pilihan Source of Business - popup Change SOB (tiket 33). Syarat = perintah work owner
// 03-10-2026 ("StatusActive=1, AgentType2 != "LIFE INSURANCE", ClientID is not null ... tanpa
// memperhatikan huruf besar dan kecil"), SAMA dengan RD korpus `[terverifikasi]`
// `NB FacIn\ReportDefinition\BrowseAgentNonLife_RD.xml` (kelas ASM-FW-GISFW-Int-AGENT,
// dipakai `Section\SourceHierarki.xml`): logika `B AND E AND A AND C` - B `.AgentType2 !=
// "LIFE INSURANCE"`, A `.StatusActive = "1"` (Text), C `.ClientID IS NOT NULL`, E `.ClientName
// Contains` (tidak peka huruf); urut `.ID` DESC. RD menerjemahkan `!=` ke `<>` SQL, jadi baris
// ber-AgentType2 NULL juga tidak ikut di Pega `[dugaan]`. ⛔ Keputusan work owner 03-10-2026
// (butir 79.1): baris ber-AgentType2 NULL TETAP ikut - yang dibuang hanya 'LIFE INSURANCE'.
//
// Tabel `[terverifikasi]` DDL `D:\migrasi\RNM\DDL\AGENT.txt` (POOLDATA.AGENT, 28 kolom, tanpa PK;
// semua kolom di bawah VARCHAR2(1000 BYTE)). ⚠️ Kolom tipe-2 di DDL dieja **AGENTTPYE2** (salah eja
// di basis data, juga AGENTTPYE) - BUKAN AGENTTYPE2; `.AgentType2` -> AGENTTPYE2 `[dugaan]` (satu-
// satunya kolom tipe-2; nama itu tidak muncul di korpus Pega).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelAgent - tabel sumber SOB (POOLDATA.AGENT; skema lewat db.Qualify, pola akun.go).
	TabelAgent = "AGENT"
	// Kolom AGENT `[terverifikasi]` DDL - lihat kepala berkas.
	kolomAgentID       = "ID"
	kolomAgentClientID = "CLIENTID"
	kolomAgentNama     = "CLIENTNAME"
	kolomAgentStatus   = "STATUSACTIVE"
	kolomAgentTipe2    = "AGENTTPYE2" // ejaan DDL apa adanya (bukan AGENTTYPE2)
	// StatusAgentAktif, TipeAgentDibuang - nilai syarat (RD A dan B), teks.
	StatusAgentAktif = "1"
	TipeAgentDibuang = "LIFE INSURANCE"
)

// PembacaSOB - satu halaman pilihan SOB beserta cacah seluruhnya.
type PembacaSOB interface {
	CariSOB(ctx context.Context, cari string, offset, ukuran int) ([]models.SOB, int, error)
}

// SOBOracle - PembacaSOB atas Oracle.
type SOBOracle struct{ db *db.DB }

// NewSOBOracle merakit pembaca tabel agent.
func NewSOBOracle(d *db.DB) *SOBOracle { return &SOBOracle{db: d} }

// syaratSOB - tiga syarat work owner; :1 status, :2 tipe dibuang. Tipe NULL ikut (butir 79.1).
var syaratSOB = fmt.Sprintf("%s = :1 AND (%s IS NULL OR %s <> :2) AND %s IS NOT NULL",
	kolomAgentStatus, kolomAgentTipe2, kolomAgentTipe2, kolomAgentClientID)

// ErrSOBTidakSah - kode SOB yang dikirim tidak ada di AGENT atau tidak lolos syarat tiket 33.
var ErrSOBTidakSah = errors.New("repository: kode SOB tidak ditemukan di AGENT atau tidak lolos syarat")

// saringSOB - syarat + (bila cari) "mengandung" tidak peka huruf atas ID / ClientID / nama.
func saringSOB(cari bool) string {
	if !cari {
		return " WHERE " + syaratSOB
	}
	return fmt.Sprintf(` WHERE %s AND (UPPER(%s) LIKE :3 ESCAPE '\' OR UPPER(%s) LIKE :4 ESCAPE '\' OR UPPER(%s) LIKE :5 ESCAPE '\')`,
		syaratSOB, kolomAgentID, kolomAgentClientID, kolomAgentNama)
}

func sqlCariSOB(tabel string, cari bool) string {
	n := 3
	if cari {
		n = 6
	}
	return fmt.Sprintf("SELECT %s, %s, %s FROM %s%s ORDER BY %s DESC OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY",
		kolomAgentID, kolomAgentClientID, kolomAgentNama, tabel, saringSOB(cari), kolomAgentID, n, n+1)
}

func sqlCacahSOB(tabel string, cari bool) string {
	return "SELECT COUNT(*) FROM " + tabel + saringSOB(cari)
}

// sqlCekSOB - kode yang dikirim layar: cacah baris yang lolos syarat + namanya (MAX - tabel
// tanpa PK yang diketahui).
func sqlCekSOB(tabel string) string {
	return fmt.Sprintf("SELECT COUNT(*), MAX(%s) FROM %s WHERE %s AND %s = :3", kolomAgentNama, tabel, syaratSOB, kolomAgentID)
}

// CariSOB - lihat PembacaSOB.
func (r *SOBOracle) CariSOB(ctx context.Context, cari string, offset, ukuran int) ([]models.SOB, int, error) {
	q, err := r.db.Qualify(TabelAgent)
	if err != nil {
		return nil, 0, err
	}
	pola := PolaCari(cari)
	arg := []any{StatusAgentAktif, TipeAgentDibuang}
	if pola != "" {
		arg = append(arg, pola, pola, pola)
	}
	var total int
	if err := r.db.QueryRowContext(ctx, sqlCacahSOB(q, pola != ""), arg...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository: cacah %s: %w", TabelAgent, err)
	}
	baris, err := r.db.QueryContext(ctx, sqlCariSOB(q, pola != ""), append(arg, offset, ukuran)...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: baca %s: %w", TabelAgent, err)
	}
	defer baris.Close()
	hasil := []models.SOB{}
	for baris.Next() {
		var id, clientID, nama sql.NullString
		if err := baris.Scan(&id, &clientID, &nama); err != nil {
			return nil, 0, fmt.Errorf("repository: %s: %w", TabelAgent, err)
		}
		hasil = append(hasil, models.SOB{ID: id.String, ClientID: clientID.String, Name: nama.String})
	}
	if err := baris.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository: %s: %w", TabelAgent, err)
	}
	return hasil, total, nil
}

// namaAgent - nama AGENT (CLIENTNAME) untuk kode `id` yang lolos syarat tiket 33 - dipakai SOB
// (tiket 33) dan Ceding Co (tiket 34); ErrSOBTidakSah bila tidak ada.
func namaAgent(ctx context.Context, d *db.DB, tx *db.Tx, id string) (string, error) {
	q, err := d.Qualify(TabelAgent)
	if err != nil {
		return "", err
	}
	teks := sqlCekSOB(q)
	if err := db.PeriksaSQL(teks); err != nil {
		return "", err
	}
	var n int
	var nama sql.NullString
	if err := tx.QueryRowContext(ctx, teks, StatusAgentAktif, TipeAgentDibuang, id).Scan(&n, &nama); err != nil {
		return "", fmt.Errorf("repository: cek SOB %s: %w", TabelAgent, err)
	}
	if n == 0 {
		return "", ErrSOBTidakSah
	}
	return nama.String, nil
}
