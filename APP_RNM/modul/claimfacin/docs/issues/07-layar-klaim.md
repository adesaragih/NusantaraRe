# 07: Layar klaim — apa yang tampil dan apa yang dapat disunting

**Status:** ready-for-agent
**Blocked by:** 01 (registrasi) · 03 (estimasi) · 05 (penyesuaian)
**Menutup:** AC 99 · 100 · 101 · 102 *(4 AC)* — US 6 · 7

## Hasil & nilai pengguna

Hari ini Penilai **belum punya layar** untuk melihat klaim, objeknya, estimasinya, dan penyesuaiannya dalam satu tempat.

Sesudah tiket ini, Penilai melihat klaim **utuh dalam satu layar** — objek, item, estimasi, pembagian, penyesuaian — dan dapat menyunting yang memang boleh disunting.

## Area codebase

- Antarmuka klaim
- Lapisan layanan: pembacaan klaim utuh

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Layar klaim | rule layar klaim fakultatif masuk |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0012** — wewenang eksplisit, bukan efek samping layar

## Acceptance criteria

- [ ] **AC 99–102** — layar
- [ ] ⛔ Layar **tidak menjadi penjaga wewenang** — ia hanya menyembunyikan; penegakan di tiket 08

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **11** | ⚠️ **29 medan bergantung jenis objek** — menjadi kolom, atau dibaca dari polis | ⚠️ **menahan bentuk layar objek** |

## Perintah verifikasi

1. Buka satu klaim — ⭐ objek, item, estimasi, pembagian, dan penyesuaian tampil.
2. Coba sunting medan yang tidak boleh disunting — ⭐ layar menolaknya.
3. ⭐ Panggil lapisan layanan **langsung** untuk menyunting medan itu — ⛔ **juga ditolak** *(tiket 08)*.
