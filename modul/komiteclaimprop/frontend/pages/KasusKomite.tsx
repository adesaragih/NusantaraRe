// Layar komite - flow action `ViewTransferDtl`, Section `ShowTransfer` wajah TT 2 (ADJUSTMENT). Bagian, label, dan
// syarat tampil disusun server VERBATIM; layar merender lalu mengirim isian keputusan lewat tombol "Submit"
// (`finishAssignment` -> `KomitePost`). "Cancel" kembali ke pemanggil; "View more details" membuka klaim induk Claim
// Prop hanya-baca (`PropsRute.onLihatBerkas`).
//
// Tata letak (permintaan work owner 09-10-2026 "layout komite diperbaiki, lebih enak dilihat dan userfriendly"):
// ringkasan di kepala; bagian server dikelompokkan ke kartu berjudul (`susunKomite.ts`); di BAWAH rincian (work owner
// "ini biarkan tetep dibawah") langkah tangga ("Committe Accept Status") lalu kartu keputusan - berdampingan di layar
// lebar, bertumpuk di layar sempit (urutan ShowTransfer: grid tangga sebelum isian keputusan).

import { useCallback, useEffect, useState } from 'react'

import { ApiFailure } from '../../../../inti/frontend/klien'
import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { bukaKasus, putuskan, type Bagian, type Grid, type Keputusan, type Layar, type Medan } from '../api'
import { KCP, TATA_KCP } from '../labels'
import { isianKirim, KEPUTUSAN, periksa, tampil, tampilCatatanSubjectivity, tampilSubjectivity } from '../nilai'
import { kelompokkanBagian, langkahTangga, ringkasan, type Kelompok, type KeadaanLangkah } from '../susunKomite'

/** Nilai kosong ditandai; angka dan tanggal diformat. */
function Nilai({ m }: { m: Medan }) {
  const v = tampil(m.nilai, m.jenis)
  if (v === '') return <span className="komiteclaimprop__kosong">—</span>
  return <>{v}</>
}

/** Pasangan label-nilai satu bagian; teks panjang selebar kartu. */
function DaftarNilai({ medan }: { medan: Medan[] }) {
  return (
    <dl className="komiteclaimprop__dl">
      {medan.map((m, i) => (
        <div
          key={i}
          className={
            'komiteclaimprop__dl-baris' +
            // teks panjang berisi = blok selebar kartu; kosong = baris biasa bertanda "—"
            (m.jenis === 'teksPanjang' && m.nilai.trim() !== '' ? ' komiteclaimprop__dl-baris--panjang' : '') +
            (m.jenis === 'angka' ? ' komiteclaimprop__angka' : '')
          }
        >
          <dt className="komiteclaimprop__label">{m.label}</dt>
          <dd className="komiteclaimprop__nilai">
            <Nilai m={m} />
          </dd>
        </div>
      ))}
    </dl>
  )
}

function TabelGrid({ g }: { g: Grid }) {
  return (
    <div className="komiteclaimprop__grid">
      {g.judul && <h4 className="komiteclaimprop__subjudul">{g.judul}</h4>}
      <div className="komiteclaimprop__gulir">
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
                <td colSpan={g.kolom.length} className="komiteclaimprop__kosong">
                  {TATA_KCP.tanpaData}
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
    </div>
  )
}

/** Tabel bebas (blok deductible 6 x 4): baris pertama dan sel berlabel tanpa nilai = kepala. */
function TabelSel({ sel }: { sel: Medan[][] }) {
  return (
    <div className="komiteclaimprop__gulir">
      <table className="komiteclaimprop__tabel">
        <tbody>
          {sel.map((baris, i) => (
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
    </div>
  )
}

function BagianIsi({ b, judulKartu }: { b: Bagian; judulKartu: string }) {
  return (
    <div className="komiteclaimprop__bagian">
      {b.judul && b.judul !== judulKartu && <h4 className="komiteclaimprop__subjudul">{b.judul}</h4>}
      {b.medan && b.medan.length > 0 && <DaftarNilai medan={b.medan} />}
      {(b.grid ?? []).map((g, i) => (
        <TabelGrid key={i} g={g} />
      ))}
      {b.sel && <TabelSel sel={b.sel} />}
    </div>
  )
}

/** Satu kartu kelompok; dua bagian medan saja (Claim Analysis kiri / kanan, Payment / Bank) berdampingan. */
function Kartu({ k }: { k: Kelompok }) {
  const berdampingan = k.bagian.length > 1 && k.bagian.every((b) => (b.medan?.length ?? 0) > 0 && !b.grid && !b.sel)
  const isi = k.bagian.map((b) => <BagianIsi key={b.kunci} b={b} judulKartu={k.judul} />)
  return (
    <section className="panel komiteclaimprop__kartu" aria-label={k.judul || undefined}>
      {k.judul && <h3 className="panel__title komiteclaimprop__kartu-judul">{k.judul}</h3>}
      {berdampingan ? <div className="komiteclaimprop__dua">{isi}</div> : isi}
    </section>
  )
}

const LABEL_LANGKAH: Record<KeadaanLangkah, string> = {
  setuju: TATA_KCP.langkahSetuju,
  tolak: TATA_KCP.langkahTolak,
  berjalan: TATA_KCP.langkahBerjalan,
  menunggu: TATA_KCP.langkahMenunggu,
}

const TANDA_LANGKAH: Record<KeadaanLangkah, string> = { setuju: '✓', tolak: '✕', berjalan: '', menunggu: '' }

export default function KasusKomite({
  id,
  onKembali,
  onLihatBerkas,
}: {
  id: string
  onKembali: (pesan?: string[]) => void
  /** `PropsRute.onLihatBerkas` - tombol View more details. */
  onLihatBerkas?: (modul: string, id: string) => boolean
}) {
  const [layar, setLayar] = useState<Layar | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [isi, setIsi] = useState<Keputusan | null>(null)
  const [salah, setSalah] = useState<Partial<Record<keyof Keputusan, string>>>({})
  const [sibuk, setSibuk] = useState(false)
  const [pesanLihat, setPesanLihat] = useState<string | null>(null)

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
    const p = periksa(isi, isian.terbuka)
    setSalah(p)
    if (Object.keys(p).length > 0) return
    setSibuk(true)
    setGalat(null)
    putuskan(id, isianKirim(isi, isian.terbuka)).then(
      (h) => {
        onKembali(h.galatDokumen ? [h.galatDokumen] : undefined)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalat(g)
      },
    )
  }

  const pesanGalat = galat instanceof ApiFailure ? (galat.detail.message ?? null) : null
  const { kelompok, tangga } = kelompokkanBagian(layar.bagian)
  const r = ringkasan(layar)
  const langkah = langkahTangga(layar.kasus.tangga ?? [])
  const judulTangga = tangga?.grid?.[0]?.judul ?? ''
  const status = layar.bolehKerja
    ? { teks: TATA_KCP.statusGiliranAnda, kelas: 'komiteclaimprop__status--anda' }
    : r.tingkat
      ? { teks: `${TATA_KCP.statusMenunggu} ${r.tingkat.jabatan}`, kelas: 'komiteclaimprop__status--tunggu' }
      : { teks: TATA_KCP.statusSelesai, kelas: 'komiteclaimprop__status--selesai' }
  const tombolLihat = layar.tombol.find((t) => t.aksi === 'lihat')
  const ringkas: [string, string, boolean][] = [
    [TATA_KCP.ringkasNoKlaim, r.noKlaim, false],
    [TATA_KCP.ringkasPolis, r.polis, false],
    [TATA_KCP.ringkasTertanggung, r.tertanggung, false],
    [
      TATA_KCP.ringkasAdjustment,
      r.adjustment ? `${r.adjustment.mataUang} ${tampil(r.adjustment.nilai, 'angka')}` : '',
      true,
    ],
    [TATA_KCP.ringkasTingkat, r.tingkat ? `${String(r.tingkat.ke)} / ${String(r.tingkat.dari)}` : '', false],
  ]

  return (
    <section className="inbox komiteclaimprop__akar">
      <header className="komiteclaimprop__kepala">
        <div className="komiteclaimprop__kepala-baris">
          <div>
            <h2 className="inbox__judul komiteclaimprop__judul">
              {layar.judul.map((j, i) => (
                <span key={i}>{j}</span>
              ))}
            </h2>
            <div className="komiteclaimprop__kepala-sub">
              <span className="komiteclaimprop__lencana">{layar.kasus.id}</span>
              <span className={'komiteclaimprop__status ' + status.kelas}>{status.teks}</span>
            </div>
          </div>
          {tombolLihat && (
            // View more details = harness ViewClaimFormKomite atas klaim induk: berkas Claim Prop dibuka hanya-baca di
            // jendela di atas layar ini (`PropsRute.onLihatBerkas`). Modul tidak dipasang = pesan.
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={!tombolLihat.aktif}
              title={tombolLihat.alasan}
              onClick={() => {
                setPesanLihat(
                  onLihatBerkas?.('claimprop', layar.kasus.klaimId) === true ? null : KCP.berkasTakTerpasang,
                )
              }}
            >
              {tombolLihat.label}
            </button>
          )}
        </div>
        <dl className="komiteclaimprop__ringkas">
          {ringkas.map(([label, nilai, angka]) => (
            <div key={label} className="komiteclaimprop__ringkas-isi">
              <dt>{label}</dt>
              <dd className={angka ? 'komiteclaimprop__angka-teks' : undefined}>
                {nilai === '' ? <span className="komiteclaimprop__kosong">—</span> : nilai}
              </dd>
            </div>
          ))}
        </dl>
      </header>
      {pesanLihat && <div className="alert alert--info">{pesanLihat}</div>}

      <div className="komiteclaimprop__tata">
        <div className="komiteclaimprop__utama">
          {kelompok.map((k) => (
            <Kartu key={k.kunci} k={k} />
          ))}
        </div>

        <div className="komiteclaimprop__bawah">
          {langkah.length > 0 && (
            <section className="panel komiteclaimprop__kartu" aria-label={judulTangga || undefined}>
              {judulTangga && <h3 className="panel__title komiteclaimprop__kartu-judul">{judulTangga}</h3>}
              <ol className="komiteclaimprop__langkah">
                {langkah.map((l) => (
                  <li key={l.urut} className={'komiteclaimprop__langkah-isi komiteclaimprop__langkah--' + l.keadaan}>
                    <span className="komiteclaimprop__langkah-tanda" aria-hidden="true">
                      {TANDA_LANGKAH[l.keadaan] || String(l.urut)}
                    </span>
                    <div className="komiteclaimprop__langkah-teks">
                      <div className="komiteclaimprop__langkah-kepala">
                        <strong>{l.jabatan}</strong>
                        <span className={'komiteclaimprop__status komiteclaimprop__status--' + l.keadaan}>
                          {LABEL_LANGKAH[l.keadaan]}
                        </span>
                      </div>
                      {l.tanggal && (
                        <span className="komiteclaimprop__langkah-tgl">{tampil(l.tanggal, 'tanggalJam')}</span>
                      )}
                      {l.komentar && <p className="komiteclaimprop__langkah-catatan">{l.komentar}</p>}
                    </div>
                  </li>
                ))}
              </ol>
            </section>
          )}

          <section
            className="panel komiteclaimprop__kartu komiteclaimprop__keputusan"
            aria-label={TATA_KCP.keputusanAnda}
          >
            <h3 className="panel__title komiteclaimprop__kartu-judul">{TATA_KCP.keputusanAnda}</h3>
            {!layar.bolehKerja && <div className="alert alert--info">{KCP.hanyaLihat}</div>}

            <div className="komiteclaimprop__isian-blok">
              <span className="komiteclaimprop__label" id="komiteclaimprop-terima">
                {isian.label.acceptStatus} <span className="field__req">*</span>
              </span>
              <div className="komiteclaimprop__pilihan" role="radiogroup" aria-labelledby="komiteclaimprop-terima">
                {isian.pilihanTerima.map((p) => {
                  const dipilih = isi.acceptStatus === p.nilai
                  const jenis = p.nilai === KEPUTUSAN.tolak ? 'tolak' : 'setuju'
                  return (
                    <button
                      key={p.nilai}
                      type="button"
                      role="radio"
                      aria-checked={dipilih}
                      disabled={kunci}
                      className={
                        'komiteclaimprop__opsi komiteclaimprop__opsi--' +
                        jenis +
                        (dipilih ? ' komiteclaimprop__opsi--dipilih' : '')
                      }
                      onClick={() => ubah('acceptStatus', p.nilai)}
                    >
                      {p.label}
                    </button>
                  )
                })}
              </div>
              {salah.acceptStatus && <span className="field__error">{salah.acceptStatus}</span>}
            </div>

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
              <label className="komiteclaimprop__isian-blok">
                <span className="komiteclaimprop__label">
                  {isian.label.subjectivityNote} <span className="field__req">*</span>
                </span>
                <select
                  className="field__input"
                  value={isi.subjectivityNote}
                  disabled={kunci || !terbuka}
                  onChange={(e) => ubah('subjectivityNote', e.target.value)}
                >
                  <option value="">{KCP.pilih}</option>
                  {isian.pilihanSubjectivityNote.map((p) => (
                    <option key={p.nilai} value={p.nilai}>
                      {p.label}
                    </option>
                  ))}
                </select>
                {salah.subjectivityNote && <span className="field__error">{salah.subjectivityNote}</span>}
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
            <label className="komiteclaimprop__isian-blok">
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
          </section>
        </div>
      </div>
    </section>
  )
}
