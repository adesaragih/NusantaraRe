# 04: Rumus premi PA — keempat bentuk terverifikasi

**What to build:** Premi PA dihitung benar untuk **keempat bentuk rumus** yang benar-benar ada di
sistem lama, bukan satu bentuk yang mewakili. Keempatnya hidup berdampingan di satu rule dan sepakat
bahwa rate PA berskala **per mille** — dua memberi pembagi 1.000 langsung, dua memberi 100.000 yang
terurai menjadi 1.000 × 100.

Asal, `CalculatePremiPA_FacIn`:

```
L713   (@Math.divide((.TSI * .Rate),100000,4) * ProRatePercent) - .Discount
L858   (@Math.divide((.TSI*.Rate),1000,4) * @Math.divide((.PctShortPeriod),100,4)) - .Discount
L1003  (@Math.divide((.TSI*.Rate),1000,4)*1) - .Discount
L1146  @Math.divide((.TSI*ProRatePercent*.Rate),100000,20) - .Discount
```

⚠️ **Presisi berbeda di dalam rule yang sama** — `L1146` membulatkan **20 desimal**, tiga lainnya
**4 desimal**. Itu bukan kekeliruan yang perlu dirapikan; itu perilaku terekam. Presisi ditulis
**literal di tempatnya**, bukan disentralkan ke registry: satu rule tidak punya satu presisi.

**Urutan operasi dipertahankan persis** — `round(a,4) * b`, bukan `round(a*b,4)`. Ini yang menentukan
digit terakhir, dan digit terakhir adalah ukuran keberhasilan.

**Blocked by:** 03

**Status:** ready-for-human — sebagian: keempat rumus diport; L713/L1003 cocok kasus nyata, L858 (metode 2) dan diskon persen tanpa kasus nyata — menunggu DBA (paket D-1)

- [ ] Keempat bentuk diimplementasikan sebagai cabang internal di balik satu pintu masuk perhitungan — **bukan** empat seam
- [ ] Presisi tiap pembulatan ditulis **literal di tempatnya**, disertai komentar yang menyebut rule asal dan nomor langkahnya
- [ ] Presisi **20 desimal** pada bentuk `L1146` direproduksi apa adanya, tidak diseragamkan ke 4
- [ ] Urutan operasi tiap bentuk identik dengan sumbernya
- [ ] **Rekonsiliasi eksak** per bentuk: keempat fixture cocok sampai digit terakhir, tanpa toleransi
- [ ] Pengurangan diskon terjadi di posisi yang sama seperti sumbernya

## Comments

### 2026-10-01 — temuan dari tiket 01 (agent)

⚠️ **"Keempat bentuk hidup berdampingan" tidak berlaku untuk premi akhir.** `[terverifikasi]` Di
`CalculatePremiPA_FacIn` langkah 7 (L1146, `pyStepsPreCondition=false`) berjalan sesudah langkah 4/5/6
(L713/L858/L1003) dan menimpa `.Premium` yang sama. Menurut P-11 langkah itu tetap jalan tanpa syarat,
sehingga L713, L858, dan L1003 tidak pernah keluar sebagai premi akhir. Satuan ‰ keempatnya (K-018)
tetap benar. Kriteria "rekonsiliasi eksak per bentuk" perlu ditinjau ulang work owner setelah butir 3
`../PERTANYAAN-NB01.md` terjawab. Langkah 2–3 (penghitungan ulang `.Discount` dari `.Premium` lama)
juga belum diport.

### 2026-10-01 — RALAT entri di atas (agent)

⛔ **Entri sebelumnya salah.** Ia menyatakan langkah 7 (L1146) selalu menimpa `.Premium`, sehingga L713,
L858, dan L1003 "tidak pernah keluar sebagai premi akhir". `[terverifikasi]` Empat kasus PA nyata
membantahnya: metode 1 cocok eksak dengan **L713**, tiga kasus metode 3 cocok eksak dengan **L1003**, dan
L1146 cocok dengan tidak satu pun. Keputusan work owner 01-10-2026: langkah 7 tidak dijalankan
(`../KEPUTUSAN-30-09-2026.md` butir 13).

Jadi "keempat bentuk" tiket ini praktis **tiga bentuk hidup + satu langkah yang tidak berpengaruh**.
Dua sudah diport dan direkonsiliasi di tiket 01 (L713, L1003). Sisa tiket ini: **L858** (short period,
butuh `.PctShortPeriod` dan kasus metode 2 untuk rekonsiliasi), langkah 2–3 (diskon), dan keputusan
tentang L1146 setelah `../PERTANYAAN-NB01.md` butir 5 terjawab.

### 2026-10-01 — implementasi (agent)

Kode: `premium.go` (`premiPA`, `premiPAPeriodePendek`). Bentuk yang hidup kini lengkap:

| Bentuk | Metode | Keadaan |
| --- | :-: | --- |
| L713 prorata | `'1'` | ✅ rekonsiliasi 1 kasus nyata (tiket 01) |
| L858 short period | `'2'` | ✅ diport; **contoh hitung saja** — tidak ada kasus metode 2 di `DDL\CONTOH` (`../PERTANYAAN-LANJUTAN.md` butir 2) |
| L1003 fixrate | `'3'` | ✅ rekonsiliasi 3 kasus nyata (tiket 01) |
| L1146 | — | ⛔ tidak dijalankan (butir 13) |

Langkah 2 (`DiscountType == "Percent"`: `.Discount = round20(DiscountPercentage × .Premium LAMA ÷ 100)`)
**diport**, dengan masukan baru `PremiSebelumnya`, dan berjalan sebelum rumus premi. Langkah 3
(`"Amount"`) hanya menulis `.DiscountPercentage` — tidak diport, belum ada pemakainya. Uji mutasi: 3/3.
