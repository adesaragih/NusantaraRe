# 71: Menu View Policy (Facultative) — cari polis jadi, detail baca-saja (sesuai sistem baru)

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 09-10-2026 "jadi untuk cari case
> yang udh selesai darimana? dipega ada menu view policy, fungsinya untuk cari data polis yang udh jadi". XML diserahkan work
> owner ke `D:\migrasi\RNM\DDL\` (ViewOpenPolicyMarineInbox, GetListViewPolicy_act, GridOpenPolicyMarine,
> ViewEndorsementData_Act, BackViewPolicy_Dt, ResetFilterValue_DT). Riset: dua subagen sesi 0f, klaim inti diverifikasi ulang.

**Status:** kontrak backend dikirim ke sesi c3 09-10-2026; frontend sesi 0f.

## Bukti `[terverifikasi]`

- `DDL\ViewOpenPolicyMarineInbox.xml` (`DATA-PORTAL!VIEWOPENPOLICYMARINEINBOX`): mode daftar bertab Facultative (L1224) / Treaty
  (L17980) / Life (L35445); tab Facultative blok "List Polis": "View Option" `.FilterForViewRetro` ("Choose", wajib) + "Filter"
  `.FilterTermForPolicy` (placeholder "Policy No / ID"; Enter / tombol Filter → `GetListViewPolicy_act` SEARCHON "POLIFAC",
  L5691/L5969) + grid `GridOpenPolicyMarine` (L8594). Mode detail (container `IsShowOpenPolicyMarine`): Back, Refresh, "View
  Data" (versi; tampil bila >1), `Periode_IsUW` (L47714) + `ViewInwardFacultativeDtl_IsUW` baca-saja, View Facretro, riwayat
  akseptasi, Attachment File.
- `DDL\GetListViewPolicy_act.xml` blok 8 "POLIS FACIN": sama dengan `GetListViewFollowing_act` (tiket 70) KECUALI opsi 5 / 6 / 7
  cocok AWALAN: `INSUREDNAME LIKE 'x%'` (L2419), `SOB LIKE 'x%'` (L2581), `UPPER(QQNAME) LIKE 'x%'` (L2743); tanda kutip dibuang
  dari filter; tanpa cadangan json_following. SQL `BrowseDataViewPolis` (L2934) TIDAK ada di korpus — `BrowseDataViewPolis1`
  (FACINPRODUCTION) adalah salinan Save-As-nya `[dugaan kuat]`.
- `DDL\GridOpenPolicyMarine.xml`: `TempListPolicy.pxResults`, 20 per halaman, kolom ID · Policy Number · Reff. Number · Insured
  Name · Source of Business · QQ Name · Start Period · End Period · link "View Details" (sama dengan 9 kolom pertama
  GridViewFollowingNB); View Details → `SetDetailViewOpenPolicyMarine_PreDT` (TIDAK ada) → `GetPolicyData_ACT(prodke "0")` →
  `BrowsePolicyJSONMarine_Act`.
- "View Data": `GetProdKeOldData_SQL` = `select PRODKE … from pooldata.json_polis where NOPOLIS = … order by TGL_INPUT desc`;
  label = `substring(IDPEGA, 18)` (`GetPolicyData_ACT` L1262) = ID case versi itu; tampil bila versi > 1. Isi detail tiap versi
  di Pega = `JSON_POLIS.DATA_JSON` (adoptJSONObject) — JSON.
- Polis lahir di Pega TIDAK punya baris di tabel case sistem baru (`T_*`); datanya di `FACINPRODUCTION` / `CEDING_FACINPRODUCTION` /
  `JSON_POLIS`.

## Keputusan work owner 09-10-2026

- **K71-1** Detail = **campuran**: polis/versi yang case-nya ada di sistem baru → layar case baca-saja penuh (semua tab, riwayat,
  lampiran); polis lama Pega → ringkasan dari `FACINPRODUCTION` + `CEDING_FACINPRODUCTION`. Tanpa DATA_JSON (K70-1).
- **K71-2** "View Data": "tampilkan saja, kalau memang tidak ada endorsmentnya kan, brarti listnya cuma 1, tidak ada dropdown" —
  daftar versi dari kolom relasional `JSON_POLIS` (NOPOLIS / PRODKE / IDPEGA / TGL_INPUT; bukan DATA_JSON); dropdown hanya bila
  versi > 1.
- **K71-3** Akses: semua pengguna NB FacIn (baca saja).
- Lanjutan K70-3: View Option 3 No Slip (ARASAPAS) tetap ditunda.

## Keputusan agent (sesi 0f)

- **A71-1** Cakupan: tab Facultative "List Polis" saja; Treaty / Life (modul lain) dan "List Offer" tidak.
- **A71-2** Opsi 7 QQ Name: `UPPER(QQNAME) LIKE UPPER(filter)%` — Pega hanya meng-UPPER kolom (input huruf kecil tak pernah cocok);
  maksudnya tak peka huruf. Opsi 5 / 6 awalan apa adanya seperti Pega.
- **A71-3** Ringkasan polis lama: kepala (nomor polis, endorsement, ID case, tertanggung, QQ, SOB, ceding, periode, policy type,
  reff number, group, status bisnis, master polis) + baris FACINPRODUCTION (objek, lokasi risiko, item, coverage, akumulasi,
  mata uang, TSI 100 %, TSI RNM, rate, premi RNM, treaty, % share, RI commission) + daftar ceding; teks kolom sistem baru.
- **A71-4** Detail case sistem baru = halaman case yang sudah ada (case Resolved otomatis baca-saja), tombol Back kembali ke
  View Policy.
- **A71-5** View Facretro tidak dibangun (FacRetro = tiket 65).

## Backend sesi c3 (09-10-2026)

Berkas: `services/polis.go`, `repository/polis.go` (+ `CariPolisQQ` di `repository/following.go`), `models/polis.go`,
`handlers/polis.go`; layanan `DenganPolis(NewPolisOracle)`.

| Endpoint | Jawaban |
| --- | --- |
| `GET /api/nbfacin/polis?opsi=&filter=&halaman=` | `{baris, total, halaman, ukuran: 20}` (baris = BarisFollowing tiket 70) |
| `GET /api/nbfacin/polis/versi?nopolis=` | `{versi: [{prodke, caseId, sistemBaru}]}` |
| `GET /api/nbfacin/polis/ringkasan?nopolis=&prodke=` | `{kepala (14), baris (13), ceding [{cedingName, pctShare}]}` |

- **K71-B1 Cari** `[terverifikasi]` `DDL\GetListViewPolicy_act.xml`:
  - langkah 6 membuang `"` dan `'` dari filter, langkah 7 filter kosong -> keluar (lalu pangkas spasi);
  - blok 8 = klausa tiket 70 berurutan, KECUALI 8.7 `INSUREDNAME LIKE 'x%'`, 8.8 `SOB LIKE 'x%'` (awalan) dan 8.9 opsi 7
    `UPPER(QQNAME) LIKE UPPER(:1)` + `x%` (A71-2);
  - opsi 4 tetap `%x%`; opsi 3 / lain 400; bind;
  - SQL `BrowseDataViewPolis` tidak ada di korpus -> `BrowseDataViewPolis1` (FACINPRODUCTION) `[dugaan kuat]`;
  - tanpa penjaga case: hanya login (K71-3).
- **K71-B2 Versi** (GetProdKeOldData_SQL): kolom relasional JSON_POLIS NOPOLIS / PRODKE / IDPEGA, urut TGL_INPUT turun (lalu
  PRODKE turun, `[keputusan agent]`). caseId = SUBSTR(IDPEGA, 19) (= Java `substring(IDPEGA, 18)` GetPolicyData_ACT L1262).
  sistemBaru = baris T_WORK_POLIS ber-ID itu, LINI FAC. Polis tanpa baris JSON_POLIS -> satu versi "0" dari FACINPRODUCTION
  (MIN(IDPEGA), NOENDORS kosong `[keputusan agent]`). Tidak ada -> 404 "Polis tidak ditemukan.".
- **K71-B3 Ringkasan** (A71-3):
  - IDPEGA versi = JSON_POLIS (NOPOLIS, PRODKE; ganda -> TGL_INPUT terbaru); polis tanpa JSON_POLIS + prodke "0" ->
    FACINPRODUCTION MIN(IDPEGA);
  - baris = FACINPRODUCTION ber-IDPEGA + NOPOLIS itu, urut OBJECTNO, IDX_OBJECT_ITEM, COVERAGENOTE, JN_REAS `[keputusan agent]`;
  - kepala dari baris pertama; reffNumber = RISLIPCEDING terakhir menurut SPREAD_DATE se-NOPOLIS (rumus BrowseDataViewPolis1);
  - policyType = kode (kodePolicyType, label lama ikut);
  - RATE (kolom VARCHAR2; TO_NUMBER implisit dapat menulis ".15") dinormalkan ke desimal bertitik bila terurai, selain itu apa
    adanya;
  - angka lewat TO_CHAR TM9 -> apd -> teks; tanggal DD-MM-YYYY; NULL -> "";
  - ceding = CEDING_FACINPRODUCTION ber-IDPEGA itu, urut CEDINGID, CEDINGNAME;
  - nol DATA_JSON; versi / baris tidak ada -> 404; nopolis / prodke kosong -> 400.
- **K71-B4 Detail case sistem baru**: GET kasus / submit (riwayat) / lampiran tidak punya penjaga workbasket -> akun NB FacIn tanpa
  workbasket posisi membaca tanpa 403. GET submit memberi `bolehSubmit` false + alasan, bukan 403. Lampiran: hanya login.
- **K71-B5 Menu**: menu modul DATAR (satu baris `M_NAV_MENU` per modul; 962 sudah menyalakan `DIMIGRASI`). Halaman
  `nbfacin-view-policy` cukup di `HALAMAN_NBFACIN` frontend; tanpa migrasi, tanpa baris menu. Penjaga backend
  (`inti/backend/penjaga`) TIDAK membaca daftar halaman frontend.
- MODUL.md baris `JSON_POLIS` (tabel warisan) diperbarui: ditulis 69a-2, dibaca kolom relasional, DATA_JSON tidak dibaca.

**Uji**:
- services: `TestJenisCariPolis` (10 kasus + kutip), `TestCariPolis` (akun tanpa peran), `TestVersiDanRingkasanPolis`;
- repository: `TestSQLPolis` (27 kolom dihitung dari SQL = 27 indeks Scan, bind, nol DATA_JSON);
- handlers: `TestPolisHandler` (VersiPolis / RingkasanPolis + blok `kepala` 14 / `baris` 13 = api.ts dua arah + urutan; 400 / 401 /
  404).

**Mutasi 11 / 11 mati.** ⚠️ `PolisOracle` (urutan baca: JSON_POLIS -> cadangan FACINPRODUCTION) hanya teruji lewat teks SQL.

## Menu (09-10-2026)

- **K71-4** "View Policy dibuatin di menu baru. dibawah menu home" (gambar layar sidebar): tombol sidebar **View Policy** tepat
  di bawah Home, di luar golongan. Dibangun generik di inti: `MenuModul.menuAtas` (aditif, pola `antreanBeranda`) →
  `frontend/daftar.ts` `ENTRI_MENU` (`atas: true`) → `susunMenu` (`MenuTersusun.atas`; tampil bila tombol modul pemiliknya
  tampil; halamannya tidak menyalakan tombol NB FacIn) → `Shell` (ikon `IkonDokumen`) + palet Ctrl+K (sesudah Home).
  Pengecualian sadar atas "satu modul satu menu" (30-09-2026). Tanpa baris `M_NAV_MENU` / migrasi. Tombol sementara di portal
  NB FacIn dicabut.
- Uji inti diperbarui: `frontend/daftar.sinkron.test.ts` (palet = Home + tombol atas + tombol modul), `frontend/Shell.test.ts`
  (tombol atas bukan tombol modul). Kegagalan lama di luar tiket ini tidak berubah (premiumlistlife, treatycontractout,
  nbtreatyin, masterproductnamelife, Shell "menu profil").
