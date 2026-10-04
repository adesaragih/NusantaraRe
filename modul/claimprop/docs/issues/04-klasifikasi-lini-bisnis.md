# 04: Klasifikasi lini bisnis — satu kunci, satu mekanisme

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi klaim dan nomor polis)
**Menutup:** AC 71 · 72 · 73 · 74 *(4 AC)*

## Hasil & nilai pengguna

Klaim terklasifikasi ke lini bisnis yang benar dengan **satu** mekanisme, bukan dua yang saling
bertentangan. Klasifikasi itu menentukan template teks objek pertanggungan yang dilihat Claim Admin,
sehingga tiket **05 (Insured Interest)** bergantung padanya.

## Area codebase

Pemetaan kelompok treaty ke kelas lini bisnis · template teks objek pertanggungan.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/SetValueToClaim_Act.xml` | 8 | penautan klasifikasi ke kelompok treaty |
| `Activity/SetValueToClaim_Act.xml` | 9 · 10 · 11 · 12 | empat template teks objek pertanggungan menurut kelas |
| `RDBList/GetTreatyGroupID.xml` | — | sumber nilai kelompok treaty |
| `When/` (61 berkas) | — | ⭐ **52 dari 61** menguji properti pada halaman yang **tidak ada** pada objek kerja Claim Prop — **mati di sini, hidup di Claim Fac In**. ⛔ Teks lama: ~~sisa impor model Fac In~~ *(diralat 2026-09-19)* |
| `When/IsCustomBonds.xml` | — | berkas terbaru; **menambah** barisan sisa itu, bukan mengurangi |

## ADR terkait

Tidak ada ADR yang mengikat langsung.

## Acceptance criteria

- [ ] `[terverifikasi]` + `[data DBA]` Klasifikasi lini bisnis memakai **satu kunci: kelompok treaty**; tidak ada mekanisme kedua *(AC 71 spec)*
- [ ] ⭐ `[keputusan work owner]` **2026-09-19** Ke-61 rule klasifikasi lini **dipindahkan sekali sebagai satu himpunan**; **hidup-matinya tidak ikut dipindahkan**. Di Claim Prop **52 dari 61** menguji halaman yang tidak ada di sini sehingga **tidak pernah bernilai benar**; di Claim Fac In halaman itu ada dan rule yang sama **hidup**. **Satu salinan, dua nasib.** `[data DBA]` data klaim nyata memuat nama dan kode bisnis, bukan properti yang diuji rule-rule itu. ✅ **Penggolongannya SUDAH diputuskan: penerapan aturan berdiri**, bukan penyimpangan — label `**Alasan menyimpang:**` **dicabut** *(AC 72 spec)*
      > ⛔ **RALAT 2026-09-19** — teks lamanya **dikutip utuh, tidak dihapus**: *"⚠️ Rule klasifikasi warisan **tidak dimigrasikan**. **Alasan menyimpang:** 52 dari 61 rule `When` menguji properti pada halaman yang tidak ada pada objek kerja Claim Prop; dihidupkan pun semuanya bernilai salah. … ⚠️ `[terbuka]` **penggolongannya belum diputuskan** — lihat catatan di bawah; label `**Alasan menyimpang:**` dibiarkan apa adanya sampai dijawab."*
- [ ] ⚠️ **Cacat template dibawa apa adanya, dan itu keputusan sadar.** Marine Cargo memakai template proyek konstruksi yang sama dengan keranjang Aneka, **termasuk typo pada label**. **Catatan paritas:** ini **paritas, bukan penyimpangan** — dicatat di sini agar bila kelak diperbaiki, perbaikannya dicatat sebagai penyimpangan sadar tersendiri *(AC 73 spec)*
- [ ] ⚠️ `[terverifikasi]` Kode pada tabel **kelompok treaty** dan kode bernilai sama pada tabel **jenis treaty** adalah **dua enum berbeda**; test yang menyatukannya **gagal**. **Jebakan:** kode yang sama pernah dibaca sebagai jenis treaty dan ternyata salah — dua tabel kode, dua arti *(AC 74 spec)*

## ✅ Catatan `[terbuka]` ringan — **SUDAH DIJAWAB 2026-09-19**

> ✅ **`[keputusan work owner]` 2026-09-19 — butir ini DITUTUP.** Penggolongannya **penerapan
> aturan berdiri**. ⭐ **Sebabnya: pertanyaan di bawah terbukti salah pertanyaan.** Selama
> pilihannya dibingkai *"dimigrasikan atau tidak"*, ketiga aturan induk memang tidak menjawabnya.
> Keputusan **A5-6 Claim Fac In** — *rule dipisahkan dari kehidupannya* — membingkainya ulang:
> rule-nya **dipindahkan**; yang tidak dipindahkan adalah **anggapan bahwa ia hidup di sini**.
> ⛔ **Seluruh teks di bawah dibiarkan apa adanya sebagai jejak.**

⚠️ `[terbuka]` **Penggolongan AC 72 belum dapat diputuskan dari aturan berdiri.** Ketiga aturan induk
menjawab tiga hal lain: gerbang yang flag-nya mati (WHEN tidak ditulis, STEP tetap ditulis) · langkah
ber-remark (STEP tidak ditulis) · **elemen UI** selalu-salah (elemen tidak dibuat). **Tidak satu pun
menyebut sebuah rule utuh** yang selalu bernilai salah karena menguji halaman yang tidak ada.

Pertanyaannya: tidak memigrasikan 52 dari 61 rule `When` itu **penerapan aturan berdiri** — sejenis
dengan elemen selalu-salah yang tidak dibuat — atau **penyimpangan tersendiri** yang wajib dihitung?
Pemilik: **work owner**. **Tidak memblokir** — perilaku yang ditiru sama saja dalam kedua bacaan,
karena rule-rule itu bernilai salah di kedua dunia. Yang bergantung padanya hanya **penggolongan**,
bukan hasil kerjanya.

### ✅ Jalur catatan pengembang — **sudah diperiksa, tertutup** *(ronde 4, 2026-09-19)*

`[terverifikasi]` Ronde 3 menduga ada *"tiga rule `When` yang catatannya menyebut perubahan
penggolongan"* yang perlu diperiksa. Sensus 100% atas `pyMemo` menemukan **empat**, dan
**tidak satu pun mengubah pembacaan tiket ini**:

| Rule (`pxInsName`) | Catatan pengembang | Hasil |
| --- | --- | --- |
| `ASM-FW-GISFW-DATA!ISANEKA` | *"ganti IsBonding jadi IsBondingAndCustomBonds"* | ✅ **sudah dikerjakan** — rujukan yang berlaku `IsBondingAndCustomBonds`; `IsBonding` polos hanya tersisa di dalam teks catatan itu sendiri |
| `ASM-FW-GISFW-WORK!ISEDMADJSHARECEDANT` | *"save as ganti value"* | ⛔ bukan klasifikasi lini bisnis |
| `@BASECLASS!ISPEGAPROD` | *"ganti"* | ⛔ penanda lingkungan, bukan klasifikasi |
| `ASM-FW-GCNMFW-WORK-PNC!ISPA_PNC` | *"tambah pyWorkPage dan hapus policy"* | ⛔ kelas `WORK-PNC` — termasuk ke-52 rule yang **mati di modul ini** *(diralat 2026-09-19; teks lama: ~~sisa impor~~)* |

⛔ **Nol AC berubah.** Dicatat supaya jalur ini tidak diperiksa ulang.

### ✅ Lima rule `When` terakhir — **sapu bersih selesai**, nol yang hidup *(ronde 6, 2026-09-19)*

`[terverifikasi]` Kelima rule `When` yang belum pernah dibuka telah dibaca satu per satu:

| Rule (`pxInsName`) | Halaman yang diuji | Hidup? |
| --- | --- | --- |
| `ASM-FW-GISFW-DATA!ISMBD` | `pyWorkPage.Quotation.BusinessType` | ⛔ tidak |
| `ASM-FW-GISFW-DATA!ISMBU` | `pyWorkPage.OfferFacIn.QuotationData.BusinessType` | ⛔ tidak |
| `ASM-FW-GISFW-DATA!ISMARINEHULL` | `pyWorkPage.Quotation.BusinessType` | ⛔ tidak |
| `@BASECLASS!ISBILLBOARDNEONSYARIAH` | `pyWorkPage.Quotation.BusinessType` + `BusinessCode` | ⛔ tidak |
| `@BASECLASS!ISMAINTENANCE` | `pyWorkPage.Quotation.BusinessType` | ⛔ tidak |

Kelimanya menguji **halaman yang tidak ada** pada objek kerja Claim Prop — termasuk kelompok 52
rule sisa impor. ⛔ **AC 72 tidak berubah**, dan penggolongannya tetap `[terbuka]` di tangan work
owner.

⚠️ **Jebakan parser yang hampir mengenai:** `isMaintenance` tampak **tanpa kondisi** bila dibaca
dari medan *viewer*, yang berisi pola kosong. Kondisi sebenarnya ada di `pyConditionString`.
Aturan *"isi viewer adalah sisa basi"* berlaku dan terbukti lagi.

## Perintah verifikasi

```
jalankan test "kelompok treaty properti -> template Property"
jalankan test "kelompok treaty motor -> template Motor Vehicle"
jalankan test "kelompok treaty marine cargo -> template proyek konstruksi (paritas)"
jalankan test "kelompok treaty lain -> keranjang Aneka"
jalankan test "kode kelompok treaty dan kode jenis treaty tidak saling tertukar"
cari rule klasifikasi warisan yang dimigrasikan           -> nihil
```
