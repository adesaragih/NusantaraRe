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

**Status:** wontfix — Layering tidak dipakai di sistem baru (butir 30)

- [ ] Rumus diimplementasikan sebagai cabang internal di balik pintu masuk perhitungan yang sama
- [ ] Skala **per mille** diambil dari resolver (tiket 03)
- [ ] Pembagi `100000` diuraikan sebagai `1.000 × 100`, bukan diperlakukan sebagai satu satuan
- [ ] **Rekonsiliasi eksak**, tanpa toleransi
- [ ] Label tampilan daftar Layer diperbaiki dari `(%)` menjadi `(‰)` — **perbaikan label saja**, dicatat bahwa nol angka berubah sehingga rekonsiliasi tidak terganggu

## Comments

### 2026-10-01 — implementasi (agent)

Kode: `premium/layering.go`. `[terverifikasi]` `GenerateLayerList_ACT` langkah 2 `.Premium = 0`, langkah 3
diulang 0…`param.Layer − 1`, langkah 3.2 menambah satu lapisan berpremi
`round4(.TSI × .Rate × ProRatePercent ÷ 100000)` (coverage yang sama di setiap putaran) lalu
`.Premium += premi lapisan`. Masukan baru: `JumlahLapisan` (`param.Layer`).

| Kriteria | Keadaan |
| --- | --- |
| Cabang di pintu yang sama | ✅ `Calculate` lini `Layering` |
| Per mille dari resolver; 100000 = 1.000 × 100 | ✅ |
| Rekonsiliasi eksak | ⛔ **terbuka** — tidak ada kasus nyata ber-`LayerList` terisi di `DDL\CONTOH` (semua tag kosong) (`../PERTANYAAN-LANJUTAN.md` butir 1) |
| Label `(%)` → `(‰)` | ⏸ label layar; belum ada layar |

Uji mutasi: 3/3.

### 2026-10-01 — DITUTUP (keputusan work owner, butir 30)

"Layering tidak dipakai lagi di sistem baru." `layering.go`, `layering_test.go`, `Input.JumlahLapisan`,
`ErrJumlahLapisan`, dan `LiniLayering` dihapus; keputusan agent A6 gugur. Rumus L909 tetap tercatat di
K-018 sebagai fakta korpus. Permintaan kasus Layering nyata (`../PERTANYAAN-LANJUTAN.md` butir 1) dicabut.
