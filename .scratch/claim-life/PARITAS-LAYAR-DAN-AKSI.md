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
| 3 | `SaveInsuredClaim_Act` | `Property-Set`, `Page-Remove` | tidak | — | — | — | **A3 — Register** |
| 4 | `SaveAdjustment_Act` | `Property-Set`, `RDB-List` | tidak¹ | `Akseptasi.SimpanAdjustment` | `POST …/peserta/{pid}/akseptasi` | tombol **Save Adjustment** di `KlaimLife.tsx` | **ada** |
| 5 | `RejectOSClaimLife_Act` | `Property-Set`, `RDB-List`, `Obj-Save` | **ya** | `Status.Tolak` | `POST …/adjustment/{aid}/tolak` | tombol Tolak | **ada** |
| 6 | `CreateKMTLife_Act` | `Property-Set`, `Call pxRetrieveReportData`, `Call pxAddChildWork`, `Obj-Refresh-And-Lock`, `Obj-Save`, `Call SendEmailKlaimLF` | **ya** | `Penyerahan.Serahkan` + `BuatKasusKomite` | `POST …/adjustment/{aid}/komite` | tombol Serahkan | **ada** *(sisi induk belum — lihat §4)* |
| 7 | `GetListKomiteLife` | `Property-Remove`, `Property-Set`, `Property-Set-Messages`, `Call pxRetrieveReportData`, `Obj-Save` | **ya** | `RosterKomiteOracle` | *(di dalam Serahkan)* | — | **ada** |
| 8 | `UpdateDateClaimLife_Act` | `Property-Set`, `RDB-List` | **ya²** | `Tanggal.Ubah` | `PUT …/peserta/{pid}/tanggal-kejadian` | — | **ada** *(kontrol React belum)* |
| 9 | `ValidasiDOL_Act` | `Page-Clear-Messages`, `Property-Set`, `Property-Set-Messages` | tidak | `PeriksaDOL` | *(di dalam `Tanggal.Ubah`)* | — | **ada** |
| 10 | `DeletePesertaClaimLife` | `Property-Set`, `Obj-Save` | **ya** | — | — | — | **A3 — Register** |
| 11 | `LoadDataPeserta_Act` | `Property-Set`, `RDB-List` | tidak | `PesertaPolis.AmbilUntukKlaim` | `GET /api/peserta-life` | `RegisterKlaim.tsx` | **ada** |
| 12 | `LoadDataPesertaSpesifik_Act` | `Page-Remove`, `Property-Set`, `RDB-List` | tidak | — | — | — | **A3 — Register** |
| 13 | `SelectAllClaimLife_act` | `Property-Set` | tidak | — | — | — | **A3 — Register** |
| 14 | `SearchPolicyHolder_act` | `Property-Set` | tidak | — | — | — | **A3 — Register** |
| 15 | `ValidasiClaimReceived_Act` | `RDB-List`, `Property-Set` | tidak | — | — | — | **A3 — Register** |
| 16 | `ValidasiSTNC_Act` | `RDB-List`, `Property-Set` | tidak | — | — | — | **A3 — Register** |
| 17 | `CountClaimAmountLife_Act` | `Page-Clear-Messages`, `Property-Set`, `Property-Set-Messages`, `call SpreadingClaimLife_Act` | tidak | `HitungSpreading` | — | — | **ada** *(rute + kontrol belum)* |
| 18 | `CheckTotalAdjustmentClaim` | — | — | — | — | — | ⛔ **tidak ada di korpus** *(Koreksi 2)* |
| 19 | `SetIndexAdjustmentList` | `Property-Set` | tidak | `IsCheck` peserta | — | — | **ada** *(dalam kontrak API)* |
| 20 | `SetClaimXOL_Act` | `Page-Remove`, `Property-Set` | tidak | — | — | — | **A3 — Outstanding** |
| 21 | `SendtoMedical_Act` | `Property-Set` | tidak | `Tahap.KirimKeMedis` | — | — | **ada** *(rute + kontrol belum)* |
| 22 | `SendtoAdmin_Act` | `Property-Set` | tidak | `Tahap.KembaliKeAdmin` | — | — | **ada** *(rute + kontrol belum)* |
| 23 | `SendtoAdmin_Act1` | `Property-Set` | tidak | — | — | — | **A3 — Medis** *(beda dengan no. 22 dibaca di kelompok itu)* |
| 24 | `SearchDiagnose_act` | `Property-Set` | tidak | — | — | — | **A3 — Medis** |
| 25 | `SetDisease` | `Property-Set`, `Obj-Save` | **ya** | — | — | — | **A3 — Medis** *(bukti **al**)* |
| 26 | `DownloadDocumentClaim` | `Property-Set`, `Call GetUrlGoogleStorage_Act` | tidak | `EfekBerkas` *(gagal terang)* | — | — | **A3 — Dokumen** |
| 27 | `ProtectCloseClaim_act` | `Property-Set`, `Page-Set-Messages`, `Call FinishAssignment` | tidak | — | — | — | **A3 — Tutup & lihat** |
| 28 | `setDetailClaim_act` | `Obj-Open-By-Handle`, `Property-Set`, **`Obj-Save`** | **ya** | — | — | — | **A3 — Tutup & lihat** ⛔ *(Koreksi 3: ia MENULIS)* |
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
| 12 | `UploadCSV_ClaimLife` | `UploadCSVClaimLife_Act`, `pxUploadCSVResults` | — | **A3 — Register** |
| 13 | `CloseClaim` | `CloseClaim_Section` | — | **A3 — Tutup & lihat** |
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
