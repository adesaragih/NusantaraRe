package repository

// Pembaca `T_TREATY_HAZARD_LIMIT` — isi tab **Event Limits** (Non-Prop).
//
// ---------------------------------------------------------------------
// ⛔ CELAH YANG DITUTUP 8 Oktober 2026
// ---------------------------------------------------------------------
// Tabelnya berdiri sejak migrasi `446`, peta pendaratan mengenalnya, dan Save
// MENULISNYA. Yang tidak pernah ada: pembacanya. Tab Event Limits karena itu
// selalu kosong — kedelapan medannya menampilkan `—` walau kontraknya punya
// nilai.
//
// ⚠️ Komentar lama di `TabEventLimits.tsx` berbunyi *"NILAI AKARNYA BELUM
// DAPAT DIBACA … tidak ada kolom pendaratan yang memuatnya"*. Itu BENAR saat
// ditulis 6 Oktober dan sudah KEDALUWARSA sejak `446` — dan keterangan
// kedaluwarsa yang terdengar pasti adalah sebab cacat ini bertahan dua hari.
//
// ⭐ Bentuknya sengaja sama dengan `BacaRevisiPendaratan`: kolom yang belum
// terpasang disaring, baris yang tidak ada bukan galat, dan `NULL` DILEWATI
// alih-alih menjadi teks kosong — itulah yang membedakan "tidak punya" dari
// "belum diisi".

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// TabelBatasBahaya - tabel akar kedua, di luar `T_TREATY_REVISION`.
const TabelBatasBahaya = "T_TREATY_HAZARD_LIMIT"

// kolomBatasBahaya - kolom ↔ ejaan properti dokumen.
//
// ⛔ Urutan dan ejaannya mengikuti `PetaPendaratan`; `TestBatasBahayaIkutPeta`
// mengadu keduanya supaya pembaca dan penulis tidak berselisih diam-diam.
var kolomBatasBahaya = [][2]string{
	{"RSMDLIMIT", "RSMDLimit"},
	{"CURRENCYRSMD", "CurrencyRSMD"},
	{"EARTHQUAKE", "Earthquake"},
	{"CURRENCYEARTHQUAKE", "CurrencyEarthquake"},
	{"FLOODJAB", "FloodJab"},
	{"CURRENCYFLOODJAB", "CurrencyFloodJab"},
	{"FLOODNATION", "FloodNation"},
	{"CURRENCYFLOODNAT", "CurrencyFloodNat"},
	// ⭐ Kedua batas Co-Ins ikut di tabel yang sama (tab Co-Ins Scale).
	{"MAXCOGROUP", "MaxCoGroup"},
	{"MAXCONONGROUP", "MaxCoNonGroup"},
}

// BacaBatasBahaya membaca medan akar `T_TREATY_HAZARD_LIMIT` satu kontrak.
//
// ⚠️ Baris yang TIDAK ADA bukan galat — kontrak yang belum didaratkan tetap
// harus terbuka.
func (g *Gudang) BacaBatasBahaya(ctx context.Context, masterID string) (map[string]string, error) {
	out := map[string]string{}
	terpasang, err := g.kolomTerpasang(ctx, TabelBatasBahaya)
	if err != nil {
		return nil, err
	}
	if terpasang == nil {
		return out, nil // tabelnya belum ada — bukan galat
	}
	var pakai [][2]string
	for _, p := range kolomBatasBahaya {
		if terpasang[p[0]] {
			pakai = append(pakai, p)
		}
	}
	if len(pakai) == 0 {
		return out, nil
	}
	nama, err := g.db.Qualify(TabelBatasBahaya)
	if err != nil {
		return nil, err
	}
	kolom := make([]string, len(pakai))
	for i, p := range pakai {
		kolom[i] = p[0]
	}
	q := fmt.Sprintf("SELECT %s FROM %s WHERE MASTERID = :1", strings.Join(kolom, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	sel := make([]sql.NullString, len(kolom))
	tuju := make([]any, len(kolom))
	for i := range sel {
		tuju[i] = &sel[i]
	}
	err = g.db.QueryRowContext(ctx, q, masterID).Scan(tuju...)
	if err == sql.ErrNoRows {
		return out, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s kontrak %s: %w", TabelBatasBahaya, masterID, err)
	}
	for i, p := range pakai {
		// ⛔ `NULL` dilewati, bukan dimasukkan sebagai teks kosong.
		if sel[i].Valid {
			out[p[1]] = sel[i].String
		}
	}
	return out, nil
}
