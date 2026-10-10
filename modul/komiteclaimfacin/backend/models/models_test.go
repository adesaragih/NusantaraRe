package models_test

// Uji model murni Komite Claim Fac In: JSON halaman Pega bersarang (DATA_JSON OS DEV), perluasan tangga KCF-02 dan pita
// SPV B KCF-01 (ApprovalKomite_Act S3-S8), pra-proses SetValueKomite (total per mata uang + Total in IDR), rencana
// KomitePost_Adjustment (KomiteCount S14 / S24, IsKomiteLoop), wewenang tingkat berjalan, subjek email. Fixture `UJI-`.

import (
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimfacin/backend/models"
	"nusantarare/modul/komiteclaimfacin/backend/tiruan"
)

var saat = time.Date(2026, 10, 10, 9, 30, 0, 0, models.Jakarta)

func TestJSONHalamanPegaBersarang(t *testing.T) {
	got := models.JSONHalamanPega(models.HalamanJSON{Nilai: map[string]string{"pxObjClass": "K", "Type": "1",
		"AcceptedNo": "A", "Kosong": ""}, Daftar: []models.DaftarJSON{{Nama: "CurrencyList", Isi: []models.HalamanJSON{
		{Nilai: map[string]string{"Currency": "IDR"}, Daftar: []models.DaftarJSON{{Nama: "AcceptationList",
			Isi: []models.HalamanJSON{{Nilai: map[string]string{"Value": "1"}}, {Nilai: map[string]string{"Value": "2"}}}}}}}}}})
	want := "{\n\"AcceptedNo\":\"A\"\n,\"pxObjClass\":\"K\"\n,\"Type\":\"1\"\n,\"CurrencyList\":[ \n{\n\"Currency\":\"IDR\"" +
		"\n,\"AcceptationList\":[ \n{\n\"Value\":\"1\"\n}\n,{\n\"Value\":\"2\"\n}\n] \n}\n] \n}\n"
	if got != want {
		t.Fatalf("JSON halaman Pega:\n got %q\nwant %q", got, want)
	}
}

func TestNomorAkseptasiTanpaAkhiranLini(t *testing.T) {
	no := models.RakitNomorAkseptasi("UJI-A", "12", "10.2026", 7)
	if no != "UJI-A12.10.2026.00007" || !models.PanjangNoAksepCLM(no) || strings.Contains(no, ".TP") || strings.Contains(no, ".TX") {
		t.Fatalf("nomor %q", no)
	}
}

func kasusTT2(tangga ...models.Anggota) models.Kasus {
	return models.Kasus{ID: "KMT-UJI1", TransferType: models.TransferAdjustment, Count: 1, Loop: len(tangga), Tangga: tangga}
}

func spv() models.Anggota {
	return models.Anggota{ID: "L1", Urut: 1, OperatorID: models.WorkbasketSPVA, Jabatan: "Claim Supervisor",
		Keputusan: models.KeputusanMenunggu}
}

func total(v int64) models.PraProses { return models.PraProses{TotalAdj: apd.New(v, 0)} }

func TestPerluasTanggaMenurutBatasBawah(t *testing.T) {
	r := tiruan.RosterUji()
	for _, c := range []struct {
		nilai int64
		mau   []string
	}{
		{20000000, nil},
		{57750001, nil}, // LIMIT_BOTTOM < total (ketat)
		{57750002, []string{"Claim Dept. Head"}},
		{200000000, []string{"Claim Dept. Head", "Technic Div. Head"}},
		{700000000, []string{"Claim Dept. Head", "Technic Div. Head", "Operational Director", "Technical Director"}},
	} {
		baru, err := models.PerluasTangga(kasusTT2(spv()), total(c.nilai), r, false)
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for i, a := range baru {
			got = append(got, a.Jabatan)
			if a.Urut != i+2 || a.Keputusan != models.KeputusanMenunggu {
				t.Fatalf("%d: anggota %+v", c.nilai, a)
			}
		}
		if strings.Join(got, ",") != strings.Join(c.mau, ",") {
			t.Fatalf("%d: perluasan %v mau %v", c.nilai, got, c.mau)
		}
	}
}

func TestPerluasTanggaPitaSPVB(t *testing.T) {
	r := tiruan.RosterUji()
	for _, c := range []struct {
		nilai int64
		spvb  bool
		mau   int
	}{
		{30000000, true, 0}, // > 30 jt eksklusif
		{30000001, true, 1},
		{57750000, true, 1}, // <= 57,75 jt inklusif
		{57750001, true, 0}, // di luar pita: aturan S4 (LIMIT_BOTTOM 57.750.001 tidak < total)
		{40000000, false, 0},
	} {
		baru, err := models.PerluasTangga(kasusTT2(spv()), total(c.nilai), r, c.spvb)
		if err != nil {
			t.Fatal(err)
		}
		if len(baru) != c.mau || (c.mau == 1 && baru[0].OperatorID != models.WorkbasketDeptHead) {
			t.Fatalf("%d spvb=%v: %+v", c.nilai, c.spvb, baru)
		}
	}
}

func TestPerluasTanggaTidakBerlaku(t *testing.T) {
	r := tiruan.RosterUji()
	pr := total(200000000)
	retro := pr
	retro.Retro = true
	if b, _ := models.PerluasTangga(kasusTT2(spv()), retro, r, false); b != nil { // S3
		t.Fatalf("Fac Retro diperluas: %+v", b)
	}
	k := kasusTT2(spv(), models.Anggota{ID: "L2", Urut: 2, OperatorID: models.WorkbasketDeptHead})
	if b, _ := models.PerluasTangga(k, pr, r, false); b != nil { // S7 tangga sudah dua baris
		t.Fatalf("tangga dua baris diperluas: %+v", b)
	}
	k = kasusTT2(spv())
	k.Count = 2
	if b, _ := models.PerluasTangga(k, pr, r, false); b != nil { // S4 KomiteCount != 1
		t.Fatalf("tingkat 2 diperluas: %+v", b)
	}
	// S6.2: calon ber-OPERATOR_ID = KomiteList(1).KomiteID dibuang
	satu := models.Anggota{ID: "L1", Urut: 1, OperatorID: models.WorkbasketDeptHead, Jabatan: "Claim Dept. Head",
		Keputusan: models.KeputusanMenunggu}
	b, _ := models.PerluasTangga(kasusTT2(satu), pr, r, false)
	if len(b) != 1 || b[0].OperatorID != "ReasClaimTechDivHead" {
		t.Fatalf("S6.2: %+v", b)
	}
}

func TestSetValueKomiteTotal(t *testing.T) {
	n, d := tiruan.KlaimFacInUji("20000000")
	d[models.DaftarAdj(1, 1)] = append(d[models.DaftarAdj(1, 1)],
		map[string]string{"ID": "UJI-FEE", "PaymentType": "4", "Currency": "USD", "CurrencyID": "USD-ID",
			"GrossAdjustment": "100", "VATValue": "10", "AdjusterFeeValue": "25"},
		map[string]string{"ID": "UJI-TOLAK", "PaymentType": "1", "Currency": "IDR", "CurrencyID": "10026",
			"GrossAdjustment": "999", "AdjustmentValue": "999", "AcceptanceStatus": "2"})
	kl := kontrak.KlaimFacIn{Nilai: n, Daftar: d, Objek: 1, Item: 2, Adjustment: 1}
	k := models.Kasus{ID: tiruan.KomiteUji, AdjustmentID: tiruan.AdjUji, TransferType: models.TransferAdjustment}
	pr, err := models.SetValueKomite(k, kl, map[string]string{"10026": "1", "USD-ID": "15000"})
	if err != nil {
		t.Fatal(err)
	}
	if len(pr.AdjKomite) != 1 || pr.AdjKomite[0]["pyNote"] != "ObjectItem ke 2, Adjustment ke 1" ||
		pr.TotalAdj.String() != "20000000" || pr.Inisial != "UJI Admin" || pr.Okupasi != "UJI OKUPASI" {
		t.Fatalf("DataTempAdj / TotalAdj: %+v total %s", pr.AdjKomite, pr.TotalAdj)
	}
	tot := map[string]map[string]string{}
	for _, b := range pr.Total {
		tot[b["Currency"]] = b
	}
	// IDR: 8 jt (UJI-ADJ-11) + 20 jt, baris status 2 dikecualikan; USD: gross 100 + VAT 10, RNM = fee 25
	if tot["IDR"]["AdjustmentGross"] != "28000000" || tot["IDR"]["AdjustmentValue"] != "22000000" ||
		tot["USD"]["AdjustmentGross"] != "110" || tot["USD"]["AdjustmentValue"] != "25" {
		t.Fatalf("total per mata uang: %+v", tot)
	}
	idr := tot[models.LabelTotalIDR]
	if idr["AdjustmentGross"] != "29650000" || idr["AdjustmentValue"] != "22375000" || pr.Total[len(pr.Total)-1]["Currency"] != models.LabelTotalIDR {
		t.Fatalf("Total in IDR: %+v", idr)
	}
	// TT3 / TT4: S1 keluar
	k.TransferType = models.TransferReject
	if pr, _ := models.SetValueKomite(k, kl, nil); len(pr.AdjKomite) != 0 || len(pr.Total) != 0 {
		t.Fatalf("pra-proses TT3: %+v", pr)
	}
}

func TestRencanaHitungTingkat(t *testing.T) {
	n, d := tiruan.KlaimFacInUji("100000000")
	kl := kontrak.KlaimFacIn{Nilai: n, Daftar: d, Objek: 1, Item: 2, Adjustment: 1}
	dua := models.Anggota{ID: "L2", Urut: 2, OperatorID: models.WorkbasketDeptHead, Jabatan: "Claim Dept. Head",
		Keputusan: models.KeputusanMenunggu}
	k := kasusTT2(spv(), dua)
	kep := models.Keputusan{AcceptStatus: models.KeputusanSetuju, Comment: "UJI", UsulTutup: true}
	r := models.Rencanakan(k, kl, kep, "UJI-SPVA", saat)
	if r.Selesai || r.Kepala.Count != 2 || r.Posisi != models.WorkbasketDeptHead || r.TerbitkanNomor || r.SimpanOS {
		t.Fatalf("tingkat 1 dari 2: %+v", r)
	}
	if r.Klaim.Header["ClaimData.IsCloseFile"] != "1" || r.Kepala.UsulTutup != "1" || r.Klaim.Header["AktifButton"] != "" {
		t.Fatalf("usul tingkat 1: %+v %+v", r.Klaim.Header, r.Kepala)
	}
	k.Count, k.UsulTutup, k.Tangga[0].Keputusan = 2, "1", models.KeputusanSetuju
	r = models.Rencanakan(k, kl, models.Keputusan{AcceptStatus: models.KeputusanSetuju, Comment: "UJI2"}, "UJI-HEAD", saat)
	if !r.Selesai || r.Kepala.Count != 3 || !r.TerbitkanNomor || !r.SimpanOS || r.StatusRiwayat != "ACCEPT" ||
		r.SubProgres != "Accepted" || r.Email != models.EmailPembuatSetuju || r.Klaim.Header["ClaimData.IsCloseFile"] != "1" {
		t.Fatalf("tingkat akhir: %+v", r)
	}
	// tolak di tingkat 1: KomiteCount := KomiteLoop lalu +1, sisa tangga ditolak tanpa komentar
	k = kasusTT2(spv(), dua)
	r = models.Rencanakan(k, kl, models.Keputusan{AcceptStatus: models.KeputusanTolak, Comment: "UJI"}, "UJI-SPVA", saat)
	if !r.Selesai || r.Kepala.Count != 3 || len(r.Tangga) != 2 || r.Tangga[1].IsiKomentar || r.SubProgres != "Rejected" ||
		r.Klaim.Adjustment["IsApproved"] != "" || r.Email != models.EmailPembuatTolak {
		t.Fatalf("tolak tingkat 1: %+v", r)
	}
}

func TestPemegangTingkatBerjalan(t *testing.T) {
	k := kasusTT2(spv(), models.Anggota{ID: "L2", Urut: 2, OperatorID: models.WorkbasketDeptHead,
		Keputusan: models.KeputusanMenunggu})
	for _, c := range []struct {
		akun  string
		peran []string
		mau   bool
	}{
		{"UJI-A", []string{models.WorkbasketSPVA}, true},
		{"UJI-B", []string{models.WorkbasketSPVB}, true}, // cadangan SPV A (KCF-01)
		{"UJI-H", []string{models.WorkbasketDeptHead}, false},
		{"UJI-X", nil, false},
		{models.WorkbasketSPVA, nil, true}, // KomiteID = akun
	} {
		if got := k.Pemegang(c.akun, c.peran); got != c.mau {
			t.Fatalf("%s %v: %v mau %v", c.akun, c.peran, got, c.mau)
		}
	}
	k.Count, k.Tangga[0].Keputusan = 2, models.KeputusanSetuju
	if k.Pemegang("UJI-B", []string{models.WorkbasketSPVB}) || !k.Pemegang("UJI-H", []string{models.WorkbasketDeptHead}) {
		t.Fatal("SPV B hanya cadangan tingkat ber-KomiteID SPV A")
	}
	if p := models.PeranKerja([]string{models.WorkbasketSPVB}); len(p) != 2 || p[1] != models.WorkbasketSPVA {
		t.Fatalf("PeranKerja: %v", p)
	}
}

func TestSubjekSurel(t *testing.T) {
	for _, c := range []struct{ tt, jenis, mau string }{
		{models.TransferAdjustment, models.EmailPenyetujuBerikut, "Pengajuan Akseptasi : CLM-1/KMT-1 UJI DOL 05 Febuari 2026"},
		{models.TransferAdjustment, models.EmailPembuatSetuju, "(Approval) Pengajuan Akseptasi : CLM-1/KMT-1 UJI DOL 05 Febuari 2026"},
		{models.TransferAdjustment, models.EmailPembuatTolak, "(Reject) Pengajuan Akseptasi : CLM-1/KMT-1 UJI DOL 05 Febuari 2026"},
		{models.TransferReject, models.EmailPembuatSetuju, "(Approval) Pengajuan Reject Klaim : CLM-1/KMT-1 UJI DOL 05 Febuari 2026"},
		{models.TransferClose, models.EmailPembuatTolak, "(Reject) Pengajuan Close Klaim : CLM-1/KMT-1 UJI DOL 05 Febuari 2026"},
	} {
		if got := models.SubjekSurel(c.tt, c.jenis, "CLM-1", "KMT-1", "UJI", "2026-02-05"); got != c.mau {
			t.Fatalf("%s/%s: %q", c.tt, c.jenis, got)
		}
	}
	if got := models.TanggalSurat("20260715"); got != "15 July 2026" {
		t.Fatalf("bulan VERBATIM: %q", got)
	}
}

func TestOSTolakPerEstimasiCetak(t *testing.T) {
	n, d := tiruan.KlaimFacInUji("10")
	d[models.DaftarDiItem(1, 2, models.AnakEstimasi)][0]["PrintFaceClaim"] = "1"
	d[models.DaftarDiItem(1, 2, models.AnakEstimasi)][0]["EstimationValue"] = ""
	n["pyID"] = "CLM-UJI"
	rows := models.SusunOSTolak(kontrak.KlaimFacIn{Nilai: n, Daftar: d}, saat)
	if len(rows) != 2 || rows[0].StsKonversi != "" || rows[1].StsKonversi != "1" || rows[0].StsDLA != "8" {
		t.Fatalf("OS STS 2: %+v", rows)
	}
	if !strings.Contains(rows[1].DataJSON, `"ObjectItemName":"UJI ITEM 2"`) || !strings.Contains(rows[0].DataJSON,
		`"PersenRNM":"25"`) || !strings.Contains(rows[0].DataJSON, `"EstimationDate":"20261010"`) {
		t.Fatalf("DATA_JSON STS 2:\n%s\n%s", rows[0].DataJSON, rows[1].DataJSON)
	}
}

func TestPeriodePolisBerformatTanggal(t *testing.T) {
	// ShowTransfer LS8: StartDateTime "-" EndDateTime, keduanya FormatType=date (dd-mm-yyyy di layar)
	kl := kontrak.KlaimFacIn{Nilai: map[string]string{"OfferFacIn.PolicyData.StartDateTime": "2026-01-01",
		"OfferFacIn.PolicyData.EndDateTime": "2026-12-31"}, Daftar: map[string][]map[string]string{}}
	ly := models.SusunLayar(models.Kasus{ID: "KMT-UJI", TransferType: models.TransferAdjustment, Count: 1, Loop: 1},
		kl, models.PraProses{}, nil, nil, "", nil)
	for _, b := range ly.Bagian {
		for _, m := range b.Medan {
			if m.Label == "Period" {
				if m.Nilai != "01-01-2026 - 31-12-2026" {
					t.Fatalf("Period %q", m.Nilai)
				}
				return
			}
		}
	}
	t.Fatal("medan Period tidak ada")
}
