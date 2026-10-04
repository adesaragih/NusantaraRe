# 05: Penyesuaian nilai klaim, pembagian dan quota share di atasnya

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 03 (estimasi) · 04 (pembagian)
**Menutup:** AC 34 · 35 · 36 · 37 · 38 *(5 AC)* — US 9 · 10

## Hasil & nilai pengguna

Hari ini Baris penyesuaian **belum punya tempat**, dan pembagian serta quota share di atasnya belum terpisah dari pembagian tingkat item objek.

Sesudah tiket ini, Penilai dapat **mengajukan penyesuaian nilai klaim**, dengan pembagian dan quota share yang melekat **pada penyesuaian itu** — ⭐ terpisah dari pembagian tingkat item objek.

## Area codebase

- Lapisan layanan klaim: penyesuaian
- Pembagian dan quota share tingkat penyesuaian

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Penyesuaian | daftar penyesuaian di dalam item objek, ditambah bahan pertimbangan komite |
| Pembagian atas penyesuaian | dua daftar terpisah, ditulis **langkah bertetangga** di berkas yang sama |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0007** — jejak audit

## Acceptance criteria

- [ ] **AC 34–38** — penyesuaian nilai klaim
- [ ] ⭐ Pembagian tingkat **penyesuaian** terpisah dari pembagian tingkat **item objek**
- [ ] ⭐ Enam medan bahan pertimbangan komite tersimpan **bersama penyesuaian**, bukan tabel sendiri

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **13** | Daftar lokasi di dalam retro fakultatif — larik bersarang | tidak menahan |

## Perintah verifikasi

1. Buat satu penyesuaian, isi pembagian dan quota share-nya — ⭐ keduanya tersimpan terpisah.
2. Hapus penyesuaian — ⭐ keduanya **ikut terhapus**.
3. Bandingkan pembagian tingkat penyesuaian dengan tingkat item objek — ⭐ **baris berbeda**.
