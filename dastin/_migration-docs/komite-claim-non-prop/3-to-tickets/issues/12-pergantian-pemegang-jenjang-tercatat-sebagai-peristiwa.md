---
status: tertahan
---

# 12: Pergantian pemegang jenjang tercatat sebagai peristiwa

*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Pemilik proses mengganti pemegang sebuah baris roster, atau mencatat
delegasi tetap, dan pergantian itu tersimpan sebagai peristiwa dengan **pemegang sebelum dan
sesudah, pelaku, alasan, dan waktunya**. Auditor dapat menelusuri siapa memegang jenjang apa
pada tanggal berapa tanpa bertanya kepada siapa pun — dan tanpa bergantung pada baris roster
yang sudah tertimpa.

`PERISTIWA_ROSTER` beserta sequence dan constraint-nya; penulisan peristiwa di dalam
transaksi `KelolaRoster` yang sama dengan perubahan barisnya.

**Persyaratan:** `S-052`.

**Tidak termasuk:** **pengalihan sesaat.** `H-4` menetapkan ia tetap tidak diekspos: modelnya
menampung lewat `E-3`, permukaannya tidak ada. Kolom objek ini **tidak kurang** — yang tidak
ada adalah jalur yang menulisnya. Jangan menambahkan kolom yang dikira terlupakan.

**Jalur gagal:** peristiwa tidak dapat ditulis → **penyimpanan roster batal seluruhnya**;
baris roster tidak berubah, dan tidak ada peristiwa separuh yang tertinggal · baris yang
kedua kolom pemegangnya kosong sekaligus → ditolak constraint, karena baris semacam itu tidak
mencatat pergantian apa pun · pemanggil tanpa peran administratif → `403`, nol peristiwa
lahir.

**Uji:** BARU — pergantian pemegang menghasilkan tepat satu peristiwa yang membawa nilai
sebelum **dan** sesudah · BARU — kegagalan penulisan peristiwa membatalkan perubahan
rosternya, diperiksa dengan membandingkan baris roster sebelum dan sesudah · BARU — uji DDL
untuk `CK_PERISTIWA_ROSTER_2` dan untuk ketiadaan kolom pengalihan sesaat.

**Menggantikan:** `CreateChildKomiteCNP_Act`·26.8.3 — substitusi satu orang ke orang lain yang
ditulis **di dalam kode** (`F-19`). Cabangnya dibuang oleh keputusan beku no. 5; kebutuhan
yang ditambalnya digantikan dua hal — delegasi tetap sebagai atribut roster (`S-051`, tiket
`01`) dan pergantiannya tercatat sebagai peristiwa di sini (`S-052`). Di sistem lama nilai
`KomiteID` sebelum ditimpa **hilang begitu ditimpa**; itu yang objek ini kembalikan.

**Blocked by:**

- `01` — Roster dan tabel seleksi berdiri sebagai data yang dikelola


**Dasar:** DECIDED(`E-3`, `H-3`, `H-4`, keputusan beku no. 4, `ADR-0019`, `K5-6`).
EVIDENCED: `F-19` — satu substitusi orang ditulis di dalam rule,
`CreateChildKomiteCNP_Act`·26.8.3.

- [ ] Tiap pergantian pemegang menghasilkan **tepat satu** peristiwa dengan pelaku, alasan,
      dan waktunya.
- [ ] Peristiwa membawa **pemegang sebelum dan sesudah**; nilai sebelum tidak hilang ketika
      baris roster ditimpa.
- [ ] Peristiwa dan perubahan baris roster berada dalam **satu transaksi**; keduanya batal
      bersama.
- [ ] Baris yang kedua kolom pemegangnya kosong sekaligus **tidak dapat ditulis** — ditolak
      constraint, bukan diperiksa kode.
- [ ] Kosong pada kolom pemegang berarti **tidak ada pemegang**, bukan nol; nol kolom di objek
      ini yang kosong dan nolnya dapat tertukar (`K5-6`, `ADR-0019`).
- [ ] Peristiwa roster **tidak ikut terbaca** saat sirkulasi dibaca — ia hidup di seam
      `KelolaRoster`, bukan `BacaSirkulasi`.
- [ ] Nol kolom untuk pengalihan sesaat, dan nol jalur yang menulisnya (`H-4`).
- [ ] Cap waktu berzona `Asia/Jakarta` (aturan kerja `G`); nama di bawah 30 byte, constraint
      berbentuk `<peran>_<tabel>[_n]` yang tidak mengeja kolom.

**Ketidakpastian:** Tidak ada.
