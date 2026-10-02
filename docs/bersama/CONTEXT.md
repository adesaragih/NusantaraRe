# Claim — Life

Siklus penanganan klaim untuk lini **life** pada bisnis reasuransi Nusantara Re: dari pendaftaran
klaim hingga keputusan akseptasi. Konteks ini **tidak** mencakup tangga persetujuan komite — itu
konteks luar (lihat **Komite Life**).

Dibuat 2026-09-14 dari grilling Ronde 1 (`.scratch/claim-life/grilling-ronde-1.md`), jawaban
**work owner**, dan bukti FASE A di `discovery/`.

> Setiap entri membawa baris `_Bukti_` berisi `path + rule` — perluasan wajib dari aturan proyek
> (`CLAUDE.md` §4). Istilah yang artinya belum dijawab work owner **tidak dimasukkan** ke sini;
> ia tetap di `discovery/open-questions.md`.

---

## Tahapan siklus

**Register**:
Tahap pendaftaran klaim life yang baru masuk. Tahap pertama siklus.
_Avoid_: Registrasi, Input Klaim, Lodgement
_Bukti_: `Claim Life/Flow/Register_Flow.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` / `RULE-OBJ-FLOW`), shape `Assignment2` "Input Register"

**Outstanding**:
Tahap klaim yang sudah terdaftar dan masih berjalan — belum diaksep maupun ditolak.
_Avoid_: Pending, Open Claim, Dalam Proses
_Bukti_: shape `Assignment1` "Outstanding Claim" pada `Register_Flow`; nilai `STS_REJECT = 0`

**Medical Check**:
Tahap penelaahan medis atas klaim life.
_Avoid_: Pemeriksaan Medis, Medical Review, Underwriting Medis
_Bukti_: shape `Assignment3` "Medical Check" pada `Register_Flow`; FlowAction `MEDICALCHECK`

**Claim Analis**:
Tahap keputusan akseptasi klaim. Tahap terakhir sebelum siklus selesai.
_Avoid_: Claim Analyst, Analisa Klaim, Assessment
_Bukti_: shape `Assignment4` "Claim Analis" pada `Register_Flow`; FlowAction `AKSEPTASICLAIMLIFE`

**Akseptasi**:
Keputusan menerima klaim. Dipakai konsisten di seluruh korpus untuk "penerimaan setelah
penelaahan", bukan sekadar persetujuan administratif.
_Avoid_: Approval, Persetujuan, Acceptance
_Bukti_: FlowAction `AKSEPTASICLAIMLIFE`; tabel `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`

---

## Peran

`[terverifikasi work owner 2026-09-14]` Ketiganya **daftar lengkap** peran di siklus klaim Life.
Rangkap peran **tidak diperbolehkan**, kecuali akses ditambahkan eksplisit pada role akun.

**ReasLifeAdmin**:
Peran yang mengerjakan tahap **Register** dan **Outstanding**. Menginput baris **AdjustmentList**
pertama dan menyimpannya ke Outstanding (`STS_REJECT = 0`). **Dapat menolak adjustment yang ia input
sendiri, tanpa Komite** — penolakan itu **membatalkan baris itu saja**, klaim tidak tertutup, dan ia
lalu menginput baris baru. Tidak mengirim ke Komite — **kecuali** untuk klaim ber-`Type` `TP`/`TR`,
yang bebas peran (**ADR-0012**). Tujuan pengembalian kasus dari dua peran lain.
_Avoid_: Admin, Administrator Klaim
_Bukti_: 21 kemunculan `'ReasLifeAdmin'` di 6 berkas `Claim Life`; gerbang tombol "Reject Outstanding" di `Claim Life/Section/AdjustmentDetail_Section.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION`) = `pyWorkPage.pyPosition =='ReasLifeAdmin' && …CLAIM_NO !='' && .STS_REJECT == 0`. Peran per tahap dikonfirmasi work owner 2026-09-14

**ReasLifeMedicalAdvisor**:
Peran yang mengerjakan tahap **Medical Check**. Tidak mengubah `STS_REJECT`.
_Avoid_: Dokter, Medical Officer, Penasihat Medis
_Bukti_: 13 kemunculan `'ReasLifeMedicalAdvisor'` di 5 berkas `Claim Life`

**ReasLifeSPV**:
Peran yang mengerjakan tahap **Claim Analis**. **Satu-satunya peran yang mengirim ke Komite —
untuk `Type` selain `TP`/`TR`.** Setelah Komite menolak, ia menambah baris **AdjustmentList** baru
dan menyimpannya ke Outstanding.
_Avoid_: Supervisor, SPV, Atasan
_Bukti_: 12 kemunculan `'ReasLifeSPV'` di 3 berkas `Claim Life`; gerbang jalur Komite (`GetListKomiteLife`) dan tombol "Save to Outstanding" di `AdjustmentDetail_Section.xml`. Kewenangan dikonfirmasi work owner 2026-09-14

⚠️ **Pengecualian `Type` `TP`/`TR`** `[keputusan work owner 2026-09-14]`: untuk klaim ber-`Type`
**`TP`** atau **`TR`**, pengiriman ke Komite **tidak dibatasi peran** — **ReasLifeAdmin pun dapat
mengirim langsung**. Untuk tipe lain, hanya SPV. **Perilaku ini dibawa apa adanya (paritas)**;
risikonya dicatat di **ADR-0012**.
_Bukti_ `[terverifikasi]`: gerbang jalur Komite =
`pyWorkPage.pyPosition =='ReasLifeSPV' || pyWorkPage.Type = 'TP' || pyWorkPage.Type = 'TR'` —
seluruhnya `||`, tanpa `&&`, sehingga tidak ada ambiguitas presedensi. `=` di korpus ini adalah
**pembanding**, bukan assignment: pada `<pyCondition>` bentuk `=` justru mayoritas (2.282 vs 954
`==`), dan banyak ekspresi mencampur keduanya dalam satu baris — mis.
`.ACCEPTEDNO=="" && .IsCheck = true && .STS_REJECT == "0"`

⚠️ `[terverifikasi]` **Penegakan peran di sistem lama tidak seragam.** `InputOSClaimLife.xml`,
`RejectOSClaimLife_Sec.xml`, `ClaimComite.xml`, dan `Harness/Committe_Life.xml` **tidak** memuat
`pyPosition` sama sekali; penegakan bertumpu pada penugasan tahap di `Register_Flow`. Sistem baru
menegakkan peran di lapisan layanan (**ADR-0002**), bukan meniru ketidakseragaman ini.

---

## Status dan kode

**STS_REJECT**:
Status **satu baris `AdjustmentList`**: `0` = **Outstanding**, `1` = **Aksep** (diterima),
`2` = **Reject** (ditolak).
⚠️ **Nama field ini menyesatkan** — nilai `1` berarti *diaksep*, bukan ditolak. Jangan membaca
namanya sebagai artinya.
_Avoid_: Reject Status, Status Penolakan, Status Klaim
_Bukti_: 50 kemunculan di `Claim Life`; pola `STS_REJECT=='1' || STS_REJECT=='2'` (7×) = "sudah selesai diproses". Arti nilai dikonfirmasi work owner 2026-09-14

**Mesin status** `[keputusan work owner 2026-09-14]` — selengkapnya di **ADR-0011**:

| # | Langkah | Peran | Akibat pada baris |
| --- | --- | --- | --- |
| 1 | Input baris `AdjustmentList` pertama → **Save ke OS** | `ReasLifeAdmin` | `0` |
| 1b | **Reject Outstanding** atas adjustment yang ia input — tanpa Komite. **Membatalkan baris itu saja**; Admin lalu input baris baru. | `ReasLifeAdmin` | `2` |
| 2 | Submit ke Medical Check | → `ReasLifeMedicalAdvisor` | tetap `0` |
| 3 | Submit ke SPV | → `ReasLifeSPV` | tetap `0` |
| 4 | **Send ke Komite** | `ReasLifeSPV` — bebas peran bila `Type` `TP`/`TR` (**ADR-0012**) | tetap `0` |
| 5 | Komite memutus | Komite Life | aksep → `1`; tolak → `2` |
| 6 | Setelah tolak: tambah baris **baru** → Save ke OS → send Komite. Berulang. | `ReasLifeSPV` | baris baru `0` |

**Dua sumber penolakan**, bukan satu: Admin menolak langsung tanpa Komite, dan Komite menolak lewat
SPV. **Aksep final selalu lewat Komite.**

**Arti tunggal nilai `2`** `[keputusan work owner 2026-09-14]`: `STS_REJECT = 2` **selalu** berarti
"baris ini ditolak", **tidak pernah** "klaim selesai". Kedua sumber penolakan menulis nilai yang
sama dengan makna operasional **setara pada tingkat baris**; keduanya diikuti baris baru, bukan
penutupan klaim.

**Kefinalan**: baris **terminal** — sekali `1`/`2`, tidak berubah. **Klaim tidak terminal** —
setelah tolak dibuat baris baru (oleh SPV bila penolaknya Komite, oleh Admin bila ia sendiri yang
menolak). **Revisi = baris baru**, bukan pengubahan baris lama.

_Bukti_ `[terverifikasi]`: `Komite Claim Life/Activity/KomitePostAdjustment.xml` (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`) — 6 `Property-Set` bernilai `1`, 2 bernilai `2`, digerbangi `.ACCEPTEDNO=="" && .IsCheck = true && .STS_REJECT == "0"`; `Claim Life/Activity/RejectOSClaimLife_Act.xml` menulis `2` dan dipicu tombol "Reject Outstanding" bergerbang `ReasLifeAdmin`; tidak ada rule yang menulis `0` setelah `1`/`2`. Korpus **tidak membedakan** kedua sumber penolakan — sesuai keputusan di atas, karena memang tidak seharusnya berbeda

**PremiumListDetail**:
Baris daftar premi pada sebuah klaim. `STS_REJECT` di tingkat ini adalah **cerminan** hasil baris
`AdjustmentList` terakhir, **bukan** unit keputusan tersendiri. Header klaim mengikuti adjustment
terakhir.
_Avoid_: Premium List, Detail Premi, Baris Premi
_Bukti_ `[terverifikasi]`: setiap rule yang memutus menulis **dua tingkat berbarengan** dengan nilai sama — `RejectOSClaimLife_Act` (`.STS_REJECT = 2` dan `…PremiumListDetail(idx).STS_REJECT = 2`), `KomitePostAdjustment` (pasangan `1`/`1` dan `2`/`2`). `[keputusan work owner]`: bahwa yang tercermin adalah baris **terakhir** — korpus hanya menunjukkan baris yang sedang diputus

`[terverifikasi]` **`STS_REJECT` hidup di tiga tingkat**: klaim (`pyWorkPage.ClaimData.STS_REJECT`),
`PremiumListDetail(idx)`, dan `AdjustmentList(idx)`. Penyalinan terjadi **dua arah** —
`Claim Life/Activity/SetSTS_Reject.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!SETSTS_REJECT`)
menurunkan `.STS_REJECT = Primary.STS_REJECT`, sementara
`Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` menaikkan
`pyWorkPage.ClaimData.STS_REJECT = .STS_REJECT`. Otoritasnya ditetapkan **ADR-0011**: baris
adjustment (OQ-061 ditutup untuk konteks ini).

⚠️ `[keputusan work owner 2026-09-14]` **`Claim Life/Activity/SaveAdjustment_Act.xml` adalah dead
rule — jangan dimigrasikan.** Penulis `STS_REJECT = 1` yang berlaku adalah rule sisi Komite,
`KomitePostAdjustment`. Status "dead" ini **tidak dapat diverifikasi dari korpus**: yang terbaca
justru sebaliknya — `[terverifikasi]` rule itu **terpasang di UI** (dirujuk 2× dari
`Claim Life/Section/ClaimLifeDetailGCNM.xml`) dan menulis `.AdjustmentList(<LAST>).STS_REJECT = 1`
tanpa precondition Komite. Lihat **ADR-0011** §`SaveAdjustment_Act`.

**AdjustmentList**:
Daftar baris penyesuaian pada sebuah klaim. **Unit keputusan status** — yang benar-benar diaksep
atau ditolak adalah barisnya, bukan klaimnya. Baris pertama diinput **ReasLifeAdmin**. Baris
berikutnya ditambah **ReasLifeSPV** setelah Komite menolak, atau **ReasLifeAdmin** setelah ia
sendiri menolak barisnya.
_Avoid_: Adjustment, Penyesuaian, Daftar Adjustment
_Bukti_: class `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` (dipakai 8 rule `Claim Life`); `Claim Life/Section/AdjustmentDetail_Section.xml` menggerbangi wewenang terhadap `.STS_REJECT` **tingkat baris**; `Claim Life/Activity/SetIndexAdjustmentList.xml` (`…!SETINDEXADJUSTMENTLIST`) menyalin 8 kolom dari `AdjustmentList(1)` ke `AdjustmentList(<LAST>)` — `SHARE_NUSANTARA_RE`, `CEDING_RETENTION`, `SUM_REASURED`, `SUM_INSURED`, `SHARE_RETRO`, `RETROCEDED_SHARE`, `CURRENCYID`, `CURRENCY` — **tanpa** `STS_REJECT`; 50 kemunculan `.AdjustmentList`. Peran sebagai unit keputusan dikonfirmasi work owner 2026-09-14

⚠️ `[terverifikasi 2026-09-14]` **`AdjustmentList` tidak punya tabel fisik sendiri.** Kelasnya
`ASM-FW-GISFW-Data-AdjustmentLife` berawalan **`Data-`** — embedded, tanpa tabel — dengan induk
`ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL`. Ia **page-list in-memory** sebelum disimpan; hasilnya
di-persist sebagai **baris pada `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`**, yang kolomnya identik dengan
kolom baris adjustment.
_Bukti_: `Claim Life/Activity/SaveOutStandingLife_Act.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`) — langkah `Property-Set` dengan `pyStepsObjectName = .AdjustmentList` dan `pyStepsClassName = ASM-FW-GISFW-Data-AdjustmentLife`
Ini **tidak membatalkan ADR-0011**: unit *keputusan* tetap baris adjustment. Yang kini diketahui
adalah **di mana baris itu mendarat**.

**AcceptStatus**:
Hasil keputusan **Komite Life**: `1` = diaksep, `2` = reject. **Kosakata konteks luar** — bukan
status internal Claim — Life. Ia **dipetakan ke `STS_REJECT` di batas kontrak** dan tidak disimpan
sebagai status kedua (**ADR-0001**).
_Avoid_: Approval Status, Status Persetujuan
_Bukti_: `Komite Claim Life` (28 kemunculan) — precondition `pyWorkPage.AcceptStatus = 1 && pyWorkPage.KomiteCount == pyWorkPage.KomiteLoop` dan padanannya `==2`; di `Claim Life` hanya 1 kemunculan (`Activity/SendEmailKlaimLF.xml`). Arti dan pemetaan dikonfirmasi work owner 2026-09-14

**KomiteLoop**:
Jumlah tingkat tangga persetujuan Komite untuk satu kasus. **Ditentukan Claim — Life**, dikirim ke
Komite sebagai bagian muatan penyerahan. `STS_REJECT` hanya berubah pada tingkat terakhir
(`KomiteCount == KomiteLoop`).
`[terverifikasi 2026-09-14]` **Nilainya dihitung, bukan dikonfigurasi:**
**KomiteLoop = COUNT baris roster `EMAILKOMITE` yang aktif dan ber-`LIMIT_BOTTOM <= CLAIM_AMOUNT`.**
Tidak ada konstanta di mana pun. (OQ-032 **tertutup**.)
_Avoid_: Komite Level, Tingkat Komite, Loop
_Bukti_: `Claim Life/Activity/GetListKomiteLife.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `GETLISTKOMITELIFE` / `RULE-OBJ-ACTIVITY`) set `Local.IsADj = .CLAIM_AMOUNT` → panggil report `FilterEmailKomiteWithLimit` → hasil ke `.KomiteList`; `Claim Life/ReportDefinition/FilterEmailKomiteWithLimit.xml` (`ASM-FW-GCNMFW-INT-EMAILKOMITE` / `FILTEREMAILKOMITEWITHLIMIT` / `RULE-OBJ-REPORT-DEFINITION`) filter `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM && .STS_KLAIM = Param.STS_KLAIM && .STS_AKTIF = "1"`; `Claim Life/Activity/CreateKMTLife_Act.xml` set `childPageKomite.KomiteLoop = @Utilities.SizeOfPropertyList(childPageKomite.KomiteList)`

⚠️ `[terverifikasi]` Ambang yang dikirim ke roster adalah **nilai mutlak** klaim —
`Param.LIMIT_BOTTOM = @if(Local.IsADj<0, Local.IsADj * -1, Local.IsADj)`. Nilai klaim negatif
dihilangkan tandanya sebelum pencarian. Perilaku ini harus direplikasi apa adanya.

**EMAILKOMITE**:
Tabel roster anggota komite beserta **pita nilai**-nya. Menentukan siapa yang masuk tangga
persetujuan sebuah klaim, dan — lewat jumlah barisnya — **berapa tingkat** tangga itu.
_Avoid_: Daftar Komite, Committee Roster, Email Komite
_Bukti_ `[data DBA 2026-09-14]`: `POOLDATA.EMAILKOMITE` — PK `ID INTEGER`; `LIMIT_BOTTOM INTEGER`, `LIMIT_TOP INTEGER`, `DEGREE VARCHAR2(150)`, `TYPE_KOMITE INTEGER`, `STS_AKTIF VARCHAR2(150)`, `STS_KLAIM VARCHAR2(10)`, `STS_REJECT VARCHAR2(15)`
⚠️ Roster **tidak punya kolom mata uang** → pita berlaku atas **satu mata uang implisit**.
⚠️ `STS_REJECT` di sini bertipe **`VARCHAR2(15)`**, berbeda dari `STS_REJECT NUMBER(38)` di tabel
klaim. Nama sama, tabel berbeda, tipe berbeda — jangan disatukan.

**M_LINK_SERVICE**:
Tabel yang memetakan **kunci kategori** `(KATEGORI_1, KATEGORI_2)` ke **alamat endpoint** keluar.
Alamat layanan adalah **data**, bukan konfigurasi kode — pemisahan dev/prod ditentukan isi tabel
per-database.
_Avoid_: Service Registry, Daftar Endpoint, Konfigurasi URL
_Bukti_: `Claim Life/Activity/GetLinkService.xml` (`ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY`) — `Obj-Browse` `.KATEGORI_1 = Param.Kategori_1 && .KATEGORI_2 = Param.Kategori_2`, ambil `.URL`, lalu `Connect-REST`. Kunci Claim — Life: `Kategori_1 = "Klaim"`, `Kategori_2 = "insertClaimLife"` (`Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml`). Lihat **ADR-0013**

**Type**:
Tipe polis yang mendasari klaim. Nilainya `QP`, `QR`, `TP`, `TR`. **Bukan sekadar data** — ia
menggerbangi **wewenang** (siapa boleh kirim ke Komite, **ADR-0012**) dan **validasi** (jendela
Date of Loss). Sumber otoritatifnya `PolicyDataLife.Type`; `pyWorkPage.Type` adalah salinan kerja.
_Avoid_: Jenis, Tipe Klaim, Policy Type
_Bukti_ `[terverifikasi]`: 24 kemunculan `pyWorkPage.Type` dan 22 `pyWorkPage.PolicyDataLife.Type` di `Claim Life`; penyalinannya di `Claim Life/Activity/LoadDataPeserta_Act.xml` — `pyWorkPage.Type = pyWorkPage.PolicyDataLife.Type`

**TP** (Payable), **TR** (Receivable):
Dua nilai `Type`. ⚠️ **Kepanjangan ini tidak ada di korpus** — sumbernya work owner.
Keduanya membuat pengiriman ke Komite **bebas peran** (**ADR-0012**) dan mengarahkan validasi DOL
ke jendela **retrosesi**.
_Avoid_: Treaty Proportional, Treaty Retro (kepanjangan yang pernah diduga — **bukan** artinya)
_Bukti_ `[keputusan work owner 2026-09-14]` untuk kepanjangan. `[terverifikasi]` untuk perilakunya: gerbang Komite di `Claim Life/Section/AdjustmentDetail_Section.xml`; cabang `Type=="TR"||Type=="TP"` → `RETROCESSION_VALUATION_BEGIN_DATE`/`_EXPIRED_DATE` di `Claim Life/Activity/ValidasiDOL_Act.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `VALIDASIDOL_ACT`)

**QP**, **QR**:
Dua nilai `Type` lainnya. **Artinya belum dijawab** = **OQ-020** — tidak dimasukkan sebagai istilah
sampai dikonfirmasi.
⚠️ **KOREKSI 2026-09-16:** sempat ditandai "QP=Payable, QR=Receivable [terverifikasi] dari dropdown
EDM Life" — itu **KELIRU** dan sudah dibatalkan. Sensus korpus: `Receivable`/`Payable` **NOL berkas**
di `Endorsement Life/`; label itu hanya di `pyLocalizedValue` payload `DATA_JSON` runtime (instance
data), **bukan rule korpus**. Arti QP/QR **tetap OQ-020 terbuka**.
_Bukti_ `[terverifikasi]` untuk perilakunya saja: cabang `Type=="QR"||Type=="QP"` di `ValidasiDOL_Act` memakai jendela **gross** (`GROSS_VALUATION_BEGIN_DATE`/`_EXPIRED_DATE`) dengan pergeseran tanggal **nol**, sedangkan cabang `TP`/`TR` memakai jendela retrosesi dengan pergeseran **+1 hari** (`@addCalendar(.DATE_OF_LOSS,0,0,0,1,0,0,0)`)

**SendtoAdmin**:
Penanda pengembalian kasus ke **ReasLifeAdmin**, dari **ReasLifeMedicalAdvisor** atau
**ReasLifeSPV**. Bernilai `1` saat aktif.
_Avoid_: Return to Admin, Kembalikan
_Bukti_: `Claim Life/When/IsSendtoAdmin.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOADMIN` / `RULE-OBJ-WHEN`) — `pyWorkPage.SendtoAdmin = 1`

**SendtoMedical**:
Penanda pengembalian kasus dari **ReasLifeSPV** ke **ReasLifeMedicalAdvisor**. Bernilai `1` saat
aktif.
_Avoid_: Return to Medical, Kembalikan ke Medis
_Bukti_: `Claim Life/When/IsSendtoMedical.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOMEDICAL` / `RULE-OBJ-WHEN`) — kondisi tidak terbaca dari tag; arti dikonfirmasi work owner 2026-09-14

**ContentNote**:
**Jenis klaim**, diturunkan dari `BusinessCode`. Nilai: `DEATH`, `HEALTH`, `CI`, `TPD`, `TI`.
_Avoid_: Catatan, Note, Keterangan
_Bukti_: `Claim Life/Activity/SaveOutStandingLife_Act.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT`) — gerbang `RDB-List` + `Property-Set` pada nilai `"DEATH"`; turunan dari `BusinessCode` dikonfirmasi work owner 2026-09-14

**BusinessCode**:
Kode produk life. Menentukan **ContentNote** (jenis klaim). Daftar lengkap di bawah.
_Avoid_: Product Code, Kode Lini, LOB Code
_Bukti_: `Claim Life/Activity/SaveOutStandingLife_Act.xml` — satu precondition menguji `L1`…`L11` berderet; daftar penuh `L1`–`L21` dari work owner 2026-09-14

### Daftar BusinessCode → produk → jenis klaim

`[terverifikasi work owner 2026-09-14]`

| Kode | Nama produk | ContentNote |
| --- | --- | --- |
| `L1` | INDIVIDUAL TERM LIFE | `DEATH` |
| `L2` | GROUP TERM LIFE | `DEATH` |
| `L3` | INDIVIDUAL WHOLE LIFE | `DEATH` |
| `L4` | GROUP PA | `DEATH` |
| `L5` | GROUP LEVEL TERM LIFE | `DEATH` |
| `L6` | GROUP DECREASING TERM LIFE | `DEATH` |
| `L7` | INDIVIDU ENDOWMENT LIFE | `DEATH` |
| `L8` | INDIVIDU INCREASING TERM LIFE | `DEATH` |
| `L9` | GROUP ENDOWMENT LIFE | `DEATH` |
| `L10` | GROUP INCREASING TERM LIFE | `DEATH` |
| `L11` | INDIVIDU PA | `DEATH` |
| `L12` | INDIVIDU EXPENSE HEALTH | `HEALTH` |
| `L13` | INDIVIDU DISABILITY HEALTH | `HEALTH` |
| `L14` | GROUP EXPENSE HEALTH | `HEALTH` |
| `L15` | GROUP DISABILITY HEALTH | `HEALTH` |
| `L16` | INDIVIDU CRITICAL ILLNESS | `CI` |
| `L17` | INDIVIDU TPD | `TPD` |
| `L18` | INDIVIDU HOSPITAL CASH PLAN | `HEALTH` |
| `L19` | GROUP CRITICAL ILLNESS | `CI` |
| `L20` | GROUP TPD | `TPD` |
| `L21` | GROUP TERMINAL ILLNESS | `TI` |

`[terverifikasi]` Modul `Claim Life` **hanya menguji `L1`–`L11`** (seluruhnya `DEATH`) dalam satu
precondition. `L12`–`L21` ada di master produk, bukan di modul klaim ini.

**CI**, **TPD**, **TI**:
Nilai `ContentNote` untuk *critical illness*, *total & permanent disability*, dan *terminal
illness*. **Kepanjangan ditulis apa adanya dari work owner**; tidak dijabarkan di korpus.
_Bukti_: daftar `BusinessCode` dari work owner 2026-09-14

---

## Batas konteks

**Komite Life**:
Tangga persetujuan komite untuk klaim life. **Konteks/sistem luar** bagi Claim — Life. Claim — Life
menyerahkan kasus kepadanya sebagai *child work*; isi dan alur internalnya **tidak
dispesifikasikan** di konteks ini.
_Avoid_: Komite, Committee, Approval Board
_Bukti_: class `ASM-FW-GCNMFW-Work-KomiteLife`; `Claim Life/Activity/CreateKMTLife_Act.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`) merujuknya 15×

**Penyerahan ke Komite**:
Tindakan membuat *child work* berkelas **Komite Life** dari layar klaim. Dipicu dari UI, bukan dari
graf alur. Berlaku untuk baris `AdjustmentList` yang masih `STS_REJECT = 0`.
Pelakunya bergantung pada `Type` klaim `[keputusan work owner 2026-09-14]`:
**`Type` `TP` atau `TR`** → **bebas peran**, `ReasLifeAdmin` pun dapat mengirim langsung.
**`Type` lain** → **hanya `ReasLifeSPV`**.
_Avoid_: Eskalasi, Escalation, Submit to Committee
_Bukti_: `Claim Life/Activity/CreateKMTLife_Act.xml` — langkah `Call pxAddChildWork` → `Obj-Refresh-And-Lock` → `Obj-Save`; dipicu lewat `<pyActivity>CreateKMTLife_Act</pyActivity>` di `Claim Life/Section/ClaimComite.xml` (2×) dan `Claim Life/Harness/Committe_Life.xml` (2×); gerbang jalur Komite (`GetListKomiteLife`) di `Claim Life/Section/AdjustmentDetail_Section.xml` = `pyWorkPage.pyPosition =='ReasLifeSPV' || pyWorkPage.Type = 'TP' || pyWorkPage.Type = 'TR'`

`[pertanyaan terbuka]` **Kapan** penyerahan wajib dipicu — apakah otomatis di atas nilai tertentu
atau keputusan manual — masih terbuka di **OQ-039** (butir yang tersisa untuk 3 modul Claim lain;
untuk Claim — Life akibat rejectnya sudah dijawab).

**Retro**:
Retrosesi — penempatan ulang sebagian risiko yang sudah diterima. Muncul sebagai percabangan
penomoran (retro vs non-retro) di beberapa konteks.
_Avoid_: Retrocession, Retrosesi
_Bukti_: `Komite Claim Life` — `Generate_NoAccept_KMT_Life` vs `Generate_NoAccept_KMT_LifeRetro`

---

## Catatan pemeliharaan

- Berkas ini **glossary saja**. Keputusan arsitektur → `docs/adr/`. Spesifikasi → `.scratch/claim-life/spec.md`.
- Seed awalnya `discovery/glossary.md` (169 entri korpus-wide). Berkas ini memuat **hanya istilah
  konteks Claim — Life yang sudah dikonfirmasi work owner**.
- Istilah yang artinya belum dijawab **tidak dimasukkan** — ia tetap terbuka di
  `discovery/open-questions.md`.

---

# Lampiran — OQ pemblokir Claim Life: semua tertutup (2026-09-14)

> Bagian ini **bukan glossary**; ia catatan status. Kosakata ada di bagian-bagian di atas.
> Rincian penuh tiap OQ ada di `discovery/open-questions.md`.

| OQ | Verdict untuk Claim — Life | Sumber |
| --- | --- | --- |
| **OQ-001** DDL Oracle | **TERTUTUP** — DDL **12 tabel/view** mencakup seluruh persistensi Claim Life; uang = `NUMBER` tanpa presisi; `STS_REJECT NUMBER(38)`; tiga master adalah **view atas `JSONDATA`**; `AdjustmentList` **tanpa tabel fisik** | `[data DBA]` + `[terverifikasi]` |
| **OQ-002** format nomor & penomoran | **TERTUTUP** — `RNML-KL1.08.2026.00936`; prefix di-**lookup** dari `KODE_PRODUKSI`; sequence per `(class, jenis, tahun)` dengan `SELECT … FOR UPDATE` | `[data DBA]` + `[terverifikasi]` |
| **OQ-013** batas transaksi | **TERTUTUP** — proc jalur Life **tidak** commit sendiri; **Go memegang batas transaksi** | `[data DBA]` + `[terverifikasi]` + `[keputusan work owner]` |
| **OQ-018** host/URL ter-hardcode | **TERTUTUP** — **nol** URL bisnis di modul ini; lingkungan dibedakan `pzProductionLevel == "5"` | `[terverifikasi]` |
| **OQ-032** tingkat & anggota komite | **TERTUTUP** — `KomiteLoop` = **COUNT roster aktif** ber-`LIMIT_BOTTOM <= CLAIM_AMOUNT`; tidak ada konstanta | `[terverifikasi]` |
| **OQ-037** ambang roster | **TERTUTUP** — ambang adalah **data di tabel `EMAILKOMITE`**, bukan hardcode | `[terverifikasi]` + `[data DBA]` |
| **OQ-047** daftar endpoint | **TERTUTUP** — resolusi **runtime lookup** `M_LINK_SERVICE` lewat `(KATEGORI_1, KATEGORI_2)`; kunci Life = `Klaim` / `insertClaimLife` | `[terverifikasi]` + `[data DBA]` |
| **OQ-060** cakupan `CURRENCY` | **TERTUTUP** — `(amount, currency)` per baris, **invariant: satu klaim satu mata uang** | `[terverifikasi]` + `[keputusan work owner]` |

**Masih terbuka untuk modul lain:** OQ-001 (DDL di luar Claim Life), OQ-002 (66 procedure lain),
OQ-013 (proc yang commit sendiri, mis. `PEGA_JSON_POLIS_TREATYIN`), OQ-018 (URL literal di 15 modul),
OQ-037 (pita hardcode di Komite Claim FacIn), OQ-060 (cakupan mata uang konteks lain).
**OQ-032** dan **OQ-047** tertutup **tanpa sisa cakupan**.

## Tiga keputusan desain non-korpus `[keputusan work owner]`

1. **Aturan "klaim selesai"** (tiket 04). Status klaim adalah **turunan**, bukan kolom tersimpan.
   Unit keputusan = **baris `AdjustmentList`** (**ADR-0011**); **header klaim mengikuti adjustment
   terakhir**. Sistem lama memang tidak menyimpan keadaan ini dalam bentuk apa pun — jadi ia
   ditetapkan, bukan ditiru.

2. **DOL paritas** (tiket 06). Validasi *Date of Loss* per `Type` **dipertahankan sesuai perilaku
   Pega**. `[terbuka — catatan kecil, non-pemblokir]` Satuan pergeseran tanggal pada cabang
   `TP`/`TR` (`@addCalendar(.DATE_OF_LOSS,0,0,0,1,0,0,0)`) belum dikonfirmasi Product+UW; definisi
   `@addCalendar` tidak ada di korpus. Pergeseran dibawa apa adanya; **tidak memblokir tiket**.

3. **Aturan cutover `02/01/2026`** dari `PROC_GENERATE_SEQUENCE_NUMBER`: periode sebelum atau pada
   2 Januari 2026 → **`12.2025`**. Ini **murni logika penomoran periode**, bukan aturan bisnis
   klaim, dan **direplikasi apa adanya** di layanan penomoran Go. Bergantung pada
   `POOLDATA.TANGGAL_CLOSING` (`TANGGAL VARCHAR2(10)` — **string, bukan `DATE`**).

---

# Lampiran — Grilling Komite Claim Life, Ronde 1 (frontier kosong, 2026-09-15)

> Catatan status, bukan glossary. Rincian dan bukti di
> `.scratch/komite-claim-life/grilling-ronde-1.md`.

| # | Keputusan | Sumber |
| --- | --- | --- |
| 1 | **`TransferType` = dead code → dibuang.** Tidak pernah diisi (nol `Property-Set`); hanya dibaca di precondition `KomiteRouter` dan visible-when `ShowTransfer`. Selalu FALSE. UI yang bergantung padanya ikut dibuang. Routing tingkat = baris roster pertama ber-`KomiteAproval == 0`, sasaran `param.AssignTo = .KomiteID` | `[terverifikasi + keputusan work owner]` |
| 2 | **Penegakan pemutus per tingkat.** Hanya pemilik `KomiteList(KomiteCount).KomiteID` pada tingkat berjalan yang boleh memutuskan; lainnya ditolak **di lapisan layanan**. Pega hanya *menempatkan* tugas di worklist | `[keputusan work owner]` → **ADR-0014** |
| 3 | **Keputusan komite = enum tertutup `{1, 2}`.** Setuju = `1`, Tolak = `2`. Tidak ada opsi ketiga; nilai lain ditolak terang-terangan | `[keputusan work owner]` |
| 4 | **Cutover 7 Feb 2025 tidak dipakai lagi.** Blok `ProdDateTime < "20250207T000000.000 GMT"` sudah di-remark; jangan direplikasi. Identitas retro cukup dua precondition aktif (`Type TP/TR`, `SecurityReinsurer` terisi) | `[keputusan work owner]` |
| 5 | **Dua blok tulis kembar disatukan** menjadi satu jalur simpan berparameter status. Isi identik, beda hanya nilai status. Penyimpangan sadar dari paritas struktural | `[terverifikasi + keputusan work owner]` |
| 6 | **Nomor & rekam akseptasi dibuat sekali, di tingkat final.** Tingkat bukan-terakhir hanya mencatat `KomiteAproval = AcceptStatus` lalu `KomiteCount + 1` — naik tingkat tanpa nomor | `[terverifikasi]` |
| 7 | **Semua efek keluar wajib berhasil** (at-least-once, transactional outbox + retry). **Menyimpang dari ADR-0008**; pemicunya **Kasir** (integrasi keuangan). ⚠️ *Ralat bukti 28-09-2026: REST Kasir ter-remark di korpus (b3033) — premisnya dikembalikan ke work owner, OQ-K-06* | `[keputusan work owner]` → **ADR-0015** |
| 8 | **Tiga identitas retro ter-hardcode dibuang** (`1000013`, `L0000141`, `L0000134`) — tidak direplikasi sebagai konstanta | `[keputusan work owner]` |

## Kosakata baru dari ronde ini

**KomiteAproval**:
Penanda per baris roster: sudah menyetujui atau belum. `0` = belum. **Inilah yang menggerakkan
tangga** — tugas selalu dirutekan ke baris roster pertama yang masih `0`, sehingga graf cukup satu
Assignment yang di-loop, bukan satu shape per tingkat.
_Avoid_: Approval, Status Persetujuan Komite
_Bukti_ `[terverifikasi]`: `Komite Claim Life/Activity/KomiteRouter.xml` (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY`) — precondition `.KomiteAproval == 0` (baris ~382) menggerbangi `param.AssignTo = .KomiteID` (baris ~294). Nilainya diisi `KomitePostAdjustment` per tingkat (baris ~792, ~908)

**KomiteCount**:
Pencacah tingkat tangga yang sedang berjalan. Naik satu tiap tingkat selesai
(`KomiteCount = KomiteCount + 1`). Tangga final saat `KomiteCount == KomiteLoop`.
_Avoid_: Level, Tingkat, Urutan Komite
_Bukti_ `[terverifikasi]`: `Komite Claim Life/When/IsKomiteLoop.xml` (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `ISKOMITELOOP` / `RULE-OBJ-WHEN`) — `.AcceptStatus = "1" && .KomiteCount <= .KomiteLoop`; kenaikan di `KomitePostAdjustment.xml` baris ~9020

**AcceptStatus** *(pelengkap entri yang sudah ada di atas)*:
Di Komite Claim Life ia **tidak ditulis rule mana pun** — ia diisi lewat **dropdown wajib** di layar.
Enum tertutup `{1, 2}`.
_Bukti_ `[terverifikasi]`: nol `<PropertiesName>…AcceptStatus</PropertiesName>` di 47 berkas modul; kontrol di `Komite Claim Life/Section/ShowTransfer.xml` baris 32607 — `pyValue = .AcceptStatus`, `pyFormat = pxDropdown`, `pyRequired = true`, `pyRequiredNew = always`

**Kasir**:
Integrasi pembayaran yang *dimaksudkan* dipanggil setelah keputusan komite final. **Tidak ada di
Claim Life** — khas Komite. Ia yang membuat efek keluar Komite menuntut jaminan lebih kuat daripada
Claim Life. ⛔ *Ralat 28-09-2026 (sensus remark):* di korpus `Connect-REST`-nya (langkah 11 b3021)
**ter-remark** (`//` b3033); efek hidupnya satu baris `POOLDATA.DIRECTTOKASIR_LOG` (langkah 13 b3395)
berisi JSON pembayaran (b2774) — OQ-K-06.
_Avoid_: Cashier, Payment Service
_Bukti_ `[terverifikasi]`: `Komite Claim Life/Activity/HitServiceToKasirKMTLife_Act.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `HITSERVICETOKASIRKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`), dipanggil `KomitePostAdjustment` step 12 (b8819), digerbangi `AcceptStatus = 1 && KomiteCount == KomiteLoop` (b8887), `Type` bukan TP/TR (b8910), `IsKPR=="KPR"` (b8939), `IsPEGAPROD` (b8962) *(diralat 28-09-2026: dulu "step 11 (~8878)")*

---

# Lampiran — Grilling Komite Claim Life, Ronde 2 (2026-09-15)

> ⚠️ **Catatan audit.** `Komite Claim Life/Activity/KomitePostAdjustment.xml` **diperbarui work
> owner 2026-09-15** (515.675 byte, stempel simpan `20260915T023136`). Beberapa langkah yang
> sebelumnya di-remark kini aktif. Seluruh catatan di bawah dari versi terbaru.

## Lima keputusan Ronde 2

| # | Keputusan | Sumber |
| --- | --- | --- |
| 1 | **Idempotensi wajib.** Keempat efek keluar boleh di-retry; tiap kiriman membawa **ID idempoten unik**. Khusus **Email** dan **Kasir**: **cek status "sudah terkirim sukses" sebelum kirim ulang** — anti email/pembayaran dobel | `[keputusan work owner]` → **ADR-0015** |
| 2 | **Gerbang EXIT retro dibuang.** `1000013`, `L0000141`, `L0000134` + gerbang step 9 dihapus. Konsekuensi diterima: klaim yang tadinya di-EXIT kini menjalankan **seluruh** efek termasuk Kasir. **Semua klaim menjalankan keempat efek** | `[keputusan work owner]` → OQ-064 tertutup |
| 3 | **`SetInformationData` (step 14) bukan efek keluar** — hanya temporary penampung hasil submit (`"Approve"`/`"Reject"` + `AcceptedNo` + `CLAIM_NO`). Tidak wajib-berhasil | `[keputusan work owner]` |
| 4 | **Pemantau kegagalan** — status eksplisit di UI Komite + laporan harian | `[keputusan work owner]` → **ADR-0015** |
| 5 | **Eskalasi naik satu tingkat.** Inbox hanya menampilkan kasus sesuai posisi roster. Bila anggota absen, kasus dipindahkan **naik satu tingkat** (aksi manual admin) — pengecualian sah atas penegakan per `KomiteID` | `[keputusan work owner]` → **ADR-0014** |

## Daftar final efek keluar AKTIF `[terverifikasi]`

Empat efek, seluruhnya **setelah `Obj-Save` (step 6)**, berurutan:

| Urutan | Step | Efek |
| ---: | ---: | --- |
| 1 | **8** | `Call InsertJsonClaimLife_Act` |
| 2 | **10** | `Call serviceInsertArasapasClaimLife_act` — endpoint via `M_LINK_SERVICE`, kunci `Klaim`/`insertClaimLife` (**ADR-0013**) |
| 3 | **11** | `Call SendEmailKlaimLife` |
| 4 | **12** | `Call HitServiceToKasirKMTLife_Act` — Kasir/pembayaran |

Step **9** = gerbang EXIT (berdiri di antara efek 1 dan 2; dibuang di sistem baru).
Step **14** `SetInformationData` = **bukan** efek keluar.

## Lima langkah ter-REMARK `[terverifikasi]`

| Step | Deskripsi |
| --- | --- |
| **4.4** | Generate No Akseptasi (QP,QR) |
| **4.5** | Generate No Akseptasi (TP,TR) |
| **4.6** | Set Nilai Akseprtasi |
| **4.14** | Tukar SecurityReinsurer dengan RetroName — *jalur aksep* |
| **5.5** | Tukar SecurityReinsurer dengan RetroName — *jalur reject* |

⚠️ **4.6, 4.14, dan 5.5 melebihi daftar yang diperkirakan.** Ketiga precondition retro — termasuk
`Type=="TP"||"TR"` dan `SecurityReinsurerID!=""` yang sempat dianggap masih aktif — berada **di
dalam** langkah 4.14/5.5 yang ter-remark. Akibatnya `RetroName`/`RetroID` **tidak lagi ditukar**
sebelum rekam akseptasi ditulis → **OQ-065**.

## Jalur penomoran akseptasi yang AKTIF `[terverifikasi]`

Dengan 4.4/4.5 mati, nomor akseptasi dibuat lewat rantai aktif di dalam step 4:
**4.7** `GetKodeProdLife_SQL` (AMBIL KODE PROD) → **4.9** `GetSequenceNumber_SQL`
(generate MM.YYYY DAN SEQUENCE) → **4.11/4.12** percabangan `QR,QP` / `TR,TP` → **4.15**
`UpdateOsAkseptasiClaimLife_sql` (Insert ke OS).

**Jalur yang sama persis dengan Claim Life** (**ADR-0006**).

**"Nomor & rekam akseptasi dibuat sekali, di tingkat final" TETAP BENAR** — step 4 digerbangi
`AcceptStatus = 1 && KomiteCount == KomiteLoop` (baris 5695), step 5 oleh `AcceptStatus==2 && …`
(baris 8119), step 10 dan 12 oleh gerbang yang sama (baris 8648, 8887).

## Penutup — OQ-065 tertutup (2026-09-15)

`[keputusan work owner]` Kelima langkah ter-remark di `KomitePostAdjustment.xml` (4.4, 4.5, 4.6,
4.14, 5.5) **sengaja dimatikan**: nilai `RetroID`/`RetroName` **sudah di-set di langkah sebelumnya**,
mentah dari data policy. Nilai retro yang masuk ke `OS_AKSEPTASI_KLAIM_LIFE` adalah **nilai apa
adanya, tanpa penukaran**. Sistem baru **tidak mereplikasi logika penukaran** apa pun.

**Dengan ini Komite Claim Life MATANG — nol OQ pemblokir.**

---

# Lampiran — Grilling PremiumList Life + Endorsement Life, Ronde 1 (2026-09-15)

## ⚠️ Pemisahan konteks (2026-09-15) `[keputusan work owner]`

Lampiran ini ditulis ketika PremiumList Life dan Endorsement Life masih diperlakukan sebagai **satu**
konteks. **Itu sudah tidak berlaku.** Keduanya kini **dua konteks/menu terpisah**:

| | PremiumList Life | Endorsement Life |
| --- | --- | --- |
| Artefak | `.scratch/premiumlist-life/` | `.scratch/endorsement-life/` |
| Class work `[terverifikasi]` | `ASM-FW-GISFW-WORK-LIFE` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` |
| Flow masuk `[terverifikasi]` | `PremiumList Life/InputPolicyHolder.xml` (`INPUTPOLICYHOLDER` / `RULE-OBJ-FLOW`) | `Endorsement Life/Flow/InputEDMLife.xml` (`INPUTEDMLIFE` / `RULE-OBJ-FLOW`) |
| Alur inti | penawaran baru → premium list → simpan polis | **pilih polis NB yang sudah ada → muat data lamanya → endorse** |

**Alasan pemisahan** `[keputusan work owner]`: Endorsement punya alur inti — memilih polis new
business yang sudah ada lalu memuat data lamanya untuk di-endorse — yang **tidak ada** di
PremiumList Life. Modul, class, dan flow-nya memang sudah terpisah di korpus.

**Mesin yang dipakai BERSAMA** — ini **komponen**, bukan alasan menyatukan menu:

| Mesin bersama | Identitas | Sifat |
| --- | --- | --- |
| Penulis detail | `SaveMasterLPDet` — `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL` | **satu rule identik** di kedua modul `[terverifikasi]` |
| Penulis summary | `InsertPLSummary` — `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL` | **satu rule identik** `[terverifikasi]` |
| Tabel | `POOLDATA.M_LIFE_PREMIUM_DETAIL`, `M_LIFE_PREMIUM_SUMMARY`, `JSON_POLIS` | dipakai kedua jalur |
| Penomoran & periode | `PROC_GENERATE_SEQUENCE_NUMBER`, `KODE_PRODUKSI`, `TANGGAL_CLOSING` | dipakai kedua jalur |
| Orkestrator simpan | `InsertJsonPolisLife_Act` | **nama sama, class berbeda → rule BERBEDA** `[terverifikasi]` |

⚠️ Baris terakhir adalah jebakan klasik: `InsertJsonPolisLife_Act` **bukan** mesin bersama. Nama
identik di class berbeda = **dua rule berbeda** dengan isi berbeda (327.332 byte / 17 langkah di
`ASM-FW-GISFW-WORK-LIFE`; 444.634 byte / 16 langkah di `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE`).

> Konteks **hulu** domain Life: sumber `PremiumListSummary` / `PremiumListDetail` yang dikonsumsi
> Claim Life. Rincian dan bukti di `.scratch/premiumlist-life/grilling-ronde-1.md`.

## Enam verdict

| # | Keputusan | Sumber |
| --- | --- | --- |
| 1 | **Keputusan alur dibuat MANUAL** oleh inputor/admin, bukan formula. Yang direplikasi adalah **akibatnya**: `Confirm` → naik tahap; `Reject` → **selalu balik ke input**; `Decline` → case ditutup. `IsFlagOnGoingPolicy`: `1`=offer → berhenti di penawaran, `2`=premium → lanjut Input Premium List Detail | `[keputusan work owner]` — OQ-043 tertutup |
| 2 | **Ambang tutup buku dari tabel `POOLDATA.TANGGAL_CLOSING`**, bukan hardcode. `25` hanya **fallback** bila tabel kosong. Transaksi setelah tanggal closing dibukukan ke periode bulan berikutnya (tgl 1, 05:00 GMT = 12:00 WIB) | `[terverifikasi + keputusan work owner]` — OQ-030 tertutup |
| 3 | **Empat gerbang treaty ID `1000032`–`1000035` = logika polis lama, dibuang.** ⚠️ Di korpus masih **aktif** (blockname kosong) — dead code hidup. Jenis treaty diambil dari data **`TypeCeding`** | `[keputusan work owner + terverifikasi]` — OQ-031 tertutup |
| 4 | **Dua efek keluar aktif:** `serviceInsertArasapasLife_act` (efek bisnis) dan `SendEmailNotification` (**jalur alarm**, bukan notifikasi bisnis). `ConvertJsonNusareToProduction` **tidak dipakai** — dibuang | `[terverifikasi + keputusan work owner]` |
| 5 | **Go memegang batas transaksi eksplisit.** `Commit` di korpus ter-remark; jangan andalkan commit implisit Pega | `[terverifikasi + keputusan desain]` — konsisten OQ-013 |
| 6 | **`PL_NUMBER_EDM` = jalur penomoran EDM yang SEJAJAR dengan `PL_NUMBER` NB** — dibuat terpisah agar tidak saling timpa; **bukan pointer ke PL asal**. Field yang boleh diubah endorsement ditentukan kondisi `When` **per form**, tidak ada daftar global | `[keputusan work owner]` |

## Dead code — tidak direplikasi

`PremiumList Life/Activity/InsertJsonPolisLife_Act.xml` (`ASM-FW-GISFW-WORK-LIFE` /
`INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY`, **17 langkah**, bukan 21):

| Bagian | Penanda korpus | Kenyataan |
| --- | --- | --- |
| Step 6-7, empat gerbang ID polis lama | **AKTIF** (blockname kosong) | **tidak dipakai** `[keputusan work owner]` |
| Step 16 `Commit` | **REMARK** | mati `[terverifikasi]` |
| Step 17 `Connect-REST` `ConvertJsonNusareToProduction` | **REMARK** | mati `[terverifikasi]` |

⚠️ **Penanda `blockname` tidak dapat dipercaya sendirian di modul ini** → **OQ-066**.

## Kosakata baru

**TypeCeding**:
Kode jenis treaty pada polis Life. **`1` = QS (Quota Share), `2` = SURPLUS, `3` = QS + SURPLUS,
`4` = XOL.** Inilah sumber jenis treaty yang sah di sistem baru — bukan konstanta ID.
_Avoid_: Tipe Ceding, Jenis Cession
_Bukti_ `[terverifikasi]`: `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml` (`ASM-FW-GISFW-WORK-LIFE` / `INSERTJSONPOLISLIFE_ACT`) baris ~3787 — `@if(TypeCeding="1","QS",@if(="2","SURPLUS",@if(="3","QS + SURPLUS",@if(="4","XOL",""))))`

**PL_NUMBER / PL_NUMBER_EDM**:
Dua **jalur penomoran premium list yang sejajar**: `PL_NUMBER` untuk **NB** (new business),
`PL_NUMBER_EDM` untuk **EDM** (endorsement). Dibuat terpisah agar tidak saling timpa.
**Bukan** relasi induk-anak.
_Avoid_: Nomor PL Induk, PL Turunan
_Bukti_ `[terverifikasi]`: `PremiumList Life/RDBList/InsertPLSummary.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL`) memanggil `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY(...)` dengan keduanya sebagai **parameter terpisah**. `[keputusan work owner]` untuk sifat sejajarnya

⚠️ **Keduanya dibentuk dengan cara yang berbeda** `[terverifikasi]`:
`PL_NUMBER` lahir dari `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` (sequence terpusat, **ADR-0006**).
`PL_NUMBER_EDM` **tidak** — ia dirakit sebagai **`<nomor polis>` + `/` + `<PRODKE terakhir + 1>`**
oleh `Endorsement Life/RDBList/Generate_NoEndorsmentLife.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!GENERATE_NOENDORSMENTLIFE` / `RULE-CONNECT-SQL`):
`SELECT NOPOLIS||'/'||{InputData.CARI14} … FROM POOLDATA.JSON_POLIS … ORDER BY PRODKE DESC`.
Nomornya **lahir sekali** — `Endorsement Life/Activity/GenerateNoEDM_Life.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GENERATENOEDM_LIFE` / `RULE-OBJ-ACTIVITY`) menjaga step 5
dan 6 dengan precondition `pyWorkPage.PremiumListSummary.PL_NUMBER_EDM==""`.

**Penulisan detail peserta (`M_LIFE_PREMIUM_DETAIL`)**:
Baris peserta premium list yang **dibaca Claim Life**. Ditulis lewat **satu** `RULE-CONNECT-SQL`
untuk kedua jalur: `SaveMasterLPDet` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET`).
Yang berbeda hanyalah **pemicunya** — dan di sistem baru perbedaan itu **dihapus**.
_Avoid_: Peserta Batch, Detail Job, Tabel Peserta Klaim
_Bukti_ `[terverifikasi]`: `PremiumList Life/RDBList/SaveMasterLPDet.xml` dan `Endorsement Life/RDBList/SaveMasterLPDet.xml` adalah **rule yang sama** — identitas, `pxUpdateDateTime` (`20260211T064412.717 GMT`), dan ukuran identik; diff dua ekspor yang dinormalisasi menyisakan 8 baris cap waktu ekspor

| Jalur | Pemicu di Pega | Pemicu di sistem baru |
| --- | --- | --- |
| **NB** (new business) | **job/batch** — `InsertLifePremiumDetail_act` (`ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT` / `RULE-OBJ-ACTIVITY`) step 2 `Call Rule-Obj-Report-Definition.pxRetrieveReportData` atas `SelectNoJsonPolis_RD`, lalu loop step 3 | **INLINE saat simpan polis** |
| **EDM** (endorsement) | **inline** — `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT`) step 11.6 | **INLINE saat simpan polis** (tetap) |

⚠️ **Penyimpangan sadar** `[keputusan work owner]`: **pemicu job dihilangkan.** Logika penulisannya
**tidak** diganti — tetap `InsertLifePremiumDetail_act` — hanya waktu pemanggilannya yang disamakan
dengan EDM. Akibatnya data peserta **langsung tersedia untuk klaim tanpa jeda job**, dan kedua jalur
berperilaku seragam. Step 2 dan loop step 3 tidak dimigrasikan; yang dimigrasikan adalah **badan
per-kasus** (3.1–3.6).

`[terverifikasi]` Penjaga idempotensi sudah ada di korpus: langkah insert (step 3.3.4) dijaga
precondition `hasilDetail.pxResults(1).PL_NUMBER==""` (baris 3820) — **wajib dipertahankan** agar
pemanggilan inline yang terulang tidak menggandakan baris.

⚠️ `[terverifikasi]` `SaveMasterLPDet` **commit sendiri** (`COMMIT;` baris 252) — ia **titik potong**
transaksi, bukan bagian transaksi summary. Tunduk pada aturan urutan transaksi campuran (**OQ-013**).


**Offer / Premium**:
Dua jalur lanjutan setelah penawaran diterima. **Offer** = berhenti di tahap penawaran;
**Premium** = lanjut mengisi Premium List Detail. Ditentukan `IsFlagOnGoingPolicy` (`1`/`2`).
_Avoid_: Quotation, Penawaran Saja
_Bukti_ `[terverifikasi]` nilai keluaran `ASM-FW-GISFW-WORK-LIFE` / `ISFLAGONGOINGPOLICY`; `[keputusan work owner]` arti `1`/`2`

## Kontrak hulu ke Claim Life `[terverifikasi]`

PremiumList Life dan Endorsement Life **menghasilkan** baris premium summary/detail lewat stored
procedure `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY`, berkunci `PL_NUMBER` (NB) atau `PL_NUMBER_EDM`
(endorsement). Claim Life **mengonsumsinya** sebagai `PremiumListSummary` / `PremiumListDetail` —
struktur yang dipakai mesin status klaim (**ADR-0011**).

Class integrasi `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` adalah **tulang punggung domain Life**:
Claim Life 46 rule, Endorsement Life 25, PremiumList Life 19, Komite Claim Life 11.

⚠️ **Batas pengetahuan:** badan `PEGA_M_LIFE_PREMIUM_SUMMARY`, `INSERTJSONPOLISLIFE`, dan
`INSERTJSONOFFERLIFE` **tidak ada di korpus** — keluarga **OQ-002**.

## Penutup — OQ-002 tertutup untuk PremiumList Life + Endorsement Life (2026-09-15)

`[data DBA]` Body **tiga procedure premium** diserahkan:

| Procedure | Sasaran | Commit sendiri? |
| --- | --- | --- |
| `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY` | `INSERT` ke `M_LIFE_PREMIUM_SUMMARY` (PK dari `M_LIFE_PREMIUM_SUMMARY_SEQ`), **37 kolom** | **tidak** |
| `POOLDATA.INSERTJSONPOLISLIFE` | **upsert** `JSON_POLIS` by `IDPEGA`; fallback `JSON_POLIS_ERROR` | **ya** |
| `POOLDATA.INSERTJSONOFFERLIFE` | **upsert** `JSON_OFFER_LIFE` by (`IDPEGA`, `STATUS`) | **ya** |

⚠️ **Seluruh parameter `PEGA_M_LIFE_PREMIUM_SUMMARY` bertipe `VARCHAR2` — termasuk kolom uang.**
Uang menyeberang batas sebagai **teks**. Memperkuat **ADR-0003**: konversi ke desimal presisi
arbitrer dilakukan **eksplisit di satu batas**, tidak pernah lewat `float`.

⚠️ **Batas transaksi jalur Life CAMPURAN.** `PEGA_M_LIFE_PREMIUM_SUMMARY` dan
`PROC_GENERATE_SEQUENCE_NUMBER` tidak commit sendiri; `INSERTJSONPOLISLIFE` dan
`INSERTJSONOFFERLIFE` **commit sendiri** — keduanya **memutus transaksi Go**. Data yang harus
atomik **tidak boleh** dipisahkan oleh procedure yang commit sendiri. Berbeda dari Claim Life, yang
transaksinya dapat utuh.

**Sisa OQ-001 untuk konteks ini:** DDL fisik `M_LIFE_PREMIUM_SUMMARY`, `JSON_POLIS`,
`JSON_OFFER_LIFE` — tipe, presisi, PK, index. **Kolomnya sudah terbaca; hanya tiket migrasi yang
menunggu.**

**Dengan ini PremiumList Life + Endorsement Life MATANG — nol OQ pemblokir.**

---

# Lampiran — Grilling Endorsement Life, Ronde 1–2 (2026-09-15)

> Konteks **terpisah** dari PremiumList Life `[keputusan work owner]`. Class work
> `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE`, flow masuk `Endorsement Life/Flow/InputEDMLife.xml`.
> Rincian dan bukti di `.scratch/endorsement-life/grilling-ronde-1.md` dan `grilling-ronde-2.md`.

## Dua belas verdict Ronde 1

| # | Keputusan | Sumber |
| --- | --- | --- |
| 1 | **Titik masuk = nomor polis** (`NOPOLIS` / `TempWork.PolicyNo`). `BrowsePremiumList_RD` hanya alat bantu pencarian, bukan sumber kunci | `[keputusan work owner]` |
| 2 | **Lima gerbang kelayakan** berlaku (`SetErrorBatalEndorsement_Act` — **namanya menyesatkan, diganti**): nomor polis kosong; polis tak ada di `JSON_POLIS`; **sudah ada EDM belum resolve**; sudah pernah EDM batal; sudah ada pembayaran + `EdmType="3"` | `[keputusan work owner + terverifikasi]` |
| 3 | **`EdmType`: `1` = Perubahan Data, `3` = Batal**; `2`/`4` tidak dipakai — enum efektif `{1,3}` | `[keputusan work owner]`; nilai `3` `[terverifikasi]` |
| 4 | **`IVD_JR_ID = '5'` = pembayaran/pelunasan.** Sistem baru tetap membaca `ARASAPAS.DETAIL_INVOICE` langsung, **dikurung di satu repository** bertanda batas lintas sistem | `[keputusan work owner + desain]` |
| 5 | **Endorsement tidak punya `Reject`** — hanya `Confirm`/`Decline`, beda dari PremiumList Life | `[keputusan work owner + terverifikasi]` |
| 6 | **`GetOldDetail_EDM` tidak dimigrasikan** sebagai rule. ⚠️ Perilakunya tetap terpakai (lihat Ronde 2) | `[keputusan work owner]` |
| 7 | **`EDMStatus` = status per baris**; baris polis NB **tidak dihapus fisik**. ⚠️ **Dikoreksi V13:** `"Delete"` bukan sekadar penanda — ia **juga** menghasilkan nilai negatif (minus selektif per peserta) | `[keputusan work owner + terverifikasi]` |
| 8 | **Pembacaan `PRODKE` disatukan ke `ORDER BY PRODKE DESC`** di kedua pembaca — penyimpangan sadar | `[keputusan desain]` |
| 9 | Case EDM dibuat **setelah lima gerbang lolos**; **membatalkan endorsement = `Decline`** (`Resolved-Rejected`), polis bebas di-endorse ulang. `CancelCreateCaseEDML` bukan mekanisme batal | `[keputusan work owner]` |
| 10 | **Alarm dihidupkan** di jalur endorsement (deteksi separuh + email), seragam dengan new business | `[keputusan work owner]` |
| 11 | **Penjaga anti-dobel `(NOPOLIS, PRODKE)`** sebelum insert `JSON_POLIS` — penyimpangan sadar | `[keputusan desain]` |
| 12 | **Keempat jenis `QP`/`QR`/`TP`/`TR` didukung**; `QR` memakai section induk `ViewOldPolicy_EDM` | `[terverifikasi]` (dikuatkan Ronde 2) |

## Penyimpangan sadar Endorsement Life

| # | Penyimpangan | Alasan |
| --- | --- | --- |
| 1 | `ORDER BY PRODKE DESC` di **kedua** pembaca | Pega punya dua urutan berbeda (`PRODKE DESC` vs `TGL_INPUT desc`) → nomor EDM bisa salah |
| 2 | **Penjaga anti-dobel** `(NOPOLIS, PRODKE)` sebelum insert `JSON_POLIS` | `InsertJsonPolisEDM` adalah `INSERT` polos + commit sendiri → pengulangan menggandakan rekam |
| 3 | **Alarm dihidupkan** | Pega menjalankan endorsement tanpa alarm; jalur NB memilikinya |
| 4 | Baca Arasapas **terkurung di satu repository** bertanda batas lintas sistem | mencegah kueri lintas skema tersebar |
| 5 | Rule gerbang **diganti nama** (bukan "SetErrorBatal…") | nama Pega menyesatkan: ia gerbang kelayakan, bukan pembatalan |

## Kosakata baru

**EdmType**:
Maksud sebuah endorsement Life. **`1` = Perubahan Data**, **`3` = Batal** (pembatalan polis).
Nilai `2` dan `4` tidak dipakai — enum efektif `{1, 3}`. Disimpan **di dalam CLOB**
`POOLDATA.JSON_POLIS.DATA_JSON`, **bukan** sebagai kolom.
_Avoid_: Tipe EDM, Jenis Endorsemen, Endorsement Type
_Bukti_ `[terverifikasi]` nilai `3`: `Endorsement Life/Activity/SetPremi_EDM.xml` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETPREMI_EDM` / `RULE-OBJ-ACTIVITY`) baris 1273 — precondition `.EdmBatal=="True" || pyWorkPage.EdmType==3` di bawah deskripsi "Set 0 jika EDM Batal"; dan `Endorsement Life/RDBList/GetEdmTypeLife.xml` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETEDMTYPELIFE` / `RULE-CONNECT-SQL`) — `SELECT A.DATA_JSON.EdmType … FROM POOLDATA.JSON_POLIS A`. `[keputusan work owner]` untuk arti nilai `1` dan ketidakterpakaian `2`/`4`

**EDMStatus**:
Status **per baris detail peserta** di dalam sebuah endorsement. **Empat** nilai, dan bedanya adalah
**cakupan minus — bukan jenis**: **`"Old"`** (warisan polis new business, nilai tidak diubah),
**`"New"`** (peserta ditambah — hanya pada Perubahan Data `EdmType=1`, nilai positif),
**`"Delete"`** (peserta dihapus dalam Perubahan Data — **nilainya diminuskan, selektif PER
PESERTA**), **`"Batal"`** (lewat EDM Batal `EdmType=3` — **seluruh peserta otomatis batal, nilai
diminuskan MENYELURUH**).
⚠️ **`Delete` bukan sekadar penanda** — ia **juga** menghasilkan nilai negatif, lewat jurnal balik
`× -1` yang sama di `SetPremi_EDM`. Yang tetap benar: baris polis NB **tidak dihapus fisik**.
_Avoid_: Status Baris, Flag Detail, Row Status
_Bukti_ `[terverifikasi]`: `"Old"` — `MappingEDMLife` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `MAPPINGEDMLIFE`) baris 2608; `"New"` — `SaveCSVEDMLife` (`…` / `SAVECSVEDMLIFE`) baris 2714; `"Delete"` — `SetPremi_EDM` baris 1384; `"Batal"` — `SetPremi_EDM` baris 1506. `[keputusan work owner]` untuk cakupan minus (OQ-070 tertutup)

**Jurnal balik endorsement**:
Pembatalan atau penghapusan peserta **tidak menimpa** baris lama dan **tidak menulis nol** — ia
menulis **baris bernilai negatif** ke tabel yang **sama**. Baris positif asli dan baris negatif
hidup berdampingan; **net akunting** diperoleh dari penjumlahan. ⚠️ Deskripsi langkah Pega
**"Set 0 jika EDM Batal" menyesatkan**; kodenya mengalikan **32 kolom uang dengan `-1`**.
_Avoid_: Reversal, Void, Set Nol
_Bukti_ `[terverifikasi]`: `Endorsement Life/Activity/SetPremi_EDM.xml` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETPREMI_EDM` / `RULE-OBJ-ACTIVITY`) step 2.1, precondition `.EdmBatal=="True" || pyWorkPage.EdmType==3` (baris 1273), 32 penugasan `<kolom> = <kolom> * -1`. `[keputusan work owner]` untuk sifat akuntingnya (OQ-071 tertutup)

**Peserta hidup (bagi Claim Life)**:
Peserta yang **belum** dibatalkan dan **belum** dihapus endorsement. ⚠️ **Kontrak batas:** peserta
ber-`EDMSTATUS` **`"Batal"`** atau **`"Delete"`** **TIDAK BOLEH MUNCUL** di Claim Life — jalur baca
klaim **wajib menyaringnya keluar**. Akuntansi melihat **seluruh** baris; klaim hanya peserta hidup.
**Satu tabel, dua sudut pandang — dan itu disengaja.**

⚠️ **Baris new business masuk dengan `EDMSTATUS` kosong/NULL** dan **tetap peserta hidup**.
`[terverifikasi]` Jalur NB tidak pernah mengisi kolom itu: sensus `InsertLifePremiumDetail_act`
(`ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT` / `RULE-OBJ-ACTIVITY`) — **nol**
kemunculan `EDMStatus`. Karena itu penyaring naif `EDMSTATUS NOT IN ('Delete','Batal')` **membuang
seluruh peserta new business** di Oracle. Aturan yang benar: hidup = **kosong/NULL, `Old`, atau
`New`**; mati = `Delete` atau `Batal`.

⚠️ `[terverifikasi]` **Kontrak ini belum ditegakkan di korpus.** `GetPesertaClaim_sql1` menyaring
hanya dengan `PL_NUMBER`, `CERTIFICATE_NO`, dan `NAME_OF_INSURED` — **tanpa** penyaring status;
sensus modul `Claim Life/` menemukan **nol** kemunculan `EDMSTATUS`. Apakah Pega menyaring di
lapisan lain **tidak terbukti dari korpus** — **jangan menebak mekanismenya**; tegakkan aturannya di
sistem baru.
_Avoid_: Peserta Aktif, Peserta Valid
_Bukti_ `[keputusan work owner]` verdict V14 grilling Endorsement Life (OQ-071); konsumennya `Claim Life/RDBList/GetPesertaClaim_sql1.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `RNM!GETPESERTACLAIM_SQL1` / `RULE-CONNECT-SQL`) lewat `Claim Life/Activity/LoadDataPesertaSpesifik_Act.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `LOADDATAPESERTASPESIFIK_ACT` / `RULE-OBJ-ACTIVITY`, baris 545), dipasang di `Claim Life/Section/InputRegisterClaimLife.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `INPUTREGISTERCLAIMLIFE`)

**Kolom status di `M_LIFE_PREMIUM_DETAIL`**:
Tabel ini punya **tiga** kolom berakhiran status, dan **hanya satu** yang menandai hidup/mati:

| Kolom | Diisi dari | Nilai | Penanda hidup/mati? |
| --- | --- | --- | --- |
| **`EDMSTATUS`** | `TempValue.EDMStatus` | `Old` / `New` / `Delete` / `Batal`, atau kosong untuk NB | ✅ **ya** |
| `STATUS` | `CARI48` | `0` untuk `QR`/`QP`, `1` untuk `TP`/`TR` | ❌ — penanda **jenis transaksi** |
| `STATUSOLD` | `CARI47` | `1` bila "Statusold = Old", `0` bila bukan | ❌ |
_Avoid_: Flag Status, Status Baris Detail
_Bukti_ `[terverifikasi]`: daftar kolom `INSERT` dan klausa `VALUES` di `Endorsement Life/RDBList/SaveMasterLPDet.xml` dan `PremiumList Life/RDBList/SaveMasterLPDet.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL`); nilai `STATUS` jalur NB `@if(Type=="QR","0",@if(Type=="QP","0","1"))` di `InsertLifePremiumDetail_act`, sejalan dengan sub-langkah "Status QR/QP" dan "Status TR/TP" di `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml`
⚠️ `[terbuka]` **Tipe dan nullability `EDMSTATUS` belum terbaca** — `NULL` atau string kosong untuk baris NB. Masuk **OQ-001 (sisa)**, pemilik **DBA**

**IVD_JR_ID = '5'**:
Penanda **pembayaran/pelunasan** pada `ARASAPAS.DETAIL_INVOICE`. Dipakai gerbang kelayakan
endorsement: polis yang **sudah dibayar tidak boleh dibatalkan** (`EdmType="3"`).
_Avoid_: Kode Jurnal 5, Status Bayar
_Bukti_ `[terverifikasi]` kuerinya: `Endorsement Life/RDBList/SearcStatusBayarArasaps_SQL.xml` (`ASM-FW-GISFW-INT-POLICYJSON` / `ASM!SEARCSTATUSBAYARARASAPS_SQL` / `RULE-CONNECT-SQL`) — `select * from ARASAPAS.DETAIL_INVOICE where inv_inv_no = {InputData.CARI18} and IVD_JR_ID = '5'`. `[keputusan work owner]` untuk artinya

**EditInput / EditInput1**:
**Kunci satu arah** atas field endorsement. Begitu data polis lama dipetakan atau premi ditetapkan,
`.EditInput = 1` dan **nomor polis, jenis batal, serta centang batal terkunci permanen** pada case
itu; `.EditInput1 = 1` mengunci kelompok field setelah CSV disimpan. **Tidak ada** rule yang
mengembalikannya ke `0` — jalan keluarnya `Decline` lalu mulai baru.
_Avoid_: Flag Edit, Lock Input
_Bukti_ `[terverifikasi]`: penulis tunggal ke nilai `1` — `MappingEDMLife` baris 2956, `SetPremi_EDM` baris 1338, `SaveCSVEDMLife` baris 2899 (`.EditInput1`); pemakainya `<pyDisabledWhen>` di `Endorsement Life/Section/InputEDMLife.xml` (`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` / `RULE-OBJ-HTML-SECTION`) baris 3166 (`.PolicyNo`), 7537 (`.EdmTypeBatal`), 15763 (`.EdmBatal`), 8965 & 10403 (`.EditInput1`)

## Mesin bersama PremiumList Life ↔ Endorsement Life — **daftar diperbarui**

Melengkapi daftar di §"Pemisahan konteks":

| Mesin bersama | Identitas | Sifat |
| --- | --- | --- |
| Penulis detail | `SaveMasterLPDet` — `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL` | satu rule identik |
| Penulis summary | `InsertPLSummary` — `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL` | satu rule identik |
| Tabel | `POOLDATA.M_LIFE_PREMIUM_DETAIL`, `M_LIFE_PREMIUM_SUMMARY`, `JSON_POLIS` | dipakai kedua jalur |
| Penomoran & periode | `PROC_GENERATE_SEQUENCE_NUMBER`, `KODE_PRODUKSI`, `TANGGAL_CLOSING` | dipakai kedua jalur |
| Orkestrator simpan | `InsertJsonPolisLife_Act` | **nama sama, class berbeda → rule BERBEDA** |

⚠️ **`Calculate1_Act` BUKAN mesin bersama.** `[keputusan work owner]` Meski berkasnya **identik** di
kedua folder modul (`ASM-FW-GISFW-WORK-LIFE` / `CALCULATE1_ACT` / `RULE-OBJ-ACTIVITY`, 302.597 byte,
diff ternormalisasi nol baris), **di jalur ENDORSEMENT ia TIDAK dijalankan** — ia milik
**PremiumList Life (new business)** saja. Perhitungan endorsement dikerjakan **`SetPremi_EDM`**
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETPREMI_EDM`).

⚠️ **Pelajaran metodologi, seiring OQ-066: berkas identik ≠ dipakai.** Sebagaimana
`<pyStepsBlockName>` tidak boleh dipakai sendirian untuk menyimpulkan hidup/mati, **kesamaan berkas
antar folder tidak boleh dipakai sendirian untuk menyimpulkan "mesin bersama"**. Rujukan
`<RequestType>` / `Call` pun hanya menunjukkan **kemungkinan** pemanggilan — konfirmasi work owner
tetap diperlukan.

## Enam verdict Ronde 2 — **frontier kosong, OQ-070/071/072 tertutup**

| # | Keputusan | Sumber |
| --- | --- | --- |
| 13 | **`EDMStatus` empat nilai; bedanya CAKUPAN MINUS** — `Delete` minus **selektif per peserta**, `Batal` minus **menyeluruh**. ⚠️ `Delete` **juga** menghasilkan nilai negatif | `[keputusan work owner]` — OQ-070 tertutup |
| 14 | **Baris negatif ke tabel yang SAMA** (`M_LIFE_PREMIUM_DETAIL` / `SUMMARY`), berdampingan dengan baris positif → **net akunting**. ⚠️ **Peserta batal/delete TIDAK muncul di Claim Life** — jalur baca klaim wajib menyaring | `[keputusan work owner]` — OQ-071 tertutup |
| 15 | **Popup polis lama TETAP ADA** (QR→induk, QP/TP/TR→varian). Rule `GetOldDetail_EDM` tidak ditiru, **perilakunya wajib ada** | `[keputusan work owner]` |
| 16 | **Kunci field permanen** (`.PolicyNo`, `.EdmTypeBatal`, `.EdmBatal`) + **pesan penjelas** | `[keputusan work owner + desain]` |
| 17 | **Tanpa batas jumlah baris CSV** — 50.000 dibuang; unggahan besar diproses **bertahap** | `[keputusan work owner]` — OQ-072 tertutup |
| 18 | **`Calculate1_Act` TIDAK dipakai di endorsement** — bukan mesin bersama | `[keputusan work owner]` |

⚠️ **Editabilitas per form bukan rule `When`.** `[terverifikasi]` Modul ini hanya punya dua rule
`When` (`@BASECLASS` / `ISPEGAPROD` dan `@BASECLASS` / `RECORDEVENT`), keduanya bukan soal
editabilitas. Mekanismenya adalah **kondisi sebaris `<pyDisabledWhen>`** di Section.

**Endorsement Life MATANG** — nol OQ pemblokir. Sisa yang menyentuh konteks ini hanyalah
**OQ-001 (sisa)**: DDL fisik tabel premium, yang memblokir **tiket migrasi** saja.
