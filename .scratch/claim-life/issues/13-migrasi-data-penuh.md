# 13: Migrasi data penuh Claim — Life

**Status:** ready-for-agent

**Blocked by:** 11 (kontrak Komite — jalur balik), **14 (skema relasional klaim — PREFACTOR)**

> ⚠️ **Dipersempit 2026-09-16 — revisi penyimpanan.** **Pemindahan bentuk data klaim
> (JSON + tabel flat warisan → enam tabel relasional, beserta normalisasi atribut polis) berpindah
> ke tiket 14**, yang menjadi PREFACTOR. Tiket ini **tinggal** menangani hal yang tidak menyentuh
> bentuk: keutuhan riwayat setelah cutover, penomoran, pembacaan master, dan kolom yang tidak
> ditafsirkan. AC yang berpindah ditandai di bawah.

## Hasil & nilai pengguna

Sebagai **work owner**, saya ingin seluruh data klaim Life dipindahkan — beserta **semua barisnya** —
sehingga tidak ada pekerjaan tertinggal di Pega dan riwayat putaran Komite tidak hilang.
*(User story 40, 41, 42 di spec)*

## Area codebase

`internal/repository` (pembacaan sumber + penulisan sasaran), skrip migrasi + DDL produksi di dalam
`OUTPUT_HASIL_RNM/`. Tidak menyentuh `handlers`/`frontend`.

## Rule Pega sumber

| Rule | Identitas | Dipakai untuk |
| --- | --- | --- |
| `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `[terverifikasi]` 55 nama kolom `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` terbaca langsung dari blok PL/SQL |
| `Claim Life/RDBList/GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETSEQUENCENUMBER_SQL` / `RULE-CONNECT-SQL` | penomoran yang tidak boleh melompat/mengulang setelah migrasi |

`[terverifikasi]` Tabel akseptasi Life **terpisah** dari lini lain: `OS_AKSEPTASI_KLAIM_LIFE`
(7 kemunculan, hanya `Claim Life` + `Komite Claim Life`) versus `OS_AKSEPTASI_KLAIM` (191 kemunculan,
enam modul non-Life). Migrasi ini **tidak menyentuh** tabel non-Life.

## ADR terkait

**ADR-0009** (migrasi penuh; koeksistensi **ditolak**; cutover memutus, bukan bertahap),
**ADR-0011** (seluruh **baris** ikut pindah, bukan hanya keadaan terakhir),
**ADR-0006** (penomoran), **ADR-0007** (jejak audit tidak dapat direkonstruksi ke belakang).

## Acceptance criteria

⚠️ **Berpindah ke tiket 14** (jangan dikerjakan di sini): seluruh klaim + seluruh baris adjustment
terbawa; jumlah baris per klaim utuh; presisi uang; desimal presisi arbitrer; migrasi dapat
dijalankan ulang; normalisasi atribut polis. Rujukan silang: AC 51–54 spec.

- [ ] Klaim yang sedang berada di tengah siklus terbawa **beserta statusnya** dan posisi tahapnya.
      ⚠️ **Diselaraskan:** posisi tahap mendarat di **`T_WORK_CLAIM`** (spec §2b), bukan di header
      klaim. *(AC 46 spec; penyimpangan sadar 6)*
- [ ] Setelah migrasi, penomoran klaim **tidak melompat dan tidak mengulang**. *(AC 42 spec)*
- [ ] Tidak ada periode dua penulis ke rekam akseptasi Life dari sisi Claim — Life. *(ADR-0009)*
      ⚠️ **Diselaraskan:** setelah cutover, penulis satu-satunya adalah keenam tabel klaim baru;
      `OS_AKSEPTASI_KLAIM_LIFE` **tetap ditulis** — ⚠️ **koreksi 2026-09-16**
      `[keputusan work owner]`: setelah cutover, penulis klaim adalah **kedelapan tabel relasional
      baru** *dan* `INSERT` flat ke `OS_AKSEPTASI_KLAIM_LIFE`, karena hilir masih membaca dari sana.
      Yang **dibuang hanya JSON**. Test yang **menolak** penulisan `OS_AKSEPTASI_KLAIM_LIFE` justru
      **gagal**. *(AC 32 spec)*
- [ ] ⚠️ Master `RATE_LIFE`, `PRODUCTINWARD_LIFE`, `CURRENCY` diperlakukan sebagai **view atas
      `JSONDATA`**, bukan sebagai tabel relasional biasa. **Pengecualian 2026-09-16:** **produk**
      dibaca dari **`product_life` relasional** (hasil migrasi Master Product Name Life), bukan dari
      `m_product_life.JSONDATA`. *(AC 38 spec)*
- [ ] `NO_SEQ` terjaga **per kombinasi `(CLASS, JENIS, TAHUN)`**, bukan sebagai penghitung global.
- [ ] Kolom `LAYER_1`…`LAYER_4` dipindahkan apa adanya **tanpa ditafsirkan** — perannya belum
      terverifikasi.

## Catatan penutupan (2026-09-14)

**OQ-001 TERTUTUP** untuk Claim — Life `[data DBA]` — DDL **12 tabel/view** diserahkan, mencakup
**seluruh** persistensi Claim Life. Daftar lengkap dan tipenya di `discovery/open-questions.md`
OQ-001. Yang mengikat tiket ini:

| Temuan | Konsekuensi untuk migrasi |
| --- | --- |
| Kolom uang = Oracle **`NUMBER` tanpa presisi** | Go **wajib** desimal presisi arbitrer; **`float64` dilarang** (**ADR-0003**) |
| `STS_REJECT NUMBER(38)` di tabel klaim | ⚠️ berbeda dari `STS_REJECT VARCHAR2(15)` di `EMAILKOMITE` — nama sama, tabel berbeda, tipe berbeda |
| `RATE_LIFE`, `PRODUCTINWARD_LIFE`, `CURRENCY` = **VIEW atas `JSONDATA`** | master sesungguhnya **JSON**; membacanya sebagai tabel relasional biasa akan menyesatkan |
| `TANGGAL_CLOSING.TANGGAL VARCHAR2(10)` | tanggal disimpan sebagai **string**, bukan `DATE` |
| `LAYER_1`…`LAYER_4 VARCHAR2(10)` | ⚠️ **peran belum terverifikasi** — ikut dipindah apa adanya, jangan ditafsirkan |

**`AdjustmentList` tidak punya tabel fisik** `[terverifikasi]` — kelasnya
`ASM-FW-GISFW-Data-AdjustmentLife` (berawalan `Data-` = embedded, tanpa tabel), induk
`ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL`. Bukti:
`Claim Life/Activity/SaveOutStandingLife_Act.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` /
`SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`). **Baris adjustment di-persist ke
`POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`** — jadi "memindahkan seluruh baris" berarti memindahkan seluruh
baris tabel itu, bukan mencari tabel adjustment yang tidak pernah ada.

**OQ-002 TERTUTUP** — sequence per `(class, jenis, tahun)` di `GENERATE_SEQUENCE_NUMBER`
(PK komposit, `TAHUN VARCHAR2(5)`, `NO_SEQ NUMBER`, `MM_YYYY VARCHAR2(10)`). Migrasi nomor kini
dapat direncanakan: yang harus dijaga adalah **`NO_SEQ` per kombinasi kunci**, bukan satu penghitung
global.

`[terbuka]` **OQ-018** — lingkungan `pega-nusre` belum dinyatakan (IT-infra). **Tidak memblokir
penulisan skrip migrasi**; memblokir **pemilihan sumber data** saat eksekusi.

`[terbuka]` **OQ-013 untuk modul non-Life** — tidak menyentuh migrasi Claim Life.

## Catatan

`[keputusan work owner 2026-09-14]` Opsi koeksistensi — sistem baru hanya menerima klaim baru
sementara klaim berjalan diselesaikan di Pega — **ditolak**. Penolakan itu dicatat di **ADR-0009**
supaya tidak diusulkan ulang.

`[keputusan work owner]` Jejak audit lama **tidak dapat direkonstruksi**: data sebelum cutover hanya
punya `CREATEOPNAME` + empat kolom tanggal. Riwayat transisi lengkap hanya ada setelah cutover.

## Perintah verifikasi

```
go test ./internal/...
make check
```

Tambah verifikasi khusus migrasi: hitungan klaim dan hitungan baris per klaim sebelum dan sesudah
pemindahan harus sama.
