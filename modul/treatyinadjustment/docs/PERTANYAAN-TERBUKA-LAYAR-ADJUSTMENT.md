# Pertanyaan terbuka — layar Adjustment

Dibuka 5 Oktober 2026. Tiap butir: apa yang tidak diketahui, apa yang layar lakukan
sementara, dan apa yang mengubahnya. Ditujukan ke pemilik proses Treaty In Adjustment.

## 1 · Teks pilihan `Adjustment Type` (`EDMState`) dan `Material Type` (`EDMMaterialType`)

Kedua dropdown memakai `pyListSource=associated` — daftar milik rule properti yang
**tidak ikut diekspor**. Petunjuk makna yang ADA di korpus (judul pemilih, bukan teks
dropdown):

| Kode | Petunjuk | Sumber |
|---|---|---|
| `EDMState 1` | *Choose Master to Create Internal Revision* | `Section/PickerTreatyInMasterRevisi.xml` |
| `EDMState 2` | *Choose Master to Create External Revision* | idem |
| `EDMState 3` | *Choose Master to create Premium Adjustment*; label `Adjustment Premium` menggantikan dropdown | `PickerTreatyInMaster.xml`, `InputTreatyInAdjustment.xml` @82379 |
| `EDMMaterialType 1` / `2` | tidak ada | — |

**Sementara:** layar menampilkan KODE bertanda "teks pilihannya tidak ada di ekspor",
kecuali `EDMState 3` yang memang tampil sebagai label. Grid daftar (kolom Type,
Material Type) juga menampilkan kode — selnya mengikat `.CARI1` sebuah halaman bantu
(@728375, @734798) yang isinya tidak tersimpan.
**Yang mengubahnya:** teks pilihan resmi dari pemilik proses, atau ekspor rule
properti `EDMState`/`EDMMaterialType`.

## 2 · Pengenal lampiran penyesuaian

`WorkAttachments` membaca data page `KategoriDocument` berparameter `IDDoc`; nilainya
tidak terlihat di ekspor. `M_ATTACHMENTTREATY_2` punya **0** baris berpengenal
penyesuaian (`…/Rnn`). **Sementara:** panel dibaca dengan `TreatyIn.ID` (pengenal
revisi) dan tampil `No items`. **Pertanyaan:** apakah lampiran penyesuaian semestinya
memperlihatkan lampiran kontrak asalnya (`OLDID`)?

## 3 · Syarat peran tombol `Actions`

`OperatorID.pyPosition= 'IT Developer' || OperatorID.pyWorkBasketList(1).pyWorkBasketName = 'ReasTreatyInSecHead' || …`
(@155994). Pemetaan workbasket Pega ke peran aplikasi ini belum ada. **Sementara:**
tombol tampil, mati.

## 4 · ✅ DIJAWAB — desimal per kolom mengikuti ekspor

Keputusan pemilik proses 5 Oktober 2026: desimal ditentukan PER KOLOM dari sumbernya, nol di
ekor dipertahankan. `ProRatePercent` kini 4, uang di grid 2 bila ekspor menyatakannya. Tabelnya
di `LAYAR-ADJUSTMENT.md` §6.2.

## 5 · Desimal kolom yang TIDAK dinyatakan ekspor

Tercatat di `LAYAR-ADJUSTMENT.md` §6.2 (baris "tidak dinyatakan"): antara lain `Conversion`,
`AmountIDR`, `Limit`/`MDP`/`Deductible`/`NetPremi` di grid Summary, `Value` grid RNM Share,
`RNMShare`, `BrokeragePercent`, `FacultativeShare*`, `RSMDLimit` dan kawan-kawan Event Limits.
**Sementara:** aturan lama (uang 4, persen 2, nol ekor dibuang). **Yang mengubahnya:**
`pyDecimalPlaces` rule Property-nya — ditagih bersama butir 1.

## 6 · ✅ KEPUTUSAN — panel Old nol medan tersunting

`TeritorialScope` @43977 dan `BordereauxNote` @56021 di `TreatyInNONProportionalOldData.xml`
ber-`pyReadOnly = false` di ekspor. Layar ini membuat **seluruh** panel Old baca-saja, dan
pemilik proses 5 Oktober 2026 **membenarkannya**: Old Data yang dapat diketik adalah cara
tercepat merusak jejak perubahan. ⛔ Yang membaca ekspor nanti: ini BUKAN kelalaian salin —
jangan "diperbaiki".

## 8 · Golongan `MaxCoGroup` / `MaxCoNonGroup`

Tab Co-Ins Scale P: `Max Co-Insurance Panel (Non Group/Group)`, `pxNumber` tanpa desimal.
Terukur hanya bernilai `10`. Uang, cacah, atau persen — tidak dinyatakan. **Sementara:**
apa adanya (golongan `teks`).

## 9 · Teks pilihan dropdown di dalam tab

`OptionLimit` (kode `1`/`2`), `ReportingPeriod` (`quarter`), `AccumulationPeriod`
(`none`/`quarter`): daftar pilihannya tidak ada di ekspor. **Sementara:** nilai tersimpan.

## 10 · Blok berjudul `hiddden`

Tab Share NP dan Value Difference memuat blok `pyHeaderType=BAR` berjudul `hiddden`
(ejaan ekspor) berisi empat grid total. Kepala BAR menampilkan judulnya, jadi layar ini
menampilkannya apa adanya. **Pertanyaan:** apakah di layar lama judul itu memang terlihat?

## 11 · Rincian baris `…!pyGridRowDetails`

Grid Installment (`ASM-FW-GISFW-Data-TreatyInInstallment!pyGridRowDetails`) dan Limits
membuka rincian lewat rule yang TIDAK ada di ekspor. Data `InstallmentList` ada di dokumen,
tetapi kolomnya tidak dikarang. `Installments.xml` di ekspor gridnya mati (`1=2`) dan bukan
rule yang dirujuk.

## 12 · Activity penghitung total

`Total Retention`, `Total EGNPI`, `Total All Layers`, `Total ROL`, dan sembilan total RNM
Share menampilkan nilai TERSIMPAN. Activity-nya belum ditelusuri ronde ini.

## 7 · `ProportionType` di grid daftar

Kolom `Reinsurance Type` menampilkan nilai tersimpan apa adanya (`NonProportional`,
tanpa spasi). Modul `treatyin` menerjemahkannya menjadi `Non Proportional` di
services-nya; modul ini tidak boleh mengimpornya, dan menulis penerjemah kedua
dilarang. **Pertanyaan:** perlukah aturan tampil itu dipindah ke tempat bersama?
