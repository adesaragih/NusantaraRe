// Tab Inw Fac Cedant Panels kasus FIRE (tiket 49) - port `NB FacIn\Section\InputInwardFacultativeDtl.xml` tab "Inw Fac
// Cedant Panels" (baris 60510-67588):
// - Source of Business (`.OfferFacIn.QuotationData.SobName`, tampil);
// - Share Cedant Type (`.OfferFacIn.ShareCedantType`, radio Gross "0" / Share RNM "1" = `DDL\ShareCedantType.xml`); wajib
//   bila proteksi berlaku; `SetShareOfCeding` (ShareOfCeding "100%" / PercentShare + "%") dijalankan backend saat Save;
// - % Share RNM, Total TSI RNM, Total Premi RNM (tampil - hitungan backend);
// - grid `.OfferFacIn.CedingCedantList`: Ceding · % Share · Add (`AddCedantList_act`: salin Ceding Co tab General) /
//   Delete; tombol "Choose Ceding" Pega (`Section\CedingCedant.xml`) TIDAK dibuat (work owner 05-10-2026, C-6);
// - Save -> `PUT …/cedant`; proteksi `ProtectShareCedant_Act` (pesan verbatim) diperiksa juga di layar.
//
// Detail baris = `Section\CedingCedant.xml` (gambar layar Pega work owner 05-10-2026): Ceding; % Share*
// "of" ShareOfCeding; grid Currency · TSI Cedant · Premium Cedant (`SetTSIPremiCedant_Act` langkah 4 atas
// `OfferFacIn.CurrencyList` - langkah 5 `TotalSpreadAll` ber-label `//` = mati; ralat temuan sesi c3): nilai =
// ShareCeding x dasar / 100, dasar per tipe dari backend (`dasarMataUang`), perkalian eksak di layar).
//
// Keputusan agent (tiket 49): C-1 (diralat 05-10-2026) detail baris dibuka di bawah baris grid (pola tab Object), bukan
// popup; proteksi "Total premi share cedant must be equal total premi RNM!" tetap ditunda. C-2 % Share RNM baca-saja di
// sini - diubah di tab Spreading (satu jalur hitung). C-3 proteksi hanya bila backend menyatakan `wajib` (gerbang
// `ProtectShareCedant_Act`: IsGroup == "Group" atau FlagOnGoingPolicy == 2).

import { useEffect, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { desimalSah, jumlahDesimal, pangkasNolEkor } from '../../../../inti/frontend/lib/desimal'
import { angka4 } from '../angka4'
import { ambilCedant, ambilSpreading, simpanCedant, type BarisCedant, type DasarCedant, type TampilanCedant, type TampilanSpreading } from '../api'
import { CEDANT as C, OPSI_SHARE_CEDANT_TYPE, TEKS_CEDANT as T } from '../labels'
import IsianAngka4 from './IsianAngka4'
import { lebihDari100, RingkasanSpreading } from './TabSpreading'

/**
 * Add = salin Ceding Co tab General (`cedingUmum` = `QuotationData.CedingCoList`, atau CedingCo tunggal) seperti
 * `AddCedantList_act`; tombol Choose Ceding dihapus (work owner 05-10-2026, C-6 - membatalkan C-5). Grid sudah berisi ->
 * hanya ceding General yang belum ada di grid (kode) yang ditambahkan, bukan baris kosong: tanpa Choose Ceding baris
 * kosong tidak dapat diisi.
 */
export function tambahCedant(baris: BarisCedant[], cedingUmum: BarisCedant[]): BarisCedant[] {
  const ada = new Set(baris.map((b) => b.cedingCo.trim()))
  const baru = cedingUmum.filter((c) => !ada.has(c.cedingCo.trim())).map((c) => ({ cedingCo: c.cedingCo, cedingCoName: c.cedingCoName, shareCeding: '' }))
  return [...baris, ...baru]
}

/** % Share sah: desimal > 0 dan <= 100 (`.ShareCeding<=0||.ShareCeding>100` ditolak); kosong = tidak sah. */
export function shareSah(v: string): boolean {
  const t = v.trim()
  if (!desimalSah(t) || t.startsWith('-')) return false
  return /[1-9]/.test(t) && !lebihDari100(t)
}

/**
 * Σ % Share cedant > 100 - selalu ditolak, tanpa gerbang `wajib` (work owner 05-10-2026: "ini masa bisa lebih dari
 * 100%"). Pega hanya menahannya lewat proteksi total premi di balik gerbang; pesan agent (C-7).
 */
export function totalShareLebih(baris: BarisCedant[]): boolean {
  const sah = baris.map((b) => b.shareCeding.trim()).filter((s) => desimalSah(s))
  return lebihDari100(jumlahDesimal(sah).total)
}

/** Ada % Share > 100 per baris (`SetTSIPremiCedant_Act` - ditolak backend selalu, tanpa gerbang `wajib`). */
export function adaShareLebih(baris: BarisCedant[]): boolean {
  return baris.some((b) => desimalSah(b.shareCeding.trim()) && lebihDari100(b.shareCeding.trim()))
}

/** Teks "of" detail baris: SetShareOfCeding - Gross "100%", Share RNM PercentShare + "%"; belum dipilih -> tersimpan. */
export function teksShareOfCeding(type: string, percentShare: string, tersimpan: string): string {
  if (type === '0') return '100%'
  if (type === '1') return `${desimalSah(percentShare.trim()) ? pangkasNolEkor(percentShare.trim(), 0) : percentShare.trim()}%`
  return tersimpan
}

/** a x b / 100, eksak (BigInt); "" bila salah satu bukan desimal. */
export function kaliPersen(a: string, b: string): string {
  const x = a.trim()
  const y = b.trim()
  if (!desimalSah(x) || !desimalSah(y)) return ''
  const pecah = (t: string): [bigint, number] => {
    const negatif = t.startsWith('-')
    const [bulat = '0', pecahan = ''] = (negatif || t.startsWith('+') ? t.slice(1) : t).split('.')
    const n = BigInt((bulat === '' ? '0' : bulat) + pecahan)
    return [negatif ? -n : n, pecahan.length]
  }
  const [na, sa] = pecah(x)
  const [nb, sb] = pecah(y)
  const n = na * nb
  const skala = sa + sb + 2
  const negatif = n < 0n
  const digit = (negatif ? -n : n).toString().padStart(skala + 1, '0')
  return (negatif ? '-' : '') + digit.slice(0, digit.length - skala) + '.' + digit.slice(digit.length - skala)
}

/** CurrencyList satu cedant (`SetTSIPremiCedant_Act`): % Share tidak sah atau tipe belum dipilih -> kosong. */
export function mataUangCedant(share: string, type: string, dasar: DasarCedant[]): { currency: string; tsi: string; premium: string }[] {
  if (!shareSah(share) || (type !== '0' && type !== '1')) return []
  return dasar.map((d) => ({
    currency: d.currency,
    tsi: kaliPersen(share, type === '1' ? d.tsiRnm : d.tsiGross),
    premium: kaliPersen(share, type === '1' ? d.premiumRnm : d.premiumGross),
  }))
}

/** Pesan proteksi `ProtectShareCedant_Act` (urutan langkah Pega, tanpa pesan kembar). */
export function galatCedant(baris: BarisCedant[]): string[] {
  if (baris.length === 0) return [T.listKosong]
  const pesan: string[] = []
  if (baris.some((b) => !shareSah(b.shareCeding))) pesan.push(T.shareTidakSah)
  if (baris.some((b) => b.cedingCoName.trim() === '')) pesan.push(T.cedingKosong)
  return pesan
}

export default function TabCedant({ caseId }: { caseId: string }) {
  const [v, setV] = useState<TampilanCedant | null>(null)
  const [sp, setSp] = useState<TampilanSpreading | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [cobaSimpan, setCobaSimpan] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)
  // Add tidak menambah apa pun (semua Ceding Co General sudah ada / General kosong).
  const [tanpaTambahan, setTanpaTambahan] = useState(false)
  const [terbuka, setTerbuka] = useState<number[]>([])

  useEffect(() => {
    let batal = false
    ambilCedant(caseId).then(
      (h) => {
        if (!batal) setV(h)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    // Total per mata uang (Dtl layout S62 `IsFire`: grid `TotalSpreadAll.pxResults`) = ringkasan tab Spreading.
    ambilSpreading(caseId).then(
      (h) => {
        if (!batal) setSp(h)
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

  const ubah = (baru: Partial<TampilanCedant>) => {
    setV({ ...v, ...baru })
    setTersimpan(false)
  }
  const ubahBaris = (n: number, b: Partial<BarisCedant>) => ubah({ cedant: v.cedant.map((x, k) => (k === n ? { ...x, ...b } : x)) })
  const typeKosong = v.wajib && v.shareCedantType === ''
  const pesan = v.wajib ? galatCedant(v.cedant) : []
  const shareLebih = adaShareLebih(v.cedant)
  const totalLebih = totalShareLebih(v.cedant)
  const ofCeding = teksShareOfCeding(v.shareCedantType, v.percentShare, v.shareOfCeding)
  const hapusBaris = (n: number) => {
    ubah({ cedant: v.cedant.filter((_, k) => k !== n) })
    setTanpaTambahan(false)
    setTerbuka((t) => t.filter((x) => x !== n).map((x) => (x > n ? x - 1 : x)))
  }

  async function simpan() {
    if (v === null) return
    setCobaSimpan(true)
    if (typeKosong || pesan.length > 0 || shareLebih || totalLebih) return
    setMenyimpan(true)
    setGalat(null)
    try {
      setV(await simpanCedant(caseId, { shareCedantType: v.shareCedantType, cedant: v.cedant }))
      setTersimpan(true)
      setCobaSimpan(false)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  return (
    <div className="nbf-objek nbf-cedant">
      <Gagal galat={galat} />
      {tersimpan && <div className="alert alert--ok">{T.tersimpan}</div>}

      {/* Blok atas dirapikan (work owner 05-10-2026 "biar lebih user friendly"): Source of Business + Share Cedant Type
          sebaris, tiga angka RNM sebagai kartu kecil. */}
      <section className="nbf-lapis nbf-cedant__atas">
        <div className="nbf-cedant__info">
          <div className="nbf-cedant__medan">
            <span className="nbf-cedant__label">{C.sob}</span>
            <span className="nbf-cedant__nilai">{v.sobName || '—'}</span>
          </div>
          <div className="nbf-cedant__medan">
            <span className="nbf-cedant__label">
              {C.shareCedantType}
              {v.wajib && <span className="field__req">*</span>}
            </span>
            <div className="nbf-cedant__tipe" role="radiogroup" aria-label={C.shareCedantType}>
              {OPSI_SHARE_CEDANT_TYPE.map((o) => (
                <label key={o.value} className={'nbf-cedant__opsi' + (v.shareCedantType === o.value ? ' nbf-cedant__opsi--aktif' : '')}>
                  <input
                    type="radio"
                    name="shareCedantType"
                    value={o.value}
                    checked={v.shareCedantType === o.value}
                    onChange={() => ubah({ shareCedantType: o.value })}
                  />
                  {o.label}
                </label>
              ))}
            </div>
            {cobaSimpan && typeKosong && <div className="field__error">{T.typeWajib}</div>}
          </div>
        </div>
        <dl className="nbf-cedant__angka">
          <div>
            <dt>{C.percentShare}</dt>
            <dd>{v.percentShare === '' ? '—' : angka4(v.percentShare)}</dd>
          </div>
          <div>
            <dt>
              {C.totalTsi} <span className="nbf-cedant__uang">{T.mataUangTotal}</span>
            </dt>
            <dd>{v.totalTsiRnm === '' ? '—' : angka4(v.totalTsiRnm)}</dd>
          </div>
          <div>
            <dt>
              {C.totalPremi} <span className="nbf-cedant__uang">{T.mataUangTotal}</span>
            </dt>
            <dd>{v.totalPremiRnm === '' ? '—' : angka4(v.totalPremiRnm)}</dd>
          </div>
        </dl>
        {/* Total RNM dijumlah dalam IDR lewat kurs (C-8); mata uang tanpa kurs disebut di sini. */}
        {v.kursTidakAda.length > 0 && (
          <div className="alert alert--warn">
            {T.kursTidakAda}
            {v.kursTidakAda.join(', ')}
          </div>
        )}
      </section>

      <section className="nbf-lapis nbf-cedant__grid">
        <div className="nbf-cov-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                <th scope="col" />
                <th scope="col">{C.kolomCeding}</th>
                <th scope="col">{C.kolomShare}</th>
                <th scope="col" className="table__actions">
                  <button
                    type="button"
                    className="btn btn--sm"
                    onClick={() => {
                      const hasil = tambahCedant(v.cedant, v.cedingUmum)
                      setTanpaTambahan(hasil.length === v.cedant.length)
                      ubah({ cedant: hasil })
                      setTerbuka((t) => [...t, ...hasil.slice(v.cedant.length).map((_, k) => v.cedant.length + k)])
                    }}
                  >
                    {C.tambah}
                  </button>
                </th>
              </tr>
            </thead>
            <tbody>
              {v.cedant.length === 0 && (
                <tr>
                  <td colSpan={4}>{T.kosong}</td>
                </tr>
              )}
              {v.cedant.map((b, n) => {
                const buka = terbuka.includes(n)
                const uang = mataUangCedant(b.shareCeding, v.shareCedantType, v.dasarMataUang)
                return [
                  <tr key={`b-${n}`} className={buka ? 'nbf-baris--terbuka' : undefined}>
                    <td>
                      <button
                        type="button"
                        className={'btn btn--ghost btn--sm nbf-buka' + (buka ? ' nbf-buka--terbuka' : '')}
                        aria-label={C.bukaBaris}
                        aria-expanded={buka}
                        onClick={() => setTerbuka((t) => (buka ? t.filter((x) => x !== n) : [...t, n]))}
                      >
                        ▸
                      </button>
                    </td>
                    <td>{b.cedingCoName || '—'}</td>
                    <td className="nbf-angka">{b.shareCeding === '' ? '—' : angka4(b.shareCeding)}</td>
                    <td className="table__actions">
                      <button type="button" className="btn btn--sm" onClick={() => hapusBaris(n)}>
                        {C.hapus}
                      </button>
                    </td>
                  </tr>,
                  buka && (
                    <tr key={`d-${n}`} className="nbf-objek__detail">
                      <td colSpan={4}>
                        <div className="nbf-cedant__rinci">
                          <div className="nbf-cedant__baris">
                            <span className="nbf-cedant__label">{C.kolomCeding}</span>
                            <span className="nbf-cedant__nilai">{b.cedingCoName || '—'}</span>
                          </div>
                          <div className="nbf-cedant__baris">
                            <label className="nbf-cedant__label" htmlFor={`share-cedant-${n}`}>
                              {C.kolomShare}
                              <span className="field__req">*</span>
                            </label>
                            <IsianAngka4
                              id={`share-cedant-${n}`}
                              value={b.shareCeding}
                              onChange={(x) => ubahBaris(n, { shareCeding: x })}
                            />
                            <span className="nbf-cedant__dari">
                              {C.dari} <strong>{ofCeding}</strong>
                            </span>
                          </div>
                          <div className="nbf-cov-wrap">
                            <table className="nbf-tabel">
                              <thead>
                                <tr>
                                  {C.kolomMataUang.map((k) => (
                                    <th key={k} scope="col">
                                      {k}
                                    </th>
                                  ))}
                                </tr>
                              </thead>
                              <tbody>
                                {uang.length === 0 && (
                                  <tr>
                                    <td colSpan={3}>{T.kosong}</td>
                                  </tr>
                                )}
                                {uang.map((u) => (
                                  <tr key={u.currency}>
                                    <td>{u.currency}</td>
                                    <td className="nbf-angka">{angka4(u.tsi)}</td>
                                    <td className="nbf-angka">{angka4(u.premium)}</td>
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
        {tanpaTambahan && <div className="muted">{v.cedingUmum.length === 0 ? T.generalKosong : T.sudahSemua}</div>}
        {shareLebih && <div className="field__error">{T.nilaiLebih}</div>}
        {totalLebih && <div className="field__error">{T.totalLebih}</div>}
        {cobaSimpan &&
          pesan.map((p) => (
            <div key={p} className="field__error">
              {p}
            </div>
          ))}
      </section>

      {sp !== null && <RingkasanSpreading v={sp} hanyaMataUang />}

      <div className="nbf-objek__kaki">
        <button type="button" className="btn btn--primary" onClick={() => void simpan()} disabled={menyimpan}>
          {menyimpan ? T.menyimpan : T.simpan}
        </button>
      </div>
    </div>
  )
}
