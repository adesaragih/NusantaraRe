# 21: Pulang-pergi dan dua bentuk dokumen

**Status:** ready-for-agent
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

- [ ] **AC 49** — polis proporsional pulang-pergi menghasilkan nilai yang sama
- [ ] **AC 50** — polis non-proporsional pulang-pergi menghasilkan nilai yang sama
- [ ] **AC 51** — urutan baris anak saat dibaca sama dengan urutan nomor urutnya
- [ ] **AC 52** — bentuk daftar angsuran **datar** terurai benar
- [ ] **AC 53** — bentuk daftar angsuran **bersarang** terurai benar
- [ ] **AC 54** — cakupan yang hanya satu bentuk dinyatakan **tidak memadai**
