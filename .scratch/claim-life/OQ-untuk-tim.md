# Pertanyaan untuk Tim — pembuka tiket Claim — Life

Dokumen bantu untuk mengumpulkan jawaban dari DBA, Product+UW, Finance, dan IT-infra.
Tujuannya membuka 6 tiket berstatus `needs-info` menjadi `ready-for-agent`.

**Cara pakai:** bawa ke tiap pemilik peran, isi kolom "Jawaban", lalu jawaban ini akan dicatat ke
CONTEXT.md/ADR/register lewat skill (bukan diketik langsung ke artefak). Yang tak terjawab tetap
OQ terbuka — jangan ditebak.

Sumber pertanyaan: `discovery/open-questions.md` (nomor OQ otoritatif di sana). Tiket yang
diblokir disebut di tiap butir.

---

## Untuk DBA

### OQ-002 — Kontrak `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` (memblokir tiket 02, 12)
Penomoran klaim Life memanggil procedure ini; badannya tidak ada di korpus.
Yang dibutuhkan:
1. Format nomor klaim Life yang dihasilkan (retro dan non-retro) — bentuk string persisnya.
2. Apa yang membedakan nomor retro vs non-retro.
3. Sequence di-reset per tahun, per lini bisnis, atau berjalan terus (global)?
4. (untuk tiket 12) Kontrak `GET_TOKEN_STORAGE` — parameter masuk/keluar, dan apakah commit sendiri.

> **Jawaban:**

### OQ-001 — DDL Oracle produksi (memblokir tiket 13; TIDAK memblokir tiket 01)
Tidak ada DDL di korpus. Untuk migrasi data nyata (tiket 13) dibutuhkan definisi tabel:
`OS_AKSEPTASI_KLAIM_LIFE`, `JSON_KLAIM`, `M_LIFE_PREMIUM_DETAIL`, dan tabel klaim Life terkait —
nama kolom, tipe, presisi, nullability, PK/FK, index.
Catatan: tiket 01 (skema uji provisional) TIDAK menunggu ini — nama kolom sudah terbaca dari SQL,
tipe longgar. Yang menunggu hanya migrasi produksi.

> **Jawaban:**

### OQ-013 — Batas transaksi & identitas pengguna di dalam stored procedure (memblokir tiket 13)
Procedure `POOLDATA.*` melakukan `COMMIT` di dalam blok PL/SQL, dan `{OperatorID.pyUserIdentifier}`
dikirim sebagai parameter. Yang dibutuhkan:
1. Apakah semua procedure `POOLDATA.*` commit sendiri?
2. Identitas pengguna dipakai untuk apa di dalam procedure — audit trail, otorisasi, atau keduanya?

> **Jawaban:**

### OQ-018 — pxHostId production (memblokir tiket 13 untuk cutover; sebagian sudah dijawab)
Sudah dijawab untuk Claim Life: jboss1073 = production, jboss117 = dev (mirroring).
Sisa untuk cutover: konfirmasi pxHostId `pega-nusre` dan satu id-hash — lingkungan apa?

> **Jawaban:**

### OQ-047 — Daftar endpoint di tabel `M_LINK_SERVICE` (memblokir tiket 12)
Alamat integrasi (Google Storage, email, Arasapas, konversi) dibaca dari tabel `M_LINK_SERVICE`,
bukan hanya SystemSettings. Yang dibutuhkan: daftar endpoint aktual (untuk jadi env var), atau
konfirmasi bahwa semuanya via tabel itu dan strukturnya.

> **Jawaban:**

---

## Untuk Product + Underwriting

### OQ-032 — Sumber nilai `KomiteLoop` (jumlah tingkat tangga komite) (memblokir tiket 10)
`KomiteLoop` menentukan berapa tingkat persetujuan komite. Nilainya ditentukan Claim Life, tapi
sumber/aturan penentuannya tidak ada di korpus.
Pertanyaan: apa yang menentukan jumlah tingkat komite untuk sebuah klaim? (nilai klaim, jenis, dll)

> **Jawaban:**

### OQ-037 — Ambang nominal roster komite (memblokir tiket 10; Finance + Product+UW)
Komposisi roster komite ditentukan ambang nominal ter-hardcode (mis. batas 30jt/50jt di modul lain).
Pertanyaan: apa aturan ambang nominal yang menentukan siapa masuk roster komite untuk klaim Life?
Apakah mata uangnya IDR? Apakah dapat dikonfigurasi?

> **Jawaban:**

### Tinjauan aturan turunan "klaim selesai" (memblokir tiket 04)
Spec §3 menetapkan (bukan temuan korpus, tapi konsekuensi logis jawaban Anda):
> Klaim SELESAI bila tidak ada lagi baris AdjustmentList bernilai 0 DAN ada minimal satu baris
> bernilai 1 (aksep). Bila tak ada 0 dan tak ada 1 → ditolak seluruhnya, tapi masih bisa
> dilanjutkan dengan baris baru.
Pertanyaan: apakah aturan ini benar menurut proses bisnis Anda?

> **Jawaban:**

### Satuan pergeseran tanggal DOL (memblokir tiket 06)
Di ValidasiDOL_Act, cabang TP/TR memakai argumen `@addCalendar(...,1,...)` — [dugaan] +1 hari.
Definisi fungsi tidak ada di korpus.
Pertanyaan: untuk tipe TP/TR, apakah jendela validasi Date of Loss digeser +1 hari? Bila ya,
mengapa (aturan bisnis) — atau ini kekeliruan lama?

> **Jawaban:**

### OQ-060 — Cakupan kolom CURRENCY pada rekam akseptasi Life (memblokir bentuk tipe uang, tiket terkait)
CURRENCY ada di tingkat baris dan disalin dari baris 1. Pertanyaan: apakah satu klaim Life selalu
satu mata uang, atau boleh campur antar baris?

> **Jawaban:**

---

## Ringkasan pemetaan OQ → tiket

| OQ | Pemilik | Memblokir tiket |
| --- | --- | --- |
| OQ-002 | DBA | 02, 12 |
| OQ-001 | DBA | 13 (bukan 01) |
| OQ-013 | DBA | 13 |
| OQ-018 | IT-infra/DBA | 13 (cutover) |
| OQ-047 | DBA/Platform | 12 |
| OQ-032 | Product+UW | 10 |
| OQ-037 | Finance + Product+UW | 10 |
| aturan "klaim selesai" | Product+UW (work owner) | 04 |
| satuan DOL | Product+UW | 06 |
| OQ-060 | Product+UW + DBA | bentuk tipe uang |

**Setelah terjawab:** jawaban dibawa kembali ke sesi Claude, dicatat ke CONTEXT.md/ADR/register
lewat alur skill, lalu tiket yang bersangkutan dinaikkan `needs-info` → `ready-for-agent`.
Yang tetap tak terjawab: tiket tetap `needs-info`, jangan ditebak.
```

## 27 September 2026 — tiga pertanyaan untuk pemilik ekspor / pengembang Pega

**OQ-A.** Apakah ada **harness portal/navigasi** untuk Claim Life *(padanan `SFAPortalOpportunities`
di NB Treaty In)* beserta tombol Create-nya? Bila ada, mohon diekspor. Sampai terjawab, menu sidebar
Claim Life adalah **minimum berbukti**: `Inbox Claim Life` *(label `[tidak ada di korpus]`)* dan
`Register` *(`Register_Flow.xml:155`)*.

**OQ-B.** `Activity\setDetailClaim_act.xml` tampak **residu uji pengembang**: `Obj-Open-By-Handle`
pada satu handle literal *(b284)*, `Property-Set` tanggal literal Januari–Februari 2026 *(b478,
b711)*, prasyarat yang membandingkan `.NAME_OF_INSURED` dengan satu nama literal *(b872)*, lalu
`Obj-Save` *(b960)*. Ia terpasang pada tombol `Choose` popup pencarian polis
*(`SearchPolicy_Section.xml` 3337–3356, 3468)*. **Mohon konfirmasi** apakah ia memang tidak
dipakai di produksi; kami **tidak menirunya**.

**OQ-C.** `SendtoAdmin_Act` dan `SendtoAdmin_Act1` keduanya berprasyarat
`pyWorkPage.pyPosition=="ReasLifeMedicalAdvisor"` *(WhenTrue=2, WhenFalse=3 = lewati)*, padahal
tombol pemanggilnya — `Send Back to Register` dan `Send to Medical Check` — berdiri di layar
**Outstanding**, yang `pyPosition`-nya `ReasLifeAdmin`. Akibatnya menurut XML apa adanya: kedua
tombol itu **tidak menulis apa pun** pada posisi Admin. **Mohon konfirmasi** apakah prasyaratnya
memang demikian di produksi, atau salah tempel. Keputusan kami *(butir aw)* mengikuti **maksud**
yang terang dari label tombol dan penyambung alurnya, dan cacatnya dilaporkan di sini.

**OQ-D** *(urutan modul, untuk work owner)*. Sebelas medan layar Register terikat
`.PolicyDataLife.*` dan terisi dari kasus **PremiumList Life**. Register Claim Life karena itu
**menunggu modul PremiumList Life** untuk lengkap. Mohon konfirmasi urutan pengerjaan modul.

## 27 September 2026 — dua pertanyaan tambahan (paket Register(2))

**OQ-E** *(cacat halus, untuk pengembang Pega)*. `RDBList/GetPesertaClaim_sql1.xml:85` memasang
**kedua** `LIKE` tanpa syarat:

> `AND CERTIFICATE_NO LIKE '%'||{SearchPolicyHolder.CARI2}||'%'`
> `AND UPPER(NAME_OF_INSURED) LIKE '%'||{SearchPolicyHolder.CARI3}||'%'`

Di Oracle, `X LIKE '%'` bernilai **FALSE** ketika `X` NULL. Akibatnya kotak pencarian `Find
Insured` yang dibiarkan **kosong** pun membuang setiap peserta yang `NAME_OF_INSURED`- atau
`CERTIFICATE_NO`-nya NULL — baris yang hilang tanpa seorang pun memintanya, dan tanpa pesan.

**Mohon konfirmasi** apakah itu memang dikehendaki. Kami **tidak menirunya**: penyaring hanya
terpasang untuk kotak yang terisi, sehingga kotak kosong berarti *"jangan saring"* — yang memang
dibaca orang dari kotak kosong.

**OQ-F** *(untuk pemilik ekspor)*. `Activity/UploadCSVClaimLife_Act.xml` hanya **tiga** langkah:
`Page-Remove TempWorkPage` *(b250)*, `Page-New TempWorkPage` *(b340)*, lalu
`Call pxUploadCSVResults` *(b488)* — yaitu mesin unggah **bawaan platform**.

Pemetaan kolom CSV ke medan peserta karena itu **tidak ada di korpus**: ia tersimpan di
konfigurasi gadget, bukan di rule yang diekspor. **Mohon kirimkan** definisi pemetaan kolomnya
bila fitur unggah CSV memang dipakai; tanpa itu fitur ini tidak dapat ditiru tanpa mengarang, dan
kami tidak mengarang.

## 27 September 2026 — OQ-G (paket Detail & Tutup 1)

**OQ-G** *(dua cacat rule, untuk pengembang Pega)*. `ValidasiClaimReceived_Act` dan
`ValidasiSTNC_Act` berbentuk sama persis, dan keduanya memuat dua hal yang tampak keliru:

**G1 — pesan galat disetel tetapi tidak pernah dipasang.** Keduanya menyetel `local.errmsg`
*(`"Max Claim invalid"` b451; `"STNC invalid"` b479)* lalu **tidak punya langkah
`Property-Set-Messages`** — tidak seperti `ValidasiDOL_Act`, yang punya di b842. Menurut XML apa
adanya, pemakai **tidak pernah melihat** sebab penolakannya; yang tampak hanya `.MAXCLAIM_RECEIVED`
atau `.STNC` yang tiba-tiba berisi tanggal.

**Mohon konfirmasi** apakah pesan itu memang tidak dimaksudkan tampil. Kami memakai kalimatnya
**verbatim** supaya teks yang sama dapat dicari di kedua sistem.

**G2 — selisih negatif selalu lolos.** Perbandingannya `@if(selisih <= ambang, "", …)` *(b582,
b627)* tanpa lantai bawah. Klaim yang **diterima sebelum tanggal kejadiannya** menghasilkan
selisih negatif, dan negatif selalu `<=` ambang — sehingga ia terbaca **sah**.

Kami **menirunya apa adanya** dan menjaganya dengan `TestLubangSelisihNegatif`, supaya perilaku itu
terlihat dan tidak berubah diam-diam. **Mohon putuskan** apakah lantai bawah perlu ditambahkan;
bila ya, test itulah yang gagal lebih dulu dan menagih keputusannya.

**Catatan ambang** *(bukan pertanyaan — penjelasan)*: `MAXEXPIREDCLAIM` dan `MAXDATARECEIVE` dibaca
`RDBList/GetProductName.xml:84` dari tabel produk warisan lewat
`pyWorkPage.PolicyDataLife.ProductNameID`. `PolicyDataLife` **menunggu modul PremiumList Life**
*(keputusan av)*, dan tabelnya sendiri berstatus `[data DBA]` OQ-001. Karena itu aturannya ditulis
sebagai fungsi murni yang **menerima** ambang; ambang kosong menjawab **galat**, bukan nol.

⚠️ Perhatikan bedanya nama, sebab ia mudah tertukar: propertinya `.MAXDATARECEIVED` *(ber-D)*,
kolom sumbernya `MAXDATARECEIVE` *(tanpa D)*.

**G3 — nilai antara diparkir di halaman BERSAMA, hanya pada salah satu dari dua rule kembar.**

Kedua rule itu beraturan sama, tetapi menyimpan selisih harinya di tempat yang berbeda:

| Rule | Baris | Tempat selisih disimpan | Sifat |
| --- | ---: | --- | --- |
| `ValidasiClaimReceived_Act` | 561 → dibaca 582 | **`TempDetail.CARI2`** | properti pada halaman `TempDetail` |
| `ValidasiSTNC_Act` | 598 → dibaca 627 | `Local.DateDif` | local sejati |

`TempDetail` **bukan** halaman gores milik rule ini sendiri: `SaveInsuredClaim_Act` memakai halaman
yang sama untuk menumpuk peserta terpilih *(`TempDetail.pxResults(<APPEND>)`, b535 dst)*. Nilai
antara yang diparkir di sana karena itu dapat tertimpa di antara penulisan *(b561)* dan pembacaannya
*(b582)* bila ada rule lain yang menyentuh halaman itu di sela keduanya.

**Mohon konfirmasi** apakah `TempDetail.CARI2` memang dimaksudkan, atau ia sisa penyuntingan
— rule kembarnya memakai local, dan local memang yang tepat untuk nilai antara.

⚠️ Kami **tidak menirunya**: di sisi Go selisihnya peubah lokal dan tidak pernah meninggalkan
fungsinya, sehingga tidak ada halaman bersama yang dapat tertimpa.

⚠️ Ralat susunan kata G2 di atas: yang dibandingkan b582 secara harfiah adalah `TempDetail.CARI2`,
bukan sebuah local bernama "selisih". Isinya tetap selisih hari itu, dan kesimpulan G2 tidak berubah.

## 27 September 2026 — OQ-H (rujukan menggantung pada LIMA total uang)

**OQ-H** *(untuk pemilik ekspor — menahan lima medan uang)*.
`Section/ClaimLifeDetailGCNM.xml` memanggil activity **`CheckTotalAdjustmentClaim`** sepuluh kali,
tetapi activity itu **tidak punya satu pun berkas rule di seluruh korpus**.

Diperiksa dua cara: `find . -iname "*CheckTotalAdjustment*"` → **nol** hasil; dan
`grep -rl` → satu-satunya berkas yang menyebutnya adalah section itu sendiri, dengan **10**
kemunculan.

Kesepuluh pemanggilan itu menempel pada **lima medan total**, dua pemanggilan per medan
*(`postValue` + `refresh`)*:

| Medan | `pyLabelPreview` | Pemanggilan |
| --- | ---: | --- |
| `Total Share Nusantara Re` | b20914 | b20970, b21091 |
| `Total Sum Insured` | b21201 | b21262, b21377 |
| `Total Sum Reasured` | b21488 | b21546, b21664 |
| `Total Share Retro` | b21775 | b21836, b21951 |
| `Total Claim Amount` | b22063 | b22120, b22238 |

⛔ **Akibatnya kami tidak dapat meniru kelima angka itu.** "Total" terdengar seperti penjumlahan
kolomnya, tetapi pertanyaan yang menentukan tidak terjawab oleh apa pun yang kami punya:

- baris yang **mana** yang ikut dijumlahkan — seluruhnya, atau hanya yang `IsCheck`?
- apakah baris ber-`STS_REJECT = 2` *(ditolak)* ikut?
- apakah ia menjumlahkan, atau memeriksa silang terhadap nilai yang sudah tersimpan?

Jalur tolak **mencabut `IsCheck`** *(`RejectOSClaimLife_Act`)*, jadi jawabannya berpengaruh nyata —
dan ini **angka uang**. Menebaknya melanggar ADR-U-0003, dan angka uang yang ditebak jauh lebih
berbahaya daripada angka yang dinyatakan belum ada.

**Mohon kirimkan** `CheckTotalAdjustmentClaim`, atau konfirmasikan bahwa ia memang sudah dihapus
dan kelima medan itu tidak lagi berisi apa-apa di produksi.

**Sementara itu**, layar Detail menampilkan kelima medannya **dengan penanda "belum tersedia"** dan
menyebut nama rule-nya — bukan dihilangkan *(paritas yang tampak lengkap padahal tidak adalah
paritas yang tidak akan dicari lagi)*, dan bukan dijumlahkan sendiri. Ada uji yang akan gagal bila
seseorang menambahkan penjumlahan, dan uji lain yang akan gagal bila ekspornya kelak dilengkapi —
gagalnya yang terakhir itu **kabar baik**: aturannya sudah dapat ditiru.

### Ralat OQ-H dan laporan tutup — 27 September 2026 *(sesudah telaah paritas)*

Dua klaim kami sendiri terbukti keliru dan diralat di sini, bukan didiamkan.

**R1 — "dua pemanggilan per medan (`postValue` + `refresh`)" KELIRU.** Blok `postValue`
*(mis. b20940)* **tidak membawa `pyActivity` sama sekali** — hanya `pyActivityClass`. Kedua
kemunculan per medan adalah aksi **`refresh` yang SAMA**, terserialisasi dua kali oleh Pega:
sekali di `pyModes`, sekali di `pyActionSets`. Jadi yang benar: **lima medan, lima pemanggilan,
sepuluh kemunculan teks**. Kesimpulan OQ-H **tidak berubah** — rule-nya tetap tidak ada, dan
kelima total tetap tidak dapat ditiru.

**R2 — "tombol `Close Claim` tidak punya aksi lain" KELIRU.** Satu klik menjalankan **dua** aksi:
`refresh` → `ProtectCloseClaim_act` *(b1101)* **dan** `closeContainer` *(b1129)*. Ada pula
**konfirmasi** yang sebelumnya terlewat: `Are you sure want to Close Claim?` *(b499, `pyCaption`
b1499)*. Gerbangnya tetap seperti yang dilaporkan; yang salah adalah kalimat "tidak punya aksi
lain".

⚠️ Keduanya lolos karena bacaan pertama berhenti pada **aksi yang membawa activity**, lalu
menyimpulkan tentang **seluruh** tombol. Aksi tanpa activity tetap aksi.

---

### ⛔ Ralat kedua atas OQ-H — 27 September 2026 *(kesimpulannya DICABUT)*

**Fakta OQ-H tetap berdiri:** `CheckTotalAdjustmentClaim` dirujuk **sepuluh kali** oleh
`Section/ClaimLifeDetailGCNM.xml` dan **nol** berkas rule-nya ada di seluruh korpus.

**Kesimpulan yang digantungkan padanya DICABUT.** OQ-H menutup dengan *"kami tidak dapat meniru
kelima angka itu; baris yang mana yang ikut belum terjawab"*, dan layar karena itu menampilkan enam
medan uang sebagai *"belum tersedia"*. Itu **keliru**. Yang hilang hanya pemanggil **refresh** di
layar. Yang **menghitung** nilainya ada di korpus — dua kali, dengan rumus yang sama persis:

| Activity | Dipanggil dari | Langkah yang menghitung |
| --- | --- | --- |
| `Activity/SavePesertaClaim.xml` | `InputRegisterClaimLife.xml` `Submit` b27369 | **8** b4002 ULANG peserta *(`HasRepeat EMBEDDED` b4843, prasyarat kosong b4835)*; reset enam `local.Total*` = `0` b4027–b4159 · **8.1** b4221 ULANG `.AdjustmentList` b4226 *(`EMBEDDED` b4570, **prasyarat kosong** b4562)*; `local.Total… = .KOLOM + local.Total…` b4246, b4292, b4312, b4332, b4352, b4372 · **8.2** b4592 *(tanpa ulang, b4801)*; tulis ke halaman **peserta** b4616–b4743 |
| `Activity/SaveOutStandingLife_Act.xml` | `InputOSClaimLife.xml` `Save to RNM` b21102 | **23** b10638 · **23.1** b10841 *(`EMBEDDED` b11046)* · **23.2** b11067 → b11091–b11217 |

**Jawaban atas pertanyaan yang kami sebut "tidak terjawab":** ikut dijumlah adalah **seluruh baris
`.AdjustmentList` peserta itu** — tanpa memandang `STS_REJECT`, tanpa memandang `IsCheck`. Kedua
langkah penjumlahnya **berprasyarat kosong**. Tidak ada penyaringan sama sekali.

**Dan totalnya ENAM, bukan lima.** `Total Ceding Retention` b20629 / `.TotalCedingRetention` b20636
luput dari daftar kami karena kami mencacah lewat **rujukan** `CheckTotalAdjustmentClaim` — dan ia
satu-satunya total yang **tidak punya aksi refresh**, sehingga ia tidak ikut tercacah.

**Sebab salahnya**, dicatat: kami berhenti pada rule yang **namanya tertulis di section**, lalu
menyimpulkan tentang **medannya**, tanpa menanyakan *siapa lagi yang menulis properti itu*.

#### Yang masih kami tanyakan — OQ-H versi sempit

> `CheckTotalAdjustmentClaim` dirujuk sepuluh kali di `ClaimLifeDetailGCNM` tetapi berkas rule-nya
> tidak ikut dalam ekspor. **Apakah ia hanya me-refresh keenam medan dengan nilai yang sudah
> dihitung `SavePesertaClaim` / `SaveOutStandingLife_Act`, atau ia menghitung ulang dengan aturan
> lain?** Bila yang kedua, mohon kirimkan berkasnya.

**Sementara itu** keenam total dihitung dengan rumus di atas dan disajikan saat Detail dibaca.
⚠️ **Penyimpangan sadar yang kami catat**: Pega menghitungnya saat `Submit` / `Save to RNM` lalu
**menyimpan** hasilnya ke halaman peserta, sehingga layar Pega dapat **basi** sesudah putaran atau
akseptasi sampai tombol simpan ditekan lagi. Di sistem baru ia dihitung saat dibaca, jadi tidak
pernah basi. Arah selisihnya disengaja — dan bila produksi justru mengandalkan angka yang **tersimpan**
*(mis. untuk laporan yang harus cocok dengan cetakan lama)*, beri tahu kami: itu mengubah keputusan.

---

## 27 September 2026 — OQ-I (`Close Claim` dari Outstanding: selesai, atau pindah?)

**OQ-I** *(untuk pemilik ekspor / pengembang Pega — kami sudah memutuskan, dan ingin dikoreksi bila salah)*.

`Close Claim` bukan konektor alur. `Flow/Register_Flow.xml` dibaca utuh: **12 konektor**, dan tidak
satu pun bernama `CloseClaim`. Ia **local action** — `pyLocalAction>CloseClaim` muncul di tepat dua
section: `Section/InputOSClaimLife.xml` **b22837, b22988** dan
`Section/InputAkseptasiClaimLife.xml` **b21457, b21602**.

Activity-nya `ProtectCloseClaim_act` berakhir dengan `Call FinishAssignment` **b838** yang **seluruh
parameternya kosong** *(`TaskStatus`, `PerformFormName`, `ReviewFormName`, `pyHarness`,
`DisplayHarness false`)*. Satu-satunya shape yang menetapkan status kerja adalah **End1**:
`<rowdata REPEATINGINDEX="End1">` **b883**, `pyMOId End1` **b885**, `Data-MO-Event-End` **b901**,
`<pyWorkStatus>Resolved-Completed</pyWorkStatus>` **b899**. Sembilan shape lain ber-`pyWorkStatus`
kosong.

⛔ **Yang ekspor tidak jawab:** dari `Assignment1` *(Outstanding Claim)* tidak ada jalur ke End1
tanpa melewati Medical Check dan Claim Analis. Perilaku mesin Pega untuk `FinishAssignment` yang
dipanggil dari local action **tanpa konektor senama** tidak dapat diturunkan dari ekspor.

> **Pertanyaannya:** di produksi, ketika `Close Claim` ditekan pada layar **Outstanding**, kasusnya
> **selesai** — atau ia justru berpindah ke **Medical Check** mengikuti konektor `Else`?

**Yang kami bangun sementara ini** *(keputusan bb, dan ia dapat diveto)*: `Close Claim` **menutup**
kasus. `T_WORK_CLAIM.STATUS_WORK` diisi `Resolved-Completed` VERBATIM, `TAHAP` **dikosongkan**, dan
setiap rute pengubah menolak kasus itu sesudahnya. Dasarnya niat nyata tombol tersebut: labelnya,
konfirmasi b499 *"Are you sure want to Close Claim?"*, dan gerbang **"is not approved yet"** yang
menahan setiap peserta yang belum diaksep. Ketiganya hanya masuk akal bila tombolnya menutup —
tombol yang sekadar memindahkan kasus tidak perlu menuntut seluruh peserta sudah diaksep.

⚠️ Bila jawabannya "pindah ke Medical Check", yang berubah kecil: satu migrasi mundur dan satu rute.
Mohon dijawab sebelum modul ini dipakai di produksi.

---

## 27 September 2026 — OQ-J (saringan dokumen `KATEGORI_1 = .DOCUMENT` tidak dapat ditiru)

**OQ-J** *(untuk pemilik ekspor / DBA)*.

`Activity/LoadDocumentLife_ACT.xml` memuat daftar dokumen dengan `Obj-Browse` *(langkah 1.2
**b495**)*: `ObjClass` `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM`, `PageName` `DOCUMENT_CLAIM`,
`RowKey` `ID`, dan **satu** saringan — `Field` **`.KATEGORI_1`**, `Condition` **`=`**,
`Value` **`.DOCUMENT`**.

Activity itu berkelas `Int-LIFE_PREMIUM_DETAIL`, jadi `.DOCUMENT` adalah properti **peserta**.
Langkah 1.4 **b904** memperkuatnya: prasyaratnya **b1404** `.DOCUMENT==""` dengan `WhenTrue=3`
*(lewati)* — peserta tanpa nilai `DOCUMENT` tidak memuat dokumen apa pun.

⛔ **Masalahnya:** kolom `DOCUMENT` **tidak ada** di `T_CLAIMLF_PREMIUMLIST_DETAIL`, dan tidak ada
di daftar kolom mana pun yang kami terima. Kami karena itu **tidak dapat meniru saringan itu**.

**Yang kami pakai sebagai gantinya:** kunci tamu `T_CLAIMLF_DOCUMENT.PREMIUM_LIST_DETAIL_ID` —
relasi yang di model baru memang memiliki dokumen tersebut. Untuk data yang lahir di sistem baru
keduanya sama; untuk data warisan, **belum tentu**.

> **Pertanyaannya:** (a) kolom apa `\.DOCUMENT` itu pada baris peserta, dan (b) apakah
> `KATEGORI_1` benar-benar dipakai sebagai penghubung ke peserta, atau ia penggolong jenis dokumen
> yang kebetulan bernilai sama?

⚠️ **Satu penyimpangan lagi, disengaja dan dicatat.** Langkah 1.4.2 **b1129** berprasyarat **b1310**
`DataImage.URLImage==""` dengan `WhenTrue=3` — baris yang **URL penyimpanannya kosong tidak ikut
ditambahkan** ke daftar. Karena penyambungan ke Google Storage belum dilakukan, URL-nya selalu
kosong, sehingga meniru gerbang itu membuat daftar **selalu kosong** dan layar berkata *"tidak ada
dokumen"* untuk peserta yang dokumennya lengkap. Barisnya karena itu **tetap ditampilkan**, dengan
penanda bahwa pranalanya menunggu penyambungan.

### ⛔ Ralat OQ-J — 27 September 2026, beberapa jam sesudah ditulis

**Pertanyaan (a) sudah terjawab, dan yang menjawabnya adalah XML yang belum kami baca saat menulis
OQ-J.** `Activity/SaveAttachLife.xml` langkah 1.2 **b595-596**:

```
Primary.DOCUMENT = @if(Primary.DOCUMENT == "",
                       "DL-" + @pxReplaceAllViaRegex(@CurrentDateTime(),"[^0-9]",""),
                       Primary.DOCUMENT)
Local.IDDoc      = Primary.DOCUMENT                                        b615-616
```

dan langkah 1.7 **b1467** meneruskan `KATEGORI_1 = Local.IDDoc` ke `InsertDocument_Act`.

Jadi `.DOCUMENT` adalah **kunci kelompok dokumen milik peserta** — `DL-` diikuti seluruh angka dari
waktu saat lampiran PERTAMA disimpan, dibuat sekali lalu **tidak pernah ditimpa**. Setiap dokumen
peserta itu menyimpan kunci tersebut di `KATEGORI_1`, dan itulah yang disaring
`LoadDocumentLife_ACT`. Jawaban (b): `KATEGORI_1` **memang penghubung ke peserta**, bukan penggolong
jenis dokumen — penggolong jenisnya `KATEGORI_2`, yang `DocumentLife.xml` tampilkan.

⚠️ **Sebab kami salah:** OQ-J ditulis sesudah membaca activity **pembaca** saja. Aktivitas
**penulisnya** yang memegang jawabannya. Bentuk yang sama dengan ralat OQ-H hari ini — berhenti pada
rule yang namanya tertulis, tanpa menanyakan siapa yang MENULIS medan itu. Dua kali dalam satu hari.

**Yang masih ditanyakan, dan lebih sempit:**

> `T_CLAIMLF_PREMIUMLIST_DETAIL` **tidak punya** kolom `DOCUMENT`. Untuk data yang lahir di sistem
> baru, FK `PREMIUM_LIST_DETAIL_ID` dan kunci kelompok `DL-…` menunjuk himpunan yang sama, jadi
> keduanya setara. **Untuk data warisan yang akan dimigrasi (tiket 13), apakah `KATEGORI_1` selalu
> terisi dan selalu cocok dengan satu peserta?** Bila ada dokumen warisan ber-`KATEGORI_1` yang tidak
> menunjuk peserta mana pun, ia tidak akan terbawa oleh FK — dan kami perlu tahu sebelum migrasi,
> bukan sesudah.

---

## 27 September 2026 — OQ-K (diagnosa: satu kolom lawan satu grid, dan tiga nama kolom)

**OQ-K.1 — `DISEASE_LIFE`: nama kolomnya tidak ada di ekspor** *(untuk DBA)*.

`ReportDefinition/BrowseDiseaseLife_RD.xml` menyebut **properti Pega** — `.Number` *(label `ID`)*,
`.Disease`, `.ICD_Code` — tetapi ekspor ini **tidak memuat `Rule-Obj-Class`** untuk
`ASM-FW-GISFW-Int-DISEASE_LIFE`, sehingga pemetaan kelas-ke-tabelnya tidak ada.

| Properti | Kolom yang kami pakai | Dasar |
| --- | --- | --- |
| `.Disease` | `DISEASE` | `[dugaan kuat]` — nama yang sama persis dengan `T_CLAIMLF_PREMIUMLIST_DETAIL` *(migrasi 003)* |
| `.ICD_Code` | `ICD_CODE` | `[dugaan kuat]` — sama |
| `.Number` | `NUMBER_` | ⛔ `[terbuka]` — **`NUMBER` kata cadangan Oracle**, jadi kolomnya pasti bernama lain, dan nama itu tidak ada di korpus |

> **Mohon dipastikan ketiganya.** Ketiganya terkumpul di satu tempat *(`repository/penyakit.go`,
> blok `kolomNomorPenyakit` dst.)* supaya koreksi Anda adalah **satu suntingan**, bukan perburuan.

**OQ-K.2 — satu diagnosa atau banyak?** *(untuk pemilik ekspor / work owner)*

`Section/ClaimLifeDetailGCNM.xml` **b3923** `pyPageListProperty` **`.DiagnoseList`**, disajikan
sebagai **`RepeatGrid`** — jadi satu peserta dapat punya **banyak** diagnosa.
`Activity/SetDisease.xml` menulis `.DISEASE` b260 dan `.ICDCODE` b307 pada halaman berkelas
`ASM-FW-GISFW-Data-DiagnoseLife`, lalu `Obj-Save pyWorkPage` b390.

⛔ Tetapi `T_CLAIMLF_PREMIUMLIST_DETAIL` hanya punya **`DISEASE`** dan **`ICD_CODE` TUNGGAL**
*(migrasi 003)*. **Satu lawan banyak.**

> **Pertanyaannya:** apakah di produksi seorang peserta benar-benar dapat memiliki lebih dari satu
> diagnosa, atau grid itu hanya cara menampilkan satu baris? Bila benar banyak, diperlukan tabel
> anak *(mis. `T_CLAIMLF_DIAGNOSA`)* dan itu keputusan skema — bukan keputusan executor.

**Sementara itu:** pencarian diagnosanya **sudah ada dan berbatas** *(`Find Disease` → dua kotak
`ICD Code` dan `Disease`, disambung `AND`, batas 500 / halaman 50 dari rule)*, tetapi tombol
**`Choose` b2509 dinyatakan belum tersedia**. Memilih diagnosa ke kolom tunggal berarti **membuang
diagnosa kedua dan seterusnya** — diam-diam, dan tanpa ada yang tahu.

### ✅ OQ-K.1 DITUTUP — 27 September 2026 *(katalog DEV)*

`[data DBA — katalog DEV, akun POOLDATA, agregat saja]` `POOLDATA.DISEASE_LIFE` berkolom:

| Kolom | Tipe | Isi terpanjang |
| --- | --- | ---: |
| `ID` | `VARCHAR2(100)` | — |
| `ICD_CODE` | `VARCHAR2(100)` | 7 |
| `DISEASE` | `VARCHAR2(1000)` | 290 |

97.586 baris. `repository/penyakit.go` kini memakai ketiganya apa adanya.

⛔ **Dan tebakan kami keliru, dengan cara yang layak dicatat.** Kami menulis `NUMBER_` beralasan
*"`NUMBER` kata cadangan Oracle, jadi kolomnya pasti bernama lain — dan nama itu tidak ada di
korpus"*. Separuh pertamanya benar. Separuh keduanya **salah**: namanya **ada** di korpus, dan kami
membacanya tanpa melihatnya —

```
ReportDefinition/BrowseDiseaseLife_RD.xml
  b598  <pyFieldName>.Number</pyFieldName>
  b599  <pyFieldLabel>ID</pyFieldLabel>      <- namanya, di baris berikutnya
  b676-677  pasangan yang sama, kedua kalinya
```

Baris b599 **ada di keluaran grep kami sendiri** saat membaca report definition itu, dan kami
memperlakukannya sebagai **label layar**. Pelajarannya sempit dan tajam: sebelum menyatakan sesuatu
*"tidak ada di korpus"*, periksa apa yang **sudah terbaca** — bukan hanya apa yang sudah dicari.

### ✅ OQ-K.2 DITUTUP — 27 September 2026 *(XML, bukan katalog)*

**Banyak diagnosa per peserta.** Jawabannya ada di `Section/ClaimLifeDetailGCNM.xml`, dua baris di
bawah tempat kami berhenti membaca:

```
b3914  kelas grid   ASM-FW-GISFW-Data-DiagnoseLife
b3923  pyPageListProperty  .DiagnoseList
b4690  <pyLabel>Add</pyLabel>     -> addRow    b4700/b4841
b6160  <pyLabel>Delete</pyLabel>  -> deleteRow b6170
```

⚠️ **Sebab kami salah**: kami melihat **grid** dan menyimpulkan *"satu lawan banyak, tidak dapat
diputuskan"*. Grid memang dapat berarti tampilan satu baris — tetapi grid ber-`Add` **dan**
ber-`Delete` tidak dapat. Kami berhenti pada bentuk tampilannya tanpa membaca **tombolnya**.

Tabel `T_CLAIMLF_DIAGNOSE` lahir di migrasi `018` *(keputusan **bd**)*.

⛔ **Satu hal TETAP terbuka, dan lebih sempit — OQ-L:** daftar pilihan `GROUPDIAGNOSE`.
Dropdown b5860/b5863 ber-`pyListSource` **`associated`**, artinya daftarnya hidup pada **rule
properti** `.GROUPDIAGNOSE` di kelas `Data-DiagnoseLife` — dan rule itu **tidak ada di ekspor**:
`GROUPDIAGNOSE` muncul di **tepat satu** berkas korpus, yaitu section itu sendiri.

> **Mohon kirimkan** ekspor `Rule-Obj-Property GROUPDIAGNOSE` kelas
> `ASM-FW-GISFW-Data-DiagnoseLife` **beserta daftar lokalnya** — atau, bila lebih mudah, daftar
> nilai yang pernah dipakai di produksi. Kolomnya sudah dibuat *(migrasi 018, `VARCHAR2(255)`)*;
> nilainya **tidak dikarang**, dan layar menyatakan daftarnya belum ada.

**Keputusan sementara — butir bf** *(27-09-2026)*: backend menerima `groupDiagnose` sebagai teks
≤ 255 **tanpa** memeriksa daftar apa pun; frontend menampilkan `BelumTersedia` bernama kepala kolom
b4490 `GROUP DIAGNOSE`. Nilai yang **sudah** tersimpan tetap ditampilkan — yang belum ada hanyalah
cara memilihnya. `pyLabelPreview` sel itu *(b5854)* **kosong**, jadi kepala kolom adalah satu-satunya
kata yang sah untuknya.

⚠️ Ketika OQ-L dijawab, pemeriksaan daftar **menggantikan** pemeriksaan panjang — dan
panjangnya tetap berlaku.

## 28 September 2026 — OQ-M (sensus akhir: keputusan yang bukan milik executor)

Sensus §3.1 memutuskan keempat puluh satu baris *"tercatat, belum berkode"* (bab *Sensus akhir*
di `PARITAS-LAYAR-DAN-AKSI.md`). Sembilan butir di bawah **tidak dapat dibangun tanpa keputusan**:
masing-masing menyentuh skema *(migrasi `020`+ hanya dari keputusan tercatat)*, desain data, atau
wewenang. Baris XML dari pecahan `sed 's/></>\n</g'`.

**OQ-M1** *(untuk work owner)* — separuh gerbang dialog Edit Date. Keempat isian
`EditDateClaimLife_Section` berprasyarat baca-saja
`pyPosition!='ReasLifeAdmin' || PremiumListSummary.CLAIM_NO!=''` *(b1000, b1313, b1550, b1788)*.
Separuh pertama ditiru *(tahap Outstanding + peran Admin, `models.TahapBolehUbahTanggalKlaim`)*.
Separuh kedua **tidak**: Pega menomori klaim di Save Outstanding *(`GetSequenceNumber_SQL` dari
`SaveOutStandingLife_Act` b8057)*, aplikasi di **pendaftaran**. Meniru hurufnya mengunci keempat
isian pada SETIAP klaim. Maksudnya — *"sesudah Save Outstanding, tanggal tidak diubah lagi"* —
menunggu padanan Save Outstanding, yang belum punya rute.
> Apakah tanggal klaim boleh diubah Admin selama kasus di Outstanding, atau harus terkunci sejak
> suatu peristiwa (yang mana)?

**OQ-M2** *(untuk DBA / pemilik hilir)* — cermin warisan `UpdateDateClaimLife_SQL` b85-90:
`UPDATE OS_AKSEPTASI_KLAIM_LIFE SET LAPSE_DATE = <DOL>, CLAIM_RECEIVED_DATE, COMPLETE_DATE,
CONFIRMATION_DATE WHERE CASEID AND NAME_OF_INSURED AND CERTIFICATE_NO`. Dua hal: (1) baris warisan
yang kita tulis membiarkan `NAME_OF_INSURED` NULL *(nol nama orang)*, sehingga WHERE itu mengenai
**nol** baris; (2) kolom `LAPSE_DATE` di tabel warisan diisi **DOL**, bukan tanggal lapse.
> Apakah hilir membaca keempat tanggal ini dari `OS_AKSEPTASI_KLAIM_LIFE`? Bila ya, dengan kunci
> apa baris cermin dicari tanpa nama tertanggung, dan apakah `LAPSE_DATE = DOL` disengaja?

**OQ-M3** *(untuk work owner + DBA)* — `CountClaimAmountLife_Act`. Rumusnya kini aturan murni yang
teruji (`models/klaimbayar.go`): `CLAIM_PAID = CLAIM_GROSS × @divide(PCTClaim,100,5)` *(b702/b723)*,
galat b526 bila gross > share *(b949)* kecuali L12/L13/L14/L18 *(b972)*, lalu keluar activity
*(kode 6, b865)* sebelum Spreading. **Tidak tersambung**: `T_CLAIMLF_ADJUSTMENT` (migrasi 004) tidak
punya kolom `PCT_CLAIM` / `CLAIM_PAID`, dan tidak ada rute yang menyunting baris adjustment.
> Mohon keputusan migrasi dua kolom itu (tipe, presisi) dan rute sunting adjustment-nya.

**OQ-M4** *(untuk work owner — desain)* — Upload CSV peserta di Register (`UploadCSVClaimLife_Act`
+ `SetClaimXOL_Act`). ⛔ **Ralat OQ-F**: pemetaannya ADA — 28 kolom berkas + `PL_NUMBER` dari polis *(b478-b1065)*. Tetapi CSV itu
membawa uang *(`SUM_INSURED`, `CLAIM_AMOUNT`, share)*, **jendela valuasi**, dan `NAME_OF_INSURED`
dari berkas klien, lalu bendera XOL *(b1169)* mencegah pembacaan dari polis menimpanya *(`LoadDataPeserta_Act` b282/b327)*.
Aplikasi ini sengaja membaca peserta ULANG dari sumbernya saat pendaftaran, dan `ValidasiDOL`
memvalidasi terhadap jendela yang TERSIMPAN.
> Apakah klaim XOL boleh membawa peserta di luar PremiumList dari berkas? Bila ya: kolom mana yang
> dipercaya dari berkas, dan siapa yang berwenang mengunggahnya?

**OQ-M5** *(untuk work owner)* — dialog Reject Outstanding (`RejectOSClaimLife_Sec`): Date b790, PIC
b975, Remarks b1687, Submit b3117. `RejectOSClaimLife_Act` menambah baris `KomiteList`
*(`IDKomite="Claim Admin"`, `KomiteAproval=2`, `KomiteComment` = Remarks — b2173-b2263)*. Aplikasi
menolak baris tanpa dialog dan tanpa alasan; `T_CLAIMLF_JEJAK` tidak punya kolom komentar, dan
`T_KOMITE_KOMITELIST` terikat FK ke `T_GENERAL_KOMITE` (baris "Claim Admin" bukan tingkat komite).
⚠️ Dialognya sendiri SENGAJA tidak dibangun lebih dulu: dialog yang menerima Remarks lalu
membuangnya menipu pemakai — ia mengira alasannya tercatat.
> Di mana alasan penolakan disimpan: kolom komentar di jejak, atau baris riwayat Komite?

**OQ-M6** *(untuk work owner)* — `DeletePesertaClaimLife`. Activity-nya nol hapus *(b249, b337,
Obj-Save b443)*, tetapi tombol `DELETE` `InputOSClaimLife` b17865 lebih dulu menjalankan `deleteRow`
b17874 atas grid `PremiumListDetail`: di Pega peserta **dicabut** dari kasus. ADR-U-0031 menyatakan
hapus = penanda.
> Mencabut peserta: penanda (dan layar menyembunyikannya), atau baris benar-benar dilepas?

**OQ-M7** *(untuk work owner)* — `RetroDetailClaimLife`, panel rincian grid treaty-year
(`AdjustmentDetail` b10388-b10389): Reinsurer, Currency, Percent Share, Claim Retro *(b1325-b2498)*.
Grid induknya belum tampil, dan tabel yang akan ditampilkannya KOSONG untuk setiap klaim baru:
`HitungSpreading` (`services/spreading.go`) nol pemanggil produksi karena masukan rate-nya tanpa
pembaca. Di Pega rate dicari lewat `OUTWARDRATEID` produk — hasil mengurai
`M_PRODUCT_LIFE.JSONDATA` (`GetProductLife` b84-86, Java b1016, b1416-b1604), yang AC 38 larang —
atas `RATE_LIFE`, view master yang katalognya `[data DBA]` dan tidak ditiru skema uji.
> Dari mana `OUTWARDRATEID` dibaca tanpa mengurai `JSONDATA`, dan bolehkah `RATE_LIFE` dibaca
> (view master, seperti izin sempit butir bh)? Tanpa keduanya Spreading tidak dapat berjalan dan
> panel ini tidak punya isi.

**OQ-M8** *(untuk work owner — wewenang, temuan)* — rute DOL yang sudah ada
(`PUT …/tanggal-kejadian`, tiket 06) **tidak** menegakkan gerbang b1000 yang sama: siapa pun yang
beridentitas dapat mengubah DOL pada kasus terbuka di tahap apa pun. Rute tiga tanggal yang baru
menegakkannya. Menambahkannya ke rute DOL adalah **perubahan authz** — menunggu persetujuan, tidak
dilakukan sendiri. Bedanya dengan rute baru bukan standar ganda: rute baru LAHIR dengan gerbang
yang XML-nya tuntut (posisi kasus → tahap; perutean penugasan ke Admin → peran pelaku, pola
`TahapLayanan.Pindah`), sedangkan rute DOL sudah dipakai — mempersempitnya mengubah siapa yang
hari ini boleh mengubah DOL.
> Samakan gerbang DOL dengan tiga tanggal lainnya (Admin, tahap Outstanding)?

✅ **DITUTUP 28-09-2026 — butir bj `[DIPUTUSKAN; veto work owner]`**: disamakan menurut XML;
perubahan authz dicatat bertanggal di tiket 07.

**OQ-M9** *(untuk work owner)* — `ValidasiClaimReceived_Act` b1120, dipicu perubahan
`CLAIM_RECEIVED_DATE` di dialog Edit Date. Aturannya (`models.PenandaBatasHari`) dan ambangnya
(`services.AmbangKlaim.Hitung`, butir ba/bh — `MAXEXPIREDCLAIM` dari view produk) sudah ada, tetapi
**nol pemanggil**. Menyambungkannya ke simpan tiga tanggal menuntut dua hal yang belum diputuskan:
(1) penandanya `.MAXCLAIM_RECEIVED` tidak punya kolom — disimpan, atau hanya ditampilkan?;
(2) `Hitung` menghitung KEDUA penanda sekaligus dan gagal bila STNC tidak dapat dihitung, padahal
dialog ini hanya menyentuh yang pertama. Di Pega penandanya tidak memblokir (OQ-G).
> Penanda Claim Received: kolom baru, atau tampil-saja sesudah simpan?

✅ **DITUTUP 28-09-2026 — butir bk `[DIPUTUSKAN; veto work owner]`**: nol penulis tabel di korpus
(hanya `Property-Set` b582, dibaca sel b12131) → **dihitung saat baca**, tanpa kolom. Bukti di
tiket 06.

## 28 September 2026 — OQ-N (Save to RNM, GILIRAN-11 paket 1)

`Activity/SaveOutStandingLife_Act.xml` dibaca UTUH sebagai pohon (29 langkah teratas, pecahan
`sed 's/></>\n</g'`) dan dibangun sebagai `POST /api/klaim-life/{id}/outstanding`
(`services/simpanrnm.go`). Empat hal di bawah tidak dapat diputuskan executor.

**OQ-N1** *(untuk work owner — skema)* — bendera `pyWorkPage.Save`. Langkah 21 menyetelnya `""`,
langkah 24 `1` bila seluruh `TempError.CARIx==0` dan `Attachment.pxResults(1).CountAttach!=0` (b11448);
ia **mematikan** tombol `Save to RNM` (`pyDisabledWhen` b21095) dan **membuka** kontainer
*Participant Details* (b15490; di Register kebalikannya, b16855/b23236). Tidak ada kolom untuknya di
migrasi 001–020. Akibatnya tombol tetap hidup dan dapat ditekan ulang — **aman**, sebab tulisannya
idempoten: nomor hanya bila `CLAIM_NO` kosong, `STS_REJECT=0` hanya pada baris tanpa status.
⛔ Karena itu pula langkah 22.1.3.2 **tidak** ditiru hurufnya: ia menulis `.STS_REJECT=0` pada
**setiap** baris tanpa precondition; di Pega itu aman hanya karena bendera mematikan tombolnya
sesudah simpan pertama. Tanpa bendera, menirunya berarti membatalkan penolakan Admin diam-diam.
> Kolom bendera simpan di `T_WORK_CLAIM` (migrasi baru), atau cukup keadaan turunan?

**OQ-N2** *(untuk work owner / DBA)* — klaim ganda (langkah 11.2–11.6) membaca
`OS_AKSEPTASI_KLAIM_LIFE` dengan kunci `CEDINGCO`, `NAME_OF_INSURED`, `DOB`, `CERTIFICATE_NO`,
`PL_NUMBER`. Dua hal: (1) baris warisan yang **aplikasi ini** tulis saat pendaftaran membiarkan
`NAME_OF_INSURED`, `DOB`, `CEDINGCO` NULL, sehingga klaim ganda antarklaim **baru** tidak pernah
tertangkap — hanya baris era Pega; (2) SQL health membandingkan `LAPSE_DATE = TO_DATE({CARI6},
'DD/MM/YYYY')` padahal `CARI6 = .DATE_OF_LOSS` mentah (b3189) — ditiru maksudnya (DOL sama).
Nama dan tanggal lahir tidak pernah meninggalkan basis data: pencocokannya menggabung baris sumber
`M_LIFE_PREMIUM_DETAIL` di dalam SQL (`repository/gandawarisan.go`).
> Bolehkah cermin warisan mengisi ketiga kolom itu, supaya klaim ganda antarklaim baru tertangkap?

**OQ-N3** *(untuk work owner)* — gerbang langkah 27: `RetroID=="L0000141" ||
SecurityReinsurerID=="L0000134"` (b11794) dan `RetroID=="1000013"` (b11817) KELUAR sebelum Arasapas.
Modul Komite **membuang** gerbang yang sama *(OQ-064, `KomitePostAdjustment` langkah 9)*; keputusan
itu tidak menyebut Claim Life, jadi di sini **XML yang menang** (`ArasapasDilewatiRetro`).
> Berlakukah OQ-064 juga untuk Save to RNM?

**OQ-N4** *(untuk pemilik ekspor — temuan)* — tiga residu yang tidak ditiru, dengan buktinya:
(a) langkah 5 `@contains(.Protect,"1")` → pesan *"Claim gross tidak boleh lebih besar dari Share
Nusantara Re"*: properti `.Protect` **nol penulis** di seluruh korpus Claim Life, jadi pesannya tidak
pernah muncul; (b) `ADJUSTMENT_DATE` dan `PrintFaceClaim` (langkah 22.1.3.2) tanpa kolom — penanda
"sudah tersimpan" per baris digantikan status baris itu sendiri; (c) gerbang dokumen pertama menyaring
`.IsAccept=="true"` (b1181), yang tidak punya kolom; `IS_CHECK` dipakai sebagai padanan terdekat
(`[terbuka]` sejak tiket 03).

