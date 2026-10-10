# 12: Jejak audit dan kronologi klaim

**Status:** ready-for-agent
**Blocked by:** 01 (registrasi) · 06 (uang)
**Menutup:** AC 68 · 69 · 70 · 71 *(4 AC)* — US 24 · 25 · 26

## Hasil & nilai pengguna

Hari ini Perubahan pada klaim **tidak meninggalkan jejak yang dapat dibaca**, sehingga tidak ada yang dapat menjawab siapa mengubah apa dan kapan.

Sesudah tiket ini, Setiap perubahan penting meninggalkan **catatan kronologi**, dan ⛔ **catatan itu tidak ikut terhapus** ketika klaimnya dihapus.

## Area codebase

- Lapisan layanan: penulisan kronologi
- ⚠️ Perilaku hapus — ⛔ **JANGAN berantai**

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"⚠️ Perilaku hapus — ⛔ **JANGAN berantai**"* →
> **berantai (CASCADE)**: migrasi Claim Prop `532` membuat `T_VIEW_SUGGEST.CLAIM_ID` → `T_GENERAL_CLAIM` **ON DELETE
> CASCADE** (`FK_VS_CLAIM`), dan Claim Fac In memakai tabel bersama itu apa adanya. Modul ini tidak menghapus baris
> `T_GENERAL_CLAIM` (penutupan hanya mengisi `STATUS_WORK`). Pemetaan medan kronologi ke kolom bersama: OQ-CFI-19.

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Kronologi klaim | daftar usulan pandangan pada data klaim; dua penulis |
| Penulis kedua | rule transformasi data penyisip kronologi |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0007** — jejak audit

## Acceptance criteria

- [ ] **AC 68–71** — jejak audit
- [ ] ⭐ Catatan kronologi mencatat **akun**, **jabatan saat itu**, keputusan, dan waktu
- [ ] ⛔⛔ Menghapus klaim **TIDAK menghapus** kronologinya — ⭐ jejak yang ikut terhapus berhenti menjadi jejak

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"⛔⛔ Menghapus klaim **TIDAK menghapus**
> kronologinya"* → **kebalikannya** (CASCADE, Claim Prop `532`) — lihat RALAT 10-10-2026 di *Area codebase*.

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi daftar jabatan | ⚠️ menahan **pencatatan jabatan**, tidak menahan pencatatan akun |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"⚠️ menahan **pencatatan jabatan**, tidak menahan
> pencatatan akun"* → **tidak menahan lagi** (butir 6 = register **spec**, terjawab oleh ADR-0030 + workbasket). Jabatan
> di kronologi diambil dari roster `EMAILKOMITE` FACIN, menggantikan tiga nama orang yang tertulis mati di
> `SetchronologyKlaimFacIn` 1.1 (`PARITAS.md` §8 butir 14); disimpan di `T_VIEW_SUGGEST` menurut pemetaan OQ-CFI-19.
>
> ⛔ Kalimat perintah verifikasi lamanya juga dikutip utuh: *"Ubah satu nilai klaim — ⭐ kronologi bertambah satu
> baris."* dan *"Hapus klaimnya — ⛔ kronologinya **tetap ada**."* → langkah 1 **di luar XML**: kronologi ditulis
> **hanya pada titik-titik yang ada di sistem lama** (spec AC 69), bukan pada setiap perubahan nilai —
> **`[penyimpangan sadar]` DIBUANG, tidak dibangun**. Langkah 2 **terbalik**: kronologi ikut terhapus (CASCADE `532`).

## Perintah verifikasi

1. Ubah satu nilai klaim — ⭐ kronologi bertambah satu baris.
2. Hapus klaimnya — ⛔ kronologinya **tetap ada**.
3. Naikkan jabatan pengubahnya — ⭐ catatan lama **tidak berubah**.
