// Layar satu kasus Claim Non Prop - FlowAction OutstandingClaim (Assignment2) atau InputAcceptation (Assignment1). Isinya
// pohon tata server; aksi dikirim bersama nilai semua medan terbuka. Local action / harness (GeneratePLACNP,
// CloseClaimMD, CloseClaimNP, KomiteCNP) dibuka sebagai Modal; expand pane (AdjustmentDetailNP, InputDtlInterest,
// ShowDetailXOL) di bawah barisnya; pop-up pemilih / tampilan sebagai `Popup`. Pola
// `modul/claimprop/frontend/components/LayarKasus.tsx` (disalin, bukan impor).

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
import { CNP } from '../labels'
import { ambil, masukan, semuaTata, setel } from '../nilai'
import PanelLampiran from './PanelLampiran'
import Popup, { type JenisPopup, type SumberMaster } from './Popup'
import { barisModal, DAFTAR_ADJ, DAFTAR_INTEREST, jalurXOL, type RincianGrid } from './rincian'
import { pisahKakiModal } from './susun'
import TataView, { BarisAdjustment, LayarTata, Tombol, type KonteksTata } from './TataView'

/** Aksi yang hanya membuka pop-up harness (tanpa activity server sebelum pop-up terbuka). */
const POPUP: Record<string, { jenis: JenisPopup; sumber?: SumberMaster }> = {
  'PilihMaster:IN': { jenis: 'master', sumber: 'IN' },
  'PilihMaster:INEDM': { jenis: 'master', sumber: 'INEDM' },
  PilihPolis: { jenis: 'daftarPolis' },
  PilihSebab: { jenis: 'sebab' },
  BukaKatastrofe: { jenis: 'katastrofe' },
  LihatSelisihAktual: { jenis: 'selisihAktual' },
  LihatAlokasiLama: { jenis: 'alokasiLama' },
  LihatLampiranBayar: { jenis: 'lampiranBayar' },
  LihatRiwayatMaster: { jenis: 'riwayatMaster' },
}

/** Tombol aktif beraksi `aksi` di layout utama (tab Lampiran: Save = tombol Save layar, pola Claim Prop). */
function adaTombol(ts: readonly Tata[], aksi: string): boolean {
  return ts.some((t) => (t.jenis === 'tombol' && t.aksi === aksi && !t.nonaktif) || adaTombol(t.anak ?? [], aksi))
}

/** Aksi tombol Save Outstanding Claim (`SaveDataToJClaim_Act`); Input Acceptation tanpa Save layar (`1=2`). */
const AKSI_SIMPAN = 'SaveDataToJClaim'

/** Aksi yang parameternya nilai pilihan (bukan jalur medan): rekening = `AccountNo|NameOfBank`. */
const PARAM_DARI_NILAI = new Set(['PilihRekening', 'PilihRekening2'])

const JUDUL_MODAL = (m: string): string => {
  if (m === 'pla') return CNP.popPLA
  if (m === 'tutupKlaim') return CNP.popTutup
  if (m === 'cwp') return CNP.popCWP
  if (m.startsWith('komite:')) return CNP.popKomite
  return m
}

/** Satu aksi layar beserta ubahan medan yang memicunya. */
interface Permintaan extends AksiDiminta {
  baris: number
  ubahan: Record<string, string>
  param: string
  /** Tombol kaki modal: modal ditutup bila aksinya berhasil tanpa pesan. */
  tutupModal: boolean
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
  /** Tampilan saja walau pemegang. */
  hanyaLihat?: boolean
}) {
  const [layar, setLayar] = useState<Layar | null>(null)
  const [h, setH] = useState<Halaman | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [info, setInfo] = useState<string | null>(null)
  const [modal, setModal] = useState<string | null>(null)
  const [popup, setPopup] = useState<{ jenis: JenisPopup; sumber?: SumberMaster } | null>(null)
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

  // Treaty Type Spreading List (ListTreaty.pxResults) bergantung pada master klaim; dibaca ulang saat master berganti.
  const idMaster = h?.nilai['ClaimData.IDMaster'] ?? ''
  useEffect(() => {
    if (idMaster === '') return
    pilihanKasus<Pilihan[]>(id, 'treaty').then(
      (d) => setOpsiKasus((o) => ({ ...o, treaty: d ?? [] })),
      () => undefined,
    )
  }, [id, idMaster])
  // Specify (Payable To = Others) = dropdown ClientName.pxResults - dibaca saat dibutuhkan.
  const payable = h?.nilai['ClaimData.Payable'] ?? ''
  useEffect(() => {
    if (payable !== '3') return
    pilihanKasus<Pilihan[]>(id, 'klien').then(
      (d) => setOpsiKasus((o) => ({ ...o, klien: d ?? [] })),
      () => undefined,
    )
  }, [id, payable])

  // Aksi yang sedang berjalan dan satu aksi yang menunggu (`antreAksi.ts`).
  const berjalan = useRef<Permintaan | null>(null)
  const menunggu = useRef<Permintaan | null>(null)

  const jalankan = useCallback(
    function jalan(p: Permintaan, ly: Layar, hKini: Halaman) {
      const { aksi, indeks, baris, ubahan, param, tutupModal } = p
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
        baris,
        param: prm,
        tahap: ly.kasus.tahap,
        masukan: masukan(h2, semua),
        mode: ly.mode,
      }).then(
        (l) => {
          terima(l)
          if (l.info) setInfo(l.info)
          if (tutupModal && (l.pesan ?? []).length === 0) setModal(null)
          if (l.bukaModal) setModal(l.bukaModal)
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
    (aksi: string, indeks = 0, ubahan: Record<string, string> = {}, param = '', baris = 0, tutupModal = false) => {
      if (!layar || !h) return
      const p: Permintaan = { aksi, indeks, baris, ubahan, param, tutupModal }
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

  // Tombol View polis (bawaan OQ-CNP-13 = jendela Modal berkas NB / EDM Treaty In, pola Claim Prop): berkas polis ini
  // dicari (`GET /berkas-polis`) lalu tampil di atas layar ini lewat `PropsRute.onLihatBerkas`.
  const noPolis = h ? ambil(h, 'ClaimData.PolicyData.PolicyNo') : ''
  const [pesanView, setPesanView] = useState<string | null>(null)
  const lihatPolis = useCallback(() => {
    if (noPolis === '') return
    berkasPolis(noPolis).then(
      (b) => {
        if (onLihatBerkas?.(b.modul, b.kasus) !== true) setPesanView(CNP.berkasTakTerpasang)
      },
      (g: unknown) => {
        if (g instanceof ApiFailure && g.status === 404) setPesanView(CNP.polisTanpaBerkas)
        else setGalat(g)
      },
    )
  }, [noPolis, onLihatBerkas])

  const aksiLayar = useCallback(
    (nama: string, indeks = 0, ubahan: Record<string, string> = {}, baris = 0, tutupModal = false) => {
      const pop = POPUP[nama]
      if (pop) {
        setPopup(pop)
        return
      }
      if (nama === 'BukaPLA') return setModal('pla')
      if (nama === 'BukaTutupKlaim') return setModal('tutupKlaim')
      if (nama === 'TutupModal') return setModal(null)
      if (nama === 'LihatPolis') return lihatPolis()
      kirim(nama, indeks, ubahan, '', baris, tutupModal)
    },
    [kirim, lihatPolis],
  )

  const saran = useCallback(
    (sumber: string, indeks: number, cari: string) => {
      if (!['polis', 'adjuster', 'provinsi', 'rekening'].includes(sumber)) return
      pilihanKasus<Pilihan[]>(id, sumber, indeks, cari).then(
        (d) => setSaranPeta((s) => ({ ...s, [`${sumber}|${indeks}`]: d ?? [] })),
        () => undefined,
      )
    },
    [id],
  )

  const opsi = useCallback(
    (sumber: string, indeks: number): Pilihan[] => {
      if (sumber === 'mataUang') return acuan?.mataUang ?? []
      if (sumber === 'lossAlloc') return acuan?.lossAlloc ?? []
      if (sumber.startsWith('kode:')) {
        const p = sumber.slice(5)
        return (acuan?.kode[p] ?? []).map((k) => ({ nilai: k, label: acuan?.labelKode?.[p]?.[k] ?? k }))
      }
      if (opsiKasus[sumber]) return opsiKasus[sumber] ?? []
      return saranPeta[`${sumber}|${indeks}`] ?? []
    },
    [acuan, opsiKasus, saranPeta],
  )

  const nAkseptasi = h?.daftar[DAFTAR_ADJ]?.length ?? 0
  const k: KonteksTata | null = useMemo(() => {
    if (!layar || !h) return null
    const kk: KonteksTata = {
      h,
      ubah: (j, v) => setH((x) => (x ? setel(x, j, v) : x)),
      aksi: (nama, indeks, ubahan, baris) => aksiLayar(nama, indeks, ubahan, baris),
      opsi,
      saran,
      pesanMedan: layar.pesanMedan ?? {},
      sibuk,
    }
    const panel = (tata: Layar['tata'] | undefined) => (
      <div className="claimnonprop__rinci">
        <TataView tata={tata ?? []} k={kk} />
      </div>
    )
    const rincian: RincianGrid[] = [
      // Acceptation List -> AdjustmentDetailNP (pyExpandMultipleRows); baris terbaru terbuka (pola Claim Prop).
      { daftar: DAFTAR_ADJ, bukaAwal: true, nomorAkseptasi: true, isi: (n) => panel(layar.adjustment?.[String(n)]) },
      // Insured Interests 100 % Outstanding -> InputDtlInterest.
      { daftar: DAFTAR_INTEREST, isi: (i) => panel(layar.modal?.[`interest:${i}`]) },
    ]
    // XOL Allocation panel akseptasi ke-n -> ShowDetailXOL (ReinstatementPremiumDetails).
    for (let n = 1; n <= nAkseptasi; n++) {
      rincian.push({ daftar: jalurXOL(n), isi: (i) => panel(layar.modal?.[`reinstatement:${n}:${i}`]) })
    }
    kk.rincian = rincian
    return kk
  }, [layar, h, aksiLayar, opsi, saran, sibuk, nAkseptasi])

  // Pesan halaman tampil di atas layar; aksi dari bagian bawah menggulir ke kotak pesan supaya penolakannya terlihat.
  const kotakPesan = useRef<HTMLDivElement>(null)
  const pesanLayar = layar?.pesan
  useEffect(() => {
    if ((pesanLayar ?? []).length > 0) kotakPesan.current?.scrollIntoView({ block: 'nearest', behavior: 'smooth' })
  }, [pesanLayar])

  if (!layar || !h || !k) {
    return (
      <section className="inbox claimnonprop__akar">
        {galat ? <Gagal galat={galat} /> : <Memuat pesan={CNP.memuat} />}
      </section>
    )
  }

  const pesanGalat = galat instanceof ApiFailure ? galat.detail.message : null
  // Tombol di akhir isi modal ke kaki `Modal`; tombol penutup section = tombol batal `Modal` (bukan dua Cancel).
  const isiModal = modal ? pisahKakiModal(layar.modal?.[modal] ?? []) : null
  const kModal: KonteksTata = { ...k, rincian: undefined }
  const kKaki: KonteksTata = {
    ...kModal,
    aksi: (nama, indeks, ubahan, baris) => aksiLayar(nama, indeks, ubahan, baris, true),
  }
  return (
    <section className="inbox claimnonprop__akar">
      <header className="inbox__kepala">
        <button type="button" className="btn btn--ghost btn--sm" onClick={onKembali}>
          {CNP.kembali}
        </button>
        <h2 className="inbox__judul">
          {layar.kasus.id} - {layar.label || layar.kasus.statusWork}
        </h2>
      </header>
      {!layar.bolehKerja && <div className="alert">{CNP.hanyaLihat}</div>}
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
        <Modal judul={CNP.view} onTutup={() => setPesanView(null)} labelBatal={CNP.tutup}>
          <p role="alert">{pesanView}</p>
        </Modal>
      )}
      <LayarTata tata={layar.tata} k={k} />
      <PanelLampiran
        id={id}
        hanyaLihat={hanyaLihat}
        bolehSimpan={layar.bolehKerja && adaTombol(layar.tata, AKSI_SIMPAN)}
        onSimpan={() => kirim(AKSI_SIMPAN)}
      />
      {modal && isiModal && (
        <BarisAdjustment.Provider value={barisModal(modal)}>
          <Modal
            judul={JUDUL_MODAL(modal)}
            onTutup={() => setModal(null)}
            lebar
            labelBatal={isiModal.batal}
            aksi={isiModal.kaki.map((t, i) => (
              <Tombol key={i} t={t} k={kKaki} utama />
            ))}
          >
            <TataView tata={isiModal.isi} k={kModal} />
          </Modal>
        </BarisAdjustment.Provider>
      )}
      {popup && (
        // Satu Popup per jenis: data jenis lama tidak pernah dirender dengan bentuk jenis baru.
        <Popup
          key={`${popup.jenis}:${popup.sumber ?? ''}`}
          id={id}
          jenis={popup.jenis}
          sumber={popup.sumber}
          pelaku={pelaku}
          sts={{ sts: ambil(h, 'ClaimData.StsKatastrofe'), non: ambil(h, 'ClaimData.NonKatastrofeType') }}
          mataUang={acuan?.mataUang ?? []}
          onPilih={(a, p) => kirim(a, 0, {}, p)}
          onTutup={() => setPopup(null)}
        />
      )}
    </section>
  )
}
