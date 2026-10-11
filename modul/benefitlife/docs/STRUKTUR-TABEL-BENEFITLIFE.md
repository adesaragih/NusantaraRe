# Struktur tabel — Benefit

Modul `benefitlife` (keputusan work owner 08-10-2026 K1-K6, `MODUL.md`). SATU tabel: tabel Pega `M_BENEFIT_LIFE`
**berganti nama** menjadi `BENEFIT_LIFE` (nama view lama yang dibuang) dan menjadi tabel flat tanpa JSONDATA - migrasi
inti 942-944; nol tabel baru. Nama objek dan kolom ditulis SEKALI di `backend/repository/bnfl_tabel.go`.

| Objek | Jenis | Ditulis | Dibaca | Pembaca lain |
| --- | --- | --- | --- | --- |
| `BENEFIT_LIFE` | tabel (dulu `M_BENEFIT_LIFE`), PK `ID` (`SYS_C009031`, ikut RENAME), 2 kolom, 10 baris DEV | Save Add / Edit | grid, ID terpakai | tidak ada (fakta WO 08-10-2026; medan `Benefit` di `masterproductnamelife` adalah teks bebas) |
| `M_BENEFIT_LIFE_SEQ` | sequence warisan (last_number 12 DEV) | NEXTVAL | — | prosedur `PEGA_M_BENEFIT_LIFE` (INVALID sesudah 944, K3) |
| `M_BENEFIT`, `MBENEFIT*`, `VJ_*BENEFIT*`, `VH_*BENEFIT*`, `HCD_*BENEFIT*`, `T_BENEFIT`, `DET_GROUP_BENEFIT`, `M_PLAN_BENEFIT`, `M_PROPERTY_BENEFIT` | milik aplikasi lain | — | — | TIDAK disentuh (uji `TestMigrasiBenefitTanpaTabelLain`, `TestPeriksaTulis`) |

**ID baru** (rumus prosedur Pega, ditiru persis): `'1' || LPAD(M_BENEFIT_LIFE_SEQ.NEXTVAL, 5, '0')` (huruf `1` tetap,
bukan id_site; ≤ 6 karakter dari VARCHAR2(10) - nomor 6 angka DITOLAK, LPAD Oracle memotongnya). ID DEV: 100001,
100002, 100004-100011. NEXTVAL yang menghasilkan ID yang sudah ada = 409 berkalimat (K3). `models.BentukID`, uji
`TestBentukIDRumusProsedur`.

## BENEFIT_LIFE

Tabel Pega yang BENTUKNYA diubah migrasi inti 942-944 (karena itu tidak dinyatakan "Tabel warisan" di `MODUL.md`).
Tabel di bab ini = kolom yang DIBUAT migrasi (`943_benefit_life_kolom.sql`, `TestKolomDDLCocokDenganStruktur`); nama =
kolom view lama. Diisi 944 dari kunci JSON `Benefit` apa adanya (ekspresi = teks view), diperiksa (K2), lalu JSONDATA
dibuang.

| Kolom | Tipe | DDL | Bukti lebar dan isi |
| --- | --- | --- | --- |
| `BENEFIT` | teks | VARCHAR2(200) | **[terverifikasi data DEV 08-10-2026]** panjang maksimum 48 byte, ada di 10 dari 10 baris. Medan "Benefit" `.Benefit` b1172 (pxTextArea b1175, wajib b1137, `SetUpperCase_DT` b1211, tanpa batas panjang di XML); lebar = pola ririsklife `USEDBY` 936 (K1) dan `models.BatasBenefit` |

### Kolom warisan BENEFIT_LIFE (tidak dibuat migrasi mana pun)

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(10) NOT NULL, PK `SYS_C009031` (ikut RENAME) | `'1' || LPAD(M_BENEFIT_LIFE_SEQ, 5)`; grid "ID" b3431 = `.Number` b3858, form "Number / ID" b964 disabled b1013 |
| `JSONDATA` | CLOB, `ENSURE_M_BENEFIT_LIFE_JSON` | **DIBUANG 944** (`CASCADE CONSTRAINTS`); kunci `Benefit` di kolom BENEFIT; `Number` (5 baris, = ID), `pxCreateDateTime`, `pxCreateOperator`, `pxCreateOpName`, `pxCreateSystemID`, `pxObjClass`, `pyRuleHarness` hanya di cadangan P3 |

## Riwayat — view yang dibuang

- `BENEFIT_LIFE` (VIEW): `SELECT a.ID, a.JSONDATA.Benefit FROM M_BENEFIT_LIFE a` (kolom `ID`, `BENEFIT`) - dibuang 942,
  dibuat ulang hanya oleh 942_down (teks sama, ber-`{skema}`).
