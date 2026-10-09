package models_test

import (
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimprop/backend/models"
)

func tgl(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

// AC 123: urutan TANGGAL naik lalu AcceptedNo turun (kosong lebih dulu, seperti DESC Oracle); baris berlaku = baris
// TERAKHIR urutan itu, sisanya riwayat.
func TestBarisBerlakuUrutanAC123(t *testing.T) {
	rows := []models.BarisOSLama{
		{CaseID: "UJI-C", AcceptedNo: "UJI-A1", Tanggal: tgl("2026-05-02 08:00:00"), StsReject: "2"},
		{CaseID: "UJI-C", AcceptedNo: "", Tanggal: tgl("2026-05-02 08:00:00"), StsReject: "0"},
		{CaseID: "UJI-C", AcceptedNo: "UJI-A2", Tanggal: tgl("2026-05-02 08:00:00"), StsReject: "4"},
		{CaseID: "UJI-C", AcceptedNo: "UJI-A9", Tanggal: tgl("2026-05-01 08:00:00"), StsReject: "0"},
	}
	b, riwayat := models.BarisBerlaku(rows)
	if b.AcceptedNo != "UJI-A1" || b.StsReject != "2" {
		t.Fatalf("baris berlaku %+v, mau AcceptedNo UJI-A1 (tanggal terakhir, AcceptedNo terkecil)", b)
	}
	if riwayat != 3 {
		t.Fatalf("riwayat %d, mau 3", riwayat)
	}
	if _, n := models.BarisBerlaku(rows[:1]); n != 0 {
		t.Fatalf("satu baris: riwayat %d, mau 0", n)
	}
}

func TestTahapLama(t *testing.T) {
	for _, c := range []struct {
		sts, tahap string
		tutup, ok  bool
	}{
		{models.StsOSEstimasi, models.TahapOutstanding, false, true},
		{models.StsOSAkseptasi, models.TahapAcceptation, false, true},
		{models.StsOSTutupBerkas, models.TahapAcceptation, true, true},
		{models.StsOSAkseptasiKomite, models.TahapAcceptation, false, true}, // keputusan work owner 08-10-2026
		{"3", "", false, false},
		{"", "", false, false},
	} {
		tahap, tutup, ok := models.TahapLama(c.sts)
		if tahap != c.tahap || tutup != c.tutup || ok != c.ok {
			t.Errorf("TahapLama(%q) = %q %v %v", c.sts, tahap, tutup, ok)
		}
	}
}

func sebabUntuk(d []models.MedanDibuang, jalur string) string {
	for _, m := range d {
		if m.Jalur == jalur {
			return m.Sebab
		}
	}
	return ""
}

const dokumenLama = `{
  "pxObjClass": "UJI-KELAS", "pzInsKey": "UJI-KUNCI",
  "NoClaim": "UJI-K01.05.2026.T00001", "IDMaster": "UJI-M1", "DateOfLoss": "20260501T013000.000 GMT",
  "ReportDate": "20260502", "CauseOfLoss": "UJI-SEBAB", "DaftarObjek": "UJI-OBJEK",
  "PolicyData": {"PolicyNo": "UJI-POLIS", "pxObjClass": "UJI"},
  "EstimationList": [
    {"pxCreateDateTime": "1", "EstimationDate": "20260503", "GrossEstimationPct": "1000.5", "CurrencyID": "UJI-IDR",
     "ClaimData": {"PolicyData": {"PolicyNo": "UJI-POLIS"}}},
    {"EstimationDate": "20260504", "GrossEstimationPct": "2000"}
  ],
  "SpreadingRisk": [{"TreatyName": "UJI-T"}],
  "InterestListDtl": [{"ObjectName": "UJI"}],
  "PaymentData": {"Nilai": "1"},
  "SuggestList": [{"CommentSuggest": "UJI-RIWAYAT"}],
  "Attachment": [{"IMAGEID": "UJI-GAMBAR"}],
  "ClaimComitee": [{"KomiteID": "UJI-K"}],
  "AdjustmentList": [
    {"pyExpanded": "true", "Type": "1", "GrossAdjustment": "500", "DataCommitteeTreaty": {"Remarks": "UJI-CATATAN"},
     "LossAllocation": [{"TreatyName": "UJI-T", "SharePercentage": "40"}],
     "ComiteeClaim": [{"IDKomite": "UJI-JABATAN"}]}
  ]
}`

func TestUraiKlaimLama(t *testing.T) {
	h, dibuang, err := models.UraiKlaimLama(dokumenLama)
	if err != nil {
		t.Fatal(err)
	}
	cd := models.CD
	sama(t, "NoClaim", h.Ambil(cd+"NoClaim"), "UJI-K01.05.2026.T00001")
	sama(t, "PolicyNo", h.Ambil(cd+"PolicyData.PolicyNo"), "UJI-POLIS")
	sama(t, "DateOfLoss apa adanya", h.Ambil(cd+"DateOfLoss"), "20260501T013000.000 GMT")
	est := h.AmbilDaftar(models.DaftarEstimasi)
	if len(est) != 2 || est[0]["GrossEstimationPct"] != "1000.5" {
		t.Fatalf("EstimationList %+v", est)
	}
	if _, ada := est[0]["pxCreateDateTime"]; ada {
		t.Fatalf("kunci px ikut ke baris: %+v", est[0])
	}
	adj := h.AmbilDaftar(models.DaftarAdjustment)
	if len(adj) != 1 || adj[0]["DataCommitteeTreaty.Remarks"] != "UJI-CATATAN" {
		t.Fatalf("AdjustmentList %+v", adj)
	}
	if _, ada := adj[0]["pyExpanded"]; ada {
		t.Fatalf("pyExpanded ikut disimpan (AC 11)")
	}
	if la := h.AmbilDaftar(models.JalurAdj(1, models.AnakLossAllocation)); len(la) != 1 || la[0]["SharePercentage"] != "40" {
		t.Fatalf("LossAllocation baris adjustment %+v", la)
	}
	for _, j := range []string{cd + "SpreadingRisk", cd + "InterestListDtl", cd + "PaymentData"} {
		if len(h.AmbilDaftar(j)) != 0 {
			t.Errorf("%s diimpor", j)
		}
	}
	for jalur, mau := range map[string]string{
		cd + "pxObjClass":                                      models.SebabInternalPega,
		cd + "PolicyData.pxObjClass":                           models.SebabInternalPega,
		cd + "AdjustmentList(1).pyExpanded":                    models.SebabInternalPega,
		cd + "SpreadingRisk(1).TreatyName":                     models.SebabTidakDiimpor,
		cd + "InterestListDtl(1).ObjectName":                   models.SebabTidakDiimpor,
		cd + "PaymentData.Nilai":                               models.SebabTidakDiimpor,
		cd + "Attachment(1).IMAGEID":                           models.SebabLampiran,
		cd + "ClaimComitee(1).KomiteID":                        models.SebabKomite,
		models.JalurAdj(1, "ComiteeClaim(1).IDKomite"):         models.SebabKomite,
		cd + "DaftarObjek":                                     models.SebabTanpaKolom,
		cd + "EstimationList(1).ClaimData.PolicyData.PolicyNo": models.SebabTanpaKolom,
	} {
		if s := sebabUntuk(dibuang, jalur); s != mau {
			t.Errorf("%s: sebab %q, mau %q", jalur, s, mau)
		}
	}
	for _, m := range dibuang {
		if m.Jalur == cd+"NoClaim" || m.Jalur == cd+"EstimationList(1).GrossEstimationPct" {
			t.Errorf("medan berkolom dilaporkan dibuang: %+v", m)
		}
	}
	if r := h.AmbilDaftar(models.DaftarRiwayat); len(r) != 1 || r[0]["CommentSuggest"] != "UJI-RIWAYAT" ||
		r[0][models.PropRiwayatBaru] != "1" {
		t.Fatalf("SuggestList lama mau dibawa sebagai riwayat baru: %+v", r)
	}
	if s := sebabUntuk(dibuang, cd+"SuggestList(1).CommentSuggest"); s != "" {
		t.Errorf("riwayat dilaporkan dibuang: %q", s)
	}
	if g := models.GalatNilaiKatalog(h); len(g) != 0 {
		t.Fatalf("dua format tanggal Pega dan angka teks mau sah (AC 12): %+v", g)
	}
	if _, _, err := models.UraiKlaimLama(`[1, 2]`); err == nil {
		t.Fatalf("dokumen bukan objek mau galat")
	}
}

func TestGalatNilaiKatalog(t *testing.T) {
	h := models.HalamanBaru()
	h.Setel(models.CD+"DateOfLoss", "UJI-BUKAN-TANGGAL")
	h.Setel(models.CD+"PolicyData.PolicyNo", strings.Repeat("P", 65))
	h.SetelDaftar(models.DaftarEstimasi, []models.Baris{{"GrossEstimationPct": "1.000,50"}})
	g := models.GalatNilaiKatalog(h)
	if len(g) != 3 {
		t.Fatalf("galat %+v, mau 3 (tanggal, panjang, angka)", g)
	}
	for _, m := range g {
		if m.Nilai != "" && strings.Contains(m.Sebab, m.Nilai) {
			t.Errorf("sebab memuat nilai: %+v", m)
		}
	}
}

// Kasus yang hanya punya baris OS: kunci header dipetakan, medan baris OS tetap di OS_AKSEPTASI_KLAIM.
func TestHalamanDariOSLama(t *testing.T) {
	b := models.BarisOSLama{CaseID: models.KunciPegaLama("UJI-CLMP-7"), NoClaim: "UJI-NO", NoPolis: "UJI-POLIS",
		MasterID: "UJI-M", DataJSON: `{"NoClaim":"UJI-NO-JSON","CauseOfLoss":"UJI-SEBAB","CauseOfLossID":"UJI-S1",
		"Value":"10","pzInsKey":"UJI"}`}
	h, dibuang, err := models.HalamanDariOSLama(b)
	if err != nil {
		t.Fatal(err)
	}
	sama(t, "NoClaim kolom datar", h.Ambil(models.CD+"NoClaim"), "UJI-NO")
	sama(t, "PolicyNo", h.Ambil(models.CD+"PolicyData.PolicyNo"), "UJI-POLIS")
	sama(t, "IDMaster", h.Ambil(models.CD+"IDMaster"), "UJI-M")
	sama(t, "CauseOfLoss", h.Ambil(models.CD+"CauseOfLoss"), "UJI-SEBAB")
	sama(t, "CauseOfLossID", h.Ambil(models.CD+"CauseOfLossID"), "UJI-S1")
	if s := sebabUntuk(dibuang, "Value"); s != models.SebabBarisOS {
		t.Errorf("Value: sebab %q", s)
	}
	if s := sebabUntuk(dibuang, "pzInsKey"); s != models.SebabInternalPega {
		t.Errorf("pzInsKey: sebab %q", s)
	}
	if _, _, err := models.HalamanDariOSLama(models.BarisOSLama{DataJSON: "{"}); err == nil {
		t.Fatalf("JSON rusak mau galat")
	}
}
