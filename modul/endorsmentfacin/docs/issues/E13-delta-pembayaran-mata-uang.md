# E13: Selisih pembayaran tingkat mata uang

**What to build:** Untuk setiap mata uang pada kasus endorsement, hitung **selisih pembayaran** —
nilai sesudah dikurangi nilai sebelum.

`[terverifikasi]` Rumus kanoniknya tunggal dan berulang di belasan penugasan. Nilai "sebelum"
diambil dari **lapis A** (dokumen polis lama), bukan dari nilai lama per baris.

⛔ **Netto menambahkan PPh dan PPN, tidak mengurangkannya.** Terlihat berlawanan intuisi pajak;
diport apa adanya.

⚠️ **Urutan mengikat:** transformasi perhitungan premi **menimpa** porsi periode yang disetel modul
before-image sebelum jalur produksi membacanya (K-048).

**Asal (Pega).** `Activity\CountEndorsementData` · `CountDataEDMElse` · `CountPaymentEdm_Act` ·
`CountPaymentEdmTSIObj_Act` · `DataTransform\CountPremiEDM_DT`

**Keputusan.** K-048 · K-046 · K-010/K-012

**Blocked by:** E12 · E06

**Status:** sebagian — 09-10-2026. `CountPaymentEdm_Act` diport dari langkah 1 sampai blok 15 kecuali
langkah 2 (`SetRIComm_Act`, query master pajak — E17 dilewati) dan 6 (`GetEdmProdKe_Act`, query; hasilnya
masukan `HasilProdKe`): `porsipembayaran.go` (porsi, langkah 1/3/4/7/9/11, blok 13–15), `pembayaran.go`
(pembayaran per mata uang 13.2/14.4/15.5 + `CountPaymentEdmTSIObj_Act` 1.4/1.5/2.2/langkah 3 cicilan),
`predikatkasus.go` (`KasusSesudahRISlip`, langkah 5), langkah 10 di `premiedmdt.go` (E14). `[terverifikasi]`
Rekonsiliasi kasus nyata NB-15 blok 13: komponen lewat lokal `Decimal` EKSAK; lewat lokal `double` sama
sesudah pembulatan double (A19 — Pega menghitung `preminew`/`oldpayment`/`premiins`/`ricommins` sebagai
`double`; port tetap desimal eksak, CLAUDE.md §7). ⛔ Menyimpang dari tiket, mengikuti korpus: "mata uang
baru → seluruh nilai sebelum dipaksa nol" — korpus 13.2/15.5 tidak menolkan apa pun, 14.4.4 tiga medan,
TSIObj 1.4.4 empat medan, `CountPremiEDM_DT` EdmType 4 tiga medan (EdmType 1/2 lewat penolan Param);
komponen lain TERBAWA dari baris sebelumnya (A20, `BarisTerbawa`). Langkah 10 dirangkai di tempatnya lewat
`MasukanPorsiPembayaran.LangkahSepuluh` (`TestRangkaianLangkahSepuluh`). Langkah 5 menulis `Type = 0` ke
halaman kerja di Pega; di sini hanya salinan untuk predikat — penulisan balik belum diport.

- [ ] Selisih pembayaran = nilai sesudah − nilai sebelum, per mata uang
- [ ] Netto = premi − komisi − brokerage − potongan **+ PPh + PPN**
- [ ] Nilai sebelum bersumber **lapis A**, bukan lapis B
- [ ] Lima komponen "lama" tingkat pembayaran ikut terisi: komisi · PPN · PPh · brokerage · potongan
- [ ] Mata uang yang **tidak ada** di polis lama → seluruh nilai dianggap **baru**, nilai sebelum dipaksa nol
- [ ] Seluruh nilai bertipe `Money` dengan mata uangnya
- [ ] ⚠️ Urutan eksekusi terhadap modul before-image **dijaga** (K-048)
- [ ] **K-046** `K046_NettoPremi_PPhPPNDitambahkan`
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
