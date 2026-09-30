# 00: Kolom EDM + `PARENT_ID` pada tabel polis — **PREFACTOR (ringan)**

**Status:** ready-for-agent

**Blocked by:** **PremiumList Life tiket `00`** (skema tujuh tabel — tabel yang di-ALTER harus ada
lebih dulu)

⚠️ **PREFACTOR dan tiket PERTAMA konteks ini.** Diberi nomor `00` supaya berada di depan tanpa
menomori ulang dua belas tiket yang sudah terbit. **Seluruh tiket 01–12 kini memblokir pada tiket
ini.**

⚠️ **Ini tiket RINGAN — `ALTER`, bukan skema baru.** `[keputusan work owner]` Endorsement **berbagi
tujuh tabel PremiumList Life**; ia **tidak** punya tabel sendiri. Yang ditambahkan hanya kolom
pembeda versi dan satu pointer.

## Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin tabel polis mampu menampung **versi endorsement** — kolom yang
membedakan endorsement dari new business, dan pointer yang memasangkan tiap peserta dengan dirinya
di versi sebelumnya — sehingga endorsement dapat disimpan tanpa tabel tambahan dan tanpa kehilangan
jejak. *(Spec §16)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | `ALTER TABLE` dua tabel + index; backfill `PRODKE` untuk baris lama |
| `internal/models` | Bentuk polis bertambah atribut versi; peserta bertambah pointer induk |

## Yang ditambahkan

**`T_PREMIUM_LIST`** — kolom EDM, **seluruhnya nullable**, kosong pada baris new business
`[terverifikasi dari pyFields `DATA_JSON` EDM nyata]`:

| Kolom | Tipe | Catatan |
| --- | --- | --- |
| `EDM_TYPE` | teks | **`1`=Perubahan Data, `3`=Batal** — dropdown `EDMTYPE`. Nilai `2` **memang tidak ada** |
| `OLD_POLICY_NO` | teks | nomor polis yang di-endorse |
| `EDM_DATE` | **DATE** | |
| `EDM_NOTE` | teks | |
| `TYPE_CEDING` | teks | `1`=QS, `2`=SURPLUS, `3`=QS+SURPLUS, `4`=XOL — dropdown korpus |
| `PREMI_PROPOSED`, `UANG_PERTANGGUNGAN` | **NUMBER** | **ADR-0003** |
| `JENIS_PRODUK`, `SISTEM_REASURANSI` | teks | |
| `STATUSS`, `STATUS_UPDATE`, `STATUS_SERVICE` | teks | status proses endorsement |
| `START_DATE`, `END_DATE` | **DATE** | |
| `NOENDORS` | teks | `[terverifikasi]` `Generate_NoEndorsmentLife` |
| `PL_NUMBER_EDM` | teks | |
| **`PRODKE`** | bilangan bulat | ⚠️ **versi berjalan = `PRODKE` terbesar** `[terverifikasi]` |
| `EDMSTATUS`, `STATUSOLD` | teks | `Old` / `New` / `Delete` / `Batal` `[terverifikasi]` |

**`T_PREMIUM_LIST_DETAIL`** — satu kolom:

| Kolom | Tipe | Catatan |
| --- | --- | --- |
| **`PARENT_ID`** | FK self-reference → `T_PREMIUM_LIST_DETAIL.ID`, **nullable**, **ber-index** | ⚠️ pointer ke peserta versi sebelumnya; `NULL` untuk new business dan peserta baru |

⚠️ **Hanya tabel peserta yang mendapat `PARENT_ID`.** `[keputusan work owner]`
`T_PREMIUM_LIST_SPREADING` dan `_SPREADING_RETRO` **tidak** — keduanya ikut tersalin di bawah
peserta versi baru, dan FK ke induknya sudah cukup. `T_PREMIUM_LIST_SUMMARY` tidak disalin sama
sekali.

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `MappingEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `MAPPINGEDMLIFE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/MappingEDMLife.xml` | ⚠️ pencocokan **loop indeks** yang digantikan `PARENT_ID` |
| `Generate_NoEndorsmentLife` | `Endorsement Life/` | | sumber `NOENDORS` |
| `InsertJsonPolisEDM` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLISEDM` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/InsertJsonPolisEDM.xml` | ⚠️ **dibuang** — penulis JSON |

⚠️ **Penyimpangan sadar 10 — `PARENT_ID` menggantikan loop indeks.** `[terverifikasi]`
`MappingEDMLife` mencocokkan peserta lama↔baru dengan membandingkan
`PremiumListDetail(Local.idxLocationEDM)` terhadap `(Param.EDMList = .pxListSubscript)` — **rawan
salah pasang** begitu urutan baris bergeser, dan **kunci pencocokannya tidak terbaca dari korpus**.
Pointer eksplisit menghapus kelas bug itu.

## ADR terkait

**ADR-0003** (uang non-float), **ADR-0006** (identitas lewat sequence), **ADR-0009** (migrasi penuh).

## Acceptance criteria

- [ ] ⚠️ `T_PREMIUM_LIST` memuat **seluruh kolom EDM di atas**, **seluruhnya nullable**. Baris new
      business tetap dapat ditulis tanpa mengisi satu pun di antaranya. *(AC 57 spec)*
- [ ] `EDM_TYPE` hanya menerima **`1`** dan **`3`**; nilai lain **ditolak** di lapisan layanan.
      *(AC 58 spec; `[terverifikasi]` dropdown `EDMTYPE`)*
- [ ] ⚠️ `T_PREMIUM_LIST_DETAIL` memuat **`PARENT_ID`** sebagai FK self-reference, **nullable**, dan
      **ber-index**. *(AC 60 spec; penyimpangan sadar 10)*
- [ ] ⚠️ `T_PREMIUM_LIST_SPREADING`, `_SPREADING_RETRO`, dan `_SUMMARY` **tidak** mendapat
      `PARENT_ID`. Test yang menemukannya **gagal**. *(spec §16)*
- [ ] ⚠️ **Tidak ada tabel baru** dibuat untuk endorsement. Test yang menemukan tabel header khusus
      endorsement **gagal**. *(AC 55 spec; penyimpangan sadar 9)*
- [ ] Versi berjalan sebuah polis dapat ditemukan sebagai baris ber-**`PRODKE` terbesar**; seluruh
      versi **hidup berdampingan**. *(AC 56 spec)*
- [ ] Baris polis lama hasil migrasi mendapat `PRODKE` yang benar, sehingga "versi berjalan"
      terbaca sejak hari pertama.
- [ ] `PARENT_ID` yang menunjuk peserta **tidak ada** ditolak basis data — integritas rujukan
      ditegakkan, bukan diandaikan.
- [ ] Nilai uang tambahan (`PREMI_PROPOSED`, `UANG_PERTANGGUNGAN`) bertipe **desimal presisi
      arbitrer** dan **menerima nilai negatif** — jurnal balik menulis minus. *(**ADR-0003**;
      spec §6, §16)*
- [ ] Tanggal tambahan (`EDM_DATE`, `START_DATE`, `END_DATE`) bertipe **`DATE`**.
- [ ] `ALTER` dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.

## Blocker

**Tidak ada pemblokir.** **OQ-001 ditutup 2026-09-16** — tabel yang di-ALTER dirancang sendiri di
PremiumList Life tiket `00`.

⚠️ `[terbuka]` **tidak memblokir:** format `NOENDORS` (`Generate_NoEndorsmentLife`, DBA) dan format
identitas kerja **`EDMLF-<n>`** (`[fakta bisnis — work owner]`, **nol kecocokan di korpus**).
Keduanya kolom teks; formatnya ditetapkan di tiket **04**.

## Catatan

⚠️ **Mengapa ini ringan padahal PREFACTOR.** `[keputusan work owner]` Endorsement **bukan entitas
baru** — ia **versi** dari polis yang sama. Karena itu tidak ada tabel yang dibuat, tidak ada data
yang dipindah, dan tidak ada bentuk yang berubah. Yang ditambahkan hanya kolom pembeda dan satu
pointer. Ini kebalikan dari tiket `00` PremiumList Life, yang membangun tujuh tabel dari nol.

⚠️ **Nama berbohong (OQ-066).** `[terverifikasi]` `InsertJsonPolisEDM` benar-benar menulis blob
JSON — **dibuang**. Bandingkan dengan Claim Life, di mana `InsertJsonKlaimLife_sql` justru `INSERT`
flat. **Baca kodenya, jangan namanya.**

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — integritas rujukan `PARENT_ID` dan
nullability kolom EDM **hanya berperilaku benar pada basis data sungguhan**.

```
go test ./internal/...
cd frontend && npm test
make check
```
