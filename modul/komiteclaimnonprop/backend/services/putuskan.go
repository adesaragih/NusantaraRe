package services

// Untuk apa berkas ini: SUBMIT - tombol "Submit" (`finishAssignment`) flow action `ViewTransferDtl`, pasca-proses
// `KomitePostAdjustment` (korpus `Komite Claim Non Prop`), dalam SATU transaksi. Rencana langkah disusun
// `models.Rencanakan` (murni); berkas ini menjalankan langkah yang menyentuh basis data, kontrak Claim Non Prop, dan
// outbox:
//
//	S5           kunci klaim induk (kontrak)                  S16-S18  IsKomite / IsSubjectivity (tulis balik)
//	S6-S7, S11.3 tangga                                       S17      OS_AKSEPTASI_SUBJECTIVITY
//	S14.8-S14.11 nomor akseptasi                              S19.2    email (outbox, IsPEGAPROD)
//	S14.13       AcceptedNo / AcceptedDate / AcceptanceStatus S19.3    Kasir (outbox, IsPEGAPROD)
//	S14.18       OS_AKSEPTASI_KLAIM + CLAIMXOL2               S22-S23  HISTORYAKSEPTASIPEGA
//	S14.22       konversi Arasapas (outbox, IsPEGAPROD)       S24      JSON_KLAIM
//	S8-S20       tulis balik klaim induk (kontrak, sekali)    S20/S29  kepala kasus; S28 tutup kasus
//
// PENYIMPANGAN SADAR (PARITAS, pola Komite Claim Prop): di Pega setiap rule menyimpan sendiri (`COMMIT` di RDB, S27
// Commit). Di sini semuanya satu transaksi aplikasi: keputusan yang gagal di tengah batal UTUH. Panggilan keluar
// (Arasapas, Kasir, email) diantre di transaksi yang sama dan baru berjalan sesudah commit (outbox).

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimnonprop/backend/models"
)

// HasilKeputusan - jawaban satu Submit.
type HasilKeputusan struct {
	KasusID string `json:"kasusId"`
	// Selesai - S28 `KomiteCount >= KomiteLoop`: kasus Resolved-Completed.
	Selesai     bool   `json:"selesai"`
	KomiteCount int    `json:"komiteCount"`
	AcceptedNo  string `json:"acceptedNo,omitempty"`
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
		if pesan := models.PeriksaIsian(k, kep); len(pesan) > 0 {
			return &GalatValidasi{Pesan: pesan}
		}
		nama, err := l.a.NamaPelaku(ctx, p.AkunID) // S8-S10 / S22 OperatorID.pyUserName
		if err != nil {
			return err
		}
		j := &jalan{l: l, ctx: ctx, tx: tx, k: k, kl: kl, akun: p.AkunID, nama: nama, saat: saat}
		r := models.Rencanakan(k, kl, kep, p.AkunID, nama, saat)
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
	nama  string
	saat  time.Time
	nomor string
}

func (j *jalan) laksanakan(r *models.Rencana) error {
	ctx, tx, k, kl := j.ctx, j.tx, j.k, j.kl
	for _, u := range r.Tangga { // S6 / S7, S11.3
		if err := j.l.g.TulisAnggota(ctx, tx, k.ID, u); err != nil {
			return err
		}
	}
	adj := models.AdjustmentKlaim(kl)
	noAksep := adj["AcceptedNo"]
	if r.TerbitkanNomor { // S14.8-S14.11
		bn, err := j.l.g.UrutNomorAkseptasi(ctx, tx, j.saat)
		if err != nil {
			return err
		}
		j.nomor = models.RakitNomorAkseptasi(bn.Jenis, kl.Nilai["OfferFacIn.QuotationData.BusinessOldId"], bn.MMYYYY,
			bn.Urut)
		noAksep = j.nomor
	}
	kunci := models.KunciInstans(k.KlaimID)
	noPolis := kl.Nilai["ClaimData.PolicyData.PolicyNo"]
	if r.SimpanOS { // S14.13 lalu S14.18 InsertOSKlaimCNP (S6 OS, S7 InsertXOLKlaimCNP)
		r.AdjustmentDiterima(noAksep, j.saat)
		b, err := models.SusunOSAkseptasi(kl, k.KlaimID, noAksep, j.saat)
		if err != nil {
			return err
		}
		if err := j.l.g.SisipOS(ctx, tx, b, j.saat); err != nil {
			return err
		}
		for _, x := range models.SusunXOL2(kl, k.KlaimID, r.Subjectivity) {
			if err := j.l.g.SisipXOL2(ctx, tx, x, j.saat); err != nil {
				return err
			}
		}
	}
	if r.Konversi { // S14.22 KonversiKlaim_Act (Connect-REST bergerbang IsPEGAPROD; STSREJECT "1" tertulis mati)
		if err := j.antre(JenisEfekKonversi, k.KlaimID, map[string]string{"CASEID": kunci, "NOPOLIS": noPolis,
			"STS_REJECT": models.StsRejectKonversi}); err != nil {
			return err
		}
	}
	if r.SimpanOSSubjectivity { // S17 InsertOSSubjectivityCNP
		b, err := models.SusunOSSubjectivity(kl, k.KlaimID, noAksep, r.Subjectivity, j.saat)
		if err != nil {
			return err
		}
		if err := j.l.g.SimpanOSSubjectivity(ctx, tx, b, j.saat); err != nil {
			return err
		}
	}
	if err := j.email(r); err != nil { // S19.2
		return err
	}
	if r.Kasir { // S19.3 HitServiceToKasirKMT_Act
		if err := j.kasirEfek(r, adj, noAksep); err != nil {
			return err
		}
	}
	if err := j.l.klaim.TulisBalikKlaimTreaty(ctx, tx, k.KlaimID, k.AdjustmentID, r.Klaim); err != nil {
		return galatKontrak(err)
	}
	if err := j.l.g.CatatRiwayatAkseptasi(ctx, tx, models.RiwayatAkseptasi{IDPega: kunci,
		IDKomite: models.KunciInstans(k.ID), Status: r.StatusRiwayat, Username: j.nama,
		Workbasket: models.WorkbasketRiwayat}, j.saat); err != nil { // S22-S23
		return err
	}
	// S24 (dan S14.21 / InsertOSKlaimCNP S8 / InsertOSSubjectivityCNP S5 - INSERT hanya sekali per IDPEGA).
	if err := j.l.g.SalinJSONKlaim(ctx, tx, kunci, kl.Nilai["ClaimData.NoClaim"], noPolis, j.saat); err != nil {
		return err
	}
	if err := j.l.g.SimpanKepala(ctx, tx, k.ID, k.Count, r.Kepala()); err != nil { // S20 / S29
		return err
	}
	return j.l.g.TutupKasus(ctx, tx, k.ID, r.Selesai, r.Posisi, j.saat) // S28; POSITION tingkat berikut
}

// kasirEfek = HitServiceToKasirKMT_Act jalur CLMNP. S2: DirectToKasir dicentang dan StatusKasir kosong. S3
// `getStatusKonversi_Act`: hanya dibaca di produksi (IsPEGAPROD). S9-S10 muatan per mata uang Spreading In (hanya bila
// panjang AcceptedNo 23 / 24), S10.2 email ceding, REST S10.7 = outbox. S7 NoAccount angka saja, S12-S13 IDOfBank. S10.8 StatusKasir dan S11
// IsPrintAccept bergantung jawaban REST - tidak ditulis di Submit (pola Komite Claim Prop).
func (j *jalan) kasirEfek(r *models.Rencana, adj map[string]string, noAksep string) error {
	if !(adj["DirectToKasir"] == "true" && adj["StatusKasir"] == "") { // S2
		return nil
	}
	sts, err := j.l.a.StatusKonversi(j.ctx, strings.ReplaceAll(noAksep, ".", "")) // S3 (hanya IsPEGAPROD)
	if err != nil {
		return err
	}
	if sts != models.StatusKonversiSudah { // S3: belum dikonversi -> keluar
		return nil
	}
	b := make(map[string]string, len(adj))
	for p, v := range adj {
		b[p] = v
	}
	for p, v := range r.Klaim.Adjustment {
		b[p] = v
	}
	b["AcceptedNo"] = noAksep
	if akun := models.RekeningAngka(b["NoAccount"]); akun != b["NoAccount"] { // S7 (tersimpan S26)
		b["NoAccount"] = akun
		r.Klaim.Adjustment = tambahUbahan(r.Klaim.Adjustment, "NoAccount", akun)
	}
	if models.PanjangNoAksepCNP(noAksep) { // S9-S10
		email, err := j.l.a.EmailCeding(j.ctx, j.kl.Nilai["TreatyInMaster.CedingID"]) // S10.1-S10.2
		if err != nil {
			return err
		}
		m, err := models.SusunMuatanKasir(j.kl, b, email, false, j.l.kasir, j.akun) // IsPEGASyariah = OQ-CNP-40
		if err != nil {
			return err
		}
		if len(m) > 0 {
			if err := j.antre(JenisEfekKasir, noAksep, map[string]any{"TAllPaymentData": m}); err != nil {
				return err
			}
		}
	}
	if b["IDOfBank"] == "" { // S12-S13
		id, err := j.l.a.IDBankRekening(j.ctx, b["NameOfBank"], b["BranchOfBank"], b["NoAccount"])
		if err != nil {
			return err
		}
		if id != "" {
			r.Klaim.Adjustment = tambahUbahan(r.Klaim.Adjustment, "IDOfBank", id)
		}
	}
	return nil
}

// tambahUbahan - satu properti akseptasi ditambahkan ke ubahan (peta nil dibuat).
func tambahUbahan(m map[string]string, prop, nilai string) map[string]string {
	if m == nil {
		m = map[string]string{}
	}
	m[prop] = nilai
	return m
}

// email = SendEmailKlaim_KMT (S19 bergerbang IsPEGAPROD; tanpa SpreadingAdjustment -> keluar). MUATAN hanya pengenal:
// badan `EmailKlaim_HTML_KMT`, subjek, alamat, akun, dan CC dirakit saat dikirim (`SusunEmailKomite`). BCC pribadi XML
// tidak disalin.
func (j *jalan) email(r *models.Rencana) error {
	if !j.l.produksi || r.Email == "" || !models.AdaSpreadingAdjustment(j.kl) {
		return nil
	}
	penerima := j.k.PembuatID                    // pembuat kasus
	if r.Email == models.EmailPenyetujuBerikut { // ComiteeClaim(KomiteCount + 1)
		penerima = ""
		for _, a := range j.k.Tangga {
			if a.Urut == j.k.Count+1 {
				penerima = a.OperatorID
			}
		}
	}
	anggota := ""
	for _, u := range r.Tangga { // S6 / S7: baris yang diputuskan (pertama di rencana)
		if u.Pemutus != "" {
			anggota = u.ID
			break
		}
	}
	return j.antre(JenisEfekEmailKomite, j.k.ID, map[string]string{IsiJenis: r.Email, IsiPenerima: penerima,
		IsiAnggota: anggota})
}

// Jenis efek outbox Komite Claim Non Prop.
const (
	JenisEfekKonversi    = "konversi-klaim" // KonversiKlaim_Act -> KonversiKlaimNonLife (Klaim / insertClaimAccept)
	JenisEfekKasir       = "kasir"          // HitServiceToKasirKMT_Act -> SendAcceptationToKasir (Kasir / insertAllPaymentKasir)
	JenisEfekEmailKomite = "email-komite"   // SendEmailKlaim_KMT -> SendEmailWithAttachments
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
