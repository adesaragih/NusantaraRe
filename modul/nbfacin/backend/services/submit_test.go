package services

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// submitTiruan - PenyimpanSubmit dalam memori. RW: kode pos di peta = ada (berikut pasangan kotanya); kategori okupasi
// per kode.
type submitTiruan struct {
	konteks    map[string]models.KonteksSubmit
	akun       models.AkunSubmit
	riwayat    []models.RiwayatAkseptasi
	banding    []string
	sebab      []models.BarisCauseOfDecline
	rw         map[string]models.HasilCekRW
	kategori   map[string]string
	tulis      []models.TulisanSubmit
	galatTulis error
	idBanding  string
}

func (s *submitTiruan) KonteksSubmit(_ context.Context, id string) (models.KonteksSubmit, error) {
	k, ok := s.konteks[id]
	if !ok {
		return k, repository.ErrKasusTidakAda
	}
	return k, nil
}
func (s *submitTiruan) AkunSubmit(context.Context, string) (models.AkunSubmit, error) {
	return s.akun, nil
}
func (s *submitTiruan) RiwayatAkseptasi(context.Context, string) ([]models.RiwayatAkseptasi, error) {
	return s.riwayat, nil
}
func (s *submitTiruan) OpsiBanding(_ context.Context, idPega, _ string) ([]string, error) {
	s.idBanding = idPega
	return s.banding, nil
}
func (s *submitTiruan) SebabDecline(context.Context) ([]models.BarisCauseOfDecline, error) {
	return s.sebab, nil
}
func (s *submitTiruan) CekRW(_ context.Context, zip, _, _, _ string) (models.HasilCekRW, error) {
	return s.rw[zip], nil
}
func (s *submitTiruan) KategoriOkupasi(_ context.Context, kode []string) (map[string]string, error) {
	h := map[string]string{}
	for _, k := range kode {
		if c, ada := s.kategori[k]; ada {
			h[k] = c
		}
	}
	return h, nil
}
func (s *submitTiruan) TulisSubmit(_ context.Context, _ *db.Tx, t models.TulisanSubmit) error {
	if s.galatTulis != nil {
		return s.galatTulis
	}
	s.tulis = append(s.tulis, t)
	return nil
}

// lokasiSah - satu lokasi FIRE yang lolos seluruh proteksi lokasi.
func lokasiSah() models.ObjekFire {
	it := itemSah()
	it.Coverages = []models.CoverageObjek{{Coverage: "FLEXAS", CoverageNote: "UJI", OldID: "FLEXAS", AccumulationCode: "UJI-AKU"}}
	return models.ObjekFire{Province: "UJI PROV", City: "UJI KOTA", District: "UJI KEC", ZipCode: "12345",
		Occupations: []models.OkupasiObjek{{OccupationID: "2001"}}, LossRecords: []models.CatatanKerugian{{Currency: "IDR"}},
		Items: []models.ItemObjek{it}}
}

var pelakuMkt = inti.Pelaku{AkunID: "UJI-MKT", Peran: []string{"ReasFacInMarketing"}}

// layananSubmit - kasus UJI-NB-1 tahap marketing, Scoring tersimpan, MOID terisi, satu lokasi sah, payment 100 % satu
// cicilan per mata uang (IDR / USD, layananPayment).
func layananSubmit(st *submitTiruan, lokasi ...models.ObjekFire) *Service {
	if lokasi == nil {
		lokasi = []models.ObjekFire{lokasiSah()}
	}
	kasus := &models.Kasus{CaseID: "UJI-NB-1", Opportunity: models.Opportunity{InsuredID: "UJI-ORG"},
		General: models.General{MarketingID: "UJI-MO"}, RiskScoring: models.RingkasanScoring{FinalScore: d("61")}}
	cicilan := func() []models.CicilanPayment {
		return []models.CicilanPayment{{InstallmentNo: 1, DueDate: "20261007", Persen: d("100")}}
	}
	p := &paymentTiruan{ada: models.KasusPayment{Installment: "1", Cicilan: map[string][]models.CicilanPayment{"IDR": cicilan(), "USD": cicilan()}}}
	n := 0
	jam := func() time.Time { return time.Date(2026, 10, 7, 15, 4, 5, 0, time.UTC) } // 22:04:05 WIB
	return layananPayment(p).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": kasus}}).
		DenganObjek(objekTiruan{ada: map[string][]models.ObjekFire{"UJI-NB-1": lokasi}}).
		DenganPilihanItem(pilihanItemTiruan{&n}).DenganSubmit(st).DenganJam(jam)
}

func submitUji() *submitTiruan {
	return &submitTiruan{
		konteks: map[string]models.KonteksSubmit{"UJI-NB-1": {Position: "Offer", FlagOnGoingPolicy: "0", MOID: "UJI-MO",
			MarketingName: "Uji - Pemasar", TeamGroup: "1", MOLeader: "UJI-LEADER", InsuredName: "UJI TERTANGGUNG",
			SobName: "UJI SOB", ClassOfBusiness: "UJI COB", PercentShare: d("12.5"), RICommision: d("10")}},
		akun:     models.AkunSubmit{Nama: "Uji Pengguna", Divisi: "UJI-DIV"},
		banding:  []string{"SENIORUW", "SENIORUW", "JUW_B", "MANAGERTEKNIK"},
		sebab:    barisCauseOfDeclineCSV,
		rw:       map[string]models.HasilCekRW{"12345": {Zip: true, City: true, District: true, Province: true}},
		kategori: map[string]string{"2001": "Preferred Risk"},
	}
}

// TestProteksiAcceptLolos - kasus sah: tanpa pesan (instrumen: setiap pesan di TestProteksiAcceptGagal muncul hanya bila
// datanya dirusak).
func TestProteksiAcceptLolos(t *testing.T) {
	p, err := layananSubmit(submitUji()).ProteksiAccept(context.Background(), "UJI-NB-1")
	if err != nil || len(p) != 0 {
		t.Fatalf("%v %q", err, p)
	}
}

// TestProteksiAcceptGagal - pesan verbatim per langkah, urutan langkah Pega, kembar sekali.
func TestProteksiAcceptGagal(t *testing.T) {
	ctx := context.Background()
	rusak := lokasiSah()
	rusak.Province, rusak.City, rusak.ZipCode, rusak.District = "", "", "99999", "UJI KEC"
	rusak.Occupations = []models.OkupasiObjek{{OccupationID: ""}}
	rusak.LossRecords = nil
	it := rusak.Items[0]
	it.Currency, it.ItemTypeID = "", ""
	it.TotalNetRate = d("0.5")
	it.Coverages = []models.CoverageObjek{{Coverage: "FLEXAS", OldID: "FLEXAS", AccumulationCode: "UNKNOWNID"}}
	asing := itemSah()
	asing.Currency = "UJI-XYZ" // tidak di CURRENCY
	rusak.Items = []models.ItemObjek{it, asing}
	st := submitUji()
	svc := layananSubmit(st, lokasiSah(), rusak)
	k := svc.kasus.(kasusTiruan).ada["UJI-NB-1"]
	k.General.MarketingID, k.RiskScoring.FinalScore = "", nil
	p, err := svc.ProteksiAccept(ctx, "UJI-NB-1")
	mau := []string{
		"Harus memilih dari daftar nama marketing yang tersedia",
		"Scoring Risk Cannot Empty !!",
		"Currency Pada Lokasi 2 Object Item 2 Bermasalah, Coba Pilih Ulang Currency",
		"Province on Location 2can't be null!",
		"City on Location 2can't be null!",
		"Loss record on location 2 can't be empty !",
		"ZipCode on Location 2 doesn’t exist in Master!",
		"CITYNAME  on Location 2 doesn’t exist in Master!",
		"DISTRICT UJI KEC on Location 2 doesn’t exist in Master!",
		"PROVINCE  on Location 2 doesn’t exist in Master!",
		"Occupation ID on location2can't be null!",
		"Currency on location 2 object Item 1 must be choose!",
		"Choose Object Item on location 2 !",
		" Coverage 1 on location 2 object item 1 can't be empty!",
		" Accumulation on location 2 object item 1 can't be empty!",
		"coverage on location 2 object item 1 can't be Net Rate!",
	}
	if err != nil || !reflect.DeepEqual(p, mau) {
		t.Fatalf("%v\n%q\nmau\n%q", err, p, mau)
	}
	// Langkah 2: > 400 lokasi / InsuredID dikecualikan -> tanpa proteksi.
	banyak := make([]models.ObjekFire, 401)
	if p, err := layananSubmit(submitUji(), banyak...).ProteksiAccept(ctx, "UJI-NB-1"); err != nil || len(p) != 0 {
		t.Errorf("> 400 lokasi: %v %q", err, p)
	}
	svc = layananSubmit(submitUji(), rusak)
	svc.kasus.(kasusTiruan).ada["UJI-NB-1"].Opportunity.InsuredID = "ASM-SFAGIS-WORK-ORG ORG-7889"
	if p, err := svc.ProteksiAccept(ctx, "UJI-NB-1"); err != nil || len(p) != 0 {
		t.Errorf("ORG-7889: %v %q", err, p)
	}
}

// TestProteksiKategoriOkupasi - IsPreferredRisk kosong -> pesan kategori terakhir; Preferred di salah satu okupasi -> lolos.
func TestProteksiKategoriOkupasi(t *testing.T) {
	ctx := context.Background()
	o := lokasiSah()
	o.Occupations = []models.OkupasiObjek{{OccupationID: "UJI-9"}, {OccupationID: "UJI-8"}}
	p, err := layananSubmit(submitUji(), o).ProteksiAccept(ctx, "UJI-NB-1")
	if err != nil || !reflect.DeepEqual(p, []string{"Kategori Okupasi ID UJI-8 lokasi1 tidak boleh kosong"}) {
		t.Errorf("%v %q", err, p)
	}
	o.Occupations = append(o.Occupations, models.OkupasiObjek{OccupationID: "2001"})
	if p, err := layananSubmit(submitUji(), o).ProteksiAccept(ctx, "UJI-NB-1"); err != nil || len(p) != 0 {
		t.Errorf("preferred: %v %q", err, p)
	}
}

// TestProteksiAcceptPayment - tanpa cicilan: "List Installment can't be null!" sekali (kembar per mata uang) DAN "Total premium
// RNM ..." (Pega memulai jumlah cicilan dari 0 - premi 0 <> premi mata uang; ProtectPremiPolicy_Act 5.1 / 5.5).
func TestProteksiAcceptPayment(t *testing.T) {
	svc := layananSubmit(submitUji())
	svc.payment.(*paymentTiruan).ada.Cicilan = map[string][]models.CicilanPayment{}
	p, err := svc.ProteksiAccept(context.Background(), "UJI-NB-1")
	if err != nil || !reflect.DeepEqual(p, []string{"List Installment can't be null!", "Total premium RNM on installment didn't match with premium RNM!"}) {
		t.Errorf("%v %q", err, p)
	}
}

// TestTampilanSubmit - boleh / alasan (posisi, flag, workbasket), riwayat, opsi Banding unik ber-pzInsKey, sebab Decline
// ber-opsi kosong, RI Comm.
func TestTampilanSubmit(t *testing.T) {
	ctx := context.Background()
	st := submitUji()
	st.riwayat = []models.RiwayatAkseptasi{{Tanggal: "07-10-2026 22:04", PIC: "UJI", Approval: "Accept", Status: "Offer", Comment: "ok"}}
	svc := layananSubmit(st)
	v, err := svc.TampilanSubmit(ctx, pelakuMkt, "UJI-NB-1")
	if err != nil || !v.BolehSubmit || v.Alasan != "" || len(v.Riwayat) != 1 || st.idBanding != "ASM-FW-GISFW-WORK UJI-NB-1" ||
		!reflect.DeepEqual(v.OpsiBanding, []string{"SENIORUW", "JUW_B"}) || len(v.SebabTc) != 14 ||
		v.SebabTc[0].Opsi == nil || len(v.SebabTc[0].Opsi) != 0 || len(v.SebabNonTc) != 7 || v.PctRnmRIComm != "10" || v.PctRnmRate != "" {
		t.Fatalf("%v %+v", err, v)
	}
	if v, _ := svc.TampilanSubmit(ctx, inti.Pelaku{AkunID: "UJI-UW", Peran: []string{"ReasFacInUnderwriting"}}, "UJI-NB-1"); v.BolehSubmit ||
		!strings.Contains(v.Alasan, "workbasket ReasFacInMarketing") {
		t.Errorf("bukan marketing: %+v", v)
	}
	k := st.konteks["UJI-NB-1"]
	k.FlagOnGoingPolicy = "1"
	st.konteks["UJI-NB-1"] = k
	if v, _ := svc.TampilanSubmit(ctx, pelakuMkt, "UJI-NB-1"); v.BolehSubmit || !strings.Contains(v.Alasan, "tahap marketing") {
		t.Errorf("flag 1: %+v", v)
	}
	if _, err := svc.TampilanSubmit(ctx, inti.Pelaku{}, "UJI-NB-1"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("401: %v", err)
	}
	if _, err := svc.TampilanSubmit(ctx, pelakuMkt, "UJI-NB-X"); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("404: %v", err)
	}
}

// TestAjukanSubmit - empat Email Type: posisi, NB_STATUS, baris riwayat (jam 24 WIB), FACINPERFORMANCE hanya Decline.
func TestAjukanSubmit(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		isi            models.IsiSubmit
		posisi, status string
	}{
		{models.IsiSubmit{EmailType: "1", Comment: " ok "}, "ReasFacInTeamLeader", "NB IS IN TEAM LEADER'S INBOX"},
		{models.IsiSubmit{EmailType: "2", Comment: "tolak"}, "Offer", "NB IS IN UJI PEMASAR'S OPPORTUNITY"},
		{models.IsiSubmit{EmailType: "4", Comment: "banding", BandingTo: "SENIORUW"}, "ReasFacInSeniorUnderwriting", "NB IS IN SENIORUW'S INBOX"},
		{models.IsiSubmit{EmailType: "4", Comment: "banding", BandingTo: "JUW_B"}, "ReasFacInJuniorUnderwriting", "NB IS IN JUNIOR UNDERWRITER (B)'S INBOX"},
		{models.IsiSubmit{EmailType: "7", Comment: "diabaikan", SebabTc: []models.PilihanSebabTc{{TermCondition: "Rate", Prioritas: "I", PctRNM: "1.5", PctSOB: "2"},
			{TermCondition: "Occupancy"}, {TermCondition: "Risk Exposure", Prioritas: "II", Description1: "Others", Description2: "UJI"}}}, "Decline",
			"NB WAS DECLINED BY  UJI PENGGUNA"},
	} {
		st := submitUji()
		h, err := layananSubmit(st).AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", c.isi)
		if err != nil || h.Position != c.posisi || h.NBStatus != c.status || len(h.Pesan) != 0 || len(st.tulis) != 1 {
			t.Fatalf("%s: %v %+v", c.isi.EmailType, err, h)
		}
		w := st.tulis[0]
		label := map[string]string{"1": "Accept", "2": "Reject", "4": "Banding", "7": "Decline"}[c.isi.EmailType]
		mauRiwayat := models.RiwayatProduksi{TypePolis: "NB", Posisi: "Offer", PIC: "Uji Pengguna", TglInp: "07/10/2026 22:04:05",
			Div: "UJI-DIV", Type: "F", Putaran: "1", Approval: label, Keterangan: w.Riwayat.Keterangan, AksesLogin: "UJI-MKT",
			PercentRNM: w.Riwayat.PercentRNM}
		if w.IDPega != "ASM-FW-GISFW-WORK UJI-NB-1" || w.PosisiLama != "Offer" || w.Riwayat != mauRiwayat ||
			w.Riwayat.PercentRNM.String() != "12.5" || w.Pega != (models.RiwayatPega{Status: strings.ToUpper(label), Username: "Uji Pengguna",
			Workbasket: "ReasFacInMarketing"}) || (w.Performance != nil) != (c.isi.EmailType == "7") {
			t.Errorf("%s tulisan: %+v", c.isi.EmailType, w)
		}
		switch c.isi.EmailType {
		case "1":
			if w.Riwayat.Keterangan != "ok" {
				t.Errorf("comment dipangkas: %q", w.Riwayat.Keterangan)
			}
		case "7":
			if w.StatusWorkBaru != "Resolved-Decline" || w.Riwayat.Keterangan != "Rate   1.5 % RNM 2 % SOB \nRisk Exposure Others UJI " || *w.Performance != (models.PerformanceFacIn{PositionDoc: "OFFER",
				StatusDoc: "DECLINE", Username: "UJI-MKT", InsuredName: "UJI TERTANGGUNG", SOB: "UJI SOB", MOLeader: "UJI-LEADER", COB: "UJI COB"}) {
				t.Errorf("decline: %q %+v", w.Riwayat.Keterangan, w.Performance)
			}
		}
	}
	// Decline: teks baris T/C TIDAK dipangkas (Comment = komentarDecline popup).
	st0 := submitUji()
	if _, err := layananSubmit(st0).AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "7",
		SebabTc: []models.PilihanSebabTc{{TermCondition: "Occupancy", Prioritas: "I", Description2: " gudang "}}}); err != nil ||
		st0.tulis[0].Riwayat.Keterangan != "Occupancy   gudang  \n" {
		t.Errorf("tanpa pangkas: %v %q", err, st0.tulis[0].Riwayat.Keterangan)
	}
	// FACINPERFORMANCE dilewati bila kolom wajib kosong.
	st := submitUji()
	k := st.konteks["UJI-NB-1"]
	k.MOLeader = ""
	st.konteks["UJI-NB-1"] = k
	if _, err := layananSubmit(st).AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "7", SebabNonTc: "Accumulation / Capacity"}); err != nil ||
		st.tulis[0].Performance != nil || st.tulis[0].Riwayat.Keterangan != "Accumulation / Capacity\n" {
		t.Errorf("performance dilewati: %v %+v", err, st.tulis)
	}
}

// TestAjukanSubmitDitolak - 401 / 404 / 403 / 409 / 400 / 422 tanpa satu pun tulisan.
func TestAjukanSubmitDitolak(t *testing.T) {
	ctx := context.Background()
	sah := models.IsiSubmit{EmailType: "1", Comment: "ok"}
	for nama, c := range map[string]struct {
		pelaku inti.Pelaku
		id     string
		isi    models.IsiSubmit
		ubah   func(*submitTiruan, *Service)
		mau    error
	}{
		"401":            {inti.Pelaku{}, "UJI-NB-1", sah, nil, inti.ErrTanpaIdentitas},
		"404":            {pelakuMkt, "UJI-NB-X", sah, nil, ErrKasusTidakAda},
		"403 workbasket": {inti.Pelaku{AkunID: "UJI-UW", Peran: []string{"ReasFacInTeamLeader"}}, "UJI-NB-1", sah, nil, inti.ErrTanpaWewenang},
		"409 posisi": {pelakuMkt, "UJI-NB-1", sah, func(st *submitTiruan, _ *Service) {
			k := st.konteks["UJI-NB-1"]
			k.Position = "ReasFacInTeamLeader"
			st.konteks["UJI-NB-1"] = k
		}, ErrSubmitBukanTahapMarketing},
		"409 serentak":   {pelakuMkt, "UJI-NB-1", sah, func(st *submitTiruan, _ *Service) { st.galatTulis = repository.ErrPosisiBerubah }, ErrSubmitBukanTahapMarketing},
		"400 ask":        {pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "3", Comment: "x"}, nil, ErrMasukanSubmit},
		"400 comment":    {pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "2", Comment: "  "}, nil, ErrMasukanSubmit},
		"400 banding":    {pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "4", Comment: "x", BandingTo: "DIREKTURTEKNIK"}, nil, ErrMasukanSubmit},
		"400 managertek": {pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "4", Comment: "x", BandingTo: "MANAGERTEKNIK"}, nil, ErrMasukanSubmit},
		"400 decline":    {pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "7", SebabTc: []models.PilihanSebabTc{{TermCondition: "Rate", Prioritas: "I"}}}, nil, ErrMasukanSubmit},
		"403 sebelum 409": {inti.Pelaku{AkunID: "UJI-UW", Peran: []string{"ReasFacInTeamLeader"}}, "UJI-NB-1", sah, func(st *submitTiruan, _ *Service) {
			k := st.konteks["UJI-NB-1"]
			k.Position = "Decline"
			st.konteks["UJI-NB-1"] = k
		}, inti.ErrTanpaWewenang},
		"400 option asing": {pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "7", SebabTc: []models.PilihanSebabTc{{TermCondition: "Period Of Insurance (POI)", Prioritas: "I", Description1: "Power plant"}}}, nil, ErrMasukanSubmit},
		"400 tc asing":     {pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "7", SebabTc: []models.PilihanSebabTc{{TermCondition: "UJI asing", Prioritas: "I"}}}, nil, ErrMasukanSubmit},
		"400 sebab asing":  {pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "7", SebabNonTc: "UJI asing"}, nil, ErrMasukanSubmit},
		"400 dua I":        {pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "7", SebabTc: []models.PilihanSebabTc{{TermCondition: "Rate", Prioritas: "I", PctRNM: "1", PctSOB: "1"}, {TermCondition: "Occupancy", Prioritas: "I", Description2: "x"}}}, nil, ErrMasukanSubmit},
		"422 proteksi": {pelakuMkt, "UJI-NB-1", sah, func(_ *submitTiruan, svc *Service) {
			svc.kasus.(kasusTiruan).ada["UJI-NB-1"].RiskScoring.FinalScore = nil
		}, ErrProteksiSubmit},
	} {
		st := submitUji()
		svc := layananSubmit(st)
		if c.ubah != nil {
			c.ubah(st, svc)
		}
		_, err := svc.AjukanSubmit(ctx, c.pelaku, c.id, c.isi)
		if !errors.Is(err, c.mau) || (c.mau != ErrSubmitBukanTahapMarketing || nama != "409 serentak") && len(st.tulis) != 0 {
			t.Errorf("%s: %v (tulis %d)", nama, err, len(st.tulis))
		}
		var g GalatProteksi
		if nama == "422 proteksi" && (!errors.As(err, &g) || !reflect.DeepEqual(g.Pesan, []string{"Scoring Risk Cannot Empty !!"})) {
			t.Errorf("422 pesan: %v", err)
		}
	}
}

// TestPeriksaSubmit - Accept menjalankan proteksi; Reject / Banding / Decline kosong; emailType asing 400.
func TestPeriksaSubmit(t *testing.T) {
	ctx := context.Background()
	svc := layananSubmit(submitUji())
	svc.kasus.(kasusTiruan).ada["UJI-NB-1"].RiskScoring.FinalScore = nil
	for e, mau := range map[string]int{"1": 1, "2": 0, "4": 0, "7": 0} {
		h, err := svc.PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", e)
		if err != nil || len(h.Pesan) != mau || h.Pesan == nil {
			t.Errorf("%s: %v %+v", e, err, h)
		}
	}
	if _, err := svc.PeriksaSubmit(ctx, pelakuMkt, "UJI-NB-1", "9"); !errors.Is(err, ErrMasukanSubmit) {
		t.Errorf("400: %v", err)
	}
}

// TestSusunSebabDecline - 41 baris CAUSEOFDECLINE.csv: 14 T/C urut ID ("Clauses" dipangkas), opsi POI / Extension of
// coverage / Risk Exposure (teks TERMCONDITION, ID 41 "Others" sesudah ID 33), 7 Non - T/C; jawaban diketahui dari data.
func TestSusunSebabDecline(t *testing.T) {
	tc, nonTc := susunSebabDecline(barisCauseOfDeclineCSV)
	var nama []string
	opsi := map[string][]string{}
	for _, x := range tc {
		nama = append(nama, x.TermCondition)
		opsi[x.TermCondition] = x.Opsi
	}
	if !reflect.DeepEqual(nama, []string{"Rate", "R/I Comm, Deduction", "OJK Discount", "Loss Ratio (net)", "Occupancy",
		"Period Of Insurance (POI)", "TSI (in IDR)", "Limit Of Liability (LOL)", "Deductible", "Clauses", "Extension of coverage",
		"Risk Exposure", "Lack of data for acceptance", "Others T/C"}) {
		t.Errorf("T/C %q", nama)
	}
	for term, mau := range map[string][]string{
		"Period Of Insurance (POI)": {"Backdated", "Short Period", "> 12 months (Non CAR/EAR)", "> 60 months (CAR/EAR)"},
		"Extension of coverage":     {"Terrorism and Sabotage", "Sub Limit MB", "Overseas Jurisdiction", "Others"},
		"Risk Exposure": {"Flood prone", "Earthquake prone", "Building under construction", "Wet Risk", "Minning risk", "Offshore risk",
			"Oil and gas", "Power plant", "Progress of Work > 30%  (CAR/EAR)", "Conveyange (MC)", "Interest Insured (MC)", "Others"},
		"Rate": {},
	} {
		if !reflect.DeepEqual(opsi[term], mau) {
			t.Errorf("opsi %s %q", term, opsi[term])
		}
	}
	if !reflect.DeepEqual(nonTc, []string{"Accumulation / Capacity", "SOB lost business", "No response from SOB", "Cancel - double input",
		"Cancel - incorect main data", "Retain by cedant", "Late confirmation"}) {
		t.Errorf("Non - T/C %q", nonTc)
	}
	// Baris ber-TERMCONDITION2 bukan opsi walau DESCRIPTION-nya nama T/C (tidak terjadi di data; syarat dijaga).
	tc, _ = susunSebabDecline([]models.BarisCauseOfDecline{{ID: "1", TermCondition: "Rate"},
		{ID: "2", TermCondition: "UJI", Description: "Rate", TermCondition2: "UJI"}})
	if len(tc) != 1 || len(tc[0].Opsi) != 0 {
		t.Errorf("opsi ber-TERMCONDITION2: %+v", tc)
	}
}

// TestPotongKomentar - 3990 karakter, lalu <= 4000 bita.
func TestPotongKomentar(t *testing.T) {
	if s := potongKomentar(strings.Repeat("a", 4100)); len(s) != 3990 {
		t.Errorf("ascii %d", len(s))
	}
	if s := potongKomentar(strings.Repeat("é", 3000)); len(s) != 4000 || !strings.HasPrefix(s, "é") {
		t.Errorf("multibita %d", len(s))
	}
}
