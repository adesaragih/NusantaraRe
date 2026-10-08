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

// NilaiRosterKomite = AddKomiteTreatyChild_ACT langkah 18-19: batas roster = ValueAdjustment; subjectivity = 0
// (hanya jenjang terbawah).
func NilaiRosterKomite(b Baris) string {
	if b["IsSubjectivity"] == "true" {
		return "0"
	}
	return b["ValueAdjustment"]
}
