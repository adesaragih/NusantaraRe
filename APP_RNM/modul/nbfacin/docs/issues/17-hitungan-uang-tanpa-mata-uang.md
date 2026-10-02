# 17: Hitungan nilai uang tanpa mata uang di pemuat `repository`

**What to build:** Saat pemuat `repository` NB membaca data, ia **menghitung** berapa nilai uang yang tiba
tanpa mata uang dan melaporkannya — tidak memberi mata uang bawaan, tidak berhenti.

Asal: tiket 02 ditutup work owner 01-10-2026 (butir 34, A9) karena digantikan butir 3 (ikuti
`inti/backend/uang` apa adanya). Dua kewajibannya tidak ikut gugur dan dipindah ke sini.
`[terverifikasi]` berkas kasus PA dan MBU nyata menyimpan `Currency` kosong pada coverage (tiket 02 bab
Comments); pengukuran produksi D1 mencatat 4 dari 691.925 baris tanpa mata uang (dikutip tiket 02).

**Blocked by:** pemuat `repository` NB (menunggu tabel flat)

**Status:** needs-info

- [ ] Pemuat menghitung nilai uang tanpa mata uang per medan, dan hitungannya terlihat (log/laporan), tanpa default diam-diam
- [ ] Nilai tanpa mata uang tetap dibawa apa adanya (butir 3), seperti `TestPremiPAMataUangKosongTidakDiberiBawaan`
- [ ] Jalur tulis Oracle gagal keras bila diminta menulis uang tanpa mata uang — kewajiban jalur tulis produksi, tidak dihapus

## Comments
