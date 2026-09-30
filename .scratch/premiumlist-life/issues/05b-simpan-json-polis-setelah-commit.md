# 05b: ~~Simpan JSON polis dan rekam produksi~~ — ⛔ **DIBATALKAN 2026-09-16**

**Status:** wontfix — `JSON_POLIS`/`JSON_OFFER_LIFE` dibuang (spec §12); sisa non-JSON (`Utility1`, `finishAssignment`) dibangun di `2b8cadc`

**Blocked by:** —

> ⛔ **JANGAN KERJAKAN TIKET INI.** `[keputusan work owner]` **Seluruh JSON dibuang** (spec §12).
> Tiket ini seluruhnya tentang menulis `JSON_POLIS` dan `JSON_OFFER_LIFE` **sesudah** commit, dan
> tentang menangani titik potong yang ditimbulkan keduanya. Kedua procedure itu **tidak lagi
> dipanggil**, sehingga **seluruh alasan keberadaan tiket ini hilang**.
>
> **Ke mana isinya pindah:**
>
> | Yang dulu di sini | Sekarang |
> | --- | --- |
> | Tulis `JSON_POLIS` / `JSON_OFFER_LIFE` | **dibuang** — tidak ada padanannya (AC 32, 34 spec) |
> | Urutan "sesudah commit" & titik potong | **lenyap** — satu polis = **satu transaksi** (tiket 05a, AC 35 spec) |
> | Penulisan peserta ke `M_LIFE_PREMIUM_DETAIL` | ke **`T_PREMIUM_LIST_DETAIL`**, di dalam transaksi yang sama (tiket 03 + 05a) |
> | `SaveLifeinProduction_SQL` sebagai titik potong | `LIFEINPRODUCTION` **tidak ditulis** lagi; hanya dibaca saat migrasi (tiket 00, AC 33 spec) |
> | Deteksi keadaan separuh & pengulangan | **tidak perlu** — atomisitas menggantikannya (AC 35 spec) |
> | Efek keluar Arasapas | tetap di tiket **06** |
>
> Berkas dipertahankan sebagai jejak keputusan, **bukan** sebagai pekerjaan.

---

<details>
<summary>Isi asli (sudah tidak berlaku)</summary>

**Blocked by (asli):** 05a (transaksi summary harus commit lebih dulu — itulah inti tiket ini)

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin representasi JSON polis dan rekam produksi ditulis **setelah**
premium list aman tersimpan, dan saya ingin penulisan itu **dapat diulang tanpa menggandakan data**
bila gagal di tengah — supaya kegagalan jaringan atau basis data tidak pernah meninggalkan premium
list yang benar dengan polis yang tidak pernah terbit. *(User story 33–36 di spec)*

## Area codebase

`internal/repository` (pemanggilan procedure yang commit sendiri, di luar transaksi Go; deteksi
keadaan separuh), `internal/services` (perakitan payload JSON polis; orkestrasi urutan + pengulangan),
`internal/handlers` (status "tersimpan tetapi polis tertunda" terbaca API).

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InsertJsonPolisLife_Act` | `ASM-FW-GISFW-WORK-LIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml` (327.332 byte, tersimpan `20260728T024422`, **17 langkah**) | orkestrator |
| `InsertJsonPolis` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLIS` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/InsertJsonPolis.xml` | `POOLDATA.INSERTJSONPOLISLIFE(p_IDPEGA, p_NOPOLIS, P_JSONDATA, P_TGL_PROD, P_USERNAME, {OutputData.HASIL1 out})` + `COMMIT;` (102) |
| `SaveLifeinProduction_SQL` | `ASM-FW-GISFW-WORK-LIFE` / `ASM!SAVELIFEINPRODUCTION_SQL` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/SaveLifeinProduction_SQL.xml` | `INSERT INTO POOLDATA.LIFEINPRODUCTION` + `COMMIT;` (164) |
| `GetNopolisByIDPega` | `ASM-FW-GISFW-…` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetNopolisByIDPega.xml` | baca balik untuk verifikasi |
| `GetJsonProductLife` | `ASM-FW-GISFW-…` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetJsonProductLife.xml` | data produk untuk payload |

`[terverifikasi]` Rantai rujukan `<RequestType>` di `InsertJsonPolisLife_Act`: `GetJsonProductLife`
(1444) → `InsertPLSummary` (3112, step **8**) → `InsertJsonPolis` (4094, step **10**) →
`SaveLifeinProduction_SQL` (4271, step **11**) → `GetNopolisByIDPega` (4449, step **12**).
**`SaveOfferJsonLife_SQL` tidak ada di rantai ini** — lihat tiket 01 dan **OQ-067**.

`[data DBA]` `POOLDATA.INSERTJSONPOLISLIFE` — **upsert** ke `POOLDATA.JSON_POLIS` berkunci `IDPEGA`
(`UPDATE` bila ada, `INSERT` bila belum). Kolom: `IDPEGA`, `DATA_JSON` (**CLOB**), `TGL_INPUT`,
`NOPOLIS`, `PRODKE`, `TGL_PROD`, `USERNAME`. Gagal → jatuh ke `JSON_POLIS_ERROR`.
**COMMIT di dalam procedure.**

`[terverifikasi]` Step **12** "Cek sudah masuk atau blm datanya" membaca balik lewat
`GetNopolisByIDPega`; step **13** menyalakan penanda bila `OutDataLife.pxResults(1).PL_NUMBER==""`
(baris 5105). Inilah **deteksi keadaan separuh** milik Pega — ia dipertahankan sebagai konsep, dan
menjadi pemicu alarm di tiket **06**.

## ⚠️ Urutan yang mengikat `[keputusan desain]`

1. Transaksi Go (penomoran + summary) **sudah commit** — tiket **05a**.
2. **Baru** `INSERTJSONPOLISLIFE` dipanggil, **di luar** transaksi Go, karena ia commit sendiri.
3. `SaveLifeinProduction_SQL` juga commit sendiri (baris 164) — ia **titik potong ketiga**.
4. **Penulisan detail peserta NB ke `M_LIFE_PREMIUM_DETAIL` termasuk di sini** — dipanggil **INLINE**
   sebagai bagian alur simpan polis, lewat logika `InsertLifePremiumDetail_act`
   (`ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT` / `RULE-OBJ-ACTIVITY`) →
   `SaveMasterLPDet` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET`), yang
   **commit sendiri** (`COMMIT;` baris 252) — **titik potong keempat**. Rincian dan AC-nya di tiket
   **08**. ⚠️ `[keputusan work owner]` **Tanpa job** — di Pega, NB dipicu batch; di sistem baru
   inline, meniru EDM.

**Jangan** menempatkan data yang harus atomik pada dua sisi procedure yang commit sendiri.

## ADR terkait

**ADR-0015** (batas transaksi dipegang Go; efek yang tidak boleh hilang ditangani eksplisit),
**ADR-0003** (uang non-float di dalam payload JSON), **ADR-0013** (endpoint di-resolve runtime —
berlaku pada efek keluar di tiket 06), **ADR-0011**.

## Acceptance criteria

- [ ] `INSERTJSONPOLISLIFE` dipanggil **setelah** transaksi summary commit, tidak pernah di dalamnya.
      *(AC 21 spec)* — belum: wontfix — `INSERTJSONPOLISLIFE` dibuang (spec §12)
- [ ] Kegagalan pada `INSERTJSONPOLISLIFE` meninggalkan nomor + rekam summary **utuh**; API tetap
      melaporkan premium list tersimpan, dengan penanda bahwa polis **tertunda**. *(AC 23 spec)* — belum: wontfix — tidak ada procedure JSON; simpan satu transaksi
- [ ] Pemanggilan ulang `INSERTJSONPOLISLIFE` untuk `IDPEGA` yang sama **tidak** menggandakan baris —
      dibuktikan dengan memanggilnya dua kali dan menghitung baris `JSON_POLIS`. — belum: wontfix — `JSON_POLIS` tidak ditulis
- [ ] Keadaan separuh **terdeteksi**: setelah penulisan, sistem membaca balik dan menandai kasus yang
      belum lengkap. Penanda itu terbaca lewat API dan menjadi masukan tiket 06. — belum: wontfix — digantikan atomisitas (langkah 12–13 tidak ditiru)
- [ ] `SaveLifeinProduction_SQL` diperlakukan sebagai **titik potong** tersendiri; kegagalan
      sesudahnya tidak merusak apa yang sudah ter-commit dan dapat diulang dengan aman. — belum: wontfix — `LIFEINPRODUCTION` tidak ditulis
- [ ] Uang di dalam payload JSON ditulis sebagai desimal presisi arbitrer — **tidak** lewat `float`
      dan tidak lewat pembulatan diam. *(AC 14 spec; **ADR-0003**)* — belum: wontfix — tidak ada payload JSON
- [ ] Ada test yang **gagal** bila urutan 05a → 05b dibalik. *(AC 24 spec)* — belum: wontfix — urutan 05a → 05b lenyap; urutan simpan-lalu-tutup dikunci `TestSimpanSebelumTutupDalamSatuTransaksi`
- [ ] Kegagalan yang jatuh ke `JSON_POLIS_ERROR` **terlihat**: sistem tidak menganggapnya sukses. — belum: wontfix — `JSON_POLIS_ERROR` tidak dipakai
- [x] Penulisan detail peserta NB ke `M_LIFE_PREMIUM_DETAIL` terjadi **di dalam alur simpan ini**,
      bukan dijadwalkan. Setelah respons simpan berhasil, baris untuk `PL_NUMBER` itu **sudah ada**.
      Tidak ada job/cron/worker terjadwal di jalur ini. *(AC lengkap di tiket 08;
      `[keputusan work owner]`)* — bukti: `services/polis_summary.go:simpanDalam` (`warisan.Ganti` di transaksi simpan, nol penjadwal); uji `TestSimpanDalamUrutanTerkunci`

## Blocker

**Tidak ada pemblokir.** Catatan terbuka yang tidak memblokir: **OQ-067** (penempatan
`INSERTJSONOFFERLIFE` — ditulis di tahap penawaran, tiket 01).

## Catatan — kode mati yang tidak dimigrasikan

`[terverifikasi]` Di `InsertJsonPolisLife_Act`, hanya **dua** `<pyStepsBlockName>` berisi `//`:

| Bagian | Baris | Nasib |
| --- | ---: | --- |
| step 16 `Commit` eksplisit | 5294 | **mati** — Go memegang transaksi (**ADR-0015**) |
| step 17 `Connect-REST` `InsertLifePremiumDetail` *(dulu tertulis `ConvertJsonNusareToProduction` — diralat 28-09-2026, b5422)* | 5383 | **mati** — tidak dipakai |

`[keputusan work owner]` **Juga tidak dimigrasikan** meski di korpus **AKTIF**: empat gerbang treaty
ID pada step 6–7 — `@contains(.ID,"1000032")` **QS** (1714), `"1000033"` **2nd QS** (1856),
`"1000034"` **SURPLUS** (1998), `"1000035"` **2nd SURPLUS** (2140). Itu logika polis lama. Jenis
treaty diambil dari data **`TypeCeding`** (`1`=QS, `2`=SURPLUS, `3`=QS+SURPLUS, `4`=XOL — ekspresi
baris ~3787). *(AC 30–31 spec)*

`[keputusan work owner]` Step **4** dengan precondition `@toDecimal(Local.currentdate)>25` (baris
1170) **tidak** direplikasi — ambang dibaca dari tabel, lihat tiket **02**.

⚠️ **OQ-066**: penanda `<pyStepsBlockName>` **tidak dapat dipercaya sendirian** di modul ini —
berkas ini justru buktinya. Konfirmasikan ke work owner sebelum menyimpulkan langkah lain.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

</details>

---

## ⚠️ RALAT 28-09-2026 (giliran 10) — tiket ini TIDAK seluruhnya kosong: `Utility1` dan `finishAssignment`

Pembatalan 2026-09-16 benar untuk **JSON**: `JSON_POLIS`/`JSON_OFFER_LIFE` tetap dibuang (pl1, spec §6,
§12). Tetapi kalimat *"seluruh alasan keberadaan tiket ini hilang"* **dibantah korpus** di dua titik,
dan keduanya perilaku yang **tidak dikerjakan tiket mana pun** sampai giliran ini:

1. **`Utility1`.** `InputPolicyHolder.xml`: `Decision2 --Transition7 [Confirm b2090]--> Utility1
   (InsertJsonPolisLife b765) --Transition8--> END52`. Tiket 01 menutup kasus Resolved-Completed pada
   `Confirm` di tahap detail **tanpa** menjalankan apa pun di antaranya — padahal activity itu yang
   memanggil `InsertPLSummary` (langkah 8). Akibatnya: kasus selesai tanpa rekap tersimpan.
2. **`finishAssignment`.** `Section/ShowLifePremiumSummary.xml`: tombol `Submit` b27471 →
   `InsertJsonPolisLife_Act` b26414 → `finishAssignment` b26442; `Assignment1 --Transition2
   [ShowLifePremiumSummary b1881]--> END52` (Resolved-Completed b947).

**Status:** `wontfix` → **selesai sebagian 28-09-2026** — sisa non-JSON dibangun, JSON tetap dibuang.

### Tujuh belas langkah `InsertJsonPolisLife_Act`, diputuskan satu per satu

Perintah audit (nomor baris = `sed -e 's/></>\n</g'`):

```
grep -n "<pyStepsActivityName>\|<pyStepsDescription>[^<]\|<RequestType>\|<pyStepsBlockName>[^<]" \
  Activity/InsertJsonPolisLife_Act.xml
```

| # | Metode · deskripsi | Putusan |
| ---: | --- | --- |
| 1 | `Obj-Refresh-And-Lock` | ⚙️ padanan: transaksi Go + `FOR UPDATE` penghitung nomor |
| 2 | `Property-Set` | ➖ penetapan halaman kerja; tanpa kolom |
| 3 | `set prodatetime` | ➖ milik `JSON_POLIS.TGL_PROD` — dibuang bersama JSON |
| 4 | `set prodatetime > 25` | ⛔ tidak ditiru — `>25` tertanam (pl5, tiket 02) |
| 5 | `Pega to jsondata` | ⛔ dibuang (JSON) |
| 6 | `RDB-List GetJsonProductLife` (b1444) | ⛔ bahan payload JSON — dibuang |
| 7 | QS / 2nd QS / SURPLUS / 2nd SURPLUS (`@contains(.ID,"100003x")`) | ⛔ tidak dimigrasikan `[keputusan work owner]` (catatan asli di bawah) |
| 8 | `Insert to table summary` → `InsertPLSummary` (b3112) | ✅ **dibangun** — `T_PREMIUM_LIST_SUMMARY` (tiket 05a); `M_LIFE_PREMIUM_SUMMARY` ✅ sejak GILIRAN-18 (OQ-PL-09 ditutup, `SummaryWarisan.Ganti`) |
| 9 | `Pega to json_offer & lifeinproduction` | ⛔ dibuang (JSON) |
| 10 | `InsertJsonPolis` (b4094) | ⛔ dibuang (JSON) |
| 11 | `SaveLifeinProduction_SQL` (b4271) | ⛔ `LIFEINPRODUCTION` tidak ditulis (AC 33 spec) |
| 12 | `Cek sudah masuk atau blm datanya` → `GetNopolisByIDPega` | ➖ digantikan atomisitas: satu commit, tidak ada keadaan separuh untuk dibaca balik |
| 13 | penanda bila `PL_NUMBER==""` | ➖ idem |
| 14 | `SendEmailNotification` — *"Kalau blm, email errornya"* | → tiket **06** (outbox) |
| 15 | `serviceInsertArasapasLife_act` | → tiket **06** (outbox + stub) |
| 16 | `Commit` | ➖ `//` mati (b5294); Go memegang transaksi |
| 17 | `Connect-REST InsertLifePremiumDetail` *(diralat 28-09-2026)* | ➖ `//` mati (b5383) |

Salinan peserta ke `M_LIFE_PREMIUM_DETAIL` (`InsertLifePremiumDetail_act` → `SaveMasterLPDet`, pl2)
ikut di jalur yang sama.

### Yang dibangun

- `models.AkibatKeputusan.SimpanPolis` — **hanya** `Transition7` (detail + Confirm → `Utility1`) dan
  `models.PenyelesaianSummary()` (`Transition2`). ⛔ `Offer` (`Transition11`) juga berakhir di `END52`
  tetapi **tanpa** `Utility1` — ia tidak menyimpan (`TestHanyaUtility1YangMenyimpanPolis`).
- `Penawaran.terapkan`: bila `SimpanPolis`, `SummaryPremiumList.simpanDalam` (nomor → rekap → ganti
  rekap → sumber → ganti warisan) berjalan **di transaksi yang sama, sebelum** `TutupKasus` dan jejak
  (`TestSimpanSebelumTutupDalamSatuTransaksi`). Kasus tertutup tanpa rekap tidak mungkin lagi.
- `POST /api/polis-life/{id}/summary` kini = `Submit` + `finishAssignment`: **hanya** dari tahap
  `Input Premium Summary` (409 selain itu, `ErrSubmitBukanTahapSummary`), menutup Resolved-Completed
  dengan jejak `(Submit)`. Layar mematikan tombolnya sesudah berhasil dan mengatakan kasusnya tertutup.

### ⚠️ Akibat yang harus diketahui work owner

`Confirm` di tahap detail kini **menuntut** polis siap disimpan: punya peserta, `Type` QR/QP/TP/TR,
bahan nomor lengkap. Polis tanpa peserta yang dulu dapat ditutup dengan `Confirm` kini dijawab **409**
dengan pesan "unggah rincian peserta lebih dahulu". Itu perilaku Pega yang sebenarnya (activity itu
berjalan di antara keputusan dan penutupan), tetapi **berbeda** dari perilaku tiket 01 yang sudah
ter-commit.

### Angka

Go **534 PASS · 0 FAIL** tingkat atas; vet (+`-tags db`), gofmt bersih · vitest **347** · tsc bersih.

### Temuan `/code-review` 28-09-2026 atas `d7610fb..4559ff3` (tiket 05b–09) — diperbaiki atau dicatat

Dua sumbu (Standards, Spec) paralel. Spec memverifikasi ulang dari korpus: `Transition7 → Utility1 →
END52` dan `Transition11 → END52` tanpa `Utility1` (✅; `Transition11` keluar dari `Decision3`, bukan
`Decision2`), kunci Arasapas b844–845 dan tiga pengenal permintaan REST (✅), 17 langkah
`InsertJsonPolisLife_Act` (✅), gerbang tahap summary (✅ dengan catatan di bawah).

| Sumbu | Temuan | Tindakan |
| --- | --- | --- |
| Standards 1 | ⛔ `sqlPungutEfek` tanpa penyaring `MODUL`: worker Claim Life akan memungut baris `PREMIUMLISTLIFE`, pelaksananya tidak mengenal `arasapas-polis`, lalu menandainya **gagal permanen** — dan komentar `AntreanEfekOracleModul` mengklaim "worker yang sama memungut keduanya" tanpa bukti | **diperbaiki**: `AND MODUL = :3`, `PekerjaEfek.modul = CLAIMLIFE`; klaim komentar diganti yang benar (baris polis menunggu pekerja PremiumList yang belum ada). Dikunci `TestPungutEfekHanyaModulnyaSendiri` |
| Standards 4 | ⛔ TOCTOU: tahap dibaca di luar transaksi, `TutupKasus` hanya menolak kasus tertutup — kasus yang berpindah tahap tetap ditutup lewat konektor tahap lamanya | **diperbaiki**: `sqlTutupPolis` kini menuntut `STATUS = tahap yang dibaca` (pola `sqlPindahTahapPolis`). Dikunci `TestTutupPolisDikunciTahapYangDibaca`. Berlaku juga untuk `Confirm`/`Decline` |
| Standards 2 | jejak kegagalan polis masuk `T_CLAIMLF_JEJAK` dengan `KLAIM_ID` = id polis, tak terbedakan dari jejak klaim | **dicatat** — pola sudah ada sejak tiket 01 PremiumList (`PerekamJejakOracle` untuk keputusan polis); tabel jejak polis adalah keputusan skema `[terbuka — work owner]` |
| Standards 3 | `Putuskan` membuang `EfekKeluar` pada jalur `Confirm`/`Utility1` | **dicatat** — layar penawaran menerima `jawabanAkibat` saja; ringkasan efek tampil di jalur summary |
| Standards (smell) | `Penawaran.terapkan` ↔ `SummaryPremiumList` saling memanggil metode privat; setter `DenganJejak`/`DenganPenyalur` kembar | **dicatat** — satu tempat penutupan kasus sengaja dipilih; disatukan bila keduanya disentuh lagi |
| Spec (c) | ⚠️ `POST /summary` kini hanya menerima tahap `Input Premium Summary`, dan **tidak satu pun konektor menuju `Assignment1`** — jadi di keadaan sekarang **tombol Submit summary selalu 409** untuk kasus baru; simpan berjalan lewat `Confirm` di tahap detail (`Utility1`) | **dinyatakan terang di sini**: itu bentuk korpus (tiket 01: `Assignment1` nol konektor masuk, `[terbuka — work owner]`). Layar summary dan rutenya menunggu jawaban itu |
| Spec (c) | muatan outbox membawa `Waktu` = saat keputusan, sedangkan kontrak REST memetakan `.pxCreateDateTime` (saat kasus lahir) | **dicatat di OQ-PL-11** — muatan panggilan nyata ditetapkan bersama jawabannya |
| Spec (a) | log per panggilan tanpa nomor polis; log keberhasilan belum ada | **dicatat** (tiket 06, AC "sebagian") |
