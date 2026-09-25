# 07: Uang — presisi penuh, persentase bukan uang, pajak brokerage apa adanya

**Status:** ready-for-agent
**Blocked by:** 01
**Menutup:** AC 18 · 23 · 24 · 25 · 26 · 27 · 28 · 85 · 86 *(9 AC)* — US 21 · 22 · 36 · 38 · 39

## Hasil & nilai pengguna

Hari ini nilai uang ditampilkan di layar apa adanya seperti tersimpan, ⭐ **bukan dihitung ulang**.
⚠️ `[terverifikasi]` Empat medan yang tampak seperti uang sebenarnya **persentase** — `12.5` berarti
12,5 persen. Dan potongan brokerage dibagi **1,022** bila jenis pajaknya *inclusive*, ⛔ dengan
pembanding teks yang **persis**: huruf kecil atau spasi di belakang menggeser angkanya **2,2 %**
tanpa pesan galat.

Sesudah tiket ini, nilai uang tersimpan **berpresisi penuh**, ditampilkan beserta mata uangnya, dan
⭐ rumus pajak brokerage berjalan **persis** seperti sistem lama — termasuk ketidakseragamannya.

## Area codebase

- Lapisan repository: pembacaan nilai uang dan mata uang pasangannya
- Fungsi murni: rumus potongan pajak brokerage
- Lapisan handler: penyajian

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Rumus pajak brokerage | `@if(TypeTax == "Inclusive", Deduction / (102.2/100), Deduction)` — **12 tempat pada 7 Activity** |
| Nilai bawaan jenis pajak | `"Inclusive"`; ⛔ nilai lawannya **tidak pernah tertulis** di mana pun |
| Medan uang | **kolom tersimpan** pada sumber relasional; ⛔ nol `Activity` atau `DataTransform` menghasilkannya |

## ADR terkait

- ⭐ **ADR-0003** — uang **tidak** direpresentasikan sebagai `float`

## Acceptance criteria

- [ ] **AC 23** — nilai uang disimpan **berpresisi penuh**
- [ ] **AC 24** — pembulatan **hanya** di titik penyajian, ⛔ tidak pernah di repository
- [ ] **AC 25** — uang **tidak pernah** diwakili tipe pecahan biner
- [ ] **AC 26** — empat medan itu dibaca sebagai **persentase**, bukan jumlah uang
- [ ] **AC 27** — jenis pajak *inclusive* ⇒ potongan dibagi **1,022**
- [ ] **AC 28** — nilai lain — ⭐ **termasuk kosong, huruf kecil, atau berspasi** — ⇒ potongan
      dipakai apa adanya
- [ ] **AC 18** — kedelapan medan uang **ditampilkan apa adanya**, ⛔ tidak dihitung ulang
- [ ] **AC 85** — angka uang ditampilkan beserta **kode mata uang pasangannya**
- [ ] **AC 86** — ⛔ format penyajian **belum ditetapkan**; penyajian tidak dibangun dengan format
      yang ditebak

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | tipe penyimpanan kolom uang — sisi Pega menyimpannya sebagai **teks** | ⚠️ menahan **penguraian**, bukan rumusnya |
| **7** | format penyajian — desimal dan pemisah ribuan | ⚠️ menahan **penyajian** saja |
| **8** | satu medan persentase **dari apa** | tidak menahan |
| **9** | ketidakseragaman presisi 4 lawan 8 desimal | tidak menahan |

## Perintah verifikasi

1. Jenis pajak `"Inclusive"` — ⭐ potongan **dibagi 1,022**.
2. Jenis pajak `"inclusive"` huruf kecil — ⭐ potongan **tidak** dibagi; selisihnya **2,2 %**.
3. Jenis pajak **kosong** — ⭐ sama seperti butir 2.
4. Simpan nilai berdesimal panjang, baca kembali — ⭐ **tidak berubah**.

## Catatan

⭐ **Penyajian uang TIDAK bergantung pada P18.** `[keputusan work owner]` P43 menetapkan angka
**disalin apa adanya**, sehingga bagian terbesar uang dapat dibangun sekarang.

> ⚠️ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat penutup semula berbunyi:
> > *"⛔ Yang menunggu P18 adalah **perhitungannya**, ada di tiket **13**."*
>
> ⭐ **P18 ditarik** — perhitungannya pun tidak menunggu P18. Rumus pajak brokerage terbaca penuh
> di ekspor: `@if(TypeTax="Inclusive", @divide(Deduction,@divide(102.2,100,8),8), Deduction)`,
> **PPH 2 %**, **PPN 2,2 %**. Perhitungannya ada di tiket **13**, yang kini tertahan oleh **P30**
> dan **P8** saja.
