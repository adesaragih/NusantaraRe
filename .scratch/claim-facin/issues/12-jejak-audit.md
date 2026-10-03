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

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi daftar jabatan | ⚠️ menahan **pencatatan jabatan**, tidak menahan pencatatan akun |

## Perintah verifikasi

1. Ubah satu nilai klaim — ⭐ kronologi bertambah satu baris.
2. Hapus klaimnya — ⛔ kronologinya **tetap ada**.
3. Naikkan jabatan pengubahnya — ⭐ catatan lama **tidak berubah**.
