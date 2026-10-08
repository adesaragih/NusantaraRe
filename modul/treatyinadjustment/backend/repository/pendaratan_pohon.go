package repository

// POHON satu sisi penyesuaian — larik bersarang dari tabel ANAK pendaratan.
//
// ---------------------------------------------------------------------
// ⭐ MENGAPA DUA BENTUK (larik datar DAN pohon)
// ---------------------------------------------------------------------
// Grid kerangka membaca baris DATAR (`Larik`). Rumus Limits/Share/Installment
// — rute `/hitung/*` modul Treaty In, yang Activity-nya identik di kedua
// korpus — membaca SIMPUL bersarang: `Limits[].TreatyGroupList[]`,
// `.Detail[].COBList[]`, `.MDPList[]`, `Share[].DeductionList[]`,
// `.GrossPremiumList[]` … Pohon ini memberinya bentuk itu, tanpa mengubah
// larik datar yang grid baca.
//
// ⛔ Dikemudikan peta yang SAMA (`petaPendaratanPenyesuaian`), dan urutan
// peta menaruh setiap induk sebelum anaknya — urutan yang sama dengan
// pemuatnya. Anak yang induknya tidak ditemukan DIHITUNG, bukan ditelan.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// punyaAnak - tabel yang menjadi `Induk` entri lain.
var punyaAnak = func() map[string]bool {
	m := map[string]bool{}
	for _, pd := range petaPendaratanPenyesuaian {
		if pd.Induk != "" {
			m[pd.Induk] = true
		}
	}
	return m
}()

// larikAnak - nama larik yang dimiliki simpul satu tabel induk.
var larikAnak = func() map[string][]string {
	m := map[string][]string{}
	for _, pd := range petaPendaratanPenyesuaian {
		if pd.Induk == "" {
			continue
		}
		if pd.KunciAnak != "" {
			m[pd.Induk] = append(m[pd.Induk], pd.KunciAnak)
		}
		m[pd.Induk] = append(m[pd.Induk], pd.Gabung...)
	}
	return m
}()

type simpulPohon = map[string]any

// bacaPohon merangkai pohon satu sisi: larik AKAR yang punya anak → simpul.
//
// Kembalian kedua: cacah baris anak yang induknya tidak ditemukan — nol bila
// data utuh, dan dijaga uji db.
func (g *Gudang) bacaPohon(ctx context.Context, masterID string) (map[string][]simpulPohon, int, error) {
	puncak := map[string][]simpulPohon{}
	perTabel := map[string]map[int64]simpulPohon{}
	yatim := 0
	for _, asli := range petaPendaratanPenyesuaian {
		if asli.Akar || (asli.Induk == "" && !punyaAnak[asli.Tabel]) {
			continue
		}
		// Hanya kolom yang SUDAH terpasang; tabel yang belum ada = nol baris.
		pd, ada, err := g.larikTerpasang(ctx, asli)
		if err != nil {
			return nil, 0, err
		}
		if !ada {
			continue
		}
		baris, err := g.bacaBarisPohon(ctx, pd, masterID)
		if err != nil {
			return nil, 0, err
		}
		for _, b := range baris {
			s := b.nilai
			if punyaAnak[pd.Tabel] {
				// ⛔ Larik anak KOSONG, bukan tidak ada: rumus membaca
				// `.length` dan menapakinya tanpa memeriksa keberadaannya.
				for _, n := range larikAnak[pd.Tabel] {
					s[n] = []simpulPohon{}
				}
				if perTabel[pd.Tabel] == nil {
					perTabel[pd.Tabel] = map[int64]simpulPohon{}
				}
				perTabel[pd.Tabel][b.id] = s
			}
			if pd.Induk == "" {
				puncak[pd.Larik] = append(puncak[pd.Larik], s)
				continue
			}
			induk := perTabel[pd.Induk][b.induk]
			if induk == nil {
				yatim++
				continue
			}
			nama := pd.KunciAnak
			if nama == "" {
				nama = b.jenis
			}
			daftar, _ := induk[nama].([]simpulPohon)
			induk[nama] = append(daftar, s)
		}
	}
	return puncak, yatim, nil
}

type barisPohon struct {
	id, induk int64
	jenis     string
	nilai     simpulPohon
}

// bacaBarisPohon membaca satu tabel beserta `ID`, `IDINDUK` (tabel anak) dan
// `JENIS` (tabel gabungan), urut `URUTAN` — urutan larik di dokumen. URUTAN
// dihitung per induk; menempelkan baris berurutan ke induknya masing-masing
// mempertahankan urutan itu.
func (g *Gudang) bacaBarisPohon(ctx context.Context, pd larikPendaratan, masterID string) ([]barisPohon, error) {
	nama, err := g.db.Qualify(pd.Tabel)
	if err != nil {
		return nil, err
	}
	pilih := []string{"ID"}
	if pd.Induk != "" {
		pilih = append(pilih, "IDINDUK")
	} else {
		pilih = append(pilih, "0")
	}
	if pd.Induk != "" && len(pd.Gabung) > 0 {
		pilih = append(pilih, "JENIS")
	} else {
		pilih = append(pilih, "NULL")
	}
	q := fmt.Sprintf("SELECT %s, %s FROM %s WHERE MASTERID = :1 ORDER BY URUTAN",
		strings.Join(pilih, ", "), strings.Join(pd.Kolom, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, masterID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s %s: %w", pd.Tabel, masterID, err)
	}
	defer func() { _ = rows.Close() }()

	var out []barisPohon
	for rows.Next() {
		var id, induk sql.NullInt64
		var jenis sql.NullString
		sel := make([]sql.NullString, len(pd.Kolom))
		tuju := []any{&id, &induk, &jenis}
		for i := range sel {
			tuju = append(tuju, &sel[i])
		}
		if err := rows.Scan(tuju...); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", pd.Tabel, err)
		}
		s := simpulPohon{}
		for i, kunci := range pd.Kunci {
			// ⛔ `NULL` dilewati, sama dengan larik datar: kunci yang tidak ada
			// berbeda arti dari kunci yang ada bernilai kosong.
			if sel[i].Valid {
				s[kunci] = sel[i].String
			}
		}
		out = append(out, barisPohon{id: id.Int64, induk: induk.Int64, jenis: jenis.String, nilai: s})
	}
	return out, rows.Err()
}
