# Spesifikasi tombol Treaty In — dibaca dari ekspor Pega, bukan dikarang

**Disapu 6 Oktober 2026 dari `D:\XML_NURE\Treaty In`.**
Cakupan: seluruh 52 berkas `Section/*.xml` + `Harness/*.xml`.
Inventaris mentahnya: [`lampiran/INVENTARIS-TOMBOL-EKSPOR.txt`](lampiran/INVENTARIS-TOMBOL-EKSPOR.txt).

---

## ⚠️ Dua aturan baca — ditambahkan 6 Oktober 2026 sore, sesudah keduanya memakan korban

1. **`1=2` di XML belum berarti mati.** Yang menentukan PASANGANNYA:

   | Tingkat | Pilihan visibilitas | Syarat | Artinya |
   | --- | --- | --- | --- |
   | sel | `pyVisible = OTHER` | `pyCondition = 1=2` | **mati** |
   | sel | `pyVisible = ALWAYS` | `pyCondition = 1=2` | **hidup** — syaratnya sisa |
   | wadah | `pyIsVisibilityOption = ALWAYS` | `pyContainerVisibleWhen = 1=2` | **hidup** — syaratnya sisa |

   Kedua baris terakhir nyata: `DetailLimits` sel 3 `.TreatyGroupID`, dan grid `TreatyIn.CoInScale`
   — yang sempat dinyatakan mati lalu dicabut, sampai gambar 18 (tangkapan layar Pega) memperlihatkannya
   tampil. ⭐ Ke-69 tombol mati di §0 diperiksa ulang terhadap aturan ini: **seluruhnya `OTHER`** —
   angkanya bertahan.

2. **Sapu tombol berdasarkan API aksi, bukan label.** Tombol ikon berlabel `Button`, bukan `Add`/`Delete`.
   Grid `TreatyIn.CoInScale` punya `AddRow` (sel 262) dan `DeleteRow` (sel 266) yang hidup, dan sapuan yang
   menyaring label `Add|Delete|Remove` melewatkan keduanya. Inventaris di lampiran tidak terkena —
   ia mencatat setiap sel ber-`Embed-SelectedContextAPI-*` apa pun labelnya.

## 0 · Angka pokok, dan dari mana ia datang

Penyapu membuang `pyIncludedRuleXML` dengan **hitung kedalaman**, lalu mengurai hasilnya
sebagai XML sungguhan (`ElementTree`) — bukan regex. Saringannya: sel ber-`pyFormat` =
`pxButton` **atau** `pyControlDisplayTitle` = `Button`, **dan** punya sekurangnya satu
`Embed-SelectedContextAPI-*`.

| | Cacah |
| --- | ---: |
| Tombol beraksi, seluruh modul | **203** |
| — **KODE MATI** (`1=2` · `NEVER` · `FALSE &&`) | **69** |
| — khusus pengembang (`OperatorID.pyUserName = 'ALDO SAPUTRA'`, …) | 11 |
| — **hidup** | **123** |

⛔ Jadi **sepertiga tombol di ekspor ini tidak pernah tampil**. Membangunnya berarti
membuat layar baru yang punya tombol yang layar lama tidak punya.

---

## 1 · ⭐ `Save` — jalurnya terlacak sampai ke prosedur Oracle

Ini temuan terpenting ronde ini, dan ia menentukan apakah `Save` dapat dibangun sama sekali.

```
tombol Save            Section/TreatyInActionButtons.xml
   pyCondition         TreatyIn.IsEditData !='1' && TreatyIn.Status…
   ↓ pyActivity
Activity/SaveTreatyIn_Act.xml                                  11 langkah
   1  PROPERTY-SET     TreatyIn.FacultativeShare          := TreatyIn.FacShare
                       TreatyIn.FacultativeShareBrokerage := TreatyIn.FacShareBrokerage
   2  CALL             CheckDuplicateOffer
   3  PROPERTY-SET     InputParam.DATAPEGA := @ASM.GetPageJSONString()   ← SELURUH dokumen
                       InputParam.IDPEGA   := TreatyIn.ID
   4  RDB-LIST         RequestType=SaveTreatyIn · ClassName=ASM-FW-GISFW-Int-TREATY_IN
   5  PROPERTY-SET     Param.IDPEGAOUT := OutputParam.IDPEGAOUT
   7  PROPERTY-SET     OutputParam.ERRMSG := "Terjadi kesalahan saat menyimpan data"
                       when OutputParam.STSSAVE == 0
   8  CALL             SaveTreatyInDetail_Act   when StatusAkseptasi=="Resolve Complete"
   9  CALL             SaveTreatyInOffer_Act    when StatusAkseptasi=="Resolve Complete"
  10  PROPERTY-SET     TreatyIn.ID := Param.IDPEGAOUT
   ↓ RDBList/SaveTreatyIn.xml  (pyBrowseSQL)
BEGIN
  POOLDATA.PEGA_TREATY_IN(
    {InputParam.DATAPEGA}, {TreatyIn.ID}, {TreatyIn.ProportionType},
    {TreatyIn.TreatyContractName}, {TreatyIn.TeritorialScope},
    {TreatyIn.Commencement}, {TreatyIn.Termination}, {TreatyIn.ClassofBusiness},
    {TreatyIn.LeadingReinsSource}, {TreatyIn.LeadingReinsSourceID},
    {TreatyIn.Ceding}, {TreatyIn.CedingID}, {TreatyIn.LeadingReinsID},
    {TreatyIn.NusareSharePct}, {TreatyIn.BrokeragePct}, {TreatyIn.Information},
    {TreatyIn.PositionUsername}, {TreatyIn.Position}, {TreatyIn.StatusAkseptasi},
    {TreatyIn.ChooseStatusAkseptasi}, {TreatyIn.TreatyYear},
    {OutputParam.ERRMSG out}, {OutputParam.IDPEGAOUT out}, {OutputParam.STSSAVE out});
  COMMIT;
END;
```

Catatan ekspor pada aturan itu: *"ganti query insert n update menggunakan procedure pada
flat tabel TreatyIn"*.

### 1.1 Keempat prosedurnya ADA di POOLDATA — diperiksa, bukan diduga

`SELECT OBJECT_NAME, OBJECT_TYPE, STATUS FROM ALL_OBJECTS WHERE OWNER='POOLDATA' AND
OBJECT_NAME LIKE 'PEGA%'` → 100 objek. Yang berkaitan:

| Prosedur | Status | Dipanggil oleh |
| --- | --- | --- |
| `PEGA_TREATY_IN` | **VALID** | `Save` (Treaty In) |
| `PEGA_M_TREATY_IN_DETAIL` | VALID | `SaveTreatyInDetail` |
| `PEGA_M_TREATY_IN_EDM` | VALID | `Save` (Adjustment) |
| `PEGA_M_TREATY_IN_DETAIL_EDM` | VALID | detail Adjustment |
| `PEGA_M_TREATY_IN` | **INVALID** | — tidak dipakai jalur mana pun di atas |

### 1.2 ⛔ DAN DI SINILAH IA TERBENTUR LARANGAN

Isi prosedurnya **dibaca**, bukan diduga — `SELECT TEXT FROM ALL_SOURCE WHERE OWNER =
'POOLDATA' AND NAME = …`, lalu disapu `INSERT INTO` / `UPDATE … SET` / `DELETE FROM` /
`MERGE INTO`:

| Prosedur | Bita sumber | Tabel yang DITULISNYA |
| --- | ---: | --- |
| `PEGA_TREATY_IN` | 6.346 | **`M_TREATY_IN`** · **`POOLDATA.TREATY_IN`** |
| `PEGA_M_TREATY_IN_EDM` | 6.640 | **`M_TREATY_IN_EDM`** · **`POOLDATA.TREATY_IN_EDM`** |

⛔ Jadi `Save` di Pega menulis ke **kedua** tabel yang ronde ini larang: `M_TREATY_IN`
(larangan 2) dan `TREATY_IN` (larangan 5, tabel warisan — dan `TestWarisanHanyaDibaca`
menolak `INSERT`/`UPDATE`/`DELETE` terhadapnya dari naskah SQL modul ini).

### 1.3 Syarat tampil `Save`, lengkap tanpa dipotong

```
Save (Treaty In)   TreatyIn.IsEditData !='1' && TreatyIn.StatusAkseptasi != 'Resolve Complete'
Save (Adjustment)  TreatyIn.ViewState  !='1' && TreatyIn.StatusAkseptasi != 'Resolve Complete'
Save ROL Profile   TreatyIn.ViewState  =='1' && TreatyIn.ProportionType == 'NonProportional'
                   && TreatyIn.StatusAkseptasi = 'Resolve Complete'
                   && OperatorID.pyWorkBasketList(2).pyWorkBasketName = 'ReasTreatyInAdmin'
Close              selalu
Submit             TreatyIn.ViewState !='1' && 1=2        ⛔ KODE MATI
```

⚠️ Jadi `Save` **menghilang** begitu status menjadi `Resolve Complete` — kontrak yang
sudah selesai tidak dapat disimpan ulang lewat tombol itu. Satu-satunya tulisan sesudah
itu adalah `Save ROL Profile`, dan hanya bagi `ReasTreatyInAdmin` pada kontrak
non-proporsional.

### 1.4 Akibatnya

Larangan ronde ini berbunyi: *"`M_TREATY_IN` DILARANG KERAS DIPAKAI APLIKASI"* dan
*"`TREATY_IN` tabel WARISAN"*.

Ketiganya tidak dapat berlaku bersamaan. **Tiga jalan keluar, dan memilih di antaranya
bukan keputusan teknis:**

| Jalan | Artinya | Biayanya |
| --- | --- | --- |
| **A** panggil `POOLDATA.PEGA_TREATY_IN` apa adanya | Save bekerja persis seperti sistem lama, hari itu juga | melanggar larangan 2 dan 5 secara harfiah; aplikasi menulis ke tabel warisan |
| **B** tulis ke tabel pendaratan `T_TREATY_*` | tidak melanggar larangan | ⛔ **tulisannya hilang** — pemuat `MuatKontrak` MENGOSONGKAN lalu menimpa tiap kontrak dari dokumen Pega |
| **C** tulis ke model acuan (`KONTRAK`/`VERSI_KONTRAK`/…) | sasaran akhir yang benar | model itu **nol baris**, dan `RENCANA-PINDAH-SKEMA.md` Tahap 2–7 masih terblokir pada bentuk akar kontrak induk |

⛔ **Tidak dipilih di sini.** Pertanyaannya ada di daftar kebutuhan, butir 1.

⚠️ **Pembaruan 6 Oktober 2026 malam.** `KEPUTUSAN-SASARAN-TULIS.md` (ditulis sesi lain,
16:29) mencatat keputusan pemilik proses yang sama dengan **jalan B**: Save/Submit menulis
ke tabel pendaratan masing-masing, jalur `PEGA_TREATY_IN` ditinggalkan, dan isian tidak
masuk DB sebelum tombolnya ditekan. Berkas ini tidak memutuskan ulang — rujuk berkas itu.
Biaya jalan B di tabel atas (pemuat mengosongkan lalu menimpa) **belum terjawab**; di sana
ia menjadi §2, dan jalur tulis tidak boleh di-`Commit` sebelum salah satu dari tiga
jalannya dipilih.

---

## 2 · Tombol lain, hidup dan mati

### 2.1 Yang HIDUP dan rumusnya sudah terbaca

| Tombol | Syarat tampil | Activity | Isi rumusnya |
| --- | --- | --- | --- |
| `Add` tiap grid Prop | `TreatyIn.IsEditData!='1'` / `ViewState !='1'` | `TreatyInPropAdd` | satu `PROPERTY-SET` per jenis: `TreatyIn.<Larik>(<APPEND>).<ruas> := ""`; parameternya `param.Type` (`"portfolio"`, …) |
| `Delete` tiap grid | sama | — | `Embed-SelectedContextAPI-DeleteRow` murni, nol Activity |
| `Apply` Reporting Period | `TreatyIn.IsEditData !='1'` | `TreatyInSetReport` | sudah dibangun di `services/periode_pelaporan.go`, diukur 99,5% cocok atas 4.548 baris |
| `Refresh` (tab Share prop) | `TreatyIn.ViewState !='1'` | `TreatyInPropshare` | menyusun ulang `Share[]` dari Limits |
| `Add` Fac Retro | `ViewState !='1'` | `AddFacRetroProp` | |
| `Add`/`Remove` pohon Limits | `ViewState !='1'` | `AddValue` + `LimitCalculation` | lihat §2.2 |
| `Add`/`Remove` Deduction | `ViewState !='1'` | `AddDeduction` + `CalculateDeduction` | |
| `Update Summary` (non-prop) | `ViewState !='1'` | `TreatyInNonAddItem` | 41 langkah; lihat §2.3 |
| `Update Summary` (Actual Share) | `ViewState !='1'` | `TreatyInActualUpdateValueShare` | |
| `Close` | selalu | `TreatyInInputVis` | |
| `Edit` / `View` / `Copy` (daftar) | hak akses | `SetTreatyIn_Act` + `TreatyInCheckCedingBlacklist` + `GetMasterTreatyCategory_Act` | |
| `Delete` lampiran | `ViewState !='1' \|\| Revis…` | `Delete_act` | |
| ganti kategori lampiran | `StatusDoc.CARI30` | `ChangeDokument_Act` | |
| `Save ROL Profile` | gabungan `ViewState`/`Pr…` | `TreatyInSaveROL` | |

### 2.2 Pohon Limits — rumus `LimitCalculation` apa adanya

```
langkah 1   keluar bila param.autocalculate == false
QUOTA SHARE (param.kindoftreaty == "qs")
  .RetentionPct := 100 - .QSPct
  .CessionPct   := 100 - .RetentionPct
  RetentionList(<APPEND>).Value := .Value * @divide(RetentionPct, 100, 4)
  CessionList  (<APPEND>).Value := .Value * @divide(CessionPct,   100, 4)
  (Currency dan CurrencyID disalin apa adanya)
SURPLUS (param.kindoftreaty == "surplus")
  .IOOPct     := .Surplus * 100
  .CessionPct := .Surplus * 100
  manual (param.add=="man"):
     IOOLimitList(<APPEND>).Value := .Value * .Surplus
     CessionList (<APPEND>).Value := .Value * .Surplus
  otomatis (param.add!="man" && .Value > 0): sama, lewat loop per TreatyGroup
```

⚠️ `@divide(x, 100, 4)` — **empat angka di belakang koma**, dan itu bagian rumusnya.

`AddValue` hanya menambah baris kosong, satu cabang per `param.type`:
`pla` · `cashloss` · `claimcoop` · `epi` · `reserve` → `.<Jenis>List(<APPEND>).ID := ""`.

### 2.3 `Update Summary` non-prop — `TreatyInNonAddItem`, 41 langkah

Cabang per `param.Type`: `retention` · `egnpi` · `actualpremium` · `limits` · `share` ·
`sharefac`. Yang bukan sekadar menambah baris:

```
egnpi / actualpremium  baris baru MEWARISI mata uang baris retensi PERTAMA:
                       EGNPI(<LAST>).Currency   := TreatyIn.Retention(1).Currency
                       EGNPI(<LAST>).CurrencyID := TreatyIn.Retention(1).CurrencyID
limits                 Limits(<APPEND>).ID := @SizeOfPropertyList(TreatyIn.Limits)
                       lalu CALL CopyLastLimitNP
share / sharefac       PROPERTY-REMOVE share lama lebih dulu
                       galat bila TreatyIn.RNMShare == 0
                       Local.ShareCalc := TreatyIn.RNMShare
                       bila FacultativeShare > 0:
                         Local.ShareCalc        := RNMShare - FacultativeShare
                         TreatyIn.RnmShareDeducted := Local.ShareCalc
                       FacultativeShareList(<LAST>).RnmLimitList.Value
                           := .Limit * @divide(TreatyIn.FacultativeShare, 100, 4)
                       GrossPremiumList.Value
                           := .Value * @divide(TreatyIn.FacultativeShare, 100, 4)
```

⚠️ Langkah 13 bersyarat `.Limit2 < 1` dengan keterangan *"When currency 2 is empty, do not
append it"* — **syaratnya terbalik terhadap keterangannya** di ekspor. Disalin apa adanya;
jangan "dibetulkan" tanpa keputusan.

### 2.4 ⛔ Yang TIDAK boleh dibangun — KODE MATI di Pega

| Tombol | Di seksi | Syaratnya |
| --- | --- | --- |
| **`Download All`** | `ShowAttachmentTreaty` | `FALSE` |
| **`Submit`** | `TreatyInActionButtons` | `TreatyIn.ViewState !='1' && 1=2` |
| **`Submit Revision`** | `TreatyInActionButtons` | `FALSE && …` |
| **`Update Value in Limits`** | `TreatyInTabsProportional` | `NEVER && TreatyMasterInEDM` |
| `Update Summary` | kedua `TreatyInTabsNPValueDifference*` | `FALSE && ViewState !='1'` |
| `Add CoB` / `Delete` | `CoBList`, `DetailLimits` | `ViewState !='1' && 1=2` |
| `Refresh Spreading` | `Share`, `ShareOldData`, `ShareRetro`, `DetailShare`, `TreatyRetroList` | `NEVER` |
| 40 tombol tanpa label | `TreatyInTabsNonProportional`, `…Actual*`, `…AdjustPremi` | `1=2` |

⭐ **Dua di antaranya menutup pekerjaan yang repo ini daftarkan sebagai "rumus belum
dibaca":** `Update Value in Limits` dan `Update Summary` di `TOMBOL_TAMBAH_BELUM_BERGRID`.
Jawabannya bukan "rumusnya belum ketemu" melainkan **"di Pega pun ia tidak pernah tampil"**
— kecuali varian `TreatyInTabsNonProportional`-nya, yang HIDUP (§2.3).

⚠️ **Jebakan yang nyaris memakan korban di sini:** label yang sama hidup di satu seksi dan
mati di seksi lain. `Update Summary` ada empat kali — dua mati, dua hidup. Memvonis dari
labelnya saja akan salah pada separuhnya.

### 2.5 `Upload file` dan `Download All` — keduanya lewat Google Cloud Storage

```
unggah   Activity/TreatySaveAttachment.xml   15 langkah
   5  CALL  ASM-FW-GISFW-INT-T_STORAGE_IMAGE.INSERTGOOGLESTORAGE_ACT
   7  RDB-LIST  (InsertAttachment2_Sql)       when Datain.CARI51==""
   9  CALL  Data-Portal.LoadAttachment
  13  StatusDoc.CARI19 := 1 · StatusDoc.CARI12 := TreatyIn.ID
  14  StatusDoc.CARI10 := @ASM.GetPageJSONString()
  15  RDB-LIST

unduh    Activity/DownloadAll_Act.xml        12 langkah   ⛔ TOMBOLNYA MATI (cond=FALSE)
   6  CALL  GETURLGOOGLESTORAGE_ACT
   8  JAVA  ambil URL → stream → base64
  11  JAVA  ByteArrayOutputStream + ZipOutputStream  (zip di memori)
```

⚠️ Jadi lampiran **tidak disimpan di Oracle**; Oracle hanya memegang metadata dan
tokennya. Repo ini sudah punya lapisan `PELAKSANA_STORAGE` + `STORAGE_TOKEN_SALT`;
apakah sasarannya bucket Google yang SAMA adalah pertanyaan, bukan anggapan — daftar
kebutuhan butir 4.

---

## 3 · Adjustment — nol rute tulis, dan sebabnya sama dengan §1.2

`Save` sisi Adjustment memanggil `SaveTreatyIn_EDM_Act` → RDBList `SaveTreatyInEDM` →
`POOLDATA.PEGA_M_TREATY_IN_EDM(25 masuk, 3 keluar)`, dengan `{InputParam.DATAPEGA}` yang
sama. Tambahan ruasnya `{TreatyIn.OLDID}`, `{InputParam.STSCARI1}`, `{TreatyIn.EDMState}`,
`{TreatyIn.EDMMaterialType}`.

Penghapusannya terbuka dan sederhana:

```
RDBList/RemoveTreatyInEDM.xml     DELETE FROM M_TREATY_IN_EDM WHERE ID = {InputData.CARI1}
RDBList/RemoveTreatyInDetail.xml  DECLARE IDPEGA VARCHAR2(32767); BEGIN
                                    IDPEGA := {TreatyIn.ID};
                                    DELETE FROM TREATYINDETAIL a WHERE a.TREATYID = IDPEGA;
                                    DELETE FROM M_TREATY_IN_DETAIL b WHERE b.JSONDATA.TREATYID = IDPEGA;
                                    COMMIT; END;
```

⛔ Keduanya menyentuh tabel warisan. Terbentur larangan yang sama.
