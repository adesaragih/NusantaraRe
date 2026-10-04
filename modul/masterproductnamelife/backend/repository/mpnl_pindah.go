package repository

// Pindah JSON → flat dengan rekonsiliasi (tiket 01 bab 02-10-2026, T2/T3; alat `backend/alat/pindahflat`).
//
// ⛔ BUKAN berkas migrasi SQL: di skema uji tabel warisan dibuat SESUDAH migrasi, dan teks JSON path tidak lolos
// penjaga inti. Pengurainya kodek yang sudah teruji (`UraiProduk`), pemetaannya `NormalkanFlat`.
//
//	-uji       (bawaan) hanya SELECT atas kedua tabel JSON; nol tulisan; laporan.
//	-terima-normalisasi  (bawaan mati) normalisasi teks yang SUDAH diputuskan work owner (OQ-FLAT-07) tidak menahan
//	           -jalankan; tetap dicetak laporan. Kegagalan tidak pernah dapat diterima.
//	-jalankan  menolak IS_PEGA_PROD=true dan rekonsiliasi yang tidak lolos; SATU transaksi: isi tabel flat dihapus lalu
//	           diisi ulang, dibaca ulang dan dibandingkan sebelum ditutup. Aman diulang - dan MENOLAK bila tabel flat
//	           memuat produk yang berbeda dari sumber JSON (tulisan baru sesudah peralihan tidak pernah ditimpa).
//
// Rekonsiliasi - SETIAP medan setiap produk, teks demi teks (JSON → model → bentuk flat): beda yang bukan K3/K4 =
// GAGAL; normalisasi teks (mis. nol ekor angka) dicatat per kolom dan menghentikan `-jalankan` untuk keputusan work
// owner. ⛔ Laporan AGREGAT: ID produk + nama kolom/kunci, tidak pernah nilainya.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

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
	TerimaNormalisasi   bool // -terima-normalisasi (OQ-FLAT-07)
	Ditulis             bool
	ContohNormalisasi   []Temuan
}

// Lolos - nol kegagalan DAN nol normalisasi (normalisasi menunggu keputusan work owner, brief T3).
func (l LaporanPindah) Lolos() bool { return len(l.Gagal) == 0 && len(l.Normalisasi) == 0 }

// BolehDitulis - nol kegagalan, dan normalisasi nol ATAU diterima operator sesudah keputusan work owner (OQ-FLAT-07).
func (l LaporanPindah) BolehDitulis() bool {
	return len(l.Gagal) == 0 && (len(l.Normalisasi) == 0 || l.TerimaNormalisasi)
}

func laporanBaru() LaporanPindah {
	return LaporanPindah{Baris: map[string]int{}, Normalisasi: map[string]int{}, MedanTanpaKolom: map[string]int{},
		DibuangD2: map[string]int{}}
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

// k3Dibolehkan - K3 keputusan work owner 02-10-2026: hanya dua bentuk nilai tidak sah yang di-NULL-kan.
func k3Dibolehkan(m MasalahNilai) bool {
	return (m.Tabel == TabelFlatInduk && m.Kolom == "MATURE" && m.Jenis == MasalahBukanTanggal) ||
		(m.Tabel == TabelFlatUWLimit && m.Kolom == "MAXINSURED" && m.Jenis == MasalahBukanAngka)
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
			if b.id == u.id || productIDDari(b) == u.id {
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
		hitungDibuang(&lap, u.isi, in.isi, p)
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
		for _, k := range MedanTanpaKolom(p) {
			lap.MedanTanpaKolom[k]++
			gagal("%s: %s terisi tetapi tidak punya kolom flat", u.id, k)
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
			lap.Normalisasi[b.Tabel+"."+b.Kolom+" ("+JenisNormalisasi(nilaiMedan(p, b))+")"]++
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

// hitungDibuang - yang tidak pindah karena D2 (bukan data): kunci internal Pega halaman dan baris, `IsView`, `Comment`.
func hitungDibuang(lap *LaporanPindah, jsonUmum, jsonInward string, p models.Produk) {
	if obj, err := uraiObjek(jsonUmum); err == nil {
		for k := range obj {
			if !kunciUmumDikelola[k] {
				lap.DibuangD2["kunci halaman umum "+k]++
			}
		}
		if _, ada := obj[kunciIsView]; ada {
			lap.DibuangD2["IsView (keadaan layar)"]++
		}
	}
	if obj, err := uraiObjek(jsonInward); err == nil {
		for k := range obj {
			if !kunciInwardDikelola[k] {
				lap.DibuangD2["kunci halaman inward "+k]++
			}
		}
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
		if len(c.asli) > 0 {
			lap.DibuangD2["baris "+c.daftar+" berkunci internal Pega"] += len(c.asli)
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
	tulisPeta(&b, "  normalisasi teks (menunggu keputusan work owner)", l.Normalisasi)
	for _, t := range l.ContohNormalisasi {
		fmt.Fprintf(&b, "    contoh: %s\n", t)
	}
	tulisPeta(&b, "  medan tanpa kolom yang terisi", l.MedanTanpaKolom)
	tulisPeta(&b, "  dibuang D2 (bukan data)", l.DibuangD2)
	if l.TulisanFlatBerbeda > 0 {
		fmt.Fprintf(&b, "  tabel flat sudah memuat %d produk yang berbeda dari sumber JSON\n", l.TulisanFlatBerbeda)
	}
	fmt.Fprintf(&b, "  gagal: %d\n", len(l.Gagal))
	for _, g := range l.Gagal {
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

func sqlHapusSemua(tabel string) string { return fmt.Sprintf(`DELETE FROM %s`, tabel) }

// bacaSumberJSON - seluruh baris kedua tabel JSON warisan (SELECT saja).
func (g *Gudang) bacaSumberJSON(ctx context.Context) (umum, inward []barisJSON, err error) {
	for _, s := range []struct {
		tabel string
		ke    *[]barisJSON
	}{{TabelProduk, &umum}, {TabelInward, &inward}} {
		q, err := g.siapkan(s.tabel, sqlSemuaJSON)
		if err != nil {
			return nil, nil, err
		}
		baris, err := g.bacaBaris(ctx, nil, q)
		if err != nil {
			return nil, nil, fmt.Errorf("repository: reading %s: %w", s.tabel, err)
		}
		*s.ke = baris
	}
	return umum, inward, nil
}

// PindahFlat - satu putaran alat pindah (`jalankan` = `-jalankan`, selain itu `-uji`; `terimaNormalisasi` =
// `-terima-normalisasi`).
func (g *Gudang) PindahFlat(ctx context.Context, jalankan, terimaNormalisasi bool) (LaporanPindah, error) {
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
	lap.Mode, lap.TerimaNormalisasi = mode, terimaNormalisasi
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
			lap.TulisanFlatBerbeda++
			continue
		}
		f, err := g.bacaFlat(ctx, tx, id, true)
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
	// Anak lebih dulu (FK), lalu induk.
	for i := len(DaftarTabelFlat) - 1; i >= 0; i-- {
		if err := g.eksekusi(ctx, tx, DaftarTabelFlat[i], sqlHapusSemua, false); err != nil {
			return lap, err
		}
	}
	for _, p := range produk {
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
