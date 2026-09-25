# FINDING-003 — Baris Retensi Cedant hadir di `SpreadingRisk` saat rule tanpa penyaring membacanya

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

**Jenis**: laporan kondisi sistem lama
**Status**: KEHADIRAN **MUNGKIN**, BUKAN PASTI — bergantung urutan tindakan pengguna (lihat 3.3, koreksi atas versi pertama)
**Ruang lingkup bukti**: `D:\XML_NURE\Claim Non Prop` saja
**Tanggal**: 18 September 2026

---

## 1. Koreksi atas pernyataan saya sendiri

Ronde lalu saya menutup temuan ini dengan *"itu membuktikan tidak adanya penyaring, bukan adanya salah hitung"*. Batas itu benar pada saat itu, tetapi saya berhenti satu langkah terlalu awal. Urutan pemanggilan dapat ditelusuri seluruhnya dari XML dan dari dokumen struktur milik klien — tanpa Oracle sama sekali. Berikut hasilnya.

Sekaligus saya luruskan dua pernyataan saya yang tampak bertentangan: Ronde 2 saya menyebut tiga rule **menyaring** `TreatyName=="UR"`; Ronde 3 saya menyebut sembilan rule **tidak menyebut** `"UR"`. Keduanya benar dan tidak bertabrakan — sebagian menyaring, sebagian tidak. Yang belum saya kerjakan waktu itu adalah memeriksa apakah yang tidak menyaring memang melihat baris itu.

## 2. Kapan baris `"UR"` masuk ke PageList

`Activity\CountLossAllocation_act.xml` baris 7015–7120:

```
pyWorkPage.ClaimData.SpreadingRisk(<APPEND>).TreatyType      = "UR"
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).TreatyName        = "UR"
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).ClaimEstimation   = Local.UR
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).ClaimAmountAdjust = Local.UR
```

Retensi Cedant bukan entitas terpisah — ia baris tambahan di PageList yang sama dengan Layer.

## 3. Urutan pemanggilan — DIBUKTIKAN ULANG DARI XML

**Koreksi.** Versi pertama bagian ini memakai `Struktur_Flow_TreatyIn.xlsx` sebagai bukti. Berkas itu turunan, bukan sumber, dan tidak boleh berdiri sebagai bukti. Seluruh bagian ini ditulis ulang dari XML. **Hasilnya berbeda dari versi pertama, dan versi pertama salah pada satu hal pokok** — lihat sub-bagian 3.3.

### 3.1 Siapa memanggil apa (dari `pyStepsActivityName` di XML)

Indeks panggil activity ke activity, dibaca dari `<pyStepsActivityName>Call ...</pyStepsActivityName>` beserta `<pyStepPageReference>`:

| Dipanggil | Oleh | Langkah |
|---|---|---|
| `CountLossAllocation_act` | `CountClaimTNP_Act` | 15 |
| `CountLossAllocation_act` | `AddAkseptasiCNP_Act` | 18 |
| `CountLossAllocation_act` | `InputAkseptasi_PreAct` | 5 |
| `GenerateCFS_act` | `SaveToOS` | 7 |
| `CountTotalInsterest_Act` | `SetCurencyInterest_act` | 6 |
| `CountSpreadingXOL` | `AdjClaimCNP_Act` | 16 |
| `CopyOldataCurr_act` | `CountTotalInsterest_Act` | 11 |
| `SendEmailKlaim` | `CreateChildKomiteCNP_Act` | 35 |

### 3.2 Temuan pokok: urutannya tidak ditentukan kode, melainkan layar

**Tidak ada satu pun activity di folder ini yang memanggil `CountClaimTNP_Act`.** Rule itu dipicu dari **Section**, lewat `<pyActivity>`, sebanyak 42 rujukan di tiga berkas:

| Section | Activity yang dipicu, dalam urutan dokumen |
|---|---|
| `Section/AdjustmentDetailNP.xml` | `CountClaimTNP_Act`, `CountLossAllocation_act`, `AdjClaimCNP_Act` |
| `Section/InputAcceptation.xml` | `SetCurencyInterest_act`, `CountTotalInsterest_Act`, `CountClaimTNP_Act`, `CountLossAllocation_act`, `SaveToOS`, `AddAkseptasiCNP_Act`, `DeleteAkseptasi_Act` |
| `Section/OutstandingClaim(1).xml` | `CountClaimTNP_Act`, `CountLossAllocation_act`, `GenerateCFS_act` |
| `Section/ViewDetailDeptHeadTreatyIn_UW.xml` | `CountSpreading_Act` |

Ini **kontrol UI yang berdiri sendiri-sendiri pada satu layar**, bukan rantai berurutan. Urutan dokumen di dalam XML Section adalah urutan tata letak, **bukan** urutan eksekusi.

Konsekuensinya, dan inilah temuan sebenarnya:

> **Urutan jalannya perhitungan alokasi ditentukan oleh kontrol mana yang ditekan pengguna, bukan oleh rantai pemanggilan di dalam kode.** Tidak ada satu tempat pun di XML yang menetapkan bahwa `GenerateCFS_act` berjalan sesudah `CountLossAllocation_act`.

### 3.3 Yang harus saya cabut

Versi pertama menyatakan: *"baris Retensi Cedant sudah ada di dalam `SpreadingRisk` ketika ketiga rule itu membacanya - ini sekarang fakta, bukan dugaan."*

**Pernyataan itu saya cabut.** Ia bersandar pada penomoran outline spreadsheet. Dari XML, yang dapat dinyatakan hanyalah:

- `CountLossAllocation_act` (yang meng-`<APPEND>` baris `"UR"`) dan `GenerateCFS_act` dapat dipicu dari **layar yang sama** (`OutstandingClaim(1)`), begitu pula `CountTotalInsterest_Act` dan `DeleteAkseptasi_Act` pada `InputAcceptation`.
- Karena itu baris `"UR"` **dapat** sudah ada saat rule-rule itu membaca `SpreadingRisk`, tergantung urutan penekanan kontrol oleh pengguna.
- Apakah ia **selalu** ada: **TIDAK DAPAT DISIMPULKAN DARI XML.** Itu bergantung perilaku pengguna, dan hanya terbaca dari data produksi.

Yang justru menguat dari koreksi ini: ketergantungan pada urutan klik itu sendiri adalah cacat rancangan, terlepas dari apakah salah hitungnya pernah terjadi.

## 4. Akibatnya berbeda per rule — dan satu di antaranya ternyata aman

Kehadiran baris itu tidak otomatis berarti salah hitung. Saya periksa satu per satu:

### 4.1 `CountTotalInsterest_Act` — **TIDAK ADA SALAH HITUNG**

```
Local.SizeLoss = @SizeOfPropertyList(pyWorkPage.ClaimData.SpreadingRisk)
Local.SizeEst  = @SizeOfPropertyList(pyWorkPage.ClaimData.SpreadingRisk)
```
Cacahan itu memang ikut menghitung baris `"UR"`. Tetapi pemakaiannya hanya sebagai penjaga kosong-tidaknya:
```
Local.SizeLoss > 0 && Local.PropCount == 1
Local.SizeEst  > 0 && Local.PropCount == 1
```
Satu baris tambahan tidak mengubah hasil `>0`. **Rule ini bersih.** Saya catat ini karena dugaan awal saya mengarah ke sebaliknya.

### 4.2 `GenerateCFS_act` — **TUNTAS: baris UR ikut tersalin, bila ia ada**

Rantai `Local.XOL` / `Local.Check` selesai ditelusuri.

**Bentuk rantainya** — berulang empat kali di dalam rule, dengan pola yang sama persis:

```
Local.Check = 0                        <- penanda direset
Local.XOL   = .TreatyName              <- disalin dari BARIS YANG SEDANG DIBACA
Local.Check = 1                        <- penanda dinyalakan
... precondition: Local.Check == 0     <- penjaga agar satu TreatyName diproses sekali
... precondition: Local.XOL == .TreatyName  <- pencocokan baris berikutnya
```

Empat penetapan `Local.XOL`, empat belas penetapan `Local.Check`, sebelas precondition yang membacanya.

**Yang menentukan**: `Local.XOL = .TreatyName` menyalin nama treaty **dari baris yang sedang dibaca, apa adanya**. Tidak ada penyaring. Dan sapuan sebelumnya sudah memastikan `GenerateCFS_act` **tidak menyebut `"UR"` satu kali pun** di seluruh 66 rujukan `SpreadingRisk`-nya.

**Kesimpulan**: bila baris `TreatyName = "UR"` ada di `SpreadingRisk` saat rule ini berjalan, maka `Local.XOL` akan bernilai `"UR"`, barisnya lolos seluruh precondition, dan ia **tersalin ke `TempDataOutStanding.ClaimData.SpreadingRisk(<APPEND>)` beserta `RetroList`-nya** — yaitu ke struktur retrosesi, diperlakukan persis seperti Layer.

**Syaratnya tetap satu, dan itu tidak berubah**: baris `"UR"` harus sudah ada saat rule ini dipicu. Karena urutan pemicuan ditentukan kontrol yang ditekan pengguna (bagian 3.2), kehadirannya **mungkin, bukan pasti** — persis batas yang ditetapkan di bagian 3.3.

Jadi yang tadinya *"belum dapat dipastikan apakah tersalin"* sekarang menjadi *"pasti tersalin bila ada"*. Yang tersisa hanya pertanyaan kehadiran, dan itu pertanyaan data — bukan lagi pertanyaan kode.

### 4.3 Enam rule sisanya — **selesai, dan hasilnya menggugurkan sebagian temuan ini**

`CountSpreading_Act`, `CountSpreadingXOL`, `DeleteAkseptasi_Act`, `SetAccoutNo_Act`, `SendEmailKlaim`, `CountClaimTNP_Act`: **nol rujukan ke PageList `SpreadingRisk`.** Seluruh kemunculan teks `SpreadingRisk` pada keenamnya adalah nama class di metadata langkah, bukan pembacaan daftar.

`CopyOldataCurr_act` punya dua rujukan nyata: `@LengthOfPageList(...)>0` sebagai penjaga, dan satu `Property-Set` `repeat=EMBEDDED` yang menyetel `.Currency = Local.CurrencyNew` pada setiap baris. Baris `"UR"` ikut terkena, tetapi baris itu memang membawa `Currency` sendiri (disetel `CountLossAllocation_act`), jadi memperbaruinya bersama baris lain konsisten — **bukan cacat**.

**Ringkas**: dari sembilan rule yang semula saya sebut terdampak, **satu terbukti terdampak** (`GenerateCFS_act`), **dua terbukti tidak**, dan **enam tidak membaca daftarnya sama sekali**. Angka "16–66 rujukan" pada versi pertama keliru — lihat koreksi di `BLUEPRINT.md` §2.4.

## 5. Yang menutup temuan ini

Seluruhnya pekerjaan XML. **Tidak satu pun butuh Oracle.**

1. Telusuri rantai `Local.XOL` / `Local.Check` di `GenerateCFS_act`.
2. Periksa enam rule sisanya dengan cara yang sama: apakah baris `"UR"` ikut terbaca, dan apakah pembacaan itu mengubah angka.
3. Tetapkan mana yang memang **seharusnya** menyertakan Retensi Cedant — sebagian mungkin benar demikian. Itu pertanyaan maksud bisnis, bukan pertanyaan kode → Round 3.
