# 06: Loss allocation dan spreading — share wajib tepat 100 %

**Status:** sebagian 07-10-2026 — loss allocation dan spreading dibangun; tambah / hapus baris Spreading Claim nonaktif, termasuk ikon grid bawaan (keputusan work owner 07-10-2026); kasus baru tanpa baris spreading ditolak Save to issue RNM (`ProteksiData_act` 5) — RALAT 07-10-2026 (semula `ready-for-agent`)
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi klaim dan nomor polis) · 05 (Insured Interest / TSI)
**Menutup:** AC 29 · 30 · 31 · 32 · 33 *(5 AC)* — US 20–23, 30–32

## Hasil & nilai pengguna

Nilai kerugian terbagi ke beberapa treaty menurut persentase share, disegmentasi per mata uang, dan
**tidak ada nilai yang hilang atau terhitung dua kali** — total share ditolak baik saat kurang dari
100 % maupun lebih. Finance mendapat jaminan bahwa jumlah seluruh baris spreading sama persis dengan
nilai induknya.

## Area codebase

Entitas loss allocation · spreading tingkat klaim · spreading quota share · validasi total share.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/CountSpreading_Act.xml` | 6.2.3 | **satu-satunya** pemeriksaan total share di 329 berkas — hanya menjaga sisi **lebih dari** 100 |
| `Activity/SetTreatyNameSpreading_Act.xml` | 10 | ⚠️ step **di-remark** → tidak ditulis sama sekali |
| `Activity/AddAdjustment_Act.xml` | 5 | loss allocation disalin sebagai **snapshot** ke baris adjustment saat baris dibuat |

⚠️ `[terverifikasi]` Enam properti spreading berbagi satu class di Pega; di Oracle mereka **tabel
berbeda dengan peran berbeda** (lihat tiket 00, AC 3). Jangan disatukan hanya karena class-nya sama.

## ADR terkait

**ADR-0003** (uang non-float) · **ADR-0011** (unit keputusan = baris `AdjustmentList`).

## Acceptance criteria

- [ ] `[terverifikasi]` Loss allocation membagi nilai kerugian ke **treaty**, disegmentasi **per mata uang** — tanpa dimensi tahun maupun coverage *(AC 29 spec)*
- [ ] ⚠️ `[keputusan work owner]` **Total share wajib tepat 100 %**, ditolak pada kurang dari 100 % **dan** lebih dari 100 %. **Alasan menyimpang:** Pega hanya menjaga sebelah — pemeriksaan tunggalnya tidak pernah menangkap kurang dari 100 *(AC 30 spec)*
- [ ] `[terverifikasi]` **Alokasi sisa pembulatan tidak diperlukan** — karena total 100 % dan presisi disimpan penuh, penjumlahan baris spreading tepat sampai digit terakhir *(AC 31 spec)*
- [ ] Spreading di tingkat klaim terpisah dari spreading pada baris adjustment *(AC 32 spec)*
- [ ] Loss allocation disalin sebagai **snapshot** ke baris adjustment saat baris dibuat *(AC 33 spec)*

## Perintah verifikasi

```
jalankan test "total share 99,99 % -> DITOLAK"
jalankan test "total share 100,01 % -> DITOLAK"
jalankan test "total share tepat 100 % -> diterima"
jalankan test "dua mata uang -> masing-masing disegmentasi dan masing-masing 100 %"
jalankan test "jumlah baris spreading == nilai induk sampai digit terakhir"
jalankan test "snapshot loss allocation pada baris adjustment tidak berubah saat sumbernya berubah"
```
