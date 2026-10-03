# 21: Pulang-pergi dan dua bentuk dokumen

**Status:** selesai — tertahan hanya pihak luar: AC 49–51 uji `db` ditulis, belum dijalankan (K11) *(putaran 2, konsolidasi P10 04-10-2026 — rincian `docs/HASIL-IMPLEMENTASI.md` bab 9; semula: sebagian, implementasi 2026-10-03; awalnya ready-for-agent)*
**Blocked by:** **18** · **19**
**Menutup:** NB AC **49–54** *(6 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` Bab 7

## Hasil & nilai pengguna

Dokumen yang disimpan lalu dibaca kembali menghasilkan nilai yang sama — pada **kedua** bentuk
dokumen, bukan hanya yang kebetulan diuji lebih dulu.

⛔⛔ **Bentuk dokumen tidak seragam.** Daftar angsuran **bersarang** pada satu bentuk dan **datar**
pada bentuk lain; keduanya nyata di sistem lama. Pengurai yang menganggapnya seragam akan patah.

## Yang dibangun

Rangkaian uji pulang-pergi yang mencakup kedua bentuk, dan **aturan kecukupan cakupan**:
⛔ test suite yang hanya mencakup satu bentuk **dinyatakan tidak memadai** dan tidak boleh dihitung
sebagai cakupan.

⚠️ **Pulang-pergi saja tidak cukup.** Bila tulis dan baca sama-sama salah secara simetris, ujinya
tetap hijau — sebagian test karena itu memeriksa **nilai kolom langsung**.

## Batas — yang TIDAK termasuk

⛔ Bentuk dokumen pada endorsemen — tiket **27**.
⛔ Contoh dokumen produksi **tidak disimpan** di dalam repositori; yang disimpan cetakan bentuknya.

## Cara mengujinya

Lewat seam `repository`, dengan data buatan — ⛔ **bukan cuplikan produksi**.

## Acceptance criteria

- [ ] 🟡 **AC 49** — polis proporsional pulang-pergi menghasilkan nilai yang sama
- [ ] 🟡 **AC 50** — polis non-proporsional pulang-pergi menghasilkan nilai yang sama
- [ ] 🟡 **AC 51** — urutan baris anak saat dibaca sama dengan urutan nomor urutnya
- [x] **AC 52** — bentuk daftar angsuran **datar** terurai benar *(putaran 2, pemuat tiket 22: `models/dokumenlama_test.go` `TestPecahDokumenProporsionalDatar`)*
- [x] **AC 53** — bentuk daftar angsuran **bersarang** terurai benar *(putaran 2: `TestPecahDokumenNonProporsionalBersarang`; ditulis juga lewat repository — `TestPemuatLamaNonProporsionalBersarang` bertag `db`, belum dijalankan)*
- [x] **AC 54** — cakupan yang hanya satu bentuk dinyatakan **tidak memadai** *(putaran 2: `TestUjiPemecahMencakupDuaBentuk` gagal bila fixture hanya memuat satu bentuk)*

## Catatan implementasi 2026-10-03

Pulang-pergi diuji di `repository/polis_db_test.go` (`-tags db`) — **belum pernah dijalankan**: env skema
uji tidak tersedia di sesi implementasi. AC 52-54 (dua bentuk dokumen lama) milik pemuat tiket 22.

⛔ **RALAT putaran 2 (03-10-2026).** Bunyi lama: *"AC 52-54 (dua bentuk dokumen lama) milik pemuat tiket 22."*
Bunyi baru: pemuat tiket 22 sudah dibangun; kedua bentuk `ListInstallment` diuji di seam fungsi murni
pemecah (`models.PecahDokumenLama`) dengan fixture fiktif `UJI-` berbentuk contoh 2 (datar) dan contoh 3
(bersarang). `[penyimpangan sadar]` atas *"Lewat seam `repository`"* — lihat tiket 22 RALAT seam uji.
