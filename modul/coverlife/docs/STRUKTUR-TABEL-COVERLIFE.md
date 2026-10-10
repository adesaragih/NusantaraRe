# Struktur tabel — Cover Life

Modul `coverlife` (keputusan work owner 08-10-2026 K0, C1-C4, K5, `MODUL.md`). SATU tabel: tabel Pega `M_COVER_LIFE`
**dijadikan flat di tempat - nama TETAP, TANPA RENAME** - migrasi MODUL 085-086 (K0); JSONDATA dan view `COVER_LIFE`
dibuang; nol tabel baru. Nama objek dan kolom ditulis SEKALI di `backend/repository/cvl_tabel.go`.

| Objek | Jenis | Ditulis | Dibaca | Pembaca lain |
| --- | --- | --- | --- | --- |
| `M_COVER_LIFE` | tabel, PK `ID` (`SYS_C009203`, tetap), 3 kolom sesudah 086, 4 baris DEV | Save Add / Edit | grid, ID terpakai, Cover kembar | — (tidak ada pembaca lain di repo) |
| `M_COVER_LIFE_SEQ` | sequence warisan (last_number 5 DEV) | NEXTVAL | — | prosedur `PEGA_M_COVER_LIFE` (INVALID sesudah 086, C2) |
| `COVER_LIFE` | VIEW lama | — | — | **DIBUANG 086** (terakhir); kode modul tidak menyebutnya (`TestKodeTidakMenyebutViewLama`) |
| `COVERAGE`, `COVERAGE_FACIN`, `COVERAGETRAVEL*`, `COVERNOTE*` | milik aplikasi lain | — | — | TIDAK disentuh (uji `TestMigrasiTanpaTabelLainDanTanpaRename`, `TestPeriksaTulis`) |

**ID baru** (rumus prosedur Pega, ditiru persis, C2): `'1' || LPAD(M_COVER_LIFE_SEQ.NEXTVAL, 5, '0')` (huruf `1` tetap;
≤ 6 karakter dari VARCHAR2(10) - nomor 6 angka DITOLAK, LPAD Oracle memotongnya). ID DEV: 100001-100004; ID berikutnya
100005 (belum terpakai). NEXTVAL yang menghasilkan ID yang sudah ada = 409 berkalimat. `models.BentukID`, uji
`TestBentukIDRumusProsedur`.

## M_COVER_LIFE

Tabel Pega yang BENTUKNYA diubah migrasi modul 085-086 (karena itu tidak dinyatakan "Tabel warisan" di `MODUL.md`).
Tabel di bab ini = kolom yang DIBUAT migrasi (`085_cover_life_kolom.sql`, `TestKolomDDLCocokDenganStruktur`); nama =
kolom view lama, urutan sama. Diisi 086 dari kunci JSON `Cover` / `Note` apa adanya (ekspresi = teks view), diperiksa
(C1.3), lalu JSONDATA dan view dibuang.

| Kolom | Tipe | DDL | Bukti lebar dan isi |
| --- | --- | --- | --- |
| `COVER` | teks | VARCHAR2(200) NULLABLE | **[fakta WO 08-10-2026]** maks 17 byte (PERSONAL ACCIDENT), 4/4 terisi, huruf besar. Medan "Cover" `.Cover` b877 (pxTextInput b880, wajib b843 / b895, tanpa batas panjang di XML); lebar = C1.1 dan `models.BatasCover` |
| `NOTE` | teks | VARCHAR2(1000) NULLABLE | **[fakta WO 08-10-2026]** NULL di 4/4 baris (kunci `Note` tidak ada di JSONDATA; view tetap mengeluarkannya). Medan "Note" `.Note` b1055 (pxTextArea b1058, tidak wajib b1020 / b1070, tanpa pyMaxChars) → 1000 (C1.1 "lebar dari XML atau 1000") dan `models.BatasNote` |

### Kolom warisan M_COVER_LIFE (tidak dibuat migrasi mana pun)

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(10) NOT NULL, PK `SYS_C009203` (tetap) | `'1' || LPAD(M_COVER_LIFE_SEQ, 5)`; grid "ID" b3023 = `.ID` b3449; form TANPA medan ID |
| `JSONDATA` | CLOB, `ENSURE_M_COVER_LIFE_JSON` | **DIBUANG 086** (`CASCADE CONSTRAINTS`); kunci `Cover` di kolom COVER, `Note` (tidak ada di DEV) di NOTE; `pxObjClass`, `pyRuleHarness` hanya di cadangan P3 |

## Riwayat — view yang dibuang

- `COVER_LIFE` (VIEW): `SELECT ID, a.JSONDATA.Cover, a.JSONDATA.Note FROM M_COVER_LIFE a` (kolom `ID`, `COVER`, `NOTE`)
  - dibuang 086 (blok terakhir, berpelindung ALL_VIEWS), dibuat ulang hanya oleh 086_down (teks sama, ber-`{skema}`).
