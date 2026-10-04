# nbtreatyin: modul NB Treaty In (realisasi treaty masuk, Prop + NonProp)

> Urutan merge: **(1)** PR inti `inti/nbtreatyin-uji-menu` (`PR-INTI-UJI-MENU.md`) → **(2)** PR ini. Migrasi
> dijalankan work owner, bukan oleh merge — dan 320–327 baru sesudah C10 (tabrakan nama `T_GENERAL_POLIS`) beres.

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
│   └── migrations/    # 320–327 CREATE TABLE (8 tabel diagram) · 968 menu
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

Tabel: **tepat 8 `CREATE TABLE`** menurut diagram grilling sheet *NB Treaty In Prop/NonProp* — `T_GENERAL_POLIS`
(320, tabel Treaty In sendiri, shared PK `T_WORK_POLIS`), `T_POLIS_QUOTATION`, `_CEDING`, `_INSTALMENT`,
`_INSTALMENT_DETAIL`, `_SPREADING`, `_XOL`, `_XOL_LAYER` (321–327). Nol tabel lain. Uang dan persen `NUMBER(38,10)`
(diagram F20: skala minimal 9). Daftar kolom lawan diagram: `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`; kepatuhan sheet
*NB Treaty In Prop* butir demi butir: `docs/KEPATUHAN-SHEET-NB-TREATY-IN-PROP.md`.

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
TestPilihBisnisDiLuarDaftarPopupDitolak           → 422
TestBentukProporsionalMenolakXOLDanRincian        → Prop tanpa XOL dan tanpa rincian angsuran (diagram J56, J71)
TestTabelDanKolomMengikutiDiagramGrilling         → tepat 8 CREATE TABLE, kolom = diagram
TestSkalaUangPersenMinimalSembilanDiSemuaTabel    → NUMBER(38,10) seragam (diagram F20)
TestPraTerbangMenolakTGeneralPolisBentukFacIn     → -migrate berhenti sebelum apa pun bila T_GENERAL_POLIS FacIn ada
```

## Merge Danger

**Door:** two-way untuk kode (modul baru, rute baru, menu 968 sudah tercatat di `POOLDATA`); **one-way** untuk
migrasi 320–327 bila dijalankan (delapan tabel baru; mundur = `DROP`).

**Blast Radius:** basis-data-bersama.

Penahan sebelum migrasi dijalankan WO:

- **K18 / C10 — tabrakan nama:** `POOLDATA.T_GENERAL_POLIS` sudah ada, dibuat migrasi `nbfacin` `182_t_general_polis`
  (7 kolom FacIn, di luar repo). Keputusan WO 04-10-2026: `T_GENERAL_POLIS` **tetap tabel Treaty In sendiri** (diagram
  F9); **pemilik `nbfacin` perlu mengganti nama tabel FacIn**. Sampai itu, `-migrate` berhenti di pra-terbang inti
  (`praTerbangBentuk`) sebelum satu pernyataan pun dikirim — langkah semua modul yang belum tercatat ikut tertahan.
- **A5:** satu entri `NUMBER(38,10)` di penjaga inti `presisiSah` (commit `inti:` tersendiri) — mohon ditinjau tim inti.
- **K13:** `968_menu_nbtreatyin` sudah tercatat di `POOLDATA` oleh pihak lain.
- **K11:** uji db belum pernah berjalan (A4 `Makefile test-db`, C9 skema uji).
- **K12:** pemetaan peran IAM kosong — tempat `ProductionDate` dan 12 tempat tiket 05 tertutup.
- **F7:** pemuatan dokumen lama: migrasi → uji-kering di skema uji → produksi oleh WO/DBA (`MODUL.md`).
- Commit `88c303de` (`inti:`) ikut di cabang ini dan identik dengan PR inti — merge PR inti dulu.
