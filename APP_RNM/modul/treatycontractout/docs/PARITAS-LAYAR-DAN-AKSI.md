# Paritas layar dan aksi — Treaty Contract Out

> Disusun 28-09-2026 (sesi modul, tiket 01 →). Satu baris per tombol/aksi korpus → rute/kontrol sistem
> baru. Nomor baris = `sed -e 's/></>\n</g'` atas berkas korpus `D:\XML\RNM_BRD\Treaty Contract Out\`.
> Keadaan: ✅ dibangun · ⚠️ belum · ➖ sengaja tidak dibawa (kode mati / penyimpangan sadar) · 🔜 tiket berikut.

> ⛔ **Ralat bertanggal 29-09-2026 — tco4/tco5** `[keputusan work owner]`: seluruh baris di bawah kini menulis/membaca
> **tabel warisan** (nama `T_*` di baris lama dibaca sebagai nama warisannya; migrasi 300–307 dibuang), **nol jejak
> modul** ("+ jejak" di baris lama gugur), security berkunci `(REAS_ID, TRIM(REAS_SECURITY))`, UPDATE business lima
> kolom procedure, lampiran di `M_ATTACHMENTTREATY_2` + `T_STORAGE_IMAGE`, dan menu SATU butir `Treaty Contract Out`
> (ReinsType/Description = popup tombol form, b20778/b22196).

> ⛔ **Ralat bertanggal 29-09-2026 — lanjutan 4** `[asisten dari data DEV; veto work owner]`: tanggal tahun `YYYYMMDD`,
> tanggal reinsurer tidak ditulis (OQ-TCO-01); `USERID`/`TGLUPDATE` reinsurer/business kosong seperti Pega (OQ-TCO-25);
> `GetUrlGoogleStorage_Act` kini juga menjalankan padanan `Update_T_Storage_SQL` (OQ-TCO-26). Baris di bawah yang
> menyebut "stempel 00:00 WIB" atau "diisi layanan" dibaca menurut ralat ini.

## Menu — tiga harness portal (kelas `Data-Portal`)

| Korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `Harness/InboxTreatyContract.xml` b151 `<pyLabel>InboxTreatyContract</pyLabel>` → section `GridTreatyContract` (b931) → `InputTreatyContract` (`GridTreatyContract.xml` b2482) | sidebar **Treaty Contract Out → InboxTreatyContract** | 🔜 tiket 03 |
| `Harness/InboxTreatyContractReinsType.xml` b151 `<pyLabel>InboxTreatyContractReinsType</pyLabel>` → judul `ReinsType` b1670 → section `PanggilReinsType` (b1825) → `InputTreatyContractReinsType` (`PanggilReinsType.xml` b1300) | sidebar **Treaty Contract Out → InboxTreatyContractReinsType** | ✅ tiket 04 (pemilih tahun + editor) |
| `Harness/InboxTreatyContractDescription.xml` b359 `<pyLabel>InboxTreatyContractDescription</pyLabel>` → `NitipKurs` b3882, `SubViewDetailDescription` b11366, `ViewDetailDescriptionProp` b12583, `ViewDetailDescriptionNonProp` b13547 | sidebar **Treaty Contract Out → InboxTreatyContractDescription** (pilih tahun → `PanelKlausulTahun`) | ✅ tiket 08 |

⚠️ Kelompok **Treaty Contract Out** BELUM ada di sidebar (`labels.ts` `MODUL` memuat 17 kelompok; folder
korpus 20 — ralat §8 `PROMPT-EKSEKUSI-HULU-HILIR.md`). Ditambahkan ADITIF bersama layar pertama (tiket 03).

## Tiket 01 — skema + migrasi data (tanpa layar)

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `RDBList/SaveMasterTreatyYear_SQL.xml` → `POOLDATA.PEGA_TREATYYEAR` (10 param) | tabel `TREATYYEAR` (migrasi 300), sequence `TREATYYEAR_SEQ` | ✅ skema; penulisnya 🔜 tiket 03 |
| `RDBList/SaveMasterTreatyContract_SQL.xml` → `PEGA_TREATYCONTRACT` (8 param) | `TREATYCONTRACT` (301) | ✅ skema; 🔜 tiket 04 |
| `RDBList/SaveMasterTreatyReinsurer_SQL.xml` → `PEGA_TREATYREINSURER` (19 param) | `TREATYREINSURER` (302) | ✅ skema; 🔜 tiket 05 |
| `RDBList/InsertToMTreatySecurity.xml` INSERT posisional 7 nilai, `UpdateMTreatySecurity.xml` kunci `trim()` | `MTREATYSECURITY` (303) PK surrogate, kolom bernama, tanpa `trim()` | ✅ skema; ✅ tiket 06 |
| `RDBList/SaveMasterTreatyBusiness_SQL.xml` → `PEGA_TREATYBUSINESS` (12 param) | `TREATYBUSINESS` (304) | ✅ skema; 🔜 tiket 07 |
| `RDBList/SaveMasterProportionalArrg.xml` (35) + `SaveMasterProportionalArrgChild.xml` (26) → satu tabel | `PROPORTIONALARRG` (305) — 35 kolom, anak NULL di 9 kolom induk | ✅ skema; 🔜 tiket 08 |
| — (tidak ada di korpus; ADR-0007) | ~~`T_TREATYCO_JEJAK`~~ (dibuang tco4) | ✅ skema; penulis 🔜 tiket 03+ |
| enam tabel warisan `POOLDATA.TREATYYEAR` … `PROPORTIONALARRG` | `-migrate-data-treaty-contract-out` → `repository.MigrasiTCO.Pindahkan` (baca teks → konversi → tulis → rekonsiliasi tepat → sequence) | ✅ kode + uji `db` (SKIP tanpa skema uji); ⛔ tidak dijalankan di DEV (tco2) |
| kueri hilir `Claim Prop/RDBList/GetLimitPLATreatyin.xml`, `GetListRetro_Sql.xml`, `GetTreatyGroupID.xml`; `Claim Fac In/RDBList/GetLimitPLADLA_Sql.xml`, `GetQuotaShare.xml`, `GetTreatyGroup_Sql.xml`, `GetDataTreatyLimit_Sql.xml`; `Komite Claim Prop/RDBList/GetListRetro_Sql.xml` | `repository.KontrakHilirTCO` — 6 pembaca read-only berkolom VERBATIM; uji tiga sisi (DDL, korpus, nol tulis) | ✅ tco3 |
| `M_PROPORTIONALARRG`, `M_TREATYCONTRACT`, `M_TREATYYEAR`, `M_TREATYBUSINESS` + seluruh kueri `FROM m_*` | — | ➖ MATI (penyimpangan sadar 1); penjaga `TestTCONolTabelDokumenWarisan` |
| `POOLDATA.PROSESCOPY`, `RDBList/SaveMasterCopyData_SQL.xml`, `Activity/BrowseCopyData.xml`, `SaveTreatyYearMultiple_Act.xml`, `NewInputTreatyYearMultiple_Act.xml`, bagian salin `BrowseDeleteRowTreatyInContract.xml`; form From/To `InputTreatyContract.xml` b2374/b3818, tombol `Proces` b5104 → `BrowseCopyData` b5128 | — | ➖ fitur salin DIBUANG (penyimpangan sadar 3, AC 72) |
| kolom warisan `PROPORTIONALLIST`, `OBJECT` | dicacah laporan migrasi (`KolomMatiBerisi`), tidak dibawa | ➖ AC 70 |

## Tiket 02 — jenis reasuransi: master dibaca + saringan non-life

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `ReportDefinition/BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull.xml` — `.ID NotStartsWith` 12 nilai b573/b581, `.Flag="active"` b586, `.Type` 1/2/3 b603, sort `.Note` b667; dipakai 11 grid klausul | `GET /api/treaty-contract-out/jenis-reasuransi` → `repository.MasterJenisReasuransi.DaftarNonLife` (`FLAG=`, `TYPE IN`, `ID NOT LIKE` ×12, `ORDER BY NOTE`) | ✅ nama jujur (nol "Old") |
| master `POOLDATA.REINSURANCETYPE` (`ID`, `NOTE`, `TYPE`, `FLAG`) | dibaca saja; penjaga `TestTCOWarisanHanyaDibaca` (`masterDibacaSajaTCO`) | ✅ AC 12 |
| master kosong / tersaring habis | 503 `ErrMasterJenisReasuransiKosong` menyebut `REINSURANCETYPE` | ✅ ADR-0015 |
| pemilih `ReinsType` (`InputTreatyContractReinsType.xml` b2652) / `Reinsurance Type` (`InputTreatyContract.xml` b7925) | `components/treaty-contract-out/PilihJenisReasuransi.tsx` (`Pilih` ui/dasar; nilai `.ID` teks, label `.Note`) | ✅ komponen; dipasang di layar 🔜 tiket 03/04/08 |
| RD non-Old `BrowseReinsuranceType_RD.xml` (param `Flag="active"` b2768, pemilih form kontrak) | daftar tersaring yang sama (AC 5 tiket 02) | ⚠️ OQ-TCO-06 |
| `RDBList/GetMasterReinsTypeContract.xml` (`a.JSONDATA`) | — | ➖ MATI |

## Tiket 03 — layar tahun treaty (`InboxTreatyContract` → `GridTreatyContract` → `InputTreatyContract` + `InputDtlTreatyContact`)

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| menu `Harness/InboxTreatyContract.xml` b151 | sidebar **Treaty Contract Out → InboxTreatyContract** (`labels.ts` `MODUL.treatyContractOut` +1, `Shell.tsx` kelompok ke-18, `daftarMenu.ts` `tco-tahun`) | ✅ keterangan "belum dimigrasi" dicabut untuk kelompok ini |
| judul `TREATY CONTRACT OUT` (`GridTreatyContract.xml` b1057) | `TAHUN_TCO.judul` | ✅ |
| grid `BrowseTreatyYear_RD` b17355: `Underwriting Year` b17378 / `Transaction Year` b17531 / `StartDate` b17684 / `EndDate` b17837 / `Treaty Group` b17990 / `Reinsurance Type` b18136 (sel b18964/b19125/b19286/b19446/b19605/b19724); sort `.ID DESC` b672 | `GET /api/treaty-contract-out/tahun` → tabel enam kolom, `Halaman` 20/halaman | ✅ (OQ-TCO-05 label bersilang, dibawa apa adanya) |
| `Add` b16387 → `NewInputTreatyYear_Act` | tombol `Add` → `formKosong()` | ✅ |
| `Edit` b19939 → `SetTreatyYear_Act` b20030 | tombol `Edit` → `formDari(baris)` | ✅ |
| `ReinsType` b20778 → harness 2 | tombol membuka `PanelKontrakTahun` untuk baris itu | ✅ tiket 04 |
| `List Description` b22196 → harness 3 | tombol membuka `PanelKlausulTahun` tahun itu | ✅ tiket 08 |
| `Copy` b20459 + `From`/`To` b2374/b3818 + `Proces` b5104 → `BrowseCopyData` | — | ➖ AC 72 |
| form `Input New Data` (`InputDtlTreatyContact.xml` b5437): `ID` b6379 · `Treaty Group` b6560 (`BrowseTreatyGroup_RD` b6624) · `Reinsurance Type` b6800 → `.Proportion` b6829 · `Start Date` b7532 · `End Date` b7816 · `Underwriting Year` b8004 → `.TreatyYear` · `Transaction Year` b8284 → `.UnderwritingYear` · `Modified Date` b9097 · `Username` b9282 | `Field`/`FieldTanggal`/`Pilih` (grup dari `GET /grup-treaty`) / `PilihJenisReasuransi`; ID, Modified Date, Username hanya dibaca | ✅ |
| `Save` b10332 → `SaveTreatyYear_Act` (b280 UserID, b327 TglUpdate, prasyarat b388/b411/b434, RDB `SaveMasterTreatyYear_SQL` → `PEGA_TREATYYEAR`) | `POST /tahun` (baru, ID dari `TREATYYEAR_SEQ`) / `PUT /tahun/{id}` (seluruh medan) → `TREATYYEAR` + jejak `T_TREATYCO_JEJAK`, satu transaksi | ✅ logika ditiru, procedure tidak dipanggil |
| `Cancel` b10622 → `CancelActivityTreatyContract` | tombol `Cancel` menutup form | ✅ |
| `CheckYear` b335 `isNumber(TreatyYear)` (`pyMessageLabel CheckYearly` — teks tidak diekspor) | `ErrTahunTreatyBukanAngka` → 422 | ✅ (teks pesan kosakata kami) |
| — (tidak ada di Pega) | gerbang periode terbalik (AC 9) → 422; anti-dobel (AC 73) → 409 menyebut ID baris lain | ✅ tambahan sadar |
| form kedua `Input New Data` `InputTreatyContract.xml` b6351 (tanpa Start/End/UW Year; `ReinsuranceType` b7954 bukan parameter procedure) | — | ⚠️ `[terbuka]` mana yang tampil di Pega; form lengkap yang dibangun |
| `GridTreatyArrangementAttachment` (`Attachment for` b11721) | `PanelLampiranTahun` di form tahun ber-ID | ✅ tiket 12 |

## Tiket 12 — lampiran tahun treaty (`GridTreatyArrangementAttachment`, FITUR BARU)

Nomor baris = `Section/GridTreatyArrangementAttachment.xml` kecuali disebut lain.

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `InputTreatyContract.xml` b11721 `Attachment for`, b13074 include panel | `PanelLampiranTahun` di form tahun treaty yang sudah ber-ID; tahun baru menampilkan catatan "simpan dulu" | ✅ |
| b1785 `For Treaty Contract Out` | subjudul panel | ✅ |
| `Add attachment` b578 → `SetCategory_act` b596 → `TreatyOutAttachContent` b642 (`pyAttachmentScreen`) → `TreatyOutSaveAttachment` → `InsertAtatchment_Sql` → `PEGA_M_ATTACHMENT` | pemilih `Type` (master `CATEGORY_ATTACH_REAS`) + kotak berkas → `POST /api/treaty-contract-out/tahun/{id}/lampiran` (multipart) → rekam `M_ATTACHMENTTREATY_2` + efek outbox `storage-unggah` + jejak, satu transaksi | ✅ prosedur tidak dipanggil |
| `Refresh` b1023 → `LoadAttachmentTreatyOut` | `GET /tahun/{id}/lampiran` | ✅ |
| `Download All` b2659 → `TreatyOutDownloadAll_Act` | `GET /tahun/{id}/lampiran/semua` (zip) lewat `fetch` berheader identitas | ✅ |
| `Download` b2391 → `DownloadAll_Act` | — | ➖ tombol kedua untuk aksi yang sama |
| sel `.pyFileName` b3428 → `TreatyOutDownloadOne` b3488 | tombol nama berkas → `GET /tahun/{id}/lampiran/{lid}/isi` lewat `fetch` berheader identitas | ✅ |
| kolom `File Name` b3032 · `Type` b3170 (`.pyCategory` b3705) | kolom tabel panel | ✅ |
| `Delete` b3897 → `DeleteAttachmentTreaty` → `DeleteAttachment2_Sql` | `DELETE /tahun/{id}/lampiran/{lid}` → rekam + efek `storage-hapus` + jejak; berkas yang sudah tidak ada tidak menggagalkan | ✅ |
| `TreatyOutSaveAttachment.xml` b376 `Tidak ada file yg diattach` | 400 dengan teks VERBATIM | ✅ |
| `GetAllAttachment2_Sql` / `GetAttachment2_Sql` (`M_ATTACHMENTTREATY_2 where treatyid = {TreatyIn.ID}`) | — | ➖ kunci treaty inward tidak dibawa (penjaga Go + JS) |
| — (tidak ada di Pega) | status terkirim / tertunda / gagal + galat terakhir; `Ulangi` (`POST …/{lid}/ulangi`); `Periksa keselarasan` (`GET …/selaras`) | ✅ tambahan AC 55, 58, 61 |
| `ConnectREST/ServiceGoogle.xml`, `LinkService`, `GetTokenStorage_SQL`, `InsertGoogleStorage_Act` / `GetUrlGoogleStorage_Act` / `DeleteGoogleStorage_Act` | `PenyimpananLampiranTCO`: bawaan stub lokal `PenyimpananLokalTCO`; `PELAKSANA_STORAGE=nyata` → `PenyimpananJarakJauhTCO` (resolver `M_LINK_SERVICE` saat jalan, `CacheTokenTCO`, token `GCP_IMAGE`/garam env) + transport HTTP `UploadDoc` (delete dengan jalur objek penuh, b1091) | ✅ OQ-TCO-08 (keputusan work owner 29-09-2026); `Folder`/`Durasi` OQ-TCO-22 |
| — (tidak ada di Pega) | pekerja latar antrean lampiran `JalankanPekerja` dari `cmd/api`, interval `TCO_PEKERJA_LAMPIRAN_INTERVAL` (bawaan mati) | ✅ OQ-TCO-09 |

## Tiket 04 — kontrak treaty di dalam tahun (`InputTreatyContractReinsType`)

Nomor baris = `Section/InputTreatyContractReinsType.xml` kecuali disebut lain.

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| harness `InboxTreatyContractReinsType` b151, judul `ReinsType` b1670 | butir menu kedua → `InboxTreatyContractReinsType.tsx` (pemilih tahun, `[tidak ada di korpus]`) → `PanelKontrakTahun` | ✅ |
| popup dari `InputTreatyContract.xml` b20778 (`BrowseReinsTypeYear` + `showHarness` + `KirimTahunGroupID`) | tombol `ReinsType` baris tahun → `PanelKontrakTahun` tahun itu | ✅ |
| kepala `Underwriting Year` b1145 · `ReinsType` b1358 (nilai `.TreatyGroupName` b1386) | dua medan baca-saja + catatan label bersilang | ✅ |
| `ReinsType` b2652 (RD non-Old, `Flag "active"` b2768) | `PilihJenisReasuransi` (daftar tersaring tiket 02) | ✅ OQ-TCO-06 tetap terbuka |
| `Start Date` b2905 → `SetTanggalTreatyContract` b3007 | isi tanggal mulai → `GET /tahun/{id}/kontrak/akhir-bawaan` mengisi tanggal akhir = mulai + 1 tahun kalender (`ADD_MONTHS(…,12)`) | ✅ OQ-TCO-10 — penyimpangan sadar (keputusan work owner 29-09-2026) |
| `End Date` b3244 (`IsEndDate==1`) | ubah tanggal akhir saja | ✅ |
| `Save` b3618 → `SaveTreatyContract_Act` → `SaveMasterTreatyContract_SQL` → `PEGA_TREATYCONTRACT` | `POST /tahun/{id}/kontrak` / `PUT /tahun/{id}/kontrak/{kid}` → `TREATYCONTRACT` + jejak, satu transaksi | ✅ prosedur tidak dipanggil |
| `Undo` b5343 → `UndoOperation` | kembalikan isian terakhir yang dimuat | ✅ |
| `Information` b6400 (`OutputData.HASIL1`) | baris status sesudah simpan | ✅ |
| `Modified Date` b4363 · `Username` b4547 · ID b2478 | baca-saja (Username = penulis terakhir rekam) | ✅ ralat 7 |
| `Add` b8528 → `NewInputTreatyContract_Act` | form kosong | ✅ |
| grid `BrowseTreatyContract_RD` (tidak diekspor): `Reins Type` b9164 · `Treaty Start` b9304 · `Treaty End` b9444 | `GET /tahun/{id}/kontrak`, `ID DESC` | ✅ |
| `Edit` b10519 → `SetUbahTreatyContract` | form dari baris | ✅ |
| `Business List` b10842 | membuka `PanelBusinessKombinasi` | ✅ tiket 07 |
| `Reinsurer List` b11308 | membuka `PanelReinsurerKombinasi` | ✅ tiket 05 |
| `Delete` b11809 → `BrowseDeleteRowTreatyInContract` | popup Ya/Batal berjumlah → kaskade satu transaksi | ✅ tiket 10 |
| `ViewDetailTreatyReinsurerGrid1` b13311 | `PanelReinsurerKombinasi` | ✅ tiket 05 |
| `ViewDetailTreatyBusinessGrid` b14064 | `PanelBusinessKombinasi` | ✅ tiket 07 |
| grid security `SelectSecurityReinsurer` b16069 | `PanelSecurityReinsurer` | ✅ tiket 06 |
| gerbang `SaveTreatyContract_Act` langkah 2–6, 12 (DIKOMENTARI) | jenis wajib dari daftar, tanggal wajib, periode (AC 9), dobel 409 `Data sudah pernah di Input` | ✅ ralat 1–2, OQ-TCO-11 |
| gerbang `SetTanggalTreatyContract` langkah 4 (HIDUP) | tahun mulai = tahun treaty → 422 | ✅ |

## Tiket 05 — reinsurer pada kombinasi (`ViewDetailTreatyReinsurerGrid1`)

Nomor baris = `Section/ViewDetailTreatyReinsurerGrid1.xml`.

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `Add` b1730 → `NewTreatyReinsurerDetail_Act` | form kosong | ✅ |
| `Tambah` b1996 (aksi sama) | — | ➖ tombol kedua untuk aksi yang sama |
| grid `BrowseDetailTreatyReisurer_RD` (kombinasi, `ID ASC`): `ReinsID` b2418 · `Reinsurer` b2560 · `%Share` b2702 · `%Comm` b2844 · `Rating` b2990 · `Operator Name` b3138 | `GET /tahun/{id}/kontrak/{kid}/reinsurer` | ✅ |
| `Total Share -->>` b6186 (`InputTreatyReinsurer.TotalShare`) | kaki tabel, desimal persis dari server | ✅ AC 15 |
| `Edit` b4491 → `SetUbahTreatyReinsurerList_Act` | form dari baris | ✅ |
| `Delete` b4936 → `DeleteTreatyReins_Act` | popup Ya/Batal (jumlah security) → hapus satu transaksi | ✅ tiket 10 |
| `Security Reinsurer` b5277 | tombol membuka `PanelSecurityReinsurer` reinsurer itu | ✅ tiket 06 |
| form `ID` b7842 · `Reins.ID` b8042 · `Reinsurer` b8226 (pemilih `BrowseAgentReinsSOA_RD`) · `%Share` b8522 · `%Comm` b8800 · `Rating` b9076 · `Operator Name` b11100 | ID/Reins.ID/Operator Name baca-saja; kotak cari + pemilih master aktif | ✅ OQ-TCO-12 |
| delapan medan tersembunyi `pyCondition 1=2` | tidak diterima dari klien; dipertahankan server | ✅ |
| `%Share`/`%Comm` → `SetErrorMessageReinsurer` (koma → titik, 0..100) | `models.UraiPersenMasukTCO` di batas masukan | ✅ |
| `Save` b11405 → `SaveTreatyReinsurerDetail1_Act` → `SaveMasterTreatyReinsurer_SQL` → `PEGA_TREATYREINSURER` | `POST`/`PUT .../reinsurer` → `TREATYREINSURER` + jejak; total > 100 → 422 VERBATIM | ✅ prosedur tidak dipanggil |
| `Error` b12131 · `Informasi` b12868 | pita galat + baris status | ✅ |

## Tiket 07 — business pada kombinasi (`ViewDetailTreatyBusinessGrid`)

Nomor baris = `Section/ViewDetailTreatyBusinessGrid.xml`.

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| kepala `Business List` b1988 (`.ReinsTypeName` b2019) | judul panel | ✅ |
| `Add` b2785 → `NewTreatyBusinessDetail_Act` | form kosong | ✅ |
| grid `BrowseTreatyBusiness_RD` (saringan `.IsActive = 1`): `Treaty Group` b3422 · `Business ID` b3566 (ID baris) · `Business Name` b3710 | `GET .../business` — SELURUH baris + kolom `Active` | ✅ ralat 1 (AC 22) |
| `Edit` b4547 → `SetUbahTreatyBusinessList_Act` | form dari baris | ✅ |
| `Delete` b4826 → `DeleteRowBusiness` → `DeleteRowBusinessList` (dua SQL) | `DELETE .../business/{bid}` satu tabel; pesan `Data Dengan ID … Berhasil di Hapus` | ✅ AC 63/64 |
| `Business Name` b6241 (pemilih `BrowseFilterBusiness_RD`) | `Pilih` dari `GET /business-master` | ✅ OQ-TCO-13 |
| `Active` b6499 (radio wajib) | `Pilih` Aktif/Nonaktif (`1`/`0`) | ✅ OQ-TCO-13 |
| `Business Code` b6680 (tersembunyi) | baca-saja | ✅ |
| `Save` b6966 → `SaveTreatyBusinessDetail_Act` → `SaveMasterTreatyBusiness_SQL` → `PEGA_TREATYBUSINESS` | `POST`/`PUT .../business` → `TREATYBUSINESS` SELURUH medan + jejak | ✅ AC 23, prosedur tidak dipanggil |
| `Information` b9319/b10064 (`ERRMSG4` "Data sudah pernah di Input") | baris status; 409 dobel | ✅ ralat 4 |
| `Close List` b10889 → `CancelActivity` | tombol tutup panel | ✅ |

## Tiket 08 — klausul (`InboxTreatyContractDescription` → 25 section `GridTreatyArrangement*`)

Nomor baris = `Harness/InboxTreatyContractDescription.xml` kecuali disebut lain.

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| kepala b2003–b3382 (tujuh medan) | tujuh `Field` baca-saja | ✅ OQ-TCO-05 |
| `For Non XOL` b4880 / `For XOL` b7971 (`BrowseTreatyDesc_RD`, `IsXOL` 0/1) · `ID` b5440 · `Description Name` b5549 | `GET /jenis-klausul?isXol=` dari `TREATYDESC` (baca-saja) | ✅ AC 26/27 |
| `Show` b6059 → `BrowseDescriptionLimit` + `testingKurs` + `SetKirimIDDesc` + `PanggilID` | membuka `PanelJenisKlausul`; beberapa jenis boleh terbuka bersamaan | ✅ AC 29; kurs ✅ tiket 11 |
| 18 section induk: `Add` / `Edit` / `Save` (`GridTreatyArrangementEpi.xml` b8980/b10917/b5501) → `SaveTreatyArr*_Act` → `SaveMasterProportionalArrg` | `POST`/`PUT /tahun/{id}/klausul` → `PROPORTIONALARRG`, aturan per jenis dari server | ✅ AC 24/33/35, prosedur tidak dipanggil |
| 7 section anak: `Show Child` b11200 → `Browse*ParentList`; `Close Child` (`GridTreatyArrTreatyEpiList.xml` b8657) → `Save*List_Act` → `SaveMasterProportionalArrgChild` | grid anak per induk; Rp/Usd turunan; total Pct + peringatan; ReinsType anak ketujuh grid = `TreatyContractSetReinsTypeList` atas nama ReinsType induk (`QS`→QS (OR)/QS (R/I)/ORS, `SPL`→SPL (OR)/SPL (RI)/ORS, `XOL`→QS (OR)/QS (R/I)/XL, `ORS`→ORS), dapat difilter [keputusan work owner 02-10-2026] | ✅ AC 25/28 |
| `HitungRpUsd` | `RpUsdAnakTCO` di server | ✅ |
| `TreatyTestChildTotal_Act` (TreatyLimitChild) | peringatan `Please make sure spreading is 100%` sesudah simpan | ✅ ralat 4 |
| ExclutionTreaty: empat sub-bagian (`Param.Type`) + pemilih `BrowseOccupationFIRE_RD` / `BrowseFireClauseFacIn_RD` | satu jenis, empat subjenis — tampil sebagai **tab** (`StripTab`); ID Occupation / ID Clause = satu dropdown yang dapat difilter (`PilihMasterKlausul`, tanpa kotak Search) atas `GET /klausul-pilihan/{occupation,clause}` [keputusan work owner 02-10-2026] | ✅ OQ-TCO-14/16 |
| CoinsPanel `GridTreatyArrangementCoins.xml`: judul `Co-Ins Scale` b917; dua grid `Risk with Sum Insured less than` b8708 / `more than USD 100.000.000` b12971 (`Param.Type` "Less Than" b10055 / "More Than" b14320); kolom `Co Insurance Share` (`DetailCoinsShare`: `.CoIns_Min` - `.CoIns_Max`) + `Treaty Limit`; `Add` b10554/b14815 → `NewTreatyArrCoins` (`.SpreadingOrder = Param.Type` b406); `Edit` b11254 → `SetTreatyArrExclustionCoins_Act`; RD `BrowseTreatyArrangement_CoinsPanel_RD` `.SpreadingOrder = Param.Type` b610/b612, `.TreatyLimit DESC` b653/b655 | satu jenis, dua subjenis `Less Than` / `More Than` disimpan server di `SPREADINGORDER` (klien tidak mengirimnya); dua grid **bertumpuk** yang dapat dilipat, `Add` di kepala kolom aksi, `No items` di dalam tabel, urut `TreatyLimit` turun; baris tidak pindah grid; kunci dobel per grid `[asumsi]` [keputusan work owner 02-10-2026: tampilan seperti gambar] | ✅ |
| `SaveTreatyArrLimitMB_Act`, `SaveTreatyArrPortfolio_Act` | ditahan: 422 + alasan di layar | ⏸ AC 36 (Product + UW) |
| 16 `CancelActivity*` | `Cancel` per panel membuang isian panel itu saja | ✅ AC 29 |
| `NitipKurs` b3882 / `testingKurs` | baris kurs berlaku di panel jenis berkurs (`GET /tahun/{id}/kurs`) | ✅ tiket 11 |

## Tiket 06 — security di bawah reinsurer (`InputTreatyContractReinsType` bagian `HASILD21`)

Nomor baris = `Section/InputTreatyContractReinsType.xml`.

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `Security Reinsurer` (`ViewDetailTreatyReinsurerGrid1.xml` b5277) → `SetSecurityReinsurer` | tombol baris reinsurer → `PanelSecurityReinsurer` | ✅ |
| grid: `Reas Security` b16088 · `Security Name` b16228 · `Percent Share` b16368 | `GET .../reinsurer/{rid}/security` (nama dari `AGENT`) | ✅ ralat 5 |
| `Add` b15459 → `InputNewSecurityReinsurer` | form kosong | ✅ |
| `Edit` b17252 → `ShowEditSecurityReinsurer` | form dari baris (ID tetap) | ✅ |
| `Delete` b17559 → `DeleteSecurityReinsurer` (kunci `trim(nama)`) | `DELETE .../security/{sid}` satu ID | ✅ ralat 2 |
| form `Security ID` b19468 (nonaktif) · `Security Name` b19648 (pemilih `BrowseAgentReinsSOA_RD`) · `%Share` b19888 | baca-saja · `Pilih` dari `GET /reinsurer-master` · teks desimal wajib 0..100 | ✅ ralat 4, OQ-TCO-17 |
| `Save` b20246 → `SaveSecurityReinsurer_Act` → `InsertToMTreatySecurity` / `UpdateMTreatySecurity` | `POST`/`PUT .../security` → `MTREATYSECURITY` kolom bernama + jejak | ✅ AC 18–20, ralat 1/3 |
| `Error` b20980 · `Informasi` b21717 | galat / baris status | ✅ |
| `DeleteTreatyReins_Act` → `DeleteFromTreatyReinsurer_Act` | security lalu reinsurer, eksplisit + FK `ON DELETE CASCADE` | ✅ tiket 10 |

## Tiket 11 — kurs USD → IDR (`testingKurs`, `HitungRpUsd_depan`, `CalculateTSIExcludeTreaty`)

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `testingKurs` (Show Non XOL b6138 `StartDate`; Show XOL b9256 `TreatyYear`) → `GetMasterKursList` | `GET /tahun/{id}/kurs` — StartDate tahun untuk kedua grid | ✅ ralat 1 |
| `IDCURRENCY = '10001'` literal di SQL | pengenal dari master `CURRENCY` lewat kode `USD` | ✅ AC 47 |
| `HitungRpUsd_depan` onchange Rp tujuh form induk; `Usd` hanya dibaca | `Usd` turunan server `Rp ÷ Kurs` (skala 8); pratinjau `GET /kurs/konversi` | ✅ ralat 3/4 |
| `CalculateTSIExcludeTreaty` (exclusion Occupation, `Curr` IDR/USD) | pratinjau dua arah dari server; kedua nilai tersimpan terpisah | ✅ AC 49 |
| 14 `NewTreatyArr*`: "Tidak ada Nilai Kurs di Tahun : " + TreatyYear, form tidak tampil | 422 VERBATIM saat simpan; `Add` nonaktif + pesan di panel | ✅ ADR-0015 |
| `RefreshKurs` (mengosongkan `Kurs`) | tidak dibawa: kurs dibaca ulang tiap panel dibuka | ✅ |
| `SetTreatyArrangementDesc_Act` (5/6 langkah di-remark) | tidak dibawa | ✅ AC 50 |

## Tiket 09 — simpan atomik lintas enam tabel

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `COMMIT;` di enam rule Connect-SQL (mis. `SaveMasterProportionalArrg.xml` b123) | nol COMMIT di teks SQL; commit sekali per permintaan | ✅ |
| `Save` per panel (kontrak, reinsurer, security, business, 25 klausul), masing-masing COMMIT | tetap per panel, masing-masing satu transaksi + jejak | ✅ paritas |
| — (tidak ada di Pega) | ~~`POST /tahun/{id}/kontrak-utuh`, `PUT /tahun/{id}/kontrak/{kid}/utuh`~~ — dibuang (OQ-TCO-19, keputusan work owner 29-09-2026) | ✖ wontfix |
| — | tombol simpan tunggal di layar | ✖ tidak perlu (OQ-TCO-19) |
| `StsSimpan` 1/0 `[data DBA]` | prosedur tidak dipanggil; jawaban simpan utuh dibuang (OQ-TCO-19) | ✖ wontfix |

## Tiket 10 — kaskade hapus, popup, klausul yang tetap hidup

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `Delete` kontrak b11809 (tanpa konfirmasi) | `GET .../kontrak/{kid}/dampak-hapus` → popup (reinsurer, security, business; klausul TIDAK terhapus) → `DELETE .../kontrak/{kid}?reinsurer=&security=&business=` | ✅ AC 42/43, penyimpangan sadar 4 |
| `DeleteFromTREATYCONTRACT_SQL` b80–b94 (empat DELETE + COMMIT) | empat DELETE di `T_`, satu transaksi, nol COMMIT; bisnis tahan `TREATYYEARID` NULL | ✅ AC 45, 63/64 |
| `PROPORTIONALARRG` tidak ikut | klausul tidak disentuh; jumlahnya (dari induknya, OQ-TCO-20) masih dihitung server tetapi tidak lagi disebut di popup (keputusan work owner 01-10-2026); langkah kaskade diuji tanpa klausul | ✅ AC 44 |
| kaskade menghapus seluruh anak kombinasi walau dipakai kontrak tahun lain | seperti Pega (OQ-TCO-21); popup memperingatkan cacah kontrak lain, cacahnya ikut dikonfirmasi | ✅ keputusan work owner 29-09-2026 |
| `Data Berhasil di Hapus` b762 | pesan sukses VERBATIM | ✅ |
| `Delete` reinsurer b4936 → `DeleteFromTreatyReinsurer_Act` | popup (jumlah security) → `DELETE .../reinsurer/{rid}?security=` | ✅ |
| `DetailTreatyExclustion` → `DetailTreatyExclustion_Sec` (panel rinci exclusion Occupation, `pyEditAction` b8801) | form baris exclusion di `PanelJenisKlausul`; tetap hidup sesudah kontrak dihapus | ✅ |
| jalur salin `SaveMasterCopyData_SQL` (b922, di-remark) | tidak dibawa (AC 72) | ✅ |
