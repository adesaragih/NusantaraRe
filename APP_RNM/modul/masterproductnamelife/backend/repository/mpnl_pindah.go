package repository

// Pindah JSON → flat dengan rekonsiliasi (tiket 01 bab 02-10-2026, T2/T3; alat `backend/alat/pindahflat`).
//
// ⛔ BUKAN berkas migrasi SQL: di skema uji tabel warisan dibuat SESUDAH migrasi, dan teks JSON path tidak lolos
// penjaga inti. Pengurainya kodek yang sudah teruji (`UraiProduk`), pemetaannya `NormalkanFlat`.
//
//	-uji       (bawaan) hanya SELECT atas kedua tabel JSON; nol tulisan; laporan.
//	-terima-normalisasi=<jenis,…>  (bawaan kosong) JENIS normalisasi teks yang SUDAH diputuskan work owner (OQ-FLAT-07)
//	           tidak menahan -jalankan; jenis lain tetap menahan; tetap dicetak laporan. Kegagalan tidak pernah diterima.
//	-jalankan  menolak IS_PEGA_PROD=true dan rekonsiliasi yang tidak lolos; SATU transaksi yang lebih dulu MENGUNCI tabel
//	           induk (`LOCK TABLE … IN EXCLUSIVE MODE` - aplikasi yang berjalan tidak dapat menulis produk selama pindah):
//	           produk bersumber JSON dihapus lalu diisi ulang, dibaca ulang dan dibandingkan sebelum ditutup. Aman diulang.
//	           Produk yang HANYA ada di tabel flat (tulisan baru aplikasi) dibiarkan utuh; produk bersumber JSON yang
//	           isinya di tabel flat sudah BERBEDA (diubah sesudah peralihan) menolak seluruh putaran - tidak pernah ditimpa.
//
// Rekonsiliasi - SETIAP medan setiap produk, teks demi teks (JSON → model → bentuk flat): beda yang bukan K3/K4 =
// GAGAL; normalisasi teks (mis. nol ekor angka) dicatat per kolom dan menghentikan `-jalankan` untuk keputusan work
// owner. ⛔ Laporan AGREGAT: ID produk + nama kolom/kunci, tidak pernah nilainya.

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

var (
	// ErrPindahDiProduksi - `-jalankan` di lingkungan IS_PEGA_PROD=true.
	ErrPindahDiProduksi = errors.New("repository: moving products to the flat tables is refused when IS_PEGA_PROD=true")
	// ErrPindahTidakLolos - rekonsiliasi gagal atau memuat normalisasi; nol tulisan.
	ErrPindahTidakLolos = errors.New("repository: reconciliation did not pass; nothing was written")
	// ErrTulisanFlatBaru - tabel flat memuat produk yang berbeda dari sumber JSON (sudah ada tulisan baru).
	ErrTulisanFlatBaru = errors.New("repository: the flat tables hold products that do not match the JSON source; refusing to overwrite them")
	// ErrBacaUlangBeda - hasil baca ulang sesudah menulis berbeda dari yang ditulis; transaksi dibatalkan.
	ErrBacaUlangBeda = errors.New("repository: the flat rows read back do not match the rows written; the transaction was rolled back")
)

// Temuan - satu nilai yang dilaporkan, TANPA nilainya.
type Temuan struct {
	ProdukID, Tabel, Kolom string
	Urut                   int
	Jenis                  string
}

func (t Temuan) String() string {
	letak := t.Tabel + "." + t.Kolom
	if t.Urut > 0 {
		letak = fmt.Sprintf("%s[%d].%s", t.Tabel, t.Urut, t.Kolom)
	}
	if t.Jenis == "" {
		return t.ProdukID + " " + letak
	}
	return t.ProdukID + " " + letak + ": " + t.Jenis
}

// LaporanPindah - laporan agregat satu putaran alat pindah.
type LaporanPindah struct {
	Mode                string
	Produk, BarisInward int
	ProdukTanpaInward   int
	InwardBerIDLain     int // ID baris inward ≠ ID produk (dicari lewat PRODUCTID; mis. pasangan bersilang DEV)
	InwardYatim         int // baris inward yang tidak dipilih produk mana pun - tidak dipindah
	Baris               map[string]int
	K3                  []Temuan // nilai tidak sah yang di-NULL-kan (K3)
	K4OutwardKosong     int      // objek OutwardList kosong yang dibuang (K4)
	Normalisasi         map[string]int
	MedanTanpaKolom     map[string]int
	DibuangD2           map[string]int
	Gagal               []string
	TulisanFlatBerbeda  int
	ProdukFlatSaja      int      // produk yang hanya ada di tabel flat (tulisan baru aplikasi) - dibiarkan utuh
	TerimaNormalisasi   []string // jenis normalisasi yang diterima operator (-terima-normalisasi, OQ-FLAT-07)
	TanggalDikonversi   []Temuan // tanggal inward berbentuk lain yang dikonversi (OQ-FLAT-09)
	NormalisasiJenis    map[string]int
	Ditulis             bool
	ContohNormalisasi   []Temuan
	// NormalisasiProduk - ID produk → jenis normalisasi → cacah (Copy Old menilai per produk).
	NormalisasiProduk map[string]map[string]int
}

// Lolos - nol kegagalan DAN nol normalisasi (normalisasi menunggu keputusan work owner, brief T3).
func (l LaporanPindah) Lolos() bool { return len(l.Gagal) == 0 && len(l.Normalisasi) == 0 }

// BolehDitulis - nol kegagalan, dan setiap JENIS normalisasi yang ditemukan diterima operator sesudah keputusan work
// owner (OQ-FLAT-07). Menerima satu jenis tidak pernah menerima jenis lain (mis. `koma desimal` dapat berarti pemisah
// ribuan).
func (l LaporanPindah) BolehDitulis() bool {
	if len(l.Gagal) > 0 {
		return false
	}
	diterima := map[string]bool{}
	for _, j := range l.TerimaNormalisasi {
		diterima[strings.TrimSpace(j)] = true
	}
	for j := range l.NormalisasiJenis {
		if !diterima[j] {
			return false
		}
	}
	return true
}

func laporanBaru() LaporanPindah {
	return LaporanPindah{Baris: map[string]int{}, Normalisasi: map[string]int{}, MedanTanpaKolom: map[string]int{},
		DibuangD2: map[string]int{}, NormalisasiJenis: map[string]int{}, NormalisasiProduk: map[string]map[string]int{}}
}

// Beda - satu medan yang berbeda antara dua produk.
type Beda struct {
	Tabel, Kolom string
	Urut         int
}

func bedaMedan[T any](kolom []KolomFlat[T], tabel string, urut int, a, b *T, hasil *[]Beda) {
	for _, k := range kolom {
		if *k.ambil(a) != *k.ambil(b) {
			*hasil = append(*hasil, Beda{Tabel: tabel, Kolom: k.Nama, Urut: urut})
		}
	}
}

func bedaAnak[T any](a AnakFlat[T], x, y *models.Produk, hasil *[]Beda) {
	dx, dy := *a.daftar(x), *a.daftar(y)
	if len(dx) != len(dy) {
		*hasil = append(*hasil, Beda{Tabel: a.Tabel, Kolom: fmt.Sprintf("jumlah baris %d lawan %d", len(dx), len(dy))})
		return
	}
	for i := range dx {
		bedaMedan(a.Kolom, a.Tabel, i+1, &dx[i], &dy[i], hasil)
	}
}

// BedaProduk - medan BERKOLOM yang berbeda antara dua produk (teks demi teks), termasuk pemegang polis sisi umum.
func BedaProduk(a, b models.Produk) []Beda {
	var hasil []Beda
	bedaMedan(KolomFlatInduk, TabelFlatInduk, 0, &a, &b, &hasil)
	if a.Umum.IsORS != b.Umum.IsORS {
		hasil = append(hasil, Beda{Tabel: TabelFlatInduk, Kolom: KolomIsORS})
	}
	if a.Umum.PolicyHolder != b.Umum.PolicyHolder || a.Umum.PolicyHolderName != b.Umum.PolicyHolderName {
		hasil = append(hasil, Beda{Tabel: TabelFlatInduk, Kolom: "POLICYHOLDER (sisi umum)"})
	}
	bedaAnak(AnakLien, &a, &b, &hasil)
	bedaAnak(AnakDokumen, &a, &b, &hasil)
	bedaAnak(AnakPlan, &a, &b, &hasil)
	bedaAnak(AnakFinUW, &a, &b, &hasil)
	bedaAnak(AnakUWLimit, &a, &b, &hasil)
	bedaAnak(AnakOutward, &a, &b, &hasil)
	bedaAnak(AnakKomentar, &a, &b, &hasil)
	return hasil
}

// nilaiMedan - nilai satu medan berkolom menurut Beda (untuk menggolongkan normalisasi; tidak pernah dicetak).
func nilaiMedan(p models.Produk, b Beda) string {
	if b.Tabel == TabelFlatInduk {
		for _, k := range KolomFlatInduk {
			if k.Nama == b.Kolom {
				return *k.ambil(&p)
			}
		}
		return ""
	}
	for _, cari := range []func() (string, bool){
		func() (string, bool) { return nilaiAnak(AnakLien, &p, b) },
		func() (string, bool) { return nilaiAnak(AnakDokumen, &p, b) },
		func() (string, bool) { return nilaiAnak(AnakPlan, &p, b) },
		func() (string, bool) { return nilaiAnak(AnakFinUW, &p, b) },
		func() (string, bool) { return nilaiAnak(AnakUWLimit, &p, b) },
		func() (string, bool) { return nilaiAnak(AnakOutward, &p, b) },
		func() (string, bool) { return nilaiAnak(AnakKomentar, &p, b) },
	} {
		if v, ada := cari(); ada {
			return v
		}
	}
	return ""
}

func nilaiAnak[T any](a AnakFlat[T], p *models.Produk, b Beda) (string, bool) {
	daftar := *a.daftar(p)
	if a.Tabel != b.Tabel || b.Urut < 1 || b.Urut > len(daftar) {
		return "", false
	}
	for _, k := range a.Kolom {
		if k.Nama == b.Kolom {
			return *k.ambil(&daftar[b.Urut-1]), true
		}
	}
	return "", false
}

// JenisNormalisasi - golongan perubahan teks nilai ASAL menjadi bentuk kanonik, untuk keputusan work owner tanpa
// mencetak nilainya: koma desimal, spasi tepi, nol ekor desimal, titik tanpa pecahan, tanpa nol depan, nol depan,
// tanda plus, atau lain.
func JenisNormalisasi(asal string) string {
	t := strings.TrimSpace(asal)
	switch {
	case t != asal:
		return "spasi tepi"
	case strings.Contains(t, ","):
		return "koma desimal"
	case strings.HasPrefix(t, "+"):
		return "tanda plus"
	case strings.HasSuffix(t, "."):
		return "titik tanpa pecahan"
	case strings.Contains(t, ".") && strings.HasSuffix(t, "0"):
		return "nol ekor desimal"
	case strings.HasPrefix(t, ".") || strings.HasPrefix(t, "-."):
		return "tanpa nol depan"
	case len(t) > 1 && strings.HasPrefix(strings.TrimPrefix(t, "-"), "0") && !strings.HasPrefix(strings.TrimPrefix(t, "-"), "0."):
		return "nol depan"
	default:
		return "lain"
	}
}

// k3Dibolehkan - K3 keputusan work owner 02-10-2026: hanya dua bentuk nilai tidak sah yang di-NULL-kan - dan hanya bila
// nilainya memang tidak terbaca. Tanggal berbentuk lain (`1/3/2027`, `01-03-2027`, `2027/03/01`) dan angka berpemisah
// ribuan (`1.000.000`) DAPAT diselamatkan: keduanya GAGAL (menunggu keputusan), tidak di-NULL-kan diam-diam (temuan
// /code-review 02-10-2026).
func k3Dibolehkan(m MasalahNilai) bool {
	switch {
	case m.Tabel == TabelFlatInduk && m.Kolom == "MATURE" && m.Jenis == MasalahBukanTanggal:
		return !tanggalBentukLain(m.Nilai)
	case m.Tabel == TabelFlatUWLimit && m.Kolom == "MAXINSURED" && m.Jenis == MasalahBukanAngka:
		return !angkaBerpemisah(m.Nilai)
	}
	return false
}

// k3Bentuk - kolom dan jenis masalahnya bentuk K3 (terlepas dari dapat-tidaknya diselamatkan).
func k3Bentuk(m MasalahNilai) bool {
	return (m.Tabel == TabelFlatInduk && m.Kolom == "MATURE" && m.Jenis == MasalahBukanTanggal) ||
		(m.Tabel == TabelFlatUWLimit && m.Kolom == "MAXINSURED" && m.Jenis == MasalahBukanAngka)
}

// bentukTanggalLain - bentuk lazim selain `dd/MM/yyyy` (urutan hari-bulan seperti Pega; tahun-bulan-hari bila tahun
// di depan). Nama bentuk dicetak laporan (bukan nilainya).
var bentukTanggalLain = []struct{ tata, nama string }{
	{"2/1/2006", "d/M/yyyy"}, {"02-01-2006", "dd-MM-yyyy"}, {"2-1-2006", "d-M-yyyy"}, {"2006/01/02", "yyyy/MM/dd"},
	{"2006/1/2", "yyyy/M/d"}, {"2006-1-2", "yyyy-M-d"}, {"02.01.2006", "dd.MM.yyyy"}, {"2/1/06", "d/M/yy"},
	{"02/01/06", "dd/MM/yy"},
}

// konversiTanggalBentukLain - teks tanggal dalam bentuk lain → `YYYY-MM-DD` dan nama bentuknya; ok = terkonversi.
// Kosong atau sudah `YYYY-MM-DD` = tidak ada yang dikonversi.
func konversiTanggalBentukLain(v string) (string, string, bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", "", false
	}
	if _, err := time.Parse(bentukTanggalAPI, v); err == nil {
		return "", "", false
	}
	for _, b := range bentukTanggalLain {
		if t, err := time.Parse(b.tata, v); err == nil && t.Year() >= 1 {
			return t.Format(bentukTanggalAPI), b.nama, true
		}
	}
	return "", "", false
}

// tanggalBentukLain - teks yang terbaca sebagai tanggal dalam salah satu bentuk lazim selain `dd/MM/yyyy`.
func tanggalBentukLain(v string) bool {
	_, _, ok := konversiTanggalBentukLain(v)
	return ok
}

// angkaBerpemisah - teks yang menjadi bilangan bila titik, koma, dan spasi pemisahnya dibuang.
func angkaBerpemisah(v string) bool {
	t := strings.TrimPrefix(strings.NewReplacer(".", "", ",", "", " ", "").Replace(strings.TrimSpace(v)), "-")
	if t == "" {
		return false
	}
	for _, c := range t {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// kunciInternalPega - kunci sistem Pega (`px*`, `py*`, `pz*`): bukan data (D2).
func kunciInternalPega(k string) bool {
	return strings.HasPrefix(k, "px") || strings.HasPrefix(k, "py") || strings.HasPrefix(k, "pz")
}

// rawTerisi - nilai JSON berisi (bukan absen, null, teks kosong, larik/objek kosong).
func rawTerisi(raw json.RawMessage) bool {
	switch strings.TrimSpace(string(raw)) {
	case "", "null", `""`, "[]", "{}":
		return false
	}
	return true
}

// outwardKosong - objek OutwardList tanpa satu pun medan terisi (K4).
func outwardKosong(b models.BarisOutward) bool {
	for _, k := range AnakOutward.Kolom {
		if strings.TrimSpace(*k.ambil(&b)) != "" {
			return false
		}
	}
	return true
}

// cacahLarik - jumlah elemen larik JSON (satu objek tunggal = 1, absen/null = 0), seperti uraiBaris.
func cacahLarik(raw json.RawMessage) int {
	t := strings.TrimSpace(string(raw))
	if t == "" || t == "null" {
		return 0
	}
	var larik []json.RawMessage
	if err := json.Unmarshal(raw, &larik); err != nil {
		return 1
	}
	return len(larik)
}

// kunciDikelola - kunci halaman yang dipetakan kodek (selain itu: kunci internal Pega, dibuang D2).
func kunciDikelola[T any](medan []medanTeks[T], tambahan ...string) map[string]bool {
	h := map[string]bool{}
	for _, m := range medan {
		h[m.kunci] = true
	}
	for _, k := range tambahan {
		h[k] = true
	}
	return h
}

var (
	kunciUmumDikelola = kunciDikelola(medanUmum, kunciID, kunciIsORS, kunciIsView, kunciLien, kunciDokumen, kunciPlan,
		kunciFinUW, kunciUWLimit, kunciOutward, kunciKomentar)
	kunciInwardDikelola = kunciDikelola(medanInward, kunciID, kunciProductID)
)

// rekonsiliasi - murni: kedua tabel JSON (semua baris) → laporan + produk berbentuk flat yang akan ditulis.
func rekonsiliasi(umum, inward []barisJSON) (LaporanPindah, []models.Produk) {
	lap := laporanBaru()
	lap.Produk, lap.BarisInward = len(umum), len(inward)
	sort.Slice(umum, func(i, j int) bool { return umum[i].id < umum[j].id })
	gagal := func(format string, a ...any) { lap.Gagal = append(lap.Gagal, fmt.Sprintf(format, a...)) }
	dipilih := map[int]bool{}
	sudah := map[string]bool{}
	// PRODUCTID setiap baris inward diurai SEKALI (dulu sekali per pasangan produk × baris).
	pid := make([]string, len(inward))
	for i, b := range inward {
		pid[i] = productIDDari(b)
	}
	var hasil []models.Produk
	for _, u := range umum {
		if sudah[u.id] {
			gagal("%s: ID produk ganda di %s", u.id, TabelProduk)
			continue
		}
		sudah[u.id] = true
		var kandidat []barisJSON
		var indeks []int
		for i, b := range inward {
			if b.id == u.id || pid[i] == u.id {
				kandidat, indeks = append(kandidat, b), append(indeks, i)
			}
		}
		in, ada := pilihInward(u.id, kandidat)
		for j, b := range kandidat {
			if ada && b == in {
				dipilih[indeks[j]] = true
			}
		}
		switch {
		case !ada:
			lap.ProdukTanpaInward++
		case in.id != u.id:
			lap.InwardBerIDLain++
		}
		p, err := UraiProduk(u.id, u.isi, in.id, in.isi)
		if err != nil {
			gagal("%s: JSON tidak terbaca", u.id)
			continue
		}
		hitungDibuang(&lap, u.id, u.isi, in.isi, p)
		// Kolom datar `M_PRODUCT_LIFE` (`SaveProductNameLIfeFlat` b84) sama dengan kunci JSON-nya - bila beda, mana yang
		// benar adalah keputusan, bukan pilihan diam-diam alat ini.
		if u.adaDatar && (u.datar[0] != p.Umum.RIRiskID || u.datar[1] != p.Umum.RIRisk) {
			gagal("%s: kolom datar RIRISKID/RIRISK %s berbeda dari kunci JSON-nya", u.id, TabelProduk)
		}
		// K4: objek outward kosong tidak dipindah.
		var outward []models.BarisOutward
		for _, b := range p.OutwardList {
			if outwardKosong(b) {
				lap.K4OutwardKosong++
				continue
			}
			outward = append(outward, b)
		}
		if outward == nil {
			outward = []models.BarisOutward{}
		}
		p.OutwardList = outward
		// OQ-FLAT-09 (keputusan work owner 02-10-2026 "konversi"): tanggal inward berbentuk lain dikonversi ke tanggal
		// yang sama, dicatat per produk dan kolom (bentuknya, bukan nilainya) - tidak di-NULL-kan.
		for _, k := range []struct {
			kolom string
			v     *string
		}{{"BEGIN_DATE", &p.Inward.Begin}, {"STNC", &p.Inward.STNC}, {"MATURE", &p.Inward.Mature}} {
			if api, bentuk, ok := konversiTanggalBentukLain(*k.v); ok {
				*k.v = api
				lap.TanggalDikonversi = append(lap.TanggalDikonversi, Temuan{ProdukID: u.id, Tabel: TabelFlatInduk,
					Kolom: k.kolom, Jenis: "dari bentuk " + bentuk})
			}
		}
		for _, k := range MedanTanpaKolom(p) {
			lap.MedanTanpaKolom[k]++
			gagal("%s: %s %s", u.id, k, pesanTanpaKolom)
		}
		n, masalah := NormalkanFlat(p)
		dijelaskan := map[Beda]bool{}
		for _, m := range masalah {
			t := Temuan{ProdukID: u.id, Tabel: m.Tabel, Kolom: m.Kolom, Urut: m.Urut, Jenis: m.Jenis}
			dijelaskan[Beda{Tabel: m.Tabel, Kolom: m.Kolom, Urut: m.Urut}] = true
			if k3Dibolehkan(m) {
				lap.K3 = append(lap.K3, t)
				continue
			}
			if k3Bentuk(m) {
				gagal("%s (bentuk lain yang dapat diselamatkan - bukan K3, menunggu keputusan)", t)
				continue
			}
			gagal("%s", t)
		}
		for _, b := range BedaProduk(p, n) {
			if dijelaskan[b] {
				continue
			}
			if b.Kolom == "POLICYHOLDER (sisi umum)" {
				gagal("%s: pemegang polis sisi umum berbeda dari sisi inward", u.id)
				continue
			}
			jenis := JenisNormalisasi(nilaiMedan(p, b))
			lap.Normalisasi[b.Tabel+"."+b.Kolom+" ("+jenis+")"]++
			lap.NormalisasiJenis[jenis]++
			if lap.NormalisasiProduk[u.id] == nil {
				lap.NormalisasiProduk[u.id] = map[string]int{}
			}
			lap.NormalisasiProduk[u.id][jenis]++
			if len(lap.ContohNormalisasi) < 20 {
				lap.ContohNormalisasi = append(lap.ContohNormalisasi, Temuan{ProdukID: u.id, Tabel: b.Tabel,
					Kolom: b.Kolom, Urut: b.Urut})
			}
		}
		periksaCacahAnak(&lap, u.id, u.isi, n)
		lap.Baris[TabelFlatInduk]++
		lap.Baris[TabelFlatLien] += len(n.LienClause)
		lap.Baris[TabelFlatDokumen] += len(n.DocumentClaim)
		lap.Baris[TabelFlatPlan] += len(n.PlanList)
		lap.Baris[TabelFlatFinUW] += len(n.FinancialUnderwriting)
		lap.Baris[TabelFlatUWLimit] += len(n.UnderwritingLimit)
		lap.Baris[TabelFlatOutward] += len(n.OutwardList)
		lap.Baris[TabelFlatKomentar] += len(n.CommentList)
		hasil = append(hasil, n)
	}
	for i := range inward {
		if !dipilih[i] {
			lap.InwardYatim++
		}
	}
	if lap.InwardYatim > 0 {
		gagal("%d baris %s tidak dimiliki produk mana pun (yatim) - tidak dipindah", lap.InwardYatim, TabelInward)
	}
	return lap, hasil
}

// periksaCacahAnak - jumlah baris anak = jumlah elemen larik JSON (OutwardList: dikurangi objek kosong K4).
func periksaCacahAnak(lap *LaporanPindah, id, jsonUmum string, n models.Produk) {
	obj, err := uraiObjek(jsonUmum)
	if err != nil {
		return
	}
	kosong := 0
	if raw := obj[kunciOutward]; len(raw) > 0 {
		if baris, err := uraiBaris(raw, kodekOutward); err == nil {
			for _, b := range baris {
				if outwardKosong(b) {
					kosong++
				}
			}
		}
	}
	for _, c := range []struct {
		kunci  string
		dapat  int
		kurang int
	}{
		{kunciLien, len(n.LienClause), 0}, {kunciDokumen, len(n.DocumentClaim), 0}, {kunciPlan, len(n.PlanList), 0},
		{kunciFinUW, len(n.FinancialUnderwriting), 0}, {kunciUWLimit, len(n.UnderwritingLimit), 0},
		{kunciOutward, len(n.OutwardList), kosong}, {kunciKomentar, len(n.CommentList), 0},
	} {
		if mau := cacahLarik(obj[c.kunci]) - c.kurang; mau != c.dapat {
			lap.Gagal = append(lap.Gagal, fmt.Sprintf("%s: %s %d baris JSON, %d baris flat", id, c.kunci, mau, c.dapat))
		}
	}
}

// hitungDibuang - kunci yang tidak dipindah. Kunci internal Pega (`px*`/`py*`/`pz*`), `IsView` (keadaan layar), dan
// `Comment` (masukan popup) = dibuang D2, dicacah. Kunci LAIN yang tidak dikelola kodek dan BERISI = data tanpa kolom:
// GAGAL (temuan /code-review 02-10-2026 - dulu seluruh kunci baris tak dikelola dicacah "internal Pega").
func hitungDibuang(lap *LaporanPindah, id, jsonUmum, jsonInward string, p models.Produk) {
	periksaKunci := func(asal string, obj map[string]json.RawMessage, dikelola map[string]bool) {
		for k, v := range obj {
			switch {
			case dikelola[k]:
			case kunciInternalPega(k):
				lap.DibuangD2[asal+" "+k]++
			case rawTerisi(v):
				lap.MedanTanpaKolom[asal+" "+k]++
				lap.Gagal = append(lap.Gagal, fmt.Sprintf("%s: %s %s %s", id, asal, k, pesanTanpaKolom))
			default:
				lap.DibuangD2[asal+" "+k+" (kosong)"]++
			}
		}
	}
	if obj, err := uraiObjek(jsonUmum); err == nil {
		periksaKunci("kunci halaman umum", obj, kunciUmumDikelola)
		if _, ada := obj[kunciIsView]; ada {
			lap.DibuangD2["IsView (keadaan layar)"]++
		}
	}
	if obj, err := uraiObjek(jsonInward); err == nil {
		periksaKunci("kunci halaman inward", obj, kunciInwardDikelola)
	}
	if strings.TrimSpace(p.Umum.Comment) != "" {
		lap.DibuangD2["Comment (masukan popup; isinya sudah baris CommentList)"]++
	}
	for _, c := range []struct {
		daftar string
		asli   []string
	}{
		{kunciLien, asliDari(p.LienClause, func(b models.BarisLien) string { return b.Asli })},
		{kunciDokumen, asliDari(p.DocumentClaim, func(b models.BarisDokumen) string { return b.Asli })},
		{kunciPlan, asliDari(p.PlanList, func(b models.BarisPlan) string { return b.Asli })},
		{kunciFinUW, asliDari(p.FinancialUnderwriting, func(b models.BarisFinUW) string { return b.Asli })},
		{kunciUWLimit, asliDari(p.UnderwritingLimit, func(b models.BarisUWLimit) string { return b.Asli })},
		{kunciOutward, asliDari(p.OutwardList, func(b models.BarisOutward) string { return b.Asli })},
		{kunciKomentar, asliDari(p.CommentList, func(b models.BarisKomentar) string { return b.Asli })},
	} {
		for _, a := range c.asli {
			if obj, err := uraiObjek(a); err == nil {
				periksaKunci("baris "+c.daftar+" kunci", obj, map[string]bool{})
			}
		}
	}
}

func asliDari[T any](daftar []T, asli func(T) string) []string {
	var hasil []string
	for _, b := range daftar {
		if a := asli(b); a != "" {
			hasil = append(hasil, a)
		}
	}
	return hasil
}

// Batas baris laporan (laporan tetap terbaca): kegagalan unik dan contoh kegagalan "tanpa kolom flat".
const (
	MaksBarisGagalDicetak = 25
	MaksContohTanpaKolom  = 8
)

// pesanTanpaKolom - akhir kalimat kegagalan nilai yang tidak punya kolom flat.
const pesanTanpaKolom = "terisi tetapi tidak punya kolom flat"

// Teks - laporan untuk manusia (agregat; tidak satu nilai data pun).
func (l LaporanPindah) Teks() string {
	var b strings.Builder
	fmt.Fprintf(&b, "pindahflat -%s\n", l.Mode)
	fmt.Fprintf(&b, "  sumber: %s %d baris, %s %d baris\n", TabelProduk, l.Produk, TabelInward, l.BarisInward)
	fmt.Fprintf(&b, "  produk tanpa inward: %d; inward ber-ID lain (dicari lewat PRODUCTID): %d; inward yatim: %d\n",
		l.ProdukTanpaInward, l.InwardBerIDLain, l.InwardYatim)
	b.WriteString("  baris yang (akan) ditulis:\n")
	for _, t := range DaftarTabelFlat {
		fmt.Fprintf(&b, "    %-30s %d\n", t, l.Baris[t])
	}
	fmt.Fprintf(&b, "  K3 nilai tidak sah di-NULL-kan: %d\n", len(l.K3))
	for _, t := range l.K3 {
		fmt.Fprintf(&b, "    %s\n", t)
	}
	fmt.Fprintf(&b, "  K4 objek OutwardList kosong dibuang: %d\n", l.K4OutwardKosong)
	fmt.Fprintf(&b, "  tanggal bentuk lain dikonversi (OQ-FLAT-09): %d\n", len(l.TanggalDikonversi))
	for _, t := range l.TanggalDikonversi {
		fmt.Fprintf(&b, "    %s\n", t)
	}
	tulisPeta(&b, "  normalisasi teks (menunggu keputusan work owner)", l.Normalisasi)
	if len(l.TerimaNormalisasi) > 0 {
		fmt.Fprintf(&b, "    jenis diterima operator (-terima-normalisasi): %s\n", strings.Join(l.TerimaNormalisasi, ", "))
	}
	for _, t := range l.ContohNormalisasi {
		fmt.Fprintf(&b, "    contoh: %s\n", t)
	}
	tulisPeta(&b, "  medan tanpa kolom yang terisi", l.MedanTanpaKolom)
	tulisPeta(&b, "  dibuang D2 (bukan data)", l.DibuangD2)
	if l.TulisanFlatBerbeda > 0 {
		fmt.Fprintf(&b, "  tabel flat sudah memuat %d produk bersumber JSON yang isinya berbeda (diubah sesudah peralihan)\n",
			l.TulisanFlatBerbeda)
	}
	if l.ProdukFlatSaja > 0 {
		fmt.Fprintf(&b, "  produk yang hanya ada di tabel flat (tulisan baru aplikasi, dibiarkan): %d\n", l.ProdukFlatSaja)
	}
	// Kegagalan unik lebih dulu; kegagalan "tanpa kolom flat" berulang per baris dan sudah dicacah per kunci di atas,
	// jadi hanya contohnya yang dicetak.
	var unik, tanpaKolom []string
	for _, g := range l.Gagal {
		if strings.HasSuffix(g, pesanTanpaKolom) {
			tanpaKolom = append(tanpaKolom, g)
		} else {
			unik = append(unik, g)
		}
	}
	fmt.Fprintf(&b, "  gagal: %d (%d tanpa kolom flat)\n", len(l.Gagal), len(tanpaKolom))
	for i, g := range unik {
		if i == MaksBarisGagalDicetak {
			fmt.Fprintf(&b, "    ... dan %d lainnya\n", len(unik)-MaksBarisGagalDicetak)
			break
		}
		fmt.Fprintf(&b, "    %s\n", g)
	}
	for i, g := range tanpaKolom {
		if i == MaksContohTanpaKolom {
			fmt.Fprintf(&b, "    ... dan %d lainnya tanpa kolom flat (cacah per kunci di atas)\n", len(tanpaKolom)-MaksContohTanpaKolom)
			break
		}
		fmt.Fprintf(&b, "    %s\n", g)
	}
	status := "LOLOS"
	switch {
	case !l.Lolos() && l.BolehDitulis():
		status = "LOLOS DENGAN NORMALISASI DITERIMA OPERATOR (-terima-normalisasi, OQ-FLAT-07)"
	case !l.Lolos():
		status = "TIDAK LOLOS"
	}
	fmt.Fprintf(&b, "  rekonsiliasi: %s; ditulis: %v\n", status, l.Ditulis)
	return b.String()
}

func tulisPeta(b *strings.Builder, judul string, m map[string]int) {
	kunci := make([]string, 0, len(m))
	jumlah := 0
	for k, v := range m {
		kunci = append(kunci, k)
		jumlah += v
	}
	sort.Strings(kunci)
	fmt.Fprintf(b, "%s: %d\n", judul, jumlah)
	for _, k := range kunci {
		fmt.Fprintf(b, "    %-58s %d\n", k, m[k])
	}
}

// --- Oracle --------------------------------------------------------------------------

func sqlSemuaJSON(tabel string) string {
	return fmt.Sprintf(`SELECT ID, JSONDATA FROM %s ORDER BY ID ASC`, tabel)
}

// sqlSemuaJSONUmum - `M_PRODUCT_LIFE` beserta kolom datarnya (direkonsiliasi dengan kunci JSON-nya).
func sqlSemuaJSONUmum(tabel string) string {
	return fmt.Sprintf(`SELECT ID, JSONDATA, RIRISKID, RIRISK FROM %s ORDER BY ID ASC`, tabel)
}

// sqlKunciInduk - tabel induk dikunci selama pindah: penulis aplikasi mulai dari induk, jadi ia menunggu.
func sqlKunciInduk(tabel string) string { return fmt.Sprintf(`LOCK TABLE %s IN EXCLUSIVE MODE`, tabel) }

// sqlHapusInduk - satu produk; anak ikut terhapus (FK `ON DELETE CASCADE`, migrasi 141–147).
func sqlHapusInduk(tabel string) string { return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1`, tabel) }

// bacaSumberJSON - seluruh baris kedua tabel JSON warisan (SELECT saja).
func (g *Gudang) bacaSumberJSON(ctx context.Context) (umum, inward []barisJSON, err error) {
	q, err := g.siapkan(TabelProduk, sqlSemuaJSONUmum)
	if err != nil {
		return nil, nil, err
	}
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, nil, fmt.Errorf("repository: reading %s: %w", TabelProduk, err)
	}
	for rows.Next() {
		var id, isi, riskID, risk sql.NullString
		if err := rows.Scan(&id, &isi, &riskID, &risk); err != nil {
			_ = rows.Close()
			return nil, nil, fmt.Errorf("repository: reading %s: %w", TabelProduk, err)
		}
		umum = append(umum, barisJSON{id: id.String, isi: isi.String, datar: [2]string{riskID.String, risk.String},
			adaDatar: true})
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, nil, fmt.Errorf("repository: reading %s: %w", TabelProduk, err)
	}
	_ = rows.Close()
	qi, err := g.siapkan(TabelInward, sqlSemuaJSON)
	if err != nil {
		return nil, nil, err
	}
	if inward, err = g.bacaBaris(ctx, nil, qi); err != nil {
		return nil, nil, fmt.Errorf("repository: reading %s: %w", TabelInward, err)
	}
	return umum, inward, nil
}

// PindahFlat - satu putaran alat pindah (`jalankan` = `-jalankan`, selain itu `-uji`; `terima` = jenis normalisasi
// `-terima-normalisasi`).
func (g *Gudang) PindahFlat(ctx context.Context, jalankan bool, terima []string) (LaporanPindah, error) {
	mode := "uji"
	if jalankan {
		mode = "jalankan"
		if g.db.PegaProduksi() {
			return LaporanPindah{Mode: mode}, ErrPindahDiProduksi
		}
	}
	umum, inward, err := g.bacaSumberJSON(ctx)
	if err != nil {
		return LaporanPindah{Mode: mode}, err
	}
	lap, produk := rekonsiliasi(umum, inward)
	lap.Mode, lap.TerimaNormalisasi = mode, terima
	if !jalankan {
		return lap, nil
	}
	if !lap.BolehDitulis() {
		return lap, ErrPindahTidakLolos
	}
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return lap, err
	}
	defer func() { _ = tx.Rollback() }()
	// ⛔ Kunci induk LEBIH DULU (temuan /code-review 02-10-2026): tanpa ini produk yang dibuat aplikasi di antara
	// pemeriksaan dan penulisan dapat hilang. Penulis aplikasi menunggu sampai transaksi ini ditutup.
	if err := g.eksekusi(ctx, tx, TabelFlatInduk, sqlKunciInduk, false); err != nil {
		return lap, err
	}
	menurutID := map[string]models.Produk{}
	for _, p := range produk {
		menurutID[p.ID] = p
	}
	ada, err := g.semuaIDFlat(ctx, tx)
	if err != nil {
		return lap, err
	}
	for _, id := range ada {
		p, dariJSON := menurutID[id]
		if !dariJSON {
			lap.ProdukFlatSaja++
			continue
		}
		f, err := g.bacaFlat(ctx, tx, id, false)
		if err != nil {
			return lap, err
		}
		if len(BedaProduk(p, f)) > 0 {
			lap.TulisanFlatBerbeda++
		}
	}
	if lap.TulisanFlatBerbeda > 0 {
		return lap, fmt.Errorf("%w (%d products)", ErrTulisanFlatBaru, lap.TulisanFlatBerbeda)
	}
	// Hanya produk bersumber JSON: produk yang hanya ada di tabel flat tidak disentuh.
	for _, p := range produk {
		if err := g.eksekusi(ctx, tx, TabelFlatInduk, sqlHapusInduk, false, p.ID); err != nil {
			return lap, err
		}
		if err := g.tulisFlat(ctx, tx, p, true); err != nil {
			return lap, err
		}
	}
	for _, p := range produk {
		f, err := g.bacaFlat(ctx, tx, p.ID, false)
		if err != nil {
			return lap, err
		}
		if beda := BedaProduk(p, f); len(beda) > 0 {
			return lap, fmt.Errorf("%w: %s %s.%s", ErrBacaUlangBeda, p.ID, beda[0].Tabel, beda[0].Kolom)
		}
	}
	if err := tx.Commit(); err != nil {
		return lap, err
	}
	lap.Ditulis = true
	return lap, nil
}

// --- Copy Old (permintaan work owner 03-10-2026) -----------------------------------------------------------------
//
// Tombol `Copy Old` di samping `Add` membuka popup berisi produk tabel JSON warisan yang BELUM ada di induk flat; yang
// dicentang disalin lewat `Process Copy`. Aturannya SAMA dengan alat pindah - rekonsiliasi seluruh sumber, K3,
// keputusan OQ-FLAT-07 (jenis normalisasi diterima) dan OQ-FLAT-09 (tanggal dikonversi) - tetapi per produk: satu
// transaksi per produk, tulis lalu baca ulang. Kedua tabel JSON hanya DIBACA (penjaga TestMPNLAplikasiHanyaTabelFlat:
// berkas ini satu-satunya kode aplikasi yang menyebutnya).

// NormalisasiDiputuskan - jenis normalisasi yang diterima work owner 02-10-2026 (OQ-FLAT-07, "ikuti rekomendasi").
var NormalisasiDiputuskan = []string{"koma desimal", "nol depan"}

// ErrSudahDiFlat - ID produk lama sudah ada di induk flat (disalin sebelumnya, oleh pemakai lain, atau alat pindah).
var ErrSudahDiFlat = errors.New("repository: the product is already in the flat tables")

// statusLama - murni: rekonsiliasi seluruh sumber JSON → baris popup untuk produk yang BELUM ada di induk flat
// (`diFlat`), urut ID, dan bentuk flat produk yang boleh disalin. Alasan dan catatan menyebut tabel dan kolom saja.
func statusLama(umum, inward []barisJSON, diFlat map[string]bool) ([]models.ProdukLama, map[string]models.Produk) {
	lap, produk := rekonsiliasi(umum, inward)
	menurutID := make(map[string]models.Produk, len(produk))
	for _, p := range produk {
		menurutID[p.ID] = p
	}
	diterima := map[string]bool{}
	for _, j := range NormalisasiDiputuskan {
		diterima[j] = true
	}
	daftar := []models.ProdukLama{}
	siap := map[string]models.Produk{}
	sudah := map[string]bool{}
	for _, u := range umum {
		if diFlat[u.id] || sudah[u.id] {
			continue
		}
		sudah[u.id] = true
		d := models.ProdukLama{ID: u.id, Alasan: []string{}, Catatan: []string{}}
		for _, g := range lap.Gagal {
			// Kalimat kegagalan berawalan ID produk: `<id>: …` atau `<id> <tabel>.<kolom>…` (Temuan).
			if sisa, ok := strings.CutPrefix(g, u.id+":"); ok {
				d.Alasan = append(d.Alasan, strings.TrimSpace(sisa))
			} else if sisa, ok := strings.CutPrefix(g, u.id+" "); ok {
				d.Alasan = append(d.Alasan, strings.TrimSpace(sisa))
			}
		}
		jenis := make([]string, 0, len(lap.NormalisasiProduk[u.id]))
		for j := range lap.NormalisasiProduk[u.id] {
			jenis = append(jenis, j)
		}
		sort.Strings(jenis)
		for _, j := range jenis {
			n := lap.NormalisasiProduk[u.id][j]
			if diterima[j] {
				d.Catatan = append(d.Catatan, fmt.Sprintf("angka dirapikan (%s), %d nilai - diputuskan OQ-FLAT-07", j, n))
			} else {
				d.Alasan = append(d.Alasan, fmt.Sprintf("angka berbentuk %s (%d nilai) belum diputuskan work owner (OQ-FLAT-07)", j, n))
			}
		}
		for _, t := range lap.K3 {
			if t.ProdukID == u.id {
				d.Catatan = append(d.Catatan, "nilai tidak sah dikosongkan (K3) - "+strings.TrimPrefix(t.String(), u.id+" "))
			}
		}
		for _, t := range lap.TanggalDikonversi {
			if t.ProdukID == u.id {
				d.Catatan = append(d.Catatan, "tanggal dikonversi (OQ-FLAT-09) - "+strings.TrimPrefix(t.String(), u.id+" "))
			}
		}
		p, ada := menurutID[u.id]
		if ada {
			d.ProductName, d.Ceding, d.TreatyNumber = p.Umum.ProductName, p.Umum.Ceding, p.Umum.TreatyNumber
			d.InwardName, d.CreateOp, d.UpdateOp = p.Umum.InwardName, p.Umum.CreateOp, p.Umum.UpdateOp
		} else if len(d.Alasan) == 0 {
			d.Alasan = append(d.Alasan, "produk tidak dapat dibentuk dari JSON-nya")
		}
		d.BolehDisalin = ada && len(d.Alasan) == 0
		if d.BolehDisalin {
			siap[u.id] = p
		}
		daftar = append(daftar, d)
	}
	sort.Slice(daftar, func(i, j int) bool { return daftar[i].ID < daftar[j].ID })
	return daftar, siap
}

// SiapkanProdukLama - Copy Old: baris popup dan bentuk flat produk yang boleh disalin. SELECT saja (kedua tabel JSON
// dan ID induk flat).
func (g *Gudang) SiapkanProdukLama(ctx context.Context) ([]models.ProdukLama, map[string]models.Produk, error) {
	umum, inward, err := g.bacaSumberJSON(ctx)
	if err != nil {
		return nil, nil, err
	}
	ids, err := g.semuaIDFlat(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	diFlat := make(map[string]bool, len(ids))
	for _, id := range ids {
		diFlat[id] = true
	}
	daftar, siap := statusLama(umum, inward, diFlat)
	return daftar, siap, nil
}

// SalinProdukLama - Copy Old: satu produk lama (bentuk flat dari SiapkanProdukLama) ditulis ke tabel flat di transaksi
// pemanggil, lalu dibaca ulang dan dibandingkan. ID yang sudah ada di induk flat = ErrSudahDiFlat, nol tulisan.
func (g *Gudang) SalinProdukLama(ctx context.Context, tx *db.Tx, p models.Produk) error {
	if !tx.Terisi() {
		return errors.New("repository: copying an old product requires a transaction")
	}
	q, err := g.siapkan(TabelFlatInduk, sqlCacahID)
	if err != nil {
		return err
	}
	var c int
	if err := tx.QueryRowContext(ctx, q, p.ID).Scan(&c); err != nil {
		return fmt.Errorf("repository: checking %s %s: %w", TabelFlatInduk, p.ID, err)
	}
	if c > 0 {
		return fmt.Errorf("%w: %s", ErrSudahDiFlat, p.ID)
	}
	if err := g.tulisFlat(ctx, tx, p, true); err != nil {
		// Disalin pemakai lain di antara pemeriksaan dan penulisan: PK induk menolaknya.
		if strings.Contains(err.Error(), "ORA-00001") {
			return fmt.Errorf("%w: %s", ErrSudahDiFlat, p.ID)
		}
		return err
	}
	f, err := g.bacaFlat(ctx, tx, p.ID, false)
	if err != nil {
		return err
	}
	if beda := BedaProduk(p, f); len(beda) > 0 {
		return fmt.Errorf("%w: %s %s.%s", ErrBacaUlangBeda, p.ID, beda[0].Tabel, beda[0].Kolom)
	}
	return nil
}

// idLamaTerpakai - penerbitan ID produk baru: ID yang masih dipakai kedua tabel JSON warisan tidak diterbitkan ulang,
// karena produk lama itu masih dapat disalin (Copy Old / alat pindah) dengan ID-nya - tanpa ini produk baru dapat
// merebut nomornya. Tabel warisan yang sudah dibuang DBA bukan galat (ORA-00942 dicatat log). SELECT saja.
func (g *Gudang) idLamaTerpakai(ctx context.Context, tx *db.Tx, id string) (bool, error) {
	for _, tabel := range []string{TabelProduk, TabelInward} {
		q, err := g.siapkan(tabel, sqlCacahID)
		if err != nil {
			return false, err
		}
		var c int
		if err := tx.QueryRowContext(ctx, q, id).Scan(&c); err != nil {
			if strings.Contains(err.Error(), "ORA-00942") {
				// ⚠️ ORA-00942 juga berarti hak SELECT dicabut: dicatat, supaya pemeriksaan yang dilewati terlihat DBA.
				log.Printf("master product name life: %s not readable (%v); its IDs are not checked for %s", tabel, err, id)
				continue
			}
			return false, fmt.Errorf("repository: checking %s %s: %w", tabel, id, err)
		}
		if c > 0 {
			return true, nil
		}
	}
	return false, nil
}
