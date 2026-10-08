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

## 10 · Blok berjudul `hiddden` — ✅ TERJAWAB 7 Oktober 2026

Tab Share NP dan Value Difference memuat blok `pyHeaderType=BAR` berjudul `hiddden`
(ejaan ekspor) berisi empat grid total. **Jawaban dari ekspor:** blok itu `pyIncludeHeader = false`
— judulnya TIDAK tampil; gridnya tampil. Penentunya `pyIncludeHeader`, bukan
`pyContainerFormat = NOHEADER` (blok Exclusions ber-NOHEADER dan judulnya tampil di gambar Pega 39).
Lihat `LAYAR-ADJUSTMENT.md` §9.

## 11 · Rincian baris `…!pyGridRowDetails`

~~Grid Installment (`ASM-FW-GISFW-Data-TreatyInInstallment!pyGridRowDetails`) dan Limits
membuka rincian lewat rule yang TIDAK ada di ekspor. Data `InstallmentList` ada di dokumen,
tetapi kolomnya tidak dikarang. `Installments.xml` di ekspor gridnya mati (`1=2`) dan bukan
rule yang dirujuk.~~

⛔ **RALAT 7 Oktober 2026 — catatan di atas SALAH, dua kali.** (1) `…!pyGridRowDetails` hanya
templat bingkai; isinya flow action di `pyGridProps/pyEditAction` (`Installments`, `Layers`,
`Share`, `MaxRetention`, `DetailEGNPI`, …), dan flow action itu ADA di ekspor, menunjuk Section
lewat `pySectionReference`. (2) `1=2` di `Installments.xml` adalah `pyContainerVisibleWhen`
wadah luar yang DITIMPA `pyIsVisibilityOption = ALWAYS` — gridnya tampil (gambar Pega 38).
Catatan itu ditulis sebelum aturan ALWAYS diketahui. Kini dibangkitkan: `LAYAR-ADJUSTMENT.md` §9.

## 12 · Activity penghitung total

`Total Retention`, `Total EGNPI`, `Total All Layers`, `Total ROL`, dan sembilan total RNM
Share menampilkan nilai TERSIMPAN. Activity-nya belum ditelusuri ronde ini.

## 7 · `ProportionType` di grid daftar

Kolom `Reinsurance Type` menampilkan nilai tersimpan apa adanya (`NonProportional`,
tanpa spasi). Modul `treatyin` menerjemahkannya menjadi `Non Proportional` di
services-nya; modul ini tidak boleh mengimpornya, dan menulis penerjemah kedua
dilarang. **Pertanyaan:** perlukah aturan tampil itu dipindah ke tempat bersama?


## 13 · Tombol Add dan `Choose` — draf, bukan simpanan (7 Oktober 2026)

Rantai ekspornya: `Add Revision` @500554 / `Add Adjustment Premium` @515429 →
DataTransform `TreatyCreateEDM` (EDMState `"1"` / `"3"`) → modal picker → `Choose` →
Activity `TreatyInEDMSetValue` → `closeContainer`.

**Keputusan pemilik proses:** `Choose` menyusun **draf di layar**. Langkah [8]
`TreatyRevisionCopyAttachment` dan [9] `SaveTreatyIn_EDM_Act` — keduanya tulisan — menunggu
jalur Save. Panel Attachment draf memperlihatkan lampiran **asalnya** (yang [8] akan salin).

Yang disalin **apa adanya**, dan dicatat supaya tidak terbaca sebagai cacat aplikasi:

| Hal | Ekspor | Di sini |
| --- | --- | --- |
| Pengenal baru, belum ada revisi | `ID + "/R01"` ([4]) | sama |
| Pengenal baru, sudah ada | `@If(n < 10, dasar+"/R0"+(n+1), dasar+"/R"+(n+1))` ([5]) | sama — `R09` → `R010`; master 7 aksara yang sudah punya revisi → `R01` lagi. Pemeriksa duplikat [7] ber-`//`. Diukur: 263 `R01` ber-OLDID master, 16 `R02` + 1 `R03` ber-OLDID revisi sebelumnya |
| Ada revisi? | `GetTreatyRevisionID` atas `M_TREATY_IN_EDM.ID` | `TREATY_IN_EDM.ID` — pengenal sama (280/280), tanpa tabel ber-JSON |
| `EDMMaterialType` | `Param.MaterialType` dan `@if(InternalType=3,"1",…)` | `Param.MaterialType` — hasil keduanya sama, picker Premi mengirim `'1'` |
| `OLDDATA` | salinan `TreatyIn` SESUDAH [3]–[5], SEBELUM pengenal berganti | sama |

**Yang masih terbuka:**

1. **Syarat tampil `Add Adjustment Premium`** — `OperatorID.pyWorkBasketList(2).pyWorkBasketName =
   'ReasTreatyInAdmin'`. Pemetaan workbasket ke peran aplikasi belum ada (sama dengan §3), jadi
   tombolnya tampil untuk semua.
2. **Teks radio.** `Internal` / `External` diambil dari judul picker di Section yang sama
   (`EDMState = '1'` → "…Internal Revision", `'2'` → "…External Revision"); daftar `associated`
   aslinya tidak diekspor. `Material Type` tampil sebagai kode `1` / `2` (§1).
3. **Kepala dokumen yang belum mendarat.** Tabel pendaratan sedang dikosongkan; kunci kepala
   yang tidak ada di dokumen dilengkapi dari kolom `TREATY_IN` / `TREATY_IN_EDM` — tanpa itu
   `ProportionType` kosong dan panel Old/New tidak terbuka.

## 14 · Keputusan ronde rumus rincian (7 Oktober 2026)

1. **Penggabung baris `pyActionConditions`.** `pyLogic` tersimpan di setiap baris (bawaan `And`)
   dan arah sambungnya tidak terbaca dari korpus (claimprop grilling-ronde-2 R7-baru). Dibaca:
   semua baris membandingkan properti yang SAMA dengan `=` → `||`, karena
   `.Note = 'SURPLUS' && .Note = '2ND SURPLUS'` tidak pernah benar; selain itu `&&`. Treaty In
   membacanya sama (`SURPLUS_MATA_UANG`/`syaratSurplus`; langkah EDM `7897987 && EDMState != 3`).
2. **`SetCurrName_Act` mengisi pasangan nama ↔ pengenal.** Activity-nya menulis
   `MstBankAccount.CURRENCY` bila `Param.CURRID == .ID`, yaitu halaman lain, dan `CURRID` tidak
   dikirim sel mana pun. Atas barisnya sendiri Activity itu tidak mengubah apa pun. Data tersimpan
   Pega tetap berpasangan, jadi pasangannya diisi medan tambahan dropdown yang tidak diekspor. Di
   sini pasangan diisi dari daftar mata uang — sama dengan Treaty In (`PilihKode`).
3. **Kind of Treaty = SOA Name** (Name bila kosong) — keputusan pemakai 6 Oktober 2026 untuk Treaty
   In (§20 di sana), diterapkan juga di panel New karena ia form yang sama. Menyimpang dari langkah
   3 `SetTreatyTypeName_Act` (`.Note`).
4. **`AddDeduction` BERBEDA antarkorpus.** Adjustment: `Currency = .Currency`; Treaty In:
   `Currency = .CurrencyIOOLimit`. Panel ini mengikuti korpusnya sendiri.
5. **`TreatyInXOLAddSpreadingDetailActual` belum dibangun.** Activity ini menulis halaman
   `TreatyIn.ActualValue.Share(idx)` (EDM Actual), yang tidak ada di model layar. Ia bersyarat
   `TreatyIn.EDMState = 3`. Bila syarat itu benar, `% RNM Share` rincian Share membatalkan seluruh
   rantainya dengan pesan; bila salah, langkah itu dilewati seperti di Pega.
6. **Master belum dimuat.** Nama Treaty Group/Type/mata uang tidak diisi dan pesan tampil, bukan
   dikosongkan.

## 15 · Keputusan ronde sisa (7 Oktober 2026, lanjutan)

1. **Payment Date — tidak berbukti, diputuskan.** WPC kosong/bukan bilangan bulat = 0 hari
   (`+1` tetap); Due Date kosong/tak terbaca → Payment Date DIBIARKAN. `DueDate`/`PaymentDate`
   bertipe tanggal 8 aksara, jadi `+WPC+1` dibaca harfiah (tanpa geseran zona waktu GMT).
2. **Halaman sesi di akar panel.** `SearchData.*`, `FlagExcel.*`, `TempQuarter*` disimpan berkunci
   utuh di akar panel New. ⚠️ Jalur Save (sesi lain) harus MENYARINGNYA — Treaty In sudah
   (`simpan_kontrak.go` `kunciSesi`).
3. **Cabang Adjust Premium dipilih oleh `EDMState = 3`** pada akar panel New, bunyi syarat
   @566980. Panel Old tetap `TreatyInTabsNonProportionalOldData`.
4. **Submit Achievement tidak dibangun di sini.** Ia menulis `POOLDATA.LOG_ACHIEVEMENT`; menulis
   langsung ke skema warisan dari layar baru bertentangan dengan arah migrasi 423 dan dengan
   kepemilikan jalur tulis (Save, sesi lain). ✅ **KEPUTUSAN PEMILIK PROSES 7 Oktober 2026:
   tombolnya TETAP MATI** sampai tabel log diputuskan bersama jalur Save. Angka Achievement
   (`AchievementLists`) sendiri tersimpan lewat Save (`T_TREATY_LIMIT_ACHIEVEMENT`). Tombolnya
   menyebut alasan ini (`PENYESUAIAN.achievementMenunggu`).
