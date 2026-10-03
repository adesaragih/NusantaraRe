# 13: Rantai perhitungan uang — karantina P18 DIBUKA, sisa dua penahan lain

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Kedua penahan gugur. **P30** di luar lingkup Treaty In — seluruh rantai pemanggilnya bekerja pada `OfferFacIn`, nol `PolicyTreatyIn`. **P8** dicabut — bentuk muatannya terbaca dari saudara kelas induk: empat medan, dua panggilan keluar, surat bila gagal.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 6 dan 7.

---


**Status:** selesai *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)* — ⚠️ **bukan lagi karena P18**)*
~~**Blocked by:** ⛔ **P30** *(aturan penjumlah yang hilang)* · ⛔ **P8** *(muatan efek keluar)*~~ ⛔ **penahan gugur 23-09-2026**

> ⭐⭐ **SISI P18 DIBUKA.** `[penyimpangan sadar]` 2026-09-22 — dua baris ini semula berbunyi:
> > *"# 13: Rantai perhitungan uang — dikarantina sampai bahannya lengkap · **Status:** blocked ·
> > **Blocked by:** ⛔ **P18** *(isi 268 langkah)* · ⛔ **P8** *(muatan efek keluar)*"*
>
> **P18 ditarik oleh tim migrasi.** Isi 268 langkah ada di ekspor — **707 pasangan nama=nilai
> terisi pada 269 langkah**, nol tanda terpotong. Lihat `VERIFIKASI-P18.md`.
>
> ⚠️ **Tiket ini TETAP `blocked`**, tetapi sebabnya berubah: **P30** dan **P8**, bukan P18.
> `[terverifikasi]` `SumTSIPremiSpreadRNMMultiCob_Act` **tetap tidak ada di satu pun dari 21 modul**
> korpus — diperiksa ulang 2026-09-22. ⛔ **P8 tidak dibuka di ronde ini**: keputusan
> "ikuti apa adanya" sudah ada, tetapi sisa `[terbuka]`-nya — ARASAPAS itu sistem apa — masih milik
> `[Product+Underwriting]`.
**Menutup:** AC 79 · 87 *(2 AC)* — US 46

## Hasil & nilai pengguna

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kedua paragraf berikut dikutip utuh lalu
> ditarik:
> > *"Hari ini seluruh aritmetika uang … tinggal di dalam **langkah penetapan nilai** yang ⛔
> > **isinya tidak ikut dalam ekspor** yang diterima tim migrasi. ⭐ Bentuknya diketahui;
> > **rumusnya tidak**. … ⛔ **Tiket ini belum dapat dikerjakan**, dan ⛔ **tidak boleh dikerjakan
> > dengan rumus yang ditebak.**"*

⭐⭐ **Rumusnya diketahui.** `[terverifikasi]` Seluruh aritmetika uang tinggal di **langkah
penetapan nilai**, dan isinya **terbaca penuh** di ekspor — termasuk pembagi **102,2** untuk
`TypeTax="Inclusive"`, **PPH 2 %**, **PPN 2,2 %**, dan **presisi 8 angka** pada tiap `@divide`.

⚠️ **Yang masih menahan tinggal dua, dan keduanya bukan tentang rumus:** satu aturan penjumlah
yang **hilang dari seluruh korpus** *(P30)*, dan **muatan** efek keluar terakhir *(P8)*.
⛔ **Tetap tidak boleh dikerjakan dengan rumus yang ditebak** — sekarang tidak perlu menebak.

## Blocker

| Butir | Yang ditunggu | Pemilik |
| --- | --- | --- |
| ⭐ ~~**P18**~~ | ~~isi **268** langkah penetapan nilai yang terjangkau~~ — ⭐ **DITARIK 2026-09-22** | ~~`[pemilik export Pega]`~~ |
| ⛔ **P30** | satu aturan penjumlah yang **dipanggil tetapi tidak ada di seluruh korpus 21 modul** | `[pemilik export Pega]` |
| ⛔ **P8** | **muatan** yang dikirim ke layanan luar di akhir alur | `[Product+Underwriting]` |

## Yang SUDAH diketahui bentuknya

⭐ **Disebut supaya tidak dicari ulang.** ⚠️ Bab ini semula berjudul maksud *"apa yang TIDAK
menunggu P18"*; sesudah P18 ditarik, **tidak ada lagi yang menunggu P18.**

| Hal | Keadaan |
| --- | --- |
| medan | delapan medan uang **berpasangan mata uang** |
| sumber | ⭐ **kolom tersimpan**, bukan hasil hitung |
| arah aliran | ⭐ **dibaca, bukan dihasilkan** |
| ⭐ penyajian | ⭐ **sudah dapat dibangun** — tiket **07** |
| ⭐ rumus pajak brokerage | ⭐ **terbaca penuh** — tiket **07** |

## Area codebase

- Lapisan service: rantai perhitungan uang
- Lapisan service: efek keluar di akhir alur

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Langkah penetapan nilai | **268** terjangkau dari titik masuk; ⛔ isinya **tidak terekspor** |
| Aturan yang hilang | dipanggil `Activity\SumTSIPremiSpreadedRNM_Act.xml`; ⛔ **nol salinan** di 21 modul |
| Efek keluar | `Flow\InputRealizationTreatyIn.xml` — `Utility2`, lewat `Activity\serviceInsertArasapas_act.xml` **18.013 B, satu langkah**; ⛔ implementasinya ada di **dua modul lain** dan **tidak boleh dipinjam** |

## ADR terkait

- ⭐ **ADR-0003** — uang tidak `float`
- **ADR-0008** — efek keluar asinkron dengan antre-ulang
- **ADR-0015** — efek keluar wajib berhasil, transactional outbox

## Acceptance criteria

- [x] **AC 79** — ⭐ rantai perhitungan **dibangun dari isi langkah yang terbaca di ekspor**;
      ⛔ nol rumus yang ditebak *(AC 79 dicabut lalu diganti 2026-09-22 — lihat `spec.md`)*
- [x] **AC 87** — pengiriman ke layanan luar **tidak dibangun** sebelum muatannya diketahui

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| ⭐ ~~**P18**~~ | ~~isi 268 langkah~~ | ⭐ **DITARIK 2026-09-22** — tidak menahan |
| ⛔ **4** | di bagian mana aturan penjumlah yang hilang tinggal | ⛔ **MENAHAN** |
| ⛔ **5** | muatan efek keluar | ⛔ **MENAHAN** |

## Catatan

⛔⛔ **Menebak rumus di sini berarti membangun angka keuangan dari tebakan** — dan sekarang tidak
ada alasan menebak: rumusnya terbaca.

⚠️ Bab karantina di `spec.md` **§8.2** sudah **DIBUKA**, dan **keempat akibat** yang dulu
didalilkan bila P18 tidak dijawab **gugur seluruhnya**. Bab itu dipertahankan sebagai catatan.

⭐ **Yang dapat diambil sekarang:** seluruh rantai aritmetika uang. ⛔ **Yang masih ditahan:**
pemanggilan `SumTSIPremiSpreadRNMMultiCob_Act` *(P30)* dan pengiriman ke layanan luar *(P8, AC 87)*.

## ⭐ Penerapan KEPUTUSAN-RONDE-12 dan RALAT — 2026-10-03

- **Butir 6 (P30 dicabut — milik Fac In)** dan **butir 7 (P8 dicabut — ikuti yang terbaca)**: tiket
  ini tidak tertahan. Bunyi lama baris Blocker P30/P8 tetap dikutip di atas, tidak dihapus.
- **AC 87 — RALAT.** Bunyi lama: *"pengiriman ke layanan luar tidak dibangun sebelum muatannya
  diketahui"*. Digantikan butir 7: muatan **4 medan** (`CARI1` pzInsKey, `CARI2`/`CARI3` PolicyNo,
  `CARI21` ProdDateTime `dd/MM/yyyy HH:mm:ss`) dibangun di `services/konversi.go`, dijalankan sesudah
  commit, gagal **tidak** membatalkan simpan (`FlagErrorKonversi` verbatim ke layar), dilewati di luar
  produksi. Sambungan nyatanya (kunci `M_LINK_SERVICE`, persetujuan) `[terbuka]` — PERMINTAAN-TIM-INTI C1.
- **Butir 3/3b**: `BreakDownSpreading_Act` tidak dimigrasi (CountSpreading langkah 6).
