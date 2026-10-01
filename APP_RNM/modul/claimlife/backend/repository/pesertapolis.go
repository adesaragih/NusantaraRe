package repository

// Pembaca peserta polis dari tabel warisan - tiket 02.
//
// Untuk apa berkas ini: layar Register memilih peserta dari
// POOLDATA.M_LIFE_PREMIUM_DETAIL, tabel warisan milik modul PremiumList Life.
//
// ⛔ TABEL ITU BERISI 66,8 JUTA BARIS. Katalog instance pengembangan
// (26-09-2026) menyebut index pada PL_NUMBER, CERTIFICATE_NO, dan POLICY_NO -
// dan TIDAK ada index yang berawalan EDMSTATUS. Query tanpa penyaring
// ber-index karena itu bukan "agak lambat", melainkan pemindaian penuh atas
// 66 juta baris yang menahan basis data produksi.
//
// Dua aturan yang lahir dari itu, dan dijaga test statik:
//  1. setiap query ke tabel ini menyaring dengan PL_NUMBER atau CERTIFICATE_NO
//  2. setiap query berbatas hasil (FETCH FIRST :n ROWS ONLY)
//
// ⛔ VERSI TERAKHIR (keputusan work owner 01-10-2026, OQ-N14): Endorsement Life menulis versi baru polis
// ke tabel yang SAMA (`PL_NUMBER` tetap, `PL_NUMBER_EDM` = `<polis>/01`, `/02`, …) tanpa mengubah baris versi
// lama. Yang dapat dipilih dan diklaim adalah baris VERSI TERAKHIR tiap sertifikat - bila itu `Delete`/`Batal`,
// sertifikatnya tidak tampil. Penyimpangan sadar dari `RDBList/GetPesertaClaim_sql1.xml` (nol saringan versi
// maupun status). Aturannya SATU tempat: sqlVersiPeserta, peringkatVersi, pilihVersiHidup di bawah.
//
// Dibaca sesudah: pohonklaim.go.

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/uang"
	"nusantarare/modul/claimlife/backend/models"
)

// namaTabelPeserta adalah tabel warisan peserta polis.
const namaTabelPeserta = "M_LIFE_PREMIUM_DETAIL"

// batasHasilBawaan dipakai bila pemanggil tidak menyebut batas.
const batasHasilBawaan = 200

// batasHasilTertinggi menahan permintaan yang batasnya kelewat besar.
const batasHasilTertinggi = 1000

// PesertaPolis membaca calon peserta klaim dari tabel warisan.
type PesertaPolis struct{ db *db.DB }

// NewPesertaPolis menyusun pembacanya.
func NewPesertaPolis(db *db.DB) *PesertaPolis { return &PesertaPolis{db: db} }

// CalonPeserta adalah satu baris hasil pencarian.
//
// Sengaja TIDAK memuat seluruh 85 kolom tabel itu: yang dibawa hanya yang
// dipakai layar untuk memilih, dan ⛔ kolom KTP tidak pernah ikut.
type CalonPeserta struct {
	NomorPremiList  string `json:"nomorPremiList"`
	NomorPolis      string `json:"nomorPolis"`
	NomorSertifikat string `json:"nomorSertifikat"`
	NamaTertanggung string `json:"namaTertanggung"`
	MataUang        string `json:"mataUang"`
	EDMStatus       string `json:"edmStatus"`
}

// penyaringHidup menyingkirkan peserta yang sudah batal atau dihapus lunak.
//
// `[terverifikasi]` verdict V14 grilling Endorsement Life: baris negatif hasil
// jurnal balik ditulis ke tabel yang SAMA dan hidup berdampingan dengan baris
// positifnya - akuntansi melihat keduanya, klaim hanya melihat yang masih
// hidup. Kontraknya berbunyi "peserta yang sudah EDM Batal atau soft-delete
// TIDAK BOLEH MUNCUL di Claim Life".
//
// ⚠️ Penyaring naif `EDMSTATUS = ”` KELIRU dan itu bukan kehati-hatian
// berlebihan: agregat instance pengembangan menghitung 59,1 juta baris
// ber-EDMSTATUS NULL dan NOL baris bernilai teks kosong. Peserta new business
// justru yang NULL, sehingga penyaring naif membuang hampir seluruh tabel.
//
// ⚠️ TRIM di kedua sisi. Tanpa itu, SQL dan PesertaHidup menjadi DUA aturan
// yang berbeda: "Batal " berspasi lolos SQL tetapi ditolak Go, dan AC 29 -
// "penyaringan terjadi di satu tempat" - dilanggar tanpa satu pun test gagal.
const penyaringHidup = `(EDMSTATUS IS NULL OR TRIM(EDMSTATUS) NOT IN ('Batal','Delete'))`

// sqlVersiPeserta - VERSI satu baris peserta: angka sesudah `<PL_NUMBER>/` di `PL_NUMBER_EDM`
// (`Generate_NoEndorsmentLife` b84 `NOPOLIS||'/'||CARI14`); `PL_NUMBER_EDM` kosong, atau tidak berawalan
// nomor polisnya, = versi 0 (new business).
//
// ⚠️ Diperketat dari rumus brief `REGEXP_SUBSTR(PL_NUMBER_EDM, '[0-9]+$')`: rumus itu membaca angka di akhir
// NOMOR POLIS sendiri sebagai versi bila baris new business ber-`PL_NUMBER_EDM` = `PL_NUMBER` (mis.
// `…-2024` → versi 2024, menang atas `/01`). Di sini hanya akhiran sesudah `<PL_NUMBER>/` yang terhitung.
const sqlVersiPeserta = `NVL(TO_NUMBER(CASE WHEN SUBSTR(TRIM(PL_NUMBER_EDM), 1, LENGTH(PL_NUMBER) + 1) = PL_NUMBER || '/' ` +
	`THEN REGEXP_SUBSTR(SUBSTR(TRIM(PL_NUMBER_EDM), LENGTH(PL_NUMBER) + 2), '^[0-9]+$') END), 0)`

// peringkatVersi - peringkat baris dalam SATU sertifikat: versi terbesar, lalu `TGL_INPUT` terbaru (NULL
// terakhir), lalu `ID` terbesar secara ANGKA (`LENGTH` dulu: `ID` teks berisi angka sequence, `'99' > '100'`
// secara teks). ⚠️ Tidak satu pun penulis korpus (`SaveMasterLPDet`, PremiumList, Endorsement) mengisi
// `TGL_INPUT` - di baris mereka seri versi praktis diputus `ID`.
//
// ⛔ Hanya boleh dipakai di subkueri yang `WHERE`-nya sudah `PL_NUMBER = :n` (aturan §1 butir 4): jendelanya
// menyentuh baris SATU polis, nol pemindaian penuh atas 66,8 juta baris.
const peringkatVersi = `ROW_NUMBER() OVER (PARTITION BY CERTIFICATE_NO ORDER BY ` + sqlVersiPeserta +
	` DESC, TGL_INPUT DESC NULLS LAST, LENGTH(ID) DESC, ID DESC)`

// pilihVersiHidup - baris versi terakhir (`RN_VERSI` dari peringkatVersi) yang masih hidup.
const pilihVersiHidup = `RN_VERSI = 1 AND ` + penyaringHidup

// kolomCari - daftar pilih `Find Insured` (kontrak kolom `TestKolomBacaClaimLifeDiisiPenulisPremiumList`).
const kolomCari = `PL_NUMBER, POLICY_NO, CERTIFICATE_NO, NAME_OF_INSURED, CURRENCY, EDMSTATUS`

// sqlCariPeserta menyusun query pencarian peserta beserta nilai bind-nya.
//
// Dipisah dari Cari supaya bentuk SQL-nya dapat diuji tanpa Oracle: kedua
// pagarnya - penyaring ber-index dan batas hasil - adalah hal yang paling
// mahal bila hilang, dan test yang memerlukan basis data tidak pernah jalan
// di mesin pengembang.
//
// Ditiru dari `RDBList/GetPesertaClaim_sql1.xml:85`, dengan SATU penyimpangan
// yang disengaja dan dilaporkan (OQ-E):
//
//	Pega memasang KEDUA `LIKE` tanpa syarat. Di Oracle `X LIKE '%'` bernilai
//	FALSE ketika X NULL, sehingga kotak pencarian yang dibiarkan KOSONG pun
//	diam-diam membuang setiap peserta yang NAME_OF_INSURED-nya NULL. Kita
//	memasang `LIKE` hanya untuk kotak yang terisi, sehingga kotak kosong
//	berarti "jangan saring" - yang memang dibaca orang dari kotak kosong.
//
// ⚠️ Isi kotak TIDAK di-escape dari `%` dan `_`, sama seperti Pega. Pemakai
// yang mengetik `%` memperluas pencariannya sendiri, dan hasilnya tetap
// terkurung PL_NUMBER dan batas hasil.
//
// ⛔ Versi terakhir (OQ-N14): subkueri memeringkat baris SATU polis (`PL_NUMBER = :1`); saringan
// sertifikat ikut di dalamnya (menyaring seluruh kelompok satu sertifikat - peringkat di dalamnya tidak
// berubah), saringan NAMA dan status di luar, atas baris versi terakhir saja: nama lama dari versi yang sudah
// diganti tidak lagi menemukan pesertanya.
func sqlCariPeserta(tabel, nomorPremiList, sertifikat, nama string, batas int) (
	string, []any) {
	dalam := []string{"PL_NUMBER = :1"}
	luar := []string{pilihVersiHidup}
	arg := []any{nomorPremiList}

	if s := strings.TrimSpace(sertifikat); s != "" {
		arg = append(arg, s)
		dalam = append(dalam,
			"CERTIFICATE_NO LIKE '%'||:"+strconv.Itoa(len(arg))+"||'%'")
	}
	if n := strings.TrimSpace(nama); n != "" {
		// b405 `@toUpperCase`. Huruf besarnya dikerjakan di Go, bukan lewat
		// UPPER(:bind) di SQL, supaya nilai yang dikirim dan nilai yang
		// dibandingkan adalah satu hal yang sama dan terlihat di log bind.
		// ⛔ Penampung muncul URUT (`:1`, `:2` di dalam, `:3` di luar): godror
		// mengikat menurut urutan kemunculan.
		arg = append(arg, strings.ToUpper(n))
		luar = append(luar,
			"UPPER(NAME_OF_INSURED) LIKE '%'||:"+strconv.Itoa(len(arg))+"||'%'")
	}

	q := `SELECT ` + kolomCari + `
		        FROM (SELECT ` + kolomCari + `, ` + peringkatVersi + ` AS RN_VERSI
		                FROM ` + tabel + `
		               WHERE ` + strings.Join(dalam, " AND ") + `)
		        WHERE ` + strings.Join(luar, " AND ") + `
		        ORDER BY CERTIFICATE_NO
		        FETCH FIRST ` + strconv.Itoa(batas) + ` ROWS ONLY`
	return q, arg
}

// sqlAmbilPesertaKlaim - baris yang disalin ke klaim untuk SATU sertifikat: baris versi terakhirnya, hanya
// bila hidup (OQ-N14). Subkueri memeringkat baris (PL_NUMBER, CERTIFICATE_NO) itu saja - keduanya
// ber-index (`INDEX4`, `INDEX8`) - lalu barisnya dibaca menurut `ID` terpilih.
// Penampung :1/:3 nomor premium list, :2/:4 sertifikat (urut kemunculan).
func sqlAmbilPesertaKlaim(tabel string) string {
	return fmt.Sprintf(`SELECT %s FROM %s
		        WHERE PL_NUMBER = :1 AND CERTIFICATE_NO = :2
		          AND ID = (SELECT ID FROM (SELECT ID, EDMSTATUS, %s AS RN_VERSI FROM %s
		                                     WHERE PL_NUMBER = :3 AND CERTIFICATE_NO = :4)
		                     WHERE %s)
		        FETCH FIRST 1 ROWS ONLY`, kolomSalin, tabel, peringkatVersi, tabel, pilihVersiHidup)
}

// Cari mengembalikan calon peserta satu premium list.
//
// Batas hasil WAJIB: layar tidak pernah memerlukan 66 juta baris, dan query
// tanpa batas adalah cara paling mudah menahan basis data tanpa sengaja.
//
// `sertifikat` dan `nama` adalah kotak pencarian `Find Insured` layar Register
// (`LoadDataPesertaSpesifik_Act`). Keduanya BOLEH kosong; yang kosong tidak
// memasang penyaring sama sekali - lihat sqlCariPeserta.
func (r *PesertaPolis) Cari(ctx context.Context, nomorPremiList, sertifikat, nama string,
	batas int) ([]CalonPeserta, error) {
	if strings.TrimSpace(nomorPremiList) == "" {
		return nil, fmt.Errorf("repository: mencari peserta tanpa nomor premium list; " +
			"tabel peserta hanya ber-index pada PL_NUMBER, CERTIFICATE_NO, dan POLICY_NO")
	}
	if batas <= 0 {
		batas = batasHasilBawaan
	}
	if batas > batasHasilTertinggi {
		batas = batasHasilTertinggi
	}
	tabel, err := r.db.Qualify(namaTabelPeserta)
	if err != nil {
		return nil, err
	}
	q, arg := sqlCariPeserta(tabel, nomorPremiList, sertifikat, nama, batas)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: mencari peserta %s: %w", nomorPremiList, err)
	}
	defer func() { _ = baris.Close() }()

	var out []CalonPeserta
	for baris.Next() {
		var pl, polis, sertifikat, nama, mataUang, edm sql.NullString
		if err := baris.Scan(&pl, &polis, &sertifikat, &nama, &mataUang, &edm); err != nil {
			return nil, fmt.Errorf("repository: membaca peserta %s: %w", nomorPremiList, err)
		}
		// ⛔ Disaring DUA KALI dengan sengaja, dan itu bukan pemborosan.
		// Penyaring SQL ada supaya Oracle tidak mengirim 66 juta baris;
		// PesertaHidup adalah ATURANNYA, dan ia yang menentukan. Tanpa
		// pemanggil produksi, aturan di Go dan aturan di SQL dapat berselisih
		// diam-diam - dan yang terlihat hanya salah satunya.
		if !PesertaHidup(edm.String) {
			continue
		}
		out = append(out, CalonPeserta{
			NomorPremiList:  pl.String,
			NomorPolis:      polis.String,
			NomorSertifikat: sertifikat.String,
			NamaTertanggung: nama.String,
			MataUang:        mataUang.String,
			EDMStatus:       edm.String,
		})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca peserta %s: %w", nomorPremiList, err)
	}
	return out, nil
}

// PesertaHidup memutuskan apakah satu nilai EDMSTATUS masih boleh tampil.
//
// Dipisah dari SQL supaya aturannya dapat diuji tanpa Oracle, dan supaya
// penyaring itu punya SATU rumah - AC 29 tiket 02 menuntut penyaringan terjadi
// di satu tempat, dan test yang menemukan jalur kedua gagal.
func PesertaHidup(edmStatus string) bool {
	switch strings.TrimSpace(edmStatus) {
	case "Batal", "Delete":
		return false
	default:
		// NULL, teks kosong, Old, New, dan nilai tak dikenal tetap tampil.
		// ⛔ Nilai tak dikenal sengaja LOLOS: menyembunyikan peserta karena
		// status yang belum kita pahami jauh lebih berbahaya daripada
		// menampilkan satu baris berlebih, sebab yang hilang tidak terlihat.
		return true
	}
}

// kolomSalin adalah kolom M_LIFE_PREMIUM_DETAIL yang disalin ke klaim.
//
// ⛔ NAME_OF_INSURED dan POLICY_HOLDER SENGAJA TIDAK ADA di sini. Nama
// tertanggung hanya diperlukan layar saat memilih, dan itu dilayani
// CalonPeserta. Menyalinnya ke tabel klaim berarti menduplikasi data pribadi
// tanpa satu pun AC yang memintanya. ⛔ Kolom KTP tidak pernah dibaca sama
// sekali.
//
// Seluruh nama di sini `[terverifikasi]` ada di katalog instance pengembangan
// (KATALOG-TABEL-PESERTA-DAN-TREATY.md), bukan diturunkan dari nama tabel
// warisan yang mirip.
//
// ⛔ STNC dibungkus TO_CHAR meski namanya tidak berbunyi seperti tanggal:
// katalog menyebutnya DATE (kolom 29). Membacanya apa adanya membuat bentuknya
// bergantung NLS_DATE_FORMAT sesi - jebakan yang sama dengan seluruh tanggal
// lain di proyek ini. ⚠️ [terbuka] Kolom tujuannya di 003 bernama STNC_TREATY
// dan bertipe VARCHAR2(64); sumbernya DATE. Selisih tipe itu belum diputuskan
// siapa pun, dan executor tidak mengubah DDL tanpa keputusan.
const kolomSalin = `ID, PL_NUMBER, POLICY_NO, CERTIFICATE_NO, CURRENCY, ` +
	`TO_CHAR(STNC, 'YYYY-MM-DD HH24:MI:SS'), ` +
	`TO_CHAR(GROSS_VALUATION_BEGIN_DATE, 'YYYY-MM-DD HH24:MI:SS'), ` +
	`TO_CHAR(GROSS_VALUATION_EXPIRED_DATE, 'YYYY-MM-DD HH24:MI:SS'), ` +
	`TO_CHAR(RETRO_VALUATION_BEGIN_DATE, 'YYYY-MM-DD HH24:MI:SS'), ` +
	`TO_CHAR(RETRO_VALUATION_EXPIRED_DATE, 'YYYY-MM-DD HH24:MI:SS'), ` +
	`TO_CHAR(WPC, 'YYYY-MM-DD HH24:MI:SS'), ` +
	`TO_CHAR(BEGIN_DATE, 'YYYY-MM-DD HH24:MI:SS'), ` +
	`TO_CHAR(EFFECTIVE_DATE, 'YYYY-MM-DD HH24:MI:SS'), ` +
	`TO_CHAR(LAPSE_DATE, 'YYYY-MM-DD HH24:MI:SS'), ` +
	`TO_CHAR(EXPIRED_DATE, 'YYYY-MM-DD HH24:MI:SS'), ` +
	`TO_CHAR(SUM_INSURED, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), ` +
	`TO_CHAR(SUM_REASURED, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), ` +
	`TO_CHAR(GROSS_PREMIUM, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), ` +
	`TO_CHAR(NET_PREMIUM, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), ` +
	`TO_CHAR(CEDING_RETENTION, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), ` +
	`TO_CHAR(SHARE_NUSANTARA_RE, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), ` +
	`TO_CHAR(SHARE_RETRO, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), ` +
	`TO_CHAR(RETROCEDED_SHARE, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), ` +
	`TO_CHAR(EM_PERCENT, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), ` +
	// Keempat kolom di bawah ini TIDAK disalin apa adanya: keempatnya
	// bahan bagi dua aturan pemilihan di pilihpeserta.go. Dibaca di sini
	// supaya pilihannya terjadi dalam SATU baris yang sama - membacanya
	// lewat query kedua membuka celah baris berubah di antaranya.
	`TO_CHAR(SHARE_NUSANTARA_RE_GROSS, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), ` +
	`TO_CHAR(AGE, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), ` +
	`TO_CHAR(ENTRY_AGE, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), ` +
	`TO_CHAR(CURRENT_AGE, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), ` +
	// ⭐ GILIRAN-14 butir bp - di EKOR, supaya nol posisi lain bergeser.
	// `SavePesertaClaim` 7.7 b3280 dan 7.8 b3828 menyalin `.CLAIM_AMOUNT`
	// sumber ke peserta DAN ke baris adjustment pertamanya. `[terverifikasi]`
	// kolom 67 katalog M_LIFE_PREMIUM_DETAIL; penulis PremiumList mengisinya
	// (`polis_warisan.go`).
	`TO_CHAR(CLAIM_AMOUNT, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')`

// cacahKolomSalin adalah cacah ekspresi kolomSalin - posisi yang dibaca
// salinKePeserta. Dikunci TestUrutanKolomSalinDikunci.
const cacahKolomSalin = 29

// NamaKolomSalinPeserta mengeluarkan nama kolom yang dibaca kolomSalin.
//
// Diekspor untuk penjaga tiruan skema uji: tiruan yang tertinggal dari
// pembacanya membuat setiap uji db pendaftaran gagal tanpa ada yang tahu.
func NamaKolomSalinPeserta() []string {
	return namaKolomDaftarPilih(kolomSalin)
}

// polaNamaKolomPilih mengenali nama kolom satu butir daftar pilih.
var polaNamaKolomPilih = regexp.MustCompile(`^(?:TO_CHAR\()?\s*([A-Z][A-Z0-9_]*)`)

// namaKolomDaftarPilih mengambil nama kolom dari satu daftar pilih (kolom
// hasil sebuah kueri).
//
// Pemecahnya koma di tingkat teratas; bentuk yang ditemui hanya dua - nama
// telanjang dan `TO_CHAR(NAMA, …)`.
func namaKolomDaftarPilih(daftar string) []string {
	var butir []string
	tingkat, awal := 0, 0
	for i, c := range daftar {
		switch c {
		case '(':
			tingkat++
		case ')':
			tingkat--
		case ',':
			if tingkat == 0 {
				butir = append(butir, daftar[awal:i])
				awal = i + 1
			}
		}
	}
	butir = append(butir, daftar[awal:])
	var keluar []string
	for _, b := range butir {
		if m := polaNamaKolomPilih.FindStringSubmatch(strings.TrimSpace(b)); m != nil {
			keluar = append(keluar, m[1])
		}
	}
	return keluar
}

// AmbilUntukKlaim membaca peserta terpilih dan menyiapkannya untuk disalin.
//
// ⛔ Kenapa server membaca ulang alih-alih memercayai badan permintaan: nilai
// polis - uang, tanggal valuasi, share - menentukan angka klaim dan jendela
// DOL. Menerimanya dari klien berarti siapa pun yang dapat mengirim permintaan
// dapat menentukannya. Klien hanya menyebut NOMOR SERTIFIKAT; sisanya dibaca
// server dari sumbernya.
//
// Penyaringnya sama dengan Cari: ber-index pada (PL_NUMBER, CERTIFICATE_NO),
// berbatas hasil, dan hanya baris VERSI TERAKHIR yang hidup yang ikut - sertifikat yang
// versi terakhirnya batal/delete ditolak dengan pesan "tidak ada" yang sama (OQ-N14).
func (r *PesertaPolis) AmbilUntukKlaim(ctx context.Context, nomorPremiList string,
	sertifikat []string) ([]models.Peserta, error) {
	if strings.TrimSpace(nomorPremiList) == "" {
		return nil, fmt.Errorf("repository: menyalin peserta tanpa nomor premium list")
	}
	if len(sertifikat) == 0 {
		return nil, nil
	}
	tabel, err := r.db.Qualify(namaTabelPeserta)
	if err != nil {
		return nil, err
	}

	var out []models.Peserta
	// Satu sertifikat satu query: keduanya ber-index, dan daftar IN yang
	// panjang membuat Oracle membuat rencana baru untuk tiap panjang daftar.
	for _, no := range sertifikat {
		q := sqlAmbilPesertaKlaim(tabel)
		if err := db.PeriksaSQL(q); err != nil {
			return nil, err
		}
		baris := r.db.QueryRowContext(ctx, q, nomorPremiList, no, nomorPremiList, no)
		sel := make([]sql.NullString, cacahKolomSalin)
		tujuan := make([]any, len(sel))
		for i := range sel {
			tujuan[i] = &sel[i]
		}
		if err := baris.Scan(tujuan...); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf(
					"repository: peserta sertifikat %q pada premium list %q tidak ada "+
						"atau sudah batal/delete", no, nomorPremiList)
			}
			return nil, fmt.Errorf("repository: membaca peserta %q: %w", no, err)
		}
		p, err := salinKePeserta(sel)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

// salinKePeserta menyusun models.Peserta dari satu baris kolomSalin.
func salinKePeserta(sel []sql.NullString) (models.Peserta, error) {
	teks := func(i int) string { return sel[i].String }
	p := models.Peserta{
		SumberID:        teks(0),
		NomorPremiList:  teks(1),
		NomorPolis:      teks(2),
		NomorSertifikat: teks(3),
		MataUang:        teks(4),
		STNC:            teks(5),
		// ⭐ Dipilih pengguna - inilah penyimpan aturan mana peserta yang
		// diklaim (AC 8 tiket 03).
		//
		// ⛔ RALAT 27-09-2026: baris ini menuliskan "1", sedangkan gerbang
		// akseptasi menuntut "true" (`SavePesertaClaim.xml` b3631 langkah 7.7
		// menguji `.IsCheck=="true"`, `SetIndexAdjustmentList.xml` b328 menuliskan
		// `true`). Akibatnya setiap peserta yang baru didaftarkan ditolak
		// saat hendak diaksep. Nilainya kini datang dari satu konstanta.
		IsCheck:             models.PenandaDipilih,
		ValuasiGrossMulai:   teks(6),
		ValuasiGrossSelesai: teks(7),
		ValuasiRetroMulai:   teks(8),
		ValuasiRetroSelesai: teks(9),
		WPC:                 teks(10),
		TanggalMulai:        teks(11),
		TanggalEfektif:      teks(12),
		TanggalLapse:        teks(13),
		TanggalExpired:      teks(14),
		Baris:               []models.BarisAdjustment{},
	}
	// ⛔ Share dipilih SEBAGAI TEKS, sebelum diurai jadi Money. XML
	// membandingkan teksnya dengan "0" (pilihpeserta.go), dan
	// membandingkannya sesudah penguraian berarti membandingkan desimal -
	// aturan yang berbeda, yang menganggap "0.00" nol pula.
	//
	// ⚠️ Hasilnya disimpan di PEUBAH SENDIRI, tidak ditulis balik ke sel[20].
	// Menulis balik berarti fungsi bernama "salin..." diam-diam mengubah
	// slice milik pemanggilnya, dan pemanggil berikutnya membaca nilai yang
	// bukan lagi isi basis data.
	share := sql.NullString{
		String: ShareNusantaraReTeks(teks(20), teks(24)),
		Valid:  true,
	}
	p.Umur = UmurPeserta(teks(25), teks(26), teks(27))

	uang := []struct {
		v  sql.NullString
		ke *uang.Money
	}{
		{sel[15], &p.SumInsured}, {sel[16], &p.SumReasured}, {sel[17], &p.GrossPremium},
		{sel[18], &p.NetPremium}, {sel[19], &p.CedingRetention}, {share, &p.ShareNusantaraRe},
		{sel[21], &p.ShareRetro}, {sel[22], &p.RetrocededShare},
		{sel[28], &p.JumlahKlaim},
	}
	for _, u := range uang {
		m, err := uraiUang(p.NomorSertifikat, "kolom polis", u.v, p.MataUang)
		if err != nil {
			return models.Peserta{}, err
		}
		*u.ke = m
	}
	rasio, err := uraiRasio(p.NomorSertifikat, "EM_PERCENT", sel[23])
	if err != nil {
		return models.Peserta{}, err
	}
	p.EMPercent = rasio
	return p, nil
}
