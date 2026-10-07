# 41: Sub-tab FEA (Fire Extinguisher Availability)

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah work owner 03-10-2026 (urutan
> sub-tab objek, lihat tiket 39). Belum ada tangkapan layar; tampilan diturunkan dari XML.

**What to build:** sub-tab **FEA** di baris objek FIRE. Isinya grid `.FEAList` (APAR · Sprinkler · Smoke Detector &
Alarm · Hydrant · Private Truck Brigade · Others Info) dengan Tambah / Hapus. Baris yang dibuka berisi form jumlah unit
dan tiga dropdown. Datanya ikut Save tab Object.

**Blocked by:** daftar Private Team Fire Brigade / Team & SOP Safety / Team & SOP Risk Management (aturan properti tidak
ada di korpus).

**Status:** frontend selesai 03-10-2026 (uji hijau); backend selesai 03-10-2026 (sesi c3; migrasi 190 ditulis, BELUM dijalankan).

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `FlowAction\InputFEA.xml`: pyLabel "Input Fire Extinguisher Availability" (FEA = Fire Extinguisher Availability);
  section `InputFEA`, kelas `Data-OfferFacIn-OfferFEAList`; tanpa pra/pasca-proses.
- `Section\FEAList.xml`: tombol Save (`SaveFacIn_Act`) di atas; grid `.FEAList` (master-detail, flow action
  `InputFEA`, 10 baris, tanpa total); kepala sel 15–20; Add = addRow, Delete = deleteRow (tanpa activity).
- Section `OfferFEAList!InputFEA` TIDAK ada di korpus. Padanannya `Section\InputFEA_IsUW.xml` (kelas sama, baca-saja)
  sel 12–21: APAR / Sprinkler / Smoke Detector & Alarm / Hydrant / Private Truck Brigade (Unit) (pxNumber), Private Team
  Fire Brigade / Team & SOP Safety / Team & SOP Risk Management (dropdown associated), Others Info (textarea).
- Data contoh: semua `FEAList` kosong; rancangan flat tidak punya tabel FEA.

## Keputusan agent

- **M-1** Form = medan `InputFEA_IsUW`, dibuat dapat diisi.
- **M-2** Jumlah unit = bilangan bulat ≥ 0; galat menahan Save tab Object.
- **M-3** Tombol Save di atas grid FEA tidak diulang (Save tab Object sudah menyimpan seluruh objek).
- **M-4** Tiga dropdown kosong sampai aturan properti datang.

## Kontrak

- `ObjekFire.fea: BarisFEA[]` (`apar`, `sprinkler`, `smokeDetector`, `hydrant`, `privateTruckBrigade`,
  `privateFireBrigade`, `teamSopSafety`, `teamSopRiskManagement`, `info`), ikut `GET`/`PUT …/objek`.

## Acceptance criteria

- [x] Grid + Tambah / Hapus + form; label diuji ke korpus.
- [x] Backend: tabel baru `T_FEALIST` (migrasi 190), simpan / baca.
- [ ] Daftar tiga dropdown.

## Backend (sesi c3, 03-10-2026) — disusun agent

- `GET`/`PUT …/objek`: tiap baris objek membawa `fea` (kunci persis `BarisFEA`; selalu larik). Disimpan ke tabel BARU
  `T_FEALIST` (induk `T_LOCATIONLIST` — `.FEAList` milik baris lokasi, bukan Property), urut `SEQ_NO`; dihapus-sisip
  ulang bersama objek (FEA dihapus tepat sebelum lokasi). Dibaca kueri ketiga per case, dipasangkan lewat
  `T_LOCATIONLIST.ID`.
- Validasi PUT (400 `baris[n].fea[m].<medan>`): lima jumlah unit kosong atau `^[0-9]+$` (M-2, pola lantai A111); lebar
  kolom; tiga dropdown tanpa enumerasi (M-4).
- Loader: tabel `T_FEALIST` + jalur `LocationList/FEAList` (induk T_LOCATIONLIST) + lipatan halaman `.DataFEA` ke baris
  FEA (medan sendiri, pola V-39) — skema 80 tabel / **1.433** kolom / 150 jalur; uji `TestFlattenFEA`.
- Migrasi **190** (`190_t_fealist.sql` + `_down`). ⛔ Ditulis, **tidak dijalankan** agent.

**Bukti `[terverifikasi]`:** `Section\FEAList.xml` = kelas `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance`; fixture
`nb-fire-1.json` `LocationList[0]` berkelas sama dan berkunci `FEAList`; nol jalur `FEA` di rancangan (`skema_gen.go`, dua
cara: pencarian `FEA` di nama tabel / medan — satu-satunya kena `RateLifeAverage`, bukan FEA — dan di daftar jalur = 0).

**Keputusan agent — DISETUJUI work owner 03-10-2026** (butir 82, diteruskan sesi `nusantarare-0f`: *"setuju keputusan
agent"*):

| # | Keputusan | Dasar |
| --- | --- | --- |
| A142 | Tabel baru `T_FEALIST`, induk `T_LOCATIONLIST`, kolom sistem pola tabel berulang | rancangan tidak punya tabel FEA; pola T_ADDITIONALSHIP (butir 70) |
| A143 | Lima jumlah unit `VARCHAR2(50)` teks; tiga kode dropdown `VARCHAR2(50)`; Others Info `VARCHAR2(500)` | angka teks pola A129 / NUMBER_OF_FLOOR; kode / catatan pola V-6 |
| A144 | Halaman tertanam `.DataFEA` dilipat ke baris FEA (bukan tabel anak) | empat medan tunggal; pola lipatan V-39 (PolicyData) |

**Aturan properti (03-10-2026, diteruskan sesi 0f):** `DDL\PrivateFireBrigade.xml`, `TeamSOPSafety.xml`,
`TeamSOPRiskManagement.xml` = Have / Not Have / No Info (nilai = label) dipakai frontend; ≤ 11 bita, muat di kolom
VARCHAR2(50). Backend tetap tanpa validasi enumerasi.
