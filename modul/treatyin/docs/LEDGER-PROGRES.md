# Ledger progres — seluruh 77 tiket

Diperbarui 3 Oktober 2026. **Tiap kali sesuatu mendarat, berkas ini ikut berubah.**

Berkas ini ada karena enam belas tiket berskema-tanpa-aplikasi **hilang dari pandangan selama dua
ronde**: tidak ada satu pun tempat yang memperlihatkan tiket mana sampai lapisan mana, sehingga
laporan tiap ronde hanya memuat yang mendarat, dan yang ke-skip tidak terlihat.

| Lambang | Artinya |
| --- | --- |
| OK | mendarat |
| — | belum — **kolom catatan wajib menyebut sebabnya** |
| n/a | tidak berlaku bagi tiket ini |

**`uji Oracle` berarti dijalankan terhadap Oracle sungguhan**, bukan terhadap teks DDL. Sembilan
constraint terbukti menolak dan dua jalur positif terbukti diterima pada 2 Oktober 2026
([`VERIFIKASI-ORACLE-2026-10-02.md`](VERIFIKASI-ORACLE-2026-10-02.md)); ia **belum** otomatis —
`make test-db` menunggu skema uji, dan akun aplikasi tidak punya `CREATE USER`.

⛔ **Nol tiket berstatus `selesai`**, dan itu bukan kelalaian: status adalah medan di dalam berkas
tiket, milik pemilik proses. Ledger ini mencatat **lapisan**, bukan status.

## Treaty In — batch 1 (`14`–`44`)

| tiket | skema | constraint | kaskade | migrasi | uji bentuk | uji Oracle | services | handler | catatan |
|---|---|---|---|---|---|---|---|---|---|
| `14` | OK | OK | OK | 401 | OK | OK | OK | OK | satu-satunya yang utuh sampai handler. INV-53 + INV-29 ditegakkan di services |
| `15` | OK | OK | OK | 400 | OK | OK | OK | OK | jalur baca enam tabel acuan |
| `16` | n/a | n/a | n/a | n/a | n/a | — | OK | OK | peringatan kunci alami ganda **menyebut pembandingnya**; simpan tetap BERHASIL (ADR-0040 §2). Uji Oracle menunggu skema uji |
| `17` | n/a | n/a | n/a | n/a | n/a | — | OK | OK | cari lewat nomor warisan; hasil kosong = jawaban, baris non-warisan tidak ikut terambil |
| `18` | n/a | n/a | n/a | n/a | n/a | — | OK | OK | kelima ruas beku ditolak, pesannya menyebut SELURUH ruas yang menyimpang; `UPDATE` SQL-nya pun tidak menyentuhnya |
| `19` | n/a | n/a | n/a | n/a | n/a | — | OK | OK | versi menyimpang ditolak, pesannya menyebut ruas + nilai yang berlaku. INV-59: nol kolom beku di `VERSI_KONTRAK` |
| `20` | OK | OK | OK | 403 | OK | **OK** | — | — | INV-38 dan INV-57 belum ditegakkan — pernyataan keputusan di 403 **Terbukti di Oracle ronde 8** — TestTiket20KursMataUangGandaDitolak · TestTiket20DuaMataUangPadaSatuVersiDiterima. |
| `21` | n/a | n/a | n/a | n/a | n/a | — | OK | — | kegagalan hitung mengembalikan **nil + keterangan**, bukan nol dan bukan kurs satu (ADR-0035). **Handler belum** — belum ada rute yang menghitung |
| `22` | OK | OK | OK | 404 | OK | **OK** | — | — | INV-39/40 lubang, bukan penundaan — lihat 404 **Terbukti di Oracle ronde 8** — TestTiket22RetensiCedantGandaDitolak · TestTiket22RetensiDuaMataUangDiterima. |
| `23` | OK | OK | OK | 405 | OK | **OK** | — | — | **Terbukti di Oracle ronde 8** — TestTiket23EgnpiGandaDitolak · TestTiket23EgnpiDuaMataUangDanKelasBisnisKosongDiterima. |
| `24` | OK | OK | OK | 406 | OK | **OK** | — | — | **Terbukti di Oracle ronde 8** — TestTiket24PortofolioGandaDitolak · TestTiket24PortofolioEmpatKombinasiDiterima. |
| `25` | OK | OK | OK | 407 | OK | **OK** | — | — | INV-55 menunggu jalur simpan **Terbukti di Oracle ronde 8** — TestTiket25PeriodePelaporanGandaDitolak · TestTiket25EmpatPeriodePelaporanDiterima. |
| `26` | OK | OK | OK | 408 | OK | **OK** | — | — | INV-56 menunggu jalur simpan **Terbukti di Oracle ronde 8** — TestTiket26PeriodeAkumulasiGandaDitolak · TestTiket26AkumulasiKeduaDanKolomOpsionalKosongDiterima. |
| `27` | OK | OK | OK | 409 | OK | **OK** | — | — | **Terbukti di Oracle ronde 8** — TestTiket27TerminGandaDitolak · TestTiket27TerminNomorSamaDuaMataUangDiterima. |
| `28` | OK | OK | OK | 410 | OK | **OK** | — | — | **Terbukti di Oracle ronde 8** — TestTiket28SkalaKoasuransiGandaDitolak · TestTiket28TigaBarisSkalaKoasuransiDiterima. |
| `29` | OK | OK | OK | 411 | OK | OK | — | — | uji positif "bahaya kesembilan" ada di `invarian_db_test.go` |
| `30` | OK | OK | OK | 412 | OK | **OK** | — | — | **Terbukti di Oracle ronde 8** — TestTiket30DokumenGandaDitolak · TestTiket30DokumenSamaPadaDuaVersiDiterima. |
| `31` | OK | OK | OK | 413 | OK | OK | — | — | INV-05 terbukti menegakkan — klaim sebaliknya dicabut |
| `32` | OK | OK | **sebagian** | 413 | OK | OK | OK | OK | **mendarat ronde 5.** Tarif berbeda per pemulihan diterima, di luar 0-100 ditolak menyebut baris + nilainya. ⚠️ `kaskade` sebagian: `UNIQUE (ID_LAYER, NOMOR_URUT_PEMULIHAN)` **TIDAK dipasang** — `SPEC-INVARIAN.md` berhenti di `INV-71` dan penjaga melarang `UNIQUE` tanpa nomor. Kembar hanya ditolak di services |
| `33` | OK | OK | OK | 416 | OK | OK | — | — | INV-33 tertahan `F-13` |
| `34` | OK | OK | OK | 417 | OK | **OK** | — | — | INV-32 tertahan `F-13` **Terbukti di Oracle ronde 8** — TestTiket34DetailProporsionalGandaDitolak · TestTiket34DuaKelompokDalamSatuLayerDiterima. |
| `35` | n/a | n/a | n/a | n/a | n/a | — | OK | — | XOR quota-share/surplus + kesesuaian `JENIS_TREATY`; lima jalur negatif, dua positif. **Handler belum** |
| `36` | n/a | n/a | n/a | n/a | n/a | — | OK | — | surplus tanpa quota share ditolak, pesannya menyebut apa yang kurang. **Handler belum** |
| `37` | OK | OK | OK | 418 | OK | OK | — | — | `CK_POTONGAN_INDUK` terbukti dua arah; INV-52 tertahan `F-13`, INV-63 tinjauan kode |
| `38` | — | — | — | — | — | — | — | — | tertahan `Uji AD`. `L-3` sudah tutup, jadi penahannya tinggal satu |
| `39` | OK | OK | OK | 414 | OK | **OK** | — | — | entitasnya saja; **mekanisme pengisian** belum — menunggu jalur simpan **Terbukti di Oracle ronde 8** — TestTiket39JejakYatimDitolakDanVersiBerjejakTidakDapatDihapus · TestTiket39DuaJejakRuasSamaPadaVersiSamaDiterima. |
| `40` | n/a | n/a | n/a | n/a | n/a | — | OK | OK | **mendarat ronde 5.** Versi dasar dibaca lewat `ID_VERSI_KONTRAK_DASAR`; uji positifnya MENGUBAH versi dasar di tengah jalan — rancangan yang menyalin gagal di situ. Versi pertama menjawab 200 + `"ada": false`, bukan 404 |
| `41` | n/a | n/a | n/a | n/a | n/a | — | OK | OK | **mendarat ronde 5.** Diturunkan saat dibaca — nol tabel, nol kolom, nol cache (INV-58). Kontrak sistem baru menjawab `nomorLamaAda: false`, **bukan** nomor karangan |
| `42` | OK | OK | OK | 425 | OK | **OK** | OK | OK | **mendarat ronde 6**, sesudah tabelnya berdiri — urutannya tidak dibalik. `INV-61` diwujudkan sebagai BENTUK: `models.BuktiArsip` tidak punya ruas muatan, dan `TestArsipTidakPunyaJalurBaca` menyapu **34 berkas Go, 0 pembaca** — angkanya dicetak, seperti tiket 42 tuntut **Terbukti di Oracle ronde 8** — TestTiket42ArsipYatimDitolakDanKontrakBerarsipTidakDapatDihapus · TestTiket42MuatanRusakDanPengirimanKeduaDiterima. |
| `43` | n/a | n/a | n/a | n/a | n/a | — | OK | — | sakelar pemindahan: keadaan terlihat tanpa basis data, alasan wajib, menyalakan **memeriksa ulang** baris yang masuk. **Handler belum** |
| `44` | — | — | — | — | — | — | — | — | isi tabel acuan dari sistem lama; **bertenggat** — ia memuat data pertama |

## Treaty In — batch 2 (`45`–`64`)

**Nol dari dua puluh punya tabelnya.** `CATATAN_PERSETUJUAN` (tiket `54`) belum berdiri, dan ia yang
`49`–`56` gantungi.

| tiket | skema | penahan | catatan |
|---|---|---|---|
| `45` | — | — | daftar keadaan siklus hidup. **Pembuka batch** — sendirian melepas `46 47 48 54 62`. Kolomnya berdiri sejak `401`; **artinya** belum. Tidak dikerjakan ronde 4 — di luar daftar ronde ini |
| `46` | — | `45` | versi terminal beku menjangkau seluruh anaknya |
| `47` | — | `45` | paling banyak satu versi tak-terminal per kontrak |
| `48` | — | `45` | pengajuan menolak dua syarat, memperingatkan enam |
| `49` | — | `48` · `54` | SH menyetujui, versi pindah ke antrian DH |
| `50` | — | `49` | DH menyetujui, versi pindah ke antrian DR |
| `51` | — | `50` | DR menyetujui, versi menjadi `DISETUJUI` |
| `52` | — | `51` | pengaju tidak menyetujui; tiap tingkat orang berbeda |
| `53` | — | `49` | pengembalian ke `DRAFT` dengan alasan wajib |
| `54` | — | `45` | **`CATATAN_PERSETUJUAN` — PEMBUAT PERTAMA**, tabelnya belum berdiri. `49`–`56` menggantungnya |
| `55` | — | `49` · `54` | penolakan versi meninggalkan catatan |
| `56` | — | `45` · `54` | pembatalan draf oleh pembuatnya sendiri |
| `57` | — | `46` · `20` | angka rupiah pada versi disetujui tidak bergeser |
| `58` | — | `51` · `75` · **dokumen batas wewenang** | **berstatus `tertahan`** — dan penahannya **DUA**: dokumennya belum ada, dan tabelnya pun belum (`BATAS_WEWENANG`, tiket `75`) |
| `59` | — | `45` · `44` | kepala kontrak warisan dipindahkan beserta keadaannya |
| `60` | — | `59` | keadaan warisan tanpa padanan → `WARISAN_TAK_TERPETAKAN` |
| `61` | — | `60` · `71` · `73` | baris warisan diperbaiki ke keadaan sah; penelusurannya menuntut `MIGRASI_KORELASI` (`71`) dan daftar kerjanya `MIGRASI_NILAI_DITOLAK` (`73`) — **keduanya belum berdiri** |
| `62` | — | `45` | kontrak dicari termasuk menurut keadaan siklus hidupnya |
| `63` | — | `14` · **`Uji X-2`** | **berstatus `tertahan`** — cara pembukuan XOL |
| `64` | — | `15` · **penetapan pemilik** | **berstatus `tertahan`** — tabel acuan pembagian kapasitas |

**Dua puluh tiket, nol berskema.** `CATATAN_PERSETUJUAN` (tiket `54`) adalah satu-satunya tabel baru
yang batch ini butuhkan, dan ia belum berdiri.

> ⚠️ **RALAT.** Bagian ini pernah meringkas tiga belas tiket menjadi dua baris rentang
> (*"`46`–`53`, `55`–`58`"*), sehingga ledger yang mengaku mencakup **64 tiket** sebenarnya hanya
> memuat **51**. Itu terjadi di berkas yang justru dibuat supaya yang ke-skip terlihat — dan
> peringkasan itu menyembunyikan persis apa yang hendak diperlihatkannya.

## Treaty In Adjustment (`01`–`13`)

| tiket | skema | constraint | kaskade | migrasi | uji bentuk | uji Oracle | services | handler | catatan |
|---|---|---|---|---|---|---|---|---|---|
| `01` | OK | OK | OK | 440 | OK | **OK** | — | — | **Terbukti di Oracle ronde 8** — `TestTiket01VersiDasarTakAdaDitolakDanVersiDasarTidakDapatDihapus` · `TestTiket01VersiPertamaTanpaDasarDiterima`. Keempat penolakan lain menunggu jalur simpan |
| `02` | sebagian | — | n/a | 401 | — | — | — | — | kolomnya ada sejak 401; **dua nilainya belum dinyatakan** — sisa SKEMA |
| `03` | sebagian | — | n/a | 401 | — | — | — | — | kolomnya ada; **bawaan tanggal mulai belum ada** — sisa SKEMA |
| `04` | — | — | — | — | — | — | — | — | `DOKUMEN_ADDENDUM`; tertahan `DB-16a`, belum dikirim |
| `05` | OK | OK | n/a | 441 | OK | OK | n/a | n/a | `NOMOR_URUT_VERSI` nullable — terbukti di Oracle |
| `06` | n/a | — | — | — | — | — | — | — | **tidak lagi BLOKIR**: `NILAI_SELISIH` berdiri lewat migrasi `426`. Tiket ini kini tinggal perilaku padanan kunci bisnisnya, dan **dapat dimulai** |
| `07` | — | — | — | — | — | — | — | — | titik beku materialitas; tertahan `DB-20` |
| `08` | — | — | — | — | — | — | — | — | tanggal berlaku dokumen; tertahan `DB-16b` |
| `09` | — | — | — | — | — | — | — | — | nomor dokumen warisan; menunggu `04` |
| `10` | n/a | — | n/a | n/a | — | — | — | — | pengisian nomor urut warisan; **tidak dikerjakan ronde 3** |
| `11` | — | — | — | — | — | — | — | — | INV-69/70 + pemantau kebasian; menunggu `06`, dan `06` menunggu `76` |
| `12` | n/a | — | n/a | n/a | — | — | — | — | kerutkan urutan; menunggu `10` |
| `13` | n/a | — | — | — | — | — | — | — | **tidak lagi BLOKIR** oleh `NILAI_SELISIH` — tabelnya berdiri (`426`). Masih menunggu `11` |
| `76` | **OK** | sebagian | OK | **426** | OK | — | — | — | **skema mendarat ronde 6.** `VERSI_KONTRAK` ikut hapus (`ERD.md` §2.6). ⚠️ `constraint` **sebagian**: kunci asing KEDUA yang §2.6 tuntut tidak dapat dipasang — `BESARAN_DAPAT_DISESUAIKAN` tidak ada di mana pun, jadi `KODE_BESARAN` berdiri sebagai teks. ⚠️ Berkasnya di modul `treatyin`: penjaga `TestNolTabelBaru` menolak tabel di modul Adjustment, dan penolakan itu benar |
| `77` | **OK** | OK | OK | **426** | OK | — | — | — | **skema mendarat ronde 6.** Induknya `VERSI_KONTRAK`, bukan `KONTRAK`: ERD menulis `TREATY_IN`, dan `TREATY_IN` dipecah dua oleh `ADR-0040`. Yang memutuskan kolom `MEMAKAI_PRORATA` — ia berdiri di `VERSI_KONTRAK` |

## Treaty In — batch 3 (`65`–`75`) · entitas dari rekonsiliasi 2 Oktober 2026

**Enam dari sebelas kini punya tabelnya** (`65` `67` `68` `69` `70` `74`, migrasi `422`–`425`).
Kesebelasnya ditiketkan 2 Oktober 2026 dari
`4-erd-dan-tabel-datar/REKONSILIASI-XML-VS-DDL.md` — 18 entitas yang sapuan 329 XML punya dan
`ddl-usulan/` tidak. Ledger mencatatnya **sejak hari pertama** justru supaya kesebelasnya tidak
mengulang nasib enam belas tiket berskema-tanpa-aplikasi yang hilang dari pandangan dua ronde.

| tiket | entitas | skema | penahan | catatan |
|---|---|---|---|---|
| `65` | `KELAS_BISNIS_LAYER` | **OK** `422` | — | **skema mendarat ronde 6.** FK ganda: `DETAIL_PROPORSIONAL` ikut hapus, `KELAS_BISNIS` tolak. ⚠️ **nol `UNIQUE`** — kunci alaminya belum bernomor |
| `66` | `BESARAN_LAYER` | — | `31` · `64` | **TIDAK dikerjakan ronde 6**, dan itu satu-satunya tiket Tugas 2 yang tertinggal. Ia menuntut tabel acuan **peran besaran**; tiket `64` yang memilikinya, dan `64` berstatus `tertahan` (penetapan pemilik). Membuat acuannya sendiri = mengambil keputusan yang `64` tahan |
| `67` | `KELOMPOK_LAYER` | **OK** `422` | — | **skema mendarat ronde 6.** `LAYER` ikut hapus, `KELOMPOK_TREATY` tolak. ⚠️ nol `UNIQUE` |
| `68` | `KELAS_BISNIS_KELOMPOK` | **OK** `422` | — | **skema mendarat ronde 6.** Kedalaman keempat, berdiri **terpisah** dari `65`. ⚠️ nol `UNIQUE` |
| `69` | `PENCAPAIAN` | **OK** `423` | — | **skema mendarat ronde 6.** Induk `KONTRAK`, hapus **tolak** — ERD HTML baris 6 menulis `LIMIT_DETAIL`+CASCADE dengan bukti yang ia sendiri tandai **DAUN-RELATIF**; SQL korpus menang. `CASH_CALL_KLAIM` ikut, tiga turunan tidak |
| `70` | `RINCIAN_ANGSURAN` | **OK** `424` | — | **skema mendarat ronde 6.** `TERMIN` berhenti menjadi induk tanpa anak. ⚠️ keempat kolomnya tanpa penulis hidup — penghalang tetap |
| `71` | `MIGRASI_KORELASI` | **OK** `428` | — | **skema + bukti Oracle ronde 8.** `TestTiket71KorelasiTanpaCapWaktuDitolak` · `TestTiket71KorelasiBertahanTanpaKunciAsing`. Nol kunci asing, dan uji positifnya membuktikan keputusan itu: baris korelasi **bertahan** sesudah kontraknya dihapus |
| `72` | `MIGRASI_PENDARATAN` | **OK** `428` | — | **skema + bukti Oracle ronde 8.** `TestTiket72PendaratanTanpaMuatanDitolak` · `TestTiket72MuatanRusakTetapMendarat` — muatan yang **tidak dapat diurai sama sekali** tetap mendarat |
| `73` | `MIGRASI_NILAI_DITOLAK` | **OK** `428` | — | **skema + bukti Oracle ronde 8.** `TestTiket73NilaiDitolakYatimDanPendaratanTakDapatDihapus` · `TestTiket73NilaiBukanAngkaTersimpanApaAdanya`. ⚠️ Perilaku hapus **`tolak`**, bukan `ikut hapus` — tiket 73 DIKOREKSI; ERD baris 38 menulis "di Go" |
| `74` | `ARSIP_MUATAN_KELUAR` | **OK** `425` | — | **skema mendarat ronde 6**, dan ia prasyarat `42` yang ikut mendarat. Hapus **tolak**: ERD baris 39 menulis "di Go", jadi ERD tidak meresepkan aturan basis data |
| `75` | `BATAS_WEWENANG` | — | `64` · **dokumen batas wewenang** | **tetap `tertahan`**, dan penahannya **tidak diubah**. Sapuan korpus 2 Okt 2026 atas `authority`/`BatasWewenang`/`ApprovalLimit`: **nol kemunculan**. ⚠️ Ia **ADA di ERD** — daftar *"7 tabel tanpa jalur Pega"* — jadi "tidak ada di ERD" bukan penahan yang sah |

> ⛔ **`71` `72` `73` dikerjakan SEBELUM `44`.** Tiket `44` memuat data pertama; sesudahnya tidak
> ada lagi jalan menelusuri baris baru kembali ke asalnya, dan tiket `61` menuntut justru
> penelusuran itu.

## ⭐ TIKET YANG KINI MEMENUHI UKURAN SELESAI — daftar untuk pemilik proses

**Ditulis 3 Oktober 2026, dan ini daftar yang belum pernah ada.** Ukuran selesai papan:
*"sebuah tiket selesai bila perilakunya **dapat gagal** dan **dapat diperiksa**"*, dengan uji
negatif **dan** positif.

⛔ **Nol tiket ditandai `selesai` di sini.** `status:` adalah medan di berkas tiket, milik pemilik
proses. Yang daftar ini lakukan: menunjukkan tiket mana yang buktinya kini ADA, berikut nama
ujinya, supaya keputusannya tinggal dibaca.

Seluruh uji di bawah berjalan terhadap **POOLDATA** lewat `make test-db-treatyin`, satu transaksi
per test yang diakhiri `Rollback`. **46 uji lulus, nol dilewati, nol gagal**, dan sesudahnya
**nol baris menetap**.

### A · Lapisan skema terbukti, dan kriteria tiketnya TIDAK menyebut lapisan aplikasi

Kesebelas tiket ini kriterianya seluruhnya tentang skema dan constraint. Buktinya kini lengkap.

| # | Uji NEGATIF | Uji POSITIF |
|---|---|---|
| `20` | `TestTiket20KursMataUangGandaDitolak` | `TestTiket20DuaMataUangPadaSatuVersiDiterima` |
| `22` | `TestTiket22RetensiCedantGandaDitolak` | `TestTiket22RetensiDuaMataUangDiterima` |
| `23` | `TestTiket23EgnpiGandaDitolak` | `TestTiket23EgnpiDuaMataUangDanKelasBisnisKosongDiterima` |
| `24` | `TestTiket24PortofolioGandaDitolak` | `TestTiket24PortofolioEmpatKombinasiDiterima` |
| `27` | `TestTiket27TerminGandaDitolak` | `TestTiket27TerminNomorSamaDuaMataUangDiterima` ⭐ |
| `28` | `TestTiket28SkalaKoasuransiGandaDitolak` | `TestTiket28TigaBarisSkalaKoasuransiDiterima` |
| `30` | `TestTiket30DokumenGandaDitolak` | `TestTiket30DokumenSamaPadaDuaVersiDiterima` |
| `34` | `TestTiket34DetailProporsionalGandaDitolak` | `TestTiket34DuaKelompokDalamSatuLayerDiterima` |
| `38` | `TestTiket38PenyebaranTanpaIndukDanBerindukGandaDitolak` · `TestTiket38PenyebaranJenisGandaDitolak` | `TestTiket38KeduaPelekatanPenyebaranDiterima` · `TestTiket38DuaJenisPadaSatuBagianDiterima` · `TestTiket38RantaiIkutHapusDuaTingkat` · `TestTiket38JenisReasuransiTerpakaiTidakDapatDihapus` |
| `39` | `TestTiket39JejakYatimDitolakDanVersiBerjejakTidakDapatDihapus` | `TestTiket39DuaJejakRuasSamaPadaVersiSamaDiterima` |
| `54` | `TestTiket54CatatanYatimDitolakDanVersiBercatatTidakDapatDihapus` | `TestTiket54DuaKeputusanOrangSamaPadaVersiSamaDiterima` |
| `01` ⟦Adjustment⟧ | `TestTiket01VersiDasarTakAdaDitolakDanVersiDasarTidakDapatDihapus` | `TestTiket01VersiPertamaTanpaDasarDiterima` |

⭐ **`27` layak dibaca dua kali.** Uji positifnya — termin nomor 1 dalam IDR dan termin nomor 1
dalam USD pada versi yang sama, **keduanya diterima** — adalah uji yang **gagal** pada rumusan
`INV-12` yang lama. Papan menuntutnya dijalankan, bukan diargumentasikan. Ia dijalankan.

⭐ **`38` membuktikan rantai `ikut hapus` DUA TINGKAT**: menghapus `PENYEBARAN` membawa serta
`RINCIAN_PENYEBARAN` dan `NILAI_PENYEBARAN`. Rantai yang putus di tengah meninggalkan nilai yatim
yang tidak terlihat sampai seseorang menjumlahkannya.

### B · Lapisan skema terbukti, tetapi kriterianya JUGA menyebut lapisan aplikasi

Ketiga tiket ini punya butir kriteria yang menyentuh jalur simpan atau jalur baca. Skemanya
terbukti; sisanya belum.

| # | Uji NEGATIF | Uji POSITIF | Yang MASIH kurang |
|---|---|---|---|
| `25` | `TestTiket25PeriodePelaporanGandaDitolak` | `TestTiket25EmpatPeriodePelaporanDiterima` | `INV-55` menunggu jalur simpan |
| `26` | `TestTiket26PeriodeAkumulasiGandaDitolak` | `TestTiket26AkumulasiKeduaDanKolomOpsionalKosongDiterima` | `INV-56` menunggu jalur simpan |
| `42` | `TestTiket42ArsipYatimDitolakDanKontrakBerarsipTidakDapatDihapus` | `TestTiket42MuatanRusakDanPengirimanKeduaDiterima` | lapisan aplikasinya **sudah** mendarat ronde 6; yang tersisa pemuatan arsip nyata, tiket `44` |

### C · Perkakas pemindahan — terbukti, dan BERTENGGAT

| # | Uji NEGATIF | Uji POSITIF |
|---|---|---|
| `71` | `TestTiket71KorelasiTanpaCapWaktuDitolak` | `TestTiket71KorelasiBertahanTanpaKunciAsing` ⭐ |
| `72` | `TestTiket72PendaratanTanpaMuatanDitolak` | `TestTiket72MuatanRusakTetapMendarat` ⭐ |
| `73` | `TestTiket73NilaiDitolakYatimDanPendaratanTakDapatDihapus` | `TestTiket73NilaiBukanAngkaTersimpanApaAdanya` |

⛔ **Ketiganya harus BERDIRI sebelum tiket `44` memuat data pertama** — dan kini ketiganya berdiri.

### D · Yang sudah punya buktinya sebelum ronde ini

`14` `15` `29` `31` `32` `33` `37` `05` — delapan, dan tidak berubah.

## Dua entitas yang TIDAK ditiketkan, dan itu keputusan

`RINGKASAN_LIMIT` (`T_TREATY_LIMIT_SUMMARY`) dan `REKAP_KONTRAK` (`T_TREATY_TOTAL`) diperiksa
2 Oktober 2026 dan diputuskan **turunan** — menyimpannya melanggar `INV-58`. Bukti yang menentukan:
sapuan atas setiap nama tabel yang disebut SQL di seluruh korpus mengeluarkan **enam belas** tabel,
dan **tidak satu pun** bernama demikian. Rinciannya berikut **syarat pembalikannya** di
`REKONSILIASI-XML-VS-DDL.md` §4.

## Dua lubang spec yang TIDAK ditambal

~~**`NILAI_SELISIH`**~~ — **ditambal 2 Oktober 2026: tiket `76` papan Adjustment.** Ia pernah
berdiri di sini sebagai lubang sebab `ERD.md` §2.6 menyatakannya mengikat sementara `ddl-usulan/`
dan `KAMUS-KOLOM.md` tidak memuatnya. Tiket `76` merancang bentuknya dari
`T_TREATY_VALUE_DIFFERENCE` dengan **mata uang di dalam kunci** (`ADR-0048` butir 3), dan menuntut
entrinya masuk `KAMUS-KOLOM.md` lewat definisi §10 — bukan dengan menyunting hasil bangkitannya.
Satu penghalang ikut lahir bersamanya: `BESARAN_DAPAT_DISESUAIKAN`, induk kunci asing keduanya,
**tidak ada di mana pun**.

**`PERISTIWA_KONTRAK`** — tabelnya ada di `ddl-usulan/` dan relasinya di `ERD.md` §2.3b, tetapi
**nol tiket dan nol kemampuan `P-nn` menghasilkannya**. Tidak disentuh.

## Sembilan tabel PENDARATAN — di LUAR ke-77 tiket

Ronde tabel tab Treaty In, 3 Oktober 2026 — delapan siang, kesembilan sore.
⚠️ **Sengaja tidak diberi nomor tiket**, sebab
kedelapannya bukan entitas model baru: ia mendaratkan larik di dalam `M_TREATY_IN.JSONDATA` supaya
layar berhenti mengurai CLOB. Menaruhnya di tabel ke-77 akan membuat penyebutnya berbohong.

| Tabel | Migrasi | DDL di Oracle | pemuat | uji bentuk | uji Oracle | rekonsiliasi |
| --- | --- | --- | --- | --- | --- | --- |
| `M_TREATYIN_REPORTINGPERIOD` | `430` `431` | **OK** | OK | OK | OK | OK |
| `M_TREATYIN_PORTFOLIO` | `430` `431` | **OK** | OK | OK | OK | OK |
| `M_TREATYIN_ACCUMULATION` | `430` `431` | **OK** | OK | OK | OK | OK |
| `M_TREATYIN_EGNPI` | `430` `431` | **OK** | OK | OK | OK | OK |
| `M_TREATYIN_RETENTION` | `430` `431` | **OK** | OK | OK | OK | OK |
| `M_TREATYIN_INSTALLMENT` | `430` `431` | **OK** | OK | OK | OK | OK |
| `M_TREATYIN_INSTALLMENTITEM` | `430` `431` | **OK** | OK | OK | OK | OK |
| `M_TREATYIN_COMMENT` | `430` `431` | **OK** | OK | OK | OK | OK |
| `M_TREATYIN_COINSCALE` | `432` `433` | **OK** | OK | OK | OK | OK |

⭐ **"DDL di Oracle = OK" di sini berarti BENAR-BENAR BERDIRI.** `go run ./cmd/api -migrate`
dijalankan dua kali 3 Oktober 2026: pukul 16:05 *"2 langkah dijalankan, 93 dilewati, 17
pernyataan"* (migrasi `430` `431`), dan pukul 17:28 *"2 langkah dijalankan, 95 dilewati, 2
pernyataan"* (migrasi `432` `433`). Dibuktikan ulang lewat katalog: sembilan tabel, sembilan
sequence `cycle=N`, satu kunci asing `FK_MTI_INSTALLMENTITEM_1` `delete=CASCADE`.

⚠️ **Ini MENGOREKSI catatan ronde 8 di bawah** — *"nol migrasi baru dijalankan terhadap
POOLDATA"* — yang benar **pada ronde itu** dan berhenti benar sekarang. Migrasi `422`–`429` juga
sudah tercatat di `T_MIGRASI` ketika `430` dijalankan; totalnya kini **95 langkah**.

**Tabel kesembilan yang TIDAK dibuat:** `M_TREATYIN_ACHIEVEMENT`. `AchievementLists` muncul **nol
kali** sebagai kunci puncak di 1.854 dokumen; jalurnya `Limits[].Detail[].AchievementLists`, dan
induknya di `M_TREATY_IN2`. Lihat `STRUKTUR-TABEL-TREATY-IN.md` dan
`4-erd-dan-tabel-datar/KOREKSI-ERD-VERSUS-POOLDATA.md` §3.

**Isi tabelnya hari ini: 27.238 baris, dari ke-1.854 kontrak** — 26.536 dari kedelapan tabel
pertama, ditambah **702** dari `M_TREATYIN_COINSCALE`. Dimuat dengan
`go run modul/treatyin/backend/pemuat/jalankan.go -semua -ikat`, dan tiap kali **didahului
jalannya yang kering atas 1.854 yang sama** dengan hasil identik: *"cocok 1.854 kontrak, 27.238
baris · selisih 0 · gagal 0"*. Dicocokkan sesudahnya lewat `-cocokkan` — kesembilan angkanya sama
dengan sapuan Python yang berdiri sendiri. Nol kunci JSON tanpa kolom ditemukan.

**Tiga tab sudah BERHENTI mengurai CLOB.** `ReportingPeriodList`, `Portfolio`, dan
`AccumulationList` kini dibaca `repository/pendaratan_baca.go` dari tabelnya;
`TestLarikYangSudahPunyaTabelTidakDiuraiLagi` menolak kembalinya ke `jsonWarisan`, dan
`TestIsiTabelSamaDenganIsiDokumen` membandingkan **219 nilai** dari 30 kontrak terpadat satu per
satu terhadap dokumennya. `CurrencyList` **tetap** di JSON — lihat
`4-erd-dan-tabel-datar/KOREKSI-ERD-VERSUS-POOLDATA.md` §12.

**DUA BELAS tab kini berisi.** Tujuh dari tabel pendaratan (Reporting Period · Portfolio ·
Accumulation · EGNPI · Maximum Retention · Installment · Information & Submit), **satu** dari
tabel pendaratan kesembilan (Co-Ins Scale), **empat** dari `M_TREATY_IN2` tanpa tabel baru
(Limits · Share · Event Limits · RNM Share), dan Rate of Exchange masih dari dokumen.

**EMPAT BELAS tab berisi sejak 4 Oktober 2026** — dua tab teks menyusul lewat Jalan B
(`KEPUTUSAN` §15): Exclusions dan Special Conditions, dibaca dari `JSONDATA`, **nol tabel baru**.

**Tiga tab masih `.trin__belum`, dan sebabnya BERBEDA:** Retro (**ditunda**, `KEPUTUSAN` §14 —
kedua kontraknya ditandai: `1000493`, `1000755`), Value Difference (objek, bukan teks — milik
tiket `06`/`11`/`13`), Achievement In IDR (induknya `Limits[].Detail[]`, bukan kontrak).

⚠️ **Ronde 4 Oktober 2026 nol migrasi dan nol pemuatan.** Ketiga jawaban pemilik proses tidak
menuntut tabel; yang berubah tampilan, pemilihan ejaan, dan penandaan.

## Ringkasan lapisan

| Lapisan | Cacah dari 77 | Perubahan ronde 8 |
| --- | ---: | --- |
| skema berdiri | **34** | **+7** — `38` ×3 tabel, `54`, `71` `72` `73` |
| kaskade sesuai ERD | **27** | **+8**; `TestPerilakuHapusSesuaiERD` kini mengadu **41** kunci asing |
| uji bentuk | **27** | +8 |
| uji Oracle (manual, belum otomatis) | 6 | — · **nol migrasi baru dijalankan** terhadap POOLDATA ronde ini |
| **services** | **14** | **+1** — tiket `42` |
| **handler** | **10** | **+1** — tiket `42` |
| **tiket di papan** | **77** | — |
| berstatus `selesai` | **0** | status milik pemilik proses |

⚠️ **27 dari 77, dan ke-27 tabelnya ada di Oracle hanya SEBAGIAN.** Ke-29 tabel lama berdiri di
POOLDATA sejak 2 Oktober; **enam tabel baru ronde ini BELUM dijalankan** — migrasi `422`–`426`
ditulis, diuji bentuknya, dan **tidak dijalankan**. Kolom "skema berdiri" menghitung **tiket yang
punya migrasinya**, bukan tabel yang berdiri di instans.

⛔ **Enam kunci alami menunggu nomor invarian**, dan keenamnya **bertenggat terhadap tiket `44`**:
`UNIQUE` yang dipasang di atas data yang sudah kembar akan gagal. Daftarnya di
`SPEC-MODEL-DATA.md` §13 — tagihan baru yang lahir ronde ini.

⚠️ **Ronde 5 menaikkan penyebut tanpa menaikkan pembilang** — 13 entitas yang hilang mendapat
tiketnya, nol mendapat tabelnya. **Ronde 6 menaikkan pembilangnya**: delapan dari 13 kini punya
migrasinya, dan `ERD-TREATY-IN-DAN-EDM.html` yang menjadi acuan strukturnya.
