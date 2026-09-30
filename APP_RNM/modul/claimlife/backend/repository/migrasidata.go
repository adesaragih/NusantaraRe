package repository

// Pembongkar data lama - dari satu tabel flat ke pohon klaim yang baru.
//
// Untuk apa berkas ini: mengubah baris `OS_AKSEPTASI_KLAIM_LIFE` (tabel datar
// warisan, 55 kolom) menjadi pohon klaim bertingkat, sambil MELAPORKAN apa pun
// yang tidak dapat dipindahkan dengan yakin.
//
// Dibaca sesudah: migrasi.go (bentuk tabel barunya) dan models/pohonklaim.go.
//
// Istilah:
//   - baris lama   : satu baris tabel datar warisan. Satu baris = satu peserta.
//   - rekonsiliasi : membandingkan yang masuk dengan yang keluar, digit demi
//                    digit, lalu melaporkan selisihnya.
//
// Fungsi di berkas ini MURNI: ia tidak menyentuh basis data sama sekali.
// Masukannya potongan data, keluarannya pohon dan laporan. Karena itu ia dapat
// diuji sepenuhnya tanpa instance Oracle.
//
// Aturan yang dijaga:
//   ADR-U-0011  SELURUH baris adjustment ikut pindah, bukan hanya keadaan
//               terakhir.
//   ADR-U-0021  nilai uang dibandingkan TEPAT, bukan dengan toleransi.
//   ADR-U-0022  konversi tipe sekali saat masuk; kode dan penanda tetap teks.
//   ADR-U-0027  kolom kosong tetap kosong, tidak menjadi nol.

import (
	"fmt"
	"sort"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/claimlife/backend/models"
)

// BarisLama adalah satu baris `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`.
//
// Kelima puluh lima medan di bawah `[terverifikasi]` berasal dari rule
// `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL`
// bertipe `Rule-Connect-SQL` - namanya menyebut "Update" tetapi isinya INSERT.
//
// Seluruhnya bertipe teks, apa adanya seperti tersimpan. Konversi ke desimal
// dan tanggal dilakukan SEKALI, di sini, saat membongkar (ADR-U-0022).
type BarisLama struct {
	CASEID, NO_CLAIM, POLICY_NO, POLICY_HOLDER, CERTIFICATE_NO                  string
	NAME_OF_INSURED, SEX, DOB, AGE, PLAN                                        string
	BEGIN_DATE, LAPSE_DATE, EXPIRED_DATE, STATUS, EM_PERCENT                    string
	CURRENCY, SUM_INSURED, CEDING_RETENTION, SUM_REASURED, SHARE_NUSANTARA_RE   string
	CLAIM_AMOUNT, WPC, PL_NUMBER, DISEASE, ICD_CODE                             string
	NOTES, CEDINGCO, CEDINGCONAME, SOB, SOBNAME                                 string
	BUSINESSID, BUSINESSNAME, SHARE_RETRO, KETERANGAN, ACCEPTATION_DATE         string
	NO_ACCEPTATION, STS_REJECT, CLAIM_RETRO, RETROID, RETRONAME                 string
	SECURITYREINSURERID, SECURITYREINSURER, TYPECEDING, TYPE, CONFIRMATION_DATE string
	CLAIM_RECEIVED_DATE, COMPLETE_DATE, ID, NAME_OF_BANK, IDBANK                string
	ACCOUNTNO, CREATEOPNAME, PRODUCTNAMEID, PRODUCTNAME, RETROCEDED_SHARE       string
}

// Jenis temuan rekonsiliasi.
const (
	TemuanTanggalTakTerurai = "tanggal tidak dapat diurai"
	TemuanUangTakTerurai    = "nilai uang tidak dapat diurai"
	TemuanAtributBerbeda    = "atribut polis berbeda antar baris peserta"
	TemuanNilaiHardcode     = "nilai berasal dari hardcode sumber"
	TemuanIdentitasKosong   = "identitas baris kosong"
)

// Temuan adalah satu hal yang tidak dapat dipindahkan dengan yakin.
//
// Ia dilaporkan, TIDAK didiamkan, dan TIDAK ditebak.
type Temuan struct {
	Jenis   string
	Sumber  string // pengenal baris lama, supaya dapat ditelusuri
	Medan   string
	Nilai   string
	Catatan string
}

func (t Temuan) String() string {
	return fmt.Sprintf("%s | %s | %s=%q | %s", t.Jenis, t.Sumber, t.Medan, t.Nilai, t.Catatan)
}

// LaporanRekonsiliasi adalah keluaran pembongkaran, di samping pohonnya.
//
// Cacahnya dipakai untuk membuktikan bahwa jumlah baris sesudah migrasi SAMA
// dengan sebelumnya.
type LaporanRekonsiliasi struct {
	BarisMasuk          int
	KlaimTerbentuk      int
	PesertaTerbentuk    int
	AdjustmentTerbentuk int
	// BarisHardcode mencacah baris yang nilainya diketahui berasal dari
	// hardcode di sumbernya, bukan dari keadaan sebenarnya.
	BarisHardcode int
	Temuan        []Temuan
}

// Bersih menyatakan tidak ada satu pun temuan.
func (l LaporanRekonsiliasi) Bersih() bool { return len(l.Temuan) == 0 }

// atributKlaim adalah medan yang di bentuk baru disimpan SEKALI di tingkat
// klaim, padahal di tabel datar ia berulang di setiap baris peserta.
//
// Bila nilainya berbeda antar baris dalam satu klaim, migrasi MELAPORKANNYA dan
// tidak diam-diam memilih salah satu.
//
// ⚠️ Dokumen proyek menyebut "14 atribut polis berulang" tanpa mengurutkannya.
// Yang dibandingkan di sini adalah atribut yang bentuk barunya memang menyimpan
// di tingkat klaim - bukan tebakan atas daftar yang tidak pernah ditulis.
var atributKlaim = []string{"NO_CLAIM", "POLICY_NO", "BUSINESSNAME", "CLAIM_RETRO"}

func ambilAtributKlaim(b BarisLama, nama string) string {
	switch nama {
	case "NO_CLAIM":
		return b.NO_CLAIM
	case "POLICY_NO":
		return b.POLICY_NO
	case "BUSINESSNAME":
		return b.BUSINESSNAME
	case "CLAIM_RETRO":
		return b.CLAIM_RETRO
	}
	return ""
}

// BongkarBarisLama mengubah baris tabel datar menjadi pohon klaim.
//
// Pengelompokan: CASEID menentukan klaim, CERTIFICATE_NO menentukan peserta di
// dalam klaim itu. Setiap baris lama menjadi SATU baris adjustment - sehingga
// dua baris untuk peserta yang sama menghasilkan dua baris adjustment, bukan
// satu yang menimpa yang lain (ADR-U-0011).
func BongkarBarisLama(baris []BarisLama) ([]models.PohonKlaim, LaporanRekonsiliasi) {
	lap := LaporanRekonsiliasi{BarisMasuk: len(baris)}

	// urutan klaim dan peserta dijaga stabil supaya hasilnya dapat diulang
	urutKlaim := []string{}
	perKlaim := map[string][]BarisLama{}
	for _, b := range baris {
		if strings.TrimSpace(b.CASEID) == "" {
			lap.Temuan = append(lap.Temuan, Temuan{
				Jenis: TemuanIdentitasKosong, Sumber: b.ID, Medan: "CASEID",
				Catatan: "baris tidak dapat digantung ke klaim mana pun; tidak dipindahkan",
			})
			continue
		}
		if _, ada := perKlaim[b.CASEID]; !ada {
			urutKlaim = append(urutKlaim, b.CASEID)
		}
		perKlaim[b.CASEID] = append(perKlaim[b.CASEID], b)
	}

	var pohon []models.PohonKlaim
	for _, caseID := range urutKlaim {
		rows := perKlaim[caseID]
		lap.Temuan = append(lap.Temuan, periksaAtributKlaim(caseID, rows)...)

		p := models.PohonKlaim{
			Work: models.WorkClaim{
				CaseID:       caseID,
				Lini:         inti.LiniLife,
				Type:         rows[0].TYPE,
				CreateOpName: rows[0].CREATEOPNAME,
			},
			Klaim: models.Klaim{
				NomorKlaim: rows[0].NO_CLAIM,
				NomorPolis: rows[0].POLICY_NO,
				NamaBisnis: rows[0].BUSINESSNAME,
				// Mata uangnya datang dari baris adjustment: header warisan
				// tidak punya kolom mata uang sendiri (butir w2).
				ClaimRetro: uang.Money{Currency: rows[0].CURRENCY},
			},
		}

		// ⛔ Mata uang header diambil dari baris PERTAMA, dan itu hanya benar
		// bila seluruh baris klaim menyepakatinya. CURRENCY bukan anggota
		// atributKlaim - jadi tanpa pemeriksaan ini, klaim bermata-uang campur
		// akan mendapat mata uang sembarang TANPA ada yang tahu. Uang yang
		// mata uangnya dipilih diam-diam adalah persis hal yang ADR-F-0004
		// dan seluruh disiplin uang proyek ini cegah.
		mataUang := map[string]bool{}
		for _, b := range rows {
			mataUang[strings.TrimSpace(b.CURRENCY)] = true
		}
		if len(mataUang) > 1 {
			lap.Temuan = append(lap.Temuan, Temuan{
				Jenis: TemuanAtributBerbeda, Sumber: caseID, Medan: "CURRENCY",
				Catatan: fmt.Sprintf("baris satu klaim memakai %d mata uang berbeda; "+
					"mata uang ClaimRetro diambil dari baris pertama", len(mataUang)),
			})
		}

		// CLAIM_RETRO adalah UANG di tingkat header (butir w2). Diurai dengan
		// pola yang sama seperti jumlah klaim: yang gagal DILAPORKAN, tidak
		// ditebak dan tidak membuat proses berhenti.
		if teks := strings.TrimSpace(rows[0].CLAIM_RETRO); teks != "" {
			d, err := utils.ParseDecimal(teks)
			if err != nil {
				lap.Temuan = append(lap.Temuan, Temuan{
					Jenis: TemuanUangTakTerurai, Sumber: rows[0].ID, Medan: "CLAIM_RETRO",
					Catatan: fmt.Sprintf("nilai %q tidak dapat diurai: %v", teks, err),
				})
			} else {
				p.Klaim.ClaimRetro.Amount = d
			}
		}

		urutPeserta := []string{}
		perPeserta := map[string][]BarisLama{}
		for _, b := range rows {
			kunci := b.CERTIFICATE_NO
			if strings.TrimSpace(kunci) == "" {
				lap.Temuan = append(lap.Temuan, Temuan{
					Jenis: TemuanIdentitasKosong, Sumber: b.ID, Medan: "CERTIFICATE_NO",
					Catatan: "peserta tanpa nomor sertifikat; dikelompokkan sendiri per baris",
				})
				kunci = "\x00baris:" + b.ID
			}
			if _, ada := perPeserta[kunci]; !ada {
				urutPeserta = append(urutPeserta, kunci)
			}
			perPeserta[kunci] = append(perPeserta[kunci], b)
		}

		for _, kunci := range urutPeserta {
			br := perPeserta[kunci]
			ps := models.Peserta{
				NomorPremiList:  br[0].PL_NUMBER,
				NomorPolis:      br[0].POLICY_NO,
				NomorSertifikat: br[0].CERTIFICATE_NO,
				MataUang:        br[0].CURRENCY,
				Baris:           make([]models.BarisAdjustment, 0, len(br)),
			}
			for _, b := range br {
				adj, temuan := barisAdjustmentDari(b)
				lap.Temuan = append(lap.Temuan, temuan...)
				if b.STS_REJECT == kontrak.KodeOutstanding ||
					strings.TrimSpace(b.ACCEPTATION_DATE) != "" {
					lap.BarisHardcode++
				}
				ps.Baris = append(ps.Baris, adj)
				lap.AdjustmentTerbentuk++
			}
			p.Klaim.Peserta = append(p.Klaim.Peserta, ps)
			lap.PesertaTerbentuk++
		}
		pohon = append(pohon, p)
		lap.KlaimTerbentuk++
	}

	if lap.BarisHardcode > 0 {
		lap.Temuan = append(lap.Temuan, Temuan{
			Jenis: TemuanNilaiHardcode,
			Medan: "ACCEPTATION_DATE, STS_REJECT",
			Nilai: fmt.Sprintf("%d baris", lap.BarisHardcode),
			Catatan: "INSERT warisan meng-hardcode ACCEPTATION_DATE = SYSDATE dan " +
				"STS_REJECT = 0; keduanya tidak di-bind. Nilai ini menggambarkan waktu " +
				"insert dan bukan keadaan sebenarnya. Migrasi membawanya apa adanya dan " +
				"TIDAK mengarang tanggal akseptasi.",
		})
	}
	return pohon, lap
}

// periksaAtributKlaim melaporkan atribut tingkat klaim yang berbeda antar baris
// peserta. Ia TIDAK memilih salah satu.
func periksaAtributKlaim(caseID string, rows []BarisLama) []Temuan {
	var out []Temuan
	for _, nama := range atributKlaim {
		nilai := map[string]bool{}
		for _, b := range rows {
			nilai[ambilAtributKlaim(b, nama)] = true
		}
		if len(nilai) <= 1 {
			continue
		}
		daftar := make([]string, 0, len(nilai))
		for v := range nilai {
			daftar = append(daftar, v)
		}
		sort.Strings(daftar)
		out = append(out, Temuan{
			Jenis: TemuanAtributBerbeda, Sumber: "CASEID=" + caseID, Medan: nama,
			Nilai:   strings.Join(daftar, " | "),
			Catatan: "nilai berbeda di dalam satu klaim; tidak dipilih sendiri",
		})
	}
	return out
}

// barisAdjustmentDari mengubah satu baris lama menjadi satu baris adjustment.
func barisAdjustmentDari(b BarisLama) (models.BarisAdjustment, []Temuan) {
	var temuan []Temuan
	adj := models.BarisAdjustment{
		ID:             b.ID,
		KodeStatus:     b.STS_REJECT,
		NomorAkseptasi: b.NO_ACCEPTATION,
		// Ketiga kolom bank pindah dari nama warisannya.
		// NAME_OF_BANK tetap, IDBANK -> ID_BANK, ACCOUNTNO -> ACCOUNT_NO.
		// Dipindah apa adanya sebagai teks: nomor rekening berawalan nol
		// adalah hal biasa, dan mengubahnya menjadi angka menghilangkan nol
		// itu (ADR-U-0022).
		NamaBank:      b.NAME_OF_BANK,
		IDBank:        b.IDBANK,
		NomorRekening: b.ACCOUNTNO,
		JumlahKlaim:   uang.Money{Currency: b.CURRENCY},
	}

	if teks := strings.TrimSpace(b.CLAIM_AMOUNT); teks != "" {
		d, err := utils.ParseDecimal(teks)
		if err != nil {
			temuan = append(temuan, Temuan{
				Jenis: TemuanUangTakTerurai, Sumber: b.ID, Medan: "CLAIM_AMOUNT",
				Nilai: teks, Catatan: "dibiarkan kosong; tidak dibulatkan dan tidak ditebak",
			})
		} else {
			adj.JumlahKlaim.Amount = d
		}
	}

	if teks := strings.TrimSpace(b.ACCEPTATION_DATE); teks != "" {
		t, err := utils.ParseTanggal(teks)
		if err != nil {
			temuan = append(temuan, Temuan{
				Jenis: TemuanTanggalTakTerurai, Sumber: b.ID, Medan: "ACCEPTATION_DATE",
				Nilai: teks, Catatan: "dibiarkan kosong; tidak ditebak",
			})
		} else {
			adj.TanggalAkseptasi = t
		}
	}
	return adj, temuan
}

// BarisLamaDari menyusun baris tabel datar dari pohon klaim.
//
// Ini arah SEBALIKNYA dari BongkarBarisLama, dan ia diperlukan karena sistem
// baru menulis DUA tempat: tabel relasional baru DAN INSERT datar ke
// OS_AKSEPTASI_KLAIM_LIFE, sebab hilir masih membaca dari sana.
// Yang dibuang hanya blob JSON.
//
// Satu baris adjustment menjadi satu baris datar, sehingga pulang-pergi
// pohon -> datar -> pohon mempertahankan cacah barisnya.
func BarisLamaDari(p models.PohonKlaim) []BarisLama {
	var out []BarisLama
	for _, ps := range p.Klaim.Peserta {
		for _, adj := range ps.Baris {
			out = append(out, BarisLama{
				ID:             adj.ID,
				CASEID:         p.Work.CaseID,
				NO_CLAIM:       p.Klaim.NomorKlaim,
				POLICY_NO:      p.Klaim.NomorPolis,
				BUSINESSNAME:   p.Klaim.NamaBisnis,
				CERTIFICATE_NO: ps.NomorSertifikat,
				PL_NUMBER:      ps.NomorPremiList,
				CURRENCY:       adj.JumlahKlaim.Currency,
				CLAIM_AMOUNT:   utils.FormatDecimal(adj.JumlahKlaim.Amount),
				// Nilai header disalin ke TIAP baris datar - begitulah tabel
				// warisan menyimpannya, dan rekonsiliasi memang menuntut
				// seluruh baris satu klaim menyepakatinya (atributKlaim).
				CLAIM_RETRO:    utils.FormatDecimal(p.Klaim.ClaimRetro.Amount),
				STS_REJECT:     adj.KodeStatus,
				NO_ACCEPTATION: adj.NomorAkseptasi,
				// Tanggal ditulis dalam satu bentuk yang sama dengan yang
				// dapat dibaca kembali ParseTanggal (ADR-U-0022).
				ACCEPTATION_DATE: utils.FormatTanggal(adj.TanggalAkseptasi),
				TYPE:             p.Work.Type,
				CREATEOPNAME:     p.Work.CreateOpName,
				// Kembali ke nama warisannya, arah berlawanan dengan
				// barisAdjustmentDari di atas.
				NAME_OF_BANK: adj.NamaBank,
				IDBANK:       adj.IDBank,
				ACCOUNTNO:    adj.NomorRekening,
			})
		}
	}
	return out
}
