package repository

// Tabel akun untuk popup ChooseAccount (tiket 27). Skema HANYA dari DDL
// `D:\migrasi\RNM\DDL\T_M_ACCOUNT.txt` `[terverifikasi]`: ID VARCHAR2(255 CHAR) NOT
// NULL, GROUPBUSINESSID VARCHAR2(32 CHAR), GROUPBUSINESS VARCHAR2(64 CHAR), INSUREDID
// VARCHAR2(255 CHAR) NOT NULL, INSUREDNAME VARCHAR2(64 CHAR). Rule Pega ChooseAccount
// tidak ada di korpus; kolom layar dan "mengandung" dari jawaban work owner.
//
// Hanya membaca kelima kolom itu; tidak ada `SELECT *`, nol tulisan. Masukan pencarian
// selalu parameter terikat - tidak pernah disambung ke teks SQL.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// TabelAkun - nama tabel, ejaan DDL.
const TabelAkun = "T_M_ACCOUNT"

// PembacaAkun - sumber baris akun berhalaman.
type PembacaAkun interface {
	// CariAkun - satu halaman (offset, ukuran) baris yang "mengandung" cari (lihat
	// PolaCari), beserta cacah seluruh baris yang cocok. cari kosong = tanpa saringan.
	CariAkun(ctx context.Context, cari string, offset, ukuran int) ([]models.Akun, int, error)
}

// AkunOracle - PembacaAkun atas Oracle.
type AkunOracle struct{ db *db.DB }

// NewAkunOracle merakit pembaca tabel akun.
func NewAkunOracle(d *db.DB) *AkunOracle { return &AkunOracle{db: d} }

// PolaCari - "mengandung", TIDAK peka huruf besar-kecil (keputusan work owner
// 03-10-2026, butir 75 - "itu buatin tanpa liat huruf besar atau kecil", membatalkan
// 73.1): huruf besar, wildcard LIKE (`\` `%` `_`) diloloskan, lalu diapit `%`. Escape
// sepola PolaCari masterproductnamelife (disalin; modul tidak saling mengimpor).
// Kosong (sesudah dipangkas) = "" = tanpa saringan. ⚠️ Huruf besar dibuat
// strings.ToUpper Go, kolom dengan UPPER Oracle: untuk huruf non-ASCII keduanya dapat
// berbeda (`[dugaan]`, belum diuji ke Oracle).
func PolaCari(kata string) string {
	k := strings.ToUpper(strings.TrimSpace(kata))
	if k == "" {
		return ""
	}
	k = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(k)
	return "%" + k + "%"
}

// saringAkun - A69: tiga kolom yang tampil di layar, OR, UPPER di sisi kolom (tidak peka
// huruf, butir 75); :1..:3 bernilai pola yang sama.
const saringAkun = ` WHERE (UPPER(INSUREDID) LIKE :1 ESCAPE '\' OR UPPER(INSUREDNAME) LIKE :2 ESCAPE '\' OR UPPER(GROUPBUSINESS) LIKE :3 ESCAPE '\')`

// sqlCariAkun - satu halaman. A72: urut INSUREDID lalu ID (INSUREDID tidak dijamin unik).
func sqlCariAkun(tabel string, saring bool) string {
	w, n := "", 1
	if saring {
		w, n = saringAkun, 4
	}
	return fmt.Sprintf("SELECT ID, GROUPBUSINESSID, GROUPBUSINESS, INSUREDID, INSUREDNAME FROM %s%s"+
		" ORDER BY INSUREDID, ID OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY", tabel, w, n, n+1)
}

// sqlCacahAkun - cacah seluruh baris yang cocok (dibaca terpisah: total bukan panjang halaman).
func sqlCacahAkun(tabel string, saring bool) string {
	w := ""
	if saring {
		w = saringAkun
	}
	return "SELECT COUNT(*) FROM " + tabel + w
}

// CariAkun - lihat PembacaAkun.
func (r *AkunOracle) CariAkun(ctx context.Context, cari string, offset, ukuran int) ([]models.Akun, int, error) {
	q, err := r.db.Qualify(TabelAkun)
	if err != nil {
		return nil, 0, err
	}
	pola := PolaCari(cari)
	saring := pola != ""
	var arg []any
	if saring {
		arg = []any{pola, pola, pola}
	}
	var total int
	if err := r.db.QueryRowContext(ctx, sqlCacahAkun(q, saring), arg...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository: cacah %s: %w", TabelAkun, err)
	}
	baris, err := r.db.QueryContext(ctx, sqlCariAkun(q, saring), append(arg, offset, ukuran)...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: baca %s: %w", TabelAkun, err)
	}
	defer baris.Close()
	hasil := []models.Akun{}
	for baris.Next() {
		var id, gbID, gb, insID, insNama sql.NullString
		if err := baris.Scan(&id, &gbID, &gb, &insID, &insNama); err != nil {
			return nil, 0, fmt.Errorf("repository: %s: %w", TabelAkun, err)
		}
		hasil = append(hasil, models.Akun{ID: id.String, GroupBusinessID: gbID.String, GroupBusiness: gb.String,
			InsuredID: insID.String, InsuredName: insNama.String})
	}
	if err := baris.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository: %s: %w", TabelAkun, err)
	}
	return hasil, total, nil
}
