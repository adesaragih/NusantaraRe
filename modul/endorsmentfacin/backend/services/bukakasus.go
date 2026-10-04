package services

// Alur masuk endorsement - SEAM 5 (K-049), tiket E03.
//
// Untuk apa berkas ini: dari kasus portal (kelas `ASM-SFAGIS-Work-Endorsement`,
// tempat petugas memilih polis) lahir kasus endorsement (kelas
// `ASM-FW-GISFW-Work-Endorsement`, berawalan `EDM-`) - atau kasus itu ditolak
// dengan alasan bisnis.
//
// Dibaca sebelum: beforeimage.go - keluaran berkas ini masukan Seam 4.
//
// ⚠️ A06 - keputusan agent, menunggu konfirmasi: gerbang penolakan
// (gerbangtolak.go) dan validasi tanggal (validasitanggal.go) TIDAK dipanggil
// dari sini. `[terverifikasi]` Di Pega ketiganya dipanggil terpisah dari layar:
// `SetErrorBatalEndorsement_Act` dan `CheckEDMPolisDate` dari
// `Section/WorkPrimaryDetails.xml`, `SetValueToEDMWork` dari
// `Section/crmNewHarnessButtons.xml`; gerbang mematikan tombol lewat
// `InputData.CARIGROUP`. Pemanggil (layar/handler) wajib menjalankan keduanya
// lebih dulu dan tidak memanggil OpenCase bila salah satunya menolak.
//
// Asal (Pega): `Endorsment Fac In/Activity/SetValueToEDMWork.xml` langkah 1-13
// dan 22-25 · `DataTransform/DataToEDM.xml` · `Flow/InputAddendumFacultativeIn.xml`
// konektor `Start2 → Assignment7` dan shape `Assignment7`.
//
// ⛔ Penolakan adalah KELUARAN BISNIS (`HasilBukaKasus.Ditolak`), bukan galat.
// Galat (`error`) hanya untuk masukan yang tidak memungkinkan kasus dibuat sama
// sekali - kegagalan sistem.
//
// ⛔ `[terverifikasi]` Langkah 1-13 dan `DataToEDM` tidak memanggil satu pun
// rule `When`; Seam 5 karena itu tidak menunggu registry predikat EDM (E01).
//
// Tidak di sini: langkah 14-15 (Seam 4, beforeimage.go), 16 (selisih, E13),
// 17-20 (penyimpanan, repository), 21 (cabang EDM retro - `EDMRetro_Act`, fac
// out), 26 (di-remark `//`).

import (
	"errors"
	"time"

	"nusantarare/modul/endorsmentfacin/backend/models"
)

const (
	// KelasKasusEndorsement - `param.classname` langkah 7.
	KelasKasusEndorsement = "ASM-FW-GISFW-Work-Endorsement"
	// PrefiksIDEndorsement - `param.IDPrefix` langkah 7.
	PrefiksIDEndorsement = "EDM-"
	// FasePolicy - `.IsCedingConfirm = Policy` di konektor `Start2 → Assignment7`
	// (kata tanpa kutip dibaca sebagai teks, sikap registry NB butir 28).
	// Endorsement lahir LANGSUNG di fase polis: flow ini tidak punya konektor
	// yang menyetel fase penawaran atau binding.
	FasePolicy = "Policy"
	// pesanSudahBatal - `Local.Error` langkah 1 (teks yang sama dipakai
	// `SetErrorBatalEndorsement_Act` langkah 2 `Local.errmsg`).
	pesanSudahBatal = "Sudah Di endorsement Batal"
	// medanPolicyNo - `Field = .PolicyNo` semua `Property-Set-Messages` alur masuk.
	medanPolicyNo = "PolicyNo"
	// workbasketMarketing - `PositionNote` konektor Start2 dan workbasket `Assignment7`.
	workbasketMarketing = "ReasFacInMarketing"
	// nbStatusEDMBaru - `.NBStatus` / `.NBStatusNew` konektor Start2.
	nbStatusEDMBaru = "NEW EDM"
	// penandaSalah - `InputCari.CARI25` langkah 5-6.
	penandaSalah = "SALAH"
)

// ErrIdentitasKasusKosong - pembuat kasus tidak memberi identitas kasus baru.
var ErrIdentitasKasusKosong = errors.New("endorsement: identitas kasus baru (pxInsName/pzInsKey) kosong")

// QuotationPortal - `.Quotation` kasus portal.
type QuotationPortal struct {
	EdmType    models.JenisEndorsemen
	EdmTypeNew string
	Type       models.JenisPenyesuaian
}

// KasusPortal - kasus portal (`ASM-SFAGIS-Work-Endorsement`), sebatas yang
// dibaca alur masuk.
type KasusPortal struct {
	PzInsKey                 string
	PolicyNo                 string
	Note                     string
	EndorsementSource        string
	EndorsementSourceNote    string
	EndorsementDate          time.Time
	EndorsementStatus        string
	AdminCharge              string
	Survey                   string
	EndorsementInternalRetro string
	Quotation                QuotationPortal
}

// DaftarPolis - daftar nomor polis dari KONFIGURASI.
//
// ⛔ Nomor polis produksi tidak pernah literal di kode (CLAUDE.md §4.10, spec
// alur masuk §4.1): langkah 5 membypass dua nomor polis tertanam; di sini
// daftar itu masukan yang dapat diaudit, isinya tetap sama persis.
type DaftarPolis []string

// Memuat - apakah nomor polis ada di daftar.
func (d DaftarPolis) Memuat(nomor string) bool {
	for _, n := range d {
		if n == nomor {
			return true
		}
	}
	return false
}

// IdentitasKasus - identitas yang diberikan pembuat kasus (`svcAddWorkObject`;
// penomoran ID milik repository).
type IdentitasKasus struct {
	PxInsName, PzInsKey string
}

// MasukanBukaKasus - satu permintaan endorsement dari portal.
type MasukanBukaKasus struct {
	Portal KasusPortal
	// StatusEDMPolis - `OutputData1.pxResults(1).CARI20` hasil langkah 2
	// ("Get EdmStatus=2 / Batal Internal"). Kosong = query tanpa baris.
	StatusEDMPolis string
	// PengecualianBatal - konfigurasi pengganti dua literal langkah 5.
	PengecualianBatal DaftarPolis
	KasusBaru         IdentitasKasus
}

// PesanMedan - `Property-Set-Messages`: pesan pada satu medan kasus portal.
type PesanMedan struct {
	Medan, Teks string
}

// InfoAssignment - assignment pertama kasus (`Assignment7` "MARKETING").
type InfoAssignment struct {
	Workbasket, Tiket string
}

// HasilBukaKasus - keluaran Seam 5.
type HasilBukaKasus struct {
	// Ditolak - langkah 6 keluar activity: kasus TIDAK lahir.
	Ditolak     bool
	AlasanTolak string
	// Pesan - pesan yang dipasang pada kasus portal (langkah 4), juga ketika
	// kasus tetap lahir.
	Pesan []PesanMedan
	Kasus models.KasusEndorsement
	// PortalEDMHandle - `.EDMHandle` portal = `pzInsKey` kasus baru (`DataToEDM`).
	PortalEDMHandle string
	// PortalEdmStatus - `curWorkPage.Quotation.EdmStatus` (langkah 3).
	PortalEdmStatus string
	Assignment      InfoAssignment
}

// OpenCase - SEAM 5.
func OpenCase(m MasukanBukaKasus) (HasilBukaKasus, error) {
	p := m.Portal
	// 1
	cari17, cari25 := p.PolicyNo, ""
	// 3 (`pyStepsPreCondition=false`: tetap jalan, P-11)
	h := HasilBukaKasus{PortalEdmStatus: m.StatusEDMPolis}
	// 4 `CARI20==1` - literal TANPA kutip: pembanding angka, lewat `samaKode`
	// (teks terbaca angka sama tetapi teks berbeda → ErrTafsirKodeBerbeda,
	// butir 20). `WhenTrue` kosong = lanjut → jalan.
	batal, err := samaKode(m.StatusEDMPolis, models.StatusEDMPolisBatal)
	if err != nil {
		return HasilBukaKasus{}, err
	}
	if batal {
		h.Pesan = append(h.Pesan, PesanMedan{Medan: medanPolicyNo, Teks: pesanSudahBatal})
	}
	// 5 `CARI20=="1"` - literal BERKUTIP: pembanding teks; dan nomor polis
	// bukan pengecualian (konfigurasi).
	if m.StatusEDMPolis == models.StatusEDMPolisBatal && !m.PengecualianBatal.Memuat(cari17) {
		cari25 = penandaSalah
	}
	// 6 `CARI25=="SALAH"` → keluar activity (transisi 6).
	if cari25 == penandaSalah {
		h.Ditolak, h.AlasanTolak = true, pesanSudahBatal
		return h, nil
	}

	if m.KasusBaru.PxInsName == "" || m.KasusBaru.PzInsKey == "" {
		return HasilBukaKasus{}, ErrIdentitasKasusKosong
	}
	// 7-8 - kasus lahir.
	k := models.KasusEndorsement{
		Kelas: KelasKasusEndorsement, PrefiksID: PrefiksIDEndorsement,
		PxInsName: m.KasusBaru.PxInsName, PzInsKey: m.KasusBaru.PzInsKey,
		PolicyNumber: p.PolicyNo,
	}
	// 11 - `DataToEDM`, urutan korpus.
	// ⚠️ `[dugaan]` Baris-barisnya ada di dalam `UPDATE_PAGE newWorkPage`; di sini
	// sumber berawalan titik (`.PolicyNo`, `.Note`, `.Quotation.EdmType`, …)
	// dibaca dari kasus PORTAL. Konteks titik di dalam Update Page belum
	// terverifikasi - bila ia halaman yang di-update, sumber-sumber itu dibaca
	// dari kasus baru. Petunjuk ke arah sebaliknya: langkah 12-13 menyetel ulang
	// `OldPolicyNo` dan `EndorsmentReason` dari portal. `[pertanyaan terbuka]`.
	k.PyLabel = k.PxInsName
	k.EndorsementID = p.PzInsKey
	o := &k.OfferFacIn
	q := &o.QuotationData
	o.IsBanding = "false"
	q.OldPolicyNo = p.PolicyNo
	k.Policy.PolicyNo = p.PolicyNo
	q.EdmSource = p.EndorsementSource
	q.EdmSourceNote = p.EndorsementSourceNote
	q.EdmNote = p.Note
	q.EdmType = p.Quotation.EdmType
	q.EdmTypeNew = p.Quotation.EdmTypeNew
	q.Type = p.Quotation.Type
	q.EdmDate = p.EndorsementDate
	q.EdmStatus = p.EndorsementStatus
	q.EdmChargeFee = p.AdminCharge
	q.EdmSurvey = p.Survey
	q.EndorsementInternalRetro = p.EndorsementInternalRetro
	// Penanda siklus endorsement - yang dibaca `IsEDM` (K-029: 3 = EDM).
	q.StatusBusiness = models.StatusBusinessEDM
	o.EndorsmentReason = p.Note
	h.PortalEDMHandle = k.PzInsKey
	// 12-13
	q.OldPolicyNo = p.PolicyNo
	o.EndorsmentReason = p.Note

	// Flow `InputAddendumFacultativeIn`: konektor `Start2 → Assignment7`.
	// `Assignment7` memicu `SetNBStatus_Act`, yang `[terverifikasi]` langsung
	// keluar (langkah 2) untuk NBStatus berisi "NEW" di `ReasFacInMarketing` -
	// NBStatus tidak tertimpa.
	k.FlagOnGoingPolicy = "1"
	k.IsCedingConfirm = FasePolicy
	k.Position = "1"
	k.PositionNote = workbasketMarketing
	k.NBStatus = nbStatusEDMBaru
	k.NBStatusNew = nbStatusEDMBaru
	// Shape `Assignment7` "MARKETING": workbasket `ReasFacInMarketing`, tiket
	// `AdminPolicy`.
	h.Assignment = InfoAssignment{Workbasket: workbasketMarketing, Tiket: "AdminPolicy"}

	h.Kasus = k
	return h, nil
}

// Assignment - satu baris `Assign-WorkBasket` / `Assign-Worklist`.
type Assignment struct {
	RefObjectKey, InsKey string
}

// TautanAssignment - langkah 22-25, tautan KETIGA: `Primary.EDMHandle2` =
// `pzInsKey` assignment yang `pxRefObjectKey`-nya = `EDMHandle`.
//
// Langkah 24 mengulang hasil `Assign-WorkBasket`, langkah 25 hasil
// `Assign-Worklist`, keduanya menulis medan yang sama - jadi yang menang
// baris COCOK TERAKHIR, dengan worklist sesudah workbasket. Kosong bila tidak
// ada yang cocok.
func TautanAssignment(edmHandle string, workbasket, worklist []Assignment) string {
	hasil := ""
	for _, daftar := range [][]Assignment{workbasket, worklist} {
		for _, a := range daftar {
			if a.RefObjectKey == edmHandle {
				hasil = a.InsKey
			}
		}
	}
	return hasil
}
