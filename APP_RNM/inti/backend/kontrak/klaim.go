package kontrak

// Kontrak Claim Life untuk Komite Claim Life - refactor bentuk B.
//
// Komite memutus baris adjustment milik Claim Life, dan keputusan akhirnya
// (aksep/tolak) menulis kembali ke tabel Claim Life DI DALAM transaksi Komite
// sendiri - tidak ada saat ketika tangga Komite sudah tertutup sementara baris
// klaimnya belum. Maka setiap operasi tulis di bawah menerima `*db.Tx` milik
// pemanggil. Disediakan Claim Life (`services.KlaimUntukKomite`), dipakai
// Komite; `cmd/api` yang menyambungnya.

import (
	"context"
	"errors"
	"time"

	"nusantarare/inti/backend/db"
)

// KlaimKomite adalah yang Komite butuhkan dari Claim Life (km3).
type KlaimKomite interface {
	// TypeKlaim membaca `TYPE` kasus klaim (`T_WORK_CLAIM`).
	TypeKlaim(ctx context.Context, klaimID string) (string, error)
	// KodeBisnisKlaim membaca `BUSINESSID` header klaim; teks kosong bila
	// header tidak ada atau kodenya belum tersimpan.
	KodeBisnisKlaim(ctx context.Context, klaimID string) (string, error)
	// PerbaruiStatusBaris memindah satu baris adjustment dari `kodeLama` ke
	// `kode`, beserta nomor dan tanggal akseptasinya.
	PerbaruiStatusBaris(ctx context.Context, tx *db.Tx, pesertaID, adjID, kodeLama, kode,
		nomorAksep string, tglAksep time.Time) error
	// CerminkanHeader menyalin status baris terakhir ke header klaim.
	CerminkanHeader(ctx context.Context, tx *db.Tx, klaimID, kode, nomorAksep string) error
	// CabutPenandaDipilih mencabut penanda `IS_CHECK` peserta.
	CabutPenandaDipilih(ctx context.Context, tx *db.Tx, pesertaID string) error
	// NomorAkseptasiDipakai menjawab apakah nomor akseptasi sudah dipakai
	// baris Claim Life mana pun.
	NomorAkseptasiDipakai(ctx context.Context, tx *db.Tx, nomor string) (bool, error)
	// PastikanKasusTerbuka menolak kasus klaim yang sudah ditutup dengan
	// `ErrKasusSudahTertutup`.
	PastikanKasusTerbuka(ctx context.Context, klaimID string) error
}

// StatusBaris adalah keadaan satu baris AdjustmentList.
//
// Baris AdjustmentList - bukan klaim, bukan peserta - adalah unit keputusan
// mesin status (ADR-U-0011).
type StatusBaris int

const (
	// StatusTidakDiketahui dipakai untuk kolom kosong maupun kode di luar
	// ketiga nilai yang tertulis di spec. Artinya TIDAK ditebak.
	StatusTidakDiketahui StatusBaris = iota
	StatusOutstanding
	StatusAksep
	StatusDitolak
)

// String menulis status sebagai kata yang dibaca pengguna.
//
// Tiket 01 AC-3 / spec.md US-26: status ditampilkan sebagai kata, bukan sebagai
// nama field `STS_REJECT` dan bukan sebagai angka.
func (s StatusBaris) String() string {
	switch s {
	case StatusOutstanding:
		return "Outstanding"
	case StatusAksep:
		return "Aksep"
	case StatusDitolak:
		return "Ditolak"
	default:
		return "Tidak diketahui"
	}
}

// Diketahui membedakan status yang benar-benar tertulis di spec dari yang tidak.
func (s StatusBaris) Diketahui() bool { return s != StatusTidakDiketahui }

// Kode mentah kolom STS_REJECT, ditulis dan dibandingkan sebagai TEKS
// (ADR-U-0022). Menamainya di satu tempat membuat penulisnya dapat dicari:
// nol yang tersebar sebagai literal di banyak berkas tidak dapat ditelusuri.
//
// ⛔ Kode "4" SENGAJA tidak punya nama di sini. Ia ada di data warisan,
// artinya belum diputuskan work owner, dan sistem baru tidak pernah
// menulisnya. Nama akan membuatnya tampak seperti pilihan yang sah.
const (
	KodeOutstanding = "0"
	KodeAksep       = "1"
	KodeDitolak     = "2"
)

// StatusBarisDariKode menerjemahkan nilai kolom `STS_REJECT`.
//
// ⚠️ Namanya menyesatkan: nilai "1" berarti **DIAKSEP**, bukan ditolak.
// spec.md bab Problem butir 3, dan tabel "Penyimpangan sadar 4":
//
//	ReasLifeAdmin insert ke Outstanding -> "0"
//	ReasLifeAdmin reject langsung       -> "2"
//	ReasLifeSPV tambah baris Outstanding-> "0"
//
// dan `[terverifikasi]` nilai "1" = diaksep (spec.md US-26, CONTEXT.md).
//
// Perbandingan dilakukan atas TEKS, tidak pernah lewat bilangan (ADR-U-0022):
// "006" yang dibaca sebagai 6 lolos pulang-pergi dan memecahkan penggolong.
// Karena itu "00" dan "01" BUKAN "0" dan "1".
func StatusBarisDariKode(kode string) StatusBaris {
	switch kode {
	case KodeOutstanding:
		return StatusOutstanding
	case KodeAksep:
		return StatusAksep
	case KodeDitolak:
		return StatusDitolak
	default:
		return StatusTidakDiketahui
	}
}

// StatusWorkSelesai adalah SATU-SATUNYA nilai yang pernah ditulis ke
// `T_WORK_CLAIM.STATUS_WORK`.
//
// ⛔ VERBATIM dari `Flow/Register_Flow.xml` baris 899, di dalam
// `<rowdata REPEATINGINDEX="End1">` (baris 883) yang ber-`pyMOId` End1
// (885) dan `pxObjClass` `Data-MO-Event-End` (901). Sembilan shape lainnya
// ber-`pyWorkStatus` KOSONG, jadi hanya End1 yang menetapkan status kerja.
//
// ⛔ Status Pega untuk kasus yang sedang BERJALAN (`New`, `Open`,
// `Pending-…`) tidak ada di ekspor ini dan TIDAK DIKARANG: kasus terbuka
// berkolom kosong. Kosong berarti "belum ditutup", bukan "tidak diketahui"
// (ADR-U-0027).
const StatusWorkSelesai = "Resolved-Completed"

// KasusTertutup menjawab apakah baris work sudah ditutup.
//
// ⛔ Dibandingkan PERSIS, tanpa merapikan spasi - alasan yang sama dengan
// gerbang di atas: pembanding yang lebih longgar daripada aslinya akan
// memperlakukan kasus yang di sistem lama masih terbuka sebagai tertutup,
// dan menolak setiap perubahan atasnya tanpa jalan keluar yang terlihat.
func KasusTertutup(statusWork string) bool { return statusWork == StatusWorkSelesai }

// approvalAwal adalah nilai `KOMITE_APPROVAL` saat tangga dibentuk.
//
// `[terverifikasi]` `GetListKomiteLife.xml` b1197: `= 0` (disalin
// `CreateKMTLife_Act` b1188). TEKS, bukan bilangan (ADR-U-0022).
const ApprovalKomiteAwal = "0"

// ErrNomorAkseptasiBerganda - nomor rakitan sudah dipakai baris lain.
var ErrNomorAkseptasiBerganda = errors.New(
	"repository: nomor akseptasi sudah dipakai")

var (
	// ErrKodeBisnisBelumTersimpan - kode bisnis tidak ada di model relasional.
	//
	// ⛔ TEMUAN AUDIT A0. Nomor akseptasi memuat kode bisnis
	// (`'RNML-A'||{pyWorkPage.BusinessCode}||…`), tetapi model relasional kita
	// TIDAK menyimpannya: `T_GENERAL_CLAIM` punya `BUSINESS_NAME` saja, dan
	// `BUSINESSID` bukan salah satu dari 18 kolom datar warisan yang `Simpan`
	// tulis. Ia dipakai sekali saat pendaftaran lalu hilang.
	//
	// Gagal terang, bukan dikarang: menebak kode bisnis berarti menerbitkan
	// nomor akseptasi di seri yang salah, dan nomor itu tercetak di dokumen.
	// Kolomnya lahir di A1.
	ErrKodeBisnisBelumTersimpan = errors.New(
		"services: kode bisnis klaim belum tersimpan di model relasional (butir A1)")
)

// ErrKasusSudahTertutup - perubahan atas kasus yang sudah ditutup.
var ErrKasusSudahTertutup = errors.New("services: kasus sudah ditutup")
