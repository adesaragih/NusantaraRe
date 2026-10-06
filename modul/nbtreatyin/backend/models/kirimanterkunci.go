package models

// Untuk apa berkas ini: SATU POLA "kiriman terkunci" (F4, keputusan WO
// 04-10-2026: IKUTI XML) - medan yang di layar TIDAK dapat diketik, hanya
// ditulis tombol/popup yang hasilnya DIPEGANG layar sampai Save/Submit
// (refresh tanpa Obj-Save). Server menerima kiriman medan itu hanya bila sama
// dengan salah satu hasil tombol/popup yang DIHITUNG ULANG di server; selain
// itu `GalatKiriman` (services menjawabnya 422 lewat satu jalur).
//
// Dipakai:
//
//	Enable / Disable Input Type  `KirimanEnableDisable` (layar.go)
//	Select Source Of Business    `KirimanSumberBisnis`  (sumberbisnis.go)
//	Choose popup bisnis          `PeriksaPilihanBisnis` (pilihbisnis.go) - id
//	                             permintaan, bukan medan halaman: hanya galatnya

// GalatKiriman - kiriman layar memuat nilai yang tidak mungkin dihasilkan sel
// / tombol / popup layar itu (pola F4). Services menjawabnya 422.
type GalatKiriman struct{ Pesan string }

func (e *GalatKiriman) Error() string { return "models: kiriman layar tidak sah: " + e.Pesan }

// KirimanTerkunci - medan terkunci yang hanya ditulis satu tombol/popup.
// `T` = nilai seluruh medan yang ditulis tombol itu sekaligus.
type KirimanTerkunci[T comparable] struct {
	// Jalur - jalur halaman yang ditulis tombol/popup.
	Jalur []string
	// Tampil - tombolnya tampil di layar saat ini. Tidak tampil = medan
	// TERKUNCI: kiriman diabaikan, nilai halaman server dipakai (pola AC 49-51).
	Tampil bool
	// Baca / Tulis - nilai `T` dari / ke halaman.
	Baca  func(*Halaman) T
	Tulis func(*Halaman, T)
	// Berubah - kiriman memuat pilihan baru dibanding halaman server; nil =
	// `kiriman != server`.
	Berubah func(kiriman, server T) bool
	// Pesan - pesan 422 bila kiriman bukan hasil tombol/popup.
	Pesan func(kiriman T) string
	// HitungUlang - hasil tombol/popup yang dijalankan ulang di server
	// (dipanggil hanya bila kiriman berubah).
	HitungUlang func() ([]T, error)
}

// TerimaKirimanTerkunci menerapkan pola F4 atas halaman server `h`:
//
//   - kiriman = `h` ditimpa nilai `Jalur` yang DIKIRIM layar (`masuk`);
//     jalur yang tidak dikirim tetap bernilai halaman server;
//   - tidak berubah, atau tombol tidak tampil -> tidak berbuat apa pun;
//   - selain itu `HitungUlang` (hasil tombol/popup di server) dijalankan dan
//     kiriman diterima - ditulis ke `h` - hanya bila SAMA dengan salah satu
//     hasilnya; tidak sama -> `GalatKiriman`.
//
// `diterima` = kiriman ditulis ke `h`.
func TerimaKirimanTerkunci[T comparable](h, masuk *Halaman, a KirimanTerkunci[T]) (diterima bool, err error) {
	if masuk == nil || !a.Tampil {
		return false, nil
	}
	gabung := h.Salin()
	dikirim := false
	for _, j := range a.Jalur {
		if v, ada := masuk.Nilai[j]; ada {
			gabung.Setel(j, v)
			dikirim = true
		}
	}
	if !dikirim {
		return false, nil
	}
	kiriman, server := a.Baca(gabung), a.Baca(h)
	berubah := kiriman != server
	if a.Berubah != nil {
		berubah = a.Berubah(kiriman, server)
	}
	if !berubah {
		return false, nil
	}
	hasil, err := a.HitungUlang()
	if err != nil {
		return false, err
	}
	for _, x := range hasil {
		if x == kiriman {
			a.Tulis(h, kiriman)
			return true, nil
		}
	}
	return false, &GalatKiriman{Pesan: a.Pesan(kiriman)}
}
