# 16: Kerangka rekonsiliasi eksak — tahap 1

**What to build:** Sebuah pembanding yang menjalankan kasus terekam lewat perhitungan sistem baru dan
menyatakan, tanpa ruang tafsir, apakah hasilnya **identik sampai digit terakhir** dengan sistem lama.

Ini ukuran keberhasilan migrasi, bukan pelengkap. Karena seluruh jalur tulis produksi melewati
prosedur basis data yang isinya baru sebagian kami miliki, **membaca kode saja tidak dapat
membuktikan port-nya benar** — hanya membandingkan keluaran dua sistem atas masukan yang sama yang
bisa.

**Tahap 1 = perhitungan murni atas masukan terekam.** Tidak menyentuh alur, tidak menyentuh jalur
tulis, tidak menunggu ekspor produksi tunggal. Tahap berikutnya (per-modul dengan fixture tangga, lalu
end-to-end) menyusul setelah prasyaratnya tiba.

⛔ **Toleransi ditolak.** Bila urutan operasi dan presisi per-langkah direproduksi apa adanya, nilai
desimal bersifat deterministik dan hasilnya **harus** identik. Selisih sekecil apa pun berarti ada
salah-port — dan toleransi hanya menyembunyikannya, justru pada fase yang dirancang untuk
menemukannya.

⛔ **Kerangka ini tidak pernah menyentuh berkas kasus mentah.** Masukannya **hanya** fixture hasil
tiket 15. Menunjuk langsung ke berkas asal akan menarik data pelanggan ke dalam jalur pengujian.

**Blocked by:** 15, 04, 05, 06, 11

**Status:** ready-for-human — kerangka tahap 1 dibangun 01-10-2026; A29 dikonfirmasi (butir 49)

- [ ] Pembanding membaca **fixture ter-de-identifikasi saja**; tidak ada jalur yang membuka berkas kasus mentah
- [ ] Perbandingan **eksak, tanpa toleransi**; tidak ada parameter epsilon di mana pun
- [ ] Selisih dilaporkan dengan menyebut **kasus, lini bisnis, dan langkah** tempat angkanya mulai berbeda — bukan hanya "tidak cocok"
- [ ] Mencakup lini bisnis yang rumusnya sudah ada (tiket 04–06) dan nilai dasar akseptasi (tiket 11)
- [ ] Kasus yang lini bisnisnya belum punya rumus dilaporkan sebagai **belum tercakup**, bukan lulus
- [ ] Hasilnya dapat dijalankan ulang dan deterministik
- [ ] Dicatat bahwa ini **tahap 1**, dengan prasyarat tahap berikutnya disebut eksplisit

## Comments

### 2026-10-01 — ditahan (agent)

Bergantung pada fixture NB-15 (belum boleh masuk repositori: 9 medan bocor menunggu tinjauan,
`../PERTANYAAN-LANJUTAN.md` butir 4) dan nilai dasar akseptasi NB-11 (ditahan). Yang sudah ada sebagai
rekonsiliasi per lini, tanpa toleransi: PA (4 kasus, tiket 01) dan MBU (93 baris, empat jadi fixture,
tiket 05). Pembanding umum yang menyebut kasus/lini/langkah tempat angka mulai berbeda belum dibangun.

### 2026-10-01 — fixture lima kasus tersedia

`backend/services/premium/testdata/kasus/` (tiket 15): NB FIRE, NB Asuransi Kredit, NB MARINE CARGO,
RNW FIRE, EDM FIRE. Belum dipakai uji rekonsiliasi mana pun.

### 2026-10-01 — kerangka tahap 1 dibangun

Kode: `backend/services/rekonsiliasi/` — `DaftarIzin(isi)`, `BerkasKasus(nama, isi)`, `Ringkas(hasil)`. Data PA/MBU
nyata dipindah dari literal Go ke `backend/services/premium/testdata/daftarizin/premi.json` (diekstrak dengan skrip,
tanpa mengubah angka) dan dibaca bersama oleh `premium.TestRekonsiliasiDaftarIzin` dan kerangka ini.

| Kriteria | Keadaan |
| --- | --- |
| Hanya fixture ter-de-identifikasi | ✅ fungsi menerima isi, bukan jalur; `BerkasKasus` menolak berkas yang gagal `deidentifikasi.SudahBersih` (`TestBerkasMentahDitolak`) |
| Eksak, tanpa toleransi | ✅ kesamaan teks `utils.FormatDecimal`; mutasi "toleransi digit terakhir" tertangkap |
| Selisih menyebut kasus, lini, langkah | ✅ `Hasil{Kasus, Lini, Langkah, SistemLama, SistemBaru}`; langkah dari `premium.AsalRumus` (mis. `CalculatePremiPA_FacIn langkah 6 L1003`). ⚠️ "langkah tempat angka mulai berbeda" hanya setingkat rumus: sistem lama hanya merekam premi akhir, bukan nilai antara |
| Lini berumus (04–06) dan nilai dasar akseptasi (11) | ✅ PA dan MBU; ⏸ nilai dasar akseptasi dilaporkan **belum tercakup** — nilai lama (`CARID2`, `LetterNo`) tidak terekam di fixture |
| Lini tanpa rumus = belum tercakup | ✅ `TestLiniBelumDiportBelumTercakup`, `TestBerkasKasusP5` |
| Deterministik | ✅ `TestDeterministik` |
| Dicatat tahap 1 + prasyarat tahap berikutnya | ✅ komentar paket `rekonsiliasi.go` |

**Hasil sekarang:** 8 cocok (4 PA, 4 MBU) · 10 belum tercakup (5 kasus P-5 × premi + nilai dasar; lini FIRE ×3,
ANEKA, MARINE CARGO dikenali lewat `LiniDariPredikat`, A29) · 0 selisih · 0 galat. Uji mutasi 7/7 (termasuk satu
digit data premi diubah).

### 2026-10-01 — agregat MBU per mata uang ditambahkan (butir 52–53)

`AgregatMBU(isi)` membaca `premium/testdata/daftarizin/mbu_mata_uang.json` (tiket 07) dan melaporkan satu baris per
mata uang per kasus, langkah `FillPremiMBU_FacIn langkah 2.6 CurrencyList.Premium <mata uang>`. Kesamaan **nilai**
desimal eksak (A31; `CurrencyList.Premium` bertipe Decimal, butir 52), bukan teks. Mata uang yang hanya ada di satu
sisi = selisih.

**Hasil sekarang (seluruh kerangka):** 12 cocok (4 PA, 4 MBU coverage, 4 agregat MBU) · 1 selisih (agregat kasus MBU
NB #5, `[pertanyaan terbuka]` di register) · 10 belum tercakup · 0 galat. Uji mutasi kerangka 10/10.


### 2026-10-01 — kelima kasus P-5 tercakup per coverage (tiket 18)

`BerkasKasus` kini menurunkan satu baris per coverage dari jalur pemilih rumus (`jalurCoverage`) dan membandingkannya
sebagai teks persis. **Hasil:** 159 cocok (FIRE 45 + 5 + 4, ANEKA 1, MARINE CARGO 104) · 0 selisih · 0 galat · 5 belum
tercakup (nilai dasar akseptasi, satu per kasus — `CARID2`/`LetterNo` lama tetap tidak terekam). Kasus EDM ditandai
catatan `StatusBusiness 3` (A39). Ditambah sebelumnya: 8 cocok PA/MBU dan agregat MBU per mata uang. Kasus renewal
(`rnw-fire-1`) 4/4 cocok — masukan tiket R08 rnwfacin.
