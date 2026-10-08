// Satu panel form Treaty In — dipakai DUA KALI berdampingan oleh layar
// Adjustment: Old Data (baca-saja) dan New Data.
//
// ⛔ SATU komponen untuk kedua panel. Yang membedakan keduanya hanya DATA —
// spesifikasi kepala form (`labelsPenyesuaian.ts`), kerangka tab bangkitan
// (`ekspor/kerangka.gen.ts`, satu Section per sisi × cabang) — dan satu
// bendera: `bacaSaja`.
//
// ⛔ Pola tata letaknya DITIRU dari `modul/treatyin/frontend/pages/
// FormKontrakTreatyIn.tsx`; nol impor dari modul itu.

import { useEffect, useRef, useState } from 'react'

import { Area, Field, Panel, Pilih, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilAgenTreatyIn,
  ambilIndukSpreadingTreatyIn,
  ambilKelasBisnisTreatyIn,
  ambilOpsiKepalaTreatyIn,
  ambilOpsiLimitsTreatyIn,
  type BarisBersarang,
  type SisiPenyesuaian,
} from '../api'
import { GRID_KURS, KERANGKA_TAB } from '../ekspor/kerangka.gen'
import { syaratTerpenuhi } from '../ekspor/syarat'
import { KURS, PENYESUAIAN, PILIH_AGEN, PILIHAN_MEDAN, PRO_RATA, type MedanForm } from '../labelsPenyesuaian'
import { GridEkspor, RenderKerangka, type KonteksKerangka } from './KerangkaTab'
import type { JenisTulis } from './aksiTombol'
import { gabungPohon, gabungTigaArah } from './baris'
import { ikutBerubah } from './rumusKepala'
import { FieldTanggalKetik } from './TanggalKetik'
import { opsiUntuk, type DataOpsi, type Opsi, type PemuatOpsi } from './pilihan'
import PilihAgen from './PilihAgen'
import { AreaBacaSaja, BelumDibangun, Centang, MedanTakAda, TanggalBacaSaja, selNilai } from './medan'

export { selNilai } from './medan'

/**
 * Tab yang DISEMBUNYIKAN dari strip — keputusan pemilik proses 7 Oktober
 * 2026: *"untuk sementara retro di hide dari tab sampai ada perintah dari
 * developer"*. Berlaku bagi kedua panel dan ketiga cabang (`Retro`, `Actual
 * Retro` Adjust Premium). Kerangka dan syaratnya tetap dibangkitkan.
 */
export const TAB_DISEMBUNYIKAN: ReadonlySet<string> = new Set(['Retro', 'Actual Retro'])

/** `ViewState` sesi — `0` Edit, `1` View. */
export type ModeLayar = '0' | '1'

const tanpaAksi = () => undefined

/** Rute daftar BERPARAMETER — rute Treaty In yang SAMA dengan layar Treaty In. */
const PEMUAT_OPSI: PemuatOpsi = {
  kelasBisnis: ambilKelasBisnisTreatyIn,
  indukSpreading: ambilIndukSpreadingTreatyIn,
}

/**
 * Pilihan dropdown kepala — rute `opsi-kepala` Treaty In (nilai tersimpan
 * berlabel tampil, urutan yang SAMA dengan form Treaty In), atau daftar
 * cadangan `PILIHAN_MEDAN` selama rute itu belum menjawab.
 */
export function opsiKepala(kunci: string, kepala: DataOpsi['kepala']): Opsi[] {
  const dariRute: Readonly<Record<string, readonly Opsi[] | undefined>> = {
    Bordeaux: kepala?.bordereaux,
    AccountingMode: kepala?.caraPembukuan,
    AccountingModeNonProp: kepala?.caraPembukuanNonProp,
  }
  return [...(dariRute[kunci] ?? (PILIHAN_MEDAN[kunci] ?? []).map((v) => ({ value: v, label: v })))]
}

/** Satu medan kepala form menurut spesifikasinya. */
function Medan({
  spek,
  nilai,
  ada,
  bacaSaja,
  onUbah,
  opsi,
}: {
  spek: MedanForm
  nilai: string
  ada: boolean
  bacaSaja: boolean
  onUbah: (v: string) => void
  /** Pilihan medan `pilih` — `opsiKepala`. */
  opsi?: Opsi[]
}) {
  if (!ada) return <MedanTakAda label={spek.label} />
  if (bacaSaja) {
    switch (spek.bentuk) {
      case 'area':
        return <AreaBacaSaja label={spek.label} nilai={nilai} />
      case 'tanggal':
        return <TanggalBacaSaja label={spek.label} nilai={nilai} />
      case 'centang':
        return <Centang label={spek.label} nilai={nilai} bacaSaja onUbah={onUbah} />
      default:
        return <Field label={spek.label} value={nilai} onChange={tanpaAksi} readOnly />
    }
  }
  switch (spek.bentuk) {
    case 'area':
      return <Area label={spek.label} value={nilai} onChange={onUbah} />
    case 'tanggal':
      // ⛔ `nilai` adalah nilai TERSIMPAN (`20180101`) atau bentuk kabel
      // (`DD-MM-YYYY`) yang medan ini sendiri kembalikan — keduanya bentuk
      // yang `keInputTanggal` kenali. Nol terjemahan tampil masuk ke sini.
      // ⭐ Dapat DIKETIK atau dipilih dari kalender — `TanggalKetik.tsx`.
      return <FieldTanggalKetik label={spek.label} value={nilai} onChange={onUbah} />
    case 'pilih':
      return (
        <Pilih
          label={spek.label}
          value={nilai}
          onChange={onUbah}
          opsi={opsi ?? opsiKepala(spek.kunci, undefined)}
        />
      )
    case 'centang':
      return <Centang label={spek.label} nilai={nilai} bacaSaja={false} onUbah={onUbah} />
    case 'tampil':
      // `pxDisplayText` — tampil saja di sistem lama pula.
      return <Field label={spek.label} value={nilai} onChange={tanpaAksi} readOnly />
    default:
      return <Field label={spek.label} value={nilai} onChange={onUbah} />
  }
}

/**
 * Blok Pro Rate kepala form New — `IsProRate == true` @264580.
 *
 * ⛔ Ketiganya `pyReadOnly=true` tanpa syarat (@287433 @306188 @327367).
 * Hari adalah CACAH — tidak diformat. `ProRatePercent` 4 desimal: ekspor
 * menyatakannya (@327220), dan desimal per kolom mengikuti ekspor.
 */
function BlokProRata({ medan }: { medan: Readonly<Record<string, string>> }) {
  const persen = medan.ProRatePercent
  return (
    <div className="tria__prorata">
      <span className="tria__nilai-label">{PRO_RATA.judul}</span>
      <span>{medan.ProRateDays ?? ''}</span>
      <span>{PRO_RATA.pemisah}</span>
      <span>{medan.ProRateTotalDays ?? ''}</span>
      <span className="tria__nilai-label">{PRO_RATA.persen}</span>
      <span className="tria__angka">{persen === undefined ? '' : selNilai('persen', persen, PRO_RATA.desimal)}</span>
      <span>{PRO_RATA.tandaPersen}</span>
    </div>
  )
}

export interface SisiFormProps {
  judul: string
  /** Halaman yang panel ini tampilkan — `OLDDATA` atau akar. */
  sisi: SisiPenyesuaian
  /** Halaman AKAR — beberapa ikatan panel Old menunjuknya, dan SEMUA syarat dinilai di sana. */
  akar: SisiPenyesuaian
  /**
   * Panel New: halaman `OLDDATA` (sisi Old) — rumus Adjustment membacanya,
   * mis. `TreatyInSetValueInstallment` langkah 10.
   */
  lama?: SisiPenyesuaian
  /** ⛔ Panel Old: `true`, tanpa pengecualian. */
  bacaSaja: boolean
  mode: ModeLayar
  cabang: string
  medanKiri: readonly MedanForm[]
  medanKanan: readonly MedanForm[]
  /** Section tab sisi × cabang, mis. `TreatyInTabsNonProportionalOldData`. */
  bagian: string
  tab: readonly string[]
  /**
   * ⭐ Panel New: keadaan TERKINI dilaporkan ke form — tombol Save/Submit
   * mengirimnya. Suntingan tetap hidup di panel ini.
   */
  onKini?: (s: SisiPenyesuaian) => void
  /** Submit / Decline offer EDM di tab Information & Submit. */
  tulis?: (jenis: JenisTulis) => void
  sibukTulis?: boolean
}

export default function SisiForm({
  judul,
  sisi,
  akar,
  lama,
  bacaSaja,
  mode,
  cabang,
  medanKiri,
  medanKanan,
  bagian,
  tab,
  onKini,
  tulis,
  sibukTulis,
}: SisiFormProps) {
  // Suntingan panel New hidup di sini SAJA — kepala, medan tab, dan baris
  // grid (Add/Delete). Isian tidak masuk basis data sebelum Save, dan Save
  // masih mati. Pemanggil memberi `key` per penyesuaian × mode sehingga ia
  // lahir ulang dari data tersimpan.
  // ⭐ Pohon digabung ke larik SEKALI saat lahir — baris membawa larik
  // anaknya sendiri (`komponen/baris.ts`).
  const [kini, setKini] = useState<SisiPenyesuaian>(() => gabungPohon(sisi))
  useEffect(() => {
    if (!bacaSaja) onKini?.(kini)
  }, [kini, bacaSaja, onKini])
  const tampilSisi = bacaSaja ? gabungPohon(sisi) : kini
  const medan = tampilSisi.medan
  // ⛔ Syarat dinilai atas halaman AKAR dengan `ViewState` SESI.
  const akarKini: Record<string, string> = bacaSaja ? akar.medan : medan
  const halaman: Record<string, string> = { ...akarKini, ViewState: mode }
  // ⭐ Medan yang punya event di ekspor menjalankan DataTransform-nya:
  // Commencement → Treaty Year/Termination, Effective/Is Pro Rate → Pro
  // Rate. Lihat `rumusKepala.ts`.
  const ubahMedan = (kunci: string, v: string) => {
    setKini((x) => {
      const medanBaru = { ...x.medan, [kunci]: v }
      return { ...x, medan: { ...medanBaru, ...ikutBerubah(kunci, medanBaru) } }
    })
  }
  const ubahMedanBanyak = (pasangan: Readonly<Record<string, string>>) => {
    setKini((x) => {
      let medanBaru = { ...x.medan }
      for (const [kunci, v] of Object.entries(pasangan)) {
        medanBaru = { ...medanBaru, [kunci]: v }
        medanBaru = { ...medanBaru, ...ikutBerubah(kunci, medanBaru) }
      }
      return { ...x, medan: medanBaru }
    })
  }
  const ubahLarik = (larik: string, baris: BarisBersarang[]) => {
    setKini((x) => ({ ...x, larik: { ...x.larik, [larik]: baris } }))
  }
  // ⭐ Daftar BERPARAMETER (Class of Business, Spreading Type) — dimuat per
  // parameter saat sel/medannya pertama dirender, sekali per kunci.
  const [muatan, setMuatan] = useState<Record<string, readonly Opsi[]>>({})
  const sedangDimuat = useRef(new Set<string>())
  const muat = (kunci: string, ambil: () => Promise<Opsi[]>) => {
    if (sedangDimuat.current.has(kunci)) return
    sedangDimuat.current.add(kunci)
    ambil()
      .then((xs) => {
        setMuatan((m) => ({ ...m, [kunci]: xs }))
      })
      .catch(() => {
        // Rute gagal — daftar kosong, bukan galat layar (pola daftar lain).
        setMuatan((m) => ({ ...m, [kunci]: [] }))
      })
  }
  // ⭐ Daftar dropdown/autocomplete dimuat SEKALI, hanya bila panel ini
  // dapat disunting. Rute yang gagal meninggalkan medannya kotak teks —
  // bukan galat layar.
  const dapatSunting = !bacaSaja && mode === '0'
  const [dataOpsi, setDataOpsi] = useState<DataOpsi>({})
  const [pilihAgen, setPilihAgen] = useState<'reinsured' | 'source' | null>(null)
  useEffect(() => {
    if (!dapatSunting) return
    let dibuang = false
    void Promise.allSettled([ambilOpsiLimitsTreatyIn(), ambilOpsiKepalaTreatyIn(), ambilAgenTreatyIn()]).then(([l, kp, ag]) => {
      if (dibuang) return
      setDataOpsi({
        limits: l.status === 'fulfilled' ? l.value : undefined,
        kepala: kp.status === 'fulfilled' ? kp.value : undefined,
        agen: ag.status === 'fulfilled' ? ag.value : undefined,
      })
    })
    return () => {
      dibuang = true
    }
  }, [dapatSunting])
  const k: KonteksKerangka = {
    // ⭐ Larik akar ikut — daftar `pageList` halaman SESI (Achievement).
    opsi: (sp, kunci, tempat) =>
      opsiUntuk(sp, kunci, { ...dataOpsi, halaman: tampilSisi.larik, muatan, muat }, tempat, PEMUAT_OPSI),
    sisi: tampilSisi,
    akar,
    lama: lama === undefined ? undefined : gabungPohon(lama),
    halaman,
    // ⭐ Mode Edit panel New — `ViewState 0`. Panel Old tidak pernah.
    ubah: dapatSunting,
    ubahMedan,
    ubahMedanBanyak,
    ubahMedanAkar: ubahMedan,
    ubahLarik,
    // ⭐ Akar panel — Section rincian memperpanjang jalur ini (`konteksBaris`).
    jalur: [],
    // ⭐ Panel New saja — panel Old tidak pernah menulis.
    tulis: bacaSaja ? undefined : tulis,
    sibukTulis,
    master: dataOpsi.limits,
    terapkan: (h) => {
      // ⭐ Hasil RANTAI digabung TIGA ARAH atas keadaan TERKINI: yang rumus
      // ubah menimpa, isian yang diketik selama rute menjawab TIDAK hilang
      // (`gabungTigaArah`, `komponen/baris.ts`).
      const { awal, akhir } = h
      if (awal !== undefined && akhir !== undefined) {
        setKini((x) => gabungTigaArah(awal, x, akhir))
        return
      }
      setKini((x) => ({ medan: { ...x.medan, ...h.medan }, larik: { ...x.larik, ...h.larik } }))
    },
  }

  const kerangka = (t: string) => KERANGKA_TAB[`${bagian}#${t}`]
  const tabTampil = tab.filter((t) => !TAB_DISEMBUNYIKAN.has(t) && syaratTerpenuhi(kerangka(t)?.syarat ?? [], halaman))
  const [tabAktif, setTabAktif] = useState<string>(tabTampil[0] ?? '')
  const aktif = tabTampil.includes(tabAktif) ? tabAktif : (tabTampil[0] ?? '')
  const isi = kerangka(aktif)

  const render = (m: MedanForm, i: number) => {
    if (m.syarat !== undefined && m.syarat !== cabang) return null
    const sumber = m.dari === 'akar' ? akar.medan : medan
    const ada = Object.prototype.hasOwnProperty.call(sumber, m.kunci)
    const kunciBaca = bacaSaja || m.selaluBacaSaja === true || (m.bacaSajaBila?.(halaman) ?? false)
    // ⭐ Medan yang diisi lewat jendela pencarian: nilai tampil + tombol
    // `Choose …` (sel 27/29, `TreatyIn.ViewState != 1`).
    if (m.pilihAgen !== undefined && dapatSunting) {
      return (
        <div key={`${m.kunci}-${i}`} className="field">
          <label className="field__label">{m.label}</label>
          <span className="tria__nilai-pilih">{sumber[m.kunci] ?? ''}</span>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            onClick={() => {
              setPilihAgen(m.pilihAgen ?? null)
            }}
          >
            {m.pilihAgen === 'reinsured' ? PILIH_AGEN.tombolCeding : PILIH_AGEN.tombolSob}
          </button>
        </div>
      )
    }
    return (
      <Medan
        key={`${m.kunci}-${i}`}
        spek={m}
        nilai={sumber[m.kunci] ?? ''}
        ada={ada}
        bacaSaja={kunciBaca}
        opsi={m.bentuk === 'pilih' ? opsiKepala(m.kunci, dataOpsi.kepala) : undefined}
        onUbah={(v) => {
          ubahMedan(m.kunci, v)
        }}
      />
    )
  }

  // ⚠️ `CurrencyID` → `Currency` hanya untuk TAMPIL — lihat `KURS`.
  const kursAsli = GRID_KURS[bacaSaja ? 'lama' : 'baru']
  // ⛔ Hanya untuk TAMPIL baca-saja. Di mode Edit sel itu dropdown
  // `BrowseCurrency_RD` bernilai `.ID` — menukarnya ke nama membuat nilai
  // tersimpan terbaca "di luar daftar".
  const kurs = dapatSunting ? kursAsli : { ...kursAsli, kunci: kursAsli.kunci.map((x) => KURS.gantiKunciTampil[x] ?? x) }

  return (
    <section className={'tria__sisi' + (bacaSaja ? ' tria__sisi--lama' : '')} aria-label={judul}>
      {pilihAgen !== null && (
        <PilihAgen
          judul={pilihAgen === 'reinsured' ? PILIH_AGEN.judulCeding : PILIH_AGEN.judulSob}
          onTutup={() => {
            setPilihAgen(null)
          }}
          onPilih={(nama, id) => {
            // `TreatyInSetReinsured`: type reinsured → Ceding/CedingID,
            // source → LeadingReinsSource/LeadingReinsSourceID.
            const [kNama, kId] = pilihAgen === 'reinsured' ? ['Ceding', 'CedingID'] : ['LeadingReinsSource', 'LeadingReinsSourceID']
            ubahMedan(kNama, nama)
            ubahMedan(kId, id)
          }}
        />
      )}
      <Panel judul={judul}>
        <div className="tria__dwikolom">
          <div className="tria__kolom">{medanKiri.map(render)}</div>
          <div className="tria__kolom">
            {medanKanan.map(render)}
            {!bacaSaja && medan.IsProRate === 'true' && <BlokProRata medan={medan} />}
          </div>
        </div>
        <h5 className="tria__subjudul">{KURS.judul}</h5>
        <GridEkspor g={kurs} k={k} />
      </Panel>

      {tabTampil.length > 0 && <StripTab tab={tabTampil} aktif={aktif} onPilih={setTabAktif} />}
      {aktif !== '' && (
        <Panel judul={aktif}>
          {isi === undefined ? (
            <BelumDibangun judul={PENYESUAIAN.belumDibangun} />
          ) : (
            <div className="tria__blok">
              <RenderKerangka isi={isi.isi} k={k} />
            </div>
          )}
        </Panel>
      )}
    </section>
  )
}
