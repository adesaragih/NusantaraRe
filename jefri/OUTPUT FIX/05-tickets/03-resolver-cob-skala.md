# 03: Resolver lini bisnis → skala rasio, dan aturan pembagi komposit

**What to build:** Setiap rasio yang masuk perhitungan mendapat skalanya dari **satu resolver
terpusat** yang memetakan lini bisnis (COB) → per mille atau persen. Tidak ada tempat kedua yang bisa
menyimpang.

Properti rate yang sama memakai **dua satuan berbeda menurut lini bisnis**. Menyeragamkannya
menggeser premi **satu ordo besaran 10** pada separuh portofolio — jadi bahayanya bukan salah memilih
satuan, melainkan menyamakannya.

**Aturan penguraian pembagi komposit — inti tiket ini.** Pembagi gabungan pada satu operasi pembagian
adalah **hasil kali** sumbangan tiap faktor bersatuan, bukan satu satuan tunggal. Satuan rate dibaca
dari **faktor yang menempel padanya saja**, tidak pernah dari pembagi total:

```
@Math.divide((.TSI * .Rate * ProRatePercent), 100000, 4)
        100000  =  1000 (rate ber-‰)  ×  100 (ProRatePercent ber-%)
```

Pembagi `100000` karena itu **menegaskan** rate = per mille — bukan membantahnya. Siapa pun yang
membacanya sebagai "satuan lebih kecil dari per mille" akan salah satu ordo besaran.

Resolver mengikuti **rumus**, bukan label layar: tiga label layar terbukti bertentangan dengan
rumusnya, dan rumus yang menang di ketiganya.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] Peta skala lengkap: **PA · Layering · FIRE = per mille**; **MBU · ANEKA · BONDING · GOLF · MARINE CARGO = persen**
- [ ] Aturan pembagi komposit diterapkan: sumbangan per faktor, bukan pembagi total
- [ ] Kasus uji membuktikan pembagi `100000` menghasilkan **per mille**, bukan satuan lain
- [ ] **Lini bisnis yang belum ada di peta → gagal keras**, bukan skala default
- [ ] Skala **tidak** dibaca dari tabel konfigurasi yang dapat diubah tanpa deployment — ia fakta struktural korpus, bukan parameter bisnis
- [ ] Resolver adalah **satu-satunya** pengisi skala; tidak ada pemanggil yang menyetel skala sendiri
