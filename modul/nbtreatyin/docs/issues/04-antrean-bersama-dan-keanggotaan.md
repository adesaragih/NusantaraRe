# 04: Antrean bersama dan pemeriksaan keanggotaan — bukan nomor urut daftar

**Status:** selesai *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ready-for-agent)*
**Blocked by:** 03
**Menutup:** AC 11 · 14 · 92 *(3 AC)* — US 3 · 8 · 19

## Hasil & nilai pengguna

Hari ini pekerjaan menunggu di **antrean bersama**, dan itu memang yang diinginkan. ⛔ Tetapi
wewenang di beberapa tempat ditentukan dengan bertanya *"antrean nomor dua Anda namanya apa"* —
`[terverifikasi]` menunjuk antrean **menurut posisi dalam daftar**, bukan menurut namanya.
⚠️ Menambah seorang pengguna ke antrean baru **mengubah wewenangnya** tanpa ada yang menyentuh
aturan.

Sesudah tiket ini, pekerjaan tetap menunggu di antrean bersama, dan ⭐ pemeriksaan wewenang
bertanya **"apakah pengguna ini anggota antrean X"** — jawabannya tidak berubah ketika daftar
diurutkan ulang.

## Area codebase

- Lapisan service: penugasan ke antrean
- Lapisan service: pemeriksaan keanggotaan antrean

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Antrean bersama | `Flow\InputRealizationTreatyIn.xml` — **enam** Assignment memakai `ToWorkBasket`, **nol** `ToWorklist` |
| ⛔ Penunjukan menurut posisi | `Section\SFAPortal_OpportunitiesList.xml` *(posisi 2)* · `DataTransform\InputPolicyTreatyIn_preDT.xml` *(posisi 2)* · `When\IsUW.xml` *(posisi 1)* |

## ADR terkait

- **ADR-0002** — RBAC memakai peran yang sudah ada

## Acceptance criteria

- [x] **AC 11** — setiap penugasan masuk ke **antrean bersama**; ⛔ tidak ada kotak masuk pribadi
- [x] **AC 14** — keanggotaan antrean diperiksa **menurut nama antrean**, ⛔ bukan menurut nomor urut
- [x] **AC 92** — berkas menunggu **posisi**, bukan orang

## Perintah verifikasi

1. Tambahkan pengguna ke satu antrean lain, lalu urutkan ulang daftarnya — ⭐ wewenangnya
   **tidak berubah**.
2. Ambil berkas sebagai pengguna kedua yang memegang posisi sama — ⭐ **berhasil**.

## Catatan

⚠️ `[terverifikasi]` Pola penunjukan-menurut-posisi ada di **41 berkas pada 9 modul** di seluruh
korpus. ⭐ Perubahan yang sama berlaku di sana ketika modul itu digarap — **catat, jangan kerjakan
di tiket ini**.

## Hasil implementasi 2026-10-03

Keanggotaan diperiksa menurut NAMA workbasket (`inti.Pelaku.Peran`) — `services.anggota`. Penunjukan
menurut nomor urut tidak dibangun: `InputPolicyTreatyIn_preDT` langkah 12-13 (`pyWorkBasketList(2)`),
`When\IsUW` (`pyWorkBasketList(1)` = ReasFacIn* — tidak pernah benar bagi antrean treaty).

## ⭐ Putaran 2 — paket P8: gerbang daftar portal (03-10-2026)

Dasar: PROMPT-NB-TREATY-IN-PUTARAN-2 bab 0 butir 4–9, bab 2 K12; temuan audit P6.

Hasil implementasi 2026-10-03 di atas tidak menyebut tempat penunjukan-menurut-posisi yang ketiga
(`Section\SFAPortal_OpportunitiesList.xml`, posisi 2, tabel *Rule Pega sumber*). Bunyi lama, dikutip:
*"Penunjukan menurut nomor urut tidak dibangun: `InputPolicyTreatyIn_preDT` langkah 12-13
(`pyWorkBasketList(2)`), `When\IsUW` (`pyWorkBasketList(1)` = ReasFacIn* — tidak pernah benar bagi antrean
treaty)."* Bunyi itu tetap benar untuk kedua rule tersebut; tempat ketiga kini **dibangun menurut nama**:

| XML (`Section\SFAPortal_OpportunitiesList.xml`) | Dibangun |
| --- | --- |
| baris saringan `.FilterTermForOpportunity` + tombol Filter — tampil selalu | `PortalNBTreatyIn.tsx` (tetap) |
| wadah `pyContainerVisibleWhen` `OperatorID.pyWorkGroup!='ReasLife' && OperatorID.pyWorkBasketList(2).pyWorkBasketName=='ReasTreatyInAdmin'` ⊃ wadah `!IsOperatorLife` ⊃ SATU grid `pyGridProps/pyRDName = GetListOpportunity` | `services.DaftarKasus`: anggota `ReasTreatyInAdmin` **menurut nama** (AC 14) melihat grid — semua kasus terbuka, opsional per posisi |
| Sec Head / Dept Head: grid tersembunyi; tugasnya dirutekan `Flow\InputRealizationTreatyIn` ke workbasket (Assignment4/6 `ReasTreatyInSecHead`, Assignment3 `ReasTreatyInDeptHead`, `ToWorkBasket`). ⛔ Daftar kerja workbasket Pega **tidak ada di korpus** | portal tunggal (bab 0 butir 7): **hanya** kasus yang menunggu di posisi tangga yang dipegang pelaku (AC 11, 92); `?posisi=` di luar posisinya ⇒ 403 |
| pelaku tanpa posisi tangga (termasuk GroupLeader/Director yang dibuang, AC 10) | 403 `ErrBukanAnggotaAntrean` |
| klausa `pyWorkGroup!='ReasLife'` (= `!IsOperatorLife`) | ⛔ **tidak dibangun** — work group Pega tak berpadanan di `inti.Pelaku` (AkunID + workbasket) maupun M_LOGIN_GO; K12 kosong ⇒ butir terbuka di tiket 05 |

Perubahan perilaku sejak putaran 1: daftar portal semula terbuka bagi **setiap** pelaku beridentitas (semua
kasus terbuka); kini mengikuti wadah XML di atas. Membuka kasus lewat ID (`GET /kasus/{id}`, hanya-baca bagi
bukan anggota) tidak berubah.

Uji seam HTTP: `backend/handlers/portal_test.go` `TestGerbangDaftarPortal` (admin di urutan 1 dan 2, Sec Head,
Dept Head, keduanya, saring posisi, cari, bukan anggota, GroupLeader, tanpa peran);
`backend/handlers/rute_test.go` `TestAlurHTTP` (daftar tanpa antrean ⇒ 403);
`backend/repository/kolom_test.go` `TestSQLDaftarKasusPenampungUnik` (`POSITION_NOTE IN (...)`, penampung unik).

## ⭐ Putaran 2 — P9 (04-10-2026): pencarian dan batas baris ikut RD `GetListOpportunity`

Dasar: temuan paket P8 (*"BATAS BARIS 200 vs XML 500 & PENCARIAN lebih luas dari XML"*), tinjauan P9. XML
`ReportDefinition\GetListOpportunity.xml`: logika `B AND D AND E AND F AND A AND (C OR G) AND F1`;
C = `.Name Contains Param.Search`, G = `.TextNoQuotation Contains Param.Search` (keduanya `pyCaseInsensitive=true`);
`pyMaxRecords` **500**.

| Hal | Semula (putaran 1) | Kini |
| --- | --- | --- |
| Medan pencarian | pengenal kasus ATAU nama bisnis ATAU nama tertanggung | hanya filter **G** = pengenal kasus `NB-<n>` (`.TextNoQuotation`, filter B `Contains "NB-"`), tanpa beda huruf besar/kecil (`models.CocokCariPortal`, `repository.sqlDaftarKasus`) |
| Filter C `.Name` | — | ⛔ tidak dibangun: `.Name` milik kelas CRM `ASM-FW-SFAGISFW-Work-Opportunity`, ditulis NOL rule korpus dan tak berkolom di diagram grilling (bab 0 butir 11) — butir terbuka |
| Batas baris | 200 | **500** (`models.BatasDaftarPortal`) |

RALAT bunyi uji di bab P8 di atas, dikutip: *"`backend/repository/kolom_test.go` `TestSQLDaftarKasusPenampungUnik`
(`POSITION_NOTE IN (...)`, penampung unik)"* → kini `TestSQLDaftarKasusPenampungSamaDenganArgumen` (perilaku
pengikatan: penampung unik = argumen, pencarian mengikat satu nilai; bukan teks SQL). Uji seam HTTP:
`backend/handlers/portal_test.go` TestDaftarPortalSesuaiGetListOpportunity (cari huruf kecil menemukan pengenal;
nama bisnis/tertanggung tidak; 501 kasus → 500 baris); uji db `repository/polis_db_test.go`
TestDaftarKasusPortalMenurutGetListOpportunity (ditulis, tidak dijalankan — K11).
