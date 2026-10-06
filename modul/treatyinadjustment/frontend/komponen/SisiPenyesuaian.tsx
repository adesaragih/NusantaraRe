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

import { useState } from 'react'

import { Area, Field, FieldTanggal, Panel, Pilih, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import type { SisiPenyesuaian } from '../api'
import { GRID_KURS, KERANGKA_TAB } from '../ekspor/kerangka.gen'
import { syaratTerpenuhi } from '../ekspor/syarat'
import { KURS, PENYESUAIAN, PILIHAN_MEDAN, PRO_RATA, type MedanForm } from '../labelsPenyesuaian'
import { GridEkspor, RenderKerangka, type KonteksKerangka } from './KerangkaTab'
import { AreaBacaSaja, BelumDibangun, Centang, MedanTakAda, TanggalBacaSaja, selNilai } from './medan'

export { selNilai } from './medan'

/** `ViewState` sesi — `0` Edit, `1` View. */
export type ModeLayar = '0' | '1'

const tanpaAksi = () => undefined

/** Satu medan kepala form menurut spesifikasinya. */
function Medan({
  spek,
  nilai,
  ada,
  bacaSaja,
  onUbah,
}: {
  spek: MedanForm
  nilai: string
  ada: boolean
  bacaSaja: boolean
  onUbah: (v: string) => void
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
      // ⛔ `nilai` adalah nilai TERSIMPAN (`20180101`) atau yang
      // `FieldTanggal` sendiri kembalikan (`DD-MM-YYYY`) — keduanya bentuk
      // yang `keInputTanggal` kenali. Nol terjemahan tampil masuk ke sini.
      return <FieldTanggal label={spek.label} value={nilai} onChange={onUbah} />
    case 'pilih':
      return (
        <Pilih
          label={spek.label}
          value={nilai}
          onChange={onUbah}
          opsi={(PILIHAN_MEDAN[spek.kunci] ?? []).map((v) => ({ value: v, label: v }))}
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
  /** ⛔ Panel Old: `true`, tanpa pengecualian. */
  bacaSaja: boolean
  mode: ModeLayar
  cabang: string
  medanKiri: readonly MedanForm[]
  medanKanan: readonly MedanForm[]
  /** Section tab sisi × cabang, mis. `TreatyInTabsNonProportionalOldData`. */
  bagian: string
  tab: readonly string[]
}

export default function SisiForm({
  judul,
  sisi,
  akar,
  bacaSaja,
  mode,
  cabang,
  medanKiri,
  medanKanan,
  bagian,
  tab,
}: SisiFormProps) {
  // Suntingan kepala New hidup di sini SAJA — Save mati, jadi nol terkirim.
  // Pemanggil memberi `key` per penyesuaian × mode sehingga ia lahir ulang.
  const [ubahan, setUbahan] = useState<Record<string, string>>({})
  const medan: Record<string, string> = { ...sisi.medan, ...ubahan }
  // ⛔ Syarat dinilai atas halaman AKAR dengan `ViewState` SESI.
  const akarKini: Record<string, string> = bacaSaja ? akar.medan : medan
  const halaman: Record<string, string> = { ...akarKini, ViewState: mode }
  const k: KonteksKerangka = { sisi, akar, halaman }

  const kerangka = (t: string) => KERANGKA_TAB[`${bagian}#${t}`]
  const tabTampil = tab.filter((t) => syaratTerpenuhi(kerangka(t)?.syarat ?? [], halaman))
  const [tabAktif, setTabAktif] = useState<string>(tabTampil[0] ?? '')
  const aktif = tabTampil.includes(tabAktif) ? tabAktif : (tabTampil[0] ?? '')
  const isi = kerangka(aktif)

  const render = (m: MedanForm, i: number) => {
    if (m.syarat !== undefined && m.syarat !== cabang) return null
    const sumber = m.dari === 'akar' ? akar.medan : medan
    const ada = Object.prototype.hasOwnProperty.call(sumber, m.kunci)
    const kunciBaca = bacaSaja || m.selaluBacaSaja === true || (m.bacaSajaBila?.(halaman) ?? false)
    return (
      <Medan
        key={`${m.kunci}-${i}`}
        spek={m}
        nilai={sumber[m.kunci] ?? ''}
        ada={ada}
        bacaSaja={kunciBaca}
        onUbah={(v) => {
          setUbahan((u) => ({ ...u, [m.kunci]: v }))
        }}
      />
    )
  }

  // ⚠️ `CurrencyID` → `Currency` hanya untuk TAMPIL — lihat `KURS`.
  const kursAsli = GRID_KURS[bacaSaja ? 'lama' : 'baru']
  const kurs = { ...kursAsli, kunci: kursAsli.kunci.map((x) => KURS.gantiKunciTampil[x] ?? x) }

  return (
    <section className={'tria__sisi' + (bacaSaja ? ' tria__sisi--lama' : '')} aria-label={judul}>
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
