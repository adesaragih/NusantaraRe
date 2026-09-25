> Modul  : Treaty In Adjustment · cabang K — pemetaan ke artefak to-spec induk
> Dibuat : 2026-09-26
> Sifat  : **PEMETAAN SAJA. To-spec belum dimulai, dan tidak boleh dimulai dari berkas ini.**
> Dasar  : GRL-01 — spesifikasi Adjustment ditulis **ke dalam** spesifikasi induk, bukan di sampingnya

# Cabang K — ke mana setiap keputusan mendarat

## 0. Apa yang berkas ini bukan

Ia **bukan** spesifikasi, **bukan** draf, dan **bukan** izin memulai to-spec. Ia daftar alamat:
untuk setiap artefak to-spec induk, keputusan GRL dan bahan to-spec mana yang akan ditulis ke sana.
Gunanya agar sesi to-spec **tinggal menulis, bukan mencari**.

Larangan sesi grilling tetap berlaku sepenuhnya: tidak ada spesifikasi, DDL, struct Golang,
komponen React, rancangan endpoint, maupun tiket yang ditulis dari sini.

## 1. Bagaimana to-spec induk disusun

Dibaca dari `SPEC-MODEL-DATA.md`, yang menyatakan dirinya **"BENTUK AWAL — hasil langkah 0"** dan
menuangkan langkah 1–7; **langkah 8 (invarian), 9 (peta telusur), dan 10 (ERD) belum dijalankan**
untuk induk. Artefaknya lima, dan urutannya mengikat:

| # | Artefak | Perannya |
|---|---|---|
| 1 | `SPEC-MODEL-DATA.md` | entitas, atribut, tipe, boleh-kosong |
| 2 | `SPEC-INVARIAN.md` | aturan yang harus selalu benar, beserta bentuk penegakannya |
| 3 | `UJI-NEGATIF-INVARIAN.md` | uji yang **harus gagal** untuk tiap invarian |
| 4 | `PETA-TELUSUR-JSON.md` | telusur jalur `JSONDATA` lama ke atribut baru |
| 5 | `4-erd-dan-tabel-datar/` | ERD, struktur data, tabel datar |

## 2. Pemetaan — per artefak

### 2.1 `SPEC-MODEL-DATA.md`

| Yang mendarat | Dari | Bentuk mendaratnya |
|---|---|---|
| addendum adalah `VERSI_KONTRAK`; `PENYESUAIAN` tidak kembali | **GRL-01**, **GRL-04** | tidak ada entitas baru; §10.2 dipakai apa adanya |
| `JENIS_ADDENDUM` — bernilai **dua** untuk versi baru; `EDMState` warisan dibawa apa adanya | **GRL-13**, **GRL-18** | nilai enum pada §10.2, dan kolom nilai warisan |
| materialitas — **atribut MASUKAN**, bukan turunan; beku sejak `AJUKAN`, berjejak selama `DRAFT` | **GRL-20**, **KTV-1** | `SIFAT_MATERIAL_ADDENDUM` pada §10.2 **kembali**; `GRL-18` bagian yang menghapusnya **gugur** |
| `ActualValue` — **dibuang seluruhnya**; premi aktual menjadi nilai versinya sendiri | **GRL-14** | **tidak ada** entitas anak baru; 108 jalur cermin tidak melahirkan kolom |
| `ID_VERSI_KONTRAK_DASAR` wajib terisi pada versi penyesuaian | **B-3** | kolom "boleh kosong" pada §10.2 |
| versi berlaku **tidak** menjadi kolom di `KONTRAK` | **GRL-11** | catatan turunan, bukan atribut |
| cedant disimpan sebagai rujukan | KONFIRMASI — `STRUKTUR-DATA.md` §3 `CEDANT [luar]` | tidak ada yang berubah |
| kolom asal-usul untuk baris warisan | **GRL-17** | `NOMOR_URUT_VERSI` warisan **diberikan ulang menurut kronologi**; pengenal lamanya dipelihara |
| **`DOKUMEN_ADDENDUM` — entitas BARU di atas versi** | **GRL-19** | **entitas baru**; `VERSI_KONTRAK` memperoleh rujukan **boleh kosong** |
| tanggal berlaku dokumen | **KTV-2** | kolom **nullable** pada `DOKUMEN_ADDENDUM`; dicabut bila `DB-16b` dibantah |
| **tidak ada** penanda "dikirim ke luar" | **KTV-3** | tidak ada kolom; satu nilai enum ditambahkan bila `DB-16a` dibantah |

### 2.2 `SPEC-INVARIAN.md`

| Yang mendarat | Dari | Nomor yang disentuh |
|---|---|---|
| nomor urut versi yang bentrok **ditolak**, pesannya menyebut nomornya | **B-1** | memperkuat **INV-04** |
| pengenal tampilan versi baru tidak boleh sama dengan pengenal warisan mana pun di kontrak yang sama | **B-2** | invarian baru |
| `ID_VERSI_KONTRAK_DASAR` wajib terisi pada versi penyesuaian, wajib kosong pada versi pertama | **B-3** | invarian baru |
| yang ditunjuk harus versi berlaku terakhir, bukan `DITOLAK`/`DIBATALKAN` | **B-4** | invarian baru |
| keadaan bertambah `DIBATALKAN`; tidak ada perpindahan keluar dari keadaan terminal | ADR-0055 perubahan 24 Sep | **INV-20**, **INV-23**, **INV-25** |
| rumus versi berlaku dan kasus yang dapat gagal | **GRL-11** | invarian baru + uji |
| kunci padanan baris selisih | **E1** (belum terkunci) | invarian per entitas anak |

### 2.3 `UJI-NEGATIF-INVARIAN.md`

| Uji yang harus gagal | Dari |
|---|---|
| simpan dua versi bernomor urut sama dalam satu kontrak | **B-1** |
| simpan versi penyesuaian tanpa `ID_VERSI_KONTRAK_DASAR` | **B-3** |
| tunjuk versi yang sudah `DITOLAK` sebagai dasar | **B-4** |
| kontrak versi 3 `DISETUJUI` + versi 4 `DITOLAK` → versi berlaku harus **3** | **GRL-11** |
| kontrak yang hanya punya versi pertama `DISETUJUI` → versi berlaku = versi pertama | **GRL-11** |
| ubah cedant, tanggal mulai, atau tanggal berakhir pada versi penyesuaian | ADR-0040 butir 3 + **NA-17b** |

### 2.4 `PETA-TELUSUR-JSON.md`

| Yang mendarat | Dari |
|---|---|
| 33 akar `ValueDifference` dan 165 jalur daunnya | `PENGETAHUAN.md` §5.6 |
| `OLDDATA` **tidak** menjadi pohon kedua; sisi lama lewat `ID_VERSI_KONTRAK_DASAR` | **GRL-01**, `STRUKTUR-DATA.md` §1.1 |
| `EDMState`, `EDMMaterialType`, `RevisionState`, `ViewState`, `IsEditData`, `Position`, `StatusAkseptasi` — nasibnya | **GRL-08**; `Position` dan `StatusAkseptasi` **dibaca saat migrasi, tidak disimpan** |
| `OLDID` warisan disimpan apa adanya | **GRL-09**, **GRL-10** |

### 2.5 `4-erd-dan-tabel-datar/`

| Yang mendarat | Dari |
|---|---|
| tidak ada relasi baru antarmodul; `ID_VERSI_KONTRAK_DASAR` sudah ada | KONFIRMASI — `ERD.md` §2.6 |
| tidak ada percabangan rantai versi | **GRL-10** |
| `DOKUMEN_KONTRAK` dirujuk, tidak dimiliki | KONFIRMASI — `STRUKTUR-DATA.md` §1.2 |
| penguncian optimistis dengan penanda versi baris | **R-F1** |

### 2.2a Invarian baru ronde D

| Yang mendarat | Dari | Bentuk mendaratnya |
|---|---|---|
| **`INV-69`** — versi **Non Material** tidak boleh punya baris `NILAI_SELISIH` bertipe uang atau porsi | **GRL-20** | baris invarian di `SPEC-INVARIAN.md`, golongan **CONSTRAINT** lewat *materialized view* — ia lintas baris |
| **`INV-70`** — versi **Material** tidak mengubah teks kesepakatan dan identitas kontrak terhadap versi dasarnya | **GRL-20** | idem |
| uji negatif **dan positif** `N-7a` … `N-7c+` | **GRL-20** | `UJI-NEGATIF-INVARIAN.md`, bagian baru |

> **`N-7c+` adalah uji positif yang menuntut versi WARISAN yang melanggar tetap dapat DIMUAT.**
> Tanpa ia, invariannya akan menolak data lama pada hari peralihan — dan `UA-3` yang mengukur
> berapa banyak.

### 2.6 Bahan to-spec yang sudah terkunci — peta lengkapnya

| Bahan | Dari | Artefak induk · bagian |
|---|---|---|
| `T-1` … `T-3` | `GRL-18` | `SPEC-MODEL-DATA.md` §10.2 — jenis sebagai masukan, titik bekunya, jejaknya |
| `B-1`, `B-2` | `GRL-09` | §10.2 — pengenal tampilan; bentuk turunan untuk versi baru, warisan apa adanya |
| `B-3`, `B-4` | `GRL-10` | §10.2 — `ID_VERSI_KONTRAK_DASAR`, disimpan **tepat di satu tempat** |
| `C-1` | `GRL-12` ~~batal~~ → **`GRL-20`** | §10.2 — `SIFAT_MATERIAL_ADDENDUM` **kembali** sebagai masukan |
| `C-2` | `GRL-13` | §10.2 — enum jenis bernilai dua |
| **`C-3`** | `GRL-14` | §10 — **tidak ada** entitas `ActualValue`; dan `TDA-16` hilang bersamanya |
| `I-1`, `I-2` | `GRL-11`, `GRL-17` | §10.2 dan rencana migrasi — versi berlaku turunan; nomor urut warisan diberikan ulang |
| `E3a` | `GRL-15` | §10.3 — atribut pro rata dibawa, **mesinnya tidak dibangun** |
| `E3b` | `GRL-16` | `PETA-TELUSUR-JSON.md` — kelas cacat share fakultatif |

## 3. Yang mendarat di luar artefak to-spec

| Tujuan | Isi |
|---|---|
| `DAFTAR-ESKALASI-MANAJEMEN.md` induk | butir 1 sudah dikoreksi; pertanyaan wewenang tombol "(dev)" |
| `USULAN-REVISI-ADR.md` | `REV-1`, `REV-2`, `REV-3`, `IND-1`, `IND-2`, `IND-3` |
| daftar untuk dibantah | `DB-1` … `DB-14` |
| kueri data | `UA-1` … `UA-17` |
| rencana peralihan | `R-D1`, `R-I1` |

### 3.1 `DOKUMEN_ADDENDUM` — kolom yang KOSONG untuk seluruh baris warisan

| | |
|---|---|
| **Apa** | `GRL-19` melahirkan entitas dokumen. `TD-01` membuktikan sistem lama **tidak pernah merekam nomornya** — nol atas 28 nama calon, kalibrasi lulus |
| **Kenapa ia TIDAK mendarat di artefak to-spec** | tidak ada yang perlu ditulis di spesifikasi selain kolomnya sendiri. Yang perlu direncanakan adalah **pengisiannya**, dan itu **pekerjaan orang dari arsip kertas** — bukan pekerjaan migrasi, bukan pekerjaan sistem |
| **Ke mana ia mendarat** | **rencana peralihan**, bukan spesifikasi |
| **Yang menagih** | laporan pertama yang mengelompokkan versi menurut dokumennya, dan mendapati **seluruh warisan** berkelompok *"tanpa dokumen"* |

### 3.2 Tiga kolom `TDA-17` yang tidak punya rumah — dan ini milik INDUK

`DEDUCTIBLE2`, `PREMIUM_EARNED`, `ROL_PCT` **tidak punya kolom di skema baru**, dan sebabnya `L-8`:
ketiganya muncul di peta telusur **hanya di dalam pohon cermin**, tidak pernah di pohon utama.

> **Tidak dapat diputuskan di modul Adjustment.** Ia lubang pada §10 **induk**, dan pemiliknya sesi
> to-spec Treaty In. Dicatat di sini supaya sesi to-spec Adjustment **tidak mengira** ketiganya
> sengaja dibuang.

## 4. Prasyarat masuk to-spec

1. ~~**Ketujuh butir GRILL terkunci** — C1, C2, C3, E1, E3, I1 ✅, I2.~~
   **TERPENUHI 24 September 2026** — seluruhnya terkunci lewat `GRL-11` … `GRL-20`, ditambah ronde D
   yang membatalkan `GRL-04` dan `GRL-12` dan menggantinya.
2. **Paket `REV` DISERAHKAN** kepada pemilik ADR.

   > **BUNYINYA DIUBAH 24 September 2026 — `KTV-4`.** Semula berbunyi *"diserahkan **dan
   > ditanggapi**"*. Keenam REV diperiksa terhadap pertanyaan *"apakah ia menentukan letak
   > kolom?"* dan **nol dari enam** lolos: seluruhnya merevisi **alasan** dan **daftar keadaan**,
   > bukan bentuk data. `REV-4` bahkan sudah didahului `GRL-20`, dan `REV-5` oleh `GRL-18`.
   >
   > **Tanggapannya DITAGIH sebelum to-ticket**, bukan sebelum to-spec. Bila pemilik ADR menolak
   > salah satu REV **dengan akibat pada bentuk data**, bagian to-spec yang bersandar padanya
   > dibuka kembali — dan yang paling mungkin `REV-3`, sebab keadaan **adalah** kolom.
   >
   > Dasar dan syarat pembalikannya lengkap di
   > [`KEPUTUSAN-TANPA-VERIFIKASI-ADJUSTMENT.md`](KEPUTUSAN-TANPA-VERIFIKASI-ADJUSTMENT.md) `KTV-4`.
3. **`EXP-1` ditutup** ✅ (26 September 2026).
4. Setiap butir KONFIRMASI, REKOMENDASI, dan DITUNDA sudah punya tempat — dipenuhi
   `PEMILAHAN-SISA-GRILLING.md`.

Butir 2 adalah satu-satunya yang **tidak berada di tangan sesi ini**.
