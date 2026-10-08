// Layar komite - flow action `ViewTransferDtl`, Section `ShowTransfer` wajah TT 2 (ADJUSTMENT). Bagian, label, dan
// syarat tampil disusun server VERBATIM; layar merender lalu mengirim isian keputusan lewat tombol "Submit"
// (`finishAssignment` -> `KomitePost`). "Cancel" kembali ke daftar kerja; "View more details" nonaktif (harness
// `ViewClaimFormKomite` tidak diekspor).

import { useCallback, useEffect, useState } from 'react'

import { ApiFailure } from '../../../../inti/frontend/klien'
import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { bukaKasus, putuskan, type Bagian, type Grid, type Keputusan, type Layar, type Medan } from '../api'
import { KCP } from '../labels'
import { isianKirim, periksa, tampil, tampilCatatanSubjectivity, tampilSubjectivity } from '../nilai'

function NilaiMedan({ m }: { m: Medan }) {
  const v = tampil(m.nilai, m.jenis)
  const kelas =
    'field__input field__input--readonly komiteclaimprop__nilai' +
    (m.jenis === 'angka' ? ' komiteclaimprop__angka' : '') +
    (m.jenis === 'teksPanjang' ? ' komiteclaimprop__panjang' : '')
  return <div className={kelas}>{v === '' ? ' ' : v}</div>
}

function DaftarMedan({ medan }: { medan: Medan[] }) {
  return (
    <div className="komiteclaimprop__medan">
      {medan.map((m, i) => (
        <div key={i} className="komiteclaimprop__baris-medan">
          <span className="komiteclaimprop__label">{m.label}</span>
          <NilaiMedan m={m} />
        </div>
      ))}
    </div>
  )
}

function TabelGrid({ g }: { g: Grid }) {
  return (
    <div className="komiteclaimprop__grid">
      <h4 className="komiteclaimprop__subjudul">{g.judul}</h4>
      <table className="komiteclaimprop__tabel">
        <thead>
          <tr>
            {g.kolom.map((k, i) => (
              <th key={i} className={k.jenis === 'angka' ? 'komiteclaimprop__angka' : undefined}>
                {k.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {g.baris.length === 0 && (
            <tr>
              <td colSpan={g.kolom.length} className="muted">
                —
              </td>
            </tr>
          )}
          {g.baris.map((b, i) => (
            <tr key={i}>
              {g.kolom.map((k, j) => (
                <td key={j} className={k.jenis === 'angka' ? 'komiteclaimprop__angka' : undefined}>
                  {tampil(b[k.properti] ?? '', k.jenis)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function BagianView({ b }: { b: Bagian }) {
  return (
    <div className="komiteclaimprop__bagian">
      {b.judul && <h3 className="panel__title">{b.judul}</h3>}
      {b.medan && b.medan.length > 0 && <DaftarMedan medan={b.medan} />}
      {(b.grid ?? []).map((g, i) => (
        <TabelGrid key={i} g={g} />
      ))}
      {b.sel && (
        <table className="komiteclaimprop__tabel">
          <tbody>
            {b.sel.map((baris, i) => (
              <tr key={i}>
                {baris.map((s, j) =>
                  i === 0 || (s.label !== '' && s.nilai === '') ? (
                    <th key={j}>{s.label}</th>
                  ) : (
                    <td key={j} className={s.jenis === 'angka' ? 'komiteclaimprop__angka' : undefined}>
                      {tampil(s.nilai, s.jenis)}
                    </td>
                  ),
                )}
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}

export default function KasusKomite({ id, onKembali }: { id: string; onKembali: (pesan?: string[]) => void }) {
  const [layar, setLayar] = useState<Layar | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [isi, setIsi] = useState<Keputusan | null>(null)
  const [salah, setSalah] = useState<Partial<Record<keyof Keputusan, string>>>({})
  const [sibuk, setSibuk] = useState(false)

  useEffect(() => {
    let aktif = true
    bukaKasus(id).then(
      (l) => {
        if (!aktif) return
        setLayar(l)
        setIsi(l.isian.nilai)
      },
      (g: unknown) => {
        if (aktif) setGalat(g)
      },
    )
    return () => {
      aktif = false
    }
  }, [id])

  const ubah = useCallback(<K extends keyof Keputusan>(k: K, v: Keputusan[K]) => {
    setIsi((x) => (x ? { ...x, [k]: v } : x))
  }, [])

  if (galat !== null && layar === null) return <Gagal galat={galat} />
  if (layar === null || isi === null) return <Memuat pesan={KCP.memuat} />

  const { isian } = layar
  const terbuka = isian.terbuka && layar.bolehKerja
  const kunci = !layar.bolehKerja || sibuk

  const kirim = () => {
    const p = periksa(isi, isian.terbuka, layar.kasus.komiteLoop)
    setSalah(p)
    if (Object.keys(p).length > 0) return
    setSibuk(true)
    setGalat(null)
    putuskan(id, isianKirim(isi, isian.terbuka)).then(
      () => {
        onKembali()
      },
      (g: unknown) => {
        setSibuk(false)
        setGalat(g)
      },
    )
  }

  const pesanGalat = galat instanceof ApiFailure ? (galat.detail.message ?? null) : null
  const [kiri, kanan, ...lain] = layar.bagian

  return (
    <section className="inbox komiteclaimprop__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul komiteclaimprop__judul">
          {layar.judul.map((j, i) => (
            <span key={i}>{j}</span>
          ))}
        </h2>
        <span className="muted">{layar.kasus.id}</span>
      </header>
      {!layar.bolehKerja && <div className="alert alert--info">{KCP.hanyaLihat}</div>}

      <div className="panel komiteclaimprop__dua">
        {kiri && <BagianView b={kiri} />}
        <div>
          {kanan && <BagianView b={kanan} />}
          {layar.tombol
            .filter((t) => t.aksi === 'lihat')
            .map((t) => (
              <button key={t.aksi} type="button" className="btn btn--ghost btn--sm" disabled title={t.alasan}>
                {t.label}
              </button>
            ))}
        </div>
      </div>
      {lain.map((b) => (
        <div key={b.kunci} className="panel">
          <BagianView b={b} />
        </div>
      ))}

      <div className="panel komiteclaimprop__isian">
        <label className="komiteclaimprop__baris-medan">
          <span className="komiteclaimprop__label">
            {isian.label.acceptStatus} <span className="field__req">*</span>
          </span>
          <select
            className="field__input"
            value={isi.acceptStatus}
            disabled={kunci}
            onChange={(e) => ubah('acceptStatus', e.target.value)}
          >
            <option value="">{KCP.pilih}</option>
            {isian.pilihanTerima.map((p) => (
              <option key={p.nilai} value={p.nilai}>
                {p.label}
              </option>
            ))}
          </select>
          {salah.acceptStatus && <span className="field__error">{salah.acceptStatus}</span>}
        </label>
        {tampilSubjectivity(isi) && (
          <label className="komiteclaimprop__centang">
            <input
              type="checkbox"
              checked={isi.isSubjectivity}
              disabled={kunci || !terbuka}
              onChange={(e) => ubah('isSubjectivity', e.target.checked)}
            />
            {isian.label.isSubjectivity}
          </label>
        )}
        {tampilCatatanSubjectivity(isi) && (
          <label className="komiteclaimprop__baris-medan">
            <span className="komiteclaimprop__label">
              {isian.label.subjectivityNote} <span className="field__req">*</span>
            </span>
            <input
              className="field__input"
              value={isi.subjectivityNote}
              disabled={kunci || !terbuka}
              onChange={(e) => ubah('subjectivityNote', e.target.value)}
            />
            {salah.subjectivityNote && <span className="field__error">{salah.subjectivityNote}</span>}
            {salah.isSubjectivity && <span className="field__error">{salah.isSubjectivity}</span>}
          </label>
        )}
        <label className="komiteclaimprop__centang">
          <input
            type="checkbox"
            checked={isi.usulTutup}
            disabled={kunci || !terbuka}
            onChange={(e) => ubah('usulTutup', e.target.checked)}
          />
          {isian.label.usulTutup}
        </label>
        <label className="komiteclaimprop__centang">
          <input
            type="checkbox"
            checked={isi.usulCadang}
            disabled={kunci || !terbuka}
            onChange={(e) => ubah('usulCadang', e.target.checked)}
          />
          {isian.label.usulCadang}
        </label>
        <label className="komiteclaimprop__baris-medan">
          <span className="komiteclaimprop__label">
            {isian.label.comment} <span className="field__req">*</span>
          </span>
          <textarea
            className="field__input komiteclaimprop__area"
            value={isi.comment}
            disabled={kunci}
            onChange={(e) => ubah('comment', e.target.value)}
          />
          {salah.comment && <span className="field__error">{salah.comment}</span>}
        </label>
        {pesanGalat && <div className="alert alert--error">{pesanGalat}</div>}
        <div className="komiteclaimprop__aksi">
          {layar.tombol
            .filter((t) => t.aksi !== 'lihat')
            .map((t) =>
              t.aksi === 'batal' ? (
                <button
                  key={t.aksi}
                  type="button"
                  className="btn btn--ghost"
                  disabled={sibuk}
                  onClick={() => onKembali()}
                >
                  {t.label}
                </button>
              ) : (
                <button
                  key={t.aksi}
                  type="button"
                  className="btn btn--primary"
                  disabled={!t.aktif || sibuk}
                  title={t.alasan}
                  onClick={kirim}
                >
                  {sibuk ? KCP.memproses : t.label}
                </button>
              ),
            )}
        </div>
      </div>
    </section>
  )
}
