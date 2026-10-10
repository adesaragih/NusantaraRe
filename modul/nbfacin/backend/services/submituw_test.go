package services

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// ---- Validity Of Acceptance ------------------------------------------------------------------------------------

// TestHariBawaanNB - GetDateValidity_ACT langkah 2: Begindate lampau / kosong -> 0; selisih dipotong, +1, paling banyak 30.
func TestHariBawaanNB(t *testing.T) {
	now := time.Date(2026, 10, 7, 15, 4, 5, 0, time.UTC)
	for mulai, mau := range map[string]int{
		"": 0, "20261001T170000.000 GMT": 0, "bukan": 0,
		"20261017T170000.000 GMT": 11, // 9,08 hari -> 9 + 1? (17 Okt 17:00 GMT - 7 Okt 15:04 GMT = 10,08 hari) -> 10 + 1
		"20261008T150405.000 GMT": 2,  // tepat 1 hari -> 1 + 1
		"20261008T150404.000 GMT": 0,  // 0,99 hari -> 0 (tanpa +1)
		"20270101T170000.000 GMT": 30,
	} {
		if got := hariBawaanNB(now, mulai); got != mau {
			t.Errorf("%q: %d, mau %d", mulai, got, mau)
		}
	}
}

// TestHitungValiditas - DateValidity = hari ini WIB + hari (00:00 WIB, teks Pega UTC); WaitingBind L701 (bulan Calendar.add).
func TestHitungValiditas(t *testing.T) {
	now := time.Date(2026, 10, 7, 18, 0, 0, 0, time.UTC) // 8 Okt 01:00 WIB
	aksep := time.Date(2026, 1, 31, 10, 0, 0, 0, WIB)
	v := isianValiditas{sekarang: now, mulai: "20261019T170000.000 GMT", aksepTerakhir: aksep, adaAksepTerakhir: true}
	h := hitungValiditas(v, nil)
	// mulai 20 Okt 00:00 WIB - 8 Okt 01:00 WIB = 11,96 hari -> 11 + 1 = 12; 8 Okt + 12 = 20 Okt.
	if h.hari != "12" || h.tanggal != "20261019T170000.000 GMT" || h.waitingBind != "20261119T170000.000 GMT" {
		t.Errorf("bawaan (keUW kosong <= mulai -> Begin + 1 bulan): %+v", h)
	}
	n := 45
	v.keUW = "20261020T000000.000 GMT" // > mulai -> LastAksep + 1 bulan (31 Jan -> 28 Feb)
	if h := hitungValiditas(v, &n); h.hari != "45" || h.tanggal != "20261121T170000.000 GMT" || h.waitingBind != "20260227T170000.000 GMT" {
		t.Errorf("edit: %+v", h)
	}
	v.keUW = v.mulai // DateToUW == Begindate -> cabang "<=" (Begin + 1 bulan)
	if h := hitungValiditas(v, &n); h.waitingBind != "20261119T170000.000 GMT" {
		t.Errorf("keUW = mulai: %+v", h)
	}
	v.keUW, v.tba = "", true // IsTBA (tidak ada di sistem; jalur diuji): LastAksep + 3 bulan (30 Apr 00:00 WIB)
	if h := hitungValiditas(v, &n); h.waitingBind != "20260429T170000.000 GMT" {
		t.Errorf("tba: %+v", h)
	}
	v.tba, v.adaAksepTerakhir, v.keUW = false, false, "20991231T000000.000 GMT"
	if h := hitungValiditas(v, &n); h.waitingBind != "" {
		t.Errorf("tanpa riwayat: %+v", h)
	}
	v.keUW, v.mulai = "", ""
	if h := hitungValiditas(v, nil); h.hari != "0" || h.tanggal != "20261007T170000.000 GMT" || h.waitingBind != "" {
		t.Errorf("mulai kosong: %+v", h)
	}
}

// TestUraiHariValiditas - "" mode bawaan; 0-9999; selain itu 400.
func TestUraiHariValiditas(t *testing.T) {
	for teks, mau := range map[string]int{"0": 0, " 7 ": 7, "9999": 9999} {
		if n, err := uraiHariValiditas(teks); err != nil || n == nil || *n != mau {
			t.Errorf("%q: %v %v", teks, n, err)
		}
	}
	if n, err := uraiHariValiditas(""); n != nil || err != nil {
		t.Errorf("kosong: %v %v", n, err)
	}
	for _, teks := range []string{"10000", "-1", "+1", "1.5", "a", "1e3"} {
		if _, err := uraiHariValiditas(teks); !errors.Is(err, ErrMasukanSubmit) {
			t.Errorf("%q: %v", teks, err)
		}
	}
}

// ---- jenjang Accept ---------------------------------------------------------------------------------------------

func barisLimit(tabel, tg, jab, bottom string) models.BarisLimitA {
	return models.BarisLimitA{Tabel: tabel, TeamGroup: tg, Jabatan: jab, LimitBottom: d(bottom)}
}

const tabelPRY = "M_LIMIT_PROPERTYY"

// limitUji - team 1 lengkap; team 3 tanpa UNDERWRITER (seperti DDL); tabel lain berisi jebakan.
var limitUji = []models.BarisLimitA{
	barisLimit(tabelPRY, "1", "JUW_B", "0"), barisLimit(tabelPRY, "1", "UNDERWRITER", "1"),
	barisLimit(tabelPRY, "1", "SENIORUW", "300"), barisLimit(tabelPRY, "1", "SENIORUW", "300"),
	barisLimit(tabelPRY, "1", "DEPHEADUNDERWRITER", "1000"), barisLimit(tabelPRY, "1", "KADIVFACULTATIVE", "5000"),
	barisLimit(tabelPRY, "3", "JUW_B", "0"), barisLimit(tabelPRY, "3", "SENIORUW", "1"), barisLimit(tabelPRY, "3", "DEPHEADUNDERWRITER", "255"),
	barisLimit("M_LIMIT_PROPERTY_NON_PREFERREDD", "1", "SENIORUW", "2"),
	barisLimit(tabelPRY, "2", "UNDERWRITER", "2"),
}

// TestJabatanBerikutUW - GetLimitAkseptasiJUWA_SQL: jenjang LIMIT_BOTTOM terkecil di atas jabatan sendiri dan di bawah TSI,
// per tabel + team; jabatan sendiri tidak ada / ganda beda -> ""; seri -> urutan tabel; batas TSI eksklusif.
func TestJabatanBerikutUW(t *testing.T) {
	for _, c := range []struct {
		tabel, tg, jab, tsi, mau string
	}{
		{tabelPRY, "1", "JUW_B", "500", "UNDERWRITER"},
		{tabelPRY, "1", "UNDERWRITER", "500", "SENIORUW"},
		{tabelPRY, "1", "UNDERWRITER", "300", ""}, // LIMIT_BOTTOM < CARID2 eksklusif
		{tabelPRY, "1", "UNDERWRITER", "1001", "SENIORUW"},
		{tabelPRY, "1", "JUW_B", "1", ""},           // TSI di bawah semua jenjang -> Else
		{tabelPRY, "3", "JUW_B", "500", "SENIORUW"}, // team tanpa UNDERWRITER -> SENIORUW (keputusan work owner)
		{tabelPRY, "1", "JUW_A", "500", ""},         // DEV: tanpa baris JUW_A -> Else
		{"M_LIMIT_PROPERTY_NON_PREFERREDD", "1", "JUW_B", "500", ""},
		{tabelPRY, " 1 ", "JUW_B", "500", "UNDERWRITER"},
	} {
		if got, _ := jabatanBerikutUW(limitUji, c.tabel, c.tg, c.jab, d(c.tsi)); got != c.mau {
			t.Errorf("%+v: %q", c, got)
		}
	}
	// TeamGroup kosong -> Else, walau ada baris TEAM_GROUP NULL yang "cocok".
	null := append(append([]models.BarisLimitA{}, limitUji...), barisLimit(tabelPRY, "", "JUW_B", "0"), barisLimit(tabelPRY, "", "UNDERWRITER", "1"))
	if got, alasan := jabatanBerikutUW(null, tabelPRY, " ", "JUW_B", d("500")); got != "" || alasan != "TeamGroup kosong" {
		t.Errorf("TeamGroup kosong: %q %q", got, alasan)
	}
	ganda := append(append([]models.BarisLimitA{}, limitUji...), barisLimit(tabelPRY, "1", "JUW_B", "5"))
	if got, alasan := jabatanBerikutUW(ganda, tabelPRY, "1", "JUW_B", d("500")); got != "" || !strings.Contains(alasan, "ganda") {
		t.Errorf("ganda beda: %q %q", got, alasan)
	}
	seri := []models.BarisLimitA{barisLimit(tabelPRY, "1", "JUW_B", "0"), barisLimit(tabelPRY, "1", "B", "7"), barisLimit(tabelPRY, "1", "A", "7")}
	if got, _ := jabatanBerikutUW(seri, tabelPRY, "1", "JUW_B", d("500")); got != "B" {
		t.Errorf("seri: %q", got)
	}
}

// ---- Submit tahap uw --------------------------------------------------------------------------------------------

// spreadingSah - kasusSatuCoverage PS 50 dengan satu baris FACOUT 100 % bernilai Save (lolos M1-M5); TotalTSINusaRe 500.
func spreadingSah() models.KasusSpreading {
	k := kasusSatuCoverage("")
	k.Lokasi[0].Items[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: "10015", TreatyName: "FACOUT", SharePercentage: d("100")}}
	hitungSaveSpreading(&k, d("50"))
	return k
}

// layananSubmitUW - kasus UJI-NB-1 di posisi `posisi` (flag 0), spreading sp, tabel limit limitUji, riwayat terakhir
// 06-10-2026 10:00 WIB, Begin date 20 Okt 2026.
func layananSubmitUW(st *submitTiruan, posisi string, sp models.KasusSpreading) *Service {
	k := st.konteks["UJI-NB-1"]
	k.Position, k.CreateOp, k.StartDateTime = posisi, "UJI-PEMBUAT", "20261019T170000.000 GMT"
	st.konteks["UJI-NB-1"] = k
	st.riwayat = []models.RiwayatAkseptasi{{Tanggal: "06-10-2026 10:00", PIC: "UJI", Approval: "Accept"}}
	svc := layananSubmit(st).DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{sp}, idUang: map[string]string{"IDR": idMataUangIDR}})
	svc.limit = limitTiruan{a: limitUji}
	kasus := &models.Kasus{CaseID: "UJI-NB-1", Opportunity: models.Opportunity{InsuredID: "UJI-ORG", ClassOfBusiness: "UJI KELAS",
		GroupBusinessID: "1"}, General: models.General{MarketingID: "UJI-MO"}, RiskScoring: models.RingkasanScoring{FinalScore: d("61")}}
	return svc.DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": kasus}}).
		DenganTableOfLimit(tolTiruan{oldID: []string{oldIDUji}})
}

// oldIDUji - BUSINESS.OLDID tiruan di luar daftar When ISLIMIT* (bukan bentuk B).
const oldIDUji = "UJI-OLD"

// denganOldID - ganti BUSINESS.OLDID tiruan (nil = 0 baris).
func denganOldID(svc *Service, old []string) *Service {
	return svc.DenganTableOfLimit(tolTiruan{oldID: old})
}

func pelakuPosisi(p string) inti.Pelaku { return inti.Pelaku{AkunID: "UJI-UW", Peran: []string{p}} }

// TestAjukanSubmitUW - matriks posisi × Email Type × hasil jenjang (tabel tiruan): posisi / NB status / tulisan.
func TestAjukanSubmitUW(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		posisi, email      string
		ubahTabel          func(*Service)
		mauPosisi, mauNB   string
		perf, perfLanjut   string // StatusDoc / "" = tidak
		statusWork, isiTip string
	}{
		{posisiJuniorUWB, "1", nil, posisiUWUtama, "NB IS IN UNDERWRITER 'S INBOX", "ACCEPT", "OUTSTANDING", "", ""},
		{posisiUWUtama, "1", nil, posisiSeniorUW, "NB IS IN SENIOR UNDERWRITER'S INBOX", "ACCEPT", "OUTSTANDING", "", ""},
		{posisiJuniorUWB, "1", func(s *Service) {
			s.limit = limitTiruan{a: []models.BarisLimitA{barisLimit(tabelPRY, "1", "JUW_B", "0"), barisLimit(tabelPRY, "1", "JUW_A", "2"),
				barisLimit(tabelPRY, "1", "DEPHEADUNDERWRITER", "3")}}
		}, posisiJuniorUWA, "NB IS IN JUNIOR UNDERWRITER (A)'S INBOX", "ACCEPT", "OUTSTANDING", "", ""},
		{posisiJuniorUWB, "1", func(s *Service) {
			s.limit = limitTiruan{a: []models.BarisLimitA{barisLimit(tabelPRY, "1", "JUW_B", "0"), barisLimit(tabelPRY, "1", "DEPHEADUNDERWRITER", "3")}}
		}, posisiDepHeadUW, "NB IS IN DEPHEAD UNDERWRITER'S INBOX", "ACCEPT", "OUTSTANDING", "", ""},
		{posisiJuniorUWA, "1", func(s *Service) {
			s.limit = limitTiruan{a: []models.BarisLimitA{barisLimit(tabelPRY, "1", "JUW_A", "0"), barisLimit(tabelPRY, "1", "UNDERWRITER", "2")}}
		}, posisiUWUtama, "NB IS IN UNDERWRITER 'S INBOX", "ACCEPT", "OUTSTANDING", "", ""},
		{posisiJuniorUWB, "2", nil, "ReasFacInMarketing", "NB IS IN Uji - Pemasar'S INBOX", "", "", "", ""},
		{posisiUWUtama, "2", nil, "ReasFacInMarketing", "NB IS IN Uji - Pemasar'S INBOX", "", "", "", ""},
		{"ReasFacInUnderwritingFinancial", "7", nil, posisiDecline, "NB WAS DECLINED BY  UJI PENGGUNA", "DECLINE", "", statusDecline, ""},
		{posisiSeniorUW, "3", nil, "ReasFacInMarketing", "NB IS IN Uji - Pemasar'S INBOX", "", "", "", ""},
		{posisiJuniorUWA, "7", nil, posisiDecline, "NB WAS DECLINED BY  UJI PENGGUNA", "DECLINE", "", statusDecline, ""},
		{posisiDepHeadUW, "7", nil, posisiDecline, "NB WAS DECLINED BY  UJI PENGGUNA", "DECLINE", "", statusDecline, ""},
	} {
		st := submitUji()
		svc := layananSubmitUW(st, c.posisi, spreadingSah())
		if c.ubahTabel != nil {
			c.ubahTabel(svc)
		}
		h, err := svc.AjukanSubmit(ctx, pelakuPosisi(c.posisi), "UJI-NB-1", models.IsiSubmit{EmailType: c.email, Comment: " UJI komentar "})
		if err != nil || h.Position != c.mauPosisi || h.NBStatus != c.mauNB || len(st.tulis) != 1 {
			t.Fatalf("%s/%s: %v %+v", c.posisi, c.email, err, h)
		}
		w := st.tulis[0]
		if w.PosisiLama != c.posisi || w.Pega.Workbasket != c.posisi || w.StatusWorkBaru != c.statusWork || w.Riwayat.Keterangan != "UJI komentar" ||
			w.Validitas == nil || *w.Validitas != (models.ValiditasUW{Days: "13", Date: "20261019T170000.000 GMT",
			WaitingBind: "20261119T170000.000 GMT", Comment: "UJI komentar"}) {
			t.Errorf("%s/%s tulisan: %+v / %+v", c.posisi, c.email, w, w.Validitas)
		}
		statusPerf := func(f *models.PerformanceFacIn) string {
			if f == nil {
				return ""
			}
			return f.StatusDoc
		}
		if statusPerf(w.Performance) != c.perf || statusPerf(w.PerformanceLanjut) != c.perfLanjut {
			t.Errorf("%s/%s performance: %+v / %+v", c.posisi, c.email, w.Performance, w.PerformanceLanjut)
		}
		if w.PerformanceLanjut != nil && (*w.PerformanceLanjut != models.PerformanceFacIn{PositionDoc: "BINDING", StatusDoc: "OUTSTANDING",
			Username: "UJI-PEMBUAT", InsuredName: "UJI TERTANGGUNG", SOB: "UJI SOB", MOLeader: "UJI-LEADER", COB: "UJI COB"} ||
			w.Performance.PositionDoc != "OFFER" || w.Performance.Username != "UJI-UW") {
			t.Errorf("%s/%s baris performance: %+v / %+v", c.posisi, c.email, w.Performance, w.PerformanceLanjut)
		}
		if w.KeUW != "" || w.Deductible != nil {
			t.Errorf("%s/%s isian TL ikut tertulis", c.posisi, c.email)
		}
	}
}

// TestAjukanSubmitUWDitahan - Accept ke jalur akhir / limit tak terhitung -> 422 tanpa tulisan; isian 400; workbasket 403;
// flag 1 409; Validity edit.
func TestAjukanSubmitUWDitahan(t *testing.T) {
	ctx := context.Background()
	terima := models.IsiSubmit{EmailType: "1", Comment: "x"}
	for nama, c := range map[string]struct {
		posisi string
		pl     inti.Pelaku
		isi    models.IsiSubmit
		ubah   func(*submitTiruan, *Service)
		mau    error
		pesan  string
	}{
		// Tiket 62: Decision24 - spreadingSah memuat baris FACOUT (IsFacRetro) -> jalur FacRetro ditahan.
		"Senior FacRetro":       {posisiSeniorUW, pelakuPosisi(posisiSeniorUW), terima, nil, ErrProteksiSubmit, pesanFacRetro},
		"DepHead FacRetro":      {posisiDepHeadUW, pelakuPosisi(posisiDepHeadUW), terima, nil, ErrProteksiSubmit, pesanFacRetro},
		"Financial FacRetro":    {"ReasFacInUnderwritingFinancial", pelakuPosisi("ReasFacInUnderwritingFinancial"), terima, nil, ErrProteksiSubmit, pesanFacRetro},
		"JUW A tanpa baris":     {posisiJuniorUWA, pelakuPosisi(posisiJuniorUWA), terima, nil, ErrProteksiSubmit, pesanFacRetro},
		"preferred tak dikenal": {posisiJuniorUWB, pelakuPosisi(posisiJuniorUWB), terima, func(st *submitTiruan, _ *Service) { st.kategori = nil }, ErrProteksiSubmit, pesanFacRetro},
		"TSI di bawah jenjang": {posisiJuniorUWB, pelakuPosisi(posisiJuniorUWB), terima, func(_ *submitTiruan, s *Service) {
			sp := spreadingSah()
			hitungSaveSpreading(&sp, d("0.1")) // TotalTSINusaRe 1 -> tidak ada LIMIT_BOTTOM > 0 dan < 1 -> Decision24
			s.DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{sp}, idUang: map[string]string{"IDR": idMataUangIDR}})
		}, ErrProteksiSubmit, pesanFacRetro},
		"kurs tidak ada": {posisiJuniorUWB, pelakuPosisi(posisiJuniorUWB), terima, func(_ *submitTiruan, s *Service) {
			sp := spreadingSah()
			sp.Lokasi[0].Items[0].Currency = "USD"
			s.DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{sp}, idUang: map[string]string{"IDR": idMataUangIDR}})
		}, ErrProteksiSubmit, "Kurs USD ke IDR tidak ada di TREATYEXCHANGEYEARLY - limit akseptasi tidak dapat dihitung"},
		"proteksi spreading": {posisiJuniorUWB, pelakuPosisi(posisiJuniorUWB), terima, func(_ *submitTiruan, s *Service) {
			sp := spreadingSah()
			sp.Lokasi[0].Items[0].Coverages[0].Spreading = nil
			s.DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{sp}, idUang: map[string]string{"IDR": idMataUangIDR}})
		}, ErrProteksiSubmit, "Total premium spreaded on location 1 object item 1 coverage 1 is not equals premium RNM!"},
		"banding 400":       {posisiJuniorUWB, pelakuPosisi(posisiJuniorUWB), models.IsiSubmit{EmailType: "4", Comment: "x", BandingTo: "SENIORUW"}, nil, ErrMasukanSubmit, ""},
		"comment 400":       {posisiJuniorUWB, pelakuPosisi(posisiJuniorUWB), models.IsiSubmit{EmailType: "7", Comment: " "}, nil, ErrMasukanSubmit, ""},
		"comment lebar 400": {posisiJuniorUWB, pelakuPosisi(posisiJuniorUWB), models.IsiSubmit{EmailType: "2", Comment: strings.Repeat("é", 250) + "a"}, nil, ErrMasukanSubmit, ""},
		"hari 400":          {posisiJuniorUWB, pelakuPosisi(posisiJuniorUWB), models.IsiSubmit{EmailType: "2", Comment: "x", DaysValidity: "abc"}, nil, ErrMasukanSubmit, ""},
		"403 workbasket":    {posisiJuniorUWB, pelakuPosisi(posisiUWUtama), terima, nil, inti.ErrTanpaWewenang, ""},
		"409 flag 1": {posisiJuniorUWB, pelakuPosisi(posisiJuniorUWB), terima, func(st *submitTiruan, _ *Service) {
			k := st.konteks["UJI-NB-1"]
			k.FlagOnGoingPolicy = "1"
			st.konteks["UJI-NB-1"] = k
		}, ErrKasusBacaSaja, ""},
	} {
		st := submitUji()
		svc := layananSubmitUW(st, c.posisi, spreadingSah())
		if c.ubah != nil {
			c.ubah(st, svc)
		}
		_, err := svc.AjukanSubmit(ctx, c.pl, "UJI-NB-1", c.isi)
		var g GalatProteksi
		if !errors.Is(err, c.mau) || len(st.tulis) != 0 || (c.pesan != "" && (!errors.As(err, &g) || g.Pesan[len(g.Pesan)-1] != c.pesan)) {
			t.Errorf("%s: %v (tulis %d)", nama, err, len(st.tulis))
		}
	}
	// Comment tepat 500 bita lolos; daysValidity edit tersimpan.
	st := submitUji()
	if _, err := layananSubmitUW(st, posisiJuniorUWB, spreadingSah()).AjukanSubmit(ctx, pelakuPosisi(posisiJuniorUWB), "UJI-NB-1",
		models.IsiSubmit{EmailType: "2", Comment: strings.Repeat("é", 250), DaysValidity: "45"}); err != nil ||
		st.tulis[0].Validitas.Days != "45" || st.tulis[0].Validitas.Date != "20261120T170000.000 GMT" { // 7 Okt WIB + 45 = 21 Nov
		t.Errorf("edit: %v %+v", err, st.tulis)
	}
}

// TestProteksiSpreadingUW - CheckSpreadingProtect(Fire)_ACT M1-M5 verbatim, indeks mulai 1, urutan Pega; layering dilewati;
// premi dibandingkan 4 desimal.
func TestProteksiSpreadingUW(t *testing.T) {
	ctx := context.Background()
	periksa := func(ubah func(*models.KasusSpreading)) []string {
		sp := spreadingSah()
		ubah(&sp)
		svc := layananSubmitUW(submitUji(), posisiJuniorUWB, sp)
		p, err := svc.proteksiSpreadingUW(ctx, "UJI-NB-1")
		if err != nil {
			t.Fatal(err)
		}
		return p
	}
	if p := periksa(func(*models.KasusSpreading) {}); len(p) != 0 {
		t.Fatalf("sah: %q", p)
	}
	letak := "location 1 object item 1 coverage 1"
	for nama, c := range map[string]struct {
		ubah func(*models.KasusSpreading)
		mau  []string
	}{
		"kosong": {func(k *models.KasusSpreading) { k.Lokasi[0].Items[0].Coverages[0].Spreading = nil },
			[]string{"% Total spreading on " + letak + " must be 100%!", "Total premium spreaded on " + letak + " is not equals premium RNM!"}},
		"share kosong": {func(k *models.KasusSpreading) {
			c := &k.Lokasi[0].Items[0].Coverages[0]
			c.Spreading = append(c.Spreading, models.BarisSpreading{TreatyType: "10007"})
		}, []string{"% spreading on " + letak + " can't be empty!"}},
		"treaty kosong": {func(k *models.KasusSpreading) { k.Lokasi[0].Items[0].Coverages[0].Spreading[0].TreatyType = " " },
			[]string{"% spreading on " + letak + " can't be empty!", pesanJenisTreaty}},
		"premi beda": {func(k *models.KasusSpreading) {
			k.Lokasi[0].Items[0].Coverages[0].Spreading[0].PremiumSpreaded = d("3.99")
		},
			[]string{"Total premium spreaded on " + letak + " is not equals premium RNM!"}},
		"premi beda di desimal 5": {func(k *models.KasusSpreading) {
			k.Lokasi[0].Items[0].Coverages[0].Spreading[0].PremiumSpreaded = d("4.00004")
		}, nil},
		"treaty asing": {func(k *models.KasusSpreading) { k.Lokasi[0].Items[0].Coverages[0].Spreading[0].TreatyType = "X9" }, []string{pesanJenisTreaty}},
		"layering": {func(k *models.KasusSpreading) {
			c := &k.Lokasi[0].Items[0].Coverages[0]
			c.CoverageBasis, c.Spreading = "5", nil
		}, nil},
	} {
		if got := periksa(c.ubah); !reflect.DeepEqual(got, append([]string{}, c.mau...)) {
			t.Errorf("%s: %q", nama, got)
		}
	}
}

// TestPeriksaDanTampilanSubmitUW - periksa Accept = proteksi + jenjang; Reject kosong; GET tahap uw + validitas (bawaan /
// tersimpan); POST validitas (edit, penjaga).
func TestPeriksaDanTampilanSubmitUW(t *testing.T) {
	ctx := context.Background()
	st := submitUji()
	svc := layananSubmitUW(st, posisiSeniorUW, spreadingSah())
	pl := pelakuPosisi(posisiSeniorUW)
	if h, err := svc.PeriksaSubmit(ctx, pl, "UJI-NB-1", IsiPeriksa{EmailType: "1"}); err != nil || !reflect.DeepEqual(h.Pesan, []string{pesanFacRetro}) {
		t.Errorf("periksa Accept Senior: %v %q", err, h.Pesan)
	}
	if _, err := svc.PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", IsiPeriksa{EmailType: "1"}); !errors.Is(err, inti.ErrTanpaWewenang) {
		t.Errorf("periksa tanpa workbasket posisi: %v", err)
	}
	sp := spreadingSah()
	sp.Lokasi[0].Items[0].Coverages[0].Spreading[0].TreatyType = "X9"
	juwB := layananSubmitUW(submitUji(), posisiJuniorUWB, sp)
	if h, err := juwB.PeriksaSubmit(ctx, pelakuPosisi(posisiJuniorUWB), "UJI-NB-1", IsiPeriksa{EmailType: "1"}); err != nil ||
		!reflect.DeepEqual(h.Pesan, []string{pesanJenisTreaty}) {
		t.Errorf("periksa Accept JUW B (M5, jenjang lolos): %v %q", err, h.Pesan)
	}
	if h, err := svc.PeriksaSubmit(ctx, pl, "UJI-NB-1", IsiPeriksa{EmailType: "3"}); err != nil || len(h.Pesan) != 0 {
		t.Errorf("periksa Ask: %v %q", err, h.Pesan)
	}
	if _, err := svc.PeriksaSubmit(ctx, pl, "UJI-NB-1", IsiPeriksa{EmailType: "4"}); !errors.Is(err, ErrMasukanSubmit) {
		t.Errorf("periksa Banding: %v", err)
	}
	if _, err := svc.PeriksaSubmit(ctx, pl, "UJI-NB-1", IsiPeriksa{EmailType: "2", DaysValidity: "x"}); !errors.Is(err, ErrMasukanSubmit) {
		t.Errorf("periksa hari: %v", err)
	}
	v, err := svc.TampilanSubmit(ctx, pl, "UJI-NB-1")
	if err != nil || v.Tahap != TahapUW || !v.BolehSubmit || v.DaysValidity != "" || v.DateValidity != "" || v.WaitingBindDate != "" {
		t.Errorf("tampilan belum tersimpan (tanpa hitung): %v %+v", err, v)
	}
	// Mode bawaan lewat POST validitas daysValidity "" (klik radio Email Type).
	if r, err := svc.HitungValiditasSubmit(ctx, pl, "UJI-NB-1", ""); err != nil ||
		r != (ValiditasSubmit{DaysValidity: "13", DateValidity: "20-10-2026", WaitingBindDate: "20-11-2026"}) {
		t.Errorf("validitas bawaan: %v %+v", err, r)
	}
	k := st.konteks["UJI-NB-1"]
	k.DaysValidity, k.DateValidity, k.WaitingBindDate = "5", "20261011T170000.000 GMT", ""
	st.konteks["UJI-NB-1"] = k
	if v, _ := svc.TampilanSubmit(ctx, pl, "UJI-NB-1"); v.DaysValidity != "5" || v.DateValidity != "12-10-2026" || v.WaitingBindDate != "" {
		t.Errorf("tampilan tersimpan: %+v", v)
	}
	if r, err := svc.HitungValiditasSubmit(ctx, pl, "UJI-NB-1", "45"); err != nil ||
		r != (ValiditasSubmit{DaysValidity: "45", DateValidity: "21-11-2026", WaitingBindDate: "20-11-2026"}) {
		t.Errorf("validitas: %v %+v", err, r)
	}
	for nama, c := range map[string]struct {
		pl   inti.Pelaku
		hari string
		mau  error
	}{
		"rusak": {pl, "-1", ErrMasukanSubmit}, "huruf": {pl, "x", ErrMasukanSubmit}, "403": {pelakuMkt, "1", inti.ErrTanpaWewenang},
		"401": {inti.Pelaku{}, "1", inti.ErrTanpaIdentitas},
	} {
		if _, err := svc.HitungValiditasSubmit(ctx, c.pl, "UJI-NB-1", c.hari); !errors.Is(err, c.mau) {
			t.Errorf("validitas %s: %v", nama, err)
		}
	}
	k.FlagOnGoingPolicy = "1" // tiket 69a: UW flag 1 = tahap uwPolicy (Validity Of Acceptance tetap tampil)
	st.konteks["UJI-NB-1"] = k
	if v, _ := svc.TampilanSubmit(ctx, pl, "UJI-NB-1"); v.Tahap != TahapUWPolicy || !v.BolehSubmit || v.DaysValidity == "" {
		t.Errorf("tampilan flag 1: %+v", v)
	}
	k.FlagOnGoingPolicy = "2"
	st.konteks["UJI-NB-1"] = k
	if v, _ := svc.TampilanSubmit(ctx, pl, "UJI-NB-1"); v.Tahap != "" || v.BolehSubmit || v.DaysValidity != "" {
		t.Errorf("tampilan flag 2: %+v", v)
	}
	mkt := layananSubmit(submitUji())
	if _, err := mkt.HitungValiditasSubmit(ctx, pelakuMkt, "UJI-NB-1", "1"); !errors.Is(err, ErrKasusBacaSaja) {
		t.Errorf("validitas marketing: %v", err)
	}
	if v, _ := mkt.TampilanSubmit(ctx, pelakuMkt, "UJI-NB-1"); v.DaysValidity != "" || v.DateValidity != "" || v.WaitingBindDate != "" {
		t.Errorf("tampilan marketing: %+v", v)
	}
}

// TestCarid2UW - langkah 2 / 9: tanpa top risk = TotalTSINusaRe IDR (kurs TERBARU; IDR = 1); ada top risk = TotalTSITopRisk
// (Σ TSILiability coverage pertama item lokasi top risk × PS / 100, TANPA kurs - ikut Pega), dibulatkan 0 desimal.
func TestCarid2UW(t *testing.T) {
	ctx := context.Background()
	hitung := func(ubah func(*models.KasusSpreading), kurs map[string]*apd.Decimal) (string, string) {
		sp := spreadingSah() // PS 50, IDR, TSILiability 1000
		ubah(&sp)
		svc := layananSubmitUW(submitUji(), posisiJuniorUWB, sp)
		svc.DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{sp}, idUang: map[string]string{"IDR": idMataUangIDR, "USD": "10001"},
			kurs: kurs})
		c, pesan, err := svc.carid2UW(ctx, "UJI-NB-1")
		if err != nil {
			t.Fatal(err)
		}
		if c == nil {
			return "", pesan
		}
		return utils.FormatDecimal(c), pesan
	}
	if c, p := hitung(func(*models.KasusSpreading) {}, nil); c != "500" || p != "" {
		t.Errorf("IDR: %s %q", c, p)
	}
	usd := func(k *models.KasusSpreading) { k.Lokasi[0].Items[0].Currency = "USD" }
	if c, _ := hitung(usd, map[string]*apd.Decimal{"10001": d("15000.4")}); c != "7500200" {
		t.Errorf("USD kurs terbaru: %s", c)
	}
	// K73-33 (keputusan work owner 09-10-2026): TotalTSITopRisk DENGAN kurs (CountTotalTSIPremiNusaRe_Act L9060 / L9254); semula
	// tiket 61 tanpa kurs (500).
	topUSD := func(k *models.KasusSpreading) { usd(k); k.Lokasi[0].IsTopRisk = true }
	if c, p := hitung(topUSD, map[string]*apd.Decimal{"10001": d("15000.4")}); c != "7500200" || p != "" {
		t.Errorf("top risk dengan kurs: %s %q", c, p)
	}
	if c, p := hitung(topUSD, nil); c != "" || !strings.HasPrefix(p, "Kurs USD ke IDR tidak ada") {
		t.Errorf("top risk tanpa kurs tersedia: %s %q", c, p)
	}
	if c, _ := hitung(func(k *models.KasusSpreading) { k.PercentShare = d("33.333") }, nil); c != "333" {
		t.Errorf("bulat 0: %s", c)
	}
}

// ---- tiket 62: jenjang lanjutan (Decision24 / GetLimitAkseptasi_ActFlow / Decision23) ---------------------------

// spreadingTanpaFacRetro - spreadingSah dengan ORS (10007) menggantikan FACOUT: IsFacRetro tidak, lolos M1-M5; TSILiability
// `tsi` (PS 50 -> TotalTSINusaRe = tsi / 2).
func spreadingTanpaFacRetro(tsi string) models.KasusSpreading {
	k := kasusSatuCoverage("")
	c := &k.Lokasi[0].Items[0].Coverages[0]
	c.TSILiability, c.TSI = d(tsi), d(tsi)
	c.Spreading = []models.BarisSpreading{{TreatyType: "10007", TreatyName: "ORS", SharePercentage: d("100")}}
	c.OldID = "FLEXAS" // sejalan dengan coverage lokasiSah (FACINOFFER memasangkan objek dan spreading, tiket 63)
	hitungSaveSpreading(&k, d("50"))
	return k
}

// limitLanjut - tabel PRY team 1 untuk tangga lanjutan (LIMIT_BOTTOM berbeda semua).
func limitLanjut(tambahan ...models.BarisLimitA) []models.BarisLimitA {
	b := []models.BarisLimitA{barisLimit(tabelPRY, "1", "JUW_B", "0"), barisLimit(tabelPRY, "1", "UNDERWRITER", "1"),
		barisLimit(tabelPRY, "1", "SENIORUW", "300"), barisLimit(tabelPRY, "1", "DEPHEADUNDERWRITER", "1000"),
		barisLimit(tabelPRY, "1", "MANAGERTEKNIK", "2000"), barisLimit(tabelPRY, "1", "KADIVFACULTATIVE", "2500"),
		barisLimit(tabelPRY, "1", "KADIVTEKNIK", "3000"), barisLimit(tabelPRY, "1", "DIREKTURTEKNIK", "4000")}
	return append(b, tambahan...)
}

// nbBinding - Transition38: prefix + " IS IN " + MarketingName (apa adanya) + "'S INBOX".
const nbBinding = "NB IS IN Uji - Pemasar'S INBOX"

// TestAcceptJenjangLanjut - matriks posisi × hasil tangga: Decision24 (Senior / DepHead / Financial), langsung ActFlow (A1 /
// A8 / A11 / A6 / A4 / A5), langkah 18 membuang KADIVFACULTATIVE / MANAGERTEKNIK, LOOP dua putaran, jabatan pengguna dari
// POSITION (butir 33), Else / jabatan tanpa baris / seri -> ditahan; NB status K62-2.
func TestAcceptJenjangLanjut(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		nama, posisi, tsi string
		baris             []models.BarisLimitA
		mauPosisi, mauNB  string
		pesan             string
		oldID             []string // nil = oldIDUji
	}{
		{"Senior -> DepHead", posisiSeniorUW, "4000", limitLanjut(), posisiDepHeadUW, "NB IS IN DEPHEAD UNDERWRITER'S INBOX", "", nil},
		// Putaran 2 (LOOP): DepHead -> MANAGERTEKNIK 2000 & KADIVFACULTATIVE 2500 dibuang langkah 18 -> KADIVTEKNIK.
		{"DepHead -> Kadiv Teknik", posisiDepHeadUW, "12000", limitLanjut(), posisiKadivTeknik, "NB IS IN KADIV TEKNIK'S INBOX", "", nil},
		{"Kadiv Teknik -> Direktur Teknik", posisiKadivTeknik, "12000", limitLanjut(), posisiDirTeknik, "NB IS IN DIREKTUR TEKNIK'S INBOX", "", nil},
		// Tiket 63: Decision23 Else -> SAVE JSON_OFFER -> Transition38 (marketing Binding).
		{"Direktur Teknik -> Binding", posisiDirTeknik, "12000", limitLanjut(), repository.WorkbasketMarketing, nbBinding, "", nil},
		{"Senior TSI kecil -> Binding", posisiSeniorUW, "700", limitLanjut(), repository.WorkbasketMarketing, nbBinding, "", nil},
		{"Financial tanpa jabatan -> Binding", posisiUWFinancial, "12000", limitLanjut(), repository.WorkbasketMarketing, nbBinding, "", nil},
		{"Kadiv Fin tanpa baris sendiri -> Binding", posisiKadivFin, "12000", limitLanjut(), repository.WorkbasketMarketing, nbBinding, "", nil},
		{"Kadiv Fac -> Kadiv Fin", posisiKadivFac, "12000", limitLanjut(barisLimit(tabelPRY, "1", "KADIVFINANCIAL", "2700")), posisiKadivFin,
			"NB IS IN KADIV FINANCIAL'S INBOX", "", nil},
		{"Kadiv Fin -> Direktur Marketing", posisiKadivFin, "12000", limitLanjut(barisLimit(tabelPRY, "1", "KADIVFINANCIAL", "2700"),
			barisLimit(tabelPRY, "1", "DIREKTURMARKETING", "2800")), posisiDirMarketing, "NB IS IN DIREKTUR MARKETING'S INBOX", "", nil},
		{"Direktur Marketing -> Kadiv Teknik", posisiDirMarketing, "12000", limitLanjut(barisLimit(tabelPRY, "1", "DIREKTURMARKETING", "2800")),
			posisiKadivTeknik, "NB IS IN KADIV TEKNIK'S INBOX", "", nil},
		{"Manager Teknik -> Kadiv Teknik", posisiManagerTeknik, "12000", limitLanjut(), posisiKadivTeknik, "NB IS IN KADIV TEKNIK'S INBOX", "", nil},
		{"seri -> ditahan", posisiSeniorUW, "12000", limitLanjut(barisLimit(tabelPRY, "1", "DIREKTURMARKETING", "1000")), "", "",
			pesanTanggaTakPasti, nil},
		// Decision42 Else (TSI di bawah UNDERWRITER) -> Decision24 tanpa FacRetro -> ActFlow JUW_B -> Decision23 Else -> Binding.
		{"JUW B Decision42 Else -> Binding", posisiJuniorUWB, "1", limitLanjut(), repository.WorkbasketMarketing, nbBinding, "", nil},
		{"JUW B Decision42 tetap", posisiJuniorUWB, "4000", limitLanjut(), posisiUWUtama, "NB IS IN UNDERWRITER 'S INBOX", "", nil},
		// BusinessOldId: When IsLimitCreditCL (C2) -> bentuk B -> NextFinancial (K73-34; semula ditahan): FINANCIALINS tanpa baris ->
		// LetterNo "" -> Decision23 Else -> Binding. Tidak diketahui (0 / 2 baris) -> ErrTafsirBerbeda -> ditahan.
		{"bentuk B FINANCIALINS kosong -> Binding", posisiSeniorUW, "4000", limitLanjut(), repository.WorkbasketMarketing, nbBinding, "", []string{"C2"}},
		{"OLDID 0 baris -> ditahan", posisiSeniorUW, "4000", limitLanjut(), "", "", pesanTanggaTakPasti, []string{}},
		{"OLDID ganda -> ditahan", posisiKadivTeknik, "12000", limitLanjut(), "", "", pesanTanggaTakPasti, []string{"UJI-OLD", "UJI-OLD2"}},
	} {
		st := submitUji()
		svc := layananSubmitUW(st, c.posisi, spreadingTanpaFacRetro(c.tsi))
		if c.oldID != nil {
			svc = denganOldID(svc, c.oldID)
		}
		svc.limit = limitTiruan{a: c.baris}
		h, err := svc.AjukanSubmit(ctx, pelakuPosisi(c.posisi), "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "x"})
		if c.pesan != "" {
			var g GalatProteksi
			if !errors.As(err, &g) || g.Pesan[len(g.Pesan)-1] != c.pesan || len(st.tulis) != 0 {
				t.Errorf("%s: %v (tulis %d)", c.nama, err, len(st.tulis))
			}
			continue
		}
		if err != nil || h.Position != c.mauPosisi || h.NBStatus != c.mauNB || len(st.tulis) != 1 || st.tulis[0].Pega.Workbasket != c.posisi {
			t.Errorf("%s: %v %+v", c.nama, err, h)
			continue
		}
		if b := st.tulis[0].Binding; (b != nil) != (c.mauPosisi == repository.WorkbasketMarketing) || (b != nil && len(b.Offer) != 1) {
			t.Errorf("%s: binding %+v", c.nama, b)
		}
	}
}

// TestSubmitJenjangLanjut - tahap uw meliputi enam posisi baru (pemegang workbasket, flag 0): Reject / Ask / Decline seperti
// UW (Decline juga di A11 - penyimpangan sadar); Spreading / Akumulasi TETAP hanya IsUWSpread.
func TestSubmitJenjangLanjut(t *testing.T) {
	ctx := context.Background()
	for _, pos := range []string{posisiKadivTeknik, posisiKadivFac, posisiKadivFin, posisiDirMarketing, posisiDirTeknik, posisiManagerTeknik} {
		tk := models.TahapKasus{Position: pos, FlagOnGoingPolicy: "0"}
		if !tahapSubmitUW(tk) || tahapUW(tk) || tahapSubmitUW(models.TahapKasus{Position: pos, FlagOnGoingPolicy: "1"}) {
			t.Errorf("%s: tahap submit / spreading salah", pos)
		}
		for email, mau := range map[string]string{"2": "ReasFacInMarketing", "3": "ReasFacInMarketing", "7": posisiDecline} {
			st := submitUji()
			h, err := layananSubmitUW(st, pos, spreadingTanpaFacRetro("1000")).AjukanSubmit(ctx, pelakuPosisi(pos), "UJI-NB-1",
				models.IsiSubmit{EmailType: email, Comment: "x"})
			if err != nil || h.Position != mau || st.tulis[0].Pega.Workbasket != pos {
				t.Errorf("%s/%s: %v %+v", pos, email, err, h)
			}
		}
		if _, err := layananSubmitUW(submitUji(), pos, spreadingTanpaFacRetro("1000")).AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1",
			models.IsiSubmit{EmailType: "2", Comment: "x"}); !errors.Is(err, inti.ErrTanpaWewenang) {
			t.Errorf("%s 403: %v", pos, err)
		}
	}
}

// TestAcceptUWGroup - ralat tiket 61 (keputusan work owner 08-10-2026 "Selaraskan, ikut Pega"): GetLimitAkseptasi_JUW_UW
// langkah 8 When IsGroup -> LetterNo "" -> Decision42 Else -> Decision24 (FacRetro ditahan) / ActFlow (langkah 3 exit group) ->
// Decision23 Else -> Decision16 IsGroup -> jalur Policy group ditahan (pesanPolicyGroup, tiket 64). Dulu JUW B / UW group naik jenjang.
func TestAcceptUWGroup(t *testing.T) {
	ctx := context.Background()
	grup := KonfigGrup{Login: []string{"UJI-UW"}}
	for _, c := range []struct {
		nama, posisi string
		sp           models.KasusSpreading
		g            KonfigGrup
		mauPosisi    string
		pesan        string
	}{
		{"JUW B bukan group", posisiJuniorUWB, spreadingTanpaFacRetro("4000"), KonfigGrup{}, posisiUWUtama, ""},
		{"UW bukan group", posisiUWUtama, spreadingTanpaFacRetro("4000"), KonfigGrup{}, posisiSeniorUW, ""},
		{"JUW B group", posisiJuniorUWB, spreadingTanpaFacRetro("4000"), grup, "", pesanPolicyGroup},
		{"UW group", posisiUWUtama, spreadingTanpaFacRetro("4000"), grup, "", pesanPolicyGroup},
		{"JUW B group pembuat case", posisiJuniorUWB, spreadingTanpaFacRetro("4000"), KonfigGrup{Login: []string{"UJI-PEMBUAT"}}, "",
			pesanPolicyGroup},
		{"JUW B group FacRetro", posisiJuniorUWB, spreadingSah(), grup, "", pesanFacRetro},
	} {
		st := submitUji()
		svc := layananSubmitUW(st, c.posisi, c.sp).DenganGrup(c.g, false)
		svc.limit = limitTiruan{a: limitLanjut()}
		h, err := svc.AjukanSubmit(ctx, pelakuPosisi(c.posisi), "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "x"})
		if c.pesan != "" {
			var g GalatProteksi
			if !errors.As(err, &g) || g.Pesan[len(g.Pesan)-1] != c.pesan || len(st.tulis) != 0 {
				t.Errorf("%s: %v (tulis %d)", c.nama, err, len(st.tulis))
			}
			continue
		}
		if err != nil || h.Position != c.mauPosisi || len(st.tulis) != 1 {
			t.Errorf("%s: %v %+v", c.nama, err, h)
		}
	}
}

// TestAcceptBentukBFinancial - K73-34: Bond / Kredit di Senior ke atas memakai M_LIMIT_FINANCIALINS (NextFinancial): Kadiv
// Financial, Kredit CL (BusinessOldId C2), baris "DIREKTUR MARKETING" ber-limit <= CARID2 -> Direktur Marketing.
func TestAcceptBentukBFinancial(t *testing.T) {
	st := submitUji()
	svc := denganOldID(layananSubmitUW(st, posisiKadivFin, spreadingTanpaFacRetro("12000")), []string{"C2"})
	svc.limit = limitTiruan{a: limitLanjut(), b: []models.BarisLimitB{{Jabatan: "DIREKTUR MARKETING", LimitCreditCL: d("100")}}}
	h, err := svc.AjukanSubmit(context.Background(), pelakuPosisi(posisiKadivFin), "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "x"})
	if err != nil || h.Position != posisiDirMarketing || h.NBStatus != "NB IS IN DIREKTUR MARKETING'S INBOX" {
		t.Errorf("%v %+v", err, h)
	}
}
