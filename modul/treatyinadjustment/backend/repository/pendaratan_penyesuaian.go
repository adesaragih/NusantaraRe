package repository

// Layar Penyesuaian dibaca dari TABEL PENDARATAN — bukan lagi dari
// `M_TREATY_IN_EDM.JSONDATA`.
//
// ---------------------------------------------------------------------
// ⛔ KEPUTUSAN PEMILIK PROSES, 6 Oktober 2026
// ---------------------------------------------------------------------
//
//	"jangan ada dri jsondata lagi, begitu juga treaty in adjustment
//	 kemudian hubungkan ke backend agar table di aplikasi tidak membaca
//	 dari json data lagi"
//
// ---------------------------------------------------------------------
// ⭐ TABEL YANG SAMA DENGAN TREATY IN, DAN ITU BUKAN KEBETULAN
// ---------------------------------------------------------------------
// `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx` menandai setiap kotaknya
// `BERSAMA -> Prop, Non Prop, EDM Prop, EDM Non Prop`. Dokumen Adjustment
// berbentuk sama dengan dokumen Treaty In, jadi ia mendarat di tabel yang
// sama — dan modul ini MEMBACA tabel milik modul Treaty In, tanpa menulis
// apa pun ke sana.
//
// ---------------------------------------------------------------------
// ⭐ DUA SISI, SATU TABEL: pembedanya `MASTERID`
// ---------------------------------------------------------------------
// Dokumen Adjustment membawa DUA halaman — akar (sisi `New`) dan `OLDDATA`
// (sisi `Old`). Keduanya mendarat di tabel yang sama; yang membedakan
// akhiran `MASTERID`, dan tetapannya hidup di SATU tempat
// (`treatyin/repository.AkhiranSisiLama`) supaya pemuat dan pembaca tidak
// dapat berselisih.
//
// ⚠️ `OLDDATA` DI DALAM `OLDDATA` tetap TIDAK ditelusuri — aturan yang sama
// seperti sebelum penukaran ini. Section hanya mengikat satu tingkat
// (`TreatyIn.OLDDATA.*`); tingkat kedua tidak tampil di layar mana pun.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyinadjustment/backend/models"
)

// BacaPenyesuaianPendaratan menyusun kedua sisi satu penyesuaian dari tabel
// pendaratan.
//
// ⛔ Dikemudikan `treatyin.PetaPendaratan` — peta yang SAMA yang memuatnya.
// Dua daftar yang memerikan hal yang sama akan berselisih suatu hari, dan
// yang basi selalu yang kedua; di sini daftar kedua itu tidak ada.
func (g *Gudang) BacaPenyesuaianPendaratan(ctx context.Context, id string) (models.Penyesuaian, error) {
	p := models.Penyesuaian{ID: id, Baru: sisiKosong(), Lama: sisiKosong()}
	baru, err := g.bacaSisi(ctx, id)
	if err != nil {
		return models.Penyesuaian{}, err
	}
	lama, err := g.bacaSisi(ctx, id+AkhiranSisiLama)
	if err != nil {
		return models.Penyesuaian{}, err
	}
	// ⭐ Halaman `ActualValue` (MASTERID `id + AkhiranSisiAktual`) kembali ke
	// sisi New sebagai kunci BERTITIK — bentuk yang kerangka dan rumus layar
	// ikat (`ActualValue.EGNPI`, `ActualValue.TotalEgnpiAmount`, …).
	aktual, err := g.bacaSisi(ctx, id+AkhiranSisiAktual)
	if err != nil {
		return models.Penyesuaian{}, err
	}
	gabungAktual(baru, aktual)
	// ⭐ Halaman `ValueBeforeProrate` (MASTERID `id +
	// AkhiranSisiSebelumProrata`) — pola yang sama, kunci bertitik.
	prorata, err := g.bacaSisi(ctx, id+AkhiranSisiSebelumProrata)
	if err != nil {
		return models.Penyesuaian{}, err
	}
	gabungSebelumProrata(baru, prorata)
	// ⭐ Kunci kepala yang TIDAK ada di pendaratan dilengkapi dari kolom
	// `TREATY_IN_EDM` — cara yang sama dengan `BacaDokumenMaster` (picker
	// Add). Tanpa itu penyesuaian yang belum (atau tidak lagi) mendarat tidak
	// punya `ProportionType`, panel Old/New tidak pernah terbuka, dan layar
	// Edit tidak dapat disunting sama sekali (laporan 8 Oktober 2026,
	// `1001540/R01`). Hanya sisi New: `OLDDATA` adalah salinan saat
	// penyesuaian dibuat, dan kepala master HARI INI bukan salinan itu.
	kepala, _, err := g.bacaKepala(ctx, TabelWarisanPenyesuaian,
		append(append([][2]string{}, kolomKepalaMaster...), kolomKepalaEDM...), id)
	if err != nil {
		return models.Penyesuaian{}, err
	}
	for k, v := range kepala {
		if _, sudah := baru.Medan[k]; !sudah {
			baru.Medan[k] = v
		}
	}
	// ⭐ Rate of Exchange — `TREATYEXCHANGEYEARLY` per Treaty Year sisi itu
	// sendiri, bukan larik dokumen. Lihat `kurs_tahunan.go`.
	for _, sisi := range []models.SisiPenyesuaian{baru, lama} {
		if err := g.isiKurs(ctx, sisi); err != nil {
			return models.Penyesuaian{}, err
		}
	}
	p.Baru, p.Lama = baru, lama
	p.Terdarat = baru.Terdarat
	// `OLDID` dibaca dari sisi yang punya — sisi `New` lebih dahulu.
	if v, ada := baru.Medan["OLDID"]; ada && v != "" {
		p.IDAsal = v
	} else if v, ada := lama.Medan["OLDID"]; ada {
		p.IDAsal = v
	}
	return p, nil
}

// bacaSisi membaca SATU sisi — satu `MASTERID` — menjadi medan dan larik.
func (g *Gudang) bacaSisi(ctx context.Context, masterID string) (models.SisiPenyesuaian, error) {
	sisi := sisiKosong()
	for _, asli := range petaPendaratanPenyesuaian {
		// ⭐ Hanya kolom yang SUDAH terpasang (`kolom_terpasang.go`) — peta
		// boleh mendahului migrasinya; tabel yang belum ada = nol baris.
		pd, ada, err := g.larikTerpasang(ctx, asli)
		if err != nil {
			return sisi, err
		}
		if !ada {
			continue
		}
		switch {
		case pd.Induk != "":
			// Tabel ANAK — dirangkai ke `Pohon` oleh `bacaPohon`, tidak ke
			// larik datar (nama larik anaknya bukan larik akar).
			continue
		case pd.Akar:
			medan, err := g.bacaMedanAkar(ctx, pd, masterID)
			if err != nil {
				return sisi, err
			}
			// Baris akar ber-MASTERID ini ada → sisi ini TERDARAT.
			sisi.Terdarat = sisi.Terdarat || len(medan) > 0
			for k, v := range medan {
				sisi.Medan[k] = v
			}
		case pd.Larik != "":
			baris, err := g.bacaLarik(ctx, pd, masterID)
			if err != nil {
				return sisi, err
			}
			if len(baris) > 0 {
				sisi.Larik[pd.Larik] = baris
			}
		case len(pd.Gabung) > 0:
			// ⛔ Beberapa larik akar DI SATU tabel — dipisah kembali menurut
			// kolom `JENIS`, yang pemuat isi dengan nama larik asalnya.
			perJenis, err := g.bacaLarikGabung(ctx, pd, masterID)
			if err != nil {
				return sisi, err
			}
			for nama, baris := range perJenis {
				sisi.Larik[nama] = baris
			}
		}
	}
	pohon, _, err := g.bacaPohon(ctx, masterID)
	if err != nil {
		return sisi, err
	}
	sisi.Pohon = pohon
	return sisi, nil
}

// BacaPohonUntukUji membuka perangkai pohon beserta cacah baris yatimnya —
// untuk uji db saja.
func (g *Gudang) BacaPohonUntukUji(ctx context.Context, masterID string) (map[string][]map[string]any, int, error) {
	return g.bacaPohon(ctx, masterID)
}

func (g *Gudang) bacaMedanAkar(ctx context.Context, pd larikPendaratan, masterID string) (map[string]string, error) {
	nama, err := g.db.Qualify(pd.Tabel)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf("SELECT %s FROM %s WHERE MASTERID = :1",
		strings.Join(pd.Kolom, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	sel := make([]sql.NullString, len(pd.Kolom))
	tuju := make([]any, len(pd.Kolom))
	for i := range sel {
		tuju[i] = &sel[i]
	}
	err = g.db.QueryRowContext(ctx, q, masterID).Scan(tuju...)
	if err == sql.ErrNoRows {
		// Sisi yang tidak ada bukan galat: penyesuaian tanpa `OLDDATA`
		// memang hanya punya sisi `New`.
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s %s: %w", pd.Tabel, masterID, err)
	}
	out := map[string]string{}
	for i, kunci := range pd.Kunci {
		// ⛔ `NULL` dilewati: kunci yang TIDAK ADA di dokumen berbeda arti
		// dari kunci yang ada bernilai kosong, dan layar menyatakan bedanya.
		if sel[i].Valid {
			out[kunci] = sel[i].String
		}
	}
	return out, nil
}

func (g *Gudang) bacaLarik(ctx context.Context, pd larikPendaratan, masterID string) ([]map[string]string, error) {
	baris, _, err := g.bacaBarisPendaratan(ctx, pd, masterID, false)
	return baris, err
}

func (g *Gudang) bacaLarikGabung(ctx context.Context, pd larikPendaratan, masterID string) (map[string][]map[string]string, error) {
	baris, jenis, err := g.bacaBarisPendaratan(ctx, pd, masterID, true)
	if err != nil {
		return nil, err
	}
	out := map[string][]map[string]string{}
	for i, b := range baris {
		out[jenis[i]] = append(out[jenis[i]], b)
	}
	return out, nil
}

// bacaBarisPendaratan membaca satu tabel larik, urut `URUTAN`.
func (g *Gudang) bacaBarisPendaratan(ctx context.Context, pd larikPendaratan,
	masterID string, pakaiJenis bool) ([]map[string]string, []string, error) {
	nama, err := g.db.Qualify(pd.Tabel)
	if err != nil {
		return nil, nil, err
	}
	kolom := append([]string{}, pd.Kolom...)
	if pakaiJenis {
		kolom = append(kolom, "JENIS")
	}
	q := fmt.Sprintf("SELECT %s FROM %s WHERE MASTERID = :1 ORDER BY URUTAN",
		strings.Join(kolom, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, masterID)
	if err != nil {
		return nil, nil, fmt.Errorf("repository: membaca %s %s: %w", pd.Tabel, masterID, err)
	}
	defer func() { _ = rows.Close() }()

	var out []map[string]string
	var jenis []string
	for rows.Next() {
		sel := make([]sql.NullString, len(kolom))
		tuju := make([]any, len(kolom))
		for i := range sel {
			tuju[i] = &sel[i]
		}
		if err := rows.Scan(tuju...); err != nil {
			return nil, nil, fmt.Errorf("repository: membaca baris %s: %w", pd.Tabel, err)
		}
		el := map[string]string{}
		for i, kunci := range pd.Kunci {
			if sel[i].Valid {
				el[kunci] = sel[i].String
			}
		}
		out = append(out, el)
		if pakaiJenis {
			jenis = append(jenis, sel[len(kolom)-1].String)
		}
	}
	return out, jenis, rows.Err()
}

// AkhiranSisiLama membedakan sisi `Old` dari sisi `New` di dalam `MASTERID`.
//
// ⛔ NILAINYA WAJIB SAMA dengan `treatyin/repository.AkhiranSisiLama` — yang
// menulisnya pemuat milik modul itu, yang membacanya berkas ini. Kesamaannya
// dibuktikan `uji/lintasmodul/peta_pendaratan_test.go`; dua tempat yang
// menulis akhiran berbeda terbaca sebagai "sisi Old kosong", bukan sebagai
// galat.
const AkhiranSisiLama = "#LAMA"

// AkhiranSisiAktual - akhiran `MASTERID` halaman `ActualValue` (cabang
// Adjust Premium). Keputusan pemakai 8 Oktober 2026: `ActualValue` mendarat
// di tabel `T_TREATY_*` yang sama, ber-MASTERID sendiri.
//
// ⛔ NILAINYA WAJIB SAMA dengan `treatyin/repository.AkhiranSisiAktual`
// (penulisnya) — `uji/lintasmodul/peta_pendaratan_test.go` mengadunya.
const AkhiranSisiAktual = "#AKTUAL"

// AkhiranSisiSebelumProrata - akhiran `MASTERID` halaman `ValueBeforeProrate`
// (salinan `ValueDifference` sebelum pro-rate, `TreatyEDMProRateCalculation`).
//
// ⛔ NILAINYA WAJIB SAMA dengan `treatyin/repository.
// AkhiranSisiSebelumProrata` (penulisnya) —
// `uji/lintasmodul/peta_pendaratan_test.go` mengadunya.
const AkhiranSisiSebelumProrata = "#PRORATA"
