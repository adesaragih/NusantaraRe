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

> ⛔ **RALAT 10-10-2026.** Sel lamanya dikutip utuh, tidak dihapus: *"dua daftar terpisah, ditulis **langkah
> bertetangga** di berkas yang sama"* → benar untuk `CekExGratia` 6 / 8, tetapi **bukan satu-satunya penulis**: yang
> ditiru kode — `Adjustment.SpreadingAdjustment` → `T_CLAIM_ADJ_SPREADING` (`CountTotalEstimasi_Act` 17.2,
> `CheckCurrency_ACT` 4, `CekExGratia` 6, `SetSpreadingAjsutement_Act`); `Adjustment.SpreadingQuotaShare` →
> `T_CLAIM_ADJ_QUOTA_SHARE` (Break QS item bermata uang sama; `ExGratia = 1` mengosongkannya); Break QS **item** →
> `T_CLAIM_BREAK_QS` dari `CheckLimit_Act1` 8. Lihat RALAT 10-10-2026 di `STRUKTUR-TABEL-CLAIM-FACIN.md` §2b
> (`T_CLAIM_BREAK_QS`). Penyesuaian disimpan di `T_CLAIM_ADJUSTMENT` dengan `CLAIM_ID` **dan** `OBJECT_ITEM_ID`
> (OQ-CFI-01, migrasi `566`).

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

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"Daftar lokasi di dalam retro fakultatif — larik
> bersarang"* → ⚠️ **ruang nomor**: butir **13** di sini = register **`STRUKTUR-TABEL-CLAIM-FACIN.md` §6**, bukan butir
> 13 spec (apakah pesan lama menghalangi penyerahan). Isinya tidak berubah.

## Perintah verifikasi

1. Buat satu penyesuaian, isi pembagian dan quota share-nya — ⭐ keduanya tersimpan terpisah.
2. Hapus penyesuaian — ⭐ keduanya **ikut terhapus**.
3. Bandingkan pembagian tingkat penyesuaian dengan tingkat item objek — ⭐ **baris berbeda**.
