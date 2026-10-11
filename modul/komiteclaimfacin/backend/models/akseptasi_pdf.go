package models

// Untuk apa berkas ini: ISI DOKUMEN AKSEPTASI - `PrintPDFAccep_MultiAksep_KMT` S1-S22 (susunan halaman `TempAcceptedNo`
// + `NameInput`) dan syarat `pega:when` stream `AcceptanceNotePDF` (urut stream). Pemetaan per medan: scratchpad tahap 2
// PEMETAAN.md; ringkasan di docs/PARITAS.md.
//
// PERBAIKAN kelainan XML (OQ-CFI-03, maksud pasti):
//   - S10.3.6 / S10.3.9 `.PaymentType=="6" || .PaymentType=="4" && Local.SizeTempUang..` dibaca `(6 || 4) && ...`;
//   - `TemporaryUang` segar per cetak (Pega tidak pernah membersihkannya - terbawa antarcetak satu sesi), bertambah
//     selama iterasi seperti S10.3.1-S10.3.3;
//   - Property-Remove di dalam For Each atas daftar yang sama (S9.2, S10.3.11, S10.4) = saringan biasa;
//   - lini (IsFire / IsAneka / ...) dinilai atas halaman KLAIM (Pega menilai `pyWorkPage` = kasus komite, yang tanpa
//     OfferFacIn - stream bisa kosong);
//   - "Swift Code" hanya bila terisi (`SwiftCode != ""` Java = perbandingan referensi; preseden Komite Claim Prop).
//
// Dipertahankan VERBATIM (maksud tidak pasti, preseden Komite Claim Prop): baris "Location of Loss" berisi CauseOfLoss
// (S14 `ClaimData.Location := ClaimData.CauseOfLoss`); nominal + Payable To + tabel spreading hanya untuk item TERAKHIR
// yang lolos (stream `Subscript == CountObject`); tabel "BreakDown Spreading (QS)" membaca `ObjectItemList(i).
// SpreadingQuotaShare` yang tanpa penulis di tingkat item - tidak pernah tercetak. "Premium Paid On" = PaymentData
// (REST getPremiumPaidOn, OQ-CFI-14) tidak disimpan - kosong.
//
// `[penyimpangan sadar]` saringan S10.3.11 `.IsPrintAccept == ""` tidak dipakai: PDF dibuat SESUDAH keputusan tersimpan
// (S12 sudah menyetel IsPrintAccept 1) - adjustment dipilih menurut status 1 + nomor akseptasi.

import (
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/kontrak"
)

// Nama stream dokumen akseptasi per lini (S16-S22; langkah terakhir yang cocok yang menang).
const (
	StreamAkseptasiMBU    = "AcceptanceNotePDFMBU"
	StreamAkseptasiTravel = "AcceptanceNotePDFTravel"
	StreamAkseptasiPA     = "AcceptanceNotePDFPA"
)

// StreamAkseptasiLini = S16-S22: stream dokumen akseptasi lini klaim `v`; kosong = tanpa lini cocok (Pega: stream
// kosong, PDF gagal). Hanya `StreamAkseptasi` yang diekspor korpus (`Claim Fac In/AcceptanceNotePDF.xml`).
func StreamAkseptasiLini(v map[string]string) string {
	s := ""
	for _, c := range []struct {
		cocok  bool
		stream string
	}{{IsFire(v), StreamAkseptasi}, {IsAneka(v), StreamAkseptasi}, {IsGolfInsurance(v), StreamAkseptasi},
		{IsMBU(v), StreamAkseptasiMBU}, {IsMarineCargo(v), StreamAkseptasi}, {IsTravel(v), StreamAkseptasiTravel},
		{IsPA(v), StreamAkseptasiPA}} {
		if c.cocok {
			s = c.stream
		}
	}
	return s
}

// GroupPanelPA - `OfferFacIn.QuotationData.GroupPanel` lini PA (`pega:when ...GroupPanel='002'`).
const GroupPanelPA = "002"

// itemDokumen - satu item yang lolos saringan S9-S10: adjustment-nya (`pos` = nomor baris asal) dan AdjustmentGross
// masing-masing.
type itemDokumen struct {
	n     int
	b     map[string]string
	adj   []map[string]string
	pos   []int
	gross []*apd.Decimal
}

// SusunAcceptanceNote = PrintPDFAccep_MultiAksep_KMT S1-S15 + isi stream `AcceptanceNotePDF`: dokumen akseptasi nomor
// `noAksep` (Param.NoAkseptasi) atas objek kasus komite (`kl.Objek` = Param.idxObj). `saat` = NameInput.CARI24;
// `teknik` = `OperatorID.pyUserName` (Technical PIC); `namaReas` = REINSURANCETYPE (CallSpreadingView, TreatyName).
func SusunAcceptanceNote(kl kontrak.KlaimFacIn, noAksep string, saat time.Time, teknik string,
	namaReas map[string]string) (DokumenPDF, error) {
	v, o := kl.Nilai, kl.Objek
	if o < 1 {
		return DokumenPDF{}, fmt.Errorf("models: dokumen akseptasi tanpa objek")
	}
	var h hitung
	uang := map[string]bool{} // S10.3.1-S10.3.3 TemporaryUang (segar per cetak, bertambah selama iterasi)
	adjKum := apd.New(0, 0)   // Local.ADJ (tidak direset antaritem)
	var items []itemDokumen
	for i, it := range kl.Daftar[DaftarItem(o)] { // S9.2 / S10.3 / S10.4
		x := itemDokumen{n: i + 1, b: it}
		for a, b := range kl.Daftar[DaftarAdj(o, i+1)] {
			if b["AcceptedNo"] != noAksep { // S10.3.1 pre (IsPrintAccept: lihat kepala berkas)
				continue
			}
			uang[b["CurrencyID"]] = true
			if b["AcceptanceStatus"] != KeputusanSetuju { // S10.3.4-S10.3.9 pre; S10.3.11 membuang
				continue
			}
			nilai := h.dari("AdjustmentValue", b["AdjustmentValue"])
			switch b["PaymentType"] {
			case "3":
				nilai = h.dari("SalvageValue", b["SalvageValue"])
			case "4", "6":
				nilai = h.dari("AdjusterFeeValue", b["AdjusterFeeValue"])
			}
			if len(uang) == 1 {
				adjKum = h.tambah(adjKum, nilai)
			} else {
				adjKum = nilai
			}
			x.adj, x.pos, x.gross = append(x.adj, b), append(x.pos, a+1), append(x.gross, adjKum)
		}
		if len(x.adj) > 0 {
			items = append(items, x)
		}
	}
	if h.err != nil {
		return DokumenPDF{}, h.err
	}
	if len(items) == 0 {
		return DokumenPDF{}, fmt.Errorf("models: tidak ada adjustment berakseptasi %q untuk dokumen", noAksep)
	}
	akhir := items[len(items)-1]
	adjAkhir := akhir.adj[len(akhir.adj)-1] // S14 ObjectList(<LAST>).ObjectItemList(<LAST>).Adjustment(<LAST>)
	ob := baris(kl, DaftarObjek, o)
	travel, pa := v[OQ+"BusinessType"] == "Travel", v[OQ+"GroupPanel"] == GroupPanelPA

	d := DokumenPDF{Nomor: adjAkhir["AcceptedNo"], Penutup: "Jakarta, " + saat.In(Jakarta).Format("02 January 2006")}
	if !travel && !pa {
		d.Judul = append(d.Judul, "ACCEPTED CLAIM INSURANCE")
	}
	if travel {
		d.Judul = append(d.Judul, "ACCEPTED CLAIM SIMAS TRAVEL INSURANCE")
	}
	if pa {
		d.Judul = append(d.Judul, "ACCEPTED PERSONAL ACCIDENT CLAIM")
	}
	tambah := func(label string, nilai ...string) { d.Baris = append(d.Baris, BarisPDF{Label: label, Nilai: nilai}) }
	tambah("Policy No", v["OfferFacIn.PolicyData.PolicyNo"])
	tambah("Line of Business", v[OQ+"BusinessName"])
	tambah("Claim No / Claim ID", v["ClaimData.NoClaim"]+"/"+v["ClaimData.ClaimNo"])
	tambah("Name of Insured", v[OQ+"InsuredName"])
	if v[OQ+"QQName"] != "" {
		tambah("QQ Name", v[OQ+"QQName"])
	}
	tambah("SOB Name", v[OQ+"SobName"])
	var ceding []string // CedingCedantList item / adjustment pertama sesudah saring; kosong -> CedingCoName
	for _, c := range kl.Daftar[DaftarDiAdj(o, items[0].n, items[0].pos[0], AnakCedant)] {
		ceding = append(ceding, c["CedingCoName"])
	}
	if len(ceding) == 0 {
		tambah("Ceding Co Name", v[OQ+"CedingCoName"])
	} else {
		tambah("Ceding Co Name", strings.Join(ceding, " "))
	}
	tambah("Period of Policy", tanggalGaris(v["OfferFacIn.PolicyData.StartDateTime"])+"- "+
		tanggalGaris(v["OfferFacIn.PolicyData.EndDateTime"])) // S14 pyNote / pyCategory
	namaItem := ob["ObjectName"] + " "
	var tsi []string
	for _, it := range items {
		namaItem += " - " + it.b["ObjectItemName"]
		t, err := AngkaPega("TSINusare", it.b["TSINusare"])
		if err != nil {
			return DokumenPDF{}, err
		}
		tsi = append(tsi, strings.TrimSpace(it.b["ObjectItemName"]+" "+it.b["Currency"]+" "+t))
	}
	tambah("Object Item Name", namaItem) // ObjectName + 2 nbsp + per item "- " ObjectItemName
	tambah("TSI RNM", tsi...)
	tambah("Location of Loss", v["ClaimData.CauseOfLoss"]) // S14 VERBATIM: Location := CauseOfLoss
	if travel && !pa {
		cov := ""
		if c := kl.Daftar[anak(DaftarObjek, o, AnakCoverage)]; len(c) > 0 {
			cov = c[0]["CoverageNote"]
		}
		tambah("Policy Condition", cov)
	}
	if pa {
		tambah("Name of Interest", ob["ObjectName"])
	}
	tambah("Date of Loss", tanggalGaris(v["ClaimData.DateOfLoss"]))
	if pa || (!travel && !pa) {
		tambah("Cause of Loss", v["ClaimData.CauseOfLoss"])
	}
	if travel {
		d.Baris = append(d.Baris, BarisPDF{Tabel: &TabelPDF{Kepala: []string{"Jenis Jaminan",
			"Jumlah Kerugian Yang Diajukan ()", "Jumlah Maksimum Penggantian(USD)", "Jumlah Klaim Yang Dibayar ()"},
			Baris: [][]string{{"", "", "", ""}}}})
	}
	tambah("Premium Paid On", "") // PaymentData tidak disimpan (OQ-CFI-14)
	tambah("PIC Name / Technical PIC", adjAkhir["pxCreateOpName"]+"/"+teknik)
	spread := kl.Daftar[DaftarDiItem(o, akhir.n, AnakSpreadKlaim)]
	for x, a := range akhir.adj { // stream: blok per adjustment item terakhir
		g, err := AngkaPega("AdjustmentGross", akhir.gross[x].Text('f'))
		if err != nil {
			return DokumenPDF{}, err
		}
		label := "Accepted Claim"
		switch a["PaymentType"] {
		case "3":
			label = "Accepted Salvage"
		case "4":
			label = "Accepted Adjuster Fee"
		case "6":
			label = "Accepted Consultant Fee"
		}
		tambah(label, a["Currency"]+" "+g)
		bayar := []string{a["PayableTo"], "Name of Bank : " + a["NameOfBank"]}
		if strings.TrimSpace(a["SwiftCode"]) != "" {
			bayar = append(bayar, "Swift Code : "+a["SwiftCode"])
		}
		bayar = append(bayar, "Branch of Bank : "+a["BranchOfBank"], "Account No : "+a["NoAccount"])
		tambah("Payable To", bayar...)
		tabel := &TabelPDF{Kepala: []string{"Treaty Name", "Share (%)", "Claim Spread"}}
		for _, s := range spread { // S12.2 ClaimSpreaded := @divide(Local.ADJ * Share / 100, 1, 4)
			cs := h.bagiBulat(h.kali(adjKum, h.dari("SharePercentage", s["SharePercentage"])), apd.New(100, 0), 4)
			share, err := AngkaPega("SharePercentage", s["SharePercentage"])
			if err != nil {
				return DokumenPDF{}, err
			}
			if h.err != nil {
				return DokumenPDF{}, h.err
			}
			t, err := AngkaPega("ClaimSpreaded", cs.Text('f'))
			if err != nil {
				return DokumenPDF{}, err
			}
			tabel.Baris = append(tabel.Baris, []string{namaReas[s["TreatyType"]], share, t})
		}
		d.Baris = append(d.Baris, BarisPDF{Label: "Spreading Adjustment", Tabel: tabel})
	}
	return d, nil
}

// tanggalGaris - tanggal halaman ("2006-01-02..." / "yyyyMMdd..." / DateTime Pega "yyyyMMddTHHmmss.SSS GMT")
// sebagai "dd/MM/yyyy" menurut Asia/Jakarta (`@FormatDateTime(..., "dd/MM/yyyy", "Asia/Jakarta")`, `format="date"` -
// preseden Komite Claim Prop); lainnya apa adanya.
func tanggalGaris(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, " GMT") {
		if t, err := time.Parse("20060102T150405.000 MST", s); err == nil {
			return t.In(Jakarta).Format("02/01/2006")
		}
	}
	if len(s) >= 10 && s[4] == '-' {
		if t, err := time.Parse("2006-01-02", s[:10]); err == nil {
			return t.Format("02/01/2006")
		}
	}
	if len(s) >= 8 {
		if t, err := time.Parse("20060102", s[:8]); err == nil {
			return t.Format("02/01/2006")
		}
	}
	return s
}
