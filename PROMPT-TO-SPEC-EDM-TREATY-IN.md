# PROMPT TO-SPEC — EDM Treaty In

> Salin seluruh isi berkas ini sebagai prompt ke sesi eksekutor.
>
> ✅ **Prasyarat sudah terpenuhi.** `.scratch\edm-treaty-in\KEADAAN-EDM-TREATY-IN.md` ada, dan
> ronde yang menulisnya **sekaligus menjadi ronde verifikasinya** — seluruh angka diuji dua cara di
> dalam berkas itu sendiri. **Tidak ada berkas `VERIFIKASI-*` terpisah untuk EDM, dan memang tidak
> diperlukan.**
>
> ⚠️ **Skill `to-spec` tidak dapat dipanggil sendiri oleh agen** (CLAUDE.md §8). Ketikkan
> `/to-spec` sebagai manusia, lalu berkas ini menjadi briefnya. Bila skill tidak dipakai, tulis
> spec langsung dalam bentuk rumah yang dijelaskan di Bab 3.

---

## 0. LINGKUP — DIKUNCI

Hanya modul **`D:\XML\RNM_BRD\EDM Treaty In`**. `NB Treaty In` **hanya sebagai pembanding**.

`D:\XML\RNM_BRD\` adalah korpus **READ-ONLY**. Menulis hanya ke `OUTPUT_HASIL_RNM\`.
`D:\XML\nusantara-re\` **terlarang** — jangan dibaca, dikutip, dibandingkan, atau dijadikan sasaran.

Keluaran satu berkas: **`OUTPUT_HASIL_RNM\.scratch\edm-treaty-in\spec.md`**.

⛔ `grilling-ronde-1.md` tersegel — jangan disunting. Koreksinya sudah hidup di berkas keadaan Bab 6.

---

## 1. SUMBER DAN URUTAN KEWENANGAN

Bila dua sumber bertentangan, yang di atas menang. **Tuliskan pertentangan yang Anda temukan.**

| # | Sumber | Sifat |
| ---: | --- | --- |
| 1 | `.scratch\edm-treaty-in\PERTANYAAN-RONDE-1.md` — **P50–P60**, seluruhnya terjawab | keputusan work owner, mengikat |
| 2 | `.scratch\nb-treaty-in\PERTANYAAN-untuk-*.md` — enam lembar, **P1–P49** | keputusan work owner, mengikat; berlaku lintas modul |
| 3 | `.scratch\edm-treaty-in\KEADAAN-EDM-TREATY-IN.md` | keadaan terukur, sudah diuji dua cara |
| 4 | `.scratch\nb-treaty-in\KEADAAN-NB-TREATY-IN.md` Bab 4 | dua belas ketetapan; berlaku **kecuali terbukti sebaliknya** |
| 5 | `.scratch\edm-treaty-in\grilling-ronde-1.md` | latar; ⛔ **empat pernyataannya sudah terbukti keliru** — lihat keadaan Bab 6 |
| 6 | korpus XML | selalu boleh dipakai untuk membuktikan ulang |

⚠️ **Jangan menyalin angka dari nomor 5.** Empat pernyataan ronde 1 dikoreksi, dan dua angkanya
*(378 langkah, 1.062 pasangan)* **ditarik dari peredaran karena tidak dapat direproduksi**.

Contoh bentuk rumah: sebelas `spec.md` di `.scratch\*\`, ditambah `.scratch\nb-treaty-in\spec.md`
*(60 KB, 48 user story, 96 AC)* yang paling dekat kerabatnya.

---

## 2. KEADAAN YANG DIWARISI

Ringkas, supaya tidak perlu ditemukan ulang. Rinciannya di sumber nomor 3.

| Hal | Angka |
| --- | ---: |
| Berkas `.xml` dalam modul | **163** *(+1 `.xlsx`, bukan aturan)* |
| `Activity` terjangkau / yatim | ⭐ **66 / 0** |
| Langkah `Property-Set` terjangkau / yatim | ⭐ **380 / 0** |
| Langkah yang membawa isinya | **379** *(99,7 %)* |
| Pasangan nama=nilai terisi | **1.309** |
| `RDBList` punya naskah SQL | **36 dari 36** *(100 %)* |
| Pertanyaan EDM terjawab | **11 dari 11** |
| Penahan tersisa | ⭐ **nol** |

### ⭐⭐ Yang paling penting dinyatakan di spec

**Lingkup EDM tidak menyusut.** NB kehilangan 28 aturan dan 561 langkah karena foldernya adalah
*dependency closure*; folder EDM **bukan**. Setiap aturan EDM terjangkau dari titik masuk nyata.

⇒ **380 langkah dikerjakan seluruhnya.** Jangan ada yang mengharapkan penyusutan 60 % seperti NB.

### Tiga hal yang tidak dimiliki NB

| Hal | Bukti |
| --- | --- |
| tipe rule **`ConnectREST`** | `convertJsonNusareToProduction`, kelas `ASM-FW-GISFW-Work` |
| tipe rule **`SystemSettings`** | `LinkService` |
| ⭐ **kelas kerja sendiri** | `Flow\InputAddendumTreatyIn`, kelas `…Work-EndorsementTreaty` |

⭐ **Endorsemen adalah jenis kasus tersendiri**, bukan tahap di dalam kasus polis baru. Alurnya
memuat 8 kotak keputusan, 4 penugasan, 22 sambungan, dan menyebut **tiga antrean** —
`ReasTreatyInAdmin` · `ReasTreatyInSecHead` · `ReasTreatyInDeptHead` — sejalan **P13**.

⚠️ **Kedalaman sarang langkah mencapai enam tingkat.** Pembaca yang hanya menelusuri satu tingkat
kehilangan 215 dari 380 langkah.

---

## 3. BENTUK SPEC

Bab wajib, berurutan:

```
Cara membaca berkas ini      Problem Statement      Solution
User Stories                 Implementation Decisions
Testing Decisions            Acceptance Criteria
Out of Scope                 Butir [terbuka] — daftar penuh
Further Notes                Lampiran
```

Jangan mengejar ukuran — kejar kelengkapan. Tetapi ingat lingkupnya **penuh**, tidak menyusut.

### Blok ringkasan di "Cara membaca berkas ini"

Wajib memuat cacah: user story · acceptance criteria · sebaran penanda · butir `[terbuka]` aktif ·
butir `[penyimpangan sadar]`.

**Dihitung dua cara yang berbeda**, dan sebutkan jendela hitungnya. Bila Anda menulis angka sebelum
mengukurnya lalu memperbaikinya — **kutip angka lamanya, jangan dihapus** (CLAUDE.md §4a).

### Aturan menulis Acceptance Criteria

Bab itu **tidak memutuskan apa pun**. Ia menyatakan ulang keputusan yang sudah ada, dalam bentuk
yang dapat diuji dari luar. Bila sebuah butir terasa seperti keputusan baru, ia salah tulis.

```
N. `[terverifikasi]` <pernyataan>. Test yang menemukan <keadaan sebaliknya> **gagal**. *(Bab X)*
```

Setiap butir membawa **penanda** dan **rujukan bab**. Butir tanpa penanda adalah cacat.

---

## 4. YANG WAJIB ADA DI ACCEPTANCE CRITERIA

### 4.1 Sebelas ketetapan EDM sendiri — P50 sampai P60

| Butir | Ketetapan yang harus dapat diuji |
| --- | --- |
| **P50** | EDM tidak punya FacOut. Pengiriman kedua **tidak pernah berjalan** dan **tidak dibangun** |
| **P51** | Data produksi dihapus bila konversi gagal — ⭐ **bersyarat**, hanya bila `STS_KONVERSI` belum `1`. Penghapusan **hanya lewat procedure**, dan pengguna **diberi tahu** lebih dulu |
| **P52** | Langkah 10 `Page-Set-Messages` **dihidupkan**; langkah 18 `ASMForceCaseClose` **tetap mati** |
| **P53** | Pengiriman ke produksi **tanpa autentikasi dari sisi kita** — jalur internal; autentikasi milik penerima. Alamat tujuan **variabel lingkungan, bukan harfiah** |
| **P54** | `TANGGAL_CLOSING` dibaca `WHERE ROWNUM = 1` — **satu baris berlaku global**; perubahan tanggal mengubah seluruhnya |
| **P55** | Penomoran **satu jalur**: `NoPolis + "/E" + ProdKe` dengan `ProdKe` dibubuhi nol di depan bila `< 10`. Adendum dan premi tambahan **ikut jalur yang sama** |
| **P56** | Pembatalan adalah **jenis endorsemen**, dipilih di awal — bukan cabang tersendiri |
| **P57** | Endorsemen **berlapis**; selisih dihitung terhadap keadaan **tepat sebelumnya**, bukan terhadap polis asli |
| **P58** | Data lama **disalin dan beku** saat berkas dibuat *(`CreateEDMT` langkah 10, sumbernya dibuang langkah 13)* |
| **P59** | Lingkungan uji dan produksi memakai **basis data berbeda** |
| **P60** | `CountSpreading_Act` EDM: **8 langkah** lawan 7 di NB; langkah tambahan `.SpreadingRiskList(1).SplitRNMSharePct = 100`; presisi pembagian **20** lawan 10 di NB |

### 4.2 Dua belas ketetapan NB — berlaku juga di sini

Seluruhnya dari `KEADAAN-NB-TREATY-IN.md` Bab 4. ⚠️ **Satu di antaranya dilanggar sistem lama** —
lihat Bab 5.2 berkas keadaan EDM: hanya 16 dari 31 pernyataan `FROM` memakai skema `POOLDATA.`,
`JSON_POLIS` muncul dua ejaan, dan `DATAPEGA` adalah skema kedua. Ketetapannya **tetap berlaku
untuk sistem baru**; yang berubah hanya catatan bahwa sistem lama tidak menaatinya.

### 4.3 ⭐ Penomoran EDM — wajib satu AC tersendiri

```
PolicyTreatyIn.EDMNo = OldData.PolicyNo + "/E"
                     + @If(ProdKe < 10, "0" + ProdKe, ProdKe)
```

Ditetapkan `Activity\SetEDMTNoPolis`. Bentuknya: nomor polis diikuti `/E` lalu dua digit.

⛔ **`GenerateNoEDMTreaty` adalah aturan mati.** Ia dijaga syarat `EDMNo` kosong yang **tidak
pernah benar**, sebab `SetEDMTNoPolis` sudah mengisinya lebih dulu. **Tidak dimigrasi.** Masuk
Out of Scope.

### 4.4 ⛔⛔ Bentuk dokumen tidak seragam — wajib beberapa AC

**Penentunya** `QuotationData.ProportionalType` *(`Proportional` / `NonProportional`)* dan
`IsNewPolicyNonProp` *(`0` / `1`)*. Korpus menyebutnya 170 dan 87 kali.

⛔ **Yang paling berbahaya:** `ListInstallment` **bersarang** pada satu bentuk, **datar** pada
bentuk lain. Keduanya nyata di korpus — 58 rujukan `ListInstallment(n).InstallmentList` di 6 berkas,
101 rujukan `ListInstallment(n).<medan skalar>` di 11 berkas.

⇒ **Pengurai yang menganggap bentuknya seragam akan patah.** Wajib satu AC yang mengujinya, dan
**test yang hanya mencakup satu bentuk dinyatakan tidak memadai**.

Tiga contoh nyata memberi 37, 64, dan 47 medan skalar tingkat atas; **23 medan ada di ketiganya**,
gabungannya **74**.

---

## 5. ENAM STORED PROCEDURE — NASKAHNYA KINI LENGKAP

⭐ **Butir `[terbuka]` nomor 1 pada keadaan Bab 9 sudah TERTUTUP.** Work owner menyerahkan kedua
naskah yang kurang pada 22 September 2026. Tercatat di jawaban **P1**.

| # | Procedure | `COMMIT` di dalam |
| ---: | --- | :---: |
| 1 | `PEGA_TREATY_IN` — menulis **dua** tabel | ya |
| 2 | `PEGA_JSON_POLIS_TREATYIN` | ya |
| 3 | `PROC_GENERATE_SEQUENCE_NUMBER` | **tidak** |
| 4 | `PEGA_DELETE_ERROR_KONVERSI` | ya |
| 5 | `INSERTJSONPOLISMONITORING` | ya |
| 6 | `InsertUpdateAchievment` | **tidak** |

### Empat hal dari kedua naskah baru yang wajib masuk spec

1. ⛔ **`INSERTJSONPOLISMONITORING` bukan "insert-or-update".** Bila baris ber-`IDPEGA` sama sudah
   ada, ia **berhenti tanpa berbuat apa-apa**. Namanya menjanjikan pembaruan; perilakunya tidak.
   `PRODKE` **selalu diisi `0`**, bukan dari parameter.
2. ⛔ **Pesan galatnya memuat markah HTML** dan membaca penanda lingkungan `p_ProductionLevel`.
   Di sistem baru pesan galat **tidak membawa markah**; penyajian milik lapisan tampilan.
   Menyentuh **ADR-0005**.
3. ⭐ **`InsertUpdateAchievment` hanya menyisip** — tidak ada jalur pembaruan. Parameter uangnya
   bertipe **`NUMBER`**, satu-satunya procedure yang menerima uang sebagai bilangan. Di sinilah
   presisi ditetapkan saat melintasi batas.
4. ⛔⛔ **Penjaga anti-duplikatnya dipatahkan galat uang.** Pemeriksaannya membandingkan **seluruh
   17 kolom termasuk enam kolom uang**. Selisih `2,76 x 10^-7` — yang **terbukti ada di data
   produksi** — membuat baris yang seharusnya sama dinilai berbeda dan tersisip dua kali.
   ⇒ **Sistem baru tidak meniru pemeriksaan duplikat berbasis nilai uang.** Kuncinya medan
   pengenal. ⛔ **Butir `[terbuka]`, pemutusnya work owner** — jangan diputuskan di spec.

---

## 6. YANG WAJIB MASUK OUT OF SCOPE

⚠️ **Tidak ada aturan yatim di EDM.** Out of Scope di sini berisi hal yang **mati**, **tidak
pernah berjalan**, atau **sengaja tidak ditiru** — bukan aturan yang tak terjangkau.

| Yang dikeluarkan | Sebab | Butir |
| --- | --- | --- |
| jalur konversi **FacOut** | EDM tidak punya FacOut; tidak pernah berjalan | **P50** |
| langkah 18 `Call ASMForceCaseClose` | tetap mati, keputusan work owner | **P52** |
| `Activity\GenerateNoEDMTreaty` | aturan mati; syaratnya tak pernah benar | **P55** |
| `OldData` **di dalam** `OldData` | dibuat `CreateEDMT` langkah 10, **nol rujukan** di korpus | keadaan Bab 9 |
| pencarian berdasarkan nomor urut antrean | diganti pemeriksaan keanggotaan | **P25** |
| pemeriksaan duplikat berbasis nilai uang | dipatahkan galat presisi | Bab 5 di atas |
| markah HTML di dalam pesan galat basis data | penyajian milik lapisan tampilan | Bab 5 di atas |
| `Struktur_InputAddendumTreatyIn.xlsx` | lembar kerja, bukan aturan Pega | keadaan Bab 2.1 |

Setiap butir menyebut **sebabnya** dan **butir asalnya**. Jangan ada yang dikeluarkan tanpa alasan
tertulis.

### ⛔ Bab karantina galat uang — tulis tersendiri

EDM **menghitung selisih**, dan selisih dua angka bergalat membawa **jumlah** galatnya. Endorsemen
berlapis *(P57)* menumpuknya lagi. Bab ini memuat tiga hal:

1. apa yang **sudah** diketahui — galatnya lahir di rantai perhitungan, bukan di penyimpanan;
2. apa yang **menunggu keputusan work owner** — tiga pilihan pada **P29**: ikuti apa adanya ·
   bulatkan saat migrasi · bulatkan saat dihitung;
3. apa **akibatnya bila tidak pernah diputuskan** — termasuk penjaga duplikat yang gagal.

⭐ **ADR-0003 tetap ditegakkan apa pun pilihannya** — sistem baru tidak memakai `float`, sehingga
tidak menambah galat baru. ⛔ **Tidak ada ADR baru.**

---

## 7. DISIPLIN

Setiap pernyataan membawa bukti: **jalur berkas + tipe rule + nama rule**.

Penanda: `[terverifikasi]` · `[dugaan]` · `[terbuka]` · `[data DBA]` · `[keputusan work owner]` ·
`[penyimpangan sadar]`.

**Jangan menutup pertanyaan terbuka sendiri.** Penutupan milik work owner, DBA, Product &
Underwriting, Aktuaria, Finance, atau IAM.

**Jangan menyalin nilai berupa nama orang.** Nomor polis **tidak** ditulis apa adanya. Contoh JSON
**tidak disimpan apa adanya** — ketiganya memuat nama orang. Nol rahasia, nol token, nol data
nasabah. Alamat dan hos internal adalah **variabel lingkungan**, tidak pernah harfiah.

Uang tidak pernah `float` (CLAUDE.md §7). Arah ketergantungan `handlers → services → repository` (§5).

### Enam jebakan yang sudah menjerat proyek ini

1. **Menebak nama tag.** ⭐ `Activity` memakai `PropertiesName`/`PropertiesValue` **tanpa** awalan
   `py`; `Flow` memakai **dengan** awalan. Kebalikan satu sama lain. `pyPropRef`,
   `pyPropertiesValue`, `pyStepsPage`, `pyCriteriaValue`, `pyResult`, `pyShapeName` **tidak ada**.
2. `<rowdata REPEATINGINDEX="n"/>` yang menutup sendiri tidak tertangkap pola
   `<rowdata...>(.*?)</rowdata>`.
3. `pyRowNum` pada `pyOrConditions` berbasis **nol**; baris tabel berbasis satu.
4. Memasangkan dua daftar menurut **urutan**, bukan menurut blok `rowdata` yang sama.
5. Membaca `pyConditionString` — keterangan manusia — alih-alih bentuk yang dijalankan *(P23)*.
6. Menyatakan sesuatu nihil tanpa menyebut lingkup penelusuran.

⚠️ **Buang blok `pyExpressionGadget` lebih dulu** — peninggalan layar penyunting.
⚠️ **Jangan `html.unescape` sebelum mencocokkan pola** — entitas `&lt;` berubah menjadi tanda
kurung sudut dan memecah pola `<PropertiesValue>…</PropertiesValue>`. Jebakan ini sudah memakan
satu sesi penuh.

Jalur korpus memuat spasi — pakai Python, bukan loop shell.

---

## 8. BAB WAJIB — TELEMETRI EKSEKUSI

`spec.md` ditutup dengan bab `## TELEMETRI EKSEKUSI`.

Angka token sejati tidak terlihat dari dalam sesi. Ukur dari luar:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Catat token keluaran · cache-read · jumlah panggilan alat · durasi · biaya · byte dibaca.
Bila pengukuran dari luar tidak dilakukan, **katakan begitu** — jangan menaksir lalu menyajikannya
sebagai angka terukur.

---

## 9. YANG MENANDAKAN RONDE INI BERHASIL

- `spec.md` lengkap sebelas bab, seluruh AC berpenanda dan berujuk bab.
- **P50–P60 seluruhnya terwakili** di Acceptance Criteria, satu butir atau lebih masing-masing.
- Dua belas ketetapan NB diperiksa keberlakuannya, dan **pelanggaran skema SQL dinyatakan**, bukan
  didiamkan.
- **Dua bentuk `ListInstallment` diuji terpisah** — test satu bentuk saja dinyatakan tidak memadai.
- Penomoran endorsemen punya AC sendiri, dan `GenerateNoEDMTreaty` dinyatakan mati di Out of Scope.
- Keenam stored procedure tercatat, dengan **empat temuan dari dua naskah baru** masuk spec.
- Galat uang dikarantina dalam babnya sendiri — **tidak diputuskan, tidak dihilangkan**.
- **Lingkup penuh dinyatakan terang-terangan** — nol yatim, 380 langkah dikerjakan seluruhnya.
- Blok ringkasan dihitung dua cara, dengan jendela hitung disebutkan.
- Nol butir ditutup sendiri, nol ADR baru, nol nama orang tersalin, nol contoh JSON tersimpan.
- Bab telemetri terisi angka terukur, bukan taksiran.

---

*Disusun 22 September 2026, sesudah ronde keadaan EDM mengoreksi empat pernyataan ronde 1, menarik
dua angka yang tidak dapat direproduksi, dan menutup penahan terakhir — naskah dua stored procedure
yang tidak pernah diminta.*
