// Pop-up harness Claim Fac In yang isinya BUKAN tata halaman kasus: CauseofLoss_Harness, CatastrofeList, Outstanding,
// RetroList_Harnness, dan grid hasil pencarian ViewPolis (di dalam modal Choose Polis). Judul = pyWindowName tombol
// pembuka / pyLabel rule; judul kolom VERBATIM section XML; baris dibaca dari server (`GET .../pilihan/{jenis}`) dan
// pilihan dikirim balik sebagai aksi (server membaca ulang barisnya). Pola
// `modul/claimnonprop/frontend/components/Popup.tsx` (disalin, bukan impor).

import { useEffect, useState, type ReactNode } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import {
  pilihanKasus,
  type Baris,
  type BarisKatastrofe,
  type BarisPolisCari,
  type BarisSebab,
  type RingkasanOutstanding,
} from '../api'
import { tampilTanggal } from '../ketikTanggal'
import { CFI } from '../labels'
import { tampilAngka } from '../nilai'
import { kunciPolis } from './kunciBaris'
import type { JenisPopup } from './sesudahAksi'
import { jumlahHalaman, potongHalaman } from './susun'

const JUDUL: Record<JenisPopup, string> = {
  sebab: CFI.popSebab,
  katastrofe: CFI.popKatastrofe,
  outstanding: CFI.popOutstanding,
  retro: CFI.popRetro,
}

/** Baris per halaman grid hasil ViewPolis (pyGridPaginator); server membatasi 500 baris (`models.BatasCariPolis`). */
const PER_HALAMAN_POLIS = 10

/** Sel angka rata kanan. */
function Angka({ v }: { v: string | undefined }) {
  return <td className="claimfacin__angka">{tampilAngka(v ?? '')}</td>
}

/** Tabel hanya-baca: judul kolom VERBATIM + baris. */
function Tabel({ kolom, baris, kaki }: { kolom: readonly string[]; baris: ReactNode[]; kaki?: ReactNode }) {
  return (
    <table className="claimfacin__tabel">
      <thead>
        <tr>
          {kolom.map((k, i) => (
            <th key={i}>{k}</th>
          ))}
        </tr>
      </thead>
      <tbody>
        {baris.length === 0 ? (
          <tr>
            <td className="muted" colSpan={kolom.length}>
              —
            </td>
          </tr>
        ) : (
          baris
        )}
      </tbody>
      {kaki}
    </table>
  )
}

/**
 * Grid hasil SearchPolis_act di modal Choose Polis (ViewPolis): tautan Policy Number = `CopyNB_Act` (param
 * "PolicyNo|Prodke", dibaca ulang server).
 */
export function HasilPolis({
  baris,
  sibuk,
  onPilih,
}: {
  baris: readonly BarisPolisCari[]
  sibuk: boolean
  onPilih: (param: string) => void
}) {
  const [hal, setHal] = useState(1)
  const nHal = jumlahHalaman(baris.length, PER_HALAMAN_POLIS)
  const halIni = Math.min(hal, nHal)
  return (
    <div className="claimfacin__grid">
      {baris.length > PER_HALAMAN_POLIS && (
        <div className="claimfacin__pager">
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={halIni <= 1}
            onClick={() => setHal(halIni - 1)}
            aria-label="Previous"
          >
            ‹
          </button>
          <span>
            Page {halIni} of {nHal}
          </span>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={halIni >= nHal}
            onClick={() => setHal(halIni + 1)}
            aria-label="Next"
          >
            ›
          </button>
        </div>
      )}
      <Tabel
        kolom={CFI.polisKolom}
        baris={potongHalaman(baris, PER_HALAMAN_POLIS, halIni).map((b) => (
          <tr key={kunciPolis(b)}>
            <td>
              <button
                type="button"
                className="claimfacin__tautan"
                disabled={sibuk}
                onClick={() => onPilih(`${b.policyNo}|${b.prodke}`)}
              >
                {b.policyNo}
              </button>
            </td>
            <td>{b.customerName}</td>
            <td>{b.sourceOfBusinessName}</td>
            <td>{b.cedingCoName}</td>
            <td>{b.qq}</td>
            <td>{tampilTanggal(b.startDateTime, 'tanggal')}</td>
            <td>{tampilTanggal(b.endDateTime, 'tanggal')}</td>
            <td>{b.prodke}</td>
            <td>{b.businessName}</td>
          </tr>
        ))}
      />
    </div>
  )
}

export default function Popup({
  id,
  jenis,
  indeks,
  pelaku,
  sts,
  retro,
  peringatan,
  onPilih,
  onTutup,
}: {
  id: string
  jenis: JenisPopup
  /** Outstanding: baris objek (1..n) tombol "Outstanding Summary". */
  indeks: number
  pelaku: string
  /** Catastrophe / NonKatastrofeType klaim (form katastrofe baru, read-only). */
  sts: { sts: string; non: string }
  /** `OfferFacIn.FacRetroList` halaman polis (RetroList_SC). */
  retro: readonly Baris[]
  /** Galat / pesan validasi aksi pilihan (mis. Note katastrofe kosong). */
  peringatan?: ReactNode
  onPilih: (aksi: string, param: string) => void
  onTutup: () => void
}) {
  const [cari, setCari] = useState('')
  const [data, setData] = useState<unknown>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [formBaru, setFormBaru] = useState(false)
  const [catatan, setCatatan] = useState('')

  useEffect(() => {
    if (jenis === 'retro') return
    // Jawaban basi diabaikan: permintaan saringan lama tidak boleh menimpa hasil saringan terbaru.
    let aktif = true
    const t = setTimeout(
      () => {
        pilihanKasus<unknown>(id, jenis, { cari, indeks: jenis === 'outstanding' ? indeks : 0 }).then(
          (d) => {
            if (!aktif) return
            setData(d)
            setGalat(null)
          },
          (g: unknown) => {
            if (aktif) setGalat(g)
          },
        )
      },
      jenis === 'outstanding' ? 0 : 300,
    )
    return () => {
      aktif = false
      clearTimeout(t)
    }
  }, [id, jenis, cari, indeks])

  const pakaiCari = jenis === 'sebab' || jenis === 'katastrofe'
  let isi
  if (jenis === 'retro') {
    // RetroList_SC: Reinsurer Name | Share (%) | R/I Comm - halaman polis yang dibaca ulang (OQ-CFI-16).
    isi = (
      <Tabel
        kolom={CFI.retroKolom}
        baris={retro.map((b, i) => (
          <tr key={i}>
            <td>{b.ReinsurerName}</td>
            <Angka v={b.PctShareAllObj} />
            <Angka v={b.RiCommAllObj} />
          </tr>
        ))}
      />
    )
  } else if (galat) isi = <Gagal galat={galat} />
  else if (data === null) isi = <Memuat pesan={CFI.memuat} />
  else if (jenis === 'sebab') {
    isi = (
      <>
        <div className="claimfacin__label">{CFI.judulSebab}</div>
        <table className="claimfacin__tabel">
          <thead>
            <tr>
              <th>{CFI.kolomSebab}</th>
              <th />
            </tr>
          </thead>
          <tbody>
            {(data as BarisSebab[]).map((b) => (
              <tr key={b.id}>
                <td>{b.description}</td>
                <td>
                  <button
                    type="button"
                    className="btn btn--sm btn--primary"
                    onClick={() => onPilih('GetNameCauseofLoss', b.id)}
                  >
                    {CFI.pilih}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </>
    )
  } else if (jenis === 'katastrofe') {
    isi = formBaru ? (
      <div className="form-grid">
        <label className="field">
          <span className="field__label">Catastrophe</span>
          <input className="field__input field__input--readonly" readOnly tabIndex={-1} value={sts.sts} />
        </label>
        {sts.sts === 'Non-Catastrophe' && (
          <label className="field">
            <span className="field__label">&nbsp;</span>
            <input className="field__input field__input--readonly" readOnly tabIndex={-1} value={sts.non} />
          </label>
        )}
        <label className="field field--lebar">
          <span className="field__label">{CFI.note}</span>
          <textarea
            className="field__input claimfacin__area"
            value={catatan}
            onChange={(e) => setCatatan(e.target.value)}
          />
        </label>
        <label className="field">
          <span className="field__label">{CFI.userInput}</span>
          <input className="field__input field__input--readonly" readOnly tabIndex={-1} value={pelaku} />
        </label>
        <div className="field--lebar claimfacin__tombol claimfacin__tombol--akhir">
          <button type="button" className="btn btn--ghost btn--sm" onClick={() => setFormBaru(false)}>
            {CFI.cancel}
          </button>
          <button type="button" className="btn btn--sm btn--primary" onClick={() => onPilih('SaveCatasrtope', catatan)}>
            {CFI.save}
          </button>
        </div>
      </div>
    ) : (
      <>
        <div className="claimfacin__grid-alat">
          <button type="button" className="btn btn--sm" onClick={() => setFormBaru(true)}>
            {CFI.addNew}
          </button>
        </div>
        <table className="claimfacin__tabel">
          <thead>
            <tr>
              <th />
              <th />
              <th>{CFI.note}</th>
              <th />
            </tr>
          </thead>
          <tbody>
            {(data as BarisKatastrofe[]).map((b) => (
              <tr key={b.id}>
                <td>{b.stsKatastrofe}</td>
                <td>{b.nonKatastrofeType}</td>
                <td>{b.note}</td>
                <td>
                  <button
                    type="button"
                    className="btn btn--sm btn--primary"
                    onClick={() => onPilih('SetCatastrope', b.id)}
                  >
                    {CFI.pilih}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </>
    )
  } else {
    // Outstanding_SC: kolom nama per lini (`labelNama`), nilai akseptasi negatif, kurs standar, total dalam IDR.
    const r = data as RingkasanOutstanding
    isi = (
      <Tabel
        kolom={[
          CFI.osNoClaim,
          CFI.osPolis,
          CFI.osAksep,
          r.labelNama,
          CFI.osCoverage,
          CFI.osCurrency,
          CFI.osTSI,
          CFI.osNilai,
          CFI.osKurs,
          CFI.osKonversi,
        ]}
        baris={(r.baris ?? []).map((b, i) => (
          <tr key={i}>
            <td>{b.noClaim}</td>
            <td>{b.noPolis}</td>
            <td>{b.acceptedNo}</td>
            <td>{b.nama}</td>
            <td>{b.coverage}</td>
            <td>{b.currency}</td>
            <Angka v={b.tsiRnm} />
            <Angka v={b.nilai} />
            <Angka v={b.kurs} />
            <Angka v={b.konversi} />
          </tr>
        ))}
        kaki={
          <tfoot>
            <tr>
              <td colSpan={8} className="claimfacin__label">
                {CFI.osTotal}
              </td>
              <td>{r.mataUang}</td>
              <Angka v={r.total} />
            </tr>
          </tfoot>
        }
      />
    )
  }

  return (
    <Modal judul={JUDUL[jenis]} onTutup={onTutup} labelBatal={CFI.tutup} lebar>
      {peringatan}
      {pakaiCari && !formBaru && (
        <input
          className="field__input claimfacin__cari"
          placeholder={CFI.cariPopup}
          aria-label={CFI.cariPopup}
          value={cari}
          onChange={(e) => setCari(e.target.value)}
        />
      )}
      {isi}
    </Modal>
  )
}
