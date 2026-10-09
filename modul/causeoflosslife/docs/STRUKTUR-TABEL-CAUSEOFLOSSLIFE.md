# Struktur tabel — Cause Of Loss Life

Modul `causeoflosslife` (keputusan work owner 08-10-2026 K0-K6, `MODUL.md`). SATU tabel: tabel Pega
`M_CAUSEOFLOSS_LIFE` **berganti nama** menjadi `CAUSEOFLOSS_LIFE` (nama view lama yang dibuang) dan menjadi tabel flat
tanpa JSONDATA - migrasi MODUL 090-092 (K0); nol tabel baru. Nama objek dan kolom ditulis SEKALI di
`backend/repository/coll_tabel.go`.

| Objek | Jenis | Ditulis | Dibaca | Pembaca lain |
| --- | --- | --- | --- | --- |
| `CAUSEOFLOSS_LIFE` | tabel (dulu `M_CAUSEOFLOSS_LIFE`), PK `ID` (`SYS_C008825`, ikut RENAME), 2 kolom, 4 baris DEV | Save Add / Edit | grid, ID terpakai, nama kembar | `masterproductnamelife` (pemilih Cause Of Loss: `ID`, `CAUSEOFLOSS`, `ORDER BY ID ASC`) - nama dan kolom tetap |
| `M_CAUSEOFLOSS_LIFE_SEQ` | sequence warisan (last_number 5 DEV) | NEXTVAL | — | prosedur `PEGA_M_CAUSEOFLOSS_LIFE` (INVALID sesudah 092, K3) |
| `M_CAUSE_OF_LOSS`, `D_CAUSE_OF_LOSS`, `V_M_CAUSE_OF_LOSS`, `V_D_CAUSE_OF_LOSS*`, `T_LISTCAUSEOFLOSS`, `PEGA_M_CAUSE_OF_LOSS`, `PEGA_D_CAUSE_OF_LOSS` + sequence-nya | milik aplikasi lain | — | — | TIDAK disentuh (uji `TestMigrasiTanpaTabelLain`, `TestPeriksaTulis`) |

**ID baru** (rumus prosedur Pega, ditiru persis): `'1' || LPAD(M_CAUSEOFLOSS_LIFE_SEQ.NEXTVAL, 5, '0')` (huruf `1`
tetap; ≤ 6 karakter dari VARCHAR2(10) - nomor 6 angka DITOLAK, LPAD Oracle memotongnya). ID DEV: 100001-100004; ID
berikutnya 100005. NEXTVAL yang menghasilkan ID yang sudah ada = 409 berkalimat (K3). `models.BentukID`, uji
`TestBentukIDRumusProsedur`.

## CAUSEOFLOSS_LIFE

Tabel Pega yang BENTUKNYA diubah migrasi modul 090-092 (karena itu tidak dinyatakan "Tabel warisan" di `MODUL.md`).
Tabel di bab ini = kolom yang DIBUAT migrasi (`091_causeofloss_life_kolom.sql`, `TestKolomDDLCocokDenganStruktur`);
nama = kolom view lama. Diisi 092 dari kunci JSON `CauseofLoss` apa adanya (ekspresi = teks view), diperiksa (K1.4),
lalu JSONDATA dibuang.

| Kolom | Tipe | DDL | Bukti lebar dan isi |
| --- | --- | --- | --- |
| `CAUSEOFLOSS` | teks | VARCHAR2(200) NULLABLE | **[terverifikasi data DEV 08-10-2026]** maks 9 byte, semua huruf besar, tanpa kembar; 100001 = NULL (`"CauseofLoss":""`). Medan "Cause of Loss" `.CauseofLoss` b1027 (pxTextInput b1029, wajib b996, tanpa batas panjang di XML); lebar = pola Benefit 943 (K1) dan `models.BatasCauseOfLoss` |

### Kolom warisan CAUSEOFLOSS_LIFE (tidak dibuat migrasi mana pun)

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(10) NOT NULL, PK `SYS_C008825` (ikut RENAME) | `'1' || LPAD(M_CAUSEOFLOSS_LIFE_SEQ, 5)`; grid "ID" b2989 = `.ID` b3415, form "ID" b823 disabled b872 |
| `JSONDATA` | CLOB, `ENSURE_M_CAUSEOFLOSS_LIFE_JSON` | **DIBUANG 092** (`CASCADE CONSTRAINTS`); kunci `CauseofLoss` di kolom CAUSEOFLOSS; `pxObjClass`, `pyRuleHarness` hanya di cadangan P3 |

## Riwayat — view yang dibuang

- `CAUSEOFLOSS_LIFE` (VIEW): `SELECT a.ID, a.JSONDATA.CauseofLoss FROM M_CAUSEOFLOSS_LIFE a` (kolom `ID`,
  `CAUSEOFLOSS`) - dibuang 090, dibuat ulang hanya oleh 090_down (teks sama, ber-`{skema}`).
