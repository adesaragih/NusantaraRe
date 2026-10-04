# Komite Claim Fac In — Grilling Ronde 3
## Tangga penyetuju, dua penyimpan keputusan, dan sebuah tuduhan palsu yang saya cegat sendiri

⭐ **Ronde ini MEMBACA.** Satu-satunya yang diputuskan adalah **dua keputusan work owner** di §7,
dan itu pun keputusan yang **sudah diambil**, bukan yang lahir di sini.

**Jendela sensus seluruh ronde ini:** 114 berkas XML di `Komite Claim FacIn`, dibandingkan bila
perlu dengan 482 berkas `Claim Fac In` dan 80 berkas `Komite Claim Prop`. ⛔ Tidak ada berkas
di luar ketiga korpus itu yang dibuka.

> **SENSUS BERKAS INI**
>
> ⛔ **Jendelanya disebut lebih dulu: seluruh berkas ini DIKURANGI blok sensus ini sendiri**
> — **927 baris**. ⭐ Blok ini dikecualikan dengan sengaja, sebab ia mengutip penanda sungguhan
> di dalam RALAT-nya, sehingga menghitung dirinya sendiri membuat angkanya berkejaran tanpa henti.
>
> `[terverifikasi]` **29** · `[keputusan work owner]` **2** · `[terbuka]` **9** ·
> `[penyimpangan sadar]` **1** · `[data DBA]` **3**
> ⚠️ **82** · ⭐ **172** · ⛔ **98** · ✅ **23** · ❓ **6** · RALAT **5** · tabel **35**
>
> ⭐ **Dihitung DUA CARA:** *(a)* cacah penanda atas badan berkas sebagai satu teks utuh;
> *(b)* cacah ulang **baris demi baris**. ✅ **Keduanya sepakat pada kedua belas angka.**
>
> ⛔ **RALAT terhadap kepala sensus yang sempat saya tulis SEBELUM berkasnya jadi.** Kalimatnya
> dikutip: *"terverifikasi **27** · keputusan work owner **2** · terbuka **13** · peringatan **58** ·
> bintang **71** · larangan **44** · centang **11** · RALAT **3** · tabel **16**"*.
> ⚠️ **Sembilan angka, delapan di antaranya salah.** ⭐ Ditulis mendahului berkasnya — persis
> kegagalan yang aturan sensus ada untuk mencegah. Dicatat, tidak dihapus.

---

## §1 — ⭐⭐⭐ `KomiteCount`: tangga penyetuju TERBACA PENUH

⛔ **Ronde 2 menutup bab ini dengan kalimat *"tangga penyetuju tidak terbaca sama sekali"*.**
⭐ **Kalimat itu sekarang batal.** Tangganya terbaca ujung ke ujung — penaiknya memang **bukan** di
`KomiteRouter`, dan di situlah ronde 2 tersesat.

### 1a — Cacah DUA CARA

**CARA-1 — sisiran teks mentah**, 114 · 482 · 80 berkas:

| | `KomiteCount` | `KomiteLoop` | `TotalKomite` |
| --- | ---: | ---: | ---: |
| Komite Claim FacIn | 69 kali / 9 berkas | 50 kali / 6 berkas | 7 kali / 2 berkas |
| Claim Fac In | 10 kali / 5 berkas | 10 kali / 5 berkas | 8 kali / 4 berkas |
| Komite Claim Prop | 64 kali / 7 berkas | 42 kali / 4 berkas | 21 kali / 1 berkas |

**CARA-2 — parser XML, hanya PENUGASAN** *(sasaran di kiri tanda sama dengan)*:

| | penugasan |
| --- | ---: |
| Komite Claim FacIn | **7** |
| Claim Fac In | **10** |
| Komite Claim Prop | **6** |

⭐ **Kedua cara sepakat pada hal yang penting:** `Claim Fac In` — **modul induk** — **menyentuh
`KomiteCount`**, dan ronde 2 tidak pernah melihat ke sana.

### 1b — ⭐⭐ Penyemainya ada di INDUK, bukan di modul ini

`[terverifikasi]` **Tiga rule di `Claim Fac In` menyemai pencacah, lalu membuat kasus anak komite:**

| Rule induk | Langkah | Yang disemai | Sesudahnya |
| --- | --- | --- | --- |
| ⭐ `CreateKMTNo_Act` | 9 | `childPageKomite.KomiteLoop = @Utilities.SizeOfPropertyList(childPageKomite.KomiteList)` · `childPageKomite.KomiteCount = 1` · `.TotalKomite = @Utilities.SizeOfPropertyList(…ObjectList…)` | lgk 10 `Call pxAddChildWork` |
| `SendCloseClaimToKomite` | 5.3 | `KomiteLoop = 1` · `KomiteCount = 1` | lgk 5.4 `Call pxAddChildWork` |
| `SendRejectClaimToKomite2` | 7.3 | `KomiteLoop = 1` · `KomiteCount = 1` | `pxAddChildWork` |

⭐ **Artinya tiga hal yang langsung dapat dibawa ke rancangan:**

1. **`KomiteLoop` = banyaknya anggota komite**, dihitung dari panjang daftar `KomiteList`.
2. **`KomiteCount` mulai dari 1**, bukan 0.
3. ⭐ **Jalur "tutup klaim" dan "tolak klaim" selalu berkomite SATU ORANG** — `KomiteLoop`
   ditanam `1`. ⚠️ Hanya jalur **Adjustment** yang bertangga banyak.

### 1c — Penaik pencacah di modul ini — **empat tempat, satu mati**

`[terverifikasi]` Jendela: seluruh 114 berkas, medan `PropertiesName` di dalam `pyParamArray`
pada setiap langkah, termasuk langkah bersarang.

| Rule | Lgk | Penugasan | Bendera gerbang | Keadaan |
| --- | --- | --- | --- | --- |
| ⭐ `KomitePost_Adjustment` | **24** | `KomiteCount = KomiteCount + 1` | ⚠️ **`false`** — gerbang `KomiteCount < KomiteLoop` **MATI** | **HIDUP, tanpa syarat** |
| ⭐ `KomitePost_Adjustment` | **14** | `KomiteCount = KomiteLoop` | **`true`** — gerbang `AcceptStatus=="2"` **HIDUP** | **HIDUP, bersyarat** |
| `KomitePost_CloseClaim` | 13 | `KomiteCount = KomiteCount +1` | *(tanpa gerbang sama sekali)* | **HIDUP, tanpa syarat** |
| `KomitePost_Reject` | 16 | `KomiteCount = KomiteCount +1` | *(tanpa gerbang sama sekali)* | **HIDUP, tanpa syarat** |
| ⛔ `KomiteRouter` | 5 | `.KomiteCount = .KomiteCount+1` | — | ⛔ **BER-REMARK, mati** |
| `ApprovalKomite_Act` | 8 | `KomiteLoop = @LengthOfPageList(pyWorkPage.KomiteList)` | *(tanpa gerbang)* | **HIDUP** — ⭐ dihitung **ulang** saat penyetujuan |
| `SetDataForInformation_Act` | 2 | `Local.TotalKomite = pyWorkCover.ClaimData.…` | *(tanpa gerbang)* | HIDUP, hanya lokal |

⚠️ **Dua hal di tabel itu wajib dibaca pelan.**

**Pertama — lgk 24.** ⭐ Benderanya `false`, dan menurut aturan yang sudah mapan itu **mematikan
GERBANG, bukan LANGKAH**. ⛔ Maka syarat `KomiteCount < KomiteLoop` **tidak pernah ditegakkan**, dan
pencacah **naik setiap kali rule ini berjalan** — termasuk ketika penyetuju terakhir sudah selesai.
⚠️ **`KomiteCount` karenanya dapat melampaui `KomiteLoop`.** Efek akhir yang bergantung pada
`KomiteCount==KomiteLoop` berjalan **lebih dulu** di langkah 7.x–23, jadi efeknya tetap terbit
sekali; ⛔ **tetapi nilai yang tersimpan sesudah itu bukan `KomiteLoop`, melainkan `KomiteLoop+1`.**

**Kedua — lgk 14.** ⭐ Ini **jalan pintas penolakan**: begitu seorang penyetuju menolak
(`AcceptStatus=="2"`), pencacah **dilompatkan langsung ke ujung**. ⚠️ Digabung dengan lgk 24 yang
tanpa syarat, sebuah penolakan meninggalkan `KomiteCount = KomiteLoop + 1`.

### 1d — ⭐⭐⭐ Pemutus putaran ternyata BUKAN pencacah

⛔ **Ini temuan terpenting ronde 3, dan ia mengubah bentuk modulnya.**

`[terverifikasi]` Alur `Komite_Flow.xml`:

| Unsur | Isi |
| --- | --- |
| bentuk `ASSIGNMENT63` | nama **`KomiteRouter`**, `pyImplementation = WorkList`, `pyRouteTo = Custom` |
| bentuk `Decision1` | nama **`KomiteLoop`**, `pyDecisionClass = Data-MO-Gateway-DataXOR`, `pyCompareOperator = "="` |
| penghubung `Transition2` | nama **`IsKomiteLoop`**, `pyConditionType = **When**`, `pyExpression = IsKomiteLoop`, bobot **100**, tujuan **`ASSIGNMENT63`** ⭐ *(balik ke router — inilah putarannya)* |
| penghubung `Transition3` | nama **`NoLoop`**, `pyConditionType = **Else**`, bobot **0**, tujuan **`END52`** |

`[terverifikasi]` Isi rule `IsKomiteLoop`, dibaca dari medan mentahnya:
`pyConditionString = "Status Penerimaan = 1"`, dan medan pendukungnya
`pyUnmodifiedPath = .AcceptStatus`, `pyLabel = [.AcceptStatus][=]["1"]`.
⭐ Jadi syarat berputar itu berbunyi: **`.AcceptStatus == "1"`**.

⛔⛔ **Maka putaran persetujuan berputar selama keputusan terakhir adalah TERIMA, dan berhenti
hanya ketika keputusan terakhir adalah TOLAK.** ⚠️ **`KomiteCount` sama sekali tidak ikut
memutuskan kapan putaran berhenti.** Ia hanya menjadi **penunjuk giliran** dan **penjaga efek
akhir** *(`KomiteCount==KomiteLoop`)*.

⚠️⚠️ **Dan di sinilah lubang yang tidak dapat saya tutup dengan membaca saja.** Dua-duanya
ber-remark:

| Rule | Lgk | Penugasan | Keadaan |
| --- | --- | --- | --- |
| `KomitePostAct` | 6 | `.AcceptStatus = ""` | ⛔ **BER-REMARK** |
| `KomiteRouter` | 5 | `.AcceptStatus = ""` | ⛔ **BER-REMARK** |

⛔ **Tidak ada satu pun langkah hidup di 114 berkas yang mengosongkan `.AcceptStatus`.**
⚠️ Dan penutup paksa kasus pada jalur Adjustment — `KomitePost_Adjustment` lgk **28**
`Call ASMForceCaseClose`, bergerbang `KomiteCount==KomiteLoop || AcceptStatus=="2"` — ⛔ **juga
BER-REMARK.** *(Jalur `CloseClaim` lgk 14 dan `Reject` lgk 17 memang memanggil `pxForceCaseClose`
dan keduanya hidup — tetapi itu jalur berkomite satu orang.)*

⛔ **Saya BERHENTI di sini dan tidak menyimpulkan bahwa ini cacat.** Yang dapat saya nyatakan:
*(a)* syarat berputar adalah `.AcceptStatus=="1"`; *(b)* tidak ada penulis hidup yang
mengosongkannya; *(c)* penutup paksa jalur Adjustment ber-remark. ⚠️ **Bagaimana putaran jalur
Adjustment berakhir pada keadaan "semua menyetujui" TIDAK TERBACA dari korpus ini.**
⭐ **`[terbuka]` BARU butir 19 — MEMBLOKIR.**

### 1e — Diadu dengan Komite Claim Prop

`[terverifikasi]` **Bentuk tangganya IDENTIK** — sampai ke nomor langkah yang sejajar:

| | Komite Claim FacIn | Komite Claim Prop |
| --- | --- | --- |
| jalan pintas tolak | `KomitePost_Adjustment` lgk 14, bendera `true` | `KomitePostAdjustment` lgk 25, bendera `true` |
| penaik tanpa syarat | lgk 24, bendera `false` | lgk 40, bendera `false` |
| penaik jalur tutup | `KomitePost_CloseClaim` lgk 13 | `KomitePost_Close` lgk 16 |
| penaik jalur tolak | `KomitePost_Reject` lgk 16 | `KomitePost_Reject` lgk 15 |
| penaik di router | lgk 5 ⛔ ber-remark | lgk 5 ⛔ ber-remark |
| bentuk keputusan alur | `Decision1` / `KomiteLoop` / `DataXOR` | **sama persis** |

⛔ **TETAPI identitas empat bagiannya BERBEDA — tidak satu pun rule penaik yang dibagi:**

| Rule | `pxInsName` | Kelas | Ruleset | Versi |
| --- | --- | --- | --- | --- |
| FacIn | `…WORK-KOMITE!KOMITEPOST_ADJUSTMENT` | `…Work-Komite` | GCNMFW | 01-01-27 |
| Prop | `…WORK-KOMITETREATY!KOMITEPOSTADJUSTMENT` | `…Work-KomiteTreaty` | GCNMFW | 01-01-28 |
| FacIn | `…WORK-KOMITE!KOMITEROUTER` | `…Work-Komite` | GCNMFW | 01-01-15 |
| Prop | `…WORK-KOMITETREATY!KOMITEROUTER` | `…Work-KomiteTreaty` | GCNMFW | 01-01-15 |

⭐ **Ini kabar BAIK, dan ia mengoreksi nada ronde 2.** Ronde 2 menutup §4a dengan *"daur hidup
SAMA"* dan memakai itu sebagai alasan **memperlakukan kedua modul serupa**. ⭐ Yang benar:
**rancangannya satu, rule-nya dua.** ⚠️ Maka mengubah tangga Fac In **tidak** menyentuh Prop, dan
sebaliknya — ⛔ **tetapi juga berarti setiap perbaikan harus dikerjakan DUA KALI**, dan sejarah
korpus menunjukkan keduanya sudah **melenceng** *(01-01-27 lawan 01-01-28)*.

---

## §2 — Dua rule penyimpan keputusan

### 2a — Bentuk luarnya

`[terverifikasi]` Keduanya berkelas **`ASM-FW-GCNMFW-Data-ObjectItem`**, ruleset **GCNMFW 01-01-22**,
dan ⭐ **keduanya berasal dari salinan produksi `pegaprdnusare`** — bukan dari server pengembangan.

| | `SaveAccept_ACT` | `SaveReject_ACT_KMT` |
| --- | ---: | ---: |
| langkah | **51** | **20** |
| ber-remark | **2** | **0** |
| bendera `true` | 33 | *(tidak dicacah terpisah)* |
| bendera kosong | 18 | — |
| parameter rule | `ObjectIndex` · `ObjectItem` · `Adjustment` | `Object` · `ObjectItem` · `Estimasi` |

⚠️ **Parameter keduanya BERBEDA nama** meski berkelas sama — `ObjectIndex` lawan `Object`,
`Adjustment` lawan `Estimasi`. ⛔ **Bukan sepasang rule kembar.**

### 2b — ⭐⭐ Yang disimpan, dan ke mana

`[terverifikasi]` **`SaveAccept_ACT`** membangun **ringkasan akseptasi** `osAkseptasi` di dalam
halaman kerja sementara, bercabang per **lini usaha**:

| Lgk | Gerbang | Yang dikerjakan |
| --- | --- | --- |
| 3 | `IsFire` ⭐ kode 5 · `IsAneka` ⭐ kode 5 · `isGolfInsurance` | cabang api/aneka/golf |
| 4 | `IsMarineCargo` | cabang marine |
| 5 | `IsMBU` | cabang MBU |
| 6 | `IsTravel` | cabang travel |
| 7 | `IsPA` | cabang PA |

⭐ **Dua gerbang di lgk 3 memakai kode arah 5** — ⚠️ **rantai DAN berubah menjadi ATAU**. Jadi
cabang itu berbunyi `IsFire OR IsAneka OR isGolfInsurance`, **bukan** ketiganya sekaligus.

**Halaman yang disentuh:** `TempOpenPage` · `TempOSAkseptasi` · `InputData` · `Local`.
⛔ **NOL penugasan ke `Primary.*`, `pyWorkCover.*`, maupun `pyWorkPage.*`.**

⛔⛔ **DAN INILAH TEMUANNYA:** kedua langkah yang akan **menyimpannya ke basis data** —

| Lgk | Metode | Rule SQL | Keadaan |
| --- | --- | --- | --- |
| 10 | `RDB-List` | `SaveOSClaim_SQL` *(ber-`COMMIT;`)* | ⛔ **BER-REMARK** |
| 11 | `RDB-List` | `SaveOSClaim_SQL` *(ber-`COMMIT;`)* | ⛔ **BER-REMARK** |

⚠️ **`SaveAccept_ACT` menyusun ringkasan akseptasi seluruhnya di memori dan TIDAK PERNAH
menyimpannya.** ⛔ Lima puluh satu langkah kerja yang berakhir tanpa tulisan.

`[terverifikasi]` **`SaveReject_ACT_KMT`** berkebalikan:

| Lgk | Yang dikerjakan | Keadaan |
| --- | --- | --- |
| 1 | `pyWorkPage.Quotation = TempOpenPage.Quotation` ⭐ salinan **satu halaman utuh** | HIDUP |
| 2–15 | cabang lini usaha yang sama bentuknya | HIDUP |
| 19 | `InputData.CARI1 = @ASM.GetPageJSONString()` | HIDUP |
| ⭐ **20** | `RDB-List` → **`SaveOSClaim_SQL`** *(ber-`COMMIT;`)* | ⭐ **HIDUP, tanpa bendera, tanpa gerbang** |

⛔ **Jadi PENOLAKAN tersimpan permanen; PENERIMAAN tidak.** ⚠️ Asimetri itu **tidak dapat saya
jelaskan dari korpus**, dan ia menyentuh catatan akseptasi — bahan laporan. ⭐ **`[terbuka]` BARU
butir 22.**

⚠️ **Catatan yang wajib menyertainya:** ronde 2 sudah membuktikan modul ini punya **NOL jalur
kegagalan pada 25 langkah `Obj-*`**. ⛔ Langkah 20 di atas membawa `COMMIT;` **tanpa gerbang dan
tanpa jalur gagal** — bila ia gagal, tidak ada apa pun yang menangkapnya.

### 2c — Wewenang di dalamnya: ⛔ **NOL**

`[terverifikasi]` Jendela: **seluruh** `pyStepsPreCondParamsWhen` dan `pyStepsTransParamsWhen`
pada kedua rule, disaring dengan pola `pyPosition|KomiteID|USERNAME|pyUserIdentifier|OPERATOR_ID|
pyUserName|Proteksi|CARI1|KomiteAproval|Komite`.

> **Hasil: 0 gerbang berwewenang di `SaveAccept_ACT`. 0 di `SaveReject_ACT_KMT`.**

⛔ **Kedua rule yang menyimpan keputusan komite TIDAK memeriksa siapa yang memanggilnya.**
⭐ Itu **menguatkan**, bukan melemahkan, vonis ronde 2: penegakan wewenang memang **tidak ada di
lapisan penyimpan**.

⚠️ **Satu catatan tentang `InputData.CARI1` di lgk 9/19:** namanya mirip `ProteksiKomite.CARI1`,
⛔ **tetapi ia sama sekali bukan wewenang** — isinya `@ASM.GetPageJSONString()`, yaitu **muatan JSON**
yang hendak dikirim ke SQL. ⭐ Persis jebakan penampung generik yang §3b di bawah ini uji.

### 2d — Siapa pemanggilnya

`[terverifikasi]`

| Rule | Pemanggil | Lgk | Keadaan |
| --- | --- | --- | --- |
| `SaveAccept_ACT` | `KomitePost_Adjustment` **(modul ini)** | 7.2.1.17 | HIDUP, bendera `true` |
| `SaveAccept_ACT` | `SaveAcceptation` *(Claim Fac In)* | 9.1 | ⛔ **BER-REMARK** |
| `SaveReject_ACT_KMT` | `KomitePost_Reject` **(modul ini)** | 12.2.2.2 | HIDUP, bendera `true` |
| `SaveAcceptation_KMT` | `KomitePost_Adjustment` | 7.2.1.16 | HIDUP, bendera `true` |

⭐ **Keduanya dipanggil HANYA dari dalam modul komite.** Satu-satunya pemanggil dari induk
ber-remark. ⛔ Jadi meski kelasnya `Data-ObjectItem` — kelas yang dipakai bersama — **pemakaiannya
khas komite.**

⚠️ Gerbang pemanggil `SaveAccept_ACT` di lgk 7.2.1.17 berbunyi tiga syarat berantai:
`AcceptStatus=="1"` **DAN** `KomiteCount==KomiteLoop` **DAN** `.AcceptanceStatus=="1"`.
⭐ **Artinya ringkasan akseptasi hanya disusun oleh penyetuju TERAKHIR.**

---

## §3 — ⭐⭐⭐ `ProteksiKomite` disisir habis — vonis ADR-0014 **BERTAHAN**

⚠️ Ronde 2 menyebut bab ini **kesimpulannya yang paling rawan**, sebab vonisnya bersandar pada
**satu pembacaan**. ⭐ Sekarang ia bersandar pada **sisiran menyeluruh**.

### 3a — Setiap tag, 114 berkas

`[terverifikasi]` Jendela: **setiap elemen berisi teks** pada seluruh 114 berkas, dicocokkan dengan
kata `ProteksiKomite` — **tanpa** menyaring nama tag lebih dulu.

> **Hasil: 2 berkas · 5 kemunculan · 4 nama tag.**

| Berkas | Tag | Isi |
| --- | --- | --- |
| `Activity\SetProteksiSubmiteKomite` | `PropertiesName` | `ProteksiKomite.CARI1` |
| `Activity\SetProteksiSubmiteKomite` | `PropertiesName` | `ProteksiKomite.CARI1` |
| `Activity\SetProteksiSubmiteKomite` | `pyPagesAndClassesPage` | `ProteksiKomite` |
| `Activity\SetProteksiSubmiteKomite` | `pxPageName` | `ProteksiKomite` |
| ⭐ `Section\ShowTransfer` | `pyDisabledWhen` | `ProteksiKomite.CARI1==0 \|\| pyWorkPage.Adjustment.AcceptedNo != ''` |

⭐ **Satu penulis. Satu pembaca. Tidak ada yang ketiga.** ⛔ Tidak ada alias halaman, tidak ada
parameter yang membawanya, tidak ada rujukan dari `Flow`, `FlowAction`, `DataTransform`, maupun
`When`.

### 3b — ⚠️ Jebakan `CARI<n>` diuji, dan ia memang jebakan

`[terverifikasi]` Kata utuh `CARI1` — **tanpa awalan halaman** — muncul di **19 berkas**. Yang
dipakainya untuk hal yang **sama sekali berbeda**:

| Berkas | Halaman pembawa | Untuk apa |
| --- | --- | --- |
| `HitServiceToKasirKMT_Act` | `TempCeding` · `TempKasir` · `Tanggal` | muatan kirim ke kasir |
| `InsertLogServiceClaim` | `ParamLog` | muatan log |
| `KomitePost_Adjustment` | `ParamSeq` · `InsertHistory` | ⭐ parameter pembangkit nomor |
| `KomitePost_CloseClaim` · `SaveAccept_ACT` · `SaveReject_ACT_KMT` | `InputData` | muatan JSON |
| `KomitePost_Reject` | `DataIn` | muatan JSON |
| `getStatusKonversi_Act` | `ParamCekKonversi` | cek konversi |
| `GetEmailCeding_SQL` | `TempCeding` | parameter SQL surel |
| ⭐ `SetProteksiSubmiteKomite` | **`ProteksiKomite`** | ⭐ **satu-satunya yang berarti wewenang** |

⛔ **`CARI1` sendirian TIDAK BERARTI APA-APA.** ⭐ Yang memberinya arti adalah **halaman
pembawanya**. ⚠️ Seandainya sisiran dilakukan atas `CARI1` saja, hasilnya akan **19 berkas** dan
vonis ADR-0014 akan runtuh karena alasan yang **salah**.

### 3c — ⭐⭐ Pemeriksaan giliran KEDUA yang ronde 2 lewatkan

`[terverifikasi]` Jendela: **setiap elemen berisi teks**, 114 berkas, disaring pola
`OperatorID.|pyUserIdentifier|pyUserName|Local.USERNAME|OPERATOR_ID|.pyPosition|pxRequestor`.
**Hasil: 12 berkas · 84 kemunculan.** Yang **membandingkan akun dengan pemilik giliran**: **dua**.

| # | Rule | Lgk | Ungkapan | Keadaan |
| --- | --- | --- | --- | --- |
| **1** | `SetProteksiSubmiteKomite` | — | `Local.USERNAME == .KomiteID` | ⭐ sudah dikenal ronde 2 |
| ⭐ **2** | **`ApprovalKomite_Act`** | **6.2** | **`pyWorkPage.KomiteList(1).KomiteID == .OPERATOR_ID`** → metode **`Property-Remove`** | ⭐ **BARU** |

⭐⭐ **Yang kedua itu adalah aturan "tidak boleh menyetujui diri sendiri":** bila anggota komite
**pertama** ternyata orang yang sedang membuka kasus, ⭐ **ia DIKELUARKAN dari daftar** — dan lgk 8
langsung **menghitung ulang `KomiteLoop`** dari daftar yang sudah pendek itu.

⚠️ **Tetapi indeksnya DITANAM `1`.** ⛔ Hanya anggota **pertama** yang diperiksa. Bila orang yang
sama berada di urutan kedua atau ketiga, **ia tetap di daftar dan tetap boleh menyetujui klaim yang
ia ajukan sendiri.** ⭐ **`[terbuka]` BARU butir 24.**

### 3d — ⚠️ Nama akun pegawai yang **ditanam keras**

`[terverifikasi]` Jendela: pola `pyUserIdentifier == "<teks>"` pada seluruh berkas ketiga modul.

| Modul | Akun ditanam (unik) | Berkas |
| --- | ---: | ---: |
| Komite Claim FacIn | **5** | 4 |
| Komite Claim Prop | **4** | 2 |
| Claim Fac In | **0** | 0 |

Berkasnya: `KomitePost_Adjustment` · `KomitePost_CloseClaim` · `KomitePost_Reject` *(masing-masing 4
akun)* dan `SetProteksiSubmiteKomite` *(2 akun)*.

⛔ **Saya sengaja TIDAK menuliskan daftar orangnya di sini** — yang dibutuhkan rancangan adalah
**bentuk dan akibatnya**, bukan rosternya.

**Bentuknya** — pada `SetProteksiSubmiteKomite`, medan `Local.USERNAME` **tidak** berisi akun
masuk apa adanya, melainkan hasil sebuah ungkapan `@if` sepanjang 146 aksara dengan **2
pembandingan** yang **memetakan satu akun menjadi nama lain**, dan mengembalikan akun asli untuk
selebihnya. Pada tiga rule `KomitePost_*`, pola yang sama dipakai untuk menyusun teks
`"Accepted by …"` / `"Rejected by …"`.

⛔⛔ **Akibatnya berat, dan ia menyentuh §3c nomor 1 secara langsung:** yang dibandingkan dengan
`.KomiteID` **bukan identitas akun**, melainkan **alias hasil pemetaan yang ditanam di dalam kode**.
⚠️ Bila seorang pegawai berganti akun, atau alias itu dipakai orang lain, **pemeriksaan pemilik
giliran satu-satunya yang dipunyai modul ini ikut bergeser** — tanpa ada yang mengubah aturan apa
pun. ⭐ **`[terbuka]` BARU butir 21.**

### 3e — ⭐ Vonis ulang ADR-0014

> ⭐ **BERTAHAN pada MENENTANG SEBAGIAN** — ⚠️ **dengan dua penguatan dan satu pelemahan.**

**Menguatkan vonis lama:**
1. ⭐ Sisiran menyeluruh membuktikan pembaca `ProteksiKomite` memang **hanya satu, dan ia di layar**
   *(§3a)*.
2. ⭐ Kedua rule **penyimpan keputusan** terbukti punya **NOL gerbang wewenang** *(§2c)* — jadi
   lapisan layanan memang **tidak menegakkan apa pun**.

**Melemahkan vonis lama:**
3. ⚠️ Ternyata **ada pemeriksaan giliran kedua** di `ApprovalKomite_Act` *(§3c)* — ⛔ ronde 2 tidak
   melihatnya, dan ronde 2 menulis seolah hanya ada satu. ⭐ **Itu RALAT**, dicatat di §8.4.

⛔ **Yang tidak berubah:** kedua pemeriksaan itu **berdiri di lapisan tampilan**, keduanya dapat
**dilewati** dengan memanggil rule penyimpan secara langsung, dan yang pertama **dapat dimatikan
seluruhnya** oleh gerbang berkode arah 5. ⭐ Maka ADR-0014 — *penegakan wewenang milik lapisan
layanan* — tetap **ditentang sebagian** oleh modul ini, dan tetap wajib **ditegakkan ulang** di
sistem baru.

---

## §4 — Lima `Section`, lima `FlowAction`, dan uji instrumen

### 4a — Daftar kesepuluhnya

`[terverifikasi]`

| Jenis | Nama | Kelas | Gerbang |
| --- | --- | --- | --- |
| Section | `DetailAdjustmentFac` | `ASM-FW-GCNMFW-Data-Adjustment` | — |
| Section | ⚠️ `ShowRetro_Sec` | **`ASM-FW-GISFW-Data-SpreadingRisk`** | — |
| Section | ⚠️ `ShowSecurityReinsurer` | **`ASM-FW-GISFW-Data-FacOffer`** | — |
| ⭐ Section | **`ShowTransfer`** | **`ASM-FW-GCNMFW-Work-Komite`** | ⭐ **3 × `pyDisabledWhen`** |
| Section | `SpreadingDetail` | `ASM-FW-GCNMFW-Data-ObjectItem` | — |
| FlowAction | `DetailAdjustmentFac` | `ASM-FW-GCNMFW-Data-Adjustment` | — |
| FlowAction | ⚠️ `ShowRetro` | **`ASM-FW-GISFW-Data-SpreadingRisk`** | — |
| FlowAction | ⚠️ `ShowSecurityReinsurer` | **`ASM-FW-GISFW-Data-FacOffer`** | — |
| FlowAction | `SpreadingDetail` | `ASM-FW-GCNMFW-Data-ObjectItem` | — |
| ⭐ FlowAction | **`ViewTransferDtl`** | **`ASM-FW-GCNMFW-Work-Komite`** | — |

⭐ **Hanya DUA dari sepuluh yang benar-benar milik kelas kerja komite** — `ShowTransfer` dan
`ViewTransferDtl`. ⚠️ Delapan sisanya **dipinjam** dari kelas lain, dan **dua di antaranya dari
kerangka `GISFW`, bukan `GCNMFW`** — ⛔ kerangka yang berbeda. `[terbuka]` ringan: apakah kedua
layar `GISFW` itu memang dipakai dari konteks komite, atau ikut terbawa ekspor.

⛔ **Delapan berkas itu NOL gerbang** — tidak ada `pyVisibleWhen`, `pyDisabledWhen`,
`pyReadOnlyWhen`, maupun `pyRequiredWhen`. ⭐ Layar baca-saja.

### 4b — ⭐⭐ `ShowTransfer` dibaca dalam — **temuan baru**

`[terverifikasi]` Ruleset **`ADESAMUEL@` 01-01-01**, sistem `pegadevnusare2`. ⭐ Inilah berkas
tunggal beruleset pribadi itu — sudah dikenal ronde 2 *(butir 2)*.

**Ketiga gerbangnya:**

| # | `pyDisabledWhen` |
| --- | --- |
| ⭐ **1** | **`.KomiteCount!='1'`** |
| ⭐ **2** | **`.KomiteCount!='1'`** |
| 3 | `ProteksiKomite.CARI1==0 \|\| pyWorkPage.Adjustment.AcceptedNo != ''` |

⛔⛔ **Dua gerbang pertama BARU — ronde 2 hanya melihat yang ketiga.**

⭐ **Artinya sebuah aturan bisnis yang belum pernah tercatat:** **medan yang dijaga kedua gerbang
itu hanya dapat disunting oleh penyetuju PERTAMA** *(`KomiteCount == 1`)*. Penyetuju kedua dan
seterusnya **melihatnya terkunci**.

⚠️ **Dan ia bertaut langsung dengan §1c.** Karena lgk 24 menaikkan pencacah **tanpa syarat**,
`KomiteCount` **tidak pernah kembali ke 1**. ⭐ Maka penguncian itu **searah dan permanen** dalam
satu daur hidup kasus — sesuatu yang harus ditiru dengan sadar, bukan disimpulkan ulang dari kode
Go.

⚠️ Pembandingnya **`'1'` bertanda kutip** — ⛔ **teks, bukan bilangan**, padahal `KomiteCount`
diisi bilangan `1` oleh `CreateKMTNo_Act`. ⭐ Pega memaafkan itu; **Go tidak akan**. Catat sebagai
jebakan alih.

⛔ Medan terikat properti **tidak terbaca** dari berkas ini dengan penyaring `pyReference` /
`pyPropertyName` / `pyFieldName`: **hasilnya nol**. ⚠️ **Tidak ketemu di medan `pyReference`,
`pyPropertyName`, `pyFieldName`, `pyFieldType`** — jadi **daftar medan yang dapat disunting belum
punya data**, dan saya tidak menebaknya. ⭐ `[terbuka]` BARU butir 23.

### 4c — ⭐ UJI INSTRUMEN penyaring gerbang mati *(butir 10)*

⛔ Ronde 1 memakai penyaring ini **tanpa pernah mengujinya**. Diuji sekarang, dengan kasus yang
jawabannya sudah diketahui lebih dulu:

| Langkah | Bendera | Instrumen membaca | Harapan | Hasil |
| --- | --- | --- | --- | --- |
| `KomitePost_Adjustment` lgk 24 | `false` | **MATI** | MATI | ✅ **LULUS** |
| `KomitePost_Adjustment` lgk 14 | `true` | **HIDUP** | HIDUP | ✅ **LULUS** |
| `KomitePost_Adjustment` lgk 13 | `false` | **MATI** | MATI | ✅ **LULUS** |
| `KomitePost_Adjustment` lgk 15 | `true` | **HIDUP** | HIDUP | ✅ **LULUS** |

**Uji kelima — pembacaan remark:** `KomitePost_Adjustment` melaporkan **11 langkah ber-remark**
pada nomor `7.2.1.2.1` · `.2` · `.3` · `.8` · `.9` · `.10` · `.11` · `.12` · `.13` · `.14` · `28`
— ⭐ **termasuk langkah BERSARANG lima tingkat**, yang membuktikan penelusur langkah tidak
kehilangan anak.

> ✅ **Instrumen LULUS 4 dari 4. Butir 10 DITUTUP.**

### 4d — Aksi lokal bernama

⛔ **NOL** pada seluruh sepuluh berkas — tidak ada `pyLocalActionName` berisi.
⚠️ **Tidak ketemu di medan `pyLocalActionName`, `pyActivityName`, `pyFlowName`.**

⭐ **Bandingkan:** pada Komite Claim Prop, bentuk `ASSIGNMENT63` membawa
`pyLocalActionsString = (3)` — **tiga aksi lokal**. ⛔ Pada modul ini medan itu **tidak ada sama
sekali**. ⚠️ Jadi kedua modul **tidak sama persis** di titik ini — RALAT ringan terhadap §4a ronde 2.

---

## §5 — ⭐ Empat puluh tujuh rule `When` bersama, diadu

### 5a — Angkanya, dua cara

`[terverifikasi]` Jendela: folder `When` ketiga modul.

| | CARA-1 *(nama berkas)* | CARA-2 *(identitas empat bagian)* |
| --- | ---: | ---: |
| Komite Claim FacIn | **49** | **49** |
| Claim Fac In | 60 | 60 |
| Komite Claim Prop | 6 | 6 |

⭐ **Kedua cara sepakat** — tidak ada dua berkas beridentitas sama.

| Irisan | Jumlah |
| --- | ---: |
| ⭐ identik empat bagian dengan **Claim Fac In** | ⭐ **47 dari 49** |
| identik empat bagian dengan **Komite Claim Prop** | **5 dari 49** |
| ⭐ **khas modul ini** | ⭐ **2** — `IsCustomBonds` · `IsKomiteLoop` |

> ✅ **Angka 47 dari brief TERVERIFIKASI ULANG dan BENAR.**

⚠️ Catatan: `IsKomiteLoop` Komite Claim Prop **ada**, tetapi berkelas `…Work-KomiteTreaty` ruleset
`01-01-07` — ⛔ **bukan identitas yang sama**, meski `pyConditionString`-nya **persis sama**:
`"Status Penerimaan = 1"`. ⭐ Sekali lagi: **rancangan satu, rule dua.**

### 5b — ⭐⭐ Apa yang sebenarnya diuji 49 rule itu

`[terverifikasi]` Jendela: **setiap** medan bersyarat pada tiap berkas `When` —
`pyLabel` · `pyCondition` · `pyConditionString` · `pyExpression` · `pyLeftTerm` · `pyRightTerm` ·
`pyLeftValue` · `pyRightValue` · `pyWhenExpression` · `pyFilterValue` · `pyValue`.
⛔ **Penyaring TIDAK mewajibkan awalan halaman** — itu pelajaran ronde 1.

> **Ungkapan terbaca: 49 dari 49. Kosong: 0.**

| Awalan yang diuji | Berapa `When` |
| --- | ---: |
| ⭐ `pyWorkPage` | **41** |
| ⭐ `Quotation` | **40** |
| `OfferFacIn` | 3 |
| `QuotationData` | 3 |
| `pxProcess` | 1 |

⭐⭐ **Seluruh pustaka `When` modul ini sebenarnya satu benda: PENGGOLONG LINI USAHA**, yang hampir
selalu membaca `pyWorkPage.Quotation.*`.

**Kelas tempat mereka didefinisikan:**

| Kelas | Jumlah |
| --- | ---: |
| `@baseclass` | 26 |
| `ASM-FW-GISFW-Data` | 13 |
| `ASM-FW-GISFW-Work` | 6 |
| `ASM-FW-GCNMFW-Work` | 2 |
| ⭐ `ASM-FW-GCNMFW-Work-Komite` | ⭐ **1** — hanya `IsKomiteLoop` |
| `Data-Party-Person` | 1 |

⛔ **Empat puluh delapan dari 49 didefinisikan di kelas yang BUKAN kelas kerja komite.**

### 5c — ⭐⭐⭐ Inti pertanyaannya: apakah medannya TERISI di `Work-Komite`

`[terverifikasi]` Yang mengisi `pyWorkPage.Quotation` pada objek kerja komite:

| Rule | Lgk | Yang disalin |
| --- | --- | --- |
| ⭐ `SetValueKomite` | 12 | ⭐ **HANYA DUA MEDAN** — `Quotation.BusinessType` dan `Quotation.GroupPanel`, dari `pyWorkCover.OfferFacIn.Quotation…` |
| `KomitePost_Reject` | 2 | `pyWorkPage.Quotation = TempOpenPage.Quotation` ⭐ **halaman utuh** |
| `SaveReject_ACT_KMT` | 1 | `pyWorkPage.Quotation = TempOpenPage.Quotation` ⭐ **halaman utuh** |

⚠️⚠️ **Salinan utuh HANYA terjadi di jalur penolakan.** ⛔ Di jalur penerimaan, `Quotation` pada
objek kerja komite **hanya punya dua medan**.

**Sekarang adu dengan apa yang diuji:**

| Medan `Quotation.*` yang diuji | Berapa `When` | Tersedia di `Work-Komite`? |
| --- | ---: | --- |
| ⭐ `BusinessType` | **38** | ⭐ **YA** — disalin `SetValueKomite` |
| ⚠️ `BusinessCode` | **8** | ⛔ **TIDAK DISALIN** |
| ⚠️ `BusinessName` | **1** | ⛔ **TIDAK DISALIN** |

| Golongan | Jumlah | Daftar |
| --- | ---: | --- |
| ✅ **AMAN** — hanya menguji medan yang disalin | **32** | `IsAllRisk` `IsAneka` `IsAviationHull` `IsBoiler` `IsBondingAndCustomBonds` `IsBondingKBG` `IsBurglary` `IsCIS` `IsCIT` `IsContractorsPlantMachinery` `IsCustomBond` `IsCustomBonds` `IsEar` `IsFidelity` `IsFire` `IsFireStyle1` `IsFireStyle2` `IsGlass` `IsGrowingTrees` `IsHE` `IsKPR` `IsLandRig` `IsMBD` `IsMarineCargo` `IsMarineHull` `IsOilGas` `IsPA` `IsWorkmenCompensation` `isBillboardNeonSyariah` `isElectronicEquipment` `isGolfInsurance` `isMaintenance` |
| ⚠️ **RAGU** — menguji medan yang **tidak** disalin | **9** | `IsBillboardNeon`→`BusinessName` · `IsCrime` `IsEnvironmental` `IsExclusion` `IsMBU` `IsProductsLiability` `IsProfessionalLiability` `IsTravel` `IsYieldShortfall` → semuanya `BusinessCode` |
| — tidak menguji `Quotation` | **8** | `IsCAR` `IsCLM` `IsCLMNP` `IsCLMP` `IsKomiteLoop` `IsLiability` `IsPEGAPROD` `IsPEGASyariah` |

`32 + 9 + 8 = 49` ✅

### 5d — ⭐⭐ Vonis: berapa yang HIDUP, dan bahaya yang tersembunyi

`[terverifikasi]` Jendela: **seluruh 114 berkas di luar folder `When`**, tag
`pyStepsPreCondParamsWhen` · `pyStepsTransParamsWhen` · `pyDisabledWhen` · `pyVisibleWhen` ·
`pyWhenName` · `pyConditionString` · `pyReadOnlyWhen` · `pyRequiredWhen` · ⭐ **`pyExpression`** ·
⭐ **`pyPreviousRule`** · `pyActionWhen` · `pyWhen` · `pyValidateWhen`.

> ⭐ **Dipakai: 13 dari 49.** ⭐ **Dipakai dan HIDUP: 13 dari 49.**
> **Tidak dipakai sama sekali: 36.**

| `When` | Dipakai | Hidup | Di mana |
| --- | ---: | ---: | --- |
| `IsPEGAPROD` | 11 | 10 | penjaga "hanya di produksi" — surel, Google Storage, kasir |
| `IsCLMNP` | 8 | 8 | `HitServiceToKasirKMT_Act` · `SendEmailKlaim_KMT` |
| `IsAneka` | 6 | 5 | `SaveAccept_ACT` · `SaveReject_ACT_KMT` · cetak PDF |
| `IsFire` | 6 | 5 | idem |
| ⚠️ **`IsMBU`** | 6 | **5** | ⚠️ idem — **RAGU** |
| `IsMarineCargo` | 6 | 5 | idem |
| `IsPA` | 6 | 5 | idem |
| ⚠️ **`IsTravel`** | 6 | **5** | ⚠️ idem — **RAGU** |
| `isGolfInsurance` | 5 | 4 | idem |
| `IsCLM` | 3 | 3 | kasir · surel |
| `IsCLMP` | 3 | 3 | kasir · surel |
| `IsPEGASyariah` | 3 | 3 | `HitServiceToKasirKMT_Act` |
| ⭐ `IsKomiteLoop` | 2 | 2 | ⭐ **`Flow\Komite_Flow`** — `pyExpression` + `pyPreviousRule` |

⛔⛔ **DAN INILAH BAHAYANYA:** **`IsMBU` dan `IsTravel` ADA di golongan RAGU, dan keduanya
DIPAKAI HIDUP lima kali** — di `SaveAccept_ACT` lgk 5 dan 6, `SaveReject_ACT_KMT` lgk 6/13 dan 7/14,
`SetDataForInformation_Act` lgk 10 dan 11, serta `PrintPDFAccep_MultiAksep_KMT` lgk 19 dan 21.

⚠️ Keduanya menguji `pyWorkPage.Quotation.BusinessCode`, yang **`SetValueKomite` tidak salin**.
⛔ **Maka pada objek kerja komite keduanya membaca medan kosong dan bernilai SALAH — selalu,
diam-diam, tanpa pesan galat.** ⭐ Akibatnya **cabang MBU dan cabang Travel pada penyusunan
ringkasan akseptasi tidak pernah terbit** dari konteks komite.

⚠️ **Satu peringatan atas kesimpulan saya sendiri:** ini **belum saya buktikan dengan menjalankan
apa pun**. Yang terbukti adalah: *(a)* kedua `When` menguji `BusinessCode`; *(b)* `SetValueKomite`
tidak menyalinnya; *(c)* salinan halaman utuh hanya ada di jalur tolak. ⛔ **Mungkin saja ada
penulis `BusinessCode` yang belum saya temukan** — saya menyisir `PropertiesName`, `pyParamArray`,
dan `Page-Copy`/`Property-Set-Page` di kedua modul dan **tidak menemukannya**. ⭐ **`[terbuka]` BARU
butir 20 — MEMBLOKIR.**

⭐ **Sisi baiknya: 36 dari 49 rule `When` adalah BEBAN MATI di modul ini.** Mereka ikut terbawa
ekspor karena kelasnya dibagi, bukan karena modul komite memakainya. ⛔ **Jangan dialihkan.**

---

## §6 — Pembangkit nomor akseptasi

### 6a — Cacah menyeluruh

`[terverifikasi]` Jendela: seluruh berkas ketiga modul, medan `pyStepsActivityName` dan
`RequestType` pada setiap langkah.

| Pembangkit | Komite Claim FacIn | Claim Fac In | Komite Claim Prop |
| --- | --- | --- | --- |
| `GenerateNoAccept` | 5 pemanggilan · **1 ber-remark** | **0** | **0** |
| `GenerateNoAcceptNonFire` | 10 pemanggilan · **6 ber-remark** | **0** | **0** |
| ⭐ `GetSequenceNumber_SQL` | 5 pemanggilan · ⭐ **0 ber-remark** | 4 · 0 ber-remark | 2 · 0 ber-remark |

⚠️ **Bacalah angka itu dengan hati-hati.** Pemanggilan yang tercatat *"hidup"* pada lgk `7`, `7.2`,
`7.2.1`, `7.2.1.2` adalah **langkah INDUK yang membungkus** langkah ber-remark di dalamnya —
penelusur mencatat induknya karena teks nama anaknya ada di dalam badan induk.
⛔ **Langkah yang benar-benar memanggil pembangkit per lini — `7.2.1.2.8` sampai `7.2.1.2.14` —
BER-REMARK tujuh-tujuhnya.** ⭐ Sedangkan `GetSequenceNumber_SQL` dipanggil di **`7.2.1.2.6`**,
⭐ **hidup, bendera `true`, gerbang `OutputData.START_DATE==""`**.

> ⭐ **Vonis T6a: TIDAK ADA pengganti pembangkit per lini di rule lain.** Seluruh 114 berkas
> disisir; ketujuhnya memang mati dan **tidak digantikan oleh apa pun bernama serupa**.
> ⚠️ **Tidak ketemu di medan `pyStepsActivityName`, `RequestType`, maupun teks mentah berkas.**

### 6b — ⭐⭐ Satu rangkaian, bukan per lini

`[terverifikasi]` `RDBList\GetSequenceNumber_SQL` — ⭐ **identitas empat bagian SAMA** di Komite
Claim FacIn dan Claim Fac In *(ruleset `GISFW` 01-01-91)*.

Prosedurnya `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`, dengan **tiga parameter**:
`{ParamSeq.CARI1}` · `{ParamSeq.CARI2}` · `{ParamSeq.CARI3}`.

**Yang diisikan pemanggil** — `KomitePost_Adjustment` lgk **7.2.1.2.5**:

| Parameter | Diisi dengan |
| --- | --- |
| ⭐ `ParamSeq.CARI1` | ⭐ **`pyWorkCover.pxObjClass`** — **kelas kasus INDUK** |
| `ParamSeq.CARI2` | `ParamSeq.HASIL3 + "A"` — hasil langkah sebelumnya, berimbuhan `"A"` |
| `ParamSeq.CARI3` | ⛔ **tidak ketemu diisi di mana pun** dalam 114 berkas |

⭐⭐ **Jawaban T6b: SATU RANGKAIAN, dikunci pada KELAS KASUS INDUK — bukan per lini usaha.**
⚠️ Itulah **penggantinya**: ketujuh pembangkit per lini di-remark, dan tugasnya diambil alih oleh
**satu pembangkit yang dikunci pada kelas**.

⛔ **Yang TIDAK terbaca:** apa yang dikerjakan prosedur tersimpan itu di dalam basis data, dan apa
yang terjadi bila `CARI3` masuk kosong. ⭐ **Itu pertanyaan DBA**, bukan pertanyaan korpus —
`[data DBA]`, disambungkan ke butir 1 yang sudah terbuka.

### 6c — Diadu dengan Claim Fac In

⭐ Induknya memakai `GetSequenceNumber_SQL` **empat kali** — `CLaimFaceSheet_Act` lgk 14
*(bendera `true`)*, `DLAFacintoTreaty_Act` lgk 7 *(⚠️ bendera `false` — gerbang mati)*,
`GenerateDLAFacin_Act` lgk 9 *(⚠️ bendera `false`)*, `GeneratePLATreaty_Act` lgk 7 *(bendera `true`)*.

⚠️ **Dua dari empat pemanggilan di induk bergerbang MATI** — artinya keduanya **membangkitkan nomor
lebih sering** daripada yang tertulis di layar rancangan. ⛔ Itu menyentuh penomoran dokumen DLA dan
PLA, dan **belum pernah tercatat di spec Claim Fac In**. ⭐ Dicatat di sini sebagai **rujukan
silang**, ⛔ **bukan sebagai suntingan ke berkas modul itu.**

---

## §7 — ✅ Dua keputusan work owner — 2026-09-19

### 7.1 — Q1 → butir terbuka **4**: ⭐ **PINTU BELAKANG TIDAK DIBAWA**

> ⭐ **`[keputusan work owner]` 2026-09-19 — Gerbang `OperatorID.pyPosition=="IT Developer"` pada
> wewenang menyimpan keputusan komite TIDAK dibawa ke sistem baru. Wewenang memutuskan hanya
> pemilik `KomiteID` pada tingkat yang sedang berjalan.**

⚠️ **`[penyimpangan sadar]`** — dan ⭐ **dicatat sebagai LUBANG YANG DITUTUP, bukan fitur yang
dibuang.**

**Alasannya, ditulis lengkap:**

1. `[terverifikasi]` Gerbang itu memakai **kode arah 5**, sehingga **memotong dua syarat sisanya** —
   termasuk `Local.USERNAME == .KomiteID`, ⭐ pemeriksaan pemilik giliran yang paling langsung yang
   dipunyai modul ini.
2. ⛔ Ia bergantung pada **teks jabatan**, yaitu data kepegawaian yang **berubah tanpa sepengetahuan
   sistem**.
3. ⭐ Catatan pengembang pada rule itu — medan `pyDeleteMemo` berbunyi
   *"Add kondisi when OperatorID.pyPosition==\"IT Developer\""* — menunjukkan ia **ditambahkan
   belakangan**, bukan bagian rancangan asli.
4. ⭐ **Bukti baru ronde 3 yang menguatkan:** yang dibandingkan dengan `.KomiteID` bukan akun
   masuk apa adanya, melainkan **alias hasil pemetaan yang ditanam di kode** *(§3d)*. ⛔ Menambah
   pintu belakang di atas pemeriksaan yang **sudah rapuh** membuatnya rapuh dua lapis.

✅ Sejalan keputusan **Q4 Claim Prop** 2026-09-19 *(izin eksplisit lewat peran, bukan teks jabatan)*.

⚠️ **Yang WAJIB ikut ditulis:** kebutuhan pengembang menembus keadaan macet **tidak hilang** — ia
dipindahkan ke **jalur administratif yang tercatat**, bukan dihapus. ⛔ **Bentuk jalur itu belum
diputuskan.** ⭐ **`[terbuka]` BARU butir 25.**

⚠️ **Dan satu penemuan ronde 3 yang menambah pekerjaan pada keputusan ini:** teks `"IT Developer"`
muncul **juga** di `DataTransform\ChronologyInsertion_DT`, pada medan `pyPropertiesName` berbunyi
`OperatorID.pyPosition != "IT Developer"`. ⭐ Itu **kemunculan yang berbeda sifatnya** — ia mengatur
**catatan kronologi**, bukan wewenang keputusan — ⛔ **dan keputusan Q1 ini TIDAK menyentuhnya.**
⚠️ Dicatat supaya tidak ikut terhapus karena kesamaan kata.

### 7.2 — Q2 → butir terbuka **15**: ⚠️ **PITA NILAI DITIRU APA ADANYA**

> ⭐ **`[keputusan work owner]` 2026-09-19 — Ambang `30.000.000,00` dan `57.750.000,00` berikut
> jabatan `"SPV B"` ditiru APA ADANYA sebagai tetapan di dalam kode, bukan dijadikan baris tabel.**

⛔ **Ini BERBEDA dari rekomendasi asisten**, yang mengusulkan menjadikannya data. ⭐ **Perbedaan itu
disengaja dan dicatat apa adanya** — bukan kelalaian, dan bukan pula sesuatu yang perlu
diperdebatkan lagi.

**Letaknya** `[terverifikasi]`: `ApprovalKomite_Act` lgk **5**, dua gerbang berantai —
`OperatorID.pyPosition=="SPV B"` **DAN** `Local.TotalAdj > 30000000.00 && Local.TotalAdj <= 57750000.00`.

⚠️ **Akibat yang WAJIB ditulis, bukan disembunyikan:**

1. ⛔ **Setiap perubahan kebijakan ambang menuntut rilis perangkat lunak**, sebab angkanya hidup di
   dalam kode. ⭐ **Itu diterima sadar.**
2. ⚠️ Jabatan `"SPV B"` ikut tertanam — jadi perubahan **nama jabatan** pun menuntut rilis.
   ⛔ Ini **menarik ke arah yang berlawanan** dengan keputusan Q1, yang justru membuang
   ketergantungan pada teks jabatan. ⭐ **Keduanya tetap dijalankan sebagaimana diputuskan**; yang
   dicatat di sini hanyalah bahwa **satu modul akan memuat dua sikap berbeda terhadap teks jabatan.**

⭐ **Satu akibat yang meringankan:** butir **15** *(asal-usul angka `57.750.000`)* **BERHENTI
MEMBLOKIR** — ⛔ **tetapi TIDAK ditutup.** ⚠️ Ia tetap perlu diketahui agar orang yang kelak
mengubahnya tahu apa yang ia ubah. ✅ **Statusnya diturunkan menjadi: terbuka, tidak memblokir.**

---

## §8 — Penutup

### 8.1 — ⛔ Kesimpulan yang PALING RAWAN SALAH

> ⚠️⚠️ **§5d — vonis bahwa cabang MBU dan Travel tidak pernah terbit dari konteks komite.**

**Kenapa ia rawan.** Ia adalah **kesimpulan tentang KETIADAAN** — bahwa tak ada satu pun penulis
`Quotation.BusinessCode` pada jalur penerimaan. ⛔ Ketiadaan adalah jenis pernyataan yang paling
mudah salah, dan malam ini instrumen saya **sudah sekali** hampir memproduksi tuduhan palsu dengan
cara yang persis sama *(§8.4 nomor 3)*.

**Apa yang sudah saya sisir:** `PropertiesName` di dalam `pyParamArray` seluruh langkah,
`Page-Copy` / `Property-Set-Page` / `Page-New`, dan teks mentah, pada **Komite Claim FacIn dan
Claim Fac In**. ⛔ **Yang BELUM:** `DataTransform` sebagai penulis, `Declare-Expression`, pengisian
lewat basis data saat `Obj-Open`, dan **modul di luar kedua korpus itu**.

⭐ **Cara mematahkannya bila saya salah:** temukan satu penulis hidup `Quotation.BusinessCode` yang
berjalan sebelum `SaveAccept_ACT`. **Satu saja cukup.**

**Peringkat kedua:** §1d — bahwa putaran jalur Adjustment tidak punya pemutus hidup. Sama bentuknya,
tetapi **lebih kuat buktinya**, sebab kedua penulis `.AcceptStatus = ""` **terbaca dan terbukti
ber-remark** — bukan tidak ditemukan.

### 8.2 — ❓ Pertanyaan untuk work owner

❓ **Q1** — **Putaran persetujuan jalur Adjustment: bagaimana ia berakhir?**
`[terverifikasi]` Penghubung `IsKomiteLoop` mengembalikan alur ke `KomiteRouter` selama
`.AcceptStatus == "1"`, dan ⛔ **tidak ada langkah hidup yang mengosongkan `.AcceptStatus`**,
serta penutup paksa `ASMForceCaseClose` **ber-remark**. Pertanyaannya:
**(a)** di produksi, apakah kasus komite Fac In yang **semua anggotanya menyetujui** benar-benar
menutup sendiri? **(b)** bila ya, apa yang menutupnya?

➡️ **Rekomendasi saya: jangan tiru bentuk ini.** ⭐ Di sistem baru, **jadikan pencacah sebagai
pemutus** — putaran berhenti ketika `KomiteCount == KomiteLoop`, bukan ketika keputusan terakhir
kebetulan "tolak". ⛔ Ini **penyimpangan sadar** yang saya usulkan, sebab bentuk Pega-nya
**bergantung pada langkah yang sudah mati** dan karenanya **tidak dapat ditiru dengan jujur**.

---

❓ **Q2** — **Cabang MBU dan Travel: memang mati, atau cacat yang tak pernah ketahuan?**
`[terverifikasi]` Keduanya menguji `Quotation.BusinessCode`, yang tidak disalin ke objek kerja
komite pada jalur penerimaan. Pertanyaannya: **apakah klaim fakultatif masuk berlini MBU atau
Travel pernah melewati komite, dan apakah ringkasan akseptasinya benar?**

➡️ **Rekomendasi saya: perlakukan sebagai CACAT, dan perbaiki saat alih** — yaitu **salin
`BusinessCode` dan `BusinessName` bersama `BusinessType`** di sistem baru, sehingga kesembilan
`When` golongan RAGU membaca medan yang terisi. ⭐ Alasannya: ⛔ **tidak ada satu pun catatan
yang menyatakan cabang itu sengaja dimatikan** — yang ada hanyalah medan yang **lupa disalin**.
⚠️ Bila work owner memastikan MBU dan Travel memang **tidak pernah berkomite**, maka
rekomendasinya berbalik: **buang kesembilan cabang itu**, jangan dialihkan.

---

❓ **Q3** — **`SaveAccept_ACT` tidak menyimpan apa pun. Disengaja?**
`[terverifikasi]` Kedua langkah `SaveOSClaim_SQL` di dalamnya ber-remark, sedangkan padanannya di
`SaveReject_ACT_KMT` hidup tanpa gerbang. Pertanyaannya: **ringkasan akseptasi untuk klaim yang
DITERIMA — tersimpan di mana?**

➡️ **Rekomendasi saya: telusuri lebih dulu, jangan putuskan sekarang.** ⭐ `KomitePost_Adjustment`
lgk 7.2.1.16 memanggil **`SaveAcceptation_KMT`** — rule lain yang **belum saya baca** — dan itu
kandidat terkuat penyimpannya. ⛔ Menjawab Q3 sebelum rule itu dibaca berisiko memutuskan atas
ketiadaan yang belum terbukti.

---

❓ **Q4** — **Alias akun yang ditanam keras: apa yang menggantikannya?**
`[terverifikasi]` `Local.USERNAME` — satu-satunya sisi kiri pemeriksaan pemilik giliran — adalah
**hasil pemetaan akun yang ditanam di dalam kode**, bukan akun masuk apa adanya. Pertanyaannya:
**apa alasan pemetaan itu ada, dan apakah ia masih berlaku?**

➡️ **Rekomendasi saya: bandingkan identitas akun, bukan alias.** ⭐ Di sistem baru, pemilik giliran
disimpan sebagai **rujukan ke pengguna**, dan pembandingannya memakai rujukan itu — sehingga
pergantian nama tampil tidak menyentuh wewenang sama sekali. ⚠️ Bila ada alasan historis yang
masih hidup, ia perlu menjadi **baris data pemetaan**, ⛔ bukan cabang `@if` di dalam kode.

---

❓ **Q5** — **Penyuntingan hanya oleh penyetuju pertama: aturan, atau kebetulan?**
`[terverifikasi]` `ShowTransfer` mengunci medan saat `.KomiteCount != '1'`. Pertanyaannya: **apakah
memang hanya penyetuju pertama yang boleh mengubah isian, dan penyetuju berikutnya murni menilai?**

➡️ **Rekomendasi saya: tegaskan sebagai aturan, dan tiru.** ⭐ Bentuknya masuk akal untuk komite —
**satu orang menyiapkan, sisanya menilai hal yang sama**. ⚠️ Tanpa penegasan, ia akan terbaca
sebagai cacat layar dan hilang saat alih.

### 8.3 — Register `[terbuka]` sesudah ronde 3

**✅ DITUTUP — empat:**

| # | Butir | Ditutup oleh |
| --- | --- | --- |
| **4** | Isi `SetProteksiSubmiteKomite` | §7.1 — dinaikkan jadi Q1 dan **diputuskan** |
| ⭐ **10** | Penyaring gerbang mati belum diuji instrumennya | §4c — ✅ **LULUS 4 dari 4** |
| ⭐⭐ **16** | **Di mana `KomiteCount` dinaikkan** | §1c · §1d — ⭐ **tangga terbaca penuh** |
| **17** | Isi `SaveAccept_ACT` dan `SaveReject_ACT_KMT` | §2a–§2d |

**⭐ BARU — tujuh:**

| # | Butir | Pemilik | Memblokir? |
| --- | --- | --- | --- |
| ⚠️⚠️ **19** | **Bagaimana putaran jalur Adjustment berakhir** — pemutusnya `.AcceptStatus`, dan tak ada penulis hidup yang mengosongkannya | work owner | ⛔ **YA** — Q1 |
| ⚠️⚠️ **20** | **`IsMBU` dan `IsTravel` membaca `BusinessCode` yang tidak disalin** — cabang MBU & Travel diduga tak pernah terbit | work owner | ⛔ **YA** — Q2 |
| ⚠️ **21** | **Alias akun ditanam keras** menjadi sisi kiri pemeriksaan pemilik giliran | work owner | tidak |
| ⚠️ **22** | **`SaveAccept_ACT` tidak menyimpan; `SaveReject_ACT_KMT` menyimpan** | asisten → lalu work owner | tidak — ⭐ tunggu `SaveAcceptation_KMT` dibaca |
| **23** | **Medan `ShowTransfer` yang dapat disunting** belum terbaca | asisten | tidak |
| **24** | **`ApprovalKomite_Act` lgk 6.2 hanya memeriksa `KomiteList(1)`** — larangan menyetujui diri sendiri hanya berlaku untuk anggota pertama | work owner | tidak |
| **25** | **Bentuk jalur administratif pengganti pintu belakang** | work owner | tidak |

⛔ **RALAT dalam berkas ini sendiri.** Judul tabel di atas sempat saya tulis
*"⭐ **BARU — enam:**"* sementara isinya **tujuh baris** — butir 25 terlewat dari cacahan karena
ia lahir di §7.1, bukan di bab ini. ⭐ **Yang benar: BARU = 7.** Judulnya sudah diperbaiki;
kekeliruannya dicatat di sini, tidak dihapus.

**Sisa yang masih terbuka dari ronde sebelumnya: 1 · 6 · 11 · 12 · 13 · 14 · 15.**
*(15 tetap terbuka, ⭐ statusnya turun menjadi **tidak memblokir** — §7.2.)*

⭐ **Total, dihitung DUA CARA:**
*(a)* **cacah baris** ⇒ sisa `7` + baru `7` = **14**;
*(b)* **delta** ⇒ `11 − 4 ditutup + 7 baru` = **14**.
✅ **Keduanya sepakat: 14 butir terbuka.**

⛔ **Yang MEMBLOKIR: 1 · 6 · 19 · 20** — empat butir.
*(Butir 1 = kolom `T_QUOTATIONDATA`, `[data DBA]`; butir 6 = isi rule peran, lintas tiga modul.)*

### 8.4 — ⛔ RALAT

| # | Yang diralat | Kalimat LAMA, dikutip | Yang benar |
| --- | --- | --- | --- |
| ⭐ **1** | ronde 2 §7.5 tentang `KomiteCount` | *"⚠️ tanpa itu **tangga penyetuju tidak terbaca**"*, dan ronde 2 §7.3 butir 16 *"tangga tak terbaca"* | ⭐ **Tangganya TERBACA.** Ronde 2 hanya menyisir modul ini; ⛔ **penyemainya ada di `Claim Fac In`**, dan penaiknya ada di `KomitePost_*` — bukan di `KomiteRouter` |
| ⭐ **2** | ronde 2 §4a tentang daur hidup | *"daur hidup **SAMA**, sampai nama penghubung"* | ⭐ **Benar untuk BENTUKNYA, keliru bila dibaca sebagai rule yang sama.** Identitas empat bagian **semuanya berbeda** *(kelas `Work-Komite` lawan `Work-KomiteTreaty`)*, versinya sudah melenceng *(01-01-27 lawan 01-01-28)*, dan ⚠️ bentuk `ASSIGNMENT63` Prop membawa `pyLocalActionsString = (3)` yang **tidak ada** di sini |
| ⛔⛔ **3** | **instrumen SAYA SENDIRI, ronde ini** | Di tengah ronde ini saya menghitung *"`IsKomiteLoop` tidak dipakai di mana pun"* dan *"37 dari 49 `When` tidak terpakai"* | ⛔ **SALAH, dan nyaris menjadi tuduhan palsu terhadap ronde 1 dan 2.** Penyaring saya menyisir `pyWhenName`; ⭐ alur menyimpannya di **`pyExpression`** dengan `pyConditionType = "When"`. Angka yang benar: **dipakai 13**, **tidak dipakai 36** |

⚠️ **RALAT nomor 3 wajib dibaca sebagai pelajaran, bukan sebagai catatan kaki.** ⭐ Ia adalah
**kejadian kedua dalam modul ini** di mana penyaring yang mengandaikan **satu nama tag saja**
menghasilkan kesimpulan yang keliru — yang pertama adalah penyaring beawalan halaman di ronde 1.
⛔ **Aturan yang berlaku sejak sekarang: sebelum menyatakan sesuatu TIDAK ADA, sebut tag mana saja
yang disisir — dan sisir lebih dari satu.**

✅ Kalimat-kalimat lama di atas **tidak dihapus dari berkas ronde 1 dan 2**; ia dikutip di sini
berdampingan dengan koreksinya, sesuai kebiasaan yang berlaku.

### 8.5 — ⭐ Vonis kesiapan spec

> ⛔ **BELUM SIAP.**

⭐ **Yang sudah cukup untuk ditulis sekarang** — dan ini banyak:

| Sudah cukup | Dari |
| --- | --- |
| tangga penyetuju: penyemaian, kenaikan, jalan pintas tolak, penjaga efek akhir | §1b · §1c |
| bentuk keputusan alur dan syarat berputar | §1d |
| dua rule penyimpan keputusan, isi dan pemanggilnya | §2 |
| vonis wewenang dan hubungannya dengan ADR-0014 | §3 |
| daftar layar dan aturan "hanya penyetuju pertama yang menyunting" | §4 |
| pustaka `When`: 13 hidup, 36 beban mati | §5d |
| penomoran akseptasi: satu rangkaian dikunci kelas induk | §6b |

⛔ **Yang menahan — empat butir, dan semuanya mengubah bentuk yang akan dibangun:**

1. ⚠️⚠️ **Butir 19** — tanpa jawabannya, **daur hidup kasus komite tidak dapat ditulis**. ⛔ Ini
   bukan detail; ini **titik akhir modul**.
2. ⚠️⚠️ **Butir 20** — menentukan apakah **sembilan cabang lini usaha dialihkan atau dibuang**.
3. **Butir 6** *(isi rule peran)* — sudah menahan **tiga modul**; di sini ia menentukan siapa yang
   berhak memutuskan sesudah pintu belakang ditutup oleh Q1.
4. **Butir 1** *(kolom `T_QUOTATIONDATA`)* — `[data DBA]`, dan ⭐ ronde ini **menambah beban
   padanya**: perilaku `PROC_GENERATE_SEQUENCE_NUMBER` dan nasib `CARI3` yang masuk kosong *(§6b)*.

⭐ **Usulan urutan langkah berikutnya:**

**USULAN 1** — ⭐ **Ronde 4, pendek dan terarah**: baca `SaveAcceptation_KMT` *(menjawab butir 22
dan Q3)*, `KomitePostAct` dan `SetValueKomite` selengkapnya *(mempersempit butir 20)*, serta
`HitServiceToKasirKMT_Act` — rule terpanjang yang belum dibaca dan **satu-satunya yang menyentuh
uang keluar**. ⛔ **Empat rule, bukan satu ronde besar.**

**USULAN 2** — Bersamaan dengan itu, **ajukan keempat pertanyaan §8.2 kepada work owner**, sebab
Q1 dan Q2 **tidak dapat dijawab oleh pembacaan korpus sama sekali** — keduanya menuntut
pengetahuan tentang **apa yang benar-benar terjadi di produksi**.

⛔ **Yang saya TIDAK usulkan:** langsung menulis spec. ⚠️ Menulis daur hidup di atas butir 19 yang
masih terbuka berarti **mengarang titik akhir modul** — dan itu jenis kesalahan yang tidak
tertangkap oleh pembacaan ulang mana pun.

---

## Lampiran — bukti berkas lain tidak disentuh

| Berkas | Baris | Keadaan |
| --- | ---: | --- |
| `komite-claim-facin\grilling-ronde-1.md` | 673 | ✅ utuh |
| `komite-claim-facin\grilling-ronde-2.md` | 630 | ✅ utuh |
| `claim-facin\spec.md` | 1538 | ✅ utuh |
| `claim-prop\spec.md` | 2119 | ✅ utuh |
| `komite-claim-prop\spec.md` | 2034 | ✅ utuh |
| `OUTPUT_HASIL_RNM\CLAUDE.md` | 279 | ✅ utuh |
| `docs\adr\` | 15 berkas | ✅ utuh |
| korpus `Komite Claim FacIn` · `Claim Fac In` · `Komite Claim Prop` | 114 · 482 · 80 | ✅ dibuka hanya untuk dibaca |

⛔ Kode Go/React **NOL** · `CREATE TABLE` **NOL** · DDL **NOL** · nomor baris XML **NOL** ·
nilai kredensial **NOL** · revisi ADR **NOL**.

⚠️ **Peringatan kredensial, diulang:** `ConnectREST\SendAcceptationToKasir` membawa kredensial di
dalam ekspor. ⛔ **Nilainya tidak dibaca, tidak dicetak, dan tidak disalin ke mana pun dalam ronde
ini.** ⭐ Karena ia beredar di dalam berkas ekspor, **ia wajib diputar ulang sesudah alih** — dan
itu tanggung jawab tim pemilik layanan.
