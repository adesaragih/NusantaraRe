# Permintaan untuk tim inti — modul `nbtreatyin`

> Wilayah sunting implementasi ini hanya `modul/nbtreatyin/` dan
> `inti/backend/daftar/modul_nbtreatyin_gen.go`. Setiap perubahan di luar itu yang dibutuhkan modul
> ini ditulis di sini, beserta sebabnya. ⛔ Tidak satu pun dikerjakan sendiri — kecuali bagian A
> (keputusan work owner K10, satu commit `inti:` terpisah).
>
> Cabang: `modul/nbtreatyin/implementasi` (dari `dev` 51d9b9cc). Dibuat 2026-10-03; ⭐ **daftar terkini per
> 04-10-2026** (konsolidasi putaran 2, paket P10). Satu tempat untuk SEMUA permintaan keluar modul ini:
> tim inti (A, B, E), pihak lain (C), berkas di luar wilayah (D), dan keputusan work owner (F).
> Rujukan butir dari `docs/HASIL-IMPLEMENTASI.md` dan tiket memakai nomor di sini (mis. "C5", "F3").

| Bagian | Isi | Butir |
| --- | --- | --- |
| A | sudah dikerjakan — **mohon ditinjau** | A1, A2 |
| B | konstanta dan kepemilikan bersama di `inti` | B1–B4 |
| C | sambungan dan data yang menunggu pihak lain (DBA, IAM, operator lingkungan, tim inti) | C1–C8 |
| D | berkas di luar wilayah yang termodifikasi pihak lain | D |
| E | kontrak `PembacaMasterTreaty` dari modul `treatyin` | E1–E4 |
| F | keputusan yang menunggu work owner | F1–F8 |

## A · Uji di luar wilayah yang berubah merah karena modul ini menyala — SUDAH DIKERJAKAN, mohon ditinjau

> ⭐ **Keduanya SUDAH DIKERJAKAN di commit `88c303de`** (`inti: sesuaikan dua uji lintas modul sesudah
> nbtreatyin menyala`) — **mohon ditinjau**. Dasar: keputusan work owner **K10** (PROMPT putaran 2
> bab 2): perbaikan minimal dalam **satu commit terpisah** berawalan `inti:`. Kolom *Permintaan* di
> bawah kini berarti **yang mohon ditinjau**, bukan yang mohon dikerjakan.

| # | Berkas | Yang terjadi | Permintaan |
| ---: | --- | --- | --- |
| A1 | `frontend/daftar.menuTabel.test.ts` baris 90-92 | Uji gigit "modul frontend berbaris `DIMIGRASI='0'`" memakai **nbtreatyin** sebagai contoh modul yang belum dimigrasi. Modul ini sekarang menyalakan barisnya (slot menu `968`), sehingga hasilnya bukan lagi selisih yang diharapkan. | ✅ **Sudah dikerjakan di commit `88c303de`, mohon ditinjau**: contohnya diganti `edmtreatyin` (berbaris `DIMIGRASI='0'`, tanpa slot menu). *(semula: "Ganti contohnya dengan modul yang barisnya masih `DIMIGRASI='0'`.")* |
| A2 | `modul/claimlife/backend/repository/batasanpemakaian_test.go:250` (`TestSetiapPemanggilBukaMemeriksaBolehDilewati`) | Sensus pemanggil `skemauji.Buka()` di seluruh repo: **17**, angka tertanam **16**. Pemanggil baru: `modul/nbtreatyin/backend/repository/polis_db_test.go` (uji repository lawan Oracle yang diminta brief; ia MEMERIKSA `BolehDilewati`). | ✅ **Sudah dikerjakan di commit `88c303de`, mohon ditinjau**: `const mau = 17` beserta komentar pemanggil barunya. *(semula: "Naikkan angkanya menjadi 17 (pemilik claimlife / tim inti).")* |

## B · Konstanta dan kepemilikan bersama

| # | Hal | Sebab | Permintaan |
| ---: | --- | --- | --- |
| B1 | `inti/backend/lini.go` hanya mengenal `LiniLife = "LIFE"` | Kasus NB Treaty In ditulis ke `T_WORK_POLIS.LINI`. Modul memakai konstanta lokal `models.LiniKasus = "NONLIFE"`, sejalan `KODE_PRODUKSI.TYPE = 'NONLIFE'` yang dibaca `GeneratePolicyNoTreaty_Act` (RDB `GetKodeProdNonLife_SQL`). | Tetapkan konstanta lini Non-Life bersama di `inti`; modul akan beralih memakainya. |
| B2 | `T_WORK_POLIS` + `SEQ_WORK_POLIS` milik **premiumlistlife** (migrasi 050/057/059) | Tabel kasus lintas-lini; migrasi 320 `nbtreatyin` membuat kunci tamu ke sana. Modul ini karenanya **membutuhkan premiumlistlife aktif** (MODUL_AKTIF) dan migrasinya berjalan lebih dulu (rentang 050 < 320). | Pertimbangkan memindah `T_WORK_POLIS` ke inti, atau nyatakan ketergantungan antar-modul di perakit. |
| B3 | `inti.Pelaku` tidak membawa nama tampilan | AC 40/42 menuntut `OperatorName`/`USERNAME` = **nama tampilan** (`OperatorID.pyUserName`). Modul membaca `M_LOGIN_GO.NAME` sendiri (pola modul `marketingofficer`). | Bila diinginkan satu sumber, sediakan nama tampilan di `inti` (mis. `Pelaku.Nama` dari sesi). |
| B4 | ⭐ *(baru, P10)* `inti.Pelaku` tidak membawa divisi operator | `HISTORYAKSEPTASIPRODUCTION.DIV` = `OperatorID.pyOrgDivision` (`SaveViewSuggest` langkah 2 `CARI6`, `InsertViewSuggest_SQL`). `inti.Pelaku` hanya `AkunID` + `Peran` ⇒ modul menulis **NULL** (`backend/models/usulan.go`; tiket 10, 19). | Sediakan divisi operator di `inti` (sumbernya — mis. kolom `M_LOGIN_GO` — ditetapkan tim inti / IAM), atau nyatakan `DIV` tetap NULL. |

## C · Sambungan dan data yang menunggu pihak lain

| # | Hal | Pihak | Catatan |
| ---: | --- | --- | --- |
| C1 | Sambungan konversi Arasapas treaty (`serviceInsertArasapas_act`) | tim inti + `[pemilik export Pega]` | Muatan 4 medan dibangun (KEPUTUSAN-RONDE-12 butir 7). Kunci `M_LINK_SERVICE` rule treaty tidak terekspor; pemanggilan layanan luar menuntut persetujuan (pola `outbox.ErrArasapasBelumDisetujui`). Pengirim bawaan modul GAGAL TERANG di produksi, dilewati di luar produksi. Antre-ulang (`outbox.Antrean`) dan email kegagalan belum dipasang — keduanya `BelumDiputuskan` di inti. |
| C2 | Pemetaan tempat → peran → arah — ⛔ ~~`M_NBTRIN_PERAN_TEMPAT` (migrasi 330)~~ **konstanta `backend/models/peran_tempat.go` `PemetaanPeranTempat`** (K16, 03-10-2026: tabel di luar diagram grilling, dihapus) | `[IAM]` + work owner — **K12** | Kosong sampai IAM menjawab (K12; tiket 05 `needs-info`). ⭐ **P9 (04-10-2026): kedua belas tempat tiket 05 terdaftar** di `models.DaftarTempat` (DetailDeptHeadTreatyIn_UW 3 tombol Submit, GeneralDeptHeadTreatyIn_UW 3, ListSuggest `.ProductionDate` tampil ×2 + wajib ×2, DetailPoliciesNonProportional LABEL "NON EDM"/"EDM"); yang dibaca layanan/layar: empat tempat ProductionDate dan dua label (enam tombol Submit diganti posisi kasus). Isian per tempat berikut ARAH `MUNCUL`/`KECUALI` — tiket 05; peran pengguna dari `inti.Pelaku.Peran`. *(Bunyi P1 dikutip: "Isian per tempat (saat ini satu tempat terpakai: `LISTSUGGEST_PRODUCTIONDATE`)".)* *(Bunyi lama dikutip: "`M_NBTRIN_PERAN_TEMPAT` (migrasi 330) … Nol baris dari migrasi.")* |
| C3 | Tipe kolom view `TREATYINDETAILJOINEDM`, tabel `TREATYINDETAIL`, kolom `AGENT.STS_PKP` / `STATUSACTIVE` | ~~`[DBA]`~~ ⭐ **katalog Oracle** (view dan `AGENT`) · `[DBA]` (tabel `TREATYINDETAIL` saja) | ⭐ **Dijawab dari katalog `ALL_TAB_COLUMNS`, dicek 03-10-2026 (PROMPT putaran 2 bab 1) — tidak perlu ditanyakan ke DBA lagi:** (1) **`POOLDATA.TREATYINDETAILJOINEDM`** — kolom nilai `LIMITVALUE`, `RETENTIONVALUE`, `EPIVALUE`, `NETPREMIVALUE`, `SHAREVALUE`, `MDPVALUE`, `DEDUCTION1`, `DEDUCTION2`, `RIOGR`, `RIONR`, `RNM_SHARE`, `BROKERAGE` = `NUMBER` **tanpa presisi dan skala**; `COMMENCEMENT`, `TERMINATION` = `DATE`; `INSTALLMENTNO` = `VARCHAR2(1000)` (teks di view, dikonversi ke bilangan bulat di modul); kolom lain `VARCHAR2(1000)`, `ID` `VARCHAR2(100)`. (2) **`POOLDATA.AGENT`** — `STS_PKP` dan `STATUSACTIVE` **ada**, keduanya `VARCHAR2(1000)`; nama kolom dari properti RD terbukti benar. Repository tetap membaca tipe view dari `SYS.ALL_TAB_COLUMNS` saat berjalan (`repository/acuan.go`). ⛔ Yang **belum** tercakup fakta katalog itu: tipe kolom tabel `TREATYINDETAIL` (grid popup) — tetap `[DBA]`. *(Bunyi lama dikutip: "Repository membaca tipe dari `SYS.ALL_TAB_COLUMNS` saat berjalan; nama kolom AGENT diambil dari nama properti RD (belum dikonfirmasi DBA).")* |
| C4 | `make test-db` — skema uji Oracle **K11** | operator lingkungan + `[DBA]` (nama skema uji, **bukan** `POOLDATA`) | Uji `-tags db` modul ini MELEWATI tanpa `ORACLE_SCHEMA` + `ORACLE_SKEMA_UJI=true`; **belum pernah dijalankan** lawan skema uji (K11 kosong). Yang menunggu: `repository/polis_db_test.go`, `repository/lama_db_test.go`, `repository/masterxol_db_test.go` (spec AC 23, 29, 31, 68; spec-penyimpanan AC 1, 6, 13, 46, 49–51, 55). Skema uji perlu memuat tabel warisan yang dibaca/ditulis (`HISTORYAKSEPTASIPRODUCTION`, `HISTORYAKSEPTASIPEGA`, view `TREATYINDETAILJOINEDM`, `M_TREATY_IN`) — uji melewati tabel warisan yang tidak ada. |
| C5 | ⭐ *(baru, P10; dari tiket 22)* `SEQ_WORK_POLIS` sesudah pemuat dokumen lama | tim inti / pemilik premiumlistlife + operator yang menjalankan pemuat | ID kasus lama = pyID `NB-<n>`; sequence **wajib dimajukan** melewati nomor `NB-` terbesar yang dicetak ringkasan pemuat sebelum kasus baru dibuat (pola OQ-PL-15). Sequence milik premiumlistlife — modul ini tidak mengubahnya. |
| C6 | ⭐ *(baru, P10; dari P8)* Klausa wadah portal `OperatorID.pyWorkGroup!='ReasLife'` (= `!When/IsOperatorLife`) | `[IAM]` — **K12** | `Section/SFAPortal_OpportunitiesList` menampilkan grid `GetListOpportunity` hanya bila `pyWorkGroup!='ReasLife' && pyWorkBasketList(2)=='ReasTreatyInAdmin'`. Klausa kedua dibangun menurut nama antrean (AC 14); klausa pertama **tidak**: work group Pega tak berpadanan di `inti.Pelaku` maupun `M_LOGIN_GO`. Mohon padanan work group (atau pernyataan bahwa klausa itu gugur) — tiket 04, 05; `docs/issues/04-*.md`. |
| C7 | ⭐ *(baru, P10; dari tiket 22)* Empat angka dokumen lama (KEPUTUSAN-RONDE-12 butir 5) | `[DBA]` | cacah baris `JSON_POLIS` generasi NB, tahun terawal, ukuran, `PRODKE` tertinggi — dibutuhkan untuk merencanakan jalankan pemuat; pemuat sendiri mencetak cacahnya saat uji-kering. |
| C8 | ⭐ *(baru, P10; dari P1)* Tipe kolom fisik `POOLDATA.HISTORYAKSEPTASIPRODUCTION` (terutama `NOURUT`, `PERCENT_RNM`, `TGL_INP`) | `[DBA]` / katalog `ALL_TAB_COLUMNS` | belum dicek ke katalog; pembacaan memakai `TO_NUMBER(NOURUT)` (`repository/usulan.go`). Bila `NOURUT` bertipe angka, `TO_NUMBER` tetap aman; bila teks berisi nilai bukan angka, pembacaan gagal — mohon tipe fisiknya. `docs/PERBANDINGAN-KOLOM-DIAGRAM.md` bab 11 butir 2. |

## D · Berkas di luar wilayah yang termodifikasi oleh pihak lain

`package.json` dan `package-lock.json` sudah termodifikasi di salinan kerja sebelum implementasi ini
dimulai (vite `^8.3.1` → `^7.3.6`, `@vitejs/plugin-react` `^4.7.0`). ⛔ Tidak di-commit dan tidak
dikembalikan oleh implementasi ini.

## E · Kontrak `PembacaMasterTreaty` dari modul `treatyin` (K8, 03-10-2026)

`[keputusan work owner]` K8: jalur NB NonProporsional / XOL membaca master kontrak treaty BACA-SAJA dari
`JSONDATA` `M_TREATY_IN` / `M_TREATY_IN_EDM` (`RDBList\BrowseTreatyIn`, `BrowseTreatyInJoinEDM`) — sebuah
`[penyimpangan sadar]` atas P29, sebab tabel master relasional modul `treatyin` (`KONTRAK`, `LAYER`,
`BAGIAN`, `PEMULIHAN_LIMIT`, `TERMIN`, `POTONGAN`) masih **nol baris** (dicek 03-10-2026) dan medannya tidak
ada di view `TREATYINDETAILJOINEDM`. Pembacaan itu terisolasi di SATU fungsi
(`modul/nbtreatyin/backend/repository/masterxol.go`, `MasterXOLDariJSON`) di balik antarmuka
`services.PembacaMasterTreaty`.

| # | Permintaan | Pihak | Rincian |
| ---: | --- | --- | --- |
| E1 | Sediakan kontrak `PembacaMasterTreaty` dari modul `treatyin` (baris `Kontrak disediakan` di `MODUL.md`-nya) | pemilik `treatyin` + tim inti (perakit kontrak) | Satu metode: `MasterXOL(ctx, noKontrak string) (MasterXOL, error)` — master satu kontrak menurut nomor kontrak (`TREATYID` = `PolicyTreatyIn.NoOffer`); nol master = galat "tidak ada", bukan master kosong. Bentuk dan daftar medannya: `modul/nbtreatyin/backend/models/masterxol.go` (`SkalarMasterXOL`, `DaftarMasterXOL` — Share() beserta GrossPremiumList/NetPremiumList/DeductionList/DeductionTotalList/RnmLimitList/SpreadingListXOL, Installment().InstallmentList, FacultativeShare(List), RNMShare, EDMState, ProportionType, Limits() pemulihan, ringkasan dan total layer). Angka sebagai teks desimal (nol float), tanggal `YYYY-MM-DD`. |
| E2 | Pemetaan medan master lawan tabel relasional `treatyin` | pemilik `treatyin` | Tiap medan di E1 dipetakan ke kolom `KONTRAK`/`LAYER`/`BAGIAN`/`PEMULIHAN_LIMIT`/`TERMIN`/`POTONGAN` — atau dinyatakan tidak ada (mis. total dan ringkasan layer yang di Pega disimpan di dokumen). |
| E3 | Begitu E1 tersedia | `nbtreatyin` | `penyimpanOracle.MasterXOL` (`services/gudang.go`) beralih ke kontrak; `repository.MasterXOLDariJSON` dan uji urainya DIHAPUS; `MODUL.md` baris `Kontrak dipakai` diisi. Nol perubahan di `models` dan `handlers`. |
| E4 | ⭐ *(baru, P10; dari P5)* Satu contoh `JSONDATA` `M_TREATY_IN` (dan `M_TREATY_IN_EDM`) **disamarkan**, untuk mencocokkan bentuk yang diurai | work owner / `[DBA]` | Pengurai (`repository/masterxol.go`, `uraiMasterXOL`) mengikuti nama properti XML (`adoptJSONObject` atas halaman `TreatyIn`) dan diuji dengan dokumen fiktif `UJI-` (`repository/masterxol_test.go`); bentuknya **belum dicocokkan dengan contoh nyata**. Cap waktu ` GMT` di master dibaca jam dinding Asia/Jakarta (`models.TanggalMasterPega`). Satu contoh per tabel cukup — nama orang dan nomor polis disamarkan sebelum diserahkan. |

⛔ Sampai E1 tersedia, `M_TREATY_IN` / `M_TREATY_IN_EDM` tetap dibaca (baca-saja) dan **tidak pernah
ditulis** dari modul ini. Treaty keluar (`M_TREATY_OUT`, `RDBList\BrowseTreatyOut`) tetap tidak dibaca
(K8 butir 4).

## F · Keputusan yang menunggu work owner

Butir di bawah bukan pekerjaan tim inti, tetapi dicatat di sini supaya satu tempat memuat seluruh permintaan
keluar modul ini. F1–F5 **sudah dibangun dengan tafsiran yang ditulis** dan menunggu konfirmasi; F6–F8 menunggu
keputusan sebelum ada yang dibangun / dijalankan.

| # | Hal | Tafsiran yang dibangun / keadaan | Yang dibutuhkan |
| ---: | --- | --- | --- |
| F1 | Ukuran pengecualian **K8** (baca-saja JSON master) | K8 = medan master yang **dibaca rule terjangkau jalur NB NonProp** — setiap medan di `backend/models/masterxol.go` berkutip langkah XML (lebih luas dari daftar harfiah K8: `RNMShare`, `RnmShareDeducted`, `EDMState`, `ProportionType`, `Limits()`, ringkasan dan total layer); tabel tetap hanya `M_TREATY_IN`/`M_TREATY_IN_EDM`; alasan status.json dinilai ulang dengan ukuran ini (tiket 01 bab P9) | `[menunggu konfirmasi WO]`: setujui ukuran ini, atau tetapkan daftar harfiah (akibatnya jalur NonProp kehilangan medan yang dibaca XML) |
| F2 | Tiga penyimpangan **K4** catatan usulan → `HISTORYAKSEPTASIPRODUCTION` | (1) ditulis di submit KETIGA jenjang (XML: hanya pasca-submit admin) — akibat langsung K4; (2) `TGL_INP` jam 24 (XML `hh`); (3) `NOURUT` dari repository (XML `.pxListSubscript`; terbukti bernilai sama, uji di tiket 10) | `[penyimpangan sadar — menunggu konfirmasi WO]` untuk ketiganya (tiket 10 bab P9) |
| F3 | **AC 59 / K17**: CSV medan tak dikenal pemuat dokumen lama tidak dapat 0 secara struktur | pemuat menulis medan tanpa kolom ke CSV, kode keluar ≠ 0 | keputusan per medan dokumen lama tanpa kolom (22 pola menurut data guide + `IsOJKNopolis`, `InstallmentList().PPN/PPh`) dan `SuggestList` dokumen lama (salin ke `HISTORYAKSEPTASIPRODUCTION` atau dibuang) — daftar lengkap di tiket 22 bab P9 |
| F4 | ⭐ *(baru, P10; dari P3)* Pemilih **Source Of Business** (ClaimType `XOL Retro`) menyimpan saat klik | Pega memegang hasil `SearchHierarkiSourceBizAgent_PostDT` di clipboard sampai Save/Submit; di sini `POST /kasus/{id}/pilih-sumber-bisnis` menyimpan `Quotation.SourceOfBusiness` dalam satu transaksi, sebab layar tidak boleh mengirim `Quotation.*` (`services/sumberbisnis.go:41`, `[penyesuaian sadar]`) | `[menunggu konfirmasi WO]` |
| F5 | ⭐ *(baru, P10; dari P5)* Keanehan Pega jalur XOL **ditiru apa adanya** | `InsertToTreatyXOLList`: `DueTo` selalu kosong, `Currency` induk dari layer terakhir, potongan kumulatif antar layer (`models/nonprop.go:133`); `InsertToTreatyXOLListRetroShare`: nilai CARI terbawa antar layer, langkah 2.3.7.4 menambah net/potongan layer (`:220`); `TreatyNonPropSetSpreading`: `PremiumSpreaded = local.netpremi` (0) lalu dihitung ulang langkah 10 (`:340`); `InputDetailNonProp` 16.2.2/17/19 (`models/nonprop_detail.go:111`); preACT 18 (`:242`) | `[menunggu konfirmasi WO]`: tetap ditiru (patokan XML), atau diperbaiki dengan keputusan tertulis per butir |
| F6 | ⭐ *(baru, P10; dari tiket 22)* Penanda **`SUMBER='PEGA'`** baris hasil pemuat (KEPUTUSAN-RONDE-12 butir 5 poin 2) | ⛔ tidak ada kolom di diagram grilling — tidak dibuat; baris hasil pemuat dikenali dari `IDPEGA` berbentuk `<kelas> <pyID>` dan status `Resolved-Completed` | tambah kolom lewat RALAT diagram grilling, atau nyatakan penanda itu gugur |
| F7 | ⭐ *(baru, P10; dari tiket 22)* **Pemuatan dokumen lama di produksi** | `-jalankan` ditolak bila `IS_PEGA_PROD=true` (pola `pindahflat`, `backend/alat/pemuatlama/main.go:66`); uji-kering boleh | keputusan kapan dan oleh siapa pemindahan di produksi dijalankan (sesudah migrasi 320–327, C5, F3) |
| F8 | ⭐ *(baru, P10)* **`GENERATE_SEQUENCE_NUMBER`** — diagram grilling (Prop F103, NonProp F118) menyebutnya **"dibaca saja"** | penerbitan nomor polis memakai penomor bersama `inti/backend/penomor` (`UrutNomorBerikut`: `SELECT … FOR UPDATE`, lalu `INSERT`/`UPDATE` baris deret) — padanan procedure `PROC_GENERATE_SEQUENCE_NUMBER` yang di Pega menulis tabel yang sama; spec-penyimpanan AC 48 melarang memanggil procedure | konfirmasi bahwa "dibaca saja" di diagram berarti "tidak ditulis langsung oleh rule NB" (penulisnya procedure/penomor bersama), atau RALAT diagram |
