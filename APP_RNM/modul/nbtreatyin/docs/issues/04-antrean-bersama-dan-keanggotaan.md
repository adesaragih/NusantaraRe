# 04: Antrean bersama dan pemeriksaan keanggotaan — bukan nomor urut daftar

**Status:** ready-for-agent
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

- [ ] **AC 11** — setiap penugasan masuk ke **antrean bersama**; ⛔ tidak ada kotak masuk pribadi
- [ ] **AC 14** — keanggotaan antrean diperiksa **menurut nama antrean**, ⛔ bukan menurut nomor urut
- [ ] **AC 92** — berkas menunggu **posisi**, bukan orang

## Perintah verifikasi

1. Tambahkan pengguna ke satu antrean lain, lalu urutkan ulang daftarnya — ⭐ wewenangnya
   **tidak berubah**.
2. Ambil berkas sebagai pengguna kedua yang memegang posisi sama — ⭐ **berhasil**.

## Catatan

⚠️ `[terverifikasi]` Pola penunjukan-menurut-posisi ada di **41 berkas pada 9 modul** di seluruh
korpus. ⭐ Perubahan yang sama berlaku di sana ketika modul itu digarap — **catat, jangan kerjakan
di tiket ini**.
