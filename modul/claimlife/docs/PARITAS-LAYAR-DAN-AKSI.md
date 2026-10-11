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
| 2 | `SavePesertaClaim` | `Page-Clear-Messages`, `Property-Set`, `Property-Set-Messages`, `Property-Remove`, `Obj-Browse`, `RDB-List`, `Page-Set-Messages`, `call FinishAssignment` | tidak | `Pendaftaran.Daftar` *(pemilihan peserta)* | `GET /api/peserta-life` | `RegisterKlaim.tsx` | **ada** *(ralat sensus 28-09-2026: "precondition 1812 belum ditiru" dicabut — b1812 WHEN langkah 7.1 yang ber-precondition `false`; langkah 3, 4, 7.2–7.4, 7.6, 9 ter-remark)* |
| 3 | `SaveInsuredClaim_Act` | `Property-Set`, `Page-Remove` | tidak | `UmurPeserta`, `ShareNusantaraReTeks` *(`repository/pilihpeserta.go`)* | *(di dalam `POST /api/klaim-life`)* | — | **ada** — 36 dari 38 medannya salinan yang backend sudah salin; dua sisanya **pilihan**: `AGE` b661, `SHARE_NUSANTARA_RE` b601/b1412 |
| 4 | `SaveAdjustment_Act` | `Property-Set`, `RDB-List` | tidak¹ | `Akseptasi.SimpanAdjustment` | `POST …/peserta/{pid}/akseptasi` | tombol **Save Adjustment** di `KlaimLife.tsx` | **ada** |
| 5 | `RejectOSClaimLife_Act` | `Property-Set`, `RDB-List`, `Obj-Save` | **ya** | `Status.Tolak` | `POST …/adjustment/{aid}/tolak` | tombol Tolak | **ada** |
| 6 | `CreateKMTLife_Act` | `Property-Set`, ~~`Call pxRetrieveReportData`~~ *(langkah 1–3 ter-remark, sensus 28-09-2026)*, `Call pxAddChildWork`, `Obj-Refresh-And-Lock`, `Obj-Save`, `Call SendEmailKlaimLF` | **ya** | `Penyerahan.Serahkan` + `BuatKasusKomite` | `POST …/adjustment/{aid}/komite` | tombol Serahkan | **ada** *(sisi induk belum — lihat §4)* |
| 7 | `GetListKomiteLife` | `Property-Remove`, `Property-Set`, `Property-Set-Messages`, `Call pxRetrieveReportData`, ~~`Obj-Save`~~ | **tidak** *(ralat sensus 28-09-2026: satu-satunya `Obj-Save`, langkah 8 b1508, ter-remark)* | `RosterKomiteOracle` | *(di dalam Serahkan)* | — | **ada** |
| 8 | `UpdateDateClaimLife_Act` | `Property-Set`, `RDB-List` | **ya²** | `Tanggal.Ubah` | `PUT …/peserta/{pid}/tanggal-kejadian` | — | **ada** *(kontrol React belum)* |
| 9 | `ValidasiDOL_Act` | `Page-Clear-Messages`, `Property-Set`, `Property-Set-Messages` | tidak | `PeriksaDOL` | *(di dalam `Tanggal.Ubah`)* | — | **ada** |
| 10 | `DeletePesertaClaimLife` | `Property-Set` ×2 *(dua loop `EMBEDDED`)*, `Obj-Save` | **ya** | — | — | — | ⛔ **RALAT** — ia **tidak menghapus**; ia mengindeks ulang `.AdjustmentList(*).IndexPremiumList` lalu menyimpan *(b225, b314, b337, b417, b583, b443)*. Yang menghapus barisnya **klien**: `pyAction = deleteRow` b18017, tanpa konfirmasi b18021. Tombolnya di layar **Outstanding** *(b17909, b18039)*, dan hanya tampil saat `CLAIM_NO` kosong *(b18082)*. Rute `DELETE` bergerbang tahap **tidak dibangun** |
| 11 | `LoadDataPeserta_Act` | `Property-Set`, `RDB-List` | tidak | `PesertaPolis.AmbilUntukKlaim` | `GET /api/peserta-life` | `RegisterKlaim.tsx` | **ada** — baris **versi polis terakhir** yang hidup, penyimpangan sadar **OQ-N14** (01-10-2026) |
| 12 | `LoadDataPesertaSpesifik_Act` | `Page-Remove`, `Property-Set`, `RDB-List`, `Property-Set` ber-ULANG | tidak | `PesertaPolis.Cari` | `GET /api/peserta-life?pl=&sertifikat=&nama=` | kotak `Certificate No` + `Name of Insured` + tombol `Search` | **ada** — **tiga** kriteria *(`GetPesertaClaim_sql1.xml:85`)*; `+7 jam` b738—b904 **tidak ditiru** *(tambalan zona waktu JDBC)*; `LIKE` bersyarat = penyimpangan sadar **OQ-E**; hanya **versi polis terakhir** tiap sertifikat, `Delete`/`Batal` tidak tampil = penyimpangan sadar **OQ-N14** (01-10-2026) |
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
| 5 | `RejectOSClaimLife` | `RejectOSClaimLife_Sec` | tombol Tolak → dialog Date/PIC/Remarks/Submit (GILIRAN-17) | **ya** |
| 6 | `AkseptasiClaimLife` | `InputAkseptasiClaimLife`, `PreCaimLife_Act` | `Save Adjustment` · `POST …/tahap/{tujuan}` · `POST …/tutup` | ⚠️ **SEBAGIAN, LEBIH LENGKAP** *(27-09-2026)*. Keempat tombol layar itu kini terpeta: `Save` b18543 *(belum)*, `Send Back to Admin` b20221 → `outstanding`, `Send Back to Medical` b20467 → `medical-check`, `Close Claim` b21428 → penutupan *(butir bb)*. Transition5 → Decision2 *(`IsSendtoMedical` / `IsSendtoAdmin` / `Else` → End1)* karena itu seluruhnya terwakili |
| 7 | `MedicalCheck` | `MedicalCheckClaimLife`, `PreCaimLife_Act` | `POST /api/klaim-life/{id}/tahap/{tujuan}` | ✅ **PERPINDAHANNYA ADA** *(27-09-2026)*. `Send Back to Admin` b20256 → `outstanding`; `Send to Claim Analyst` b21151 → `claim-analis` *(dua aksi pada satu klik: `SendtoAdmin_Act1` b21174 **dan** `finishAssignment` b21202)*. Layar Detail menawarkannya HANYA pada tahap ini. `Save` b18572 dan telaah medisnya — **A3 lanjutan** |
| 8 | `SendtoMedical` | `SendtoMedical_Section` | `POST …/tahap/medical-check` | ✅ **ADA** — dipicu `Send Back to Medical` b20467 pada layar Claim Analis |
| 9 | `SendtoAdmin` | `SendtoAdmin_Section` | `POST …/tahap/outstanding` | ✅ **ADA** — dipicu `Send Back to Admin` b20256 *(Medical Check)* dan b20221 *(Claim Analis)* |
| 10 | `AttachDocumentLife` | `AttachDocScreenLife`, `NewAttachLife`, `SaveAttachLife`, `InsertDocument_Act`, `GetMimeType` | `POST /api/klaim-life/{id}/peserta/{pesertaId}/dokumen` | ✅ **ADA** *(butir **be**)*. Tabel MIME **48 baris + `otherwise`** dan **pemanggil menang** *(b586 `WhenTrue=2` LANJUT)* kini benar-benar dipanggil; kunci kelompok `DL-` b596 dan pengenal cap waktu b648 pula. ⚠️ *Sensus 28-09-2026:* b596 **memakai ulang** `Primary.DOCUMENT` peserta; kolom itu tidak ada (OQ-J), jadi setiap unggahan membuat kunci baru — celah tercatat, bukan tiruan penuh. `SaveAttachLife` 1.1, 1.3–1.6, 2, 3 ter-remark. ⛔ **Penyimpangan sadar**: prasyarat b1283 menahan penyimpanan baris selama `T_STORAGE_ID` kosong; dengan outbox jawabannya belum ada, jadi barisnya ditulis lebih dulu dan efek `storage-unggah` yang mengisinya. Yang hilang: gerbang b1283. Yang didapat: berkas yang sudah mendarat tidak pernah tanpa catatan | ⛔ **RALAT 28-09-2026**: `IMAGEID` semula memakai pengenal dokumen; rumus sebenarnya `GenerateImageID_SQL.xml` b85 — `STANDARD_HASH('ASMPP'||FF9||SYS_GUID(),'MD5')`, heksa huruf besar. Kunci penyimpanan yang dapat ditebak dari waktu unggah bukan kunci. `TANGGAL_UPLOAD` *(`Update_T_Storage_SQL.xml` b89)* pun terlewat — migrasi **020**
| 11 | `ConfirmDeleteAttachment` | `DeleteDocument_Act` | `DELETE /api/klaim-life/{id}/dokumen/{dokId}` | ✅ **ADA** *(butir **be**)*. Urutannya **dari rule dan berlawanan dengan unggah**: `Obj-Delete` b513 berjalan **tanpa prasyarat** — barisnya hilang seketika — sedangkan panggilan penyimpanan bersyarat *(b472 `WhenTrue=3` LEWATI bila `T_STORAGE_ID` kosong)* dan menyusul lewat outbox `storage-hapus`. Efek yang berjalan dua kali atas berkas yang sudah terhapus **berhasil**, bukan berputar sampai jatah percobaannya habis |
| 12 | `UploadCSV_ClaimLife` | `UploadCSVClaimLife_Act`, `pxUploadCSVResults` | — | ⛔ **tidak ditiru — MESIN BAWAAN PLATFORM**. Activity-nya hanya tiga langkah *(b250, b340, b488)*, dan langkah ketiganya memanggil `pxUploadCSVResults`. Pemetaan kolom CSV-nya **tidak ada di korpus** *(OQ-F)* |
| 13 | `CloseClaim` | `CloseClaim_Section` b1081 → `ProtectCloseClaim_act` b1101 **dan** `closeContainer` b1129 | `GET /api/klaim-life/{id}/boleh-tutup` *(memeriksa)* · `POST /api/klaim-life/{id}/tutup` *(menutup)* | ✅ **ADA** *(diralat 27-09-2026, butir bb)*. Ronde sebelumnya menulis *"sisi `Call FinishAssignment` b836 belum ada sebab tahap tujuannya belum dibaca"* — alurnya kini **sudah dibaca utuh**: 12 konektor, nol bernama `CloseClaim`, dan hanya shape `End1` yang menetapkan status kerja. Konfirmasi b499 VERBATIM; berhasil → jendela ditutup *(`closeContainer` b1129)*. Pesannya memakai **nomor sertifikat**, bukan nama — penyimpangan sadar, `kolomSalin` tiket 02. Perilaku mesin Pega untuk local action tanpa konektor senama tetap **OQ-I** |
| 13b | `DocumentLife` | `Section/DocumentLife.xml` + `LoadDocumentLife_ACT` | *(daftar di panel peserta)* | ✅ **DAFTARNYA ADA** *(27-09-2026)*. `PanelDokumenPeserta` menampilkan dokumen per **peserta** dari `T_CLAIMLF_DOCUMENT` *(FK `PREMIUM_LIST_DETAIL_ID`)*. ⛔ **Satu penyimpangan sadar**: saringan Pega `KATEGORI_1 = .DOCUMENT` **tidak dapat ditiru** — kolom `DOCUMENT` tidak ada di `T_CLAIMLF_PREMIUMLIST_DETAIL` *(OQ-J)*. *(Ralat sensus 28-09-2026: "penyimpangan" kedua — baris ber-URL kosong dibuang, b1310 — dicabut; ia milik langkah 1.4.2, anak 1.4 yang ter-remark `//` b913. Yang berjalan 1.3 b797: seluruh baris, tanpa URL. Menampilkan semua baris adalah paritas.)* `Refresh` b611, `Add attachment` b1245, `View Office Online` b3502, `Delete` b4288 dinyatakan `BelumTersedia` → **A3 — Dokumen** |
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
| ~~precondition `SavePesertaClaim` 1812 *(cacat kurung)*~~ — *dicabut sensus 28-09-2026: WHEN langkah 7.1, precondition `false` (b1587)* | lanjutan 5 §3 |
| `Protect.pxResults` di 2555 *(2704 milik 7.6, ter-remark; `Protect` hanya diisi 7.2/7.4 yang ter-remark — sensus 28-09-2026)* | lanjutan 5 §3 |
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
| ~~478~~ · 711 | `Property-Set` | tanggal literal Januari–Februari 2026 — *hanya b711 (3.1); b478 milik 2.1, di bawah langkah 2 yang ter-remark (`//` b444), dan ia menyalin BEGIN/EXPIRED, bukan tanggal literal (sensus 28-09-2026)* |
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
`InputRegisterClaimLife` b6763. Yang tersisa bukan data yang hilang melainkan keputusan — OQ-M4, **ditutup
29-09-2026 (GILIRAN-17): klaim XOL tidak membawa peserta dari berkas; tidak dibangun**.

⛔ **Empat baris ✅ di luar 41 DIRALAT (temuan /code-review, 28-09-2026)** — bukan dibiarkan:
`SaveOutStandingLife_Act` dan `Adjustment_Detail` tertulis ✅ hanya berdasar label (`labels.ts`),
`SpreadingClaimLife_Act` dan `ValidasiClaimReceived_Act` berdasar aturan murni yang nol pemanggil.
Keadaannya kini ditulis apa adanya di barisnya. `SaveOutStandingLife_Act` adalah **celah terbesar
yang tersisa** — rute simpan Outstanding — dan tidak dibangun giliran ini: ia bukan salah satu
dari 41 baris, dan separuh isinya menunggu OQ-M1/M3.

✅ **Paritas av-2 — DIBANGUN GILIRAN-11 paket 2.** Seluruh ikatan `.PolicyDataLife.*` di layar
Register dan Outstanding kini tampil di `PanelDataPolis`, urut section; `ClaimLifeDetailGCNM.xml`
mengikat **nol**. Label `pyLabelPreview` VERBATIM (uji korpus `labels.test.ts`); sumbernya
`T_PREMIUM_LIST` lewat `PolisRingkas` (dilebarkan aditif), kontrak JSON dikunci dua sisi.

| # | Properti | Label | Register | Outstanding | Sumber / keadaan |
| ---: | --- | --- | ---: | ---: | --- |
| — | `PremiumListSummary.PL_NUMBER` | PL Number | b3604 | b3085 | nomor PL klaim (sudah tampil) |
| 1 | `Type` | Type | b9110 | b6905 | `TYPE` |
| 2 | `TypeCeding` | System Reinsurance | b9307 | b7102 | `TYPE_CEDING_NAME` (kode `TYPE_CEDING` bila kosong) |
| 3 | `ProRateType` | Premium Method | b9700 | b7495 | `PRO_RATE_TYPE` |
| 4 | `MarketingName` | Marketing Officer | b9897 | b7692 | `MARKETING_NAME` |
| 5 | `WPC` | WPC | b10156 | b7953 | `WPC` |
| 6 | `RetroName` | Retro Name | b10342 | b8159 | `RETRO_NAME` — tampil bila Type TP/TR (`pyCondition` b10541 / b8363) |
| 7 | `SecurityReinsurer` | Security Reinsurer | b10623 | b8465 | `SECURITY_REINSURER` — tampil bila Type TP/TR (b10824 / b8671) |
| 8 | `CedingCoName` | Ceding | b11548 | b9414 | `CEDING_CO_NAME` |
| 9 | `PolicyHolderName` | Policy Holder | b11743 | b9609 | `POLICY_HOLDER_NAME` |
| 10 | `SobName` | SOB | b11936 | b9802 | `SOB_NAME` |
| 11 | `BusinessName` | Class of Business | b12177 | b10043 | `BUSINESS_NAME` |
| 12 | `ProductNameID` | Product Name ID | b12458 | b10324 | `PRODUCT_NAME_ID` |
| 13 | `ProductName` | Product Name | b12656 | b10522 | `PRODUCT_NAME` |
| 14 | `DateReceived` | Date Received Email | b13390 | b11256 | `DATE_RECEIVED` |
| 15 | `TanggalRespon` | Response Date | b13596 | b11462 | ⛔ tanpa kolom — penanda bernama |
| 16 | `TanggalKonfirmasi` | Confirmation Date | b13802 | b11668 | ⛔ tanpa kolom — penanda bernama |
| 17 | `Status` | Status | b14010 | b11876 | `STATUSS` |
| 18 | `TanggalKonfirmasiBalik` | Confirmation Reserved | b14205 | b12071 | ⛔ tanpa kolom — penanda bernama |
| 19 | `StatusUpdate` | Updated Status | b14413 | b12279 | `STATUS_UPDATE` |
| 20 | `TanggalRealisasi` | Realization Date | b14608 | b12474 | ⛔ tanpa kolom — penanda bernama |
| 21 | `KetentuanUnderwriting` | Underwriter Note | b14813 | b12679 | ⛔ tanpa kolom — penanda bernama |

⚠️ Dua hal yang ikut terbaca: layar Outstanding sebelumnya memanggil `<PanelDataPolis />` **tanpa
polis** — seluruh medannya tampil kosong di sana; kini ia membaca klaimnya lalu polisnya. Dan label
`pyLabelFieldValue` sel `RetroName` di Outstanding adalah *"Billing Name"* (b8130) sedangkan
`pyLabelPreview`-nya *"Retro Name"* (b8153) — yang tampil adalah `pyLabelPreview`, sama dengan
seluruh label panel.

| Rule | Folder | Dipanggil dari | Padanan | Keadaan |
| --- | --- | --- | --- | --- |
| `CountClaimAmountLife_Act` | Activity | `Section/AdjustmentDetail_Section.xml:1709` | `internal/models/klaimbayar.go:66`, `:87` | 🔧 **celah → aturan dibangun; sambungan ❓ OQ-M3** — `CLAIM_PAID = CLAIM_GROSS × @divide(PCTClaim,100,5)` b702/b723; galat b526 bila gross > share (b949) kecuali L12/L13/L14/L18 (b972), lalu KELUAR activity (kode 6 b865) sebelum Spreading. ⚠️ Pemicunya b1709 belum hidup: kolom PCT/CLAIM_PAID tidak ada (migrasi `020`+ hanya dari keputusan), dan rute sunting adjustment tanpa kolom itu tidak punya apa pun untuk ditulis — **OQ-M3** |
| `CreateKMTLife_Act` | Activity | `Section/ClaimComite.xml:7051`, `Harness/Committe_Life.xml:8013` | `internal/repository/kasuskomite.go:8`, `internal/repository/roster.go:35` | ✅ **ada** |
| `DeleteDocument_Act` | Activity | `FlowAction/ConfirmDeleteAttachment.xml:145` | `frontend/src/services/api.ts:1216`, `internal/models/dokumenbaru.go:61` | ✅ **ada** |
| `DeleteGoogleStorage_Act` | Activity | `Activity/DeleteDocument_Act.xml:384` | `internal/models/dokumenbaru.go:19`, `internal/models/mimedokumen.go:8` | ✅ **ada** |
| `DeletePesertaClaimLife` | Activity | `Section/InputOSClaimLife.xml:17909` | `internal/services/cabutpeserta.go`, migrasi 022 | ✅ **ditiru sebagai penanda — OQ-M6 ditutup 29-09-2026 (GILIRAN-17)**; dahulu: activity-nya nol hapus (b249, b337, Obj-Save b443), tetapi tombol `DELETE` `InputOSClaimLife` b17865 lebih dulu `deleteRow` b17874 atas grid `PremiumListDetail`: di Pega peserta DICABUT. Bertentangan dengan ADR-U-0031 (hapus = penanda) |
| `DownloadDocumentClaim` | Activity | `Section/DocumentLife.xml:3006` | `frontend/src/assets/labels.ts:369`, `frontend/src/components/PanelDokumenPeserta.tsx:18` | ✅ **ada** |
| `GetLinkService` | Activity | `Activity/DeleteGoogleStorage_Act.xml:1304`, `Activity/GetUrlGoogleStorage_Act.xml:1668` | `internal/repository/linkservice.go:8`, `internal/services/efekkeluar.go:65` | ✅ **ada** |
| `GetListKomiteLife` | Activity | `Section/AdjustmentDetail_Section.xml:15564` | `internal/services/komite.go:33`, `internal/services/komite_test.go:28` | ✅ **ada** |
| `getMaxPagination_Act` | Activity | `Section/DetailPolisLife.xml:2174`, `Section/InputRegisterClaimLife.xml:21958` | — | ⛔ **MILIK MODUL PREMIUMLIST** — Dipanggil `DetailPolisLife.xml` b2174 dan `InputRegisterClaimLife.xml` b21958 — paginasi grid **polis**, bukan grid klaim. Layarnya milik PremiumList Life; membangunnya di sini berarti dua modul memegang satu layar |
| `GetUrlGoogleStorage_Act` | Activity | `Activity/DownloadDocumentClaim.xml:369`, `Activity/LoadDocumentLife_ACT.xml:1010` | `frontend/src/components/PanelDokumenPeserta.test.ts:42`, `frontend/src/services/api.ts:70` | ✅ **ada** |
| `InsertDocument_Act` | Activity | `Activity/SaveAttachLife.xml:1466` | `internal/models/dokumenbaru.go:44`, `internal/models/imageid.go:7` | ✅ **ada** |
| `InsertGoogleStorage_Act` | Activity | `Activity/InsertDocument_Act.xml:1022` | `internal/models/dokumenbaru.go:19`, `internal/models/imageid.go:17` | ✅ **ada** |
| `InsertJsonClaimLife_Act` | Activity | `Activity/SaveOutStandingLife_Act.xml:11496` | `internal/services/simpanrnm.go` (`PolisRetro`) | ✅ **tidak perlu — kecuali langkah 2** — `@ASM.GetPageJSONString()` ClaimData b1381 lalu procedure b1608; JSON dibuang keputusan 2026-09-16 (`services/efekkeluar.go:10-11`) dan prinsip o. ⛔ *Diralat 28-09-2026:* langkah 2 (tidak ter-remark) **menukar** `RetroID`/`SecurityReinsurerID` pada salinan yang dibaca gerbang langkah 27 Save to RNM (b1152–1195, WHEN b1268/b1291/b1314); tukarnya ditiru di `ArasapasDilewatiRetro`, `ProdDateTime` tanpa sumber — OQ-N5 |
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
| `SaveOutStandingLife_Act` | Activity | `Section/InputOSClaimLife.xml:21126` | `frontend/src/assets/labels.test.ts:255`, `frontend/src/assets/labels.ts:305` | ✅ **ada — GILIRAN-11 paket 1** — `POST /api/klaim-life/{id}/outstanding` (`services/simpanrnm.go` `Simpan` + `PeriksaSimpanRNM`), tombol `Save to RNM` di `OutstandingClaimLife.tsx`. Pohon 29 langkah dipetakan di kepala berkasnya; bendera `Save` tanpa kolom (OQ-N1), klaim ganda antarklaim baru (OQ-N2), gerbang retro 27 (OQ-N3), residu (OQ-N4). ⛔ *Diralat 28-09-2026:* delapan langkah ter-remark (`pyStepsBlockName = //` — 11.3, 11.9, 11.11, 12, 13–15, 23); gerbang STNC dan dokumen lengkap yang sempat ditegakkan dibuang, penandaan Outstanding kini baris adjustment SAJA (`TandaiBarisOutstanding`), gerbang 27 membaca nilai sesudah tukar (OQ-N5); kelengkapan per kategori — OQ-N6, ditutup butir bl (GILIRAN-12): tidak ditiru, `PeriksaDokumenLengkap` dibuang |
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
| `SetClaimXOL_Act` | Activity | `Section/InputRegisterClaimLife.xml:6763` | — | ✅ **tidak dibangun — keputusan work owner 29-09-2026 (OQ-M4 ditutup)**; dahulu: pemetaan CSV ADA — 28 kolom berkas + `PL_NUMBER` dari polis (b478-b1065; **ralat OQ-F**), tetapi memuat uang (`SUM_INSURED`, `CLAIM_AMOUNT`, share), jendela valuasi, dan `NAME_OF_INSURED` dari berkas klien — pendaftaran membaca peserta ULANG dari sumber, dan `ValidasiDOL` bersandar pada jendela yang tersimpan |
| `SetCurrencyID_Act` | Activity | `Section/AdjustmentDetail_Section.xml:1214` | — | ✅ **ADA sejak 28-09-2026** — `pyPreDataTransform` `AdjustmentDetail_Section.xml` b1206. Ia menerjemahkan `.CURRENCY` → `.CURRENCYID` lewat `GetCurrencyID` b419 **saat layar dimuat**. Ditiru di jalur **baca** (`services/klaimlife.go`), meniru letaknya bukan hanya hasilnya. ⚠️ Sebelumnya `CURRENCYID` hanya DIBAWA `WarisiKolom` dan tidak pernah DITERBITKAN — baris **pertama** peserta lahir tanpa pengenal, dan `HitungTotalPeserta` menolaknya dengan *"mata uang beragam"*: kalimat yang benar tentang hal yang salah |
| `setDetailClaim_act` | Activity | `Section/SearchPolicy_Section.xml:3356`, `Harness/SearchPolicy_Harness.xml:4334` | — | ✅ **tidak perlu** (residu) — `Obj-Open-By-Handle` atas SATU pengenal kasus berkode keras (b284/b332) |
| `SetDisease` | Activity | `Section/Diagnose_Section.xml:2528`, `Harness/Diagnose_Harness.xml:3506` | `frontend/src/assets/labels.ts:471`, `frontend/src/components/CariDiagnosa.tsx:24` | ✅ **ada** |
| `SetIndexAdjustmentList` | Activity | `Section/ClaimLifeDetailGCNM.xml:17991` | `internal/models/klaimlife.go:97`, `internal/models/totalpeserta.go:51` | ✅ **ada** |
| `SetMOClaim_Act` | Activity | `Activity/PreCaimLife_Act.xml:889` | — | ✅ **tidak perlu** (JSON dibuang) — menulis `ClaimData.MarketingData.*` (b276-b408, b1182-b1314); nol pembaca `MarketingData` di korpus di luar JSON b1381 |
| `SetSTS_Reject` | Activity | `Section/ClaimLifeDetailGCNM.xml:3224` | `frontend/src/components/GridDiagnosa.tsx:4`, `internal/models/diagnosa.go:48` | ✅ **ada** |
| `setVisibility_Act` | Activity | `Section/InputRegisterClaimLife.xml:7134` | — | ✅ **tidak perlu** (keadaan layar, ditiru) — `Find Insured` b7111 menyetel `CARI2=1`, `CARI3=""`, mengosongkan TempDetail (b244-b380); kotak pencariannya `RegisterKlaim.tsx` |
| `SpreadingClaimLife_Act` | Activity | `Activity/CountClaimAmountLife_Act.xml:1012` | `internal/models/pohonklaim.go:92`, `internal/repository/migrations/005_t_claimlf_adjustment_spreading.sql:33` | ⚠️ **RALAT 28-09-2026: aturan ada, sambungan tidak** — `services/spreading.go` `HitungSpreading` nol pemanggil produksi; masukan rate-nya tanpa pembaca — **OQ-M7** |
| `UpdateDateClaimLife_Act` | Activity | `Section/EditDateClaimLife_Section.xml:1929` | `internal/services/dol.go:281`, `internal/handlers/dol.go:137`, `frontend/src/components/claimlife/PanelTanggalKlaim.tsx:51` | 🔧 **celah → dibangun** — tiga tanggal `CLAIM_RECEIVED_DATE` b1076, `COMPLETE_DATE` b1387, `CONFIRMATION_DATE` b1626; `Save` b1910 → `PUT …/peserta/{pesertaId}/tanggal-klaim`. Gerbang b1000/b1313/b1550/b1788 ditiru separuh: `pyPosition` = tahap Outstanding + peran Admin; `CLAIM_NO != ''` **OQ-M1** |
| `UploadCSVClaimLife_Act` | Activity | `FlowAction/UploadCSV_ClaimLife.xml:144` | — | ✅ **tidak dibangun — keputusan work owner 29-09-2026 (OQ-M4 ditutup)**; tiga langkah `pxUploadCSVResults` (b490-b539), tombol tampil bila `PL_NUMBER != '' && CARI2 != 1` (b6386); pemetaannya `SetClaimXOL_Act` (baris itu) |
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
| `CountPesertaAkseptasiLife_SQL` | RDBList | `Activity/SaveOutStandingLife_Act.xml:3355`, `Activity/SavePesertaClaim.xml:1907` | `internal/repository/gandawarisan.go` (`StatusWarisanTerakhir`) | ✅ **ada — GILIRAN-11 paket 1** *(diralat 28-09-2026)*: Save to RNM langkah 11.2/11.4 (`services/simpanrnm.go`). Pemanggil `SavePesertaClaim.xml:1907` ada di langkah 7.2 yang **ter-remark** (`//` b1863) — tidak perlu. *(Teks lama: "CELAH"; prasyarat yang dikutipnya, b1997, milik langkah ter-remark itu.)* Namanya "Count" tetapi ia **SELECT**: mencari baris akseptasi yang sudah ada untuk orang yang sama (`CEDINGCO`+`NAME_OF_INSURED`+`DOB`+`CERTIFICATE_NO`+`PL_NUMBER`). ⛔ Kolom `NAME_OF_INSURED` dan `DOB` adalah data orang: dibaca saat jalan, **tidak pernah** disalin ke fixture |
| `CountPesertaAkseptasiLifeHealth_SQL` | RDBList | `Activity/SaveOutStandingLife_Act.xml:3964`, `Activity/SavePesertaClaim.xml:2323` | `internal/repository/gandawarisan.go` (`AdaWarisanSamaDOL`) | ✅ **ada — GILIRAN-11 paket 1** *(diralat 28-09-2026)*: Save to RNM langkah 11.5/11.6. Kembaran di atas, ditambah `LAPSE_DATE`; pemanggil `SavePesertaClaim` langkah 7.4 ter-remark (`//` b2279) |
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
| `GetPesertaClaim_sql1` | RDBList | `Activity/LoadDataPesertaSpesifik_Act.xml:545` | `frontend/src/services/api.ts:447`, `internal/repository/caripeserta_test.go:15` | ✅ **ada** — ⚠️ menyimpang sadar: versi polis terakhir (OQ-N14, 01-10-2026) |
| `GetProductLife` | RDBList | `Activity/SpreadingClaimLife_Act.xml:673` | — | ✅ **tidak ditiru sebagai rule** (AC 38) — mengurai `M_PRODUCT_LIFE.JSONDATA` b84-86 dengan Java b1016; pengurai JSON produk dilarang (`repository/ambangproduk.go:21`). ⚠️ NILAI yang dicarinya (`OUTWARDRATEID`, b1416-b1604) tetap dibutuhkan Spreading — sumber penggantinya **OQ-M7** |
| `GetProductName` | RDBList | `Activity/ValidasiClaimReceived_Act.xml:306`, `Activity/ValidasiSTNC_Act.xml:333` | `internal/models/validasitanggal.go:29`, `internal/models/validasitanggal_test.go:9` | ✅ **ada** |
| `GetRateRetro` | RDBList | `Activity/SpreadingClaimLife_Act.xml:1750` | `internal/repository/ratelife.go` (`RateLife.Baca`, GILIRAN-17), `internal/services/spreading.go:107` | ✅ **ada** — pembaca sempit lima kolom (OQ-M7); nol pemanggil sampai `OUTWARDRATEID` ada |
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
| `FilterEmailKomiteWithLimit` | ReportDefinition | ~~`Activity/CreateKMTLife_Act.xml:437`~~ *(langkah 1, ter-remark b422)*, `Activity/GetListKomiteLife.xml:722` *(langkah 4, hidup)* | `internal/repository/roster.go:14`, `internal/repository/roster_test.go:14` | ✅ **ada** |
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
| `RejectOSClaimLife_Sec` | Section | `FlowAction/RejectOSClaimLife.xml:88` | `frontend/src/components/claimlife/DialogTolakOutstanding.tsx`, migrasi 021 | ✅ **dibangun — OQ-M5 ditutup 29-09-2026 (GILIRAN-17): Remarks di `T_CLAIMLF_JEJAK.KOMENTAR`**; dahulu: dialog Date b790, PIC read-only b975, Remarks b1687, Submit b3117 → `RejectOSClaimLife_Act` menambah baris `KomiteList` (`IDKomite="Claim Admin"`, `KomiteAproval=2`, `KomiteComment`=Remarks, b2173-b2263). Aplikasi menolak tanpa dialog; tempat simpan Remarks belum diputuskan |
| `RetroDetailClaimLife` | Section | `FlowAction/RetroClaimLife.xml:94` | `internal/repository/ratelife.go` (pembaca rate) | ⏸ **OQ-M7 ditutup 29-09-2026 (GILIRAN-17): `RATE_LIFE` dibaca sempit; panel menunggu `OUTWARDRATEID` (DBA)** — expand-pane grid treaty-year (`AdjustmentDetail` b10388-b10389): `RetroLifeList` Reinsurer/Currency/Percent Share/Claim Retro (b1325-b2498). Grid induknya pun tidak tampil; `AmbilSpreading`/`HitungSpreading` nol pemanggil produksi |
| `SearchPolicy_Section` | Section | `Section/SendtoAdmin_Section.xml:212`, `Section/SendtoMedical_Section.xml:215` | — | ✅ **tidak perlu** (residu / milik PremiumList) — "pemanggil" b212/b215 hanya entri pyPagesAndClasses: nol pyInclude, satu-satunya pyStreamName dirinya sendiri (b74/b77) |
| `SendtoAdmin_Section` | Section | `FlowAction/SendtoAdmin.xml:90`, `Section/SendtoMedical_Section.xml:227` | `frontend/src/assets/labels.claimlife.ts:399`, `frontend/src/components/claimlife/PanelPindahTahap.tsx:47` | 🔧 **celah → dibangun** — "Send Back to Admin?" b566 + `Submit` b1229 sebelum pindah tahap, untuk KETIGA tombol ber-local action `SendtoAdmin` (b21433, b20285, b20250); tombol maju tetap langsung |
| `SendtoMedical_Section` | Section | `FlowAction/SendtoMedical.xml:93` | `frontend/src/assets/labels.claimlife.ts:399`, `frontend/src/components/claimlife/PanelPindahTahap.tsx:47` | 🔧 **celah → dibangun** — "Send Back to Medical?" b577 + `Submit` b1267, tombol `Send Back to Medical` b20496 |
| `LinkService` | SystemSettings | `Activity/DeleteGoogleStorage_Act.xml:1304`, `Activity/GetLinkService.xml:6` | `internal/repository/linkservice.go:8`, `internal/services/antrean.go:318` | ✅ **ada** |
| `IsPEGAPROD` | When | `Activity/InsertGoogleStorage_Act.xml:1115`, `Activity/SaveOutStandingLife_Act.xml:11919` | `internal/services/efekkeluar.go:108`, `internal/services/efekkeluar_statik_test.go:222` | ✅ **ada** |
| `IsSendtoAdmin` | When | — *(nol pemanggil)* | `internal/services/tahap.go:38`, `internal/services/tahap_test.go:153` | ✅ **ada** |
| `IsSendtoMedical` | When | — *(nol pemanggil)* | `internal/services/tahap.go:39`, `internal/services/tahap_test.go:155` | ✅ **ada** |

## Sensus remark 28-09-2026 (GILIRAN-12 paket 2)

Setiap activity Claim Life dicetak beserta `pyStepsBlockName`; `//` = langkah ter-remark (beserta
anaknya). 51 activity, **8** ber-remark, **31** langkah `//`. Label lain (`EXIT`, `send`) adalah
label lompatan, bukan remark.

**Cara hitung** `[terverifikasi]` — dua cara, jendela = seluruh `Claim Life/Activity/*.xml`:

| Cara | Perintah | activity | ber-remark | langkah `//` |
| --- | --- | ---: | ---: | ---: |
| A — pengurai XML (expat), per `rowdata` ber-`pyStepPageReference` yang `pyStepsBlockName` tepat `//` | skrip sensus sesi (cetak activity, nomor langkah, bNNN) | 51 | 8 | 31 |
| B — grep mentah per berkas | `grep -o "<pyStepsBlockName>//</pyStepsBlockName>" "<berkas>" \| wc -l`, dijumlah | 51 | 8 | 31 |

Alat diuji lebih dulu pada item yang sudah diketahui: `Claim Life/Activity/SaveOutStandingLife_Act.xml`
— 8 langkah `//` (dibaca ulang pemeriksa independen GILIRAN-11) — kedua cara menjawab 8. Label
lompatan lain (`EXIT`, `send`, `AA`, …) bukan remark dan tidak terhitung.

⚠️ **Cakupan celah langkah hidup DINYATAKAN:** sensus ini mendaftar langkah hidup yang tidak ditiru
untuk activity **ber-remark** saja. Activity tanpa remark tidak diaudit ulang langkah demi langkah di
giliran ini; statusnya tetap seperti baris paritasnya (bila ada) — **belum** audit penuh.
 Nol langkah mati yang **ditiru** di kode; yang salah adalah **bukti**
yang mengutip baris mati — diralat di tempat.

| Activity | Langkah `//` (b) | Ditiru? | Keputusan |
| --- | --- | --- | --- |
| `SaveOutStandingLife_Act` | 11.3 b3495 · 11.9 b4632 · 11.11 b5009 · 12 b6178 · 13 b7074 · 14 b7293 · 15 b7512 · 23 b10649 | STNC & dokumen lengkap **sempat** — dibuang GILIRAN-11 (`9406c8f`) dan bl (`a170268`) | selesai |
| `CreateKMTLife_Act` | 1 b422 · 2 b619 · 3 b804 (+3.1) | tidak — tangga dari `GetListKomiteLife` 6.1 + langkah 4 b1188 | kutipan 866…1041 diralat (kode Komite, tiket 10) |
| `GetListKomiteLife` | 8 b1508 `Obj-Save` | tidak | — |
| `LoadDocumentLife_ACT` | 1.4 b913 (+1.4.1, 1.4.2) | tidak — yang hidup 1.3 b797 | "penyimpangan sadar URL" dicabut (`repository/klaimlife.go`, `PanelDokumenPeserta.tsx`, baris 13b); sumber URL = `DownloadDocumentClaim` langkah 2 b369 (`api.ts`) |
| `SaveAttachLife` | 1.1 b350 · 1.3 b718 · 1.4 b985 · 1.5 b1114 · 1.6 b1257 · 2 b1709 · 3 b1841 | tidak — yang hidup 1.2 + 1.7 | `models/dokumenbaru.go` menandai `Commit` ter-remark. **Celah hidup:** 1.2 b596 memakai ulang `Primary.DOCUMENT` (OQ-J) — dicatat di `services/unggahan.go` dan baris 10 |
| `SavePesertaClaim` | 3 b777 · 4 b922 · 7.2 b1863 · 7.3 b2039 · 7.4 b2279 · 7.6 b2606 · 9 b4875 | tidak (STNC 3/4, klaim ganda 7.x, `FinishAssignment` 9) | "precondition 1812" dicabut (baris 2, bab lanjutan 5); bukti `IsCheck` sudah dipindah ke b3631/b3919 (`9406c8f`) |
| `SendEmailKlaimLF` | 1 b333 · 2 b467 · 3 b658 (CC) | tidak — nol kode CC; `EfekEmail` stub | — |
| `setDetailClaim_act` | 2 b444 (+2.1 b478) | tidak (residu) | b478 dikoreksi di baris `setDetailClaim_act` dan OQ-B |

Activity Claim Life lain: **nol** remark — sensus ini tidak mengubah statusnya; celah langkah
hidupnya tetap seperti baris-baris di atas.

## 29-09-2026 (GILIRAN-15) — `Delete` baris adjustment TIDAK BERLAKU (OQ-N7 ditutup)

`Section/ClaimLifeDetailGCNM.xml` `Delete` b19120 → `deleteRow` b19130 (tampil bila `.PrintFaceClaim == '' &&
pyPosition=='ReasLifeSPV'` b19399) **tidak ditiru dan tidak dirender** — keputusan work owner. Sebabnya: `deleteRow` adalah
hapus fisik, ADR-U-0031 melarangnya di jalur pengguna, dan `T_CLAIMLF_ADJUSTMENT` tidak punya kolom penanda hapus untuk
padanan logisnya. GILIRAN-13 sempat merendernya sebagai tombol mati; dibuang, sebab tombol yang tidak pernah dapat
ditekan menjanjikan aksi yang tidak ada. Sunting sel baris adjustment pun tidak berlaku (butir br — semua sel
`Read-only`).

## 01-10-2026 — peserta versi polis terakhir: penyimpangan sadar dari `GetPesertaClaim_sql1` (OQ-N14)

`[keputusan work owner 01-10-2026]` `Claim Life/RDBList/GetPesertaClaim_sql1.xml` b85 *(nomor baris dari `sed -e
's/></>\n</g' GetPesertaClaim_sql1.xml | grep -n …`)*:

```sql
SELECT * FROM POOLDATA.M_LIFE_PREMIUM_DETAIL WHERE PL_NUMBER = {…PremiumListSummary.PL_NUMBER}
  AND CERTIFICATE_NO LIKE '%'||{SearchPolicyHolder.CARI2}||'%'  AND UPPER(NAME_OF_INSURED) LIKE '%'||{SearchPolicyHolder.CARI3}||'%'
```

**Nol saringan versi dan nol saringan `EDMSTATUS`** — di Pega setiap baris versi satu sertifikat (new business,
`/01`, `/02`, …) tampil, termasuk yang `Delete`/`Batal`. Sistem ini menampilkan dan menyalin **hanya baris versi
terakhir** tiap sertifikat, dan menyembunyikannya bila baris itu `Delete`/`Batal` (`repository/pesertapolis.go`:
`sqlCariPeserta`, `sqlAmbilPesertaKlaim`). Versi = angka sesudah `<PL_NUMBER>/` di `PL_NUMBER_EDM` (rumus penulisnya
`Endorsement Life/RDBList/Generate_NoEndorsmentLife.xml` b84 `NOPOLIS||'/'||CARI14`); seri → `TGL_INPUT` terbaru, lalu
`ID` terbesar secara angka. Baris `11` dan `12` tabel aksi serta baris `GetPesertaClaim_sql1` sensus di atas ditandai.
Rincian, SQL lama/baru, dan sepuluh kasus ujinya: tiket 02, bab *"Keputusan bertanggal — 1 Oktober 2026 (OQ-N14)"*.

## 10-10-2026 — Register banyak premium list: satu case per PL number + mata uang (penyimpangan sadar)

Keputusan work owner 10-10-2026: *"saat user ceklis 2 orang dari PL A, lalu 4 orang dari PL B, claim membuat 2 case
(CLMLF) berdasarkan PL number"*. **Menyimpang dari Pega**: `InputRegisterClaimLife` mendaftarkan SATU polis per register
(`SaveOutStandingLife_Act` / `SavePesertaClaim` atas satu `PolicyDataLife`). Kini:

- Layar `RegisterKlaim.tsx`: pilihan peserta **bertahan antar pencarian** (kunci PL + polis + sertifikat), blok *Selected insured*
  dikelompokkan per PL + mata uang dengan kalimat *"{n} claim case(s) will be created"*; penolakan pilihan campuran polis /
  mata uang dibuang; hasil = daftar case (nomor, PL, mata uang, jumlah peserta).
- `POST /api/klaim-life` menerima `kelompok[]` (bentuk lama satu kelompok tetap diterima). Server **membentuk ulang** case dari
  baris sumber (`AmbilUntukKlaim`): kunci = `PL_NUMBER` + `CURRENCY` sumber; kelompok dan mata uang klien diabaikan.
- Cross-check polis: `NO_POLIS` header diturunkan dari `POLICY_NO` baris sumber; seluruh peserta satu case wajib sepolis;
  `nomorPolis` klien yang berbeda = **422**. Satu PL + mata uang berisi dua polis = 422 (satu case = satu polis).
- **Semua atau tidak sama sekali**: validasi seluruh kelompok sebelum nomor diambil; N nomor dalam satu transaksi pendek
  (AC tiket 02: kuncian penghitung tidak ditahan); semua case dalam SATU transaksi penulisan. ⚠️ Penghitung tidak ikut
  rollback: bila transaksi penulisan gagal, N nomor terbakar (urutan berlubang) - dinyatakan di `Pendaftaran.Daftar`.
- Type dan kode bisnis sama untuk semua case (sampai layar Register dirombak). Batas 50 case / 500 peserta per daftar.
- Ganda di dalam permintaan (sertifikat sama pada PL sama) ditolak 400. Pendaftaran tetap tanpa jejak dan tanpa pemeriksa
  klaim ganda warisan (`gandawarisan.go` berjalan di `Save to RNM` langkah 11) - tidak berubah.
- Uji: `services/pendaftaran_banyak_test.go`, `handlers/register_banyak_test.go`, `frontend/pilihanPeserta.test.tsx`.

## 10-10-2026 — Register menolak peserta yang masih aktif di klaim lain (aturan in-app, penyimpangan sadar)

Keputusan work owner 10-10-2026: *"peserta berada di case lain, dan di case itu peserta masih outstanding atau sudah di
accept. tapi jika peserta di case lain di reject. peserta itu bisa di input lagi di case yang sama atau di case lain"*.
**Menyimpang dari Pega**: Pega memeriksa klaim ganda hanya di `Save to RNM` langkah 11 atas tabel datar warisan
(`gandawarisan.go`, TIDAK berubah - baris aplikasi berstatus NULL sengaja dilewati di sana). Kini Register juga memeriksa:

- `repository.PesertaAktifDiKlaimLain` (`pesertainapp.go`) atas `T_CLAIMLF_PREMIUMLIST_DETAIL` + `T_CLAIMLF_ADJUSTMENT`, kunci
  `PL_NUMBER` + `CERTIFICATE_NO` (TRIM) - bukan `SOURCE_ID`, yang berganti per versi polis (OQ-N14). Satu kueri per PL.
- Peserta AKTIF bila status baris pesertanya bukan Ditolak ATAU ia punya baris adjustment bukan Ditolak (tombol `Add`
  membuka putaran baru tanpa mengubah baris peserta). Rejected di semua tingkat dan baris tercabut (`STS_HAPUS`) boleh
  didaftarkan lagi. Status header tidak menentukan. Kode dari `kontrak` (diikat, tidak literal).
- Diperiksa SEBELUM nomor diambil; seluruh bentrok dilaporkan sekaligus, **409**:
  `Certificate No "X" (PL Y) is already in claim CLMLF-… (Outstanding|Accepted)`. Multi-PL dengan satu bentrok = semua batal.
- ⚠️ Kinerja: tabel peserta klaim hanya ber-index `CLAIM_ID`; index `(PL_NUMBER, CERTIFICATE_NO)` butuh migrasi (belum).
- Uji: `repository/pesertainapp_test.go`, `services/pendaftaran_banyak_test.go`, `-tags db` `TestDaftarPesertaAktifDiKlaimLain`.

## 10-10-2026 — Data polis otomatis dengan kunci PL; Type dan kode bisnis dari polis (penyimpangan sadar)

Keputusan work owner 10-10-2026: *"buat data polis ini otomatis terisi saat menceklis peserta, ambil dari data polis
peserta itu; kolom type dan kode bisnis hilangkan saja"*.

- **Kunci polis = PL NUMBER peserta** (`services/kuncipolis.go` `RingkasPolisKlaim`): `T_PREMIUM_LIST.NO_POLIS` berisi PL
  number, sedangkan `POLICY_NO` peserta adalah nomor polis CEDING (`docs/spec.md` b359). Fallback `POLICY_NO` untuk polis
  warisan. Dipakai panel Data Polis Register dan Outstanding, ambang MAXEXPIREDCLAIM (detail), `AmbangKlaim`, Save to RNM,
  pengisian `CEDINGCO` cermin warisan, dan Register.
- Arti kolom tidak berubah: `T_GENERAL_CLAIM.POLICY_NO` dan `T_CLAIMLF_PREMIUMLIST_DETAIL.POLICY_NO` tetap POLICY_NO ceding.
- Register: isian **Type** dan **Kode bisnis** dihapus. Per case Type = `PolisRingkas.Type`, kode bisnis = `BusinessCode`
  (sama dengan Save to RNM); polis tak ditemukan / tanpa Type / tanpa kode bisnis = 422 menyebut PL, nol nomor.
  **Menyimpang dari Pega**: di `InputRegisterClaimLife` `.PolicyDataLife.Type` adalah dropdown (b9104).
- Aturan "satu case satu POLICY_NO" dan cross-check nomor polis klien dibuang - case = PL + mata uang saja.
- Panel: kotak catatan teknis (`.PolicyDataLife.*`, daftar medan tanpa kolom) dihapus; polis tak ditemukan =
  "Policy data not found for PL X.".

- Uji: `services/pendaftaran_banyak_test.go`, `repository/gandawarisan_test.go`, `frontend/pilihanPeserta.test.tsx`.

## 10-10-2026 (sore) — Data polis klaim dari IDPEGA peserta (MENGGANTIKAN kunci PL bab sebelumnya)

Keputusan work owner 10-10-2026: *"ambil ke tabel T_PREMIUM_LIST berdasarkan idpega peserta nya saja"*.

- Header polis = baris UTAMA `T_PREMIUM_LIST` ber-`ID = ID_PEGA = IDPEGA` peserta (`M_LIFE_PREMIUM_DETAIL.IDPEGA` = pengenal
  kasus versi polis `NBLF-`/`EDML-`). Kunci PL / POLICY_NO / NO_POLIS dan fallback-nya DIBUANG dari Claim Life.
- PremiumList: `RingkasDariIdPega` (+ `GET /api/polis-life/ringkas?idpega=`); `Ringkas(NO_POLIS)` lama tidak berubah.
  `kontrak.PembacaPolis` (inti) tidak diubah - Claim Life menegaskan antarmuka sempit `PembacaPolisIdPega`.
- Claim Life membaca IDPEGA bersama peserta (`kolomSalin`, `Cari`); `T_CLAIMLF_PREMIUMLIST_DETAIL` TANPA kolom IDPEGA,
  jadi detail / Save to RNM / ambang / CEDINGCO cermin menurunkannya dari baris sumber lewat `SOURCE_ID`.
- Register per case: Type dan kode bisnis dari header IDPEGA (versi terbaru bila peserta case ber-IDPEGA beda);
  header tak ditemukan / peserta tanpa IDPEGA = 422, nol nomor. Panel: "Policy data not found for this insured (<IDPEGA>)."

## 10-10-2026 (malam) — IDPEGA disimpan (migrasi 025) dan polis wajib Resolved-Completed

Keputusan work owner 10-10-2026: (1) *"iya tambahkan"* - kolom `T_CLAIMLF_PREMIUMLIST_DETAIL.IDPEGA` (migrasi 025, VARCHAR2(50)
= lebar sumber) + isi ulang baris lama lewat `SOURCE_ID`; Register menulisnya, tahap lanjut membacanya (cadangan lewat
`SOURCE_ID` hanya bila NULL). (2) *"tidak, polis harus sudah resolve completed"* - Register menolak (422, semua case,
nol nomor) header IDPEGA yang status efektifnya bukan `Resolved-Completed`: `T_WORK_POLIS.STATUS_WORK` kasus NB, selain
itu `T_PREMIUM_LIST.STATUSS` (ditulis EDM Confirm). Find Insured tetap menampilkan peserta; panel Data Polis memperingatkan
"Policy not Resolved-Completed (status: X)". **Menyimpang dari Pega**: Pega tidak menyaring status polis di Register.

## 10-10-2026 — Tampilan Register dan panel Data Polis dirapikan (tampilan saja)

Keputusan work owner 10-10-2026: *"rapihkan tampilan, agar lebih enak di lihat; tanggal menggunakan dd/mm/yyyy; system
reinsurance, premium method tidak perlu di tampilkan; bagian Selected insured dirapikan, menggunakan tabel dan user friendly"*.

- Panel Data Polis (Register DAN Outstanding): System Reinsurance (b9301) dan Premium Method (b9694) disembunyikan - datanya tetap;
  tanggal dd/mm/yyyy (`tanggalTampil`, di atas `formatDate` inti); grid label-di-atas-nilai; baris teknis "Polis … versi …"
  diganti chip versi ("New business" / "Endorsement /NN"); peringatan sebagai `alert alert--warn`. **Menyimpang dari Pega**:
  himpunan medan section tidak lagi utuh di layar (penjaga cacah medan diperbarui dengan alasan).
- Register: judul "Register Claim Life"; kartu pencarian (Premium List No, Certificate No, Name of Insured, Search); tabel hasil
  berbaris-klik dengan sorot; Selected insured = tabel per PL + mata uang (No, Certificate No, Name of Insured, Policy No, Remove),
  ringkasan + Clear selection + Register Claim sejajar kanan; label layar seragam bahasa Inggris (`REGISTER_BANYAK`, bukan korpus).
  Logika dan backend tidak berubah.

## 10-10-2026 — Register: kolom Date of Birth menggantikan Policy No; konfirmasi sebelum Register Claim

Keputusan work owner 10-10-2026: *"policy no ini ganti saja dengan dob (dd/mm/yyyy); dan saat register claim di klik, ada
konfirmasi dahulu"*.

- `GET /api/peserta-life` membawa `tanggalLahir` (`TO_CHAR(DOB, 'YYYY-MM-DD')` dari baris versi terpilih M_LIFE_PREMIUM_DETAIL,
  kebal NLS); tabel hasil dan Selected insured menampilkan Date of Birth dd/mm/yyyy ("—" bila kosong). Nomor polis ceding
  tetap di respons, tidak ditampilkan. Catatan: aturan "nama/DOB tidak keluar dari basis data" (gandawarisan) berlaku
  untuk pencocokan klaim ganda; layar Find Insured sudah menampilkan nama, dan DOB kini ditampilkan atas keputusan ini.
- Klik Register Claim membuka modal konfirmasi (Modal inti): jumlah case + tabel PL/Currency/Insured, "Continue?",
  Cancel / Register Claim; saat proses tombol nonaktif "Registering…" dan modal tak dapat ditutup; galat tampil di layar.

## 10-10-2026 — Inbox Claim Life dirapikan; tab Input Register disembunyikan selama kosong

Keputusan work owner 10-10-2026: *"rapikan tampilan inbox claim life ini, dan tab input register hilangkan saja, karena sudah
ada tombol register"*.

- Tab Input Register TIDAK tampil selama antriannya kosong; ia muncul lagi otomatis bila ada kasus di sana, sebab
  `Send Back to Register` (b21404, PanelPindahTahap) masih menaruh kasus di tahap itu. Tahap, rute, dan model tidak diubah.
  Tab bawaan = Outstanding Claim. **Menyimpang dari Pega**: Pega selalu menampilkan keempat antrian.
- Label tab Inggris saja (subjudul Indonesia dibuang); badge jumlah kasus di setiap tab (halaman 1 tiap antrian).
- Register = btn--primary, Export xlsx = btn--ghost (tetap nonaktif saat nol baris - sengaja, butir bg).
- Waktu buat dd/mm/yyyy HH:mm; kosong "No cases in this queue."; paginasi komponen `Halaman` inti (seperti inbox lain).

## 10-10-2026 — Perbaikan sesudah uji DEV Register (CLMLF-, nomor di Save to RNM, detail, status inbox)

- Case ID baru `CLMLF-xxxxxx` (keputusan work owner; dulu `CLM-`); pengenal lama tetap terbaca (tidak ada pengurai awalan).
- Nomor klaim lahir di Save to RNM, bukan Register (RALAT tiket 02); inbox menampilkan "—" sampai bernomor.
- Detail klaim: handler kini mencatat galat aslinya ke log (dulu tidak) dan menjawab "Failed to load the claim."; galat
  membaca polis/ambang produk untuk penanda MAX CLAIM RECEIVED menjadi alasan, tidak lagi menggagalkan layar Detail.
- Work Status kotak masuk berbahasa Inggris dari STS_REJECT header: NULL/0 Outstanding, 1 Accepted, 2 Rejected, selain itu Unknown.
- Modul Claim Life memakai BahasaUI "en" seperti Endorsement Life (paginasi dan teks komponen bersama berbahasa Inggris).

## 10-10-2026 — Sesudah Register berhasil langsung ke Inbox Claim Life

Keputusan work owner 10-10-2026: *"saat klik tombol register claim ini, langsung keluar dari menu atau langsung ke inbox
claim life yang paling depan"*.

- Berhasil: modal tertutup, pilihan dikosongkan, rute pindah ke Inbox Claim Life tab tampil pertama (Outstanding Claim).
  Inbox dipasang ulang, jadi jumlah badge dan isi tab dimuat ulang.
- Inbox menampilkan alert hijau "{n} claim case(s) registered: CLMLF-…" yang bisa ditutup (Close). Case ID dikirim lewat
  state modul di rute.tsx (bukan localStorage); hilang saat ditutup, saat membuka Register lagi, atau saat muat ulang halaman.
- Gagal: tetap di Register, modal tertutup, alert galat, pilihan peserta utuh.
- Tabel hasil (HasilKasus) di layar Register dihapus. **Menyimpang dari Pega**: Pega tetap di layar flow sesudah Register.

## 10-10-2026 — OS_AKSEPTASI_KLAIM_LIFE ditulis saat Save to RNM, bukan saat Register

- Register tidak lagi menulis baris datar OS_AKSEPTASI_KLAIM_LIFE. Save to RNM menulisnya (UPDATE bila ada, INSERT bila
  belum) di transaksi utamanya, sesudah nomor klaim lahir. **Sama dengan Pega** (`InsertJsonKlaimLife_sql` di
  `SaveOutStandingLife_Act` 22.1.3.1). Penyimpangan lama (cermin ditulis saat Register) dihapus.
- Klaim ganda: baris datar klaim yang belum Save to RNM memang belum ada. Dulu barisnya ada dengan STS_REJECT NULL dan
  dilewati pemeriksa, jadi hasilnya tidak berubah.
- Hapus klaim / cabut peserta sebelum Save to RNM: nol baris datar bukan galat.

## 10-10-2026 — Detail klaim "bermata uang beragam: IDR dan 10026"

- Perbaikan pembanding total peserta (`HitungTotalPeserta`): KODE mata uang (CURRENCY / mata uang peserta, mis. IDR)
  dibandingkan dengan kode, dan PENGENAL (CURRENCY_ID, mis. 10026 = IDR) dengan pengenal. Dulu kode peserta dibandingkan
  dengan CURRENCY_ID sehingga setiap case hasil Register gagal dibaca. Mata uang yang benar-benar berbeda tetap ditolak.
- Penulisnya sudah benar: Register menulis CURRENCY = kode dan CURRENCY_ID = NULL, lalu pengenal diisi saat baca
  (`GetCurrencyID` / lengkapiPengenalMataUang). Data DEV tidak perlu dibetulkan.

## 10-10-2026 — Layar Outstanding Claim dirombak mengikuti Pega

Permintaan work owner 10-10-2026: *"kenapa tidak ada detail peserta nya seperti yang di pega?"* (screenshot CLMLF-9293).
Korpus: `kelvin/InboxClaimLife (Done)/Section/InputOSClaimLife.xml` (bNNN = nomor baris berkas itu);
`FlowAction/OSClaimLife.xml` → section itu.

| Bagian Pega | Bukti | Aplikasi |
| --- | --- | --- |
| Judul "Claim Life (<Case ID>)", Claim No baca saja | b939/b958 `.ClaimData.PremiumListSummary.CLAIM_NO` | kepala halaman; "—" sampai Save to RNM |
| Tab Register Claim / Policy Details / Claim Details | b1956, b5909, b14576 `pyTitle` | `StripTab` inti, urutan sama |
| PL Number (baca saja) | b3054 | dari peserta pertama |
| Catastrophe (radio, nonaktif bila CLAIM_NO terisi) | b3475, b3504, b3515; bawaan "Non-Catastrophe" b3728 | `T_GENERAL_CLAIM.STS_KATASTROFE`; server 409 bila diubah sesudah bernomor |
| Is KPR (radio) | b3781, b3810 | `IS_KPR`; ⚠️ nilai pilihan **tidak ada di korpus** — Yes/No sementara, perlu konfirmasi |
| STNC Treaty (baca saja) | b4726, b4757 | `STNC_CLAIM` |
| Catastrophe Note (tampil bila Catastrophe, baca saja bila CLAIM_NO terisi) | b4905, b4934, b5059 | `KATASTROFE_NOTE`; server 409 bila diubah sesudah bernomor |
| Policy Details | b5909–b12679 | `PanelDataPolis` (keputusan WO sebelumnya: tanpa System Reinsurance / Premium Method, tanggal dd/mm/yyyy) |
| Grid DetailPolisLife (View Detail + grid peserta polis) | b13267 | **BELUM ditiru** — milik modul PremiumList; perlu rute baca-saja lintas modul |
| Participant Details: PL NUMBER, NAME OF INSURED, DATE OF BIRTH | b15524, b15924, b16101, b16256, b16406 | nama dan DOB **dibaca saat baca** dari `M_LIFE_PREMIUM_DETAIL` (SOURCE_ID), tidak disalin ke tabel klaim |
| Accept (`.IsAccept`) + SELECT ALL | b17582, b17597, b16609 → `SelectAllClaimLife_act` b247/b429 | **migrasi 026** `IS_ACCEPT` (disetujui work owner 10-10-2026); penjungkit `inti/frontend/lib/pilihSemua.ts` |
| DELETE (tampil bila CLAIM_NO kosong) | b16839, b17865, b18082 → `DeletePesertaClaimLife` | cabut peserta (penanda STS_HAPUS). **Menyimpang**: aplikasi menambah konfirmasi "Remove this insured from the claim?" (Pega tanpa konfirmasi, b18021) |
| Baris dibuka → rincian peserta | b18252 `ViewClaimDetailLifeGCNM` → `ClaimLifeDetailGCNM` | `RincianPeserta` — blok per peserta layar Detail **dipindah** ke komponen bersama |
| Save | b19413 | `PUT /api/klaim-life/{id}/outstanding/isian` (Accept + medan Register Claim), satu transaksi |
| Save to RNM (nonaktif sesudah tersimpan) | b21102, b21095 `pyWorkPage.Save = 1`, b21126 | padanan bendera: `sudahSaveRNM` (OQ-N1, tanpa kolom) |
| Send Back to Register (tampil bila CLAIM_NO kosong) | b21404, b21433, b21699 | server 409 bila sudah bernomor |
| Send to Medical Check (nonaktif bila CLAIM_NO kosong; tampil bila QP/QR) | b21839, b21832, b21863, b22123 | server 409 tanpa CLAIM_NO atau Type selain QP/QR — **menutup risiko Komite** (klaim tanpa Save to RNM tidak punya baris datar) |
| Close Claim (selalu tampil) | b22808, b22837 | `PanelTutupKlaim` — logika layar Detail dipindah ke komponen bersama |

- **Gerbang dokumen Save to RNM** langkah 3 kini membaca `IS_ACCEPT` (`SaveOutStandingLife_Act` b1181
  `.IsAccept=="true"`). Dulu `IS_CHECK` dipakai sebagai padanan terdekat — ralat. Akibat: peserta yang belum di-Accept
  tidak lagi dituntut dokumen saat Save to RNM (sama dengan Pega).
- **Tombol "Detail klaim" dihapus** untuk kasus tahap Outstanding: seluruh isinya ada di tab Claim Details. Kasus tahap lain
  yang terbuka di layar ini diarahkan ke layar Detail lewat catatan.
- Nomor akseptasi tetap tidak dibuat Save to RNM (`SaveOutStandingLife_Act` hanya CLAIM_NO).

## 11-10-2026 — Is KPR, View Detail peserta polis, dan temuan Save to RNM

Keputusan work owner 11-10-2026.

- **Is KPR**: pilihan **KPR / Non-KPR**, teks disimpan persis (pola Catastrophe / Non-Catastrophe). Server menolak nilai lain
  (400), KECUALI nilai yang sama dengan yang sudah tersimpan: data warisan di luar dua pilihan itu ditampilkan apa adanya
  ("Stored value: …", tanpa radio tercentang) dan Save yang tidak menyentuhnya tidak gagal. Data DEV belum dapat dibaca dari
  sini (tanpa DSN); belum diketahui apakah ada nilai lain.
- **View Detail peserta polis** (`Section/DetailPolisLife.xml`, InputOSClaimLife b13267): tombol View Detail b631 (hilang
  sesudah ditekan, b792), grid `TempDetail.pxResults` b4537 dengan kolom No., POLICY NO, CERTIFICATE NO, NAME OF INSURED, SEX,
  DATE OF BIRTH, AGE, PLAN, BEGIN DATE, EFFECTIVE DATE, EXPIRED DATE (b4700–b6082), 500 baris per halaman
  (`LoadDataPeserta_Act` b572/b619, `getMaxPagination_sql`). Rute baca-saja claimlife
  `GET /api/klaim-life/{id}/peserta-polis?halaman=N`; PremiumList tidak disentuh (rute premiumlistlife membaca tabelnya
  sendiri, bukan `M_LIFE_PREMIUM_DETAIL` yang dibaca Pega). Paginasi NextPrev ⟪ ⟨ ⟩ ⟫ diganti komponen `Halaman` inti.
  **Menyimpang**: (1) Pega menyaring `PL_NUMBER` saja (`GetPesertaClaim_sql`, semua versi polis termasuk Batal/Delete);
  di sini `PL_NUMBER` (INDEX4) + `IDPEGA` = versi polis yang sama dengan header klaim, tanpa saringan status EDM seperti
  Pega; (2) Pega memberi `ROWNUM` sebelum `ORDER BY` sehingga halamannya tidak terurut — di sini urutan Pega diterapkan
  sungguh-sungguh (+ID pemutus seri); (3) Pega menambah 7 jam ke tanggal untuk tampilan (b1097 dst) — di sini DATE dibaca
  sebagai tanggal kalender, hasil tampil sama.
- **Temuan (belum diubah)**: wadah Participant Details berkondisi `pyContainerVisibleWhen pyWorkPage.Save = 1`
  (InputOSClaimLife b15490), dan `pyWorkPage.Save` hanya ditulis `SaveOutStandingLife_Act` (`""` b8549 langkah 21, `1`
  b11367 langkah 24). Di Pega grid peserta karena itu baru tampil SESUDAH Save to RNM; aplikasi menampilkannya selalu
  (permintaan work owner 10-10-2026). Menunggu keputusan.
- **Temuan (belum diubah)**: tombol Edit Date (ClaimLifeDetailGCNM b14115) tampil hanya bila Admin DAN CLAIM_NO terisi
  (b14357), padahal keempat tanggal baca-saja bila CLAIM_NO terisi (b11339, b11577, b12857, b13997; EditDate b1000/b1313/
  b1550/b1788). Aplikasi: tanggal dapat diubah sebelum Save to RNM, terkunci sesudahnya (OQ-M1).

## 11-10-2026 — Layar Outstanding dirapikan mengikuti screenshot Pega; rincian peserta dan isian adjustment

Keputusan work owner 11-10-2026: *"rapikan tampilannya dengan referensi dari pega (4 screenshot CLMLF-9293). Tombol send
back to register, check readiness claim, dan close claim hilangkan saja. Saat pesertanya diklik bisa menampilkan detail
pesertanya dan bisa isi adjustment-nya seperti di pega."*

- **Tombol bawah**: hanya Save to RNM dan Send to Medical Check. Send Back to Register (b21404), Check readiness dan Close
  Claim (b22808) **tidak ditampilkan** di layar Outstanding. **Menyimpang dari XML** (keduanya ada di InputOSClaimLife);
  rute/aksi server TETAP (Close Claim masih di layar Detail dan tahap Claim Analis; Pindah Outstanding → Input Register
  masih diterima server bila belum bernomor).
- **Register Claim**: form dua kolom label-kiri; Catastrophe radio tegak; Is KPR bertanda wajib (*).
- **Policy Details**: tiga kolom seperti screenshot; System Reinsurance (`.TypeCeding` b7102) dan Premium Method
  (`.ProRateType` b7495) **ditampilkan di Outstanding** (screenshot 11-10-2026) — panel Register tetap tanpa keduanya
  (keputusan 10-10-2026). Type tampil sebagai KODE: daftar nama (`pyListSource associated` b6905) tidak ada di korpus.
- **Participant Details**: grid bernomor, baris terbuka tersorot, saringan nama, SELECT ALL berupa tautan, Save kecil di kanan
  bawah. Klik baris membuka rincian DI BAWAH baris itu (b18252) dan menutup baris lain.
- **Rincian peserta** (`ClaimLifeDetailGCNM.xml`, label VERBATIM): tiga kolom urut screenshot; ADMIN NOTES (b2588)
  **migrasi 027** `ADMIN_NOTES` (disetujui work owner 11-10-2026), dibekukan bila peserta diputus (b2628); CLAIM ANALYST =
  `NOTES` baca saja; MAX EXPIRED CLAIM = `MAXEXPIREDCLAIM` produk (ValidasiClaimReceived_Act b540) saat baca; SEX/PLAN
  dibaca dari baris sumber. **Belum ditiru**: kotak LOCATION OF DEATH (ada di screenshot, TIDAK ada di korpus, tanpa kolom —
  baca saja "—"); STNC DATA dan MAX DATA RECEIVED (`ValidasiSTNC_Act` b577/b627); OutStanding Date (`.ADJUSTMENT_DATE`, ditulis
  SaveOutStandingLife_Act b10403, tanpa kolom). Diagnosa hanya di tahap Medical Check (b4908 `pyPosition` MedicalAdvisor).
  Sub-tab Adjustment (b16543) | Document (b23293).
- **Isian adjustment** (`AdjustmentDetail_Section`): keputusan work owner 11-10-2026 — SEBELUM Save to RNM (klaim belum
  bernomor, baris belum berstatus) **Claim RNM** (`.CLAIM_GROSS` b3012 = kolom `CLAIM_AMOUNT`) dapat diisi; **Percent Claim**
  dan **Claim Paid** TIDAK diisi pemakai — diturunkan saat baca (Percent = Claim RNM / Share RNM × 100, pecahan 5 desimal
  seperti b702; Claim Paid = Claim RNM), tanpa kolom; Claim RNM > Share Nusantara Re ditolak 422 dengan kalimat b526. SESUDAH
  Save to RNM seluruhnya baca saja. **Menyimpang dari XML** (`pyEditOptions Read-only` di semua medan; keputusan br/OQ-N8
  29-09-2026 kini diubah untuk Claim RNM saja). Pengecualian kode bisnis L12/L13/L14/L18 (`CountClaimAmountLife_Act` b972)
  **tidak** dipakai karena work owner menyatakan batasnya mutlak — perlu konfirmasi. Spreading ulang sesudah Claim RNM berubah
  (`CountClaimAmountLife_Act` langkah 5 → `SpreadingClaimLife_Act`) **belum** dipanggil.

## 11-10-2026 (lanjutan) — Tanggal peserta diisi pemakai, Claim Analyst, View Detail dicabut, adjustment sesuai Pega

Keputusan work owner 11-10-2026 (lima screenshot).

- **Lima tanggal peserta** (DATE OF LOSS*, CLAIM RECEIVED DATE*, DOCUMENT COMPLETE DATE, MAX DATA RECEIVED, CONFIRMATION
  DATE*) menjadi isian langsung di kolom kanan rincian peserta, sebelum Save to RNM; sesudahnya baca saja. Disimpan tombol
  **Save grid** (satu transaksi bersama Accept / Admin Notes / Catastrophe) — satu jalur, satu tombol. DOL tetap melewati
  `ValidasiDOL` (ValidasiDOL_Act); gerbang OQ-M1 (`BolehUbahTanggalKlaim`) dan CLAIM_NO kosong di layanan dan di SQL.
  MAX DATA RECEIVED: **migrasi 028** `MAX_DATA_RECEIVED DATE` (disetujui work owner 11-10-2026); di Pega dihitung
  `ValidasiSTNC_Act` b577 — kini diisi pemakai (**menyimpang**). Tombol Edit Date tidak tampil di Outstanding (masih di layar
  Detail). `ValidasiClaimReceived_Act` tetap bukan gerbang (OQ-M9; penanda saat baca).
- **Save to RNM wajib**: Claim Received Date dan Confirmation Date kosong ditolak 422 ("Claim Received Date is required for
  insured number N", "Confirmation Date is required for insured number N") — SEMUA peserta, seperti gerbang 11.x. DOL kosong
  sudah ditolak langkah 11.10 dengan kalimat XML ("DOL cannot be blank or outside the valuation period No N").
- **View Detail** dicabut dari tab Policy Details. Rute baca `GET …/peserta-polis` tetap (tanpa pemakai layar; dapat dibuang).
- **CLAIM ANALYST** hanya tampil di tahap Claim Analis (Pega: `pyPosition` ReasLifeSPV b15474); disembunyikan bagi Admin.
- **Grid adjustment**: Status sebagai teks ('' → "—", 0 → "OutStanding", 1/2 → "Accepted"/"Rejected"; daftar Pega tidak ada
  di korpus); **OutStanding Date** diturunkan dari `T_CLAIMLF_JEJAK` — `MIN(WAKTU)` baris itu dengan `KE='0'` dan `DARI` NULL
  (transisi '' → '0' yang direkam Save to RNM), tanpa migrasi; "—" sebelum Save to RNM.
- **Isian adjustment** tiga kolom seperti screenshot, angka bergaya Indonesia lewat `formatNumber` inti (94.695,705). Data
  bank: Payable To baca saja (XML); Name of Bank, Branch of Bank, Account No **isian teks** sebelum Save to RNM (kolom 004/011)
  — master autocomplete `Int-BANKACCOUNT` belum diketahui (keputusan work owner: teks dulu). Grid Treaty Type | Claim RNM dari
  spreading (`.TREATYTYPENAME`/`.RetrocadedShare`); grid komite tampil "No items" (data modul Komite belum disambungkan).
  Reject Outstanding dan Send To Komite di kaki isian dengan kondisi XML (b15399). Rute isian baris:
  `PUT …/adjustment/{adjId}/isian` (Claim RNM + bank).

## 11-10-2026 (lanjutan 2) — Total peserta dibuang, bagian Spreading, satu tombol Save, kerapian tata letak

- **Total peserta**: wadah enam total di ClaimLifeDetailGCNM berkondisi `pyContainerVisibleWhen 1=2` (b20403) — di Pega TIDAK
  pernah tampil di tahap mana pun. Layar Outstanding membuangnya (prop `tanpaTotal`); layar Detail (KlaimLife, dipakai
  kasus tahap lain) MASIH menampilkannya sampai diputuskan. Perhitungan backend (`HitungTotalPeserta`, gerbang mata uang) tetap.
- **Spreading** menggantikan grid Treaty Type dan grid komite di isian adjustment: bentuk tab Spreading PremiumList Life
  (disalin bentuknya) — per lapisan (Reins Type, Retroceded Share) lalu per reinsurer (% Share, Amount, Rate, Commision %,
  Ovr Comm %, Gross, Net), dan rekap per reinsurer (Currency, Premium Spreaded Gross/Net) yang dijumlahkan BACKEND
  (Money.Add). Sumber: T_CLAIMLF_ADJUSTMENT_SPREADING + T_CLAIMLF_ADJ_SPREADING_RETRO apa adanya. **Spreading ulang saat
  Claim RNM berubah BELUM disambungkan**: `services.HitungSpreading` (SpreadingClaimLife_Act) ada dan teruji tetapi tanpa
  pemanggil — belum ada pembaca TREATYYEAR_LIFE / RETROCESSIONLIFE, kapasitas IDR/USD treaty-year tidak punya kolom di DEV
  (`ErrKapasitasBelumDiketahui`), dan pemetaan IDUSEDBY rate belum diputuskan. Kosong → "No spreading yet. It is calculated
  when Claim RNM is saved." Grid komite dibuang dari isian.
- **Satu tombol simpan**: tombol "Save to Outstanding" dicabut; Claim RNM + data bank ikut `PUT …/outstanding/isian`
  (`IsianOutstanding.adjustment`, `services.PeriksaIsianBaris`) dalam satu transaksi bersama Accept, Admin Notes, lima
  tanggal, Catastrophe/KPR. Rute `PUT …/adjustment/{adjId}/isian` dicabut. Aturan tetap: Claim RNM <= Share RNM (422), hanya
  sebelum Save to RNM (409).
- **Kerapian**: sel rincian di grid kembali membungkus (inti `.inbox__tabel td` nowrap), label 24ch nowrap, nilai kode
  nowrap, teks panjang membungkus, isian tanggal/angka berlebar ch, radio satu baris sejajar (radiogroup menggantikan
  fieldset/legend), angka rata kanan, tabel diagnosa/dokumen dalam pembungkus gulir. Dijaga `gaya.test.ts`.

## 11-10-2026 (lanjutan 3) — Spreading diambil dari PremiumList Life; total peserta dibuang dari semua layar

- **Spreading** (keputusan work owner: "datanya ambil saja dari premium list, tidak perlu perhitungan lagi"): bagian Spreading
  rincian peserta membaca `T_PREMIUM_LIST_SPREADING` + `_RETRO` (milik PremiumList Life, baca saja; SQL disalin bentuknya dari
  polis_spreading.go, paketnya tidak diimpor) lewat `GET …/peserta/{pesertaId}/spreading-polis`, dimuat saat rincian dibuka.
  Kunci: IDPEGA peserta (= `PREMIUM_LIST_ID`, IDX_PLD_PL) + `TRIM(CERTIFICATE_NO)`, tanpa baris EDM `Delete`; lalu
  `DETAIL_ID` (IDX_PLS_DETAIL) dan `SPREADING_ID` (IDX_PLSR_SPR) — seluruhnya berindeks. Bukan lewat ID: `M_LIFE_PREMIUM_DETAIL.ID`
  (SOURCE_ID) adalah sequence sendiri. Versi EDM menyalin spreading (endorsementlife edm_buat.go, AC 59). Rekap per reinsurer
  dijumlahkan backend (Money.Add) dari baris peserta itu (rekap tersimpan PremiumList per POLIS, bukan per peserta). Kosong →
  "No spreading data for this insured in the policy." Tanpa hitung, tanpa tulis.
- **Spreading klaim lama**: `services.HitungSpreading` dan `PohonKlaim.AmbilSpreading` kini tanpa pemanggil produksi;
  `PohonKlaim.Simpan` (Register) masih punya jalur INSERT spreading/retro klaim, tetapi Register tidak pernah mengisi
  `Baris[].Spreading` (BarisPendaftaran) — praktis nol baris ditulis. Dibiarkan, menunggu keputusan.
- **Total peserta** dibuang dari SEMUA layar (keputusan work owner "buang saja dari semua halaman"); komponen
  PanelTotalPeserta dan tesnya dihapus. Backend `HitungTotalPeserta` dan gerbang mata uang tetap (JSON `total` tetap terkirim).

## 11-10-2026 (lanjutan 4) — Spreading Detail/Summary, konsistensi isian, validasi Claim RNM, Register Claim

- **Spreading**: subjudul "Spreading Detail" (lapisan/reinsurer) dan "Spreading Summary" (rekap per reinsurer) + baris TOTAL
  per mata uang (JSON `total`, dijumlahkan backend dengan Money.Add).
- **Isian konsisten**: semua isian layar Outstanding (Claim RNM, bank, lima tanggal, Admin Notes, saringan nama, Catastrophe
  Note) memakai kelas inti `field__input`; tiga kolom isian adjustment berkolom nilai selebar tetap (20ch) sehingga angka (rata
  kanan), teks (rata kiri), dan isian berbaris rapi; isian bank sama lebar.
- **Claim RNM**: hanya digit + satu pemisah desimal (titik/koma, `inputMode=decimal`), tanpa pemisah ribuan saat diketik;
  pesan di bawah isian "Claim RNM must be a number." / "Claim RNM cannot be greater than Share Nusantara Re (<share berformat>).";
  isian bertanda `field__input--error`; tombol Save grid nonaktif selama ada isian tidak sah. Perbandingan desimal TEKS berbasis
  digit (`frontend/desimal.ts`, inti tidak punya pembanding desimal), tanpa float. Server tetap menolak 422.
- **Register Claim**: tata letak yang SAMA dengan Policy Details (`os__form` + `dl.os__kolom` + baris `dt`/`dd`); nilai kode tidak
  lagi nowrap (meluber ke kolom sebelah) melainkan membungkus; radio Catastrophe dan Is KPR bertumpuk dengan jarak seragam;
  catatan "Saved with the Save button…" dihapus.

## 11-10-2026 (lanjutan 5) — Bawaan Non-Catastrophe / Non-KPR, Catastrophe Note di kanan, Percent/Claim Paid langsung

- **Bawaan**: Register menulis header klaim BARU dengan `STS_KATASTROFE = Non-Catastrophe` dan `IS_KPR = Non-KPR` (INSERT header
  `PohonKlaim.Simpan` kini memuat kedua kolom). Klaim LAMA berkolom kosong: layar memilih bawaan itu dan mengirimkannya pada
  Save grid berikutnya; nilai warisan asing tetap "Stored value". Tidak ada UPDATE massal; SQL opsional untuk work owner (TIDAK
  dijalankan) di laporan 11-10-2026.
- **Catastrophe Note** pindah ke kolom kanan, tepat di bawah STNC Treaty (tetap hanya bila Catastrophe; baca saja sesudah bernomor).
- **Percent Claim / Claim Paid langsung** saat Claim RNM diketik: Claim Paid = Claim RNM; Percent = `persenKlaimTeks`
  (frontend/desimal.ts, BigInt berskala) dengan aturan IDENTIK backend `models.HitungPersenKlaim`: pecahan dibulatkan
  SETENGAH-KE-ATAS pada 5 desimal (b702; apd.BaseContext half-up, diperiksa dengan Go) lalu × 100. Kasus kesamaan sama di Go dan
  TS (50/100, 1/3, 2/3, 665/5.850.038,75 = 0,011, 1/200000, 1/400000, 94.695,705/100.000). Isian tidak sah → "—".
  Catatan: 665 / 5.850.038,75 menghasilkan 0,011% (bukan 0,01137%) karena pecahannya dibulatkan 5 desimal seperti Pega.
- Kolom kiri isian adjustment memakai kolom nilai tetap 20ch yang sama dengan kolom tengah/kanan (angka rata kanan).

## 11-10-2026 (lanjutan 6) — ORA-12899 STS_KATASTROFE: migrasi 029 memperlebar kolom

- Bug DEV: pengisian bawaan "Non-Catastrophe" (15 byte, InputOSClaimLife b3728) ke `T_GENERAL_CLAIM.STS_KATASTROFE VARCHAR2(8)`
  (migrasi 002, tanpa bukti lebar) menghasilkan ORA-12899; Register dan Save grid 11-10-2026 akan gagal sama.
- **Migrasi 029** (disetujui work owner 11-10-2026): `STS_KATASTROFE` dan `IS_KPR` -> VARCHAR2(20); _down kembali ke 8 (gagal terang
  ORA-01441 bila ada nilai > 8, tanpa memotong). **Wajib dijalankan sebelum kode Register / Save grid 11-10-2026.** Jatah migrasi
  Claim Life 001-029 HABIS.
- Pengaman Go: Catastrophe dan Is KPR > 20 byte ditolak 400 sebelum bind. Audit lebar kolom teks yang ditulis aplikasi: lihat
  laporan 11-10-2026 (seluruh batas Go <= lebar DDL).

## 11-10-2026 (lanjutan 7) — Save to RNM: semua kolom wajib, semua peserta, satu penolakan

Keputusan work owner: "kolom claim rnm dibuat wajib. saat ingin save to rnm pastikan semua kolom wajib harus diisi untuk semua
peserta. dan jika ada yang belum ada isi tidak bisa lanjut dan berikan pesan".

- ⛔ **MENYIMPANG dari Pega**: SaveOutStandingLife_Act berhenti di pelanggaran PERTAMA (langkah 11.x per `local.idx`). Kini
  `services.PeriksaSimpanRNM` mengumpulkan SELURUH kekurangan kolom wajib semua peserta dan menolak SEKALI (422,
  `{galat, langkah:"wajib", kekurangan:[{nomorPeserta, pesertaId, sertifikat, nama, kolom, pesan}]}`). Kalimat per kolom tetap
  kalimat XML bila ada (DOL 11.10, DOB 11.12, Begin 11.13, Expired 11.14, Policy 11.15, Certificate 11.16, dokumen langkah 3-4).
- Urutan: (1) Type tidak dikenal — berhenti; (2) klaim ganda 11.2-11.6 — berhenti di peserta pertama (bukan "kolom kosong");
  (3) kumpulkan: Is KPR (header, baru), lalu per peserta Document (peserta di-Accept, non-TP/TR), Date of Loss (kosong ATAU di
  luar jendela, satu kalimat XML 11.10), DOB, Begin, Expired, Policy, Certificate, Claim Received Date, Confirmation Date,
  Claim RNM (satu per peserta). Nomor klaim diambil SESUDAH seluruh gerbang lolos (uji sumber
  `TestSimpanRNMGerbangSebelumPenomoran`).
- **Claim RNM wajib > 0** (`services.KlaimRNMMinimal` = frontend `KLAIM_RNM_MINIMAL` = "0"), menggantikan 11.17.1 "Claim Gross No N
  can't null" (Pega hanya menolak kosong). Kalimat: "Claim RNM is required for insured number N"; di layar "Claim RNM is required."
  di bawah isian, label bertanda *. Save grid BOLEH menyimpan Claim RNM kosong (CLAIM_AMOUNT NULL); wajibnya hanya di Save to RNM.
- Layar: pre-check kolom kosong yang sama SEBELUM memanggil server (tanpa permintaan bila ada kekurangan); satu kotak
  "Please complete the required fields before Save to RNM:" dengan satu baris per peserta ("No 1 · NAMA: Date of Loss, Claim
  RNM"; tingkat klaim = "Register Claim: Is KPR"); tanggal yang kurang merah, baris adjustment dengan Claim RNM kosong dibuka dan
  ditandai; peserta bermasalah pertama dibuka di tab Claim Details. Jendela DOL, DOB sumber, klaim ganda tetap hanya di server.
- Perubahan belum disimpan: Save to RNM NONAKTIF dengan "Save your changes first." (pilihan (a)); (b) "simpan lalu Save to RNM"
  tidak dipakai karena penomoran bertransaksi tersendiri. Dibanding nilai MENTAH tersimpan: klaim lama dengan Catastrophe / Is KPR
  kosong harus di-Save sekali dulu (bawaan Non-Catastrophe / Non-KPR belum tersimpan).
- Bug alert ganda diperbaiki: `Gagal` sudah menampilkan pesan server; `<p>{pesanGalat(...)}</p>` salinannya dibuang (dua tempat).

## 11-10-2026 (lanjutan 8) — DOL dalam periode valuation, alur Save to RNM, galat tunggal, rincian simetris

- **DOL harus di dalam periode valuation** (keputusan work owner): validasi LANGSUNG di layar (`galatDOL`, RincianPeserta.tsx) dengan
  jendela yang SAMA dengan `services.ValidasiDOL`: QR/QP jendela GROSS (BEGIN, EXPIRED]; TR/TP jendela RETRO dengan DOL digeser +1
  hari (ValidasiDOL_Act b698). Type lain / valuasi kosong: tidak diperiksa di layar (server). Pesan di bawah isian: "Date of Loss
  must be within the valuation period (dd/mm/yyyy – dd/mm/yyyy)."; isian merah; Save grid NONAKTIF selama tidak sah.
  JSON peserta kini juga membawa `valuasiRetroMulai` / `valuasiRetroSelesai`.
- Server: "Invalid DOL" (`local.errmsg` ValidasiDOL_Act) TETAP sebagai kode (`GalatDOL.Error()`, medan `kode` rute Edit Date);
  yang ditampilkan `GalatDOL.Kalimat()` = "Insured number N: Date of Loss must be within the valuation period (…)." — Save grid,
  rute Edit Date, dan Save to RNM (DOL TERISI di luar jendela). DOL KOSONG di Save to RNM tetap kalimat XML 11.10. Save grid
  memvalidasi DOL urut peserta (bukan urutan peta).
- ⚠️ Batas hari: Save to RNM (langkah 11.8, b4464) TIDAK menggeser DOL retro, Save grid (ValidasiDOL) menggeser +1 hari — untuk
  TR/TP DOL = RETRO BEGIN lolos Save grid tetapi ditolak Save to RNM, DOL = RETRO EXPIRED sebaliknya. Perbedaan dari Pega,
  dipertahankan; rentang yang ditampilkan = tanggal valuasi apa adanya (awal gross bersifat eksklusif).
- **Alur Save to RNM**: selama ada perubahan belum disimpan Save to RNM nonaktif, pre-check TIDAK berjalan, kotak kekurangan
  tidak muncul (hanya teks bantuan "Save your changes first." di atas deretan tombol). Isian berubah / Save grid / Delete
  membersihkan kotak lama; Save grid gagal = hanya galat Save grid.
- **Galat tunggal "seperti premium list"**: komponen `components/PitaGalat.tsx` — satu tempat di ATAS isi (di bawah kepala,
  sebelum tab), kelas inti `alert alert--error` / `<Gagal>`, dapat ditutup, dibersihkan setiap aksi. Acuan (baca saja):
  premiumlistlife `PremiumListDetail.tsx` / `InputOffer.tsx` (Gagal di atas panel, sebelum tabel/tombol) dan
  `FormDataPolis.tsx` (daftar "Required: …"). Dipakai di Outstanding (galat muat, Save grid, Save to RNM, Delete, kekurangan),
  Register, Inbox, Detail (KlaimLife). Kabar berhasil (Saved / Saved to RNM) ikut di tempat yang sama.
- **Rincian peserta**: tiga kelompok berjudul sama lebar — Insured (Certificate, Policy, Plan, Age, Sex, DOB), Coverage (Begin,
  Expired, Effective, Valuation Begin/Expired, Received Date, STNC Treaty), Claim Dates (DOL*, Claim Received*, Confirmation*,
  Document Complete, Max Data Received, lalu baca saja Max Expired Claim, Max Claim Received, STNC Data); di bawahnya Location of
  Death | Admin Notes sama tinggi. Label 24ch di semua kolom; Certificate/Policy No tidak dipecah (satu baris, menggulir); 3 → 2
  (≤1280px) → 1 (≤760px) kolom. Menyimpang dari urutan Pega dengan sengaja ("hanya jadikan referensi").

## 11-10-2026 (lanjutan 9) — Jendela DOL inklusif dan seragam

Keputusan work owner: "tanggal awal valuation masih termasuk".

- ⛔ **MENYIMPANG SADAR dari Pega**: jendela DOL kini `VALUATION BEGIN <= DOL <= VALUATION EXPIRED`, inklusif di kedua ujung, per
  tanggal kalender, untuk QR/QP (gross) dan TR/TP (retro). SATU sumber backend `services.DOLDalamJendela`, dipakai Save grid / Edit
  Date (`ValidasiDOL`) dan Save to RNM 11.7-11.8 (`dolDalamJendelaRNM`); frontend `galatDOL` cermin aturan yang sama.
- Dibuang: pergeseran DOL retro +1 hari ValidasiDOL_Act b698 (beserta `pergeseranDOLRetro` / `PergeseranDOLRetro`) dan
  perbandingan ketat `@CompareDates` (`bandingKetat`). Akibatnya ketidakkonsistenan lama hilang: Save grid dan Save to RNM memberi
  hasil SAMA untuk tanggal yang sama (diuji di `TestSimpanRNMDOLJendelaTanpaGeserRetro`).
- Pesan tetap menampilkan tanggal valuasi apa adanya — kini memang rentang yang berlaku. Catatan "batas hari" di lanjutan 8
  tidak berlaku lagi.

## 11-10-2026 (lanjutan 10) — Rincian tanpa gulir, Admin Notes di bawah Location of Death, aturan edit

- Certificate No / Policy No: gulir horizontal dibuang ("ada slide bar buat jelek"); nilai MEMBUNGKUS di kolom nilai
  (`overflow-wrap: anywhere`, hanya bila perlu). Kolom tetap sama lebar (memperlebar kolom Insured merusak kesimetrisan).
  Penjaga gaya: tidak ada `overflow: auto|scroll` pada aturan `rincian`.
- Susunan lama dikembalikan: LOCATION OF DEATH lalu ADMIN NOTES DI BAWAHNYA, bertumpuk di kolom kiri di bawah data Insured.
- Aturan edit (keputusan work owner "hanya bisa diedit oleh admin dan jika masih di inbox admin"): Admin Notes dapat diisi
  hanya bila pelaku ReasLifeAdmin DAN kasus terbuka di tahap Outstanding (`bolehUbahCatatanAdmin`; peran dikirim rute.tsx);
  selain itu tampil sebagai teks baca saja (tanpa kotak isian). Peserta diputus (b2628 `.STS_REJECT=='1'||'2'`) tetap
  terkunci. Server Save grid sudah menegakkan 409 (bukan tahap Outstanding) dan 403 (bukan pemegang tahap = ReasLifeAdmin)
  sebelum menulis apa pun. Sesudah Save to RNM: XML ClaimLifeDetailGCNM b2588/b2628 TIDAK mengunci Admin Notes (hanya
  `pyDisabledWhen .STS_REJECT`; kondisi `CLAIM_NO!=''` milik isian lain) → tetap dapat diisi selama masih di tahap Outstanding.
- Location of Death: dicari ulang di korpus kelvin/InboxClaimLife (Done) ("Place of Death", LocationOfDeath, PlaceOfDeath,
  Country, Province, City, HOSPITAL) — TIDAK ditemukan di file mana pun ("DEATH" hanya muncul sebagai ContentNote dan label "Double claim Death"). Tetap baca
  saja "—"; usulan kolom menunggu keputusan work owner + rentang migrasi baru (jatah 001-029 habis).

## 11-10-2026 (lanjutan 11) — Save = simpan progres, bank baca saja, spreading terlipat, tombol di tab Claim Details

- ⛔ **Keputusan WO terbaru, MEMBALIK sebagian lanjutan 8**: tombol Save grid "hanya save progress saja, tidak ada menjalankan
  fungsi apapun". DOL di luar periode valuation, Claim RNM > Share Nusantara Re, dan kolom wajib kosong TETAP tampil merah di
  isian sebagai peringatan, tetapi Save tetap aktif dan menyimpan. Yang masih memblokir Save: format rusak (Claim RNM bukan
  angka), teks melebihi lebar kolom, dan aturan kunci (sesudah Save to RNM, bukan Admin/pemegang tahap, bukan tahap
  Outstanding, baris berstatus). Sesudah Save: pita "Progress saved." di tempat galat/kabar yang sama.
- Server diselaraskan: `ValidasiDOL` dicabut dari Save grid (`SimpanOutstanding.Simpan`); `PeriksaClaimGrossTerhadapShare`
  dicabut dari `PeriksaIsianBaris`. Save to RNM (`PeriksaSimpanRNM`) menolak keduanya dalam daftar kekurangan, dengan
  kalimat layar: "Insured number N: Date of Loss must be within the valuation period (…)." dan "Insured number N: Claim RNM
  cannot be greater than Share Nusantara Re (5.077.429,5)." (share kosong tidak dibandingkan, sama dengan layar). Rute Edit
  Date (layar Detail) tetap memakai ValidasiDOL. Save grid tidak menomori, tidak spreading, tidak outbox/Arasapas (diuji).
- Bank: "semua read only kecuali name of bank". Payable To, Branch of Bank, Account No baca saja; hanya Name of Bank isian.
  Server: `IsianBarisAdjustment` hanya {jumlah, namaBank} — kunci cabangBank/nomorRekening dari klien diabaikan, dan
  UPDATE T_CLAIMLF_ADJUSTMENT tidak lagi menyentuh BRANCH_OF_BANK / ACCOUNT_NO (nilai tersimpan tidak tertimpa).
- Spreading isian adjustment TERLIPAT bawaan ("View Spreading ▸" / "Hide Spreading ▾"); dimuat sekali per peserta saat pertama
  dibuka. Isi tetap Spreading Detail + Spreading Summary + TOTAL.
- Save to RNM + Send to Medical Check (dan petunjuk perubahan belum disimpan) hanya di tab Claim Details, di kaki tab sesudah
  grid; petunjuk menyebut bagian yang berubah, mis. "You have unsaved changes (Register Claim). Click Save first.".
- Register Claim: radio Catastrophe / KPR MENDATAR (jarak seragam, bulat, label sejajar baris pertama); Catastrophe Note
  berlabel dengan placeholder "Describe the catastrophe…", lebar penuh kolom.
- Location of Death: ditunda (pilihan a) — tetap "—" baca saja.

## 11-10-2026 (lanjutan 12) — Tombol committee, konfirmasi Save to RNM, isi OS_AKSEPTASI_KLAIM_LIFE

- "Send Claim to Committee" DISEMBUNYIKAN di layar Outstanding (keputusan work owner "hilangkan saja"); rute/aksi backend
  tetap. Masih tampil di layar Detail (KlaimLife → RincianPeserta) untuk tahap selain Outstanding bila baris berstatus
  Outstanding dan belum diserahkan (bolehSerahkanDiLayar).
- Save to RNM: klik → pre-check → (lolos) Modal konfirmasi "Save to RNM" ("Are you sure…?" + alert--warn "This action cannot
  be undone…", ringkas jumlah peserta + PL, Cancel / "Yes, Save to RNM"); hanya tombol modal yang memanggil server, satu kali,
  nonaktif "Saving…" dan tak dapat ditutup selama proses. Pre-check gagal = kotak kekurangan tanpa modal; galat server di
  tempat galat biasa sesudah modal tertutup.
- **OS_AKSEPTASI_KLAIM_LIFE saat Save to RNM** (keputusan work owner; tabel warisan, tanpa DDL). INSERT dan UPDATE memakai SATU
  daftar kolom (`repository.KolomCerminSaveRNM`); nilai disusun `services.SusunCerminSaveRNM` (murni). UPDATE MENIMPA (bukan
  NVL) — Save to RNM membekukan nilai; kecuali kolom keputusan STS_REJECT dan ACCEPTATION_DATE (hanya saat INSERT).
  Teks melebihi lebar katalog = galat jelas (`PeriksaLebarWarisan`), tidak dipotong. Tanggal lewat TO_DATE berformat tetap;
  tanggal tak terbaca = galat. Di transaksi UTAMA, sesudah transaksi penomoran.

| Kolom | Sumber |
|---|---|
| ID / CASEID | ID baris adjustment / ID klaim (kunci) |
| NO_CLAIM | nomor klaim (baru lahir) |
| NO_ACCEPTATION | tanggal outstanding = saat Save to RNM (WIB), teks `DD/MM/YYYY` (kolom VARCHAR2(100)) |
| STS_REJECT / ACCEPTATION_DATE | INSERT saja: kosong (status '0' oleh SetelCerminOutstanding) / tanggal akseptasi baris (kosong) |
| POLICY_NO | header klaim; kosong → polis NO_POLIS |
| POLICY_HOLDER | polis POLICY_HOLDER_NAME (T_PREMIUM_LIST via IDPEGA) |
| CERTIFICATE_NO, PL_NUMBER | peserta |
| NAME_OF_INSURED, DOB | SQL dari baris sumber M_LIFE_PREMIUM_DETAIL (IsiTertanggungCermin, OQ-N2) |
| SEX, AGE | peserta (baris sumber) |
| PLAN | polis BUSINESS_NAME (aturan PremiumList: PLAN = BUSINESS_NAME) |
| BEGIN_DATE / EXPIRED_DATE | peserta |
| LAPSE_DATE | DATE OF LOSS peserta |
| EM_PERCENT | peserta |
| CURRENCY, CLAIM_AMOUNT | baris adjustment (Claim RNM) |
| SUM_INSURED, CEDING_RETENTION, SUM_REASURED, SHARE_NUSANTARA_RE, SHARE_RETRO, RETROCEDED_SHARE | baris adjustment; kosong → peserta |
| CLAIM_RETRO | header klaim |
| WPC, CEDINGCONAME, SOBNAME, BUSINESSNAME, PRODUCTNAMEID, PRODUCTNAME | polis |
| BUSINESSID | polis BUSINESS_CODE (CARI31 PolicyDataLife.BusinessCode) |
| SOB | polis T_PREMIUM_LIST.SOB (dibaca terpisah, `KlaimLife.SOBPolis`; PolisRingkas inti tak membawanya) |
| CEDINGCO | SQL dari polis CEDING_CO (tetap seperti sebelumnya) |
| RETROID, RETRONAME, SECURITYREINSURERID, SECURITYREINSURER | polis; ditukar bila Type TP/TR dan Security Reinsurer terisi (b9671) |
| TYPECEDING | polis TYPE_CEDING 1/2/3/4 → QS / SURPLUS / QS + SURPLUS / XOL (CARI40) |
| TYPE | Type header klaim |
| CONFIRMATION_DATE, CLAIM_RECEIVED_DATE, COMPLETE_DATE | peserta |
| CREATEOPNAME | pembuat klaim |
| NAME_OF_BANK, IDBANK, ACCOUNTNO | baris adjustment |
| STATUS, DISEASE, ICD_CODE, NOTES, KETERANGAN | TIDAK terpetakan (tidak ditulis) |
| STS_KONVERSI, TGL_KONVERSI, LAYER_1..4, INDEXLIST | di luar 55 kolom rule warisan; tidak ditulis |

- Menyimpang dari Pega: InsertJsonKlaimLife_sql tidak menulis NO_ACCEPTATION dan menulis ACCEPTATION_DATE = SYSDATE; di sini
  NO_ACCEPTATION = tanggal outstanding (keputusan WO) dan ACCEPTATION_DATE dibiarkan milik akseptasi/Komite.

## 11-10-2026 (lanjutan 13) — Kolom Medical, ACCEPTATION_DATE, Komite hanya Claim Analis, MAX CLAIM RECEIVED

- STATUS, DISEASE, ICD_CODE, NOTES, KETERANGAN (OS_AKSEPTASI_KLAIM_LIFE): keputusan work owner "sumbernya nanti saat di
  medical" — diisi pada tahap Medical Check (PEKERJAAN BERIKUTNYA, belum dikerjakan). Save to RNM: INSERT membiarkannya NULL,
  UPDATE tidak menyentuhnya (tidak ada di KolomCerminSaveRNM), sehingga nilai yang diisi tahap lain tidak terhapus (diuji).
- ACCEPTATION_DATE: keputusan work owner "iya gunakan tanggal outstanding masuk ke table" — Save to RNM menulisnya = tanggal
  outstanding (waktu Save to RNM, WIB, DATE), INSERT dan UPDATE baris yang belum berkeputusan (STS_REJECT NULL/0). STS_REJECT
  tetap '0' lewat SetelCerminOutstanding. Perilaku yang ada sesudahnya: Komite `RekamAkhirWarisan` (komiteclaimlife
  komite_keputusan.go) MENIMPA STS_REJECT, NO_ACCEPTATION, ACCEPTATION_DATE baris itu saat keputusan akhir (1/2); Akseptasi /
  status baris Claim Life (`sqlIkutkanStatusCermin`) hanya menulis STS_REJECT cermin (ACCEPTATION_DATE miliknya di
  T_CLAIMLF_ADJUSTMENT).
- "Send Claim to Committee" HANYA di tahap Claim Analis (keputusan work owner "iya itu hanya di claim analis saja"): layar
  merendernya hanya bila tahap = Claim Analis; server `Penyerahan.Serahkan` menolak 409 (`ErrSerahkanBukanClaimAnalis`) di
  tahap lain — gerbang BARU (sebelumnya hanya wewenang per Type, kasus terbuka, status baris, rekening).
- MAX CLAIM RECEIVED: teks galat internal tidak lagi dirender — layar menampilkan "—" dengan tooltip "Not available: product
  threshold not set". Akar: Claim Life membaca ambang dari view warisan PRODUCTINWARD_LIFE (JSON), padahal Master Product Name
  Life sejak 02-10-2026 menulis tabel DATAR M_PRODUCTNAME_LIFE (view DEV tidak dibangun ulang). Kini `ProdukLife.Ambang` membaca
  M_PRODUCTNAME_LIFE.MAXEXPIREDCLAIM ("Max Notification Claim Expired") / MAXDATARECEIVE berkunci ID = ProductNameID polis
  lebih dulu; view hanya cadangan bila produknya belum ada di tabel datar. Rumus tetap XML ValidasiClaimReceived_Act
  (b562/b583): MAX CLAIM RECEIVED = penanda — kosong bila selisih hari DOL→CLAIM RECEIVED DATE ≤ ambang, tanggal CLAIM RECEIVED
  (dd/MM/yyyy) bila melewatinya; MAX EXPIRED CLAIM = nilai ambang produk. Ambang kosong di Master = "—" (bukan nol). Pesan usang
  "menunggu modul PremiumList Life" diganti `models.ErrAmbangHariKosong`.

## 11-10-2026 (lanjutan 14) — Inbox PL Number & Update Operator, Policy Details, Reject, dokumen

- Inbox: "Policy No" diganti "PL Number" = PL peserta aktif klaim (subkueri skalar berindeks IX_PLD_CLAIM_ID, STS_HAPUS IS
  NULL; tanpa peserta = "—"). Urutan kolom (tabel dan ekspor xlsx): Case ID · Claim No · PL Number · Create Operator Name ·
  Update Operator Name · Create Date/Time · Work Status.
- Update Operator Name: T_WORK_CLAIM / T_GENERAL_CLAIM TIDAK punya kolom operator-pengubah (hanya TGL_UPDATE). Tanpa migrasi
  (jatah habis) ia DITURUNKAN dari T_CLAIMLF_JEJAK: AKUN_ID jejak terakhir klaim (KEEP DENSE_RANK LAST, indeks
  IX_CLAIMLF_JEJAK_KLAIM), cadangan CREATE_OP_NAME (padanan pxUpdateOpName = pembuat). ⚠️ Jejak hanya mencatat perubahan
  status baris (Save to RNM, Reject, akseptasi, Komite) — Save grid dan pindah tahap tidak tercatat. Kolom sungguhan
  (mis. T_WORK_CLAIM.UPDATE_OP + diisi setiap aksi pengubah) menunggu keputusan work owner + rentang migrasi baru.
- Policy Details: label "PL Number" (bukan PL_NUMBER); nilai kode satu baris (os__nilai-pl nowrap) dengan kolom panel
  minmax(380px) dan label 18ch, sehingga "RNML-QRL2.10.26.00043" muat di 1280 (3 kolom) dan 1024 (2 kolom), tanpa gulir.
- Reject Outstanding: modal konfirmasi — "Are you sure you want to reject this outstanding claim for <nama / sertifikat>?" +
  alert--warn "This action cannot be undone. …"; isian Remarks Pega (RejectOSClaimLife_Sec, wajib) dipertahankan; Cancel +
  "Yes, Reject" (danger), nonaktif "Rejecting…" dan tak dapat ditutup selama proses; penjaga klik ganda.
- Dokumen — akar "tidak bisa upload":
  1. Isian "Kategori" teks BEBAS, dan kotak berkas NONAKTIF sampai teks diisi (Choose File tampak mati).
  2. Server menolak kategori di luar CATEGORY_ATTACH_CLAIMLIFE (422 "kategori tidak dikenal") — teks bebas hampir selalu ditolak.
  3. Server menuntut env UNGGAHAN_DIR (folder berkas); kosong = 503 "UNGGAHAN_DIR belum disetel".
  4. Claim Life TIDAK punya pekerja latar outbox (Modul.JalankanPekerja = TanpaPekerja): efek "storage-unggah" tidak pernah
     dijalankan, T_STORAGE_ID tetap kosong, baris tampil "Processing — waiting for storage" tanpa View/Delete.
  Diperbaiki di modul: Category = dropdown dari rute baru GET /api/klaim-life/dokumen/kategori (sumber yang sama dengan
  pemeriksa unggahan), pilih berkas lalu tombol "Add attachment", galat tampil sebagai alert, label Inggris.
  Butuh keputusan/konfigurasi (TIDAK dikerjakan): (a) set UNGGAHAN_DIR di server DEV; (b) menyalakan pekerja outbox Claim
  Life — PekerjaEfek memungut SEMUA efek modul CLAIMLIFE, sedangkan PelaksanaBerkasLokal menolak jenis lain sebagai gagal
  permanen; perlu saringan jenis di outbox (inti) atau pelaksana gabungan, plus interval di config (inti) — keputusan tim inti
  / work owner; (c) pelaksana penyimpanan nyata (Google Storage) menuntut persetujuan + kredensial.

## 11-10-2026 (lanjutan 15) — Dokumen dapat dihapus, View lokal, tanggal

- Delete ("buat document bisa dihapus"): setiap baris dokumen — termasuk yang masih "Processing — waiting for storage" —
  punya View + Delete; Delete lewat Modal "Delete this document?" + nama berkas + alert--warn "This cannot be undone.",
  Cancel / "Yes, Delete" (nonaktif "Deleting…", tak dapat ditutup selama proses).
- Server `Unggahan.Hapus`: satu transaksi — hapus baris T_CLAIMLF_DOCUMENT; bila T_STORAGE_ID kosong (efek unggah belum
  jalan) efek "storage-unggah" dokumen itu yang masih ANTRE ditutup "selesai" dengan catatan `CatatanBatalUnggah` (outbox
  inti tidak punya status "dibatalkan"; tidak masuk daftar perlu-intervensi), lalu berkas lokal di UNGGAHAN_DIR dibuang
  SESUDAH commit (best-effort, gagal dicatat di log). Bila T_STORAGE_ID terisi: jalur lama DeleteDocument_Act /
  DeleteGoogleStorage_Act — efek "storage-hapus" diantre (pelaksana menghapus kartu berkas + berkas lokal). Pengaman ganda:
  pelaksana efek unggah melewati dokumen yang sudah tidak ada (efek yang sedang "jalan" saat dihapus).
- Gerbang (XML): DeleteDocument_Act b513 Obj-Delete TANPA prasyarat (hanya panggilan penyimpanan b472 bersyarat
  T_STORAGE_ID); DocumentLife / ClaimLifeDetailGCNM tidak punya syarat peran, tahap, atau kunci (tidak terkunci sesudah Save
  to RNM maupun peserta diputus). Diterapkan minimal: identitas (401), kasus terbuka (409, butir bb), dokumen milik klaim
  (404, SatuDokumen berbatas klaim). Dampak: dokumen terakhir peserta yang di-Accept dihapus → Save to RNM menolak lagi
  (gerbang dokumen) — memang benar.
- View: rute baru GET /api/klaim-life/{id}/dokumen/{dokId}/isi (berbatas klaim; dokumen klaim lain 404). Selama penyimpanan
  belum berjalan, berkas dilayani dari salinan LOKAL di UNGGAHAN_DIR (jalur disusun server dari pengenal dokumen, tidak dari
  klien; Content-Disposition attachment, nosniff; berkas hilang = 404). Rute lama /api/dokumen/{dokId}/isi tetap.
- Kolom Date tampil dd/mm/yyyy HH:mm (tanggalDokumen), bukan ISO mentah.

## 11-10-2026 (lanjutan 16) — Claim Documents: kategori dari produk, penyimpanan segera

- Judul "Claim Documents"; kosong "No claim documents for this insured yet.".
- Kategori = Document Claim Master Product Name Life (M_PRODUCTNAME_LIFE_DOCCLAIM: PRODUCTID, URUT, DOCUMENT VARCHAR2(500);
  baca saja, SQL disalin, tanpa impor) untuk ProductNameID polis klaim (header polis via IDPEGA), urutan URUT. Rute
  GET /api/klaim-life/{id}/dokumen/kategori; validasi unggah memakai sumber yang SAMA (422 di luar daftar). Nilai disimpan
  = teks DOCUMENT (master tidak punya ID dokumen selain URUT posisional, yang berubah bila daftar diubah). Produk tanpa daftar:
  dropdown kosong + "No claim document types are set for this product in Master Product Name Life.". CATEGORY_ATTACH_CLAIMLIFE
  (butir ar1) tidak lagi dipakai jalur mana pun (gerbang dokumen Save to RNM hanya menghitung ada-nya dokumen). Dokumen lama
  tetap tampil apa adanya.
- ⛔ BERHENTI (butuh migrasi): T_CLAIMLF_DOCUMENT.KATEGORI_2 hanya VARCHAR2(64), DOCUMENT master VARCHAR2(500); nama seperti
  "Copy identitas diri KTP/SIM/Paspor (…)" > 64 byte. Sementara nama yang terlalu panjang ditolak terang 422
  (ErrKategoriDokumenTerlaluPanjang) dan opsinya dinonaktifkan di layar dengan keterangan — tidak dipotong. Pelebaran ke
  VARCHAR2(500) menunggu persetujuan work owner + rentang migrasi baru.
- Penyimpanan: master (masterproductnamelife mpnl_lampiran.go) merekam + mengantre outbox dalam satu transaksi lalu MENCOBA
  MENGIRIM SEGERA di permintaan yang sama (kirimLampiran) — tanpa pekerja latar; target penyimpanan nyata Google Storage
  (ServiceGoogle via M_LINK_SERVICE + token inti/layanan, garam STORAGE_TOKEN_SALT), objek dicatat di tabel BERSAMA
  T_STORAGE_IMAGE; berkas antre lokal di UNGGAHAN_DIR/master-product-name-life/antre/<IMAGEID>.
  Diterapkan di Claim Life: MEKANISME kirim segera — Unggah memanggil kirimSegera sesudah commit (pungut efek
  "storage-unggah" dokumen itu, jalankan pelaksana Claim Life, tuntaskan / jadwal ulang / gagal permanen), sehingga dokumen
  baru langsung berstatus tersimpan (View/Delete) tanpa menunggu pekerja. Pelaksananya tetap PelaksanaBerkasLokal (kartu
  T_CLAIMLF_STORAGE, berkas di UNGGAHAN_DIR).
  ⛔ BELUM (butuh keputusan): mengirim ke Google Storage nyata dan mencatat ke T_STORAGE_IMAGE seperti master — panggilan
  layanan luar + tabel bersama lintas aplikasi; perlu persetujuan work owner (dan kredensial/garam di server). Efek
  "storage-unggah" yang sudah antre di DEV sebelum perubahan ini TETAP antre (tidak dipungut siapa pun): dokumennya tetap
  dapat di-View (salinan lokal) dan di-Delete (efeknya dibatalkan).

## 11-10-2026 (lanjutan 17) — Rentang migrasi baru (berhenti) dan Google Storage dokumen klaim

- A. RENTANG MIGRASI BARU — BERHENTI (butuh perubahan INTI). Penjaga inti (`inti/backend/penjaga/modulmd_test.go`
  `polaRentang ^(\d{3})-(\d{3})$`, `jatahDari`) hanya menerima SATU rentang tiga digit per modul. Seluruh 001–899 sudah
  terbagi (claimlife 001–029, komiteclaimlife 030–049 dengan hanya 030 terpakai, … bordereaux 890–899; masterprovince 880–883
  bertabrakan lama dengan aggregate 880–889), 900–999 slot menu: TIDAK ADA blok bebas. Perluasan bersambung mustahil (030
  dipakai komite). Opsi: (1) penjaga mendukung multi-rentang / empat digit — perubahan inti; (2) komiteclaimlife melepas
  sebagian 031–049 lalu claimlife memakai rentang kedua — tetap butuh multi-rentang; (3) mengubah rentang modul lain — butuh
  pemiliknya. Tidak ada migrasi yang dibuat (KATEGORI_2 500, UPDATE_OP/UPDATE_OP_NAME tetap tertunda).
- B. GOOGLE STORAGE seperti Master Product Name Life (keputusan work owner): pola mpnl_storage.go disalin ke
  `services/penyimpanangoogle.go` (tanpa impor): XML Claim Life memakai aktivitas yang SAMA (InsertDocument_Act → InsertGoogleStorage_Act,
  Folder "Claim", Durasi 1800; Folder "Claim/Doc/YYYY/MM/", Namafile "yyyyMMdd-hhmmss-S - nama"; GetUrlGoogleStorage_Act;
  DeleteGoogleStorage_Act). Alamat M_LINK_SERVICE ("Google"/upload|geturl|delete) lewat inti/layanan, token GCP_IMAGE dipakai
  ulang dulu lalu dirakit dengan STORAGE_TOKEN_SALT (seperti master). Objek DITULIS ke T_STORAGE_IMAGE (IMAGEID, URLPUBLIC,
  APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE='standard'; TANGGAL_UPLOAD saat URL diperbarui), APPNAME dari T_FOLDER_IMAGE;
  T_CLAIMLF_DOCUMENT.T_STORAGE_ID = IMAGEID.
  Alur: unggah → rekam + antre efek (satu transaksi) → kirim segera (PelaksanaGoogle) → T_STORAGE_IMAGE + T_STORAGE_ID; gagal →
  efek dijadwal ulang / gagal permanen, dokumen "Processing". View: objek T_STORAGE_IMAGE → URL bertanda tangan (dipakai selama
  EXPDATE belum lewat, selain itu geturl + Update_T_Storage_SQL) diunduh backend (https, tanpa pengalihan); selain itu salinan
  lokal. Delete: objek Google dihapus LEBIH DULU (gagal = rekam tetap), lalu satu transaksi baris dokumen + baris T_STORAGE_IMAGE
  (atau kartu stub lama / pembatalan efek), lalu berkas lokal.
  T_CLAIMLF_STORAGE: dipensiunkan untuk dokumen baru; dokumen lama berkartu stub tetap di-View dari salinan lokal dan di-Delete
  (kartu dihapus). PelaksanaBerkasLokal tidak dipakai lagi oleh unggahan baru.
  Konfigurasi DEV: STORAGE_TOKEN_SALT TIDAK ada di .env (UNGGAHAN_DIR ada) — token baru tidak dapat dirakit; bila GCP_IMAGE punya
  token berlaku ia dipakai ulang. Tanpa alamat/APPNAME/token: gagal tertutup "Storage is not configured: …" (503 di View;
  unggahan tetap tersimpan lokal, "Processing").
- ⛔ RALAT lanjutan 17 (hari yang sama): transport HTTP Google DICABUT lagi dari modul — penjaga inti TestNolAlamatLayananDiKode
  (lintasaplikasi_test.go) hanya mengizinkan klien HTTP keluar di berkas yang terdaftar di daftar izin MILIK INTI (mpnl_storage.go,
  tco_pengirim_storage.go). Menambah `modul/claimlife/backend/services/penyimpanangoogle_transport.go` ke daftar itu = perubahan
  inti → menunggu tim inti / work owner. Yang TETAP terpasang: repository T_STORAGE_IMAGE / T_FOLDER_IMAGE, antarmuka
  PenyimpananDokumen, PelaksanaGoogle (rekam T_STORAGE_IMAGE + T_STORAGE_ID), View (objek → penyimpanan, selain itu lokal),
  Delete (objek → T_STORAGE_IMAGE → baris). Tanpa transport: unggahan memakai PelaksanaBerkasLokal seperti sebelumnya (dokumen
  langsung tersedia), dan objek T_STORAGE_IMAGE tidak lahir. Transport + ujinya (server HTTPS tiruan) tersimpan siap pasang.

## 11-10-2026 (lanjutan 18) — Proteksi Send to Medical Check

- Keputusan work owner: "hanya bisa send to medical jika sudah sukses save to rnm dan semua peserta yang dicentang accept minimal
  punya 1 document. untuk peserta yang tidak dicentang accept tidak apa".
- Server (`TahapLayanan.Pindah`, Outstanding → Medical Check), urutan: PeriksaPindahDariOutstanding (CLAIM_NO terisi b21832, Type
  QP/QR b22123; 409) → BARU `SudahSaveRNM` (baris sudah Outstanding; 409 ErrPindahMenungguNomorKlaim) → BARU
  `PeriksaDokumenAcceptMedis`: jumlah dokumen per peserta AKTIF ber-IS_ACCEPT 'true' lewat SATU kueri berkelompok
  (`CacahDokumenPesertaAccept`, LEFT JOIN T_CLAIMLF_DOCUMENT, STS_HAPUS IS NULL); dokumen berstatus apa pun dihitung; peserta
  tanpa Accept tidak diperiksa; SEMUA kekurangan → 422 `{galat, langkah:"dokumen-medis", kekurangan:[…]}` (pola kekurangan
  Save to RNM: permintaan sah, data belum lengkap), pesan "Insured number N (NAMA): at least one claim document is required
  because the insured is accepted.". Tanpa peserta Accept = lolos (WO tidak melarang).
- PARITAS: XML SendtoMedical (FlowAction/Section/Activity) hanya menggerbangi CLAIM_NO dan Type; syarat dokumen peserta Accept di
  XML ada di Save to RNM (SaveOutStandingLife_Act langkah 3-4), bukan di Send to Medical — aturan ini TAMBAHAN keputusan WO.
- Layar: Send to Medical Check tetap nonaktif sebelum Save to RNM; perubahan belum disimpan (termasuk sesudah Save to RNM)
  menonaktifkannya dengan petunjuk; klik → pra-periksa (Accept tersimpan + dokumen yang dimuat) → kotak PitaGalat berjudul
  "Please upload at least one claim document for every accepted insured before Send to Medical Check:" + daftar per peserta,
  peserta pertama dibuka di sub-tab Document; lolos → pindah (layar Outstanding tidak punya dialog konfirmasi Send to Medical;
  tombolnya memindah langsung seperti sebelumnya). 422 server berdaftar tampil di kotak yang sama.

## 11-10-2026 (lanjutan 19) — Konfirmasi Send to Medical Check

- Keputusan work owner "iya perlu pop up konfirmasi": klik Send to Medical Check → pra-periksa dokumen (lanjutan 18); kurang =
  kotak kekurangan tanpa modal; lolos = Modal inti "Send to Medical Check" ("Are you sure you want to send this claim to Medical
  Check?" + alert--warn "After sending, the claim moves to the Medical Advisor queue and can no longer be edited from the
  Outstanding screen.", ringkas Claim No / jumlah peserta / jumlah peserta Accept), Cancel / "Yes, Send to Medical Check"
  (nonaktif "Sending…", tak dapat ditutup selama proses, penjaga klik ganda).
- Satu pemanggil pindahTahap di layar Outstanding (tombol modal; PanelPindahTahap tidak lagi memindah untuk Medical Check).
  Galat server (422 berdaftar → kotak kekurangan; lainnya → PitaGalat) sesudah modal ditutup.
- Sesudah berhasil: layar kembali ke Inbox Claim Life (perilaku lama, onPindah) dan kini membawa pita sukses
  "Claim CLMLF-… sent to Medical Check." (state rute `kabar`, pola pemberitahuan Register → Inbox; ditutup dengan tombol atau
  saat membuka kasus).

## 11-10-2026 (lanjutan 20) — Layar Medical Check, kotak MEDICAL ADVISOR, Recommendation

- Navigasi (work owner: "dari Inbox, setiap tahap membuka layar yang BENAR"): Inbox membuka `OutstandingClaimLife` untuk SEMUA
  tahap; pita kuning peringatan tahap dan tombol "Open claim detail" DIHAPUS (prop `onDetail` dicabut dari layar dan rute).
  Tata letak sama (kepala + Claim No, tab Register Claim / Policy Details / Claim Details, grid peserta + rincian); isian Admin
  hanya hidup di Outstanding (`bolehIsianOutstanding`), jadi di tahap lain data klaim baca saja. Tombol bawah per tahap
  (`TombolBawahTahap`): Medical Check = Save · Send Back to Admin · Send to Claim Analyst (`FlowAction/MedicalCheckClaimLife`);
  Claim Analis = Send Back to Admin · Send Back to Medical · Close Claim (PanelTutupKlaim). Kabar Inbox per tujuan
  (`kabarPindah`, `kabarTutup`). `PanelPindahTahap.sesudahPindah` kini menerima tujuannya.
- `pages/KlaimLife.tsx` (layar Detail): TIDAK lagi dibuka dari Inbox; masih terpasang sebagai halaman menu `detail`
  (HALAMAN_CLAIMLIFE, pencarian per ID). Tidak dihapus — keputusan WO.
- Kotak MEDICAL ADVISOR (RincianPeserta, tahap Medical Check dan Claim Analis): grid diagnosa Find Disease (tombol per baris,
  pencarian kini di Modal inti) | DIAGNOSE | ICD CODE | GROUP DIAGNOSE (dropdown) | Delete, "+ Add" di kepala (dan kaki bila
  ada baris); teks Indonesia diganti English ringkas. Dapat diubah HANYA tahap Medical Check + ReasLifeMedicalAdvisor + kasus
  terbuka (`bolehUbahMedis`) + peserta belum diputus (STS_REJECT b4682); selain itu baca saja (3 kolom teks).
  Server: `DiagnosaPeserta.pagari` kini menolak tahap selain Medical Check (409 `ErrDiagnosaBukanMedicalCheck`) sebelum gerbang
  pemegang tahap (ReasLifeMedicalAdvisor, 403).
- GROUP DIAGNOSE "tidak dapat dimuat" — akarnya: dropdown b5863 ber-`pyListSource associated` pada rule properti
  `.GROUPDIAGNOSE` kelas `Data-DiagnoseLife`, TIDAK ada di ekspor korpus (OQ-L, butir bf); layar dulu menampilkan
  `BelumTersedia`. Kini dropdown dengan daftar `PILIHAN_KELOMPOK_DIAGNOSA` (KOSONG, tidak dikarang) → nonaktif + "Group Diagnose
  list is not configured yet."; nilai tersimpan tetap tampil. Perlu dari WO/DBA: tabel/daftar sumbernya, mis. cek
  `SELECT owner, table_name FROM all_tables WHERE table_name LIKE '%DIAG%' OR table_name LIKE '%GROUP%DIS%'`.
- RECOMMENDATION (b7316/b7323 `.RECOMMENDATION`, pyDisabledWhen STS_REJECT 1/2): kolom
  `T_CLAIMLF_PREMIUMLIST_DETAIL.RECOMMENDATION VARCHAR2(1000)` SUDAH ada (migrasi 003) — tanpa migrasi baru. Dibaca lewat
  `kolomPeserta` (JSON `rekomendasi`), textarea di bawah grid diagnosa, ditulis tombol Save Medical Check:
  `PUT /api/klaim-life/{id}/medis` `{rekomendasi:{pesertaId:teks}}` → `SimpanMedis.Simpan` (identitas, kasus terbuka bb,
  tahap Medical Check 409, pemegang tahap 403, peserta aktif, ≤ 1000 byte 422, UPDATE ber-syarat STS_HAPUS IS NULL dan
  STS_REJECT bukan 1/2 → 409 bila diputus), satu transaksi. Perubahan belum disimpan menonaktifkan tombol perpindahan
  ("Save your changes first.").
- Usulan (BELUM diterapkan) pengisian OS_AKSEPTASI_KLAIM_LIFE dari Medical Check — di XML, DISEASE / ICD_CODE / NOTES /
  KETERANGAN / STATUS ditulis dari properti datar peserta (`.DISEASE`, `.ICD_CODE`, `.NOTES`, `.KETERANGAN`, `.STATUS`) oleh
  SaveAdjustment_Act / SaveOutStandingLife_Act / RejectOSClaimLife_Act; `SetDisease` hanya menulis baris DiagnoseList
  (`.DISEASE`, `.ICDCODE`), jadi tidak ada aturan XML yang menyalin daftar diagnosa ke tabel datar. Usul: saat Send to Claim
  Analyst (Medical → Claim Analis, transaksi `Pindah`), UPDATE baris OS_AKSEPTASI peserta klaim: DISEASE = nama diagnosa
  digabung "; " urut URUTAN (dipotong aman ke 1000 byte, PeriksaLebarWarisan); ICD_CODE = ICD diagnosa PERTAMA yang terisi
  (kolom VARCHAR2(10) tidak muat gabungan); STATUS, NOTES (catatan Claim Analyst) tidak disentuh; KETERANGAN = Recommendation
  hanya bila WO menyetujui (tidak ada padanan XML). Menunggu keputusan WO.
