package repository

// Untuk apa berkas ini: SATU-SATUNYA pembacaan dokumen JSON di modul ini - master
// kontrak treaty untuk jalur NB NonProporsional / XOL.
//
// ⭐ `[keputusan work owner]` K8 (03-10-2026) - pengecualian SEMPIT atas P29:
// data master jalur XOL (`TreatyIn.Share()`, `Installment()`,
// `FacultativeShareList()`, `Limits()`, ...) tidak ada di view
// `TREATYINDETAILJOINEDM`, dan tabel master relasional modul `treatyin`
// (`KONTRAK`, `LAYER`, `BAGIAN`, `PEMULIHAN_LIMIT`, `TERMIN`, `POTONGAN`) masih nol
// baris. Maka `JSONDATA` `M_TREATY_IN` / `M_TREATY_IN_EDM` dibaca:
//
//   - BACA-SAJA (satu SELECT; nol INSERT/UPDATE, nol penulisan JSON di mana pun);
//   - HANYA medan master daftar K8 (`models.SkalarMasterXOL`,
//     `models.DaftarMasterXOL`) - selebihnya dibuang di sini, tidak pernah sampai
//     ke halaman;
//   - di SATU fungsi (`MasterXOLDariJSON`) di balik `services.PembacaMasterTreaty`,
//     supaya kelak diganti kontrak modul `treatyin` begitu tabel masternya terisi
//     (PERMINTAAN-TIM-INTI bagian E).
//
// `[penyimpangan sadar]` atas P29 - dicatat di tiket 01 (RALAT K8) dan spec RALAT
// AC 15/62. Treaty KELUAR (`BrowseTreatyOut`, `M_TREATY_OUT`) TIDAK dibaca (K8
// butir 4).

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbtreatyin/backend/models"
)

// Tabel master warisan yang dibaca (baca-saja, tidak dibuat ulang - diagram
// grilling "dibaca saja": M_TREATY_IN_EDM; M_TREATY_IN ditulis modul treaty in).
const (
	tabelMasterTreaty    = "M_TREATY_IN"
	tabelMasterTreatyEDM = "M_TREATY_IN_EDM"
)

var (
	// ErrMasterXOLTidakAda - nol baris master untuk nomor kontrak itu.
	ErrMasterXOLTidakAda = errors.New("repository: master treaty (M_TREATY_IN / M_TREATY_IN_EDM) kontrak ini tidak ditemukan")
	// ErrMasterXOLRusak - dokumen master tidak dapat diurai. Di Pega galatnya
	// hanya masuk log (`catch ... oLog.error`) dan halaman tetap kosong; di sini
	// ditampilkan (AC 36-38, spec §5.9).
	ErrMasterXOLRusak = errors.New("repository: dokumen master treaty tidak dapat diurai")
)

// MasterXOLDariJSON membaca master satu kontrak (`noOffer` =
// `PolicyTreatyIn.NoOffer` = `TREATYID` baris view) - padanan PERSIS dua RDB-List:
//
//	RDBList\BrowseTreatyIn (SetTreatyIn_Act 4, dipanggil TreatyRealizationCheckXOLList):
//	  select * from (
//	  select JSONDATA as ClassofBusiness from pooldata.M_TREATY_IN where ID={TreatyIn.ID}
//	  union all
//	  select JSONDATA as ClassofBusiness from pooldata.M_TREATY_IN_edm where ID={TreatyIn.ID}
//	  )
//	RDBList\BrowseTreatyInJoinEDM (InputPolicyTreatyInDetail_NonProp 5):
//	  select a.id, a.JSONDATA as CLASSOFBUSINESS from pooldata.m_treaty_in a
//	  where a.ID = {pyWorkPage.PolicyTreatyIn.NoOffer}
//	  union all
//	  select b.id, b.JSONDATA as CLASSOFBUSINESS from pooldata.m_treaty_in_edm b
//	  where b.ID = {pyWorkPage.PolicyTreatyIn.NoOffer}
//
// (`TreatyIn.ID` = `Param.ID` = `PolicyTreatyIn.NoOffer`, CheckXOLList langkah 3 -
// `pyPassCurrentParameterPage=true`, sehingga ID harfiah `3000004` di larik
// parameter Call tidak dipakai.) Setiap baris diurai berurutan dan kunci tingkat
// atasnya menimpa baris sebelumnya - Java langkah 6 (`adoptJSONObject` per baris).
func (g *Gudang) MasterXOLDariJSON(ctx context.Context, noOffer string) (models.MasterXOL, error) {
	a, err := g.nama(tabelMasterTreaty)
	if err != nil {
		return models.MasterXOL{}, err
	}
	b, err := g.nama(tabelMasterTreatyEDM)
	if err != nil {
		return models.MasterXOL{}, err
	}
	q := fmt.Sprintf(`SELECT JSONDATA FROM %s WHERE ID = :1 UNION ALL SELECT JSONDATA FROM %s WHERE ID = :2`, a, b)
	if err := db.PeriksaSQL(q); err != nil {
		return models.MasterXOL{}, err
	}
	rows, err := g.db.QueryContext(ctx, q, noOffer, noOffer)
	if err != nil {
		return models.MasterXOL{}, fmt.Errorf("repository: membaca master treaty: %w", err)
	}
	defer rows.Close()
	var isi []string
	for rows.Next() {
		var v sql.NullString
		if err := rows.Scan(&v); err != nil {
			return models.MasterXOL{}, fmt.Errorf("repository: membaca master treaty: %w", err)
		}
		isi = append(isi, v.String)
	}
	if err := rows.Err(); err != nil {
		return models.MasterXOL{}, fmt.Errorf("repository: membaca master treaty: %w", err)
	}
	if len(isi) == 0 {
		return models.MasterXOL{}, fmt.Errorf("%w: %s", ErrMasterXOLTidakAda, noOffer)
	}
	return uraiMasterXOL(isi)
}

// uraiMasterXOL mengurai dokumen berurutan dan menyaring medan K8.
func uraiMasterXOL(isi []string) (models.MasterXOL, error) {
	m := models.MasterXOL{Nilai: map[string]string{}, Daftar: map[string][]models.Baris{}}
	for n, s := range isi {
		dec := json.NewDecoder(strings.NewReader(s))
		dec.UseNumber() // angka tetap TEKS - nol float
		var obj map[string]any
		if err := dec.Decode(&obj); err != nil {
			return models.MasterXOL{}, fmt.Errorf("%w: baris %d: %w", ErrMasterXOLRusak, n+1, err)
		}
		for _, k := range models.SkalarMasterXOL {
			if v, ada := obj[k]; ada {
				if t, ok := teksSkalar(v); ok {
					m.Nilai[k] = t
				}
			}
		}
		for nama, skema := range models.DaftarMasterXOL {
			v, ada := obj[nama]
			if !ada {
				continue
			}
			for k := range m.Daftar { // daftar diganti utuh beserta anaknya
				if k == nama || strings.HasPrefix(k, nama+"(") {
					delete(m.Daftar, k)
				}
			}
			var baris []models.Baris
			for i, r := range barisObjek(v) {
				baris = append(baris, ambilAnggota(r, nama, skema.Anggota))
				for anak, anggota := range skema.Anak {
					var sub []models.Baris
					for _, rr := range barisObjek(r[anak]) {
						sub = append(sub, ambilAnggota(rr, nama+"."+anak, anggota))
					}
					if len(sub) > 0 {
						m.Daftar[models.JalurAnak(nama, i+1, anak)] = sub
					}
				}
			}
			m.Daftar[nama] = baris
		}
	}
	return m, nil
}

// ambilAnggota menyalin anggota skalar yang disebut skema saja.
func ambilAnggota(r map[string]any, jalur string, anggota []string) models.Baris {
	b := models.Baris{}
	for _, a := range anggota {
		v, ada := r[a]
		if !ada {
			continue
		}
		if t, ok := teksSkalar(v); ok {
			if models.TanggalMasterXOL[jalur+"."+a] {
				t = tanggalPega(t)
			}
			b[a] = t
		}
	}
	return b
}

// barisObjek - PageList JSON: larik objek; satu objek = satu baris.
func barisObjek(v any) []map[string]any {
	switch x := v.(type) {
	case []any:
		var out []map[string]any
		for _, e := range x {
			if o, ok := e.(map[string]any); ok {
				out = append(out, o)
			}
		}
		return out
	case map[string]any:
		return []map[string]any{x}
	}
	return nil
}

// teksSkalar - nilai skalar JSON sebagai teks properti Pega. Angka lewat
// desimal (bukan float) ke bentuk bertitik tanpa notasi ilmiah.
func teksSkalar(v any) (string, bool) {
	switch x := v.(type) {
	case nil:
		return "", true
	case string:
		return x, true
	case bool:
		if x {
			return "true", true
		}
		return "false", true
	case json.Number:
		d, err := utils.ParseDecimal(x.String())
		if err != nil {
			return x.String(), true
		}
		return d.Text('f'), true
	}
	return "", false
}

// tanggalPega menormalkan tanggal dokumen ke "2006-01-02" (diagram F21: dua
// format masuk - `YYYYMMDD` dan cap waktu Pega bersufiks GMT). Cap waktu GMT
// dibaca di zona Asia/Jakarta (`[tafsiran]` zona bisnis). Bentuk lain apa adanya -
// penyimpanan menolaknya dengan pesan jalur medannya.
func tanggalPega(s string) string {
	t := strings.TrimSpace(s)
	if len(t) == 8 && utils.AngkaSaja(t) {
		if d, err := time.Parse("20060102", t); err == nil {
			return d.Format("2006-01-02")
		}
	}
	if strings.HasSuffix(t, " GMT") {
		for _, l := range []string{"20060102T150405.000 MST", "20060102T150405 MST"} {
			if d, err := time.Parse(l, t); err == nil {
				return d.In(zonaWIB()).Format("2006-01-02")
			}
		}
	}
	if len(t) >= 10 {
		if d, err := time.Parse("2006-01-02", t[:10]); err == nil {
			return d.Format("2006-01-02")
		}
	}
	return s
}

func zonaWIB() *time.Location {
	if l, err := time.LoadLocation("Asia/Jakarta"); err == nil {
		return l
	}
	return time.FixedZone("WIB", 7*3600)
}
