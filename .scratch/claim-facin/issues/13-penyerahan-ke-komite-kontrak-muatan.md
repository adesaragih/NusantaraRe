# 13: Penyerahan ke komite — kontrak muatan

**Status:** ready-for-agent
**Blocked by:** 05 (penyesuaian) · 06 (uang) · 08 (wewenang)
**Menutup:** AC 43 · 44 · 45 · 46 · 47 · 48 · 49 · 50 *(8 AC)* — US 13 · 14

## Hasil & nilai pengguna

Hari ini Penyesuaian di atas kewenangan penilai **belum dapat diserahkan ke komite**, dan tidak ada kesepakatan tentang **apa yang diserahkan**.

Sesudah tiket ini, Penilai dapat **menyerahkan penyesuaian ke komite**, dan ⭐ **muatan yang diserahkan lengkap** — termasuk data kutipan **utuh**, sehingga penggolongan lini usaha di sisi komite tidak pincang.

## Area codebase

- Lapisan layanan klaim: penyerahan ke komite
- ⭐ Kontrak muatan — apa yang disalin dan seberapa lengkap

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pembentukan kasus komite | rule pembuat nomor komite; menyemai jumlah jenjang dan giliran mulai |
| Jalur satu jenjang | rule kirim tutup klaim dan kirim tolak klaim — ⭐ keduanya **berkomite satu** |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0014** — wewenang ditegakkan di lapisan layanan

## Acceptance criteria

- [ ] **AC 43–50** — penyerahan ke komite
- [ ] ⭐ **Data kutipan disalin UTUH**, ⛔ bukan daftar medan bernama — ⚠️ daftar bernama itulah yang melahirkan cacat MBU dan Travel
- [ ] ⭐ Jalur **tutup klaim** dan **tolak klaim** membentuk komite **satu jenjang**
- [ ] ⛔ Tiket ini **berhenti di kontrak muatan** — ⭐ perilaku komite ada di putaran sisi komite

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi daftar jabatan dan susunan jenjang | ⛔⛔ **MENAHAN** penentuan jumlah jenjang |

## Perintah verifikasi

1. Serahkan satu penyesuaian ke komite — ⭐ kasus komite lahir dengan muatan lengkap.
2. Periksa data kutipan di sisi komite — ⭐ **utuh**, bukan dua medan.
3. Serahkan lewat jalur tutup klaim — ⭐ komitenya **satu jenjang**.
