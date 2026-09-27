package services

// Kontrak penyerahan baris adjustment ke Komite - tiket 10.
//
// Untuk apa berkas ini: BATAS menuju Komite, bukan isinya. Komite adalah
// sistem luar; yang dibangun di sini gerbang-gerbangnya, muatan yang
// menyeberang, dan rujukan baliknya - bukan roster, bukan keputusan per
// anggota.
//
// Dibaca sesudah: wewenang.go, statusbaris.go.
//
// Istilah:
//   - roster  : baris `POOLDATA.EMAILKOMITE` yang aktif dan menutup nilai klaim.
//   - tingkat : cacah baris roster yang cocok - tinggi tangga komite.
//   - ambang  : nilai klaim MUTLAK yang dipakai mencari roster.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/pkg/utils"
)

// StatusKlaimRoster menyaring roster ke lini Life.
//
// `[terverifikasi]` `Claim Life/Activity/GetListKomiteLife.xml` pecahan baris
// 810-811: `Param.STS_KLAIM = "LIFE"`.
const StatusKlaimRoster = "LIFE"

var (
	// ErrRekeningBelumLengkap - gerbang rekening pembayaran.
	//
	// ⛔ Teksnya PERSIS seperti korpus, bukan parafrase: `[terverifikasi]`
	// `GetListKomiteLife.xml` pecahan baris 449-450 menyetel
	// `local.errmsg = "Name of bank cannot be empty"`, lalu baris 587-589
	// memasangnya ke `.NameOfBank` sebagai `pyMessageLabel`. Pengguna lama
	// mengenali kalimat ini; menerjemahkannya memutus pengenalan itu.
	//
	// ⚠️ Kalimatnya hanya menyebut nama bank, padahal gerbangnya menjaga
	// ketiga medan. Itu bukan kekeliruan saya melainkan bunyi aslinya -
	// dipertahankan apa adanya, dan medan mana yang kosong dibawa terpisah.
	ErrRekeningBelumLengkap = errors.New("Name of bank cannot be empty")
	// ErrBarisBukanOutstanding - hanya baris Outstanding yang diserahkan.
	ErrBarisBukanOutstanding = errors.New(
		"services: hanya baris Outstanding yang dapat diserahkan ke Komite")
	// ErrBarisSudahDiserahkan - baris ber-KOMITE_ID tidak diserahkan dua kali.
	ErrBarisSudahDiserahkan = errors.New(
		"services: baris sudah pernah diserahkan ke Komite")
	// ErrRosterKomiteKosong - penjaga defensif, lihat PeriksaTingkatKomite.
	ErrRosterKomiteKosong = errors.New(
		"services: nol baris roster Komite menutup nilai klaim ini")
	// ErrRosterBelumDiputuskan - sumber roster belum disahkan work owner.
	ErrRosterBelumDiputuskan = errors.New(
		"services: sumber roster Komite belum diputuskan work owner (butir af)")
	// ErrKasusKomiteBelumDiputuskan - penyimpan kasus Komite belum disahkan.
	ErrKasusKomiteBelumDiputuskan = errors.New(
		"services: penyimpan kasus Komite belum diputuskan work owner (butir af)")
	// ErrMataUangKlaimCampur - invariant OQ-060.
	ErrMataUangKlaimCampur = errors.New(
		"services: satu klaim wajib bermata uang tunggal sebelum diserahkan")
)

// PeriksaRekening menolak baris yang rekening pembayarannya belum lengkap.
//
// `[terverifikasi]` `GetListKomiteLife.xml` pecahan baris 655:
// `.NameOfBank==""||.NoAccount==""||.IDOfBank==""`, dengan
// `pyStepsPreCondParamsWhenTrue=2` (LANJUT ke langkah yang memasang galat) dan
// `WhenFalse=3` (LEWATI). Ketiganya setara - satu kosong sudah cukup.
//
// ⚠️ Ini PARITAS, bukan penyimpangan: gerbang yang sama sudah ada di sistem
// lama, dan tiket 10 AC 57 memintanya dipertahankan.
// ⛔ Mengembalikan sentinelnya TELANJANG. Ronde pertama membungkusnya dalam
// `GalatRekening{AdjustmentID, Medan}` - yang `Error()`-nya membuang kedua
// medan itu, dan yang nol pemanggil baca. Pipa yang tidak pernah dibaca
// dibuang, seperti medan `jejak` di tiket 15.
func PeriksaRekening(b models.BarisAdjustment) error {
	for _, isi := range []string{b.NamaBank, b.IDBank, b.NomorRekening} {
		if strings.TrimSpace(isi) == "" {
			return ErrRekeningBelumLengkap
		}
	}
	return nil
}

// PeriksaBolehDiserahkan menjaga keadaan baris sebelum penyerahan.
func PeriksaBolehDiserahkan(b models.BarisAdjustment) error {
	if strings.TrimSpace(b.KomiteID) != "" {
		return fmt.Errorf("%w: baris %q sudah tertaut ke %q",
			ErrBarisSudahDiserahkan, b.ID, b.KomiteID)
	}
	if b.KodeStatus != models.KodeOutstanding {
		return fmt.Errorf("%w: baris %q berkode %q",
			ErrBarisBukanOutstanding, b.ID, b.KodeStatus)
	}
	return nil
}

// AmbangRoster mengubah nilai klaim menjadi ambang pencarian roster.
//
// `[terverifikasi]` `GetListKomiteLife.xml` pecahan baris 790:
// `@if(Local.IsADj<0, Local.IsADj* -1, Local.IsADj)` - nilai MUTLAK. Klaim
// bernilai negatif dicari dengan tandanya dihilangkan, bukan ditolak dan bukan
// dianggap nol.
//
// ⛔ Mata uangnya IKUT. Ambang tanpa mata uang adalah angka telanjang, dan
// angka telanjang yang dibandingkan dengan pita roster adalah persis cara
// pembandingan lintas mata uang terjadi tanpa ada yang sadar.
func AmbangRoster(nilai models.Money) (models.Money, error) {
	if nilai.Kosong() {
		return models.Money{}, fmt.Errorf(
			"%w: nilai klaim kosong, ambang roster tidak dapat dihitung",
			ErrPermintaanTidakSah)
	}
	mutlak := new(apd.Decimal).Set(nilai.Amount)
	mutlak.Negative = false
	return models.Money{Amount: mutlak, Currency: nilai.Currency}, nil
}

// PeriksaTingkatKomite menolak tangga komite tanpa tingkat.
//
// ⚠️ `[keputusan work owner]` Ini PENJAGA DEFENSIF, bukan alur normal: roster
// dijamin >= 1 secara bisnis, sebab pita limitnya berjenjang dan selalu
// menutup nilai klaim.
//
// ⛔ Ia SENGAJA BERBEDA dari Pega. `[terverifikasi]`
// `Claim Life/Activity/CreateKMTLife_Act.xml` pecahan baris 1322-1323
// menyetel `KomiteLoop = @Utilities.SizeOfPropertyList(KomiteList)` dan baris
// 1377-1378 `.TotalKomite` dengan cara yang sama, sedangkan baris 1398-1399
// menyetel `KomiteCount` ke literal `1` - tanpa prasyarat (langkahnya
// ber-`pyStepsPreCondParamsWhen` KOSONG, `WhenTrue=2`/`WhenFalse=2`).
//
// ⚠️ Yang XML NYATAKAN hanya itu: `KomiteCount` konstan 1 sementara dua
// saudaranya mengikuti ukuran daftar. Bahwa akibatnya "tangga satu tingkat
// dari roster kosong" adalah TAFSIR saya, bukan bunyi korpus - `KomiteCount`
// bersebelahan dengan `IndexAdjustment = .pxListSubscript` dan dapat pula
// terbaca sebagai kursor berbasis satu. Yang dijaga di sini tetap: roster
// kosong tidak boleh lewat diam-diam. `[dugaan]` atas tafsirnya.
func PeriksaTingkatKomite(tingkat int) error {
	if tingkat <= 0 {
		return fmt.Errorf("%w: cacah tingkat %d", ErrRosterKomiteKosong, tingkat)
	}
	return nil
}

// PeriksaSatuMataUang menegakkan invariant OQ-060.
//
// ⚠️ Roster `POOLDATA.EMAILKOMITE` `[data DBA]` TANPA kolom mata uang,
// sehingga pita limitnya berlaku atas satu mata uang implisit. Klaim bermata
// uang campur yang diserahkan akan dibandingkan dengan pita yang bukan
// miliknya - tanpa satu pun galat muncul. Karena itu ia dijaga di sini.
func PeriksaSatuMataUang(peserta []models.Peserta) error {
	// ⛔ DUA kolom, dan keduanya dijaga. `CURRENCY_ID` (`CurrencyID`) dan
	// `CURRENCY` (`JumlahKlaim.Currency`) adalah kolom TERPISAH di
	// `T_CLAIMLF_ADJUSTMENT`. Yang ikut ke Komite lewat `AmbangRoster` adalah
	// yang KEDUA.
	//
	// Ronde pertama hanya menjaga yang pertama - sehingga klaim ber-`CURRENCY`
	// campur di balik `CURRENCY_ID` seragam lolos utuh, dan ambangnya
	// dibandingkan dengan pita roster yang bukan miliknya. Persis kegagalan
	// yang komentar fungsi ini sendiri peringatkan.
	for _, kolom := range []struct {
		nama  string
		ambil func(models.BarisAdjustment) string
	}{
		{"CURRENCY", func(b models.BarisAdjustment) string { return b.JumlahKlaim.Currency }},
		{"CURRENCY_ID", func(b models.BarisAdjustment) string { return b.CurrencyID }},
	} {
		pertama, dari := "", ""
		for _, p := range peserta {
			for _, b := range p.Baris {
				nilai := strings.TrimSpace(kolom.ambil(b))
				if nilai == "" {
					continue
				}
				if pertama == "" {
					pertama, dari = nilai, b.ID
					continue
				}
				if nilai != pertama {
					return fmt.Errorf(
						"%w: kolom %s - baris %q bermata uang %q, baris %q bermata uang %q",
						ErrMataUangKlaimCampur, kolom.nama, dari, pertama, b.ID, nilai)
				}
			}
		}
	}
	return nil
}

// SumberRoster mencacah tingkat komite yang menutup sebuah ambang.
//
// `[terverifikasi]` `ReportDefinition/FilterEmailKomiteWithLimit.xml` pecahan
// baris 662 `pyFilterLogic = A AND C AND B`. Ketiga penyaringnya:
// A (670-680) `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM`, **C** (683-697)
// `.STS_KLAIM = Param.STS_KLAIM`, dan **B** (701-715) `.STS_AKTIF = "1"`.
//
// ⚠️ Ronde pertama menukar label B dan C. Tidak mengubah hasil - ketiganya
// disambung AND - tetapi kutipan yang labelnya salah adalah kutipan yang
// salah, dan `<pyLogicLabel>` di dalam rentang yang saya kutip sendiri
// menyanggahnya.
//
// ⛔ Antarmuka, bukan query. `POOLDATA.EMAILKOMITE` tabel PRODUKSI: membacanya
// menuntut persetujuan manusia, dan butir af masih `[USULAN]`.
type SumberRoster interface {
	AmbilAnggota(ctx context.Context, ambang models.Money, lini string) ([]AnggotaKomite, error)
}

// AnggotaKomite adalah satu anggota tangga yang berhak memutuskan.
//
// ⛔ RALAT RANCANGAN A2. Antarmuka ini semula mengembalikan CACAH saja - dan
// cacah tidak dapat membangun tangga: `CreateKMTLife_Act` mengisi
// `KomiteList(<APPEND>)` dengan `KomiteID`, `IDKomite`, dan `KomiteEmail` per
// anggota. Yang kurang baru terlihat ketika penulis kasusnya ditulis.
//
// ⚠️ `Email` DATA ORANG: ia menyeberang ke pengirim email dan ke
// `T_KOMITE_KOMITELIST`, dan tidak ke mana pun lagi.
type AnggotaKomite struct {
	// KomiteID diisi dari `OPERATOR_ID` roster.
	//
	// ⛔ NAMA KOLOM KORPUS MENIPU KE DUA ARAH, dan ini bacaan bukan tebakan:
	// `[terverifikasi]` `CreateKMTLife_Act.xml` pecahan 866/972
	// `.KomiteID = .OPERATOR_ID`, dan 952/1041 `.IDKomite = .JABATAN`.
	// Medan di bawah karena itu dinamai menurut ISInya, dan pemetaan ke nama
	// kolom terjadi di satu tempat - penulis kasus komite.
	KomiteID string
	// IDKomite diisi dari `JABATAN` roster.
	IDKomite string
	Email    string
	// Urut adalah `DEGREE` - jenjang tangga, menaik.
	Urut int
}

// RosterBelumDiputuskan gagal terang selama butir af belum disahkan.
type RosterBelumDiputuskan struct{}

// AmbilAnggota selalu gagal, dan menyebut apa yang ditunggu.
func (RosterBelumDiputuskan) AmbilAnggota(context.Context, models.Money, string) ([]AnggotaKomite, error) {
	return nil, ErrRosterBelumDiputuskan
}

// MuatanKomite adalah yang menyeberang ke Komite.
//
// ⭐ Tiga medannya TIDAK ada di sistem lama - `[keputusan work owner]`,
// AC 24: nilai klaim, mata uangnya, dan status baris SAAT penyerahan. Tanpa
// ketiganya Komite harus membaca balik ke sistem ini untuk tahu apa yang
// sedang ia putuskan.
//
// ⚠️ Bentuk muatan ini KONTRAK LINTAS KONTEKS. Mengubahnya bukan perubahan
// lokal: ia menuntut kesepakatan dengan konteks Komite Claim Life lebih dulu.
type MuatanKomite struct {
	KlaimID      string
	AdjustmentID string
	// Lini dan Type diperlukan work object anaknya (`pxAddChildWork`).
	Lini string
	Type string
	// Anggota adalah tangga yang roster tentukan, urut menaik.
	//
	// ⛔ Ia ADA di muatan sebab `CreateKMTLife_Act` menulis tangganya BERSAMA
	// kasusnya - satu transaksi. Menyerahkan cacahnya saja membuat penulis
	// kasus harus membaca roster untuk kedua kalinya, dan dua bacaan dapat
	// berbeda.
	Anggota       []AnggotaKomite
	JumlahKlaim   models.Money
	KodeStatus    string
	TingkatKomite int
	AkunID        string
	Waktu         time.Time
}

// PembuatKasusKomite melahirkan kasus anak Komite dan mengembalikan ID-nya.
//
// ⛔ Antarmuka, bukan penulis. Tabel `T_GENERAL_KOMITE` dan
// `T_KOMITE_KOMITELIST` milik konteks Komite Claim Life, dan butir af yang
// mengesahkannya masih `[USULAN]`.
type PembuatKasusKomite interface {
	Buat(ctx context.Context, tx *repository.Tx, m MuatanKomite) (string, error)
}

// KasusKomiteBelumDiputuskan gagal terang selama butir af belum disahkan.
type KasusKomiteBelumDiputuskan struct{}

// Buat selalu gagal, dan menyebut apa yang ditunggu.
func (KasusKomiteBelumDiputuskan) Buat(context.Context, *repository.Tx,
	MuatanKomite) (string, error) {
	return "", ErrKasusKomiteBelumDiputuskan
}

// Penyerahan menyerahkan baris adjustment ke Komite.
type Penyerahan struct {
	svc      *Service
	roster   SumberRoster
	kasus    PembuatKasusKomite
	jejak    Jejak
	penyalur *Penyalur
	terakhir HasilSalur
}

// Komite menyusun layanan itu dengan ketiga ketergantungan yang gagal terang.
func (s *Service) Komite() *Penyerahan {
	return &Penyerahan{
		svc:    s,
		roster: RosterBelumDiputuskan{},
		kasus:  KasusKomiteBelumDiputuskan{},
		jejak:  JejakBelumDiputuskan{},
		// Lingkungannya datang dari Service, yang menerimanya sekali dari
		// `config.IsPegaProd`. Bawaan Service sendiri BUKAN produksi.
		penyalur: NewPenyalur(s.lingkungan, AntreanBelumDiputuskan{},
			EfekKeluarClaimLife(ResolverBelumDiputuskan{})...),
	}
}

// DenganRoster mengganti sumber rosternya.
func (p *Penyerahan) DenganRoster(r SumberRoster) *Penyerahan {
	salin := *p
	salin.roster = r
	return &salin
}

// DenganKasus mengganti pembuat kasus komitenya.
func (p *Penyerahan) DenganKasus(k PembuatKasusKomite) *Penyerahan {
	salin := *p
	salin.kasus = k
	return &salin
}

// DenganPenyalur mengganti penyalur efek keluarnya.
func (p *Penyerahan) DenganPenyalur(s *Penyalur) *Penyerahan {
	salin := *p
	salin.penyalur = s
	return &salin
}

// DenganJejak mengganti perekam jejaknya.
func (p *Penyerahan) DenganJejak(j Jejak) *Penyerahan {
	salin := *p
	salin.jejak = j
	return &salin
}

// Serahkan menyerahkan satu baris adjustment ke Komite.
//
// Meniru `Claim Life/Activity/CreateKMTLife_Act.xml` beserta
// `GetListKomiteLife.xml` yang mendahuluinya, tetapi TIDAK meniru bentuk
// penautannya: `[terverifikasi]` Pega menaut lewat `Call pxAddChildWork`
// (pecahan baris 1543) beserta `IndexAdjustment = .pxListSubscript` (baris
// 1440-1441) - indeks POSISI. Di sini penautannya `KOMITE_ID`, ID stabil;
// penyimpangan sadar (AC 62, tiket 14).
//
// ⛔ Urutan gerbangnya bukan gaya. Identitas dulu, lalu wewenang, baru apa pun
// dibaca: permintaan tanpa wewenang tidak berhak tahu apakah klaimnya ada.
// Sesudahnya keadaan baris, lalu rekening, baru roster - sebab memanggil
// roster untuk baris yang tidak boleh diserahkan adalah pekerjaan sia-sia yang
// menyentuh tabel produksi.
func (p *Penyerahan) Serahkan(ctx context.Context, pelaku Pelaku,
	klaimID, pesertaID, adjID string, saat time.Time) error {

	if err := WajibIdentitas(pelaku); err != nil {
		return err
	}
	if strings.TrimSpace(klaimID) == "" || strings.TrimSpace(pesertaID) == "" ||
		strings.TrimSpace(adjID) == "" {
		return fmt.Errorf("%w: pengenal klaim, peserta, dan baris wajib diisi",
			ErrPermintaanTidakSah)
	}
	if !p.svc.PunyaDatabase() {
		return repository.ErrTanpaOracle
	}

	// ⛔ BUTIR bb: kasus yang sudah ditutup tidak dapat diubah lagi.
	// Satu pintu untuk seluruh rute pengubah - lihat
	// services.PastikanKasusTerbuka, yang pemanggilannya ditagih penjaga
	// statik. Diperiksa SESUDAH wewenang: pemanggil yang tidak berhak tidak
	// berhak pula tahu keadaan kasusnya.
	if err := p.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		return err
	}

	baca := repository.NewKlaimLife(p.svc.db)

	// ⛔ Wewenangnya bergantung Type, dan Type hanya ada di baris klaim. Ini
	// pemanggil PRODUKSI pertama WajibWewenangKomite - tiket 07 melahirkannya
	// tanpa pemanggil, dan AC-nya sengaja dibiarkan terbuka sampai di sini.
	tipe, err := baca.TypeKlaim(ctx, klaimID)
	if err != nil {
		return err
	}
	if err := WajibWewenangKomite(pelaku, tipe); err != nil {
		return err
	}

	perBaris, err := baca.AmbilBaris(ctx, klaimID)
	if err != nil {
		return err
	}
	baris, ada := cariBarisPeserta(perBaris, pesertaID, adjID)
	if !ada {
		return fmt.Errorf("%w: baris %q bukan milik peserta %q pada klaim %q",
			ErrPermintaanTidakSah, adjID, pesertaID, klaimID)
	}
	if err := PeriksaSatuMataUang(pesertaDariPeta(perBaris)); err != nil {
		return err
	}
	if err := PeriksaBolehDiserahkan(baris); err != nil {
		return err
	}
	if err := PeriksaRekening(baris); err != nil {
		return err
	}

	ambang, err := AmbangRoster(baris.JumlahKlaim)
	if err != nil {
		return err
	}
	anggota, err := p.roster.AmbilAnggota(ctx, ambang, StatusKlaimRoster)
	if err != nil {
		return err
	}
	// ⛔ Cacahnya DITURUNKAN dari anggotanya, bukan dibaca terpisah: dua
	// sumber untuk satu angka akhirnya berbeda.
	tingkat := len(anggota)
	if err := PeriksaTingkatKomite(tingkat); err != nil {
		return err
	}

	// Kelahiran kasus komite, penautan baris, dan jejaknya berada dalam SATU
	// transaksi. Penunjuk dua arah yang ditulis di dua transaksi dapat berakhir
	// menunjuk sebelah pihak saja.
	if err := p.svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		komiteID, err := p.kasus.Buat(ctx, tx, MuatanKomite{
			KlaimID:       klaimID,
			AdjustmentID:  adjID,
			Lini:          models.LiniLife,
			Type:          tipe,
			Anggota:       anggota,
			JumlahKlaim:   baris.JumlahKlaim,
			KodeStatus:    baris.KodeStatus,
			TingkatKomite: tingkat,
			AkunID:        pelaku.AkunID,
			Waktu:         saat,
		})
		if err != nil {
			return err
		}
		if err := baca.PerbaruiKomiteID(ctx, tx, adjID, komiteID); err != nil {
			return err
		}
		return p.jejak.Rekam(ctx, tx, CatatanJejak{
			AdjustmentID: adjID,
			KlaimID:      klaimID,
			Dari:         baris.KodeStatus,
			Ke:           "diserahkan ke Komite " + komiteID,
			AkunID:       pelaku.AkunID,
			Waktu:        saat,
		})
	}); err != nil {
		return err
	}

	// ⛔ Efek keluar berjalan SESUDAH transaksinya selesai, dan hasilnya TIDAK
	// menjadi galat penyerahan. `[terverifikasi]` urutan yang sama di
	// `Claim Life/Activity/CreateKMTLife_Act.xml`: `Obj-Save` di pecahan baris
	// 2081, lalu `Call SendEmailKlaimLF` di baris 2213 - email menyusul simpan,
	// bukan mendahuluinya.
	//
	// ⚠️ Konsekuensi yang diterima (ADR-U-0008): penyerahan dapat berhasil
	// sementara email kepada Komite masih tertunda di antrean.
	hasil := p.penyalur.Salurkan(ctx, MuatanEfek{
		KlaimID:      klaimID,
		AdjustmentID: adjID,
		AkunID:       pelaku.AkunID,
		Waktu:        saat,
	})
	// ⛔ Hasilnya TIDAK dibuang. Ronde pertama membuangnya utuh - padahal
	// komentar `GagalDiantre` sendiri berbunyi "dilaporkan, bukan ditelan".
	// Kegagalan yang bahkan tidak dapat DIANTRE hilang untuk selamanya, dan
	// itu kegagalan sistem, bukan keadaan normal.
	//
	// ⚠️ Ia tetap TIDAK menggagalkan penyerahan (ADR-U-0008): yang
	// dikembalikan hasilnya, bukan galat.
	p.terakhir = hasil
	return nil
}

// HasilEfekTerakhir menyebut apa yang terjadi pada efek keluar penyerahan
// terakhir yang dijalankan instance ini.
//
// ⚠️ Dibaca pemanggil yang ingin MELAPORKAN, bukan yang ingin menggagalkan.
func (p *Penyerahan) HasilEfekTerakhir() HasilSalur { return p.terakhir }

// cariBarisPeserta mencari satu baris milik peserta tertentu.
func cariBarisPeserta(perBaris map[string][]models.BarisAdjustment,
	pesertaID, adjID string) (models.BarisAdjustment, bool) {

	for _, b := range perBaris[pesertaID] {
		if b.ID == adjID {
			return b, true
		}
	}
	return models.BarisAdjustment{}, false
}

// pesertaDariPeta menyusun ulang peserta beserta barisnya untuk pemeriksaan
// invariant yang menuntut SELURUH klaim, bukan satu baris.
func pesertaDariPeta(perBaris map[string][]models.BarisAdjustment) []models.Peserta {
	out := make([]models.Peserta, 0, len(perBaris))
	for id, daftar := range perBaris {
		out = append(out, models.Peserta{ID: id, Baris: daftar})
	}
	return out
}

// rosterOracle membaca roster dari `EMAILKOMITE` - butir af, A2.
type rosterOracle struct{ pohon *repository.PohonKlaim }

// RosterKomiteOracle menyusun pembaca roster yang memakai Oracle.
func RosterKomiteOracle(svc *Service) SumberRoster {
	return rosterOracle{pohon: repository.NewPohonKlaim(svc.db)}
}

// AmbilAnggota membaca anggota yang menutup ambang, urut menaik.
//
// ⚠️ Ambangnya diserahkan sebagai TEKS desimal - uang tidak pernah menjadi
// float (ADR-U-0003), dan repository membandingkannya lewat `TO_NUMBER`.
func (r rosterOracle) AmbilAnggota(ctx context.Context, ambang models.Money,
	lini string) ([]AnggotaKomite, error) {

	if ambang.Kosong() {
		return nil, fmt.Errorf("%w: ambang roster kosong", ErrPermintaanTidakSah)
	}
	baris, err := r.pohon.AmbilRosterKomite(ctx, utils.FormatDecimal(ambang.Amount), lini)
	if err != nil {
		return nil, err
	}
	out := make([]AnggotaKomite, 0, len(baris))
	for _, b := range baris {
		out = append(out, AnggotaKomite{
			KomiteID: b.OperatorID,
			IDKomite: b.Jabatan,
			Email:    b.Email,
			Urut:     b.Degree,
		})
	}
	return out, nil
}

// kasusOracle menulis kasus komite beserta tangganya - butir af, A2.
type kasusOracle struct{ pohon *repository.PohonKlaim }

// KasusKomiteOracle menyusun penulis kasus yang memakai Oracle.
func KasusKomiteOracle(svc *Service) PembuatKasusKomite {
	return kasusOracle{pohon: repository.NewPohonKlaim(svc.db)}
}

// Buat melahirkan kasus komite dan mengembalikan pengenalnya.
func (k kasusOracle) Buat(ctx context.Context, tx *repository.Tx,
	m MuatanKomite) (string, error) {

	anggota := make([]repository.AnggotaTangga, 0, len(m.Anggota))
	for _, a := range m.Anggota {
		anggota = append(anggota, repository.AnggotaTangga{
			Urut:       a.Urut,
			OperatorID: a.KomiteID,
			Jabatan:    a.IDKomite,
			Email:      a.Email,
		})
	}
	return k.pohon.BuatKasusKomite(ctx, tx,
		m.KlaimID, m.AdjustmentID, m.Lini, m.Type, anggota, m.Waktu)
}
