# PROMPT — LANJUTAN 3 MODUL **Treaty Contract Out** *(folder `OUTPUT_HASIL_RNM`, cabang `main` @ `5cc123f` atau lebih baru)*: **NOL TABEL BARU — modul menulis dan membaca tabel yang SUDAH ADA, persis seperti XML**

> Keputusan work owner 29-09-2026, kutipan: *"ini dihapus · untuk modul ini pakai table yang sudah ada · baca XML lagi aja ·
> khusus modul treaty contract out tidak ada tabel baru sama sekali!!!"*. Brief modul dan lanjutan 1–2 tetap berlaku **kecuali**
> yang dibantah di sini. Hanya konteks Treaty Contract Out. `App.tsx` berisi suntingan work owner yang belum di-commit — jangan
> di-commit atau dibuang. Worktree `.worktrees/desain-template` milik sesi lain — jangan disentuh. Setiap pembacaan activity dan
> RDB mencetak `pyStepsBlockName`.

## 0. YANG SUDAH DILAKUKAN ASISTEN DI DEV *(29-09-2026)*

| Hal | Keadaan |
| --- | --- |
| Delapan tabel migrasi `300`–`307` + delapan sequence | **dihapus** dari `POOLDATA` DEV *(`T_TREATYYEAR`, `T_TREATYCONTRACT`, `T_TREATYREINSURER`, `T_MTREATYSECURITY`, `T_TREATYBUSINESS`, `T_PROPORTIONALARRG`, `T_TREATYCO_JEJAK`, `T_TREATYYEAR_LAMPIRAN` beserta `SEQ_T_*`)*; sisa objek = 0 |
| `T_MIGRASI` | baris `300`–`307` dihapus; kini **29** langkah, terakhir `057` |
| Data | satu baris uji di `T_TREATYYEAR` + satu jejak sempat ada; dicadangkan asisten di luar repo sebelum dihapus |
| Tabel warisan | **utuh**: `TREATYYEAR`, `TREATYCONTRACT`, `TREATYREINSURER`, `MTREATYSECURITY`, `TREATYBUSINESS`, `PROPORTIONALARRG` |

⛔ Sampai paket 1 di bawah ter-commit, **`-migrate` tidak boleh dijalankan** — berkas `300`–`307` masih ada di repo dan akan
membuat ulang tabel itu.

## 1. KEPUTUSAN **tco4** — menggantikan **tco1** `[DIPUTUSKAN work owner 29-09-2026]`

| Unsur | Isi |
| --- | --- |
| Tabel | **nol** `CREATE TABLE`/`CREATE SEQUENCE` untuk modul ini. Seluruh tulis dan baca ke tabel **yang sudah ada** di `POOLDATA`, dengan **nama tabel dan kolom VERBATIM** |
| Sumber kebenaran | RDB `D:\XML\RNM_BRD\Treaty Contract Out\RDBList\`: penulis `SaveMasterTreatyYear_SQL`, `SaveMasterTreatyContract_SQL`, `SaveMasterTreatyReinsurer_SQL`, `SaveMasterTreatyBusiness_SQL` *(→ `PEGA_TREATYBUSINESS`)*, `SaveMasterProportionalArrg` *(→ `PEGA_PROPORTIONALARRG`)*, `SaveMasterProportionalArrgChild` *(→ `PEGA_M_PROPORTIONALARRG_CHILD`)*, `InsertToMTreatySecurity`, `UpdateMTreatySecurity`; penghapus `DeleteFromTREATYCONTRACT_SQL` *(→ `TREATYCONTRACT`, `TREATYREINSURER`, `MTREATYSECURITY`, `TREATYBUSINESS`)*, `DeleteFromTreatyReinsurer_Act`, `DeleteRowBusinessList`, `DeleteSecurityReinsurer`; pembaca `GetMaster*` *(antara lain `M_PROPORTIONALARRG`, `M_TREATYYEAR`, `TREATYEXCHANGE`)*; lampiran `InsertAtatchment_Sql`, `GetAttachment2_Sql`/`GetAllAttachment2_Sql` *(`M_ATTACHMENTTREATY_…`)*, `DeleteAttachment2_Sql`, `GetLinkStorage_SQL`/`Update_T_Storage_SQL`/`DeleteStorage_SQL` *(`T_STORAGE_IMAGE`)*, `GetTokenStorage_SQL`, `CategoryAttach_SQL` |
| Prosedur | keputusan **o** tetap: prosedur **tidak dipanggil**; isinya ditiru di Go dari `.scratch/treaty-contract-out/dba-procedures.md` *(`ALL_SOURCE` yang tercatat)* — tabel tujuan, kolom, **format nilai** *(tabel warisan ber-`VARCHAR2` untuk tanggal dan angka: tulis dengan bentuk teks yang sama persis dengan prosedur)*, pengenal *(sequence warisan yang dipakai prosedur, mis. `TREATYCONTRACT_SEQ`)*; nol `COMMIT` di teks SQL *(ADR-U-0029)*; bila `ALL_SOURCE` untuk satu prosedur tidak tercatat → `[terbuka — DBA]` berbukti, bukan dikarang |
| Tabel JSON `M_*` | dibaca/ditulis **hanya bila** RDB yang hidup *(bukan ter-remark)* membaca/menulisnya; bila spec menyatakannya "mati", buktikan dari rule — ralat bertanggal bila XML membantah spec |
| Jejak audit | **tidak** ada tabel jejak modul *(`T_TREATYCO_JEJAK` dibuang)*; Pega tidak mencatat jejak modul ini. Efek keluar lampiran tetap memakai outbox bersama `T_LOG_SERVICE_RNM` yang **sudah ada** *(milik aplikasi, bukan tabel baru modul ini)* |
| Lampiran | tabel lampiran warisan `M_ATTACHMENTTREATY_…` + `T_STORAGE_IMAGE` sesuai RDB; `T_TREATYYEAR_LAMPIRAN` dibuang |
| Kontrak hilir | lebih sederhana: Claim Prop, Komite Claim Prop, Claim Fac In **sudah** membaca tabel warisan ini; repository kontrak hilir membaca tabel yang sama |
| Penyimpangan sadar yang tetap | keputusan work owner lanjutan 2 *(OQ-10 tanggal +1 tahun kalender, OQ-17 `%Share` 0..100, OQ-18 `KURS`, OQ-21 hapus seperti Pega, OQ-20 induk)* tetap berlaku, diterapkan pada tabel warisan |

## 2. URUTAN — satu commit per paket

| # | Paket | Isi |
| ---: | --- | --- |
| 1 | **Buang migrasi 300–307** | hapus 16 berkas `300_*`–`307_*` *(+ `_down`)*; penghitung penjaga `migrasi_test.go` kembali ke **51**; penjaga STRUKTUR/kolom/kaskade yang merujuk tabel `T_` modul ini disesuaikan; `STRUKTUR-TABEL-TREATY-CONTRACT-OUT.md` ditulis ulang menjadi **peta tabel warisan yang dipakai** *(tabel → kolom VERBATIM → tipe katalog → RDB penulis/pembaca)*; uji hijau. Commit `treaty-contract-out: tco4 — nol tabel baru, migrasi 300-307 dibuang` |
| 2 | **Repository ke tabel warisan** | seluruh `tco_*` repository menulis/membaca tabel warisan sesuai §1; model mengikuti tipe warisan di tepi repository *(konversi teks ⇄ `apd.Decimal`/tanggal di satu tempat, diuji dua arah dengan contoh bentuk teks prosedur)*; uji murni + uji `db` *(SKIP tanpa skema uji)* |
| 3 | **Lampiran ke tabel warisan** | `M_ATTACHMENTTREATY_…` + `T_STORAGE_IMAGE`; pekerja latar dan pelaksana penyimpanan lanjutan 2 tetap |
| 4 | **Tiket dan dokumen** | tiket 01 *(PREFACTOR)* diralat bertanggal: skema relasional baru **dibatalkan** oleh tco4; tiket 03–12 diralat bila menyebut tabel `T_`; `spec.md`, PARITAS, LAPORAN-GILIRAN; OQ-TCO baru bila ada |
| 5 | **Uji penuh + `/code-review` singkat** | termasuk penjaga statik baru: **nol** `CREATE TABLE`/`CREATE SEQUENCE` di rentang `300`–`319`, dan nol nama `T_TREATY`/`T_MTREATY`/`T_PROPORTIONAL` di kode |

## 3. LAPORAN

Satu pesan: tabel **paket → commit → RDB XML → tabel warisan → rute/komponen**; ralat tiket; OQ; angka uji tiap commit dengan dan tanpa
tag `db`; bab **TELEMETRI EKSEKUSI**. Sesudah paket 1 ter-commit, beri tahu work owner bahwa `-migrate` aman dijalankan lagi *(ia
tidak akan membuat tabel apa pun untuk modul ini)*.

---

*Disusun 29 September 2026 sesudah penghapusan delapan tabel dan sequence migrasi 300–307 dari DEV (sisa objek 0, `T_MIGRASI` 29
langkah, enam tabel warisan utuh) dan sensus 36 RDB `Treaty Contract Out\RDBList\` (tabel dan prosedur yang disentuh).*
