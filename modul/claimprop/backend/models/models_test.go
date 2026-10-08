package models_test

// Uji seam 1 (aturan murni): angka harapan dihitung tangan dari rumus Pega yang dikutip di komentar port, bukan dari
// kode yang diuji.

import (
	"context"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/tiruan"
)

func halamanUji() *models.Halaman {
	h := models.HalamanBaru()
	h.Setel(models.CD+"ShareCeding", "50")
	h.Setel(models.CD+"NetDeductibleValue", "100")
	h.Setel(models.TM+"RNMShareP", "10")
	h.SetelDaftar(models.DaftarInterest, []models.Baris{{"ObjectName": "UJI-OBJ", "CurrencyID": "UJI-A",
		"Currency": "UJA", "KursObjectItem": "2", "TSIPerObject": "1000"}})
	h.SetelDaftar(models.DaftarClaimAmount, []models.Baris{{"CurrencyID": "UJI-A", "Currency": "UJA", "IDR": "2",
		"ClaimAmount": "1000"}})
	h.SetelDaftar(models.DaftarLossAlloc, []models.Baris{{"CurrencyID": "UJI-A", "SharePercentage": "25",
		"PremiumSpreaded": "2", "TreatyType": "UJI-QS"}})
	h.SetelDaftar(models.DaftarEstimasi, []models.Baris{{"CurrencyID": "UJI-A", "Currency": "UJA", "KursValue": "2",
		"GrossEstimationPct": "1000", "Type": "1", "TypeLossID": "UJI-QS"}})
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{{"CurrencyID": "UJI-A", "SharePercentage": "60"},
		{"CurrencyID": "UJI-A", "SharePercentage": "40"}})
	h.SetelDaftar(models.DaftarBreakQS, []models.Baris{{"CurrencyID": "UJI-A", "SharePercentage": "50"}})
	return h
}

func sama(t *testing.T, nama, dapat, mau string) {
	t.Helper()
	if dapat != mau {
		t.Errorf("%s = %q, mau %q", nama, dapat, mau)
	}
}

func TestHitungTurunanRantaiKlaim(t *testing.T) {
	h := halamanUji()
	if err := models.HitungTurunan(h); err != nil {
		t.Fatal(err)
	}
	in := h.AmbilDaftar(models.DaftarInterest)[0]
	sama(t, "TSIPerObjectIDR", in["TSIPerObjectIDR"], "2000") // 2 x 1000
	sama(t, "TotalSumInsuredIDR", h.Ambil(models.CD+"TotalSumInsuredIDR"), "2000")
	sama(t, "TotalInterestInsured", h.AmbilDaftar(models.DaftarTotalTSI)[0]["Value"], "1000")
	ca := h.AmbilDaftar(models.DaftarClaimAmount)[0]
	sama(t, "ClaimAmount.Value", ca["Value"], "400") // 50% x 1000 - 100 (rumus 2024)
	sama(t, "ClaimAmount.USD", ca["USD"], "800")
	sama(t, "TotalListClaimAmountIDR", h.Ambil(models.CD+"TotalListClaimAmountIDR"), "800")
	la := h.AmbilDaftar(models.DaftarLossAlloc)[0]
	sama(t, "LossAlloc.ClaimSpreaded", la["ClaimSpreaded"], "100") // 25% x 400
	sama(t, "LossAlloc.ClaimEstimation", la["ClaimEstimation"], "200")
	es := h.AmbilDaftar(models.DaftarEstimasi)[0]
	sama(t, "EstimationValue", es["EstimationValue"], "100") // 10% x 1000
	sama(t, "ConvertValue", es["ConvertValue"], "200")
	sama(t, "ConvertGrossEstimasi", es["ConvertGrossEstimasi"], "2000")
	sama(t, "TotalEstimasiIDR", h.Ambil(models.CD+"TotalEstimasiIDR"), "200")
	sama(t, "TotalGrossEstimateTreaty", h.Ambil(models.CD+"TotalGrossEstimateTreaty"), "1000")
	sp := h.AmbilDaftar(models.DaftarSpreading)
	sama(t, "Spreading1", sp[0]["ClaimSpreaded"], "60") // 60% x 100
	sama(t, "Spreading2", sp[1]["ClaimSpreaded"], "40")
	sama(t, "BreakQS", h.AmbilDaftar(models.DaftarBreakQS)[0]["ClaimSpreaded"], "20") // 50% x spread TERAKHIR (40)
}

func TestBarisBekuTidakDihitungUlang(t *testing.T) {
	h := halamanUji()
	h.AmbilDaftar(models.DaftarClaimAmount)[0]["Note"] = "Yes"
	h.AmbilDaftar(models.DaftarClaimAmount)[0]["Value"] = "123"
	h.AmbilDaftar(models.DaftarEstimasi)[0]["PrintFaceClaim"] = "1"
	h.AmbilDaftar(models.DaftarEstimasi)[0]["EstimationValue"] = "7"
	if err := models.HitungTurunan(h); err != nil {
		t.Fatal(err)
	}
	sama(t, "ClaimAmount beku", h.AmbilDaftar(models.DaftarClaimAmount)[0]["Value"], "123")
	sama(t, "Estimasi terkirim", h.AmbilDaftar(models.DaftarEstimasi)[0]["EstimationValue"], "7")
	sama(t, "LossAlloc dari nilai beku", h.AmbilDaftar(models.DaftarLossAlloc)[0]["ClaimSpreaded"], "30.75") // 25% x 123
}

func TestPresisiPenuhKaliDuluBagiTerakhir(t *testing.T) {
	h := models.HalamanBaru()
	h.Setel(models.TM+"RNMShareP", "33.3333333333")
	h.SetelDaftar(models.DaftarEstimasi, []models.Baris{{"CurrencyID": "UJI-A", "KursValue": "1", "GrossEstimationPct": "3"}})
	if err := models.HitungTurunan(h); err != nil {
		t.Fatal(err)
	}
	sama(t, "EstimationValue", h.AmbilDaftar(models.DaftarEstimasi)[0]["EstimationValue"], "0.999999999999")
}

func TestFormatNomor(t *testing.T) {
	sama(t, "klaim", models.RakitNomorKlaim("UJI-K", "12", "05.2026", 7), "UJI-K12.05.2026.T00007")
	sama(t, "sementara", models.RakitNomorSementara("2026", 12), "KT.2026-00012")
	sama(t, "CFS", models.RakitNomorCFS("", "05", "2026", 1168), "RNM-K.05.2026.T01168")
	sama(t, "PLA", models.RakitNomorPLA("UJI-H", "12", "05.2026", 3), "UJI-H12.05.2026.00003")
}

func TestRencanaNomor(t *testing.T) {
	h := models.HalamanBaru()
	if r := models.RencanakanNomor(h); !r.Sementara || r.Klaim {
		t.Errorf("tanpa polis: %+v, mau nomor sementara saja", r)
	}
	h.Setel(models.CD+"PolicyData.PolicyNo", "UJI-RNM-Q1")
	if r := models.RencanakanNomor(h); r.Sementara || !r.Klaim {
		t.Errorf("berpolis: %+v, mau nomor klaim saja", r)
	}
	h.Setel(models.CD+"NoClaim", "UJI-K1")
	if r := models.RencanakanNomor(h); r.Sementara || r.Klaim {
		t.Errorf("bernomor: %+v, mau tanpa nomor baru", r)
	}
}

func TestTglBolehBayar(t *testing.T) {
	sama(t, "hari <= 25", models.TglBolehBayar("20260520", 5), "20-06-2026")
	sama(t, "hari > 25 Desember", models.TglBolehBayar("20261228", 12), "01-01-2027")
	sama(t, "November", models.TglBolehBayar("20261120", 11), "20-12-2026")
}

func konteksUji(a *tiruan.Acuan) *models.Konteks {
	return &models.Konteks{Ctx: context.Background(), Acuan: a, Pelaku: "UJI-ADMIN",
		Sekarang: time.Date(2026, 5, 20, 10, 0, 0, 0, models.Jakarta)}
}

func TestKirimEstimasiMenandaiDanMelahirkanBarisOS(t *testing.T) {
	h := halamanUji()
	h.Setel(models.CD+"NoClaim", "UJI-K1")
	h.AmbilDaftar(models.DaftarEstimasi)[0]["EstimationValue"] = "100"
	h.TambahBaris(models.DaftarEstimasi, models.Baris{"CurrencyID": "UJI-A", "PrintFaceClaim": "1"})
	rows := models.KirimEstimasi(konteksUji(tiruan.AcuanBaru()), h, "CLMP-000001")
	if len(rows) != 1 {
		t.Fatalf("baris OS = %d, mau 1 (estimasi terkirim dilewati)", len(rows))
	}
	sama(t, "STS_REJECT", rows[0].StsReject, models.StsOSEstimasi)
	sama(t, "NOCLAIM", rows[0].NoClaim, "UJI-K1")
	sama(t, "VALUE", rows[0].Value, "100")
	sama(t, "PrintFaceClaim", h.AmbilDaftar(models.DaftarEstimasi)[0]["PrintFaceClaim"], "1")
	models.BekukanDataLama(h)
	sama(t, "Note", h.AmbilDaftar(models.DaftarClaimAmount)[0]["Note"], "Yes")
	sama(t, "IsOldData", h.AmbilDaftar(models.DaftarSpreading)[1]["IsOldData"], "Yes")
	sama(t, "IsAdjVal", h.AmbilDaftar(models.DaftarInterest)[0]["IsAdjVal"], "Yes")
}

func TestDuplikatDOLMemblokirKecualiKlaimDitolak(t *testing.T) {
	a := tiruan.AcuanBaru()
	a.Riwayat["UJI-RNM-Q1"] = []models.RiwayatKlaimPolis{{IDKasus: "UJI-CLMP-9", DateOfLoss: "2026-05-01 08:00:00"}}
	h := models.HalamanBaru()
	h.Setel("pyID", "CLMP-000001")
	h.Setel(models.CD+"PolicyData.PolicyNo", "UJI-RNM-Q1")
	h.Setel(models.CD+"DateOfLoss", "2026-05-01 15:00:00")
	if err := models.PeriksaDOLSerupa(konteksUji(a), h); err != nil {
		t.Fatal(err)
	}
	if h.Ambil(models.PropBlokirDOL) != "1" || !strings.Contains(strings.Join(h.SemuaPesan(), ";"), models.PesanDOLSerupa) {
		t.Fatalf("duplikat DOL tidak memblokir: %v", h.SemuaPesan())
	}
	a.Riwayat["UJI-RNM-Q1"][0].Ditolak = true
	h.BersihkanPesan()
	if err := models.PeriksaDOLSerupa(konteksUji(a), h); err != nil {
		t.Fatal(err)
	}
	if h.Ambil(models.PropBlokirDOL) != "" || h.AdaPesan() {
		t.Fatalf("klaim lama yang ditolak tetap memblokir: %v", h.SemuaPesan())
	}
}

func TestShareSpreadingTepatSeratus(t *testing.T) {
	h := models.HalamanBaru()
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{{"CurrencyID": "UJI-A", "Currency": "UJA", "SharePercentage": "60"},
		{"CurrencyID": "UJI-A", "Currency": "UJA", "SharePercentage": "30"}})
	p, err := models.PeriksaShareTepat100(h)
	if err != nil || len(p) != 1 || !strings.Contains(p[0], models.PesanShareKurang100) {
		t.Fatalf("90 %%: %v %v, mau pesan kurang dari 100", p, err)
	}
	h.AmbilDaftar(models.DaftarSpreading)[1]["SharePercentage"] = "40"
	if p, _ := models.PeriksaShareTepat100(h); len(p) != 0 {
		t.Fatalf("100 %%: %v, mau lolos", p)
	}
	h.AmbilDaftar(models.DaftarSpreading)[1]["SharePercentage"] = "41"
	if p, _ := models.PeriksaShareTepat100(h); len(p) != 1 || !strings.Contains(p[0], models.PesanShareLebih100) {
		t.Fatalf("101 %%: %v, mau pesan lebih dari 100", p)
	}
}

func TestTataMedanTerkunciTidakTerbuka(t *testing.T) {
	h := models.HalamanBaru()
	h.Setel("IsOutstanding", "1")
	ts := models.Evaluasi(h, models.LayarOutstanding(), false)
	terbuka := models.MedanTerbuka(ts)
	if terbuka[models.CD+"NoClaim"] || terbuka[models.CD+"PolicyData.PolicyNo"] {
		t.Errorf("medan read-only selalu terbuka: %v", terbuka)
	}
	if terbuka[models.CD+"ReporterName"] {
		t.Errorf("Reporter Name read-only jika IsOutstanding==1, tetapi terbuka")
	}
	if !terbuka[models.CD+"PostalCode"] {
		t.Errorf("Zip Code (tanpa kondisi read-only) harus terbuka")
	}
	if models.AksiTerbuka(ts, "PilihMaster", 0) {
		t.Errorf("Choose Master disabled jika IsOutstanding==1, tetapi terbuka")
	}
	if models.AksiTerbuka(ts, "SaveOutstanding", 0) {
		t.Errorf("Save to issue RNM disabled jika IsCFS = '', tetapi terbuka")
	}
	if !models.AksiTerbuka(models.Evaluasi(h, models.LayarOutstanding(), false), "CheckNopolicy", 0) {
		t.Errorf("Send to Acceptation tampil jika IsOutstanding = 1 dan aktif selama IsAcceptation != 1")
	}
	if models.AksiTerbuka(models.Evaluasi(h, models.LayarOutstanding(), true), "CheckNopolicy", 0) {
		t.Errorf("layar terkunci tetap membuka aksi")
	}
}

func TestTataBarisGridMengikutiBarisnya(t *testing.T) {
	h := models.HalamanBaru()
	h.SetelDaftar(models.DaftarInterest, []models.Baris{{"ObjectName": "UJI-1"}, {"ObjectName": "UJI-2", "IsAdjVal": "Yes"}})
	ts := models.Evaluasi(h, models.LayarOutstanding(), false)
	terbuka := models.MedanTerbuka(ts)
	if !terbuka[models.JalurAnak(models.DaftarInterest, 1, "TSIPerObject")] {
		t.Errorf("TSI baris 1 harus terbuka")
	}
	if terbuka[models.JalurAnak(models.DaftarInterest, 2, "TSIPerObject")] {
		t.Errorf("TSI baris 2 (IsAdjVal Yes) harus terkunci")
	}
	if models.AksiTerbuka(ts, "DeleteInterest", 2) || !models.AksiTerbuka(ts, "DeleteInterest", 1) {
		t.Errorf("Delete baris interest mengikuti IsAdjVal")
	}
}

func TestCountValueAdjustment(t *testing.T) {
	a := tiruan.AcuanBaru()
	h := models.HalamanBaru()
	h.Setel(models.CD+"ShareCeding", "50")
	h.SetelDaftar(models.DaftarAdjustment, []models.Baris{{"Type": "1", "CurrencyID": "UJI-A", "KursIDR": "2",
		"PersenRNM": "10", "GrossAdjustment": "1000", "IndividualRiskType": "1", "IndividualRiskPercentage": "10",
		"TotalEstimasiValue": "1000"}})
	h.SetelDaftar(models.JalurAdj(1, models.AnakLossAllocation), []models.Baris{{"CurrencyID": "UJI-A", "SharePercentage": "40"}})
	if err := models.CountGrossAdjTreaty(konteksUji(a), h, 1); err != nil {
		t.Fatal(err)
	}
	b := h.AmbilDaftar(models.DaftarAdjustment)[0]
	sama(t, "GrossValue", b["GrossValue"], "20")                    // 1000 x 10% x 40% x 50%
	sama(t, "IndividualRiskValue", b["IndividualRiskValue"], "100") // 1000 x 10%
	sama(t, "ProposeAdjustmentValue", b["ProposeAdjustmentValue"], "900")
	sama(t, "AdjustmentValue", b["AdjustmentValue"], "18") // 900 x 10% x 40% x 50%
	sama(t, "ValueAdjustment", b["ValueAdjustment"], "36") // kurs 2
	// CountGrossAdjTreaty_Act 5.2 berjalan SEBELUM CountValueADJTreaty_Act menghitung individual risk: TreatyGross = 1000
	// - IndividualRiskValue lama (kosong) -> 1000 x 50% x 40%.
	sama(t, "LossAlloc.ClaimSpreaded", h.AmbilDaftar(models.JalurAdj(1, models.AnakLossAllocation))[0]["ClaimSpreaded"], "200")
}

// cariTata - simpul pertama ber-ID / berjalur `kunci` di pohon tata.
func cariTata(ts []models.Tata, kunci string) *models.Tata {
	for i := range ts {
		if ts[i].ID == kunci || (ts[i].Jenis == models.JenisGrid && ts[i].Jalur == kunci) {
			return &ts[i]
		}
		if t := cariTata(ts[i].Anak, kunci); t != nil {
			return t
		}
	}
	return nil
}

// Add / Delete Spreading Claim AKTIF (keputusan work owner 08-10-2026, membatalkan 07-10-2026 "non aktifkan"): Add
// selalu aktif, Delete nonaktif bila `.IsOldData='Yes'` (pyDisabledWhen XML). Ikon grid bawaan tetap tidak dibangun.
// Treaty Type = dropdown spreading polis, terkunci bila baris sudah ber-TreatyType (dari polis, atau sudah dipilih
// sesudah Add) atau data lama (work owner 08-10-2026 "ini disable aja").
func TestSpreadingTambahHapusDanTreatyTypeTerkunci(t *testing.T) {
	for _, layar := range []func() []models.Unsur{models.LayarOutstanding, models.LayarAkseptasi} {
		h := models.HalamanBaru()
		h.SetelDaftar(models.DaftarSpreading, []models.Baris{
			{"TreatyType": "UJI-INDUK"}, {"TreatyType": ""}, {"TreatyType": "UJI-LAMA", "IsOldData": "Yes"},
		})
		ts := models.Evaluasi(h, layar(), false)
		if !models.AksiTerbuka(ts, "AddSpreading", 0) {
			t.Fatal("Add spreading mau aktif")
		}
		if models.AksiTerbuka(ts, models.DaftarSpreading+"#tambah", 0) {
			t.Fatal("ikon grid bawaan tidak dibangun")
		}
		if !models.AksiTerbuka(ts, "DeleteSpreading", 1) || !models.AksiTerbuka(ts, "DeleteSpreading", 2) ||
			models.AksiTerbuka(ts, "DeleteSpreading", 3) {
			t.Fatal("Delete mau aktif kecuali IsOldData Yes")
		}
		// SetTreatyNameSpreading = aksi kolom Treaty Type: hanya baris yang Treaty Type-nya masih kosong
		if models.AksiTerbuka(ts, "SetTreatyNameSpreading", 1) || !models.AksiTerbuka(ts, "SetTreatyNameSpreading", 2) ||
			models.AksiTerbuka(ts, "SetTreatyNameSpreading", 3) {
			t.Fatal("Treaty Type mau terkunci bila sudah terisi atau data lama")
		}
	}
}

// rosterUji - satu anggota roster EMAILKOMITE lini PROP berbatas 0.
func rosterUji() *tiruan.Acuan {
	a := tiruan.AcuanBaru()
	a.Roster = append(a.Roster, tiruan.AnggotaRoster{Batas: "0", Sts: models.STSKlaimProp,
		AnggotaKomite: models.AnggotaKomite{OperatorID: "UJI-K1", Jabatan: "UJI-JABATAN"}})
	return a
}

func halamanKomiteUji() *models.Halaman {
	h := models.HalamanBaru()
	h.Setel(models.CD+"Occupation", "UJI-OKUPASI")
	h.SetelDaftar(models.DaftarAdjustment, []models.Baris{{"Type": "1", "ProposeAdjustmentValue": "10",
		"ValueAdjustment": "10", "DataCommitteeTreaty.CircumCauseOfLoss": "UJI", "DataCommitteeTreaty.Remarks": "UJI"}})
	return h
}

// Send Claim to Committee (keputusan work owner 07-10-2026, opsi B): tampil bila gerbang lampiran / premi lolos dan
// isian lengkap, nonaktif selama Payable kosong (dis `pyWorkPage.ClaimData.Payable = ”`).
func TestPenyerahanKomiteMengikutiIsian(t *testing.T) {
	h := halamanKomiteUji()
	ts := models.Evaluasi(h, models.LayarKomite(1, true), false)
	if tb := cariTata(ts, "SendClaimToCommittee"); tb == nil || !tb.Nonaktif || models.AksiTerbuka(ts, "AddKomiteTreatyChild", 1) {
		t.Fatalf("Payable kosong: tombol mau tampil nonaktif: %+v", tb)
	}
	h.Setel(models.CD+"Payable", "2")
	if !models.AksiTerbuka(models.Evaluasi(h, models.LayarKomite(1, true), false), "AddKomiteTreatyChild", 1) {
		t.Fatalf("isian lengkap: penyerahan komite mau terbuka")
	}
	if tb := cariTata(models.Evaluasi(h, models.LayarKomite(1, false), false), "SendClaimToCommittee"); tb != nil {
		t.Fatalf("gerbang lampiran / premi gagal: tombol mau tidak tampil")
	}
}

// Grid "Committe Accept Status": roster calon sebelum diserahkan, keputusan tangga sesudahnya.
func TestGridKomiteRosterLaluTangga(t *testing.T) {
	h := halamanKomiteUji()
	k := konteksUji(rosterUji())
	if err := models.SusunKomiteAdjustment(k, h, 1, nil); err != nil {
		t.Fatal(err)
	}
	rows := h.AmbilDaftar(models.JalurAdj(1, "ComiteeClaim"))
	if len(rows) != 1 || rows[0]["IDKomite"] != "UJI-JABATAN" || rows[0][models.PropKeputusanAnggota] != models.ApprovalKomiteMenunggu {
		t.Fatalf("roster calon %+v", rows)
	}
	g := cariTata(models.Evaluasi(h, models.LayarAdjustment(1), false), models.JalurAdj(1, "ComiteeClaim"))
	if g == nil || len(g.Kolom) != 4 || g.Kolom[1].Jalur != models.PropKeputusanAnggota || g.Kolom[3].Jalur != models.PropCatatanKeputusan {
		t.Fatalf("grid komite mau empat kolom Committee Name / Status / Date Approve / Comment: %+v", g)
	}
	tangga := []models.AnggotaKomite{{OperatorID: "UJI-K1", Jabatan: "UJI-JABATAN", Approval: "1", Comment: "UJI-SETUJU",
		TanggalSetuju: "2026-05-04 10:00:00"}}
	if err := models.SusunKomiteAdjustment(k, h, 1, tangga); err != nil {
		t.Fatal(err)
	}
	rows = h.AmbilDaftar(models.JalurAdj(1, "ComiteeClaim"))
	if len(rows) != 1 || rows[0][models.PropKeputusanAnggota] != "1" || rows[0][models.PropCatatanKeputusan] != "UJI-SETUJU" ||
		rows[0][models.PropTanggalKeputusan] != "2026-05-04 10:00:00" {
		t.Fatalf("keputusan tangga %+v", rows)
	}
}

// Kerangka layout Section OutstandingClaim (XML `pyLayoutOtherFormat` / layout group; work owner 08-10-2026 "berikut
// layout lama, ikuti dan rapihkan"): kepala judul, Claim Treaty Inline grid double 8|5, baris Inline "Quarter/Year",
// layout group Tab Interest / Estimation / Spreading, Claim History paging 5, ikon pi-plus / pi-trash / pyEditIcon.
func TestLetakOutstandingIkutXML(t *testing.T) {
	ts := models.Evaluasi(models.HalamanBaru(), models.LayarOutstanding(), false)
	if ts[0].Letak != models.LetakJudul || ts[0].Anak[0].Label != "Outstanding Claim" {
		t.Fatalf("kepala layar: %+v", ts[0])
	}
	var cari func(ts []models.Tata, ok func(models.Tata) bool) *models.Tata
	cari = func(ts []models.Tata, ok func(models.Tata) bool) *models.Tata {
		for i := range ts {
			if ok(ts[i]) {
				return &ts[i]
			}
			if x := cari(ts[i].Anak, ok); x != nil {
				return x
			}
		}
		return nil
	}
	treaty := cari(ts, func(x models.Tata) bool { return x.Label == "Claim Treaty" })
	kolom := cari(treaty.Anak, func(x models.Tata) bool { return x.Letak == models.LetakDua })
	if kolom == nil || len(kolom.Anak) != 2 || len(kolom.Anak[0].Anak) != 8 || len(kolom.Anak[1].Anak) != 5 {
		t.Fatalf("Claim Treaty bukan Inline grid double 8|5: %+v", kolom)
	}
	qy := cari(ts, func(x models.Tata) bool { return x.Letak == models.LetakSebaris && x.Label == "Quarter/Year" })
	if qy == nil || len(qy.Anak) != 6 {
		t.Fatalf("baris Quarter/Year: %+v", qy)
	}
	tab := cari(ts, func(x models.Tata) bool { return x.Letak == models.LetakTab })
	var judul []string
	for _, a := range tab.Anak {
		judul = append(judul, a.Label)
	}
	if strings.Join(judul, "|") != "Interest|Estimation|Spreading" {
		t.Fatalf("tab: %v", judul)
	}
	if g := cariTata(ts, models.DaftarRiwayatTampil); g == nil || g.PerHalaman != 5 {
		t.Fatalf("Claim History paging: %+v", g)
	}
	if g := cariTata(ts, models.DaftarInterest); g.Tambah.Ikon != models.IkonTambah || g.Kolom[len(g.Kolom)-1].Ikon != models.IkonHapus {
		t.Fatalf("ikon grid Interest: %+v / %+v", g.Tambah, g.Kolom[len(g.Kolom)-1])
	}
	if b := cariTata(ts, "EditRNMShare"); b == nil || b.Ikon != models.IkonUbah {
		t.Fatalf("ikon EditRNMShare: %+v", b)
	}
}

// Generasi polis: PRODKE 0 = NB Treaty In, PRODKE >= 1 = generasi endorsemen EDM Treaty In (T_GENERAL_POLIS_TREATY).
func TestModulBerkasPolis(t *testing.T) {
	for prodke, mau := range map[int]string{0: models.ModulNBTreatyIn, 1: models.ModulEDMTreatyIn, 3: models.ModulEDMTreatyIn} {
		if got := models.ModulBerkasPolis(prodke); got != mau {
			t.Fatalf("PRODKE %d: %s, mau %s", prodke, got, mau)
		}
	}
}

// Perbaikan Claim Information (work owner 08-10-2026): Date of Loss / Received Date tanggal saja, nomor telepon hanya
// angka, Report Type 1-5 dan Reporter Status 1-3 berlabel, "Policy Period TBA ?".
func TestClaimInformationPerbaikanWO(t *testing.T) {
	ts := models.Evaluasi(models.HalamanBaru(), models.LayarOutstanding(), false)
	var cari func(ts []models.Tata, jalur string) *models.Tata
	cari = func(ts []models.Tata, jalur string) *models.Tata {
		for i := range ts {
			if ts[i].Jenis == models.JenisMedan && ts[i].Jalur == jalur {
				return &ts[i]
			}
			if x := cari(ts[i].Anak, jalur); x != nil {
				return x
			}
		}
		return nil
	}
	for jalur, mau := range map[string]string{
		models.CD + "DateOfLoss":   models.KTanggal,
		models.CD + "DateReceived": models.KTanggal,
		models.CD + "ReporterTelp": models.KTelepon,
	} {
		if m := cari(ts, jalur); m == nil || m.Kendali != mau {
			t.Fatalf("%s: kendali %+v, mau %s", jalur, m, mau)
		}
	}
	if m := cari(ts, models.CD+"PeriodPolicyTBA"); m == nil || m.Label != "Policy Period TBA ?" {
		t.Fatalf("label PeriodPolicyTBA: %+v", m)
	}
	label := func(p string) []string {
		var out []string
		for _, k := range models.KodePilihan[p] {
			out = append(out, k+" "+models.LabelKode[p][k])
		}
		return out
	}
	if got := strings.Join(label("ReportType"), "|"); got != "1 Direct|2 Via Email|3 Via Fax|4 via Postal Mail/Courier|5 Via Telephone" {
		t.Fatalf("ReportType: %s", got)
	}
	if got := strings.Join(label("ReporterStatus"), "|"); got != "1 Ceding Co Name|2 SOB Name|3 Others" {
		t.Fatalf("ReporterStatus: %s", got)
	}
}

// Tombol "+" Consultant / Adjuster (MstAdjusterConsultant; work owner 08-10-2026): aksi layar tambah master, nonaktif
// selama ID hanya-baca (IsOutstanding = 1).
func TestTombolTambahAdjuster(t *testing.T) {
	for _, outstanding := range []string{"", "1"} {
		h := models.HalamanBaru()
		h.Setel("IsOutstanding", outstanding)
		ts := models.Evaluasi(h, models.LayarOutstanding(), false)
		for id, aksi := range map[string]string{"AdjusterConsultantBaru1": models.AksiTambahKonsultan,
			"AdjusterConsultantBaru2": models.AksiTambahAdjuster} {
			b := cariTata(ts, id)
			if b == nil || b.Aksi != aksi || b.Ikon != models.IkonTambah {
				t.Fatalf("%s: %+v", id, b)
			}
			if b.Nonaktif != (outstanding == "1") {
				t.Fatalf("%s IsOutstanding=%q: nonaktif=%v", id, outstanding, b.Nonaktif)
			}
		}
	}
}
