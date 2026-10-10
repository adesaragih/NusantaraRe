package services

// Untuk apa berkas ini: SUBMIT - tombol "Submit" (`finishAssignment`) flow action `ViewTransferDtl`, pasca-proses
// `KomitePostAct` (korpus `Komite Claim FacIn`), dalam SATU transaksi. Rencana langkah disusun `models.Rencanakan`
// (murni); berkas ini menjalankan langkah yang menyentuh basis data, kontrak Claim Fac In, dan outbox:
//
//	TT2 KomitePost_Adjustment                         TT3 KomitePost_Reject / TT4 KomitePost_CloseClaim
//	S2         kunci klaim induk (kontrak)            S1         kunci klaim induk
//	(ApprovalKomite_Act S4-S8) perluasan tangga       S6 / S5    keputusan tangga
//	S7.2.1.2   nomor akseptasi (penomor)              S3-S5      kronologi (kontrak)
//	S7.2.1.4-9 tangga + adjustment (kontrak)          S11 / S10  JSON_KLAIM
//	S7.2.1.15  email (outbox, IsPEGAPROD)             S12        OS STS 2 per estimasi / OS STS 4
//	S7.2.1.16  SaveAcceptation_KMT (kontrak)          S13 / S11  PDF (stream tidak diekspor, OQ) + email (outbox)
//	S8         OS_AKSEPTASI_KLAIM                     S14        CLAIMREJECTED (TT3)
//	S12-S13    penanda cetak / usul (kontrak)         S15 / S12.4 konversi Arasapas (outbox)
//	S16        JSON_KLAIM                             S16 / S13  KomiteCount + 1
//	S17-S19    konversi (outbox) + log AKSEPATSI      S17 / S14  tutup klaim induk (kontrak)
//	S20-S23    HISTORYAKSEPTASIPEGA, SUBPROGRESSCLAIM
//	S24        KomiteCount + 1;  S25 Kasir (outbox)
//
// PENYIMPANGAN SADAR (prompt §5 butir 6, pola Komite Claim Prop): di Pega setiap rule menyimpan sendiri (`COMMIT` di
// RDB, Obj-Save per iterasi). Di sini semuanya satu transaksi aplikasi: keputusan yang gagal di tengah batal UTUH.
// Panggilan keluar (Arasapas, Kasir, email) diantre di transaksi yang sama dan baru berjalan sesudah commit (outbox,
// hanya produksi). Dokumen PDF tidak dibuat (stream tidak diekspor, OQ-KCFI-01).

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimfacin/backend/models"
)

// OQDokumenPDF - pesan info Submit yang menerbitkan dokumen PDF di Pega (stream tidak diekspor).
const OQDokumenPDF = "OQ-KCFI-01: dokumen PDF komite (stream AcceptanceNotePDF / CommitteReject_CC / " +
	"CommitteCloseClaim) tidak diekspor korpus - berkas belum dibuat"

// StatusKlaimTerbuka - `TempOpenPage.pyStatusWork` kasus klaim yang masih terbuka saat CLAIMREJECTED ditulis
// (KomitePost_Reject S14 sebelum S17 pxForceCaseClose; DEV CLAIMREJECTED `CLM-` STATUSWORK "New").
const StatusKlaimTerbuka = "New"

// LabelKlaim - `TempOpenPage.pyLabel` (DEV CLAIMREJECTED `CLM-` LABEL "CLM").
const LabelKlaim = "CLM"

// HasilKeputusan - jawaban satu Submit.
type HasilKeputusan struct {
	KasusID string `json:"kasusId"`
	// Selesai - decision IsKomiteLoop salah: kasus Resolved-Completed.
	Selesai     bool   `json:"selesai"`
	KomiteCount int    `json:"komiteCount"`
	AcceptedNo  string `json:"acceptedNo,omitempty"`
	Info        string `json:"info,omitempty"`
}

// Putuskan menjalankan satu Submit `kep` oleh `p` atas kasus komite `id`.
func (l *Layanan) Putuskan(ctx context.Context, p inti.Pelaku, id string, kep models.Keputusan) (HasilKeputusan, error) {
	if err := l.siap(p); err != nil {
		return HasilKeputusan{}, err
	}
	saat := l.jam()
	var hasil HasilKeputusan
	err := l.g.Transaksi(ctx, func(tx *db.Tx) error {
		k, kl, err := l.muat(ctx, tx, id, true)
		if err != nil {
			return err
		}
		if k.Tertutup() {
			return ErrKasusTertutup
		}
		if !k.Pemegang(p.AkunID, p.Peran) {
			return ErrBukanPemegang
		}
		if pesan := models.PeriksaIsian(kep); len(pesan) > 0 {
			return &GalatValidasi{Pesan: pesan}
		}
		j := &jalan{l: l, ctx: ctx, tx: tx, k: k, kl: kl, akun: p.AkunID, saat: saat}
		if k.TransferType == models.TransferAdjustment {
			if kl.Adjustment == 0 {
				return ErrKlaimIndukTidakAda
			}
			if st := models.Adjustment(kl)["AcceptanceStatus"]; st != "" && st != models.KeputusanMenunggu {
				// LS48 Submit `NA pyWorkPage.Adjustment.AcceptedNo != ''`; S7.2.1 hanya adjustment belum diputus
				return ErrKasusTertutup
			}
			if err := j.perluas(p.Peran); err != nil {
				return err
			}
		}
		if j.nama, err = l.a.NamaPelaku(ctx, p.AkunID); err != nil { // S20 OperatorID.pyUserName
			return err
		}
		r := models.Rencanakan(j.k, kl, kep, p.AkunID, saat)
		if err := j.laksanakan(&r); err != nil {
			return err
		}
		hasil = HasilKeputusan{KasusID: k.ID, Selesai: r.Selesai, KomiteCount: r.Kepala.Count, AcceptedNo: j.nomor}
		if r.PDF != "" {
			hasil.Info = OQDokumenPDF
		}
		return nil
	})
	if err != nil {
		return HasilKeputusan{}, err
	}
	return hasil, nil
}

// jalan - keadaan satu Submit.
type jalan struct {
	l     *Layanan
	ctx   context.Context
	tx    *db.Tx
	k     models.Kasus
	kl    kontrak.KlaimFacIn
	akun  string
	nama  string
	saat  time.Time
	nomor string
}

// perluas = ApprovalKomite_Act S4-S8 disimpan saat tingkat 1 Submit (KCF-02).
func (j *jalan) perluas(peran []string) error {
	pr, err := j.l.praProses(j.ctx, j.k, j.kl)
	if err != nil {
		return err
	}
	baru, err := j.l.perluasan(j.ctx, j.k, pr, peran)
	if err != nil || len(baru) == 0 {
		return err
	}
	tersimpan, err := j.l.g.TambahAnggota(j.ctx, j.tx, j.k.ID, baru)
	if err != nil {
		return err
	}
	j.k.Tangga = append(j.k.Tangga, tersimpan...)
	j.k.Loop = len(j.k.Tangga) // S8
	return nil
}

func (j *jalan) laksanakan(r *models.Rencana) error {
	ctx, tx, k, kl := j.ctx, j.tx, j.k, j.kl
	for _, u := range r.Tangga {
		if err := j.l.g.TulisAnggota(ctx, tx, k.ID, u); err != nil {
			return err
		}
	}
	kunci := models.KunciInstans(k.KlaimID)
	noPolis := kl.Nilai["OfferFacIn.PolicyData.PolicyNo"]
	noKlaim := kl.Nilai["ClaimData.NoClaim"]
	if r.TerbitkanNomor { // S7.2.1.2.4-S7.2.1.2.7
		bn, err := j.l.g.UrutNomorAkseptasi(ctx, tx, j.saat)
		if err != nil {
			return err
		}
		j.nomor = models.RakitNomorAkseptasi(bn.Jenis, kl.Nilai[models.OQ+"BusinessOldId"], bn.MMYYYY, bn.Urut)
		r.Klaim.Adjustment["AcceptedNo"] = j.nomor
	}
	adj := models.Adjustment(kl)
	sesudah := make(map[string]string, len(adj))
	for p, v := range adj {
		sesudah[p] = v
	}
	for p, v := range r.Klaim.Adjustment {
		sesudah[p] = v
	}
	if r.SimpanOS { // S7.2.1.17 + S8
		retro := models.Objek(kl)["IsFacretro"]
		if v, ada := r.Klaim.Objek["IsFacretro"]; ada {
			retro = v
		}
		if err := j.l.g.SisipOS(ctx, tx, models.SusunOSAkseptasi(kl, sesudah, retro, k.ID, j.saat), j.saat); err != nil {
			return err
		}
	}
	if r.OSTolak { // KomitePost_Reject S12 (SaveReject_ACT_KMT)
		for _, b := range models.SusunOSTolak(kl, j.saat) {
			if err := j.l.g.SisipOS(ctx, tx, b, j.saat); err != nil {
				return err
			}
		}
	}
	if r.OSTutup { // KomitePost_CloseClaim S12.1-S12.3
		if err := j.l.g.SisipOS(ctx, tx, models.SusunOSTutup(kl), j.saat); err != nil {
			return err
		}
	}
	if r.JSONKlaim { // S16 / Reject S11 / Close S10
		if err := j.l.g.SalinJSONKlaim(ctx, tx, kunci, noKlaim, noPolis, j.saat); err != nil {
			return err
		}
	}
	if r.Konversi != "" { // S17 / Reject S15 / Close S12.4 KonversiKlaim_Act (Connect-REST bergerbang IsPEGAPROD)
		if err := j.antre(JenisEfekKonversi, k.KlaimID, map[string]string{"CASEID": kunci, "NOPOLIS": noPolis,
			"STS_REJECT": r.Konversi}); err != nil {
			return err
		}
	}
	if r.LogAkseptasi { // S18-S19 (bukan Fac Retro); status / jawaban layanan menyusul (outbox)
		if err := j.l.g.CatatLogLayanan(ctx, tx, models.LogLayanan{IDPega: kunci,
			Parameter: kunci + " / " + noPolis + " / " + models.StsOSAkseptasi, JenisService: models.JenisLogAkseptasi,
			NoAkseptasi: j.nomor}, j.saat); err != nil {
			return err
		}
	}
	if r.Kasir { // S25 HitServiceToKasirKMT_Act
		if err := j.kasir(r, sesudah); err != nil {
			return err
		}
	}
	if r.KlaimDitolak { // KomitePost_Reject S14
		if err := j.l.g.SisipKlaimDitolak(ctx, tx, models.KlaimDitolak{InsKey: kunci, ID: k.KlaimID, InsName: k.KlaimID,
			Label: LabelKlaim, StatusWork: StatusKlaimTerbuka, PembuatNama: kl.Nilai[kontrak.JalurNamaPembuatKlaim],
			PembuatID: kl.Nilai[kontrak.JalurPembuatKlaim], Kelas: models.KelasKlaim, Diperbarui: j.saat,
			PengubahNama: j.nama, PengubahID: j.akun, Remark: kl.Nilai["ClaimData.Remark"]}); err != nil {
			return err
		}
	}
	if err := j.email(r); err != nil {
		return err
	}
	if err := j.l.klaim.TulisBalikKlaimFacIn(ctx, tx, k.KlaimID, k.AdjustmentID, bersihkan(r.Klaim)); err != nil {
		return galatKontrak(err)
	}
	if r.StatusRiwayat != "" { // S20-S21
		if err := j.l.g.CatatRiwayatAkseptasi(ctx, tx, models.RiwayatAkseptasi{IDPega: kunci,
			IDKomite: models.KunciInstans(k.ID), Status: r.StatusRiwayat, Username: j.nama,
			Workbasket: models.WorkbasketRiwayat}, j.saat); err != nil {
			return err
		}
	}
	if r.SubProgres != "" { // S22-S23
		if err := j.l.g.UbahSubProgres(ctx, tx, k.ID, r.SubProgres); err != nil {
			return err
		}
	}
	if r.TutupKlaim != "" { // Reject S17 / Close S14 pxForceCaseClose
		if err := j.l.klaim.TutupKlaimFacIn(ctx, tx, k.KlaimID, k.ID, r.TutupKlaim, j.saat); err != nil {
			return galatKontrak(err)
		}
	}
	if err := j.l.g.SimpanKepala(ctx, tx, k.ID, k.Count, r.Kepala); err != nil {
		return err
	}
	return j.l.g.TutupKasus(ctx, tx, k.ID, r.Selesai, r.Posisi, j.saat)
}

// bersihkan - peta ubahan kosong dijadikan nil (kontrak tidak menuntut posisi adjustment bila tanpa ubahan pohon).
func bersihkan(u kontrak.UbahanKlaimFacIn) kontrak.UbahanKlaimFacIn {
	if len(u.Header) == 0 {
		u.Header = nil
	}
	if len(u.Adjustment) == 0 {
		u.Adjustment = nil
	}
	if len(u.Objek) == 0 {
		u.Objek = nil
	}
	if len(u.Item) == 0 {
		u.Item = nil
	}
	return u
}

// kasir = HitServiceToKasirKMT_Act jalur CLM (S2 gerbang, S3 status konversi, S12-S13 IDOfBank, S14.2 panjang nomor +
// muatan, S14.4 REST = outbox). S14.5 StatusKasir / S14.6 DIRECTTOKASIR_LOG bergantung jawaban REST - ditulis pelaksana
// (menunggu persetujuan, OQ-CFI-26); dedupe kasir tetap OQ-CFI-26.
//
// ⚠️ S17 konversi kini efek outbox (asinkron, pola Komite Claim Prop): di produksi status konversi saat Submit bisa belum
// "1" sehingga kasir tidak diantre - tombol "Acceptation" Claim Fac In (HitServiceToKasir_Act) mengirimnya kemudian.
func (j *jalan) kasir(r *models.Rencana, adj map[string]string) error {
	if adj["AcceptanceStatus"] != models.KeputusanSetuju || !models.BolehKasir(adj) { // S25.2.1 / S2
		return nil
	}
	sts, err := j.l.a.StatusKonversi(j.ctx, strings.ReplaceAll(adj["AcceptedNo"], ".", "")) // S3 (hanya IsPEGAPROD)
	if err != nil {
		return err
	}
	if sts != "1" { // S3 transisi pasca-langkah: status konversi 1 lanjut (T=2), selainnya keluar (F=6)
		return nil
	}
	if adj["IDOfBank"] == "" { // S12-S13
		id, err := j.l.a.IDBankRekening(j.ctx, adj["NameOfBank"], adj["BranchOfBank"], adj["NoAccount"])
		if err != nil {
			return err
		}
		if id != "" {
			adj["IDOfBank"] = id
			r.Klaim.Adjustment["IDOfBank"] = id
		}
	}
	if !models.PanjangNoAksepCLM(adj["AcceptedNo"]) { // S14.2 F -> 6 keluar
		return nil
	}
	email, err := j.l.a.EmailCeding(j.ctx, models.KunciCedingKasir(j.kl)) // S14.2.2.1.1.2 (hanya produksi)
	if err != nil {
		return err
	}
	m, err := models.SusunMuatanKasir(j.kl, adj, email, j.akun, j.l.kasir, int(j.saat.In(models.Jakarta).Month()))
	if err != nil {
		return err
	}
	return j.antre(JenisEfekKasir, adj["AcceptedNo"], m)
}

// Kunci isi MUATAN efek email-komite.
const (
	IsiJenis    = "jenis"    // EmailPenyetujuBerikut / EmailPembuatSetuju / EmailPembuatTolak
	IsiPenerima = "penerima" // KomiteID penerima (S12 tingkat berikut) atau akun pembuat kasus komite
)

// email = SendEmailKlaim_KMT (TT2, bergerbang IsPEGAPROD; S1 tanpa SpreadingAdjustment -> keluar) / KomitePost_Reject
// S13.10 (TT3 / TT4). MUATAN hanya pengenal: subjek, alamat, akun, dan CC dirakit saat dikirim (`SusunEmailKomite`).
func (j *jalan) email(r *models.Rencana) error {
	if !j.l.produksi || r.Email == "" {
		return nil
	}
	k, kl := j.k, j.kl
	if k.TransferType == models.TransferAdjustment &&
		len(kl.Daftar[models.DaftarDiAdj(kl.Objek, kl.Item, kl.Adjustment, models.AnakSpread)]) == 0 { // S1
		return nil
	}
	penerima := k.PembuatID                      // S13 Obj-Browse pyUserIdentifier = pyWorkPage.pxCreateOperator
	if r.Email == models.EmailPenyetujuBerikut { // S12 ComiteeClaim(KomiteCount + 1)
		penerima = models.KomiteIDTingkat(k, k.Count+1)
	}
	return j.antre(JenisEfekEmailKomite, k.ID, map[string]string{IsiJenis: r.Email, IsiPenerima: penerima})
}

// Jenis efek outbox Komite Claim Fac In.
const (
	JenisEfekKonversi    = "konversi-klaim" // KonversiKlaim_Act -> KonversiKlaimNonLife (Klaim / insertClaimAccept)
	JenisEfekKasir       = "kasir"          // HitServiceToKasirKMT_Act -> SendAcceptationToKasir (Kasir / insertAllPaymentKasir)
	JenisEfekEmailKomite = "email-komite"   // SendEmailKlaim_KMT / KomitePost_Reject S13.10 -> SendEmailWithAttachments
)

// MuatanOutbox - isi baris T_LOG_SERVICE_RNM.MUATAN.
type MuatanOutbox struct {
	Kategori1 string `json:"kategori1,omitempty"`
	Kategori2 string `json:"kategori2,omitempty"`
	KlaimID   string `json:"klaimId"`
	KomiteID  string `json:"komiteId"`
	Isi       any    `json:"isi,omitempty"`
}

// kunciEfek - kunci M_LINK_SERVICE per jenis (Kategori_1/Kategori_2 VERBATIM dari activity).
var kunciEfek = map[string][2]string{
	JenisEfekKonversi: {"Klaim", "insertClaimAccept"},
	JenisEfekKasir:    {"Kasir", "insertAllPaymentKasir"},
}

// antre - efek keluar hanya di produksi (`When/IsPEGAPROD` di setiap Connect-REST / email).
func (j *jalan) antre(jenis, rujukan string, isi any) error {
	if !j.l.produksi {
		return nil
	}
	k := kunciEfek[jenis]
	teks, err := json.Marshal(MuatanOutbox{Kategori1: k[0], Kategori2: k[1], KlaimID: j.k.KlaimID, KomiteID: j.k.ID, Isi: isi})
	if err != nil {
		return err
	}
	_, err = j.l.g.AntreEfek(j.ctx, j.tx, jenis, rujukan, string(teks), j.saat)
	return err
}
