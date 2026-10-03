package models

// Untuk apa berkas ini: KOMISI OGP DARI KONTRAK saat pilih bisnis - port
// `Activity/TreatyInputPctCommSpreading` (dipanggil
// `InputPolicyTreatyInDetail_preACT` langkah 17; tiket 01, 13; CATATAN-PUTARAN-2
// §2e). Sumber nilainya baris view relasional `POOLDATA.TREATYINDETAILJOINEDM`,
// dibaca repository (`KomisiKontrak`); fungsi di sini menerapkan langkah aslinya.
//
//	1        Call FetchMasterTreatyIn - JSON `M_TREATY_IN`/`M_TREATY_IN_EDM`
//	         `where ID = PolicyTreatyIn.NoOffer` (RDB BrowseTreatyInJoinEDM).
//	         ⭐ Pengganti relasional: baris view ber-`TREATYID = NoOffer`; kolom
//	         `TREATYTYPE`, `TREATYGROUP`, `RIOGR`, `RIONR` adalah medan
//	         `TreatyIn.Limits(n).TreatyType` / `.Detail(m).TreatyGroup/RIOGR/RIONR`
//	         yang didatarkan view (39 kolom, fakta katalog 03-10-2026).
//	2        kalang `pyWorkPage.TreatyIn.Limits`
//	2.1      syarat `pyWorkPage.PolicyTreatyIn.TreatyType==.TreatyType`
//	2.1.1    kalang `.Detail`
//	2.1.1.1  syarat `.TreatyGroup==pyWorkPage.PolicyTreatyIn.TreatyGroupName`
//	         RiCommOgp = @replaceAll(.RIOGR,",",".")
//	         RiCommOgp = @replaceAll(.RIONR,",",".")
//	         SpreadingRiskList(1).SharePercentage = .SpreadingTotalPct
//	         SpreadingRiskList(1).TreatyType      = .SpreadingTypeID
//	         SpreadingRiskList(1).TreatyName      = .SpreadingType
//	3        Call BreakDownSpreading_Act - ⛔ tidak (K9, KEPUTUSAN-RONDE-12 butir 3/3b)
//	4        berlabel `//` - dinonaktifkan di rule
//
// ⚠️ DITIRU APA ADANYA: RIOGR ditulis lalu langsung DITIMPA RIONR di langkah
// yang sama - hasil akhirnya RIONR (komisi ONR masuk ke medan OGP), dan
// RiCommOnp tidak disentuh. Bukan diperbaiki.
//
// ⛔ TIDAK DIBANGUN, alasan (c) bab 4 prompt putaran 2: tiga baris
// `SpreadingRiskList(1)` - `SpreadingTotalPct`, `SpreadingTypeID`,
// `SpreadingType` bukan kolom view (39 kolom: 33 kolom RD + BROKERAGE,
// COMMENCEMENT, TERMINATION, RIOGR, RIONR, RNM_SHARE); hanya ada di JSON master,
// dan pengecualian baca-JSON K8 terbatas pada medan jalur NonProp/XOL.
//
// ⚠️ Urutan baris: Pega menimpa di setiap baris yang syaratnya benar, menurut
// urutan Limits/Detail di dokumen JSON. Urutan itu tidak ada di view; baris
// diberikan repository berurut ID (penyimpangan kecil, sama dengan
// `BisnisDariKunci`). Bila hanya satu baris cocok, hasilnya sama persis.

import "strings"

// Kolom view yang dibaca langkah 2-2.1.1.1 - VERBATIM nama kolom view.
const (
	KolomJenisTreaty = "TREATYTYPE"
	KolomGrupTreaty  = "TREATYGROUP"
	KolomRIOGR       = "RIOGR"
	KolomRIONR       = "RIONR"
)

// LangkahKomisiProporsional = syarat `InputPolicyTreatyInDetail_preACT`
// langkah 17: `pyWorkPage.Quotation.ProportionalType=="NonProportional"` benar
// -> lewati; selain itu (termasuk kosong) aktivitas dijalankan.
func LangkahKomisiProporsional(h *Halaman) bool {
	return h.Ambil(HalamanQuotation+".ProportionalType") != JenisNonProporsional
}

// TreatyInputPctCommSpreading = langkah 2-2.1.1.1 atas `baris` (pengganti
// halaman master `TreatyIn.Limits(n).Detail(m)`, lihat kepala berkas).
func TreatyInputPctCommSpreading(h *Halaman, baris []BarisKontrak) {
	jenis := h.Ambil(HalamanPolis + ".TreatyType")
	grup := h.Ambil(HalamanPolis + ".TreatyGroupName")
	for _, b := range baris {
		// 2.1 / 2.1.1.1 - pembanding teks persis (`==` Pega atas teks)
		if b[KolomJenisTreaty] != jenis || b[KolomGrupTreaty] != grup {
			continue
		}
		h.Setel(HalamanPolis+".RiCommOgp", strings.ReplaceAll(b[KolomRIOGR], ",", "."))
		h.Setel(HalamanPolis+".RiCommOgp", strings.ReplaceAll(b[KolomRIONR], ",", "."))
	}
}
