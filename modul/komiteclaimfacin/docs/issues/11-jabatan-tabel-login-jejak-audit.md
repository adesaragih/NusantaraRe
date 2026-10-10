# 11: Jabatan, tabel login, dan jejak audit

**Status:** ready-for-agent
**Blocked by:** 05 · `claim-facin\issues\12` *(jejak audit klaim)*
**Menutup:** AC 77 · 78 · 79 · 80 · 81 · 82 · 83 · 84 · 85 *(9 AC)* — US 21 · 22 · 23 · 26

> ⛔ **RALAT 10-10-2026 — TIKET GUGUR, jangan dibangun.** Kepala lamanya dikutip utuh, tidak dihapus: *"**Status:**
> ready-for-agent"* dan *"**Menutup:** AC 77 · 78 · 79 · 80 · 81 · 82 · 83 · 84 · 85 *(9 AC)* — US 21 · 22 · 23 · 26"*
> → **K13 "kolom jabatan di tabel login" digantikan workbasket** (prompt tahap 2 §3 "Gugur", KCF-01). Tidak ada kolom
> jabatan di tabel login, tidak ada daftar induk jabatan, dan tidak ada jejak audit perubahan jabatan di modul ini.
> Susunan jenjang tetap roster `EMAILKOMITE` STS_KLAIM FACIN (DEGREE, JABATAN, LIMIT_BOTTOM / LIMIT_TOP;
> `ApprovalKomite_Act` L4), dan `OPERATOR_ID`-nya di-UPDATE di tempat menjadi workbasket lewat migrasi komiteclaimfacin
> 640–679 (pola claimprop 537 / claimnonprop 611). Pergantian orang = pergantian anggota workbasket (Kelola User), tanpa
> rilis. Nasib AC-nya ada di RALAT spec AC 77–83 dan AC 85: AC 77 / 78 / 82 / 83 gugur; AC 79 / 80 / 81 / 84 isinya
> tetap dengan sumber roster + KCF-01; AC 85 alasannya diganti workbasket. **Uji salinan jabatan pindah ke tiket 05.**

## Hasil & nilai pengguna

Hari ini Jabatan **ditanam di dalam kode** sebagai rantai pencocokan ID pengguna perorangan, dan ⚠️ **orang yang sama tertulis dengan jabatan berbahasa berbeda antar modul**. ⛔ Mengubah susunan jenjang menuntut **rilis**.

Sesudah tiket ini, ⭐ Jabatan menjadi **kolom berisi KODE di tabel login**, susunan jenjang komite ada di **daftar terpisah**, dan ⭐ **perubahan jabatan masuk jejak audit** — sebab jabatan **menentukan siapa boleh menyetujui**.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| ⚠️ Tabel jabatan lama | ⛔ rantai pencocokan **ID pengguna perorangan** dengan jabatan |
| ⚠️ Dua bahasa | ⛔ orang yang sama tertulis berbeda antar modul, dan **kuncinya beda medan** |
| Roster lama | ⛔ menyimpan **nama orang**, bukan jabatan — ⭐ itulah akar penukaran akun |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"⚠️ Tabel jabatan lama | ⛔ rantai pencocokan
> **ID pengguna perorangan** dengan jabatan"* dan *"Roster lama | ⛔ menyimpan **nama orang**, bukan jabatan — ⭐ itulah
> akar penukaran akun"* → rantai pencocokan jabatan di modul ini adalah **sisa editor** (spec RALAT lintas-ronde #2);
> jalur hidupnya memakai `IDKomite` = `.JABATAN` roster. Roster `EMAILKOMITE` FACIN menyimpan **JABATAN** dan **akun
> orang** di `OPERATOR_ID` (DEV 10-10-2026), bukan nama orang. Akar penukaran akun (`SetProteksiSubmiteKomite` L1)
> dibuang (prompt §5 #1) dan `OPERATOR_ID` diganti workbasket (KCF-01).

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
