// Layar satu kasus Claim Prop - FlowAction OutstandingClaim (Assignment2) atau InputAcceptation (Assignment1). Isinya
// pohon tata server; aksi dikirim bersama nilai semua medan terbuka. Local action / harness (PrintFile, GeneratePLA,
// GenerateDLATreaty, PreventRejectClaimProp, CommitteeTreaty) dibuka sebagai Modal; pop-up pemilih sebagai `Popup`.

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ApiFailure } from '../../../../inti/frontend/klien'
import {
  aksiKasus,
  ambilAcuan,
  berkasPolis,
  bukaKasus,
  pilihanKasus,
  type AcuanStatis,
  type Halaman,
  type Layar,
  type Pilihan,
  type Tata,
} from '../api'
import { putuskanAksi, type AksiDiminta } from '../antreAksi'
import { CP } from '../labels'
import { ambil, masukan, semuaTata, setel } from '../nilai'
import PanelLampiran from './PanelLampiran'
import Popup, { type JenisPopup } from './Popup'
import { barisModal, DAFTAR_RINCI } from './rincian'
import { pisahKakiModal } from './susun'
import TambahAdjuster from './TambahAdjuster'
import TataView, { BarisAdjustment, LayarTata, Tombol, type KonteksTata } from './TataView'

/** Aksi yang hanya membuka pop-up harness. */
const POPUP: Record<string, JenisPopup> = {
  PilihMaster: 'master',
  PilihPolis: 'polis',
  PilihSebab: 'sebab',
  BukaKatastrofe: 'katastrofe',
  RingkasanOS: 'ringkasanOS',
}

/** Tombol aktif beraksi `aksi` di layout utama (tab Lampiran: Save = tombol Save layar, aksi `Simpan`). */
function adaTombol(ts: readonly Tata[], aksi: string): boolean {
  return ts.some((t) => (t.jenis === 'tombol' && t.aksi === aksi && !t.nonaktif) || adaTombol(t.anak ?? [], aksi))
}

/** Aksi submit local action: modal ditutup bila berhasil. */
const SUBMIT_MODAL = new Set(['TryMakePLA', 'PrintDLATreatyIn', 'CloseClaimProp', 'AddKomiteTreatyChild'])

/** Aksi yang parameternya nilai pilihan (bukan jalur medan). */
const PARAM_DARI_NILAI = new Set(['PilihRekening', 'PilihAllocation', 'DisableEditRNMShare'])

const JUDUL_MODAL = (m: string): string => {
  if (m === 'pla') return CP.popPLA
  if (m === 'tutupKlaim') return CP.popTutup
  if (m.startsWith('dla:')) return CP.popDLA
  if (m.startsWith('komite:')) return CP.popKomite
  return m
}

/** Satu aksi layar beserta ubahan medan yang memicunya. */
interface Permintaan extends AksiDiminta {
  ubahan: Record<string, string>
  param: string
}

interface RekeningBank {
  ClientName: string
  NameOfBank: string
  BranchOfBank: string
  AccountNo: string
}

export default function LayarKasus({
  id,
  pelaku,
  onKembali,
  onLihatBerkas,
  hanyaLihat = false,
}: {
  id: string
  pelaku: string
  onKembali: () => void
  /** `PropsRute.onLihatBerkas` - tombol View polis. */
  onLihatBerkas?: (modul: string, id: string) => boolean
  /** Tampilan saja walau pemegang (View more details Komite Claim Prop, keputusan work owner 09-10-2026). */
  hanyaLihat?: boolean
}) {
  const [layar, setLayar] = useState<Layar | null>(null)
  const [h, setH] = useState<Halaman | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [info, setInfo] = useState<string | null>(null)
  const [modal, setModal] = useState<string | null>(null)
  const [popup, setPopup] = useState<JenisPopup | null>(null)
  // Tombol "+" Consultant / Adjuster: medan dan aksi yang diisi ID master baru.
  const [tambahAdj, setTambahAdj] = useState<{ jalur: string; aksi: string } | null>(null)
  const [acuan, setAcuan] = useState<AcuanStatis | null>(null)
  const [opsiKasus, setOpsiKasus] = useState<Record<string, Pilihan[]>>({})
  const [saranPeta, setSaranPeta] = useState<Record<string, Pilihan[]>>({})

  const terima = useCallback((l: Layar) => {
    setLayar(l)
    setH(l.halaman)
    setGalat(null)
  }, [])

  useEffect(() => {
    bukaKasus(id, hanyaLihat).then(terima, (g: unknown) => setGalat(g))
    ambilAcuan().then(setAcuan, (g: unknown) => setGalat(g))
  }, [id, hanyaLihat, terima])

  const idMaster = h?.nilai['ClaimData.IDMaster'] ?? ''
  useEffect(() => {
    if (idMaster === '') return
    for (const s of ['limits', 'shareRNM']) {
      pilihanKasus<Pilihan[]>(id, s).then(
        (d) => setOpsiKasus((o) => ({ ...o, [s]: d ?? [] })),
        () => undefined,
      )
    }
  }, [id, idMaster])
  // Treaty Type spreading = spreading polis klaim (koreksi 08-10-2026), jadi dibaca ulang saat polis berganti.
  const polisKasus = h?.nilai['ClaimData.PolicyData.PolicyNo'] ?? ''
  useEffect(() => {
    if (polisKasus === '') return
    pilihanKasus<Pilihan[]>(id, 'spreading').then(
      (d) => setOpsiKasus((o) => ({ ...o, spreading: d ?? [] })),
      () => undefined,
    )
  }, [id, polisKasus])

  // Aksi yang sedang berjalan dan satu aksi yang menunggu (`antreAksi.ts`).
  const berjalan = useRef<Permintaan | null>(null)
  const menunggu = useRef<Permintaan | null>(null)

  const jalankan = useCallback(
    function jalan(p: Permintaan, ly: Layar, hKini: Halaman) {
      const { aksi, indeks, ubahan, param } = p
      let h2 = hKini
      for (const [j, v] of Object.entries(ubahan)) h2 = setel(h2, j, v)
      setH(h2)
      const prm = param !== '' ? param : PARAM_DARI_NILAI.has(aksi) ? (Object.values(ubahan)[0] ?? '') : ''
      const semua = semuaTata(ly.tata, [ly.adjustment ?? {}, ly.modal ?? {}])
      berjalan.current = p
      setSibuk(true)
      setInfo(null)
      aksiKasus(id, {
        aksi,
        indeks,
        param: prm,
        tahap: ly.kasus.tahap,
        masukan: masukan(h2, semua),
        mode: ly.mode,
      }).then(
        (l) => {
          terima(l)
          if (l.info) setInfo(l.info)
          if (SUBMIT_MODAL.has(aksi)) setModal(null)
          if (aksi === 'BukaKomite') {
            const ok = l.halaman.nilai['Protect.CARI1'] === '1' && l.halaman.nilai['Protect.CARI2'] === '1'
            setModal(ok ? `komite:${indeks}` : null)
          }
          setPopup(null)
          berjalan.current = null
          const lanjut = menunggu.current
          menunggu.current = null
          if (lanjut) jalan(lanjut, l, l.halaman)
          else setSibuk(false)
        },
        (g: unknown) => {
          berjalan.current = null
          menunggu.current = null
          setSibuk(false)
          setGalat(g)
        },
      )
    },
    [id, terima],
  )

  const kirim = useCallback(
    (aksi: string, indeks = 0, ubahan: Record<string, string> = {}, param = '') => {
      if (!layar || !h) return
      const p: Permintaan = { aksi, indeks, ubahan, param }
      const putusan = putuskanAksi(berjalan.current, p)
      if (putusan === 'abaikan') return
      if (putusan === 'antre') {
        menunggu.current = p
        return
      }
      jalankan(p, layar, h)
    },
    [layar, h, jalankan],
  )

  // Tombol View polis (perintah work owner 08-10-2026 "jangan tab baru ... biarkan di layar utama", "hanya tampilan
  // polisnya aja"): berkas NB / EDM Treaty In polis ini dicari (`GET /berkas-polis`) lalu tampil di jendela di atas
  // layar ini lewat `PropsRute.onLihatBerkas` - tanpa menu, tanpa tab baru. Tanpa berkas / modul tidak dipasang = pesan.
  const noPolis = h ? ambil(h, 'ClaimData.PolicyData.PolicyNo') : ''
  const [pesanView, setPesanView] = useState<string | null>(null)
  const lihatPolis = useCallback(() => {
    if (noPolis === '') return
    berkasPolis(noPolis).then(
      (b) => {
        if (onLihatBerkas?.(b.modul, b.kasus) !== true) setPesanView(CP.berkasTakTerpasang)
      },
      (g: unknown) => {
        if (g instanceof ApiFailure && g.status === 404) setPesanView(CP.polisTanpaBerkas)
        else setGalat(g)
      },
    )
  }, [noPolis, onLihatBerkas])

  const aksi = useCallback(
    (nama: string, indeks = 0, ubahan: Record<string, string> = {}) => {
      if (POPUP[nama]) {
        setPopup(POPUP[nama] ?? null)
        return
      }
      if (nama === 'BukaPLA') return setModal('pla')
      if (nama === 'BukaTutupKlaim') return setModal('tutupKlaim')
      if (nama === 'BukaDLA') return setModal(`dla:${indeks}`)
      if (nama === 'TutupModal') return setModal(null)
      if (nama === 'LihatPolis') return lihatPolis()
      if (nama === 'TambahKonsultan') return setTambahAdj({ jalur: 'ClaimData.ConsultantID', aksi: 'SetConsultant' })
      if (nama === 'TambahAdjuster') return setTambahAdj({ jalur: 'ClaimData.AppointedADJID', aksi: 'SetAdjsuter' })
      kirim(nama, indeks, ubahan)
    },
    [kirim, lihatPolis],
  )

  const saran = useCallback(
    (sumber: string, indeks: number, cari: string) => {
      if (!['adjuster', 'provinsi', 'rekening'].includes(sumber)) return
      pilihanKasus<unknown[]>(id, sumber, indeks, cari).then(
        (d) => {
          const op: Pilihan[] =
            sumber === 'rekening'
              ? ((d ?? []) as RekeningBank[]).map((r) => ({
                  nilai: `${r.AccountNo}|${r.NameOfBank}`,
                  label: `${r.NameOfBank} - ${r.BranchOfBank} - ${r.AccountNo} (${r.ClientName})`,
                }))
              : ((d ?? []) as Pilihan[])
          setSaranPeta((s) => ({ ...s, [`${sumber}|${indeks}`]: op }))
        },
        () => undefined,
      )
    },
    [id],
  )

  const opsi = useCallback(
    (sumber: string, indeks: number): Pilihan[] => {
      if (!h) return []
      if (sumber === 'mataUang') return acuan?.mataUang ?? []
      if (sumber === 'jenisReas') return acuan?.jenisReas ?? []
      if (sumber === 'jenisReas4') return acuan?.jenisReas4 ?? []
      if (sumber.startsWith('kode:')) {
        const p = sumber.slice(5)
        return (acuan?.kode[p] ?? []).map((k) => ({ nilai: k, label: acuan?.labelKode?.[p]?.[k] ?? k }))
      }
      if (sumber === 'mataUangAdj') {
        return (h.daftar[`ClaimData.AdjustmentList(${indeks}).CurencyAdjustment`] ?? []).map((b) => ({
          nilai: b.CurrencyID ?? '',
          label: b.Currency ?? '',
        }))
      }
      if (sumber === 'allocation') {
        return (h.daftar[`ClaimData.AdjustmentList(${indeks}).LossAllocation`] ?? []).map((b) => ({
          nilai: b.TreatyName ?? '',
          label: `${b.TreatyName ?? ''} (${b.SharePercentage ?? ''})`,
        }))
      }
      if (opsiKasus[sumber]) return opsiKasus[sumber] ?? []
      return saranPeta[`${sumber}|${indeks}`] ?? []
    },
    [h, acuan, opsiKasus, saranPeta],
  )

  const k: KonteksTata | null = useMemo(() => {
    if (!layar || !h) return null
    const dasar: KonteksTata = {
      h,
      ubah: (j, v) => setH((x) => (x ? setel(x, j, v) : x)),
      aksi,
      opsi,
      saran,
      pesanMedan: layar.pesanMedan ?? {},
      sibuk,
    }
    return {
      ...dasar,
      rincian: {
        daftar: DAFTAR_RINCI,
        isi: (n: number) => (
          <div className="claimprop__rinci">
            <TataView tata={layar.adjustment?.[String(n)] ?? []} k={dasar} />
          </div>
        ),
      },
    }
  }, [layar, h, aksi, opsi, saran, sibuk])

  // Pesan halaman tampil di atas layar; aksi dari bagian bawah (Add Spreading List yang ditolak "spreading sama")
  // menggulir ke kotak pesan supaya penolakannya terlihat.
  const kotakPesan = useRef<HTMLDivElement>(null)
  const pesanLayar = layar?.pesan
  useEffect(() => {
    if ((pesanLayar ?? []).length > 0) kotakPesan.current?.scrollIntoView({ block: 'nearest', behavior: 'smooth' })
  }, [pesanLayar])

  if (!layar || !h || !k) {
    return (
      <section className="inbox claimprop__akar">
        {galat ? <Gagal galat={galat} /> : <Memuat pesan={CP.memuat} />}
      </section>
    )
  }

  const pesanGalat = galat instanceof ApiFailure ? galat.detail.message : null
  // Tombol di akhir isi modal ke kaki `Modal`; tombol penutup section = tombol batal `Modal` (bukan dua Cancel).
  const isiModal = modal ? pisahKakiModal(layar.modal?.[modal] ?? []) : null
  const kModal = { ...k, rincian: undefined }
  return (
    <section className="inbox claimprop__akar">
      <header className="inbox__kepala">
        <button type="button" className="btn btn--ghost btn--sm" onClick={onKembali}>
          {CP.kembali}
        </button>
        <h2 className="inbox__judul">
          {layar.kasus.id} - {layar.label || layar.kasus.statusWork}
        </h2>
      </header>
      {!layar.bolehKerja && <div className="alert">{CP.hanyaLihat}</div>}
      {galat !== null &&
        (pesanGalat ? <div className="alert alert--error">{pesanGalat}</div> : <Gagal galat={galat} />)}
      {(layar.pesan ?? []).length > 0 && (
        <div ref={kotakPesan} className="alert alert--error">
          <ul>
            {(layar.pesan ?? []).map((p, i) => (
              <li key={i}>{p}</li>
            ))}
          </ul>
        </div>
      )}
      {info && (
        <div className="alert alert--ok" role="status">
          {info}
        </div>
      )}
      {pesanView && (
        <Modal judul={CP.view} onTutup={() => setPesanView(null)} labelBatal={CP.tutup}>
          <p role="alert">{pesanView}</p>
        </Modal>
      )}
      <LayarTata tata={layar.tata} k={k} />
      <PanelLampiran
        id={id}
        hanyaLihat={hanyaLihat}
        bolehSimpan={layar.bolehKerja && adaTombol(layar.tata, 'Simpan')}
        onSimpan={() => kirim('Simpan')}
      />
      {modal && isiModal && (
        <BarisAdjustment.Provider value={barisModal(modal)}>
          <Modal
            judul={JUDUL_MODAL(modal)}
            onTutup={() => setModal(null)}
            lebar
            labelBatal={isiModal.batal}
            aksi={isiModal.kaki.map((t, i) => (
              <Tombol key={i} t={t} k={kModal} utama />
            ))}
          >
            <TataView tata={isiModal.isi} k={kModal} />
          </Modal>
        </BarisAdjustment.Provider>
      )}
      {tambahAdj && (
        <TambahAdjuster
          onTersimpan={(idBaru) => {
            const t = tambahAdj
            setTambahAdj(null)
            kirim(t.aksi, 0, { [t.jalur]: idBaru })
          }}
          onTutup={() => setTambahAdj(null)}
        />
      )}
      {popup && (
        <Popup
          id={id}
          jenis={popup}
          pelaku={pelaku}
          sts={{ sts: ambil(h, 'ClaimData.StsKatastrofe'), non: ambil(h, 'ClaimData.NonKatastrofeType') }}
          onPilih={(a, p) => kirim(a, 0, {}, p)}
          onTutup={() => setPopup(null)}
        />
      )}
    </section>
  )
}
