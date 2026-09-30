---
status: accepted
tanggal: 2026-09-14
sumber: grilling Ronde 1 Q1 + Ronde 2 Q8/Q9/Q15 (`.scratch/claim-life/grilling-ronde-1.md`, `grilling-ronde-2.md`), keputusan work owner
---

# Claim — Life dan Komite Life adalah dua konteks terpisah, dihubungkan kontrak child work

**Claim — Life** dispesifikasikan sebagai bounded context tersendiri; **Komite Life** diperlakukan
sebagai **konteks/sistem luar**. Isi dan alur internal Komite **tidak** dispesifikasikan di sini.
Hubungan keduanya berupa **tiga kontrak eksplisit**: penyerahan kasus sebagai *child work*, jalur
balik hasil keputusan, dan tabel akseptasi bersama.

## Considered Options

- **(a) Claim Life saja, Komite sebagai konteks luar** — dipilih
- (b) Claim Life + Komite Claim Life sebagai satu konteks
- (c) Claim Life saja, tetapi rule tulis bersama ikut dispesifikasikan di dalamnya

(b) ditolak karena Komite membawa **13 OQ pemblokir** yang hampir seluruhnya RBAC
(`discovery/D3-D4-CLOSING-REPORT.md` §4.2); menariknya masuk akan mengubah konteks tersiap
(5 pemblokir, cakupan bukti `full`) menjadi yang paling terblokir.

## Kontrak 1 — penyerahan kasus ke Komite `[terverifikasi]`

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

### Pemilihan roster komite `[terverifikasi]`

Langkah 2 memanggil `pxRetrieveReportData` dengan report
**`ASM-FW-GCNMFW-INT-EMAILKOMITE!FILTEREMAILKOMITEWITHLIMIT`**
(`Claim Life/ReportDefinition/FilterEmailKomiteWithLimit.xml`, `RULE-OBJ-REPORT-DEFINITION`,
67.827 byte), lewat `Param.pyReportName` / `Param.pyReportClass`.

Report itu menerima **`Param.LIMIT_BOTTOM`** dan **`Param.STS_KLAIM`** — artinya **roster dipilih
berdasarkan pita nilai**, pola yang sama dengan `Komite Claim FacIn/Activity/ApprovalKomite_Act.xml`
(OQ-037). **Ambang dan aturannya sendiri belum terverifikasi.**

### Muatan yang menyeberang `[terverifikasi]`

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

### Tambahan pada muatan — keputusan Ronde 2 Q15

`[terverifikasi work owner 2026-09-14]` Daftar di atas **benar** sebagai muatan sekarang, tetapi
**tidak cukup**. Kontrak baru **menambahkan tiga hal**:

| Tambahan | Alasan |
| --- | --- |
| **nilai klaim** (`CLAIM_AMOUNT`) | Komite memilih roster berdasarkan pita nilai (`Param.LIMIT_BOTTOM`); nilainya sekarang tidak ikut menyeberang |
| **`CURRENCY`** | nilai uang tanpa mata uang tidak dapat dibandingkan terhadap pita nilai — lihat **ADR-0003** |
| **`STS_REJECT`** saat penyerahan | keadaan klaim pada saat diserahkan, agar jalur balik punya titik awal yang eksplisit |

Ini **penyimpangan sadar** dari paritas, sejenis dengan **ADR-0007**: alurnya tidak berubah, yang
bertambah adalah apa yang terekam menyeberang.

`[terverifikasi work owner 2026-09-14]` **`KomiteLoop` ditentukan oleh Claim — Life**, bukan oleh
Komite — sesuai arah muatan di atas (`childPageKomite.KomiteLoop` diisi induk). **Apa yang
menentukan nilainya masih OQ-032.**

## Kontrak 2 — jalur balik dari Komite `[terverifikasi]`

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

## Kontrak 3 — tabel akseptasi bersama `[terverifikasi]`

`OS_AKSEPTASI_KLAIM_LIFE` ditulis oleh **satu rule yang sama** dari kedua sisi:
`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` /
`RULE-CONNECT-SQL`, hash ternormalisasi **`c50bfd9a12`** identik di
`Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` dan
`Komite Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml`, dan **tidak terdaftar** di register
konflik `discovery/inventory/_oq011-konflik-isi.md`.

Jadi **bukan dua penulis independen** — satu rule dipanggil dua sisi. 55 kolomnya terbaca di
`discovery/flows/Komite Claim Life.md` §3.1.

## Consequences

- Spesifikasi Claim — Life berhenti di batas: membuat child work, menerima hasil keputusan, dan
  menulis tabel akseptasi.
- `AcceptStatus` adalah **kosakata konteks Komite**, bukan status internal Claim — Life. Ia dipetakan
  ke `STS_REJECT` **di batas kontrak**, bukan disimpan sebagai status kedua (Ronde 2 Q8).
- Perubahan pada struktur `KomiteList`, pada muatan penyerahan, atau pada
  `OS_AKSEPTASI_KLAIM_LIFE` adalah **perubahan kontrak lintas konteks**, bukan perubahan internal.
- Nilai uang yang menyeberang harus memakai representasi yang sama di kedua sisi (**ADR-0003**).
- Setiap transisi lewat jalur balik ini wajib merekam siapa + kapan (**ADR-0007**).
- `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` adalah **satu-satunya salinan di
  korpus** dan dipanggil Komite Claim Life (**OQ-035**) — arah ketergantungan ini berlawanan dengan
  kontrak di atas dan perlu ditetapkan saat Komite dispesifikasikan.

## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-039** (dipersempit 2026-09-14) | **Kapan** penyerahan dipicu. Akibat reject Admin sudah terjawab Ronde 4; `SaveAdjustment_Act` dinyatakan dead. Sisa: pemicu Adjustment/Close/Reject di 3 modul Claim lain. Pemilik: Product+UW |
| **OQ-035** | Activity dipanggil lintas modul tetapi salinannya hanya ada di satu modul |

**Tertutup untuk Claim — Life sejak 2026-09-14** — tidak lagi menyentuh ADR ini:

| OQ | Verdict |
| --- | --- |
| **OQ-032** | **`KomiteLoop` = COUNT roster `EMAILKOMITE` aktif ber-`LIMIT_BOTTOM <= CLAIM_AMOUNT`** — data-driven, tanpa konstanta. Rinci di **ADR-0012** |
| **OQ-037** | Ambang adalah **data di tabel** `EMAILKOMITE`, bukan hardcode |
| **OQ-060** | Nilai uang menyeberang sebagai `(amount, currency)` per baris, invariant satu klaim satu mata uang (**ADR-0003**) |
| **OQ-013** | `COMMIT` ada di **pembungkus Pega**, bukan di procedure. Sistem baru: **Go memegang batas transaksi** |
| **OQ-061** | Unit keputusan = baris `AdjustmentList` (**ADR-0011**); `PremiumListDetail` dan header klaim adalah cerminan |
