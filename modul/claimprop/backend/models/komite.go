package models

// Untuk apa berkas ini: BATAS KOMITE di model Claim Prop - penyerahan baris adjustment ke kasus komite TKMT-
// (`AddKomiteTreatyChild_ACT`) dan grid "Committe Accept Status" (Section ComiteeClaimTreaty / AdjustmentDetail).
//
// ⛔ Nama properti keputusan anggota komite HANYA ditulis di berkas ini; berkas lain memakai konstantanya. Penjaga batas
// Claim Life `komite_statik_test.go` mengecualikan berkas ini (keputusan work owner 07-10-2026, opsi B: roster dan
// keputusan tetap milik konteks Komite, Claim Prop menulis tangga awal dan membacanya untuk tampilan).

import (
	"strconv"

	"github.com/cockroachdb/apd/v3"
)

// Properti baris `.ComiteeClaim` yang membawa keputusan anggota (SetKomiteTreaty_ACT 5.1, Section ComiteeClaimTreaty).
const (
	PropKeputusanAnggota = "KomiteAproval"
	PropTanggalKeputusan = "DateApprove"
	PropCatatanKeputusan = "KomiteComment"
)

// ApprovalKomiteMenunggu - `KomiteAproval = 0`: roster calon dan anggota tangga yang belum memutuskan
// (AddKomiteTreatyChild_ACT 22.1).
const ApprovalKomiteMenunggu = "0"

// TeksKirimKomite - riwayat AddKomiteTreatyChild_ACT langkah 31 (VERBATIM).
const TeksKirimKomite = "Send to Commite"

// SusunKomiteAdjustment - grid "Committe Accept Status" baris adjustment `idx`. Baris yang sudah diserahkan (`tangga`
// tidak nil, dibaca dari tangga kasus komitenya) menampilkan keputusan anggotanya; yang belum menampilkan roster calon
// (SetKomiteTreaty_ACT 5.1: `KomiteAproval = 0`, `IDKomite = .JABATAN`).
func SusunKomiteAdjustment(k *Konteks, h *Halaman, idx int, tangga []AnggotaKomite) error {
	b, err := adj(h, idx)
	if err != nil {
		return err
	}
	anggota := tangga
	if tangga == nil {
		if b["ProposeAdjustmentValue"] == "" || b["ValueAdjustment"] == "" {
			h.SetelDaftar(JalurAdj(idx, "ComiteeClaim"), nil)
			b["TotalKomite"] = ""
			return nil
		}
		if anggota, err = k.Acuan.RosterKomite(k.Ctxt(), b["ValueAdjustment"], STSKlaimProp); err != nil {
			return err
		}
	}
	var rows []Baris
	for _, a := range anggota {
		appr := a.Approval
		if appr == "" {
			appr = ApprovalKomiteMenunggu
		}
		rows = append(rows, Baris{"KomiteID": a.OperatorID, "IDKomite": a.Jabatan, PropKeputusanAnggota: appr,
			PropTanggalKeputusan: a.TanggalSetuju, PropCatatanKeputusan: a.Comment})
	}
	h.SetelDaftar(JalurAdj(idx, "ComiteeClaim"), rows)
	b["TotalKomite"] = strconv.Itoa(len(rows))
	return nil
}

// IsiEstimasiRetro - AddKomiteTreatyChild_ACT langkah 7-8: `Local.TotalEstimasireas` = spread klaim baris
// SpreadingClaim TERAKHIR (7.2 prakondisi nonaktif); FacRetroList `.TotalEstimasiReas = total x share / 100`.
func IsiEstimasiRetro(h *Halaman) error {
	var kal Kalkulator
	tot := apd.New(0, 0)
	for _, s := range h.AmbilDaftar(DaftarSpreading) {
		tot = kal.B(s, "ClaimSpreaded")
	}
	for _, f := range h.AmbilDaftar(DaftarFacRetro) {
		f["TotalEstimasiReas"] = Teks(kal.Persen(tot, kal.B(f, "PctShareAllObj")))
	}
	return kal.Galat()
}

// TandaiKirimKomite = AddKomiteTreatyChild_ACT langkah 1, 3-4, 31: baris `idx` terkunci komite, riwayat.
// `FlagOnGoingCommitte` (langkah 3) TIDAK disimpan - nama kolomnya dilarang MODUL.md; KomiteNo baris subjectivity
// dikosongkan.
func TandaiKirimKomite(k *Konteks, h *Halaman, idx int) error {
	b, err := adj(h, idx)
	if err != nil {
		return err
	}
	b["IsKomite"] = "1"
	if b["IsSubjectivity"] == "true" {
		b["KomiteNo"] = ""
	}
	k.Riwayat(h, TeksKirimKomite)
	return IsiEstimasiRetro(h)
}

// BolehSerahKomite - baris adjustment boleh diserahkan ke komite: belum pernah (KomiteID kosong), atau diserahkan ULANG
// karena disetujui bersyarat (keputusan work owner 08-10-2026, OQ-KCP-06 "a": Komite S23 `IsKomite := 0` membuka lagi
// tombol "Send to Committe", AddKomiteTreatyChild_ACT S16-S19 / S31 menyerahkannya ke jenjang terbawah). `IsKomite = 1`
// = masih di komite.
func BolehSerahKomite(b Baris) bool {
	return b["IsKomite"] != "1" && (b[PropKomiteID] == "" || b["IsSubjectivity"] == "true")
}

// ---------------------------------------------------------------- Close Without Payment (SendCloseClaimToKomite)

// Close Without Payment - perintah work owner 10-10-2026 (OQ-CP-06 ditutup lewat kolom bersama migrasi
// komiteclaimfacin 641-643).
const (
	// TransferTutup - `childPageKomite.TransferType := 4` (SendCloseClaimToKomite 5.3).
	TransferTutup = "4"
	// PropKronologiTutup - isian "Chronology" pop-up PreventRejectClaimProp (`TempCommiteClaim.CircumtansesCouseOfLoss`,
	// 5.2 -> `childPageKomite.Komite.CircumtansesCouseOfLoss`).
	PropKronologiTutup = "TempCommiteClaim.CircumtansesCouseOfLoss"
	// PropCentangTutupTanpaBayar - kotak centang "Close Without Payment" pop-up PreventRejectClaimProp.
	PropCentangTutupTanpaBayar = "TempCommiteClaim.AllocationShareSalvage"
	// JabatanTutupTanpaBayar - satu-satunya penyetuju kasus komite close (5.2-5.3 `IDKomite := "Claim Dept. Head"`).
	// XML menanam akun orang; di sini baris roster PROP berjabatan itu (workbasket ReasClaimDeptHead, migrasi 537).
	JabatanTutupTanpaBayar = "Claim Dept. Head"
	// BatasRosterSemua - batas roster yang memuat setiap baris aktif (LIMIT_BOTTOM <= batas): penyetuju close dicari
	// menurut JABATAN, bukan menurut nilai (langkah 5.2-5.3 tanpa batas nilai).
	BatasRosterSemua = "999999999999999"
	// TeksMintaTutupTanpaBayar - riwayat SendCloseClaimToKomite 5.6 (VERBATIM): + ID kasus komite.
	TeksMintaTutupTanpaBayar = "Request close claim without payment "
	// PesanTutupTanpaBayarBerjalan - `[tidak ada di XML]` penjaga kiriman ganda: Pega membiarkan Close Without Payment
	// dikirim lagi selama kasus komite close sebelumnya masih menunggu (penyimpangan sadar, PARITAS).
	PesanTutupTanpaBayarBerjalan = "Close without payment for this claim is already waiting for committee decision"
)

// PeriksaTutupTanpaBayar = SendCloseClaimToKomite langkah 3-4: setiap baris adjustment ber-AcceptanceStatus kosong
// -> "Can not close claim, there is adjustment in comitee!" (3.1 hanya `.AcceptanceStatus==""`, berbeda dengan
// CloseClaimProp yang juga menolak "0").
func PeriksaTutupTanpaBayar(h *Halaman) bool {
	for _, b := range h.AmbilDaftar(DaftarAdjustment) {
		if b["AcceptanceStatus"] == "" {
			h.TambahPesan("", PesanKomiteMasihJalan)
			return false
		}
	}
	return true
}

// TandaiTutupTanpaBayar = SendCloseClaimToKomite langkah 2: Remark / Remark_Close <- Remarks pop-up.
func TandaiTutupTanpaBayar(h *Halaman, remarks string) {
	h.Setel(CD+"Remark", remarks)
	h.Setel(CD+"Remark_Close", remarks)
}

// PilihPenyetujuTutup - baris roster PROP berjabatan `JabatanTutupTanpaBayar` (langkah 5.2-5.3); false bila tidak ada.
func PilihPenyetujuTutup(roster []AnggotaKomite) (AnggotaKomite, bool) {
	for _, a := range roster {
		if a.Jabatan == JabatanTutupTanpaBayar {
			return a, true
		}
	}
	return AnggotaKomite{}, false
}

// NilaiRosterKomite = AddKomiteTreatyChild_ACT langkah 18-19: batas roster = ValueAdjustment; subjectivity = 0
// (hanya jenjang terbawah).
func NilaiRosterKomite(b Baris) string {
	if b["IsSubjectivity"] == "true" {
		return "0"
	}
	return b["ValueAdjustment"]
}
