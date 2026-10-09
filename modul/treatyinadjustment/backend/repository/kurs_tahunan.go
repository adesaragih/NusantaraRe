package repository

// Grid Rate of Exchange layar Adjustment — dari `TREATYEXCHANGEYEARLY`.
//
// ---------------------------------------------------------------------
// ⛔ BUKAN `CurrencyList` dokumen, dan bukan `T_TREATY_CURRENCY`
// ---------------------------------------------------------------------
// Pemilik proses menegaskannya dua kali (4 dan 6 Oktober 2026): Rate of
// Exchange dibaca dari `TREATYEXCHANGEYEARLY`, tabel WARISAN yang hidup.
// `CurrencyList` di dokumen hanyalah SALINAN tabel itu
// (`KOREKSI-ERD-VERSUS-POOLDATA.md` §1.2), dan ia tidak pernah dimuat ke
// tabel pendaratan mana pun — itulah sebab grid kedua panel layar ini
// selalu "No items" sebelum berkas ini ada.
//
// ⚠️ Disaring menurut `TREATYYEAR` sisi itu sendiri, sama dengan pembaca
// Treaty In (`treatyin/backend/repository.BacaKursTahunan`) — modul ini
// tidak boleh mengimpornya, jadi aturannya disalin, bukan dipanggil. Panel
// Old memakai Treaty Year Old, panel New memakai Treaty Year New: revisi
// yang memindah tahun kontrak menampilkan kurs tahun masing-masing.
//
// Kolom → kunci Pega grid (`kerangka.gen.ts` GRID_KURS):
//
//	CURRENCY   -> Currency     IDCURRENCY -> CurrencyID
//	TOIDR      -> Conversion   STARTDATE  -> PeriodStart   ENDDATE -> PeriodEnd

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyinadjustment/backend/models"
)

// TabelKursTahunan - master kurs warisan. Dibaca, tidak pernah ditulis.
const TabelKursTahunan = "TREATYEXCHANGEYEARLY"

// LarikKurs - nama larik grid Rate of Exchange di halaman Pega.
const LarikKurs = "CurrencyList"

// bacaKursTahunan membaca baris kurs satu Treaty Year sebagai larik berkunci
// Pega. Tahun kosong → nol baris, bukan galat: kontrak tanpa tahun tidak
// punya apa pun untuk disaring.
func (g *Gudang) bacaKursTahunan(ctx context.Context, tahun string) ([]map[string]string, error) {
	tahun = strings.TrimSpace(tahun)
	if tahun == "" {
		return nil, nil
	}
	nama, err := g.db.Qualify(TabelKursTahunan)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT CURRENCY, TO_CHAR(IDCURRENCY), TO_CHAR(TOIDR), STARTDATE, ENDDATE
		FROM %s WHERE TREATYYEAR = :1 ORDER BY CURRENCY`, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, tahun)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s tahun %s: %w", TabelKursTahunan, tahun, err)
	}
	defer func() { _ = rows.Close() }()
	var out []map[string]string
	for rows.Next() {
		var cur, idCur, toidr, mulai, akhir sql.NullString
		if err := rows.Scan(&cur, &idCur, &toidr, &mulai, &akhir); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", TabelKursTahunan, err)
		}
		out = append(out, map[string]string{
			"Currency":    cur.String,
			"CurrencyID":  idCur.String,
			"Conversion":  toidr.String,
			"PeriodStart": mulai.String,
			"PeriodEnd":   akhir.String,
		})
	}
	return out, rows.Err()
}

// TabelKursKontrak — penghubung kontrak ↔ baris `TREATYEXCHANGEYEARLY`
// (migrasi 455, 8 Oktober 2026). Disalin dari Treaty In
// (`treatyin/backend/repository/kurs_kontrak.go`), bukan diimpor.
const TabelKursKontrak = "T_TREATY_KURS"

// akarKontrak — ID kontrak ASAL sebuah sisi penyesuaian: `1002306/R01` →
// `1002306` (kurs milik kontrak dicatat atas ID asalnya).
func akarKontrak(sisi models.SisiPenyesuaian) string {
	id := strings.TrimSpace(sisi.Medan["ID"])
	if id == "" {
		id = strings.TrimSpace(sisi.Medan["OLDID"])
	}
	return strings.TrimSpace(strings.SplitN(id, "/", 2)[0])
}

// bacaKursKontrak — HANYA kurs milik kontrak `masterID` (laporan pemakai
// 8 Oktober 2026: Add Revision menampilkan tujuh kurs padahal kontraknya
// hanya punya dua — *"biarkan apa yg di input user yang tampil"*).
// `punya` false bila kontrak itu BELUM PERNAH punya catatan (atau 455 belum
// terpasang) — hanya itu yang boleh jatuh ke kurs tahun. Catatan penanda
// `-` (grid dikosongkan pemakai, Treaty In `kursKontrakKosong`) = punya,
// nol baris.
func (g *Gudang) bacaKursKontrak(ctx context.Context, masterID, tahun string) (baris []map[string]string, punya bool, err error) {
	tahun = strings.TrimSpace(tahun)
	if masterID == "" || tahun == "" {
		return nil, false, nil
	}
	kol, err := g.kolomTerpasang(ctx, TabelKursKontrak)
	if err != nil || kol == nil {
		return nil, false, err
	}
	hub, err := g.db.Qualify(TabelKursKontrak)
	if err != nil {
		return nil, false, err
	}
	nama, err := g.db.Qualify(TabelKursTahunan)
	if err != nil {
		return nil, false, err
	}
	cacah := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE MASTERID = :1`, hub)
	if err := db.PeriksaSQL(cacah); err != nil {
		return nil, false, err
	}
	var n int
	if err := g.db.QueryRowContext(ctx, cacah, masterID).Scan(&n); err != nil {
		return nil, false, fmt.Errorf("repository: menghitung kurs kontrak %s: %w", masterID, err)
	}
	if n == 0 {
		return nil, false, nil
	}
	q := fmt.Sprintf(`SELECT k.CURRENCY, TO_CHAR(k.IDCURRENCY), TO_CHAR(k.TOIDR), k.STARTDATE, k.ENDDATE
		FROM %s h JOIN %s k ON k.ID = h.IDKURS
		WHERE h.MASTERID = :1 AND k.TREATYYEAR = :2 ORDER BY h.URUTAN`, hub, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, false, err
	}
	// ⚠️ ID kurs TIDAK unik antartahun — kuncinya (ID, TREATYYEAR).
	rows, err := g.db.QueryContext(ctx, q, masterID, tahun)
	if err != nil {
		return nil, false, fmt.Errorf("repository: membaca kurs kontrak %s: %w", masterID, err)
	}
	defer func() { _ = rows.Close() }()
	var out []map[string]string
	for rows.Next() {
		var cur, idCur, toidr, mulai, akhir sql.NullString
		if err := rows.Scan(&cur, &idCur, &toidr, &mulai, &akhir); err != nil {
			return nil, false, fmt.Errorf("repository: membaca baris kurs kontrak %s: %w", masterID, err)
		}
		out = append(out, map[string]string{
			"Currency":    cur.String,
			"CurrencyID":  idCur.String,
			"Conversion":  toidr.String,
			"PeriodStart": mulai.String,
			"PeriodEnd":   akhir.String,
		})
	}
	return out, true, rows.Err()
}

// isiKurs memasang grid kurs satu sisi: kurs MILIK kontrak asalnya bila
// tercatat (455), selain itu kurs Treaty Year sisi itu.
func (g *Gudang) isiKurs(ctx context.Context, sisi models.SisiPenyesuaian) error {
	baris, punya, err := g.bacaKursKontrak(ctx, akarKontrak(sisi), sisi.Medan["TreatyYear"])
	if err != nil {
		return err
	}
	if !punya {
		if baris, err = g.bacaKursTahunan(ctx, sisi.Medan["TreatyYear"]); err != nil {
			return err
		}
	}
	if len(baris) > 0 {
		sisi.Larik[LarikKurs] = baris
	}
	return nil
}
