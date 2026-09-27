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
| 24 | `SearchDiagnose_act` | `Property-Set` | tidak | — | — | — | **A3 — Medis** |
| 25 | `SetDisease` | `Property-Set`, `Obj-Save` | **ya** | — | — | — | **A3 — Medis** *(bukti **al**)* |
| 26 | `DownloadDocumentClaim` | `Property-Set`, `Call GetUrlGoogleStorage_Act` | tidak | `EfekBerkas` *(gagal terang)* | — | — | **A3 — Dokumen** |
| 27 | `ProtectCloseClaim_act` | `Property-Set`, `Page-Set-Messages`, `Call FinishAssignment` | tidak | — | — | — | **A3 — Tutup & lihat** |
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
| 6 | `AkseptasiClaimLife` | `InputAkseptasiClaimLife`, `PreCaimLife_Act` | tombol Save Adjustment | **sebagian** |
| 7 | `MedicalCheck` | `MedicalCheckClaimLife`, `PreCaimLife_Act` | — | **A3 — Medis** |
| 8 | `SendtoMedical` | `SendtoMedical_Section` | — | **A3 — Medis** |
| 9 | `SendtoAdmin` | `SendtoAdmin_Section` | — | **A3 — Medis** |
| 10 | `AttachDocumentLife` | `AttachDocScreenLife`, `NewAttachLife`, `SaveAttach…` | — | **A3 — Dokumen** |
| 11 | `ConfirmDeleteAttachment` | `DeleteDocument_Act` | — | **A3 — Dokumen** |
| 12 | `UploadCSV_ClaimLife` | `UploadCSVClaimLife_Act`, `pxUploadCSVResults` | — | ⛔ **tidak ditiru — MESIN BAWAAN PLATFORM**. Activity-nya hanya tiga langkah *(b250, b340, b488)*, dan langkah ketiganya memanggil `pxUploadCSVResults`. Pemetaan kolom CSV-nya **tidak ada di korpus** *(OQ-F)* |
| 13 | `CloseClaim` | `CloseClaim_Section` b1081 → `ProtectCloseClaim_act` b1101 | `GET /api/klaim-life/{id}/boleh-tutup` | **sebagian** — GERBANGnya ada dan lengkap *(prasyarat b608 `.STS_REJECT!=1`, b992 `CARI1==""`)*; sisi `Call FinishAssignment` b836 **dinyatakan belum ada** sebab tahap tujuannya belum dibaca dari `Flow/`. Pesannya memakai **nomor sertifikat**, bukan nama — penyimpangan sadar, `kolomSalin` tiket 02 |
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
