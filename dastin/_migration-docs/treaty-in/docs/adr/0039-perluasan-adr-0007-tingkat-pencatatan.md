# ADR-0039 — Perluasan ADR-0007: tingkat pencatatan ikut tersimpan

**Status:** diterima, 23 September 2026
**Sifat:** ini **perluasan ADR-0007**, bukan keputusan yang berdiri sendiri. Ia menambah satu
dimensi pada paket yang sama.

## Konteks

ADR-0007 menetapkan nilai uang disimpan sebagai paket: nilai, mata uang, kurs, tanggal atau sumber
kurs, dan nilai IDR.

Rumus pencapaian di sistem lama membuka dimensi kedua yang belum pernah dinyatakan:

```
AchievementPct = ( TotalAchPremium / (RNMShareP/100) ) / EPI x 100
```

Premi dibagi bagian NuRe untuk **dinaikkan ke tingkat 100%**, lalu dibandingkan dengan EPI. Itu
bukan kebetulan aritmetika — itu pernyataan model: EPI disimpan pada tingkat 100% treaty,
sementara premi disimpan pada bagian NuRe. Dua besaran yang dibandingkan satu sama lain, disimpan
pada tingkat berbeda, didamaikan oleh pembagian yang tidak tertulis di mana pun.

## Keputusan

> Setiap nilai uang membawa **tingkatnya** — 100% treaty atau bagian NuRe.
> Bila ia besaran bagian, ia juga membawa **bagian yang dipakai** untuk menghasilkannya.

**Tingkat adalah sifat besarannya, bukan konvensi yang kita pilih.** Limit treaty *adalah* besaran
100%: ia disepakati dengan cedant tanpa peduli berapa bagian yang diambil NuRe, dan tetap angka itu
meski bagian NuRe berubah. Premi yang diterima *adalah* besaran bagian NuRe: itu yang masuk
rekening, dan tidak ada tingkat lain yang pernah terjadi.

Menyeragamkan keduanya bukan penyederhanaan — ia menghilangkan informasi pada satu arah dan
mengarangnya pada arah lain.

## Konsekuensi

- Penaikan dan penurunan tingkat berhenti menjadi pembagian tak tertulis yang tersebar di rumus,
  dan menjadi **konversi yang tercatat**, sejajar dengan konversi mata uang.
- Konversi tingkat memakai bagian **yang berlaku pada saat jumlah itu terjadi**, bukan bagian hari
  ini. Karena setiap jumlah membawa bagiannya sendiri, aturan ini terpenuhi dengan sendirinya.
- Besaran tingkat 100% dan besaran tingkat bagian **tidak berbagi kolom**.
- Rumus pencapaian tidak lagi memerlukan pembagian yang tidak dijelaskan, karena kedua pembandingnya
  sudah menyatakan tingkatnya sendiri.

## Tingkat EPI — diputuskan tanpa verifikasi

**EPI dicatat pada tingkat 100% treaty.**

**Dasar.** Rumus pencapaian membagi premi dengan bagian NuRe *sebelum* membandingkannya dengan EPI;
pembagian itu hanya masuk akal bila EPI berada di tingkat 100%. Dan bila EPI ternyata di tingkat
bagian, seluruh angka pencapaian selama ini berlebih sebesar satu per bagian — untuk bagian 20%,
lima kali lipat. Selisih sebesar itu terhadap target tahunan tidak mungkin luput bertahun-tahun.

**Syarat pembalikan.** Underwriting menyatakan sebaliknya. Bila itu terjadi, ini **bukan perubahan
model** melainkan **penyajian ulang angka historis**, dan langsung naik ke daftar eskalasi.
