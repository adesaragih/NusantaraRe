# 04: Layar komite — 92 medan, dua kotak centang terkunci

**Status:** ready-for-agent
**Blocked by:** 00 · 02
**Menutup:** AC 30 · 31 · 32 · 33 · 34 · 35 *(6 AC)* — US 6–11

## Hasil & nilai pengguna

Hari ini Anggota komite **belum punya layar** untuk menilai penyesuaian, dan aturan siapa boleh menyunting apa belum ditetapkan.

Sesudah tiket ini, Anggota komite melihat rekapitulasi penyesuaian **dalam mata uang asli dan dalam rupiah**, dapat menuliskan catatan, dan ⭐ **dua kotak centang usulan hanya dapat disunting penyetuju pertama**.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Hari ini Anggota komite **belum punya layar**
> untuk menilai penyesuaian, dan aturan siapa boleh menyunting apa belum ditetapkan."* → cara layar dicapai kini
> ditetapkan pola Komite Prop / Non Prop sesudah commit `dcd1522e` (perintah work owner 09-10-2026: menu komite dihapus,
> "anggap menu itu tidak pernah ada"):
>
> - Modul **tanpa menu**: `frontend/layar.ts` (`PENDAFTARAN_LAYAR`), bukan `menu.ts`; slot menu **`—`** (penjaga
>   `modulTanpaMenu`). Baris menu `komiteclaimfacin` sudah dibuang migrasi inti 949.
> - **Tabel komite di bawah inbox Claim Fac In**: hanya kasus `KMT-` yang menunggu workbasket / akun pelaku, tanpa
>   switch, tidak ikut tab; tampil bila `GET /api/claim-fac-in/hak` → `komite`.
> - Klik baris membuka layar komite **di tempat** (`onBukaModul`, `MODUL_DIPINJAM` di `frontend/App.tsx`, `ruteDipinjam`
>   di `cmd/api/rakit.go`); Back / Submit kembali ke inbox Claim Fac In. Prefix rute `/api/komite-claim-fac-in`.
> - **View more details** = klaim hanya-baca (`bukaKasus.hanyaLihat`, `GET …/kasus/{id}?lihat=1`).

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Layar komite | 92 medan; ⭐ **tepat dua** terkunci bagi jenjang selain yang pertama |
| Dua kotak centang | usul **menutup klaim** dan usul **mencadangkan** |
| ⚠️ Jebakan alih | ⛔ pembanding di Pega bertipe **teks**, padahal pencacahnya **bilangan** |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"Layar komite | 92 medan; ⭐ **tepat dua**
> terkunci bagi jenjang selain yang pertama"* → isinya dirinci dari Section `ShowTransfer` (FlowAction `ViewTransferDtl`,
> pra-proses `SetValueKomite`, pasca `KomitePostAct`):
>
> - Judul **"CLAIM COMMITTEE -"** + **ADJUSTMENT** (`TransferType == '2'`) / **REJECT** (`'3'`) / **CLOSE** (`'4'`).
> - Blok: Policy Detail(s), Object Detail + `SpreadingDetail`, Claim Details, grid adjustment KMT + `DetailAdjustmentFac`
>   (View Retro 10015), History Adjustment, Total Adjustment (+ "Total in IDR", `SetValueKomite` S13), teks komite, dan
>   **List of Committee** (hanya `TransferType == 2`; kolom Committee / Status / Date Approve / Comment).
> - Isian: `AcceptStatus` (dropdown, wajib), `Adjustment.IsProposeClose` / `Adjustment.IsPropReserved` (tampil hanya
>   TT2, **nonaktif bila `KomiteCount != '1'`**), Note (`Comment`, wajib); tombol Cancel / Submit.
> - Tata letak pola Komite Prop: ubin ringkasan, kartu berjudul, pasangan label-nilai; **tangga dan kartu keputusan di
>   bawah semua rincian** (keputusan work owner 09-10-2026 untuk Komite Prop).

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K4** — ⭐ **Ditiru apa adanya** — usulan tindak lanjut dibuat **penyetuju pertama**; jenjang berikutnya **menilai**, ⛔ tidak mengganti

## Yang harus diuji

- [ ] Penyetuju **pertama** dapat menyunting kedua kotak centang usulan
- [ ] Penyetuju **kedua ke atas** melihat keduanya **terkunci**
- [ ] ⭐ **90 medan lain TIDAK terkunci** oleh aturan itu — ⛔ jangan menguncinya karena salah membaca
- [ ] Layar menampilkan rekapitulasi **dalam mata uang asli** dan **total dalam rupiah**
- [ ] Anggota komite dapat menuliskan **catatan** pada keputusannya
- [ ] ⛔ Pembandingan memakai **bilangan dengan bilangan** — ⚠️ Pega memaafkan teks, sistem baru tidak

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **11** | ⚠️ 29 medan bergantung jenis objek | ⚠️ menahan bentuk tampilan objek |

## Seam & verifikasi

**Seam:** lapisan layanan komite untuk data; layar hanya menampilkan.
1. Buka layar sebagai penyetuju **pertama** ⇒ ⭐ kedua kotak centang **dapat disunting**.
2. Buka sebagai penyetuju **kedua** ⇒ ⭐ keduanya **terkunci**.
3. Periksa 90 medan lain pada penyetuju kedua ⇒ ⛔ **tidak ikut terkunci**.
4. Bandingkan total rupiah dengan hitungan manual ⇒ ⭐ cocok.
