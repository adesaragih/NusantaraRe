package models

// Untuk apa berkas ini: BATAS KOMITE di model Claim Non Prop - penyerahan akseptasi ke kasus komite KMTNP-
// (`CreateChildKomiteCNP_Act`), aturan tangga (OQ-CNP-01), dan grid "Committe Accept Status" (Section
// AdjustmentDetailNP / InputAcceptation).
//
// ⛔ Nama properti keputusan anggota komite HANYA ditulis di berkas ini (pola Claim Prop); berkas lain memakai
// konstantanya. Penjaga batas Claim Life `komite_statik_test.go` mengecualikan berkas ini (izin work owner 09-10-2026).
// Kasus komite diputuskan modul tahap 2 (`komiteclaimnonprop`); berkas ini hanya menulis tangga awal dan membacanya.

import (
	"strconv"

	"github.com/cockroachdb/apd/v3"
)

// Properti baris `.ComiteeClaim` (Section AdjustmentDetailNP S36 / InputAcceptation "Committe Accept Status").
const (
	PropKeputusanAnggota = "KomiteAproval"
	PropTanggalKeputusan = "DateApprove"
	PropCatatanKeputusan = "KomiteComment"
)

// ApprovalKomiteMenunggu - `KomiteAproval = 0` (CreateChildKomiteCNP_Act 26.8.1).
const ApprovalKomiteMenunggu = "0"

// STSKlaimNonProp - `Param.STS_KLAIM = "NONPROP"` roster EMAILKOMITE (CreateChildKomiteCNP_Act 26.3).
const STSKlaimNonProp = "NONPROP"

// Batas tangga komite (CreateChildKomiteCNP_Act langkah 10, 12; keputusan work owner 09-10-2026 OQ-CNP-01 "ikuti XML"):
// RNM Share <= 30 DAN ValueAdjustment akseptasi terakhir <= 30.000.000 -> hanya tingkat roster ber-LIMIT_BOTTOM <= 0
// (DEGREE 1). Batas Div Head 50.000.000 (langkah 13-14) ber-REMARK.
const (
	BatasNilaiKomite  = "30000000"
	BatasPersenKomite = "30"
)

// StatusKomiteDeptHead - `pyWorkPage.CNPStatusCase` sesudah penyerahan (CreateChildKomiteCNP_Act 31,
// CreateChildKomiteCloseNP_Act 12).
const StatusKomiteDeptHead = "COMITEE ACCEPTANCE (DEPT. HEAD)"

// PesanGrossSpreadingIn - CreateChildKomiteCNP_Act langkah 38 (VERBATIM).
const PesanGrossSpreadingIn = "Nilai Gross Value tidak sesuai dengan Spreading In"

// HanyaTingkat1 = CreateChildKomiteCNP_Act langkah 11-12 + 26.4 / 26.6: tangga hanya tingkat 1 bila `Flagkomite = 1`
// (RNM Share <= 30 dan ValueAdjustment akseptasi TERAKHIR <= 30.000.000) atau akseptasi subjectivity.
// ⚠️ `ValueAdjustment` tidak punya penulis di korpus Claim Non Prop (nilai kosong = 0) - OQ.
func HanyaTingkat1(h *Halaman, n int) bool {
	b, err := adj(h, n)
	if err != nil {
		return false
	}
	if b["IsSubjectivity"] == "true" {
		return true
	}
	var kal Kalkulator
	nilai := apd.New(0, 0)
	for _, a := range h.AmbilDaftar(DaftarAdjustment) { // 11
		nilai = kal.B(a, "ValueAdjustment")
	}
	persen := kal.H(h, TM+"RNMShare")
	if kal.Galat() != nil {
		return false
	}
	return !Lebih(persen, kal.D(BatasPersenKomite)) && !Lebih(nilai, kal.D(BatasNilaiKomite))
}

// SusunKomiteAkseptasi - grid "Committe Accept Status" akseptasi `n` + `.TotalKomite`. Akseptasi yang sudah diserahkan
// (`tangga` tidak nil, dibaca dari tangga kasus komitenya) menampilkan keputusan anggotanya; yang belum menampilkan
// roster calon menurut aturan tangga (`calon`).
// ⚠️ `[inferensi]` (PARITAS, OQ): TotalKomite hanya ditulis CreateChildKomiteCNP_Act 26.11 SESUDAH penyerahan, padahal
// tombol "Send to Committe" tampil hanya bila TotalKomite terisi - di XML tombol itu tak pernah tampil untuk akseptasi
// baru. Di sini TotalKomite = cacah tangga (calon atau tersimpan), pola Claim Prop (SetKomiteTreaty_ACT).
func SusunKomiteAkseptasi(h *Halaman, n int, tangga, calon []AnggotaKomite) error {
	b, err := adj(h, n)
	if err != nil {
		return err
	}
	anggota := tangga
	if tangga == nil {
		anggota = calon
	}
	var rows []Baris
	for _, a := range anggota {
		appr := a.Approval
		if appr == "" {
			appr = ApprovalKomiteMenunggu
		}
		rows = append(rows, Baris{"KomiteID": a.OperatorID, "KomitePost": a.Jabatan, "IDKomite": a.ID,
			PropKeputusanAnggota: appr, PropTanggalKeputusan: a.TanggalSetuju, PropCatatanKeputusan: a.Comment})
	}
	h.SetelDaftar(JalurAdj(n, AnakKomite), rows)
	if len(rows) == 0 {
		b["TotalKomite"] = ""
	} else {
		b["TotalKomite"] = strconv.Itoa(len(rows))
	}
	return nil
}

// gridKomite - grid "Committe Accept Status" (`.ComiteeClaim`, Section AdjustmentDetailNP S36).
func gridKomite(jalur string) Unsur {
	return bagian("Committe Accept Status", Unsur{Jenis: JenisGrid, Jalur: jalur, Bernomor: true,
		Kolom: []Unsur{kRO(kol("KomitePost", "Committee Name", KTampil)), kRO(kol("Initial", "", KTampil)),
			kSumber(kRO(kol(PropKeputusanAnggota, "Status", KPilih)), kode(PropKeputusanAnggota)),
			kRO(kol(PropTanggalKeputusan, "Date", KWaktu)), kRO(kol(PropCatatanKeputusan, "Comment", KTampil))}})
}

// BolehKirimKomite - akseptasi boleh diserahkan: belum di komite (`IsKomite != 1`), atau subjectivity diserahkan ulang.
func BolehKirimKomite(b Baris) bool {
	return b["IsKomite"] != "1" && (b[PropKomiteID] == "" || b["IsSubjectivity"] == "true")
}

// ValidasiGrossSpreadingIn = CreateChildKomiteCNP_Act langkah 21-22 (hanya bila Previously Calculated kosong): per mata
// uang Σ GrossValue layer XoL (= ClaimSpreaded baris non-UR) dibandingkan dengan Spreading In BERINDEKS SAMA
// (`SpreadingAdjustment(idx)`, mata uang sama); beda -> pesan langkah 38-39.
//
// PERBAIKAN OQ-CNP-05 butir 4 (PARITAS `[penyimpangan sadar]`): XML menyetel pesan ini SESUDAH kasus komite dibuat
// (langkah 29 pxAddChildWork mendahului 38-39 pada jalur transisi langkah 22.1). Di sini validasi jalan SEBELUM kasus
// dibuat dan aksi dibatalkan seluruhnya - tidak ada yang tersimpan.
func ValidasiGrossSpreadingIn(h *Halaman, n int) (bool, error) {
	if len(h.AmbilDaftar(JalurAdj(n, AnakXOLDibayar))) > 0 {
		return true, nil
	}
	var kal Kalkulator
	type gross struct {
		cur   string
		nilai *apd.Decimal
	}
	var temp []*gross
	for _, s := range h.AmbilDaftar(JalurAdj(n, AnakXOL)) {
		if s["TreatyName"] == TreatyUR {
			continue
		}
		var g *gross
		for _, t := range temp {
			if t.cur == s["Currency"] {
				g = t
			}
		}
		if g == nil {
			g = &gross{cur: s["Currency"], nilai: apd.New(0, 0)}
			temp = append(temp, g)
		}
		g.nilai = kal.Tambah(g.nilai, kal.B(s, "ClaimSpreaded"))
	}
	in := h.AmbilDaftar(JalurAdj(n, AnakSpreadIn))
	for i, g := range temp {
		if i >= len(in) || in[i]["Currency"] != g.cur {
			continue
		}
		if Banding(g.nilai, kal.B(in[i], "ClaimSpreaded")) != 0 {
			return false, kal.Galat()
		}
	}
	return true, kal.Galat()
}

// TandaiKirimKomite = CreateChildKomiteCNP_Act langkah 8 (IsKomite), 20 / 32 / 33 (baris terkunci CNPFlagOuts), 31
// (KomiteNo, CNPStatusCase, AcceptanceStatus 0).
func TandaiKirimKomite(h *Halaman, n int, komiteNo string) error {
	b, err := adj(h, n)
	if err != nil {
		return err
	}
	b["IsKomite"] = "1"
	for _, anak := range []string{AnakXOL, AnakClaimAccept, AnakLossAlloc} {
		for _, r := range h.AmbilDaftar(JalurAdj(n, anak)) {
			r["CNPFlagOuts"] = "1"
		}
	}
	b["KomiteNo"] = komiteNo
	b[PropKomiteID] = komiteNo
	b["AcceptanceStatus"] = "0"
	h.Setel("CNPStatusCase", StatusKomiteDeptHead)
	return nil
}
