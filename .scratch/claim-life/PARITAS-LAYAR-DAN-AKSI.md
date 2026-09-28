# Paritas layar dan aksi — Claim Life

**Paket A3, brief lanjutan 4 §3 / lanjutan 5 §3.** Disusun 27 September 2026.

Dokumen ini adalah **daftar hutang yang jujur**, bukan laporan kemenangan: setiap baris menyebut
keadaannya apa adanya, dan baris yang belum dikerjakan tetap berdiri di sini sampai ia dikerjakan.

---

## 0. Sensus — dan tiga koreksi terhadap angka brief

**Cara pertama** *(brief lanjutan 4 §3)*: 3 harness · 20 section · 16 flow action · **31 aksi layar**.

**Cara kedua** *(sensus sendiri, 27-09-2026)*:

```
grep -ho "<pyActivity>[^<]*</pyActivity>" Section/*.xml Harness/*.xml | sed 's/<[^>]*>//g' | sort -u
```

| Yang dihitung | Cara 1 | Cara 2 | Sepakat? |
| --- | ---: | ---: | --- |
| harness *(berkas `Harness\`)* | 3 | 3 | ✅ |
| section *(berkas `Section\`)* | 20 | 20 | ✅ |
| flow action *(berkas `FlowAction\`)* | 16 | 16 | ✅ |
| aksi layar *(string `<pyActivity>` unik)* | 31 | 31 | ✅ |

Keempatnya sepakat. Tetapi menghitung **string** bukan menghitung **aksi**, dan dua hal muncul
begitu daftarnya dibandingkan dengan isi `Activity\`:

### ⛔ Koreksi 1 — `SaveOutStandingLife_Act` dan `SaveOutstandingLife_Act` adalah SATU aksi

Berkasnya hanya satu: `Activity\SaveOutStandingLife_Act.xml` *(642.787 byte)*. Ejaan huruf kecil
`t` dirujuk dari `Section\AdjustmentDetail_Section.xml`; ejaan huruf besar `S` dari section lain.
Nama rule Pega tidak peka huruf besar-kecil, jadi keduanya memanggil aktivitas yang sama.

**31 string = 30 aksi berbeda.**

### ⚠️ Ralat 3 — 27-09-2026: Koreksi 2 benar, KESIMPULANNYA keliru

Koreksi 2 dan ralatnya di baris 18 sama-sama benar pada FAKTANYA: `CheckTotalAdjustmentClaim`
dirujuk sepuluh kali dan berkas rule-nya nol. Yang keliru adalah apa yang kami simpulkan darinya —
bahwa keenam angka total karena itu **tidak dapat ditiru**, sehingga layar menampilkan penanda
"belum tersedia" di enam medan uang.

Yang hilang ternyata hanya pemanggil **refresh** di layar. Yang **menghitung** ada di korpus, dua
kali, dengan rumus yang sama persis:

| Activity | Dipanggil dari | Langkah |
| --- | --- | --- |
| `Activity/SavePesertaClaim.xml` | `InputRegisterClaimLife.xml` `Submit` b27369 | 8 b4002 *(ULANG peserta, EMBEDDED b4843)* · 8.1 b4221 *(ULANG `.AdjustmentList`, EMBEDDED b4570, prasyarat KOSONG b4562)* · 8.2 b4592 *(tulis ke halaman peserta b4616–b4743)* |
| `Activity/SaveOutStandingLife_Act.xml` | `InputOSClaimLife.xml` `Save to RNM` b21102 | 23 b10638 · 23.1 b10841 *(EMBEDDED b11046)* · 23.2 b11067 *(b11091–b11217)* |

**Sebab salahnya, dicatat supaya tidak terulang:** kami berhenti pada rule yang **namanya tertulis
di section**, lalu menyimpulkan tentang **medannya**. Pertanyaan yang tidak ditanyakan: *siapa lagi
yang menulis properti itu?* Bentuk kesalahan yang sama dengan *"Close Claim tidak punya aksi lain"*
— berhenti di aksi pertama yang membawa activity, lalu menyimpulkan tentang seluruh tombol.

**Akibat kedua:** totalnya **enam**, bukan lima. `Total Ceding Retention` b20629 / `.TotalCedingRetention`
b20636 tidak ikut tercacah karena kami mencacah lewat **rujukan** `CheckTotalAdjustmentClaim`, dan ia
satu-satunya total yang tidak punya aksi refresh. Mencacah lewat pemanggil, bukan lewat label.

### ⛔ Koreksi 2 — `CheckTotalAdjustmentClaim` TIDAK ADA di korpus

Dirujuk `Section\ClaimLifeDetailGCNM.xml`, tetapi `find . -iname "*CheckTotalAdjustment*"` atas
seluruh `D:\XML\RNM_BRD\` **nol hasil**. Ia rujukan menggantung — tombol yang menunjuk aktivitas
yang tidak ikut terekspor, atau yang memang sudah tidak ada.

**30 aksi berbeda = 29 yang berkasnya dapat dibaca.**

`[terbuka — work owner]` apakah tombolnya masih hidup di Pega berjalan. Tidak memblokir: yang tidak
dapat dibaca tidak dapat ditiru, dan menebak isinya lebih buruk daripada mencatat ketiadaannya.

### ⚠️ Koreksi 3 — "menulis data" dihitung dari metode langkah, bukan dari nama

Brief mengizinkan sebuah aksi dinyatakan "tidak ditiru" **hanya** bila ada bukti ia tidak menulis
data. Buktinya diambil dari `<pyStepsActivityName>` tiap langkah, bukan dari namanya:

| Metode langkah | Menulis basis data? |
| --- | --- |
| `Obj-Save`, `RDB-Save`, `Obj-Delete`, `pxAddChildWork` | **ya** |
| `Property-Set`, `Page-Remove`, `Page-Clear-Messages`, `Page-Set-Messages` | tidak — halaman di memori |
| `RDB-List`, `Obj-Browse`, `Obj-Open-By-Handle` | tidak *(baca)*, **kecuali** SQL-nya sendiri menulis |

⛔ Dan satu aksi yang namanya terdengar seperti mekanik UI ternyata **menulis**: `setDetailClaim_act`
memanggil `Obj-Save`. Brief menyebutnya kandidat "boleh tidak ditiru *bila hanya menata halaman*" —
ia tidak hanya menata halaman.

---

## 1. Tiga puluh satu aksi layar — verdict per aksi

Kolom **keadaan**: `ada` · `A3` *(dikerjakan paket ini)* · `tidak ditiru + bukti` · `terbuka`.

| # | Aksi *(`Claim Life\Activity\`)* | Langkah yang dijalankan | Tulis? | `services` | Rute | Kontrol React | Keadaan |
| ---: | --- | --- | :---: | --- | --- | --- | --- |
| 1 | `SaveOutStandingLife_Act` | `Page-Clear-Messages`, `Property-Set`, `Page-Set-Messages`, `Page-Remove`, `RDB-List`, `Obj-Browse`, `Java`, `Call InsertJsonClaimLife_Act`, `Call serviceInsertArasapasClaimLife_act`, `Obj-Save` | **ya** | `Pendaftaran.Daftar` + `Status` | `POST /api/klaim-life` | `RegisterKlaim.tsx` | **ada** *(sebagian — jalur `ContentNote != DEATH` belum)* |
| 2 | `SavePesertaClaim` | `Page-Clear-Messages`, `Property-Set`, `Property-Set-Messages`, `Property-Remove`, `Obj-Browse`, `RDB-List`, `Page-Set-Messages`, `call FinishAssignment` | tidak | `Pendaftaran.Daftar` *(pemilihan peserta)* | `GET /api/peserta-life` | `RegisterKlaim.tsx` | **ada** *(precondition 1812 belum ditiru)* |
| 3 | `SaveInsuredClaim_Act` | `Property-Set`, `Page-Remove` | tidak | `UmurPeserta`, `ShareNusantaraReTeks` *(`repository/pilihpeserta.go`)* | *(di dalam `POST /api/klaim-life`)* | — | **ada** — 36 dari 38 medannya salinan yang backend sudah salin; dua sisanya **pilihan**: `AGE` b661, `SHARE_NUSANTARA_RE` b601/b1412 |
| 4 | `SaveAdjustment_Act` | `Property-Set`, `RDB-List` | tidak¹ | `Akseptasi.SimpanAdjustment` | `POST …/peserta/{pid}/akseptasi` | tombol **Save Adjustment** di `KlaimLife.tsx` | **ada** |
| 5 | `RejectOSClaimLife_Act` | `Property-Set`, `RDB-List`, `Obj-Save` | **ya** | `Status.Tolak` | `POST …/adjustment/{aid}/tolak` | tombol Tolak | **ada** |
| 6 | `CreateKMTLife_Act` | `Property-Set`, `Call pxRetrieveReportData`, `Call pxAddChildWork`, `Obj-Refresh-And-Lock`, `Obj-Save`, `Call SendEmailKlaimLF` | **ya** | `Penyerahan.Serahkan` + `BuatKasusKomite` | `POST …/adjustment/{aid}/komite` | tombol Serahkan | **ada** *(sisi induk belum — lihat §4)* |
| 7 | `GetListKomiteLife` | `Property-Remove`, `Property-Set`, `Property-Set-Messages`, `Call pxRetrieveReportData`, `Obj-Save` | **ya** | `RosterKomiteOracle` | *(di dalam Serahkan)* | — | **ada** |
| 8 | `UpdateDateClaimLife_Act` | `Property-Set`, `RDB-List` | **ya²** | `Tanggal.Ubah` | `PUT …/peserta/{pid}/tanggal-kejadian` | — | **ada** *(kontrol React belum)* |
| 9 | `ValidasiDOL_Act` | `Page-Clear-Messages`, `Property-Set`, `Property-Set-Messages` | tidak | `PeriksaDOL` | *(di dalam `Tanggal.Ubah`)* | — | **ada** |
| 10 | `DeletePesertaClaimLife` | `Property-Set` ×2 *(dua loop `EMBEDDED`)*, `Obj-Save` | **ya** | — | — | — | ⛔ **RALAT** — ia **tidak menghapus**; ia mengindeks ulang `.AdjustmentList(*).IndexPremiumList` lalu menyimpan *(b225, b314, b337, b417, b583, b443)*. Yang menghapus barisnya **klien**: `pyAction = deleteRow` b18017, tanpa konfirmasi b18021. Tombolnya di layar **Outstanding** *(b17909, b18039)*, dan hanya tampil saat `CLAIM_NO` kosong *(b18082)*. Rute `DELETE` bergerbang tahap **tidak dibangun** |
| 11 | `LoadDataPeserta_Act` | `Property-Set`, `RDB-List` | tidak | `PesertaPolis.AmbilUntukKlaim` | `GET /api/peserta-life` | `RegisterKlaim.tsx` | **ada** |
| 12 | `LoadDataPesertaSpesifik_Act` | `Page-Remove`, `Property-Set`, `RDB-List`, `Property-Set` ber-ULANG | tidak | `PesertaPolis.Cari` | `GET /api/peserta-life?pl=&sertifikat=&nama=` | kotak `Certificate No` + `Name of Insured` + tombol `Search` | **ada** — **tiga** kriteria *(`GetPesertaClaim_sql1.xml:85`)*; `+7 jam` b738—b904 **tidak ditiru** *(tambalan zona waktu JDBC)*; `LIKE` bersyarat = penyimpangan sadar **OQ-E** |
| 13 | `SelectAllClaimLife_act` | `Property-Set`, `Property-Set` ber-ULANG `EMBEDDED` | tidak | — | — | `lib/pilihSemua.ts` *(belum ada pemanggil)* | ⛔ **RALAT layar**: **Outstanding** *(b16633, b16710, b24489)*, bukan Register — nol `Select All` dan nol `IsAccept` di section Register. Penjungkit **tiga** keadaan b247 sudah ditulis + diuji; menunggu **grid peserta Outstanding** |
| 14 | `SearchPolicyHolder_act` | `Property-Set` | tidak | — | — | — | **A3 — Register** |
| 15 | `ValidasiClaimReceived_Act` | `RDB-List`, `Property-Set` | tidak | `models.PenandaBatasHari` | — | — | **dipindah ke A3 — Detail & Tutup** *(sheet b115)*. Aturannya **ada**; ambang `MAXEXPIREDCLAIM` **menunggu modul PremiumList Life** *(`GetProductName.xml:84` lewat `PolicyDataLife.ProductNameID`)*. Tiga cacat rule dilaporkan **OQ-G** |
| 16 | `ValidasiSTNC_Act` | `RDB-List`, `Property-Set` | tidak | `models.PenandaBatasHari` *(fungsi yang SAMA)* | — | — | **dipindah ke A3 — Detail & Tutup** *(sheet b118)*. Penanda `.STNC` → kolom `STNC_CLAIM` *(migrasi 002)*; ambang `MAXDATARECEIVE` **menunggu modul** |
| 17 | `CountClaimAmountLife_Act` | `Page-Clear-Messages`, `Property-Set`, `Property-Set-Messages`, `call SpreadingClaimLife_Act` | tidak | `HitungSpreading` | — | — | **ada** *(rute + kontrol belum)* |
| 18 | `CheckTotalAdjustmentClaim` | *(tidak terbaca — berkasnya tidak ada)* | ? | — | — | `PanelTotalPeserta` *(enam medan berangka)* | ⚠️ **RUJUKAN MENGGANTUNG, TETAPI ANGKANYA ADA** *(diralat 27-09-2026; lihat Ralat 3 di bawah)*. `ClaimLifeDetailGCNM.xml` memanggilnya **10 kali**, dan berkas rule-nya **nol** di seluruh korpus — itu tetap benar dan tetap **OQ-H**. Yang **keliru** adalah kesimpulan lamanya *("kelima total tidak dihitung; baris mana yang ikut belum terjawab")*: yang hilang hanya pemanggil **refresh**. **Nilainya** dihitung `SavePesertaClaim.xml` langkah 8 b4002 / 8.1 b4221 / 8.2 b4592 dan `SaveOutStandingLife_Act.xml` langkah 23 b10638 / 23.1 b10841 / 23.2 b11067, keduanya **berprasyarat kosong** — jadi jawabannya **seluruh baris `.AdjustmentList` peserta**, termasuk yang `STS_REJECT = 2`. Dan totalnya **enam**, bukan lima: `Total Ceding Retention` b20629 luput karena pencacahannya memakai rujukan, dan hanya total itu yang tidak punya aksi refresh. Keputusan **bc** |
| 19 | `SetIndexAdjustmentList` | `Property-Set` | tidak | `IsCheck` peserta | — | — | **ada** *(dalam kontrak API)* |
| 20 | `SetClaimXOL_Act` | `Page-Remove`, `Property-Set` | tidak | — | — | — | **A3 — Outstanding** |
| 21 | `SendtoMedical_Act` | `Property-Set` | tidak | `Tahap.KirimKeMedis` | — | — | **ada** *(rute + kontrol belum)* |
| 22 | `SendtoAdmin_Act` | `Property-Set` | tidak | `Tahap.KembaliKeAdmin` | — | — | **ada** *(rute + kontrol belum)* |
| 23 | `SendtoAdmin_Act1` | `Property-Set` | tidak | — | — | — | **A3 — Medis** *(beda dengan no. 22 dibaca di kelompok itu)* |
| 24 | `SearchDiagnose_act` | `Property-Set` | tidak | `Diagnosa.Cari` | `GET /api/penyakit-life` | `CariDiagnosa.tsx` | ✅ **ADA** *(27-09-2026)*. Kedua kriteria dihurufbesarkan **di Go** — padanan `@toUpperCase` b255 dan b302 — lalu disambung `A AND B`. Batas **500** dan halaman **50** dari rule, dijepit di kedua sisi |
| 25 | `SetDisease` | `Property-Set`, `Obj-Save` | **ya** | `Diagnosa.Ubah` | `PUT /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa/{diagId}` | `GridDiagnosa.tsx` + `CariDiagnosa.tsx` | ✅ **ADA** *(butir **bd**, 27-09-2026)*. **Dua** kolom dari rule — `.DISEASE` b260 dan `.ICDCODE` b307 — dan `GROUPDIAGNOSE` dari dropdown-nya sendiri b5886, bukan dari rule ini. b389 `Obj-Save` **`pyWorkPage`**: daftarnya menetap bersama klaimnya, dan itulah bukti kaskade 018 |
| 26 | `DownloadDocumentClaim` | `Property-Set`, `Call GetUrlGoogleStorage_Act` | tidak | `Unggahan.Unduh` | `GET /api/dokumen/{dokId}/isi` | tautan baris `PanelDokumenPeserta` | ✅ **ADA** *(butir **be**, 27-09-2026)*. Jalur unduhnya tidak menyebut klaim — meniru `GetLinkStorage_SQL.xml` b91 yang mencari dengan `imageid` saja — dan batas klaimnya ditegakkan services. ⛔ Berkas dikirim `attachment` + `nosniff`: isinya datang dari pengunggah, dan `inline` akan menjalankannya di asal kita |
| 27 | `ProtectCloseClaim_act` | `Property-Set`, `Page-Set-Messages`, `Call FinishAssignment` | tidak | — | — | `services/tutup.go` · `POST /api/klaim-life/{id}/tutup` | ✅ **ADA** *(butir bb, 27-09-2026)*. Gerbang b608 `.STS_REJECT!=1` → **409** berisi SELURUH penghalang; lolos → satu transaksi: `STATUS_WORK` = `Resolved-Completed` VERBATIM *(`Register_Flow.xml` b899, shape End1 b883)*, `TAHAP` **dikosongkan**, `TGL_UPDATE`, jejak. Hanya dari tahap **Outstanding Claim** dan **Claim Analis** — cacah berkas, `pyLocalAction>CloseClaim` ada di tepat dua section. Sesudah tutup setiap rute pengubah ditolak *(satu pintu `PastikanKasusTerbuka`, penjaga statik menagih pemanggilannya)*. **OQ-I** terbuka |
| 28 | `setDetailClaim_act` | `Obj-Open-By-Handle`, `Property-Set`, **`Obj-Save`** | ya | — | — | — | ⛔ **tidak ditiru — RESIDU UJI PENGEMBANG** *(lihat §6)* |
| 29 | `NextPrev` | `Property-Set`, `Call LoadDataPeserta_Act` | tidak | — | — | — | **tidak ditiru** — paginasi grid; nol metode tulis, dan satu-satunya panggilannya adalah pembaca no. 11 |
| 30 | `setVisibility_Act` | `Property-Set`, `Page-Remove` | tidak | — | — | — | **tidak ditiru** — nol metode tulis; React menghitung keterlihatan dari keadaannya sendiri |
| 31 | `SaveOutstandingLife_Act` | *(sama dengan no. 1)* | — | — | — | — | ⛔ **ejaan lain no. 1** *(Koreksi 1)* |

¹ `SaveAdjustment_Act` tidak memanggil `Obj-Save`; penulisannya lewat `RDB-List` atas
`UpdateOsAkseptasiClaimLife_sql`. **Ia tetap menulis** — lewat SQL, bukan lewat mesin objek Pega.

² sama: `UpdateDateClaimLife_Act` menulis lewat `UpdateDateClaimLife_SQL`.

**Rekapitulasi:** 31 string → 30 aksi → 29 berkas terbaca → **9 punya rute dan/atau kontrol hari
ini**, 2 dinyatakan **tidak ditiru dengan bukti**, 1 **tidak ada**, sisanya pekerjaan A3.

---

## 2. Enam belas flow action — satu layar React per baris

| # | Flow action | Section yang dipakainya | Halaman React | Keadaan |
| ---: | --- | --- | --- | --- |
| 1 | `InputRegisterClaimLife` | *(section sendiri)* | `RegisterKlaim.tsx` | **ada** |
| 2 | `OSClaimLife` | `InputOSClaimLife`, `PreCaimLife_Act` | — | **A3 — Outstanding** |
| 3 | `Adjustment_Detail` | `AdjustmentDetail_Section` | *(panel di `KlaimLife.tsx`)* | **sebagian** |
| 4 | `RetroClaimLife` | `RetroDetailClaimLife` | — | **A3 — Outstanding** |
| 5 | `RejectOSClaimLife` | `RejectOSClaimLife_Sec` | tombol Tolak | **sebagian** |
| 6 | `AkseptasiClaimLife` | `InputAkseptasiClaimLife`, `PreCaimLife_Act` | `Save Adjustment` · `POST …/tahap/{tujuan}` · `POST …/tutup` | ⚠️ **SEBAGIAN, LEBIH LENGKAP** *(27-09-2026)*. Keempat tombol layar itu kini terpeta: `Save` b18543 *(belum)*, `Send Back to Admin` b20221 → `outstanding`, `Send Back to Medical` b20467 → `medical-check`, `Close Claim` b21428 → penutupan *(butir bb)*. Transition5 → Decision2 *(`IsSendtoMedical` / `IsSendtoAdmin` / `Else` → End1)* karena itu seluruhnya terwakili |
| 7 | `MedicalCheck` | `MedicalCheckClaimLife`, `PreCaimLife_Act` | `POST /api/klaim-life/{id}/tahap/{tujuan}` | ✅ **PERPINDAHANNYA ADA** *(27-09-2026)*. `Send Back to Admin` b20256 → `outstanding`; `Send to Claim Analyst` b21151 → `claim-analis` *(dua aksi pada satu klik: `SendtoAdmin_Act1` b21174 **dan** `finishAssignment` b21202)*. Layar Detail menawarkannya HANYA pada tahap ini. `Save` b18572 dan telaah medisnya — **A3 lanjutan** |
| 8 | `SendtoMedical` | `SendtoMedical_Section` | `POST …/tahap/medical-check` | ✅ **ADA** — dipicu `Send Back to Medical` b20467 pada layar Claim Analis |
| 9 | `SendtoAdmin` | `SendtoAdmin_Section` | `POST …/tahap/outstanding` | ✅ **ADA** — dipicu `Send Back to Admin` b20256 *(Medical Check)* dan b20221 *(Claim Analis)* |
| 10 | `AttachDocumentLife` | `AttachDocScreenLife`, `NewAttachLife`, `SaveAttachLife`, `InsertDocument_Act`, `GetMimeType` | `POST /api/klaim-life/{id}/peserta/{pesertaId}/dokumen` | ✅ **ADA** *(butir **be**)*. Tabel MIME **48 baris + `otherwise`** dan **pemanggil menang** *(b586 `WhenTrue=2` LANJUT)* kini benar-benar dipanggil; kunci kelompok `DL-` b596 dan pengenal cap waktu b648 pula. ⛔ **Penyimpangan sadar**: prasyarat b1283 menahan penyimpanan baris selama `T_STORAGE_ID` kosong; dengan outbox jawabannya belum ada, jadi barisnya ditulis lebih dulu dan efek `storage-unggah` yang mengisinya. Yang hilang: gerbang b1283. Yang didapat: berkas yang sudah mendarat tidak pernah tanpa catatan | ⛔ **RALAT 28-09-2026**: `IMAGEID` semula memakai pengenal dokumen; rumus sebenarnya `GenerateImageID_SQL.xml` b85 — `STANDARD_HASH('ASMPP'||FF9||SYS_GUID(),'MD5')`, heksa huruf besar. Kunci penyimpanan yang dapat ditebak dari waktu unggah bukan kunci. `TANGGAL_UPLOAD` *(`Update_T_Storage_SQL.xml` b89)* pun terlewat — migrasi **020**
| 11 | `ConfirmDeleteAttachment` | `DeleteDocument_Act` | `DELETE /api/klaim-life/{id}/dokumen/{dokId}` | ✅ **ADA** *(butir **be**)*. Urutannya **dari rule dan berlawanan dengan unggah**: `Obj-Delete` b513 berjalan **tanpa prasyarat** — barisnya hilang seketika — sedangkan panggilan penyimpanan bersyarat *(b472 `WhenTrue=3` LEWATI bila `T_STORAGE_ID` kosong)* dan menyusul lewat outbox `storage-hapus`. Efek yang berjalan dua kali atas berkas yang sudah terhapus **berhasil**, bukan berputar sampai jatah percobaannya habis |
| 12 | `UploadCSV_ClaimLife` | `UploadCSVClaimLife_Act`, `pxUploadCSVResults` | — | ⛔ **tidak ditiru — MESIN BAWAAN PLATFORM**. Activity-nya hanya tiga langkah *(b250, b340, b488)*, dan langkah ketiganya memanggil `pxUploadCSVResults`. Pemetaan kolom CSV-nya **tidak ada di korpus** *(OQ-F)* |
| 13 | `CloseClaim` | `CloseClaim_Section` b1081 → `ProtectCloseClaim_act` b1101 **dan** `closeContainer` b1129 | `GET /api/klaim-life/{id}/boleh-tutup` *(memeriksa)* · `POST /api/klaim-life/{id}/tutup` *(menutup)* | ✅ **ADA** *(diralat 27-09-2026, butir bb)*. Ronde sebelumnya menulis *"sisi `Call FinishAssignment` b836 belum ada sebab tahap tujuannya belum dibaca"* — alurnya kini **sudah dibaca utuh**: 12 konektor, nol bernama `CloseClaim`, dan hanya shape `End1` yang menetapkan status kerja. Konfirmasi b499 VERBATIM; berhasil → jendela ditutup *(`closeContainer` b1129)*. Pesannya memakai **nomor sertifikat**, bukan nama — penyimpangan sadar, `kolomSalin` tiket 02. Perilaku mesin Pega untuk local action tanpa konektor senama tetap **OQ-I** |
| 13b | `DocumentLife` | `Section/DocumentLife.xml` + `LoadDocumentLife_ACT` | *(daftar di panel peserta)* | ✅ **DAFTARNYA ADA** *(27-09-2026)*. `PanelDokumenPeserta` menampilkan dokumen per **peserta** dari `T_CLAIMLF_DOCUMENT` *(FK `PREMIUM_LIST_DETAIL_ID`)*. ⛔ **Dua penyimpangan sadar**: (1) saringan Pega `KATEGORI_1 = .DOCUMENT` **tidak dapat ditiru** — kolom `DOCUMENT` tidak ada di `T_CLAIMLF_PREMIUMLIST_DETAIL` *(OQ-J)*; (2) Pega **MEMBUANG** baris yang URL penyimpanannya kosong *(prasyarat b1310 `DataImage.URLImage==""`, `WhenTrue=3`)*, dan karena penyambungan Google Storage belum ada, meniru itu membuat daftar **selalu kosong** — barisnya tetap tampil dengan penanda. `Refresh` b611, `Add attachment` b1245, `View Office Online` b3502, `Delete` b4288 dinyatakan `BelumTersedia` → **A3 — Dokumen** |
| 13c | `Diagnose_Harness` / `Diagnose_Section` | `SearchDiagnose_act`, `BrowseDiseaseLife_RD`, `SetDisease`, `SetSTS_Reject` | `GET /api/penyakit-life` · `POST`/`PUT`/`DELETE /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa[/{diagId}]` | ✅ **ADA** *(butir **bd**, 27-09-2026 — ralat atas baris ini sendiri)*. Ronde sebelumnya menulis *"`.DiagnoseList` b3923 adalah RepeatGrid (banyak) sedangkan kolomnya tunggal — **OQ-K**"* dan berhenti di situ. Jawabannya **dua baris di bawah tempat saya berhenti**: `Add` b4690 → `addRow` b4700 dan `Delete` b6160 → `deleteRow` b6170. Grid **dapat** berarti tampilan satu baris; grid ber-`Add` **dan** ber-`Delete` **tidak dapat**. OQ-K.2 ditutup, migrasi **018** dibuat. Kini: `Choose` b2509 → `SetDisease` menulis `.DISEASE` b260 dan `.ICDCODE` b307 *(keduanya `Read-only` di grid — b5374, b5566)*; `Add`/`Delete` bekerja; keempat gerbang `.STS_REJECT=='1'||=='2'` *(b4682, b5059, b5870, b6152)* ditiru sebagai **409**; `SetSTS_Reject` b241/b257 mencerminkan keputusan peserta ke **setiap** diagnosa dalam satu transaksi. ⛔ `GROUPDIAGNOSE` **tetap kosong pilihannya** — daftar b5863 hidup pada rule properti yang tidak diekspor dan tidak ada di katalog DEV: butir **bf**, **OQ-L** |
| 13d | `ClaimComite` *(`pyRuleName` `ClaimComiteeLife` b139)* | `CreateKMTLife_Act`, `GetListKomiteLife`, `SendEmailKlaimLF` | `POST …/adjustment/{adjId}/komite` | ⚠️ **SEBAGIAN**. Penyerahannya ada sejak tiket 10. ⛔ **RALAT 27-09-2026**: labelnya sempat karangan *(`Send ke Komite`, tidak ada di korpus mana pun)*; kini VERBATIM `Send Claim to Committee` **b7033**, dikunci uji yang membaca baris itu **langsung dari korpus**. `Cancel` b7715. ⛔ `SendEmailKlaimLF` menunggu persetujuan penyambungan email; `EMAILKOMITE` berisi nama dan alamat orang — dibaca saat jalan, tidak pernah disalin |
| 14 | `ShowEditClaimLife` | `EditDateClaimLife_Section` | — | **A3 — Tutup & lihat** |
| 15 | `ViewClaimDetailLifeGCNM` | `ClaimLifeDetailGCNM`, `ObjSave_Act` | *(panel detail)* | **sebagian** |
| 16 | `PL_DetailAction_ViewPolis` | `PL_DetailViewPolis_Sec` | — | **A3 — Tutup & lihat** |

⚠️ `ViewClaimDetailLifeGCNM` diluncurkan lewat `pyEditAction` dari **ketiga** layar tahap
*(Outstanding, Medis, Akseptasi)*, bergerbang `pyWorkPage.Save = 1`. Di React itu **satu panel
detail yang dipakai tiga halaman**, bukan tiga salinan.

---

## 3. Tujuh kelompok kerja A3

| Kelompok | Aksi / layar yang belum ada | Commit |
| --- | --- | --- |
| Register | 3, 10, 12, 13, 14, 15, 16 · flow action 12 | `claim-life: A3 — Register` |
| Outstanding | 20 · flow action 2, 4 | `claim-life: A3 — Outstanding` |
| Dokumen | 26 · flow action 10, 11 | `claim-life: A3 — Dokumen` |
| Medis | 23, 24, 25 · flow action 7, 8, 9 | `claim-life: A3 — Medis` |
| Akseptasi | *(inti sudah ada A0)* — sisa layar | `claim-life: A3 — Akseptasi` |
| Komite | *(inti sudah ada A2)* — sisa layar | `claim-life: A3 — Komite` |
| Tutup & lihat | 27, 28 · flow action 13, 14, 16 | `claim-life: A3 — Tutup & lihat` |

---

## 4. Hutang yang dicatat, bukan disembunyikan

| Hutang | Asal |
| --- | --- |
| sisi induk `CreateKMTLife_Act`: `.IsKomite`, `.KomiteNo`, `.TotalKomite`, `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` | brief menyebut "sepuluh langkah utuh"; A2 menutup tiga penulisan tabel |
| jalur `ContentNote != DEATH` *(baris 4054, 4199)* | A0 titik 9 — belum dijawab |
| precondition `SavePesertaClaim` 1812 *(cacat kurung)* | lanjutan 5 §3 |
| `Protect.pxResults` di 2555/2704 | lanjutan 5 §3 |
| `CheckTotalAdjustmentClaim` tidak ada di korpus | Koreksi 2 |
| awalan akseptasi `RNML-A`/`RNML-AR` tidak ada di `KODE_PRODUKSI` | A2 |

---

## 5. A3 kelompok 1 — REGISTER *(27 September 2026)*

### 5.1 Himpunan medan `Section\InputRegisterClaimLife.xml`, dibaca sebagai pohon

Lima belas label `<pyLabelPreview>`, masing-masing dengan properti yang diikatnya `<pyValue>` dan
bentuk kontrolnya `<pyFormat>`. **Tiga belas medan data terikat ke `.PolicyDataLife.*`** — halaman
POLIS, bukan isian bebas: layar Pega mengisinya ketika polis dipilih.

| Baris | Label | Properti | Kontrol | Keadaan di layar kita |
| ---: | --- | --- | --- | --- |
| 3776 | `Choose Policy No` | `.pyTemplateButton` | `pxButton` | **ada** — judul kolom polis di tabel peserta |
| 7057 | `Find Insured` | `.pyTemplateButton` | `pxButton` | **ada** — tombol pencarian |
| 20008 | `Select Insured` | `.pyTemplateButton` | `pxButton` | **ada** — judul kolom centang |
| 9104 | `Type` | `.PolicyDataLife.Type` | `pxDropdown` | **ada** *(isian; sumber dropdown-nya `[terbuka]`)* |
| 16064 | `Name of Insured` | — | `pxTextInput` | **ada** — dari `GET /api/peserta-life` |
| 9890 | `Marketing Officer` | `.PolicyDataLife.MarketingName` | `pxAutoComplete` | **panel, belum bersumber** |
| 11541 | `Ceding` | `.PolicyDataLife.CedingCoName` | `pxTextInput` | **panel, belum bersumber** |
| 11736 | `Policy Holder` | `.PolicyDataLife.PolicyHolderName` | `pxTextInput` | **panel, belum bersumber** |
| 12171 | `Class of Business` | `.PolicyDataLife.BusinessName` | `pxAutoComplete` | **panel, belum bersumber** |
| 13384 | `Date Received Email` | `.PolicyDataLife.DateReceived` | `pxDateTime` | **panel, belum bersumber** |
| 13590 | `Response Date` | `.PolicyDataLife.TanggalRespon` | `pxDateTime` | **panel, belum bersumber** |
| 13795 | `Confirmation Date` | `.PolicyDataLife.TanggalKonfirmasi` | `pxDateTime` | **panel, belum bersumber** |
| 14002 | `Status` | `.PolicyDataLife.Status` | `pxTextInput` | **panel, belum bersumber** |
| 14406 | `Updated Status` | `.PolicyDataLife.StatusUpdate` | `pxTextInput` | **panel, belum bersumber** |
| 14601 | `Realization Date` | `.PolicyDataLife.TanggalRealisasi` | `pxDateTime` | **panel, belum bersumber** |

⛔ **"Belum bersumber" DINYATAKAN di layar, bukan dihilangkan.** Medannya berdiri di panel Data Polis
dengan nilai `—` dan gaya berbeda, ditambah catatan yang menyebut sebabnya. Paritas yang tampak
lengkap padahal tidak adalah paritas yang tidak akan pernah dicari lagi.

⚠️ `[terbuka]` sumber kesepuluh medan itu. Ticket 02 `[keputusan work owner 2026-09-18]` melarang
membuat `T_CLAIMLF_POLICY`/`T_CLAIMLF_MARKETING`: data polis dan marketing **dibaca hidup** dari
tabel polis modul PremiumList Life lewat penunjuk di `T_GENERAL_CLAIM`. Jalan bacanya belum ada;
membuatnya menuntut membaca modul PremiumList Life, dan itu di luar Claim Life.

### 5.2 Tiga dropdown — sumbernya report definition, bukan daftar karangan

| Isian | Report definition | Kelas | Tabel `[data DBA — katalog DEV]` |
| --- | --- | --- | --- |
| `Ceding` | `BrowseCedingCoLife_RD.xml` | `ASM-FW-GISFW-Int-AGENT` | `POOLDATA.AGENT` 28 kolom, 429 baris |
| `Class of Business` | `BrowseBusinessLife_RD.xml` | `ASM-FW-GISFW-Int-BUSINESS` | `POOLDATA.BUSINESS` 18 kolom, 177 baris |
| `Marketing Officer` | `BrowseMarketingOfficer_RD.xml` | `ASM-FW-GISFW-Int-marketingofficer` | `POOLDATA.MARKETINGOFFICER` 14 kolom, 74 baris |

Rute `GET /api/rujukan/{jenis}?cari=` melayani ketiganya. Jenisnya **himpunan tertutup** dan
dipetakan ke nama tabel **di kode** — nama objek yang datang dari pemakai adalah injeksi lewat pintu
yang tidak dijaga penanda `:1`. Setiap kueri **berbatas 50** dan menuntut **dua huruf**; ketikan
lebih pendek dijawab daftar kosong, bukan galat.

⛔ `AGENT.CLIENTNAME` dan `MARKETINGOFFICER.CLIENTNAME` memuat **nama**; dibaca saat jalan, nol baris
disalin ke fixture/tiket/log.

⚠️ `[terbuka]` `BUSINESS` punya `NOTE` **dan** `NOTEINA`; mana yang layar Pega tampilkan belum
terbaca dari pohon section. Dipakai `NOTE`, dan pilihannya dicatat di sini alih-alih ditebak diam.

### 5.3 Aksi yang bergeser

| # | Aksi | Sebelum | Sesudah |
| ---: | --- | --- | --- |
| 14 | `SearchPolicyHolder_act` | A3 — Register | ⛔ **RALAT** — activity ini hanya `Property-Set SearchPolicyHolder.CARI1 = @toUpperCase(…)`: teks pencarian popup `SearchPolicy_Harness`, **bukan** ceding. Ia menunggu popup `Choose Policy No`, yaitu **menunggu modul PremiumList Life** |
| 11 | `LoadDataPeserta_Act` | ada | tetap — kini tombolnya berlabel VERBATIM `Find Insured` |

**Cacah aksi yang punya rute DAN kontrol: 10 dari 30** *(sebelumnya 9)*.

### 5.4 Yang BELUM dikerjakan di kelompok ini, dan sebabnya

| Aksi | Sebab |
| --- | --- |
| 3 `SaveInsuredClaim_Act` · 10 `DeletePesertaClaimLife` · 12 `LoadDataPesertaSpesifik_Act` · 13 `SelectAllClaimLife_act` | menuntut pembacaan pohon masing-masing activity; belum dibaca utuh giliran ini |
| 15 `ValidasiClaimReceived_Act` · 16 `ValidasiSTNC_Act` | sama |
| flow action 12 `UploadCSV_ClaimLife` | menuntut jalur unggah berkas, yang bergantung penyimpanan (`[terbuka]` persetujuan manusia) |

---

## 6. Ralat 27 September 2026 — butir av dan residu uji

### 6.1 Tiga dropdown Register: TIDAK ADA

§5.2 ronde pertama menyebut `Ceding`, `Class of Business`, dan `Marketing Officer` sebagai isian
ber-autocomplete. **Keliru**: keempat medan *(bersama `Type`)* `pyReadOnly` `true`, `pyEditOptions`
`Read-only`, `pyLabelFor` menunjuk `CedingCoName` / `BusinessName` / `MarketingName` / `Type`, dan
terikat `.PolicyDataLife.*`. Baris buktinya di tiket 02 blok ralat.

Sebabnya: blok kontrol berdiri **sebelum** labelnya di DOM, dan jendela pembacaan menghadap **ke
depan**. Seluruh kode yang lahir dari salah baca itu **dibuang** di tiga lapis.

`[DIPUTUSKAN — av]` sumbernya modul **PremiumList Life**; sebelas medan ditandai *menunggu modul*
di `PanelDataPolis`, cacahnya dikunci uji.

### 6.2 `setDetailClaim_act` — residu uji pengembang, bukan logika bisnis

`Activity\setDetailClaim_act.xml`, tombol `Choose` popup *(`SearchPolicy_Section.xml` 3337–3356,
3468)*:

| Baris | Langkah | Yang membuatnya residu |
| ---: | --- | --- |
| 284 | `Obj-Open-By-Handle` | handle **satu instance literal** tertulis mati di rule |
| 478 · 711 | `Property-Set` | tanggal literal Januari–Februari 2026 |
| 872 | prasyarat | `.NAME_OF_INSURED == "<satu nama literal>"` |
| 960 | `Obj-Save` | menulis kembali instance itu |

Aktivitas yang hanya berjalan untuk **satu nama orang tertentu** dan membuka **satu kasus tertentu**
bukan aturan bisnis. **Tidak ditiru.** Nilai literalnya **tidak dikutip ke mana pun** — ia memuat
nama orang dan pengenal kasus. Dilaporkan ke `OQ-untuk-tim.md` untuk pengembang Pega.

⚠️ Koreksi 3 di §0 tetap berlaku sebagai **metode**: "menulis atau tidak" dijawab dari metode
langkah. Yang berubah adalah kesimpulan untuk aksi INI — ia menulis, tetapi yang ditulisnya bukan
sesuatu yang boleh ditiru.

---

## 7. Pembaruan 27 September 2026 — giliran lanjutan 10

Baris **3, 10, 12, 13, 15, 16** tabel aksi dan baris **12** tabel alur diperbarui di tempat.

**Dua aksi berpindah kelompok** menurut XML: `ValidasiClaimReceived_Act` dan `ValidasiSTNC_Act`
**bukan** milik Register — keduanya membandingkan tanggal terhadap ambang produk dan dipakai di
panel detail, sejalan dengan sheet baris 115 dan 118.

**Dua ralat layar**, keduanya ditemukan dengan **menaiki** pohon dan bukan membaca maju:
`DeletePesertaClaimLife` dan `SelectAllClaimLife_act` berdiri di `InputOSClaimLife.xml`
*(Outstanding)*, bukan `InputRegisterClaimLife.xml`.

**Satu aksi dinyatakan tidak ditiru**: `UploadCSV_ClaimLife` memanggil mesin bawaan platform, dan
pemetaan kolomnya tidak ada di korpus. Itu menjawab tuntutan brief §1 *("bila hanya UI → dinyatakan
di PARITAS")*: ia bukan sekadar UI, melainkan **gadget platform** — dan itu pun dinyatakan.

### 7.1 Pembaruan kedua — giliran lanjutan 11

Baris **13** tabel alur *(`CloseClaim`)* diperbarui: gerbangnya ada, penyelesaian penugasannya
belum — dan itu **dinyatakan**, bukan didiamkan.

Dua aksi dibaca tetapi **sengaja belum dibangun**, sebab pemiliknya kelompok lain:
`SetIndexAdjustmentList` *(b542—b743: putaran baru **mewarisi** delapan angka dari
`.AdjustmentList(1)`)* milik **Akseptasi**; `SetSTS_Reject` *(b233: halaman langkahnya
**`.DiagnoseList`**, bukan AdjustmentList)* milik **Medis**, yang daftar diagnosisnya belum ada.


---

## Sensus akhir 28-09-2026 — satu baris per rule korpus

**135 rule** di dua belas folder `D:\XML\RNM_BRD\Claim Life\`, dicacah dengan skrip, bukan
dengan ingatan. Kolom **dipanggil dari** diisi dari pencarian nama rule di seluruh korpus *(di
luar berkasnya sendiri)*; kolom **padanan** dari pencarian yang sama di `APP_RNM/`.

| Keadaan | Cacah |
| --- | ---: |
| ✅ ada padanan kode | 85 |
| dihakimi tangan sesudah rule-nya dibaca *(delapan celah sensus)* | 9 |
| ⛔ residu *(nol pemanggil, nol padanan)* | 0 |
| **41 baris "tercatat, belum berkode" — diputuskan 28-09-2026 (giliran 10 §3.1):** | |
| ↳ ✅ ada / tidak perlu, dengan sebab dari XML | 31 |
| ↳ 🔧 celah → dibangun *(tiga utuh; `CountClaimAmountLife_Act` aturannya saja, sambungan OQ-M3)* | 4 |
| ↳ ❓ pertanyaan terbuka — keputusan work owner *(skema, desain, atau authz)* | 6 |

⛔ **Nol baris "nanti".** Kedelapan celah dibaca satu per satu dan dihakimi: **dua dibangun**
*(`SetCurrencyID_Act` + `GetCurrencyID`)*, **satu ditolak dengan sebab** *(`SetDisableAddButton`)*,
**tiga milik modul PremiumList**, **dua bagian cabang `ContentNote != DEATH`** yang PARITAS baris 1
sudah tandai.

⛔ **RALAT 28-09-2026 (giliran 10 §3.1) — keempat puluh satu baris kini berkeputusan.** Ronde
`eb9bcef` menulis *"tercatat, belum berkode — dijelaskan di …"*: itu PENUNJUK ke dokumen, bukan
keputusan. Tiap baris kini dibaca ulang sebagai pohon (`sed 's/></>\n</g'`, bNNN) dan kolom
keadaannya berisi salah satu dari tiga:

- ✅ **ada / tidak perlu** — dengan sebabnya dari XML: tanpa pemanggil, residu, milik modul lain,
  invarian sudah ditiru rule lain, atau **tidak dibangun karena gerbang persetujuan manusia**
  *(sistem luar: SMTP, Google — CLAUDE.md)*. Yang terakhir BUKAN "nanti": padanannya pengirim stub
  yang sudah ada, dan pemanggilan nyatanya menunggu manusia, bukan kode.
- 🔧 **celah → dibangun** di giliran ini: dialog Edit Date tiga tanggal, konfirmasi dua jalur
  balik, dan aturan murni Claim Paid.
- ❓ **pertanyaan terbuka** — yang membangunnya berarti memutuskan skema *(migrasi `020`+ hanya
  dari keputusan tercatat)*, desain *(data peserta dari berkas klien)*, atau penghapusan
  *(ADR-U-0031)*. Masing-masing bernomor **OQ-M1..M7** (ditambah temuan wewenang OQ-M8 dan penanda Claim Received OQ-M9) di `OQ-untuk-tim.md` dengan barisnya.

⚠️ **Ralat OQ-F.** OQ-F menyatakan pemetaan kolom CSV *"tidak ada di korpus"*. Keliru: pemetaannya
ada di `Activity/SetClaimXOL_Act.xml` b453-b1065 (26 kolom), aksi kedua tombol `Upload CSV`
`InputRegisterClaimLife` b6763. Yang tersisa bukan data yang hilang melainkan keputusan — OQ-M4.

⛔ **Empat baris ✅ di luar 41 DIRALAT (temuan /code-review, 28-09-2026)** — bukan dibiarkan:
`SaveOutStandingLife_Act` dan `Adjustment_Detail` tertulis ✅ hanya berdasar label (`labels.ts`),
`SpreadingClaimLife_Act` dan `ValidasiClaimReceived_Act` berdasar aturan murni yang nol pemanggil.
Keadaannya kini ditulis apa adanya di barisnya. `SaveOutStandingLife_Act` adalah **celah terbesar
yang tersisa** — rute simpan Outstanding — dan tidak dibangun giliran ini: ia bukan salah satu
dari 41 baris, dan separuh isinya menunggu OQ-M1/M3.

⚠️ **Temuan paritas av-2.** Layar Register dan Outstanding mengikat **21** medan `.PolicyDataLife.*`
(`InputRegisterClaimLife` b9110-b14813, `InputOSClaimLife` b6905-b12679); `PanelDataPolis`
menampilkan sebelas. Sepuluh sisanya — `TypeCeding`, `ProRateType`, `WPC`, `RetroName` *(Retro
Name / Billing Name)*, `SecurityReinsurer`, `SobName`, `ProductNameID`, `ProductName`,
`TanggalKonfirmasiBalik`, `KetentuanUnderwriting` — belum dibaca kondisi tampilnya. `RETRO_NAME`
ADA di `T_PREMIUM_LIST` (migrasi 051). Butir av di brief menyebut sepuluh medan dan sudah
dipenuhi; kesepuluh ini di luar cakupannya, dicatat supaya tidak hilang.

| Rule | Folder | Dipanggil dari | Padanan | Keadaan |
| --- | --- | --- | --- | --- |
| `CountClaimAmountLife_Act` | Activity | `Section/AdjustmentDetail_Section.xml:1709` | `internal/models/klaimbayar.go:66`, `:87` | 🔧 **celah → aturan dibangun; sambungan ❓ OQ-M3** — `CLAIM_PAID = CLAIM_GROSS × @divide(PCTClaim,100,5)` b702/b723; galat b526 bila gross > share (b949) kecuali L12/L13/L14/L18 (b972), lalu KELUAR activity (kode 6 b865) sebelum Spreading. ⚠️ Pemicunya b1709 belum hidup: kolom PCT/CLAIM_PAID tidak ada (migrasi `020`+ hanya dari keputusan), dan rute sunting adjustment tanpa kolom itu tidak punya apa pun untuk ditulis — **OQ-M3** |
| `CreateKMTLife_Act` | Activity | `Section/ClaimComite.xml:7051`, `Harness/Committe_Life.xml:8013` | `internal/repository/kasuskomite.go:8`, `internal/repository/roster.go:35` | ✅ **ada** |
| `DeleteDocument_Act` | Activity | `FlowAction/ConfirmDeleteAttachment.xml:145` | `frontend/src/services/api.ts:1216`, `internal/models/dokumenbaru.go:61` | ✅ **ada** |
| `DeleteGoogleStorage_Act` | Activity | `Activity/DeleteDocument_Act.xml:384` | `internal/models/dokumenbaru.go:19`, `internal/models/mimedokumen.go:8` | ✅ **ada** |
| `DeletePesertaClaimLife` | Activity | `Section/InputOSClaimLife.xml:17909` | — | ❓ **[pertanyaan terbuka] OQ-M6** — activity-nya nol hapus (b249, b337, Obj-Save b443), tetapi tombol `DELETE` `InputOSClaimLife` b17865 lebih dulu `deleteRow` b17874 atas grid `PremiumListDetail`: di Pega peserta DICABUT. Bertentangan dengan ADR-U-0031 (hapus = penanda) |
| `DownloadDocumentClaim` | Activity | `Section/DocumentLife.xml:3006` | `frontend/src/assets/labels.ts:369`, `frontend/src/components/PanelDokumenPeserta.tsx:18` | ✅ **ada** |
| `GetLinkService` | Activity | `Activity/DeleteGoogleStorage_Act.xml:1304`, `Activity/GetUrlGoogleStorage_Act.xml:1668` | `internal/repository/linkservice.go:8`, `internal/services/efekkeluar.go:65` | ✅ **ada** |
| `GetListKomiteLife` | Activity | `Section/AdjustmentDetail_Section.xml:15564` | `internal/services/komite.go:33`, `internal/services/komite_test.go:28` | ✅ **ada** |
| `getMaxPagination_Act` | Activity | `Section/DetailPolisLife.xml:2174`, `Section/InputRegisterClaimLife.xml:21958` | — | ⛔ **MILIK MODUL PREMIUMLIST** — Dipanggil `DetailPolisLife.xml` b2174 dan `InputRegisterClaimLife.xml` b21958 — paginasi grid **polis**, bukan grid klaim. Layarnya milik PremiumList Life; membangunnya di sini berarti dua modul memegang satu layar |
| `GetUrlGoogleStorage_Act` | Activity | `Activity/DownloadDocumentClaim.xml:369`, `Activity/LoadDocumentLife_ACT.xml:1010` | `frontend/src/components/PanelDokumenPeserta.test.ts:42`, `frontend/src/services/api.ts:70` | ✅ **ada** |
| `InsertDocument_Act` | Activity | `Activity/SaveAttachLife.xml:1466` | `internal/models/dokumenbaru.go:44`, `internal/models/imageid.go:7` | ✅ **ada** |
| `InsertGoogleStorage_Act` | Activity | `Activity/InsertDocument_Act.xml:1022` | `internal/models/dokumenbaru.go:19`, `internal/models/imageid.go:17` | ✅ **ada** |
| `InsertJsonClaimLife_Act` | Activity | `Activity/SaveOutStandingLife_Act.xml:11496` | — | ✅ **tidak perlu** — `@ASM.GetPageJSONString()` ClaimData b1381 lalu procedure b1608; JSON dibuang keputusan 2026-09-16 (`services/efekkeluar.go:10-11`) dan prinsip o |
| `InsertLogServiceClaim` | Activity | `Activity/serviceInsertArasapasClaimLife_act.xml:1020`, `RDBList/InsertLogServiceClaim.xml:5` | — | ✅ **tidak perlu** (keputusan ax) — INSERT `monitoring_klaim_log` (`RDBList/InsertLogServiceClaim.xml` b79-83); digantikan outbox `T_LOG_SERVICE_RNM` (migrasi 015, `services/antrean.go`) |
| `LoadDataPeserta_Act` | Activity | `Activity/NextPrev.xml:74`, `Section/DetailPolisLife.xml:650` | `internal/models/satutype_test.go:9`, `internal/services/dol.go:52` | ✅ **ada** |
| `LoadDataPesertaSpesifik_Act` | Activity | `Section/InputRegisterClaimLife.xml:16576` | `frontend/src/assets/labels.ts:202`, `frontend/src/pages/RegisterKlaim.tsx:155` | ✅ **ada** |
| `LoadDocumentLife_ACT` | Activity | `Section/DocumentLife.xml:249` | `frontend/src/components/PanelDokumenPeserta.tsx:4`, `frontend/src/pages/KlaimLife.tsx:522` | ✅ **ada** |
| `NewAttachLife` | Activity | `FlowAction/AttachDocumentLife.xml:31` | — | ✅ **tidak perlu** (residu platform) — Page-Remove `dragDropFileUpload` b226, Page-New `pyAttachmentPage` b369; nol efek bisnis. Formulirnya `PanelDokumenPeserta.tsx` |
| `NextPrev` | Activity | `Section/DetailPolisLife.xml:1234` | — | ✅ **tidak perlu** (milik PremiumList) — paginasi grid `DetailPolisLife` `PagiNasi.CARI1` ±1 (b353, b643) lalu `LoadDataPeserta_Act` b1693 |
| `ObjSave_Act` | Activity | `FlowAction/ViewClaimDetailLifeGCNM.xml:143` | — | ✅ **tidak perlu** (invarian ditiru) — satu `Obj-Save pyWorkPage` b222; tiap rute menyimpan di transaksinya sendiri (`DalamTransaksi`) |
| `PreCaimLife_Act` | Activity | `FlowAction/AkseptasiClaimLife.xml:24`, `FlowAction/MedicalCheck.xml:27` | — | ✅ **tidak perlu** (tanpa pembaca) — langkah 1 (b282-b758) mengisi `PolicyDataLife.BusinessID` yang hanya dibaca JSON ClaimData b1381 dan autocomplete read-only (`InputOSClaimLife` b9992); `SaveAdjustment_Act` b1638 membaca `pyWorkPage.BusinessID`, properti lain. Langkah 2 = `SetMOClaim_Act` b889 |
| `ProtectCloseClaim_act` | Activity | `Section/CloseClaim_Section.xml:1101` | `frontend/src/assets/labels.test.ts:208`, `frontend/src/assets/labels.ts:326` | ✅ **ada** |
| `RejectOSClaimLife_Act` | Activity | `Section/RejectOSClaimLife_Sec.xml:3117` | `internal/repository/klaimlife.go:290`, `internal/services/statusbaris.go:91` | ✅ **ada** |
| `SaveAdjustment_Act` | Activity | `Section/ClaimLifeDetailGCNM.xml:22665` | `frontend/src/assets/labels.ts:319`, `frontend/src/pages/KlaimLife.tsx:560` | ✅ **ada** |
| `SaveAttachLife` | Activity | `FlowAction/AttachDocumentLife.xml:179` | `internal/models/dokumenbaru.go:38`, `internal/models/mimedokumen.go:32` | ✅ **ada** |
| `SaveInsuredClaim_Act` | Activity | `Section/InputRegisterClaimLife.xml:20084` | `internal/models/klaimlife.go:222`, `internal/models/validasitanggal.go:20` | ✅ **ada** |
| `SaveOutStandingLife_Act` | Activity | `Section/InputOSClaimLife.xml:21126` | `frontend/src/assets/labels.test.ts:255`, `frontend/src/assets/labels.ts:305` | ✅ **ada — GILIRAN-11 paket 1** — `POST /api/klaim-life/{id}/outstanding` (`services/simpanrnm.go` `Simpan` + `PeriksaSimpanRNM`), tombol `Save to RNM` di `OutstandingClaimLife.tsx`. Pohon 29 langkah dipetakan di kepala berkasnya; bendera `Save` tanpa kolom (OQ-N1), klaim ganda antarklaim baru (OQ-N2), gerbang retro 27 (OQ-N3), residu (OQ-N4) |
| `SavePesertaClaim` | Activity | `Section/InputRegisterClaimLife.xml:27393` | `frontend/src/assets/labels.test.ts:255`, `frontend/src/assets/labels.ts:304` | ✅ **ada** |
| `SearchDiagnose_act` | Activity | `Section/Diagnose_Section.xml:568`, `Harness/Diagnose_Harness.xml:1557` | `frontend/src/components/CariDiagnosa.tsx:17`, `internal/models/penyakit.go:20` | ✅ **ada** |
| `SearchPolicyHolder_act` | Activity | `Section/SearchPolicy_Section.xml:886`, `Harness/SearchPolicy_Harness.xml:1877` | — | ✅ **tidak perlu** (milik PremiumList) — satu langkah huruf besar `SearchPolicyHolder.CARI1` b247-248 di pemilih polis; server sudah menghurufbesarkan cari nama |
| `SelectAllClaimLife_act` | Activity | `Section/InputOSClaimLife.xml:16633` | `frontend/src/lib/pilihSemua.test.ts:5`, `frontend/src/lib/pilihSemua.ts:1` | ✅ **ada** |
| `SendEmailKlaimLF` | Activity | `Activity/CreateKMTLife_Act.xml:12` | `internal/services/efekkeluar.go:351`, `internal/services/komite.go:469` | ✅ **ada** |
| `SendEmailWithAttachments` | Activity | `Activity/SendEmailKlaimLF.xml:2505` | — | ✅ **tidak dibangun — gerbang persetujuan manusia** (sistem luar: SMTP) — `SendEmailKlaimLF` b2505, host literal b2556/b2574; padanannya pengirim stub `EfekEmail` (`services/efekkeluar.go`) |
| `SendtoAdmin_Act` | Activity | `Activity/SendtoAdmin_Act1.xml:5`, `Section/InputOSClaimLife.xml:21863` | `frontend/src/assets/labels.ts:280`, `frontend/src/pages/OutstandingClaimLife.tsx:12` | ✅ **ada** |
| `SendtoAdmin_Act1` | Activity | `Section/InputOSClaimLife.xml:21863`, `Section/MedicalCheckClaimLife.xml:21174` | `frontend/src/assets/labels.ts:280`, `frontend/src/pages/OutstandingClaimLife.tsx:12` | ✅ **ada** |
| `SendtoMedical_Act` | Activity | `Section/SendtoMedical_Section.xml:1290` | `internal/models/tahap.go:166`, `internal/repository/klaimlife.go:545` | ✅ **ada** (invarian ditiru) — satu Property-Set `SendtoMedical="1"` b260-261 = `JalurBalikTahap` + `PerbaruiTahap`; konfirmasinya kini dibangun (baris `SendtoMedical_Section`). Cacat prasyarat b339 = OQ-C |
| `serviceInsertArasapasClaimLife_act` | Activity | `Activity/SaveOutStandingLife_Act.xml:11857` | `internal/repository/klaimlife.go:296`, `internal/services/efekkeluar.go:36` | ✅ **ada** |
| `SetClaimXOL_Act` | Activity | `Section/InputRegisterClaimLife.xml:6763` | — | ❓ **[pertanyaan terbuka] OQ-M4** — pemetaan CSV ADA — 28 kolom berkas + `PL_NUMBER` dari polis (b478-b1065; **ralat OQ-F**), tetapi memuat uang (`SUM_INSURED`, `CLAIM_AMOUNT`, share), jendela valuasi, dan `NAME_OF_INSURED` dari berkas klien — pendaftaran membaca peserta ULANG dari sumber, dan `ValidasiDOL` bersandar pada jendela yang tersimpan |
| `SetCurrencyID_Act` | Activity | `Section/AdjustmentDetail_Section.xml:1214` | — | ✅ **ADA sejak 28-09-2026** — `pyPreDataTransform` `AdjustmentDetail_Section.xml` b1206. Ia menerjemahkan `.CURRENCY` → `.CURRENCYID` lewat `GetCurrencyID` b419 **saat layar dimuat**. Ditiru di jalur **baca** (`services/klaimlife.go`), meniru letaknya bukan hanya hasilnya. ⚠️ Sebelumnya `CURRENCYID` hanya DIBAWA `WarisiKolom` dan tidak pernah DITERBITKAN — baris **pertama** peserta lahir tanpa pengenal, dan `HitungTotalPeserta` menolaknya dengan *"mata uang beragam"*: kalimat yang benar tentang hal yang salah |
| `setDetailClaim_act` | Activity | `Section/SearchPolicy_Section.xml:3356`, `Harness/SearchPolicy_Harness.xml:4334` | — | ✅ **tidak perlu** (residu) — `Obj-Open-By-Handle` atas SATU pengenal kasus berkode keras (b284/b332) |
| `SetDisease` | Activity | `Section/Diagnose_Section.xml:2528`, `Harness/Diagnose_Harness.xml:3506` | `frontend/src/assets/labels.ts:471`, `frontend/src/components/CariDiagnosa.tsx:24` | ✅ **ada** |
| `SetIndexAdjustmentList` | Activity | `Section/ClaimLifeDetailGCNM.xml:17991` | `internal/models/klaimlife.go:97`, `internal/models/totalpeserta.go:51` | ✅ **ada** |
| `SetMOClaim_Act` | Activity | `Activity/PreCaimLife_Act.xml:889` | — | ✅ **tidak perlu** (JSON dibuang) — menulis `ClaimData.MarketingData.*` (b276-b408, b1182-b1314); nol pembaca `MarketingData` di korpus di luar JSON b1381 |
| `SetSTS_Reject` | Activity | `Section/ClaimLifeDetailGCNM.xml:3224` | `frontend/src/components/GridDiagnosa.tsx:4`, `internal/models/diagnosa.go:48` | ✅ **ada** |
| `setVisibility_Act` | Activity | `Section/InputRegisterClaimLife.xml:7134` | — | ✅ **tidak perlu** (keadaan layar, ditiru) — `Find Insured` b7111 menyetel `CARI2=1`, `CARI3=""`, mengosongkan TempDetail (b244-b380); kotak pencariannya `RegisterKlaim.tsx` |
| `SpreadingClaimLife_Act` | Activity | `Activity/CountClaimAmountLife_Act.xml:1012` | `internal/models/pohonklaim.go:92`, `internal/repository/migrations/005_t_claimlf_adjustment_spreading.sql:33` | ⚠️ **RALAT 28-09-2026: aturan ada, sambungan tidak** — `services/spreading.go` `HitungSpreading` nol pemanggil produksi; masukan rate-nya tanpa pembaca — **OQ-M7** |
| `UpdateDateClaimLife_Act` | Activity | `Section/EditDateClaimLife_Section.xml:1929` | `internal/services/dol.go:281`, `internal/handlers/dol.go:137`, `frontend/src/components/claimlife/PanelTanggalKlaim.tsx:51` | 🔧 **celah → dibangun** — tiga tanggal `CLAIM_RECEIVED_DATE` b1076, `COMPLETE_DATE` b1387, `CONFIRMATION_DATE` b1626; `Save` b1910 → `PUT …/peserta/{pesertaId}/tanggal-klaim`. Gerbang b1000/b1313/b1550/b1788 ditiru separuh: `pyPosition` = tahap Outstanding + peran Admin; `CLAIM_NO != ''` **OQ-M1** |
| `UploadCSVClaimLife_Act` | Activity | `FlowAction/UploadCSV_ClaimLife.xml:144` | — | ❓ **[pertanyaan terbuka] OQ-M4** — tiga langkah `pxUploadCSVResults` (b490-b539), tombol tampil bila `PL_NUMBER != '' && CARI2 != 1` (b6386); pemetaannya `SetClaimXOL_Act` (baris itu) |
| `ValidasiClaimReceived_Act` | Activity | `Section/ClaimLifeDetailGCNM.xml:11457`, `Section/EditDateClaimLife_Section.xml:1120` | `internal/models/validasitanggal.go:9`, `internal/models/validasitanggal_test.go:5` | ⚠️ **RALAT 28-09-2026: aturan + ambang ada, sambungan tidak** — `services.AmbangKlaim.Hitung` nol pemanggil — **OQ-M9** |
| `ValidasiDOL_Act` | Activity | `Section/ClaimLifeDetailGCNM.xml:11177`, `Section/EditDateClaimLife_Section.xml:837` | `frontend/src/pages/KlaimLife.tsx:61`, `frontend/src/services/api.ts:906` | ✅ **ada** |
| `ValidasiSTNC_Act` | Activity | `Section/ClaimLifeDetailGCNM.xml:845` | `internal/models/validasitanggal.go:11`, `internal/models/validasitanggal_test.go:6` | ✅ **ada** |
| `convertJsonNusareToProductionClaimLife` | ConnectREST | `Activity/serviceInsertArasapasClaimLife_act.xml:660` | — | ✅ **tidak perlu** (JSON dibuang) — dipanggil `serviceInsertArasapas…` b660; REST Arasapas lewat `services/efekkeluar.go` tanpa JSON produksi |
| `ServiceGoogle` | ConnectREST | `Activity/DeleteGoogleStorage_Act.xml:1478`, `Activity/GetUrlGoogleStorage_Act.xml:1839` | — | ✅ **tidak dibangun — gerbang persetujuan manusia** (sistem luar: Google Storage) — Connect-REST `DeleteGoogleStorage_Act` b1420; pelaksana stub (`services/unggahan.go`) |
| `SetDisableAddButton` | DataTransform | `Section/ClaimLifeDetailGCNM.xml:19169` | — | ⛔ **TIDAK DITIRU — sebab dinyatakan** — DataTransform kelas `Int-LIFE_PREMIUM_DETAIL` yang melakukan **satu** hal: `.IsCheck = false` (b139–141). Ia `pyPreDataTransform` section `ClaimLifeDetail` (`ClaimLifeDetailGCNM.xml` b19169, section b19208). Artinya **membuka layar mencabut penanda "dipilih"**. Invariannya sudah kami punya lewat jalur yang benar — `RejectOSClaimLife_Act` mencabut `IsCheck` saat baris **ditolak** (`cabutPenanda`, tiket 05). Meniru mekanismenya berarti setiap orang yang MELIHAT layar meng-unselect pesertanya, dan itu perubahan keadaan tanpa ada yang memintanya |
| `GetMimeType` | DecisionTable | `Activity/InsertGoogleStorage_Act.xml:761` | `internal/models/mimedokumen.go:19` | ✅ **ada** |
| `Adjustment_Detail` | FlowAction | `Section/ClaimLifeDetailGCNM.xml:19583` | `frontend/src/assets/labels.test.ts:63`, `frontend/src/assets/labels.ts:76` | ⚠️ **RALAT 28-09-2026: sebagian** — panelnya ada (`KlaimLife.tsx`), penyuntingan jumlah klaim tidak — **OQ-M3** |
| `AkseptasiClaimLife` | FlowAction | `Activity/RejectOSClaimLife_Act.xml:2028`, `Activity/SaveAdjustment_Act.xml:2429` | `frontend/src/assets/labels.test.ts:48`, `frontend/src/assets/labels.ts:57` | ✅ **ada** |
| `AttachDocumentLife` | FlowAction | `Section/DocumentLife.xml:1273` | `frontend/src/assets/labels.test.ts:65`, `frontend/src/assets/labels.ts:80` | ✅ **ada** |
| `CloseClaim` | FlowAction | `Activity/ProtectCloseClaim_act.xml:6`, `Section/CloseClaim_Section.xml:6` | `frontend/src/assets/labels.test.ts:50`, `frontend/src/assets/labels.ts:61` | ✅ **ada** |
| `ConfirmDeleteAttachment` | FlowAction | `Section/ConfirmDeleteAttachment.xml:7`, `Section/DocumentLife.xml:4317` | `frontend/src/assets/labels.test.ts:67`, `frontend/src/assets/labels.ts:84` | ✅ **ada** |
| `InputRegisterClaimLife` | FlowAction | `Section/InputRegisterClaimLife.xml:8` | `frontend/src/assets/labels.test.ts:61`, `frontend/src/assets/labels.ts:72` | ✅ **ada** |
| `MedicalCheck` | FlowAction | `Section/MedicalCheckClaimLife.xml:8` | `frontend/src/assets/labels.test.ts:68`, `frontend/src/assets/labels.ts:86` | ✅ **ada** |
| `OSClaimLife` | FlowAction | `Activity/RejectOSClaimLife_Act.xml:6`, `FlowAction/RejectOSClaimLife.xml:24` | `frontend/src/App.tsx:48`, `frontend/src/assets/labels.test.ts:62` | ✅ **ada** |
| `PL_DetailAction_ViewPolis` | FlowAction | `Section/DetailPolisLife.xml:8193` | `frontend/src/assets/labels.test.ts:76`, `frontend/src/assets/labels.ts:102` | ✅ **ada** |
| `RejectOSClaimLife` | FlowAction | `Activity/RejectOSClaimLife_Act.xml:6`, `Section/AdjustmentDetail_Section.xml:15183` | `frontend/src/assets/labels.test.ts:72`, `frontend/src/assets/labels.ts:94` | ✅ **ada** |
| `RetroClaimLife` | FlowAction | `Section/AdjustmentDetail_Section.xml:10389` | `frontend/src/assets/labels.test.ts:64`, `frontend/src/assets/labels.ts:78` | ✅ **ada** |
| `SendtoAdmin` | FlowAction | `Activity/SendtoAdmin_Act.xml:5`, `Activity/SendtoAdmin_Act1.xml:5` | `frontend/src/assets/labels.test.ts:70`, `frontend/src/assets/labels.ts:90` | ✅ **ada** |
| `SendtoMedical` | FlowAction | `Activity/SendtoMedical_Act.xml:5`, `Section/InputAkseptasiClaimLife.xml:20496` | `frontend/src/assets/labels.test.ts:69`, `frontend/src/assets/labels.ts:88` | ✅ **ada** |
| `ShowEditClaimLife` | FlowAction | `Section/ClaimLifeDetailGCNM.xml:14144` | `frontend/src/assets/labels.test.ts:74`, `frontend/src/assets/labels.ts:98` | ✅ **ada** |
| `UploadCSV_ClaimLife` | FlowAction | `Section/InputRegisterClaimLife.xml:6692` | `frontend/src/assets/labels.test.ts:66`, `frontend/src/assets/labels.ts:82` | ✅ **ada** |
| `ViewClaimDetailLifeGCNM` | FlowAction | `Section/InputAkseptasiClaimLife.xml:17387`, `Section/InputOSClaimLife.xml:18252` | `frontend/src/assets/labels.test.ts:75`, `frontend/src/assets/labels.ts:100` | ✅ **ada** |
| `Committe_Life` | Harness | `Section/AdjustmentDetail_Section.xml:15558` | `internal/services/komite.go:361` | ✅ **ada** (A2) — `Penyerahan.Serahkan`, tombolnya `KlaimLife.tsx` |
| `Diagnose_Harness` | Harness | `Section/ClaimLifeDetailGCNM.xml:5081` | `frontend/src/assets/labels.ts:315`, `frontend/src/components/CariDiagnosa.tsx:3` | ✅ **ada** |
| `SearchPolicy_Harness` | Harness | `Section/InputRegisterClaimLife.xml:3847` | — | ✅ **tidak perlu** (milik PremiumList) — `Choose Policy No` b3827 membuka harness b3837 atas `InboxPremiumList_Claim`; `GET /api/polis-life` kini di `main` |
| `CountPesertaAkseptasiLife_SQL` | RDBList | `Activity/SaveOutStandingLife_Act.xml:3355`, `Activity/SavePesertaClaim.xml:1907` | — | ⚠️ **CELAH — cabang `ContentNote != DEATH`** — Namanya "Count" tetapi ia **SELECT**: mencari baris akseptasi yang sudah ada untuk orang yang sama (`CEDINGCO`+`NAME_OF_INSURED`+`DOB`+`CERTIFICATE_NO`+`PL_NUMBER`). Prasyarat pemanggilnya `SavePesertaClaim.xml` **b1997**: `.IsCheck=="true" && Business.pxResults(1).ContentNote="DEATH"`. Jadi ia bagian dari cabang `ContentNote` yang PARITAS baris 1 sudah tandai belum dibangun — bukan celah baru. ⛔ Kolom `NAME_OF_INSURED` dan `DOB` adalah data orang: dibaca saat jalan, **tidak pernah** disalin ke fixture |
| `CountPesertaAkseptasiLifeHealth_SQL` | RDBList | `Activity/SaveOutStandingLife_Act.xml:3964`, `Activity/SavePesertaClaim.xml:2323` | — | ⚠️ **CELAH — cabang `ContentNote != DEATH`** — Kembaran di atas, ditambah `LAPSE_DATE`. Sama sebabnya, sama batasnya |
| `DeleteStorage_SQL` | RDBList | `Activity/DeleteGoogleStorage_Act.xml:1635` | `internal/repository/dokumenunggah.go:205` | ✅ **ada** (ditiru) — `delete T_STORAGE_IMAGE where imageid` b85 = `HapusKartuBerkas`, dipanggil `services/unggahan.go:586` |
| `Generate_NoAccept_Life` | RDBList | `Activity/SaveAdjustment_Act.xml:747`, `RDBList/Generate_NoAccept_LifeRetro.xml:4` | `internal/models/klaimlife.go:436`, `internal/repository/migrations/014_kolom_business_code.sql:6` | ✅ **ada** |
| `Generate_NoAccept_LifeRetro` | RDBList | `Activity/SaveAdjustment_Act.xml:940` | `internal/services/akseptasi.go:64`, `internal/services/akseptasi_test.go:21` | ✅ **ada** |
| `GenerateImageID_SQL` | RDBList | `Activity/InsertGoogleStorage_Act.xml:2226` | `internal/models/imageid.go:11`, `internal/models/imageid_test.go:37` | ✅ **ada** |
| `GetAcceptedNoCL` | RDBList | `Activity/SaveAdjustment_Act.xml:2244` | `internal/repository/nomorakseptasi.go:44`, `internal/services/akseptasi.go:190` | ✅ **ada** |
| `GetAppName_SQL` | RDBList | `Activity/InsertGoogleStorage_Act.xml:1025` | — | ✅ **tidak perlu** — `SELECT APPNAME FROM T_FOLDER_IMAGE` b58 membaca nama aplikasi LAMA; aplikasi baru menulis namanya SENDIRI, `NamaAplikasiBerkas` (`services/unggahan.go:93`, `[data DBA]` `GCP_IMAGE.APPNAME`, "nilainya milik kita"). Tokennya bagian Google (gerbang persetujuan) |
| `GetCurrencyID` | RDBList | `Activity/SetCurrencyID_Act.xml:419` | — | ✅ **ADA sejak 28-09-2026** — `repository.MataUang.Pengenal` — `SELECT ID FROM CURRENCY WHERE CURRENCY = :1`, berbatas satu baris (rule aslinya membaca `pxResults(1)`). Tabel warisan, **dibaca saja** |
| `GetJsonProductLife` | RDBList | `Activity/SpreadingClaimLife_Act.xml:496` | `internal/services/spreading.go:73`, `internal/services/spreading_test.go:54` | ✅ **ada** |
| `GetKodeProdLife_SQL` | RDBList | `Activity/SaveOutStandingLife_Act.xml:7710` | `internal/repository/penomor.go:206`, `internal/services/pendaftaran.go:270` | ✅ **ada** |
| `GetLinkStorage_SQL` | RDBList | `Activity/DeleteGoogleStorage_Act.xml:692`, `Activity/GetUrlGoogleStorage_Act.xml:736` | `frontend/src/services/api.ts:1232`, `internal/handlers/dokumen.go:13` | ✅ **ada** |
| `getMaxPagination_sql` | RDBList | `Activity/getMaxPagination_Act.xml:283` | — | ⛔ **MILIK MODUL PREMIUMLIST** — Dipanggil `getMaxPagination_Act` b283 saja |
| `GetPesertaClaim_sql` | RDBList | `Activity/LoadDataPesertaSpesifik_Act.xml:545`, `Activity/LoadDataPeserta_Act.xml:892` | `frontend/src/services/api.ts:447`, `internal/repository/caripeserta_test.go:15` | ✅ **ada** |
| `GetPesertaClaim_sql1` | RDBList | `Activity/LoadDataPesertaSpesifik_Act.xml:545` | `frontend/src/services/api.ts:447`, `internal/repository/caripeserta_test.go:15` | ✅ **ada** |
| `GetProductLife` | RDBList | `Activity/SpreadingClaimLife_Act.xml:673` | — | ✅ **tidak ditiru sebagai rule** (AC 38) — mengurai `M_PRODUCT_LIFE.JSONDATA` b84-86 dengan Java b1016; pengurai JSON produk dilarang (`repository/ambangproduk.go:21`). ⚠️ NILAI yang dicarinya (`OUTWARDRATEID`, b1416-b1604) tetap dibutuhkan Spreading — sumber penggantinya **OQ-M7** |
| `GetProductName` | RDBList | `Activity/ValidasiClaimReceived_Act.xml:306`, `Activity/ValidasiSTNC_Act.xml:333` | `internal/models/validasitanggal.go:29`, `internal/models/validasitanggal_test.go:9` | ✅ **ada** |
| `GetRateRetro` | RDBList | `Activity/SpreadingClaimLife_Act.xml:1750` | `internal/services/spreading.go:107`, `internal/services/spreading_test.go:222` | ✅ **ada** |
| `GetRetroLife_SQL` | RDBList | `Activity/SpreadingClaimLife_Act.xml:2778` | `internal/services/spreading.go:89` | ✅ **ada** |
| `GetSequenceNumber_SQL` | RDBList | `Activity/SaveOutStandingLife_Act.xml:8057` | `internal/repository/penomor.go:188` | ✅ **ada** (ditiru) — `PROC_GENERATE_SEQUENCE_NUMBER` b87 ditulis ulang di Go (FOR UPDATE). ⚠️ Saatnya bergeser: Pega menomori di Save Outstanding b8057, aplikasi di pendaftaran — akibatnya **OQ-M1** |
| `GETTanggalClosing_SQL` | RDBList | `Activity/SaveOutStandingLife_Act.xml:1910` | `internal/repository/penomor.go:145`, `:118` | ✅ **ada** (ditiru) — `SELECT * FROM TANGGAL_CLOSING` b58, aturan NextMonth b2321 = `HariClosing` + `HitungPeriodeNomor` |
| `GetTokenStorage_SQL` | RDBList | `Activity/DeleteGoogleStorage_Act.xml:875`, `Activity/GetUrlGoogleStorage_Act.xml:1119` | `internal/services/efekkeluar.go:331` | ✅ **ada** |
| `Insert_T_Storage_SQL` | RDBList | `Activity/InsertGoogleStorage_Act.xml:2645` | `internal/models/imageid.go:18`, `internal/repository/dokumenunggah.go:154` | ✅ **ada** |
| `InsertJsonClaimLifeGCNM` | RDBList | `Activity/InsertJsonClaimLife_Act.xml:1608` | — | ✅ **tidak perlu** (JSON + procedure) — `PEGA_JSON_KLAIM_PNC(...)` b85-95 |
| `InsertJsonKlaimLife_sql` | RDBList | `Activity/SaveOutStandingLife_Act.xml:10182` | `internal/repository/pohonklaim.go:348` | ✅ **ada** (ditiru) — namanya menipu: INSERT `OS_AKSEPTASI_KLAIM_LIFE` b85-194 dari Save Outstanding b10182; 18 dari 55 kolom, 37 NULL dinyatakan (tiket 02/03) |
| `InsertLogServiceClaim` | RDBList | `Activity/InsertLogServiceClaim.xml:5`, `Activity/serviceInsertArasapasClaimLife_act.xml:1020` | — | ✅ **tidak perlu** (keputusan ax) — sama dengan baris Activity-nya |
| `Update_T_Storage_SQL` | RDBList | `Activity/GetUrlGoogleStorage_Act.xml:2427` | `internal/repository/dokumenunggah.go:172`, `internal/repository/migrations/019_t_claimlf_storage.sql:5` | ✅ **ada** |
| `UpdateDateClaimLife_SQL` | RDBList | `Activity/UpdateDateClaimLife_Act.xml:522` | — | ❓ **[pertanyaan terbuka] OQ-M2** — cermin warisan `OS_AKSEPTASI_KLAIM_LIFE` WHERE `CASEID` + `NAME_OF_INSURED` + `CERTIFICATE_NO` (b90): baris warisan kita membiarkan `NAME_OF_INSURED` NULL, jadi UPDATE itu mengenai NOL baris; dan `LAPSE_DATE` diisi DOL (b86) |
| `UpdateOsAkseptasiClaimLife_sql` | RDBList | `Activity/RejectOSClaimLife_Act.xml:2028`, `Activity/SaveAdjustment_Act.xml:2429` | `internal/repository/migrasidokumen.go:81`, `internal/services/tolak.go:72` | ✅ **ada** |
| `BrowseBusinessLife_RD` | ReportDefinition | `Section/InputAkseptasiClaimLife.xml:9901`, `Section/InputOSClaimLife.xml:10082` | — | ✅ **tidak perlu** (sel read-only) — `Class of Business` b9992, nilainya `PolicyDataLife.BusinessName` b10043, tampil di `PanelDataPolis` lewat `/api/polis-life/ringkas` |
| `BrowseCedingCoLife_RD` | ReportDefinition | `Section/InputAkseptasiClaimLife.xml:8013`, `Section/InputOSClaimLife.xml:8197` | — | ✅ **tidak perlu** (sel read-only, RD tidak dipakai) — `Billing Name` b8109 = `PolicyDataLife.RetroName` b8159. ⚠️ Medannya sendiri BELUM tampil — bukan urusan RD ini, melainkan temuan paritas **av-2** di bawah |
| `BrowseDiseaseLife_RD` | ReportDefinition | `Section/Diagnose_Section.xml:1751`, `Harness/Diagnose_Harness.xml:2734` | `frontend/src/components/CariDiagnosa.tsx:12`, `internal/handlers/penyakit.go:7` | ✅ **ada** |
| `BrowseFilterBusiness_RD` | ReportDefinition | `Activity/PreCaimLife_Act.xml:339` | — | ✅ **tidak perlu** — satu-satunya pemakainya `PreCaimLife_Act` b339, yang keluarannya tanpa pembaca |
| `BrowseMarketingOfficer_RD` | ReportDefinition | `Section/InputAkseptasiClaimLife.xml:7545`, `Section/InputOSClaimLife.xml:7731` | — | ✅ **tidak perlu** (sel read-only) — b7641 `pyEditOptions Read-only`, nilainya `PolicyDataLife.MarketingName` b7692, di `PanelDataPolis` |
| `FilterEmailKomiteWithLimit` | ReportDefinition | `Activity/CreateKMTLife_Act.xml:437`, `Activity/GetListKomiteLife.xml:722` | `internal/repository/roster.go:14`, `internal/repository/roster_test.go:14` | ✅ **ada** |
| `InboxPremiumList` | ReportDefinition | `Section/SearchPolicy_Section.xml:1641`, `Harness/SearchPolicy_Harness.xml:2630` | `frontend/src/assets/labels.ts:155`, `frontend/src/pages/InboxClaimLife.test.ts:92` | ✅ **ada** |
| `InboxPremiumList_Claim` | ReportDefinition | `Section/SearchPolicy_Section.xml:2362`, `Harness/SearchPolicy_Harness.xml:3346` | — | ⛔ **MILIK MODUL PREMIUMLIST** — Report definition pencarian polis, dipakai `SearchPolicy_Section.xml` b2362 dan `SearchPolicy_Harness.xml` b3346 |
| `AdjustmentDetail_Section` | Section | `FlowAction/Adjustment_Detail.xml:135` | `frontend/src/assets/labels.test.ts:46`, `frontend/src/assets/labels.ts:53` | ✅ **ada** |
| `AttachDocScreenLife` | Section | `FlowAction/AttachDocumentLife.xml:112` | `frontend/src/components/claimlife/PanelDokumenPeserta.tsx:182` | ✅ **ada** (ditiru) — berkas b812, nama berkas b2573, Category b2815; unggah banyak sekaligus (`pzMultiFilePath`) tidak ditiru |
| `ClaimComite` | Section | `Harness/Committe_Life.xml:235` | `frontend/src/assets/labels.test.ts:212`, `frontend/src/assets/labels.ts:419` | ✅ **ada** |
| `ClaimLifeDetailGCNM` | Section | `FlowAction/ViewClaimDetailLifeGCNM.xml:90`, `Section/EditDateClaimLife_Section.xml:1983` | `frontend/src/assets/labels.test.ts:45`, `frontend/src/assets/labels.ts:51` | ✅ **ada** |
| `CloseClaim_Section` | Section | `FlowAction/CloseClaim.xml:135` | `frontend/src/assets/labels.test.ts:50`, `frontend/src/assets/labels.ts:61` | ✅ **ada** |
| `ConfirmDeleteAttachment` | Section | `FlowAction/ConfirmDeleteAttachment.xml:26`, `Section/DocumentLife.xml:4317` | `frontend/src/assets/labels.test.ts:67`, `frontend/src/assets/labels.ts:84` | ✅ **ada** |
| `DetailPolisLife` | Section | `Section/InputAkseptasiClaimLife.xml:13083`, `Section/InputOSClaimLife.xml:13267` | — | ⚠️ **butir av — menunggu PremiumList tiket 04** — Section detail polis, dibuka `InputAkseptasiClaimLife.xml` b13083 dan `InputOSClaimLife.xml` b13267. Ia yang `PanelDataPolis` tiru; sepuluh medannya menunggu `repository.PolisRingkas` (pl4) menyatu ke `main` — §4 brief giliran ini |
| `Diagnose_Section` | Section | `Harness/Diagnose_Harness.xml:235` | `frontend/src/assets/labels.ts:471`, `frontend/src/components/CariDiagnosa.tsx:3` | ✅ **ada** |
| `DocumentLife` | Section | `Activity/LoadDocumentLife_ACT.xml:5`, `FlowAction/AttachDocumentLife.xml:58` | `frontend/src/assets/labels.test.ts:65`, `frontend/src/assets/labels.ts:80` | ✅ **ada** |
| `EditDateClaimLife_Section` | Section | `FlowAction/ShowEditClaimLife.xml:91` | `frontend/src/pages/KlaimLife.tsx:468`, `frontend/src/services/api.ts:905` | ✅ **ada** |
| `InputAkseptasiClaimLife` | Section | `FlowAction/AkseptasiClaimLife.xml:86` | `frontend/src/assets/labels.test.ts:48`, `frontend/src/assets/labels.ts:57` | ✅ **ada** |
| `InputOSClaimLife` | Section | `FlowAction/OSClaimLife.xml:86` | `frontend/src/assets/labels.ts:63`, `frontend/src/components/PanelPindahTahap.test.ts:16` | ✅ **ada** |
| `InputRegisterClaimLife` | Section | `FlowAction/InputRegisterClaimLife.xml:45` | `frontend/src/assets/labels.test.ts:61`, `frontend/src/assets/labels.ts:72` | ✅ **ada** |
| `MedicalCheckClaimLife` | Section | `FlowAction/MedicalCheck.xml:91` | `frontend/src/assets/labels.ts:381`, `frontend/src/components/PanelPindahTahap.test.ts:24` | ✅ **ada** |
| `PL_DetailViewPolis_Sec` | Section | `FlowAction/PL_DetailAction_ViewPolis.xml:90` | — | ✅ **tidak perlu** (milik PremiumList) — panel rincian baris grid `DetailPolisLife` b8193 = `pages/premiumlist/PremiumListDetail.tsx` |
| `RejectOSClaimLife_Sec` | Section | `FlowAction/RejectOSClaimLife.xml:88` | — | ❓ **[pertanyaan terbuka] OQ-M5** — dialog Date b790, PIC read-only b975, Remarks b1687, Submit b3117 → `RejectOSClaimLife_Act` menambah baris `KomiteList` (`IDKomite="Claim Admin"`, `KomiteAproval=2`, `KomiteComment`=Remarks, b2173-b2263). Aplikasi menolak tanpa dialog; tempat simpan Remarks belum diputuskan |
| `RetroDetailClaimLife` | Section | `FlowAction/RetroClaimLife.xml:94` | — | ❓ **[pertanyaan terbuka] OQ-M7** — expand-pane grid treaty-year (`AdjustmentDetail` b10388-b10389): `RetroLifeList` Reinsurer/Currency/Percent Share/Claim Retro (b1325-b2498). Grid induknya pun tidak tampil; `AmbilSpreading`/`HitungSpreading` nol pemanggil produksi |
| `SearchPolicy_Section` | Section | `Section/SendtoAdmin_Section.xml:212`, `Section/SendtoMedical_Section.xml:215` | — | ✅ **tidak perlu** (residu / milik PremiumList) — "pemanggil" b212/b215 hanya entri pyPagesAndClasses: nol pyInclude, satu-satunya pyStreamName dirinya sendiri (b74/b77) |
| `SendtoAdmin_Section` | Section | `FlowAction/SendtoAdmin.xml:90`, `Section/SendtoMedical_Section.xml:227` | `frontend/src/assets/labels.claimlife.ts:399`, `frontend/src/components/claimlife/PanelPindahTahap.tsx:47` | 🔧 **celah → dibangun** — "Send Back to Admin?" b566 + `Submit` b1229 sebelum pindah tahap, untuk KETIGA tombol ber-local action `SendtoAdmin` (b21433, b20285, b20250); tombol maju tetap langsung |
| `SendtoMedical_Section` | Section | `FlowAction/SendtoMedical.xml:93` | `frontend/src/assets/labels.claimlife.ts:399`, `frontend/src/components/claimlife/PanelPindahTahap.tsx:47` | 🔧 **celah → dibangun** — "Send Back to Medical?" b577 + `Submit` b1267, tombol `Send Back to Medical` b20496 |
| `LinkService` | SystemSettings | `Activity/DeleteGoogleStorage_Act.xml:1304`, `Activity/GetLinkService.xml:6` | `internal/repository/linkservice.go:8`, `internal/services/antrean.go:318` | ✅ **ada** |
| `IsPEGAPROD` | When | `Activity/InsertGoogleStorage_Act.xml:1115`, `Activity/SaveOutStandingLife_Act.xml:11919` | `internal/services/efekkeluar.go:108`, `internal/services/efekkeluar_statik_test.go:222` | ✅ **ada** |
| `IsSendtoAdmin` | When | — *(nol pemanggil)* | `internal/services/tahap.go:38`, `internal/services/tahap_test.go:153` | ✅ **ada** |
| `IsSendtoMedical` | When | — *(nol pemanggil)* | `internal/services/tahap.go:39`, `internal/services/tahap_test.go:155` | ✅ **ada** |
