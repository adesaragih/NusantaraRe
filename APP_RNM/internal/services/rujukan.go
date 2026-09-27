package services

// Dropdown layar Register - A3 kelompok Register.
//
// Untuk apa berkas ini: gerbang ketiga sumber dropdown. Yang dijaga bukan
// isinya melainkan SIAPA yang boleh membacanya dan seberapa luas.
//
// Dibaca sesudah: pendaftaran.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/internal/repository"
)

// ErrJenisRujukanTidakDikenal - nama dropdown di luar ketiga yang ada.
var ErrJenisRujukanTidakDikenal = errors.New(
	"services: jenis rujukan tidak dikenal")

// Jenis rujukan - ketiganya, dan tidak lebih.
//
// ⛔ Himpunan TERTUTUP, dan namanya dipetakan ke tabel DI KODE - bukan
// diteruskan dari permintaan HTTP ke nama tabel. Nama objek yang datang dari
// pemakai adalah injeksi lewat pintu yang tidak dijaga `:1`.
const (
	RujukanCeding    = "ceding"
	RujukanBisnis    = "bisnis"
	RujukanMarketing = "marketing"
)

// panjangCariMinimum adalah huruf tersedikit sebelum pencarian dijalankan.
//
// ⛔ Pencarian kosong menarik 50 baris pertama menurut abjad - daftar yang
// hampir tidak pernah memuat yang dicari, dan yang menyeret nama orang ke
// layar tanpa ada yang memintanya. Dua huruf adalah ambang termurah yang
// membuat hasilnya bermakna.
const panjangCariMinimum = 2

// Rujukan melayani ketiga dropdown.
type Rujukan struct{ svc *Service }

// SumberRujukan menyusun layanannya.
func (s *Service) SumberRujukan() *Rujukan { return &Rujukan{svc: s} }

// Cari membaca pilihan sebuah dropdown.
//
// ⚠️ Gerbangnya IDENTITAS, bukan peran: ketiga daftar ini dipakai layar
// Register yang sudah bergerbang peran di jalur simpannya, dan menutup
// dropdown per peran akan membuat layar tampak rusak bagi orang yang
// sebenarnya berhak melihatnya. Yang ditolak hanyalah permintaan anonim -
// sebab isinya memuat nama.
func (rj *Rujukan) Cari(ctx context.Context, pelaku Pelaku, jenis, cari string) (
	[]repository.BarisRujukan, error) {

	if err := WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	bersih := strings.TrimSpace(cari)
	if len(bersih) < panjangCariMinimum {
		// ⚠️ Daftar KOSONG, bukan galat: pemakai yang baru mengetik satu huruf
		// belum salah melakukan apa pun.
		return nil, nil
	}
	// ⛔ PERMINTAANNYA diperiksa lebih dulu, baru infrastrukturnya.
	// Ronde pertama memeriksa basis data dulu - dan jenis yang salah
	// ketik dijawab "database belum dikonfigurasi", kalimat yang
	// menyuruh orang memperbaiki hal yang sama sekali berbeda.
	if jenis != RujukanCeding && jenis != RujukanBisnis &&
		jenis != RujukanMarketing {
		return nil, fmt.Errorf("%w: %q", ErrJenisRujukanTidakDikenal, jenis)
	}
	if !rj.svc.PunyaDatabase() {
		return nil, repository.ErrTanpaOracle
	}
	pohon := repository.NewPohonKlaim(rj.svc.db)
	switch jenis {
	case RujukanBisnis:
		return pohon.CariBisnis(ctx, bersih)
	case RujukanMarketing:
		return pohon.CariMarketing(ctx, bersih)
	default:
		// Hanya `RujukanCeding` yang tersisa - sudah dijaga di atas.
		return pohon.CariCeding(ctx, bersih)
	}
}

// BarisRujukan diekspor ulang supaya handlers tidak perlu mengimpor
// repository.
//
// ⛔ Arah ketergantungan `handlers → services → repository` ditegakkan penjaga
// `TestHandlersTidakMengimporRepository`, dan ia menangkap pelanggaran ini
// saat berkas handler-nya pertama ditulis. Alias, bukan salinan struct: tipe
// kembar berarti dua bentuk yang harus berubah bersama dan tidak ada yang
// memaksanya.
type BarisRujukan = repository.BarisRujukan
