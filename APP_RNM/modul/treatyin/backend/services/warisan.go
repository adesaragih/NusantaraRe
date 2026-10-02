package services

// Jalur baca warisan — tiket 40 dan 41.
//
//	40  nilai versi sebelumnya di-SELECT lewat ID_VERSI_KONTRAK_DASAR, BUKAN disalin
//	41  identitas kontrak dalam bentuk lama, DITURUNKAN saat dibaca
//
// ⛔ Keduanya jalur BACA, dan keduanya TIDAK menyimpan apa pun. Itu bukan
// kebetulan: `INV-58` melarang menyimpan turunan dan `INV-59` melarang
// menyalin nilai dari entitas lain. Siapa pun yang kelak menambahkan kolom
// cache di sini melanggar keduanya sekaligus.

import (
	"context"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// NilaiVersiSebelumnya membaca versi yang menjadi DASAR sebuah versi.
//
// ⛔ Dibaca lewat `ID_VERSI_KONTRAK_DASAR`, dan nilainya diambil dari baris
// versi dasar ITU SENDIRI pada saat pembacaan. Sistem lama menyalin seluruh
// halaman clipboard versi dasar ke pohon `OLDDATA` di dalam `JSONDATA` versi
// baru; ketika versi dasarnya dibetulkan, salinannya tidak ikut berubah, dan
// selisih yang dihitung terhadap salinan basi tidak menghasilkan galat - ia
// menghasilkan angka yang salah dengan tenang.
//
// ⚠️ Versi PERTAMA menghasilkan `Ada: false`, BUKAN galat dan BUKAN nilai
// kosong yang menyerupai nol. "Tidak ada versi sebelumnya" adalah jawaban yang
// sah atas pertanyaan ini; memperlakukannya sebagai galat memaksa pemanggilnya
// membedakan galat yang wajar dari galat yang tidak - dan pembedaan itu selalu
// hilang pada pemanggil ketiga.
func (l *Layanan) NilaiVersiSebelumnya(ctx context.Context, p inti.Pelaku, idVersi int64) (models.VersiSebelumnya, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.VersiSebelumnya{}, err
	}
	if idVersi <= 0 {
		return models.VersiSebelumnya{}, fmt.Errorf("%w: pengenal versi %d bukan pengenal yang sah", ErrMasukanTidakSah, idVersi)
	}
	dasar, err := l.gudang.BacaVersiDasar(ctx, idVersi)
	if err != nil {
		return models.VersiSebelumnya{}, err
	}
	if dasar == nil {
		return models.VersiSebelumnya{Ada: false}, nil
	}
	return models.VersiSebelumnya{Ada: true, Versi: dasar}, nil
}

// IdentitasBentukLama menyajikan identitas kontrak dalam bentuk sistem lama.
//
// ⛔ DITURUNKAN dari `KONTRAK` pada saat dibaca - nol tabel, nol kolom, nol
// cache. `ADR-0051` memerintahkan menganggap ada konsumen hilir yang belum
// diketahui: tabel datar `TREATYINDETAIL` ditulis pada setiap penyimpanan di
// sistem lama dan tidak satu pun aturan di dalam ekspor membacanya kembali.
// Diam bukan bukti ketiadaan pembaca; ia hanya bukti pembacanya tidak ada di
// dalam ekspor.
//
// ⚠️ Kontrak yang LAHIR DI SISTEM BARU tidak punya nomor lama, dan jawabannya
// `NomorLamaAda: false` - bukan nomor yang dibangkitkan. Membangkitkannya
// berarti memberi hilir sebuah pengenal yang tidak pernah ada di sistem mana
// pun, dan hilir tidak punya cara membedakannya dari yang asli.
func (l *Layanan) IdentitasBentukLama(ctx context.Context, p inti.Pelaku, idKontrak int64) (models.IdentitasLama, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.IdentitasLama{}, err
	}
	if idKontrak <= 0 {
		return models.IdentitasLama{}, fmt.Errorf("%w: pengenal kontrak %d bukan pengenal yang sah", ErrMasukanTidakSah, idKontrak)
	}
	kv, err := l.gudang.BacaKontrak(ctx, idKontrak)
	if err != nil {
		return models.IdentitasLama{}, err
	}
	k := kv.Kontrak
	if k.ID == 0 {
		return models.IdentitasLama{}, fmt.Errorf("%w: pengenal %d", ErrKontrakTidakAda, idKontrak)
	}
	nomor := strings.TrimSpace(k.NomorKontrakWarisan)
	return models.IdentitasLama{
		NomorLamaAda:    nomor != "",
		NomorLama:       nomor,
		IDKontrak:       k.ID,
		IDCedant:        k.IDCedant,
		IDAsalBisnis:    k.IDAsalBisnis,
		SifatProporsi:   k.SifatProporsi,
		TanggalMulai:    k.TanggalMulai,
		TanggalBerakhir: k.TanggalBerakhir,
	}, nil
}
