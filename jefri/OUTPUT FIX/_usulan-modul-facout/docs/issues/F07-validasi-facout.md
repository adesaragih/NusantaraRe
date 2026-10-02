# F07: Validasi Fac Out sebelum penawaran dikirim

**What to build:** Sebelum penawaran retrosesi berjalan, sistem memeriksa bahwa daftar objek Fac Out
konsisten dengan objek Fac In induknya dan bahwa reasuradur tujuannya sudah terisi — dan menolak
dengan pesan bila tidak.

`[terverifikasi]` Dua pemeriksa terpisah, keduanya `Rule-Obj-Activity`:

| Pemeriksa | Gerbang yang terbaca | Akibat |
| --- | --- | --- |
| `CheckDataFacOut_Act` | `Local.ListObjectFacout != Local.ListObjectFacIn` · `Local.FlagError==1` · `pyWorkPage.OfferFacIn.FlagSaveFO=="1"` | `Page-Set-Messages` — pesan galat |
| `CheckProtectFacout_Act` | `.ReinsurerName==""` · `Local.reas<1` · `ProtectReasFacOut.CARI2=="1"` · `…FacRetroDetails.BackUpStatus==3` | `Property-Set-Messages` + `Obj-Save`; menyetel `QuotationData.FacOutStatus` |

⚠️ **Perbandingan daftar objek dilakukan sebagai perbandingan dua nilai terangkai**, bukan sebagai
himpunan. `[terverifikasi]` Keduanya variabel lokal yang dibandingkan langsung dengan `!=`. Bentuk
perakitannya **diport apa adanya** — mengubahnya menjadi perbandingan himpunan akan mengubah kapan
galat muncul (urutan dan duplikat mulai berpengaruh).

`[terverifikasi]` **Ambang jumlah lokasi.** `Activity\OfferFacOut_PreAct` memuat gerbang
`@Utilities.SizeOfPropertyList(.OfferFacIn.LocationList)>50`. Ambang **50** diport apa adanya; artinya
**`[pertanyaan terbuka]`** — korpus tidak menjelaskan mengapa 50.

⚠️ `.FlagSaveFO=="1"` — angka dibandingkan **sebagai string**, sejalan dengan pola yang sudah dicatat
di F05. Kandidat K-046.

**Asal (Pega).** `Activity\CheckDataFacOut_Act` · `CheckProtectFacout_Act` · `SetErrorMessageFacOut` ·
gerbang ambang lokasi di `OfferFacOut_PreAct`

**Keputusan.** **K-046** · K-053 · `CLAUDE.md` §1, §4.6

**Blocked by:** F05

**Status:** blocked

- [ ] Ketidaksesuaian daftar objek Fac Out vs Fac In **menghasilkan pesan galat**, bukan penolakan diam
- [ ] Reasuradur kosong menghalangi lanjut, dengan pesan yang dapat dibaca pengguna
- [ ] Ambang **50 lokasi** diport apa adanya; komentar menandainya `[pertanyaan terbuka]`
- [ ] Perbandingan daftar objek **mempertahankan bentuk aslinya**, tidak diubah menjadi perbandingan himpunan
- [ ] **K-046** `K046_FlagSaveFO_AngkaSebagaiString`
- [ ] Pesan galat **tidak memuat** nama orang, nomor polis, atau data pelanggan (K-025)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
