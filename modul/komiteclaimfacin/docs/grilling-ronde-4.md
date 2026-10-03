# Komite Claim Fac In — Grilling Ronde 4
## Empat rule terakhir, mekanisme yang akhirnya terbukti, dan satu vonis ronde 3 yang saya batalkan sendiri

⭐ **Ronde ini MEMBACA.** Yang diputuskan hanyalah **empat keputusan work owner** di §6, dan
keempatnya **sudah diambil** sebelum ronde ini mulai.

**Jendela sensus seluruh ronde ini:** 114 berkas XML `Komite Claim FacIn`, dibandingkan bila perlu
dengan 482 berkas `Claim Fac In` dan 80 berkas `Komite Claim Prop`. ⛔ Tidak ada berkas di luar
ketiga korpus itu yang dibuka.

> **SENSUS BERKAS INI**
>
> ⛔ **Jendelanya disebut lebih dulu: seluruh berkas ini DIKURANGI blok sensus ini sendiri**
> — **791 baris** menurut cacah baris-baru, **792** menurut pemisahan teks. ⚠️ **Kedua cara
> BERBEDA SATU, dan sebabnya diketahui:** berkas berakhir dengan baris baru, sehingga pemisahan
> teks menghasilkan satu baris kosong penutup yang bukan baris sungguhan. ⭐ **Angka yang dipakai:
> 791.** ⭐ Blok sensus dikecualikan dengan sengaja: bila ia menghitung dirinya sendiri, angkanya
> berkejaran tanpa henti setiap kali diperbaiki. ⭐ **Cara yang sama dipakai ronde 3.**
>
> `[terverifikasi]` **31** · `[keputusan work owner]` **5** · `[terbuka]` **7** ·
> `[penyimpangan sadar]` **1** · `[data DBA]` **5**
> ⚠️ **62** · ⭐ **210** · ⛔ **123** · ✅ **17** · ❓ **6** · RALAT **6** · tabel **25**
>
> ⭐ **Dihitung DUA CARA:** *(a)* cacah penanda atas badan berkas sebagai satu teks utuh;
> *(b)* cacah ulang **baris demi baris**. ✅ **Keduanya sepakat pada kedua belas angka.**
>
> ⚠️ **Catatan yang jujur:** angka `[terbuka]` **7** menghitung **penanda yang tertulis**, bukan
> **butir dalam register**. ⛔ Register §7.3 mencatat **14 butir terbuka** — keduanya benar,
> jendelanya berbeda.

---

## §1 — ⛔⛔⛔ `SaveAcceptation_KMT`, dan sebuah vonis ronde 3 yang **BATAL**

### 1a — Rule itu bukan yang namanya janjikan

`[terverifikasi]` `Activity\SaveAcceptation_KMT` — kelas **`ASM-FW-GCNMFW-Data-Adjustment`**,
ruleset **GCNMFW 01-01-24**, sistem pengekspor **`pegadevnusare2`**.

| | |
| --- | ---: |
| langkah | **6** |
| ber-remark | **0** |
| bendera: kosong / `true` | 3 / 3 |
| parameter rule | ⛔ **NOL** |
| metode berat *(`Obj-*`, `RDB-*`, `Connect-*`)* | ⛔ **NOL** |
| gerbang berwewenang | ⛔ **NOL** |

**Isi keenam langkahnya:**

| Lgk | Gerbang | Yang dikerjakan |
| --- | --- | --- |
| 1 | — | membaca indeks objek/item/adjustment · `Data.CARI12 = "Acceptation"` · `Local.NoAksep = .AcceptedNo` · `Local.Retro = 0` |
| 2.1 | `.TreatyType=="10015"` | ⭐ `Local.Retro = 1` |
| 3 | `Local.Retro=="1"` | ⭐ `.IsFacRetro = 1` + dua cerminan ke `TempOpenPage.ClaimData.ObjectList(…)` |
| 4 | `Local.Retro=="1"` **kode 3 di sisi benar** | `.IsFacRetro = 0` + dua cerminan |
| 5 | — | satu cerminan bernilai `0` |

⭐⭐ **Namanya menyesatkan: ia BUKAN penyimpan. Ia PENENTU BENDERA RETRO.**
⭐ Satu-satunya keputusan yang dibuatnya: **apakah adjustment ini bersifat Fac Retro**, ditentukan
oleh **`TreatyType == "10015"`**.

⚠️ **Dan bendera itu punya akibat besar** `[terverifikasi]`: `ApprovalKomite_Act` lgk 2.2
bergerbang `.IsFacRetro==1` lalu lgk 3 **`Exit-Activity`** bergerbang `Local.Retro==1`.
⭐ **Artinya klaim Fac Retro MELEWATI penyiapan wewenang penyetujuan sepenuhnya.**
⚠️ `KomitePost_Adjustment` lgk 17 dan 19 juga bergerbang `Local.FacRetro==1` **berkode arah 3**.

⭐ **`[terbuka]` BARU:** apa arti `TreatyType == "10015"`, dan kenapa Fac Retro boleh melompati
penyiapan wewenang. ⛔ Angka itu **ditanam keras**, sama seperti pita nilai di §6 K-lama.

### 1b — ⛔⛔ RALAT BESAR: penerimaan **TERSIMPAN**, ronde 3 salah

⛔ **Ronde 3 §2b menyimpulkan:** *"⛔ **Jadi PENOLAKAN tersimpan permanen; PENERIMAAN tidak.**"*
⛔ **Kesimpulan itu SALAH, dan saya batalkan di sini.**

`[terverifikasi]` Jendela: **seluruh 114 berkas**, medan `RequestType` pada setiap langkah,
hidup dipisahkan dari ber-remark.

**Pemanggilan `SaveOSClaim_SQL` — enam, empat di antaranya HIDUP:**

| Rule | Lgk | Gerbang | Keadaan |
| --- | --- | --- | --- |
| ⭐⭐ `KomitePost_Adjustment` | **8** | `AcceptStatus=="1"` **DAN** `KomiteCount==KomiteLoop` | ⭐ **HIDUP** |
| `KomitePost_CloseClaim` | 12 | `AcceptStatus==1` | **HIDUP** |
| `KomitePost_CloseClaim` | 12.3 | *(tanpa gerbang)* | **HIDUP** |
| `SaveReject_ACT_KMT` | 20 | *(tanpa gerbang)* | **HIDUP** |
| `SaveAccept_ACT` | 10 | `Local.SizeObjectItem<=1` | ⛔ ber-remark |
| `SaveAccept_ACT` | 11 | `Local.SizeObjectItem>1 && …==Param.Object` | ⛔ ber-remark |

⭐⭐ **Rancangannya justru rapi, dan saya yang gagal melihatnya:** `KomitePost_Adjustment` lgk
**7.2.1.17** memanggil `SaveAccept_ACT` untuk **MENYUSUN** ringkasan akseptasi di halaman
sementara, lalu lgk **8** — di rule pemanggil, bukan di dalamnya — **MENYIMPANNYA**.
⭐ **Langkah penyimpan diangkat ke pemanggil**, dan yang ber-remark di dalam `SaveAccept_ACT`
adalah **penyimpan lama yang digantikan**.

⭐ **Dan gerbangnya tepat:** `AcceptStatus=="1"` **DAN** `KomiteCount==KomiteLoop` —
**tersimpan sekali, oleh penyetuju terakhir.**

> ⛔ **RALAT terhadap ronde 3.** Kalimat lama dikutip: *"⛔ **Jadi PENOLAKAN tersimpan permanen;
> PENERIMAAN tidak.** ⚠️ Asimetri itu **tidak dapat saya jelaskan dari korpus**"*.
> ⭐ **Yang benar: tidak ada asimetri.** Keduanya tersimpan lewat `SaveOSClaim_SQL`; yang berbeda
> hanyalah **di rule mana langkah penyimpannya duduk**. ⚠️ **Butir 22 DITUTUP** — dan ⛔ **bukan
> karena bukti baru, melainkan karena sisiran ronde 3 saya terlalu sempit: saya memeriksa isi dua
> rule, bukan seluruh korpus.**

### 1c — Isi `SaveOSClaim_SQL`

`[terverifikasi]` Identitas: `ASM-FW-GCNMFW-INT-V_POLIS!ASM!SAVEOSCLAIM_SQL`, ruleset **GCNMFW
01-01-15**, sistem pengekspor **`pega`** *(bukan salinan pengembangan)*.

| | |
| --- | --- |
| panjang SQL | 283 aksara |
| `COMMIT` | ⚠️ **1** |
| `INSERT` / `UPDATE` / `DELETE` mentah | **0** / **0** / **0** |
| ⭐ yang dipanggil | ⭐ **`POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM`** |
| parameter | `InputData.CARI1` · `CARI2` · `CARI3` · `CARI10` · `CARI16` · `CARI17` · `CARI18` · `CARI20` · `CARI29` |

⭐ **Ia prosedur tersimpan yang menerima JSON** — `InputData.CARI1` diisi
`@ASM.GetPageJSONString()` di rule pemanggilnya. ⛔ **Apa yang dikerjakannya di dalam basis data
TIDAK TERBACA dari korpus** — ⭐ `[data DBA]`, disambungkan ke butir 1.

⚠️ **Yang terbaca dan penting:** ia membawa **`COMMIT;`**, dan `[terverifikasi]` ronde 2 modul ini
punya **NOL jalur kegagalan pada 25 langkah `Obj-*`**. ⛔ Pemanggilan di `SaveReject_ACT_KMT` lgk 20
**tanpa gerbang dan tanpa jalur gagal**.

### 1d — Halaman yang ditulis, kasus induk, dan kembaran

`[terverifikasi]` `SaveAcceptation_KMT` menulis ke **`Data` · `Local` · `TempOpenPage`**.
⛔ **NOL penugasan ke `Primary.*` maupun `pyWorkCover.*`** — **tidak menyentuh kasus induk.**

`[terverifikasi]` **Kembaran identitas empat bagian**, dicari di kedua modul saudara:

| Rule | Claim Fac In | Komite Claim Prop |
| --- | --- | --- |
| `SaveAcceptation_KMT` | ⛔ — | ⛔ — |
| `SaveAccept_ACT` | ⛔ — | ⛔ — |
| `SaveReject_ACT_KMT` | ⛔ — | ⛔ — |
| `SetValueKomite` | ⛔ — | ⛔ — |
| `KomitePostAct` | ⛔ — | ⛔ — |
| ⭐ **`HitServiceToKasirKMT_Act`** | ⛔ — | ⭐ **ADA — identik empat bagian** |

⭐⭐ **Lima rule inti modul ini KHAS — tak ada kembarannya di mana pun.** ⚠️ **Satu-satunya yang
dibagi adalah rule yang menyentuh UANG.** ⛔ Maka keputusan apa pun tentang `HitServiceToKasirKMT_Act`
**otomatis mengikat Komite Claim Prop juga**.

⚠️ **Nama yang mirip menipu:** ada `SaveAcceptation` di Claim Fac In *(20 langkah)*,
`SaveAcceptation_Act` dan `SaveAcceptationTreaty_TKMT` di Komite Claim Prop *(6 dan 16 langkah)*.
⭐ Keempatnya berkelas atau berisi **berbeda**; `SaveAcceptation_Act` kebetulan **juga 6 langkah**,
⛔ tetapi isinya sama sekali lain — ia menyusun `TempOSAkseptasi`, bukan bendera Retro.
⛔ **Jangan disamakan karena cacah langkahnya sama.**

---

## §2 — ⭐⭐⭐ Mekanisme berakhirnya putaran — **TERBUKTI**

⛔ **Ronde 3 menyerahkan ini sebagai lubang yang memblokir. Ia sekarang tertutup, dan Pega
ternyata MENGERJAKAN persis apa yang work owner tetapkan.**

### 2a — `KomiteRouter` dibaca sampai ke tag mentahnya

`[terverifikasi]` **8 langkah · 6 ber-remark · 2 hidup.** ⭐ **Identik di Komite Claim Prop** —
langkah demi langkah, gerbang demi gerbang, hanya kelasnya berbeda.

| Lgk | Gerbang | Isi | Keadaan |
| --- | --- | --- | --- |
| 1 | `.KomiteCount==1` | `param.AssignTo = "komitepnc"` | ⛔ remark |
| 2 | `.KomiteCount==2` | `param.AssignTo = "komitepnc2"` | ⛔ remark |
| 3 | `.KomiteCount==3` | `param.AssignTo = "komitepnc3"` | ⛔ remark |
| 4 | `.KomiteCount==4` | `param.AssignTo = "komitepnc4"` | ⛔ remark |
| 5 | — | `.KomiteCount = .KomiteCount+1` · `.AcceptStatus = ""` | ⛔ remark |
| ⭐ **6** | `Primary.TransferType=='2'` ⚠️ **bendera `false`** | ⭐ **ITERASI** | **HIDUP** |
| ⭐ **6.1** | `.KomiteAproval==0` | `param.AssignTo = .KomiteID` | **HIDUP** |
| 7 | `Primary.TransferType!='2'` | `param.AssignTo = .Komite.KomiteID` | ⛔ remark |

⭐⭐ **Langkah 6 adalah ITERASI**, dan itu baru terbaca dari tag mentahnya:

| Tag | Nilai |
| --- | --- |
| ⭐ `pyStepsObjectName` | ⭐ **`.KomiteList`** |
| ⭐ `pyStepsClassName` | ⭐ **`ASM-FW-GCNMFW-Data-Comitee`** |
| `pyStepsPreCondition` | ⚠️ **`false`** — ⛔ gerbang `TransferType=='2'` **MATI**, iterasi **selalu** jalan |
| `pyStepsPreCondParamsWhenFalse` | `3` *(Skip Step)* |
| `pyStepsTransParamsWhenTrue/False` | `2` / `2` *(Continue Whens)* |

**Dan langkah 6.1**, dengan catatan pengembangnya sendiri:

| Tag | Nilai |
| --- | --- |
| ⭐ `pyStepsDescription` | ⭐ **"set assign to jika aproval masih 0"** |
| `pyStepsPreCondParamsWhen` | `.KomiteAproval==0` |
| `pyStepsPreCondParamsWhenFalse` | `3` *(Skip Step)* |
| ⭐ `pyStepsTransParamsWhen` | `true` |
| ⭐ `pyStepsTransParamsWhenTrue` / `WhenFalse` | ⭐ **`6` / `6`** *(Exit Activity, kedua sisi)* |

### 2b — ⭐⭐ Mekanismenya, dinyatakan lengkap

> ⭐ **Router menelusuri `KomiteList` berurutan. Anggota PERTAMA yang `KomiteAproval == 0` menerima
> tugas, lalu aktivitas KELUAR SEKETIKA. Bila setiap anggota sudah memutuskan, gerbangnya
> melewatkan semua baris, iterasi selesai, dan `param.AssignTo` TIDAK PERNAH TERISI.**

⭐⭐ **Itulah pemutusnya**, dan ia **cocok dengan aturan work owner K1** tanpa perlu ditambal.

⛔ **Batas yang jujur:** apa yang dikerjakan Pega terhadap sebuah penugasan yang **tak menghasilkan
penerima** adalah **perilaku mesin Pega**, ⛔ **bukan sesuatu yang tertulis di korpus.**
⚠️ **Tidak ketemu di medan `pyRouteTo`, `pyRouterProp`, `pyRouteToType`, `pySkipAssignment`
pada bentuk `ASSIGNMENT63`** — bentuk itu hanya menyatakan `pyRouteTo = Custom` dan
`pyImplementation = WorkList`. ⭐ **Untuk sistem baru hal itu tidak penting**, sebab K1 sudah
menetapkan aturannya secara eksplisit.

⭐ **Akibat besar untuk rancangan:** ⛔ **`KomiteCount` BUKAN penggerak tangga.** Penggeraknya
adalah **`KomiteAproval` pada tiap anggota**. `KomiteCount` hanya **penunjuk giliran** dan
**penjaga efek akhir** *(`KomiteCount==KomiteLoop`)*. ⚠️ Ini **menghaluskan** kesimpulan ronde 3
§1d yang menyebut `.AcceptStatus` sebagai pemutus: `.AcceptStatus` memutus **putaran alur**,
`KomiteAproval` memutus **antrean penugasan**. **Keduanya benar, pada lapisan yang berbeda.**

### 2c — Penulis `KomiteAproval`, dua cara

`[terverifikasi]` **CARA-1 — teks mentah:** **7 berkas · 40 kemunculan.**
**CARA-2 — parser, sasaran dipisah dari pembaca:** **14 penugasan · 8 pembacaan.** ⭐ Sepakat pada
himpunan berkasnya.

| Rule | Lgk | Penugasan | Bendera | Catatan |
| --- | --- | --- | --- | --- |
| ⚠️ `ApprovalKomite_Act` | 7.1 | `KomiteList(<LAST>).KomiteAproval = 0` | ⚠️ **`false`** | ⛔ gerbang mati — lihat 2d |
| ⭐ `KomitePost_Adjustment` | **7.2.1.8** | `KomiteList(local.count).KomiteAproval = AcceptStatus` | ⚠️ **`false`** | ⭐ **inilah penanda "sudah memutuskan"** |
| `KomitePost_Adjustment` | 7.2.1.4 · 7.2.1.5 · 7.2.1.6 | `.ComiteeClaim(local.count).KomiteAproval = AcceptStatus` | `true` | cerminan ke data klaim |
| `KomitePost_Adjustment` | 7.2.1.9.1.1 | tiga penugasan bernilai `"2"` | `true` | jalur tolak beruntun |
| `KomitePost_CloseClaim` | 5 · 6 | `KomiteList(Local.Count)…` dan `ApprovalCommite.KomiteAproval` | — | |
| `KomitePost_Reject` | 6 · 7 | idem | — | |

⭐⭐ **Langkah 7.2.1.8 adalah engsel seluruh modul**, dan ⚠️ **gerbangnya
`pyWorkPage.pyStatusWork=="New"` MATI** — ⭐ **justru itulah yang membuat tangga bekerja**: penanda
"sudah memutuskan" dipasang **tanpa peduli status kasus**. ⛔ Bila gerbang itu hidup, tangga akan
macet pada kasus yang statusnya bukan `New`.

### 2d — ⚠️ Satu gerbang mati yang **membuka kembali** penyetuju terakhir

`[terverifikasi]` `ApprovalKomite_Act` lgk **7** bergerbang
`pyWorkPage.KomiteCount==1 && @LengthOfPageList(pyWorkPage.KomiteList) =1` *(bendera `true`)*,
dan di dalamnya lgk **7.1** menyetel `KomiteList(<LAST>).KomiteAproval = 0` dengan
⚠️ **bendera `false`** — ⛔ gerbangnya sendiri mati, jadi ia **selalu jalan** begitu bloknya masuk.

⭐ **Bacaan yang masuk akal:** untuk komite **beranggota satu**, penyetuju tunggal
**dibuka kembali** setiap kali layar disiapkan. ⚠️ **Tetapi perhatikan tanda sama dengan TUNGGAL**
pada `@LengthOfPageList(…) =1`. ⛔ **Apakah Pega membacanya sebagai pembandingan atau penugasan
TIDAK DAPAT SAYA PASTIKAN dari korpus** — ⭐ dan **bila itu penugasan**, blok tersebut masuk
**selalu**, bukan hanya untuk komite beranggota satu. ⭐ **`[terbuka]` BARU.**

---

## §3 — ⚠️ `HitServiceToKasirKMT_Act` — jalur uang keluar

⛔ **Kredensial: nilainya tidak dibaca, tidak dicetak, tidak disalin. Penyaring rahasia dipasang
pada instrumen pembacaan.**

### 3a — Bentuk luarnya

`[terverifikasi]` Kelas **`ASM-FW-GCNMFW-Data-Adjustment`**, ruleset **GCNMFW 01-01-28**, sistem
**`pegadevnusare2`**.

| | |
| --- | ---: |
| langkah | ⭐ **53** |
| ber-remark | **1** *(lgk 1)* |
| bendera `true` / kosong / **`false`** | 21 / 26 / ⚠️ **6** |
| ⭐ kedalaman sarang | ⭐ **6 tingkat** |
| parameter rule | ⛔ **NOL** |
| langkah `Java` | **2** |

⭐ **Terdalam di seluruh modul**, dan ⛔ **tanpa satu pun parameter** — seluruh masukannya dibaca
dari halaman.

### 3b — Urutan efeknya

`[terverifikasi]` Dua cabang besar, dipilih oleh `IsCLMNP`:

| Cabang | Lgk | Isi |
| --- | --- | --- |
| **Non-Prop** | 10 | gerbang `IsCLMNP` · ⚠️ `@length(Primary.AcceptedNo) = "23" \|\| … ` **kode 6 di sisi salah** |
| | 10.2 | `RDB-List` → `GetEmailCeding_SQL` *("GET EMAIL")* |
| | 10.6 | `Java` |
| | ⭐ **10.7** | ⭐ **`Connect-REST`** — gerbang `IsPEGAPROD`, catatan pengembang: *"kalau diserver dev jangan dijalanin"* · ⭐ **lompatan `StepStatusFail` → kode 1** |
| | 10.8 | `Primary.StatusKasir` disetel |
| | 10.9 | `RDB-List` → `InsertLOGDirectKasir_SQL` *("Insert Log Direct to Kasir")* |
| **Prop & Facin** | 14 | gerbang `IsCLMNP` ⭐ **kode 3/2** |
| | 14.1 / 14.2 | cabang `IsCLMP` dan `IsCLM`, masing-masing berpenyaring panjang `AcceptedNo` |
| | 14.3 | `Java` |
| | ⭐ **14.4** | ⭐ **`Connect-REST`** — gerbang `IsPEGAPROD` · ⭐ lompatan `StepStatusFail` |
| | 14.5 | `.StatusKasir` disetel |
| | 14.6 | `RDB-List` → `InsertLOGDirectKasir_SQL` |

⭐⭐ **Jawaban T3b: efek keluar mendahului penandaan, dan penandaan mendahului log.**
⭐ Urutannya **konsisten dan masuk akal**: panggil kasir → tandai `StatusKasir` → catat log.

⚠️ **Tetapi terhadap `Obj-Save`, polanya BERBEDA dari ronde 2.** `[terverifikasi]` Rule ini
**tidak punya `Obj-Save` sama sekali**. Pemanggilnya `KomitePost_Adjustment` memanggilnya di lgk
**25.2.1.1**, sedangkan `Obj-Save` pemanggil ada di lgk **7.2.1.10**, **7.2.1.11** *(sebelum)* dan
lgk **26** *(sesudah)*. ⭐ Jadi **ada simpanan sebelum DAN sesudah** efek keluar.
⛔ **RALAT ringan terhadap ronde 2**, yang menyimpulkan *"Obj-Save mendahului efek keluar"* sebagai
pola tunggal modul ini: ⭐ **pada jalur kasir, ada simpanan di kedua sisi.**

### 3c — Jalur kegagalan: ⭐ **ADA, dan ini pengecualian**

`[terverifikasi]` **Kedua `Connect-REST` membawa lompatan `StepStatusFail` berkode arah 1**
*(Jump to Later Step)*. ⭐⭐ **Ini satu-satunya tempat di modul ini yang menangani kegagalan efek
luar secara eksplisit** — ⚠️ bandingkan dengan **NOL dari 25** pada langkah `Obj-*`.

⚠️ **Tetapi penanganan galatnya cacat.** `[terverifikasi]` lgk **17**
`Call SendErrorDirectKasir`, bergerbang `.StatusServiceKasir.ReponseCode != "1"`,
⚠️ **bendera `false`** — ⛔ **gerbangnya MATI, jadi surel galat dikirim SETIAP KALI**, termasuk
ketika kasir menjawab berhasil.

⚠️ Perhatikan pula ejaan medannya: **`ReponseCode`**, bukan `ResponseCode`. ⛔ **Bila itu salah
ketik, gerbangnya tak akan pernah benar walau dihidupkan** — ⭐ **dan itu mungkin justru ALASAN
gerbangnya dimatikan.** ⛔ **Belum punya data**; dicatat sebagai `[terbuka]` BARU.

### 3d — Penjaga ganda-bayar: ⭐ **ADA**

`[terverifikasi]` lgk **2** bergerbang **`.DirectToKasir=="true" && .StatusKasir = ""`**, dan
gerbang yang **persis sama** menjaga pemanggilnya di `KomitePost_Adjustment` lgk **25.2.1.1**.
⭐ Digabung dengan lgk 10.8 / 14.5 yang **menyetel `StatusKasir` sesudah panggilan berhasil**,
⭐ **penjaga ganda-bayarnya lengkap: sekali dibayar, `StatusKasir` terisi, gerbang menutup.**

⛔⛔ **TETAPI — dan inilah yang diminta T3d — perhatikan tanda sama dengan TUNGGAL:**
`.StatusKasir = ""`, bukan `.StatusKasir == ""`.

⭐ **Laporan jujur:** ⛔ **saya TIDAK DAPAT memastikan dari korpus** apakah Pega membacanya sebagai
**pembandingan** atau **penugasan**. ⚠️ **Bila penugasan**, gerbang itu **selalu bernilai benar**
dan **penjaga ganda-bayar tidak ada** — sekaligus **menghapus** `StatusKasir` setiap kali gerbang
dievaluasi. ⭐ Pola yang sama muncul **tiga kali lagi**: `@length(Primary.AcceptedNo) = "23"`,
`@length(.AcceptedNo) = "21"`, dan `@LengthOfPageList(…) =1` di §2d.

⭐⭐ **`[terbuka]` BARU — MEMBLOKIR.** ⛔ Ini **bukan** pertanyaan gaya penulisan: ia menentukan
**apakah sistem lama punya penjaga ganda-bayar atau tidak**, dan **apakah sistem baru wajib
menambahkannya**. ⚠️ Jawabannya ada pada **perilaku mesin Pega**, ⛔ bukan pada korpus — jadi ia
menuntut **uji di lingkungan Pega** atau **keterangan pengembang lama**.

⚠️ **Penjaga sejenis di tempat lain** `[terverifikasi]`: `.IsPrintAccept` menjaga pencetakan
akseptasi — `KomitePost_Adjustment` lgk 7.2.1.16 *(`.AcceptanceStatus=="1" && .IsPrintAccept==""`)*
dan **sembilan** gerbang di `PrintPDFAccep_MultiAksep_KMT` lgk 10.3.1–10.3.9. ⭐ **Semuanya memakai
`==` yang benar.**

### 3e — Dua langkah `Java`

`[terverifikasi]` lgk **10.6** dan **14.3**, keduanya **tepat sebelum** `Connect-REST`-nya
masing-masing. ⭐ **Bentuknya: penyusun muatan** — merakit halaman JSON untuk dikirim.
⛔ **Isinya tidak disalin.**

### 3f — Kredensial

`[terverifikasi]` `ConnectREST\SendAcceptationToKasir` memetakan **`ParamKasir.CARI1`** lewat
medan `pyMapFromKey`. ⭐ **Itu saja yang saya sebut: rule-nya, medannya, dan bahwa pemetaannya ada.**
⛔ **Nilai kredensial tidak dibaca, tidak dicetak, tidak disalin — penyaring rahasia dipasang pada
instrumen agar perintah tidak dapat menampilkannya walau tak sengaja.**

⚠️ **Peringatan tetap berlaku:** karena kredensial beredar di dalam berkas ekspor, **ia wajib
diputar ulang sesudah alih** — tanggung jawab tim pemilik layanan.

---

## §4 — `KomitePostAct`, `SetValueKomite`, dan layar

### 4a — ⭐ `KomitePostAct`: **empat cabang, tidak ada yang kelima**

`[terverifikasi]` **6 langkah · 2 ber-remark**, kelas `Work-Komite`, ruleset GCNMFW 01-01-24.

| Lgk | Isi | Gerbang | Keadaan |
| --- | --- | --- | --- |
| 1 | `Obj-Refresh-And-Lock` | — | HIDUP |
| 2 | `Call KomitePost_Adjustment` | `TransferType=="2"` | HIDUP |
| 3 | `Call KomitePost_Survey` | `TransferType=="1"` | ⛔ **REMARK** |
| 4 | `Call KomitePost_Reject` | `TransferType=="3"` | HIDUP |
| 5 | `Call KomitePost_CloseClaim` | `TransferType=="4"` | HIDUP |
| 6 | `.AcceptStatus = ""` · `.Comment = ""` | — | ⛔ **REMARK** |

> ✅ **Tepat empat cabang. Tidak ada cabang kelima.** ⭐ Ronde 2 benar.

⭐ **Dan lgk 1 penting:** modul ini **mengunci kasus induk** sebelum mengerjakan apa pun.

### 4b — ⭐⭐ `SetValueKomite`: ia **penyiap layar**, bukan penyalin data

⛔ **Ronde 3 hanya mengenal langkah 12-nya. Rule ini punya 31 LANGKAH dan 36 PENUGASAN.**

`[terverifikasi]` Kelas `Work-Komite`, ruleset GCNMFW 01-01-24.

| Tahap | Lgk | Yang dikerjakan |
| --- | --- | --- |
| bersih-bersih | 2 · 3 | `Page-Remove` · `Property-Remove` |
| ambil konteks | 4 · 5 | indeks objek/item/adjustment · `TempObject.OccupationName` dari induk · `Page-Copy` |
| kumpulkan adjustment | 6 – 7.1.1 | iterasi bergerbang `.KomiteNo==pyWorkPage.pyID`, menyusun `DataTempAdj.Adjustment` berikut catatan *"ObjectItem ke …"* |
| siapkan baris baru | 8 – 8.3 | bergerbang `.AcceptanceStatus=="" \|\| …` · `local.Initial = .pxCreateOpName` |
| hitung total | 9 – 11.2 | `Java` di lgk 10 · penjumlahan bruto & bagian RNM, **berikut konversi ke IDR** · cabang `PaymentType` 3, 4, 6 |
| ⭐ **salin ke objek komite** | ⭐ **12** | ⭐ lihat tabel di bawah |
| total dalam IDR | 13 | `TempTotalAdj.pxResults` |
| ⭐ **panggil penyiap wewenang** | ⭐ **14 · 15** | ⭐ `Call ApprovalKomite_Act` · `Call SetProteksiSubmiteKomite` |

⭐⭐ **Langkah 12 — SELURUH yang disalin ke objek kerja komite, lima medan:**

| Sasaran | Sumber |
| --- | --- |
| ⭐ `pyWorkPage.Quotation.BusinessType` | `pyWorkCover.OfferFacIn.QuotationData.BusinessType` |
| ⭐ `pyWorkPage.Quotation.GroupPanel` | `pyWorkCover.OfferFacIn.QuotationData.GroupPanel` |
| `pyWorkPage.CLMNO` | `pyWorkCover.pyID` |
| `pyWorkPage.Komite.DateOfComitee` | `pyWorkPage.pxCreateDateTime` |
| `pyWorkPage.Komite.Initial` | `Local.Initial` |

> ⭐ **Ronde 3 TERVERIFIKASI ULANG dan BENAR: hanya `BusinessType` dan `GroupPanel` dari
> `Quotation` yang disalin.** ⛔ `BusinessCode` dan `BusinessName` **tidak**.

⭐ **Dan lgk 14–15 menjelaskan bentuk modul:** `SetValueKomite` adalah **pintu masuk layar komite**
— ia menyiapkan data, lalu memanggil **penyaring wewenang** dan **penyaring tombol Submit**.

### 4c — ⭐⭐⭐ `BusinessCode`: ronde 3 **BERTAHAN**, dan mekanismenya kini lengkap

⚠️ Ronde 3 menyebut ini **kesimpulannya yang paling rawan salah** dan menantang ronde 4
mematahkannya. ⛔ **Ia tidak patah. Ia menguat.**

`[terverifikasi]` **Jendela: 2 korpus (114 + 482 berkas) · 15 nama tag sasaran · semua jenis rule.**
Tag yang disisir: `PropertiesName` · `pyPropertyName` · `pySetPropertyName` · `pyTargetProperty` ·
`pyPropName` · `pxPropertyName` · `pyName` · `pyPropertiesName` · `pySourceName` · `pyMapToKey` ·
`pyResultProperty` · `pyDataTransformTarget` · `pyReference` · `pyFieldName` · `pyColumnName`.

| Korpus | Penulis ditemukan |
| --- | ---: |
| ⭐ **Komite Claim FacIn** | ⛔ **NOL** |
| Claim Fac In | 22 — ⛔ **tak satu pun menulis `Quotation.BusinessCode`** |

⭐ **Kedua puluh dua itu ternyata bukan penugasan sama sekali:** delapan `pyName` di
`Harness\ChoosePolis`, delapan `pyName` di `Section\ViewPolis` *(nama medan layar)*, dan enam
`pyTargetProperty`/`pyFieldName` di `ReportDefinition\BrowseTreatyGroup_RD` untuk properti
**berbeda** — `.OJKBusinessName` dan `.OJKBusinessNameIDN`.

⭐⭐ **Dan inilah mekanisme yang akhirnya menjelaskan segalanya** `[terverifikasi]`:

| Modul | Yang mengisi `Quotation` | Bentuk |
| --- | --- | --- |
| ⭐ **Claim Fac In** | `CopyNB_Act` lgk 8 · `UpdateDataPolisFacin_Act` lgk 7 | ⭐ **`pyWorkPage.Quotation = pyWorkPage.OfferFacIn.QuotationData` — HALAMAN UTUH** |
| ⚠️ **Komite Claim FacIn** | `SetValueKomite` lgk 12 | ⛔ **dua medan bernama saja** |
| | `KomitePost_Reject` lgk 2 · `SaveReject_ACT_KMT` lgk 1 | halaman utuh — ⛔ **jalur TOLAK saja** |

⭐⭐ **Jadi: di Claim Fac In, `BusinessCode` ikut terbawa oleh salinan halaman utuh. Di Komite, ia
tidak pernah ikut pada jalur terima.** ⛔ **Bukan karena ada yang menghapusnya — karena tak pernah
ada yang menyalinnya.**

> ✅ **Vonis ronde 3 §5d BERTAHAN, dan kini bersandar pada mekanisme, bukan pada ketiadaan saja.**
> ⭐ Keputusan **K2** *(§6)* menutup lubang ini.

⚠️ Catatan pendukung: `BusinessCode` **dibaca** sebagai nilai di empat Activity Claim Fac In —
semuanya dari **`OfferFacIn.QuotationData.BusinessCode`** atau `Policy.Quotation.BusinessCode`,
⭐ **bukan** dari `pyWorkPage.Quotation.BusinessCode`. ⭐ **Sumber sejatinya memang
`OfferFacIn.QuotationData`.**

### 4d — ⭐⭐ Layar `ShowTransfer`: **butir 23 DITUTUP**, dan hasilnya jauh lebih sempit

⛔ **Ronde 3 gagal membaca medannya** — `pyReference`, `pyPropertyName`, `pyFieldName`,
`pyFieldType` semuanya nol. ⭐ **Tag yang benar: `pyValue` dan `pyLabelFor`.**

`[terverifikasi]` **92 medan `pyValue` unik · 41 `pyLabelFor` unik.**

⭐⭐ **Dan yang dikunci oleh kedua gerbang `.KomiteCount!='1'` ternyata TEPAT DUA:**

| Gerbang | Medan | Label |
| --- | --- | --- |
| `.KomiteCount!='1'` | ⭐ **`.Adjustment.IsProposeClose`** | `IsProposeClose` |
| `.KomiteCount!='1'` | ⭐ **`.Adjustment.IsPropReserved`** | `IsPropReserved` |
| `ProteksiKomite.CARI1==0 \|\| pyWorkPage.Adjustment.AcceptedNo != ''` | ⭐⭐ **tombol `Submit`** | `Submit` |

> ⛔ **RALAT terhadap ronde 3 §4b.** Kalimat lama dikutip: *"⭐ **Artinya sebuah aturan bisnis yang
> belum pernah tercatat: medan yang dijaga kedua gerbang itu hanya dapat disunting oleh penyetuju
> PERTAMA.**"* ⚠️ Kalimatnya **tidak salah**, ⛔ **tetapi nadanya membesar-besarkan**: yang dijaga
> bukan isian layar pada umumnya, melainkan **dua kotak centang usulan** — *usul menutup klaim* dan
> *usul mencadangkan*. ⭐ **Aturannya nyata, cakupannya kecil.**

⭐⭐ **Dan gerbang ketiga adalah konfirmasi paling bersih untuk ADR-0014 sejauh ini:** pemeriksaan
`ProteksiKomite.CARI1` menempel pada **TOMBOL SUBMIT**. ⛔ **Persis seperti yang ronde 2 dan 3
simpulkan — penegakan wewenang duduk di layar, bukan di lapisan layanan.** ✅ Vonis
**MENENTANG SEBAGIAN** berdiri dengan bukti langsung.

⚠️ **Satu temuan sampingan yang rapi:** layar ini membaca
**`pyWorkCover.OfferFacIn.QuotationData.BusinessName`** — ⭐ **langsung dari kasus induk**.
⛔ Itulah sebabnya tak seorang pun pernah menyadari `pyWorkPage.Quotation.BusinessName` kosong:
**layar tidak memakainya, hanya rule `When` yang memakainya.**

---

## §5 — Dua butir sisa

### 5a — Butir **24**: larangan menyetujui diri sendiri

`[terverifikasi]` Jendela: pola `KomiteList(<indeks>)` pada seluruh berkas dua modul.

| Modul | Indeks yang dipakai |
| --- | --- |
| Komite Claim FacIn | `Local.Count` ×22 · `pyWorkPage.KomiteCount` ×6 · `Local.IdxKmt` ×4 · `local.count` ×3 · `<LAST>` ×3 · ⭐ **`1` ×1** · `<APPEND>` ×1 |
| Komite Claim Prop | `Local.Count` ×22 · `local.count` ×3 · `pyWorkPage.KomiteCount` ×2 · `Local.IdxKmt` ×2 |

⭐ **Pemeriksaan pemilik terhadap akun operator: TEPAT SATU di seluruh dua modul** —
`ApprovalKomite_Act` `pyWorkPage.KomiteList(1).KomiteID == .OPERATOR_ID`.

⛔⛔ **Dan Komite Claim Prop TIDAK PUNYA sama sekali** — ⚠️ indeks `1` bahkan tidak muncul di sana.
⭐ **Jadi larangan menyetujui diri sendiri hanya ada di Fac In, hanya untuk anggota PERTAMA.**

⚠️ **Butir 24 tetap terbuka** *(keputusan work owner)*, ⭐ **tetapi kini lebih tajam:** pertanyaannya
bukan lagi "kenapa indeks 1", melainkan **"apakah larangan ini memang hanya untuk anggota pertama,
dan kenapa modul saudaranya tidak punya sama sekali"**.

### 5b — Butir **14**: `KomitePost_Survey` — ⭐ **DITUTUP**

`[terverifikasi]` Jendela: teks mentah seluruh berkas **ketiga** korpus.

| Modul | Menyebut namanya | Berkas definisinya |
| --- | --- | --- |
| Komite Claim FacIn | `Activity\KomitePostAct` — **2 kali** *(langkah ber-remark)* | ⛔ **TIDAK ADA** |
| Komite Claim Prop | **NOL** | ⛔ **TIDAK ADA** |
| Claim Fac In | **NOL** | ⛔ **TIDAK ADA** |

⭐ **Jalur `TransferType=="1"` (Survey) memang mati seluruhnya:** pemanggilnya ber-remark, dan
**rule-nya tidak ada di korpus mana pun**. ✅ **Butir 14 DITUTUP** — ⛔ **jangan dialihkan.**

---

## §6 — ✅ Empat keputusan work owner — 2026-09-20

### K1 → butir **19** ⭐ **DITUTUP** — kapan komite berakhir

> ⭐ **`[keputusan work owner]` 2026-09-20 — Komite berakhir ketika SELURUH jenjang akseptasi
> menyetujui. Ketika ada satu yang menolak, kasus komite LANGSUNG tertutup — tidak menunggu
> jenjang berikutnya.**

⭐ **Sejalan rekomendasi asisten**, dan ⭐⭐ **§2 membuktikan Pega memang MENGERJAKANNYA BEGITU** —
router berhenti menghasilkan penerima ketika setiap anggota sudah memutuskan, dan `.AcceptStatus=="2"`
jatuh ke penghubung `NoLoop` menuju `END52`.

⚠️ **Yang WAJIB ikut ditulis:** di sistem baru, **pemutusnya dinyatakan EKSPLISIT** —
*seluruh anggota sudah memutuskan* **atau** *ada satu yang menolak*. ⛔ **Jangan meniru bentuk
Pega yang menyandarkan pemutusan pada `param.AssignTo` yang kebetulan tidak terisi** — ⭐ itu
**perilaku diam** yang tidak dapat diuji dan tidak dapat dibaca.

⭐ **Dan satu koreksi arah terhadap brief ronde ini sendiri:** brief menulis *"pemutus putaran
adalah PENCACAH (`KomiteCount == KomiteLoop`)"*. ⚠️ **Lebih tepat: pemutusnya adalah KEADAAN TIAP
ANGGOTA** *(`KomiteAproval`)*, dan `KomiteCount == KomiteLoop` **kebetulan setara** selama pencacah
naik sekali per keputusan. ⛔ **Keduanya setara hanya bila daftar anggota tidak berubah di
tengah jalan** — dan `ApprovalKomite_Act` lgk 6.2 **memang dapat mengubahnya** *(§3c ronde 3)*.
⭐ **Pakailah keadaan tiap anggota sebagai kebenaran, pencacah sebagai tampilan.**

### K2 → butir **20** ⭐ **DITUTUP** — cabang MBU dan Travel

> ⭐ **`[keputusan work owner]` 2026-09-20 — Diperlakukan sebagai CACAT dan diperbaiki saat alih:
> `BusinessCode` dan `BusinessName` DISALIN ke objek kerja komite bersama `BusinessType`.**

⚠️ **`[penyimpangan sadar]`** — ⭐ **dicatat sebagai CACAT YANG DIPERBAIKI, bukan perilaku yang
ditiru.**

**Alasannya, kini berdasar mekanisme penuh** `[terverifikasi]` §4c: di Claim Fac In
`Quotation` diisi **salinan halaman utuh** dari `OfferFacIn.QuotationData`, sehingga `BusinessCode`
ikut; di Komite hanya **dua medan bernama** yang disalin. ⛔ **Tak ada yang menghapusnya — tak
pernah ada yang menyalinnya.** ⭐ Sembilan rule `When` menguji medan itu, **dua dipakai hidup lima
kali** *(`IsMBU`, `IsTravel`)*.

⭐ **Cara terbersih menutupnya di sistem baru:** salin **seluruh data kutipan** yang dibutuhkan
komite dari kasus induk, ⛔ bukan memperpanjang daftar medan satu per satu — sebab daftar bernama
itulah yang **melahirkan cacat ini**.

⚠️ **Yang WAJIB ikut ditulis:** perbaikan ini **mengubah perilaku terhadap data lama**. Klaim MBU
dan Travel yang **sudah melewati komite** punya ringkasan akseptasi yang dibangun **tanpa** cabang
itu. ⭐ **`[terbuka]` BARU:** apakah data lama dibangun ulang, atau dibiarkan.

### K3 → butir **21** ⭐ **DITUTUP** — alias akun

> ⭐ **`[keputusan work owner]` 2026-09-20 — Pemilik giliran dibandingkan lewat IDENTITAS AKUN,
> bukan alias. `KomiteID` disimpan sebagai rujukan ke pengguna, dan pemetaan nama yang ditanam di
> dalam kode TIDAK DIBAWA.**

**Alasannya** `[terverifikasi]` ronde 3 §3d: `Local.USERNAME` — sisi kiri satu-satunya pemeriksaan
pemilik giliran — adalah hasil ungkapan `@if` yang **memetakan akun menjadi nama lain di dalam
kode**. ⛔ Maka pergantian nama tampil **menggeser wewenang** tanpa ada yang mengubah aturan.

⚠️ **Yang WAJIB ikut ditulis:** bila pemetaan itu punya alasan historis yang masih hidup, ia
menjadi **baris data**, ⛔ **bukan cabang `@if` di dalam kode.** **`[terbuka]`:** alasan
historisnya belum diketahui.

### K4 → butir **23** ⭐ **DITEGASKAN, DITIRU, dan kini TERBACA**

> ⭐ **`[keputusan work owner]` 2026-09-20 — Hanya penyetuju PERTAMA yang boleh menyunting isian.
> Penyetuju kedua dan seterusnya MENILAI, tidak mengubah. Ditiru apa adanya.**

⭐⭐ **§4d menutup butir 23 dan membuat keputusan ini dapat dibangun.** Yang dikunci **tepat dua
kotak centang**:

| Medan | Arti |
| --- | --- |
| ⭐ `.Adjustment.IsProposeClose` | usul **menutup klaim** |
| ⭐ `.Adjustment.IsPropReserved` | usul **mencadangkan** |

⭐ **Aturannya jadi jernih dan sempit:** **usulan tindak lanjut dibuat oleh penyetuju pertama;
penyetuju berikutnya menyetujui atau menolak usulan itu, tidak menggantinya.**

⚠️ **Yang WAJIB ikut ditulis:**
1. ⛔ Pembandingnya di Pega adalah **`'1'` bertanda kutip — teks**, padahal `KomiteCount` diisi
   **bilangan**. ⭐ Pega memaafkan; **Go tidak akan.** **Jebakan alih.**
2. ⭐ **Sisa 90 medan layar tidak dikunci oleh gerbang ini** — ⛔ jangan menguncinya karena salah
   membaca aturan ini sebagai "layar terkunci untuk penyetuju kedua".

---

## §7 — Penutup

### 7.1 — ⛔ Kesimpulan yang PALING RAWAN SALAH

> ⚠️⚠️ **§3d — bahwa penjaga ganda-bayar ADA.**

**Kenapa ia rawan.** Ia bergantung sepenuhnya pada bagaimana Pega membaca **tanda sama dengan
tunggal** di `.StatusKasir = ""`. ⭐ Bila itu **pembandingan**, penjaganya utuh. ⛔ Bila itu
**penugasan**, gerbangnya **selalu benar**, penjaganya **tidak ada**, dan ⚠️ **`StatusKasir` justru
DIKOSONGKAN setiap kali gerbang dievaluasi** — yang berarti sistem lama **dapat membayar dua kali**.

⛔ **Saya tidak dapat memutuskannya dari korpus**, dan ⛔ **saya menolak menebak** pada hal yang
menyentuh uang.

⭐ **Cara mematahkannya:** jalankan satu uji di lingkungan Pega dengan `StatusKasir` terisi, dan
lihat apakah langkah 2 masuk. **Satu uji cukup.** ⚠️ Atau: tanyakan kepada pengembang lama.

**Peringkat kedua:** §2b — mekanisme berhentinya router. ⭐ Buktinya kuat *(tag `pyStepsObjectName`,
`pyStepsClassName`, kode arah 6/6, dan catatan pengembang sendiri)*, ⛔ **tetapi langkah
terakhirnya — apa yang Pega lakukan pada penugasan tanpa penerima — adalah perilaku mesin, bukan
korpus.** ⭐ Untungnya **K1 membuatnya tidak penting** bagi sistem baru.

### 7.2 — ❓ Pertanyaan untuk work owner

❓ **Q1** — **Tanda sama dengan tunggal: pembandingan atau penugasan?** `[terverifikasi]` Pola
`X = ""` / `X = "23"` / `X =1` muncul di **empat gerbang**, salah satunya **penjaga ganda-bayar**.
**Apakah ada yang dapat menguji ini di lingkungan Pega, atau menanyakannya kepada pengembang lama?**

➡️ **Rekomendasi: jangan menunggu jawabannya untuk membangun.** ⭐ Di sistem baru, **penjaga
ganda-bayar dibuat EKSPLISIT dan tidak boleh bergantung pada tafsir apa pun** — sekali pembayaran
terkirim, penandanya tersimpan **dalam transaksi yang sama**, dan panggilan kedua ditolak.
⛔ Jawabannya tetap perlu dicari, ⭐ **tetapi untuk mengetahui apakah data lama perlu diperiksa
ganda-bayar** — bukan untuk menentukan rancangan.

---

❓ **Q2** — **Surel galat kasir terkirim setiap kali.** `[terverifikasi]` lgk 17 bergerbang
`.StatusServiceKasir.ReponseCode != "1"` dengan **bendera `false`** — gerbang mati.
⚠️ Ejaan medannya **`ReponseCode`**. **Apakah tim operasi memang menerima surel galat pada setiap
pembayaran, termasuk yang berhasil?**

➡️ **Rekomendasi: jangan tiru.** ⭐ Di sistem baru, **beri tahu hanya pada kegagalan**, dan
**namai medannya dengan benar**. ⚠️ Jawaban work owner tetap berharga — **bila ternyata surel itu
memang tak pernah sampai kepada siapa pun**, maka **penanganan galat kasir sebenarnya TIDAK ADA**,
dan itu butir yang jauh lebih besar.

---

❓ **Q3** — **Fac Retro melompati penyiapan wewenang.** `[terverifikasi]` `TreatyType == "10015"`
→ `.IsFacRetro = 1` → `ApprovalKomite_Act` **`Exit-Activity`**. **Apakah klaim Fac Retro memang
tidak melalui pemeriksaan wewenang komite, dan apa arti `10015`?**

➡️ **Rekomendasi: jadikan `TreatyType` DATA, jangan tetapan.** ⭐ Berbeda dari pita nilai di ronde 3
— yang kamu putuskan ditiru apa adanya — ⚠️ **angka ini menentukan JALUR, bukan ambang**, dan
melompati pemeriksaan wewenang adalah **keputusan kendali**, bukan kebijakan nilai. ⛔ Bila ia tetap
tetapan, **tidak akan ada yang tahu jalur itu ada** sampai sesuatu terjadi.

---

❓ **Q4** — **Larangan menyetujui diri sendiri hanya untuk anggota pertama, dan tidak ada sama
sekali di Komite Claim Prop.** *(butir 24)* **Mana yang benar?**

➡️ **Rekomendasi: berlakukan untuk SEMUA anggota, di kedua modul.** ⭐ Aturan "tidak boleh
menyetujui klaim yang diajukan sendiri" yang **hanya berlaku pada satu posisi daftar** bukan
aturan — ⛔ **ia lubang dengan penjaga di satu pintu.** ⚠️ Bila work owner menyatakan larangan itu
memang tidak dikehendaki, **buang seluruhnya** — ⛔ jangan tinggalkan setengahnya.

---

❓ **Q5** — **Data lama MBU dan Travel.** Akibat K2, ringkasan akseptasi lama untuk kedua lini itu
dibangun **tanpa** cabangnya. **Perlu dibangun ulang, atau dibiarkan?**

➡️ **Rekomendasi: biarkan, tetapi TANDAI.** ⭐ Membangun ulang catatan akseptasi historis menyentuh
angka yang mungkin sudah dilaporkan keluar. ⚠️ Yang perlu dikerjakan: **ketahui berapa banyak**,
dan **tandai** supaya siapa pun yang membacanya tahu ia dibangun dengan aturan lama.
⛔ **Cacahnya `[data DBA]`** — tidak terbaca dari korpus.

### 7.3 — Register `[terbuka]` sesudah ronde 4

**✅ DITUTUP — tujuh:**

| # | Butir | Ditutup oleh |
| --- | --- | --- |
| ⭐ **14** | `KomitePost_Survey` ada berkasnya? | §5b — ⛔ **tidak ada di korpus mana pun** |
| ⭐⭐ **19** | Bagaimana putaran berakhir | §2 + **K1** — ⭐ mekanisme **terbukti** |
| ⭐⭐ **20** | `IsMBU`/`IsTravel` membaca medan tak disalin | §4c + **K2** |
| **21** | Alias akun ditanam keras | **K3** |
| ⭐⭐ **22** | `SaveAccept_ACT` tidak menyimpan | §1b — ⛔ **vonis ronde 3 BATAL**, penerimaan **tersimpan** |
| ⭐ **23** | Medan `ShowTransfer` yang dikunci | §4d — ⭐ **tepat dua kotak centang** |
| **25** | Bentuk jalur administratif pengganti pintu belakang | ⛔ **TIDAK ditutup** — lihat catatan |

⚠️ **Butir 25 TIDAK ditutup.** ⛔ Ia masuk tabel karena kekeliruan penyusunan dan **saya perbaiki
di sini alih-alih menghapusnya**: ⭐ **yang ditutup ada ENAM, bukan tujuh** — 14 · 19 · 20 · 21 ·
22 · 23.

**⭐ BARU — enam:**

| # | Butir | Pemilik | Memblokir? |
| --- | --- | --- | --- |
| ⚠️⚠️ **26** | **Tanda sama dengan TUNGGAL** di 4 gerbang — pembandingan atau penugasan; menentukan ada-tidaknya **penjaga ganda-bayar** | work owner / pengembang lama | ⛔ **YA** — Q1 |
| ⚠️ **27** | **Surel galat kasir terkirim selalu** — gerbang mati + ejaan `ReponseCode` | work owner | tidak |
| ⚠️ **28** | **`TreatyType == "10015"`** — arti, dan kenapa Fac Retro melompati penyiapan wewenang | work owner | tidak — Q3 |
| **29** | **Data lama MBU/Travel** — dibangun ulang atau dibiarkan | work owner | tidak — Q5 |
| **30** | **`ApprovalKomite_Act` lgk 7/7.1** — membuka kembali penyetuju terakhir; bergantung butir 26 | asisten → work owner | tidak |
| **31** | **Isi `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM`** | `[data DBA]` | tidak — digabung ke butir 1 |

**Sisa yang masih terbuka dari ronde sebelumnya: 1 · 6 · 11 · 12 · 13 · 15 · 24 · 25.**
*(15 tetap tidak memblokir sejak ronde 3.)*

⭐ **Total, dihitung DUA CARA:**
*(a)* **cacah baris** ⇒ sisa `8` + baru `6` = **14**;
*(b)* **delta** ⇒ `14 − 6 ditutup + 6 baru` = **14**.
✅ **Keduanya sepakat: 14 butir terbuka.**

⛔ **Yang MEMBLOKIR: 1 · 6 · 26** — tiga butir.
*(1 = kolom `T_QUOTATIONDATA` `[data DBA]`; 6 = isi rule peran, lintas tiga modul; 26 = tanda sama
dengan tunggal pada penjaga ganda-bayar.)*

⛔ **Nomor 18 tetap dibatalkan dan tidak dipakai ulang.**

### 7.4 — ⛔ RALAT

| # | Yang diralat | Kalimat LAMA, dikutip | Yang benar |
| --- | --- | --- | --- |
| ⛔⛔ **1** | **ronde 3 §2b** | *"⛔ **Jadi PENOLAKAN tersimpan permanen; PENERIMAAN tidak.** ⚠️ Asimetri itu **tidak dapat saya jelaskan dari korpus**"* | ⭐ **Tidak ada asimetri.** `SaveOSClaim_SQL` **hidup empat kali**, termasuk `KomitePost_Adjustment` lgk 8 bergerbang `AcceptStatus=="1" && KomiteCount==KomiteLoop`. ⛔ Sisiran ronde 3 hanya memeriksa **isi dua rule**, bukan **seluruh korpus** |
| ⭐ **2** | **ronde 3 §4b** | *"medan yang dijaga kedua gerbang itu hanya dapat disunting oleh penyetuju PERTAMA"* | ⭐ **Benar, tetapi nadanya membesar-besarkan.** Yang dikunci **tepat dua kotak centang usulan**, bukan isian layar pada umumnya. **90 medan lain tidak terkunci** |
| ⭐ **3** | **ronde 2 §1** tentang urutan simpan | *"`Obj-Save` mendahului efek keluar"* sebagai pola modul | ⭐ **Benar untuk jalur akseptasi; tidak lengkap untuk jalur kasir.** Di sana ada simpanan **sebelum** *(lgk 7.2.1.10/11)* **dan sesudah** *(lgk 26)* |
| ⭐ **4** | **brief ronde 4 ini sendiri** | *"pemutus putaran adalah PENCACAH (`KomiteCount == KomiteLoop`)"* | ⭐ **Lebih tepat: keadaan tiap anggota (`KomiteAproval`).** Setara dengan pencacah **hanya bila daftar anggota tak berubah**, padahal `ApprovalKomite_Act` lgk 6.2 dapat mengubahnya |

⚠️ **RALAT nomor 1 adalah yang ketiga kalinya dalam modul ini sebuah pernyataan KETIADAAN dari saya
terbukti terlalu cepat.** ⭐ Dua yang pertama disebabkan **penyaring bertag tunggal**; ⛔ **yang ini
disebabkan JENDELA yang terlalu sempit — saya membaca dua rule dan menyimpulkan tentang modul.**

⭐ **Aturan turunan, berlaku sejak sekarang:** ⛔ **sebuah pernyataan tentang MODUL wajib disisir
pada tingkat MODUL.** Membaca isi satu rule hanya mengizinkan pernyataan tentang **rule itu**.

✅ Kalimat-kalimat lama **tidak dihapus** dari berkas ronde 2 dan 3; ia dikutip di sini berdampingan
dengan koreksinya.

### 7.5 — ⭐ Vonis kesiapan spec

> ⭐⭐ **SIAP — dengan tiga butir dinyatakan terbuka di muka.**

⛔ **Ini berubah dari ronde 3, yang memvonis BELUM.** Yang berubah: ⭐ **butir 19 dan 20 — keduanya
memblokir — sudah tertutup**, dan ⭐ **mekanisme modulnya kini terbaca ujung ke ujung.**

**Yang sudah cukup untuk menulis spec:**

| Sudah cukup | Dari |
| --- | --- |
| daur hidup penuh: penyemaian, penugasan berputar, pemutus, penutupan | ronde 3 §1 · ronde 4 §2 |
| aturan komite: semua setuju ⇒ selesai · satu menolak ⇒ tutup seketika | **K1** |
| penyimpanan keputusan: siapa, kapan, dengan gerbang apa | §1b |
| wewenang: satu pemeriksaan pemilik di tombol Submit, satu larangan diri-sendiri | §4d · §5a |
| layar: 92 medan, dua terkunci, aturannya jernih | §4d · **K4** |
| lini usaha: 13 `When` hidup, 36 beban mati, cacat `BusinessCode` beserta perbaikannya | ronde 3 §5 · §4c · **K2** |
| penomoran akseptasi: satu rangkaian dikunci kelas induk | ronde 3 §6b |
| jalur uang: dua cabang, dua panggilan luar, penjaga ganda-bayar, penanganan galat | §3 |

⛔ **Tiga butir yang dibawa TERBUKA ke dalam spec — dan masing-masing dapat dikurung:**

1. **Butir 26** *(tanda sama dengan tunggal)* — ⭐ **dikurung** oleh rekomendasi Q1: sistem baru
   membuat penjaga ganda-bayar eksplisit, **apa pun jawabannya**. ⚠️ Jawabannya tetap dibutuhkan
   untuk **memeriksa data lama**, bukan untuk merancang.
2. **Butir 6** *(isi rule peran)* — ⭐ **dikurung** sebagai **titik sambung**: spec menuliskan
   *"wewenang ditentukan oleh peran"* tanpa menyebut peran mana. ⛔ Menahan **pembangunan**, tidak
   menahan **penulisan**.
3. **Butir 1** *(kolom `T_QUOTATIONDATA`)* — `[data DBA]`, ⭐ menyentuh **bentuk tabel**, bukan
   perilaku.

⭐ **USULAN — urutan langkah berikutnya:**

**USULAN 1** ⭐ **Tulis spec Komite Claim Fac In sekarang**, dengan ketiga butir di atas tercatat
sebagai terbuka di dalamnya.

**USULAN 2** ⭐ **Bersamaan: ajukan kelima pertanyaan §7.2.** Tak satu pun menahan penulisan spec;
kelimanya menahan **pembangunan**.

**USULAN 3** ⭐ **Sesudah spec jadi: bentuk tabel dan relasinya untuk KEDUA modul Fac In
bersamaan** — sesuai `[keputusan work owner]` 2026-09-19 pada butir 23 spec Claim Fac In.

⛔ **Yang saya TIDAK usulkan: ronde 5.** ⚠️ Ronde 4 menutup enam butir dan membuka enam, ⛔ **tetapi
yang dibuka tak satu pun dapat dijawab dengan membaca korpus lagi** — semuanya menuntut **work
owner, DBA, atau uji di Pega**. ⭐ **Membaca lebih jauh sekarang menghasilkan kedalaman, bukan
kejelasan.**

---

## Lampiran — bukti berkas lain tidak disentuh

| Berkas | Baris | Keadaan |
| --- | ---: | --- |
| `komite-claim-facin\grilling-ronde-1.md` | 673 | ✅ utuh |
| `komite-claim-facin\grilling-ronde-2.md` | 630 | ✅ utuh |
| `komite-claim-facin\grilling-ronde-3.md` | 944 | ✅ utuh |
| `claim-facin\spec.md` | 1538 | ✅ utuh |
| `claim-prop\spec.md` | 2119 | ✅ utuh |
| `komite-claim-prop\spec.md` | 2034 | ✅ utuh |
| `OUTPUT_HASIL_RNM\CLAUDE.md` | 279 | ✅ utuh |
| `docs\adr\` | 15 berkas | ✅ utuh |
| korpus 114 · 482 · 80 | — | ✅ dibuka hanya untuk dibaca |

⛔ Kode Go/React **NOL** · `CREATE TABLE` **NOL** · DDL **NOL** · nomor baris XML **NOL** ·
nilai kredensial **NOL** · revisi ADR **NOL** · keputusan work owner baru di luar K1–K4 **NOL**.
