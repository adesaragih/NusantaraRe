// Layar B dan C — form kontrak Treaty In beserta strip tabnya.
//
// ⛔ DISALIN dari ekspor Pega 2026-09:
//   tata letak + ikatan medan  `Section/TreatyInNONProportional.xml`
//   kepala (ID, Reinsurance Type)  `Section/InputTreatyInOffer.xml`
//   strip tab  `Section/TreatyInTabsProportional.xml` dan
//              `Section/TreatyInTabsNonProportional.xml`
// Jejak per medan — nama rule dan posisi bitanya — ada di `../labels.ts`.
//
// ⚠️ DUA himpunan tab, dan keduanya TIDAK sama. Yang diperkirakan: satu daftar
// sepuluh tab yang dialihkan radio. Yang ekspor katakan: cabang proporsional
// punya sebelas tab, non-proporsional punya dua belas, dan hanya ENAM namanya
// yang sama. Membangun satu daftar akan memperlihatkan tab yang di sistem lama
// tidak pernah ada pada cabang itu.
//
// ⛔ BLOK BERGARIS MATI TIDAK DIBANGUN. Dua di `TreatyInTabsProportional.xml`
// berkondisi `pyContainerVisibleWhen = 1=2`: satu blok tanpa judul di dalam
// Co-Ins Scale (@1095716), dan satu SALINAN KEDUA "Account Reporting Period"
// (@1300558). Yang dibangun salinan hidupnya, @45828.

import { useEffect, useId, useState } from 'react'

import { Area, Field, Gagal, Memuat, Modal, Panel, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { PilihCari as Pilih } from '../../../../inti/frontend/components/ui/pilihSaring'
import { ambilSesiSaya } from '../../../../inti/frontend/klien'
import {
  ambilDaftarAsalBisnis,
  ambilDaftarCedant,
  ambilDraftSalinan,
  ambilKontrakWarisan,
  ambilOpsiKepala,
  ambilOpsiLimits,
  kirimKontrak,
  simpanKontrak,
  simpanSalinan,
  type HasilSimpan,
  type MasukanKirim,
  type BarisEgnpi,
  type BarisRetensi,
  type KontrakWarisan,
  type LimitsAkar,
  type OpsiKepala,
  type PilihanWarisan,
  type ShareNP,
  type SimpulLimit,
} from '../api'
import type { ModeForm } from '../mode'
import { PesanMedanAgen, useDaftarNegatifAgen } from '../components/PesanDaftarNegatif'
import {
  FORM_KONTRAK,
  SYARAT_TAB_NON_PROPORSIONAL,
  SYARAT_TAB_PROPORSIONAL,
  type SyaratTabKontrak,
  TAB_NON_PROPORSIONAL,
  TAB_PROPORSIONAL,
  GRID_TAMBAH,
} from '../labels'

// ⭐ KOMPONEN DIPECAH 5 Oktober 2026 — pemindahan MURNI, nol perubahan
// perilaku. Contohnya `modul/masterproductnamelife`, yang komponen
// terbesarnya 372 baris sementara layar ini dahulu 1.826 dalam satu berkas.
import PanelHistory from '../components/PanelHistory'
import PanelPolisProduksi from '../components/PanelPolisProduksi'
import TanggalRedup from '../components/TanggalRedup'
import { KotakTanggalKetik } from '../components/TanggalKetik'
import { DropdownDaftar } from '../components/IsianAuto'
import { selAngka } from '../components/angka'
import { formatDate } from '../../../../inti/frontend/lib/format'
import { akhirSetahunSesudah, keKabel, keSimpan, tahunDariMulai } from '../components/tanggalIso'
import PanelLampiran from '../components/PanelLampiran'
import TabLimitsNonProp from '../components/TabLimitsNonProp'
import TabLimitsProp from '../components/TabLimitsProp'
import TabShareNonProp from '../components/TabShareNonProp'
import TabShareProp from '../components/TabShareProp'
import TabInfoSubmit from '../components/TabInfoSubmit'
import TabAchievement from '../components/TabAchievement'
import TabEgnpi from '../components/TabEgnpi'
import TabRetensi from '../components/TabRetensi'
import TabPortofolio from '../components/TabPortofolio'
import TabCoInsScale from '../components/TabCoInsScale'
import TabTeksPanjang from '../components/TabTeksPanjang'
import TabEventLimits from '../components/TabEventLimits'
import DropdownWarisan from '../components/DropdownWarisan'
import TabReportingPeriod from '../components/TabReportingPeriod'
import TabAkumulasi from '../components/TabAkumulasi'
import { saringAngka } from '../components/saringAngka'
import { TataPegaBlok } from '../components/tataPega'
import { bacaProperti, PenyediaHalaman, usePenampungHalaman } from '../halaman'
import { bolehActions, bolehSave, DIVISI_IT, kursBerubah, susunDokumen } from '../simpanDokumen'
import { PILIHAN_AKSEPTASI, TOMBOL_TULIS } from '../labelsTulis'
import TabAngsuran from '../components/TabAngsuran'

// ⛔ DIEKSPOR ULANG, bukan didefinisikan di sini: uji yang mengimpornya dari
// halaman ini tetap berjalan tanpa disunting — itulah bukti pemindahannya
// murni.
export { padankanDesimal, selAngka } from '../components/angka'





/**
 * Susunan dua kolom form, **dalam urutan rujukan**.
 *
 * ⛔ Dipisahkan dari komponennya supaya ia dapat diuji tanpa merender — dan
 * supaya panjang kolom yang BERBEDA menjadi hal yang gagal bila seseorang
 * kelak meratakannya menjadi grid pengisi baris. Kolom kanan memang lebih
 * panjang; itu bentuk rujukannya, bukan kelalaian.
 */
export const tataLetakKolom = {
  kiri: [
    'Treaty Contract Name',
    'Contract Ref No',
    'Teritorial Scope',
    'Bordereaux',
    'Bordereaux Note',
  ],
  kanan: [
    'Commencement',
    'Termination',
    'Treaty Year',
    'Accounting Mode',
    'Ceding',
    'RNM as Treaty Leader',
    'Source of Business',
  ],
} as const

/** Kedua nilai radio `TreatyIn.ProportionType`. */
export const PROPORSIONAL = 'Proportional'
export const NON_PROPORSIONAL = 'Non Proportional'

/**
 * Himpunan tab yang berlaku bagi sebuah cabang.
 *
 * ⛔ Dipisahkan sebagai fungsi supaya ia dapat diuji tanpa merender apa pun —
 * dan supaya perbedaan kedua himpunan menjadi hal yang gagal bila seseorang
 * kelak meratakannya.
 */
export function tabUntuk(jenis: string, syarat?: SyaratTabKontrak): readonly string[] {
  const nonProp = jenis === NON_PROPORSIONAL
  const semua = nonProp ? TAB_NON_PROPORSIONAL : TAB_PROPORSIONAL
  // ⛔ Syaratnya BERCABANG. `Retro` ada di kedua daftar dan hanya yang
  // proporsional bersyarat — lihat `SYARAT_TAB_*` di `labels.ts`.
  const syaratCabang = nonProp ? SYARAT_TAB_NON_PROPORSIONAL : SYARAT_TAB_PROPORSIONAL
  // ⛔ Kontrak BARU (nol `syarat`) memperlihatkan SELURUH tab.
  //
  // Syarat tampil dibaca dari dokumen kontrak; kontrak yang belum punya
  // dokumen tidak punya jawabannya, dan menyembunyikan tab karena
  // pertanyaannya belum dapat dijawab menyembunyikannya SELAMANYA — tab
  // yang tersembunyi tidak pernah diisi, dan yang tidak pernah diisi tidak
  // pernah memenuhi syaratnya.
  if (syarat === undefined) return semua
  return semua.filter((t) => {
    const uji = syaratCabang[t]
    return uji === undefined || uji(syarat)
  })
}

/** Satu baris grid Rate of Exchange; kosong adalah keadaan awal yang sah. */
interface BarisKurs {
  mataUang: string
  /** `.CurrencyID` — nilai dropdown sel Currency. */
  mataUangID: string
  nilaiKeIDR: string
  /** Bentuk TAMPIL (`dd/mm/yy`) — mode lihat. */
  berlakuDari: string
  berlakuSampai: string
  /** Bentuk TERSIMPAN (`YYYYMMDD`) — kotak tanggal mode ubah. */
  berlakuDariAsli: string
  berlakuSampaiAsli: string
}

export interface FormKontrakProps {
  /**
   * Pengenal kontrak yang dibuka; TEKS KOSONG berarti kontrak baru
   * (tombol `Add`).
   *
   * ⛔ Teks, bukan bilangan: layar daftar membaca `POOLDATA.TREATY_IN`, dan
   * kolom `ID` di sana `VARCHAR2(100)`.
   */
  idKontrak: string
  /**
   * ⭐ `ubah` dari tombol `Edit` (dan `Add`), `lihat` dari `View`. Mode lihat
   * mematikan seluruh isian lewat `<fieldset disabled>` dan tidak merender
   * tombol `Add`/`Delete` tabel.
   */
  mode?: ModeForm
  onKembali: () => void
  /**
   * Kontrak BARU tersimpan — pengenalnya baru lahir di server; rute membuka
   * ulang form dengan pengenal itu.
   */
  onTersimpan?: (id: string) => void
  /**
   * ⭐ Tombol `Copy` daftar — pengenal kontrak SUMBER. Bersama `idKontrak`
   * kosong ia berarti DRAF salinan (`Activity/TreatyInCopy.xml`: `ID =
   * "UnknownId"`, `OLDID` = sumber): isinya dimuat dari sumbernya, dan
   * Save/Submit/Decline melahirkan kontrak BARU lewat `simpanSalinan`.
   */
  salinDari?: string
}

export default function FormKontrakTreatyIn({ idKontrak, mode: modeRute = 'lihat', onKembali, onTersimpan, salinDari }: FormKontrakProps) {
  // Draf `Copy` yang BELUM disimpan — sesudah Save, rute membuka pengenal barunya.
  const sumberSalin = idKontrak === '' ? (salinDari ?? '') : ''
  // Pilihan dropdown kepala — nilai TERSIMPAN ↔ label, disusun services.
  // Kontrak yang dibuka membawanya; kontrak BARU memintanya sendiri.
  const [opsi, setOpsi] = useState<OpsiKepala | null>(null)
  useEffect(() => {
    if (idKontrak !== '') return
    let dibuang = false
    ambilOpsiKepala()
      .then((o) => {
        if (!dibuang) setOpsi(o)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [idKontrak])
  // Kontrak warisan yang sedang dibuka; null selama memuat, dan tetap null
  // untuk kontrak BARU (pengenal kosong) — yang memang tidak punya warisan.
  const [warisan, setWarisan] = useState<KontrakWarisan | null>(null)
  // ⛔ BERKAS TUNTAS = BACA (9 Oktober 2026) — laporan pemakai: sesudah
  // approve, status `Resolve Complete` tetapi isian masih dapat diubah. Di
  // Pega halaman tertutup sesudah akseptasi dan kontrak tuntas/ditolak tidak
  // dapat dibuka sunting (Edit daftar & Save tersembunyi). Revision
  // mengosongkan `StatusAkseptasi`, jadi jalurnya tidak terkunci.
  const statusMuat = warisan?.statusAkseptasi ?? ''
  const statusTuntas = statusMuat === 'Resolve Complete' || statusMuat === 'Decline'
  // ⭐ Force Edit (9 Oktober 2026) — padanan `Force Edit (dev)` Pega
  // (`TreatyInActionButtons`, DT `TreatyInForceEdit`: `ViewState = 0` →
  // refresh section): form yang terkunci (View / posisi bukan milik akun /
  // Resolve Complete / Decline) menjadi dapat diubah. Khusus divisi IT.
  //
  // ⛔ MENYIMPANG dari Pega untuk `Resolve Complete` (keputusan pemakai
  // 9 Oktober 2026): Save Pega ber-syarat `StatusAkseptasi != 'Resolve
  // Complete'`, jadi Force Edit di sana tidak dapat menyimpan. Di sini Save
  // tampil sesudah Force Edit, dan server menerimanya dari akun IT tanpa
  // mengubah status maupun posisi.
  const [paksaUbah, setPaksaUbah] = useState(false)
  const mode: ModeForm = paksaUbah ? 'ubah' : statusTuntas ? 'lihat' : modeRute
  const bisaUbah = mode === 'ubah'
  const [galat, setGalat] = useState<unknown>(null)
  const [memuat, setMemuat] = useState(idKontrak !== '')

  const [jenis, setJenis] = useState(PROPORSIONAL)
  const [nama, setNama] = useState('')
  const [rujukan, setRujukan] = useState('')
  const [wilayah, setWilayah] = useState('')
  const [bordereaux, setBordereaux] = useState('')
  const [bordereauxNote, setBordereauxNote] = useState('')
  const [mulai, setMulai] = useState('')
  const [berakhir, setBerakhir] = useState('')
  // ⭐ `Treaty Year` BERKEADAAN, bukan turunan tampilan.
  //
  // `DataTransform/TreatyInSetTreatyYear.xml` menyetelnya pada event
  // `change` milik `TreatyIn.Commencement` — jadi ia diisi dari kolom saat
  // kontrak dibuka, lalu DIHITUNG ULANG setiap Commencement berubah.
  const [tahunTreaty, setTahunTreaty] = useState('')
  const [pembukuan, setPembukuan] = useState('')
  // ⛔ MEDAN KEDUA, bukan medan yang sama. `AccountingModeNonProp` berdiri
  // sendiri di ekspor dengan `pyCondition` cabangnya sendiri; satu keadaan
  // untuk keduanya akan membuat berpindah cabang menimpa nilai cabang yang
  // ditinggalkan.
  const [pembukuanNonProp, setPembukuanNonProp] = useState('')
  const [cedant, setCedant] = useState('')
  const [asalBisnis, setAsalBisnis] = useState('')
  // ⭐ PENGENAL ikut disimpan, 6 Oktober 2026. Pega mengisi KEDUANYA saat
  // memilih — `DataTransform/TreatyInSetReinsured.xml`:
  //   1.1 TreatyIn.Ceding               := param.name
  //   1.2 TreatyIn.CedingID             := param.id
  //   2.1 TreatyIn.LeadingReinsSource   := param.name
  //   2.2 TreatyIn.LeadingReinsSourceID := param.id
  // Layar ini dahulu hanya menjalankan 1.1 dan 2.1, sehingga memilih ulang
  // meninggalkan pengenal LAMA di samping nama BARU — dan Submit Pega
  // menolak kontrak yang pengenalnya kosong.
  const [idCedant, setIdCedant] = useState('')
  const [idAsalBisnis, setIdAsalBisnis] = useState('')
  // `TreatyInCheckCedingBlacklist` — tombol `Edit` menjalankannya sesudah
  // `SetTreatyIn_Act`: sekali, saat kontrak lama selesai dimuat di mode ubah.
  const pesanAgen = useDaftarNegatifAgen(bisaUbah && idKontrak !== '' && !memuat, idKontrak, idCedant, idAsalBisnis)
  const [pemimpin, setPemimpin] = useState(false)
  // ⭐ Grid Rate of Exchange dibaca dari dokumen warisan kontrak ini, bukan
  // dari keadaan kosong. `CurrencyList` berisi di 297 dari 300 dokumen yang
  // disapu — "No items" selama ini bukan karena datanya tidak ada, melainkan
  // karena tidak ada yang membacanya.
  const [kurs, setKurs] = useState<BarisKurs[]>([])
  // ⭐ KEADAAN BERSAMA tab Limits Non-Prop dan Share Non-Prop — clipboard
  // Pega (`TreatyIn.Limits`, `TreatyIn.Share`) hidup selama kasus terbuka,
  // dan `Update Summary` tab Share MEMBACA layer yang sedang disunting di tab
  // Limits. Tab di sini dirender ulang tiap pindah tab, jadi keadaannya
  // disimpan di form; `null` = belum disentuh, pakai yang dimuat.
  // ⭐ PENAMPUNG HALAMAN `TreatyIn` (`../halaman.tsx`) — padanan clipboard
  // Pega untuk tab cabang Proporsional: isian tab bertahan saat pindah tab,
  // dan tab yang saling bergantung (Limits ↔ Share ↔ Achievement In IDR,
  // Reporting Period → Accumulation) membaca properti yang SAMA. Hidup di
  // form, di atas strip tab — bukan tabel; Save/Submit kelak membacanya.
  const penampung = usePenampungHalaman()
  const [limitsNP, setLimitsNP] = useState<{ layers: SimpulLimit[]; akar: LimitsAkar } | null>(null)
  const [shareNP, setShareNP] = useState<ShareNP | null>(null)
  // Sel grid kurs yang sedang diketik (`baris:kolom`), atau null. Dipakai
  // supaya pemformat tidak melawan pengetik — lihat catatan di gridnya.
  const [selDiketik, setSelDiketik] = useState<string | null>(null)
  // Pilihan dropdown Currency grid kurs — `BrowseCurrency_RD` (CURRENCY,
  // `Currency != 'ITL'`), rute yang SAMA dengan dropdown mata uang tab
  // Limits. Hanya dibutuhkan selama gridnya dapat disunting.
  const [mataUangKurs, setMataUangKurs] = useState<PilihanWarisan[]>([])
  useEffect(() => {
    if (!bisaUbah) return
    let dibuang = false
    ambilOpsiLimits()
      .then((o) => {
        if (!dibuang) setMataUangKurs(o.mataUang)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [bisaUbah])
  const ubahKurs = (i: number, sebagian: Partial<BarisKurs>) => {
    setKurs(kurs.map((x, j) => (j === i ? { ...x, ...sebagian } : x)))
  }

  // ---------------------------------------------------------------------
  // ⭐ TOMBOL TULIS — Save, Submit, Actions, Decline offer (keputusan
  // pemilik proses 6–7 Oktober 2026). Isian masuk basis data HANYA di sini,
  // dari tombolnya; sasarannya tabel masing-masing (`T_TREATY_*`, kepala
  // `TREATY_IN`, kurs `TREATYEXCHANGEYEARLY`).
  // ---------------------------------------------------------------------
  const [muatUlang, setMuatUlang] = useState(0)
  // ⭐ GENERASI DATA — naik SETIAP kali kontrak yang dimuat diterapkan
  // (termasuk muat ulang sesudah Save), dan ikut `key` fieldset tab.
  //
  // ⛔ Laporan pemakai 8 Oktober 2026: Currency 100% Limit yang SUDAH
  // tersimpan (`T_TREATY_LIMIT_AMOUNT.CURRENCY = IDR`) tampil "Choose" lagi
  // sesudah Save. Sebabnya: `useProperti` membekukan nilai awal tab
  // (`useState(awal)`) saat tab PERTAMA dirender. Muat ulang mengosongkan
  // penampung, tetapi `key` lama (`id|isi|mode`) tidak berubah — tab tidak
  // dilahirkan ulang dan menyemai ulang data SEBELUM Save. Save berikutnya
  // lalu menulis data basi itu kembali.
  const [generasi, setGenerasi] = useState(0)
  const [sibukTulis, setSibukTulis] = useState(false)
  const [hasilTulis, setHasilTulis] = useState<{ galat: boolean; pesan: string; takTersimpan: string[] } | null>(null)
  // Workbasket pemakai — `OperatorID.pyWorkBasketList`; syarat tampil Actions
  // dan Save. Divisi — pengecualian Force Edit (IT).
  const [workbasket, setWorkbasket] = useState<string[]>([])
  const [divisi, setDivisi] = useState('')
  useEffect(() => {
    let dibuang = false
    ambilSesiSaya()
      .then((p) => {
        if (dibuang) return
        setWorkbasket(p.peran)
        setDivisi(p.divisi)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [])
  // Modal `TreatyInAction` — radio `ChooseStatusAkseptasi` + `Comment`.
  const [actionsBuka, setActionsBuka] = useState(false)
  const [pilihanActions, setPilihanActions] = useState('Accept')
  const [komentarActions, setKomentarActions] = useState('')

  /** Isi layar → dokumen `TreatyIn` ejaan Pega (`simpanDokumen.ts`). */
  const isiTombol = (tambahan: Record<string, unknown> = {}) => ({
    idKontrak,
    dokumen: {
      ...susunDokumen(
        {
          nonProporsional: jenis === NON_PROPORSIONAL,
          nama, rujukan, wilayah, bordereaux, bordereauxNote, mulai, berakhir, tahunTreaty,
          pembukuan, pembukuanNonProp, cedant, idCedant, asalBisnis, idAsalBisnis, pemimpin,
        },
        penampung.halaman,
        limitsNP,
        shareNP,
      ),
      ...tambahan,
    },
    kurs: kursBerubah(kurs, warisan?.kurs ?? [], tahunTreaty),
  })

  /** Satu penekanan tombol tulis; sesudah berhasil, data dibaca ulang dari tabel. */
  const tekanTulis = (jalan: () => Promise<HasilSimpan>) => {
    setSibukTulis(true)
    setHasilTulis(null)
    jalan()
      .then((h) => {
        setHasilTulis({ galat: false, pesan: h.pesan, takTersimpan: h.kunciTakTersimpan })
        if (idKontrak === '') onTersimpan?.(h.id)
        else setMuatUlang((n) => n + 1)
      })
      .catch((e: unknown) => {
        setHasilTulis({ galat: true, pesan: e instanceof Error ? e.message : String(e), takTersimpan: [] })
      })
      .finally(() => {
        setSibukTulis(false)
      })
  }
  /** Draf `Copy`: isi layar + sumbernya — server menyusun clipboard salinan. */
  const isiSalinan = (aksi: '' | 'submit' | 'decline', tambahan: Record<string, unknown> = {}) => {
    const { dokumen, kurs: k } = isiTombol(tambahan)
    return { idSumber: sumberSalin, dokumen, kurs: k, aksi }
  }
  const tekanSave = () => {
    tekanTulis(() => (sumberSalin !== '' ? simpanSalinan(isiSalinan('')) : simpanKontrak(isiTombol())))
  }
  const tekanKirim = (aksi: MasukanKirim['aksi'], pilihan = '', tambahan: Record<string, unknown> = {}) => {
    tekanTulis(() =>
      sumberSalin !== '' && aksi !== 'akseptasi'
        ? simpanSalinan(isiSalinan(aksi, tambahan))
        : kirimKontrak({ ...isiTombol(tambahan), aksi, pilihan }),
    )
  }
  const statusKini = warisan?.statusAkseptasi ?? ''
  // `TreatyIn.RevisionState` TERSIMPAN (migrasi 448) — tombol Revision daftar.
  const tersimpanPenampung: Record<string, string> = warisan?.penampung ?? {}
  const revisi = tersimpanPenampung.RevisionState === '1'
  // Satu-satunya tab yang jalur revisi hidupkan di mode lihat.
  const TAB_REVISI = 'Information & Submit'
  // Syarat tampil Save — `TreatyIn.IsEditData !='1' && StatusAkseptasi != 'Resolve Complete'`.
  // ⭐ 9 Oktober 2026 — DAN akun memegang workbasket posisi berkas (`bolehSave`).
  const saveTampil = bisaUbah && (paksaUbah || bolehSave(warisan?.posisi ?? '', workbasket, statusKini, divisi))
  const actionsTampil = bolehActions(warisan?.posisi ?? '', workbasket, statusKini)
  // Force Edit tampil hanya bagi divisi IT, saat form TERKUNCI (mode View,
  // posisi berkas bukan workbasket akun, Resolve Complete, Decline), bukan
  // untuk kontrak baru.
  const terkunci = !bisaUbah || !bolehSave(warisan?.posisi ?? '', workbasket, statusKini, '')
  const paksaTampil = divisi === DIVISI_IT && idKontrak !== '' && terkunci && !paksaUbah

  // ⛔ Mengisi SELURUH medan dari kontrak nyata, termasuk radio jenisnya.
  // Radio itu memilih strip tab, dan tiket ronde ini menyebutnya: ia
  // disambungkan ke NILAI NYATA, bukan dibiarkan pada bawaan yang pemakai
  // ubah sendiri — kontrak non-proporsional yang membuka strip proporsional
  // memperlihatkan sebelas tab yang tidak satu pun miliknya.
  useEffect(() => {
    if (idKontrak === '' && sumberSalin === '') {
      setMemuat(false)
      return
    }
    let dibuang = false
    setMemuat(true)
    setGalat(null)
    setLimitsNP(null)
    setShareNP(null)
    penampung.kosongkan()
    // ⭐ Draf `Copy` dimuat dari SUMBERNYA, sudah melewati `TreatyInCopy`
    // di server (status/posisi/riwayat). NOL tulisan.
    const muat = sumberSalin !== '' ? ambilDraftSalinan(sumberSalin) : ambilKontrakWarisan(idKontrak)
    muat
      .then((k) => {
        if (dibuang) return
        setWarisan(k)
        setLimitsNP(null)
        setShareNP(null)
        // ⛔ Halaman kasus BARU: tab yang sempat dirender selama memuat sudah
        // menyemai nilai kosong — dibuang supaya tab menyemai dari kontrak.
        penampung.kosongkan()
        // ⭐ Properti penampung yang TERSIMPAN (migrasi `448`) disemai lebih
        // dulu — tab yang dirender sesudahnya membaca nilai ini, bukan
        // bawaan kosongnya sendiri.
        for (const [kunci, v] of Object.entries(k.penampung ?? {})) {
          penampung.ubah(kunci, () => v)
        }
        // ⭐ Larik total yang TERSIMPAN (`T_TREATY_TOTAL`) — tanpa ini total
        // Share Prop / EGNPI / Installment "No items" sesudah Save sampai
        // Refresh ditekan (laporan pemakai 8 Oktober 2026).
        for (const [kunci, v] of Object.entries(k.penampungLarik ?? {})) {
          penampung.ubah(kunci, () => v)
        }
        // Tab dilahirkan ulang atas data INI — lihat `generasi`.
        setGenerasi((n) => n + 1)
        setOpsi(k.opsiKepala)
        // ⛔ `?? []` BUKAN hiasan, dan bukan pula ketidakpercayaan pada
        // backend. Ia lapis kedua dari galat yang menghentikan halaman
        // 6 Oktober 2026 — `Cannot read properties of null (reading
        // 'length')` — ketika jawaban memuat `"kurs": null`. Seluruh medan
        // larik lain di layar ini sudah dibaca lewat `?? []`; yang satu ini
        // lewat `useState`, dan itulah sebabnya hanya ia yang jatuh.
        setKurs(k.kurs ?? [])
        setJenis(k.sifatProporsi === NON_PROPORSIONAL ? NON_PROPORSIONAL : PROPORSIONAL)
        setNama(k.namaKontrak)
        setRujukan(k.nomorRujukan)
        setWilayah(k.lingkupWilayah)
        // ⛔ NILAI TERSIMPAN, bukan label tampil — dropdown memegang nilai dan
        // menampilkan label dari `opsiKepala`. Mengisinya dengan label
        // (`Underwriting Year`) memberi "tidak ada di daftar referensi".
        setBordereaux(k.bordereauxAsli)
        setBordereauxNote(k.bordereauxCatatan)
        // ⛔ Dari `…Asli` (`YYYYMMDD`) ke bentuk KABEL `FieldTanggal`
        // (`DD-MM-YYYY`), BUKAN dari `tanggalMulai` (`01/01/2025`).
        //
        // Bentuk BACA ditolak medan tanggal TANPA BERSUARA, dan bentuk ISO
        // membuat kedua turunannya (`Treaty Year`, `Termination`) kosong
        // begitu pemakai memilih tanggal — sebab `FieldTanggal`
        // mengembalikan `DD-MM-YYYY`, bukan ISO. Satu bentuk kabel saja.
        setMulai(keKabel(k.tanggalMulaiAsli))
        setBerakhir(keKabel(k.tanggalBerakhirAsli))
        setTahunTreaty(k.tahunTreaty)
        setPembukuan(k.caraPembukuanAsli)
        setPembukuanNonProp(k.caraPembukuanNonPropAsli)
        setCedant(k.cedant)
        setAsalBisnis(k.asalBisnis)
        // ⛔ Dari KOLOM `TREATY_IN`, apa adanya. Terukur 6 Oktober 2026:
        // 0 dari 1.854 kontrak bernama tanpa pengenal — jadi pencarian balik
        // nama->ID TIDAK diperlukan di sini, dan sengaja tidak dipasang.
        setIdCedant(k.idCedant)
        setIdAsalBisnis(k.idAsalBisnis)
        // `TreatyLeader` tersimpan sebagai TEKS `"true"`/`"false"`, bukan
        // boolean JSON — sapuan menemukan 103 `true` dan 556 `false`.
        setPemimpin(k.pemimpinTreaty === 'true')
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
      .finally(() => {
        if (!dibuang) setMemuat(false)
      })
    return () => {
      dibuang = true
    }
    // ⭐ `muatUlang` — sesudah tombol tulis berhasil, kontrak dibaca ULANG dari
    // tabel: layar memperlihatkan yang TERSIMPAN, bukan yang diketik.
  }, [idKontrak, sumberSalin, muatUlang])

  // ⭐ `TreatyIn.EDMMaterialType = 1` — mengunci `Contract Ref No`
  // (`pyDisabledWhen`) dan `Bordereaux Note` (`pyReadOnlyCondition`).
  const edmMateri = (warisan?.edmJenisMaterial ?? '') === '1'
  const idBordereauxNote = useId()

  // ⭐ Syarat tampil tab — dari dokumen kontrak, nol untuk kontrak baru.
  const syaratTab: SyaratTabKontrak | undefined =
    warisan === null
      ? undefined
      : {
          retroBerganda: warisan.retroBerganda,
          edmState: warisan.edmState,
          edmJenisMaterial: warisan.edmJenisMaterial,
        }
  // ⭐ Syarat tab dinilai `tabUntuk`; yang DISEMBUNYIKAN pemilik proses
  // (Retro) dibuang di sini, di strip — syaratnya tetap teruji.
  const tab = tabUntuk(jenis, syaratTab)
  const [tabAktif, setTabAktif] = useState<string>(tab[0] ?? '')
  const tabTampil = tab.includes(tabAktif) ? tabAktif : (tab[0] ?? '')

  /**
   * `Treaty Year` — baca-saja, dan TIDAK DIFORMAT.
   *
   * ⭐ PERTANYAAN TERBUKA RONDE SEBELUMNYA DITUTUP oleh gambar `01`: medan
   * ini berbunyi `2025` — tahun polos. Dugaan `05/1` sebagai turunan SALAH;
   * itu artefak form yang belum tersimpan, dan pemilik proses menyatakannya
   * sendiri 5 Oktober 2026.
   *
   * ⛔ `selAngka` TIDAK dipanggil di sini, dan itu disengaja: pemisah ribuan
   * atas sebuah tahun memberi `2.025`.
   *
   * RUMUSNYA juga terbaca, dan ia tetap tidak dihitung di sini.
   *
   * ⭐ `DataTransform/TreatyInSetTreatyYear.xml` (`pyRuleAvailable = Yes`,
   * memo *"termination automatically add 1 year from start date"*):
   *
   *     TreatyIn.TreatyYear  := @substring(TreatyIn.Commencement,0,4)
   *     TreatyIn.Termination := @addCalendar(TreatyIn.Commencement,"1",0,0,0,0,0,0)
   *
   * ⚠️ Dan itu rumus SAAT MASUK, bukan aturan yang terus berlaku. Diadu
   * dengan POOLDATA, 5 Oktober 2026:
   *
   *     TREATYYEAR = SUBSTR(COMMENCEMENT,1,4)       1.849 dari 1.854
   *     TERMINATION = COMMENCEMENT + 12 bulan           9 dari 1.854
   *
   * Baris kedua yang memutuskan: kalau `Termination` sungguh terikat pada
   * `Commencement`, ia akan cocok pada ribuan baris, bukan sembilan. Jadi
   * kedua langkah itu NILAI AWAL yang pemakai ubah sesudahnya — dan
   * menghitung ulang `TreatyYear` di layar akan menimpa lima kontrak yang
   * tahunnya sengaja berbeda dari tahun mulainya.
   *
   * ⛔ Karena itu yang ditampilkan tetap KOLOMNYA apa adanya.
   *
   */
  // ⭐ `TREATYYEAR` adalah KOLOM di `TREATY_IN` — dibaca, bukan dihitung.
  // ⭐ Aturan `TreatyInSetTreatyYear`, disalin APA ADANYA dari ekspor —
  // dijalankan pada `change` Commencement, persis seperti Pega
  // (`pyEvent change` + `postValue` + `refresh`, pemiliknya
  // `TreatyIn.Commencement`).
  //
  // ⚠️ Ia MENIMPA `Termination`, dan itu memang perilaku ekspornya. Diadu
  // dengan POOLDATA 5 Oktober 2026, `TERMINATION = COMMENCEMENT + 12 bulan`
  // hanya cocok pada 9 dari 1.854 — bukan karena aturannya salah, melainkan
  // karena ia NILAI AWAL yang pemakai ubah sesudahnya. Menimpa saat
  // Commencement berubah lalu membiarkannya disunting adalah justru yang
  // menghasilkan sebaran itu.
  const ubahMulai = (v: string) => {
    setMulai(v)
    setTahunTreaty(tahunDariMulai(v))
    setBerakhir(akhirSetahunSesudah(v))
  }

  // ⭐ KETERGANTUNGAN ANTARTAB Non-Prop — 7 Oktober 2026. Satu sumber NILAI
  // AWAL per page list (tab pemiliknya menyemai penampung darinya), dan
  // nilai TERKINI dari penampung untuk tab yang membacanya:
  //   Retention → EGNPI   `TreatyInNonAddItem(egnpi)` memakai `Retention(1)`
  //   EGNPI → Limits      `TotalEgnpi` (`EgnpiTotalList` tiap layer)
  //   Share → Installment `TreatyIn.TotalShareNetNP`
  // Dahulu ketiganya membaca DATA KONTRAK yang dimuat, sehingga isian tab
  // sumbernya tidak pernah sampai ke tab tujuan.
  const retensiAwal: BarisRetensi[] = (warisan?.retensi ?? []).map((b) => ({
    ID: '',
    TreatyGroup: b.kelompokTreaty,
    TreatyGroupID: '',
    Currency: b.mataUang,
    CurrencyID: '',
    Amount: b.jumlah,
    ClassOfBusiness: b.kelasBisnis,
    ClassOfBusinessID: '',
    Note: b.keterangan,
  }))
  const egnpiAwal: BarisEgnpi[] = (warisan?.egnpi ?? []).map((b) => ({
    ID: '',
    TreatyGroup: b.kelompokTreaty,
    TreatyGroupID: '',
    AsDate: b.perTanggal,
    Proportion: b.proporsi,
    Currency: b.mataUang,
    CurrencyID: '',
    Amount: b.jumlah,
    AmountIDR: b.jumlahIDR,
    ClassOfBusiness: b.kelasBisnis,
    ClassOfBusinessID: '',
    Note: b.keterangan,
  }))
  const retensiKini = (bacaProperti(penampung.halaman, 'Retention') as BarisRetensi[] | undefined) ?? retensiAwal
  const egnpiKini = (bacaProperti(penampung.halaman, 'EGNPI') as BarisEgnpi[] | undefined) ?? egnpiAwal
  const netPremiumKini = (shareNP ?? warisan?.shareNP)?.Total.TotalShareNetNP ?? []

  return (
    <>
    <div className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{FORM_KONTRAK.judul}</h2>
        <button type="button" className="btn btn--ghost" onClick={onKembali}>
          {FORM_KONTRAK.kembali}
        </button>
      </header>

      {galat !== null && <Gagal galat={galat} />}
      {memuat && <Memuat />}

      {/* ⭐ MODE LIHAT: `<fieldset disabled>` mematikan SELURUH isian di
          dalamnya — medan, dropdown, kotak centang, tombol pemilih. Strip
          tab dan `Close` berdiri di luarnya, jadi tetap hidup. */}
      <fieldset className="trin__mode" disabled={!bisaUbah}>
      <Panel judul={FORM_KONTRAK.judul}>
        {/* ⭐ TATA LETAK PEGA (8 Oktober 2026) — `Section/InputTreatyInOffer.xml`
            `Inline grid double` L1394: blok ID + Reinsurance Type (L1694)
            di separuh KIRI, panel "Existing Policy for Master ID" (L2663)
            mulai di tengah — separuh KANAN, bukan didorong ke tepi kanan.
            Kelas `trin__tata--g2` menimpa `display: flex` `.trin__kepala`
            (aturannya lebih belakang di `treatyin.css`). */}
        <div className="trin__kepala trin__tata trin__tata--g2">
          {/* ⭐ ID dan Reinsurance Type BERTUMPUK di kiri atas, bukan sebaris.
              Di gambar Pega keduanya blok padat di pojok, dan tabel
              "Existing Policy" berdiri sendiri di kanan. Tanpa kolom ini
              keduanya ikut ditengahkan terhadap tabel yang jauh lebih
              tinggi, dan melayang di tengah layar. */}
          <div className="trin__kepala-kiri">
          <span className="trin__id">
            {FORM_KONTRAK.id}: {idKontrak === '' ? <span className="trin__redup">—</span> : idKontrak}
          </span>
          {/* Radio, bukan daftar pilihan: ekspor memakai dua pilihan sejajar,
              dan ia yang menentukan tab mana yang tampil. */}
          <fieldset className="trin__radio">
            <legend>{FORM_KONTRAK.jenisReasuransi}</legend>
            {[
              { nilai: PROPORSIONAL, label: FORM_KONTRAK.proporsional },
              { nilai: NON_PROPORSIONAL, label: FORM_KONTRAK.nonProporsional },
            ].map((o) => (
              <label key={o.nilai}>
                <input
                  type="radio"
                  name="jenis-reasuransi"
                  value={o.nilai}
                  checked={jenis === o.nilai}
                  onChange={() => {
                    setJenis(o.nilai)
                    // ⭐ DT `TreatyInSetPeriod` — satu langkah, dan itu
                    // seluruh isinya:
                    //
                    //   TreatyIn.ReportingPeriod := "quarter"
                    //
                    // Ia terpasang pada peristiwa `change` radio inilah
                    // (`Section/InputTreatyInOffer.xml`, aksi `refresh`
                    // ber-pra-DT), bukan pada tab Reporting Period.
                    //
                    // ⚠️ TANPA SYARAT, termasuk ketika pemakai memilih
                    // kembali jenis yang sama: ekspor nol memeriksa nilai
                    // lama, dan menambah pemeriksaan itu akan membuat
                    // perilaku kita berbeda pada satu-satunya kasus yang
                    // membedakannya.
                    penampung.ubah('ReportingPeriod', () => 'quarter')
                  }}
                />
                {o.label}
              </label>
            ))}
          </fieldset>
          </div>

          {/* ⭐ Kanan atas, sejajar blok ID dan Reinsurance Type — persis
              letaknya di gambar 01 dan 26. */}
          <PanelPolisProduksi baris={warisan?.polisProduksi ?? []} />
        </div>

        {/* ⛔ DUA KOLOM YANG MENGALIR SENDIRI, bukan `.form-grid`.

            `.form-grid` dikunci dua kolom dari `treatyin.css`, dan itu perlu
            tetapi TIDAK cukup: ia mengisi BARIS DEMI BARIS, sementara kolom
            kanan rujukan lebih panjang (tujuh butir lawan lima) dan tiga
            butir terakhirnya berdiri tanpa pasangan di kiri.

            Bentuk sebelumnya menirunya dengan `<span />` pengisi. Pengisi itu
            bekerja hari ini dan akan bergeser diam-diam pada medan berikutnya
            yang ditambahkan — tanpa satu uji pun berbunyi. Dua kolom yang
            mengalir sendiri tidak punya cacat itu, dan `tataLetakKolom` di
            bawah membuat susunannya DAPAT DIUJI tanpa merender.

            ⭐ TATA LETAK PEGA (8 Oktober 2026) — `Section/TreatyInNONProportional.xml`
            `Inline grid double` L543 → kolom kiri `Stacked with labels left`
            L845, kolom kanan `Stacked with labels left` L2575 (sel ketiga
            `Spacer` `1=2`). `.trin__dwikolom` = dua kolom sama lebar,
            `.trin__kolom` = label di kiri medan. */}
        <div className="trin__dwikolom">
          <div className="trin__kolom">
            <Field label={FORM_KONTRAK.namaKontrak} value={nama} onChange={setNama} />
            {/* ⭐ RALAT 7 Oktober 2026 — ketiga medan kepala yang dahulu
                menjadi medan mati "Tidak ada di dokumen sistem lama"
                kini SELALU dirender, persis ekspor
                `Section/TreatyInNONProportional.xml`: sel `ContractRefNo`
                dan `BordereauxNote` `pyVisible ALWAYS`, sel `TreatyLeader`
                tanpa syarat tampil sama sekali. Ada-tidaknya nilai di
                dokumen lama bukan syarat tampil di Pega; tabel `T_TREATY_*`
                yang kosong pun membuat ketiganya mati pada SETIAP kontrak.

                `pxTextInput` · `pyReadOnlyCondition TreatyIn.IsEditData= 1`
                · `pyDisabledWhen TreatyIn.EDMMaterialType = 1`.
                `IsEditData` didekati mode lihat (`<fieldset disabled>`),
                preseden `TabPortofolio.tsx`/`TabRetensi.tsx`. */}
            <Field label={FORM_KONTRAK.nomorRujukan} value={rujukan} onChange={setRujukan} readOnly={edmMateri} />
            {/* `pyWidth=0` — mengisi kolomnya; satu-satunya pembedaan lebar
                yang ekspor nyatakan di form ini. */}
            <Area label={FORM_KONTRAK.lingkupWilayah} value={wilayah} onChange={setWilayah} />
              {/* ⛔ Pilihannya NILAI YANG TERSIMPAN, bukan yang layar lama
                tampilkan. Sapuan menemukan domainnya: `reporting` (1.164)
                dan `nonreporting` (687), huruf kecil. Layar lama menulis
                "Reporting"; pemetaan dari yang tersimpan ke yang tampil
                TIDAK ada di ekspor mana pun, jadi ia tidak dikarang. */}
            {/* ⛔ PROPORSIONAL SAJA — `pyCondition
                `TreatyIn.ProportionType='Proportional'` @60649, sel yang
                sama dengan `TreatyIn.Bordeaux` @56974 di
                `Section/TreatyInNONProportional.xml`.

                ⚠️ Dan itu aturan TAMPIL, bukan aturan data: kuncinya ADA di
                772 dari 775 dokumen non-proporsional, bernilai `reporting`
                (711) atau `nonreporting` (61). Layar lama tidak
                memperlihatkannya; layar ini mengikuti layar lama, dan
                nilainya tetap utuh di dokumen. */}
            {jenis !== NON_PROPORSIONAL && (
              <Pilih
                label={FORM_KONTRAK.bordereaux}
                value={bordereaux}
                onChange={setBordereaux}
                opsi={opsi?.bordereaux ?? []}
              />
            )}
            {/* `pxTextArea`, `pyWidth=0` — idem. Baca-saja bila
                `TreatyIn.IsEditData= 1 || TreatyIn.EDMMaterialType = 1`.
                ⛔ `textarea` lokal bermarkup `Area`: `Area` inti nol punya
                `readOnly`, dan ia tidak diubah dari sini. */}
            <div className="field field--lebar">
              <label className="field__label" htmlFor={idBordereauxNote}>
                {FORM_KONTRAK.bordereauxCatatan}
              </label>
              <textarea
                id={idBordereauxNote}
                className="field__input"
                rows={4}
                value={bordereauxNote}
                readOnly={edmMateri}
                onChange={(e) => {
                  setBordereauxNote(e.target.value)
                }}
              />
            </div>
          </div>

          <div className="trin__kolom">
            <TanggalRedup label={FORM_KONTRAK.mulai} value={mulai} onChange={ubahMulai} />
            <TanggalRedup label={FORM_KONTRAK.berakhir} value={berakhir} onChange={setBerakhir} />
            <Field
              label={FORM_KONTRAK.tahunTreaty}
              value={tahunTreaty}
              onChange={() => undefined}
              readOnly
            />
            {/* ⭐ DUA PROPERTI, satu label. Cabangnya yang memilih, dan
                pemilihannya DIBACA: `pyCondition
                `TreatyIn.ProportionType='Proportional'` @111777 untuk
                `AccountingMode` @108019, dan `…='NonProportional'` @118211
                untuk `AccountingModeNonProp` @114500.

                ⛔ Sampai 5 Oktober 2026 layar ini hanya punya yang pertama,
                jadi kontrak non-proporsional memperlihatkan `accounting` /
                `underwriting` — nilai cabang seberang — di tempat `loss` /
                `risk` seharusnya berdiri. */}
            {jenis === NON_PROPORSIONAL ? (
              <Pilih
                label={FORM_KONTRAK.caraPembukuan}
                value={pembukuanNonProp}
                onChange={setPembukuanNonProp}
                opsi={opsi?.caraPembukuanNonProp ?? []}
              />
            ) : (
              <Pilih
                label={FORM_KONTRAK.caraPembukuan}
                value={pembukuan}
                onChange={setPembukuan}
                opsi={opsi?.caraPembukuan ?? []}
              />
            )}
            {/* ⭐ TATA LETAK PEGA — `Inline grid double` L4994: Ceding
                (`Stacked with labels left` L5292) dan RNM as Treaty Leader
                (L5995) BERDAMPINGAN dalam kolom kanan. */}
            <TataPegaBlok tata="g2">
            <div>
            <DropdownWarisan
              label={FORM_KONTRAK.cedant}
              ambil={ambilDaftarCedant}
              nilai={cedant}
              onPilih={(b) => {
                setCedant(b.nama)
                setIdCedant(b.id)
              }}
            />
            {/* ⛔ TERSEMBUNYI, sebab Pega pun tidak menampilkannya: kepala
                hanya merender `TreatyIn.Ceding` (sel 23, baca-saja). Yang
                memperlihatkan pengenal adalah DAFTAR pencariannya — di sini
                `nama — id` untuk nama kembar. */}
            <input type="hidden" name="cedingId" value={idCedant} />
            {/* `TreatyInCheckCedingBlacklist` [2]: pesan pada `TreatyIn.Ceding`. */}
            <PesanMedanAgen pesan={pesanAgen.cedant} />
            </div>
            {/* `pxCheckbox` · `pyCheckboxCaption` "RNM as Treaty Leader" ·
                `pyIncludeLabel=false` — keterangan di SAMPING kotak, tanpa
                label medan. `pyDisabledWhen TreatyIn.ViewState = 1` = mode
                lihat (`<fieldset disabled>`). */}
            <label className="trin__centang">
              <input
                type="checkbox"
                checked={pemimpin}
                onChange={(e) => {
                  setPemimpin(e.target.checked)
                }}
              />
              {FORM_KONTRAK.pemimpinTreaty}
            </label>
            </TataPegaBlok>
            <DropdownWarisan
              label={FORM_KONTRAK.asalBisnis}
              ambil={ambilDaftarAsalBisnis}
              nilai={asalBisnis}
              onPilih={(b) => {
                setAsalBisnis(b.nama)
                setIdAsalBisnis(b.id)
              }}
            />
            {/* Sel 28 `TreatyIn.LeadingReinsSource` baca-saja; pengenalnya
                tidak tampil di Pega. */}
            <input type="hidden" name="leadingReinsSourceId" value={idAsalBisnis} />
            {/* `TreatyInCheckCedingBlacklist` [3]: pesan pada `TreatyIn.LeadingReinsSource`. */}
            <PesanMedanAgen pesan={pesanAgen.asalBisnis} />
          </div>
        </div>
      </Panel>

            {/* ⭐ 8 Oktober 2026 — permintaan pemakai: *"untuk rate of change
          disamakan tablenya seperti lampiran saya ini mengikuti yg ada di
          share"*. Grid berbentuk SAMA dengan grid tab (`GridPega`): tombol
          "Add" di SEL KEPALA kolom tombol, "Delete" per baris di kolom yang
          sama, kosong = satu baris "No items" polos. */}
      <section className="panel trin__kurs">
        <h4 className="panel__title">{FORM_KONTRAK.kurs}</h4>
        <div className="table-wrap">
          {/* ⛔ LEBARNYA TIDAK LAGI DARI EKSPOR — lihat `treatyin.css`.

              Ekspor memberi `.Currency` W=193, `.Value` W=349, `.Date` W=80
              (`Section/InputTreatyInOffer.xml`), yaitu 27 : 48 : 12,5 : 12,5.
              Pemilik proses 7 Oktober 2026 memerintahkan kotak `Value to IDR`
              dikecilkan, jadi rasionya kini 14 : 28 : 21 : 21 — ralat pemilik,
              bukan pergeseran diam-diam.

              ⚠️ HANYA EMPAT `col` untuk LIMA sel di mode ubah: sel tombol
              Remove sengaja nol `col`. Keempatnya berjumlah 84%, jadi sisanya
              jatuh ke sel tombol itu — inilah sebabnya jumlahnya bukan 100%,
              dan mengapa penggulir mendatar yang dulu muncul kini hilang. */}
          <table className="trin__tabel trin__tabel--pega trin__tabel-kurs">
            <colgroup>
              <col className="trin__kol-mata-uang" />
              <col className="trin__kol-nilai" />
              <col className="trin__kol-tanggal" />
              <col className="trin__kol-tanggal" />
            </colgroup>
            <thead>
              <tr>
                <th scope="col">{FORM_KONTRAK.kursMataUang}</th>
                <th scope="col">{FORM_KONTRAK.kursKeIDR}</th>
                <th scope="col">{FORM_KONTRAK.kursBerlakuDari}</th>
                <th scope="col">{FORM_KONTRAK.kursBerlakuSampai}</th>
                {bisaUbah && (
                  <th scope="col">
                    <button
                      type="button"
                      className="btn btn--sm"
                      onClick={() => {
                        // `TreatyInAddCurrency` — `CurrencyList(<APPEND>).CurrencyID = ""`.
                        setKurs([...kurs, { mataUang: '', mataUangID: '', nilaiKeIDR: '', berlakuDari: '', berlakuSampai: '', berlakuDariAsli: '', berlakuSampaiAsli: '' }])
                      }}
                    >
                      {FORM_KONTRAK.tambah}
                    </button>
                  </th>
                )}
              </tr>
            </thead>
            <tbody>
              {kurs.length === 0 && (
                <tr>
                  <td colSpan={bisaUbah ? 5 : 4} className="trin__kosong-pega">
                    {FORM_KONTRAK.tanpaBaris}
                  </td>
                </tr>
              )}
              {bisaUbah &&
                kurs.map((b, i) => (
                  <tr key={i}>
                    {/* ⭐ Sel @439485 `.CurrencyID` — `pxDropdown`, sumber
                        `BrowseCurrency_RD` (nilai `.ID`, teks `.Currency`).
                        Saat berubah Pega menjalankan
                        `SetCurrNameMasterTreaty_Act`: `Obj-Browse CURRENCY`
                        lalu `.Currency = nama`. Dropdown ini mengisi
                        keduanya sekaligus. ⛔ Dahulu kotak teks bebas. */}
                    <td>
                      <DropdownDaftar
                        label=""
                        nilai={b.mataUang}
                        pilihan={mataUangKurs}
                        bisaUbah
                        onPilih={(nama, id) => {
                          ubahKurs(i, { mataUang: nama, mataUangID: id })
                        }}
                      />
                    </td>
                    <td>
                      {/* Sel @454068 `.Conversion` — `pxTextInput`.

                          ⛔ ISIAN PUN DIFORMAT — cacat nyata 6 Oktober 2026:
                          mode Edit merender nilai MENTAH (`10584.39`)
                          sementara mode View merender `10.584,39`. Satu
                          layar, dua bentuk, dan yang menyuntingnya mengira
                          angkanya memang berbeda.

                          ⭐ Diformat HANYA saat sel itu tidak sedang
                          diketik. Memformat di tiap ketukan membuat koma
                          desimal mustahil diketik: `10,` berubah menjadi
                          `10` sebelum angka berikutnya sempat masuk. */}
                      <input
                        className="field__input"
                        type="text"
                        value={selDiketik === `${String(i)}:nilaiKeIDR` ? b.nilaiKeIDR : selAngka(['uang', 2], b.nilaiKeIDR)}
                        aria-label={FORM_KONTRAK.kursKeIDR}
                        onFocus={() => {
                          setSelDiketik(`${String(i)}:nilaiKeIDR`)
                        }}
                        onBlur={() => {
                          setSelDiketik(null)
                        }}
                        onChange={(e) => {
                          ubahKurs(i, { nilaiKeIDR: saringAngka(e.target.value) })
                        }}
                      />
                    </td>
                    {/* Sel @460497 `.PeriodStart` / @466983 `.PeriodEnd` —
                        `pxDateTime`: kotak tanggal yang dapat diketik atau
                        dipilih dari kalender, bukan teks bebas.

                        ⛔ RALAT 7 Oktober 2026 — kotak ini SEMPAT KOSONG
                        walau basis datanya berisi. Ia diberi `berlakuDari`,
                        bentuk TAMPIL `dd/mm/yy`; kotak tanggal menerima
                        bentuk kabel `DD-MM-YYYY`, dan tahun dua digit
                        beserta garis miring itu tidak terbaca olehnya.
                        Sekarang yang diberikan `…Asli` (`YYYYMMDD`).

                        ⚠️ Kosongnya BERBAHAYA, bukan sekadar jelek: Save
                        menulis isi kotak ini apa adanya ke `STARTDATE`,
                        jadi membuka lalu menyimpan kontrak akan MENGHAPUS
                        tanggal yang sudah ada. */}
                    <td>
                      <KotakTanggalKetik
                        label={FORM_KONTRAK.kursBerlakuDari}
                        value={b.berlakuDariAsli}
                        onChange={(v) => {
                          ubahKurs(i, { berlakuDariAsli: v })
                        }}
                      />
                    </td>
                    <td>
                      <KotakTanggalKetik
                        label={FORM_KONTRAK.kursBerlakuSampai}
                        value={b.berlakuSampaiAsli}
                        onChange={(v) => {
                          ubahKurs(i, { berlakuSampaiAsli: v })
                        }}
                      />
                    </td>
                    <td>
                      <button
                        type="button"
                        className="btn btn--sm"
                        onClick={() => {
                          setKurs(kurs.filter((_, j) => j !== i))
                        }}
                      >
                        {GRID_TAMBAH.hapus}
                      </button>
                    </td>
                  </tr>
                ))}
              {!bisaUbah && kurs.map((b, i) => (
                <tr key={i}>
                  <td>{b.mataUang}</td>
                  {/* ⭐ DUA desimal, nol di ekor DIPERTAHANKAN — gambar 01
                      berbunyi `1,00` dan `16.000,00`, gambar 26 `15.500,00`.
                      Aturan lama membuangnya dan memberi `1`. */}
                  <td>{selAngka(['uang', 2], b.nilaiKeIDR)}</td>
                  {/* `pxDateTime` — tanggal, bukan stempel mentah
                      `20250701T140000.000 GMT`. Dibaca dari `…Asli`
                      (`YYYYMMDD`): `berlakuDari` berbentuk `dd/mm/yy`, dan
                      `formatDate` membacanya bulan-dulu — tanggal tertukar
                      atau kosong (laporan 8 Oktober 2026). */}
                  <td>{formatDate(b.berlakuDariAsli)}</td>
                  <td>{formatDate(b.berlakuSampaiAsli)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* Keterangan kaki, satu baris redup — bukan paragraf di tengah layar. */}

      </fieldset>

      {/* Strip tab — himpunannya berganti bersama radio di atas. */}
      <StripTab tab={tab} aktif={tabTampil} onPilih={setTabAktif} />

      {/* Isi tab — mode lihat mematikannya seperti kepala. `key` melahirkan
          ulang grid ketika data kontrak tiba, sebab grid menyalin barisnya
          SEKALI. */}
      <PenyediaHalaman penampung={penampung}>
      {/* ⭐ Kontrak REVISI: tab Information & Submit tidak dimatikan — Comment
          dan Submit revisinya hidup di mode lihat (`TreatyInfoSubmit` cell 9
          dan 21); tab itu sendiri yang mematikan medan lainnya. */}
      {/* ⭐ `trin__isi-tab` — seluruh isi tab berbentuk PEGA DATAR
          (`treatyin.css` "BENTUK PEGA — SELURUH ISI TAB", 8 Oktober 2026). */}
      <fieldset
        className="trin__mode trin__isi-tab"
        disabled={!bisaUbah && !(revisi && tabTampil === TAB_REVISI)}
        key={`${idKontrak}|${warisan === null ? '-' : 'isi'}|${mode}|${generasi}`}
      >

      {/* ⭐ DUA BELAS TAB KINI BERISI. Tujuh dari tabel pendaratan
          `M_TREATYIN_*`, empat dari `M_TREATY_IN2` (satu tabel warisan,
          empat proyeksi atas baris layer yang sama), dan Rate of Exchange
          dari dokumen. Sisanya tetap menyatakan "belum ada KODE" — dan
          bedanya dijaga: tab berisi yang kosong memakai `Kosong` ("belum
          ada DATA"), tab tanpa kode memakai `.trin__belum`. */}
      {tabTampil === 'Reporting Period' ? (
        <TabReportingPeriod baris={warisan?.periodePelaporan ?? []} opsiPeriode={opsi?.periodePelaporan ?? []} mode={mode} />
      ) : tabTampil === 'Portfolio' ? (
        // ⛔ BUKAN `TabGridWarisan` lagi, sejak 6 Oktober 2026. Ekspor
        // memberi tab ini dua daftar pilihan dan satu area teks, dan grid
        // umum hanya tahu kotak teks. Pemetaan kolomnya juga DIBETULKAN di
        // langkah yang sama: kolom 1 `TypePortfolio`, kolom 2 `Type` —
        // sebelumnya tertukar. Lihat `TabPortofolio.tsx`.
        <TabPortofolio
          baris={warisan?.portofolio ?? []}
          petunjukKosong={FORM_KONTRAK.petunjukPortofolio}
          mode={mode}
        />
      ) : tabTampil === 'EGNPI' ? (
        /* ⛔ GANTI 6 Oktober 2026 — grid hanya-baca menjadi tab berumus.
           Ekspor (`Section/TreatyInTabsNonProportional.xml`, wadah TABBED
           ke-3) memberi tab ini `Add`/`Delete`, dua panel total, dan dua
           tombol berumus (`Update Total`, `Update EGNPI Value`). Nol di
           antaranya ada sebelum ini, sehingga `Amount in IDR` dan
           `Proportion %` nol pernah terisi — padahal tab Limits MEMBACA
           keduanya lewat `TotalEgnpi`.

           ⚠️ EGNPI hanya ada di `TAB_NON_PROPORSIONAL`; cabang prop nol
           pernah sampai ke sini, jadi nol percabangan jenis di bawah. */
        <TabEgnpi
          baris={egnpiAwal}
          kurs={kurs.map((k) => ({ Currency: k.mataUang, Conversion: k.nilaiKeIDR }))}
          // ⭐ Retention TERKINI (tab Maximum Retention, penampung halaman).
          retensi={retensiKini.map((r) => ({ Currency: r.Currency, CurrencyID: r.CurrencyID }))}
          mode={mode}
        />
      ) : tabTampil === 'Maximum Retention' ? (
        /* ⛔ GANTI 6 Oktober 2026 — grid hanya-baca + panel bertombol MATI
           menjadi tab berumus. Ekspor (`Section/TreatyInTabsNonProportional
           .xml`, wadah TABBED ke-1) memberi tab ini `Add`/`Delete`, rincian
           baris (`Section/MaxRetention.xml`), dan `Update Total`
           (`TreatyInNPSetTotal` type=retention).

           ⚠️ SATU tombol, bukan dua seperti EGNPI: retensi nol punya
           kolom `Amount in IDR`, jadi nol konversi kurs untuk dijalankan. */
        <TabRetensi
          baris={retensiAwal}
          totalAwal={(warisan?.totalRetensi ?? []).map((b) => ({
            Currency: b.mataUang,
            CurrencyID: '',
            Value: b.nilai,
          }))}
          edmJenisMaterial={warisan?.edmJenisMaterial ?? ''}
          mode={mode}
        />
      ) : tabTampil === 'Installment' ? (
        /* ⭐ 7 Oktober 2026 — tab BERUMUS dari ekspor (`TabAngsuran.tsx`):
           Installment + Update Value (`TreatyInSetValueInstallment`), grid per
           mata uang dengan rincian `Installments`, Update Total. Nilainya
           `TotalShareNetNP` tab Share TERKINI.

           ⛔ Add/Delete tetap TIDAK ada — Pega tidak punya keduanya (nol sel
           tombol bergrid `TreatyIn.Installment` / `.InstallmentList` di 52
           seksi); barisnya LAHIR dari Update Value. */
        <TabAngsuran
          baris={warisan?.angsuran ?? []}
          netPremium={netPremiumKini}
          edmState={warisan?.edmState ?? ''}
          edmJenisMaterial={warisan?.edmJenisMaterial ?? ''}
          mode={mode}
          // Kepala TERKINI — Due Date dari Commencement, dibagi rata.
          commencement={keSimpan(mulai) || (warisan?.tanggalMulaiAsli ?? '')}
          termination={keSimpan(berakhir) || (warisan?.tanggalBerakhirAsli ?? '')}
        />
      ) : tabTampil === 'Information & Submit' ? (
        /* ⛔ RALAT 6 Oktober 2026 — tab ini FORM, bukan grid riwayat.
           `Section/TreatyInfoSubmit.xml`: dua `Text area`
           (`TreatyIn.Information`, `TreatyIn.Comment`) dan dua tombol.

           ⚠️ Yang dirender sebelumnya bukan sekadar salah, ia SALINAN —
           riwayat sudah punya panelnya sendiri di kaki layar, membaca tabel
           yang sama. Layar menampilkan daftar yang sama dua kali, dan tab
           yang seharusnya tempat MENGIRIM justru tempat membaca. */
        <TabInfoSubmit
          mode={mode}
          statusAkseptasi={warisan?.statusAkseptasi ?? ''}
          revisi={revisi}
          sibuk={sibukTulis}
          onKirim={() => {
            tekanKirim('submit')
          }}
          onTolak={() => {
            tekanKirim('decline')
          }}
        />
      ) : tabTampil === 'Limits' ? (
        jenis === NON_PROPORSIONAL ? (
          /* ⭐ Cabang NP: grid layer → `Layers` → Summary → Total All
             Layers, dari ekspor — `labelsLimitsNP.ts`. EGNPI dan kurs yang
             rumusnya baca diambil dari tab EGNPI dan grid Rate of Exchange
             kepala (keadaan TERKINI grid itu, seperti clipboard Pega). */
          <TabLimitsNonProp
            petunjukKosong={FORM_KONTRAK.petunjukLayer}
            pohon={limitsNP?.layers ?? warisan?.limitsPohon ?? []}
            akar={limitsNP?.akar ?? warisan?.limitsAkar}
            onUbah={(layers, akar) => setLimitsNP({ layers, akar })}
            // ⭐ EGNPI TERKINI (tab EGNPI, penampung halaman) — `TotalEgnpi`.
            egnpi={egnpiKini.map((e) => ({
              TreatyGroup: e.TreatyGroup,
              Currency: e.Currency,
              CurrencyID: e.CurrencyID,
              Amount: e.Amount,
            }))}
            kurs={kurs.map((k) => ({ Currency: k.mataUang, Conversion: k.nilaiKeIDR }))}
            edmState={warisan?.edmState ?? ''}
            edmJenisMaterial={warisan?.edmJenisMaterial ?? ''}
            mode={mode}
          />
        ) : (
          /* ⭐ Cabang P: tiga tingkat dari ekspor — `labelsLimitsProp.ts`. */
          <TabLimitsProp
            pohon={warisan?.limitsPohon ?? []}
            mode={mode}
            idKontrak={idKontrak}
            edmJenisMaterial={warisan?.edmJenisMaterial ?? ''}
          />
        )
      ) : tabTampil === 'Share' ? (
        /* ⭐ `RNM Share` adalah SUB-TAB di dalam `Share`, bukan tab setara —
           dan itu kini terbukti dari dokumen desain, bukan disimpulkan.

           Gambar `16`/`17` (prop) dan `34` (non-prop) memperlihatkan tab
           `Share` berisi panel atasnya, lalu satu strip sub-tab berjudul
           `RNM Share` di bawahnya. Pengukuran ekspor ronde sebelumnya
           mengatakan hal yang sama dari sisi lain: kesebelas wadah `TABBED`
           di `Section/TreatyInTabsNonProportional.xml` tidak memuat
           `RNM Share`, dan kedua judulnya (@2.093.710, @2.133.893) jatuh di
           dalam wilayah tab `Share`.

           ⛔ Ini BUKAN penghapusan tab — §6 melarangnya, dan nol tab
           dihapus: `RNM Share` tetap dirender, dengan kolom dan data yang
           sama, satu tingkat di dalam tempat ekspor dan gambar menaruhnya. */
        /* ⭐ DUA BENTUK, DUA CABANG — tangkapan layar Pega 6 Oktober 2026.
           Cabang proporsional memakai panel `Total Share` (Refresh, % RNM
           Share, % Brokerage, Option) lalu sub-tab `RNM Share` berisi grid
           `Kind of Treaty` dan tiga grid total. Cabang non-proporsional
           memakai grid per layer, dan kolomnya tidak pernah ada di sini. */
        jenis === PROPORSIONAL ? (
          <TabShareProp
            pohon={warisan?.limitsPohon ?? []}
            petunjukKosong={FORM_KONTRAK.petunjukLayer}
            mode={mode}
            // Commencement kepala TERKINI — saringan RD spreading.
            commencement={keSimpan(mulai) || (warisan?.tanggalMulaiAsli ?? '')}
            edmJenisMaterial={warisan?.edmJenisMaterial ?? ''}
          />
        ) : (
          /* ⭐ Cabang NON-PROP — tangkapan layar Pega pemakai 7 Oktober 2026
             dan `TreatyInTabsNonProportional.xml` @1695720: panel Share,
             grid Reinsurer / Facultative Reinsurers, sub-tab RNM Share
             (grid per layer + rincian), Summarry, Total All Layers. Isi dari
             `warisan.shareNP` (pendaratan `T_TREATY_SHARE*`); rumus di
             services. Lihat `TabShareNonProp.tsx`. */
          <TabShareNonProp
            key={`${idKontrak}|${warisan === null ? '-' : 'isi'}`}
            share={shareNP ?? warisan?.shareNP}
            layers={limitsNP?.layers ?? warisan?.limitsPohon ?? []}
            onUbah={setShareNP}
            mode={mode}
            idKontrak={idKontrak}
            // Commencement medan kepala TERKINI (yang sedang diisi).
            commencement={keSimpan(mulai) || (warisan?.tanggalMulaiAsli ?? '')}
            edmState={warisan?.edmState ?? ''}
            edmJenisMaterial={warisan?.edmJenisMaterial ?? ''}
          />
        )
      ) : tabTampil === 'Event Limits' ? (
        /* ⭐ EMPAT BARIS BERLABEL, SATU SET PER KONTRAK — gambar 28.
           ⛔ Tidak lagi dari `warisan.layer`: tab Non-Prop mengikat properti
           AKAR (`TreatyIn.RSMDLimit` …), bukan `Detail[]`. Lihat
           `TabEventLimits.tsx`. */
        <TabEventLimits mode={mode} />
      ) : tabTampil === 'Achievement In IDR' ? (
        /* ⛔ Cabang tab utama `RNM Share` (Non-Prop) DICABUT 8 Oktober 2026 —
           pemilik proses: *"di non prop tab RNM SHARE itu tidak ada"*. Ia
           hanya sub-tab di dalam tab `Share` (`TabShareNonProp`). */
        /* ⭐ TAB BARU 7 Oktober 2026. Namanya sudah ada di
           `TAB_PROPORSIONAL` sejak lama, tetapi NOL cabang merendernya —
           membukanya menampilkan isi tab pertama.

           Bentuknya dari `Section/AchievementCombine.xml`: grid enam kolom
           hanya-baca + tiga sel kaki. Angkanya dari
           `Activity/GetAchievement.xml`, rumus yang SAMA yang sudah dipakai
           sub-tab Achievement di dalam Limits Prop. Lihat
           `TabAchievement.tsx`. */
        <TabAchievement idKontrak={idKontrak} pohon={warisan?.limitsPohon ?? []} />
      ) : tabTampil === 'Co-Ins Scale' ? (
        // ⭐ Grid `TreatyIn.CoInScale` (dua kolom, Add/Delete hidup sel
        // 262/266) DITAMBAH dua medan `Max Co-Insurance Panel` sel 277/278
        // — urutan gambar 18. Lihat `TabCoInsScale.tsx`.
        <TabCoInsScale baris={warisan?.skalaKoasuransi ?? []} mode={mode} />
      ) : tabTampil === 'Exclusions' ? (
        <TabTeksPanjang
          judul={tabTampil}
          // `TreatyIn.ExclusionsP` (Prop) / `TreatyIn.Exclusions` (Non-Prop).
          properti={jenis === NON_PROPORSIONAL ? 'Exclusions' : 'ExclusionsP'}
          tab={warisan?.pengecualian}
          petunjukKosong={FORM_KONTRAK.petunjukTeksPengecualian}
          mode={mode}
        />
      ) : tabTampil === 'Special Conditions' ? (
        <TabTeksPanjang
          judul={tabTampil}
          properti={jenis === NON_PROPORSIONAL ? 'SpecialConditions' : 'SpecialConditionsP'}
          tab={warisan?.syaratKhusus}
          petunjukKosong={FORM_KONTRAK.petunjukTeksSyarat}
          mode={mode}
        />
      ) : tabTampil === 'Accumulation' ? (
        /* ⭐ 7 Oktober 2026 — tab BERUMUS, bukan grid umum: Period →
           `TreatyInSetAccountReport` (membaca Start/End tab Reporting Period
           dari penampung halaman), Reporting Date / Submission Days →
           `TreatyInAccumulationSetSubDue`, Add → `TreatyInAddAccumulation`.
           Lihat `TabAkumulasi.tsx`. */
        <TabAkumulasi baris={warisan?.akumulasi ?? []} mode={mode} edmJenisMaterial={warisan?.edmJenisMaterial ?? ''} />
      ) : (
        /* ⛔ SENGAJA BUKAN `Kosong`. `Kosong` menyatakan *"belum ada DATA"*;
           tab ini menyatakan *"belum ada KODE"*. Keduanya bukan hal yang
           sama, dan yang salah membacanya akan menagih hal yang salah —
           menagih pemuatan data padahal yang kurang layarnya. Bentuknya
           dibedakan (tepi putus-putus, nol ikon keranjang) supaya bedanya
           terlihat sebelum teksnya dibaca. */
        <Panel judul={tabTampil}>
          {/* ⛔ KEADAAN KEEMPAT untuk Retro: JARANG, bukan "belum ada kode".
              Keduanya memakai kotak bertepi putus-putus yang sama — yang
              berbeda teksnya, sebab yang berbeda memang alasannya, bukan
              bentuknya. `KEPUTUSAN §17`. */}
          <div className="trin__belum" role="note">
            <span className="trin__belum-judul">
              {FORM_KONTRAK.belumDibangun}
            </span>
            <span className="trin__belum-petunjuk">
              {FORM_KONTRAK.belumDibangunPetunjuk}
            </span>
          </div>
        </Panel>
      )}
      </fieldset>
      </PenyediaHalaman>
    </div>
    {/* ⭐ WADAH KEDUA — permintaan pemakai 8 Oktober 2026 (gambar Pega:
        Attachment berdiri di wadah putih TERSENDIRI di bawah wadah tab):
        *"Attachment itu dipisahkan dulu dynamic layout nya jangan gabung
        dengan container inputan, krn menjadi makan tempat"*. */}
    <div className="inbox trin__inbox-lampiran">
      {/* ⭐ ATTACHMENT · tombol · HISTORY — di BAWAH strip tab, bukan di
          dalamnya. Begitu layar lama menyusunnya: panel Attachment, deret
          tombol Save/Close/Actions, lalu panel History. Ketiganya berlaku
          untuk kontraknya, bukan untuk tab yang kebetulan terbuka. */}
      <PanelLampiran
        kategori={warisan?.kategoriLampiran ?? []}
        berkas={warisan?.lampiran ?? []}
        idKontrak={idKontrak}
        bisaUnggah={bisaUbah || revisi}
        statusAkseptasi={statusKini}
      />

      {/* ⭐ HIDUP sejak 7 Oktober 2026 — keputusan pemilik proses: Save/Submit
          menyimpan ke tabel masing-masing (`T_TREATY_*`, kepala `TREATY_IN`,
          kurs `TREATYEXCHANGEYEARLY`), tidak lewat Pega.

            Save     `IsEditData !='1' && StatusAkseptasi != 'Resolve Complete'`
                     → DT `TreatyInAddNew` → `SaveTreatyIn_Act`
            Close    selalu
            Actions  pemegang workbasket `TreatyIn.Position` penyetuju
                     (SecHead/DeptHead/Director) → modal `TreatyInAction` →
                     `TreatyInAkseptasi_Act`

          ⛔ Tombol yang syarat tampilnya tidak terpenuhi TIDAK dirender —
          persis Pega. */}
      <div className="trin__aksi trin__aksi--kaki" role="group" aria-label={FORM_KONTRAK.judul}>
        {saveTampil && (
          <button type="button" className="btn btn--primary" disabled={sibukTulis} onClick={tekanSave}>
            {sibukTulis ? TOMBOL_TULIS.menyimpan : TOMBOL_TULIS.simpan}
          </button>
        )}
        <button type="button" className="btn" onClick={onKembali}>
          {TOMBOL_TULIS.tutup}
        </button>
        {paksaTampil && (
          <button
            type="button"
            className="btn"
            title={TOMBOL_TULIS.petunjukPaksaEdit}
            onClick={() => {
              setPaksaUbah(true)
            }}
          >
            {TOMBOL_TULIS.paksaEdit}
          </button>
        )}
        {actionsTampil && (
          <button
            type="button"
            className="btn"
            disabled={sibukTulis}
            onClick={() => {
              setPilihanActions('Accept')
              setKomentarActions('')
              setActionsBuka(true)
            }}
          >
            {TOMBOL_TULIS.aksi}
          </button>
        )}
      </div>
      {hasilTulis !== null && (
        <div className={hasilTulis.galat ? 'alert alert--error' : 'alert alert--info'} role={hasilTulis.galat ? 'alert' : 'status'}>
          {hasilTulis.pesan}
          {hasilTulis.takTersimpan.length > 0 && (
            <div className="trin__redup">
              {TOMBOL_TULIS.takTersimpan} {hasilTulis.takTersimpan.join(', ')}
            </div>
          )}
        </div>
      )}
      {actionsBuka && (
        <Modal
          judul={TOMBOL_TULIS.judulAksi}
          onTutup={() => {
            setActionsBuka(false)
          }}
          labelBatal={TOMBOL_TULIS.batal}
          onKirim={() => {
            setActionsBuka(false)
            tekanKirim('akseptasi', pilihanActions, { Comment: komentarActions })
          }}
          aksi={
            <button type="submit" className="btn btn--primary" disabled={sibukTulis}>
              {TOMBOL_TULIS.kirim}
            </button>
          }
        >
          <p>
            {TOMBOL_TULIS.dariID} <strong>{idKontrak}</strong>
          </p>
          <div className="trin__radio" role="radiogroup" aria-label={TOMBOL_TULIS.pilihan}>
            {PILIHAN_AKSEPTASI.map((v) => (
              <label key={v}>
                <input
                  type="radio"
                  name="trin-pilihan-akseptasi"
                  value={v}
                  checked={pilihanActions === v}
                  onChange={() => {
                    setPilihanActions(v)
                  }}
                />
                {v}
              </label>
            ))}
          </div>
          <Area label={TOMBOL_TULIS.komentar} value={komentarActions} onChange={setKomentarActions} baris={4} />
        </Modal>
      )}

      <PanelHistory baris={warisan?.catatan ?? []} />

    </div>
    </>
  )
}

/** Dipakai uji — menjaga cabang `Gagal`/`Memuat` tetap ada dan tidak dipakai salah. */
export const cabangKeadaan = { Gagal, Memuat }
