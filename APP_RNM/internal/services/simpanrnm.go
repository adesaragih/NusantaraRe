package services

// `Save to RNM` - tiket 03, GILIRAN-11 paket 1.
//
// Untuk apa berkas ini: tombol `Save to RNM` layar Outstanding
// (`Section/InputOSClaimLife.xml` b21102 -> `refresh` b21112 ->
// `SaveOutStandingLife_Act` b21126), ditiru dari `Activity/
// SaveOutStandingLife_Act.xml` yang dibaca UTUH sebagai pohon (29 langkah
// teratas; pecahan `sed 's/></>\n</g'`, bNNN).
//
// ⛔ CARA MEMBACA - ditetapkan dari korpus: baris `WHEN` hanya berlaku bila
// langkahnya ber-`pyStepsPreCondition=true`; baris `TRANS` hanya bila
// `pyStepsTransition=true`. Kode aksi: 2 lanjut, 3 lewati langkah, 6 keluar
// activity. Langkah 11.1, 22.1.1, dan 22.1.3 ber-WHEN tetapi TANPA bendera -
// WHEN-nya mati, dan langkahnya selalu berjalan.
//
// ⛔ `pyStepsBlockName = //` berarti langkahnya DI-REMARK dan tidak pernah
// jalan, apa pun gerbangnya (`.scratch/claim-prop/grilling-ronde-2.md` Aturan
// 2). Delapan langkah activity ini ber-remark: 11.3 b3495, 11.9 b4632, 11.11
// b5009, 12 b6178 (beserta anaknya), 13 b7074, 14 b7293, 15 b7512, 23 b10649.
// RALAT GILIRAN-11 (temuan /code-review): paket 1 membaca 11.9/11.11 (STNC)
// dan 12 (dokumen lengkap) sebagai gerbang hidup - pembaca pohon kami tidak
// mencetak medan itu - sehingga dua penolakan yang tidak pernah terjadi di
// sistem lama ikut ditegakkan. Keduanya dibuang.
//
// Peta langkah -> di sini:
//
//	1       Page-Clear-Messages                           -
//	2-4     dokumen ADA (bukan TP/TR, peserta IsAccept)    gerbang 1, keluar b1249
//	5       `.Protect` memuat "1" -> pesan b1614           RESIDU: `.Protect` NOL penulis di korpus
//	6       Page-Remove halaman kerja                       -
//	7-9     periode penomoran                              milik penomor (tiket 02)
//	10      Business per BusinessCode                      ContentNoteDari
//	11.1-6  klaim ganda DEATH / HEALTH (tabel warisan)      gerbang 2 - BACA SAJA; 11.3 `//`
//	11.7-10 DOL dalam jendela, TANPA geser retro            gerbang 3
//	11.9/11 STNC                                           `//` - TIDAK ditiru
//	11.12-17 DOB, BEGIN, EXPIRED, POLICY, CERT, CLAIM_GROSS gerbang 4
//	12      dokumen LENGKAP                                 `//` - TIDAK ditiru
//	13-15   nomor lewat `Generate_NoKlaim_Life*`            `//` - TIDAK ditiru
//	16-20   nomor klaim bila CLAIM_NO kosong                penomor yang sama dengan pendaftaran
//	21, 24  bendera `pyWorkPage.Save`                       ⚠️ TANPA kolom - OQ-N1
//	22      insert baris warisan + STS_REJECT=0             warisan BACA SAJA (brief); STS_REJECT
//	                                                        hanya baris tanpa status (tiket 04)
//	23      enam total peserta                              `//` - totalnya dari SavePesertaClaim
//	                                                        langkah 8 (HitungTotalPeserta, saat baca)
//	25      InsertJsonClaimLife_Act: 1-2 salin + TUKAR      tukarnya dibaca langkah 27 - PolisRetro
//	        retro; 3-4 JSON                                 JSON DIBUANG (keputusan 2026-09-16)
//	26      transformasi tampilan JSON                      DIBUANG (keputusan 2026-09-16)
//	27      keluar bila RETROID/SECURITYREINSURERID tertentu ArasapasDilewatiRetro
//	28      Arasapas bila IsPEGAPROD                        penyalur, SESUDAH commit (ADR-U-0008)
//	29      Obj-Save                                        transaksi tunggal
//
// Dibaca sesudah: dokumen.go, dol.go, statusbaris.go.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
	"nusantarare/inti/jejak"
	"nusantarare/inti/layanan"
	"nusantarare/inti/outbox"
	"nusantarare/inti/utils"
)

// ErrSimpanRNMDitolak - salah satu gerbang XML menolak; kalimatnya di
// PelanggaranRNM.Pesan.
var ErrSimpanRNMDitolak = errors.New("services: Save to RNM ditolak")

// PelanggaranRNM adalah alasan Save to RNM berhenti - kalimat XML VERBATIM.
//
// ⛔ `Error()` mengembalikan kalimatnya saja, tanpa awalan: pemakai sistem
// lama mencari kalimat yang sama, dan pesan galat adalah logika bisnis.
type PelanggaranRNM struct {
	// Langkah XML asal pesannya, untuk yang menelusuri.
	Langkah string
	Pesan   string
}

func (p *PelanggaranRNM) Error() string { return p.Pesan }

// Unwrap membuat `errors.Is(err, ErrSimpanRNMDitolak)` benar.
func (p *PelanggaranRNM) Unwrap() error { return ErrSimpanRNMDitolak }

// PesertaRNM adalah seorang peserta beserta tiga fakta yang hanya dapat dibaca
// dari tabel sumber dan tabel warisan.
//
// ⛔ Nama tertanggung dan tanggal lahir TIDAK dibawa: yang menyeberang hanya
// JAWABANNYA (kosong? ganda?). Pencocokannya terjadi di SQL - lihat
// repository/gandawarisan.go.
type PesertaRNM struct {
	models.Peserta
	// DOBKosong - `M_LIFE_PREMIUM_DETAIL.DOB` peserta ini NULL (langkah 11.12).
	DOBKosong bool
	// StatusWarisanTerakhir - `STS_REJECT` baris warisan TERBARU bertertanggung
	// sama (`CountPesertaAkseptasiLife_SQL`, ORDER BY ACCEPTATION_DATE DESC);
	// kosong bila tidak ada. Dipakai bila ContentNote DEATH (langkah 11.4).
	StatusWarisanTerakhir string
	// AdaKlaimSehatSamaDOL - ada baris warisan bertertanggung sama DAN
	// `LAPSE_DATE` = DOL (`CountPesertaAkseptasiLifeHealth_SQL`, langkah 11.6).
	AdaKlaimSehatSamaDOL bool
}

// MasukanRNM adalah seluruh masukan aturan murni Save to RNM.
type MasukanRNM struct {
	// Tipe - `PolicyDataLife.Type`.
	Tipe string
	// ContentNote - `Business.pxResults(1).ContentNote` (langkah 10).
	ContentNote string
	Peserta     []PesertaRNM
}

// Kalimat VERBATIM, termasuk salah eja dan apostrof lengkungnya.
const (
	// b1054 `local.Errmsg10`, bertumpuk per peserta.
	pesanDokumenBelumDiunggah = "The document hasn’t been uploaded person number "
	// b2974 `local.Errmsg7`.
	formatSudahDiaksep = "Person number %d has already been accepted."
	// b2827 `Local.Errmsg`.
	formatDOLDiLuar = "DOL cannot be blank or outside the valuation period No %d"
	// ⚠️ b2996 `local.Errmsg8` ("Begin date exceed STNC") disusun langkah 11
	// tetapi hanya ditampilkan 11.11, yang ter-remark - karena itu tidak ada.
	// b2869, b2890, b2911, b2932, b2953, b2848.
	formatDOBKosong     = "DOB cannnot be blank No %d"
	formatBeginKosong   = "Begin Date cannnot be blank No %d"
	formatExpiredKosong = "Expired Date cannnot be blank No %d"
	formatPolisKosong   = "Policy No cannnot be blank No %d, please contact IT"
	formatSertKosong    = "Certificate No cannnot be blank No %d, please contact IT"
	formatGrossKosong   = "Claim Gross No %d can't null"
)

// contentNoteDeath - langkah 11.2-11.4 `ContentNote = "DEATH"`; 11.5-11.6 `!=`.
const contentNoteDeath = "DEATH"

func tolakRNM(langkah, pesan string) error { return &PelanggaranRNM{Langkah: langkah, Pesan: pesan} }

func kosongTeks(s string) bool { return strings.TrimSpace(s) == "" }

// PeriksaSimpanRNM menjalankan seluruh gerbang Save to RNM, URUT XML - MURNI.
//
// Tiap gerbang yang gagal KELUAR (kode 6), jadi hanya pelanggaran pertama yang
// dilaporkan - kecuali gerbang 1, yang menumpuk pesannya untuk seluruh peserta
// lebih dulu (langkah 3 berputar, langkah 4 keluar).
func PeriksaSimpanRNM(m MasukanRNM) error {
	peserta := make([]models.Peserta, len(m.Peserta))
	for i, p := range m.Peserta {
		peserta[i] = p.Peserta
	}

	// Gerbang 1 - langkah 3, 3.1, 4.
	if !tipeTanpaGerbangDokumen[strings.ToUpper(strings.TrimSpace(m.Tipe))] {
		kurang := nomorPesertaYangGagal(peserta, func(p models.Peserta) bool { return len(p.Dokumen) == 0 })
		if len(kurang) > 0 {
			baris := make([]string, len(kurang))
			for i, n := range kurang {
				baris[i] = pesanDokumenBelumDiunggah + strconv.Itoa(n)
			}
			return tolakRNM("4", strings.Join(baris, "\n"))
		}
	}

	// Langkah 11.7/11.8 hanya punya cabang untuk keempat Type ini.
	if !TypeDikenal(m.Tipe) {
		return fmt.Errorf("%w: %q", ErrTypeTidakDikenal, m.Tipe)
	}
	gross := m.Tipe == TypeQR || m.Tipe == TypeQP

	for i, p := range m.Peserta {
		if err := periksaPesertaRNM(m, p, i+1, gross); err != nil {
			return err
		}
	}
	// ⛔ Langkah 12 "Proteksi Attachment Harus Lengkap" ter-remark (b6178):
	// Save to RNM TIDAK menuntut dokumen lengkap - dan `GetCategoryLife_SQL`
	// tidak disebut rule lain mana pun di seluruh korpus, jadi sistem lama
	// tidak menegakkannya di tempat lain pula.
	// Butir bl (OQ-N6 ditutup): gerbang itu dibuang, bukan ditunda.
	return nil
}

// periksaPesertaRNM adalah langkah 11 untuk satu peserta (`local.idx`).
func periksaPesertaRNM(m MasukanRNM, p PesertaRNM, idx int, gross bool) error {
	// Gerbang 2 - klaim ganda, 11.2-11.6.
	if strings.TrimSpace(m.ContentNote) == contentNoteDeath {
		// 11.4: `CARI1=="" || ==2` lewati; `==0 || ==1` pesan.
		if s := strings.TrimSpace(p.StatusWarisanTerakhir); s == "0" || s == "1" {
			return tolakRNM("11.4", fmt.Sprintf(formatSudahDiaksep, idx))
		}
	} else if p.AdaKlaimSehatSamaDOL {
		return tolakRNM("11.6", fmt.Sprintf(formatSudahDiaksep, idx))
	}

	// Gerbang 3 - DOL, 11.7/11.8 lalu 11.10. ⛔ TANPA pergeseran retro: b4464
	// `@addCalendar(.DATE_OF_LOSS,0,0,0,0,0,0,0)`, beda dengan ValidasiDOL_Act.
	if !dolDalamJendelaRNM(p.Peserta, gross) {
		return tolakRNM("11.10", fmt.Sprintf(formatDOLDiLuar, idx))
	}

	// ⛔ 11.9 dan 11.11 (STNC) ter-remark (b4632, b5009) - tidak ada gerbang.

	// Gerbang 4 - medan kosong, 11.12-11.17, urut.
	switch {
	case p.DOBKosong:
		return tolakRNM("11.12", fmt.Sprintf(formatDOBKosong, idx))
	case kosongTeks(p.TanggalMulai):
		return tolakRNM("11.13", fmt.Sprintf(formatBeginKosong, idx))
	case kosongTeks(p.TanggalExpired):
		return tolakRNM("11.14", fmt.Sprintf(formatExpiredKosong, idx))
	case kosongTeks(p.NomorPolis):
		return tolakRNM("11.15", fmt.Sprintf(formatPolisKosong, idx))
	case kosongTeks(p.NomorSertifikat):
		return tolakRNM("11.16", fmt.Sprintf(formatSertKosong, idx))
	}
	for _, b := range p.Baris {
		// `.CLAIM_GROSS` - kolom kita `CLAIM_AMOUNT` (tiket 03 catatan 7).
		// ⚠️ [sementara - menunggu OQ-N11 pemilik ekspor] OQ-N12 (a), work owner
		// 29-09-2026: korpus tanpa penulis `.CLAIM_GROSS` (4 pembaca, 0 penulis;
		// dugaan Declare Expression yang tidak diekspor). Karena 7.8 membaca sumber
		// kosong sebagai 0, gerbang ini praktis tidak menolak baris Register.
		if b.JumlahKlaim.Kosong() {
			return tolakRNM("11.17.1", fmt.Sprintf(formatGrossKosong, idx))
		}
	}
	return nil
}

// dolDalamJendelaRNM - `local.Begin==false || local.Expired==true` -> tolak.
//
// DOL kosong ditolak (kalimatnya sendiri "cannot be blank"), begitu pula
// jendela yang tidak terbaca: `@CompareDates` atas nilai kosong tidak pernah
// benar, jadi `local.Begin` salah dan pesannya keluar.
func dolDalamJendelaRNM(p models.Peserta, gross bool) bool {
	if kosongTeks(p.TanggalKejadian) {
		return false
	}
	dol, err := utils.ParseTanggal(strings.TrimSpace(p.TanggalKejadian))
	if err != nil {
		return false
	}
	awal, akhir, err := p.JendelaValuasi(gross)
	if err != nil {
		return false
	}
	return sesudah(dol, awal) && !sesudah(dol, akhir)
}

// PolisRetro adalah medan polis yang menentukan gerbang langkah 27.
//
// ⛔ Langkah 27 membaca `pyWorkPage.ClaimData.PolicyDataLife`, BUKAN
// `pyWorkPage.PolicyDataLife`. Salinan itu diisi langkah 25
// (`InsertJsonClaimLife_Act`): langkah 1 menyalin apa adanya, langkah 2
// "Tukar SecurityReinsurer dengan RetroName" (TIDAK ter-remark) menukar kedua
// pengenal - b1152-1153 `ClaimData.RetroID = SecurityReinsurerID`, b1194-1195
// `ClaimData.SecurityReinsurerID = RetroID` - bila WHEN-nya benar:
//
//	b1268 `Type=="TP"||Type=="TR"`
//	b1291 `SecurityReinsurerID!="" && SecurityReinsurer!=""`
//
// ⛔ OQ-N5 DITUTUP 29-09-2026 (GILIRAN-17) `[keputusan work owner]`: keputusan
// Komite "cutover 7 Feb 2025 tidak dipakai" (CONTEXT.md, Komite ronde 1 #4)
// berlaku JUGA di Claim Life. WHEN ketiga, b1314 `OfferFacIn.PolicyData.
// ProdDateTime < "20250207T000000.000 GMT"`, tidak dipakai - dua WHEN cukup.
type PolisRetro struct {
	Tipe                string
	RetroID             string
	SecurityReinsurerID string
	SecurityReinsurer   string
}

// kodeRetroKeluar - b11794 `RetroID=="L0000141" || SecurityReinsurerID==
// "L0000134"` dan b11817 `RetroID=="1000013"`, keduanya WhenTrue 6.
func kodeRetroKeluar(retroID, securityReinsurerID string) bool {
	r, s := strings.TrimSpace(retroID), strings.TrimSpace(securityReinsurerID)
	return r == "L0000141" || s == "L0000134" || r == "1000013"
}

// ArasapasDilewatiRetro adalah langkah 27 - MURNI: true berarti activity
// KELUAR sebelum Arasapas langkah 28. Ia membaca nilai SESUDAH tukar
// `InsertJsonClaimLife_Act` langkah 2 (dua WHEN, OQ-N5).
//
// ⚠️ Modul Komite MEMBUANG gerbang yang sama `[keputusan work owner, OQ-064]`
// untuk `KomitePostAdjustment`. Di Claim Life gerbangnya DIPERTAHANKAN
// (OQ-N3 ditutup GILIRAN-17): XML Claim Life hidup.
func ArasapasDilewatiRetro(p PolisRetro) bool {
	tipe := strings.ToUpper(strings.TrimSpace(p.Tipe))
	tukar := (tipe == TypeTP || tipe == TypeTR) &&
		!kosongTeks(p.SecurityReinsurerID) && !kosongTeks(p.SecurityReinsurer)
	if tukar {
		return kodeRetroKeluar(p.SecurityReinsurerID, p.RetroID)
	}
	return kodeRetroKeluar(p.RetroID, p.SecurityReinsurerID)
}

// ErrSimpanRNMBukanOutstanding - tombol `Save to RNM` hanya ada di layar
// Outstanding (`InputOSClaimLife.xml` b21102), yang dipegang Admin.
var ErrSimpanRNMBukanOutstanding = errors.New(
	"services: Save to RNM hanya tersedia pada tahap Outstanding Claim")

// HasilSimpanRNM adalah hasil Save to RNM yang berhasil.
type HasilSimpanRNM struct {
	NomorKlaim string `json:"nomorKlaim"`
	// NomorBaru - nomor diterbitkan di sini (langkah 16-20), bukan saat pendaftaran.
	NomorBaru bool `json:"nomorBaru"`
	// BarisDitandai - baris tanpa status yang kini Outstanding (langkah 22.1.3.2).
	BarisDitandai int `json:"barisDitandai"`
	// Arasapas - keadaan efek langkah 28, sebagai KATA.
	Arasapas string `json:"arasapas"`
}

// SimpanRNM adalah layanan Save to RNM.
type SimpanRNM struct {
	svc      *Service
	penomor  Penomor
	jejak    jejak.Jejak
	penyalur *outbox.Penyalur
}

// SimpanRNM menyusun layanannya dengan ketergantungan yang GAGAL TERANG.
func (s *Service) SimpanRNM() *SimpanRNM {
	return &SimpanRNM{
		svc: s, penomor: PenomorBelumDiputuskan{}, jejak: jejak.JejakBelumDiputuskan{},
		penyalur: outbox.NewPenyalur(s.Lingkungan(), outbox.AntreanBelumDiputuskan{},
			outbox.EfekArasapas{Resolver: layanan.ResolverBelumDiputuskan{}}),
	}
}

// SimpanRNMOracle menyusun layanannya dengan ketergantungan Oracle - satu
// tempat, supaya handler tidak merakit sendiri.
func SimpanRNMOracle(s *Service) *SimpanRNM {
	return &SimpanRNM{
		svc: s, penomor: PenomorCounterOracle(s), jejak: jejak.PerekamJejakOracle(s),
		penyalur: outbox.NewPenyalur(s.Lingkungan(), AntreanEfekOracle(s),
			outbox.EfekArasapas{Resolver: layanan.ResolverLinkServiceOracle(s)}),
	}
}

// tanggalSaja memotong `YYYY-MM-DD HH24:MI:SS` menjadi `YYYY-MM-DD`.
func tanggalSaja(teks string) string {
	t := strings.TrimSpace(teks)
	if len(t) >= 10 {
		return t[:10]
	}
	return ""
}

// kataHasilSalur menulis hasil efek langkah 28 sebagai kata (km5).
func kataHasilSalur(h outbox.HasilSalur) string {
	switch {
	case h.Dilewati:
		return "dilewati: lingkungan bukan produksi (IsPEGAPROD)"
	case len(h.GagalDiantre) > 0:
		return "gagal, dan TIDAK dapat diantre ulang"
	case len(h.Gagal) > 0:
		return "gagal, diantre ulang"
	default:
		return "terkirim"
	}
}

// Simpan menjalankan Save to RNM atas satu klaim.
//
// Urutan: identitas, pengenal, kasus terbuka, tahap Outstanding + pemegangnya,
// lalu SELURUH gerbang XML (PeriksaSimpanRNM) SEBELUM satu tulisan pun, lalu
// satu transaksi, lalu Arasapas sesudah commit.
func (x *SimpanRNM) Simpan(ctx context.Context, pelaku inti.Pelaku, klaimID string,
	saat time.Time) (HasilSimpanRNM, error) {

	var hasil HasilSimpanRNM
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return hasil, err
	}
	if strings.TrimSpace(klaimID) == "" {
		return hasil, fmt.Errorf("%w: pengenal klaim wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if x == nil || x.svc == nil || !x.svc.PunyaDatabase() {
		return hasil, db.ErrTanpaOracle
	}
	// ⛔ BUTIR bb - kasus tertutup tidak dapat diubah lagi.
	if err := x.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		// Klaim yang tidak ada dijawab 404, bukan 500 - handler tidak
		// mengimpor repository, jadi galatnya diterjemahkan di sini.
		if errors.Is(err, repository.ErrWorkTidakAda) {
			return hasil, fmt.Errorf("%w: %w", ErrKlaimTidakAda, err)
		}
		return hasil, err
	}

	baca := repository.NewKlaimLife(x.svc.DB())
	kolomTahap, peranPemegang, err := baca.TahapDanPeran(ctx, klaimID)
	if err != nil {
		return hasil, err
	}
	// ⚠️ Aturan cadangannya dari models.TahapBerlaku (satu sumber), tetapi
	// tahap tak dikenal di sini BUKAN galat tersendiri: ia bukan Outstanding,
	// dan itulah yang dijawab (409) - perilaku yang sudah diuji tetap.
	tahap := models.TahapBerlaku(kolomTahap, peranPemegang)
	if tahap != models.TahapOutstanding {
		return hasil, fmt.Errorf("%w: tahap %s", ErrSimpanRNMBukanOutstanding, tahap)
	}
	// Orang menyimpan pekerjaan yang sedang ia pegang (pola Pindah).
	if err := WajibPemegangTahap(pelaku, tahap); err != nil {
		return hasil, err
	}

	klaim, err := x.svc.KlaimLife().Ambil(ctx, klaimID)
	if err != nil {
		return hasil, err
	}
	tipe, err := baca.TypeKlaim(ctx, klaimID)
	if err != nil {
		return hasil, err
	}
	polis, err := repository.NewRingkasPolisLife(x.svc.DB()).Ringkas(ctx, klaim.NomorPolis)
	if err != nil {
		return hasil, err
	}
	note, err := ContentNoteDari(polis.BusinessCode)
	if err != nil {
		return hasil, err
	}
	m := MasukanRNM{Tipe: tipe, ContentNote: note}
	// OQ-N2 (GILIRAN-17): baris cermin klaim ini sendiri dikecualikan dari
	// pemeriksaan klaim ganda - di Pega ia belum ada saat langkah 11.x.
	caseID, err := baca.CaseIDKlaim(ctx, klaimID)
	if err != nil {
		return hasil, err
	}
	if strings.TrimSpace(caseID) == "" {
		// CASEID klaim aplikasi = pengenal work (butir ae1). Tanpa ini
		// `o.CASEID <> NULL` akan menyaring SELURUH baris warisan.
		caseID = klaimID
	}
	for _, p := range klaim.Peserta {
		if strings.TrimSpace(p.SumberID) == "" {
			return hasil, fmt.Errorf("%w: peserta %q tanpa SOURCE_ID; DOB dan klaim ganda "+
				"tidak dapat diperiksa", galat.ErrPermintaanTidakSah, p.ID)
		}
		k := repository.KunciPesertaSumber{PLNumber: p.NomorPremiList,
			Sertifikat: p.NomorSertifikat, SumberID: p.SumberID}
		pr := PesertaRNM{Peserta: p}
		if pr.DOBKosong, err = baca.DOBSumberKosong(ctx, k); err != nil {
			return hasil, err
		}
		if strings.TrimSpace(note) == contentNoteDeath {
			if pr.StatusWarisanTerakhir, err = baca.StatusWarisanTerakhir(ctx, k, polis.CedingCo, caseID); err != nil {
				return hasil, err
			}
		} else if dol := tanggalSaja(p.TanggalKejadian); dol != "" {
			if pr.AdaKlaimSehatSamaDOL, err = baca.AdaWarisanSamaDOL(ctx, k, polis.CedingCo, dol, caseID); err != nil {
				return hasil, err
			}
		}
		m.Peserta = append(m.Peserta, pr)
	}
	if err := PeriksaSimpanRNM(m); err != nil {
		return hasil, err
	}

	hasil.NomorKlaim = klaim.NomorKlaim
	err = x.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		// Langkah 16-20 (13-15 ter-remark) - hanya bila CLAIM_NO masih kosong.
		if strings.TrimSpace(klaim.NomorKlaim) == "" {
			nomor, err := x.penomor.NomorBerikut(ctx, tx, polis.BusinessCode, saat)
			if err != nil {
				return err
			}
			if err := baca.IsiNomorKlaimKosong(ctx, tx, klaimID, nomor); err != nil {
				return err
			}
			hasil.NomorKlaim, hasil.NomorBaru = nomor, true
		}
		// Langkah 22.1.3.2 - `.STS_REJECT = 0` bagi baris yang BELUM pernah
		// tersimpan. ⛔ Baris berstatus TIDAK ditimpa: tanpa bendera `Save`
		// (OQ-N1) tombolnya dapat ditekan lagi, dan menimpa baris yang sudah
		// ditolak Admin berarti membatalkan penolakan diam-diam.
		// ⛔ HANYA baris adjustment (sensus penulis STS_REJECT,
		// repository.PerbaruiStatusBaris): langkah ini tidak menulis status
		// peserta, nomor akseptasi, atau tanggalnya - karena itu bukan
		// PerbaruiStatusBaris, yang menulis keempatnya.
		for i := range klaim.Peserta {
			p := &klaim.Peserta[i]
			for j := range p.Baris {
				b := &p.Baris[j]
				if strings.TrimSpace(b.KodeStatus) != "" {
					continue
				}
				if err := baca.TandaiBarisOutstanding(ctx, tx, b.ID); err != nil {
					return err
				}
				b.KodeStatus = models.KodeOutstanding
				if err := x.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
					AdjustmentID: b.ID, KlaimID: klaimID, Dari: "", Ke: models.KodeOutstanding,
					AkunID: pelaku.AkunID, Waktu: saat,
				}); err != nil {
					return err
				}
				hasil.BarisDitandai++
			}
		}
		// ⚠️ SELISIH SADAR dengan XML, dicatat: activity ini tidak menulis
		// header. `ClaimData.STS_REJECT/ACCEPTEDNO` hanya ditulis
		// `serviceInsertArasapasClaimLife_act` langkah 1.1.1 (b371, b417),
		// yang dicapai lewat langkah 28 - produksi saja, dan tidak bila
		// langkah 27 keluar. Di aplikasi ini SETIAP penulis status baris
		// mencerminkan header dalam transaksinya sendiri (PerbaruiStatusBaris,
		// tiket 04), karena Arasapas di sini efek keluar yang dapat gagal
		// atau dilewati; di sini pun sama, termasuk saat Arasapas dilewati
		// (OQ-N3 - dipertahankan di transaksi simpan, GILIRAN-17).
		if hasil.BarisDitandai > 0 {
			akhir := BarisTerakhir(klaim)
			if err := baca.CerminkanHeader(ctx, tx, klaimID, akhir.KodeStatus,
				akhir.NomorAkseptasi); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return hasil, err
	}

	// Langkah 27-28 - SESUDAH commit, dan kegagalannya bukan galat simpan
	// (ADR-U-0008).
	//
	// ⚠️ Di sistem lama keluar di langkah 27 (kode 6) melewati Obj-Save
	// langkah 29 juga; di sini tulisan langkah 16-22 sudah di-commit lebih
	// dulu. Lihat OQ-N3.
	lewat := ArasapasDilewatiRetro(PolisRetro{Tipe: tipe, RetroID: polis.RetroID,
		SecurityReinsurerID: polis.SecurityReinsurerID, SecurityReinsurer: polis.SecurityReinsurer})
	switch {
	case lewat:
		hasil.Arasapas = "dilewati: kode retro langkah 27"
	default:
		hasil.Arasapas = kataHasilSalur(x.penyalur.Salurkan(ctx,
			outbox.MuatanEfek{KlaimID: klaimID, AkunID: pelaku.AkunID, Waktu: saat}))
	}
	return hasil, nil
}
