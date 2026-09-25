# F09: Dua generator nomor retrosesi — tidak disatukan

**What to build:** Penawaran retrosesi mendapat nomor slip yang unik dan berurut, memakai **dua skema
penomoran berbeda** yang hidup berdampingan di sistem lama.

⛔ **DUA generator, bukan satu.** `[terverifikasi]` Skemanya berbeda, urutannya berbeda, dan
**sequence-nya berbeda**:

| | Generator 1 | Generator 2 |
| --- | --- | --- |
| Berkas | `DDL\GENERATE_FACRETRO_NO.txt` (fungsi Oracle) | `RDBList\GenerateOurRefFacOut_SQL` (`Rule-Connect-SQL`, tag `pyBrowseSQL`) |
| Bentuk | `RNM-` + FacCode + BusinessCode + `.` + bulan + `.` + tahun + `.` + `LPAD(seq,5,'0')` | `'Y'` + `{InputData.CARI21}` + `{…BusinessOldId}` + `'RNM'` + `to_char(sysdate,'yy')` + seq |
| Sequence | `FACRETRO_SEQ` | **`T_FAC_OFFER_SEQ`** |
| Padding | `LPAD(…,5,'0')` | **tanpa padding** |
| Posisi `'Y'` | di tengah (sebagai FacCode) | **di depan** |

`[terverifikasi]` Pemetaan FacCode: `'Y'` bila `FacType = 'FAKULTATIF RETROSESI'`; `'F'` bila
`'FAKULTATIF INWARD'`. Jadi retrosesi berawalan **`RNM-Y…`**, berbeda dari Fac In `RNM-F…`.

⛔ **Tidak ada cabang `ELSE`.** `[terverifikasi]` Bila `FacType` bukan salah satu dari kedua nilai itu,
`FacCode` tetap **NULL**. Perilaku itu **diport apa adanya** — jangan ditambahi cabang default.

⚠️ **`FacRetroSlipNumber` dideklarasikan `VARCHAR2(21)`.** `[terverifikasi]` Panjang keluarannya
bergantung panjang `BusinessCode`; bila melebihi 21 karakter, fungsi Oracle gagal saat runtime.
**Diport apa adanya**, dengan batas panjangnya diuji. `[pertanyaan terbuka]` apakah `BusinessCode`
dijamin pendek — `GLOSARIUM.md` mencatat 98 kode dengan arti **belum terverifikasi**.

📌 `[pertanyaan terbuka]` Korpus **tidak menjelaskan kapan** masing-masing generator dipakai.
Keduanya diport; pemilihannya mengikuti pemanggil, bukan ditebak.

**Asal (Pega).** `DDL\GENERATE_FACRETRO_NO.txt` · `RDBList\GenerateOurRefFacOut_SQL`

**Keputusan.** K-055 · K-056 · `CLAUDE.md` §1, §4.3, §4.6

**Blocked by:** F08

**Status:** blocked

- [ ] **Kedua generator** hidup terpisah; tidak diseragamkan menjadi satu
- [ ] Bentuk nomor retro `RNM-Y…` berbeda dari Fac In `RNM-F…`
- [ ] ⛔ **Tanpa cabang `ELSE`** — `FacType` tak dikenal menghasilkan kode kosong, **bukan** nilai default
- [ ] Padding 5 digit hanya pada generator 1; generator 2 **tanpa padding**
- [ ] Nomor urut diambil dari sequence basis data, **bukan** dihitung aplikasi
- [ ] Batas panjang `VARCHAR2(21)` diuji sebagai **batas nyata**, bukan diabaikan
- [ ] **K-046** `K046_GenerateFacRetroNo_TanpaElse`
- [ ] Nomor polis produksi **tidak pernah** muncul di test; memakai nilai sintetis (K-025)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
