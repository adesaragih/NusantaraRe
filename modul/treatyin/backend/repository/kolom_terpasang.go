package repository

// Kolom tabel pendaratan yang SUDAH terpasang di skema.
//
// ⭐ Peta pendaratan boleh mendahului migrasinya. Migrasi dijalankan pihak
// lain (POOLDATA tidak boleh dimigrasi dari sesi pengembangan), jadi selama
// rentang antara berkas migrasi ditulis dan dipasang, kolom/tabel yang peta
// sebut belum ada. Tanpa penyaring ini setiap Save gagal `ORA-00904` dan
// setiap pembacaan dokumen gagal bersamanya.
//
// ⛔ Yang dilewati DILAPORKAN, tidak ditelan: `KunciBelumTerpasang`
// menyebut properti dokumen yang tidak dapat disimpan karena kolomnya belum
// ada, dan jalur Save meneruskannya ke layar.

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"nusantarare/inti/backend/db"
)

// Berapa lama hasil katalog dipercaya. Migrasi yang dipasang saat aplikasi
// berjalan terlihat paling lambat sesudah selang ini.
const umurKatalogKolom = time.Minute

type katalogKolom struct {
	kolom map[string]bool // nil = tabelnya belum ada
	waktu time.Time
}

var (
	kunciKatalog sync.Mutex
	katalog      = map[string]katalogKolom{}
)

// kolomTerpasang - kolom (huruf besar) tabel itu, atau nil bila tabelnya
// belum ada di skema.
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

// petaTerpasang - satu entri peta disaring ke kolom yang sudah ada; `ada`
// palsu bila tabelnya sendiri belum ada.
func (g *Gudang) petaTerpasang(ctx context.Context, p Pendaratan) (Pendaratan, bool, error) {
	kolom, err := g.kolomTerpasang(ctx, p.Tabel)
	if err != nil {
		return Pendaratan{}, false, err
	}
	if kolom == nil {
		return p, false, nil
	}
	saring := p
	saring.Kunci, saring.Kolom = nil, nil
	for i, k := range p.Kolom {
		if kolom[strings.ToUpper(k)] {
			saring.Kunci = append(saring.Kunci, p.Kunci[i])
			saring.Kolom = append(saring.Kolom, k)
		}
	}
	return saring, true, nil
}

// KunciBelumTerpasang - properti dokumen yang peta kenal tetapi TIDAK dapat
// disimpan karena kolom atau tabelnya belum terpasang (migrasi menunggu).
func (g *Gudang) KunciBelumTerpasang(ctx context.Context, doc map[string]any) ([]string, error) {
	var out []string
	sudah := map[string]bool{}
	tambah := func(k string) {
		if k != "" && !sudah[k] {
			sudah[k] = true
			out = append(out, k)
		}
	}
	elemen := elemenPerTabel(doc)
	for _, p := range PetaPendaratan {
		if p.Induk != "" {
			// ⭐ Tabel ANAK (mis. `T_TREATY_LIMIT_DETAIL`, migrasi `449`):
			// kolom yang belum ada dilaporkan per kunci yang SUNGGUH dikirim;
			// tabel yang belum ada dilaporkan dengan nama lariknya.
			el := elemen[p.Tabel]
			if len(el) == 0 {
				continue
			}
			kolom, err := g.kolomTerpasang(ctx, p.Tabel)
			if err != nil {
				return nil, err
			}
			if kolom == nil {
				tambah(p.KunciAnak)
				for _, n := range p.LarikGabung {
					tambah(n)
				}
				continue
			}
			for i, k := range p.Kunci {
				if kolom[strings.ToUpper(p.Kolom[i])] {
					continue
				}
				for _, e := range el {
					if _, dikirim := e[k]; dikirim {
						tambah(k)
						break
					}
				}
			}
			continue
		}
		kolom, err := g.kolomTerpasang(ctx, p.Tabel)
		if err != nil {
			return nil, err
		}
		switch {
		case p.Akar:
			for i, k := range p.Kunci {
				if _, dikirim := doc[strings.SplitN(k, ".", 2)[0]]; !dikirim {
					continue
				}
				if kolom == nil || !kolom[strings.ToUpper(p.Kolom[i])] {
					tambah(k)
				}
			}
		case kolom == nil:
			for _, n := range append([]string{p.Larik}, p.LarikGabung...) {
				if _, dikirim := doc[n]; dikirim && n != "" {
					tambah(n)
				}
			}
		}
	}
	return out, nil
}
