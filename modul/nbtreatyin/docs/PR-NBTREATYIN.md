# nbtreatyin: modul NB Treaty In (realisasi treaty masuk, Prop + NonProp)

> Urutan merge: **(1)** PR inti `inti/nbtreatyin-uji-menu` (`PR-INTI-UJI-MENU.md`) → **(2)** migrasi `nbfacin`
> `182_t_general_polis` masuk `dev` → **(3)** PR ini. Migrasi dijalankan work owner, bukan oleh merge.

## Summary

Satu halaman portal (`nbtreatyin-portal`, slot menu 968), satu layar kasus, tangga tiga jenjang Admin → Sec Head →
Dept Head, logika diambil dari korpus XML Pega `NB Treaty In`.

```text
modul/nbtreatyin/
├── backend/
│   ├── models/        # port activity/DT/When XML (fungsi murni, apd, nol float)
│   ├── repository/    # satu-satunya SQL; JSON master hanya MasterXOLDariJSON (baca-saja, K8)
│   ├── services/      # tangga, kunci kasus, satu transaksi per aksi
│   ├── handlers/      # /api/nb-treaty-in — seam uji utama
│   ├── alat/pemuatlama/   # pemuat dokumen lama (dijalankan manusia)
│   └── migrations/    # 320 ALTER T_GENERAL_POLIS (bersama FacIn) · 321–327 CREATE T_POLIS_* · 968 menu
├── frontend/          # portal + LayarKasus + popup pilih bisnis / SOB / subsection NonProp
└── docs/              # spec, tiket 00–23, HASIL, INVENTARIS-XML, AUDIT-SILANG, PERMINTAAN
```

```text
buka kasus → pra-proses (hari tutup buku dari TANGGAL_CLOSING)
  pilih bisnis  (view TREATYINDETAILJOINEDM; NonProp: + master XOL baca-saja)
  hitung        (urutan action set XML; server menghitung ulang medan terkunci)
  simpan / submit admin → Sec Head → Dept Head (K2: selalu naik)
    Dept Head: nomor polis (inti/backend/penomor) → selesai → konversi
  catatan usulan → POOLDATA.HISTORYAKSEPTASIPRODUCTION (K4)
```

Tabel: **7 `CREATE TABLE`** (`T_POLIS_QUOTATION`, `_CEDING`, `_INSTALMENT`, `_INSTALMENT_DETAIL`, `_SPREADING`,
`_XOL`, `_XOL_LAYER`) + **`ALTER TABLE T_GENERAL_POLIS ADD`** (tabel bersama FacIn/Treaty, keputusan WO 04-10-2026).
Nol tabel lain. Daftar kolom lawan diagram: `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`.

## Evidence

Dari akar repo, lawan baseline `origin/dev`:

| Perintah | Before (`origin/dev`) | After (cabang ini) |
| --- | --- | --- |
| `go vet ./...` | bersih | bersih |
| `go test ./...` | gagal 3 paket claimlife (baseline) | gagal **3 paket claimlife yang sama**; 74 paket ok |
| `npm run typecheck` | 0 galat | 0 galat |
| `npm test` | 14 gagal / 947 lulus | **14 gagal yang sama** / 1025 lulus |
| `go test -tags=db ./modul/nbtreatyin/...` | — | **belum dijalankan** (K11: tanpa skema uji) |

AC: `spec.md` ✅ 80 · 🟡 10 · ⛔ 2 · 📄 4 (96); spec penyimpanan ✅ 53 · 🟡 12 · 📄 1 (66). Rincian dan bukti per AC:
`docs/HASIL-IMPLEMENTASI.md`. Audit silang independen: `docs/AUDIT-SILANG-PUTARAN-3.md`.

Contoh uji yang mengunci keputusan:

```text
TestSubmitTanpaApprovalDitolakDiSetiapJenjang     → 422 (K6)
TestDaftarPortal… baris LINI='FAC' tidak tampil   → FacIn tersaring (tabel bersama)
TestPilihBisnisDiLuarDaftarPopupDitolak           → 422
migrasi_test: 320 nol CREATE TABLE; ALTER + ID/IDPEGA = katalog; _down tanpa kolom FacIn
kolom_test: tepat 7 CREATE TABLE + ALTER 320, kolom = diagram
```

## Merge Danger

**Door:** two-way untuk kode (modul baru, rute baru, menu 968 sudah tercatat di `POOLDATA`); **one-way** untuk
migrasi 320–327 bila dijalankan — 320 mengubah tabel bersama `T_GENERAL_POLIS` milik FacIn.

**Blast Radius:** basis-data-bersama.

Penahan sebelum migrasi dijalankan WO:

- **K18 / C10:** migrasi `nbfacin` 182 wajib di `dev` lebih dulu; tanpa itu 320 gagal ORA-00942 dan `skemauji.Pasang`
  gagal untuk semua modul. `IDPEGA` treaty bisa 53 karakter > 50 kolom dasar — pelebaran oleh pemilik nbfacin. Bila
  182 sudah punya FK pada `ID` → ORA-02275; bila baris FAC yang ada tidak di `T_WORK_POLIS` → ORA-02298. `PRODKE DEFAULT
  0` mengisi baris FAC. 320 tidak aman diulang bila gagal di tengah (ORA-01430, perlu DBA).
- **K13:** `968_menu_nbtreatyin` sudah tercatat di `POOLDATA` oleh pihak lain.
- **K11:** uji db belum pernah berjalan (A4 `Makefile test-db`, C9 skema uji).
- **K12:** pemetaan peran IAM kosong — tempat `ProductionDate` dan 12 tempat tiket 05 tertutup.
- **F7:** pemuatan dokumen lama: migrasi → uji-kering di skema uji → produksi oleh WO/DBA (`MODUL.md`).
- Commit `88c303de` (`inti:`) ikut di cabang ini dan identik dengan PR inti — merge PR inti dulu.
