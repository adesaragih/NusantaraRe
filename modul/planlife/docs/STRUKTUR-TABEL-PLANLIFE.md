# Struktur tabel — Plan

Modul `planlife` (keputusan work owner 08-10-2026 K1-K7, `MODUL.md`). SATU tabel: tabel Pega `M_PRODUCT_TYPE_LIFE`
**berganti nama** menjadi `PRODUCT_TYPE_LIFE` (nama view lama yang dibuang), menjadi tabel flat tanpa JSONDATA, dan
mendapat `PK_PRODUCT_TYPE_LIFE (ID)` - migrasi inti 946-948; nol tabel baru. Nama objek dan kolom ditulis SEKALI di
`backend/repository/plnl_tabel.go`.

| Objek | Jenis | Ditulis | Dibaca | Pembaca lain |
| --- | --- | --- | --- | --- |
| `PRODUCT_TYPE_LIFE` | tabel (dulu `M_PRODUCT_TYPE_LIFE`, tanpa PK), PK `PK_PRODUCT_TYPE_LIFE` sejak 948, 6 kolom, 32 baris DEV | Save Add / Edit | grid, Plan Name kembar, ID terpakai | `masterproductnamelife` (autocomplete `Plan Name`, `SELECT ID, COVERNAME, BUSINESS, BENEFIT` - kolom sama) |
| `BUSINESS` | tabel warisan (dibaca juga `edmtreatyin`, `mastercontractretrolife`) | — | pilihan Business `GROUPPANEL = '009'` | TIDAK diubah |
| `BENEFIT_LIFE` | tabel modul `benefitlife` | — | pilihan Benefit | TIDAK diubah |
| `M_PRODUCT_TYPE_LIFE_SEQ` | sequence warisan (last_number 44 DEV) | NEXTVAL | — | prosedur `PEGA_M_PRODUCT_TYPE_LIFE` (INVALID sesudah 948, K3) |
| `PLAN_LIFE_SUMMARY`, `PEGA_M_PLAN_LIFE_SUMMARY`, `M_BUSINESS` | milik objek lain | — | — | TIDAK disentuh |

**ID baru**: `'1' || LPAD(M_PRODUCT_TYPE_LIFE_SEQ.NEXTVAL, 5, '0')` (≤ 6 karakter = VARCHAR2(6); nomor 6 angka
DITOLAK). NEXTVAL yang menghasilkan ID yang sudah ada = 409 berkalimat (K3).

## PRODUCT_TYPE_LIFE

Tabel di bab ini = kolom yang DIBUAT migrasi (`947_product_type_life_kolom.sql`, `TestKolomDDLCocokDenganStruktur`);
nama = kolom view lama, urutan sama. Diisi 948 dari JSONDATA apa adanya (ekspresi = teks view), diperiksa (K1.4 / K2).

| Kolom | Tipe | DDL | Bukti lebar dan isi |
| --- | --- | --- | --- |
| `COVERNAME` | teks | VARCHAR2(200) | **[terverifikasi data DEV 08-10-2026]** maks 34 byte, unik. "Plan Name" `.CoverName` b841 (pxTextInput b844, tanpa batas); lebar = pola benefitlife 943 |
| `BUSINESS` | teks | VARCHAR2(200) | maks 29 byte; "Business" `.Business` b1105 (pxAutoComplete) = `BUSINESS.NOTE` (b1187 / b1231) di 32/32 baris |
| `BUSINESSID` | teks | VARCHAR2(10) | `BUSINESS.ID` (`pyPropertyTarget` .BusinessID b1268); lebar = master (bukti `plan_bukti_k1.sql` E) |
| `BENEFIT` | teks | VARCHAR2(200) | maks 48 byte; "Benefit" `.Benefit` b1461 = `BENEFIT_LIFE.BENEFIT` (b1541 / b1552) |
| `BENEFITID` | teks | VARCHAR2(10) | `BENEFIT_LIFE.ID` (`.Number` → .BenefitID b1563-b1564) = VARCHAR2(10) |

### Kolom warisan PRODUCT_TYPE_LIFE (tidak dibuat migrasi mana pun)

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(6), dulu NULLABLE tanpa PK; sejak 948 PK `PK_PRODUCT_TYPE_LIFE` (NOT NULL) | `'1' || LPAD(seq, 5)`; = `JSONDATA.ID` di 32/32 baris (948 BERHENTI bila tidak) |
| `JSONDATA` | CLOB, `ENSURE_M_PRODUCT_TYPE_LIFE` | **DIBUANG 948** (`CASCADE CONSTRAINTS`); kunci `px*` / `pyRuleHarness` hanya di cadangan P3; `Note` tidak pernah tersimpan |

## Riwayat — view yang dibuang

- `PRODUCT_TYPE_LIFE` (VIEW): `SELECT a.JSONDATA.ID, a.JSONDATA.CoverName, a.JSONDATA.Business, a.JSONDATA.BusinessID,
  a.JSONDATA.Benefit, a.JSONDATA.BenefitID FROM M_PRODUCT_TYPE_LIFE a` - dibuang 946, dibuat ulang hanya oleh 946_down.
