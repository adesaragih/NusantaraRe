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

**Status:** sebagian — 01-10-2026, `backend/services/validasitanggal.go` (`ValidasiTanggalEndorsemen`); test
`validasitanggal_test.go`. Query `GetStartDate` milik repository (E17). Tiga `[dugaan]` dari bentuk
aturannya sendiri — format teks DateTime Pega, `@CompareDates(a,b)` = a sesudah b, `@toDate` = tanggal
kalender — dicatat di kode.

- [x] Tanggal di luar periode polis → "EDM date cannot be outside the period"
- [x] Daftar pengecualian menjadi **konfigurasi** (`DaftarPolis`); pola awalan `RNML` tetap literal (pola, bukan nomor polis)
- [x] **Perilaku tidak berubah** — pengecualian melompat ke label `END` (transisi `1` → `END`, `[terverifikasi]`)
- [x] ⛔ Nomor polis **tidak pernah** muncul di kode, tiket, test, fixture, maupun log — `grep -rniE 'RNM-[A-Z0-9]' APP_RNM/modul/endorsmentfacin` → nol kecocokan (01-10-2026); `CheckEDMPolisDate` memuat 2 kemunculan tetapi hanya 1 nomor unik (prakondisi + deskripsi langkah 2; `grep -oi 'RNM-[A-Z0-9.-]*' CheckEDMPolisDate.xml | sort | uniq -c`) — spec §4.1 menyebut "2 nomor polis literal"
- [x] Test memakai **nomor sintetis**
- [ ] **K-046** perbandingan substring 8 karakter — ⛔ **premis tidak sepenuhnya didukung korpus**: `@substring(…,0,8)` hanya di langkah 4.2 "temporary utk monitoring" (CARI17) dan deskripsi 4.1; keputusan memakai CARI20 (tanggal `dd/MM/yyyy`) dan `@CompareDates`. Kejanggalan yang SEBENARNYA diport: tanggal mulai dibaca dari tanggal kalender **GMT** (sehari lebih awal dari WIB untuk polis 17:00 GMT) — `TestTanggalSamaDenganTanggalMulaiLolos`
- [x] Menyebut rule Pega asalnya dalam komentar (§4.6)
