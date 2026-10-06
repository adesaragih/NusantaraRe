// Sub-tab Occupation baris objek FIRE (tiket 40).
//
// Port `NB FacIn\Section\OccupationList.xml`: grid `.Property.OccupationList` (Occupation ID · Occupation Name · Class
// Of Construction; Add di kepala, Delete per baris; baris dibuka = `OccupationItemFacIn_Section`): tombol Choose
// Occupation (popup `ChooseOccupation`: Search Name/ID, kolom ID · Name, Choose -> `SetDataOccupation`), Occupation ID
// dan Name baca-saja, tombol Choose Class of Construction (popup `ChooseClassofContraction`: RD
// `BrowseTableOfLimit_RD` dengan Tahun / Bizcode case dan Category baris, kolom Description, Choose ->
// `SetDataClassofConstraction`), Class of Construction baca-saja bertanda wajib.
//
// `SetDataOccupation`: OccupationId = OldID, OccupationName = Name, TableOfLimit.Category = KDRiskExposure
// "03" -> "III", "02" -> "II", "01" -> "I", lainnya kosong. `SetDataClassofConstraction`: Description + PctLimit.
//
// Keputusan agent (tiket 40): L-1 `GetLowestPctLimit_ACT` (batas persen terendah lintas objek -> parameter case untuk
// kapasitas treaty) dihitung di backend / tahap kapasitas, bukan di layar. L-2 tanda wajib Class of Construction hanya
// penanda - Save Pega (`SaveFacIn_Act`) tidak memvalidasi; penolakan = tahap Submit. L-3 popup Occupation memuat
// seluruh okupasi FIRE saat dibuka (Pega: grid tanpa saringan, 20 per halaman) dan menyaring saat mengetik.
// Choose Class of Construction hanya ditahan bila Occupation belum dipilih (Category kosong tetap dibuka).
// L-4 popup tampil MENGGANTIKAN satu sama lain (pola E-6) - tidak bertumpuk.

import { useEffect, useRef, useState } from 'react'

import { Field, Gagal, Halaman, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { cariOccupation, cariTableOfLimit, type BarisOccupation, type BarisTableOfLimit, type OkupasiObjek } from '../api'
import { FORM_OKUPASI as F, GRID_OBJEK, GRID_OKUPASI, POPUP_KONSTRUKSI, POPUP_OKUPASI, TEKS_INWARD, TEKS_OBJEK, TEKS_OKUPASI } from '../labels'

const JEDA_MS = 400

/** Baris Occupation kosong (Add). */
export const okupasiBaru = (): OkupasiObjek => ({ occupationId: '', occupationName: '', category: '', constructionClass: '', pctLimit: '' })

/** KDRiskExposure -> Category romawi (`SetDataOccupation`). */
export function kategoriDari(kd: string | undefined): string {
  return kd === '03' ? 'III' : kd === '02' ? 'II' : kd === '01' ? 'I' : ''
}

/** Choose di popup Occupation. Class of Construction tidak disentuh (sama dengan Pega). */
export function pilihOkupasi(o: OkupasiObjek, b: BarisOccupation): OkupasiObjek {
  return { ...o, occupationId: b.oldId, occupationName: b.name, category: kategoriDari(b.kdRiskExposure) }
}

/** Choose di popup Class of Construction. */
export function pilihKonstruksi(o: OkupasiObjek, b: BarisTableOfLimit): OkupasiObjek {
  return { ...o, constructionClass: b.description, pctLimit: b.pctLimit }
}

function Tampil({ label, nilai, wajib }: { label: string; nilai: string; wajib?: boolean }) {
  return (
    <div className="field">
      <span className="field__label">
        {label}
        {wajib && <span className="field__req">*</span>}
      </span>
      <div className="nbf-inward__teks">{nilai}</div>
    </div>
  )
}

function PopupOkupasi({ onTutup, onPilih }: { onTutup: () => void; onPilih: (b: BarisOccupation) => void }) {
  const [kotak, setKotak] = useState('')
  const [hasil, setHasil] = useState<BarisOccupation[] | null>(null)
  const [halaman, setHalaman] = useState(1)
  const [galat, setGalat] = useState<unknown>(null)
  const nomor = useRef(0)

  useEffect(() => {
    const n = ++nomor.current
    const jadwal = window.setTimeout(() => {
      setGalat(null)
      cariOccupation(kotak).then(
        (h) => {
          if (n === nomor.current) {
            setHasil(h.baris)
            setHalaman(1)
          }
        },
        (err: unknown) => {
          if (n === nomor.current) {
            setHasil(null)
            setGalat(err)
          }
        },
      )
    }, kotak === '' ? 0 : JEDA_MS)
    return () => window.clearTimeout(jadwal)
  }, [kotak])

  const u = POPUP_OKUPASI.ukuran
  const tampil = (hasil ?? []).slice((halaman - 1) * u, halaman * u)
  return (
    <Modal judul={F.pilihOkupasi.label} onTutup={onTutup} lebar>
      <div className="nbf-popup__kotak">
        <Field label={POPUP_OKUPASI.cari.label} value={kotak} onChange={setKotak} />
      </div>
      {hasil && hasil.length > u && <Halaman halaman={halaman} ukuran={u} total={hasil.length} onPindah={setHalaman} />}
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              {POPUP_OKUPASI.kolom.map((k) => (
                <th key={k.sel} scope="col">
                  {k.label}
                </th>
              ))}
              <th scope="col" />
            </tr>
          </thead>
          {tampil.length > 0 && (
            <tbody>
              {tampil.map((b) => (
                <tr key={b.oldId}>
                  <td>{b.oldId}</td>
                  <td>{b.name}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => onPilih(b)}>
                      {POPUP_OKUPASI.pilih.label}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          )}
        </table>
      </div>
      {hasil === null && galat === null && <Memuat />}
      {hasil !== null && hasil.length === 0 && <Kosong pesan={TEKS_OKUPASI.tanpaHasil} />}
      <Gagal galat={galat} />
    </Modal>
  )
}

function PopupKonstruksi({
  caseId,
  category,
  onTutup,
  onPilih,
}: {
  caseId: string
  category: string
  onTutup: () => void
  onPilih: (b: BarisTableOfLimit) => void
}) {
  const [hasil, setHasil] = useState<BarisTableOfLimit[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let batal = false
    cariTableOfLimit(caseId, category).then(
      (h) => {
        if (!batal) setHasil(h.baris)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [caseId, category])

  return (
    <Modal judul={F.pilihKonstruksi.label} onTutup={onTutup} lebar>
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              <th scope="col">{POPUP_KONSTRUKSI.kolom.label}</th>
              <th scope="col" />
            </tr>
          </thead>
          {hasil && hasil.length > 0 && (
            <tbody>
              {hasil.map((b, i) => (
                <tr key={`${i}-${b.description}`}>
                  <td>{b.description}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => onPilih(b)}>
                      {POPUP_KONSTRUKSI.pilih.label}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          )}
        </table>
      </div>
      {hasil === null && galat === null && <Memuat />}
      {hasil !== null && hasil.length === 0 && <Kosong pesan={TEKS_OKUPASI.tanpaHasil} />}
      <Gagal galat={galat} />
    </Modal>
  )
}

export default function SubTabOkupasi({
  caseId,
  occupations,
  ubah,
}: {
  caseId: string
  occupations: OkupasiObjek[]
  ubah: (o: OkupasiObjek[]) => void
}) {
  const [terbuka, setTerbuka] = useState<number[]>([])
  // Popup yang terbuka: jenis + indeks baris.
  const [popup, setPopup] = useState<{ jenis: 'okupasi' | 'konstruksi'; n: number } | null>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const ganti = (n: number, o: OkupasiObjek) => ubah(occupations.map((x, k) => (k === n ? o : x)))

  return (
    <div className="nbf-objek__isi">
      {pesan && <div className="alert alert--warn">{pesan}</div>}
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              <th scope="col" />
              {GRID_OKUPASI.map((k) => (
                <th key={k.sel} scope="col">
                  {k.label}
                </th>
              ))}
              <th scope="col" className="table__actions">
                <button
                  type="button"
                  className="btn btn--ghost btn--sm"
                  onClick={() => {
                    ubah([...occupations, okupasiBaru()])
                    setTerbuka((t) => [...t, occupations.length])
                  }}
                >
                  {GRID_OBJEK.tambah}
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            {occupations.length === 0 && (
              <tr>
                <td colSpan={5}>{TEKS_INWARD.kosong}</td>
              </tr>
            )}
            {occupations.map((o, n) => [
              <tr key={`b-${n}`}>
                <td>
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    aria-label={TEKS_OBJEK.bukaBaris}
                    aria-expanded={terbuka.includes(n)}
                    onClick={() => setTerbuka((t) => (t.includes(n) ? t.filter((x) => x !== n) : [...t, n]))}
                  >
                    {terbuka.includes(n) ? '▾' : '▸'}
                  </button>
                </td>
                <td>{o.occupationId}</td>
                <td>{o.occupationName}</td>
                <td>{o.constructionClass}</td>
                <td className="table__actions">
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      ubah(occupations.filter((_, x) => x !== n))
                      setTerbuka((t) => t.filter((x) => x !== n).map((x) => (x > n ? x - 1 : x)))
                    }}
                  >
                    {GRID_OBJEK.hapus}
                  </button>
                </td>
              </tr>,
              terbuka.includes(n) && (
                <tr key={`d-${n}`} className="nbf-objek__detail">
                  <td colSpan={5}>
                    <div className="nbf-objek__isi nbf-ringkas">
                      <div className="nbf-ringkas__grid">
                        <div className="nbf-opp__tombol nbf-ringkas__lebar">
                          <button type="button" className="btn btn--sm" onClick={() => setPopup({ jenis: 'okupasi', n })}>
                            {F.pilihOkupasi.label}
                          </button>
                        </div>
                        <Tampil label={F.occupationId.label} nilai={o.occupationId} />
                        <div className="nbf-ringkas__dua">
                          <Tampil label={F.occupationName.label} nilai={o.occupationName} />
                        </div>
                        <div className="nbf-opp__tombol nbf-ringkas__lebar">
                          <button
                            type="button"
                            className="btn btn--sm"
                            onClick={() => {
                              // Hanya ditahan bila Occupation belum dipilih sama sekali. Category kosong (KDRiskExposure
                              // selain 01/02/03) TETAP membuka popup - Pega pun membukanya; backend membuang saringan
                              // kategori yang kosong (perbaikan 03-10-2026).
                              if (o.occupationId === '') {
                                setPesan(TEKS_OKUPASI.pilihOkupasiDulu)
                                return
                              }
                              setPesan(null)
                              setPopup({ jenis: 'konstruksi', n })
                            }}
                          >
                            {F.pilihKonstruksi.label}
                          </button>
                        </div>
                        <div className="nbf-ringkas__dua">
                          <Tampil label={F.konstruksi.label} nilai={o.constructionClass} wajib />
                        </div>
                      </div>
                    </div>
                  </td>
                </tr>
              ),
            ])}
          </tbody>
        </table>
      </div>

      {popup?.jenis === 'okupasi' && (
        <PopupOkupasi
          onTutup={() => setPopup(null)}
          onPilih={(b) => {
            ganti(popup.n, pilihOkupasi(occupations[popup.n]!, b))
            setPopup(null)
          }}
        />
      )}
      {popup?.jenis === 'konstruksi' && (
        <PopupKonstruksi
          caseId={caseId}
          category={occupations[popup.n]!.category}
          onTutup={() => setPopup(null)}
          onPilih={(b) => {
            ganti(popup.n, pilihKonstruksi(occupations[popup.n]!, b))
            setPopup(null)
          }}
        />
      )}
    </div>
  )
}
