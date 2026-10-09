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

import { FieldAngka, Panel } from '../../../../inti/frontend/components/ui/dasar'
import { formatNumber, formatPersen } from '../../../../inti/frontend/lib/format'
import { hitungAngsuran, type Angsuran, type BarisAngsuranWarisan, type MasukanAngsuran, type NilaiShare } from '../api'
import { useProperti } from '../halaman'
import { ANGSURAN, KOLOM_RINCIAN_ANGSURAN } from '../labelsAngsuran'
import type { ModeForm } from '../mode'
import { angkaMurni, padankanDesimal, selAngka } from './angka'
import { TombolNavigasi } from './navigasi'
import { TataPegaBlok } from './tataPega'
import { PemicuUbah, usePemicuUbah } from './pemicuUbah'
import { saringAngka } from './saringAngka'
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

/**
 * ⭐ DESIMAL MENURUT XML — permintaan pemakai 9 Oktober 2026 ("perbaiki
 * format angka yang ada di tab installment"): `Amount` tampil mentah
 * `10213025.1974774898103125`, dan Total ber-4 desimal.
 *
 * `Section/Installments`: `% Installment` (@81277) dan `Amount` (@91479)
 * kontrol Number `pyDecimalPlaces` 2; `Total` (`AmountTotal`) 2; grid
 * `TotalInstallmentNP` kolom Value 2; `% Total` (`PctTotal`) 4.
 *
 * ⛔ `selAngka(['uang', …])` memformat dengan `DESIMAL_UANG` (4) dan sengaja
 * tidak memotong — aturan seluruh modul. Tab ini mengikuti XML-nya sendiri:
 * dibulatkan ke 2 desimal untuk TAMPIL; nilai tersimpan tidak disentuh.
 */
const DESIMAL_ANGSURAN = 2
const uang = (v: string) => padankanDesimal(formatNumber(v, DESIMAL_ANGSURAN), DESIMAL_ANGSURAN)
const persenAngsuran = (v: string) =>
  angkaMurni(v) ? padankanDesimal(formatPersen(v, DESIMAL_ANGSURAN), DESIMAL_ANGSURAN) : v

export default function TabAngsuran({
  baris,
  netPremium,
  edmState = '',
  edmJenisMaterial = '',
  mode = 'lihat',
  commencement = '',
  termination = '',
}: {
  baris: readonly BarisAngsuranWarisan[]
  /** `TreatyIn.TotalShareNetNP` TERKINI — tab Share Non-Prop. */
  netPremium: readonly NilaiShare[]
  edmState?: string
  edmJenisMaterial?: string
  mode?: ModeForm
  /** Commencement / Termination kepala TERKINI (bentuk simpan) — Due Date
   *  dibagi rata dari Commencement (keputusan pemakai 9 Oktober 2026). */
  commencement?: string
  termination?: string
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
      commencement,
      termination,
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
  // `Installment` — peristiwa `change` Pega: hanya bila jumlahnya BERUBAH
  // (Enter tetap memicu langsung). `pemicuUbah.tsx`.
  const pemicuNo = usePemicuUbah(no, terkunci ? undefined : () => {
    if (no !== '') nilai('', no)
  })
  const ubahBaris = (i: number, j: number, ubahan: Partial<Angsuran['InstallmentList'][number]>): Angsuran => {
    const a = angsuran[i] ?? { Currency: '', AmountTotal: '', PctTotal: '', InstallmentList: [] }
    const baru = { ...a, InstallmentList: a.InstallmentList.map((r, y) => (y === j ? { ...r, ...ubahan } : r)) }
    setAngsuran((as) => as.map((x, k) => (k === i ? { ...x, InstallmentList: x.InstallmentList.map((r, y) => (y === j ? { ...r, ...ubahan } : r)) } : x)))
    return baru
  }

  return (
    <Panel judul={ANGSURAN.judul}>
      {/* `Inline grid quadruple` @118888: [ `Stacked with labels left`
          @119186 Installment | Update Value ] — kolom 1 dan 2 dari empat. */}
      <TataPegaBlok tata="g4">
        <TataPegaBlok tata="kiri">
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
                setNo(saringAngka(e.target.value, true))
              }}
              // `change` (+ Enter) → `TreatyInSetValueInstallment` TANPA status.
              onFocus={pemicuNo.masuk}
              onBlur={pemicuNo.keluar}
              onKeyDown={(e) => {
                if (e.key === 'Enter' && !terkunci) nilai('', e.currentTarget.value)
              }}
            />
          </div>
        </TataPegaBlok>
        {bisaUbah && (
          <div>
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
          </div>
        )}
      </TataPegaBlok>

      <div className="table-wrap">
        <table className="trin__tabel trin__tabel--pega">
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
                    {/* Navigasi, bukan `<button>` — tetap hidup di mode lihat. */}
                    <TombolNavigasi
                      className="trin__buka"
                      terbuka={buka.has(i)}
                      label={`${ANGSURAN.rincian} ${String(i + 1)}`}
                      onKlik={() => {
                        setBuka((x) => {
                          const y = new Set(x)
                          if (y.has(i)) y.delete(i)
                          else y.add(i)
                          return y
                        })
                      }}
                    >
                      {buka.has(i) ? '▾' : '▸'}
                    </TombolNavigasi>
                  </td>
                  <td>{a.Currency}</td>
                </tr>
                {buka.has(i) && (
                  <tr className="trin__rincian">
                    <td colSpan={2}>
                      <table className="trin__tabel trin__tabel--pega">
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
                                  // `change` → TreatyInUpdatePaymentDate_Act (semua halaman) — hanya bila berubah.
                                  <PemicuUbah
                                    nilai={r.WPC}
                                    aksi={() => {
                                      tanggalBayar(true, i, ubahBaris(i, j, { WPC: r.WPC }))
                                    }}
                                  >
                                    <input
                                      className="field__input"
                                      type="text"
                                      inputMode="numeric"
                                      aria-label={KOLOM_RINCIAN_ANGSURAN[2]}
                                      value={r.WPC}
                                      onChange={(e) => {
                                        ubahBaris(i, j, { WPC: saringAngka(e.target.value, true) })
                                      }}
                                    />
                                  </PemicuUbah>
                                )}
                              </td>
                              <td>{tanggalTampil(r.PaymentDate)}</td>
                              <td className="trin__angka">
                                {terkunci ? (
                                  persenAngsuran(r.InstallmentPct)
                                ) : (
                                  // `change` → SetTotalInstallment(status=editpercentage) — hanya bila berubah.
                                  <PemicuUbah
                                    nilai={r.InstallmentPct}
                                    aksi={() => {
                                      totalHalaman(i, 'editpercentage', ubahBaris(i, j, { InstallmentPct: r.InstallmentPct }))
                                    }}
                                  >
                                    {/* Kontrol Number 2 desimal (@81277): berpemisah,
                                        huruf ditolak, mentah saat diketik. */}
                                    <FieldAngka
                                      label=""
                                      value={r.InstallmentPct}
                                      desimal={DESIMAL_ANGSURAN}
                                      onChange={(x) => {
                                        ubahBaris(i, j, { InstallmentPct: x })
                                      }}
                                    />
                                  </PemicuUbah>
                                )}
                              </td>
                              <td className="trin__angka">
                                {terkunci ? (
                                  uang(r.Amount)
                                ) : (
                                  // `change` → SetTotalInstallment (tanpa parameter): hanya total — hanya bila berubah.
                                  <PemicuUbah
                                    nilai={r.Amount}
                                    aksi={() => {
                                      totalHalaman(i, '', ubahBaris(i, j, { Amount: r.Amount }))
                                    }}
                                  >
                                    {/* Kontrol Number 2 desimal (@91479). */}
                                    <FieldAngka
                                      label=""
                                      value={r.Amount}
                                      desimal={DESIMAL_ANGSURAN}
                                      onChange={(x) => {
                                        ubahBaris(i, j, { Amount: x })
                                      }}
                                    />
                                  </PemicuUbah>
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

      {/* Wadah tanpa kepala @128066 (`Default` @128259): grid Total
          Installment lalu tombol Update Total di bawahnya. */}
      <TataPegaBlok tata="tumpuk">
        <div className="table-wrap trin__share-total">
          <table className="trin__tabel trin__tabel--pega">
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
      </TataPegaBlok>
      {pesan.length > 0 && (
        <p className="trin__galat" role="alert">
          {pesan.join(' · ')}
        </p>
      )}
    </Panel>
  )
}
