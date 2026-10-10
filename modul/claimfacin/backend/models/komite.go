package models

// Untuk apa berkas ini: BATAS KOMITE di model Claim Fac In (pola Claim Prop) - grid "Committee Accept Status"
// (`.ComiteeClaim` InputAdjustment: roster calon `SetListKomite_act`, atau tangga kasus komite tersimpan), proteksi
// "Send to Committe" (`SendPICProtect_Act` + `CekPremiLunas_Act`), pop-up Comittee (`ClaimComite`, pra-proses
// `SetRemarksKomite`), dan penandaan kelahiran kasus komite TT2 (`CreateKMTNo_Act`).
//
// ⛔ Nama properti keputusan anggota komite HANYA ditulis di berkas ini; berkas lain memakai konstantanya. Penjaga batas
// Claim Life `komite_statik_test.go` mengecualikan berkas ini (izin work owner 09-10-2026, pola Claim Prop / Non Prop):
// keputusan tetap DITULIS konteks Komite (tahap 2), Claim Fac In menulis tangga awal dan membacanya untuk tampilan.

import (
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// Properti baris `.ComiteeClaim` yang membawa keputusan anggota (InputAdjustment LS52-LS53).
const (
	PropKeputusanAnggota = "KomiteAproval"
	PropTanggalKeputusan = "DateApprove"
	PropCatatanKeputusan = "KomiteComment"
)

// ApprovalKomiteMenunggu - `KomiteAproval = 0`: roster calon dan anggota tangga yang belum memutuskan.
const ApprovalKomiteMenunggu = "0"

// TransferType kasus komite (`childPageKomite.TransferType`, kolom T_GENERAL_KOMITE.TRANSFER_TYPE migrasi
// komiteclaimfacin 642): 2 adjustment (CreateKMTNo_Act 6.9), 3 Reject Claim (SendRejectClaimToKomite2 7.3), 4 Close
// Without Payment (SendCloseClaimToKomite 7.3).
const (
	TransferAdjustment = "2"
	TransferTolak      = "3"
	TransferTutup      = "4"
)

// JabatanBatasFacOut - JABATAN roster yang LIMIT_TOP-nya membatasi tangga objek retro (SetListKomite_act 4.1).
const JabatanBatasFacOut = "Technic Div. Head"

func init() {
	// ASM-FW-GCNMFW-Data-Comitee.KomiteAproval - urutan dan label ekspor work owner `Komite Claim Prop/KomiteAproval.xml`
	// (kelas yang sama, pola Claim Prop / Non Prop).
	KodePilihan[PropKeputusanAnggota] = []string{"1", "2", ApprovalKomiteMenunggu}
	LabelKode[PropKeputusanAnggota] = map[string]string{"1": "Approved", "2": "Reject", ApprovalKomiteMenunggu: "Waiting"}
}

// gridKomite - grid "Committee Accept Status" (InputAdjustment LS52-LS53, hanya-baca).
func gridKomite(daftar string) Unsur {
	return Unsur{Jenis: JenisGrid, Jalur: daftar, Label: "Committee Accept Status", Kolom: []Unsur{
		kRO(kol(PropNomorBaris, "", KTampil)),
		kRO(kol("IDKomite", "Committee Name", KTampil)),
		kSumber(kRO(kol(PropKeputusanAnggota, "Status", KPilih)), kode(PropKeputusanAnggota)),
		kRO(kol(PropTanggalKeputusan, "Date Approve", KWaktu)),
		kRO(kol(PropCatatanKeputusan, "Comment", KTampil)),
	}}
}

// RosterKomiteCalon = `SetListKomite_act` (lewat SetKomiteList_ACT / CreateKMTNo_Act 8): roster FACIN aktif
// ber-LIMIT_BOTTOM <= batas; batas 1; objek retro -> batas = ValueAdjustment, tanpa saringan batas bila melampaui
// LIMIT_TOP baris "Technic Div. Head" (4.3 `Param.LIMIT_BOTTOM := ""`: RD mengabaikan saringan berparameter kosong).
func RosterKomiteCalon(h *Halaman, o, i, a int, roster []AnggotaKomite) ([]AnggotaKomite, error) {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return nil, err
	}
	var kk Kalkulator
	batas := apd.New(1, 0)
	tanpaBatas := false
	if ob, err := Objek(h, o); err == nil && ob["IsFacretro"] == "1" { // 4
		var top string
		for _, r := range roster {
			if r.Jabatan == JabatanBatasFacOut {
				top = r.LimitTop
				break
			}
		}
		batas = kk.B(b, "ValueAdjustment")
		if lewat, err := lebihDari(b["ValueAdjustment"], top); err == nil && lewat { // 4.3
			tanpaBatas = true
		}
	}
	var out []AnggotaKomite
	for _, r := range roster { // 6
		if !tanpaBatas {
			if strings.TrimSpace(r.LimitBottom) == "" {
				continue // RD `LIMIT_BOTTOM <= param`: NULL tidak lolos
			}
			if Lebih(kk.Teks("LIMIT_BOTTOM", r.LimitBottom), batas) {
				continue
			}
		}
		out = append(out, r)
	}
	return out, kk.Galat()
}

// SusunKomiteAdjustment - `.ComiteeClaim` adjustment a: `tangga` (dibaca dari kasus komite tersimpan) bila adjustment
// sudah diserahkan, selainnya roster calon (`RosterKomiteCalon`, KomiteAproval 0, IDKomite = JABATAN - SetListKomite_act
// 6.1). Total komite = cacah barisnya (7).
func SusunKomiteAdjustment(h *Halaman, o, i, a int, roster, tangga []AnggotaKomite) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	anggota := tangga
	if tangga == nil {
		if anggota, err = RosterKomiteCalon(h, o, i, a, roster); err != nil {
			return err
		}
	}
	var rows []Baris
	for _, r := range anggota {
		appr := r.Approval
		if appr == "" {
			appr = ApprovalKomiteMenunggu
		}
		rows = append(rows, Baris{"KomiteID": r.OperatorID, "IDKomite": r.Jabatan, "KomiteEmail": r.Email,
			PropKeputusanAnggota: appr, PropTanggalKeputusan: r.TanggalSetuju, PropCatatanKeputusan: r.Comment})
	}
	h.SetelDaftar(DaftarDiAdj(o, i, a, AnakKomiteAdj), rows)
	b["TotalKomite"] = strconv.Itoa(len(rows))
	return nil
}

// KomiteMenunggu - ValidationAdjustmentKomite 2.1.2.5.1: cacah anggota `.ComiteeClaim` adjustment yang belum memutuskan
// (KomiteAproval kosong / 0).
func KomiteMenunggu(komite []Baris) int {
	n := 0
	for _, k := range komite {
		if v := k[PropKeputusanAnggota]; v == "" || v == ApprovalKomiteMenunggu {
			n++
		}
	}
	return n
}

// ---------------------------------------------------------------- Send to Committe (TT2)

// Kategori lampiran wajib (`AttachCategory.pxResults.ID`, SendPICProtect_Act 3).
const (
	LampiranLOD     = "LOD"
	LampiranDLA     = "DLA"
	LampiranSPGR    = "SPGR"
	LampiranADU     = "ADU"
	LampiranInvoice = "Invoice"
	LampiranSalvage = "Salvage"
)

// Pesan VERBATIM SendPICProtect_Act 2 / CekPremiLunas_Act 6.
var pesanLampiran = map[string]string{
	LampiranLOD: "Please Upload Attachment LOD", LampiranDLA: "Please Upload Attachment DLA",
	LampiranSPGR: "Please Upload Attachment SPGR", LampiranADU: "Please Upload Attachment Approval Direktur Utama",
	LampiranInvoice: "Please Upload Attachment Invoice", LampiranSalvage: "Please Upload Attachment Salvage",
}

// Pesan proteksi VERBATIM.
const (
	PesanBankKosong      = "Any Bank Data is Null, Please Check Again / Re-Select Bank Account"
	PesanPremiBelumLunas = "Akseptasi tidak dapat dilanjutkan dikarenakan Premi belum Lunas"
)

// ProteksiKomite = `SendPICProtect_Act` 1-18 (adjustment a): lampiran wajib per Payment Type, rekening, batas Direktur
// Utama (ADU). `lampiran` = cacah berkas per kategori (`AttachCategory.pxResults`, master FAC x DOCUMENT_CLAIM);
// `limitDirut` = LIMIT_BOTTOM baris Direktur Utama (GetLimitDirekturUtama_SQL - nama orang dibuang, pola Claim Prop
// DEGREE 6). Mengembalikan `Protect.CARI1`.
//
// Langkah 7-13 hanya pesan pada medan IsError (tidak menghentikan); pop-up Comittee tetap terbuka dan tombol "Send
// Claim to Committee" tersembunyi bila CARI1 / CARI2 bukan 1 (ClaimComite LS7).
func ProteksiKomite(h *Halaman, o, i, a int, lampiran map[string]int, limitDirut string) (bool, error) {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return false, err
	}
	h.BersihkanPesan() // 1
	ada := func(kat string) bool { return lampiran[kat] > 0 }
	var kk Kalkulator
	total := apd.New(0, 0) // 4.1
	if b["AcceptanceStatus"] != "2" {
		total = kk.B(b, "ValueAdjustment")
	}
	pt := ""
	if b["AcceptanceStatus"] == "" { // 4.2
		pt = b["PaymentType"]
	}
	bank := b["NameOfBank"] != "" && b["NoAccount"] != "" && b["IDOfBank"] != "" // 4.3
	lewat := false                                                               // 5-6 `TotalAdjustment >= LIMIT_BOTTOM`
	if strings.TrimSpace(limitDirut) != "" {
		lewat = !Lebih(kk.Teks("LIMIT_BOTTOM", limitDirut), total)
	}
	if err := kk.Galat(); err != nil {
		return false, err
	}
	pesan := func(kat string) { h.TambahPesan(JalurIsError, pesanLampiran[kat]) }
	dasar := pt == BayarFinal || pt == BayarInterim
	if dasar && !ada(LampiranLOD) { // 7
		pesan(LampiranLOD)
	}
	if (dasar || pt == BayarSalvage || pt == BayarFee || pt == BayarExpense) && !ada(LampiranDLA) { // 8
		pesan(LampiranDLA)
	}
	if dasar && !ada(LampiranSPGR) { // 9
		pesan(LampiranSPGR)
	}
	if lewat && dasar && !ada(LampiranADU) { // 10
		pesan(LampiranADU)
	}
	if (pt == BayarFee || pt == BayarExpense) && !ada(LampiranInvoice) { // 11
		pesan(LampiranInvoice)
	}
	if pt == BayarSalvage && !ada(LampiranSalvage) { // 12
		pesan(LampiranSalvage)
	}
	if !bank { // 13
		h.TambahPesan(JalurIsError, PesanBankKosong)
	}
	lolos := false
	switch { // 14-17 (ADU 14: `(!L && !A) || (L && A) || A` - lewat batas Direktur Utama menuntut ADU)
	case dasar:
		lolos = ada(LampiranLOD) && ada(LampiranDLA) && ada(LampiranSPGR) &&
			((!lewat && !ada(LampiranADU)) || (lewat && ada(LampiranADU)) || ada(LampiranADU))
	case pt == BayarSalvage:
		lolos = ada(LampiranSalvage) && ada(LampiranDLA)
	case pt == BayarFee || pt == BayarExpense:
		lolos = ada(LampiranInvoice) && ada(LampiranDLA)
	case pt == BayarAdjust:
		lolos = true
	}
	if lolos { // 18 `Protect.CARI1==0 [T=3]`: CARI1 yang sudah 1 ditimpa penanda rekening
		lolos = bank
	}
	return lolos, nil
}

// PerluCekPremi - SendPICProtect_Act 19.1: CekPremiLunas_Act hanya bila adjustment `.IsKomite` kosong; selainnya
// ParamData.HASIL1 kosong sehingga 19.2 menyetel CARI2 = 1.
func PerluCekPremi(b Baris) bool { return b["IsKomite"] == "" }

// PesanPremi = CekPremiLunas_Act 8.
func PesanPremi(h *Halaman) { h.TambahPesan("", PesanPremiBelumLunas) }

// KunciPremi - CekPremiLunas_Act 3 (jalur `CLM-`): nomor polis tanpa titik, mata uang adjustment, nomor polis.
func KunciPremi(h *Halaman, b Baris) (invoice, cur, nopolis string) {
	nopolis = h.Ambil(JalurNoPolis)
	return strings.ReplaceAll(nopolis, ".", ""), b["CurrencyID"], nopolis
}

// LayarKomite - Section ClaimComite (harness Comittee) untuk adjustment a. `boleh` = `Protect.CARI1 = 1 &&
// Protect.CARI2 = 1` dihitung server (SendPICProtect_Act) setiap evaluasi.
func LayarKomite(o, i, a int, boleh Kondisi) []Unsur {
	p := JalurAdj(o, i, a) + "."
	j := func(prop string) string { return p + "DataCommitteFacin." + prop }
	pt := func(h *Halaman) string { return AmbilJalur(h, p+"PaymentType") }
	salvage := func(h *Halaman) bool { return pt(h) == BayarSalvage }
	fee := func(h *Halaman) bool { t := pt(h); return t == BayarFee || t == BayarExpense }
	isi := func(prop string) func(h *Halaman) bool {
		return func(h *Halaman) bool { return AmbilJalur(h, j(prop)) != "" }
	}
	adaRemarks := isi("Remarks")
	kirim := func(h *Halaman) bool { // LS7 `[r1]` VIS
		t := pt(h)
		switch {
		case t == BayarSalvage:
			return isi("Salvage")(h) && adaRemarks(h)
		case t == BayarFee || t == BayarExpense:
			return isi("AdjusterFee")(h) && adaRemarks(h)
		}
		return adaRemarks(h)
	}
	return []Unsur{
		sebaris("",
			naJika(medan(JalurTanggalKomite, "Date", KTanggal), selalu),
			ro(medan(JalurInisialKomite, "Initial", KTeks)),
		),
		tampil(label("Adjustment"), func(h *Halaman) bool {
			t := pt(h)
			return t == BayarFinal || t == BayarInterim || t == BayarAdjust || t == "7"
		}),
		tampil(label("Salvage"), salvage),
		tampil(label("Adjuster / Consultant Fee"), fee),
		medan(j("CircumCauseOfLoss"), "Chronology", KArea),
		medan(j("ExtentOfLoss"), "Extent of Loss", KArea),
		medan(j("LegalLiability"), "Legal Liability", KArea),
		tampil(wajibU(medan(j("Salvage"), "Salvage", KArea)), salvage),
		tampil(wajibU(medan(j("AdjusterFee"), "Adjuster / Consultant Fee", KArea)), fee),
		wajibU(medan(j("Remarks"), "Remarks", KArea)),
		tampil(tampil(tombol("KirimKomite", "Send Claim to Committee", "KirimKomite"), kirim), boleh),
	}
}

// Jalur medan pop-up Comittee (halaman requestor `TempCommiteClaim` / `TreatyExchangeYearly`, tidak disimpan).
const (
	JalurTanggalKomite = "TempCommiteClaim.DateOfComitee"
	JalurInisialKomite = "TreatyExchangeYearly.UserName"
)

// SetRemarksKomite = pra-proses harness Comittee (`SetRemarksKomite` 1, adjustment bukan baris pertama): Chronology dan
// Remarks disalin dari adjustment sebelumnya. Inisial = pembuat adjustment (CreateKMTNo_Act 7), tanggal = hari ini
// (ClaimComite LS2 DEF).
func SetRemarksKomite(k *Konteks, h *Halaman, o, i, a int) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	if a > 1 {
		lalu, err := Adj(h, o, i, a-1)
		if err != nil {
			return err
		}
		b["DataCommitteFacin.CircumCauseOfLoss"] = lalu["DataCommitteFacin.CircumCauseOfLoss"]
		b["DataCommitteFacin.Remarks"] = lalu["DataCommitteFacin.Remarks"]
	}
	h.Setel(JalurTanggalKomite, k.Hari())
	h.Setel(JalurInisialKomite, b["pxCreateOpName"])
	return nil
}

// PesanPembayaran = CreateKMTNo_Act 15 `Local.Messages` (Progres1 sub-progress "Waiting Committee").
func PesanPembayaran(pt string) string {
	switch pt {
	case BayarFinal:
		return "Payment - Final"
	case BayarInterim:
		return "Payment - Claim"
	case BayarSalvage:
		return "Payment - Salvage"
	case BayarFee:
		return "Payment - Adjuster Fee"
	case BayarAdjust:
		return "Payment - Adjustment"
	}
	return "Payment - Consultant Fee"
}

// Teks VERBATIM CreateKMTNo_Act.
const (
	AwalanKirimKomite  = "Send to Committee - " // 15 Data.CARI12
	ProgresMenungguKmt = "Waiting Committee"    // 16 Progres2
	ProgresAutoCreate  = "Auto Create"          // 16 Comment
)

// TandaiKirimKomite = CreateKMTNo_Act 12-15, 17 sesudah kasus komite `kmt` lahir: item `.IsKomite` 0 (tombol "+"
// terbuka lagi), adjustment `.IsKomite` 1 dan tertaut (KOMITE_ID = `KomiteNo`, 13.1.1), kronologi.
//
// Tidak ditiru (PARITAS): 4 memotong tanggal polis / DateOfLoss / pxCreateDateTime ke 8 karakter lalu menyimpannya
// (perbaikan prompt §5 butir 6); 11 menimpa TempCommiteClaim.CircumtansesCouseOfLoss dengan pzInsKey (§5 butir 7);
// 6.8 menyertakan adjustment PT 2 lain yang belum berkomite ke KMT yang sama - indeks unik
// `UQ_CLAIM_ADJUSTMENT_KOMITE` menolak dua adjustment bertaut ke satu KMT (OQ-CFI-28, tanpa MODIFY); 6.2-6.5
// TotalEstimasiReas / TFAllObj FacRetroList polis (halaman polis tidak disimpan, pembacanya hanya HTML).
func TandaiKirimKomite(k *Konteks, h *Halaman, o, i, a int, kmt string) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	it["IsKomite"] = "0" // 12
	b["IsKomite"] = "1"
	b[PropKomiteID] = kmt                 // 13.1.1
	k.Kronologi(h, AwalanKirimKomite+kmt) // 15, 17
	return nil
}

// ModalKomite - prefiks kunci pop-up Comittee (kunci = `KunciPanel(ModalKomite, DaftarAdj(o, i), a)`).
const ModalKomite = "komite"

// Penanda SendPICProtect_Act (halaman requestor `Protect`, dibawa ModeLayar; dihitung ulang server sebelum
// penyerahan).
const (
	JalurProtect1 = "Protect.CARI1"
	JalurProtect2 = "Protect.CARI2"
)

// PremiBelumLunas = CekPremiLunas_Act 5-6: saldo (koma -> titik) > 0 = "BELUM LUNAS"; kosong = lunas.
func PremiBelumLunas(saldo string) (bool, error) {
	saldo = strings.ReplaceAll(strings.TrimSpace(saldo), ",", ".")
	if saldo == "" {
		return false, nil
	}
	var kk Kalkulator
	v := kk.Teks("HASIL1", saldo)
	if err := kk.Galat(); err != nil {
		return false, err
	}
	return Lebih(v, apd.New(0, 0)), nil
}
