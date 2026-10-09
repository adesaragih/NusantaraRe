// Halaman awal Claim Non Prop - pola Claim Prop (prompt Claim Non Prop §7; disalin dari
// `modul/claimprop/frontend/pages/ClaimProp.tsx`, bukan impor): dua tab, Process dan Resolve. Tab Process bawaan =
// worklist pembuat (Assignment2 "Outstanding Claim", `Flow_TreatyIn` ToCurrentOperator - tanpa cek workbasket); switch
// Teknik = Assignment1 "Input Acceptation" (workbasket ReasKlaimTeknik), dapat dinyalakan hanya oleh anggota workbasket
// itu. Add Claim (Start1 -> Assignment2; harness New tidak diekspor - OQ-CNP-23) hanya saat switch Teknik mati. Tabel
// komite di bawah inbox = tahap 2 (`komiteclaimnonprop`).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { ambilHak, buatKasus, daftarKasus, type HakPelaku, type RingkasanKasus } from '../api'
import LayarKasus from '../components/LayarKasus'
import { CNP } from '../labels'
import { bolehTambahKlaim, jenisDaftar, switchTeknikAktif, TAB_INBOX, type TabInbox } from './inbox'

const LABEL_TAB: Record<TabInbox, string> = { proses: CNP.tabProses, selesai: CNP.tabResolve }
const JEDA_CARI_MS = 300

export default function ClaimNonProp({
  pelaku,
  onLihatBerkas,
  bukaKasus,
  onBeranda,
}: {
  pelaku: string
  /** `PropsRute.onLihatBerkas` - tombol View polis (jendela Modal NB / EDM Treaty In, OQ-CNP-13). */
  onLihatBerkas?: (modul: string, id: string) => boolean
  /** `PropsRute.bukaKasus` - satu berkas dibuka langsung; `hanyaLihat` = tampilan saja. */
  bukaKasus?: { id: string; ketuk: number; hanyaLihat?: boolean }
  /** `PropsRute.onBeranda` - Back berkas yang dibuka lewat `bukaKasus` kembali ke pemanggil. */
  onBeranda?: () => void
}) {
  const [tab, setTab] = useState<TabInbox>('proses')
  const [teknik, setTeknik] = useState(false)
  const [hak, setHak] = useState<HakPelaku | null>(null)
  const [kata, setKata] = useState('')
  const [cari, setCari] = useState('')
  const [daftar, setDaftar] = useState<RingkasanKasus[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [buka, setBuka] = useState<string | null>(null)
  const [segar, setSegar] = useState(0)
  const [dariLuar, setDariLuar] = useState(false)
  const [hanyaLihat, setHanyaLihat] = useState(false)

  const ketukBuka = bukaKasus?.ketuk
  const idBuka = bukaKasus?.id
  const lihatBuka = bukaKasus?.hanyaLihat === true
  useEffect(() => {
    if (ketukBuka !== undefined && idBuka !== undefined) {
      setBuka(idBuka)
      setDariLuar(true)
      setHanyaLihat(lihatBuka)
    }
  }, [ketukBuka, idBuka, lihatBuka])

  useEffect(() => {
    const t = setTimeout(() => setCari(kata), JEDA_CARI_MS)
    return () => clearTimeout(t)
  }, [kata])

  useEffect(() => {
    ambilHak().then(setHak, () => setHak(null))
  }, [])

  const bolehTeknik = switchTeknikAktif(hak)
  const teknikNyala = teknik && bolehTeknik

  useEffect(() => {
    if (buka) return
    // Jawaban basi (tab / switch / kata cari lama) diabaikan - tidak menimpa daftar terbaru.
    let aktif = true
    setDaftar(null)
    daftarKasus(jenisDaftar(tab, teknikNyala), cari).then(
      (d) => {
        if (!aktif) return
        setDaftar(d)
        setGalat(null)
      },
      (g: unknown) => {
        if (aktif) setGalat(g)
      },
    )
    return () => {
      aktif = false
    }
  }, [tab, teknikNyala, cari, buka, segar])

  const baru = useCallback(() => {
    buatKasus().then(
      (l) => setBuka(l.kasus.id),
      (g: unknown) => setGalat(g),
    )
  }, [])

  if (buka) {
    return (
      <LayarKasus
        id={buka}
        pelaku={pelaku}
        onLihatBerkas={onLihatBerkas}
        hanyaLihat={hanyaLihat}
        onKembali={() => {
          setBuka(null)
          setHanyaLihat(false)
          setSegar((s) => s + 1)
          if (dariLuar) {
            setDariLuar(false)
            onBeranda?.()
          }
        }}
      />
    )
  }

  return (
    <section className="inbox claimnonprop__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{CNP.judul}</h2>
      </header>
      <StripTab tab={TAB_INBOX} aktif={tab} onPilih={setTab} label={(t) => LABEL_TAB[t]} />

      <div className="toolbar">
        {tab === 'proses' && (
          <span className="claimnonprop__saklar-bungkus" title={bolehTeknik ? undefined : CNP.teknikTanpaHak}>
            <span className="claimnonprop__saklar-label">{CNP.inbox}</span>
            <button
              type="button"
              role="switch"
              aria-checked={teknikNyala}
              aria-labelledby="claimnonprop-saklar-teknik"
              className={'claimnonprop__saklar' + (teknikNyala ? ' claimnonprop__saklar--nyala' : '')}
              disabled={!bolehTeknik}
              onClick={() => setTeknik(!teknikNyala)}
            >
              <span className="claimnonprop__saklar-kenop" />
            </button>
            <span id="claimnonprop-saklar-teknik" className="claimnonprop__saklar-label">
              {CNP.wbTeknik}
            </span>
          </span>
        )}
        <input
          className="field__input claimnonprop__cari"
          type="search"
          aria-label={CNP.cari}
          placeholder={CNP.cari}
          value={kata}
          onChange={(e) => setKata(e.target.value)}
        />
        <span className="toolbar__spacer" />
        {bolehTambahKlaim(tab, teknikNyala) && (
          <button type="button" className="btn btn--primary" onClick={baru}>
            {CNP.tambahKlaim}
          </button>
        )}
      </div>

      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={CNP.memuat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={CNP.kosong} />}
      {daftar !== null && daftar.length > 0 && (
        <table className="inbox__tabel claimnonprop__tabel">
          <thead>
            <tr>
              <th>{CNP.kolomID}</th>
              <th>{CNP.kolomTahap}</th>
              <th>{CNP.kolomNoClaim}</th>
              <th>{CNP.kolomPolis}</th>
              <th>{CNP.kolomTreaty}</th>
              <th>{CNP.kolomPembuat}</th>
              <th>{CNP.kolomTanggal}</th>
            </tr>
          </thead>
          <tbody>
            {daftar.map((k) => (
              <tr key={k.id} className="inbox__baris" onClick={() => setBuka(k.id)}>
                <td>{k.id}</td>
                <td>{k.statusWork || k.label}</td>
                <td>{k.noClaim || k.claimNoTemp || <span className="muted">—</span>}</td>
                <td>{k.policyNo || <span className="muted">—</span>}</td>
                <td>{k.treatyName || <span className="muted">—</span>}</td>
                <td>{k.pembuatNama}</td>
                <td>{k.tglCreate}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}
