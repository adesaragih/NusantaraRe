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
//	11.7-10 DOL dalam [BEGIN, EXPIRED] inklusif (WO 11-10-2026) gerbang 3
//	11.9/11 STNC                                           `//` - TIDAK ditiru
//	11.12-17 DOB, BEGIN, EXPIRED, POLICY, CERT, CLAIM_GROSS gerbang 4
//	12      dokumen LENGKAP                                 `//` - TIDAK ditiru
//	13-15   nomor lewat `Generate_NoKlaim_Life*`            `//` - TIDAK ditiru
//	16-20   nomor klaim bila CLAIM_NO kosong                penomor yang sama dengan pendaftaran
//	21, 24  bendera `pyWorkPage.Save`                       ⚠️ TANPA kolom - OQ-N1
//	22      insert baris warisan + STS_REJECT=0             STS_REJECT hanya baris tanpa status
//	                                                        (tiket 04); baris datar `OS_AKSEPTASI_KLAIM_LIFE`
//	                                                        LAHIR di sini (RALAT 10-10-2026, TulisCerminSaveRNM:
//	                                                        UPDATE bila ada, INSERT bila baris belum disimpan)
//	                                                        lalu disetel '0' (b176, OQ-N13 GILIRAN-18)
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

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/jejak"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
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
	// Kekurangan - SELURUH kolom wajib yang kosong / tidak sah, semua peserta (keputusan work owner 11-10-2026). Kosong
	// untuk gerbang yang berhenti dini (klaim ganda). Pesan = gabungan baris Pesan tiap kekurangan.
	Kekurangan []KekuranganRNM
}

// KekuranganRNM - satu kolom wajib yang belum terisi / tidak sah. NomorPeserta 0 = tingkat klaim (mis. Is KPR).
type KekuranganRNM struct {
	NomorPeserta int    `json:"nomorPeserta"`
	PesertaID    string `json:"pesertaId"`
	Sertifikat   string `json:"sertifikat"`
	Nama         string `json:"nama"`
	Kolom        string `json:"kolom"`
	Pesan        string `json:"pesan"`
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
	// IsKPR - header klaim (wajib sebelum Save to RNM, keputusan work owner 11-10-2026).
	IsKPR string
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
	// Keputusan work owner 11-10-2026: DATE OF LOSS, CLAIM RECEIVED DATE, CONFIRMATION DATE WAJIB sebelum Save to RNM.
	// DOL kosong sudah ditolak langkah 11.10 (formatDOLDiLuar, kalimat XML); dua lainnya aturan baru aplikasi.
	formatTerimaKlaimKosong = "Claim Received Date is required for insured number %d"
	formatKonfirmasiKosong  = "Confirmation Date is required for insured number %d"
	// Keputusan work owner 11-10-2026: Claim RNM WAJIB (> 0) menggantikan 11.17.1 formatGrossKosong (yang hanya menolak
	// kosong); Is KPR (header) wajib.
	formatClaimRNMKosong = "Claim RNM is required for insured number %d"
	pesanIsKPRKosong     = "Is KPR is required."
	// Claim RNM > Share Nusantara Re (11-10-2026, dipindah dari Save grid) - kalimat layar galatClaimRNM + nomor peserta.
	formatClaimRNMMelebihi = "Insured number %d: Claim RNM cannot be greater than Share Nusantara Re (%s)."
)

// KlaimRNMMinimal - Claim RNM harus LEBIH BESAR dari nilai ini (keputusan work owner 11-10-2026: wajib; 0 belum boleh).
// ⚠️ Pega hanya menolak kosong (11.17.1 "Claim Gross No N can't null") dan Register 7.8 mengisi 0 - bila work owner
// kelak membolehkan 0, ubah aturan di bawah (claimRNMTerisi) menjadi >= 0.
const KlaimRNMMinimal = "0"

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
	// Urutan (11-10-2026, keputusan work owner "pastikan semua kolom wajib harus diisi untuk semua peserta"):
	//  1. Type harus dikenal (cabang DOL 11.7/11.8) - berhenti dini, galat aplikasi.
	//  2. Klaim ganda 11.2-11.6 (bukan "kolom kosong") - berhenti dini pada peserta pertama, kalimat XML.
	//  3. SELURUH kekurangan kolom wajib semua peserta dikumpulkan lalu dikembalikan SEKALIGUS (422):
	//     dokumen peserta di-Accept (langkah 3-4), Is KPR (header), DOL (kosong / di luar jendela 11.10), DOB 11.12,
	//     Begin 11.13, Expired 11.14, Policy 11.15, Certificate 11.16, Claim Received, Confirmation, Claim RNM tiap baris
	//     (11.17.1, kini > 0).
	// ⛔ MENYIMPANG dari Pega, yang berhenti di pelanggaran pertama. Nomor klaim tetap diambil SESUDAH seluruh gerbang lolos.
	if !TypeDikenal(m.Tipe) {
		return fmt.Errorf("%w: %q", ErrTypeTidakDikenal, m.Tipe)
	}
	gross := m.Tipe == TypeQR || m.Tipe == TypeQP

	for i, p := range m.Peserta {
		if err := periksaGandaRNM(m, p, i+1); err != nil {
			return err
		}
	}

	var kurang []KekuranganRNM
	tambah := func(idx int, p *PesertaRNM, kolom, pesan string) {
		k := KekuranganRNM{NomorPeserta: idx, Kolom: kolom, Pesan: pesan}
		if p != nil {
			k.PesertaID, k.Sertifikat, k.Nama = p.ID, p.NomorSertifikat, p.NamaTertanggung
		}
		kurang = append(kurang, k)
	}
	if kosongTeks(m.IsKPR) {
		tambah(0, nil, "Is KPR", pesanIsKPRKosong)
	}
	dokumenWajib := !tipeTanpaGerbangDokumen[strings.ToUpper(strings.TrimSpace(m.Tipe))]
	for i := range m.Peserta {
		p := &m.Peserta[i]
		idx := i + 1
		// Langkah 3-4: peserta di-Accept tanpa dokumen (kalimat XML).
		if dokumenWajib && diterimaDiOutstanding(p.Peserta) && len(p.Dokumen) == 0 {
			tambah(idx, p, "Document", pesanDokumenBelumDiunggah+strconv.Itoa(idx))
		}
		// DOL kosong: kalimat XML 11.10; DOL TERISI di luar jendela: kalimat layar yang sama dengan Save grid, dengan
		// rentang valuasi (11-10-2026).
		// Jendela [BEGIN, EXPIRED] inklusif, SAMA dengan Save grid (DOLDalamJendela, keputusan work owner 11-10-2026).
		if !dolDalamJendelaRNM(p.Peserta, gross) {
			k, ok := kalimatDOLRNM(p.Peserta, gross, idx)
			if !ok {
				k = fmt.Sprintf(formatDOLDiLuar, idx)
			}
			tambah(idx, p, "Date of Loss", k)
		}
		// ⛔ 11.9 dan 11.11 (STNC) ter-remark (b4632, b5009) - tidak ada gerbang.
		if p.DOBKosong {
			tambah(idx, p, "DOB", fmt.Sprintf(formatDOBKosong, idx))
		}
		if kosongTeks(p.TanggalMulai) {
			tambah(idx, p, "Begin Date", fmt.Sprintf(formatBeginKosong, idx))
		}
		if kosongTeks(p.TanggalExpired) {
			tambah(idx, p, "Expired Date", fmt.Sprintf(formatExpiredKosong, idx))
		}
		if kosongTeks(p.NomorPolis) {
			tambah(idx, p, "Policy No", fmt.Sprintf(formatPolisKosong, idx))
		}
		if kosongTeks(p.NomorSertifikat) {
			tambah(idx, p, "Certificate No", fmt.Sprintf(formatSertKosong, idx))
		}
		if kosongTeks(p.TanggalTerimaKlaim) {
			tambah(idx, p, "Claim Received Date", fmt.Sprintf(formatTerimaKlaimKosong, idx))
		}
		if kosongTeks(p.TanggalKonfirmasi) {
			tambah(idx, p, "Confirmation Date", fmt.Sprintf(formatKonfirmasiKosong, idx))
		}
		for _, b := range p.Baris {
			// `.CLAIM_GROSS` = kolom CLAIM_AMOUNT = Claim RNM (tiket 03 catatan 7). Satu kekurangan per peserta cukup.
			if !claimRNMTerisi(b) {
				tambah(idx, p, "Claim RNM", fmt.Sprintf(formatClaimRNMKosong, idx))
				break
			}
			// Claim RNM <= Share Nusantara Re - dipindah dari Save grid (keputusan work owner 11-10-2026: Save = simpan
			// progres); kalimat SAMA dengan layar (galatClaimRNM). Baris berstatus sudah lolos saat tersimpan.
			if strings.TrimSpace(b.KodeStatus) == "" && claimRNMMelebihiShare(b) {
				tambah(idx, p, "Claim RNM", fmt.Sprintf(formatClaimRNMMelebihi, idx, angkaLayar(b.ShareNusantaraRe)))
				break
			}
		}
	}
	if len(kurang) > 0 {
		baris := make([]string, len(kurang))
		for i, k := range kurang {
			baris[i] = k.Pesan
		}
		return &PelanggaranRNM{Langkah: "wajib", Pesan: strings.Join(baris, "\n"), Kekurangan: kurang}
	}
	// ⛔ Langkah 12 "Proteksi Attachment Harus Lengkap" ter-remark (b6178): Save to RNM TIDAK menuntut dokumen lengkap.
	// Butir bl (OQ-N6 ditutup): gerbang itu dibuang, bukan ditunda.
	return nil
}

// claimRNMTerisi - Claim RNM terisi dan > KlaimRNMMinimal (keputusan work owner 11-10-2026).
func claimRNMTerisi(b models.BarisAdjustment) bool {
	if b.JumlahKlaim.Kosong() {
		return false
	}
	batas, _, err := apd.NewFromString(KlaimRNMMinimal)
	if err != nil {
		return false
	}
	return b.JumlahKlaim.Amount.Cmp(batas) > 0
}

// claimRNMMelebihiShare - Claim RNM > Share Nusantara Re (aturan models.PeriksaClaimGrossTerhadapShare tanpa pengecualian
// kode bisnis, keputusan work owner 11-10-2026). Mata uang baris = mata uang share.
func claimRNMMelebihiShare(b models.BarisAdjustment) bool {
	// Share kosong: tidak dibandingkan (sama dengan layar galatClaimRNM, yang melewati share tak terbaca).
	if b.JumlahKlaim.Kosong() || b.ShareNusantaraRe.Kosong() {
		return false
	}
	share := b.ShareNusantaraRe
	share.Currency = b.JumlahKlaim.Currency
	return errors.Is(models.PeriksaClaimGrossTerhadapShare(b.JumlahKlaim, share, ""), models.ErrClaimGrossMelebihiShare)
}

// angkaLayar - angka bergaya Indonesia seperti layar (formatNumber tanpa batas desimal): "5.077.429,5"; kosong = "0".
func angkaLayar(m uang.Money) string {
	if m.Kosong() {
		return "0"
	}
	t := m.Amount.Text('f')
	neg := strings.HasPrefix(t, "-")
	t = strings.TrimPrefix(t, "-")
	bulat, pecahan, _ := strings.Cut(t, ".")
	pecahan = strings.TrimRight(pecahan, "0")
	var sb strings.Builder
	for i, c := range bulat {
		if i > 0 && (len(bulat)-i)%3 == 0 {
			sb.WriteByte('.')
		}
		sb.WriteRune(c)
	}
	out := sb.String()
	if pecahan != "" {
		out += "," + pecahan
	}
	if neg {
		out = "-" + out
	}
	return out
}

// periksaGandaRNM - gerbang 2 klaim ganda, 11.2-11.6, satu peserta (`local.idx`). Berhenti dini (bukan kolom kosong).
func periksaGandaRNM(m MasukanRNM, p PesertaRNM, idx int) error {
	if strings.TrimSpace(m.ContentNote) == contentNoteDeath {
		// 11.4: `CARI1=="" || ==2` lewati; `==0 || ==1` pesan.
		if s := strings.TrimSpace(p.StatusWarisanTerakhir); s == "0" || s == "1" {
			return tolakRNM("11.4", fmt.Sprintf(formatSudahDiaksep, idx))
		}
	} else if p.AdaKlaimSehatSamaDOL {
		return tolakRNM("11.6", fmt.Sprintf(formatSudahDiaksep, idx))
	}
	return nil
}

// dolDalamJendelaRNM - `local.Begin==false || local.Expired==true` -> tolak.
//
// DOL kosong ditolak (kalimatnya sendiri "cannot be blank"), begitu pula
// jendela yang tidak terbaca: `@CompareDates` atas nilai kosong tidak pernah
// benar, jadi `local.Begin` salah dan pesannya keluar.
// kalimatDOLRNM - kalimat "Insured number N: Date of Loss must be within the valuation period (…)" untuk DOL TERISI di luar
// jendela; false bila DOL kosong / tak terurai atau jendela tidak terbaca (kalimat XML 11.10 yang dipakai).
func kalimatDOLRNM(p models.Peserta, gross bool, nomor int) (string, bool) {
	if kosongTeks(p.TanggalKejadian) {
		return "", false
	}
	if _, err := utils.ParseTanggal(strings.TrimSpace(p.TanggalKejadian)); err != nil {
		return "", false
	}
	awal, akhir, err := p.JendelaValuasi(gross)
	if err != nil {
		return "", false
	}
	return GalatDOL{PesertaID: p.ID, Nomor: nomor, Awal: tanggalLayar(awal), Akhir: tanggalLayar(akhir)}.Kalimat(), true
}

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
	// Aturan yang SAMA dengan Save grid (ValidasiDOL): [BEGIN, EXPIRED] inklusif, tanpa geser (keputusan WO 11-10-2026).
	return DOLDalamJendela(dol, awal, akhir)
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
	// CerminDisetel - baris cermin warisan yang kini '0' (b176, OQ-N13).
	CerminDisetel int `json:"cerminDisetel"`
	// CerminDitulis - baris datar OS_AKSEPTASI_KLAIM_LIFE yang ditulis/diperbarui Save to RNM (RALAT 10-10-2026).
	CerminDitulis int `json:"cerminDitulis"`
	// Arasapas - keadaan efek langkah 28, sebagai KATA.
	Arasapas string `json:"arasapas"`
}

// langkahCerminSaveRNM memutuskan apa yang dikerjakan Save to RNM atas satu
// baris: baris tanpa status DITANDAI (22.1.3.2) dan dicerminkan (b176); baris
// yang sudah Outstanding HANYA dicerminkan - cermin NULL milik klaim yang
// disimpan sebelum N13 pulih saat tombolnya ditekan ulang (OQ-N1, tombol
// tetap hidup); baris berkeputusan tidak disentuh sama sekali.
//
// Temuan /code-review GILIRAN-18. `SetelCerminOutstanding` sendiri hanya
// menyentuh cermin yang masih NULL, jadi mencerminkan ulang baris '0' aman.
func langkahCerminSaveRNM(kode string) (tandai, cermin bool) {
	switch strings.TrimSpace(kode) {
	case "":
		return true, true
	case kontrak.KodeOutstanding:
		return false, true
	}
	return false, false
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
	// Kunci polis = IDPEGA peserta (kuncipolis.go, 10-10-2026); KlaimLife().Ambil sudah menurunkannya dari SOURCE_ID.
	polis, kunciPolis, err := RingkasPolisIdPega(ctx, x.svc.PembacaPolis(), IdPegaPeserta(klaim.Peserta))
	if err != nil {
		return hasil, err
	}
	note, err := ContentNoteDari(polis.BusinessCode)
	if err != nil {
		return hasil, err
	}
	m := MasukanRNM{Tipe: tipe, ContentNote: note, IsKPR: klaim.IsKPR}
	// OQ-N2 (GILIRAN-17): baris cermin klaim ini sendiri dikecualikan dari
	// pemeriksaan klaim ganda - di Pega ia belum ada saat langkah 11.x.
	// CASEID klaim = ID klaim (butir ae1; kolom CASE_ID dibuang migrasi 023,
	// keputusan work owner 01-10-2026) - tidak pernah kosong, jadi
	// `o.CASEID <> NULL` tidak dapat menyaring seluruh baris warisan.
	caseID := klaimID
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
	// Langkah 16-20 (13-15 ter-remark) - hanya bila CLAIM_NO masih kosong; sejak 10-10-2026 SATU-SATUNYA tempat nomor
	// klaim lahir (Register tidak lagi menomori - keputusan work owner).
	//
	// ⛔ TRANSAKSI PENOMORAN TERSENDIRI, di-commit segera (AC tiket 02 "commit terjadi segera setelah nomor terbentuk,
	// sehingga lock `SELECT … FOR UPDATE` pada `GENERATE_SEQUENCE_NUMBER` tidak menahan pendaftar lain"): nomor DAN
	// CLAIM_NO-nya ditulis bersama lalu di-commit, sebelum transaksi Save to RNM. Bila transaksi utama kemudian gagal,
	// klaim TETAP bernomor dan Save berikutnya tidak menomori ulang (`IsiNomorKlaimKosong` hanya mengisi yang kosong) -
	// tanpa nomor hangus. Seluruh gerbang Save to RNM (langkah 1-11) sudah lolos di atas.
	if strings.TrimSpace(klaim.NomorKlaim) == "" {
		if err := x.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
			nomor, err := x.penomor.NomorBerikut(ctx, tx, polis.BusinessCode, saat)
			if err != nil {
				return err
			}
			if err := baca.IsiNomorKlaimKosong(ctx, tx, klaimID, nomor); err != nil {
				return err
			}
			hasil.NomorKlaim, hasil.NomorBaru = nomor, true
			return nil
		}); err != nil {
			return hasil, err
		}
	}
	// ⛔ RALAT 10-10-2026 (keputusan work owner "saat register claim, data peserta belum di insert ke tabel
	// os_akseptasi_klaim_life; itu hanya terjadi saat save to rnm bersamaan dengan no klaim dan no akseptasi"): baris
	// datar disusun dari klaim yang SUDAH bernomor (transaksi penomoran di atas) dan ditulis di transaksi UTAMA di
	// bawah - bila transaksi utama gagal, tidak ada baris datar setengah jadi. Nomor akseptasi TIDAK lahir di Save to
	// RNM (ia lahir di Akseptasi / Komite); NO_ACCEPTATION diisi TANGGAL OUTSTANDING (keputusan work owner 11-10-2026).
	pembuat, err := baca.NamaPembuatKlaim(ctx, klaimID)
	if err != nil {
		return hasil, err
	}
	klaimDatar := *klaim
	klaimDatar.NomorKlaim = hasil.NomorKlaim
	klaimDatar.Type = tipe
	// SOB id polis (OS_AKSEPTASI_KLAIM_LIFE.SOB) - kunci IDPEGA yang sama dengan header polis di atas.
	sob, err := baca.SOBPolis(ctx, kunciPolis)
	if err != nil {
		return hasil, err
	}
	dasar := map[string]repository.BarisLama{}
	for _, b := range repository.BarisLamaDari(models.PohonKlaim{
		Work: models.WorkClaim{ID: caseID, CreateOpName: pembuat}, Klaim: klaimDatar}) {
		// Status cermin hanya lewat SetelCerminOutstanding (satu jalur, '0' b176).
		b.STS_REJECT = ""
		dasar[b.ID] = b
	}
	// Keputusan work owner 11-10-2026: nilai polis, peserta, dan baris ikut ditulis (SusunCerminSaveRNM); NO_ACCEPTATION =
	// tanggal outstanding (saat Save to RNM, waktu yang sama dengan jejak OutStanding Date).
	datar := map[string]repository.BarisLama{}
	for _, p := range klaim.Peserta {
		for _, b := range p.Baris {
			if d, ada := dasar[b.ID]; ada {
				datar[b.ID] = SusunCerminSaveRNM(d, p, b, polis, sob, saat)
			}
		}
	}
	pohon := repository.NewPohonKlaim(x.svc.DB())
	err = x.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
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
				tandai, cermin := langkahCerminSaveRNM(b.KodeStatus)
				if !cermin {
					continue
				}
				if tandai {
					if err := baca.TandaiBarisOutstanding(ctx, tx, b.ID); err != nil {
						return err
					}
				}
				// Langkah 22.1.3.1 `InsertJsonKlaimLife_sql` - baris datar baris ini (RALAT 10-10-2026): UPDATE bila
				// sudah ada (klaim lama yang sempat ditulis Register, atau tombol ditekan ulang), INSERT hanya bila
				// baris BELUM pernah disimpan. Idempoten: kunci ID + CASEID.
				kunci := repository.KunciPesertaSumber{PLNumber: p.NomorPremiList,
					Sertifikat: p.NomorSertifikat, SumberID: p.SumberID}
				ada, err := pohon.TulisCerminSaveRNM(ctx, tx, datar[b.ID], kunci,
					repository.CerminBolehDisisip(b.KodeStatus))
				if err != nil {
					return err
				}
				if ada {
					hasil.CerminDitulis++
				}
				// OQ-N13 (GILIRAN-18) - langkah 22.1.3.1 `InsertJsonKlaimLife_sql`
				// b176: cermin baris ini '0', supaya klaim ganda antarklaim baru
				// tertangkap tanpa menunggu Komite. `caseID` yang SAMA dengan
				// pengecualian pemeriksa ganda di atas.
				n, err := baca.SetelCerminOutstanding(ctx, tx, b.ID, caseID)
				if err != nil {
					return err
				}
				hasil.CerminDisetel += n
				if !tandai {
					continue
				}
				b.KodeStatus = kontrak.KodeOutstanding
				if err := x.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
					AdjustmentID: b.ID, KlaimID: klaimID, Dari: "", Ke: kontrak.KodeOutstanding,
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
