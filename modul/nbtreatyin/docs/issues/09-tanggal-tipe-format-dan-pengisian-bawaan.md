# 09: Tanggal — satu tipe, satu format, dan pengisian bawaan yang ditiru apa adanya

**Status:** sebagian *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ready-for-agent)*
**Blocked by:** —
**Menutup:** AC 32 · 33 · 34 · 35 · 69 *(5 AC)* — US 2 · 35 · 44

## Hasil & nilai pengguna

Hari ini tanggal disimpan **sebagai teks**, dan ⛔ **dua susunan berbeda dipakai di berkas yang
sama** — hari-bulan untuk tanggal mulai, bulan-hari untuk tanggal laporan. `[terverifikasi]`
⚠️ Sebagian nilai lama karenanya **ambigu secara mutlak**: satu nilai dapat berarti dua tanggal
berbeda, dan kekeliruannya **tidak menimbulkan pesan galat**.

Sesudah tiket ini, tanggal disimpan sebagai **tipe tanggal**, dengan **satu format** di seluruh
sistem, dan formatnya diurus di lapisan layar — ⭐ urutan dan perbandingan tanggal menjadi benar.

## Area codebase

- Lapisan repository: tipe tanggal
- Lapisan service: pengisian tanggal bawaan
- Lapisan handler: format tampilan

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Dua format berdampingan | `DataTransform\InputPolicyTreatyIn_preDT.xml` — dua susunan berbeda di berkas yang sama |
| Pengisian bila kosong | gerbang tanggal-mulai, tanggal-akhir, dan tanggal-laporan kosong ⇒ **tanggal hari ini** |
| Aturan satu tahun | `DataTransform\SystemSetOneYear_DT.xml` — menambah satu tahun ke tanggal mulai |

## ADR terkait

- **ADR-0009** — migrasi penuh

## Acceptance criteria

- [x] **AC 32** — tanggal disimpan sebagai **tipe tanggal**, ⛔ bukan teks
- [x] **AC 33** — **satu format** dipakai di seluruh sistem; ⛔ tidak ada dua format berdampingan
- [x] **AC 34** — tanggal akhir kosong diisi **tanggal hari ini**, ⛔ bukan ditambah satu tahun
- [x] **AC 35** — tanggal mulai dan tanggal laporan kosong diisi tanggal hari ini
- [ ] ⛔ **AC 69** — migrasi menghasilkan kontrak ber-tanggal-akhir **sama dengan tanggal mulai**
      untuk berkas yang tanggal akhirnya kosong

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| ⚠️ **15** | data lama bertanggal **ambigu mutlak** — satu nilai, dua arti | ⚠️ menahan **migrasi**, bukan perilaku baru |
| **12** | kapan aturan satu-tahun berjalan, dan apakah ia menimpa pengisian hari-ini | tidak menahan |

> ⛔ **RALAT 2026-10-03 (audit ulang P6) — aturan satu-tahun `SystemSetOneYear_DT`.**
> Bunyi lama (status putaran 1, `docs/alat/status.json`): *"tidak dibangun — [keputusan work owner]
> P35/AC 34: EndDate kosong diisi HARI INI, bukan +1 tahun"*; butir 12 di atas: *"kapan aturan
> satu-tahun berjalan, dan apakah ia menimpa pengisian hari-ini"*.
> ⭐ Bunyi baru, bukti XML: `DataTransform/SystemSetOneYear_DT` adalah `pyPreDataTransform` aksi
> refresh `change` sel `.StartDate` di `Section/DetailPolicyTreatyIn` (layar admin, dapat disunting);
> di `Section/DetailDeptHeadTreatyIn_UW` sel yang sama `pyReadOnly 1==1`. Langkah 1: `.EndDate =
> @DateTime.addCalendar(.StartDate,1,0,0,0,0,0,0)`. Jadi pengisian tanggal-akhir-kosong (P35/AC 34,
> `InputPolicyTreatyIn_preDT`) berjalan saat layar dibuka, sedangkan aturan satu-tahun berjalan setiap
> kali admin mengubah tanggal mulai dan **menimpa** `.EndDate`. P35 tidak melarangnya — jawabannya
> menulis peran rule ini `[terbuka]` — sehingga alasan "tidak dibangun" itu tidak sah (bab 4 PROMPT
> putaran 2). **Dibangun:** `models.SystemSetOneYear` (29 Februari → 28 Februari, semantik
> `Calendar.add`), aksi hitung `SystemSetOneYear`, `.StartDate` admin memicunya (`frontend/medan.ts`).
> AC 34 tetap ✅ (tanggal akhir **kosong** tetap diisi hari ini). Butir 12 terjawab (spec §9.2).

## Perintah verifikasi

1. Kosongkan tanggal akhir, simpan — ⭐ terisi **tanggal hari ini**, bukan setahun kemudian.
2. Simpan tanggal, baca kembali, urutkan — ⭐ urutannya **benar secara kronologis**.
3. Ubah format tampilan — ⭐ nilai tersimpan **tidak berubah**.

## Catatan

⚠️ `[penyimpangan sadar]` Tipe tanggal dan format tunggal **berbeda** dari Pega. ⭐ Tetapi
pengisian tanggal-akhir-kosong-jadi-hari-ini **ditiru apa adanya** *(P35)*, walau hasilnya kontrak
bermasa berlaku nol hari — ⛔ **disengaja, bukan cacat migrasi.**

## ⛔ Penyimpangan sadar — 2026-10-03

`InputPolicyTreatyInPre_Act` langkah 9 menanam `@substring(StatementDate,6,2)>25`, sedangkan
`GeneratePolicyNoTreaty_Act` langkah 5.3 membaca hari dari `POOLDATA.TANGGAL_CLOSING`. Preseden
`[keputusan work owner]` PremiumList Life (dua aturan sama: "ikuti yang dari DB") dan penjaga repo
`TestNolAmbangTutupBukuTertanam` ⇒ hari tutup buku dibaca dari tabel di kedua tempat. Langkah 3-4
(StatementDate = ProductionDate = sysdate) kotak When-nya TIDAK dicentang — berjalan setiap pra-proses.

## Tampilan tanggal — putaran 2, 2026-10-03 (paket P2)

`[penyimpangan sadar]` AC 33 (satu format): setiap sel tanggal hanya-baca layar realisasi, grid
angsuran (`.DueDate`), dan daftar Suggest (`.Date`) ditampilkan lewat `formatDate` inti (`DD-MM-YYYY`,
tanpa jam — `frontend/sajian.ts`, sajian `'tanggal'`). XML memakai beberapa format: `dd/MM/yyyy`
(`.StartDate`, `.EndDate`, `.StatementDate` admin, `.DueDate` atasan), `Date-Short-Custom-YYYY`
(`TreatyIn.Commencement`/`Termination`, `.ProductionDate` atasan), `dd/MM/yyyy hh:mm`
(`.StatementDate` atasan), `dd/MM/yyyy HH:mm` (`ListSuggest .Date`), `Date-Short`
(`ListSuggest .ProductionDate`) — justru "dua format berdampingan" yang AC 33 larang. Nilai
tersimpan tidak berubah.
