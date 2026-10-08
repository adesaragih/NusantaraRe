# Alur akseptasi `Submit` — dibaca dari ekspor Pega

**Keputusan pemilik proses 6 Oktober 2026**, jawaban atas tiga pertanyaan
terbuka di [`KEPUTUSAN-SASARAN-TULIS.md`](KEPUTUSAN-SASARAN-TULIS.md):

> 1. "sasarannya tetap TREATYEXCHANGEYEARLY"
> 2. "save dan submit beda, save hanya menyimpan di db sedangkan submit adalah
>    alur dalam akseptasi silahkan baca xml untuk tiap alurnya kemana saja"
> 3. "siapa saja yang terdaftar seperti yang di xml sesuaikan seperti yg di xml"

Berkas ini menjawab (2) dan (3) dari ekspor. Untuk (1) lihat §4 — ia membawa
akibat yang harus diketahui sebelum dibangun.

---

## 1 · `Save` ≠ `Submit`

| | Yang dilakukan |
|---|---|
| **Save** | menyimpan ke tabel, selesai. Nol perubahan status, nol perpindahan posisi |
| **Submit** | **menyimpan DAN menjalankan tangga akseptasi** |

Rantai `Submit` dari `Activity/TreatyInSubmit.xml` (manifes `pxRuleReferences`,
sesudah `pyIncludedRuleXML` bersarang dibuang dengan hitung kedalaman):

```
TreatyInSubmit
  1  TreatyInInputVis      kesiapan isian
  2  TreatyInCheckError    pemeriksaan → properti ERRMSG
  3  SaveTreatyIn_Act      penyimpanan
  4  Akseptasi_DT          TANGGA AKSEPTASI          ← inti alurnya
  5  AddCommentList_Act    menambah baris riwayat
```

⛔ Jadi `Submit` **memanggil `Save`**, bukan menggantikannya. Membangun
`Submit` tanpa `Save` lebih dulu berarti membangun langkah 4 di atas langkah 3
yang belum ada.

---

## 2 · Tangga akseptasi — `DataTransform/Akseptasi_DT.xml`, **sesudah koreksi**

Tiga medan yang berubah di setiap langkah: `TreatyIn.Position`,
`TreatyIn.StatusAkseptasi`, `TreatyIn.PositionUsername`.

### 2.0 ⭐ KOREKSI PEMILIK PROSES — ia yang berlaku, bukan ekspor

**6 Oktober 2026**, atas bentuk pertama berkas ini:

> "untuk SPVTREATY1 · TREATY1 jangan digunakan dulu. alur akseptasinya itu
> dari ReasTreatyInAdmin > ReasTreatyInSecHead > ReasTreatyInDeptHead >
> ReasTreatyInDirector seperti ini"

Dua hal berubah, dan keduanya menyempit — bukan menambah:

| | Ekspor Pega | Yang dibangun |
|---|---|---|
| **Penentu langkah pertama** | `OperatorID.pyTelephone` berisi kode regu `SPVTREATY1`/`SPVTREATY2`/`TREATY1`/`TREATY2` | **PERAN** pemakai — `ReasTreatyInAdmin` |
| **Jenjang** | empat posisi, `GroupLeader` ikut dibaca | **empat anak tangga berurut**, `GroupLeader` di luar |

⛔ Jadi cabang `pyTelephone` **tidak dibangun**. Ia tetap tertulis di §2.4
sebagai catatan ekspor, supaya ronde berikutnya tidak menemukannya di XML dan
mengira ia terlewat.

### 2.1 Jalur biasa (`RevisionState != 1`) — tangga yang dibangun

```
ReasTreatyInAdmin → ReasTreatyInSecHead → ReasTreatyInDeptHead → ReasTreatyInDirector → Resolve Complete
```

| Position sekarang | Accept → | Reject → | Decline → |
|---|---|---|---|
| `ReasTreatyInAdmin` (atau `""`) | `ReasTreatyInSecHead` · `Accept` | `—` | `—` |
| `ReasTreatyInSecHead` | `ReasTreatyInDeptHead` · `Accept` | `ReasTreatyInAdmin` · `Reject` | `""` · `Decline` |
| `ReasTreatyInDeptHead` | `ReasTreatyInDirector` · `Accept` | `ReasTreatyInAdmin` · `Reject` | `""` · `Decline` |
| `ReasTreatyInDirector` | `""` · **`Resolve Complete`** | `ReasTreatyInAdmin` · `Reject` | `""` · `Decline` |

⭐ Kolom `Reject` dan `Decline` **tetap dari ekspor** — koreksi pemilik proses
menyebut jalur naiknya saja, dan dua cabang turun itu tidak disentuhnya.

⭐ Pada **Reject**, `PositionUsername` diisi
`TreatyIn.CommentList(1).OperatorName` — dikembalikan kepada **penulis
komentar pertama**, yaitu yang mengajukannya.

⭐ Pada **Decline**, ketiganya dikosongkan: `Position` `""`,
`PositionUsername` `""`, status `Decline`.

⛔ **`PositionUsername` pada jalur NAIK menjadi pertanyaan terbuka.** Di
ekspor ia diisi nama orang (`IRVANDY`, `YOHANESKRISTIAWAN`, `NANDINA`) yang
hanya dicapai lewat cabang kode regu yang kini dibuang. Tanpa cabang itu, nol
sumber mengisinya. Tiga kemungkinan — **tidak satu pun ditebak di sini**:
dikosongkan (berkas menunggu POSISI, bukan orang), diisi dari daftar pemakai
berperan itu, atau tetap nama tertanam. Jalur naik tidak dapat dibangun
sebelum ini dijawab; yang lain dapat.

⭐ Pada **Reject**, `PositionUsername` diisi
`TreatyIn.CommentList(1).OperatorName` — dikembalikan kepada **penulis
komentar pertama**, yaitu yang mengajukannya.

⭐ Pada **Decline**, ketiganya dikosongkan: `Position` `""`,
`PositionUsername` `""`, status `Decline`.

### 2.2 Jalur REVISI (`RevisionState == 1`)

Tangganya **lebih pendek**, dan itu disengaja — revisi tidak mengulang seluruh
persetujuan:

| Position sekarang | Hasil |
|---|---|
| `""` atau `ReasTreatyInAdmin` | → `ReasTreatyInSecHead` · `Accept` · ~~`AGUNGPUTRAANDALAS`~~ |
| `ReasTreatyInSecHead` **Accept** | → `""` · **`Resolve Complete`** · `RevisionState=""` · `ViewState=""` |
| `ReasTreatyInSecHead` **Reject** | → `ReasTreatyInAdmin` · `Reject` · `CommentList(<LAST>).OperatorName` · **`RevisionState="1"`** |
| `ReasTreatyInSecHead` **Decline** | → `""` · `Decline` · `RevisionState=""` |

⚠️ Nama `AGUNGPUTRAANDALAS` dicoret: ia nama tertanam yang sama dengan
yang dibuang di §2.0, dan ia jatuh di bawah pertanyaan terbuka
`PositionUsername` yang sama.

⚠️ Dua beda halus yang mudah terlewat:
- Revisi **langsung tuntas** di Sec Head — tidak lewat Dept Head maupun
  Director.
- `Reject` di jalur revisi memakai komentar **TERAKHIR** (`<LAST>`),
  sementara jalur biasa memakai yang **PERTAMA** (`(1)`).

### 2.3 Pagar luar

Seluruh tangga dibungkus
`TreatyIn.StatusAkseptasi != "Resolve Complete"` — kontrak yang sudah tuntas
tidak dapat dikirim ulang. Itu pula syarat yang menyembunyikan tombol `Submit`
di `Section/TreatyInfoSubmit.xml`.

---

### 2.4 Yang ADA DI EKSPOR tetapi TIDAK dibangun

Dicatat supaya ronde berikutnya tidak menemukannya di XML dan mengira ia
terlewat — ketiganya dibuang dengan sengaja, bukan luput:

| Di ekspor | Mengapa tidak dibangun |
|---|---|
| cabang `OperatorID.pyTelephone` → `SPVTREATY1` · `SPVTREATY2` · `TREATY1` · `TREATY2` | pemilik proses: *"jangan digunakan dulu"*. Langkah pertama ditentukan PERAN, bukan kode regu |
| nama tertanam `IRVANDY` · `AGUNGPUTRAANDALAS` · `YOHANESKRISTIAWAN` · `NANDINA` | hanya dicapai lewat cabang di atas. Lihat pertanyaan terbuka `PositionUsername` di §2.1 |
| `ReasTreatyInGroupLeader` sebagai sumber langkah | di luar tangga yang dinyatakan pemilik proses |

⚠️ `pyTelephone` dipakai ekspor sebagai **kode regu**, bukan nomor telepon.
Dicatat sebab siapa pun yang membacanya sebagai nomor telepon akan membangun
cabang yang tidak pernah berjalan — dan cabang mati lebih sukar ditemukan
daripada cabang yang tidak ada.

---

## 3 · Siapa yang terdaftar (jawaban pertanyaan 3)

Pemilik proses menyebut **empat peran, berurut**, dan daftar peran aplikasi
memuat keempatnya:

```
1  ReasTreatyInAdmin        ← pengaju
2  ReasTreatyInSecHead
3  ReasTreatyInDeptHead
4  ReasTreatyInDirector     → Resolve Complete
```

⚠️ **Dua peran lain ADA di daftar peran aplikasi tetapi BUKAN anak tangga:**
`ReasTreatyInGroupLeader` dan `ReasTreatyInUnderwriting`. Keduanya tidak
dinyatakan pemilik proses dan tidak dibangun sebagai tujuan.

⛔ Itu **bukan** berarti keduanya boleh dihapus dari daftar peran — peran
dapat dipakai untuk hak lihat atau hak sunting tanpa pernah menjadi tempat
berkas menunggu. Yang dinyatakan di sini hanya: nol cabang tangga menuju ke
sana.

⭐ Di ekspor, `ReasTreatyInGroupLeader` memang **tidak pernah menjadi
tujuan** di mana pun pula — nol cabang menyetel `Position` ke nilai itu,
ia hanya dibaca. Jadi keputusan pemilik proses dan ekspor sepakat di titik ini.

### 3.1 ⚠️ Modul tetangga memutuskan tangga yang LEBIH PENDEK

`modul/nbtreatyin/backend/models/tangga.go` sudah memakai nama posisi yang
sama persis, tetapi **berhenti di Dept Head**:

```go
PosisiTangga = []string{PosisiAdmin, PosisiSecHead, PosisiDeptHead}
```

Komentarnya menyebut `ReasTreatyInGroupLeader` dan `ReasTreatyInDirector`
*"tidak dibangun (AC 10)"*.

⛔ **Keduanya tidak saling membatalkan**, dan tidak boleh disamakan
diam-diam: itu modul lain (New Business Treaty In) dengan pemilik tiket,
spec, dan AC sendiri. Modul tidak boleh saling impor, jadi tetapan posisi
Treaty In ditulis sendiri — bukan diambil dari sana.

⭐ Yang perlu diketahui: **nama posisinya identik**, jadi satu basis data
pemakai melayani keduanya; yang berbeda hanya berapa anak tangga yang dinaiki.
Treaty In naik empat, NB Treaty In naik tiga.

## 4 · ⛔ Jawaban (1) membawa akibat yang harus diketahui

Pemilik proses: *"sasarannya tetap TREATYEXCHANGEYEARLY"*.

Itu dapat dibangun. Dua hal yang harus diketahui sebelum dibangun:

**a. Tabelnya BERKUNCI TAHUN, bukan kontrak.** `TREATYEXCHANGEYEARLY`
berkolom `TREATYYEAR` + `CURRENCY` — 140 baris untuk 25 mata uang. Nol
kolomnya menunjuk kontrak. Jadi menyimpan kurs dari satu kontrak mengubah
kurs **setiap kontrak pada tahun yang sama**.

**b. Tabelnya DIPAKAI BERSAMA modul lain.** `aggregate` membacanya untuk
konversi, `treatycontractout` lewat `MasterKursTahunanTCO`, dan `nbfacin`
menyebutnya 21 kali. Perubahan dari layar Treaty In sampai ke ketiganya.

**c. Penjaga `TestWarisanHanyaDibaca` melarangnya hari ini.** Tabel warisan
hanya boleh dibaca; menulis ke sana menuntut penjaga itu dilonggarkan dengan
sengaja, beserta sebabnya tertulis.

⛔ Ketiganya bukan penolakan — keputusan sudah diambil. Ketiganya hal yang
harus dinyatakan supaya yang menekan `Save` di grid itu tahu bahwa ia sedang
mengubah data bersama, bukan data kontraknya sendiri.

---

## 5 · Yang masih menghalangi pembangunan

[`KEPUTUSAN-SASARAN-TULIS.md`](KEPUTUSAN-SASARAN-TULIS.md) §2 — pemuat
menghapus lalu mengisi ulang setiap kontrak. Sampai jalan keluarnya dipilih,
**nol `Commit` ke tabel pendaratan**.


## 6 · ⭐ Dibangun 7 Oktober 2026

Tangga di atas kini DIIKAT ke pemakai (`services/simpan_kontrak.go`, `KirimKontrak`):

- Peran pemakai = WORKBASKET-nya (`M_LOGIN_GO_WORKBASKET`, sama dengan `OperatorID.pyWorkBasketList`).
  Penekan wajib memegang workbasket `Position` berkas.
- `PositionUsername` jalur naik — keputusan pemilik proses: pemegang workbasket posisi berikutnya di
  Kelola User (`LOGIN_ID`, dipisah koma). Data 7 Oktober: SecHead/DeptHead/Director dipegang
  `Dastin`, `JEFRIHARI`, `SUPERADMIN`; Admin juga `ADESAMUEL`.
- §5 di atas TERJAWAB — lihat `KEPUTUSAN-SASARAN-TULIS.md` §5.
