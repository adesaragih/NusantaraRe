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
