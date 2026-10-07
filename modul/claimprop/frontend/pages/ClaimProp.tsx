// Halaman awal Claim Prop - pintu masuk Pega yang BERBUKTI (`Flow/Flow_TreatyIn.xml`): Assignment2 "Outstanding Claim"
// ke worklist operator (ToCurrentOperator), Assignment1 "Input Acceptation" ke workbasket (ReasKlaimTeknik), dan
// pembuatan kasus (Start1 -> Assignment2; harness New tidak diekspor - OQ-CP-13). Kasus selesai di tab ketiga.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { buatKasus, daftarKasus, type JenisDaftar, type RingkasanKasus } from '../api'
import LayarKasus from '../components/LayarKasus'
import { CP } from '../labels'

const TAB: readonly JenisDaftar[] = ['saya', 'workbasket', 'selesai']
const LABEL_TAB: Record<JenisDaftar, string> = {
  saya: CP.tabSaya,
  workbasket: CP.tabWorkbasket,
  selesai: CP.tabSelesai,
}
const JEDA_CARI_MS = 300

export default function ClaimProp({ pelaku }: { pelaku: string }) {
  const [tab, setTab] = useState<JenisDaftar>('saya')
  const [kata, setKata] = useState('')
  const [cari, setCari] = useState('')
  const [daftar, setDaftar] = useState<RingkasanKasus[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [buka, setBuka] = useState<string | null>(null)
  const [segar, setSegar] = useState(0)

  useEffect(() => {
    const t = setTimeout(() => setCari(kata), JEDA_CARI_MS)
    return () => clearTimeout(t)
  }, [kata])

  useEffect(() => {
    if (buka) return
    setDaftar(null)
    daftarKasus(tab, cari).then(
      (d) => {
        setDaftar(d)
        setGalat(null)
      },
      (g: unknown) => setGalat(g),
    )
  }, [tab, cari, buka, segar])

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
        onKembali={() => {
          setBuka(null)
          setSegar((s) => s + 1)
        }}
      />
    )
  }

  return (
    <section className="inbox claimprop__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{CP.judul}</h2>
        <span className="toolbar__spacer" />
        <button type="button" className="btn btn--primary" onClick={baru}>
          {CP.baru}
        </button>
      </header>
      <StripTab tab={TAB} aktif={tab} onPilih={setTab} label={(t) => LABEL_TAB[t]} />
      <input
        className="claimprop__input claimprop__cari"
        placeholder={CP.cari}
        value={kata}
        onChange={(e) => setKata(e.target.value)}
      />
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={CP.memuat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={CP.kosong} />}
      {daftar !== null && daftar.length > 0 && (
        <table className="claimprop__tabel">
          <thead>
            <tr>
              <th>{CP.kolomID}</th>
              <th>{CP.kolomTahap}</th>
              <th>{CP.kolomNoClaim}</th>
              <th>{CP.kolomPolis}</th>
              <th>{CP.kolomTreaty}</th>
              <th>{CP.kolomPembuat}</th>
              <th>{CP.kolomTanggal}</th>
            </tr>
          </thead>
          <tbody>
            {daftar.map((k) => (
              <tr key={k.id} className="claimprop__baris-rinci" onClick={() => setBuka(k.id)}>
                <td>{k.id}</td>
                <td>{k.statusWork || k.label}</td>
                <td>{k.noClaim || k.claimNoTemp}</td>
                <td>{k.policyNo}</td>
                <td>{k.treatyName}</td>
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
