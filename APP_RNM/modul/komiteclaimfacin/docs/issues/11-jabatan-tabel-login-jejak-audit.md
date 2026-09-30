# 11: Jabatan, tabel login, dan jejak audit

**Status:** ready-for-agent
**Blocked by:** 05 · `claim-facin\issues\12` *(jejak audit klaim)*
**Menutup:** AC 77 · 78 · 79 · 80 · 81 · 82 · 83 · 84 · 85 *(9 AC)* — US 21 · 22 · 23 · 26

## Hasil & nilai pengguna

Hari ini Jabatan **ditanam di dalam kode** sebagai rantai pencocokan ID pengguna perorangan, dan ⚠️ **orang yang sama tertulis dengan jabatan berbahasa berbeda antar modul**. ⛔ Mengubah susunan jenjang menuntut **rilis**.

Sesudah tiket ini, ⭐ Jabatan menjadi **kolom berisi KODE di tabel login**, susunan jenjang komite ada di **daftar terpisah**, dan ⭐ **perubahan jabatan masuk jejak audit** — sebab jabatan **menentukan siapa boleh menyetujui**.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| ⚠️ Tabel jabatan lama | ⛔ rantai pencocokan **ID pengguna perorangan** dengan jabatan |
| ⚠️ Dua bahasa | ⛔ orang yang sama tertulis berbeda antar modul, dan **kuncinya beda medan** |
| Roster lama | ⛔ menyimpan **nama orang**, bukan jabatan — ⭐ itulah akar penukaran akun |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K13** — ⭐⭐ **Empat syarat:** *(1)* isinya **KODE** menunjuk daftar induk jabatan · *(2)* **urutan jenjang di daftar terpisah** · *(3)* catatan keputusan **menyimpan SALINAN kode jabatan** · *(4)* perubahan kolom jabatan **masuk jejak audit**
- **K12** — ⭐ **Roster berbasis JABATAN**, ⛔ bukan nama orang
- **K11** — ⭐ Jabatan dibaca dari **data pengguna**; ⛔ **nilai bawaan tidak dibawa**

## Yang harus diuji

- [ ] ⭐ Kolom jabatan berisi **kode**, menunjuk **daftar induk jabatan** — ⛔ bukan teks bebas
- [ ] ⭐ **Susunan jenjang komite** disimpan **terpisah** dari tabel login
- [ ] ⭐ Mengubah susunan jenjang adalah **mengubah baris data** — ⛔ bukan rilis
- [ ] ⭐ **Perubahan kolom jabatan masuk jejak audit**: siapa mengubah, kapan, dari kode apa ke kode apa
- [ ] ⛔ Nama orang **tidak menjadi kunci** di mana pun
- [ ] ⭐ Tulisan jabatan **seragam** — ⚠️ di sistem lama berbeda bahasa antar modul

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | ⛔⛔ **Isi daftar induk jabatan dan susunan jenjang BELUM ADA** | ⛔⛔ **MENAHAN PEMBANGUNAN** |
| **33-lama** | Alasan historis pemetaan akun menjadi nama orang lain | tidak menahan |

## Seam & verifikasi

**Seam:** lapisan layanan — pembacaan jabatan dan penulisan jejak audit.
⭐⭐ **UJI SALINAN JABATAN — inti tiket ini:**
   *(a)* catat satu keputusan komite;
   *(b)* **naikkan jabatan** orang yang memutuskannya;
   *(c)* ⭐ **catatan lama TIDAK berubah**.
⚠️ **Ini bukan uji fungsi** — ia yang **membedakan SALINAN dari RUJUKAN**. ⛔ Tanpa uji ini, pengembang berikutnya akan menggantinya dengan rujukan karena terlihat lebih rapi.
1. Ubah nama tampil orangnya ⇒ ⭐ catatan lama **ikut nama baru** — ⭐ itu memang dikehendaki.
2. Ubah kolom jabatan seseorang ⇒ ⭐ perubahannya **tercatat di jejak audit**.
3. Ubah susunan jenjang ⇒ ⭐ cukup **mengubah baris**, tanpa rilis.
