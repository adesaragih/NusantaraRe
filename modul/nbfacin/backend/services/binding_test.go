package services

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// TestKasusGroup - When IsGroup (tiket 63 D5): A = pembuat case (CREATE_OP) di daftar login, tidak peka huruf; B-D = akun yang
// memproses di daftar login; E = MARKETINGOFFICER.CLIENTID persis; konfigurasi kosong / login kosong tidak pernah group. Nilai
// sintetis.
func TestKasusGroup(t *testing.T) {
	g := KonfigGrup{Login: []string{"UJI-GRUP"}, MarketingCode: []string{"UJI-MC-1"}}
	k := models.KonteksSubmit{CreateOp: "uji-grup"}
	lain := inti.Pelaku{AkunID: "UJI-LAIN"}
	for nama, c := range map[string]struct {
		g      KonfigGrup
		k      models.KonteksSubmit
		pelaku inti.Pelaku
		mau    bool
	}{
		"A pembuat":          {g, k, lain, true},
		"B akun memproses":   {g, models.KonteksSubmit{CreateOp: "UJI-X"}, inti.Pelaku{AkunID: " uji-grup "}, true},
		"E kode marketing":   {g, models.KonteksSubmit{MarketingCode: "UJI-MC-1"}, lain, true},
		"E beda huruf":       {g, models.KonteksSubmit{MarketingCode: "uji-mc-1"}, lain, false},
		"bukan group":        {g, models.KonteksSubmit{CreateOp: "UJI-X", MarketingCode: "UJI-MC-2"}, lain, false},
		"konfigurasi kosong": {KonfigGrup{}, k, inti.Pelaku{AkunID: "uji-grup"}, false},
		"login kosong":       {KonfigGrup{Login: []string{"UJI-GRUP"}, MarketingCode: []string{""}}, models.KonteksSubmit{}, inti.Pelaku{}, false},
	} {
		if got := c.g.kasusGroup(c.k, c.pelaku); got != c.mau {
			t.Errorf("%s: %v", nama, got)
		}
	}
	if !g.pembuatGroup(k) || g.pembuatGroup(models.KonteksSubmit{CreateOp: "UJI-X"}) {
		t.Error("pembuatGroup")
	}
}

// TestAcceptGroupDitahan - case group: ActFlow langkah 3 "exit kalau group" -> Decision23 Else -> Decision16 IsGroup ->
// SubProcess POLICY (tiket 64) DITAHAN pesanPolicyGroup, tanpa tulisan; tanpa konfigurasi group Senior TSI 4000 -> DepHead.
func TestAcceptGroupDitahan(t *testing.T) {
	ctx := context.Background()
	for nama, c := range map[string]struct {
		g       KonfigGrup
		ubah    func(*submitTiruan)
		ditahan bool
	}{
		"A pembuat case": {KonfigGrup{Login: []string{"uji-pembuat"}}, nil, true},
		"B akun":         {KonfigGrup{Login: []string{"UJI-UW"}}, nil, true},
		"E marketing": {KonfigGrup{MarketingCode: []string{"UJI-MC-1"}}, func(st *submitTiruan) {
			k := st.konteks["UJI-NB-1"]
			k.MarketingCode = "UJI-MC-1"
			st.konteks["UJI-NB-1"] = k
		}, true},
		"bukan group": {KonfigGrup{}, nil, false},
	} {
		st := submitUji()
		if c.ubah != nil {
			c.ubah(st)
		}
		svc := layananSubmitUW(st, posisiSeniorUW, spreadingTanpaFacRetro("4000")).DenganGrup(c.g, false)
		svc.limit = limitTiruan{a: limitLanjut()}
		h, err := svc.AjukanSubmit(ctx, pelakuPosisi(posisiSeniorUW), "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "x"})
		var g GalatProteksi
		if c.ditahan {
			if !errors.As(err, &g) || g.Pesan[len(g.Pesan)-1] != pesanPolicyGroup || len(st.tulis) != 0 {
				t.Errorf("%s: %v", nama, err)
			}
			continue
		}
		if err != nil || h.Position != posisiDepHeadUW || len(st.tulis) != 1 || st.tulis[0].Binding != nil {
			t.Errorf("%s: %v %+v", nama, err, h)
		}
	}
}

// objekFacInOffer / spreadingFacInOffer - satu lokasi, dua okupasi, satu item, dua coverage (yang kedua tanpa spreading),
// coverage pertama dua baris spreading.
func objekFacInOffer() []models.ObjekFire {
	return []models.ObjekFire{{ObjectNo: "1", ObjectName: "UJI OBJEK", IsTopRisk: true, RoadType: "JL.", RoadName: "UJI RAYA",
		BuildingNo: "7", ZipCode: "12345", RiskLocation: "JL. UJI RAYA 7", Territory: "UJI KEL", City: "UJI KOTA", District: "UJI KEC",
		Province: "UJI PROV", Occupations: []models.OkupasiObjek{{OccupationID: "2001"}, {OccupationID: "2002"}},
		Items: []models.ItemObjek{{ItemType: "BUILDING", Currency: "IDR", Coverages: []models.CoverageObjek{
			{Coverage: "100815", OldID: "FLEXAS", CoverageNote: "FLEXAS (PSAKI)", AccumulationCode: "UJI-AKU", Rate: d("0.25")},
			{Coverage: "100825", OldID: "4.3", CoverageNote: "FLOOD"}}}}}}
}

func spreadingFacInOffer() models.KasusSpreading {
	return models.KasusSpreading{Lokasi: []models.LokasiSpreading{{Items: []models.ItemSpreading{{Coverages: []models.CoverageSpreading{
		{OldID: "FLEXAS", Spreading: []models.BarisSpreading{
			{TreatyType: "10007", SharePercentage: d("60"), TSISpreaded: d("600"), PremiumSpreaded: d("6")},
			{TreatyType: "10011", SharePercentage: d("40"), TSISpreaded: d("400"), PremiumSpreaded: d("4")}}},
		{OldID: "4.3"}}}}}}}
}

// TestBarisFacInOffer - SaveOfferProduction_Act cabang FIRE: satu baris per spreading (coverage tanpa spreading tidak menulis),
// ROADNAME = RoadType + RoadName (tanpa spasi), BUILDINGNO "NO. " + nomor, OCCUPATIONCODE yang baru di depan, tanggal WIB
// 24-jam, TSI / Premi spreading; objek dan spreading tidak sejalan -> ErrBindingTidakSejalan.
func TestBarisFacInOffer(t *testing.T) {
	k := models.KonteksSubmit{MOID: "UJI-MO", InsuredName: "UJI TERTANGGUNG", SobName: "UJI SOB", PercentShare: d("12.5"),
		RICommision: d("10")}
	g := models.General{ReffNumber: "UJI-SLIP", CedingCoName: "UJI CEDING", GroupName: "UJI GRUP",
		StartDateTime: "20261019T170000.000 GMT", EndDateTime: "20271020T080000.000 GMT"}
	baris, err := barisFacInOffer("UJI-NB-1", k, g, "UJI-BIZ", "2", objekFacInOffer(), spreadingFacInOffer())
	if err != nil || len(baris) != 2 {
		t.Fatalf("%v %d", err, len(baris))
	}
	b := baris[0]
	mau := models.BarisFacInOffer{IDPega: "UJI-NB-1", BusinessCode: "UJI-BIZ", CedingCo: "UJI CEDING", SOB: "UJI SOB",
		BeginDate: "20/10/2026 00:00:00", EndDate: "20/10/2027 15:00:00", Province: "UJI PROV", Kabupaten: "UJI KOTA",
		Kecamatan: "UJI KEC", Kelurahan: "UJI KEL", RoadName: "JL.UJI RAYA", BuildingNo: "NO. 7", ZipCode: "12345", ObjectNo: "1",
		ObjectName: "UJI OBJEK", RiskLocation: "JL. UJI RAYA 7", ObjectItem: "BUILDING", Coverage: "100815",
		CoverageNote: "FLEXAS (PSAKI)", Rate: "0.25", JnReas: "10007", CurrID: "IDR", AccumulationCode: "UJI-AKU",
		OccupationCode: "2002 || 2001", MarketingOfficerCode: "UJI-MO", InsuredName: "UJI TERTANGGUNG", RISlipCeding: "UJI-SLIP",
		GroupType: "2", GroupName: "UJI GRUP", TopRisk: "true"}
	salin := b
	salin.PercentRNM, salin.PercentShareSpreaded, salin.TSI, salin.Premi, salin.PercentRIComm = nil, nil, nil, nil, nil
	if !reflect.DeepEqual(salin, mau) || t8(b.PercentRNM) != "12.5" || t8(b.PercentRIComm) != "10" || t8(b.PercentShareSpreaded) != "60" ||
		t8(b.TSI) != "600" || t8(b.Premi) != "6" {
		t.Errorf("%+v", b)
	}
	if baris[1].JnReas != "10011" || t8(baris[1].TSI) != "400" || baris[1].Coverage != "100815" {
		t.Errorf("baris 2 %+v", baris[1])
	}
	beda := spreadingFacInOffer()
	beda.Lokasi[0].Items[0].Coverages[1].OldID = "EQ"
	kurang := spreadingFacInOffer()
	kurang.Lokasi[0].Items[0].Coverages = kurang.Lokasi[0].Items[0].Coverages[:1]
	for nama, sp := range map[string]models.KasusSpreading{"oldid": beda, "coverage": kurang, "lokasi": {}} {
		if _, err := barisFacInOffer("UJI-NB-1", k, g, "", "2", objekFacInOffer(), sp); !errors.Is(err, ErrBindingTidakSejalan) {
			t.Errorf("%s: %v", nama, err)
		}
	}
	if tanggalFacInOffer("") != "" || tanggalFacInOffer("bukan") != "" || kodeOkupasiFacInOffer(nil) != "" {
		t.Error("kosong")
	}
}

// TestSubmitBinding - Accept Direktur Teknik (Decision23 Else) menulis Binding: flag 2, JSON minimal {"IDNewBisnis": id}, bulan
// nomor offer = bulan WIB jam Submit, baris FACINOFFER dari objek + spreading case; email hanya di produksi dengan SMTP.
func TestSubmitBinding(t *testing.T) {
	ctx := context.Background()
	for _, produksi := range []bool{false, true} {
		st := submitUji()
		st.pembuat, st.ccBinding = "uji@nusantarare.com", []string{"uji-mo@contoh.test"}
		sur := &surelTiruan{}
		svc := layananSubmitUW(st, posisiDirTeknik, spreadingTanpaFacRetro("12000")).DenganGrup(KonfigGrup{}, produksi).DenganPengirimSurel(sur)
		bulan := "10.2026" // jam layananSubmit 07-10-2026 22:04 WIB
		if produksi {
			// 31-10-2026 18:00 UTC = 01-11-2026 01:00 WIB: bulan nomor offer menurut WIB.
			svc.DenganJam(func() time.Time { return time.Date(2026, 10, 31, 18, 0, 0, 0, time.UTC) })
			bulan = "11.2026"
		}
		svc.DenganTableOfLimit(tolTiruan{kode: []string{"UJI-BIZ"}, oldID: []string{oldIDUji}})
		svc.limit = limitTiruan{a: limitLanjut()}
		h, err := svc.AjukanSubmit(ctx, pelakuPosisi(posisiDirTeknik), "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "x"})
		if err != nil || h.Position != repository.WorkbasketMarketing || h.NBStatus != nbBinding || len(st.tulis) != 1 {
			t.Fatalf("%v %+v", err, h)
		}
		b := st.tulis[0].Binding
		if b == nil || b.FlagBaru != "2" || b.JSONData != `{"IDNewBisnis":"UJI-NB-1"}` || b.BulanNoOffer != bulan || len(b.Offer) != 1 {
			t.Fatalf("%+v", b)
		}
		o := b.Offer[0]
		if o.IDPega != "UJI-NB-1" || o.BusinessCode != "UJI-BIZ" || o.JnReas != "10007" || o.Coverage != "FLEXAS" || o.CurrID == "" ||
			o.AccumulationCode != "UJI-AKU" || o.OccupationCode != "2001" || o.BuildingNo != "NO. " || o.MarketingOfficerCode != "UJI-MO" ||
			t8(o.PercentRNM) != "12.5" || t8(o.PercentRIComm) != "10" || o.GroupType != "2" || o.InsuredName != "UJI TERTANGGUNG" {
			t.Errorf("%+v", o)
		}
		if !produksi {
			if len(sur.kirim) != 0 {
				t.Errorf("DEV mengirim email")
			}
			continue
		}
		if len(sur.kirim) != 1 || !reflect.DeepEqual(sur.kirim[0].To, []string{"uji@nusantarare.com"}) ||
			!reflect.DeepEqual(sur.kirim[0].Cc, []string{"uji-mo@contoh.test"}) ||
			sur.kirim[0].Subject != "Waiting for Bind ID Pega : UJI-NB-1 - UJI TERTANGGUNG" {
			t.Errorf("%+v", sur.kirim)
		}
	}
}

// TestSurelBinding - SendEmailPolicy langkah 5: pembuat bukan alamat internal / kosong -> To = CC; keduanya kosong -> tidak
// dikirim (langkah 12); badan di-escape.
func TestSurelBinding(t *testing.T) {
	k := models.KonteksSubmit{InsuredName: "UJI <B>"}
	for nama, c := range map[string]struct {
		pembuat string
		cc      []string
		mauTo   []string
		ok      bool
	}{
		"internal":       {"uji@nusantarare.com", []string{"a@contoh.test"}, []string{"uji@nusantarare.com"}, true},
		"bukan internal": {"uji@contoh.test", []string{"a@contoh.test", "b@contoh.test"}, []string{"a@contoh.test", "b@contoh.test"}, true},
		"pembuat kosong": {"", []string{"a@contoh.test"}, []string{"a@contoh.test"}, true},
		"tanpa penerima": {"", nil, nil, false},
	} {
		s, ok := surelBinding("UJI-NB-1", k, "NB IS IN UJI'S INBOX", c.pembuat, c.cc)
		if ok != c.ok || (ok && (!reflect.DeepEqual(s.To, c.mauTo) || !reflect.DeepEqual(s.Cc, c.cc) ||
			s.HTML != "<p>Waiting for Bind ID Pega : UJI-NB-1 - UJI &lt;B&gt;</p><p>NB IS IN UJI'S INBOX</p>")) {
			t.Errorf("%s: %v %+v", nama, ok, s)
		}
	}
}
