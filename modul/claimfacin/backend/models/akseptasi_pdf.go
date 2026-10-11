package models

// Untuk apa berkas ini: ISI DOKUMEN AKSEPTASI - `PrintPDFAccep_MultiAksep` S1-S21 (SaveAcceptation 12, tombol
// "Acceptation"; susunan halaman `TempAcceptedNo` + `NameInput`) dan syarat `pega:when` stream `AcceptanceNotePDF`
// (urut stream). DISALIN dari Komite Claim Fac In `models/akseptasi_pdf.go` (`PrintPDFAccep_MultiAksep_KMT`, bukan
// impor); beda activity: sumber = halaman klaim di memori (pyWorkPage), periode dari `.Policy` (= salinan
// OfferFacIn.PolicyData), saringan IsPrintAccept VERBATIM (PDF disusun SEBELUM langkah 13 menyetelnya). Ringkasan di
// docs/PARITAS.md.
//
// PERBAIKAN kelainan XML (OQ-CFI-03, maksud pasti):
//   - S9.3.6 / S9.3.9 `.PaymentType=="6" || .PaymentType=="4" && Local.SizeTempUang..` dibaca `(6 || 4) && ...`;
//   - `TemporaryUang` segar per cetak (Pega tidak pernah membersihkannya - terbawa antarcetak satu sesi), bertambah
//     selama iterasi seperti S9.3.1-S9.3.3;
//   - Property-Remove di dalam For Each atas daftar yang sama (S8.2, S9.3.11, S9.4) = saringan biasa;
//   - "Swift Code" hanya bila terisi (`SwiftCode != ""` Java = perbandingan referensi; preseden Komite Claim Prop).
//
// Dipertahankan VERBATIM (maksud tidak pasti, preseden Komite Claim Prop): baris "Location of Loss" berisi CauseOfLoss
// (S13 `ClaimData.Location := ClaimData.CauseOfLoss`); nominal + Payable To + tabel spreading hanya untuk item TERAKHIR
// yang lolos (stream `Subscript == CountObject`); tabel "BreakDown Spreading (QS)" membaca `ObjectItemList(i).
// SpreadingQuotaShare` yang tanpa penulis di tingkat item - tidak pernah tercetak. "Premium Paid On" = PaymentData
// (REST getPremiumPaidOn, OQ-CFI-14) tidak disimpan - kosong.

import (
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"
)

// Kategori / ekstensi dokumen akseptasi (S14-S20 `param.AttachmentCategory`, S25 `MIME="pdf"`).
const (
	KategoriDokumenAkseptasi = "AcceptanceNote"
	MIMEDokumenAkseptasi     = "pdf"
)

// NamaBerkasAkseptasi = S14-S20 `"Persetujuan " + "  " + ClaimData.ClaimNo + " AcceptNo " + TempAcceptedNo.Policy.
// AccountNo + ".pdf"` (`ClaimData.ClaimNo` = ID kasus klaim, `Policy.AccountNo` = AcceptedNo).
func NamaBerkasAkseptasi(claimNo, acceptedNo string) string {
	return "Persetujuan " + "  " + claimNo + " AcceptNo " + acceptedNo + ".pdf"
}

// Nama stream dokumen akseptasi per lini (S14-S20; langkah terakhir yang cocok yang menang).
const (
	StreamAkseptasi       = "AcceptanceNotePDF"
	StreamAkseptasiMBU    = "AcceptanceNotePDFMBU"
	StreamAkseptasiTravel = "AcceptanceNotePDFTravel"
	StreamAkseptasiPA     = "AcceptanceNotePDFPA"
)

// StreamAkseptasiLini = S14-S20: stream dokumen akseptasi lini klaim `h`; kosong = tanpa lini cocok (Pega: stream
// kosong, PDF gagal). Hanya `StreamAkseptasi` yang diekspor korpus (`Claim Fac In/AcceptanceNotePDF.xml`).
func StreamAkseptasiLini(h *Halaman) string {
	s := ""
	for _, c := range []struct {
		cocok  bool
		stream string
	}{{IsFire(h), StreamAkseptasi}, {IsAneka(h), StreamAkseptasi}, {IsGolfInsurance(h), StreamAkseptasi},
		{IsMBU(h), StreamAkseptasiMBU}, {IsMarineCargo(h), StreamAkseptasi}, {IsTravel(h), StreamAkseptasiTravel},
		{IsPA(h), StreamAkseptasiPA}} {
		if c.cocok {
			s = c.stream
		}
	}
	return s
}

// InfoTanpaStream - pesan info Acceptation yang di Pega menerbitkan PDF dengan stream `stream` yang tidak diekspor
// korpus (OQ-CFI-20: hanya `AcceptanceNotePDF` yang diberikan work owner 10-10-2026); kosong = lini tanpa stream.
func InfoTanpaStream(stream string) string {
	if stream == "" {
		return "OQ-CFI-20: lini klaim tanpa stream dokumen akseptasi - berkas belum dibuat"
	}
	return "OQ-CFI-20: stream PDF " + stream + " tidak diekspor korpus - berkas belum dibuat"
}

// PesanDokumenGagal - akseptasi TERSIMPAN, tetapi PDF akseptasi gagal dibuat / diunggah / dicatat (teks sama dengan
// Komite Claim Prop / Komite Claim Fac In).
const PesanDokumenGagal = "The decision was saved, but the acceptance note PDF could not be stored. Please contact the administrator."

// GroupPanelPA - `OfferFacIn.QuotationData.GroupPanel` lini PA (`pega:when ...GroupPanel='002'`).
const GroupPanelPA = "002"

// itemDokumen - satu item yang lolos saringan S8-S9: adjustment-nya (`pos` = nomor baris asal) dan AdjustmentGross
// masing-masing.
type itemDokumen struct {
	n     int
	b     Baris
	adj   []Baris
	pos   []int
	gross []*apd.Decimal
}

// SusunAcceptanceNote = PrintPDFAccep_MultiAksep S1-S13 + isi stream `AcceptanceNotePDF`: dokumen akseptasi nomor
// `noAksep` (Param.NoAkseptasi = AcceptedNo adjustment yang diklik) atas objek `o` (Param.idxObj) halaman klaim `h`.
// `claimNo` = ClaimData.ClaimNo (ID kasus); `saat` = NameInput.CARI24; `teknik` = `OperatorID.pyUserName` (Technical
// PIC, penekan tombol); `namaReas` = REINSURANCETYPE (CallSpreadingView, TreatyName).
func SusunAcceptanceNote(h *Halaman, o int, noAksep, claimNo string, saat time.Time, teknik string,
	namaReas map[string]string) (DokumenPDF, error) {
	if o < 1 {
		return DokumenPDF{}, fmt.Errorf("models: dokumen akseptasi tanpa objek")
	}
	var kk Kalkulator
	uang := map[string]bool{} // S9.3.1-S9.3.3 TemporaryUang (segar per cetak, bertambah selama iterasi)
	adjKum := apd.New(0, 0)   // Local.ADJ (tidak direset antaritem)
	var items []itemDokumen
	for i, it := range h.AmbilDaftar(DaftarItem(o)) { // S8.2 / S9.3 / S9.4
		x := itemDokumen{n: i + 1, b: it}
		for a, b := range h.AmbilDaftar(DaftarAdj(o, i+1)) {
			if b["IsPrintAccept"] != "" || b["AcceptedNo"] != noAksep { // S9.3.1 pre
				continue
			}
			uang[b["CurrencyID"]] = true
			if b["AcceptanceStatus"] != "1" { // S9.3.4-S9.3.9 pre; S9.3.11 membuang
				continue
			}
			nilai := kk.B(b, "AdjustmentValue")
			switch b["PaymentType"] {
			case "3":
				nilai = kk.B(b, "SalvageValue")
			case "4", "6":
				nilai = kk.B(b, "AdjusterFeeValue")
			}
			if len(uang) == 1 {
				adjKum = kk.Tambah(adjKum, nilai)
			} else {
				adjKum = nilai
			}
			x.adj, x.pos, x.gross = append(x.adj, b), append(x.pos, a+1), append(x.gross, adjKum)
		}
		if len(x.adj) > 0 {
			items = append(items, x)
		}
	}
	if err := kk.Galat(); err != nil {
		return DokumenPDF{}, err
	}
	if len(items) == 0 {
		return DokumenPDF{}, fmt.Errorf("models: tidak ada adjustment berakseptasi %q untuk dokumen", noAksep)
	}
	akhir := items[len(items)-1]
	adjAkhir := akhir.adj[len(akhir.adj)-1] // S13 ObjectList(<LAST>).ObjectItemList(<LAST>).Adjustment(<LAST>)
	ob, err := Objek(h, o)
	if err != nil {
		return DokumenPDF{}, err
	}
	travel, pa := h.Ambil(OQ+"BusinessType") == "Travel", h.Ambil(OQ+"GroupPanel") == GroupPanelPA

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
	tambah("Policy No", h.Ambil(JalurNoPolis))
	tambah("Line of Business", h.Ambil(OQ+"BusinessName"))
	tambah("Claim No / Claim ID", h.Ambil(CD+"NoClaim")+"/"+claimNo)
	tambah("Name of Insured", h.Ambil(OQ+"InsuredName"))
	if h.Ambil(OQ+"QQName") != "" {
		tambah("QQ Name", h.Ambil(OQ+"QQName"))
	}
	tambah("SOB Name", h.Ambil(OQ+"SobName"))
	var ceding []string // CedingCedantList item / adjustment pertama sesudah saring; kosong -> CedingCoName
	for _, c := range h.AmbilDaftar(DaftarDiAdj(o, items[0].n, items[0].pos[0], AnakCedingCedant)) {
		ceding = append(ceding, c["CedingCoName"])
	}
	if len(ceding) == 0 {
		tambah("Ceding Co Name", h.Ambil(OQ+"CedingCoName"))
	} else {
		tambah("Ceding Co Name", strings.Join(ceding, " "))
	}
	tambah("Period of Policy", tanggalGaris(h.Ambil(JalurMulaiPolis))+"- "+
		tanggalGaris(h.Ambil(JalurAkhirPolis))) // S13 pyNote / pyCategory (.Policy.Start/EndDateTime)
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
	tambah("Location of Loss", h.Ambil(CD+"CauseOfLoss")) // S13 VERBATIM: Location := CauseOfLoss
	if travel && !pa {
		cov := ""
		if c := h.AmbilDaftar(JalurAnak(DaftarObjek, o, "CoverageList")); len(c) > 0 {
			cov = c[0]["CoverageNote"]
		}
		tambah("Policy Condition", cov)
	}
	if pa {
		tambah("Name of Interest", ob["ObjectName"])
	}
	tambah("Date of Loss", tanggalGaris(h.Ambil(CD+"DateOfLoss")))
	if pa || (!travel && !pa) {
		tambah("Cause of Loss", h.Ambil(CD+"CauseOfLoss"))
	}
	if travel {
		d.Baris = append(d.Baris, BarisPDF{Tabel: &TabelPDF{Kepala: []string{"Jenis Jaminan",
			"Jumlah Kerugian Yang Diajukan ()", "Jumlah Maksimum Penggantian(USD)", "Jumlah Klaim Yang Dibayar ()"},
			Baris: [][]string{{"", "", "", ""}}}})
	}
	tambah("Premium Paid On", "") // PaymentData tidak disimpan (OQ-CFI-14)
	tambah("PIC Name / Technical PIC", adjAkhir["pxCreateOpName"]+"/"+teknik)
	spread := h.AmbilDaftar(DaftarDiItem(o, akhir.n, AnakSpreadKlaim))
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
		for _, s := range spread { // S11.2 ClaimSpreaded := @divide(Local.ADJ * Share / 100, 1, 4)
			cs := bulat4(kk.Bagi(kk.Kali(adjKum, kk.B(s, "SharePercentage")), apd.New(100, 0)))
			if err := kk.Galat(); err != nil {
				return DokumenPDF{}, err
			}
			share, err := AngkaPega("SharePercentage", s["SharePercentage"])
			if err != nil {
				return DokumenPDF{}, err
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

// tanggalGaris - tanggal halaman (`UraiTanggal`: "2006-01-02", "yyyyMMdd", DateTime Pega GMT -> Asia/Jakarta) sebagai
// "dd/MM/yyyy" (`@FormatDateTime(..., "dd/MM/yyyy", "Asia/Jakarta")`, `format="date"`); lainnya apa adanya.
func tanggalGaris(s string) string {
	if t, ok := UraiTanggal(s); ok {
		return t.Format("02/01/2006")
	}
	return strings.TrimSpace(s)
}
