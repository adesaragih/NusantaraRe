# F13: Cetak R/I Slip dan kirim surat retrosesi ⛔ BLOCKED EKSTERNAL

**What to build:** Setelah tangga retro selesai, sistem mencetak R/I Slip retrosesi dan mengirimkannya
ke reasuradur tujuan.

⛔ **BLOCKED — menunggu isi `M_LINK_SERVICE`.** `CLAUDE.md` §4.4: daftar endpoint sesungguhnya ada di
tabel Oracle `M_LINK_SERVICE` (kunci `KATEGORI_1` + `KATEGORI_2`), yang **isinya tidak ada di korpus**.
Tanpa itu, tujuan kirim tidak dapat ditentukan, dan endpoint **tidak boleh** ditulis literal.

⚠️ **Tiket ini tidak dihapus.** Lingkupnya sudah diketahui; menghilangkannya akan menyembunyikan
pekerjaan yang pasti ada. Ini **pola yang sama dengan `..\rnw\R07-konversi-produksi-renewal.md`** —
tiket ditulis penuh, hanya eksekusinya menunggu.

`[terverifikasi]` Yang ada di korpus:

| Bagian | Bukti |
| --- | --- |
| Shape alur | `DDL\OfferFacRetro.xml` `<pyMOName>`: `PRINT R/I SLIP`, `Send Email Print RI Slip`, `IsRISlip` (2×), `PrintRISlip_FlowAction` |
| Activity PDF | `GeneratePDFFacOutMemoPlacing` · `GeneratePDFFacOutOfferStatus` · `GeneratePDFFacOutViewOffer` · `DeletePDFFacOutMemoPlacing` |
| Section cetak | `FacOutPrintRISlipSectionInside` · `PrintRISlip` · `Harness\PrintRISlips` |
| Data slip | `.OfferFacIn.FacRetroDetails.DocumentPosition` (F05) · nomor slip dari **F09** |
| Tipe email | `OfferFacOut_PostAct` membaca `EmailTypeRetro` / `EmailTypeRetroSlip` bernilai `1` / `2` / `7` |

⚠️ `[pertanyaan terbuka]` Arti `1` / `2` / `7` pada tipe email **tidak dijelaskan korpus**. Nilainya
diport apa adanya; nilai di luar ketiganya → **`panic`** (`CLAUDE.md` §4.5).

⛔ **Empat activity PDF ada, tetapi isi templatnya tidak.** `[dugaan]` Templat dokumen kemungkinan
berada di luar ekspor rule. Selama belum terbukti ada, bentuk keluaran PDF **tidak boleh ditebak**.

⛔ **Alamat email tujuan tidak pernah disalin** ke tiket, test, fixture, log, maupun dokumen (K-025).
Yang dicatat hanya **jumlah dan mekanismenya**.

**Asal (Pega).** `DDL\OfferFacRetro.xml` (shape cetak & kirim) · `Activity\GeneratePDFFacOutMemoPlacing` ·
`GeneratePDFFacOutOfferStatus` · `GeneratePDFFacOutViewOffer` · `DeletePDFFacOutMemoPlacing` ·
`Section\FacOutPrintRISlipSectionInside` · `Activity\OfferFacOut_PostAct`

**Keputusan.** **K-061** (ditulis penuh, ditandai menunggu, pola R07 — **tidak dihapus**) · K-055 ·
K-062 (pemanggil `PrintRISlipPre_act` menghasilkan kosong) · K-025 · `CLAUDE.md` §4.4, §4.5, §4.6

**Blocked by:** F08 · ⛔ **eksternal — isi `M_LINK_SERVICE`**

**Status:** blocked

- [ ] ⛔ Menunggu isi `M_LINK_SERVICE`; **tidak ada endpoint literal** di kode
- [ ] R/I Slip memuat nomor slip dari **F09** dan posisi dokumen dari **F05**
- [ ] Enumerasi tipe email `1`/`2`/`7` diport apa adanya; nilai lain → **`panic`**
- [ ] ⛔ Bentuk keluaran PDF **tidak ditebak** selama templatnya belum terbukti ada di korpus
- [ ] ⛔ **Alamat email dan nama orang tidak pernah** muncul di artefak mana pun (K-025)
- [ ] Endpoint dan kredensial dari konfigurasi (§4.4)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

✅ **K-061 sudah diputuskan: pola R07.** Tiket ditulis penuh, **ditandai menunggu**, **tidak dihapus**.

⛔ **Blocker eksternalnya, dinyatakan tegas:** isi tabel Oracle **`M_LINK_SERVICE`** (kunci
`KATEGORI_1` + `KATEGORI_2`) **tidak ada di korpus**. Tanpa itu tujuan kirim tidak dapat ditentukan,
dan `CLAUDE.md` §4.4 melarang endpoint literal. **Pemiliknya DBA.** Sampai itu tersedia, F13
**tidak dapat dieksekusi** — hanya ditulis.

⚠️ **K-062 menyentuh tiket ini:** `Activity\PrintRISlipPre_act` adalah salah satu dari tiga pemanggil
`GetOPFacOut_Sql`. Klep itu **dihapus**, sehingga pemanggil di jalur cetak ini **wajib menghasilkan
nilai kosong, bukan galat** — lihat **F12**.
