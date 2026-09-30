# Modul `treatycontractout` — Treaty Contract Out

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di
sini (struktur tim satu folder per modul, keputusan work owner 30-09-2026). Commit Anda menyentuh
folder ini saja; berkas di luarnya milik tim inti (`.github/CODEOWNERS`).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `treatycontractout` |
| Folder korpus | `Treaty Contract Out` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-TREATYCONTRACTOUT` |
| Status | dimigrasi |
| Rentang migrasi | `300-319` |
| Slot menu | `956-957` |
| Prefix rute API | `/api/treaty-contract-out` (+ pekerja latar antrean lampiran) |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda; akun sebenarnya diisi work owner di `.github/CODEOWNERS`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `models/` `repository/` `services/` `handlers/` `modul.go` — paket Go `nusantarare/modul/treatycontractout/backend/...`; **tanpa** `migrations/` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `menu.ts` `rute.tsx` `proporsi.ts` dan berkas `*.test.ts` |
| `docs/` | spec, tiket (`issues/`), grilling, PARITAS, LAPORAN, OQ — dulu `.scratch/treaty-contract-out/` |

## Migrasi

Rentang `300-319` **tetap kosong**: tco4 *(keputusan work owner 29-09-2026)* — modul ini menulis dan
membaca tabel warisan, nol tabel baru (`TestTCONolTabelBaru`). Slot menu `956-957` untuk butir
`M_NAV_MENU` baru bila kelak ada (bentuk SQL-nya di `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 6);
berkas slot tidak membuat tabel, jadi tidak melanggar tco4.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/treatycontractout/...
go test -tags db ./modul/treatycontractout/...    # tanpa ORACLE_DSN: uji db SKIP dengan pesan
npx vitest run modul/treatycontractout
```

`npx tsc --noEmit` memeriksa seluruh aplikasi sekaligus — tipe tidak dapat diperiksa per folder.
Sebelum pull request, jalankan gerbang lengkap (`docs/bersama/PANDUAN-TIM-PER-MODUL.md` di akar repo, bab 8).

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per
butir. Penjaganya berlaku untuk setiap modul; yang KHUSUS modul ini dinyatakan di sini, supaya
mengubahnya tidak pernah menyunting berkas di luar folder ini. Judul yang tidak ada berarti modul
ini tidak menyatakan apa pun untuk jenis itu. Nilai di dalam `` ` `` dibaca apa adanya.

### Tabel warisan: dibaca, tidak dibuat

Tabel yang dokumen STRUKTUR modul ini gambarkan tetapi SENGAJA tidak dibuat migrasi mana pun
(`TestKolomDDLCocokDenganStruktur`, `TestTabelBukanMilikKitaTidakDibuat`). Mencabut satu baris =
kepemilikan tabel berpindah — keputusan work owner.

| Tabel | Alasan |
| --- | --- |
| `TREATYYEAR` | tabel warisan POOLDATA yang Treaty Contract Out tulis dan baca tanpa membuatnya (tco4, keputusan work owner 29-09-2026) |
| `TREATYCONTRACT` | tabel warisan POOLDATA yang Treaty Contract Out tulis dan baca tanpa membuatnya (tco4, keputusan work owner 29-09-2026) |
| `TREATYREINSURER` | tabel warisan POOLDATA yang Treaty Contract Out tulis dan baca tanpa membuatnya (tco4, keputusan work owner 29-09-2026) |
| `MTREATYSECURITY` | tabel warisan POOLDATA yang Treaty Contract Out tulis dan baca tanpa membuatnya (tco4, keputusan work owner 29-09-2026) |
| `TREATYBUSINESS` | tabel warisan POOLDATA yang Treaty Contract Out tulis dan baca tanpa membuatnya (tco4, keputusan work owner 29-09-2026) |
| `PROPORTIONALARRG` | tabel warisan POOLDATA yang Treaty Contract Out tulis dan baca tanpa membuatnya (tco4, keputusan work owner 29-09-2026) |
| `M_ATTACHMENTTREATY_2` | tabel warisan POOLDATA yang Treaty Contract Out tulis dan baca tanpa membuatnya (tco4, keputusan work owner 29-09-2026) |
| `T_STORAGE_IMAGE` | tabel warisan POOLDATA yang Treaty Contract Out tulis dan baca tanpa membuatnya (tco4, keputusan work owner 29-09-2026) |

### Penyuntikan wajib di handler

Penyusun layanan di `backend/handlers/<berkas>` yang bawaannya ber-stub, dan penyuntikan yang wajib
menyertainya di berkas yang sama (`TestHandlerMenyuntikkanImplementasiNyata`).

| Berkas | Penyusun | Wajib | Catatan |
| --- | --- | --- | --- |
| `rute_treaty_contract_out.go` | `svc.JenisReasuransiTreaty()` | `DenganPembaca(services.PembacaJenisReasuransiOracle(svc))` | Treaty Contract Out tiket 02 (aditif 28-09-2026): pembaca master jenis reasuransi bawaannya gagal terang; handler memasang Oracle-nya. |
| `rute_treaty_contract_out.go` | `svc.GrupTreaty()` | `DenganPembaca(services.PembacaGrupTreatyOracle(svc))` | Tiket 03: master grup treaty dan gudang tahun treaty. |
| `rute_treaty_contract_out.go` | `svc.TahunTreatyTCO()` | `DenganGudang(services.GudangTahunTreatyOracle(svc))`, `DenganGrup(services.PembacaGrupTreatyOracle(svc))` |  |
| `tco_kontrak.go` | `svc.KontrakTreatyTCO()` | `DenganGudang(services.GudangKontrakOracle(svc))`, `DenganTahun(services.GudangTahunTreatyOracle(svc))`, `DenganJenis(services.PembacaJenisReasuransiOracle(svc))` | Tiket 04 (aditif 29-09-2026): gudang kontrak, tahun induk, daftar jenis. |
| `tco_reinsurer.go` | `svc.ReinsurerTCO()` | `DenganGudang(services.GudangReinsurerOracle(svc))`, `DenganKontrak(services.PemegangKontrakOracle(svc))`, `DenganTahun(services.GudangTahunTreatyOracle(svc))`, `DenganMaster(services.MasterReinsurerOracle(svc))` | Tiket 05 (aditif 29-09-2026): gudang reinsurer, kontrak, tahun, master AGENT. |
| `tco_business.go` | `svc.BusinessTCO()` | `DenganGudang(services.GudangBusinessOracle(svc))`, `DenganKontrak(services.PemegangKontrakOracle(svc))`, `DenganTahun(services.GudangTahunTreatyOracle(svc))`, `DenganMaster(services.MasterBusinessOracle(svc))` | Tiket 07 (aditif 29-09-2026): gudang bisnis, kontrak, tahun, master BUSINESS. |
| `tco_klausul.go` | `svc.KlausulTCO()` | `DenganGudang(services.GudangKlausulOracle(svc))`, `DenganMaster(services.MasterKlausulOracle(svc))`, `DenganTahun(services.PengunciTahunOracle(svc))`, `DenganJenis(services.PembacaJenisReasuransiOracle(svc))`, `DenganKurs(services.PembacaKursOracle(svc))` | Tiket 08 (aditif 29-09-2026): gudang klausul, tiga master, tahun, jenis reasuransi. |
| `tco_kaskade.go` | `svc.KaskadeTCO()` | `DenganKaskade(services.KaskadeOracle(svc))`, `DenganKontrak(services.PemegangKontrakOracle(svc))`, `DenganTahun(services.GudangTahunTreatyOracle(svc))`, `DenganReinsurer(services.GudangReinsurerOracle(svc))` | Tiket 10 (aditif 29-09-2026): pelaksana kaskade, kontrak, tahun, reinsurer, jejak. tco4: nol jejak modul (T_TREATYCO_JEJAK dibuang; Pega tidak mencatatnya). |
| `tco_kurs.go` | `svc.KursTCO()` | `DenganTahun(services.GudangTahunTreatyOracle(svc))`, `DenganMaster(services.MasterKursOracle(svc))`, `DenganMataUang(services.MataUangOracle(svc))` | Tiket 11 (aditif 29-09-2026): tahun, master kurs, master mata uang. |
| `tco_security.go` | `svc.SecurityTCO()` | `DenganGudang(services.GudangSecurityOracle(svc))`, `DenganReinsurer(services.GudangReinsurerOracle(svc))`, `DenganKontrak(services.PemegangKontrakOracle(svc))`, `DenganTahun(services.GudangTahunTreatyOracle(svc))`, `DenganMaster(services.MasterReinsurerOracle(svc))` | Tiket 06 (aditif 29-09-2026): gudang security, reinsurer induk, kontrak, tahun, master AGENT. |
| `tco_lampiran.go` | `svc.LampiranTahunTCO()` | `DenganGudang(services.GudangLampiranOracle(svc))`, `DenganKategori(services.KategoriLampiranOracle(svc))`, `DenganAntrean(services.AntreanLampiranOracle(svc))`, `DenganPenyimpanan(services.PenyimpananLampiranTCO(svc))`, `DenganTahun(services.GudangTahunTreatyOracle(svc))` | Tiket 12 (aditif 29-09-2026): lima pasangan lampiran; bawaannya gagal terang. OQ-TCO-08: pemilih stub/nyata, bukan stub mati. |

## Brief acuan

`PROMPT-IMPLEMENTASI-MODUL-TREATY-CONTRACT-OUT.md`, `PROMPT-LANJUTAN-TREATY-CONTRACT-OUT-*.md` (folder `OUTPUT_HASIL_RNM\`).
