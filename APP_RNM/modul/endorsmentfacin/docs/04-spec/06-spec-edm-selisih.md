# Bahan spec — Modul 2: Perhitungan Selisih Endorsement (EDM)

> **Ini BAHAN untuk `/to-spec`, bukan spec.** `/to-spec` ber-`disable-model-invocation: true` dan
> **belum dijalankan**. Work owner yang menjalankannya manual.
>
> **Sumber:** `07-edm\04-e3-before-image-selisih.md` §3 · `07-edm\06-e5-jalur-produksi.md` §3.3 ·
> arsip `01-flow\03-alur-endorsement.md` §3. Prior art bentuk: `04-spec\05-spec-edm-before-image.md`
> (modul 1) dan `04-spec\03-spec-modul-terverifikasi.md`.
>
> **Keputusan mengikat:** K-010/K-012 (uang = `Money`) · K-018 (COB→skala lewat resolver) ·
> K-027 (koma desimal) · K-044 · K-046 (kode usang diport apa adanya, **kecuali A.5**) ·
> **K-047 (seam: modul ini TIDAK dapat seam baru)**.
>
> ⚠️ **Dua temuan baru** muncul saat merapikan bahan — keduanya butuh keputusan work owner. Lihat §8.
> Tidak saya putuskan sendiri.

---

## 0. Posisi modul ini

Modul 1 (before-image) menghasilkan **nilai lama**. Modul 2 memakainya untuk menghasilkan **selisih**.

`[terverifikasi]` Sambungannya nyata di korpus, bukan sekadar konseptual: `EDMOldPayment` diisi dari
**`pyWorkPage.OfferFacIn.OldData.CurrencyList(<CURRENT>).SumTotalPayment`** — yaitu **lapis A** milik
modul 1 — pada **12 dari 21** penugasannya.

⛔ Keluaran modul ini adalah **selisih**. Bahwa tangga akseptasi memakainya (arsip §4.3) **tetap
`[belum diuji]`** dan **tidak boleh** dinyatakan di spec — arsip §4 belum ditinjau.

---

## 1. Dua tingkat delta yang berbeda tujuan

`[terverifikasi]` Korpus menghitung selisih pada **dua tingkat berbeda**, dengan berkas, rumus, dan
konsumen berbeda. Menyatukannya akan menghasilkan angka yang salah di salah satu sisi.

| | **Tingkat A — delta PEMBAYARAN** | **Tingkat B — delta BARIS PRODUKSI** |
| --- | --- | --- |
| **Satuan** | per **mata uang** (`CurrencyList`) | per **baris spreading** (per objek × coverage × treaty) |
| **Dihitung oleh** | `CountEndorsementData` · `CountDataEDMElse` · `CountPaymentEdm_Act` · `CountPaymentEdmTSIObj_Act` · `ReCountPremiLifeEDM` · `CountPremiEDM_DT` | enam `SaveFacinProd*EDM*_Act` (Fire, Aneka, Golf, Marine Cargo, MBU, PA) |
| **Mengisi** | `CurrencyList.Policy.Payment.EDM*` dan `CurrencyList.SumTotalPayment` | parameter `Datain.*` / `Datain1.*` → kolom `facinproduction` |
| **Konsumen** | layar pembayaran + ringkasan endorsement | tabel produksi (modul berikutnya) |
| **Nilai lama dari** | lapis A — `OldData.CurrencyList(...).SumTotalPayment` | lapis A — `OldData...SpreadingList` dengan **`TreatyType` sama** |

📌 Keduanya **sama-sama** membaca lapis A, bukan lapis B. Properti `*Old` per baris (lapis B) melayani
**layar**, bukan kedua perhitungan ini.

### 1.1 Keenam berkas tingkat A — terverifikasi ada

`[terverifikasi]`

| Berkas | `pxObjClass` | `pyClassName` |
| --- | --- | --- |
| `Activity\CountEndorsementData` | `Rule-Obj-Activity` | `ASM-SFAGIS-Work-Endorsement` |
| `Activity\CountDataEDMElse` | `Rule-Obj-Activity` | `ASM-SFAGIS-Work-Endorsement` |
| `Activity\CountPaymentEdm_Act` | `Rule-Obj-Activity` | `ASM-FW-GISFW-Work` |
| `Activity\CountPaymentEdmTSIObj_Act` | `Rule-Obj-Activity` | `ASM-FW-GISFW-Work` |
| `Activity\ReCountPremiLifeEDM` | `Rule-Obj-Activity` | **`Data-Party-Person`** |
| **`DataTransform\CountPremiEDM_DT`** | **`Rule-Obj-Model`** | `ASM-FW-GISFW-Work` |

⚠️ Yang keenam bertipe **`Rule-Obj-Model`**, bukan Activity — strukturnya `pyProperties` bersarang,
bukan `pySteps`. Dipanggil dari `Activity\CountPaymentEdm_Act` dan `Section\ViewOldDataPayment_IsUW`.

---

## 2. Rumus kanonik tingkat A

### 2.1 Selisih pembayaran

`[terverifikasi]` **19 penugasan** memakai bentuk yang sama (12 berspasi + 7 tanpa spasi — identik
secara semantik):

> ### **`SumTotalPayment = EDMPremiMenjadi − EDMOldPayment`**

Tiga penyimpangan `[terverifikasi]`, diport apa adanya:

| Berkas · baris | Bentuk |
| --- | --- |
| `ReCountPremiLifeEDM` L3248 | `.Policy.Payment.Premium − Local.NetPremiOld` |
| `ReCountPremiLifeEDM` L627 | `0` |
| `CountDataEDMElse` L6299 | `CurrencyList(<LAST>).Policy.Payment.Premium − CurrencyList(<LAST>).Policy.Payment.Commision` |

Ditambah **dua bentuk berprorata** di `CountPremiEDM_DT` — lihat §2.4.

### 2.2 Definisi netto — `EDMNewPremi`

`[terverifikasi]` Bentuk utama, **4 kemunculan** (`CountPaymentEdm_Act` L3285 · L5614 · L8481 ·
`CountPaymentEdmTSIObj_Act` L1743), dan bentuk setara di `CountPremiEDM_DT` 10.2.2 / 12.3.6:

```
EDMNewPremi = Premium − Commision − BrokerageFee − Deduction2 + PPh + PPN
```

`[terverifikasi]` Varian lain: **11× disetel `0`** (blok penihilan) · `0 − OldData.CurrencyList(<CURRENT>).SumTotalPayment`
(`CountDataEDMElse` L4341) · `Local.TotalNewPremiNusantaraRe` (`CountDataEDMElse` L6213).

⛔ **PPh dan PPN DITAMBAHKAN, bukan dikurangkan.** Terlihat berlawanan intuisi pajak; diport apa
adanya.

### 2.3 Nilai lama — `EDMOldPayment`

`[terverifikasi]` Lima bentuk:

| Bentuk | Cacah | Sumber |
| --- | ---: | --- |
| `OldData.CurrencyList(<CURRENT>).SumTotalPayment` | **12** | **lapis A** |
| `Local.oldpayment` | 5 | lokal, diisi dari lapis A di langkah sebelumnya |
| `0` | 2 | mata uang tidak ada di polis lama |
| `Local.NetPremiOld` | 1 | jalur Life |
| `OldData.CurrencyList(Local.Index).SumTotalPayment` | 1 | lapis A, indeks eksplisit |

### 2.4 Dua rasio prorata — peran berbeda

`[terverifikasi]` `ProrateEDMEnd` = porsi **sesudah** tanggal endorsement, pengali umum.
`ProrateStartEDM` = porsi **sebelum**, dipakai terbatas.

⛔ **Tetapi keduanya dihitung di DUA tempat berbeda dengan rumus berbeda** — lihat §8.1. Ini temuan
baru yang butuh keputusan.

---

## 3. Tabel `EDMPremiMenjadi` per (lini × `EdmType` × jalur)

⛔ **Inilah inti "tidak seragam".** `[terverifikasi]` **22 penugasan · 12 rumus literal · 11 rumus
semantik** di enam berkas. **JANGAN diseragamkan** — modelkan sebagai tabel, bukan satu fungsi.

Arti `EdmType` mengikuti **K-029**: `1` = Batal Sejak Semula · `2` = Batal Prorata ·
`4` = Penambahan/Pengurangan/Perubahan · `3` usang (cabang tetap diport).

### 3.1 Kelompok I — jalur "batal" (`CountEndorsementData` / `CountDataEDMElse`)

| Berkas | Langkah | Lini (gerbang) | Gerbang `EdmType` | Label langkah | Rumus |
| --- | --- | --- | --- | --- | --- |
| `CountEndorsementData` | 1.2.1 | `IsFire` | `==2` | *"batal sejak semula"* | `EDMOldPremi + EDMNewPremi` |
| `CountEndorsementData` | **1.3.1** | `IsFire` | `==2` | *"batal"* | ⚠️ **`EDMOldPayment + EDMNewPremi`** |
| `CountEndorsementData` | 2.2.1 · 2.3.1 | `IsAneka` | `==2` | keduanya | `EDMOldPremi + EDMNewPremi` |
| `CountEndorsementData` | 3.2.1 · 3.3.1 | `isGolfInsurance` | `==2` | keduanya | `EDMOldPremi + EDMNewPremi` |
| `CountDataEDMElse` | 1.2.1 · 1.3.1 | `IsMarineCargo` | `==2` | keduanya | `EDMOldPremi + EDMNewPremi` |
| `CountDataEDMElse` | 2.2.1 · 2.3.1 | ⚠️ `IsMarineCargo` **(label "MBU business")** | `==2` | keduanya | `EDMOldPremi + EDMNewPremi` |
| `CountDataEDMElse` | 3.2.1 | `IsLife` | ⚠️ **`==1 \|\| ==2`** | *"batal sejak semula"* | `EDMOldPremi + EDMNewPremi` |
| `CountDataEDMElse` | **3.3.1** | `IsLife` | `==2` | *"batal"* | ⚠️ **`EDMOldPayment + EDMNewPremi`** |
| `CountDataEDMElse` | 3.4.1.3 | `IsLife` | `==2` + `IsSameCurrency==1` | *"batal"* | `Local.TotalNewPremiNusantaraRe + Local.TotalOldPremiNusantaraRe` |

📌 **Pembagian `EDMOldPayment` vs `EDMOldPremi` bukan sifat lini.** Ia berlaku **hanya pada cabang
`x.3` (*"batal"*)**, dan di sana hanya **FIRE dan LIFE** yang memakai `EDMOldPayment`. Pada cabang
`x.2` ketujuh lini memakai `EDMOldPremi`.

### 3.2 Kelompok II — jalur per **jenis endorsement** (`CountPaymentEdm_Act`)

`[terverifikasi]` Bercabang menurut predikat jenis endorsement, **bukan** `EdmType`:

| Langkah | Gerbang luar | Gerbang dalam | Rumus |
| --- | --- | --- | --- |
| 12.1 | `EdmTypeNew==4` *("EDM RI SLIP")* | — | `OldData.CurrencyList(Index).SumTotalPayment` |
| 13.2.3 | `IsEdmExtendPeriod` | `IsFire` ⛔ prakondisi nonaktif | netto × `@Math.divide(datedif, day, 20)` |
| 14.4.3 | `IsEdmAdjRate` | `IsFire` ⛔ nonaktif | idem **`+ (oldpayment × @Math.divide(datedifbefore, day, 20))`** |
| 14.4.4 | `IsEdmAdjRate` | `IsFire` + `Local.check==0` | `(preminew − Commision) × @Math.divide(datedif, day, 20)` |
| 15.5.3 | `IsEdmAdjPeriod` | `IsFire` ⛔ nonaktif | netto × `@Math.divide(edmdate, startdate, 20)` |

Predikat jenis endorsement yang dipakai: `IsEdmExtendPeriod` · `IsEdmAdjRate` · `IsEdmAdjPeriod` ·
`IsEdmAdjShareCedant` · `IsEdmAdjTSI` · `IsEdmAdjInsured`.

### 3.3 Kelompok III — sisanya

| Berkas | Langkah | Gerbang | Rumus |
| --- | --- | --- | --- |
| `CountPaymentEdmTSIObj_Act` | 1.4.3 | `IsEdmAdjTSI` + `CalcBrokerFee=="true"` | netto × `@Math.divide(datedif, day, 20)` |
| `CountPaymentEdmTSIObj_Act` | 2.2.3 | `IsEdmAdjInsured` | `Local.oldpayment` |
| `ReCountPremiLifeEDM` | 4.1.9.3 | **(tanpa prakondisi)** | `.Policy.Payment.Premium` |
| `CopyAllObj_ACT` | 9.2.2.1 | `Local.idxOcc==.Name` | `EDMPremiMenjadi + Local.Ujrah` — **akumulator** |

### 3.4 Kelompok IV — `CountPremiEDM_DT` (baru, belum ada di E-3)

`[terverifikasi]` Bercabang murni menurut `EdmType`:

| Langkah | `EdmType` | `EDMPremiMenjadi` | `SumTotalPayment` |
| --- | --- | --- | --- |
| 10.2.7 | `==4` | `Param.actualpremi` (= netto) | `Param.Premi` = `(actualpremi × prorateedmend) − (OldPremi × prorateedmend)` |
| 10.2.16.1 | `==4` + `Param.Check=="0"` | `EDMNewPremi` | `(EDMPremiMenjadi × prorateedmend) − (EDMOldPayment × prorateedmend)` |
| 11.4.10 | `==1` | netto penuh | `(EDMPremiMenjadi × prorateedmend) − (EDMOldPayment × prorateedmend)` |
| 12.3.10 | `==2` | **`Param.OldPremiProrate × Param.prorateedmstart`** | `EDMPremiMenjadi − EDMOldPayment` |
| 1.1.8 | `==4` | `Param.PremiMenjadi` | ⛔ **seluruh langkah 1 `pyDisabled=true`** |

⛔ **`Param.Check`** adalah penanda ketemu-tidaknya mata uang yang sama di `OldData.CurrencyList`.
`"0"` = tidak ketemu → seluruh nilai dianggap **baru**, `EDMOldPayment` dipaksa `0`.

⛔ **Langkah 1 seluruhnya `pyDisabled=true`** — implementasi lama `EdmType==4` yang digantikan
langkah 10. Diport apa adanya sebagai cabang mati? → **pertanyaan §8.2**.

### 3.5 Komponen "lama" yang belum tercatat modul 1

`[terverifikasi]` `CountPremiEDM_DT` mengisi **lima properti lama tingkat pembayaran** di luar
`EDMOldPayment`/`EDMOldPremi`:

```
EDMOldCommision · EDMOldPPN · EDMOldPPh · EDMOldBrokerage · EDMOldDeduction2
```

Kelimanya diisi dari mata uang lama yang cocok namanya (`Param.Curr == .Name`). Ini **tingkat
pembayaran**, bukan lapis B per baris — jadi milik modul 2, bukan modul 1.

---

## 4. Delta baris produksi (tingkat B)

### 4.1 Rumus

`[terverifikasi]` `Activity\SaveFacinProdEDMFire_Act.xml` — **57 penugasan delta berbentuk
pengurangan**:

```
# 1) nilai baru diprorata lebih dulu
Local.TsiSpreadEDM   = .TSISpreaded     × pyWorkPage.OfferFacIn.ProrateEDMEnd
Local.PremiSpreadEDM = .PremiumSpreaded × pyWorkPage.OfferFacIn.ProrateEDMEnd

# 2) lalu dikurangi nilai lama berjenis treaty SAMA
Local.TsiSpreadEDM   = Local.TsiSpreadEDM   − Local.TsiSpreadNB
Local.PremiSpreadEDM = Local.PremiSpreadEDM − Local.PremiSpreadNB
Datain1.CARI11       = Local.NewRIComm        − Local.OldRIComm
Datain1.CARI15       = Local.NewBrokerageFree − Local.OldBrokerageFree
```

> ### **`SELISIH = (nilai_baru × ProrateEDMEnd) − nilai_lama_dengan_TreatyType_sama`**

⚠️ Sufiks **`NB`** pada `TsiSpreadNB`/`PremiSpreadNB` berarti **nilai lama**, bukan siklus New
Business. Penamaan menyesatkan; jangan diikuti di sistem baru.

### 4.2 Pencocokan pasangan

`[terverifikasi]` Loop `…4.3` beriterasi di `OfferFacIn.OldData.LocationList(...).SpreadingList`;
sub-langkah `…4.3.2` mencocokkan **`TreatyType`** baru terhadap lama sebelum mengisi
`Local.TsiSpreadNB` / `Local.PremiSpreadNB`.

⛔ **Pencocokan menurut `TreatyType`, bukan indeks posisi.** Penambahan atau penghapusan baris tidak
boleh menggeser pasangan.

### 4.3 Pemetaan parameter → kolom `facinproduction`

`[terverifikasi]` Tabel **82 kolom**; rule penulisnya `RDBList\InsertTreatyProduction_Sql` —
`pxObjClass` = **`Rule-Connect-SQL`** ⚠️ (bukan tipe RDBList meski di folder `RDBList\`).

| Kolom "menjadi" | Parameter | Kolom `*_SELISIH` | Parameter |
| --- | --- | --- | --- |
| `tsi_menjadi` | `Datain.CARI24` | `tsi_selisih` | `Datain.CARI26` |
| `premi_menjadi` | `Datain.CARI25` | `premi_selisih` | `Datain.CARI27` |
| `LOL_MENJADI` | `Datain.CARI50` | `LOL_SELISIH` | `Datain.CARI51` |
| `TSI100_MENJADI` | `Datain1.CARI48` | `TSI100_SELISIH` | `Datain1.CARI49` |
| `BROKERAGE_FEE_MENJADI` | `Datain1.CARI14` | `BROKERAGE_FEE_SELISIH` | `Datain1.CARI15` |
| `ricomm` *(tanpa sufiks)* | `Datain1.CARI10` | `ricomm_selisih` | `Datain1.CARI11` |
| `percent_ri_comm` *(ejaan beda)* | `Datain.CARI42` | `pct_ri_comm_selisih` | `Datain1.CARI9` |
| **`PCT_BROKERGARE_FEE`** *(ejaan asli)* | `Datain.CARI52` | `PCT_BROKERAGE_FEE_SELISIH` | `Datain.CARI53` |
| `prorate` | `Datain1.CARI8` | — | — |

⚠️ **5 kolom `*_MENJADI` vs 8 kolom `*_SELISIH`** — tiga pasangan memakai kolom "sesudah" **tanpa
sufiks**. Diport apa adanya.

⛔ Seluruh kolom nilai dibungkus `To_number(Replace({…},',','.'))` — **koma desimal**, kontrak K-027.

`[terverifikasi]` Keenam cabang lini mengisi **seluruh** parameter delta; bedanya cacah titik
penugasan: Fire 10× `CARI26`/`CARI27`, Aneka 7×, PA 6×, Golf/MC/MBU 5×.

### 4.4 Varian menyimpang

| Kondisi | Perlakuan |
| --- | --- |
| `Type=="4"` (Adjustment Spreading) | `…4.4–4.6` agregasi `TempSpread` per `TreatyType`; `…4.16` → **`CARI26 = CARI24`**, **`CARI27 = CARI25`**, `CARI51 = CARI50` — selisih **= nilai penuh** |
| `Type=="7"` + mata uang berubah | idem (`…4.17`) |
| `Type=="7"` + mata uang sama | `…4.18` → `CARI26 = TsiSpreadEDM`, `CARI27 = PremiSpreadEDM` |
| `Type=="3"` (Adj Rate) | `Local.PremiStart = @if(Type=="3", ProrateStart × .PremiumSpreaded, .PremiumSpreaded)`; `…4.8` → `PremiSpreadEDM = (PremiumSpreaded × CARI8 + PremiStart) − PremiSpreadNB` |
| Coverage hilang dari data baru | `…4.19`/`…4.20` → `PremiSpreadEDM = .PremiumSpreaded × −1` |
| `EdmType ∈ {1,2,3}` | `…4.12` → ketiga selisih persentase **dikali −1** |
| `Type!="4" && Type!="7"` | `…4.21` → `CARI26 = TsiSpreadEDM`, `CARI27 = PremiSpreadEDM` |

---

## 5. Kontrak masukan / keluaran modul

### 5.1 Masukan

| Dari | Isi |
| --- | --- |
| **Modul 1 (Seam 4)** — lapis A | dokumen polis lama utuh: `OldData.CurrencyList(...).SumTotalPayment`, `OldData...SpreadingList` beserta `TreatyType`, `OldData.PolicyData.{StartDateTime, EndDateTime}` |
| **Data kerja** | nilai baru setelah diubah pengguna: `Premium`, `Commision`, `BrokerageFee`, `Deduction2`, `PPh`, `PPN`, `TSISpreaded`, `PremiumSpreaded`, `TreatyType` |
| **Rasio prorata** | `ProrateEDMEnd`, `ProrateStartEDM` — ⚠️ sumbernya ambigu, lihat §8.1 |
| **Diskriminator** | `QuotationData.EdmType`, `QuotationData.Type`, `EdmTypeNew`, predikat lini bisnis (Seam 1), predikat jenis endorsement |

### 5.2 Keluaran

| Tingkat | Keluaran |
| --- | --- |
| **A — pembayaran** | per mata uang: `EDMNewPremi` · `EDMPremiMenjadi` · `EDMOldPayment` · `EDMOldPremi` · `EDMOldCommision` · `EDMOldPPN` · `EDMOldPPh` · `EDMOldBrokerage` · `EDMOldDeduction2` · **`SumTotalPayment`** |
| **B — baris produksi** | per baris spreading: `tsi_menjadi`/`tsi_selisih` · `premi_menjadi`/`premi_selisih` · `LOL` · `TSI100` · `BROKERAGE_FEE` · `ricomm` · `percent_ri_comm` · `PCT_BROKERGARE_FEE` · `prorate` |

### 5.3 Tipe

| Kelompok | Tipe |
| --- | --- |
| Seluruh nilai premi, TSI, komisi, brokerage, deduction, PPh, PPN, dan seluruh `EDMOld*` | **`Money{Amount decimal, Currency}`** (K-010/K-012) |
| `ProrateEDMEnd` · `ProrateStartEDM` · persentase komisi/brokerage | **`Ratio`**, skala per lini lewat resolver (K-018) |

⛔ `Money` + `Ratio` **gagal saat kompilasi**; jembatan tunggal **`Money × Ratio → Money`**.
⛔ **A.5 konsisten:** premi Nusantara Re memakai **angka** `0`, bukan string.

### 5.4 Seam — tidak ada yang baru

✅ **K-047:** modul ini **tidak mendapat seam baru**. Diuji lewat **Seam 3
(`services/premium.Calculate`)** yang **diperluas** menerima masukan endorsement (nilai lama dari
lapis A + kedua rasio prorata). Nilai lama datang dari **Seam 4** (modul 1).

**Total tetap 4 seam.** Setelah menelusuri `CountPremiEDM_DT`, saya **tidak menemukan alasan kuat**
mengusulkan seam ke-5: kedua tingkat delta sama-sama masuk/keluar lewat perhitungan premi, dan
memisahkannya akan mengunci pembagian internal yang belum tentu bertahan.

---

## 6. Kejanggalan K-046 — kasus uji bernama, bukan diperbaiki

⛔ Seluruhnya **diport apa adanya**. Nama test **wajib menyebut K-046** agar tidak "diluruskan" oleh
pembaca berikutnya.

| # | Kejanggalan | Bukti | Nama test yang disarankan |
| ---: | --- | --- | --- |
| 1 | `EDMPremiMenjadi` **tidak seragam** — 11 rumus semantik; tabel (lini × `EdmType` × jalur) | §3 | `K046_EDMPremiMenjadi_TidakSeragam_PerLiniDanEdmType` |
| 2 | Blok `EdmType==1` **hanya jalan untuk Life** — 11 dari 12 cabang perhitungan bergerbang `==2` | §3.1 | `K046_EdmType1_HanyaLife_LiniLainTidakDihitungUlang` |
| 3 | `CountDataEDMElse` langkah 2 berlabel *"MBU business"* tetapi bergerbang **`IsMarineCargo`** (C.9) — **abaikan label, ikuti kondisi** | §3.1 | `K046_LabelMBU_GerbangIsMarineCargo_IkutiKondisi` |
| 4 | **`PCT_BROKERGARE_FEE`** ejaan asli — kolom persen yang **benar** (B.6) | §4.3 | `K046_EjaanKolomPCT_BROKERGARE_FEE_Dipertahankan` |
| 5 | Ketidaksetangkupan `*_MENJADI`/`*_SELISIH` — 5 vs 8, tiga pakai kolom tanpa sufiks | §4.3 | `K046_PasanganMenjadiSelisih_TidakSetangkup` |
| 6 | **PPh dan PPN ditambahkan**, bukan dikurangkan, dalam netto | §2.2 | `K046_NettoPremi_PPhPPNDitambahkan` |
| 7 | Paradoks label *"batal sejak semula"* (= `EdmType==1` per K-029) bergerbang `==2` di 5 dari 6 lini | §3.1 | `K046_LabelBatalSejakSemula_Gerbang2_BukanSatu` |
| 8 | Sufiks **`NB`** pada variabel nilai lama berarti "lama", bukan New Business | §4.1 | — (catatan penamaan, bukan test) |

**A.5** tetap satu-satunya perbaikan sadar, diwarisi dari modul 1: premi Nusantara Re memakai angka
`0`. Catatan rekonsiliasi modul 1 berlaku penuh di sini.

---

## 7. Out of Scope

1. **Penulisan ke tabel produksi** — modul ini menghasilkan nilai kolom; menuliskannya adalah modul
   berikutnya, dan itu masih menunggu tabel flat + `ALL_SOURCE` dari DBA.
2. **Tangga akseptasi endorsement** — termasuk klaim arsip §4.3 *"nilai dasar akseptasi = selisih
   TSI"*, yang tetap **`[belum diuji]`**.
3. **Jalur fac out** (`FACOUTPRODUCTION`, `InsertTreatyProd_Sql`) — delta-nya berpola sama tetapi
   tabel dan kolomnya berbeda; modul tersendiri.
4. **`CountPremiEDMFacOut_DT`** — padanan fac out dari `CountPremiEDM_DT`, belum ditelusuri.
5. **Jalur produksi Life** (`SaveFacinLive_Act`, `SaveFacinSpreadLife_Sql`) — belum ditelusuri.
6. **Lapisan tampilan** — `Section\ViewOldDataPayment_IsUW` memanggil `CountPremiEDM_DT`; merancang
   layarnya terpisah.

---

## 8. ⚠️ Dua hal yang butuh keputusan work owner

### 8.1 `ProrateEDMEnd` dihitung di dua tempat dengan rumus berbeda — dan saling menimpa

`[terverifikasi]` Properti `pyWorkPage.OfferFacIn.ProrateEDMEnd` **ditulis oleh 5 rule (12
penugasan)** dan **dibaca oleh 13 rule**:

| Penulis | Cacah | Rumus |
| --- | ---: | --- |
| `Activity\SetValueToEDMWork` langkah 15 **(modul 1)** | 1 | `@Math.divide(Local.EdmToEnd, Local.TotalPeriod, 20)` — selisih **DateTime** periode polis lama |
| `DataTransform\CountPremiEDM_DT` langkah 9 | 1 | `Param.prorateedmend` = `@Math.divide(Param.edmdate, Param.startdate, 20)` — selisih **hari** via `@DateTimeDifference(…,"D")` |
| `CountPremiEDM_DT` langkah 11.3 | 1 | dipaksa **`1`** bila `EdmType==1` |
| `CountPremiEDM_DT` langkah 12.1 | 1 | **`Param.prorateedmstart`** bila `EdmType==2` ⚠️ rasio **awal** ditulis ke properti **akhir** |
| `Activity\CountPaymentEdm_Act` | 6 | — |
| `Activity\CountPaymentEdmTSIObj_Act` | 4 | — |
| `DataTransform\CountPremiEDMFacOut_DT` | 1 | — |

Pembacanya termasuk **keenam cabang produksi** `SaveFacinProd*EDM*_Act` (Fire 14×, lima lainnya 6×
masing-masing) dan `SetReinsurerEndorsement_Act` (15×).

⛔ **Artinya:** nilai yang disetel modul 1 **ditimpa** modul 2 sebelum dibaca jalur produksi. Selain
itu `CountPremiEDM_DT` menambah dua guard yang tidak ada di modul 1: pemaksaan `prorateedmend = 1`
bila `ProtectSpreading.CARI2=="FIX RATE"` (langkah 8), dan `Param.startdate == 0 → 1` (langkah 5).

**Pertanyaan:** apakah `ProrateEDMEnd` dimodelkan sebagai **satu properti yang ditimpa berurutan**
(port apa adanya, rapuh terhadap urutan) atau sebagai **beberapa rasio berbeda dengan nama berbeda**
(lebih jelas, tetapi menyimpang dari struktur lama)? Ini mengubah kontrak antar-modul 1 dan 2.

⚠️ Bila dipilih "port apa adanya", **spec modul 1 perlu catatan tambahan** bahwa nilainya bersifat
sementara.

### 8.2 Langkah 1 `CountPremiEDM_DT` seluruhnya `pyDisabled=true`

`[terverifikasi]` Langkah 1 beserta seluruh 11 sub-langkahnya ber-`pyDisabled=true` — implementasi
lama `EdmType==4` yang digantikan langkah 10. Isinya **berbeda** dari langkah 10 (mis. `EDMNewPremi`
= `SumTotalPayment × ProrateEDMEnd`, bukan netto).

**Pertanyaan:** `pyDisabled=true` pada DataTransform — apakah setara "langkah tidak dieksekusi"?
Bila ya, cabang ini **kode mati** dan cukup dicatat; bila tidak, ia ikut berjalan dan hasilnya
ditimpa langkah 10. Berbeda dari `pyStepsPreCondition=false` pada Activity yang sudah ditutup.

⚠️ Mengikuti **K-006**, saya **tidak** menyatakannya usang. Diajukan sebagai pertanyaan.

---

## 9. Kesiapan

✅ **Bahan lengkap untuk `/to-spec`:**

| Butir | Status |
| --- | --- |
| Dua tingkat delta dibedakan + keenam berkas terverifikasi | §1 |
| Rumus kanonik `SumTotalPayment` (19 penugasan) + 3 penyimpangan | §2.1 |
| Definisi netto `EDMNewPremi` + catatan PPh/PPN ditambahkan | §2.2 |
| Lima bentuk `EDMOldPayment`, 12 dari lapis A | §2.3 |
| **Tabel `EDMPremiMenjadi`** 4 kelompok — termasuk `CountPremiEDM_DT` yang belum ada di E-3 | §3 |
| Lima komponen "lama" tingkat pembayaran yang baru tercatat | §3.5 |
| Delta baris spreading + pencocokan `TreatyType` + 7 varian | §4 |
| Pemetaan parameter → 9 pasangan kolom | §4.3 |
| Kontrak masukan/keluaran + tipe + **Seam 3 diperluas, tanpa seam baru** | §5 |
| 7 kejanggalan K-046 dengan **nama test** | §6 |
| Out of scope | §7 |
| **2 pertanyaan untuk work owner** | §8 |

⛔ **`/to-spec` TIDAK dijalankan.** ⚠️ §8.1 sebaiknya diputuskan **sebelum** spec ditulis — ia
mengubah kontrak antar-modul.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
