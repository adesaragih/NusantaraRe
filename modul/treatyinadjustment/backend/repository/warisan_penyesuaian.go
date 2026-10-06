// ⛔⛔ `BacaPenyesuaianWarisan` DI BERKAS INI TIDAK LAGI DIPAKAI JALUR BACA,
// sejak 6 Oktober 2026.
//
// Pemilik proses melarang keras menarik nilai dari `JSONDATA`. Layar
// Penyesuaian kini dibaca dari TABEL PENDARATAN lewat
// `pendaratan_penyesuaian.go`; nol kode produksi memanggil pengurai di
// berkas ini.
//
// ⚠️ DITINGGALKAN satu ronde dengan sengaja: tabel pendaratannya masih NOL
// BARIS sampai pemuat dijalankan atas `M_TREATY_IN_EDM`, sehingga jalur
// barunya belum pernah terbukti di layar. Membuang jalur lama sebelum
// penggantinya terbukti berarti membuang satu-satunya pembanding ketika
// hasilnya berselisih.
//
// ⭐ `DaftarPenyesuaianWarisan` TETAP DIPAKAI dan TIDAK terkena larangan:
// ia membaca KOLOM `TREATY_IN_EDM`, bukan `JSONDATA`.

package repository

// Baca penyesuaian WARISAN — `TREATY_IN_EDM` berpasangan `M_TREATY_IN_EDM`.
//
// ⛔ BACA SAJA, kedua tabel. Nol `INSERT`/`UPDATE`/`DELETE`/`MERGE`, nol DDL,
// nol penyebutan di migrasi modul ini.
//
// =====================================================================
// KENAPA BUKAN `VERSI_KONTRAK` (migrasi 440)
// =====================================================================
//   Migrasi 440 memasang `ID_VERSI_KONTRAK_DASAR` — jalur versi dasar model
//   baru. Diukur 5 Oktober 2026: `KONTRAK` dan `VERSI_KONTRAK` NOL baris.
//   Penyesuaian yang benar-benar ada hidup di tabel warisan ini: 280 baris,
//   dan ke-280 dokumennya membawa `OLDDATA` — nilai versi sebelumnya yang
//   panel Old Data tampilkan. Jalur 440 tetap milik model baru; jalur ini
//   milik data yang sudah ada.
//
// =====================================================================
// KENAPA DAFTAR-IZIN KUNCI, BUKAN SELURUH DOKUMEN
// =====================================================================
//   Satu dokumen terbesar 2.475.844 bita. Yang layar pakai beberapa puluh
//   kunci; sisanya tidak dikirim. Daftar-izin ini SELEKSI, bukan tafsir:
//   nilai yang terpilih dibawa apa adanya, termasuk ejaannya.

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyinadjustment/backend/models"
)

// TabelWarisanPenyesuaian - kepala penyesuaian sistem lama.
const TabelWarisanPenyesuaian = "TREATY_IN_EDM"

// TabelWarisanPenyesuaianJSON - dokumen halaman `TreatyIn` penyesuaian itu.
const TabelWarisanPenyesuaianJSON = "M_TREATY_IN_EDM"

// ErrPenyesuaianTidakAda - pengenalnya tidak menunjuk baris mana pun.
var ErrPenyesuaianTidakAda = errors.New("penyesuaian warisan tidak ada")

// ErrJSONPenyesuaianRusak - dokumennya ada tetapi tidak dapat diurai.
var ErrJSONPenyesuaianRusak = errors.New("dokumen JSON penyesuaian tidak dapat diurai")

// MedanDibaca - kunci SKALAR yang layar baca dari tiap halaman.
//
// ⛔ Kunci bertitik (`ValueDifference.RNMShare`) adalah medan halaman
// tertanam satu tingkat — persis cara Section mengikatnya
// (`TreatyIn.ValueDifference.RNMShare`, `TreatyInTabsNPValueDifferenceProRate`
// @137035).
//
// ⚠️ `frontend/layar.test.ts` membaca daftar ini dan menolak kunci yang
// layar pakai tetapi tidak ada di sini — kunci yang lupa ditambahkan akan
// tampil sebagai "tidak ada di sistem lama", padahal ada.
var MedanDibaca = []string{
	"ID", "OLDID", "ProportionType", "EDMState", "EDMMaterialType", "StatusAkseptasi",
	"TreatyContractName", "ContractRefNo", "TeritorialScope", "Bordeaux", "BordereauxNote",
	"Commencement", "Termination", "TreatyYear", "AccountingMode", "AccountingModeNonProp",
	"Ceding", "LeadingReinsSource", "TreatyLeader", "EDMEffective",
	"IsProRate", "ProRateDays", "ProRateTotalDays", "ProRatePercent", "IsMultipleRetro", "IsEditData",
	"Exclusions", "ExclusionsP", "SpecialConditions", "SpecialConditionsP",
	"TotalEgnpiAmount", "TotalEgnpiProportion", "TotalLimitsROL",
	"FacultativeShare", "FacultativeShareBrokerage",
	"ValueDifference.RNMShare", "ValueDifference.BrokeragePercent",
}

// LarikDibaca - larik yang layar baca dari tiap halaman. Elemennya dibawa
// sebagai peta medan SKALAR; medan bersarang di dalam elemen tidak dibawa.
var LarikDibaca = []string{
	"CurrencyList", "CommentList",
	"Retention", "TotalRetentionAmountNP",
	"EGNPI", "TotalEgnpiAmountNP",
	"Limits", "LimitSummaryList",
	"TotalLimitIOONP", "TotalLimitDeductblNP", "TotalLimitPremiEarnNP", "TotalLimitMDPNP",
	"ReportingPeriodList", "Portfolio", "AccumulationList",
}

// DaftarPenyesuaianWarisan membaca seluruh kepala penyesuaian.
//
// ⛔ `ORDER BY ID`. Pengenalnya berbentuk `<asal>/R<nn>`, jadi urutan teks
// menaruh revisi satu kontrak bersebelahan dan berurut — dan tanpa klausa
// ini Oracle bebas mengembalikan urutan apa pun.
func (g *Gudang) DaftarPenyesuaianWarisan(ctx context.Context) ([]models.BarisPenyesuaian, error) {
	nama, err := g.db.Qualify(TabelWarisanPenyesuaian)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT ID, OLDID, EDMSTATE, EDMMATERIALTYPE, TREATYCONTRACTNAME,
		PROPORTIONTYPE, LEADINGREINSSOURCE, CEDING, COMMENCEMENT, TERMINATION,
		POSITIONUSERNAME, STATUSAKSEPTASI
		FROM %s ORDER BY ID`, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", TabelWarisanPenyesuaian, err)
	}
	defer func() { _ = rows.Close() }()

	out := []models.BarisPenyesuaian{}
	for rows.Next() {
		var k [12]sql.NullString
		if err := rows.Scan(&k[0], &k[1], &k[2], &k[3], &k[4], &k[5], &k[6], &k[7],
			&k[8], &k[9], &k[10], &k[11]); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", TabelWarisanPenyesuaian, err)
		}
		// ⛔ Apa adanya — nol tafsir di repository.
		out = append(out, models.BarisPenyesuaian{
			ID: k[0].String, IDAsal: k[1].String,
			JenisPenyesuaian: k[2].String, JenisMaterial: k[3].String,
			NamaKontrak: k[4].String, SifatProporsi: k[5].String,
			AsalBisnis: k[6].String, Cedant: k[7].String,
			TanggalMulai: k[8].String, TanggalBerakhir: k[9].String,
			Posisi: k[10].String, StatusAkseptasi: k[11].String,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", TabelWarisanPenyesuaian, err)
	}
	return out, nil
}

// BacaPenyesuaianWarisan membaca satu penyesuaian beserta dokumennya.
//
// ⚠️ `LEFT JOIN`: diukur 280 lawan 280 dan nol dokumen NULL, tetapi
// penyesuaian yang kehilangan dokumennya harus tetap TERBUKA dengan kedua
// sisinya kosong — bukan galat yang menyembunyikan barisnya.
func (g *Gudang) BacaPenyesuaianWarisan(ctx context.Context, id string) (models.Penyesuaian, error) {
	tKepala, err := g.db.Qualify(TabelWarisanPenyesuaian)
	if err != nil {
		return models.Penyesuaian{}, err
	}
	tJSON, err := g.db.Qualify(TabelWarisanPenyesuaianJSON)
	if err != nil {
		return models.Penyesuaian{}, err
	}
	q := fmt.Sprintf(`SELECT t.ID, t.OLDID, m.JSONDATA
		FROM %s t LEFT JOIN %s m ON m.ID = t.ID
		WHERE t.ID = :1`, tKepala, tJSON)
	if err := db.PeriksaSQL(q); err != nil {
		return models.Penyesuaian{}, err
	}

	var sid, asal, dok sql.NullString
	err = g.db.QueryRowContext(ctx, q, id).Scan(&sid, &asal, &dok)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Penyesuaian{}, fmt.Errorf("%w: %s", ErrPenyesuaianTidakAda, id)
	}
	if err != nil {
		return models.Penyesuaian{}, fmt.Errorf("repository: membaca penyesuaian %s: %w", id, err)
	}

	p := models.Penyesuaian{ID: sid.String, IDAsal: asal.String, Baru: sisiKosong(), Lama: sisiKosong()}
	if !dok.Valid || dok.String == "" {
		return p, nil
	}
	baru, lama, err := UraiDokumenPenyesuaian([]byte(dok.String))
	if err != nil {
		// ⛔ Galat yang MENYEBUT pengenalnya — "invalid character" tanpa
		// pengenal membuat yang menyelidiki memeriksa 280 dokumen.
		return models.Penyesuaian{}, fmt.Errorf("%w: %s: %v", ErrJSONPenyesuaianRusak, id, err)
	}
	p.Baru, p.Lama = baru, lama
	return p, nil
}

func sisiKosong() models.SisiPenyesuaian {
	return models.SisiPenyesuaian{Medan: map[string]string{}, Larik: map[string][]map[string]string{}}
}

// UraiDokumenPenyesuaian memisah satu dokumen menjadi sisi New (halaman
// akar) dan sisi Old (halaman `OLDDATA`).
//
// ⛔ `OLDDATA` di dalam `OLDDATA` TIDAK ditelusuri. Diukur: dokumen
// menyimpan `OLDDATA` bersarang, tetapi Section hanya mengikat SATU tingkat
// (`TreatyIn.OLDDATA.*`). Tingkat kedua tidak tampil di layar mana pun.
//
// Diekspor supaya uji tanpa Oracle dapat mengujinya atas dokumen tiruan.
func UraiDokumenPenyesuaian(dok []byte) (baru, lama models.SisiPenyesuaian, err error) {
	akar, err := uraiObjek(dok)
	if err != nil {
		return baru, lama, err
	}
	baru = petik(akar)
	lama = sisiKosong()
	if mentah, ada := akar["OLDDATA"]; ada && !nihil(mentah) {
		old, err := uraiObjek(mentah)
		if err != nil {
			return baru, lama, fmt.Errorf("OLDDATA: %w", err)
		}
		lama = petik(old)
	}
	return baru, lama, nil
}

func uraiObjek(b []byte) (map[string]json.RawMessage, error) {
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	var m map[string]json.RawMessage
	if err := d.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func nihil(b json.RawMessage) bool {
	return len(bytes.TrimSpace(b)) == 0 || string(bytes.TrimSpace(b)) == "null"
}

// LarikBersarangDibaca - larik DI DALAM elemen larik yang sel grid ikat
// lewat indeks (`.RnmLimitListDisplay(1).Value`, sel grid RNM Share
// `TreatyInTabsNonProportional.xml` @2232609…). Elemennya dibawa sebagai
// kunci `Nama(i).Medan`, i mulai 1 — bentuk ikatan Pega apa adanya.
//
// ⛔ Daftar-izin, bukan seluruh larik bersarang: elemen `Share` memuat
// belasan larik bersarang (terukur 969 elemen), dan hanya dua yang tampil.
var LarikBersarangDibaca = []string{"RnmLimitListDisplay", "RnmGrossPremiDisplay"}

// petik mengambil kunci daftar-izin dari satu halaman.
//
// Daftar tangan (`MedanDibaca`, `LarikDibaca`) melayani kepala form; daftar
// bangkitan (`MedanKerangka`, `LarikKerangka`, `kunci_kerangka_gen.go`)
// melayani tab yang dirender dari kerangka ekspor.
func petik(hal map[string]json.RawMessage) models.SisiPenyesuaian {
	s := sisiKosong()
	tertanam := map[string]map[string]json.RawMessage{}
	cari := func(kunci string) (json.RawMessage, bool) {
		induk, anak, bertitik := strings.Cut(kunci, ".")
		if !bertitik {
			r, ada := hal[kunci]
			return r, ada
		}
		h, sudah := tertanam[induk]
		if !sudah {
			if r, ok := hal[induk]; ok {
				h, _ = uraiObjek(r) // bukan objek = tidak ada medan di dalamnya
			}
			tertanam[induk] = h
		}
		r, ada := h[anak]
		return r, ada
	}
	for _, kunci := range gabung(MedanDibaca, MedanKerangka) {
		mentah, ada := cari(kunci)
		if !ada {
			continue
		}
		if v, skalar := teksSkalar(mentah); skalar {
			s.Medan[kunci] = v
		}
	}
	for _, kunci := range gabung(LarikDibaca, LarikKerangka) {
		mentah, ada := cari(kunci)
		if !ada {
			continue
		}
		var elemen []json.RawMessage
		if err := json.Unmarshal(mentah, &elemen); err != nil {
			continue // bukan larik = bukan larik yang dicari
		}
		baris := make([]map[string]string, 0, len(elemen))
		for _, e := range elemen {
			obj, err := uraiObjek(e)
			if err != nil {
				continue
			}
			baris = append(baris, ratakan(obj))
		}
		// ⛔ Urutan di dokumen DIPERTAHANKAN — itulah urutan layar lama.
		s.Larik[kunci] = baris
	}
	return s
}

// ratakan membawa medan skalar satu elemen, ditambah larik bersarang
// daftar-izin sebagai `Nama(i).Medan`.
func ratakan(obj map[string]json.RawMessage) map[string]string {
	b := map[string]string{}
	for k, r := range obj {
		if v, skalar := teksSkalar(r); skalar {
			b[k] = v
		}
	}
	for _, nama := range LarikBersarangDibaca {
		r, ada := obj[nama]
		if !ada {
			continue
		}
		var anak []json.RawMessage
		if json.Unmarshal(r, &anak) != nil {
			continue
		}
		for i, a := range anak {
			o, err := uraiObjek(a)
			if err != nil {
				continue
			}
			for k, rr := range o {
				if v, skalar := teksSkalar(rr); skalar {
					b[fmt.Sprintf("%s(%d).%s", nama, i+1, k)] = v
				}
			}
		}
	}
	return b
}

// gabung menyatukan dua daftar kunci tanpa kembar, urutan dipertahankan.
func gabung(a, b []string) []string {
	sudah := map[string]bool{}
	out := make([]string, 0, len(a)+len(b))
	for _, x := range append(append([]string{}, a...), b...) {
		if !sudah[x] {
			sudah[x] = true
			out = append(out, x)
		}
	}
	return out
}

// teksSkalar mengubah nilai JSON skalar menjadi teks APA ADANYA: string
// tanpa kutipnya, angka dengan digit aslinya (bukan lewat float), boolean
// `true`/`false`, `null` menjadi teks kosong. Objek dan larik bukan skalar.
func teksSkalar(r json.RawMessage) (string, bool) {
	t := bytes.TrimSpace(r)
	if len(t) == 0 {
		return "", false
	}
	switch t[0] {
	case '{', '[':
		return "", false
	case '"':
		var s string
		if err := json.Unmarshal(t, &s); err != nil {
			return "", false
		}
		return s, true
	case 'n':
		return "", true
	default:
		return string(t), true
	}
}
