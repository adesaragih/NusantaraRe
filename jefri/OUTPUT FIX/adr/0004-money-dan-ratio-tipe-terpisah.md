---
status: accepted
---

# `Money` dan `Ratio` adalah dua tipe yang tidak dapat dijumlahkan; skala melekat pada nilai

Nilai uang memakai `Money{Amount decimal, Currency}`; rate dan persen memakai
`Ratio{Value decimal, Scale}` yang **membawa satuannya sendiri**. Keduanya tipe berbeda: `Money +
Ratio` tidak dapat dikompilasi, dan satu-satunya jembatan adalah operasi eksplisit `Money × Ratio →
Money`. Skala sebuah `Ratio` diisi oleh **satu resolver terpusat** yang memetakan lini bisnis (COB) →
skala.

## Mengapa, dengan bukti

`[terverifikasi]` Properti `.Rate` yang sama memakai **dua satuan berbeda** menurut COB: FIRE dan PA
per mille (‰); ANEKA, BONDING, GOLF, MARINE CARGO, dan MBU persen (%). Buktinya menutup kemungkinan
"salah satu rule bug" — satu berkas memuat tiga COB berdampingan dengan dua satuan, dan label layar
menyebut satuannya harfiah (`‰ Standard Rate` dengan karakter U+2030 pada layar FIRE).

Menyeragamkan satuan menggeser premi **satu ordo besaran 10** pada separuh portofolio. Tipe yang
membawa skalanya sendiri mengubah cacat itu dari kesalahan runtime yang senyap menjadi **kegagalan
kompilasi**.

## Considered Options

- **Satu tipe `Money` untuk semuanya, persen sebagai `decimal` telanjang** — ditolak. Justru
  menghapus proteksi yang paling dibutuhkan.
- **Empat tipe (`Money`, `Rate`, `Percent`, `Limit`)** — ditolak. `Limit` berperilaku persis seperti
  `Money`; tipe tambahan tanpa perilaku tambahan hanya menambah friksi.
- **Skala disimpan di tabel konfigurasi Oracle** — ditolak, dan ini penolakan yang paling penting
  untuk diingat. Skala satuan adalah **fakta struktural yang terverifikasi dari korpus**, bukan
  parameter bisnis. Menaruhnya di tabel yang dapat diubah tanpa deployment mengulang persis
  kerentanan yang sudah kita temukan pada `M_PROMPT_AI`: perilaku sistem berubah tanpa jejak
  version control.

## Consequences

### Aturan penguraian pembagi komposit — wajib dipakai saat membaca rumus premi

Ini satu-satunya cara membaca satuan dengan benar, dan mengabaikannya menghasilkan kesimpulan yang
terbalik.

**Pembagi gabungan pada satu `@Math.divide` adalah hasil kali faktor-faktor satuan, bukan satu
satuan tunggal.** Satuan `.Rate` disimpulkan dari **faktor yang menempel padanya saja** — tidak
pernah dari pembagi total.

| Faktor dalam ekspresi | Sumbangan ke pembagi |
| --- | ---: |
| `.Rate` ber-satuan ‰ | 1.000 |
| `.Rate` ber-satuan % | 100 |
| `ProRatePercent` (persen) | 100 |

**Contoh terverifikasi** — `GenerateLayerList_ACT.xml` L909:

```
@Math.divide((.TSI * .Rate * pyWorkPage.OfferFacIn.ProRatePercent),100000,4)
        100000  =  1000 (.Rate ber-‰)  ×  100 (ProRatePercent ber-%)
```

Pembagi `100000` karena itu **menegaskan** `.Rate` = ‰ — bukan membantahnya. Siapa pun yang membaca
`100000` sebagai "satuan yang lebih kecil dari ‰" akan salah satu ordo besaran.

⚠️ **Struktur ekspresi MBU membuktikan aturan ini secara langsung.**
`FillPremiMBU_FacIn.xml` L1144 memakai bentuk **bersarang**, sehingga keterikatan faktor→pembagi
tertulis eksplisit dan bukan tafsir:

```
@Math.divide( @Math.divide((.TSI*(.Rate+.Loading)),100,4) * @if(.ProRatePercent=="",100,.ProRatePercent), 100, 4)
                                                  ^^^                                                    ^^^
                                      pembagi DALAM menempel .Rate (%)              pembagi LUAR menempel ProRatePercent
```

### Peta skala per COB — terkunci

Nilainya ditetapkan **K-018** (`00-KEPUTUSAN-WORK-OWNER.md`), berikut kutipan baris sumbernya:

| COB | Skala `.Rate` | Pembagi yang menempel |
| --- | :-: | ---: |
| FIRE · PA · **Layering** | **‰** | 1.000 |
| MBU · ANEKA · BONDING · GOLF · MARINE CARGO | **%** | 100 |

### Lain-lain

- Resolver COB → skala mengikuti **rumus**, bukan label layar. `[terverifikasi]` tiga label UI
  bertentangan dengan rumusnya — dan K-018 menetapkan **rumus yang menang** di ketiganya. Perbaikan
  labelnya tidak mengubah satu angka pun, sehingga tidak mengganggu rekonsiliasi.
- COB yang belum ada di peta resolver → **`panic`**, bukan skala default. Menebak skala adalah persis
  cacat 10× yang tipe ini dirancang untuk mencegah.
