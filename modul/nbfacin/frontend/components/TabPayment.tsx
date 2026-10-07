// Tab Payment kasus FIRE (tiket 51) - port `NB FacIn\Section\InputInwardFacultativeDtl.xml` tab "Payment" non-Life ->
// `PaymentCurrencyList` -> `InputDtlPayment_FacIn`, tata letak = gambar layar Pega `DDL\Blok Currency, Payment,
// Installment.JPG`:
// - Installment* (`.Policy.Payment.Installment`) diubah (change / Enter) -> `POST …/payment/hitung` aksi "installment"
//   (`FillPaymentInstallment`: daftar cicilan dibangun ulang, % rata, Due Date bawaan);
// - grid per mata uang Currency · Premium RNM · Total Payment, dibuka -> "Installment Data Information": No · Due Date* ·
//   % Installment · Premium RNM · Deduction · Deduction 2 · Total Brokerage Fee · PPH n % · PPN n % · Discount · Total
//   Payment (Policy Fee hanya `IsHealthIndividu` - tidak untuk FIRE);
// - % Installment diubah -> aksi "persen" (`SetValidateInstallment_Act`); Σ % per mata uang diperiksa juga di layar
//   (pesan verbatim);
// - Save -> `PUT …/payment` (`IsThereAnyObjectLocation_Act` + `ProtectPremiPolicy_Act`; pesan proteksi dari backend).
// Seluruh hitungan uang di backend (desimal eksak); layar hanya menampilkan.
//
// Keputusan agent (tiket 51): P-1 isian dihitung saat lepas fokus / Enter (event change / enter Pega). P-2 Due Date
// diubah tanpa memanggil server (Pega: tanpa aksi); tersimpan lewat Save.

import { useEffect, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { desimalSah, jumlahDesimal } from '../../../../inti/frontend/lib/desimal'
import { angka4 } from '../angka4'
import {
  ambilPayment,
  hitungPayment,
  simpanPayment,
  type AksiPayment,
  type CicilanPayment,
  type TampilanPayment,
} from '../api'
import { PAYMENT as P, TEKS_FORM_OPPORTUNITY, TEKS_PAYMENT as T } from '../labels'
import IsianAngka4 from './IsianAngka4'
import { kepalaPajak } from './TabDeduction'
import { lebihDari100 } from './TabSpreading'
import TanggalDMY from './TanggalDMY'

const angka = (v: string | undefined) => angka4(v ?? '')

/** Pesan Σ % Installment satu mata uang (`SetValidateInstallment_Act`): > 100, < 100, atau "" bila tepat 100. */
export function pesanPersenCicilan(cicilan: CicilanPayment[]): string {
  if (cicilan.length === 0) return ''
  const sah = cicilan.map((c) => c.installmentPercentage.trim()).filter((s) => desimalSah(s))
  const total = jumlahDesimal(sah).total
  if (lebihDari100(total)) return T.persenLebih
  const [bulat = '0', pecahan = ''] = total.replace(/^\+/, '').split('.')
  const tepat100 = !bulat.startsWith('-') && BigInt(bulat === '' ? '0' : bulat) === 100n && !/[1-9]/.test(pecahan)
  return tepat100 ? '' : T.persenKurang
}

/** Jumlah cicilan sah: bilangan bulat 1..99 (`pyMax=2`). */
export const installmentSah = (v: string) => /^[1-9]\d?$/.test(v.trim())

export default function TabPayment({ caseId }: { caseId: string }) {
  const [v, setV] = useState<TampilanPayment | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [terbuka, setTerbuka] = useState<string[]>([])
  const [menghitung, setMenghitung] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)

  useEffect(() => {
    let batal = false
    ambilPayment(caseId).then(
      (h) => {
        if (!batal) setV(h)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [caseId])

  if (v === null) return galat ? <Gagal galat={galat} /> : <Memuat />

  const ubah = (baru: Partial<TampilanPayment>) => {
    setV({ ...v, ...baru })
    setTersimpan(false)
  }
  const ubahCicilan = (cur: string, n: number, c: Partial<CicilanPayment>) =>
    ubah({
      mataUang: v.mataUang.map((m) => (m.currency !== cur ? m : { ...m, cicilan: m.cicilan.map((x, k) => (k === n ? { ...x, ...c } : x)) })),
    })

  async function hitung(aksi: AksiPayment, currency: string, idx: number, isi: TampilanPayment) {
    setMenghitung(true)
    setGalat(null)
    try {
      setV(await hitungPayment(caseId, aksi, currency, idx, isi))
    } catch (err) {
      setGalat(err)
    } finally {
      setMenghitung(false)
    }
  }

  async function simpan() {
    if (v === null) return
    setMenyimpan(true)
    setGalat(null)
    try {
      setV(await simpanPayment(caseId, v))
      setTersimpan(true)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  /** Lepas fokus / Enter = event change / enter Pega (P-1). */
  const saatSelesai = (jalankan: () => void) => ({
    onBlur: jalankan,
    onKeyDown: (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'Enter') jalankan()
    },
  })

  return (
    <div className="nbf-objek nbf-pay">
      <Gagal galat={galat} />
      {tersimpan && <div className="alert alert--ok">{T.tersimpan}</div>}
      {v.pesan.map((p) => (
        <div key={p} className="alert alert--warn">
          {p}
        </div>
      ))}

      <h4 className="nbf-pay__judul">{P.judul}</h4>
      <div className="nbf-pay__installment">
        <label htmlFor="pay-installment">
          {P.installment}
          <span className="field__req">*</span>
        </label>
        <input
          id="pay-installment"
          className="nbf-angka"
          inputMode="numeric"
          maxLength={2}
          value={v.installment}
          onChange={(e) => ubah({ installment: e.target.value })}
          {...saatSelesai(() => {
            if (installmentSah(v.installment)) void hitung('installment', '', 0, v)
          })}
        />
        {menghitung && <span className="muted">{T.menghitung}</span>}
        {v.installment.trim() !== '' && !installmentSah(v.installment) && <div className="field__error">{T.installmentTidakSah}</div>}
      </div>

      <section className="nbf-pay__luar">
        <div className="nbf-cov-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                <th scope="col" />
                {P.kolomMataUang.map((k) => (
                  <th key={k} scope="col">
                    {k}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {v.mataUang.length === 0 && (
                <tr>
                  <td colSpan={4}>{T.kosong}</td>
                </tr>
              )}
              {v.mataUang.map((m) => {
                const buka = terbuka.includes(m.currency)
                const pesan = pesanPersenCicilan(m.cicilan)
                return [
                  <tr key={m.currency} className={buka ? 'nbf-baris--terbuka' : undefined}>
                    <td>
                      <button
                        type="button"
                        className={'btn btn--ghost btn--sm nbf-buka' + (buka ? ' nbf-buka--terbuka' : '')}
                        aria-label={T.bukaBaris}
                        aria-expanded={buka}
                        onClick={() => setTerbuka((t) => (buka ? t.filter((x) => x !== m.currency) : [...t, m.currency]))}
                      >
                        ▸
                      </button>
                    </td>
                    <td>{m.currency}</td>
                    <td className="nbf-angka">{angka(m.premium)}</td>
                    <td className="nbf-angka">{angka(m.sumTotalPayment)}</td>
                  </tr>,
                  buka && (
                    <tr key={`${m.currency}-d`} className="nbf-objek__detail">
                      <td colSpan={4}>
                        <div className="nbf-pay__cicilan">
                          <h5>{P.judulCicilan}</h5>
                          {pesan !== '' && <div className="field__error">{pesan}</div>}
                          <div className="nbf-cov-wrap">
                            <table className="nbf-tabel">
                              <thead>
                                <tr>
                                  <th scope="col">{P.no}</th>
                                  <th scope="col">
                                    {P.dueDate}
                                    <span className="field__req">*</span>
                                  </th>
                                  <th scope="col">{P.persen}</th>
                                  <th scope="col">{P.premium}</th>
                                  <th scope="col">{P.deduction}</th>
                                  <th scope="col">{P.deduction2}</th>
                                  <th scope="col">{P.brokerage}</th>
                                  <th scope="col">{kepalaPajak(P.pph, v.tarifPph)}</th>
                                  <th scope="col">{kepalaPajak(P.ppn, v.tarifPpn)}</th>
                                  <th scope="col">{P.discount}</th>
                                  <th scope="col">{P.totalPayment}</th>
                                </tr>
                              </thead>
                              <tbody>
                                {m.cicilan.length === 0 && (
                                  <tr>
                                    <td colSpan={11}>{T.kosong}</td>
                                  </tr>
                                )}
                                {m.cicilan.map((c, n) => (
                                  <tr key={n}>
                                    <td>{c.installmentNo}</td>
                                    <td>
                                      <TanggalDMY
                                        label={`${P.dueDate} ${m.currency} ${c.installmentNo}`}
                                        value={c.dueDate}
                                        onChange={(d) => ubahCicilan(m.currency, n, { dueDate: d })}
                                        labelKalender={TEKS_FORM_OPPORTUNITY.kalender}
                                        pesanFormat={TEKS_FORM_OPPORTUNITY.formatTanggal}
                                      />
                                    </td>
                                    <td>
                                      <IsianAngka4
                                        ariaLabel={`${P.persen} ${m.currency} ${c.installmentNo}`}
                                        value={c.installmentPercentage}
                                        onChange={(x) => ubahCicilan(m.currency, n, { installmentPercentage: x })}
                                        onSelesai={() => void hitung('persen', m.currency, n, v)}
                                      />
                                    </td>
                                    <td className="nbf-angka">{angka(c.premium)}</td>
                                    <td className="nbf-angka">{angka(c.riCommision)}</td>
                                    <td className="nbf-angka">{angka(c.deduction2)}</td>
                                    <td className="nbf-angka">{angka(c.brokerageFee)}</td>
                                    <td className="nbf-angka">{angka(c.pph)}</td>
                                    <td className="nbf-angka">{angka(c.ppn)}</td>
                                    <td className="nbf-angka">{angka(c.discount)}</td>
                                    <td className="nbf-angka">{angka(c.paymentTotal)}</td>
                                  </tr>
                                ))}
                              </tbody>
                            </table>
                          </div>
                        </div>
                      </td>
                    </tr>
                  ),
                ]
              })}
            </tbody>
          </table>
        </div>
      </section>

      <div className="nbf-objek__kaki">
        <button type="button" className="btn btn--primary" onClick={() => void simpan()} disabled={menyimpan || menghitung}>
          {menyimpan ? T.menyimpan : T.simpan}
        </button>
      </div>
    </div>
  )
}
