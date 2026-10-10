# 03: Wewenang komite — ditegakkan di lapisan layanan

**Status:** ready-for-agent
**Blocked by:** 00 · `claim-facin\issues\08` *(wewenang klaim)*
**Menutup:** AC 23 · 24 · 25 · 26 · 27 · 28 · 29 *(7 AC)* — US 19–21

## Hasil & nilai pengguna

Hari ini ⛔ `[terverifikasi]` Hanya **satu** pemeriksaan pemilik giliran ada di korpus, dan ia menempel pada **tombol Submit di layar**. ⚠️ Di modul saudaranya **tidak ada sama sekali**. ⛔ Siapa pun yang dapat memanggil lapisan layanan dapat menyimpan keputusan untuk jenjang mana pun.

Sesudah tiket ini, ⭐ **Hanya akun beku pada jenjang berjalan** yang dapat menyimpan keputusan, dan penolakannya terjadi **di lapisan layanan** — ⛔ bukan di layar.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Sesudah tiket ini, ⭐ **Hanya akun beku pada
> jenjang berjalan** yang dapat menyimpan keputusan, dan penolakannya terjadi **di lapisan layanan** — ⛔ bukan di
> layar."* → **pemegang = anggota workbasket tingkat berjalan**, bukan akun beku. Layanan memakai `Pemegang(akun, peran)`
> (pola Komite Prop 09-10-2026, `modul/komiteclaimprop/backend/models/kasus.go`): tingkat 1 = anggota `ReasClaimSPVA`
> **atau** `ReasClaimSPVB` (KCF-01, cadangan SPV A seperti Pega); tingkat lain = anggota workbasket `KomiteID` baris
> `KomiteCount`; TT3 / TT4 = anggota `ReasClaimDeptHead` (KCF-03). Selainnya **403**. Submit juga ditolak bila kasus
> tertutup atau tingkatnya sudah memutus. **Tanpa larangan rangkap** (pola Komite Prop 09-10-2026). Penegakan di lapisan
> layanan tetap (K10, ADR-0014).

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Pemeriksaan pemilik giliran | satu-satunya di korpus; menempel pada **tombol Submit** |
| ⚠️ Penukaran identitas | ⛔ **dua akun ditukar menjadi akun ketiga SEBELUM pemeriksaan** — ⭐ **tidak dibawa** |
| Modul saudara | ⛔ **nol pemeriksaan** — tidak di aktivitas, tidak di layar |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"Pemeriksaan pemilik giliran | satu-satunya di
> korpus; menempel pada **tombol Submit**"* dan *"⚠️ Penukaran identitas | ⛔ **dua akun ditukar menjadi akun ketiga
> SEBELUM pemeriksaan** — ⭐ **tidak dibawa**"* → rinciannya dari XML: `SetProteksiSubmiteKomite` (dipanggil
> `SetValueKomite` S15, hanya TT2 karena S1 keluar bila `TransferType != 2`) L1 memetakan akun pelaku (pemetaan akun
> orang tertulis mati), lalu L2.1 menyetel `ProteksiKomite.CARI1 := 1` bila `pyPosition == "IT Developer"` **atau**
> (`KomiteAproval == 0` **dan** akun == `KomiteID`) pada baris tangga **mana saja**. Tombol Submit `ShowTransfer` nonaktif
> bila `ProteksiKomite.CARI1 == 0 || pyWorkPage.Adjustment.AcceptedNo != ''`. Arah prompt §5 #1: pemutus = anggota
> workbasket **tingkat berjalan**; **pintu belakang "IT Developer" dan pemetaan akun gugur** (prompt §3 "Gugur").

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K10** — ⭐ **ADR-0014 BERDIRI** — penegakan pemilik giliran **tetap di lapisan layanan**
- **K3 · K12** — ⭐ Pembandingan memakai **identitas akun**, ⛔ bukan alias

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"**K10** — ⭐ **ADR-0014 BERDIRI** — penegakan
> pemilik giliran **tetap di lapisan layanan**"* → tetap berlaku. ADR kini di
> `OUTPUT_HASIL_RNM/docs/bersama/adr/0014-pemutus-komite-ditegakkan-per-komiteid-tingkat-berjalan.md`. `KomiteID` tingkat
> berjalan kini berisi workbasket, jadi "pemilik" = anggota workbasket itu (ADR-0030, KCF-01). K3 · K12 tetap: akun
> pelaku dibandingkan apa adanya terhadap keanggotaan workbasket, tanpa penukaran.

## Yang harus diuji

- [ ] Hanya **akun beku pada jenjang berjalan** dapat menyimpan keputusan
- [ ] Percobaan oleh akun lain **ditolak dengan galat** — ⛔ bukan diabaikan diam-diam
- [ ] Percobaan yang ditolak **terekam** di jejak audit
- [ ] ⛔ Penukaran akun menjadi akun ketiga **tidak dibawa**
- [ ] ⭐ Layar boleh menyembunyikan tindakan tak berwenang — ⛔ **itu bukan penegakan**

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | ⛔⛔ **Isi daftar jabatan dan susunan jenjang BELUM ADA** — korpus tidak memuat satu pun rule otorisasi; 13 medan privilese seluruhnya kosong | ⛔⛔ **MENAHAN PEMBANGUNAN tiket ini.** ⭐ Tiketnya tetap ditulis dan jahitannya disebut; ⛔ ia **tidak dapat selesai** sebelum work owner memberikan daftarnya |

> ⛔ **RALAT 10-10-2026.** Sel lamanya dikutip utuh, tidak dihapus: *"⛔⛔ **MENAHAN PEMBANGUNAN tiket ini.** ⭐ Tiketnya
> tetap ditulis dan jahitannya disebut; ⛔ ia **tidak dapat selesai** sebelum work owner memberikan daftarnya"* → **butir
> 6 tertutup, tidak menahan.** Penutupnya ADR-0030 *"Aturan peran ditetapkan sekali dan berlaku lintas modul"*
> (`OUTPUT_HASIL_RNM/docs/bersama/adr/0030-aturan-peran-ditetapkan-sekali-lintas-modul.md`) + keanggotaan workbasket +
> KCF-01 (roster FACIN DEGREE 1–5 → `ReasClaimSPVA` / `ReasClaimDeptHead` / `ReasClaimTechDivHead` / `ReasClaimOpsDir` /
> `ReasClaimTechDir`; `ReasClaimSPVB` cadangan tingkat 1). Wewenang klaim (`claim-facin\issues\08`) sudah dibangun
> claimfacin tahap 1 dengan aturan yang sama.

## Seam & verifikasi

**Seam:** lapisan layanan komite — ⭐ **uji lewat sini, bukan lewat layar**.
2. Panggil sebagai **bukan pemegang giliran** ⇒ ⛔ **ditolak**.
3. Sembunyikan tombolnya di layar, lalu panggil layanan **langsung** ⇒ ⛔ **tetap ditolak**.
4. Periksa jejak audit ⇒ ⭐ percobaan yang ditolak **tercatat**.
5. Panggil sebagai akun yang dulu 'menyamar' ⇒ ⛔ **ditolak** — ⭐ penukaran tidak dibawa.

> ⛔ **RALAT 10-10-2026.** Butir dan langkah lamanya dikutip utuh, tidak dihapus: *"Hanya **akun beku pada jenjang
> berjalan** dapat menyimpan keputusan"*, *"Percobaan yang ditolak **terekam** di jejak audit"* dan *"4. Periksa jejak
> audit ⇒ ⭐ percobaan yang ditolak **tercatat**."* → butir pertama dibaca menurut RALAT di kepala tiket (anggota
> workbasket tingkat berjalan). **"Percobaan ditolak terekam" gugur** (prompt tahap 2 §3: di luar XML); langkah 4 tidak
> diuji, AC 25 gugur. Tambahan uji: anggota `ReasClaimSPVB` memutus tingkat 1 ⇒ diterima; anggota workbasket tingkat 2
> mencoba memutus saat giliran tingkat 1 ⇒ 403; akun berjabatan "IT Developer" yang bukan anggota ⇒ 403.
