> Modul  : Komite Claim Non Prop · Tahap spesifikasi · 2026-09-21
> Peran  : auditor
> Ronde  : penutupan spesifikasi, 2026-09-21 — §7 ditambahkan; empat rule masuk §2.4
> Masukan: 338 XML ekspor Pega · 49 DDL · `PENGETAHUAN.md` · `PUTUSAN-01.md` · `FAKTA-A1B.md` · `GRILL-05/*` · `GRILL-06/*` · penyuntingan penutup ronde 6 · `SPEC-KOMITE-01.md` dan `AUDIT-SPEC-01.md` (2026-09-21)
> Status : TERBUKA
> Sifat  : HIDUP

# INVENTARIS BUKTI

> Disusun 20 September 2026, ronde 4. Seluruh cacah dihitung dari berkas, bukan dari ingatan.
> Sebab berkas ini ada: premis pembekuan A-1b ternyata salah — rule pemagarnya sudah di repo sebelum pagarnya dipasang. Itu kejadian **kelima** dari pola yang sama (rincian di §3).

---

## 1. Kumpulan yang dipegang

| Kumpulan | Lokasi | Cacah | Asal & tanggal | Apa yang diliput | Apa yang TIDAK diliput |
|---|---|---:|---|---|---|
| **XML rule Komite** | `Komite Claim Non Prop/` | **59** | ekspor Pega, berkas bertanggal 2026-09-09; rule termutakhir di dalamnya `pxUpdateDateTime` 2026-08-13 | 12 rule milik class Komite + 47 rule pinjaman; langkah, pra-syarat, transisi, parameter `Property-Set`, source Java, SQL RDB, definisi Section | **Parameter metode** `RDB-List`, `Connect-REST`, `Obj-Open-By-Handle`, `Obj-Save`, `Apply-DataTransform`, `Property-Set-HTML`, `Page-New/Copy` — nol dari 46 langkah. Harness, Correspondence, Field Value, HTML stream |
| **XML rule Claim** | `Claim Non Prop/` | **279** | ekspor Pega, berkas 2026-09-08 s/d 2026-09-09 | Activity, Flow, FlowAction, Harness, RDBList, ReportDefinition, Section, When, ConnectREST, DataTransform, DecisionTable, SystemSettings | sama seperti di atas; ditambah rule ruleset lain yang dipanggil (`px*`, `pz*`) |
| **DDL Oracle sistem lama** | `_migration-docs/claim-non-prop/pengetahuan/ddl/` | **49** | reverse-engineer Toad dari `DDL_Script_ClaimNonProp.xls`, ditempatkan 2026-09-18 | 36 tabel, 5 view, 6 prosedur/fungsi, termasuk `OS_AKSEPTASI_KLAIM`, `DIRECTTOKASIR_LOG`, `EMAILKOMITE`, `PROC_GENERATE_SEQUENCE_NUMBER`, `PEGA_JSON_OS_AKSEP_KLAIMTNP` | Tiga prosedur akseptasi lain, `TANGGAL_CLOSING`, `HISTORYAKSEPTASIPEGA`, `OS_AKSEPTASI_SUBJECTIVITY` — daftar penuh di §2 |
| **Pohon panggilan flow** | `Komite Claim Non Prop/Struktur_KomiteTreaty_Flow.xlsx` | 1 berkas, 574 sel | buatan tim, 2026-09-09 | hierarki bernomor rule di bawah `KomiteTreaty_Flow`, cocok persis dengan isi folder | **bukan keluaran Pega** — indeks buatan manusia; tidak menyatakan langkah mana memanggil rule mana |
| **ddl-usulan (sistem baru)** | `_migration-docs/claim-non-prop/ddl-usulan/` | **36** | tulisan tim, 2026-09-19 | usulan skema `KLAIMNP` sisi Claim: 22 tabel, 10 view, 4 berkas hak/isian | **Tidak ada satu pun objek Komite** — itulah lubang yang modul ini isi |
| **ADR** | `…/claim-non-prop/docs/adr/` | **29** | 2026-09-18 s/d 2026-09-19 | keputusan sisi Claim, termasuk ADR-0003, ADR-0006, ADR-0016 | **nol ADR untuk modul Komite**; sepuluh ketetapan D/E/G belum tertulis sebagai ADR |
| **FINDING sisi Claim** | `…/claim-non-prop/FINDING-00*.md` | **8** | 2026-09-18/19 | cacat sisi Claim (ambang mata uang, percabangan identitas, dua rumus reinstatement, kurs=1, dst.) | temuan modul Komite (K-01…K-12, G-01…G-20) hidup di berkas lain |
| **SAPUAN** | `…/claim-non-prop/SAPUAN-*.md` | **3** | 2026-09-18/19 | sapuan ulang XML sisi Claim | tidak menyapu folder Komite |
| **Dokumen akar sisi Claim** | `…/claim-non-prop/*.md` | **28** | 2026-09-18 s/d 2026-09-20 | `SPEC-MODEL-DATA.md` (959 baris), `BLUEPRINT.md` (941), `PENGETAHUAN.md` (3.723), `KAMUS-KOLOM.md` (717), `TICKETS.md` (723), `CONTEXT.md` (206), dll. | isi modul Komite hanya muncul sebagai "field pendaratan" (SPEC §9) |
| **Blueprint dua modul** | `MEMORI_PEMAHAMAN.MD` | 1 berkas, 1.331 baris | 2026-09-17 | pemahaman AS-IS Claim **dan** Komite, termasuk §5.3 pembentukan case anak | ditulis sebelum folder Komite dibedah berkas-per-berkas; §11 berstatus parkir |
| **Berkas tahap 1–1.5 modul Komite** | `…/komite-claim-non-prop/` | **2** | `PENGETAHUAN.md` 2026-09-19 · `PUTUSAN-01.md` 2026-09-20 | pengetahuan AS-IS 59 XML; adjudikasi 20 temuan | `GRILL-01.md` **dihapus** atas permintaan (ronde 1 diulang dengan skill). `PENILAIAN-01.md` **tidak pernah ada** |
| **Jawaban tiap ronde** | — | **0 berkas** | ronde 1–3, 2026-09-20 | D-1…D-5, E-1…E-5, G-1…G-4, F-1…F-12 | **Hanya hidup di transkrip percakapan.** Tidak satu pun tersimpan sebagai berkas. Ini lubang inventaris terbesar: sepuluh ketetapan mengikat tanpa berkas |
| **Alat bantu** | `…/claim-non-prop/alat/` + scratchpad | 7 + 1 | 2026-09-19/20 | pembaca XML aktivitas Pega, pembangun ERD, pemeriksa penamaan | pembaca langkah tidak dapat memulihkan parameter metode yang memang tidak diekspor |

---

## 2. Lubang inventaris

Dirujuk oleh bukti yang kita pegang, tetapi berkasnya tidak ada. Hasil sapuan sistematis atas tiga arah: rule yang memanggil rule, SQL yang merujuk objek Oracle, dan DDL yang merujuk DDL lain.

### 2.1 Objek Oracle dirujuk rule Pega, DDL tidak ada (7)

| Objek | Dirujuk oleh | Akibat |
|---|---|---|
| `PEGA_JSON_OS_AKSEP_KLAIM` | `SaveOSClaim_SQL` | arti 9 argumen akseptasi normal tidak terbaca |
| `PEGA_JSON_OS_AKSEP_SUBJECTIVITY` | `SaveOSSubjectivity_SQL` | idem, jalur bersyarat |
| `XOL2_AKSEP_KLAIM` | `SaveXOLClaim_SQL` | sasaran rincian layer tidak terbaca |
| `HISTORYAKSEPTASIPEGA` | `InsertHistoryAkseptasiPega_Sql` | satu-satunya jejak persetujuan di tabel |
| `gl.f_get_email` | `GetEmailCeding_SQL` | sumber email cedant, skema `gl` |
| `reinsurance.trloss_detail_t` | `getStatusKonversi_SQL` | gerbang status konversi |
| `TREATY_OUT` | `GetLimitTONPPLA` (sisi Claim) | di luar cakupan modul ini |

### 2.2 Objek dirujuk DI DALAM DDL yang kita pegang, DDL-nya tidak ada (12)

`OS_AKSEPTASI_SUBJECTIVITY` (dirujuk `VIEW_CLAIMXOL`) · `TANGGAL_CLOSING` (dirujuk `PROC_GENERATE_SEQUENCE_NUMBER` — **penentu periode buku**) · `GCP_IMAGE` (`GET_TOKEN_STORAGE`) · `MST_USER_TEKNIS` (`V_MST_USER_TEKNIS`) · `M_CURRENCY`, `M_CURRENCYSTANDARD`, `M_NATION`, `M_PROVINCE` (view referensi) · `M_SITE_DATABASE`, `D_CAUSE_OF_LOSS`, `M_CAUSE_OF_LOSS` (prosedur sebab kerugian) · `BRANCH`/`CITYINPUT`/`DISTRICTINPUT`/`RWINPUT` (`VIEW_CITY`).
Ditambah satu **database link lintas sistem**: `hrdasm.v_hrd_mst@asmd.sinarmas.co.id`, dirujuk `V_MST_USER_TEKNIS`.

### 2.3 Rule dipanggil, berkasnya tidak ada (3 kustom + 9 bawaan)

Kustom: **`PostEmailKomiteCNP`** (`KomitePostAdjustment` 19.1 — penerima surat jenjang berikutnya) · **`SendEmailWithAttachments`** (×5) · **`SetProtectionEstimation`** (`ProteksiSendKomiteCNP_Act` langkah 21).
Bawaan Pega, tidak perlu dikejar: `pxShowReport`, `pxRetrieveReportData`, `pxAddChildWork`, `pzCheckWorkObjectID`, `pzGetDescendants`, `pzUpdateAndDeleteAssignments`, `pzRemoveAttachments`, `CallVirusCheck`, `UpdateWorkObject`.

### 2.4 Yang TIDAK lagi termasuk lubang — sudah di repo, pernah salah didaftarkan

| Rule | Lokasi sebenarnya | Pernah diminta sebagai |
|---|---|---|
| `CreateChildKomiteCNP_Act` | `Claim Non Prop/Activity/` | pemblokir A-1b |
| `CreateChildKomiteCloseNP_Act` | `Claim Non Prop/Activity/` | — |
| `ProteksiSendKomiteCNP_Act` | `Claim Non Prop/Activity/` | pemblokir A-1b |
| `FilterEmailKomiteWithLimit` | `Claim Non Prop/ReportDefinition/` | pemblokir A-1b |
| `ASMForceCaseClose` | `Claim Non Prop/Activity/` | "ekspor, admin Pega" |
| `HTMLToPDF` | `Claim Non Prop/Activity/` | "ekspor, admin Pega" |
| `GetUrlGoogleStorage_Act` | `Claim Non Prop/Activity/` | "ekspor, admin Pega" |
| `GetBase64Attachment` | dua tempat sekaligus | "rule bawaan, tidak dipindahkan" |
| `KonversiKlaim_Act` | `Komite Claim Non Prop/Activity/` | lubang §2.3 dan pemagar pagar baru, diminta 2026-09-21 — lihat §7.1 |
| `KonversiKlaimNonLife` | `Komite Claim Non Prop/ConnectREST/` | — |
| `InsertJsonClaimTreatyNonProp_act` | `Komite Claim Non Prop/Activity/` | lubang §2.3 dan pemagar pagar baru, diminta 2026-09-21 — lihat §7.1 |
| `InsertClaimPNC` | `Komite Claim Non Prop/RDBList/` | — |

### 2.5 Data produksi yang belum diambil — ditambahkan ronde 5, 2026-09-20

Bagian ini ada karena `AUDIT-TUTUP-01` menemukan cacat bentuk pada inventaris: seluruh
berkas ini menginventarisasi **berkas**, sementara sembilan pagar A-4 dan A-5 menunjuk
ketiadaan **data**. Tanpa bagian ini, sembilan pagar itu tidak punya baris yang dapat
ditunjuk, dan menurut §4 mereka tidak sah.

| # | Data | Objek | Dirujuk oleh | Berapa yang dipegang | Pembukanya |
|---|---|---|---|---|---|
| 1 | Log kiriman ke Kasir | `POOLDATA.DIRECTTOKASIR_LOG` | pagar PG-04 (A-4); G-03, G-04, K-12, Q-10 | **nol baris** — DDL-nya dipegang, isinya tidak | satu `SELECT` baca-saja |
| 2 | Baris akseptasi OS | `OS_AKSEPTASI_KLAIM`, kolom `DATA_JSON` | pagar PG-05 (A-5); Q-5, Q-12, Q-14 | **nol baris** | satu `SELECT` baca-saja |
| 3 | Baris komite pada roster | `EMAILKOMITE` | F-14, N-09 | **nol baris** — DDL-nya dipegang | satu `SELECT` baca-saja |
| 4 | Isi badan prosedur pada basis data **berjalan** | `PROC_GENERATE_SEQUENCE_NUMBER` | pagar PG-03 (A-3); C-01, D-4, K-06, G-05, Q-4, Q-11 | **nol** — yang dipegang hasil reverse-engineer Toad dari sebuah Excel, bukan ekspor langsung | `DBMS_METADATA` atau source dari DBA |
| 5 | Distribusi nilai `.Type` pada adjustment yang pernah masuk komite | `ADJUSTMENT` sisi Claim | N-04, sisa Q-2 | **nol baris** | satu `SELECT` baca-saja |
| 6 | Cacah klaim yang `.KomiteNo`-nya tidak cocok dengan sirkulasi mana pun | `KLAIM` × case Komite | `N-07`, `ADR-0031` | **nol baris** | satu `SELECT` baca-saja · **DBA** |

Empat dari lima baris ini dibuka oleh satu `SELECT` baca-saja. Tidak satu pun memerlukan
Tracer.

Baris 6 ditambahkan ronde 6. Sebab ia ada: bila pembuatan sirkulasi gagal dan langkah
berikutnya tetap menulis `.KomiteNo` dari `pxCoveredInsKeys(<LAST>)` lalu menyimpan, ada klaim
di produksi yang membawa nomor komite **bukan miliknya**. **Migrasi data tidak boleh
mempercayai `.KomiteNo` sebagai tali ke sirkulasi**; tali yang sah adalah relasi dari sisi
sirkulasi. Cacahnya mengubah ukuran pekerjaan migrasi data, bukan bentuk spesifikasi — karena
itu ia tidak memblokir, tetapi ia diambil sebelum migrasi dirancang.

---

## 3. Pola yang sudah lima kali terjadi

> Ditambah satu kejadian keenam pada §5.4.

1. `PENGETAHUAN.md` §0 menyatakan DDL tabel lama tidak ikut diekspor — 49 berkasnya sudah di repo sejak 18 September.
2. `GRILL-01.md` meminta source prosedur sekuens lewat Tracer — berkasnya sudah di repo, dan menjawab pertanyaannya utuh.
3. `GRILL-01.md` meminta ekspor `ASMForceCaseClose` dan `HTMLToPDF` — keduanya sudah di repo.
4. A-1b dibekukan dengan syarat "dibuka oleh ekspor 2 rule" — keempat rule-nya sudah di repo.
5. Daftar pengambilan bukti ronde 3 masih memuat `GetUrlGoogleStorage_Act` sebagai rule absen — sudah di repo.

Penyebabnya satu dan sama: **pernyataan tentang ketiadaan bukti dibuat tanpa memeriksa inventaris**, lalu diwarisi berkas berikutnya sebagai premis.

### Tambahan ronde 6 — pola ketujuh, kedelapan, dan kesembilan

7. **`QF-1` dan `QF-2` sudah dijawab; jawabannya tidak pernah masuk berkas.** `AUDIT-TUTUP-01` dan `GRILL-05/07-AUDIT.md` mengajukan keduanya kembali sebagai perkara terbuka yang menunggu pemilik proses, padahal keputusannya sudah diambil. Ronde 6 mencatatnya sebagai `K6-1` dan `K6-2`.
8. **Status `G-03` dikutip dari berkas yang menyebutnya, bukan dari berkas yang menetapkannya.** Instruksi ronde 6 menyatakan `G-03` berstatus `RAGU`; `PUTUSAN-01.md` baris 57 mencantumkannya pada daftar "Diterima tanpa perubahan", dan yang berstatus `RAGU` adalah bacaan gabungan `G-03`+`G-04` (baris 31).
9. **Keluaran ronde 6 dikira berberkas padahal hanya teks di chat.** Prasyarat `to-spec-komite-v4` menyebut `CONTEXT.md` dan `ADR-0030/0031/0032` sebagai sudah ada; pemeriksaan ke disk menunjukkan tidak satu pun dari sembilan keluaran ronde 6 pernah mendarat.

Arah pola ketujuh **terbalik** dari enam sebelumnya. Enam yang pertama menyatakan bukti tidak
ada padahal ada. Yang ketujuh menyatakan keputusan belum diambil padahal sudah. Pelakunya
bukan interogator melainkan **penilai dan auditor**, yang justru bertugas memeriksa yang
lain. Penyebabnya satu keluarga dengan enam sebelumnya: keadaan sebuah perkara dibaca dari
berkas terakhir yang menyebutnya, bukan dari berkas yang menetapkannya.

Diagnosis pola kesembilan, dicatat apa adanya:

> Delapan kejadian sebelumnya adalah salah membaca **keadaan** — bukti dikira tidak ada
> padahal ada, keputusan dikira belum diambil padahal sudah. Yang kesembilan adalah salah
> membaca **medium**: keluaran dikira berberkas padahal hanya teks di chat. Penyebabnya
> bukan kelalaian melainkan aturan yang tidak lengkap — larangan membuat berkas dipasang
> agar penyalinan dikerjakan pemilik proyek, tetapi tidak ada satu pun ronde yang berakhir
> dengan pemeriksaan bahwa keluarannya sudah mendarat di disk. **Pelakunya penilai.**

Delapan kejadian sebelumnya dicatat tanpa menyamarkan pelakunya; yang kesembilan tidak jadi
pengecualian karena kebetulan pelakunya yang memegang palu.

Aturan kerja yang lahir dari ketiganya ada di `KETETAPAN.md` bagian 7, Tambahan ronde 6.

---

## 4. Aturan pagar

Mulai sekarang, **setiap pagar (aliran dinyatakan beku) dan setiap permintaan pengambilan bukti wajib menunjuk baris inventaris ini** yang membuktikan buktinya memang belum dipegang — nomor bagian dan nama objeknya, bukan kalimat "belum ada".

Pagar atau permintaan yang tidak menunjuk baris inventaris **tidak sah**, dan siapa pun yang membacanya berhak mengabaikannya sampai rujukannya dilengkapi. Bila rujukan itu ternyata menunjuk baris yang isinya ada di repo, pagar itu gugur dengan sendirinya — tanpa perlu ronde, tanpa perlu keputusan.

Inventaris ini dimutakhirkan setiap kali satu berkas bukti masuk atau satu lubang tertutup. Berkas yang tidak dimutakhirkan berhenti menjadi dasar yang sah bagi pagar mana pun.


---

## 5. Pemutakhiran ronde 5 — 2026-09-20

### 5.1 Bukti yang masuk

| Bukti | Cara masuk | Menutup |
|---|---|---|
| `KomiteRouter.xml` dibaca sampai sandi transisinya | dibaca dari ekspor yang sudah dipegang | K-02, Q-2 sebagian, N-01 s/d N-03, N-12 |
| `CreateChildKomiteCNP_Act.xml` 39 langkah, `CreateChildKomiteCloseNP_Act.xml` 14 langkah | dibaca dari ekspor yang sudah dipegang | N-04 s/d N-11, sidang F-13 … F-24 |
| `When/IsKomiteLoop.xml` | dibaca | mekanika gerbang lingkar jenjang |
| `KomiteTreaty_Flow.xml`, bagian `pyTicketShapes` | dibaca | Q-7 |
| `Section/ReinstatementPremiumDetails.xml` | dibaca | Q-3 |
| Transkrip sesi `5696a074-…jsonl` | berkas di disk, dibaca | bunyi D-1…D-5, E-1…E-5, J-1…J-4, H-1…H-7, C-01…C-05 dipulihkan ke `KETETAPAN.md` |

**Lubang inventaris terbesar tertutup.** Baris "Jawaban tiap ronde — 0 berkas" pada §1
berbunyi: sepuluh ketetapan mengikat tanpa berkas, hanya hidup di transkrip percakapan.
Sejak ronde 5, seluruhnya berada di `KETETAPAN.md`. Baris §1 itu **tidak dihapus**; ia
diberi status di sini: **TERTUTUP 2026-09-20 oleh `KETETAPAN.md`.**

### 5.2 Lubang yang tertutup

| Lubang | Status baru | Oleh |
|---|---|---|
| Q-2 — nilai `TransferType` | **TERTUTUP** untuk apa yang dipakai rule; sisa distribusi produksi ada di §2.5 baris 5 | N-04 |
| Q-3 — arti `FlagProrate` | **TERTUTUP sebagian** — empat nilai terbukti dipakai; artinya menunggu ekspor Field Value | N-15 |
| Q-6 — pintu masuk lain `KomitePostAdjustmentCWP` | **TERTUTUP** — penyelidikan dihentikan | N-13, `06-PUTUSAN.md` §2 |
| Q-7 — pemicu `komiteAccept_ticket` | **TERTUTUP sebagian** — sasarannya terbukti; penaiknya di luar ekspor | N-14 |

### 5.3 Lubang baru

| Lubang | Bagian | Keterangan |
|---|---|---|
| `Rule-Obj-FieldValue` untuk `.Type` dan `FlagProrate` | §1 kolom "TIDAK diliput" | sudah tercatat sebagai jenis rule yang tidak diekspor; ronde 5 menunjukkan dua nilai yang artinya benar-benar dibutuhkan |
| Penaik `komiteAccept_ticket` di ruleset lain | baru | tak satu pun dari 338 berkas menaikkannya; pencarian perlu diperluas ke ekspor lain |
| Data produksi | **§2.5 baru** | lima baris, empat di antaranya dibuka satu `SELECT` |

### 5.4 Pola kelima — dan yang keenam

§3 mencatat pola yang sudah lima kali terjadi. Ronde 5 menambah satu kejadian lagi, dengan
bentuk berbeda: **pagar A-1a** yang menahan mekanika keputusan "sampai parameter metode dua
activity inti diekspor" ternyata tidak pernah menyentuh apa yang ditahannya — `KomiteRouter`
seluruhnya `Property-Set`, yang memang terekspor. Pagar itu dicatat gugur sebagai PG-07 di
`REGISTER-PAGAR.md` §2.

Penyebabnya tetap satu dan sama, dengan satu tambahan: pernyataan tentang ketiadaan bukti
dibuat tanpa memeriksa inventaris — **dan** tanpa memeriksa apakah bukti yang hilang itu
memang dibutuhkan oleh perkara yang ditahan.

---

## 6. Pemutakhiran ronde 6 — 2026-09-20

### 6.1 Bukti yang masuk

| Bukti | Cara masuk | Menutup |
|---|---|---|
| `HitServiceToKasirKMT_Act.xml` langkah 10.3 dan 10.7, dibaca dengan penambat langkah | dibaca dari ekspor yang sudah dipegang | P6-1 · `G-03` dikukuhkan |
| `InsertXOLKlaimCNP.xml` langkah 1 dan 4.1 | dibaca | P6-2 · `G-08` dikukuhkan |
| `KomitePostAdjustment.xml` langkah 3, 4, 6, 7 | dibaca | P6-3 · `K-01`; P6-4 · `K-07` |
| `alat/dump_act.py` diperbaiki | perbaikan alat sebelum pembacaan | menutup sebab `M5-01` |

### 6.2 Lubang yang tertutup

| Lubang | Status baru | Oleh |
|---|---|---|
| Bunyi penuh `H-4` | **TERTUTUP** — tercatat di `KETETAPAN.md` bagian 4 | ronde 6 |
| `QF-1`, `QF-2`, `QF-3` tanpa berkas | **TERTUTUP** — `K6-1`, `K6-2`, `K6-3` | ronde 6 |
| Risiko misatribusi transisi pada empat temuan lama | **TERTUTUP** — nol dari empat runtuh | `GRILL-06/01-PEMBACAAN.md` |
| Istilah modul belum dibakukan | **TERTUTUP** — `CONTEXT.md` | ronde 6 |
| Nol ADR untuk modul ini | **TERTUTUP** — `ADR-0030`, `ADR-0031`, `ADR-0032` | ronde 6 |
| Seluruh keluaran ronde 6 belum terbit | **TERTUTUP** — dua belas berkas ditulis ke disk | penerbitan ronde 6, 2026-09-20 |

### 6.3 Lubang baru

| Lubang | Bagian | Keterangan |
|---|---|---|
| Cacah `.KomiteNo` yatim | **§2.5 baris 6** | satu `SELECT`; mengubah ukuran migrasi data |
| Bunyi penuh `H-2` | §10.1 `KETETAPAN.md` | tetap `BUNYI HILANG` — tidak terserap ketetapan atau ADR mana pun |
| Bunyi penuh `H-1`, `H-3`, `H-5`, `H-6`, `H-7` | §10.1 `KETETAPAN.md` | berstatus `RINGKAS — MENGIKAT APA ADANYA`, masing-masing menunjuk nomor yang menyerapnya |

---

## 7. Pemutakhiran tahap spesifikasi — 2026-09-21

### 7.1 Dua rule yang hampir didaftarkan sebagai lubang — dan tidak

Instruksi `tutup-spec-komite` memerintahkan menambahkan `KonversiKlaim_Act` dan
`InsertJsonClaimTreatyNonProp_act` ke **§2.3**, lalu memasang pagar baru di atasnya.

**Perintah itu tidak dijalankan, karena premisnya tidak bertahan.** Keempat berkasnya ada di
repo, diperiksa sebelum ditulis:

| Rule | Lokasi sebenarnya | Byte | Pernah diminta sebagai |
|---|---|---:|---|
| `KonversiKlaim_Act` | `Komite Claim Non Prop/Activity/` | 50.135 | lubang §2.3, pemagar pagar baru |
| `KonversiKlaimNonLife` | `Komite Claim Non Prop/ConnectREST/` | 14.056 | — |
| `InsertJsonClaimTreatyNonProp_act` | `Komite Claim Non Prop/Activity/` | 54.465 | lubang §2.3, pemagar pagar baru |
| `InsertClaimPNC` | `Komite Claim Non Prop/RDBList/` | 6.919 | — |

Keempatnya karena itu masuk **§2.4**, bukan §2.3. Pagar yang hendak dipasang di atasnya akan
menunjuk sesuatu yang ada di repo, dan menurut §4 ia **gugur dengan sendirinya** — tanpa
ronde, tanpa keputusan. Memasangnya lebih dulu lalu menunggu ia gugur adalah pekerjaan yang
hasilnya sudah diketahui.

**Ini kejadian kesepuluh dari pola §3, dan ketujuh dari keluarga yang pertama** — pernyataan
tentang ketiadaan bukti dibuat tanpa memeriksa inventaris. Ia dicatat dengan cara yang sama
seperti sembilan sebelumnya: tanpa menyamarkan pelakunya. Kali ini pelakunya **instruksi
tahap spesifikasi**, dan yang menangkapnya adalah aturan §4 yang dipasang justru untuk ini.

Yang **benar** dari perintah itu tetap berlaku: kedua rule memang dipanggil dari jalur
keputusan sistem lama dan memang **tidak terpeta ke satu pun dari lima port hilir** yang
`SPEC-KOMITE-01.md` namai. Itu lubang **rancangan**, bukan lubang bukti — dan perbedaan itu
yang menentukan siapa yang menutupnya.

Pemanggilnya, sejauh terbaca: keduanya berada di bawah `KomitePostAdjustment` pada pohon
panggilan `PENGETAHUAN.md` §2.3. **Nomor langkahnya tidak terbaca** — pohon itu bernomor
hierarkis penelusuran, dan §1 berkas ini sudah mencatat bahwa ia "tidak menyatakan langkah
mana memanggil rule mana". Menyebut nomor langkah dari nomor pohon akan menaikkan indeks
buatan manusia menjadi bukti langkah.

### 7.2 Yang ditambahkan ke §2.4

Empat baris di atas ditambahkan ke tabel §2.4 sebagai rule yang **tidak** termasuk lubang.

### 7.3 Lubang yang tetap terbuka

Tidak ada lubang baru. §2.1 sampai §2.5 tidak berubah isinya ronde ini, dan **nol pagar baru
dipasang**. Enam baris §2.5 tetap nol baris data.

### 7.4 Satu pembacaan yang belum dilakukan — dan mengapa ia bukan lubang

`Section/ShowTransfer.xml` ada di repo. Dua dari enam medan yang `F-2` cacah sebagai
dapat-disunting belum dinamai di berkas dokumentasi mana pun, dan namanya terbaca dari sana.
Itu **pembacaan yang belum dilakukan**, bukan bukti yang belum dipegang — karena itu ia tidak
masuk §2 dan tidak boleh dipagari. Dicatat di sini agar tidak kelak didaftarkan sebagai
lubang oleh pembaca berikutnya.
