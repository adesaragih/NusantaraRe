# E03: Tracer — alur masuk endorsement sampai kasus berfase Policy

**What to build:** Seorang petugas memilih **nomor polis + jenis endorsement + tanggal**, dan sebuah
kasus endorsement tercipta — atau **ditolak dengan alasan spesifik**.

Ini irisan **tertipis yang lengkap** untuk alur masuk. Begitu hijau, premisnya terbukti: endorsement
lahir **langsung di fase polis**, tanpa fase penawaran maupun binding.

`[terverifikasi]` Kasus lahir dari kelas portal menuju kelas kerja berprefiks `EDM-`, dengan **tautan
tiga arah** antara kasus portal, kasus endorsement, dan handle assignment. Tidak ada satu pun
konektor di flow yang menetapkan fase penawaran atau binding.

⛔ **Penolakan adalah keluaran bisnis yang sah**, bukan kegagalan teknis — keduanya harus dapat
dibedakan.

**Asal (Pega).** `Activity\SetValueToEDMWork` langkah 7–13 · `DataTransform\DataToEDM` ·
`Flow\InputAddendumFacultativeIn` konektor `Start2 → Assignment7` · `Activity\SetEdmType`

**Keputusan.** K-049 (Seam 5) · K-029 · K-044

**Blocked by:** `..\08-registry-rules-eval.md` · E01
⛔ Tiket NB itu **belum dikerjakan** — lihat §0 indeks.

**Pembaruan 11-10-2026 (irisan I3, W-EDM-08…12):** `OpenCase` kini dipanggil dari layar lewat `POST
/api/endorsmentfacin/kasus` (`backend/services/kasus.go`, `repository/kasus.go`): periksa ulang → OpenCase →
satu transaksi `SEQ_WORK_POLIS_EDM` (260) + `T_WORK_POLIS` (LINI `FACEDM`, POSITION `ReasFacInMarketing`,
`Pending-Policy`, flag `1`) + `T_EDM_FACIN_KASUS` (262). Pengecualian langkah 5 dari tabel 261 `JENIS='BATAL'`.
Tautan `DataToEDM` tidak ditulis (A43). Keputusan agent A42–A49: `../LAPORAN-IMPLEMENTASI-EDM-11-10-2026.md` §5.

**Status:** sebagian — 01-10-2026, `backend/services/bukakasus.go` (`OpenCase`, `TautanAssignment`); test
`bukakasus_test.go`. Penghalang NB-08 gugur; E01 **tidak** diperlukan: `[terverifikasi]` langkah 1-13
`SetValueToEDMWork` dan `DataToEDM` tidak memanggil rule `When`. Daftar pengecualian langkah 5 (dua
literal nomor polis) menjadi konfigurasi `DaftarPolis`. Sisa: pembuatan ID kasus (`svcAddWorkObject`,
repository) dan cabang EDM retro langkah 21 (fac out).

- [x] **Seam 5 hidup** — `services.OpenCase`
- [x] Kasus lahir berkelas kerja endorsement, berprefiks **`EDM-`**
- [x] **Tautan tiga arah** — `EndorsementID` (EDM → portal), `PortalEDMHandle` (portal → EDM), `TautanAssignment` (langkah 22-25: worklist menang atas workbasket)
- [x] Nomor polis lama dan alasan endorsement terisi
- [x] Kasus **berfase Policy** (`IsCedingConfirm = Policy`); kasus uji negatif `TestOpenCaseTidakPernahFasePenawaran`
- [x] Ditempatkan di workbasket `ReasFacInMarketing` dengan tiket `AdminPolicy`
- [x] Jenis endorsement dari portal dibawa apa adanya (`EdmType`, `EdmTypeNew`, `Type`); tidak ada cabang yang dihapus
- [x] Penolakan **terbedakan** dari kegagalan sistem, dengan alasan spesifik per gerbang — `HasilBukaKasus.Ditolak` + `AlasanTolak` vs `error`; gerbang lain (E04) masing-masing berpesan sendiri
- [x] Tiap transisi menyebut rule Pega asalnya dalam komentar (§4.6)
- [ ] Pembuatan identitas kasus (`svcAddWorkObject`, penomoran `EDM-`) — milik repository; seam menerimanya sebagai masukan
- [ ] Cabang EDM retro langkah 21 (`EDMRetro_Act`) — jalur fac out, belum diport
