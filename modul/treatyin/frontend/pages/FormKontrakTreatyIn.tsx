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

import { useEffect, useState } from 'react'

import {
  Area,
  Field,
  FieldTanggal,
  Gagal,
  Kosong,
  Memuat,
  Panel,
  Pilih,
  StripTab,
} from '../../../../inti/frontend/components/ui/dasar'
import { formatNumber, formatPersen } from '../../../../inti/frontend/lib/format'
import {
  ambilKontrakWarisan,
  type BarisLayerWarisan,
  type BarisPeriodeWarisan,
  type KontrakWarisan,
  type TabTeksWarisan,
} from '../api'
import {
  FORM_KONTRAK,
  REPORTING_PERIOD,
  TAB_NON_PROPORSIONAL,
  TAB_PROPORSIONAL,
  KOLOM_PORTOFOLIO,
  KOLOM_AKUMULASI,
  KOLOM_EGNPI,
  KOLOM_RETENSI,
  KOLOM_ANGSURAN,
  KOLOM_CATATAN,
  KOLOM_LIMITS,
  KOLOM_SHARE,
  KOLOM_EVENT_LIMITS,
  KOLOM_RNM_SHARE,
  KOLOM_COIN_SCALE,
  JENIS_LIMITS,
  JENIS_SHARE,
  JENIS_EVENT_LIMITS,
  JENIS_RNM_SHARE,
  JENIS_COIN_SCALE,
  DESIMAL_UANG,
  DESIMAL_PERSEN,
  DESIMAL_PERSEN_SHARE,
  type JenisAngka,
} from '../labels'

/**
 * Satu sel grid, diformat menurut GOLONGAN kolomnya.
 *
 * ⛔ Pemformatnya `inti/frontend/lib/format.ts` — dipanggil, tidak ditulis
 * ulang. Ia bekerja pada DIGIT, bukan pada float, sehingga nilai uang
 * berdigit banyak (`LIMIT_100` mencapai 1,8 triliun) tidak kehilangan angka
 * terakhirnya diam-diam. Teks yang bukan angka dikembalikan APA ADANYA —
 * itulah yang menjaga `>=30% up to < 50%` dan nilai aneh warisan tetap
 * terlihat mentah alih-alih berpura-pura nol.
 *
 * ⛔ Pemformatan HANYA di lapis tampilan. Nilai tersimpan tetap teks apa
 * adanya; `repository` dan `services` tidak menyentuhnya.
 */
export function selAngka(jenis: JenisAngka, nilai: string): string {
  switch (jenis) {
    case 'uang':
      return formatNumber(nilai, DESIMAL_UANG)
    case 'persen':
      if (!angkaMurni(nilai)) return nilai
      return formatPersen(nilai, DESIMAL_PERSEN)
    case 'persenShare':
      // ⛔ DELAPAN desimal — keputusan pemilik proses 4 Oktober 2026,
      // `KEPUTUSAN-PENYELARASAN-REPO.md` §13. Sama persis dengan batas
      // penyimpanan `NUMBER(38,8)`: menampilkan lebih berarti mengaku lebih
      // teliti daripada yang sistem simpan, menampilkan kurang menutupi
      // selisih yang orang cari ketika memeriksa.
      if (!angkaMurni(nilai)) return nilai
      return formatPersen(nilai, DESIMAL_PERSEN_SHARE)
    default:
      return nilai
  }
}

/**
 * ⛔ Nilai yang BUKAN angka tidak boleh diberi tanda `%`.
 *
 * `formatPersen` mengembalikan teks bukan-angka apa adanya lalu MENEMPELKAN
 * `%` padanya — benar untuk nilai kosong, salah untuk pita seperti
 * `>=30% up to < 50%`, yang menjadi `>=30% up to < 50%%` dengan dua tanda.
 *
 * ⚠️ Ini BUKAN pemformat kedua: ia hanya memutuskan APAKAH `formatPersen`
 * dipanggil. `format.ts` tidak disentuh — aturannya dipenuhi lewat argumen
 * dan lewat pemanggilan, persis seperti yang ronde ini tuntut.
 *
 * Hari ini hanya `CoInShare` yang berbentuk pita, dan ia digolongkan `teks`
 * sehingga tidak melewati jalur ini sama sekali. Penjaga ini ada untuk
 * salah-golong BERIKUTNYA — satu huruf di `JENIS_*` sudah cukup.
 */
function angkaMurni(nilai: string): boolean {
  return /^[+-]?\d*(?:\.\d*)?$/.test(nilai.trim()) && /\d/.test(nilai)
}

/**
 * Tab yang isinya SATU medan teks panjang — Exclusions, Special Conditions.
 *
 * ⛔ Jalan B, keputusan §15: teksnya dibaca dari `JSONDATA`, nol tabel
 * pendaratan. Yang dikerjakan layar hanya menampilkannya.
 *
 * ⛔ TIGA keadaan, dan ketiganya terlihat berbeda:
 *
 *   ada isinya            teksnya, dapat digulir, plus ejaan asalnya.
 *   ejaan cabang kosong   `Kosong` — "belum ada DATA". TIDAK diisi dari
 *                         ejaan lain: isinya BERBEDA, bukan salinan basi.
 *   ejaan lain juga ada   peringatan yang MENYEBUT ejaannya, supaya
 *                         pembacanya tahu ada teks yang ia tidak lihat.
 */
function TabTeksPanjang({
  judul,
  tab,
  petunjukKosong,
}: {
  judul: string
  tab: TabTeksWarisan | undefined
  petunjukKosong: string
}) {
  const isi = tab?.isi ?? ''
  const lain = tab?.ejaanLain ?? []
  return (
    <Panel judul={judul}>
      {lain.length > 0 && (
        <span className="trin__teks-lain" role="note">
          {FORM_KONTRAK.ejaanLainBerisi} {lain.join(', ')}
        </span>
      )}
      {isi === '' ? (
        <Kosong pesan={FORM_KONTRAK.tanpaTeks} petunjuk={petunjukKosong} />
      ) : (
        <>
          <span className="trin__teks-asal">
            {FORM_KONTRAK.ejaanDipakai} {tab?.ejaan}
          </span>
          <div className="trin__teks" tabIndex={0}>
            {isi}
          </div>
        </>
      )}
    </Panel>
  )
}

/** Menyusun baris grid dari baris layer, satu nilai per kolom. */
function barisLayer(
  layer: readonly BarisLayerWarisan[],
  ambil: (b: BarisLayerWarisan) => readonly string[],
  jenis: readonly JenisAngka[],
): string[][] {
  return layer.map((b) => ambil(b).map((v, i) => selAngka(jenis[i] ?? 'teks', v)))
}

/**
 * Medan yang KUNCINYA tidak ada di dokumen warisan kontrak ini.
 *
 * ⛔ MATI dengan keterangan, bukan kotak kosong. Kotak kosong terbaca
 * "belum diisi"; medan mati terbaca "tidak ada di sistem lama". Bedanya
 * menentukan apa yang orang tagih — dan sapuan 3 Oktober 2026 menemukan ia
 * BUKAN kasus langka: `ContractRefNo` tidak ada di 1.112 dari 1.854 kontrak,
 * `TreatyLeader` di 1.195.
 *
 * Pola yang sama sudah dipakai tombol `Choose Ceding`, dan sebabnya tertulis
 * di sana.
 */
function MedanTakAda({ label }: { label: string }) {
  return (
    <div className="field">
      <label className="field__label">{label}</label>
      <input className="field__input" type="text" value="" readOnly disabled />
      <span className="trin__redup">{FORM_KONTRAK.takAdaDiWarisan}</span>
    </div>
  )
}

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
export function tabUntuk(jenis: string): readonly string[] {
  return jenis === NON_PROPORSIONAL ? TAB_NON_PROPORSIONAL : TAB_PROPORSIONAL
}

/** Satu baris grid Rate of Exchange; kosong adalah keadaan awal yang sah. */
interface BarisKurs {
  mataUang: string
  nilaiKeIDR: string
  berlakuDari: string
  berlakuSampai: string
}

/**
 * Grid satu tab yang isinya dibaca dari dokumen warisan.
 *
 * ⛔ Kolomnya DISEBUT pemanggilnya, tidak disimpulkan dari barisnya: baris
 * boleh nol, dan kepala kolom harus tetap terlihat. Layar yang
 * menyembunyikan kolomnya ketika kosong tidak memperlihatkan rancangannya
 * sendiri.
 */
function TabGridWarisan({
  judul,
  kolom,
  baris,
  petunjukKosong,
}: {
  judul: string
  kolom: readonly string[]
  baris: readonly (readonly string[])[]
  petunjukKosong: string
}) {
  return (
    <Panel judul={judul}>
      <div className="table-wrap">
        <table className="trin__tabel">
          <thead>
            <tr>
              {kolom.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {baris.length === 0 && (
              <tr>
                <td colSpan={kolom.length}>
                  <Kosong pesan={FORM_KONTRAK.tanpaBaris} petunjuk={petunjukKosong} />
                </td>
              </tr>
            )}
            {baris.map((b, i) => (
              <tr key={i}>
                {b.map((v, j) => (
                  <td key={j}>{v}</td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Panel>
  )
}

/** Isi tab Reporting Period — satu-satunya tab yang dibangun ronde ini. */
function TabReportingPeriod({ baris: barisWarisan }: { baris: readonly BarisPeriodeWarisan[] }) {
  const [mulai, setMulai] = useState('')
  const [akhir, setAkhir] = useState('')
  const [periode, setPeriode] = useState('')
  const [penyerahan, setPenyerahan] = useState('')
  const [konfirmasi, setKonfirmasi] = useState('')
  const [pelunasan, setPelunasan] = useState('')
  const [galat, setGalat] = useState('')
  // ⭐ Grid ini DIISI dari dokumen warisan, dan tombol Apply menambah ke
  // atasnya. `ReportingPeriodList` berisi di 225 dari 300 dokumen yang
  // disapu; keenam kolomnya sama persis dengan yang layar lama tampilkan.
  const barisAwal: readonly string[][] = barisWarisan.map((b) => [
    b.periode,
    b.hitungOtomatis,
    b.tanggalAwal,
    b.jatuhTempoKirim,
    b.jatuhTempoKonfirmasi,
    b.jatuhTempoBayar,
  ])
  const [baris, setBaris] = useState<readonly string[][]>(barisAwal)

  // ⛔ Pesan galatnya DISALIN apa adanya dari ekspor; lihat `labels.ts`.
  // Jangan diperhalus — pemakai yang mengenali kalimat lamanya tahu ia sedang
  // ditolak oleh aturan yang sama, bukan oleh aturan baru yang mirip.
  function terapkan() {
    if (mulai === '' || penyerahan === '') {
      setGalat(REPORTING_PERIOD.galatKosong)
      return
    }
    setGalat('')
    // ⚠️ Perhitungan jatuh temponya BELUM dibangun: ia menuntut tabel
    // PERIODE_PELAPORAN terisi (tiket 25) dan aturan hari kerjanya, dan
    // keduanya belum ada. Baris sengaja tidak dikarang.
    setBaris(barisAwal)
  }

  return (
    <Panel judul={REPORTING_PERIOD.judul}>
      <div className="form-grid">
        <FieldTanggal label={REPORTING_PERIOD.mulai} value={mulai} onChange={setMulai} />
        <FieldTanggal label={REPORTING_PERIOD.akhir} value={akhir} onChange={setAkhir} />
        <Pilih
          label={REPORTING_PERIOD.periode}
          value={periode}
          onChange={setPeriode}
          opsi={[{ value: REPORTING_PERIOD.periodeNilai, label: REPORTING_PERIOD.periodeNilai }]}
        />
        <Field
          label={`${REPORTING_PERIOD.penyerahan} (${REPORTING_PERIOD.hari})`}
          value={penyerahan}
          onChange={setPenyerahan}
        />
        <Field
          label={`${REPORTING_PERIOD.konfirmasi} (${REPORTING_PERIOD.hari})`}
          value={konfirmasi}
          onChange={setKonfirmasi}
        />
        <Field
          label={`${REPORTING_PERIOD.pelunasan} (${REPORTING_PERIOD.hari})`}
          value={pelunasan}
          onChange={setPelunasan}
        />
      </div>

      <div className="trin__aksi">
        <button type="button" className="btn btn--primary" onClick={terapkan}>
          {REPORTING_PERIOD.terapkan}
        </button>
      </div>

      {/* ⛔ Galat masukan, BUKAN cabang <Gagal>. `Gagal` menyatakan backend
          menolak atau tidak terjangkau; ini aturan layar. Mencampurnya adalah
          salah diagnosa yang `dasar.tsx` baris 344 catat. */}
      {galat !== '' && (
        <p className="trin__galat" role="alert">
          {galat}
        </p>
      )}

      {baris.length === 0 ? (
        <Kosong
          pesan={REPORTING_PERIOD.tanpaBaris}
          petunjuk="Perhitungan jatuh tempo menuntut PERIODE_PELAPORAN terisi (tiket 25); barisnya tidak dikarang."
        />
      ) : (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">{REPORTING_PERIOD.kolomPeriode}</th>
                <th scope="col">{REPORTING_PERIOD.kolomOtomatis}</th>
                <th scope="col">{REPORTING_PERIOD.kolomTanggalAwal}</th>
                <th scope="col">{REPORTING_PERIOD.kolomPenyerahan}</th>
                <th scope="col">{REPORTING_PERIOD.kolomKonfirmasi}</th>
                <th scope="col">{REPORTING_PERIOD.kolomPelunasan}</th>
              </tr>
            </thead>
            <tbody>
              {baris.map((b, i) => (
                <tr key={i}>
                  {b.map((c, j) => (
                    <td key={j}>{c}</td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Panel>
  )
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
  onKembali: () => void
}

export default function FormKontrakTreatyIn({ idKontrak, onKembali }: FormKontrakProps) {
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
  const [cedant, setCedant] = useState('')
  const [asalBisnis, setAsalBisnis] = useState('')
  const [pemimpin, setPemimpin] = useState(false)
  // ⭐ Grid Rate of Exchange dibaca dari dokumen warisan kontrak ini, bukan
  // dari keadaan kosong. `CurrencyList` berisi di 297 dari 300 dokumen yang
  // disapu — "No items" selama ini bukan karena datanya tidak ada, melainkan
  // karena tidak ada yang membacanya.
  const kurs: BarisKurs[] = warisan?.kurs ?? []

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
        setJenis(k.sifatProporsi === NON_PROPORSIONAL ? NON_PROPORSIONAL : PROPORSIONAL)
        setNama(k.namaKontrak)
        setRujukan(k.nomorRujukan)
        setWilayah(k.lingkupWilayah)
        setBordereaux(k.bordereaux)
        setBordereauxNote(k.bordereauxCatatan)
        setMulai(k.tanggalMulai)
        setBerakhir(k.tanggalBerakhir)
        setPembukuan(k.caraPembukuan)
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

  const tab = tabUntuk(jenis)
  const [tabAktif, setTabAktif] = useState<string>(tab[0] ?? '')
  const tabTampil = tab.includes(tabAktif) ? tabAktif : (tab[0] ?? '')

  /**
   * `Treaty Year` — baca-saja, DITURUNKAN dari tanggal mulai.
   *
   * ⚠️ Rumusnya BELUM diverifikasi di ekspor: `TreatyIn.TreatyYear` ada
   * sebagai properti, tetapi aktivitas yang menghitungnya belum ditelusuri.
   * Yang ditampilkan tahun dari `Commencement` apa adanya; begitu rumusnya
   * terbaca, inilah satu tempat yang berubah.
   */
  // ⭐ `TREATYYEAR` adalah KOLOM di `TREATY_IN` — dibaca, bukan dihitung.
  // Bentuk sebelumnya memotong empat aksara pertama `Commencement`, yang
  // kini salah dua kali: tanggalnya sudah berbentuk `dd/mm/yy`, dan
  // tahunnya memang tersimpan sendiri.
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

      <Panel judul={FORM_KONTRAK.judul}>
        <div className="trin__kepala">
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
            <Pilih
              label={FORM_KONTRAK.bordereaux}
              value={bordereaux}
              onChange={setBordereaux}
              opsi={FORM_KONTRAK.bordereauxNilai.map((v) => ({ value: v, label: v }))}
            />
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
            <FieldTanggal label={FORM_KONTRAK.mulai} value={mulai} onChange={setMulai} />
            <FieldTanggal label={FORM_KONTRAK.berakhir} value={berakhir} onChange={setBerakhir} />
            <Field
              label={FORM_KONTRAK.tahunTreaty}
              value={tahunTreaty}
              onChange={() => undefined}
              readOnly
            />
            <Pilih
              label={FORM_KONTRAK.caraPembukuan}
              value={pembukuan}
              onChange={setPembukuan}
              opsi={FORM_KONTRAK.caraPembukuanNilai.map((v) => ({ value: v, label: v }))}
            />
            <div className="trin__pilih-luar">
              <Field label={FORM_KONTRAK.cedant} value={cedant} onChange={setCedant} />
              <button type="button" className="btn btn--ghost btn--sm" disabled>
                {FORM_KONTRAK.pilihCedant}
              </button>
            </div>
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
            <div className="trin__pilih-luar">
              <Field label={FORM_KONTRAK.asalBisnis} value={asalBisnis} onChange={setAsalBisnis} />
              <button type="button" className="btn btn--ghost btn--sm" disabled>
                {FORM_KONTRAK.pilihAsalBisnis}
              </button>
            </div>
          </div>
        </div>
      </Panel>

      {/* ⚠️ Kedua tombol "Choose …" DIMATIKAN, bukan disembunyikan. Keduanya
          membuka pemilih atas tabel yang `ERD.md` §2.8 tempatkan DI LUAR skema
          ini (`CEDANT`, `ASAL_BISNIS`); tanpa modul pemiliknya tidak ada yang
          dapat dipilih. Tombol yang hilang terbaca sebagai layar yang berbeda;
          tombol yang mati terbaca sebagai kemampuan yang belum datang. */}
      {/* ⛔ TOMBOL "Add" DI BARIS KEPALA, bukan di bawah kotak kosong.
          Di ekspor ia duduk sebaris dengan kepala grid; menaruhnya di bawah
          area kosong membuat panel menjulur tinggi dan tombolnya terbaca
          sebagai milik keadaan kosong, bukan milik gridnya. */}
      <section className="panel">
        <div className="trin__panel-kepala">
          <h4 className="panel__title">{FORM_KONTRAK.kurs}</h4>
          <button type="button" className="btn btn--ghost btn--sm" disabled>
            {FORM_KONTRAK.tambah}
          </button>
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
              {kurs.map((b, i) => (
                <tr key={i}>
                  <td>{b.mataUang}</td>
                  <td>{b.nilaiKeIDR}</td>
                  <td>{b.berlakuDari}</td>
                  <td>{b.berlakuSampai}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* Keterangan kaki, satu baris redup — bukan paragraf di tengah layar. */}
      <p className="trin__kaki">{FORM_KONTRAK.catatanPilihLuar}</p>

      {/* Strip tab — himpunannya berganti bersama radio di atas. */}
      <StripTab tab={tab} aktif={tabTampil} onPilih={setTabAktif} />

      {/* ⭐ DUA BELAS TAB KINI BERISI. Tujuh dari tabel pendaratan
          `M_TREATYIN_*`, empat dari `M_TREATY_IN2` (satu tabel warisan,
          empat proyeksi atas baris layer yang sama), dan Rate of Exchange
          dari dokumen. Sisanya tetap menyatakan "belum ada KODE" — dan
          bedanya dijaga: tab berisi yang kosong memakai `Kosong` ("belum
          ada DATA"), tab tanpa kode memakai `.trin__belum`. */}
      {tabTampil === 'Reporting Period' ? (
        <TabReportingPeriod baris={warisan?.periodePelaporan ?? []} />
      ) : tabTampil === 'Portfolio' ? (
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_PORTOFOLIO}
          baris={(warisan?.portofolio ?? []).map((b) => [b.jenis, b.jenisPortfolio, b.keterangan])}
          petunjukKosong={FORM_KONTRAK.petunjukPortofolio}
        />
      ) : tabTampil === 'EGNPI' ? (
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_EGNPI}
          baris={(warisan?.egnpi ?? []).map((b) => [
            b.jumlah, b.jumlahIDR, b.perTanggal, b.kelasBisnis,
            b.mataUang, b.keterangan, b.proporsi, b.kelompokTreaty,
          ])}
          petunjukKosong={FORM_KONTRAK.petunjukTabel}
        />
      ) : tabTampil === 'Maximum Retention' ? (
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_RETENSI}
          baris={(warisan?.retensi ?? []).map((b) => [
            b.jumlah, b.kelasBisnis, b.mataUang, b.keterangan, b.kelompokTreaty,
          ])}
          petunjukKosong={FORM_KONTRAK.petunjukTabel}
        />
      ) : tabTampil === 'Installment' ? (
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_ANGSURAN}
          baris={(warisan?.angsuran ?? []).map((b) => [
            b.angsuran, b.mataUang, b.jumlah, b.persen,
            b.jatuhTempo, b.tanggalBayar, b.wpc,
          ])}
          petunjukKosong={FORM_KONTRAK.petunjukTabel}
        />
      ) : tabTampil === 'Information & Submit' ? (
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_CATATAN}
          baris={(warisan?.catatan ?? []).map((b) => [
            b.tanggal, b.operator, b.disetujui, b.catatan,
          ])}
          petunjukKosong={FORM_KONTRAK.petunjukTabel}
        />
      ) : tabTampil === 'Limits' ? (
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_LIMITS}
          baris={barisLayer(
            warisan?.layer ?? [],
            (b) => [
              b.layer, b.jenisLayer, b.dasarCover, b.jenisTreaty, b.mataUang,
              b.limit100, b.retensiCedant, b.mdp, b.rasioMDP, b.rol,
              b.adjRate, b.premiEarned, b.relasiMataUang, b.mataUangLimit,
            ],
            JENIS_LIMITS,
          )}
          petunjukKosong={FORM_KONTRAK.petunjukLayer}
        />
      ) : tabTampil === 'Share' ? (
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_SHARE}
          baris={barisLayer(
            warisan?.layer ?? [],
            (b) => [
              b.layer, b.persenCession, b.jenisPenyebaran, b.persenBrokerage,
              b.cessionKeRI, b.qsor, b.qsri, b.liabilityQSOR, b.liabilityQSRI,
              b.epi100, b.riogr,
            ],
            JENIS_SHARE,
          )}
          petunjukKosong={FORM_KONTRAK.petunjukLayer}
        />
      ) : tabTampil === 'Event Limits' ? (
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_EVENT_LIMITS}
          baris={barisLayer(
            warisan?.layer ?? [],
            (b) => [b.layer, b.mataUang, b.gempa],
            JENIS_EVENT_LIMITS,
          )}
          petunjukKosong={FORM_KONTRAK.petunjukEventLimits}
        />
      ) : tabTampil === 'RNM Share' ? (
        <TabGridWarisan
          judul={tabTampil}
          kolom={KOLOM_RNM_SHARE}
          baris={barisLayer(
            warisan?.layer ?? [],
            (b) => [
              b.layer, b.rnmShare, b.liabilityRNM, b.mdpRNM100,
              b.epiRNMQS100, b.rnmRetainedPremi, b.rnmQSPremi,
            ],
            JENIS_RNM_SHARE,
          )}
          petunjukKosong={FORM_KONTRAK.petunjukLayer}
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
          <div className="trin__belum" role="note">
            <span className="trin__belum-judul">{FORM_KONTRAK.belumDibangun}</span>
            <span className="trin__belum-petunjuk">{FORM_KONTRAK.belumDibangunPetunjuk}</span>
          </div>
        </Panel>
      )}
    </div>
  )
}

/** Dipakai uji — menjaga cabang `Gagal`/`Memuat` tetap ada dan tidak dipakai salah. */
export const cabangKeadaan = { Gagal, Memuat }
