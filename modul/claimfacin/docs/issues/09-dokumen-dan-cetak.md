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

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"**AC 51 · 52** — dokumen"* → kedua AC
> **diralat di spec**:
>
> - **AC 51** — jenis dokumen di XML bukan tiga: **empat** (PLA, Draft DLA, DLA fac / treaty, nota akseptasi
>   `PrintPDFAccep_MultiAksep` — pencetak di *Rule Pega sumber* di atas) **ditambah Claim Face Sheet per lini**. Aliran
>   HTML/PDF-nya **tidak diekspor** (OQ-CFI-20): tahap 1 menulis nomor, data, dan penanda (`IsCFS`, `PrintFaceClaim`,
>   `PlaStatus`, `DLAStatus`, `IsPrintAccept`) seperti XML, **tanpa berkas PDF** (pesan info di layar).
> - **AC 52** — penomoran lewat `inti/backend/penomor` di aplikasi (huruf **K** klaim · **G** PLA treaty · **P** DLA fac ·
>   **S** DLA treaty; periode dari tanggal tutup buku), bukan procedure — ADR-0043 meng-*supersede* ADR-0006.
>
> Penjaga cetak ganda **ada di XML** dan dibangun: `SaveAcceptation` langkah 1 keluar bila `IsPrintAccept = 1`, dan
> tombol Acceptation nonaktif menurut `IsPrintAccept` (`backend/models/dla.go`). "Ditolak" di sini = **tanpa efek**,
> seperti XML, bukan pesan galat.

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | tidak menahan |

## Perintah verifikasi

1. Terbitkan dokumen akseptasi — ⭐ angkanya cocok dengan penyesuaian.
2. Terbitkan ulang untuk akseptasi yang sama — ⛔ **ditolak**.
3. Periksa penanda tercetak — ⭐ terisi sesudah cetak pertama.
