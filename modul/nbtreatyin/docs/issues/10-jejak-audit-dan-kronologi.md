# 10: Jejak audit dan kronologi — identitas akses terpisah dari nama tampilan

**Status:** selesai — penahan tersisa hanya pihak luar: **K11** (skema uji Oracle — uji bertag `db` AC 72 sudah ditulis, belum dijalankan), **F2** (tiga penyimpangan K4, menunggu konfirmasi WO), **B4** (`DIV` tanpa sumber di `inti.Pelaku`), **C8** (tipe kolom fisik `HISTORYAKSEPTASIPRODUCTION`, DBA) *(putaran 2, paket P11 04-10-2026; semula: sebagian — konsolidasi P10 04-10-2026; sebagian, implementasi 2026-10-03; awalnya ready-for-agent)*
**Blocked by:** 03
**Menutup:** AC 39 · 40 · 41 · 42 · 43 · 44 · 71 · 72 *(8 AC)* — US 15 · 16 · 20 · 40 · 41 · 42

## Hasil & nilai pengguna

Hari ini riwayat akseptasi mencatat siapa memutuskan apa. ⛔ Tetapi penampung identitas operator
**tidak pernah diisi** dari aplikasi — `[terverifikasi]` **nol dari 1.151** naskah SQL di 21 modul
menyebutnya sebagai kolom. ⚠️ Dan medan "nama operator" diisi dari **dua sumber berbeda** di dua
tahap berdekatan: pengenal akun di satu tahap, nama tampilan di tahap lain.

Sesudah tiket ini, setiap tindakan mencatat **identitas akses login**-nya terpisah dari **nama
tampilan**-nya, dan ⭐ keduanya tidak lagi tertukar — jejaknya tetap benar walau nama tampilan
seseorang berubah.

## Area codebase

- Lapisan repository: penulisan riwayat akseptasi
- Lapisan service: kronologi dan catatan pengguna

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pemisahan yang sudah benar | `RDBList\InsertViewSuggest_SQL.xml` — penampung akses-login diisi dari pengenal akun, terpisah dari penampung nama tampilan |
| ⛔ Pengisian yang **bug** | `DataTransform\DeptHeadTreatyInUW_preDT.xml` — medan nama operator diisi dari **pengenal akun** |
| Pengisian yang benar | `DataTransform\InputPolicyTreatyIn_preDT.xml` — dari **nama tampilan** |
| Mekanisme kronologi | `DataTransform\AddToListCommentsPolicyTreatyIn_DT.xml` — empat medan: catatan, penanda persetujuan, tanggal, operator |
| ⛔ Nama orang di dalam teks pesan | `DataTransform\DeptHeadTreatyIn_UW_postDT.xml` — **1** kemunculan; pesan lain di berkas sama mengambilnya dari data |

## ADR terkait

- ⭐ **ADR-0007** — jejak audit merekam siapa + kapan untuk setiap transisi **dan setiap jalur balik**

## Acceptance criteria

- [x] **AC 39** — penampung identitas operator diisi dari **identitas akses login**
- [x] **AC 40** — penampung nama diisi dari **nama tampilan**
- [x] **AC 41** — penampung identitas operator **terisi** pada setiap penulisan riwayat
- [x] **AC 42** — medan nama operator diisi dari **nama tampilan di setiap tahap**, termasuk tahap
      jenjang ketiga
- [x] **AC 43** — setiap perpindahan tahap menulis **satu baris riwayat**
- [x] **AC 44** — pemberitahuan menyebut nama orang **dari data**; ⛔ tidak tertanam di dalam teks
- [x] **AC 71** — catatan pengguna tersimpan bersama tanggal dan operatornya
- [ ] 🟡 **AC 72** — riwayat dapat dibaca **berurutan waktu** *(P11: `repository/riwayat_db_test.go` `TestDaftarRiwayatBerurutWaktu` — K11)*

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **16** | berapa baris yang penampung identitas operatornya **sudah terisi** sekarang | tidak menahan — ⚠️ bila bukan nol, ada penulis **di luar aplikasi** yang belum terpeta |

## Perintah verifikasi

1. Lakukan satu tindakan di tiap jenjang — ⭐ tiga baris riwayat, **masing-masing berisi kedua
   identitas**.
2. Ubah nama tampilan seorang pengguna — ⭐ riwayat lama **tidak berubah**.
3. Baca pemberitahuan — ⭐ namanya **dari data**, dan berubah mengikuti data.

## Catatan

⚠️ `[penyimpangan sadar]` Mengisi penampung identitas operator adalah **perbaikan jejak audit**,
bukan peniruan — kolomnya selama ini kosong. ⭐ Dan pengisian nama operator dari pengenal akun
adalah **bug**, ⛔ bukan perbedaan maksud antar tahap *(P33)*.

## ⛔ RALAT implementasi 2026-10-03

1. **NBStatus.** Tujuh connector Flow menanam nama orang (`NB IS IN <nama>'S INBOX`). P40: nama dari
   data. Berkas menunggu POSISI, bukan orang (AC 92) ⇒ teksnya memakai **nama posisi tujuan**
   (`NB IS IN REASTREATYINSECHEAD'S INBOX`). Cabang DT yang memang memakai data (`pyUserName`,
   `pxCreateOpName`) memakai nama tampilan dari `M_LOGIN_GO.NAME`. ~~⚠️ Tafsiran — mohon konfirmasi.~~
   ⇒ **Dijawab K5 (PROMPT-NB-TREATY-IN-PUTARAN-2 bab 2, 03-10-2026): Ya** — `NBStatus` memakai nama
   posisi. Uji: `handlers/alur_test.go` TestTanggaPenuhDanNomorPolisSekali
   (`NB IS IN REASTREATYINSECHEAD'S INBOX`), `models/tangga_test.go` TestTeksNBStatus.
2. `HISTORYAKSEPTASIPEGA.ID_PEGA` = `pzInsKey` (`ASM-FW-GISFW-WORK-NB NB-<n>`); `OPERATORID` = identitas
   login; `USERNAME` = nama tampilan; ditulis di transaksi submit (AC 83). Nama tampilan kosong
   menghasilkan kosong — tanpa jatuh-balik ke ID login (AC 40, 42).
3. `HISTORYAKSEPTASIPRODUCTION` tidak ditulis kasus treaty: kedua kalang `SaveViewSuggest` bersyarat
   `Quotation.BusinessFac == "F"`; treaty = "T". Catatan disimpan `T_POLIS_SUGGEST` (tiket 19).

## ⭐ Putaran 2 — paket penyimpanan (03-10-2026)

Dasar: PROMPT-NB-TREATY-IN-PUTARAN-2 bab 0 butir 11–12, bab 2 K4/K16/K17; rincian kolom `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`.

⛔ **RALAT atas RALAT butir 3** di atas. Bunyi lama, dikutip: *"`HISTORYAKSEPTASIPRODUCTION` tidak ditulis
kasus treaty: kedua kalang `SaveViewSuggest` bersyarat `Quotation.BusinessFac == "F"`; treaty = "T". Catatan
disimpan `T_POLIS_SUGGEST` (tiket 19)."*

Bunyi baru (`[keputusan work owner]` **K4**): `T_POLIS_SUGGEST` di luar diagram grilling — **dihapus**.
Catatan `SuggestList` **ditulis ke `POOLDATA.HISTORYAKSEPTASIPRODUCTION`** (pengganti
`SaveViewSuggest` langkah 2 → `RDBList\InsertViewSuggest_SQL`) dan **dibaca balik** untuk grid
`Section\ListSuggest`:

| Kolom | Isi | XML |
| --- | --- | --- |
| `IDPEGA` | `pzInsKey` (`ASM-FW-GISFW-WORK-NB NB-<n>`) | `{pyWorkPage.pzInsKey}` |
| `TYPE_POLIS` | `NB` | `@replaceAll(pyWorkIDPrefix,"-","")` |
| `NOURUT` | berikutnya per `IDPEGA` di bawah kunci kasus | `.pxListSubscript` |
| `POSISI` · `PUTARAN` | `Policy` · `2` | CARI3 · CARI8 kalang "UNTUK TREATY" |
| `PIC` | **nama tampilan** penulis catatan (AC 40, 44; P33) | `.OperatorName` |
| `TGL_INP` | `.Date` catatan | CARI5 |
| `TYPE` | `Quotation.BusinessFac` (`T`) | CARI7 |
| `APPROVAL` | `1` → `Accept`, `0` → `Reject`, selain itu kosong — **per baris** (AC 41, 42) | CARI9 |
| `KETERANGAN` | 3990 karakter pertama `.Suggest` (AC 40) | `substr(CARI10,0,3990)` |
| `AKSES_LOGIN` | **identitas login** penulis (AC 39, 43; P4) | `OperatorID.pyUserIdentifier` |
| `BUSINESS_CODE` | `Quotation.BusinessCode` | |
| `DIV` · `B2B` · `PERCENT_RNM` | NULL | `pyOrgDivision` tanpa sumber di inti; `OfferFacIn` halaman Fac |

`[penyimpangan sadar]` (K4, grilling ID-31/AC 39): syarat `BusinessFac == "F"` tidak ditiru; baris ditulis
pada submit **ketiga** jenjang yang menambahkannya (XML: hanya `InputPolicyTreatyInPost_Act` langkah 4);
`TGL_INP` jam 24 (XML `hh` → `HH24` menyimpan jam sore sebagai pagi). Uji: `models/usulan_test.go`,
`handlers/alur_test.go` (`TestAdminMenolakDiselesaikanDitolak`, `TestTanggaPenuhDanNomorPolisSekali` — lima
baris `NOURUT` 1..5 dengan `AKSES_LOGIN` tiap jenjang), `repository/kolom_test.go`
`TestSQLRiwayatProduksiMengikutiInsertViewSuggest`, `polis_db_test.go` `TestRiwayatProduksiPulangPergi`
(`-tags db`, belum dijalankan — K11).

AC 39, 40, 41, 42, 43, 44, 71: ✅ (seam HTTP + fungsi murni). AC 72 tetap 🟡 (pengurutan lawan Oracle).
Butir terbuka: `DIV` tanpa sumber; tipe fisik `NOURUT` tabel lama belum dicek katalog.

## ⭐ Putaran 2 — P9 (04-10-2026): tiga penyimpangan K4 `[penyimpangan sadar — menunggu konfirmasi WO]`

Dasar: tinjauan spec P9 (temuan 5). Bunyi P1 di atas, dikutip: *"`[penyimpangan sadar]` (K4, grilling ID-31/AC 39):
syarat `BusinessFac == "F"` tidak ditiru; baris ditulis pada submit **ketiga** jenjang yang menambahkannya (XML:
hanya `InputPolicyTreatyInPost_Act` langkah 4); `TGL_INP` jam 24 (XML `hh` → `HH24` menyimpan jam sore sebagai
pagi)."* RALAT penanda: yang berdasar keputusan WO (**K4**, grilling ID-31/AC 39) hanya syarat `BusinessFac`. Tiga
penyimpangan berikut **dipertahankan** tetapi berstatus `[penyimpangan sadar — menunggu konfirmasi WO]`:

| # | XML | Sistem baru | Dasar |
| ---: | --- | --- | --- |
| 1 | `SaveViewSuggest` hanya dari `InputPolicyTreatyInPost_Act` langkah 4 (pasca-submit **admin**) | baris yang ditambahkan pasca DT ditulis di submit **admin, Sec Head, Dept Head** | akibat langsung K4: `SuggestList` tidak punya tempat simpan lain di antara langkah (Pega menyimpan halaman utuh; di sini hanya delapan tabel diagram) — tanpa ini catatan jenjang atasan hilang (AC 71) |
| 2 | `CARI5 = @FormatDateTime(.Date,"dd/MM/yyyy hh:mm:ss",…)` lalu `To_date(…,'DD/MM/YYYY HH24:MI:SS')` — jam sore tersimpan sebagai pagi | `TGL_INP` jam 24 apa adanya | layar (riwayat catatan) membaca balik tabel ini |
| 3 | `CARI2 = .pxListSubscript` | `NOURUT` = MAX+1 per IDPEGA di bawah kunci kasus (`repository.CatatUsulan`) | dua submit serentak tidak berbagi nomor. **Nilainya SAMA dengan `.pxListSubscript`**: `SuggestList` dibangun ulang dari tabel berurut NOURUT dan catatan baru ditambahkan di ujung — uji `models/usulan_test.go` TestCatatanBaruBerposisiNourutBerikut, `handlers/alur_test.go` TestNourutUsulanSamaDenganSubskripSuggestList (lima submit tolak-naik-setuju: baris NOURUT j selalu di pxListSubscript j), `repository/polis_db_test.go` TestRiwayatProduksiPulangPergi (tag db) |

Dicatat juga di PERMINTAAN-TIM-INTI bagian F2. Kode: `backend/models/usulan.go`, `backend/services/tindakan.go`
(`Kirim`). AC tidak berubah status.

## ⭐ Putaran 2 — paket P11 (04-10-2026): uji `db` AC 72

`backend/repository/riwayat_db_test.go` `TestDaftarRiwayatBerurutWaktu` (bertag `db`, **belum dijalankan** — K11):
tiga baris ditulis `CatatRiwayat` (`TGL_TRANSFER = SYSDATE`, padanan `sysdate` `InsertHistoryAkseptasiPega_Sql`),
lalu waktunya digeser sehingga urutan TULIS berlawanan dengan urutan WAKTU (pertama ditulis = terakhir menurut
waktu); `DaftarRiwayat` wajib mengembalikan urutan waktu `2026-10-01 09:00:00`, `2026-10-02 08:30:00`,
`2026-10-03 10:00:00` — pembacaan berurut sisip/ROWID saja gagal. Tabel warisan dipakai bila ada di skema uji,
selain itu tiruan (lihat tiket 08 bab P11).

Status: **selesai** — sisa penahan hanya pihak luar (K11, F2, B4, C8).
