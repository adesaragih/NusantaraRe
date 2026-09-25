# E05: Validasi tanggal endorsement

**What to build:** Tanggal endorsement wajib berada **di dalam periode polis** yang di-endors.
Sejumlah nomor polis tertentu dikecualikan dari validasi ini.

`[terverifikasi]` Periode dibaca dari tabel produksi; tanggal di luar rentangnya memunculkan galat.
Pengecualiannya berupa **daftar nomor polis yang tertanam di kode lama**.

⛔ **Nomor polis produksi tidak boleh menjadi literal di kode target.** Ia menjadi **data konfigurasi
yang dapat diaudit** — itu keputusan rancangan, **bukan** perubahan perilaku. Daftar
pengecualiannya tetap sama persis.

**Asal (Pega).** `Activity\CheckEDMPolisDate` (6 langkah) · `RDBList\GetStartDate`

**Keputusan.** K-046 · `CLAUDE.md` §3.5, §4.6

**Blocked by:** E03

**Status:** blocked

- [ ] Tanggal di luar periode polis → galat dengan pesan yang dapat dibaca petugas
- [ ] Daftar pengecualian menjadi **konfigurasi**, bukan literal di kode
- [ ] **Perilaku tidak berubah** — kasus yang dikecualikan sistem lama tetap dikecualikan
- [ ] ⛔ Nomor polis **tidak pernah** muncul di kode, tiket, test, fixture, maupun log
- [ ] Test memakai **nomor sintetis**, bukan nomor produksi
- [ ] **K-046** `K046_ValidasiTanggal_PerbandinganSubstring8Karakter` — sistem lama membandingkan **potongan string tanggal**, bukan objek tanggal; diport apa adanya
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
