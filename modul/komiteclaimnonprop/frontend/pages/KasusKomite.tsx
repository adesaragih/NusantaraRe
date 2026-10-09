// Layar komite Non Prop - flow action `ViewTransferDtl`, Section `ShowTransfer` (korpus `Komite Claim Non Prop`). Bagian,
// label, dan syarat tampil disusun server VERBATIM; layar merender lalu mengirim isian keputusan lewat tombol "Submit"
// (`finishAssignment` -> `KomitePostAdjustment`). "Cancel" kembali ke pemanggil; "View Claim" (harness
// `ViewDtlClaimKmt`) membuka klaim induk Claim Non Prop hanya-baca (`PropsRute.onLihatBerkas`).
//
// Tata letak pola Komite Claim Prop (perintah work owner 09-10-2026 "ikuti tampilan klaim prop"): ringkasan di kepala;
// bagian server dikelompokkan ke kartu berjudul (`susunKomite.ts`); di bawah rincian langkah tangga ("Committe Accept
// Status") lalu kartu keputusan - berdampingan di layar lebar, bertumpuk di layar sempit.

import { useCallback, useEffect, useState } from 'react'

import { ApiFailure } from '../../../../inti/frontend/klien'
import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { bukaKasus, putuskan, type Bagian, type Grid, type Keputusan, type Layar, type Medan } from '../api'
import { KCNP, TATA_KCNP } from '../labels'
import { isianKirim, KEPUTUSAN, periksa, tampil, tampilCatatanSubjectivity, tampilSubjectivity } from '../nilai'
import { kelompokkanBagian, langkahTangga, ringkasan, type Kelompok, type KeadaanLangkah } from '../susunKomite'

/** Nilai kosong ditandai; angka dan tanggal diformat. */
function Nilai({ m }: { m: Medan }) {
  const v = tampil(m.nilai, m.jenis)
  if (v === '') return <span className="komiteclaimnonprop__kosong">—</span>
  return <>{v}</>
}

/** Pasangan label-nilai satu bagian; teks panjang selebar kartu. */
function DaftarNilai({ medan }: { medan: Medan[] }) {
  return (
    <dl className="komiteclaimnonprop__dl">
      {medan.map((m, i) => (
        <div
          key={i}
          className={
            'komiteclaimnonprop__dl-baris' +
            // teks panjang berisi = blok selebar kartu; kosong = baris biasa bertanda "—"
            (m.jenis === 'teksPanjang' && m.nilai.trim() !== '' ? ' komiteclaimnonprop__dl-baris--panjang' : '') +
            (m.jenis === 'angka' ? ' komiteclaimnonprop__angka' : '')
          }
        >
          <dt className="komiteclaimnonprop__label">{m.label}</dt>
          <dd className="komiteclaimnonprop__nilai">
            <Nilai m={m} />
          </dd>
        </div>
      ))}
    </dl>
  )
}

function TabelGrid({ g }: { g: Grid }) {
  return (
    <div className="komiteclaimnonprop__grid">
      {g.judul && <h4 className="komiteclaimnonprop__subjudul">{g.judul}</h4>}
      <div className="komiteclaimnonprop__gulir">
        <table className="komiteclaimnonprop__tabel">
          <thead>
            <tr>
              {g.kolom.map((k, i) => (
                <th key={i} className={k.jenis === 'angka' ? 'komiteclaimnonprop__angka' : undefined}>
                  {k.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {g.baris.length === 0 && (
              <tr>
                <td colSpan={g.kolom.length} className="komiteclaimnonprop__kosong">
                  {TATA_KCNP.tanpaData}
                </td>
              </tr>
            )}
            {g.baris.map((b, i) => (
              <tr key={i}>
                {g.kolom.map((k, j) => (
                  <td key={j} className={k.jenis === 'angka' ? 'komiteclaimnonprop__angka' : undefined}>
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

function BagianIsi({ b, judulKartu }: { b: Bagian; judulKartu: string }) {
  return (
    <div className="komiteclaimnonprop__bagian">
      {b.judul && b.judul !== judulKartu && <h4 className="komiteclaimnonprop__subjudul">{b.judul}</h4>}
      {b.medan && b.medan.length > 0 && <DaftarNilai medan={b.medan} />}
      {(b.grid ?? []).map((g, i) => (
        <TabelGrid key={i} g={g} />
      ))}
    </div>
  )
}

/** Satu kartu kelompok; dua bagian medan saja (Claim Analysis kiri / kanan, Payment / Bank) berdampingan. */
function Kartu({ k }: { k: Kelompok }) {
  const berdampingan = k.bagian.length > 1 && k.bagian.every((b) => (b.medan?.length ?? 0) > 0 && !b.grid)
  const isi = k.bagian.map((b) => <BagianIsi key={b.kunci} b={b} judulKartu={k.judul} />)
  return (
    <section className="panel komiteclaimnonprop__kartu" aria-label={k.judul || undefined}>
      {k.judul && <h3 className="panel__title komiteclaimnonprop__kartu-judul">{k.judul}</h3>}
      {berdampingan ? <div className="komiteclaimnonprop__dua">{isi}</div> : isi}
    </section>
  )
}

const LABEL_LANGKAH: Record<KeadaanLangkah, string> = {
  setuju: TATA_KCNP.langkahSetuju,
  tolak: TATA_KCNP.langkahTolak,
  berjalan: TATA_KCNP.langkahBerjalan,
  menunggu: TATA_KCNP.langkahMenunggu,
}

const TANDA_LANGKAH: Record<KeadaanLangkah, string> = { setuju: '✓', tolak: '✕', berjalan: '', menunggu: '' }

export default function KasusKomite({
  id,
  onKembali,
  onLihatBerkas,
}: {
  id: string
  onKembali: () => void
  /** `PropsRute.onLihatBerkas` - tombol View Claim. */
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
  if (layar === null || isi === null) return <Memuat pesan={KCNP.memuat} />

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
      () => onKembali(),
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
    ? { teks: TATA_KCNP.statusGiliranAnda, kelas: 'komiteclaimnonprop__status--anda' }
    : r.tingkat
      ? { teks: `${TATA_KCNP.statusMenunggu} ${r.tingkat.jabatan}`, kelas: 'komiteclaimnonprop__status--tunggu' }
      : { teks: TATA_KCNP.statusSelesai, kelas: 'komiteclaimnonprop__status--selesai' }
  const tombolLihat = layar.tombol.find((t) => t.aksi === 'lihat')
  const ringkas: [string, string, boolean][] = [
    [TATA_KCNP.ringkasNoKlaim, r.noKlaim, false],
    [TATA_KCNP.ringkasPolis, r.polis, false],
    [TATA_KCNP.ringkasTertanggung, r.tertanggung, false],
    [TATA_KCNP.ringkasTingkat, r.tingkat ? `${String(r.tingkat.ke)} / ${String(r.tingkat.dari)}` : '', false],
  ]

  return (
    <section className="inbox komiteclaimnonprop__akar">
      <header className="komiteclaimnonprop__kepala">
        <div className="komiteclaimnonprop__kepala-baris">
          <div>
            <h2 className="inbox__judul komiteclaimnonprop__judul">
              {layar.judul.map((j, i) => (
                <span key={i}>{j}</span>
              ))}
            </h2>
            <div className="komiteclaimnonprop__kepala-sub">
              <span className="komiteclaimnonprop__lencana">{layar.kasus.id}</span>
              <span className={'komiteclaimnonprop__status ' + status.kelas}>{status.teks}</span>
            </div>
          </div>
          {tombolLihat && (
            // View Claim = harness ViewDtlClaimKmt atas klaim induk: berkas Claim Non Prop dibuka hanya-baca di jendela
            // di atas layar ini (`PropsRute.onLihatBerkas`). Modul tidak dipasang = pesan.
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={!tombolLihat.aktif}
              title={tombolLihat.alasan}
              onClick={() => {
                setPesanLihat(
                  onLihatBerkas?.('claimnonprop', layar.kasus.klaimId) === true ? null : KCNP.berkasTakTerpasang,
                )
              }}
            >
              {tombolLihat.label}
            </button>
          )}
        </div>
        <dl className="komiteclaimnonprop__ringkas">
          {ringkas.map(([label, nilai, angka]) => (
            <div key={label} className="komiteclaimnonprop__ringkas-isi">
              <dt>{label}</dt>
              <dd className={angka ? 'komiteclaimnonprop__angka-teks' : undefined}>
                {nilai === '' ? <span className="komiteclaimnonprop__kosong">—</span> : nilai}
              </dd>
            </div>
          ))}
        </dl>
      </header>
      {pesanLihat && <div className="alert alert--info">{pesanLihat}</div>}

      <div className="komiteclaimnonprop__tata">
        <div className="komiteclaimnonprop__utama">
          {kelompok.map((k) => (
            <Kartu key={k.kunci} k={k} />
          ))}
        </div>

        <div className="komiteclaimnonprop__bawah">
          {langkah.length > 0 && (
            <section className="panel komiteclaimnonprop__kartu" aria-label={judulTangga || undefined}>
              {judulTangga && <h3 className="panel__title komiteclaimnonprop__kartu-judul">{judulTangga}</h3>}
              <ol className="komiteclaimnonprop__langkah">
                {langkah.map((l) => (
                  <li
                    key={l.urut}
                    className={'komiteclaimnonprop__langkah-isi komiteclaimnonprop__langkah--' + l.keadaan}
                  >
                    <span className="komiteclaimnonprop__langkah-tanda" aria-hidden="true">
                      {TANDA_LANGKAH[l.keadaan] || String(l.urut)}
                    </span>
                    <div className="komiteclaimnonprop__langkah-teks">
                      <div className="komiteclaimnonprop__langkah-kepala">
                        <strong>{l.jabatan}</strong>
                        <span className={'komiteclaimnonprop__status komiteclaimnonprop__status--' + l.keadaan}>
                          {LABEL_LANGKAH[l.keadaan]}
                        </span>
                      </div>
                      {l.tanggal && (
                        <span className="komiteclaimnonprop__langkah-tgl">{tampil(l.tanggal, 'tanggalJam')}</span>
                      )}
                      {l.komentar && <p className="komiteclaimnonprop__langkah-catatan">{l.komentar}</p>}
                    </div>
                  </li>
                ))}
              </ol>
            </section>
          )}

          <section
            className="panel komiteclaimnonprop__kartu komiteclaimnonprop__keputusan"
            aria-label={TATA_KCNP.keputusanAnda}
          >
            <h3 className="panel__title komiteclaimnonprop__kartu-judul">{TATA_KCNP.keputusanAnda}</h3>
            {!layar.bolehKerja && <div className="alert alert--info">{KCNP.hanyaLihat}</div>}

            <div className="komiteclaimnonprop__isian-blok">
              <span className="komiteclaimnonprop__label" id="komiteclaimnonprop-terima">
                {isian.label.acceptStatus} <span className="field__req">*</span>
              </span>
              <div
                className="komiteclaimnonprop__pilihan"
                role="radiogroup"
                aria-labelledby="komiteclaimnonprop-terima"
              >
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
                        'komiteclaimnonprop__opsi komiteclaimnonprop__opsi--' +
                        jenis +
                        (dipilih ? ' komiteclaimnonprop__opsi--dipilih' : '')
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
              <label className="komiteclaimnonprop__centang">
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
              <label className="komiteclaimnonprop__isian-blok">
                <span className="komiteclaimnonprop__label">
                  {isian.label.subjectivityNote} <span className="field__req">*</span>
                </span>
                <select
                  className="field__input"
                  value={isi.subjectivityNote}
                  disabled={kunci || !terbuka}
                  onChange={(e) => ubah('subjectivityNote', e.target.value)}
                >
                  <option value="">{KCNP.pilih}</option>
                  {isian.pilihanSubjectivityNote.map((p) => (
                    <option key={p.nilai} value={p.nilai}>
                      {p.label}
                    </option>
                  ))}
                </select>
                {salah.subjectivityNote && <span className="field__error">{salah.subjectivityNote}</span>}
              </label>
            )}
            <label className="komiteclaimnonprop__centang">
              <input
                type="checkbox"
                checked={isi.usulTutup}
                disabled={kunci || !terbuka}
                onChange={(e) => ubah('usulTutup', e.target.checked)}
              />
              {isian.label.usulTutup}
            </label>
            <label className="komiteclaimnonprop__centang">
              <input
                type="checkbox"
                checked={isi.usulCadang}
                disabled={kunci || !terbuka}
                onChange={(e) => ubah('usulCadang', e.target.checked)}
              />
              {isian.label.usulCadang}
            </label>
            <label className="komiteclaimnonprop__isian-blok">
              <span className="komiteclaimnonprop__label">
                {isian.label.comment} <span className="field__req">*</span>
              </span>
              <textarea
                className="field__input komiteclaimnonprop__area"
                value={isi.comment}
                disabled={kunci}
                onChange={(e) => ubah('comment', e.target.value)}
              />
              {salah.comment && <span className="field__error">{salah.comment}</span>}
            </label>
            {pesanGalat && <div className="alert alert--error">{pesanGalat}</div>}
            <div className="komiteclaimnonprop__aksi">
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
                      {sibuk ? KCNP.memproses : t.label}
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
