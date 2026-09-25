# 09: Dokumen akseptasi dan pencetakannya

**Status:** ready-for-agent
**Blocked by:** 05 (penyesuaian) · 06 (uang)
**Menutup:** AC 51 · 52 *(2 AC)* — US 28 · 32

## Hasil & nilai pengguna

Hari ini Dokumen akseptasi **belum dapat terbit**, sehingga penilai tidak punya bukti tertulis atas penyesuaian yang disetujui.

Sesudah tiket ini, Dokumen akseptasi **tercetak dengan angka yang benar**, dan ⭐ **tidak tercetak dua kali** untuk akseptasi yang sama.

## Area codebase

- Lapisan layanan: penerbitan dokumen
- Penanda dokumen sudah tercetak

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pencetakan akseptasi | rule pencetak dokumen akseptasi bermata-uang-banyak |
| Penjaga cetak ganda | penanda dokumen sudah dicetak, diuji di sembilan gerbang |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0007** — jejak audit

## Acceptance criteria

- [ ] **AC 51 · 52** — dokumen
- [ ] ⭐ Dokumen tercetak **sekali** per akseptasi; percobaan kedua **ditolak**

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | tidak menahan |

## Perintah verifikasi

1. Terbitkan dokumen akseptasi — ⭐ angkanya cocok dengan penyesuaian.
2. Terbitkan ulang untuk akseptasi yang sama — ⛔ **ditolak**.
3. Periksa penanda tercetak — ⭐ terisi sesudah cetak pertama.
