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

**Status:** ready-for-agent

- [ ] Keempat bentuk diimplementasikan sebagai cabang internal di balik satu pintu masuk perhitungan — **bukan** empat seam
- [ ] Presisi tiap pembulatan ditulis **literal di tempatnya**, disertai komentar yang menyebut rule asal dan nomor langkahnya
- [ ] Presisi **20 desimal** pada bentuk `L1146` direproduksi apa adanya, tidak diseragamkan ke 4
- [ ] Urutan operasi tiap bentuk identik dengan sumbernya
- [ ] **Rekonsiliasi eksak** per bentuk: keempat fixture cocok sampai digit terakhir, tanpa toleransi
- [ ] Pengurangan diskon terjadi di posisi yang sama seperti sumbernya
