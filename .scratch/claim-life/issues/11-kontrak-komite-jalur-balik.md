# 11: Terima & tampilkan hasil keputusan Komite

> ⚠️ **Cakupan diubah 2026-09-15.** Judul lama: *"Kontrak Komite — jalur balik hasil keputusan"*.
> Tiket ini **tidak lagi menulis `STS_REJECT`**. `[terverifikasi]` Penulisannya milik **Komite Claim
> Life**, bukan Claim — Life: `Komite Claim Life/Activity/KomitePostAdjustment.xml`
> (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`) yang menulis ke
> dua tingkat baris. Implementasinya ada di **tiket Komite 05**.
> Tiket ini kini **membaca dan menampilkan** hasil itu, lalu melanjutkan siklus klaim.

**Status:** ready-for-agent

**Blocked by:** 10 (kontrak Komite — penyerahan) · **Komite 05** (jalur balik — penulisan
`STS_REJECT`)

## Hasil & nilai pengguna

Sebagai **ReasLifeSPV**, saya melihat hasil keputusan Komite pada baris adjustment saya — sehingga
saya tahu apakah baris itu diaksep atau ditolak — dan bila ditolak, saya dapat menambah baris
adjustment baru untuk mengajukan ulang dengan angka yang diperbaiki.
*(User story 20, 21, 22 di spec)*

Inilah yang membuat **klaim tidak terminal**: penolakan menghasilkan putaran berikutnya, bukan akhir.

## Area codebase

`internal/services` (pembacaan hasil; pembuatan baris lanjutan; perhitungan status klaim turunan),
`internal/handlers` (endpoint status klaim + tambah baris), `frontend/` (tampilan hasil dan riwayat
putaran).

**Tidak** menulis `STS_REJECT` — itu milik Komite 05.

## Rule Pega sumber

| Rule | Identitas | Peran di tiket ini |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` **penulis** hasil — 6 `Property-Set` bernilai `1`, 2 bernilai `2`, digerbangi `KomiteCount == KomiteLoop`. Tiket ini **membacanya**, tidak menjalankannya |
| `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `[terverifikasi]` **satu rule bersama** — hash ternormalisasi `c50bfd9a12` identik di kedua modul |
| `Claim Life/Activity/SetIndexAdjustmentList.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SETINDEXADJUSTMENTLIST` / `RULE-OBJ-ACTIVITY` | pewarisan 8 kolom ke baris lanjutan |

## ADR terkait

**ADR-0001** (Kontrak 2 — jalur balik; `AcceptStatus` dipetakan ke `STS_REJECT` **di batas**, tidak
disimpan sebagai status kedua), **ADR-0011** (terminal per baris; klaim tidak terminal; revisi =
baris baru), **ADR-0007** (jejak audit).

## Acceptance criteria

- [ ] Hasil keputusan Komite **terbaca** pada baris `AdjustmentList` yang diserahkan — Aksep atau
      Ditolak. *(AC 4 spec Claim Life)*
- [ ] Klaim **tetap** dapat menerima baris adjustment baru setelah penolakan Komite.
- [ ] Baris baru yang ditambahkan setelah penolakan berstatus Outstanding dan **mewarisi delapan
      kolom** dari baris pertama **tanpa** mewarisi status. *(AC 5 spec Claim Life)*
- [ ] `PremiumListDetail` dan header klaim **mencerminkan** baris terakhir setelah hasil diterapkan.
      *(AC 7 spec Claim Life)*
- [ ] Status klaim "selesai" dihitung sebagai keadaan **turunan** dari kumpulan baris — bukan kolom
      tersimpan.
- [ ] `AcceptStatus` **tidak** disimpan sebagai status kedua di konteks ini. *(**ADR-0001**)*
- [ ] Hasil keputusan **hanya diterapkan** ketika putaran Komite sudah mencapai **tingkat terakhir**;
      hasil dari tingkat antara **tidak** mengubah status baris mana pun di konteks ini.
      *(AC 6 spec)* ⚠️ **Penegakannya milik Komite Claim Life tiket 05** — tiket ini hanya wajib
      **tidak menerapkan lebih awal**. Dicatat agar AC 6 punya jejak pemilik, bukan tampak terlewat.
- [ ] Tiket ini **tidak** menulis `STS_REJECT` — diverifikasi dengan tidak adanya jalur tulis status
      baris di konteks Claim — Life.

## Catatan — mengapa cakupan diubah

`[terverifikasi]` Penulisan `STS_REJECT` dilakukan rule di modul **`Komite Claim Life`** dengan class
`ASM-FW-GCNMFW-Work-KomiteLife`. Menempatkannya di tiket Claim — Life akan membuat **dua tiket
`ready-for-agent` sama-sama mengklaim penulisan yang sama** — dua agent dapat mengimplementasikannya
berdua.

`[keputusan work owner 2026-09-15]` Kepemilikan ditetapkan: **Komite menulis, Claim Life membaca.**

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
