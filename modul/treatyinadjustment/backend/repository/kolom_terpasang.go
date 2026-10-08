package repository

// Kolom tabel pendaratan yang SUDAH terpasang — padanan
// `treatyin/repository/kolom_terpasang.go` (modul tidak saling impor).
//
// ⭐ Salinan peta ini wajib sama dengan peta Treaty In (`uji/lintasmodul`),
// dan peta itu boleh mendahului migrasinya (mis. `448`). Tanpa penyaring ini
// pembukaan penyesuaian gagal `ORA-00904` sampai migrasinya dipasang.

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"nusantarare/inti/backend/db"
)

const umurKatalogKolom = time.Minute

type katalogKolom struct {
	kolom map[string]bool // nil = tabelnya belum ada
	waktu time.Time
}

var (
	kunciKatalog sync.Mutex
	katalog      = map[string]katalogKolom{}
)

func (g *Gudang) kolomTerpasang(ctx context.Context, tabel string) (map[string]bool, error) {
	nama, err := g.db.Qualify(tabel)
	if err != nil {
		return nil, err
	}
	kunciKatalog.Lock()
	k, ada := katalog[nama]
	kunciKatalog.Unlock()
	if ada && time.Since(k.waktu) < umurKatalogKolom {
		return k.kolom, nil
	}
	pemilik, nTabel := "", strings.ToUpper(nama)
	if i := strings.IndexByte(nama, '.'); i > 0 {
		pemilik, nTabel = strings.ToUpper(nama[:i]), strings.ToUpper(nama[i+1:])
	}
	q := "SELECT COLUMN_NAME FROM USER_TAB_COLUMNS WHERE TABLE_NAME = :1"
	args := []any{nTabel}
	if pemilik != "" {
		q = "SELECT COLUMN_NAME FROM ALL_TAB_COLUMNS WHERE OWNER = :1 AND TABLE_NAME = :2"
		args = []any{pemilik, nTabel}
	}
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca katalog kolom %s: %w", tabel, err)
	}
	defer func() { _ = rows.Close() }()
	var kolom map[string]bool
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, fmt.Errorf("repository: membaca katalog kolom %s: %w", tabel, err)
		}
		if kolom == nil {
			kolom = map[string]bool{}
		}
		kolom[strings.ToUpper(c)] = true
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	kunciKatalog.Lock()
	katalog[nama] = katalogKolom{kolom: kolom, waktu: time.Now()}
	kunciKatalog.Unlock()
	return kolom, nil
}

// larikTerpasang - entri peta disaring ke kolom yang sudah ada; `ada` palsu
// bila tabelnya sendiri belum ada.
func (g *Gudang) larikTerpasang(ctx context.Context, pd larikPendaratan) (larikPendaratan, bool, error) {
	kolom, err := g.kolomTerpasang(ctx, pd.Tabel)
	if err != nil {
		return larikPendaratan{}, false, err
	}
	if kolom == nil {
		return pd, false, nil
	}
	saring := pd
	saring.Kunci, saring.Kolom = nil, nil
	for i, k := range pd.Kolom {
		if kolom[strings.ToUpper(k)] {
			saring.Kunci = append(saring.Kunci, pd.Kunci[i])
			saring.Kolom = append(saring.Kolom, k)
		}
	}
	return saring, true, nil
}
