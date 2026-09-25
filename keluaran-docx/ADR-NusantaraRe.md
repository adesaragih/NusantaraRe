# Keputusan Arsitektur - Nusantara Re

> Seluruh keputusan arsitektur proyek migrasi, dimuat **utuh**. ADR pendek; memangkasnya
> menghilangkan alasannya.

## Cara membaca berkas ini

Keputusan lahir di **tiga rangkaian kerja yang berjalan berdampingan**, dan ketiganya
memulai penomoran dari `0001`. Nomor yang sama karena itu dapat berarti keputusan yang
berbeda. Di dalam dokumen ini tiap rangkaian diberi **awalan seri** supaya rujukan tidak
bertabrakan. **Nol berkas diganti nama** - awalan hanya hidup di dokumen ini.

| Bagian | Seri | Awalan | Jumlah |
| --- | --- | --- | ---: |
| I | keputusan lintas-modul | `ADR-U-nnnn` | **42** |
| II | keputusan per modul, rangkaian kedua | `ADR-D-<modul>-nnnn` | **56** |
| III | keputusan Fakultatif | `ADR-F-nnnn` | **7** |
| | | **Jumlah** | **105** |

**Lampiran 1** memuat konkordansi - nomor asli, awalan seri, judul, dan berkas asalnya.
**Lampiran 2** menyandingkan keputusan dari seri berbeda yang membahas hal yang sama.
[!] **Lampiran 2 tidak memutuskan mana yang berlaku.**

---

# Bagian I - Keputusan lintas-modul

Jumlah: **42** keputusan.

## ADR-U-0001 - Claim — Life dan Komite Life adalah dua konteks terpisah, dihubungkan kontrak child work

**Claim — Life** dispesifikasikan sebagai bounded context tersendiri; **Komite Life** diperlakukan
sebagai **konteks/sistem luar**. Isi dan alur internal Komite **tidak** dispesifikasikan di sini.
Hubungan keduanya berupa **tiga kontrak eksplisit**: penyerahan kasus sebagai *child work*, jalur
balik hasil keputusan, dan tabel akseptasi bersama.

#### Considered Options

- **(a) Claim Life saja, Komite sebagai konteks luar** — dipilih
- (b) Claim Life + Komite Claim Life sebagai satu konteks
- (c) Claim Life saja, tetapi rule tulis bersama ikut dispesifikasikan di dalamnya

(b) ditolak karena Komite membawa **13 OQ pemblokir** yang hampir seluruhnya RBAC
(`discovery/D3-D4-CLOSING-REPORT.md` §4.2); menariknya masuk akan mengubah konteks tersiap
(5 pemblokir, cakupan bukti `full`) menjadi yang paling terblokir.

#### Kontrak 1 — penyerahan kasus ke Komite `[terverifikasi]`

Claim Life menyerahkan kasus dengan **membuat child work** berkelas
`ASM-FW-GCNMFW-Work-KomiteLife`.

| Fakta | Bukti |
| --- | --- |
| Activity pembuat | `Claim Life/Activity/CreateKMTLife_Act.xml` → `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE!CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`, 121.652 byte |
| Sepuluh langkah | `Property-Set` → **`Call pxRetrieveReportData`** → `Property-Set` ×3 → **`Call pxAddChildWork`** → `Obj-Refresh-And-Lock` → `Property-Set` → **`Obj-Save`** → **`Call SendEmailKlaimLF`** |
| Kelas anak | `ASM-FW-GCNMFW-Work-KomiteLife` (9 rujukan `<pyRuleName>`, 15 kemunculan total) |
| Halaman induk | `pyWorkPage` = `ASM-FW-GCNMFW-Work-ClaimLife` |
| **Dipicu dari UI, bukan dari Flow** | `<pyActivity>CreateKMTLife_Act</pyActivity>` di `Claim Life/Section/ClaimComite.xml` (2×) dan `Claim Life/Harness/Committe_Life.xml` (2×) |

Inilah sebabnya penyerahan **tidak terlihat di graf flow** `Claim Life/Flow/Register_Flow.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` / `RULE-OBJ-FLOW`) — ia bukan shape.

##### Pemilihan roster komite `[terverifikasi]`

Langkah 2 memanggil `pxRetrieveReportData` dengan report
**`ASM-FW-GCNMFW-INT-EMAILKOMITE!FILTEREMAILKOMITEWITHLIMIT`**
(`Claim Life/ReportDefinition/FilterEmailKomiteWithLimit.xml`, `RULE-OBJ-REPORT-DEFINITION`,
67.827 byte), lewat `Param.pyReportName` / `Param.pyReportClass`.

Report itu menerima **`Param.LIMIT_BOTTOM`** dan **`Param.STS_KLAIM`** — artinya **roster dipilih
berdasarkan pita nilai**, pola yang sama dengan `Komite Claim FacIn/Activity/ApprovalKomite_Act.xml`
(OQ-037). **Ambang dan aturannya sendiri belum terverifikasi.**

##### Muatan yang menyeberang `[terverifikasi]`

Properti yang diisi `CreateKMTLife_Act` sebelum `pxAddChildWork`:

```
childPageKomite.CLMNO            childPageKomite.KomiteCount
childPageKomite.KomiteLoop       childPageKomite.IndexAdjustment
childPageKomite.IndexPremiumList childPageKomite.KomiteList(<APPEND>).KomiteID
childPageKomite.KomiteList(<LAST>).IDKomite
childPageKomite.KomiteList(<LAST>).KomiteAproval
childPageKomite.KomiteList(<LAST>).KomiteEmail
```

Di sisi induk: `.IsKomite`, `.KomiteNo`, `.TotalKomite`.

Perintah audit:
```
grep -ohE "<PropertiesName>[^<]*" "Claim Life/Activity/CreateKMTLife_Act.xml" | sed 's/<[^>]*>//' | sort -u
```

##### Tambahan pada muatan — keputusan Ronde 2 Q15

`[terverifikasi work owner 2026-09-14]` Daftar di atas **benar** sebagai muatan sekarang, tetapi
**tidak cukup**. Kontrak baru **menambahkan tiga hal**:

| Tambahan | Alasan |
| --- | --- |
| **nilai klaim** (`CLAIM_AMOUNT`) | Komite memilih roster berdasarkan pita nilai (`Param.LIMIT_BOTTOM`); nilainya sekarang tidak ikut menyeberang |
| **`CURRENCY`** | nilai uang tanpa mata uang tidak dapat dibandingkan terhadap pita nilai — lihat **ADR-U-0003** |
| **`STS_REJECT`** saat penyerahan | keadaan klaim pada saat diserahkan, agar jalur balik punya titik awal yang eksplisit |

Ini **penyimpangan sadar** dari paritas, sejenis dengan **ADR-U-0007**: alurnya tidak berubah, yang
bertambah adalah apa yang terekam menyeberang.

`[terverifikasi work owner 2026-09-14]` **`KomiteLoop` ditentukan oleh Claim — Life**, bukan oleh
Komite — sesuai arah muatan di atas (`childPageKomite.KomiteLoop` diisi induk). **Apa yang
menentukan nilainya masih OQ-032.**

#### Kontrak 2 — jalur balik dari Komite `[terverifikasi]`

Kontraknya **dua arah**. Hasil keputusan Komite kembali ke Claim — Life sebagai perubahan
`STS_REJECT`:

`[terverifikasi work owner 2026-09-14]` Admin menyisipkan baris `AdjustmentList` dengan
`STS_REJECT = 0`. Saat Claim Analis mengirim ke Komite: bila Komite **aksep** (`AcceptStatus = 1`)
maka `STS_REJECT` menjadi `1`; bila **reject** (`AcceptStatus = 2`) maka `STS_REJECT` menjadi `2`.

`[terverifikasi]` **Dikuatkan bukti korpus.** `Komite Claim Life/Activity/KomitePostAdjustment.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE!KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`, 513.661 byte) memuat:

| Yang ditemukan | Jumlah |
| --- | ---: |
| `Property-Set` → `…STS_REJECT` bernilai literal `1` | 6 |
| `Property-Set` → `…STS_REJECT` bernilai literal `2` | 2 |
| Precondition `pyWorkPage.AcceptStatus = 1 && pyWorkPage.KomiteCount == pyWorkPage.KomiteLoop` | 4 |
| Precondition `pyWorkPage.AcceptStatus==2 && pyWorkPage.KomiteCount == pyWorkPage.KomiteLoop` | 1 |
| Gerbang baris yang boleh diubah: `.ACCEPTEDNO=="" && .IsCheck = true && .STS_REJECT == "0"` | 4 |

Perintah audit:
```
f="Komite Claim Life/Activity/KomitePostAdjustment.xml"
grep -n "AcceptStatus" "$f"
grep -n "<PropertiesName>.*STS_REJECT<" "$f"   # nilai ada di baris berikutnya
```

**`[dugaan]` yang tersisa:** pemasangan tepat precondition mana menggerbangi penulisan mana tidak
dapat dibuktikan dari urutan XML — Pega menserialkan parameter langkah dalam `rowdata` bersarang,
sehingga urutan berkas bukan urutan langkah. Yang **terbukti** adalah keberadaan pasangan nilai
`{1, 2}` di kedua sisi dan gerbangnya. Pemetaannya sendiri berasal dari work owner.

**Dua fakta penting yang ikut terbaca:**

1. `[terverifikasi]` Perubahan `STS_REJECT` oleh Komite hanya terjadi pada **rung terakhir**
   (`KomiteCount == KomiteLoop`) — tingkat antara tidak mengubah status klaim.
2. `[terverifikasi]` **Komite bukan satu-satunya penulis `STS_REJECT = 1`.**
   `Claim Life/Activity/SaveAdjustment_Act.xml`
   (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!SAVEADJUSTMENT_ACT` / `RULE-OBJ-ACTIVITY`, 179.221 byte)
   punya satu langkah `Property-Set` **tanpa precondition** yang menulis, pada baris yang baru
   ditambahkan:

   ```
   .AdjustmentList(<LAST>).ACCEPTEDNO        = InputData.pxResults(1).HASIL1
   .AdjustmentList(<LAST>).STS_REJECT        = 1
   .AdjustmentList(<LAST>).ACCEPTATION_DATE  = @CurrentDateTime()
   ```

   Artinya Claim — Life dapat mengaksep **tanpa** melalui Komite. Kapan jalur itu dipakai dan kapan
   harus lewat Komite adalah **OQ-039** — pertanyaan yang sama, kini dengan bukti bahwa dua jalur
   memang ada berdampingan.

#### Kontrak 3 — tabel akseptasi bersama `[terverifikasi]`

`OS_AKSEPTASI_KLAIM_LIFE` ditulis oleh **satu rule yang sama** dari kedua sisi:
`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` /
`RULE-CONNECT-SQL`, hash ternormalisasi **`c50bfd9a12`** identik di
`Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` dan
`Komite Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml`, dan **tidak terdaftar** di register
konflik `discovery/inventory/_oq011-konflik-isi.md`.

Jadi **bukan dua penulis independen** — satu rule dipanggil dua sisi. 55 kolomnya terbaca di
`discovery/flows/Komite Claim Life.md` §3.1.

#### Consequences

- Spesifikasi Claim — Life berhenti di batas: membuat child work, menerima hasil keputusan, dan
  menulis tabel akseptasi.
- `AcceptStatus` adalah **kosakata konteks Komite**, bukan status internal Claim — Life. Ia dipetakan
  ke `STS_REJECT` **di batas kontrak**, bukan disimpan sebagai status kedua (Ronde 2 Q8).
- Perubahan pada struktur `KomiteList`, pada muatan penyerahan, atau pada
  `OS_AKSEPTASI_KLAIM_LIFE` adalah **perubahan kontrak lintas konteks**, bukan perubahan internal.
- Nilai uang yang menyeberang harus memakai representasi yang sama di kedua sisi (**ADR-U-0003**).
- Setiap transisi lewat jalur balik ini wajib merekam siapa + kapan (**ADR-U-0007**).
- `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` adalah **satu-satunya salinan di
  korpus** dan dipanggil Komite Claim Life (**OQ-035**) — arah ketergantungan ini berlawanan dengan
  kontrak di atas dan perlu ditetapkan saat Komite dispesifikasikan.

#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-039** (dipersempit 2026-09-14) | **Kapan** penyerahan dipicu. Akibat reject Admin sudah terjawab Ronde 4; `SaveAdjustment_Act` dinyatakan dead. Sisa: pemicu Adjustment/Close/Reject di 3 modul Claim lain. Pemilik: Product+UW |
| **OQ-035** | Activity dipanggil lintas modul tetapi salinannya hanya ada di satu modul |

**Tertutup untuk Claim — Life sejak 2026-09-14** — tidak lagi menyentuh ADR ini:

| OQ | Verdict |
| --- | --- |
| **OQ-032** | **`KomiteLoop` = COUNT roster `EMAILKOMITE` aktif ber-`LIMIT_BOTTOM <= CLAIM_AMOUNT`** — data-driven, tanpa konstanta. Rinci di **ADR-U-0012** |
| **OQ-037** | Ambang adalah **data di tabel** `EMAILKOMITE`, bukan hardcode |
| **OQ-060** | Nilai uang menyeberang sebagai `(amount, currency)` per baris, invariant satu klaim satu mata uang (**ADR-U-0003**) |
| **OQ-013** | `COMMIT` ada di **pembungkus Pega**, bukan di procedure. Sistem baru: **Go memegang batas transaksi** |
| **OQ-061** | Unit keputusan = baris `AdjustmentList` (**ADR-U-0011**); `PremiumListDetail` dan header klaim adalah cerminan |

## ADR-U-0002 - RBAC Claim — Life memakai tiga peran yang sudah ada; rangkap peran ditolak

Model peran untuk Claim — Life **dirancang**, bukan dimigrasikan — korpus Pega tidak memuat satu pun
rule identitas atau otorisasi (**OQ-007**). Bahannya diambil dari tiga **kode peran** yang sudah
dipakai modul ini, dikonfirmasi work owner sebagai **daftar lengkap**. **Rangkap peran tidak
diperbolehkan**, kecuali akses ditambahkan eksplisit pada role akun.

#### Peran dan pemetaan tahap

`[terverifikasi work owner 2026-09-14]`

| Tahap (shape pada `Register_Flow`) | Peran |
| --- | --- |
| Register (`Assignment2`) + Outstanding (`Assignment1`) | **`ReasLifeAdmin`** |
| Medical Check (`Assignment3`) | **`ReasLifeMedicalAdvisor`** |
| Claim Analis (`Assignment4`) | **`ReasLifeSPV`** |

Jalur balik yang melekat pada peran: `SendtoAdmin = 1` mengembalikan kasus dari
`ReasLifeMedicalAdvisor` **atau** `ReasLifeSPV` ke `ReasLifeAdmin`; `SendtoMedical = 1`
mengembalikan dari `ReasLifeSPV` ke `ReasLifeMedicalAdvisor`.

#### Mengapa tiga kode ini, bukan peran baru

`[terverifikasi]` Claim — Life **satu-satunya konteks di korpus yang otorisasinya sudah berbasis
peran, bukan identitas orang**:

- `pyPosition` dibandingkan terhadap `'ReasLifeAdmin'`, `'ReasLifeSPV'`,
  `'ReasLifeMedicalAdvisor'` di **17 berkas** `Claim Life`, mis.
  `pyPosition != 'ReasLifeAdmin' || pyWorkPage.ClaimData.…` (8 kemunculan),
  `pyPosition == 'ReasLifeMedicalAdvisor'` (5), `pyPosition == 'ReasLifeSPV'` (6).
- `OperatorID.pyUserIdentifier` (2 berkas) dan `pyUserName` (4 berkas) dipakai sebagai **data jejak**,
  **bukan** guard terhadap literal nama orang:
  ```
  grep -rhoE "(pyUserIdentifier|pyUserName)[^<]{0,45}" "Claim Life" --include="*.xml" \
    | sed 's/&amp;#61;/=/g;s/\]\[/ /g;s/[][]//g' | grep -E "[=!]"      # -> kosong
  ```

Bandingkan konteks lain: identitas orang ter-hardcode di `NB FacIn` (43 berkas),
`RNW Fac In` (36), `Endorsment Fac In` (33), `NB Treaty In` (12), `Komite Claim FacIn` (4),
`Komite Claim Prop` (3) — dan **OQ-053**, di mana 5 identitas orang justru *ditetapkan* sebagai
pemilik tugas berikutnya.

**Koreksi artefak:** nilai `IT Developer` yang tercatat di D1 untuk `pyPosition` berasal dari sapuan
korpus-wide — **tidak ada di `Claim Life`**.

#### Consequences

- Otorisasi di sistem baru memakai **tiga peran ini sebagai himpunan kanonik**; penambahan peran
  adalah keputusan bisnis, bukan detail implementasi.
- Larangan rangkap peran harus **ditegakkan sistem**, bukan sekadar konvensi — dan pengecualiannya
  (“akses ditambahkan eksplisit di role akun”) perlu jalur yang terdefinisi.
- Pemetaan tahap → peran ini **tidak dapat digeneralisasi** ke konteks lain: facultative dan treaty
  inward memakai nama workbasket (`ReasFacIn*`, `ReasTreatyIn*`), bukan `pyPosition`.

#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-007** | Tidak ada rule identitas/otorisasi di korpus — model peran korpus-wide tetap harus dirancang |
| **OQ-021** (terjawab untuk Claim — Life) | Identitas orang ter-hardcode di 6 modul lain **tetap terbuka** |
| **OQ-024** (maju sebagian) | Pemetaan Assignment → nama workbasket di facultative & treaty tetap tidak terbaca |
| **OQ-028** | `WorkList` vs `WorkBasket` — Claim Life memakai **campuran**: Register/Outstanding `WorkList`, Medical Check/Claim Analis `WorkBasket`. Arti perbedaannya belum terverifikasi |

## ADR-U-0003 - Uang di Claim — Life tidak direpresentasikan sebagai `float`

Delapan kolom nilai pada tabel akseptasi klaim Life adalah **uang** dan **tidak** boleh
direpresentasikan sebagai `float` di sistem baru. Satu kolom adalah **persen**, bukan uang.
Representasi uang yang dipakai harus mempertahankan presisi desimal secara eksak
(mis. tipe desimal berskala tetap, atau bilangan bulat dalam satuan minor) — **bukan** biner
floating-point.

#### Klasifikasi kolom

`[terverifikasi work owner 2026-09-14]` Kolom pada `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`:

| Kolom | Sifat |
| --- | --- |
| `SUM_INSURED` | **uang** |
| `CEDING_RETENTION` | **uang** |
| `SUM_REASURED` | **uang** |
| `SHARE_NUSANTARA_RE` | **uang** |
| `CLAIM_AMOUNT` | **uang** |
| `SHARE_RETRO` | **uang** |
| `CLAIM_RETRO` | **uang** |
| `RETROCEDED_SHARE` | **uang** |
| `EM_PERCENT` | **persen** |
| `CURRENCY` | kode mata uang |

`[terverifikasi]` Nama kolom **tidak dapat dipakai untuk menebak sifatnya**. Empat kolom bernama
`SHARE_*` / `*_SHARE` ternyata **uang**, bukan rasio — padahal namanya menyarankan sebaliknya.
Korpus ini sudah terbukti memakai nama yang menipu (lihat `STS_REJECT` di `CONTEXT.md`).
Klasifikasi di atas berasal dari work owner, **bukan dari penamaan**.

#### Bukti sumber kolom

`[terverifikasi]` Tabel ditulis oleh
`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL`
(`Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml`) — blok PL/SQL `INSERT INTO … COMMIT`,
sehingga **55 nama kolomnya terbaca langsung**, bukan tersembunyi di stored procedure.

Perintah audit:
```
awk '/<pyBrowseSQL>/,/<\/pyBrowseSQL>/' "Komite Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml"
```

#### Consequences

- Aritmetika atas delapan kolom uang wajib memakai tipe eksak di seluruh lapisan —
  `handlers` → `services` → `repository` — dan di kontrak API.
- `EM_PERCENT` diperlakukan terpisah dari uang; mencampurnya dalam tipe yang sama akan
  menyembunyikan perbedaan makna.
- Claim — Life **sadar mata uang**: ada kolom `CURRENCY` pada rekam akseptasi. Ini berbeda dari
  nilai ter-hardcode di konteks lain yang **tidak menyebut mata uang** sama sekali
  (`Claim Non Prop`: `LimitMax = 30000000.00`; `Komite Claim FacIn`: pita
  `> 30000000.00 && <= 57750000.00`; `NB FacIn`: pangsa `0.45`/`0.05`, ambang `"3000000000"`
  dibandingkan sebagai **string**) → OQ-046.
- Nilai uang menyeberang ke Komite lewat kontrak **ADR-U-0001** (`CLAIM_AMOUNT` + `CURRENCY`
  ditambahkan atas keputusan Ronde 2 Q15) — representasinya harus konsisten di kedua sisi batas.


#### Penguatan dari DDL Oracle (2026-09-14) `[data DBA]`

OQ-001 ditutup untuk Claim — Life; DDL `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` kini diketahui.
Temuannya **mengunci** keputusan non-float menjadi syarat yang lebih tajam:

> **Kedelapan kolom uang bertipe Oracle `NUMBER` — tanpa presisi dan tanpa skala.**

`NUMBER` tanpa presisi berarti Oracle menyimpan desimal **presisi arbitrer** (hingga 38 digit
signifikan) apa adanya. Konsekuensinya mengikat:

- Go **wajib** memakai tipe desimal **presisi arbitrer**. **`float64` dilarang** — ia hanya punya
  ~15–17 digit signifikan, sehingga nilai yang sah di Oracle dapat **berubah diam-diam** saat
  dibaca.
- Larangan ini berlaku di **seluruh lapisan** dan di **kontrak API**, termasuk saat nilai menyeberang
  ke Komite (**ADR-U-0001**) dan saat data dipindahkan (**ADR-U-0009**).
- Karena kolomnya tanpa skala, **tidak ada pembulatan yang boleh diasumsikan**. Pembulatan apa pun
  harus keputusan eksplisit, bukan efek samping tipe.

Kolom penyerta: `CURRENCY VARCHAR2(100)`.

**Invariant mata uang** `[keputusan work owner]` — dari penutupan **OQ-060**: bentuk nilai uang
adalah `(amount, currency)` **per baris**, dengan syarat **semua baris satu klaim wajib bermata uang
sama**. `[terverifikasi]` Konsisten dengan `Claim Life/Activity/SetIndexAdjustmentList.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SETINDEXADJUSTMENTLIST` / `RULE-OBJ-ACTIVITY`), yang
menyalin `CURRENCY` dan `CURRENCYID` dari `AdjustmentList(1)` ke `AdjustmentList(<LAST>)`.

⚠️ `[data DBA]` Roster komite `POOLDATA.EMAILKOMITE` **tidak punya kolom mata uang**, sehingga pita
`LIMIT_BOTTOM`/`LIMIT_TOP` berlaku atas **satu mata uang implisit**. Selama invariant di atas
dipegang, pembandingan pita tetap sahih; bila kelak invariant itu dilonggarkan, pembandingan pita
menjadi tidak sahih.

#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-001** (modul lain) | DDL di luar persistensi Claim — Life belum diserahkan; tipe kolom uang konteks lain belum diketahui |
| **OQ-060** | **TERTUTUP untuk Claim — Life** (2026-09-14) — invariant satu klaim satu mata uang. Cakupan konteks lain tetap terbuka |

`[terjawab 2026-09-14]` OQ-060 kini **tertutup untuk Claim — Life**: bentuk tipe uang adalah
`(amount, currency)` **per baris**, dengan invariant **semua baris satu klaim bermata uang sama**.

## ADR-U-0004 - Alamat layanan keluar menjadi env var, bukan replikasi lookup `M_LINK_SERVICE`

> ⛔ **ADR INI SUDAH TIDAK BERLAKU.** Ia digantikan **ADR-U-0013**, yang memutuskan sebaliknya:
> alamat layanan keluar **wajib** di-resolve runtime dari `M_LINK_SERVICE`, dan **dilarang**
> ditanam sebagai env var. ADR-U-0004 ditulis saat isi `M_LINK_SERVICE` belum diketahui (OQ-047
> masih terbuka); `[data DBA]` isinya kini diketahui — 19 baris, dipakai bersama banyak modul.
> Teks di bawah dipertahankan untuk jejak audit.




Di Pega, alamat layanan keluar diambil **saat runtime dari tabel Oracle** `M_LINK_SERVICE`.
Di sistem baru, alamat menjadi **konfigurasi lingkungan (env var)**; mekanisme lookup itu
**tidak direplikasi**.

Keputusan ini terbatas pada **dari mana alamat berasal**. Apakah efek keluar dijalankan sama sekali
adalah keputusan terpisah — lihat **ADR-U-0005**.

#### Keadaan sekarang `[terverifikasi]`

Dua rule `RULE-CONNECT-REST` di `Claim Life`, keduanya `pyBaseURLSelectionType = SETTING` dengan
`pyBaseURLSetting = LinkService!LinkService` — **tanpa URL literal**:

| Berkas | `pyServiceName` |
| --- | --- |
| `Claim Life/ConnectREST/ServiceGoogle.xml` | `ServiceGoogle` |
| `Claim Life/ConnectREST/convertJsonNusareToProductionClaimLife.xml` | `convertJsonNusareToProductionClaimLife` |

Alamat sesungguhnya diresolusi lewat `Activity/GetLinkService.xml`
(`ASM-FW-GISFW-INT-M_LINK_SERVICE!GETLINKSERVICE` / `RULE-OBJ-ACTIVITY`), yang melakukan
`Obj-Browse` atas class `ASM-FW-GISFW-Int-M_LINK_SERVICE` dengan kunci **`KATEGORI_1`** dan
**`KATEGORI_2`** — pola yang sama di seluruh korpus (`discovery/understanding-report.md` §3.1).

**Isi tabel `M_LINK_SERVICE` tidak ada di korpus** → **OQ-047**. Artinya daftar endpoint yang
sesungguhnya **tidak diketahui**, dari korpus maupun dari ADR ini.

#### Considered Options

- **Env var per layanan** — dipilih
- Replikasi tabel `M_LINK_SERVICE` + lookup runtime — ditolak: memindahkan konfigurasi ke dalam
  data operasional membuat alamat tidak terbaca dari konfigurasi deployment, dan **isi tabelnya
  sendiri tidak diketahui** sehingga tidak dapat direplikasi dengan setia
- URL literal di kode — ditolak oleh aturan proyek (`CLAUDE.md` §10)

#### Consequences

- Daftar endpoint menjadi bagian konfigurasi deployment, terbaca dan dapat berbeda per lingkungan.
- **Migrasi memerlukan isi `M_LINK_SERVICE` terlebih dahulu** (OQ-047) — tanpa itu, tidak diketahui
  berapa endpoint yang harus disediakan, atau apa arti `KATEGORI_1`/`KATEGORI_2`.
- Aturan proyek tetap berlaku: endpoint & host internal **bukan literal** di kode, dan tidak masuk
  prompt, tiket, atau test.

#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-047** | Isi tabel `M_LINK_SERVICE` — daftar endpoint, arti `KATEGORI_1`/`KATEGORI_2`. Pemilik: DBA + Platform |
| **OQ-018** (terjawab untuk Claim — Life) | `jboss1073` = production, `jboss117` = dev; lingkungan `pega-nusre` **belum dinyatakan** |

## ADR-U-0005 - `IsPEGAPROD` tidak ditiru sebagai rule; diganti flag lingkungan eksplisit

Di Pega, tiga efek keluar Claim — Life digerbangi rule `When` bernama `IsPEGAPROD` yang
**kondisinya tidak terbaca dari korpus**. Di sistem baru, gerbang itu menjadi **flag lingkungan
eksplisit** (`ENV=production`) yang menggerbangi ketiga efek yang sama — sehingga perilakunya
terbaca, bukan tersembunyi di dalam rule.

Keputusan ini tentang **apakah efek keluar dijalankan**. Dari mana alamatnya berasal adalah
keputusan terpisah — lihat **ADR-U-0013** (alamat di-lookup runtime dari `M_LINK_SERVICE`; ini
menggantikan ADR-U-0004). Keduanya **sengaja tidak digabung**.

#### Keadaan sekarang `[terverifikasi]`

`Claim Life/When/IsPEGAPROD.xml` → `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN`.
`<pyLabel>`-nya hanya berisi template kosong `[first value][relation][second value]` — **apa yang
diujinya tidak diketahui** (**OQ-029**). Identitas ini juga **terdaftar berkonflik** di
`discovery/inventory/_oq011-konflik-isi.md` entri **#305** — isinya berbeda antar modul.

Tiga efek keluar yang digerbanginya di modul ini:

| Efek | Rule |
| --- | --- |
| Unggah berkas | `Claim Life/Activity/InsertGoogleStorage_Act.xml` |
| Simpan utama | `Claim Life/Activity/SaveOutStandingLife_Act.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT`, 642.787 byte) |
| Email | `Claim Life/Activity/SendEmailKlaimLF.xml` |

Perintah audit:
```
grep -rl "IsPEGAPROD" "Claim Life" --include="*.xml"
```

#### Considered Options

- **Flag lingkungan eksplisit** (`ENV=production`) — dipilih
- Migrasikan `IsPEGAPROD` apa adanya — **tidak mungkin**: kondisinya tidak terbaca, dan rule-nya
  berkonflik antar modul (OQ-011 #305). Menyalin salah satu varian berarti menebak
- Jalankan efek keluar tanpa gerbang — ditolak: menghapus perilaku yang jelas disengaja

#### Consequences

- Perilaku non-production menjadi **terbaca dari konfigurasi**, bukan dari isi rule.
- Ketiga efek keluar digerbangi **satu** flag yang sama, sebagaimana keadaan sekarang. Bila
  kemudian ternyata ketiganya perlu gerbang berbeda, itu perubahan berikutnya — bukan yang
  diputuskan di sini.
- `SaveOutStandingLife_Act` adalah **simpan utama**, bukan sekadar efek samping. Menggerbanginya
  dengan flag lingkungan berarti **di non-production data tidak tersimpan lewat jalur itu** —
  konsekuensi ini perlu disadari saat menyiapkan lingkungan uji.

#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-029** (terjawab untuk Claim — Life) | Kondisi asli `IsPEGAPROD` tetap tidak diketahui — keputusan ini **menghindari**, bukan menjawab. Varian di 13 modul lain tetap terbuka |
| **OQ-011** entri **#305** | Versi `IsPEGAPROD` mana yang berlaku di production |
| **OQ-018** (terjawab untuk Claim — Life) | `jboss1073` = production, `jboss117` = dev, sistem mirroring |

## ADR-U-0006 - Penomoran klaim Life memanggil `PROC_GENERATE_SEQUENCE_NUMBER`; dua SQL lama tidak dimigrasi

Nomor klaim Life **tidak dihitung ulang di aplikasi**. Sistem baru memanggil stored procedure
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` lewat jalur yang sama seperti sekarang. Dua rule penomoran
lama — `Generate_NoKlaim_Life` dan `Generate_NoKlaim_LifeRetro` — **sudah tidak dipakai
(di-remark)** dan **tidak dimigrasi**.

#### Bukti

`[terverifikasi work owner 2026-09-14]` Kedua SQL lama sudah di-remark; penomoran sekarang lewat
`GetSequenceNumber_SQL` → `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`.

`[terverifikasi]` **Dikuatkan bukti korpus** — asimetri indeks rujukan rule di
`Claim Life/Activity/SaveOutStandingLife_Act.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`, 642.787 byte):

| RequestType | muncul sebagai `<RequestType>` | muncul di `pxRuleReferences` (`<pyRuleName>`) |
| --- | ---: | ---: |
| `Generate_NoKlaim_Life` | 1 | **0** |
| `Generate_NoKlaim_LifeRetro` | 1 | **0** |
| `GetSequenceNumber_SQL` | 1 | **2** |

Keduanya tertinggal sebagai parameter langkah tetapi **tidak terindeks sebagai rujukan aktif**.

```
f="Claim Life/Activity/SaveOutStandingLife_Act.xml"
grep -c '<pyRuleName>.*Generate_NoKlaim_Life<' "$f"        # 0
grep -c '<pyRuleName>.*GetSequenceNumber_SQL<' "$f"        # 2
```

**Koreksi artefak FASE A:** `discovery/flows/Claim Life.md` §3 mendaftar kedua rule lama sebagai
dipanggil — itu berdasarkan kehadiran `<RequestType>` saja. Status sesungguhnya: **tidak aktif**.
Koreksi ini sudah direkam di register OQ-002.

#### Considered Options

- **Panggil `PROC_GENERATE_SEQUENCE_NUMBER`** — dipilih
- Replikasi logika penomoran di aplikasi — ditolak: **isi procedure tidak ada di korpus**
  (OQ-002), sehingga replikasi berarti menebak; dan penomoran yang bercabang di dua tempat berisiko
  menghasilkan nomor bentrok selama masa transisi
- Migrasikan dua SQL lama — ditolak: sudah tidak dipakai

#### Consequences

- Sistem baru **bergantung pada Oracle** untuk penomoran klaim — ketergantungan runtime yang harus
  disadari saat merencanakan lingkungan dan pengujian.
- Format nomor klaim **tidak dinyatakan** di aplikasi; ia ditentukan procedure.
- Bila kelak penomoran dipindah ke aplikasi, itu keputusan baru yang menggantikan ADR ini.

#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-002** (dipersempit untuk Claim — Life, 2026-09-14) | **Kontrak `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`**: format nomor yang dihasilkan, dan apakah sequence di-reset per tahun / per lini / global. **Pemilik: DBA.** 66 stored procedure lain di korpus tetap terbuka |
| **OQ-013** | Apakah procedure melakukan `COMMIT` sendiri — batas transaksi di sisi database |

Pembanding yang menunjukkan kontrak semacam ini **bisa** terbaca bila SQL-nya bukan procedure:
format nomor polis treaty terbaca penuh di `NB Treaty In/RDBList/GenerateNoPolicy.xml`
(`ASM-FW-GISFW-INT-POLISTREATYIN` / `ASM!GENERATENOPOLICY` / `RULE-CONNECT-SQL`).

## ADR-U-0007 - Jejak audit merekam siapa + kapan untuk setiap transisi status dan setiap jalur balik

Sistem baru merekam **siapa** dan **kapan** untuk **setiap** transisi status klaim **dan setiap
jalur balik** (`SendtoAdmin`, `SendtoMedical`). Ini **penyimpangan sadar dari sistem lama** —
sebuah perbaikan, bukan paritas.

#### Keadaan sekarang, dan apa yang hilang `[terverifikasi]`

Yang direkam Pega pada rekam akseptasi (`POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`, ditulis oleh
`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL`):

| Kolom | Isi |
| --- | --- |
| `CREATEOPNAME` | satu nama operator |
| `ACCEPTATION_DATE`, `CONFIRMATION_DATE`, `CLAIM_RECEIVED_DATE`, `COMPLETE_DATE` | empat tanggal |

Identitas pengguna juga dipakai sebagai **data jejak** di
`Claim Life/Activity/RejectOSClaimLife_Act.xml` dan `Claim Life/Activity/SendEmailKlaimLF.xml`
(`OperatorID.pyUserIdentifier`, `pyUserName`), dan ada pencatatan layanan di
`Claim Life/RDBList/InsertLogServiceClaim.xml` → `INSERT INTO pooldata.monitoring_klaim_log`.

**Yang tidak terekam sekarang:** siapa yang mengembalikan kasus lewat `SendtoAdmin` atau
`SendtoMedical`. Kedua penanda hanya menyimpan **nilai `1`**, tanpa pelaku dan tanpa waktu:

- `Claim Life/When/IsSendtoAdmin.xml` → `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOADMIN` /
  `RULE-OBJ-WHEN` — `pyWorkPage.SendtoAdmin = 1`
- `Claim Life/When/IsSendtoMedical.xml` → `…` / `ISSENDTOMEDICAL` / `RULE-OBJ-WHEN`

Jalur balik itu bukan kejadian langka: `SendtoAdmin` menggerbangi pengembalian dari **tiga titik**
di `Claim Life/Flow/Register_Flow.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` /
`RULE-OBJ-FLOW`).

#### Considered Options

- **Rekam siapa + kapan untuk setiap transisi dan setiap jalur balik** — dipilih
- Paritas dengan Pega (`CREATEOPNAME` + empat tanggal) — ditolak: pengembalian kasus menjadi tidak
  dapat ditelusuri, padahal ia jalur yang sering dipakai

#### Consequences

- Ini **penyimpangan sadar** dari non-goal "tidak merapikan alur" (Ronde 1 Q2). Alurnya tetap sama;
  yang bertambah adalah perekamannya.
- Jejak audit menjadi **riwayat transisi**, bukan sekadar beberapa kolom tanggal pada satu rekam.
  Bentuk penyimpanannya adalah keputusan implementasi, bukan ditetapkan di sini.
- Peran pelaku terbatas pada tiga peran di **ADR-U-0002**; pencatatan pelaku memakai identitas akun,
  bukan nama ter-hardcode.
- Rekam akseptasi lama tetap ditulis sebagaimana adanya — kontrak dengan Komite (**ADR-U-0001**)
  tidak berubah karena keputusan ini.

#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-001** | Tidak ada DDL — struktur `monitoring_klaim_log` dan tipe kolom tanggal tidak diketahui |
| **OQ-013** | `COMMIT` di dalam blok PL/SQL — batas transaksi ada di sisi database, sehingga atomisitas "tulis akseptasi + tulis jejak audit" perlu ditetapkan bersama DBA |

## ADR-U-0008 - Efek keluar Claim — Life dijalankan asinkron, tidak memblokir alur, dengan antre-ulang

Keempat efek keluar Claim — Life bersifat **asinkron**: kegagalan **dicatat**, **tidak memblokir**
alur klaim, dan disediakan **antre-ulang (retry)**. Perilaku "tidak memblokir" dipertahankan dari
sistem lama; **keandalannya ditambah**.

#### Efek keluar yang tercakup `[terverifikasi]`

| Efek | Rule |
| --- | --- |
| Unggah berkas ke Google Storage | `Claim Life/Activity/InsertGoogleStorage_Act.xml` |
| Email | `Claim Life/Activity/SendEmailKlaimLF.xml` |
| Arasapas | `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` — **satu-satunya salinan di korpus** (OQ-035) |
| Konversi ke produksi | `Claim Life/ConnectREST/convertJsonNusareToProductionClaimLife.xml` (`pyServiceName` = `convertJsonNusareToProductionClaimLife`, `pyBaseURLSelectionType = SETTING`) |

#### Mengapa "tidak memblokir" adalah paritas, bukan perubahan

`[terverifikasi]` Di `Claim Life` ditemukan **pencatatan**, bukan gerbang keberhasilan:
`Claim Life/RDBList/InsertLogServiceClaim.xml` → `INSERT INTO pooldata.monitoring_klaim_log`.

Bandingkan konteks facultative, yang **punya** gerbang keberhasilan eksplisit —
`When/IsSuccessHitService.xml` menguji `.StatusService.StsKonversiFacIn` / `StsKonversiFacOut`
dan mengarahkan alur ke `SetToInbox_ACT` bila gagal. **Pola itu tidak ada di `Claim Life`.**

Jadi keputusan "tidak memblokir" **mempertahankan perilaku terbaca**, bukan mengarangnya.
Yang ditambahkan adalah **antre-ulang**, yang di Pega tidak terbaca ada.

#### Consequences

- Alur klaim (Register → Outstanding → Medical Check → Claim Analis) **tidak pernah tertahan** oleh
  kegagalan layanan luar.
- Kegagalan harus **terlihat** — jejaknya masuk ke jalur audit **ADR-U-0007**, tidak cukup hanya ke
  log layanan.
- Efek keluar digerbangi flag lingkungan **ADR-U-0005**; di non-production ketiganya tidak berjalan,
  sehingga antre-ulang pun tidak aktif di sana.
- Alamat layanan **di-lookup runtime** dari `M_LINK_SERVICE` (**ADR-U-0013**, menggantikan ADR-U-0004);
  kegagalan **kunci kategori tidak ditemukan** harus dibedakan dari kegagalan jaringan agar
  antre-ulang tidak berputar sia-sia.
- **Konsekuensi yang perlu disadari:** karena tidak memblokir, sebuah klaim dapat mencapai status
  akhir sementara efek keluarnya masih tertunda. Apakah keadaan itu boleh dianggap "selesai"
  adalah pertanyaan bisnis yang **belum dijawab**.

#### Pengecualian: Komite Claim Life (2026-09-15)

⚠️ ADR ini berlaku untuk **Claim — Life**. Konteks **Komite Claim Life** memutuskan sebaliknya —
efek keluarnya **wajib berhasil** (transactional outbox, at-least-once) — karena ia memuat integrasi
**Kasir** (pembayaran) yang tidak ada di Claim Life. Lihat **ADR-U-0015**.

Keduanya berlaku bersamaan pada konteks masing-masing; ini **bukan** pembatalan.

#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-025 / OQ-035** | `serviceInsertArasapasClaimLife_act` adalah satu-satunya salinan di korpus dan dipanggil juga oleh `Komite Claim Life` — kepemilikan dan kontraknya belum ditetapkan |
| **OQ-047** | Isi tabel `M_LINK_SERVICE` — daftar endpoint yang sebenarnya tidak diketahui |
| **OQ-002** | `POOLDATA.GET_TOKEN_STORAGE` (dipakai jalur unggah berkas) tanpa body di korpus |
| **OQ-029** | Kondisi asli `IsPEGAPROD` tidak terbaca — gerbang lingkungan menghindarinya, bukan menjawabnya |

## ADR-U-0009 - Seluruh data Claim — Life dipindahkan; tidak ada koeksistensi dua penulis

**Semua data dipindah** ke sistem baru. Opsi koeksistensi — sistem baru hanya menerima klaim baru
sementara klaim berjalan diselesaikan di Pega — **ditolak**.

#### Considered Options

- **Migrasi penuh seluruh data** — dipilih
- **Koeksistensi**: hanya klaim baru di sistem baru; klaim berjalan (Outstanding, Medical Check,
  Claim Analis) diselesaikan di Pega — **ditolak**

Koeksistensi sempat diusulkan karena secara teknis mungkin: kedua sistem menulis ke tabel yang sama,
`POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`, lewat rule yang identik
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL`, hash ternormalisasi
`c50bfd9a12`, `[terverifikasi]` sama persis di `Claim Life` dan `Komite Claim Life`).
Penolakannya dicatat di sini supaya tidak diusulkan ulang.

#### Consequences

- **Tidak ada periode dua penulis** ke `OS_AKSEPTASI_KLAIM_LIFE` dari sisi Claim Life — cutover
  bersifat memutus, bukan bertahap.
- Klaim yang sedang berada di tengah siklus **harus terbawa beserta statusnya**. Statusnya adalah
  `STS_REJECT` (`0` = Outstanding) ditambah posisi tahap (Register / Outstanding / Medical Check /
  Claim Analis) — lihat `CONTEXT.md` §Status dan kode. Kasus yang sudah diserahkan ke Komite dan
  belum kembali membawa keadaan lintas batas (**ADR-U-0001**), sehingga pemindahannya perlu
  disepakati dengan pemilik konteks Komite Life.
- Penomoran klaim tetap memanggil `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` (**ADR-U-0006**).
  **Migrasi tidak boleh membuat sequence melompat atau mengulang** — perilaku procedure terhadap
  data yang dipindahkan perlu dipastikan bersama DBA.
- Jejak audit baru (**ADR-U-0007**) **tidak dapat direkonstruksi ke belakang**: data lama hanya
  memiliki `CREATEOPNAME` dan empat kolom tanggal. Riwayat transisi lengkap hanya ada untuk kejadian
  setelah cutover.

#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-001** | Tidak ada DDL — struktur dan tipe kolom sumber tidak diketahui, sehingga rencana pemindahan belum dapat disusun rinci |
| **OQ-002** (dipersempit) | Kontrak `PROC_GENERATE_SEQUENCE_NUMBER` — apakah sequence di-reset per tahun/lini/global menentukan aman-tidaknya migrasi nomor |
| **OQ-018** (terjawab untuk Claim — Life) | `jboss1073` = production, `jboss117` = dev, sistem mirroring. Lingkungan `pega-nusre` (73 berkas di modul ini) **belum dinyatakan** — perlu dipastikan sebelum menentukan sumber data migrasi |
| **OQ-060** | Cakupan `CURRENCY` — memengaruhi bentuk data uang yang dipindahkan |

## ADR-U-0010 - Penyimpanan berkas klaim tetap di Google Storage

Berkas klaim Life **tetap disimpan di Google Storage**. Migrasi tidak memindahkan penyimpanan
berkas; yang berganti hanya aplikasi yang memanggilnya.

#### Jalur yang ada sekarang `[terverifikasi]`

Delapan berkas `Claim Life` menyentuh penyimpanan:

| Rule | Peran |
| --- | --- |
| `Claim Life/Activity/InsertGoogleStorage_Act.xml` (`ASM-FW-GISFW-INT-T_STORAGE_IMAGE!INSERTGOOGLESTORAGE_ACT`, 160.027 byte) | unggah |
| `Claim Life/Activity/GetUrlGoogleStorage_Act.xml` | ambil URL |
| `Claim Life/Activity/DeleteGoogleStorage_Act.xml` | hapus |
| `Claim Life/Activity/InsertDocument_Act.xml`, `DeleteDocument_Act.xml`, `LoadDocumentLife_ACT.xml`, `DownloadDocumentClaim.xml` | pembungkus tingkat dokumen |
| `Claim Life/RDBList/GetTokenStorage_SQL.xml` (`ASM-FW-GISFW-INT-T_STORAGE_IMAGE!RNM!GETTOKENSTORAGE_SQL`) | **penerbitan token** |

Token **tidak** diterbitkan aplikasi. Ia datang dari Oracle:

```sql
BEGIN
  pooldata.GET_TOKEN_STORAGE ( {UploadDoc.App}, {OperatorID.pyUserIdentifier},
                               {UploadDoc.Kodestring OUT}, {DocAPI.ResponseMsg OUT});
  COMMIT;
END;
```

`[terverifikasi]` Dua hal terbaca langsung dari tanda tangan itu:

1. Procedure menerima **`OperatorID.pyUserIdentifier`** — identitas pengguna ikut masuk, sehingga
   wewenang atas berkas **mungkin** ditentukan per pengguna. Tanpa body procedure hal ini tidak
   dapat dipastikan (**OQ-002**).
2. Procedure menerima **`UploadDoc.App`** — token tampaknya dicakup per aplikasi, bukan global.

`[terverifikasi]` Ini **bukan milik Claim — Life**: 97 berkas di **15 modul** memakai jalur
penyimpanan yang sama, dan kelasnya `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` berada di luar kelas modul.
Penyimpanan berkas adalah **infrastruktur bersama korpus**.

#### Considered Options

- **Tetap di Google Storage** — dipilih
- Pindah ke penyimpanan lain (S3-compatible, on-prem) — ditolak: berkas yang sudah ada harus ikut
  pindah, dan penyimpanan dipakai bersama 15 modul yang **tidak** ikut dimigrasi pada gelombang ini

#### Consequences

- Sistem baru harus dapat memanggil `POOLDATA.GET_TOKEN_STORAGE` — **ketergantungan Oracle kedua**
  di luar penomoran klaim (**ADR-U-0006**). Keduanya bukan sekadar penyimpanan data; keduanya
  **logika** yang berada di database.
- Kontrak procedure harus diperoleh dari **DBA** sebelum unggah/unduh dapat dispesifikasikan:
  masa berlaku token, cakupan (`App`), dan apakah `pyUserIdentifier` menentukan wewenang.
- Unggah berkas termasuk **efek keluar asinkron** (**ADR-U-0008**) — kegagalannya dicatat, tidak
  memblokir alur klaim, dan diantre ulang.
- Karena penyimpanan dipakai 15 modul yang masih di Pega, selama masa transisi **dua aplikasi**
  membaca dan menulis berkas yang sama. Ini berbeda dari data klaim, yang dipindah penuh tanpa
  koeksistensi (**ADR-U-0009**).
- Alamat layanan penyimpanan **di-lookup runtime** dari `M_LINK_SERVICE` (**ADR-U-0013**, menggantikan
  ADR-U-0004); kredensial dan koneksi database tetap env var. Gerbang lingkungan **ADR-U-0005**.
#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-002** | Body `POOLDATA.GET_TOKEN_STORAGE` tidak ada di korpus — kontrak token tidak terbaca. **Pemilik: DBA** |
| **OQ-018** | Apakah bucket/proyek berbeda antara production dan dev belum dinyatakan |
| **OQ-047** | Isi `M_LINK_SERVICE` — daftar endpoint yang sebenarnya tidak diketahui |

## ADR-U-0011 - Unit status Claim — Life adalah **baris `AdjustmentList`**, bukan klaim

Status klaim Life diputuskan **per baris `AdjustmentList`**. `STS_REJECT` pada
`PremiumListDetail` dan pada header klaim adalah **cerminan** baris adjustment terakhir, bukan unit
keputusan tersendiri. Revisi keputusan dilakukan dengan **menambah baris baru**, bukan mengubah
baris lama.

Ini menggantikan pembacaan awal saya bahwa Claim — Life memiliki satu status di tingkat klaim.

#### Mesin status `[keputusan work owner 2026-09-14]`

| # | Langkah | Peran | Akibat pada `STS_REJECT` baris |
| --- | --- | --- | --- |
| 1 | Input baris `AdjustmentList` pertama, lalu **Save ke OS** | `ReasLifeAdmin` | `0` |
| 1b | **Reject Outstanding** atas adjustment yang ia input sendiri — tanpa Komite. **Membatalkan baris itu saja**; Admin lalu input baris baru. | `ReasLifeAdmin` | `2` |
| 2 | Submit ke Medical Check | `ReasLifeAdmin` → `ReasLifeMedicalAdvisor` | tetap `0` |
| 3 | Submit ke SPV | `ReasLifeMedicalAdvisor` → `ReasLifeSPV` | tetap `0` |
| 4 | **Send ke Komite**, untuk semua adjustment | `ReasLifeSPV` — **kecuali `Type` `TP`/`TR`**, lihat ralat di bawah | tetap `0` |
| 5 | Komite memutus | Komite Life (konteks luar) | aksep → `1`; tolak → `2` |
| 6 | Setelah tolak: **tambah baris `AdjustmentList` baru** → Save ke OS → send Komite. Berulang. | `ReasLifeSPV` | baris baru mulai dari `0` |

**Aksep final selalu lewat Komite.**

##### Ralat langkah 4 — `Type` `TP`/`TR` bebas peran

`[keputusan work owner 2026-09-14, Ronde 3 lanjutan]` Aturan "send ke Komite hanya SPV" berlaku
**hanya untuk tipe selain `TP` dan `TR`**. Untuk klaim ber-`Type` **`TP`** atau **`TR`**, pengiriman
ke Komite **tidak dibatasi peran** — **`ReasLifeAdmin` pun dapat mengirim langsung**.

Ini **meralat** kalimat "Admin tidak mengirim ke Komite" pada versi pertama ADR ini.

`[terverifikasi]` Gerbang visibilitasnya di `Claim Life/Section/AdjustmentDetail_Section.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE!ADJUSTMENTDETAIL_SECTION`):

```
pyWorkPage.pyPosition =='ReasLifeSPV' || pyWorkPage.Type = 'TP' || pyWorkPage.Type = 'TR'
```

Seluruh operatornya `||`, tanpa `&&` — jadi **tidak ada ambiguitas presedensi**: kedua nilai `Type`
itu memang meloloskan gerbang tanpa memeriksa peran.

`[keputusan work owner 2026-09-14]` **`TP` = Payable, `TR` = Receivable.** Definisi ini **tidak ada
di korpus**; sumbernya work owner.

**Perilaku ini dibawa apa adanya (paritas), tidak diperketat** — keputusan itu beserta risikonya
dicatat terpisah di **ADR-U-0012**, karena ia menyangkut wewenang, bukan mesin status.

##### Kefinalan

- **Baris terminal.** Sekali sebuah baris bernilai `1` atau `2`, nilainya tidak berubah lagi.
- **Klaim tidak terminal.** Setelah penolakan, baris baru dibuat; siklus berulang — oleh SPV bila
  penolaknya Komite, oleh Admin bila ia sendiri yang menolak.
- **Revisi = baris baru**, bukan pengubahan baris lama.

##### Arti tunggal nilai `2` `[keputusan work owner 2026-09-14]`

**`STS_REJECT = 2` SELALU berarti "baris ini ditolak" — TIDAK PERNAH "klaim selesai".**

Kedua sumber penolakan menulis nilai yang **sama** dengan **makna operasional setara pada tingkat
baris**:

| Sumber | Rule penulis | Akibat pada klaim |
| --- | --- | --- |
| **Admin** — "Reject Outstanding", tanpa Komite | `Claim Life/Activity/RejectOSClaimLife_Act.xml` | klaim **tidak** tertutup; Admin input baris baru |
| **Komite** — lewat SPV | `Komite Claim Life/Activity/KomitePostAdjustment.xml` | klaim **tidak** tertutup; SPV input baris baru |

Konsekuensi untuk spesifikasi: **tidak ada nilai `STS_REJECT` yang menandakan klaim selesai.**
Selesainya sebuah klaim bukan status tersimpan — ia keadaan turunan dari kumpulan barisnya.
`[pertanyaan terbuka]` **Aturan turunannya belum ditetapkan** dan bukan fakta yang hilang dari
korpus, melainkan keputusan desain — dicatat di `.scratch/claim-life/kesiapan-to-spec.md`, bukan
sebagai OQ.

`[terverifikasi]` Korpus memang tidak membedakan keduanya: kedua rule menulis `2` ke tingkat baris
**dan** ke `PremiumListDetail(idx)` dalam bentuk yang sama persis. Perbedaannya memang tidak ada —
sesuai keputusan di atas, karena memang tidak seharusnya ada.

#### Bukti korpus yang menguatkan `[terverifikasi]`

##### 1. Gerbang peran memang ada, dan bekerja pada `.STS_REJECT` tingkat baris

`Claim Life/Section/AdjustmentDetail_Section.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE!ADJUSTMENTDETAIL_SECTION`, 557.563 byte):

| Kontrol | `<pyCondition>` |
| --- | --- |
| **"Reject Outstanding"** | `pyWorkPage.pyPosition =='ReasLifeAdmin' && pyWorkPage.ClaimData.PremiumListSummary.CLAIM_NO !='' && .STS_REJECT == 0` |
| Jalur Komite (`GetListKomiteLife`) | `pyWorkPage.pyPosition =='ReasLifeSPV' \|\| pyWorkPage.Type = 'TP' \|\| pyWorkPage.Type = 'TR'` |
| **"Save to Outstanding"** (`SaveOutstandingLife_Act`) | `pyWorkPage.pyPosition =='ReasLifeSPV'` |

Gerbang Reject berbunyi **`.STS_REJECT == 0`** — relatif, yaitu **tingkat baris**. Inilah bukti
terkuat bahwa unit keputusan adalah baris: wewenang diuji terhadap status baris, bukan status klaim.

Perintah audit:
```
grep -n "pyPosition" "Claim Life/Section/AdjustmentDetail_Section.xml"
```

##### 2. `RejectOSClaimLife_Act` **bukan** dead rule

Dirujuk 2× sebagai `<pyActivity>` dari `Claim Life/Section/RejectOSClaimLife_Sec.xml`, dan
kontrol "Reject Outstanding" di atas digerbangi peran `ReasLifeAdmin`. Ia menulis:

```
.STS_REJECT                                                    = 2
pyWorkPage.ClaimData.PremiumListSummary.PremiumListDetail(idx).STS_REJECT = 2
```

**Koreksi:** pada laporan sebelumnya saya menyebut adanya "dua jalur aksep/tolak" tanpa dapat
menjelaskan yang mana milik siapa. Sekarang jelas — jalur reject langsung adalah **milik Admin**.

##### 3. Pencerminan ke `PremiumListDetail` ditulis berbarengan, bukan disalin belakangan

Setiap rule yang memutus menulis **dua tingkat sekaligus** dengan nilai sama:

| Rule | Tingkat baris adjustment | Tingkat `PremiumListDetail` |
| --- | --- | --- |
| `Claim Life/Activity/RejectOSClaimLife_Act.xml` | `.STS_REJECT = 2` | `= 2` |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `1` (6×) / `2` (2×) | nilai sama, pasangan langsung |

Jadi "cerminan" bukan tafsir — ia terbaca sebagai penulisan berpasangan.

`[keputusan work owner]` Bahwa yang tercermin adalah baris **terakhir** berasal dari work owner;
korpus hanya menunjukkan yang tercermin adalah **baris yang sedang diputus**.

##### 4. Baris baru diwarisi dari baris pertama

`Claim Life/Activity/SetIndexAdjustmentList.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!SETINDEXADJUSTMENTLIST`, 47.748 byte) menyalin **delapan**
kolom dari `AdjustmentList(1)` ke `AdjustmentList(<LAST>)`:

```
SHARE_NUSANTARA_RE  CEDING_RETENTION  SUM_REASURED  SUM_INSURED
SHARE_RETRO         RETROCEDED_SHARE  CURRENCYID    CURRENCY
```

`STS_REJECT` **tidak** termasuk — baris baru tidak mewarisi status. Ini persis yang dibutuhkan
langkah 6 (baris baru mulai dari `0`).

##### 5. Tidak ada jalur buka-ulang pada baris yang sama

Sensus penulis lengkap (register **OQ-061**): **tidak ada satu pun rule** di `Claim Life` maupun
`Komite Claim Life` yang menulis `0` setelah `1` atau `2`. Satu-satunya penulis nilai `0` adalah
`SaveOutStandingLife_Act`, yang dipakai saat baris **baru** disimpan.

#### Considered Options

- **Unit status = baris `AdjustmentList`; header mencerminkan baris terakhir** — dipilih
- Satu status di tingkat klaim, baris hanya rincian — ditolak: gerbang wewenang di korpus menguji
  `.STS_REJECT` **tingkat baris**, dan penolakan Komite menghasilkan baris baru, bukan perubahan
  status klaim
- Status ganda (klaim dan baris masing-masing otoritatif) — ditolak: menghasilkan dua sumber
  kebenaran untuk pertanyaan yang sama

#### Consequences

- **Model data**: entitas keputusan adalah baris adjustment. Status klaim adalah **turunan**
  (*derived*), bukan kolom yang ditulis sendiri — meskipun sistem lama menyimpannya sebagai kolom.
- **Riwayat**: menolak lalu mengajukan ulang menghasilkan **beberapa baris** pada satu klaim.
  Jumlah baris = jumlah putaran Komite. Ini yang membuat klaim tidak terminal.
- **Jejak audit** (**ADR-U-0007**) direkam **per baris**, bukan per klaim — transisi yang perlu
  siapa+kapan adalah transisi baris.
- **Kontrak Komite** (**ADR-U-0001** Kontrak 2) memang bekerja pada tingkat baris; `IndexAdjustment`
  dan `IndexPremiumList` yang menyeberang kini masuk akal sebagai penunjuk baris.
- **Migrasi** (**ADR-U-0009**) harus membawa **seluruh baris**, bukan hanya keadaan terakhir —
  kalau tidak, riwayat putaran Komite hilang.
- **Uang** (**ADR-U-0003**) melekat pada baris; `CURRENCY` dan `CURRENCYID` pun ada di tingkat baris
  dan disalin dari baris pertama.

#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-060** | Apakah keseragaman mata uang per klaim adalah **aturan** atau kebetulan implementasi — `CURRENCY` ada di tingkat baris dan disalin dari baris 1 |
| **OQ-032** | Apa yang menentukan `KomiteLoop` — berapa putaran sebelum keputusan final |
| **OQ-001** | Tidak ada DDL — bentuk penyimpanan baris adjustment di Oracle tidak diketahui |

**Sudah tertutup untuk Claim — Life** (2026-09-14, work owner): **OQ-039** (akibat reject Admin),
**OQ-061** (unit keputusan), **OQ-062** (kefinalan), **OQ-063** (gerbang `TP`/`TR` — lihat
**ADR-U-0012**).

#### `SaveAdjustment_Act` — dead rule, tidak dimigrasikan

`[keputusan work owner 2026-09-14]` `Claim Life/Activity/SaveAdjustment_Act.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!SAVEADJUSTMENT_ACT` / `RULE-OBJ-ACTIVITY`, 179.221 byte)
**sudah tidak dipakai**. **Jangan dimigrasikan.** Penulis `STS_REJECT = 1` yang berlaku adalah rule
sisi Komite, `KomitePostAdjustment`.

⚠️ **Status "dead" ini TIDAK dapat diverifikasi dari korpus — ia keputusan work owner.** Yang
terbaca justru sebaliknya:

`[terverifikasi]` Rule itu **terpasang di UI** — dirujuk 2× sebagai `<pyActivity>` dari
`Claim Life/Section/ClaimLifeDetailGCNM.xml` (824.562 byte) — dan menulis
`.AdjustmentList(<LAST>).STS_REJECT = 1` bersama `ACCEPTEDNO` dan `ACCEPTATION_DATE`, **tanpa
precondition Komite**. Ia juga hanya ada di `Claim Life`, tidak di modul lain.

Perbedaan ini dicatat dengan sengaja: bila kelak ditemukan bahwa jalur itu masih dipakai di
produksi, keputusan "jangan dimigrasikan" perlu ditinjau ulang, dan bukti di atas adalah titik
mulanya.

Perintah audit:
```
grep -rn "<pyActivity>SaveAdjustment_Act</pyActivity>" "Claim Life" --include="*.xml"
grep -n "AdjustmentList(&lt;LAST&gt;).STS_REJECT" "Claim Life/Activity/SaveAdjustment_Act.xml"
```

## ADR-U-0012 - Wewenang kirim ke Komite bergantung `Type` — dibawa apa adanya, dengan risiko RBAC tercatat

Di Claim — Life, **siapa yang boleh mengirim kasus ke Komite ditentukan oleh `Type` klaim**, bukan
oleh peran saja:

| `Type` | Siapa yang boleh mengirim |
| --- | --- |
| **`TP`** (Payable) atau **`TR`** (Receivable) | **siapa pun** — `ReasLifeAdmin` termasuk |
| lainnya (`QP`, `QR`) | **hanya `ReasLifeSPV`** |

**Perilaku ini dibawa apa adanya ke sistem baru (paritas). Tidak diperketat pada migrasi ini.**
Risikonya dicatat di bawah untuk ditinjau saat konteks **Komite Life** / **IAM** digarap.

#### Arti kode `[keputusan work owner 2026-09-14]`

**`TP` = Payable. `TR` = Receivable.**

⚠️ Definisi ini **tidak ada di korpus** — sumbernya work owner. `[terverifikasi]` Pencarian korpus
sudah tuntas dan nihil: tidak ada label, caption, atau opsi mana pun yang memasangkan kode itu
dengan teks; satu-satunya tempat nilainya **ditetapkan** adalah
`PremiumList Life/Activity/SubmitPremiumList_Act.xml` (`ASM-FW-GISFW-WORK-LIFE!SUBMITPREMIUMLIST_ACT`,
214.154 byte), dan di sana keempat kode hanya diteruskan ke slot generik `InputData.CARI20` dengan
precondition `.Type=="<kode itu sendiri>"` — melingkar, tidak mendefinisikan.

`QP` dan `QR` **belum** dijawab dan tetap terbuka di **OQ-020**.

#### Bukti: `Type` menggerbangi dua hal `[terverifikasi]`

##### 1. Wewenang kirim ke Komite

`Claim Life/Section/AdjustmentDetail_Section.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE!ADJUSTMENTDETAIL_SECTION`, 557.563 byte), kontrol jalur Komite
(`GetListKomiteLife`):

```
pyWorkPage.pyPosition =='ReasLifeSPV' || pyWorkPage.Type = 'TP' || pyWorkPage.Type = 'TR'
```

Seluruhnya `||`, tanpa `&&` — tidak ada ambiguitas presedensi. Dan `=` di korpus ini adalah
**pembanding**, bukan assignment: pada `<pyCondition>` bentuk `=` justru mayoritas (2.282 vs 954
`==`), dan banyak ekspresi mencampur keduanya dalam satu baris di tempat yang akan rusak kasatmata
kalau `=` berarti assignment — mis. `.ACCEPTEDNO=="" && .IsCheck = true && .STS_REJECT == "0"`.

`[terverifikasi]` Gerbang seperti ini **hanya ada di `Claim Life`** — tidak ada di 19 modul lain.

##### 2. Jendela validasi Date of Loss

`Claim Life/Activity/ValidasiDOL_Act.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!VALIDASIDOL_ACT`, 59.747 byte):

| Cabang | Jendela yang dipakai | Pergeseran tanggal |
| --- | --- | --- |
| `Type=="QR" \|\| Type=="QP"` | `GROSS_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `@addCalendar(.DATE_OF_LOSS,0,0,0,0,0,0,0)` — **nol** |
| `Type=="TR" \|\| Type=="TP"` | `RETROCESSION_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `@addCalendar(.DATE_OF_LOSS,0,0,0,1,0,0,0)` — **+1 hari** |

Gagal → `local.errmsg = "Invalid DOL"`, dengan `local.Begin==false || local.Expired==true`.

##### 3. Kedua gerbang membaca **dua salinan** properti yang sama

`[terverifikasi]` Gerbang Komite membaca `pyWorkPage.Type`; `ValidasiDOL_Act` membaca
`pyWorkPage.PolicyDataLife.Type`. Keduanya berisi nilai yang sama karena disalin:

```
Claim Life/Activity/LoadDataPeserta_Act.xml
  pyWorkPage.Type = pyWorkPage.PolicyDataLife.Type
```

Sumber otoritatifnya adalah **`PolicyDataLife.Type`** — tipe polis; `pyWorkPage.Type` hanya salinan
kerja. Di sistem baru keduanya menjadi **satu** field; yang perlu disadari adalah di Pega mereka
dua, sehingga secara teori dapat berbeda bila `LoadDataPeserta_Act` tidak berjalan.

#### Considered Options

- **Paritas — bawa apa adanya, catat risikonya** — dipilih
- Perketat menjadi `SPV && (TP || TR)` — ditolak pada migrasi ini: itu **mengubah perilaku**, dan
  wewenang lintas konteks (Claim Life ↔ Komite Life ↔ IAM) belum ditetapkan menyeluruh.
  **OQ-021** menunjukkan RBAC adalah blok terbesar konteks Komite

Membiarkan pilihan ini tidak tercatat bukan opsi: seorang pembaca di kemudian hari akan menyangka
ini kekeliruan porting, lalu "memperbaikinya" tanpa tahu bahwa ia disengaja.

#### Risiko yang diterima

⚠️ **Wewenang bergantung pada atribut data, bukan pada peran.** Siapa pun yang dapat membuat klaim
ber-`Type` `TP`/`TR` dengan sendirinya memperoleh wewenang mengirim ke Komite. Peran bukan lagi
satu-satunya penentu akses.

⚠️ **`Type` adalah arah akuntansi (Payable/Receivable), bukan tingkat kewenangan.** Tidak terbaca —
di korpus maupun dari definisinya — mengapa arah akuntansi menentukan siapa boleh mengeskalasi.
Kemungkinannya: kelonggaran operasional lama, atau `Type` membawa makna lain yang belum dinyatakan.
**Tidak ditebak.**

⚠️ **Ditinjau ulang saat Komite Life / IAM digarap**, bukan sekarang. Bila peninjauan itu
memutuskan memperketat, ADR ini yang diganti.

#### Consequences

- Lapisan layanan harus menegakkan aturan bergantung-`Type` ini **secara eksplisit** — bukan
  mewarisinya sebagai efek samping visibilitas UI, sebagaimana di Pega. Penegakan di lapisan layanan
  mengikuti **ADR-U-0002**.
- `Type` menjadi **field yang menyentuh keamanan**, bukan sekadar data. Perubahan nilainya harus
  masuk jejak audit (**ADR-U-0007**).
- Mesin status (**ADR-U-0011**) tidak berubah karenanya — yang berbeda hanya **siapa** yang boleh
  melakukan langkah 4.
- Validasi DOL mewarisi percabangan yang sama; keduanya harus memakai **sumber `Type` yang sama**
  agar tidak dapat berbeda seperti di Pega.

#### Tangga komite sepenuhnya data-driven (2026-09-14) `[terverifikasi]`

Penutupan **OQ-032** melengkapi ADR ini dari sisi lain: `Type` menentukan **siapa** yang boleh
menyerahkan; **roster** menentukan **berapa tingkat** tangga persetujuannya.

> **`KomiteLoop` = COUNT baris roster `EMAILKOMITE` yang aktif (`STS_AKTIF = "1"`) dan
> ber-`LIMIT_BOTTOM <= CLAIM_AMOUNT`.** Tidak ada konstanta di mana pun.

Tiga bukti berantai:

1. `Claim Life/Activity/GetListKomiteLife.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` /
   `GETLISTKOMITELIFE` / `RULE-OBJ-ACTIVITY`) — `Local.IsADj = .CLAIM_AMOUNT`, lalu memanggil report
   `FilterEmailKomiteWithLimit` (`Param.pyReportClass = "ASM-FW-GCNMFW-Int-EMAILKOMITE"`), hasilnya
   ke `Primary.KomiteList(<APPEND>)`.
2. `Claim Life/ReportDefinition/FilterEmailKomiteWithLimit.xml` (`ASM-FW-GCNMFW-INT-EMAILKOMITE` /
   `FILTEREMAILKOMITEWITHLIMIT` / `RULE-OBJ-REPORT-DEFINITION`) — filter
   `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM AND .STS_KLAIM = Param.STS_KLAIM AND .STS_AKTIF = "1"`.
3. `Claim Life/Activity/CreateKMTLife_Act.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` /
   `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`) —
   `childPageKomite.KomiteLoop = @Utilities.SizeOfPropertyList(childPageKomite.KomiteList)`.

⚠️ **Detail yang wajib direplikasi:** ambang yang dikirim ke roster adalah **nilai mutlak** klaim —
`Param.LIMIT_BOTTOM = @if(Local.IsADj<0, Local.IsADj * -1, Local.IsADj)`. Klaim bernilai negatif
dicari dengan tandanya dihilangkan.

**Akibat:** menambah atau menonaktifkan satu baris roster **mengubah jumlah tingkat persetujuan**
klaim yang sedang berjalan. Roster adalah **data operasional yang menyentuh alur keputusan**, bukan
sekadar daftar alamat email — perubahannya layak masuk jejak audit (**ADR-U-0007**).

⚠️ `[data DBA]` Roster **tidak punya kolom mata uang**; pita nilainya berlaku atas satu mata uang
implisit. Sahih selama invariant **ADR-U-0003** (satu klaim satu mata uang) dipegang.

#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-020** | Arti `QP` dan `QR` belum dijawab. `[dugaan]` Bukti §2 memperlihatkan huruf **pertama** memisahkan gross (`Q*`) dari retrocession (`T*`); bila huruf **kedua** memang Payable/Receivable seperti pada `TP`/`TR`, maka `QP`/`QR` mengikuti pola yang sama — **belum dikonfirmasi, jangan dipakai sebagai fakta** |
| **OQ-021** | RBAC lintas konteks belum ditetapkan; peninjauan ulang ADR ini bergantung padanya |
| **OQ-001** | Tidak ada DDL — nilai `Type` yang sah menurut basis data tidak diketahui |

## ADR-U-0013 - Alamat layanan keluar di-resolve runtime dari `M_LINK_SERVICE` — bukan env var

Resolusi URL keluar **harus** berupa **lookup runtime** ke tabel `M_LINK_SERVICE` lewat kunci
`(KATEGORI_1, KATEGORI_2)`, **setiap kali** layanan dipanggil.

**DILARANG** menanam alamat sebagai literal, sebagai konstanta, **maupun sebagai env var.**
Yang boleh menjadi konstanta di kode hanyalah **kunci kategori** — identifier endpoint, bukan
alamatnya.

Pemisahan dev–prod ditangani **isi tabel per-database**: database dev memuat baris dev, database
production memuat baris production. Kode tidak tahu bedanya, dan tidak perlu tahu.

> ⚠️ **ADR ini menggantikan ADR-U-0004**, yang memutuskan hal sebaliknya (alamat menjadi env var,
> lookup tidak direplikasi). Lihat §"Mengapa ADR-U-0004 dibatalkan".

#### Kontrak yang ditiru `[terverifikasi]`

`Claim Life/Activity/GetLinkService.xml`
(`ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY`) melakukan
`Obj-Browse` atas tabel `M_LINK_SERVICE` dengan

```
.KATEGORI_1 = Param.Kategori_1  AND  .KATEGORI_2 = Param.Kategori_2
```

mengambil kolom `.URL`, lalu memanggil `Connect-REST`.

`[terverifikasi]` Kunci untuk Claim — Life terbaca langsung di
`Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml`:
`Kategori_1 = "Klaim"`, `Kategori_2 = "insertClaimLife"`.

`[terverifikasi]` Lima berkas `Claim Life` memakai jalur ini: `GetLinkService.xml`,
`InsertGoogleStorage_Act.xml`, `GetUrlGoogleStorage_Act.xml`, `DeleteGoogleStorage_Act.xml`,
`serviceInsertArasapasClaimLife_act.xml`.

`[data DBA]` Isi tabel **19 baris**. Endpoint Life:
`<konfigurasi alamat layanan>

#### Mengapa ADR-U-0004 dibatalkan

ADR-U-0004 ditulis **sebelum** isi `M_LINK_SERVICE` diketahui (OQ-047 masih terbuka). Alasannya waktu
itu: "mekanisme lookup tidak perlu direplikasi; env var lebih sederhana." Dua hal membatalkannya:

1. **`[data DBA]` Tabelnya nyata dan dipakai bersama** — 19 baris melayani banyak modul, bukan hanya
   Claim Life. Memindahkan alamat Claim Life ke env var akan membuat **satu endpoint punya dua
   sumber kebenaran** selama modul lain masih membacanya dari tabel.
2. **`[keputusan work owner]`** Pemisahan dev–prod di organisasi ini memang dikerjakan lewat isi
   database, bukan lewat konfigurasi aplikasi. Env var akan memindahkan tanggung jawab itu ke tempat
   yang tidak dikelola tim yang sama.

#### Considered Options

- **Runtime lookup `M_LINK_SERVICE` setiap panggilan** — dipilih
- Env var per endpoint (**ADR-U-0004**) — ditolak, alasan di atas
- Lookup sekali saat start lalu di-cache selamanya — ditolak: perubahan alamat di tabel tidak akan
  terlihat tanpa restart, sehingga menghidupkan kembali persoalan dua sumber kebenaran.
  *(Cache ber-TTL pendek boleh dipakai sebagai optimisasi, asal kebenarannya tetap tabel.)*

#### Consequences

- **Ketergantungan Oracle bertambah satu**: memanggil layanan luar kini menuntut database hidup.
  Ini ketergantungan ketiga di luar data, setelah penomoran (**ADR-U-0006**) dan token penyimpanan.
- `internal/config` **tidak** memuat URL. Ia memuat koneksi database; alamat datang dari
  `internal/repository`.
- **Kunci kategori menjadi bagian kontrak** — mengubah `("Klaim", "insertClaimLife")` adalah
  perubahan yang memutus, setara mengubah nama endpoint.
- Efek keluar (**ADR-U-0008**) harus membedakan **kunci tidak ditemukan** dari **jaringan gagal**:
  yang pertama kesalahan konfigurasi data dan tidak layak diantre ulang berkali-kali; yang kedua
  layak.
- **ADR-U-0005** (flag lingkungan) tetap berlaku untuk *apakah* efek keluar berjalan; ADR ini hanya
  mengatur *ke mana* ia pergi.

#### Bukti pendukung: tidak ada alamat ter-hardcode untuk ditinggalkan `[terverifikasi]`

Penutupan **OQ-018** memastikan tidak ada yang perlu "dibersihkan" lebih dulu — modul `Claim Life`
memang tidak pernah menanam alamat bisnis:

| URL yang ditemukan di modul | Jumlah | Sifat |
| --- | ---: | --- |
| `https://community.pega.com/help_v88/…` | 132+ | `pyHelpURI`, tautan bantuan |
| `https://view.officeapps.live.com/op/view.aspx?src=` | 3 | penampil dokumen Office |

Pembedaan lingkungan di Pega memakai `Claim Life/When/IsPEGAPROD.xml`
(`@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN`): `pzProductionLevel = "5"` — sebuah **flag**, bukan
hostname tertanam.

#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-035** | `serviceInsertArasapasClaimLife_act` satu salinan dipakai dua konteks — kepemilikan kunci kategorinya belum ditetapkan |
| **OQ-018** (modul lain) | URL literal di 15 modul lain belum disapu; ADR ini menetapkan arah, penerapannya menyusul per modul |

## ADR-U-0014 - Keputusan komite hanya boleh diambil pemilik `KomiteID` pada tingkat berjalan

Pada setiap tingkat tangga persetujuan, **hanya pemilik
`KomiteList(KomiteCount).KomiteID`** yang boleh menyimpan keputusan. Pengguna lain **ditolak di
lapisan layanan**, meskipun ia dapat membuka kasusnya.

Ini **penyimpangan sadar**: Pega hanya *menempatkan* tugas, ia tidak *menegakkan* siapa yang
memutuskan.

#### Keadaan sekarang `[terverifikasi]`

Penempatan tugas dilakukan router:

`Komite Claim Life/Activity/KomiteRouter.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY`, 26.387 byte):

```
Property-Set   param.AssignTo = .KomiteID        (baris ~294)
  precondition .KomiteAproval == 0               (baris ~382)
```

Flow memakai `<pyImplementation>WorkList` dengan `<pyRouteTo>Custom`. Artinya kasus **muncul di
antrean** pemilik `KomiteID` — penempatan.

**Yang tidak ada:** sapuan 47 berkas modul tidak menemukan satu pun pemeriksaan bahwa pengguna yang
menyimpan keputusan adalah pemilik `KomiteID` tingkat berjalan. Keputusan masuk lewat dropdown wajib
di `Komite Claim Life/Section/ShowTransfer.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `SHOWTRANSFER` / `RULE-HTML-SECTION`) baris 32607 — dan
`AcceptStatus` **tidak ditulis rule mana pun** (nol `<PropertiesName>…AcceptStatus</PropertiesName>`).

Jadi di sistem lama, **siapa pun yang dapat membuka layar dapat menyimpan keputusan** untuk tingkat
itu. Tidak ditemukan bukti bahwa itu disengaja.

#### Considered Options

- **Tegakkan di lapisan layanan: hanya pemilik `KomiteID` tingkat berjalan** — dipilih
- Paritas (mengandalkan penempatan worklist saja) — ditolak: memindahkan lubang wewenang ke sistem
  baru, dan di Go tidak ada "worklist" yang secara kebetulan membatasi — tanpa penegakan eksplisit
  endpoint terbuka bagi siapa pun yang terautentikasi
- Tegakkan di UI saja — ditolak: sama dengan tidak menegakkan; API tetap terbuka

#### Consequences

- `internal/services` memeriksa identitas pemanggil terhadap `KomiteList(KomiteCount).KomiteID`
  **sebelum** menyimpan keputusan; gagal → tolak, bukan abaikan diam-diam.
- **Roster menjadi sumber wewenang**, bukan sekadar daftar penerima email. Menambah/menonaktifkan
  baris `EMAILKOMITE` mengubah **siapa yang boleh memutuskan** — perubahannya layak masuk jejak
  audit (**ADR-U-0007**).
- Penegakan bergantung pada pemetaan **`KomiteID` → identitas akun**. Bentuk pemetaan itu belum
  ditetapkan; ia bagian dari Identity & Access yang `[terverifikasi]` **ABSENT** dari korpus
  (`discovery/context-map.md`).
- Sejalan semangat **ADR-U-0002** (peran ditegakkan di lapisan layanan, bukan di visibilitas layar)
  dan **ADR-U-0012** (wewenang eksplisit, bukan efek samping UI).
- Konsekuensi operasional: bila pemilik `KomiteID` berhalangan, kasus akan macet — **kecuali** lewat
  jalur eskalasi manual yang ditetapkan di §"Pengecualian sah" di bawah. Di sistem lama, orang lain
  bisa menyelesaikannya tanpa jejak; kini perpindahannya eksplisit dan terekam.

#### Pengecualian sah: eskalasi manual naik satu tingkat (2026-09-15)

`[keputusan work owner]` Penegakan per `KomiteID` **tetap berlaku**, dengan **satu** pengecualian
yang ditetapkan eksplisit:

> **Bila anggota komite pada tingkat berjalan berhalangan, kasus dipindahkan NAIK SATU TINGKAT**
> untuk diaksep oleh tingkat di atasnya. Perpindahan itu adalah **aksi manual admin**, bukan
> otomatis, dan bukan "siapa saja boleh".

Ini menjawab konsekuensi operasional yang dicatat di atas — kasus **tidak** macet ketika pemilik
`KomiteID` absen.

**Yang tetap dilarang:** pemutus di tingkat **sama** yang bukan pemilik `KomiteID`, dan eskalasi
**turun** tingkat.

##### Consequences tambahan

- **Inbox komite hanya menampilkan kasus sesuai posisi** — roster tingkat itu. Bukan seluruh antrean
  komite.
- Eskalasi adalah **tindakan yang direkam**: siapa yang memindahkan, kapan, dari tingkat mana ke
  tingkat mana. Ia mengubah siapa yang berwenang, jadi ia masuk jejak audit (**ADR-U-0007**).
- `KomiteCount` **ikut naik** saat eskalasi — konsekuensinya tingkat yang dilewati **tidak pernah
  memberi keputusan**, sehingga entri `KomiteAproval` untuk tingkat itu kosong. Bentuk pencatatannya
  adalah keputusan implementasi.
- ⚠️ Eskalasi **memperpendek** tangga persetujuan yang sebenarnya. Bila kebijakan menuntut jumlah
  persetujuan minimum, eskalasi melanggarnya — `[pertanyaan terbuka]`, tetapi **tidak memblokir**:
  perilaku lama pun tidak memaksakan minimum.

#### OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-021** | RBAC lintas konteks belum ditetapkan; pemetaan `KomiteID` → akun bergantung padanya |
| **OQ-007** | Otorisasi Komite secara umum — korpus tidak memuat rule otorisasi |

## ADR-U-0015 - Efek keluar Komite Claim Life **wajib berhasil** — transactional outbox, at-least-once

Keempat efek keluar keputusan komite **wajib berhasil**. Kegagalan **tidak boleh diam**: keputusan
komite **tidak dianggap tuntas** sampai seluruh efek berhasil terkirim, atau ditandai **perlu
intervensi**.

Jaminan: **at-least-once**. Pola: **transactional outbox** — keputusan dan daftar efeknya disimpan
dalam **satu transaksi database**; worker terpisah mengirim tiap efek dengan retry sampai sukses.

> ⚠️ **Ini MENYIMPANG dari ADR-U-0008**, yang untuk **Claim — Life** memutuskan efek keluar bersifat
> asinkron dan **tidak memblokir** alur. Kedua ADR berlaku bersamaan pada konteks yang berbeda —
> lihat §"Mengapa Komite berbeda dari Claim".

#### Keempat efek keluar `[terverifikasi]`

`Komite Claim Life/Activity/KomitePostAdjustment.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`, versi **2026-09-15**,
515.675 byte). Nomor langkah dibaca dari `<pyStepPageReference>RH_1.pySteps(n)` — yang muncul
**setelah** `<pyStepsActivityName>` milik langkahnya:

| Step | Efek | Catatan |
| ---: | --- | --- |
| 6 | `Obj-Save` | **persist** — seluruh efek berjalan sesudahnya |
| 8 | `Call InsertJsonClaimLife_Act` | |
| 9 | *(gerbang EXIT retro)* | dibuang di sistem baru — lihat §Interaksi |
| 10 | `Call serviceInsertArasapasClaimLife_act` | endpoint di-lookup runtime dari `M_LINK_SERVICE`, kunci `Klaim` / `insertClaimLife` (**ADR-U-0013**) |
| 11 | `Call SendEmailKlaimLife` | |
| 12 | `Call HitServiceToKasirKMTLife_Act` | **Kasir — pembayaran.** `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `HITSERVICETOKASIRKMTLIFE_ACT` |
| 14 | `Call SetInformationData` | **di luar lingkup** — temporary, lihat §Lingkup final |

`[terverifikasi]` Step 10 dan step 12 digerbangi `AcceptStatus = 1 && KomiteCount == KomiteLoop`
(baris **8648** dan **8887**) — keduanya hanya berjalan pada **keputusan aksep di tingkat final**.

#### Mengapa Komite berbeda dari Claim

**Pemicunya Kasir.** `[terverifikasi]` `HitServiceToKasirKMTLife_Act` **tidak ada di Claim Life** —
Claim Life punya empat efek keluar, Komite punya lima (termasuk `SetInformationData`).

Kasir adalah **integrasi keuangan**: ia memicu jalur pembayaran atas klaim yang baru saja disetujui
komite tingkat tertinggi. Kegagalan diam di sini berarti **klaim disetujui tetapi tidak pernah
sampai ke pembayaran** — dan tidak ada yang tahu.

Bandingkan Claim — Life (**ADR-U-0008**): di sana efek keluar adalah unggah berkas, email, Arasapas,
dan konversi — tidak satu pun memindahkan uang, dan `[terverifikasi]` Pega sendiri hanya
**mencatat** kegagalannya (`InsertLogServiceClaim` → `monitoring_klaim_log`) tanpa gerbang
keberhasilan. Di sana "tidak memblokir" adalah paritas yang masuk akal.

Di Komite, taruhannya berbeda. Karena itu jaminannya dinaikkan.

#### Considered Options

- **Transactional outbox + retry, at-least-once** — dipilih
- Paritas dengan ADR-U-0008 (asinkron, tidak memblokir) — ditolak: kegagalan Kasir menjadi tak terlihat
- Panggil sinkron di dalam transaksi keputusan — ditolak: keputusan komite akan gagal hanya karena
  layanan luar sedang mati, dan transaksi database menggantung selama panggilan jaringan

#### Consequences

- **Keputusan komite punya dua keadaan yang harus dibedakan**: *tersimpan* dan *tuntas*. UI dan
  laporan tidak boleh menyamakan keduanya.
- Diperlukan **keadaan "perlu intervensi"** untuk efek yang gagal berulang — beserta tempat orang
  melihatnya. Tanpa itu, "wajib berhasil" hanya menggeser kegagalan diam ke antrean.
- **At-least-once berarti penerima harus tahan duplikat.** Arasapas, Kasir, dan email dapat menerima
  kiriman yang sama lebih dari sekali. `[terbuka]` Apakah ketiganya idempoten **belum diketahui** —
  ini risiko nyata yang dibawa keputusan ini.
- Outbox disimpan dalam transaksi yang sama dengan keputusan → sejalan **OQ-013** (Go memegang batas
  transaksi untuk jalur Life).
- Retry harus membedakan **kunci kategori tidak ditemukan** di `M_LINK_SERVICE` (kesalahan
  konfigurasi data, tidak layak diulang berkali-kali) dari **jaringan gagal** (layak) — **ADR-U-0013**.

#### Lingkup final dan syarat pengaman (2026-09-15)

##### Lingkup: **empat** efek keluar `[terverifikasi]`

Dibaca ulang dari `KomitePostAdjustment.xml` versi **2026-09-15** (515.675 byte). Nomor langkah dari
`<pyStepPageReference>`, status dari `<pyStepsBlockName>` (kosong = aktif):

| Urutan | Step | Efek | Status |
| ---: | ---: | --- | --- |
| 1 | **8** | `Call InsertJsonClaimLife_Act` | AKTIF |
| 2 | **10** | `Call serviceInsertArasapasClaimLife_act` | AKTIF — endpoint via `M_LINK_SERVICE`, kunci `Klaim`/`insertClaimLife` (**ADR-U-0013**) |
| 3 | **11** | `Call SendEmailKlaimLife` | AKTIF |
| 4 | **12** | `Call HitServiceToKasirKMTLife_Act` | AKTIF — **Kasir/pembayaran** |

Seluruhnya **setelah `Obj-Save` (step 6)**.

**Di luar lingkup** `[keputusan work owner]`: step **14** `Call SetInformationData` — hanya
**temporary**, menampung hasil submit (`"Approve"`/`"Reject"` + `AcceptedNo` + `CLAIM_NO`). Bukan
efek keluar, bukan penyimpanan permanen, **tidak** wajib-berhasil.

##### Syarat pengaman wajib `[keputusan work owner]`

Keempat efek **boleh di-retry berkali-kali**. Karena itu:

1. **ID idempoten unik per kiriman.** Tiap kiriman membawa identifier yang memungkinkan penerima
   **menolak duplikat**. Tanpa ini, at-least-once berarti duplikasi, bukan keandalan.
2. **Cek status sebelum kirim ulang — wajib untuk Email dan Kasir.** Sebelum retry, sistem
   memeriksa apakah kiriman sebelumnya **sudah sukses**. Email dobel mengganggu; **pembayaran dobel
   merugikan**.
3. **Status "perlu intervensi" terlihat.** Kegagalan yang tidak pulih setelah retry ditampilkan
   sebagai **status eksplisit di UI Komite** dan masuk **laporan harian**. Tanpa tempat orang
   melihatnya, "wajib berhasil" hanya memindahkan kegagalan diam ke antrean.

Butir 1 dan 2 **menjawab** risiko idempotensi yang dicatat di §Consequences: pertanyaannya bukan
lagi "apakah penerima idempoten" melainkan **"sistem baru wajib membuatnya aman"** — dengan ID unik
di sisi pengirim dan pemeriksaan status sebelum kirim ulang.

⚠️ Sisa risiko yang diterima: bila penerima **tidak** menghormati ID idempoten, butir 2 adalah
pertahanan terakhir — dan ia bergantung pada catatan status milik kita sendiri, bukan konfirmasi
penerima.

#### Interaksi dengan pembuangan identitas retro ter-hardcode

⚠️ `[terverifikasi]` Di korpus versi 2026-09-15, **step 9 adalah langkah gerbang tersendiri** —
tanpa `pyStepsActivityName`, membawa `<pyStepsPreCondition>true` dan
`<pyStepsDescription>EXIT JIKA RETROID "L0000141"</pyStepsDescription>` (baris 8483), dengan
precondition `RetroID=="L0000141" || SecurityReinsurerID=="L0000134"` (baris 8525) dan
`RetroID=="1000013"` (baris 8548). Ia berdiri **di antara** efek 1 (step 8) dan efek 2 (step 10),
sehingga saat memicu, **step 10–14 tidak berjalan**: Arasapas, email, Kasir, dan
`SetInformationData` semuanya di-skip.

`[dugaan]` Bahwa kode transisi numeriknya berarti "Exit Activity" **tidak dapat dibuktikan dari
ekspor** — nilai yang terbaca (`2`) muncul juga di step 8 yang bukan EXIT. Yang terbukti adalah
label penulisnya dan bentuk langkahnya.

`[keputusan work owner]` Ketiga identitas itu **dibuang** (**OQ-064**). Akibatnya: klaim yang selama
ini dikecualikan akan **mulai menerima keempat efek keluar** — termasuk **Kasir**.

Digabung dengan ADR ini ("wajib berhasil"), perubahan itu berarti klaim-klaim tersebut kini
**menuntut keberhasilan pembayaran** yang sebelumnya tidak pernah dipicu. **Ini perubahan perilaku
yang menyentuh uang**, dan layak dikonfirmasi Product+UW sebelum rilis — lihat **OQ-064**.

#### OQ yang masih terbuka dan menyentuh ADR ini
| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-064** | **TERTUTUP** 2026-09-15 — gerbang EXIT retro dibuang; semua klaim menjalankan keempat efek. Arti ketiga identitas tetap tidak diketahui dan tidak dicari lagi |
| **OQ-065** (baru) | Langkah "Tukar SecurityReinsurer dengan RetroName" (step 4.14 & 5.5) **ter-remark** — apa yang kini masuk ke kolom retro pada rekam akseptasi |
| **OQ-035** | `serviceInsertArasapasClaimLife_act` satu salinan dipakai dua konteks |
| **OQ-002** | Kontrak layanan **Kasir** tidak ada di korpus — bentuk permintaan, makna jawaban, dan apakah ia menghormati ID idempoten belum diketahui |

## ADR-U-0016 - Kolom uang memakai `NUMBER(38,8)`

Seluruh kolom uang di sistem baru memakai **`NUMBER(38,8)`** — tiga puluh digit di depan koma,
delapan di belakang. Berlaku **seluruh sistem**, bukan satu modul.

ADR-U-0003 menetapkan bahwa uang **tidak** direpresentasikan sebagai `float`. Ia tidak menetapkan
presisinya. Catatan ini menutup lubang itu.

#### Kenapa delapan desimal di belakang koma

`[terverifikasi]` Nilai berdesimal terpanjang yang ditemukan pada dokumen produksi membawa **24
angka** di belakang koma. Seluruh selisih yang bermakna jatuh pada desimal **ketujuh atau
sebelumnya**; sisanya ekor galat penjumlahan.

Nilai berdesimal lebih dari delapan **dibulatkan saat dimuat, bukan ditolak**.

#### Kenapa tiga puluh digit di depan koma — bukan dua belas

Angka semula yang diusulkan adalah `NUMBER(20,8)`, yaitu dua belas digit di depan koma. Sapuan
korpus membuktikan itu **tidak cukup**:

| Nilai di korpus | Digit di depan koma | Sifat |
| --- | ---: | --- |
| `99999999999999999.99` | **17** | sentinel "tanpa batas" |
| `181500000000.00` | **12** | ambang kapasitas treaty |
| `150000000000.00` | **12** | ambang kapasitas treaty |
| `200000000.00` | 9 | batas premi neto |

⭐ Ambang dagang yang nyata sudah **tepat di batas** dua belas digit — nol ruang lega. Penjumlahan
lintas mata uang dapat melewatinya, dan menaikkan presisi **sesudah** data terisi jauh lebih mahal
daripada menaikkannya sekarang.

#### Kenapa pelebaran ini tidak berbiaya

`NUMBER` di Oracle berpanjang **berubah-ubah**: hanya digit bermakna yang tersimpan. Kolom berisi
`1500` memakan ruang yang sama entah dideklarasi `(20,8)` atau `(38,8)`. Deklarasi menetapkan
**batas**, bukan ukuran tersimpan.

⭐ Karena itu tidak ada alasan menahan presisi di angka yang pas-pasan. **38 adalah presisi
tertinggi yang diterima Oracle.**

#### Akibat

1. Seluruh kolom uang pada modul mana pun dideklarasi `NUMBER(38,8)`. Tidak ada modul yang
   menetapkan presisi sendiri.
2. Kolom **persen** tidak termasuk — ia bukan uang, dan tetap mengikuti ketetapan modulnya.
3. Lapisan `services` membulatkan pada desimal kedelapan **saat memuat**, bukan tiap kali membaca.
4. Test yang menemukan kolom uang bertipe `float`, atau berpresisi selain `(38,8)`, **gagal**.

#### Yang catatan ini TIDAK putuskan

⛔ Bukan tentang pembulatan ke presisi mata uang untuk **tampilan**. Itu urusan penyajian.

⛔ Bukan tentang satuan minor. Sistem ini memakai desimal berskala tetap, bukan bilangan bulat
dalam satuan minor.

## ADR-U-0017 - Daftar medan disusun dari aturan, bukan dari panduan bentuk dokumen

Daftar kolom untuk memecah dokumen JSON menjadi baris disusun dengan **menyapu aturan Pega** —
Activity, Section, DataTransform, Flow, RDBList — **bukan** dari keluaran `JSON_DATAGUIDE`.

Berlaku untuk setiap modul yang memecah dokumen JSON tersimpan, bukan hanya modul tempat
keputusan ini lahir.

#### Kenapa panduan bentuk dokumen tidak dipakai sebagai dasar

`[terverifikasi]` Panduan yang diterima memuat 378 jalur. **Satu dokumen produksi tunggal memuat
95 jalur yang tidak ada di dalamnya** — dan bukan medan pinggiran: daftar perusahaan ceding,
seluruh daftar pecahan penyebaran, nilai klaim, nilai kelebihan kerugian, dua kolom hasil, nilai
sisa, pengenal induk dokumen.

⭐ Artinya panduan itu **sah** membuktikan sebuah medan **ada**, tetapi **tidak sah** membuktikan
sebuah medan **tidak ada**. Dasar yang hanya bisa membuktikan satu arah tidak dapat dipakai untuk
mengunci daftar kolom.

#### Kedudukan panduan bentuk dokumen: penambal

Ia tetap berguna untuk dua hal:

1. **Panjang maksimum per medan** — menghapus kebutuhan menebak ukuran kolom.
2. **Medan yang ditulis sistem**, yang karena itu tidak muncul di aturan mana pun. Pada sapuan
   pertama ada **24 nama** semacam itu.

⇒ Yang dipakai adalah **gabungan**: sapuan aturan sebagai dasar, panduan sebagai penambal.

#### Dua titik buta yang wajib ditutup saat menyapu

`[terverifikasi]` Keduanya sudah terbukti meleset **empat kali** dalam proyek ini.

**1 · Dua konvensi tag yang berlawanan.**
`Activity` menyimpan penugasan langkah pada `PropertiesName` dan `PropertiesValue` — **tanpa**
awalan `py`. `DataTransform` dan `Flow` memakai `pyPropertiesName` dan `pyPropertiesValue` —
**dengan** awalan. Sapuan yang hanya mengenal satu konvensi kehilangan separuh isi.

**2 · Rujukan relatif di dalam loop.**
Sapuan yang berjangkar pada nama halaman *(`PolicyTreatyIn.…`)* tidak melihat medan yang ditulis
relatif di dalam perulangan *(`.PremiumSpreaded`, `.CedingCoName`)*. Sapuan wajib menangkap
keduanya.

⚠️ Dua jebakan teknis lain: buang `pyExpressionGadget` sebelum mencocokkan, dan **jangan**
meng-unescape entitas HTML sebelum mencocokkan pola struktur — `&lt;` memecahkan pola tag.

#### Akibat

1. Setiap modul yang memecah dokumen JSON menyapu aturannya sendiri lebih dulu, dan menyimpan
   hasilnya sebagai daftar medan yang dapat diperiksa.
2. ⭐ **Cacah hasil sapuan adalah batas bawah, bukan total.** Pemecah dokumen **wajib** menyediakan
   penampung medan tak dikenal, dan penampung itu **wajib kosong** sebelum pekerjaan dinyatakan
   selesai. Itu jaring pengaman yang menggantikan kepastian yang memang tidak ada.
3. Test yang menemukan medan hilang **diam-diam** — tanpa masuk penampung — **gagal**.

#### Yang catatan ini TIDAK putuskan

⛔ Bukan tentang tipe kolom. Itu ADR-U-0016 untuk uang, dan ketetapan modul untuk sisanya.

⛔ Bukan larangan memakai `JSON_DATAGUIDE` sama sekali — ia tetap dipakai sebagai penambal.

## ADR-U-0018 - Generasi polis adalah rantai baris, bukan penyuntingan di tempat

Endorsemen **tidak menyunting** baris polis yang ada. Ia menambah **baris baru** yang menunjuk
baris sebelumnya. Riwayat polis adalah rantai baris, dan baris lampau **tidak boleh disunting**.

#### Bentuknya

| Unsur | Ketetapan |
| --- | --- |
| Penunjuk generasi sebelumnya | kolom penunjuk ke baris induk, **nullable** |
| Polis baru | penunjuk **kosong**, nomor generasi `0` |
| Endorsemen | penunjuk **terisi**, nomor generasi `>= 1` |
| Kunci alami | `(nomor polis, nomor generasi)` - **unik** |
| Larangan percabangan | penunjuk generasi sebelumnya juga **unik** |

Keunikan pada **penunjuk** itulah yang melarang percabangan: satu generasi hanya boleh punya satu
penerus. Tanpanya, dua endorsemen dapat menunjuk induk yang sama dan riwayat bercabang tanpa ada
yang tahu mana yang sah.

#### Kenapa bukan menyunting di tempat

1. Angka selisih endorsemen dihitung sebagai `generasi baru - generasi yang ditunjuk`. Menyunting
   di tempat **menghapus operand pengurangnya**.
2. Nilai lama wajib dapat dibaca kembali apa adanya untuk pemeriksaan. Penyuntingan membuatnya
   hilang tanpa jejak.
3. Nomor generasi di sistem lama disimpan pada medan **dua digit**, batas 99. Sistem baru **tidak
   memakai batas itu** - nomor generasi adalah **bilangan bulat**, bukan teks dua digit. Urutannya
   harus benar secara angka: sebagai teks, `"100"` jatuh sebelum `"99"`.

#### Akibat

1. Lapisan `repository` tidak menyediakan jalur pembaruan untuk baris generasi lampau.
2. Pembekuan itu ditegakkan **di lapisan layanan**, bukan hanya diandalkan pada disiplin pemakai.
3. Test yang berhasil menyunting baris generasi lampau, atau berhasil membuat dua generasi dengan
   induk yang sama, **gagal**.
4. Berlaku untuk setiap modul yang punya endorsemen - bukan hanya tempat keputusan ini lahir.

#### Yang catatan ini TIDAK putuskan

Berapa kali satu polis boleh di-endorse. Tidak ada batas dagang; batas teknis lama tidak
dipertahankan.

Bagaimana baris antar generasi dipasangkan - itu ADR-U-0019.

## ADR-U-0019 - Baris anak dipasangkan antar generasi lewat nomor urut, bukan kunci dagang

Setiap tabel anak membawa **nomor urut baris** di dalam induknya. Nomor urut itulah yang
memasangkan baris generasi lama dengan baris generasi baru saat selisih dihitung.

**Bukan kunci dagang.** Kunci dagang - nama treaty, mata uang, nomor angsuran, lapisan - **boleh
berubah** di dalam endorsemen. Memakainya sebagai pemasang berarti pasangan putus persis ketika
endorsemen mengubahnya, yaitu ketika selisih paling dibutuhkan.

#### Dua perilaku yang berbeda, sengaja

| | Polis baru | Endorsemen |
| --- | --- | --- |
| Menghapus baris | **boleh** | **tidak bisa** |
| Nomor urut sesudah penghapusan | **dinomori ulang rapat** - 3 baris, yang ke-2 dihapus, sisanya jadi 1 dan 2 | tidak berlaku |
| Baris baru | di belakang | di belakang, `maksimum + 1` |

Perbedaan ini bukan ketidakkonsistenan. Pada polis baru belum ada generasi sebelumnya, sehingga
penomoran ulang tidak memutus pasangan apa pun. Pada endorsemen, penomoran ulang **akan**
memutusnya.

#### Aturan keutuhan

Generasi `n+1` wajib memuat **setiap nomor urut** yang ada pada generasi `n`. Generasi yang
kehilangan satu nomor urut **ditolak** - karena baris yang hilang berarti selisih yang tidak dapat
dihitung.

#### Akibat

1. Nomor urut ditetapkan lapisan `repository` saat menulis, bukan dikirim dari layar.
2. Lapisan layanan menolak generasi yang kehilangan nomor urut, dan menolak penghapusan baris pada
   endorsemen.
3. Test wajib mencakup: tiga baris, yang tengah dihapus pada polis baru - sisanya bernomor 1 dan 2;
   dan generasi endorsemen yang kehilangan satu nomor urut - **ditolak**.

## ADR-U-0020 - Tabel proyeksi selisih bukan tabel sumber

Angka selisih endorsemen disimpan pada tabel **proyeksi** supaya dapat dibaca lewat SQL tanpa
menjalankan perhitungan. Tabel itu **bukan** sumber kebenaran - ia salinan hasil hitung.

Tiga aturan berikut yang membuatnya bukan tabel sumber. Ketiganya wajib, tidak boleh diringkas.

| | Aturan |
| ---: | --- |
| 1 | Ia **hanya ditulis** oleh sistem baru, di dalam **transaksi yang sama** dengan generasinya |
| 2 | Ia **boleh dihapus total dan dibangun ulang** - **tetapi hanya baris hasil hitung sistem baru** |
| 3 | Bila isinya berbeda dari hasil hitung ulang, **tabelnya yang salah**, bukan operannya |

#### Kolom penanda asal baris

Setiap baris membawa penanda asal: hasil **migrasi** atau hasil hitung **sistem baru**.

Baris hasil migrasi **beku**. Perintah bangun ulang tidak boleh menyentuhnya - angkanya berasal
dari sistem lama dan tidak dapat dihitung ulang, karena rumus lama punya varian yang sengaja tidak
ditiru.

**Tanpa aturan 2, satu perintah bangun ulang menimpa seluruh angka historis dan tidak dapat
dikembalikan.** Itu sebabnya aturan ini ditulis sebagai catatan keputusan, bukan sekadar komentar
di kode.

#### Kenapa bukan view basis data

View pernah diusulkan lalu **dibatalkan**. Alasannya: rumus selisih hidup di lapisan layanan
*(ADR-U-0021)*. View berarti rumus yang sama ditulis **dua kali** - sekali di kode, sekali di SQL -
dan cepat atau lambat keduanya bercabang tanpa ada yang tahu.

**Jangan dihidupkan kembali diam-diam.**

#### Akibat

1. Tabel proyeksi ditulis di dalam transaksi yang sama dengan barisnya - tidak menyusul.
2. Perintah bangun ulang menyaring pada penanda asal. Test yang menemukan baris hasil migrasi ikut
   terhapus **gagal**.
3. Kunci penyaring disertakan pada tabel proyeksi supaya pembaca SQL tidak perlu join balik.

## ADR-U-0021 - Rumus selisih hidup di lapisan layanan, dan nilai lama tidak pernah dihitung ulang

Selisih endorsemen dihitung di lapisan **layanan**, bukan di basis data. Lapisan `repository`
hanya mengambil **dua baris**: baris generasi ini dan baris yang ditunjuknya.

#### Satu rumus untuk semua

```
selisih.X = baris_ini.X - baris_yang_ditunjuk.X
```

Tanpa memandang endorsemen pertama atau berlapis, proporsional atau non-proporsional.

`[terverifikasi]` Terbukti dari dokumen produksi pada enam medan sekaligus - nilai baru nol, nilai
lama positif, selisih negatif sebesar nilai lama.

Empat aturan turunannya:

| Golongan | Perlakuan |
| --- | --- |
| Uang | **dikurangi** |
| Persen | **disalin, tidak dikurangi** |
| Kunci | disalin |
| Penanda arah | diturunkan dari **tanda jumlah**, bukan dari satu baris |

#### Nilai lama tidak pernah dihitung ulang

`[keputusan work owner]` Angka selisih yang berasal dari sistem lama **disalin apa adanya** saat
migrasi. Ia **tidak** dihitung ulang dengan rumus baru.

Sebabnya: sistem lama memuat **varian rumus** yang sengaja **tidak ditiru** - varian yang
mengurangi terhadap nilai yang ternyata kosong, sehingga aritmetiknya tampak ganjil. Menghitung
ulang berarti mengubah angka historis yang sudah dilaporkan.

`[penyimpangan sadar]` Perbedaan antara rumus baru dan varian lama itu **diterima dengan sadar**,
dan ditandai pada barisnya supaya terlihat, bukan disembunyikan.

#### Akibat

1. Nol rumus selisih di basis data - nol view, nol kolom terhitung, nol trigger.
2. `repository` tidak menghitung apa pun; ia mengambil dua baris dan mengembalikannya.
3. Test membandingkan hasil lapisan layanan terhadap angka dokumen produksi, bukan terhadap tabel
   proyeksi - sebab tabel proyeksi adalah salinan, bukan sumber *(ADR-U-0020)*.

## ADR-U-0022 - Konversi tipe terjadi sekali saat masuk, bukan tiap kali dibaca

Dokumen JSON menyimpan hampir seluruh nilai sebagai **teks**. Konversi ke tipe yang benar terjadi
**satu kali saat dimuat**, dan kolom disimpan dengan tipe aslinya - bukan disimpan sebagai teks
lalu dikonversi tiap kali dibaca.

| Golongan | Tipe kolom |
| --- | --- |
| Uang, persen | desimal berskala tetap *(ADR-U-0003, ADR-U-0016)* |
| Tanggal | tipe tanggal |
| Bilangan | tipe bilangan |
| Kode dan penanda | **tetap teks** |

#### Kenapa kode dan penanda tetap teks

`[terverifikasi]` **Nol di depan membawa makna.** Sebuah kode bernilai `"006"` yang disimpan
sebagai bilangan lalu dibaca kembali menjadi `"6"` akan **lolos uji pulang-pergi** tetapi
memecahkan penggolongnya.

Penanda juga tetap teks, sebab di sistem lama ia **dibandingkan sebagai teks**.

Karena itu **uji pulang-pergi saja tidak cukup.** Sebagian test wajib memeriksa **nilai kolom
langsung**.

#### Dua hal yang mudah terlewat

1. **Dua format tanggal dalam satu dokumen** - bentuk ringkas dan cap waktu lengkap. Pemuat wajib
   mengenali keduanya.
2. **Teks kosong pada kolom angka atau tanggal menjadi kosong**, bukan nol dan bukan tanggal nol.

#### Akibat

1. Konversi adalah tanggung jawab pemuat di `repository`, bukan tersebar di lapisan layanan.
2. Perbandingan nilai uang **tidak boleh sama-persis biner** - ia dibandingkan sebagai desimal.
3. Test memeriksa nilai kolom langsung untuk kode berawalan nol.

## ADR-U-0023 - Pemecah dokumen wajib punya penampung medan tak dikenal

Setiap pemecah dokumen JSON menyediakan **penampung** bagi medan yang tidak dikenal daftar
kolomnya. Medan yang tidak dikenali **masuk penampung**, tidak dibuang.

Penampung itu **wajib kosong** sebelum pekerjaan dinyatakan selesai.

#### Kenapa

ADR-U-0017 menetapkan daftar medan disusun dari sapuan aturan. Sapuan itu jujur tetapi **tidak dapat
membuktikan kelengkapan** - cacahnya adalah **batas bawah**, bukan total. Medan yang tidak pernah
dirujuk aturan mana pun tetap dapat hadir pada dokumen tersimpan.

Penampung adalah jaring pengaman yang menggantikan kepastian yang memang tidak ada.

Tanpa penampung, medan yang tidak dikenali **hilang diam-diam**. Uji pulang-pergi tidak akan
menangkapnya, sebab yang hilang tidak pernah masuk untuk dibandingkan.

#### Akibat

1. Pemecah dokumen menulis medan tak dikenal ke penampung, bukan mengabaikannya dan bukan gagal.
2. Isi penampung diperiksa sebagai bagian dari kriteria selesai. Penampung berisi = pekerjaan
   belum selesai.
3. Test yang menemukan medan hilang tanpa masuk penampung **gagal**.
4. Berlaku untuk setiap modul yang memecah dokumen JSON tersimpan.

## ADR-U-0024 - Pembatalan adalah generasi bernilai nol, bukan penghapusan

Membatalkan polis **tidak menghapus** apa pun. Ia menambah **generasi baru** yang seluruh kolom
nilainya **nol**.

Tiket atau kode yang membuat jalur penghapusan untuk pembatalan **salah**.

#### Kenapa

Pembatalan adalah peristiwa dagang yang harus terbaca di riwayat, bukan ketiadaan data. Sebagai
generasi bernilai nol ia ikut aturan rantai *(ADR-U-0018)* dan menghasilkan selisih yang benar
dengan sendirinya: `0 - nilai lama = -nilai lama`, tanpa jalur perhitungan khusus.

Sistem lama pun tidak menghapus. Perilaku ini **ditiru**, bukan diciptakan.

#### Sesudah dibatalkan

`[keputusan work owner]` Polis yang sudah dibatalkan **masih boleh di-endorse lagi**. Generasi baru
boleh ditumpuk di atas generasi pembatalan.

**Nol gerbang peran**, nol alur persetujuan tambahan. Layar memberi **peringatan** sebelum pengguna
melanjutkan - peringatan, bukan penghalang.

| Lapisan | Tugasnya |
| --- | --- |
| layanan | **tidak** menolak generasi di atas pembatalan |
| `repository` | mengembalikan **jenis generasi sebelumnya** supaya layar tahu kapan memperingatkan |
| antarmuka | menulis peringatannya |

Keputusan ini diambil dengan kata *"dulu"*. Bila kelak gerbang peran diperlukan, keputusan itu
diambil terpisah - dan rancangan ini dibuat supaya penambahannya murah.

#### Akibat

1. Nol jalur `DELETE` untuk pembatalan di lapisan mana pun.
2. Test wajib: rantai **generasi 1 - 2 - pembatalan - 4** tersimpan utuh dan terbaca utuh.

## ADR-U-0025 - Tabel inti berbagi kunci utama dengan tabel kerja

Hubungan antara **tabel kerja** dan **tabel inti** di bawahnya **tidak memakai kolom penyambung**.
Keduanya memakai **kunci utama yang sama persis**.

#### Bentuknya

| | |
| --- | --- |
| Tabel kerja | satu baris per objek kerja, memegang identitas |
| Tabel inti | satu baris per objek kerja, berbagi kunci utama yang sama |
| Kolom penyambung | **tidak ada** |

`[keputusan work owner]` Pola ini dipakai pada Claim Life dan pada Treaty In. Ia bukan kebetulan
dua modul, melainkan bentuk yang sama untuk persoalan yang sama.

#### Kenapa

1. Hubungannya **satu lawan satu dan wajib**. Kolom penyambung terpisah hanya menambah tempat
   untuk tidak sinkron, tanpa menambah keterangan apa pun.
2. Penggabungan menjadi murah dan tidak mungkin salah pasang.
3. Menghapus pertanyaan "kunci mana yang dipakai" pada setiap tabel anak di bawahnya - seluruhnya
   memakai kunci yang sama.

#### Akibat

1. Tabel inti tidak punya sequence sendiri. Identitas lahir di tabel kerja.
2. Baris tabel inti tidak dapat ada tanpa baris tabel kerja. Urutan penulisan mengikutinya.
3. Test yang menemukan kolom penyambung terpisah antara keduanya **gagal**.

#### Yang catatan ini TIDAK putuskan

Hubungan tabel inti ke tabel **anak** di bawahnya. Itu kunci tamu biasa, bukan kunci bersama.

## ADR-U-0026 - Tabel kerja adalah akar lintas-lini, bukan tabel di samping

Satu **tabel kerja** menjadi akar bagi lebih dari satu lini pekerjaan. Ia tidak berdiri di samping
pohon sebagai pelengkap - ia **akarnya**.

Baris dari lini yang berbeda hidup di tabel yang sama, dibedakan oleh **kolom pembeda lini**.

#### Contoh yang sudah berjalan

| Tabel kerja | Menaungi |
| --- | --- |
| tabel kerja klaim | klaim **dan** kasus komite, dibedakan kolom pembeda |
| tabel kerja polis | polis treaty **dan** daftar premi |

`[terverifikasi]` Keduanya sudah dipakai lebih dari satu modul sebelum keputusan ini ditulis.
Catatan ini mengesahkan yang sudah terjadi, bukan mengusulkan yang baru.

#### Kenapa

1. Objek kerja punya daur hidup, pemilik, dan jejak audit yang **bentuknya sama** lintas lini.
   Menduplikasinya per lini berarti menduplikasi daur hidup itu.
2. Penyerahan pekerjaan antar lini - misalnya klaim diserahkan ke komite - menjadi hubungan induk
   ke anak di dalam **satu** tabel, bukan jembatan antar tabel.
3. Laporan lintas lini tidak perlu menggabungkan tabel yang bentuknya berbeda-beda.

#### Akibat

1. **Setiap modul baru memeriksa dulu** apakah tabel kerja yang ada sudah menaunginya, sebelum
   mengusulkan tabel kerja sendiri.
2. Kolom pembeda lini **wajib** ada dan **wajib** terisi. Baris tanpa pembeda tidak sah.
3. Penghapusan baris anak di dalam tabel kerja ditegakkan di lapisan layanan, bukan oleh kaskade
   basis data - sebab satu tabel menaungi lebih dari satu arti.
4. Test yang menemukan baris tanpa kolom pembeda **gagal**.

## ADR-U-0027 - Kolom basis data nullable; kewajiban isi ditegakkan di kode

`[keputusan work owner]` **Seluruh kolom dideklarasi nullable.** Kewajiban mengisi ditegakkan di
lapisan aplikasi, bukan oleh batasan `NOT NULL`.

`[penyimpangan sadar]` Ini **penyimpangan yang disadari** dari kebiasaan umum, dan dicatat sebagai
penyimpangan, bukan disamarkan sebagai praktik terbaik.

#### Kenapa

1. **Data lama tidak lengkap.** Migrasi membawa baris dari sistem yang tidak pernah mewajibkan
   pengisian. Batasan `NOT NULL` akan menolak data historis yang sah, dan memaksa mengarang nilai
   pengganti - yang lebih berbahaya daripada kosong.
2. **Kewajiban isi bergantung keadaan.** Sebuah kolom wajib pada satu langkah alur dan tidak pada
   langkah lain. Batasan basis data tidak mengenal langkah alur; kode mengenalnya.
3. Pesan galat dari aplikasi dapat dibaca pemakai. Pesan galat dari batasan basis data tidak.

#### Harga yang dibayar

Basis data **tidak lagi menjadi jaring pengaman terakhir**. Bila kode lalai, kolom kosong akan
tersimpan tanpa ada yang menolak.

Karena itu dua hal menjadi wajib:

1. Pemeriksaan kewajiban isi terletak di **satu tempat** per objek, bukan tersebar di banyak jalur.
2. Test memeriksa penolakan itu secara langsung - bukan mengandalkan basis data menolaknya.

#### Akibat

1. Nol `NOT NULL` pada kolom selain kunci utama.
2. Pemuat migrasi tidak perlu mengarang nilai pengganti untuk data historis yang kosong.
3. Test yang menemukan kolom wajib-isi lolos tersimpan kosong lewat lapisan layanan **gagal**.

## ADR-U-0028 - Sistem hilir membaca tabel, bukan muatan JSON

`[keputusan work owner]` Sistem hilir - produksi, pelaporan, layanan luar - membaca **langsung
dari tabel** sistem baru.

Seluruh penyusun muatan JSON keluar yang ada di sistem lama **tidak dimigrasikan**.

#### Kenapa

1. Muatan JSON di sistem lama adalah **akibat** dari dokumen JSON sebagai bentuk simpan. Begitu
   data tersimpan relasional, muatan itu kehilangan alasan keberadaannya.
2. Menyusun ulang JSON dari tabel berarti **menulis bentuk yang sama dua kali** - sekali sebagai
   tabel, sekali sebagai muatan - dan keduanya akan bercabang.
3. Pembaca hilir yang membaca tabel mendapat kolom bertipe benar. Pembaca muatan JSON mendapat
   teks, dan harus mengurai ulang.

#### Yang tetap dipertahankan

Kontrak efek keluar - kapan dikirim, ke mana, apa yang menandakan berhasil - **tidak berubah**.
Yang berubah hanya **bentuk** bacaannya.

Tabel hilir yang sudah datar dan sudah dipakai sistem lain **tidak dibuat ulang**. Ia diisi apa
adanya.

#### Akibat

1. Rule penyusun muatan JSON keluar tidak masuk lingkup migrasi. Ia masuk daftar yang **sengaja
   tidak dibangun**.
2. Tabel yang dibaca hilir membawa **kunci penyaring** supaya pembaca tidak perlu menggabungkan
   balik.
3. Perubahan bentuk tabel yang dibaca hilir adalah perubahan kontrak - diperlakukan sebagai
   perubahan luar, bukan perubahan dalam.

## ADR-U-0029 - Batas transaksi dipegang aplikasi; `COMMIT` tidak tertanam di teks SQL

`[keputusan work owner]` Aplikasi yang membuka dan menutup transaksi. **Nol `COMMIT` di dalam teks
SQL aplikasi.**

Di sistem lama `COMMIT` tertanam di dalam blok anonim di badan rule - `BEGIN ... ; COMMIT; END;`.
Di sistem baru ia **dikeluarkan oleh aplikasi**, sesudah panggilan.

Objek basis data yang dipanggil **tetap sama**. Logika penulisannya **tidak direplikasi** di
aplikasi - mereplikasinya berarti menebak isi yang tidak terbaca.

#### Bahaya yang membuat urutan pemanggilan menjadi wajib

`[data DBA]` Sebagian stored procedure memasang penanganan galat yang menjalankan **pembatalan
transaksi di dalam dirinya sendiri**.

Pembatalan telanjang di Oracle membatalkan **seluruh** transaksi - termasuk apa pun yang aplikasi
tulis **sebelumnya** dalam transaksi yang sama. Pembatalan sampai titik simpan **tidak menolong**,
sebab yang dipasang bukan itu.

Karena itu:

| | Aturan |
| ---: | --- |
| 1 | Panggilan procedure semacam itu adalah **langkah terakhir sebelum penutupan transaksi** |
| 2 | **Tidak boleh ada pekerjaan penting yang masih menggantung** sebelumnya dalam transaksi yang sama |
| 3 | Aplikasi **memeriksa penanda keberhasilan** yang dikembalikan procedure, dan tidak menganggap ketiadaan galat sebagai keberhasilan |

#### Akibat

1. Seluruh urutan penyimpanan satu objek kerja dibungkus **satu transaksi**.
2. Test memeriksa bahwa kegagalan di tengah tidak meninggalkan baris separuh jadi.
3. Pemeriksaan penanda keberhasilan adalah bagian dari jalur, bukan pemeriksaan tambahan yang
   boleh dilewati.

## ADR-U-0030 - Aturan peran ditetapkan sekali dan berlaku lintas modul

`[keputusan work owner]` Aturan peran sistem baru **ditetapkan sekali** dan berlaku lintas modul.
Ia **tidak** dirancang ulang per modul.

Penegakannya di **lapisan layanan**, bukan di layar.

#### Kenapa ini bukan soal kerapian

`[terverifikasi]` Sensus korpus menemukan lubang yang sama pada **dua modul beruntun**: seluruh
medan hak akses ada di layar tetapi **kosong seluruhnya**. Ratusan kemunculan nama medan hak
akses, **nol** yang terisi. Satu-satunya yang terisi hanya menyebut **nama kelas**, bukan nama hak.

Artinya gerbang wewenang di sistem lama **hidup di layar, dan sebagian besar tidak hidup sama
sekali**. Merancang ulang per modul berarti menyalin lubang itu berkali-kali, dan menemukannya
kembali berkali-kali.

#### Bentuk penegakan

1. Wewenang diperiksa **di lapisan layanan**, pada setiap jalur yang mengubah keadaan. Layar boleh
   menyembunyikan tombol, tetapi penyembunyian **bukan** penegakan.
2. Sumber peran adalah **satu tabel**, bukan nama orang yang ditanam di kode.
3. `[penyimpangan sadar]` Nama orang yang ditanam di kode sistem lama **dibuang**. Perilaku tidak
   berubah - orang yang sama tetap mendapat tingkat yang sama - tetapi pergantian pemegang jabatan
   cukup lewat baris tabel.

#### Akibat

1. Nol nama orang di dalam kode maupun di dalam teks SQL.
2. Test wajib mencakup jalur yang **ditolak** karena wewenang, bukan hanya jalur yang berhasil.
3. Modul baru memakai aturan peran yang sudah ada. Bila ia menuntut peran baru, peran itu
   ditambahkan ke aturan bersama - bukan dibuat sendiri di dalam modul.

## ADR-U-0031 - Penghapusan adalah penanda dan nilai balik, bukan hapus fisik

`[keputusan work owner]` Baris yang "dihapus" oleh pengguna **tetap tersimpan**. Yang berubah
hanya **penandanya**, dan - bila baris itu membawa nilai - ditambahkan **nilai balik** yang
menolkan pengaruhnya.

**Tidak pernah hapus fisik.**

#### Bentuknya

| Peristiwa | Penanda | Nilai |
| --- | --- | --- |
| nilai diubah | penanda ubah | selisih terhadap nilai lama |
| baris baru | penanda baru | penuh, tanpa pengurang |
| baris ditandai keluar | penanda hapus | **pengurang penuh** |
| seluruh objek dibatalkan | penanda batal | **pengurang penuh untuk setiap baris** |

#### Kenapa

1. Baris tetap **dapat diperiksa**. Yang hilang secara fisik tidak dapat dijelaskan kemudian hari.
2. Jumlah akhir tetap benar tanpa jalur perhitungan khusus - nilai balik menolkannya dengan
   sendirinya.
3. Sistem lama pun tidak menghapus. Perilaku ini **ditiru**, bukan diciptakan.

#### Hubungannya dengan pembatalan generasi

ADR-U-0024 menetapkan pembatalan sebagai generasi bernilai nol. Catatan ini adalah bentuk yang sama
pada tingkat **baris**, bukan tingkat objek. Keduanya sejalan: tidak ada yang hilang, yang berubah
hanya penanda dan nilainya.

#### Akibat

1. Nol perintah hapus fisik pada jalur pengguna di lapisan mana pun.
2. Pembaca hilir **wajib** menyaring pada penanda. Pembacaan tanpa penyaring akan menghitung baris
   yang sudah ditandai keluar.
3. Test yang berhasil menghapus baris secara fisik lewat jalur pengguna **gagal**.

## ADR-U-0032 - Nama yang menyesatkan dibetulkan, disertai tabel pemetaan nama lama ke nama benar

`[keputusan work owner]` Kolom hasil pembacaan **diberi nama sesuai isinya**. Nama lama yang
menyesatkan tidak diwariskan.

Setiap pembetulan **wajib disertai tabel pemetaan** nama lama ke nama benar.

#### Masalah yang ditemukan

`[terverifikasi]` Dua rule pembaca riwayat memberi nama samaran pada kolom hasilnya, dan
**memberi nama samaran dengan pemetaan yang berbeda satu sama lain**. Tanggal kejadian dinamai
*tanggal mulai*; nomor klaim dinamai *nama cabang*; penyebab kerugian dinamai *kode bisnis* - dan
pada rule kedua, pemetaannya lain lagi.

`[terverifikasi]` Pemanggilnya satu-satu, dan **nol pemanggil memakai keduanya**, sehingga kedua
pemetaan yang menyesatkan itu tidak pernah bertemu dalam satu jalur. Tetapi **keduanya membaca
tabel yang sama**.

#### Kenapa tabel pemetaan itu wajib

Tanpa tabel pemetaan, orang yang membandingkan keluaran lama dan keluaran baru akan mengira
**datanya berubah**, padahal hanya **namanya** yang dibetulkan.

Itu kekeliruan yang mahal: ia memicu penyelidikan atas kerusakan yang tidak pernah terjadi, dan
pada kasus terburuk memicu pembatalan migrasi yang sebenarnya benar.

#### Akibat

1. Nama kolom di batas pembacaan mengikuti isinya, bukan mengikuti nama samaran lama.
2. Tabel pemetaan disimpan bersama spec modul yang bersangkutan, bukan di dalam kode.
3. `[penyimpangan sadar]` Pembetulan nama adalah penyimpangan yang disadari dari perilaku lama, dan
   dicatat sebagai penyimpangan.
4. Berapa tepatnya nama samaran yang menyesatkan **belum tentu diketahui seluruhnya**. Setiap yang
   ditemukan kemudian ditambahkan ke tabel pemetaan yang sama.

## ADR-U-0033 - Nama skema ditulis eksplisit pada setiap query

`[keputusan work owner]` Setiap query menyebut **nama skema secara eksplisit**. Tidak ada query
yang bergantung pada skema bawaan sesi.

#### Kenapa

`[terverifikasi]` Basis data yang dipakai memuat **lebih dari satu skema**, dan sebagian besar
query di sistem lama menulis nama tabel **tanpa awalan skema**. Query semacam itu benar hanya
selama sesi kebetulan menunjuk skema yang tepat.

Ketergantungan itu tidak terlihat di teks query. Ia baru muncul sebagai kesalahan ketika
penyetelan sambungan berubah - biasanya saat pindah lingkungan, dan biasanya jauh dari orang yang
menulisnya.

#### Akibat

1. Nol nama tabel tanpa awalan skema di dalam teks query.
2. Nama skema adalah **konfigurasi**, bukan tulisan tetap yang tersebar - sehingga satu tempat
   yang berubah bila skema dipindah.
3. Test yang menemukan query tanpa awalan skema **gagal**.

## ADR-U-0034 - Uang melintasi batas stored procedure sebagai teks, dan dikembalikan ke desimal di dalam aplikasi

`[data DBA]` Sebagian stored procedure penulis menerima **seluruh** parameternya bertipe teks,
**termasuk kolom uang**.

ADR-U-0003 dan ADR-U-0016 tetap berlaku. Yang berubah hanya **di mana** bentuk desimalnya hidup:

| Tempat | Bentuk uang |
| --- | --- |
| Di dalam aplikasi | **desimal berskala tetap** |
| Di kolom basis data | **desimal berskala tetap** |
| Melintasi batas procedure | **teks**, karena tanda tangannya menuntut demikian |

#### Aturannya

1. Perubahan bentuk terjadi **tepat di batas pemanggilan**, bukan lebih awal. Nilai tidak disimpan
   sebagai teks di dalam aplikasi hanya karena nanti akan dikirim sebagai teks.
2. Perubahan bentuk itu **satu fungsi**, dipakai seluruh pemanggil - bukan diulang per pemanggil,
   sebab format pemisah desimal yang berbeda-beda adalah cara paling mudah merusak angka uang.
3. Nilai yang kembali dari procedure diubah balik ke desimal **sebelum** dipakai untuk apa pun.

#### Akibat

1. Nol tipe pecahan biner di jalur mana pun, termasuk saat melintasi batas procedure.
2. Test memeriksa perjalanan pulang-pergi satu nilai uang lewat procedure: nilai yang kembali sama
   persis dengan yang dikirim.

## ADR-U-0035 - Tetapan operasional dibaca dari tabel, bukan ditanam di kode

Tanggal, ambang, dan tetapan operasional lain **dibaca dari tabel** yang dipelihara pengguna -
bukan ditulis tetap di dalam kode atau di dalam teks query.

`[terverifikasi]` Contoh yang menjadi dasar keputusan: tanggal tutup buku dibaca dari tabel
tetapan. Bila tanggal berjalan melewatinya, cap waktu produksi digeser ke awal bulan berikutnya.

#### Kenapa

1. Tetapan semacam ini **berubah menurut jadwal dagang**, bukan menurut jadwal rilis. Menanamnya
   di kode berarti setiap pergeseran tanggal menuntut penerbitan versi baru.
2. Orang yang tahu nilainya adalah pengguna, bukan pengembang.

#### Batasnya

Ini **bukan** izin memindahkan aturan dagang ke tabel. Yang dibaca dari tabel adalah **nilainya**;
aturan yang memakai nilai itu tetap hidup di kode, tempat ia dapat diuji.

`[terverifikasi]` Di sistem lama ditemukan juga tetapan yang ditanam sebagai nama orang dan nomor
tertentu di dalam teks query. Semua yang seperti itu **dibuang** - lihat ADR-U-0030.

#### Akibat

1. Nol tanggal, ambang, atau nama tetap di dalam kode maupun teks query.
2. Perilaku saat tetapan **belum terisi** ditetapkan eksplisit, bukan dibiarkan menjadi galat yang
   tidak terbaca.
3. Test memakai nilai tetapan yang disuntikkan, bukan nilai produksi.

## ADR-U-0036 - Medan kosong tidak hadir di dokumen JSON, dan ketidakhadirannya bukan kegagalan

`[terverifikasi]` Dokumen JSON sistem lama **tidak memuat medan yang nilainya kosong**. Medan itu
bukan bernilai kosong - ia **tidak ada sama sekali**.

`[keputusan work owner]` Pemuat migrasi membaca setiap medan dari dokumen dan menghasilkan
**kosong bila medan absen**. Ketidakhadiran **bukan kegagalan**.

#### Kenapa ini perlu ditulis

Pemuat yang menganggap medan absen sebagai kesalahan akan menolak sebagian besar dokumen lama -
bukan karena dokumennya rusak, melainkan karena pengisiannya memang tidak lengkap sejak awal.

Sebaliknya, pemuat yang diam-diam melewati medan absen tanpa mencatat akan menyamarkan dokumen
yang benar-benar rusak.

#### Aturannya

1. Medan absen menghasilkan nilai kosong, dan pemuatan **dilanjutkan**.
2. Medan yang hadir tetapi tidak dikenali masuk **penampung medan tak dikenal** (ADR-U-0023). Itu
   persoalan yang berbeda dan tidak boleh dicampur.
3. Kewajiban isi tidak ditegakkan oleh pemuat migrasi. Data lama masuk apa adanya; kewajiban isi
   berlaku bagi masukan baru (ADR-U-0027).

#### Akibat

1. Test memuat dokumen yang kehilangan sebagian medan, dan memastikan pemuatan berhasil dengan
   kolom kosong - bukan gagal.
2. Cacah medan yang absen dicatat sebagai keterangan pemuatan, bukan sebagai galat.

## ADR-U-0037 - Kolom dipilih dari medan yang benar-benar diisi, bukan dari yang sekadar tampil

`[keputusan work owner]` Daftar kolom disusun dari **medan yang di-set oleh aturan penyimpan**.
Itulah data yang benar-benar terisi.

Medan yang **hanya** muncul di layar dengan **kondisi tampil yang selalu salah** - sisa rancangan
yang tidak pernah hidup - **tidak menjadi kolom**.

#### Hubungannya dengan ADR-U-0017

ADR-U-0017 menetapkan daftar medan disusun dari sapuan aturan, mencakup Activity **dan** Section.
Catatan ini menyaring hasil sapuan itu, dan keduanya sejalan:

| Langkah | Sumber | Gunanya |
| ---: | --- | --- |
| 1 | sapuan Activity **dan** Section | menemukan **calon** medan - jangan ada yang terlewat |
| 2 | penyaringan di catatan ini | menentukan calon mana yang **menjadi kolom** |

Sapuan yang luas dan penyaringan yang ketat bukan hal yang bertentangan. Melewatkan langkah 1
membuat medan hilang; melewatkan langkah 2 membuat tabel penuh kolom yang tidak pernah terisi.

#### Akibat

1. Medan yang gugur di langkah 2 **dicatat beserta alasannya**, bukan dihilangkan diam-diam -
   supaya keputusan itu dapat diperiksa ulang bila kemudian ternyata keliru.
2. Medan yang terbukti diisi aturan penyimpan menjadi kolom, meskipun tidak tampil di layar mana
   pun.
3. Test membaca dan menulis kembali dokumen nyata, dan memastikan nol medan terisi yang hilang.

## ADR-U-0038 - Tangga persetujuan berakhir karena keadaan jenjang, bukan karena pencacah

`[keputusan work owner]` Tangga persetujuan berakhir bila **seluruh jenjang menyetujui**. Satu
menolak, kasus **tutup seketika**.

Syarat berakhir itu **dinyatakan eksplisit** dan dapat diuji.

#### Apa yang diganti

`[terverifikasi]` Di sistem lama penugasan berhenti karena penentu penerima **tidak menghasilkan
siapa pun** ketika setiap jenjang sudah memutuskan. Itu **perilaku diam**: benar dalam praktik,
tetapi tidak dapat diuji dan tidak dapat dijelaskan.

`[penyimpangan sadar]` Sistem baru menyatakan syaratnya, bukan menyimpulkannya dari ketiadaan
penerima.

#### Pencacah jenjang bukan kebenaran

Pencacah jumlah jenjang **boleh ada sebagai tampilan**, tetapi **bukan** dasar keputusan.

`[terverifikasi]` Alasannya nyata: daftar jenjang **dapat berubah di awal** - pengaju dikeluarkan
dari daftar - dan pencacah tidak ikut tahu. Tangga yang berhenti berdasarkan pencacah akan
berhenti di tempat yang salah persis pada kasus yang paling perlu benar.

#### Akibat

1. Keputusan berakhir diambil dari **keadaan setiap jenjang**, bukan dari cacah.
2. Test wajib mencakup: daftar jenjang yang menyusut di awal, dan tangga yang tetap berakhir pada
   jenjang yang benar.
3. Penolakan menutup kasus seketika - bukan melanjutkan ke jenjang berikutnya.

## ADR-U-0039 - Rujuk atau salin: identitas dirujuk, jabatan disalin, data kutipan disalin utuh

`[keputusan work owner]` Aturan pembeda yang mengikat seluruh pencatatan:

> **Identitas orang DIRUJUK. Jabatan DISALIN. Tulisan jabatan DIRUJUK.**

| Yang dicatat | Caranya | Alasannya |
| --- | --- | --- |
| Akun orang | **rujukan** | orangnya tetap orang yang sama; bila namanya berubah, catatan **memang seharusnya** ikut nama baru |
| Jabatan saat itu | **salinan** | jabatan saat itu adalah **fakta sejarah** - tidak boleh ikut naik ketika orangnya naik jabatan |
| Tulisan jabatan | **rujukan** ke daftar induk | bila label yang sama ditulis ulang, catatan lama boleh ikut; yang dilarang adalah catatan lama **berpindah ke jabatan yang berbeda** |

#### Data kutipan disalin utuh, bukan medan bernama satu per satu

`[keputusan work owner]` `[penyimpangan sadar]` Seluruh data kutipan yang dibutuhkan disalin dari
kasus induk - **bukan** daftar medan yang disebut satu per satu.

`[terverifikasi]` Sebab daftar bernama itulah yang melahirkan cacat nyata: objek kerja hanya
menerima **dua** medan, sedangkan **sembilan** penggolong menguji medan lain. Dua di antaranya
dipakai hidup lima kali, sehingga dua cabang lini usaha **tidak pernah terbit sama sekali**.

Ini dicatat sebagai **cacat yang diperbaiki**, bukan perilaku yang ditiru.

#### Akibat

1. Penyalinan data kutipan bersifat menyeluruh. Menambah medan baru di kasus induk tidak menuntut
   perubahan di sisi penyalin.
2. Kolom jabatan pada catatan sejarah tidak pernah diperbarui belakangan.
3. Test wajib: orang naik jabatan, dan catatan lama **tetap** menunjukkan jabatan lamanya.

## ADR-U-0040 - Pengaju dikeluarkan dari jenjang pertama - dan lubang di jenjang berikutnya diterima dengan sadar

`[keputusan work owner]` Saat kasus persetujuan dibentuk, bila **pemegang jenjang pertama adalah
pengaju**, ia **dikeluarkan dari daftar** dan jumlah jenjang dihitung ulang.

#### Lubang yang dibawa masuk dengan sadar

`[penyimpangan sadar]` Keputusan ini **berbeda dari rekomendasi asisten**, dan perbedaannya
disengaja. Akibatnya ditulis terang-terangan, bukan disamarkan:

> **Pemegang jenjang kedua ke bawah tetap dapat menyetujui pekerjaan yang ia ajukan sendiri.**

Pada modul saudaranya, **tidak ada pemeriksaan sama sekali** - sehingga lubang itu berlaku pada
jenjang mana pun.

Ini **lubang yang dibawa masuk dengan sadar**, bukan kelalaian. Ia dicatat di sini supaya siapa
pun yang menemukannya kemudian tahu bahwa ia sudah dilihat, ditimbang, dan diterima - dan supaya
keputusan untuk menutupnya dapat diambil sebagai keputusan tersendiri.

#### Akibat

1. Pemeriksaan pengaju hanya berlaku pada **jenjang pertama**, bukan seluruh jenjang.
2. Jumlah jenjang dihitung **sesudah** pengaju dikeluarkan - lihat ADR-U-0038 tentang mengapa
   pencacah tidak boleh menjadi dasar keputusan.
3. Test wajib mencakup kasus pengaju adalah pemegang jenjang pertama, dan memastikan jumlah
   jenjang menyusut dengan benar.

## ADR-U-0041 - Endorsemen Life berbagi tabel dengan new business dan dibedakan kolom versi

`[keputusan work owner]` Endorsemen **tidak punya tabel sendiri**. Ia berbagi tabel dengan new
business, dan dibedakan oleh **kolom pembeda versi**.

Pencocokan baris antar versi lewat **penunjuk ke baris induk**, bukan lewat urutan indeks.

#### Kenapa bukan tabel terpisah

1. Bentuk datanya sama. Tabel terpisah berarti bentuk yang sama dipelihara di dua tempat.
2. Pembaca hilir membaca satu tempat, bukan menggabungkan dua tabel yang bentuknya kebetulan sama.
3. Rantai versi menjadi hubungan di dalam satu tabel, bukan jembatan antar tabel.

#### Kenapa penunjuk induk, bukan indeks

`[penyimpangan sadar]` Urutan baris **dapat berubah** antar versi. Pencocokan lewat indeks putus
persis ketika urutannya bergeser - dan pergeseran itu tidak menimbulkan galat apa pun, hanya angka
yang salah pasang.

#### Hubungannya dengan modul lain

ADR-U-0018 menetapkan bentuk rantai generasi lewat penunjuk, dan ADR-U-0019 menetapkan pemasangan
baris anak lewat nomor urut. Catatan ini adalah bentuk yang **berbeda** untuk persoalan yang sama,
pada modul yang berbeda.

Perbedaan itu **disengaja dan dicatat**: modul ini berbagi tabel dengan new business, sedangkan
modul treaty punya tabel sendiri untuk tiap generasi. Menyeragamkan keduanya sesudah keduanya
berjalan akan lebih mahal daripada mencatat perbedaannya di sini.

#### Data lama tidak disimpan sebagai salinan

`[keputusan work owner]` Halaman kerja yang menampilkan keadaan sebelum dan sesudah **tidak
di-persist**. Data lama dibaca dari **versi sebelumnya** lewat penunjuk induk.

#### Akibat

1. Nol tabel salinan data lama.
2. Setiap pembacaan wajib menyaring pada kolom pembeda versi.
3. Test wajib: urutan baris bergeser antar versi, dan pasangannya **tetap benar**.

## ADR-U-0042 - Skema relasional dirancang baru bila ia tidak pernah ada

`[terverifikasi]` Pada sebagian modul, skema relasional **tidak dapat direkayasa-balik dari
korpus - sebab ia belum pernah ada**. Persistensi terjadi lewat penyimpanan objek kerja sebagai
dokumen, bukan lewat tabel.

`[keputusan work owner]` Untuk modul semacam itu, skema **dirancang baru**, dan perancangannya
dinyatakan sebagai pekerjaan tersendiri - bukan disamarkan sebagai hasil penelusuran.

#### Bahan perancangan yang sah

| Sumber | Yang diambil darinya |
| --- | --- |
| Bentuk dokumen nyata `[data DBA]` | nama dan tipe medan |
| Hierarki daftar bersarang pada layar | bentuk induk-anak |
| Sensus kolom tabel proyeksi yang sudah ada | medan yang terbukti dipakai hilir |

#### Kenapa ini perlu dinyatakan

Tanpa pernyataan ini, pembaca spec akan mengira bentuk tabel **ditemukan** di sistem lama, dan
memperlakukannya sebagai fakta yang tidak boleh diganggu.

Kenyataannya ia **rancangan**, dan rancangan boleh diperdebatkan dengan alasan. Menyamarkan
rancangan sebagai temuan menghilangkan hak orang berikutnya untuk memperbaikinya.

#### Akibat

1. Modul yang skemanya dirancang baru menandainya di spec, dengan bukti bahwa skema lama memang
   tidak ada.
2. Tabel proyeksi yang sudah dipakai sistem hilir **diikuti apa adanya** - itu kontrak luar, bukan
   rancangan bebas.
3. Setiap tabel hasil rancangan baru menyebut **dari mana** daftar kolomnya berasal (ADR-U-0017,
   ADR-U-0037).

---

# Bagian II - Keputusan per modul - rangkaian kedua

Jumlah: **56** keputusan.

## Modul claim-non-prop

### ADR-D-CNP-0001 - Aggregate root adalah Klaim, satu klaim satu kejadian kerugian

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Struktur XML mendukung tafsir ini tetapi tidak memaksakannya: seluruh properti tanggal dan penyebab kerugian di `.ClaimData` bersifat skalar, ada 18 PageList di bawahnya, dan sapuan atas 114 activity menunjukkan **tidak ada satu pun rule yang menulis `.DateOfLoss`** — nilainya hanya masuk lewat input Registrasi. Kami menetapkan satu Klaim = satu kejadian kerugian atas satu polis treaty non-proporsional, dengan Adjustment sebagai transaksi di bawahnya, sehingga Klaim menjadi aggregate root dan ke-18 PageList menjadi tabel anak.

##### Consequences

Natural key Klaim belum ditetapkan. Tidak ada rule yang mencegah duplikat `PolicyNo` + `DateOfLoss` + `IDMaster`; `CheckDateDOL_Act` hanya menampilkan riwayat klaim atas polis yang sama sebagai informasi, tidak menolak. Sampai ada konfirmasi bahwa pencegahan duplikat dilakukan secara prosedural, kunci teknis tetap identitas klaim yang digenerasi sistem, dan kombinasi di atas diperlakukan sebagai indeks pendukung, bukan unique constraint.

### ADR-D-CNP-0002 - Klaim yang sudah ditutup tidak dapat dibuka kembali

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Di seluruh 279 berkas folder ini hanya ada satu status akhir, `Resolved-Completed`, dan tidak ditemukan satu pun rule yang mengembalikan klaim dari status akhir ke status terbuka. Kami menetapkan tidak ada mekanisme reopen: koreksi atas klaim yang sudah ditutup dilakukan melalui Adjustment baru, bukan dengan membuka ulang klaim.

##### Consequences

Kesimpulan ini **disimpulkan dari ketiadaan rule, bukan dari adanya rule**. Ketiadaan jalur reopen di XML tidak membuktikan ketiadaan reopen di produksi — pembukaan ulang secara manual lewat database tidak meninggalkan jejak di rule. Verifikasinya tercatat sebagai EXTERNAL di `_selesai/OPEN-QUESTIONS.md`. Bila ternyata reopen manual memang dilakukan, keputusan ini harus ditinjau ulang sebelum model status dikunci, karena konsekuensinya menyentuh keterlacakan audit.

### ADR-D-CNP-0003 - Presisi tinggi sepanjang rantai perhitungan, pembulatan hanya di tepi

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Sistem lama tidak punya aturan pembulatan: dari 673 ekspresi aritmetika, **518 (77%) tidak menyatakan skala sama sekali** dan bergantung pada perilaku desimal bawaan Pega, sementara 155 sisanya memakai tujuh skala berbeda. Kami menetapkan aturan baru: nilai moneter dihitung dan disimpan pada satu presisi tinggi seragam sepanjang rantai perhitungan, dan pembulatan hanya dilakukan di dua tepi — saat ditampilkan ke pengguna, dan saat dikirim ke sistem hilir.

##### Considered Options

Menyalin skala per ekspresi apa adanya ditolak: skalanya tidak konsisten bahkan di dalam satu rule (`CountClaimTNP_Act` memakai 10, 20, dan 5; `CountLossAllocation_act` memakai 10 dan 20 dalam porsi berimbang), sehingga menyalinnya berarti mengabadikan ketidaksengajaan.

##### Consequences

Skala ternyata **bukan semata kebiasaan penulisan**, melainkan mengikuti batas modul: inti klaim non-proporsional memakai 20, sisi treaty/premi memakai 4, pemformatan tampilan memakai 0 dan 2 (selalu lewat "bagi dengan 1"), dan tarif pajak di `SetPPNPPH` memakai 8. Dugaan awal bahwa ini murni kebiasaan developer meleset dan sudah dikoreksi.

Angka presisi final **belum ditetapkan** dan ADR ini tidak boleh menyebut angka apa pun sebagai final sebelum hasil profil data Oracle diterima. Risiko yang harus diputuskan lebih dulu: bila profil menunjukkan kolom uang di produksi benar-benar menyimpan lebih dari dua desimal, maka aturan "bulatkan dua desimal saat kirim ke hilir" justru akan mengubah angka yang selama ini diterima akuntansi.

##### Catatan 18 September 2026 — tetap `proposed`, dan alasannya berubah

Status **tidak** dinaikkan. Dua hal yang perlu tercatat:

1. **`NUMBER(20,4)` pada `TREATYINPRODUCTION` diperlakukan sebagai konvensi satu modul**, bukan standar perusahaan, sampai akuntansi menyatakan sebaliknya. Ia **tidak** disebarkan ke tabel lain dan tidak dipakai sebagai dasar keputusan presisi.
2. **Alasan ADR ini masih draft sudah berubah.** Semula: datanya belum ditarik. Sekarang: **basis data memang tidak menyimpan presisi untuk nilai klaim sama sekali** — nilainya tidak punya kolom, tersimpan sebagai teks di dalam JSON, dan bahkan view `CLAIMXOL` mengeluarkannya sebagai `varchar2`. Profil data tidak akan mengubah kenyataan itu.

Yang menghambat sekarang adalah keputusan kebijakan, bukan ketersediaan data — `ASK-AKUNTANSI.md` pertanyaan 1.

##### Naik ke `accepted` — 18 September 2026

Status naik lewat **AK-1** di `REGISTER-RATIFIKASI.md`, dan **bukan** karena data profil akhirnya datang — catatan di atas sudah menjelaskan bahwa data itu tidak akan pernah menjawabnya.

Yang berubah: angka presisinya kini diturunkan dari **skala yang dipakai sistem lama itu sendiri**, bukan dari profil kolom yang tidak ada. Skala 20 pada 126 ekspresi di seluruh rule inti klaim non-proporsional (`BLUEPRINT.md` §6.3).

| Kelompok | Tipe |
|---|---|
| U1 nilai uang antara · P1 persentase dan porsi · K1 kurs · L1 limit dan premi deposit | `NUMBER(38,20)` |
| P2 tarif pajak dan brokerage | `NUMBER(11,8)` |
| C1 cacah dan nomor urut | `NUMBER(9)` |

**Tidak ada kelompok untuk "nilai uang final".** Pembulatan dua desimal terjadi hanya di tepi — lapisan view kompatibilitas dan arsip muatan keluar. Alasannya `TempKasir.CARI9 = .TotalClaim - .PremiumSpreaded`: `.TotalClaim` adalah nilai antara sekaligus bahan instruksi bayar, dan satu kolom tidak dapat berskala 20 dan 2 sekaligus.

Ratifikasi akuntansi dicatat di `REGISTER-RATIFIKASI.md` dan **tidak menahan pekerjaan**.

##### Bukti yang menguatkan — sapuan S10, 18 September 2026

Ditambahkan sesudah ADR ini `accepted`. Ia **menguatkan**, tidak menutup apa pun dan tidak mengubah putusan.

Dua rantai nilai uang ditelusuri dari lahir sampai keluar, bukan dicacah per ekspresi:

| Rantai | Skala yang dinyatakan | Berkas yang dilewati |
|---|---|---|
| `.GrossValue` | **tidak ada, di satu titik pun** | `CountLossAllocation_act`, `CreateChildKomiteCNP_Act`, `SaveCNPLayerList_Act`, `SaveDataToOSAksep_Act`, `SaveToOS`, `GetSelisihActual_Act` |
| `.AdjusterFeeValue` | **20**, di setiap pembagian | rule yang sama, sebagian besar |
| `Local.URLimit` | **10** | `CountLossAllocation_act` 6038, 6065 |

**Tiga skala dalam satu modul, dua di antaranya di berkas yang sama.** Angka 77% di badan ADR ini adalah statistik; ini contoh berjalan: sebuah nilai uang yang benar-benar dikirim ke sistem hilir menempuh enam berkas tanpa satu pun pernyataan skala, sementara nilai di sebelahnya menyatakannya setiap kali.

Menguatkan pernyataan *"skala mengikuti batas modul"* di bagian Consequences dengan satu koreksi halus: pada rantai ini skala mengikuti **besaran yang dihitung**, bukan rule yang menghitungnya — `CountLossAllocation_act` memakai 10 dan 20 untuk besaran berbeda di dalam dirinya sendiri.

Tercatat sebagai **E1**, **E2**, **E12**, **B1**, **B2**, dan **D24**. Rinci di `SAPUAN-S10-S16.md` §3.

##### Penutupan B1, B2, E1, E2, E12 — 19 September 2026

**Catatan status**: ADR ini **sudah** `accepted` sejak 18 September 2026 lewat AK-1; ia tidak dinaikkan lagi hari ini. Yang ditambahkan bagian ini adalah **penutupan lima butir register** yang selama ini masih menunggunya, beserta alasan mengapa menunggu lebih lama tidak akan mengubah apa pun.

###### Dasarnya sudah lengkap dan tidak akan berubah oleh data

Empat lapisan diperiksa, dan keempatnya sepakat:

| Lapisan | Keadaan |
|---|---|
| Tipe kolom | **88 dari 120** kolom `NUMBER` di DDL tanpa presisi maupun skala |
| Muatan JSON dan blob | nilai uang tersimpan sebagai **teks**; view `CLAIMXOL` mengeluarkannya sebagai `varchar2` |
| Ekspresi Pega | **77%** tanpa skala; rantai `.GrossValue` menempuh enam berkas tanpa satu pun pernyataan |
| Java tertanam | **nol `double`, nol `float`, nol `BigDecimal`**; uang melintas sebagai `String` 92 kali lewat `.toString()` (**D39**) |

Lapisan keempat adalah yang menutup pertanyaannya. Selama hanya tiga lapisan terbaca, masih mungkin berharap kendali skala hidup di tempat yang belum dilihat. Java adalah tempat terakhir itu, dan ia **tidak menghitung uang sama sekali** — ia memindahkannya sebagai teks.

> **Tidak ada presisi yang dapat dipulihkan, karena tidak pernah ada presisi yang ditegakkan.**

Profil data Oracle tidak dapat mengubah kalimat itu. Karena itu **B1 dan B2 tidak lagi menunggu REQ-001**.

###### Maka presisi adalah keputusan maju, bukan temuan

1. **Skala kanonik 20 sepanjang rantai hitung.** Angka ini bukan selera: ia **skala tertinggi yang benar-benar dipakai sistem lama** — inti klaim non-proporsional, dan setiap pembagian di rantai `.AdjusterFeeValue`. Memakainya karena itu **tidak dapat kehilangan apa pun yang pernah ada**.
2. **Pembulatan hanya di tepi** — saat ditampilkan, saat dikirim ke hilir, saat dibukukan. **Tidak di tengah rantai.**
3. **Migrasi mengurai teks apa adanya** dan menyimpan seluruh digit yang ada. Nilai yang **tidak terurai tidak dibulatkan dan tidak dibuang** — ia masuk jalur tiket `14` dan tercatat beserta asalnya.
4. **Skala 10 dan 4 yang terbaca di sistem lama dicatat sebagai perilaku lama di tabel perbedaan, bukan ditiru.** `Local.URLimit` memakai 10; `TREATYINPRODUCTION` memakai `NUMBER(20,4)`. Keduanya fakta tentang sistem lama, bukan standar yang diwarisi.

###### Apa yang ini tutup, dan apa yang tetap berdiri

| Butir | Isinya | Penutupnya |
|---|---|---|
| **B1** | presisi dan skala final tiap kolom nilai | butir 1 di atas |
| **B2** | presisi tidak dijaga di tiga (kini empat) lapis | tidak ada yang dapat dipulihkan; keputusan maju menggantikannya |
| **E1** | dua rule mencampur skala pembagian | skala mengikuti besaran, bukan rule — tidak ditiru |
| **E2** | 77% ekspresi tanpa skala | sama |
| **E12** | tiga skala dalam satu modul | sama |

**REQ-001 tetap diminta, dan kedudukannya berubah**: ia **memverifikasi** apa yang benar-benar tersimpan di blob, dan **tidak lagi menahan**. Bila jawabannya kelak memperlihatkan digit yang lebih banyak daripada dugaan, skala 20 tetap menampungnya; bila lebih sedikit, tidak ada yang hilang.

### ADR-D-CNP-0004 - Tambalan per-case tidak ikut dimigrasi; angkanya dipindahkan sebagai data

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Sapuan seluruh folder menemukan **9 identitas klaim atau master treaty yang ditulis langsung di dalam kode kalkulasi**, tersebar di 4 rule dan 14 langkah — antara lain `CLMNP-975` yang menimpa Premi Pemulihan dengan angka mati `881928.966808370` (IDR) dan `806851161.1895` (USD), serta `CLMNP-232` yang memetakan Biaya Penilaian secara manual per indeks. Kami memutuskan tidak satu pun dari cabang ini dibawa ke sistem baru: kodenya dibuang, tetapi nilai akhirnya dipertahankan sebagai data sehingga saldo klaim yang bersangkutan tetap benar.

##### Consequences

Daftar lengkap ke-9 identitas beserta nilai yang ditimpa dan rumus normal yang di-bypass ada di `BLUEPRINT.md` §7, dan menjadi instruksi bagi tim migrasi data — bukan bagi tim pembangun aplikasi.

Salah satu temuan memperkuat bahwa ini memang tambalan, bukan aturan: pada `InputOutStandingClmTNP_PreAct` langkah 5, deskripsi langkah berbunyi `pyWorkPage.pyID=="CLMNP-50"` sementara kondisi eksekusinya `pyWorkPage.pyID=="CLMNP-232"` — deskripsi dan kondisi tidak sinkron, ciri khas salin-tempel.

Pemeriksaan ulang terhadap `MEMORI_PEMAHAMAN.MD` tidak menemukan keterangan bahwa treaty `1000393` punya perlakuan bisnis khusus. Karena itu ia diperlakukan sebagai tambalan, bukan sebagai dimensi yang hilang dari model data. Bila pemilik proses kemudian menyatakan sebaliknya, keputusan ini harus ditinjau ulang dan model data perlu menambah atribut yang menjelaskan perbedaan perlakuan tersebut.

##### Tambahan 18 September 2026 — persetujuan eksplisit dan mekanisme pengganti

Persetujuan diberikan tanpa menunggu status terbuka/tutup tiap klaim: **tidak satu pun dari 29 langkah tambalan dimigrasi.**

Penggantinya satu mekanisme tunggal: **koreksi bernilai tercatat** — `nilai_sebelum`, `nilai_sesudah`, `alasan`, `pelaku`, `waktu`. Satu bentuk untuk semua kasus, bukan cabang per klaim di dalam kode.

Status tiap klaim yang disebut di dalam tambalan (`CLMNP-232`, `CLMNP-975`, `CLMNP-861`, `CLMNP-50`, `CLMNP-367`, `CLMNP-382`, `IDMaster 1000393`, `1001130`) diukur lewat REQ-015, dan itu **bukan penghalang** keputusan ini.

Satu hal yang dapat disimpulkan tanpa menunggu query: `CLMNP-975` ditambal **2026-07-16**, dua bulan sebelum ekspor ini. Tambalan tidak ditulis untuk klaim yang sudah tutup, jadi klaim itu hampir pasti masih terbuka dan tambalannya masih aktif.

**Karena itu `CLMNP-975` dijadikan kasus uji utama shadow-run.** Bila sistem baru menghasilkan angka yang benar untuk klaim itu **tanpa** tambalan apa pun, tesis "ini tambalan data, bukan aturan bisnis" terbukti untuk seluruh kelompoknya sekaligus — dan 29 langkah itu gugur bersama-sama, bukan satu per satu.

### ADR-D-CNP-0005 - Claim dan Komite satu unit cutover, dengan shadow-run sebagai bukti paritas

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Claim Non Prop dan Komite Claim Non Prop tidak dapat dipisahkan pada saat peralihan: keduanya menulis ke record yang sama dan saling mengunci (`pxAddChildWork` dari sisi Claim, penulisan balik ke `AdjustmentList` dari sisi Komite), sehingga memotong di antara keduanya berarti menciptakan penguncian lintas sistem atas baris Oracle yang sama. Kami menetapkan keduanya beralih sebagai satu unit; sebelum peralihan dijalankan shadow-run, yaitu sistem baru menghitung ulang klaim produksi yang sudah selesai dan hasilnya dibandingkan angka per angka dengan sistem lama; klaim yang sedang berjalan diselesaikan di sistem lama dan sistem baru hanya melayani klaim baru.

##### Consequences

**Tidak ada satu pun bukti di XML yang mendukung maupun menolak strategi ini** — XML hanya membuktikan keterkopelannya, bukan cara memindahkannya. Ini murni keputusan kami.

Yang terbukti dari XML adalah titik koplingnya, dan itu tercatat lengkap di `BLUEPRINT.md` §8 sebagai Boundary Contract Draft: dua pemanggilan `pxAddChildWork`, kunci relasi `pxCoveredInsKeys`/`pxCoverInsKey`, dan `Adjustment.IndexObject` sebagai penunjuk posisi numerik ke `AdjustmentList` induk.

`Adjustment.IndexObject` adalah risiko tersendiri: ia menyimpan posisi dalam daftar, bukan kunci surrogate. Bila urutan `AdjustmentList` berubah, kaitan antara keputusan Komite dan Adjustment yang dimaksud menjadi salah. Sistem baru wajib mengganti ini dengan kunci yang stabil, dan migrasi data harus memetakan posisi lama ke kunci baru.

##### Tambahan 18 September 2026 — prasyarat baseline shadow-run

Shadow-run mengandaikan baseline yang sehat. Dua hal dapat merusaknya, dan keduanya harus diselesaikan **sebelum** perbandingan dimulai:

1. **Baris ganda akibat hitung ulang yang tidak idempoten** (ADR-D-CNP-0011). Case terdampak dibersihkan atau dikeluarkan dari perbandingan, dan didaftar terpisah. Diukur lewat REQ-014.
2. **Riwayat persetujuan yang tidak dapat direkonstruksi** karena hanya tersimpan sebagai teks komentar (ADR-D-CNP-0009, FINDING-002 bagian 7). Case terdampak tidak dapat dibandingkan pada dimensi "siapa menyetujui". Diukur lewat REQ-013.

Kasus uji utama: **`CLMNP-975`**, lihat ADR-D-CNP-0004.

##### Tambahan 18 September 2026 — komposisi sampel dan aturan kegagalan

###### Komposisi sampel ditetapkan di muka

Sampel shadow-run **ditetapkan sekarang, bukan dipilih saat menjalankan**. Perbandingan atas 500 klaim rupiah satu layer akan lulus dengan mudah dan tidak membuktikan apa pun.

Sampel wajib memuat, dengan proporsi minimum terhadap keseluruhan sampel:

| Kelompok | Minimum |
|---|---|
| Klaim multi-mata-uang (lebih dari satu `Currency` pada satu klaim) | 15% |
| Klaim multi-layer (lebih dari satu baris `SpreadingRisk` selain `"UR"`) | 20% |
| Klaim dengan reinstatement (`CNPReinstatement` bukan nol) | 10% |
| Klaim dengan lebih dari satu Adjustment | 15% |
| Klaim yang melibatkan baris `"UR"` | 20% |
| Kedelapan klaim bertambalan | **100% — seluruhnya wajib ikut** |

Kelompok boleh bertumpang tindih; yang tidak boleh adalah ada kelompok yang kosong. Bila populasi produksi tidak menyediakan cukup kasus untuk satu kelompok, itu sendiri temuan dan dicatat, bukan diam-diam diturunkan ambangnya.

###### Aturan kegagalan

| Jenis nilai | Bila selisih melampaui ambang |
|---|---|
| **Nilai yang masuk jurnal akuntansi** | **Cutover berhenti.** Tanpa negosiasi, tanpa daftar pengecualian. |
| **Nilai antara** (alokasi per layer sebelum pembulatan akhir) | Boleh masuk daftar pengecualian, tetapi **setiap baris wajib punya penjelasan dan persetujuan bernama**. Pengecualian tanpa penjelasan tidak dihitung sebagai pengecualian — ia dihitung sebagai kegagalan. |

Aturan ini ditetapkan **sebelum** ada pihak yang punya kepentingan atas hasilnya. Itu alasan ia ditulis sekarang dan bukan nanti.

Angka toleransi relatif untuk nilai antara masih menunggu akuntansi (lihat ADR-D-CNP-0003 yang berstatus draft).

###### Prasyarat baseline

Selain dua hal yang sudah dicatat di atas, bertambah satu yang baru terbukti dari DDL: `PEGA_JSON_OS_AKSEP_KLAIMTNP` **selalu `INSERT` dan tidak pernah `UPDATE`** — logika upsert-nya ada tetapi dikomentari seluruhnya. Bersama ketiadaan primary key pada `OS_AKSEPTASI_KLAIM`, baris ganda pada tabel itu mungkin terjadi dan harus diukur sebelum perbandingan dimulai.

##### Tambahan — nilai IDR lama tidak dipakai sebagai pembanding

Konversi mata uang di sistem lama **tidak dapat direproduksi**: `POOLDATA.GETCURRENCYSTANDARD` mengabaikan parameter tanggalnya dan selalu memakai kurs terbaru sampai `sysdate` (`FINDING-001` bagian 7.1). Angka IDR yang dihitung tahun lalu tidak dapat dihasilkan ulang hari ini.

Karena itu, untuk nilai hasil konversi, **shadow-run membandingkan nilai mata uang asli, bukan nilai IDR-nya.** Nilai IDR lama diterima apa adanya sebagai nilai historis dan tidak dijadikan pembanding.

Pengecualian ini berlaku hanya untuk kelompok nilai hasil konversi. Nilai yang sejak awal berdenominasi rupiah tetap dibandingkan seperti biasa, dengan aturan dua ambang di atas.

### ADR-D-CNP-0006 - RBAC sistem baru dirancang dari nol, bukan diwarisi

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Modul ini tidak memiliki otorisasi di tingkat rule: `pyPrivilegeName` kosong di seluruh 279 berkas, tidak ada satu pun `Rule-Access-When`, dan seluruh elemen UI bervisibilitas `ALWAYS` — termasuk kedua jalur penutupan klaim. Kontrol akses yang benar-benar berjalan bertumpu pada access group dan routing flow yang tidak ikut ter-export. Kami menetapkan model otorisasi sistem baru dirancang dari nol bersama pemilik proses, dan tidak ada satu pun aturan akses yang boleh diklaim sebagai warisan sistem lama.

##### Consequences

Ini **menambah cakupan pekerjaan**, bukan sekadar catatan. Perancangan RBAC menjadi pekerjaan tersendiri dengan pemilik keputusan di sisi bisnis, bukan turunan otomatis dari analisis XML.

Empat label peran memang terbukti ada di `DataTransform\InsertChronology_DT.xml` — `Claim Admin`, `Claim Dept. Head`, `Operational Director`, `Technical Director` — tetapi itu label yang dicatatkan ke jejak audit, bukan aturan yang menegakkan siapa boleh melakukan apa. Pemetaannya ke orang pun bersifat hardcode per nama operator, sehingga tidak dapat dijadikan dasar model peran.

Sampai RBAC baru disepakati, setiap pernyataan tentang "siapa boleh melakukan apa" dalam dokumen mana pun berstatus hipotesis dan wajib berlabel EXTERNAL.

### ADR-D-CNP-0007 - Setiap nilai uang disimpan berpasangan, dan ambang kewenangan dibandingkan terhadap nilai IDR

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Klaim non-proporsional berjalan dalam banyak mata uang sekaligus: `ListClaimAmount`, `SpreadingRisk`, dan `ListTotalEstimation` semuanya ber-kunci `Currency`, sementara properti bersufiks `IDR` menunjukkan IDR dipakai sebagai mata uang penyetaraan. Kami menetapkan setiap nilai uang disimpan sebagai satu paket — **nilai asli, mata uang, kurs, tanggal/sumber kurs, dan nilai IDR hasil konversi** — dan setiap perbandingan terhadap ambang kewenangan Komite dilakukan terhadap nilai IDR hasil konversi tersebut.

##### Consequences

**Tanggal dan sumber kurs adalah bagian wajib dari paket**, bukan pelengkap. Tanpa keduanya, angka IDR tidak dapat diaudit ulang: nilai yang sama bisa dibenarkan atau disalahkan tergantung kurs kapan yang dipakai, dan rekonsiliasi dengan akuntansi menjadi tidak dapat diselesaikan. Sumber kurs yang terlihat di sistem lama adalah fungsi `getcurrencystandard` dan `TreatyInMaster.CurrencyList.Conversion`; keduanya perlu dicatat identitasnya, bukan hanya hasilnya.

Keputusan tentang ambang ini **tidak menunggu** hasil investigasi kondisi sistem lama yang tercatat di `FINDING-001-threshold-currency.md`. Apa pun temuan di sana, sistem baru membandingkan terhadap IDR.

Presisi dan skala setiap kolom dalam paket ini belum ditetapkan dan menunggu `pengetahuan/SCHEMA-ACTUAL.csv` — lihat ADR-D-CNP-0003 yang masih berstatus draft.

### ADR-D-CNP-0008 - Nilai hasil suntingan manual bertahan terhadap hitung ulang, dan selalu terlihat

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Ketika petugas mengganti nilai hasil perhitungan dengan nilai suntingan sendiri, nilai itu **dikunci** dan tidak ditimpa oleh perhitungan berikutnya. Layar menampilkan keduanya berdampingan — nilai hitungan yang digantikan dan nilai suntingan yang berlaku — beserta siapa yang menyunting dan kapan.

Alternatif yang ditolak: mengikuti perilaku sistem lama, yaitu suntingan selalu kalah terhadap hitung ulang.

##### Consequences

Model data memerlukan konsep **nilai terkunci** sejak awal: setiap nilai yang dapat disunting menyimpan pasangan `nilai_hitungan` dan `nilai_disunting`, plus `pelaku` dan `waktu`. Menambahkan ini setelah tabel terbentuk berarti membongkar setiap baris yang sudah ada, karena nilai lama tidak dapat dipisahkan kembali menjadi dua asal yang berbeda.

Kombinasi terburuk yang dihindari keputusan ini adalah suntingan diam-diam yang dapat hilang diam-diam — persis keadaan yang terbaca di sistem lama dan dilaporkan di `FINDING-004-suntingan-manual-tertimpa.md`.

Keputusan ini **tidak menunggu** kepastian hipotesis F1/F2 tentang `.IsEditClaim`. Apa pun yang ditemukan di modul Komite nanti, sistem baru mengunci suntingan.

Konsekuensi migrasi data bergantung pada hipotesis mana yang benar, dan itu diuraikan di FINDING-004 bagian 4, bukan di sini.

### ADR-D-CNP-0009 - Keputusan alur disimpan sebagai field terstruktur; komentar tidak pernah dibaca mesin

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Setiap keputusan yang menggerakkan alur disimpan sebagai tiga field terpisah — **tindakan**, **pelaku**, **waktu** — bukan disimpulkan dari isi kolom komentar. Kolom komentar tetap ada sebagai catatan manusia, dan **tidak pernah** menjadi masukan bagi kondisi mana pun.

Sistem lama melakukan sebaliknya: `SethistoryKlaimTreaty` mengambil keputusan lewat `@contains(.CommentSuggest,"Accepted by <nama>")` pada enam pasang kondisi. Uraiannya di `FINDING-002-percabangan-identitas.md` bagian 6.

##### Consequences

Tiga cacat sekaligus hilang, dan ketiganya nyata di sistem lama: satu huruf besar berbeda membuat kondisi gagal tanpa pesan galat; siapa pun yang boleh mengisi komentar dapat memenuhi kondisi alur; dan orang yang sama ditulis dalam dua ejaan sehingga sebagian cabang tidak pernah terpicu.

**Harga yang dibayar ada di migrasi data, bukan di rancangan.** Riwayat lama menyimpan keputusannya hanya di dalam teks, dan sebagian teks itu sudah ditimpa oleh langkah penerjemahan di sistem lama sendiri. Kolom "siapa menyetujui" pada sistem baru tidak dapat diisi lengkap dari data lama secara otomatis; sebagian akan kosong. Besarnya bagian yang kosong belum diketahui dan diukur lewat REQ-013.

Konsekuensi itu diterima secara sadar di muka. Ia bukan temuan yang muncul saat UAT.

### ADR-D-CNP-0010 - Layer dan Retensi Cedant adalah dua entitas terpisah

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Susunan lapisan XOL dan Retensi Cedant disimpan di **dua tabel berbeda**. Sistem lama menyimpan keduanya di satu PageList `SpreadingRisk`, dengan Retensi Cedant sebagai baris ber-`TreatyName="UR"` yang di-`APPEND` oleh `CountLossAllocation_act`.

##### Consequences

Alasan yang menentukan bukan soal penyaring, melainkan **perilaku baris itu sendiri**: `ClaimPercentage` selalu `0` dan `ClaimSpreaded` selalu `0`. Porsi reasuradur untuk baris itu tidak pernah ada. Sesuatu yang selalu nol di dalam daftar alokasi bukan anggota daftar itu — ia konteks bagi daftar tersebut.

Argumen bahwa "kalau ia lapisan nol, tidak akan ada rule yang perlu menyaringnya" **tidak dipakai** sebagai dasar, karena penyaring bisa ada karena sebab lain.

Setelah dipisah, pertanyaan "Retensi Cedant ikut dihitung atau tidak" menjadi eksplisit pada setiap query — bukan bergantung pada apakah penulisnya ingat menuliskan penyaring. Di sistem lama, sembilan rule menyentuh `SpreadingRisk` antara 16 dan 66 kali tanpa pernah menyebut `"UR"`.

Ini mengunci model data inti dan sangat mahal diubah kemudian: setiap query alokasi, setiap agregasi, dan setiap laporan ikut berubah bentuknya.

##### Tambahan 18 September 2026 — satu kolom yang selalu nol ikut masuk ke tabel ini

Baris Retensi Cedant menerima `TotalClaim = Local.TotalUR` (`Activity\CountLossAllocation_act.xml` baris 7015–7120).

Dan `Local.TotalUR` pada cabang `Local.Currency == .Currency` dihitung sebagai `(.Deductible * 0 * Local.ProrateClaim/100)` — **dikalikan nol**, sehingga selalu nol. Penetapan `Local.UR` pada cabang yang sama, di berkas yang sama, tidak memakai `* 0`.

Artinya: untuk seluruh klaim bermata uang sama, tabel Retensi Cedant yang baru akan memuat kolom `TotalClaim` yang **selalu bernilai nol** — bila perilaku lama dipertahankan.

Apakah perilaku itu dipertahankan atau diperbaiki **belum diputuskan**, dan bukan keputusan teknis: ia dibawa ke akuntansi sebagai `ASK-AKUNTANSI.md` pertanyaan **2b**. Keputusan ADR ini tentang **bentuk tabel** tetap berlaku apa pun jawabannya.

### ADR-D-CNP-0011 - Menjalankan ulang perhitungan alokasi harus menghasilkan keadaan yang identik

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Perhitungan alokasi bersifat **idempoten**: menjalankannya dua kali atas masukan yang sama menghasilkan keadaan yang sama persis. Ini diuji otomatis, bukan sekadar dijanjikan.

##### Consequences

Alasannya tidak punya tandingan: bila menekan tombol hitung dua kali menghasilkan angka berbeda, tidak ada satu pun angka di sistem itu yang bisa dipertanggungjawabkan.

**Konsekuensi terhadap shadow-run, dan ini yang paling mahal.** Bila sistem lama ternyata tidak idempoten, ada data produksi dengan baris ganda. Itu bukan sekadar urusan migrasi data — ia **merusak dasar perbandingan**. Hitungan sistem baru tidak dapat dibandingkan terhadap baris yang sudah terlanjur ganda, karena selisihnya akan terbaca sebagai kesalahan sistem baru padahal berasal dari kerusakan baseline.

Karena itu, sebelum shadow-run dimulai: baseline dibersihkan lebih dulu, **atau** case yang terdampak dikeluarkan dari perbandingan dan didaftar terpisah. Salah satu harus dipilih; tidak boleh dibiarkan tercampur. Lihat juga ADR-D-CNP-0005.

Yang membuat risiko ini nyata di sistem lama: `CountLossAllocation_act` melakukan `<APPEND>` baris `"UR"`, dan **tidak ditemukan di XML** penghapusan PageList `SpreadingRisk` sebelum `APPEND` itu. Satu-satunya penghapusan yang menyasar `SpreadingRisk` di seluruh folder adalah `Property-Remove` atas properti tunggal `.SpreadingRisk(Param.idx).AdjClaimValue` di dua rule — bukan atas daftarnya. Selain itu `CountLossAllocation_act` langkah 10 memuat `Property-Remove` yang **parameter sasarannya kosong** di ekspor ini.

Besarnya kerusakan diukur, bukan ditebak — REQ-014 menghitung berapa case yang memiliki lebih dari satu baris `"UR"`. Keputusan idempoten ini **tidak menunggu** hasil pengukuran itu.

### ADR-D-CNP-0012 - Tambalan baru dideteksi lewat sapuan ulang, bukan lewat register

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Tambalan per-case dan per-identitas dibekukan sejak tanggal kesepakatan. **Selama pekerjaan migrasi ini berjalan, pembekuan berlaku mutlak — tidak ada pengecualian.** Yang menjadi **alat deteksi** bukan register, melainkan **sapuan ulang otomatis** atas setiap ekspor XML baru, dibandingkan terhadap `BLUEPRINT.md` bagian 7.5. Register pengecualian tetap dibuat, tetapi statusnya pelengkap.

##### Consequences

Pembekuan total akan dilanggar diam-diam; yang tercatat lebih berguna daripada yang dilarang. Tetapi register yang diisi manusia punya cacat yang sama dengan `CNPStatusCase` — bisa tidak pernah diisi, dan ketiadaannya tidak menghasilkan sinyal apa pun. **Kebasian analisis tidak boleh bergantung pada kepatuhan orang.**

Sapuan ulang mendeteksi tambalan baru apakah registernya diisi atau tidak. Selisih antara hasil sapuan dan bagian 7.5 adalah daftar tambalan yang masuk sesudah ekspor terakhir.

Konsekuensi kedua: **setiap artefak membawa stempel asal** — tanggal ekspor XML yang menjadi dasarnya dan jumlah berkas yang disapu. Tanpa itu, tidak ada yang dapat menilai dokumen ini berlaku untuk keadaan kapan. Stempel sudah dipasang di seluruh artefak per 18 September 2026 (279 berkas, ekspor 2026-09-08/09).

##### Pengecualian dihapuskan — 19 September 2026 (menutup A13)

Versi pertama memberi pengecualian untuk *"tambalan yang menahan pembayaran atau menghentikan operasi"*, lalu menyisakan pertanyaan terbuka: **siapa berwenang menyetujuinya** (`_selesai/OPEN-QUESTIONS.md` A13). Pertanyaan itu ditutup dengan **menghapus pengecualiannya**, bukan dengan menjawabnya.

Selama pekerjaan ini berjalan, pembekuan berlaku mutlak. Tidak ada pengecualian, jadi tidak ada yang perlu berwenang menyetujuinya, dan tidak ada yang perlu ditanyakan ke luar.

**Register pengecualian tetap ada, dan tetap kosong.** Ia tidak dibongkar: bila kelak ada yang mengisinya, saat itulah wewenang penyetujunya ditetapkan — oleh orang yang mengisinya, pada saat ia mengisinya. Register yang kosong adalah pernyataan yang dapat diperiksa; ketiadaan register bukan.

Yang **tidak** berubah: alat deteksinya tetap sapuan ulang, bukan register. Pembekuan mutlak justru mempertajam alasan itu — semakin tegas larangannya, semakin kecil kemungkinan pelanggarnya mencatat sendiri.

##### Tambahan 18 September 2026 — sapuan mencakup tiga lapisan, bukan satu

Versi pertama menyebut "sapuan ulang atas setiap ekspor XML baru". **Itu terlalu sempit.** Tiga jalur yang menggerakkan sistem tidak berjejak di XML sama sekali (lihat `BLUEPRINT.md` bagian 14), termasuk tanggal tertanam di dalam `PROC_GENERATE_SEQUENCE_NUMBER` yang merupakan kelas tambalan yang sama dengan 29 langkah di bagian 7.

Sapuan berkala mencakup:

| Lapisan | Yang dicari |
|---|---|
| **Ekspor XML** | literal `CLMNP-…`, `IDMaster`, identitas operator/orang, pola `@contains` atas teks bebas |
| **DDL** | tanggal dan periode tertanam, nilai mati di dalam procedure, kolom baru, constraint yang hilang |
| **Daftar database link** | ketergantungan lintas basis data yang baru muncul |

Lapisan ketiga ditambahkan setelah `V_MST_USER_TEKNIS` terbaca. Sebelum itu tidak ada yang tahu perlu menyapunya — yang justru menjadi alasan terkuat mengapa daftar lapisan ini sendiri harus ditinjau setiap kali sumber baru masuk.

### ADR-D-CNP-0013 - Perhitungan dijalankan sebagai turunan dari data, bukan sebagai akibat penekanan tombol

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Setiap nilai hasil perhitungan adalah **turunan** dari masukannya. Begitu masukan berubah, turunannya ikut berubah. Tidak ada tombol "hitung", tidak ada keadaan setengah jadi, dan tidak ada urutan tindakan yang harus diingat pengguna.

##### Consequences

Sistem lama melakukan sebaliknya, dan itu terbukti dari XML: **tidak ada satu pun activity yang memanggil `CountClaimTNP_Act`** — ia dipicu dari Section lewat `<pyActivity>`, 42 rujukan di tiga berkas. Urutan jalannya `CountClaimTNP_Act`, `CountLossAllocation_act`, `GenerateCFS_act`, dan `SaveToOS` tidak ditetapkan di kode mana pun; ia ditentukan kontrol mana yang ditekan.

**Sapuan `MEMORI_PEMAHAMAN.MD` atas urutan kerja, SOP, dan langkah petugas menghasilkan nihil.** Urutan baku itu **tidak terdokumentasi di sumber mana pun** — tidak di XML, tidak di memori. Itu sendiri temuan: urutan yang menentukan hasil perhitungan hanya hidup di kepala petugas, dan akan hilang bersama orangnya.

Keputusan ini menutup jalur itu, bukan memperbaikinya. Konsekuensinya pada rancangan: perhitungan tidak boleh ditempatkan di lapisan tampilan; ia menjadi fungsi murni atas data, dipanggil ulang kapan pun masukannya berubah, dan hasilnya tidak disimpan sebagai keadaan yang dapat menyimpang dari masukannya.

**Keputusan ini mengikat bersama ADR-D-CNP-0011 dan ADR-D-CNP-0012, bukan berdiri sendiri.** Uraian rantai sebab-akibatnya ada di `BLUEPRINT.md` bagian 12: karena hasil bergantung urutan klik, selalu ada klaim yang keluar jalur; karena selalu ada yang keluar jalur, perbaikan termurah adalah menambal klaim itu satu per satu; dan begitulah 29 langkah tambalan menumpuk selama delapan tahun. Mengambil satu dari tiga keputusan tanpa dua lainnya akan mengembalikan pola yang sama dalam bentuk baru.

### ADR-D-CNP-0014 - Aturan penguraian teks ke angka, dan perlakuan atas kurs yang tidak ada

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.
> Berlaku hanya untuk keadaan sistem pada kedua ekspor itu.

Nilai uang di sistem lama tersimpan sebagai teks di dalam `DATA_JSON`, dan bahkan view `CLAIMXOL` mengeluarkannya sebagai `varchar2`. Aturan penguraiannya ditetapkan **tertulis dan di muka**, sebelum uji coba migrasi mana pun dijalankan.

##### Aturan penguraian

| Hal | Perlakuan |
|---|---|
| Pemisah ribuan | ditolak bila ada; angka yang sah tidak memakainya |
| Tanda desimal | titik saja; koma sebagai desimal ditolak dan masuk daftar |
| Spasi di awal/akhir | dipangkas, tidak dianggap galat |
| Notasi ilmiah (`1.2E5`) | ditolak dan masuk daftar |
| Posisi tanda minus | hanya di depan; `123-` ditolak |
| Presisi melebihi ADR-D-CNP-0003 | ditolak dan masuk daftar, **tidak dibulatkan** |

##### Tiga nilai kosong yang tidak boleh disamakan

Ini yang paling mudah terlewat, dan kerusakannya tidak dapat dipulihkan setelah migrasi:

| Bentuk | Arti yang dipegang | Perlakuan |
|---|---|---|
| `""` (string kosong) | **tidak pernah diisi** | dipetakan ke "belum diisi", bukan nol |
| `null` / kunci tidak ada di JSON | **tidak ada** | dipetakan ke null |
| `"0"` | **diisi nol** | dipetakan ke nol |

Menyamakan ketiganya menghancurkan informasi yang tidak dapat dipulihkan. Apakah ketiga makna itu benar-benar dibedakan oleh sistem lama **belum terbukti dari XML** — bila ternyata tidak dapat dibedakan maknanya, itu pertanyaan tersendiri yang dicatat di `_selesai/OPEN-QUESTIONS.md`, bukan diputuskan sepihak di sini.

##### Kurs yang tidak ada

Konversi tanpa kurs **tidak menghasilkan angka**. Nilai IDR dibiarkan kosong dan ditandai "menunggu kurs". Sistem lama memilih `RETURN 1` — lihat `FINDING-006` — dan itu tidak diwarisi.

##### Consequences

Ambang penghentian migrasi **berbasis nilai, bukan cacah baris**. Satu baris bermasalah senilai lima miliar lebih berat daripada lima ratus baris senilai seratus ribu.

- Setiap baris yang tidak dapat diurai diselesaikan satu per satu, berapa pun jumlahnya. Tidak ada baris yang dibuang karena "cuma sedikit".
- Migrasi dihentikan bila **nilai total yang belum terselesaikan** melewati ambang yang disepakati akuntansi, atau bila **ada satu baris** yang nilainya melewati ambang material.

Kedua angka ambang itu belum ditetapkan dan menunggu akuntansi — `ASK-AKUNTANSI.md` pertanyaan 4.

##### Tambahan 18 September 2026 — aturan diturunkan dari kotoran yang benar-benar ada

Daftar aturan di atas semula disusun dari kemungkinan umum. Sapuan seluruh 279 berkas menggantinya dengan yang terbukti.

| | Jumlah |
|---|---|
| Seluruh pemanggilan `toDecimal` | **127** |
| Didahului pembersihan `replaceAll` | **9** — seluruhnya atas `.Deductible2` |
| **Tanpa pembersihan apa pun** | **118** |

**Koma sebagai pemisah desimal bukan hipotesis.** `@toDecimal(@replaceAll(.Deductible2,",","."))` di `Activity\CountLossAllocation_act.xml` membuktikan seseorang pernah menemuinya di produksi dan menambalnya. Tidak ada yang menulis pembersihan untuk masalah yang tidak pernah terjadi.

**Pembersihan lain yang ditemukan, dan sasarannya bukan angka**:

| Ekspresi | Sasaran |
|---|---|
| `@replaceAll(pyWorkPage.ClaimData.InsuredName,",","")` | menghapus koma dari **nama**, bukan dari angka — diduga demi keamanan muatan |
| `@replaceAll(.AcceptedNo,".","")` | menghapus titik dari **nomor dokumen** |
| `@pxReplaceAllViaRegex(.NoAccount,"[^0-9]","")` | menyisakan hanya digit pada **nomor rekening** |

Ketiganya memperlihatkan pola yang sama: pembersihan dilakukan **per tempat, saat masalahnya muncul**, bukan sebagai aturan yang berlaku menyeluruh.

**Yang paling perlu diperhatikan**: dua dari 118 pemanggilan tanpa pembersihan adalah `@toDecimal(TempKasir.CARI20)` dan `@toDecimal(TempKasir.CARI21)` — muatan menuju sistem kasir, yaitu nilai yang masuk jurnal menurut `ASK-AKUNTANSI.md` bagian 0.2.

Aturan penguraian di sistem baru karena itu **berlaku di satu tempat untuk semua nilai**, bukan ditambahkan per ekspresi saat kegagalannya ditemukan.

##### Tambahan — keadaan "menunggu kurs" menular

Nilai yang belum terkonversi menandai **seluruh nilai turunannya**. Apa pun yang dihitung dari nilai bertanda "menunggu kurs" ikut bertanda sama.

Nilai bertanda itu **tidak boleh** masuk penjumlahan, tidak boleh dibandingkan terhadap ambang persetujuan, dan tidak boleh muncul di laporan — baik sebagai nol maupun sebagai nilai apa adanya.

Tanpa aturan penularan ini, keputusan "gagalkan, jangan mengarang angka" hanya berlaku satu tingkat: nilai pertama ditandai, lalu penjumlahan di tingkat berikutnya memperlakukannya sebagai nol dan cacatnya kembali muncul dalam bentuk yang lebih sulit dilihat — persis pola `RETURN 1` yang sedang kita tinggalkan (`FINDING-006`).

Ini melengkapi `ADR-D-CNP-0019` lapis 2: keadaan di tingkat baris mencatat *bahwa* perhitungan gagal; aturan penularan ini menentukan *apa yang terjadi pada baris-baris sesudahnya*.

### ADR-D-CNP-0015 - Akseptasi adalah satu entitas dengan keadaan, bukan dua tabel

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.
> Berlaku hanya untuk keadaan sistem pada kedua ekspor itu.

Akseptasi biasa dan akseptasi bersyarat (Subjectivity) adalah **satu entitas** dengan keadaan yang berbeda, bukan dua entitas. Keadaannya: biasa, bersyarat, dan bersyarat-gugur — beserta daftar syarat dan status pemenuhannya.

Sistem lama memisahkannya jadi dua tabel: `OS_AKSEPTASI_KLAIM` dan `OS_AKSEPTASI_SUBJECTIVITY`.

##### Consequences

Daur hidup yang berbeda adalah **keadaan**, bukan identitas. Dua hal menjadi entitas berbeda bila identitasnya berbeda; bila satu akseptasi bersyarat yang syaratnya terpenuhi berubah menjadi akseptasi biasa **yang sama**, maka sejak awal ia satu benda yang sedang berada di keadaan tertentu.

Bukti yang menguatkan datang dari sistem lama sendiri: view `CLAIMXOL` harus meng-`UNION ALL` kedua tabel untuk mendapat gambaran utuh. Sesuatu yang harus selalu digabungkan kembali tidak seharusnya dipisah.

Syarat yang gugur menjadi **transisi keadaan yang tercatat**, bukan penghapusan atau perpindahan tabel.

**Temuan yang menyertainya dan belum terjawab**: di folder Claim, `.IsSubjectivity` hanya **dibaca** (`==true`), disalin ke `Local.Subjectivity`, dan **tidak ada satu pun rule yang menggugurkan, membatalkan, atau mengakhiri akseptasi bersyarat yang syaratnya tidak pernah terpenuhi.** Delapan belas kemunculan `Subjectivity` di 279 berkas, tidak satu pun berupa transisi keadaan. Jadi dugaan "menggantung selamanya" konsisten dengan isi folder ini — *disimpulkan dari ketiadaan rule, bukan dari adanya rule*. Pemastiannya `DEFERRED-TO-KOMITE-SESSION`.

### ADR-D-CNP-0016 - Tanpa database link ke sistem lain — dan tanpa data pegawai sama sekali

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.
> Berlaku hanya untuk keadaan sistem pada kedua ekspor itu.

Sistem baru **tidak mewarisi ketergantungan lintas basis data**.

**Modul Claim Non Prop tidak mengambil data HRD sama sekali** — tidak lewat database link, tidak lewat API, tidak lewat salinan tersinkron. Pelaku sebuah tindakan disimpan sebagai **potret**: pengenal operator dan nama **sebagaimana tercatat pada saat tindakan terjadi**. Tidak ada pencarian ke HRD, tidak ada penyegaran, tidak ada salinan yang disinkronkan.

Sebabnya bukan sekadar lingkup. **Pelaku suatu tindakan adalah fakta yang terjadi pada satu saat, dan fakta itu tidak berubah ketika orangnya pindah bagian atau berhenti.** Menyimpannya sebagai rujukan hidup ke HRD membuat catatan lama ikut berubah ketika sumbernya berubah — itu merusak audit, bukan memperkayanya.

Sistem lama memakai `hrdasm.v_hrd_mst@asmd.sinarmas.co.id` — database link langsung ke basis data sistem HRD, di dalam view `V_MST_USER_TEKNIS`.

##### Consequences

Tiga alasan, tidak satu pun soal selera:

1. Database link **melewati kontrol akses aplikasi HRD sepenuhnya**.
2. Ia **mengikat kita pada struktur tabel internal** sistem lain. Bila mereka mengubah kolom, sistem kita rusak tanpa ada yang memberi tahu.
3. Bila tautan putus, view mengembalikan **kosong** — daftar pengguna menghilang tanpa pesan galat. Itu pola yang sama persis dengan `RETURN 1` pada kurs (`FINDING-006` bagian 3): kegagalan yang menyamar menjadi hasil yang sah.

Bila HRD tidak dapat menyediakan API maupun mekanisme sinkronisasi, itu percakapan yang dibawa ke luar. Rancangannya **tidak menunggu** jawaban itu, dan tidak dibuat seolah database link adalah salah satu pilihan.

Apakah tautan yang ada sekarang resmi disepakati **tidak lagi ditanyakan**: ia tautan milik sistem lama, dan sistem baru tidak memerlukan data di ujungnya.

##### Penyelesaian K1 — 19 September 2026, sebagai **K1a**

`_selesai/OPEN-QUESTIONS.md` K1 mencatat bahwa ADR ini melarang sesuatu yang sudah berjalan produksi: `V_MST_USER_TEKNIS` memakai `hrdasm.v_hrd_mst@asmd.sinarmas.co.id`, satu-satunya pemakaian `@<host>` di keempat puluh sembilan berkas DDL (D20). Dua cabang ditulis sejajar: **K1a** (ADR tetap, yang lama dibiarkan) dan **K1b** (ADR direvisi dengan pengecualian bernama).

**Dipilih K1a, dan alasannya menghapus biaya yang semula melekat padanya.**

K1a semula mahal karena ia menuntut *"tetapkan cara lain memenuhi kebutuhan data HRD"* — pekerjaan baru yang belum dilingkupi ADR mana pun. Biaya itu **tidak ada**, karena kebutuhan itu tidak ada: dengan aturan potret di atas, modul ini tidak pernah memerlukan data HRD. Larangan ADR ini karena itu **tidak pernah bertabrakan dengan kebutuhan kita**.

| | Keadaan |
|---|---|
| `V_MST_USER_TEKNIS` | **tidak dimiliki, tidak dimigrasi.** View HRD, bukan entitas klaim — ADR-D-CNP-0026 (batas kepemilikan mengikuti nama class) |
| View lama beserta link-nya | dibiarkan hidup di sistem lama; bukan milik kita, dan tidak ada yang perlu kita gantikan |
| ADR-D-CNP-0023 dan ADR-D-CNP-0028 | **tidak ikut ditinjau.** Peninjauan itu konsekuensi K1b, dan K1b tidak dipilih |
| D20 | tetap tercatat sebagai fakta; kedudukannya kini **catatan atas sistem lama**, bukan pekerjaan |

**Rujukan ADR-D-CNP-0024 dicabut.** Versi pertama bagian atas ADR ini menyebut *"ADR-D-CNP-0024 tentang pemisahan pengguna, jabatan, dan keanggotaan komite"*. ADR-D-CNP-0024 adalah **kunci alami akseptasi**; rujukan itu salah sejak ditulis, dan kalimat yang memuatnya sudah hilang bersama janji HRD.

### ADR-D-CNP-0017 - Data aplikasi hanya diubah lewat aplikasi

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.
> Berlaku hanya untuk keadaan sistem pada kedua ekspor itu.

Tidak ada jalur tulis ke data aplikasi selain lewat aplikasi itu sendiri. Proses batch yang perlu mengubah data memanggil layanan yang sama dengan yang dipakai layar, sehingga aturan dan jejaknya berlaku sama.

Di sistem lama, schema `POOLDATA` memiliki akses ke tabel work Pega `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` — artinya data case dapat diubah oleh prosedur basis data, tanpa melewati aplikasi dan tanpa jejak di dalamnya.

##### Consequences

Alasannya berdasar rekam jejak analisis ini sendiri. Sudah ditemukan **tiga** mekanisme yang hasilnya tidak dapat dilacak ke pelakunya:

| Mekanisme | Dokumen |
|---|---|
| Penggantian identitas `DARTO` menjadi `CHRISTINEANGELINA` | `FINDING-002` bagian 9 |
| Keputusan alur diambil dari teks komentar bebas | `FINDING-002` bagian 6 |
| Suntingan manual tertimpa hitung ulang tanpa peringatan | `FINDING-004` |

Pintu tulis kedua akan menambah yang keempat.

**Keputusan desainnya final; ukurannya belum, dan itu risiko lingkup.** Bila REQ-021 menunjukkan ada proses produksi yang menulis lewat jalur itu, proses-proses tersebut harus ditulis ulang **sebelum** cutover — pekerjaan yang belum masuk perkiraan mana pun.

Karena itu **REQ-021 ditandai BLOCKER perkiraan biaya, bukan sekadar BLOCKER teknis.** Bila hasilnya besar, ia dibawa ke manajemen sebagai pokok tersendiri, bukan diselipkan ke dalam laporan teknis.

##### Batas yang baru terlihat — 19 September 2026, saat tiket 01 dikerjakan

ADR ini menyatakan tidak ada jalur tulis ke data aplikasi selain lewat aplikasi. Penegakannya dirancang lewat `GRANT` dan `REVOKE` atas objek, dan itulah isi `ddl-usulan/V00_HAK_AKSES.sql`.

**Penegakan itu berlaku pada tingkat objek saja, dan itu belum pernah dikatakan.**

Oracle mengenal hak sistem berakhiran `ANY` — `SELECT ANY TABLE`, `INSERT ANY TABLE`, `UPDATE ANY TABLE`, `DELETE ANY TABLE`, `ALTER ANY TABLE`. **Hak `ANY` mengatasi hak objek.** Satu akun yang memegangnya menulis ke tabel mana pun di skema mana pun, dan seluruh `REVOKE` di `V00_HAK_AKSES.sql` tidak menyentuhnya sedikit pun.

Yang membuat ini bukan kekhawatiran teoretis: premis ADR ini sendiri adalah bahwa `POOLDATA` sudah memegang hak tulis yang tidak ia perlukan. Pertanyaan apakah ia juga memegang hak `ANY` **belum pernah diajukan** — tidak di ADR ini, tidak di ADR-D-CNP-0028, tidak di `V00_HAK_AKSES.sql`, dan tidak di satu pun REQ sebelum hari ini.

Diajukan sebagai **REQ-037**, ditandai BLOCKER. Pencabutannya bukan pekerjaan pemilik skema; ia pekerjaan DBA di tingkat instance.

**Sampai REQ-037 kembali, klaim ADR ini dibaca dengan batas ini terpasang**: satu pintu tulis ditegakkan pada tingkat objek, dan belum diketahui apakah ada pintu di tingkat sistem yang melewatinya. Keputusannya tidak berubah; yang berubah adalah apa yang boleh diklaim sudah tertutup.

##### Apa yang benar-benar dapat dijanjikan skema ini — rumusan diperbaiki 19 September 2026

**Keputusannya tidak berubah.** Yang berubah rumusan klaimnya, dan itu diubah di sini supaya tidak ada yang mengutip janji yang tidak dapat ditepati.

###### Yang tidak dapat dijanjikan

*"Tidak ada jalur tulis selain lewat aplikasi"* — tanpa syarat — **tidak dapat dipertahankan** selama REQ-037 terbuka. Enam jalur melewati hak objek, dan tidak satu pun dapat ditutup dari dalam skema:

| Jalur | Kenapa skema tidak dapat menutupnya |
|---|---|
| Hak sistem berakhiran `ANY` | dicabut di tingkat instance, oleh DBA |
| Prosedur definer's rights | pemiliknya di skema lain; **pola ini sudah dipakai** — `PEGA_JSON_OS_AKSEP_KLAIMTNP` |
| Hibah ke `PUBLIC` | tidak terlihat saat memeriksa akun satu per satu |
| `CREATE ANY TRIGGER` | menaruh penulis di dalam tabel kita sendiri |
| Peran yang memuat hak `ANY` | `DBA`, `IMP_FULL_DATABASE`, dan peran buatan lokal |
| `GRANT ANY PRIVILEGE` | **memulihkan jalur tulis besok**, sesudah pemeriksaan hari ini bersih |

Jalur terakhir mengubah sifat persoalan: REQ-037 mengukur **satu titik waktu**, sementara ADR ini bicara tentang **keadaan yang bertahan**. Jawaban bersih hari ini tidak menjamin keadaan besok.

###### Yang dapat dijanjikan, dan ditegakkan

> **Skema ini tidak dapat mencegah tulisan dari luar pintu. Tetapi tidak ada tulisan dari luar pintu yang tidak meninggalkan jejak.**

Ditegakkan `ddl-usulan/Z01_PENGAWASAN_TULIS.sql`: setiap `INSERT`, `UPDATE`, dan `DELETE` atas tabel skema ini yang datang dari akun selain `KLAIMNP_APP` tercatat beserta siapa, kapan, dan tabel mana.

Syaratnya memakai `SESSION_USER`, bukan `CURRENT_USER`. Pembedaan itu yang membuatnya bekerja: di dalam prosedur definer's rights, `CURRENT_USER` menjadi pemilik prosedur dan jejaknya hilang — yaitu persis jalur kedua di tabel atas. `SESSION_USER` tidak berubah oleh prosedur siapa pun.

Kebijakan itu **dijalankan DBA, bukan pemilik skema**, dan `KLAIMNP` sengaja tidak diberi `AUDIT_ADMIN`: pengawasan yang dapat dimatikan oleh yang diawasi bukan pengawasan.

###### Apa yang tetap terbuka

Pengawasan memperlihatkan, tidak mencegah. Menutup jalurnya tetap pekerjaan DBA di tingkat instance, dan itu tetap **REQ-037, BLOCKER**. Yang berubah: klaim ADR ini kini dapat dipertahankan apa adanya sambil REQ-037 berjalan, alih-alih menunggu REQ-037 untuk berarti.

### ADR-D-CNP-0018 - Klaim yang terdampak penjaga tanggal yang salah dimigrasi apa adanya, dan didaftar

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.
> Berlaku hanya untuk keadaan sistem pada kedua ekspor itu.

Klaim yang hasilnya berbeda bila penjaga `CheckDateDOL_Act` dijalankan dengan perbandingan yang benar **dimigrasi apa adanya**. Tidak ada koreksi otomatis, tidak ada perubahan status.

Sebagai gantinya: penjaga yang benar dijalankan atas seluruh data lama, hasilnya menjadi daftar, dan keputusan per klaim diambil manusia.

##### Consequences

Klaim yang sudah dibayar tidak boleh berubah statusnya karena sistem baru menghitung ulang sebuah penjaga — itu bukan koreksi, itu mengubah angka yang sudah diterima akuntansi. Membiarkannya tanpa daftar berarti masalahnya ikut pindah tanpa ada yang tahu.

Tiga syarat yang mengikat:

1. **Daftar dihasilkan sebelum cutover.** Ia prasyarat, bukan laporan pasca-migrasi.
2. **Dua daftar terpisah, tidak digabung.** Yang seharusnya tertolak tetapi lolos; dan yang seharusnya lolos tetapi tertolak. Kelompok kedua lebih sensitif — itu klaim yang mungkin ditolak secara keliru.
3. **Keputusan per klaim dicatat.** Tidak ada koreksi massal.

**Keterjangkauan data — diperiksa dan hasilnya baik.** Sempat dikhawatirkan daftar ini tidak dapat dibuat karena `.EndDateTreaty` tersimpan IN-BLOB. Pemeriksaan menunjukkan sisi lainnya terjangkau: `.ClaimData.EndDateTreaty` diisi dari `pyWorkPage.TreatyInMaster.Termination`, dan **`POOLDATA.TREATYINDETAIL.TERMINATION` adalah kolom bertipe `DATE`**. `POOLDATA.TREATYCONTRACT.TREATYENDDATE` juga `DATE`. Sisi kiri perbandingan, `DATEOFLOSS`, adalah kolom `VARCHAR2(8)` di tabel work.

Jadi **REQ-020 dapat dijalankan dengan SQL biasa** — tanpa membongkar blob — dengan `TO_DATE(DATEOFLOSS,'YYYYMMDD')` dibandingkan terhadap kolom `DATE` itu. Yang masih perlu dipastikan hanyalah kunci penghubungnya; `MASTERID` termasuk kolom yang di-expose di tabel work, tetapi pasangannya di `TREATYINDETAIL` belum diverifikasi.

### ADR-D-CNP-0019 - NULL bukan nol, dan keadaan perhitungan disimpan di tingkat baris

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Tiga lapis, disusun dari yang termurah. Lapis berikutnya hanya ditambahkan di tempat lapis sebelumnya tidak cukup.

##### Lapis 1 — `NULL` berarti belum dihitung, nol berarti dihitung dan hasilnya nol

Keduanya **tidak boleh pernah disamakan**, di mana pun, oleh kode mana pun. Tidak ada `NVL(x,0)` yang diam-diam menghapus perbedaan itu, tidak ada kolom nilai ber-`DEFAULT 0`, dan tidak ada layar yang menampilkan kosong sebagai `0`.

Biayanya nol — hanya disiplin — dan ia menangkap sebagian besar manfaat dari seluruh keputusan ini.

##### Lapis 2 — keadaan di tingkat baris, bukan tingkat kolom

Setiap baris hasil perhitungan membawa tiga kolom: `dihitung_pada`, `status_perhitungan`, dan `versi_aturan` yang dipakai. Satu perhitungan gagal maka statusnya terbaca di baris itu, dan kolom yang gagal bernilai `NULL`.

Tiga kolom, bukan ratusan.

##### Lapis 3 — penanda per-nilai, hanya di satu tempat

Hanya untuk **nilai uang hasil konversi mata uang**, karena di sana satu mata uang dapat gagal sementara mata uang lain berhasil **di dalam baris yang sama** — satu-satunya kasus yang tidak tertangkap lapis 2.

##### Consequences

Usulan pertama saya — satu kolom penanda untuk setiap nilai turunan — ditolak, dan alasan penolakannya benar: ratusan kolom tambahan akan diisi asal saat implementasi, lalu tidak ada yang percaya isinya. **Penanda yang tidak dipercaya lebih buruk daripada tidak ada penanda.**

Urutannya juga dibalik dari usulan saya. Bukan *"kalau terlalu mahal, batasi ke uang"*, melainkan *"mulai dari yang murah dan menyeluruh, lalu perdalam hanya di tempat yang paling mahal kalau salah"*.

**Masalah yang diselesaikan** — tiga mekanisme di sistem lama, semuanya berakar pada tidak adanya cara membedakan nol dari gagal:

| Mekanisme | Kegagalan | Terlihat sebagai |
|---|---|---|
| `GETCURRENCYSTANDARD` | kurs tidak ada | kurs bernilai 1 (`FINDING-006`) |
| Database link HRD | tautan putus | daftar pengguna kosong (`ADR-D-CNP-0016`) |
| Precondition tidak cocok | kondisi gagal | langkah dilewati tanpa pesan (`FINDING-002` bagian 2) |

Ketiganya menghasilkan jawaban yang bentuknya benar dan isinya salah.

### ADR-D-CNP-0020 - Fakultatif dan Treaty dua entitas berbeda; sistem ini hanya memiliki Treaty

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Polis fakultatif dan polis treaty adalah **dua entitas berbeda**. Sistem ini **hanya memiliki Treaty**; Fakultatif dibaca sebagai rujukan, tidak dimodelkan, tidak dimiliki. Bila kelak ada modul yang memerlukannya, modul itu yang memilikinya.

##### Consequences

**Alasan pemisahannya adalah domain, bukan bentuk penyimpanan.** Fakultatif dinegosiasikan per risiko, satu per satu; treaty adalah kontrak payung atas satu portofolio. Beda cara lahirnya, beda dokumennya, beda kewajiban para pihaknya.

Argumen bentuk data — bahwa `V_POLIS` bercabang pada `QuotationData.BusinessFac` dan membaca jalur JSON yang berbeda (`$.PolicyData.StartDateTime` versus `$.StartDate`) — adalah **akibat, bukan sebab**. Ia pendukung yang baik dan bukan bukti utama. Menyimpulkan identitas dari bentuk penyimpanan akan bertentangan dengan `ADR-D-CNP-0015`, yang justru menyatukan dua tabel karena identitasnya sama.

**Yang membatasi kepemilikan**: lingkup modul ini adalah klaim non-proporsional, dan itu berjalan di jalur Treaty. Memodelkan Fakultatif berarti membangun entitas yang tidak pernah ditulis dan tidak pernah dihitung — biaya tanpa penerima manfaat.

**Pemeriksaan yang diminta, dan hasilnya tidak sebagaimana diharapkan.** Pertanyaannya: adakah klaim non-prop yang pernah berjalan di cabang `'F'`? Sapuan 279 berkas menemukan **`BusinessFac` nol kemunculan** — modul Claim **tidak pernah membacanya sama sekali**.

Artinya modul Claim tidak dapat membedakan kedua cabang; ia menerima apa pun yang dikeluarkan `V_POLIS`. Jadi pertanyaan itu **tidak dapat dijawab dari XML** dan menjadi pertanyaan data: apakah ada baris `json_polis` ber-`BusinessFac='F'` yang pernah dirujuk sebuah klaim non-prop. → **REQ-030**.

Konsekuensi yang perlu dicatat sekarang: bila ternyata ada, klaim itu menerima `STARTDATETIME`/`ENDDATETIME` berpresisi penuh, sementara jalur Treaty menerima bentuk yang terpotong di jam (`FINDING-005` bagian 6.1) — dua bentuk teks berbeda panjang masuk ke perbandingan yang sama.

Batas kepemilikan ini ditulis di `CONTEXT.md`.

### ADR-D-CNP-0021 - Pengenal polis dan klaim diberi nama menurut isinya; CASEID dipensiunkan

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Istilah `CASEID` **dipensiunkan** dari model internal dan masuk daftar `_Avoid_` di `CONTEXT.md`. Setiap pengenal diberi nama yang menyebut isinya dan pemiliknya.

Sebabnya: nama yang sama membawa arti berbeda di dua tempat. Di `V_POLIS`, `CASEID` = `DATA_JSON.IDNewBisnis`, identitas **polis**. Di `OS_AKSEPTASI_KLAIM`, `CASEID` berbentuk `'ASM-FW-GCNMFW-WORK CLMNP-…'`, identitas **case Pega**.

##### Peta pengenal — lebih dari dua

Penamaan tidak ditetapkan sebelum peta ini lengkap, karena berhenti di dua akan mengulang kesalahan yang sama dalam bentuk lebih halus.

| Pengenal | Ditemukan di | Dugaan pemilik | Status |
|---|---|---|---|
| `CASEID` (V_POLIS) = `DATA_JSON.IDNewBisnis` | view `V_POLIS` | kita | `IDNewBisnis` **nol kemunculan** di 279 XML — hanya ada di view |
| `CASEID` (OS_AKSEPTASI_KLAIM) | tabel akseptasi | Pega | berbentuk `<class> <pyID>` |
| `NOPOLIS` | `JSON_KLAIM`, `JSON_POLIS`, `OS_AKSEPTASI_KLAIM`, `TREATYINPRODUCTION`, 2 prosedur | kita | kolom, 4 tabel |
| `.ClaimData.PolicyNo` / `.ClaimData.PolicyData.PolicyNo` | XML, 115 kemunculan | kita | apakah keduanya sama **belum terbukti** |
| `QuotationData.PolicyNoSinarmas` | XML | **perusahaan lain dalam grup** | bentuk keempat, menyebut nama perusahaan |
| `POLICY_CEDING` | tabel work Pega | **cedant** | nomor polis milik cedant |
| `DLANO_CEDING` / `DLANO_SOB` | `OS_AKSEPTASI_KLAIM` | cedant / SOB | nomor DLA, bukan nomor polis |

**Sedikitnya lima bentuk pengenal polis hidup berdampingan, dan sedikitnya dua di antaranya milik pihak lain.** Nomor polis cedant dan nomor polis kita adalah dua hal berbeda; `POLICY_CEDING` memperlihatkan keduanya memang disimpan berdampingan.

Pemetaan mana yang sama dan mana yang berbeda **belum selesai** dan menunggu DDL tabel work serta REQ-023.

##### Consequences

**Penggantian nama berhenti di batas integrasi.** Payload REST ke Arasapas dan ke sistem kasir membawa `CASEID` apa adanya. Nama internal boleh berubah; **kontrak ke luar tidak**. Pemetaan nama internal ke nama payload ditulis eksplisit di Boundary Contract, supaya penggantian istilah tidak merambat ke payload dan memutus integrasi yang selama ini berjalan.

Satu nama untuk dua hal tidak menimbulkan galat, hanya salah paham — dan salah pahamnya baru muncul saat dua orang membicarakan hal berbeda dengan kata yang sama.

### ADR-D-CNP-0022 - Masa berlaku treaty disimpan sebagai tanggal, dan batasnya inklusif

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Masa berlaku treaty disimpan sebagai **tanggal**, bukan tanggal-waktu. Kerugian yang terjadi **tepat pada** tanggal akhir treaty **tertutup** — batasnya inklusif.

##### Alasan menyimpan sebagai tanggal — bukan yang saya tulis pertama kali

Alasan versi pertama saya: *"untuk masa berlaku treaty, jam tidak membawa arti"*. **Itu salah untuk reasuransi.** Slip treaty lazim menyebut waktu attachment secara eksplisit — misalnya *"at 00:01 hours Local Standard Time, 1st January"* — dan jam bisa sangat berarti untuk kerugian katastrofa yang terjadi di pergantian periode.

Alasan yang benar lebih sederhana dan lebih kuat:

> **Data sumbernya tidak sanggup menyimpan waktu yang dapat dipercaya.** Jalur Treaty di `V_POLIS` berhenti di jam — tanpa menit dan detik — dan nilainya sendiri hasil `SUBSTR` atas teks.

Perbedaan ini penting: bila kelak ada sumber yang membawa waktu attachment sungguhan, alasan versi pertama akan dipakai untuk menolaknya. Alasan yang benar justru mengundangnya.

##### Batas inklusif — dan apa yang sebenarnya dilakukan sistem lama

Seluruh perbandingan di sistem lama memakai operator `>` tegas, tidak pernah `>=`:

```
pyWorkPage.ClaimData.DateOfLoss > pyWorkPage.ClaimData.EndDateTreaty
pyWorkPage.ClaimData.PolicyData.EndDateTime > pyWorkPage.ClaimData.EndDateTreaty
@FormatDateTime(...StartDateTime...) > pyWorkPage.ClaimData.EndDateTreaty
```

Karena penolakan terjadi ketika tanggal kerugian **lebih besar** dari tanggal akhir, kerugian yang jatuh **tepat pada** tanggal akhir tidak ditolak. **Maksud kode lama juga inklusif**, dan itu sejalan dengan praktik reasuransi.

Yang tidak dapat dipastikan adalah apakah maksud itu benar-benar terlaksana, karena perbandingannya dilakukan atas dua format teks yang berbeda (`FINDING-005`). Jadi keputusan ini menegaskan maksud yang sudah ada, bukan mengubahnya.

##### Consequences

**Ditetapkan sebelum REQ-020 dijalankan, dan itu bukan formalitas.** Klaim yang tanggal kerugiannya jatuh tepat di batas akan **berpindah sisi** tergantung pilihan ini — masuk atau keluar dari kedua daftar di `ADR-D-CNP-0018`. Bila ditetapkan setelah query berjalan, daftarnya harus dibuat ulang.

Perbandingan di REQ-020 karena itu memakai `TO_DATE(DATEOFLOSS,'YYYYMMDD') > TERMINATION`, bukan `>=`.

### ADR-D-CNP-0023 - Data akseptasi disimpan dalam satu bentuk kanonik; sistem hilir diberi view, bukan salinan

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Tidak ada data yang disimpan dua kali. Bila sistem hilir memerlukan bentuk kolom, ia diberi **view atau projeksi** atas bentuk kanonik — bukan salinan yang ditulis terpisah.

Sistem lama menyimpannya dua kali. `POOLDATA.OS_AKSEPTASI_KLAIM` punya **66 kolom**, dan `PEGA_JSON_OS_AKSEP_KLAIMTNP` hanya mengisi **8** di antaranya (`CASEID`, `NOCLAIM`, `MASTERID`, `DATA_JSON`, `TANGGAL`, `NOPOLIS`, `STS_REJECT`, `STS_KONVERSI`). **59 kolom sisanya** — `GROSSVALUE`, `ADJUSTERFEE`, `KURSIDR`, `CNPREINSTATEMENT`, `TOTALXOL`, `SALVAGE`, dan seterusnya — diisi oleh sesuatu yang lain.

`DATA_JSON` adalah sumber kebenaran; kolom skalar adalah keluaran proses hilir.

##### Consequences

**Temuan bypass Arasapas menghasilkan ramalan yang dapat diuji, dan itu memperkuat keputusan ini, bukan mengubahnya.** Untuk case yang dibuat `VINCENTVERNANDO_1`, konversi tidak berjalan (`FINDING-002` bagian 8.6). Maka barisnya seharusnya memiliki `DATA_JSON` terisi tetapi **kolom skalarnya kosong**.

Bila ternyata terisi, berarti ada penulis lain yang belum diketahui — dan seluruh pemahaman tentang siapa menulis apa ke tabel ini harus ditinjau ulang. Digabungkan ke **REQ-016**.

Siapa penulis dan siapa pembaca 59 kolom itu tetap **REQ-023**. Keputusan di atas **tidak menunggunya**: penyimpanan ganda ditolak apa pun jawabannya, karena dua salinan yang dapat menyimpang satu sama lain adalah cacat terlepas dari siapa yang menulisnya.

### ADR-D-CNP-0024 - Satu akseptasi per klaim per layer per mata uang

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Kunci alami baris akseptasi adalah **(klaim, layer, mata uang)**, ditegakkan sebagai `UNIQUE`, dengan primary key surrogate.

Sumbernya bukan tebakan: blok yang **dikomentari** di dalam `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` memakai persis kombinasi itu —

```sql
SELECT count(1) INTO id_count FROM OS_AKSEPTASI_KLAIM a
 WHERE CASEID = PegaID AND a.data_json.TypeLoss = LayerT AND a.data_json.Currency = Currency;
```

Orang yang menulis blok itu sedang menjawab pertanyaan yang sama dan sampai pada kesimpulan yang sama. Kodenya dimatikan; alasannya tidak.

##### Penguat dari bentuk penyimpanan — `CLAIMXOL`

Dipindahkan ke sini dari `BLUEPRINT.md` §13.2 pada 18 September 2026, karena tempatnya di sini.

View `POOLDATA.CLAIMXOL` membongkar JSON akseptasi dengan `JSON_TABLE`:

```sql
json_table (data_json, '$.CNPLayerList[*]'
  columns( XOL varchar2 path '$.XOL',
    nested path '$.CNPCurrencyList[*]' columns(
      TotalXOLGross varchar2 path '$.TotalXOLGross', ... )))
```

**`$.CNPLayerList[*]` membungkus `$.CNPCurrencyList[*]`** — layer di tingkat luar, mata uang di tingkat dalam. Jadi bentuk penyimpanan lama sendiri sudah bersusun dua tingkat persis menurut dua dimensi kunci ini, dan urutannya menetapkan arah sarangnya: satu layer memuat banyak mata uang, bukan sebaliknya.

Ini penguat, bukan dasar. Dasarnya tetap blok yang dikomentari di atas, karena di sanalah kombinasinya dinyatakan sebagai **kunci**, bukan sekadar sebagai susunan.

##### Consequences

Sistem lama tidak menegakkannya sama sekali: `OS_AKSEPTASI_KLAIM` **tidak punya primary key maupun unique constraint**, dan prosedurnya **selalu `INSERT`, tidak pernah `UPDATE`** karena logika upsert-nya dikomentari seluruhnya. Baris ganda bukan kemungkinan teoretis.

**Karena itu pelanggaran kunci akan ditemukan, dan itu bukan bukti kuncinya salah.** Saat REQ-018 kembali, hasilnya dipisah tiga:

| Kelompok | Perlakuan |
|---|---|
| Baris **identik di semua nilai** | duplikat murni — digabung |
| Baris **berbeda, dengan pola waktu** | ambil yang terakhir, **catat yang dibuang** |
| Baris **berbeda, tanpa pola waktu** | bawa contohnya untuk ditinjau — berarti ada dimensi pembeda yang belum tertangkap |

Kelompok ketiga adalah satu-satunya yang dapat menggugurkan kunci ini. Dua kelompok pertama adalah akibat bug selalu-`INSERT`, bukan bukti tentang model data.

### ADR-D-CNP-0025 - Tanggal tutup buku disimpan bertanggal berlaku, bukan satu baris tanpa riwayat

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Aturan penggeseran periode dipertahankan: nomor yang diterbitkan setelah hari tutup buku masuk periode berikutnya. Yang dirancang ulang adalah **penyimpanannya**.

Tanggal tutup buku disimpan dalam tabel **bertanggal berlaku** — setiap perubahan menjadi baris baru dengan periode berlakunya sendiri. Ditambah kolom **lingkup**, dengan nilai bawaan global.

##### Consequences

Sistem lama menyimpannya sebagai **satu baris tanpa riwayat**: `SELECT TO_NUMBER(tanggal) INTO v_day_closing FROM POOLDATA.TANGGAL_CLOSING WHERE ROWNUM = 1`.

**Dan itulah sebabnya tambalan lahir.** `PROC_GENERATE_SEQUENCE_NUMBER` memuat:

```sql
IF TRUNC(v_now) <= TO_DATE('02/01/2026','DD/MM/YYYY') THEN
    v_mm_yyyy := '12.2025';  v_tahun := 2025;
```

Perlakuan khusus untuk satu periode tidak punya tempat di data, jadi ia menjadi cabang di dalam kode. Dengan tanggal berlaku, kasus seperti itu cukup menjadi **satu baris data** — dan kelas tambalan ini kehilangan alasan untuk lahir lagi.

Kolom lingkup disediakan sekarang meski belum ada yang memerlukannya; biayanya hampir nol di tahap ini dan mahal ditambahkan setelah tabel terisi.

**Satu masalah yang tidak diselesaikan oleh keputusan ini** dan harus ditangani terpisah: aturan tanggal 25 ada di **dua sumber kebenaran**. `Activity\HitServiceToKasir_Act.xml` menuliskan `@if(TempKasir.CARI19 > 25, ...)` langsung di kode Pega, sementara prosedur membacanya dari tabel. Mengubah tabel tidak mengubah kode. Dibawa ke akuntansi sebagai `ASK-AKUNTANSI.md` pertanyaan 6.3.

### ADR-D-CNP-0026 - Batas kepemilikan mengikuti nama class: -Work- dan -Data- dimiliki, -Int- tidak

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); `pengetahuan/DDL_Script_ClaimNonProp.xls` versi 2026-09-18 10:36 (48 objek); dan `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` (ditempel pengguna 2026-09-18).

Modul klaim **memiliki** klaim, akseptasi, alokasi, dan adjustment. Selebihnya **dibaca, tidak dimiliki**.

Aturannya tidak perlu ditimbang per tabel — `MEMORI_PEMAHAMAN.MD` §2.3 sudah memberinya: **integration class (`-Int-`) adalah pemetaan langsung ke tabel atau view Oracle.**

| Pola class | Kepemilikan |
|---|---|
| `…-Work-…` | **dimiliki** |
| `…-Data-…` | **dimiliki** |
| `…-Int-…` | **dibaca, tidak dimiliki** |

Dengan itu `V_POLIS`, `T_STORAGE_IMAGE`, `EMAILKOMITE`, dan `M_LINK_SERVICE` seluruhnya di luar kepemilikan — tanpa perlu memutuskan satu per satu.

##### Consequences

**Tabel yang bukan milik klaim tidak dimigrasikan, meskipun ada di daftar tarikan.** Daftar itu untuk **memahami**, bukan untuk memindahkan. Ini memotong lingkup migrasi data secara langsung.

Bila modul lain belum siap pada saat cutover, klaim **membaca dari basis data lama lewat antarmuka yang disepakati** — bukan menyalin datanya. Konsekuensinya cutover dapat dilakukan per modul, tidak harus sekaligus.

Yang tetap perlu dipahami meski tidak dimiliki: 298 kolom di DDL tanpa pasangan properti Pega (`BLUEPRINT.md` §19.2) sebagian besar milik modul lain. Memahaminya perlu; memindahkannya tidak.

### ADR-D-CNP-0027 - Dokumen klaim dirujuk, tidak disimpan ulang

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); `pengetahuan/DDL_Script_ClaimNonProp.xls` versi 2026-09-18 10:36 (48 objek); dan `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` (ditempel pengguna 2026-09-18).

Sistem baru **merujuk** dokumen klaim, tidak menyimpannya. Tidak ada migrasi berkas fisik.

##### Consequences

`MEMORI_PEMAHAMAN.MD` §5.7: berkas sesungguhnya berada di **Google Cloud Storage**. `POOLDATA.T_STORAGE_IMAGE` hanya menyimpan **metadata dan URL beserta `EXPDATE`**, yang disegarkan lewat `GetUrlGoogleStorage_Act`.

Jadi kekhawatiran tentang volume berkas yang harus dipindahkan **gugur** — tidak ada berkas fisik di dalam basis data.

Yang tetap perlu dirancang: URL berumur terbatas (`EXPDATE`), sehingga sistem baru memerlukan mekanisme penyegaran yang setara dan tidak boleh menyimpan URL sebagai nilai tetap.

`POOLDATA.GET_TOKEN_STORAGE` dan class `Link-Attachment` termasuk lapisan rujukan ini, bukan lapisan penyimpanan.

### ADR-D-CNP-0028 - Basis data tujuan Oracle, skema baru bersebelahan dengan POOLDATA

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**; dan `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` (ditempel pengguna 2026-09-18).
> Keputusan diambil 18 September 2026. Ini keputusan arsitektur pertama sistem baru; seluruh ADR di bawahnya mewarisinya.

> **Status `accepted` sejak 19 September 2026.**
>
> Versi pertama dokumen ini berstatus `proposed` dengan alasan yang ditulis terang: isinya disusun dari pilihan A yang tertulis lengkap sementara pilihan B dibiarkan kosong — **pembacaan atas formulir, bukan kalimat orang**. Syaratnya konfirmasi eksplisit beserta tanggal.
>
> **Konfirmasi itu datang sebagai kalimat manusia pada 19 September 2026**, dengan alasan yang dinyatakan penuturnya — dua, bukan satu:
>
> 1. **Pindah mesin membuat ADR-D-CNP-0016 kehilangan arti dan ADR-D-CNP-0023 mustahil.** Foreign data wrapper pada dasarnya database link dengan nama lain, dan view lintas mesin bukan view melainkan salinan yang menyamar. Keduanya harus ditulis ulang lebih dulu, dan tidak ada DDL yang boleh berdiri di atas ADR yang sudah tidak berlaku.
> 2. **Mengganti mesin berbarengan dengan migrasi menambah sumber perbedaan yang tidak dapat dipisahkan dari perbedaan yang sedang diukur ADR-D-CNP-0005.** Shadow-run hanya bermakna bila selisihnya punya satu sebab; dua perubahan serentak membuat setiap selisih ambigu.
>
> Alasan kedua tidak ada di versi pertama dokumen ini. Ia datang dari penutur, dan dicatat apa adanya.

Basis data tujuan adalah **Oracle**. Skema baru berdiri sebagai **satu skema tersendiri di instance yang sama dengan `POOLDATA`**.

Versi ditulis sebagai **sekurangnya 12.1, angka pastinya terbuka**, dan itu **tidak ikut naik** bersama status ini. ~~**Gerbang DDL terangkat; gerbang penamaan tidak.** Sampai REQ-032 kembali, tidak ada nama constraint, nama index, maupun singkatan yang boleh ditulis — batas panjang pengenal belum diketahui.~~ **Gerbang penamaan ikut dicabut 19 September 2026**, dan bukan karena REQ-032 kembali: batas **30 byte ditetapkan sebagai aturan tetap**, sah di setiap versi Oracle, sehingga jawaban REQ-032 tidak dapat mengubah satu nama pun. REQ-032 turun jadi verifikasi. Aturan penamaan final: `SPEC-MODEL-DATA.md` bagian 16. Lantai itu berdasar: `BLUEPRINT.md` §11 menyimpulkan basis data **sumber** berjalan di 12c ke atas dari tiga hal sekaligus — `CHECK (x IS JSON)`, notasi titik `a.data_json.Field`, dan klausa `CREATE OR REPLACE EDITIONABLE` yang baru ada sejak 12.1. **Itu lantai sumber, bukan versi tujuan**; keduanya tidak boleh dicampur, dan angka pasti tujuan masih ditunggu.

##### Considered Options

Basis data lain tidak dipilih. Alasannya bukan selera melainkan biaya pada ADR yang sudah berdiri: pindah mesin membuat **ADR-D-CNP-0016** kehilangan arti — foreign data wrapper pada dasarnya database link dengan nama lain — dan membuat **ADR-D-CNP-0023** mustahil, karena view lintas mesin bukan view melainkan salinan yang menyamar. Keduanya harus ditulis ulang lebih dulu, dan tidak ada DDL yang boleh ditulis di atas ADR yang sudah tidak berlaku.

Instance terpisah juga tidak dipilih, atas alasan yang sama dalam bentuk lebih kecil: ADR-D-CNP-0023 mewajibkan sistem hilir diberi **view**, bukan salinan. Antar instance, view menuntut database link, dan itu dilarang ADR-D-CNP-0016.

##### Consequences

###### Satu pintu tulis tidak bertahan sebagai kesepakatan — ia ditegakkan lewat grant

Menaruh skema baru bersebelahan dengan `POOLDATA` berarti menaruhnya di tempat yang sistem lain sudah terbiasa menulisinya lewat SQL langsung. Itu bukan kekhawatiran teoretis: DDL tabel work Pega memuat

```sql
GRANT ALTER, DELETE, INDEX, INSERT, REFERENCES, SELECT, UPDATE, ... TO POOLDATA;
```

`POOLDATA` memegang `INSERT`, `UPDATE`, `DELETE`, bahkan `ALTER` atas tabel milik Pega. Apakah hak itu benar-benar dipakai belum diketahui (REQ-021, BLOCKER) — tetapi hak yang ada akan dipakai cepat atau lambat, dan **ADR-D-CNP-0017 yang hanya berupa kesepakatan akan runtuh di tempat seperti ini.**

Maka penegakannya masuk DDL, bukan catatan:

| Pihak | Hak |
|---|---|
| Akun pemilik skema | seluruhnya — hanya dipakai saat pemasangan dan migrasi |
| Akun aplikasi | `INSERT`, `UPDATE`, `DELETE`, `SELECT` atas tabel kanonik |
| Seluruh akun lain, termasuk `POOLDATA` | `SELECT` **hanya atas view** ADR-D-CNP-0023 — tidak ada satu pun hak atas tabel kanonik |

`GRANT` dan `REVOKE`-nya ditulis di DDL sebagai objek tersendiri, bukan dijalankan tangan. Hak yang tidak tertulis di DDL adalah hak yang tidak dapat diaudit.

###### Versi 12.1 sebagai dasar, dengan penanda di tempat yang lebih ringkas di atasnya

DDL ditulis supaya berjalan di 12.1. Setiap tempat yang lebih ringkas di versi lebih baru ditandai di spec, tidak dipakai diam-diam. Yang menggigit paling awal adalah **batas panjang pengenal**: 30 byte sampai 12.1, 128 byte sejak 12.2. Nama constraint dan index menabraknya lebih dulu daripada nama tabel.

Selama versi pasti belum diketahui, spec wajib memuat **tabel singkatan tertutup** — daftar tetap, ditulis sekali, tidak boleh ditambah saat menulis DDL. Singkatan ad-hoc adalah cara `UR` dan `MDP` lahir, dan itu persis yang glosarium tutup.

###### Tidak ada kolom JSON untuk data yang dimiliki

Seluruh alasan migrasi ini adalah memberi tiap nilai kolomnya sendiri. Sistem lama menyimpan nilai uang sebagai **teks di dalam JSON**, dan bahkan view `CLAIMXOL` mengeluarkannya sebagai `varchar2` (`BLUEPRINT.md` §13.2, §13.5).

JSON hanya boleh muncul di **tabel pendaratan migrasi**. Itu menjadikan tabel pendaratan satu-satunya tempat bentuk lama boleh masuk utuh, dan membuat aturan penguraian ADR-D-CNP-0014 bekerja di satu batas yang jelas — antara pendaratan dan kanonik — bukan tersebar di banyak tempat.

###### Yang diwarisi ADR lain

**ADR-D-CNP-0016** tetap berlaku apa adanya: data pegawai masuk lewat API atau salinan tersinkron. Instance yang sama tidak mengubahnya — `hrdasm.v_hrd_mst@asmd.sinarmas.co.id` yang ditemukan di `V_MST_USER_TEKNIS` adalah link ke mesin **lain**, dan itu yang dilarang.

**ADR-D-CNP-0023** menjadi dapat dilaksanakan: view lintas skema di satu instance adalah view sungguhan, tanpa salinan dan tanpa link.

**ADR-D-CNP-0017** berpindah dari kesepakatan menjadi konfigurasi yang bisa diperiksa, lewat tabel grant di atas.

##### Akibat yang lahir bersama kenaikan status ini

###### D20 naik dari catatan menjadi pekerjaan

Selama dokumen ini `proposed`, temuan **D20** — `POOLDATA` sudah memakai database link ke instance lain (`hrdasm.v_hrd_mst@asmd.sinarmas.co.id` di `V_MST_USER_TEKNIS`) — hanya catatan tentang sistem lama. Sekarang skema baru benar-benar berdiri di instance itu, dan ADR-D-CNP-0016 melarang apa yang sudah berjalan di sebelahnya.

~~**K1 harus diputuskan sebelum tiket `35` dikerjakan**, dan kedua cabangnya dibawa ke pemilik keputusan — **tidak dipilih sendiri**.~~

**Gerbang itu gugur 19 September 2026.** K1 ditutup sebagai **K1a**: ADR-D-CNP-0016 tetap berlaku tanpa revisi, karena modul ini **tidak mengambil data HRD sama sekali** — sehingga larangan database link tidak pernah bertabrakan dengan kebutuhan kita. Kekhawatiran yang dicatat di sini, bahwa K1b *"melemahkan alasan yang menopang ADR-D-CNP-0023 dan ADR-D-CNP-0028"*, **tidak terjadi**: K1b tidak dipilih, dan kedua ADR itu tidak ditinjau ulang. Rinciannya di badan ADR-D-CNP-0016.

###### Penegakan ADR-D-CNP-0017 lewat grant adalah pekerjaan, bukan catatan

Satu pintu tulis tidak bertahan sebagai kesepakatan. Penegakannya lewat `GRANT` dan `REVOKE` masuk lingkup tiket `01` **sebagai bagian pekerjaannya**, bukan catatan operasional yang menyusul. Alasannya satu kalimat: **hak yang tidak tertulis tidak dapat diaudit.**

### ADR-D-CNP-0029 - Kurs yang dipakai ikut tersimpan bersama nilai yang dikonversinya

> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, empat lapisan; `pengetahuan/ddl/` 49 berkas. Sapuan S4 dan S10, dijalankan 18 September 2026 — `SAPUAN-S3-S9.md` dan `SAPUAN-S10-S16.md`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut.

> **Naik ke `accepted` 19 September 2026, dan kedua cabangnya dicabut.** Versi pertama berstatus `proposed` **justru karena** dua cabang terbuka (A dan B) yang menunggu jawaban akuntansi atas `ASK-AKUNTANSI` butir 7. Cabang itu tidak dijawab — ia **dihapuskan** oleh keputusan seragam di bawah, sehingga alasan satu-satunya ADR ini berstatus `proposed` ikut hilang. Butir 7 tetap hidup sebagai **verifikasi**, dan tidak menahan ADR ini maupun tiket mana pun.

**Setiap nilai uang, tanpa pengecualian, disimpan lengkap dengan lima hal di baris yang sama**: nilai asli, kode mata uangnya, nilai rupiahnya, kurs yang dipakai menghasilkan nilai rupiah itu, dan asal-usul kurs tersebut.

Nilai rupiah tanpa kurs tidak dapat lahir; kurs tanpa nilai rupiah boleh ada; kolom rupiah yang kosong boleh ada.

##### Mengapa bukan sekadar kehati-hatian

Tiga bukti yang berdiri sendiri-sendiri.

**Pertama, sistem lama sudah melakukannya di satu tempat.** `POOLDATA.CLAIMXOL2` memuat `KursIDR` **berdampingan** dengan nilai yang dikonversinya, di satu baris. Begitu pula `SpreadingRisk` dan `CNPCurrencyList` di sisi Pega, yang membawa `KursIDR` bersama `ClaimAmountIDR`.

**Kedua, dan ini yang menentukan: kolom rupiah dapat terisi tanpa kurs pernah dipakai.** Satu kolom yang sama diisi oleh dua rule dengan perlakuan berbeda:

| Rule | Perlakuan |
|---|---|
| `CountLossAllocation_act` baris 9585 | menguji mata uang lebih dulu, lalu mengalikan kurs bila perlu |
| `SetActualPremium_ACT` baris 681 | **menyalin nilai apa adanya**, tanpa menguji mata uang, tanpa mengalikan kurs |

Bila nilai sumbernya bukan rupiah, kolom berlabel rupiah memuat angka yang bukan rupiah — dan **tidak ada apa pun di baris itu yang memperlihatkannya**. Menyimpan kursnya membuat keadaan itu terbaca alih-alih tersembunyi. (D23)

**Ketiga, rantai nilai melewati beberapa berkas tanpa menyatakan skala.** Rantai `.GrossValue` ditelusuri dari lahir sampai keluar dan tidak menyatakan skala di satu titik pun, sementara rantai `.AdjusterFeeValue` menyatakan skala di setiap pembagian. Nilai yang tidak membawa jejak perlakuannya tidak dapat direkonsiliasi belakangan. (D24)

##### Considered Options

**Menyimpan nilai rupiah saja** ditolak: itu keadaan sistem lama, dan D23 memperlihatkan keadaan itu sudah menghasilkan kolom yang tidak dapat dipercaya tanpa membaca kode yang mengisinya.

**Menyimpan kurs di tabel kurs terpisah dan menghubungkannya lewat tanggal** ditolak sebagai penggantinya, meski tabel kurs bertanggal tetap dibutuhkan untuk keperluan lain. Alasannya: kurs yang **dipakai** pada satu perhitungan belum tentu kurs yang **berlaku** pada tanggal itu. Yang perlu direkonsiliasi adalah yang dipakai, bukan yang seharusnya.

##### Keberlakuan seragam — keputusan 19 September 2026, menggantikan cabang A dan B

Sapuan S4 menemukan **tujuh baris, sembilan nama** properti bernilai uang yang **tidak punya padanan rupiah sama sekali**: `AdjusterFee`, `AdjusterFeeValue`, `Salvage`, `SalvageValue`, `CNPOthersFee`, `OthersFee`, `AdjustmentValue`, `TotalClaim`, `PremiumSpreaded`. Penamaan alternatif (`Rupiah`, `Idr`, `_IDR`, `Rp`) ikut disapu dan mengembalikan nol; yang ada hanya sufiks `RNM`, yaitu porsi pihak, bukan mata uang. Lima dari sembilan masuk hitungan instruksi bayar, dan rantai `.AdjusterFeeValue` terbukti tidak melewati satu pun titik konversi.

Versi pertama membaca temuan itu sebagai dua cabang — **A**: ketiadaan itu kesengajaan, sehingga ADR ini tidak berlaku untuk kelas itu; **B**: ketiadaan itu kekurangan, sehingga kelas itu memerlukan padanan rupiah. **Keduanya dicabut.** ADR ini berlaku untuk **seluruh** nilai uang, termasuk kesembilan nama di atas.

**Alasannya perbandingan biaya, bukan pengetahuan baru.** Kedua cabang tetap tidak terjawab dari XML; yang berubah adalah kesadaran bahwa keduanya **tidak berbiaya setara**:

| Bila ternyata… | dan kita menyediakan kolomnya | dan kita tidak menyediakannya |
|---|---|---|
| rupiahnya **tidak pernah** dibutuhkan | kolom kosong — tidak ada yang rugi | benar, kebetulan |
| rupiahnya **dibutuhkan** | benar | **migrasi kedua** |

Satu sisi berbiaya kolom kosong; sisi lain berbiaya migrasi kedua. Menunggu jawaban akuntansi untuk memilih di antara keduanya berarti menahan seluruh model demi menghindari biaya yang lebih kecil dari biaya menunggunya.

Ini bentuk yang sama dengan ADR-D-CNP-0025, yang menyediakan kolom lingkup sebelum ada yang memerlukannya.

**`ASK-AKUNTANSI` butir 7 tetap berdiri**, dan kedudukannya berubah: ia **memverifikasi** apakah kesembilan kolom itu terisi atau tetap kosong. Ia tidak lagi menahan ADR ini, tidak menahan `SPEC-MODEL-DATA.md`, dan tidak menahan tiket mana pun.

##### Consequences

**Tidak ada lagi rancangan yang ditunda oleh ADR ini.** Pembatasan versi pertama — *"tidak ada tabel, kolom, aturan, atau uji yang boleh dirancang dengan mengandaikan cabang A maupun cabang B"* — dicabut bersama cabangnya.

Yang berlaku sekarang seragam: **setiap** nama bernilai uang mendapat kelima medan, baik yang hari ini sudah punya padanan rupiah — `ClaimAmountIDR`, `GrossValueIDR`, `ValueIDR`, `AmountIDR`, `TotalSumInsuredIDR` — maupun kesembilan yang hari ini tidak punya.

**Dua akibat yang harus dilihat terpisah.** Yang pertama: kolom rupiah untuk kesembilan nama itu akan **kosong pada data migrasi**, karena nilainya memang tidak pernah ada di sistem lama — kosong, bukan nol (ADR-D-CNP-0019). Yang kedua: struktur menerima nilai itu bila kelak ada yang menghitungnya, tanpa perubahan skema.

ADR ini juga menutup pertanyaan satuan pada `.ValueAdjustment` yang datang dari luar: nilai yang menyeberang batas **wajib membawa mata uang dan kursnya**, dan nilai tanpa mata uang **ditolak di batas**, bukan diterima lalu ditebak. Itu yang membuat H1/H2 berhenti menjadi masalah rancangan — lihat `BLUEPRINT.md` §8.7.

ADR ini juga tidak menutup D23. Kolom rupiah yang terisi tanpa konversi adalah cacat sistem lama; apakah nilai lamanya diperbaiki, dibiarkan, atau didaftar mengikuti **AK-2a** dan **AK-2b**, dan itu putusan tersendiri.

## Modul komite-claim-non-prop

### ADR-D-KCNP-0030 - Satu penentu jenjang aktif

> Modul  : Komite Claim Non Prop · Ronde 06 · 2026-09-20
> Peran  : juru catat
> Masukan: KETETAPAN.md K5-1 · GRILL-05/01-TEMUAN.md N-01, N-02 · GRILL-05/06-PUTUSAN.md T5-01, T5-02
> Status : TERBUKA
> Sifat  : HIDUP



**Status:** Diterima · 2026-09-20 · sumber normatif `KETETAPAN.md` `K5-1`

##### Konteks

Sistem lama memakai **dua** penentu jenjang aktif di dalam satu rule yang sama,
`KomiteRouter`:

- **Hitungan.** Langkah 1–4 memilih sasaran penugasan dari `.KomiteCount`, sebuah bilangan
  yang lahir bernilai 1 dan bertambah satu di `KomitePostAdjustment` langkah 29.
- **Keadaan baris.** Langkah 6 beriterasi atas daftar jenjang, berhenti pada baris pertama
  ber-`KomiteAproval==0` (transisi `6/6`, keluar iterasi), dan menugaskan ke pemegangnya.

Keduanya tidak pernah saling memeriksa, dan langkah 6 menang karena berjalan belakangan.

Hitungan itu juga dipakai untuk hal ketiga: penutupan case. `KomitePostAdjustment` langkah
20 menetapkan `KomiteCount = KomiteLoop` ketika hasilnya penolakan — bukan untuk mencatat
apa pun, melainkan agar pra-syarat langkah 28 (`KomiteCount >= KomiteLoop`) terpenuhi dan
`ASMForceCaseClose` terpanggil. Penutupan sirkulasi dengan demikian bergantung pada
perbandingan dua bilangan, dan salah satunya sengaja dipalsukan.

Bilangan pembandingnya, `KomiteLoop`, adalah jumlah baris **hasil laporan** roster, bukan
jumlah jenjang yang benar-benar terbentuk. Pada jalur akseptasi bersyarat daftar jenjang
tidak dibersihkan sebelum diisi, sehingga jumlah anggota tumbuh sementara `KomiteLoop`
tidak. Sejak itu kedua bilangan berbeda, dan tiga hal ikut salah sekaligus: kapan giliran
berhenti, siapa yang tidak pernah mendapat giliran, dan apakah case tertutup.

Keadaan keputusan sendiri bersifat sesaat: `KomiteRouter` langkah 5 mengosongkan
`.AcceptStatus` tanpa pra-syarat, sedangkan gerbang lingkar menuntutnya bernilai `"1"`.

##### Keputusan

Jenjang aktif ditentukan oleh **satu hal saja**: jenjang berderajat terendah yang belum
memutuskan.

Hitungan jenjang adalah **turunan** yang dihitung dari daftar jenjang. Ia tidak pernah
menjadi syarat penugasan, tidak pernah menjadi syarat penutupan, dan tidak pernah disimpan
sebagai penyimpan kedua yang dapat berselisih dengan daftar jenjangnya.

Sirkulasi ditutup dari **keadaan keputusan** — tidak ada jenjang tersisa yang belum
memutuskan, atau sebuah keputusan mengakhiri sirkulasi menurut jenisnya — bukan dari
perbandingan dua bilangan.

Keputusan tiap jenjang adalah **catatan tetap** pada jenjang itu, bukan satu medan yang
dipakai bergantian oleh seluruh sirkulasi.

##### Konsekuensi

- Urutan giliran mengikuti derajat naik, dan itu **paritas** — sistem lama sudah berperilaku
  demikian. Uji `P5-01` dan `P5-02` dirancang lulus.
- Keadaan "jenjang aktif menurut hitungan berbeda dari jenjang aktif menurut daftar" menjadi
  **tidak dapat direpresentasikan**. Uji `P5-02b` karena itu dirancang gagal, dan itu
  dinyatakan, bukan disembunyikan.
- Empat nama keranjang yang ditulis di dalam rule hilang bersama penentu keduanya; sasaran
  penugasan menjadi data per jenjang.
- Riwayat keputusan menjadi dapat dibaca ulang setelah sirkulasi selesai — kemampuan yang
  sistem lama tidak punya.
- Migrasi data harus menghitung ulang jenjang aktif dari daftar jenjang, dan **tidak boleh**
  mempercayai `KomiteCount` yang tersimpan.

##### Alternatif yang ditolak

**Membawa kedua penentu dan menambahkan pemeriksaan konsistensi.** Ditolak: pemeriksaan
konsistensi hanya memberi tahu bahwa keduanya sudah berselisih, tanpa memberi tahu yang mana
yang benar. Sistem lama sudah menjalankan percobaan itu selama bertahun-tahun.

**Membawa hitungan sebagai penyimpan dan menurunkan daftar jenjang darinya.** Ditolak: arah
turunannya terbalik terhadap sumber kebenarannya. Siapa yang berhak memutus adalah fakta
tentang orang dan derajat, bukan tentang bilangan.

**Menutup sirkulasi dengan penanda "sudah selesai" yang ditulis pemutus terakhir.** Ditolak:
ia mengulang pola yang sama — satu fakta disimpan di dua tempat — dan menyerahkan kebenaran
penutupan kepada langkah yang dapat gagal.

### ADR-D-KCNP-0031 - Pembuatan sirkulasi dan akibatnya adalah satu transaksi

> Modul  : Komite Claim Non Prop · Ronde 06 · 2026-09-20
> Peran  : juru catat
> Masukan: KETETAPAN.md K5-2, K5-3, K5-7, H-5 · GRILL-05/01-TEMUAN.md N-07, N-08 · GRILL-05/06-PUTUSAN.md T5-03
> Status : TERBUKA
> Sifat  : HIDUP



**Status:** Diterima · 2026-09-20 · sumber normatif `KETETAPAN.md` `K5-2`, didampingi
`K5-3` dan `K5-7`

##### Konteks

Pada kedua activity pembuat sirkulasi, pemanggilan `pxAddChildWork` dijaga pra-syarat
`@hasMessages(myStepPage)` dengan aksi **lewati langkah** — bukan keluar activity.

Akibatnya, ketika pembuatan tidak jadi, langkah-langkah sesudahnya tetap berjalan seolah ia
berhasil. Keduanya membaca `pxCoveredInsKeys(<LAST>)` — case tercakup terakhir — dan menulis
potongannya ke klaim sebagai nomor komite, lalu menyimpan klaim, lalu mengirim surat. Bila
klaim pernah punya sirkulasi sebelumnya, yang tertulis adalah nomor **sirkulasi lain**. Bila
belum pernah, yang tertulis kosong.

Jalur kegagalan kedua lebih senyap lagi. `CreateChildKomiteCNP_Act` langkah 28 melompat ke
label `END` ketika daftar penyebaran risiko kosong; label itu menetapkan
`ClaimData.IsFlagError = "true"` dan menyimpan klaim. Nama properti itu muncul di **satu**
berkas dari 338 — berkas yang menulisnya. Tidak ada satu pun rule dalam ekspor yang
membacanya. Pengguna menekan tombol, tidak ada sirkulasi yang lahir, dan tidak ada apa pun
yang memberitahunya.

Persoalan ketiga menyertai keduanya: penanda maksud ditulis ke klaim **pada langkah yang
sama** yang membuat halaman anak, yaitu sebelum komite memutuskan apa pun. Klaim membawa
`IsCloseFile` atau `IsReject` sejak pengajuan — dan tetap membawanya ketika sirkulasinya
gagal lahir.

##### Keputusan

Pembuatan sirkulasi, penulisan nomornya ke klaim, penyimpanan klaim, dan pengiriman
pemberitahuan **berhasil bersama atau gagal bersama**. Tidak ada langkah sesudah pembuatan
yang boleh berjalan ketika pembuatan tidak jadi.

Setiap jalur pembentukan yang berakhir tanpa sirkulasi berakhir dengan **kegagalan yang
terlihat** oleh pengguna dan tercatat. Penanda yang tidak dibaca siapa pun bukan penanganan
galat.

**Maksud dan akibat adalah dua penanda berbeda.** "Diajukan untuk ditutup" tercatat sejak
pengajuan; "ditutup" hanya lahir dari keputusan komite lewat fungsi akibat. Sirkulasi yang
gagal lahir meninggalkan maksudnya tercatat dan akibatnya tidak.

Tali antara klaim dan sirkulasi adalah relasi **dari sisi sirkulasi**, bukan nomor yang
disalin ke klaim.

##### Konsekuensi

- Validasi menolak kasus-guna pembuatan, bukan menandai halaman — `H-5` sudah menetapkan
  ini; ADR ini memperluasnya ke seluruh urutan yang bergantung pada pembuatan.
- Pengguna yang gagal membentuk sirkulasi mendapat pesan, bukan kesunyian. Uji `P5-07` dan
  `P5-08` dirancang gagal.
- Klaim tidak lagi menyimpan nomor komite sebagai tali. Migrasi data **tidak boleh
  mempercayai** `.KomiteNo` yang tersimpan; berapa banyak klaim di produksi membawa nomor
  bukan miliknya adalah cacah yang belum diambil (`INVENTARIS-BUKTI.md` §2.5 baris 6).
- Penanda maksud dan penanda akibat menempati kolom berbeda, sehingga klaim yang pernah
  diajukan untuk ditutup tetapi tidak jadi ditutup dapat dibedakan dari klaim yang ditutup.

##### Alternatif yang ditolak

**Menjaga tiap langkah sesudah pembuatan dengan pra-syarat "bila pembuatan berhasil".**
Ditolak: itu memindahkan kewajiban kepada penulis tiap langkah berikutnya, dan langkah
berikutnya akan selalu bertambah. Sistem lama menjalankan bentuk ini dengan satu penjaga,
dan penjaganya menjaga satu langkah.

**Membiarkan penulisan nomor tetap berjalan dan memperbaikinya belakangan lewat pekerjaan
rekonsiliasi.** Ditolak: klaim yang tersimpan dengan nomor milik sirkulasi lain sudah
mengirim surat atas nomor itu. Rekonsiliasi memperbaiki kolom, bukan surat.

**Menampilkan `IsFlagError` di layar dan menganggapnya selesai.** Ditolak: yang salah bukan
bahwa penandanya tak terlihat, melainkan bahwa jalur itu berakhir dengan penanda alih-alih
dengan penolakan.

### ADR-D-KCNP-0032 - Komite adalah modul di dalam konteks Klaim, bukan konteks berjaringan

> Modul  : Komite Claim Non Prop · Ronde 06 · 2026-09-20
> Peran  : juru catat
> Masukan: KETETAPAN.md D-2, D-3 · GRILL-05/06-PUTUSAN.md T5-04 · ADR-D-CNP-0016 (sisi Claim)
> Status : TERBUKA
> Sifat  : HIDUP



**Status:** Diterima · 2026-09-20 · sumber normatif `KETETAPAN.md` `D-2`

##### Konteks

Di sistem lama, komite adalah **case anak** dari klaim. Keterhubungan keduanya bukan
panggilan antar-sistem melainkan penyuntingan halaman bersama: activity keputusan membuka
klaim induk dengan kunci (`Obj-Open-By-Handle` atas `pxCoverInsKey`), memutakhirkan
`TempMainWork.ClaimData.AdjustmentList(...)`, lalu menyimpan dan melakukan commit.

Pola itu tidak berdiri di satu berkas. Kedua activity pembuat sirkulasi menjalankan
`Obj-Refresh-And-Lock` atas klaim induk **dua kali** — sekali di awal dan sekali menjelang
akhir — dengan penyimpanan di antaranya. Batas antara "apa yang milik komite" dan "apa yang
milik klaim" tidak pernah digambar; yang ada hanyalah halaman yang keduanya sunting.

Sementara itu langkah commit yang sama menjalankan enam efek ke luar berturut-turut tanpa
kompensasi, dan `ADR-D-CNP-0016` sisi Claim sudah melarang tautan basis data lintas sistem.

Menjadikan komite konteks berjaringan akan berarti menggambar batas di tempat yang sistem
lama tidak pernah gambar, lalu menanggung konsistensi akhir pada perkara yang selama ini
bertransaksi bersama.

##### Keputusan

Komite adalah **modul di dalam konteks Klaim**, dengan agregat sendiri.

Sirkulasi, jenjang, dan keputusan membentuk satu agregat yang berubah bersama atau tidak
sama sekali. Perubahan yang menyeberang ke klaim — penanda maksud, penanda akibat, keadaan
klaim — terjadi di dalam batas transaksi yang sama, bukan lewat pesan yang mungkin sampai.

Batas transaksi itu menutupi **kedua** activity pembuat sirkulasi dan activity keputusan,
bukan hanya yang terakhir.

Efek ke luar yang tidak dapat dibatalkan (keputusan beku no. 6) berada di luar batas
transaksi. Jenis dan penanganannya per aliran ditetapkan setelah `PG-04`, `PG-05`, dan
`PG-06` dibuka.

##### Konsekuensi

- Tidak ada tautan basis data, tidak ada panggilan sinkron antar-konteks, dan tidak ada
  konsistensi akhir antara sirkulasi dan klaim. `ADR-D-CNP-0016` tetap berlaku penuh.
- Skema `KLAIMNP` menampung objek komite; tidak ada skema terpisah dan tidak ada basis data
  terpisah.
- Penyuntingan halaman bersama digantikan oleh kepemilikan yang tegas: sirkulasi memiliki
  jenjang dan keputusan; klaim memiliki usulan dan penandanya. Tidak ada satu pun medan
  yang dimiliki keduanya.
- Pola "buka dengan kunci, sunting, buka dengan kunci lagi" tidak dibawa. Perkara `K-11`
  tidak dibuka kembali; yang berubah adalah bahwa batas transaksinya kini diputuskan,
  bukan diwarisi.
- Enam efek luar keluar dari langkah commit — itu memang sudah diperintahkan keputusan beku
  no. 6; ADR ini menyatakan bahwa mereka berada di luar batas transaksi, bukan ke mana
  mereka pindah.

##### Alternatif yang ditolak

**Komite sebagai konteks berjaringan dengan pesan antar-konteks.** Ditolak: batasnya harus
digambar di tempat yang sistem lama tidak pernah gambar, dan perkara yang selama ini
bertransaksi bersama akan menanggung konsistensi akhir tanpa satu pun kebutuhan yang
memintanya.

**Komite sebagai bagian dari agregat klaim, tanpa agregat sendiri.** Ditolak: sebuah klaim
dapat memiliki banyak sirkulasi yang hidup bersamaan, dan menjadikan seluruhnya satu agregat
membuat dua sirkulasi atas dua usulan berbeda saling mengunci tanpa sebab.

**Menunda keputusan sampai beban nyata terukur.** Ditolak: batas agregat menentukan bentuk
tabel dan bentuk transaksi, dan keduanya ditulis sebelum beban apa pun terukur. Menunda
berarti memilih tanpa menyatakannya.

### ADR-D-KCNP-0033 - Jalur menutup dan menolak klaim

> Modul  : Komite Claim Non Prop · Tahap 3 · 2026-09-21
> Peran  : juru catat
> Masukan: KETETAPAN.md K6-4 · KETETAPAN.md §8 F-20 · REGISTER-DEVIASI.md bagian 6 (deviasi 23, 24, 25) · SPEC-KOMITE-01.md S-001 dan S-008
> Status : TERBUKA
> Sifat  : HIDUP



**Status:** Diterima · 2026-09-21 · sumber normatif `KETETAPAN.md` `K6-4`

##### Konteks

Modul komite punya dua jalur yang mekanismenya identik sampai ke tingkat langkah tetapi
akibatnya berbeda: mengedarkan **usulan pembayaran**, dan mengedarkan **usulan menutup atau
menolak klaim**. Jalur kedua lahir dari rule tersendiri, `CreateChildKomiteCloseNP_Act`, dan
cara ia menentukan siapa memutus tidak punya kemiripan dengan jalur pertama.

Jalur pembayaran menyusun jenjangnya dari daftar anggota di basis data, disaring oleh ambang
kewenangan. Jalur menutup dan menolak tidak melakukan keduanya. Ia **tidak membaca daftar
anggota sama sekali** dan **tidak mengenal ambang**. Sebagai gantinya, langkah 8 dan 9 menulis
pemegangnya langsung di dalam rule: **nama orang, alamat surel, inisial, dan dua sebutan
jabatan yang berbeda untuk satu kedudukan yang sama** (`F-20`). Empat bentuk pengenal untuk
satu orang, seluruhnya tetapan di dalam kode.

Akibatnya berlipat. Ketika orang itu pindah jabatan, cuti, atau berhenti, tidak ada tempat
untuk mengubahnya selain kode — dan karena satu kedudukan ditulis dengan dua sebutan jabatan
yang berbeda, mengubahnya di satu tempat meninggalkan yang lain. Keadaan "tidak ada pemegang"
juga **tidak dapat terjadi** di jalur ini: pemegangnya tetapan, sehingga sirkulasi selalu
lahir, bahkan ketika orangnya sudah tidak berwenang.

Jenis sirkulasinya pun tidak dinyatakan. Langkah 12 membedakan menutup dari menolak dengan
membaca **nilai sebuah kolom analisis** — kolom yang diisi untuk keperluan lain, dan yang
dapat berselisih dengan tindakan yang benar-benar dipilih pengaju tanpa ada satu pun jalur
yang memeriksanya. Dua maksud berbeda karena itu dapat tertukar diam-diam.

Ketetapan ronde 4 yang mengatur jalur ini, `H-2`, **bunyinya hilang** bersama berkas yang
memuatnya. Yang tersisa hanyalah catatan bahwa keputusan itu pernah diambil. Itu sendiri
alasan untuk menuliskannya sebagai ADR: keputusan ini sudah sekali kehilangan bunyinya.

##### Keputusan

Sirkulasi menutup klaim dan menolak klaim **tetap satu jenjang**. Kedalamannya tidak
disamakan dengan jalur pembayaran.

Pemegang jenjang itu **diselesaikan dari daftar anggota berdasarkan peran**, bukan dari nama
yang ditulis di dalam aturan. Peran adalah atribut baris daftar anggota; nama orang tidak
muncul di jalur mana pun.

Bila **tidak ada baris daftar anggota yang aktif** memegang peran itu, sirkulasi **ditolak
saat dibuat**, dengan galat yang menyebut peran yang kosong.

**Jenis sirkulasi dinyatakan oleh pemanggil.** Ia tidak diturunkan dari nilai kolom analisis
mana pun.

##### Konsekuensi

- Mutasi jabatan menjadi penyuntingan data. Tidak ada permintaan perubahan kode ketika orang
  yang menandatangani penutupan klaim berganti.
- Satu kedudukan punya **satu** sebutan. Dua sebutan berbeda untuk satu kedudukan tidak dapat
  lahir kembali, karena sebutannya kini kolom, bukan dua tetapan di dua langkah.
- Keadaan "tidak ada pemegang" menjadi **dapat terjadi dan terlihat**. Sistem lama tidak dapat
  menghasilkannya; sistem baru menolak pembentukan alih-alih melahirkan sirkulasi yang tidak
  ada yang berwenang memutuskannya.
- Dua maksud berbeda **tidak dapat tertukar**. Menutup dan menolak dibedakan oleh tindakan
  yang dipilih pengaju, bukan oleh kolom yang diisi untuk keperluan lain.
- Tiga perbedaan perilaku terhadap sistem lama, dan ketiganya **disengaja**: deviasi 23
  (jenis dinyatakan pemanggil), 24 (pemegang dari peran), dan 25 (penolakan saat peran
  kosong). Masing-masing membawa uji yang **dirancang gagal** terhadap sistem lama —
  `PS-01`, `PS-02`, `PS-03`. Uji yang lulus paritas adalah tanda deviasinya tidak terpasang.
- Jalur ini karena itu **tidak dapat dibangun sebelum daftar anggota berdiri sebagai data**.
  Urutan itu mengikat.
- `H-2` digantikan. Ia tetap tercatat `BUNYI HILANG`; ADR ini tidak memulihkannya, melainkan
  menggantikannya dengan bunyi yang ditulis dari langkah yang terbaca.

##### Alternatif yang ditolak

**Menyamakan kedalamannya dengan jalur pembayaran — banyak jenjang, disaring ambang.**
Ditolak: menutup klaim adalah keputusan yang **membatalkan**, bukan yang mengeluarkan uang.
Kedalaman berjenjang ada untuk menjaga pengeluaran, dan memakainya di tempat yang tidak
mengeluarkan apa pun menambah persetujuan tanpa menambah penjagaan. Sistem lama pun
memperlakukannya satu jenjang; yang kami ubah adalah **siapa** jenjang itu, bukan berapa.

**Menurunkan jenis sirkulasi dari kolom analisis, sebagaimana sistem lama — paritas.**
Ditolak: paritas terhadap perilaku yang dapat berselisih dengan maksud pengaju adalah paritas
terhadap cacat. Kolom itu diisi untuk keperluan lain, dan tidak ada satu pun jalur yang
memeriksa apakah isinya sejalan dengan tindakan yang dipilih. Mempertahankannya berarti
mewarisi kemungkinan tertukar tanpa mewarisi satu pun manfaat.

**Memindahkan nama orang dari kode ke berkas konfigurasi.** Ditolak: ia memindahkan masalahnya
tanpa menyelesaikannya. Yang salah bukan **tempat** nama itu ditulis, melainkan bahwa
kewenangan melekat pada **orang** alih-alih pada **kedudukan**. Konfigurasi berisi nama orang
tetap memerlukan penyuntingan setiap mutasi, tetap dapat berselisih dengan daftar anggota
yang sebenarnya, dan tetap melanggar larangan bercabang berdasarkan identitas orang.

**Membiarkan sirkulasi lahir tanpa pemegang, lalu menugaskannya belakangan.** Ditolak: ia
mengulang persis cacat yang sudah ada di jalur pembayaran — berkas yang lahir lalu diam,
tanpa ada yang tahu ia menunggu. Sirkulasi tanpa jenjang ditolak saat dibuat, dan jalur ini
tidak dikecualikan.

## Modul treaty-in

### ADR-D-TI-0034 - Arsip JSON sistem lama tidak punya jalur baca

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In

##### Konteks

Sistem lama menyimpan seluruh halaman clipboard sebagai satu dokumen JSON di `M_TREATY_IN.JSONDATA`
(dan `M_TREATY_IN_EDM.JSONDATA` untuk addendum). Dokumen itu adalah bentuk kanonik sistem lama.

Saat migrasi, dokumen-dokumen itu akan tetap ada. Godaannya jelas: menyimpannya "untuk jaga-jaga"
dan membiarkan aplikasi baru membacanya ketika ada yang tidak ketemu di model baru.

##### Keputusan

Arsip JSON sistem lama **tidak punya jalur baca apa pun** di aplikasi hasil migrasi. Bukan
"dihindari", bukan "hanya untuk kasus khusus" — **tidak ada kode yang membacanya**.

Arsip itu punya **masa simpan** yang ditetapkan, dan setelah masa itu lewat ia dihapus.

##### Alasan

Jalur baca cadangan ke bentuk lama menciptakan bentuk kanonik kedua, yang dilarang ADR-D-CNP-0023. Lebih
buruk lagi, ia menyembunyikan kegagalan migrasi: data yang tidak ikut pindah tidak akan pernah
ketahuan selama aplikasi masih bisa mengambilnya dari arsip.

Tanpa jalur baca, setiap data yang dibutuhkan **harus** ada di model baru, dan ketidakhadirannya
muncul sebagai kegagalan yang terlihat.

##### Konsekuensi

- Inventaris struktur data harus lengkap sebelum migrasi, karena tidak ada jaring pengaman.
- Arsip disimpan di luar jangkauan aplikasi, dengan akses bernama untuk keperluan forensik saja.
- Masa simpan harus ditetapkan oleh pemilik proses dan kepatuhan, bukan oleh tim teknis.

### ADR-D-TI-0035 - Kegagalan tidak pernah disamarkan menjadi nilai

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In, dan diusulkan untuk seluruh modul

##### Konteks

Lima temuan terpisah di modul ini ternyata satu penyakit:

| Keadaan | Yang dilakukan sistem lama |
|---|---|
| Limit nol saat menghitung Rate on Line | `ROLPct` diisi `9989998`, tersimpan sebagai `998999800` |
| Mata uang di luar slot IDR/USD pada reinstatement | hasil diisi `0` |
| Kurs valas tidak ditemukan | pengali tinggal kosong, nilai IDR menjadi `0` |
| Prosedur Oracle gagal | `ErrMsg` dan `StsSave` tidak diisi; Pega menerima kosong |
| Persentase tidak terurai dari `VARCHAR2(1000)` | nilai lolos apa adanya |

Nilai samarannya selalu tampak masuk akal — nol, satu, atau angka besar — sehingga hilir tidak
punya cara membedakannya dari hasil sungguhan.

##### Keputusan

1. Besaran yang **tidak dapat dihitung** tidak pernah diganti dengan nilai pengganti. Ia bernilai
   tidak-ada, dan tipe datanya harus **memaksa** pembaca menangani ketidakadaan itu.
2. Ketidakadaan itu **membawa sebabnya**: limit nol, kurs tidak ditemukan, mata uang tidak dikenal,
   nilai tidak terurai. Sebab yang berbeda menuntut tindakan yang berbeda.
3. Nilai pengganti yang netral — nol, satu, string kosong — **dilarang secara khusus**, karena
   justru nilai netral yang paling mudah lolos tanpa disadari.
4. Di perbatasan tempat masukan tidak bisa ditolak, kegagalan **diangkat ke pemanggil**, tidak
   ditelan.

##### Catatan penerapan

Aturan ini **tidak** berarti setiap ketidakmampuan menghitung menghentikan pekerjaan orang. Kontrak
yang limitnya belum diisi memang belum bisa punya Rate on Line, dan itu normal selama penyusunan.

Yang dilarang adalah **menyimpan angka palsu** untuk keadaan itu. Bedakan:

- *belum bisa dihitung, dan itu wajar di tahap ini* — bukan angka, tidak menghalangi akseptasi;
- *gagal dihitung padahal seharusnya bisa* — bukan angka, **menghalangi** akseptasi.

##### Instans yang tertutup oleh ADR ini

Selain kelima di atas: pengambilan `pxResults(1)` tanpa memeriksa ada berapa hasil; pembagian
dengan bagian NuRe tanpa memeriksa ia terisi; dan penjaga rasio kerugian yang memeriksa premi
bruto sementara pembaginya premi neto. Semuanya dicatat sebagai instans, bukan temuan baru.

### ADR-D-TI-0036 - Beku saat disetujui, hitung saat dibaca

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In, dan diusulkan untuk seluruh modul

##### Konteks

Bentuk masalah yang sama muncul tiga kali dalam satu sesi:

1. **Penyebaran** — baris penyebaran dihapus dan disusun ulang dari tabel susunan baku pada setiap
   perhitungan. Tabel itu disunting orang. Kontrak yang sudah melewati persetujuan empat tingkat
   berubah pembagian kapasitasnya begitu ia disentuh lagi.
2. **Angka yang dilihat approver** — tidak ada tempat yang mencatat berapa nilai kontrak pada saat
   disetujui. Yang tersimpan hanya hasil terakhir.
3. **Kurs** — konversi memakai kurs tahun berjalan saat perhitungan dijalankan, bukan kurs yang
   berlaku pada peristiwanya. Nilai rupiah kontrak lama bergeser setiap kali dihitung ulang di
   tahun berbeda.

Ditambah satu lagi yang ditemukan belakangan: **bagian NuRe** dipakai sebagai faktor konversi
tingkat, dan bagian yang dipakai adalah bagian hari ini, bukan bagian yang berlaku saat jumlah itu
terjadi.

##### Keputusan

> Angka yang menjadi dasar suatu persetujuan **dibekukan pada saat persetujuan itu**, bersama
> seluruh masukan yang membentuknya.
>
> Angka yang menggambarkan **posisi terkini** dihitung saat dibaca dan tidak pernah disimpan
> sebagai atribut benda yang dilaporkannya.

Penerapan yang mengikat:

- Penyebaran adalah **turunan** selama kontrak belum disetujui, dan menjadi **fakta tercatat**
  begitu disetujui. Sesudah akseptasi, perubahan pada susunan baku **tidak menjalar**.
- Menyimpan **tidak pernah** menarik ulang dari data acuan. Tidak ada mode yang berpindah sendiri.
- Hasil yang dibekukan membawa **penunjuk** ke masukan yang dipakai: susunan mana, periode mana,
  kurs berapa dari sumber apa, bagian berapa, dan kapan dihitung — serta **asal-usulnya**
  (disemai dari acuan baku, diisi tangan, atau disemai lalu disunting).
- Bagian NuRe **boleh berubah** lewat addendum, tetapi perubahannya **berlaku ke depan**; ia tidak
  menulis ulang tingkat jumlah yang sudah terjadi.
- Keterlambatan jatuh tempo dihitung saat dibaca. **Penanda terlambat tidak pernah disimpan.**

##### Kenapa ini bukan sekadar kebersihan model

Pemilik proses menetapkan bahwa Treaty In Adjustment akan **mengambil data lama lewat SELECT**,
bukan menyalinnya. Itu hanya aman bila data lama tidak bergerak. Kalau nilai kontrak masih bisa
berubah sesudah akseptasi, selisih yang dihitung hari ini dan bulan depan atas kontrak yang sama
akan berbeda tanpa ada yang mengubah apa pun.

Keputusan ini karena itu adalah **prasyarat teknis** bagi cara kerja Adjustment yang sudah
ditetapkan, bukan pilihan rancangan yang bisa ditunda.

##### Konsekuensi

- Setiap versi kontrak yang disetujui harus punya identitas yang stabil dan bisa ditunjuk.
- Jejak perubahan (ADR-D-TI-0045) menjadi bukti bahwa pembekuan benar-benar terjadi. Tanpa jejak,
  pembekuan adalah janji, bukan fakta.
- Sebagian besar angka pelaporan tidak punya kolom sama sekali.

### ADR-D-TI-0037 - Masukan versus turunan; tidak ada "mode manual"

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In, dan diusulkan untuk seluruh modul

##### Konteks

Pola yang sama muncul tiga kali:

- **Penyebaran** punya jalur otomatis dan jalur manual, dipilih oleh keterisian sebuah field yang
  bisa mengosongkan dirinya sendiri.
- **Reinstatement** punya empat field yang saling memicu dalam satu cincin: persen menggerakkan
  jumlah, jumlah menggerakkan persen yang lain.
- **Kapasitas proporsional** punya jalur "man" untuk Surplus yang diketik dan jalur turunan,
  dengan sebuah sakelar yang di dua pertiga titik pemanggilan diberi makan nomor layer.

Sistem lama tidak pernah memutuskan mana masukan dan mana turunan. Setiap kali sebuah turunan
ternyata perlu disepakati, jawabannya adalah **membuka suntingan** — dan lahirlah cincin, jalur
manual, dan penimpaan. Semuanya gejala dari satu keputusan yang tidak pernah diambil.

##### Keputusan

1. Besaran yang **dinegosiasikan** adalah **masukan** dan diketik orang.
2. Besaran yang **dihitung** adalah **turunan**, tidak bisa disunting, dan tidak pernah menulis
   balik ke masukannya.
3. **"Mode manual" tidak ada.** Bila sebuah besaran boleh disepakati, ia **dinaikkan menjadi
   masukan** — bukan dibiarkan sebagai turunan yang boleh disunting.

##### Penerapan yang sudah ditetapkan

| Besaran | Sifat |
|---|---|
| Quota Share: persentase QS | masukan (satu bilangan; retensi dan sesi turunannya) |
| Surplus: retensi sebagai jumlah uang, dan jumlah lines | **dua masukan**; kapasitas turunan |
| Reinstatement: porsi limit yang dipulihkan, tarif premi reinstatement | dua masukan, dinegosiasikan terpisah; hubungannya perkalian |
| Jumlah limit dipulihkan, premi tambahan | turunan |
| Potongan berbasis premi bruto | turunan |
| Potongan berjumlah disepakati | masukan, dan **tidak diprorata** |
| Penyebaran per baris | masukan bila diisi tangan; turunan bila disemai dari acuan baku |
| Faktor prorata addendum | turunan dari tanggal, bukan centang |
| "Boleh diubah" | turunan dari keadaan siklus hidup |

##### Konsekuensi

- **Prorata mengenai dasarnya, bukan hasilnya.** Ada satu besaran yang menyusut — paparan waktu
  atas premi. Prorata itu sekali, lalu seluruh turunan dihitung ulang dari dasar yang baru. Dua
  puluh blok prorata terpisah menjadi satu perkalian, dan pertanyaan "apa yang ikut menyusut"
  lenyap bersama bentuk yang melahirkannya.
- Cincin `RetentionPct`/`CessionPct` tidak diangkut. Satu-satunya akibat nyata dari cincin itu
  adalah membuka pintu suntingan.
- `CessionPct` **pecah**: di Quota Share ia persentase sesi yang bermakna; di Surplus yang ada
  adalah jumlah lines. Keduanya tidak berbagi kolom maupun nama. Pemecahan ini berlaku **tanpa
  syarat**, bukan bergantung pada jawaban soal manual.

### ADR-D-TI-0038 - Aturan yang bisa berubah disimpan sebagai data

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In

##### Konteks

Tiga kali dalam sesi yang sama, jawaban atas pertanyaan yang berbeda berbentuk sama:

1. **Perutean persetujuan** — sistem lama menanam nama orang di dalam aturan, dan memakai kolom
   nomor telepon sebagai kode peran.
2. **Kontrak baca ke acuan luar** — bentuk yang dipakai adapter untuk membaca susunan baku NuRe.
3. **Kelengkapan transisi** — apa yang harus sudah ada untuk berpindah dari satu keadaan ke
   keadaan berikutnya.

Ketiganya berubah tanpa mengubah arti benda yang diaturnya.

##### Keputusan

Aturan yang dapat berubah **tanpa mengubah arti benda yang diaturnya** disimpan sebagai **data**,
bukan ditanam di dalam kode.

Termasuk di dalamnya: perutean dan jenjang persetujuan, batas wewenang, daftar kelengkapan per
transisi, dan pemetaan adapter ke acuan luar.

##### Yang TIDAK termasuk

Invarian — aturan yang menjaga catatan tetap bermakna — **bukan** data. Ia bagian dari arti
bendanya dan ditegakkan oleh model. Contoh: sebuah layer tidak boleh punya limit tanpa mata uang;
periode berakhir tidak boleh mendahului periode mulai; potongan tidak boleh melebihi bruto yang
menjadi dasarnya; persentase penyebaran harus berjumlah sama dengan bagian NuRe.

Perbedaan keduanya:

| | Invarian | Kelengkapan transisi |
|---|---|---|
| Berlaku | kapan saja data disimpan | hanya pada perpindahan keadaan tertentu |
| Melanggarnya berarti | catatan tidak berarti apa-apa | catatan sah tapi belum siap maju |
| Disimpan sebagai | model | data |

##### Konsekuensi

- Daftar kelengkapan per transisi harus bisa diubah orang bisnis tanpa rilis.
- Penugasan peran bertanggal (ADR-D-TI-0044) adalah data, dan harus menyimpan **sejak kapan** sebuah
  peran melekat pada seseorang — bukan hanya keadaan hari ini.

### ADR-D-TI-0039 - Perluasan ADR-D-CNP-0007: tingkat pencatatan ikut tersimpan

**Status:** diterima, 23 September 2026
**Sifat:** ini **perluasan ADR-D-CNP-0007**, bukan keputusan yang berdiri sendiri. Ia menambah satu
dimensi pada paket yang sama.

##### Konteks

ADR-D-CNP-0007 menetapkan nilai uang disimpan sebagai paket: nilai, mata uang, kurs, tanggal atau sumber
kurs, dan nilai IDR.

Rumus pencapaian di sistem lama membuka dimensi kedua yang belum pernah dinyatakan:

```
AchievementPct = ( TotalAchPremium / (RNMShareP/100) ) / EPI x 100
```

Premi dibagi bagian NuRe untuk **dinaikkan ke tingkat 100%**, lalu dibandingkan dengan EPI. Itu
bukan kebetulan aritmetika — itu pernyataan model: EPI disimpan pada tingkat 100% treaty,
sementara premi disimpan pada bagian NuRe. Dua besaran yang dibandingkan satu sama lain, disimpan
pada tingkat berbeda, didamaikan oleh pembagian yang tidak tertulis di mana pun.

##### Keputusan

> Setiap nilai uang membawa **tingkatnya** — 100% treaty atau bagian NuRe.
> Bila ia besaran bagian, ia juga membawa **bagian yang dipakai** untuk menghasilkannya.

**Tingkat adalah sifat besarannya, bukan konvensi yang kita pilih.** Limit treaty *adalah* besaran
100%: ia disepakati dengan cedant tanpa peduli berapa bagian yang diambil NuRe, dan tetap angka itu
meski bagian NuRe berubah. Premi yang diterima *adalah* besaran bagian NuRe: itu yang masuk
rekening, dan tidak ada tingkat lain yang pernah terjadi.

Menyeragamkan keduanya bukan penyederhanaan — ia menghilangkan informasi pada satu arah dan
mengarangnya pada arah lain.

##### Konsekuensi

- Penaikan dan penurunan tingkat berhenti menjadi pembagian tak tertulis yang tersebar di rumus,
  dan menjadi **konversi yang tercatat**, sejajar dengan konversi mata uang.
- Konversi tingkat memakai bagian **yang berlaku pada saat jumlah itu terjadi**, bukan bagian hari
  ini. Karena setiap jumlah membawa bagiannya sendiri, aturan ini terpenuhi dengan sendirinya.
- Besaran tingkat 100% dan besaran tingkat bagian **tidak berbagi kolom**.
- Rumus pencapaian tidak lagi memerlukan pembagian yang tidak dijelaskan, karena kedua pembandingnya
  sudah menyatakan tingkatnya sendiri.

##### Tingkat EPI — diputuskan tanpa verifikasi

**EPI dicatat pada tingkat 100% treaty.**

**Dasar.** Rumus pencapaian membagi premi dengan bagian NuRe *sebelum* membandingkannya dengan EPI;
pembagian itu hanya masuk akal bila EPI berada di tingkat 100%. Dan bila EPI ternyata di tingkat
bagian, seluruh angka pencapaian selama ini berlebih sebesar satu per bagian — untuk bagian 20%,
lima kali lipat. Selisih sebesar itu terhadap target tahunan tidak mungkin luput bertahun-tahun.

**Syarat pembalikan.** Underwriting menyatakan sebaliknya. Bila itu terjadi, ini **bukan perubahan
model** melainkan **penyajian ulang angka historis**, dan langsung naik ke daftar eskalasi.

### ADR-D-TI-0040 - Identitas kontrak, addendum, dan kunci alami yang memperingatkan

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In

##### Konteks

Sistem lama merangkai kontrak dengan satu kolom `OLDID` yang dipakai untuk dua hubungan sekaligus:
salinan sebuah kontrak, dan addendum atas sebuah kontrak. Addendum tinggal di tabel terpisah
(`M_TREATY_IN_EDM`, `TREATY_IN_EDM`) dan muncul sebagai baris setara di daftar.

Tidak ada deteksi duplikat yang berjalan: seluruh logikanya di `CheckDuplicateOffer` dimatikan.

##### Keputusan

1. **`OLDID` pecah menjadi dua hubungan yang dinamai berbeda.** "Salinan dari" dan "addendum atas"
   adalah dua hal berbeda dan tidak boleh berbagi kolom. "Addendum atas" berhenti menjadi kolom
   dan menjadi hubungan versi.
2. **Kunci alami** sebuah kontrak adalah cedant + source of business + periode + `ProportionType`.
   Ia **memperingatkan, tidak melarang** — ada keadaan sah di mana dua kontrak berbagi kunci, dan
   sistem tidak boleh memutuskan itu untuk penggunanya.
3. **Lapisan beku addendum.** Untuk setiap addendum, lima hal tidak boleh berubah: cedant, source
   of business, `ProportionType`, tanggal mulai, dan tanggal berakhir. Perubahan atas salah satunya
   berarti kontrak lain, bukan addendum.
4. **Perpanjangan di tengah periode bukan addendum.** Ia operasi tersendiri dengan namanya sendiri.
5. **Materialitas addendum diturunkan**, tidak diketik: bila perubahan itu menghasilkan baris
   selisih, ia material. Kecuali untuk baris warisan — lihat ADR-D-TI-0042.
6. **Bagian NuRe boleh berubah** lewat addendum, dan perubahannya berlaku ke depan.

##### Konsekuensi

- Deteksi duplikat adalah **kemampuan baru**, bukan pelestarian. Jangan dihitung gratis.
- Peringatan kunci alami harus bisa diabaikan oleh orang, dan pengabaiannya tercatat.
- Riwayat versi kontrak menjadi struktur pertama-kelas, bukan tabel terpisah yang dirangkai kolom.

##### Penandaan 23 September 2026 — keputusan (1) adalah PERUBAHAN, bukan pelestarian

Dibuktikan saat merancang sambungan ke modul Adjustment, dan ditandai di sini supaya tidak
ditemukan ulang sebagai "spesifikasi yang keliru".

**Di sistem lama, addendum yang disetujui tidak menggantikan apa pun.** Aktivitas
`SaveTreatyIn_EDM_Act` memuat dua langkah yang akan mempromosikannya — satu menulis ke
`M_TREATY_IN`, satu memanggil penyimpan `TREATYINDETAIL` — dan **keduanya ber-blok `//`, mati**.
Yang hidup hanya penyimpan `TREATYINDETAILEDM`. Addendum karena itu tersimpan di tabelnya sendiri,
dan baris kontrak tetap seperti semula.

**Di sistem baru, versi yang disetujui menggantikan pendahulunya**, dan hilir membaca versi yang
berlaku. Tidak ada baris "kontrak saat ini" yang terpisah untuk ditulis, sehingga tidak ada langkah
penyalin yang dapat dimatikan.

Keputusan ini **tidak bersandar pada bukti sistem lama** — bukti sistem lama justru melawannya.
Ia keputusan rancangan, diambil sadar. Seberapa banyak data yang terdampak diukur **Uji Z**.

### ADR-D-TI-0041 - Satu fakta, satu penulis, selalu

**Status:** diterima, 23 September 2026
**Berlaku untuk:** masa berdampingan Pega dan sistem baru

##### Konteks

Gelombang rilis ditetapkan sebagai **urutan merilis**, bukan urutan membangun. Maka ada masa di
mana sebagian fungsi sudah dilayani sistem baru sementara sisanya masih dilayani Pega, atas kontrak
yang sama.

Pilihan "kedua sistem sama-sama menulis, dibagi per kontrak" ditolak: Pega menulis dokumen JSON
sebagai bentuk kanoniknya, sistem baru menulis skema relasional sebagai bentuk kanoniknya. Dua
bentuk kanonik yang sama-sama ditulis melanggar ADR-D-CNP-0023. Itu bukan pilihan berisiko — ia pilihan
yang melanggar keputusan yang sudah mengikat.

##### Keputusan

> Pada setiap saat, **tepat satu** sistem berwenang menulis sebuah fakta. Sistem yang lain boleh
> **membacanya**, tidak pernah menulisnya.

Penegakannya **bukan kebijakan, melainkan izin**: hak tulis Pega dicabut di basis data atas objek
yang fungsinya sudah pindah.

Urutan operasional tiap gelombang:

1. cabut hak tulis Pega atas objek yang akan pindah;
2. **verifikasi tidak ada penulisan lagi** — aturan Pega yang gagal di sini justru menemukan
   penulis yang tidak kita ketahui. Langkah ini adalah **alat penemuan**, bukan sekadar pemeriksaan;
3. pindahkan datanya;
4. verifikasi hasil pemindahan (ADR-D-TI-0043);
5. buka hak tulis sistem baru.

##### Konsekuensi

- Masalah ketiadaan kolom waktu **lenyap**: tidak pernah ada dua catatan atas fakta yang sama untuk
  diperbandingkan umurnya.
- Tabrakan urutan ID dan nomor offer tidak mungkin terjadi, karena hanya satu sistem yang
  menerbitkan.
- Batas gelombang menjadi **batas kewenangan tulis** yang bisa diperiksa, bukan dipercaya.
- Sistem baru harus menyediakan **bentuk baca yang cocok dengan Pega** untuk fungsi yang belum
  pindah — view atas bentuk kanonik baru, bukan salinan.
- View itu **harus memaparkan ID lama dalam format lama**. Pencapaian dan produksi dicocokkan lewat
  pemotongan tujuh karakter pertama nomor offer; begitu ID berganti bentuk, pencocokan itu putus.
  Karena itu setiap kontrak yang dimigrasi **menyimpan ID lamanya secara permanen**, bukan
  sementara: orang, dokumen, surat, dan sistem lain merujuk padanya.

### ADR-D-TI-0042 - Sejarah pindah apa adanya: baris warisan dan aturan sentuh-perbaiki

**Status:** diterima, 23 September 2026
**Berlaku untuk:** migrasi data Treaty In

##### Konteks

Model baru menetapkan invarian yang tidak dipenuhi sebagian data lama — antara lain "potongan tidak
boleh melebihi bruto yang menjadi dasarnya". Pilihannya: bersihkan sejarah supaya memenuhi,
tinggalkan sejarah di sistem lama, atau pindahkan apa adanya.

Membersihkan sejarah agar memenuhi invarian berarti **mengubah angka yang pernah disetujui dan
dibukukan** — persis penyakit yang seluruh sesi ini dipakai untuk mendiagnosis. Menyembuhkannya
dengan melakukannya sendiri, dalam skala besar, sekali jalan, atas seluruh sejarah, tidak bisa
dibenarkan.

Meninggalkan sejarah di sistem lama juga gugur: Treaty In Adjustment ditetapkan mengambil data lama
lewat SELECT, dan ADR-D-CNP-0016 melarang database link. Kontrak lama yang tidak ikut pindah tidak akan
bisa dijangkau sama sekali.

##### Keputusan

Seluruh kontrak dan addendum pindah **lengkap dan apa adanya, termasuk cacatnya**, dengan tiga
aturan:

1. **Penanda warisan.** Baris hasil migrasi ditandai sebagai warisan, dan penanda itu menyatakan:
   *diterima apa adanya, tidak divalidasi terhadap invarian yang berlaku sekarang*. Penandanya
   terbaca oleh siapa pun yang membaca datanya, bukan tersembunyi di tabel teknis.
2. **Invarian ditegakkan pada penulisan, bukan pada baris.** Baris warisan ditulis sekali oleh
   migrasi dan tidak bisa disunting, jadi ia tidak pernah melanggar apa pun.
3. **Aturan sentuh-perbaiki.** Begitu sebuah kontrak warisan **disentuh** — addendum, revisi,
   adjustment, apa pun yang menulis — ia harus memenuhi invarian saat itu juga, dan penanda
   warisannya dicabut. Yang tidak disentuh tetap warisan selamanya.

##### Materialitas addendum historis

Aturan penurunan materialitas yang baru **tidak dijalankan** atas baris warisan. Menjalankannya
akan **mengklasifikasi ulang sejarah**, dan itu bentuk lain dari mengubah masa lalu.

- `EDMMATERIALTYPE` diambil **apa adanya** sebagai nilai warisan yang tidak diverifikasi.
- Kolom itu membawa **asal-usulnya**: diturunkan oleh aturan, atau diimpor sebagai warisan.

Satu kolom memang akan menyimpan sesuatu yang aturannya sendiri tidak menghasilkan, dan itu benar
selama kolom itu jujur menyatakan dari mana nilainya datang. Yang salah bukan menyimpan nilai
warisan; yang salah adalah menyimpannya seolah ia hasil aturan.

##### Konsekuensi

- **Tidak ada anggaran pembersihan di muka.** Biayanya tersebar dan hanya dibayar untuk kontrak
  yang memang masih dipakai.
- Yang diperbaiki adalah yang **diperiksa orang**, bukan yang ditebak skrip.
- Ada kemungkinan **tiga generasi data**, bukan dua: sebagian data sekarang sendiri hasil migrasi
  sebelumnya dari tabel yang lebih lama. Penanda warisan perlu menyatakan generasi mana.

### ADR-D-TI-0043 - Migrasi memindahkan; menghitung ulang adalah peristiwa bisnis tersendiri

**Status:** diterima, 23 September 2026
**Berlaku untuk:** migrasi data dan penerimaan kalkulator baru

##### Konteks

Godaan terbesar sebuah migrasi adalah menjalankan kalkulator baru sambil memindahkan data, karena
terasa efisien. Akibatnya kedua peristiwa tercampur dan tidak ada yang bisa menjawab apakah
perbedaan angka berarti migrasi gagal atau perbaikan berhasil.

##### Keputusan

Keduanya adalah peristiwa terpisah, dan kriteria keberhasilannya **berlawanan**:

| | Kriteria |
|---|---|
| **Migrasi** — memindahkan angka apa adanya | **Setiap angka identik.** Bukan mirip, bukan dalam toleransi. Selisih satu rupiah berarti migrasi gagal dan diulang. |
| **Menghitung ulang** dengan kalkulator baru | Hasilnya **berbeda**, dan perbedaannya **bisa dijelaskan** oleh cacat yang sudah dikenali. Perbedaan yang tidak bisa dijelaskan berarti kalkulator baru yang salah. |

Migrasi menulis baris warisan apa adanya. **Sentuhan pertama** (ADR-D-TI-0042) yang menghitung ulang,
dan hitung ulang itu adalah peristiwa yang dibukukan, bukan hasil migrasi.

##### Uji paritas

Kalkulator baru dijalankan atas seluruh data yang sudah dimigrasi, dan hasilnya **dikeluarkan
sebagai laporan selisih**. Tidak dituliskan ke mana pun.

Satu run itu punya dua kegunaan sekaligus:

1. **Uji terima kalkulator baru** — setiap selisih harus jatuh ke salah satu kelas cacat yang sudah
   dikenali. Selisih yang tidak jatuh ke kelas mana pun adalah cacat **baru** di kalkulator baru.
2. **Pengukuran kesalahan historis** — jumlahnya dalam rupiah, per kelas cacat, per tahun treaty.
   Itu angka yang dibutuhkan daftar eskalasi, dan datang gratis dari run yang memang harus jalan.

**Syarat yang tidak boleh dilewati:** tuliskan lebih dulu, per cacat yang sudah dikenali, **bentuk
selisih apa** yang seharusnya ia hasilkan — arahnya, besarannya, dan populasi terdampak. Tulis itu
**sebelum** run dijalankan. Tanpa itu, pembacanya akan merasionalisasi apa pun yang muncul, dan uji
terima berubah jadi latihan menjelaskan. Daftar harapan itu adalah keluaran tersendiri.

##### Ukuran keberhasilan migrasi

1. Jumlah kontrak dan addendum cocok.
2. Setiap ID lama punya pasangan di sistem baru, **dan** setiap ID baru punya asal. Migrasi yang
   menciptakan baris dari ketiadaan sama buruknya dengan yang menghilangkan baris.
3. Nilai kunci dipindahkan **tanpa dihitung ulang** lalu dibandingkan — kriteria identik di atas.
4. **Jumlah baris anak per kontrak cocok**: layer, detail, penyebaran, installment, periode
   pelaporan. Kontrak yang headernya pindah utuh tetapi kehilangan satu baris layer akan lolos
   ketiga ukuran pertama, dan itu kegagalan yang paling sering terjadi pada migrasi dari dokumen
   JSON ke tabel relasional.

### ADR-D-TI-0044 - Wewenang: keadaan dan peran, terpisah tegas

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In

##### Konteks

Di seluruh ekspor Pega modul ini, **tidak ada satu pun nama hak akses yang terisi** — mekanisme
peran bawaan Pega hadir 2.202 kali sebagai tempat kosong dan tidak pernah dipakai sekali pun.

Yang benar-benar mengatur siapa boleh apa hanya dua, dan keduanya bukan peran: sebuah bendera pada
kontraknya yang diperiksa sekitar 660 kali dengan delapan ejaan berbeda, dan nama orang yang
ditanam langsung di dalam aturan persetujuan.

Jadi yang kita lihat bukan keputusan bahwa wewenang milik kontrak. Kita melihat sistem yang tidak
pernah punya cara menyatakan "siapa", sehingga pertanyaan itu tidak pernah bisa diajukan.

##### Keputusan

> **Keadaan** menjawab: apakah tindakan ini mungkin sekarang.
> **Peran** menjawab: apakah orang ini boleh melakukannya.
> Setiap larangan harus punya **tepat satu** sebab.

**Uji rancangannya konkret:** setiap penolakan harus bisa menyebutkan satu alasan saja. Bila sebuah
penolakan harus berbunyi "karena kontraknya sudah disetujui **atau** karena Anda bukan direktur",
rancangannya gagal — dua mekanisme sudah saling meniru. Ini uji rancangan, bukan pedoman menulis
pesan.

##### Konsekuensi

- **Entitas peran dan penugasan bertanggal masuk gelombang 1.** Penugasan menyimpan sejak kapan
  sebuah peran melekat pada seseorang, bukan hanya keadaan hari ini — karena jejak perubahan harus
  mencatat peran yang berlaku **saat itu**.
- Cakupannya bukan hanya persetujuan, melainkan seluruh tindakan.
- Nama orang tidak pernah muncul di dalam aturan.

### ADR-D-TI-0045 - Jejak perubahan adalah fakta mesin, terpisah dari catatan manusia

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In, gelombang 1

##### Konteks

Satu-satunya jejak di sistem lama adalah daftar komentar, berisi tanggal, pengenal operator, status
akseptasi, dan teks bebas. Ia ditulis **hanya pada perpindahan persetujuan**. Penyimpanan biasa —
termasuk yang mengubah angka — tidak meninggalkan jejak siapa maupun kapan. Tidak ada kolom waktu
sama sekali di tabel kontrak maupun tabel addendum.

##### Keputusan

Jejak perubahan masuk **gelombang 1**, dan alasannya bukan kepatuhan:

> **Jejak perubahan adalah bukti bahwa pembekuan benar-benar terjadi.**

ADR-D-TI-0036 menetapkan angka yang disetujui dibekukan. Pemilik proses menetapkan Adjustment mengambil
data lama lewat SELECT, bukan menyalinnya. Keduanya hanya bisa dipercaya bila ada cara membuktikan
data lama tidak berubah di antara dua pembacaan. Tanpa jejak, pembekuan adalah **janji**, bukan
fakta — dan selisih yang dibukukan Adjustment berdiri di atas janji.

**Isi minimal:** siapa, kapan, apa yang berubah dari nilai apa ke nilai apa, dan **di bawah peran
apa** — peran yang berlaku saat itu, bukan peran orangnya hari ini. Butir terakhir menyambung ke
penugasan peran bertanggal (ADR-D-TI-0044), dan itu sebabnya keduanya harus lahir bersamaan.

##### Jangan digabung dengan catatan komentar

Keduanya benda berbeda:

| | Catatan komentar | Jejak perubahan |
|---|---|---|
| Isinya | narasi manusia: alasan, pertimbangan | fakta mesin |
| Kapan | pada perpindahan persetujuan | pada setiap penulisan |
| Ditulis oleh | orang | sistem |
| Bisa disunting | ya | **tidak** |
| Teks bebas | ya | tidak ada |

Menggabungkannya menghasilkan yang terburuk dari keduanya: catatan yang bisa disunting dan karena
itu tidak membuktikan apa-apa.

### ADR-D-TI-0046 - Satu keadaan siklus hidup kontrak

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In

##### Konteks

Lima properti tumpang tindih menjawab pertanyaan yang sama — apakah kontrak ini boleh diubah, dan
di mana ia berada dalam alurnya:

| Properti | Rujukan di ekspor |
|---|---|
| `Position` | 5.634 |
| `ViewState` | 661 |
| `StatusAkseptasi` | 165 |
| `IsEditData` | 125 |
| `RevisionState` | 21 |

`ViewState` disetel **imperatif** di lima tempat, termasuk satu aturan yang namanya sendiri sudah
mengakui bendanya: `TreatyInForceEdit`. Padahal isinya sepenuhnya bisa diturunkan dari posisi
persetujuan dan mode revisi kontrak itu. Ia turunan yang disimpan dan disetel oleh klik.

Perbandingannya ditulis dengan **delapan ejaan berbeda** untuk maksud yang sama, sehingga satu
salah ketik membuka atau mengunci sebuah field tanpa ada yang tahu.

##### Keputusan

Kelima bendera dilebur menjadi **satu keadaan siklus hidup** kontrak: satu nilai bertipe tegas
dengan himpunan nilai tertutup.

**"Boleh diubah" adalah turunan** dari keadaan itu, dihitung saat ditanya, dan tidak pernah
disimpan sebagai bendera tersendiri.

`ViewState` **tidak dinamai ulang — ia tidak ada di model baru.** Menamainya ulang hanya akan
melestarikan bendanya dan memperbaiki labelnya.

##### Konsekuensi

- Delapan ejaan perbandingan tidak direkonsiliasi dan tidak dicari mana yang benar — tidak ada yang
  diangkut.
- Salah ketik tidak mungkin lolos, karena himpunan nilainya tertutup.
- Ini instans ADR-D-TI-0037: turunan tidak disimpan dan tidak disetel oleh klik.

### ADR-D-TI-0047 - Data uji: identitas boleh disamarkan, angka tidak

**Status:** diterima, 23 September 2026
**Berlaku untuk:** migrasi dan pengujian Treaty In

##### Keputusan pemisah

> **Identitas boleh disamarkan. Angka tidak boleh disentuh sama sekali.**

| Boleh disamarkan | Tidak boleh diubah walau satu digit |
|---|---|
| nama cedant, leading reinsurer, broker, retrosesioner | seluruh nilai uang |
| nama operator pada catatan komentar | seluruh persentase |
| nama berkas lampiran | seluruh tanggal |
| | mata uang, jumlah lines, jumlah reinstatement, struktur layer |

Alasannya dari ADR-D-TI-0043: kriteria migrasi adalah **identik**, nol toleransi, dan uji paritas
membandingkan hasil kalkulator baru terhadap angka lama. Penyamaran yang menyentuh angka
menghancurkan kedua uji itu sekaligus, dan tidak ada cara mengetahuinya sampai terlambat.

**Teks bebas pada catatan komentar dibuang seluruhnya, bukan disamarkan.** Ia berisiko tertinggi —
orang menulis apa saja di kolom komentar — dan paling tidak diperlukan, karena tidak satu pun uji
paritas membacanya.

##### Kejujuran yang harus dicatat

**Menyamarkan nama tidak menganonimkan data treaty.** Sebuah treaty dengan limit tertentu, periode
tertentu, susunan layer tertentu, dan leading reinsurer tertentu dapat dikenali kembali oleh siapa
pun yang bekerja di pasar reasuransi Indonesia. Data ini kecil dan pesertanya saling kenal.

Jangan membangun rasa aman dari transformasi. Perlindungan yang nyata ada pada **kendali akses dan
lingkungan**:

- data uji tinggal di lingkungan terkendali dengan akses bernama, tidak disalin ke laptop;
- siapa yang mengaksesnya tercatat;
- salinan penuh dipakai untuk run paritas lalu dibersihkan, bukan dibiarkan menganggur.

##### Dua kumpulan, bukan satu

| | Kumpulan kurasi | Salinan paritas |
|---|---|---|
| Ukuran | kecil, satu kontrak per kelas cacat | lengkap |
| Angka | boleh disamarkan berat, sebagian boleh dibuat | **asli**, tidak disentuh |
| Teks bebas | boleh dibuat | dibuang |
| Dipakai untuk | membangun dan menguji sehari-hari | run paritas saja |
| Letaknya | boleh di laptop orang | lingkungan terkendali |

Dengan pemisahan itu, pekerjaan harian tidak pernah menyentuh salinan penuh, dan salinan penuh
tidak pernah kehilangan presisinya.

##### Kumpulan kurasi adalah KELUARAN uji, bukan disusun tangan

Setiap uji A sampai V **sudah menghasilkan daftar ID kontrak** yang memenuhi kriterianya — itu
justru inti kerjanya. Kumpulan kurasi adalah keluaran uji-uji itu: ambil satu atau dua ID dari
setiap daftar.

Kelas yang harus terwakili: layer bermata uang ganda; kontrak berbagian per baris; addendum
berprorata; addendum yang mengubah bagian NuRe; penyebaran manual; penyebaran otomatis yang
masternya sudah berubah; kontrak bertanda kegagalan Rate on Line; kontrak FAC-OUT dengan sisa
daftar penyebaran; kontrak berpenanda konversi generasi lama; addendum ber-`EDMSTATE` 3; baris
reinstatement yang persennya bukan 100; dan kontrak dengan rasio potongan terhadap bruto yang
menyimpang.

**Konsekuensi jadwal, direvisi 23 September 2026.** Kumpulan kurasi tetap merupakan **keluaran**,
bukan susunan tangan — tetapi bukan keluaran uji A–V, melainkan keluaran **run paritas**. Run
paritas sendiri **tidak menunggu apa pun**: ia berjalan atas seluruh data yang dimigrasi, dan
justru run itulah yang **menamai** kontrak mana yang masuk kumpulan kurasi.

Urutannya: **migrasi → run paritas → kumpulan kurasi.** Bukan sebaliknya, dan tidak ada yang
tertahan di dalamnya.

Yang tetap harus ditulis **sebelum** run dijalankan adalah **bentuk** selisih yang diharapkan per
kelas cacat — arah, rumus penyimpangan, tanda pengenal. Hanya **populasinya** yang datang dari run.
Lihat `KEPUTUSAN-TANPA-VERIFIKASI.md` §5.

### ADR-D-TI-0048 - Sambungan (seam) ke Treaty In Adjustment

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In
**Batas:** modul Treaty In Adjustment sendiri **berada di bawah embargo** dan tidak dibedah. ADR
ini hanya menetapkan apa yang harus **disediakan** Treaty In.

##### Arahan pemilik proses yang mengikat

1. Saat sebuah Treaty In Adjustment dibuat, ID lama yang disimpan padanya adalah **ID Treaty In
   yang disesuaikan**. Sambungannya eksplisit, bukan disimpulkan.
2. **Data lama tidak boleh disimpan ulang.** Ia harus diambil lewat SELECT dari data lama. Tidak
   ada penyalinan nilai kontrak ke dalam catatan Adjustment.
3. **Selisih = nilai sekarang − nilai lama**, dan ia disimpan di **tabel tersendiri**.

Aturan (2) adalah ADR-D-CNP-0023 yang diterapkan oleh pemilik prosesnya sendiri: satu bentuk kanonik,
hilir diberi rujukan bukan salinan.

##### Yang harus disediakan Treaty In, dan hanya ini

1. **Identitas versi kontrak** yang stabil, tidak berubah, dan bisa ditunjuk dari luar konteks.
2. **Jaminan** bahwa versi yang sudah disetujui tidak pernah berubah nilainya (ADR-D-TI-0036, dibuktikan
   oleh ADR-D-TI-0045).
3. **Bentuk baca** yang bisa dipakai konteks lain untuk mengambil nilai kontrak pada versi
   tertentu — lengkap dengan tingkat, mata uang, kurs, dan bagian yang menyertainya (ADR-D-TI-0039).

Tabel selisih, layarnya, dan alur kerjanya **bukan urusan sesi ini dan tidak dirancang**.

##### Sifat yang diinginkan, dicatat sekarang

Selisih adalah **turunan yang dibukukan**, jadi ia fakta tercatat — bentuk yang sama dengan potret
penyebaran. Karena itu baris selisih membawa **penunjuk ke versi kontrak mana** ia dihitung dan
**kapan**. Tanpa itu ia angka tanpa asal.

Dan karena nilai lama bisa diambil lewat SELECT sementara nilai sekarang diketahui, **selisih yang
tersimpan dapat direkonsiliasi ulang kapan saja**. Ketidakcocokan antara selisih tersimpan dan
selisih terhitung bukan masalah — **ia alat deteksi**. Rancang supaya rekonsiliasi itu mungkin.

##### Kenapa pembekuan menjadi prasyarat, bukan kebersihan

Aturan (2) hanya aman bila data lama **tidak bergerak**. Selama nilai kontrak masih bisa berubah
sesudah akseptasi — lewat penyusunan ulang penyebaran, suntingan baris acuan, kurs tahun berjalan,
atau bagian NuRe yang berubah — maka selisih yang dihitung hari ini dan yang dihitung bulan depan
atas kontrak yang sama akan berbeda tanpa ada yang mengubah apa pun, dan selisih yang dibukukan
tidak bisa dipertanggungjawabkan.

ADR-D-TI-0036 karena itu bukan pilihan rancangan yang bisa ditunda. Ia **prasyarat teknis** bagi cara
kerja Adjustment yang sudah ditetapkan.

### ADR-D-TI-0049 - Jenis addendum adalah satu sumbu; materialitas diturunkan

**Status:** diterima **tanpa verifikasi**, 23 September 2026
**Berlaku untuk:** Treaty In

##### Konteks

Sistem lama menyimpan dua penanda pada setiap addendum: sebuah penanda keadaan dengan tiga nilai,
dan sebuah penanda materialitas. Keduanya di-AND di dalam kondisi aturan, yang membuat *kode*
memperlakukannya sebagai dua besaran independen.

Perdebatan selama penggalian: apakah dua kolom itu dua gagasan bisnis, atau satu gagasan yang
tumbuh dua kali. Uji data yang dirancang untuk memutuskannya (Uji A, B, C) **tidak dijalankan**.

##### Keputusan

**Satu sumbu.** Jenis addendum punya tiga nilai:

| Jenis | Materialitas turunannya |
|---|---|
| perubahan estimasi | material |
| penyesuaian ke nilai aktual | material |
| administratif | non-material |

**Kolom materialitas tidak ada** di model baru. Ia diturunkan dari jenis.

##### Dasar — penalaran, bukan verifikasi

Kombinasi *(administratif, material)* tidak punya arti bisnis apa pun, dan tidak ada apa pun di
sistem lama yang mencegahnya terbentuk. Dua kolom yang dapat saling bertentangan, dengan sedikitnya
satu kombinasi yang tak bermakna, adalah tanda **satu gagasan yang tumbuh dua kali** — bukan dua
sumbu yang saling melengkapi.

Bahwa kode meng-AND keduanya membuktikan kode memperlakukannya independen. Itu **bukan** bukti
bisnis punya dua gagasan.

##### Data warisan

Nilai materialitas lama **disimpan sebagai nilai warisan** dengan asal-usulnya (ADR-D-TI-0042). Ia tidak
dipakai menurunkan apa pun dan tidak menjadi kolom hidup. Aturan penurunan yang baru **tidak
dijalankan** atas baris warisan — menjalankannya akan mengklasifikasi ulang sejarah.

##### Syarat pembalikan

Bisnis menyebut **satu kasus nyata** di mana addendum administratif harus dinyatakan material, atau
menyebut pembeda material yang bukan tentang akibat premi.

**Bila terbalik:** kolom materialitas kembali sebagai sumbu kedua yang harus dijaga konsisten
selamanya. Ini salah satu dari enam keputusan yang **mengubah bentuk model**.

### ADR-D-TI-0050 - Acuan kapasitas adalah sumber luar yang tidak dipercaya

**Status:** diterima **tanpa verifikasi**, 23 September 2026
**Berlaku untuk:** Treaty In

##### Konteks

Pembagian kapasitas NuRe diambil dari sebuah tabel acuan. Uji D dan E — siapa pemeliharanya, dan
apakah periode susunan bisa tumpang tindih — **tidak dijalankan**.

##### Keputusan

Treaty In **tidak memiliki** tabel itu. Ia dibaca lewat **adapter**, dan diperlakukan sebagai
sumber luar yang tidak dipercaya:

1. Bertipe tegas di perbatasan; nilai yang tidak terurai menyebabkan **kegagalan keras**, bukan
   nilai pengganti (ADR-D-TI-0035).
2. Bila **lebih dari satu susunan cocok** untuk satu tanggal, **adapter gagal**. Tidak memilih yang
   pertama, tidak memilih yang terbaru.

> Ambiguitas adalah keadaan yang **dilaporkan**, bukan yang diselesaikan diam-diam.

##### Dasar — sebagian artefak, sebagian penalaran

**Dari artefak:** tidak ada satu pun aturan di seluruh ekspor Pega yang menulis ke tabel itu —
Treaty In murni pembaca. Bentuk kolomnya adalah bentuk kontrak berperiode, bukan lookup sederhana.
Seluruh kolomnya bertipe teks bebas, tanpa primary key dan tanpa indeks.

**Dari penalaran:** ketiadaan kunci dan indeks berarti tumpang tindih periode **memang mungkin**
terjadi, dan tidak ada apa pun yang mencegahnya. Adapter yang memilih diam-diam akan menghasilkan
kapasitas yang berbeda-beda tanpa ada yang tahu mana yang benar.

##### Syarat pembalikan

Pemelihara tabel ditemukan **dan** menyanggupi constraint yang menjamin ketunggalan. Saat itu
adapter boleh dilonggarkan menjadi pemilihan yang dinyatakan.

**Bila terbalik:** aturan pemilihan susunan masuk ke model. Ini salah satu dari enam keputusan yang
mengubah bentuk model.

### ADR-D-TI-0051 - Anggap ada setidaknya satu konsumen hilir yang tidak dikenal

**Status:** diterima **tanpa verifikasi**, 23 September 2026
**Berlaku untuk:** Treaty In

##### Konteks

Tidak ada satu pun pembaca tabel terbitan di dalam ekspor Pega — konsumennya, bila ada, berada di
luar. Uji F dan G, yang akan menyebutkan siapa mereka, **tidak dijalankan**.

Pencocokan antar-sistem yang terlihat di sistem lama memakai **tujuh karakter pertama** nomor
offer, yang terikat pada format ID `kode_situs || nomor urut enam digit`.

##### Keputusan

1. **Anggap ada setidaknya satu konsumen yang tidak dikenal.**
2. Bentuk terbitan hilir dipertahankan sebagai **turunan** dari bentuk kanonik baru — view, bukan
   salinan (ADR-D-CNP-0023).
3. View itu **memaparkan ID lama dalam format lama**, supaya pemotong tujuh karakter tetap bekerja.
4. Setiap kontrak yang dimigrasi **menyimpan ID lamanya secara permanen**, bukan sementara. Orang,
   surat, dan dokumen merujuk padanya.

##### Dasar — ongkos yang tidak setangkup

Ini bukan taruhan atas kemungkinan, melainkan atas **ongkos salahnya ke dua arah**:

| Bila salah | Akibatnya |
|---|---|
| Menganggap ada konsumen, padahal tidak | memelihara sebuah view yang tidak dibaca siapa pun — **murah** |
| Menganggap tidak ada, padahal ada | integrasi putus, dan baru ketahuan **di produksi** — **mahal** |

Bila ongkosnya tidak setangkup, ambil sisi yang murah tanpa menunggu bukti.

##### Syarat pembalikan

Sensus konsumen menunjukkan tidak ada pembaca sama sekali — view boleh dipensiunkan.
**Pemensiunan mudah, pemulihan tidak**, dan itu sebabnya urutannya begini dan bukan sebaliknya.

### ADR-D-TI-0052 - Satu rantai persetujuan empat tingkat, tabel perutean kosong dari pengecualian

**Status:** diterima **tanpa verifikasi**, 23 September 2026
**Berlaku untuk:** Treaty In

##### Konteks

Sistem lama punya dua penyimpangan dari rantai empat tingkat: sebuah jalur alternatif lewat
penerima tugas kelompok, dan sebuah jalur revisi yang memendekkan empat tingkat menjadi dua.
Pemilihan jalur revisi dilakukan lewat tombol, tanpa aturan apa pun yang mengaturnya.

Uji H, yang akan merekonstruksi jalur yang benar-benar ditempuh, **tidak dijalankan**.

##### Keputusan

1. **Satu rantai bawaan empat tingkat.**
2. Tabel peruteannya **kosong dari pengecualian**.
3. Jalur alternatif kelompok **tidak dibawa**.
4. Jalur revisi **tidak dipendekkan**. Berlaku aturan: *jalur untuk suatu perubahan tidak boleh
   lebih pendek daripada jalur yang diperlukan untuk nilai hasilnya.*

##### Dasar — arah pembatalan yang tidak setangkup

Tidak ada bukti jalur alternatif itu pernah disahkan sebagai kebijakan, dan pemendekan jalur revisi
dipilih lewat tombol oleh orang yang sama yang mengisi kontraknya.

| Pilihan | Bila ternyata salah |
|---|---|
| **Membuang keduanya** | dibatalkan dengan **menambah satu baris tabel** |
| Mempertahankan keduanya | tidak dapat dibatalkan tanpa **mengaudit ulang seluruh kontrak yang sudah melewatinya** |

##### Syarat pembalikan

Satu baris di tabel perutean — karena aturan perutean berbentuk data, bukan kode (ADR-D-TI-0038).

Ini **bukan** salah satu keputusan yang mengubah bentuk model.

### ADR-D-TI-0053 - Mata uang sebagai daftar; `CurrencyRelation` disimpan tetapi tidak menggerakkan angka

**Status:** diterima **tanpa verifikasi**, 23 September 2026
**Berlaku untuk:** Treaty In

##### Konteks

Sistem lama menyimpan mata uang layer dalam **dua slot** (`Currency`/`Currency2`,
`Limit`/`Limit2`), dan sebuah atribut `CurrencyRelation` bernilai `AND` atau `OR` yang membuat
rumus limit bercabang. Tidak ada satu pun aturan, memo, maupun label yang menjelaskan apa arti
bisnis kedua nilai itu.

Uji I dan J, yang akan menunjukkan apakah atribut itu masih berarti dan apakah asumsi slot keliru,
**tidak dijalankan**.

##### Keputusan

###### a. Bentuk: daftar, bukan dua slot

Mata uang per layer **tidak dibatasi**. Bila batas dua kelak dikonfirmasi bisnis, ia menjadi
**aturan validasi**, bukan bentuk tabel.

###### b. `CurrencyRelation`: disimpan, ditampilkan, **tidak dipakai menghitung**

Ia bertahan sebagai atribut bernilai bernama dan tetap terlihat pengguna, tetapi **tidak dipakai
dalam perhitungan apa pun** sampai bisnis menamai artinya.

> Atribut yang artinya tidak diketahui **tidak boleh menggerakkan angka.** Menyimpannya aman;
> membiarkannya bercabang di rumus adalah mewarisi perilaku yang tidak bisa dijelaskan siapa pun.

###### c. Penyebut Rate on Line dihitung sekali

Setiap mata uang dikonversi sekali lalu dijumlahkan. Tidak ada percabangan atas `CurrencyRelation`,
dan tidak ada penjumlahan limit per mata uang EGNPI.

##### Dasar

**Bagian a — penalaran.** Asimetri di cabang lama (satu sisi menambahkan limit tanpa konversi, sisi
lain dengan konversi) lebih mungkin merupakan jejak **asumsi bahwa slot pertama adalah IDR**
daripada semantik AND/OR. Bentuk daftar menampung kedua kemungkinan tanpa kehilangan apa pun, dan
tidak memaksa kita memilih sebelum tahu.

**Bagian b — penalaran, dan ini yang terpenting.** Membawa percabangan yang tidak bisa dijelaskan
berarti memindahkan perilaku yang tidak ada pemiliknya ke sistem baru, di mana ia akan tampak
disengaja.

**Bagian c — artefak.** Penjumlahan limit di sistem lama berada **di dalam** perulangan daftar
EGNPI, sehingga limit terjumlah sekali untuk setiap mata uang EGNPI; dan pada relasi `AND`, dua
langkah berurutan menambahkan limit yang sama dua kali dalam satu iterasi. Keduanya cacat, bukan
kebijakan.

##### Syarat pembalikan

Bisnis menamai arti `AND` dan `OR`. Saat itu `CurrencyRelation` **boleh naik menjadi masukan
perhitungan**, dan rumus limit bercabang atasnya.

**Bila terbalik:** ini salah satu dari enam keputusan yang mengubah bentuk model — atribut pasif
menjadi masukan aktif.

### ADR-D-TI-0054 - Keadaan warisan yang tidak ada padanannya

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In
**Menyelesaikan:** tabrakan antara ADR-D-TI-0042 dan ADR-D-TI-0046

##### Konteks

Dua keputusan yang sudah diambil saling bertabrakan pada satu titik:

| ADR | Isinya |
|---|---|
| ADR-D-TI-0042 | data warisan dipindahkan **apa adanya**, termasuk cacatnya; sejarah tidak dibersihkan |
| ADR-D-TI-0046 | keadaan siklus hidup adalah **satu nilai bertipe tegas dengan himpunan tertutup** |

Titik tabrakannya nyata, bukan hipotetis. Aturan `TreatyInSetValue` di sistem lama menyetel
`StatusAkseptasi` ke nilai `"test"`, dan aturan itu terpasang di enam layar. Nilai `"test"` tidak
punya padanan di antara keadaan sah mana pun.

Memindahkannya apa adanya melanggar himpunan tertutup. Memetakannya ke salah satu keadaan sah
berarti membersihkan sejarah, yang ADR-D-TI-0042 larang. Keduanya tidak bisa sekaligus dipenuhi tanpa
aturan tambahan.

Ini tabrakan pertama di antara keduanya dan hampir pasti bukan yang terakhir, karena `"test"`
bukan satu-satunya nilai liar yang mungkin ada di dua puluh tahun data. Karena itu ia diselesaikan
sebagai **aturan**, bukan sebagai penanganan satu kasus.

##### Keputusan

Himpunan keadaan memuat **satu anggota tambahan yang menyatakan dirinya sendiri sebagai warisan tak
terpetakan**: `WARISAN_TAK_TERPETAKAN`.

Ia bukan "lainnya" dan bukan "tidak diketahui". Artinya tepat satu hal: *nilai warisan yang tidak
ada padanannya di antara keadaan sah, dan nilai aslinya tercatat.*

Tiga ketentuan yang menyertainya:

1. **Nilai asli disimpan.** Atribut `KEADAAN_WARISAN_ASLI` memuat teks keadaan lama apa adanya.
   Ia terisi **hanya** pada baris berkeadaan `WARISAN_TAK_TERPETAKAN`, dan tidak pernah dibaca oleh
   perhitungan mana pun — ia catatan, bukan masukan.

2. **Tidak ada perpindahan masuk selain migrasi.** Sistem berjalan tidak pernah bisa menghasilkan
   keadaan ini. Satu-satunya pintu masuknya adalah pemindahan data lama.

3. **Satu-satunya perpindahan keluar adalah `PERBAIKAN_WARISAN`.** Perpindahan ini menetapkan
   keadaan sah yang **dipilih secara eksplisit** oleh orang yang memperbaikinya — bukan ditebak
   sistem — dan tercatat di jejak perubahan (ADR-D-TI-0045) beserta siapa dan kapan. Sesudah itu baris
   tersebut berjalan seperti baris lain.

##### Konsekuensi

- Himpunan keadaan **tetap tertutup**: tidak ada nilai di luar daftar, karena yang di luar daftar
  punya satu tempat bernama.
- Sejarah **tetap tidak dibersihkan**: nilai aslinya ada, dapat dibaca, dan tidak ada yang
  menyulapnya menjadi keadaan yang tidak pernah dicapainya.
- Baris bermasalah **tidak bisa menyelinap ke alur kerja biasa**, karena ia tidak punya perpindahan
  keluar selain perbaikan yang sadar.
- Ini instans dari aturan sentuh-perbaiki pada ADR-D-TI-0042: invarian baru ditegakkan saat sebuah
  kontrak disentuh, bukan lewat pembersihan massal.

##### Berapa banyak, dan apakah itu mengubah keputusan

Jumlahnya belum diketahui; ia diukur oleh **Uji S-7** pada berkas permintaan DBA.

Hasilnya **tidak mengubah keputusan ini**. Bila nol, aturannya tetap ditulis, karena ia menutup
kelas persoalan dan bukan satu nilai. Yang berubah hanya perkiraan berapa banyak pekerjaan
perbaikan yang menunggu di masa depan.

### ADR-D-TI-0055 - Daftar keadaan siklus hidup dan perpindahan yang sah

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In
**Melaksanakan:** ADR-D-TI-0046 (satu keadaan), ADR-D-TI-0052 (satu rantai persetujuan), ADR-D-TI-0054 (warisan)

##### Konteks

ADR-D-TI-0046 menetapkan **bahwa** ada satu keadaan bertipe tegas dengan himpunan tertutup. Ia tidak
menetapkan **apa saja isinya**. ADR ini mengisinya, dan mengisinya dalam bentuk daftar — karena
keadaan yang tidak punya perpindahan masuk maupun keluar adalah keadaan yang tidak pernah terjadi,
dan itu hanya ketahuan dari daftar, tidak dari uraian.

Dua hal ditemukan justru karena daftarnya disusun lebih dulu, dan keduanya argumen dari data:

**`Position` kosong menandai dua keadaan yang berlawanan.** Di sistem lama, `Resolve Complete`
(disetujui) dan `Decline` (ditolak) sama-sama menyetel `Position` ke nilai kosong. Hanya
pasangannya dengan `StatusAkseptasi` yang membedakan kontrak yang disetujui dari kontrak yang
ditolak. Itu bukan alasan selera untuk menyatukan kelima bendera menjadi satu keadaan — itu alasan
dari data.

**Yang menggantikan jalan pintas revisi bukan aturan baru, melainkan bentuk.** Di sistem lama,
`RevisionState == 1` membuat persetujuan Sec Head langsung menghasilkan `Resolve Complete`,
memotong dua tingkat. ADR-D-TI-0052 menghapus jalan pintas itu. Yang menggantikannya bukan aturan
pengganti: perubahan atas versi yang sudah disetujui **melahirkan versi baru yang mulai dari
`DRAFT`**. Versi baru itu melewati keempat tingkat karena ia memang versi baru, bukan karena ada
aturan yang memaksanya.

##### Keputusan

###### Daftar keadaan — enam, ditambah satu keadaan warisan

| Keadaan | Padanan lama (`Position` + `StatusAkseptasi`) | Terminal |
|---|---|---|
| `DRAFT` | kosong atau `ReasTreatyInAdmin`, status kosong **atau** `Reject` | tidak |
| `MENUNGGU_SEC_HEAD` | `ReasTreatyInSecHead` + `Accept` | tidak |
| `MENUNGGU_DEPT_HEAD` | `ReasTreatyInDeptHead` + `Accept` | tidak |
| `MENUNGGU_DIREKTUR` | `ReasTreatyInDirector` + `Accept` | tidak |
| `DISETUJUI` | kosong + `Resolve Complete` | **ya** |
| `DITOLAK` | kosong + `Decline` | **ya** |
| `WARISAN_TAK_TERPETAKAN` | nilai apa pun di luar daftar di atas (ADR-D-TI-0054) | tidak |

###### Perpindahan yang sah — dua belas  → **TIGA BELAS sejak perubahan 24 September 2026**

> Tabel di bawah adalah daftar sebagaimana diputuskan 23 September 2026. Satu perpindahan
> ditambahkan kemudian: **`DRAFT` — `BATALKAN` → `DIBATALKAN`**. Lihat bagian
> "Perubahan 24 September 2026" di kaki berkas ini.

| Dari | Peristiwa | Ke |
|---|---|---|
| (versi baru lahir) | `LAHIR` | `DRAFT` |
| `DRAFT` | `AJUKAN` | `MENUNGGU_SEC_HEAD` |
| `MENUNGGU_SEC_HEAD` | `SETUJUI` | `MENUNGGU_DEPT_HEAD` |
| `MENUNGGU_SEC_HEAD` | `KEMBALIKAN` | `DRAFT` |
| `MENUNGGU_SEC_HEAD` | `TOLAK` | `DITOLAK` |
| `MENUNGGU_DEPT_HEAD` | `SETUJUI` | `MENUNGGU_DIREKTUR` |
| `MENUNGGU_DEPT_HEAD` | `KEMBALIKAN` | `DRAFT` |
| `MENUNGGU_DEPT_HEAD` | `TOLAK` | `DITOLAK` |
| `MENUNGGU_DIREKTUR` | `SETUJUI` | `DISETUJUI` |
| `MENUNGGU_DIREKTUR` | `KEMBALIKAN` | `DRAFT` |
| `MENUNGGU_DIREKTUR` | `TOLAK` | `DITOLAK` |
| `WARISAN_TAK_TERPETAKAN` | `PERBAIKAN_WARISAN` | keadaan sah yang dipilih eksplisit (ADR-D-TI-0054) |

Setiap keadaan punya sedikitnya satu perpindahan masuk dan, kecuali yang terminal, sedikitnya satu
perpindahan keluar. Tidak ada keadaan yatim.

##### Empat keputusan yang tertanam di daftar itu

###### 1. `ReasTreatyInGroupLeader` tidak ada, dan itu dibuktikan bukan diduga

Di sistem lama, tingkat Group Leader punya **tiga perpindahan keluar dan nol perpindahan masuk**.
Sapuan menyeluruh atas ekspor menunjukkan **tidak ada satu pun aturan** yang pernah menyetel
`Position` ke nilai itu. Ia keadaan yang tidak pernah terjadi.

ADR-D-TI-0052 sudah memutuskan tidak membawanya atas dasar ketiadaan bukti. Sekarang dasarnya lebih
kuat: bukan tidak ditemukan pemakaiannya, melainkan **tidak ada jalan masuknya sama sekali**.

###### 2. `DIKEMBALIKAN` bukan keadaan — ia `DRAFT` yang punya riwayat penolakan

Sistem lama punya nilai `Reject` yang berdiri di samping `ReasTreatyInAdmin`. Pertanyaannya bukan
apakah nilai itu ada, melainkan apakah ia berbeda dalam hal apa pun **yang bukan sejarah**. Tiga
hal diperiksa, dan ketiganya dijawab oleh sapuan menyeluruh atas ekspor:

| Yang diperiksa | Hasil |
|---|---|
| apakah aturan kelengkapan yang berlaku berbeda | tidak — tidak ada aturan kelengkapan yang menyebut `Reject` |
| apakah siapa yang boleh menyuntingnya berbeda | tidak — keduanya `ReasTreatyInAdmin` |
| apakah apa yang boleh diubah berbeda | tidak — keterbukaan field dikendalikan `ViewState`, yang tidak pernah membaca `Reject` |

Kata `Reject` hanya muncul di **dua berkas** pada seluruh ekspor: aturan persetujuan
`Akseptasi_DT` dan satu kondisi tombol di `TreatyInActionButtons` — dan pada kondisi tombol itu ia
diperlakukan **identik** dengan `Accept`.

Maka `DIKEMBALIKAN` tidak disimpan sebagai keadaan. Riwayat penolakannya sudah tersimpan di
`CATATAN_PERSETUJUAN`; menyimpannya lagi sebagai keadaan berarti menyimpan satu fakta di dua
tempat, yang ADR-D-TI-0041 larang. Ini bentuk yang sama dengan `ViewState`: sesuatu yang bisa diturunkan,
disimpan sebagai penanda tersendiri — dan prinsip yang sama membuangnya.

**Yang tetap bisa dilakukan:** menyaring daftar kerja untuk "kontrak yang dikembalikan kepada saya"
tetap mungkin. Ia **turunan** dari keputusan terakhir di `CATATAN_PERSETUJUAN`, dihitung saat
ditanya (ADR-D-TI-0037), bukan bendera yang disimpan.

**Apa yang membatalkan keputusan ini:** bila bisnis menyatakan kontrak yang dikembalikan menanggung
kewajiban yang tidak ditanggung draft baru — misalnya wajib menyertakan tanggapan atas alasan
pengembalian sebelum boleh diajukan ulang — maka perbedaan itu nyata, `DIKEMBALIKAN` menjadi
keadaan tersendiri, dan **kewajiban itulah yang ditulis sebagai alasannya**. Yang tidak boleh:
dua keadaan bertahan hanya karena sistem lama punya dua nilai.

###### 3. `DITOLAK` terminal bagi versinya

**Ini keputusan perancang, bukan warisan sistem lama.** Di sistem lama `Decline` menyetel `Position`
ke kosong, sehingga pengajuan berikutnya masuk lagi dari Sec Head: penawaran yang sudah ditolak
bisa dihidupkan kembali tanpa jejak bahwa ia pernah ditolak.

Membuatnya terminal memaksa kebangkitan itu menjadi **versi baru yang terlihat**.

**Konsekuensi yang dinyatakan sadar, bukan akibat samping:** sebuah kontrak yang seluruh versinya
`DITOLAK` **tetap ada sebagai kontrak**. Ia punya identitas, punya kunci alami, dan **peringatan
duplikat akan menyebutnya** ketika penawaran serupa masuk lagi. Itu memang yang dikehendaki —
penawaran yang pernah ditolak untuk cedant, asal bisnis, periode, dan sifat proporsi yang sama
justru hal yang paling layak diperingatkan kepada orang yang memasukkannya.

**Apa yang membatalkan keputusan ini:** bila bisnis menyatakan penolakan boleh dicabut pada versi
yang sama, perpindahan `DITOLAK` — `AJUKAN ULANG` — `MENUNGGU_SEC_HEAD` ditambahkan. Satu baris,
bukan perubahan bentuk.

###### 4. Tidak ada pintu yang membatalkan keadaan terminal

Di sistem lama, `DISETUJUI` terminal di dalam mesin persetujuannya — kedua cabang luar `Akseptasi_DT`
menuntut `StatusAkseptasi != 'Resolve Complete'` — tetapi ada **empat aturan di luar mesin** yang
membatalkannya dari samping:

| Aturan | Yang dilakukannya | Terlihat oleh |
|---|---|---|
| `TreatyInForceEdit` | membuka kunci kontrak yang sudah disetujui | **setiap pengguna layar penawaran** |
| `TreatyInForceResolveComplete` | menyetel selesai disetujui tanpa approver | dua nama pengembang |
| `TreatyInReturntoInputor` | mengembalikan ke admin, status dikosongkan | dua nama pengembang |
| `TreatyInSetToDirector` | melompati dua tingkat | lewat `TreatyInTestAgent` di layar penawaran |

Di model baru tidak ada satu pun jalan menyetel keadaan selain melalui perpindahan di daftar di
atas. Keadaan bukan kolom yang bisa ditulis; ia hasil dari perpindahan yang sah.

##### Konsekuensi

- Daftar ini adalah sumber bagi kelengkapan per-perpindahan di `SPEC-INVARIAN.md`: setiap
  perpindahan punya daftar syarat yang harus terpenuhi sebelum ia boleh terjadi.
- Keadaan disimpan pada **versi kontrak**, bukan pada kontrak. Kontrak tidak punya keadaan; ia
  punya versi-versi yang masing-masing punya keadaan (ADR-D-TI-0040).
- `ViewState`, `IsEditData`, `RevisionState`, `Position` dan `StatusAkseptasi` tidak ada di model
  baru, tidak dinamai ulang, dan tidak direkonsiliasi.

---

#### Perubahan 24 September 2026 — keadaan kedelapan: `DIBATALKAN`

**Status:** diterima, 24 September 2026. Diputuskan pemilik proses.
**Sebab:** lubang L-7, ditemukan saat menyusun `5-tiket/DAFTAR-PEKERJAAN.md`.

##### Lubang yang ditutup

Daftar dua belas perpindahan di atas memberi `DRAFT` **tepat satu** jalan keluar: `AJUKAN`.
Digabung INV-25 — paling banyak satu versi tak-terminal per kontrak — akibatnya:

> Pengisi kontrak yang membuat versi karena salah pencet, atau memulai addendum yang ternyata tidak
> jadi, **tidak punya jalan keluar**. Kontraknya terkunci: versi itu tidak dapat dibuang, dan versi
> kedua tidak boleh dibuat.

Satu-satunya jalan yang tersisa adalah **mengajukannya supaya ada yang menolaknya** — memakai jalur
persetujuan sebagai tempat sampah.

**Sistem lama juga tidak punya jalan keluar yang bersih**, dan yang dipakainya merusak:
`TreatyInDeclineConfirmation_postactEDM` langkah 6 dan 7 — keduanya **hidup** — **menghapus baris**
di `M_TREATY_IN_EDM` dan `TREATY_IN_EDM` setelah addendum ditolak. Jalur kontrak biasa
(`TreatyInDeclineConfirmation_postact`, empat langkah hidup) tidak menghapus apa pun.

##### Keputusan

**Keadaan kedelapan `DIBATALKAN`, terminal, dengan satu perpindahan masuk.**

| Dari | Peristiwa | Ke |
|---|---|---|
| `DRAFT` | `BATALKAN` | `DIBATALKAN` |

Daftar keadaan menjadi **tujuh keadaan sah + satu keadaan warisan**; perpindahan menjadi
**tiga belas**.

###### Kenapa bukan memakai `DITOLAK` yang sudah ada

Meski itu tidak menambah keadaan: **`DITOLAK` berarti seseorang yang berwenang menolak.** Pengisi
yang membuang drafnya sendiri tidak ditolak siapa pun. Memakai satu keadaan untuk keduanya membuat
**setiap hitungan penolakan tercemar draf yang dibuang**, dan tidak ada cara memisahkannya kemudian.

Keadaan menyatakan apa yang terjadi; ia harus menyatakan yang benar.

###### Kenapa bukan melonggarkan INV-25

Melonggarkannya berarti **membuang invarian untuk menyelesaikan masalah alur kerja**. Hasilnya draf
terlantar yang menumpuk selamanya tanpa ada yang tahu mana yang sungguhan.

##### Tiga syarat yang mengikat

| # | Syarat | Catatan |
|---|---|---|
| **a** | **Baris tidak dihapus.** Versi yang dibatalkan tetap tersimpan, tetap membawa nomornya, tetap membawa jejaknya | jawaban langsung atas apa yang dilakukan sistem lama |
| **b** | **Nomor revisinya tidak dipakai ulang.** Lompatan penomoran itu jujur; nomor yang dipakai dua kali tidak | **tidak menuntut mekanisme baru** — lihat di bawah |
| **c** | **Hanya dari `DRAFT`, dan hanya oleh pembuatnya** | membatalkan sesudah diajukan adalah hal yang **berbeda** |

###### Syarat (b) jatuh sendiri dari (a)

**INV-04 sudah menetapkan `NOMOR_URUT_VERSI` unik di dalam satu `KONTRAK`.** Selama barisnya tidak
dihapus — syarat (a) — nomor yang sudah terpakai tetap menempati tempatnya, dan constraint yang
sudah ada menolak pemakaian ulangnya.

Jadi (b) **bukan aturan tambahan**; ia akibat (a) di bawah invarian yang sudah berlaku. Tidak ada
trigger baru, tidak ada kolom baru. Yang muncul hanyalah **lompatan** dalam deret nomor, dan
lompatan itu memang yang dikehendaki.

###### Apa yang TIDAK diputuskan di sini

**Penarikan sesudah diajukan** — membatalkan versi yang sudah berada di antrian persetujuan —
**bukan bagian keputusan ini.** Ia kemampuan tersendiri: pelakunya berbeda, dan akibatnya menyentuh
orang lain yang sudah mulai menilai. Bila ia diperlukan, ia pertanyaan tersendiri dengan jawabannya
sendiri.

Menggabungkannya sekarang adalah larangan *"satu tiket memuat hal yang pasti dan hal yang belum
diputuskan"*, diterapkan pada keputusan alih-alih pada tiket.

##### Konsekuensi yang harus dikerjakan sesi to-spec

| Berkas | Yang berubah |
|---|---|
| `SPEC-INVARIAN.md` **INV-20** | himpunan keadaan menjadi **delapan** nilai, bukan tujuh |
| `SPEC-INVARIAN.md` **INV-23** | tidak ada perpindahan keluar dari `DISETUJUI`, `DITOLAK`, **maupun `DIBATALKAN`** |
| `SPEC-INVARIAN.md` **INV-25** | `DIBATALKAN` terhitung **terminal**, sehingga versi yang dibatalkan tidak lagi menahan pembuatan versi berikutnya |
| `SPEC-INVARIAN.md` §3 | baris kelengkapan baru: `DRAFT` → `DIBATALKAN` — pelakunya **pembuat versi itu sendiri** |
| `SPEC-MODEL-DATA.md` §10.2 | `KEADAAN_SIKLUS_HIDUP` bertambah satu nilai |
| `5-tiket/DAFTAR-PEKERJAAN.md` | satu kemampuan baru, bergolongan **BARU** |

##### Satu akibat yang BELUM diputuskan, dan sengaja dibiarkan terbuka

**Kontrak yang seluruh versinya `DIBATALKAN`.** Bila versi pertama sebuah kontrak dibatalkan dan
tidak pernah ada versi kedua, kontrak itu ada tetapi tidak pernah punya isi yang berlaku.

Pertanyaannya — apakah ia muncul di pencarian (`5-tiket/DAFTAR-PEKERJAAN.md` P-56), apakah ia
terhitung sebagai kontrak dalam laporan apa pun — **tidak dijawab di sini**, karena ia pertanyaan
bisnis, bukan pertanyaan bentuk. Dicatat supaya tidak ditemukan belakangan sebagai kejutan.

### ADR-D-TI-0056 - Tidak ada stored procedure di sistem baru

**Status:** diterima, 24 September 2026 — arahan pemilik proses (K-4)
**Berlaku untuk:** seluruh skema `TREATY_MASUK`, sejak gelombang pertama
**Menggantikan:** tidak ada · **Bersandar pada:** ADR-D-CNP-0023, ADR-D-TI-0035, ADR-D-TI-0041, ADR-D-TI-0051

##### Konteks

Sistem lama menyimpan Treaty In lewat lima *stored procedure* di `POOLDATA`. Kelimanya dibedah
baris demi baris; hasilnya `2-to-spec/PEMETAAN-PROCEDURE.md`. Yang terbaca **bukan** pilihan
arsitektur melainkan satu cetakan yang disalin empat kali — dan cacatnya ikut tersalin empat kali.

Tiga bentuk cacat yang muncul berulang, dan ketiganya **tidak mungkin ada** kalau logikanya berada
di lapisan aplikasi yang dapat diuji:

1. **Pekerjaan yang tidak dilakukan, dilaporkan berhasil.** `PEGA_TREATY_IN` dipanggil dari layar
   addendum menghasilkan `UPDATE … WHERE ID = 'XXXXXXX/Rnn'` yang mengenai **nol baris**; Oracle
   tidak menimbulkan galat, `COMMIT` berhasil, `StsSave := 1`, dan pesannya berbunyi *"Data Sudah
   Disimpan Dengan ID : …"*. Kedua procedure rinci **tidak punya cabang `ELSE`** sama sekali, dan
   tetap menyetel `StsSave := 1` — jalur perbarui **tidak melakukan apa pun**, dengan `IDPegaOut`
   keluar `NULL`.
2. **`EXCEPTION WHEN OTHERS` tanpa `RAISE`, di tiga tingkat, di kelima procedure.** Setiap galat
   ditelan. Satu di antaranya (`PEGA_M_TREATY_IN_DETAIL`) tidak pernah memberi nilai kepada
   `ErrMsg` sama sekali, walaupun mendeklarasikannya `OUT`.
3. **Penyulihan teks buta atas dokumen JSON** — `replace(DataPega,'UnknownId',id)` mengganti setiap
   kemunculan kata itu di mana pun di dalam dokumen, tanpa dibatasi pada ruas pengenal.

Ketiganya bertahan bertahun-tahun karena **tidak ada satu pun yang dapat dijalankan sebagai uji**.
Procedure tidak punya seam: ia menulis, meng-`COMMIT`, dan menelan galatnya sendiri di dalam satu
benda yang hanya dapat diperiksa dengan menjalankannya terhadap basis data sungguhan.

##### Keputusan

> **Aplikasi Golang tidak memanggil stored procedure, dan skema `TREATY_MASUK` tidak memuat
> satu pun `CREATE PROCEDURE`, `CREATE FUNCTION`, maupun trigger yang membawa aturan bisnis.**

Oracle tetap dipakai sebagai basis data. Yang dipindahkan adalah **tempat aturannya tinggal**:

| Yang boleh tinggal di Oracle | Yang tidak |
|---|---|
| tipe, `NOT NULL`, `CHECK` domain, kunci utama dan asing, keunikan | keputusan cabang, penerbitan pengenal, pengambilan nilai induk, pemilihan jalur sisip-versus-perbarui |
| *materialized view* untuk invarian lintas baris (ADR sebelumnya) | pesan untuk pengguna, penanganan galat, batas transaksi |

**Batas transaksi dimiliki pemanggil.** Tidak ada `COMMIT` di dalam objek basis data.

**Kegagalan adalah kegagalan** (ADR-D-TI-0035): perintah yang tidak mengubah baris mana pun ketika ia
seharusnya mengubah satu baris **adalah galat**, bukan keberhasilan yang sunyi. Lapisan penyimpan
memeriksa jumlah baris terpengaruh dan menolak bila bukan yang diharapkan.

##### Alternatif yang ditimbang

| Pilihan | Kenapa tidak |
|---|---|
| **Pertahankan procedure, perbaiki cacatnya** | memperbaiki ketiga cacat menuntut uji, dan uji menuntut seam. Membangun seam di sekitar PL/SQL berarti membangun kerangka uji basis data untuk satu modul — ongkos yang sama dengan memindahkannya, tanpa memperoleh keterujian di tempat lain |
| **Pindahkan sekarang, biarkan procedure lama tetap ada untuk Pega** | tidak ditolak, dan **memang itu yang terjadi selama masa berdampingan**. ADR-D-TI-0041 sudah mengaturnya: hak tulis Pega dicabut per gelombang. ADR ini tentang **skema baru**, bukan tentang mencabut procedure lama lebih cepat |
| **Pakai trigger untuk integritas** | trigger dapat berhenti bekerja tanpa gagal, dan penegakan semacam itu menuntut pemantau kebasian tersendiri. Constraint tidak dapat berhenti diam-diam |

##### Konsekuensi

- **Kelima cacat §PEMETAAN-PROCEDURE lenyap karena bentuknya, bukan karena ditambal.** Tidak ada
  tempat bagi "berhasil tanpa melakukan apa pun" ketika jumlah baris terpengaruh diperiksa dan
  cabangnya ditulis sebagai kode yang punya uji.
- **Dua proyeksi datar tidak dibawa.** `TREATYINDETAIL` dan `TREATYINDETAILEDM` adalah turunan dari
  dokumen JSON; di skema baru bentuk kanoniknya relasional, sehingga proyeksinya tidak perlu
  ditulis ulang oleh siapa pun. Ketimpangan tujuh kolom di antara keduanya ikut lenyap.
- **Konsumen hilir yang membaca kedua tabel datar itu harus dilayani** — lewat *view* atas bentuk
  kanonik baru, bukan lewat procedure yang menyalin (ADR-D-TI-0041, ADR-D-TI-0051). **View itu tidak
  mengembalikan ketimpangan tujuh kolom**: ia memaparkan besaran yang sama untuk kontrak maupun
  addendum.
- **Penerbitan pengenal pindah ke aplikasi**, dan bersamanya kesempatan memperbaiki `lpad(...,6,'0')`
  yang membuat skema pengenal lama pecah. Bentuk penggantinya bukan urusan ADR ini.
- **Satu hal yang TIDAK diputuskan di sini:** apa yang terjadi pada kelima procedure lama. Mereka
  tetap melayani Pega sampai gelombangnya pindah, dan pencabutannya diatur ADR-D-TI-0041 langkah 1.
- **Satu temuan diangkat keluar dari lingkup Treaty In.** `GET_TOKEN_STORAGE` menerbitkan token
  akses dari `STANDARD_HASH('ASMAPP' || timestamp,'MD5')` — **dapat ditebak sepenuhnya, tanpa unsur
  acak**. Itu bukan utang migrasi; ia terbuka sekarang, di sistem yang berjalan. Ia naik ke
  `DAFTAR-ESKALASI-MANAJEMEN.md` dan **tidak** ditambal oleh sesi ini.

##### Syarat pembalikan

Keputusan ini dibalik bila salah satu terbukti:

1. ada konsumen yang **memanggil** procedure ini secara langsung — bukan membaca tabelnya — dan
   tidak dapat diubah; maka yang dibutuhkan **pembungkus**, dan bentuknya diputuskan tersendiri;
2. ada kewajiban kepatuhan yang menuntut aturan tertentu ditegakkan di dalam basis data, di luar apa
   yang dapat dinyatakan sebagai constraint.

**Tidak satu pun ditandai terverifikasi.** Keduanya belum diperiksa, dan pemeriksaannya menunggu
daftar konsumen hilir yang sampai sekarang tidak ada (ADR-D-TI-0051).

---

# Bagian III - Keputusan Fakultatif

Jumlah: **7** keputusan.

## ADR-F-0001 - Rekonsiliasi paralel run eksak, dicapai bertahap

Ukuran keberhasilan migrasi adalah **rekonsiliasi paralel run dengan nol selisih sampai digit
terakhir** — bukan toleransi. Karena seluruh jalur tulis produksi melewati 35 stored procedure
`POOLDATA` yang isinya tidak kita miliki, membaca kode saja tidak dapat membuktikan port-nya benar;
hanya membandingkan keluaran dua sistem atas masukan yang sama yang bisa. Paralel run **penuh** baru
mungkin setelah ekspor produksi tunggal tiba, sehingga sampai saat itu dijalankan bertahap.

#### Considered Options

- **Toleransi seragam (mis. ±0,01)** — ditolak. Bila urutan operasi dan presisi per-langkah
  direproduksi apa adanya, `decimal` bersifat deterministik dan hasilnya **harus** identik. Selisih
  sekecil apa pun berarti ada salah-port, dan toleransi justru menyembunyikannya. Toleransi baru
  masuk akal bila pembulatan sengaja diseragamkan — dan itu perbaikan, bukan migrasi
  (`CLAUDE.md` §1).
- **Tanpa paralel run, cutover langsung dengan UAT** — ditolak. Tanpa isi 35 stored procedure, tidak
  ada dasar lain untuk membuktikan kebenaran jalur produksi.

#### Consequences

Tahapan yang disepakati:

1. **Sekarang** — perhitungan murni saja (premi, spreading, nilai dasar akseptasi) atas masukan yang
   direkam. Tidak menyentuh alur, tidak butuh ekspor produksi.
2. **Setelah mesin akseptasi ditulis** (lihat ADR-F-0003) — per-modul, termasuk tangga akseptasi
   dengan fixture tabel limit.
3. **Setelah `ALL_SOURCE` 35 prosedur tiba** — kasus nyata end-to-end.

Konsekuensi yang mengikat porting: urutan operasi dan presisi pembulatan **tidak boleh** diubah
"agar lebih rapi". Lihat ADR-F-0005.

## ADR-F-0002 - Satu flag fase: `panic` saat paralel run, `decline`+log saat produksi

Baris keputusan `DecisionTable/IsUWAccepted` tidak ikut terekspor, sehingga pemetaan **status →
konektor flow** tidak terbaca meskipun arti tiap nilai status sudah terverifikasi. Untuk jalur yang
belum terverifikasi, sistem berperilaku berbeda menurut fase: **`panic` selama paralel run** (agar
setiap jalur yang belum terbukti muncul ke permukaan — itu memang tujuan paralel run), dan
**`decline` + catatan log saat produksi** (karena `decline` adalah perilaku terekam sistem lama,
dan menjatuhkan case produksi tanpa alasan bisnis lebih berbahaya). Keduanya dikendalikan **satu
flag fase eksplisit**, bukan dua basis kode.

#### Considered Options

- **`panic` di kedua fase** — ditolak untuk jalur tipe (a) di bawah. Paralel run akan berhenti pada
  kasus yang di sistem lama berjalan mulus sebagai `decline`, dan di produksi ia menjatuhkan case
  yang seharusnya lanjut.
- **`decline` di kedua fase** — ditolak. Menyembunyikan celah pengetahuan kita justru pada fase yang
  seharusnya menemukannya.

#### Consequences

Dua jenis celah dibedakan tegas, dan **hanya yang pertama** mengikuti flag fase:

| | Situasi | Paralel run | Produksi |
| --- | --- | --- | --- |
| **(a)** | Status **dikenali** artinya, baris keputusannya belum terverifikasi | `panic` | `decline` + log |
| **(b)** | Kondisi **tidak dikenali sama sekali** | `panic` | **`panic`** |

Transisi sebuah jalur (a) dari `panic` menjadi `decline`+log **hanya** boleh dilakukan setelah jalur
itu diverifikasi terhadap ekspor produksi tunggal — bukan karena ia sering muncul dan mengganggu.

`[terverifikasi]` `decline` adalah `<pyDefaultResult>` rule `IsUWAccepted`, jadi memilihnya di
produksi adalah reproduksi, bukan tebakan.

Mesin keadaannya sendiri dibangun **dari sisi penulis status**, yang seluruhnya terverifikasi berikut
efek sampingnya (mis. rule penulis `2` mematikan `ConfirmBinding` dan `ReceivedRiSlip`; penulis `4`
menyalakan `IsBanding` dan menyetel tujuan banding). Tabel keputusan hanya menentukan pemetaan
status → konektor, dan itu diverifikasi belakangan tanpa menahan pekerjaan.

## ADR-F-0003 - Mesin akseptasi ditulis lebih dulu; fixture tabel limit menjadi kontrak data ke DBA

`services/acceptance` ditulis **sekarang** dengan fixture tabel limit, tidak menunggu isi
`M_LIMIT_*` tiba dari DBA. Yang hilang dari tangga akseptasi adalah **datanya**, bukan **logikanya**:
query SQL-nya terbaca penuh, urutan eskalasinya terverifikasi (`LIMIT_BOTTOM` menaik, baris pertama =
approver berikutnya), cara berhentinya terverifikasi (tidak ada `To*` yang cocok → cabang `Else` →
tangga selesai, bukan galat), dan ketiga folder korpus berperilaku identik. Fixture yang kita tulis
sekaligus **menjadi kontrak bentuk data** yang dikirim ke DBA.

#### Consequences

Ini membalik arah ketergantungan menjadi menguntungkan: alih-alih menunggu data lalu menemukan
bentuknya tidak seperti dugaan, kita menyatakan bentuk yang kita harapkan lebih dulu dan
ketidakcocokan muncul **saat permintaan dikirim**, bukan di akhir.

~~Fixture wajib memuat kolom yang dibaca query: `JABATAN`, `JABATAN_ATASAN`, `LIMIT_BOTTOM`,
`LIMIT_BOTTOM2`, `MAX_LIMIT_IDR`, `MAX_LIMIT_USD`, `BATAS_WAKTU`, `LOGIN` — untuk ketujuh tabel limit
per lini bisnis.~~ → **dikoreksi amandemen di bawah: ada DUA bentuk, bukan satu, dan enam tabel, bukan tujuh.**

~~⚠️ Satu hal yang fixture **tidak** dapat tebak dan harus datang dari DBA: **ejaan pasti nilai kolom
`JABATAN`**.~~ → **sudah datang; lihat amandemen.** Token routing di rule ditulis tanpa spasi
(`DIREKTURTEKNIK`); bila ejaan Oracle berbeda, tidak ada approver yang pernah ditemukan dan tangga
macet total.

Keputusan ini **mencabut** sikap sebelumnya yang menunda `services/acceptance` sampai tabel limit
tiba (lihat K-009).

---

#### Amandemen — 17 September 2026, setelah DDL + isi tabel limit diterima

**Keputusan intinya tidak berubah.** Menulis fixture lebih dulu terbukti benar: ketidakcocokan bentuk
memang muncul saat data datang, dan itu tepat tujuannya. Yang berubah adalah **asumsi bentuknya**.

Sumber: `D:\migrasi\RNM\DDL\M_LIMIT_*.txt` (DDL) dan `M_LIMIT_*.xls` (isi), diverifikasi
17 September 2026.

##### ⛔ Asumsi "satu bentuk untuk ketujuh tabel" **SALAH**

`[terverifikasi]` Ada **enam** tabel limit di basis data, dalam **dua bentuk yang tidak kompatibel**.

###### Bentuk A — standar, 14 kolom, 5 tabel

`M_LIMIT_PROPERTYY` · `M_LIMIT_ENGINEERINGG` · `M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL` ·
`M_LIMIT_PROPERTY_NON_PREFERREDD` · `M_LIMIT_NONPROPANDENGG`

```
ID · JABATAN · MAX_LIMIT_IDR · LIMIT_BOTTOM · MAX_LIMIT_USD · BATAS_WAKTU · NAMA
TGL_UPDATE · EFFECTIVE_DATE · TEAM_GROUP · LOGIN · WORKBASKET · JABATAN_ATASAN · LIMIT_BOTTOM2
```

`[terverifikasi]` **Ejaan `JABATAN` TANPA SPASI** — cocok dengan token routing di rule:
`DIREKTURTEKNIK` · `DIREKTURMARKETING` · `KADIVTEKNIK` · `KADIVFACULTATIVE` · `DEPHEADUNDERWRITER` ·
`MANAGERTEKNIK` · `SENIORUW` · `UNDERWRITER` · `LEADER` · `JUW_B`.

`WORKBASKET` **terisi** dan cocok dengan ruang nama antrean (`ReasFacInUnderwriting`,
`ReasFacInSeniorUnderwriting`, …), sehingga kedua ruang nama — antrean dan kode jabatan — tersedia
langsung dari tabel ini.

###### Bentuk B — financial, 7 kolom, 1 tabel

`M_LIMIT_FINANCIALINS`

```
ID · JABATAN · LIMITBOND_BOTTOM · LIMITCREDITCL_BOTTOM · LIMITCREDITNCL_BOTTOM · LIMITTRADE_BOTTOM · NAMA
```

`[terverifikasi]` Bentuk ini **tidak punya** `WORKBASKET`, `JABATAN_ATASAN`, `TEAM_GROUP`,
`EFFECTIVE_DATE`, maupun `LOGIN`. Kolom limitnya **sama sekali berbeda**: empat ambang per jenis
pertanggungan (bond · credit CL · credit NCL · trade), bukan `MAX_LIMIT_IDR`/`LIMIT_BOTTOM`.

⛔ `[terverifikasi]` **Ejaan `JABATAN` PAKAI SPASI**: `DIREKTUR TEKNIK` · `DIREKTUR MARKETING` ·
`KADIV KEUANGAN` · `SENIOR UNDERWRITER`.

##### Konsekuensi yang mengikat implementasi

1. **Mesin akseptasi financial mencocokkan jabatan berspasi, bukan token tanpa spasi.** Bila kedua
   bentuk disamakan — satu fungsi pencocokan untuk semua lini — **tidak ada approver financial yang
   pernah ditemukan**, dan seluruh tangga lini financial macet. Ini persis kegagalan yang ADR ini
   peringatkan; ia nyata, dan sekarang terbukti.
2. **Normalisasi ejaan DILARANG.** Menghapus spasi agar "seragam" adalah perubahan perilaku, bukan
   migrasi (`CLAUDE.md` §1). Kedua ejaan direproduksi apa adanya, masing-masing pada bentuknya.
3. **Bentuk B tidak punya rantai atasan.** Tanpa `JABATAN_ATASAN`, tangga financial **tidak dapat
   menaik lewat mekanisme yang sama**. Bagaimana eskalasinya bekerja — atau apakah memang tidak ada
   eskalasi — **belum terjawab dari data ini** dan ditandai `[pertanyaan terbuka]`.
4. **Bentuk B tidak punya antrean.** Tanpa `WORKBASKET`, tidak ada token antrean yang dapat dibaca
   untuk lini financial.
5. **Fixture tetap dipakai untuk bentuk B** sampai butir 3 dan 4 terjawab; bentuk A boleh berpindah ke
   data nyata.

##### `M_LIMIT_LIFE` bukan kekurangan

`[terverifikasi]` Dikonfirmasi **tidak ada** di basis data. Target efektif **enam** tabel, bukan tujuh.

##### Kolom yang tidak boleh masuk fixture

⛔ Kolom **`NAMA`** dan **`LOGIN`** memuat **nama orang**. Keduanya **wajib dibuang** saat fixture
dibangun dari berkas isi tabel — `CLAUDE.md` §3 butir 5 berlaku penuh pada fixture dan test. Lihat
**K-025**.

## ADR-F-0004 - `Money` dan `Ratio` adalah dua tipe yang tidak dapat dijumlahkan; skala melekat pada nilai

Nilai uang memakai `Money{Amount decimal, Currency}`; rate dan persen memakai
`Ratio{Value decimal, Scale}` yang **membawa satuannya sendiri**. Keduanya tipe berbeda: `Money +
Ratio` tidak dapat dikompilasi, dan satu-satunya jembatan adalah operasi eksplisit `Money × Ratio →
Money`. Skala sebuah `Ratio` diisi oleh **satu resolver terpusat** yang memetakan lini bisnis (COB) →
skala.

#### Mengapa, dengan bukti

`[terverifikasi]` Properti `.Rate` yang sama memakai **dua satuan berbeda** menurut COB: FIRE dan PA
per mille (‰); ANEKA, BONDING, GOLF, MARINE CARGO, dan MBU persen (%). Buktinya menutup kemungkinan
"salah satu rule bug" — satu berkas memuat tiga COB berdampingan dengan dua satuan, dan label layar
menyebut satuannya harfiah (`‰ Standard Rate` dengan karakter U+2030 pada layar FIRE).

Menyeragamkan satuan menggeser premi **satu ordo besaran 10** pada separuh portofolio. Tipe yang
membawa skalanya sendiri mengubah cacat itu dari kesalahan runtime yang senyap menjadi **kegagalan
kompilasi**.

#### Considered Options

- **Satu tipe `Money` untuk semuanya, persen sebagai `decimal` telanjang** — ditolak. Justru
  menghapus proteksi yang paling dibutuhkan.
- **Empat tipe (`Money`, `Rate`, `Percent`, `Limit`)** — ditolak. `Limit` berperilaku persis seperti
  `Money`; tipe tambahan tanpa perilaku tambahan hanya menambah friksi.
- **Skala disimpan di tabel konfigurasi Oracle** — ditolak, dan ini penolakan yang paling penting
  untuk diingat. Skala satuan adalah **fakta struktural yang terverifikasi dari korpus**, bukan
  parameter bisnis. Menaruhnya di tabel yang dapat diubah tanpa deployment mengulang persis
  kerentanan yang sudah kita temukan pada `M_PROMPT_AI`: perilaku sistem berubah tanpa jejak
  version control.

#### Consequences

##### Aturan penguraian pembagi komposit — wajib dipakai saat membaca rumus premi

Ini satu-satunya cara membaca satuan dengan benar, dan mengabaikannya menghasilkan kesimpulan yang
terbalik.

**Pembagi gabungan pada satu `@Math.divide` adalah hasil kali faktor-faktor satuan, bukan satu
satuan tunggal.** Satuan `.Rate` disimpulkan dari **faktor yang menempel padanya saja** — tidak
pernah dari pembagi total.

| Faktor dalam ekspresi | Sumbangan ke pembagi |
| --- | ---: |
| `.Rate` ber-satuan ‰ | 1.000 |
| `.Rate` ber-satuan % | 100 |
| `ProRatePercent` (persen) | 100 |

**Contoh terverifikasi** — `GenerateLayerList_ACT.xml` L909:

```
@Math.divide((.TSI * .Rate * pyWorkPage.OfferFacIn.ProRatePercent),100000,4)
        100000  =  1000 (.Rate ber-‰)  ×  100 (ProRatePercent ber-%)
```

Pembagi `100000` karena itu **menegaskan** `.Rate` = ‰ — bukan membantahnya. Siapa pun yang membaca
`100000` sebagai "satuan yang lebih kecil dari ‰" akan salah satu ordo besaran.

⚠️ **Struktur ekspresi MBU membuktikan aturan ini secara langsung.**
`FillPremiMBU_FacIn.xml` L1144 memakai bentuk **bersarang**, sehingga keterikatan faktor→pembagi
tertulis eksplisit dan bukan tafsir:

```
@Math.divide( @Math.divide((.TSI*(.Rate+.Loading)),100,4) * @if(.ProRatePercent=="",100,.ProRatePercent), 100, 4)
                                                  ^^^                                                    ^^^
                                      pembagi DALAM menempel .Rate (%)              pembagi LUAR menempel ProRatePercent
```

##### Peta skala per COB — terkunci

Nilainya ditetapkan **K-018** (`00-KEPUTUSAN-WORK-OWNER.md`), berikut kutipan baris sumbernya:

| COB | Skala `.Rate` | Pembagi yang menempel |
| --- | :-: | ---: |
| FIRE · PA · **Layering** | **‰** | 1.000 |
| MBU · ANEKA · BONDING · GOLF · MARINE CARGO | **%** | 100 |

##### Lain-lain

- Resolver COB → skala mengikuti **rumus**, bukan label layar. `[terverifikasi]` tiga label UI
  bertentangan dengan rumusnya — dan K-018 menetapkan **rumus yang menang** di ketiganya. Perbaikan
  labelnya tidak mengubah satu angka pun, sehingga tidak mengganggu rekonsiliasi.
- COB yang belum ada di peta resolver → **`panic`**, bukan skala default. Menebak skala adalah persis
  cacat 10× yang tipe ini dirancang untuk mencegah.

## ADR-F-0005 - Presisi pembulatan ditulis literal di tiap langkah, bukan disentralkan

Setiap pembulatan hasil port menuliskan presisinya **sebagai literal di tempatnya**, disertai
komentar `// Asal: <rule>, langkah <n>, presisi <p>` (`CLAUDE.md` §4.6). Presisi **tidak**
disentralkan ke registry maupun dilekatkan pada nilai. Ini satu-satunya tempat dalam rancangan ini
yang memilih duplikasi ketimbang sentralisasi, dan alasannya spesifik.

#### Mengapa

`[terverifikasi]` Di sistem lama, presisi **bukan properti nilai maupun properti rule** — ia properti
**satu langkah tertentu di dalam** sebuah rule:

- Aritmetika yang sama dibulatkan **4 desimal** di satu rule dan **20 desimal** di rule lain.
- Satu pembagi muncul dengan **9 presisi berbeda**.
- Sebuah pembulatan berada **di dalam** loop akumulasi, sehingga galatnya menumpuk per iterasi.

#### Considered Options

- **Registry terpusat `nama rule → presisi`** — ditolak. Ia mengandaikan satu rule punya satu
  presisi, dan itu tidak benar.
- **Presisi melekat pada nilai dan mengalir lewat operasi** — ditolak. Presisi tidak mengalir; ia
  diterapkan di satu titik lalu selesai.

#### Consequences

Uji rekonsiliasi **wajib** memuat kasus khusus untuk pembulatan **di dalam loop akumulasi**, bukan
hanya nilai tunggal. Pembulatan per-iterasi menumpuk galat, sehingga port yang benar untuk satu nilai
masih bisa meleset untuk daftar panjang — dan itu justru bentuk kesalahan yang paling sulit terlihat.

Menyeragamkan presisi adalah **perbaikan terpisah**, bukan bagian migrasi.

## ADR-F-0006 - Mata uang yang tidak diketahui adalah keadaan eksplisit, bukan default

`Money.Currency` mengizinkan keadaan **`Unknown` yang eksplisit**. Nilai ber-mata-uang-tidak-diketahui
boleh dibaca dan ditampilkan, tetapi **`panic`** ketika dipakai dalam aritmetika lintas mata uang
atau ditulis ke Oracle. Tidak ada default diam-diam.

#### Mengapa

`[terverifikasi]` Korpus memang kehilangan informasi ini: **112 Section menampilkan nilai uang tanpa
field mata uang mana pun**, dan hanya ada 93 pengikatan properti mata uang di seluruh NB. Sekaligus
`[terverifikasi]` sistem lama **tidak pernah** menetapkan mata uang default di kode — ia selalu
datang dari data (`CurrencyList` / `CurrencyMaster`), dan kurs selalu dari hasil query.

#### Considered Options

- **Default ke IDR** — ditolak. Menebak (`CLAUDE.md` §3 butir 4), dan tebakan itu akan salah persis
  pada kasus yang paling mahal: penempatan valas.
- **`panic` tanpa kecuali bila mata uang tidak ada** — ditolak. Akan menghentikan aplikasi pada 112
  layar yang di sistem lama berjalan mulus — kesalahan yang sama bentuknya dengan mengganti
  `decline` terekam menjadi `panic` (lihat ADR-F-0002).

#### Consequences

Prinsipnya sama dengan ADR-F-0002: celah pengetahuan dibuat **eksplisit dan terlihat**, lalu gagal
keras **hanya** di titik yang benar-benar memerlukan jawabannya.

| Operasi | Mata uang `Unknown` |
| --- | --- |
| Baca, tampilkan, simpan di memori | diizinkan |
| Aritmetika dalam satu mata uang yang sama-sama `Unknown` | diizinkan |
| Aritmetika lintas mata uang | **`panic`** |
| Tulis ke Oracle | **`panic`** |

Efek samping yang disengaja: implementasi menghasilkan **hitungan berapa banyak nilai produksi yang
tiba tanpa mata uang**. Angka itu masuk kuesioner — ia mengubah pertanyaan dari dugaan menjadi
ukuran.

## ADR-F-0007 - Total yang menjumlahkan lintas mata uang bukan `Money`, melainkan tipe ketiga yang tidak dapat diaritmetikakan

Lima total skalar di akar penawaran — `TOTAL_TSI_NUSA_RE`, `TOTAL_PREMI_NUSA_RE`,
`TOTAL_TSI_NUSA_RE_SPREADING`, `TOTAL_TSI_TOP_RISK`, `SUM_TOTAL_TSI` — **tidak dimodelkan sebagai
`Money`**. Keduanya tidak dapat disatukan: **ADR-F-0004** mensyaratkan `Money` membawa mata uangnya,
dan di sini **tidak ada satu mata uang pun yang benar untuk dibawa**.

Kelimanya memakai **tipe ketiga** di samping `Money` dan `Ratio`. Sifatnya:

- menyimpan **angka apa adanya**, tanpa mata uang;
- ⛔ **tidak boleh dijumlahkan dengan `Money`**;
- ⛔ **tidak boleh dikonversi**;
- ⛔ **tidak boleh dibandingkan** dengan nilai bermata-uang.

Boleh dibaca, ditampilkan, disimpan, dan direkonsiliasi terhadap nilai lama. Tidak boleh masuk
aritmetika bisnis.

#### Mengapa, dengan bukti

`[terverifikasi]` `DDL\CONTOH\NB-173649.xml` adalah penawaran **dua mata uang** — `CurrencyList`
memuat entri **USD** dan **IDR**, dan `PropertyItemList` membawa `Currency` sendiri per baris
(**3 IDR + 2 USD**).

Pada berkas itu, dua kesetaraan terukur **persis**:

```
TotalTSINusaReSpreading  =  jumlah 2 baris TotalTSIPremiSpreadRNM.TSISpreaded
                            TreatyName baris 1 = USD
                            TreatyName baris 2 = IDR
TotalPremiNusaRe         =  jumlah 17 baris CoverageList.PremiNusantaraRe
SumTotalTSI              =  satu nilai LocationList/Property.TotalTSI
```

⛔ **Baris USD dan baris IDR dijumlahkan langsung, tanpa kurs.** Angka yang keluar karena itu bukan
jumlah uang dalam mata uang mana pun — ia jumlah aritmetis atas dua satuan berbeda.

Dua sisanya tidak cocok agregat sederhana mana pun dan **tidak ditebak**:

| Skalar | Diuji terhadap | Hasil |
| --- | --- | :-: |
| `TotalTSINusaRe` | jumlah `CoverageList.TSINusantaraRe` | **BEDA** |
| `TotalTSITopRisk` | `LocationList/Property.TotalTSI` | **BEDA** (bernilai nol di berkas ini) |

Keduanya tetap masuk tipe ketiga: asal-usulnya belum terbaca, sehingga **tidak ada dasar** untuk
menempelkan mata uang pada keduanya.

⚠️ **Yang tidak pecah hanya kelima skalar akar ini.** Struktur lain sudah menangani multi mata uang
dengan benar — `CurrencyList`, `TotalTSIList`, `TotalTSIPremiGrossList` dan
`TotalTSIPremiSpreadRNM` masing-masing **2 baris** di berkas itu, satu per mata uang. Jadi masalahnya
**bukan** bahwa sistem lama tidak mengenal multi mata uang; masalahnya kelima skalar ini meratakannya.

#### Considered Options

- **Modelkan sebagai `Money` dengan `Currency = Unknown` (ADR-F-0006).** ⛔ Ditolak, dan ini penolakan
  yang paling penting. `Unknown` di ADR-F-0006 berarti *"mata uangnya ada, kita belum tahu apa"* —
  keadaan pengetahuan yang bisa diperbaiki dengan data. Di sini keadaannya berbeda secara jenis:
  **mata uangnya memang tidak ada**, karena angkanya campuran. Memakai `Unknown` menyamarkan cacat
  struktural menjadi celah pengetahuan, dan membuka jalan bagi seseorang kelak "mengisi" mata uang
  yang benar — yang justru **mengubah angka**.
- **Modelkan sebagai `Money` ber-`Currency = IDR`.** ⛔ Ditolak. Menebak (`CLAUDE.md` §3 butir 4),
  dan tebakannya salah persis pada kasus yang paling mahal: penempatan valas.
- **Perbaiki sekarang — konversi ke satu mata uang sebelum dijumlahkan.** ⛔ Ditolak. Itu
  **perbaikan, bukan migrasi** (`CLAUDE.md` §1). Ia mengubah angka yang tersimpan, sehingga
  rekonsiliasi paralel run akan berbeda **dan menyembunyikan salah-port yang sesungguhnya**
  (ADR-F-0001).
- **`decimal` telanjang.** ⛔ Ditolak dengan alasan yang sama seperti ADR-F-0004 menolak persen sebagai
  `decimal` telanjang: menghapus justru proteksi yang paling dibutuhkan. Tanpa tipe, tidak ada yang
  mencegah `TOTAL_PREMI_NUSA_RE` dijumlahkan dengan `Money` pada layar ringkasan.

#### Consequences

**Kegagalan dipindahkan ke waktu kompilasi.** Sejalan ADR-F-0004: `Money + Ratio` tidak dapat
dikompilasi, dan kini `Money + <tipe ketiga>` juga tidak. Cacat yang tadinya senyap — menjumlahkan
total campuran dengan nilai bermata-uang — menjadi **kegagalan kompilasi**.

**Kandidat perbaikan bisnis, dicatat bukan dieksekusi.** Bahwa sistem lama menjumlahkan USD dan IDR
tanpa konversi adalah temuan yang layak dibawa ke bisnis. Sampai bisnis memutuskan, kelimanya
**diport apa adanya**.

**Batas terhadap ADR-F-0006.** `Unknown` tetap dipakai untuk kasus yang dimaksud ADR-F-0006 — nilai uang
bermata-uang-satu yang informasinya hilang dari korpus. Tipe ketiga ini **bukan** perluasan
`Unknown`; ia tipe lain, untuk angka yang **memang tidak punya** mata uang tunggal.

##### ⚠️ Akibat untuk butir 8 — yang MASIH TERBUKA dan tidak diputuskan di sini

`TotalTSINusaReSpreading` adalah **jumlah dari `TotalTSIPremiSpreadRNM`** — tabel yang **V-28**
rencanakan dibuang.

⛔ Bila tabel itu dibuang lalu nilainya dihitung ulang dari `T_SPREADINGLIST`, **hitung ulang itu
WAJIB ikut menjumlahkan lintas mata uang** agar angkanya sama. Bila tidak, selisihnya **bukan soal
pembulatan lagi melainkan beda hasil** — dan tidak ada ambang toleransi yang dapat menutupinya.

📌 `TreatyName` di bawah `TotalTSIPremiSpreadRNM` adalah **pembawa kode mata uang**. Membuang tabel
itu ikut membuang pembawanya. **Butir 8 dan butir 9 saling menyentuh**; butir 8 diputuskan terpisah.

---

# Lampiran 1 - Konkordansi

Seluruh **105** keputusan, nomor asli menjadi kode berseri.

| Kode di dokumen ini | Nomor asli | Modul | Judul | Berkas asal |
| --- | ---: | --- | --- | --- |
| `ADR-U-0001` | 0001 | lintas-modul | Claim — Life dan Komite Life adalah dua konteks terpisah, dihubu | `docs\adr\0001-batas-konteks-claim-life-komite-life.md` |
| `ADR-U-0002` | 0002 | lintas-modul | RBAC Claim — Life memakai tiga peran yang sudah ada; rangkap per | `docs\adr\0002-rbac-tiga-peran-claim-life.md` |
| `ADR-U-0003` | 0003 | lintas-modul | Uang di Claim — Life tidak direpresentasikan sebagai `float` | `docs\adr\0003-representasi-uang-non-float.md` |
| `ADR-U-0004` | 0004 | lintas-modul | Alamat layanan keluar menjadi env var, bukan replikasi lookup `M | `docs\adr\0004-endpoint-sebagai-env-var.md` |
| `ADR-U-0005` | 0005 | lintas-modul | `IsPEGAPROD` tidak ditiru sebagai rule; diganti flag lingkungan  | `docs\adr\0005-ispegaprod-menjadi-flag-lingkungan.md` |
| `ADR-U-0006` | 0006 | lintas-modul | Penomoran klaim Life memanggil `PROC_GENERATE_SEQUENCE_NUMBER`;  | `docs\adr\0006-penomoran-klaim-via-stored-procedure.md` |
| `ADR-U-0007` | 0007 | lintas-modul | Jejak audit merekam siapa + kapan untuk setiap transisi status d | `docs\adr\0007-jejak-audit-setiap-transisi-dan-jalur-balik.md` |
| `ADR-U-0008` | 0008 | lintas-modul | Efek keluar Claim — Life dijalankan asinkron, tidak memblokir al | `docs\adr\0008-efek-keluar-asinkron-dengan-antre-ulang.md` |
| `ADR-U-0009` | 0009 | lintas-modul | Seluruh data Claim — Life dipindahkan; tidak ada koeksistensi du | `docs\adr\0009-migrasi-penuh-data-klaim-life.md` |
| `ADR-U-0010` | 0010 | lintas-modul | Penyimpanan berkas klaim tetap di Google Storage | `docs\adr\0010-penyimpanan-berkas-tetap-di-google-storage.md` |
| `ADR-U-0011` | 0011 | lintas-modul | Unit status Claim — Life adalah **baris `AdjustmentList`**, buka | `docs\adr\0011-unit-status-adalah-baris-adjustmentlist.md` |
| `ADR-U-0012` | 0012 | lintas-modul | Wewenang kirim ke Komite bergantung `Type` — dibawa apa adanya,  | `docs\adr\0012-wewenang-kirim-komite-bergantung-type.md` |
| `ADR-U-0013` | 0013 | lintas-modul | Alamat layanan keluar di-resolve runtime dari `M_LINK_SERVICE` — | `docs\adr\0013-resolusi-endpoint-via-m-link-service.md` |
| `ADR-U-0014` | 0014 | lintas-modul | Keputusan komite hanya boleh diambil pemilik `KomiteID` pada tin | `docs\adr\0014-pemutus-komite-ditegakkan-per-komiteid-tingkat-berjalan.md` |
| `ADR-U-0015` | 0015 | lintas-modul | Efek keluar Komite Claim Life **wajib berhasil** — transactional | `docs\adr\0015-efek-keluar-komite-wajib-berhasil-transactional-outbox.md` |
| `ADR-U-0016` | 0016 | lintas-modul | Kolom uang memakai `NUMBER(38,8)` | `docs\adr\0016-presisi-uang-number-38-8.md` |
| `ADR-U-0017` | 0017 | lintas-modul | Daftar medan disusun dari aturan, bukan dari panduan bentuk doku | `docs\adr\0017-daftar-medan-dari-aturan-bukan-panduan-dokumen.md` |
| `ADR-U-0018` | 0018 | lintas-modul | Generasi polis adalah rantai baris, bukan penyuntingan di tempat | `docs\adr\0018-generasi-polis-sebagai-rantai-baris.md` |
| `ADR-U-0019` | 0019 | lintas-modul | Baris anak dipasangkan antar generasi lewat nomor urut, bukan ku | `docs\adr\0019-nourut-sebagai-pemasang-baris-antar-generasi.md` |
| `ADR-U-0020` | 0020 | lintas-modul | Tabel proyeksi selisih bukan tabel sumber | `docs\adr\0020-tabel-proyeksi-bukan-tabel-sumber.md` |
| `ADR-U-0021` | 0021 | lintas-modul | Rumus selisih hidup di lapisan layanan, dan nilai lama tidak per | `docs\adr\0021-rumus-selisih-di-lapisan-layanan.md` |
| `ADR-U-0022` | 0022 | lintas-modul | Konversi tipe terjadi sekali saat masuk, bukan tiap kali dibaca | `docs\adr\0022-konversi-tipe-sekali-saat-masuk.md` |
| `ADR-U-0023` | 0023 | lintas-modul | Pemecah dokumen wajib punya penampung medan tak dikenal | `docs\adr\0023-penampung-medan-tak-dikenal.md` |
| `ADR-U-0024` | 0024 | lintas-modul | Pembatalan adalah generasi bernilai nol, bukan penghapusan | `docs\adr\0024-pembatalan-sebagai-generasi-bernilai-nol.md` |
| `ADR-U-0025` | 0025 | lintas-modul | Tabel inti berbagi kunci utama dengan tabel kerja | `docs\adr\0025-kunci-utama-bersama-antara-tabel-kerja-dan-tabel-inti.md` |
| `ADR-U-0026` | 0026 | lintas-modul | Tabel kerja adalah akar lintas-lini, bukan tabel di samping | `docs\adr\0026-tabel-kerja-lintas-lini-sebagai-akar.md` |
| `ADR-U-0027` | 0027 | lintas-modul | Kolom basis data nullable; kewajiban isi ditegakkan di kode | `docs\adr\0027-kolom-nullable-wajib-isi-di-lapisan-layanan.md` |
| `ADR-U-0028` | 0028 | lintas-modul | Sistem hilir membaca tabel, bukan muatan JSON | `docs\adr\0028-hilir-membaca-tabel-bukan-muatan-json.md` |
| `ADR-U-0029` | 0029 | lintas-modul | Batas transaksi dipegang aplikasi; `COMMIT` tidak tertanam di te | `docs\adr\0029-batas-transaksi-dan-commit-dipegang-aplikasi.md` |
| `ADR-U-0030` | 0030 | lintas-modul | Aturan peran ditetapkan sekali dan berlaku lintas modul | `docs\adr\0030-aturan-peran-ditetapkan-sekali-lintas-modul.md` |
| `ADR-U-0031` | 0031 | lintas-modul | Penghapusan adalah penanda dan nilai balik, bukan hapus fisik | `docs\adr\0031-penghapusan-logis-bukan-hapus-fisik.md` |
| `ADR-U-0032` | 0032 | lintas-modul | Nama yang menyesatkan dibetulkan, disertai tabel pemetaan nama l | `docs\adr\0032-nama-menyesatkan-dibetulkan-dengan-tabel-pemetaan.md` |
| `ADR-U-0033` | 0033 | lintas-modul | Nama skema ditulis eksplisit pada setiap query | `docs\adr\0033-skema-oracle-ditulis-eksplisit.md` |
| `ADR-U-0034` | 0034 | lintas-modul | Uang melintasi batas stored procedure sebagai teks, dan dikembal | `docs\adr\0034-uang-melintasi-batas-procedure-sebagai-teks.md` |
| `ADR-U-0035` | 0035 | lintas-modul | Tetapan operasional dibaca dari tabel, bukan ditanam di kode | `docs\adr\0035-tetapan-operasional-dibaca-dari-tabel.md` |
| `ADR-U-0036` | 0036 | lintas-modul | Medan kosong tidak hadir di dokumen JSON, dan ketidakhadirannya  | `docs\adr\0036-medan-kosong-absen-dari-dokumen-bukan-kegagalan.md` |
| `ADR-U-0037` | 0037 | lintas-modul | Kolom dipilih dari medan yang benar-benar diisi, bukan dari yang | `docs\adr\0037-kolom-dipilih-dari-medan-yang-diisi-bukan-yang-tampil.md` |
| `ADR-U-0038` | 0038 | lintas-modul | Tangga persetujuan berakhir karena keadaan jenjang, bukan karena | `docs\adr\0038-pemutus-tangga-persetujuan-adalah-keadaan-jenjang.md` |
| `ADR-U-0039` | 0039 | lintas-modul | Rujuk atau salin: identitas dirujuk, jabatan disalin, data kutip | `docs\adr\0039-rujuk-atau-salin.md` |
| `ADR-U-0040` | 0040 | lintas-modul | Pengaju dikeluarkan dari jenjang pertama - dan lubang di jenjang | `docs\adr\0040-larangan-menyetujui-pekerjaan-sendiri.md` |
| `ADR-U-0041` | 0041 | lintas-modul | Endorsemen Life berbagi tabel dengan new business dan dibedakan  | `docs\adr\0041-endorsemen-berbagi-tabel-dengan-new-business.md` |
| `ADR-U-0042` | 0042 | lintas-modul | Skema relasional dirancang baru bila ia tidak pernah ada | `docs\adr\0042-skema-dirancang-baru-bila-tidak-dapat-direkayasa-balik.md` |
| `ADR-D-CNP-0001` | 0001 | claim-non-prop | Aggregate root adalah Klaim, satu klaim satu kejadian kerugian | `dastin\_migration-docs\claim-non-prop\docs\adr\0001-satu-klaim-satu-kejadian.md` |
| `ADR-D-CNP-0002` | 0002 | claim-non-prop | Klaim yang sudah ditutup tidak dapat dibuka kembali | `dastin\_migration-docs\claim-non-prop\docs\adr\0002-tanpa-reopen.md` |
| `ADR-D-CNP-0003` | 0003 | claim-non-prop | Presisi tinggi sepanjang rantai perhitungan, pembulatan hanya di | `dastin\_migration-docs\claim-non-prop\docs\adr\0003-presisi-dan-pembulatan.md` |
| `ADR-D-CNP-0004` | 0004 | claim-non-prop | Tambalan per-case tidak ikut dimigrasi; angkanya dipindahkan seb | `dastin\_migration-docs\claim-non-prop\docs\adr\0004-tambalan-per-case-tidak-dimigrasi.md` |
| `ADR-D-CNP-0005` | 0005 | claim-non-prop | Claim dan Komite satu unit cutover, dengan shadow-run sebagai bu | `dastin\_migration-docs\claim-non-prop\docs\adr\0005-strategi-cutover.md` |
| `ADR-D-CNP-0006` | 0006 | claim-non-prop | RBAC sistem baru dirancang dari nol, bukan diwarisi | `dastin\_migration-docs\claim-non-prop\docs\adr\0006-rbac-dirancang-dari-nol.md` |
| `ADR-D-CNP-0007` | 0007 | claim-non-prop | Setiap nilai uang disimpan berpasangan, dan ambang kewenangan di | `dastin\_migration-docs\claim-non-prop\docs\adr\0007-nilai-uang-berpasangan-dan-ambang-idr.md` |
| `ADR-D-CNP-0008` | 0008 | claim-non-prop | Nilai hasil suntingan manual bertahan terhadap hitung ulang, dan | `dastin\_migration-docs\claim-non-prop\docs\adr\0008-suntingan-manual-bertahan-dan-terlihat.md` |
| `ADR-D-CNP-0009` | 0009 | claim-non-prop | Keputusan alur disimpan sebagai field terstruktur; komentar tida | `dastin\_migration-docs\claim-non-prop\docs\adr\0009-keputusan-terstruktur-bukan-teks-bebas.md` |
| `ADR-D-CNP-0010` | 0010 | claim-non-prop | Layer dan Retensi Cedant adalah dua entitas terpisah | `dastin\_migration-docs\claim-non-prop\docs\adr\0010-layer-dan-retensi-cedant-dua-tabel.md` |
| `ADR-D-CNP-0011` | 0011 | claim-non-prop | Menjalankan ulang perhitungan alokasi harus menghasilkan keadaan | `dastin\_migration-docs\claim-non-prop\docs\adr\0011-hitung-ulang-idempoten.md` |
| `ADR-D-CNP-0012` | 0012 | claim-non-prop | Tambalan baru dideteksi lewat sapuan ulang, bukan lewat register | `dastin\_migration-docs\claim-non-prop\docs\adr\0012-deteksi-tambalan-baru-lewat-sapuan-bukan-register.md` |
| `ADR-D-CNP-0013` | 0013 | claim-non-prop | Perhitungan dijalankan sebagai turunan dari data, bukan sebagai  | `dastin\_migration-docs\claim-non-prop\docs\adr\0013-perhitungan-sebagai-turunan-data.md` |
| `ADR-D-CNP-0014` | 0014 | claim-non-prop | Aturan penguraian teks ke angka, dan perlakuan atas kurs yang ti | `dastin\_migration-docs\claim-non-prop\docs\adr\0014-penguraian-teks-ke-angka-dan-kurs-kosong.md` |
| `ADR-D-CNP-0015` | 0015 | claim-non-prop | Akseptasi adalah satu entitas dengan keadaan, bukan dua tabel | `dastin\_migration-docs\claim-non-prop\docs\adr\0015-satu-entitas-akseptasi-dengan-keadaan.md` |
| `ADR-D-CNP-0016` | 0016 | claim-non-prop | Tanpa database link ke sistem lain — dan tanpa data pegawai sama | `dastin\_migration-docs\claim-non-prop\docs\adr\0016-tanpa-database-link-ke-sistem-lain.md` |
| `ADR-D-CNP-0017` | 0017 | claim-non-prop | Data aplikasi hanya diubah lewat aplikasi | `dastin\_migration-docs\claim-non-prop\docs\adr\0017-satu-pintu-tulis.md` |
| `ADR-D-CNP-0018` | 0018 | claim-non-prop | Klaim yang terdampak penjaga tanggal yang salah dimigrasi apa ad | `dastin\_migration-docs\claim-non-prop\docs\adr\0018-klaim-terdampak-penjaga-tanggal-dimigrasi-apa-adanya.md` |
| `ADR-D-CNP-0019` | 0019 | claim-non-prop | NULL bukan nol, dan keadaan perhitungan disimpan di tingkat bari | `dastin\_migration-docs\claim-non-prop\docs\adr\0019-null-bukan-nol-dan-keadaan-di-tingkat-baris.md` |
| `ADR-D-CNP-0020` | 0020 | claim-non-prop | Fakultatif dan Treaty dua entitas berbeda; sistem ini hanya memi | `dastin\_migration-docs\claim-non-prop\docs\adr\0020-hanya-treaty-yang-dimiliki.md` |
| `ADR-D-CNP-0021` | 0021 | claim-non-prop | Pengenal polis dan klaim diberi nama menurut isinya; CASEID dipe | `dastin\_migration-docs\claim-non-prop\docs\adr\0021-penamaan-pengenal-polis-dan-klaim.md` |
| `ADR-D-CNP-0022` | 0022 | claim-non-prop | Masa berlaku treaty disimpan sebagai tanggal, dan batasnya inklu | `dastin\_migration-docs\claim-non-prop\docs\adr\0022-masa-berlaku-treaty-sebagai-tanggal-dan-batas-inklusif.md` |
| `ADR-D-CNP-0023` | 0023 | claim-non-prop | Data akseptasi disimpan dalam satu bentuk kanonik; sistem hilir  | `dastin\_migration-docs\claim-non-prop\docs\adr\0023-satu-bentuk-kanonik-tanpa-salinan.md` |
| `ADR-D-CNP-0024` | 0024 | claim-non-prop | Satu akseptasi per klaim per layer per mata uang | `dastin\_migration-docs\claim-non-prop\docs\adr\0024-kunci-alami-akseptasi.md` |
| `ADR-D-CNP-0025` | 0025 | claim-non-prop | Tanggal tutup buku disimpan bertanggal berlaku, bukan satu baris | `dastin\_migration-docs\claim-non-prop\docs\adr\0025-tutup-buku-bertanggal-berlaku.md` |
| `ADR-D-CNP-0026` | 0026 | claim-non-prop | Batas kepemilikan mengikuti nama class: -Work- dan -Data- dimili | `dastin\_migration-docs\claim-non-prop\docs\adr\0026-batas-kepemilikan-mengikuti-nama-class.md` |
| `ADR-D-CNP-0027` | 0027 | claim-non-prop | Dokumen klaim dirujuk, tidak disimpan ulang | `dastin\_migration-docs\claim-non-prop\docs\adr\0027-dokumen-dirujuk-bukan-dimiliki.md` |
| `ADR-D-CNP-0028` | 0028 | claim-non-prop | Basis data tujuan Oracle, skema baru bersebelahan dengan POOLDAT | `dastin\_migration-docs\claim-non-prop\docs\adr\0028-basis-data-tujuan-dan-letak-skema.md` |
| `ADR-D-CNP-0029` | 0029 | claim-non-prop | Kurs yang dipakai ikut tersimpan bersama nilai yang dikonversiny | `dastin\_migration-docs\claim-non-prop\docs\adr\0029-kurs-yang-dipakai-ikut-tersimpan.md` |
| `ADR-D-KCNP-0030` | 0030 | komite-claim-non-prop | Satu penentu jenjang aktif | `dastin\_migration-docs\komite-claim-non-prop\docs\adr\0030-satu-penentu-jenjang-aktif.md` |
| `ADR-D-KCNP-0031` | 0031 | komite-claim-non-prop | Pembuatan sirkulasi dan akibatnya adalah satu transak | `dastin\_migration-docs\komite-claim-non-prop\docs\adr\0031-pembuatan-sirkulasi-satu-transaksi.md` |
| `ADR-D-KCNP-0032` | 0032 | komite-claim-non-prop | Komite adalah modul di dalam konteks Klaim, bukan kon | `dastin\_migration-docs\komite-claim-non-prop\docs\adr\0032-komite-modul-di-dalam-konteks-klaim.md` |
| `ADR-D-KCNP-0033` | 0033 | komite-claim-non-prop | Jalur menutup dan menolak klaim | `dastin\_migration-docs\komite-claim-non-prop\docs\adr\0033-jalur-tutup-dan-tolak-klaim.md` |
| `ADR-D-TI-0034` | 0034 | treaty-in | Arsip JSON sistem lama tidak punya jalur baca | `dastin\_migration-docs\treaty-in\docs\adr\0034-arsip-json-tanpa-jalur-baca.md` |
| `ADR-D-TI-0035` | 0035 | treaty-in | Kegagalan tidak pernah disamarkan menjadi nilai | `dastin\_migration-docs\treaty-in\docs\adr\0035-kegagalan-bukan-nilai.md` |
| `ADR-D-TI-0036` | 0036 | treaty-in | Beku saat disetujui, hitung saat dibaca | `dastin\_migration-docs\treaty-in\docs\adr\0036-beku-saat-disetujui-hitung-saat-dibaca.md` |
| `ADR-D-TI-0037` | 0037 | treaty-in | Masukan versus turunan; tidak ada "mode manual" | `dastin\_migration-docs\treaty-in\docs\adr\0037-masukan-versus-turunan.md` |
| `ADR-D-TI-0038` | 0038 | treaty-in | Aturan yang bisa berubah disimpan sebagai data | `dastin\_migration-docs\treaty-in\docs\adr\0038-aturan-yang-bisa-berubah-disimpan-sebagai-data.md` |
| `ADR-D-TI-0039` | 0039 | treaty-in | Perluasan ADR-D-CNP-0007: tingkat pencatatan ikut tersimpan | `dastin\_migration-docs\treaty-in\docs\adr\0039-perluasan-adr-0007-tingkat-pencatatan.md` |
| `ADR-D-TI-0040` | 0040 | treaty-in | Identitas kontrak, addendum, dan kunci alami yang memperingatkan | `dastin\_migration-docs\treaty-in\docs\adr\0040-identitas-kontrak-dan-addendum.md` |
| `ADR-D-TI-0041` | 0041 | treaty-in | Satu fakta, satu penulis, selalu | `dastin\_migration-docs\treaty-in\docs\adr\0041-satu-fakta-satu-penulis.md` |
| `ADR-D-TI-0042` | 0042 | treaty-in | Sejarah pindah apa adanya: baris warisan dan aturan sentuh-perba | `dastin\_migration-docs\treaty-in\docs\adr\0042-sejarah-pindah-apa-adanya.md` |
| `ADR-D-TI-0043` | 0043 | treaty-in | Migrasi memindahkan; menghitung ulang adalah peristiwa bisnis te | `dastin\_migration-docs\treaty-in\docs\adr\0043-migrasi-memindahkan-hitung-ulang-peristiwa-bisnis.md` |
| `ADR-D-TI-0044` | 0044 | treaty-in | Wewenang: keadaan dan peran, terpisah tegas | `dastin\_migration-docs\treaty-in\docs\adr\0044-wewenang-keadaan-dan-peran.md` |
| `ADR-D-TI-0045` | 0045 | treaty-in | Jejak perubahan adalah fakta mesin, terpisah dari catatan manusi | `dastin\_migration-docs\treaty-in\docs\adr\0045-jejak-perubahan-sebagai-fakta-mesin.md` |
| `ADR-D-TI-0046` | 0046 | treaty-in | Satu keadaan siklus hidup kontrak | `dastin\_migration-docs\treaty-in\docs\adr\0046-satu-keadaan-siklus-hidup.md` |
| `ADR-D-TI-0047` | 0047 | treaty-in | Data uji: identitas boleh disamarkan, angka tidak | `dastin\_migration-docs\treaty-in\docs\adr\0047-data-uji.md` |
| `ADR-D-TI-0048` | 0048 | treaty-in | Sambungan (seam) ke Treaty In Adjustment | `dastin\_migration-docs\treaty-in\docs\adr\0048-seam-treaty-in-adjustment.md` |
| `ADR-D-TI-0049` | 0049 | treaty-in | Jenis addendum adalah satu sumbu; materialitas diturunkan | `dastin\_migration-docs\treaty-in\docs\adr\0049-jenis-addendum-satu-sumbu.md` |
| `ADR-D-TI-0050` | 0050 | treaty-in | Acuan kapasitas adalah sumber luar yang tidak dipercaya | `dastin\_migration-docs\treaty-in\docs\adr\0050-acuan-kapasitas-sumber-luar-tak-dipercaya.md` |
| `ADR-D-TI-0051` | 0051 | treaty-in | Anggap ada setidaknya satu konsumen hilir yang tidak dikenal | `dastin\_migration-docs\treaty-in\docs\adr\0051-anggap-ada-konsumen-hilir.md` |
| `ADR-D-TI-0052` | 0052 | treaty-in | Satu rantai persetujuan empat tingkat, tabel perutean kosong dar | `dastin\_migration-docs\treaty-in\docs\adr\0052-satu-rantai-persetujuan-tanpa-pengecualian.md` |
| `ADR-D-TI-0053` | 0053 | treaty-in | Mata uang sebagai daftar; `CurrencyRelation` disimpan tetapi tid | `dastin\_migration-docs\treaty-in\docs\adr\0053-mata-uang-daftar-dan-currencyrelation-pasif.md` |
| `ADR-D-TI-0054` | 0054 | treaty-in | Keadaan warisan yang tidak ada padanannya | `dastin\_migration-docs\treaty-in\docs\adr\0054-keadaan-warisan-tak-terpetakan.md` |
| `ADR-D-TI-0055` | 0055 | treaty-in | Daftar keadaan siklus hidup dan perpindahan yang sah | `dastin\_migration-docs\treaty-in\docs\adr\0055-daftar-keadaan-dan-perpindahan.md` |
| `ADR-D-TI-0056` | 0056 | treaty-in | Tidak ada stored procedure di sistem baru | `dastin\_migration-docs\treaty-in\docs\adr\0056-tanpa-stored-procedure.md` |
| `ADR-F-0001` | 0001 | Facultative | Rekonsiliasi paralel run eksak, dicapai bertahap | `jefri\OUTPUT FIX\adr\0001-rekonsiliasi-eksak-bertahap.md` |
| `ADR-F-0002` | 0002 | Facultative | Satu flag fase: `panic` saat paralel run, `decline`+log saat pro | `jefri\OUTPUT FIX\adr\0002-flag-fase-panic-vs-decline.md` |
| `ADR-F-0003` | 0003 | Facultative | Mesin akseptasi ditulis lebih dulu; fixture tabel limit menjadi  | `jefri\OUTPUT FIX\adr\0003-mesin-akseptasi-fixture-sebagai-kontrak.md` |
| `ADR-F-0004` | 0004 | Facultative | `Money` dan `Ratio` adalah dua tipe yang tidak dapat dijumlahkan | `jefri\OUTPUT FIX\adr\0004-money-dan-ratio-tipe-terpisah.md` |
| `ADR-F-0005` | 0005 | Facultative | Presisi pembulatan ditulis literal di tiap langkah, bukan disent | `jefri\OUTPUT FIX\adr\0005-presisi-pembulatan-literal-per-langkah.md` |
| `ADR-F-0006` | 0006 | Facultative | Mata uang yang tidak diketahui adalah keadaan eksplisit, bukan d | `jefri\OUTPUT FIX\adr\0006-mata-uang-unknown-eksplisit.md` |
| `ADR-F-0007` | 0007 | Facultative | Total yang menjumlahkan lintas mata uang bukan `Money`, melainka | `jefri\OUTPUT FIX\adr\0007-agregat-lintas-mata-uang-bukan-money.md` |

---

# Lampiran 2 - Pasangan sebidang

Keputusan dari seri berbeda yang membahas hal yang sama, disandingkan.
[AWAS] **Tidak ada yang dimenangkan di sini.** Penentuannya milik komite pengarah -
lihat dokumen Steering bagian 5.

## Kelompok 1 - Representasi uang

| Kode | Judul | Yang ditetapkannya |
| --- | --- | --- |
| `ADR-U-0003` | Uang tidak direpresentasikan sebagai `float` | delapan kolom nilai pada satu modul adalah uang, dan tidak boleh `float` |
| `ADR-U-0016` | Kolom uang memakai `NUMBER(38,8)` | **seluruh sistem** - tiga puluh digit di depan koma, delapan di belakang |
| `ADR-D-CNP-0003` | Presisi tinggi sepanjang rantai, pembulatan hanya di tepi | presisi dijaga sepanjang perhitungan; pembulatan hanya pada titik keluar |
| `ADR-D-CNP-0007` | Nilai uang disimpan berpasangan; ambang dibandingkan terhadap nilai IDR | uang selalu membawa mata uangnya; ambang kewenangan dibandingkan setelah dikonversi |
| `ADR-F-0004` | `Money` dan `Ratio` dua tipe yang tidak dapat dijumlahkan | uang dan rasio adalah tipe terpisah; skala melekat pada nilainya |
| `ADR-F-0005` | Presisi pembulatan ditulis literal di tiap langkah | presisi tidak disentralkan; tiap langkah menuliskan presisinya sendiri beserta asalnya |
| `ADR-F-0007` | Total lintas mata uang bukan `Money` | lima total skalar memakai tipe ketiga yang tidak dapat diaritmetikakan |

**Persamaan.** Ketujuhnya menolak bilangan mengambang untuk uang, dan ketujuhnya
memperlakukan kehilangan digit sebagai cacat, bukan sebagai penyederhanaan yang wajar.

**Perbedaan.** Seri U menetapkan **satu bentuk kolom untuk seluruh sistem**
*(`NUMBER(38,8)`)*; seri D menetapkan **aturan perilaku** - presisi tinggi sepanjang rantai,
pembulatan hanya di tepi - tanpa mengikat bentuk kolomnya; dan seri F menetapkan **sistem
tipe di dalam kode** yang memisahkan uang dari rasio dan menuliskan presisi per langkah,
bukan satu aturan global.

[!] **Dan satu selisih angka yang belum dipertemukan:** `ADR-U-0016` menetapkan
`NUMBER(38,8)`, sementara spec penyimpanan Treaty menetapkan `NUMBER(20,8)`. Keduanya
berskala delapan; yang berbeda **lebar digit di depan koma**.

## Kelompok 2 - Tabel akar bersama

| Kode | Judul | Yang ditetapkannya |
| --- | --- | --- |
| `ADR-U-0025` | Tabel inti berbagi kunci utama dengan tabel kerja | tidak ada kolom penyambung; keduanya memakai kunci utama yang sama persis |
| `ADR-U-0026` | Tabel kerja adalah akar lintas-lini | satu tabel kerja menjadi akar bagi lebih dari satu lini pekerjaan |
| `K-064` `K-065` | keputusan rangkaian Fakultatif | tabel akar bersama **diterima**, tetapi kolomnya dibangun **hanya dari sisi Fakultatif**; irisan kolom dan tabrakan ruang kunci diterima sebagai **risiko sadar** |
| `K-069` | keputusan rangkaian Fakultatif | kolom penanda lini **dinyatakan di luar proyek** dan **dicabut** dari rancangan |

**Persamaan.** Keduanya sepakat bahwa **satu tabel kerja menjadi akar bersama** untuk
lebih dari satu lini, dan bahwa tabel inti di bawahnya tidak memakai kolom penyambung.

**Perbedaan.** `ADR-U-0026` memperlakukan sifat lintas-lini itu sebagai **kewajiban yang
menuntut pembeda**, sedangkan `K-069` **mencabut kolom pembedanya** dan menyatakannya di
luar proyek - sehingga satu tabel akan memuat baris dari dua lini **tanpa kolom yang
membedakannya**. `K-065` menyatakan rekonsiliasi rancangan kolom antara kedua sisi
**gugur**, sementara seri U tidak pernah mencabut kewajiban itu.

## Kelompok 3 - Penamaan tabel untuk properti yang sama

Properti yang sama di sistem lama mendapat **dua nama tabel** di dua rangkaian kerja:

| Properti sistem lama | Rangkaian Treaty | Rangkaian Fakultatif |
| --- | --- | --- |
| daftar ceding | `T_POLIS_CEDING` | `T_CEDINGCOLIST` |
| daftar angsuran | `T_POLIS_INSTALMENT` | `T_LISTINSTALLMENT` |
| daftar sebaran | `T_POLIS_SPREADING` | `T_SPREADINGLIST` |
| data kuotasi | `T_POLIS_QUOTATION` | `T_QUOTATIONDATA` |

**Persamaan.** Keduanya memecah daftar bersarang sistem lama menjadi tabel anak
tersendiri, dengan kunci ke tabel induk yang sama.

**Perbedaan.** Rangkaian Treaty menamai tabel menurut **perannya di dalam polis**
*(`T_POLIS_<peran>`)*; rangkaian Fakultatif menamainya menurut **nama properti Pega asalnya**
*(`T_<NAMAPROPERTI>`)*. [!] Bila keduanya menulis ke basis data yang sama, akan ada **dua
tabel untuk satu hal**.

---

## Catatan tentang bentuk

Ketiga rangkaian memakai bentuk kepala dokumen yang berbeda - sebagian memakai blok
metadata di atas judul, sebagian menuliskan status sebagai baris tebal. Di dokumen ini
tampilannya diseragamkan; **berkas sumbernya tidak diubah**.

*Disusun 25 September 2026 dari berkas keputusan arsitektur proyek migrasi Nusantara Re.*