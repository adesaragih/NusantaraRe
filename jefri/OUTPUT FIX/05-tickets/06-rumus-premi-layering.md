# 06: Rumus premi Layering

**What to build:** Premi untuk struktur berlapis dihitung benar, dengan rate berskala **per mille**.

Asal: `GenerateLayerList_ACT` **L909**:

```
@Math.divide((.TSI * .Rate * pyWorkPage.OfferFacIn.ProRatePercent), 100000, 4)
```

Bentuknya identik dengan salah satu varian PA, dan terurai sama: `100000 = 1.000 × 100`.

⚠️ **Ini satu-satunya ekspresi ber-rate di seluruh berkas bernuansa Layer di korpus NB** — tidak ada
varian yang bertentangan. Kepastian itu penting justru karena label layar daftar Layer menulis
`Rate (%)`, yang **bertentangan dengan rumusnya**. Rumus yang menang; labelnya yang diperbaiki, dan
perbaikan label **tidak mengubah satu angka pun**.

**Blocked by:** 03

**Status:** ready-for-agent

- [ ] Rumus diimplementasikan sebagai cabang internal di balik pintu masuk perhitungan yang sama
- [ ] Skala **per mille** diambil dari resolver (tiket 03)
- [ ] Pembagi `100000` diuraikan sebagai `1.000 × 100`, bukan diperlakukan sebagai satu satuan
- [ ] **Rekonsiliasi eksak**, tanpa toleransi
- [ ] Label tampilan daftar Layer diperbaiki dari `(%)` menjadi `(‰)` — **perbaikan label saja**, dicatat bahwa nol angka berubah sehingga rekonsiliasi tidak terganggu
