// Tab **Installment** (cabang Non-Proportional) — dari ekspor Pega.
//
// ---------------------------------------------------------------------
// ⭐ 7 Oktober 2026 — MENGGANTIKAN `TabGridWarisan`
// ---------------------------------------------------------------------
// Bentuknya `Section/TreatyInTabsNonProportional.xml` tab `Installment` dan
// rincian baris `Section/Installments.xml` (FlowAction `Installments`,
// `pyEditingMode = expandPane`), gambar Pega 38:
//
//   Installment [isian]  [Update Value]
//   grid Currency ▸ → Installment · Due Date · WPC (In Days) · Payment Date ·
//                     % Installment · Amount   +  % Total · Total
//   grid Total Installment Amount · Value
//   [Update Total]
//
// Rumus (rute `/hitung/angsuran`, `hitung_angsuran.go`, diuji dengan angka
// gambar 38):
//   isian Installment (change)          → TreatyInSetValueInstallment
//   Update Value                         → TreatyInSetValueInstallment(status=update)
//   sel % Installment / Amount (change)  → SetTotalInstallment(editpercentage / —)
//   Update Total                         → TreatyInNPSetTotal(installment)
//
// ⭐ KETERGANTUNGAN ANTARTAB: nilai angsuran = `TreatyIn.TotalShareNetNP`
// tab SHARE Non-Prop — pemanggil mengirim total TERKINI tab itu.
//
// ⭐ Isian di PENAMPUNG HALAMAN (`InstallmentNo`, `Installment`,
// `TotalInstallmentNP`) — bertahan saat pindah tab.
//
// ⭐ Sel Due Date → `TreatyInUpdatePaymentDate` (halaman itu), sel WPC →
// `TreatyInUpdatePaymentDate_Act` (semua halaman): Payment Date = Due Date +
// WPC + 1 hari. Kedua Activity diunggah pemakai 7 Oktober 2026.

import { Fragment, useState } from 'react'

import { Panel } from '../../../../inti/frontend/components/ui/dasar'
import { hitungAngsuran, type Angsuran, type BarisAngsuranWarisan, type MasukanAngsuran, type NilaiShare } from '../api'
import { useProperti } from '../halaman'
import { ANGSURAN, KOLOM_RINCIAN_ANGSURAN } from '../labelsAngsuran'
import type { ModeForm } from '../mode'
import { selAngka } from './angka'
import { KotakTanggalKetik } from './TanggalKetik'
import { keSimpan, tanggalTampil } from './tanggalIso'

/**
 * Jadwal kontrak (daftar datar) → halaman `Installment` per mata uang, urut
 * kemunculan. `AmountTotal`/`PctTotal` tersimpannya tidak dikirim backend —
 * kosong sampai rumus dijalankan.
 */
export function angsuranDariWarisan(baris: readonly BarisAngsuranWarisan[]): Angsuran[] {
  const out: Angsuran[] = []
  for (const b of baris) {
    let a = out.find((x) => x.Currency === b.mataUang)
    if (a === undefined) {
      a = { Currency: b.mataUang, AmountTotal: '', PctTotal: '', InstallmentList: [] }
      out.push(a)
    }
    a.InstallmentList.push({
      Installment: b.angsuran,
      DueDate: b.jatuhTempoAsli,
      WPC: b.wpc,
      PaymentDate: b.tanggalBayarAsli,
      Currency: b.mataUang,
      InstallmentPct: b.persen,
      Amount: b.jumlah,
    })
  }
  return out
}

const uang = (v: string) => selAngka(['uang', 2], v)

export default function TabAngsuran({
  baris,
  netPremium,
  edmState = '',
  edmJenisMaterial = '',
  mode = 'lihat',
}: {
  baris: readonly BarisAngsuranWarisan[]
  /** `TreatyIn.TotalShareNetNP` TERKINI — tab Share Non-Prop. */
  netPremium: readonly NilaiShare[]
  edmState?: string
  edmJenisMaterial?: string
  mode?: ModeForm
}) {
  const [no, setNo] = useProperti('InstallmentNo', '')
  const [angsuran, setAngsuran] = useProperti<Angsuran[]>('Installment', () => angsuranDariWarisan(baris))
  const [total, setTotal] = useProperti<NilaiShare[]>('TotalInstallmentNP', [])
  const [buka, setBuka] = useState<ReadonlySet<number>>(() => new Set())
  const [pesan, setPesan] = useState<string[]>([])
  const [sibuk, setSibuk] = useState(false)
  // `TreatyIn.ViewState = 1` (mode lihat) atau `TreatyIn.EDMMaterialType = 2`.
  const bisaUbah = mode === 'ubah'
  const terkunci = !bisaUbah || edmJenisMaterial.trim() === '2'

  const kirim = (m: Pick<MasukanAngsuran, 'aksi' | 'status' | 'angsuran' | 'indeks'> & { installmentNo?: string }, terima: (h: Awaited<ReturnType<typeof hitungAngsuran>>) => void) => {
    setSibuk(true)
    setPesan([])
    hitungAngsuran({
      ...m,
      installmentNo: m.installmentNo ?? no,
      edmState,
      angsuranLama: [],
      netPremium,
    })
      .then((h) => {
        terima(h)
        setPesan(h.pesan)
      })
      .catch((e: unknown) => {
        setPesan([e instanceof Error ? e.message : String(e)])
      })
      .finally(() => {
        setSibuk(false)
      })
  }
  /** `TreatyInSetValueInstallment` — larik Installment DIGANTI utuh (langkah 2). */
  const nilai = (status: '' | 'update', n = no) => {
    kirim({ aksi: 'nilai', status, angsuran, indeks: 0, installmentNo: n }, (h) => {
      setAngsuran(h.angsuran)
      setTotal(h.TotalInstallmentNP ?? [])
      setNo(h.InstallmentNo)
      setBuka(new Set())
    })
  }
  /** `SetTotalInstallment` atas SATU halaman — hanya Amount, AmountTotal, PctTotal. */
  const totalHalaman = (i: number, status: '' | 'editpercentage', halaman: Angsuran) => {
    kirim({ aksi: 'total-baris', status, angsuran: [halaman], indeks: 0 }, (h) => {
      const hasil = h.angsuran[0]
      if (hasil === undefined) return
      setAngsuran((as) =>
        as.map((a, x) =>
          x !== i
            ? a
            : {
                ...a,
                AmountTotal: hasil.AmountTotal,
                PctTotal: hasil.PctTotal,
                InstallmentList: a.InstallmentList.map((r, y) => ({ ...r, Amount: hasil.InstallmentList[y]?.Amount ?? r.Amount })),
              },
        ),
      )
    })
  }
  /**
   * `TreatyInUpdatePaymentDate(_Act)` — HANYA `PaymentDate` yang diganti,
   * atas halaman TERKINI (isian lain yang diketik selama rute menjawab tetap).
   */
  const tanggalBayar = (semua: boolean, i: number, halaman: Angsuran) => {
    const kirimAngsuran = semua ? angsuran.map((a, x) => (x === i ? halaman : a)) : [halaman]
    kirim({ aksi: semua ? 'tanggal-bayar-semua' : 'tanggal-bayar', status: '', angsuran: kirimAngsuran, indeks: 0 }, (h) => {
      setAngsuran((as) =>
        as.map((a, x) => {
          const hasil = semua ? h.angsuran[x] : x === i ? h.angsuran[0] : undefined
          if (hasil === undefined) return a
          return { ...a, InstallmentList: a.InstallmentList.map((r, y) => ({ ...r, PaymentDate: hasil.InstallmentList[y]?.PaymentDate ?? r.PaymentDate })) }
        }),
      )
    })
  }
  const ubahBaris = (i: number, j: number, ubahan: Partial<Angsuran['InstallmentList'][number]>): Angsuran => {
    const a = angsuran[i] ?? { Currency: '', AmountTotal: '', PctTotal: '', InstallmentList: [] }
    const baru = { ...a, InstallmentList: a.InstallmentList.map((r, y) => (y === j ? { ...r, ...ubahan } : r)) }
    setAngsuran((as) => as.map((x, k) => (k === i ? { ...x, InstallmentList: x.InstallmentList.map((r, y) => (y === j ? { ...r, ...ubahan } : r)) } : x)))
    return baru
  }

  return (
    <Panel judul={ANGSURAN.judul}>
      <div className="trin__angsuran-kepala">
        <div className="field">
          <label className="field__label" htmlFor="trin-angsuran-no">
            {ANGSURAN.jumlah}
          </label>
          <input
            id="trin-angsuran-no"
            className="field__input"
            type="text"
            inputMode="numeric"
            value={no}
            readOnly={terkunci}
            onChange={(e) => {
              setNo(e.target.value)
            }}
            // `change` (+ Enter) → `TreatyInSetValueInstallment` TANPA status.
            onBlur={(e) => {
              if (!terkunci && e.target.value !== '') nilai('', e.target.value)
            }}
            onKeyDown={(e) => {
              if (e.key === 'Enter' && !terkunci) nilai('', e.currentTarget.value)
            }}
          />
        </div>
        {bisaUbah && (
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={terkunci || sibuk}
            onClick={() => {
              nilai('update')
            }}
          >
            {ANGSURAN.perbaruiNilai}
          </button>
        )}
      </div>

      <div className="table-wrap">
        <table className="trin__tabel">
          <thead>
            <tr>
              <th scope="col" className="trin__buka-sel" aria-label={ANGSURAN.rincian} />
              <th scope="col">{ANGSURAN.mataUang}</th>
            </tr>
          </thead>
          <tbody>
            {angsuran.length === 0 && (
              <tr>
                <td colSpan={2}>{ANGSURAN.tanpaBaris}</td>
              </tr>
            )}
            {angsuran.map((a, i) => (
              <Fragment key={i}>
                <tr>
                  <td className="trin__buka-sel">
                    <button
                      type="button"
                      className="trin__buka"
                      aria-expanded={buka.has(i)}
                      aria-label={`${ANGSURAN.rincian} ${String(i + 1)}`}
                      onClick={() => {
                        setBuka((x) => {
                          const y = new Set(x)
                          if (y.has(i)) y.delete(i)
                          else y.add(i)
                          return y
                        })
                      }}
                    >
                      {buka.has(i) ? '▾' : '▸'}
                    </button>
                  </td>
                  <td>{a.Currency}</td>
                </tr>
                {buka.has(i) && (
                  <tr className="trin__rincian">
                    <td colSpan={2}>
                      <table className="trin__tabel">
                        <thead>
                          <tr>
                            {KOLOM_RINCIAN_ANGSURAN.map((k) => (
                              <th key={k} scope="col">
                                {k}
                              </th>
                            ))}
                          </tr>
                        </thead>
                        <tbody>
                          {a.InstallmentList.length === 0 && (
                            <tr>
                              <td colSpan={KOLOM_RINCIAN_ANGSURAN.length}>{ANGSURAN.tanpaBaris}</td>
                            </tr>
                          )}
                          {a.InstallmentList.map((r, j) => (
                            <tr key={j}>
                              <td>{r.Installment}</td>
                              <td>
                                {terkunci ? (
                                  tanggalTampil(r.DueDate)
                                ) : (
                                  <KotakTanggalKetik
                                    label={KOLOM_RINCIAN_ANGSURAN[1]}
                                    value={r.DueDate}
                                    onChange={(v) => {
                                      // `change` → TreatyInUpdatePaymentDate (halaman ini).
                                      tanggalBayar(false, i, ubahBaris(i, j, { DueDate: keSimpan(v) }))
                                    }}
                                  />
                                )}
                              </td>
                              <td>
                                {terkunci ? (
                                  r.WPC
                                ) : (
                                  <input
                                    className="field__input"
                                    type="text"
                                    inputMode="numeric"
                                    aria-label={KOLOM_RINCIAN_ANGSURAN[2]}
                                    value={r.WPC}
                                    onChange={(e) => {
                                      ubahBaris(i, j, { WPC: e.target.value })
                                    }}
                                    // `change` → TreatyInUpdatePaymentDate_Act (semua halaman).
                                    onBlur={(e) => {
                                      tanggalBayar(true, i, ubahBaris(i, j, { WPC: e.target.value }))
                                    }}
                                  />
                                )}
                              </td>
                              <td>{tanggalTampil(r.PaymentDate)}</td>
                              <td className="trin__angka">
                                {terkunci ? (
                                  selAngka(['persenShare', 2], r.InstallmentPct)
                                ) : (
                                  <input
                                    className="field__input"
                                    type="text"
                                    inputMode="decimal"
                                    aria-label={KOLOM_RINCIAN_ANGSURAN[4]}
                                    value={r.InstallmentPct}
                                    onChange={(e) => {
                                      ubahBaris(i, j, { InstallmentPct: e.target.value })
                                    }}
                                    // `change` → SetTotalInstallment(status=editpercentage).
                                    onBlur={(e) => {
                                      totalHalaman(i, 'editpercentage', ubahBaris(i, j, { InstallmentPct: e.target.value }))
                                    }}
                                  />
                                )}
                              </td>
                              <td className="trin__angka">
                                {terkunci ? (
                                  uang(r.Amount)
                                ) : (
                                  <input
                                    className="field__input"
                                    type="text"
                                    inputMode="decimal"
                                    aria-label={KOLOM_RINCIAN_ANGSURAN[5]}
                                    value={r.Amount}
                                    onChange={(e) => {
                                      ubahBaris(i, j, { Amount: e.target.value })
                                    }}
                                    // `change` → SetTotalInstallment (tanpa parameter): hanya total.
                                    onBlur={(e) => {
                                      totalHalaman(i, '', ubahBaris(i, j, { Amount: e.target.value }))
                                    }}
                                  />
                                )}
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                      <dl className="trin__angsuran-total">
                        <dt>{ANGSURAN.totalPersen}</dt>
                        <dd>{selAngka(['persenShare', 4], a.PctTotal)}</dd>
                        <dt>{ANGSURAN.total}</dt>
                        <dd>{uang(a.AmountTotal)}</dd>
                      </dl>
                    </td>
                  </tr>
                )}
              </Fragment>
            ))}
          </tbody>
        </table>
      </div>

      <div className="table-wrap trin__share-total">
        <table className="trin__tabel">
          <thead>
            <tr>
              <th scope="col">{ANGSURAN.totalAngsuran}</th>
              <th scope="col">{ANGSURAN.nilai}</th>
            </tr>
          </thead>
          <tbody>
            {total.length === 0 && (
              <tr>
                <td colSpan={2}>{ANGSURAN.tanpaBaris}</td>
              </tr>
            )}
            {total.map((t, i) => (
              <tr key={i}>
                <td>{t.Currency}</td>
                <td className="trin__angka">{uang(t.Value)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {bisaUbah && (
        <div className="trin__aksi">
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={terkunci || sibuk}
            onClick={() => {
              kirim({ aksi: 'total', status: '', angsuran, indeks: 0 }, (h) => {
                setTotal(h.TotalInstallmentNP ?? [])
              })
            }}
          >
            {ANGSURAN.perbaruiTotal}
          </button>
        </div>
      )}
      {pesan.length > 0 && (
        <p className="trin__galat" role="alert">
          {pesan.join(' · ')}
        </p>
      )}
    </Panel>
  )
}
