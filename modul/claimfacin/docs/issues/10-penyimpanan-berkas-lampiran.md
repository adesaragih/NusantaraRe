# 10: Penyimpanan berkas lampiran

**Status:** ready-for-agent
**Blocked by:** 01 (registrasi)
**Menutup:** AC 53 · 54 · 55 · 56 *(4 AC)* — US 5

## Hasil & nilai pengguna

Hari ini Berkas lampiran klaim — surat, foto, laporan survei — **belum punya tempat simpan**, dan tautannya ke klaim belum ditetapkan.

Sesudah tiket ini, Penilai dapat **melampirkan berkas** pada klaim, dan berkas itu **tetap dapat dibuka** sesudah kasus ditutup.

## Area codebase

- Lapisan layanan: penyimpanan berkas
- Tautan berkas ke klaim

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Penyimpanan berkas | rule unggah dan pengambilan tautan berkas |
| Token penyimpanan | rule pengambil token penyimpanan |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0007** — jejak audit atas unggahan

## Acceptance criteria

- [ ] **AC 53–56** — penyimpanan berkas
- [ ] ⭐ Berkas **tetap terbuka** sesudah kasus ditutup

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"**AC 53–56** — penyimpanan berkas"* → **AC 53
> diralat di spec**: token penyimpanan diterbitkan **aplikasi**, bukan basis data — penyimpanan bersama
> `inti/backend/penyimpanan` + token `inti/backend/layanan` (meniru `GetTokenStorage_SQL`) dengan garam dari env
> **`STORAGE_TOKEN_SALT`**; garam kosong = gagal terang, dan garam tidak pernah masuk pesan galat. Lampiran klaim
> tercatat di **`DOCUMENT_CLAIM`** + **`T_STORAGE_IMAGE`** di transaksi aksi (`backend/services/lampiran.go`,
> `repository/lampiran.go`); kategori dokumen dari master `T_KATEGORI_DOC_KLAIM` baris FAC (Claim Prop `535` / `536`).
> Berkas tetap di Google Storage (AC 94).

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | tidak menahan |

## Perintah verifikasi

1. Lampirkan satu berkas, tutup kasus, buka kembali — ⭐ berkas **masih terbuka**.
2. Lampirkan berkas bernama sama dua kali — ⭐ perilakunya sesuai AC.
3. Periksa jejak audit unggahan — ⭐ tercatat siapa dan kapan.
