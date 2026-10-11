package models

import "strconv"

// Untuk apa berkas ini: REJECT CLAIM (perintah work owner 10-10-2026) - section `ClaimRejectTreaty` (isian reject),
// `SureRejectClaimTreaty` (konfirmasi) dan activity `CheckRejectClaimTreaty2` (diekspor work owner 10-10-2026 di
// akar `Claim Prop/`). Kasus komite TT 3 satu tingkat (Claim Dept. Head) tanpa adjustment; Komite Claim Prop
// memutusnya lewat `KomitePost_Reject`.
//
// Pintu masuk tidak diekspor (tidak ada section / flow action / flow yang membuka ClaimRejectTreaty; local action
// konfirmasi `SureRejectClaimT` tidak diekspor). Perintah work owner 10-10-2026: tombol "Reject Claim" di posisi admin
// (Outstanding) selama Claim No belum ada; sesudah Claim No ada, tombol itu diganti "Close Claim" posisi teknik.

// Reject Claim - jalur, teks, kode.
const (
	// TransferTolak - `childPageKomite.TransferType := 3` (CheckRejectClaimTreaty2 7.3).
	TransferTolak = "3"
	// JalurTanggalTolak / JalurPICTolak / JalurTipeAnalisis - isian hanya-baca ClaimRejectTreaty (Date Reject, PIC Name,
	// Tipe Analisis); diisi pra-proses local action yang tidak diekspor.
	JalurTanggalTolak = "TempCommiteClaim.DateOfComitee"
	JalurPICTolak     = "TempCommiteClaim.Initial"
	JalurTipeAnalisis = "TempCommiteClaim.TypeComentAnalysis"
	// TeksMintaTolak - riwayat CheckRejectClaimTreaty2 7.6 (VERBATIM): + ID kasus komite.
	TeksMintaTolak = "Request reject claim "
	// PesanTolakSudahAkseptasi - CheckRejectClaimTreaty2 4 (VERBATIM).
	PesanTolakSudahAkseptasi = "Cannot be rejected because there is already an acceptance"
	// OQTipeAnalisis - prompt values `TypeComentAnalysis` dan pra-proses pengisinya tidak diekspor.
	OQTipeAnalisis = "OQ-CP-19: pilihan Tipe Analisis (TypeComentAnalysis) dan pra-proses pengisinya tidak diekspor"
)

// NomorKlaimKosong - Claim No belum ada (`.NoClaim` dan `.ClaimNo` kosong): posisi admin menawarkan Reject Claim; sesudah
// terisi, Close Claim.
func NomorKlaimKosong(h *Halaman) bool {
	return h.Ambil(CD+"NoClaim") == "" && h.Ambil(CD+"ClaimNo") == ""
}

// PesanHapusEstimasi - CheckRejectClaimTreaty2 5.1 (VERBATIM): `"Please Delete Estimationlist " + .pxListSubscript +
// " in Location " + local.IdxObj`; `local.IdxObj` tidak pernah diisi di activity ini (salinan Claim Fac In), jadi kosong.
func PesanHapusEstimasi(n int) string {
	return "Please Delete Estimationlist " + strconv.Itoa(n) + " in Location "
}

// PeriksaTolak = CheckRejectClaimTreaty2 4-6: sudah ada akseptasi; baris estimasi yang belum terkirim
// (`.PrintFaceClaim == ""`) harus dihapus. `Local.Error` DITIMPA setiap temuan - hanya pesan TERAKHIR yang tampil.
func PeriksaTolak(h *Halaman) bool {
	pesan := ""
	if h.Ambil("IsAnyAcceptation") == "1" { // 4
		pesan = PesanTolakSudahAkseptasi
	}
	for i, e := range h.AmbilDaftar(DaftarEstimasi) { // 5.1
		if e["PrintFaceClaim"] == "" {
			pesan = PesanHapusEstimasi(i + 1)
		}
	}
	if pesan != "" {
		h.TambahPesan("", pesan)
		return false
	}
	return true
}

// SiapkanTolak - pra-proses isian hanya-baca ClaimRejectTreaty (padanan pra-proses local action yang tidak diekspor,
// pola Claim Fac In `SetRejectClaim_pre`): Date Reject = hari ini, PIC Name = nama pelaku. Tipe Analisis kosong
// (OQ-CP-19). Tidak disimpan.
func SiapkanTolak(k *Konteks, h *Halaman) error {
	if h.Ambil(JalurTanggalTolak) == "" {
		h.Setel(JalurTanggalTolak, k.Hari())
	}
	if h.Ambil(JalurPICTolak) == "" && k.Acuan != nil {
		nama, err := k.Acuan.NamaPelaku(k.Ctxt(), k.Pelaku)
		if err != nil {
			return err
		}
		h.Setel(JalurPICTolak, nama)
	}
	return nil
}

// kartuMenunggu - kartu "already submitted" pengajuan yang masih menunggu komite (jenis menurut TRANSFER_TYPE-nya) +
// tombol Close; hanya-baca.
func kartuMenunggu(id string) Unsur {
	nomor := ro(medan(PropKomiteMenunggu, LabelNomorKomite, KTeks))
	return bagian("",
		tampilKlien(letak(LetakInfo, label(JudulMenungguTolak), label(TeksMenungguTolak), nomor),
			PropJenisKomiteMenunggu, TransferTolak, false),
		tampilKlien(letak(LetakInfo, label(JudulMenungguTutup), label(TeksMenungguTutup), nomor),
			PropJenisKomiteMenunggu, TransferTolak, true),
		tombol(id, "Close", "TutupModal"),
	)
}

// kartuSukses - kartu sukses pop-up Close Without Payment / Reject Claim: judul, penjelasan, nomor kasus komite
// (`Komite.CARI1`, hanya-baca).
func kartuSukses(judul, teks string) Unsur {
	return letak(LetakSukses, label(judul), label(teks), ro(medan(PropKomiteBaru, LabelNomorKomite, KTeks)))
}

// LayarTolakKlaim - Section `ClaimRejectTreaty`: tiga isian hanya-baca, Chronology, Remarks (wajib); Submit = local
// action `SureRejectClaimT` (konfirmasi), Cancel. Selama pengajuan Reject / Close Without Payment klaim ini masih
// menunggu komite, pop-up hanya menampilkan kartu "already submitted" (perintah work owner 10-10-2026).
func LayarTolakKlaim() []Unsur {
	return []Unsur{
		// pengajuan Reject / Close Without Payment klaim ini masih menunggu komite: kartu "already submitted" saja
		tampilKlien(kartuMenunggu("MenungguTolak"), PropKomiteMenunggu, "", true),
		tampilKlien(bagian("",
			ro(medan(JalurTanggalTolak, "Reject Date", KTanggal)),
			ro(medan(JalurPICTolak, "PIC Name", KTeks)),
			catatan(ro(medan(JalurTipeAnalisis, "Analysis Type", KTeks)), OQTipeAnalisis),
			medan(PropKronologiTutup, "Chronology", KArea),
			wajibU(medan("TempCommiteClaim.Remarks", "Remarks", KArea)),
			tombol("SubmitTolak", "Submit", "SureRejectClaimT"),
			tombol("CancelTolak", "Cancel", "TutupModal"),
		), PropKomiteMenunggu, "", false),
	}
}

// LayarYakinTolak - Section `SureRejectClaimTreaty` (kondisi `pyIsClientWhen`, SyaratKlien): selama `Komite.CARI1`
// kosong pertanyaan + Yes (`CheckRejectClaimTreaty2`, tanpa closeContainer) / No; sesudahnya pesan sukses + Close.
// Seperti pop-up Close Without Payment (perintah work owner 10-10-2026), nomor kasus komite ditampilkan.
func LayarYakinTolak() []Unsur {
	return []Unsur{
		tampilKlien(bagian("",
			label(TanyaTolak),
			tombol("RejectYes", "Yes", "CheckRejectClaimTreaty2"),
			tombol("RejectNo", "No", "TutupModal"),
		), PropKomiteBaru, "", false),
		tampilKlien(bagian("",
			kartuSukses(JudulSuksesTolak, TeksSuksesTolak),
			tombol("RejectSukses", "Close", "TutupModal"),
		), PropKomiteBaru, "", true),
	}
}
