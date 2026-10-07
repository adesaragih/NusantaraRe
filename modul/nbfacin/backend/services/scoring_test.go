package services

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// pilihanUrut - 24 nomor ChechBox urut models.DaftarFaktorScoring -> peta pilihan.
func pilihanUrut(n ...int) map[string]int {
	m := map[string]int{}
	for i, f := range models.DaftarFaktorScoring {
		m[f.Kunci] = n[i]
	}
	return m
}

// TestPoinScoringPolaSampel - INSTRUMEN: pola centang lima kasus `DDL\CONTOH` ber-FinalScore (hanya nomor centang, tanpa
// data kasus) dengan Score dan persetujuan yang TERSIMPAN Pega sebagai jawaban. Occupation / TSI dipaksa dari lokasi (KDRisk
// / TSI IDR), band loss ratio dari persen yang jatuh di band tersimpan.
func TestPoinScoringPolaSampel(t *testing.T) {
	besar, kecil := apd.New(40000000000, 0), apd.New(1000000, 0)
	for _, k := range []struct {
		nama          string
		pilihan       []int
		kd            string
		tsi, share    *apd.Decimal
		persen        string
		score         string
		due, approval string // Occupation/ObjectConditions/LossRatio/Others, OpDiv/OpDir/Tech/Pres sebagai T/F
		layak         string // AcceptableNotApproval / AcceptableWithApproval / NotAcceptable
	}{
		{"EDM-13470", []int{1, 2, 2, 4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 2, 1, 2}, "01", besar, kecil, "", "98.5", "FFFF", "FFFF", "TFF"},
		{"EDM-13508", []int{1, 2, 2, 5, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 5, 5, 3, 6, 1, 1, 2, 3, 2}, "01", besar, kecil, "20", "61.0", "FFFF", "FFFF", "TFF"},
		{"EDM-13513", []int{4, 2, 2, 5, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 8, 3, 1, 1, 1, 2, 1, 2}, "04", besar, kecil, "80", "35.5", "TFTF", "TTTF", "FTF"},
		// Pega menyimpan OpDirector / TechDirector "true" karena `TotalRNMShare.CARI1 > "30000000000"` dibandingkan TEKS;
		// numerik (keputusan work owner): premi spread 1.000.000 -> false, 40.000.000.000 -> true (baris berikut).
		{"NB-181227 share kecil", []int{4, 2, 2, 5, 2, 1, 1, 1, 3, 3, 1, 1, 3, 1, 1, 2, 1, 2, 2, 1, 1, 2, 1, 2}, "04", besar, kecil, "", "69.5", "TFFF", "TFFF", "FTF"},
		{"NB-181227 share besar", []int{4, 2, 2, 5, 2, 1, 1, 1, 3, 3, 1, 1, 3, 1, 1, 2, 1, 2, 2, 1, 1, 2, 1, 2}, "04", besar, besar, "", "69.5", "TFFF", "TTTF", "FTF"},
		// Pega 69.0 = skor basi (TSI dicentang ulang ke 1 saat muat, ScorePerFactor masih +1); dihitung ulang: 67.0.
		{"NB-181231 dihitung ulang", []int{4, 2, 2, 6, 2, 1, 1, 1, 3, 3, 1, 1, 3, 1, 1, 2, 1, 2, 2, 1, 1, 1, 1, 2}, "04", kecil, kecil, "", "67.0", "TFFT", "TTTF", "FTF"},
	} {
		m := models.KolomScoring{Pilihan: pilihanUrut(k.pilihan...)}
		if k.persen != "" {
			m.PersenLossRatio = ds(k.persen)
		}
		lok := &lokasiScoring{OccupationName: "UJI OKUPASI", KDRisk: k.kd, TSIIDR: k.tsi, ShareIDR: k.share}
		sc, pesan := hitungScoring("New Business", [2]models.KolomScoring{m}, [2]*lokasiScoring{lok})
		got := sc.Kolom[0]
		for i, f := range models.DaftarFaktorScoring {
			if got.Pilihan[f.Kunci] != k.pilihan[i] {
				t.Errorf("%s: pilihan %s %d, mau %d", k.nama, f.Kunci, got.Pilihan[f.Kunci], k.pilihan[i])
			}
		}
		h := sc.Hasil
		tf := func(b ...bool) string {
			s := ""
			for _, x := range b {
				s += map[bool]string{true: "T", false: "F"}[x]
			}
			return s
		}
		if got.Total.String() != k.score || sc.FinalScore.String() != k.score || len(pesan) != 0 ||
			tf(h.Occupation, h.ObjectConditions, h.LossRatio, h.Others) != k.due ||
			tf(h.OperationalDivHead, h.OperationalDirector, h.TechnicalDirector, h.PresidentDirector) != k.approval ||
			tf(h.AcceptableNotApproval, h.AcceptableWithApproval, h.NotAcceptable) != k.layak {
			t.Errorf("%s: score %s final %s due %s approval %s layak %s pesan %v", k.nama, got.Total, sc.FinalScore,
				tf(h.Occupation, h.ObjectConditions, h.LossRatio, h.Others),
				tf(h.OperationalDivHead, h.OperationalDirector, h.TechnicalDirector, h.PresidentDirector),
				tf(h.AcceptableNotApproval, h.AcceptableWithApproval, h.NotAcceptable), pesan)
		}
	}
}

// TestPilihanLossRatio - SetLossRatio_Act: (0,10] 4, (10,25] 5, (25,50] 6, (50,75] 7, (75,100] 8, > 100 9; tanpa persen
// pilihan 4-9 dikosongkan, 1-3 / 10 tetap.
func TestPilihanLossRatio(t *testing.T) {
	for _, k := range []struct {
		persen    string
		lama, mau int
	}{{"0.01", 1, 4}, {"10", 0, 4}, {"10.5", 0, 5}, {"25", 0, 5}, {"50", 0, 6}, {"75", 0, 7}, {"100", 0, 8}, {"100.01", 0, 9},
		{"", 5, 0}, {"0", 9, 0}, {"-3", 3, 3}, {"", 1, 1}, {"", 10, 10}} {
		var p *apd.Decimal
		if k.persen != "" {
			p = ds(k.persen)
		}
		if got := pilihanLossRatio(p, k.lama); got != k.mau {
			t.Errorf("persen %q lama %d: %d, mau %d", k.persen, k.lama, got, k.mau)
		}
	}
}

// TestHitungKolomOtomatis - okupasi / TSI dikunci (dipaksa dari lokasi, termasuk tanpa lokasi -> 0), zona gempa hanya bila
// belum dipilih ("Daerah Zona n" -> n+1, zona lain -> 0), kurs tidak ada -> pesan + TSI 0, skor tiap faktor teks.
func TestHitungKolomOtomatis(t *testing.T) {
	var pesan []string
	m := models.KolomScoring{Pilihan: map[string]int{"occupation": 1, "totalSumInsured": 2, "earthquake": 0}}
	k := hitungKolom(m, &lokasiScoring{OccupationName: "UJI", OccupationID: "UJI-OKP-1", Alamat: "UJI ALAMAT", KDRisk: "05",
		Zona: "Daerah Zona 3", TSIIDR: apd.New(30000000000, 0)}, &pesan, "Top Risk")
	if k.Pilihan["occupation"] != 5 || k.Pilihan["totalSumInsured"] != 1 || k.Pilihan["earthquake"] != 4 || len(pesan) != 0 ||
		k.OccupationNote != "UJI" || k.OccupationCode != "UJI-OKP-1" || k.TopRiskLocation != "UJI ALAMAT" ||
		k.Skor["occupation"].String() != "-15" || k.Skor["earthquake"].String() != "2" || k.Skor["riot"].String() != "0" ||
		k.Total.String() != "-14" {
		t.Errorf("otomatis: %+v %v", k, pesan)
	}
	k = hitungKolom(models.KolomScoring{Pilihan: map[string]int{"earthquake": 7}}, &lokasiScoring{KDRisk: "09",
		Zona: "Daerah Zona 1", KursTidakAda: []string{"UJI-USD"}}, &pesan, "Dominant Risk")
	if k.Pilihan["occupation"] != 0 || k.Pilihan["totalSumInsured"] != 0 || k.Pilihan["earthquake"] != 7 || len(pesan) != 1 ||
		pesan[0] != "Kurs UJI-USD pada Begin date tidak ada di TREATYEXCHANGEYEARLY - Total Sum Insured Dominant Risk tidak diisi" {
		t.Errorf("kurs / pilihan pengguna: %+v %v", k.Pilihan, pesan)
	}
	k = hitungKolom(models.KolomScoring{Pilihan: map[string]int{"occupation": 2, "totalSumInsured": 1}}, nil, &pesan, "Top Risk")
	if k.Pilihan["occupation"] != 0 || k.Pilihan["totalSumInsured"] != 0 || k.OccupationNote != "" {
		t.Errorf("tanpa lokasi: %+v", k.Pilihan)
	}
	if k := hitungKolom(models.KolomScoring{}, &lokasiScoring{Zona: "Zona Lain", TSIIDR: apd.New(1, 0)}, &pesan, ""); k.Pilihan["earthquake"] != 0 {
		t.Errorf("zona tak dikenal: %d", k.Pilihan["earthquake"])
	}
}

// TestFinalScoreDuaKolom - Dominant ber-okupasi: rata-rata 2 desimal setengah ke atas; okupasi sama -> pesan verbatim;
// lokasi dominan tanpa okupasi -> kolom 2 dibuang (SetScore_Act langkah 3).
func TestFinalScoreDuaKolom(t *testing.T) {
	top := &lokasiScoring{OccupationName: "UJI A", KDRisk: "01", TSIIDR: apd.New(1, 0)}             // 15 + (-1) = 14
	dom := &lokasiScoring{OccupationName: "UJI B", KDRisk: "02", TSIIDR: apd.New(1, 0)}             // 10 - 1 + 2.25 = 11.25
	m := [2]models.KolomScoring{{}, {IndexLocation: "2", Pilihan: map[string]int{"earthquake": 3}}} // (14 + 11.25) / 2 = 12.625
	sc, pesan := hitungScoring("Renewal", m, [2]*lokasiScoring{top, dom})
	if !sc.Kolom[1].Ada || sc.FinalScore.String() != "12.63" || sc.NoteFinalScore != "Poor" || !sc.Hasil.Poor ||
		!sc.Hasil.NotAcceptable || len(pesan) != 0 || sc.Status != "Renewal" {
		t.Errorf("dua kolom: final %s %s %+v %v", sc.FinalScore, sc.NoteFinalScore, sc.Hasil, pesan)
	}
	sama := &lokasiScoring{OccupationName: "UJI A", KDRisk: "01", TSIIDR: apd.New(1, 0)}
	if _, pesan := hitungScoring("", m, [2]*lokasiScoring{top, sama}); !reflect.DeepEqual(pesan, []string{PesanOkupasiSama}) {
		t.Errorf("okupasi sama: %v", pesan)
	}
	sc, _ = hitungScoring("", m, [2]*lokasiScoring{top, {KDRisk: "02"}})
	if sc.Kolom[1].Ada || sc.FinalScore.String() != "14" {
		t.Errorf("dominan tanpa okupasi: %+v %s", sc.Kolom[1], sc.FinalScore)
	}
}

// TestKelasFinalScore - batas kelas ScoringResult langkah 1 dan 3 (100.5 = di atas semua kelas -> "Poor", tanpa kelas).
func TestKelasFinalScore(t *testing.T) {
	for _, k := range []struct{ fs, catatan, layak string }{{"100", "Good", "T"}, {"86", "Good", "T"}, {"85.99", "Above Average", "T"},
		{"76", "Above Average", "T"}, {"60", "Average", "T"}, {"30", "Below Average", "T"}, {"29.99", "Poor", "N"}, {"0", "Poor", "N"},
		{"-1", "Poor", "-"}, {"100.5", "Poor", "-"}} {
		h, catatan := hasilScoring(ds(k.fs), [2]models.KolomScoring{}, [2]*lokasiScoring{})
		layak := "-"
		if h.AcceptableNotApproval {
			layak = "T"
		} else if h.NotAcceptable {
			layak = "N"
		}
		if catatan != k.catatan || layak != k.layak || h.AcceptableWithApproval {
			t.Errorf("fs %s: %s %s %+v", k.fs, catatan, layak, h)
		}
	}
}

// ---- alur dengan tiruan -------------------------------------------------------------------------------------------

type scoringTiruan struct {
	k      models.KasusScoring
	tulis  []models.Scoring
	galatB error
}

func (s *scoringTiruan) BacaScoring(context.Context, string) (models.KasusScoring, error) {
	return s.k, s.galatB
}

func (s *scoringTiruan) TulisScoring(_ context.Context, _ *db.Tx, _ string, sc models.Scoring) error {
	s.tulis = append(s.tulis, sc)
	return nil
}

type okupasiScoringTiruan struct{ occupationTiruan }

func (okupasiScoringTiruan) RisikoOkupasiFire(_ context.Context, kode []string) (map[string]string, error) {
	out := map[string]string{}
	for _, k := range kode {
		out[k] = map[string]string{"UJI-OKP-A": "01", "UJI-OKP-B": "04"}[k]
	}
	return out, nil
}

func layananScoring(st *scoringTiruan) *Service {
	kasus := kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": {CaseID: "UJI-NB-1", InsuredName: "UJI TERTANGGUNG",
		General: models.General{StartDateTime: "20250101T170000.000 GMT", EndDateTime: "20260101T170000.000 GMT"}}}}
	objek := objekTiruan{ada: map[string][]models.ObjekFire{"UJI-NB-1": {
		{ObjectNo: "1", RiskLocation: "UJI ALAMAT 1", Occupations: []models.OkupasiObjek{{OccupationID: "UJI-OKP-A", OccupationName: "UJI A"}},
			Items: []models.ItemObjek{{Currency: "USD", TSI: ds("1000"), Coverages: []models.CoverageObjek{{Zone: "Daerah Zona 2"}}}}},
		{ObjectNo: "2", IsTopRisk: true, RiskLocation: "UJI ALAMAT 2", Occupations: []models.OkupasiObjek{{OccupationID: "UJI-OKP-B", OccupationName: "UJI B"}},
			Items: []models.ItemObjek{{Currency: "IDR", TSI: ds("31000000000"), Coverages: []models.CoverageObjek{{Zone: "Daerah Zona 1"}, {}}}}},
		{ObjectNo: "3", RiskLocation: "UJI ALAMAT 3"},
	}}}
	sp := &spreadingTiruan{baca: []models.KasusSpreading{{StartDateTime: "20250101T170000.000 GMT"}},
		idUang: map[string]string{"IDR": idMataUangIDR, "USD": "UJI-USD"}, kurs: map[string]*apd.Decimal{"UJI-USD": ds("16000")}}
	return Baru(nil).DenganScoring(st).DenganKasus(kasus).DenganObjek(objek).DenganSpreading(sp).
		DenganOccupation(okupasiScoringTiruan{}).DenganTransaksi(tanpaTx)
}

// TestTampilanScoringKasus - GET: kolom 1 = lokasi IsTopRisk (bukan lokasi pertama), isi otomatis dari lokasi itu, kurs
// USD Begin date, Status / tanggal / daftar dominan; bolehUbah = posisi "Offer" + flag "0".
func TestTampilanScoringKasus(t *testing.T) {
	st := &scoringTiruan{k: models.KasusScoring{Position: repository.PosisiAwalCaseNB, FlagOnGoingPolicy: "0"}}
	tp, err := layananScoring(st).TampilanScoringKasus(context.Background(), "UJI-NB-1")
	k := tp.Scoring.Kolom[0]
	if err != nil || !tp.BolehUbah || tp.Status != "New Business" || tp.InsuredName != "UJI TERTANGGUNG" || tp.Mulai != "02-01-2025" ||
		k.OccupationNote != "UJI B" || k.TopRiskLocation != "UJI ALAMAT 2" || k.Pilihan["occupation"] != 4 ||
		k.Pilihan["totalSumInsured"] != 2 || k.Pilihan["earthquake"] != 2 || tp.Scoring.Kolom[1].Ada ||
		!reflect.DeepEqual(tp.LokasiDominan, []LokasiDominan{{"1", "UJI A"}, {"2", "UJI B"}, {"3", ""}}) {
		t.Fatalf("%v %+v", err, tp)
	}
	for _, kasus := range []models.KasusScoring{{Position: "Offer", FlagOnGoingPolicy: "1"}, {Position: "Akseptasi", FlagOnGoingPolicy: "0"},
		{Position: "ReasFacInMarketing", FlagOnGoingPolicy: "0"}} {
		st.k = kasus
		tp, _ := layananScoring(st).TampilanScoringKasus(context.Background(), "UJI-NB-1")
		if tp.BolehUbah != (kasus.Position == "ReasFacInMarketing") {
			t.Errorf("bolehUbah %+v: %v", kasus, tp.BolehUbah)
		}
	}
}

// TestHitungDanSimpanScoring - kolom 2 = lokasi dominan pilihan (USD x 16.000 -> TSI 16.000.000 ≤ ambang), 409 baca-saja,
// 400 isian / lokasi tak dikenal / okupasi sama (PUT saja), 401 tanpa identitas, simpan tepat sekali.
func TestHitungDanSimpanScoring(t *testing.T) {
	ctx := context.Background()
	st := &scoringTiruan{k: models.KasusScoring{Position: "Offer", FlagOnGoingPolicy: "0"}}
	svc := layananScoring(st)
	m := [2]models.KolomScoring{{Pilihan: map[string]int{"riot": 1}}, {IndexLocation: "1", Pilihan: map[string]int{"flood": 2}}}
	tp, err := svc.HitungScoring(ctx, "UJI-NB-1", m)
	d := tp.Scoring.Kolom[1]
	if err != nil || !d.Ada || d.OccupationNote != "UJI A" || d.Pilihan["totalSumInsured"] != 1 || d.Pilihan["earthquake"] != 3 ||
		d.Pilihan["flood"] != 2 || d.Total.String() != "21.25" || len(st.tulis) != 0 { // 15 - 1 + 2.25 (Zona 2) + 5
		t.Fatalf("hitung: %v %+v", err, d)
	}
	akun := inti.Pelaku{AkunID: "UJI-USER"}
	if _, err := svc.SimpanScoring(ctx, akun, "UJI-NB-1", m); err != nil || len(st.tulis) != 1 || !st.tulis[0].Kolom[1].Ada {
		t.Fatalf("simpan: %v %d", err, len(st.tulis))
	}
	if _, err := svc.SimpanScoring(ctx, inti.Pelaku{}, "UJI-NB-1", m); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("401: %v", err)
	}
	for nama, salah := range map[string][2]models.KolomScoring{
		"kunci":   {{Pilihan: map[string]int{"uji": 1}}},
		"nomor":   {{Pilihan: map[string]int{"fireAlarm": 3}}},
		"negatif": {{Pilihan: map[string]int{"riot": -1}}},
		"lokasi":  {{}, {IndexLocation: "UJI-TAK-ADA"}},
		"panjang": {{}, {IndexLocation: strings.Repeat("9", 51)}},
	} {
		if _, err := svc.HitungScoring(ctx, "UJI-NB-1", salah); !errors.Is(err, ErrMasukanScoring) {
			t.Errorf("%s: %v", nama, err)
		}
	}
	sama := [2]models.KolomScoring{{}, {IndexLocation: "2"}}
	if tp, err := svc.HitungScoring(ctx, "UJI-NB-1", sama); err != nil || !reflect.DeepEqual(tp.Pesan, []string{PesanOkupasiSama}) {
		t.Errorf("hitung okupasi sama: %v %v", err, tp.Pesan)
	}
	if _, err := svc.SimpanScoring(ctx, akun, "UJI-NB-1", sama); !errors.Is(err, ErrMasukanScoring) || len(st.tulis) != 1 {
		t.Errorf("simpan okupasi sama: %v %d", err, len(st.tulis))
	}
	st.k.FlagOnGoingPolicy = "2"
	if _, err := svc.HitungScoring(ctx, "UJI-NB-1", m); !errors.Is(err, ErrScoringBacaSaja) {
		t.Errorf("409 hitung: %v", err)
	}
	if _, err := svc.SimpanScoring(ctx, akun, "UJI-NB-1", m); !errors.Is(err, ErrScoringBacaSaja) || len(st.tulis) != 1 {
		t.Errorf("409 simpan: %v", err)
	}
	st.galatB = repository.ErrKasusTidakAda
	if _, err := svc.TampilanScoringKasus(ctx, "UJI-NB-1"); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("404: %v", err)
	}
	if _, err := Baru(nil).TampilanScoringKasus(ctx, "UJI-NB-1"); !errors.Is(err, ErrScoringTanpaDatabase) {
		t.Errorf("503: %v", err)
	}
}

// TestScoringKursTerpisah - kurs mata uang yang hanya ada di premi spread (UJI-EUR) tidak mengosongkan TSI (kolom Top tetap
// TSI 2 tanpa pesan); share-nya kosong. Dominant tanpa Occupation -> kolom 2 tidak dihitung + pesan.
func TestScoringKursTerpisah(t *testing.T) {
	st := &scoringTiruan{k: models.KasusScoring{Position: "Offer", FlagOnGoingPolicy: "0"}}
	svc := layananScoring(st)
	svc.DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{{StartDateTime: "20250101T170000.000 GMT",
		Lokasi: []models.LokasiSpreading{{ObjectNo: "2", Items: []models.ItemSpreading{{Currency: "UJI-EUR",
			Coverages: []models.CoverageSpreading{{Spreading: []models.BarisSpreading{{TreatyType: "UJI-QS", PremiumSpreaded: ds("5")}}}}}}}}}},
		idUang: map[string]string{"IDR": idMataUangIDR}})
	l, err := svc.lokasiKolom(context.Background(), konteksUji(t, svc), true, "")
	if err != nil || l.TSIIDR == nil || l.TSIIDR.String() != "31000000000" || l.ShareIDR != nil || len(l.KursTidakAda) != 0 {
		t.Fatalf("%v %+v", err, l)
	}
	tp, err := svc.HitungScoring(context.Background(), "UJI-NB-1", [2]models.KolomScoring{{}, {IndexLocation: "3"}})
	if err != nil || tp.Scoring.Kolom[0].Pilihan[models.KunciTotalSumInsured] != 2 || tp.Scoring.Kolom[1].Ada ||
		!reflect.DeepEqual(tp.Pesan, []string{"Dominant Risk lokasi 3 belum punya Occupation - kolom Dominant Risk tidak dihitung"}) {
		t.Errorf("%v %+v", err, tp.Pesan)
	}
}

func konteksUji(t *testing.T, svc *Service) konteksScoring {
	t.Helper()
	k, err := svc.bacaKonteksScoring(context.Background(), "UJI-NB-1")
	if err != nil {
		t.Fatal(err)
	}
	return k
}
