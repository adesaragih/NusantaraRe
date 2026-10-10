# PR — modul baru Disease Life (`diseaselife`, MASTER TREATY) + sequence dan PK DISEASE_LIFE (migrasi modul 080-081, slot 951)

## Summary

```diff
 POOLDATA
   TABLE DISEASE_LIFE (ID VARCHAR2(100), ICD_CODE VARCHAR2(100), DISEASE VARCHAR2(1000))   kolom & nama TETAP
+  SEQUENCE SEQ_DISEASE_LIFE  START WITH ID angka tertinggi + 1 (DEV 197586) NOCACHE NOCYCLE  080
+  CONSTRAINT PK_DISEASE_LIFE PRIMARY KEY (ID)  (+ indeks unik, ID NOT NULL)                   081
+  M_NAV_MENU 'diseaselife' "Disease Life" MASTER TREATY URUTAN 15                              951
 M_DISEASE_LIFE / PEGA_M_DISEASE_LIFE / M_DISEASE_LIFE_SEQ tidak disentuh; PEGA_DISEASE_LIFE tetap VALID;
 baris uji 102051 TEST123 / Sakit dihapus WO (D2) SEBELUM -migrate lewat docs/sql, bukan migrasi
```

```text
080  blok sequence-dari-kueri (bentuk PERSIS 923): ALL_SEQUENCES dulu; awal = NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID))),0)+1
081  blok ALL_CONSTRAINTS (n = 0): ALTER TABLE DISEASE_LIFE ADD CONSTRAINT PK_DISEASE_LIFE PRIMARY KEY (ID)
     berhenti = galat ALTER itu sendiri: ORA-02437 (ID kembar) / ORA-01449 (ID NULL) / ORA-02260 (PK lain) - atomik,
     nol baris dihapus, tanpa kutip; pemaksa UPDATE 944 / 948 tidak dipakai (tidak menggigit ID NULL, menulis 97.586 baris)
951  INSERT datar M_NAV_MENU … WHERE NOT EXISTS
_down 081: blok DROP CONSTRAINT ; 080: DROP SEQUENCE ; 951: buang hak lalu baris
```

Modul (pola causeoflosslife + benefitlife): form Number / ID (disabled) / ICD Code / Disease (textarea), Save, Cancel
(saat Edit), grid 10 baris ID / ICD Code / Disease / Edit, bawaan ID menurun, urut ID / ICD Code, saring ICD Code /
Disease - **saring, urut, dan halaman di server** (97.586 baris). API `/api/disease-life`: GET, POST, PUT. Paritas:
`docs/PARITAS-LAYAR-DAN-AKSI.md`. Urutan WO (P3 di paling atas, D2 sebelum `-migrate`): `docs/LANGKAH-WO-DISEASELIFE.md`.

## Perlu persetujuan tim inti

- **Jatah nomor dipinjam (K0, keputusan work owner 08-10-2026)**: `modul/premiumlistlife/MODUL.md` `Rentang migrasi`
  `050-089` → `050-079` (`080-084` ke modul ini, `085-089` ke `coverlife`); `modul/claimlife/MODUL.md` `Slot menu`
  `950-951` → `950-950` (`951` ke modul ini). Nomor terbukti kosong sebelum ditulis (Glob: nol berkas `065`-`089` dan
  `950`/`951`/`956`/`957` di repo; isi MODUL.md: premiumlistlife terpakai `050–064`, claimlife slot tidak terpakai,
  treatycontractout slot tidak terpakai). Itu SATU-SATUNYA perubahan di folder premiumlistlife dan claimlife.
- **Penjaga menu `TestMenuBersihDuaPuluhBarisSatuPerModul`** (`inti/backend/penjaga/menu_test.go`): URUTAN kini
  diperiksa menurut NILAINYA (sort stabil) - dulu menurut urutan berkas. Sebabnya keputusan K0 / D4: baris URUTAN 15
  lahir di slot **951**, SEBELUM baris URUTAN 14 (`causeoflosslife`) di slot 955. Yang dijaga tetap sama (1, 2, 3, …
  tanpa celah, tanpa kembar per golongan). Ditinjau tim inti.
- `barisLahirDiSlot` + `modulLuarKorpus` (menu_test.go) bertambah `diseaselife`; `TestAturanSlotMenuMenggigit`
  (rentang_test.go) +2 kasus (slot satu modul tidak melahirkan baris modul lahir-di-slot lain).
- `inti/backend/migrasi/perintah_katalog.go` dan `sequence_kueri.go` TIDAK diubah: 080 dan 081 memakai bentuk yang ada.
- Lainnya: `cmd/api/gerbang_tulis_test.go`; `frontend/katalogKorpus.ts`, `daftar.menuTabel.test.ts`, `Shell.test.ts`
  (47); `inti/backend/daftar/modul_diseaselife_gen.go`; skema uji `uji/skemauji/disease_tiruan.go` + `skemauji.go`
  (tiruan `DISEASE_LIFE` tanpa PK SEBELUM pelari, dibongkar SESUDAH migrasi mundur).

## Pembaca lain (minimal)

`claimlife`: kueri pencarian diagnosa TIDAK berubah (nama dan kolom sama; PK hanya menambah indeks unik pada ID). Uji
claimlife tidak bergantung pada "tanpa PK" dan tidak punya tiruan DDL `DISEASE_LIFE` - nol perubahan kode / uji; hanya
`MODUL.md` (jatah K0).

## Evidence

- `go test` modul: models (ID dari sequence, wajib / huruf besar / batas, saring), services (Add 100006, kembar, Edit,
  404, ID terpakai 409, balapan PK, View only, saring / urut / halaman sampai ke gudang), handlers (201 / 200 / 422 / 403 /
  404 / 400 / 409 / 503, nol DELETE, query `icd/disease/urut/arah/halaman`), repository (SQL berhalaman setiap kombinasi,
  080 = bentuk 923 lewat `BacaSequenceDariKueri`, 081 lewat `BacaPerintahKatalog` tanpa kutip, 951, urutan pelari 080 →
  900 → 951 atas migrasi inti SUNGGUHAN, nol tulis data / RENAME / CREATE TABLE). `-tags=db`
  (`TestDBSequenceMulaiSesudahIDTertinggi`, `TestDBPKBerhentiBilaKembarAtauNull` - ORA-02437 / ORA-01449 nyata,
  `TestDBLayanan`, `TestDBUrutanPelari080Lalu900Lalu951`) terkompilasi, BELUM dijalankan. vitest modul + gaya.
- Bukti DEV (`docs/sql/dis_bukti.sql`), acuan, D2 - **menunggu hasil WO**.

## Pertanyaan terbuka (dengan rekomendasi)

1. **`M_DISEASE_LIFE` (3 baris JSON lama, ID 100001-100003 bentrok dengan ID flat tetapi isi beda) +
   `PEGA_M_DISEASE_LIFE` + `M_DISEASE_LIFE_SEQ`** - **rekomendasi: biarkan**, tidak dipakai modul ini maupun claimlife;
   sudah dicadangkan P3. Bila kelak dibuang: keputusan WO terpisah, berkas WO sendiri.
2. **Prosedur `PEGA_DISEASE_LIFE`** (masih VALID): rumus ID-nya `'1' || LPAD(M_DISEASE_LIFE_SEQ.NEXTVAL, 5, '0')`
   menghasilkan ID yang sudah terpakai (102052, 102053, …). Sesudah 081 setiap Add dari Pega DITOLAK PK (ORA-00001) -
   data tidak rusak lagi, tetapi Pega gagal menyimpan. **Rekomendasi: hentikan pemakaian layar Disease di Pega** (simpan
   lewat aplikasi ini); prosedurnya tidak diubah / dibuang tanpa keputusan WO. Edit dari Pega (UPDATE) tetap jalan.
3. **ICD Code wajib** padahal XML tidak mewajibkannya (PARITAS T3) - **rekomendasi: pertahankan** (keputusan D3).

## Merge Danger

**Door:** two-way - 080 / 081 dapat dimundurkan (`_down`), data tidak diubah migrasi. D2 (hapus satu baris uji) one-way
kecuali dari cadangan `DISEASE_LIFE_BARIS_UJI.csv`. **Blast Radius:** claimlife (pencarian diagnosa) membaca nama / kolom
yang sama; Pega Disease berhenti dapat MENAMBAH (PK menolak ID kembar dari rumusnya); premiumlistlife / claimlife
kehilangan nomor yang tidak pernah dipakainya.
