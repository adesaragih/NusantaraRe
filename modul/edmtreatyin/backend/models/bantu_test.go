package models

// adaPesan - ada setidaknya satu pesan galat terpasang (dulu metode Halaman.AdaPesan; dipindah ke uji 07-10-2026
// karena hanya uji yang memakainya - kode mati dibuang atas perintah work owner).
func adaPesan(h *Halaman) bool {
	if h == nil {
		return false
	}
	for _, p := range h.Pesan {
		if len(p) > 0 {
			return true
		}
	}
	return false
}
