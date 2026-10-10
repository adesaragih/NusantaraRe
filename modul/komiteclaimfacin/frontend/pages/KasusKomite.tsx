// Layar komite Fac In - flow action `ViewTransferDtl`, Section `ShowTransfer` (korpus `Komite Claim FacIn`) beserta
// rinciannya `SpreadingDetail` / `DetailAdjustmentFac` (expand pane) dan pop-up `ShowRetro` ("View Retro"). Ubin, bagian,
// label, dan syarat tampil disusun server VERBATIM; layar merender lalu mengirim isian keputusan lewat tombol "Submit"
// (`finishAssignment` -> `KomitePostAct`). "Cancel" kembali ke inbox Claim Fac In; "View more details" membuka klaim
// induk Claim Fac In hanya-baca (`PropsRute.onLihatBerkas` -> `bukaKasus.hanyaLihat`, `GET .../kasus/{id}?lihat=1`).
//
// Tata letak pola Komite Claim Prop / Non Prop (keputusan work owner 09-10-2026, disalin bukan impor): judul + ubin
// ringkasan di kepala; satu kartu berjudul per bagian server (`susunKomite.ts`), teks komite kartu terakhir; DI BAWAH
// semua rincian langkah tangga ("List of Committee") lalu kartu keputusan - berdampingan di layar lebar, bertumpuk di
// layar sempit.

import { Fragment, useCallback, useEffect, useRef, useState } from 'react'

import { ApiFailure } from '../../../../inti/frontend/klien'
import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import {
  bukaKasus,
  GalatValidasiKomite,
  putuskan,
  type Bagian,
  type Grid,
  type Keputusan,
  type KolomGrid,
  type Layar,
  type Medan,
} from '../api'
import { KCFI, TATA_KCFI } from '../labels'
import { isianKirim, KEPUTUSAN, periksa, tampil, tercentang, usulTerbuka } from '../nilai'
import {
  adaRincian,
  barisGrid,
  keadaanKasus,
  kunciModal,
  langkahTangga,
  rincianBaris,
  susunBagian,
  type KeadaanLangkah,
} from '../susunKomite'

/** Pembuka pop-up sel `tautan` (kunci `Layar.modal`). */
type BukaModal = (kunci: string) => void

/** Kotak centang hanya-baca (pxCheckbox). */
function Centang({ v, label }: { v: string; label: string }) {
  return (
    <input
      type="checkbox"
      className="komiteclaimfacin__centang-baca"
      checked={tercentang(v)}
      readOnly
      disabled
      aria-label={label || undefined}
    />
  )
}

/** Nilai kosong ditandai; angka dan tanggal diformat. */
function Nilai({ m }: { m: Medan }) {
  if (m.jenis === 'centang') return <Centang v={m.nilai} label={m.label} />
  const v = tampil(m.nilai, m.jenis)
  if (v === '') return <span className="komiteclaimfacin__kosong">—</span>
  return <>{v}</>
}

/** Pasangan label-nilai satu bagian; teks panjang selebar kartu. */
function DaftarNilai({ medan }: { medan: Medan[] }) {
  return (
    <dl className="komiteclaimfacin__dl">
      {medan.map((m, i) => (
        <div
          key={i}
          className={
            'komiteclaimfacin__dl-baris' +
            // teks panjang berisi = blok selebar kartu; kosong = baris biasa bertanda "—"
            (m.jenis === 'teksPanjang' && m.nilai.trim() !== '' ? ' komiteclaimfacin__dl-baris--panjang' : '') +
            (m.jenis === 'angka' ? ' komiteclaimfacin__angka' : '')
          }
        >
          <dt className="komiteclaimfacin__label">{m.label}</dt>
          <dd className="komiteclaimfacin__nilai">
            <Nilai m={m} />
          </dd>
        </div>
      ))}
    </dl>
  )
}

/** Baris kaki grid (Total ...): pasangan label-nilai di bawah tabel. */
function KakiGrid({ kaki }: { kaki: Medan[] }) {
  return (
    <dl className="komiteclaimfacin__kaki">
      {kaki.map((m, i) => (
        <div key={i} className="komiteclaimfacin__kaki-isi">
          <dt>{m.label}</dt>
          <dd className={m.jenis === 'angka' ? 'komiteclaimfacin__angka-teks' : undefined}>
            <Nilai m={m} />
          </dd>
        </div>
      ))}
    </dl>
  )
}

/** Kelas sel menurut jenis kolom. */
function kelasSel(jenis: string): string | undefined {
  if (jenis === 'angka') return 'komiteclaimfacin__angka'
  if (jenis === 'teksPanjang') return 'komiteclaimfacin__panjang'
  return undefined
}

/** Satu sel grid: tautan "View Retro" (pop-up), kotak centang, atau nilai terformat. */
function Sel({
  v,
  jenis,
  label,
  modal,
  onModal,
}: {
  v: string
  jenis: KolomGrid['jenis']
  label: string
  modal: Layar['modal']
  onModal: BukaModal
}) {
  if (jenis === 'tautan') {
    const k = kunciModal(v, modal)
    if (k === null) return null
    return (
      <button
        type="button"
        className="komiteclaimfacin__tautan"
        onClick={(e) => {
          e.stopPropagation() // tautan di baris berincian tidak ikut membuka / menutup panel baris
          onModal(k)
        }}
      >
        {KCFI.lihatRetro}
      </button>
    )
  }
  if (jenis === 'centang') return <Centang v={v} label={label} />
  return <>{tampil(v, jenis)}</>
}

/** Grid hanya-baca; baris ber-`rincian` dapat dibuka (expand pane, satu baris terbuka), isinya bagian bersarang. */
function TabelGrid({ g, modal, onModal }: { g: Grid; modal: Layar['modal']; onModal: BukaModal }) {
  const [buka, setBuka] = useState<number | null>(null)
  const baris = barisGrid(g)
  const panah = adaRincian(g)
  const lebar = g.kolom.length + (panah ? 1 : 0)
  return (
    <div className="komiteclaimfacin__grid">
      {g.judul && <h4 className="komiteclaimfacin__subjudul">{g.judul}</h4>}
      <div className="komiteclaimfacin__gulir">
        <table className="komiteclaimfacin__tabel">
          <thead>
            <tr>
              {panah && <th className="komiteclaimfacin__panah-kolom" aria-label={TATA_KCFI.rincian} />}
              {g.kolom.map((k, i) => (
                <th key={i} className={k.jenis === 'angka' ? 'komiteclaimfacin__angka' : undefined}>
                  {k.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {baris.length === 0 && (
              <tr>
                <td colSpan={lebar} className="komiteclaimfacin__kosong">
                  {TATA_KCFI.tanpaData}
                </td>
              </tr>
            )}
            {baris.map((b, i) => {
              const r = rincianBaris(g, i)
              const terbuka = r !== null && buka === i
              const alih = () => setBuka(terbuka ? null : i)
              return (
                <Fragment key={i}>
                  <tr className={r ? 'komiteclaimfacin__baris-buka' : undefined} onClick={r ? alih : undefined}>
                    {panah && (
                      <td className="komiteclaimfacin__panah-kolom">
                        {r && (
                          <button
                            type="button"
                            className="komiteclaimfacin__panah"
                            aria-expanded={terbuka}
                            aria-label={TATA_KCFI.rincian}
                            onClick={(e) => {
                              e.stopPropagation()
                              alih()
                            }}
                          >
                            {terbuka ? '▾' : '▸'}
                          </button>
                        )}
                      </td>
                    )}
                    {g.kolom.map((k, j) => (
                      <td key={j} className={kelasSel(k.jenis)}>
                        <Sel v={b[k.properti] ?? ''} jenis={k.jenis} label={k.label} modal={modal} onModal={onModal} />
                      </td>
                    ))}
                  </tr>
                  {terbuka && r && (
                    <tr className="komiteclaimfacin__baris-rinci">
                      <td colSpan={lebar}>
                        <div className="komiteclaimfacin__rinci">
                          {r.map((x, n) => (
                            <BagianIsi
                              key={`${x.kunci}-${String(n)}`}
                              b={x}
                              judulKartu=""
                              modal={modal}
                              onModal={onModal}
                            />
                          ))}
                        </div>
                      </td>
                    </tr>
                  )}
                </Fragment>
              )
            })}
          </tbody>
        </table>
      </div>
      {(g.kaki ?? []).length > 0 && <KakiGrid kaki={g.kaki ?? []} />}
    </div>
  )
}

/** Isi satu bagian (REKURSIF lewat rincian grid): sub-judul, pasangan label-nilai, grid. */
function BagianIsi({
  b,
  judulKartu,
  modal,
  onModal,
}: {
  b: Bagian
  judulKartu: string
  modal: Layar['modal']
  onModal: BukaModal
}) {
  return (
    <div className="komiteclaimfacin__bagian">
      {b.judul && b.judul !== judulKartu && <h4 className="komiteclaimfacin__subjudul">{b.judul}</h4>}
      {b.medan && b.medan.length > 0 && <DaftarNilai medan={b.medan} />}
      {(b.grid ?? []).map((g, i) => (
        <TabelGrid key={i} g={g} modal={modal} onModal={onModal} />
      ))}
    </div>
  )
}

/** Satu kartu = satu bagian server (judul server; kosong = kartu tanpa judul). */
function Kartu({ b, modal, onModal }: { b: Bagian; modal: Layar['modal']; onModal: BukaModal }) {
  const judul = b.judul ?? ''
  return (
    <section className="panel komiteclaimfacin__kartu" aria-label={judul || undefined}>
      {judul && <h3 className="panel__title komiteclaimfacin__kartu-judul">{judul}</h3>}
      <BagianIsi b={b} judulKartu={judul} modal={modal} onModal={onModal} />
    </section>
  )
}

const TANDA_LANGKAH: Record<KeadaanLangkah, string> = { setuju: '✓', tolak: '✕', berjalan: '', menunggu: '' }

export default function KasusKomite({
  id,
  onKembali,
  onLihatBerkas,
}: {
  id: string
  /** Kembali ke inbox Claim Fac In; `pesan` = info Submit yang tampil dulu. */
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
  const [modalBuka, setModalBuka] = useState<string | null>(null)
  // 403 / 409 sesudah Submit: kasus berpindah tingkat / tertutup - isian dikunci, kembali lewat Cancel.
  const [basi, setBasi] = useState(false)
  // Submit yang selesai sesudah layar ditinggal (menu lain dibuka) tidak boleh menarik pengguna kembali ke inbox.
  const terpasang = useRef(true)
  useEffect(() => {
    terpasang.current = true // StrictMode: efek dipasang ulang sesudah cleanup
    return () => {
      terpasang.current = false
    }
  }, [])

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
  if (layar === null || isi === null) return <Memuat pesan={KCFI.memuat} />

  const { isian } = layar
  const kunci = !layar.bolehKerja || sibuk || basi
  const usulKunci = kunci || !usulTerbuka(isian)

  const kirim = () => {
    const p = periksa(isi)
    setSalah(p)
    if (Object.keys(p).length > 0) return
    setSibuk(true)
    setGalat(null)
    putuskan(id, isianKirim(isi, isian)).then(
      (h) => {
        if (terpasang.current) onKembali(h.info ? [h.info] : undefined)
      },
      (g: unknown) => {
        if (!terpasang.current) return
        setSibuk(false)
        setGalat(g)
        if (g instanceof ApiFailure && (g.status === 403 || g.status === 409)) setBasi(true)
      },
    )
  }

  const { kartu, tangga } = susunBagian(layar.bagian ?? [])
  const langkah = langkahTangga(tangga, layar.kasus.tangga)
  const judulTangga = tangga?.judul || tangga?.grid?.[0]?.judul || ''
  const k = keadaanKasus(layar)
  const status =
    k.jenis === 'anda'
      ? { teks: TATA_KCFI.statusGiliranAnda, kelas: 'komiteclaimfacin__status--anda' }
      : k.jenis === 'tunggu'
        ? { teks: `${TATA_KCFI.statusMenunggu} ${k.jabatan}`, kelas: 'komiteclaimfacin__status--tunggu' }
        : { teks: TATA_KCFI.statusSelesai, kelas: 'komiteclaimfacin__status--selesai' }
  const tombolLihat = layar.tombol.find((t) => t.aksi === 'lihat')
  const isiModal = modalBuka !== null ? layar.modal?.[modalBuka] : undefined
  const judulModal = isiModal?.[0]?.judul || KCFI.lihatRetro

  // Galat Submit: 422 = pesan validasi server; 403 / 409 = kalimat server; selainnya komponen galat inti.
  const pesanValidasi = galat instanceof GalatValidasiKomite ? galat.pesan : null
  const pesanGalat = galat instanceof ApiFailure ? (galat.detail.message ?? null) : null
  const galatLain = galat !== null && pesanValidasi === null && pesanGalat === null

  return (
    <section className="inbox komiteclaimfacin__akar">
      <header className="komiteclaimfacin__kepala">
        <div className="komiteclaimfacin__kepala-baris">
          <div>
            <h2 className="inbox__judul komiteclaimfacin__judul">{layar.judul.join(' ')}</h2>
            <div className="komiteclaimfacin__kepala-sub">
              <span className={'komiteclaimfacin__status ' + status.kelas}>{status.teks}</span>
            </div>
          </div>
          {tombolLihat && (
            // View more details = klaim induk: berkas Claim Fac In dibuka hanya-baca di jendela di atas layar ini
            // (`PropsRute.onLihatBerkas`). Modul tidak dipasang = pesan.
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={!tombolLihat.aktif}
              title={tombolLihat.alasan}
              onClick={() => {
                setPesanLihat(
                  onLihatBerkas?.('claimfacin', layar.kasus.klaimId) === true ? null : KCFI.berkasTakTerpasang,
                )
              }}
            >
              {tombolLihat.label}
            </button>
          )}
        </div>
        {layar.ubin.length > 0 && (
          <dl className="komiteclaimfacin__ringkas">
            {layar.ubin.map((m, i) => {
              const v = tampil(m.nilai, m.jenis)
              return (
                <div key={i} className="komiteclaimfacin__ringkas-isi">
                  <dt>{m.label}</dt>
                  <dd className={m.jenis === 'angka' ? 'komiteclaimfacin__angka-teks' : undefined}>
                    {v === '' ? <span className="komiteclaimfacin__kosong">—</span> : v}
                  </dd>
                </div>
              )
            })}
          </dl>
        )}
      </header>
      {pesanLihat && <div className="alert alert--info">{pesanLihat}</div>}
      {(layar.pesan ?? []).map((p) => (
        <div key={p} className="alert alert--info">
          {p}
        </div>
      ))}

      <div className="komiteclaimfacin__tata">
        <div className="komiteclaimfacin__utama">
          {kartu.map((b, i) => (
            <Kartu key={`${b.kunci}-${String(i)}`} b={b} modal={layar.modal} onModal={setModalBuka} />
          ))}
        </div>

        <div className="komiteclaimfacin__bawah">
          {langkah.length > 0 && (
            <section className="panel komiteclaimfacin__kartu" aria-label={judulTangga || undefined}>
              {judulTangga && <h3 className="panel__title komiteclaimfacin__kartu-judul">{judulTangga}</h3>}
              <ol className="komiteclaimfacin__langkah">
                {langkah.map((l, i) => (
                  <li key={i} className={'komiteclaimfacin__langkah-isi komiteclaimfacin__langkah--' + l.keadaan}>
                    <span className="komiteclaimfacin__langkah-tanda" aria-hidden="true">
                      {TANDA_LANGKAH[l.keadaan] || l.no}
                    </span>
                    <div className="komiteclaimfacin__langkah-teks">
                      <div className="komiteclaimfacin__langkah-kepala">
                        <strong>{l.jabatan}</strong>
                        <span className={'komiteclaimfacin__status komiteclaimfacin__status--' + l.keadaan}>
                          {l.status}
                        </span>
                      </div>
                      {l.tanggal && (
                        <span className="komiteclaimfacin__langkah-tgl">{tampil(l.tanggal, 'tanggalJam')}</span>
                      )}
                      {l.komentar && <p className="komiteclaimfacin__langkah-catatan">{l.komentar}</p>}
                    </div>
                  </li>
                ))}
              </ol>
            </section>
          )}

          <section
            className="panel komiteclaimfacin__kartu komiteclaimfacin__keputusan"
            aria-label={TATA_KCFI.keputusanAnda}
          >
            <h3 className="panel__title komiteclaimfacin__kartu-judul">{TATA_KCFI.keputusanAnda}</h3>
            {!layar.bolehKerja && <div className="alert alert--info">{KCFI.hanyaLihat}</div>}

            <div className="komiteclaimfacin__isian-blok">
              <span className="komiteclaimfacin__label" id="komiteclaimfacin-terima">
                {isian.label.acceptStatus} <span className="field__req">*</span>
              </span>
              <div className="komiteclaimfacin__pilihan" role="radiogroup" aria-labelledby="komiteclaimfacin-terima">
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
                        'komiteclaimfacin__opsi komiteclaimfacin__opsi--' +
                        jenis +
                        (dipilih ? ' komiteclaimfacin__opsi--dipilih' : '')
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

            {isian.tampilUsul && (
              <>
                <label className="komiteclaimfacin__centang">
                  <input
                    type="checkbox"
                    checked={isi.usulTutup}
                    disabled={usulKunci}
                    onChange={(e) => ubah('usulTutup', e.target.checked)}
                  />
                  {isian.label.usulTutup}
                </label>
                <label className="komiteclaimfacin__centang">
                  <input
                    type="checkbox"
                    checked={isi.usulCadang}
                    disabled={usulKunci}
                    onChange={(e) => ubah('usulCadang', e.target.checked)}
                  />
                  {isian.label.usulCadang}
                </label>
              </>
            )}
            <label className="komiteclaimfacin__isian-blok">
              <span className="komiteclaimfacin__label">
                {isian.label.comment} <span className="field__req">*</span>
              </span>
              <textarea
                className="field__input komiteclaimfacin__area"
                value={isi.comment}
                disabled={kunci}
                onChange={(e) => ubah('comment', e.target.value)}
              />
              {salah.comment && <span className="field__error">{salah.comment}</span>}
            </label>
            {pesanValidasi && pesanValidasi.length > 0 && (
              <div className="alert alert--error" role="alert">
                <ul className="komiteclaimfacin__pesan">
                  {pesanValidasi.map((p) => (
                    <li key={p}>{p}</li>
                  ))}
                </ul>
              </div>
            )}
            {pesanGalat && (
              <div className="alert alert--error" role="alert">
                {pesanGalat}
              </div>
            )}
            {galatLain && <Gagal galat={galat} />}
            <div className="komiteclaimfacin__aksi">
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
                      disabled={!t.aktif || sibuk || basi}
                      title={t.alasan}
                      onClick={kirim}
                    >
                      {sibuk ? KCFI.memproses : t.label}
                    </button>
                  ),
                )}
            </div>
          </section>
        </div>
      </div>

      {isiModal && (
        // Local action ShowRetro ("View Retro"): reasuradur Fac Retro + Security Reinsurer per baris (expand pane).
        <Modal judul={judulModal} onTutup={() => setModalBuka(null)} labelBatal={KCFI.tutup} lebar>
          <div className="komiteclaimfacin__modal">
            {isiModal.map((b, i) => (
              <BagianIsi
                key={`${b.kunci}-${String(i)}`}
                b={b}
                judulKartu={judulModal}
                modal={layar.modal}
                onModal={setModalBuka}
              />
            ))}
          </div>
        </Modal>
      )}
    </section>
  )
}
