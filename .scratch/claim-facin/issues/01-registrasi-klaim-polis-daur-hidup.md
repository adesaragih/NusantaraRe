# 01: Registrasi klaim, pengikatan polis, dan daur hidup kasus

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR)
**Menutup:** AC 4 · 5 · 6 · 7 · 8 · 9 · 10 · 11 · 12 *(9 AC)* — US 1–5

## Hasil & nilai pengguna

Hari ini Sebuah klaim fakultatif masuk **belum bisa dibuat**. Tahapan yang boleh dilaluinya juga belum ditetapkan, sehingga tidak ada yang mencegah kasus melompat ke tahap yang tidak sah.

Sesudah tiket ini, Penilai dapat **mendaftarkan klaim**, mengikatnya ke polis yang benar, dan melihat kasus berjalan melalui **tahapan yang ditetapkan satu berkas alur** — ⛔ tahap di luar itu tidak dibuat.

## Area codebase

- Lapisan layanan klaim: pembuatan kasus dan pengikatan polis
- Mesin tahapan kasus — daftar tahap tertutup

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Tahapan kasus | satu berkas alur; tahap di luarnya tidak dibuat |
| Pengikatan polis | pembacaan data polis dan kutipan dari modul penawaran |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0014** — belum menyentuh tiket ini; wewenang di tiket 08

## Acceptance criteria

- [ ] **AC 4** — perilaku yang ditiru adalah perilaku **salinan PRODUKSI**; membandingkan terhadap salinan pengembangan **gagal**
- [ ] **AC 5 · 6 · 7** — tahapan kasus mengikuti berkas alur; tahap tambahan **ditolak**
- [ ] **AC 8–12** — registrasi dan pengikatan polis

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | ⚠️ menahan pengikatan polis yang lengkap |

## Perintah verifikasi

1. Buat satu klaim, periksa ia terikat ke polis yang benar.
2. Coba pindahkan kasus ke tahap yang **tidak ada** di berkas alur — ⛔ **ditolak**.
3. Tutup lalu buka kembali kasus — ⭐ tahapnya **tidak mundur**.
