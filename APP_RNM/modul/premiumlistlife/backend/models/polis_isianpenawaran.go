package models

// Isian layar Input Offer - tiket 01 bagian 3 (form penawaran). MURNI.
//
// Untuk apa berkas ini: aturan atas isian penawaran yang pemakai ketik di
// layar `Section/InputOfferLife.xml` - pilihan tertutupnya, turunannya, dan
// baris riwayat `AddHistorySuggest` yang lahir setiap kali ia disimpan.
//
// `[terverifikasi]` dibaca sebagai pohon 30-09-2026 dari salinan korpus
// `kelvin\PremiumListLife (Done)\` (pengurai sel section + pengurai langkah
// activity, keduanya di scratchpad sesi; bukan berkas repo):
//
//	Section/InputOfferLife.xml   sel yang DAPAT DIISI (pyReadOnly false):
//	  .TypeCeding   pxDropdown  "System Reinsurance"   pyRequired true
//	  .BusinessCode pxDropdown  "Class of Business"    pyRequired true
//	  .DateReceived pxDateTime  "Email Received Date"
//	  pyWorkPage.Status      pxRadioButtons "Status"  tampil FlagOnGoingPolicy = 0
//	  pyWorkPage.EmailTypePL pxRadioButtons "Status"  tampil FlagOnGoingPolicy = 1
//	  .Description  pxTextArea  "Comment"              pyRequired true
//	  .CedingCoName / .PolicyHolderName - read-only, diisi popup
//	    `Ceding_Harness` / `PolicyHolder_Harness` (setCeding_act, setPolicyHolder_act)
//	Activity/InputOfferLife_ACT  3 CARI8 (nama TypeCeding), 7 (JenisAsuransi)
//	DataTransform/SetCoBName_Act     BusinessCode -> BusinessName (21 cabang + UNKNOWN)
//	DataTransform/SetReinsuranceType TypeCeding "4" -> Non Proportional
//	DataTransform/SetStatusAkseptasi Status/EmailTypePL -> ProposalAcceptStatus
//	Activity/AddHistorySuggest   1 (Offer, Flag 0) / 2 (Bind, Flag 1)
//
// ⛔ RALAT 01-10-2026 - sel ber-`pyReadOnly true` BUKAN selalu read-only.
// Ronde pertama berkas ini menyimpulkan Age Limit, tanggal-tanggal, TBC,
// Status Update, Marketing Note tidak dapat diisi. Layar Pega lama yang
// ditunjukkan work owner membantahnya: sel-sel itu berisi ketikan. Read-only
// hanya berlaku lewat `pyReadOnlyCondition` (`pyWorkPage.FlagOnGoingPolicy='1'`,
// kasus Input Premium). Yang kini dibawa - sel yang TERISI di layar itu
// (migrasi 059, kecuali SUM_INSURED/STATUS_UPDATE yang sudah ada):
//
//	.BatasUsiaPeserta "Age Limit" · .PeriodePertanggungan "Coverage Period"
//	.SumInsured "Sum Insured" · .TanggalPenawaran "Offering Date"
//	.TanggalRespon "Response Date" · .TanggalKonfirmasi "Confirmation Date"
//	.TBC "Input TBC" · .TanggalTBC "Max TBC" (SetMaxTBCLife_Act, dihitung)
//	.StatusUpdate "Status Update" · .KeteranganMarketing "Marketing Note"
//
// Sel yang KOSONG di layar itu (Insured Name, Occupation, Underwriting Policy,
// Re-Confirmation/Realization/Binding Date, Final Status) ikut dibawa lewat
// migrasi 060 - keputusan work owner 01-10-2026: "tampilkan semua kolom di
// gambar, nanti baru saya filter".
//
// Dibaca sesudah: polis_validasi.go (gerbang ProtectAccept yang memakai isian ini).

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"
)

// Pilihan adalah satu opsi dropdown atau radio - kode yang disimpan dan teks
// yang tampil.
type Pilihan struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

// PilihanTypeCeding - dropdown `System Reinsurance`.
//
// `[terverifikasi]` pasangan kode-nama dari `InputOfferLife_ACT` langkah 3
// `TempInputDataLife.CARI8` (`@if(TypeCeding="1","QS",...)`), domain yang
// sama dengan keputusan tiket 00 Endorsement (STRUKTUR `TYPE_CEDING`).
var PilihanTypeCeding = []Pilihan{
	{Kode: "1", Nama: "QS"},
	{Kode: "2", Nama: "SURPLUS"},
	{Kode: "3", Nama: "QS + SURPLUS"},
	{Kode: "4", Nama: "XOL"},
}

// PilihanClassOfBusiness - dropdown `Class of Business`.
//
// `[terverifikasi]` `DataTransform/SetCoBName_Act.xml` langkah 1-21, VERBATIM
// termasuk ejaannya (`INDIVIDU`, `ILNESS`): nama inilah yang tersimpan di
// data warisan, dan pencarian atas nama yang "dirapikan" tidak menemukannya.
var PilihanClassOfBusiness = []Pilihan{
	{Kode: "L1", Nama: "INDIVIDUAL TERM LIFE"},
	{Kode: "L2", Nama: "GROUP TERM LIFE"},
	{Kode: "L3", Nama: "INDIVIDUAL WHOLE LIFE"},
	{Kode: "L4", Nama: "GROUP PA"},
	{Kode: "L5", Nama: "GROUP LEVEL TERM LIFE"},
	{Kode: "L6", Nama: "GROUP DECREASING TERM LIFE"},
	{Kode: "L7", Nama: "INDIVIDU ENDOWMENT LIFE"},
	{Kode: "L8", Nama: "INDIVIDU INCREASING TERM LIFE"},
	{Kode: "L9", Nama: "GROUP ENDOWMENT LIFE"},
	{Kode: "L10", Nama: "GROUP INCREASING TERM LIFE"},
	{Kode: "L11", Nama: "INDIVIDU PA"},
	{Kode: "L12", Nama: "INDIVIDU EXPENSE HEALTH"},
	{Kode: "L13", Nama: "INDIVIDU DISABILITY HEALTH"},
	{Kode: "L14", Nama: "GROUP EXPENSE HEALTH"},
	{Kode: "L15", Nama: "GROUP DISABILITY HEALTH"},
	{Kode: "L16", Nama: "INDIVIDU CRITICAL ILNESS"},
	{Kode: "L17", Nama: "INDIVIDU TPD"},
	{Kode: "L18", Nama: "INDIVIDU HOSPITAL CASH PLAN"},
	{Kode: "L19", Nama: "GROUP CRITICAL ILLNESS"},
	{Kode: "L20", Nama: "GROUP TPD"},
	{Kode: "L21", Nama: "GROUP TERMINAL ILLNESS"},
}

// NamaBusinessTakDikenal - cabang `OTHERWISE` SetCoBName_Act.
const NamaBusinessTakDikenal = "UNKNOWN"

// PilihanStatusOffer - radio `Status` (`pyWorkPage.Status`), tampil saat
// FlagOnGoingPolicy "0".
//
// ⚠️ `[dugaan]` DAFTARNYA DIRAKIT, bukan dibaca utuh: pilihan radio itu
// `ListSource=associated` - milik rule properti `Status`, yang TIDAK ada di
// ekspor. Kelima nilai ini SELURUH nilai yang korpus bandingkan dengan
// `.Status` di jalur penawaran: syarat tampil tombol `Save Offer`
// (`Accept`/`Pending`) dan tombol ber-label `.Status`
// (`Decline`/`Bind`/`Closed`) di InputOfferLife.xml, `SetStatusAkseptasi`
// (`Bind`, `Decline`, `Closed`), dan `InputOfferLife_ACT` langkah 1 (`//`,
// `Decline`/`Closed`). Urutannya urutan kemunculan itu.
var PilihanStatusOffer = []Pilihan{
	{Kode: "Accept", Nama: "Accept"},
	{Kode: "Pending", Nama: "Pending"},
	{Kode: "Bind", Nama: "Bind"},
	{Kode: "Decline", Nama: "Decline"},
	{Kode: "Closed", Nama: "Closed"},
}

// PilihanEmailTypePL - radio `Status` (`pyWorkPage.EmailTypePL`), tampil saat
// FlagOnGoingPolicy "1".
//
// `[terverifikasi]` `AddHistorySuggest` langkah 2: `EmailTypePL=1` Accept,
// `=2` Reject, selain itu Decline. ⚠️ `[dugaan]` kode "3" untuk Decline:
// rule itu tidak menyebut kodenya (cabang `else`), dan pilihan radionya milik
// rule properti yang tidak diekspor.
var PilihanEmailTypePL = []Pilihan{
	{Kode: "1", Nama: "Accept"},
	{Kode: "2", Nama: "Reject"},
	{Kode: "3", Nama: "Decline"},
}

// Nilai `Initial` riwayat - `AddHistorySuggest` langkah 1/2.
const (
	InitialSuggestOffer = "Offer"
	InitialSuggestBind  = "Bind"
)

// Pesan wajib-isi layar - kalimat PEGA tidak ada untuk `pyRequired` (Pega
// menampilkan pesan bawaan platformnya), jadi yang dipakai kalimat gerbang
// `ProtectAccept` bila ada, dan kalimat berpola sama untuk sisanya.
//
// ⚠️ Ceding Name, Policy Holder, dan Status WAJIB sejak 01-10-2026 -
// `[keputusan work owner]`, BUKAN dari korpus: di InputOfferLife.xml ketiganya
// tidak ber-pyRequired. Kalimatnya pun karangan layar baru, berpola
// `<label> can't null` milik ProtectAccept.
const (
	PesanCommentKosong      = "Comment can't null"
	PesanCedingKosong       = "Ceding Name can't null"
	PesanPolicyHolderKosong = "Policy Holder can't null"
	PesanStatusKosong       = "Status can't null"
)

var (
	// ErrIsianPenawaranTidakSah - isian di luar pilihan tertutupnya.
	ErrIsianPenawaranTidakSah = errors.New("models: isian penawaran tidak sah")
	// ErrIsianPenawaranBelumLengkap - isian wajib yang kosong.
	ErrIsianPenawaranBelumLengkap = errors.New("models: isian penawaran belum lengkap")
)

// IsianPenawaran adalah yang pemakai kirim dari layar Input Offer.
//
// ⚠️ Seluruhnya teks kecuali tanggal - pilihan dikirim sebagai KODE, dan nama
// turunannya dihitung di sini, tidak pernah dipercaya dari klien.
type IsianPenawaran struct {
	CedingCo         string
	CedingCoName     string
	PolicyHolder     string
	PolicyHolderName string
	TypeCeding       string
	BusinessCode     string
	DateReceived     *time.Time
	Description      string
	// Status - kode radio: `PilihanStatusOffer` (bendera "0") atau
	// `PilihanEmailTypePL` (bendera "1"). Boleh kosong: radionya tidak
	// ber-pyRequired.
	Status string

	// Sel penawaran yang terisi di layar lama (migrasi 059 + kolom 051).
	BatasUsiaPeserta     *int
	PeriodePertanggungan string
	// SumInsured - uang: desimal presisi arbitrer, nol `float` (ADR-0003).
	SumInsured          *apd.Decimal
	TanggalPenawaran    *time.Time
	TanggalRespon       *time.Time
	TanggalKonfirmasi   *time.Time
	TBC                 *int
	StatusUpdate        string
	KeteranganMarketing string

	// Sel sisa layar itu (migrasi 060) - kosong di layar lama yang ditunjukkan.
	QQName                 string
	JenisUsaha             string
	KetentuanUnderwriting  string
	TanggalKonfirmasiBalik *time.Time
	TanggalRealisasi       *time.Time
	TanggalBind            *time.Time
	StatusFinal            string
}

// PenawaranTersimpan adalah isian beserta turunannya, siap ditulis.
type PenawaranTersimpan struct {
	IsianPenawaran
	TypeCedingName string
	BusinessName   string
	// JenisAsuransi - "Reinsurance Type", turunan TypeCeding (`JenisAsuransi`).
	JenisAsuransi string
	// TanggalTBC - "Max TBC", lihat MaxTBC.
	TanggalTBC *time.Time
}

// MaxTBC meniru `SetMaxTBCLife_Act` langkah 2:
// `@DateTime.addCalendar(TanggalKonfirmasi, 0,0,0, TBC, 0,0,0)` - Confirmation
// Date ditambah TBC HARI.
//
// ⛔ nil bila salah satunya kosong. Pega menghitung dari tanggal kosong juga
// (hasilnya tidak bermakna); sel `Max TBC` sendiri hanya tampil bila `.TBC`
// terisi, dan `Input TBC` hanya bila `.TanggalKonfirmasi` terisi.
func MaxTBC(konfirmasi *time.Time, tbc *int) *time.Time {
	if konfirmasi == nil || tbc == nil {
		return nil
	}
	t := konfirmasi.AddDate(0, 0, *tbc)
	return &t
}

// NamaTypeCeding meniru `CARI8` InputOfferLife_ACT langkah 3 - "" di luar 1-4.
func NamaTypeCeding(kode string) string { return namaPilihan(PilihanTypeCeding, kode) }

// NamaBusiness meniru `SetCoBName_Act` - `UNKNOWN` di luar L1-L21.
func NamaBusiness(kode string) string {
	if n := namaPilihan(PilihanClassOfBusiness, kode); n != "" {
		return n
	}
	return NamaBusinessTakDikenal
}

// Nilai "Reinsurance Type" - VERBATIM `SetReinsuranceType`.
const (
	JenisAsuransiProporsional    = "Proportional"
	JenisAsuransiNonProporsional = "Non Proportional"
)

// JenisAsuransi meniru `SetReinsuranceType` / InputOfferLife_ACT langkah 7:
// `TypeCeding "4"` (XOL) -> Non Proportional, SELAIN ITU Proportional.
//
// ⛔ Kosong pun menjadi Proportional: cabangnya `OTHERWISE`, bukan
// `WHEN TypeCeding != "4"`. Pega menyimpan hasilnya (langkah 8 `Obj-Save`);
// di sini ia dihitung ulang setiap simpan dan ditulis ke `JENIS_ASURANSI`
// (migrasi 061) - TIDAK PERNAH diterima dari klien.
//
// ⚠️ RALAT 01-10-2026: sempat diubah menjadi dropdown yang dipilih pemakai;
// work owner membatalkannya di hari yang sama - kembali ke turunan Pega.
func JenisAsuransi(typeCeding string) string {
	if typeCeding == "4" {
		return JenisAsuransiNonProporsional
	}
	return JenisAsuransiProporsional
}

// PilihanStatusMenurutBendera - radio mana yang tampil.
func PilihanStatusMenurutBendera(flag string) []Pilihan {
	if flag == FlagPolisPremium {
		return PilihanEmailTypePL
	}
	return PilihanStatusOffer
}

// SusunPenawaran memeriksa isian lalu menurunkan nama-namanya.
//
// ⛔ Wajib-isi mengikuti `pyRequired` sel (TypeCeding, BusinessCode, Comment),
// dan SELURUH yang kosong dilaporkan sekaligus - pola `ProtectAccept`.
// Kode di luar pilihan tertutup ditolak lebih dahulu: ia cacat klien, bukan
// isian yang terlupa.
func SusunPenawaran(flag string, isi IsianPenawaran) (PenawaranTersimpan, error) {
	isi = rapikan(isi)
	if isi.TypeCeding != "" && NamaTypeCeding(isi.TypeCeding) == "" {
		return PenawaranTersimpan{}, fmt.Errorf("%w: System Reinsurance %q", ErrIsianPenawaranTidakSah, isi.TypeCeding)
	}
	if isi.BusinessCode != "" && namaPilihan(PilihanClassOfBusiness, isi.BusinessCode) == "" {
		return PenawaranTersimpan{}, fmt.Errorf("%w: Class of Business %q", ErrIsianPenawaranTidakSah, isi.BusinessCode)
	}
	if isi.Status != "" && namaPilihan(PilihanStatusMenurutBendera(flag), isi.Status) == "" {
		return PenawaranTersimpan{}, fmt.Errorf("%w: Status %q", ErrIsianPenawaranTidakSah, isi.Status)
	}
	// Pasangan kode-nama dipilih bersama dari popup; separuh pasangan adalah
	// cacat klien, bukan isian.
	if (isi.CedingCo == "") != (isi.CedingCoName == "") {
		return PenawaranTersimpan{}, fmt.Errorf("%w: Ceding tanpa pasangan kode-nama", ErrIsianPenawaranTidakSah)
	}
	if (isi.PolicyHolder == "") != (isi.PolicyHolderName == "") {
		return PenawaranTersimpan{}, fmt.Errorf("%w: Policy Holder tanpa pasangan kode-nama", ErrIsianPenawaranTidakSah)
	}
	// Bilangan bulat yang negatif bukan umur maupun jumlah hari. ⚠️ Pega tidak
	// memeriksanya (pxNumber tanpa validasi); ditolak di sini karena kolomnya
	// NUMBER(5) dan Max TBC yang mundur dari Confirmation Date tidak bermakna.
	for nama, v := range map[string]*int{"Age Limit": isi.BatasUsiaPeserta, "Input TBC": isi.TBC} {
		if v != nil && (*v < 0 || *v > 99999) {
			return PenawaranTersimpan{}, fmt.Errorf("%w: %s %d", ErrIsianPenawaranTidakSah, nama, *v)
		}
	}
	if kurang := KekuranganPenawaran(WajibPenawaran{
		CedingCoName: isi.CedingCoName, PolicyHolderName: isi.PolicyHolderName,
		TypeCeding: isi.TypeCeding, BusinessCode: isi.BusinessCode,
		Status: isi.Status, Description: isi.Description, PeriksaStatus: true,
	}); len(kurang) > 0 {
		return PenawaranTersimpan{}, fmt.Errorf("%w: %s", ErrIsianPenawaranBelumLengkap, GabungPesanPenawaran(kurang))
	}
	return PenawaranTersimpan{
		IsianPenawaran: isi,
		TypeCedingName: NamaTypeCeding(isi.TypeCeding),
		BusinessName:   NamaBusiness(isi.BusinessCode),
		JenisAsuransi:  JenisAsuransi(isi.TypeCeding),
		TanggalTBC:     MaxTBC(isi.TanggalKonfirmasi, isi.TBC),
	}, nil
}

// WajibPenawaran - sel wajib layar Input Offer.
//
// `PeriksaStatus` false untuk gerbang `Confirm`: Status tidak punya kolom di
// header (ia tercatat di riwayat `T_VIEW_SUGGEST`), jadi yang dapat diperiksa
// dari data tersimpan hanya kelima lainnya. Saat menyimpan ia SELALU diperiksa.
type WajibPenawaran struct {
	CedingCoName, PolicyHolderName, TypeCeding, BusinessCode, Status, Description string
	PeriksaStatus                                                                 bool
}

// KekuranganPenawaran - sel wajib yang kosong, berurutan seperti di layar.
// Dipakai saat menyimpan DAN sebelum `Confirm` tahap penawaran.
//
// System Reinsurance, Class of Business, Comment - `pyRequired` korpus.
// Ceding Name, Policy Holder, Status - keputusan work owner 01-10-2026.
func KekuranganPenawaran(w WajibPenawaran) []string {
	var kurang []string
	for _, c := range []struct {
		nilai, pesan string
		periksa      bool
	}{
		{w.CedingCoName, PesanCedingKosong, true},
		{w.PolicyHolderName, PesanPolicyHolderKosong, true},
		{w.TypeCeding, PesanTypeCedingKosong, true},
		{w.BusinessCode, PesanBusinessCodeKosong, true},
		{w.Status, PesanStatusKosong, w.PeriksaStatus},
		{w.Description, PesanCommentKosong, true},
	} {
		if c.periksa && kosong(c.nilai) {
			kurang = append(kurang, c.pesan)
		}
	}
	return kurang
}

// BarisSuggest adalah satu baris riwayat `T_VIEW_SUGGEST`.
type BarisSuggest struct {
	No              int       `json:"no"`
	DateSuggest     time.Time `json:"dateSuggest"`
	PICSuggest      string    `json:"picSuggest"`
	IsCedingConfirm string    `json:"isCedingConfirm"`
	CommentSuggest  string    `json:"commentSuggest"`
	InitialSuggest  string    `json:"initialSuggest"`
}

// SuggestBaru meniru `AddHistorySuggest` untuk satu simpanan.
//
//	bendera "0"  langkah 1: IsCedingConfirm = .Status,  Initial "Offer"
//	bendera "1"  langkah 2: IsCedingConfirm = @if(EmailTypePL=1,"Accept",
//	                        @if(EmailTypePL=2,"Reject","Decline")), Initial "Bind"
//
// ⛔ Cabang `else` langkah 2 ditiru apa adanya: EmailTypePL yang KOSONG pun
// tercatat "Decline". `No` diisi pemanggil (urutan berikut di basis data) -
// padanan `pxListSubscript` baris yang baru ditambahkan.
func SuggestBaru(flag, status, description, pic string, saat time.Time) BarisSuggest {
	b := BarisSuggest{DateSuggest: saat, PICSuggest: pic, CommentSuggest: strings.TrimSpace(description)}
	if flag == FlagPolisPremium {
		b.InitialSuggest = InitialSuggestBind
		switch status {
		case "1":
			b.IsCedingConfirm = "Accept"
		case "2":
			b.IsCedingConfirm = "Reject"
		default:
			b.IsCedingConfirm = "Decline"
		}
		return b
	}
	b.InitialSuggest = InitialSuggestOffer
	b.IsCedingConfirm = status
	return b
}

func namaPilihan(daftar []Pilihan, kode string) string {
	for _, p := range daftar {
		if p.Kode == kode {
			return p.Nama
		}
	}
	return ""
}

func rapikan(i IsianPenawaran) IsianPenawaran {
	for _, s := range []*string{&i.CedingCo, &i.CedingCoName, &i.PolicyHolder, &i.PolicyHolderName,
		&i.TypeCeding, &i.BusinessCode, &i.Description, &i.Status,
		&i.PeriodePertanggungan, &i.StatusUpdate, &i.KeteranganMarketing,
		&i.QQName, &i.JenisUsaha, &i.KetentuanUnderwriting, &i.StatusFinal} {
		*s = strings.TrimSpace(*s)
	}
	return i
}
