// Tab **Accumulation** (cabang Proporsional) — `Section/TreatyInTabsProportional.xml`,
// tab `Accumulation`; bentuknya gambar Pega `19-prop-menu-accumulation.png`.
//
// ---------------------------------------------------------------------
// ⭐ 7 Oktober 2026 — MENGGANTIKAN `TabGridWarisan`
// ---------------------------------------------------------------------
// Grid umum hanya tahu kotak teks dan baris kosong. Ekspor memberi tab ini:
//
//   medan `Period`            `TreatyIn.AccumulationPeriod` pxDropdown,
//                             change → `TreatyInSetAccountReport`
//   grid `AccumulationList`   Period (baca-saja) · Reporting Date ·
//                             Submission Days · Submission Due (baca-saja)
//                             sel Reporting Date / Submission Days:
//                             change → `TreatyInAccumulationSetSubDue`
//   Add (kepala)              DataTransform `TreatyInAddAccumulation`
//   Delete (baris)            `deleteRow`
//
// ⭐ KETERGANTUNGAN ANTARTAB: `TreatyInSetAccountReport` membaca
// `ReportingStart`/`ReportingEnd` tab REPORTING PERIOD. Keduanya dibaca dari
// penampung halaman (`../halaman.tsx`) — isian tab itu, bahkan bila tabnya
// sedang tidak tampil. Rumusnya di services (`hitung_akumulasi.go`).
//
// Baca-saja sel: `TreatyIn.ViewState = 1` (mode lihat → `<fieldset disabled>`
// form) atau `TreatyIn.EDMMaterialType = 2`.

import { useState } from 'react'

import { PemicuUbah } from './pemicuUbah'

import { Panel, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { hitungAkumulasi, type BarisAkumulasi, type BarisAkumulasiWarisan } from '../api'
import { useProperti } from '../halaman'
import { KOLOM_AKUMULASI } from '../labels'
import { AKUMULASI, AWALAN_AKUMULASI, OPSI_PERIODE_AKUMULASI } from '../labelsAkumulasi'
import type { ModeForm } from '../mode'
import TanggalRedup from './TanggalRedup'
import { keSimpan, tanggalTampil } from './tanggalIso'

/** Baris kontrak → baris halaman (bentuk TERSIMPAN). */
export const barisAkumulasi = (b: BarisAkumulasiWarisan): BarisAkumulasi => ({
  Period: b.periode,
  ReportDate: b.tanggalLaporAsli,
  SubDays: b.hariKirim,
  SubDueDate: b.jatuhTempoKirimAsli,
})

/** DataTransform `TreatyInAddAccumulation`: Period = awalan + (cacah + 1). */
export function barisAkumulasiBaru(periode: string, cacah: number): BarisAkumulasi {
  return { Period: `${AWALAN_AKUMULASI[periode] ?? ''}${String(cacah + 1)}`, ReportDate: '', SubDays: '', SubDueDate: '' }
}

export default function TabAkumulasi({
  baris,
  mode = 'lihat',
  edmJenisMaterial = '',
}: {
  baris: readonly BarisAkumulasiWarisan[]
  mode?: ModeForm
  edmJenisMaterial?: string
}) {
  const [periode, setPeriode] = useProperti('AccumulationPeriod', '')
  const [isi, setIsi] = useProperti<BarisAkumulasi[]>('AccumulationList', () => baris.map(barisAkumulasi))
  // ⭐ Milik tab Reporting Period — dibaca saja.
  const [mulai] = useProperti('ReportingStart', '')
  const [akhir] = useProperti('ReportingEnd', '')
  const [gagal, setGagal] = useState('')
  const bisaUbah = mode === 'ubah'
  const terkunci = !bisaUbah || edmJenisMaterial.trim() === '2'

  const jalankan = (aksi: 'periode' | 'jatuh-tempo', daftar: readonly BarisAkumulasi[], per: string) => {
    setGagal('')
    hitungAkumulasi({ aksi, AccumulationPeriod: per, ReportingStart: mulai, ReportingEnd: akhir, AccumulationList: daftar })
      .then((h) => {
        setIsi(h.AccumulationList)
      })
      .catch((e: unknown) => {
        setGagal(e instanceof Error ? e.message : String(e))
      })
  }
  const ubahSel = (i: number, kunci: 'ReportDate' | 'SubDays', v: string) => {
    const baru = isi.map((b, j) => (j === i ? { ...b, [kunci]: v } : b))
    setIsi(baru)
    return baru
  }

  return (
    <Panel judul={AKUMULASI.judul}>
      <div className="trin__akumulasi-periode">
        <Pilih
          label={AKUMULASI.periode}
          value={periode}
          onChange={(v) => {
            setPeriode(v)
            // `change` → DT pra-refresh `TreatyInDeleteAccumulationLists`
            // (`TreatyIn.AccumulationList` REMOVE) LALU
            // `TreatyInSetAccountReport` — daftar dikirim KOSONG, jadi
            // `none` (langkah 2, keluar) meninggalkan daftar kosong, persis
            // Pega. Dahulu daftar lama terkirim dan `none` membiarkannya.
            if (!terkunci) jalankan('periode', [], v)
          }}
          opsi={[...OPSI_PERIODE_AKUMULASI]}
        />
      </div>
      {gagal !== '' && (
        <p className="trin__galat" role="alert">
          {gagal}
        </p>
      )}
      <div className="table-wrap">
        <table className="trin__tabel">
          <thead>
            <tr>
              {KOLOM_AKUMULASI.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              {bisaUbah && (
                <th scope="col">
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    disabled={terkunci}
                    onClick={() => {
                      setIsi([...isi, barisAkumulasiBaru(periode, isi.length)])
                    }}
                  >
                    {AKUMULASI.tambah}
                  </button>
                </th>
              )}
            </tr>
          </thead>
          <tbody>
            {isi.length === 0 && (
              <tr>
                <td colSpan={KOLOM_AKUMULASI.length + (bisaUbah ? 1 : 0)}>{AKUMULASI.tanpaBaris}</td>
              </tr>
            )}
            {isi.map((b, i) => (
              <tr key={i}>
                <td>{b.Period}</td>
                <td>
                  {terkunci ? (
                    tanggalTampil(b.ReportDate)
                  ) : (
                    <TanggalRedup
                      label={KOLOM_AKUMULASI[1]}
                      value={b.ReportDate}
                      onChange={(v) => {
                        // Kotak tanggal: `change` segera → SetSubDue.
                        jalankan('jatuh-tempo', ubahSel(i, 'ReportDate', keSimpan(v)), periode)
                      }}
                    />
                  )}
                </td>
                <td>
                  {terkunci ? (
                    b.SubDays
                  ) : (
                    // `change` Pega — isian ditinggalkan DAN nilainya berubah.
                    <PemicuUbah
                      nilai={b.SubDays}
                      aksi={() => {
                        jalankan('jatuh-tempo', isi, periode)
                      }}
                    >
                      <input
                        className="field__input"
                        type="text"
                        inputMode="numeric"
                        aria-label={KOLOM_AKUMULASI[2]}
                        value={b.SubDays}
                        onChange={(e) => {
                          ubahSel(i, 'SubDays', e.target.value)
                        }}
                      />
                    </PemicuUbah>
                  )}
                </td>
                <td>{tanggalTampil(b.SubDueDate)}</td>
                {bisaUbah && (
                  <td>
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      disabled={terkunci}
                      aria-label={`${AKUMULASI.hapus} ${String(i + 1)}`}
                      onClick={() => {
                        setIsi(isi.filter((_, j) => j !== i))
                      }}
                    >
                      {AKUMULASI.hapus}
                    </button>
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Panel>
  )
}
