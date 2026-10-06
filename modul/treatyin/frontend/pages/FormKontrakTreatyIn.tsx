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

import {  useEffect, useState } from 'react'

import {
  Area,
  Field,
  Gagal,
  Kosong,
  Memuat,
  Panel,
  Pilih,
  StripTab,
} from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilDaftarAsalBisnis,
  ambilDaftarCedant,
  ambilKontrakWarisan,
  ambilOpsiKepala,
  type KontrakWarisan,
  type OpsiKepala,
} from '../api'
import type { ModeForm } from '../mode'
import {
  FORM_KONTRAK,
  SYARAT_TAB_NON_PROPORSIONAL,
  SYARAT_TAB_PROPORSIONAL,
  type SyaratTabKontrak,
  TAB_NON_PROPORSIONAL,
  TAB_PROPORSIONAL,
  KOLOM_AKUMULASI,
  KOLOM_EGNPI,
  KOLOM_RETENSI,
  KOLOM_ANGSURAN,
  KOLOM_CATATAN,
  KOLOM_COIN_SCALE,
  JENIS_COIN_SCALE,
  JENIS_EGNPI,
  JENIS_RETENSI,
  JENIS_ANGSURAN,
  GRID_TAMBAH,
} from '../labels'

// ⭐ KOMPONEN DIPECAH 5 Oktober 2026 — pemindahan MURNI, nol perubahan
// perilaku. Contohnya `modul/masterproductnamelife`, yang komponen
// terbesarnya 372 baris sementara layar ini dahulu 1.826 dalam satu berkas.
import MedanTakAda from '../components/MedanTakAda'
import PanelHistory from '../components/PanelHistory'
import PanelPolisProduksi from '../components/PanelPolisProduksi'
import PanelTotalRetensi from '../components/PanelTotalRetensi'
import TanggalRedup from '../components/TanggalRedup'
import { barisAngka, selAngka } from '../components/angka'
import PanelLampiran from '../components/PanelLampiran'
import PohonLimits from '../components/PohonLimits'
import TabLimitsProp from '../components/TabLimitsProp'
import SubTabShare, { PanelRnmShare } from '../components/TabShare'
import TabGridWarisan from '../components/TabGridWarisan'
import TabPortofolio from '../components/TabPortofolio'
import TabTeksPanjang from '../components/TabTeksPanjang'
import TabEventLimits from '../components/TabEventLimits'
import DropdownWarisan from '../components/DropdownWarisan'
import TabReportingPeriod from '../components/TabReportingPeriod'

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
  nilaiKeIDR: string
  berlakuDari: string
  berlakuSampai: string
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
}

export default function FormKontrakTreatyIn({ idKontrak, mode = 'lihat', onKembali }: FormKontrakProps) {
  const bisaUbah = mode === 'ubah'
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
  const [pembukuan, setPembukuan] = useState('')
  // ⛔ MEDAN KEDUA, bukan medan yang sama. `AccountingModeNonProp` berdiri
  // sendiri di ekspor dengan `pyCondition` cabangnya sendiri; satu keadaan
  // untuk keduanya akan membuat berpindah cabang menimpa nilai cabang yang
  // ditinggalkan.
  const [pembukuanNonProp, setPembukuanNonProp] = useState('')
  const [cedant, setCedant] = useState('')
  const [asalBisnis, setAsalBisnis] = useState('')
  const [pemimpin, setPemimpin] = useState(false)
  // ⭐ Grid Rate of Exchange dibaca dari dokumen warisan kontrak ini, bukan
  // dari keadaan kosong. `CurrencyList` berisi di 297 dari 300 dokumen yang
  // disapu — "No items" selama ini bukan karena datanya tidak ada, melainkan
  // karena tidak ada yang membacanya.
  const [kurs, setKurs] = useState<BarisKurs[]>([])
  // Sel grid kurs yang sedang diketik (`baris:kolom`), atau null. Dipakai
  // supaya pemformat tidak melawan pengetik — lihat catatan di gridnya.
  const [selDiketik, setSelDiketik] = useState<string | null>(null)

  // ⛔ Mengisi SELURUH medan dari kontrak nyata, termasuk radio jenisnya.
  // Radio itu memilih strip tab, dan tiket ronde ini menyebutnya: ia
  // disambungkan ke NILAI NYATA, bukan dibiarkan pada bawaan yang pemakai
  // ubah sendiri — kontrak non-proporsional yang membuka strip proporsional
  // memperlihatkan sebelas tab yang tidak satu pun miliknya.
  useEffect(() => {
    if (idKontrak === '') {
      setMemuat(false)
      return
    }
    let dibuang = false
    setMemuat(true)
    setGalat(null)
    ambilKontrakWarisan(idKontrak)
      .then((k) => {
        if (dibuang) return
        setWarisan(k)
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
        setMulai(k.tanggalMulai)
        setBerakhir(k.tanggalBerakhir)
        setPembukuan(k.caraPembukuanAsli)
        setPembukuanNonProp(k.caraPembukuanNonPropAsli)
        setCedant(k.cedant)
        setAsalBisnis(k.asalBisnis)
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
  }, [idKontrak])

  /** Kunci itu ADA di dokumen kontrak ini? Kontrak baru: tidak relevan. */
  const adaKunci = (kunci: string) => warisan === null || (warisan.adaDiJson[kunci] ?? false)

  // ⭐ Syarat tampil tab — dari dokumen kontrak, nol untuk kontrak baru.
  const syaratTab: SyaratTabKontrak | undefined =
    warisan === null
      ? undefined
      : {
          retroBerganda: warisan.retroBerganda,
          edmState: warisan.edmState,
          edmJenisMaterial: warisan.edmJenisMaterial,
        }
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
  const tahunTreaty = warisan?.tahunTreaty ?? ''

  return (
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
        <div className="trin__kepala">
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
            bawah membuat susunannya DAPAT DIUJI tanpa merender. */}
        <div className="trin__dwikolom">
          <div className="trin__kolom">
            <Field label={FORM_KONTRAK.namaKontrak} value={nama} onChange={setNama} />
            {adaKunci('ContractRefNo') ? (
              <Field label={FORM_KONTRAK.nomorRujukan} value={rujukan} onChange={setRujukan} />
            ) : (
              <MedanTakAda label={FORM_KONTRAK.nomorRujukan} />
            )}
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
            {/* `pyWidth=0` — idem. */}
            {adaKunci('BordereauxNote') ? (
              <Area
                label={FORM_KONTRAK.bordereauxCatatan}
                value={bordereauxNote}
                onChange={setBordereauxNote}
              />
            ) : (
              <MedanTakAda label={FORM_KONTRAK.bordereauxCatatan} />
            )}
          </div>

          <div className="trin__kolom">
            <TanggalRedup label={FORM_KONTRAK.mulai} value={mulai} onChange={setMulai} />
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
            <DropdownWarisan
              label={FORM_KONTRAK.cedant}
              ambil={ambilDaftarCedant}
              nilai={cedant}
              onPilih={(b) => {
                setCedant(b.nama)
              }}
            />
            {adaKunci('TreatyLeader') ? (
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
            ) : (
              <MedanTakAda label={FORM_KONTRAK.pemimpinTreaty} />
            )}
            <DropdownWarisan
              label={FORM_KONTRAK.asalBisnis}
              ambil={ambilDaftarAsalBisnis}
              nilai={asalBisnis}
              onPilih={(b) => {
                setAsalBisnis(b.nama)
              }}
            />
          </div>
        </div>
      </Panel>

            {/* ⛔ TOMBOL "Add" DI BARIS KEPALA, bukan di bawah kotak kosong.
          Di ekspor ia duduk sebaris dengan kepala grid; menaruhnya di bawah
          area kosong membuat panel menjulur tinggi dan tombolnya terbaca
          sebagai milik keadaan kosong, bukan milik gridnya. */}
      <section className="panel">
        <div className="trin__panel-kepala">
          <h4 className="panel__title">{FORM_KONTRAK.kurs}</h4>
          {bisaUbah && (
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              onClick={() => {
                setKurs([...kurs, { mataUang: '', nilaiKeIDR: '', berlakuDari: '', berlakuSampai: '' }])
              }}
            >
              {FORM_KONTRAK.tambah}
            </button>
          )}
        </div>
        <div className="table-wrap">
          {/* ⛔ Perbandingan lebar DARI ekspor, satu-satunya di layar ini yang
              sungguh berbeda: `.Currency` W=193, `.Value` W=349, `.Date` W=80
              (`Section/InputTreatyInOffer.xml`, sel ber-`pyValue` + `pyWidth`).
              193 : 349 : 80 : 80 -> 27% : 48% : 12,5% : 12,5%. Persen, bukan
              piksel: tata letak kita responsif dan Pega tidak. */}
          <table className="trin__tabel trin__tabel-kurs">
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
              </tr>
            </thead>
            <tbody>
              {kurs.length === 0 && (
                <tr>
                  <td colSpan={4}>
                    <Kosong pesan={FORM_KONTRAK.tanpaBaris} petunjuk={FORM_KONTRAK.kursPetunjuk} />
                  </td>
                </tr>
              )}
              {bisaUbah &&
                kurs.map((b, i) => (
                  <tr key={i}>
                    {(['mataUang', 'nilaiKeIDR', 'berlakuDari', 'berlakuSampai'] as const).map((kk) => (
                      <td key={kk}>
                        {/* ⛔ ISIAN PUN DIFORMAT — cacat nyata 6 Oktober 2026:
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
                          value={
                            selDiketik === `${String(i)}:${kk}` || kk !== 'nilaiKeIDR'
                              ? b[kk]
                              : selAngka(['uang', 2], b[kk])
                          }
                          aria-label={kk}
                          onFocus={() => {
                            setSelDiketik(`${String(i)}:${kk}`)
                          }}
                          onBlur={() => {
                            setSelDiketik(null)
                          }}
                          onChange={(e) => {
                            setKurs(kurs.map((x, j) => (j === i ? { ...x, [kk]: e.target.value } : x)))
                          }}
                        />
                      </td>
                    ))}
                    <td>
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
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
                  <td>{b.berlakuDari}</td>
                  <td>{b.berlakuSampai}</td>
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
      <fieldset className="trin__mode" disabled={!bisaUbah} key={`${idKontrak}|${warisan === null ? '-' : 'isi'}|${mode}`}>

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
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_EGNPI}
          baris={barisAngka(
            // ⛔ URUTANNYA mengikuti `KOLOM_EGNPI`, yang kini urut layar
            // (gambar 29), bukan urut abjad kunci dokumen.
            (warisan?.egnpi ?? []).map((b) => [
              b.kelompokTreaty, b.perTanggal, b.proporsi, b.mataUang,
              b.jumlah, b.jumlahIDR, b.kelasBisnis, b.keterangan,
            ]),
            JENIS_EGNPI,
          )}
          petunjukKosong={FORM_KONTRAK.petunjukTabel}
        />
      ) : tabTampil === 'Maximum Retention' ? (
        <>
          <TabGridWarisan
            judul={tabTampil}
            kolom={KOLOM_RETENSI}
            baris={barisAngka(
              // ⛔ Urut layar (gambar 26), bukan abjad.
              (warisan?.retensi ?? []).map((b) => [
                b.kelompokTreaty, b.mataUang, b.jumlah, b.kelasBisnis, b.keterangan,
              ]),
              JENIS_RETENSI,
            )}
            petunjukKosong={FORM_KONTRAK.petunjukTabel}
            mode={mode}
          />
          <PanelTotalRetensi baris={warisan?.totalRetensi ?? []} />
        </>
      ) : tabTampil === 'Installment' ? (
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_ANGSURAN}
          baris={barisAngka(
            (warisan?.angsuran ?? []).map((b) => [
              b.angsuran, b.mataUang, b.jumlah, b.persen,
              b.jatuhTempo, b.tanggalBayar, b.wpc,
            ]),
            JENIS_ANGSURAN,
          )}
          petunjukKosong={FORM_KONTRAK.petunjukTabel}
          mode={mode}
        />
      ) : tabTampil === 'Information & Submit' ? (
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_CATATAN}
          baris={(warisan?.catatan ?? []).map((b) => [
            b.tanggal, b.operator, b.disetujui, b.catatan,
          ])}
          petunjukKosong={FORM_KONTRAK.petunjukTabel}
          mode={mode}
          bisaTambah={false}
        />
      ) : tabTampil === 'Limits' ? (
        jenis === NON_PROPORSIONAL ? (
          <PohonLimits layer={warisan?.layer ?? []} nonProp />
        ) : (
          /* ⭐ Cabang P: tiga tingkat dari ekspor — `labelsLimitsProp.ts`. */
          <TabLimitsProp pohon={warisan?.limitsPohon ?? []} mode={mode} />
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
        <SubTabShare
          layer={warisan?.layer ?? []}
          petunjukKosong={FORM_KONTRAK.petunjukLayer}
          mode={mode}
        />
      ) : tabTampil === 'Event Limits' ? (
        /* ⭐ EMPAT BARIS BERLABEL, bukan grid sembilan kolom — gambar 28. */
        <TabEventLimits layer={warisan?.layer ?? []} />
      ) : tabTampil === 'RNM Share' ? (
        /* ⚠️ Cabang ini TETAP ADA walau `RNM Share` kini dirender sebagai
           sub-tab `Share`. Sebabnya: daftar tab masih memuatnya — §6
           melarang menghapus tab mana pun — jadi strip utama masih dapat
           memilihnya, dan tab yang dapat dipilih harus merender sesuatu. */
        <PanelRnmShare
          layer={warisan?.layer ?? []}
          petunjukKosong={FORM_KONTRAK.petunjukLayer}
          mode={mode}
        />
      ) : tabTampil === 'Co-Ins Scale' ? (
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_COIN_SCALE}
          baris={(warisan?.skalaKoasuransi ?? []).map((b) =>
            [b.bagianKoasuransi, b.persenLimit, b.penyusun, b.disusunPada].map((v, i) =>
              selAngka(JENIS_COIN_SCALE[i] ?? 'teks', v),
            ),
          )}
          petunjukKosong={FORM_KONTRAK.petunjukSkalaKoasuransi}
          mode={mode}
        />
      ) : tabTampil === 'Exclusions' ? (
        <TabTeksPanjang
          judul={tabTampil}
          tab={warisan?.pengecualian}
          petunjukKosong={FORM_KONTRAK.petunjukTeksPengecualian}
        />
      ) : tabTampil === 'Special Conditions' ? (
        <TabTeksPanjang
          judul={tabTampil}
          tab={warisan?.syaratKhusus}
          petunjukKosong={FORM_KONTRAK.petunjukTeksSyarat}
        />
      ) : tabTampil === 'Accumulation' ? (
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_AKUMULASI}
          baris={(warisan?.akumulasi ?? []).map((b) => [
            b.periode,
            b.tanggalLapor,
            b.hariKirim,
            b.jatuhTempoKirim,
          ])}
          petunjukKosong={FORM_KONTRAK.petunjukAkumulasi}
        />
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
              {tabTampil === 'Retro' ? FORM_KONTRAK.jarangDipakai : FORM_KONTRAK.belumDibangun}
            </span>
            <span className="trin__belum-petunjuk">
              {tabTampil === 'Retro'
                ? FORM_KONTRAK.jarangDipakaiPetunjuk
                : FORM_KONTRAK.belumDibangunPetunjuk}
            </span>
          </div>
        </Panel>
      )}
      </fieldset>
      {/* ⭐ ATTACHMENT · tombol · HISTORY — di BAWAH strip tab, bukan di
          dalamnya. Begitu layar lama menyusunnya: panel Attachment, deret
          tombol Save/Close/Actions, lalu panel History. Ketiganya berlaku
          untuk kontraknya, bukan untuk tab yang kebetulan terbuka. */}
      <PanelLampiran
        kategori={warisan?.kategoriLampiran ?? []}
        berkas={warisan?.lampiran ?? []}
      />

      {/* ⭐ `Close` HIDUP — ia satu-satunya dari ketiganya yang tidak menuntut
          jalur tulis: ia hanya menutup layar. `Save` dan `Actions` tetap mati,
          dan sebabnya berbeda satu sama lain:

            Save    — nol jalur tulis. Tombol simpan yang tidak menyimpan
                      adalah cara tercepat kehilangan suntingan tanpa seorang
                      pun tahu.
            Actions — syarat perannya belum terbaca di ekspor. Menghidupkannya
                      berarti menebak siapa yang boleh menekannya.

          ⛔ Keduanya DIMATIKAN, bukan disembunyikan: tombol hilang terbaca
          sebagai layar yang berbeda, tombol mati terbaca sebagai kemampuan
          yang belum datang. */}
      <div className="trin__aksi" role="group" aria-label={FORM_KONTRAK.judul}>
        <button type="button" className="btn" disabled>
          Save
        </button>
        <button type="button" className="btn" onClick={onKembali}>
          Close
        </button>
        <button type="button" className="btn" disabled>
          Actions
        </button>
      </div>

      <PanelHistory baris={warisan?.catatan ?? []} />

    </div>
  )
}

/** Dipakai uji — menjaga cabang `Gagal`/`Memuat` tetap ada dan tidak dipakai salah. */
export const cabangKeadaan = { Gagal, Memuat }
