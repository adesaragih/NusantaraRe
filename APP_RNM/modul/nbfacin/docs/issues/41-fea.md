# 41: Sub-tab FEA (Fire Extinguisher Availability)

> ⚠️ **Disusun agent atas perintah work owner — bukan hasil `/to-tickets`.** Perintah work owner 03-10-2026 (urutan
> sub-tab objek, lihat tiket 39). Belum ada tangkapan layar; tampilan diturunkan dari XML.

**What to build:** sub-tab **FEA** di baris objek FIRE. Isinya grid `.FEAList` (APAR · Sprinkler · Smoke Detector &
Alarm · Hydrant · Private Truck Brigade · Others Info) dengan Tambah / Hapus. Baris yang dibuka berisi form jumlah unit
dan tiga dropdown. Datanya ikut Save tab Object.

**Blocked by:** daftar Private Team Fire Brigade / Team & SOP Safety / Team & SOP Risk Management (aturan properti tidak
ada di korpus).

**Status:** frontend selesai 03-10-2026 (uji hijau); backend → sesi c3 (sesudah tiket 40).

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
- [ ] Backend: tabel baru (rancangan tidak punya), simpan / baca.
- [ ] Daftar tiga dropdown.
