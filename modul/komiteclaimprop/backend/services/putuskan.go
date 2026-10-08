package services

// Untuk apa berkas ini: SUBMIT - tombol "Submit" (`finishAssignment`) flow action `ViewTransferDtl`, pasca-proses
// `KomitePost` S1 -> `KomitePostAdjustment` (TT 2), dalam SATU transaksi. Rencana langkah disusun `models.Rencanakan`
// (murni); berkas ini menjalankan langkah yang menyentuh basis data, kontrak Claim Prop, dan outbox:
//
//	S4         kunci klaim induk (kontrak)                      S21   retro / batas DLA / IsPrintAccept
//	S6, S26.1  tangga                                           S28   JSON_KLAIM
//	S16.5-9    nomor akseptasi                                  S29   konversi Arasapas (outbox, IsPEGAPROD)
//	S17        OS_AKSEPTASI_KLAIM                               S30-31 MONITORING_KLAIM_LOG
//	S8-S27     tulis balik klaim induk (kontrak, sekali)        S32-33 HISTORYAKSEPTASIPEGA
//	S25, S40   kepala kasus                                     S34   Kasir (outbox, IsPEGAPROD)
//	KomiteLoop tutup kasus / kembali ke KomiteRouter            S35   email (outbox, IsPEGAPROD)
//
// ⚠️ PENYIMPANGAN SADAR (PARITAS): di Pega sembilan rule menyimpan sendiri (`COMMIT` di RDB, S41 Commit). Di sini
// semuanya satu transaksi aplikasi: keputusan yang gagal di tengah batal UTUH, tidak meninggalkan tangga yang sudah
// maju tanpa tulisan klaimnya. Panggilan keluar (Arasapas, Kasir, email) diantre di transaksi yang sama dan baru
// berjalan sesudah commit (outbox) - urutan efek keluar terhadap penyimpanan sengaja diubah (keputusan 14, 20).

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimprop/backend/models"
)

// HasilKeputusan - jawaban satu Submit.
type HasilKeputusan struct {
	KasusID string `json:"kasusId"`
	// Selesai - Decision `KomiteLoop` salah: kasus Resolved-Completed.
	Selesai     bool   `json:"selesai"`
	KomiteCount int    `json:"komiteCount"`
	AcceptedNo  string `json:"acceptedNo,omitempty"`
}

// OQDokumenPDF - `PrintFileAcceptance_TKMT` S9-S13: stream `FILEAcceptanceNote` tidak diekspor; berkas PDF tidak
// dikarang (prompt §6 butir 10). Penanda `IsPrintAccept` tetap ditulis.
const OQDokumenPDF = "Dokumen Acceptance Note (PDF) menunggu stream FILEAcceptanceNote dari ekspor Pega (OQ)"

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
		if !k.Pemegang(p.AkunID) {
			return ErrBukanPemegang
		}
		if strings.TrimSpace(models.AdjustmentKlaim(kl)["AcceptedNo"]) != "" {
			return ErrAksiTertutup
		}
		if pesan := models.PeriksaIsian(k, kep); len(pesan) > 0 {
			return &GalatValidasi{Pesan: pesan}
		}
		j := &jalan{l: l, ctx: ctx, tx: tx, k: k, kl: kl, akun: p.AkunID, saat: saat}
		r := models.Rencanakan(k, kl, kep, p.AkunID, saat)
		if err := j.laksanakan(&r); err != nil {
			return err
		}
		hasil = HasilKeputusan{KasusID: k.ID, Selesai: r.Selesai, KomiteCount: r.Count, AcceptedNo: j.nomor}
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
	kl    kontrak.KlaimTreaty
	akun  string
	saat  time.Time
	nomor string
}

func (j *jalan) laksanakan(r *models.Rencana) error {
	ctx, tx, k, kl := j.ctx, j.tx, j.k, j.kl
	if r.Klaim.Adjustment == nil { // retro / Kasir menambah properti adjustment
		r.Klaim.Adjustment = map[string]string{}
	}
	for _, u := range r.Tangga { // S6, S26.1
		if err := j.l.g.TulisAnggota(ctx, tx, k.ID, u); err != nil {
			return err
		}
	}
	adj := models.AdjustmentKlaim(kl)
	if r.TerbitkanNomor { // S16.5-S16.9
		bn, err := j.l.g.UrutNomorAkseptasi(ctx, tx, j.saat)
		if err != nil {
			return err
		}
		j.nomor = models.RakitNomorAkseptasi(bn.Jenis, kl.Nilai["OfferFacIn.QuotationData.BusinessOldId"], bn.MMYYYY,
			bn.Urut)
		r.AdjustmentDiterima(j.nomor, j.saat)
	}
	noAksep := adj["AcceptedNo"]
	if j.nomor != "" {
		noAksep = j.nomor
	}
	noPolis := kl.Nilai["ClaimData.PolicyData.PolicyNo"]
	stsReject := ""
	if r.SimpanAkseptasi { // S17 SaveAcceptation_Act
		stsReject = models.StsOSAkseptasi
		if err := j.l.g.SisipOS(ctx, tx, models.SusunOSAkseptasi(kl, noAksep, k.ID, j.saat), j.saat); err != nil {
			return err
		}
	}
	if r.SimpanRetro { // S21 SaveAcceptationTreaty_TKMT
		if err := j.retro(r, adj); err != nil {
			return err
		}
		// S8 PrintFileAcceptance_TKMT: dokumen ACCEPTED CLAIM INSURANCE (FILEAcceptanceNote -> HTMLToPDF ->
		// InsertDocument_Act) dirakit saat efek dikirim (`SusunDokumenAkseptasi`).
		if err := j.antre(JenisEfekDokumen, k.ID, map[string]string{IsiPelaku: j.akun,
			IsiSaat: j.saat.Format(time.RFC3339)}); err != nil {
			return err
		}
	}
	if r.Kasir { // S34 HitServiceToKasirKMT_Act
		if err := j.kasirEfek(r, adj, noAksep); err != nil {
			return err
		}
	}
	if err := j.l.klaim.TulisBalikKlaimTreaty(ctx, tx, k.KlaimID, k.AdjustmentID, r.Klaim); err != nil { // S36
		return galatKontrak(err)
	}
	kunci := models.KunciInstans(k.KlaimID)
	if r.JSONKlaim { // S28
		if err := j.l.g.SalinJSONKlaim(ctx, tx, kunci, kl.Nilai["ClaimData.NoClaim"], noPolis, j.saat); err != nil {
			return err
		}
	}
	if r.Konversi { // S29 KonversiKlaim_Act (Connect-REST bergerbang IsPEGAPROD)
		if err := j.antre(JenisEfekKonversi, k.KlaimID, map[string]string{"CASEID": kunci, "NOPOLIS": noPolis,
			"STS_REJECT": models.StsOSAkseptasi}); err != nil {
			return err
		}
	}
	if r.LogAkseptasi { // S30-S31
		if err := j.l.g.CatatLogLayanan(ctx, tx, models.SusunLogAkseptasi(k.KlaimID, noPolis, stsReject, j.nomor), j.saat); err != nil {
			return err
		}
	}
	nama, err := j.l.a.NamaPelaku(ctx, j.akun) // S32 InsertHistory.CARI4 = OperatorID.pyUserName
	if err != nil {
		return err
	}
	if err := j.l.g.CatatRiwayatAkseptasi(ctx, tx, models.RiwayatAkseptasi{IDPega: kunci,
		IDKomite: models.KunciInstans(k.ID), Status: r.StatusRiwayat, Username: nama,
		Workbasket: models.WorkbasketRiwayat}, j.saat); err != nil { // S33
		return err
	}
	if err := j.email(r); err != nil { // S35
		return err
	}
	if err := j.l.g.SimpanKepala(ctx, tx, k.ID, k.Count, r.Kepala()); err != nil {
		return err
	}
	return j.l.g.TutupKasus(ctx, tx, k.ID, r.Selesai, j.saat) // Decision KomiteLoop
}

// retro = SaveAcceptationTreaty_TKMT S1-S9 (+ PrintFileAcceptance_TKMT S4 IsPrintAccept).
func (j *jalan) retro(r *models.Rencana, adj map[string]string) error {
	kl := j.kl
	grup := kl.Nilai["ClaimData.TreatyGroupID"]
	tahun, err := j.l.a.TahunTreaty(j.ctx, grup, models.YMD(kl.Nilai["ClaimData.PolicyData.StartDateTime"])) // S1-S3
	if err != nil {
		return err
	}
	batas := ""
	var retro []models.BarisRetro
	if reins, ada := models.SpreadingAdjustmentTerakhir(kl); ada { // S4: putaran terakhir yang tersisa di LimitDla / ListInsurerDla
		if batas, err = j.l.a.LimitPLA(j.ctx, tahun, grup, reins); err != nil {
			return err
		}
		if retro, err = j.l.a.DaftarRetro(j.ctx, tahun, grup, reins); err != nil {
			return err
		}
	}
	h, err := models.SusunRetro(adj["AdjustmentValue"], batas, retro, models.AdaFacRetro(kl)) // S5-S7
	if err != nil {
		return err
	}
	r.Klaim.FacRetro = h.FacRetro
	if h.IsFacRetro { // S7: `TempOpenPage.Message` "Please Print DLA" milik halaman klaim (transien di Claim Prop,
		// yang menampilkannya sendiri dari IsFacRetro di local action PrintFileDLA) - tidak dikembalikan ke penyetuju.
		r.Klaim.Adjustment["IsFacRetro"] = "1"
	}
	r.Klaim.Adjustment["IsPrintAccept"] = "1" // S8 (PrintFileAcceptance_TKMT S4) + S9; berkas PDF = OQDokumenPDF
	return nil
}

// kasirEfek = HitServiceToKasirKMT_Act jalur CLMP. S2: DirectToKasir dicentang dan StatusKasir kosong. S3
// `getStatusKonversi_Act`: hanya dibaca di produksi (IsPEGAPROD); di luar produksi aktivitas keluar tanpa efek, persis
// XML. S12-S13 IDOfBank, S14.1 panjang AcceptedNo, S14.1.2 email ceding, muatan S14.1.3-S14.3; REST S14.4 = outbox.
// ⚠️ Konversi S29 kini efek outbox (asinkron): di produksi status konversi saat Submit bisa belum "1" sehingga Kasir
// tidak diantre - sama dengan Claim Prop, tombol "Acceptation" (`HitServiceToKasir_Act`) mengirimnya kemudian.
func (j *jalan) kasirEfek(r *models.Rencana, adj map[string]string, noAksep string) error {
	if !(adj["DirectToKasir"] == "true" && adj["StatusKasir"] == "") { // S2
		return nil
	}
	sts, err := j.l.a.StatusKonversi(j.ctx, strings.ReplaceAll(noAksep, ".", "")) // S3 (hanya IsPEGAPROD)
	if err != nil {
		return err
	}
	if sts != "1" { // S3 transisi `.StatusKonversi=="1"`, selainnya keluar
		return nil
	}
	b := make(map[string]string, len(adj))
	for p, v := range adj {
		b[p] = v
	}
	for p, v := range r.Klaim.Adjustment {
		b[p] = v
	}
	if b["IDOfBank"] == "" { // S12-S13
		id, err := j.l.a.IDBankRekening(j.ctx, b["NameOfBank"], b["BranchOfBank"], b["NoAccount"])
		if err != nil {
			return err
		}
		if id != "" {
			r.Klaim.Adjustment["IDOfBank"] = id
			b["IDOfBank"] = id
		}
	}
	if !models.PanjangNoAksepCLMP(noAksep) { // S14.1 F=6
		return nil
	}
	email, err := j.l.a.EmailCeding(j.ctx, j.kl.Nilai["TreatyInMaster.CedingID"]) // S14.1.1-S14.1.2
	if err != nil {
		return err
	}
	m := models.SusunMuatanKasir(j.kl, b, email, false, j.l.kasir, j.akun, j.saat) // IsPEGASyariah = OQ
	return j.antre(JenisEfekKasir, noAksep, map[string]any{"TAllPaymentData": []models.MuatanKasir{m}})
}

// email = SendEmailKlaim_KMT (S35 bergerbang IsPEGAPROD; S1 tanpa SpreadingAdjustment -> keluar). MUATAN hanya
// pengenal: badan `EmailKlaim_HTML_KMT`, subjek, alamat, akun, dan CC dirakit saat dikirim (`SusunEmailKomite`).
// BCC pribadi XML (S3) tidak disalin.
func (j *jalan) email(r *models.Rencana) error {
	if !j.l.produksi || r.Email == "" || !models.AdaSpreadingAdjustment(j.kl) {
		return nil
	}
	penerima := j.k.PembuatID                    // S13-S15: pembuat kasus
	if r.Email == models.EmailPenyetujuBerikut { // S12: ComiteeClaim(KomiteCount + 1)
		penerima = ""
		for _, a := range j.k.Tangga {
			if a.Urut == j.k.Count+1 {
				penerima = a.OperatorID
			}
		}
	}
	anggota := ""
	for _, u := range r.Tangga { // S6 / S7: baris yang diputuskan (S26.1 menolak sisa tanpa komentar)
		if u.IsiKomentar {
			anggota = u.ID
			break
		}
	}
	return j.antre(JenisEfekEmailKomite, j.k.ID, map[string]string{IsiJenis: r.Email, IsiPenerima: penerima,
		IsiAnggota: anggota})
}

// Jenis efek outbox Komite Claim Prop.
const (
	JenisEfekKonversi    = "konversi-klaim"    // KonversiKlaim_Act -> KonversiKlaimNonLife (Klaim / insertClaimAccept)
	JenisEfekKasir       = "kasir"             // HitServiceToKasirKMT_Act -> SendAcceptationToKasir (Kasir / insertAllPaymentKasir)
	JenisEfekEmailKomite = "email-komite"      // SendEmailKlaim_KMT -> SendEmailWithAttachments
	JenisEfekDokumen     = "dokumen-akseptasi" // PrintFileAcceptance_TKMT -> HTMLToPDF -> InsertDocument_Act
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
