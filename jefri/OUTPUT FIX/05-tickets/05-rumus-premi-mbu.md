# 05: Rumus premi MBU — bentuk bersarang

**What to build:** Premi MBU dihitung benar, dan bersamanya **aturan pembagi komposit terbukti dari
struktur ekspresinya sendiri** — bukan dari tafsir kita.

MBU memakai bentuk **bersarang**, dan di situlah nilainya. Pembagi **dalam** menempel pada rate;
pembagi **luar** menempel pada pro-rata. Keterikatan faktor→pembagi karena itu **tertulis eksplisit**,
tidak perlu disimpulkan:

```
@Math.divide( @Math.divide((.TSI*(.Rate+.Loading)),100,4) * @if(.ProRatePercent=="",100,.ProRatePercent), 100, 4)
                                          ^^^                                                            ^^^
                              pembagi DALAM → .Rate (%)                        pembagi LUAR → ProRatePercent
```

Asal: `FillPremiMBU_FacIn` **L1144** dan **L1626**.

Dua hal yang mudah terlewat: `.Loading` **dijumlahkan ke rate sebelum dikalikan** (bukan sesudah),
dan pro-rata kosong diperlakukan sebagai **100**, bukan nol — membalik keduanya mengubah angka tanpa
gejala.

**Blocked by:** 03

**Status:** ready-for-agent

- [ ] Kedua bentuk (`L1144`, `L1626`) diimplementasikan sebagai cabang internal di balik pintu masuk yang sama
- [ ] Struktur **bersarang** dipertahankan — bukan diratakan menjadi satu pembagian dengan pembagi 10.000
- [ ] `.Loading` dijumlahkan ke rate **sebelum** perkalian, sesuai sumbernya
- [ ] Pro-rata kosong → **100**, direproduksi apa adanya
- [ ] Skala MBU **persen** diambil dari resolver (tiket 03), bukan ditulis ulang di sini
- [ ] **Rekonsiliasi eksak** per bentuk, tanpa toleransi
- [ ] Kasus uji yang membuktikan pembagi dalam terikat ke rate dan pembagi luar ke pro-rata — bila keduanya tertukar, uji harus **gagal**
