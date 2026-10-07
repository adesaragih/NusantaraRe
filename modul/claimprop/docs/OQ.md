# OQ — pertanyaan terbuka modul Claim Prop

> Diperbarui 07-10-2026, sesudah implementasi tiket 00–15 (prompt `_brief/PROMPT-IMPLEMENTASI-MODUL-CLAIM-PROP.md`).
> Setiap OQ yang menahan perilaku layar muncul di aplikasi sebagai tombol / grid **nonaktif** berketerangan
> `OQ-CP-nn` (atribut `title`), sehingga paritasnya terbaca di layar dan di `PARITAS.md`.
>
> Pengelompokan per pemilik jawaban: **work owner** (keputusan), **pemilik ekspor Pega** (rule yang tidak ada di
> korpus 329 berkas), **DBA** (objek basis data dan data warisan).

## A. Work owner

| OQ | Pertanyaan | Akibat hari ini | Bukti |
| --- | --- | --- | --- |
| **OQ-CP-16** | ✅ **Selesai 07-10-2026** (opsi "b") — penyerahan ke komite aktif: kode batas komite Claim Prop dikumpulkan di `models/komite.go` dan `repository/komite.go`, keduanya masuk daftar pengecualian `modul/claimlife/backend/services/komite_statik_test.go` (satu-satunya suntingan di modul Claim Life). | Tombol Send Claim to Committee aktif menurut isian; gerbang lampiran tetap menolak selama OQ-CP-12 terbuka. | keputusan work owner |
| **OQ-CP-17** | ✅ **Selesai 07-10-2026** — "1 tabel aja gabung life dan non life": Claim History disimpan di `T_VIEW_SUGGEST` (`CLAIM_ID`, migrasi `532`), satu tabel dengan riwayat penawaran PremiumList Life. Uji PremiumList Life disesuaikan (244 kolom, `kolomModulLain`). | Dibangun. | keputusan work owner |
| **OQ-CP-18** | Pemuat data lama: (1) baris berlaku AC 123 = baris **terakhir** urutan `TANGGAL` naik, `AcceptedNo` turun `[dugaan]` — AC menetapkan urutan, bukan ujungnya; (2) tahap kasus diturunkan dari `STS_REJECT` baris berlaku (0 → Outstanding Claim, 2 → Input Acceptation, 4 → Resolved-Completed) `[dugaan]` — baris OS mencatat aksi, bukan assignment; (3) **419 kasus** berbaris berlaku `STS_REJECT = 1` — ✅ **ditunda** sampai modul Komite Claim Prop (keputusan work owner 07-10-2026: "itu nanti kan dari komite"); pemuat mencatatnya "ditunda", bukan gagal. (4) **2.445 kasus tanpa `JSON_KLAIM`** dimuat sebagai header + tahap saja: `DATA_JSON` baris OS berisi satu baris estimasi / adjustment (24 kunci), dan menyusun ulang `EstimationList` / `AdjustmentList` dari baris OS = pemetaan balik `[dugaan]` — baris itu tetap di `OS_AKSEPTASI_KLAIM` dan terbaca lewat *Summary Outstanding*. AC 132 ("seluruh data klaim dipindahkan") karena itu **sebagian**; susun ulang atau terima? | Uji-kering DEV: 2.032 kasus siap, 419 gagal (seluruhnya sebab butir 3). `-jalankan` belum pernah dijalankan. | uji-kering 07-10-2026 |
| — | ✅ **Selesai 07-10-2026** — tambah / hapus baris Spreading Claim **dinonaktifkan** seluruhnya, termasuk ikon grid bawaan (keputusan work owner "non aktifkan"). Akibat: kasus baru tanpa baris spreading ditolak Save to issue RNM (`ProteksiData_act` 5, "please Fill SpreadingList"). | Nonaktif. | keputusan work owner |
| — | **Lampiran** (prompt §6 butir 7: `inti/backend/unggah` + `T_STORAGE_IMAGE` + `DOCUMENT_CLAIM`) **belum dibangun**: section unggahan (`pyCaseAttachmentsWrapper`) dan sumber kategori lampiran tidak diekspor (OQ-CP-12), dan penyimpanan berkas Claim Prop (Google Storage `ServiceGoogle`) belum disetujui. | Gerbang lampiran Send to Committe selalu menolak. | `InsertGoogleStorage_Act`, `GetUrlGoogleStorage_Act` |
| **OQ-CP-04** | Pemeliharaan master Cause of Loss, Adjuster / Consultant, dan hapus Treaty Group milik menu lain. | Ikon pemeliharaan tampil nonaktif; pemilih tetap jalan. | prompt §2 |
| **OQ-CP-06** | "Close Without Payment" = kasus komite **tanpa** baris adjustment, padahal `T_GENERAL_KOMITE.ADJUSTMENT_ID` NOT NULL; jalurnya di XML juga memuat orang yang di-hardcode. | Tombol **Yes** (tutup tanpa bayar) nonaktif. | `PreventRejectClaimProp`, DDL Komite |
| **OQ-CP-03** | Panggilan nyata layanan REST luar (Kasir, konversi klaim) dari aplikasi. | Efek diantre di outbox hanya bila `IS_PEGA_PROD`; pelaksana berhenti terang (`…BelumDisetujui`). | prompt §6 butir 6 |
| **OQ-CP-13** | Label dan letak pintu pembuatan kasus: harness `New` / `NewSample` tidak diekspor (`pyCanCreateWorkObject=false`). | Tombol **New** di halaman awal = `Start1 → Assignment2` (`Flow_TreatyIn`), satu-satunya bukti. | `Flow/Flow_TreatyIn.xml` |
| **OQ-CP-14** | Varian Syariah (`IsPEGASyariah`) pada muatan Kasir: penanda sumbernya tidak diekspor. | Muatan kasir ber-`StsSyariah` bawaan. | `HitServiceToKasir_Act` 13.1.3 |
| **OQ-CP-15** | Alamat CC/BCC email komite yang di-hardcode XML pindah ke konfigurasi — nama kunci konfigurasinya. | Belum dipakai (email komite ikut OQ-CP-16). | `SendEmailKlaim` |
| — | ✅ **Selesai 07-10-2026** — kotak masuk Komite Claim Life (`komiteclaimlife/backend/repository/komite_inbox.go`) kini menyaring `(w.LINI = 'LIFE' OR w.LINI IS NULL)` (keputusan work owner "tambahkan!"), sehingga kasus komite `TKMT-` lini PROP tidak berdiri di inbox Komite Life. Pembaca satu kasus (`Kasus`, dipakai membuka, memutuskan, dan riwayat) disaring sama: kasus PROP dijawab "tidak ditemukan" ("ya", 07-10-2026). Kotak masuk Claim Life aman sejak awal: TAHAP Claim Prop = nama FlowAction. | Diuji: urutan penanda + isi saringan (uji modul itu), kueri dicoba baca-saja di DEV. | keputusan work owner |
| — | Penegasan `KunciInstans` = ID `T_WORK_CLAIM` apa adanya (`CLMP-000001`) untuk `OS_AKSEPTASI_KLAIM.CASEID` / `JSON_KLAIM.IDPEGA` kasus baru; kasus lama memakai `ASM-FW-GCNMFW-WORK CLMP-n` dan pembaca mencari kedua bentuk. | Diterapkan seperti keputusan NB Treaty In 06-10-2026. | `models/kasus.go` |
| — | Penyimpangan sadar yang perlu dikukuhkan: **write-through** (setiap aksi menyimpan halaman, Pega menyimpan di Save/Submit); `SetDate*` pra-proses hanya mengisi tanggal yang **kosong**. | Berjalan. | `services/aksi.go`, `services/layanan.go` |
| — | Penjaga Claim Life yang menyapu seluruh repo dan **sudah merah sebelum Claim Prop** (salinan bersih HEAD): `TestMigrasiTidakMenyimpanTotalPeserta` (Claim Prop menambah `528 TOTAL_ESTIMATION_VALUE` — potret estimasi AC 40, bukan total peserta) dan `TestKolomTakDibawaHanyaAdaDiKatalog` (Claim Prop **tidak** lagi menambah temuan: `STS_KONVERSI`/`TGL_KONVERSI` dihapus dari INSERT `JSON_KLAIM`, kolomnya nullable tanpa DEFAULT). Pemilik Claim Life mempersempit lingkupnya? | Tidak menahan Claim Prop. | uji bersih 07-10-2026 |

### Butir terbuka prompt §7 (tidak dapat diputuskan dari XML)

| Butir | Keadaan sesudah implementasi |
| --- | --- |
| Nama tabel jejak audit | `T_VIEW_SUGGEST` (keputusan work owner 07-10-2026, "1 tabel aja gabung life dan non life"). |
| `.CashLossList` | **Terjawab dari DEV**: daftar `{Currency, Value}` per `TreatyGroupID` di detail limit treaty; `CashCall` membacanya. Bukan OQ lagi. |
| Induk `DOCUMENT_CLAIM` | Ikut lampiran = OQ-CP-12. |
| `.AktifButton` | **Terjawab dari XML**: ditulis `AddAdjustment_Act` / `DeleteAjsutment_Act`, dibaca `Section/InputAcceptation_Adjs` (tombol Add Adjustment); disimpan di `AKTIF_BUTTON` (520). Bukan OQ lagi. |
| Pengganti `CoverageSubscript` | **Terjawab dari XML**: `AddKomiteTreatyChild_ACT` menyalin `childPageKomite.Adjustment.CoverageSubscript = .IndexCoverage` (indeks posisi); port ini menggantinya dengan penautan stabil `T_GENERAL_KOMITE.ADJUSTMENT_ID` + `T_CLAIM_ADJUSTMENT.KOMITE_ID` (AC 62). |

## B. Pemilik ekspor Pega

| OQ | Rule yang tidak ada di korpus | Akibat hari ini |
| --- | --- | --- |
| **OQ-CP-01** | Activity yang dipanggil tombol tetapi tidak diekspor (lihat kolom "OQ" `PARITAS.md`). | Tombolnya tampil nonaktif. |
| **OQ-CP-02** | Harness yang dirujuk tetapi tidak diekspor (`New`, `NewSample`, pembungkus local action tertentu). | Lihat OQ-CP-13. |
| **OQ-CP-05** | Stream HTML dokumen PLA / DLA / Acceptance Note (`HTMLToPDF`), dan aplikasi belum punya mesin PDF. | Nomor PLA / DLA tetap diterbitkan; berkasnya tidak. |
| **OQ-CP-07** | Prompt values properti bersumber `associated` (`.Type`, `.Payable`, `.ReportType`, `.ReporterStatus`, `.FormType`, `.TypeDeductible`, `.IndividualRiskType`, `.AcceptanceStatus`, jenis estimasi). Bukti label yang ADA di section: `.Type` 1 = Adjustment, 3 = Salvage, 4 = Adjuster Fee, 2 hanya berpasangan "2 atau 4"; 5 dan 6 tanpa label di Claim Prop. | Dropdown menampilkan kode DB apa adanya ("jangan di singkat ikuti apa yang di DB"). |
| **OQ-CP-12** | Sumber `AttachCategory.pxResults.CountAttach` (kategori lampiran) tidak diekspor; unggahan berkas klaim menunggu persetujuan storage. | Gerbang lampiran selalu menolak ("Please Upload Attachment …") — persis perilaku XML tanpa lampiran. |
| — | `KomiteTreaty_Flow` (flow kasus anak komite) tidak diekspor. | Kasus TKMT- ditulis dengan TAHAP `KomiteTreaty_Flow`; alurnya dibangun bersama modul Komite Claim Prop. |

## C. DBA

| OQ | Pertanyaan | Akibat hari ini |
| --- | --- | --- |
| **Migrasi** | Work owner menjalankan `-migrate` untuk **520–533, 980**. Sesudahnya jalankan `go test -tags ujidev -run TestSQLDiDEV ./modul/claimprop/backend/repository/`: 33 pernyataan yang hari ini "objek belum ada" harus menjadi `ok`. | 30 dari 33 menunggu migrasi ini. |
| **OQ-CP-11** | Skema luar tidak terlihat dari akun DEV: `REINSURANCE.TRLOSS_DETAIL_T` (status konversi), `GL.F_GET_EMAIL` (email ceding), dan `ARASAPAS.INVOICE` / `DETAIL_INVOICE` (sinonim publik ada, hak baca tidak). | Tiga bacaan itu tidak teruji di DEV (dua hanya berjalan di produksi). |
| **OQ-CP-10** | Baris roster Direktur Utama di `EMAILKOMITE`: XML mencarinya menurut nama orang (dibuang, AC 55); kunci pengganti `DEGREE = 6` `[dugaan]`. | Batas Direktur Utama dibaca dari baris `DEGREE 6`. |
| **OQ-CP-08** | `PEGA_JSON_OS_AKSEP_KLAIMTRT` tidak ada di DEV (sudah dibuang tiket 00 18-09); konfirmasi tidak ada pembaca lain yang menunggu barisnya. | Jalur KLAIMTRT tidak dijalankan. |
| **OQ-CP-09** | `YearOfQuartal` dibaca dari `JSON_POLIS.DATA_JSON`; polis realisasi sistem baru (NB Treaty In) menulis `JSON_POLIS` **tanpa** `DATA_JSON` → nilainya kosong untuk polis baru. Sumber datarnya? | Kosong untuk polis sistem baru. |
| — | Arti `OS_AKSEPTASI_KLAIM.STS_REJECT = 1` (**1.657 baris**; nol penulis di ekspor Claim Prop). | Lihat OQ-CP-18. |
