package services

// Uji SusunHasilPeriksa - irisan I2. Data UJI buatan: nomor polis sintetis (tanpa pola nomor polis produksi).

import (
	"errors"
	"testing"
	"time"

	"nusantarare/modul/endorsmentfacin/backend/models"
)

const (
	polisUji     = "UJI.01.0001"
	polisUjiLife = "UJI.RNML.0001"
)

// dataBersih - polis FIRE tanpa satu pun alasan tolak; tanggal belum diisi.
func dataBersih() DataPeriksa {
	return DataPeriksa{PolicyNo: polisUji, EdmType: models.EdmPerubahan, EdmTypeNew: "1", JenisBisnis: "Fire"}
}

func tanggalUji(t *testing.T, s string) *time.Time {
	t.Helper()
	w, err := time.ParseInLocation(formatKabel, s, jakarta)
	if err != nil {
		t.Fatal(err)
	}
	return &w
}

func susun(t *testing.T, d DataPeriksa) models.HasilPeriksa {
	t.Helper()
	h, err := SusunHasilPeriksa(d)
	if err != nil {
		t.Fatalf("SusunHasilPeriksa: %v", err)
	}
	return h
}

func hasilGerbang(h models.HasilPeriksa) []models.HasilGerbangPeriksa {
	out := make([]models.HasilGerbangPeriksa, len(h.Gerbang))
	for i, g := range h.Gerbang {
		out[i] = g.Hasil
	}
	return out
}

func teksPesan(h models.HasilPeriksa) []string {
	out := []string{}
	for _, p := range h.Pesan {
		out = append(out, p.Medan+"|"+p.Teks)
	}
	return out
}

func samaIrisan[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

const (
	L  = models.GerbangLolos
	T  = models.GerbangTolak
	D  = models.GerbangDilewati
	BD = models.GerbangBelumDiperiksa
)

func TestPeriksaBersih(t *testing.T) {
	h := susun(t, dataBersih())
	if h.Ditolak || h.Lengkap || len(h.Pesan) != 0 {
		t.Fatalf("ditolak=%v lengkap=%v pesan=%v", h.Ditolak, h.Lengkap, h.Pesan)
	}
	// Gerbang 2 selalu belum-diperiksa (GetListEdm tanpa body), gerbang 6 dilewati (bukan retro).
	if got := hasilGerbang(h); !samaIrisan(got, []models.HasilGerbangPeriksa{L, BD, L, L, L, D}) {
		t.Fatalf("gerbang = %v", got)
	}
	for i, g := range h.Gerbang {
		if g.No != i+1 {
			t.Fatalf("gerbang ke-%d bernomor %d", i, g.No)
		}
	}
	if h.Bisnis != (models.BisnisPeriksa{Hasil: L, JenisBisnis: "Fire"}) || h.Tanggal != nil {
		t.Fatalf("bisnis=%+v tanggal=%+v", h.Bisnis, h.Tanggal)
	}
}

func TestPeriksaGerbangSatuStatusEDM(t *testing.T) {
	for _, k := range []struct {
		status string
		tolak  bool
	}{{"1", true}, {"2", true}, {"3", false}, {"4", false}, {"", false}} {
		d := dataBersih()
		d.StatusEDM = k.status
		h := susun(t, d)
		want := L
		var pesan []string
		if k.tolak {
			want, pesan = T, []string{"PolicyNo|Sudah Di endorsement Batal"}
		}
		if h.Gerbang[0].Hasil != want || !samaIrisan(teksPesan(h), append([]string{}, pesan...)) || h.Ditolak != k.tolak {
			t.Errorf("status %q: gerbang1=%v pesan=%v", k.status, h.Gerbang[0].Hasil, teksPesan(h))
		}
		if h.Gerbang[0].Catatan != "" {
			t.Errorf("status %q: catatan %q", k.status, h.Gerbang[0].Catatan)
		}
	}
	// Teks "1.0" terbaca angka 1 tetapi teksnya berbeda → galat, tidak ditebak.
	d := dataBersih()
	d.StatusEDM = "1.0"
	if _, err := SusunHasilPeriksa(d); !errors.Is(err, ErrTafsirKodeBerbeda) {
		t.Fatalf("galat = %v", err)
	}
}

func TestPeriksaRISlipMelompatiGerbang(t *testing.T) {
	d := dataBersih()
	d.EdmType, d.EdmTypeNew = models.EdmPerubahan, models.EdmTypeNewRISlip
	d.StatusEDM, d.Klaim, d.Renewal, d.EndorsementInternalRetro = "1", []string{"0"}, true, "1"
	h := susun(t, d)
	if got := hasilGerbang(h); !samaIrisan(got, []models.HasilGerbangPeriksa{D, BD, D, D, D, D}) {
		t.Fatalf("gerbang = %v", got)
	}
	if h.Ditolak {
		t.Fatalf("RI slip ditolak: %v", teksPesan(h))
	}
	// EdmTypeNew 4 tanpa EdmType 4 BUKAN RI slip.
	d.EdmType = models.EdmBatalSejakSemula
	if h := susun(t, d); h.Gerbang[0].Hasil != T {
		t.Fatalf("EdmType 1 + EdmTypeNew 4 dianggap RI slip: %v", hasilGerbang(h))
	}
}

func TestPeriksaGerbangTigaKlaim(t *testing.T) {
	for _, k := range []struct {
		klaim []string
		want  models.HasilGerbangPeriksa
	}{{nil, L}, {[]string{"2"}, L}, {[]string{"0"}, T}, {[]string{""}, T}, {[]string{"1"}, T}} {
		d := dataBersih()
		d.Klaim = k.klaim
		h := susun(t, d)
		if h.Gerbang[2].Hasil != k.want {
			t.Errorf("klaim %v: gerbang3 = %v", k.klaim, h.Gerbang[2].Hasil)
		}
		if (k.want == T) != samaIrisan(teksPesan(h), []string{"PolicyNo|There's already a claim with this policy no"}) {
			t.Errorf("klaim %v: pesan = %v", k.klaim, teksPesan(h))
		}
	}
}

func TestPeriksaGerbangEmpatPembayaran(t *testing.T) {
	benar, salah := true, false
	for _, k := range []struct {
		edm        models.JenisEndorsemen
		pembayaran *bool
		want       models.HasilGerbangPeriksa
	}{
		{models.EdmBatalSejakSemula, &benar, T},
		{models.EdmBatalProrata, &benar, T},
		{models.EdmBatalSejakSemula, &salah, L},
		{models.EdmBatalSejakSemula, nil, BD},
		{models.EdmPerubahan, &benar, L},
		{"", &benar, L},
		{"", nil, L},
	} {
		d := dataBersih()
		d.EdmType, d.Pembayaran = k.edm, k.pembayaran
		h := susun(t, d)
		if h.Gerbang[3].Hasil != k.want {
			t.Errorf("EdmType %q pembayaran %v: gerbang4 = %v", k.edm, k.pembayaran, h.Gerbang[3].Hasil)
		}
		if (k.want == T) != h.Ditolak {
			t.Errorf("EdmType %q: ditolak = %v (%v)", k.edm, h.Ditolak, teksPesan(h))
		}
		if k.want == BD && (h.Gerbang[3].Catatan != catatanArasapas || h.Lengkap) {
			t.Errorf("belum-diperiksa: catatan %q lengkap %v", h.Gerbang[3].Catatan, h.Lengkap)
		}
	}
}

func TestPeriksaGerbangLimaEnam(t *testing.T) {
	d := dataBersih()
	d.Renewal = true
	if h := susun(t, d); h.Gerbang[4].Hasil != T || !samaIrisan(teksPesan(h), []string{"PolicyNo|There's RNW with this policy no!"}) {
		t.Fatalf("renewal: %v %v", h.Gerbang[4].Hasil, teksPesan(h))
	}
	for _, k := range []struct {
		retro  string
		facout bool
		want   models.HasilGerbangPeriksa
	}{{"1", false, T}, {"1", true, L}, {"", false, D}, {"0", false, D}} {
		d := dataBersih()
		d.EndorsementInternalRetro, d.Facout = k.retro, k.facout
		h := susun(t, d)
		if h.Gerbang[5].Hasil != k.want {
			t.Errorf("retro %q facout %v: gerbang6 = %v", k.retro, k.facout, h.Gerbang[5].Hasil)
		}
		if (k.want == T) != samaIrisan(teksPesan(h), []string{"PolicyNo|This policy is not spreading FACOUT"}) {
			t.Errorf("retro %q: pesan = %v", k.retro, teksPesan(h))
		}
	}
}

// Klep OPENPROTEKSI_EDM membuka tepat gerbangnya sendiri (1 Batal, 2 Klaim → gerbang 3, 3 Pembayaran → gerbang 4,
// 4 Renewal → gerbang 5), hanya bila syarat tolaknya benar.
func TestPeriksaKlep(t *testing.T) {
	benar := true
	tolakSemua := func() DataPeriksa {
		d := dataBersih()
		d.EdmType, d.StatusEDM, d.Klaim, d.Pembayaran, d.Renewal = models.EdmBatalSejakSemula, "1", []string{"0"}, &benar, true
		return d
	}
	if got := hasilGerbang(susun(t, tolakSemua())); !samaIrisan(got, []models.HasilGerbangPeriksa{T, BD, T, T, T, D}) {
		t.Fatalf("tanpa klep: %v", got)
	}
	for klep, gerbang := range map[KlepPembatal]int{KlepBatal: 1, KlepKlaim: 3, KlepPembayaran: 4, KlepRenewal: 5} {
		d := tolakSemua()
		d.Klep = map[KlepPembatal]bool{klep: true}
		h := susun(t, d)
		for i, g := range h.Gerbang {
			want := []models.HasilGerbangPeriksa{T, BD, T, T, T, D}[i]
			if i+1 == gerbang {
				want = models.GerbangDiklep
			}
			if g.Hasil != want {
				t.Errorf("klep %d: gerbang %d = %v, ingin %v", klep, i+1, g.Hasil, want)
			}
		}
		if h.Gerbang[gerbang-1].Catatan != catatanDiklep || len(h.Pesan) != 3 {
			t.Errorf("klep %d: catatan %q, %d pesan %v", klep, h.Gerbang[gerbang-1].Catatan, len(h.Pesan), teksPesan(h))
		}
	}
	// Seluruh klep terisi: nol pesan gerbang, tidak ditolak.
	d := tolakSemua()
	d.Klep = map[KlepPembatal]bool{KlepBatal: true, KlepKlaim: true, KlepPembayaran: true, KlepRenewal: true}
	if h := susun(t, d); h.Ditolak || !samaIrisan(hasilGerbang(h), []models.HasilGerbangPeriksa{
		models.GerbangDiklep, BD, models.GerbangDiklep, models.GerbangDiklep, models.GerbangDiklep, D}) {
		t.Fatalf("semua klep: %v %v", hasilGerbang(h), teksPesan(h))
	}
	// Klep tanpa syarat tolak: tetap lolos.
	d = dataBersih()
	d.Klep = map[KlepPembatal]bool{KlepBatal: true, KlepKlaim: true, KlepPembayaran: true, KlepRenewal: true}
	if got := hasilGerbang(susun(t, d)); !samaIrisan(got, []models.HasilGerbangPeriksa{L, BD, L, L, L, D}) {
		t.Fatalf("klep tanpa syarat: %v", got)
	}
	// Gerbang 6 tanpa klep.
	d = dataBersih()
	d.EndorsementInternalRetro, d.Klep = "1", map[KlepPembatal]bool{KlepBatal: true, KlepKlaim: true, KlepPembayaran: true, KlepRenewal: true}
	if h := susun(t, d); h.Gerbang[5].Hasil != T {
		t.Fatalf("gerbang 6 dibuka klep: %v", h.Gerbang[5])
	}
	// Tabel tidak ada: tolak + catatan.
	d = tolakSemua()
	d.KlepTakTerbaca = true
	h := susun(t, d)
	for _, no := range []int{1, 3, 4, 5} {
		if h.Gerbang[no-1].Hasil != T || h.Gerbang[no-1].Catatan != catatanKlepTakTerbaca {
			t.Errorf("klep tak terbaca, gerbang %d = %+v", no, h.Gerbang[no-1])
		}
	}
}

func TestPeriksaBisnis(t *testing.T) {
	for _, jb := range []string{"KPR", "Fire", "Liability", "MBUCar", "MarineCargo", "PA", "GolfInsurance"} {
		d := dataBersih()
		d.JenisBisnis = jb
		if h := susun(t, d); h.Bisnis.Hasil != L || h.Ditolak {
			t.Errorf("%s: bisnis %v pesan %v", jb, h.Bisnis.Hasil, teksPesan(h))
		}
	}
	d := dataBersih()
	d.JenisBisnis = "Bukan Lini"
	if h := susun(t, d); h.Bisnis.Hasil != T || !samaIrisan(teksPesan(h), []string{"PolicyNo|Invalid business, please contact IT!"}) {
		t.Fatalf("lini tak dikenal: %v %v", h.Bisnis, teksPesan(h))
	}
	d.JenisBisnis = ""
	if h := susun(t, d); h.Bisnis.Hasil != T || !samaIrisan(teksPesan(h), []string{"PolicyNo|Invalid policy no!"}) {
		t.Fatalf("tanpa jenis bisnis: %v %v", h.Bisnis, teksPesan(h))
	}
	// Langkah 8: polis berpola RNML = "Life" walau query kosong. [dugaan] IsLife EDM membaca
	// .Quotation.BusinessOldId kasus portal, yang tidak diisi rule mana pun → salah → langkah 17 menolak.
	d.PolicyNo = polisUjiLife
	h := susun(t, d)
	if h.Bisnis != (models.BisnisPeriksa{Hasil: T, JenisBisnis: "Life"}) ||
		!samaIrisan(teksPesan(h), []string{"PolicyNo|Invalid business, please contact IT!"}) {
		t.Fatalf("Life: %+v %v", h.Bisnis, teksPesan(h))
	}
}

// Baris tulisan NB FacIn: jam dinding 12:00 WIB (= 05:00 GMT, tanggal sama).
func TestPeriksaTanggal(t *testing.T) {
	periode := &PeriodeMentah{Mulai: "20250101120000", Akhir: "20260101120000"}
	for _, k := range []struct {
		tanggal string
		want    models.HasilGerbangPeriksa
	}{
		{"01-01-2025", L}, // = tanggal mulai (CARI20)
		{"31-12-2024", T}, // sehari sebelum mulai
		{"15-06-2025", L},
		{"01-01-2026", L}, // = tanggal akhir: tidak SESUDAH akhir
		{"02-01-2026", T}, // sehari sesudah akhir
	} {
		d := dataBersih()
		d.TanggalEndorsement, d.Periode = tanggalUji(t, k.tanggal), periode
		h := susun(t, d)
		if h.Tanggal == nil || h.Tanggal.Hasil != k.want {
			t.Fatalf("%s: tanggal = %+v", k.tanggal, h.Tanggal)
		}
		if h.Tanggal.Mulai != "01-01-2025" || h.Tanggal.Akhir != "01-01-2026" {
			t.Errorf("%s: periode keluar %s s/d %s", k.tanggal, h.Tanggal.Mulai, h.Tanggal.Akhir)
		}
		ditolak := samaIrisan(teksPesan(h), []string{"EndorsementDate|EDM date cannot be outside the period"})
		if (k.want == T) != ditolak || (k.want == T) != h.Ditolak {
			t.Errorf("%s: pesan = %v", k.tanggal, teksPesan(h))
		}
	}
}

// Baris berjam 00:00 WIB = 17:00 GMT hari sebelumnya: batas yang dibandingkan mundur sehari (K-046 apa adanya,
// keputusan work owner "zona waktu samain kaya NB facin").
func TestPeriksaTanggalTengahMalamWIB(t *testing.T) {
	periode := &PeriodeMentah{Mulai: "20250101000000", Akhir: "20260101000000"}
	for _, k := range []struct {
		tanggal string
		want    models.HasilGerbangPeriksa
	}{
		{"30-12-2024", T},
		{"31-12-2024", L}, // = batas mulai GMT
		{"31-12-2025", L}, // = batas akhir GMT
		{"01-01-2026", T}, // tanggal akhir WIB tersimpan, tetapi sesudah batas GMT
	} {
		d := dataBersih()
		d.TanggalEndorsement, d.Periode = tanggalUji(t, k.tanggal), periode
		h := susun(t, d)
		if h.Tanggal.Hasil != k.want {
			t.Errorf("%s: %+v", k.tanggal, h.Tanggal)
		}
		if h.Tanggal.Mulai != "31-12-2024" || h.Tanggal.Akhir != "31-12-2025" {
			t.Errorf("%s: periode keluar %s s/d %s", k.tanggal, h.Tanggal.Mulai, h.Tanggal.Akhir)
		}
	}
}

func TestPeriksaTanggalPengecualianDanRNML(t *testing.T) {
	luar := tanggalUji(t, "01-01-2030")
	periode := &PeriodeMentah{Mulai: "20250101000000", Akhir: "20260101000000"}

	d := dataBersih()
	d.TanggalEndorsement, d.Periode, d.PengecualianTanggal = luar, periode, DaftarPolis{polisUji}
	if h := susun(t, d); h.Tanggal.Hasil != D || h.Ditolak {
		t.Fatalf("pengecualian: %+v %v", h.Tanggal, teksPesan(h))
	}
	// Daftar lain tidak berlaku.
	d.PengecualianTanggal = DaftarPolis{"UJI.99.9999"}
	if h := susun(t, d); h.Tanggal.Hasil != T {
		t.Fatalf("pengecualian polis lain: %+v", h.Tanggal)
	}
	// Tabel 261 belum ada: daftar kosong + catatan.
	d.PengecualianTanggal, d.PengecualianBelumAda = nil, true
	if h := susun(t, d); h.Tanggal.Hasil != T || h.Tanggal.Catatan != catatanPengecualianBelumAda {
		t.Fatalf("tabel belum ada: %+v", h.Tanggal)
	}
	// RNML dilewati tanpa periode.
	d = dataBersih()
	d.PolicyNo, d.TanggalEndorsement = polisUjiLife, luar
	if h := susun(t, d); h.Tanggal.Hasil != D {
		t.Fatalf("RNML: %+v", h.Tanggal)
	}
}

func TestPeriksaTanggalPeriodeTidakAda(t *testing.T) {
	for _, p := range []*PeriodeMentah{nil, {Mulai: "", Akhir: "20260101000000"}, {Mulai: "20250101000000"}} {
		d := dataBersih()
		d.TanggalEndorsement, d.Periode = tanggalUji(t, "01-06-2025"), p
		h := susun(t, d)
		if h.Tanggal.Hasil != T || !h.Ditolak ||
			!samaIrisan(teksPesan(h), []string{"EndorsementDate|" + pesanPeriodeTidakAda}) {
			t.Errorf("periode %+v: %+v %v", p, h.Tanggal, teksPesan(h))
		}
	}
	d := dataBersih()
	d.TanggalEndorsement, d.Periode = tanggalUji(t, "01-06-2025"), &PeriodeMentah{Mulai: "2025-01-01", Akhir: "20260101000000"}
	if _, err := SusunHasilPeriksa(d); !errors.Is(err, ErrDateOracleTakTerurai) {
		t.Fatalf("teks DATE rusak: %v", err)
	}
}

// Urutan pesan: gerbang (urut langkah), bisnis, lalu tanggal.
func TestPeriksaUrutanPesan(t *testing.T) {
	benar := true
	d := DataPeriksa{
		PolicyNo: polisUji, EdmType: models.EdmBatalSejakSemula, EndorsementInternalRetro: "1",
		StatusEDM: "1", Klaim: []string{"0"}, Pembayaran: &benar, Renewal: true, JenisBisnis: "",
		TanggalEndorsement: tanggalUji(t, "01-01-2030"),
		Periode:            &PeriodeMentah{Mulai: "20250101000000", Akhir: "20260101000000"},
	}
	want := []string{
		"PolicyNo|Sudah Di endorsement Batal",
		"PolicyNo|There's already a claim with this policy no",
		"PolicyNo|There's already payment with this policy no",
		"PolicyNo|There's RNW with this policy no!",
		"PolicyNo|This policy is not spreading FACOUT",
		"PolicyNo|Invalid policy no!",
		"EndorsementDate|EDM date cannot be outside the period",
	}
	if got := teksPesan(susun(t, d)); !samaIrisan(got, want) {
		t.Fatalf("pesan =\n%v\nwant\n%v", got, want)
	}
}

// Rincian per gerbang SEPAKAT dengan pesan PeriksaGerbangPenolakan atas seluruh kombinasi masukan gerbang.
func TestRincianSepakatDenganPesanGerbang(t *testing.T) {
	benar, salah := true, false
	n := 0
	for _, edm := range []models.JenisEndorsemen{"", "1", "2", "4"} {
		for _, baru := range []string{"", "4"} {
			for _, status := range []string{"", "1", "3"} {
				for _, klaim := range [][]string{nil, {"2"}, {"0"}} {
					for _, bayar := range []*bool{nil, &benar, &salah} {
						for _, rnw := range []bool{false, true} {
							for _, retro := range []string{"", "1"} {
								for _, fo := range []bool{false, true} {
									// Klep: kosong, ganjil (1 dan 3), genap (2 dan 4).
									for _, klep := range []map[KlepPembatal]bool{nil, {KlepBatal: true, KlepPembayaran: true},
										{KlepKlaim: true, KlepRenewal: true}} {
										d := DataPeriksa{PolicyNo: polisUji, EdmType: edm, EdmTypeNew: baru, StatusEDM: status,
											Klaim: klaim, Pembayaran: bayar, Renewal: rnw, EndorsementInternalRetro: retro,
											Facout: fo, JenisBisnis: "Fire", Klep: klep}
										h, err := SusunHasilPeriksa(d)
										if err != nil {
											t.Fatalf("%+v: %v", d, err)
										}
										tolak := 0
										for _, g := range h.Gerbang {
											if g.Hasil == T {
												tolak++
											}
										}
										if tolak != len(h.Pesan) {
											t.Fatalf("%+v: %d gerbang tolak, %d pesan %v", d, tolak, len(h.Pesan), teksPesan(h))
										}
										n++
									}
								}
							}
						}
					}
				}
			}
		}
	}
	if n != 4*2*3*3*3*2*2*2*3 {
		t.Fatalf("kombinasi %d", n)
	}
}
