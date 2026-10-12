package services

// Tombol `Save` layar Outstanding dan gerbang tombol bawahnya - 10-10-2026 (layar Outstanding dirombak mengikuti Pega,
// permintaan work owner "kenapa tidak ada detail peserta nya seperti yang di pega?").
//
// Korpus: `kelvin/InboxClaimLife (Done)/Section/InputOSClaimLife.xml` (bNNN = nomor baris berkas itu).
//
//	b3504  .ClaimData.StsKatastrofe   radio, pyDisabledWhen CLAIM_NO !="" b3515, bawaan "Non-Catastrophe" b3728
//	b3810  .ClaimData.IsKPR           radio, tanpa syarat; pilihan KPR / Non-KPR (keputusan work owner 11-10-2026)
//	b4934  .ClaimData.KatastrofeNote  textarea, pyReadOnlyCondition CLAIM_NO !="" b5059
//	b17582 .IsAccept                  checkbox per peserta; Select All b16633 (SelectAllClaimLife_act)
//	b19413 Save                       menyimpan kasus (Obj-Save halaman kerja)
//	b21699 Send Back to Register      tampil bila CLAIM_NO ==''
//	b21832 Send to Medical Check      pyDisabledWhen CLAIM_NO =='' ; tampil bila Type QP/QR b22123
//
// Dibaca sesudah: simpanrnm.go, tahap.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

// Nilai radio Catastrophe - `[keputusan work owner]` claimprop grilling ronde 2 P12 (`StsKatastrofe = {Catastrophe,
// Non-Catastrophe}`), kode properti yang sama di Claim Life.
const (
	Katastrofe    = "Catastrophe"
	NonKatastrofe = "Non-Catastrophe"
	// Is KPR - keputusan work owner 11-10-2026: dua pilihan, teks disimpan persis (pola Catastrophe).
	KPR             = "KPR"
	NonKPR          = "Non-KPR"
	batasCatatanKat = 1000 // T_GENERAL_CLAIM.KATASTROFE_NOTE VARCHAR2(1000), migrasi 002
	// T_GENERAL_CLAIM.STS_KATASTROFE / IS_KPR VARCHAR2(20) sejak migrasi 029 (11-10-2026; dulu 8 -> ORA-12899 untuk
	// "Non-Catastrophe" 15 byte). Diperiksa SEBELUM bind supaya tidak pernah jatuh ke ORA-12899.
	batasIsKPR      = 20
	batasKatastrofe = 20
)

var (
	// ErrIsianBukanOutstanding - tombol Save layar ini hanya ada di tahap Outstanding Claim.
	ErrIsianBukanOutstanding = errors.New("services: Save layar Outstanding hanya tersedia pada tahap Outstanding Claim")
	// ErrKatastrofeTerkunci - Catastrophe dan Catastrophe Note terkunci sesudah klaim bernomor (b3515, b5059).
	ErrKatastrofeTerkunci = errors.New("services: Catastrophe dan Catastrophe Note tidak dapat diubah sesudah klaim bernomor")
	// ErrPindahMenungguNomorKlaim - Send to Medical Check menunggu Save to RNM (b21832 CLAIM_NO =='' -> nonaktif).
	ErrPindahMenungguNomorKlaim = errors.New("services: Send to Medical Check menunggu Save to RNM (klaim belum bernomor)")
	// ErrPindahTypeTanpaMedis - Send to Medical Check hanya untuk Type QP/QR (b22123).
	ErrPindahTypeTanpaMedis = errors.New("services: Send to Medical Check hanya untuk Type QP/QR")
	// ErrPindahSudahBernomor - Send Back to Register hanya sebelum klaim bernomor (b21699).
	ErrPindahSudahBernomor = errors.New("services: Send Back to Register hanya tersedia sebelum klaim bernomor")
)

// PelanggaranPindah - perpindahan ditolak karena data belum lengkap, beserta SELURUH kekurangannya (422, pola
// PelanggaranRNM Save to RNM). Error() = gabungan pesan.
type PelanggaranPindah struct {
	Langkah    string
	Pesan      string
	Kekurangan []KekuranganRNM
}

func (p *PelanggaranPindah) Error() string { return p.Pesan }

// formatDokumenMedisKurang - keputusan work owner 11-10-2026.
const formatDokumenMedisKurang = "Insured number %d (%s): at least one claim document is required because the insured is accepted."

// PeriksaDokumenAcceptMedis - proteksi Send to Medical Check (keputusan work owner 11-10-2026: "hanya bisa send to medical
// jika sudah sukses save to rnm dan semua peserta yang dicentang accept minimal punya 1 document. untuk peserta yang tidak
// dicentang accept tidak apa"), MURNI. `peserta` = peserta AKTIF klaim, urutan layar (nomor 1..n); `cacah` = jumlah dokumen
// per peserta Accept (CacahDokumenPesertaAccept). SELURUH peserta Accept tanpa dokumen dikumpulkan; tanpa peserta Accept = lolos
// (WO tidak melarang; XML SendtoMedical hanya menggerbangi CLAIM_NO b21832 dan Type b22123).
func PeriksaDokumenAcceptMedis(peserta []models.Peserta, cacah map[string]int) error {
	var kurang []KekuranganRNM
	for i, p := range peserta {
		if p.IsAccept != "true" || cacah[p.ID] > 0 {
			continue
		}
		nama := strings.TrimSpace(p.NamaTertanggung)
		if nama == "" {
			nama = strings.TrimSpace(p.NomorSertifikat)
		}
		kurang = append(kurang, KekuranganRNM{NomorPeserta: i + 1, PesertaID: p.ID, Sertifikat: p.NomorSertifikat,
			Nama: p.NamaTertanggung, Kolom: "Document", Pesan: fmt.Sprintf(formatDokumenMedisKurang, i+1, nama)})
	}
	if len(kurang) == 0 {
		return nil
	}
	baris := make([]string, len(kurang))
	for i, k := range kurang {
		baris[i] = k.Pesan
	}
	return &PelanggaranPindah{Langkah: "dokumen-medis", Pesan: strings.Join(baris, "\n"), Kekurangan: kurang}
}

// PeriksaPindahDariOutstanding - gerbang tombol bawah layar Outstanding, MURNI. Tahap asal selain Outstanding tidak
// diperiksa di sini (nil).
func PeriksaPindahDariOutstanding(asal, ke models.Tahap, nomorKlaim, tipe string) error {
	if asal != models.TahapOutstanding {
		return nil
	}
	bernomor := strings.TrimSpace(nomorKlaim) != ""
	switch ke {
	case models.TahapMedicalCheck:
		if !bernomor {
			return ErrPindahMenungguNomorKlaim
		}
		if t := strings.ToUpper(strings.TrimSpace(tipe)); t != TypeQP && t != TypeQR {
			return fmt.Errorf("%w: Type %q", ErrPindahTypeTanpaMedis, tipe)
		}
	case models.TahapInputRegister:
		if bernomor {
			return ErrPindahSudahBernomor
		}
	}
	return nil
}

// IsianOutstanding adalah isi tombol Save layar Outstanding.
type IsianOutstanding struct {
	StsKatastrofe  string `json:"stsKatastrofe"`
	IsKPR          string `json:"isKpr"`
	KatastrofeNote string `json:"katastrofeNote"`
	// Accept - pengenal peserta -> "true" / "false" (`.IsAccept` b17582). Peserta yang tidak disebut tidak disentuh.
	Accept map[string]string `json:"accept"`
	// CatatanAdmin - pengenal peserta -> ADMIN NOTES (ClaimLifeDetailGCNM b2588, migrasi 027; 11-10-2026). Peserta yang
	// tidak disebut tidak disentuh.
	CatatanAdmin map[string]string `json:"catatanAdmin"`
	// Tanggal - pengenal peserta -> kelima tanggal yang diisi pemakai sebelum Save to RNM (keputusan work owner 11-10-2026).
	// Peserta yang tidak disebut tidak disentuh; medan kosong = dikosongkan.
	Tanggal map[string]IsianTanggal `json:"tanggal"`
	// Adjustment - pengenal baris adjustment -> Claim RNM + data bank (dialog AdjustmentDetail). Dilebur ke tombol Save grid
	// 11-10-2026 (keputusan work owner "tombol save to outstanding juga hilangkan saja"): SATU jalur simpan.
	Adjustment map[string]IsianBarisAdjustment `json:"adjustment"`
}

// IsianTanggal - kelima tanggal satu peserta, teks `YYYY-MM-DD` ("" = kosong).
type IsianTanggal struct {
	Kejadian       string `json:"kejadian"`
	TerimaKlaim    string `json:"terimaKlaim"`
	DokumenLengkap string `json:"dokumenLengkap"`
	MaksData       string `json:"maksData"`
	Konfirmasi     string `json:"konfirmasi"`
}

// uraiIsianTanggal mengurai satu isian; galat = 400.
func uraiIsianTanggal(id string, t IsianTanggal) (repository.TanggalPeserta, error) {
	var out repository.TanggalPeserta
	for _, x := range []struct {
		nama string
		teks string
		ke   **time.Time
	}{
		{"Date of Loss", t.Kejadian, &out.Kejadian}, {"Claim Received Date", t.TerimaKlaim, &out.TerimaKlaim},
		{"Document Complete Date", t.DokumenLengkap, &out.DokumenLengkap}, {"Max Data Received", t.MaksData, &out.MaksData},
		{"Confirmation Date", t.Konfirmasi, &out.Konfirmasi},
	} {
		v := strings.TrimSpace(x.teks)
		if v == "" {
			continue
		}
		d, err := time.Parse("2006-01-02", v)
		if err != nil {
			return out, fmt.Errorf("%w: %s peserta %q bukan tanggal YYYY-MM-DD", galat.ErrPermintaanTidakSah, x.nama, id)
		}
		*x.ke = &d
	}
	return out, nil
}

// ErrCatatanAdminTerkunci - ADMIN NOTES dibekukan bila peserta sudah diputus (b2628).
var ErrCatatanAdminTerkunci = errors.New("services: Admin Notes tidak dapat diubah untuk peserta yang sudah diputus")

// PeriksaIsianOutstanding - bentuk dan kunci isian, MURNI. `lama` = header tersimpan; `pesertaAktif` = pengenal
// peserta aktif klaim itu.
func PeriksaIsianOutstanding(is IsianOutstanding, lama models.Klaim, pesertaAktif map[string]bool) error {
	if len(is.StsKatastrofe) > batasKatastrofe {
		return fmt.Errorf("%w: Catastrophe melebihi %d byte", galat.ErrPermintaanTidakSah, batasKatastrofe)
	}
	switch is.StsKatastrofe {
	case "", Katastrofe, NonKatastrofe:
	default:
		return fmt.Errorf("%w: Catastrophe %q di luar {%s, %s}", galat.ErrPermintaanTidakSah, is.StsKatastrofe,
			Katastrofe, NonKatastrofe)
	}
	if len(is.KatastrofeNote) > batasCatatanKat {
		return fmt.Errorf("%w: Catastrophe Note melebihi %d byte", galat.ErrPermintaanTidakSah, batasCatatanKat)
	}
	if len(is.IsKPR) > batasIsKPR {
		return fmt.Errorf("%w: Is KPR melebihi %d byte", galat.ErrPermintaanTidakSah, batasIsKPR)
	}
	// Di luar {KPR, Non-KPR} ditolak (keputusan work owner 11-10-2026) - KECUALI nilai yang SAMA dengan yang tersimpan:
	// data warisan (mis. Yes/Y/1) ditampilkan apa adanya, dan Save yang tidak menyentuhnya tidak boleh gagal karenanya.
	switch is.IsKPR {
	case "", KPR, NonKPR, lama.IsKPR:
	default:
		return fmt.Errorf("%w: Is KPR %q di luar {%s, %s}", galat.ErrPermintaanTidakSah, is.IsKPR, KPR, NonKPR)
	}
	for id, v := range is.Accept {
		if v != models.PenandaDipilih && v != "false" {
			return fmt.Errorf("%w: Accept peserta %q bernilai %q, mau true/false", galat.ErrPermintaanTidakSah, id, v)
		}
		if !pesertaAktif[id] {
			return fmt.Errorf("%w: peserta %q bukan peserta aktif klaim %q", galat.ErrPermintaanTidakSah, id, lama.ID)
		}
	}
	for id, v := range is.CatatanAdmin {
		if len(v) > batasCatatanKat {
			return fmt.Errorf("%w: Admin Notes peserta %q melebihi %d byte", galat.ErrPermintaanTidakSah, id, batasCatatanKat)
		}
		if !pesertaAktif[id] {
			return fmt.Errorf("%w: peserta %q bukan peserta aktif klaim %q", galat.ErrPermintaanTidakSah, id, lama.ID)
		}
	}
	if strings.TrimSpace(lama.NomorKlaim) != "" &&
		(is.StsKatastrofe != lama.StsKatastrofe || is.KatastrofeNote != lama.KatastrofeNote) {
		return ErrKatastrofeTerkunci
	}
	return nil
}

// SimpanOutstanding adalah layanan tombol Save layar Outstanding.
type SimpanOutstanding struct{ svc *Service }

// SimpanOutstanding menyusun layanannya.
func (s *Service) SimpanOutstanding() *SimpanOutstanding { return &SimpanOutstanding{svc: s} }

// Simpan menulis isian layar Outstanding dalam SATU transaksi.
//
// Gerbang: identitas, pengenal, kasus terbuka, tahap Outstanding + pemegangnya (Admin), bentuk isian, lalu tulis.
// ⛔ Status baris, nomor, dan tahap TIDAK disentuh - ini bukan Save to RNM.
func (x *SimpanOutstanding) Simpan(ctx context.Context, pelaku inti.Pelaku, klaimID string, is IsianOutstanding,
	_ time.Time) error {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if strings.TrimSpace(klaimID) == "" {
		return fmt.Errorf("%w: pengenal klaim wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if x == nil || x.svc == nil || !x.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
	}
	if err := x.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		if errors.Is(err, repository.ErrWorkTidakAda) {
			return fmt.Errorf("%w: %w", ErrKlaimTidakAda, err)
		}
		return err
	}
	baca := repository.NewKlaimLife(x.svc.DB())
	tahap, err := tahapKasus(ctx, baca, klaimID)
	if err != nil {
		return err
	}
	if tahap != models.TahapOutstanding {
		return fmt.Errorf("%w: tahap %s", ErrIsianBukanOutstanding, tahap)
	}
	if err := WajibPemegangTahap(pelaku, tahap); err != nil {
		return err
	}
	lama, err := baca.AmbilHeader(ctx, klaimID)
	if err != nil {
		return err
	}
	if lama == nil {
		return ErrKlaimTidakAda
	}
	peserta, err := baca.AmbilPeserta(ctx, klaimID)
	if err != nil {
		return err
	}
	aktif := map[string]bool{}
	for _, p := range peserta {
		aktif[p.ID] = true
	}
	if err := PeriksaIsianOutstanding(is, *lama, aktif); err != nil {
		return err
	}
	bernomor := strings.TrimSpace(lama.NomorKlaim) != ""
	// Kelima tanggal (11-10-2026): hanya sebelum Save to RNM (OQ-M1, BolehUbahTanggalKlaim) dan klaim belum bernomor;
	// DOL tetap melewati ValidasiDOL (jendela valuasi, ValidasiDOL_Act) atas peserta TERSIMPAN.
	tanggal := map[string]repository.TanggalPeserta{}
	if len(is.Tanggal) > 0 {
		sudah, err := baca.SudahSaveRNM(ctx, klaimID)
		if err != nil {
			return err
		}
		if bernomor || !models.BolehUbahTanggalKlaim(tahap, sudah) {
			return ErrTanggalTerkunciSesudahSaveRNM
		}
		for id, t := range is.Tanggal {
			if !aktif[id] {
				return fmt.Errorf("%w: peserta %q bukan peserta aktif klaim %q", galat.ErrPermintaanTidakSah, id, klaimID)
			}
			tp, err := uraiIsianTanggal(id, t)
			if err != nil {
				return err
			}
			tanggal[id] = tp
		}
		// ⛔ DOL di luar periode valuation TIDAK lagi menolak Save grid (keputusan work owner 11-10-2026: Save = simpan
		// progres). Jendelanya ditegakkan Save to RNM (PeriksaSimpanRNM -> DOLDalamJendela) dengan kalimat yang sama.
	}
	// Claim RNM + bank (11-10-2026): hanya sebelum Save to RNM (409), Claim RNM <= Share RNM (422).
	adjustment := map[string]repository.IsianBaris{}
	if len(is.Adjustment) > 0 {
		if bernomor {
			return ErrAdjustmentTerkunci
		}
		perBaris, err := baca.AmbilBaris(ctx, klaimID)
		if err != nil {
			return err
		}
		for adjID, isi := range is.Adjustment {
			var baris *models.BarisAdjustment
			for pid, bs := range perBaris {
				for i := range bs {
					if bs[i].ID == adjID && aktif[pid] {
						baris = &bs[i]
					}
				}
			}
			if baris == nil {
				return fmt.Errorf("%w: baris %q bukan milik klaim %q", galat.ErrPermintaanTidakSah, adjID, klaimID)
			}
			hasil, err := PeriksaIsianBaris(*baris, isi)
			if err != nil {
				return err
			}
			adjustment[adjID] = hasil
		}
	}
	return x.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		ok, err := baca.SetelMedanRegisterClaim(ctx, tx, klaimID, bernomor, is.StsKatastrofe, is.IsKPR, is.KatastrofeNote)
		if err != nil {
			return err
		}
		if !ok {
			// Save to RNM menyela di antara pembacaan dan penulisan: klaim kini bernomor.
			return ErrKatastrofeTerkunci
		}
		for id, v := range is.Accept {
			if err := baca.SetelAcceptPeserta(ctx, tx, klaimID, id, v); err != nil {
				return err
			}
		}
		for adjID, isi := range adjustment {
			ok, err := baca.SetelClaimRNM(ctx, tx, klaimID, adjID, isi)
			if err != nil {
				return err
			}
			if !ok {
				return ErrAdjustmentTerkunci
			}
		}
		for id, tp := range tanggal {
			ok, err := baca.SetelTanggalPeserta(ctx, tx, klaimID, id, tp)
			if err != nil {
				return err
			}
			if !ok {
				return ErrTanggalTerkunciSesudahSaveRNM
			}
		}
		for id, v := range is.CatatanAdmin {
			ok, err := baca.SetelCatatanAdmin(ctx, tx, klaimID, id, v)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("%w: peserta %q", ErrCatatanAdminTerkunci, id)
			}
		}
		return nil
	})
}

// PesertaPolis - View Detail tab Policy Details (keputusan work owner 11-10-2026): satu halaman peserta versi polis
// yang SAMA dengan header klaim (PL_NUMBER + IDPEGA peserta pertama yang ber-IDPEGA; cadangan lewat SOURCE_ID seperti
// KlaimLife.Ambil). Baca saja.
func (x *SimpanOutstanding) PesertaPolis(ctx context.Context, klaimID string, halaman int) (repository.HalamanPesertaPolis, error) {
	if strings.TrimSpace(klaimID) == "" {
		return repository.HalamanPesertaPolis{}, fmt.Errorf("%w: pengenal klaim wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if x == nil || x.svc == nil || !x.svc.PunyaDatabase() {
		return repository.HalamanPesertaPolis{}, db.ErrTanpaOracle
	}
	baca := repository.NewKlaimLife(x.svc.DB())
	hdr, err := baca.AmbilHeader(ctx, klaimID)
	if err != nil {
		return repository.HalamanPesertaPolis{}, err
	}
	if hdr == nil {
		return repository.HalamanPesertaPolis{}, ErrKlaimTidakAda
	}
	peserta, err := baca.AmbilPeserta(ctx, klaimID)
	if err != nil {
		return repository.HalamanPesertaPolis{}, err
	}
	pl, idPega := "", ""
	for _, p := range peserta {
		id := strings.TrimSpace(p.IdPega)
		if id == "" {
			if id, err = baca.IdPegaSumber(ctx, repository.KunciPesertaSumber{PLNumber: p.NomorPremiList,
				Sertifikat: p.NomorSertifikat, SumberID: p.SumberID}); err != nil {
				return repository.HalamanPesertaPolis{}, err
			}
		}
		if strings.TrimSpace(id) != "" {
			pl, idPega = p.NomorPremiList, id
			break
		}
	}
	return baca.PesertaPolisDetail(ctx, pl, idPega, halaman)
}

// IsianBarisAdjustment - isi dialog AdjustmentDetail yang dapat diisi sebelum Save to RNM (11-10-2026).
//
// ⛔ Keputusan work owner 11-10-2026 "semua read only kecuali name of bank": Branch of Bank dan Account No TIDAK lagi bagian
// isian. Kunci `cabangBank` / `nomorRekening` dalam badan permintaan DIABAIKAN (tidak ada medannya, json.Decoder membuang
// kunci asing) dan kolomnya tidak ikut UPDATE (repository.sqlSetelClaimRNM) - nilai tersimpan tidak tertimpa.
type IsianBarisAdjustment struct {
	Jumlah   string `json:"jumlah"`
	NamaBank string `json:"namaBank"`
}

// ErrAdjustmentTerkunci - Claim RNM hanya dapat diisi sebelum Save to RNM (keputusan work owner 11-10-2026: "saat belum
// di save to rnm (belum generate no claim, dan belum ber status outstanding) data adjustmen masih bisa di edit, setelah
// sudah save to rnm tidak bisa lagi").
var ErrAdjustmentTerkunci = errors.New("services: adjustment tidak dapat diubah sesudah Save to RNM")

// PeriksaIsianBaris - Claim RNM dan Name of Bank satu baris (dialog AdjustmentDetail), MURNI.
//
// Save grid = SIMPAN PROGRES (keputusan work owner 11-10-2026: "hanya save progress saja, tidak ada menjalankan fungsi
// apapun"). Yang tetap MENOLAK karena tidak dapat disimpan: baris sudah berstatus (409 ErrAdjustmentTerkunci), Claim RNM
// bukan desimal >= 0 (400), Name of Bank melebihi lebar kolom (400). Claim RNM > Share Nusantara Re TIDAK lagi ditolak di
// sini - aturan bisnis itu ditegakkan Save to RNM (PeriksaSimpanRNM, kalimat yang sama dengan layar). Kosong = NULL.
func PeriksaIsianBaris(baris models.BarisAdjustment, is IsianBarisAdjustment) (repository.IsianBaris, error) {
	if strings.TrimSpace(baris.KodeStatus) != "" {
		return repository.IsianBaris{}, ErrAdjustmentTerkunci
	}
	if err := periksaLebarBank(is); err != nil {
		return repository.IsianBaris{}, err
	}
	hasil := repository.IsianBaris{NamaBank: strings.TrimSpace(is.NamaBank)}
	if strings.TrimSpace(is.Jumlah) == "" {
		return hasil, nil
	}
	d, err := utils.ParseDecimal(strings.TrimSpace(is.Jumlah))
	if err != nil || d.Negative {
		return repository.IsianBaris{}, fmt.Errorf("%w: Claim RNM %q bukan angka >= 0", galat.ErrPermintaanTidakSah, is.Jumlah)
	}
	hasil.Jumlah = utils.FormatDecimal(d)
	return hasil, nil
}

// periksaLebarBank - Name of Bank dalam lebar kolom 004 (400 sebelum bind).
func periksaLebarBank(is IsianBarisAdjustment) error {
	for _, x := range []struct {
		nama  string
		nilai string
		batas int
	}{{"Name of Bank", is.NamaBank, 255}} {
		if len(strings.TrimSpace(x.nilai)) > x.batas {
			return fmt.Errorf("%w: %s melebihi %d byte", galat.ErrPermintaanTidakSah, x.nama, x.batas)
		}
	}
	return nil
}
