# PROMPT — GILIRAN 3 · SESI C **Komite Claim Life** *(folder `OUTPUT_HASIL_RNM\.worktrees\komite-claim-life`, cabang `modul/komite-claim-life` @ `4f40272`)*: **ralat "ShowTransfer mati" → tiket 01 → 09 seluruhnya, dari XML, dalam satu giliran**

> Baca `PROMPT-INDUK-TIGA-MODUL.md` *(§0.4)*, `PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-LIFE.md` *(§0–§3, km1–km5 tetap berlaku
> kecuali yang diralat di §1 di bawah)*, tiket `.scratch/komite-claim-life/issues/`. Tiket **00 sudah tuntas dan menyatu**
> *(`7709c0c` → `df1353d`; migrasi `030`)*. Mekanisme giliran = lanjutan 8 §1. Sesi ini **tidak** merge ke `main` — asisten
> menyatukan sesudah memverifikasi log. Bila sesi tunggal: kerjakan di `main` sesudah brief B.

---

## 0. KONTEKS YANG SUDAH PASTI

| Hal | Isi |
| --- | --- |
| Cabang dan port | `modul/komite-claim-life` @ `4f40272` *(fast-forward asisten 28-09-2026)*; backend `:8083`, Vite `5175`; `.env` ada |
| Katalog DEV *(27-09-2026)* | `T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST` **0 baris** → `030 MODIFY` aman; `T_MIGRASI` belum memuat `030` *(work owner)* |
| Keputusan | km1–km5; decision o; ADR-U-0029; kontrak masuk/keluar dengan Claim Life = A2 *(`serahkanKomite` ada di `main`: `POST /api/klaim-life/{id}/peserta/{pesertaId}/adjustment/{adjId}/komite`)*; `STS_REJECT` ke baris Claim Life lewat fungsi `services` yang **ada** *(km3)* |
| Berkas milik sesi ini | `komite_*`, `handlers/rute_komite.go`, `frontend/src/pages/komite/`, `labels.komite.ts`; Shell hanya **ditambah** kelompok menu; berkas Claim Life/PremiumList **tidak** disunting |
| Migrasi | `031`+ hanya dari keputusan tercatat; `-migrate` **tidak** dijalankan executor |
| Menu `[dari bukti]` | kelompok **Komite Claim Life** → satu butir **`Inbox Komite`** `[tidak ada di korpus — worklist `KomiteRouter`]`. Layar keputusan dibuka **dari baris** |

## 1. ⛔ RALAT — `ShowTransfer` BUKAN kode mati; ia LAYAR KEPUTUSAN

`FlowAction/ViewTransferDtl.xml` menampilkan **`ShowTransfer`** *(b91, b165 `flowActionUIRef`, b322, b563)*. Yang mati hanyalah
**bagian bergerbang `.TransferType==2` b29177** di dalamnya. Ralat bertanggal di `spec.md` *("kode mati dibuang: `TransferType`
+ `ShowTransfer`")*, brief modul §1 baris 3, dan tiket 02/09 bila mengulanginya. Isi `ShowTransfer` *(dibaca 28-09-2026)*:

| Unsur | Baris | Isi |
| --- | --- | --- |
| Medan keputusan | b30623 `.KomiteAproval` *(dropdown — sumber daftarnya dibaca; bila tidak ada di ekspor → `[terbuka]` seperti OQ-L, nilai tidak dikarang)*, b31001 `.KomiteComment`, b27860 `Remarks` | wajib diisi sesuai `pyRequired`/syarat yang tertulis |
| Tombol | `Submit` b34722 → `finishAssignment` b34732/b34844 *(→ activity flow action `KomitePostAdjustment`)*; `Cancel` b33880; `Refresh` b22882; `Add attachment` b23511 → local action `AttachDocumentLife` b23539/b23731 *(→ `SaveAttachLife`)*; tautan `DownloadDocumentClaim` b25128/b25366/b25646/b25804; `View Office Online` b25629 *(syarat b25896: MIME xls/xlsx/doc/docx/ppt/pptx **dan** `T_STORAGE_ID != ''`)*; `SetRemarkKomiteLife` b32664/b32800 | VERBATIM |
| Syarat tampil | `pyWorkCover.IsTreatyIn==0` ×12 *(b3019…b9559)*; `pyWorkCover.PolicyDataLife.Type=='TP' \|\| 'TR'` ×3 *(b15429, b16998, b17166)*; `AdjustmentDetail.SwiftCode != ''` b18945; `.RetrocadedShare!=0` b21314; `.TransferType==2` b29177 *(mati)* | tiap syarat → gerbang di services **dan** tampilan |

**Flow `KomiteLife_Flow.xml`**: Start2 → Assignment1 `[Always]` b970/b978; Assignment1 → Decision1 lewat **`ViewTransferDtl`**
b763/b769; Decision1 → Assignment1 `IsKomiteLoop` b906/b914 *(putaran tingkat berikut)*, → End1 `NoLoop` b834/b812;
utility `KomiteLoop` b538, router `KomiteRouter` b588 *(`pyRouteTo Custom` b619)*; `pyWorkTypeName KomiteLife` b220.
Peta konektor utuh ditulis di tiket 02 sebelum kode.

**`KomitePostAdjustment.xml`** *(activity `Submit`)*: `Obj-Open-By-Handle` b456 → `TempOpenPage` *(klaim induk)*; `Property-Set`
b596/b772; prasyarat `…AdjustmentList(Local.AdjusmentID).IsSubjectivity == true` b988; `RDB-List` b1077 *(`TglProd`)*, b1712, b1931,
b2391, b2736, b5184; cabang `PolicyDataLife.Type=="QP"||"QR"` b1883/b3334 lawan `"TP"||"TR"` b2100/b3612/b5098; gerbang
`.ACCEPTEDNO=="" && .IsCheck = true && .STS_REJECT == "0"` b4896/b5341; `SecurityReinsurerID!=""` b5121; **tanggal potong warisan**
`OfferFacIn.PolicyData.ProdDateTime<"20250207T000000.000 GMT"` b5144 *(ditiru sebagai konstanta bertanggal + dicatat OQ, bukan
dihapus)*. Dibaca **utuh** *(langkah 1–akhir, tiap `RDB-List` → SQL-nya di `RDBList/`)* sebelum tiket 02–04b.

Activity lain *(18)*: `DownloadDocumentClaim`, `GetLinkService`, `GetUrlGoogleStorage_Act`, `HitServiceToKasirKMTLife_Act`,
`InsertDocument_Act`, `InsertGoogleStorage_Act`, `InsertJsonClaimLife_Act`, `KomiteRouter`, `LoadDocumentKomiteLife_ACT`,
`LoadDocumentLife_ACT`, `NewAttachLife`, `PreKomiteLife`, `PrintAkseptasiPDF`, `SaveAttachLife`, `SendEmailKlaimLife`,
`SetInformationData`, `SetRemarkKomiteLife`. RD: `GetEmailUser_RD`, `GetAllAttachmentsLife`, `BrowseDiseaseLife_RD`.
Dokumen komite memakai jalur dokumen Claim Life yang **sudah ada di `main`** *(be: `T_CLAIMLF_DOCUMENT`, `T_CLAIMLF_STORAGE`,
outbox + stub)* — dipanggil, tidak disalin.

## 2. URUTAN — satu commit per tiket `komite-claim-life: tiket NN — <judul>`, tanpa pesan di antaranya

| Tiket | Layar/rule XML yang menjadi pohonnya |
| --- | --- |
| **01** terima kasus + inbox per posisi | `serahkanKomite` *(A2)* → baris `T_WORK_CLAIM` `KMTLF-`; `KomiteRouter` *(worklist per tingkat roster, km2)*; **Inbox Komite** kolom Pega-standar *(dipinjam `InboxPremiumList`, ditandai)* + nomor klaim induk, tingkat `KomiteLoop`, nilai + `CURRENCY` |
| **02** mesin tangga | peta konektor `KomiteLife_Flow`; `KomiteLoop`/`IsKomiteLoop`; `KomitePostAdjustment` utuh; layar **`ShowTransfer`** *(§1)* dengan `KomiteAproval`/`KomiteComment`/`Remarks`, `Submit`, `Cancel`, syarat tampil |
| **03** wewenang + eskalasi | roster `EMAILKOMITE` *(dibaca saat jalan; nol baris disalin)*; eskalasi manual satu tingkat — satu-satunya pengecualian, diuji |
| **04a** nomor akseptasi sekali di tingkat final | `GetAcceptedNoCL`, counter `(ASM-FW-GCNMFW-Work-KomiteLife, RNML-A)` ditiru lewat `PenomorCounter` *(o2/o3)*; keunikan lintas jalur dengan `ACCEPTATIONNOLIFE_SEQ` Claim Life diuji; bila XML tidak menjawab → `[terbuka — work owner]` berbukti |
| **04b** rekam akseptasi satu jalur | `InsertJsonClaimLife_Act`, `PrintAkseptasiPDF`, `SetInformationData`; `OS_AKSEPTASI_KLAIM_LIFE` ditulis rule yang sama dari dua sisi |
| **05** jalur balik `STS_REJECT` dua tingkat | km3: fungsi `services` Claim Life yang ada *(dipanggil)*; kaskade 030 dipakai |
| **06** outbox | `HitServiceToKasirKMTLife_Act`, `SendEmailKlaimLife`, `GetLinkService` → `T_LOG_SERVICE_RNM` *(km4)* |
| **07** worker retry anti-dobel | kunci idempoten; pelaksana stub di DEV |
| **08** status perlu intervensi + laporan harian | kata di UI, kode di `models` *(km5)* |
| **09** riwayat tangga di UI | `T_KOMITE_KOMITELIST` + jejak; label VERBATIM dari `ShowTransfer` |

## 3. LAPORAN · TELEMETRI

Satu pesan di akhir *(atau batas tiket bila konteks menipis; pohon bersih; SHA)*: tabel **tiket → commit → menu/tombol XML →
rute/kontrol**; ralat *(termasuk §1)*; kontrak yang dipakai/diusulkan *(`KONTRAK-USULAN.md` bila perlu berubah)*; kode bersama
aditif; angka uji tiap commit; bab **TELEMETRI EKSEKUSI** per tiket.

---

*Disusun 28 September 2026 dari `KomiteLife_Flow.xml`, `FlowAction/ViewTransferDtl.xml`, `Section/ShowTransfer.xml` (tombol,
medan, syarat), `Activity/KomitePostAdjustment.xml` (langkah, prasyarat), daftar 18 activity dan 3 RD, tiket 00–09, dan katalog
DEV (cacah tabel komite).*
